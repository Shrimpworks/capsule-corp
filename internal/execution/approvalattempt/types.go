// Package approvalattempt contains passive, unwired approval-consumption and
// attempt-creation candidate contracts. It does not persist authority, consume
// an approval, create an attempt, call a lifecycle, or authorize a backend.
package approvalattempt

import (
	"bytes"
	"errors"
	"fmt"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	ApprovalGrantObjectType = "capsule.approval-grant"
	ApprovalGrantPurpose    = "capsule.plan.approve"
	ApprovalGrantAudience   = "capsule.execution-supervisor"
)

// Classification is the complete fixed internal vocabulary from ADR-0024.
// It is not a public protocol error or an execution result.
type Classification string

const (
	ClassificationAuthentication   Classification = "AUTHENTICATION"
	ClassificationMalformed        Classification = "MALFORMED"
	ClassificationUnsupported      Classification = "UNSUPPORTED"
	ClassificationSchema           Classification = "SCHEMA"
	ClassificationDomain           Classification = "DOMAIN"
	ClassificationBinding          Classification = "BINDING"
	ClassificationStale            Classification = "STALE"
	ClassificationReplay           Classification = "REPLAY"
	ClassificationCapacity         Classification = "CAPACITY"
	ClassificationTrustState       Classification = "TRUST_STATE"
	ClassificationLocalFailure     Classification = "LOCAL_FAILURE"
	ClassificationRecoveryRequired Classification = "RECOVERY_REQUIRED"
)

var fixedClassifications = [...]Classification{
	ClassificationAuthentication,
	ClassificationMalformed,
	ClassificationUnsupported,
	ClassificationSchema,
	ClassificationDomain,
	ClassificationBinding,
	ClassificationStale,
	ClassificationReplay,
	ClassificationCapacity,
	ClassificationTrustState,
	ClassificationLocalFailure,
	ClassificationRecoveryRequired,
}

func Classifications() []Classification {
	return append([]Classification(nil), fixedClassifications[:]...)
}

type contractError struct {
	classification Classification
	code           string
}

func (e *contractError) Error() string {
	return fmt.Sprintf("%s: %s", e.classification, e.code)
}

func (e *contractError) Classification() Classification { return e.classification }

func classified(classification Classification, code string) error {
	return &contractError{classification: classification, code: code}
}

// ErrorClassification extracts a fixed content-free classification.
func ErrorClassification(err error) (Classification, bool) {
	var candidateError *contractError
	if !errors.As(err, &candidateError) {
		return "", false
	}
	return candidateError.classification, true
}

type IdentifierDomain string

const (
	DomainApprovalID   IdentifierDomain = "approval"
	DomainAttemptID    IdentifierDomain = "attempt"
	DomainAttemptNonce IdentifierDomain = "attempt-nonce"
)

// DomainIdentifier carries the separately trusted semantic role used by the
// conformance boundary. The bytes themselves do not encode their nominal role.
type DomainIdentifier struct {
	domain IdentifierDomain
	bytes  [16]byte
}

func NewDomainIdentifier(domain IdentifierDomain, value []byte) (DomainIdentifier, error) {
	if domain != DomainApprovalID && domain != DomainAttemptID && domain != DomainAttemptNonce {
		return DomainIdentifier{}, classified(ClassificationDomain, "unknown-identifier-domain")
	}
	if len(value) != 16 {
		return DomainIdentifier{}, classified(ClassificationSchema, "identifier-length")
	}
	var result DomainIdentifier
	result.domain = domain
	copy(result.bytes[:], value)
	if result.bytes == ([16]byte{}) {
		return DomainIdentifier{}, classified(ClassificationSchema, "identifier-zero")
	}
	return result, nil
}

func (identifier DomainIdentifier) Domain() IdentifierDomain { return identifier.domain }

func (identifier DomainIdentifier) Bytes() []byte { return bytes.Clone(identifier.bytes[:]) }

type ApprovalID [16]byte
type AttemptID [16]byte
type AttemptNonce [16]byte

