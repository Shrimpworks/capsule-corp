//go:build darwin

package installationowner

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
	"golang.org/x/sys/unix"
)

const (
	ownerHelperModeEnvironment       = "CAPSULE_OWNER_HELPER_MODE"
	ownerHelperRootEnvironment       = "CAPSULE_OWNER_HELPER_ROOT"
	ownerHelperEnrollmentEnvironment = "CAPSULE_OWNER_HELPER_ENROLLMENT"
	ownerHelperMarkersEnvironment    = "CAPSULE_OWNER_HELPER_MARKERS"
)

type ownerFixture struct {
	stateRootPath string
	lockPath      string
	enrollment    OwnerLockEnrollment
	root          *TrustedStateRoot
}

func newOwnerFixture(t *testing.T) ownerFixture {
	t.Helper()
	stateRootPath := filepath.Join(t.TempDir(), "supervisor-state")
	if err := os.Mkdir(stateRootPath, 0o700); err != nil {
		t.Fatal(err)
	}
	lockPath := filepath.Join(stateRootPath, "supervisor.owner")
	file, err := os.OpenFile(lockPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, os.FileMode(OwnerLockMode))
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	var rootStat unix.Stat_t
	if err := unix.Stat(stateRootPath, &rootStat); err != nil {
		t.Fatal(err)
	}
	var lockStat unix.Stat_t
	if err := unix.Lstat(lockPath, &lockStat); err != nil {
		t.Fatal(err)
	}
	enrollment, err := NewOwnerLockEnrollment(OwnerLockEnrollmentView{
		InstallationID:  repeated16[v0candidate.InstallationID](0x11),
		SupervisorID:    repeated16[v0candidate.SupervisorID](0x55),
		ExpectedUID:     uint32(unix.Geteuid()),
		StateRootDevice: darwinDeviceID(rootStat.Dev),
		StateRootInode:  rootStat.Ino,
		LockEntryName:   filepath.Base(lockPath),
		LockDevice:      darwinDeviceID(lockStat.Dev),
		LockInode:       lockStat.Ino,
	})
	if err != nil {
		t.Fatal(err)
	}
	root, err := OpenDarwinTrustedStateRoot(context.Background(), stateRootPath, enrollment)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if root.fd >= 0 {
			_ = root.Close()
		}
	})
	return ownerFixture{stateRootPath: stateRootPath, lockPath: lockPath, enrollment: enrollment, root: root}
}

func TestDarwinOwnerAcceptsOnlyEnrolledRegularObject(t *testing.T) {
	fixture := newOwnerFixture(t)
	rootDescriptorFlags, err := unix.FcntlInt(uintptr(fixture.root.fd), unix.F_GETFD, 0)
	if err != nil || rootDescriptorFlags != unix.FD_CLOEXEC {
		t.Fatalf("state-root descriptor flags = %#x, %v", rootDescriptorFlags, err)
	}
	rootStatusFlags, err := unix.FcntlInt(uintptr(fixture.root.fd), unix.F_GETFL, 0)
	if err != nil || rootStatusFlags&unix.O_ACCMODE != unix.O_RDONLY {
		t.Fatalf("state-root status flags = %#x, %v", rootStatusFlags, err)
	}
	operations := &faultDarwinOperations{darwinOperations: fixture.root.operations}
	fixture.root.operations = operations
	owner, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment)
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}
	if owner.OwnerSessionID().IsZero() {
		t.Fatal("acquisition issued a zero owner-session ID")
	}
	if err := owner.CheckHeld(context.Background()); err != nil {
		t.Fatalf("check held: %v", err)
	}
	concrete := owner.(*darwinInstallationOwner)
	descriptorFlags, err := unix.FcntlInt(uintptr(concrete.fd), unix.F_GETFD, 0)
	if err != nil || descriptorFlags != unix.FD_CLOEXEC {
		t.Fatalf("owner descriptor flags = %#x, %v", descriptorFlags, err)
	}
	statusFlags, err := unix.FcntlInt(uintptr(concrete.fd), unix.F_GETFL, 0)
	if err != nil || statusFlags&unix.O_ACCMODE != unix.O_RDONLY {
		t.Fatalf("owner status flags = %#x, %v", statusFlags, err)
	}
	if len(operations.openatFlags) != 3 || operations.openatFlags[0] != ownerLockOpenFlags ||
		operations.openatFlags[1] != ownerLockOpenFlags || operations.openatFlags[2] != ownerLockOpenFlags {
		t.Fatalf("openat flags = %#v, want three %#x", operations.openatFlags, ownerLockOpenFlags)
	}
	if len(operations.flockOperations) != 1 || operations.flockOperations[0] != unix.LOCK_EX|unix.LOCK_NB {
		t.Fatalf("flock operations = %#v", operations.flockOperations)
	}
	if err := owner.CloseAfterShutdown(context.Background()); err != nil {
		t.Fatalf("close after shutdown: %v", err)
	}
	if !owner.OwnerSessionID().IsZero() {
		t.Fatal("closed owner retained a usable owner-session ID")
	}
	if err := owner.CheckHeld(context.Background()); !errors.Is(err, ErrOwnerClosed) {
		t.Fatalf("closed check = %v", err)
	}
	if err := owner.CloseAfterShutdown(context.Background()); !errors.Is(err, ErrOwnerClosed) {
		t.Fatalf("repeat close = %v", err)
	}
}

