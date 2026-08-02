package registeredlifecycle

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/registrationstate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const ordinaryPlanPath = "../../../schemas/conformance/v0/execution-plan/ordinary.cbor"

func TestDriveResolvesCommittedSliceBAttemptAndExposesNoSuccessResult(t *testing.T) {
	harness := newHarness(t, nil)
	if harness.backend.CreatesGuest() {
		t.Fatal("fake backend reported guest creation")
	}

	snapshot, err := harness.component.Drive(context.Background(), harness.attemptID)
	if err != nil {
		t.Fatalf("drive fake lifecycle: %v", err)
	}
	if snapshot.AttemptID != harness.attemptID || snapshot.RegistrationID != harness.registrationID ||
		snapshot.State != StateDestroyed || snapshot.CleanupRequired {
		t.Fatalf("lifecycle disposition = %+v, want bound destroyed attempt", snapshot)
	}
	if snapshot.Failure != "" || snapshot.FailureAt != "" {
		t.Fatalf("ordinary fake lifecycle recorded a job result/failure: %+v", snapshot)
	}
	backendSnapshot := harness.backend.Snapshot(harness.attemptID)
	if !backendSnapshot.Prepared || !backendSnapshot.Created || !backendSnapshot.Started ||
		!backendSnapshot.Observed || !backendSnapshot.Stopped || !backendSnapshot.Destroyed {
		t.Fatalf("incomplete fake lifecycle: %+v", backendSnapshot)
	}

	componentType := reflect.TypeOf((*Component)(nil))
	method, ok := componentType.MethodByName("Drive")
	if !ok || method.Type.NumIn() != 3 || method.Type.In(2) != reflect.TypeOf(approvalattempt.AttemptID{}) {
		t.Fatalf("Drive inputs = %v, want receiver, context, AttemptID only", method.Type)
	}
	if _, exists := componentType.MethodByName("Execute"); exists {
		t.Fatal("obsolete RegistrationID-keyed Execute method remains")
	}
	for _, prohibited := range []string{"Configure", "Command", "Import", "Launch", "Result", "Success"} {
		if _, exists := reflect.TypeOf(harness.backend).MethodByName(prohibited); exists {
			t.Fatalf("fake backend exposes prohibited method %q", prohibited)
		}
	}
}

func TestDriveRejectsMissingMutatedAndCrossLinkedAttemptsBeforePrepare(t *testing.T) {
	tests := []struct {
		name  string
		alter func(*registrationstate.CreatedAttempt)
		want  Classification
	}{
		{
			name: "wrong plan role binding",
			alter: func(created *registrationstate.CreatedAttempt) {
				created.PlanRoleBindings.RuntimeBundleManifestDigest[0] ^= 0xff
			},
			want: ClassificationDomain,
		},
		{
			name: "wrong attempt identity",
			alter: func(created *registrationstate.CreatedAttempt) {
				created.Attempt.AttemptID[0] ^= 0xff
			},
			want: ClassificationBinding,
		},
		{
			name: "non created attempt state",
			alter: func(created *registrationstate.CreatedAttempt) {
				created.Attempt.State = ""
			},
			want: ClassificationBinding,
		},
		{
			name: "cross linked approval",
			alter: func(created *registrationstate.CreatedAttempt) {
				created.Approval.ApprovalID[0] ^= 0xff
			},
			want: ClassificationBinding,
		},
		{
			name: "mutated exact plan",
			alter: func(created *registrationstate.CreatedAttempt) {
				created.Registration.ExactPlanBytes[len(created.Registration.ExactPlanBytes)-1] ^= 0xff
			},
			want: ClassificationBinding,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			harness := newHarness(t, nil)
			resolver := &alteringResolver{base: harness.attempts, alter: test.alter}
			component := mustLifecycle(t, resolver, harness.attempts, harness.store, harness.backend, nil)
			_, err := component.Drive(context.Background(), harness.attemptID)
			assertLifecycleClassification(t, err, test.want)
			if _, err := harness.store.Snapshot(context.Background(), harness.attemptID); err == nil {
				t.Fatal("invalid resolved attempt created lifecycle state")
			}
			if calls := harness.backend.Snapshot(harness.attemptID).CallCounts; len(calls) != 0 {
				t.Fatalf("invalid resolved attempt reached fake backend: %v", calls)
			}
		})
	}

	harness := newHarness(t, nil)
	unknown := harness.attemptID
	unknown[15] ^= 0xff
	_, err := harness.component.Drive(context.Background(), unknown)
	assertLifecycleClassification(t, err, ClassificationBinding)
	if calls := harness.backend.Snapshot(unknown).CallCounts; len(calls) != 0 {
		t.Fatalf("missing attempt reached fake backend: %v", calls)
	}
	_, err = harness.component.Drive(context.Background(), approvalattempt.AttemptID{})
	assertLifecycleClassification(t, err, ClassificationBinding)
}

