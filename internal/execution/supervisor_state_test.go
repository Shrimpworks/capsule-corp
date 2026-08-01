package execution

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

var errInjectedCheckpoint = errors.New("injected checkpoint stop")

func TestRegistrationRetainsExactImmutableBytes(t *testing.T) {
	fixture := newSupervisorFixture(t, nil)
	plan := []byte{0xa1, 0x01, 0x02}
	registration, err := fixture.core.RegisterPlan(context.Background(), plan)
	if err != nil {
		t.Fatalf("register plan: %v", err)
	}
	plan[0] = 0xff

	got, err := fixture.core.GetRegisteredPlan(context.Background(), registration.ID)
	if err != nil {
		t.Fatalf("get registered plan: %v", err)
	}
	if got[0] != 0xa1 || registration.ExactPlanBytes[0] != 0xa1 {
		t.Fatalf("registered plan changed with caller buffer: %x / %x", got, registration.ExactPlanBytes)
	}
	got[0] = 0xee
	gotAgain, err := fixture.core.GetRegisteredPlan(context.Background(), registration.ID)
	if err != nil {
		t.Fatalf("get registered plan again: %v", err)
	}
	if gotAgain[0] != 0xa1 {
		t.Fatalf("stored plan changed with returned buffer: %x", gotAgain)
	}
	if registration.PlanDigest != DigestBytes([]byte{0xa1, 0x01, 0x02}) {
		t.Fatal("registration did not bind the exact original bytes")
	}
}

func TestApprovalBindingsFailClosed(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*VerifiedApproval)
		want   error
	}{
		{"wrong installation", func(value *VerifiedApproval) { value.InstallationID = repeatedID(71) }, ErrApprovalBinding},
		{"wrong epoch", func(value *VerifiedApproval) { value.EpochDigest = repeatedDigest(72) }, ErrApprovalBinding},
		{"wrong registration", func(value *VerifiedApproval) { value.RegistrationID = repeatedID(73) }, ErrApprovalBinding},
		{"wrong plan", func(value *VerifiedApproval) { value.PlanDigest = repeatedDigest(74) }, ErrApprovalBinding},
		{"wrong supervisor", func(value *VerifiedApproval) { value.SupervisorID = repeatedID(75) }, ErrApprovalBinding},
		{"zero nonce", func(value *VerifiedApproval) { value.AttemptNonce = OpaqueID{} }, ErrApprovalBinding},
		{"wrong purpose", func(value *VerifiedApproval) { value.Purpose = "capsule.execution.attest" }, ErrApprovalBinding},
		{"wrong audience", func(value *VerifiedApproval) { value.Audience = "capsule.daemon" }, ErrApprovalBinding},
		{"future issue", func(value *VerifiedApproval) { value.IssuedAt = testNow.Add(time.Second) }, ErrApprovalExpired},
		{"expired", func(value *VerifiedApproval) { value.ExpiresAt = testNow }, ErrApprovalExpired},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newSupervisorFixture(t, nil)
			registration := fixture.register(t, "binding-plan")
			approval := fixture.validApproval(registration, "binding-payload")
			test.mutate(&approval)
			fixture.approval = approval
			_, err := fixture.core.SubmitApproval(
				context.Background(),
				registration.ID,
				[]byte("signed-envelope"),
			)
			if !errors.Is(err, test.want) {
				t.Fatalf("got %v, want %v", err, test.want)
			}
		})
	}
}

func TestEquivalentSignatureCannotCreateSecondApproval(t *testing.T) {
	fixture := newSupervisorFixture(t, nil)
	registration := fixture.register(t, "duplicate-plan")
	fixture.approval = fixture.validApproval(registration, "same-canonical-payload")

	if _, err := fixture.core.SubmitApproval(
		context.Background(), registration.ID, []byte("low-s-envelope"),
	); err != nil {
		t.Fatalf("submit first approval: %v", err)
	}
	if _, err := fixture.core.SubmitApproval(
		context.Background(), registration.ID, []byte("high-s-envelope"),
	); !errors.Is(err, ErrDuplicateApproval) {
		t.Fatalf("equivalent approval got %v, want duplicate denial", err)
	}
}

