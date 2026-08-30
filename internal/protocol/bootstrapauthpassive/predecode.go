package bootstrapauthpassive

import (
	"encoding/binary"
	"errors"

	"capsule.local/capsule/internal/protocol/cborscan"
)

type predecodeProfile struct {
	maxBytes         int
	maxDepth         int
	maxItems         int
	maxMapEntries    uint64
	maxArrayElements uint64
	allowedTag       uint64
	allowTag         bool
	allowFalse       bool
}

func (p predecodeProfile) sharedProfile() cborscan.Profile {
	return cborscan.Profile{
		MaxBytes:         p.maxBytes,
		MaxDepth:         p.maxDepth,
		MaxItems:         p.maxItems,
		MaxMapEntries:    p.maxMapEntries,
		MaxArrayElements: p.maxArrayElements,
		AllowedTag:       p.allowedTag,
		AllowTag:         p.allowTag,
		AllowFalse:       p.allowFalse,
	}
}

func predecode(value []byte, profile predecodeProfile) error {
	if err := cborscan.Predecode(value, profile.sharedProfile()); err != nil {
		return rejected(ClassificationMalformed, predecodeRefusalCode(err))
	}
	return nil
}

// predecodeRefusalCode maps a cborscan.Error's Reason to this package's
// original predecoder code vocabulary. Every predecode refusal this package
// observes is classified ClassificationMalformed, matching the predecoder
// this replaces.
func predecodeRefusalCode(err error) string {
	var scanErr *cborscan.Error
	if !errors.As(err, &scanErr) {
		return "cbor-predecode"
	}
	switch scanErr.Reason {
	case cborscan.ReasonEmptyPayload:
		return "empty-cbor"
	case cborscan.ReasonRawByteLimit:
		return "raw-byte-limit"
	case cborscan.ReasonTrailingData:
		return "trailing-cbor"
	case cborscan.ReasonDepthLimit:
		return "cbor-depth"
	case cborscan.ReasonItemLimit:
		return "cbor-items"
	case cborscan.ReasonTruncatedItem:
		return "truncated-cbor"
	case cborscan.ReasonUnsafeInteger:
		return "unsafe-integer"
	case cborscan.ReasonUnsafeNegativeInteger:
		return "unsafe-negative-integer"
	case cborscan.ReasonTruncatedString:
		return "truncated-string"
	case cborscan.ReasonDecodedStringLimit:
		return "decoded-string-limit"
	case cborscan.ReasonInvalidUTF8:
		return "invalid-utf8"
	case cborscan.ReasonArrayElementLimit:
		return "array-elements"
	case cborscan.ReasonMapEntryLimit:
		return "map-entries"
	case cborscan.ReasonDuplicateMapKeyOrder:
		return "duplicate-or-noncanonical-map-key"
	case cborscan.ReasonSemanticTag:
		return "semantic-tag"
	case cborscan.ReasonSimpleOrFloat:
		return "simple-or-float"
	case cborscan.ReasonMajorType:
		return "cbor-major-type"
	case cborscan.ReasonInvalidItem:
		return "cbor-item"
	case cborscan.ReasonIndefiniteOrReserved:
		return "indefinite-or-reserved"
	case cborscan.ReasonTruncatedArgument:
		return "truncated-argument"
	case cborscan.ReasonNonpreferredArgument:
		return "nonpreferred-argument"
	default:
		return "cbor-predecode"
	}
}

type envelopeParts struct {
	protectedStart int
	protectedEnd   int
	payloadStart   int
	payloadEnd     int
	signatureStart int
	signatureEnd   int
}

// frameEnvelope parses only the fixed Sign1 framing on caller-owned bytes.
// It performs the nested raw caps before the caller bytes are copied.
func frameEnvelope(received []byte, envelopeCap, protectedCap, payloadCap int) (envelopeParts, error) {
	if len(received) == 0 || len(received) > envelopeCap {
		return envelopeParts{}, rejected(ClassificationMalformed, "envelope-raw-byte-limit")
	}
	if len(received) < 2 || received[0] != 0xd2 || received[1] != 0x84 {
		return envelopeParts{}, rejected(ClassificationMalformed, "sign1-framing")
	}
	offset := 2
	protectedStart, protectedEnd, next, err := framedByteString(received, offset, protectedCap)
	if err != nil {
		return envelopeParts{}, err
	}
	offset = next
	if offset >= len(received) || received[offset] != 0xa0 {
		return envelopeParts{}, rejected(ClassificationMalformed, "unprotected-header")
	}
	offset++
	payloadStart, payloadEnd, next, err := framedByteString(received, offset, payloadCap)
	if err != nil {
		return envelopeParts{}, err
	}
	offset = next
	signatureStart, signatureEnd, next, err := framedByteString(received, offset, 64)
	if err != nil {
		return envelopeParts{}, err
	}
	if signatureEnd-signatureStart != 64 {
		return envelopeParts{}, rejected(ClassificationMalformed, "signature-shape")
	}
	if next != len(received) {
		return envelopeParts{}, rejected(ClassificationMalformed, "trailing-cose")
	}
	return envelopeParts{protectedStart, protectedEnd, payloadStart, payloadEnd, signatureStart, signatureEnd}, nil
}

func framedByteString(value []byte, offset, capBytes int) (int, int, int, error) {
	if offset >= len(value) {
		return 0, 0, 0, rejected(ClassificationMalformed, "truncated-byte-string")
	}
	initial := value[offset]
	if initial>>5 != 2 {
		return 0, 0, 0, rejected(ClassificationMalformed, "expected-byte-string")
	}
	offset++
	additional := initial & 0x1f
	var length uint64
	switch {
	case additional < 24:
		length = uint64(additional)
	case additional == 24:
		if offset >= len(value) {
			return 0, 0, 0, rejected(ClassificationMalformed, "truncated-byte-string-length")
		}
		length = uint64(value[offset])
		offset++
		if length < 24 {
			return 0, 0, 0, rejected(ClassificationMalformed, "nonpreferred-byte-string-length")
		}
	case additional == 25:
		if offset+2 > len(value) {
			return 0, 0, 0, rejected(ClassificationMalformed, "truncated-byte-string-length")
		}
		length = uint64(binary.BigEndian.Uint16(value[offset : offset+2]))
		offset += 2
		if length < 256 {
			return 0, 0, 0, rejected(ClassificationMalformed, "nonpreferred-byte-string-length")
		}
	default:
		return 0, 0, 0, rejected(ClassificationMalformed, "unsupported-byte-string-length")
	}
	if length > uint64(capBytes) { // #nosec G115 -- capBytes is a reviewed positive object-profile constant.
		return 0, 0, 0, rejected(ClassificationMalformed, "nested-raw-byte-limit")
	}
	if length > uint64(len(value)-offset) { // #nosec G115 -- the nonnegative remaining slice length widens losslessly.
		return 0, 0, 0, rejected(ClassificationMalformed, "truncated-byte-string")
	}
	end := offset + int(length)
	return offset, end, end, nil
}