func TestRecoveryFencedAttemptStoreRefusesBeforePrepare(t *testing.T) {
	harness := newHarness(t, []fixtureSpec{
		{nonce: 0x66, variant: 0},
		{nonce: 0x67, variant: 1},
	})
	harness.authorityStore.InjectFailure(
		registrationstate.FaultTimeHighWaterIndeterminatePreState,
		errors.New("simulated indeterminate time commit"),
	)
	_, err := harness.attempts.SubmitApproval(
		context.Background(), registrationstate.AuthenticatedCallContext{
			Authenticated: true, Role: registrationstate.CallerBroker,
			Purpose: registrationstate.SubmitApprovalPurpose,
		}, harness.registrationID, harness.vectors[1].EnvelopeBytes,
	)
	if err == nil {
		t.Fatal("indeterminate write did not fence the authority store")
	}
	_, err = harness.component.Drive(context.Background(), harness.attemptID)
	assertLifecycleClassification(t, err, ClassificationRecoveryRequired)
	if calls := harness.backend.Snapshot(harness.attemptID).CallCounts; len(calls) != 0 {
		t.Fatalf("recovery-fenced attempt reached fake backend: %v", calls)
	}
}

func TestConcurrentAndSequentialDriveNeverRedriveEffects(t *testing.T) {
	harness := newHarness(t, nil)
	const workers = 32
	var group sync.WaitGroup
	errorsSeen := make(chan error, workers)
	for range workers {
		group.Add(1)
		go func() {
			defer group.Done()
			_, err := harness.component.Drive(context.Background(), harness.attemptID)
			errorsSeen <- err
		}()
	}
	group.Wait()
	close(errorsSeen)
	for err := range errorsSeen {
		if err != nil {
			t.Fatalf("concurrent drive: %v", err)
		}
	}
	if _, err := harness.component.Drive(context.Background(), harness.attemptID); err != nil {
		t.Fatalf("sequential replay: %v", err)
	}
	backendSnapshot := harness.backend.Snapshot(harness.attemptID)
	for _, operation := range []Operation{
		OperationPrepare, OperationCreate, OperationStart,
		OperationObserve, OperationStop, OperationDestroy,
	} {
		if backendSnapshot.CallCounts[operation] != 1 {
			t.Fatalf("%s calls = %d, want exactly 1", operation, backendSnapshot.CallCounts[operation])
		}
	}
}

func TestTwoApprovalsForOneRegistrationHaveIndependentAttemptLifecycles(t *testing.T) {
	harness := newHarness(t, []fixtureSpec{
		{nonce: 0x66, variant: 0},
		{nonce: 0x67, variant: 1},
	})
	second := harness.createAttempt(t, 1)
	if second == harness.attemptID {
		t.Fatal("two approvals produced one attempt identity")
	}
	firstSnapshot, err := harness.component.Drive(context.Background(), harness.attemptID)
	if err != nil {
		t.Fatalf("drive first attempt: %v", err)
	}
	secondSnapshot, err := harness.component.Drive(context.Background(), second)
	if err != nil {
		t.Fatalf("drive second attempt: %v", err)
	}
	if firstSnapshot.RegistrationID != secondSnapshot.RegistrationID ||
		firstSnapshot.AttemptID == secondSnapshot.AttemptID ||
		firstSnapshot.ApprovalID == secondSnapshot.ApprovalID {
		t.Fatalf("attempt lifecycles were conflated: first=%+v second=%+v", firstSnapshot, secondSnapshot)
	}
	if !harness.backend.Snapshot(harness.attemptID).Destroyed || !harness.backend.Snapshot(second).Destroyed {
		t.Fatal("independent fake instances did not both destroy")
	}
}

