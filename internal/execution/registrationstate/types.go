// Package registrationstate implements the unwired Phase 2 candidate
// registration, approval-ledger, and attempt-creation boundary. It is backend
// independent and creates no guest.
package registrationstate

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	RegistrationLifetimeSeconds = uint64(300)
	MaxUnexpiredRegistrations   = 256
	MaxRetainedRegistrations    = 4_096
	StorageFormatVersion        = uint64(0)
	RetentionStateRetained      = "retained"
	RegisterPlanPurpose         = "register-plan"
	SubmitApprovalPurpose       = "submit-approval"
	RequestAttemptPurpose       = "request-attempt"
)

type Classification string

const (
	ClassificationMalformed      Classification = "MALFORMED"
	ClassificationUnsupported    Classification = "UNSUPPORTED"
	ClassificationSchema         Classification = "SCHEMA"
	ClassificationPolicy         Classification = "POLICY"
	ClassificationBinding        Classification = "BINDING"
	ClassificationDomain         Classification = "DOMAIN"
	ClassificationAuthentication Classification = "AUTHENTICATION"
	ClassificationStale          Classification = "STALE"
	ClassificationCapacity       Classification = "CAPACITY"
	ClassificationLocalFailure   Classification = "LOCAL_FAILURE"
	ClassificationTrustState     Classification = "TRUST_STATE"
)

type stateError struct {
	classification Classification
	code           string
}

func (e *stateError) Error() string {
	return fmt.Sprintf("%s: %s", e.classification, e.code)
}

func (e *stateError) Classification() Classification { return e.classification }

func classified(classification Classification, code string) error {
	return &stateError{classification: classification, code: code}
}

// ErrorClassification returns the fixed internal conformance classification.
// These values are deliberately not public protocol errors.
func ErrorClassification(err error) (Classification, bool) {
	var registrationError *stateError
	if errors.As(err, &registrationError) {
		return registrationError.classification, true
	}
	if classification, ok := v0candidate.ErrorClassification(err); ok {
		return Classification(classification), true
	}
	return "", false
}

type CallerRole string

const (
	CallerDaemon  CallerRole = "daemon"
	CallerBroker  CallerRole = "broker"
	CallerUpdater CallerRole = "updater"
)

// AuthenticatedCallContext is supplied by trusted local IPC. Authentication
// identifies the peer; the exact role and purpose still have to match.
type AuthenticatedCallContext struct {
	Authenticated bool
	Role          CallerRole
	Purpose       string
}

type TrustPhase string

const (
	TrustStable           TrustPhase = "stable"
	TrustTransitionFenced TrustPhase = "transition-fenced"
	TrustRepairRequired   TrustPhase = "repair-required"
)

const sequenceExhaustedReason = "registration-sequence-exhausted"

// InitialState is the installation-global state required by the fixed first
// store. Records are intentionally not accepted here: this unwired store can
// only create them through Component.RegisterPlan.
type InitialState struct {
	InstallationID           v0candidate.InstallationID
	SupervisorID             v0candidate.SupervisorID
	EpochSequence            v0candidate.UInt53
	EpochDigest              v0candidate.TrustEpochDigest
	TrustPhase               TrustPhase
	TrustReason              string
	Quarantined              bool
	TimeHighWaterUnixSeconds v0candidate.UInt53
	LastRegistrationSequence v0candidate.UInt53
	AttemptsDisabled         bool
}

// PlanRegistration is only the deterministic-CBOR wire response. It never
// contains stored plan bytes or recovery metadata.
type PlanRegistration struct {
	bytes []byte
	view  v0candidate.PlanRegistration
}

func (registration PlanRegistration) Bytes() []byte {
	return bytes.Clone(registration.bytes)
}

func (registration PlanRegistration) View() v0candidate.PlanRegistration {
	return registration.view
}

// StoredRegistration is the exact six-field retained record from ADR-0023's
// Task 2.4 addendum. Byte slices returned to callers are always defensive
// copies of store-owned data.
type StoredRegistration struct {
	WireRegistrationBytes   []byte
	ExactPlanBytes          []byte
	RecomputedPlanDigest    v0candidate.ExecutionPlanDigest
	RegisteredAtUnixSeconds v0candidate.UInt53
	StorageFormatVersion    uint64
	RetentionState          string
}

func cloneStoredRegistration(record StoredRegistration) StoredRegistration {
	record.WireRegistrationBytes = bytes.Clone(record.WireRegistrationBytes)
	record.ExactPlanBytes = bytes.Clone(record.ExactPlanBytes)
	return record
}

type registrationIndex struct {
	RegistrationID       v0candidate.RegistrationID
	RegistrationSequence v0candidate.PositiveUInt53
	InstallationID       v0candidate.InstallationID
	EpochSequence        v0candidate.UInt53
	EpochDigest          v0candidate.TrustEpochDigest
	SupervisorID         v0candidate.SupervisorID
	ExpiresAt            v0candidate.UInt53
}

type registrationEntry struct {
	Index        registrationIndex
	PlanBindings v0candidate.ExecutionPlanRoleBindings
	Record       StoredRegistration
}

