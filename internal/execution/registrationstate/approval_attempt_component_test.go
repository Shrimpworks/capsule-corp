package registrationstate

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

type lockedClock struct {
	mu          sync.Mutex
	observation uint64
	err         error
}

func (clock *lockedClock) ObserveUnixSeconds(context.Context) (uint64, error) {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	return clock.observation, clock.err
}

func (clock *lockedClock) set(value uint64) {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	clock.observation = value
}

type approvalIDSequence struct {
	mu   sync.Mutex
	next uint64
	err  error
}

func (source *approvalIDSequence) NewApprovalID(context.Context) (approvalattempt.ApprovalID, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	if source.err != nil {
		return approvalattempt.ApprovalID{}, source.err
	}
	var result approvalattempt.ApprovalID
	result[0] = 0xa1
	binary.BigEndian.PutUint64(result[8:], source.next)
	source.next++
	return result, nil
}

type attemptIDSequence struct {
	mu   sync.Mutex
	next uint64
	err  error
}

func (source *attemptIDSequence) NewAttemptID(context.Context) (approvalattempt.AttemptID, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	if source.err != nil {
		return approvalattempt.AttemptID{}, source.err
	}
	var result approvalattempt.AttemptID
	result[0] = 0xb2
	binary.BigEndian.PutUint64(result[8:], source.next)
	source.next++
	return result, nil
}

type storeIntegrityAssessor struct {
	store StateStore
	err   error
	alter func(*RuntimeIntegrityAssessment)
}

func (assessor storeIntegrityAssessor) Assess(
	ctx context.Context,
	preflight IntegrityPreflight,
) (RuntimeIntegrityAssessment, error) {
	if assessor.err != nil {
		return RuntimeIntegrityAssessment{}, assessor.err
	}
	state, err := assessor.store.snapshot(ctx)
	if err != nil {
		return RuntimeIntegrityAssessment{}, err
	}
	result := RuntimeIntegrityAssessment{
		Preflight: preflight, AssessedAt: state.TimeHighWaterUnixSeconds, Permitted: true,
	}
	if assessor.alter != nil {
		assessor.alter(&result)
	}
	return result, nil
}

type approvalHarness struct {
	path           string
	store          *FixedFileStore
	clock          *lockedClock
	registrationID v0candidate.RegistrationID
	registration   PlanRegistration
	planDigest     v0candidate.ExecutionPlanDigest
}

func newApprovalHarness(t *testing.T) approvalHarness {
	t.Helper()
	path := filepath.Join(t.TempDir(), "supervisor-state.json")
	store, err := NewFixedFileStore(path, ordinaryInitialState())
	if err != nil {
		t.Fatalf("new fixed store: %v", err)
	}
	clock := &lockedClock{observation: 1_785_456_000}
	registrationID := repeated16[v0candidate.RegistrationID](0x33)
	registrationComponent := mustComponent(t, store, clock, identifierSourceFunc(
		func(context.Context) (v0candidate.RegistrationID, error) { return registrationID, nil },
	))
	registration, err := registrationComponent.RegisterPlan(
		context.Background(),
		AuthenticatedCallContext{Authenticated: true, Role: CallerDaemon, Purpose: RegisterPlanPurpose},
		readBytes(t, filepath.Join(conformanceRoot, "execution-plan/ordinary.cbor")),
		ordinaryPlanBindings(),
	)
	if err != nil {
		t.Fatalf("register plan: %v", err)
	}
	return approvalHarness{
		path: path, store: store, clock: clock, registrationID: registrationID,
		registration: registration, planDigest: registration.View().PlanDigest,
	}
}

func fixtureVector(
	t *testing.T,
	harness approvalHarness,
	nonceByte byte,
	issuedAt uint64,
	expiresAt uint64,
	variant byte,
) approvalattempt.FixtureVector {
	t.Helper()
	envelope := bytes.Clone(readBytes(t, filepath.Join(conformanceRoot, "approval-grant/ordinary.cose")))
	payload := bytes.Clone(readBytes(t, filepath.Join(conformanceRoot, "approval-grant/ordinary.payload.cbor")))
	protected := bytes.Clone(readBytes(t, filepath.Join(conformanceRoot, "approval-grant/ordinary.protected.cbor")))
	if variant != 0 {
		envelope[len(envelope)-1] ^= variant
		payload[len(payload)-1] ^= variant
	}
	return approvalattempt.FixtureVector{
		EnvelopeBytes: envelope, PayloadBytes: payload, ProtectedHeaderBytes: protected,
		ProtectedKeyID: []byte("approval-test-key"), View: approvalattempt.ApprovalGrant{
			ObjectType: approvalattempt.ApprovalGrantObjectType, ObjectVersion: 0,
			InstallationID: harness.registration.View().InstallationID,
			EpochDigest:    harness.registration.View().EpochDigest,
			RegistrationID: harness.registrationID, PlanDigest: harness.planDigest,
			SupervisorID: harness.registration.View().SupervisorID,
			AttemptNonce: repeated16[approvalattempt.AttemptNonce](nonceByte),
			Purpose:      approvalattempt.ApprovalGrantPurpose, Audience: approvalattempt.ApprovalGrantAudience,
			IssuedAt: v0candidate.UInt53(issuedAt), ExpiresAt: v0candidate.UInt53(expiresAt),
		},
		ResolvedEpochSequence: harness.registration.View().EpochSequence,
		AuthorizationIdentity: repeated32[approvalattempt.ApprovalKeyAuthorizationIdentity](0x99),
		SignatureAccepted:     true,
	}
}

func mustApprovalAttemptComponent(
	t *testing.T,
	harness approvalHarness,
	vectors []approvalattempt.FixtureVector,
	approvalIDs ApprovalIdentifierSource,
	attemptIDs AttemptIdentifierSource,
	checkpoint ApprovalAttemptCheckpointHook,
) *ApprovalAttemptComponent {
	t.Helper()
	verifier, err := approvalattempt.NewFixtureVerifier(vectors)
	if err != nil {
		t.Fatalf("new fixture verifier: %v", err)
	}
	component, err := NewApprovalAttempt(ApprovalAttemptOptions{
		Store: harness.store, Clock: harness.clock, Verifier: verifier,
		ApprovalIdentifiers: approvalIDs, AttemptIdentifiers: attemptIDs,
		Integrity: storeIntegrityAssessor{store: harness.store}, Checkpoint: checkpoint,
	})
	if err != nil {
		t.Fatalf("new approval/attempt component: %v", err)
	}
	return component
}