func TestApprovalRetainsExactPayloadAndEnvelopeBytes(t *testing.T) {
	fixture := newSupervisorFixture(t, nil)
	registration := fixture.register(t, "approval-bytes-plan")
	payload := []byte("canonical-payload")
	envelope := []byte("signed-envelope")
	fixture.approval = fixture.validApproval(registration, string(payload))
	fixture.approval.CanonicalPayload = payload

	approval, err := fixture.core.SubmitApproval(context.Background(), registration.ID, envelope)
	if err != nil {
		t.Fatalf("submit approval: %v", err)
	}
	payload[0] = 'X'
	envelope[0] = 'X'
	approval.CanonicalPayload[0] = 'Y'
	approval.SignedEnvelope[0] = 'Y'

	stored := fixture.snapshot(t).Approvals[approval.ID]
	if string(stored.CanonicalPayload) != "canonical-payload" ||
		string(stored.SignedEnvelope) != "signed-envelope" {
		t.Fatalf("stored approval bytes were mutated: %+v", stored)
	}
}

func TestApprovalConsumptionCreatesAtMostOneAttemptConcurrently(t *testing.T) {
	fixture := newSupervisorFixture(t, nil)
	registration, approval := fixture.registerAndApprove(t, "concurrent-plan")

	start := make(chan struct{})
	results := make(chan error, 2)
	var workers sync.WaitGroup
	for range 2 {
		workers.Add(1)
		go func() {
			defer workers.Done()
			<-start
			_, err := fixture.core.RequestAttempt(
				context.Background(), registration.ID, approval.ID,
			)
			results <- err
		}()
	}
	close(start)
	workers.Wait()
	close(results)

	var succeeded, consumed int
	for err := range results {
		switch {
		case err == nil:
			succeeded++
		case errors.Is(err, ErrApprovalConsumed):
			consumed++
		default:
			t.Fatalf("unexpected concurrent result: %v", err)
		}
	}
	if succeeded != 1 || consumed != 1 {
		t.Fatalf("got succeeded=%d consumed=%d, want one of each", succeeded, consumed)
	}
	state := fixture.snapshot(t)
	if len(state.Attempts) != 1 {
		t.Fatalf("one approval created %d attempts", len(state.Attempts))
	}
}

func TestApprovalForPlanACannotCreateAttemptForPlanB(t *testing.T) {
	fixture := newSupervisorFixture(t, nil)
	registrationA, approvalA := fixture.registerAndApprove(t, "plan-a")
	registrationB := fixture.register(t, "plan-b")

	if _, err := fixture.core.RequestAttempt(
		context.Background(), registrationB.ID, approvalA.ID,
	); !errors.Is(err, ErrApprovalBinding) {
		t.Fatalf("cross-plan attempt got %v, want binding denial", err)
	}
	state := fixture.snapshot(t)
	if state.Approvals[approvalA.ID].State != ApprovalUnused || len(state.Attempts) != 0 {
		t.Fatalf("cross-plan denial consumed authority: %+v", state)
	}
	if registrationA.PlanDigest == registrationB.PlanDigest {
		t.Fatal("test plans unexpectedly share a digest")
	}
}

func TestIntegrityPreflightFailureDoesNotConsumeApproval(t *testing.T) {
	fixture := newSupervisorFixture(t, nil)
	registration, approval := fixture.registerAndApprove(t, "integrity-denied")
	fixture.integrityErr = errors.New("debugged peer")

	if _, err := fixture.core.RequestAttempt(
		context.Background(), registration.ID, approval.ID,
	); err == nil {
		t.Fatal("integrity failure unexpectedly created an attempt")
	}
	state := fixture.snapshot(t)
	if state.Approvals[approval.ID].State != ApprovalUnused || len(state.Attempts) != 0 {
		t.Fatalf("preflight failure consumed authority: %+v", state)
	}
}

