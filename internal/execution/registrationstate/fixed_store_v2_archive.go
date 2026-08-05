package registrationstate

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"syscall"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/archivestate"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// ArchiveOwner is the narrow owner capability required by the fixed F3
// conformance transaction. The selected Darwin installation owner satisfies
// this interface; tests use an owned local assertion with the same lifetime
// semantics.
type ArchiveOwner interface {
	OwnerSessionID() lifecyclestate.OwnerSessionID
	CheckHeld(context.Context) error
}

type ArchiveFault string

const (
	FaultArchiveAfterSegmentTempCreate ArchiveFault = "archive-after-segment-temp-create"
	FaultArchiveAfterSegmentTempWrite  ArchiveFault = "archive-after-segment-temp-write" //nolint:gosec // G101: fixed fault label.
	FaultArchiveAfterSegmentTempSync   ArchiveFault = "archive-after-segment-temp-sync"
	FaultArchiveBeforeSegmentPublish   ArchiveFault = "archive-before-segment-publish"
	FaultArchiveAfterSegmentPublish    ArchiveFault = "archive-after-segment-publish"
	FaultArchiveAfterSegmentDirSync    ArchiveFault = "archive-after-segment-directory-sync"
	FaultArchiveAfterActiveTempCreate  ArchiveFault = "archive-after-active-temp-create"
	FaultArchiveAfterActiveTempWrite   ArchiveFault = "archive-after-active-temp-write" //nolint:gosec // G101: fixed fault label.
	FaultArchiveAfterActiveTempSync    ArchiveFault = "archive-after-active-temp-sync"
	FaultArchiveBeforeActivation       ArchiveFault = "archive-before-activation"
	FaultArchiveAfterActivation        ArchiveFault = "archive-after-activation"
	FaultArchiveAfterActiveDirSync     ArchiveFault = "archive-after-active-directory-sync"
	FaultArchiveAfterReopen            ArchiveFault = "archive-after-reopen"
)

type ArchiveFaultInjector interface {
	FailArchiveAt(ArchiveFault) error
}

var (
	ErrArchiveOwnerRequired        = errors.New("fixed supervisor archive requires the installation owner")
	ErrArchiveOwnerSessionMismatch = errors.New("fixed supervisor archive owner session mismatch")
	ErrArchiveStaleTransaction     = errors.New("fixed supervisor archive transaction is stale")
	ErrArchiveNoEligibleCohort     = errors.New("fixed supervisor archive has no eligible cohort")
	ErrArchiveOutcomeIndeterminate = errors.New("fixed supervisor archive outcome is indeterminate")
)

type FixedStoreV2ArchivePlan struct {
	plan  archivestate.ArchivePlan
	owner lifecyclestate.OwnerSessionID
}

func (plan FixedStoreV2ArchivePlan) View() archivestate.ArchivePlanView { return plan.plan.View() }

type PreparedFixedStoreV2Archive struct {
	plan          FixedStoreV2ArchivePlan
	segmentBytes  []byte
	segmentDigest archivestate.ArchiveSegmentDigest
	candidate     diskEnvelopeV2
}

func (prepared PreparedFixedStoreV2Archive) SegmentBytes() []byte {
	return bytes.Clone(prepared.segmentBytes)
}

func (prepared PreparedFixedStoreV2Archive) SegmentDigest() archivestate.ArchiveSegmentDigest {
	return prepared.segmentDigest
}

type VerifiedFixedStoreV2Archive struct {
	prepared PreparedFixedStoreV2Archive
}

func (verified VerifiedFixedStoreV2Archive) SegmentBytes() []byte {
	return verified.prepared.SegmentBytes()
}

func (verified VerifiedFixedStoreV2Archive) SegmentDigest() archivestate.ArchiveSegmentDigest {
	return verified.prepared.segmentDigest
}

type archiveCheckpointDisk struct {
	PreviousCheckpoint       archivestate.ArchiveCheckpointReference `json:"previousCheckpoint"`
	NewSegmentDigest         archivestate.ArchiveSegmentDigest       `json:"newSegmentDigest"`
	DescriptorSetDigest      archivestate.ArchiveDescriptorSetDigest `json:"descriptorSetDigest"`
	ArchiveIndexDigest       archivestate.ArchiveCombinedIndexDigest `json:"archiveIndexDigest"`
	HotSetDigests            archivestate.HotSetDigests              `json:"hotSetDigests"`
	SourceSnapshotGeneration v0candidate.PositiveUInt53              `json:"sourceSnapshotGeneration"`
	ResultSnapshotGeneration v0candidate.PositiveUInt53              `json:"resultSnapshotGeneration"`
	ArchiveGeneration        archivestate.ArchiveGeneration          `json:"archiveGeneration"`
	InstallationID           v0candidate.InstallationID              `json:"installationId"`
	SupervisorID             v0candidate.SupervisorID                `json:"supervisorId"`
	EpochSequence            v0candidate.UInt53                      `json:"epochSequence"`
	EpochDigest              v0candidate.TrustEpochDigest            `json:"epochDigest"`
	DurableTimeHighWater     v0candidate.UInt53                      `json:"durableTimeHighWater"`
	Digest                   archivestate.ArchiveCheckpointDigest    `json:"digest"`
}

type archiveSegmentCohortDisk struct {
	Registration      registrationEntry                  `json:"registration"`
	Approvals         []approvalattempt.ApprovalRecord   `json:"approvals"`
	Attempts          []approvalattempt.ExecutionAttempt `json:"attempts"`
	Lifecycles        []lifecycleRecordDisk              `json:"lifecycles"`
	EncodedByteLength uint64                             `json:"encodedByteLength"`
	CohortDigest      archivestate.ArchiveCohortDigest   `json:"cohortDigest"`
}

type archiveSegmentDiskV0 struct {
	FormatVersion            *uint64                              `json:"formatVersion"`
	Ordinal                  archivestate.ArchiveOrdinal          `json:"ordinal"`
	InstallationID           v0candidate.InstallationID           `json:"installationId"`
	SupervisorID             v0candidate.SupervisorID             `json:"supervisorId"`
	EpochSequence            v0candidate.UInt53                   `json:"epochSequence"`
	EpochDigest              v0candidate.TrustEpochDigest         `json:"epochDigest"`
	SourceSnapshotGeneration v0candidate.PositiveUInt53           `json:"sourceSnapshotGeneration"`
	ArchiveGeneration        archivestate.ArchiveGeneration       `json:"archiveGeneration"`
	DurableTimeHighWater     v0candidate.UInt53                   `json:"durableTimeHighWater"`
	PriorCheckpointDigest    archivestate.ArchiveCheckpointDigest `json:"priorCheckpointDigest"`
	Cohorts                  []archiveSegmentCohortDisk           `json:"cohorts"`
	DerivedIndexes           archiveIndexesDisk                   `json:"derivedIndexes"`
	SetDigests               archivestate.ArchiveSetDigests       `json:"setDigests"`
	Counts                   archivestate.ArchiveCounts           `json:"counts"`
	EncodedByteLength        uint64                               `json:"encodedByteLength"`
	SegmentDigest            string                               `json:"segmentDigest"`
}

type loadedArchiveSegmentV0 struct {
	Disk             archiveSegmentDiskV0
	State            installationState
	Lifecycles       []lifecyclestate.Record
	EffectTombstones []archivestate.EffectIndexEntry
	Segment          archivestate.ArchiveSegment
	Bytes            []byte
}

func cloneLoadedArchiveSegments(segments []loadedArchiveSegmentV0) []loadedArchiveSegmentV0 {
	result := make([]loadedArchiveSegmentV0, len(segments))
	for index, segment := range segments {
		result[index] = segment
		result[index].State = cloneState(segment.State)
		result[index].Lifecycles = cloneLifecycleRecords(segment.Lifecycles, segment.State.TimeHighWaterUnixSeconds)
		result[index].EffectTombstones = cloneEffectTombstones(segment.EffectTombstones)
		result[index].Bytes = bytes.Clone(segment.Bytes)
	}
	return result
}

var fixedStoreV2ArchiveMu sync.Mutex

// PlanArchive performs the F3 pure deterministic selection against a freshly
// reopened complete v2 genesis world. It writes no file and accepts no caller
// record, path, or lifecycle value.
func (store *FixedFileStoreV2) PlanArchive(
	ctx context.Context,
	owner ArchiveOwner,
	limits archivestate.ArchiveLimits,
) (FixedStoreV2ArchivePlan, error) {
	loaded, session, err := store.archiveEntry(ctx, owner)
	if err != nil {
		return FixedStoreV2ArchivePlan{}, err
	}
	plan, err := planArchiveFromLoaded(loaded, session, limits)
	if err != nil {
		return FixedStoreV2ArchivePlan{}, err
	}
	return FixedStoreV2ArchivePlan{plan: plan, owner: session}, nil
}

// PrepareArchive constructs the complete immutable segment and successor v2
// snapshot in memory. The returned value is sealed and owns every byte.
func (store *FixedFileStoreV2) PrepareArchive(
	ctx context.Context,
	owner ArchiveOwner,
	plan FixedStoreV2ArchivePlan,
) (PreparedFixedStoreV2Archive, error) {
	loaded, session, err := store.archiveEntry(ctx, owner)
	if err != nil {
		return PreparedFixedStoreV2Archive{}, err
	}
	if session != plan.owner || session != plan.plan.View().OwnerSessionID {
		return PreparedFixedStoreV2Archive{}, ErrArchiveOwnerSessionMismatch
	}
	want, err := planArchiveFromLoaded(loaded, session, plan.plan.View().Limits)
	if err != nil || !reflect.DeepEqual(want.View(), plan.plan.View()) {
		return PreparedFixedStoreV2Archive{}, ErrArchiveStaleTransaction
	}
	segmentBytes, segment, candidate, err := buildArchiveCandidate(loaded, plan.plan)
	if err != nil {
		return PreparedFixedStoreV2Archive{}, fmt.Errorf("prepare fixed supervisor archive: %w", err)
	}
	return PreparedFixedStoreV2Archive{
		plan: plan, segmentBytes: bytes.Clone(segmentBytes), segmentDigest: segment.Segment.Digest(),
		candidate: candidate,
	}, nil
}

// VerifyPreparedArchive independently closed-decodes the prepared bytes,
// reconstructs the complete successor world, and requires deterministic byte
// equality. It performs no filesystem mutation.
func (store *FixedFileStoreV2) VerifyPreparedArchive(
	ctx context.Context,
	owner ArchiveOwner,
	prepared PreparedFixedStoreV2Archive,
) (VerifiedFixedStoreV2Archive, error) {
	loaded, session, err := store.archiveEntry(ctx, owner)
	if err != nil {
		return VerifiedFixedStoreV2Archive{}, err
	}
	if session != prepared.plan.owner {
		return VerifiedFixedStoreV2Archive{}, ErrArchiveOwnerSessionMismatch
	}
	if err := verifyPreparedAgainstLoaded(loaded, prepared); err != nil {
		return VerifiedFixedStoreV2Archive{}, err
	}
	prepared.segmentBytes = bytes.Clone(prepared.segmentBytes)
	return VerifiedFixedStoreV2Archive{prepared: prepared}, nil
}

