package registeredlifecycle

import (
	"context"
	"sync"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

type fakeHandle uint64

type lifecycleRecord struct {
	Snapshot
	handle fakeHandle
}

// MemoryStore is a bounded, single-process development store. Reusing it
// across Component instances exercises recovery, but it makes no durability,
// power-loss, encryption, archival, or multi-process coordination claim.
type MemoryStore struct {
	mu      sync.Mutex
	records map[v0candidate.RegistrationID]lifecycleRecord
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{records: make(map[v0candidate.RegistrationID]lifecycleRecord)}
}

func (store *MemoryStore) begin(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
	planDigest v0candidate.ExecutionPlanDigest,
) error {
	if err := ctx.Err(); err != nil {
		return classified(ClassificationLocalFailure, "lifecycle-context-cancelled")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if store.records == nil {
		store.records = make(map[v0candidate.RegistrationID]lifecycleRecord)
	}
	if _, exists := store.records[registrationID]; exists {
		return classified(ClassificationStale, "lifecycle-registration-already-used")
	}
	if len(store.records) >= MaxLifecycleRecords {
		return classified(ClassificationCapacity, "lifecycle-record-capacity")
	}
	store.records[registrationID] = lifecycleRecord{Snapshot: Snapshot{
		RegistrationID:  registrationID,
		PlanDigest:      planDigest,
		State:           StatePreparing,
		CleanupRequired: true,
		TransitionCount: 1,
	}}
	return nil
}

func (store *MemoryStore) snapshot(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
) (lifecycleRecord, error) {
	if err := ctx.Err(); err != nil {
		return lifecycleRecord{}, classified(ClassificationLocalFailure, "lifecycle-context-cancelled")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	record, exists := store.records[registrationID]
	if !exists {
		return lifecycleRecord{}, classified(ClassificationBinding, "lifecycle-record-not-found")
	}
	return record, nil
}

func (store *MemoryStore) mutate(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
	mutation func(*lifecycleRecord),
) error {
	if err := ctx.Err(); err != nil {
		return classified(ClassificationLocalFailure, "lifecycle-context-cancelled")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	record, exists := store.records[registrationID]
	if !exists {
		return classified(ClassificationBinding, "lifecycle-record-not-found")
	}
	mutation(&record)
	record.TransitionCount++
	store.records[registrationID] = record
	return nil
}

func (store *MemoryStore) transition(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
	state State,
) error {
	return store.mutate(ctx, registrationID, func(record *lifecycleRecord) {
		record.State = state
	})
}

func (store *MemoryStore) retainHandle(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
	handle fakeHandle,
) error {
	return store.mutate(ctx, registrationID, func(record *lifecycleRecord) {
		record.handle = handle
		record.State = StateCreated
	})
}

func (store *MemoryStore) noteFailure(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
	classification Classification,
	operation Operation,
) error {
	return store.mutate(ctx, registrationID, func(record *lifecycleRecord) {
		if record.Failure == "" {
			record.Failure = classification
			record.FailureAt = operation
		}
	})
}

func (store *MemoryStore) markDestroyed(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
) error {
	return store.mutate(ctx, registrationID, func(record *lifecycleRecord) {
		record.State = StateDestroyed
		record.CleanupRequired = false
		record.handle = 0
	})
}

func (store *MemoryStore) markUnresolved(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
) error {
	return store.mutate(ctx, registrationID, func(record *lifecycleRecord) {
		record.State = StateUnresolved
		record.CleanupRequired = true
	})
}

// Snapshot returns one passive defensive copy. No exact plan bytes or backend
// handle are exposed.
func (store *MemoryStore) Snapshot(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
) (Snapshot, error) {
	record, err := store.snapshot(ctx, registrationID)
	if err != nil {
		return Snapshot{}, err
	}
	return record.Snapshot, nil
}
