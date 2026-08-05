package v0candidate

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"unicode/utf8"
)

const MJSMainSourceMaxBytes = 262_144

// MJSMainSource retains an exact private copy of strict UTF-8 main.mjs bytes.
// It does not parse or execute the source.
type MJSMainSource struct {
	bytes  []byte
	digest SourceContentDigest
}

func NewMJSMainSource(received []byte) (*MJSMainSource, error) {
	if len(received) > MJSMainSourceMaxBytes {
		return nil, semantic("main.mjs exceeds 262144 bytes")
	}
	exact := bytes.Clone(received)
	if bytes.HasPrefix(exact, []byte{0xef, 0xbb, 0xbf}) {
		return nil, malformed("main.mjs must not begin with a UTF-8 BOM")
	}
	if !utf8.Valid(exact) {
		return nil, malformed("main.mjs must contain strict UTF-8")
	}
	digest := sha256.Sum256(exact)
	return &MJSMainSource{bytes: exact, digest: SourceContentDigest(digest)}, nil
}

func (s *MJSMainSource) AuthoritativeBytes() []byte {
	return bytes.Clone(s.bytes)
}

func (s *MJSMainSource) Digest() SourceContentDigest {
	return s.digest
}

func (s *MJSMainSource) ByteLength() UInt53 {
	return UInt53(len(s.bytes))
}

// DecodedSourceManifest retains defensive exact bytes, digest, decoded view,
// and a private copy of the separately supplied source bytes.
type DecodedSourceManifest struct {
	authoritativeBytes []byte
	digest             SourceManifestDigest
	view               SourceManifest
	source             *MJSMainSource
}

func (m *DecodedSourceManifest) AuthoritativeBytes() []byte {
	return bytes.Clone(m.authoritativeBytes)
}

func (m *DecodedSourceManifest) Digest() SourceManifestDigest {
	return m.digest
}

func (m *DecodedSourceManifest) View() SourceManifest {
	view := m.view
	view.Members = append([]SourceManifestMember(nil), m.view.Members...)
	return view
}

func (m *DecodedSourceManifest) SourceBytes() []byte {
	return m.source.AuthoritativeBytes()
}

func ValidateMJSSourceProfile(value string) error {
	if value != MJSSourceProfile {
		return unsupported("unsupported source profile")
	}
	return nil
}

func ValidateMJSSourceMediaType(value string) error {
	if value != MJSSourceMediaType {
		return unsupported("unsupported source-member media type")
	}
	return nil
}

func ValidateSourceManifestMediaType(value string) error {
	if value != SourceManifestMediaType {
		return unsupported("unsupported SourceManifest media type")
	}
	return nil
}

// EncodeSourceManifest returns the deterministic-CBOR SourceManifest v0 for
// already validated exact source bytes. It does not parse or execute source.
func EncodeSourceManifest(source *MJSMainSource) ([]byte, error) {
	if source == nil {
		return nil, schema("source is required")
	}
	result := []byte{0xa5, 0x01, 0x77}
	result = append(result, []byte(SourceManifestObjectType)...)
	result = append(result, 0x02, 0x00, 0x03, 0x68)
	result = append(result, []byte(MJSMainPath)...)
	result = append(result, 0x04, 0x81, 0x83, 0x68)
	result = append(result, []byte(MJSMainPath)...)
	result = append(result, 0x58, 0x20)
	digest := source.Digest()
	result = append(result, digest[:]...)
	result = appendCBORUnsigned(result, uint64(source.ByteLength()))
	result = append(result, 0x05)
	result = appendCBORUnsigned(result, uint64(source.ByteLength()))
	return result, nil
}

func appendCBORUnsigned(destination []byte, value uint64) []byte {
	switch {
	case value < 24:
		return append(destination, byte(value))
	case value <= 0xff:
		return append(destination, 0x18, byte(value))
	case value <= 0xffff:
		return binary.BigEndian.AppendUint16(append(destination, 0x19), uint16(value))
	case value <= 0xffffffff:
		return binary.BigEndian.AppendUint32(append(destination, 0x1a), uint32(value))
	default:
		return binary.BigEndian.AppendUint64(append(destination, 0x1b), value)
	}
}