func TestOneApprovalCannotDriveTwoAttemptIDs(t *testing.T) {
	harness := newHarness(t, nil)
	if _, err := harness.component.Drive(context.Background(), harness.attemptID); err != nil {
		t.Fatalf("drive committed attempt: %v", err)
	}
	aliasID := attemptIDFor(99)
	resolver := aliasingResolver{base: harness.attempts, sourceID: harness.attemptID}
	component := mustLifecycle(t, resolver, harness.attempts, harness.store, harness.backend, nil)
	_, err := component.Drive(context.Background(), aliasID)
	assertLifecycleClassification(t, err, ClassificationBinding)
	if calls := harness.backend.Snapshot(aliasID).CallCounts; len(calls) != 0 {
		t.Fatalf("approval alias reached fake backend: %v", calls)
	}
}

func TestFakeFaultMatrixEndsDestroyedOrUnresolved(t *testing.T) {
	operations := []Operation{
		OperationPrepare, OperationCreate, OperationStart,
		OperationObserve, OperationStop, OperationDestroy,
	}
	for _, operation := range operations {
		for _, moment := range []FaultMoment{FaultBeforeEffect, FaultAfterEffect} {
			t.Run(string(operation)+"/"+string(moment), func(t *testing.T) {
				harness := newHarness(t, nil)
				if err := harness.backend.InjectFault(harness.attemptID, operation, moment); err != nil {
					t.Fatalf("inject fault: %v", err)
				}
				snapshot, err := harness.component.Drive(context.Background(), harness.attemptID)
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
			})
		}
	}
}

func TestPostSideEffectInterruptionRecoversByAttemptAcrossComponentRestart(t *testing.T) {
	checkpoints := []Checkpoint{
		CheckpointAfterPrepareEffect, CheckpointAfterCreateEffect,
		CheckpointAfterStartEffect, CheckpointAfterObserveEffect,
		CheckpointAfterStopEffect, CheckpointAfterDestroyEffect,
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
		t.Run(string(checkpoint), func(t *testing.T) {
			triggered := false
			harness := newHarness(t, nil)
			harness.component.checkpoint = func(
				_ context.Context, actual Checkpoint, _ approvalattempt.AttemptID,
			) error {
				if actual == checkpoint && !triggered {
					triggered = true
					return errors.New("simulated process loss with untrusted detail")
				}
				return nil
			}
			interrupted, err := harness.component.Drive(context.Background(), harness.attemptID)
			assertLifecycleClassification(t, err, ClassificationLocalFailure)
			if interrupted.State == StateDestroyed {
				t.Fatal("checkpoint did not leave recovery work")
			}

			restarted := mustLifecycle(t, harness.attempts, harness.attempts, harness.store, harness.backend, nil)
			recoveredSet, err := restarted.RecoverCreatedAttempts(context.Background())
			if err != nil {
				t.Fatalf("recover checkpoint %s: %v", checkpoint, err)
			}
			if len(recoveredSet) != 1 {
				t.Fatalf("startup recovery returned %d attempts, want 1", len(recoveredSet))
			}
			recovered := recoveredSet[0]
			if recovered.State != StateDestroyed || recovered.CleanupRequired ||
				recovered.Failure != ClassificationLocalFailure ||
				recovered.FailureAt != checkpointOperations[checkpoint] {
				t.Fatalf("recovered disposition = %+v", recovered)
			}
		})
	}
}

