// Package cborscan is the shared allocation-independent RFC-8949 deterministic
// CBOR predecoder backing internal/protocol/v0candidate and
// internal/protocol/bootstrapauthpassive. It performs one bounded,
// non-allocating pass over caller-owned bytes — depth/item/map/array caps,
// UTF-8 validation, canonical map-key ordering, non-preferred-length
// rejection, and safe-integer bounds — and reports only a caller-opaque
// Reason on refusal. Each caller package owns its own error type, message
// text, and classification vocabulary; this package never formats a message
// or assigns a classification itself, so consolidating the scan logic here
// changes no caller's observable refusal text or accept/reject decision.
package cborscan

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unicode/utf8"
)

// maxSafeInteger is the largest integer this predecoder accepts for a
// major-0/1 (unsigned/negative) integer head, matching the IEEE-754 double
// safe-integer bound every caller profile relies on.
const maxSafeInteger = uint64(9_007_199_254_740_991)

// Profile bounds one predecode pass: total byte/depth/item ceilings, the
// per-container map/array caps, and the optional major-4/6/7 extensions a
// caller profile may opt into. Every extension defaults closed, so the zero
// value rejects every array, every tag, and every simple/float value; an
// object set carrying any of them must say so explicitly.
type Profile struct {
	MaxBytes         int
	MaxDepth         int
	MaxItems         int
	MaxMapEntries    uint64
	MaxArrayElements uint64

	// AllowTag, when set, accepts exactly one semantic tag (major type 6)
	// whose argument equals AllowedTag, recursing into its tagged content.
	// Every other major-6 item is refused.
	AllowTag   bool
	AllowedTag uint64

	// AllowFalse, when set, accepts exactly the simple value `false`
	// (major type 7, argument 20). Every other major-7 item is refused.
	AllowFalse bool

	// AllowArray, when set, accepts an array (major type 4) whose element
	// count is within MaxArrayElements. Every array is otherwise refused
	// with ReasonMajorType, including an empty one — MaxArrayElements: 0
	// admits a zero-length array and so cannot express "no arrays at all".
	// Object sets that carry no array must leave this unset rather than
	// relying on the element cap.
	AllowArray bool
}

// Reason identifies why Predecode (or ReadHead) refused the input. It
// carries no message text: each caller package maps the Reason values it can
// observe to its own wording and classification.
type Reason int

// Reason values enumerate every refusal Predecode/ReadHead can produce; see
// the Reason type's own comment for how callers are expected to use them.
const (
	ReasonEmptyPayload Reason = iota
	ReasonRawByteLimit
	ReasonTrailingData
	ReasonDepthLimit
	ReasonItemLimit
	ReasonTruncatedItem
	ReasonUnsafeInteger
	ReasonUnsafeNegativeInteger
	ReasonTruncatedString
	ReasonDecodedStringLimit
	ReasonInvalidUTF8
	ReasonArrayElementLimit
	ReasonMapEntryLimit
	ReasonDuplicateMapKeyOrder
	ReasonSemanticTag
	ReasonSimpleOrFloat
	ReasonMajorType
	ReasonInvalidItem
	ReasonIndefiniteOrReserved
	ReasonTruncatedArgument
	ReasonNonpreferredArgument
)

// reasonNames gives each Reason a stable identifier. It is indexed by the
// Reason's own iota value, so a Reason added to the block above without a
// name here fails TestReasonNamesCoverEveryReason rather than silently
// formatting as "unknown".
var reasonNames = [...]string{
	ReasonEmptyPayload:          "empty-payload",
	ReasonRawByteLimit:          "raw-byte-limit",
	ReasonTrailingData:          "trailing-data",
	ReasonDepthLimit:            "depth-limit",
	ReasonItemLimit:             "item-limit",
	ReasonTruncatedItem:         "truncated-item",
	ReasonUnsafeInteger:         "unsafe-integer",
	ReasonUnsafeNegativeInteger: "unsafe-negative-integer",
	ReasonTruncatedString:       "truncated-string",
	ReasonDecodedStringLimit:    "decoded-string-limit",
	ReasonInvalidUTF8:           "invalid-utf8",
	ReasonArrayElementLimit:     "array-element-limit",
	ReasonMapEntryLimit:         "map-entry-limit",
	ReasonDuplicateMapKeyOrder:  "duplicate-map-key-order",
	ReasonSemanticTag:           "semantic-tag",
	ReasonSimpleOrFloat:         "simple-or-float",
	ReasonMajorType:             "major-type",
	ReasonInvalidItem:           "invalid-item",
	ReasonIndefiniteOrReserved:  "indefinite-or-reserved",
	ReasonTruncatedArgument:     "truncated-argument",
	ReasonNonpreferredArgument:  "nonpreferred-argument",
}

