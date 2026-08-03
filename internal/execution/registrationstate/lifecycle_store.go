package registrationstate

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// LifecycleStoreFault names every E3 lifecycle transaction boundary. Abort
// faults are confirmed before rename. Pre-state indeterminate faults model a
// lost outcome before a candidate is known durable. Commit-indeterminate
// faults occur after rename and require reopen before any further mutation.
type LifecycleStoreFault string

const (
	FaultLifecycleEnsureAbort                 LifecycleStoreFault = "lifecycle-ensure-confirmed-abort"
	FaultLifecycleEnsureIndeterminatePreState LifecycleStoreFault = "lifecycle-ensure-indeterminate-pre-state"
	FaultLifecycleEnsureIndeterminate         LifecycleStoreFault = "lifecycle-ensure-commit-indeterminate"

	FaultLifecycleIntentAbort                 LifecycleStoreFault = "lifecycle-intent-confirmed-abort"
	FaultLifecycleIntentIndeterminatePreState LifecycleStoreFault = "lifecycle-intent-indeterminate-pre-state"
	FaultLifecycleIntentIndeterminate         LifecycleStoreFault = "lifecycle-intent-commit-indeterminate"

	FaultLifecycleAfterEffectAbort                 LifecycleStoreFault = "lifecycle-after-effect-confirmed-abort"
	FaultLifecycleAfterEffectIndeterminatePreState LifecycleStoreFault = "lifecycle-after-effect-indeterminate-pre-state"
	FaultLifecycleAfterEffectIndeterminate         LifecycleStoreFault = "lifecycle-after-effect-commit-indeterminate"

	FaultLifecycleRecordIndeterminateAbort               LifecycleStoreFault = "lifecycle-record-indeterminate-confirmed-abort"
	FaultLifecycleRecordIndeterminatePreState            LifecycleStoreFault = "lifecycle-record-indeterminate-pre-state"
	FaultLifecycleRecordIndeterminateCommitIndeterminate LifecycleStoreFault = "lifecycle-record-indeterminate-commit-indeterminate"

	FaultLifecycleReconciliationEligibilityAbort                 LifecycleStoreFault = "lifecycle-reconciliation-eligibility-confirmed-abort"
	FaultLifecycleReconciliationEligibilityIndeterminatePreState LifecycleStoreFault = "lifecycle-reconciliation-eligibility-indeterminate-pre-state"
	FaultLifecycleReconciliationEligibilityIndeterminate         LifecycleStoreFault = "lifecycle-reconciliation-eligibility-commit-indeterminate"

	FaultLifecycleReconciliationResultAbort                 LifecycleStoreFault = "lifecycle-reconciliation-result-confirmed-abort"
	FaultLifecycleReconciliationResultIndeterminatePreState LifecycleStoreFault = "lifecycle-reconciliation-result-indeterminate-pre-state"
	FaultLifecycleReconciliationResultIndeterminate         LifecycleStoreFault = "lifecycle-reconciliation-result-commit-indeterminate"
)

var (
	ErrLifecycleNotFound        = errors.New("fixed supervisor lifecycle record not found")
	ErrLifecycleVersionConflict = errors.New("fixed supervisor lifecycle record version conflict")
	ErrLifecycleBindingMismatch = errors.New("fixed supervisor lifecycle binding mismatch")
)

// LifecycleEffectIDSource is the trusted local generator used immediately
// before an intent transaction. It cannot supply any backend configuration or
// execute authority.
type LifecycleEffectIDSource interface {
	NewEffectID(context.Context) (lifecyclestate.EffectID, error)
}

type CryptoLifecycleEffectIDSource struct{}

func (CryptoLifecycleEffectIDSource) NewEffectID(ctx context.Context) (lifecyclestate.EffectID, error) {
	if err := ctx.Err(); err != nil {
		return lifecyclestate.EffectID{}, err
	}
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return lifecyclestate.EffectID{}, err
	}
	domain, err := lifecyclestate.NewDomainIdentifier(lifecyclestate.DomainEffectID, value)
	if err != nil {
		return lifecyclestate.EffectID{}, err
	}
	return lifecyclestate.NewEffectID(domain)
}

type FixedFileStoreV1Options struct {
	EffectIDs      LifecycleEffectIDSource
	OwnerSessionID lifecyclestate.OwnerSessionID
}

func resolveFixedFileStoreV1Options(options FixedFileStoreV1Options) (FixedFileStoreV1Options, error) {
	if options.EffectIDs == nil {
		options.EffectIDs = CryptoLifecycleEffectIDSource{}
	}
	if options.OwnerSessionID.IsZero() {
		value := make([]byte, 16)
		if _, err := rand.Read(value); err != nil {
			return FixedFileStoreV1Options{}, errors.New("generate lifecycle owner session")
		}
		domain, err := lifecyclestate.NewDomainIdentifier(lifecyclestate.DomainOwnerSessionID, value)
		if err != nil {
			return FixedFileStoreV1Options{}, err
		}
		owner, err := lifecyclestate.NewOwnerSessionID(domain)
		if err != nil {
			return FixedFileStoreV1Options{}, err
		}
		options.OwnerSessionID = owner
	}
	return options, nil
}

var _ DurableLifecycleStore = (*FixedFileStoreV1)(nil)

// InjectLifecycleFailure installs one deterministic one-shot E3 fault.
func (store *FixedFileStoreV1) InjectLifecycleFailure(point LifecycleStoreFault, err error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	if err == nil {
		err = errors.New("injected lifecycle-store failure")
	}
	store.faults[point] = err
}

func (store *FixedFileStoreV1) LifecycleRecoveryFenced() bool {
	store.mu.Lock()
	defer store.mu.Unlock()
	return store.recoveryRequired || store.repairRequired
}

// OwnerSessionID identifies the sealed permit owner for this open store
// handle. It is used only to bind the injected E4 in-process coordinator to
// the same owner session; it grants no execute or backend authority.
func (store *FixedFileStoreV1) OwnerSessionID() lifecyclestate.OwnerSessionID {
	store.mu.Lock()
	defer store.mu.Unlock()
	return store.ownerSessionID
}

