package bootstrapauthpassive

import (
	"bytes"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/fxamacker/cbor/v2"
)

type requestWire struct {
	ObjectType                          string `cbor:"1,keyasint"`
	ObjectVersion                       uint64 `cbor:"2,keyasint"`
	Purpose                             string `cbor:"3,keyasint"`
	Audience                            string `cbor:"4,keyasint"`
	InstallationID                      []byte `cbor:"5,keyasint"`
	InstallationRootPublicKey           []byte `cbor:"6,keyasint"`
	SignerKeyID                         []byte `cbor:"7,keyasint"`
	SignerAuthorizationIdentity         []byte `cbor:"8,keyasint"`
	SupervisorID                        []byte `cbor:"9,keyasint"`
	ExpectedEffectiveUID                uint64 `cbor:"10,keyasint"`
	ContainingReleaseDigest             []byte `cbor:"11,keyasint"`
	I1ProfileCandidateDigest            []byte `cbor:"12,keyasint"`
	ComponentProfileDigest              []byte `cbor:"13,keyasint"`
	InstallationManifestCandidateDigest []byte `cbor:"14,keyasint"`
	TrustEpochSequence                  uint64 `cbor:"15,keyasint"`
	TrustEpochCandidateDigest           []byte `cbor:"16,keyasint"`
	StateRootClass                      string `cbor:"17,keyasint"`
	StateRootName                       string `cbor:"18,keyasint"`
	RequestEntryName                    string `cbor:"19,keyasint"`
	RecordEntryName                     string `cbor:"20,keyasint"`
	OwnerLockEntryName                  string `cbor:"21,keyasint"`
	OwnerLockMode                       uint64 `cbor:"22,keyasint"`
	OwnerLockLinkCount                  uint64 `cbor:"23,keyasint"`
	OwnerLockProfile                    string `cbor:"24,keyasint"`
	StoreEntryName                      string `cbor:"25,keyasint"`
	StoreFormat                         string `cbor:"26,keyasint"`
	StoreProfile                        string `cbor:"27,keyasint"`
	AttemptsEnabled                     bool   `cbor:"28,keyasint"`
	FakeBackendCreatesGuest             bool   `cbor:"29,keyasint"`
	Nonce                               []byte `cbor:"30,keyasint"`
	IssuedAt                            uint64 `cbor:"31,keyasint"`
	NotBefore                           uint64 `cbor:"32,keyasint"`
	ExpiresAt                           uint64 `cbor:"33,keyasint"`
	ReplayDisposition                   string `cbor:"34,keyasint"`
}