// ActivateArchive publishes the already verified immutable segment before it
// atomically commits the version-specific active reference. Every
// post-rename fault is indeterminate and requires a fresh reopen.
func (store *FixedFileStoreV2) ActivateArchive(
	ctx context.Context,
	owner ArchiveOwner,
	verified VerifiedFixedStoreV2Archive,
	faults ArchiveFaultInjector,
) (*FixedFileStoreV2, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	fixedStoreV2ArchiveMu.Lock()
	defer fixedStoreV2ArchiveMu.Unlock()
	loaded, session, err := store.archiveEntry(ctx, owner)
	if err != nil {
		return nil, err
	}
	if session != verified.prepared.plan.owner {
		return nil, ErrArchiveOwnerSessionMismatch
	}
	if err := verifyPreparedAgainstLoaded(loaded, verified.prepared); err != nil {
		return nil, err
	}
	segmentPath, err := publishArchiveSegment(ctx, store.path, owner, verified.prepared, faults)
	if err != nil {
		return nil, err
	}
	activeBytes, err := encodeEnvelopeV2(verified.prepared.candidate)
	if err != nil {
		return nil, err
	}
	if err := commitArchiveActivation(ctx, store.path, activeBytes, owner, session, faults); err != nil {
		return nil, err
	}
	if err := owner.CheckHeld(ctx); err != nil {
		return nil, fmt.Errorf("%w: before archive reopen: %v", ErrArchiveOwnerRequired, err)
	}
	if owner.OwnerSessionID() != session {
		return nil, ErrArchiveOwnerSessionMismatch
	}
	reopened, err := OpenFixedFileStoreV2(store.path)
	if err != nil {
		return nil, fmt.Errorf("reopen activated fixed supervisor archive: %w", err)
	}
	if len(reopened.segments) != 1 || reopened.segments[0].Segment.Digest() != verified.prepared.segmentDigest ||
		archiveSegmentPath(store.path, verified.prepared.segmentDigest) != segmentPath {
		return nil, fmt.Errorf("%w: reopened archive successor mismatch", ErrStoreRepairRequired)
	}
	if err := failArchive(faults, FaultArchiveAfterReopen, true); err != nil {
		return nil, err
	}
	return reopened, nil
}

func (store *FixedFileStoreV2) archiveEntry(
	ctx context.Context,
	owner ArchiveOwner,
) (loadedV2State, lifecyclestate.OwnerSessionID, error) {
	if err := ctx.Err(); err != nil {
		return loadedV2State{}, lifecyclestate.OwnerSessionID{}, err
	}
	if owner == nil {
		return loadedV2State{}, lifecyclestate.OwnerSessionID{}, ErrArchiveOwnerRequired
	}
	if err := owner.CheckHeld(ctx); err != nil {
		return loadedV2State{}, lifecyclestate.OwnerSessionID{}, fmt.Errorf("%w: entry check: %v", ErrArchiveOwnerRequired, err)
	}
	session := owner.OwnerSessionID()
	if session.IsZero() {
		return loadedV2State{}, lifecyclestate.OwnerSessionID{}, ErrArchiveOwnerSessionMismatch
	}
	loaded, err := loadV2State(store.path)
	if err != nil {
		return loadedV2State{}, lifecyclestate.OwnerSessionID{}, fmt.Errorf("%w: %v", ErrStoreRepairRequired, err)
	}
	if loaded.Active.View().CurrentCheckpoint != store.active.View().CurrentCheckpoint ||
		loaded.Active.View().SnapshotGeneration != store.active.View().SnapshotGeneration {
		return loadedV2State{}, lifecyclestate.OwnerSessionID{}, ErrArchiveStaleTransaction
	}
	return loaded, session, nil
}

func planArchiveFromLoaded(
	loaded loadedV2State,
	session lifecyclestate.OwnerSessionID,
	limits archivestate.ArchiveLimits,
) (archivestate.ArchivePlan, error) {
	view := loaded.Active.View()
	if view.CurrentCheckpoint.Kind != archivestate.ArchiveCheckpointMigrationGenesis || len(loaded.Segments) != 0 {
		return archivestate.ArchivePlan{}, ErrArchiveStaleTransaction
	}
	selectorState, err := selectorStateForHotWorld(loaded, session)
	if err != nil {
		return archivestate.ArchivePlan{}, fmt.Errorf("%w: %v", ErrStoreRepairRequired, err)
	}
	plan, err := archivestate.SelectClosedRegistrationCohorts(selectorState, limits)
	if err != nil {
		return archivestate.ArchivePlan{}, err
	}
	if len(plan.View().Selected) == 0 {
		return archivestate.ArchivePlan{}, ErrArchiveNoEligibleCohort
	}
	return plan, nil
}

func selectorStateForHotWorld(
	loaded loadedV2State,
	session lifecyclestate.OwnerSessionID,
) (archivestate.ValidatedV1OrV2State, error) {
	state := loaded.State
	lifecycleByAttempt := make(map[approvalattempt.AttemptID]lifecyclestate.Record, len(loaded.Lifecycles))
	for _, record := range loaded.Lifecycles {
		lifecycleByAttempt[record.Bindings().View().AttemptID] = record
	}
	approvalsByRegistration := make(map[v0candidate.RegistrationID][]approvalattempt.ApprovalRecord)
	for _, approval := range state.Approvals {
		approvalsByRegistration[approval.RegistrationID] = append(approvalsByRegistration[approval.RegistrationID], approval)
	}
	attemptsByRegistration := make(map[v0candidate.RegistrationID][]approvalattempt.ExecutionAttempt)
	for _, attempt := range state.Attempts {
		attemptsByRegistration[attempt.RegistrationID] = append(attemptsByRegistration[attempt.RegistrationID], attempt)
	}
	registrations := append([]registrationEntry(nil), state.Registrations...)
	sort.Slice(registrations, func(left, right int) bool {
		if registrations[left].Index.RegistrationSequence != registrations[right].Index.RegistrationSequence {
			return registrations[left].Index.RegistrationSequence < registrations[right].Index.RegistrationSequence
		}
		return bytes.Compare(registrations[left].Index.RegistrationID[:], registrations[right].Index.RegistrationID[:]) < 0
	})
	cohorts := make([]archivestate.RegistrationCohortCandidate, 0, len(registrations))
	recovery := make([]approvalattempt.AttemptID, 0)
	for _, registration := range registrations {
		registrationExact, err := json.Marshal(registration)
		if err != nil {
			return archivestate.ValidatedV1OrV2State{}, err
		}
		registrationDigest, err := archivestate.DigestRecord(archivestate.RecordRegistration, registrationExact)
		if err != nil {
			return archivestate.ValidatedV1OrV2State{}, err
		}
		approvals := append([]approvalattempt.ApprovalRecord(nil), approvalsByRegistration[registration.Index.RegistrationID]...)
		sort.Slice(approvals, func(left, right int) bool {
			return bytes.Compare(approvals[left].ApprovalID[:], approvals[right].ApprovalID[:]) < 0
		})
		approvalCandidates := make([]archivestate.ApprovalCandidate, len(approvals))
		attempts := append([]approvalattempt.ExecutionAttempt(nil), attemptsByRegistration[registration.Index.RegistrationID]...)
		sort.Slice(attempts, func(left, right int) bool {
			return bytes.Compare(attempts[left].AttemptID[:], attempts[right].AttemptID[:]) < 0
		})
		attemptCandidates := make([]archivestate.AttemptCandidate, len(attempts))
		lifecycleCandidates := make([]archivestate.LifecycleCandidate, len(attempts))
		encodedLength := uint64(len(registrationExact))
		for index, approval := range approvals {
			exact, marshalErr := json.Marshal(approval)
			if marshalErr != nil {
				return archivestate.ValidatedV1OrV2State{}, marshalErr
			}
			digest, digestErr := archivestate.DigestRecord(archivestate.RecordApproval, exact)
			if digestErr != nil {
				return archivestate.ValidatedV1OrV2State{}, digestErr
			}
			approvalCandidates[index] = archivestate.ApprovalCandidate{
				ApprovalID: approval.ApprovalID, AttemptNonce: approval.AttemptNonce,
				RegistrationID: approval.RegistrationID, RegistrationSequence: approval.RegistrationSequence,
				ExpiresAt: approval.ExpiresAt, State: approval.State, ConsumedAttemptID: approval.ConsumedAttemptID,
				ExactRecordBytes: exact, RecordDigest: digest,
			}
			encodedLength += uint64(len(exact))
		}
		for index, attempt := range attempts {
			exact, marshalErr := json.Marshal(attempt)
			if marshalErr != nil {
				return archivestate.ValidatedV1OrV2State{}, marshalErr
			}
			digest, digestErr := archivestate.DigestRecord(archivestate.RecordAttempt, exact)
			if digestErr != nil {
				return archivestate.ValidatedV1OrV2State{}, digestErr
			}
			attemptCandidates[index] = archivestate.AttemptCandidate{
				AttemptID: attempt.AttemptID, ApprovalID: attempt.ApprovalID,
				RegistrationID: attempt.RegistrationID, RegistrationSequence: attempt.RegistrationSequence,
				ExactRecordBytes: exact, RecordDigest: digest,
			}
			encodedLength += uint64(len(exact))
			record, present := lifecycleByAttempt[attempt.AttemptID]
			if !present {
				return archivestate.ValidatedV1OrV2State{}, errors.New("archive selector cannot invent missing lifecycle history")
			}
			disk := lifecycleRecordToDisk(record)
			lifecycleExact, marshalErr := json.Marshal(disk)
			if marshalErr != nil {
				return archivestate.ValidatedV1OrV2State{}, marshalErr
			}
			lifecycleDigest, digestErr := archivestate.DigestRecord(archivestate.RecordLifecycle, lifecycleExact)
			if digestErr != nil {
				return archivestate.ValidatedV1OrV2State{}, digestErr
			}
			recordView := record.View()
			lifecycleCandidates[index] = archivestate.LifecycleCandidate{
				AttemptID: attempt.AttemptID, State: recordView.State, CleanupRequired: recordView.CleanupRequired,
				TerminalAt: recordView.TerminalAt, LastReconciliation: recordView.LastReconciliation,
				RecoveryFence: recordView.RecoveryFence, NextRecoveryAt: recordView.NextRecoveryAt,
				AutomaticRecoveryExhausted: recordView.AutomaticRecoveryExhausted,
				ExactRecordBytes:           lifecycleExact, RecordDigest: lifecycleDigest,
			}
			encodedLength += uint64(len(lifecycleExact))
			if recordView.State != lifecyclestate.StateDestroyed || recordView.CleanupRequired ||
				recordView.RecoveryFence != lifecyclestate.RecoveryFenceNone || recordView.NextRecoveryAt.Present {
				recovery = append(recovery, attempt.AttemptID)
			}
		}
		cohorts = append(cohorts, archivestate.RegistrationCohortCandidate{
			Registration: archivestate.RegistrationCandidate{
				RegistrationID:       registration.Index.RegistrationID,
				RegistrationSequence: registration.Index.RegistrationSequence, ExpiresAt: registration.Index.ExpiresAt,
				ExactRecordBytes: registrationExact, RecordDigest: registrationDigest,
			},
			Approvals: approvalCandidates, Attempts: attemptCandidates, Lifecycles: lifecycleCandidates,
			EncodedByteLength: encodedLength,
		})
	}
	sort.Slice(recovery, func(left, right int) bool { return bytes.Compare(recovery[left][:], recovery[right][:]) < 0 })
	fence := archivestate.ArchiveFenceNone
	switch {
	case state.TrustPhase == TrustTransitionFenced:
		fence = archivestate.ArchiveFenceTrustTransition
	case state.TrustPhase == TrustRepairRequired:
		fence = archivestate.ArchiveFenceRepairRequired
	case state.Quarantined:
		fence = archivestate.ArchiveFenceQuarantined
	case state.AttemptsDisabled:
		fence = archivestate.ArchiveFenceRecovery
	}
	return archivestate.NewValidatedV1OrV2State(archivestate.ValidatedV1OrV2StateView{
		SnapshotGeneration: loaded.Active.View().SnapshotGeneration,
		ArchiveGeneration:  loaded.Active.View().ArchiveGeneration,
		CheckpointHead:     loaded.Active.View().CurrentCheckpoint.Digest,
		OwnerSessionID:     session, DurableTimeHighWater: state.TimeHighWaterUnixSeconds,
		Fence: fence, Cohorts: cohorts, RecoveryAttemptIDs: recovery,
	})
}

