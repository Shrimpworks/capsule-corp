//go:build darwin

package installationowner

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"sync"

	"capsule.local/capsule/internal/execution/lifecyclestate"
	"golang.org/x/sys/unix"
)

const (
	stateRootOpenFlags = unix.O_RDONLY | unix.O_DIRECTORY | unix.O_NOFOLLOW | unix.O_CLOEXEC
	ownerLockOpenFlags = unix.O_RDONLY | unix.O_NOFOLLOW | unix.O_CLOEXEC
)

type darwinOperations interface {
	Open(string, int, uint32) (int, error)
	Openat(int, string, int, uint32) (int, error)
	Fstat(int, *unix.Stat_t) error
	FcntlInt(uintptr, int, int) (int, error)
	Flock(int, int) error
	Close(int) error
	EffectiveUID() uint32
	ReadRandom([]byte) (int, error)
}

type systemDarwinOperations struct{}

func (systemDarwinOperations) Open(path string, flags int, mode uint32) (int, error) {
	return unix.Open(path, flags, mode)
}

func (systemDarwinOperations) Openat(directoryFD int, name string, flags int, mode uint32) (int, error) {
	return unix.Openat(directoryFD, name, flags, mode)
}

func (systemDarwinOperations) Fstat(fd int, stat *unix.Stat_t) error {
	return unix.Fstat(fd, stat)
}

func (systemDarwinOperations) FcntlInt(fd uintptr, command int, argument int) (int, error) {
	return unix.FcntlInt(fd, command, argument)
}

func (systemDarwinOperations) Flock(fd int, operation int) error {
	return unix.Flock(fd, operation)
}

func (systemDarwinOperations) Close(fd int) error { return unix.Close(fd) }

func (systemDarwinOperations) EffectiveUID() uint32 {
	uid := unix.Geteuid()
	if uid < 0 {
		return 0
	}
	return uint32(uid) //nolint:gosec // Darwin UIDs are nonnegative and the negative sentinel was rejected above.
}

func (systemDarwinOperations) ReadRandom(destination []byte) (int, error) {
	return rand.Read(destination)
}

// TrustedStateRoot retains an already validated bootstrap directory descriptor
// and its closed lock entry name. It exposes neither the descriptor nor a path.
// The bootstrap path is consumed only by OpenDarwinTrustedStateRoot; owner
// acquisition itself is descriptor-relative.
type TrustedStateRoot struct {
	mu         sync.Mutex
	fd         int
	enrollment OwnerLockEnrollment
	operations darwinOperations
}

// OpenDarwinTrustedStateRoot resolves the bootstrap path supplied by a future
// installation composition once, validates its enrolled directory identity,
// and retains only the resulting descriptor and closed entry name. G1 does not
// authenticate that path or enrollment provenance. It never creates the root.
func OpenDarwinTrustedStateRoot(
	ctx context.Context,
	installationConfiguredPath string,
	enrollment OwnerLockEnrollment,
) (*TrustedStateRoot, error) {
	return openDarwinTrustedStateRoot(ctx, installationConfiguredPath, enrollment, systemDarwinOperations{})
}

func openDarwinTrustedStateRoot(
	ctx context.Context,
	installationConfiguredPath string,
	enrollment OwnerLockEnrollment,
	operations darwinOperations,
) (*TrustedStateRoot, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if installationConfiguredPath == "" || operations == nil {
		return nil, ErrInvalidEnrollment
	}
	if err := validateEnrollmentView(enrollment.view); err != nil {
		return nil, err
	}
	fd, err := operations.Open(installationConfiguredPath, stateRootOpenFlags, 0)
	if err != nil {
		return nil, repair("state-root-open")
	}
	if err := validateStateRootDescriptor(fd, enrollment, operations); err != nil {
		_ = operations.Close(fd)
		return nil, err
	}
	return &TrustedStateRoot{fd: fd, enrollment: enrollment, operations: operations}, nil
}

// Close releases only the bootstrap directory descriptor. It has no effect on
// an acquired owner-lock descriptor, which is separately lifetime-scoped.
func (root *TrustedStateRoot) Close() error {
	if root == nil {
		return ErrOwnerClosed
	}
	root.mu.Lock()
	defer root.mu.Unlock()
	if root.fd < 0 {
		return ErrOwnerClosed
	}
	if err := validateStateRootIdentity(root.fd, root.enrollment, root.operations); err != nil {
		root.fd = -1
		return invariant("state-root-close-identity")
	}
	fd := root.fd
	root.fd = -1
	flagsErr := validateDescriptorFlags(fd, root.operations, "state-root")
	if err := root.operations.Close(fd); err != nil {
		return localFailure("state-root-close")
	}
	if flagsErr != nil {
		return invariant("state-root-close-flags")
	}
	return nil
}

type darwinInstallationOwner struct {
	mu             sync.Mutex
	fd             int
	root           *TrustedStateRoot
	enrollment     OwnerLockEnrollment
	ownerSessionID lifecyclestate.OwnerSessionID
	operations     darwinOperations
}

func (*darwinInstallationOwner) sealedInstallationOwner() {}