func submitCall() AuthenticatedCallContext {
	return AuthenticatedCallContext{Authenticated: true, Role: CallerBroker, Purpose: SubmitApprovalPurpose}
}

func attemptCall() AuthenticatedCallContext {
	return AuthenticatedCallContext{Authenticated: true, Role: CallerDaemon, Purpose: RequestAttemptPurpose}
}

func TestApprovalAttemptDurableIdempotentAtomicHappyPath(t *testing.T) {
	harness := newApprovalHarness(t)
	first := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
	equivalent := first
	equivalent.EnvelopeBytes = bytes.Clone(first.EnvelopeBytes)
	equivalent.EnvelopeBytes[len(equivalent.EnvelopeBytes)-1] ^= 0x80
	component := mustApprovalAttemptComponent(
		t, harness, []approvalattempt.FixtureVector{first, equivalent},
		&approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil,
	)

	callerEnvelope := bytes.Clone(first.EnvelopeBytes)
	submission, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, callerEnvelope)
	if err != nil {
		t.Fatalf("submit approval: %v", err)
	}
	callerEnvelope[0] ^= 0xff
	replay, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, equivalent.EnvelopeBytes)
	if err != nil {
		t.Fatalf("equivalent replay: %v", err)
	}
	if replay.Reference != submission.Reference || replay.State != approvalattempt.ApprovalUsable {
		t.Fatalf("equivalent replay = %#v, want original usable reference", replay)
	}
	record, err := component.Approval(context.Background(), submission.Reference)
	if err != nil {
		t.Fatalf("read approval: %v", err)
	}
	if !bytes.Equal(record.ExactEnvelopeBytes, first.EnvelopeBytes) {
		t.Fatal("equivalent replay or caller mutation replaced the first exact envelope")
	}
	returned := record.ExactEnvelopeBytes
	returned[0] ^= 0xff
	again, _ := component.Approval(context.Background(), submission.Reference)
	if !bytes.Equal(again.ExactEnvelopeBytes, first.EnvelopeBytes) {
		t.Fatal("approval inspection did not return a defensive copy")
	}

	created, err := component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
	if err != nil {
		t.Fatalf("request attempt: %v", err)
	}
	replayedAttempt, err := component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
	if err != nil || replayedAttempt != created {
		t.Fatalf("attempt replay = %#v, %v; want %#v", replayedAttempt, err, created)
	}
	consumed, _ := component.Approval(context.Background(), submission.Reference)
	attempt, _ := component.Attempt(context.Background(), created.Reference)
	if consumed.State != approvalattempt.ApprovalConsumed || consumed.ConsumedAttemptID != attempt.AttemptID ||
		attempt.ApprovalID != consumed.ApprovalID || attempt.AttemptNonce != consumed.AttemptNonce ||
		attempt.RegistrationID != consumed.RegistrationID || attempt.RegistrationSequence != consumed.RegistrationSequence ||
		attempt.PlanDigest != consumed.PlanDigest || attempt.InstallationID != consumed.InstallationID ||
		attempt.EpochSequence != consumed.EpochSequence || attempt.EpochDigest != consumed.EpochDigest ||
		attempt.SupervisorID != consumed.SupervisorID || attempt.Purpose != consumed.Purpose ||
		attempt.Audience != consumed.Audience || attempt.ApprovalPayloadDigest != consumed.PayloadDigest ||
		attempt.AuthorizationIdentity != consumed.AuthorizationIdentity || attempt.CreatedAt != 1_785_456_000 {
		t.Fatalf("atomic copied bindings mismatch:\napproval=%#v\nattempt=%#v", consumed, attempt)
	}
	if bytes.Equal(attempt.AttemptID[:], attempt.AttemptNonce[:]) {
		t.Fatal("attempt ID was not independent of the signed nonce")
	}

	reopened, err := NewFixedFileStore(harness.path, InitialState{})
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	harness.store = reopened
	restarted := mustApprovalAttemptComponent(
		t, harness, []approvalattempt.FixtureVector{first, equivalent},
		&approvalIDSequence{next: 99}, &attemptIDSequence{next: 99}, nil,
	)
	afterRestart, err := restarted.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
	if err != nil || afterRestart != created {
		t.Fatalf("restart replay = %#v, %v; want %#v", afterRestart, err, created)
	}
	createdReferences, err := restarted.CreatedAttempts(context.Background())
	if err != nil || len(createdReferences) != 1 || createdReferences[0] != created.Reference {
		t.Fatalf("startup enumeration = %#v, %v", createdReferences, err)
	}
}

func TestResolveCreatedRefusesCrossLinkedDurableState(t *testing.T) {
	harness := newApprovalHarness(t)
	vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
	component := mustApprovalAttemptComponent(
		t, harness, []approvalattempt.FixtureVector{vector},
		&approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil,
	)
	submission, err := component.SubmitApproval(
		context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes,
	)
	if err != nil {
		t.Fatalf("submit approval: %v", err)
	}
	created, err := component.RequestAttempt(
		context.Background(), attemptCall(), harness.registrationID, submission.Reference,
	)
	if err != nil {
		t.Fatalf("request attempt: %v", err)
	}
	harness.store.mu.Lock()
	harness.store.state.Attempts[0].ApprovalID[0] ^= 0xff
	harness.store.mu.Unlock()

	_, err = component.ResolveCreated(context.Background(), created.Reference.AttemptID())
	assertApprovalClassification(t, err, approvalattempt.ClassificationRecoveryRequired)
}

