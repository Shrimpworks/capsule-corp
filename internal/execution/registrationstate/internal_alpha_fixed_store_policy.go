package registrationstate

import (
	"context"
	"errors"
	"fmt"
	"time"

	"capsule.local/capsule/internal/execution/archivestate"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	// InternalAlphaFixedStorePosture labels this finite policy oracle. It is not
	// a production engine, installed-state claim, or continuity mechanism.
	InternalAlphaFixedStorePosture = "owner-only-disposable-alpha-oracle"

	InternalAlphaMaxCumulativeAttempts  = uint64(128)
	InternalAlphaMaxActiveStoreBytes    = uint64(8 << 20)
	InternalAlphaMaxArchiveSegments     = uint16(16)
	InternalAlphaMaxStartupVerification = 2 * time.Second
	InternalAlphaMaxDurableCommitP95    = 250 * time.Millisecond
)

var ErrInternalAlphaAttemptRefused = errors.New("owner-only internal-alpha fixed-store attempt refused")

type InternalAlphaStoreStopReason string

const (
	InternalAlphaStoreStopUnknownState        InternalAlphaStoreStopReason = "unknown-state"
	InternalAlphaStoreStopCumulativeAttempts  InternalAlphaStoreStopReason = "cumulative-attempts"
	InternalAlphaStoreStopActiveStoreBytes    InternalAlphaStoreStopReason = "active-store-bytes"
	InternalAlphaStoreStopArchiveSegments     InternalAlphaStoreStopReason = "archive-segments"
	InternalAlphaStoreStopStartupVerification InternalAlphaStoreStopReason = "startup-full-verification"
	InternalAlphaStoreStopDurableCommitP95    InternalAlphaStoreStopReason = "durable-commit-p95"
)

// InternalAlphaStartupVerification is a sealed, re-evaluated observation issued
// only after an owner-held full verification. It is not a persistent
// installation trip latch; a later owner session performs and times a new full
// verification. The zero value is unknown and refuses.
type InternalAlphaStartupVerification struct {
	observed       bool
	duration       time.Duration
	ownerSessionID lifecyclestate.OwnerSessionID
	installationID v0candidate.InstallationID
	supervisorID   v0candidate.SupervisorID
	epochSequence  v0candidate.UInt53
	epochDigest    v0candidate.TrustEpochDigest
}

// InternalAlphaDurableCommitP95 is a sealed, re-evaluated policy input. This
// passive slice compares the supplied observation with ADR-0040's threshold but
// does not select or authenticate its source, sample window, aggregation
// lifetime, or persistence. It is not a persistent installation trip latch.
type InternalAlphaDurableCommitP95 struct {
	observed bool
	duration time.Duration
}

func ObserveInternalAlphaDurableCommitP95(duration time.Duration) (InternalAlphaDurableCommitP95, error) {
	if duration < 0 {
		return InternalAlphaDurableCommitP95{}, errors.New("internal-alpha durable-commit p95 cannot be negative")
	}
	return InternalAlphaDurableCommitP95{observed: true, duration: duration}, nil
}

// InternalAlphaStoreAdmissionReport is a fixed, content-free passive checker
// result. Values are observed, never clamped or rewritten. The report is not an
// authority-bearing permit or evidence of a persistent installation trip.
type InternalAlphaStoreAdmissionReport struct {
	Posture             string
	Allowed             bool
	StopReason          InternalAlphaStoreStopReason
	CumulativeAttempts  uint64
	ActiveStoreBytes    uint64
	ArchiveSegments     uint16
	StartupVerification time.Duration
	DurableCommitP95    time.Duration
	SnapshotGeneration  v0candidate.PositiveUInt53
	ArchiveGeneration   archivestate.ArchiveGeneration
	CurrentCheckpoint   archivestate.ArchiveCheckpointReference
}

type internalAlphaStoreFacts struct {
	known                       bool
	cumulativeAttempts          uint64
	activeStoreBytes            uint64
	archiveSegments             uint16
	startupVerification         time.Duration
	startupVerificationObserved bool
	durableCommitP95            time.Duration
	durableCommitP95Observed    bool
	snapshotGeneration          v0candidate.PositiveUInt53
	archiveGeneration           archivestate.ArchiveGeneration
	currentCheckpoint           archivestate.ArchiveCheckpointReference
}

type internalAlphaClock interface {
	Now() time.Time
}

type systemInternalAlphaClock struct{}

func (systemInternalAlphaClock) Now() time.Time { return time.Now() }

