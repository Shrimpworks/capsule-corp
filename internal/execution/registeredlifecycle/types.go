package registeredlifecycle

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/execution/registrationstate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// Classification is a fixed, content-free internal development oracle. It is
// not a public protocol error or an execution result.
type Classification string

const (
	ClassificationMalformed         Classification = "MALFORMED"
	ClassificationUnsupported       Classification = "UNSUPPORTED"
	ClassificationSchema            Classification = "SCHEMA"
	ClassificationDomain            Classification = "DOMAIN"
	ClassificationBinding           Classification = "BINDING"
	ClassificationStale             Classification = "STALE"
	ClassificationCapacity          Classification = "CAPACITY"
	ClassificationLocalFailure      Classification = "LOCAL_FAILURE"
	ClassificationLifecycleFailure  Classification = "LIFECYCLE_FAILURE"
	ClassificationCleanupUnresolved Classification = "CLEANUP_UNRESOLVED"
	ClassificationTrustState        Classification = "TRUST_STATE"
	ClassificationRecoveryRequired  Classification = "RECOVERY_REQUIRED"
)

type lifecycleError struct {
	classification Classification
	code           string
}

func (e *lifecycleError) Error() string {
	return fmt.Sprintf("%s: %s", e.classification, e.code)
}

func classified(classification Classification, code string) error {
	return &lifecycleError{classification: classification, code: code}
}

// ErrorClassification extracts the fixed internal classification without
// exposing plan bytes, caller text, or backend-controlled strings.
func ErrorClassification(err error) (Classification, bool) {
	var lifecycleErr *lifecycleError
	if errors.As(err, &lifecycleErr) {
		return lifecycleErr.classification, true
	}
	return "", false
}

type Operation string

const (
	OperationPrepare   Operation = "prepare"
	OperationCreate    Operation = "create"
	OperationStart     Operation = "start"
	OperationObserve   Operation = "observe"
	OperationStop      Operation = "stop"
	OperationDestroy   Operation = "destroy"
	OperationReconcile Operation = "reconcile"
)

type FaultMoment string

const (
	FaultBeforeEffect FaultMoment = "before-effect"
	FaultAfterEffect  FaultMoment = "after-effect"
)

type State = lifecyclestate.LifecycleState

const (
	StatePreparing    State = lifecyclestate.StatePreparePending
	StateCreateIntent State = lifecyclestate.StateCreateIntent
	StateCreated      State = lifecyclestate.StateCreated
	StateStarting     State = lifecyclestate.StateStartIntent
	StateObserving    State = lifecyclestate.StateObserveIntent
	StateStopping     State = lifecyclestate.StateStopIntent
	StateDestroying   State = lifecyclestate.StateDestroyIntent
	StateDestroyed    State = lifecyclestate.StateDestroyed
	StateUnresolved   State = lifecyclestate.StateUnresolved
	StateQuarantined  State = lifecyclestate.StateQuarantined
)

// Snapshot is the complete caller-visible lifecycle state. It contains no
// plan bytes, backend handle, content, guest data, success result, or free-form
// failure text.
type Snapshot struct {
	AttemptID                  approvalattempt.AttemptID
	ApprovalID                 approvalattempt.ApprovalID
	AttemptNonce               approvalattempt.AttemptNonce
	RegistrationID             v0candidate.RegistrationID
	RegistrationSequence       v0candidate.PositiveUInt53
	PlanDigest                 v0candidate.ExecutionPlanDigest
	InstallationID             v0candidate.InstallationID
	EpochSequence              v0candidate.UInt53
	EpochDigest                v0candidate.TrustEpochDigest
	SupervisorID               v0candidate.SupervisorID
	ApprovalPurpose            string
	ApprovalAudience           string
	ApprovalPayloadDigest      approvalattempt.ApprovalPayloadDigest
	AuthorizationIdentity      approvalattempt.ApprovalKeyAuthorizationIdentity
	CreatedAt                  v0candidate.UInt53
	State                      State
	CleanupRequired            bool
	Failure                    Classification
	FailureAt                  Operation
	TransitionCount            uint64
	RecordVersion              lifecyclestate.RecordVersion
	OperationSequence          v0candidate.PositiveUInt53
	EffectID                   lifecyclestate.EffectID
	EffectStatus               lifecyclestate.EffectStatus
	InstanceDigest             lifecyclestate.BackendInstanceDigest
	LastReconciliation         lifecyclestate.ReconciliationStatus
	AutomaticRecoveryCount     v0candidate.UInt53
	NextRecoveryAt             lifecyclestate.OptionalUnixSeconds
	RecoveryFence              lifecyclestate.RecoveryFenceReason
	AutomaticRecoveryExhausted bool
}

