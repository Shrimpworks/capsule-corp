package registeredlifecycle

import (
	"bytes"
	"context"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"math/big"
	"path/filepath"
	"sync"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/execution/registrationstate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// This known answer was produced once by the explicitly authorized local
// fixture-generation task. Only the public P-256 coordinates and signed
// envelope are retained; no private key or signer is present in the repository.
const productionIntegrationPublicXHex = "d5bca91fa7a3f3af865093e3f6f9cf6fb9e21a528c961ce4e1b2742fbfd58660"
const productionIntegrationPublicYHex = "f8e52aece2aebe079e3cd20bc3a4cf43d8dd6ab05b3de3dda4c9e94b419024f8"
const productionIntegrationEnvelopeHex = "d2845854a3012603782b6170706c69636174696f6e2f63617073756c652e617070726f76616c2d6772616e742b63626f723b763d30045820854430da1c275d31ff9afa888be6ff420fcc5d955c68470d611d11a8aac5fd2fa058eaac017663617073756c652e617070726f76616c2d6772616e7402000350111111111111111111111111111111110458202222222222222222222222222222222222222222222222222222222222222222055033333333333333333333333333333333065820ef268a0b829adc1ce1307203f4b805f63379954ccf41e8e20a7487b6e5acf241075055555555555555555555555555555555085066666666666666666666666666666666097463617073756c652e706c616e2e617070726f76650a781c63617073756c652e657865637574696f6e2d73757065727669736f720b1a6a6be5800c1a6a6be6ac584028ab827a3338baef1edce5362f8b70b9523436122ef1a2f228978f7da32721851cab905555302628c6642500bf6d73e643ac1f1ecda59d952677f3ee8a6f30b8"

type productionRegistrationIDSequence struct {
	mu    sync.Mutex
	index int
	ids   []v0candidate.RegistrationID
}

func (source *productionRegistrationIDSequence) NewRegistrationID(context.Context) (v0candidate.RegistrationID, error) {
	source.mu.Lock()
	defer source.mu.Unlock()
	if source.index >= len(source.ids) {
		return v0candidate.RegistrationID{}, errors.New("production integration registration IDs exhausted")
	}
	result := source.ids[source.index]
	source.index++
	return result, nil
}

type productionIntegrationHarness struct {
	path           string
	clock          *testClock
	store          *registrationstate.FixedFileStore
	verifier       *approvalattempt.ProductionShapedVerifier
	attempts       *registrationstate.ApprovalAttemptComponent
	backend        *FakeBackend
	registrationA  v0candidate.PlanRegistration
	registrationB  v0candidate.PlanRegistration
	envelope       []byte
	approvalIDs    *approvalIDSequence
	attemptIDs     *attemptIDSequence
	effectIDs      *effectIDSequence
	checkpointSeen map[registrationstate.ApprovalAttemptCheckpoint]int
}

type failAfterLifecycleCommitStore struct {
	registrationstate.DurableLifecycleStore
	after func(approvalattempt.AttemptID) error
}

func (store failAfterLifecycleCommitStore) EnsureLifecycle(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	binding lifecyclestate.BackendBinding,
) (lifecyclestate.Record, bool, error) {
	record, created, err := store.DurableLifecycleStore.EnsureLifecycle(ctx, attemptID, binding)
	if err != nil || store.after == nil {
		return record, created, err
	}
	return record, created, store.after(attemptID)
}

func newProductionIntegrationHarness(t *testing.T, loseCommitResponses bool) *productionIntegrationHarness {
	t.Helper()
	path := filepath.Join(t.TempDir(), "supervisor-state.json")
	initial := registrationstate.InitialState{
		InstallationID:           repeated16[v0candidate.InstallationID](0x11),
		SupervisorID:             repeated16[v0candidate.SupervisorID](0x55),
		EpochSequence:            7,
		EpochDigest:              repeated32[v0candidate.TrustEpochDigest](0x22),
		TrustPhase:               registrationstate.TrustStable,
		TimeHighWaterUnixSeconds: 1_785_456_000,
	}
	store, err := registrationstate.NewFixedFileStore(path, initial)
	if err != nil {
		t.Fatal(err)
	}
	clock := &testClock{value: 1_785_456_000}
	registrationIDs := &productionRegistrationIDSequence{ids: []v0candidate.RegistrationID{
		repeated16[v0candidate.RegistrationID](0x33),
		repeated16[v0candidate.RegistrationID](0x34),
	}}
	registrations, err := registrationstate.New(registrationstate.Options{
		Store: store, Clock: clock, Identifiers: registrationIDs,
	})
	if err != nil {
		t.Fatal(err)
	}
	register := func(path string) v0candidate.PlanRegistration {
		issued, registerErr := registrations.RegisterPlan(
			context.Background(),
			registrationstate.AuthenticatedCallContext{
				Authenticated: true, Role: registrationstate.CallerDaemon,
				Purpose: registrationstate.RegisterPlanPurpose,
			},
			mustRead(t, path), ordinaryBindings(),
		)
		if registerErr != nil {
			t.Fatalf("register %s: %v", path, registerErr)
		}
		return issued.View()
	}
	registrationA := register(ordinaryPlanPath)
	registrationB := register("../../../schemas/conformance/v0/execution-plan/expires-one-second.cbor")
	if registrationA.PlanDigest != v0candidate.ExecutionPlanDigest(sha256.Sum256(mustRead(t, ordinaryPlanPath))) ||
		registrationA.PlanDigest == registrationB.PlanDigest {
		t.Fatal("production integration plans did not retain distinct exact-byte digests")
	}

	decode := func(value string) []byte {
		result, decodeErr := hex.DecodeString(value)
		if decodeErr != nil {
			t.Fatal(decodeErr)
		}
		return result
	}
	var publicKey approvalattempt.ApprovalP256PublicKey
	copy(publicKey.X[:], decode(productionIntegrationPublicXHex))
	copy(publicKey.Y[:], decode(productionIntegrationPublicYHex))
	authorization, err := approvalattempt.NewPassiveApprovalKeyAuthorization(
		publicKey, initial.InstallationID, initial.EpochSequence, initial.EpochDigest,
		1_785_455_900, 1_785_456_400,
	)
	if err != nil {
		t.Fatal(err)
	}
	verifier, err := approvalattempt.NewProductionShapedVerifier([]approvalattempt.ApprovalKeyAuthorization{authorization})
	if err != nil {
		t.Fatal(err)
	}
	bindings := ordinaryBindings()
	backendBinding, err := lifecyclestate.NewBackendBinding(lifecyclestate.BackendBindingView{
		Kind:            lifecyclestate.BackendFakeNoGuest,
		ProtocolVersion: lifecyclestate.FakeBackendProtocolVersion,
		ImplementationIdentityDigest: lifecyclestate.BackendImplementationDigest(
			sha256.Sum256([]byte("fake-no-guest-production-approval-integration")),
		),
		BackendConfigurationDigest:    bindings.BackendConfigurationDigest,
		BackendValidationRecordDigest: bindings.BackendValidationRecordDigest,
		CreatesGuest:                  false,
	})
	if err != nil {
		t.Fatal(err)
	}
	backend, err := NewFakeBackend(backendBinding)
	if err != nil || backend.CreatesGuest() {
		t.Fatalf("new no-guest backend: %v", err)
	}
	harness := &productionIntegrationHarness{
		path: path, clock: clock, store: store, verifier: verifier, backend: backend,
		registrationA: registrationA, registrationB: registrationB,
		envelope:    decode(productionIntegrationEnvelopeHex),
		approvalIDs: &approvalIDSequence{next: 1}, attemptIDs: &attemptIDSequence{next: 1},
		effectIDs: &effectIDSequence{next: 1}, checkpointSeen: make(map[registrationstate.ApprovalAttemptCheckpoint]int),
	}
	var checkpoint registrationstate.ApprovalAttemptCheckpointHook
	if loseCommitResponses {
		checkpoint = func(_ context.Context, point registrationstate.ApprovalAttemptCheckpoint) error {
			harness.checkpointSeen[point]++
			if (point == registrationstate.CheckpointAfterApprovalCommit ||
				point == registrationstate.CheckpointAfterAttemptCommit) && harness.checkpointSeen[point] == 1 {
				return errors.New("simulated committed response loss")
			}
			return nil
		}
	}
	harness.attempts = harness.newApprovalAttempt(t, store, harness.approvalIDs, harness.attemptIDs, checkpoint)
	return harness
}

func (harness *productionIntegrationHarness) newApprovalAttempt(
	t *testing.T,
	store registrationstate.StateStore,
	approvalIDs registrationstate.ApprovalIdentifierSource,
	attemptIDs registrationstate.AttemptIdentifierSource,
	checkpoint registrationstate.ApprovalAttemptCheckpointHook,
) *registrationstate.ApprovalAttemptComponent {
	t.Helper()
	component, err := registrationstate.NewApprovalAttempt(registrationstate.ApprovalAttemptOptions{
		Store: store, Clock: harness.clock, Verifier: harness.verifier,
		ApprovalIdentifiers: approvalIDs, AttemptIdentifiers: attemptIDs,
		Integrity: fixedIntegrity{assessedAt: 1_785_456_000}, Checkpoint: checkpoint,
	})
	if err != nil {
		t.Fatal(err)
	}
	return component
}

func productionSubmitCall() registrationstate.AuthenticatedCallContext {
	return registrationstate.AuthenticatedCallContext{
		Authenticated: true, Role: registrationstate.CallerBroker,
		Purpose: registrationstate.SubmitApprovalPurpose,
	}
}

func productionAttemptCall() registrationstate.AuthenticatedCallContext {
	return registrationstate.AuthenticatedCallContext{
		Authenticated: true, Role: registrationstate.CallerDaemon,
		Purpose: registrationstate.RequestAttemptPurpose,
	}
}

func complementaryProductionSignature(t *testing.T, envelope []byte) []byte {
	t.Helper()
	if len(envelope) < 64 {
		t.Fatal("production approval envelope has no raw signature")
	}
	result := bytes.Clone(envelope)
	signatureStart := len(result) - 64
	s := new(big.Int).SetBytes(result[signatureStart+32:])
	s.Sub(elliptic.P256().Params().N, s)
	copy(result[signatureStart+32:], s.FillBytes(make([]byte, 32)))
	return result
}

func productionSignatureForms(t *testing.T, envelope []byte) ([]byte, []byte) {
	t.Helper()
	complementary := complementaryProductionSignature(t, envelope)
	signatureStart := len(envelope) - 64
	originalS := new(big.Int).SetBytes(envelope[signatureStart+32:])
	complementaryS := new(big.Int).SetBytes(complementary[signatureStart+32:])
	if originalS.Cmp(complementaryS) < 0 {
		return bytes.Clone(envelope), complementary
	}
	return complementary, bytes.Clone(envelope)
}

func derProductionSignature(t *testing.T, envelope []byte) []byte {
	t.Helper()
	if len(envelope) < 66 || envelope[len(envelope)-66] != 0x58 || envelope[len(envelope)-65] != 0x40 {
		t.Fatal("production approval envelope does not end in a 64-byte signature")
	}
	signature := envelope[len(envelope)-64:]
	der, err := asn1.Marshal(struct {
		R *big.Int
		S *big.Int
	}{new(big.Int).SetBytes(signature[:32]), new(big.Int).SetBytes(signature[32:])})
	if err != nil || len(der) > 255 {
		t.Fatalf("encode DER signature: %v", err)
	}
	result := bytes.Clone(envelope[:len(envelope)-66])
	result = append(result, 0x58, byte(len(der)))
	return append(result, der...)
}

func assertApprovalIntegrationClassification(
	t *testing.T,
	err error,
	want approvalattempt.Classification,
) {
	t.Helper()
	classification, ok := approvalattempt.ErrorClassification(err)
	if !ok || classification != want {
		t.Fatalf("classification = %q (%t), want %q: %v", classification, ok, want, err)
	}
}

func TestProductionVerifiedApprovalCannotAuthorizeAnotherRegisteredPlan(t *testing.T) {
	harness := newProductionIntegrationHarness(t, false)
	lowS, highS := productionSignatureForms(t, harness.envelope)
	_, err := harness.attempts.SubmitApproval(
		context.Background(), productionSubmitCall(), harness.registrationB.RegistrationID, lowS,
	)
	assertApprovalIntegrationClassification(t, err, approvalattempt.ClassificationBinding)
	if snapshot := harness.backend.Snapshot(approvalattempt.AttemptID{}); len(snapshot.CallCounts) != 0 {
		t.Fatalf("binding refusal reached fake backend: %+v", snapshot)
	}

	submission, err := harness.attempts.SubmitApproval(
		context.Background(), productionSubmitCall(), harness.registrationA.RegistrationID, lowS,
	)
	if err != nil || submission.State != approvalattempt.ApprovalUsable {
		t.Fatalf("matching production approval = %+v, %v", submission, err)
	}
	equivalent, err := harness.attempts.SubmitApproval(
		context.Background(), productionSubmitCall(), harness.registrationA.RegistrationID,
		highS,
	)
	if err != nil || equivalent != submission {
		t.Fatalf("equivalent signature replay = %+v, %v; want %+v", equivalent, err, submission)
	}
}

func TestProductionIntegrationRefusesNonProfileSign1BeforeDurableState(t *testing.T) {
	tests := []struct {
		name     string
		mutate   func(*testing.T, []byte) []byte
		expected approvalattempt.Classification
	}{
		{"untagged", func(_ *testing.T, envelope []byte) []byte { return bytes.Clone(envelope[1:]) }, approvalattempt.ClassificationMalformed},
		{"der-signature", derProductionSignature, approvalattempt.ClassificationMalformed},
		{"non-equivalent-signature", func(_ *testing.T, envelope []byte) []byte {
			result := bytes.Clone(envelope)
			result[len(result)-1] ^= 0x01
			return result
		}, approvalattempt.ClassificationBinding},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			harness := newProductionIntegrationHarness(t, false)
			before := mustRead(t, harness.path)
			_, err := harness.attempts.SubmitApproval(
				context.Background(), productionSubmitCall(), harness.registrationA.RegistrationID,
				test.mutate(t, harness.envelope),
			)
			assertApprovalIntegrationClassification(t, err, test.expected)
			if after := mustRead(t, harness.path); !bytes.Equal(after, before) {
				t.Fatal("strict predecode or signature refusal changed durable authority state")
			}
			if attempts, listErr := harness.attempts.CreatedAttempts(context.Background()); listErr != nil || len(attempts) != 0 {
				t.Fatalf("refused envelope created attempts: %+v, %v", attempts, listErr)
			}
			if calls := harness.backend.Snapshot(approvalattempt.AttemptID{}).CallCounts; len(calls) != 0 {
				t.Fatalf("refused envelope reached fake backend: %v", calls)
			}
		})
	}
}