// AcquireDarwinInstallationOwner opens the closed enrolled entry relative to
// the trusted state-root descriptor, validates, acquires nonblocking exclusive
// BSD flock, and revalidates before issuing a nonzero in-process owner session.
func AcquireDarwinInstallationOwner(
	ctx context.Context,
	root *TrustedStateRoot,
	enrollment OwnerLockEnrollment,
) (InstallationOwner, error) {
	if root == nil {
		return nil, ErrInvalidEnrollment
	}
	root.mu.Lock()
	defer root.mu.Unlock()
	return acquireDarwinInstallationOwner(ctx, root, enrollment)
}

func acquireDarwinInstallationOwner(
	ctx context.Context,
	root *TrustedStateRoot,
	enrollment OwnerLockEnrollment,
) (InstallationOwner, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := validateEnrollmentView(enrollment.view); err != nil {
		return nil, err
	}
	if root.fd < 0 || root.operations == nil {
		return nil, ErrOwnerClosed
	}
	if enrollment != root.enrollment {
		return nil, ErrInvalidEnrollment
	}
	if err := validateStateRootDescriptor(root.fd, enrollment, root.operations); err != nil {
		return nil, err
	}

	fd, err := root.operations.Openat(root.fd, enrollment.view.LockEntryName, ownerLockOpenFlags, 0)
	if err != nil {
		return nil, repair("owner-open")
	}
	closeOnFailure := true
	defer func() {
		if closeOnFailure {
			// A close failure leaves availability fail-closed until process exit.
			_ = root.operations.Close(fd)
		}
	}()

	if err := validateOwnerDescriptor(fd, enrollment, root.operations); err != nil {
		return nil, err
	}
	if err := root.operations.Flock(fd, unix.LOCK_EX|unix.LOCK_NB); err != nil {
		if errors.Is(err, unix.EWOULDBLOCK) || errors.Is(err, unix.EAGAIN) {
			return nil, ErrDuplicateOwner
		}
		return nil, localFailure("owner-flock")
	}
	if err := validateOwnerDescriptor(fd, enrollment, root.operations); err != nil {
		return nil, err
	}
	if err := validateStateRootDescriptor(root.fd, enrollment, root.operations); err != nil {
		return nil, err
	}
	if err := validateNamedOwnerEntry(root.fd, enrollment, root.operations); err != nil {
		return nil, err
	}

	ownerSessionID, err := newOwnerSessionID(root.operations)
	if err != nil {
		return nil, err
	}
	closeOnFailure = false
	return &darwinInstallationOwner{
		fd: fd, root: root, enrollment: enrollment, ownerSessionID: ownerSessionID, operations: root.operations,
	}, nil
}

func validateNamedOwnerEntry(
	stateRootFD int,
	enrollment OwnerLockEnrollment,
	operations darwinOperations,
) error {
	fd, err := operations.Openat(stateRootFD, enrollment.view.LockEntryName, ownerLockOpenFlags, 0)
	if err != nil {
		return repair("owner-reopen")
	}
	if err := validateOwnerDescriptor(fd, enrollment, operations); err != nil {
		_ = operations.Close(fd)
		return err
	}
	if err := operations.Close(fd); err != nil {
		return localFailure("owner-reopen-close")
	}
	return nil
}

func validateStateRootDescriptor(
	fd int,
	enrollment OwnerLockEnrollment,
	operations darwinOperations,
) error {
	if err := validateStateRootIdentity(fd, enrollment, operations); err != nil {
		return err
	}
	return validateDescriptorFlags(fd, operations, "state-root")
}

func validateStateRootIdentity(
	fd int,
	enrollment OwnerLockEnrollment,
	operations darwinOperations,
) error {
	if operations.EffectiveUID() != enrollment.view.ExpectedUID {
		return repair("effective-uid")
	}
	var stat unix.Stat_t
	if err := operations.Fstat(fd, &stat); err != nil {
		return repair("state-root-fstat")
	}
	if uint32(stat.Mode)&unix.S_IFMT != unix.S_IFDIR || stat.Uid != enrollment.view.ExpectedUID ||
		darwinDeviceID(stat.Dev) != enrollment.view.StateRootDevice || stat.Ino != enrollment.view.StateRootInode {
		return repair("state-root-identity")
	}
	return nil
}

func validateOwnerDescriptor(
	fd int,
	enrollment OwnerLockEnrollment,
	operations darwinOperations,
) error {
	if err := validateOwnerIdentity(fd, enrollment, operations); err != nil {
		return err
	}
	return validateDescriptorFlags(fd, operations, "owner")
}

func validateOwnerIdentity(
	fd int,
	enrollment OwnerLockEnrollment,
	operations darwinOperations,
) error {
	var stat unix.Stat_t
	if err := operations.Fstat(fd, &stat); err != nil {
		return repair("owner-fstat")
	}
	mode := uint32(stat.Mode)
	if mode&unix.S_IFMT != unix.S_IFREG || mode&0o7777 != OwnerLockMode ||
		stat.Uid != enrollment.view.ExpectedUID || uint64(stat.Nlink) != OwnerLockLinkCount ||
		darwinDeviceID(stat.Dev) != enrollment.view.LockDevice || stat.Ino != enrollment.view.LockInode {
		return repair("owner-identity")
	}
	return nil
}

