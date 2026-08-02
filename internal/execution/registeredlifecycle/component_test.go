package registeredlifecycle

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"capsule.local/capsule/internal/execution/registrationstate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const ordinaryPlanPath = "../../../schemas/conformance/v0/execution-plan/ordinary.cbor"

func TestExecuteUsesTask4BRegistrationAndExposesNoSuccessResult(t *testing.T) {
	harness := newHarness(t, nil)
	if harness.backend.CreatesGuest() {
		t.Fatal("fake backend reported guest creation")
	}

	snapshot, err := harness.component.Execute(context.Background(), harness.registrationID)
	if err != nil {
		t.Fatalf("execute fake lifecycle: %v", err)
	}
	if snapshot.State != StateDestroyed || snapshot.CleanupRequired {
		t.Fatalf("lifecycle disposition = %+v, want destroyed with no cleanup", snapshot)
	}
	if snapshot.Failure != "" || snapshot.FailureAt != "" {
		t.Fatalf("ordinary fake lifecycle recorded a job result/failure: %+v", snapshot)
	}
	backendSnapshot := harness.backend.Snapshot(harness.registrationID)
	if !backendSnapshot.Prepared || !backendSnapshot.Created || !backendSnapshot.Started ||
		!backendSnapshot.Observed || !backendSnapshot.Stopped || !backendSnapshot.Destroyed {
		t.Fatalf("incomplete fake lifecycle: %+v", backendSnapshot)
	}

	method, ok := reflect.TypeOf((*Component)(nil)).MethodByName("Execute")
	if !ok {
		t.Fatal("Execute method missing")
	}
	if method.Type.NumIn() != 3 || method.Type.In(2) != reflect.TypeOf(v0candidate.RegistrationID{}) {
		t.Fatalf("Execute inputs = %v, want receiver, context, RegistrationID only", method.Type)
	}
	for _, prohibited := range []string{"Configure", "Command", "Import", "Launch", "Result", "Success"} {
		if _, exists := reflect.TypeOf(harness.backend).MethodByName(prohibited); exists {
			t.Fatalf("fake backend exposes prohibited method %q", prohibited)
		}
	}
}

func TestExecuteRejectsWrongTrustedBindingsBeforeLifecycleState(t *testing.T) {
	bindings := ordinaryBindings()
	bindings.RuntimeBundleManifestDigest[0] ^= 0xff
	harness := newHarness(t, &bindings)

	_, err := harness.component.Execute(context.Background(), harness.registrationID)
	classification, ok := ErrorClassification(err)
	if !ok || classification != ClassificationDomain {
		t.Fatalf("wrong binding classification = %q (%t), want DOMAIN: %v", classification, ok, err)
	}
	if _, err := harness.store.Snapshot(context.Background(), harness.registrationID); err == nil {
		t.Fatal("wrong trusted bindings created lifecycle state")
	}
	if calls := harness.backend.Snapshot(harness.registrationID).CallCounts; len(calls) != 0 {
		t.Fatalf("wrong trusted bindings reached fake backend: %v", calls)
	}
}

func TestConcurrentExecuteRunsOneLifecyclePerRegistration(t *testing.T) {
	harness := newHarness(t, nil)
	const workers = 32
	var group sync.WaitGroup
	errorsSeen := make(chan error, workers)
	for range workers {
		group.Add(1)
		go func() {
			defer group.Done()
			_, err := harness.component.Execute(context.Background(), harness.registrationID)
			errorsSeen <- err
		}()
	}
	group.Wait()
	close(errorsSeen)

	succeeded := 0
	stale := 0
	for err := range errorsSeen {
		if err == nil {
			succeeded++
			continue
		}
		classification, ok := ErrorClassification(err)
		if ok && classification == ClassificationStale {
			stale++
			continue
		}
		t.Fatalf("unexpected concurrent result: %v", err)
	}
	if succeeded != 1 || stale != workers-1 {
		t.Fatalf("concurrent results = %d complete, %d stale; want 1 and %d", succeeded, stale, workers-1)
	}
	backendSnapshot := harness.backend.Snapshot(harness.registrationID)
	for _, operation := range []Operation{
		OperationPrepare,
		OperationCreate,
		OperationStart,
		OperationObserve,
		OperationStop,
		OperationDestroy,
	} {
		if backendSnapshot.CallCounts[operation] != 1 {
			t.Fatalf("%s calls = %d, want 1", operation, backendSnapshot.CallCounts[operation])
		}
	}
}