// AdvanceLifecycleTime durably advances the installation high-water used by
// lifecycle transitions and retry eligibility. It never moves time backward
// and does not create or widen approval, attempt, or backend authority.
func (store *FixedFileStoreV1) AdvanceLifecycleTime(
	ctx context.Context,
	observed v0candidate.UInt53,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if uint64(observed) > v0candidate.MaxSafeInteger {
		return classified(ClassificationSchema, "lifecycle-time-uint53")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.requireMutableLocked(); err != nil {
		return err
	}
	if observed <= store.state.TimeHighWaterUnixSeconds {
		return nil
	}
	return store.commitLifecycleLocked(
		ctx, "", "", "", false,
		func(state *installationState, _ *[]lifecyclestate.Record) error {
			state.TimeHighWaterUnixSeconds = observed
			return nil
		},
	)
}

func (store *FixedFileStoreV1) EnsureLifecycle(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	backend lifecyclestate.BackendBinding,
) (lifecyclestate.Record, bool, error) {
	if err := ctx.Err(); err != nil {
		return lifecyclestate.Record{}, false, err
	}
	validatedBackend, err := lifecyclestate.NewBackendBinding(backend.View())
	if err != nil {
		return lifecyclestate.Record{}, false, classified(ClassificationUnsupported, "lifecycle-backend-binding")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.requireMutableLocked(); err != nil {
		return lifecyclestate.Record{}, false, err
	}
	existing := findLifecycle(store.lifecycles, attemptID)
	bindings, err := lifecycleBindingsFromAuthority(store.state, attemptID, validatedBackend)
	if err != nil {
		if existing >= 0 || errors.Is(err, ErrStoreRepairRequired) {
			store.repairRequired = true
			return lifecyclestate.Record{}, false, fmt.Errorf("%w: lifecycle authority binding", ErrStoreRepairRequired)
		}
		return lifecyclestate.Record{}, false, err
	}
	if existing >= 0 {
		record := cloneLifecycleRecord(store.lifecycles[existing], store.state.TimeHighWaterUnixSeconds)
		if !record.Bindings().Equal(bindings) {
			store.repairRequired = true
			return lifecyclestate.Record{}, false, fmt.Errorf("%w: existing lifecycle binding", ErrStoreRepairRequired)
		}
		return record, false, nil
	}
	if len(store.lifecycles) >= MaxRetainedLifecycleRecords || activeLifecycleWork(store.state, store.lifecycles) > MaxActiveLifecycleRecords {
		return lifecyclestate.Record{}, false, classified(ClassificationCapacity, "lifecycle-capacity")
	}
	openedAt := store.state.TimeHighWaterUnixSeconds
	if openedAt < bindings.View().AttemptCreatedAt {
		store.repairRequired = true
		return lifecyclestate.Record{}, false, fmt.Errorf("%w: lifecycle time cross-link", ErrStoreRepairRequired)
	}
	record, err := lifecyclestate.NewRecord(lifecyclestate.RecordView{
		FormatVersion:           lifecyclestate.LifecycleRecordFormatVersion,
		RecordVersion:           1,
		SnapshotGeneration:      1,
		Bindings:                bindings,
		ImmutableBindingDigest:  bindings.Digest(),
		Operation:               lifecyclestate.OperationNone,
		EffectStatus:            lifecyclestate.EffectNone,
		State:                   lifecyclestate.StatePreparePending,
		CleanupRequired:         true,
		LastConfirmedCheckpoint: lifecyclestate.CheckpointNone,
		FirstFailure:            lifecyclestate.FailureNone,
		FailureOperation:        lifecyclestate.OperationNone,
		LastReconciliation:      lifecyclestate.ReconciliationNone,
		RecoveryFence:           lifecyclestate.RecoveryFenceNone,
		OpenedAt:                openedAt,
		LastTransitionAt:        openedAt,
	}, store.state.TimeHighWaterUnixSeconds)
	if err != nil {
		return lifecyclestate.Record{}, false, err
	}
	var committed lifecyclestate.Record
	err = store.commitLifecycleLocked(
		ctx,
		FaultLifecycleEnsureAbort,
		FaultLifecycleEnsureIndeterminatePreState,
		FaultLifecycleEnsureIndeterminate,
		false,
		func(_ *installationState, records *[]lifecyclestate.Record) error {
			if findLifecycle(*records, attemptID) >= 0 {
				return errNoStoreMutation
			}
			*records = append(*records, record)
			sortLifecycleRecords(*records)
			committed = record
			return nil
		},
	)
	if err != nil {
		return lifecyclestate.Record{}, false, err
	}
	return cloneLifecycleRecord(committed, store.state.TimeHighWaterUnixSeconds), true, nil
}

func (store *FixedFileStoreV1) ReadLifecycle(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (lifecyclestate.Record, error) {
	if err := ctx.Err(); err != nil {
		return lifecyclestate.Record{}, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if store.repairRequired {
		return lifecyclestate.Record{}, ErrStoreRepairRequired
	}
	index := findLifecycle(store.lifecycles, attemptID)
	if index < 0 {
		return lifecyclestate.Record{}, ErrLifecycleNotFound
	}
	return cloneLifecycleRecord(store.lifecycles[index], store.state.TimeHighWaterUnixSeconds), nil
}

func (store *FixedFileStoreV1) BeginEffect(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	expected lifecyclestate.RecordVersion,
	operation lifecyclestate.Operation,
) (lifecyclestate.EffectPermit, error) {
	if err := ctx.Err(); err != nil {
		return lifecyclestate.EffectPermit{}, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.requireMutableLocked(); err != nil {
		return lifecyclestate.EffectPermit{}, err
	}
	index := findLifecycle(store.lifecycles, attemptID)
	if index < 0 {
		return lifecyclestate.EffectPermit{}, ErrLifecycleNotFound
	}
	view := store.lifecycles[index].View()
	if permit, ok := store.issuedPermits[attemptID]; ok {
		permitView := permit.View()
		if permitView.Operation == operation && permitView.EffectID == view.EffectID &&
			permitView.RecordVersion == view.RecordVersion &&
			(expected == view.RecordVersion || uint64(expected)+1 == uint64(view.RecordVersion)) {
			return permit, nil
		}
	}
	if retryableSameEffect(view, expected, operation) {
		permit, err := store.permitForView(view)
		if err != nil {
			return lifecyclestate.EffectPermit{}, err
		}
		store.issuedPermits[attemptID] = permit
		return permit, nil
	}
	if view.RecordVersion != expected {
		return lifecyclestate.EffectPermit{}, ErrLifecycleVersionConflict
	}
	if err := validateEffectPredecessor(view, operation); err != nil {
		return lifecyclestate.EffectPermit{}, err
	}
	effectID, err := store.effectIDs.NewEffectID(ctx)
	if err != nil || effectID.IsZero() {
		return lifecyclestate.EffectPermit{}, classified(ClassificationLocalFailure, "effect-id-source")
	}
	for _, record := range store.lifecycles {
		if record.View().EffectID == effectID {
			return lifecyclestate.EffectPermit{}, classified(ClassificationLocalFailure, "effect-id-collision")
		}
	}
	if uint64(view.OperationSequence) >= v0candidate.MaxSafeInteger {
		return lifecyclestate.EffectPermit{}, classified(ClassificationCapacity, "operation-sequence")
	}
	view.OperationSequence++
	view.Operation = operation
	view.EffectID = effectID
	view.EffectStatus = lifecyclestate.EffectIntent
	view.State = intentState(operation)
	view.LastReconciliation = lifecyclestate.ReconciliationNone
	view.NextRecoveryAt = lifecyclestate.OptionalUnixSeconds{}
	view.RecoveryFence = lifecyclestate.RecoveryFenceNone
	view.AutomaticRecoveryExhausted = false
	if err := advanceLifecycleVersion(&view, store.state.TimeHighWaterUnixSeconds); err != nil {
		return lifecyclestate.EffectPermit{}, err
	}
	updated, err := lifecyclestate.NewRecord(view, store.state.TimeHighWaterUnixSeconds)
	if err != nil {
		return lifecyclestate.EffectPermit{}, err
	}
	err = store.commitLifecycleLocked(
		ctx,
		FaultLifecycleIntentAbort,
		FaultLifecycleIntentIndeterminatePreState,
		FaultLifecycleIntentIndeterminate,
		false,
		func(_ *installationState, records *[]lifecyclestate.Record) error {
			position := findLifecycle(*records, attemptID)
			if position < 0 || (*records)[position].View().RecordVersion != expected {
				return ErrLifecycleVersionConflict
			}
			(*records)[position] = updated
			return nil
		},
	)
	if err != nil {
		return lifecyclestate.EffectPermit{}, err
	}
	permit, err := store.permitForView(updated.View())
	if err != nil {
		store.repairRequired = true
		return lifecyclestate.EffectPermit{}, fmt.Errorf("%w: sealed permit", ErrStoreRepairRequired)
	}
	store.issuedPermits[attemptID] = permit
	return permit, nil
}

func (store *FixedFileStoreV1) ConfirmEffect(
	ctx context.Context,
	permit lifecyclestate.EffectPermit,
	result lifecyclestate.EffectResult,
) (lifecyclestate.Record, error) {
	if err := ctx.Err(); err != nil {
		return lifecyclestate.Record{}, err
	}
	if result.Status() == lifecyclestate.EffectResultIndeterminate {
		return lifecyclestate.Record{}, classified(ClassificationBinding, "indeterminate-result-requires-record-indeterminate")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.requireMutableLocked(); err != nil {
		return lifecyclestate.Record{}, err
	}
	view, position, err := store.validatePermitLocked(permit)
	if err != nil {
		if replay, ok := store.confirmReplayLocked(permit, result); ok {
			return replay, nil
		}
		return lifecyclestate.Record{}, err
	}
	if result.Operation() != view.Operation {
		return lifecyclestate.Record{}, ErrLifecycleBindingMismatch
	}
	if result.Status() == lifecyclestate.EffectResultApplied && result.Operation() == lifecyclestate.OperationCreate {
		for index, record := range store.lifecycles {
			if index != position && record.View().Instance.Present() &&
				record.View().Instance.Digest() == result.Instance().Digest() {
				return store.commitIdentityQuarantineLocked(ctx, position, view)
			}
		}
	}
	if result.Status() == lifecyclestate.EffectResultApplied {
		applyConfirmedEffect(&view, result.Instance())
	} else {
		applyNotAppliedDisposition(&view)
	}
	if err := advanceLifecycleVersion(&view, store.state.TimeHighWaterUnixSeconds); err != nil {
		return lifecyclestate.Record{}, err
	}
	updated, err := lifecyclestate.NewRecord(view, store.state.TimeHighWaterUnixSeconds)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	err = store.commitLifecycleLocked(
		ctx,
		FaultLifecycleAfterEffectAbort,
		FaultLifecycleAfterEffectIndeterminatePreState,
		FaultLifecycleAfterEffectIndeterminate,
		true,
		func(_ *installationState, records *[]lifecyclestate.Record) error {
			current := findLifecycle(*records, permit.View().AttemptID)
			if current < 0 || (*records)[current].View().RecordVersion != permit.View().RecordVersion {
				return ErrLifecycleVersionConflict
			}
			(*records)[current] = updated
			return nil
		},
	)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	delete(store.issuedPermits, permit.View().AttemptID)
	return cloneLifecycleRecord(updated, store.state.TimeHighWaterUnixSeconds), nil
}

func (store *FixedFileStoreV1) RecordIndeterminate(
	ctx context.Context,
	permit lifecyclestate.EffectPermit,
	classification Classification,
) (lifecyclestate.Record, error) {
	if err := ctx.Err(); err != nil {
		return lifecyclestate.Record{}, err
	}
	failure, ok := lifecycleFailureClassification(classification)
	if !ok {
		return lifecyclestate.Record{}, classified(ClassificationBinding, "indeterminate-classification")
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.requireMutableLocked(); err != nil {
		return lifecyclestate.Record{}, err
	}
	view, _, err := store.validatePermitLocked(permit)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	view.State = lifecyclestate.StateUnresolved
	view.EffectStatus = lifecyclestate.EffectIndeterminate
	view.RecoveryFence = lifecyclestate.RecoveryFenceReconcileUnknown
	setFirstFailure(&view, failure)
	if err := advanceLifecycleVersion(&view, store.state.TimeHighWaterUnixSeconds); err != nil {
		return lifecyclestate.Record{}, err
	}
	updated, err := lifecyclestate.NewRecord(view, store.state.TimeHighWaterUnixSeconds)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	err = store.commitLifecycleLocked(
		ctx,
		FaultLifecycleRecordIndeterminateAbort,
		FaultLifecycleRecordIndeterminatePreState,
		FaultLifecycleRecordIndeterminateCommitIndeterminate,
		true,
		func(_ *installationState, records *[]lifecyclestate.Record) error {
			position := findLifecycle(*records, permit.View().AttemptID)
			if position < 0 || (*records)[position].View().RecordVersion != permit.View().RecordVersion {
				return ErrLifecycleVersionConflict
			}
			(*records)[position] = updated
			return nil
		},
	)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	delete(store.issuedPermits, permit.View().AttemptID)
	return cloneLifecycleRecord(updated, store.state.TimeHighWaterUnixSeconds), nil
}

func (store *FixedFileStoreV1) BeginReconciliation(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	expected lifecyclestate.RecordVersion,
) (lifecyclestate.Record, error) {
	if err := ctx.Err(); err != nil {
		return lifecyclestate.Record{}, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.requireMutableLocked(); err != nil {
		return lifecyclestate.Record{}, err
	}
	position := findLifecycle(store.lifecycles, attemptID)
	if position < 0 {
		return lifecyclestate.Record{}, ErrLifecycleNotFound
	}
	view := store.lifecycles[position].View()
	if view.RecordVersion != expected {
		return lifecyclestate.Record{}, ErrLifecycleVersionConflict
	}
	if view.EffectStatus != lifecyclestate.EffectIntent &&
		view.EffectStatus != lifecyclestate.EffectIndeterminate &&
		view.State != lifecyclestate.StateDestroyConfirmed {
		return lifecyclestate.Record{}, ErrLifecycleBindingMismatch
	}
	if view.AutomaticRecoveryCount >= lifecyclestate.MaxAutomaticRecoveryAttempts {
		return lifecyclestate.Record{}, classified(ClassificationCapacity, "automatic-recovery-exhausted")
	}
	if view.NextRecoveryAt.Present && store.state.TimeHighWaterUnixSeconds < view.NextRecoveryAt.Value {
		return lifecyclestate.Record{}, classified(ClassificationStale, "reconciliation-not-eligible")
	}

	// Commit the eligibility/count checkpoint separately so E4 places the
	// observation strictly after this durable edge.
	eligibility := view
	eligibility.AutomaticRecoveryCount++
	eligibility.NextRecoveryAt = nextRecoveryTime(eligibility.AutomaticRecoveryCount, store.state.TimeHighWaterUnixSeconds)
	eligibility.LastReconciliation = lifecyclestate.ReconciliationNone
	if err := advanceLifecycleVersion(&eligibility, store.state.TimeHighWaterUnixSeconds); err != nil {
		return lifecyclestate.Record{}, err
	}
	eligibleRecord, err := lifecyclestate.NewRecord(eligibility, store.state.TimeHighWaterUnixSeconds)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	err = store.commitLifecycleLocked(
		ctx,
		FaultLifecycleReconciliationEligibilityAbort,
		FaultLifecycleReconciliationEligibilityIndeterminatePreState,
		FaultLifecycleReconciliationEligibilityIndeterminate,
		false,
		func(_ *installationState, records *[]lifecyclestate.Record) error {
			current := findLifecycle(*records, attemptID)
			if current < 0 || (*records)[current].View().RecordVersion != expected {
				return ErrLifecycleVersionConflict
			}
			(*records)[current] = eligibleRecord
			return nil
		},
	)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	return cloneLifecycleRecord(eligibleRecord, store.state.TimeHighWaterUnixSeconds), nil
}

func (store *FixedFileStoreV1) CompleteReconciliation(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	expected lifecyclestate.RecordVersion,
	result lifecyclestate.ReconcileResult,
) (lifecyclestate.Record, error) {
	if err := ctx.Err(); err != nil {
		return lifecyclestate.Record{}, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.requireMutableLocked(); err != nil {
		return lifecyclestate.Record{}, err
	}
	position := findLifecycle(store.lifecycles, attemptID)
	if position < 0 {
		return lifecyclestate.Record{}, ErrLifecycleNotFound
	}
	eligibleRecord := store.lifecycles[position]
	eligibleView := eligibleRecord.View()
	if eligibleView.RecordVersion != expected {
		return lifecyclestate.Record{}, ErrLifecycleVersionConflict
	}
	if eligibleView.AutomaticRecoveryCount == 0 ||
		eligibleView.LastReconciliation != lifecyclestate.ReconciliationNone ||
		result.Operation() != eligibleView.Operation ||
		(eligibleView.EffectStatus != lifecyclestate.EffectIntent &&
			eligibleView.EffectStatus != lifecyclestate.EffectIndeterminate &&
			eligibleView.State != lifecyclestate.StateDestroyConfirmed) {
		return lifecyclestate.Record{}, ErrLifecycleBindingMismatch
	}

	final := eligibleRecord.View()
	applyReconciliationResult(&final, result)
	if err := store.validateReconciledInstanceLocked(attemptID, &final, result); err != nil {
		if errors.Is(err, ErrLifecycleBindingMismatch) {
			final.State = lifecyclestate.StateQuarantined
			final.EffectStatus = lifecyclestate.EffectIndeterminate
			final.LastReconciliation = lifecyclestate.ReconciliationIdentityMismatch
			final.RecoveryFence = lifecyclestate.RecoveryFenceIdentityMismatch
			final.NextRecoveryAt = lifecyclestate.OptionalUnixSeconds{}
			setFirstFailure(&final, lifecyclestate.FailureBinding)
		} else {
			return lifecyclestate.Record{}, err
		}
	}
	if err := advanceLifecycleVersion(&final, store.state.TimeHighWaterUnixSeconds); err != nil {
		return lifecyclestate.Record{}, err
	}
	finalRecord, err := lifecyclestate.NewRecord(final, store.state.TimeHighWaterUnixSeconds)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	err = store.commitLifecycleLocked(
		ctx,
		FaultLifecycleReconciliationResultAbort,
		FaultLifecycleReconciliationResultIndeterminatePreState,
		FaultLifecycleReconciliationResultIndeterminate,
		true,
		func(state *installationState, records *[]lifecyclestate.Record) error {
			current := findLifecycle(*records, attemptID)
			if current < 0 || (*records)[current].View().RecordVersion != eligibleRecord.View().RecordVersion {
				return ErrLifecycleVersionConflict
			}
			if final.State == lifecyclestate.StateQuarantined {
				state.TrustPhase = TrustRepairRequired
				state.TrustReason = "lifecycle-instance-identity-mismatch"
				state.AttemptsDisabled = true
			}
			(*records)[current] = finalRecord
			return nil
		},
	)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	delete(store.issuedPermits, attemptID)
	return cloneLifecycleRecord(finalRecord, store.state.TimeHighWaterUnixSeconds), nil
}

// RecordReconciliation retains the E3 composite transaction oracle. The E4
// driver uses BeginReconciliation and CompleteReconciliation so the fake
// observation occurs strictly between the two durable commits.
func (store *FixedFileStoreV1) RecordReconciliation(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
	expected lifecyclestate.RecordVersion,
	result lifecyclestate.ReconcileResult,
) (lifecyclestate.Record, error) {
	eligible, err := store.BeginReconciliation(ctx, attemptID, expected)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	return store.CompleteReconciliation(ctx, attemptID, eligible.View().RecordVersion, result)
}

func (store *FixedFileStoreV1) RecoveryAttemptIDs(ctx context.Context) ([]approvalattempt.AttemptID, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if store.repairRequired {
		return nil, ErrStoreRepairRequired
	}
	result := make([]approvalattempt.AttemptID, 0, len(store.state.Attempts))
	for _, attempt := range store.state.Attempts {
		position := findLifecycle(store.lifecycles, attempt.AttemptID)
		if position < 0 {
			result = append(result, attempt.AttemptID)
			continue
		}
		view := store.lifecycles[position].View()
		if view.State != lifecyclestate.StateDestroyed || view.CleanupRequired ||
			view.State == lifecyclestate.StateUnresolved || view.State == lifecyclestate.StateQuarantined ||
			view.RecoveryFence != lifecyclestate.RecoveryFenceNone || view.NextRecoveryAt.Present {
			result = append(result, attempt.AttemptID)
		}
	}
	sort.Slice(result, func(left, right int) bool {
		return bytes.Compare(result[left][:], result[right][:]) < 0
	})
	return append([]approvalattempt.AttemptID(nil), result...), nil
}

type lifecycleMutation func(*installationState, *[]lifecyclestate.Record) error

func (store *FixedFileStoreV1) commitLifecycleLocked(
	ctx context.Context,
	abortPoint LifecycleStoreFault,
	preIndeterminatePoint LifecycleStoreFault,
	postRenamePoint LifecycleStoreFault,
	fenceOnAbort bool,
	mutation lifecycleMutation,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := store.requireMutableLocked(); err != nil {
		return err
	}
	candidateState := cloneState(store.state)
	candidateRecords := cloneLifecycleRecords(store.lifecycles, store.state.TimeHighWaterUnixSeconds)
	if err := mutation(&candidateState, &candidateRecords); err != nil {
		if errors.Is(err, errNoStoreMutation) {
			return nil
		}
		return err
	}
	sortLifecycleRecords(candidateRecords)
	candidateDigest := lifecycleSetDigest(candidateRecords)
	if err := validateV1State(candidateState, candidateRecords, candidateDigest); err != nil {
		return err
	}
	abort := store.takeLifecycleFaultLocked(abortPoint)
	if err := store.takeLifecycleFaultLocked(preIndeterminatePoint); err != nil {
		store.recoveryRequired = true
		return fmt.Errorf("%w: injected before lifecycle commit", ErrCommitOutcomeIndeterminate)
	}
	postRename := store.takeLifecycleFaultLocked(postRenamePoint)
	envelope := encodedEnvelopeV1(candidateState, candidateRecords)
	if err := persistLifecycleEnvelopeV1(store.path, envelope, abort, postRename); err != nil {
		if abort != nil && fenceOnAbort {
			store.recoveryRequired = true
		}
		if errors.Is(err, ErrCommitOutcomeIndeterminate) {
			store.state = candidateState
			store.lifecycles = candidateRecords
			store.lifecycleSetDigest = candidateDigest
			store.recoveryRequired = true
		}
		return err
	}
	store.state = candidateState
	store.lifecycles = candidateRecords
	store.lifecycleSetDigest = candidateDigest
	return nil
}

func persistLifecycleEnvelopeV1(path string, envelope diskEnvelopeV1, beforeRename, afterRename error) error {
	parent := filepath.Dir(path)
	temporary, err := os.CreateTemp(parent, ".capsule-supervisor-v1-transaction-*.tmp")
	if err != nil {
		return errors.New("create fixed supervisor v1 transaction")
	}
	temporaryPath := temporary.Name()
	committed := false
	defer func() {
		_ = temporary.Close()
		if !committed {
			_ = os.Remove(temporaryPath)
		}
	}()
	if err := temporary.Chmod(0o600); err != nil {
		return errors.New("protect fixed supervisor v1 transaction")
	}
	encoder := json.NewEncoder(temporary)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(envelope); err != nil {
		return errors.New("encode fixed supervisor v1 transaction")
	}
	if err := temporary.Sync(); err != nil {
		return errors.New("sync fixed supervisor v1 transaction")
	}
	if err := temporary.Close(); err != nil {
		return errors.New("close fixed supervisor v1 transaction")
	}
	if beforeRename != nil {
		return errors.New("fixed supervisor v1 transaction confirmed before rename failure")
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return errors.New("commit fixed supervisor v1 transaction")
	}
	committed = true
	if afterRename != nil {
		return fmt.Errorf("%w: injected after lifecycle rename", ErrCommitOutcomeIndeterminate)
	}
	directory, err := os.Open(parent) //nolint:gosec // G304: parent is the store's own installation-configured directory, not caller/network-supplied
	if err != nil {
		return fmt.Errorf("%w: open lifecycle store directory", ErrCommitOutcomeIndeterminate)
	}
	syncErr := directory.Sync()
	closeErr := directory.Close()
	if syncErr != nil || closeErr != nil {
		return fmt.Errorf("%w: sync lifecycle store directory", ErrCommitOutcomeIndeterminate)
	}
	return nil
}

func (store *FixedFileStoreV1) requireMutableLocked() error {
	if store.repairRequired {
		return ErrStoreRepairRequired
	}
	if store.recoveryRequired {
		return ErrCommitOutcomeIndeterminate
	}
	return nil
}

func (store *FixedFileStoreV1) takeLifecycleFaultLocked(point LifecycleStoreFault) error {
	if point == "" {
		return nil
	}
	err := store.faults[point]
	delete(store.faults, point)
	return err
}

func lifecycleBindingsFromAuthority(
	state installationState,
	attemptID approvalattempt.AttemptID,
	backend lifecyclestate.BackendBinding,
) (lifecyclestate.ImmutableBindings, error) {
	if attemptID == (approvalattempt.AttemptID{}) {
		return lifecyclestate.ImmutableBindings{}, classified(ClassificationSchema, "attempt-id-zero")
	}
	attemptIndex := -1
	for index, attempt := range state.Attempts {
		if attempt.AttemptID == attemptID {
			attemptIndex = index
			break
		}
	}
	if attemptIndex < 0 {
		return lifecyclestate.ImmutableBindings{}, ErrLifecycleNotFound
	}
	attempt := state.Attempts[attemptIndex]
	approvalIndex := -1
	for index, approval := range state.Approvals {
		if approval.ApprovalID == attempt.ApprovalID {
			approvalIndex = index
			break
		}
	}
	if approvalIndex < 0 || state.Approvals[approvalIndex].State != approvalattempt.ApprovalConsumed ||
		state.Approvals[approvalIndex].ConsumedAttemptID != attemptID {
		return lifecyclestate.ImmutableBindings{}, fmt.Errorf("%w: approval cross-link", ErrStoreRepairRequired)
	}
	registrationIndex := findRegistration(state.Registrations, attempt.RegistrationID)
	if registrationIndex < 0 {
		return lifecyclestate.ImmutableBindings{}, fmt.Errorf("%w: registration cross-link", ErrStoreRepairRequired)
	}
	approval := state.Approvals[approvalIndex]
	registration := state.Registrations[registrationIndex]
	backendView := backend.View()
	if backendView.BackendConfigurationDigest != registration.PlanBindings.BackendConfigurationDigest ||
		backendView.BackendValidationRecordDigest != registration.PlanBindings.BackendValidationRecordDigest {
		return lifecyclestate.ImmutableBindings{}, ErrLifecycleBindingMismatch
	}
	return lifecyclestate.NewImmutableBindings(lifecyclestate.ImmutableBindingsView{
		AttemptID:                       attempt.AttemptID,
		ApprovalID:                      attempt.ApprovalID,
		AttemptNonce:                    attempt.AttemptNonce,
		RegistrationID:                  attempt.RegistrationID,
		RegistrationSequence:            attempt.RegistrationSequence,
		PlanDigest:                      attempt.PlanDigest,
		InstallationID:                  attempt.InstallationID,
		EpochSequence:                   attempt.EpochSequence,
		EpochDigest:                     attempt.EpochDigest,
		SupervisorID:                    attempt.SupervisorID,
		ApprovalPurpose:                 attempt.Purpose,
		ApprovalAudience:                attempt.Audience,
		ApprovalPayloadDigest:           attempt.ApprovalPayloadDigest,
		AuthorizationIdentity:           attempt.AuthorizationIdentity,
		AttemptCreatedAt:                attempt.CreatedAt,
		AttemptStorageVersion:           attempt.StorageFormatVersion,
		ApprovalStorageVersion:          approval.StorageFormatVersion,
		RegistrationStorageVersion:      registration.Record.StorageFormatVersion,
		SourceManifestDigest:            registration.PlanBindings.SourceManifestDigest,
		InlineInputDigest:               registration.PlanBindings.InlineInputDigest,
		RuntimeBundleManifestDigest:     registration.PlanBindings.RuntimeBundleManifestDigest,
		ProfileReviewAttestationDigests: registration.PlanBindings.ProfileReviewAttestationDigests,
		ProfileRegistryEntryDigest:      registration.PlanBindings.ProfileRegistryEntryDigest,
		BackendValidationRecordDigest:   registration.PlanBindings.BackendValidationRecordDigest,
		BackendConfigurationDigest:      registration.PlanBindings.BackendConfigurationDigest,
		TrustSnapshotDigest:             registration.PlanBindings.TrustSnapshotDigest,
		PolicyDecisionDigest:            registration.PlanBindings.PolicyDecisionDigest,
		Backend:                         backend,
	})
}

func findLifecycle(records []lifecyclestate.Record, attemptID approvalattempt.AttemptID) int {
	position := sort.Search(len(records), func(index int) bool {
		candidate := records[index].Bindings().View().AttemptID
		return bytes.Compare(candidate[:], attemptID[:]) >= 0
	})
	if position < len(records) && records[position].Bindings().View().AttemptID == attemptID {
		return position
	}
	return -1
}

func sortLifecycleRecords(records []lifecyclestate.Record) {
	sort.Slice(records, func(left, right int) bool {
		leftID := records[left].Bindings().View().AttemptID
		rightID := records[right].Bindings().View().AttemptID
		return bytes.Compare(leftID[:], rightID[:]) < 0
	})
}

func activeLifecycleWork(state installationState, records []lifecyclestate.Record) int {
	active := len(state.Attempts)
	for _, record := range records {
		view := record.View()
		if view.State == lifecyclestate.StateDestroyed && !view.CleanupRequired {
			active--
		}
	}
	return active
}

func cloneLifecycleRecord(record lifecyclestate.Record, highWater v0candidate.UInt53) lifecyclestate.Record {
	cloned, err := restoreLifecycleRecord(lifecycleRecordToDisk(record), highWater)
	if err != nil {
		panic(fmt.Sprintf("clone validated lifecycle record: %v", err))
	}
	return cloned
}

func validateEffectPredecessor(view lifecyclestate.RecordView, operation lifecyclestate.Operation) error {
	valid := false
	switch operation {
	case lifecyclestate.OperationPrepare:
		valid = view.State == lifecyclestate.StatePreparePending
	case lifecyclestate.OperationCreate:
		valid = view.State == lifecyclestate.StatePrepared
	case lifecyclestate.OperationStart:
		valid = view.State == lifecyclestate.StateCreated
	case lifecyclestate.OperationObserve:
		valid = view.State == lifecyclestate.StateStarted
	case lifecyclestate.OperationStop:
		valid = view.State == lifecyclestate.StateCreated || view.State == lifecyclestate.StateStarted || view.State == lifecyclestate.StateObserved
	case lifecyclestate.OperationDestroy:
		valid = view.State == lifecyclestate.StateCreated || view.State == lifecyclestate.StateStarted ||
			view.State == lifecyclestate.StateObserved || view.State == lifecyclestate.StateStopped
	}
	if !valid {
		return ErrLifecycleBindingMismatch
	}
	return nil
}

func intentState(operation lifecyclestate.Operation) lifecyclestate.LifecycleState {
	switch operation {
	case lifecyclestate.OperationPrepare:
		return lifecyclestate.StatePrepareIntent
	case lifecyclestate.OperationCreate:
		return lifecyclestate.StateCreateIntent
	case lifecyclestate.OperationStart:
		return lifecyclestate.StateStartIntent
	case lifecyclestate.OperationObserve:
		return lifecyclestate.StateObserveIntent
	case lifecyclestate.OperationStop:
		return lifecyclestate.StateStopIntent
	case lifecyclestate.OperationDestroy:
		return lifecyclestate.StateDestroyIntent
	default:
		return ""
	}
}

func retryableSameEffect(view lifecyclestate.RecordView, expected lifecyclestate.RecordVersion, operation lifecyclestate.Operation) bool {
	if view.RecordVersion != expected || view.Operation != operation || view.EffectID.IsZero() ||
		view.State != lifecyclestate.StateUnresolved {
		return false
	}
	if view.LastReconciliation == lifecyclestate.ReconciliationNotApplied &&
		(view.EffectStatus == lifecyclestate.EffectIntent || view.EffectStatus == lifecyclestate.EffectIndeterminate) {
		return true
	}
	return view.EffectStatus == lifecyclestate.EffectConfirmed &&
		view.FirstFailure == lifecyclestate.FailureLifecycle &&
		view.RecoveryFence == lifecyclestate.RecoveryFenceNone
}

func (store *FixedFileStoreV1) permitForView(view lifecyclestate.RecordView) (lifecyclestate.EffectPermit, error) {
	return lifecyclestate.NewEffectPermit(lifecyclestate.EffectPermitView{
		AttemptID:              view.Bindings.View().AttemptID,
		RecordVersion:          view.RecordVersion,
		OperationSequence:      view.OperationSequence,
		Operation:              view.Operation,
		EffectID:               view.EffectID,
		ImmutableBindingDigest: view.ImmutableBindingDigest,
		OwnerSessionID:         store.ownerSessionID,
		Instance:               view.Instance,
	})
}

func (store *FixedFileStoreV1) validatePermitLocked(
	permit lifecyclestate.EffectPermit,
) (lifecyclestate.RecordView, int, error) {
	permitView := permit.View()
	if permitView.OwnerSessionID != store.ownerSessionID {
		return lifecyclestate.RecordView{}, -1, ErrLifecycleBindingMismatch
	}
	position := findLifecycle(store.lifecycles, permitView.AttemptID)
	if position < 0 {
		return lifecyclestate.RecordView{}, -1, ErrLifecycleNotFound
	}
	view := store.lifecycles[position].View()
	retryPermit := view.State == lifecyclestate.StateUnresolved &&
		view.EffectStatus == lifecyclestate.EffectConfirmed &&
		view.FirstFailure == lifecyclestate.FailureLifecycle &&
		view.RecoveryFence == lifecyclestate.RecoveryFenceNone
	if view.RecordVersion != permitView.RecordVersion || view.OperationSequence != permitView.OperationSequence ||
		view.Operation != permitView.Operation || view.EffectID != permitView.EffectID ||
		view.ImmutableBindingDigest != permitView.ImmutableBindingDigest ||
		!backendInstancesEqual(view.Instance, permitView.Instance) ||
		(view.EffectStatus != lifecyclestate.EffectIntent &&
			view.EffectStatus != lifecyclestate.EffectIndeterminate &&
			!retryPermit) {
		return lifecyclestate.RecordView{}, -1, ErrLifecycleBindingMismatch
	}
	return view, position, nil
}

func (store *FixedFileStoreV1) confirmReplayLocked(
	permit lifecyclestate.EffectPermit,
	result lifecyclestate.EffectResult,
) (lifecyclestate.Record, bool) {
	permitView := permit.View()
	if permitView.OwnerSessionID != store.ownerSessionID || result.Operation() != permitView.Operation {
		return lifecyclestate.Record{}, false
	}
	position := findLifecycle(store.lifecycles, permitView.AttemptID)
	if position < 0 {
		return lifecyclestate.Record{}, false
	}
	record := store.lifecycles[position]
	view := record.View()
	if uint64(view.RecordVersion) != uint64(permitView.RecordVersion)+1 ||
		view.OperationSequence != permitView.OperationSequence || view.Operation != permitView.Operation ||
		view.EffectID != permitView.EffectID || view.ImmutableBindingDigest != permitView.ImmutableBindingDigest {
		return lifecyclestate.Record{}, false
	}
	if result.Status() == lifecyclestate.EffectResultApplied {
		if view.EffectStatus != lifecyclestate.EffectConfirmed ||
			(result.Operation() == lifecyclestate.OperationCreate && !backendInstancesEqual(view.Instance, result.Instance())) {
			return lifecyclestate.Record{}, false
		}
	} else if view.State != lifecyclestate.StateUnresolved || view.FirstFailure == lifecyclestate.FailureNone {
		return lifecyclestate.Record{}, false
	}
	return cloneLifecycleRecord(record, store.state.TimeHighWaterUnixSeconds), true
}

func applyConfirmedEffect(view *lifecyclestate.RecordView, instance lifecyclestate.BackendInstanceIdentity) {
	view.EffectStatus = lifecyclestate.EffectConfirmed
	view.RecoveryFence = lifecyclestate.RecoveryFenceNone
	view.NextRecoveryAt = lifecyclestate.OptionalUnixSeconds{}
	switch view.Operation {
	case lifecyclestate.OperationPrepare:
		view.State = lifecyclestate.StatePrepared
		view.LastConfirmedCheckpoint = lifecyclestate.CheckpointPrepare
	case lifecyclestate.OperationCreate:
		view.State = lifecyclestate.StateCreated
		view.LastConfirmedCheckpoint = lifecyclestate.CheckpointCreate
		view.Instance = instance
	case lifecyclestate.OperationStart:
		view.State = lifecyclestate.StateStarted
		view.LastConfirmedCheckpoint = lifecyclestate.CheckpointStart
	case lifecyclestate.OperationObserve:
		view.State = lifecyclestate.StateObserved
		view.LastConfirmedCheckpoint = lifecyclestate.CheckpointObserve
	case lifecyclestate.OperationStop:
		view.State = lifecyclestate.StateStopped
		view.LastConfirmedCheckpoint = lifecyclestate.CheckpointStop
	case lifecyclestate.OperationDestroy:
		view.State = lifecyclestate.StateDestroyConfirmed
		view.LastConfirmedCheckpoint = lifecyclestate.CheckpointDestroy
	}
}

func applyNotAppliedDisposition(view *lifecyclestate.RecordView) {
	view.State = lifecyclestate.StateUnresolved
	view.EffectStatus = lifecyclestate.EffectConfirmed
	view.RecoveryFence = lifecyclestate.RecoveryFenceNone
	setFirstFailure(view, lifecyclestate.FailureLifecycle)
}

func setFirstFailure(view *lifecyclestate.RecordView, failure lifecyclestate.FailureClassification) {
	if view.FirstFailure == lifecyclestate.FailureNone {
		view.FirstFailure = failure
		view.FailureOperation = view.Operation
	}
}

func advanceLifecycleVersion(view *lifecyclestate.RecordView, highWater v0candidate.UInt53) error {
	if uint64(view.RecordVersion) >= v0candidate.MaxSafeInteger || uint64(view.SnapshotGeneration) >= v0candidate.MaxSafeInteger {
		return classified(ClassificationCapacity, "lifecycle-version-generation")
	}
	if highWater < view.LastTransitionAt {
		return classified(ClassificationStale, "lifecycle-time-rollback")
	}
	view.RecordVersion++
	view.SnapshotGeneration++
	view.LastTransitionAt = highWater
	return nil
}

func lifecycleFailureClassification(classification Classification) (lifecyclestate.FailureClassification, bool) {
	switch classification {
	case ClassificationMalformed:
		return lifecyclestate.FailureMalformed, true
	case ClassificationUnsupported:
		return lifecyclestate.FailureUnsupported, true
	case ClassificationSchema:
		return lifecyclestate.FailureSchema, true
	case ClassificationDomain:
		return lifecyclestate.FailureDomain, true
	case ClassificationBinding:
		return lifecyclestate.FailureBinding, true
	case ClassificationStale:
		return lifecyclestate.FailureStale, true
	case ClassificationCapacity:
		return lifecyclestate.FailureCapacity, true
	case ClassificationLocalFailure:
		return lifecyclestate.FailureLocal, true
	case ClassificationLifecycleFailure:
		return lifecyclestate.FailureLifecycle, true
	case ClassificationCleanupUnresolved:
		return lifecyclestate.FailureCleanupUnresolved, true
	case ClassificationTrustState:
		return lifecyclestate.FailureTrustState, true
	case ClassificationRecoveryRequired:
		return lifecyclestate.FailureRecoveryRequired, true
	default:
		return lifecyclestate.FailureNone, false
	}
}

func nextRecoveryTime(count v0candidate.UInt53, highWater v0candidate.UInt53) lifecyclestate.OptionalUnixSeconds {
	delta := uint64(0)
	switch count {
	case 1:
		delta = 1
	case 2:
		delta = 5
	default:
		return lifecyclestate.OptionalUnixSeconds{}
	}
	if uint64(highWater) > v0candidate.MaxSafeInteger-delta {
		return lifecyclestate.OptionalUnixSeconds{Present: true, Value: v0candidate.UInt53(v0candidate.MaxSafeInteger)}
	}
	return lifecyclestate.OptionalUnixSeconds{Present: true, Value: highWater + v0candidate.UInt53(delta)}
}

func applyReconciliationResult(view *lifecyclestate.RecordView, result lifecyclestate.ReconcileResult) {
	view.LastReconciliation = result.Status()
	switch result.Status() {
	case lifecyclestate.ReconciliationApplied:
		applyConfirmedEffect(view, result.Instance())
	case lifecyclestate.ReconciliationNotApplied:
		view.State = lifecyclestate.StateUnresolved
		view.RecoveryFence = lifecyclestate.RecoveryFenceNone
		view.NextRecoveryAt = lifecyclestate.OptionalUnixSeconds{}
		setFirstFailure(view, lifecyclestate.FailureLifecycle)
	case lifecyclestate.ReconciliationUnknown:
		view.State = lifecyclestate.StateUnresolved
		view.EffectStatus = lifecyclestate.EffectIndeterminate
		view.RecoveryFence = lifecyclestate.RecoveryFenceReconcileUnknown
		setFirstFailure(view, lifecyclestate.FailureCleanupUnresolved)
		if view.AutomaticRecoveryCount == lifecyclestate.MaxAutomaticRecoveryAttempts {
			view.AutomaticRecoveryExhausted = true
			view.RecoveryFence = lifecyclestate.RecoveryFenceAutomaticExhausted
			view.NextRecoveryAt = lifecyclestate.OptionalUnixSeconds{}
		}
	case lifecyclestate.ReconciliationPresent:
		view.State = lifecyclestate.StateUnresolved
		view.EffectStatus = lifecyclestate.EffectIndeterminate
		view.RecoveryFence = lifecyclestate.RecoveryFenceReconcileUnknown
		setFirstFailure(view, lifecyclestate.FailureCleanupUnresolved)
	case lifecyclestate.ReconciliationAuthoritativelyAbsent:
		if view.State == lifecyclestate.StateDestroyConfirmed ||
			(view.LastConfirmedCheckpoint == lifecyclestate.CheckpointNone && view.Operation == lifecyclestate.OperationPrepare) {
			view.State = lifecyclestate.StateDestroyed
			view.CleanupRequired = false
			view.TerminalAt = lifecyclestate.OptionalUnixSeconds{Present: true, Value: view.LastTransitionAt}
			view.RecoveryFence = lifecyclestate.RecoveryFenceNone
			view.NextRecoveryAt = lifecyclestate.OptionalUnixSeconds{}
		} else {
			view.State = lifecyclestate.StateUnresolved
			view.EffectStatus = lifecyclestate.EffectIndeterminate
			view.RecoveryFence = lifecyclestate.RecoveryFenceReconcileUnknown
			setFirstFailure(view, lifecyclestate.FailureCleanupUnresolved)
		}
	case lifecyclestate.ReconciliationIdentityMismatch:
		view.State = lifecyclestate.StateQuarantined
		view.EffectStatus = lifecyclestate.EffectIndeterminate
		view.RecoveryFence = lifecyclestate.RecoveryFenceIdentityMismatch
		view.NextRecoveryAt = lifecyclestate.OptionalUnixSeconds{}
		setFirstFailure(view, lifecyclestate.FailureBinding)
	}
}

func (store *FixedFileStoreV1) validateReconciledInstanceLocked(
	attemptID approvalattempt.AttemptID,
	view *lifecyclestate.RecordView,
	result lifecyclestate.ReconcileResult,
) error {
	instance := result.Instance()
	if result.Status() == lifecyclestate.ReconciliationApplied && result.Operation() == lifecyclestate.OperationCreate {
		if !instance.Present() {
			return ErrLifecycleBindingMismatch
		}
		for _, record := range store.lifecycles {
			other := record.View()
			if other.Bindings.View().AttemptID != attemptID && other.Instance.Present() && other.Instance.Digest() == instance.Digest() {
				return ErrLifecycleBindingMismatch
			}
		}
		return nil
	}
	if result.Status() == lifecyclestate.ReconciliationPresent || result.Status() == lifecyclestate.ReconciliationIdentityMismatch {
		if !view.Instance.Present() || !backendInstancesEqual(view.Instance, instance) {
			return ErrLifecycleBindingMismatch
		}
	}
	return nil
}

func backendInstancesEqual(left, right lifecyclestate.BackendInstanceIdentity) bool {
	if left.Present() != right.Present() || left.Kind() != right.Kind() || left.Digest() != right.Digest() {
		return false
	}
	return bytes.Equal(left.Value(), right.Value())
}

func (store *FixedFileStoreV1) commitIdentityQuarantineLocked(
	ctx context.Context,
	position int,
	view lifecyclestate.RecordView,
) (lifecyclestate.Record, error) {
	view.State = lifecyclestate.StateQuarantined
	view.EffectStatus = lifecyclestate.EffectIndeterminate
	view.RecoveryFence = lifecyclestate.RecoveryFenceIdentityMismatch
	setFirstFailure(&view, lifecyclestate.FailureBinding)
	if err := advanceLifecycleVersion(&view, store.state.TimeHighWaterUnixSeconds); err != nil {
		return lifecyclestate.Record{}, err
	}
	record, err := lifecyclestate.NewRecord(view, store.state.TimeHighWaterUnixSeconds)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	attemptID := view.Bindings.View().AttemptID
	err = store.commitLifecycleLocked(
		ctx,
		FaultLifecycleAfterEffectAbort,
		FaultLifecycleAfterEffectIndeterminatePreState,
		FaultLifecycleAfterEffectIndeterminate,
		true,
		func(state *installationState, records *[]lifecyclestate.Record) error {
			current := findLifecycle(*records, attemptID)
			if current != position {
				return ErrLifecycleVersionConflict
			}
			state.TrustPhase = TrustRepairRequired
			state.TrustReason = "lifecycle-instance-identity-collision"
			state.AttemptsDisabled = true
			(*records)[current] = record
			return nil
		},
	)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	delete(store.issuedPermits, attemptID)
	return cloneLifecycleRecord(record, store.state.TimeHighWaterUnixSeconds), ErrLifecycleBindingMismatch
}
