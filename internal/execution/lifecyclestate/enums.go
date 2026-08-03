package lifecyclestate

// Operation names only the six durable lifecycle effects. Reconciliation is
// observation-only and has its own closed result vocabulary.
type Operation string

const (
	OperationNone    Operation = "none"
	OperationPrepare Operation = "prepare"
	OperationCreate  Operation = "create"
	OperationStart   Operation = "start"
	OperationObserve Operation = "observe"
	OperationStop    Operation = "stop"
	OperationDestroy Operation = "destroy"
)

var operations = [...]Operation{
	OperationPrepare,
	OperationCreate,
	OperationStart,
	OperationObserve,
	OperationStop,
	OperationDestroy,
}

func Operations() []Operation { return append([]Operation(nil), operations[:]...) }

func validOperation(value Operation) bool {
	for _, candidate := range operations {
		if value == candidate {
			return true
		}
	}
	return false
}

// LifecycleState is the complete ADR-0025 first-slice state vocabulary.
type LifecycleState string

// RecordState is the ADR-0025 name for LifecycleState.
type RecordState = LifecycleState

const (
	StatePreparePending   LifecycleState = "prepare-pending"
	StatePrepareIntent    LifecycleState = "prepare-intent"
	StatePrepared         LifecycleState = "prepared"
	StateCreateIntent     LifecycleState = "create-intent"
	StateCreated          LifecycleState = "created"
	StateStartIntent      LifecycleState = "start-intent"
	StateStarted          LifecycleState = "started"
	StateObserveIntent    LifecycleState = "observe-intent"
	StateObserved         LifecycleState = "observed"
	StateStopIntent       LifecycleState = "stop-intent"
	StateStopped          LifecycleState = "stopped"
	StateDestroyIntent    LifecycleState = "destroy-intent"
	StateDestroyConfirmed LifecycleState = "destroy-confirmed"
	StateDestroyed        LifecycleState = "destroyed"
	StateUnresolved       LifecycleState = "unresolved"
	StateQuarantined      LifecycleState = "quarantined"
)

var lifecycleStates = [...]LifecycleState{
	StatePreparePending,
	StatePrepareIntent,
	StatePrepared,
	StateCreateIntent,
	StateCreated,
	StateStartIntent,
	StateStarted,
	StateObserveIntent,
	StateObserved,
	StateStopIntent,
	StateStopped,
	StateDestroyIntent,
	StateDestroyConfirmed,
	StateDestroyed,
	StateUnresolved,
	StateQuarantined,
}

func LifecycleStates() []LifecycleState {
	return append([]LifecycleState(nil), lifecycleStates[:]...)
}

func validLifecycleState(value LifecycleState) bool {
	for _, candidate := range lifecycleStates {
		if value == candidate {
			return true
		}
	}
	return false
}

// Checkpoint records the last durably confirmed lifecycle effect.
type Checkpoint string

const (
	CheckpointNone    Checkpoint = "none"
	CheckpointPrepare Checkpoint = "prepare"
	CheckpointCreate  Checkpoint = "create"
	CheckpointStart   Checkpoint = "start"
	CheckpointObserve Checkpoint = "observe"
	CheckpointStop    Checkpoint = "stop"
	CheckpointDestroy Checkpoint = "destroy"
)

var checkpoints = [...]Checkpoint{
	CheckpointNone,
	CheckpointPrepare,
	CheckpointCreate,
	CheckpointStart,
	CheckpointObserve,
	CheckpointStop,
	CheckpointDestroy,
}

func Checkpoints() []Checkpoint { return append([]Checkpoint(nil), checkpoints[:]...) }

func validCheckpoint(value Checkpoint) bool {
	for _, candidate := range checkpoints {
		if value == candidate {
			return true
		}
	}
	return false
}

// EffectStatus records the durable status of the current logical effect.
type EffectStatus string

const (
	EffectNone          EffectStatus = "none"
	EffectIntent        EffectStatus = "intent"
	EffectConfirmed     EffectStatus = "confirmed"
	EffectIndeterminate EffectStatus = "indeterminate"
)

var effectStatuses = [...]EffectStatus{
	EffectNone,
	EffectIntent,
	EffectConfirmed,
	EffectIndeterminate,
}

func EffectStatuses() []EffectStatus {
	return append([]EffectStatus(nil), effectStatuses[:]...)
}

func validEffectStatus(value EffectStatus) bool {
	for _, candidate := range effectStatuses {
		if value == candidate {
			return true
		}
	}
	return false
}

// FailureClassification is the first fixed lifecycle failure retained by a
// record. Empty or free-form failure text is never part of the contract.
type FailureClassification string