func TestTrustTransitionFencesApprovalsUntilEveryComponentAccepts(t *testing.T) {
	fixture := newSupervisorFixture(t, nil)
	registration := fixture.register(t, "old-epoch-plan")
	fixture.approval = fixture.validApproval(registration, "old-epoch-unused")
	approval, err := fixture.core.SubmitApproval(
		context.Background(), registration.ID, []byte("unused-envelope"),
	)
	if err != nil {
		t.Fatalf("submit approval: %v", err)
	}

	transitionID := repeatedID(80)
	targetEpoch := repeatedDigest(81)
	if err := fixture.core.FenceTrustTransition(
		context.Background(), transitionID, targetEpoch,
	); err != nil {
		t.Fatalf("fence transition: %v", err)
	}
	state := fixture.snapshot(t)
	if state.AttemptsEnabled || state.TrustPhase != TrustTransitionFenced {
		t.Fatalf("transition did not fence attempts: %+v", state)
	}
	if state.Approvals[approval.ID].State != ApprovalInvalidated {
		t.Fatal("transition did not invalidate the unused old-epoch approval")
	}
	if _, err := fixture.core.RegisterPlan(context.Background(), []byte("blocked")); !errors.Is(err, ErrTransitionFenced) {
		t.Fatalf("registration while fenced got %v", err)
	}

	if err := fixture.core.CommitEpochForComponentAcceptance(
		context.Background(), transitionID,
	); err != nil {
		t.Fatalf("commit epoch: %v", err)
	}
	for index, role := range requiredComponentRoles[:len(requiredComponentRoles)-1] {
		if err := fixture.core.AcceptCurrentComponent(
			context.Background(), transitionID, role, repeatedID(byte(90+index)),
		); err != nil {
			t.Fatalf("accept %s: %v", role, err)
		}
	}
	state = fixture.snapshot(t)
	if state.AttemptsEnabled || state.TrustPhase != TrustAwaitingComponentAcceptance {
		t.Fatal("partial component acceptance enabled execution")
	}
	lastRole := requiredComponentRoles[len(requiredComponentRoles)-1]
	if err := fixture.core.AcceptCurrentComponent(
		context.Background(), transitionID, lastRole, repeatedID(99),
	); err != nil {
		t.Fatalf("accept final component: %v", err)
	}
	state = fixture.snapshot(t)
	if !state.AttemptsEnabled || state.TrustPhase != TrustStable || state.EpochDigest != targetEpoch {
		t.Fatalf("complete acceptance did not enable target epoch: %+v", state)
	}

	fixture.approval = fixture.validApproval(registration, "old-registration-new-epoch")
	fixture.approval.EpochDigest = targetEpoch
	if _, err := fixture.core.SubmitApproval(
		context.Background(), registration.ID, []byte("wrong-epoch-registration"),
	); !errors.Is(err, ErrApprovalBinding) {
		t.Fatalf("old registration crossed epoch: %v", err)
	}
}

func TestTrustTransitionRefusesAnOpenAttempt(t *testing.T) {
	fixture := newSupervisorFixture(t, nil)
	_, _, _ = fixture.createAttempt(t, "open-attempt")
	if err := fixture.core.FenceTrustTransition(
		context.Background(), repeatedID(101), repeatedDigest(102),
	); !errors.Is(err, ErrCleanupRequired) {
		t.Fatalf("transition with open attempt got %v", err)
	}
}