func TestFakeFaultMatrixEndsDestroyedOrUnresolved(t *testing.T) {
	operations := []Operation{
		OperationPrepare,
		OperationCreate,
		OperationStart,
		OperationObserve,
		OperationStop,
		OperationDestroy,
	}
	for _, operation := range operations {
		for _, moment := range []FaultMoment{FaultBeforeEffect, FaultAfterEffect} {
			operation := operation
			moment := moment
			t.Run(string(operation)+"/"+string(moment), func(t *testing.T) {
				harness := newHarness(t, nil)
				if err := harness.backend.InjectFault(harness.registrationID, operation, moment); err != nil {
					t.Fatalf("inject fault: %v", err)
				}
				snapshot, err := harness.component.Execute(context.Background(), harness.registrationID)
				if err == nil {
					t.Fatal("injected lifecycle fault returned nil error")
				}
				classification, ok := ErrorClassification(err)
				if !ok || (classification != ClassificationLifecycleFailure &&
					classification != ClassificationCleanupUnresolved) {
					t.Fatalf("fault classification = %q (%t): %v", classification, ok, err)
				}
				expectedState := StateDestroyed
				if operation == OperationDestroy && moment == FaultBeforeEffect {
					expectedState = StateUnresolved
				}
				if snapshot.State != expectedState {
					t.Fatalf("fault disposition = %s, want %s", snapshot.State, expectedState)
				}
				if moment == FaultAfterEffect && snapshot.State != StateDestroyed {
					t.Fatalf("post-side-effect fault ended %s, want destroyed", snapshot.State)
				}
				if snapshot.Failure != ClassificationLifecycleFailure || snapshot.FailureAt != operation {
					t.Fatalf("fault record = %+v, want fixed lifecycle failure at %s", snapshot, operation)
				}
				if snapshot.State == StateDestroyed && snapshot.CleanupRequired {
					t.Fatal("destroyed fault retained cleanup obligation")
				}
				if snapshot.State == StateUnresolved && !snapshot.CleanupRequired {
					t.Fatal("unresolved fault lost cleanup obligation")
				}
			})
		}
	}
}

