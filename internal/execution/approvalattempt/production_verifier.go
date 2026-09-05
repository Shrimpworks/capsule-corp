package approvalattempt

import (
	"bytes"
	"context"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"capsule.local/capsule/internal/protocol/cborscan"
	"capsule.local/capsule/internal/protocol/v0candidate"
	"github.com/fxamacker/cbor/v2"
)

const (
	ApprovalGrantMediaType = "application/capsule.approval-grant+cbor;v=0"

	ApprovalKeyTeamID            = "3DDR84M4JS"
	ApprovalBrokerRoleID         = "capsule.role.approval-broker/v0"
	ApprovalKeyAccessGroupPrefix = "3DDR84M4JS.com.capsulecorp.capsule.broker.approval.epoch-"
	ApprovalKeyProtection        = "secure-enclave-p256"
	ApprovalKeyAccessControl     = "userPresence|privateKeyUsage"
	ApprovalKeyContextPolicy     = "fresh-nonreused-lacontext-per-operation"
	ApprovalKeyStatusActive      = "active"

	ApprovalProductionEnvelopeCalculatedMaxBytes  = 399
	ApprovalProductionPayloadCalculatedMaxBytes   = 242
	ApprovalProductionProtectedCalculatedMaxBytes = 84
	ApprovalKeyAuthorizationMaxEntries            = 32

	approvalAuthorizationDomain = "capsule.approval-key-authorization/v0\x00"
)

// ApprovalP256PublicKey is public fixture or locally authorized key material.
// No private-key representation exists in this package.
type ApprovalP256PublicKey struct {
	X [32]byte
	Y [32]byte
}

// ApprovalKeyAuthorization is a passive projection of trusted local
// authorization state. Constructing this value does not enroll or authorize a
// key; callers must supply only installation-root-authorized active state.
type ApprovalKeyAuthorization struct {
	TeamID                string
	RoleID                string
	AccessGroup           string
	KeyID                 [32]byte
	PublicKey             ApprovalP256PublicKey
	InstallationID        v0candidate.InstallationID
	EpochSequence         v0candidate.UInt53
	EpochDigest           v0candidate.TrustEpochDigest
	Purpose               string
	Audience              string
	NotBefore             v0candidate.UInt53
	ExpiresAt             v0candidate.UInt53
	Protection            string
	AccessControl         string
	ContextPolicy         string
	SoftwareFallback      bool
	Status                string
	AuthorizationIdentity ApprovalKeyAuthorizationIdentity
}

// ApprovalAccessGroup returns the exact epoch-scoped Keychain access group for
// the Broker build admitted in that security epoch.
func ApprovalAccessGroup(epoch v0candidate.UInt53) string {
	return ApprovalKeyAccessGroupPrefix + strconv.FormatUint(uint64(epoch), 10)
}

// NewPassiveApprovalKeyAuthorization builds the exact policy projection used
// by public-key-only conformance fixtures. It performs no Keychain operation.
func NewPassiveApprovalKeyAuthorization(
	publicKey ApprovalP256PublicKey,
	installationID v0candidate.InstallationID,
	epochSequence v0candidate.UInt53,
	epochDigest v0candidate.TrustEpochDigest,
	notBefore v0candidate.UInt53,
	expiresAt v0candidate.UInt53,
) (ApprovalKeyAuthorization, error) {
	keyID, err := approvalPublicKeyID(publicKey)
	if err != nil {
		return ApprovalKeyAuthorization{}, err
	}
	authorization := ApprovalKeyAuthorization{
		TeamID: ApprovalKeyTeamID, RoleID: ApprovalBrokerRoleID,
		AccessGroup: ApprovalAccessGroup(epochSequence), KeyID: keyID, PublicKey: publicKey,
		InstallationID: installationID, EpochSequence: epochSequence, EpochDigest: epochDigest,
		Purpose: ApprovalGrantPurpose, Audience: ApprovalGrantAudience,
		NotBefore: notBefore, ExpiresAt: expiresAt,
		Protection: ApprovalKeyProtection, AccessControl: ApprovalKeyAccessControl,
		ContextPolicy: ApprovalKeyContextPolicy, SoftwareFallback: false, Status: ApprovalKeyStatusActive,
	}
	authorization.AuthorizationIdentity = approvalAuthorizationIdentity(authorization)
	if err := validateApprovalKeyAuthorization(authorization); err != nil {
		return ApprovalKeyAuthorization{}, err
	}
	return authorization, nil
}