func TestApprovalAttemptConcurrentExactRequestsConverge(t *testing.T) {
	const workers = 24
	harness := newApprovalHarness(t)
	vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
	component := mustApprovalAttemptComponent(
		t, harness, []approvalattempt.FixtureVector{vector},
		&approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil,
	)
	submissions := make(chan approvalattempt.ApprovalSubmission, workers)
	errorsSeen := make(chan error, workers)
	var group sync.WaitGroup
	for range workers {
		group.Add(1)
		go func() {
			defer group.Done()
			result, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
			if err != nil {
				errorsSeen <- err
				return
			}
			submissions <- result
		}()
	}
	group.Wait()
	close(submissions)
	close(errorsSeen)
	for err := range errorsSeen {
		t.Fatalf("concurrent submission: %v", err)
	}
	var reference approvalattempt.ApprovalReference
	for result := range submissions {
		if reference == (approvalattempt.ApprovalReference{}) {
			reference = result.Reference
		} else if result.Reference != reference {
			t.Fatal("concurrent submissions returned different approvals")
		}
	}
	state, _ := harness.store.snapshot(context.Background())
	if len(state.Approvals) != 1 {
		t.Fatalf("concurrent submissions created %d approvals", len(state.Approvals))
	}

	attempts := make(chan approvalattempt.AttemptCreation, workers)
	errorsSeen = make(chan error, workers)
	for range workers {
		group.Add(1)
		go func() {
			defer group.Done()
			result, err := component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, reference)
			if err != nil {
				errorsSeen <- err
				return
			}
			attempts <- result
		}()
	}
	group.Wait()
	close(attempts)
	close(errorsSeen)
	for err := range errorsSeen {
		t.Fatalf("concurrent attempt: %v", err)
	}
	var attemptReference approvalattempt.AttemptReference
	for result := range attempts {
		if attemptReference == (approvalattempt.AttemptReference{}) {
			attemptReference = result.Reference
		} else if result.Reference != attemptReference {
			t.Fatal("concurrent requests returned different attempts")
		}
	}
	state, _ = harness.store.snapshot(context.Background())
	if len(state.Attempts) != 1 || state.Approvals[0].State != approvalattempt.ApprovalConsumed {
		t.Fatalf("concurrent consume/create split state: %#v", state)
	}
}

func TestApprovalSubmissionTimeBoundaries(t *testing.T) {
	tests := []struct {
		name    string
		now     uint64
		issued  uint64
		expires uint64
		want    approvalattempt.Classification
		accept  bool
	}{
		{name: "issued equality", now: 1_785_456_000, issued: 1_785_456_000, expires: 1_785_456_001, accept: true},
		{name: "issued future", now: 1_785_456_000, issued: 1_785_456_001, expires: 1_785_456_002, want: approvalattempt.ClassificationStale},
		{name: "expiry equality", now: 1_785_456_000, issued: 1_785_455_999, expires: 1_785_456_000, want: approvalattempt.ClassificationStale},
		{name: "expiry plus one", now: 1_785_456_000, issued: 1_785_456_000, expires: 1_785_456_001, accept: true},
		{name: "equal issue expiry", now: 1_785_456_000, issued: 1_785_456_000, expires: 1_785_456_000, want: approvalattempt.ClassificationSchema},
		{name: "lifetime exact", now: 1_785_456_000, issued: 1_785_456_000, expires: 1_785_456_300, accept: true},
		{name: "lifetime over", now: 1_785_456_000, issued: 1_785_455_999, expires: 1_785_456_300, want: approvalattempt.ClassificationStale},
		{name: "registration expiry equality", now: 1_785_456_000, issued: 1_785_456_000, expires: 1_785_456_300, accept: true},
		{name: "registration expiry over", now: 1_785_456_000, issued: 1_785_456_000, expires: 1_785_456_301, want: approvalattempt.ClassificationStale},
	}
	for index, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			harness := newApprovalHarness(t)
			harness.clock.set(test.now)
			vector := fixtureVector(t, harness, byte(0x60+index), test.issued, test.expires, byte(index+1))
			component := mustApprovalAttemptComponent(
				t, harness, []approvalattempt.FixtureVector{vector},
				&approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil,
			)
			_, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
			if test.accept {
				if err != nil {
					t.Fatalf("accepted boundary rejected: %v", err)
				}
				return
			}
			classification, ok := approvalattempt.ErrorClassification(err)
			if !ok || classification != test.want {
				t.Fatalf("classification = %q (%t), want %q: %v", classification, ok, test.want, err)
			}
			state, _ := harness.store.snapshot(context.Background())
			if len(state.Approvals) != 0 || len(state.Attempts) != 0 {
				t.Fatal("time rejection changed approval/attempt authority")
			}
		})
	}
}

