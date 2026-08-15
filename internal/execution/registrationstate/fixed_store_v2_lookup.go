package registrationstate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/archivestate"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// RetainedRegistration is the location-independent closed registration
// projection returned by F4A. Hot and archived records with the same retained
// semantics produce the same value.
type RetainedRegistration struct {
	RegistrationID       v0candidate.RegistrationID
	RegistrationSequence v0candidate.PositiveUInt53
	ExpiresAt            v0candidate.UInt53
	PlanRoleBindings     v0candidate.ExecutionPlanRoleBindings
	Record               StoredRegistration
}

// RetainedAttempt is the immutable attempt and its exact optional lifecycle
// projection. LifecyclePresent distinguishes the ordinary hot pre-lifecycle
// state from a zero value; an archived attempt always has a present terminal
// lifecycle because F3 archives only complete closed cohorts.
type RetainedAttempt struct {
	Attempt          approvalattempt.ExecutionAttempt
	LifecyclePresent bool
	Lifecycle        lifecyclestate.Record
}

// RetainedNonce binds one retained nonce tombstone to its full approval.
type RetainedNonce struct {
	AttemptNonce  approvalattempt.AttemptNonce
	PayloadDigest approvalattempt.ApprovalPayloadDigest
	Approval      approvalattempt.ApprovalRecord
}

// RetainedEffectClassification distinguishes the lifecycle record's current
// effect from an independently retained historical tombstone. It prevents a
// later lifecycle record from being presented as the issuer of an older effect.
type RetainedEffectClassification string

const (
	// RetainedEffectCurrent includes the matching live lifecycle projection;
	// RetainedEffectSupersededByCurrent returns only immutable tombstone facts.
	// Neither result proves an adapter effect occurred or grants retry authority.
	RetainedEffectCurrent RetainedEffectClassification = "current"
	// RetainedEffectSupersededByCurrent returns tombstone facts without a false lifecycle issuer.
	RetainedEffectSupersededByCurrent RetainedEffectClassification = "superseded-by-current"
)

// RetainedEffect binds one visible effect tombstone to its own immutable
// issuance facts. LifecyclePresent is true only while the tombstone still
// names the lifecycle record's current effect; historical resolution never
// mislabels the later lifecycle record as the issuing record.
type RetainedEffect struct {
	EffectID                   lifecyclestate.EffectID
	AttemptID                  approvalattempt.AttemptID
	OperationSequence          v0candidate.PositiveUInt53
	Operation                  lifecyclestate.Operation
	IssuanceSnapshotGeneration v0candidate.PositiveUInt53
	VisibleV1Seed              bool
	Classification             RetainedEffectClassification
	Attempt                    approvalattempt.ExecutionAttempt
	LifecyclePresent           bool
	Lifecycle                  lifecyclestate.Record
}

// RetainedInstance binds one instance tombstone to its full lifecycle record.
type RetainedInstance struct {
	InstanceDigest lifecyclestate.BackendInstanceDigest
	AttemptID      approvalattempt.AttemptID
	Lifecycle      lifecyclestate.Record
}

// RetainedIdentityCandidates is the passive F4A uniqueness query. Zero scalar
// identities mean "not queried". Replay identities are closed pairs: setting
// only one arm is a schema error. F4A never reserves or writes a candidate.
type RetainedIdentityCandidates struct {
	RegistrationID v0candidate.RegistrationID
	ApprovalID     approvalattempt.ApprovalID
	AttemptID      approvalattempt.AttemptID
	AttemptNonce   approvalattempt.AttemptNonce
	EffectID       lifecyclestate.EffectID
	InstanceDigest lifecyclestate.BackendInstanceDigest

	ApprovalReplayPayloadDigest         approvalattempt.ApprovalPayloadDigest
	ApprovalReplayAuthorizationIdentity approvalattempt.ApprovalKeyAuthorizationIdentity
	AttemptReplayRegistrationID         v0candidate.RegistrationID
	AttemptReplayApprovalID             approvalattempt.ApprovalID
}

// RetainedIdentityCollisions reports which candidate identities already exist
// in the fully verified retained-global indexes. It contains no record bytes or
// locations and grants no authority.
type RetainedIdentityCollisions struct {
	Registration   bool
	Approval       bool
	Attempt        bool
	Nonce          bool
	Effect         bool
	Instance       bool
	ApprovalReplay bool
	AttemptReplay  bool
}

// Any reports whether at least one candidate identity or replay key already
// exists in the fully verified retained-global world. It is a passive query,
// not a reservation; mutation must repeat collision checks transactionally.
func (collisions RetainedIdentityCollisions) Any() bool {
	return collisions.Registration || collisions.Approval || collisions.Attempt || collisions.Nonce ||
		collisions.Effect || collisions.Instance || collisions.ApprovalReplay || collisions.AttemptReplay
}

// ResolveRegistration resolves one retained registration by its nominal ID.
// It always reopens and fully validates the active snapshot and every
// referenced segment before following the retained-global typed location.
func (store *FixedFileStoreV2) ResolveRegistration(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
) (RetainedRegistration, error) {
	if registrationID == (v0candidate.RegistrationID{}) {
		return RetainedRegistration{}, approvalAttemptClassified(approvalattempt.ClassificationSchema, "retained-registration-id-zero")
	}
	loaded, err := store.loadRetainedLookupWorld(ctx)
	if err != nil {
		return RetainedRegistration{}, err
	}
	registration, err := lookupRegistrationInWorld(loaded, registrationID)
	if err != nil {
		return RetainedRegistration{}, retainedLookupFailure(err)
	}
	return cloneRetainedRegistration(registration), nil
}

