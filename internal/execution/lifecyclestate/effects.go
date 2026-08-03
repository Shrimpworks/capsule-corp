package lifecyclestate

import (
	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// EffectPermitView is the exact passive projection sealed by a future store
// after an intent commits. It contains no adapter flags, paths, plan bytes, or
// caller-provided backend configuration.
type EffectPermitView struct {
	AttemptID              approvalattempt.AttemptID
	RecordVersion          RecordVersion
	OperationSequence      v0candidate.PositiveUInt53
	Operation              Operation
	EffectID               EffectID
	ImmutableBindingDigest ImmutableBindingDigest
	OwnerSessionID         OwnerSessionID
	Instance               BackendInstanceIdentity
}

type EffectPermit struct {
	view EffectPermitView
}

func NewEffectPermit(view EffectPermitView) (EffectPermit, error) {
	if view.AttemptID == (approvalattempt.AttemptID{}) ||
		view.RecordVersion == 0 || uint64(view.RecordVersion) > v0candidate.MaxSafeInteger ||
		view.OperationSequence == 0 ||
		uint64(view.OperationSequence) > v0candidate.MaxSafeInteger ||
		!validOperation(view.Operation) || view.EffectID.IsZero() ||
		zeroDigest(view.ImmutableBindingDigest) || view.OwnerSessionID.IsZero() ||
		!view.Instance.valid() {
		return EffectPermit{}, classified(ClassificationBinding, "effect-permit")
	}
	wantInstance := operationRequiresInstance(view.Operation)
	if view.Instance.Present() != wantInstance {
		return EffectPermit{}, classified(ClassificationBinding, "effect-permit-instance")
	}
	return EffectPermit{view: cloneEffectPermitView(view)}, nil
}

func (permit EffectPermit) View() EffectPermitView { return cloneEffectPermitView(permit.view) }

func cloneEffectPermitView(view EffectPermitView) EffectPermitView {
	if view.Instance.Present() {
		instance, err := RestoreBackendInstanceIdentity(view.Instance.View())
		if err == nil {
			view.Instance = instance
		}
	}
	return view
}

func operationRequiresInstance(operation Operation) bool {
	switch operation {
	case OperationStart, OperationObserve, OperationStop, OperationDestroy:
		return true
	case OperationPrepare, OperationCreate:
		return false
	default:
		return false
	}
}

// EffectResult is a closed adapter response. Applied create is the only
// effect response that can introduce an instance identity.
type EffectResult struct {
	operation Operation
	status    EffectResultStatus
	instance  BackendInstanceIdentity
}

func NewEffectResult(
	operation Operation,
	status EffectResultStatus,
	instance BackendInstanceIdentity,
) (EffectResult, error) {
	if !validOperation(operation) || !validEffectResultStatus(status) || !instance.valid() {
		return EffectResult{}, classified(ClassificationUnsupported, "effect-result")
	}
	wantInstance := operation == OperationCreate && status == EffectResultApplied
	if instance.Present() != wantInstance {
		return EffectResult{}, classified(ClassificationBinding, "effect-result-instance")
	}
	return EffectResult{operation: operation, status: status, instance: cloneInstance(instance)}, nil
}

func (result EffectResult) Operation() Operation { return result.operation }

func (result EffectResult) Status() EffectResultStatus { return result.status }

func (result EffectResult) Instance() BackendInstanceIdentity { return cloneInstance(result.instance) }

func validEffectResultStatus(status EffectResultStatus) bool {
	switch status {
	case EffectResultApplied, EffectResultNotApplied, EffectResultIndeterminate:
		return true
	default:
		return false
	}
}

// ReconcileResult is an observation-only closed result. It authorizes no
// effect and retains only an exact bounded identity where the status requires
// one.
type ReconcileResult struct {
	operation Operation
	status    ReconciliationStatus
	instance  BackendInstanceIdentity
}

func NewReconcileResult(
	operation Operation,
	status ReconciliationStatus,
	instance BackendInstanceIdentity,
) (ReconcileResult, error) {
	if !validOperation(operation) || status == ReconciliationNone ||
		!validReconciliationStatus(status) || !instance.valid() {
		return ReconcileResult{}, classified(ClassificationUnsupported, "reconcile-result")
	}
	wantInstance := status == ReconciliationPresent ||
		status == ReconciliationIdentityMismatch ||
		(status == ReconciliationApplied && operation == OperationCreate)
	if instance.Present() != wantInstance {
		return ReconcileResult{}, classified(ClassificationBinding, "reconcile-result-instance")
	}
	return ReconcileResult{
		operation: operation, status: status, instance: cloneInstance(instance),
	}, nil
}

func (result ReconcileResult) Operation() Operation { return result.operation }

func (result ReconcileResult) Status() ReconciliationStatus { return result.status }

func (result ReconcileResult) Instance() BackendInstanceIdentity {
	return cloneInstance(result.instance)
}

func cloneInstance(instance BackendInstanceIdentity) BackendInstanceIdentity {
	if !instance.Present() {
		return BackendInstanceIdentity{}
	}
	cloned, err := RestoreBackendInstanceIdentity(instance.View())
	if err != nil {
		return BackendInstanceIdentity{}
	}
	return cloned
}
