package approvalattempt

import (
	"bytes"
	"context"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	ApprovalEnvelopeRawMaxBytes         = 512
	ApprovalPayloadRawMaxBytes          = 256
	ApprovalProtectedRawMaxBytes        = 128
	ApprovalEnvelopeCalculatedMaxBytes  = 431
	ApprovalPayloadCalculatedMaxBytes   = 242
	ApprovalProtectedCalculatedMaxBytes = 116
)

// Verifier is the bounded role-bound verification port. FixtureVerifier and
// the unwired public-key-only ProductionShapedVerifier both satisfy it.
type Verifier interface {
	Verify(context.Context, []byte, ApprovalGrantRoleBindings) (*VerifiedApproval, error)
}

// CandidateVerifier is the existing Slice B fixture injection seam. It
// resolves retained test metadata before the store applies registration and
// active-state bindings; ProductionShapedVerifier is intentionally not wired here.
type CandidateVerifier interface {
	VerifyCandidate(context.Context, []byte) (*VerifiedApproval, error)
}

// FixtureVector is retained test metadata for one exact envelope. It is not a
// general COSE representation and must never be populated from caller claims.
type FixtureVector struct {
	EnvelopeBytes         []byte
	PayloadBytes          []byte
	ProtectedHeaderBytes  []byte
	ProtectedKeyID        []byte
	View                  ApprovalGrant
	ResolvedEpochSequence v0candidate.UInt53
	AuthorizationIdentity ApprovalKeyAuthorizationIdentity
	SignatureAccepted     bool
}

func cloneFixtureVector(vector FixtureVector) FixtureVector {
	vector.EnvelopeBytes = bytes.Clone(vector.EnvelopeBytes)
	vector.PayloadBytes = bytes.Clone(vector.PayloadBytes)
	vector.ProtectedHeaderBytes = bytes.Clone(vector.ProtectedHeaderBytes)
	vector.ProtectedKeyID = bytes.Clone(vector.ProtectedKeyID)
	return vector
}

type FixtureVerifier struct {
	vectors []FixtureVector
}

func NewFixtureVerifier(vectors []FixtureVector) (*FixtureVerifier, error) {
	result := &FixtureVerifier{vectors: make([]FixtureVector, len(vectors))}
	for index, vector := range vectors {
		if len(vector.EnvelopeBytes) == 0 {
			return nil, classified(ClassificationLocalFailure, "empty-fixture-envelope")
		}
		result.vectors[index] = cloneFixtureVector(vector)
	}
	return result, nil
}

type VerifiedApproval struct {
	envelopeBytes         []byte
	payloadBytes          []byte
	protectedHeaderBytes  []byte
	protectedKeyID        []byte
	view                  ApprovalGrant
	resolvedEpochSequence v0candidate.UInt53
	authorizationIdentity ApprovalKeyAuthorizationIdentity
}

func (approval *VerifiedApproval) EnvelopeBytes() []byte {
	return bytes.Clone(approval.envelopeBytes)
}

func (approval *VerifiedApproval) PayloadBytes() []byte {
	return bytes.Clone(approval.payloadBytes)
}

func (approval *VerifiedApproval) ProtectedHeaderBytes() []byte {
	return bytes.Clone(approval.protectedHeaderBytes)
}

func (approval *VerifiedApproval) ProtectedKeyID() []byte {
	return bytes.Clone(approval.protectedKeyID)
}

func (approval *VerifiedApproval) View() ApprovalGrant { return approval.view }

func (approval *VerifiedApproval) ResolvedEpochSequence() v0candidate.UInt53 {
	return approval.resolvedEpochSequence
}

func (approval *VerifiedApproval) AuthorizationIdentity() ApprovalKeyAuthorizationIdentity {
	return approval.authorizationIdentity
}

func (verifier *FixtureVerifier) Verify(
	ctx context.Context,
	received []byte,
	bindings ApprovalGrantRoleBindings,
) (*VerifiedApproval, error) {
	if err := ctx.Err(); err != nil {
		return nil, classified(ClassificationLocalFailure, "verification-cancelled")
	}
	if len(received) == 0 || len(received) > ApprovalEnvelopeRawMaxBytes {
		return nil, classified(ClassificationMalformed, "envelope-raw-byte-limit")
	}
	envelope := bytes.Clone(received)

	var matched *FixtureVector
	for index := range verifier.vectors {
		if bytes.Equal(envelope, verifier.vectors[index].EnvelopeBytes) {
			vector := cloneFixtureVector(verifier.vectors[index])
			matched = &vector
			break
		}
	}
	if matched == nil {
		return nil, classified(ClassificationUnsupported, "unknown-fixture-envelope")
	}
	if len(matched.PayloadBytes) > ApprovalPayloadRawMaxBytes {
		return nil, classified(ClassificationMalformed, "payload-raw-byte-limit")
	}
	if len(matched.ProtectedHeaderBytes) > ApprovalProtectedRawMaxBytes {
		return nil, classified(ClassificationMalformed, "protected-raw-byte-limit")
	}
	if len(envelope) > ApprovalEnvelopeCalculatedMaxBytes ||
		len(matched.PayloadBytes) > ApprovalPayloadCalculatedMaxBytes ||
		len(matched.ProtectedHeaderBytes) > ApprovalProtectedCalculatedMaxBytes {
		return nil, classified(ClassificationSchema, "calculated-candidate-maximum")
	}
	if !matched.SignatureAccepted {
		return nil, classified(ClassificationUnsupported, "fixture-signature-rejected")
	}
	if len(matched.ProtectedKeyID) == 0 || len(matched.ProtectedKeyID) > 64 ||
		matched.AuthorizationIdentity == (ApprovalKeyAuthorizationIdentity{}) {
		return nil, classified(ClassificationSchema, "fixture-key-authorization-shape")
	}
	if err := validateGrantShape(matched.View); err != nil {
		return nil, err
	}
	if err := bindGrant(
		matched.View,
		matched.ResolvedEpochSequence,
		matched.ProtectedKeyID,
		matched.AuthorizationIdentity,
		bindings,
	); err != nil {
		return nil, err
	}
	return &VerifiedApproval{
		envelopeBytes:         envelope,
		payloadBytes:          bytes.Clone(matched.PayloadBytes),
		protectedHeaderBytes:  bytes.Clone(matched.ProtectedHeaderBytes),
		protectedKeyID:        bytes.Clone(matched.ProtectedKeyID),
		view:                  matched.View,
		resolvedEpochSequence: matched.ResolvedEpochSequence,
		authorizationIdentity: matched.AuthorizationIdentity,
	}, nil
}

