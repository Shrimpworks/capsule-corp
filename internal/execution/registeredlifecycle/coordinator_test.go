package registeredlifecycle

import (
	"context"
	"sync"
	"testing"

	"capsule.local/capsule/internal/execution/lifecyclestate"
)

func mustCoordinator(t *testing.T) *Coordinator {
	t.Helper()
	owner := lifecyclestate.OwnerSessionID{0x01}
	coordinator, err := NewCoordinator(owner)
	if err != nil {
		t.Fatalf("construct coordinator: %v", err)
	}
	return coordinator
}

// TestCoordinatorEvictsTerminalAttemptsBeforeExhaustingCapacity proves issue
// #140's fix: driving far more than maxCoordinatedAttempts distinct AttemptIDs
// to the terminal Destroyed state must not permanently exhaust the
// coordinator's capacity, because each entry is pruned once its work
// completes at Destroyed with no other caller still holding it.
func TestCoordinatorEvictsTerminalAttemptsBeforeExhaustingCapacity(t *testing.T) {
	coordinator := mustCoordinator(t)
	const driven = maxCoordinatedAttempts + 50

	for index := uint64(0); index < driven; index++ {
		attemptID := attemptIDFor(index)
		snapshot, err := coordinator.withAttempt(context.Background(), attemptID, func() (Snapshot, error) {
			return Snapshot{AttemptID: attemptID, State: lifecyclestate.StateDestroyed}, nil
		})
		if err != nil {
			t.Fatalf("attempt %d: unexpected coordinator error: %v", index, err)
		}
		if snapshot.State != lifecyclestate.StateDestroyed {
			t.Fatalf("attempt %d: snapshot state = %q, want Destroyed", index, snapshot.State)
		}
	}

	coordinator.mu.Lock()
	remaining := len(coordinator.locks)
	coordinator.mu.Unlock()
	if remaining != 0 {
		t.Fatalf("coordinator retained %d lock entries after every attempt reached Destroyed, want 0", remaining)
	}
}

// TestCoordinatorStillEnforcesCapacityForNonterminalWork proves the existing
// capacity ceiling still applies: entries that never reach Destroyed are not
// evicted, so the (maxCoordinatedAttempts+1)-th distinct nonterminal attempt
// is refused.
func TestCoordinatorStillEnforcesCapacityForNonterminalWork(t *testing.T) {
	coordinator := mustCoordinator(t)

	for index := uint64(0); index < maxCoordinatedAttempts; index++ {
		attemptID := attemptIDFor(index)
		_, err := coordinator.withAttempt(context.Background(), attemptID, func() (Snapshot, error) {
			return Snapshot{AttemptID: attemptID, State: lifecyclestate.StateCreated}, nil
		})
		if err != nil {
			t.Fatalf("attempt %d: unexpected coordinator error: %v", index, err)
		}
	}

	coordinator.mu.Lock()
	filled := len(coordinator.locks)
	coordinator.mu.Unlock()
	if filled != maxCoordinatedAttempts {
		t.Fatalf("coordinator holds %d entries after filling nonterminal work, want %d", filled, maxCoordinatedAttempts)
	}

	overflow := attemptIDFor(maxCoordinatedAttempts)
	_, err := coordinator.withAttempt(context.Background(), overflow, func() (Snapshot, error) {
		t.Fatal("work must not run once the coordinator is at capacity")
		return Snapshot{}, nil
	})
	classification, ok := ErrorClassification(err)
	if !ok || classification != ClassificationCapacity {
		t.Fatalf("classification = %v (recognized %v), want CAPACITY: %v", classification, ok, err)
	}
}

// TestCoordinatorSerializesSameAttemptIDAndEvictsExactlyOnce drives many
// concurrent callers against one AttemptID, proving mutual exclusion still
// holds under the refcounted eviction and that the final terminal outcome
// prunes the entry exactly once with no double-delete or leaked entry.
func TestCoordinatorSerializesSameAttemptIDAndEvictsExactlyOnce(t *testing.T) {
	coordinator := mustCoordinator(t)
	attemptID := attemptIDFor(1)

	const callers = 64
	var active int32
	var mu sync.Mutex
	var maxObservedActive int32
	var wg sync.WaitGroup
	wg.Add(callers)
	for i := 0; i < callers; i++ {
		go func(final bool) {
			defer wg.Done()
			_, _ = coordinator.withAttempt(context.Background(), attemptID, func() (Snapshot, error) {
				mu.Lock()
				active++
				if active > maxObservedActive {
					maxObservedActive = active
				}
				mu.Unlock()
				defer func() {
					mu.Lock()
					active--
					mu.Unlock()
				}()
				return Snapshot{AttemptID: attemptID, State: lifecyclestate.StateCreated}, nil
			})
		}(i == callers-1)
	}
	wg.Wait()

	if maxObservedActive > 1 {
		t.Fatalf("observed %d concurrently active callers for one AttemptID, want at most 1", maxObservedActive)
	}

	// Nothing reached Destroyed yet, so the entry must still be retained.
	coordinator.mu.Lock()
	beforeTerminal := len(coordinator.locks)
	coordinator.mu.Unlock()
	if beforeTerminal != 1 {
		t.Fatalf("coordinator holds %d entries before the terminal call, want 1", beforeTerminal)
	}

	if _, err := coordinator.withAttempt(context.Background(), attemptID, func() (Snapshot, error) {
		return Snapshot{AttemptID: attemptID, State: lifecyclestate.StateDestroyed}, nil
	}); err != nil {
		t.Fatalf("terminal call: unexpected error: %v", err)
	}

	coordinator.mu.Lock()
	afterTerminal := len(coordinator.locks)
	coordinator.mu.Unlock()
	if afterTerminal != 0 {
		t.Fatalf("coordinator retained %d entries after the sole caller reached Destroyed, want 0", afterTerminal)
	}
}