func TestApprovalAttemptFaultAndProcessDeathMatrix(t *testing.T) {
	t.Run("approval confirmed abort", func(t *testing.T) {
		harness := newApprovalHarness(t)
		vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
		component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
		harness.store.InjectFailure(FaultApprovalCommitAbort, errors.New("abort"))
		_, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
		assertApprovalClassification(t, err, approvalattempt.ClassificationLocalFailure)
		state, _ := harness.store.snapshot(context.Background())
		if len(state.Approvals) != 0 || len(state.Attempts) != 0 || harness.store.recoveryFenced() {
			t.Fatal("confirmed approval abort changed or fenced authority")
		}
	})

	t.Run("approval indeterminate rename", func(t *testing.T) {
		harness := newApprovalHarness(t)
		vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
		component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
		harness.store.InjectFailure(FaultApprovalCommitIndeterminate, errors.New("directory sync"))
		_, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
		assertApprovalClassification(t, err, approvalattempt.ClassificationRecoveryRequired)
		if !harness.store.recoveryFenced() {
			t.Fatal("indeterminate approval commit did not fence mutation")
		}
		_, err = component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
		assertApprovalClassification(t, err, approvalattempt.ClassificationRecoveryRequired)
		reopened, reopenErr := NewFixedFileStore(harness.path, InitialState{})
		if reopenErr != nil {
			t.Fatalf("reopen indeterminate approval: %v", reopenErr)
		}
		harness.store = reopened
		restarted := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 9}, &attemptIDSequence{next: 9}, nil)
		result, retryErr := restarted.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
		if retryErr != nil || result.State != approvalattempt.ApprovalUsable {
			t.Fatalf("retry after reopen = %#v, %v", result, retryErr)
		}
		state, _ := reopened.snapshot(context.Background())
		if len(state.Approvals) != 1 {
			t.Fatalf("indeterminate approval recovered %d records", len(state.Approvals))
		}
	})

	t.Run("approval indeterminate pre-state", func(t *testing.T) {
		harness := newApprovalHarness(t)
		vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
		component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
		harness.store.InjectFailure(FaultApprovalCommitIndeterminatePreState, errors.New("power loss"))
		_, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
		assertApprovalClassification(t, err, approvalattempt.ClassificationRecoveryRequired)
		reopened, reopenErr := NewFixedFileStore(harness.path, InitialState{})
		if reopenErr != nil {
			t.Fatalf("reopen approval pre-state: %v", reopenErr)
		}
		harness.store = reopened
		state, _ := reopened.snapshot(context.Background())
		if len(state.Approvals) != 0 || len(state.Attempts) != 0 {
			t.Fatal("indeterminate approval pre-state created authority")
		}
		restarted := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 9}, &attemptIDSequence{next: 9}, nil)
		if _, retryErr := restarted.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes); retryErr != nil {
			t.Fatalf("retry approval pre-state: %v", retryErr)
		}
	})

	t.Run("attempt abort and indeterminate", func(t *testing.T) {
		harness := newApprovalHarness(t)
		vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
		component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
		submission, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
		if err != nil {
			t.Fatal(err)
		}
		harness.store.InjectFailure(FaultAttemptCommitAbort, errors.New("abort"))
		_, err = component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
		assertApprovalClassification(t, err, approvalattempt.ClassificationLocalFailure)
		record, _ := component.Approval(context.Background(), submission.Reference)
		if record.State != approvalattempt.ApprovalUsable {
			t.Fatal("confirmed attempt abort consumed approval")
		}
		harness.store.InjectFailure(FaultAttemptCommitIndeterminate, errors.New("directory sync"))
		_, err = component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
		assertApprovalClassification(t, err, approvalattempt.ClassificationRecoveryRequired)
		reopened, reopenErr := NewFixedFileStore(harness.path, InitialState{})
		if reopenErr != nil {
			t.Fatalf("reopen indeterminate attempt: %v", reopenErr)
		}
		harness.store = reopened
		restarted := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 9}, &attemptIDSequence{next: 9}, nil)
		result, retryErr := restarted.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
		if retryErr != nil || result.State != approvalattempt.AttemptCreated {
			t.Fatalf("retry indeterminate attempt = %#v, %v", result, retryErr)
		}
		state, _ := reopened.snapshot(context.Background())
		if len(state.Attempts) != 1 || state.Approvals[0].State != approvalattempt.ApprovalConsumed {
			t.Fatal("indeterminate attempt reopened into split state")
		}
	})

	t.Run("attempt indeterminate pre-state", func(t *testing.T) {
		harness := newApprovalHarness(t)
		vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
		component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
		submission, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
		if err != nil {
			t.Fatal(err)
		}
		harness.store.InjectFailure(FaultAttemptCommitIndeterminatePreState, errors.New("power loss"))
		_, err = component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
		assertApprovalClassification(t, err, approvalattempt.ClassificationRecoveryRequired)
		if _, inspectErr := component.Approval(context.Background(), submission.Reference); inspectErr == nil {
			t.Fatal("fenced component exposed indeterminate approval state")
		}
		reopened, reopenErr := NewFixedFileStore(harness.path, InitialState{})
		if reopenErr != nil {
			t.Fatalf("reopen attempt pre-state: %v", reopenErr)
		}
		harness.store = reopened
		restarted := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 9}, &attemptIDSequence{next: 9}, nil)
		record, _ := restarted.Approval(context.Background(), submission.Reference)
		if record.State != approvalattempt.ApprovalUsable {
			t.Fatal("indeterminate attempt pre-state did not preserve usable approval")
		}
		if _, retryErr := restarted.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference); retryErr != nil {
			t.Fatalf("retry attempt pre-state: %v", retryErr)
		}
	})

	t.Run("response loss after durable commits", func(t *testing.T) {
		harness := newApprovalHarness(t)
		vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
		var mu sync.Mutex
		seen := map[ApprovalAttemptCheckpoint]int{}
		checkpoint := func(_ context.Context, point ApprovalAttemptCheckpoint) error {
			mu.Lock()
			defer mu.Unlock()
			seen[point]++
			if point == CheckpointAfterApprovalCommit && seen[point] == 1 {
				return errors.New("lost approval response")
			}
			if point == CheckpointAfterAttemptCommit && seen[point] == 1 {
				return errors.New("lost attempt response")
			}
			return nil
		}
		component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, checkpoint)
		_, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
		assertApprovalClassification(t, err, approvalattempt.ClassificationLocalFailure)
		submission, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
		if err != nil {
			t.Fatalf("approval response-loss retry: %v", err)
		}
		_, err = component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
		assertApprovalClassification(t, err, approvalattempt.ClassificationLocalFailure)
		created, err := component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
		if err != nil || created.State != approvalattempt.AttemptCreated {
			t.Fatalf("attempt response-loss retry = %#v, %v", created, err)
		}
	})
}

