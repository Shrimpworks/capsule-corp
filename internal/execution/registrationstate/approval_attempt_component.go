package registrationstate

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const ApprovalLifetimeSeconds = uint64(300)

type approvalAttemptStateError struct {
	classification approvalattempt.Classification
	code           string
}

func (err *approvalAttemptStateError) Error() string {
	return fmt.Sprintf("%s: %s", err.classification, err.code)
}

func (err *approvalAttemptStateError) ApprovalAttemptClassification() approvalattempt.Classification {
	return err.classification
}

func approvalAttemptClassified(classification approvalattempt.Classification, code string) error {
	return &approvalAttemptStateError{classification: classification, code: code}
}

type ApprovalIdentifierSource interface {
	NewApprovalID(context.Context) (approvalattempt.ApprovalID, error)
}

type AttemptIdentifierSource interface {
	NewAttemptID(context.Context) (approvalattempt.AttemptID, error)
}

type CryptoApprovalIdentifierSource struct{}

func (CryptoApprovalIdentifierSource) NewApprovalID(ctx context.Context) (approvalattempt.ApprovalID, error) {
	if err := ctx.Err(); err != nil {
		return approvalattempt.ApprovalID{}, err
	}
	var identifier approvalattempt.ApprovalID
	_, err := rand.Read(identifier[:])
	return identifier, err
}

type CryptoAttemptIdentifierSource struct{}

func (CryptoAttemptIdentifierSource) NewAttemptID(ctx context.Context) (approvalattempt.AttemptID, error) {
	if err := ctx.Err(); err != nil {
		return approvalattempt.AttemptID{}, err
	}
	var identifier approvalattempt.AttemptID
	_, err := rand.Read(identifier[:])
	return identifier, err
}

type IntegrityPreflight struct {
	ApprovalID     approvalattempt.ApprovalID
	RegistrationID v0candidate.RegistrationID
	PlanDigest     v0candidate.ExecutionPlanDigest
	InstallationID v0candidate.InstallationID
	EpochSequence  v0candidate.UInt53
	EpochDigest    v0candidate.TrustEpochDigest
	SupervisorID   v0candidate.SupervisorID
}

type RuntimeIntegrityAssessor interface {
	Assess(context.Context, IntegrityPreflight) (RuntimeIntegrityAssessment, error)
}

type RuntimeIntegrityAssessment struct {
	Preflight  IntegrityPreflight
	AssessedAt v0candidate.UInt53
	Permitted  bool
}

type ApprovalAttemptCheckpoint string

const (
	CheckpointBeforeApprovalCommit ApprovalAttemptCheckpoint = "before-approval-commit"
	CheckpointAfterApprovalCommit  ApprovalAttemptCheckpoint = "after-approval-commit"
	CheckpointBeforeAttemptCommit  ApprovalAttemptCheckpoint = "before-attempt-commit"
	CheckpointAfterAttemptCommit   ApprovalAttemptCheckpoint = "after-attempt-commit"
)

type ApprovalAttemptCheckpointHook func(context.Context, ApprovalAttemptCheckpoint) error

type ApprovalAttemptOptions struct {
	Store               StateStore
	Clock               TrustedClock
	Verifier            approvalattempt.CandidateVerifier
	ApprovalIdentifiers ApprovalIdentifierSource
	AttemptIdentifiers  AttemptIdentifierSource
	Integrity           RuntimeIntegrityAssessor
	Checkpoint          ApprovalAttemptCheckpointHook
}

// ApprovalAttemptComponent is the unwired Slice B authority boundary. It
// shares the registration store transaction and never calls a lifecycle or
// backend.
type ApprovalAttemptComponent struct {
	store               StateStore
	clock               TrustedClock
	verifier            approvalattempt.CandidateVerifier
	approvalIdentifiers ApprovalIdentifierSource
	attemptIdentifiers  AttemptIdentifierSource
	integrity           RuntimeIntegrityAssessor
	checkpoint          ApprovalAttemptCheckpointHook
}