func buildArchiveCandidate(
	loaded loadedV2State,
	plan archivestate.ArchivePlan,
) ([]byte, loadedArchiveSegmentV0, diskEnvelopeV2, error) {
	selected := make(map[v0candidate.RegistrationID]archivestate.PlannedCohort, len(plan.View().Selected))
	for _, cohort := range plan.View().Selected {
		selected[cohort.RegistrationID] = cohort
	}
	segmentState := cloneState(loaded.State)
	segmentState.Registrations = nil
	segmentState.Approvals = nil
	segmentState.Attempts = nil
	hotState := cloneState(loaded.State)
	hotState.Registrations = nil
	hotState.Approvals = nil
	hotState.Attempts = nil
	for _, registration := range loaded.State.Registrations {
		if _, exists := selected[registration.Index.RegistrationID]; exists {
			segmentState.Registrations = append(segmentState.Registrations, registration)
		} else {
			hotState.Registrations = append(hotState.Registrations, registration)
		}
	}
	for _, approval := range loaded.State.Approvals {
		if _, exists := selected[approval.RegistrationID]; exists {
			segmentState.Approvals = append(segmentState.Approvals, approvalattempt.CloneApprovalRecord(approval))
		} else {
			hotState.Approvals = append(hotState.Approvals, approvalattempt.CloneApprovalRecord(approval))
		}
	}
	for _, attempt := range loaded.State.Attempts {
		if _, exists := selected[attempt.RegistrationID]; exists {
			segmentState.Attempts = append(segmentState.Attempts, attempt)
		} else {
			hotState.Attempts = append(hotState.Attempts, attempt)
		}
	}
	segmentAttemptIDs := make(map[approvalattempt.AttemptID]struct{}, len(segmentState.Attempts))
	for _, attempt := range segmentState.Attempts {
		segmentAttemptIDs[attempt.AttemptID] = struct{}{}
	}
	segmentLifecycles := make([]lifecyclestate.Record, 0, len(segmentState.Attempts))
	hotLifecycles := make([]lifecyclestate.Record, 0, len(loaded.Lifecycles)-len(segmentState.Attempts))
	for _, record := range loaded.Lifecycles {
		if _, exists := segmentAttemptIDs[record.Bindings().View().AttemptID]; exists {
			segmentLifecycles = append(segmentLifecycles, record)
		} else {
			hotLifecycles = append(hotLifecycles, record)
		}
	}
	segmentEffectTombstones := make([]archivestate.EffectIndexEntry, 0, len(loaded.EffectTombstones))
	hotEffectTombstones := make([]archivestate.EffectIndexEntry, 0, len(loaded.EffectTombstones))
	for _, tombstone := range loaded.EffectTombstones {
		if _, exists := segmentAttemptIDs[tombstone.AttemptID]; exists {
			segmentEffectTombstones = append(segmentEffectTombstones, tombstone)
		} else {
			hotEffectTombstones = append(hotEffectTombstones, tombstone)
		}
	}
	if len(segmentState.Registrations) != len(selected) {
		return nil, loadedArchiveSegmentV0{}, diskEnvelopeV2{}, ErrArchiveStaleTransaction
	}
	if err := recomputeAuthoritySetDigests(&segmentState); err != nil {
		return nil, loadedArchiveSegmentV0{}, diskEnvelopeV2{}, err
	}
	if err := recomputeAuthoritySetDigests(&hotState); err != nil {
		return nil, loadedArchiveSegmentV0{}, diskEnvelopeV2{}, err
	}
	segmentBytes, segment, err := encodeArchiveSegmentV0(archiveSegmentBuild{
		Ordinal: 1, SourceSnapshotGeneration: loaded.Active.View().SnapshotGeneration,
		ArchiveGeneration:     archivestate.ArchiveGeneration(uint64(loaded.Active.View().ArchiveGeneration) + 1),
		PriorCheckpointDigest: loaded.Active.View().CurrentCheckpoint.Digest,
		State:                 segmentState, Lifecycles: segmentLifecycles, EffectTombstones: segmentEffectTombstones,
	})
	if err != nil {
		return nil, loadedArchiveSegmentV0{}, diskEnvelopeV2{}, err
	}
	descriptor, err := archivestate.NewArchiveDescriptor(archivestate.ArchiveDescriptorView{
		Ordinal:                  segment.Segment.View().Ordinal,
		SourceSnapshotGeneration: segment.Segment.View().SourceSnapshotGeneration,
		ArchiveGeneration:        segment.Segment.View().ArchiveGeneration,
		SegmentDigest:            segment.Segment.Digest(), EncodedByteLength: uint64(len(segmentBytes)),
		Counts: segment.Segment.View().Counts,
	})
	if err != nil {
		return nil, loadedArchiveSegmentV0{}, diskEnvelopeV2{}, err
	}
	descriptorDigest, err := archivestate.DigestArchiveDescriptorSet([]archivestate.ArchiveDescriptor{descriptor})
	if err != nil {
		return nil, loadedArchiveSegmentV0{}, diskEnvelopeV2{}, err
	}
	indexes, err := reconstructV2IndexesForWorld(hotState, hotLifecycles, hotEffectTombstones, []loadedArchiveSegmentV0{segment})
	if err != nil {
		return nil, loadedArchiveSegmentV0{}, diskEnvelopeV2{}, err
	}
	hotIndexes, err := reconstructV2Indexes(hotState, hotLifecycles, hotEffectTombstones)
	if err != nil {
		return nil, loadedArchiveSegmentV0{}, diskEnvelopeV2{}, err
	}
	hotCounts := countsForIndexView(hotIndexes.View())
	archivedCounts := segment.Segment.View().Counts
	totalCounts, err := addArchiveCounts(hotCounts, archivedCounts)
	if err != nil {
		return nil, loadedArchiveSegmentV0{}, diskEnvelopeV2{}, err
	}
	hotLifecycleDigest := lifecycleSetDigest(hotLifecycles)
	hotSetDigests := archivestate.HotSetDigests{
		Registrations: hotState.RegistrationSetDigest, Approvals: hotState.ApprovalSetDigest,
		Attempts: hotState.AttemptSetDigest, Lifecycles: [32]byte(hotLifecycleDigest),
	}
	checkpoint, err := archivestate.NewArchiveCheckpoint(archivestate.ArchiveCheckpointView{
		PreviousCheckpoint: loaded.Active.View().CurrentCheckpoint, NewSegmentDigest: segment.Segment.Digest(),
		DescriptorSetDigest: descriptorDigest, ArchiveIndexDigest: indexes.CombinedDigest(), HotSetDigests: hotSetDigests,
		SourceSnapshotGeneration: loaded.Active.View().SnapshotGeneration,
		ResultSnapshotGeneration: v0candidate.PositiveUInt53(uint64(loaded.Active.View().SnapshotGeneration) + 1),
		ArchiveGeneration:        archivestate.ArchiveGeneration(uint64(loaded.Active.View().ArchiveGeneration) + 1),
		InstallationID:           loaded.State.InstallationID, SupervisorID: loaded.State.SupervisorID,
		EpochSequence: loaded.State.EpochSequence, EpochDigest: loaded.State.EpochDigest,
		DurableTimeHighWater: loaded.State.TimeHighWaterUnixSeconds,
	})
	if err != nil {
		return nil, loadedArchiveSegmentV0{}, diskEnvelopeV2{}, err
	}
	active, err := archivestate.NewActiveStateV2(archivestate.ActiveStateV2View{
		StoreFormatVersion:     archivestate.SupervisorStoreFormatV2,
		MigrationSourceVersion: archivestate.MigrationSourceFormatV1,
		SnapshotGeneration:     checkpoint.View().ResultSnapshotGeneration,
		ArchiveGeneration:      checkpoint.View().ArchiveGeneration,
		InstallationID:         hotState.InstallationID, SupervisorID: hotState.SupervisorID,
		EpochSequence: hotState.EpochSequence, EpochDigest: hotState.EpochDigest,
		DurableTimeHighWater: hotState.TimeHighWaterUnixSeconds,
		Descriptors:          []archivestate.ArchiveDescriptor{descriptor}, DescriptorSetDigest: descriptorDigest,
		Indexes: indexes, IndexDigests: indexes.Digests(), CombinedIndexDigest: indexes.CombinedDigest(),
		EffectTombstoneCoverage:   archivestate.EffectTombstoneCoverageVisibleV1AndV2,
		VisibleV1EffectSeedCount:  loaded.Active.View().VisibleV1EffectSeedCount,
		VisibleV1EffectSeedDigest: loaded.Active.View().VisibleV1EffectSeedDigest,
		PreviousCheckpoint:        loaded.Active.View().CurrentCheckpoint, CurrentCheckpoint: checkpoint.Reference(),
		HotCounts: hotCounts, ArchivedCounts: archivedCounts, TotalCounts: totalCounts,
	})
	if err != nil {
		return nil, loadedArchiveSegmentV0{}, diskEnvelopeV2{}, err
	}
	version := archivestate.SupervisorStoreFormatV2
	sourceVersion := archivestate.MigrationSourceFormatV1
	lifecycleDisks := make([]lifecycleRecordDisk, len(hotLifecycles))
	for index, record := range hotLifecycles {
		lifecycleDisks[index] = lifecycleRecordToDisk(record)
	}
	activeView := active.View()
	checkpointView := checkpoint.View()
	candidate := diskEnvelopeV2{
		StoreFormatVersion: &version, MigrationSourceVersion: &sourceVersion,
		SnapshotGeneration: activeView.SnapshotGeneration, ArchiveGeneration: activeView.ArchiveGeneration,
		State: hotState, LifecycleSetDigest: hotLifecycleDigest, Lifecycles: lifecycleDisks,
		Descriptors: []archivestate.ArchiveDescriptorView{descriptor.View()}, DescriptorSetDigest: descriptorDigest,
		Indexes: archiveIndexesToDisk(indexes), IndexDigests: indexes.Digests(), CombinedIndexDigest: indexes.CombinedDigest(),
		EffectTombstoneCoverage:   activeView.EffectTombstoneCoverage,
		VisibleV1EffectSeedCount:  activeView.VisibleV1EffectSeedCount,
		VisibleV1EffectSeedDigest: activeView.VisibleV1EffectSeedDigest,
		PreviousCheckpoint:        activeView.PreviousCheckpoint, CurrentCheckpoint: activeView.CurrentCheckpoint,
		MigrationGenesis: loaded.Envelope.MigrationGenesis,
		ActivationCheckpoint: &archiveCheckpointDisk{
			PreviousCheckpoint: checkpointView.PreviousCheckpoint, NewSegmentDigest: checkpointView.NewSegmentDigest,
			DescriptorSetDigest: checkpointView.DescriptorSetDigest, ArchiveIndexDigest: checkpointView.ArchiveIndexDigest,
			HotSetDigests:            checkpointView.HotSetDigests,
			SourceSnapshotGeneration: checkpointView.SourceSnapshotGeneration,
			ResultSnapshotGeneration: checkpointView.ResultSnapshotGeneration,
			ArchiveGeneration:        checkpointView.ArchiveGeneration,
			InstallationID:           checkpointView.InstallationID, SupervisorID: checkpointView.SupervisorID,
			EpochSequence: checkpointView.EpochSequence, EpochDigest: checkpointView.EpochDigest,
			DurableTimeHighWater: checkpointView.DurableTimeHighWater, Digest: checkpoint.Digest(),
		},
		HotCounts: hotCounts, ArchivedCounts: archivedCounts, TotalCounts: totalCounts,
	}
	return segmentBytes, segment, candidate, nil
}