func TestDarwinOwnerRefusesSymlinkModeOwnerLinkAndIdentityMismatch(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*testing.T, ownerFixture) OwnerLockEnrollment
	}{
		{name: "missing", mutate: func(t *testing.T, fixture ownerFixture) OwnerLockEnrollment {
			if err := os.Remove(fixture.lockPath); err != nil {
				t.Fatal(err)
			}
			return fixture.enrollment
		}},
		{name: "symlink", mutate: func(t *testing.T, fixture ownerFixture) OwnerLockEnrollment {
			target := fixture.lockPath + ".target"
			if err := os.Rename(fixture.lockPath, target); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(filepath.Base(target), fixture.lockPath); err != nil {
				t.Fatal(err)
			}
			return fixture.enrollment
		}},
		{name: "mode", mutate: func(t *testing.T, fixture ownerFixture) OwnerLockEnrollment {
			if err := os.Chmod(fixture.lockPath, 0o640); err != nil {
				t.Fatal(err)
			}
			return fixture.enrollment
		}},
		{name: "hard link", mutate: func(t *testing.T, fixture ownerFixture) OwnerLockEnrollment {
			if err := os.Link(fixture.lockPath, fixture.lockPath+".link"); err != nil {
				t.Fatal(err)
			}
			return fixture.enrollment
		}},
		{name: "directory", mutate: func(t *testing.T, fixture ownerFixture) OwnerLockEnrollment {
			if err := os.Remove(fixture.lockPath); err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(fixture.lockPath, 0o600); err != nil {
				t.Fatal(err)
			}
			return fixture.enrollment
		}},
		{name: "effective uid", mutate: func(t *testing.T, fixture ownerFixture) OwnerLockEnrollment {
			view := fixture.enrollment.View()
			view.ExpectedUID++
			enrollment, err := NewOwnerLockEnrollment(view)
			if err != nil {
				t.Fatal(err)
			}
			return enrollment
		}},
		{name: "lock inode", mutate: func(t *testing.T, fixture ownerFixture) OwnerLockEnrollment {
			view := fixture.enrollment.View()
			view.LockInode++
			enrollment, err := NewOwnerLockEnrollment(view)
			if err != nil {
				t.Fatal(err)
			}
			return enrollment
		}},
		{name: "lock device", mutate: func(t *testing.T, fixture ownerFixture) OwnerLockEnrollment {
			view := fixture.enrollment.View()
			view.LockDevice++
			enrollment, err := NewOwnerLockEnrollment(view)
			if err != nil {
				t.Fatal(err)
			}
			return enrollment
		}},
		{name: "root inode", mutate: func(t *testing.T, fixture ownerFixture) OwnerLockEnrollment {
			view := fixture.enrollment.View()
			view.StateRootInode++
			enrollment, err := NewOwnerLockEnrollment(view)
			if err != nil {
				t.Fatal(err)
			}
			return enrollment
		}},
		{name: "root device", mutate: func(t *testing.T, fixture ownerFixture) OwnerLockEnrollment {
			view := fixture.enrollment.View()
			view.StateRootDevice++
			enrollment, err := NewOwnerLockEnrollment(view)
			if err != nil {
				t.Fatal(err)
			}
			return enrollment
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newOwnerFixture(t)
			enrollment := test.mutate(t, fixture)
			fixture.root.enrollment = enrollment
			beforeEntries := directoryEntries(t, fixture.stateRootPath)
			owner, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, enrollment)
			if owner != nil || (!errors.Is(err, ErrRepairRequired) && !errors.Is(err, ErrInvalidEnrollment)) {
				t.Fatalf("refusal = %T, %v", owner, err)
			}
			afterEntries := directoryEntries(t, fixture.stateRootPath)
			if strings.Join(afterEntries, "\x00") != strings.Join(beforeEntries, "\x00") {
				t.Fatalf("refusal changed state-root entries: before=%v after=%v", beforeEntries, afterEntries)
			}
		})
	}
}

