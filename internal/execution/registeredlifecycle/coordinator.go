package registeredlifecycle

import (
	"context"
	"errors"
	"sync"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
)

const maxCoordinatedAttempts = 4_096

// Coordinator serializes every Drive, Recover, and startup action for one
// AttemptID under the same sealed owner session as the open durable store.
// It is an in-process E4 mechanic, not the future platform owner lock.
type Coordinator struct {
	owner lifecyclestate.OwnerSessionID
	mu    sync.Mutex
	locks map[approvalattempt.AttemptID]*sync.Mutex
}

func NewCoordinator(owner lifecyclestate.OwnerSessionID) (*Coordinator, error) {
	if owner.IsZero() {
		return nil, errors.New("lifecycle coordinator owner session is required")
	}
	return &Coordinator{
		owner: owner,
		locks: make(map[approvalattempt.AttemptID]*sync.Mutex),
	}, nil
}

func (coordinator *Coordinator) OwnerSessionID() lifecyclestate.OwnerSessionID {
	if coordinator == nil {
		return lifecyclestate.OwnerSessionID{}
	}
	return coordinator.owner
}

func (coordinator *Coordinator) withAttempt(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	work func() (Snapshot, error),
) (Snapshot, error) {
	if err := ctx.Err(); err != nil {
		return Snapshot{}, classified(ClassificationLocalFailure, "lifecycle-context-cancelled")
	}
	coordinator.mu.Lock()
	lock := coordinator.locks[attemptID]
	if lock == nil {
		if len(coordinator.locks) >= maxCoordinatedAttempts {
			coordinator.mu.Unlock()
			return Snapshot{}, classified(ClassificationCapacity, "lifecycle-coordinator-capacity")
		}
		lock = &sync.Mutex{}
		coordinator.locks[attemptID] = lock
	}
	coordinator.mu.Unlock()
	lock.Lock()
	defer lock.Unlock()
	if err := ctx.Err(); err != nil {
		return Snapshot{}, classified(ClassificationLocalFailure, "lifecycle-context-cancelled")
	}
	return work()
}