func recomputeAuthoritySetDigests(state *installationState) error {
	registrationDigest := emptyRegistrationSetDigest()
	for _, entry := range state.Registrations {
		var err error
		registrationDigest, err = appendRegistrationSetDigest(registrationDigest, entry.Record)
		if err != nil {
			return err
		}
	}
	approvalDigest, err := approvalSetDigest(state.Approvals)
	if err != nil {
		return err
	}
	attemptDigest, err := attemptSetDigest(state.Attempts)
	if err != nil {
		return err
	}
	state.RegistrationSetDigest = registrationDigest
	state.ApprovalSetDigest = approvalDigest
	state.AttemptSetDigest = attemptDigest
	return nil
}

func addArchiveCounts(left, right archivestate.ArchiveCounts) (archivestate.ArchiveCounts, error) {
	result := left
	leftValues := []*uint64{&result.Registrations, &result.Approvals, &result.Attempts, &result.Lifecycles,
		&result.Nonces, &result.Effects, &result.Instances, &result.ApprovalReplay, &result.AttemptReplay}
	rightValues := []uint64{right.Registrations, right.Approvals, right.Attempts, right.Lifecycles,
		right.Nonces, right.Effects, right.Instances, right.ApprovalReplay, right.AttemptReplay}
	for index := range leftValues {
		if ^uint64(0)-*leftValues[index] < rightValues[index] {
			return archivestate.ArchiveCounts{}, errors.New("archive count overflow")
		}
		*leftValues[index] += rightValues[index]
	}
	return result, nil
}

type archiveSegmentBuild struct {
	Ordinal                  archivestate.ArchiveOrdinal
	SourceSnapshotGeneration v0candidate.PositiveUInt53
	ArchiveGeneration        archivestate.ArchiveGeneration
	PriorCheckpointDigest    archivestate.ArchiveCheckpointDigest
	State                    installationState
	Lifecycles               []lifecyclestate.Record
	EffectTombstones         []archivestate.EffectIndexEntry
}

func encodeArchiveSegmentV0(build archiveSegmentBuild) ([]byte, loadedArchiveSegmentV0, error) {
	cohortDisks, projections, err := segmentCohorts(build.State, build.Lifecycles)
	if err != nil {
		return nil, loadedArchiveSegmentV0{}, err
	}
	provisional := loadedArchiveSegmentV0{State: cloneState(build.State), Lifecycles: cloneLifecycleRecords(build.Lifecycles, build.State.TimeHighWaterUnixSeconds), EffectTombstones: cloneEffectTombstones(build.EffectTombstones)}
	provisional.Disk = archiveSegmentDiskV0{Ordinal: build.Ordinal, Cohorts: cohortDisks}
	archivedIndexes, err := archiveIndexesForLoadedSegment(provisional)
	if err != nil {
		return nil, loadedArchiveSegmentV0{}, err
	}
	derivedView := archivedIndexes.View()
	derivedView.Scope = archivestate.ArchiveIndexScopeSegmentDerived
	derivedView.Registrations = []archivestate.RegistrationIndexEntry{}
	derivedView.Approvals = []archivestate.ApprovalIndexEntry{}
	derivedView.Attempts = []archivestate.AttemptIndexEntry{}
	derivedIndexes, err := archivestate.NewArchiveIndexes(derivedView)
	if err != nil {
		return nil, loadedArchiveSegmentV0{}, err
	}
	version := archivestate.SupervisorArchiveFormatV0
	encodedLength := uint64(1)
	setDigests, counts, err := archivestate.DeriveArchiveSegmentSets(projections, derivedIndexes)
	if err != nil {
		return nil, loadedArchiveSegmentV0{}, err
	}
	buildAtLength := func(encodedLength uint64) ([]byte, error) {
		segment, segmentErr := archivestate.NewArchiveSegment(archivestate.ArchiveSegmentView{
			FormatVersion: version, Ordinal: build.Ordinal,
			InstallationID: build.State.InstallationID, SupervisorID: build.State.SupervisorID,
			EpochSequence: build.State.EpochSequence, EpochDigest: build.State.EpochDigest,
			SourceSnapshotGeneration: build.SourceSnapshotGeneration, ArchiveGeneration: build.ArchiveGeneration,
			DurableTimeHighWater:  build.State.TimeHighWaterUnixSeconds,
			PriorCheckpointDigest: build.PriorCheckpointDigest, Cohorts: projections,
			DerivedIndexes: derivedIndexes, SetDigests: setDigests,
			Counts: counts, EncodedByteLength: encodedLength,
		})
		if segmentErr != nil {
			return nil, segmentErr
		}
		view := segment.View()
		semanticDigest := segment.Digest()
		disk := archiveSegmentDiskV0{
			FormatVersion: &version, Ordinal: view.Ordinal,
			InstallationID: view.InstallationID, SupervisorID: view.SupervisorID,
			EpochSequence: view.EpochSequence, EpochDigest: view.EpochDigest,
			SourceSnapshotGeneration: view.SourceSnapshotGeneration, ArchiveGeneration: view.ArchiveGeneration,
			DurableTimeHighWater: view.DurableTimeHighWater, PriorCheckpointDigest: view.PriorCheckpointDigest,
			Cohorts: cohortDisks, DerivedIndexes: archiveIndexesToDisk(derivedIndexes),
			SetDigests: view.SetDigests, Counts: view.Counts, EncodedByteLength: encodedLength,
			SegmentDigest: hex.EncodeToString(semanticDigest[:]),
		}
		candidate, marshalErr := json.Marshal(disk)
		if marshalErr != nil {
			return nil, marshalErr
		}
		candidate = append(candidate, '\n')
		if int64(len(candidate)) > archivestate.MaxSupervisorArchiveBytes {
			return nil, errors.New("fixed supervisor archive segment exceeds byte capacity")
		}
		return candidate, nil
	}
	for range 8 {
		candidate, buildErr := buildAtLength(encodedLength)
		if buildErr != nil {
			return nil, loadedArchiveSegmentV0{}, buildErr
		}
		if uint64(len(candidate)) != encodedLength {
			encodedLength = uint64(len(candidate))
			continue
		}
		loaded, decodeErr := decodeArchiveSegmentV0(candidate)
		if decodeErr != nil {
			return nil, loadedArchiveSegmentV0{}, decodeErr
		}
		return candidate, loaded, nil
	}
	return nil, loadedArchiveSegmentV0{}, errors.New("fixed supervisor archive encoded length did not converge")
}

func segmentCohorts(
	state installationState,
	lifecycles []lifecyclestate.Record,
) ([]archiveSegmentCohortDisk, []archivestate.CohortProjection, error) {
	registrations := append([]registrationEntry(nil), state.Registrations...)
	sort.Slice(registrations, func(left, right int) bool {
		if registrations[left].Index.RegistrationSequence != registrations[right].Index.RegistrationSequence {
			return registrations[left].Index.RegistrationSequence < registrations[right].Index.RegistrationSequence
		}
		return bytes.Compare(registrations[left].Index.RegistrationID[:], registrations[right].Index.RegistrationID[:]) < 0
	})
	lifecycleByAttempt := make(map[approvalattempt.AttemptID]lifecyclestate.Record, len(lifecycles))
	for _, record := range lifecycles {
		lifecycleByAttempt[record.Bindings().View().AttemptID] = record
	}
	disks := make([]archiveSegmentCohortDisk, len(registrations))
	projections := make([]archivestate.CohortProjection, len(registrations))
	for cohortIndex, registration := range registrations {
		approvals := make([]approvalattempt.ApprovalRecord, 0)
		for _, approval := range state.Approvals {
			if approval.RegistrationID == registration.Index.RegistrationID {
				approvals = append(approvals, approvalattempt.CloneApprovalRecord(approval))
			}
		}
		sort.Slice(approvals, func(left, right int) bool {
			return bytes.Compare(approvals[left].ApprovalID[:], approvals[right].ApprovalID[:]) < 0
		})
		attempts := make([]approvalattempt.ExecutionAttempt, 0)
		for _, attempt := range state.Attempts {
			if attempt.RegistrationID == registration.Index.RegistrationID {
				attempts = append(attempts, attempt)
			}
		}
		sort.Slice(attempts, func(left, right int) bool {
			return bytes.Compare(attempts[left].AttemptID[:], attempts[right].AttemptID[:]) < 0
		})
		lifecycleDisks := make([]lifecycleRecordDisk, len(attempts))
		registrationExact, _ := json.Marshal(registration)
		registrationDigest, _ := archivestate.DigestRecord(archivestate.RecordRegistration, registrationExact)
		approvalRefs := make([]archivestate.ArchiveRecordReference, len(approvals))
		attemptRefs := make([]archivestate.ArchiveRecordReference, len(attempts))
		lifecycleRefs := make([]archivestate.ArchiveRecordReference, len(attempts))
		encodedLength := uint64(len(registrationExact))
		for index, approval := range approvals {
			exact, _ := json.Marshal(approval)
			digest, _ := archivestate.DigestRecord(archivestate.RecordApproval, exact)
			approvalRefs[index] = archivestate.ArchiveRecordReference{Kind: archivestate.RecordApproval, Digest: digest}
			encodedLength += uint64(len(exact))
		}
		for index, attempt := range attempts {
			exact, _ := json.Marshal(attempt)
			digest, _ := archivestate.DigestRecord(archivestate.RecordAttempt, exact)
			attemptRefs[index] = archivestate.ArchiveRecordReference{Kind: archivestate.RecordAttempt, Digest: digest}
			encodedLength += uint64(len(exact))
			record, exists := lifecycleByAttempt[attempt.AttemptID]
			if !exists {
				return nil, nil, errors.New("archive segment cohort lacks lifecycle")
			}
			lifecycleDisks[index] = lifecycleRecordToDisk(record)
			lifecycleExact, _ := json.Marshal(lifecycleDisks[index])
			lifecycleDigest, _ := archivestate.DigestRecord(archivestate.RecordLifecycle, lifecycleExact)
			lifecycleRefs[index] = archivestate.ArchiveRecordReference{Kind: archivestate.RecordLifecycle, Digest: lifecycleDigest}
			encodedLength += uint64(len(lifecycleExact))
		}
		projection, err := archivestate.NewCohortProjection(archivestate.CohortProjectionView{
			RegistrationID:       registration.Index.RegistrationID,
			RegistrationSequence: registration.Index.RegistrationSequence, ExpiresAt: registration.Index.ExpiresAt,
			Registration: archivestate.ArchiveRecordReference{Kind: archivestate.RecordRegistration, Digest: registrationDigest},
			Approvals:    approvalRefs, Attempts: attemptRefs, Lifecycles: lifecycleRefs, EncodedByteLength: encodedLength,
		})
		if err != nil {
			return nil, nil, err
		}
		disks[cohortIndex] = archiveSegmentCohortDisk{
			Registration: registration, Approvals: approvals, Attempts: attempts, Lifecycles: lifecycleDisks,
			EncodedByteLength: encodedLength, CohortDigest: projection.Digest(),
		}
		projections[cohortIndex] = projection
	}
	return disks, projections, nil
}