func TestDarwinOwnerDoesNotCreateMissingObject(t *testing.T) {
	fixture := newOwnerFixture(t)
	if err := os.Remove(fixture.lockPath); err != nil {
		t.Fatal(err)
	}
	if _, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment); !errors.Is(err, ErrRepairRequired) {
		t.Fatalf("missing result = %v", err)
	}
	if _, err := os.Lstat(fixture.lockPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing acquisition created an object: %v", err)
	}
}

func TestDarwinOwnerDuplicateProcessRefusesBeforeAnyDownstreamWork(t *testing.T) {
	fixture := newOwnerFixture(t)
	storePath := filepath.Join(fixture.stateRootPath, "supervisor-state.json")
	storeBytes := []byte("retained-store-bytes\n")
	if err := os.WriteFile(storePath, storeBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	beforeDigest := sha256.Sum256(storeBytes)
	markers := filepath.Join(fixture.stateRootPath, "markers")
	if err := os.Mkdir(markers, 0o700); err != nil {
		t.Fatal(err)
	}

	holder, holderInput := startOwnerHelper(t, "hold", fixture, "")
	defer func() {
		_ = holderInput.Close()
		if holder.ProcessState == nil {
			_ = holder.Process.Kill()
		}
		_ = holder.Wait()
	}()

	for repetition := range 20 {
		output := runOwnerHelper(t, "contend", fixture, markers)
		if firstOutputLine(output) != "BUSY" {
			t.Fatalf("contender %d = %q", repetition, output)
		}
	}
	entries, err := os.ReadDir(markers)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("duplicate owner reached downstream markers: %v", entries)
	}
	afterBytes, err := os.ReadFile(storePath)
	if err != nil {
		t.Fatal(err)
	}
	if sha256.Sum256(afterBytes) != beforeDigest {
		t.Fatal("duplicate owner changed store bytes")
	}
}

func TestDarwinOwnerProcessDeathAllowsReopen(t *testing.T) {
	fixture := newOwnerFixture(t)
	holder, holderInput := startOwnerHelper(t, "hold", fixture, "")
	_ = holderInput.Close()
	if err := holder.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = holder.Wait()
	owner, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment)
	if err != nil {
		t.Fatalf("reopen after death: %v", err)
	}
	if err := owner.CloseAfterShutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestDarwinOwnerExplicitDescriptorInheritanceAndCLOEXEC(t *testing.T) {
	t.Run("ordinary exec omits owner descriptor", func(t *testing.T) {
		fixture := newOwnerFixture(t)
		output := runOwnerHelper(t, "normal-exec", fixture, "")
		pid := parsePID(t, output)
		defer killPID(pid)
		owner, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment)
		if err != nil {
			t.Fatalf("ordinary exec retained CLOEXEC owner: %v", err)
		}
		if err := owner.CloseAfterShutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("explicit inherited description delays successor", func(t *testing.T) {
		fixture := newOwnerFixture(t)
		output := runOwnerHelper(t, "explicit-inherit", fixture, "")
		pid := parsePID(t, output)
		if _, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment); !errors.Is(err, ErrDuplicateOwner) {
			killPID(pid)
			t.Fatalf("explicit inherited description result = %v", err)
		}
		killPID(pid)
		deadline := time.Now().Add(3 * time.Second)
		for {
			owner, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment)
			if err == nil {
				if closeErr := owner.CloseAfterShutdown(context.Background()); closeErr != nil {
					t.Fatal(closeErr)
				}
				break
			}
			if !errors.Is(err, ErrDuplicateOwner) || time.Now().After(deadline) {
				t.Fatalf("successor after inherited close = %v", err)
			}
			time.Sleep(10 * time.Millisecond)
		}
	})
}

