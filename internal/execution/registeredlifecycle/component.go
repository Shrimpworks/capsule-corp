package registeredlifecycle

import (
	"context"
	"crypto/sha256"
	"errors"

	"capsule.local/capsule/internal/execution/registrationstate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// Component drives only the concrete, no-guest FakeBackend. It is not wired
// to SupervisorCore, an endpoint, approval, attempts, content, or evidence.
type Component struct {
	registrations registrationstate.RegistrationResolver
	roleBindings  PlanRoleBindingResolver
	store         *MemoryStore
	backend       *FakeBackend
	checkpoint    CheckpointHook
}

func New(options Options) (*Component, error) {
	if options.Registrations == nil || options.RoleBindings == nil ||
		options.Store == nil || options.Backend == nil {
		return nil, errors.New("registration resolver, role bindings, store, and fake backend are required")
	}
	if options.Backend.CreatesGuest() {
		return nil, errors.New("registered lifecycle requires a no-guest fake backend")
	}
	if options.Checkpoint == nil {
		options.Checkpoint = func(context.Context, Checkpoint, v0candidate.RegistrationID) error {
			return nil
		}
	}
	return &Component{
		registrations: options.Registrations,
		roleBindings:  options.RoleBindings,
		store:         options.Store,
		backend:       options.Backend,
		checkpoint:    options.Checkpoint,
	}, nil
}

// Execute accepts only a Supervisor-issued RegistrationID. It resolves and
// revalidates the Supervisor-retained exact bytes using separately trusted
// role bindings before creating any lifecycle state. It returns lifecycle
// disposition only; there is no caller-supplied or backend-supplied job result
// and no ordinary terminal-success classification.
func (component *Component) Execute(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
) (Snapshot, error) {
	if registrationID == (v0candidate.RegistrationID{}) {
		return Snapshot{}, classified(ClassificationBinding, "registration-id-required")
	}
	record, err := component.registrations.ResolveUsable(ctx, registrationID)
	if err != nil {
		return Snapshot{}, mapRegistrationFailure(err)
	}
	bindings, err := component.roleBindings.ResolveExecutionPlanRoleBindings(ctx, registrationID)
	if err != nil {
		return Snapshot{}, classified(ClassificationBinding, "trusted-plan-bindings-unavailable")
	}
	decoded, err := v0candidate.DecodeExecutionPlan(
		record.ExactPlanBytes,
		cloneRoleBindings(bindings),
	)
	if err != nil {
		return Snapshot{}, mapDecodeFailure(err)
	}
	computed := v0candidate.ExecutionPlanDigest(sha256.Sum256(record.ExactPlanBytes))
	if computed != record.RecomputedPlanDigest || decoded.Digest() != record.RecomputedPlanDigest {
		return Snapshot{}, classified(ClassificationBinding, "stored-plan-digest-mismatch")
	}
	if err := component.store.begin(ctx, registrationID, decoded.Digest()); err != nil {
		return Snapshot{}, err
	}

	if err := component.backend.prepare(ctx, registrationID); err != nil {
		return component.failAndCleanup(ctx, registrationID, OperationPrepare)
	}
	if err := component.afterEffect(
		ctx,
		CheckpointAfterPrepareEffect,
		registrationID,
		OperationPrepare,
	); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	if err := component.store.transition(ctx, registrationID, StateCreateIntent); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	handle, err := component.backend.create(ctx, registrationID)
	if err != nil || handle == 0 {
		return component.failAndCleanup(ctx, registrationID, OperationCreate)
	}
	if err := component.afterEffect(
		ctx,
		CheckpointAfterCreateEffect,
		registrationID,
		OperationCreate,
	); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	if err := component.store.retainHandle(ctx, registrationID, handle); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	if err := component.store.transition(ctx, registrationID, StateStarting); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	if err := component.backend.start(ctx, registrationID, handle); err != nil {
		return component.failAndCleanup(ctx, registrationID, OperationStart)
	}
	if err := component.afterEffect(
		ctx,
		CheckpointAfterStartEffect,
		registrationID,
		OperationStart,
	); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	if err := component.store.transition(ctx, registrationID, StateObserving); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	if err := component.backend.observe(ctx, registrationID, handle); err != nil {
		return component.failAndCleanup(ctx, registrationID, OperationObserve)
	}
	if err := component.afterEffect(
		ctx,
		CheckpointAfterObserveEffect,
		registrationID,
		OperationObserve,
	); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	return component.cleanup(ctx, registrationID)
}

// Recover accepts only the original registration ID. It never re-resolves an
// expired registration before cleanup: an existing cleanup obligation must
// survive registration expiry and be reconciled from retained lifecycle state.
func (component *Component) Recover(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
) (Snapshot, error) {
	record, err := component.store.snapshot(ctx, registrationID)
	if err != nil {
		return Snapshot{}, err
	}
	if record.State == StateDestroyed {
		return record.Snapshot, nil
	}
	if !record.CleanupRequired {
		return record.Snapshot, classified(ClassificationLocalFailure, "nonterminal-cleanup-state-invalid")
	}
	return component.cleanup(ctx, registrationID)
}

func (component *Component) failAndCleanup(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
	operation Operation,
) (Snapshot, error) {
	if err := component.store.noteFailure(
		ctx,
		registrationID,
		ClassificationLifecycleFailure,
		operation,
	); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	snapshot, cleanupErr := component.cleanup(ctx, registrationID)
	if cleanupErr != nil {
		return snapshot, cleanupErr
	}
	return snapshot, classified(ClassificationLifecycleFailure, "fake-lifecycle-operation-failed")
}

func (component *Component) cleanup(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
) (Snapshot, error) {
	record, err := component.store.snapshot(ctx, registrationID)
	if err != nil {
		return Snapshot{}, err
	}
	if record.State == StateDestroyed {
		return record.Snapshot, nil
	}
	if record.State == StateDestroying {
		observation, reconcileErr := component.backend.reconcile(ctx, registrationID)
		if reconcileErr != nil || observation.status == fakeUnknown {
			return component.unresolved(ctx, registrationID)
		}
		if observation.status == fakeAbsent {
			if err := component.store.markDestroyed(ctx, registrationID); err != nil {
				return component.snapshotWithError(ctx, registrationID, err)
			}
			return component.store.Snapshot(ctx, registrationID)
		}
	}
	if record.handle == 0 {
		observation, reconcileErr := component.backend.reconcile(ctx, registrationID)
		if reconcileErr != nil || observation.status == fakeUnknown {
			return component.unresolved(ctx, registrationID)
		}
		switch observation.status {
		case fakeAbsent:
			if err := component.store.markDestroyed(ctx, registrationID); err != nil {
				return component.snapshotWithError(ctx, registrationID, err)
			}
			return component.store.Snapshot(ctx, registrationID)
		case fakePresent:
			if observation.handle == 0 {
				return component.unresolved(ctx, registrationID)
			}
			if err := component.store.retainHandle(ctx, registrationID, observation.handle); err != nil {
				return component.snapshotWithError(ctx, registrationID, err)
			}
			record.handle = observation.handle
		default:
			return component.unresolved(ctx, registrationID)
		}
	}

	if err := component.store.transition(ctx, registrationID, StateStopping); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	stopErr := component.backend.stop(ctx, registrationID, record.handle)
	if stopErr != nil {
		_ = component.store.noteFailure(
			ctx,
			registrationID,
			ClassificationLifecycleFailure,
			OperationStop,
		)
	} else if err := component.afterEffect(
		ctx,
		CheckpointAfterStopEffect,
		registrationID,
		OperationStop,
	); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	if err := component.store.transition(ctx, registrationID, StateDestroying); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	destroyErr := component.backend.destroy(ctx, registrationID, record.handle)
	if destroyErr != nil {
		_ = component.store.noteFailure(
			ctx,
			registrationID,
			ClassificationLifecycleFailure,
			OperationDestroy,
		)
	} else if err := component.afterEffect(
		ctx,
		CheckpointAfterDestroyEffect,
		registrationID,
		OperationDestroy,
	); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}

	observation, reconcileErr := component.backend.reconcile(ctx, registrationID)
	if reconcileErr != nil || observation.status != fakeAbsent {
		return component.unresolved(ctx, registrationID)
	}
	if err := component.store.markDestroyed(ctx, registrationID); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	snapshot, err := component.store.Snapshot(ctx, registrationID)
	if err != nil {
		return Snapshot{}, err
	}
	if stopErr != nil || destroyErr != nil {
		return snapshot, classified(ClassificationLifecycleFailure, "fake-lifecycle-cleanup-operation-failed")
	}
	return snapshot, nil
}