type recordWire struct {
	ObjectType                          string `cbor:"1,keyasint"`
	ObjectVersion                       uint64 `cbor:"2,keyasint"`
	Purpose                             string `cbor:"3,keyasint"`
	Audience                            string `cbor:"4,keyasint"`
	RequestPayloadDigest                []byte `cbor:"5,keyasint"`
	RequestEnvelopeDigest               []byte `cbor:"6,keyasint"`
	RequestEnvelopeLength               uint64 `cbor:"7,keyasint"`
	RequestNonce                        []byte `cbor:"8,keyasint"`
	InstallationID                      []byte `cbor:"9,keyasint"`
	InstallationRootPublicKey           []byte `cbor:"10,keyasint"`
	SignerKeyID                         []byte `cbor:"11,keyasint"`
	SignerAuthorizationIdentity         []byte `cbor:"12,keyasint"`
	SupervisorID                        []byte `cbor:"13,keyasint"`
	ExpectedEffectiveUID                uint64 `cbor:"14,keyasint"`
	ContainingReleaseDigest             []byte `cbor:"15,keyasint"`
	I1ProfileCandidateDigest            []byte `cbor:"16,keyasint"`
	ComponentProfileDigest              []byte `cbor:"17,keyasint"`
	InstallationManifestCandidateDigest []byte `cbor:"18,keyasint"`
	TrustEpochSequence                  uint64 `cbor:"19,keyasint"`
	TrustEpochCandidateDigest           []byte `cbor:"20,keyasint"`
	StateRootClass                      string `cbor:"21,keyasint"`
	StateRootName                       string `cbor:"22,keyasint"`
	StateRootDevice                     uint64 `cbor:"23,keyasint"`
	StateRootInode                      uint64 `cbor:"24,keyasint"`
	StateRootType                       string `cbor:"25,keyasint"`
	StateRootUID                        uint64 `cbor:"26,keyasint"`
	StateRootMode                       uint64 `cbor:"27,keyasint"`
	OwnerLockEntryName                  string `cbor:"28,keyasint"`
	OwnerLockDevice                     uint64 `cbor:"29,keyasint"`
	OwnerLockInode                      uint64 `cbor:"30,keyasint"`
	OwnerLockType                       string `cbor:"31,keyasint"`
	OwnerLockUID                        uint64 `cbor:"32,keyasint"`
	OwnerLockMode                       uint64 `cbor:"33,keyasint"`
	OwnerLockLinkCount                  uint64 `cbor:"34,keyasint"`
	OwnerLockProfile                    string `cbor:"35,keyasint"`
	StoreEntryName                      string `cbor:"36,keyasint"`
	StoreDevice                         uint64 `cbor:"37,keyasint"`
	StoreType                           string `cbor:"38,keyasint"`
	StoreUID                            uint64 `cbor:"39,keyasint"`
	StoreMode                           uint64 `cbor:"40,keyasint"`
	StoreLinkCount                      uint64 `cbor:"41,keyasint"`
	StoreFormat                         string `cbor:"42,keyasint"`
	StoreProfile                        string `cbor:"43,keyasint"`
	InitialStoreByteLength              uint64 `cbor:"44,keyasint"`
	InitialStoreDigest                  []byte `cbor:"45,keyasint"`
	RequestEntryName                    string `cbor:"46,keyasint"`
	RecordEntryName                     string `cbor:"47,keyasint"`
	RequestRetentionPolicy              string `cbor:"48,keyasint"`
	RecordRetentionPolicy               string `cbor:"49,keyasint"`
	InstallationTransition              string `cbor:"50,keyasint"`
	Transition                          string `cbor:"51,keyasint"`
	FinalizedAt                         uint64 `cbor:"52,keyasint"`
	AttemptsEnabled                     bool   `cbor:"53,keyasint"`
	RuntimeAdmitted                     bool   `cbor:"54,keyasint"`
	BackendAdmitted                     bool   `cbor:"55,keyasint"`
	GuestCreated                        bool   `cbor:"56,keyasint"`
	OneUseDisposition                   string `cbor:"57,keyasint"`
	ReplayDisposition                   string `cbor:"58,keyasint"`
	RequestIssuedAt                     uint64 `cbor:"59,keyasint"`
	RequestNotBefore                    uint64 `cbor:"60,keyasint"`
	RequestExpiresAt                    uint64 `cbor:"61,keyasint"`
	RequestReplayDisposition            string `cbor:"62,keyasint"`
	RequestObjectType                   string `cbor:"63,keyasint"`
	RequestObjectVersion                uint64 `cbor:"64,keyasint"`
	RequestPurpose                      string `cbor:"65,keyasint"`
	RequestAudience                     string `cbor:"66,keyasint"`
	IssuedAt                            uint64 `cbor:"67,keyasint"`
	NotBefore                           uint64 `cbor:"68,keyasint"`
	ExpiresAt                           uint64 `cbor:"69,keyasint"`
}

type protectedWire struct {
	Algorithm   int64  `cbor:"1,keyasint"`
	ContentType string `cbor:"3,keyasint"`
	KeyID       []byte `cbor:"4,keyasint"`
}

type protectedView struct {
	Algorithm   int64
	ContentType string
	KeyID       KeyID
}

type Codec struct {
	encode cbor.EncMode
	decode cbor.DecMode
}

func NewCodec() (*Codec, error) {
	encOptions := cbor.CoreDetEncOptions()
	encOptions.TagsMd = cbor.TagsForbidden
	encOptions.BinaryMarshaler = cbor.BinaryMarshalerNone
	encOptions.TextMarshaler = cbor.TextMarshalerNone
	encode, err := encOptions.EncMode()
	if err != nil {
		return nil, fmt.Errorf("construct bootstrap deterministic encoder: %w", err)
	}
	decode, err := (cbor.DecOptions{
		DupMapKey:        cbor.DupMapKeyEnforcedAPF,
		MaxNestedLevels:  4,
		MaxArrayElements: 16,
		// The library exposes only broader collection modes; Capsule's
		// allocation-independent predecoder enforces the exact 34/69-pair
		// request/record limits before this typed decoder runs.
		MaxMapPairs:       128,
		IndefLength:       cbor.IndefLengthForbidden,
		TagsMd:            cbor.TagsForbidden,
		UTF8:              cbor.UTF8RejectInvalid,
		ExtraReturnErrors: cbor.ExtraDecErrorUnknownField,
		BignumTag:         cbor.BignumTagForbidden,
		BinaryUnmarshaler: cbor.BinaryUnmarshalerNone,
		TextUnmarshaler:   cbor.TextUnmarshalerNone,
	}).DecMode()
	if err != nil {
		return nil, fmt.Errorf("construct bootstrap typed decoder: %w", err)
	}
	return &Codec{encode: encode, decode: decode}, nil
}

