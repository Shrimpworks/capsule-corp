package teardownpassive

import "errors"

const (
	// ContractIdentity distinguishes the passive C5b18 model from every
	// existing lifecycle record, driver, and experiment identity.
	ContractIdentity = "capsule.c5b18.teardown-successor-model/v1"
	// RecordVersionV1 is the only record version accepted by this model.
	RecordVersionV1 uint16 = 1
	maxSafeTick            = uint64(1<<53 - 1)
)

var (
	// ErrBinding reports a zero, substituted, or otherwise invalid binding.
	ErrBinding = errors.New("C5B18_BINDING")
	// ErrState reports a transition that the passive state does not permit.
	ErrState = errors.New("C5B18_STATE")
	// ErrPendingWrite reports a second or mismatched storage operation.
	ErrPendingWrite = errors.New("C5B18_PENDING_WRITE")
	// ErrRecovery reports an operation refused after current-lifetime custody was lost.
	ErrRecovery = errors.New("C5B18_RECOVERY_REQUIRED")
	// ErrClock reports an invalid monotonic tick or wall-deadline anchor.
	ErrClock = errors.New("C5B18_CLOCK")
)

type (
	// InstallationID identifies the enrolled local installation.
	InstallationID [16]byte
	// TrustEpochDigest identifies the exact local trust epoch.
	TrustEpochDigest [32]byte
	// SupervisorIdentityDigest identifies the exact expected Supervisor.
	SupervisorIdentityDigest [32]byte
	// AttemptID identifies the consumed attempt.
	AttemptID [16]byte
	// ApprovalID identifies the consumed approval.
	ApprovalID [16]byte
	// RegistrationID identifies the registered plan.
	RegistrationID [16]byte
	// PlanDigest identifies the exact plan bytes.
	PlanDigest [32]byte
	// ProfileDigest identifies the exact runtime-profile bytes.
	ProfileDigest [32]byte
	// RunnerDigest identifies the exact runner candidate bytes.
	RunnerDigest [32]byte
	// OperationID identifies one Supervisor-generated storage operation.
	OperationID [16]byte
	// ProcessIdentity identifies exact current-lifetime process custody.
	ProcessIdentity [32]byte
)

// Bindings are the immutable authority and byte identities retained by the
// prepared teardown obligation. They grant cleanup only.
type Bindings struct {
	InstallationID           InstallationID
	TrustEpochDigest         TrustEpochDigest
	SupervisorIdentityDigest SupervisorIdentityDigest
	AttemptID                AttemptID
	ApprovalID               ApprovalID
	RegistrationID           RegistrationID
	PlanDigest               PlanDigest
	ProfileDigest            ProfileDigest
	RunnerDigest             RunnerDigest
	PreparationOperationID   OperationID
}

func (bindings Bindings) valid() bool {
	identifiers := [][16]byte{
		bindings.InstallationID,
		bindings.AttemptID,
		bindings.ApprovalID,
		bindings.RegistrationID,
		bindings.PreparationOperationID,
	}
	for index, value := range identifiers {
		if zero16(value) {
			return false
		}
		for prior := range index {
			if value == identifiers[prior] {
				return false
			}
		}
	}
	digests := [][32]byte{
		bindings.TrustEpochDigest,
		bindings.SupervisorIdentityDigest,
		bindings.PlanDigest,
		bindings.ProfileDigest,
		bindings.RunnerDigest,
	}
	for index, value := range digests {
		if zero32(value) {
			return false
		}
		for prior := range index {
			if value == digests[prior] {
				return false
			}
		}
	}
	return true
}

func zero16(value [16]byte) bool { return value == [16]byte{} }
func zero32(value [32]byte) bool { return value == [32]byte{} }

// ClockPolicy is the exact same-lifetime C5b18 clock relation in milliseconds.
type ClockPolicy struct {
	WallAfterStartMS  uint64
	ForceAbsenceMS    uint64
	ActionToAbsenceMS uint64
}

var fixedClockPolicy = ClockPolicy{
	WallAfterStartMS:  1000,
	ForceAbsenceMS:    1000,
	ActionToAbsenceMS: 1200,
}

// WriteKind is the closed set of durable operations modeled by C5b18.
type WriteKind string

const (
	// WritePreparation persists the cleanup-only obligation before creation.
	WritePreparation WriteKind = "preparation"
	// WriteRunnerIdentity persists exact runner identity before start.
	WriteRunnerIdentity WriteKind = "runner-identity"
	// WriteTerminalJoin persists completion-last terminal state.
	WriteTerminalJoin WriteKind = "terminal-join"
)