// VerifyInternalAlphaFixedStoreStartup performs the ADR-0040 startup full
// verification under the asserted owner and returns a sealed duration
// observation. It performs no store, archive, approval, attempt, or lifecycle
// mutation.
func VerifyInternalAlphaFixedStoreStartup(
	ctx context.Context,
	root StoreRoot,
	owner ArchiveOwner,
) (InternalAlphaStartupVerification, InternalAlphaStoreAdmissionReport, error) {
	return verifyInternalAlphaFixedStoreStartupWithClock(ctx, root, owner, systemInternalAlphaClock{})
}

func verifyInternalAlphaFixedStoreStartupWithClock(
	ctx context.Context,
	root StoreRoot,
	owner ArchiveOwner,
	clock internalAlphaClock,
) (InternalAlphaStartupVerification, InternalAlphaStoreAdmissionReport, error) {
	if clock == nil {
		return InternalAlphaStartupVerification{}, unknownInternalAlphaReport(), fmt.Errorf("%w: startup clock unavailable", ErrInternalAlphaAttemptRefused)
	}
	ownerSession, err := internalAlphaOwnerSession(ctx, owner)
	if err != nil {
		return InternalAlphaStartupVerification{}, unknownInternalAlphaReport(), err
	}
	started := clock.Now()
	loaded, loadErr := loadV2State(root.path)
	duration := clock.Now().Sub(started)
	if duration < 0 {
		return InternalAlphaStartupVerification{}, unknownInternalAlphaReport(), fmt.Errorf("%w: startup clock moved backwards", ErrInternalAlphaAttemptRefused)
	}
	if loadErr != nil {
		return InternalAlphaStartupVerification{}, unknownInternalAlphaReport(), fmt.Errorf("%w: %w: startup full verification: %v", ErrInternalAlphaAttemptRefused, ErrStoreRepairRequired, loadErr)
	}
	postVerifySession, err := internalAlphaOwnerSession(ctx, owner)
	if err != nil {
		return InternalAlphaStartupVerification{}, unknownInternalAlphaReport(), err
	}
	if postVerifySession != ownerSession {
		return InternalAlphaStartupVerification{}, unknownInternalAlphaReport(), internalAlphaRefusal(InternalAlphaStoreStopUnknownState)
	}
	view := loaded.Active.View()
	startup := InternalAlphaStartupVerification{
		observed: true, duration: duration, ownerSessionID: ownerSession,
		installationID: view.InstallationID, supervisorID: view.SupervisorID,
		epochSequence: view.EpochSequence, epochDigest: view.EpochDigest,
	}
	facts := internalAlphaFactsFromLoaded(loaded, startup, InternalAlphaDurableCommitP95{observed: true})
	// Startup does not require a p95 observation yet, so keep that dimension
	// below its threshold only for the startup-duration decision.
	facts.durableCommitP95 = 0
	report := evaluateInternalAlphaFixedStoreFacts(facts)
	if !facts.known || report.StopReason == InternalAlphaStoreStopStartupVerification {
		return InternalAlphaStartupVerification{}, report, internalAlphaRefusal(report.StopReason)
	}
	return startup, report, nil
}

// CheckInternalAlphaFixedStoreAttemptAdmission is the passive, re-evaluated
// logical checker for a future owner-only attempt transaction: full verification
// and all comparisons occur before that transaction. This slice does not wire or
// invoke RequestAttempt, persist an installation trip, or return an
// authority-bearing permit.
func CheckInternalAlphaFixedStoreAttemptAdmission(
	ctx context.Context,
	root StoreRoot,
	owner ArchiveOwner,
	startup InternalAlphaStartupVerification,
	commitP95 InternalAlphaDurableCommitP95,
) (InternalAlphaStoreAdmissionReport, error) {
	ownerSession, err := internalAlphaOwnerSession(ctx, owner)
	if err != nil {
		return unknownInternalAlphaReport(), err
	}
	if !startup.observed || startup.ownerSessionID != ownerSession {
		return unknownInternalAlphaReport(), internalAlphaRefusal(InternalAlphaStoreStopUnknownState)
	}
	loaded, err := loadV2State(root.path)
	if err != nil {
		return unknownInternalAlphaReport(), fmt.Errorf("%w: %w: pre-attempt full verification: %v", ErrInternalAlphaAttemptRefused, ErrStoreRepairRequired, err)
	}
	postVerifySession, err := internalAlphaOwnerSession(ctx, owner)
	if err != nil {
		return unknownInternalAlphaReport(), err
	}
	if postVerifySession != ownerSession {
		return unknownInternalAlphaReport(), internalAlphaRefusal(InternalAlphaStoreStopUnknownState)
	}
	facts := internalAlphaFactsFromLoaded(loaded, startup, commitP95)
	report := evaluateInternalAlphaFixedStoreFacts(facts)
	if !report.Allowed {
		return report, internalAlphaRefusal(report.StopReason)
	}
	return report, nil
}

