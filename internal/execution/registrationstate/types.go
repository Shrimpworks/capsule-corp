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
	"capsule.local/capsule/internal/execution/lifecyclestate"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	// RegistrationLifetimeSeconds is the trusted maximum lifetime admitted
	// between the Supervisor's durable effective time and a plan's expiry.
	// The remaining constants are fixed local store capacities, format labels,
	// and authenticated method purposes; callers may not override them.
	RegistrationLifetimeSeconds = uint64(300)
	// MaxUnexpiredRegistrations caps concurrently usable registration authority.
	MaxUnexpiredRegistrations = 256
	// MaxRetainedRegistrations caps retained v0 records without eviction.
	MaxRetainedRegistrations = 4_096
	// StorageFormatVersion binds the retained registration-record shape.
	StorageFormatVersion = uint64(0)
	// RetentionStateRetained is the only admitted v0 record-retention state.
	RetentionStateRetained = "retained"
	// RegisterPlanPurpose authenticates daemon-to-Supervisor registration calls.
	RegisterPlanPurpose = "register-plan"
	// SubmitApprovalPurpose authenticates Broker-to-Supervisor submission calls.
	SubmitApprovalPurpose = "submit-approval"
	// RequestAttemptPurpose authenticates daemon-to-Supervisor attempt calls.
	RequestAttemptPurpose = "request-attempt"
)

// Classification is the package's closed, content-free conformance failure
// vocabulary. It identifies the boundary that refused an operation without
// becoming a public protocol error or carrying caller-controlled prose.
type Classification string

const (
	// ClassificationMalformed through ClassificationRecoveryRequired separate
	// request defects, authority-binding failures, capacity, local failure, and
	// indeterminate recovery. Callers must handle unknown errors fail closed and
	// must not treat any classification as permission to retry an effect.
	ClassificationMalformed Classification = "MALFORMED"
	// ClassificationUnsupported rejects a recognized boundary's unavailable feature/version.
	ClassificationUnsupported Classification = "UNSUPPORTED"
	// ClassificationSchema rejects a known object with an invalid closed shape.
	ClassificationSchema Classification = "SCHEMA"
	// ClassificationPolicy rejects valid input under trusted fixed policy.
	ClassificationPolicy Classification = "POLICY"
	// ClassificationBinding rejects disagreement with retained authority facts.
	ClassificationBinding Classification = "BINDING"
	// ClassificationDomain rejects a valid-width value used in the wrong role.
	ClassificationDomain Classification = "DOMAIN"
	// ClassificationAuthentication rejects an unauthenticated or wrong-role/purpose caller.
	ClassificationAuthentication Classification = "AUTHENTICATION"
	// ClassificationStale rejects expired, fenced, superseded, or sequence-invalid state.
	ClassificationStale Classification = "STALE"
	// ClassificationCapacity rejects a fixed bound without eviction or authority change.
	ClassificationCapacity Classification = "CAPACITY"
	// ClassificationLocalFailure reports a confirmed local prerequisite failure.
	ClassificationLocalFailure Classification = "LOCAL_FAILURE"
	// ClassificationTrustState rejects quarantined or repair-required installation state.
	ClassificationTrustState Classification = "TRUST_STATE"
	// ClassificationLifecycleFailure records an unresolved lifecycle operation failure.
	ClassificationLifecycleFailure Classification = "LIFECYCLE_FAILURE"
	// ClassificationCleanupUnresolved records cleanup whose authoritative result is unknown.
	ClassificationCleanupUnresolved Classification = "CLEANUP_UNRESOLVED"
	// ClassificationRecoveryRequired fences an indeterminate or invalid durable state.
	ClassificationRecoveryRequired Classification = "RECOVERY_REQUIRED"
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
	if classification, ok := lifecyclestate.ErrorClassification(err); ok {
		return Classification(classification), true
	}
	return "", false
}

// CallerRole names the authenticated local role asserted by the future IPC
// boundary. A role is an identity classification only; each operation also
// requires its exact method purpose and repeats state/binding validation.
type CallerRole string

const (
	// CallerDaemon identifies the planning daemon role for exact registration
	// and attempt-request purposes. The value grants no authority when supplied
	// through an unauthenticated or wrong-purpose call context.
	CallerDaemon CallerRole = "daemon"
	// CallerBroker identifies the trusted approval/content Broker role.
	CallerBroker CallerRole = "broker"
	// CallerUpdater identifies the external-trust/update role; it has no registration authority here.
	CallerUpdater CallerRole = "updater"
)

// AuthenticatedCallContext is supplied by trusted local IPC. Authentication
// identifies the peer; the exact role and purpose still have to match.
type AuthenticatedCallContext struct {
	Authenticated bool
	Role          CallerRole
	Purpose       string
}

// TrustPhase records whether the local installation state may admit authority,
// is fenced during a transition, or requires repair. Only TrustStable admits
// new work; package callers cannot clear or normalize another phase.
type TrustPhase string

