package registeredlifecycle

import (
	"bytes"
	"context"
	"errors"
	"sync"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
)

var errInjectedFakeFault = errors.New("injected fake lifecycle fault")

const maxInjectedFaultKeys = 4_096 * 7 * 2

type fakeFaultKey struct {
	attemptID approvalattempt.AttemptID
	operation Operation
	moment    FaultMoment
}

type fakeEffectRecord struct {
	permit lifecyclestate.EffectPermitView
	result lifecyclestate.EffectResult
}

type fakeInstance struct {
	identity  lifecyclestate.BackendInstanceIdentity
	prepared  bool
	created   bool
	started   bool
	observed  bool
	stopped   bool
	destroyed bool
	calls     map[Operation]uint64
	applied   map[Operation]uint64
	effectIDs map[Operation]lifecyclestate.EffectID
}

// FakeBackend is a closed, no-guest lifecycle simulator. Every logical effect
// is keyed by the store-issued EffectID and is applied at most once. Apart
// from its fixed fault oracle, it has no import, command, subprocess, VMM,
// container, network, filesystem-content, runtime-launch, or guest path.
type FakeBackend struct {
	mu        sync.Mutex
	binding   lifecyclestate.BackendBinding
	instances map[approvalattempt.AttemptID]*fakeInstance
	effects   map[lifecyclestate.EffectID]fakeEffectRecord
	faults    map[fakeFaultKey]uint64
}

func NewFakeBackend(binding lifecyclestate.BackendBinding) (*FakeBackend, error) {
	validated, err := lifecyclestate.NewBackendBinding(binding.View())
	if err != nil || validated.View().CreatesGuest ||
		validated.View().Kind != lifecyclestate.BackendFakeNoGuest {
		return nil, errors.New("fake backend requires the closed no-guest binding")
	}
	return &FakeBackend{
		binding:   validated,
		instances: make(map[approvalattempt.AttemptID]*fakeInstance),
		effects:   make(map[lifecyclestate.EffectID]fakeEffectRecord),
		faults:    make(map[fakeFaultKey]uint64),
	}, nil
}

func (*FakeBackend) CreatesGuest() bool { return false }

func (backend *FakeBackend) Binding() lifecyclestate.BackendBinding {
	return backend.binding
}

