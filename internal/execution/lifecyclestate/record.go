package lifecyclestate

import "capsule.local/capsule/internal/protocol/v0candidate"

const (
	// LifecycleRecordFormatVersion is the only passive record encoding version.
	LifecycleRecordFormatVersion = uint64(0)
	// MaxAutomaticRecoveryAttempts bounds automatic observation-only recovery.
	MaxAutomaticRecoveryAttempts = v0candidate.UInt53(3)
)

// RecordVersion is the positive monotonic version of one retained lifecycle record.
type RecordVersion uint64

// OptionalUnixSeconds represents an explicit optional UInt53 timestamp. A
// zero Unix second remains distinguishable from absence.
type OptionalUnixSeconds struct {
	Present bool
	Value   v0candidate.UInt53
}

// RecordView is the complete passive lifecycle record projection. NewRecord
// validates it and retains only private, immutable copies.
type RecordView struct {
	FormatVersion      uint64
	RecordVersion      RecordVersion
	SnapshotGeneration v0candidate.PositiveUInt53

	Bindings               ImmutableBindings
	ImmutableBindingDigest ImmutableBindingDigest

	OperationSequence v0candidate.PositiveUInt53
	Operation         Operation
	EffectID          EffectID
	EffectStatus      EffectStatus

	Instance BackendInstanceIdentity

	State                   LifecycleState
	CleanupRequired         bool
	LastConfirmedCheckpoint Checkpoint
	FirstFailure            FailureClassification
	FailureOperation        Operation

	LastReconciliation         ReconciliationStatus
	AutomaticRecoveryCount     v0candidate.UInt53
	NextRecoveryAt             OptionalUnixSeconds
	RecoveryFence              RecoveryFenceReason
	AutomaticRecoveryExhausted bool

	OpenedAt         v0candidate.UInt53
	LastTransitionAt v0candidate.UInt53
	TerminalAt       OptionalUnixSeconds
}

// Record is an immutable validated record snapshot. Its collection and
// transaction semantics deliberately belong to later slices.
type Record struct {
	view RecordView
}

// NewRecord validates a complete immutable snapshot against the durable time
// high-water from the same future store transaction and retains defensive data.
func NewRecord(view RecordView, durableTimeHighWater v0candidate.UInt53) (Record, error) {
	if err := validateRecordView(view, durableTimeHighWater); err != nil {
		return Record{}, err
	}
	return Record{view: cloneRecordView(view)}, nil
}

// View returns a defensive projection. Nested bindings and instance bytes
// remain immutable and expose only copied slices through their accessors.
func (record Record) View() RecordView { return cloneRecordView(record.view) }

// Bindings returns the record's immutable validated authority/runtime binding.
func (record Record) Bindings() ImmutableBindings { return record.view.Bindings }

// ImmutableBindingDigest returns the digest cross-linked to Bindings.
func (record Record) ImmutableBindingDigest() ImmutableBindingDigest {
	return record.view.ImmutableBindingDigest
}

// Validate repeats fail-closed validation against the high-water from the
// same future committed snapshot.
func (record Record) Validate(durableTimeHighWater v0candidate.UInt53) error {
	return validateRecordView(record.view, durableTimeHighWater)
}

func cloneRecordView(view RecordView) RecordView {
	// ImmutableBindings and BackendInstanceIdentity expose no mutable slice.
	// Reconstructing them also proves that any future restore path preserves
	// their deterministic digest rather than trusting a shallow copy.
	bindings, err := RestoreImmutableBindings(view.Bindings.View(), view.Bindings.Digest())
	if err == nil {
		view.Bindings = bindings
	}
	if view.Instance.Present() {
		instance, instanceErr := RestoreBackendInstanceIdentity(view.Instance.View())
		if instanceErr == nil {
			view.Instance = instance
		}
	}
	return view
}

