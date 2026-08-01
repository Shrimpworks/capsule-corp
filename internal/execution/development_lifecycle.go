package execution

import (
	"context"
	"errors"
	"fmt"
)

// DevelopmentLifecycle is deliberately incapable of representing a real guest
// backend. The scaffold uses it to test ordering and recovery only.
type DevelopmentLifecycle interface {
	Name() string
	CreatesGuest() bool
	Prepare(context.Context, AttemptRecord) error
	Create(context.Context, OpaqueID) (string, error)
	Stage(context.Context, string) error
	Start(context.Context, string) error
	Wait(context.Context, string) error
	Collect(context.Context, string) error
	Destroy(context.Context, string) error
	Reconcile(context.Context, OpaqueID) (BackendObservation, error)
}

type BackendObservationStatus string

const (
	BackendPresent               BackendObservationStatus = "present"
	BackendAuthoritativelyAbsent BackendObservationStatus = "authoritatively-absent"
	BackendOutcomeUnknown        BackendObservationStatus = "unknown"
)

type BackendObservation struct {
	Status BackendObservationStatus
	Handle string
}

func (core *SupervisorCore) DriveDevelopmentAttempt(
	ctx context.Context,
	attemptID OpaqueID,
	backend DevelopmentLifecycle,
) error {
	if backend == nil {
		return errors.New("development backend is required")
	}
	if backend.CreatesGuest() {
		return ErrRealBackendBlocked
	}
	attempt, err := core.moveAttempt(ctx, attemptID, AttemptCreated, AttemptPreparing, nil)
	if err != nil {
		return err
	}
	if err := backend.Prepare(ctx, attempt); err != nil {
		return core.finishWithoutBackend(ctx, attemptID, TerminalBackendFailed, err)
	}
	_, err = core.moveAttempt(
		ctx,
		attemptID,
		AttemptPreparing,
		AttemptBackendCreateIntent,
		func(attempt *AttemptRecord) {
			attempt.CleanupRequired = true
			attempt.ReconciliationKey = attempt.ID
		},
	)
	if err != nil {
		return err
	}
	if err := core.checkpoint(CheckpointBackendCreateIntent); err != nil {
		return err
	}

	handle, createErr := backend.Create(ctx, attemptID)
	if createErr != nil || handle == "" {
		if createErr == nil {
			createErr = errors.New("backend returned an empty handle")
		}
		cleanupErr := core.ReconcileDevelopmentAttempt(ctx, attemptID, backend)
		return errors.Join(fmt.Errorf("development backend create: %w", createErr), cleanupErr)
	}
	_, err = core.moveAttempt(
		ctx,
		attemptID,
		AttemptBackendCreateIntent,
		AttemptBackendCreated,
		func(attempt *AttemptRecord) { attempt.BackendHandle = handle },
	)
	if err != nil {
		return err
	}
	if err := core.checkpoint(CheckpointBackendHandle); err != nil {
		return err
	}

	steps := []struct {
		from AttemptState
		to   AttemptState
		call func(context.Context, string) error
		name string
	}{
		{AttemptBackendCreated, AttemptStaging, backend.Stage, "stage"},
		{AttemptStaging, AttemptStarting, backend.Start, "start"},
		{AttemptStarting, AttemptRunning, backend.Wait, "wait"},
		{AttemptRunning, AttemptCollecting, backend.Collect, "collect"},
	}
	for _, step := range steps {
		if _, err := core.moveAttempt(ctx, attemptID, step.from, step.to, nil); err != nil {
			return err
		}
		if err := step.call(ctx, handle); err != nil {
			cleanupErr := core.destroyDevelopmentAttempt(
				ctx,
				attemptID,
				backend,
				TerminalSimulationFailed,
				fmt.Errorf("development backend %s: %w", step.name, err),
			)
			return errors.Join(fmt.Errorf("development backend %s: %w", step.name, err), cleanupErr)
		}
	}
	return core.destroyDevelopmentAttempt(
		ctx,
		attemptID,
		backend,
		TerminalSimulationPassed,
		nil,
	)
}

func (core *SupervisorCore) ReconcileDevelopmentAttempt(
	ctx context.Context,
	attemptID OpaqueID,
	backend DevelopmentLifecycle,
) error {
	if backend == nil {
		return errors.New("development backend is required")
	}
	if backend.CreatesGuest() {
		return ErrRealBackendBlocked
	}
	state, err := core.store.Snapshot(ctx)
	if err != nil {
		return err
	}
	attempt, exists := state.Attempts[attemptID]
	if !exists {
		return ErrNotFound
	}
	if !attempt.CleanupRequired {
		return nil
	}
	if attempt.BackendHandle != "" {
		return core.destroyDevelopmentAttempt(
			ctx,
			attemptID,
			backend,
			TerminalSimulationFailed,
			errors.New("recovered attempt with a durable backend handle"),
		)
	}

	observation, reconcileErr := backend.Reconcile(ctx, attempt.ReconciliationKey)
	if reconcileErr != nil || observation.Status == BackendOutcomeUnknown {
		return core.markCleanupUnresolved(ctx, attemptID, errors.Join(reconcileErr, ErrCleanupUnresolved))
	}
	switch observation.Status {
	case BackendAuthoritativelyAbsent:
		return core.finishWithoutBackend(
			ctx,
			attemptID,
			TerminalBackendFailed,
			errors.New("backend authoritatively absent after create intent"),
		)
	case BackendPresent:
		if observation.Handle == "" {
			return core.markCleanupUnresolved(
				ctx,
				attemptID,
				errors.New("backend present without a durable handle"),
			)
		}
		if _, err := core.setAttemptForReconciledHandle(ctx, attemptID, observation.Handle); err != nil {
			return err
		}
		return core.destroyDevelopmentAttempt(
			ctx,
			attemptID,
			backend,
			TerminalSimulationFailed,
			errors.New("reconciled backend after ambiguous create"),
		)
	default:
		return core.markCleanupUnresolved(
			ctx,
			attemptID,
			fmt.Errorf("unsupported backend observation %q", observation.Status),
		)
	}
}

