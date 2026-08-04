//go:build darwin

package registeredlifecycle

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
	"reflect"
	"strings"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/installationowner"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/execution/registrationstate"
	"capsule.local/capsule/internal/protocol/v0candidate"
	"golang.org/x/sys/unix"
)

const (
	ownedStartupHelperMode       = "CAPSULE_OWNED_STARTUP_HELPER_MODE"
	ownedStartupHelperRoot       = "CAPSULE_OWNED_STARTUP_HELPER_ROOT"
	ownedStartupHelperStore      = "CAPSULE_OWNED_STARTUP_HELPER_STORE"
	ownedStartupHelperEnrollment = "CAPSULE_OWNED_STARTUP_HELPER_ENROLLMENT"
)

type ownedStartupFixture struct {
	rootPath   string
	lockPath   string
	enrollment installationowner.OwnerLockEnrollment
}

type startupTraceEntry struct {
	checkpoint ownedStartupCheckpoint
	attemptID  approvalattempt.AttemptID
}

func TestOwnerRequiredStartupOrdersOwnershipStoreRecoveryAndClose(t *testing.T) {
	harness := newHarness(t, []fixtureSpec{{nonce: 0x66}, {nonce: 0x67, variant: 0x01}})
	fixture := newOwnedStartupFixture(t, filepath.Dir(harness.path))
	var trace []startupTraceEntry
	startup, err := openOwnedStartup(context.Background(), OwnedStartupOptions{
		StateRootPath: fixture.rootPath, Enrollment: fixture.enrollment, StorePath: harness.path,
		Attempts: harness.attempts, Backend: harness.backend, Clock: harness.clock, EffectIDs: harness.effectIDs,
	}, ownedStartupTestHooks{startup: func(_ context.Context, checkpoint ownedStartupCheckpoint, attemptID approvalattempt.AttemptID) error {
		trace = append(trace, startupTraceEntry{checkpoint: checkpoint, attemptID: attemptID})
		return nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	if startup.OwnerSessionID().IsZero() || startup.OwnerSessionID() != startup.store.OwnerSessionID() ||
		startup.OwnerSessionID() != startup.component.coordinator.OwnerSessionID() {
		t.Fatal("owner, store, and coordinator did not share one live owner session")
	}
	recovered := startup.InitialRecovery()
	if len(recovered) != len(harness.attemptIDs) {
		t.Fatalf("initial recovery count = %d, want %d", len(recovered), len(harness.attemptIDs))
	}
	for index, snapshot := range recovered {
		if snapshot.AttemptID != harness.attemptIDs[index] || snapshot.State != StateDestroyed || snapshot.CleanupRequired {
			t.Fatalf("initial recovery %d = %+v", index, snapshot)
		}
	}
	beforeRecover := make([]approvalattempt.AttemptID, 0, len(harness.attemptIDs))
	for _, entry := range trace {
		if entry.checkpoint == ownedStartupBeforeRecover {
			beforeRecover = append(beforeRecover, entry.attemptID)
		}
	}
	if !sortedUniqueAttemptIDs(beforeRecover) || !reflect.DeepEqual(beforeRecover, harness.attemptIDs) {
		t.Fatalf("recovery order = %x, want %x", beforeRecover, harness.attemptIDs)
	}
	wantPrefix := []ownedStartupCheckpoint{
		ownedStartupRootOpened,
		ownedStartupOwnerAcquired,
		ownedStartupOwnerChecked,
		ownedStartupStoreOpened,
		ownedStartupCoordinatorCreated,
		ownedStartupLifecycleCreated,
		ownedStartupRecoveryOwnerHeld,
		ownedStartupRecoveryEnumerated,
	}
	if len(trace) < len(wantPrefix) {
		t.Fatalf("startup trace too short: %+v", trace)
	}
	for index, checkpoint := range wantPrefix {
		if trace[index].checkpoint != checkpoint {
			t.Fatalf("startup trace %d = %s, want %s", index, trace[index].checkpoint, checkpoint)
		}
	}
	for _, attemptID := range harness.attemptIDs {
		backend := harness.backend.Snapshot(attemptID)
		for _, operation := range durableOperations {
			if backend.CallCounts[operation] != 1 || backend.ApplicationCounts[operation] != 1 {
				t.Fatalf("%x %s calls/applications = %d/%d", attemptID, operation, backend.CallCounts[operation], backend.ApplicationCounts[operation])
			}
		}
	}
	if harness.backend.CreatesGuest() {
		t.Fatal("owned startup fake backend creates a guest")
	}
	firstSession := startup.OwnerSessionID()
	if err := startup.CloseAfterShutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	wantClose := []ownedStartupCheckpoint{
		ownedStartupLifecycleStopped, ownedStartupStoreClosed,
		ownedStartupOwnerClosed, ownedStartupStateRootClosed,
	}
	if len(trace) < len(wantClose) {
		t.Fatalf("close trace too short: %+v", trace)
	}
	for index, checkpoint := range wantClose {
		actual := trace[len(trace)-len(wantClose)+index].checkpoint
		if actual != checkpoint {
			t.Fatalf("close trace %d = %s, want %s", index, actual, checkpoint)
		}
	}

	sessions := map[lifecyclestate.OwnerSessionID]struct{}{firstSession: {}}
	for repetition := range 20 {
		reopened, err := OpenOwnedStartup(context.Background(), OwnedStartupOptions{
			StateRootPath: fixture.rootPath, Enrollment: fixture.enrollment, StorePath: harness.path,
			Attempts: harness.attempts, Backend: harness.backend, Clock: harness.clock, EffectIDs: harness.effectIDs,
		})
		if err != nil {
			t.Fatalf("reopen %d: %v", repetition, err)
		}
		if len(reopened.InitialRecovery()) != 0 {
			t.Fatalf("terminal reopen %d recovery = %+v", repetition, reopened.InitialRecovery())
		}
		session := reopened.OwnerSessionID()
		if session.IsZero() {
			t.Fatalf("reopen %d produced zero owner session", repetition)
		}
		if _, duplicate := sessions[session]; duplicate {
			t.Fatalf("reopen %d repeated an owner session", repetition)
		}
		sessions[session] = struct{}{}
		if err := reopened.CloseAfterShutdown(context.Background()); err != nil {
			t.Fatalf("close reopen %d: %v", repetition, err)
		}
	}
}

func TestDuplicateOwnerAndWrongEnrollmentRefuseBeforeStoreAndFakeWork(t *testing.T) {
	t.Run("duplicate beats corrupt store", func(t *testing.T) {
		harness := newHarness(t, nil)
		fixture := newOwnedStartupFixture(t, filepath.Dir(harness.path))
		root, err := installationowner.OpenDarwinTrustedStateRoot(context.Background(), fixture.rootPath, fixture.enrollment)
		if err != nil {
			t.Fatal(err)
		}
		owner, err := installationowner.AcquireDarwinInstallationOwner(context.Background(), root, fixture.enrollment)
		if err != nil {
			t.Fatal(err)
		}
		corrupt := []byte("corrupt-before-duplicate\n")
		if err := os.WriteFile(harness.path, corrupt, 0o600); err != nil {
			t.Fatal(err)
		}
		startup, err := OpenOwnedStartup(context.Background(), OwnedStartupOptions{
			StateRootPath: fixture.rootPath, Enrollment: fixture.enrollment, StorePath: harness.path,
			Attempts: harness.attempts, Backend: harness.backend, Clock: harness.clock, EffectIDs: harness.effectIDs,
		})
		if startup != nil || !errors.Is(err, installationowner.ErrDuplicateOwner) {
			t.Fatalf("duplicate startup = %T, %v", startup, err)
		}
		assertZeroFakeCalls(t, harness.backend, harness.attemptIDs)
		if err := owner.CloseAfterShutdown(context.Background()); err != nil {
			t.Fatal(err)
		}
		if err := root.Close(); err != nil {
			t.Fatal(err)
		}
		if startup, err := OpenOwnedStartup(context.Background(), OwnedStartupOptions{
			StateRootPath: fixture.rootPath, Enrollment: fixture.enrollment, StorePath: harness.path,
			Attempts: harness.attempts, Backend: harness.backend, Clock: harness.clock, EffectIDs: harness.effectIDs,
		}); startup != nil || !errors.Is(err, registrationstate.ErrStoreRepairRequired) {
			t.Fatalf("corrupt reopen = %T, %v", startup, err)
		}
		assertZeroFakeCalls(t, harness.backend, harness.attemptIDs)
		if retained, readErr := os.ReadFile(harness.path); readErr != nil || !bytes.Equal(retained, corrupt) {
			t.Fatalf("corrupt refusal changed store bytes: %q, %v", retained, readErr)
		}
	})

	t.Run("wrong enrolled inode beats missing store", func(t *testing.T) {
		harness := newHarness(t, nil)
		fixture := newOwnedStartupFixture(t, filepath.Dir(harness.path))
		view := fixture.enrollment.View()
		view.LockInode++
		wrong, err := installationowner.NewOwnerLockEnrollment(view)
		if err != nil {
			t.Fatal(err)
		}
		startup, err := OpenOwnedStartup(context.Background(), OwnedStartupOptions{
			StateRootPath: fixture.rootPath, Enrollment: wrong,
			StorePath: filepath.Join(fixture.rootPath, "missing-store"),
			Attempts:  harness.attempts, Backend: harness.backend, Clock: harness.clock, EffectIDs: harness.effectIDs,
		})
		if startup != nil || !errors.Is(err, installationowner.ErrRepairRequired) {
			t.Fatalf("wrong enrollment startup = %T, %v", startup, err)
		}
		assertZeroFakeCalls(t, harness.backend, harness.attemptIDs)
	})
}

func TestOwnerReplacementAfterOpenFencesRecoveryUntilShutdownReopen(t *testing.T) {
	harness := newHarness(t, nil)
	fixture := newOwnedStartupFixture(t, filepath.Dir(harness.path))
	startup, err := OpenOwnedStartup(context.Background(), OwnedStartupOptions{
		StateRootPath: fixture.rootPath, Enrollment: fixture.enrollment, StorePath: harness.path,
		Attempts: harness.attempts, Backend: harness.backend, Clock: harness.clock, EffectIDs: harness.effectIDs,
	})
	if err != nil {
		t.Fatal(err)
	}
	before := harness.backend.Snapshot(harness.attemptID)
	moved := fixture.lockPath + ".replaced"
	if err := os.Rename(fixture.lockPath, moved); err != nil {
		t.Fatal(err)
	}
	replacement, err := os.OpenFile(fixture.lockPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := replacement.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := startup.RecoverCreatedAttempts(context.Background()); !errors.Is(err, registrationstate.ErrStoreOwnerFenced) {
		t.Fatalf("post-open replacement recovery = %v", err)
	}
	if !startup.store.LifecycleRecoveryFenced() {
		t.Fatal("post-open owner failure did not fence the store")
	}
	if err := startup.store.AdvanceLifecycleTime(context.Background(), 1_785_456_001); !errors.Is(err, registrationstate.ErrStoreOwnerFenced) {
		t.Fatalf("post-open replacement mutation = %v", err)
	}
	if _, err := startup.store.ReadLifecycle(context.Background(), harness.attemptID); !errors.Is(err, registrationstate.ErrStoreOwnerFenced) {
		t.Fatalf("post-open replacement read = %v", err)
	}
	after := harness.backend.Snapshot(harness.attemptID)
	if !reflect.DeepEqual(before.CallCounts, after.CallCounts) || !reflect.DeepEqual(before.ApplicationCounts, after.ApplicationCounts) {
		t.Fatalf("owner loss reached fake backend: before=%+v after=%+v", before, after)
	}
	if err := startup.CloseAfterShutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if reopened, err := OpenOwnedStartup(context.Background(), OwnedStartupOptions{
		StateRootPath: fixture.rootPath, Enrollment: fixture.enrollment, StorePath: harness.path,
		Attempts: harness.attempts, Backend: harness.backend, Clock: harness.clock, EffectIDs: harness.effectIDs,
	}); reopened != nil || !errors.Is(err, installationowner.ErrRepairRequired) {
		t.Fatalf("replacement reopen = %T, %v", reopened, err)
	}
}

func TestOwnedStartupRecoveryFaultClosesOwnerAndReopensSameFakeEffect(t *testing.T) {
	harness := newHarness(t, nil)
	fixture := newOwnedStartupFixture(t, filepath.Dir(harness.path))
	startup, err := openOwnedStartup(context.Background(), OwnedStartupOptions{
		StateRootPath: fixture.rootPath, Enrollment: fixture.enrollment, StorePath: harness.path,
		Attempts: harness.attempts, Backend: harness.backend, Clock: harness.clock, EffectIDs: harness.effectIDs,
	}, ownedStartupTestHooks{lifecycle: oneShotCheckpoint(CheckpointAfterCreateEffect)})
	if startup != nil || err == nil {
		t.Fatalf("interrupted startup = %T, %v", startup, err)
	}
	interrupted := harness.backend.Snapshot(harness.attemptID)
	if interrupted.ApplicationCounts[OperationCreate] != 1 || interrupted.CallCounts[OperationCreate] != 1 {
		t.Fatalf("interrupted create = %+v", interrupted)
	}
	reopened, err := OpenOwnedStartup(context.Background(), OwnedStartupOptions{
		StateRootPath: fixture.rootPath, Enrollment: fixture.enrollment, StorePath: harness.path,
		Attempts: harness.attempts, Backend: harness.backend, Clock: harness.clock, EffectIDs: harness.effectIDs,
	})
	if err != nil {
		t.Fatal(err)
	}
	if recovered := reopened.InitialRecovery(); len(recovered) != 1 || recovered[0].State != StateDestroyed {
		t.Fatalf("reopened recovery = %+v", recovered)
	}
	after := harness.backend.Snapshot(harness.attemptID)
	if after.ApplicationCounts[OperationCreate] != 1 ||
		!reflect.DeepEqual(interrupted.EffectIDs[OperationCreate], after.EffectIDs[OperationCreate]) {
		t.Fatalf("response-loss recovery changed create identity: before=%+v after=%+v", interrupted, after)
	}
	if err := reopened.CloseAfterShutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestOwnedStartupProcessDeathAllowsFullReopen(t *testing.T) {
	rootPath := filepath.Join(t.TempDir(), "owned-startup-process")
	if err := os.Mkdir(rootPath, 0o700); err != nil {
		t.Fatal(err)
	}
	storePath := filepath.Join(rootPath, "supervisor-state.json")
	store, err := registrationstate.NewFixedFileStore(storePath, registrationstate.InitialState{
		InstallationID: repeated16[v0candidate.InstallationID](0x11),
		SupervisorID:   repeated16[v0candidate.SupervisorID](0x55), EpochSequence: 7,
		EpochDigest: repeated32[v0candidate.TrustEpochDigest](0x22), TrustPhase: registrationstate.TrustStable,
		TimeHighWaterUnixSeconds: 1_785_456_000,
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = store
	if _, err := registrationstate.MigrateFixedFileStoreV0ToV1(
		context.Background(), storePath,
		registrationstate.V0ToV1MigrationOptions{Lock: startupOfflineLock{}},
	); err != nil {
		t.Fatal(err)
	}
	fixture := newOwnedStartupFixture(t, rootPath)
	command, input := startOwnedStartupHelper(t, fixture, storePath)
	t.Cleanup(func() {
		_ = input.Close()
		_ = command.Process.Kill()
		_ = command.Wait()
	})
	backend := ownedStartupFakeBackend(t)
	startup, err := OpenOwnedStartup(context.Background(), OwnedStartupOptions{
		StateRootPath: rootPath, Enrollment: fixture.enrollment, StorePath: storePath,
		Attempts: emptyAttemptResolver{}, Backend: backend, Clock: &testClock{value: 1_785_456_000},
	})
	if startup != nil || !errors.Is(err, installationowner.ErrDuplicateOwner) {
		t.Fatalf("live child duplicate = %T, %v", startup, err)
	}
	if err := command.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = input.Close()
	if err := command.Wait(); err == nil {
		t.Fatal("killed child exited successfully")
	}
	reopened, err := OpenOwnedStartup(context.Background(), OwnedStartupOptions{
		StateRootPath: rootPath, Enrollment: fixture.enrollment, StorePath: storePath,
		Attempts: emptyAttemptResolver{}, Backend: backend, Clock: &testClock{value: 1_785_456_000},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(reopened.InitialRecovery()) != 0 || backend.CreatesGuest() {
		t.Fatalf("process-death reopen recovery=%+v createsGuest=%t", reopened.InitialRecovery(), backend.CreatesGuest())
	}
	if err := reopened.CloseAfterShutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestOwnedStartupProcessHelper(t *testing.T) {
	if os.Getenv(ownedStartupHelperMode) != "hold" {
		return
	}
	var view installationowner.OwnerLockEnrollmentView
	if err := json.Unmarshal([]byte(os.Getenv(ownedStartupHelperEnrollment)), &view); err != nil {
		t.Fatalf("decode enrollment: %v", err)
	}
	enrollment, err := installationowner.NewOwnerLockEnrollment(view)
	if err != nil {
		t.Fatal(err)
	}
	startup, err := OpenOwnedStartup(context.Background(), OwnedStartupOptions{
		StateRootPath: os.Getenv(ownedStartupHelperRoot), Enrollment: enrollment,
		StorePath: os.Getenv(ownedStartupHelperStore), Attempts: emptyAttemptResolver{},
		Backend: ownedStartupFakeBackend(t), Clock: &testClock{value: 1_785_456_000},
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("READY")
	_, _ = io.Copy(io.Discard, os.Stdin)
	_ = startup.CloseAfterShutdown(context.Background())
}

type startupOfflineLock struct{}

func (startupOfflineLock) CheckOfflineMigrationLock(context.Context) error { return nil }

type emptyAttemptResolver struct{}

func (emptyAttemptResolver) ResolveCreated(context.Context, approvalattempt.AttemptID) (registrationstate.CreatedAttempt, error) {
	return registrationstate.CreatedAttempt{}, errors.New("empty startup resolver called")
}

func newOwnedStartupFixture(t *testing.T, rootPath string) ownedStartupFixture {
	t.Helper()
	lockPath := filepath.Join(rootPath, "supervisor.owner")
	file, err := os.OpenFile(lockPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	var rootStat unix.Stat_t
	if err := unix.Stat(rootPath, &rootStat); err != nil {
		t.Fatal(err)
	}
	var lockStat unix.Stat_t
	if err := unix.Lstat(lockPath, &lockStat); err != nil {
		t.Fatal(err)
	}
	enrollment, err := installationowner.NewOwnerLockEnrollment(installationowner.OwnerLockEnrollmentView{
		InstallationID: repeated16[v0candidate.InstallationID](0x11),
		SupervisorID:   repeated16[v0candidate.SupervisorID](0x55), ExpectedUID: uint32(unix.Geteuid()), //nolint:gosec // Local Darwin test UID is nonnegative.
		StateRootDevice: uint64(uint32(rootStat.Dev)), StateRootInode: rootStat.Ino, //nolint:gosec // Preserve Darwin dev_t bits.
		LockEntryName: filepath.Base(lockPath), LockDevice: uint64(uint32(lockStat.Dev)), LockInode: lockStat.Ino, //nolint:gosec // Preserve Darwin dev_t bits.
	})
	if err != nil {
		t.Fatal(err)
	}
	return ownedStartupFixture{rootPath: rootPath, lockPath: lockPath, enrollment: enrollment}
}

func ownedStartupFakeBackend(t *testing.T) *FakeBackend {
	t.Helper()
	bindings := ordinaryBindings()
	binding, err := lifecyclestate.NewBackendBinding(lifecyclestate.BackendBindingView{
		Kind: lifecyclestate.BackendFakeNoGuest, ProtocolVersion: lifecyclestate.FakeBackendProtocolVersion,
		ImplementationIdentityDigest:  lifecyclestate.BackendImplementationDigest(sha256.Sum256([]byte("owned-startup-fake-no-guest"))),
		BackendConfigurationDigest:    bindings.BackendConfigurationDigest,
		BackendValidationRecordDigest: bindings.BackendValidationRecordDigest,
		CreatesGuest:                  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	backend, err := NewFakeBackend(binding)
	if err != nil {
		t.Fatal(err)
	}
	return backend
}

func assertZeroFakeCalls(t *testing.T, backend *FakeBackend, attemptIDs []approvalattempt.AttemptID) {
	t.Helper()
	for _, attemptID := range attemptIDs {
		snapshot := backend.Snapshot(attemptID)
		for operation, count := range snapshot.CallCounts {
			if count != 0 {
				t.Fatalf("%x %s fake call count = %d", attemptID, operation, count)
			}
		}
		for operation, count := range snapshot.ApplicationCounts {
			if count != 0 {
				t.Fatalf("%x %s fake application count = %d", attemptID, operation, count)
			}
		}
	}
	if backend.CreatesGuest() {
		t.Fatal("fake backend creates a guest")
	}
}

func startOwnedStartupHelper(
	t *testing.T,
	fixture ownedStartupFixture,
	storePath string,
) (*exec.Cmd, io.WriteCloser) {
	t.Helper()
	encoded, err := json.Marshal(fixture.enrollment.View())
	if err != nil {
		t.Fatal(err)
	}
	command := exec.Command(os.Args[0], "-test.run=^TestOwnedStartupProcessHelper$") //nolint:gosec // Exact current test binary and fixed test selector.
	command.Env = append(os.Environ(),
		ownedStartupHelperMode+"=hold",
		ownedStartupHelperRoot+"="+fixture.rootPath,
		ownedStartupHelperStore+"="+storePath,
		ownedStartupHelperEnrollment+"="+string(encoded),
	)
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
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "READY" {
		_ = command.Process.Kill()
		_ = command.Wait()
		t.Fatalf("owned startup helper readiness = %q, stderr=%q", scanner.Text(), stderr.String())
	}
	return command, input
}