const (
	FailureNone              FailureClassification = "none"
	FailureMalformed         FailureClassification = "MALFORMED"
	FailureUnsupported       FailureClassification = "UNSUPPORTED"
	FailureSchema            FailureClassification = "SCHEMA"
	FailureDomain            FailureClassification = "DOMAIN"
	FailureBinding           FailureClassification = "BINDING"
	FailureStale             FailureClassification = "STALE"
	FailureCapacity          FailureClassification = "CAPACITY"
	FailureLocal             FailureClassification = "LOCAL_FAILURE"
	FailureLifecycle         FailureClassification = "LIFECYCLE_FAILURE"
	FailureCleanupUnresolved FailureClassification = "CLEANUP_UNRESOLVED"
	FailureTrustState        FailureClassification = "TRUST_STATE"
	FailureRecoveryRequired  FailureClassification = "RECOVERY_REQUIRED"
)

var failureClassifications = [...]FailureClassification{
	FailureNone,
	FailureMalformed,
	FailureUnsupported,
	FailureSchema,
	FailureDomain,
	FailureBinding,
	FailureStale,
	FailureCapacity,
	FailureLocal,
	FailureLifecycle,
	FailureCleanupUnresolved,
	FailureTrustState,
	FailureRecoveryRequired,
}

func FailureClassifications() []FailureClassification {
	return append([]FailureClassification(nil), failureClassifications[:]...)
}

func validFailureClassification(value FailureClassification) bool {
	for _, candidate := range failureClassifications {
		if value == candidate {
			return true
		}
	}
	return false
}

// ReconciliationStatus is the complete observation result vocabulary.
type ReconciliationStatus string

// ReconcileStatus is the short form used by the future store port.
type ReconcileStatus = ReconciliationStatus

const (
	ReconciliationNone                  ReconciliationStatus = "none"
	ReconciliationApplied               ReconciliationStatus = "applied"
	ReconciliationNotApplied            ReconciliationStatus = "not-applied"
	ReconciliationAuthoritativelyAbsent ReconciliationStatus = "authoritatively-absent"
	ReconciliationPresent               ReconciliationStatus = "present"
	ReconciliationUnknown               ReconciliationStatus = "unknown"
	ReconciliationIdentityMismatch      ReconciliationStatus = "identity-mismatch"
)

var reconciliationStatuses = [...]ReconciliationStatus{
	ReconciliationNone,
	ReconciliationApplied,
	ReconciliationNotApplied,
	ReconciliationAuthoritativelyAbsent,
	ReconciliationPresent,
	ReconciliationUnknown,
	ReconciliationIdentityMismatch,
}

func ReconciliationStatuses() []ReconciliationStatus {
	return append([]ReconciliationStatus(nil), reconciliationStatuses[:]...)
}

func validReconciliationStatus(value ReconciliationStatus) bool {
	for _, candidate := range reconciliationStatuses {
		if value == candidate {
			return true
		}
	}
	return false
}

// BackendKind is deliberately closed to the no-guest fake adapter for E1.
type BackendKind string

type BackendAdapterKind = BackendKind

const BackendFakeNoGuest BackendKind = "fake-no-guest"

func BackendKinds() []BackendKind { return []BackendKind{BackendFakeNoGuest} }

// BackendInstanceKind is deliberately closed to the bounded fake identity.
type BackendInstanceKind string

type InstanceIdentityKind = BackendInstanceKind

const BackendInstanceFake BackendInstanceKind = "fake-instance"

func BackendInstanceKinds() []BackendInstanceKind {
	return []BackendInstanceKind{BackendInstanceFake}
}

// EffectResultStatus classifies one adapter response without backend text.
type EffectResultStatus string

const (
	EffectResultApplied       EffectResultStatus = "applied"
	EffectResultNotApplied    EffectResultStatus = "not-applied"
	EffectResultIndeterminate EffectResultStatus = "indeterminate"
)

// RecoveryFenceReason is a closed explanation for why forward mutation must
// remain disabled. It does not contain an adapter or storage error string.
type RecoveryFenceReason string

const (
	RecoveryFenceNone                RecoveryFenceReason = "none"
	RecoveryFenceCommitIndeterminate RecoveryFenceReason = "commit-indeterminate"
	RecoveryFenceCommitFailure       RecoveryFenceReason = "commit-failure"
	RecoveryFenceReconcileUnknown    RecoveryFenceReason = "reconcile-unknown"
	RecoveryFenceIdentityMismatch    RecoveryFenceReason = "identity-mismatch"
	RecoveryFenceAutomaticExhausted  RecoveryFenceReason = "automatic-recovery-exhausted"
	RecoveryFenceRepairRequired      RecoveryFenceReason = "repair-required"
)

func validRecoveryFenceReason(value RecoveryFenceReason) bool {
	switch value {
	case RecoveryFenceNone,
		RecoveryFenceCommitIndeterminate,
		RecoveryFenceCommitFailure,
		RecoveryFenceReconcileUnknown,
		RecoveryFenceIdentityMismatch,
		RecoveryFenceAutomaticExhausted,
		RecoveryFenceRepairRequired:
		return true
	default:
		return false
	}
}
