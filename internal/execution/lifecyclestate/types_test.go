package lifecyclestate

import (
	"bytes"
	"encoding/hex"
	"testing"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

func TestEffectIDIsNonzeroCopiedAndDomainSeparated(t *testing.T) {
	source := bytes.Repeat([]byte{0x41}, 16)
	effectToken, err := NewDomainIdentifier(DomainEffectID, source)
	if err != nil {
		t.Fatalf("new effect token: %v", err)
	}
	source[0] ^= 0xff
	effectID, err := NewEffectID(effectToken)
	if err != nil {
		t.Fatalf("new effect ID: %v", err)
	}
	if effectID[0] != 0x41 {
		t.Fatal("effect ID retained caller-owned bytes")
	}
	ownerToken, err := NewDomainIdentifier(DomainOwnerSessionID, bytes.Repeat([]byte{0x42}, 16))
	if err != nil {
		t.Fatalf("new owner token: %v", err)
	}
	if _, err := NewEffectID(ownerToken); classification(t, err) != ClassificationDomain {
		t.Fatalf("owner token entered effect domain: %v", err)
	}
	if _, err := NewOwnerSessionID(effectToken); classification(t, err) != ClassificationDomain {
		t.Fatalf("effect token entered owner domain: %v", err)
	}
	if _, err := NewDomainIdentifier(DomainEffectID, make([]byte, 16)); classification(t, err) != ClassificationSchema {
		t.Fatalf("zero effect token accepted: %v", err)
	}
	if _, err := NewDomainIdentifier("attempt", bytes.Repeat([]byte{1}, 16)); classification(t, err) != ClassificationDomain {
		t.Fatalf("foreign identifier domain accepted: %v", err)
	}

	projected := effectToken.Bytes()
	projected[0] ^= 0xff
	if effectToken.Bytes()[0] != 0x41 {
		t.Fatal("domain identifier exposed mutable bytes")
	}
}

func TestClosedEnumProjectionsAreDefensive(t *testing.T) {
	operations := Operations()
	operations[0] = "mutated"
	if Operations()[0] != OperationPrepare {
		t.Fatal("operation projection exposed package state")
	}
	states := LifecycleStates()
	states[0] = "mutated"
	if LifecycleStates()[0] != StatePreparePending {
		t.Fatal("state projection exposed package state")
	}
	checkpoints := Checkpoints()
	checkpoints[0] = "mutated"
	if Checkpoints()[0] != CheckpointNone {
		t.Fatal("checkpoint projection exposed package state")
	}
	failures := FailureClassifications()
	failures[0] = "mutated"
	if FailureClassifications()[0] != FailureNone {
		t.Fatal("failure projection exposed package state")
	}
	reconciliation := ReconciliationStatuses()
	reconciliation[0] = "mutated"
	if ReconciliationStatuses()[0] != ReconciliationNone {
		t.Fatal("reconciliation projection exposed package state")
	}
	backends := BackendKinds()
	backends[0] = "mutated"
	if BackendKinds()[0] != BackendFakeNoGuest {
		t.Fatal("backend projection exposed package state")
	}
}

func TestImmutableBindingsCopyCrossLinkAndKnownAnswer(t *testing.T) {
	view := validBindingsView(t)
	binding, err := NewImmutableBindings(view)
	if err != nil {
		t.Fatalf("new immutable bindings: %v", err)
	}

	const wantDigest = "b0bdbe2a1d5862181789aa18e0f019e8d1ca30e76c6dd94dcbb12e3bbc5f9f37"
	bindingDigest := binding.Digest()
	if got := hex.EncodeToString(bindingDigest[:]); got != wantDigest {
		t.Fatalf("immutable binding digest = %s, want %s", got, wantDigest)
	}

	view.ProfileReviewAttestationDigests[0][0] ^= 0xff
	if binding.View().ProfileReviewAttestationDigests[0][0] != 0x31 {
		t.Fatal("immutable bindings retained caller-owned collection")
	}
	projected := binding.View()
	projected.ProfileReviewAttestationDigests[0][0] ^= 0xff
	if binding.View().ProfileReviewAttestationDigests[0][0] != 0x31 {
		t.Fatal("immutable binding view exposed mutable collection")
	}

	restored, err := RestoreImmutableBindings(binding.View(), binding.Digest())
	if err != nil || !restored.Equal(binding) {
		t.Fatalf("restore immutable bindings: %v", err)
	}
	badDigest := binding.Digest()
	badDigest[0] ^= 0xff
	if _, err := RestoreImmutableBindings(binding.View(), badDigest); classification(t, err) != ClassificationBinding {
		t.Fatalf("wrong immutable digest accepted: %v", err)
	}

	reordered := binding.View()
	reordered.ProfileReviewAttestationDigests[0], reordered.ProfileReviewAttestationDigests[1] =
		reordered.ProfileReviewAttestationDigests[1], reordered.ProfileReviewAttestationDigests[0]
	reorderedBinding, err := NewImmutableBindings(reordered)
	if err != nil {
		t.Fatalf("new reordered binding: %v", err)
	}
	if reorderedBinding.Digest() == binding.Digest() {
		t.Fatal("ordered review digest projection was order-insensitive")
	}

	wrongCrossLink := binding.View()
	backend := wrongCrossLink.Backend.View()
	backend.BackendConfigurationDigest[0] ^= 0xff
	wrongCrossLink.Backend, err = NewBackendBinding(backend)
	if err != nil {
		t.Fatalf("new mismatched backend binding: %v", err)
	}
	if _, err := NewImmutableBindings(wrongCrossLink); classification(t, err) != ClassificationBinding {
		t.Fatalf("backend/plan cross-link mismatch accepted: %v", err)
	}

	overCapacity := binding.View()
	overCapacity.ProfileReviewAttestationDigests = make(
		[]v0candidate.ProfileReviewAttestationDigest,
		MaxProfileReviewAttestationDigests+1,
	)
	for index := range overCapacity.ProfileReviewAttestationDigests {
		overCapacity.ProfileReviewAttestationDigests[index] = repeated32[v0candidate.ProfileReviewAttestationDigest](byte(index + 1))
	}
	if _, err := NewImmutableBindings(overCapacity); classification(t, err) != ClassificationCapacity {
		t.Fatalf("over-capacity review collection accepted: %v", err)
	}
}

func TestBackendBindingAndInstanceFailClosed(t *testing.T) {
	backend := validBackendBinding(t)
	view := backend.View()
	view.Kind = "real-backend"
	if _, err := NewBackendBinding(view); classification(t, err) != ClassificationUnsupported {
		t.Fatalf("unknown backend accepted: %v", err)
	}
	view = backend.View()
	view.CreatesGuest = true
	if _, err := NewBackendBinding(view); classification(t, err) != ClassificationUnsupported {
		t.Fatalf("guest-creating backend accepted: %v", err)
	}

	source := []byte{0x01, 0x02, 0x03, 0x04}
	identity, err := NewBackendInstanceIdentity(BackendInstanceFake, source)
	if err != nil {
		t.Fatalf("new instance identity: %v", err)
	}
	source[0] ^= 0xff
	if identity.Value()[0] != 0x01 {
		t.Fatal("instance identity retained caller-owned bytes")
	}
	const wantDigest = "9257a84da53e87753c4e6781d0967fcc56b4c7d07e6b756dbe15b1ca746b79c4"
	instanceDigest := identity.Digest()
	if got := hex.EncodeToString(instanceDigest[:]); got != wantDigest {
		t.Fatalf("instance digest = %s, want %s", got, wantDigest)
	}
	projected := identity.View()
	projected.Value[0] ^= 0xff
	if identity.Value()[0] != 0x01 {
		t.Fatal("instance identity exposed mutable bytes")
	}
	projected = identity.View()
	projected.Digest[0] ^= 0xff
	if _, err := RestoreBackendInstanceIdentity(projected); classification(t, err) != ClassificationBinding {
		t.Fatalf("wrong instance digest accepted: %v", err)
	}
	if _, err := NewBackendInstanceIdentity(BackendInstanceFake, make([]byte, MaxBackendInstanceIdentityBytes+1)); classification(t, err) != ClassificationSchema {
		t.Fatalf("oversized instance identity accepted: %v", err)
	}
	if _, err := NewBackendInstanceIdentity("pid", []byte{1}); classification(t, err) != ClassificationUnsupported {
		t.Fatalf("unknown instance kind accepted: %v", err)
	}
}

func TestRecordLifecycleShapesAndUInt53Times(t *testing.T) {
	for _, state := range LifecycleStates() {
		view := validRecordView(t, state)
		record, err := NewRecord(view, 200)
		if err != nil {
			t.Fatalf("state %s: %v", state, err)
		}
		if err := record.Validate(200); err != nil {
			t.Fatalf("revalidate state %s: %v", state, err)
		}
	}

	pending := validRecordView(t, StatePreparePending)
	pending.CleanupRequired = false
	if _, err := NewRecord(pending, 200); classification(t, err) != ClassificationBinding {
		t.Fatalf("pending record without cleanup accepted: %v", err)
	}
	pending = validRecordView(t, StatePreparePending)
	pending.ImmutableBindingDigest[0] ^= 0xff
	if _, err := NewRecord(pending, 200); classification(t, err) != ClassificationBinding {
		t.Fatalf("wrong record binding digest accepted: %v", err)
	}
	pending = validRecordView(t, StatePreparePending)
	pending.OpenedAt = v0candidate.UInt53(v0candidate.MaxSafeInteger + 1)
	if _, err := NewRecord(pending, v0candidate.UInt53(v0candidate.MaxSafeInteger)); classification(t, err) != ClassificationSchema {
		t.Fatalf("non-UInt53 record time accepted: %v", err)
	}
	pending = validRecordView(t, StatePreparePending)
	if _, err := NewRecord(pending, 99); classification(t, err) != ClassificationStale {
		t.Fatalf("time above snapshot high-water accepted: %v", err)
	}
	pending = validRecordView(t, StatePreparePending)
	pending.State = "future-state"
	if _, err := NewRecord(pending, 200); classification(t, err) != ClassificationUnsupported {
		t.Fatalf("unknown lifecycle state accepted: %v", err)
	}

	for name, mutate := range map[string]func(*RecordView){
		"effect status":  func(view *RecordView) { view.EffectStatus = "future-effect" },
		"checkpoint":     func(view *RecordView) { view.LastConfirmedCheckpoint = "future-checkpoint" },
		"failure":        func(view *RecordView) { view.FirstFailure = "future-failure" },
		"reconciliation": func(view *RecordView) { view.LastReconciliation = "future-reconciliation" },
		"recovery fence": func(view *RecordView) { view.RecoveryFence = "future-fence" },
	} {
		t.Run(name, func(t *testing.T) {
			candidate := validRecordView(t, StatePreparePending)
			mutate(&candidate)
			if _, err := NewRecord(candidate, 200); classification(t, err) != ClassificationUnsupported {
				t.Fatalf("unknown enum accepted: %v", err)
			}
		})
	}

	maxBindingsView := validBindingsView(t)
	maxBindingsView.AttemptCreatedAt = v0candidate.UInt53(v0candidate.MaxSafeInteger - 1)
	maxBindings, err := NewImmutableBindings(maxBindingsView)
	if err != nil {
		t.Fatalf("max UInt53 bindings: %v", err)
	}
	maxRecord := validRecordView(t, StatePreparePending)
	maxRecord.Bindings = maxBindings
	maxRecord.ImmutableBindingDigest = maxBindings.Digest()
	maxRecord.OpenedAt = v0candidate.UInt53(v0candidate.MaxSafeInteger - 1)
	maxRecord.LastTransitionAt = v0candidate.UInt53(v0candidate.MaxSafeInteger)
	if _, err := NewRecord(maxRecord, v0candidate.UInt53(v0candidate.MaxSafeInteger)); err != nil {
		t.Fatalf("maximum UInt53 times rejected: %v", err)
	}

	exhausted := validRecordView(t, StateUnresolved)
	exhausted.AutomaticRecoveryCount = MaxAutomaticRecoveryAttempts
	exhausted.AutomaticRecoveryExhausted = true
	exhausted.NextRecoveryAt = OptionalUnixSeconds{}
	exhausted.RecoveryFence = RecoveryFenceAutomaticExhausted
	if _, err := NewRecord(exhausted, 200); err != nil {
		t.Fatalf("valid exhausted recovery record: %v", err)
	}
	exhausted.AutomaticRecoveryCount--
	if _, err := NewRecord(exhausted, 200); classification(t, err) != ClassificationBinding {
		t.Fatalf("inconsistent exhaustion accepted: %v", err)
	}

	destroyedWithoutAbsence := validRecordView(t, StateDestroyed)
	destroyedWithoutAbsence.LastReconciliation = ReconciliationNone
	destroyedWithoutAbsence.AutomaticRecoveryCount = 0
	if _, err := NewRecord(destroyedWithoutAbsence, 200); classification(t, err) != ClassificationBinding {
		t.Fatalf("destroyed record without authoritative absence accepted: %v", err)
	}

	earlyAbsence := validRecordView(t, StatePreparePending)
	earlyAbsence.State = StateDestroyed
	earlyAbsence.CleanupRequired = false
	earlyAbsence.LastReconciliation = ReconciliationAuthoritativelyAbsent
	earlyAbsence.AutomaticRecoveryCount = 1
	earlyAbsence.TerminalAt = OptionalUnixSeconds{Present: true, Value: earlyAbsence.LastTransitionAt}
	if _, err := NewRecord(earlyAbsence, 200); err != nil {
		t.Fatalf("authoritative absence before create rejected: %v", err)
	}
}

func TestRecordViewsAndPermitsDefensivelyCopyInstance(t *testing.T) {
	view := validRecordView(t, StateCreated)
	record, err := NewRecord(view, 200)
	if err != nil {
		t.Fatalf("new record: %v", err)
	}
	projected := record.View()
	instanceBytes := projected.Instance.Value()
	instanceBytes[0] ^= 0xff
	reviews := projected.Bindings.View().ProfileReviewAttestationDigests
	reviews[0][0] ^= 0xff
	if record.View().Instance.Value()[0] != 0x91 ||
		record.Bindings().View().ProfileReviewAttestationDigests[0][0] != 0x31 {
		t.Fatal("record projection exposed mutable nested state")
	}

	permit, err := NewEffectPermit(EffectPermitView{
		AttemptID:              record.Bindings().View().AttemptID,
		RecordVersion:          2,
		OperationSequence:      3,
		Operation:              OperationStart,
		EffectID:               effectID(t, 0xa1),
		ImmutableBindingDigest: record.ImmutableBindingDigest(),
		OwnerSessionID:         ownerSessionID(t, 0xb1),
		Instance:               record.View().Instance,
	})
	if err != nil {
		t.Fatalf("new effect permit: %v", err)
	}
	permitBytes := permit.View().Instance.Value()
	permitBytes[0] ^= 0xff
	if permit.View().Instance.Value()[0] != 0x91 {
		t.Fatal("effect permit exposed mutable instance bytes")
	}
	badPermit := permit.View()
	badPermit.Operation = OperationCreate
	if _, err := NewEffectPermit(badPermit); classification(t, err) != ClassificationBinding {
		t.Fatalf("create permit with instance accepted: %v", err)
	}
}

func TestEffectAndReconciliationResultsRequireExactIdentityShape(t *testing.T) {
	identity := instanceIdentity(t)
	if _, err := NewEffectResult(OperationCreate, EffectResultApplied, identity); err != nil {
		t.Fatalf("create applied result: %v", err)
	}
	if _, err := NewEffectResult(OperationCreate, EffectResultApplied, BackendInstanceIdentity{}); classification(t, err) != ClassificationBinding {
		t.Fatalf("create applied without identity accepted: %v", err)
	}
	if _, err := NewEffectResult(OperationStart, EffectResultApplied, identity); classification(t, err) != ClassificationBinding {
		t.Fatalf("non-create result introduced identity: %v", err)
	}
	if _, err := NewEffectResult(OperationStart, "success", BackendInstanceIdentity{}); classification(t, err) != ClassificationUnsupported {
		t.Fatalf("unknown effect result accepted: %v", err)
	}

	if _, err := NewReconcileResult(OperationCreate, ReconciliationApplied, identity); err != nil {
		t.Fatalf("create reconciliation applied: %v", err)
	}
	if _, err := NewReconcileResult(OperationDestroy, ReconciliationAuthoritativelyAbsent, BackendInstanceIdentity{}); err != nil {
		t.Fatalf("destroy authoritative absence: %v", err)
	}
	if _, err := NewReconcileResult(OperationStart, ReconciliationPresent, BackendInstanceIdentity{}); classification(t, err) != ClassificationBinding {
		t.Fatalf("present reconciliation without identity accepted: %v", err)
	}
	if _, err := NewReconcileResult(OperationStart, ReconciliationNone, BackendInstanceIdentity{}); classification(t, err) != ClassificationUnsupported {
		t.Fatalf("none reconciliation returned as observation: %v", err)
	}
}

func validBindingsView(t *testing.T) ImmutableBindingsView {
	t.Helper()
	backend := validBackendBinding(t)
	return ImmutableBindingsView{
		AttemptID:                   repeated16[approvalattempt.AttemptID](0x11),
		ApprovalID:                  repeated16[approvalattempt.ApprovalID](0x12),
		AttemptNonce:                repeated16[approvalattempt.AttemptNonce](0x13),
		RegistrationID:              repeated16[v0candidate.RegistrationID](0x14),
		RegistrationSequence:        7,
		PlanDigest:                  repeated32[v0candidate.ExecutionPlanDigest](0x15),
		InstallationID:              repeated16[v0candidate.InstallationID](0x16),
		EpochSequence:               8,
		EpochDigest:                 repeated32[v0candidate.TrustEpochDigest](0x17),
		SupervisorID:                repeated16[v0candidate.SupervisorID](0x18),
		ApprovalPurpose:             approvalattempt.ApprovalGrantPurpose,
		ApprovalAudience:            approvalattempt.ApprovalGrantAudience,
		ApprovalPayloadDigest:       repeated32[approvalattempt.ApprovalPayloadDigest](0x19),
		AuthorizationIdentity:       repeated32[approvalattempt.ApprovalKeyAuthorizationIdentity](0x20),
		AttemptCreatedAt:            90,
		AttemptStorageVersion:       SliceBAttemptStorageVersion,
		ApprovalStorageVersion:      SliceBApprovalStorageVersion,
		RegistrationStorageVersion:  SliceBRegistrationStorageVersion,
		SourceManifestDigest:        repeated32[v0candidate.SourceManifestDigest](0x21),
		InlineInputDigest:           repeated32[v0candidate.InlineInputDigest](0x22),
		RuntimeBundleManifestDigest: repeated32[v0candidate.RuntimeBundleManifestDigest](0x23),
		ProfileReviewAttestationDigests: []v0candidate.ProfileReviewAttestationDigest{
			repeated32[v0candidate.ProfileReviewAttestationDigest](0x31),
			repeated32[v0candidate.ProfileReviewAttestationDigest](0x32),
		},
		ProfileRegistryEntryDigest:    repeated32[v0candidate.ProfileRegistryEntryDigest](0x24),
		BackendValidationRecordDigest: repeated32[v0candidate.BackendValidationRecordDigest](0x25),
		BackendConfigurationDigest:    repeated32[v0candidate.BackendConfigurationDigest](0x26),
		TrustSnapshotDigest:           repeated32[v0candidate.TrustSnapshotDigest](0x27),
		PolicyDecisionDigest:          repeated32[v0candidate.PolicyDecisionDigest](0x28),
		Backend:                       backend,
	}
}

func validBackendBinding(t *testing.T) BackendBinding {
	t.Helper()
	binding, err := NewBackendBinding(BackendBindingView{
		Kind:                          BackendFakeNoGuest,
		ProtocolVersion:               FakeBackendProtocolVersion,
		ImplementationIdentityDigest:  repeated32[BackendImplementationDigest](0x29),
		BackendConfigurationDigest:    repeated32[v0candidate.BackendConfigurationDigest](0x26),
		BackendValidationRecordDigest: repeated32[v0candidate.BackendValidationRecordDigest](0x25),
		CreatesGuest:                  false,
	})
	if err != nil {
		t.Fatalf("new backend binding: %v", err)
	}
	return binding
}

func validRecordView(t *testing.T, state LifecycleState) RecordView {
	t.Helper()
	bindings, err := NewImmutableBindings(validBindingsView(t))
	if err != nil {
		t.Fatalf("new bindings: %v", err)
	}
	view := RecordView{
		FormatVersion:           LifecycleRecordFormatVersion,
		RecordVersion:           1,
		SnapshotGeneration:      1,
		Bindings:                bindings,
		ImmutableBindingDigest:  bindings.Digest(),
		State:                   state,
		CleanupRequired:         true,
		LastConfirmedCheckpoint: CheckpointNone,
		FirstFailure:            FailureNone,
		FailureOperation:        OperationNone,
		LastReconciliation:      ReconciliationNone,
		RecoveryFence:           RecoveryFenceNone,
		OpenedAt:                100,
		LastTransitionAt:        101,
	}
	setEffect := func(operation Operation, status EffectStatus, checkpoint Checkpoint, withInstance bool) {
		view.Operation = operation
		view.EffectStatus = status
		view.LastConfirmedCheckpoint = checkpoint
		if status != EffectNone {
			view.OperationSequence = 1
			view.EffectID = effectID(t, 0x81)
		}
		if withInstance {
			view.Instance = instanceIdentity(t)
		}
	}
	switch state {
	case StatePreparePending:
		setEffect(OperationNone, EffectNone, CheckpointNone, false)
	case StatePrepareIntent:
		setEffect(OperationPrepare, EffectIntent, CheckpointNone, false)
	case StatePrepared:
		setEffect(OperationPrepare, EffectConfirmed, CheckpointPrepare, false)
	case StateCreateIntent:
		setEffect(OperationCreate, EffectIntent, CheckpointPrepare, false)
	case StateCreated:
		setEffect(OperationCreate, EffectConfirmed, CheckpointCreate, true)
	case StateStartIntent:
		setEffect(OperationStart, EffectIntent, CheckpointCreate, true)
	case StateStarted:
		setEffect(OperationStart, EffectConfirmed, CheckpointStart, true)
	case StateObserveIntent:
		setEffect(OperationObserve, EffectIntent, CheckpointStart, true)
	case StateObserved:
		setEffect(OperationObserve, EffectConfirmed, CheckpointObserve, true)
	case StateStopIntent:
		setEffect(OperationStop, EffectIntent, CheckpointObserve, true)
	case StateStopped:
		setEffect(OperationStop, EffectConfirmed, CheckpointStop, true)
	case StateDestroyIntent:
		setEffect(OperationDestroy, EffectIntent, CheckpointStop, true)
	case StateDestroyConfirmed:
		setEffect(OperationDestroy, EffectConfirmed, CheckpointDestroy, true)
	case StateDestroyed:
		setEffect(OperationDestroy, EffectConfirmed, CheckpointDestroy, true)
		view.CleanupRequired = false
		view.LastReconciliation = ReconciliationAuthoritativelyAbsent
		view.AutomaticRecoveryCount = 1
		view.TerminalAt = OptionalUnixSeconds{Present: true, Value: view.LastTransitionAt}
	case StateUnresolved:
		setEffect(OperationObserve, EffectIndeterminate, CheckpointStart, true)
		view.FirstFailure = FailureCleanupUnresolved
		view.FailureOperation = OperationObserve
		view.LastReconciliation = ReconciliationUnknown
		view.AutomaticRecoveryCount = 1
		view.NextRecoveryAt = OptionalUnixSeconds{Present: true, Value: 105}
		view.RecoveryFence = RecoveryFenceReconcileUnknown
	case StateQuarantined:
		setEffect(OperationStart, EffectIndeterminate, CheckpointCreate, true)
		view.FirstFailure = FailureBinding
		view.FailureOperation = OperationStart
		view.LastReconciliation = ReconciliationIdentityMismatch
		view.AutomaticRecoveryCount = 1
		view.RecoveryFence = RecoveryFenceIdentityMismatch
	}
	return view
}

func instanceIdentity(t *testing.T) BackendInstanceIdentity {
	t.Helper()
	identity, err := NewBackendInstanceIdentity(BackendInstanceFake, []byte{0x91, 0x92, 0x93})
	if err != nil {
		t.Fatalf("new instance identity: %v", err)
	}
	return identity
}

func effectID(t *testing.T, value byte) EffectID {
	t.Helper()
	token, err := NewDomainIdentifier(DomainEffectID, bytes.Repeat([]byte{value}, 16))
	if err != nil {
		t.Fatalf("new effect token: %v", err)
	}
	identifier, err := NewEffectID(token)
	if err != nil {
		t.Fatalf("new effect ID: %v", err)
	}
	return identifier
}

func ownerSessionID(t *testing.T, value byte) OwnerSessionID {
	t.Helper()
	token, err := NewDomainIdentifier(DomainOwnerSessionID, bytes.Repeat([]byte{value}, 16))
	if err != nil {
		t.Fatalf("new owner token: %v", err)
	}
	identifier, err := NewOwnerSessionID(token)
	if err != nil {
		t.Fatalf("new owner ID: %v", err)
	}
	return identifier
}

func classification(t *testing.T, err error) Classification {
	t.Helper()
	if err == nil {
		return ""
	}
	classification, ok := ErrorClassification(err)
	if !ok {
		t.Fatalf("unclassified error: %v", err)
	}
	return classification
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