// ProductionShapedVerifier is production-shaped but unwired. It accepts only exact
// Capsule ApprovalGrant Sign1 bytes and trusted public-key authorizations. It
// never signs, opens Keychain, prompts, persists replay state, or calls IPC.
type ProductionShapedVerifier struct {
	encode         cbor.EncMode
	decode         cbor.DecMode
	authorizations map[[32]byte]ApprovalKeyAuthorization
}

var _ Verifier = (*ProductionShapedVerifier)(nil)
var _ CandidateVerifier = (*ProductionShapedVerifier)(nil)

func NewProductionShapedVerifier(authorizations []ApprovalKeyAuthorization) (*ProductionShapedVerifier, error) {
	if len(authorizations) == 0 || len(authorizations) > ApprovalKeyAuthorizationMaxEntries {
		return nil, classified(ClassificationSchema, "approval-key-authorization-count")
	}
	encOptions := cbor.CoreDetEncOptions()
	encOptions.TagsMd = cbor.TagsForbidden
	encOptions.BinaryMarshaler = cbor.BinaryMarshalerNone
	encOptions.TextMarshaler = cbor.TextMarshalerNone
	encode, err := encOptions.EncMode()
	if err != nil {
		return nil, fmt.Errorf("construct approval deterministic encoder: %w", err)
	}
	decode, err := (cbor.DecOptions{
		DupMapKey: cbor.DupMapKeyEnforcedAPF, MaxNestedLevels: 4,
		MaxArrayElements: 16, MaxMapPairs: 16, IndefLength: cbor.IndefLengthForbidden,
		TagsMd: cbor.TagsForbidden, UTF8: cbor.UTF8RejectInvalid,
		ExtraReturnErrors: cbor.ExtraDecErrorUnknownField, BignumTag: cbor.BignumTagForbidden,
		BinaryUnmarshaler: cbor.BinaryUnmarshalerNone, TextUnmarshaler: cbor.TextUnmarshalerNone,
	}).DecMode()
	if err != nil {
		return nil, fmt.Errorf("construct approval typed decoder: %w", err)
	}
	verifier := &ProductionShapedVerifier{encode: encode, decode: decode, authorizations: make(map[[32]byte]ApprovalKeyAuthorization, len(authorizations))}
	for _, authorization := range authorizations {
		if err := validateApprovalKeyAuthorization(authorization); err != nil {
			return nil, err
		}
		if _, exists := verifier.authorizations[authorization.KeyID]; exists {
			return nil, classified(ClassificationDomain, "duplicate-approval-key-id")
		}
		verifier.authorizations[authorization.KeyID] = authorization
	}
	return verifier, nil
}

type approvalProtectedWire struct {
	Algorithm   int64  `cbor:"1,keyasint"`
	ContentType string `cbor:"3,keyasint"`
	KeyID       []byte `cbor:"4,keyasint"`
}

type approvalPayloadWire struct {
	ObjectType     string `cbor:"1,keyasint"`
	ObjectVersion  uint64 `cbor:"2,keyasint"`
	InstallationID []byte `cbor:"3,keyasint"`
	EpochDigest    []byte `cbor:"4,keyasint"`
	RegistrationID []byte `cbor:"5,keyasint"`
	PlanDigest     []byte `cbor:"6,keyasint"`
	SupervisorID   []byte `cbor:"7,keyasint"`
	AttemptNonce   []byte `cbor:"8,keyasint"`
	Purpose        string `cbor:"9,keyasint"`
	Audience       string `cbor:"10,keyasint"`
	IssuedAt       uint64 `cbor:"11,keyasint"`
	ExpiresAt      uint64 `cbor:"12,keyasint"`
}

// VerifyCandidate is the production-shaped approval-submission verification
// seam. It derives the nonce, protected key ID, and authorization identity
// only from the signed envelope and the verifier's trusted authorization set.
// It deliberately does not accept Supervisor state from the caller: the
// durable approval component applies registration, active-state, effective-
// time, replay, and nonce-uniqueness rules after this method returns.
func (verifier *ProductionShapedVerifier) VerifyCandidate(
	ctx context.Context,
	received []byte,
) (*VerifiedApproval, error) {
	verified, _, err := verifier.verifyCandidate(ctx, received)
	return verified, err
}