func (verifier *FixtureVerifier) VerifyCandidate(
	ctx context.Context,
	received []byte,
) (*VerifiedApproval, error) {
	if err := ctx.Err(); err != nil {
		return nil, classified(ClassificationLocalFailure, "verification-cancelled")
	}
	if len(received) == 0 || len(received) > ApprovalEnvelopeRawMaxBytes {
		return nil, classified(ClassificationMalformed, "envelope-raw-byte-limit")
	}
	envelope := bytes.Clone(received)
	var matched *FixtureVector
	for index := range verifier.vectors {
		if bytes.Equal(envelope, verifier.vectors[index].EnvelopeBytes) {
			vector := cloneFixtureVector(verifier.vectors[index])
			matched = &vector
			break
		}
	}
	if matched == nil {
		return nil, classified(ClassificationUnsupported, "unknown-fixture-envelope")
	}
	if len(matched.PayloadBytes) > ApprovalPayloadRawMaxBytes {
		return nil, classified(ClassificationMalformed, "payload-raw-byte-limit")
	}
	if len(matched.ProtectedHeaderBytes) > ApprovalProtectedRawMaxBytes {
		return nil, classified(ClassificationMalformed, "protected-raw-byte-limit")
	}
	if len(envelope) > ApprovalEnvelopeCalculatedMaxBytes ||
		len(matched.PayloadBytes) > ApprovalPayloadCalculatedMaxBytes ||
		len(matched.ProtectedHeaderBytes) > ApprovalProtectedCalculatedMaxBytes {
		return nil, classified(ClassificationSchema, "calculated-candidate-maximum")
	}
	if !matched.SignatureAccepted {
		return nil, classified(ClassificationUnsupported, "fixture-signature-rejected")
	}
	if len(matched.ProtectedKeyID) == 0 || len(matched.ProtectedKeyID) > 64 ||
		matched.AuthorizationIdentity == (ApprovalKeyAuthorizationIdentity{}) {
		return nil, classified(ClassificationSchema, "fixture-key-authorization-shape")
	}
	if err := validateGrantShape(matched.View); err != nil {
		return nil, err
	}
	return &VerifiedApproval{
		envelopeBytes:         envelope,
		payloadBytes:          bytes.Clone(matched.PayloadBytes),
		protectedHeaderBytes:  bytes.Clone(matched.ProtectedHeaderBytes),
		protectedKeyID:        bytes.Clone(matched.ProtectedKeyID),
		view:                  matched.View,
		resolvedEpochSequence: matched.ResolvedEpochSequence,
		authorizationIdentity: matched.AuthorizationIdentity,
	}, nil
}

func validateGrantShape(view ApprovalGrant) error {
	if view.ObjectType != ApprovalGrantObjectType {
		return classified(ClassificationUnsupported, "approval-object-type")
	}
	if view.ObjectVersion != 0 {
		return classified(ClassificationUnsupported, "approval-object-version")
	}
	if view.InstallationID == ([16]byte{}) || view.RegistrationID == ([16]byte{}) ||
		view.SupervisorID == ([16]byte{}) || view.AttemptNonce == (AttemptNonce{}) {
		return classified(ClassificationSchema, "approval-identifier-zero")
	}
	if view.Purpose != ApprovalGrantPurpose || view.Audience != ApprovalGrantAudience {
		return classified(ClassificationBinding, "approval-purpose-or-audience")
	}
	if view.IssuedAt >= view.ExpiresAt {
		return classified(ClassificationSchema, "approval-time-order")
	}
	return nil
}

func bindGrant(
	view ApprovalGrant,
	resolvedEpochSequence v0candidate.UInt53,
	protectedKeyID []byte,
	authorizationIdentity ApprovalKeyAuthorizationIdentity,
	bindings ApprovalGrantRoleBindings,
) error {
	bindings = cloneRoleBindings(bindings)
	if view.InstallationID != bindings.InstallationID ||
		resolvedEpochSequence != bindings.EpochSequence ||
		view.EpochDigest != bindings.EpochDigest ||
		view.RegistrationID != bindings.RegistrationID ||
		view.PlanDigest != bindings.PlanDigest ||
		view.SupervisorID != bindings.SupervisorID ||
		view.AttemptNonce != bindings.AttemptNonce ||
		!bytes.Equal(protectedKeyID, bindings.ProtectedKeyID) ||
		authorizationIdentity != bindings.AuthorizationIdentity {
		return classified(ClassificationBinding, "approval-role-binding")
	}
	return nil
}