func TestUnknownRecoveryRemainsUnresolvedAndRetryable(t *testing.T) {
	triggered := false
	harness := newHarness(t, nil)
	harness.component.checkpoint = func(
		_ context.Context, checkpoint Checkpoint, _ approvalattempt.AttemptID,
	) error {
		if checkpoint == CheckpointAfterCreateEffect && !triggered {
			triggered = true
			return errors.New("simulated process loss")
		}
		return nil
	}
	if _, err := harness.component.Drive(context.Background(), harness.attemptID); err == nil {
		t.Fatal("expected simulated interruption")
	}
	if err := harness.backend.InjectFault(harness.attemptID, OperationReconcile, FaultBeforeEffect); err != nil {
		t.Fatalf("inject reconcile fault: %v", err)
	}

	restarted := mustLifecycle(t, harness.attempts, harness.attempts, harness.store, harness.backend, nil)
	unresolved, err := restarted.Recover(context.Background(), harness.attemptID)
	assertLifecycleClassification(t, err, ClassificationCleanupUnresolved)
	if unresolved.State != StateUnresolved || !unresolved.CleanupRequired {
		t.Fatalf("unknown recovery disposition = %+v", unresolved)
	}
	recovered, err := restarted.Recover(context.Background(), harness.attemptID)
	if err != nil || recovered.State != StateDestroyed || recovered.CleanupRequired {
		t.Fatalf("retry disposition = %+v, %v", recovered, err)
	}
}

func TestStartupEnumerationDrivesCreatedAttemptOnceAndIgnoresExpiry(t *testing.T) {
	harness := newHarness(t, nil)
	counter := &countingEnumerator{base: harness.attempts}
	component := mustLifecycle(t, harness.attempts, counter, harness.store, harness.backend, nil)
	harness.clock.set(1_785_456_301)

	results, err := component.RecoverCreatedAttempts(context.Background())
	if err != nil || len(results) != 1 || results[0].State != StateDestroyed {
		t.Fatalf("startup recovery = %+v, %v", results, err)
	}
	if counter.calls != 1 {
		t.Fatalf("startup enumeration calls = %d, want 1", counter.calls)
	}
	before := harness.backend.Snapshot(harness.attemptID).CallCounts
	results, err = component.RecoverCreatedAttempts(context.Background())
	if err != nil || len(results) != 1 || results[0].State != StateDestroyed {
		t.Fatalf("startup replay = %+v, %v", results, err)
	}
	after := harness.backend.Snapshot(harness.attemptID).CallCounts
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("startup replay redrove effects: before=%v after=%v", before, after)
	}
}

func TestDefensiveCopiesSnapshotsAndFixedErrors(t *testing.T) {
	harness := newHarness(t, nil)
	resolved, err := harness.attempts.ResolveCreated(context.Background(), harness.attemptID)
	if err != nil {
		t.Fatalf("resolve created attempt: %v", err)
	}
	resolved.Registration.ExactPlanBytes[0] ^= 0xff
	resolved.PlanRoleBindings.ProfileReviewAttestationDigests[0][0] ^= 0xff
	again, err := harness.attempts.ResolveCreated(context.Background(), harness.attemptID)
	if err != nil || !bytes.Equal(again.Registration.ExactPlanBytes, harness.plan) ||
		again.PlanRoleBindings.ProfileReviewAttestationDigests[0][0] == resolved.PlanRoleBindings.ProfileReviewAttestationDigests[0][0] {
		t.Fatal("created-attempt resolution did not defensively copy retained bindings")
	}
	if err := harness.backend.InjectFault(harness.attemptID, OperationObserve, FaultAfterEffect); err != nil {
		t.Fatalf("inject fault: %v", err)
	}
	_, err = harness.component.Drive(context.Background(), harness.attemptID)
	if err == nil || bytes.Contains([]byte(err.Error()), harness.plan) ||
		bytes.Contains([]byte(err.Error()), []byte("simulated process loss with untrusted detail")) {
		t.Fatalf("fixed error leaked caller/backend content: %q", err)
	}
	first := harness.backend.Snapshot(harness.attemptID)
	first.CallCounts[OperationPrepare] = 999
	second := harness.backend.Snapshot(harness.attemptID)
	if second.CallCounts[OperationPrepare] != 1 {
		t.Fatal("backend snapshot did not defensively copy call counts")
	}
}

