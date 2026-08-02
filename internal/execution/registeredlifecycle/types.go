package registeredlifecycle

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/execution/registrationstate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const MaxLifecycleRecords = 256

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

type State string

const (
	StatePreparing    State = "preparing"
	StateCreateIntent State = "create-intent"
	StateCreated      State = "created"
	StateStarting     State = "starting"
	StateObserving    State = "observing"
	StateStopping     State = "stopping"
	StateDestroying   State = "destroying"
	StateDestroyed    State = "destroyed"
	StateUnresolved   State = "unresolved"
)

// Snapshot is the complete caller-visible lifecycle state. It contains no
// plan bytes, backend handle, content, guest data, success result, or free-form
// failure text.
type Snapshot struct {
	AttemptID             approvalattempt.AttemptID
	ApprovalID            approvalattempt.ApprovalID
	AttemptNonce          approvalattempt.AttemptNonce
	RegistrationID        v0candidate.RegistrationID
	RegistrationSequence  v0candidate.PositiveUInt53
	PlanDigest            v0candidate.ExecutionPlanDigest
	InstallationID        v0candidate.InstallationID
	EpochSequence         v0candidate.UInt53
	EpochDigest           v0candidate.TrustEpochDigest
	SupervisorID          v0candidate.SupervisorID
	ApprovalPurpose       string
	ApprovalAudience      string
	ApprovalPayloadDigest approvalattempt.ApprovalPayloadDigest
	AuthorizationIdentity approvalattempt.ApprovalKeyAuthorizationIdentity
	CreatedAt             v0candidate.UInt53
	State                 State
	CleanupRequired       bool
	Failure               Classification
	FailureAt             Operation
	TransitionCount       uint64
}

// AttemptResolver is the only authority seam used before a first lifecycle
// effect. It resolves an already committed created attempt by its distinct
// AttemptID and returns copied immutable bindings plus retained exact plan
// bytes. Callers cannot supply approval or plan bytes through Drive/Recover.
type AttemptResolver interface {
	ResolveCreated(context.Context, approvalattempt.AttemptID) (registrationstate.CreatedAttempt, error)
}

// CreatedAttemptEnumerator is separate from AttemptResolver so startup
// enumeration does not widen the ID-only resolution seam.
type CreatedAttemptEnumerator interface {
	CreatedAttempts(context.Context) ([]approvalattempt.AttemptReference, error)
}

var (
	_ AttemptResolver          = (*registrationstate.ApprovalAttemptComponent)(nil)
	_ CreatedAttemptEnumerator = (*registrationstate.ApprovalAttemptComponent)(nil)
)

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
	Attempts        AttemptResolver
	CreatedAttempts CreatedAttemptEnumerator
	Store           *MemoryStore
	Backend         *FakeBackend
	Checkpoint      CheckpointHook
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

func createdAttemptsEqual(left, right registrationstate.CreatedAttempt) bool {
	return left.Attempt == right.Attempt &&
		left.Approval == right.Approval &&
		left.Registration.RecomputedPlanDigest == right.Registration.RecomputedPlanDigest &&
		left.Registration.RegisteredAtUnixSeconds == right.Registration.RegisteredAtUnixSeconds &&
		left.Registration.StorageFormatVersion == right.Registration.StorageFormatVersion &&
		left.Registration.RetentionState == right.Registration.RetentionState &&
		bytes.Equal(left.Registration.WireRegistrationBytes, right.Registration.WireRegistrationBytes) &&
		bytes.Equal(left.Registration.ExactPlanBytes, right.Registration.ExactPlanBytes) &&
		reflect.DeepEqual(left.PlanRoleBindings, right.PlanRoleBindings)
}