func NewApprovalAttempt(options ApprovalAttemptOptions) (*ApprovalAttemptComponent, error) {
	if options.Store == nil || options.Clock == nil || options.Verifier == nil ||
		options.ApprovalIdentifiers == nil || options.AttemptIdentifiers == nil || options.Integrity == nil {
		return nil, errors.New("approval store, clock, verifier, identifiers, and integrity assessor are required")
	}
	if options.Checkpoint == nil {
		options.Checkpoint = func(context.Context, ApprovalAttemptCheckpoint) error { return nil }
	}
	return &ApprovalAttemptComponent{
		store:               options.Store,
		clock:               options.Clock,
		verifier:            options.Verifier,
		approvalIdentifiers: options.ApprovalIdentifiers,
		attemptIdentifiers:  options.AttemptIdentifiers,
		integrity:           options.Integrity,
		checkpoint:          options.Checkpoint,
	}, nil
}

func (component *ApprovalAttemptComponent) SubmitApproval(
	ctx context.Context,
	call AuthenticatedCallContext,
	registrationID v0candidate.RegistrationID,
	exactEnvelopeBytes []byte,
) (approvalattempt.ApprovalSubmission, error) {
	if !call.Authenticated || call.Role != CallerBroker || call.Purpose != SubmitApprovalPurpose {
		return approvalattempt.ApprovalSubmission{}, approvalAttemptClassified(
			approvalattempt.ClassificationAuthentication, "approval-caller-not-authorized")
	}
	if registrationID == (v0candidate.RegistrationID{}) {
		return approvalattempt.ApprovalSubmission{}, approvalAttemptClassified(
			approvalattempt.ClassificationSchema, "registration-id-zero")
	}
	if component.store.recoveryFenced() {
		return approvalattempt.ApprovalSubmission{}, approvalAttemptClassified(
			approvalattempt.ClassificationRecoveryRequired, "approval-store-recovery-fenced")
	}
	verified, err := component.verifier.VerifyCandidate(ctx, exactEnvelopeBytes)
	if err != nil {
		return approvalattempt.ApprovalSubmission{}, err
	}
	payloadBytes := verified.PayloadBytes()
	payloadDigest := approvalattempt.ApprovalPayloadDigest(sha256.Sum256(payloadBytes))
	state, err := component.store.snapshot(ctx)
	if err != nil {
		return approvalattempt.ApprovalSubmission{}, component.localOrRecovery(err, "approval-read-failed")
	}
	if existing, ok, replayErr := findPayloadReplay(state.Approvals, payloadDigest, payloadBytes, registrationID, verified.AuthorizationIdentity()); ok || replayErr != nil {
		if replayErr != nil {
			return approvalattempt.ApprovalSubmission{}, replayErr
		}
		return approvalSubmission(existing)
	}
	if err := component.persistObservedTime(ctx); err != nil {
		return approvalattempt.ApprovalSubmission{}, err
	}
	approvalID, err := component.approvalIdentifiers.NewApprovalID(ctx)
	if err != nil || approvalID == (approvalattempt.ApprovalID{}) {
		return approvalattempt.ApprovalSubmission{}, approvalAttemptClassified(
			approvalattempt.ClassificationLocalFailure, "approval-identifier-failed")
	}
	if err := component.checkpoint(ctx, CheckpointBeforeApprovalCommit); err != nil {
		return approvalattempt.ApprovalSubmission{}, approvalAttemptClassified(
			approvalattempt.ClassificationLocalFailure, "approval-interrupted-before-commit")
	}
	var committed approvalattempt.ApprovalRecord
	err = component.store.commitApproval(ctx, func(fresh *installationState) error {
		if existing, ok, replayErr := findPayloadReplay(fresh.Approvals, payloadDigest, payloadBytes, registrationID, verified.AuthorizationIdentity()); ok || replayErr != nil {
			if replayErr != nil {
				return replayErr
			}
			committed = approvalattempt.CloneApprovalRecord(existing)
			return errNoStoreMutation
		}
		if err := validateVerifiedApproval(fresh, registrationID, verified, fresh.TimeHighWaterUnixSeconds); err != nil {
			return err
		}
		if len(fresh.Approvals) >= MaxRetainedApprovals ||
			countUsableApprovals(fresh.Approvals, fresh.TimeHighWaterUnixSeconds) >= MaxUsableApprovals {
			return approvalAttemptClassified(approvalattempt.ClassificationCapacity, "approval-capacity")
		}
		if findApproval(fresh.Approvals, approvalID) >= 0 {
			return approvalAttemptClassified(approvalattempt.ClassificationLocalFailure, "duplicate-approval-identifier")
		}
		view := verified.View()
		if findNonce(fresh.Approvals, view.AttemptNonce) >= 0 {
			return approvalAttemptClassified(approvalattempt.ClassificationReplay, "approval-nonce-reused")
		}
		registration := fresh.Registrations[findRegistration(fresh.Registrations, registrationID)]
		committed = approvalattempt.ApprovalRecord{
			ApprovalID:            approvalID,
			AttemptNonce:          view.AttemptNonce,
			RegistrationID:        registrationID,
			RegistrationSequence:  registration.Index.RegistrationSequence,
			PlanDigest:            registration.Record.RecomputedPlanDigest,
			InstallationID:        registration.Index.InstallationID,
			EpochSequence:         registration.Index.EpochSequence,
			EpochDigest:           registration.Index.EpochDigest,
			SupervisorID:          registration.Index.SupervisorID,
			Purpose:               view.Purpose,
			Audience:              view.Audience,
			IssuedAt:              view.IssuedAt,
			ExpiresAt:             view.ExpiresAt,
			PayloadDigest:         payloadDigest,
			EnvelopeDigest:        approvalattempt.ApprovalEnvelopeDigest(sha256.Sum256(verified.EnvelopeBytes())),
			ExactEnvelopeBytes:    verified.EnvelopeBytes(),
			ExactPayloadBytes:     payloadBytes,
			ExactProtectedBytes:   verified.ProtectedHeaderBytes(),
			ProtectedKeyID:        verified.ProtectedKeyID(),
			AuthorizationIdentity: verified.AuthorizationIdentity(),
			State:                 approvalattempt.ApprovalUsable,
			StorageFormatVersion:  ApprovalStoreVersion,
		}
		fresh.Approvals = append(fresh.Approvals, approvalattempt.CloneApprovalRecord(committed))
		fresh.ApprovalSetDigest = approvalSetDigest(fresh.Approvals)
		return nil
	})
	if err != nil {
		return approvalattempt.ApprovalSubmission{}, component.storeFailure(err, "approval-commit-failed")
	}
	if err := component.checkpoint(ctx, CheckpointAfterApprovalCommit); err != nil {
		return approvalattempt.ApprovalSubmission{}, approvalAttemptClassified(
			approvalattempt.ClassificationLocalFailure, "approval-response-lost-after-commit")
	}
	return approvalSubmission(committed)
}