func validateRecordView(view RecordView, durableTimeHighWater v0candidate.UInt53) error {
	if view.FormatVersion != LifecycleRecordFormatVersion {
		return classified(ClassificationUnsupported, "record-format-version")
	}
	if view.RecordVersion == 0 || uint64(view.RecordVersion) > v0candidate.MaxSafeInteger ||
		view.SnapshotGeneration == 0 ||
		uint64(view.SnapshotGeneration) > v0candidate.MaxSafeInteger {
		return classified(ClassificationSchema, "record-version-generation")
	}
	if !view.Bindings.valid() || view.ImmutableBindingDigest != view.Bindings.Digest() {
		return classified(ClassificationBinding, "record-immutable-bindings")
	}
	if !validLifecycleState(view.State) || !validEffectStatus(view.EffectStatus) ||
		!validCheckpoint(view.LastConfirmedCheckpoint) ||
		!validFailureClassification(view.FirstFailure) ||
		!validReconciliationStatus(view.LastReconciliation) ||
		!validRecoveryFenceReason(view.RecoveryFence) {
		return classified(ClassificationUnsupported, "record-enum")
	}
	if !view.Instance.valid() {
		return classified(ClassificationBinding, "record-instance")
	}
	if uint64(durableTimeHighWater) > v0candidate.MaxSafeInteger ||
		uint64(view.OpenedAt) > v0candidate.MaxSafeInteger ||
		uint64(view.LastTransitionAt) > v0candidate.MaxSafeInteger ||
		(view.TerminalAt.Present && uint64(view.TerminalAt.Value) > v0candidate.MaxSafeInteger) ||
		(view.NextRecoveryAt.Present && uint64(view.NextRecoveryAt.Value) > v0candidate.MaxSafeInteger) {
		return classified(ClassificationSchema, "record-time-uint53")
	}
	createdAt := view.Bindings.View().AttemptCreatedAt
	if view.OpenedAt < createdAt || view.LastTransitionAt < view.OpenedAt ||
		view.OpenedAt > durableTimeHighWater ||
		view.LastTransitionAt > durableTimeHighWater ||
		(view.TerminalAt.Present && (view.TerminalAt.Value < view.LastTransitionAt ||
			view.TerminalAt.Value > durableTimeHighWater)) ||
		(view.NextRecoveryAt.Present && view.NextRecoveryAt.Value < view.LastTransitionAt) {
		return classified(ClassificationStale, "record-time-order")
	}
	if err := validateFailure(view); err != nil {
		return err
	}
	if err := validateRecovery(view); err != nil {
		return err
	}
	return validateStateShape(view)
}

func validateFailure(view RecordView) error {
	if view.FirstFailure == FailureNone {
		if view.FailureOperation != OperationNone {
			return classified(ClassificationBinding, "failure-operation-without-failure")
		}
		return nil
	}
	if !validOperation(view.FailureOperation) {
		return classified(ClassificationBinding, "failure-operation")
	}
	return nil
}

func validateRecovery(view RecordView) error {
	if view.AutomaticRecoveryCount > MaxAutomaticRecoveryAttempts {
		return classified(ClassificationCapacity, "automatic-recovery-count")
	}
	if view.AutomaticRecoveryCount == 0 && view.LastReconciliation != ReconciliationNone {
		return classified(ClassificationBinding, "reconciliation-without-attempt")
	}
	if view.LastReconciliation == ReconciliationIdentityMismatch && view.State != StateQuarantined {
		return classified(ClassificationBinding, "identity-mismatch-not-quarantined")
	}
	if view.LastReconciliation == ReconciliationAuthoritativelyAbsent && view.State != StateDestroyed {
		return classified(ClassificationBinding, "authoritative-absence-not-destroyed")
	}
	if view.AutomaticRecoveryExhausted {
		if view.AutomaticRecoveryCount != MaxAutomaticRecoveryAttempts ||
			view.LastReconciliation != ReconciliationUnknown ||
			view.State != StateUnresolved ||
			view.RecoveryFence != RecoveryFenceAutomaticExhausted ||
			view.NextRecoveryAt.Present {
			return classified(ClassificationBinding, "automatic-recovery-exhaustion")
		}
	}
	if view.NextRecoveryAt.Present && view.AutomaticRecoveryCount == 0 {
		return classified(ClassificationBinding, "next-recovery-without-attempt")
	}
	if view.RecoveryFence == RecoveryFenceIdentityMismatch && view.State != StateQuarantined {
		return classified(ClassificationBinding, "identity-fence-not-quarantined")
	}
	return nil
}

func validateStateShape(view RecordView) error {
	if view.State == StateDestroyed {
		if view.CleanupRequired || !view.TerminalAt.Present ||
			view.TerminalAt.Value != view.LastTransitionAt ||
			view.LastReconciliation != ReconciliationAuthoritativelyAbsent {
			return classified(ClassificationBinding, "destroyed-disposition")
		}
	} else {
		if !view.CleanupRequired || view.TerminalAt.Present {
			return classified(ClassificationBinding, "nonterminal-cleanup-disposition")
		}
	}

	switch view.State {
	case StatePreparePending:
		return validateEffectShape(view, OperationNone, EffectNone, CheckpointNone, false)
	case StatePrepareIntent:
		return validateIntentShape(view, OperationPrepare, CheckpointNone, false)
	case StatePrepared:
		return validateEffectShape(view, OperationPrepare, EffectConfirmed, CheckpointPrepare, false)
	case StateCreateIntent:
		return validateIntentShape(view, OperationCreate, CheckpointPrepare, false)
	case StateCreated:
		return validateEffectShape(view, OperationCreate, EffectConfirmed, CheckpointCreate, true)
	case StateStartIntent:
		return validateIntentShape(view, OperationStart, CheckpointCreate, true)
	case StateStarted:
		return validateEffectShape(view, OperationStart, EffectConfirmed, CheckpointStart, true)
	case StateObserveIntent:
		return validateIntentShape(view, OperationObserve, CheckpointStart, true)
	case StateObserved:
		return validateEffectShape(view, OperationObserve, EffectConfirmed, CheckpointObserve, true)
	case StateStopIntent:
		if !checkpointOneOf(view.LastConfirmedCheckpoint, CheckpointCreate, CheckpointStart, CheckpointObserve) {
			return classified(ClassificationBinding, "stop-intent-checkpoint")
		}
		return validateIntentShape(view, OperationStop, view.LastConfirmedCheckpoint, true)
	case StateStopped:
		return validateEffectShape(view, OperationStop, EffectConfirmed, CheckpointStop, true)
	case StateDestroyIntent:
		if !checkpointOneOf(view.LastConfirmedCheckpoint, CheckpointCreate, CheckpointStart, CheckpointObserve, CheckpointStop) {
			return classified(ClassificationBinding, "destroy-intent-checkpoint")
		}
		return validateIntentShape(view, OperationDestroy, view.LastConfirmedCheckpoint, true)
	case StateDestroyConfirmed:
		return validateEffectShape(view, OperationDestroy, EffectConfirmed, CheckpointDestroy, true)
	case StateDestroyed:
		return validateDestroyedShape(view)
	case StateUnresolved, StateQuarantined:
		return validateRetainedFailureShape(view)
	default:
		return classified(ClassificationUnsupported, "lifecycle-state")
	}
}