func TestDarwinOwnerRenameUnlinkReplacementNeverMatchesEnrollment(t *testing.T) {
	fixture := newOwnerFixture(t)
	owner, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment)
	if err != nil {
		t.Fatal(err)
	}
	moved := fixture.lockPath + ".moved"
	if err := os.Rename(fixture.lockPath, moved); err != nil {
		t.Fatal(err)
	}
	replacement, err := os.OpenFile(fixture.lockPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, os.FileMode(OwnerLockMode))
	if err != nil {
		t.Fatal(err)
	}
	if err := replacement.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment); !errors.Is(err, ErrRepairRequired) {
		t.Fatalf("replacement while held = %v", err)
	}
	if err := owner.CheckHeld(context.Background()); !errors.Is(err, ErrOwnerInvariant) {
		t.Fatalf("replaced enrolled entry check = %v", err)
	}
	if err := owner.CloseAfterShutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment); !errors.Is(err, ErrRepairRequired) {
		t.Fatalf("replacement after close = %v", err)
	}
	if err := os.Remove(fixture.lockPath); err != nil {
		t.Fatal(err)
	}
	if _, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment); !errors.Is(err, ErrRepairRequired) {
		t.Fatalf("unlink result = %v", err)
	}
}

func TestDarwinOwnerReplacementBetweenOpenAndRecheckRefuses(t *testing.T) {
	fixture := newOwnerFixture(t)
	operations := &faultDarwinOperations{darwinOperations: fixture.root.operations}
	operations.beforeOpenat = func(call int) {
		if call != 2 {
			return
		}
		moved := fixture.lockPath + ".raced"
		if err := os.Rename(fixture.lockPath, moved); err != nil {
			t.Errorf("rename during recheck: %v", err)
			return
		}
		replacement, err := os.OpenFile(
			fixture.lockPath,
			os.O_WRONLY|os.O_CREATE|os.O_EXCL,
			os.FileMode(OwnerLockMode),
		)
		if err != nil {
			t.Errorf("replace during recheck: %v", err)
			return
		}
		if err := replacement.Close(); err != nil {
			t.Errorf("close replacement: %v", err)
		}
	}
	fixture.root.operations = operations
	owner, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment)
	if owner != nil || !errors.Is(err, ErrRepairRequired) {
		t.Fatalf("replacement race result = %T, %v", owner, err)
	}
	if operations.closeCalls < 2 {
		t.Fatalf("replacement race did not close recheck and held descriptors: %d", operations.closeCalls)
	}
}

func TestDarwinOwnerRetainedRootDescriptorSurvivesPathReplacement(t *testing.T) {
	fixture := newOwnerFixture(t)
	movedRoot := fixture.stateRootPath + ".moved"
	if err := os.Rename(fixture.stateRootPath, movedRoot); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(fixture.stateRootPath, 0o700); err != nil {
		t.Fatal(err)
	}
	owner, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment)
	if err != nil {
		t.Fatalf("retained descriptor followed replacement pathname: %v", err)
	}
	if err := owner.CloseAfterShutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenDarwinTrustedStateRoot(context.Background(), fixture.stateRootPath, fixture.enrollment); !errors.Is(err, ErrRepairRequired) {
		t.Fatalf("replacement pathname matched enrolled root: %v", err)
	}
}

func TestDarwinOwnerDescriptorReuseIsDetectedWithoutClosingForeignFD(t *testing.T) {
	fixture := newOwnerFixture(t)
	capability, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment)
	if err != nil {
		t.Fatal(err)
	}
	owner := capability.(*darwinInstallationOwner)
	reusedFD := owner.fd
	if err := unix.Close(reusedFD); err != nil {
		t.Fatal(err)
	}
	foreignPath := filepath.Join(fixture.stateRootPath, "foreign")
	if err := os.WriteFile(foreignPath, []byte("foreign"), 0o600); err != nil {
		t.Fatal(err)
	}
	foreignFD, err := unix.Open(foreignPath, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		t.Fatal(err)
	}
	if foreignFD != reusedFD {
		_ = unix.Close(foreignFD)
		t.Skipf("descriptor allocator returned %d after %d; reuse oracle unavailable", foreignFD, reusedFD)
	}
	if err := capability.CheckHeld(context.Background()); !errors.Is(err, ErrOwnerInvariant) {
		t.Fatalf("reused descriptor check = %v", err)
	}
	if err := capability.CloseAfterShutdown(context.Background()); !errors.Is(err, ErrOwnerInvariant) {
		t.Fatalf("reused descriptor close = %v", err)
	}
	buffer := make([]byte, 7)
	if _, err := unix.Read(foreignFD, buffer); err != nil || string(buffer) != "foreign" {
		t.Fatalf("owner closed or changed reused foreign descriptor: %q, %v", buffer, err)
	}
	if err := unix.Close(foreignFD); err != nil {
		t.Fatal(err)
	}
}