func TestLifecycleStoreRefusesBeyondAttemptKeyedCapacity(t *testing.T) {
	harness := newHarness(t, nil)
	created, err := harness.attempts.ResolveCreated(context.Background(), harness.attemptID)
	if err != nil {
		t.Fatalf("resolve created: %v", err)
	}
	store := NewMemoryStore()
	for index := 1; index <= MaxLifecycleRecords; index++ {
		record := cloneCreatedAttempt(created)
		record.Attempt.AttemptID = attemptIDFor(uint64(index))
		record.Attempt.ApprovalID = approvalIDFor(uint64(index))
		record.Approval.ApprovalID = record.Attempt.ApprovalID
		record.Approval.ConsumedAttemptID = record.Attempt.AttemptID
		if _, began, err := store.begin(context.Background(), record); err != nil || !began {
			t.Fatalf("begin record %d: began=%t err=%v", index, began, err)
		}
	}
	overflow := cloneCreatedAttempt(created)
	overflow.Attempt.AttemptID = attemptIDFor(MaxLifecycleRecords + 1)
	overflow.Attempt.ApprovalID = approvalIDFor(MaxLifecycleRecords + 1)
	overflow.Approval.ApprovalID = overflow.Attempt.ApprovalID
	overflow.Approval.ConsumedAttemptID = overflow.Attempt.AttemptID
	if _, _, err := store.begin(context.Background(), overflow); err == nil {
		t.Fatal("fixed lifecycle capacity accepted another record")
	} else {
		assertLifecycleClassification(t, err, ClassificationCapacity)
	}
}

type testHarness struct {
	component      *Component
	attempts       *registrationstate.ApprovalAttemptComponent
	authorityStore *registrationstate.FixedFileStore
	clock          *testClock
	store          *MemoryStore
	backend        *FakeBackend
	registrationID v0candidate.RegistrationID
	attemptID      approvalattempt.AttemptID
	plan           []byte
	vectors        []approvalattempt.FixtureVector
}

type fixtureSpec struct {
	nonce   byte
	variant byte
}

func newHarness(t *testing.T, specs []fixtureSpec) *testHarness {
	t.Helper()
	if len(specs) == 0 {
		specs = []fixtureSpec{{nonce: 0x66}}
	}
	plan, err := os.ReadFile(filepath.Clean(ordinaryPlanPath))
	if err != nil {
		t.Fatalf("read ordinary plan: %v", err)
	}
	stateStore, err := registrationstate.NewFixedFileStore(
		filepath.Join(t.TempDir(), "supervisor-state.json"),
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
		t.Fatalf("new fixed store: %v", err)
	}
	clock := &testClock{value: 1_785_456_000}
	registrations, err := registrationstate.New(registrationstate.Options{
		Store: stateStore, Clock: clock,
		Identifiers: fixedRegistrationIDSource{id: repeated16[v0candidate.RegistrationID](0x33)},
	})
	if err != nil {
		t.Fatalf("new registration component: %v", err)
	}
	issued, err := registrations.RegisterPlan(
		context.Background(), registrationstate.AuthenticatedCallContext{
			Authenticated: true, Role: registrationstate.CallerDaemon,
			Purpose: registrationstate.RegisterPlanPurpose,
		}, plan, ordinaryBindings(),
	)
	if err != nil {
		t.Fatalf("register ordinary plan: %v", err)
	}
	registrationID := issued.View().RegistrationID
	vectors := make([]approvalattempt.FixtureVector, len(specs))
	for index, spec := range specs {
		vectors[index] = lifecycleFixtureVector(t, issued.View(), spec.nonce, spec.variant)
	}
	verifier, err := approvalattempt.NewFixtureVerifier(vectors)
	if err != nil {
		t.Fatalf("new fixture verifier: %v", err)
	}
	attempts, err := registrationstate.NewApprovalAttempt(registrationstate.ApprovalAttemptOptions{
		Store: stateStore, Clock: clock, Verifier: verifier,
		ApprovalIdentifiers: &approvalIDSequence{next: 1},
		AttemptIdentifiers:  &attemptIDSequence{next: 1},
		Integrity:           fixedIntegrity{assessedAt: 1_785_456_000},
	})
	if err != nil {
		t.Fatalf("new approval-attempt component: %v", err)
	}
	lifecycleStore := NewMemoryStore()
	backend := NewFakeBackend()
	component := mustLifecycle(t, attempts, attempts, lifecycleStore, backend, nil)
	harness := &testHarness{
		component: component, attempts: attempts, authorityStore: stateStore,
		clock: clock, store: lifecycleStore,
		backend: backend, registrationID: registrationID, plan: plan, vectors: vectors,
	}
	harness.attemptID = harness.createAttempt(t, 0)
	return harness
}