// InjectFault adds one deterministic, attempt-scoped fault immediately before
// or after the selected fake effect or reconciliation observation.
func (backend *FakeBackend) InjectFault(
	attemptID approvalattempt.AttemptID,
	operation Operation,
	moment FaultMoment,
) error {
	if attemptID == (approvalattempt.AttemptID{}) {
		return errors.New("fake fault attempt ID is required")
	}
	if !validOperation(operation) || (moment != FaultBeforeEffect && moment != FaultAfterEffect) {
		return errors.New("fake fault operation and moment must be closed values")
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	key := fakeFaultKey{attemptID: attemptID, operation: operation, moment: moment}
	if _, exists := backend.faults[key]; !exists && len(backend.faults) >= maxInjectedFaultKeys {
		return errors.New("fake fault table reached its fixed capacity")
	}
	if backend.faults[key] == ^uint64(0) {
		return errors.New("fake fault count reached its fixed capacity")
	}
	backend.faults[key]++
	return nil
}

func validOperation(operation Operation) bool {
	switch operation {
	case OperationPrepare, OperationCreate, OperationStart, OperationObserve,
		OperationStop, OperationDestroy, OperationReconcile:
		return true
	default:
		return false
	}
}

func (backend *FakeBackend) takeFaultLocked(
	attemptID approvalattempt.AttemptID,
	operation Operation,
	moment FaultMoment,
) bool {
	key := fakeFaultKey{attemptID: attemptID, operation: operation, moment: moment}
	remaining := backend.faults[key]
	if remaining == 0 {
		return false
	}
	if remaining == 1 {
		delete(backend.faults, key)
	} else {
		backend.faults[key] = remaining - 1
	}
	return true
}

func (backend *FakeBackend) instanceLocked(attemptID approvalattempt.AttemptID) *fakeInstance {
	instance := backend.instances[attemptID]
	if instance == nil {
		instance = &fakeInstance{
			calls:     make(map[Operation]uint64),
			applied:   make(map[Operation]uint64),
			effectIDs: make(map[Operation]lifecyclestate.EffectID),
		}
		backend.instances[attemptID] = instance
	}
	return instance
}

func operationForDurable(operation lifecyclestate.Operation) Operation {
	return Operation(operation)
}

func durableOperation(operation Operation) lifecyclestate.Operation {
	return lifecyclestate.Operation(operation)
}

func samePermit(left, right lifecyclestate.EffectPermitView) bool {
	return left.AttemptID == right.AttemptID &&
		left.OperationSequence == right.OperationSequence &&
		left.Operation == right.Operation &&
		left.EffectID == right.EffectID &&
		left.ImmutableBindingDigest == right.ImmutableBindingDigest &&
		instancesEqual(left.Instance, right.Instance)
}

func instancesEqual(
	left lifecyclestate.BackendInstanceIdentity,
	right lifecyclestate.BackendInstanceIdentity,
) bool {
	return left.Present() == right.Present() && left.Kind() == right.Kind() &&
		left.Digest() == right.Digest() && bytes.Equal(left.Value(), right.Value())
}

func (backend *FakeBackend) apply(
	ctx context.Context,
	permit lifecyclestate.EffectPermit,
) (lifecyclestate.EffectResult, error) {
	if err := ctx.Err(); err != nil {
		return lifecyclestate.EffectResult{}, err
	}
	view := permit.View()
	operation := operationForDurable(view.Operation)
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance := backend.instanceLocked(view.AttemptID)
	instance.calls[operation]++

	if cached, exists := backend.effects[view.EffectID]; exists {
		if !samePermit(cached.permit, view) {
			return lifecyclestate.EffectResult{}, errors.New("fake effect identity binding mismatch")
		}
		return cached.result, nil
	}
	if backend.takeFaultLocked(view.AttemptID, operation, FaultBeforeEffect) {
		result, err := lifecyclestate.NewEffectResult(
			view.Operation,
			lifecyclestate.EffectResultNotApplied,
			lifecyclestate.BackendInstanceIdentity{},
		)
		return result, errors.Join(errInjectedFakeFault, err)
	}

	var resultInstance lifecyclestate.BackendInstanceIdentity
	switch view.Operation {
	case lifecyclestate.OperationPrepare:
		instance.prepared = true
	case lifecyclestate.OperationCreate:
		if !instance.prepared || instance.destroyed {
			return lifecyclestate.EffectResult{}, errors.New("fake create state rejected")
		}
		if !instance.created {
			identity, err := lifecyclestate.NewBackendInstanceIdentity(
				lifecyclestate.BackendInstanceFake,
				append(append([]byte(nil), view.AttemptID[:]...), view.EffectID[:]...),
			)
			if err != nil {
				return lifecyclestate.EffectResult{}, err
			}
			instance.identity = identity
			instance.created = true
		}
		resultInstance = instance.identity
	case lifecyclestate.OperationStart:
		if !instance.created || instance.destroyed || !instancesEqual(instance.identity, view.Instance) {
			return lifecyclestate.EffectResult{}, errors.New("fake start state rejected")
		}
		instance.started = true
	case lifecyclestate.OperationObserve:
		if !instance.started || instance.destroyed || !instancesEqual(instance.identity, view.Instance) {
			return lifecyclestate.EffectResult{}, errors.New("fake observe state rejected")
		}
		instance.observed = true
	case lifecyclestate.OperationStop:
		if !instance.created || instance.destroyed || !instancesEqual(instance.identity, view.Instance) {
			return lifecyclestate.EffectResult{}, errors.New("fake stop state rejected")
		}
		instance.stopped = true
	case lifecyclestate.OperationDestroy:
		if !instance.created || instance.destroyed || !instancesEqual(instance.identity, view.Instance) {
			return lifecyclestate.EffectResult{}, errors.New("fake destroy state rejected")
		}
		instance.destroyed = true
	default:
		return lifecyclestate.EffectResult{}, errors.New("fake operation rejected")
	}
	instance.applied[operation]++
	instance.effectIDs[operation] = view.EffectID
	result, err := lifecyclestate.NewEffectResult(
		view.Operation,
		lifecyclestate.EffectResultApplied,
		resultInstance,
	)
	if err != nil {
		return lifecyclestate.EffectResult{}, err
	}
	backend.effects[view.EffectID] = fakeEffectRecord{permit: view, result: result}
	if backend.takeFaultLocked(view.AttemptID, operation, FaultAfterEffect) {
		indeterminate, indeterminateErr := lifecyclestate.NewEffectResult(
			view.Operation,
			lifecyclestate.EffectResultIndeterminate,
			lifecyclestate.BackendInstanceIdentity{},
		)
		return indeterminate, errors.Join(errInjectedFakeFault, indeterminateErr)
	}
	return result, nil
}

func (backend *FakeBackend) reconcile(
	ctx context.Context,
	record lifecyclestate.Record,
) (lifecyclestate.ReconcileResult, error) {
	if err := ctx.Err(); err != nil {
		return lifecyclestate.ReconcileResult{}, err
	}
	view := record.View()
	attemptID := view.Bindings.View().AttemptID
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance := backend.instanceLocked(attemptID)
	instance.calls[OperationReconcile]++
	unknown := func() (lifecyclestate.ReconcileResult, error) {
		result, err := lifecyclestate.NewReconcileResult(
			view.Operation,
			lifecyclestate.ReconciliationUnknown,
			lifecyclestate.BackendInstanceIdentity{},
		)
		return result, err
	}
	if backend.takeFaultLocked(attemptID, OperationReconcile, FaultBeforeEffect) {
		result, err := unknown()
		return result, errors.Join(errInjectedFakeFault, err)
	}

	status := lifecyclestate.ReconciliationNotApplied
	var resultInstance lifecyclestate.BackendInstanceIdentity
	if view.Instance.Present() {
		if !instance.created {
			status = lifecyclestate.ReconciliationUnknown
		} else if !instancesEqual(instance.identity, view.Instance) {
			status = lifecyclestate.ReconciliationIdentityMismatch
			resultInstance = instance.identity
		}
	}
	if status == lifecyclestate.ReconciliationNotApplied {
		if view.State == lifecyclestate.StateDestroyConfirmed && instance.destroyed &&
			instancesEqual(instance.identity, view.Instance) {
			status = lifecyclestate.ReconciliationAuthoritativelyAbsent
		} else if cached, exists := backend.effects[view.EffectID]; exists {
			if !samePermit(cached.permit, lifecyclestate.EffectPermitView{
				AttemptID:              attemptID,
				OperationSequence:      view.OperationSequence,
				Operation:              view.Operation,
				EffectID:               view.EffectID,
				ImmutableBindingDigest: view.ImmutableBindingDigest,
				OwnerSessionID:         cached.permit.OwnerSessionID,
				Instance:               view.Instance,
			}) {
				status = lifecyclestate.ReconciliationIdentityMismatch
				if instance.identity.Present() {
					resultInstance = instance.identity
				}
			} else {
				status = lifecyclestate.ReconciliationApplied
				if view.Operation == lifecyclestate.OperationCreate {
					resultInstance = cached.result.Instance()
				}
			}
		}
	}
	result, err := lifecyclestate.NewReconcileResult(view.Operation, status, resultInstance)
	if err != nil {
		return lifecyclestate.ReconcileResult{}, err
	}
	if backend.takeFaultLocked(attemptID, OperationReconcile, FaultAfterEffect) {
		unknownResult, unknownErr := unknown()
		return unknownResult, errors.Join(errInjectedFakeFault, unknownErr)
	}
	return result, nil
}

// BackendSnapshot is bounded passive test/development evidence. Maps are
// defensive copies; effect applications count logical effects, not retries.
type BackendSnapshot struct {
	Prepared          bool
	Created           bool
	Started           bool
	Observed          bool
	Stopped           bool
	Destroyed         bool
	InstanceDigest    lifecyclestate.BackendInstanceDigest
	CallCounts        map[Operation]uint64
	ApplicationCounts map[Operation]uint64
	EffectIDs         map[Operation]lifecyclestate.EffectID
}

func (backend *FakeBackend) Snapshot(attemptID approvalattempt.AttemptID) BackendSnapshot {
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance := backend.instances[attemptID]
	if instance == nil {
		return BackendSnapshot{
			CallCounts: make(map[Operation]uint64), ApplicationCounts: make(map[Operation]uint64),
			EffectIDs: make(map[Operation]lifecyclestate.EffectID),
		}
	}
	calls := make(map[Operation]uint64, len(instance.calls))
	for operation, count := range instance.calls {
		calls[operation] = count
	}
	applied := make(map[Operation]uint64, len(instance.applied))
	for operation, count := range instance.applied {
		applied[operation] = count
	}
	effectIDs := make(map[Operation]lifecyclestate.EffectID, len(instance.effectIDs))
	for operation, effectID := range instance.effectIDs {
		effectIDs[operation] = effectID
	}
	return BackendSnapshot{
		Prepared: instance.prepared, Created: instance.created, Started: instance.started,
		Observed: instance.observed, Stopped: instance.stopped, Destroyed: instance.destroyed,
		InstanceDigest: instance.identity.Digest(), CallCounts: calls, ApplicationCounts: applied,
		EffectIDs: effectIDs,
	}
}

// forgetAttemptState and replaceAttemptIdentity are closed local fault oracles
// used by the retained E4 missing/stale/identity tests.
func (backend *FakeBackend) forgetAttemptState(attemptID approvalattempt.AttemptID) {
	backend.mu.Lock()
	defer backend.mu.Unlock()
	delete(backend.instances, attemptID)
}

func (backend *FakeBackend) replaceAttemptIdentity(
	attemptID approvalattempt.AttemptID,
	identity lifecyclestate.BackendInstanceIdentity,
) {
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance := backend.instanceLocked(attemptID)
	instance.created = true
	instance.identity = identity
}
