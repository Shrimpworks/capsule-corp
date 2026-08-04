package v0cbor

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"

	"capsule.local/capsule/internal/protocol/v0candidate"
	"github.com/fxamacker/cbor/v2"
)

// FailureClassification is the internal first-owning-boundary classification.
// It is not a public protocol error code.
type FailureClassification string

const (
	FailureMalformed   FailureClassification = "MALFORMED"
	FailureUnsupported FailureClassification = "UNSUPPORTED"
	FailureSchema      FailureClassification = "SCHEMA"
	FailureSemantic    FailureClassification = "SEMANTIC"
	FailureDomain      FailureClassification = "DOMAIN"
)

// DecodeError reports a bounded classification without echoing received bytes
// or decoded attacker-controlled values.
type DecodeError struct {
	classification FailureClassification
	detail         string
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("%s: %s", e.classification, e.detail)
}

// Classification returns the first owning conformance classification.
func (e *DecodeError) Classification() FailureClassification {
	return e.classification
}

// ErrorClassification extracts classifications from this wrapper and from the
// retained handwritten Capsule predecoder/source validator.
func ErrorClassification(err error) (FailureClassification, bool) {
	var decodeError *DecodeError
	if errors.As(err, &decodeError) {
		return decodeError.classification, true
	}
	classification, ok := v0candidate.ErrorClassification(err)
	return FailureClassification(classification), ok
}

func classified(classification FailureClassification, detail string) error {
	return &DecodeError{classification: classification, detail: detail}
}

func malformed(detail string) error {
	return classified(FailureMalformed, detail)
}

func unsupported(detail string) error {
	return classified(FailureUnsupported, detail)
}

func schema(detail string) error {
	return classified(FailureSchema, detail)
}

func domain(detail string) error {
	return classified(FailureDomain, detail)
}

type sourceManifestMemberWire struct {
	_          struct{} `cbor:",toarray"`
	Path       string
	Digest     []byte
	ByteLength uint64
}

type sourceManifestWire struct {
	ObjectType          string                     `cbor:"1,keyasint"`
	ObjectVersion       uint64                     `cbor:"2,keyasint"`
	Entrypoint          string                     `cbor:"3,keyasint"`
	Members             []sourceManifestMemberWire `cbor:"4,keyasint"`
	AggregateByteLength uint64                     `cbor:"5,keyasint"`
}

// SourceManifestMember is the closed, typed projection of the sole main.mjs
// member.  Arrays make the digest immutable by value.
type SourceManifestMember struct {
	LogicalPath   v0candidate.MJSMainLogicalPath
	ContentDigest v0candidate.SourceContentDigest
	ByteLength    v0candidate.UInt53
}

// SourceManifestView is the closed typed projection returned by the wrapper.
type SourceManifestView struct {
	ObjectType          string
	ObjectVersion       v0candidate.ObjectVersion
	Entrypoint          v0candidate.MJSMainEntrypoint
	Member              SourceManifestMember
	AggregateByteLength v0candidate.UInt53
}

// DecodedSourceManifest owns exact canonical input bytes, separately supplied
// source bytes, and their typed projection.  All byte accessors return copies.
type DecodedSourceManifest struct {
	exactBytes  []byte
	digest      v0candidate.SourceManifestDigest
	sourceBytes []byte
	view        SourceManifestView
}

// ExactBytes returns a defensive copy of the exact received canonical bytes.
func (m *DecodedSourceManifest) ExactBytes() []byte {
	return bytes.Clone(m.exactBytes)
}

// Digest returns SHA-256 over ExactBytes.
func (m *DecodedSourceManifest) Digest() v0candidate.SourceManifestDigest {
	return m.digest
}

// SourceBytes returns a defensive copy of the separately bound main.mjs bytes.
func (m *DecodedSourceManifest) SourceBytes() []byte {
	return bytes.Clone(m.sourceBytes)
}