// String returns a stable identifier for the Reason. It is diagnostic text
// for logs and test failures, not protocol surface: callers still map Reason
// values to their own refusal codes and classifications.
func (r Reason) String() string {
	if r < 0 || int(r) >= len(reasonNames) || reasonNames[r] == "" {
		return fmt.Sprintf("unknown(%d)", int(r))
	}
	return reasonNames[r]
}

// Error reports a predecode refusal. Callers extract Reason via errors.As
// and own their own wording; the message names the reason so a refusal is
// still legible in a log line or a test failure that only has the error.
type Error struct{ Reason Reason }

func (e *Error) Error() string { return "cborscan: predecode refused: " + e.Reason.String() }

func reject(reason Reason) error { return &Error{Reason: reason} }

// Predecode applies profile's allocation-independent deterministic-CBOR
// predecode pass to received: it proves the bytes are one well-formed,
// canonically encoded CBOR item within the profile's bounds, with no
// trailing data, before any allocating strict decoder runs.
func Predecode(received []byte, profile Profile) error {
	if len(received) == 0 {
		return reject(ReasonEmptyPayload)
	}
	if len(received) > profile.MaxBytes {
		return reject(ReasonRawByteLimit)
	}
	scanner := Scanner{Value: received, Profile: profile}
	if err := scanner.Item(1); err != nil {
		return err
	}
	if scanner.Offset != len(received) {
		return reject(ReasonTrailingData)
	}
	return nil
}

// Scanner is the reusable predecode cursor over Value, bounded by Profile.
// It is exported so a caller that needs to locate an already-accepted
// item's byte range (for example, a same-package mutation test) can drive
// the identical scan logic Predecode uses, one call at a time.
type Scanner struct {
	Value   []byte
	Offset  int
	Profile Profile

	items            int
	decodedStringLen uint64
}

// Head parses one CBOR major-type/argument head at s.Offset, advancing
// s.Offset past it on success.
func (s *Scanner) Head() (major byte, argument uint64, err error) {
	major, argument, next, err := ReadHead(s.Value, s.Offset)
	if err != nil {
		return 0, 0, err
	}
	s.Offset = next
	return major, argument, nil
}