func TestPostSideEffectInterruptionRecoversAcrossComponentRestart(t *testing.T) {
	checkpoints := []Checkpoint{
		CheckpointAfterPrepareEffect,
		CheckpointAfterCreateEffect,
		CheckpointAfterStartEffect,
		CheckpointAfterObserveEffect,
		CheckpointAfterStopEffect,
		CheckpointAfterDestroyEffect,
	}
	checkpointOperations := map[Checkpoint]Operation{
		CheckpointAfterPrepareEffect: OperationPrepare,
		CheckpointAfterCreateEffect:  OperationCreate,
		CheckpointAfterStartEffect:   OperationStart,
		CheckpointAfterObserveEffect: OperationObserve,
		CheckpointAfterStopEffect:    OperationStop,
		CheckpointAfterDestroyEffect: OperationDestroy,
	}
	for _, checkpoint := range checkpoints {
		checkpoint := checkpoint
		t.Run(string(checkpoint), func(t *testing.T) {
			triggered := false
			harness := newHarnessWithCheckpoint(t, func(
				context.Context,
				Checkpoint,
				v0candidate.RegistrationID,
			) error {
				return nil
			})
			harness.component.checkpoint = func(
				_ context.Context,
				actual Checkpoint,
				_ v0candidate.RegistrationID,
			) error {
				if actual == checkpoint && !triggered {
					triggered = true
					return errors.New("simulated process loss with untrusted detail")
				}
				return nil
			}
			interrupted, err := harness.component.Execute(context.Background(), harness.registrationID)
			classification, ok := ErrorClassification(err)
			if !ok || classification != ClassificationLocalFailure {
				t.Fatalf("interruption classification = %q (%t), want LOCAL_FAILURE: %v", classification, ok, err)
			}
			if interrupted.State == StateDestroyed {
				t.Fatal("checkpoint did not leave recovery work")
			}

			restarted, err := New(Options{
				Registrations: harness.registrations,
				RoleBindings:  harness.bindingSource,
				Store:         harness.store,
				Backend:       harness.backend,
			})
			if err != nil {
				t.Fatalf("restart component: %v", err)
			}
			recovered, err := restarted.Recover(context.Background(), harness.registrationID)
			if err != nil {
				t.Fatalf("recover checkpoint %s: %v", checkpoint, err)
			}
			if recovered.State != StateDestroyed || recovered.CleanupRequired {
				t.Fatalf("recovered disposition = %+v", recovered)
			}
			if recovered.Failure != ClassificationLocalFailure ||
				recovered.FailureAt != checkpointOperations[checkpoint] {
				t.Fatalf("recovered fixed interruption evidence = %+v", recovered)
			}
		})
	}
}

func TestUnknownRecoveryRemainsUnresolvedAndRetryable(t *testing.T) {
	triggered := false
	harness := newHarnessWithCheckpoint(t, func(
		_ context.Context,
		checkpoint Checkpoint,
		_ v0candidate.RegistrationID,
	) error {
		if checkpoint == CheckpointAfterCreateEffect && !triggered {
			triggered = true
			return errors.New("simulated process loss")
		}
		return nil
	})
	if _, err := harness.component.Execute(context.Background(), harness.registrationID); err == nil {
		t.Fatal("expected simulated interruption")
	}
	if err := harness.backend.InjectFault(
		harness.registrationID,
		OperationReconcile,
		FaultBeforeEffect,
	); err != nil {
		t.Fatalf("inject reconcile fault: %v", err)
	}

	restarted, err := New(Options{
		Registrations: harness.registrations,
		RoleBindings:  harness.bindingSource,
		Store:         harness.store,
		Backend:       harness.backend,
	})
	if err != nil {
		t.Fatalf("restart component: %v", err)
	}
	unresolved, err := restarted.Recover(context.Background(), harness.registrationID)
	classification, ok := ErrorClassification(err)
	if !ok || classification != ClassificationCleanupUnresolved {
		t.Fatalf("unknown recovery classification = %q (%t): %v", classification, ok, err)
	}
	if unresolved.State != StateUnresolved || !unresolved.CleanupRequired {
		t.Fatalf("unknown recovery disposition = %+v", unresolved)
	}

	recovered, err := restarted.Recover(context.Background(), harness.registrationID)
	if err != nil {
		t.Fatalf("retry unresolved recovery: %v", err)
	}
	if recovered.State != StateDestroyed || recovered.CleanupRequired {
		t.Fatalf("retry disposition = %+v", recovered)
	}
}

