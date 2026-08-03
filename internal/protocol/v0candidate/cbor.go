package v0candidate

import (
	"bytes"
	"encoding/binary"
	"unicode/utf8"
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
)

// cborProfile bounds the allocation-independent predecoder for one object
// kind: total byte/depth/item ceilings plus per-container map/array caps.
type cborProfile struct {
	maxBytes         int
	maxDepth         int
	maxItems         int
	maxMapEntries    uint64
	maxArrayElements uint64
}

var executionPlanCBORProfile = cborProfile{
	maxBytes:         ExecutionPlanMaxCBORBytes,
	maxDepth:         ExecutionPlanMaxCBORDepth,
	maxItems:         ExecutionPlanMaxCBORItems,
	maxMapEntries:    ExecutionPlanMaxCBORMapEntries,
	maxArrayElements: ExecutionPlanMaxCBORArrayElements,
}

var planRegistrationCBORProfile = cborProfile{
	maxBytes:         PlanRegistrationMaxCBORBytes,
	maxDepth:         PlanRegistrationMaxCBORDepth,
	maxItems:         PlanRegistrationMaxCBORItems,
	maxMapEntries:    PlanRegistrationMaxCBORMapEntries,
	maxArrayElements: PlanRegistrationMaxCBORArrayElements,
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

type cborScanner struct {
	bytes            []byte
	offset           int
	items            int
	decodedStringLen uint64
	profile          cborProfile
}

func predecodeCBOR(received []byte, profile cborProfile) error {
	if len(received) == 0 {
		return malformed("empty CBOR payload")
	}
	if len(received) > profile.maxBytes {
		return malformed("CBOR raw-byte limit exceeded")
	}

	scanner := cborScanner{bytes: received, profile: profile}
	if err := scanner.scanItem(1); err != nil {
		return err
	}
	if scanner.offset != len(received) {
		return malformed("trailing CBOR data")
	}
	return nil
}

func (s *cborScanner) scanItem(depth int) error {
	if depth > s.profile.maxDepth {
		return malformed("CBOR nesting-depth limit exceeded")
	}
	s.items++
	if s.items > s.profile.maxItems {
		return malformed("CBOR data-item limit exceeded")
	}

	start := s.offset
	major, argument, err := s.readHead()
	if err != nil {
		return err
	}
	switch major {
	case 0:
		if argument > MaxSafeInteger {
			return malformed("CBOR unsigned integer exceeds UInt53")
		}
	case 1:
		if argument >= MaxSafeInteger {
			return malformed("CBOR negative integer exceeds the safe-integer range")
		}
	case 2, 3:
		if argument > uint64(len(s.bytes)-s.offset) {
			return malformed("truncated CBOR string")
		}
		if s.decodedStringLen > uint64(s.profile.maxBytes)-argument {
			return malformed("CBOR decoded-string limit exceeded")
		}
		s.decodedStringLen += argument
		end := s.offset + int(argument)
		if major == 3 && !utf8.Valid(s.bytes[s.offset:end]) {
			return malformed("invalid UTF-8 in CBOR text string")
		}
		s.offset = end
	case 4:
		if argument > s.profile.maxArrayElements {
			return malformed("CBOR array-element limit exceeded")
		}
		for index := uint64(0); index < argument; index++ {
			if err := s.scanItem(depth + 1); err != nil {
				return err
			}
		}
	case 5:
		if argument > s.profile.maxMapEntries {
			return malformed("CBOR map-entry limit exceeded")
		}
		var previousKey []byte
		for index := uint64(0); index < argument; index++ {
			keyStart := s.offset
			if err := s.scanItem(depth + 1); err != nil {
				return err
			}
			key := s.bytes[keyStart:s.offset]
			if previousKey != nil && compareDeterministicCBOR(previousKey, key) >= 0 {
				return malformed("duplicate or noncanonical CBOR map key order")
			}
			previousKey = key
			if err := s.scanItem(depth + 1); err != nil {
				return err
			}
		}
	case 6:
		return malformed("CBOR semantic tags are unsupported")
	case 7:
		return malformed("CBOR simple and floating-point values are unsupported")
	default:
		return malformed("unsupported CBOR major type")
	}
	if s.offset <= start {
		return malformed("invalid CBOR item")
	}
	return nil
}

func (s *cborScanner) readHead() (byte, uint64, error) {
	if s.offset >= len(s.bytes) {
		return 0, 0, malformed("truncated CBOR item")
	}
	initial := s.bytes[s.offset]
	s.offset++
	major := initial >> 5
	additional := initial & 0x1f
	if additional < 24 {
		return major, uint64(additional), nil
	}

	var width int
	switch additional {
	case 24:
		width = 1
	case 25:
		width = 2
	case 26:
		width = 4
	case 27:
		width = 8
	default:
		return 0, 0, malformed("indefinite or reserved CBOR argument")
	}
	if width > len(s.bytes)-s.offset {
		return 0, 0, malformed("truncated CBOR argument")
	}

	var argument uint64
	switch width {
	case 1:
		argument = uint64(s.bytes[s.offset])
	case 2:
		argument = uint64(binary.BigEndian.Uint16(s.bytes[s.offset : s.offset+width]))
	case 4:
		argument = uint64(binary.BigEndian.Uint32(s.bytes[s.offset : s.offset+width]))
	case 8:
		argument = binary.BigEndian.Uint64(s.bytes[s.offset : s.offset+width])
	}
	s.offset += width

	var minimum uint64
	switch width {
	case 1:
		minimum = 24
	case 2:
		minimum = 1 << 8
	case 4:
		minimum = 1 << 16
	case 8:
		minimum = 1 << 32
	}
	if argument < minimum {
		return 0, 0, malformed("nonpreferred CBOR integer or length encoding")
	}
	return major, argument, nil
}

// compareDeterministicCBOR orders two encoded CBOR map keys by RFC 8949 §4.2.1
// canonical/deterministic ordering: shorter encoding first, then byte-lexicographic.
func compareDeterministicCBOR(left, right []byte) int {
	if len(left) < len(right) {
		return -1
	}
	if len(left) > len(right) {
		return 1
	}
	return bytes.Compare(left, right)
}

// cborReader is used only after the allocation-independent predecoder has
// accepted the complete private copy.
type cborReader struct {
	bytes  []byte
	offset int
}

func (r *cborReader) head() (byte, uint64, error) {
	scanner := cborScanner{bytes: r.bytes, offset: r.offset}
	major, argument, err := scanner.readHead()
	if err != nil {
		return 0, 0, err
	}
	r.offset = scanner.offset
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