func TestApprovalAttemptReopenRejectsCorruptionAndCrossLinks(t *testing.T) {
	harness := newApprovalHarness(t)
	vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
	component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
	submission, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference); err != nil {
		t.Fatal(err)
	}
	base, _ := harness.store.snapshot(context.Background())
	tests := []struct {
		name   string
		mutate func(*installationState)
	}{
		{name: "consumed without attempt", mutate: func(state *installationState) {
			state.Attempts = nil
			attemptDigest, err := attemptSetDigest(state.Attempts)
			if err != nil {
				t.Fatalf("attempt set digest: %v", err)
			}
			state.AttemptSetDigest = attemptDigest
		}},
		{name: "attempt without consumption", mutate: func(state *installationState) {
			state.Approvals[0].State = approvalattempt.ApprovalUsable
			state.Approvals[0].ConsumedAttemptID = approvalattempt.AttemptID{}
			approvalDigest, err := approvalSetDigest(state.Approvals)
			if err != nil {
				t.Fatalf("approval set digest: %v", err)
			}
			state.ApprovalSetDigest = approvalDigest
		}},
		{name: "cross linked approval", mutate: func(state *installationState) {
			state.Attempts[0].ApprovalID[0] ^= 0xff
			attemptDigest, err := attemptSetDigest(state.Attempts)
			if err != nil {
				t.Fatalf("attempt set digest: %v", err)
			}
			state.AttemptSetDigest = attemptDigest
		}},
		{name: "copied nonce mismatch", mutate: func(state *installationState) {
			state.Attempts[0].AttemptNonce[0] ^= 0xff
			attemptDigest, err := attemptSetDigest(state.Attempts)
			if err != nil {
				t.Fatalf("attempt set digest: %v", err)
			}
			state.AttemptSetDigest = attemptDigest
		}},
		{name: "payload digest corruption", mutate: func(state *installationState) {
			state.Approvals[0].PayloadDigest[0] ^= 0xff
			approvalDigest, err := approvalSetDigest(state.Approvals)
			if err != nil {
				t.Fatalf("approval set digest: %v", err)
			}
			state.ApprovalSetDigest = approvalDigest
		}},
		{name: "approval record version", mutate: func(state *installationState) {
			state.Approvals[0].StorageFormatVersion++
			approvalDigest, err := approvalSetDigest(state.Approvals)
			if err != nil {
				t.Fatalf("approval set digest: %v", err)
			}
			state.ApprovalSetDigest = approvalDigest
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			state := cloneState(base)
			test.mutate(&state)
			path := filepath.Join(t.TempDir(), "corrupt.json")
			encoded, encodeErr := jsonMarshalDiskState(state)
			if encodeErr != nil {
				t.Fatalf("encode corrupt state: %v", encodeErr)
			}
			if writeErr := os.WriteFile(path, encoded, 0o600); writeErr != nil {
				t.Fatalf("write corrupt state: %v", writeErr)
			}
			if _, reopenErr := NewFixedFileStore(path, InitialState{}); reopenErr == nil {
				t.Fatal("corrupt/cross-linked store reopened")
			} else if !errors.Is(reopenErr, ErrStoreRepairRequired) {
				t.Fatalf("reopen error = %v, want ErrStoreRepairRequired", reopenErr)
			}
		})
	}
}

func TestFixedStoreRejectsUnsafeFilePermissionsWithoutRewrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "permissions-v0.json")
	if _, err := NewFixedFileStore(path, ordinaryInitialState()); err != nil {
		t.Fatalf("create store: %v", err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}
	before := mustReadFile(t, path)
	if _, err := NewFixedFileStore(path, InitialState{}); err == nil {
		t.Fatal("v0 open accepted unsafe file permissions")
	}
	if after := mustReadFile(t, path); !bytes.Equal(after, before) {
		t.Fatal("unsafe file refusal rewrote evidence")
	}
}

func TestApprovalAttemptAuthenticationAndPreflightFailuresDoNotChangeAuthority(t *testing.T) {
	harness := newApprovalHarness(t)
	vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
	component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
	before, _ := harness.store.snapshot(context.Background())
	for _, call := range []AuthenticatedCallContext{
		{},
		{Authenticated: true, Role: CallerDaemon, Purpose: SubmitApprovalPurpose},
		{Authenticated: true, Role: CallerUpdater, Purpose: SubmitApprovalPurpose},
		{Authenticated: true, Role: CallerBroker, Purpose: "wrong"},
	} {
		_, err := component.SubmitApproval(context.Background(), call, harness.registrationID, vector.EnvelopeBytes)
		assertApprovalClassification(t, err, approvalattempt.ClassificationAuthentication)
	}
	after, _ := harness.store.snapshot(context.Background())
	if !reflect.DeepEqual(before, after) {
		t.Fatal("submission authentication failure changed durable state")
	}
	submission, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
	if err != nil {
		t.Fatal(err)
	}
	for _, call := range []AuthenticatedCallContext{
		{},
		{Authenticated: true, Role: CallerBroker, Purpose: RequestAttemptPurpose},
		{Authenticated: true, Role: CallerUpdater, Purpose: RequestAttemptPurpose},
		{Authenticated: true, Role: CallerDaemon, Purpose: "wrong"},
	} {
		_, err := component.RequestAttempt(context.Background(), call, harness.registrationID, submission.Reference)
		assertApprovalClassification(t, err, approvalattempt.ClassificationAuthentication)
	}
	record, _ := component.Approval(context.Background(), submission.Reference)
	if record.State != approvalattempt.ApprovalUsable {
		t.Fatal("attempt authentication failure consumed approval")
	}

	verifier, _ := approvalattempt.NewFixtureVerifier([]approvalattempt.FixtureVector{vector})
	failing, err := NewApprovalAttempt(ApprovalAttemptOptions{
		Store: harness.store, Clock: harness.clock, Verifier: verifier,
		ApprovalIdentifiers: &approvalIDSequence{next: 9}, AttemptIdentifiers: &attemptIDSequence{next: 9},
		Integrity: storeIntegrityAssessor{store: harness.store, err: errors.New("integrity denied")},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = failing.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
	assertApprovalClassification(t, err, approvalattempt.ClassificationTrustState)
	record, _ = failing.Approval(context.Background(), submission.Reference)
	if record.State != approvalattempt.ApprovalUsable {
		t.Fatal("integrity preflight failure consumed approval")
	}
}

func TestApprovalAttemptExpiryRollbackNonceAndIdentifierRules(t *testing.T) {
	t.Run("final expiry and rollback", func(t *testing.T) {
		harness := newApprovalHarness(t)
		vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
		component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
		submission, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
		if err != nil {
			t.Fatal(err)
		}
		harness.clock.set(1_785_456_300)
		_, err = component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
		assertApprovalClassification(t, err, approvalattempt.ClassificationStale)
		harness.clock.set(1_785_456_000)
		_, err = component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
		assertApprovalClassification(t, err, approvalattempt.ClassificationStale)
		record, _ := component.Approval(context.Background(), submission.Reference)
		if record.State != approvalattempt.ApprovalUsable {
			t.Fatal("expiry or rollback rejection consumed the approval")
		}
		state, _ := harness.store.snapshot(context.Background())
		if state.TimeHighWaterUnixSeconds != 1_785_456_300 {
			t.Fatal("clock rollback reduced durable effective time")
		}
	})

	t.Run("nonce replay and ID collisions", func(t *testing.T) {
		harness := newApprovalHarness(t)
		first := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 1)
		nonceReplay := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 2)
		component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{first, nonceReplay}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
		submission, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, first.EnvelopeBytes)
		if err != nil {
			t.Fatal(err)
		}
		_, err = component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, nonceReplay.EnvelopeBytes)
		assertApprovalClassification(t, err, approvalattempt.ClassificationReplay)
		state, _ := harness.store.snapshot(context.Background())
		if len(state.Approvals) != 1 {
			t.Fatal("nonce replay created an approval")
		}

		second := fixtureVector(t, harness, 0x67, 1_785_456_000, 1_785_456_300, 3)
		verifier, _ := approvalattempt.NewFixtureVerifier([]approvalattempt.FixtureVector{first, second})
		collision, _ := NewApprovalAttempt(ApprovalAttemptOptions{
			Store: harness.store, Clock: harness.clock, Verifier: verifier,
			ApprovalIdentifiers: approvalIDSequenceFunc(func(context.Context) (approvalattempt.ApprovalID, error) {
				return submission.Reference.ApprovalID(), nil
			}),
			AttemptIdentifiers: &attemptIDSequence{next: 1}, Integrity: storeIntegrityAssessor{store: harness.store},
		})
		_, err = collision.SubmitApproval(context.Background(), submitCall(), harness.registrationID, second.EnvelopeBytes)
		assertApprovalClassification(t, err, approvalattempt.ClassificationLocalFailure)
		state, _ = harness.store.snapshot(context.Background())
		if len(state.Approvals) != 1 {
			t.Fatal("approval ID collision created a record")
		}
	})
}