func (c *Codec) VerifyRequest(received []byte, expected Request, trustedNow uint64, replay ReplayState) (*VerifiedRequest, error) {
	if c == nil {
		return nil, errors.New("bootstrap codec is required")
	}
	parts, err := frameEnvelope(received, RequestEnvelopeRawMaxBytes, RequestProtectedRawMaxBytes, RequestPayloadRawMaxBytes)
	if err != nil {
		return nil, err
	}
	if len(received) > RequestEnvelopeCalculatedMaxBytes || parts.payloadEnd-parts.payloadStart > RequestPayloadCalculatedMaxBytes || parts.protectedEnd-parts.protectedStart > RequestProtectedCalculatedMaxBytes {
		return nil, rejected(ClassificationSchema, "calculated-maximum")
	}
	exact := bytes.Clone(received)
	protected := exact[parts.protectedStart:parts.protectedEnd]
	payload := exact[parts.payloadStart:parts.payloadEnd]
	signature := exact[parts.signatureStart:parts.signatureEnd]
	if err := predecode(exact, predecodeProfile{maxBytes: RequestEnvelopeRawMaxBytes, maxDepth: 3, maxItems: 7, maxMapEntries: 0, maxArrayElements: 4, allowedTag: 18, allowTag: true}); err != nil {
		return nil, err
	}
	header, err := c.decodeProtected(protected, RequestMediaType)
	if err != nil {
		return nil, err
	}
	view, err := c.decodeRequest(payload)
	if err != nil {
		return nil, err
	}
	if header.KeyID != view.SignerKeyID {
		return nil, rejected(ClassificationBinding, "protected-payload-key-id")
	}
	if err := validateRequestShape(view); err != nil {
		return nil, err
	}
	if view != expected {
		return nil, rejected(ClassificationBinding, "request-role-binding")
	}
	if trustedNow < view.NotBefore {
		return nil, rejected(ClassificationTime, "request-future")
	}
	if trustedNow >= view.ExpiresAt {
		return nil, rejected(ClassificationTime, "request-stale-or-expired")
	}
	if err := verifySignature(view.InstallationRootPublicKey, protected, payload, signature); err != nil {
		return nil, err
	}
	decision, err := evaluateRequestReplay(sha256.Sum256(exact), view.Nonce, replay)
	if err != nil {
		return nil, err
	}
	return &VerifiedRequest{exactEnvelope: exact, exactPayload: bytes.Clone(payload), view: view, decision: decision}, nil
}

type RecordBindings struct {
	Expected        Record
	RequestEnvelope []byte
	RequestPayload  []byte
	TrustedNow      uint64
	Replay          ReplayState
}