const (
	// TrustStable is the only phase that may admit new authority after every
	// other binding and policy check passes; it is not proof of product readiness.
	TrustStable TrustPhase = "stable"
	// TrustTransitionFenced blocks admission during an incomplete trust transition.
	TrustTransitionFenced TrustPhase = "transition-fenced"
	// TrustRepairRequired blocks admission until an explicit evidence-preserving repair.
	TrustRepairRequired TrustPhase = "repair-required"
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

// Bytes returns a defensive copy of the exact deterministic-CBOR registration
// response. The bytes are a typed local response, not signed portable authority
// and never contain or replace the Supervisor-retained plan bytes.
func (registration PlanRegistration) Bytes() []byte {
	return bytes.Clone(registration.bytes)
}

// View returns the decoded, typed projection of the Supervisor-issued wire
// response. The projection is informational until its bindings are resolved
// against Supervisor-retained state and does not contain the plan bytes.
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

// TrustedClock supplies a nonnegative UInt53 Unix-second observation to the
// Supervisor-owned durable high-water rule. Implementations must report
// observation failure rather than substitute, clamp, or move time backward.
type TrustedClock interface {
	ObserveUnixSeconds(context.Context) (uint64, error)
}

// IdentifierSource supplies fresh Supervisor-owned registration identities.
// The store rejects zero and retained collisions; callers never provide a
// registration ID through this port.
type IdentifierSource interface {
	NewRegistrationID(context.Context) (v0candidate.RegistrationID, error)
}

// SystemTrustedClock is the local wall-clock implementation used by unwired
// conformance components. Durable monotonicity comes from the store high-water,
// not from this stateless source or from a platform-attestation claim.
type SystemTrustedClock struct{}

// ObserveUnixSeconds returns the current whole Unix second after cancellation
// and UInt53-range checks. It performs no persistence; the caller must commit
// the observation before relying on it for expiry or authority decisions.
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

// CryptoIdentifierSource draws opaque registration IDs from the operating
// system cryptographic random source. Uniqueness and durable non-reuse remain
// store obligations, including across all visible retained history.
type CryptoIdentifierSource struct{}

// NewRegistrationID returns one caller-independent candidate identity or an
// error. The value does not grant registration authority until committed with
// the complete Supervisor record in one store transaction.
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

// Checkpoint identifies test-only observation edges in registration handling.
// It is not a durable lifecycle checkpoint, execution permission, or evidence
// that a registration transaction committed.
type Checkpoint string

const (
	// CheckpointPlanDecoded and CheckpointTimeHighWaterDone bracket validation
	// and durable-time persistence for deterministic fault/concurrency tests.
	// Hooks must not be used as product callbacks or authority sources.
	CheckpointPlanDecoded Checkpoint = "plan-decoded"
	// CheckpointTimeHighWaterDone follows confirmed durable time persistence.
	CheckpointTimeHighWaterDone Checkpoint = "time-high-water-durable"
)

// CheckpointHook exists only for deterministic concurrency, mutation, and
// fault tests. It receives the already-decoded wrapper, whose accessors return
// defensive copies.
type CheckpointHook func(context.Context, Checkpoint, *v0candidate.DecodedExecutionPlan)

// StateStore is the package-private Supervisor transaction seam shared by
// registration and approval/attempt components. Implementations must return
// defensive state, preserve atomic consume/create semantics, and fence every
// indeterminate outcome before later mutation; it is not a public store API.
type StateStore interface {
	snapshot(context.Context) (installationState, error)
	persistTimeHighWater(context.Context, v0candidate.UInt53) error
	update(context.Context, func(*installationState) error) error
	commitApproval(context.Context, func(*installationState) error) error
	commitAttempt(context.Context, func(*installationState) error) error
	recoveryFenced() bool
}

// DurableLifecycleStore is the unwired ADR-0025 E3 transaction boundary. It
// accepts only Supervisor-issued AttemptID authority and closed lifecycle
// values; no approval bytes, plan bytes, backend flags, paths, images, mounts,
// or guest configuration can enter through this interface.
type DurableLifecycleStore interface {
	EnsureLifecycle(context.Context, approvalattempt.AttemptID, lifecyclestate.BackendBinding) (lifecyclestate.Record, bool, error)
	ReadLifecycle(context.Context, approvalattempt.AttemptID) (lifecyclestate.Record, error)
	AdvanceLifecycleTime(context.Context, v0candidate.UInt53) error
	BeginEffect(context.Context, approvalattempt.AttemptID, lifecyclestate.RecordVersion, lifecyclestate.Operation) (lifecyclestate.EffectPermit, error)
	ConfirmEffect(context.Context, lifecyclestate.EffectPermit, lifecyclestate.EffectResult) (lifecyclestate.Record, error)
	RecordIndeterminate(context.Context, lifecyclestate.EffectPermit, Classification) (lifecyclestate.Record, error)
	BeginReconciliation(context.Context, approvalattempt.AttemptID, lifecyclestate.RecordVersion) (lifecyclestate.Record, error)
	CompleteReconciliation(context.Context, approvalattempt.AttemptID, lifecyclestate.RecordVersion, lifecyclestate.ReconcileResult) (lifecyclestate.Record, error)
	RecordReconciliation(context.Context, approvalattempt.AttemptID, lifecyclestate.RecordVersion, lifecyclestate.ReconcileResult) (lifecyclestate.Record, error)
	RecoveryAttemptIDs(context.Context) ([]approvalattempt.AttemptID, error)
	OwnerSessionID() lifecyclestate.OwnerSessionID
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
