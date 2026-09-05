package v0candidate

import (
	"errors"

	"capsule.local/capsule/internal/protocol/cborscan"
)

const (
	ExecutionPlanMaxCBORBytes            = 65_536
	ExecutionPlanMaxCBORDepth            = 8
	ExecutionPlanMaxCBORItems            = 256
	ExecutionPlanMaxCBORMapEntries       = 64
	ExecutionPlanMaxCBORArrayElements    = 8
	PlanRegistrationMaxCBORBytes         = 4_096
	PlanRegistrationMaxCBORDepth         = 4
	PlanRegistrationMaxCBORItems         = 33
	PlanRegistrationMaxCBORMapEntries    = 16
	PlanRegistrationMaxCBORArrayElements = 0
	SourceManifestMinCBORBytes           = 87
	SourceManifestMaxCBORBytes           = 95
	SourceManifestMaxCBORDepth           = 4
	SourceManifestMaxCBORItems           = 15
	SourceManifestMaxCBORMapEntries      = 5
	SourceManifestMaxCBORArrayElements   = 3
)

// cborProfile bounds the allocation-independent predecoder for one object
// kind: total byte/depth/item ceilings plus per-container map/array caps. It
// is the shared cborscan.Profile — this package never allows a semantic tag
// or a simple/float value, so beyond the caps it only opts into arrays, which
// cborscan refuses unless a profile asks for them.
type cborProfile = cborscan.Profile

var executionPlanCBORProfile = cborProfile{
	MaxBytes:         ExecutionPlanMaxCBORBytes,
	MaxDepth:         ExecutionPlanMaxCBORDepth,
	MaxItems:         ExecutionPlanMaxCBORItems,
	MaxMapEntries:    ExecutionPlanMaxCBORMapEntries,
	MaxArrayElements: ExecutionPlanMaxCBORArrayElements,
	AllowArray:       true,
}

var planRegistrationCBORProfile = cborProfile{
	MaxBytes:         PlanRegistrationMaxCBORBytes,
	MaxDepth:         PlanRegistrationMaxCBORDepth,
	MaxItems:         PlanRegistrationMaxCBORItems,
	MaxMapEntries:    PlanRegistrationMaxCBORMapEntries,
	MaxArrayElements: PlanRegistrationMaxCBORArrayElements,
	AllowArray:       true,
}

var sourceManifestCBORProfile = cborProfile{
	MaxBytes:         SourceManifestMaxCBORBytes,
	MaxDepth:         SourceManifestMaxCBORDepth,
	MaxItems:         SourceManifestMaxCBORItems,
	MaxMapEntries:    SourceManifestMaxCBORMapEntries,
	MaxArrayElements: SourceManifestMaxCBORArrayElements,
	AllowArray:       true,
}

// PredecodeExecutionPlanCBOR applies the allocation-independent deterministic
// CBOR profile for an ExecutionPlan. It does not imply that the bytes have the
// ExecutionPlan object shape.
func PredecodeExecutionPlanCBOR(received []byte) error {
	return predecodeCBOR(received, executionPlanCBORProfile)
}

// PredecodePlanRegistrationCBOR applies the allocation-independent
// deterministic CBOR profile for a PlanRegistration. It does not imply that
// the bytes have the PlanRegistration object shape.
func PredecodePlanRegistrationCBOR(received []byte) error {
	return predecodeCBOR(received, planRegistrationCBORProfile)
}

// PredecodeSourceManifestCBOR applies the exact passive v0 manifest profile.
func PredecodeSourceManifestCBOR(received []byte) error {
	if len(received) < SourceManifestMinCBORBytes {
		return malformed("SourceManifest is shorter than its minimum canonical encoding")
	}
	return predecodeCBOR(received, sourceManifestCBORProfile)
}

func predecodeCBOR(received []byte, profile cborProfile) error {
	if err := cborscan.Predecode(received, profile); err != nil {
		return malformed(predecodeRefusalMessage(err))
	}
	return nil
}

