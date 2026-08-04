package installationowner

import (
	"context"
	"errors"

	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// OwnerLockMode is the exact enrolled permission mode for the pre-created
// regular owner object. Normal startup never creates or repairs the object.
const OwnerLockMode uint32 = 0o600

// OwnerLockLinkCount is the only accepted link count for the owner object.
const OwnerLockLinkCount uint64 = 1

const maxClosedEntryNameBytes = 255

var (
	// ErrUnsupportedPlatform reports that the exact Darwin owner primitive is
	// unavailable on the current build target.
	ErrUnsupportedPlatform = errors.New("installation owner is unsupported on this platform")
	// ErrInvalidEnrollment reports a closed bootstrap projection that cannot
	// identify one installation, Supervisor, state root, and lock object.
	ErrInvalidEnrollment = errors.New("installation owner enrollment is invalid")
	// ErrRepairRequired reports missing or mismatched enrolled filesystem state.
	// Acquisition never repairs, replaces, or creates that state.
	ErrRepairRequired = errors.New("installation owner state requires repair")
	// ErrDuplicateOwner is the fixed nonblocking result when another open file
	// description already owns the BSD flock.
	ErrDuplicateOwner = errors.New("installation owner is already held")
	// ErrOwnerInvariant reports loss or corruption of a live opaque capability.
	// Callers must treat it as terminal rather than ordinary readiness failure.
	ErrOwnerInvariant = errors.New("installation owner invariant failed")
	// ErrOwnerClosed reports use after orderly close.
	ErrOwnerClosed = errors.New("installation owner is closed")
	// ErrLocalFailure reports a bounded local operation failure that is not an
	// enrolled identity mismatch or duplicate owner.
	ErrLocalFailure = errors.New("installation owner local operation failed")
)

// OwnerLockEnrollmentView is the closed value projection needed by passive G1
// acquisition. A later installed bootstrap must authenticate these values; G1
// does not establish their provenance. The trusted installer, not normal
// Supervisor startup, creates and enrolls the referenced root and lock object.
type OwnerLockEnrollmentView struct {
	InstallationID  v0candidate.InstallationID
	SupervisorID    v0candidate.SupervisorID
	ExpectedUID     uint32
	StateRootDevice uint64
	StateRootInode  uint64
	LockEntryName   string
	LockDevice      uint64
	LockInode       uint64
}

// OwnerLockEnrollment is an immutable validated bootstrap projection.
type OwnerLockEnrollment struct {
	view OwnerLockEnrollmentView
}

// NewOwnerLockEnrollment validates and closes the bootstrap projection before
// any filesystem descriptor is accepted.
func NewOwnerLockEnrollment(view OwnerLockEnrollmentView) (OwnerLockEnrollment, error) {
	if err := validateEnrollmentView(view); err != nil {
		return OwnerLockEnrollment{}, err
	}
	return OwnerLockEnrollment{view: view}, nil
}

// View returns the value-only closed enrollment projection.
func (enrollment OwnerLockEnrollment) View() OwnerLockEnrollmentView {
	return enrollment.view
}

func validateEnrollmentView(view OwnerLockEnrollmentView) error {
	if view.InstallationID == (v0candidate.InstallationID{}) ||
		view.SupervisorID == (v0candidate.SupervisorID{}) ||
		view.ExpectedUID == 0 ||
		view.StateRootDevice == 0 || view.StateRootInode == 0 ||
		view.LockDevice == 0 || view.LockInode == 0 ||
		!validClosedEntryName(view.LockEntryName) {
		return ErrInvalidEnrollment
	}
	return nil
}

func validClosedEntryName(name string) bool {
	if name == "" || len(name) > maxClosedEntryNameBytes || name == "." || name == ".." {
		return false
	}
	for index := range len(name) {
		character := name[index]
		if character < 0x21 || character > 0x7e || character == '/' || character == '\\' {
			return false
		}
	}
	return true
}

// InstallationOwner is the opaque lifetime owner capability. It intentionally
// exposes no descriptor, file, duplicate, unlock, or transfer operation.
type InstallationOwner interface {
	OwnerSessionID() lifecyclestate.OwnerSessionID
	CheckHeld(context.Context) error
	CloseAfterShutdown(context.Context) error
	sealedInstallationOwner()
}