func (component *ApprovalAttemptComponent) RequestAttempt(
	ctx context.Context,
	call AuthenticatedCallContext,
	registrationID v0candidate.RegistrationID,
	reference approvalattempt.ApprovalReference,
) (approvalattempt.AttemptCreation, error) {
	if !call.Authenticated || call.Role != CallerDaemon || call.Purpose != RequestAttemptPurpose {
		return approvalattempt.AttemptCreation{}, approvalAttemptClassified(
			approvalattempt.ClassificationAuthentication, "attempt-caller-not-authorized")
	}
	approvalID := reference.ApprovalID()
	if registrationID == (v0candidate.RegistrationID{}) || approvalID == (approvalattempt.ApprovalID{}) {
		return approvalattempt.AttemptCreation{}, approvalAttemptClassified(
			approvalattempt.ClassificationSchema, "attempt-reference-zero")
	}
	if component.store.recoveryFenced() {
		return approvalattempt.AttemptCreation{}, approvalAttemptClassified(
			approvalattempt.ClassificationRecoveryRequired, "attempt-store-recovery-fenced")
	}
	state, err := component.store.snapshot(ctx)
	if err != nil {
		return approvalattempt.AttemptCreation{}, component.localOrRecovery(err, "attempt-read-failed")
	}
	approvalIndex := findApproval(state.Approvals, approvalID)
	if approvalIndex < 0 {
		return approvalattempt.AttemptCreation{}, approvalAttemptClassified(
			approvalattempt.ClassificationBinding, "approval-reference-not-found")
	}
	approval := state.Approvals[approvalIndex]
	if approval.RegistrationID != registrationID {
		return approvalattempt.AttemptCreation{}, approvalAttemptClassified(
			approvalattempt.ClassificationBinding, "approval-registration-mismatch")
	}
	if approval.State == approvalattempt.ApprovalConsumed {
		return attemptCreationForConsumed(state, approval)
	}
	if approval.State != approvalattempt.ApprovalUsable {
		return approvalattempt.AttemptCreation{}, approvalAttemptClassified(
			approvalattempt.ClassificationTrustState, "approval-not-usable")
	}
	if err := component.persistObservedTime(ctx); err != nil {
		return approvalattempt.AttemptCreation{}, err
	}
	state, err = component.store.snapshot(ctx)
	if err != nil {
		return approvalattempt.AttemptCreation{}, component.localOrRecovery(err, "attempt-read-failed")
	}
	approvalIndex = findApproval(state.Approvals, approvalID)
	if approvalIndex < 0 || state.Approvals[approvalIndex].RegistrationID != registrationID {
		return approvalattempt.AttemptCreation{}, approvalAttemptClassified(
			approvalattempt.ClassificationBinding, "approval-reference-not-found")
	}
	approval = state.Approvals[approvalIndex]
	if approval.State == approvalattempt.ApprovalConsumed {
		return attemptCreationForConsumed(state, approval)
	}
	if err := validateAttemptAdmission(&state, approval, registrationID); err != nil {
		return approvalattempt.AttemptCreation{}, err
	}
	preflight := IntegrityPreflight{
		ApprovalID: approval.ApprovalID, RegistrationID: approval.RegistrationID,
		PlanDigest: approval.PlanDigest, InstallationID: approval.InstallationID,
		EpochSequence: approval.EpochSequence, EpochDigest: approval.EpochDigest,
		SupervisorID: approval.SupervisorID,
	}
	assessment, err := component.integrity.Assess(ctx, preflight)
	if err != nil || !assessment.Permitted || assessment.Preflight != preflight ||
		assessment.AssessedAt != state.TimeHighWaterUnixSeconds {
		return approvalattempt.AttemptCreation{}, approvalAttemptClassified(
			approvalattempt.ClassificationTrustState, "runtime-integrity-preflight-failed")
	}
	attemptID, err := component.attemptIdentifiers.NewAttemptID(ctx)
	if err != nil || attemptID == (approvalattempt.AttemptID{}) {
		return approvalattempt.AttemptCreation{}, approvalAttemptClassified(
			approvalattempt.ClassificationLocalFailure, "attempt-identifier-failed")
	}
	if err := component.checkpoint(ctx, CheckpointBeforeAttemptCommit); err != nil {
		return approvalattempt.AttemptCreation{}, approvalAttemptClassified(
			approvalattempt.ClassificationLocalFailure, "attempt-interrupted-before-commit")
	}
	var created approvalattempt.ExecutionAttempt
	err = component.store.commitAttempt(ctx, func(fresh *installationState) error {
		index := findApproval(fresh.Approvals, approvalID)
		if index < 0 || fresh.Approvals[index].RegistrationID != registrationID {
			return approvalAttemptClassified(approvalattempt.ClassificationBinding, "approval-reference-not-found")
		}
		current := fresh.Approvals[index]
		if current.State == approvalattempt.ApprovalConsumed {
			attemptIndex := findAttempt(fresh.Attempts, current.ConsumedAttemptID)
			if attemptIndex < 0 {
				return approvalAttemptClassified(approvalattempt.ClassificationRecoveryRequired, "consumed-approval-attempt-missing")
			}
			created = fresh.Attempts[attemptIndex]
			return errNoStoreMutation
		}
		if err := validateAttemptAdmission(fresh, current, registrationID); err != nil {
			return err
		}
		if assessment.Preflight != (IntegrityPreflight{
			ApprovalID: current.ApprovalID, RegistrationID: current.RegistrationID,
			PlanDigest: current.PlanDigest, InstallationID: current.InstallationID,
			EpochSequence: current.EpochSequence, EpochDigest: current.EpochDigest,
			SupervisorID: current.SupervisorID,
		}) || !assessment.Permitted || assessment.AssessedAt != fresh.TimeHighWaterUnixSeconds {
			return approvalAttemptClassified(approvalattempt.ClassificationTrustState, "runtime-integrity-assessment-stale")
		}
		if len(fresh.Attempts) >= MaxRetainedAttempts || countNonterminalAttempts(fresh.Attempts) >= MaxNonterminalAttempts {
			return approvalAttemptClassified(approvalattempt.ClassificationCapacity, "attempt-capacity")
		}
		if findAttempt(fresh.Attempts, attemptID) >= 0 {
			return approvalAttemptClassified(approvalattempt.ClassificationLocalFailure, "duplicate-attempt-identifier")
		}
		created = approvalattempt.ExecutionAttempt{
			AttemptID: attemptID, ApprovalID: current.ApprovalID, AttemptNonce: current.AttemptNonce,
			RegistrationID: current.RegistrationID, RegistrationSequence: current.RegistrationSequence,
			PlanDigest: current.PlanDigest, InstallationID: current.InstallationID,
			EpochSequence: current.EpochSequence, EpochDigest: current.EpochDigest,
			SupervisorID: current.SupervisorID, Purpose: current.Purpose, Audience: current.Audience,
			ApprovalPayloadDigest: current.PayloadDigest,
			AuthorizationIdentity: current.AuthorizationIdentity,
			CreatedAt:             fresh.TimeHighWaterUnixSeconds, State: approvalattempt.AttemptCreated,
			StorageFormatVersion: AttemptStoreVersion,
		}
		fresh.Attempts = append(fresh.Attempts, created)
		fresh.Approvals[index].State = approvalattempt.ApprovalConsumed
		fresh.Approvals[index].ConsumedAttemptID = attemptID
		fresh.ApprovalSetDigest = approvalSetDigest(fresh.Approvals)
		fresh.AttemptSetDigest = attemptSetDigest(fresh.Attempts)
		return nil
	})
	if err != nil {
		return approvalattempt.AttemptCreation{}, component.storeFailure(err, "attempt-commit-failed")
	}
	if err := component.checkpoint(ctx, CheckpointAfterAttemptCommit); err != nil {
		return approvalattempt.AttemptCreation{}, approvalAttemptClassified(
			approvalattempt.ClassificationLocalFailure, "attempt-response-lost-after-commit")
	}
	return attemptCreation(created)
}