func TestIncompleteTransitionRecoversToRepairRequired(t *testing.T) {
	fixture := newSupervisorFixture(t, nil)
	transitionID := repeatedID(103)
	if err := fixture.core.FenceTrustTransition(
		context.Background(), transitionID, repeatedDigest(104),
	); err != nil {
		t.Fatalf("fence transition: %v", err)
	}
	if err := fixture.core.RecoverTrustState(context.Background()); err != nil {
		t.Fatalf("recover trust state: %v", err)
	}
	state := fixture.snapshot(t)
	if state.TrustPhase != TrustRepairRequired || state.AttemptsEnabled {
		t.Fatalf("incomplete transition did not fail closed: %+v", state)
	}
	if err := fixture.core.AcceptCurrentComponent(
		context.Background(), transitionID, ComponentDaemon, repeatedID(105),
	); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("component acceptance cleared repair state: %v", err)
	}
}

func TestDevelopmentLifecycleCompletesAndCleansUp(t *testing.T) {
	fixture := newSupervisorFixture(t, nil)
	_, _, attempt := fixture.createAttempt(t, "lifecycle-success")
	backend := newFakeDevelopmentLifecycle()

	if err := fixture.core.DriveDevelopmentAttempt(
		context.Background(), attempt.ID, backend,
	); err != nil {
		t.Fatalf("drive development attempt: %v", err)
	}
	got := fixture.snapshot(t).Attempts[attempt.ID]
	if got.State != AttemptTerminal || got.CleanupRequired ||
		got.TerminalClass != TerminalSimulationPassed {
		t.Fatalf("unexpected terminal attempt: %+v", got)
	}
	wantCalls := []string{"prepare", "create", "stage", "start", "wait", "collect", "destroy"}
	if fmt.Sprint(backend.calls) != fmt.Sprint(wantCalls) {
		t.Fatalf("backend calls = %v, want %v", backend.calls, wantCalls)
	}
}

func TestDevelopmentLifecycleRejectsGuestCreatingBackend(t *testing.T) {
	fixture := newSupervisorFixture(t, nil)
	_, _, attempt := fixture.createAttempt(t, "no-real-backend")
	backend := newFakeDevelopmentLifecycle()
	backend.createsGuest = true

	if err := fixture.core.DriveDevelopmentAttempt(
		context.Background(), attempt.ID, backend,
	); !errors.Is(err, ErrRealBackendBlocked) {
		t.Fatalf("guest-creating backend got %v", err)
	}
	if fixture.snapshot(t).Attempts[attempt.ID].State != AttemptCreated {
		t.Fatal("blocked backend changed attempt state")
	}
}

func TestPostCreateFailureStillDestroysBackend(t *testing.T) {
	fixture := newSupervisorFixture(t, nil)
	_, _, attempt := fixture.createAttempt(t, "stage-failure")
	backend := newFakeDevelopmentLifecycle()
	backend.fail["stage"] = errors.New("stage failed")

	if err := fixture.core.DriveDevelopmentAttempt(
		context.Background(), attempt.ID, backend,
	); err == nil {
		t.Fatal("stage failure unexpectedly succeeded")
	}
	got := fixture.snapshot(t).Attempts[attempt.ID]
	if got.State != AttemptTerminal || got.CleanupRequired ||
		got.TerminalClass != TerminalSimulationFailed {
		t.Fatalf("post-create failure did not clean up: %+v", got)
	}
	if backend.destroyed != 1 {
		t.Fatalf("destroy calls = %d, want 1", backend.destroyed)
	}
}