func decodeArchiveSegmentV0(data []byte) (loadedArchiveSegmentV0, error) {
	if int64(len(data)) > archivestate.MaxSupervisorArchiveBytes {
		return loadedArchiveSegmentV0{}, errors.New("fixed supervisor archive segment exceeds byte capacity")
	}
	var disk archiveSegmentDiskV0
	if err := decodeOneClosedJSON(data, &disk); err != nil {
		return loadedArchiveSegmentV0{}, errors.New("decode fixed supervisor archive segment")
	}
	if disk.FormatVersion == nil || *disk.FormatVersion != archivestate.SupervisorArchiveFormatV0 {
		return loadedArchiveSegmentV0{}, errors.New("unsupported fixed supervisor archive segment version")
	}
	if len(disk.Cohorts) == 0 || uint64(len(data)) != disk.EncodedByteLength {
		return loadedArchiveSegmentV0{}, errors.New("fixed supervisor archive segment length or cohort collection mismatch")
	}
	state := installationState{InitialState: InitialState{
		InstallationID: disk.InstallationID, SupervisorID: disk.SupervisorID,
		EpochSequence: disk.EpochSequence, EpochDigest: disk.EpochDigest,
		TrustPhase: TrustStable, TimeHighWaterUnixSeconds: disk.DurableTimeHighWater,
	}}
	lifecycles := make([]lifecyclestate.Record, 0)
	var previousSequence v0candidate.PositiveUInt53
	var previousID v0candidate.RegistrationID
	for cohortIndex, cohort := range disk.Cohorts {
		registration := cohort.Registration
		if cohort.Approvals == nil || cohort.Attempts == nil || cohort.Lifecycles == nil ||
			registration.Index.RegistrationID == (v0candidate.RegistrationID{}) {
			return loadedArchiveSegmentV0{}, errors.New("fixed supervisor archive cohort is incomplete")
		}
		if cohortIndex > 0 && (registration.Index.RegistrationSequence < previousSequence ||
			(registration.Index.RegistrationSequence == previousSequence &&
				bytes.Compare(registration.Index.RegistrationID[:], previousID[:]) <= 0)) {
			return loadedArchiveSegmentV0{}, errors.New("fixed supervisor archive cohort order mismatch")
		}
		previousSequence, previousID = registration.Index.RegistrationSequence, registration.Index.RegistrationID
		if !sort.SliceIsSorted(cohort.Approvals, func(left, right int) bool {
			return bytes.Compare(cohort.Approvals[left].ApprovalID[:], cohort.Approvals[right].ApprovalID[:]) < 0
		}) || !sort.SliceIsSorted(cohort.Attempts, func(left, right int) bool {
			return bytes.Compare(cohort.Attempts[left].AttemptID[:], cohort.Attempts[right].AttemptID[:]) < 0
		}) || !sort.SliceIsSorted(cohort.Lifecycles, func(left, right int) bool {
			return bytes.Compare(cohort.Lifecycles[left].Bindings.AttemptID[:], cohort.Lifecycles[right].Bindings.AttemptID[:]) < 0
		}) {
			return loadedArchiveSegmentV0{}, errors.New("fixed supervisor archive record order mismatch")
		}
		state.Registrations = append(state.Registrations, registration)
		state.Approvals = append(state.Approvals, cohort.Approvals...)
		state.Attempts = append(state.Attempts, cohort.Attempts...)
		for _, lifecycleDisk := range cohort.Lifecycles {
			record, err := restoreLifecycleRecord(lifecycleDisk, disk.DurableTimeHighWater)
			if err != nil {
				return loadedArchiveSegmentV0{}, errors.New("fixed supervisor archive lifecycle record is invalid")
			}
			lifecycles = append(lifecycles, record)
		}
	}
	state.LastRegistrationSequence = 0
	for _, registration := range state.Registrations {
		if v0candidate.UInt53(registration.Index.RegistrationSequence) > state.LastRegistrationSequence {
			state.LastRegistrationSequence = v0candidate.UInt53(registration.Index.RegistrationSequence)
		}
	}
	if err := recomputeAuthoritySetDigests(&state); err != nil {
		return loadedArchiveSegmentV0{}, err
	}
	if err := validateV1State(state, lifecycles, lifecycleSetDigest(lifecycles)); err != nil {
		return loadedArchiveSegmentV0{}, fmt.Errorf("fixed supervisor archive record cross-link mismatch: %w", err)
	}
	storedDerived, err := archiveIndexesFromDisk(disk.DerivedIndexes)
	if err != nil {
		return loadedArchiveSegmentV0{}, errors.New("fixed supervisor archive derived index is invalid")
	}
	provisional := loadedArchiveSegmentV0{
		Disk: disk, State: state, Lifecycles: lifecycles,
		EffectTombstones: cloneEffectTombstones(storedDerived.View().Effects), Bytes: bytes.Clone(data),
	}
	archivedIndexes, err := archiveIndexesForLoadedSegment(provisional)
	if err != nil {
		return loadedArchiveSegmentV0{}, err
	}
	derivedView := archivedIndexes.View()
	derivedView.Scope = archivestate.ArchiveIndexScopeSegmentDerived
	derivedView.Registrations = []archivestate.RegistrationIndexEntry{}
	derivedView.Approvals = []archivestate.ApprovalIndexEntry{}
	derivedView.Attempts = []archivestate.AttemptIndexEntry{}
	derived, err := archivestate.NewArchiveIndexes(derivedView)
	if err != nil {
		return loadedArchiveSegmentV0{}, err
	}
	if !reflect.DeepEqual(storedDerived.View(), derived.View()) {
		return loadedArchiveSegmentV0{}, errors.New("fixed supervisor archive derived index mismatch")
	}
	_, projections, err := segmentCohorts(state, lifecycles)
	if err != nil {
		return loadedArchiveSegmentV0{}, err
	}
	for index, projection := range projections {
		if projection.Digest() != disk.Cohorts[index].CohortDigest ||
			projection.View().EncodedByteLength != disk.Cohorts[index].EncodedByteLength {
			return loadedArchiveSegmentV0{}, errors.New("fixed supervisor archive cohort digest mismatch")
		}
	}
	segment, err := archivestate.NewArchiveSegment(archivestate.ArchiveSegmentView{
		FormatVersion: *disk.FormatVersion, Ordinal: disk.Ordinal,
		InstallationID: disk.InstallationID, SupervisorID: disk.SupervisorID,
		EpochSequence: disk.EpochSequence, EpochDigest: disk.EpochDigest,
		SourceSnapshotGeneration: disk.SourceSnapshotGeneration, ArchiveGeneration: disk.ArchiveGeneration,
		DurableTimeHighWater: disk.DurableTimeHighWater, PriorCheckpointDigest: disk.PriorCheckpointDigest,
		Cohorts: projections, DerivedIndexes: derived, SetDigests: disk.SetDigests,
		Counts: disk.Counts, EncodedByteLength: disk.EncodedByteLength,
	})
	wantSegmentDigest, digestErr := decodeArchiveSegmentDigest(disk.SegmentDigest)
	if err != nil || digestErr != nil || segment.Digest() != wantSegmentDigest {
		return loadedArchiveSegmentV0{}, errors.New("fixed supervisor archive segment digest mismatch")
	}
	provisional.Segment = segment
	return provisional, nil
}

func decodeArchiveSegmentDigest(encoded string) (archivestate.ArchiveSegmentDigest, error) {
	var digest archivestate.ArchiveSegmentDigest
	if len(encoded) != hex.EncodedLen(len(digest)) {
		return digest, errors.New("fixed supervisor archive segment digest encoding length mismatch")
	}
	decoded, err := hex.DecodeString(encoded)
	if err != nil || len(decoded) != len(digest) {
		return digest, errors.New("fixed supervisor archive segment digest encoding is invalid")
	}
	copy(digest[:], decoded)
	return digest, nil
}