// ResolveApproval resolves one retained approval by its nominal ID and
// revalidates its registration, nonce, replay, and consumed-attempt links.
func (store *FixedFileStoreV2) ResolveApproval(
	ctx context.Context,
	approvalID approvalattempt.ApprovalID,
) (approvalattempt.ApprovalRecord, error) {
	if approvalID == (approvalattempt.ApprovalID{}) {
		return approvalattempt.ApprovalRecord{}, approvalAttemptClassified(approvalattempt.ClassificationSchema, "retained-approval-id-zero")
	}
	loaded, err := store.loadRetainedLookupWorld(ctx)
	if err != nil {
		return approvalattempt.ApprovalRecord{}, err
	}
	record, err := lookupApprovalInWorld(loaded, approvalID)
	if err != nil {
		return approvalattempt.ApprovalRecord{}, retainedLookupFailure(err)
	}
	return approvalattempt.CloneApprovalRecord(record), nil
}

// ResolveAttempt resolves one retained attempt by its nominal ID and returns
// its exact optional lifecycle without making it recovery authority.
func (store *FixedFileStoreV2) ResolveAttempt(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (RetainedAttempt, error) {
	if attemptID == (approvalattempt.AttemptID{}) {
		return RetainedAttempt{}, approvalAttemptClassified(approvalattempt.ClassificationSchema, "retained-attempt-id-zero")
	}
	loaded, err := store.loadRetainedLookupWorld(ctx)
	if err != nil {
		return RetainedAttempt{}, err
	}
	result, err := lookupAttemptInWorld(loaded, attemptID)
	if err != nil {
		return RetainedAttempt{}, retainedLookupFailure(err)
	}
	return result, nil
}

// ResolveNonce resolves a retained nonce tombstone exclusively through the
// retained-global nonce index, then returns the fully validated approval.
func (store *FixedFileStoreV2) ResolveNonce(
	ctx context.Context,
	nonce approvalattempt.AttemptNonce,
) (RetainedNonce, error) {
	if nonce == (approvalattempt.AttemptNonce{}) {
		return RetainedNonce{}, approvalAttemptClassified(approvalattempt.ClassificationSchema, "retained-nonce-zero")
	}
	loaded, err := store.loadRetainedLookupWorld(ctx)
	if err != nil {
		return RetainedNonce{}, err
	}
	indexes := loaded.Active.View().Indexes.View()
	position := sort.Search(len(indexes.Nonces), func(index int) bool {
		return bytes.Compare(indexes.Nonces[index].AttemptNonce[:], nonce[:]) >= 0
	})
	if position == len(indexes.Nonces) || indexes.Nonces[position].AttemptNonce != nonce {
		return RetainedNonce{}, approvalAttemptClassified(approvalattempt.ClassificationBinding, "retained-nonce-not-found")
	}
	entry := indexes.Nonces[position]
	approval, lookupErr := lookupApprovalInWorld(loaded, entry.ApprovalID)
	if lookupErr != nil || approval.AttemptNonce != nonce || approval.PayloadDigest != entry.PayloadDigest {
		return RetainedNonce{}, retainedLookupFailure(errors.New("retained nonce approval cross-link mismatch"))
	}
	if approvalLocation, ok := approvalIndexLocation(indexes, entry.ApprovalID); !ok || approvalLocation.View() != entry.Location.View() {
		return RetainedNonce{}, retainedLookupFailure(errors.New("retained nonce location cross-link mismatch"))
	}
	return RetainedNonce{AttemptNonce: nonce, PayloadDigest: entry.PayloadDigest, Approval: approvalattempt.CloneApprovalRecord(approval)}, nil
}

// ResolveEffect resolves one effect tombstone from the independent retained
// ledger. Superseded effects return tombstone-bound issuance facts and the
// immutable attempt, but never the current lifecycle as their issuing record.
func (store *FixedFileStoreV2) ResolveEffect(
	ctx context.Context,
	effectID lifecyclestate.EffectID,
) (RetainedEffect, error) {
	if effectID.IsZero() {
		return RetainedEffect{}, approvalAttemptClassified(approvalattempt.ClassificationSchema, "retained-effect-id-zero")
	}
	loaded, err := store.loadRetainedLookupWorld(ctx)
	if err != nil {
		return RetainedEffect{}, err
	}
	indexes := loaded.Active.View().Indexes.View()
	position := sort.Search(len(indexes.Effects), func(index int) bool {
		return bytes.Compare(indexes.Effects[index].EffectID[:], effectID[:]) >= 0
	})
	if position == len(indexes.Effects) || indexes.Effects[position].EffectID != effectID {
		return RetainedEffect{}, approvalAttemptClassified(approvalattempt.ClassificationBinding, "retained-effect-not-found")
	}
	entry := indexes.Effects[position]
	attempt, lookupErr := lookupAttemptInWorld(loaded, entry.AttemptID)
	if lookupErr != nil || !attempt.LifecyclePresent {
		return RetainedEffect{}, retainedLookupFailure(errors.New("retained effect attempt or lifecycle is missing"))
	}
	view := attempt.Lifecycle.View()
	if attemptLocation, ok := attemptIndexLocation(indexes, entry.AttemptID); !ok || attemptLocation.View() != entry.AttemptLocation.View() {
		return RetainedEffect{}, retainedLookupFailure(errors.New("retained effect location cross-link mismatch"))
	}
	result := RetainedEffect{
		EffectID: entry.EffectID, AttemptID: entry.AttemptID, OperationSequence: entry.OperationSequence,
		Operation: entry.Operation, IssuanceSnapshotGeneration: entry.IssuanceSnapshotGeneration,
		VisibleV1Seed: entry.VisibleV1Seed, Attempt: attempt.Attempt,
	}
	if view.EffectID == entry.EffectID {
		if view.OperationSequence != entry.OperationSequence || view.Operation != entry.Operation ||
			entry.IssuanceSnapshotGeneration > view.SnapshotGeneration {
			return RetainedEffect{}, retainedLookupFailure(errors.New("retained current effect lifecycle cross-link mismatch"))
		}
		result.Classification = RetainedEffectCurrent
		result.LifecyclePresent = true
		result.Lifecycle = attempt.Lifecycle
		return result, nil
	}
	if entry.OperationSequence >= view.OperationSequence || entry.IssuanceSnapshotGeneration > view.SnapshotGeneration {
		return RetainedEffect{}, retainedLookupFailure(errors.New("retained historical effect lifecycle order mismatch"))
	}
	result.Classification = RetainedEffectSupersededByCurrent
	return result, nil
}

// ResolveInstance resolves one retained instance digest through its tombstone
// and returns the exact lifecycle record that owns it.
func (store *FixedFileStoreV2) ResolveInstance(
	ctx context.Context,
	digest lifecyclestate.BackendInstanceDigest,
) (RetainedInstance, error) {
	if digest == (lifecyclestate.BackendInstanceDigest{}) {
		return RetainedInstance{}, approvalAttemptClassified(approvalattempt.ClassificationSchema, "retained-instance-digest-zero")
	}
	loaded, err := store.loadRetainedLookupWorld(ctx)
	if err != nil {
		return RetainedInstance{}, err
	}
	indexes := loaded.Active.View().Indexes.View()
	position := sort.Search(len(indexes.Instances), func(index int) bool {
		return bytes.Compare(indexes.Instances[index].InstanceDigest[:], digest[:]) >= 0
	})
	if position == len(indexes.Instances) || indexes.Instances[position].InstanceDigest != digest {
		return RetainedInstance{}, approvalAttemptClassified(approvalattempt.ClassificationBinding, "retained-instance-not-found")
	}
	entry := indexes.Instances[position]
	attempt, lookupErr := lookupAttemptInWorld(loaded, entry.AttemptID)
	if lookupErr != nil || !attempt.LifecyclePresent || !attempt.Lifecycle.View().Instance.Present() ||
		attempt.Lifecycle.View().Instance.Digest() != digest {
		return RetainedInstance{}, retainedLookupFailure(errors.New("retained instance lifecycle cross-link mismatch"))
	}
	lifecycleLocation, locationErr := lifecycleLocationForAttempt(indexes, entry.AttemptID)
	if locationErr != nil || lifecycleLocation.View() != entry.Location.View() {
		return RetainedInstance{}, retainedLookupFailure(errors.New("retained instance location cross-link mismatch"))
	}
	return RetainedInstance{InstanceDigest: digest, AttemptID: entry.AttemptID, Lifecycle: attempt.Lifecycle}, nil
}

// ResolveApprovalReplay preserves the existing exact-payload replay oracle.
// A semantic digest collision with different exact bytes is REPLAY; matching
// bytes with a different authorization identity are BINDING.
func (store *FixedFileStoreV2) ResolveApprovalReplay(
	ctx context.Context,
	payloadDigest approvalattempt.ApprovalPayloadDigest,
	exactPayloadBytes []byte,
	authorization approvalattempt.ApprovalKeyAuthorizationIdentity,
) (approvalattempt.ApprovalReference, approvalattempt.ApprovalState, error) {
	if payloadDigest == (approvalattempt.ApprovalPayloadDigest{}) || len(exactPayloadBytes) == 0 ||
		authorization == (approvalattempt.ApprovalKeyAuthorizationIdentity{}) {
		return approvalattempt.ApprovalReference{}, "", approvalAttemptClassified(approvalattempt.ClassificationSchema, "retained-approval-replay-key-invalid")
	}
	loaded, err := store.loadRetainedLookupWorld(ctx)
	if err != nil {
		return approvalattempt.ApprovalReference{}, "", err
	}
	indexes := loaded.Active.View().Indexes.View()
	start := sort.Search(len(indexes.ApprovalReplay), func(index int) bool {
		return bytes.Compare(indexes.ApprovalReplay[index].PayloadDigest[:], payloadDigest[:]) >= 0
	})
	if start == len(indexes.ApprovalReplay) || indexes.ApprovalReplay[start].PayloadDigest != payloadDigest {
		return approvalattempt.ApprovalReference{}, "", approvalAttemptClassified(approvalattempt.ClassificationBinding, "retained-approval-replay-not-found")
	}
	entry := indexes.ApprovalReplay[start]
	approval, lookupErr := lookupApprovalInWorld(loaded, entry.ApprovalID)
	if lookupErr != nil {
		return approvalattempt.ApprovalReference{}, "", retainedLookupFailure(lookupErr)
	}
	if !bytes.Equal(approval.ExactPayloadBytes, exactPayloadBytes) {
		return approvalattempt.ApprovalReference{}, "", approvalAttemptClassified(approvalattempt.ClassificationReplay, "retained-approval-payload-digest-collision")
	}
	if entry.AuthorizationIdentity != authorization || approval.AuthorizationIdentity != authorization {
		return approvalattempt.ApprovalReference{}, "", approvalAttemptClassified(approvalattempt.ClassificationBinding, "retained-approval-replay-authorization-mismatch")
	}
	exactDigest, digestErr := archivestate.DigestRecord(archivestate.RecordApproval, exactPayloadBytes)
	if digestErr != nil || exactDigest != entry.ExactPayloadDigest || approval.State != entry.State {
		return approvalattempt.ApprovalReference{}, "", retainedLookupFailure(errors.New("retained approval replay tombstone mismatch"))
	}
	reference, referenceErr := approvalattempt.NewApprovalReference(entry.ApprovalID)
	if referenceErr != nil {
		return approvalattempt.ApprovalReference{}, "", retainedLookupFailure(referenceErr)
	}
	return reference, entry.State, nil
}

// ResolveAttemptReplay resolves the exact retained (registration, approval)
// replay identity to its one immutable attempt.
func (store *FixedFileStoreV2) ResolveAttemptReplay(
	ctx context.Context,
	registrationID v0candidate.RegistrationID,
	approvalID approvalattempt.ApprovalID,
) (approvalattempt.AttemptReference, approvalattempt.AttemptState, error) {
	if registrationID == (v0candidate.RegistrationID{}) || approvalID == (approvalattempt.ApprovalID{}) {
		return approvalattempt.AttemptReference{}, "", approvalAttemptClassified(approvalattempt.ClassificationSchema, "retained-attempt-replay-key-invalid")
	}
	loaded, err := store.loadRetainedLookupWorld(ctx)
	if err != nil {
		return approvalattempt.AttemptReference{}, "", err
	}
	indexes := loaded.Active.View().Indexes.View()
	position := sort.Search(len(indexes.AttemptReplay), func(index int) bool {
		entry := indexes.AttemptReplay[index]
		comparison := bytes.Compare(entry.RegistrationID[:], registrationID[:])
		return comparison > 0 || (comparison == 0 && bytes.Compare(entry.ApprovalID[:], approvalID[:]) >= 0)
	})
	if position == len(indexes.AttemptReplay) || indexes.AttemptReplay[position].RegistrationID != registrationID ||
		indexes.AttemptReplay[position].ApprovalID != approvalID {
		return approvalattempt.AttemptReference{}, "", approvalAttemptClassified(approvalattempt.ClassificationBinding, "retained-attempt-replay-not-found")
	}
	entry := indexes.AttemptReplay[position]
	attempt, lookupErr := lookupAttemptInWorld(loaded, entry.AttemptID)
	if lookupErr != nil || attempt.Attempt.RegistrationID != registrationID || attempt.Attempt.ApprovalID != approvalID ||
		attempt.Attempt.State != entry.State {
		return approvalattempt.AttemptReference{}, "", retainedLookupFailure(errors.New("retained attempt replay tombstone mismatch"))
	}
	reference, referenceErr := approvalattempt.NewAttemptReference(entry.AttemptID)
	if referenceErr != nil {
		return approvalattempt.AttemptReference{}, "", retainedLookupFailure(referenceErr)
	}
	return reference, entry.State, nil
}

// QueryRetainedIdentityCollisions passively checks candidate identities only
// against the complete retained-global indexes. It performs no reservation,
// mutation, fallback, or adapter call.
func (store *FixedFileStoreV2) QueryRetainedIdentityCollisions(
	ctx context.Context,
	candidates RetainedIdentityCandidates,
) (RetainedIdentityCollisions, error) {
	if err := validateRetainedIdentityCandidates(candidates); err != nil {
		return RetainedIdentityCollisions{}, err
	}
	loaded, err := store.loadRetainedLookupWorld(ctx)
	if err != nil {
		return RetainedIdentityCollisions{}, err
	}
	indexes := loaded.Active.View().Indexes.View()
	return RetainedIdentityCollisions{
		Registration:   containsRegistration(indexes, candidates.RegistrationID),
		Approval:       containsApproval(indexes, candidates.ApprovalID),
		Attempt:        containsAttempt(indexes, candidates.AttemptID),
		Nonce:          containsNonce(indexes, candidates.AttemptNonce),
		Effect:         containsEffect(indexes, candidates.EffectID),
		Instance:       containsInstance(indexes, candidates.InstanceDigest),
		ApprovalReplay: containsApprovalReplay(indexes, candidates.ApprovalReplayPayloadDigest, candidates.ApprovalReplayAuthorizationIdentity),
		AttemptReplay:  containsAttemptReplay(indexes, candidates.AttemptReplayRegistrationID, candidates.AttemptReplayApprovalID),
	}, nil
}

// CheckRetainedIdentityAvailability returns REPLAY when any queried identity
// is already retained. F4B independently repeats the same closed retained-
// global collision validation while deriving every atomic successor; this
// passive F4A helper alone still reserves nothing.
func (store *FixedFileStoreV2) CheckRetainedIdentityAvailability(
	ctx context.Context,
	candidates RetainedIdentityCandidates,
) error {
	collisions, err := store.QueryRetainedIdentityCollisions(ctx, candidates)
	if err != nil {
		return err
	}
	if collisions.Any() {
		return approvalAttemptClassified(approvalattempt.ClassificationReplay, "retained-identity-unavailable")
	}
	return nil
}

// RecoveryAttemptIDs preserves AttemptID-only recovery and derives work only
// from retained-global attempt entries on the hot arm. Archive entries are
// terminal history and are never returned.
func (store *FixedFileStoreV2) RecoveryAttemptIDs(ctx context.Context) ([]approvalattempt.AttemptID, error) {
	loaded, err := store.loadRetainedLookupWorld(ctx)
	if err != nil {
		return nil, err
	}
	indexes := loaded.Active.View().Indexes.View()
	result := make([]approvalattempt.AttemptID, 0, len(indexes.Attempts))
	for _, entry := range indexes.Attempts {
		if entry.Location.View().Kind == archivestate.RecordLocationArchive {
			continue
		}
		attempt, lookupErr := lookupAttemptInWorld(loaded, entry.AttemptID)
		if lookupErr != nil {
			return nil, retainedLookupFailure(lookupErr)
		}
		if !attempt.LifecyclePresent {
			result = append(result, entry.AttemptID)
			continue
		}
		view := attempt.Lifecycle.View()
		if view.State != lifecyclestate.StateDestroyed || view.CleanupRequired ||
			view.State == lifecyclestate.StateUnresolved || view.State == lifecyclestate.StateQuarantined ||
			view.RecoveryFence != lifecyclestate.RecoveryFenceNone || view.NextRecoveryAt.Present {
			result = append(result, entry.AttemptID)
		}
	}
	return append([]approvalattempt.AttemptID(nil), result...), nil
}

func (store *FixedFileStoreV2) loadRetainedLookupWorld(ctx context.Context) (loadedV2State, error) {
	if err := ctx.Err(); err != nil {
		return loadedV2State{}, err
	}
	if store == nil || store.path == "" {
		return loadedV2State{}, fmt.Errorf("%w: retained lookup store is unavailable", ErrStoreRepairRequired)
	}
	if err := validateExistingStoreFile(store.path); err != nil {
		return loadedV2State{}, fmt.Errorf("%w: retained lookup active snapshot: %v", ErrStoreRepairRequired, err)
	}
	loaded, err := loadV2State(store.path)
	if err != nil {
		return loadedV2State{}, fmt.Errorf("%w: retained lookup full verification: %v", ErrStoreRepairRequired, err)
	}
	if err := ctx.Err(); err != nil {
		return loadedV2State{}, err
	}
	return loaded, nil
}

func retainedLookupFailure(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	if _, ok := approvalattempt.ErrorClassification(err); ok {
		return err
	}
	return fmt.Errorf("%w: retained lookup cross-link: %v", ErrStoreRepairRequired, err)
}

func cloneRetainedRegistration(registration RetainedRegistration) RetainedRegistration {
	registration.PlanRoleBindings = clonePlanBindings(registration.PlanRoleBindings)
	registration.Record = cloneStoredRegistration(registration.Record)
	return registration
}

func lookupRegistrationInWorld(loaded loadedV2State, registrationID v0candidate.RegistrationID) (RetainedRegistration, error) {
	indexes := loaded.Active.View().Indexes.View()
	position := sort.Search(len(indexes.Registrations), func(index int) bool {
		return bytes.Compare(indexes.Registrations[index].RegistrationID[:], registrationID[:]) >= 0
	})
	if position == len(indexes.Registrations) || indexes.Registrations[position].RegistrationID != registrationID {
		return RetainedRegistration{}, approvalAttemptClassified(approvalattempt.ClassificationBinding, "retained-registration-not-found")
	}
	entry := indexes.Registrations[position]
	record, err := registrationAtLocation(loaded, entry.Location)
	if err != nil {
		return RetainedRegistration{}, err
	}
	digest, err := digestV2Record(archivestate.RecordRegistration, record)
	if err != nil || digest != entry.FullRecordDigest || record.Index.RegistrationID != entry.RegistrationID ||
		record.Index.RegistrationSequence != entry.RegistrationSequence || record.Index.ExpiresAt != entry.ExpiresAt ||
		record.Record.RecomputedPlanDigest != entry.PlanDigest {
		return RetainedRegistration{}, errors.New("retained registration index or record mismatch")
	}
	return RetainedRegistration{
		RegistrationID: record.Index.RegistrationID, RegistrationSequence: record.Index.RegistrationSequence,
		ExpiresAt: record.Index.ExpiresAt, PlanRoleBindings: clonePlanBindings(record.PlanBindings),
		Record: cloneStoredRegistration(record.Record),
	}, nil
}

func lookupApprovalInWorld(loaded loadedV2State, approvalID approvalattempt.ApprovalID) (approvalattempt.ApprovalRecord, error) {
	indexes := loaded.Active.View().Indexes.View()
	position := sort.Search(len(indexes.Approvals), func(index int) bool {
		return bytes.Compare(indexes.Approvals[index].ApprovalID[:], approvalID[:]) >= 0
	})
	if position == len(indexes.Approvals) || indexes.Approvals[position].ApprovalID != approvalID {
		return approvalattempt.ApprovalRecord{}, approvalAttemptClassified(approvalattempt.ClassificationBinding, "retained-approval-not-found")
	}
	entry := indexes.Approvals[position]
	record, err := approvalAtLocation(loaded, entry.Location)
	if err != nil {
		return approvalattempt.ApprovalRecord{}, err
	}
	digest, err := digestV2Record(archivestate.RecordApproval, record)
	if err != nil || digest != entry.FullRecordDigest || record.ApprovalID != entry.ApprovalID ||
		record.RegistrationID != entry.RegistrationID || record.PayloadDigest != entry.PayloadDigest ||
		record.AuthorizationIdentity != entry.AuthorizationIdentity || record.AttemptNonce != entry.AttemptNonce ||
		record.State != entry.State || record.ConsumedAttemptID != entry.ConsumedAttemptID || record.ExpiresAt != entry.ExpiresAt {
		return approvalattempt.ApprovalRecord{}, errors.New("retained approval index or record mismatch")
	}
	registration, err := lookupRegistrationInWorld(loaded, record.RegistrationID)
	if err != nil || record.RegistrationSequence != registration.RegistrationSequence ||
		record.PlanDigest != registration.Record.RecomputedPlanDigest {
		return approvalattempt.ApprovalRecord{}, errors.New("retained approval registration cross-link mismatch")
	}
	if !nonceTombstoneMatches(indexes, record, entry.Location) || !approvalReplayTombstoneMatches(indexes, record, entry.Location) {
		return approvalattempt.ApprovalRecord{}, errors.New("retained approval uniqueness or replay tombstone mismatch")
	}
	if record.State == approvalattempt.ApprovalConsumed {
		attemptPosition := searchAttempt(indexes.Attempts, record.ConsumedAttemptID)
		if attemptPosition < 0 || indexes.Attempts[attemptPosition].ApprovalID != record.ApprovalID ||
			indexes.Attempts[attemptPosition].RegistrationID != record.RegistrationID {
			return approvalattempt.ApprovalRecord{}, errors.New("retained consumed approval attempt cross-link mismatch")
		}
	}
	return approvalattempt.CloneApprovalRecord(record), nil
}

func lookupAttemptInWorld(loaded loadedV2State, attemptID approvalattempt.AttemptID) (RetainedAttempt, error) {
	indexes := loaded.Active.View().Indexes.View()
	position := searchAttempt(indexes.Attempts, attemptID)
	if position < 0 {
		return RetainedAttempt{}, approvalAttemptClassified(approvalattempt.ClassificationBinding, "retained-attempt-not-found")
	}
	entry := indexes.Attempts[position]
	attempt, err := attemptAtLocation(loaded, entry.Location)
	if err != nil {
		return RetainedAttempt{}, err
	}
	digest, err := digestV2Record(archivestate.RecordAttempt, attempt)
	if err != nil || digest != entry.FullRecordDigest || attempt.AttemptID != entry.AttemptID ||
		attempt.ApprovalID != entry.ApprovalID || attempt.RegistrationID != entry.RegistrationID ||
		attempt.CreatedAt != entry.CreatedAt || attempt.State != approvalattempt.AttemptCreated {
		return RetainedAttempt{}, errors.New("retained attempt index or record mismatch")
	}
	approval, err := lookupApprovalInWorld(loaded, attempt.ApprovalID)
	if err != nil || approval.RegistrationID != attempt.RegistrationID || approval.ConsumedAttemptID != attempt.AttemptID ||
		approval.State != approvalattempt.ApprovalConsumed {
		return RetainedAttempt{}, errors.New("retained attempt approval cross-link mismatch")
	}
	if !attemptReplayTombstoneMatches(indexes, attempt, entry.Location) {
		return RetainedAttempt{}, errors.New("retained attempt replay tombstone mismatch")
	}
	result := RetainedAttempt{Attempt: attempt}
	lifecycle := entry.Lifecycle.View()
	if lifecycle.Presence == archivestate.AttemptLifecycleAbsent {
		return result, nil
	}
	record, err := lifecycleAtLocation(loaded, lifecycle.Location)
	if err != nil {
		return RetainedAttempt{}, err
	}
	lifecycleDigest, err := digestV2Record(archivestate.RecordLifecycle, lifecycleRecordToDisk(record))
	if err != nil || lifecycleDigest != lifecycle.FullRecordDigest || record.Bindings().View().AttemptID != attempt.AttemptID ||
		record.View().State != lifecycle.State {
		return RetainedAttempt{}, errors.New("retained attempt lifecycle cross-link mismatch")
	}
	result.LifecyclePresent = true
	result.Lifecycle = record
	return result, nil
}

func registrationAtLocation(loaded loadedV2State, location archivestate.RecordLocation) (registrationEntry, error) {
	view := location.View()
	if view.RecordKind != archivestate.RecordRegistration {
		return registrationEntry{}, errors.New("retained registration location has wrong record kind")
	}
	if view.Kind == archivestate.RecordLocationHot {
		records := append([]registrationEntry(nil), loaded.State.Registrations...)
		sort.Slice(records, func(left, right int) bool {
			return bytes.Compare(records[left].Index.RegistrationID[:], records[right].Index.RegistrationID[:]) < 0
		})
		position, err := checkedOrdinal(view.HotRecordOrdinal, len(records))
		if err != nil {
			return registrationEntry{}, err
		}
		return records[position], nil
	}
	cohort, err := archiveCohortAt(loaded, view.Archive)
	if err != nil || view.Archive.RecordOrdinal != 1 {
		return registrationEntry{}, errors.New("retained registration archive location is stale")
	}
	return cohort.Registration, nil
}

func approvalAtLocation(loaded loadedV2State, location archivestate.RecordLocation) (approvalattempt.ApprovalRecord, error) {
	view := location.View()
	if view.RecordKind != archivestate.RecordApproval {
		return approvalattempt.ApprovalRecord{}, errors.New("retained approval location has wrong record kind")
	}
	if view.Kind == archivestate.RecordLocationHot {
		records := make([]approvalattempt.ApprovalRecord, len(loaded.State.Approvals))
		for index, record := range loaded.State.Approvals {
			records[index] = approvalattempt.CloneApprovalRecord(record)
		}
		sort.Slice(records, func(left, right int) bool {
			return bytes.Compare(records[left].ApprovalID[:], records[right].ApprovalID[:]) < 0
		})
		position, err := checkedOrdinal(view.HotRecordOrdinal, len(records))
		if err != nil {
			return approvalattempt.ApprovalRecord{}, err
		}
		return approvalattempt.CloneApprovalRecord(records[position]), nil
	}
	cohort, err := archiveCohortAt(loaded, view.Archive)
	if err != nil {
		return approvalattempt.ApprovalRecord{}, err
	}
	position, err := checkedArchiveOrdinal(view.Archive.RecordOrdinal, len(cohort.Approvals))
	if err != nil {
		return approvalattempt.ApprovalRecord{}, err
	}
	return approvalattempt.CloneApprovalRecord(cohort.Approvals[position]), nil
}

func attemptAtLocation(loaded loadedV2State, location archivestate.RecordLocation) (approvalattempt.ExecutionAttempt, error) {
	view := location.View()
	if view.RecordKind != archivestate.RecordAttempt {
		return approvalattempt.ExecutionAttempt{}, errors.New("retained attempt location has wrong record kind")
	}
	if view.Kind == archivestate.RecordLocationHot {
		records := append([]approvalattempt.ExecutionAttempt(nil), loaded.State.Attempts...)
		sort.Slice(records, func(left, right int) bool {
			return bytes.Compare(records[left].AttemptID[:], records[right].AttemptID[:]) < 0
		})
		position, err := checkedOrdinal(view.HotRecordOrdinal, len(records))
		if err != nil {
			return approvalattempt.ExecutionAttempt{}, err
		}
		return records[position], nil
	}
	cohort, err := archiveCohortAt(loaded, view.Archive)
	if err != nil {
		return approvalattempt.ExecutionAttempt{}, err
	}
	position, err := checkedArchiveOrdinal(view.Archive.RecordOrdinal, len(cohort.Attempts))
	if err != nil {
		return approvalattempt.ExecutionAttempt{}, err
	}
	return cohort.Attempts[position], nil
}

func lifecycleAtLocation(loaded loadedV2State, location archivestate.RecordLocation) (lifecyclestate.Record, error) {
	view := location.View()
	if view.RecordKind != archivestate.RecordLifecycle {
		return lifecyclestate.Record{}, errors.New("retained lifecycle location has wrong record kind")
	}
	if view.Kind == archivestate.RecordLocationHot {
		position, err := checkedOrdinal(view.HotRecordOrdinal, len(loaded.Lifecycles))
		if err != nil {
			return lifecyclestate.Record{}, err
		}
		return loaded.Lifecycles[position], nil
	}
	cohort, err := archiveCohortAt(loaded, view.Archive)
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	position, err := checkedArchiveOrdinal(view.Archive.RecordOrdinal, len(cohort.Lifecycles))
	if err != nil {
		return lifecyclestate.Record{}, err
	}
	return restoreLifecycleRecord(cohort.Lifecycles[position], loaded.State.TimeHighWaterUnixSeconds)
}

func archiveCohortAt(loaded loadedV2State, location archivestate.ArchiveLocation) (archiveSegmentCohortDisk, error) {
	var segment *loadedArchiveSegmentV0
	for index := range loaded.Segments {
		if loaded.Segments[index].Disk.Ordinal == location.SegmentOrdinal {
			segment = &loaded.Segments[index]
			break
		}
	}
	if segment == nil {
		return archiveSegmentCohortDisk{}, errors.New("retained archive segment ordinal is stale")
	}
	position, err := checkedArchiveOrdinal(location.CohortOrdinal, len(segment.Disk.Cohorts))
	if err != nil {
		return archiveSegmentCohortDisk{}, err
	}
	return segment.Disk.Cohorts[position], nil
}

func checkedOrdinal(ordinal v0candidate.PositiveUInt53, length int) (int, error) {
	value := uint64(ordinal)
	if length <= 0 || value == 0 || value > v0candidate.MaxSafeInteger || strconv.IntSize != 64 {
		return 0, errors.New("retained hot record ordinal is stale")
	}
	position := int(value - 1) // #nosec G115 -- UInt53 is bounded above and 64-bit int is required above.
	if position >= length {
		return 0, errors.New("retained hot record ordinal is stale")
	}
	return position, nil
}

func checkedArchiveOrdinal(ordinal archivestate.ArchiveOrdinal, length int) (int, error) {
	value := uint64(ordinal)
	if length <= 0 || value == 0 || value > v0candidate.MaxSafeInteger || strconv.IntSize != 64 {
		return 0, errors.New("retained archive record ordinal is stale")
	}
	position := int(value - 1) // #nosec G115 -- UInt53 is bounded above and 64-bit int is required above.
	if position >= length {
		return 0, errors.New("retained archive record ordinal is stale")
	}
	return position, nil
}

func searchAttempt(entries []archivestate.AttemptIndexEntry, attemptID approvalattempt.AttemptID) int {
	position := sort.Search(len(entries), func(index int) bool { return bytes.Compare(entries[index].AttemptID[:], attemptID[:]) >= 0 })
	if position == len(entries) || entries[position].AttemptID != attemptID {
		return -1
	}
	return position
}

func approvalIndexLocation(indexes archivestate.ArchiveIndexesView, approvalID approvalattempt.ApprovalID) (archivestate.RecordLocation, bool) {
	position := sort.Search(len(indexes.Approvals), func(index int) bool { return bytes.Compare(indexes.Approvals[index].ApprovalID[:], approvalID[:]) >= 0 })
	if position == len(indexes.Approvals) || indexes.Approvals[position].ApprovalID != approvalID {
		return archivestate.RecordLocation{}, false
	}
	return indexes.Approvals[position].Location, true
}

func attemptIndexLocation(indexes archivestate.ArchiveIndexesView, attemptID approvalattempt.AttemptID) (archivestate.RecordLocation, bool) {
	position := searchAttempt(indexes.Attempts, attemptID)
	if position < 0 {
		return archivestate.RecordLocation{}, false
	}
	return indexes.Attempts[position].Location, true
}

func lifecycleLocationForAttempt(indexes archivestate.ArchiveIndexesView, attemptID approvalattempt.AttemptID) (archivestate.RecordLocation, error) {
	position := searchAttempt(indexes.Attempts, attemptID)
	if position < 0 || indexes.Attempts[position].Lifecycle.View().Presence != archivestate.AttemptLifecyclePresent {
		return archivestate.RecordLocation{}, errors.New("retained lifecycle location is missing")
	}
	return indexes.Attempts[position].Lifecycle.View().Location, nil
}

func nonceTombstoneMatches(indexes archivestate.ArchiveIndexesView, record approvalattempt.ApprovalRecord, location archivestate.RecordLocation) bool {
	position := sort.Search(len(indexes.Nonces), func(index int) bool {
		return bytes.Compare(indexes.Nonces[index].AttemptNonce[:], record.AttemptNonce[:]) >= 0
	})
	return position < len(indexes.Nonces) && indexes.Nonces[position].AttemptNonce == record.AttemptNonce &&
		indexes.Nonces[position].PayloadDigest == record.PayloadDigest && indexes.Nonces[position].ApprovalID == record.ApprovalID &&
		indexes.Nonces[position].Location.View() == location.View()
}

func approvalReplayTombstoneMatches(indexes archivestate.ArchiveIndexesView, record approvalattempt.ApprovalRecord, location archivestate.RecordLocation) bool {
	exactDigest, err := archivestate.DigestRecord(archivestate.RecordApproval, record.ExactPayloadBytes)
	if err != nil {
		return false
	}
	for _, entry := range indexes.ApprovalReplay {
		if entry.PayloadDigest == record.PayloadDigest && entry.AuthorizationIdentity == record.AuthorizationIdentity {
			return entry.ExactPayloadDigest == exactDigest && entry.ApprovalID == record.ApprovalID &&
				entry.State == record.State && entry.Location.View() == location.View()
		}
	}
	return false
}

func attemptReplayTombstoneMatches(indexes archivestate.ArchiveIndexesView, attempt approvalattempt.ExecutionAttempt, location archivestate.RecordLocation) bool {
	for _, entry := range indexes.AttemptReplay {
		if entry.RegistrationID == attempt.RegistrationID && entry.ApprovalID == attempt.ApprovalID {
			return entry.AttemptID == attempt.AttemptID && entry.State == attempt.State && entry.Location.View() == location.View()
		}
	}
	return false
}

func validateRetainedIdentityCandidates(candidates RetainedIdentityCandidates) error {
	approvalReplayDigestSet := candidates.ApprovalReplayPayloadDigest != (approvalattempt.ApprovalPayloadDigest{})
	approvalReplayAuthorizationSet := candidates.ApprovalReplayAuthorizationIdentity != (approvalattempt.ApprovalKeyAuthorizationIdentity{})
	attemptReplayRegistrationSet := candidates.AttemptReplayRegistrationID != (v0candidate.RegistrationID{})
	attemptReplayApprovalSet := candidates.AttemptReplayApprovalID != (approvalattempt.ApprovalID{})
	if approvalReplayDigestSet != approvalReplayAuthorizationSet || attemptReplayRegistrationSet != attemptReplayApprovalSet {
		return approvalAttemptClassified(approvalattempt.ClassificationSchema, "retained-collision-replay-key-incomplete")
	}
	if candidates == (RetainedIdentityCandidates{}) {
		return approvalAttemptClassified(approvalattempt.ClassificationSchema, "retained-collision-query-empty")
	}
	return nil
}

func containsRegistration(indexes archivestate.ArchiveIndexesView, id v0candidate.RegistrationID) bool {
	if id == (v0candidate.RegistrationID{}) {
		return false
	}
	position := sort.Search(len(indexes.Registrations), func(index int) bool { return bytes.Compare(indexes.Registrations[index].RegistrationID[:], id[:]) >= 0 })
	return position < len(indexes.Registrations) && indexes.Registrations[position].RegistrationID == id
}

func containsApproval(indexes archivestate.ArchiveIndexesView, id approvalattempt.ApprovalID) bool {
	if id == (approvalattempt.ApprovalID{}) {
		return false
	}
	position := sort.Search(len(indexes.Approvals), func(index int) bool { return bytes.Compare(indexes.Approvals[index].ApprovalID[:], id[:]) >= 0 })
	return position < len(indexes.Approvals) && indexes.Approvals[position].ApprovalID == id
}

func containsAttempt(indexes archivestate.ArchiveIndexesView, id approvalattempt.AttemptID) bool {
	return id != (approvalattempt.AttemptID{}) && searchAttempt(indexes.Attempts, id) >= 0
}

func containsNonce(indexes archivestate.ArchiveIndexesView, nonce approvalattempt.AttemptNonce) bool {
	if nonce == (approvalattempt.AttemptNonce{}) {
		return false
	}
	position := sort.Search(len(indexes.Nonces), func(index int) bool { return bytes.Compare(indexes.Nonces[index].AttemptNonce[:], nonce[:]) >= 0 })
	return position < len(indexes.Nonces) && indexes.Nonces[position].AttemptNonce == nonce
}

func containsEffect(indexes archivestate.ArchiveIndexesView, id lifecyclestate.EffectID) bool {
	if id.IsZero() {
		return false
	}
	position := sort.Search(len(indexes.Effects), func(index int) bool { return bytes.Compare(indexes.Effects[index].EffectID[:], id[:]) >= 0 })
	return position < len(indexes.Effects) && indexes.Effects[position].EffectID == id
}

func containsInstance(indexes archivestate.ArchiveIndexesView, digest lifecyclestate.BackendInstanceDigest) bool {
	if digest == (lifecyclestate.BackendInstanceDigest{}) {
		return false
	}
	position := sort.Search(len(indexes.Instances), func(index int) bool { return bytes.Compare(indexes.Instances[index].InstanceDigest[:], digest[:]) >= 0 })
	return position < len(indexes.Instances) && indexes.Instances[position].InstanceDigest == digest
}

func containsApprovalReplay(indexes archivestate.ArchiveIndexesView, digest approvalattempt.ApprovalPayloadDigest, authorization approvalattempt.ApprovalKeyAuthorizationIdentity) bool {
	if digest == (approvalattempt.ApprovalPayloadDigest{}) {
		return false
	}
	for _, entry := range indexes.ApprovalReplay {
		if entry.PayloadDigest == digest && entry.AuthorizationIdentity == authorization {
			return true
		}
	}
	return false
}

func containsAttemptReplay(indexes archivestate.ArchiveIndexesView, registrationID v0candidate.RegistrationID, approvalID approvalattempt.ApprovalID) bool {
	if registrationID == (v0candidate.RegistrationID{}) {
		return false
	}
	for _, entry := range indexes.AttemptReplay {
		if entry.RegistrationID == registrationID && entry.ApprovalID == approvalID {
			return true
		}
	}
	return false
}