// CreatedAttempt is the narrow defensive handoff from the durable Slice B
// authority store to the no-guest lifecycle. It contains only the immutable
// created attempt, the Supervisor-retained registration record, and the
// trusted role bindings required to independently decode those exact plan
// bytes. It contains no approval envelope/reference or backend configuration.
type CreatedAttempt struct {
	Attempt          approvalattempt.ExecutionAttempt
	Approval         ConsumedApprovalBinding
	Registration     StoredRegistration
	PlanRoleBindings v0candidate.ExecutionPlanRoleBindings
}

// ConsumedApprovalBinding is the byte-free cross-link needed to prove that a
// created attempt still belongs to exactly one consumed approval. Exact signed
// approval bytes remain inside the authority store and never enter the
// lifecycle seam.
type ConsumedApprovalBinding struct {
	ApprovalID            approvalattempt.ApprovalID
	ConsumedAttemptID     approvalattempt.AttemptID
	AttemptNonce          approvalattempt.AttemptNonce
	RegistrationID        v0candidate.RegistrationID
	RegistrationSequence  v0candidate.PositiveUInt53
	PlanDigest            v0candidate.ExecutionPlanDigest
	InstallationID        v0candidate.InstallationID
	EpochSequence         v0candidate.UInt53
	EpochDigest           v0candidate.TrustEpochDigest
	SupervisorID          v0candidate.SupervisorID
	Purpose               string
	Audience              string
	PayloadDigest         approvalattempt.ApprovalPayloadDigest
	AuthorizationIdentity approvalattempt.ApprovalKeyAuthorizationIdentity
	State                 approvalattempt.ApprovalState
	StorageFormatVersion  uint64
}

func cloneCreatedAttempt(created CreatedAttempt) CreatedAttempt {
	created.Registration = cloneStoredRegistration(created.Registration)
	created.PlanRoleBindings = clonePlanBindings(created.PlanRoleBindings)
	return created
}

type installationState struct {
	InitialState
	RegistrationSetDigest [32]byte
	Registrations         []registrationEntry
	ApprovalSetDigest     [32]byte
	Approvals             []approvalattempt.ApprovalRecord
	AttemptSetDigest      [32]byte
	Attempts              []approvalattempt.ExecutionAttempt
}

// stateSnapshot is a defensive projection used by conformance and recovery
// tests. The product handoff is RegistrationResolver, not whole-store access.
type stateSnapshot struct {
	InitialState
	StoredRegistrationCount    int
	UnexpiredRegistrationCount int
	RegistrationSetDigest      [32]byte
	Records                    []StoredRegistration
}

type TrustedClock interface {
	ObserveUnixSeconds(context.Context) (uint64, error)
}

type IdentifierSource interface {
	NewRegistrationID(context.Context) (v0candidate.RegistrationID, error)
}

type SystemTrustedClock struct{}

func (SystemTrustedClock) ObserveUnixSeconds(ctx context.Context) (uint64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	now := time.Now().Unix()
	if now < 0 || uint64(now) > v0candidate.MaxSafeInteger {
		return 0, errors.New("system clock is outside UInt53 Unix seconds")
	}
	return uint64(now), nil
}

type CryptoIdentifierSource struct{}

func (CryptoIdentifierSource) NewRegistrationID(
	ctx context.Context,
) (v0candidate.RegistrationID, error) {
	if err := ctx.Err(); err != nil {
		return v0candidate.RegistrationID{}, err
	}
	var identifier v0candidate.RegistrationID
	if _, err := rand.Read(identifier[:]); err != nil {
		return v0candidate.RegistrationID{}, err
	}
	return identifier, nil
}

type Checkpoint string

const (
	CheckpointPlanDecoded       Checkpoint = "plan-decoded"
	CheckpointTimeHighWaterDone Checkpoint = "time-high-water-durable"
)

// CheckpointHook exists only for deterministic concurrency, mutation, and
// fault tests. It receives the already-decoded wrapper, whose accessors return
// defensive copies.
type CheckpointHook func(context.Context, Checkpoint, *v0candidate.DecodedExecutionPlan)

type StateStore interface {
	snapshot(context.Context) (installationState, error)
	persistTimeHighWater(context.Context, v0candidate.UInt53) error
	update(context.Context, func(*installationState) error) error
	commitApproval(context.Context, func(*installationState) error) error
	commitAttempt(context.Context, func(*installationState) error) error
	recoveryFenced() bool
}

// RegistrationResolver is the deliberately narrow handoff for the later fake
// backend. It accepts only a Supervisor-issued registration ID. Replacement
// plan bytes, backend flags, mounts, images, and guest paths are impossible to
// supply through this interface.
type RegistrationResolver interface {
	ResolveUsable(context.Context, v0candidate.RegistrationID) (StoredRegistration, error)
}

var _ RegistrationResolver = (*Component)(nil)

func emptyRegistrationSetDigest() [32]byte {
	return sha256.Sum256([]byte("[]\n"))
}

func emptyApprovalSetDigest() [32]byte { return sha256.Sum256([]byte("[]\n")) }

func emptyAttemptSetDigest() [32]byte { return sha256.Sum256([]byte("[]\n")) }

func clonePlanBindings(
	bindings v0candidate.ExecutionPlanRoleBindings,
) v0candidate.ExecutionPlanRoleBindings {
	bindings.ProfileReviewAttestationDigests = append(
		[]v0candidate.ProfileReviewAttestationDigest(nil),
		bindings.ProfileReviewAttestationDigests...,
	)
	return bindings
}