func internalAlphaFactsFromLoaded(
	loaded loadedV2State,
	startup InternalAlphaStartupVerification,
	commitP95 InternalAlphaDurableCommitP95,
) internalAlphaStoreFacts {
	view := loaded.Active.View()
	known := loaded.Orphans == 0 && loaded.State.TrustPhase == TrustStable &&
		!loaded.State.Quarantined && !loaded.State.AttemptsDisabled &&
		startup.observed && !startup.ownerSessionID.IsZero() &&
		startup.installationID == view.InstallationID && startup.supervisorID == view.SupervisorID &&
		startup.epochSequence == view.EpochSequence && startup.epochDigest == view.EpochDigest &&
		commitP95.observed
	return internalAlphaStoreFacts{
		known: known, cumulativeAttempts: view.TotalCounts.Attempts,
		activeStoreBytes:    loaded.ActiveEncodedBytes,
		archiveSegments:     uint16(len(loaded.Segments)), //nolint:gosec // full verifier caps this at 64.
		startupVerification: startup.duration, startupVerificationObserved: startup.observed,
		durableCommitP95: commitP95.duration, durableCommitP95Observed: commitP95.observed,
		snapshotGeneration: view.SnapshotGeneration, archiveGeneration: view.ArchiveGeneration,
		currentCheckpoint: view.CurrentCheckpoint,
	}
}

func evaluateInternalAlphaFixedStoreFacts(facts internalAlphaStoreFacts) InternalAlphaStoreAdmissionReport {
	report := InternalAlphaStoreAdmissionReport{
		Posture:            InternalAlphaFixedStorePosture,
		CumulativeAttempts: facts.cumulativeAttempts, ActiveStoreBytes: facts.activeStoreBytes,
		ArchiveSegments: facts.archiveSegments, StartupVerification: facts.startupVerification,
		DurableCommitP95: facts.durableCommitP95, SnapshotGeneration: facts.snapshotGeneration,
		ArchiveGeneration: facts.archiveGeneration, CurrentCheckpoint: facts.currentCheckpoint,
	}
	switch {
	case !facts.known || !facts.startupVerificationObserved || !facts.durableCommitP95Observed:
		report.StopReason = InternalAlphaStoreStopUnknownState
	case facts.cumulativeAttempts >= InternalAlphaMaxCumulativeAttempts:
		report.StopReason = InternalAlphaStoreStopCumulativeAttempts
	case facts.activeStoreBytes >= InternalAlphaMaxActiveStoreBytes:
		report.StopReason = InternalAlphaStoreStopActiveStoreBytes
	case facts.archiveSegments >= InternalAlphaMaxArchiveSegments:
		report.StopReason = InternalAlphaStoreStopArchiveSegments
	case facts.startupVerification > InternalAlphaMaxStartupVerification:
		report.StopReason = InternalAlphaStoreStopStartupVerification
	case facts.durableCommitP95 > InternalAlphaMaxDurableCommitP95:
		report.StopReason = InternalAlphaStoreStopDurableCommitP95
	default:
		report.Allowed = true
	}
	return report
}

func internalAlphaOwnerSession(ctx context.Context, owner ArchiveOwner) (lifecyclestate.OwnerSessionID, error) {
	if err := ctx.Err(); err != nil {
		return lifecyclestate.OwnerSessionID{}, err
	}
	if owner == nil {
		return lifecyclestate.OwnerSessionID{}, fmt.Errorf("%w: %w", ErrInternalAlphaAttemptRefused, ErrArchiveOwnerRequired)
	}
	if err := owner.CheckHeld(ctx); err != nil {
		return lifecyclestate.OwnerSessionID{}, fmt.Errorf("%w: %w", ErrInternalAlphaAttemptRefused, ErrStoreOwnerFenced)
	}
	session := owner.OwnerSessionID()
	if session.IsZero() {
		return lifecyclestate.OwnerSessionID{}, fmt.Errorf("%w: %w", ErrInternalAlphaAttemptRefused, ErrStoreOwnerFenced)
	}
	return session, nil
}

func unknownInternalAlphaReport() InternalAlphaStoreAdmissionReport {
	return InternalAlphaStoreAdmissionReport{Posture: InternalAlphaFixedStorePosture, StopReason: InternalAlphaStoreStopUnknownState}
}

func internalAlphaRefusal(reason InternalAlphaStoreStopReason) error {
	return fmt.Errorf("%w: %s", ErrInternalAlphaAttemptRefused, reason)
}
