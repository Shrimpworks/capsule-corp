package registeredlifecycle

import (
	"context"
	"errors"
	"sync"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
)

const maxCoordinatedAttempts = 4_096

// coordinatedLock is one AttemptID's serialization mutex plus the count of
// callers currently holding a reference to it. An entry is only removed from
// Coordinator.locks while refCount is zero, so a caller that already looked
// the entry up always finishes on the same mutex it acquired.
type coordinatedLock struct {
	mu       sync.Mutex
	refCount int
}

// Coordinator serializes every Drive, Recover, and startup action for one
// AttemptID under the same sealed owner session as the open durable store.
// It is an in-process E4 mechanic, not the future platform owner lock.
//
// An AttemptID's entry is evicted once its lifecycle reaches the terminal
// Destroyed state and no caller is still using its lock, so a long-running
// process does not permanently exhaust maxCoordinatedAttempts on attempts
// that no longer need coordination.
type Coordinator struct {
	owner lifecyclestate.OwnerSessionID
	mu    sync.Mutex
	locks map[approvalattempt.AttemptID]*coordinatedLock
}

func NewCoordinator(owner lifecyclestate.OwnerSessionID) (*Coordinator, error) {
	if owner.IsZero() {
		return nil, errors.New("lifecycle coordinator owner session is required")
	}
	return &Coordinator{
		owner: owner,
		locks: make(map[approvalattempt.AttemptID]*coordinatedLock),
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
	entry := coordinator.locks[attemptID]
	if entry == nil {
		if len(coordinator.locks) >= maxCoordinatedAttempts {
			coordinator.mu.Unlock()
			return Snapshot{}, classified(ClassificationCapacity, "lifecycle-coordinator-capacity")
		}
		entry = &coordinatedLock{}
		coordinator.locks[attemptID] = entry
	}
	entry.refCount++
	coordinator.mu.Unlock()

	entry.mu.Lock()
	snapshot, err := func() (Snapshot, error) {
		defer entry.mu.Unlock()
		if ctxErr := ctx.Err(); ctxErr != nil {
			return Snapshot{}, classified(ClassificationLocalFailure, "lifecycle-context-cancelled")
		}
		return work()
	}()

	coordinator.mu.Lock()
	entry.refCount--
	if entry.refCount == 0 && snapshot.State == lifecyclestate.StateDestroyed {
		delete(coordinator.locks, attemptID)
	}
	coordinator.mu.Unlock()

	return snapshot, err
}
