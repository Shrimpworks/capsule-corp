package bootstrapauthpassive

import (
	"bytes"
	"encoding/binary"
	"unicode/utf8"
)

const maxSafeInteger = uint64(9007199254740991)

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

func predecode(value []byte, profile predecodeProfile) error {
	if len(value) == 0 {
		return rejected(ClassificationMalformed, "empty-cbor")
	}
	if len(value) > profile.maxBytes {
		return rejected(ClassificationMalformed, "raw-byte-limit")
	}
	s := cborScanner{value: value, profile: profile}
	if err := s.item(1); err != nil {
		return err
	}
	if s.offset != len(value) {
		return rejected(ClassificationMalformed, "trailing-cbor")
	}
	return nil
}

type cborScanner struct {
	value   []byte
	offset  int
	items   int
	strings uint64
	profile predecodeProfile
}

func (s *cborScanner) item(depth int) error {
	if depth > s.profile.maxDepth {
		return rejected(ClassificationMalformed, "cbor-depth")
	}
	s.items++
	if s.items > s.profile.maxItems {
		return rejected(ClassificationMalformed, "cbor-items")
	}
	major, argument, err := s.head()
	if err != nil {
		return err
	}
	switch major {
	case 0:
		if argument > maxSafeInteger {
			return rejected(ClassificationMalformed, "unsafe-integer")
		}
	case 1:
		if argument >= maxSafeInteger {
			return rejected(ClassificationMalformed, "unsafe-negative-integer")
		}
	case 2, 3:
		if argument > uint64(len(s.value)-s.offset) { // #nosec G115 -- the nonnegative remaining slice length widens losslessly.
			return rejected(ClassificationMalformed, "truncated-string")
		}
		if s.strings > uint64(s.profile.maxBytes)-argument { // #nosec G115 -- every profile maximum is a reviewed positive constant.
			return rejected(ClassificationMalformed, "decoded-string-limit")
		}
		s.strings += argument
		end := s.offset + int(argument) // #nosec G115 -- the preceding comparison proves argument fits the remaining in-memory slice.
		if major == 3 && !utf8.Valid(s.value[s.offset:end]) {
			return rejected(ClassificationMalformed, "invalid-utf8")
		}
		s.offset = end
	case 4:
		if argument > s.profile.maxArrayElements {
			return rejected(ClassificationMalformed, "array-elements")
		}
		for i := uint64(0); i < argument; i++ {
			if err := s.item(depth + 1); err != nil {
				return err
			}
		}
	case 5:
		if argument > s.profile.maxMapEntries {
			return rejected(ClassificationMalformed, "map-entries")
		}
		var previous []byte
		for i := uint64(0); i < argument; i++ {
			start := s.offset
			if err := s.item(depth + 1); err != nil {
				return err
			}
			key := s.value[start:s.offset]
			if previous != nil && deterministicCompare(previous, key) >= 0 {
				return rejected(ClassificationMalformed, "duplicate-or-noncanonical-map-key")
			}
			previous = key
			if err := s.item(depth + 1); err != nil {
				return err
			}
		}
	case 6:
		if !s.profile.allowTag || argument != s.profile.allowedTag {
			return rejected(ClassificationMalformed, "semantic-tag")
		}
		if err := s.item(depth + 1); err != nil {
			return err
		}
	case 7:
		if !s.profile.allowFalse || argument != 20 {
			return rejected(ClassificationMalformed, "simple-or-float")
		}
	default:
		return rejected(ClassificationMalformed, "cbor-major-type")
	}
	return nil
}

func (s *cborScanner) head() (byte, uint64, error) {
	if s.offset >= len(s.value) {
		return 0, 0, rejected(ClassificationMalformed, "truncated-cbor")
	}
	initial := s.value[s.offset]
	s.offset++
	major, additional := initial>>5, initial&0x1f
	if additional < 24 {
		return major, uint64(additional), nil
	}
	width := 0
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
		return 0, 0, rejected(ClassificationMalformed, "indefinite-or-reserved")
	}
	if width > len(s.value)-s.offset {
		return 0, 0, rejected(ClassificationMalformed, "truncated-argument")
	}
	var argument uint64
	switch width {
	case 1:
		argument = uint64(s.value[s.offset])
	case 2:
		argument = uint64(binary.BigEndian.Uint16(s.value[s.offset : s.offset+2]))
	case 4:
		argument = uint64(binary.BigEndian.Uint32(s.value[s.offset : s.offset+4]))
	case 8:
		argument = binary.BigEndian.Uint64(s.value[s.offset : s.offset+8])
	}
	s.offset += width
	minimum := [...]uint64{0, 24, 1 << 8, 0, 1 << 16, 0, 0, 0, 1 << 32}[width]
	if argument < minimum {
		return 0, 0, rejected(ClassificationMalformed, "nonpreferred-argument")
	}
	return major, argument, nil
}

func deterministicCompare(left, right []byte) int {
	if len(left) < len(right) {
		return -1
	}
	if len(left) > len(right) {
		return 1
	}
	return bytes.Compare(left, right)
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