func (verifier *ProductionShapedVerifier) Verify(ctx context.Context, received []byte, bindings ApprovalGrantRoleBindings) (*VerifiedApproval, error) {
	verified, authorization, err := verifier.verifyCandidate(ctx, received)
	if err != nil {
		return nil, err
	}
	view := verified.View()
	if err := bindApprovalAuthorization(authorization, view, bindings); err != nil {
		return nil, err
	}
	if err := validateProductionApprovalTime(view, authorization, bindings); err != nil {
		return nil, err
	}
	if err := bindGrant(view, authorization.EpochSequence, verified.ProtectedKeyID(), authorization.AuthorizationIdentity, bindings); err != nil {
		return nil, err
	}
	return verified, nil
}

func (verifier *ProductionShapedVerifier) verifyCandidate(
	ctx context.Context,
	received []byte,
) (*VerifiedApproval, ApprovalKeyAuthorization, error) {
	if verifier == nil {
		return nil, ApprovalKeyAuthorization{}, classified(ClassificationLocalFailure, "production-verifier-required")
	}
	if err := ctx.Err(); err != nil {
		return nil, ApprovalKeyAuthorization{}, classified(ClassificationLocalFailure, "verification-cancelled")
	}
	parts, err := frameApprovalEnvelope(received)
	if err != nil {
		return nil, ApprovalKeyAuthorization{}, err
	}
	if len(received) > ApprovalProductionEnvelopeCalculatedMaxBytes ||
		parts.payloadEnd-parts.payloadStart > ApprovalProductionPayloadCalculatedMaxBytes ||
		parts.protectedEnd-parts.protectedStart > ApprovalProductionProtectedCalculatedMaxBytes {
		return nil, ApprovalKeyAuthorization{}, classified(ClassificationSchema, "production-calculated-maximum")
	}
	exact := bytes.Clone(received)
	protected := exact[parts.protectedStart:parts.protectedEnd]
	payload := exact[parts.payloadStart:parts.payloadEnd]
	signature := exact[parts.signatureStart:parts.signatureEnd]
	header, err := verifier.decodeProtected(protected)
	if err != nil {
		return nil, ApprovalKeyAuthorization{}, err
	}
	view, err := verifier.decodePayload(payload)
	if err != nil {
		return nil, ApprovalKeyAuthorization{}, err
	}
	authorization, ok := verifier.authorizations[header.KeyID]
	if !ok {
		return nil, ApprovalKeyAuthorization{}, classified(ClassificationBinding, "approval-key-not-authorized")
	}
	if err := verifyApprovalSignature(verifier.encode, authorization.PublicKey, protected, payload, signature); err != nil {
		return nil, ApprovalKeyAuthorization{}, err
	}
	if err := bindCandidateApprovalAuthorization(authorization, view); err != nil {
		return nil, ApprovalKeyAuthorization{}, err
	}
	if err := validateCandidateApprovalAuthorizationTime(view, authorization); err != nil {
		return nil, ApprovalKeyAuthorization{}, err
	}
	return &VerifiedApproval{
		envelopeBytes: exact, payloadBytes: bytes.Clone(payload), protectedHeaderBytes: bytes.Clone(protected),
		protectedKeyID: bytes.Clone(header.KeyID[:]), view: view,
		resolvedEpochSequence: authorization.EpochSequence,
		authorizationIdentity: authorization.AuthorizationIdentity,
	}, authorization, nil
}

type approvalProtectedView struct {
	Algorithm   int64
	ContentType string
	KeyID       [32]byte
}

