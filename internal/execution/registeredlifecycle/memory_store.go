package registeredlifecycle

import (
	"context"
	"sync"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/registrationstate"
)

type fakeHandle uint64

type lifecycleRecord struct {
	Snapshot
	created registrationstate.CreatedAttempt
	handle  fakeHandle
}

// MemoryStore is a bounded, single-process development store. Reusing it
// across Component instances exercises recovery, but it makes no durability,
// power-loss, encryption, archival, or multi-process coordination claim.
type MemoryStore struct {
	mu      sync.Mutex
	records map[approvalattempt.AttemptID]lifecycleRecord
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{records: make(map[approvalattempt.AttemptID]lifecycleRecord)}
}

func (store *MemoryStore) begin(
	ctx context.Context,
	created registrationstate.CreatedAttempt,
) (lifecycleRecord, bool, error) {
	if err := ctx.Err(); err != nil {
		return lifecycleRecord{}, false, classified(ClassificationLocalFailure, "lifecycle-context-cancelled")
	}
	created = cloneCreatedAttempt(created)
	attempt := created.Attempt
	store.mu.Lock()
	defer store.mu.Unlock()
	if store.records == nil {
		store.records = make(map[approvalattempt.AttemptID]lifecycleRecord)
	}
	if existing, exists := store.records[attempt.AttemptID]; exists {
		if !createdAttemptsEqual(existing.created, created) {
			return lifecycleRecord{}, false, classified(ClassificationBinding, "lifecycle-attempt-replay-mismatch")
		}
		return cloneLifecycleRecord(existing), false, nil
	}
	for _, existing := range store.records {
		if existing.ApprovalID == attempt.ApprovalID {
			return lifecycleRecord{}, false, classified(ClassificationBinding, "lifecycle-approval-already-owned")
		}
	}
	if len(store.records) >= MaxLifecycleRecords {
		return lifecycleRecord{}, false, classified(ClassificationCapacity, "lifecycle-record-capacity")
	}
	record := lifecycleRecord{Snapshot: Snapshot{
		AttemptID:             attempt.AttemptID,
		ApprovalID:            attempt.ApprovalID,
		AttemptNonce:          attempt.AttemptNonce,
		RegistrationID:        attempt.RegistrationID,
		RegistrationSequence:  attempt.RegistrationSequence,
		PlanDigest:            attempt.PlanDigest,
		InstallationID:        attempt.InstallationID,
		EpochSequence:         attempt.EpochSequence,
		EpochDigest:           attempt.EpochDigest,
		SupervisorID:          attempt.SupervisorID,
		ApprovalPurpose:       attempt.Purpose,
		ApprovalAudience:      attempt.Audience,
		ApprovalPayloadDigest: attempt.ApprovalPayloadDigest,
		AuthorizationIdentity: attempt.AuthorizationIdentity,
		CreatedAt:             attempt.CreatedAt,
		State:                 StatePreparing,
		CleanupRequired:       true,
		TransitionCount:       1,
	}, created: created}
	store.records[attempt.AttemptID] = record
	return cloneLifecycleRecord(record), true, nil
}

func (store *MemoryStore) snapshot(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (lifecycleRecord, error) {
	record, exists, err := store.lookup(ctx, attemptID)
	if err != nil {
		return lifecycleRecord{}, err
	}
	if !exists {
		return lifecycleRecord{}, classified(ClassificationBinding, "lifecycle-record-not-found")
	}
	return record, nil
}

func (store *MemoryStore) lookup(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (lifecycleRecord, bool, error) {
	if err := ctx.Err(); err != nil {
		return lifecycleRecord{}, false, classified(ClassificationLocalFailure, "lifecycle-context-cancelled")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	record, exists := store.records[attemptID]
	if !exists {
		return lifecycleRecord{}, false, nil
	}
	return cloneLifecycleRecord(record), true, nil
}

func (store *MemoryStore) mutate(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	mutation func(*lifecycleRecord),
) error {
	if err := ctx.Err(); err != nil {
		return classified(ClassificationLocalFailure, "lifecycle-context-cancelled")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	record, exists := store.records[attemptID]
	if !exists {
		return classified(ClassificationBinding, "lifecycle-record-not-found")
	}
	mutation(&record)
	record.TransitionCount++
	store.records[attemptID] = record
	return nil
}

func (store *MemoryStore) transition(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	state State,
) error {
	return store.mutate(ctx, attemptID, func(record *lifecycleRecord) {
		record.State = state
	})
}

func (store *MemoryStore) retainHandle(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	handle fakeHandle,
) error {
	return store.mutate(ctx, attemptID, func(record *lifecycleRecord) {
		record.handle = handle
		record.State = StateCreated
	})
}

func (store *MemoryStore) noteFailure(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	classification Classification,
	operation Operation,
) error {
	return store.mutate(ctx, attemptID, func(record *lifecycleRecord) {
		if record.Failure == "" {
			record.Failure = classification
			record.FailureAt = operation
		}
	})
}

func (store *MemoryStore) markDestroyed(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) error {
	return store.mutate(ctx, attemptID, func(record *lifecycleRecord) {
		record.State = StateDestroyed
		record.CleanupRequired = false
		record.handle = 0
	})
}

func (store *MemoryStore) markUnresolved(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) error {
	return store.mutate(ctx, attemptID, func(record *lifecycleRecord) {
		record.State = StateUnresolved
		record.CleanupRequired = true
	})
}

// Snapshot returns one passive defensive copy. No exact plan bytes or backend
// handle are exposed.
func (store *MemoryStore) Snapshot(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (Snapshot, error) {
	record, err := store.snapshot(ctx, attemptID)
	if err != nil {
		return Snapshot{}, err
	}
	return record.Snapshot, nil
}

func cloneLifecycleRecord(record lifecycleRecord) lifecycleRecord {
	record.created = cloneCreatedAttempt(record.created)
	return record
}