func archiveIndexesForLoadedSegment(segment loadedArchiveSegmentV0) (archivestate.ArchiveIndexes, error) {
	hotTombstones, err := effectTombstonesAtHotAttemptLocations(segment.State, segment.EffectTombstones)
	if err != nil {
		return archivestate.ArchiveIndexes{}, err
	}
	base, err := reconstructV2Indexes(segment.State, segment.Lifecycles, hotTombstones)
	if err != nil {
		return archivestate.ArchiveIndexes{}, err
	}
	view := base.View()
	view.Scope = archivestate.ArchiveIndexScopeRetainedGlobal
	type locations struct {
		registration map[v0candidate.RegistrationID]archivestate.RecordLocation
		approval     map[approvalattempt.ApprovalID]archivestate.RecordLocation
		attempt      map[approvalattempt.AttemptID]archivestate.RecordLocation
		lifecycle    map[approvalattempt.AttemptID]archivestate.RecordLocation
	}
	loc := locations{
		registration: make(map[v0candidate.RegistrationID]archivestate.RecordLocation),
		approval:     make(map[approvalattempt.ApprovalID]archivestate.RecordLocation),
		attempt:      make(map[approvalattempt.AttemptID]archivestate.RecordLocation),
		lifecycle:    make(map[approvalattempt.AttemptID]archivestate.RecordLocation),
	}
	for cohortIndex, cohort := range segment.Disk.Cohorts {
		archiveLocation := func(recordOrdinal int) archivestate.ArchiveLocation {
			return archivestate.ArchiveLocation{SegmentOrdinal: segment.Disk.Ordinal,
				CohortOrdinal: archivestate.ArchiveOrdinal(cohortIndex + 1), //nolint:gosec // G115: segment validation caps cohorts at 256.
				RecordOrdinal: archivestate.ArchiveOrdinal(recordOrdinal)}   //nolint:gosec // G115: segment validation caps records at 4,096.
		}
		location, locationErr := archivestate.NewArchiveRecordLocation(archivestate.RecordRegistration, archiveLocation(1))
		if locationErr != nil {
			return archivestate.ArchiveIndexes{}, locationErr
		}
		loc.registration[cohort.Registration.Index.RegistrationID] = location
		for index, approval := range cohort.Approvals {
			location, locationErr = archivestate.NewArchiveRecordLocation(archivestate.RecordApproval, archiveLocation(index+1))
			if locationErr != nil {
				return archivestate.ArchiveIndexes{}, locationErr
			}
			loc.approval[approval.ApprovalID] = location
		}
		for index, attempt := range cohort.Attempts {
			location, locationErr = archivestate.NewArchiveRecordLocation(archivestate.RecordAttempt, archiveLocation(index+1))
			if locationErr != nil {
				return archivestate.ArchiveIndexes{}, locationErr
			}
			loc.attempt[attempt.AttemptID] = location
		}
		for index, lifecycle := range cohort.Lifecycles {
			location, locationErr = archivestate.NewArchiveRecordLocation(archivestate.RecordLifecycle, archiveLocation(index+1))
			if locationErr != nil {
				return archivestate.ArchiveIndexes{}, locationErr
			}
			loc.lifecycle[lifecycle.Bindings.AttemptID] = location
		}
	}
	for index := range view.Registrations {
		view.Registrations[index].Location = loc.registration[view.Registrations[index].RegistrationID]
	}
	for index := range view.Approvals {
		view.Approvals[index].Location = loc.approval[view.Approvals[index].ApprovalID]
	}
	for index := range view.Attempts {
		entry := &view.Attempts[index]
		entry.Location = loc.attempt[entry.AttemptID]
		lifecycleView := entry.Lifecycle.View()
		if lifecycleView.Presence == archivestate.AttemptLifecyclePresent {
			entry.Lifecycle, err = archivestate.NewPresentAttemptLifecycle(lifecycleView.State, loc.lifecycle[entry.AttemptID], lifecycleView.FullRecordDigest)
			if err != nil {
				return archivestate.ArchiveIndexes{}, err
			}
		}
	}
	for index := range view.Nonces {
		view.Nonces[index].Location = loc.approval[view.Nonces[index].ApprovalID]
	}
	for index := range view.Effects {
		view.Effects[index].AttemptLocation = loc.attempt[view.Effects[index].AttemptID]
	}
	for index := range view.Instances {
		view.Instances[index].Location = loc.lifecycle[view.Instances[index].AttemptID]
	}
	for index := range view.ApprovalReplay {
		view.ApprovalReplay[index].Location = loc.approval[view.ApprovalReplay[index].ApprovalID]
	}
	for index := range view.AttemptReplay {
		view.AttemptReplay[index].Location = loc.attempt[view.AttemptReplay[index].AttemptID]
	}
	return archivestate.NewArchiveIndexes(view)
}

func effectTombstonesAtHotAttemptLocations(
	state installationState,
	entries []archivestate.EffectIndexEntry,
) ([]archivestate.EffectIndexEntry, error) {
	attempts := append([]approvalattempt.ExecutionAttempt(nil), state.Attempts...)
	sort.Slice(attempts, func(left, right int) bool {
		return bytes.Compare(attempts[left].AttemptID[:], attempts[right].AttemptID[:]) < 0
	})
	locations := make(map[approvalattempt.AttemptID]archivestate.RecordLocation, len(attempts))
	for index, attempt := range attempts {
		location, err := archivestate.NewHotRecordLocation(archivestate.RecordAttempt, v0candidate.PositiveUInt53(index+1))
		if err != nil {
			return nil, err
		}
		locations[attempt.AttemptID] = location
	}
	result := cloneEffectTombstones(entries)
	for index := range result {
		location, ok := locations[result[index].AttemptID]
		if !ok {
			return nil, errors.New("archive effect tombstone attempt is missing")
		}
		result[index].AttemptLocation = location
	}
	return result, nil
}

func reconstructV2IndexesForWorld(
	hot installationState,
	hotLifecycles []lifecyclestate.Record,
	hotEffectTombstones []archivestate.EffectIndexEntry,
	segments []loadedArchiveSegmentV0,
) (archivestate.ArchiveIndexes, error) {
	hotIndexes, err := reconstructV2Indexes(hot, hotLifecycles, hotEffectTombstones)
	if err != nil {
		return archivestate.ArchiveIndexes{}, err
	}
	view := hotIndexes.View()
	for _, segment := range segments {
		archived, archiveErr := archiveIndexesForLoadedSegment(segment)
		if archiveErr != nil {
			return archivestate.ArchiveIndexes{}, archiveErr
		}
		archivedView := archived.View()
		view.Registrations = append(view.Registrations, archivedView.Registrations...)
		view.Approvals = append(view.Approvals, archivedView.Approvals...)
		view.Attempts = append(view.Attempts, archivedView.Attempts...)
		view.Nonces = append(view.Nonces, archivedView.Nonces...)
		view.Effects = append(view.Effects, archivedView.Effects...)
		view.Instances = append(view.Instances, archivedView.Instances...)
		view.ApprovalReplay = append(view.ApprovalReplay, archivedView.ApprovalReplay...)
		view.AttemptReplay = append(view.AttemptReplay, archivedView.AttemptReplay...)
	}
	sortArchiveIndexesView(&view)
	return archivestate.NewArchiveIndexes(view)
}

func sortArchiveIndexesView(view *archivestate.ArchiveIndexesView) {
	sort.Slice(view.Registrations, func(left, right int) bool {
		return bytes.Compare(view.Registrations[left].RegistrationID[:], view.Registrations[right].RegistrationID[:]) < 0
	})
	sort.Slice(view.Approvals, func(left, right int) bool {
		return bytes.Compare(view.Approvals[left].ApprovalID[:], view.Approvals[right].ApprovalID[:]) < 0
	})
	sort.Slice(view.Attempts, func(left, right int) bool {
		return bytes.Compare(view.Attempts[left].AttemptID[:], view.Attempts[right].AttemptID[:]) < 0
	})
	sort.Slice(view.Nonces, func(left, right int) bool {
		return bytes.Compare(view.Nonces[left].AttemptNonce[:], view.Nonces[right].AttemptNonce[:]) < 0
	})
	sort.Slice(view.Effects, func(left, right int) bool {
		return bytes.Compare(view.Effects[left].EffectID[:], view.Effects[right].EffectID[:]) < 0
	})
	sort.Slice(view.Instances, func(left, right int) bool {
		return bytes.Compare(view.Instances[left].InstanceDigest[:], view.Instances[right].InstanceDigest[:]) < 0
	})
	sort.Slice(view.ApprovalReplay, func(left, right int) bool {
		comparison := bytes.Compare(view.ApprovalReplay[left].PayloadDigest[:], view.ApprovalReplay[right].PayloadDigest[:])
		if comparison != 0 {
			return comparison < 0
		}
		return bytes.Compare(view.ApprovalReplay[left].AuthorizationIdentity[:], view.ApprovalReplay[right].AuthorizationIdentity[:]) < 0
	})
	sort.Slice(view.AttemptReplay, func(left, right int) bool {
		comparison := bytes.Compare(view.AttemptReplay[left].RegistrationID[:], view.AttemptReplay[right].RegistrationID[:])
		if comparison != 0 {
			return comparison < 0
		}
		return bytes.Compare(view.AttemptReplay[left].ApprovalID[:], view.AttemptReplay[right].ApprovalID[:]) < 0
	})
}

func reconstructMigrationGenesisIndexes(
	hot installationState,
	hotLifecycles []lifecyclestate.Record,
	segments []loadedArchiveSegmentV0,
) (archivestate.ArchiveIndexes, error) {
	world := cloneState(hot)
	lifecycles := cloneLifecycleRecords(hotLifecycles, hot.TimeHighWaterUnixSeconds)
	for _, segment := range segments {
		world.Registrations = append(world.Registrations, segment.State.Registrations...)
		world.Approvals = append(world.Approvals, segment.State.Approvals...)
		world.Attempts = append(world.Attempts, segment.State.Attempts...)
		lifecycles = append(lifecycles, segment.Lifecycles...)
	}
	sort.Slice(lifecycles, func(left, right int) bool {
		leftID := lifecycles[left].Bindings().View().AttemptID
		rightID := lifecycles[right].Bindings().View().AttemptID
		return bytes.Compare(leftID[:], rightID[:]) < 0
	})
	// Migration-genesis locations are all-hot, regardless of where F3 later
	// carried the same immutable seed. Rebuild those typed locations exactly.
	seedTombstones, err := visibleV1EffectTombstones(world, lifecycles)
	if err != nil {
		return archivestate.ArchiveIndexes{}, err
	}
	return reconstructV2Indexes(world, lifecycles, seedTombstones)
}

func reconstructMigrationGenesisHotSetDigests(
	hot installationState,
	hotLifecycles []lifecyclestate.Record,
	segments []loadedArchiveSegmentV0,
) (archivestate.HotSetDigests, error) {
	world := cloneState(hot)
	lifecycles := cloneLifecycleRecords(hotLifecycles, hot.TimeHighWaterUnixSeconds)
	for _, segment := range segments {
		world.Registrations = append(world.Registrations, segment.State.Registrations...)
		world.Approvals = append(world.Approvals, segment.State.Approvals...)
		world.Attempts = append(world.Attempts, segment.State.Attempts...)
		lifecycles = append(lifecycles, segment.Lifecycles...)
	}
	sort.Slice(world.Registrations, func(left, right int) bool {
		return world.Registrations[left].Index.RegistrationSequence < world.Registrations[right].Index.RegistrationSequence
	})
	if err := recomputeAuthoritySetDigests(&world); err != nil {
		return archivestate.HotSetDigests{}, err
	}
	return archivestate.HotSetDigests{
		Registrations: world.RegistrationSetDigest, Approvals: world.ApprovalSetDigest,
		Attempts: world.AttemptSetDigest, Lifecycles: [32]byte(lifecycleSetDigest(lifecycles)),
	}, nil
}