func (verifier *ProductionShapedVerifier) decodeProtected(exact []byte) (approvalProtectedView, error) {
	if err := predecodeApprovalCBOR(exact, approvalPredecodeProfile{maxBytes: ApprovalProtectedRawMaxBytes, maxDepth: 2, maxItems: 7, maxMapEntries: 3}); err != nil {
		return approvalProtectedView{}, err
	}
	var wire approvalProtectedWire
	if err := verifier.decode.Unmarshal(exact, &wire); err != nil {
		return approvalProtectedView{}, approvalTypedDecodeError(err, "protected-header")
	}
	canonical, err := verifier.encode.Marshal(wire)
	if err != nil || !bytes.Equal(canonical, exact) {
		return approvalProtectedView{}, classified(ClassificationMalformed, "noncanonical-protected-header")
	}
	if wire.Algorithm != -7 || wire.ContentType != ApprovalGrantMediaType || len(wire.KeyID) != 32 {
		return approvalProtectedView{}, classified(ClassificationUnsupported, "protected-header-profile")
	}
	var keyID [32]byte
	copy(keyID[:], wire.KeyID)
	return approvalProtectedView{Algorithm: wire.Algorithm, ContentType: wire.ContentType, KeyID: keyID}, nil
}

func (verifier *ProductionShapedVerifier) decodePayload(exact []byte) (ApprovalGrant, error) {
	if err := predecodeApprovalCBOR(exact, approvalPredecodeProfile{maxBytes: ApprovalPayloadRawMaxBytes, maxDepth: 2, maxItems: 25, maxMapEntries: 12}); err != nil {
		return ApprovalGrant{}, err
	}
	var wire approvalPayloadWire
	if err := verifier.decode.Unmarshal(exact, &wire); err != nil {
		return ApprovalGrant{}, approvalTypedDecodeError(err, "approval-payload")
	}
	canonical, err := verifier.encode.Marshal(wire)
	if err != nil || !bytes.Equal(canonical, exact) {
		return ApprovalGrant{}, classified(ClassificationMalformed, "noncanonical-approval-payload")
	}
	installationID, err := v0candidate.NewInstallationID(wire.InstallationID)
	if err != nil {
		return ApprovalGrant{}, classified(ClassificationSchema, "approval-installation-id")
	}
	epochDigest, err := v0candidate.NewTrustEpochDigest(wire.EpochDigest)
	if err != nil {
		return ApprovalGrant{}, classified(ClassificationSchema, "approval-epoch-digest")
	}
	registrationID, err := v0candidate.NewRegistrationID(wire.RegistrationID)
	if err != nil {
		return ApprovalGrant{}, classified(ClassificationSchema, "approval-registration-id")
	}
	planDigest, err := v0candidate.NewExecutionPlanDigest(wire.PlanDigest)
	if err != nil {
		return ApprovalGrant{}, classified(ClassificationSchema, "approval-plan-digest")
	}
	supervisorID, err := v0candidate.NewSupervisorID(wire.SupervisorID)
	if err != nil {
		return ApprovalGrant{}, classified(ClassificationSchema, "approval-supervisor-id")
	}
	nonceIdentifier, err := NewDomainIdentifier(DomainAttemptNonce, wire.AttemptNonce)
	if err != nil {
		return ApprovalGrant{}, err
	}
	attemptNonce, err := NewAttemptNonce(nonceIdentifier)
	if err != nil {
		return ApprovalGrant{}, err
	}
	issuedAt, err := v0candidate.NewUInt53(wire.IssuedAt)
	if err != nil {
		return ApprovalGrant{}, classified(ClassificationSchema, "approval-issued-at")
	}
	expiresAt, err := v0candidate.NewUInt53(wire.ExpiresAt)
	if err != nil {
		return ApprovalGrant{}, classified(ClassificationSchema, "approval-expires-at")
	}
	view := ApprovalGrant{
		ObjectType: wire.ObjectType, ObjectVersion: v0candidate.ObjectVersion(wire.ObjectVersion),
		InstallationID: installationID, EpochDigest: epochDigest, RegistrationID: registrationID,
		PlanDigest: planDigest, SupervisorID: supervisorID, AttemptNonce: attemptNonce,
		Purpose: wire.Purpose, Audience: wire.Audience, IssuedAt: issuedAt, ExpiresAt: expiresAt,
	}
	if err := validateGrantShape(view); err != nil {
		return ApprovalGrant{}, err
	}
	return view, nil
}

func approvalTypedDecodeError(err error, scope string) error {
	var unknown *cbor.UnknownFieldError
	if errors.As(err, &unknown) {
		return classified(ClassificationUnsupported, scope+"-unknown-field")
	}
	return classified(ClassificationSchema, scope+"-typed-decode")
}