// AttemptResolver is the only authority seam used before a first lifecycle
// effect. It resolves an already committed created attempt by its distinct
// AttemptID and returns copied immutable bindings plus retained exact plan
// bytes. Callers cannot supply approval or plan bytes through Drive/Recover.
type AttemptResolver interface {
	ResolveCreated(context.Context, approvalattempt.AttemptID) (registrationstate.CreatedAttempt, error)
}

var _ AttemptResolver = (*registrationstate.ApprovalAttemptComponent)(nil)

type Checkpoint string

const (
	CheckpointAfterPrepareEffect Checkpoint = "after-prepare-effect"
	CheckpointAfterCreateEffect  Checkpoint = "after-create-effect"
	CheckpointAfterStartEffect   Checkpoint = "after-start-effect"
	CheckpointAfterObserveEffect Checkpoint = "after-observe-effect"
	CheckpointAfterStopEffect    Checkpoint = "after-stop-effect"
	CheckpointAfterDestroyEffect Checkpoint = "after-destroy-effect"
)

// CheckpointHook simulates process interruption after a fake side effect. A
// returned error is converted to a fixed LOCAL_FAILURE and deliberately skips
// in-process cleanup so that a new Component can exercise recovery.
type CheckpointHook func(context.Context, Checkpoint, approvalattempt.AttemptID) error

type Options struct {
	Attempts    AttemptResolver
	Store       registrationstate.DurableLifecycleStore
	Backend     *FakeBackend
	Coordinator *Coordinator
	Clock       registrationstate.TrustedClock
	Checkpoint  CheckpointHook
}

func cloneRoleBindings(
	bindings v0candidate.ExecutionPlanRoleBindings,
) v0candidate.ExecutionPlanRoleBindings {
	bindings.ProfileReviewAttestationDigests = append(
		[]v0candidate.ProfileReviewAttestationDigest(nil),
		bindings.ProfileReviewAttestationDigests...,
	)
	return bindings
}

func mapAttemptResolutionFailure(err error) error {
	classification, ok := approvalattempt.ErrorClassification(err)
	if !ok {
		return classified(ClassificationLocalFailure, "attempt-resolution-failed")
	}
	switch classification {
	case approvalattempt.ClassificationMalformed:
		return classified(ClassificationMalformed, "attempt-resolution-rejected")
	case approvalattempt.ClassificationUnsupported:
		return classified(ClassificationUnsupported, "attempt-resolution-rejected")
	case approvalattempt.ClassificationSchema:
		return classified(ClassificationSchema, "attempt-resolution-rejected")
	case approvalattempt.ClassificationDomain:
		return classified(ClassificationDomain, "attempt-resolution-rejected")
	case approvalattempt.ClassificationBinding:
		return classified(ClassificationBinding, "attempt-resolution-rejected")
	case approvalattempt.ClassificationStale:
		return classified(ClassificationStale, "attempt-resolution-rejected")
	case approvalattempt.ClassificationCapacity:
		return classified(ClassificationCapacity, "attempt-resolution-rejected")
	case approvalattempt.ClassificationTrustState:
		return classified(ClassificationTrustState, "attempt-resolution-rejected")
	case approvalattempt.ClassificationRecoveryRequired:
		return classified(ClassificationRecoveryRequired, "attempt-resolution-rejected")
	default:
		return classified(ClassificationLocalFailure, "attempt-resolution-failed")
	}
}

func cloneCreatedAttempt(created registrationstate.CreatedAttempt) registrationstate.CreatedAttempt {
	created.Registration.WireRegistrationBytes = bytes.Clone(created.Registration.WireRegistrationBytes)
	created.Registration.ExactPlanBytes = bytes.Clone(created.Registration.ExactPlanBytes)
	created.PlanRoleBindings = cloneRoleBindings(created.PlanRoleBindings)
	return created
}
