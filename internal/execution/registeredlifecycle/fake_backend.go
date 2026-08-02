package registeredlifecycle

import (
	"context"
	"errors"
	"sync"

	"capsule.local/capsule/internal/execution/approvalattempt"
)

var errInjectedFakeFault = errors.New("injected fake lifecycle fault")

const maxInjectedFaultKeys = MaxLifecycleRecords * 7 * 2

type fakeObservationStatus string

const (
	fakePresent fakeObservationStatus = "present"
	fakeAbsent  fakeObservationStatus = "authoritatively-absent"
	fakeUnknown fakeObservationStatus = "unknown"
)

type fakeObservation struct {
	status fakeObservationStatus
	handle fakeHandle
}

type fakeFaultKey struct {
	attemptID approvalattempt.AttemptID
	operation Operation
	moment    FaultMoment
}

type fakeInstance struct {
	handle    fakeHandle
	prepared  bool
	created   bool
	started   bool
	observed  bool
	stopped   bool
	destroyed bool
	calls     map[Operation]uint64
}

// FakeBackend is a closed in-memory lifecycle simulator. Apart from its fixed
// fault oracle, it has no import, link, runtime configuration, command,
// subprocess, VMM, container, network, filesystem-content, runtime-launch, or
// guest-creation path.
type FakeBackend struct {
	mu        sync.Mutex
	next      fakeHandle
	instances map[approvalattempt.AttemptID]*fakeInstance
	faults    map[fakeFaultKey]uint64
}

func NewFakeBackend() *FakeBackend {
	return &FakeBackend{
		next:      1,
		instances: make(map[approvalattempt.AttemptID]*fakeInstance),
		faults:    make(map[fakeFaultKey]uint64),
	}
}

func (*FakeBackend) CreatesGuest() bool { return false }

// InjectFault adds one deterministic, attempt-scoped fault. The fault is
// consumed immediately before or immediately after the selected fake effect.
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
	if backend.faults == nil {
		backend.faults = make(map[fakeFaultKey]uint64)
	}
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

func (backend *FakeBackend) instanceLocked(
	attemptID approvalattempt.AttemptID,
) *fakeInstance {
	if backend.instances == nil {
		backend.instances = make(map[approvalattempt.AttemptID]*fakeInstance)
	}
	if backend.next == 0 {
		backend.next = 1
	}
	instance := backend.instances[attemptID]
	if instance == nil {
		instance = &fakeInstance{calls: make(map[Operation]uint64)}
		backend.instances[attemptID] = instance
	}
	return instance
}

func (backend *FakeBackend) beforeLocked(
	attemptID approvalattempt.AttemptID,
	operation Operation,
) (*fakeInstance, error) {
	instance := backend.instanceLocked(attemptID)
	instance.calls[operation]++
	if backend.takeFaultLocked(attemptID, operation, FaultBeforeEffect) {
		return instance, errInjectedFakeFault
	}
	return instance, nil
}

func (backend *FakeBackend) afterLocked(
	attemptID approvalattempt.AttemptID,
	operation Operation,
) error {
	if backend.takeFaultLocked(attemptID, operation, FaultAfterEffect) {
		return errInjectedFakeFault
	}
	return nil
}

func (backend *FakeBackend) prepare(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(attemptID, OperationPrepare)
	if err != nil {
		return err
	}
	instance.prepared = true
	return backend.afterLocked(attemptID, OperationPrepare)
}

func (backend *FakeBackend) create(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (fakeHandle, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(attemptID, OperationCreate)
	if err != nil {
		return 0, err
	}
	if !instance.prepared || instance.destroyed {
		return 0, errors.New("fake create state rejected")
	}
	if !instance.created {
		instance.handle = backend.next
		backend.next++
		instance.created = true
	}
	if err := backend.afterLocked(attemptID, OperationCreate); err != nil {
		return 0, err
	}
	return instance.handle, nil
}

func (backend *FakeBackend) start(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	handle fakeHandle,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(attemptID, OperationStart)
	if err != nil {
		return err
	}
	if !instance.created || instance.destroyed || instance.handle != handle {
		return errors.New("fake start state rejected")
	}
	instance.started = true
	return backend.afterLocked(attemptID, OperationStart)
}

func (backend *FakeBackend) observe(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	handle fakeHandle,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(attemptID, OperationObserve)
	if err != nil {
		return err
	}
	if !instance.started || instance.destroyed || instance.handle != handle {
		return errors.New("fake observe state rejected")
	}
	instance.observed = true
	return backend.afterLocked(attemptID, OperationObserve)
}

func (backend *FakeBackend) stop(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	handle fakeHandle,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(attemptID, OperationStop)
	if err != nil {
		return err
	}
	if !instance.created || instance.destroyed || instance.handle != handle {
		return errors.New("fake stop state rejected")
	}
	instance.stopped = true
	return backend.afterLocked(attemptID, OperationStop)
}

func (backend *FakeBackend) destroy(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	handle fakeHandle,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(attemptID, OperationDestroy)
	if err != nil {
		return err
	}
	if !instance.created || instance.destroyed || instance.handle != handle {
		return errors.New("fake destroy state rejected")
	}
	instance.destroyed = true
	return backend.afterLocked(attemptID, OperationDestroy)
}

func (backend *FakeBackend) reconcile(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (fakeObservation, error) {
	if err := ctx.Err(); err != nil {
		return fakeObservation{status: fakeUnknown}, err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(attemptID, OperationReconcile)
	if err != nil {
		return fakeObservation{status: fakeUnknown}, err
	}
	observation := fakeObservation{status: fakeAbsent}
	if instance.created && !instance.destroyed {
		observation = fakeObservation{status: fakePresent, handle: instance.handle}
	}
	if err := backend.afterLocked(attemptID, OperationReconcile); err != nil {
		return fakeObservation{status: fakeUnknown}, err
	}
	return observation, nil
}

// BackendSnapshot is bounded passive test/development evidence. CallCounts is
// returned as a defensive copy.
type BackendSnapshot struct {
	Prepared   bool
	Created    bool
	Started    bool
	Observed   bool
	Stopped    bool
	Destroyed  bool
	CallCounts map[Operation]uint64
}

func (backend *FakeBackend) Snapshot(
	attemptID approvalattempt.AttemptID,
) BackendSnapshot {
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance := backend.instances[attemptID]
	if instance == nil {
		return BackendSnapshot{CallCounts: make(map[Operation]uint64)}
	}
	calls := make(map[Operation]uint64, len(instance.calls))
	for operation, count := range instance.calls {
		calls[operation] = count
	}
	return BackendSnapshot{
		Prepared:   instance.prepared,
		Created:    instance.created,
		Started:    instance.started,
		Observed:   instance.observed,
		Stopped:    instance.stopped,
		Destroyed:  instance.destroyed,
		CallCounts: calls,
	}
}