func validateProductionApprovalTime(view ApprovalGrant, authorization ApprovalKeyAuthorization, bindings ApprovalGrantRoleBindings) error {
	if bindings.EffectiveNow > v0candidate.UInt53(v0candidate.MaxSafeInteger) || bindings.RegistrationExpiresAt > v0candidate.UInt53(v0candidate.MaxSafeInteger) {
		return classified(ClassificationSchema, "approval-time-binding-range")
	}
	if view.ExpiresAt-view.IssuedAt > 300 || view.IssuedAt > bindings.EffectiveNow || bindings.EffectiveNow >= view.ExpiresAt ||
		bindings.EffectiveNow >= bindings.RegistrationExpiresAt || view.ExpiresAt > bindings.RegistrationExpiresAt {
		return classified(ClassificationStale, "approval-time-binding")
	}
	if view.IssuedAt < authorization.NotBefore || view.ExpiresAt > authorization.ExpiresAt ||
		bindings.EffectiveNow < authorization.NotBefore || bindings.EffectiveNow >= authorization.ExpiresAt {
		return classified(ClassificationStale, "approval-key-validity")
	}
	return nil
}

func validateCandidateApprovalAuthorizationTime(view ApprovalGrant, authorization ApprovalKeyAuthorization) error {
	if view.ExpiresAt-view.IssuedAt > 300 || view.IssuedAt < authorization.NotBefore ||
		view.ExpiresAt > authorization.ExpiresAt {
		return classified(ClassificationStale, "approval-key-validity")
	}
	return nil
}

func bindCandidateApprovalAuthorization(authorization ApprovalKeyAuthorization, view ApprovalGrant) error {
	if authorization.TeamID != ApprovalKeyTeamID || authorization.RoleID != ApprovalBrokerRoleID ||
		authorization.AccessGroup != ApprovalAccessGroup(authorization.EpochSequence) ||
		authorization.Protection != ApprovalKeyProtection || authorization.AccessControl != ApprovalKeyAccessControl ||
		authorization.ContextPolicy != ApprovalKeyContextPolicy || authorization.SoftwareFallback ||
		authorization.Status != ApprovalKeyStatusActive || authorization.Purpose != view.Purpose ||
		authorization.Audience != view.Audience || authorization.InstallationID != view.InstallationID ||
		authorization.EpochDigest != view.EpochDigest ||
		authorization.AuthorizationIdentity == (ApprovalKeyAuthorizationIdentity{}) {
		return classified(ClassificationBinding, "approval-key-authorization-binding")
	}
	return nil
}

func bindApprovalAuthorization(authorization ApprovalKeyAuthorization, view ApprovalGrant, bindings ApprovalGrantRoleBindings) error {
	if authorization.TeamID != ApprovalKeyTeamID || authorization.RoleID != ApprovalBrokerRoleID ||
		authorization.AccessGroup != ApprovalAccessGroup(authorization.EpochSequence) ||
		authorization.Protection != ApprovalKeyProtection || authorization.AccessControl != ApprovalKeyAccessControl ||
		authorization.ContextPolicy != ApprovalKeyContextPolicy || authorization.SoftwareFallback ||
		authorization.Status != ApprovalKeyStatusActive || authorization.Purpose != ApprovalGrantPurpose ||
		authorization.Audience != ApprovalGrantAudience || authorization.InstallationID != view.InstallationID ||
		authorization.InstallationID != bindings.InstallationID || authorization.EpochSequence != bindings.EpochSequence ||
		authorization.EpochDigest != view.EpochDigest || authorization.EpochDigest != bindings.EpochDigest ||
		authorization.AuthorizationIdentity != bindings.AuthorizationIdentity {
		return classified(ClassificationBinding, "approval-key-authorization-binding")
	}
	return nil
}

