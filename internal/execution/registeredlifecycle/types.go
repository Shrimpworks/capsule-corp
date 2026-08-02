package registeredlifecycle

import (
	"context"
	"errors"
	"fmt"

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
	RegistrationID  v0candidate.RegistrationID
	PlanDigest      v0candidate.ExecutionPlanDigest
	State           State
	CleanupRequired bool
	Failure         Classification
	FailureAt       Operation
	TransitionCount uint64
}

// PlanRoleBindingResolver supplies role-resolved trusted identities from a
// Supervisor-owned source. These bindings are construction-time authority;
// callers cannot supply them to Execute or Recover.
type PlanRoleBindingResolver interface {
	ResolveExecutionPlanRoleBindings(
		context.Context,
		v0candidate.RegistrationID,
	) (v0candidate.ExecutionPlanRoleBindings, error)
}

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
type CheckpointHook func(context.Context, Checkpoint, v0candidate.RegistrationID) error

type Options struct {
	Registrations registrationstate.RegistrationResolver
	RoleBindings  PlanRoleBindingResolver
	Store         *MemoryStore
	Backend       *FakeBackend
	Checkpoint    CheckpointHook
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

func mapRegistrationFailure(err error) error {
	classification, ok := registrationstate.ErrorClassification(err)
	if !ok {
		return classified(ClassificationLocalFailure, "registration-resolution-failed")
	}
	switch classification {
	case registrationstate.ClassificationMalformed:
		return classified(ClassificationMalformed, "registration-resolution-rejected")
	case registrationstate.ClassificationUnsupported:
		return classified(ClassificationUnsupported, "registration-resolution-rejected")
	case registrationstate.ClassificationSchema:
		return classified(ClassificationSchema, "registration-resolution-rejected")
	case registrationstate.ClassificationDomain:
		return classified(ClassificationDomain, "registration-resolution-rejected")
	case registrationstate.ClassificationBinding:
		return classified(ClassificationBinding, "registration-resolution-rejected")
	case registrationstate.ClassificationStale:
		return classified(ClassificationStale, "registration-resolution-rejected")
	case registrationstate.ClassificationCapacity:
		return classified(ClassificationCapacity, "registration-resolution-rejected")
	case registrationstate.ClassificationTrustState:
		return classified(ClassificationTrustState, "registration-resolution-rejected")
	default:
		return classified(ClassificationLocalFailure, "registration-resolution-failed")
	}
}
