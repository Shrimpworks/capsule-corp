package registeredlifecycle

import (
	"context"
	"crypto/sha256"
	"errors"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/execution/registrationstate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// Component drives only durable created attempts from the colocated v1 store
// into the concrete no-guest FakeBackend. It is not wired to a deployed
// Supervisor, endpoint, content, evidence, runtime, real backend, or guest.
type Component struct {
	attempts    AttemptResolver
	store       registrationstate.DurableLifecycleStore
	backend     *FakeBackend
	coordinator *Coordinator
	clock       registrationstate.TrustedClock
	checkpoint  CheckpointHook
}

func New(options Options) (*Component, error) {
	if options.Attempts == nil || options.Store == nil || options.Backend == nil ||
		options.Coordinator == nil || options.Clock == nil {
		return nil, errors.New("attempt resolver, durable store, fake backend, coordinator, and clock are required")
	}
	if options.Backend.CreatesGuest() || options.Backend.Binding().View().CreatesGuest {
		return nil, errors.New("registered lifecycle requires a no-guest fake backend")
	}
	if options.Coordinator.OwnerSessionID() != options.Store.OwnerSessionID() {
		return nil, errors.New("lifecycle coordinator and store owner sessions must match")
	}
	if options.Checkpoint == nil {
		options.Checkpoint = func(context.Context, Checkpoint, approvalattempt.AttemptID) error { return nil }
	}
	return &Component{
		attempts: options.Attempts, store: options.Store, backend: options.Backend,
		coordinator: options.Coordinator, clock: options.Clock, checkpoint: options.Checkpoint,
	}, nil
}

// Drive accepts only a Supervisor-issued AttemptID. It independently
// revalidates the committed attempt and exact registered plan before the
// durable lifecycle record can be established or any fake effect can run.
func (component *Component) Drive(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (Snapshot, error) {
	if attemptID == (approvalattempt.AttemptID{}) {
		return Snapshot{}, classified(ClassificationBinding, "attempt-id-required")
	}
	return component.coordinator.withAttempt(ctx, attemptID, func() (Snapshot, error) {
		return component.driveLocked(ctx, attemptID)
	})
}

func (component *Component) driveLocked(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (Snapshot, error) {
	if err := component.advanceTime(ctx); err != nil {
		return Snapshot{}, err
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
	record, _, err := component.store.EnsureLifecycle(ctx, attemptID, component.backend.Binding())
	if err != nil {
		return Snapshot{}, mapStoreFailure(err)
	}
	return component.advanceLocked(ctx, record)
}

// Recover never re-evaluates approval usability or registration expiry. A
// missing lifecycle record remains startup work for the already committed
// AttemptID; every other decision uses only durable lifecycle state.
func (component *Component) Recover(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (Snapshot, error) {
	if attemptID == (approvalattempt.AttemptID{}) {
		return Snapshot{}, classified(ClassificationBinding, "attempt-id-required")
	}
	return component.coordinator.withAttempt(ctx, attemptID, func() (Snapshot, error) {
		record, err := component.store.ReadLifecycle(ctx, attemptID)
		if errors.Is(err, registrationstate.ErrLifecycleNotFound) {
			if timeErr := component.advanceTime(ctx); timeErr != nil {
				return Snapshot{}, timeErr
			}
			return component.driveLocked(ctx, attemptID)
		}
		if err != nil {
			return Snapshot{}, mapStoreFailure(err)
		}
		if timeErr := component.advanceTime(ctx); timeErr != nil {
			return component.cleanupWithPriorTimeLocked(ctx, record, timeErr)
		}
		return component.advanceLocked(ctx, record)
	})
}

// RecoverCreatedAttempts uses the colocated store's sorted join of committed
// created attempts and nonterminal lifecycle work. It processes every entry
// and returns the first fixed error after the complete initial pass.
func (component *Component) RecoverCreatedAttempts(ctx context.Context) ([]Snapshot, error) {
	identifiers, err := component.store.RecoveryAttemptIDs(ctx)
	if err != nil {
		return nil, mapStoreFailure(err)
	}
	results := make([]Snapshot, 0, len(identifiers))
	var firstErr error
	for _, attemptID := range identifiers {
		if attemptID == (approvalattempt.AttemptID{}) {
			return nil, classified(ClassificationRecoveryRequired, "startup-attempt-id-invalid")
		}
		snapshot, recoverErr := component.Recover(ctx, attemptID)
		results = append(results, snapshot)
		if firstErr == nil && recoverErr != nil {
			firstErr = recoverErr
		}
	}
	return results, firstErr
}

// cleanupWithPriorTimeLocked preserves ADR-0025's trusted-clock failure rule:
// no new prepare/create/start/observe effect may begin, while exact
// reconciliation and destructive stop/destroy work may use the existing
// durable high-water timestamp.
func (component *Component) cleanupWithPriorTimeLocked(
	ctx context.Context,
	record lifecyclestate.Record,
	clockErr error,
) (Snapshot, error) {
	for {
		view := record.View()
		switch view.State {
		case lifecyclestate.StateDestroyed:
			return terminalSnapshot(record)
		case lifecyclestate.StateQuarantined:
			return snapshotFromRecord(record), classified(ClassificationTrustState, "lifecycle-quarantined")
		case lifecyclestate.StateCreated, lifecyclestate.StateStarted, lifecyclestate.StateObserved:
			return component.performAndContinue(ctx, record, lifecyclestate.OperationStop)
		case lifecyclestate.StateStopped:
			return component.performAndContinue(ctx, record, lifecyclestate.OperationDestroy)
		case lifecyclestate.StateUnresolved:
			if retryableConfirmedNotApplied(view) {
				if view.Operation == lifecyclestate.OperationStop || view.Operation == lifecyclestate.OperationDestroy {
					return component.performAndContinue(ctx, record, view.Operation)
				}
				return snapshotFromRecord(record), clockErr
			}
			var err error
			record, err = component.reconcileLocked(ctx, record)
			if err != nil {
				return snapshotFromRecord(record), err
			}
		case lifecyclestate.StatePrepareIntent,
			lifecyclestate.StateCreateIntent,
			lifecyclestate.StateStartIntent,
			lifecyclestate.StateObserveIntent,
			lifecyclestate.StateStopIntent,
			lifecyclestate.StateDestroyIntent,
			lifecyclestate.StateDestroyConfirmed:
			var err error
			record, err = component.reconcileLocked(ctx, record)
			if err != nil {
				return snapshotFromRecord(record), err
			}
		case lifecyclestate.StatePreparePending, lifecyclestate.StatePrepared:
			return snapshotFromRecord(record), clockErr
		default:
			return snapshotFromRecord(record), classified(ClassificationRecoveryRequired, "lifecycle-state-invalid")
		}
	}
}

func (component *Component) advanceLocked(
	ctx context.Context,
	record lifecyclestate.Record,
) (Snapshot, error) {
	for {
		view := record.View()
		switch view.State {
		case lifecyclestate.StateDestroyed:
			return terminalSnapshot(record)
		case lifecyclestate.StateQuarantined:
			return snapshotFromRecord(record), classified(ClassificationTrustState, "lifecycle-quarantined")
		case lifecyclestate.StatePreparePending:
			return component.performAndContinue(ctx, record, lifecyclestate.OperationPrepare)
		case lifecyclestate.StatePrepared:
			return component.performAndContinue(ctx, record, lifecyclestate.OperationCreate)
		case lifecyclestate.StateCreated:
			return component.performAndContinue(ctx, record, lifecyclestate.OperationStart)
		case lifecyclestate.StateStarted:
			return component.performAndContinue(ctx, record, lifecyclestate.OperationObserve)
		case lifecyclestate.StateObserved:
			return component.performAndContinue(ctx, record, lifecyclestate.OperationStop)
		case lifecyclestate.StateStopped:
			return component.performAndContinue(ctx, record, lifecyclestate.OperationDestroy)
		case lifecyclestate.StateUnresolved:
			if retryableConfirmedNotApplied(view) {
				return component.performAndContinue(ctx, record, view.Operation)
			}
			var err error
			record, err = component.reconcileLocked(ctx, record)
			if err != nil {
				return snapshotFromRecord(record), err
			}
		case lifecyclestate.StatePrepareIntent,
			lifecyclestate.StateCreateIntent,
			lifecyclestate.StateStartIntent,
			lifecyclestate.StateObserveIntent,
			lifecyclestate.StateStopIntent,
			lifecyclestate.StateDestroyIntent,
			lifecyclestate.StateDestroyConfirmed:
			var err error
			record, err = component.reconcileLocked(ctx, record)
			if err != nil {
				return snapshotFromRecord(record), err
			}
		default:
			return snapshotFromRecord(record), classified(ClassificationRecoveryRequired, "lifecycle-state-invalid")
		}
	}
}

func (component *Component) performAndContinue(
	ctx context.Context,
	record lifecyclestate.Record,
	operation lifecyclestate.Operation,
) (Snapshot, error) {
	view := record.View()
	permit, err := component.store.BeginEffect(
		ctx, view.Bindings.View().AttemptID, view.RecordVersion, operation,
	)
	if err != nil {
		return snapshotFromRecord(record), mapStoreFailure(err)
	}
	result, adapterErr := component.backend.apply(ctx, permit)
	if result.Status() == lifecyclestate.EffectResultApplied {
		if checkpointErr := component.checkpoint(
			ctx, checkpointFor(operation), permit.View().AttemptID,
		); checkpointErr != nil {
			current, readErr := component.store.ReadLifecycle(ctx, permit.View().AttemptID)
			if readErr == nil {
				return snapshotFromRecord(current), classified(ClassificationLocalFailure, "simulated-lifecycle-interruption")
			}
			return Snapshot{}, classified(ClassificationLocalFailure, "simulated-lifecycle-interruption")
		}
	}
	switch result.Status() {
	case lifecyclestate.EffectResultApplied, lifecyclestate.EffectResultNotApplied:
		confirmed, confirmErr := component.store.ConfirmEffect(ctx, permit, result)
		if confirmErr != nil {
			return component.snapshotWithError(ctx, permit.View().AttemptID, mapStoreFailure(confirmErr))
		}
		return component.advanceLocked(ctx, confirmed)
	case lifecyclestate.EffectResultIndeterminate:
		indeterminate, recordErr := component.store.RecordIndeterminate(
			ctx, permit, registrationstate.ClassificationLifecycleFailure,
		)
		if recordErr != nil {
			return component.snapshotWithError(ctx, permit.View().AttemptID, mapStoreFailure(recordErr))
		}
		return snapshotFromRecord(indeterminate), classified(ClassificationLifecycleFailure, "fake-effect-indeterminate")
	default:
		if adapterErr == nil {
			adapterErr = errors.New("fake adapter returned no closed result")
		}
		indeterminate, recordErr := component.store.RecordIndeterminate(
			ctx, permit, registrationstate.ClassificationLifecycleFailure,
		)
		if recordErr != nil {
			return component.snapshotWithError(ctx, permit.View().AttemptID, mapStoreFailure(recordErr))
		}
		_ = adapterErr
		return snapshotFromRecord(indeterminate), classified(ClassificationLifecycleFailure, "fake-effect-indeterminate")
	}
}

func (component *Component) reconcileLocked(
	ctx context.Context,
	record lifecyclestate.Record,
) (lifecyclestate.Record, error) {
	view := record.View()
	attemptID := view.Bindings.View().AttemptID
	eligible := record
	// A prior result-commit failure leaves a durable eligibility checkpoint
	// with LastReconciliation == none. Reopen repeats only the observation.
	if view.AutomaticRecoveryCount == 0 || view.LastReconciliation != lifecyclestate.ReconciliationNone {
		begun, err := component.store.BeginReconciliation(ctx, attemptID, view.RecordVersion)
		if err != nil {
			classification, ok := registrationstate.ErrorClassification(err)
			if ok && (classification == registrationstate.ClassificationStale ||
				classification == registrationstate.ClassificationCapacity) {
				return record, classified(ClassificationCleanupUnresolved, "reconciliation-not-eligible")
			}
			return record, mapStoreFailure(err)
		}
		eligible = begun
	}
	result, _ := component.backend.reconcile(ctx, eligible)
	if result.Status() == "" {
		var resultErr error
		result, resultErr = lifecyclestate.NewReconcileResult(
			eligible.View().Operation,
			lifecyclestate.ReconciliationUnknown,
			lifecyclestate.BackendInstanceIdentity{},
		)
		if resultErr != nil {
			return eligible, classified(ClassificationLocalFailure, "fake-reconciliation-result")
		}
	}
	completed, err := component.store.CompleteReconciliation(
		ctx, attemptID, eligible.View().RecordVersion, result,
	)
	if err != nil {
		return eligible, mapStoreFailure(err)
	}
	switch completed.View().State {
	case lifecyclestate.StateUnresolved:
		if retryableConfirmedNotApplied(completed.View()) ||
			completed.View().LastReconciliation == lifecyclestate.ReconciliationNotApplied {
			return completed, nil
		}
		return completed, classified(ClassificationCleanupUnresolved, "fake-reconciliation-unresolved")
	case lifecyclestate.StateQuarantined:
		return completed, classified(ClassificationTrustState, "fake-reconciliation-identity-mismatch")
	default:
		return completed, nil
	}
}

func retryableConfirmedNotApplied(view lifecyclestate.RecordView) bool {
	return view.State == lifecyclestate.StateUnresolved &&
		((view.EffectStatus == lifecyclestate.EffectConfirmed &&
			view.FirstFailure == lifecyclestate.FailureLifecycle &&
			view.RecoveryFence == lifecyclestate.RecoveryFenceNone) ||
			(view.LastReconciliation == lifecyclestate.ReconciliationNotApplied &&
				(view.EffectStatus == lifecyclestate.EffectIntent ||
					view.EffectStatus == lifecyclestate.EffectIndeterminate)))
}

func checkpointFor(operation lifecyclestate.Operation) Checkpoint {
	switch operation {
	case lifecyclestate.OperationPrepare:
		return CheckpointAfterPrepareEffect
	case lifecyclestate.OperationCreate:
		return CheckpointAfterCreateEffect
	case lifecyclestate.OperationStart:
		return CheckpointAfterStartEffect
	case lifecyclestate.OperationObserve:
		return CheckpointAfterObserveEffect
	case lifecyclestate.OperationStop:
		return CheckpointAfterStopEffect
	case lifecyclestate.OperationDestroy:
		return CheckpointAfterDestroyEffect
	default:
		return ""
	}
}

func (component *Component) advanceTime(ctx context.Context) error {
	observed, err := component.clock.ObserveUnixSeconds(ctx)
	if err != nil || observed > v0candidate.MaxSafeInteger {
		return classified(ClassificationLocalFailure, "trusted-clock-observation")
	}
	if err := component.store.AdvanceLifecycleTime(ctx, v0candidate.UInt53(observed)); err != nil {
		return mapStoreFailure(err)
	}
	return nil
}

func (component *Component) snapshotWithError(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	err error,
) (Snapshot, error) {
	record, readErr := component.store.ReadLifecycle(ctx, attemptID)
	if readErr != nil {
		return Snapshot{}, err
	}
	return snapshotFromRecord(record), err
}

func terminalSnapshot(record lifecyclestate.Record) (Snapshot, error) {
	snapshot := snapshotFromRecord(record)
	switch snapshot.Failure {
	case "":
		return snapshot, nil
	case ClassificationCleanupUnresolved:
		return snapshot, classified(ClassificationCleanupUnresolved, "fake-lifecycle-recovered-cleanup")
	default:
		return snapshot, classified(snapshot.Failure, "fake-lifecycle-retained-failure")
	}
}

func snapshotFromRecord(record lifecyclestate.Record) Snapshot {
	view := record.View()
	bindings := view.Bindings.View()
	failure := Classification("")
	failureAt := Operation("")
	if view.FirstFailure != lifecyclestate.FailureNone {
		failure = Classification(view.FirstFailure)
		failureAt = Operation(view.FailureOperation)
	}
	return Snapshot{
		AttemptID: bindings.AttemptID, ApprovalID: bindings.ApprovalID,
		AttemptNonce: bindings.AttemptNonce, RegistrationID: bindings.RegistrationID,
		RegistrationSequence: bindings.RegistrationSequence, PlanDigest: bindings.PlanDigest,
		InstallationID: bindings.InstallationID, EpochSequence: bindings.EpochSequence,
		EpochDigest: bindings.EpochDigest, SupervisorID: bindings.SupervisorID,
		ApprovalPurpose: bindings.ApprovalPurpose, ApprovalAudience: bindings.ApprovalAudience,
		ApprovalPayloadDigest: bindings.ApprovalPayloadDigest,
		AuthorizationIdentity: bindings.AuthorizationIdentity, CreatedAt: bindings.AttemptCreatedAt,
		State: view.State, CleanupRequired: view.CleanupRequired,
		Failure: failure, FailureAt: failureAt, TransitionCount: uint64(view.RecordVersion),
		RecordVersion: view.RecordVersion, OperationSequence: view.OperationSequence,
		EffectID: view.EffectID, EffectStatus: view.EffectStatus,
		InstanceDigest: view.Instance.Digest(), LastReconciliation: view.LastReconciliation,
		AutomaticRecoveryCount: view.AutomaticRecoveryCount, NextRecoveryAt: view.NextRecoveryAt,
		RecoveryFence:              view.RecoveryFence,
		AutomaticRecoveryExhausted: view.AutomaticRecoveryExhausted,
	}
}

func mapStoreFailure(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, registrationstate.ErrStoreRepairRequired) {
		return classified(ClassificationRecoveryRequired, "durable-lifecycle-repair-required")
	}
	if errors.Is(err, registrationstate.ErrCommitOutcomeIndeterminate) {
		return classified(ClassificationRecoveryRequired, "durable-lifecycle-commit-indeterminate")
	}
	if errors.Is(err, registrationstate.ErrStoreOwnerFenced) ||
		errors.Is(err, registrationstate.ErrStoreOwnerSessionMismatch) {
		return classified(ClassificationRecoveryRequired, "durable-lifecycle-owner-fenced")
	}
	if errors.Is(err, registrationstate.ErrStoreClosed) {
		return classified(ClassificationRecoveryRequired, "durable-lifecycle-store-closed")
	}
	if errors.Is(err, registrationstate.ErrLifecycleNotFound) ||
		errors.Is(err, registrationstate.ErrLifecycleBindingMismatch) ||
		errors.Is(err, registrationstate.ErrLifecycleVersionConflict) {
		return classified(ClassificationBinding, "durable-lifecycle-binding")
	}
	classification, ok := registrationstate.ErrorClassification(err)
	if !ok {
		return classified(ClassificationLocalFailure, "durable-lifecycle-store")
	}
	return classified(Classification(classification), "durable-lifecycle-store")
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
			RegistrationID: attempt.RegistrationID, PlanDigest: attempt.PlanDigest,
			InstallationID: attempt.InstallationID, EpochDigest: attempt.EpochDigest,
			SupervisorID: attempt.SupervisorID,
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