func (c *Codec) VerifyRecord(received []byte, bindings RecordBindings) (*VerifiedRecord, error) {
	if c == nil {
		return nil, errors.New("bootstrap codec is required")
	}
	parts, err := frameEnvelope(received, RecordEnvelopeRawMaxBytes, RecordProtectedRawMaxBytes, RecordPayloadRawMaxBytes)
	if err != nil {
		return nil, err
	}
	if len(received) > RecordEnvelopeCalculatedMaxBytes || parts.payloadEnd-parts.payloadStart > RecordPayloadCalculatedMaxBytes || parts.protectedEnd-parts.protectedStart > RecordProtectedCalculatedMaxBytes {
		return nil, rejected(ClassificationSchema, "calculated-maximum")
	}
	exact := bytes.Clone(received)
	protected := exact[parts.protectedStart:parts.protectedEnd]
	payload := exact[parts.payloadStart:parts.payloadEnd]
	signature := exact[parts.signatureStart:parts.signatureEnd]
	if err := predecode(exact, predecodeProfile{maxBytes: RecordEnvelopeRawMaxBytes, maxDepth: 3, maxItems: 7, maxMapEntries: 0, maxArrayElements: 4, allowedTag: 18, allowTag: true}); err != nil {
		return nil, err
	}
	header, err := c.decodeProtected(protected, RecordMediaType)
	if err != nil {
		return nil, err
	}
	view, err := c.decodeRecord(payload)
	if err != nil {
		return nil, err
	}
	if header.KeyID != view.SignerKeyID {
		return nil, rejected(ClassificationBinding, "protected-payload-key-id")
	}
	if err := validateRecordShape(view); err != nil {
		return nil, err
	}
	if view != bindings.Expected {
		return nil, rejected(ClassificationBinding, "record-role-binding")
	}
	if bindings.TrustedNow < view.NotBefore {
		return nil, rejected(ClassificationTime, "record-future")
	}
	if bindings.TrustedNow >= view.ExpiresAt {
		return nil, rejected(ClassificationTime, "record-stale-or-expired")
	}
	if len(bindings.RequestEnvelope) == 0 || len(bindings.RequestEnvelope) > RequestEnvelopeRawMaxBytes || len(bindings.RequestPayload) == 0 || len(bindings.RequestPayload) > RequestPayloadRawMaxBytes {
		return nil, rejected(ClassificationBinding, "request-binding-shape")
	}
	if view.RequestEnvelopeLength != uint64(len(bindings.RequestEnvelope)) || view.RequestEnvelopeDigest != sha256.Sum256(bindings.RequestEnvelope) || view.RequestPayloadDigest != sha256.Sum256(bindings.RequestPayload) {
		return nil, rejected(ClassificationBinding, "request-digest-binding")
	}
	boundRequest, err := c.decodeRequest(bindings.RequestPayload)
	if err != nil || view.RequestNonce != boundRequest.Nonce || view.RequestIssuedAt != boundRequest.IssuedAt || view.RequestNotBefore != boundRequest.NotBefore || view.RequestExpiresAt != boundRequest.ExpiresAt || view.RequestReplayDisposition != boundRequest.ReplayDisposition || view.RequestObjectType != boundRequest.ObjectType || view.RequestObjectVersion != boundRequest.ObjectVersion || view.RequestPurpose != boundRequest.Purpose || view.RequestAudience != boundRequest.Audience {
		return nil, rejected(ClassificationBinding, "request-stable-field-binding")
	}
	if err := verifySignature(view.InstallationRootPublicKey, protected, payload, signature); err != nil {
		return nil, err
	}
	decision, err := evaluateRecordReplay(sha256.Sum256(exact), view.RequestNonce, bindings.Replay)
	if err != nil {
		return nil, err
	}
	return &VerifiedRecord{exactEnvelope: exact, exactPayload: bytes.Clone(payload), view: view, decision: decision}, nil
}

func (c *Codec) decodeProtected(exact []byte, mediaType string) (protectedView, error) {
	if err := predecode(exact, predecodeProfile{maxBytes: 256, maxDepth: 2, maxItems: 7, maxMapEntries: 3}); err != nil {
		return protectedView{}, err
	}
	var wire protectedWire
	if err := c.decode.Unmarshal(exact, &wire); err != nil {
		return protectedView{}, typedDecodeError(err, "protected-header")
	}
	canonical, err := c.encode.Marshal(wire)
	if err != nil || !bytes.Equal(canonical, exact) {
		return protectedView{}, rejected(ClassificationMalformed, "noncanonical-protected-header")
	}
	if wire.Algorithm != -7 || wire.ContentType != mediaType || len(wire.KeyID) != 32 {
		return protectedView{}, rejected(ClassificationUnsupported, "protected-header-profile")
	}
	var keyID KeyID
	copy(keyID[:], wire.KeyID)
	return protectedView{Algorithm: wire.Algorithm, ContentType: wire.ContentType, KeyID: keyID}, nil
}

func (c *Codec) decodeRequest(exact []byte) (Request, error) {
	if err := predecode(exact, predecodeProfile{maxBytes: RequestPayloadRawMaxBytes, maxDepth: 2, maxItems: 69, maxMapEntries: 34, allowFalse: true}); err != nil {
		return Request{}, err
	}
	var wire requestWire
	if err := c.decode.Unmarshal(exact, &wire); err != nil {
		return Request{}, typedDecodeError(err, "request-payload")
	}
	canonical, err := c.encode.Marshal(wire)
	if err != nil || !bytes.Equal(canonical, exact) {
		return Request{}, rejected(ClassificationMalformed, "noncanonical-request-payload")
	}
	return requestFromWire(wire)
}