func TestApprovalAttemptHighWaterAndAttemptCapacity(t *testing.T) {
	t.Run("high-water confirmed and indeterminate failures", func(t *testing.T) {
		for _, test := range []struct {
			name  string
			fault StoreFault
			want  approvalattempt.Classification
			fence bool
		}{
			{name: "confirmed", fault: FaultTimeHighWaterWrite, want: approvalattempt.ClassificationLocalFailure},
			{name: "indeterminate pre-state", fault: FaultTimeHighWaterIndeterminatePreState, want: approvalattempt.ClassificationRecoveryRequired, fence: true},
			{name: "indeterminate", fault: FaultTimeHighWaterCommitIndeterminate, want: approvalattempt.ClassificationRecoveryRequired, fence: true},
		} {
			t.Run(test.name, func(t *testing.T) {
				harness := newApprovalHarness(t)
				vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
				component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
				harness.store.InjectFailure(test.fault, errors.New("time fault"))
				_, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
				assertApprovalClassification(t, err, test.want)
				state, _ := harness.store.snapshot(context.Background())
				if len(state.Approvals) != 0 || len(state.Attempts) != 0 || harness.store.recoveryFenced() != test.fence {
					t.Fatal("high-water fault changed authority or used the wrong fence")
				}
			})
		}
	})

	// This compact population proves the exact open-attempt ceiling without
	// retaining hundreds of opaque fixture files or performing backend work.
	t.Run("nonterminal attempt exact ceiling and cap plus one", func(t *testing.T) {
		harness := newApprovalHarness(t)
		vectors := make([]approvalattempt.FixtureVector, 0, MaxNonterminalAttempts+1)
		for index := 0; index <= MaxNonterminalAttempts; index++ {
			vectors = append(vectors, fixtureVector(
				t, harness, byte(index%255+1), 1_785_456_000, 1_785_456_300, byte(index%255+1),
			))
			// Make every fixture and replay identity unique beyond the one-byte
			// mutations while preserving the fixed verifier's bounded copies.
			binary.BigEndian.PutUint16(vectors[index].EnvelopeBytes[len(vectors[index].EnvelopeBytes)-2:], uint16(index+1))
			binary.BigEndian.PutUint16(vectors[index].PayloadBytes[len(vectors[index].PayloadBytes)-2:], uint16(index+1))
			binary.BigEndian.PutUint16(vectors[index].View.AttemptNonce[14:], uint16(index+1))
		}
		component := mustApprovalAttemptComponent(t, harness, vectors, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
		for index := 0; index < MaxNonterminalAttempts; index++ {
			submission, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vectors[index].EnvelopeBytes)
			if err != nil {
				t.Fatalf("submit %d: %v", index, err)
			}
			if _, err := component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference); err != nil {
				t.Fatalf("attempt %d: %v", index, err)
			}
		}
		extra, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vectors[MaxNonterminalAttempts].EnvelopeBytes)
		if err != nil {
			t.Fatalf("submit capacity probe: %v", err)
		}
		_, err = component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, extra.Reference)
		assertApprovalClassification(t, err, approvalattempt.ClassificationCapacity)
		record, _ := component.Approval(context.Background(), extra.Reference)
		if record.State != approvalattempt.ApprovalUsable {
			t.Fatal("attempt capacity rejection consumed approval")
		}
		state, _ := harness.store.snapshot(context.Background())
		if len(state.Attempts) != MaxNonterminalAttempts {
			t.Fatalf("attempt count = %d, want %d", len(state.Attempts), MaxNonterminalAttempts)
		}
	})
}