func (component *ApprovalAttemptComponent) Approval(
	ctx context.Context,
	reference approvalattempt.ApprovalReference,
) (approvalattempt.ApprovalRecord, error) {
	if component.store.recoveryFenced() {
		return approvalattempt.ApprovalRecord{}, approvalAttemptClassified(
			approvalattempt.ClassificationRecoveryRequired, "approval-store-recovery-fenced")
	}
	state, err := component.store.snapshot(ctx)
	if err != nil {
		return approvalattempt.ApprovalRecord{}, component.localOrRecovery(err, "approval-read-failed")
	}
	index := findApproval(state.Approvals, reference.ApprovalID())
	if index < 0 {
		return approvalattempt.ApprovalRecord{}, approvalAttemptClassified(
			approvalattempt.ClassificationBinding, "approval-reference-not-found")
	}
	return approvalattempt.CloneApprovalRecord(state.Approvals[index]), nil
}

func (component *ApprovalAttemptComponent) Attempt(
	ctx context.Context,
	reference approvalattempt.AttemptReference,
) (approvalattempt.ExecutionAttempt, error) {
	if component.store.recoveryFenced() {
		return approvalattempt.ExecutionAttempt{}, approvalAttemptClassified(
			approvalattempt.ClassificationRecoveryRequired, "attempt-store-recovery-fenced")
	}
	state, err := component.store.snapshot(ctx)
	if err != nil {
		return approvalattempt.ExecutionAttempt{}, component.localOrRecovery(err, "attempt-read-failed")
	}
	index := findAttempt(state.Attempts, reference.AttemptID())
	if index < 0 {
		return approvalattempt.ExecutionAttempt{}, approvalAttemptClassified(
			approvalattempt.ClassificationBinding, "attempt-reference-not-found")
	}
	return state.Attempts[index], nil
}

