package registeredlifecycle

import (
	"context"
	"crypto/sha256"
	"errors"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/registrationstate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// Component drives only durable created attempts from the Slice B store into
// the concrete, no-guest FakeBackend. It is not wired to a deployed
// Supervisor, an endpoint, content, evidence, a real backend, or a guest.
type Component struct {
	attempts        AttemptResolver
	createdAttempts CreatedAttemptEnumerator
	store           *MemoryStore
	backend         *FakeBackend
	checkpoint      CheckpointHook
}

func New(options Options) (*Component, error) {
	if options.Attempts == nil || options.CreatedAttempts == nil ||
		options.Store == nil || options.Backend == nil {
		return nil, errors.New("attempt resolver, created-attempt enumeration, store, and fake backend are required")
	}
	if options.Backend.CreatesGuest() {
		return nil, errors.New("registered lifecycle requires a no-guest fake backend")
	}
	if options.Checkpoint == nil {
		options.Checkpoint = func(context.Context, Checkpoint, approvalattempt.AttemptID) error {
			return nil
		}
	}
	return &Component{
		attempts:        options.Attempts,
		createdAttempts: options.CreatedAttempts,
		store:           options.Store,
		backend:         options.Backend,
		checkpoint:      options.Checkpoint,
	}, nil
}

// Drive accepts only a Supervisor-issued AttemptID. It resolves and
// independently revalidates the already committed created attempt and exact
// registered plan before creating lifecycle state or calling fake prepare.
// An existing identical lifecycle record is returned without redriving any
// effect.
func (component *Component) Drive(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (Snapshot, error) {
	if attemptID == (approvalattempt.AttemptID{}) {
		return Snapshot{}, classified(ClassificationBinding, "attempt-id-required")
	}
	created, err := component.attempts.ResolveCreated(ctx, attemptID)
	if err != nil {
		return Snapshot{}, mapAttemptResolutionFailure(err)
	}
	created = cloneCreatedAttempt(created)
	if err := validateCreatedAttempt(created, attemptID); err != nil {
		return Snapshot{}, err
	}
	decoded, err := v0candidate.DecodeExecutionPlan(
		created.Registration.ExactPlanBytes,
		cloneRoleBindings(created.PlanRoleBindings),
	)
	if err != nil {
		return Snapshot{}, mapDecodeFailure(err)
	}
	computed := v0candidate.ExecutionPlanDigest(sha256.Sum256(created.Registration.ExactPlanBytes))
	if computed != created.Registration.RecomputedPlanDigest ||
		decoded.Digest() != created.Registration.RecomputedPlanDigest ||
		decoded.Digest() != created.Attempt.PlanDigest {
		return Snapshot{}, classified(ClassificationBinding, "stored-plan-digest-mismatch")
	}
	record, began, err := component.store.begin(ctx, created)
	if err != nil {
		return Snapshot{}, err
	}
	if !began {
		return record.Snapshot, nil
	}

	if err := component.backend.prepare(ctx, attemptID); err != nil {
		return component.failAndCleanup(ctx, attemptID, OperationPrepare)
	}
	if err := component.afterEffect(
		ctx,
		CheckpointAfterPrepareEffect,
		attemptID,
		OperationPrepare,
	); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	if err := component.store.transition(ctx, attemptID, StateCreateIntent); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	handle, err := component.backend.create(ctx, attemptID)
	if err != nil || handle == 0 {
		return component.failAndCleanup(ctx, attemptID, OperationCreate)
	}
	if err := component.afterEffect(
		ctx,
		CheckpointAfterCreateEffect,
		attemptID,
		OperationCreate,
	); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	if err := component.store.retainHandle(ctx, attemptID, handle); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	if err := component.store.transition(ctx, attemptID, StateStarting); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	if err := component.backend.start(ctx, attemptID, handle); err != nil {
		return component.failAndCleanup(ctx, attemptID, OperationStart)
	}
	if err := component.afterEffect(
		ctx,
		CheckpointAfterStartEffect,
		attemptID,
		OperationStart,
	); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	if err := component.store.transition(ctx, attemptID, StateObserving); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	if err := component.backend.observe(ctx, attemptID, handle); err != nil {
		return component.failAndCleanup(ctx, attemptID, OperationObserve)
	}
	if err := component.afterEffect(
		ctx,
		CheckpointAfterObserveEffect,
		attemptID,
		OperationObserve,
	); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	return component.cleanup(ctx, attemptID)
}

// Recover accepts only the original attempt ID. If no lifecycle record exists,
// it resolves and drives the durable created attempt. Existing cleanup work is
// recovered solely from retained lifecycle state, without rechecking approval
// usability or registration expiry.
func (component *Component) Recover(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (Snapshot, error) {
	if attemptID == (approvalattempt.AttemptID{}) {
		return Snapshot{}, classified(ClassificationBinding, "attempt-id-required")
	}
	record, exists, err := component.store.lookup(ctx, attemptID)
	if err != nil {
		return Snapshot{}, err
	}
	if !exists {
		return component.Drive(ctx, attemptID)
	}
	if record.State == StateDestroyed {
		return record.Snapshot, nil
	}
	if !record.CleanupRequired {
		return record.Snapshot, classified(ClassificationLocalFailure, "nonterminal-cleanup-state-invalid")
	}
	return component.cleanup(ctx, attemptID)
}

// RecoverCreatedAttempts enumerates the durable created-attempt authority
// store once and drives or recovers each distinct AttemptID. It continues
// through fixed lifecycle failures so one broken attempt cannot hide later
// cleanup work, and returns the first fixed error after processing the set.
func (component *Component) RecoverCreatedAttempts(ctx context.Context) ([]Snapshot, error) {
	references, err := component.createdAttempts.CreatedAttempts(ctx)
	if err != nil {
		return nil, mapAttemptResolutionFailure(err)
	}
	seen := make(map[approvalattempt.AttemptID]struct{}, len(references))
	for _, reference := range references {
		attemptID := reference.AttemptID()
		if attemptID == (approvalattempt.AttemptID{}) {
			return nil, classified(ClassificationRecoveryRequired, "startup-attempt-id-invalid")
		}
		if _, duplicate := seen[attemptID]; duplicate {
			return nil, classified(ClassificationRecoveryRequired, "startup-attempt-id-duplicate")
		}
		seen[attemptID] = struct{}{}
	}
	results := make([]Snapshot, 0, len(references))
	var firstErr error
	for _, reference := range references {
		snapshot, recoverErr := component.Recover(ctx, reference.AttemptID())
		results = append(results, snapshot)
		if firstErr == nil && recoverErr != nil {
			firstErr = recoverErr
		}
	}
	return results, firstErr
}

func (component *Component) failAndCleanup(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	operation Operation,
) (Snapshot, error) {
	if err := component.store.noteFailure(
		ctx,
		attemptID,
		ClassificationLifecycleFailure,
		operation,
	); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	snapshot, cleanupErr := component.cleanup(ctx, attemptID)
	if cleanupErr != nil {
		return snapshot, cleanupErr
	}
	return snapshot, classified(ClassificationLifecycleFailure, "fake-lifecycle-operation-failed")
}

func (component *Component) cleanup(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (Snapshot, error) {
	record, err := component.store.snapshot(ctx, attemptID)
	if err != nil {
		return Snapshot{}, err
	}
	if record.State == StateDestroyed {
		return record.Snapshot, nil
	}
	if record.State == StateDestroying {
		observation, reconcileErr := component.backend.reconcile(ctx, attemptID)
		if reconcileErr != nil || observation.status == fakeUnknown {
			return component.unresolved(ctx, attemptID)
		}
		if observation.status == fakeAbsent {
			if err := component.store.markDestroyed(ctx, attemptID); err != nil {
				return component.snapshotWithError(ctx, attemptID, err)
			}
			return component.store.Snapshot(ctx, attemptID)
		}
	}
	if record.handle == 0 {
		observation, reconcileErr := component.backend.reconcile(ctx, attemptID)
		if reconcileErr != nil || observation.status == fakeUnknown {
			return component.unresolved(ctx, attemptID)
		}
		switch observation.status {
		case fakeAbsent:
			if err := component.store.markDestroyed(ctx, attemptID); err != nil {
				return component.snapshotWithError(ctx, attemptID, err)
			}
			return component.store.Snapshot(ctx, attemptID)
		case fakePresent:
			if observation.handle == 0 {
				return component.unresolved(ctx, attemptID)
			}
			if err := component.store.retainHandle(ctx, attemptID, observation.handle); err != nil {
				return component.snapshotWithError(ctx, attemptID, err)
			}
			record.handle = observation.handle
		default:
			return component.unresolved(ctx, attemptID)
		}
	}

	if err := component.store.transition(ctx, attemptID, StateStopping); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	stopErr := component.backend.stop(ctx, attemptID, record.handle)
	if stopErr != nil {
		_ = component.store.noteFailure(
			ctx,
			attemptID,
			ClassificationLifecycleFailure,
			OperationStop,
		)
	} else if err := component.afterEffect(
		ctx,
		CheckpointAfterStopEffect,
		attemptID,
		OperationStop,
	); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	if err := component.store.transition(ctx, attemptID, StateDestroying); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	destroyErr := component.backend.destroy(ctx, attemptID, record.handle)
	if destroyErr != nil {
		_ = component.store.noteFailure(
			ctx,
			attemptID,
			ClassificationLifecycleFailure,
			OperationDestroy,
		)
	} else if err := component.afterEffect(
		ctx,
		CheckpointAfterDestroyEffect,
		attemptID,
		OperationDestroy,
	); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}

	observation, reconcileErr := component.backend.reconcile(ctx, attemptID)
	if reconcileErr != nil || observation.status != fakeAbsent {
		return component.unresolved(ctx, attemptID)
	}
	if err := component.store.markDestroyed(ctx, attemptID); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	snapshot, err := component.store.Snapshot(ctx, attemptID)
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
	attemptID approvalattempt.AttemptID,
) (Snapshot, error) {
	if err := component.store.markUnresolved(ctx, attemptID); err != nil {
		return component.snapshotWithError(ctx, attemptID, err)
	}
	snapshot, err := component.store.Snapshot(ctx, attemptID)
	if err != nil {
		return Snapshot{}, err
	}
	return snapshot, classified(ClassificationCleanupUnresolved, "fake-lifecycle-cleanup-unresolved")
}

func (component *Component) afterEffect(
	ctx context.Context,
	checkpoint Checkpoint,
	attemptID approvalattempt.AttemptID,
	operation Operation,
) error {
	if err := component.checkpoint(ctx, checkpoint, attemptID); err != nil {
		_ = component.store.noteFailure(
			ctx,
			attemptID,
			ClassificationLocalFailure,
			operation,
		)
		return classified(ClassificationLocalFailure, "simulated-lifecycle-interruption")
	}
	return nil
}

func (component *Component) snapshotWithError(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	err error,
) (Snapshot, error) {
	snapshot, snapshotErr := component.store.Snapshot(ctx, attemptID)
	if snapshotErr != nil {
		return Snapshot{}, err
	}
	return snapshot, err
}

func validateCreatedAttempt(
	created registrationstate.CreatedAttempt,
	expectedAttemptID approvalattempt.AttemptID,
) error {
	attempt := created.Attempt
	approval := created.Approval
	if attempt.AttemptID != expectedAttemptID ||
		attempt.State != approvalattempt.AttemptCreated ||
		attempt.ApprovalID == (approvalattempt.ApprovalID{}) ||
		attempt.AttemptNonce == (approvalattempt.AttemptNonce{}) ||
		attempt.RegistrationID == (v0candidate.RegistrationID{}) ||
		attempt.RegistrationSequence == 0 ||
		attempt.PlanDigest == (v0candidate.ExecutionPlanDigest{}) ||
		attempt.InstallationID == (v0candidate.InstallationID{}) ||
		attempt.EpochDigest == (v0candidate.TrustEpochDigest{}) ||
		attempt.SupervisorID == (v0candidate.SupervisorID{}) ||
		attempt.Purpose != approvalattempt.ApprovalGrantPurpose ||
		attempt.Audience != approvalattempt.ApprovalGrantAudience ||
		attempt.ApprovalPayloadDigest == (approvalattempt.ApprovalPayloadDigest{}) ||
		attempt.AuthorizationIdentity == (approvalattempt.ApprovalKeyAuthorizationIdentity{}) ||
		attempt.StorageFormatVersion != registrationstate.AttemptStoreVersion {
		return classified(ClassificationBinding, "created-attempt-binding-mismatch")
	}
	if approval.State != approvalattempt.ApprovalConsumed ||
		approval.StorageFormatVersion != registrationstate.ApprovalStoreVersion ||
		approval.ApprovalID != attempt.ApprovalID ||
		approval.ConsumedAttemptID != attempt.AttemptID ||
		approval.AttemptNonce != attempt.AttemptNonce ||
		approval.RegistrationID != attempt.RegistrationID ||
		approval.RegistrationSequence != attempt.RegistrationSequence ||
		approval.PlanDigest != attempt.PlanDigest ||
		approval.InstallationID != attempt.InstallationID ||
		approval.EpochSequence != attempt.EpochSequence ||
		approval.EpochDigest != attempt.EpochDigest ||
		approval.SupervisorID != attempt.SupervisorID ||
		approval.Purpose != attempt.Purpose || approval.Audience != attempt.Audience ||
		approval.PayloadDigest != attempt.ApprovalPayloadDigest ||
		approval.AuthorizationIdentity != attempt.AuthorizationIdentity {
		return classified(ClassificationBinding, "created-approval-cross-link-mismatch")
	}
	record := created.Registration
	if record.StorageFormatVersion != registrationstate.StorageFormatVersion ||
		record.RetentionState != registrationstate.RetentionStateRetained ||
		record.RegisteredAtUnixSeconds > attempt.CreatedAt ||
		len(record.ExactPlanBytes) == 0 || len(record.WireRegistrationBytes) == 0 {
		return classified(ClassificationBinding, "created-registration-record-invalid")
	}
	decodedRegistration, err := v0candidate.DecodePlanRegistration(
		record.WireRegistrationBytes,
		v0candidate.PlanRegistrationRoleBindings{
			RegistrationID: attempt.RegistrationID,
			PlanDigest:     attempt.PlanDigest,
			InstallationID: attempt.InstallationID,
			EpochDigest:    attempt.EpochDigest,
			SupervisorID:   attempt.SupervisorID,
		},
	)
	if err != nil {
		return classified(ClassificationBinding, "created-registration-wire-invalid")
	}
	registration := decodedRegistration.View()
	if registration.RegistrationSequence != attempt.RegistrationSequence ||
		registration.EpochSequence != attempt.EpochSequence ||
		attempt.CreatedAt >= registration.ExpiresAt {
		return classified(ClassificationBinding, "created-registration-binding-mismatch")
	}
	return nil
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