// View returns the immutable typed projection.
func (m *DecodedSourceManifest) View() SourceManifestView {
	return m.view
}

// SourceManifestCodec is the only fxamacker-backed surface in this package.
// Its immutable modes are safe for concurrent use.
type SourceManifestCodec struct {
	encode cbor.EncMode
	decode cbor.DecMode
}

// NewSourceManifestCodec constructs the exact deterministic encoder and typed
// decoder modes.  No mode or generic marshal/unmarshal method is exposed.
func NewSourceManifestCodec() (*SourceManifestCodec, error) {
	encodeOptions := cbor.CoreDetEncOptions()
	encodeOptions.TagsMd = cbor.TagsForbidden
	encodeOptions.BinaryMarshaler = cbor.BinaryMarshalerNone
	encodeOptions.TextMarshaler = cbor.TextMarshalerNone
	encode, err := encodeOptions.EncMode()
	if err != nil {
		return nil, fmt.Errorf("construct SourceManifest deterministic encoder: %w", err)
	}

	// fxamacker's configurable collection minima are 16.  Capsule's scanner
	// enforces the tighter 3-element/5-pair profile before this decoder runs.
	decode, err := (cbor.DecOptions{
		DupMapKey:         cbor.DupMapKeyEnforcedAPF,
		MaxNestedLevels:   v0candidate.SourceManifestMaxCBORDepth,
		MaxArrayElements:  16,
		MaxMapPairs:       16,
		IndefLength:       cbor.IndefLengthForbidden,
		TagsMd:            cbor.TagsForbidden,
		UTF8:              cbor.UTF8RejectInvalid,
		ExtraReturnErrors: cbor.ExtraDecErrorUnknownField,
		BignumTag:         cbor.BignumTagForbidden,
		BinaryUnmarshaler: cbor.BinaryUnmarshalerNone,
		TextUnmarshaler:   cbor.TextUnmarshalerNone,
	}).DecMode()
	if err != nil {
		return nil, fmt.Errorf("construct SourceManifest typed decoder: %w", err)
	}

	return &SourceManifestCodec{encode: encode, decode: decode}, nil
}

// Encode constructs only the frozen SourceManifest v0 shape from exact
// main.mjs bytes.  Callers cannot provide fields, tags, maps, or options.
func (c *SourceManifestCodec) Encode(sourceBytes []byte) ([]byte, error) {
	if c == nil {
		return nil, errors.New("SourceManifest codec is required")
	}
	source, err := v0candidate.NewMJSMainSource(sourceBytes)
	if err != nil {
		return nil, err
	}
	digest := source.Digest()
	byteLength := uint64(source.ByteLength())
	wire := sourceManifestWire{
		ObjectType:    v0candidate.SourceManifestObjectType,
		ObjectVersion: uint64(v0candidate.CandidateObjectVersion),
		Entrypoint:    v0candidate.MJSMainPath,
		Members: []sourceManifestMemberWire{{
			Path:       v0candidate.MJSMainPath,
			Digest:     bytes.Clone(digest[:]),
			ByteLength: byteLength,
		}},
		AggregateByteLength: byteLength,
	}
	encoded, err := c.encode.Marshal(wire)
	if err != nil {
		return nil, fmt.Errorf("encode frozen SourceManifest: %w", err)
	}
	if err := v0candidate.PredecodeSourceManifestCBOR(encoded); err != nil {
		return nil, fmt.Errorf("validate generated SourceManifest: %w", err)
	}
	return encoded, nil
}