func validateActivationCheckpointDisk(
	envelope diskEnvelopeV2,
	indexes archivestate.ArchiveIndexes,
	segment loadedArchiveSegmentV0,
) (archivestate.ArchiveCheckpoint, error) {
	disk := envelope.ActivationCheckpoint
	if disk == nil || disk.ArchiveIndexDigest != indexes.CombinedDigest() ||
		disk.PreviousCheckpoint != envelope.PreviousCheckpoint || disk.NewSegmentDigest != segment.Segment.Digest() ||
		disk.DescriptorSetDigest != envelope.DescriptorSetDigest ||
		disk.SourceSnapshotGeneration+1 != disk.ResultSnapshotGeneration ||
		disk.ResultSnapshotGeneration > envelope.SnapshotGeneration || disk.ArchiveGeneration != envelope.ArchiveGeneration ||
		disk.InstallationID != envelope.State.InstallationID || disk.SupervisorID != envelope.State.SupervisorID ||
		disk.EpochSequence != envelope.State.EpochSequence || disk.EpochDigest != envelope.State.EpochDigest ||
		disk.DurableTimeHighWater > envelope.State.TimeHighWaterUnixSeconds {
		return archivestate.ArchiveCheckpoint{}, errors.New("fixed supervisor archive checkpoint cross-link mismatch")
	}
	if disk.ResultSnapshotGeneration == envelope.SnapshotGeneration {
		wantHot := archivestate.HotSetDigests{
			Registrations: envelope.State.RegistrationSetDigest, Approvals: envelope.State.ApprovalSetDigest,
			Attempts: envelope.State.AttemptSetDigest, Lifecycles: [32]byte(envelope.LifecycleSetDigest),
		}
		if disk.HotSetDigests != wantHot {
			return archivestate.ArchiveCheckpoint{}, errors.New("fixed supervisor archive checkpoint hot-set mismatch")
		}
	}
	checkpoint, err := archivestate.NewArchiveCheckpoint(archivestate.ArchiveCheckpointView{
		PreviousCheckpoint: disk.PreviousCheckpoint, NewSegmentDigest: disk.NewSegmentDigest,
		DescriptorSetDigest: disk.DescriptorSetDigest, ArchiveIndexDigest: disk.ArchiveIndexDigest,
		HotSetDigests: disk.HotSetDigests, SourceSnapshotGeneration: disk.SourceSnapshotGeneration,
		ResultSnapshotGeneration: disk.ResultSnapshotGeneration, ArchiveGeneration: disk.ArchiveGeneration,
		InstallationID: disk.InstallationID, SupervisorID: disk.SupervisorID,
		EpochSequence: disk.EpochSequence, EpochDigest: disk.EpochDigest,
		DurableTimeHighWater: disk.DurableTimeHighWater,
	})
	if err != nil || checkpoint.Digest() != disk.Digest || checkpoint.Reference() != envelope.CurrentCheckpoint {
		return archivestate.ArchiveCheckpoint{}, errors.New("fixed supervisor archive checkpoint digest mismatch")
	}
	return checkpoint, nil
}

func archiveRoot(storePath string) string { return storePath + ".archive" }

func archiveSegmentName(digest archivestate.ArchiveSegmentDigest) string {
	return "segment-" + hex.EncodeToString(digest[:]) + ".json"
}

func archiveSegmentPath(storePath string, digest archivestate.ArchiveSegmentDigest) string {
	return filepath.Join(archiveRoot(storePath), archiveSegmentName(digest))
}

func loadArchiveSegmentsForEnvelope(
	storePath string,
	envelope diskEnvelopeV2,
) ([]loadedArchiveSegmentV0, uint16, error) {
	root := archiveRoot(storePath)
	info, err := os.Lstat(root) //nolint:gosec // G703: derived solely from the installation-configured store path.
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && len(envelope.Descriptors) == 0 {
			return []loadedArchiveSegmentV0{}, 0, nil
		}
		return nil, 0, errors.New("fixed supervisor archive directory is missing")
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() || info.Mode().Perm() != 0o700 {
		return nil, 0, errors.New("fixed supervisor archive directory shape is unsafe")
	}
	entries, err := os.ReadDir(root) //nolint:gosec // G703: trusted derived archive directory.
	if err != nil {
		return nil, 0, errors.New("read fixed supervisor archive directory")
	}
	referenced := make(map[string]archivestate.ArchiveDescriptorView, len(envelope.Descriptors))
	for _, descriptor := range envelope.Descriptors {
		referenced[archiveSegmentName(descriptor.SegmentDigest)] = descriptor
	}
	segmentsByOrdinal := make(map[archivestate.ArchiveOrdinal]loadedArchiveSegmentV0, len(envelope.Descriptors))
	orphans := 0
	for _, entry := range entries {
		name := entry.Name()
		path := filepath.Join(root, name)
		if strings.HasPrefix(name, ".capsule-supervisor-archive-") && strings.HasSuffix(name, ".tmp") {
			if err := validateArchiveTemporaryFileShape(path); err != nil {
				return nil, 0, err
			}
			orphans++
			continue
		}
		descriptor, isReferenced := referenced[name]
		if !strings.HasPrefix(name, "segment-") || !strings.HasSuffix(name, ".json") {
			return nil, 0, errors.New("fixed supervisor archive directory contains unknown entry")
		}
		if err := validateArchiveFileShape(path, 0o400); err != nil {
			return nil, 0, err
		}
		data, readErr := readBoundedStoreFile(path, archivestate.MaxSupervisorArchiveBytes)
		if readErr != nil {
			return nil, 0, readErr
		}
		segment, decodeErr := decodeArchiveSegmentV0(data)
		if decodeErr != nil {
			return nil, 0, decodeErr
		}
		if name != archiveSegmentName(segment.Segment.Digest()) {
			return nil, 0, errors.New("fixed supervisor archive digest-addressed name mismatch")
		}
		if !isReferenced {
			orphans++
			continue
		}
		segmentView := segment.Segment.View()
		if descriptor.Ordinal != segmentView.Ordinal ||
			descriptor.SourceSnapshotGeneration != segmentView.SourceSnapshotGeneration ||
			descriptor.ArchiveGeneration != segmentView.ArchiveGeneration ||
			descriptor.SegmentDigest != segment.Segment.Digest() ||
			descriptor.EncodedByteLength != uint64(len(data)) || descriptor.Counts != segmentView.Counts ||
			segmentView.InstallationID != envelope.State.InstallationID ||
			segmentView.SupervisorID != envelope.State.SupervisorID ||
			segmentView.EpochSequence != envelope.State.EpochSequence ||
			segmentView.EpochDigest != envelope.State.EpochDigest ||
			segmentView.DurableTimeHighWater > envelope.State.TimeHighWaterUnixSeconds {
			return nil, 0, errors.New("fixed supervisor archive descriptor cross-link mismatch")
		}
		if _, duplicate := segmentsByOrdinal[descriptor.Ordinal]; duplicate {
			return nil, 0, errors.New("fixed supervisor archive ordinal collision")
		}
		segmentsByOrdinal[descriptor.Ordinal] = segment
	}
	if orphans > int(^uint16(0)) {
		return nil, 0, errors.New("fixed supervisor archive orphan report overflow")
	}
	segments := make([]loadedArchiveSegmentV0, 0, len(envelope.Descriptors))
	for _, descriptor := range envelope.Descriptors {
		segment, exists := segmentsByOrdinal[descriptor.Ordinal]
		if !exists {
			return nil, 0, errors.New("fixed supervisor archive referenced segment is missing")
		}
		segments = append(segments, segment)
	}
	return segments, uint16(orphans), nil
}

func validateArchiveTemporaryFileShape(path string) error {
	info, err := os.Lstat(path) //nolint:gosec // G703: exact child of trusted derived archive root.
	if err != nil {
		return errors.New("inspect fixed supervisor archive temporary file")
	}
	mode := info.Mode().Perm()
	if mode != 0o600 && mode != 0o400 {
		return errors.New("fixed supervisor archive temporary file mode is unsafe")
	}
	return validateArchiveFileShape(path, mode)
}