// WriteOutcome is the closed settlement vocabulary for one storage operation.
type WriteOutcome string

const (
	// WriteNone records that no storage outcome has settled.
	WriteNone WriteOutcome = "none"
	// WriteConfirmed records exact confirmed publication.
	WriteConfirmed WriteOutcome = "confirmed"
	// WriteFailed records an exact refused publication.
	WriteFailed WriteOutcome = "failed"
	// WriteIndeterminate records an outcome that cannot be classified safely.
	WriteIndeterminate WriteOutcome = "indeterminate"
)

func validWriteOutcome(outcome WriteOutcome) bool {
	return outcome == WriteConfirmed || outcome == WriteFailed || outcome == WriteIndeterminate
}

// Trigger names the only teardown triggers permitted by the prepared obligation.
type Trigger string

const (
	// TriggerNone records that no stop trigger has latched.
	TriggerNone Trigger = "none"
	// TriggerCancellation records an accepted attempt-bound cancellation.
	TriggerCancellation Trigger = "cancellation"
	// TriggerWallDeadline records the exact wall anchor derived from start.
	TriggerWallDeadline Trigger = "wall-deadline"
	// TriggerFatalFault records a Supervisor-detected fatal fault.
	TriggerFatalFault Trigger = "fatal-fault"
)

func validTrigger(trigger Trigger) bool {
	return trigger == TriggerCancellation || trigger == TriggerWallDeadline || trigger == TriggerFatalFault
}

// Custody records only current-Supervisor-lifetime exact process custody.
type Custody string

const (
	// CustodyNone records that the current Supervisor has no exact live custody.
	CustodyNone Custody = "none"
	// CustodyExact records exact live custody in the current Supervisor lifetime.
	CustodyExact Custody = "exact-current-lifetime"
)

// TimingDisposition classifies only retained same-lifetime clock observations.
type TimingDisposition string

const (
	// TimingUnknown records that the model lacks a complete same-lifetime relation.
	TimingUnknown TimingDisposition = "unknown"
	// TimingSatisfied records that both applicable bounds passed.
	TimingSatisfied TimingDisposition = "satisfied"
	// TimingViolated records that at least one applicable bound failed.
	TimingViolated TimingDisposition = "violated"
)

// PendingWrite binds the sole outstanding operation. Its fields form the token
// required for exact settlement.
type PendingWrite struct {
	Active      bool
	OperationID OperationID
	Generation  uint64
	Kind        WriteKind
}

// StopSnapshot retains the monotonic latch and earliest accepted action anchor.
type StopSnapshot struct {
	Latched    bool
	Trigger    Trigger
	ActionTick uint64
}

// SignalSnapshot retains only the first current-lifetime request and uncertainty.
type SignalSnapshot struct {
	Requested bool
	Uncertain bool
	Attempts  uint8
	ForceTick uint64
}

// AbsenceSnapshot is an authoritative observation for exact current custody.
type AbsenceSnapshot struct {
	Observed bool
	Tick     uint64
}

// StartSnapshot retains the sole start-token attempt and its original anchor.
type StartSnapshot struct {
	Attempted bool
	Tick      uint64
}

// Snapshot is a defensive public projection used by the independent decision oracle.
type Snapshot struct {
	Contract                string
	RecordVersion           uint16
	PreparationStarted      bool
	Prepared                bool
	CreationConsumed        bool
	RunnerWriteStarted      bool
	RunnerIdentityConfirmed bool
	TerminalWriteStarted    bool
	TerminalConfirmed       bool
	LastSettledGeneration   uint64
	LastWriteOutcome        WriteOutcome
	Pending                 PendingWrite
	Custody                 Custody
	ProcessIdentity         ProcessIdentity
	Stop                    StopSnapshot
	Signal                  SignalSnapshot
	Absence                 AbsenceSnapshot
	Start                   StartSnapshot
	Timing                  TimingDisposition
	RecoveryRequired        bool
	OutputReleased          bool
	CapacityReleased        bool
}

// Decision is the closed authorization projection derived from one snapshot.
type Decision struct {
	MayCreate                bool
	MayStart                 bool
	MaySignal                bool
	MayPublishTerminal       bool
	MayRelease               bool
	StopIndependentOfStorage bool
	RecoveryRequired         bool
}