func TestAmbiguousCreateRetainsCleanupUntilReconciled(t *testing.T) {
	fixture := newSupervisorFixture(t, nil)
	_, _, attempt := fixture.createAttempt(t, "ambiguous-create")
	backend := newFakeDevelopmentLifecycle()
	backend.createSideEffectOnError = true
	backend.fail["create"] = errors.New("connection lost after create")
	backend.forceUnknown = true

	err := fixture.core.DriveDevelopmentAttempt(context.Background(), attempt.ID, backend)
	if !errors.Is(err, ErrCleanupUnresolved) {
		t.Fatalf("ambiguous create got %v, want unresolved cleanup", err)
	}
	got := fixture.snapshot(t).Attempts[attempt.ID]
	if got.State != AttemptUnresolved || !got.CleanupRequired ||
		got.TerminalClass != TerminalCleanupUnresolved {
		t.Fatalf("ambiguous create lost cleanup obligation: %+v", got)
	}

	backend.forceUnknown = false
	if err := fixture.core.ReconcileDevelopmentAttempt(
		context.Background(), attempt.ID, backend,
	); err != nil {
		t.Fatalf("reconcile known backend: %v", err)
	}
	got = fixture.snapshot(t).Attempts[attempt.ID]
	if got.State != AttemptTerminal || got.CleanupRequired {
		t.Fatalf("reconciliation did not discharge cleanup: %+v", got)
	}
}

func TestCheckpointStopAfterDurableHandleCanBeReconciled(t *testing.T) {
	fixture := newSupervisorFixture(t, func(checkpoint Checkpoint) error {
		if checkpoint == CheckpointBackendHandle {
			return errInjectedCheckpoint
		}
		return nil
	})
	_, _, attempt := fixture.createAttempt(t, "checkpoint-recovery")
	backend := newFakeDevelopmentLifecycle()

	if err := fixture.core.DriveDevelopmentAttempt(
		context.Background(), attempt.ID, backend,
	); !errors.Is(err, errInjectedCheckpoint) {
		t.Fatalf("checkpoint stop got %v", err)
	}
	stopped := fixture.snapshot(t).Attempts[attempt.ID]
	if stopped.State != AttemptBackendCreated || !stopped.CleanupRequired ||
		stopped.BackendHandle == "" {
		t.Fatalf("durable checkpoint state is unsafe: %+v", stopped)
	}

	restarted := fixture.newCore(t, nil)
	if err := restarted.ReconcileDevelopmentAttempt(
		context.Background(), attempt.ID, backend,
	); err != nil {
		t.Fatalf("restart reconcile: %v", err)
	}
	recovered := fixture.snapshot(t).Attempts[attempt.ID]
	if recovered.State != AttemptTerminal || recovered.CleanupRequired || backend.destroyed != 1 {
		t.Fatalf("restart did not clean durable handle: %+v destroyed=%d", recovered, backend.destroyed)
	}
}

func TestAuthoritativeAbsenceClearsPreCreateCleanupIntent(t *testing.T) {
	fixture := newSupervisorFixture(t, func(checkpoint Checkpoint) error {
		if checkpoint == CheckpointBackendCreateIntent {
			return errInjectedCheckpoint
		}
		return nil
	})
	_, _, attempt := fixture.createAttempt(t, "pre-create-stop")
	backend := newFakeDevelopmentLifecycle()

	if err := fixture.core.DriveDevelopmentAttempt(
		context.Background(), attempt.ID, backend,
	); !errors.Is(err, errInjectedCheckpoint) {
		t.Fatalf("checkpoint stop got %v", err)
	}
	stopped := fixture.snapshot(t).Attempts[attempt.ID]
	if stopped.State != AttemptBackendCreateIntent || !stopped.CleanupRequired {
		t.Fatalf("create intent was not durable: %+v", stopped)
	}

	restarted := fixture.newCore(t, nil)
	if err := restarted.ReconcileDevelopmentAttempt(
		context.Background(), attempt.ID, backend,
	); err != nil {
		t.Fatalf("reconcile authoritative absence: %v", err)
	}
	recovered := fixture.snapshot(t).Attempts[attempt.ID]
	if recovered.State != AttemptTerminal || recovered.CleanupRequired ||
		recovered.TerminalClass != TerminalBackendFailed {
		t.Fatalf("authoritative absence did not close intent: %+v", recovered)
	}
}

const testPlanLimit = 1024

var testNow = time.Unix(1_785_456_100, 0).UTC()

