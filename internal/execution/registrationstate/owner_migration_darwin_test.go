//go:build darwin

package registrationstate

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"capsule.local/capsule/internal/execution/installationowner"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"golang.org/x/sys/unix"
)

type heldInstallationOwnerMigrationLock struct {
	owner installationowner.InstallationOwner
}

func (lock heldInstallationOwnerMigrationLock) CheckOfflineMigrationLock(ctx context.Context) error {
	return lock.owner.CheckHeld(ctx)
}

func TestOwnerRequiredMigrationUsesActualDarwinInstallationOwner(t *testing.T) {
	harness := newApprovalHarness(t)
	state, err := harness.store.snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	stateRootPath := filepath.Dir(harness.path)
	lockPath := filepath.Join(stateRootPath, "supervisor.owner")
	lockFile, err := os.OpenFile(
		lockPath,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		os.FileMode(installationowner.OwnerLockMode),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := lockFile.Close(); err != nil {
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
	enrollment, err := installationowner.NewOwnerLockEnrollment(installationowner.OwnerLockEnrollmentView{
		InstallationID: state.InstallationID, SupervisorID: state.SupervisorID,
		ExpectedUID:     uint32(unix.Geteuid()),                                     //nolint:gosec // Local Darwin test UID is nonnegative.
		StateRootDevice: uint64(uint32(rootStat.Dev)), StateRootInode: rootStat.Ino, //nolint:gosec // Preserve Darwin dev_t bits.
		LockEntryName: filepath.Base(lockPath), LockDevice: uint64(uint32(lockStat.Dev)), LockInode: lockStat.Ino, //nolint:gosec // Preserve Darwin dev_t bits.
	})
	if err != nil {
		t.Fatal(err)
	}
	root, err := installationowner.OpenDarwinTrustedStateRoot(context.Background(), stateRootPath, enrollment)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = root.Close() }()
	owner, err := installationowner.AcquireDarwinInstallationOwner(context.Background(), root, enrollment)
	if err != nil {
		t.Fatal(err)
	}

	_, err = MigrateFixedFileStoreV0ToV1(
		context.Background(),
		harness.path,
		V0ToV1MigrationOptions{Lock: heldInstallationOwnerMigrationLock{owner: owner}},
	)
	if err != nil {
		t.Fatalf("migration under actual owner: %v", err)
	}
	if err := owner.CloseAfterShutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	before := mustReadFile(t, harness.path)
	if _, err := MigrateFixedFileStoreV0ToV1(
		context.Background(),
		harness.path,
		V0ToV1MigrationOptions{Lock: heldInstallationOwnerMigrationLock{owner: owner}},
	); !errors.Is(err, ErrMigrationLockRequired) {
		t.Fatalf("closed owner migration = %v", err)
	}
	if after := mustReadFile(t, harness.path); string(after) != string(before) {
		t.Fatal("closed-owner refusal rewrote the store")
	}
}

func TestOwnerRequiredStoreOpenBindsSessionAndFencesReplacement(t *testing.T) {
	harness := newApprovalHarness(t)
	state, err := harness.store.snapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	stateRootPath := filepath.Dir(harness.path)
	lockPath := filepath.Join(stateRootPath, "supervisor.owner")
	lockFile, err := os.OpenFile(
		lockPath,
		os.O_WRONLY|os.O_CREATE|os.O_EXCL,
		os.FileMode(installationowner.OwnerLockMode),
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := lockFile.Close(); err != nil {
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
	enrollment, err := installationowner.NewOwnerLockEnrollment(installationowner.OwnerLockEnrollmentView{
		InstallationID: state.InstallationID, SupervisorID: state.SupervisorID,
		ExpectedUID:     uint32(unix.Geteuid()),                                     //nolint:gosec // Local Darwin test UID is nonnegative.
		StateRootDevice: uint64(uint32(rootStat.Dev)), StateRootInode: rootStat.Ino, //nolint:gosec // Preserve Darwin dev_t bits.
		LockEntryName: filepath.Base(lockPath), LockDevice: uint64(uint32(lockStat.Dev)), LockInode: lockStat.Ino, //nolint:gosec // Preserve Darwin dev_t bits.
	})
	if err != nil {
		t.Fatal(err)
	}
	root, err := installationowner.OpenDarwinTrustedStateRoot(context.Background(), stateRootPath, enrollment)
	if err != nil {
		t.Fatal(err)
	}
	owner, err := installationowner.AcquireDarwinInstallationOwner(context.Background(), root, enrollment)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := MigrateFixedFileStoreV0ToV1(
		context.Background(), harness.path,
		V0ToV1MigrationOptions{Lock: heldInstallationOwnerMigrationLock{owner: owner}},
	); err != nil {
		t.Fatal(err)
	}

	wrongDomain, err := lifecyclestate.NewDomainIdentifier(
		lifecyclestate.DomainOwnerSessionID, bytes.Repeat([]byte{0xf1}, 16),
	)
	if err != nil {
		t.Fatal(err)
	}
	wrongSession, err := lifecyclestate.NewOwnerSessionID(wrongDomain)
	if err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(stateRootPath, "missing-store")
	if _, err := OpenFixedFileStoreV1WithOwner(
		context.Background(), missing, owner,
		FixedFileStoreV1Options{OwnerSessionID: wrongSession},
	); !errors.Is(err, ErrStoreOwnerSessionMismatch) {
		t.Fatalf("wrong owner session before missing-store inspection = %v", err)
	}

	store, err := OpenFixedFileStoreV1WithOwner(
		context.Background(), harness.path, owner,
		FixedFileStoreV1Options{OwnerSessionID: owner.OwnerSessionID()},
	)
	if err != nil {
		t.Fatal(err)
	}
	if store.OwnerSessionID() != owner.OwnerSessionID() {
		t.Fatal("owner-required store created a second owner session")
	}
	moved := lockPath + ".stale"
	if err := os.Rename(lockPath, moved); err != nil {
		t.Fatal(err)
	}
	replacement, err := os.OpenFile(lockPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := replacement.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecoveryAttemptIDs(context.Background()); !errors.Is(err, ErrStoreOwnerFenced) {
		t.Fatalf("replacement recovery fence = %v", err)
	}
	if !store.LifecycleRecoveryFenced() {
		t.Fatal("owner replacement did not permanently fence the store")
	}
	if err := store.CloseAfterLifecycleShutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecoveryAttemptIDs(context.Background()); !errors.Is(err, ErrStoreClosed) {
		t.Fatalf("closed store recovery = %v", err)
	}
	if err := owner.CloseAfterShutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := root.Close(); err != nil {
		t.Fatal(err)
	}
}
