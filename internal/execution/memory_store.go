package execution

import (
	"context"
	"errors"
	"sync"
)

// MemoryStateStore provides atomic copy-on-write transactions for tests and
// development scaffolding. It is not durable and must never back a real attempt.
type MemoryStateStore struct {
	mu    sync.Mutex
	state InstallationState
}

func NewMemoryStateStore(
	installationID OpaqueID,
	supervisorID OpaqueID,
	epochDigest Digest,
) (*MemoryStateStore, error) {
	if installationID.IsZero() || supervisorID.IsZero() || epochDigest.IsZero() {
		return nil, errors.New("installation, supervisor, and epoch identities are required")
	}
	return &MemoryStateStore{
		state: InstallationState{
			InstallationID:  installationID,
			SupervisorID:    supervisorID,
			EpochDigest:     epochDigest,
			TrustPhase:      TrustStable,
			AttemptsEnabled: true,
			Registrations:   make(map[OpaqueID]PlanRegistration),
			Approvals:       make(map[OpaqueID]ApprovalRecord),
			Attempts:        make(map[OpaqueID]AttemptRecord),
		},
	}, nil
}

func (store *MemoryStateStore) Snapshot(ctx context.Context) (InstallationState, error) {
	if err := ctx.Err(); err != nil {
		return InstallationState{}, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	return cloneInstallationState(store.state), nil
}

func (store *MemoryStateStore) Update(
	ctx context.Context,
	update func(*InstallationState) error,
) error {
	if update == nil {
		return errors.New("state update function is required")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	store.mu.Lock()
	defer store.mu.Unlock()

	candidate := cloneInstallationState(store.state)
	if err := update(&candidate); err != nil {
		return err
	}
	if err := validateInstallationState(candidate); err != nil {
		return err
	}
	store.state = candidate
	return nil
}

func cloneInstallationState(state InstallationState) InstallationState {
	clone := state
	if state.Transition != nil {
		transition := *state.Transition
		transition.AcceptedIncarnations = make(map[ComponentRole]OpaqueID)
		for role, incarnation := range state.Transition.AcceptedIncarnations {
			transition.AcceptedIncarnations[role] = incarnation
		}
		clone.Transition = &transition
	}
	clone.Registrations = make(map[OpaqueID]PlanRegistration, len(state.Registrations))
	for id, registration := range state.Registrations {
		clone.Registrations[id] = cloneRegistration(registration)
	}
	clone.Approvals = make(map[OpaqueID]ApprovalRecord, len(state.Approvals))
	for id, approval := range state.Approvals {
		clone.Approvals[id] = cloneApproval(approval)
	}
	clone.Attempts = make(map[OpaqueID]AttemptRecord, len(state.Attempts))
	for id, attempt := range state.Attempts {
		clone.Attempts[id] = attempt
	}
	return clone
}

func validateInstallationState(state InstallationState) error {
	if state.InstallationID.IsZero() || state.SupervisorID.IsZero() || state.EpochDigest.IsZero() {
		return errors.New("state lost its installation trust identity")
	}
	if state.Registrations == nil || state.Approvals == nil || state.Attempts == nil {
		return errors.New("state record maps must be initialized")
	}
	if state.AttemptsEnabled != (state.TrustPhase == TrustStable) {
		return errors.New("attempt enablement must exactly match stable trust state")
	}
	if state.TrustPhase != TrustStable && state.Transition == nil {
		return errors.New("non-stable trust state requires a transition record")
	}
	for id, approval := range state.Approvals {
		if id != approval.ID || approval.ID.IsZero() || approval.RegistrationID.IsZero() {
			return errors.New("approval map identity mismatch")
		}
		if approval.State == ApprovalConsumed && approval.AttemptID.IsZero() {
			return errors.New("consumed approval is not bound to an attempt")
		}
	}
	for id, attempt := range state.Attempts {
		if id != attempt.ID || attempt.ID.IsZero() || attempt.RegistrationID.IsZero() ||
			attempt.ApprovalID.IsZero() {
			return errors.New("attempt map identity mismatch")
		}
		if attempt.State == AttemptTerminal && attempt.CleanupRequired {
			return errors.New("terminal attempt retains a cleanup obligation")
		}
	}
	return nil
}