func validateArchiveFileShape(path string, mode os.FileMode) error {
	info, err := os.Lstat(path) //nolint:gosec // G703: exact child of trusted derived archive root.
	if err != nil {
		return errors.New("inspect fixed supervisor archive file")
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() || info.Mode().Perm() != mode {
		return errors.New("fixed supervisor archive file shape is unsafe")
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok || stat.Nlink != 1 || int(stat.Uid) != os.Geteuid() {
		return errors.New("fixed supervisor archive file owner or link count is unsafe")
	}
	return nil
}

func verifyPreparedAgainstLoaded(loaded loadedV2State, prepared PreparedFixedStoreV2Archive) error {
	if prepared.plan.plan.View().SourceSnapshotGeneration != loaded.Active.View().SnapshotGeneration ||
		prepared.plan.plan.View().ArchiveGeneration != loaded.Active.View().ArchiveGeneration ||
		prepared.plan.plan.View().CheckpointHead != loaded.Active.View().CurrentCheckpoint.Digest {
		return ErrArchiveStaleTransaction
	}
	decoded, err := decodeArchiveSegmentV0(prepared.segmentBytes)
	if err != nil {
		return fmt.Errorf("verify prepared fixed supervisor archive: %w", err)
	}
	if decoded.Segment.Digest() != prepared.segmentDigest || !bytes.Equal(decoded.Bytes, prepared.segmentBytes) {
		return errors.New("prepared fixed supervisor archive byte or digest mismatch")
	}
	encodedCandidate, err := encodeEnvelopeV2(prepared.candidate)
	if err != nil {
		return err
	}
	var candidate diskEnvelopeV2
	if err := decodeOneClosedJSON(encodedCandidate, &candidate); err != nil {
		return errors.New("prepared fixed supervisor active successor is not closed")
	}
	if candidate.Lifecycles == nil || len(candidate.Descriptors) != 1 || candidate.ActivationCheckpoint == nil {
		return errors.New("prepared fixed supervisor active successor is incomplete")
	}
	records := make([]lifecyclestate.Record, len(candidate.Lifecycles))
	for index, disk := range candidate.Lifecycles {
		record, restoreErr := restoreLifecycleRecord(disk, candidate.State.TimeHighWaterUnixSeconds)
		if restoreErr != nil {
			return errors.New("prepared fixed supervisor active lifecycle is invalid")
		}
		records[index] = record
	}
	if err := validateV1State(candidate.State, records, candidate.LifecycleSetDigest); err != nil {
		return err
	}
	storedIndexes, err := archiveIndexesFromDisk(candidate.Indexes)
	if err != nil {
		return err
	}
	candidateHotTombstones, _, _, err := effectTombstonesFromEnvelope(candidate, storedIndexes)
	if err != nil {
		return err
	}
	reconstructed, err := reconstructV2IndexesForWorld(candidate.State, records, candidateHotTombstones, []loadedArchiveSegmentV0{decoded})
	if err != nil || !reflect.DeepEqual(storedIndexes.View(), reconstructed.View()) {
		return errors.New("prepared fixed supervisor archive world index mismatch")
	}
	genesisIndexes, err := reconstructMigrationGenesisIndexes(candidate.State, records, []loadedArchiveSegmentV0{decoded})
	if err != nil {
		return err
	}
	genesisHotDigests, err := reconstructMigrationGenesisHotSetDigests(candidate.State, records, []loadedArchiveSegmentV0{decoded})
	if err != nil {
		return err
	}
	genesis, err := validateMigrationGenesisDisk(candidate, genesisIndexes, genesisHotDigests)
	if err != nil || candidate.PreviousCheckpoint != genesis.Reference() {
		return fmt.Errorf("prepared fixed supervisor archive genesis mismatch: %w", err)
	}
	descriptor, err := archivestate.NewArchiveDescriptor(candidate.Descriptors[0])
	if err != nil {
		return err
	}
	active, err := archivestate.NewActiveStateV2(archivestate.ActiveStateV2View{
		StoreFormatVersion: *candidate.StoreFormatVersion, MigrationSourceVersion: *candidate.MigrationSourceVersion,
		SnapshotGeneration: candidate.SnapshotGeneration, ArchiveGeneration: candidate.ArchiveGeneration,
		InstallationID: candidate.State.InstallationID, SupervisorID: candidate.State.SupervisorID,
		EpochSequence: candidate.State.EpochSequence, EpochDigest: candidate.State.EpochDigest,
		DurableTimeHighWater: candidate.State.TimeHighWaterUnixSeconds,
		Descriptors:          []archivestate.ArchiveDescriptor{descriptor}, DescriptorSetDigest: candidate.DescriptorSetDigest,
		Indexes: storedIndexes, IndexDigests: candidate.IndexDigests, CombinedIndexDigest: candidate.CombinedIndexDigest,
		EffectTombstoneCoverage:   candidate.EffectTombstoneCoverage,
		VisibleV1EffectSeedCount:  candidate.VisibleV1EffectSeedCount,
		VisibleV1EffectSeedDigest: candidate.VisibleV1EffectSeedDigest,
		PreviousCheckpoint:        candidate.PreviousCheckpoint, CurrentCheckpoint: candidate.CurrentCheckpoint,
		HotCounts: candidate.HotCounts, ArchivedCounts: candidate.ArchivedCounts, TotalCounts: candidate.TotalCounts,
	})
	if err != nil || active.View().CurrentCheckpoint != candidate.CurrentCheckpoint {
		return errors.New("prepared fixed supervisor active projection mismatch")
	}
	if _, err := validateActivationCheckpointDisk(candidate, storedIndexes, decoded); err != nil {
		return err
	}
	return nil
}

func ensureArchiveRoot(storePath string) (string, error) {
	root := archiveRoot(storePath)
	created := false
	if err := os.Mkdir(root, 0o700); err != nil && !errors.Is(err, os.ErrExist) { //nolint:gosec // G703: derived installation store sibling.
		return "", errors.New("create fixed supervisor archive directory")
	} else if err == nil {
		created = true
	}
	info, err := os.Lstat(root) //nolint:gosec // G703: derived installation store sibling.
	if err != nil || info.Mode()&os.ModeSymlink != 0 || !info.IsDir() || info.Mode().Perm() != 0o700 {
		return "", errors.New("fixed supervisor archive directory shape is unsafe")
	}
	if created {
		parent, openErr := os.Open(filepath.Dir(root)) //nolint:gosec // G304: trusted store parent.
		if openErr != nil {
			return "", errors.New("open fixed supervisor archive parent directory")
		}
		syncErr, closeErr := parent.Sync(), parent.Close()
		if syncErr != nil || closeErr != nil {
			return "", errors.New("sync fixed supervisor archive parent directory")
		}
	}
	return root, nil
}

func publishArchiveSegment(
	ctx context.Context,
	storePath string,
	owner ArchiveOwner,
	prepared PreparedFixedStoreV2Archive,
	faults ArchiveFaultInjector,
) (string, error) {
	root, err := ensureArchiveRoot(storePath)
	if err != nil {
		return "", err
	}
	finalPath := archiveSegmentPath(storePath, prepared.segmentDigest)
	if existing, statErr := os.Lstat(finalPath); statErr == nil {
		if existing.Mode()&os.ModeSymlink != 0 {
			return "", errors.New("fixed supervisor archive final segment is a symlink")
		}
		if err := validateArchiveFileShape(finalPath, 0o400); err != nil {
			return "", err
		}
		data, readErr := readBoundedStoreFile(finalPath, archivestate.MaxSupervisorArchiveBytes)
		if readErr != nil || !bytes.Equal(data, prepared.segmentBytes) {
			return "", errors.New("fixed supervisor archive existing segment substitution")
		}
		if err := owner.CheckHeld(ctx); err != nil {
			return "", fmt.Errorf("%w: before existing segment reuse: %v", ErrArchiveOwnerRequired, err)
		}
		if owner.OwnerSessionID() != prepared.plan.owner {
			return "", ErrArchiveOwnerSessionMismatch
		}
		directory, openErr := os.Open(root) //nolint:gosec // G304: trusted derived archive directory.
		if openErr != nil {
			return "", fmt.Errorf("%w: reopen archive directory", ErrArchiveOutcomeIndeterminate)
		}
		syncErr, closeErr := directory.Sync(), directory.Close()
		if syncErr != nil || closeErr != nil {
			return "", fmt.Errorf("%w: resync archive directory", ErrArchiveOutcomeIndeterminate)
		}
		return finalPath, nil
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return "", errors.New("inspect fixed supervisor archive final segment")
	}
	temporary, err := os.CreateTemp(root, ".capsule-supervisor-archive-*.tmp")
	if err != nil {
		return "", errors.New("create fixed supervisor archive temporary segment")
	}
	temporaryPath := temporary.Name()
	published := false
	defer func() {
		_ = temporary.Close()
		if !published {
			_ = os.Remove(temporaryPath) //nolint:gosec // G703: exact path returned by CreateTemp under trusted root.
		}
	}()
	if err := temporary.Chmod(0o600); err != nil {
		return "", errors.New("protect fixed supervisor archive temporary segment")
	}
	if err := failArchive(faults, FaultArchiveAfterSegmentTempCreate, false); err != nil {
		return "", err
	}
	if _, err := temporary.Write(prepared.segmentBytes); err != nil {
		return "", errors.New("write fixed supervisor archive temporary segment")
	}
	if err := failArchive(faults, FaultArchiveAfterSegmentTempWrite, false); err != nil {
		return "", err
	}
	if err := temporary.Sync(); err != nil {
		return "", errors.New("sync fixed supervisor archive temporary segment")
	}
	if err := failArchive(faults, FaultArchiveAfterSegmentTempSync, false); err != nil {
		return "", err
	}
	if err := temporary.Close(); err != nil {
		return "", errors.New("close fixed supervisor archive temporary segment")
	}
	if err := validateArchiveFileShape(temporaryPath, 0o600); err != nil {
		return "", err
	}
	reopened, err := readBoundedStoreFile(temporaryPath, archivestate.MaxSupervisorArchiveBytes)
	if err != nil || !bytes.Equal(reopened, prepared.segmentBytes) {
		return "", errors.New("reopened fixed supervisor archive temporary segment mismatch")
	}
	if decoded, decodeErr := decodeArchiveSegmentV0(reopened); decodeErr != nil || decoded.Segment.Digest() != prepared.segmentDigest {
		return "", errors.New("reopened fixed supervisor archive temporary segment verification failed")
	}
	if err := os.Chmod(temporaryPath, 0o400); err != nil { //nolint:gosec // G703: exact CreateTemp path.
		return "", errors.New("make fixed supervisor archive segment immutable")
	}
	immutable, err := os.Open(temporaryPath) //nolint:gosec // G304: exact CreateTemp path.
	if err != nil {
		return "", errors.New("reopen immutable fixed supervisor archive segment")
	}
	syncErr, closeErr := immutable.Sync(), immutable.Close()
	if syncErr != nil || closeErr != nil {
		return "", errors.New("sync immutable fixed supervisor archive segment")
	}
	if err := failArchive(faults, FaultArchiveBeforeSegmentPublish, false); err != nil {
		return "", err
	}
	if err := owner.CheckHeld(ctx); err != nil {
		return "", fmt.Errorf("%w: before segment publication: %v", ErrArchiveOwnerRequired, err)
	}
	if owner.OwnerSessionID() != prepared.plan.owner {
		return "", ErrArchiveOwnerSessionMismatch
	}
	if err := os.Rename(temporaryPath, finalPath); err != nil { //nolint:gosec // G703: exact temp and digest-derived final path.
		return "", errors.New("publish fixed supervisor archive segment")
	}
	published = true
	if err := failArchive(faults, FaultArchiveAfterSegmentPublish, true); err != nil {
		return "", err
	}
	directory, err := os.Open(root) //nolint:gosec // G304: trusted derived archive directory.
	if err != nil {
		return "", fmt.Errorf("%w: open archive directory", ErrArchiveOutcomeIndeterminate)
	}
	syncErr, closeErr = directory.Sync(), directory.Close()
	if syncErr != nil || closeErr != nil {
		return "", fmt.Errorf("%w: sync archive directory", ErrArchiveOutcomeIndeterminate)
	}
	if err := failArchive(faults, FaultArchiveAfterSegmentDirSync, true); err != nil {
		return "", err
	}
	return finalPath, nil
}

func commitArchiveActivation(
	ctx context.Context,
	storePath string,
	activeBytes []byte,
	owner ArchiveOwner,
	expectedSession lifecyclestate.OwnerSessionID,
	faults ArchiveFaultInjector,
) error {
	parent := filepath.Dir(storePath)
	temporary, err := os.CreateTemp(parent, ".capsule-supervisor-v2-archive-activation-*.tmp")
	if err != nil {
		return errors.New("create fixed supervisor archive activation")
	}
	temporaryPath := temporary.Name()
	committed := false
	defer func() {
		_ = temporary.Close()
		if !committed {
			_ = os.Remove(temporaryPath) //nolint:gosec // G703: exact CreateTemp path.
		}
	}()
	if err := temporary.Chmod(0o600); err != nil {
		return errors.New("protect fixed supervisor archive activation")
	}
	if err := failArchive(faults, FaultArchiveAfterActiveTempCreate, false); err != nil {
		return err
	}
	if _, err := temporary.Write(activeBytes); err != nil {
		return errors.New("write fixed supervisor archive activation")
	}
	if err := failArchive(faults, FaultArchiveAfterActiveTempWrite, false); err != nil {
		return err
	}
	if err := temporary.Sync(); err != nil {
		return errors.New("sync fixed supervisor archive activation")
	}
	if err := failArchive(faults, FaultArchiveAfterActiveTempSync, false); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return errors.New("close fixed supervisor archive activation")
	}
	if err := failArchive(faults, FaultArchiveBeforeActivation, false); err != nil {
		return err
	}
	if err := owner.CheckHeld(ctx); err != nil {
		return fmt.Errorf("%w: immediately before activation: %v", ErrArchiveOwnerRequired, err)
	}
	if owner.OwnerSessionID() != expectedSession {
		return ErrArchiveOwnerSessionMismatch
	}
	if err := os.Rename(temporaryPath, storePath); err != nil { //nolint:gosec // G703: exact temp and installation-configured store path.
		return errors.New("commit fixed supervisor archive activation")
	}
	committed = true
	if err := failArchive(faults, FaultArchiveAfterActivation, true); err != nil {
		return err
	}
	directory, err := os.Open(parent) //nolint:gosec // G304: trusted store parent.
	if err != nil {
		return fmt.Errorf("%w: open activation directory", ErrArchiveOutcomeIndeterminate)
	}
	syncErr, closeErr := directory.Sync(), directory.Close()
	if syncErr != nil || closeErr != nil {
		return fmt.Errorf("%w: sync activation directory", ErrArchiveOutcomeIndeterminate)
	}
	if err := failArchive(faults, FaultArchiveAfterActiveDirSync, true); err != nil {
		return err
	}
	return nil
}

func failArchive(faults ArchiveFaultInjector, point ArchiveFault, indeterminate bool) error {
	if faults == nil {
		return nil
	}
	if err := faults.FailArchiveAt(point); err != nil {
		if indeterminate {
			return fmt.Errorf("%w: injected at %s", ErrArchiveOutcomeIndeterminate, point)
		}
		return fmt.Errorf("fixed supervisor archive confirmed failure at %s", point)
	}
	return nil
}