func validateDestroyedShape(view RecordView) error {
	if view.Instance.Present() {
		return validateEffectShape(view, OperationDestroy, EffectConfirmed, CheckpointDestroy, true)
	}
	switch {
	case view.Operation == OperationNone && view.EffectStatus == EffectNone &&
		view.LastConfirmedCheckpoint == CheckpointNone:
		return validateEffectShape(view, OperationNone, EffectNone, CheckpointNone, false)
	case view.Operation == OperationPrepare &&
		(view.EffectStatus == EffectIntent || view.EffectStatus == EffectIndeterminate) &&
		view.LastConfirmedCheckpoint == CheckpointNone:
		return validateEffectShape(view, OperationPrepare, view.EffectStatus, CheckpointNone, false)
	case view.Operation == OperationPrepare && view.EffectStatus == EffectConfirmed &&
		view.LastConfirmedCheckpoint == CheckpointPrepare:
		return validateEffectShape(view, OperationPrepare, EffectConfirmed, CheckpointPrepare, false)
	case view.Operation == OperationCreate &&
		(view.EffectStatus == EffectIntent || view.EffectStatus == EffectIndeterminate) &&
		view.LastConfirmedCheckpoint == CheckpointPrepare:
		return validateEffectShape(view, OperationCreate, view.EffectStatus, CheckpointPrepare, false)
	default:
		return classified(ClassificationBinding, "destroyed-without-exact-absence-shape")
	}
}

func validateIntentShape(
	view RecordView,
	operation Operation,
	checkpoint Checkpoint,
	wantInstance bool,
) error {
	if view.EffectStatus != EffectIntent && view.EffectStatus != EffectIndeterminate {
		return classified(ClassificationBinding, "intent-effect-status")
	}
	return validateEffectShape(view, operation, view.EffectStatus, checkpoint, wantInstance)
}

func validateEffectShape(
	view RecordView,
	operation Operation,
	status EffectStatus,
	checkpoint Checkpoint,
	wantInstance bool,
) error {
	if view.Operation != operation || view.EffectStatus != status ||
		view.LastConfirmedCheckpoint != checkpoint || view.Instance.Present() != wantInstance {
		return classified(ClassificationBinding, "state-effect-shape")
	}
	if status == EffectNone {
		if view.OperationSequence != 0 || !view.EffectID.IsZero() {
			return classified(ClassificationBinding, "empty-effect-identity")
		}
		return nil
	}
	if !validOperation(operation) || view.OperationSequence == 0 ||
		uint64(view.OperationSequence) > v0candidate.MaxSafeInteger || view.EffectID.IsZero() {
		return classified(ClassificationBinding, "effect-identity")
	}
	return nil
}

func validateRetainedFailureShape(view RecordView) error {
	if view.EffectStatus == EffectNone {
		if view.Operation != OperationNone || view.OperationSequence != 0 || !view.EffectID.IsZero() {
			return classified(ClassificationBinding, "retained-empty-effect")
		}
		return nil
	}
	if !validOperation(view.Operation) || view.OperationSequence == 0 ||
		uint64(view.OperationSequence) > v0candidate.MaxSafeInteger || view.EffectID.IsZero() {
		return classified(ClassificationBinding, "retained-effect-identity")
	}
	if view.LastConfirmedCheckpoint != CheckpointNone &&
		!validCheckpoint(view.LastConfirmedCheckpoint) {
		return classified(ClassificationBinding, "retained-checkpoint")
	}
	if checkpointAtOrAfterCreate(view.LastConfirmedCheckpoint) && !view.Instance.Present() {
		return classified(ClassificationBinding, "retained-instance-missing")
	}
	return nil
}

func checkpointOneOf(value Checkpoint, choices ...Checkpoint) bool {
	for _, choice := range choices {
		if value == choice {
			return true
		}
	}
	return false
}

func checkpointAtOrAfterCreate(checkpoint Checkpoint) bool {
	return checkpointOneOf(
		checkpoint,
		CheckpointCreate,
		CheckpointStart,
		CheckpointObserve,
		CheckpointStop,
		CheckpointDestroy,
	)
}