func TestDarwinStateRootDescriptorReuseIsDetectedWithoutClosingForeignFD(t *testing.T) {
	fixture := newOwnerFixture(t)
	root := fixture.root
	reusedFD := root.fd
	if err := unix.Close(reusedFD); err != nil {
		t.Fatal(err)
	}
	foreignPath := filepath.Join(fixture.stateRootPath, "foreign")
	if err := os.WriteFile(foreignPath, []byte("foreign"), 0o600); err != nil {
		t.Fatal(err)
	}
	foreignFD, err := unix.Open(foreignPath, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		t.Fatal(err)
	}
	if foreignFD != reusedFD {
		_ = unix.Close(foreignFD)
		t.Skipf("descriptor allocator returned %d after %d; reuse oracle unavailable", foreignFD, reusedFD)
	}
	if err := root.Close(); !errors.Is(err, ErrOwnerInvariant) {
		t.Fatalf("reused descriptor close = %v", err)
	}
	if root.fd >= 0 {
		t.Fatal("state root retained a usable descriptor after an identity-mismatch close")
	}
	buffer := make([]byte, 7)
	if _, err := unix.Read(foreignFD, buffer); err != nil || string(buffer) != "foreign" {
		t.Fatalf("state root closed or changed reused foreign descriptor: %q, %v", buffer, err)
	}
	if err := unix.Close(foreignFD); err != nil {
		t.Fatal(err)
	}
}