func TestDefensiveSnapshotsAndFixedErrors(t *testing.T) {
	harness := newHarness(t, nil)
	if err := harness.backend.InjectFault(
		harness.registrationID,
		OperationObserve,
		FaultAfterEffect,
	); err != nil {
		t.Fatalf("inject fault: %v", err)
	}
	_, err := harness.component.Execute(context.Background(), harness.registrationID)
	if err == nil {
		t.Fatal("expected injected fault")
	}
	if bytes.Contains([]byte(err.Error()), harness.plan) ||
		bytes.Contains([]byte(err.Error()), []byte("simulated process loss with untrusted detail")) {
		t.Fatalf("fixed error leaked caller/backend content: %q", err)
	}

	first := harness.backend.Snapshot(harness.registrationID)
	first.CallCounts[OperationPrepare] = 999
	second := harness.backend.Snapshot(harness.registrationID)
	if second.CallCounts[OperationPrepare] != 1 {
		t.Fatal("backend snapshot did not defensively copy call counts")
	}
	lifecycle, err := harness.store.Snapshot(context.Background(), harness.registrationID)
	if err != nil {
		t.Fatalf("lifecycle snapshot: %v", err)
	}
	mutated := lifecycle
	mutated.State = StateUnresolved
	again, err := harness.store.Snapshot(context.Background(), harness.registrationID)
	if err != nil {
		t.Fatalf("lifecycle snapshot again: %v", err)
	}
	if again.State != StateDestroyed {
		t.Fatal("caller mutation changed retained lifecycle state")
	}
}

func TestLifecycleStoreRefusesBeyondFixedCapacity(t *testing.T) {
	store := NewMemoryStore()
	digest := repeated32[v0candidate.ExecutionPlanDigest](0x42)
	for index := 1; index <= MaxLifecycleRecords; index++ {
		var registrationID v0candidate.RegistrationID
		registrationID[0] = byte(index >> 8)
		registrationID[1] = byte(index)
		if err := store.begin(context.Background(), registrationID, digest); err != nil {
			t.Fatalf("begin record %d: %v", index, err)
		}
	}
	var overflowID v0candidate.RegistrationID
	overflowID[0] = 0x7f
	if err := store.begin(context.Background(), overflowID, digest); err == nil {
		t.Fatal("fixed lifecycle capacity accepted another record")
	} else if classification, ok := ErrorClassification(err); !ok || classification != ClassificationCapacity {
		t.Fatalf("capacity classification = %q (%t), want CAPACITY: %v", classification, ok, err)
	}
}

type testHarness struct {
	component      *Component
	registrations  *registrationstate.Component
	bindingSource  *fixedBindingSource
	store          *MemoryStore
	backend        *FakeBackend
	registrationID v0candidate.RegistrationID
	plan           []byte
}

func newHarness(
	t *testing.T,
	overrideBindings *v0candidate.ExecutionPlanRoleBindings,
) testHarness {
	t.Helper()
	return newHarnessConfigured(t, overrideBindings, nil)
}

func newHarnessWithCheckpoint(t *testing.T, checkpoint CheckpointHook) testHarness {
	t.Helper()
	return newHarnessConfigured(t, nil, checkpoint)
}

func newHarnessConfigured(
	t *testing.T,
	overrideBindings *v0candidate.ExecutionPlanRoleBindings,
	checkpoint CheckpointHook,
) testHarness {
	t.Helper()
	plan, err := os.ReadFile(filepath.Clean(ordinaryPlanPath))
	if err != nil {
		t.Fatalf("read ordinary plan: %v", err)
	}
	registrationStore, err := registrationstate.NewFixedFileStore(
		filepath.Join(t.TempDir(), "registration-state.json"),
		registrationstate.InitialState{
			InstallationID:           repeated16[v0candidate.InstallationID](0x11),
			SupervisorID:             repeated16[v0candidate.SupervisorID](0x55),
			EpochSequence:            7,
			EpochDigest:              repeated32[v0candidate.TrustEpochDigest](0x22),
			TrustPhase:               registrationstate.TrustStable,
			TimeHighWaterUnixSeconds: 1_785_456_000,
		},
	)
	if err != nil {
		t.Fatalf("new Task 4B store: %v", err)
	}
	registrations, err := registrationstate.New(registrationstate.Options{
		Store:       registrationStore,
		Clock:       fixedClock(1_785_456_000),
		Identifiers: fixedRegistrationIDSource{id: repeated16[v0candidate.RegistrationID](0x33)},
	})
	if err != nil {
		t.Fatalf("new Task 4B component: %v", err)
	}
	bindings := ordinaryBindings()
	issued, err := registrations.RegisterPlan(
		context.Background(),
		registrationstate.AuthenticatedCallContext{
			Authenticated: true,
			Role:          registrationstate.CallerDaemon,
			Purpose:       registrationstate.RegisterPlanPurpose,
		},
		plan,
		bindings,
	)
	if err != nil {
		t.Fatalf("register ordinary plan: %v", err)
	}
	registrationID := issued.View().RegistrationID
	if overrideBindings != nil {
		bindings = *overrideBindings
	}
	bindingSource := &fixedBindingSource{
		registrationID: registrationID,
		bindings:       cloneRoleBindings(bindings),
	}
	store := NewMemoryStore()
	backend := NewFakeBackend()
	component, err := New(Options{
		Registrations: registrations,
		RoleBindings:  bindingSource,
		Store:         store,
		Backend:       backend,
		Checkpoint:    checkpoint,
	})
	if err != nil {
		t.Fatalf("new registered lifecycle: %v", err)
	}
	return testHarness{
		component:      component,
		registrations:  registrations,
		bindingSource:  bindingSource,
		store:          store,
		backend:        backend,
		registrationID: registrationID,
		plan:           plan,
	}
}