func (core *SupervisorCore) destroyDevelopmentAttempt(
	ctx context.Context,
	attemptID OpaqueID,
	backend DevelopmentLifecycle,
	terminal TerminalClass,
	cause error,
) error {
	state, err := core.store.Snapshot(ctx)
	if err != nil {
		return err
	}
	attempt, exists := state.Attempts[attemptID]
	if !exists {
		return ErrNotFound
	}
	if attempt.BackendHandle == "" {
		return core.markCleanupUnresolved(ctx, attemptID, errors.New("cleanup has no backend handle"))
	}
	if attempt.State != AttemptDestroying {
		if _, err := core.moveAttemptFromAny(
			ctx,
			attemptID,
			AttemptDestroying,
			func(record *AttemptRecord) { record.Failure = errorText(cause) },
		); err != nil {
			return err
		}
	}
	if err := backend.Destroy(ctx, attempt.BackendHandle); err != nil {
		observation, reconcileErr := backend.Reconcile(ctx, attempt.ReconciliationKey)
		if reconcileErr == nil && observation.Status == BackendAuthoritativelyAbsent {
			return core.finishAttempt(ctx, attemptID, terminal, cause)
		}
		return core.markCleanupUnresolved(
			ctx,
			attemptID,
			errors.Join(fmt.Errorf("destroy backend: %w", err), reconcileErr),
		)
	}
	return core.finishAttempt(ctx, attemptID, terminal, cause)
}

func (core *SupervisorCore) finishWithoutBackend(
	ctx context.Context,
	attemptID OpaqueID,
	terminal TerminalClass,
	cause error,
) error {
	return core.store.Update(ctx, func(state *InstallationState) error {
		attempt, exists := state.Attempts[attemptID]
		if !exists {
			return ErrNotFound
		}
		attempt.State = AttemptTerminal
		attempt.CleanupRequired = false
		attempt.TerminalClass = terminal
		attempt.Failure = errorText(cause)
		state.Attempts[attemptID] = attempt
		return nil
	})
}

func (core *SupervisorCore) finishAttempt(
	ctx context.Context,
	attemptID OpaqueID,
	terminal TerminalClass,
	cause error,
) error {
	return core.finishWithoutBackend(ctx, attemptID, terminal, cause)
}

func (core *SupervisorCore) markCleanupUnresolved(
	ctx context.Context,
	attemptID OpaqueID,
	cause error,
) error {
	updateErr := core.store.Update(ctx, func(state *InstallationState) error {
		attempt, exists := state.Attempts[attemptID]
		if !exists {
			return ErrNotFound
		}
		attempt.State = AttemptUnresolved
		attempt.CleanupRequired = true
		attempt.TerminalClass = TerminalCleanupUnresolved
		attempt.Failure = errorText(cause)
		state.Attempts[attemptID] = attempt
		return nil
	})
	return errors.Join(cause, updateErr, ErrCleanupUnresolved)
}

func (core *SupervisorCore) setAttemptForReconciledHandle(
	ctx context.Context,
	attemptID OpaqueID,
	handle string,
) (AttemptRecord, error) {
	return core.moveAttemptFromAny(
		ctx,
		attemptID,
		AttemptBackendCreated,
		func(attempt *AttemptRecord) { attempt.BackendHandle = handle },
	)
}

func (core *SupervisorCore) moveAttempt(
	ctx context.Context,
	attemptID OpaqueID,
	from AttemptState,
	to AttemptState,
	mutate func(*AttemptRecord),
) (AttemptRecord, error) {
	var result AttemptRecord
	err := core.store.Update(ctx, func(state *InstallationState) error {
		attempt, exists := state.Attempts[attemptID]
		if !exists {
			return ErrNotFound
		}
		if attempt.State != from {
			return ErrInvalidState
		}
		attempt.State = to
		if mutate != nil {
			mutate(&attempt)
		}
		state.Attempts[attemptID] = attempt
		result = attempt
		return nil
	})
	return result, err
}

func (core *SupervisorCore) moveAttemptFromAny(
	ctx context.Context,
	attemptID OpaqueID,
	to AttemptState,
	mutate func(*AttemptRecord),
) (AttemptRecord, error) {
	var result AttemptRecord
	err := core.store.Update(ctx, func(state *InstallationState) error {
		attempt, exists := state.Attempts[attemptID]
		if !exists {
			return ErrNotFound
		}
		if attempt.State == AttemptTerminal {
			return ErrInvalidState
		}
		attempt.State = to
		if mutate != nil {
			mutate(&attempt)
		}
		state.Attempts[attemptID] = attempt
		result = attempt
		return nil
	})
	return result, err
}

func errorText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