func TestDarwinOwnerFaultsRefuseWithoutSession(t *testing.T) {
	tests := []struct {
		name      string
		configure func(*faultDarwinOperations)
		want      error
	}{
		{name: "openat", configure: func(ops *faultDarwinOperations) { ops.failOpenatAt = 1 }, want: ErrRepairRequired},
		{name: "fstat", configure: func(ops *faultDarwinOperations) { ops.failFstatAt = 2 }, want: ErrRepairRequired},
		{name: "fcntl", configure: func(ops *faultDarwinOperations) { ops.failFcntlAt = 1 }, want: ErrRepairRequired},
		{name: "flock busy", configure: func(ops *faultDarwinOperations) { ops.flockError = unix.EWOULDBLOCK }, want: ErrDuplicateOwner},
		{name: "flock local", configure: func(ops *faultDarwinOperations) { ops.flockError = unix.EIO }, want: ErrLocalFailure},
		{name: "random", configure: func(ops *faultDarwinOperations) { ops.randomError = io.ErrUnexpectedEOF }, want: ErrLocalFailure},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newOwnerFixture(t)
			operations := &faultDarwinOperations{darwinOperations: fixture.root.operations}
			test.configure(operations)
			fixture.root.operations = operations
			owner, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment)
			if owner != nil || !errors.Is(err, test.want) {
				t.Fatalf("fault result = %T, %v", owner, err)
			}
		})
	}

	t.Run("close", func(t *testing.T) {
		fixture := newOwnerFixture(t)
		ownerInterface, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment)
		if err != nil {
			t.Fatal(err)
		}
		owner := ownerInterface.(*darwinInstallationOwner)
		operations := &faultDarwinOperations{darwinOperations: owner.operations, failCloseAt: 1}
		owner.operations = operations
		if err := owner.CloseAfterShutdown(context.Background()); !errors.Is(err, ErrLocalFailure) {
			t.Fatalf("close fault = %v", err)
		}
		if !owner.OwnerSessionID().IsZero() {
			t.Fatal("close fault retained session")
		}
	})

	t.Run("canceled", func(t *testing.T) {
		fixture := newOwnerFixture(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		owner, err := AcquireDarwinInstallationOwner(ctx, fixture.root, fixture.enrollment)
		if owner != nil || !errors.Is(err, context.Canceled) {
			t.Fatalf("canceled result = %T, %v", owner, err)
		}
	})

	t.Run("canceled close retains owner until orderly close", func(t *testing.T) {
		fixture := newOwnerFixture(t)
		owner, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment)
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if err := owner.CloseAfterShutdown(ctx); !errors.Is(err, context.Canceled) {
			t.Fatalf("canceled close = %v", err)
		}
		if owner.OwnerSessionID().IsZero() {
			t.Fatal("canceled close released the live owner")
		}
		if err := owner.CloseAfterShutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
}

func TestDarwinOwnerRepeatedAcquireClose(t *testing.T) {
	fixture := newOwnerFixture(t)
	seen := make(map[lifecyclestate.OwnerSessionID]struct{})
	for repetition := range 200 {
		owner, err := AcquireDarwinInstallationOwner(context.Background(), fixture.root, fixture.enrollment)
		if err != nil {
			t.Fatalf("acquire %d: %v", repetition, err)
		}
		session := owner.OwnerSessionID()
		if session.IsZero() {
			t.Fatalf("session %d is zero", repetition)
		}
		if _, duplicate := seen[session]; duplicate {
			t.Fatalf("session %d repeated", repetition)
		}
		seen[session] = struct{}{}
		if err := owner.CheckHeld(context.Background()); err != nil {
			t.Fatalf("check %d: %v", repetition, err)
		}
		if err := owner.CloseAfterShutdown(context.Background()); err != nil {
			t.Fatalf("close %d: %v", repetition, err)
		}
	}
}

func TestDarwinOwnerProcessHelper(t *testing.T) {
	mode := os.Getenv(ownerHelperModeEnvironment)
	if mode == "" {
		return
	}
	var view OwnerLockEnrollmentView
	if err := json.Unmarshal([]byte(os.Getenv(ownerHelperEnrollmentEnvironment)), &view); err != nil {
		fmt.Printf("ERROR enrollment %v\n", err)
		return
	}
	enrollment, err := NewOwnerLockEnrollment(view)
	if err != nil {
		fmt.Printf("ERROR enrollment %v\n", err)
		return
	}
	root, err := OpenDarwinTrustedStateRoot(context.Background(), os.Getenv(ownerHelperRootEnvironment), enrollment)
	if err != nil {
		fmt.Printf("ERROR root %v\n", err)
		return
	}
	owner, err := AcquireDarwinInstallationOwner(context.Background(), root, enrollment)
	if errors.Is(err, ErrDuplicateOwner) {
		fmt.Println("BUSY")
		_ = root.Close()
		return
	}
	if err != nil {
		fmt.Printf("ERROR acquire %v\n", err)
		_ = root.Close()
		return
	}
	switch mode {
	case "hold":
		fmt.Println("READY")
		_, _ = io.Copy(io.Discard, os.Stdin)
	case "contend":
		for _, marker := range []string{"store-open", "store-mutation", "migration", "recovery", "ipc", "archive", "adapter"} {
			// The marker root is created by this test under t.TempDir and passed only to this owned helper.
			_ = os.WriteFile(filepath.Join(os.Getenv(ownerHelperMarkersEnvironment), marker), []byte("reached\n"), 0o600) //nolint:gosec
		}
		fmt.Println("ACQUIRED")
	case "normal-exec":
		child := exec.Command("/bin/sleep", "30")
		child.Stdout = io.Discard
		child.Stderr = io.Discard
		if err := child.Start(); err != nil {
			fmt.Printf("ERROR exec %v\n", err)
			break
		}
		fmt.Println(child.Process.Pid)
	case "explicit-inherit":
		concrete := owner.(*darwinInstallationOwner)
		duplicate, duplicateErr := unix.Dup(concrete.fd)
		if duplicateErr != nil {
			fmt.Printf("ERROR dup %v\n", duplicateErr)
			break
		}
		file := os.NewFile(uintptr(duplicate), "explicit-owner-inheritance")
		child := exec.Command("/bin/sleep", "30")
		child.Stdout = io.Discard
		child.Stderr = io.Discard
		child.ExtraFiles = []*os.File{file}
		if err := child.Start(); err != nil {
			_ = file.Close()
			fmt.Printf("ERROR exec %v\n", err)
			break
		}
		_ = file.Close()
		fmt.Println(child.Process.Pid)
	}
	_ = owner.CloseAfterShutdown(context.Background())
	_ = root.Close()
}

type faultDarwinOperations struct {
	darwinOperations
	openatCalls     int
	fstatCalls      int
	fcntlCalls      int
	closeCalls      int
	openatFlags     []int
	flockOperations []int
	beforeOpenat    func(int)
	failOpenatAt    int
	failFstatAt     int
	failFcntlAt     int
	failCloseAt     int
	flockError      error
	randomError     error
}

func (operations *faultDarwinOperations) Openat(directoryFD int, name string, flags int, mode uint32) (int, error) {
	operations.openatCalls++
	operations.openatFlags = append(operations.openatFlags, flags)
	if operations.beforeOpenat != nil {
		operations.beforeOpenat(operations.openatCalls)
	}
	if operations.openatCalls == operations.failOpenatAt {
		return -1, unix.EIO
	}
	return operations.darwinOperations.Openat(directoryFD, name, flags, mode)
}

func (operations *faultDarwinOperations) Fstat(fd int, stat *unix.Stat_t) error {
	operations.fstatCalls++
	if operations.fstatCalls == operations.failFstatAt {
		return unix.EIO
	}
	return operations.darwinOperations.Fstat(fd, stat)
}

func (operations *faultDarwinOperations) FcntlInt(fd uintptr, command int, argument int) (int, error) {
	operations.fcntlCalls++
	if operations.fcntlCalls == operations.failFcntlAt {
		return 0, unix.EIO
	}
	return operations.darwinOperations.FcntlInt(fd, command, argument)
}

func (operations *faultDarwinOperations) Flock(fd int, operation int) error {
	operations.flockOperations = append(operations.flockOperations, operation)
	if operations.flockError != nil {
		return operations.flockError
	}
	return operations.darwinOperations.Flock(fd, operation)
}

func (operations *faultDarwinOperations) Close(fd int) error {
	operations.closeCalls++
	err := operations.darwinOperations.Close(fd)
	if operations.closeCalls == operations.failCloseAt {
		return unix.EIO
	}
	return err
}

func (operations *faultDarwinOperations) ReadRandom(destination []byte) (int, error) {
	if operations.randomError != nil {
		return 0, operations.randomError
	}
	return operations.darwinOperations.ReadRandom(destination)
}

func directoryEntries(t *testing.T, path string) []string {
	t.Helper()
	entries, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, len(entries))
	for index, entry := range entries {
		names[index] = entry.Name()
	}
	return names
}

