package archivestate

import (
	"bytes"
	"sort"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// ArchiveFence is a global validated-snapshot condition that refuses archive
// selection without mutating or normalizing authority state.
type ArchiveFence string

const (
	ArchiveFenceNone            ArchiveFence = "none"
	ArchiveFenceTrustTransition ArchiveFence = "trust-transition"
	ArchiveFenceQuarantined     ArchiveFence = "quarantined"
	ArchiveFenceRepairRequired  ArchiveFence = "repair-required"
	ArchiveFenceRecovery        ArchiveFence = "recovery-fenced"
)

func validArchiveFence(fence ArchiveFence) bool {
	switch fence {
	case ArchiveFenceNone, ArchiveFenceTrustTransition, ArchiveFenceQuarantined,
		ArchiveFenceRepairRequired, ArchiveFenceRecovery:
		return true
	default:
		return false
	}
}

type RegistrationCandidate struct {
	RegistrationID       v0candidate.RegistrationID
	RegistrationSequence v0candidate.PositiveUInt53
	ExpiresAt            v0candidate.UInt53
	ExactRecordBytes     []byte
	RecordDigest         ArchiveRecordDigest
}

type ApprovalCandidate struct {
	ApprovalID           approvalattempt.ApprovalID
	AttemptNonce         approvalattempt.AttemptNonce
	RegistrationID       v0candidate.RegistrationID
	RegistrationSequence v0candidate.PositiveUInt53
	ExpiresAt            v0candidate.UInt53
	State                approvalattempt.ApprovalState
	ConsumedAttemptID    approvalattempt.AttemptID
	ExactRecordBytes     []byte
	RecordDigest         ArchiveRecordDigest
}

type AttemptCandidate struct {
	AttemptID            approvalattempt.AttemptID
	ApprovalID           approvalattempt.ApprovalID
	RegistrationID       v0candidate.RegistrationID
	RegistrationSequence v0candidate.PositiveUInt53
	ExactRecordBytes     []byte
	RecordDigest         ArchiveRecordDigest
}

type LifecycleCandidate struct {
	AttemptID                  approvalattempt.AttemptID
	State                      lifecyclestate.LifecycleState
	CleanupRequired            bool
	TerminalAt                 lifecyclestate.OptionalUnixSeconds
	LastReconciliation         lifecyclestate.ReconciliationStatus
	RecoveryFence              lifecyclestate.RecoveryFenceReason
	NextRecoveryAt             lifecyclestate.OptionalUnixSeconds
	AutomaticRecoveryExhausted bool
	ExactRecordBytes           []byte
	RecordDigest               ArchiveRecordDigest
}

// RegistrationCohortCandidate contains every validated record belonging to
// one indivisible registration cohort.
type RegistrationCohortCandidate struct {
	Registration      RegistrationCandidate
	Approvals         []ApprovalCandidate
	Attempts          []AttemptCandidate
	Lifecycles        []LifecycleCandidate
	RetentionHold     bool
	EncodedByteLength uint64
}

// ValidatedV1OrV2StateView is the complete input to the F1 pure selector. A
// future store adapter may construct it only after full snapshot validation.
type ValidatedV1OrV2StateView struct {
	SnapshotGeneration   v0candidate.PositiveUInt53
	ArchiveGeneration    ArchiveGeneration
	CheckpointHead       ArchiveCheckpointDigest
	OwnerSessionID       lifecyclestate.OwnerSessionID
	DurableTimeHighWater v0candidate.UInt53
	Fence                ArchiveFence
	Cohorts              []RegistrationCohortCandidate
	RecoveryAttemptIDs   []approvalattempt.AttemptID
}

// ValidatedV1OrV2State owns defensive copies and guarantees complete cohort
// cross-links plus exact agreement with the current RecoveryAttemptIDs rule.
type ValidatedV1OrV2State struct{ view ValidatedV1OrV2StateView }

func NewValidatedV1OrV2State(view ValidatedV1OrV2StateView) (ValidatedV1OrV2State, error) {
	if !validPositive(uint64(view.SnapshotGeneration)) ||
		!validPositive(uint64(view.ArchiveGeneration)) || view.OwnerSessionID.IsZero() ||
		uint64(view.DurableTimeHighWater) > v0candidate.MaxSafeInteger ||
		!validArchiveFence(view.Fence) || view.Cohorts == nil || view.RecoveryAttemptIDs == nil {
		return ValidatedV1OrV2State{}, classified(ClassificationSchema, "validated-state-metadata")
	}
	cloned := cloneValidatedStateView(view)
	if err := validateCohorts(cloned.Cohorts); err != nil {
		return ValidatedV1OrV2State{}, err
	}
	if err := validateRecoveryBoundary(cloned.Cohorts, cloned.RecoveryAttemptIDs); err != nil {
		return ValidatedV1OrV2State{}, err
	}
	return ValidatedV1OrV2State{view: cloned}, nil
}

func (state ValidatedV1OrV2State) View() ValidatedV1OrV2StateView {
	return cloneValidatedStateView(state.view)
}

func cloneValidatedStateView(view ValidatedV1OrV2StateView) ValidatedV1OrV2StateView {
	sourceCohorts := view.Cohorts
	view.Cohorts = make([]RegistrationCohortCandidate, len(sourceCohorts))
	for index, cohort := range sourceCohorts {
		view.Cohorts[index] = cloneCandidateCohort(cohort)
	}
	view.RecoveryAttemptIDs = cloneSlice(view.RecoveryAttemptIDs)
	return view
}

func cloneCandidateCohort(cohort RegistrationCohortCandidate) RegistrationCohortCandidate {
	cohort.Registration.ExactRecordBytes = bytes.Clone(cohort.Registration.ExactRecordBytes)
	cohort.Approvals = cloneSlice(cohort.Approvals)
	for index := range cohort.Approvals {
		cohort.Approvals[index].ExactRecordBytes = bytes.Clone(cohort.Approvals[index].ExactRecordBytes)
	}
	cohort.Attempts = cloneSlice(cohort.Attempts)
	for index := range cohort.Attempts {
		cohort.Attempts[index].ExactRecordBytes = bytes.Clone(cohort.Attempts[index].ExactRecordBytes)
	}
	cohort.Lifecycles = cloneSlice(cohort.Lifecycles)
	for index := range cohort.Lifecycles {
		cohort.Lifecycles[index].ExactRecordBytes = bytes.Clone(cohort.Lifecycles[index].ExactRecordBytes)
	}
	return cohort
}

func validateCohorts(cohorts []RegistrationCohortCandidate) error {
	seenRegistrations := make(map[v0candidate.RegistrationID]struct{}, len(cohorts))
	seenApprovals := make(map[approvalattempt.ApprovalID]struct{})
	seenNonces := make(map[approvalattempt.AttemptNonce]struct{})
	seenAttempts := make(map[approvalattempt.AttemptID]struct{})
	var previous RegistrationCandidate
	for index, cohort := range cohorts {
		if index > 0 && compareRegistrationCandidate(previous, cohort.Registration) >= 0 {
			return classified(ClassificationBinding, "selector-cohort-order")
		}
		previous = cohort.Registration
		if err := validateCandidateCohort(cohort, seenRegistrations, seenApprovals, seenNonces, seenAttempts); err != nil {
			return err
		}
	}
	return nil
}

func compareRegistrationCandidate(left, right RegistrationCandidate) int {
	if left.RegistrationSequence < right.RegistrationSequence {
		return -1
	}
	if left.RegistrationSequence > right.RegistrationSequence {
		return 1
	}
	return bytes.Compare(left.RegistrationID[:], right.RegistrationID[:])
}

func validateCandidateCohort(
	cohort RegistrationCohortCandidate,
	seenRegistrations map[v0candidate.RegistrationID]struct{},
	seenApprovals map[approvalattempt.ApprovalID]struct{},
	seenNonces map[approvalattempt.AttemptNonce]struct{},
	seenAttempts map[approvalattempt.AttemptID]struct{},
) error {
	registration := cohort.Registration
	if registration.RegistrationID == (v0candidate.RegistrationID{}) {
		return classified(ClassificationSchema, "selector-registration-id")
	}
	if !validPositive(uint64(registration.RegistrationSequence)) {
		return classified(ClassificationSchema, "selector-registration-sequence")
	}
	if uint64(registration.ExpiresAt) > v0candidate.MaxSafeInteger {
		return classified(ClassificationSchema, "selector-registration-expiry")
	}
	if cohort.Approvals == nil || cohort.Attempts == nil || cohort.Lifecycles == nil {
		return classified(ClassificationRepairNeeded, "selector-registration-collections")
	}
	if cohort.EncodedByteLength == 0 || cohort.EncodedByteLength > uint64(MaxSupervisorArchiveBytes) {
		return classified(ClassificationSchema, "selector-registration-bytes")
	}
	if _, exists := seenRegistrations[registration.RegistrationID]; exists {
		return classified(ClassificationBinding, "selector-registration-collision")
	}
	seenRegistrations[registration.RegistrationID] = struct{}{}
	if err := validateExactRecord(RecordRegistration, registration.ExactRecordBytes, registration.RecordDigest); err != nil {
		return err
	}
	if len(cohort.Approvals) > MaxArchiveApprovalsPerSegment ||
		len(cohort.Attempts) > MaxArchiveAttemptsPerSegment || len(cohort.Lifecycles) > MaxArchiveLifecyclesPerSegment {
		return classified(ClassificationCapacity, "selector-cohort-record-count")
	}
	minimumBytes := uint64(len(registration.ExactRecordBytes))
	approvalsByID := make(map[approvalattempt.ApprovalID]ApprovalCandidate, len(cohort.Approvals))
	for index, approval := range cohort.Approvals {
		if approval.ApprovalID == (approvalattempt.ApprovalID{}) ||
			approval.AttemptNonce == (approvalattempt.AttemptNonce{}) ||
			approval.RegistrationID != registration.RegistrationID ||
			approval.RegistrationSequence != registration.RegistrationSequence ||
			uint64(approval.ExpiresAt) > v0candidate.MaxSafeInteger || approval.ExpiresAt > registration.ExpiresAt ||
			!validApprovalState(approval.State) {
			return classified(ClassificationBinding, "selector-approval")
		}
		if index > 0 && bytes.Compare(cohort.Approvals[index-1].ApprovalID[:], approval.ApprovalID[:]) >= 0 {
			return classified(ClassificationBinding, "selector-approval-order")
		}
		if (approval.State == approvalattempt.ApprovalConsumed) !=
			(approval.ConsumedAttemptID != (approvalattempt.AttemptID{})) {
			return classified(ClassificationBinding, "selector-approval-consumed-link")
		}
		if _, exists := seenApprovals[approval.ApprovalID]; exists {
			return classified(ClassificationBinding, "selector-approval-collision")
		}
		if _, exists := seenNonces[approval.AttemptNonce]; exists {
			return classified(ClassificationBinding, "selector-nonce-collision")
		}
		seenApprovals[approval.ApprovalID] = struct{}{}
		seenNonces[approval.AttemptNonce] = struct{}{}
		if err := validateExactRecord(RecordApproval, approval.ExactRecordBytes, approval.RecordDigest); err != nil {
			return err
		}
		minimumBytes += uint64(len(approval.ExactRecordBytes))
		approvalsByID[approval.ApprovalID] = approval
	}
	attemptsByID := make(map[approvalattempt.AttemptID]AttemptCandidate, len(cohort.Attempts))
	for index, attempt := range cohort.Attempts {
		if attempt.AttemptID == (approvalattempt.AttemptID{}) || attempt.ApprovalID == (approvalattempt.ApprovalID{}) ||
			attempt.RegistrationID != registration.RegistrationID ||
			attempt.RegistrationSequence != registration.RegistrationSequence {
			return classified(ClassificationBinding, "selector-attempt")
		}
		if index > 0 && bytes.Compare(cohort.Attempts[index-1].AttemptID[:], attempt.AttemptID[:]) >= 0 {
			return classified(ClassificationBinding, "selector-attempt-order")
		}
		approval, exists := approvalsByID[attempt.ApprovalID]
		if !exists || approval.State != approvalattempt.ApprovalConsumed || approval.ConsumedAttemptID != attempt.AttemptID {
			return classified(ClassificationRepairNeeded, "selector-attempt-approval-half")
		}
		if _, exists := seenAttempts[attempt.AttemptID]; exists {
			return classified(ClassificationBinding, "selector-attempt-collision")
		}
		seenAttempts[attempt.AttemptID] = struct{}{}
		if err := validateExactRecord(RecordAttempt, attempt.ExactRecordBytes, attempt.RecordDigest); err != nil {
			return err
		}
		minimumBytes += uint64(len(attempt.ExactRecordBytes))
		attemptsByID[attempt.AttemptID] = attempt
	}
	for _, approval := range cohort.Approvals {
		if approval.State == approvalattempt.ApprovalConsumed {
			attempt, exists := attemptsByID[approval.ConsumedAttemptID]
			if !exists || attempt.ApprovalID != approval.ApprovalID {
				return classified(ClassificationRepairNeeded, "selector-approval-attempt-half")
			}
		}
	}
	if len(cohort.Lifecycles) != len(cohort.Attempts) {
		return classified(ClassificationRepairNeeded, "selector-lifecycle-count")
	}
	seenLifecycle := make(map[approvalattempt.AttemptID]struct{}, len(cohort.Lifecycles))
	for index, lifecycle := range cohort.Lifecycles {
		if lifecycle.AttemptID == (approvalattempt.AttemptID{}) || !validLifecycleState(lifecycle.State) ||
			!validRecoveryFence(lifecycle.RecoveryFence) || !validReconciliation(lifecycle.LastReconciliation) {
			return classified(ClassificationBinding, "selector-lifecycle")
		}
		if index > 0 && bytes.Compare(cohort.Lifecycles[index-1].AttemptID[:], lifecycle.AttemptID[:]) >= 0 {
			return classified(ClassificationBinding, "selector-lifecycle-order")
		}
		if _, exists := attemptsByID[lifecycle.AttemptID]; !exists {
			return classified(ClassificationRepairNeeded, "selector-lifecycle-attempt-half")
		}
		if _, exists := seenLifecycle[lifecycle.AttemptID]; exists {
			return classified(ClassificationBinding, "selector-lifecycle-collision")
		}
		seenLifecycle[lifecycle.AttemptID] = struct{}{}
		if lifecycle.State == lifecyclestate.StateDestroyed {
			if lifecycle.CleanupRequired || !lifecycle.TerminalAt.Present ||
				lifecycle.LastReconciliation != lifecyclestate.ReconciliationAuthoritativelyAbsent ||
				lifecycle.AutomaticRecoveryExhausted {
				return classified(ClassificationRepairNeeded, "selector-destroyed-disposition")
			}
		} else if !lifecycle.CleanupRequired || lifecycle.TerminalAt.Present {
			return classified(ClassificationRepairNeeded, "selector-active-disposition")
		}
		if err := validateExactRecord(RecordLifecycle, lifecycle.ExactRecordBytes, lifecycle.RecordDigest); err != nil {
			return err
		}
		minimumBytes += uint64(len(lifecycle.ExactRecordBytes))
	}
	if minimumBytes > cohort.EncodedByteLength {
		return classified(ClassificationBinding, "selector-cohort-byte-underflow")
	}
	return nil
}

func validateExactRecord(kind RecordKind, exact []byte, digest ArchiveRecordDigest) error {
	computed, err := DigestRecord(kind, exact)
	if err != nil {
		return err
	}
	if computed != digest {
		return classified(ClassificationBinding, "selector-record-digest")
	}
	return nil
}

func validLifecycleState(state lifecyclestate.LifecycleState) bool {
	for _, candidate := range lifecyclestate.LifecycleStates() {
		if state == candidate {
			return true
		}
	}
	return false
}

func validRecoveryFence(fence lifecyclestate.RecoveryFenceReason) bool {
	switch fence {
	case lifecyclestate.RecoveryFenceNone, lifecyclestate.RecoveryFenceCommitIndeterminate,
		lifecyclestate.RecoveryFenceCommitFailure, lifecyclestate.RecoveryFenceReconcileUnknown,
		lifecyclestate.RecoveryFenceIdentityMismatch, lifecyclestate.RecoveryFenceAutomaticExhausted,
		lifecyclestate.RecoveryFenceRepairRequired:
		return true
	default:
		return false
	}
}

func validReconciliation(status lifecyclestate.ReconciliationStatus) bool {
	for _, candidate := range lifecyclestate.ReconciliationStatuses() {
		if status == candidate {
			return true
		}
	}
	return false
}

func validateRecoveryBoundary(cohorts []RegistrationCohortCandidate, supplied []approvalattempt.AttemptID) error {
	for index, attemptID := range supplied {
		if attemptID == (approvalattempt.AttemptID{}) ||
			(index > 0 && bytes.Compare(supplied[index-1][:], attemptID[:]) >= 0) {
			return classified(ClassificationBinding, "selector-recovery-order")
		}
	}
	derived := make([]approvalattempt.AttemptID, 0)
	for _, cohort := range cohorts {
		for _, lifecycle := range cohort.Lifecycles {
			if lifecycle.State != lifecyclestate.StateDestroyed || lifecycle.CleanupRequired ||
				lifecycle.RecoveryFence != lifecyclestate.RecoveryFenceNone || lifecycle.NextRecoveryAt.Present {
				derived = append(derived, lifecycle.AttemptID)
			}
		}
	}
	sort.Slice(derived, func(left, right int) bool {
		return bytes.Compare(derived[left][:], derived[right][:]) < 0
	})
	if len(derived) != len(supplied) {
		return classified(ClassificationBinding, "selector-recovery-boundary")
	}
	for index := range derived {
		if derived[index] != supplied[index] {
			return classified(ClassificationBinding, "selector-recovery-boundary")
		}
	}
	return nil
}

// PlannedCohort is the immutable digest-only selection result for one complete
// cohort. No record bytes or caller-selected path cross this boundary.
type PlannedCohort struct {
	RegistrationID       v0candidate.RegistrationID
	RegistrationSequence v0candidate.PositiveUInt53
	CohortDigest         ArchiveCohortDigest
	EncodedByteLength    uint64
	RegistrationDigest   ArchiveRecordDigest
	ApprovalDigests      []ArchiveRecordDigest
	AttemptDigests       []ArchiveRecordDigest
	LifecycleDigests     []ArchiveRecordDigest
}

type ArchivePlanView struct {
	SourceSnapshotGeneration v0candidate.PositiveUInt53
	ArchiveGeneration        ArchiveGeneration
	CheckpointHead           ArchiveCheckpointDigest
	OwnerSessionID           lifecyclestate.OwnerSessionID
	Limits                   ArchiveLimits
	Selected                 []PlannedCohort
	SelectedBytes            uint64
	Digest                   ArchivePlanDigest
}

// ArchivePlan is a sealed immutable result; its fields can only be obtained
// from the pure selector.
type ArchivePlan struct{ view ArchivePlanView }

func (plan ArchivePlan) View() ArchivePlanView { return cloneArchivePlanView(plan.view) }

func cloneArchivePlanView(view ArchivePlanView) ArchivePlanView {
	view.Selected = cloneSlice(view.Selected)
	for index := range view.Selected {
		view.Selected[index].ApprovalDigests = cloneSlice(view.Selected[index].ApprovalDigests)
		view.Selected[index].AttemptDigests = cloneSlice(view.Selected[index].AttemptDigests)
		view.Selected[index].LifecycleDigests = cloneSlice(view.Selected[index].LifecycleDigests)
	}
	return view
}

// SelectClosedRegistrationCohorts deterministically selects a stable prefix
// of complete, expired, recovery-excluded cohorts. It performs no mutation.
func SelectClosedRegistrationCohorts(state ValidatedV1OrV2State, limits ArchiveLimits) (ArchivePlan, error) {
	if err := limits.validate(); err != nil {
		return ArchivePlan{}, err
	}
	view := state.view
	if view.Fence != ArchiveFenceNone {
		return ArchivePlan{}, classified(ClassificationTrustState, "archive-selection-fenced")
	}
	selected := make([]PlannedCohort, 0, min(int(limits.MaxCohorts), len(view.Cohorts)))
	var selectedBytes uint64
	for _, cohort := range view.Cohorts {
		if !cohortEligible(cohort, view.DurableTimeHighWater) {
			continue
		}
		if len(selected) == int(limits.MaxCohorts) ||
			cohort.EncodedByteLength > limits.MaxBytes-selectedBytes {
			break
		}
		planned, err := plannedCohort(cohort)
		if err != nil {
			return ArchivePlan{}, err
		}
		selected = append(selected, planned)
		selectedBytes += cohort.EncodedByteLength
	}
	planView := ArchivePlanView{
		SourceSnapshotGeneration: view.SnapshotGeneration, ArchiveGeneration: view.ArchiveGeneration,
		CheckpointHead: view.CheckpointHead, OwnerSessionID: view.OwnerSessionID, Limits: limits,
		Selected: selected, SelectedBytes: selectedBytes,
	}
	planView.Digest = digestArchivePlan(planView)
	return ArchivePlan{view: cloneArchivePlanView(planView)}, nil
}

func cohortEligible(cohort RegistrationCohortCandidate, highWater v0candidate.UInt53) bool {
	if cohort.RetentionHold || highWater < cohort.Registration.ExpiresAt {
		return false
	}
	for _, approval := range cohort.Approvals {
		if approval.State == approvalattempt.ApprovalUsable && highWater < approval.ExpiresAt {
			return false
		}
	}
	for _, lifecycle := range cohort.Lifecycles {
		if lifecycle.State != lifecyclestate.StateDestroyed || lifecycle.CleanupRequired ||
			!lifecycle.TerminalAt.Present || lifecycle.LastReconciliation != lifecyclestate.ReconciliationAuthoritativelyAbsent ||
			lifecycle.RecoveryFence != lifecyclestate.RecoveryFenceNone || lifecycle.NextRecoveryAt.Present ||
			lifecycle.AutomaticRecoveryExhausted {
			return false
		}
	}
	return true
}

func plannedCohort(cohort RegistrationCohortCandidate) (PlannedCohort, error) {
	result := PlannedCohort{
		RegistrationID: cohort.Registration.RegistrationID, RegistrationSequence: cohort.Registration.RegistrationSequence,
		EncodedByteLength: cohort.EncodedByteLength, RegistrationDigest: cohort.Registration.RecordDigest,
		ApprovalDigests:  make([]ArchiveRecordDigest, len(cohort.Approvals)),
		AttemptDigests:   make([]ArchiveRecordDigest, len(cohort.Attempts)),
		LifecycleDigests: make([]ArchiveRecordDigest, len(cohort.Lifecycles)),
	}
	for index, approval := range cohort.Approvals {
		result.ApprovalDigests[index] = approval.RecordDigest
	}
	for index, attempt := range cohort.Attempts {
		result.AttemptDigests[index] = attempt.RecordDigest
	}
	for index, lifecycle := range cohort.Lifecycles {
		result.LifecycleDigests[index] = lifecycle.RecordDigest
	}
	projection, err := NewCohortProjection(CohortProjectionView{
		RegistrationID: result.RegistrationID, RegistrationSequence: result.RegistrationSequence,
		ExpiresAt:    cohort.Registration.ExpiresAt,
		Registration: ArchiveRecordReference{Kind: RecordRegistration, Digest: result.RegistrationDigest},
		Approvals:    references(RecordApproval, result.ApprovalDigests), Attempts: references(RecordAttempt, result.AttemptDigests),
		Lifecycles: references(RecordLifecycle, result.LifecycleDigests), EncodedByteLength: result.EncodedByteLength,
	})
	if err != nil {
		return PlannedCohort{}, classified(ClassificationBinding, "selector-cohort-projection")
	}
	result.CohortDigest = projection.Digest()
	return result, nil
}

func references(kind RecordKind, digests []ArchiveRecordDigest) []ArchiveRecordReference {
	result := make([]ArchiveRecordReference, len(digests))
	for index, digest := range digests {
		result[index] = ArchiveRecordReference{Kind: kind, Digest: digest}
	}
	return result
}

func digestArchivePlan(view ArchivePlanView) ArchivePlanDigest {
	encoder := newDigestEncoder("capsule.supervisor.archive-plan.v0")
	encoder.uint64(uint64(view.SourceSnapshotGeneration))
	encoder.uint64(uint64(view.ArchiveGeneration))
	encoder.bytes(view.CheckpointHead[:])
	encoder.bytes(view.OwnerSessionID[:])
	encoder.uint64(uint64(view.Limits.MaxCohorts))
	encoder.uint64(view.Limits.MaxBytes)
	encoder.uint64(uint64(len(view.Selected)))
	for _, cohort := range view.Selected {
		encoder.bytes(cohort.RegistrationID[:])
		encoder.uint64(uint64(cohort.RegistrationSequence))
		encoder.bytes(cohort.CohortDigest[:])
		encoder.uint64(cohort.EncodedByteLength)
	}
	encoder.uint64(view.SelectedBytes)
	return digestSum[ArchivePlanDigest](encoder)
}