// ResolveCreated resolves only a durably committed created attempt. It does
// not consult current approval usability or registration expiry, because a
// committed attempt remains authoritative recovery work after either expires.
// Every returned byte slice and role-binding collection is a defensive copy.
func (component *ApprovalAttemptComponent) ResolveCreated(
	ctx context.Context,
	attemptID approvalattempt.AttemptID,
) (CreatedAttempt, error) {
	if attemptID == (approvalattempt.AttemptID{}) {
		return CreatedAttempt{}, approvalAttemptClassified(
			approvalattempt.ClassificationSchema, "attempt-id-zero")
	}
	if component.store.recoveryFenced() {
		return CreatedAttempt{}, approvalAttemptClassified(
			approvalattempt.ClassificationRecoveryRequired, "attempt-store-recovery-fenced")
	}
	state, err := component.store.snapshot(ctx)
	if err != nil {
		return CreatedAttempt{}, component.localOrRecovery(err, "attempt-read-failed")
	}
	if err := validateState(state); err != nil {
		return CreatedAttempt{}, approvalAttemptClassified(
			approvalattempt.ClassificationRecoveryRequired, "attempt-store-state-invalid")
	}
	attemptIndex := findAttempt(state.Attempts, attemptID)
	if attemptIndex < 0 {
		return CreatedAttempt{}, approvalAttemptClassified(
			approvalattempt.ClassificationBinding, "created-attempt-not-found")
	}
	attempt := state.Attempts[attemptIndex]
	if attempt.State != approvalattempt.AttemptCreated {
		return CreatedAttempt{}, approvalAttemptClassified(
			approvalattempt.ClassificationBinding, "attempt-not-created")
	}
	approvalIndex := findApproval(state.Approvals, attempt.ApprovalID)
	registrationIndex := findRegistration(state.Registrations, attempt.RegistrationID)
	if approvalIndex < 0 || registrationIndex < 0 ||
		!attemptMatchesApproval(attempt, state.Approvals[approvalIndex]) ||
		!approvalMatchesRegistration(state.Approvals[approvalIndex], state.Registrations[registrationIndex]) ||
		attempt.PlanDigest != state.Registrations[registrationIndex].Record.RecomputedPlanDigest ||
		validateStoredRecord(state.Registrations[registrationIndex]) != nil {
		return CreatedAttempt{}, approvalAttemptClassified(
			approvalattempt.ClassificationRecoveryRequired, "created-attempt-cross-link-invalid")
	}
	return cloneCreatedAttempt(CreatedAttempt{
		Attempt:          attempt,
		Approval:         consumedApprovalBinding(state.Approvals[approvalIndex]),
		Registration:     state.Registrations[registrationIndex].Record,
		PlanRoleBindings: state.Registrations[registrationIndex].PlanBindings,
	}), nil
}