func ownerHelperCommand(t *testing.T, mode string, fixture ownerFixture, markers string) *exec.Cmd {
	t.Helper()
	encoded, err := json.Marshal(fixture.enrollment.View())
	if err != nil {
		t.Fatal(err)
	}
	// os.Args[0] is this exact test binary; no caller-controlled executable or argument is accepted.
	command := exec.Command(os.Args[0], "-test.run=^TestDarwinOwnerProcessHelper$") //nolint:gosec
	command.Env = append(os.Environ(),
		ownerHelperModeEnvironment+"="+mode,
		ownerHelperRootEnvironment+"="+fixture.stateRootPath,
		ownerHelperEnrollmentEnvironment+"="+string(encoded),
		ownerHelperMarkersEnvironment+"="+markers,
	)
	return command
}

func startOwnerHelper(t *testing.T, mode string, fixture ownerFixture, markers string) (*exec.Cmd, io.WriteCloser) {
	t.Helper()
	command := ownerHelperCommand(t, mode, fixture, markers)
	stdout, err := command.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	command.Stderr = &stderr
	input, err := command.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() || scanner.Text() != "READY" {
		_ = command.Process.Kill()
		_ = command.Wait()
		t.Fatalf("holder readiness = %q, stderr=%q", scanner.Text(), stderr.String())
	}
	return command, input
}

func runOwnerHelper(t *testing.T, mode string, fixture ownerFixture, markers string) string {
	t.Helper()
	command := ownerHelperCommand(t, mode, fixture, markers)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("helper %s: %v: %s", mode, err, output)
	}
	return string(output)
}

func parsePID(t *testing.T, output string) int {
	t.Helper()
	pid, err := strconv.Atoi(firstOutputLine(output))
	if err != nil {
		t.Fatalf("helper PID %q: %v", output, err)
	}
	return pid
}

func firstOutputLine(output string) string {
	line, _, _ := strings.Cut(strings.TrimSpace(output), "\n")
	return line
}

func killPID(pid int) {
	process, err := os.FindProcess(pid)
	if err == nil {
		_ = process.Kill()
	}
}