func validateDescriptorFlags(fd int, operations darwinOperations, subject string) error {
	descriptorFlags, err := operations.FcntlInt(uintptr(fd), unix.F_GETFD, 0)
	if err != nil || descriptorFlags != unix.FD_CLOEXEC {
		return repair(subject + "-descriptor-flags")
	}
	statusFlags, err := operations.FcntlInt(uintptr(fd), unix.F_GETFL, 0)
	if err != nil || statusFlags&unix.O_ACCMODE != unix.O_RDONLY ||
		statusFlags&(unix.O_APPEND|unix.O_ASYNC|unix.O_NONBLOCK) != 0 {
		return repair(subject + "-status-flags")
	}
	return nil
}

func newOwnerSessionID(operations darwinOperations) (lifecyclestate.OwnerSessionID, error) {
	var value [16]byte
	read, err := operations.ReadRandom(value[:])
	if err != nil || read != len(value) {
		return lifecyclestate.OwnerSessionID{}, localFailure("owner-session-random")
	}
	domain, err := lifecyclestate.NewDomainIdentifier(lifecyclestate.DomainOwnerSessionID, value[:])
	if err != nil {
		return lifecyclestate.OwnerSessionID{}, localFailure("owner-session-invalid")
	}
	ownerSessionID, err := lifecyclestate.NewOwnerSessionID(domain)
	if err != nil || ownerSessionID.IsZero() {
		return lifecyclestate.OwnerSessionID{}, localFailure("owner-session-invalid")
	}
	return ownerSessionID, nil
}

// OwnerSessionID returns the session only while the capability remains live.
func (owner *darwinInstallationOwner) OwnerSessionID() lifecyclestate.OwnerSessionID {
	if owner == nil {
		return lifecyclestate.OwnerSessionID{}
	}
	owner.mu.Lock()
	defer owner.mu.Unlock()
	if owner.fd < 0 {
		return lifecyclestate.OwnerSessionID{}
	}
	return owner.ownerSessionID
}

// CheckHeld verifies the opaque wrapper remains live, enrolled, read-only, and
// close-on-exec and that the retained state-root entry still names the held
// inode. BSD flock is advisory; this check does not authenticate peers or prove
// unrelated processes cooperate with it.
func (owner *darwinInstallationOwner) CheckHeld(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if owner == nil {
		return ErrOwnerClosed
	}
	owner.mu.Lock()
	defer owner.mu.Unlock()
	if owner.fd < 0 {
		return ErrOwnerClosed
	}
	if owner.root == nil {
		return invariant("owner-state-root-missing")
	}
	owner.root.mu.Lock()
	defer owner.root.mu.Unlock()
	if owner.root.fd < 0 || owner.root.operations == nil {
		return invariant("owner-state-root-closed")
	}
	if err := validateOwnerDescriptor(owner.fd, owner.enrollment, owner.operations); err != nil {
		return invariant("owner-check-held")
	}
	if err := validateStateRootDescriptor(owner.root.fd, owner.enrollment, owner.root.operations); err != nil {
		return invariant("owner-check-state-root")
	}
	if err := validateNamedOwnerEntry(owner.root.fd, owner.enrollment, owner.root.operations); err != nil {
		return invariant("owner-check-entry")
	}
	if owner.ownerSessionID.IsZero() {
		return invariant("owner-session-zero")
	}
	return nil
}

// CloseAfterShutdown releases the owner descriptor only after the caller has
// stopped listeners, work, stores, and queues. There is intentionally no
// earlier unlock operation.
func (owner *darwinInstallationOwner) CloseAfterShutdown(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if owner == nil {
		return ErrOwnerClosed
	}
	owner.mu.Lock()
	defer owner.mu.Unlock()
	if owner.fd < 0 {
		return ErrOwnerClosed
	}
	if err := validateOwnerIdentity(owner.fd, owner.enrollment, owner.operations); err != nil {
		owner.fd = -1
		owner.ownerSessionID = lifecyclestate.OwnerSessionID{}
		return invariant("owner-close-identity")
	}
	fd := owner.fd
	owner.fd = -1
	owner.ownerSessionID = lifecyclestate.OwnerSessionID{}
	flagsErr := validateDescriptorFlags(fd, owner.operations, "owner")
	if err := owner.operations.Close(fd); err != nil {
		return localFailure("owner-close")
	}
	if flagsErr != nil {
		return invariant("owner-close-flags")
	}
	return nil
}

func repair(code string) error { return fmt.Errorf("%w: %s", ErrRepairRequired, code) }

func invariant(code string) error { return fmt.Errorf("%w: %s", ErrOwnerInvariant, code) }

func localFailure(code string) error { return fmt.Errorf("%w: %s", ErrLocalFailure, code) }

func darwinDeviceID(device int32) uint64 {
	// Darwin exposes dev_t through x/sys as int32. Convert through uint32 to
	// preserve its kernel bit pattern without sign extension.
	return uint64(uint32(device)) //nolint:gosec // Intentional dev_t bit-pattern normalization.
}