func consumedApprovalBinding(record approvalattempt.ApprovalRecord) ConsumedApprovalBinding {
	return ConsumedApprovalBinding{
		ApprovalID:            record.ApprovalID,
		ConsumedAttemptID:     record.ConsumedAttemptID,
		AttemptNonce:          record.AttemptNonce,
		RegistrationID:        record.RegistrationID,
		RegistrationSequence:  record.RegistrationSequence,
		PlanDigest:            record.PlanDigest,
		InstallationID:        record.InstallationID,
		EpochSequence:         record.EpochSequence,
		EpochDigest:           record.EpochDigest,
		SupervisorID:          record.SupervisorID,
		Purpose:               record.Purpose,
		Audience:              record.Audience,
		PayloadDigest:         record.PayloadDigest,
		AuthorizationIdentity: record.AuthorizationIdentity,
		State:                 record.State,
		StorageFormatVersion:  record.StorageFormatVersion,
	}
}

func (component *ApprovalAttemptComponent) CreatedAttempts(ctx context.Context) ([]approvalattempt.AttemptReference, error) {
	if component.store.recoveryFenced() {
		return nil, approvalAttemptClassified(
			approvalattempt.ClassificationRecoveryRequired, "attempt-store-recovery-fenced")
	}
	state, err := component.store.snapshot(ctx)
	if err != nil {
		return nil, component.localOrRecovery(err, "attempt-read-failed")
	}
	result := make([]approvalattempt.AttemptReference, 0, len(state.Attempts))
	for _, attempt := range state.Attempts {
		if attempt.State != approvalattempt.AttemptCreated {
			continue
		}
		reference, referenceErr := approvalattempt.NewAttemptReference(attempt.AttemptID)
		if referenceErr != nil {
			return nil, approvalAttemptClassified(approvalattempt.ClassificationRecoveryRequired, "attempt-reference-invalid")
		}
		result = append(result, reference)
	}
	return result, nil
}