type fixedClock uint64

func (clock fixedClock) ObserveUnixSeconds(context.Context) (uint64, error) {
	return uint64(clock), nil
}

type fixedRegistrationIDSource struct {
	id v0candidate.RegistrationID
}

func (source fixedRegistrationIDSource) NewRegistrationID(
	context.Context,
) (v0candidate.RegistrationID, error) {
	return source.id, nil
}

type fixedBindingSource struct {
	registrationID v0candidate.RegistrationID
	bindings       v0candidate.ExecutionPlanRoleBindings
}

func (source *fixedBindingSource) ResolveExecutionPlanRoleBindings(
	_ context.Context,
	registrationID v0candidate.RegistrationID,
) (v0candidate.ExecutionPlanRoleBindings, error) {
	if registrationID != source.registrationID {
		return v0candidate.ExecutionPlanRoleBindings{}, errors.New("unknown registration")
	}
	return cloneRoleBindings(source.bindings), nil
}

func ordinaryBindings() v0candidate.ExecutionPlanRoleBindings {
	return v0candidate.ExecutionPlanRoleBindings{
		InstallationID:                  repeated16[v0candidate.InstallationID](0x11),
		EpochDigest:                     repeated32[v0candidate.TrustEpochDigest](0x22),
		SourceManifestDigest:            hex32[v0candidate.SourceManifestDigest]("e5e09b2435baedf897526a89c698c0b0531437a69472372ae426f62d801fc171"),
		InlineInputDigest:               hex32[v0candidate.InlineInputDigest]("bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72"),
		RuntimeBundleManifestDigest:     repeated32[v0candidate.RuntimeBundleManifestDigest](0x55),
		ProfileReviewAttestationDigests: []v0candidate.ProfileReviewAttestationDigest{repeated32[v0candidate.ProfileReviewAttestationDigest](0x66), repeated32[v0candidate.ProfileReviewAttestationDigest](0x67)},
		ProfileRegistryEntryDigest:      repeated32[v0candidate.ProfileRegistryEntryDigest](0x77),
		BackendValidationRecordDigest:   repeated32[v0candidate.BackendValidationRecordDigest](0x88),
		BackendConfigurationDigest:      repeated32[v0candidate.BackendConfigurationDigest](0x99),
		TrustSnapshotDigest:             repeated32[v0candidate.TrustSnapshotDigest](0xaa),
		PolicyDecisionDigest:            repeated32[v0candidate.PolicyDecisionDigest](0xbb),
	}
}

func repeated16[T ~[16]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}

func repeated32[T ~[32]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}

func hex32[T ~[32]byte](value string) T {
	decoded, err := hex.DecodeString(value)
	if err != nil || len(decoded) != 32 {
		panic("invalid test digest")
	}
	var result T
	copy(result[:], decoded)
	return result
}