func TestApprovalCapacityAndNonEvictingRetainedState(t *testing.T) {
	t.Run("usable approval exact ceiling and cap plus one", func(t *testing.T) {
		harness := newApprovalHarness(t)
		vectors := make([]approvalattempt.FixtureVector, MaxUsableApprovals+1)
		for index := range vectors {
			vectors[index] = fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 1)
			binary.BigEndian.PutUint16(vectors[index].EnvelopeBytes[len(vectors[index].EnvelopeBytes)-2:], uint16(index+1))
			binary.BigEndian.PutUint16(vectors[index].PayloadBytes[len(vectors[index].PayloadBytes)-2:], uint16(index+1))
			binary.BigEndian.PutUint16(vectors[index].View.AttemptNonce[14:], uint16(index+1))
		}
		component := mustApprovalAttemptComponent(t, harness, vectors, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
		for index := 0; index < MaxUsableApprovals; index++ {
			if _, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vectors[index].EnvelopeBytes); err != nil {
				t.Fatalf("submit %d: %v", index, err)
			}
		}
		_, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vectors[MaxUsableApprovals].EnvelopeBytes)
		assertApprovalClassification(t, err, approvalattempt.ClassificationCapacity)
		state, _ := harness.store.snapshot(context.Background())
		if len(state.Approvals) != MaxUsableApprovals || countUsableApprovals(state.Approvals, state.TimeHighWaterUnixSeconds) != MaxUsableApprovals {
			t.Fatalf("usable approval capacity changed or evicted state: %d", len(state.Approvals))
		}
	})

	t.Run("retained exact ceiling survives reopen and rejects cap plus one", func(t *testing.T) {
		harness := newApprovalHarness(t)
		vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 1)
		component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
		submission, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
		if err != nil {
			t.Fatal(err)
		}
		template, _ := component.Approval(context.Background(), submission.Reference)
		state, _ := harness.store.snapshot(context.Background())
		state.Approvals = make([]approvalattempt.ApprovalRecord, MaxRetainedApprovals)
		for index := range state.Approvals {
			record := approvalattempt.CloneApprovalRecord(template)
			record.State = approvalattempt.ApprovalInvalidated
			record.ConsumedAttemptID = approvalattempt.AttemptID{}
			record.ApprovalID = approvalattempt.ApprovalID{0xa1}
			binary.BigEndian.PutUint16(record.ApprovalID[14:], uint16(index+1))
			binary.BigEndian.PutUint16(record.AttemptNonce[14:], uint16(index+1))
			binary.BigEndian.PutUint16(record.ExactEnvelopeBytes[len(record.ExactEnvelopeBytes)-2:], uint16(index+1))
			binary.BigEndian.PutUint16(record.ExactPayloadBytes[len(record.ExactPayloadBytes)-2:], uint16(index+1))
			record.EnvelopeDigest = approvalattempt.ApprovalEnvelopeDigest(sha256.Sum256(record.ExactEnvelopeBytes))
			record.PayloadDigest = approvalattempt.ApprovalPayloadDigest(sha256.Sum256(record.ExactPayloadBytes))
			state.Approvals[index] = record
		}
		approvalDigest, err := approvalSetDigest(state.Approvals)
		if err != nil {
			t.Fatalf("approval set digest: %v", err)
		}
		state.ApprovalSetDigest = approvalDigest
		if err := validateState(state); err != nil {
			t.Fatalf("exact retained approval population invalid: %v", err)
		}
		path := filepath.Join(t.TempDir(), "retained-capacity.json")
		if err := persistState(path, state); err != nil {
			t.Fatalf("persist retained population: %v", err)
		}
		reopened, err := NewFixedFileStore(path, InitialState{})
		if err != nil {
			t.Fatalf("reopen exact retained population: %v", err)
		}
		harness.path = path
		harness.store = reopened
		newVector := fixtureVector(t, harness, 0x67, 1_785_456_000, 1_785_456_300, 2)
		component = mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{newVector}, &approvalIDSequence{next: 9_000}, &attemptIDSequence{next: 1}, nil)
		_, err = component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, newVector.EnvelopeBytes)
		assertApprovalClassification(t, err, approvalattempt.ClassificationCapacity)
		after, _ := reopened.snapshot(context.Background())
		if len(after.Approvals) != MaxRetainedApprovals || after.ApprovalSetDigest != state.ApprovalSetDigest {
			t.Fatal("retained approval cap-plus-one evicted or rewrote state")
		}
	})
}

func TestApprovalSubmissionBindingMatrix(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*approvalattempt.FixtureVector)
		want   approvalattempt.Classification
	}{
		{name: "object", mutate: func(vector *approvalattempt.FixtureVector) { vector.View.ObjectType = "capsule.execution-attempt" }, want: approvalattempt.ClassificationUnsupported},
		{name: "version", mutate: func(vector *approvalattempt.FixtureVector) { vector.View.ObjectVersion = 1 }, want: approvalattempt.ClassificationUnsupported},
		{name: "registration", mutate: func(vector *approvalattempt.FixtureVector) { vector.View.RegistrationID[0] ^= 0xff }, want: approvalattempt.ClassificationBinding},
		{name: "plan digest", mutate: func(vector *approvalattempt.FixtureVector) { vector.View.PlanDigest[0] ^= 0xff }, want: approvalattempt.ClassificationBinding},
		{name: "installation", mutate: func(vector *approvalattempt.FixtureVector) { vector.View.InstallationID[0] ^= 0xff }, want: approvalattempt.ClassificationBinding},
		{name: "epoch digest", mutate: func(vector *approvalattempt.FixtureVector) { vector.View.EpochDigest[0] ^= 0xff }, want: approvalattempt.ClassificationBinding},
		{name: "epoch sequence resolution", mutate: func(vector *approvalattempt.FixtureVector) { vector.ResolvedEpochSequence++ }, want: approvalattempt.ClassificationBinding},
		{name: "supervisor", mutate: func(vector *approvalattempt.FixtureVector) { vector.View.SupervisorID[0] ^= 0xff }, want: approvalattempt.ClassificationBinding},
		{name: "purpose", mutate: func(vector *approvalattempt.FixtureVector) { vector.View.Purpose = "capsule.execution.attest" }, want: approvalattempt.ClassificationBinding},
		{name: "audience", mutate: func(vector *approvalattempt.FixtureVector) { vector.View.Audience = "capsule.agent-daemon" }, want: approvalattempt.ClassificationBinding},
	}
	for index, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			harness := newApprovalHarness(t)
			vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, byte(index+1))
			test.mutate(&vector)
			component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
			_, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
			assertApprovalClassification(t, err, test.want)
			state, _ := harness.store.snapshot(context.Background())
			if len(state.Approvals) != 0 || len(state.Attempts) != 0 {
				t.Fatal("binding rejection created approval/attempt authority")
			}
		})
	}
}