func (component *ApprovalAttemptComponent) persistObservedTime(ctx context.Context) error {
	observation, err := component.clock.ObserveUnixSeconds(ctx)
	if err != nil || observation > v0candidate.MaxSafeInteger {
		return approvalAttemptClassified(approvalattempt.ClassificationLocalFailure, "trusted-clock-observation-failed")
	}
	if err := component.store.persistTimeHighWater(ctx, v0candidate.UInt53(observation)); err != nil {
		return component.localOrRecovery(err, "time-high-water-write-failed")
	}
	return nil
}

func (component *ApprovalAttemptComponent) storeFailure(err error, code string) error {
	if _, ok := approvalattempt.ErrorClassification(err); ok {
		return err
	}
	return component.localOrRecovery(err, code)
}

func (component *ApprovalAttemptComponent) localOrRecovery(err error, code string) error {
	if errors.Is(err, ErrCommitOutcomeIndeterminate) || component.store.recoveryFenced() {
		return approvalAttemptClassified(approvalattempt.ClassificationRecoveryRequired, code)
	}
	return approvalAttemptClassified(approvalattempt.ClassificationLocalFailure, code)
}

func validateVerifiedApproval(
	state *installationState,
	requestedRegistration v0candidate.RegistrationID,
	verified *approvalattempt.VerifiedApproval,
	effectiveNow v0candidate.UInt53,
) error {
	if state.TrustPhase != TrustStable || state.Quarantined {
		return approvalAttemptClassified(approvalattempt.ClassificationTrustState, "approval-trust-state-denied")
	}
	registrationIndex := findRegistration(state.Registrations, requestedRegistration)
	if registrationIndex < 0 {
		return approvalAttemptClassified(approvalattempt.ClassificationBinding, "approval-registration-not-found")
	}
	registration := state.Registrations[registrationIndex]
	if err := validateStoredRecord(registration); err != nil {
		return approvalAttemptClassified(approvalattempt.ClassificationRecoveryRequired, "registration-record-invalid")
	}
	view := verified.View()
	if view.RegistrationID != requestedRegistration ||
		view.PlanDigest != registration.Record.RecomputedPlanDigest ||
		view.InstallationID != registration.Index.InstallationID ||
		view.EpochDigest != registration.Index.EpochDigest ||
		verified.ResolvedEpochSequence() != registration.Index.EpochSequence ||
		view.SupervisorID != registration.Index.SupervisorID ||
		registration.Index.InstallationID != state.InstallationID ||
		registration.Index.EpochSequence != state.EpochSequence ||
		registration.Index.EpochDigest != state.EpochDigest ||
		registration.Index.SupervisorID != state.SupervisorID {
		return approvalAttemptClassified(approvalattempt.ClassificationBinding, "approval-binding-mismatch")
	}
	if view.IssuedAt > effectiveNow || effectiveNow >= view.ExpiresAt ||
		view.ExpiresAt-view.IssuedAt > v0candidate.UInt53(ApprovalLifetimeSeconds) ||
		view.ExpiresAt > registration.Index.ExpiresAt || effectiveNow >= registration.Index.ExpiresAt {
		return approvalAttemptClassified(approvalattempt.ClassificationStale, "approval-time-invalid")
	}
	if findNonce(state.Approvals, view.AttemptNonce) >= 0 {
		return approvalAttemptClassified(approvalattempt.ClassificationReplay, "approval-nonce-reused")
	}
	return nil
}