func (harness *testHarness) createAttempt(t *testing.T, vectorIndex int) approvalattempt.AttemptID {
	t.Helper()
	vector := harness.vectors[vectorIndex]
	submission, err := harness.attempts.SubmitApproval(
		context.Background(), registrationstate.AuthenticatedCallContext{
			Authenticated: true, Role: registrationstate.CallerBroker,
			Purpose: registrationstate.SubmitApprovalPurpose,
		}, harness.registrationID, vector.EnvelopeBytes,
	)
	if err != nil {
		t.Fatalf("submit approval %d: %v", vectorIndex, err)
	}
	created, err := harness.attempts.RequestAttempt(
		context.Background(), registrationstate.AuthenticatedCallContext{
			Authenticated: true, Role: registrationstate.CallerDaemon,
			Purpose: registrationstate.RequestAttemptPurpose,
		}, harness.registrationID, submission.Reference,
	)
	if err != nil {
		t.Fatalf("request attempt %d: %v", vectorIndex, err)
	}
	return created.Reference.AttemptID()
}

func mustLifecycle(
	t *testing.T,
	resolver AttemptResolver,
	enumerator CreatedAttemptEnumerator,
	store *MemoryStore,
	backend *FakeBackend,
	checkpoint CheckpointHook,
) *Component {
	t.Helper()
	component, err := New(Options{
		Attempts: resolver, CreatedAttempts: enumerator,
		Store: store, Backend: backend, Checkpoint: checkpoint,
	})
	if err != nil {
		t.Fatalf("new registered lifecycle: %v", err)
	}
	return component
}

type alteringResolver struct {
	base  AttemptResolver
	alter func(*registrationstate.CreatedAttempt)
}

type aliasingResolver struct {
	base     AttemptResolver
	sourceID approvalattempt.AttemptID
}