func TestAttemptAdmissionTrustBindingAndIdentifierFailures(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*installationState)
		want   approvalattempt.Classification
	}{
		{name: "transition fenced", mutate: func(state *installationState) { state.TrustPhase = TrustTransitionFenced }, want: approvalattempt.ClassificationTrustState},
		{name: "repair required", mutate: func(state *installationState) {
			state.TrustPhase = TrustRepairRequired
			state.TrustReason = "fixture-repair"
		}, want: approvalattempt.ClassificationTrustState},
		{name: "quarantined", mutate: func(state *installationState) { state.Quarantined = true }, want: approvalattempt.ClassificationTrustState},
		{name: "attempts disabled", mutate: func(state *installationState) { state.AttemptsDisabled = true }, want: approvalattempt.ClassificationTrustState},
		{name: "approval invalidated", mutate: func(state *installationState) {
			state.Approvals[0].State = approvalattempt.ApprovalInvalidated
			approvalDigest, err := approvalSetDigest(state.Approvals)
			if err != nil {
				t.Fatalf("approval set digest: %v", err)
			}
			state.ApprovalSetDigest = approvalDigest
		}, want: approvalattempt.ClassificationTrustState},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			harness := newApprovalHarness(t)
			vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
			component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
			submission, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
			if err != nil {
				t.Fatal(err)
			}
			if err := harness.store.commitApproval(context.Background(), func(state *installationState) error { test.mutate(state); return nil }); err != nil {
				t.Fatalf("apply trust fixture: %v", err)
			}
			_, err = component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, submission.Reference)
			assertApprovalClassification(t, err, test.want)
			record, _ := component.Approval(context.Background(), submission.Reference)
			if record.State == approvalattempt.ApprovalConsumed {
				t.Fatal("trust/binding failure consumed approval")
			}
		})
	}

	t.Run("cross registration reference", func(t *testing.T) {
		harness := newApprovalHarness(t)
		vector := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 0)
		component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{vector}, &approvalIDSequence{next: 1}, &attemptIDSequence{next: 1}, nil)
		submission, err := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, vector.EnvelopeBytes)
		if err != nil {
			t.Fatal(err)
		}
		wrong := harness.registrationID
		wrong[0] ^= 0xff
		_, err = component.RequestAttempt(context.Background(), attemptCall(), wrong, submission.Reference)
		assertApprovalClassification(t, err, approvalattempt.ClassificationBinding)
	})

	t.Run("attempt identifier collision", func(t *testing.T) {
		harness := newApprovalHarness(t)
		first := fixtureVector(t, harness, 0x66, 1_785_456_000, 1_785_456_300, 1)
		second := fixtureVector(t, harness, 0x67, 1_785_456_000, 1_785_456_300, 2)
		attemptIDs := &attemptIDSequence{next: 1}
		component := mustApprovalAttemptComponent(t, harness, []approvalattempt.FixtureVector{first, second}, &approvalIDSequence{next: 1}, attemptIDs, nil)
		one, _ := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, first.EnvelopeBytes)
		created, err := component.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, one.Reference)
		if err != nil {
			t.Fatal(err)
		}
		two, _ := component.SubmitApproval(context.Background(), submitCall(), harness.registrationID, second.EnvelopeBytes)
		collisionVerifier, _ := approvalattempt.NewFixtureVerifier([]approvalattempt.FixtureVector{first, second})
		collision, _ := NewApprovalAttempt(ApprovalAttemptOptions{
			Store: harness.store, Clock: harness.clock, Verifier: collisionVerifier,
			ApprovalIdentifiers: &approvalIDSequence{next: 99},
			AttemptIdentifiers:  attemptIDSequenceFunc(func(context.Context) (approvalattempt.AttemptID, error) { return created.Reference.AttemptID(), nil }),
			Integrity:           storeIntegrityAssessor{store: harness.store},
		})
		_, err = collision.RequestAttempt(context.Background(), attemptCall(), harness.registrationID, two.Reference)
		assertApprovalClassification(t, err, approvalattempt.ClassificationLocalFailure)
		record, _ := collision.Approval(context.Background(), two.Reference)
		if record.State != approvalattempt.ApprovalUsable {
			t.Fatal("attempt ID collision consumed approval")
		}
	})
}

type attemptIDSequenceFunc func(context.Context) (approvalattempt.AttemptID, error)

func (function attemptIDSequenceFunc) NewAttemptID(ctx context.Context) (approvalattempt.AttemptID, error) {
	return function(ctx)
}

type approvalIDSequenceFunc func(context.Context) (approvalattempt.ApprovalID, error)

func (function approvalIDSequenceFunc) NewApprovalID(ctx context.Context) (approvalattempt.ApprovalID, error) {
	return function(ctx)
}

func assertApprovalClassification(t *testing.T, err error, want approvalattempt.Classification) {
	t.Helper()
	classification, ok := approvalattempt.ErrorClassification(err)
	if !ok || classification != want {
		t.Fatalf("classification = %q (%t), want %q: %v", classification, ok, want, err)
	}
}

func jsonMarshalDiskState(state installationState) ([]byte, error) {
	return json.Marshal(diskEnvelope{StoreFormatVersion: 0, State: state})
}