func (c *Codec) decodeRecord(exact []byte) (Record, error) {
	if err := predecode(exact, predecodeProfile{maxBytes: RecordPayloadRawMaxBytes, maxDepth: 2, maxItems: 139, maxMapEntries: 69, allowFalse: true}); err != nil {
		return Record{}, err
	}
	var wire recordWire
	if err := c.decode.Unmarshal(exact, &wire); err != nil {
		return Record{}, typedDecodeError(err, "record-payload")
	}
	canonical, err := c.encode.Marshal(wire)
	if err != nil || !bytes.Equal(canonical, exact) {
		return Record{}, rejected(ClassificationMalformed, "noncanonical-record-payload")
	}
	return recordFromWire(wire)
}

func typedDecodeError(err error, scope string) error {
	var unknown *cbor.UnknownFieldError
	if errors.As(err, &unknown) {
		return rejected(ClassificationUnsupported, scope+"-unknown-field")
	}
	return rejected(ClassificationSchema, scope+"-typed-decode")
}

func decodePublicKey(exact []byte) (InstallationRootPublicKey, error) {
	if len(exact) != 77 || exact[0] != 0xa5 || exact[1] != 0x01 || exact[2] != 0x02 || exact[3] != 0x03 || exact[4] != 0x26 || exact[5] != 0x20 || exact[6] != 0x01 || exact[7] != 0x21 || exact[8] != 0x58 || exact[9] != 0x20 || exact[42] != 0x22 || exact[43] != 0x58 || exact[44] != 0x20 {
		return InstallationRootPublicKey{}, rejected(ClassificationSchema, "installation-root-public-key")
	}
	var key InstallationRootPublicKey
	key.KeyType, key.Algorithm, key.Curve = 2, -7, 1
	copy(key.X[:], exact[10:42])
	copy(key.Y[:], exact[45:77])
	encodedPoint := make([]byte, 65)
	encodedPoint[0] = 4
	copy(encodedPoint[1:33], key.X[:])
	copy(encodedPoint[33:], key.Y[:])
	if _, err := ecdh.P256().NewPublicKey(encodedPoint); err != nil {
		return InstallationRootPublicKey{}, rejected(ClassificationSchema, "installation-root-public-key-point")
	}
	return key, nil
}

func encodePublicKey(key InstallationRootPublicKey) []byte {
	result := make([]byte, 77)
	copy(result, []byte{0xa5, 0x01, 0x02, 0x03, 0x26, 0x20, 0x01, 0x21, 0x58, 0x20})
	copy(result[10:42], key.X[:])
	copy(result[42:], []byte{0x22, 0x58, 0x20})
	copy(result[45:], key.Y[:])
	return result
}

func authorizationIdentity(keyID KeyID, installationID InstallationID, epoch uint64, epochDigest Digest) AuthorizationIdentity {
	h := sha256.New()
	_, _ = h.Write([]byte(AuthorizationIdentityDomain))
	_, _ = h.Write(keyID[:])
	_, _ = h.Write(installationID[:])
	var encodedEpoch [8]byte
	binary.BigEndian.PutUint64(encodedEpoch[:], epoch)
	_, _ = h.Write(encodedEpoch[:])
	_, _ = h.Write(epochDigest[:])
	var result AuthorizationIdentity
	copy(result[:], h.Sum(nil))
	return result
}

func verifySignature(key InstallationRootPublicKey, protected, payload, signature []byte) error {
	if len(signature) != 64 {
		return rejected(ClassificationMalformed, "signature-shape")
	}
	digest := sha256.Sum256(sigStructure(protected, payload))
	public := ecdsa.PublicKey{Curve: elliptic.P256(), X: new(big.Int).SetBytes(key.X[:]), Y: new(big.Int).SetBytes(key.Y[:])}
	if !ecdsa.Verify(&public, digest[:], new(big.Int).SetBytes(signature[:32]), new(big.Int).SetBytes(signature[32:])) {
		return rejected(ClassificationSignature, "es256-verification")
	}
	return nil
}

func sigStructure(protected, payload []byte) []byte {
	result := []byte{0x84, 0x6a}
	result = append(result, "Signature1"...)
	result = appendByteString(result, protected)
	result = append(result, 0x40)
	result = appendByteString(result, payload)
	return result
}

func appendByteString(destination, value []byte) []byte {
	length := len(value)
	switch {
	case length < 24:
		destination = append(destination, 0x40|byte(length))
	case length <= 0xff:
		destination = append(destination, 0x58, byte(length))
	case length <= 0xffff:
		destination = append(destination, 0x59, byte(length>>8), byte(length)) //nolint:gosec // This branch proves length is in 0..65535.
	default:
		panic("passive bootstrap fixture exceeds bounded byte-string encoder")
	}
	return append(destination, value...)
}
