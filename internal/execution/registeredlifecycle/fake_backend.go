package registeredlifecycle

import (
	"context"
	"errors"
	"sync"

	"capsule.local/capsule/internal/protocol/v0candidate"
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
	registrationID v0candidate.RegistrationID
	operation      Operation
	moment         FaultMoment
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
	instances map[v0candidate.RegistrationID]*fakeInstance
	faults    map[fakeFaultKey]uint64
}

func NewFakeBackend() *FakeBackend {
	return &FakeBackend{
		next:      1,
		instances: make(map[v0candidate.RegistrationID]*fakeInstance),
		faults:    make(map[fakeFaultKey]uint64),
	}
}

func (*FakeBackend) CreatesGuest() bool { return false }

// InjectFault adds one deterministic, registration-scoped fault. The fault is
// consumed immediately before or immediately after the selected fake effect.
func (backend *FakeBackend) InjectFault(
	registrationID v0candidate.RegistrationID,
	operation Operation,
	moment FaultMoment,
) error {
	if registrationID == (v0candidate.RegistrationID{}) {
		return errors.New("fake fault registration ID is required")
	}
	if !validOperation(operation) || (moment != FaultBeforeEffect && moment != FaultAfterEffect) {
		return errors.New("fake fault operation and moment must be closed values")
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	if backend.faults == nil {
		backend.faults = make(map[fakeFaultKey]uint64)
	}
	key := fakeFaultKey{registrationID: registrationID, operation: operation, moment: moment}
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
	registrationID v0candidate.RegistrationID,
	operation Operation,
	moment FaultMoment,
) bool {
	key := fakeFaultKey{registrationID: registrationID, operation: operation, moment: moment}
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
	registrationID v0candidate.RegistrationID,
) *fakeInstance {
	if backend.instances == nil {
		backend.instances = make(map[v0candidate.RegistrationID]*fakeInstance)
	}
	if backend.next == 0 {
		backend.next = 1
	}
	instance := backend.instances[registrationID]
	if instance == nil {
		instance = &fakeInstance{calls: make(map[Operation]uint64)}
		backend.instances[registrationID] = instance
	}
	return instance
}

func (backend *FakeBackend) beforeLocked(
	registrationID v0candidate.RegistrationID,
	operation Operation,
) (*fakeInstance, error) {
	instance := backend.instanceLocked(registrationID)
	instance.calls[operation]++
	if backend.takeFaultLocked(registrationID, operation, FaultBeforeEffect) {
		return instance, errInjectedFakeFault
	}
	return instance, nil
}

func (backend *FakeBackend) afterLocked(
	registrationID v0candidate.RegistrationID,
	operation Operation,
) error {
	if backend.takeFaultLocked(registrationID, operation, FaultAfterEffect) {
		return errInjectedFakeFault
	}
	return nil
}

func (backend *FakeBackend) prepare(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(registrationID, OperationPrepare)
	if err != nil {
		return err
	}
	instance.prepared = true
	return backend.afterLocked(registrationID, OperationPrepare)
}

func (backend *FakeBackend) create(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
) (fakeHandle, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(registrationID, OperationCreate)
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
	if err := backend.afterLocked(registrationID, OperationCreate); err != nil {
		return 0, err
	}
	return instance.handle, nil
}

func (backend *FakeBackend) start(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
	handle fakeHandle,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(registrationID, OperationStart)
	if err != nil {
		return err
	}
	if !instance.created || instance.destroyed || instance.handle != handle {
		return errors.New("fake start state rejected")
	}
	instance.started = true
	return backend.afterLocked(registrationID, OperationStart)
}

func (backend *FakeBackend) observe(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
	handle fakeHandle,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(registrationID, OperationObserve)
	if err != nil {
		return err
	}
	if !instance.started || instance.destroyed || instance.handle != handle {
		return errors.New("fake observe state rejected")
	}
	instance.observed = true
	return backend.afterLocked(registrationID, OperationObserve)
}

func (backend *FakeBackend) stop(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
	handle fakeHandle,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(registrationID, OperationStop)
	if err != nil {
		return err
	}
	if !instance.created || instance.destroyed || instance.handle != handle {
		return errors.New("fake stop state rejected")
	}
	instance.stopped = true
	return backend.afterLocked(registrationID, OperationStop)
}

func (backend *FakeBackend) destroy(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
	handle fakeHandle,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(registrationID, OperationDestroy)
	if err != nil {
		return err
	}
	if !instance.created || instance.destroyed || instance.handle != handle {
		return errors.New("fake destroy state rejected")
	}
	instance.destroyed = true
	return backend.afterLocked(registrationID, OperationDestroy)
}

func (backend *FakeBackend) reconcile(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
) (fakeObservation, error) {
	if err := ctx.Err(); err != nil {
		return fakeObservation{status: fakeUnknown}, err
	}
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance, err := backend.beforeLocked(registrationID, OperationReconcile)
	if err != nil {
		return fakeObservation{status: fakeUnknown}, err
	}
	observation := fakeObservation{status: fakeAbsent}
	if instance.created && !instance.destroyed {
		observation = fakeObservation{status: fakePresent, handle: instance.handle}
	}
	if err := backend.afterLocked(registrationID, OperationReconcile); err != nil {
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
	registrationID v0candidate.RegistrationID,
) BackendSnapshot {
	backend.mu.Lock()
	defer backend.mu.Unlock()
	instance := backend.instances[registrationID]
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