func validateApprovalKeyAuthorization(authorization ApprovalKeyAuthorization) error {
	if authorization.InstallationID == (v0candidate.InstallationID{}) || authorization.EpochSequence == 0 ||
		uint64(authorization.EpochSequence) > v0candidate.MaxSafeInteger ||
		uint64(authorization.NotBefore) > v0candidate.MaxSafeInteger || uint64(authorization.ExpiresAt) > v0candidate.MaxSafeInteger ||
		authorization.TeamID != ApprovalKeyTeamID || authorization.RoleID != ApprovalBrokerRoleID ||
		authorization.AccessGroup != ApprovalAccessGroup(authorization.EpochSequence) ||
		authorization.Purpose != ApprovalGrantPurpose || authorization.Audience != ApprovalGrantAudience ||
		authorization.Protection != ApprovalKeyProtection || authorization.AccessControl != ApprovalKeyAccessControl ||
		authorization.ContextPolicy != ApprovalKeyContextPolicy || authorization.SoftwareFallback ||
		authorization.Status != ApprovalKeyStatusActive || authorization.NotBefore >= authorization.ExpiresAt {
		return classified(ClassificationBinding, "approval-key-policy")
	}
	keyID, err := approvalPublicKeyID(authorization.PublicKey)
	if err != nil {
		return err
	}
	if authorization.KeyID != keyID || authorization.AuthorizationIdentity != approvalAuthorizationIdentity(authorization) ||
		authorization.AuthorizationIdentity == (ApprovalKeyAuthorizationIdentity{}) {
		return classified(ClassificationBinding, "approval-key-identity")
	}
	return nil
}

func approvalPublicKeyID(publicKey ApprovalP256PublicKey) ([32]byte, error) {
	encoded := encodeApprovalCOSEKey(publicKey)
	point := make([]byte, 65)
	point[0] = 4
	copy(point[1:33], publicKey.X[:])
	copy(point[33:], publicKey.Y[:])
	if _, err := ecdh.P256().NewPublicKey(point); err != nil {
		return [32]byte{}, classified(ClassificationSchema, "approval-public-key-point")
	}
	return sha256.Sum256(encoded), nil
}

func encodeApprovalCOSEKey(publicKey ApprovalP256PublicKey) []byte {
	result := []byte{0xa5, 0x01, 0x02, 0x03, 0x26, 0x20, 0x01, 0x21, 0x58, 0x20}
	result = append(result, publicKey.X[:]...)
	result = append(result, 0x22, 0x58, 0x20)
	result = append(result, publicKey.Y[:]...)
	return result
}

func approvalAuthorizationIdentity(authorization ApprovalKeyAuthorization) ApprovalKeyAuthorizationIdentity {
	preimage := []byte(approvalAuthorizationDomain)
	appendText := func(value string) {
		preimage = append(preimage, value...)
		// Every accepted policy value is an exact ASCII constant or the fixed
		// lowercase-hex access group, so NUL is an unambiguous separator.
		preimage = append(preimage, 0)
	}
	appendText(authorization.TeamID)
	appendText(authorization.RoleID)
	appendText(authorization.AccessGroup)
	preimage = append(preimage, authorization.KeyID[:]...)
	preimage = append(preimage, encodeApprovalCOSEKey(authorization.PublicKey)...)
	preimage = append(preimage, authorization.InstallationID[:]...)
	preimage = binary.BigEndian.AppendUint64(preimage, uint64(authorization.EpochSequence))
	preimage = append(preimage, authorization.EpochDigest[:]...)
	appendText(authorization.Purpose)
	appendText(authorization.Audience)
	preimage = binary.BigEndian.AppendUint64(preimage, uint64(authorization.NotBefore))
	preimage = binary.BigEndian.AppendUint64(preimage, uint64(authorization.ExpiresAt))
	appendText(authorization.Protection)
	appendText(authorization.AccessControl)
	appendText(authorization.ContextPolicy)
	appendText(authorization.Status)
	preimage = append(preimage, 0) // softwareFallback == false is the only accepted encoding.
	return ApprovalKeyAuthorizationIdentity(sha256.Sum256(preimage))
}

func verifyApprovalSignature(encode cbor.EncMode, key ApprovalP256PublicKey, protected, payload, signature []byte) error {
	if len(signature) != 64 {
		return classified(ClassificationMalformed, "approval-signature-shape")
	}
	structure, err := encode.Marshal([]any{"Signature1", protected, []byte{}, payload})
	if err != nil {
		return classified(ClassificationLocalFailure, "approval-signature-structure")
	}
	digest := sha256.Sum256(structure)
	x := new(big.Int).SetBytes(key.X[:])
	y := new(big.Int).SetBytes(key.Y[:])
	if !ecdsa.Verify(&ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, digest[:], new(big.Int).SetBytes(signature[:32]), new(big.Int).SetBytes(signature[32:])) {
		return classified(ClassificationBinding, "approval-signature")
	}
	return nil
}