func TestProductionVerifiedApprovalCommitAndRestartConvergeBeforeFakeEffects(t *testing.T) {
	harness := newProductionIntegrationHarness(t, true)
	lowS, highS := productionSignatureForms(t, harness.envelope)
	_, err := harness.attempts.SubmitApproval(
		context.Background(), productionSubmitCall(), harness.registrationA.RegistrationID, lowS,
	)
	assertApprovalIntegrationClassification(t, err, approvalattempt.ClassificationLocalFailure)
	submission, err := harness.attempts.SubmitApproval(
		context.Background(), productionSubmitCall(), harness.registrationA.RegistrationID,
		highS,
	)
	if err != nil || submission.State != approvalattempt.ApprovalUsable {
		t.Fatalf("approval response-loss replay = %+v, %v", submission, err)
	}

	_, err = harness.attempts.RequestAttempt(
		context.Background(), productionAttemptCall(), harness.registrationA.RegistrationID, submission.Reference,
	)
	assertApprovalIntegrationClassification(t, err, approvalattempt.ClassificationLocalFailure)
	if calls := harness.backend.Snapshot(approvalattempt.AttemptID{}).CallCounts; len(calls) != 0 {
		t.Fatalf("attempt response loss reached fake backend: %v", calls)
	}

	reopenedV0, err := registrationstate.NewFixedFileStore(harness.path, registrationstate.InitialState{})
	if err != nil {
		t.Fatal(err)
	}
	restartedAttempts := harness.newApprovalAttempt(
		t, reopenedV0, &approvalIDSequence{next: 9}, &attemptIDSequence{next: 9}, nil,
	)
	created, err := restartedAttempts.RequestAttempt(
		context.Background(), productionAttemptCall(), harness.registrationA.RegistrationID, submission.Reference,
	)
	if err != nil || created.State != approvalattempt.AttemptCreated {
		t.Fatalf("attempt response-loss replay = %+v, %v", created, err)
	}
	if calls := harness.backend.Snapshot(created.Reference.AttemptID()).CallCounts; len(calls) != 0 {
		t.Fatalf("committed attempt reached fake backend before AttemptID-only drive: %v", calls)
	}
	consumedReplay, err := restartedAttempts.SubmitApproval(
		context.Background(), productionSubmitCall(), harness.registrationA.RegistrationID, lowS,
	)
	if err != nil || consumedReplay.Reference != submission.Reference ||
		consumedReplay.State != approvalattempt.ApprovalConsumed {
		t.Fatalf("consumed approval replay = %+v, %v", consumedReplay, err)
	}
	consumedEquivalent, err := restartedAttempts.SubmitApproval(
		context.Background(), productionSubmitCall(), harness.registrationA.RegistrationID, highS,
	)
	if err != nil || consumedEquivalent != consumedReplay {
		t.Fatalf("consumed complementary-S replay = %+v, %v; want %+v", consumedEquivalent, err, consumedReplay)
	}
	createdAgain, err := restartedAttempts.RequestAttempt(
		context.Background(), productionAttemptCall(), harness.registrationA.RegistrationID, submission.Reference,
	)
	if err != nil || createdAgain != created {
		t.Fatalf("created attempt replay = %+v, %v; want %+v", createdAgain, err, created)
	}
	if attempts, listErr := restartedAttempts.CreatedAttempts(context.Background()); listErr != nil ||
		len(attempts) != 1 || attempts[0] != created.Reference {
		t.Fatalf("created attempts = %+v, %v", attempts, listErr)
	}

	if _, err := registrationstate.MigrateFixedFileStoreV0ToV1(
		context.Background(), harness.path,
		registrationstate.V0ToV1MigrationOptions{Lock: offlineLock{}},
	); err != nil {
		t.Fatal(err)
	}
	storeV1, err := registrationstate.OpenFixedFileStoreV1WithOptions(
		harness.path,
		registrationstate.FixedFileStoreV1Options{
			EffectIDs: harness.effectIDs, OwnerSessionID: ownerID(t, 0x70),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	coordinator, err := NewCoordinator(storeV1.OwnerSessionID())
	if err != nil {
		t.Fatal(err)
	}
	commitBoundaryStore := failAfterLifecycleCommitStore{
		DurableLifecycleStore: storeV1,
		after: func(attemptID approvalattempt.AttemptID) error {
			if calls := harness.backend.Snapshot(attemptID).CallCounts; len(calls) != 0 {
				t.Fatalf("fake effect occurred before lifecycle-record commit returned: %v", calls)
			}
			return errors.New("simulated response loss after lifecycle-record commit")
		},
	}
	commitBoundaryLifecycle, err := New(Options{
		Attempts: restartedAttempts, Store: commitBoundaryStore, Backend: harness.backend,
		Coordinator: coordinator, Clock: harness.clock,
	})
	if err != nil {
		t.Fatal(err)
	}
	attemptID := created.Reference.AttemptID()
	if _, err := commitBoundaryLifecycle.Drive(context.Background(), attemptID); err == nil {
		t.Fatal("expected simulated response loss after lifecycle-record commit")
	}
	if calls := harness.backend.Snapshot(attemptID).CallCounts; len(calls) != 0 {
		t.Fatalf("lifecycle-record response loss reached fake backend: %v", calls)
	}
	if _, err := storeV1.ReadLifecycle(context.Background(), attemptID); err != nil {
		t.Fatalf("lifecycle record was not durable before the zero-effect boundary: %v", err)
	}

	reopenedAfterLifecycleCommit, err := registrationstate.OpenFixedFileStoreV1WithOptions(
		harness.path,
		registrationstate.FixedFileStoreV1Options{
			EffectIDs: harness.effectIDs, OwnerSessionID: ownerID(t, 0x71),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	prepareCoordinator, err := NewCoordinator(reopenedAfterLifecycleCommit.OwnerSessionID())
	if err != nil {
		t.Fatal(err)
	}
	firstEffectLifecycle, err := New(Options{
		Attempts: restartedAttempts, Store: reopenedAfterLifecycleCommit, Backend: harness.backend,
		Coordinator: prepareCoordinator, Clock: harness.clock,
		Checkpoint: oneShotCheckpoint(CheckpointAfterPrepareEffect),
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := firstEffectLifecycle.Drive(context.Background(), attemptID); err == nil {
		t.Fatal("expected simulated response loss after committed attempt and fake prepare")
	}
	prepared := harness.backend.Snapshot(attemptID)
	if prepared.ApplicationCounts[OperationPrepare] != 1 {
		t.Fatalf("prepare applications after response loss = %d", prepared.ApplicationCounts[OperationPrepare])
	}

	reopenedV1, err := registrationstate.OpenFixedFileStoreV1WithOptions(
		harness.path,
		registrationstate.FixedFileStoreV1Options{
			EffectIDs: harness.effectIDs, OwnerSessionID: ownerID(t, 0x72),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	reopenedCoordinator, err := NewCoordinator(reopenedV1.OwnerSessionID())
	if err != nil {
		t.Fatal(err)
	}
	restartedLifecycle, err := New(Options{
		Attempts: restartedAttempts, Store: reopenedV1, Backend: harness.backend,
		Coordinator: reopenedCoordinator, Clock: harness.clock,
	})
	if err != nil {
		t.Fatal(err)
	}
	recovered, err := restartedLifecycle.Recover(context.Background(), attemptID)
	if err != nil || recovered.State != StateDestroyed || recovered.CleanupRequired ||
		recovered.AttemptID != attemptID || harness.backend.CreatesGuest() {
		t.Fatalf("production approval lifecycle recovery = %+v, %v", recovered, err)
	}
	backend := harness.backend.Snapshot(attemptID)
	for _, operation := range durableOperations {
		if backend.ApplicationCounts[operation] != 1 {
			t.Fatalf("%s applications = %d, want one", operation, backend.ApplicationCounts[operation])
		}
	}
	if replayed, replayErr := restartedLifecycle.Drive(context.Background(), attemptID); replayErr != nil ||
		replayed.State != StateDestroyed {
		t.Fatalf("terminal AttemptID replay = %+v, %v", replayed, replayErr)
	}
	for _, operation := range durableOperations {
		if harness.backend.Snapshot(attemptID).ApplicationCounts[operation] != 1 {
			t.Fatalf("terminal replay redrove %s", operation)
		}
	}
}