func (resolver aliasingResolver) ResolveCreated(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (registrationstate.CreatedAttempt, error) {
	created, err := resolver.base.ResolveCreated(ctx, resolver.sourceID)
	if err != nil {
		return registrationstate.CreatedAttempt{}, err
	}
	created.Attempt.AttemptID = attemptID
	created.Approval.ConsumedAttemptID = attemptID
	return created, nil
}

func (resolver *alteringResolver) ResolveCreated(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (registrationstate.CreatedAttempt, error) {
	created, err := resolver.base.ResolveCreated(ctx, attemptID)
	if err == nil && resolver.alter != nil {
		resolver.alter(&created)
	}
	return created, err
}

type countingEnumerator struct {
	base  CreatedAttemptEnumerator
	calls int
}

func (enumerator *countingEnumerator) CreatedAttempts(
	ctx context.Context,
) ([]approvalattempt.AttemptReference, error) {
	enumerator.calls++
	return enumerator.base.CreatedAttempts(ctx)
}

type testClock struct {
	mu    sync.Mutex
	value uint64
}

func (clock *testClock) ObserveUnixSeconds(context.Context) (uint64, error) {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	return clock.value, nil
}

func (clock *testClock) set(value uint64) {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	clock.value = value
}

type fixedIntegrity struct{ assessedAt v0candidate.UInt53 }

func (integrity fixedIntegrity) Assess(
	_ context.Context,
	preflight registrationstate.IntegrityPreflight,
) (registrationstate.RuntimeIntegrityAssessment, error) {
	return registrationstate.RuntimeIntegrityAssessment{
		Preflight: preflight, AssessedAt: integrity.assessedAt, Permitted: true,
	}, nil
}

type fixedRegistrationIDSource struct{ id v0candidate.RegistrationID }

func (source fixedRegistrationIDSource) NewRegistrationID(
	context.Context,
) (v0candidate.RegistrationID, error) {
	return source.id, nil
}

type approvalIDSequence struct {
	mu   sync.Mutex
	next uint64
}

func (source *approvalIDSequence) NewApprovalID(
	context.Context,
) (approvalattempt.ApprovalID, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	var result approvalattempt.ApprovalID
	result[0] = 0xa1
	binary.BigEndian.PutUint64(result[8:], source.next)
	source.next++
	return result, nil
}

type attemptIDSequence struct {
	mu   sync.Mutex
	next uint64
}

func (source *attemptIDSequence) NewAttemptID(
	context.Context,
) (approvalattempt.AttemptID, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	result := attemptIDFor(source.next)
	source.next++
	return result, nil
}

func attemptIDFor(value uint64) approvalattempt.AttemptID {
	var result approvalattempt.AttemptID
	result[0] = 0xb2
	binary.BigEndian.PutUint64(result[8:], value)
	return result
}

func approvalIDFor(value uint64) approvalattempt.ApprovalID {
	var result approvalattempt.ApprovalID
	result[0] = 0xa1
	binary.BigEndian.PutUint64(result[8:], value)
	return result
}

func lifecycleFixtureVector(
	t *testing.T,
	registration v0candidate.PlanRegistration,
	nonceByte byte,
	variant byte,
) approvalattempt.FixtureVector {
	t.Helper()
	envelope := mustRead(t, "../../../schemas/conformance/v0/approval-grant/ordinary.cose")
	payload := mustRead(t, "../../../schemas/conformance/v0/approval-grant/ordinary.payload.cbor")
	protected := mustRead(t, "../../../schemas/conformance/v0/approval-grant/ordinary.protected.cbor")
	if variant != 0 {
		envelope[len(envelope)-1] ^= variant
		payload[len(payload)-1] ^= variant
	}
	return approvalattempt.FixtureVector{
		EnvelopeBytes: envelope, PayloadBytes: payload, ProtectedHeaderBytes: protected,
		ProtectedKeyID: []byte("approval-test-key"),
		View: approvalattempt.ApprovalGrant{
			ObjectType: approvalattempt.ApprovalGrantObjectType, ObjectVersion: 0,
			InstallationID: registration.InstallationID, EpochDigest: registration.EpochDigest,
			RegistrationID: registration.RegistrationID, PlanDigest: registration.PlanDigest,
			SupervisorID: registration.SupervisorID,
			AttemptNonce: repeated16[approvalattempt.AttemptNonce](nonceByte),
			Purpose:      approvalattempt.ApprovalGrantPurpose, Audience: approvalattempt.ApprovalGrantAudience,
			IssuedAt: 1_785_456_000, ExpiresAt: 1_785_456_300,
		},
		ResolvedEpochSequence: registration.EpochSequence,
		AuthorizationIdentity: repeated32[approvalattempt.ApprovalKeyAuthorizationIdentity](0x99),
		SignatureAccepted:     true,
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	value, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return value
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

func assertLifecycleClassification(t *testing.T, err error, want Classification) {
	t.Helper()
	classification, ok := ErrorClassification(err)
	if !ok || classification != want {
		t.Fatalf("classification = %q (%t), want %q: %v", classification, ok, want, err)
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