type approvalEnvelopeParts struct {
	protectedStart int
	protectedEnd   int
	payloadStart   int
	payloadEnd     int
	signatureStart int
	signatureEnd   int
}

func frameApprovalEnvelope(received []byte) (approvalEnvelopeParts, error) {
	if len(received) == 0 || len(received) > ApprovalEnvelopeRawMaxBytes {
		return approvalEnvelopeParts{}, classified(ClassificationMalformed, "envelope-raw-byte-limit")
	}
	if len(received) < 2 || received[0] != 0xd2 || received[1] != 0x84 {
		return approvalEnvelopeParts{}, classified(ClassificationMalformed, "approval-sign1-framing")
	}
	offset := 2
	protectedStart, protectedEnd, next, err := framedApprovalByteString(received, offset, ApprovalProtectedRawMaxBytes)
	if err != nil {
		return approvalEnvelopeParts{}, err
	}
	offset = next
	if offset >= len(received) || received[offset] != 0xa0 {
		return approvalEnvelopeParts{}, classified(ClassificationMalformed, "approval-unprotected-header")
	}
	offset++
	payloadStart, payloadEnd, next, err := framedApprovalByteString(received, offset, ApprovalPayloadRawMaxBytes)
	if err != nil {
		return approvalEnvelopeParts{}, err
	}
	offset = next
	signatureStart, signatureEnd, next, err := framedApprovalByteString(received, offset, 64)
	if err != nil {
		return approvalEnvelopeParts{}, err
	}
	if signatureEnd-signatureStart != 64 {
		return approvalEnvelopeParts{}, classified(ClassificationMalformed, "approval-signature-shape")
	}
	if next != len(received) {
		return approvalEnvelopeParts{}, classified(ClassificationMalformed, "approval-envelope-trailing")
	}
	return approvalEnvelopeParts{protectedStart, protectedEnd, payloadStart, payloadEnd, signatureStart, signatureEnd}, nil
}

func framedApprovalByteString(value []byte, offset int, capBytes uint64) (int, int, int, error) {
	if offset >= len(value) || value[offset]>>5 != 2 {
		return 0, 0, 0, classified(ClassificationMalformed, "approval-expected-byte-string")
	}
	initial := value[offset]
	offset++
	additional := initial & 0x1f
	var length uint64
	switch {
	case additional < 24:
		length = uint64(additional)
	case additional == 24:
		if offset >= len(value) {
			return 0, 0, 0, classified(ClassificationMalformed, "approval-truncated-byte-string")
		}
		length = uint64(value[offset])
		offset++
		if length < 24 {
			return 0, 0, 0, classified(ClassificationMalformed, "approval-nonpreferred-byte-string")
		}
	case additional == 25:
		if offset+2 > len(value) {
			return 0, 0, 0, classified(ClassificationMalformed, "approval-truncated-byte-string")
		}
		length = uint64(binary.BigEndian.Uint16(value[offset : offset+2]))
		offset += 2
		if length < 256 {
			return 0, 0, 0, classified(ClassificationMalformed, "approval-nonpreferred-byte-string")
		}
	default:
		return 0, 0, 0, classified(ClassificationMalformed, "approval-byte-string-length")
	}
	if length > capBytes || length > ApprovalEnvelopeRawMaxBytes {
		return 0, 0, 0, classified(ClassificationMalformed, "approval-nested-raw-byte-limit")
	}
	lengthInt := int(length)
	if lengthInt > len(value)-offset {
		return 0, 0, 0, classified(ClassificationMalformed, "approval-nested-raw-byte-limit")
	}
	end := offset + lengthInt
	return offset, end, end, nil
}

// approvalPredecodeProfile is the approval object set's view of the shared
// cborscan profile. The set carries only maps, integers, and byte/text
// strings, so the cborscan.Profile built from it leaves AllowArray, AllowTag,
// and AllowFalse unset and every array, tag, and simple/float value is
// refused as a major-type error.
type approvalPredecodeProfile struct {
	maxBytes      int
	maxDepth      int
	maxItems      int
	maxMapEntries uint64
}