func NewApprovalID(value DomainIdentifier) (ApprovalID, error) {
	if value.domain != DomainApprovalID {
		return ApprovalID{}, classified(ClassificationDomain, "approval-id-domain")
	}
	return ApprovalID(value.bytes), nil
}

func NewAttemptID(value DomainIdentifier) (AttemptID, error) {
	if value.domain != DomainAttemptID {
		return AttemptID{}, classified(ClassificationDomain, "attempt-id-domain")
	}
	return AttemptID(value.bytes), nil
}

func NewAttemptNonce(value DomainIdentifier) (AttemptNonce, error) {
	if value.domain != DomainAttemptNonce {
		return AttemptNonce{}, classified(ClassificationDomain, "attempt-nonce-domain")
	}
	return AttemptNonce(value.bytes), nil
}

type ApprovalReference struct {
	approvalID ApprovalID
}

func NewApprovalReference(approvalID ApprovalID) (ApprovalReference, error) {
	if approvalID == (ApprovalID{}) {
		return ApprovalReference{}, classified(ClassificationSchema, "approval-reference-zero")
	}
	return ApprovalReference{approvalID: approvalID}, nil
}

func (reference ApprovalReference) ApprovalID() ApprovalID { return reference.approvalID }

type AttemptReference struct {
	attemptID AttemptID
}

func NewAttemptReference(attemptID AttemptID) (AttemptReference, error) {
	if attemptID == (AttemptID{}) {
		return AttemptReference{}, classified(ClassificationSchema, "attempt-reference-zero")
	}
	return AttemptReference{attemptID: attemptID}, nil
}

func (reference AttemptReference) AttemptID() AttemptID { return reference.attemptID }

// ApprovalKeyAuthorizationIdentity is a fixture-only resolved identity. The
// production KeyAuthorization object remains blocked on ADR-0019.
type ApprovalKeyAuthorizationIdentity [32]byte

func NewApprovalKeyAuthorizationIdentity(value []byte) (ApprovalKeyAuthorizationIdentity, error) {
	if len(value) != 32 {
		return ApprovalKeyAuthorizationIdentity{}, classified(ClassificationSchema, "authorization-identity-length")
	}
	var identity ApprovalKeyAuthorizationIdentity
	copy(identity[:], value)
	if identity == (ApprovalKeyAuthorizationIdentity{}) {
		return ApprovalKeyAuthorizationIdentity{}, classified(ClassificationSchema, "authorization-identity-zero")
	}
	return identity, nil
}

// ApprovalGrant is the decoded passive view of the candidate signed payload.
type ApprovalGrant struct {
	ObjectType     string
	ObjectVersion  v0candidate.ObjectVersion
	InstallationID v0candidate.InstallationID
	EpochDigest    v0candidate.TrustEpochDigest
	RegistrationID v0candidate.RegistrationID
	PlanDigest     v0candidate.ExecutionPlanDigest
	SupervisorID   v0candidate.SupervisorID
	AttemptNonce   AttemptNonce
	Purpose        string
	Audience       string
	IssuedAt       v0candidate.UInt53
	ExpiresAt      v0candidate.UInt53
}

type ApprovalGrantRoleBindings struct {
	InstallationID        v0candidate.InstallationID
	EpochSequence         v0candidate.UInt53
	EpochDigest           v0candidate.TrustEpochDigest
	RegistrationID        v0candidate.RegistrationID
	PlanDigest            v0candidate.ExecutionPlanDigest
	SupervisorID          v0candidate.SupervisorID
	AttemptNonce          AttemptNonce
	ProtectedKeyID        []byte
	AuthorizationIdentity ApprovalKeyAuthorizationIdentity
}

func cloneRoleBindings(bindings ApprovalGrantRoleBindings) ApprovalGrantRoleBindings {
	bindings.ProtectedKeyID = bytes.Clone(bindings.ProtectedKeyID)
	return bindings
}

type ApprovalState string

const (
	ApprovalUsable      ApprovalState = "usable"
	ApprovalConsumed    ApprovalState = "consumed"
	ApprovalInvalidated ApprovalState = "invalidated"
)

type AttemptState string

const AttemptCreated AttemptState = "created"