// predecodeRefusalMessage maps a cborscan.Error's Reason to this package's
// original predecoder wording. Every predecode refusal is classified
// ClassificationMalformed here, matching the predecoder this replaces: this
// package never allows a semantic tag or a simple/float value, so it can
// only ever observe the Reason values below.
func predecodeRefusalMessage(err error) string {
	var scanErr *cborscan.Error
	if !errors.As(err, &scanErr) {
		return "malformed CBOR payload"
	}
	switch scanErr.Reason {
	case cborscan.ReasonEmptyPayload:
		return "empty CBOR payload"
	case cborscan.ReasonRawByteLimit:
		return "CBOR raw-byte limit exceeded"
	case cborscan.ReasonTrailingData:
		return "trailing CBOR data"
	case cborscan.ReasonDepthLimit:
		return "CBOR nesting-depth limit exceeded"
	case cborscan.ReasonItemLimit:
		return "CBOR data-item limit exceeded"
	case cborscan.ReasonTruncatedItem:
		return "truncated CBOR item"
	case cborscan.ReasonUnsafeInteger:
		return "CBOR unsigned integer exceeds UInt53"
	case cborscan.ReasonUnsafeNegativeInteger:
		return "CBOR negative integer exceeds the safe-integer range"
	case cborscan.ReasonTruncatedString:
		return "truncated CBOR string"
	case cborscan.ReasonDecodedStringLimit:
		return "CBOR decoded-string limit exceeded"
	case cborscan.ReasonInvalidUTF8:
		return "invalid UTF-8 in CBOR text string"
	case cborscan.ReasonArrayElementLimit:
		return "CBOR array-element limit exceeded"
	case cborscan.ReasonMapEntryLimit:
		return "CBOR map-entry limit exceeded"
	case cborscan.ReasonDuplicateMapKeyOrder:
		return "duplicate or noncanonical CBOR map key order"
	case cborscan.ReasonSemanticTag:
		return "CBOR semantic tags are unsupported"
	case cborscan.ReasonSimpleOrFloat:
		return "CBOR simple and floating-point values are unsupported"
	case cborscan.ReasonMajorType:
		return "unsupported CBOR major type"
	case cborscan.ReasonInvalidItem:
		return "invalid CBOR item"
	case cborscan.ReasonIndefiniteOrReserved:
		return "indefinite or reserved CBOR argument"
	case cborscan.ReasonTruncatedArgument:
		return "truncated CBOR argument"
	case cborscan.ReasonNonpreferredArgument:
		return "nonpreferred CBOR integer or length encoding"
	default:
		return "malformed CBOR payload"
	}
}

// cborReader is used only after the allocation-independent predecoder has
// accepted the complete private copy.
type cborReader struct {
	bytes  []byte
	offset int
}

func (r *cborReader) head() (byte, uint64, error) {
	major, argument, next, err := cborscan.ReadHead(r.bytes, r.offset)
	if err != nil {
		return 0, 0, err
	}
	r.offset = next
	return major, argument, nil
}

func (r *cborReader) unsigned() (uint64, error) {
	major, value, err := r.head()
	if err != nil {
		return 0, schema("invalid unsigned integer field")
	}
	if major != 0 || value > MaxSafeInteger {
		return 0, schema("field must be UInt53")
	}
	return value, nil
}

func (r *cborReader) byteString() ([]byte, error) {
	major, length, err := r.head()
	if err != nil || major != 2 || length > uint64(len(r.bytes)-r.offset) {
		return nil, schema("field must be a byte string")
	}
	value := r.bytes[r.offset : r.offset+int(length)]
	r.offset += int(length)
	return value, nil
}

func (r *cborReader) textString() (string, error) {
	major, length, err := r.head()
	if err != nil || major != 3 || length > uint64(len(r.bytes)-r.offset) {
		return "", schema("field must be a text string")
	}
	value := string(r.bytes[r.offset : r.offset+int(length)])
	r.offset += int(length)
	return value, nil
}

func (r *cborReader) arrayLength() (uint64, error) {
	major, length, err := r.head()
	if err != nil || major != 4 {
		return 0, schema("field must be an array")
	}
	return length, nil
}

func (r *cborReader) mapLength() (uint64, error) {
	major, length, err := r.head()
	if err != nil || major != 5 {
		return 0, schema("object must be a map")
	}
	return length, nil
}

func (r *cborReader) requiredKey(expected uint64) error {
	actual, err := r.unsigned()
	if err != nil {
		return schema("object field label must be an unsigned integer")
	}
	if actual == expected {
		return nil
	}
	if actual > expected && actual <= 24 {
		return schema("required object field is missing")
	}
	return unsupported("unknown object field")
}