func (component *Component) unresolved(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
) (Snapshot, error) {
	if err := component.store.markUnresolved(ctx, registrationID); err != nil {
		return component.snapshotWithError(ctx, registrationID, err)
	}
	snapshot, err := component.store.Snapshot(ctx, registrationID)
	if err != nil {
		return Snapshot{}, err
	}
	return snapshot, classified(ClassificationCleanupUnresolved, "fake-lifecycle-cleanup-unresolved")
}

func (component *Component) afterEffect(
	ctx context.Context,
	checkpoint Checkpoint,
	registrationID v0candidate.RegistrationID,
	operation Operation,
) error {
	if err := component.checkpoint(ctx, checkpoint, registrationID); err != nil {
		_ = component.store.noteFailure(
			ctx,
			registrationID,
			ClassificationLocalFailure,
			operation,
		)
		return classified(ClassificationLocalFailure, "simulated-lifecycle-interruption")
	}
	return nil
}

func (component *Component) snapshotWithError(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
	err error,
) (Snapshot, error) {
	snapshot, snapshotErr := component.store.Snapshot(ctx, registrationID)
	if snapshotErr != nil {
		return Snapshot{}, err
	}
	return snapshot, err
}

func mapDecodeFailure(err error) error {
	classification, ok := v0candidate.ErrorClassification(err)
	if !ok {
		return classified(ClassificationLocalFailure, "stored-plan-decode-failed")
	}
	switch classification {
	case v0candidate.ClassificationMalformed:
		return classified(ClassificationMalformed, "stored-plan-rejected")
	case v0candidate.ClassificationUnsupported:
		return classified(ClassificationUnsupported, "stored-plan-rejected")
	case v0candidate.ClassificationSchema:
		return classified(ClassificationSchema, "stored-plan-rejected")
	case v0candidate.ClassificationDomain:
		return classified(ClassificationDomain, "stored-plan-rejected")
	default:
		return classified(ClassificationLocalFailure, "stored-plan-decode-failed")
	}
}