type supervisorFixture struct {
	core         *SupervisorCore
	store        *MemoryStateStore
	clock        fixedClock
	identifiers  *sequenceIdentifiers
	approval     VerifiedApproval
	integrityErr error
}

func newSupervisorFixture(t *testing.T, checkpoint CheckpointHook) *supervisorFixture {
	t.Helper()
	store, err := NewMemoryStateStore(repeatedID(1), repeatedID(2), repeatedDigest(3))
	if err != nil {
		t.Fatalf("new memory store: %v", err)
	}
	fixture := &supervisorFixture{
		store:       store,
		clock:       fixedClock{now: testNow},
		identifiers: &sequenceIdentifiers{next: 10},
	}
	fixture.core = fixture.newCore(t, checkpoint)
	return fixture
}

func (fixture *supervisorFixture) newCore(
	t *testing.T,
	checkpoint CheckpointHook,
) *SupervisorCore {
	t.Helper()
	core, err := NewSupervisorCore(SupervisorOptions{
		Store: fixture.store,
		PlanValidator: PlanValidatorFunc(func(_ context.Context, plan []byte) error {
			if len(plan) == 0 {
				return errors.New("empty plan")
			}
			return nil
		}),
		ApprovalVerifier: ApprovalVerifierFunc(func(
			_ context.Context,
			_ []byte,
		) (VerifiedApproval, error) {
			return cloneVerifiedApproval(fixture.approval), nil
		}),
		Integrity: IntegrityAssessorFunc(func(
			_ context.Context,
			_ PlanRegistration,
		) error {
			return fixture.integrityErr
		}),
		Identifiers:      fixture.identifiers,
		Clock:            fixture.clock,
		Checkpoint:       checkpoint,
		MaxPlanBytes:     testPlanLimit,
		MaxApprovalBytes: 1024,
	})
	if err != nil {
		t.Fatalf("new supervisor core: %v", err)
	}
	return core
}

func (fixture *supervisorFixture) register(t *testing.T, plan string) PlanRegistration {
	t.Helper()
	registration, err := fixture.core.RegisterPlan(context.Background(), []byte(plan))
	if err != nil {
		t.Fatalf("register plan: %v", err)
	}
	return registration
}

func (fixture *supervisorFixture) registerAndApprove(
	t *testing.T,
	plan string,
) (PlanRegistration, ApprovalRecord) {
	t.Helper()
	registration := fixture.register(t, plan)
	fixture.approval = fixture.validApproval(registration, "approval:"+plan)
	approval, err := fixture.core.SubmitApproval(
		context.Background(), registration.ID, []byte("signed:"+plan),
	)
	if err != nil {
		t.Fatalf("submit approval: %v", err)
	}
	return registration, approval
}

func (fixture *supervisorFixture) createAttempt(
	t *testing.T,
	plan string,
) (PlanRegistration, ApprovalRecord, AttemptRecord) {
	t.Helper()
	registration, approval := fixture.registerAndApprove(t, plan)
	attempt, err := fixture.core.RequestAttempt(
		context.Background(), registration.ID, approval.ID,
	)
	if err != nil {
		t.Fatalf("request attempt: %v", err)
	}
	return registration, approval, attempt
}

func (fixture *supervisorFixture) validApproval(
	registration PlanRegistration,
	payload string,
) VerifiedApproval {
	return VerifiedApproval{
		CanonicalPayload: []byte(payload),
		InstallationID:   registration.InstallationID,
		EpochDigest:      registration.EpochDigest,
		RegistrationID:   registration.ID,
		PlanDigest:       registration.PlanDigest,
		SupervisorID:     registration.SupervisorID,
		AttemptNonce:     repeatedID(60),
		Purpose:          ApprovalPurpose,
		Audience:         ApprovalAudience,
		IssuedAt:         testNow.Add(-time.Minute),
		ExpiresAt:        testNow.Add(time.Minute),
	}
}

func (fixture *supervisorFixture) snapshot(t *testing.T) InstallationState {
	t.Helper()
	state, err := fixture.store.Snapshot(context.Background())
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	return state
}