// Decode applies media, raw-byte, depth/item/collection, and deterministic
// predecode checks before copying.  fxamacker then performs typed field decode;
// Capsule re-encodes and compares the exact input, validates the closed shape,
// binds the separately supplied source, and only then returns owned bytes.
func (c *SourceManifestCodec) Decode(received []byte, mediaType string, sourceBytes []byte) (*DecodedSourceManifest, error) {
	if c == nil {
		return nil, errors.New("SourceManifest codec is required")
	}
	if err := v0candidate.ValidateSourceManifestMediaType(mediaType); err != nil {
		return nil, err
	}
	if err := v0candidate.PredecodeSourceManifestCBOR(received); err != nil {
		return nil, err
	}
	source, err := v0candidate.NewMJSMainSource(sourceBytes)
	if err != nil {
		return nil, domain("source bytes do not satisfy the frozen main.mjs contract")
	}

	exact := bytes.Clone(received)
	var wire sourceManifestWire
	if err := c.decode.Unmarshal(exact, &wire); err != nil {
		var unknownField *cbor.UnknownFieldError
		if errors.As(err, &unknownField) {
			return nil, unsupported("SourceManifest contains an unknown field")
		}
		return nil, schema("SourceManifest typed field decoding failed")
	}
	if err := validateShape(wire); err != nil {
		return nil, err
	}
	canonical, err := c.encode.Marshal(wire)
	if err != nil {
		return nil, fmt.Errorf("re-encode typed SourceManifest: %w", err)
	}
	if !bytes.Equal(canonical, exact) {
		return nil, malformed("SourceManifest is not canonical on wire")
	}
	if err := bindSource(wire, source); err != nil {
		return nil, err
	}

	member := wire.Members[0]
	var contentDigest v0candidate.SourceContentDigest
	copy(contentDigest[:], member.Digest)
	digest := sha256.Sum256(exact)
	return &DecodedSourceManifest{
		exactBytes:  exact,
		digest:      v0candidate.SourceManifestDigest(digest),
		sourceBytes: source.AuthoritativeBytes(),
		view: SourceManifestView{
			ObjectType:    wire.ObjectType,
			ObjectVersion: v0candidate.ObjectVersion(wire.ObjectVersion),
			Entrypoint:    v0candidate.MJSMainEntrypoint(wire.Entrypoint),
			Member: SourceManifestMember{
				LogicalPath:   v0candidate.MJSMainLogicalPath(member.Path),
				ContentDigest: contentDigest,
				ByteLength:    v0candidate.UInt53(member.ByteLength),
			},
			AggregateByteLength: v0candidate.UInt53(wire.AggregateByteLength),
		},
	}, nil
}

func validateShape(wire sourceManifestWire) error {
	if wire.ObjectType != v0candidate.SourceManifestObjectType {
		return unsupported("unknown SourceManifest object type")
	}
	if wire.ObjectVersion != uint64(v0candidate.CandidateObjectVersion) {
		return unsupported("unsupported SourceManifest object version")
	}
	if wire.Entrypoint != v0candidate.MJSMainPath {
		return domain("SourceManifest entrypoint is outside the frozen profile")
	}
	if len(wire.Members) != 1 {
		return schema("SourceManifest must contain exactly one member")
	}
	member := wire.Members[0]
	if member.Path != v0candidate.MJSMainPath {
		return domain("SourceManifest member path is outside the frozen profile")
	}
	if len(member.Digest) != sha256.Size {
		return schema("SourceManifest content digest width is invalid")
	}
	if member.ByteLength > v0candidate.MJSMainSourceMaxBytes ||
		wire.AggregateByteLength > v0candidate.MJSMainSourceMaxBytes {
		return domain("SourceManifest byte length exceeds the frozen source profile")
	}
	if member.ByteLength != wire.AggregateByteLength {
		return domain("SourceManifest member and aggregate lengths disagree")
	}
	return nil
}

func bindSource(wire sourceManifestWire, source *v0candidate.MJSMainSource) error {
	member := wire.Members[0]
	digest := source.Digest()
	if !bytes.Equal(member.Digest, digest[:]) ||
		member.ByteLength != uint64(source.ByteLength()) ||
		wire.AggregateByteLength != uint64(source.ByteLength()) {
		return domain("SourceManifest does not bind the supplied exact main.mjs bytes")
	}
	return nil
}