func validateAttemptAdmission(
	state *installationState,
	approval approvalattempt.ApprovalRecord,
	registrationID v0candidate.RegistrationID,
) error {
	if state.TrustPhase != TrustStable || state.Quarantined || state.AttemptsDisabled {
		return approvalAttemptClassified(approvalattempt.ClassificationTrustState, "attempt-trust-state-denied")
	}
	if approval.State != approvalattempt.ApprovalUsable || approval.RegistrationID != registrationID {
		return approvalAttemptClassified(approvalattempt.ClassificationBinding, "approval-not-usable-for-registration")
	}
	registrationIndex := findRegistration(state.Registrations, registrationID)
	if registrationIndex < 0 || !approvalMatchesRegistration(approval, state.Registrations[registrationIndex]) {
		return approvalAttemptClassified(approvalattempt.ClassificationBinding, "attempt-registration-binding-mismatch")
	}
	registration := state.Registrations[registrationIndex]
	if err := validateStoredRecord(registration); err != nil {
		return approvalAttemptClassified(approvalattempt.ClassificationRecoveryRequired, "registration-record-invalid")
	}
	if approval.InstallationID != state.InstallationID || approval.EpochSequence != state.EpochSequence ||
		approval.EpochDigest != state.EpochDigest || approval.SupervisorID != state.SupervisorID {
		return approvalAttemptClassified(approvalattempt.ClassificationBinding, "attempt-active-binding-mismatch")
	}
	now := state.TimeHighWaterUnixSeconds
	if approval.IssuedAt > now || now >= approval.ExpiresAt || now >= registration.Index.ExpiresAt {
		return approvalAttemptClassified(approvalattempt.ClassificationStale, "attempt-time-invalid")
	}
	return nil
}

func findPayloadReplay(
	records []approvalattempt.ApprovalRecord,
	digest approvalattempt.ApprovalPayloadDigest,
	payload []byte,
	registrationID v0candidate.RegistrationID,
	authorization approvalattempt.ApprovalKeyAuthorizationIdentity,
) (approvalattempt.ApprovalRecord, bool, error) {
	for _, record := range records {
		if record.PayloadDigest != digest {
			continue
		}
		if !bytes.Equal(record.ExactPayloadBytes, payload) {
			return approvalattempt.ApprovalRecord{}, false, approvalAttemptClassified(
				approvalattempt.ClassificationReplay, "approval-payload-digest-collision")
		}
		if record.RegistrationID != registrationID || record.AuthorizationIdentity != authorization {
			return approvalattempt.ApprovalRecord{}, false, approvalAttemptClassified(
				approvalattempt.ClassificationBinding, "approval-replay-binding-mismatch")
		}
		return approvalattempt.CloneApprovalRecord(record), true, nil
	}
	return approvalattempt.ApprovalRecord{}, false, nil
}

func findApproval(records []approvalattempt.ApprovalRecord, id approvalattempt.ApprovalID) int {
	for index := range records {
		if records[index].ApprovalID == id {
			return index
		}
	}
	return -1
}

func findAttempt(records []approvalattempt.ExecutionAttempt, id approvalattempt.AttemptID) int {
	for index := range records {
		if records[index].AttemptID == id {
			return index
		}
	}
	return -1
}

func findNonce(records []approvalattempt.ApprovalRecord, nonce approvalattempt.AttemptNonce) int {
	for index := range records {
		if records[index].AttemptNonce == nonce {
			return index
		}
	}
	return -1
}

func approvalSubmission(record approvalattempt.ApprovalRecord) (approvalattempt.ApprovalSubmission, error) {
	reference, err := approvalattempt.NewApprovalReference(record.ApprovalID)
	if err != nil {
		return approvalattempt.ApprovalSubmission{}, approvalAttemptClassified(
			approvalattempt.ClassificationRecoveryRequired, "approval-reference-invalid")
	}
	return approvalattempt.ApprovalSubmission{Reference: reference, State: record.State}, nil
}

func attemptCreation(record approvalattempt.ExecutionAttempt) (approvalattempt.AttemptCreation, error) {
	reference, err := approvalattempt.NewAttemptReference(record.AttemptID)
	if err != nil {
		return approvalattempt.AttemptCreation{}, approvalAttemptClassified(
			approvalattempt.ClassificationRecoveryRequired, "attempt-reference-invalid")
	}
	return approvalattempt.AttemptCreation{Reference: reference, State: record.State}, nil
}

func attemptCreationForConsumed(
	state installationState,
	approval approvalattempt.ApprovalRecord,
) (approvalattempt.AttemptCreation, error) {
	index := findAttempt(state.Attempts, approval.ConsumedAttemptID)
	if index < 0 || state.Attempts[index].ApprovalID != approval.ApprovalID {
		return approvalattempt.AttemptCreation{}, approvalAttemptClassified(
			approvalattempt.ClassificationRecoveryRequired, "consumed-approval-attempt-missing")
	}
	return attemptCreation(state.Attempts[index])
}