type fixedClock struct {
	now time.Time
}

func (clock fixedClock) Now() time.Time {
	return clock.now
}

type sequenceIdentifiers struct {
	mu   sync.Mutex
	next byte
}

func (source *sequenceIdentifiers) NewIdentifier() (OpaqueID, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	id := repeatedID(source.next)
	source.next++
	return id, nil
}

type fakeDevelopmentLifecycle struct {
	calls                   []string
	fail                    map[string]error
	present                 map[OpaqueID]string
	createsGuest            bool
	createSideEffectOnError bool
	forceUnknown            bool
	destroyed               int
}

func newFakeDevelopmentLifecycle() *fakeDevelopmentLifecycle {
	return &fakeDevelopmentLifecycle{
		fail:    make(map[string]error),
		present: make(map[OpaqueID]string),
	}
}

func (backend *fakeDevelopmentLifecycle) Name() string { return "no-guest-fake" }

func (backend *fakeDevelopmentLifecycle) CreatesGuest() bool { return backend.createsGuest }

func (backend *fakeDevelopmentLifecycle) Prepare(_ context.Context, _ AttemptRecord) error {
	backend.calls = append(backend.calls, "prepare")
	return backend.fail["prepare"]
}

func (backend *fakeDevelopmentLifecycle) Create(_ context.Context, key OpaqueID) (string, error) {
	backend.calls = append(backend.calls, "create")
	handle := fmt.Sprintf("fake-%x", key[:4])
	if backend.fail["create"] == nil || backend.createSideEffectOnError {
		backend.present[key] = handle
	}
	return handle, backend.fail["create"]
}

func (backend *fakeDevelopmentLifecycle) Stage(_ context.Context, _ string) error {
	backend.calls = append(backend.calls, "stage")
	return backend.fail["stage"]
}

func (backend *fakeDevelopmentLifecycle) Start(_ context.Context, _ string) error {
	backend.calls = append(backend.calls, "start")
	return backend.fail["start"]
}

func (backend *fakeDevelopmentLifecycle) Wait(_ context.Context, _ string) error {
	backend.calls = append(backend.calls, "wait")
	return backend.fail["wait"]
}

func (backend *fakeDevelopmentLifecycle) Collect(_ context.Context, _ string) error {
	backend.calls = append(backend.calls, "collect")
	return backend.fail["collect"]
}

func (backend *fakeDevelopmentLifecycle) Destroy(_ context.Context, handle string) error {
	backend.calls = append(backend.calls, "destroy")
	if err := backend.fail["destroy"]; err != nil {
		return err
	}
	for key, candidate := range backend.present {
		if candidate == handle {
			delete(backend.present, key)
		}
	}
	backend.destroyed++
	return nil
}

func (backend *fakeDevelopmentLifecycle) Reconcile(
	_ context.Context,
	key OpaqueID,
) (BackendObservation, error) {
	backend.calls = append(backend.calls, "reconcile")
	if err := backend.fail["reconcile"]; err != nil {
		return BackendObservation{}, err
	}
	if backend.forceUnknown {
		return BackendObservation{Status: BackendOutcomeUnknown}, nil
	}
	if handle, exists := backend.present[key]; exists {
		return BackendObservation{Status: BackendPresent, Handle: handle}, nil
	}
	return BackendObservation{Status: BackendAuthoritativelyAbsent}, nil
}

func cloneVerifiedApproval(approval VerifiedApproval) VerifiedApproval {
	approval.CanonicalPayload = append([]byte(nil), approval.CanonicalPayload...)
	return approval
}

func repeatedID(value byte) OpaqueID {
	var id OpaqueID
	for index := range id {
		id[index] = value
	}
	return id
}

func repeatedDigest(value byte) Digest {
	var digest Digest
	for index := range digest {
		digest[index] = value
	}
	return digest
}