// DecodeSourceManifest copies, predecodes, strictly decodes, hashes, and
// independently binds the manifest to separately supplied exact main.mjs bytes.
func DecodeSourceManifest(received []byte, mediaType string, sourceBytes []byte) (*DecodedSourceManifest, error) {
	if err := ValidateSourceManifestMediaType(mediaType); err != nil {
		return nil, err
	}
	source, err := NewMJSMainSource(sourceBytes)
	if err != nil {
		return nil, domain("SourceManifest source bytes do not satisfy the exact main.mjs byte contract")
	}
	// Reject an out-of-bounds length before cloning so an oversized received
	// slice cannot force a full-size allocation ahead of rejection.
	if len(received) < SourceManifestMinCBORBytes {
		return nil, malformed("SourceManifest is shorter than its minimum canonical encoding")
	}
	if len(received) > SourceManifestMaxCBORBytes {
		return nil, malformed("CBOR raw-byte limit exceeded")
	}
	authoritative := bytes.Clone(received)
	if err := PredecodeSourceManifestCBOR(authoritative); err != nil {
		return nil, err
	}

	reader := cborReader{bytes: authoritative}
	entries, err := reader.mapLength()
	if err != nil {
		return nil, err
	}
	if entries < 5 {
		return nil, schema("SourceManifest is missing required fields")
	}
	if entries > 5 {
		return nil, unsupported("SourceManifest contains unknown fields")
	}
	if err := reader.requiredKey(1); err != nil {
		return nil, err
	}
	objectType, err := reader.textString()
	if err != nil {
		return nil, err
	}
	if objectType != SourceManifestObjectType {
		return nil, unsupported("unknown SourceManifest object type")
	}
	if err := reader.requiredKey(2); err != nil {
		return nil, err
	}
	version, err := reader.unsigned()
	if err != nil {
		return nil, err
	}
	if version != uint64(CandidateObjectVersion) {
		return nil, unsupported("unsupported SourceManifest object version")
	}
	if err := reader.requiredKey(3); err != nil {
		return nil, err
	}
	entrypoint, err := reader.textString()
	if err != nil {
		return nil, err
	}
	if entrypoint != MJSMainPath {
		return nil, domain("SourceManifest entrypoint must be exactly main.mjs")
	}
	if err := reader.requiredKey(4); err != nil {
		return nil, err
	}
	members, err := reader.arrayLength()
	if err != nil {
		return nil, err
	}
	if members != 1 {
		return nil, schema("SourceManifest must contain exactly one member")
	}
	memberFields, err := reader.arrayLength()
	if err != nil {
		return nil, err
	}
	if memberFields != 3 {
		return nil, schema("SourceManifest member must contain exactly three fields")
	}
	logicalPath, err := reader.textString()
	if err != nil {
		return nil, err
	}
	if logicalPath != MJSMainPath {
		return nil, domain("SourceManifest member path must be exactly main.mjs")
	}
	contentDigestBytes, err := reader.byteString()
	if err != nil {
		return nil, err
	}
	contentDigest, err := NewSourceContentDigest(contentDigestBytes)
	if err != nil {
		return nil, err
	}
	memberLengthValue, err := reader.unsigned()
	if err != nil {
		return nil, err
	}
	memberLength, err := NewUInt53(memberLengthValue)
	if err != nil {
		return nil, err
	}
	if err := reader.requiredKey(5); err != nil {
		return nil, err
	}
	aggregateValue, err := reader.unsigned()
	if err != nil {
		return nil, err
	}
	aggregateLength, err := NewUInt53(aggregateValue)
	if err != nil {
		return nil, err
	}
	if reader.offset != len(authoritative) {
		return nil, malformed("trailing SourceManifest data")
	}
	if contentDigest != source.Digest() {
		return nil, domain("SourceManifest content digest does not match exact main.mjs bytes")
	}
	if memberLength != source.ByteLength() || aggregateLength != source.ByteLength() {
		return nil, domain("SourceManifest lengths do not match exact main.mjs bytes")
	}

	digest := sha256.Sum256(authoritative)
	return &DecodedSourceManifest{
		authoritativeBytes: authoritative,
		digest:             SourceManifestDigest(digest),
		view: SourceManifest{
			ObjectType:    objectType,
			ObjectVersion: ObjectVersion(version),
			Entrypoint:    MJSMainEntrypoint(entrypoint),
			Members: []SourceManifestMember{{
				LogicalPath:   MJSMainLogicalPath(logicalPath),
				ContentDigest: contentDigest,
				ByteLength:    memberLength,
			}},
			AggregateByteLength: aggregateLength,
		},
		source: source,
	}, nil
}