// Item scans one complete CBOR data item (recursing into array/map/tag
// content as needed) starting at s.Offset, enforcing every Profile bound
// along the way.
func (s *Scanner) Item(depth int) error {
	if depth > s.Profile.MaxDepth {
		return reject(ReasonDepthLimit)
	}
	s.items++
	if s.items > s.Profile.MaxItems {
		return reject(ReasonItemLimit)
	}

	start := s.Offset
	major, argument, err := s.Head()
	if err != nil {
		return err
	}
	switch major {
	case 0:
		if argument > maxSafeInteger {
			return reject(ReasonUnsafeInteger)
		}
	case 1:
		if argument >= maxSafeInteger {
			return reject(ReasonUnsafeNegativeInteger)
		}
	case 2, 3:
		if argument > uint64(len(s.Value)-s.Offset) { // #nosec G115 -- the nonnegative remaining slice length widens losslessly.
			return reject(ReasonTruncatedString)
		}
		// Compare additively rather than subtracting from MaxBytes: an
		// argument larger than MaxBytes made the subtraction wrap near 2^64
		// and the check pass. Predecode cannot produce that (it refuses any
		// payload over MaxBytes first), but Scanner is exported and a caller
		// driving it directly, or supplying a small MaxBytes, could.
		if argument > uint64(max(s.Profile.MaxBytes, 0))-s.decodedStringLen { // #nosec G115 -- the max() above proves the operand is nonnegative.
			return reject(ReasonDecodedStringLimit)
		}
		s.decodedStringLen += argument
		end := s.Offset + int(argument) // #nosec G115 -- the preceding comparison proves argument fits the remaining in-memory slice.
		if major == 3 && !utf8.Valid(s.Value[s.Offset:end]) {
			return reject(ReasonInvalidUTF8)
		}
		s.Offset = end
	case 4:
		if !s.Profile.AllowArray {
			return reject(ReasonMajorType)
		}
		if argument > s.Profile.MaxArrayElements {
			return reject(ReasonArrayElementLimit)
		}
		for index := uint64(0); index < argument; index++ {
			if err := s.Item(depth + 1); err != nil {
				return err
			}
		}
	case 5:
		if argument > s.Profile.MaxMapEntries {
			return reject(ReasonMapEntryLimit)
		}
		var previousKey []byte
		for index := uint64(0); index < argument; index++ {
			keyStart := s.Offset
			if err := s.Item(depth + 1); err != nil {
				return err
			}
			key := s.Value[keyStart:s.Offset]
			if previousKey != nil && CompareDeterministic(previousKey, key) >= 0 {
				return reject(ReasonDuplicateMapKeyOrder)
			}
			previousKey = key
			if err := s.Item(depth + 1); err != nil {
				return err
			}
		}
	case 6:
		if !s.Profile.AllowTag || argument != s.Profile.AllowedTag {
			return reject(ReasonSemanticTag)
		}
		if err := s.Item(depth + 1); err != nil {
			return err
		}
	case 7:
		if !s.Profile.AllowFalse || argument != 20 {
			return reject(ReasonSimpleOrFloat)
		}
	default:
		// Unreachable through this switch: major is the top 3 bits of one
		// byte (0-7), so every value is handled by a case above. Kept as a
		// defensive refusal rather than a panic, matching both predecoders
		// this package replaces. ReasonMajorType itself is reachable — an
		// array under a profile that does not set AllowArray returns it.
		return reject(ReasonMajorType)
	}
	if s.Offset <= start {
		// Unreachable: every non-error branch above advances s.Offset past
		// start, since Head() itself always consumes at least one byte on
		// success. Kept as defensive-in-depth, matching v0candidate's
		// original scanItem.
		return reject(ReasonInvalidItem)
	}
	return nil
}

// ReadHead parses one CBOR major-type/argument head from value at offset,
// returning the major type, the decoded argument, and the offset
// immediately following the head. It performs no item-shape or profile
// validation — only the head's own truncation and non-preferred-length
// checks — so it is safe to call standalone (for example, from a reader
// that only consumes already-predecoder-accepted bytes).
func ReadHead(value []byte, offset int) (major byte, argument uint64, next int, err error) {
	if offset >= len(value) {
		return 0, 0, 0, reject(ReasonTruncatedItem)
	}
	initial := value[offset]
	offset++
	major = initial >> 5
	additional := initial & 0x1f
	if additional < 24 {
		return major, uint64(additional), offset, nil
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
		return 0, 0, 0, reject(ReasonIndefiniteOrReserved)
	}
	if width > len(value)-offset {
		return 0, 0, 0, reject(ReasonTruncatedArgument)
	}

	switch width {
	case 1:
		argument = uint64(value[offset])
	case 2:
		argument = uint64(binary.BigEndian.Uint16(value[offset : offset+width]))
	case 4:
		argument = uint64(binary.BigEndian.Uint32(value[offset : offset+width]))
	case 8:
		argument = binary.BigEndian.Uint64(value[offset : offset+width])
	}
	offset += width

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
		return 0, 0, 0, reject(ReasonNonpreferredArgument)
	}
	return major, argument, offset, nil
}

// CompareDeterministic orders two encoded CBOR map keys by RFC 8949 §4.2.1
// canonical/deterministic ordering: shorter encoding first, then
// byte-lexicographic.
func CompareDeterministic(left, right []byte) int {
	if len(left) < len(right) {
		return -1
	}
	if len(left) > len(right) {
		return 1
	}
	return bytes.Compare(left, right)
}