// sharedProfile projects this package's bounds onto cborscan's. Note that
// MaxArrayElements is deliberately left at zero: it is unreachable, because
// AllowArray being unset refuses every array before any element count is
// considered.
func (profile approvalPredecodeProfile) sharedProfile() cborscan.Profile {
	return cborscan.Profile{
		MaxBytes:      profile.maxBytes,
		MaxDepth:      profile.maxDepth,
		MaxItems:      profile.maxItems,
		MaxMapEntries: profile.maxMapEntries,
	}
}

// predecodeApprovalCBOR runs the shared deterministic-CBOR predecoder over an
// approval envelope component and maps its refusal onto this package's own
// code vocabulary. Every refusal here is ClassificationMalformed, matching the
// hand-rolled scanner this replaces.
func predecodeApprovalCBOR(value []byte, profile approvalPredecodeProfile) error {
	if err := cborscan.Predecode(value, profile.sharedProfile()); err != nil {
		return classified(ClassificationMalformed, approvalPredecodeRefusalCode(err))
	}
	return nil
}

// approvalPredecodeRefusalCode maps a cborscan Reason onto the refusal code
// this package emitted before the two predecoders were consolidated. The codes
// are contract surface, so the mapping preserves two collapses cborscan itself
// does not make: empty payload and over-limit both report a byte-limit
// refusal, and a truncated head and a truncated argument both report a
// truncation. TestApprovalPredecodeRefusalCodeCoversEveryReason fails if a
// Reason is added to cborscan without a code here.
func approvalPredecodeRefusalCode(err error) string {
	var scanErr *cborscan.Error
	if !errors.As(err, &scanErr) {
		return approvalUnmappedPredecodeCode
	}
	switch scanErr.Reason {
	case cborscan.ReasonEmptyPayload, cborscan.ReasonRawByteLimit:
		return "approval-cbor-byte-limit"
	case cborscan.ReasonTrailingData:
		return "approval-cbor-trailing"
	case cborscan.ReasonDepthLimit:
		return "approval-cbor-depth"
	case cborscan.ReasonItemLimit:
		return "approval-cbor-items"
	case cborscan.ReasonTruncatedItem, cborscan.ReasonTruncatedArgument:
		return "approval-cbor-truncated"
	case cborscan.ReasonUnsafeInteger:
		return "approval-unsafe-integer"
	case cborscan.ReasonUnsafeNegativeInteger:
		return "approval-unsafe-negative-integer"
	case cborscan.ReasonTruncatedString, cborscan.ReasonDecodedStringLimit:
		return "approval-string-limit"
	case cborscan.ReasonInvalidUTF8:
		return "approval-invalid-utf8"
	case cborscan.ReasonMapEntryLimit:
		return "approval-map-entries"
	case cborscan.ReasonDuplicateMapKeyOrder:
		return "approval-map-key-order"
	case cborscan.ReasonMajorType, cborscan.ReasonArrayElementLimit,
		cborscan.ReasonSemanticTag, cborscan.ReasonSimpleOrFloat,
		cborscan.ReasonInvalidItem:
		// Every major type the approval object set does not contain reports
		// one type refusal. ReasonArrayElementLimit cannot be reached while
		// AllowArray is unset, and ReasonInvalidItem is cborscan's own
		// documented-unreachable guard; both are mapped so that enabling an
		// array profile later cannot produce an unmapped code.
		return "approval-cbor-type"
	case cborscan.ReasonIndefiniteOrReserved:
		return "approval-cbor-indefinite"
	case cborscan.ReasonNonpreferredArgument:
		return "approval-cbor-nonpreferred"
	default:
		return approvalUnmappedPredecodeCode
	}
}

// approvalUnmappedPredecodeCode is returned only for a refusal this package
// does not recognize. It is asserted unreachable by
// TestApprovalPredecodeRefusalCodeCoversEveryReason; reaching it in production
// would mean cborscan grew a Reason without a corresponding code here.
const approvalUnmappedPredecodeCode = "approval-cbor-unmapped"
