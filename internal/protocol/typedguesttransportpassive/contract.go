package typedguesttransportpassive

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Protocol constants freeze the C5a v1 byte layout and inclusive limits.
const (
	ProtocolVersion = uint16(1)
	Method          = uint16(1)
	SourceRole      = uint16(1)
	InputRole       = uint16(2)
	CompletionRole  = uint16(3)

	InputHeaderBytes                   = 152
	CompletionHeaderBytes              = 160
	TrailerBytes                       = 64
	PayloadMaximumBytes                = 262_144
	InputFrameMaximumBytes             = 262_296
	CompletionFrameMaximumBytes        = 262_368
	CompletionRetainMaximumBytes       = 262_369
	maximumJSONDepth                   = 32
	maximumJSONNodes                   = 8_193
	maximumJSONMembers                 = 4_096
	maximumJSONObjectMembers           = 256
	maximumJSONArrayElements           = 256
	maximumJSONKeyBytes                = 128
	maximumSafeInteger           int64 = 9_007_199_254_740_991
)

var (
	sourceMagic     = [8]byte{'C', 'P', 'S', 'R', 'C', '0', '0', '1'}
	inputMagic      = [8]byte{'C', 'P', 'I', 'N', 'P', '0', '0', '1'}
	completionMagic = [8]byte{'C', 'P', 'C', 'M', 'P', '0', '0', '1'}
	trailerMagic    = [8]byte{'C', 'P', 'E', 'N', 'D', '0', '0', '1'}
)

// Status is the closed completion status encoded at completion-header offset 112.
type Status uint16

// Completion status values expose no guest-controlled failure text.
const (
	StatusSucceeded       Status = 1
	StatusWorkloadFailed  Status = 2
	StatusResultInvalid   Status = 3
	StatusChildTerminated Status = 4
)

// Disposition identifies the first refusing control in the frozen decoder order.
type Disposition string

// Frame refusal dispositions are applied in the order documented by ADR-0046.
const (
	FrameOversize         Disposition = "FRAME_OVERSIZE"
	HeaderTruncated       Disposition = "HEADER_TRUNCATED"
	MagicMismatch         Disposition = "MAGIC"
	VersionMismatch       Disposition = "PROTOCOL_VERSION"
	MethodMismatch        Disposition = "METHOD"
	RoleMismatch          Disposition = "ROLE"
	HeaderLengthMismatch  Disposition = "HEADER_LENGTH"
	FlagsReserved         Disposition = "FLAGS_RESERVED"
	IdentifierInvalid     Disposition = "IDENTIFIER"
	BindingMismatch       Disposition = "BINDING"
	PayloadLengthCap      Disposition = "PAYLOAD_LENGTH_CAP"
	FrameLengthMismatch   Disposition = "FRAME_LENGTH"
	CommitOffset          Disposition = "COMMIT_OFFSET"
	PayloadDigestMismatch Disposition = "PAYLOAD_DIGEST"
	StatusInvalid         Disposition = "STATUS"
	NonSuccessPayload     Disposition = "NON_SUCCESS_PAYLOAD"
	JSONInvalid           Disposition = "JSON"
	CommitMagicMismatch   Disposition = "COMMIT_MAGIC"
	CommitVersionMismatch Disposition = "COMMIT_VERSION"
	CommitMethodMismatch  Disposition = "COMMIT_METHOD"
	CommitRoleMismatch    Disposition = "COMMIT_ROLE"
	CommitLengthMismatch  Disposition = "COMMIT_LENGTH"
	CommitAttemptMismatch Disposition = "COMMIT_ATTEMPT"
	CommitDigestMismatch  Disposition = "COMMIT_DIGEST"
)

// Refusal reports one closed passive frame disposition.
type Refusal struct{ Disposition Disposition }

// Error returns the closed refusal disposition.
func (r *Refusal) Error() string { return string(r.Disposition) }

// Bindings supplies independently retained Supervisor facts for frame verification.
type Bindings struct {
	AttemptID             [16]byte
	RegistrationID        [16]byte
	PlanDigest            [32]byte
	RuntimeProfileDigest  [32]byte
	ExpectedPayloadBytes  *uint64
	ExpectedPayloadDigest *[32]byte
}

// Frame owns a defensive copy of one accepted passive payload.
type Frame struct {
	Role    uint16
	Status  Status
	Payload []byte
}

// Decode verifies one exact C5a v1 frame without creating an endpoint or changing state.
func Decode(exact []byte, expectedRole uint16, bindings Bindings) (Frame, error) {
	maximum := InputFrameMaximumBytes
	headerBytes := InputHeaderBytes
	if expectedRole == CompletionRole {
		maximum = CompletionFrameMaximumBytes
		headerBytes = CompletionHeaderBytes
	} else if expectedRole != SourceRole && expectedRole != InputRole {
		return Frame{}, refuse(RoleMismatch)
	}
	if len(exact) > maximum {
		return Frame{}, refuse(FrameOversize)
	}
	if len(exact) < headerBytes {
		return Frame{}, refuse(HeaderTruncated)
	}
	roleMagic := magicFor(expectedRole)
	if !bytes.Equal(exact[0:8], roleMagic[:]) {
		return Frame{}, refuse(MagicMismatch)
	}
	if binary.BigEndian.Uint16(exact[8:10]) != ProtocolVersion {
		return Frame{}, refuse(VersionMismatch)
	}
	if binary.BigEndian.Uint16(exact[10:12]) != Method {
		return Frame{}, refuse(MethodMismatch)
	}
	if binary.BigEndian.Uint16(exact[12:14]) != expectedRole {
		return Frame{}, refuse(RoleMismatch)
	}
	if binary.BigEndian.Uint16(exact[14:16]) != uint16(headerBytes) {
		return Frame{}, refuse(HeaderLengthMismatch)
	}
	if expectedRole == CompletionRole &&
		(binary.BigEndian.Uint16(exact[114:116]) != 0 || binary.BigEndian.Uint32(exact[116:120]) != 0) {
		return Frame{}, refuse(FlagsReserved)
	}
	if allZero(exact[16:32]) || allZero(exact[32:48]) || bytes.Equal(exact[16:32], exact[32:48]) {
		return Frame{}, refuse(IdentifierInvalid)
	}
	if !bytes.Equal(exact[16:32], bindings.AttemptID[:]) ||
		!bytes.Equal(exact[32:48], bindings.RegistrationID[:]) ||
		!bytes.Equal(exact[48:80], bindings.PlanDigest[:]) ||
		!bytes.Equal(exact[80:112], bindings.RuntimeProfileDigest[:]) {
		return Frame{}, refuse(BindingMismatch)
	}
	lengthOffset := 112
	digestOffset := 120
	if expectedRole == CompletionRole {
		lengthOffset = 120
		digestOffset = 128
	}
	declared := binary.BigEndian.Uint64(exact[lengthOffset : lengthOffset+8])
	if declared > PayloadMaximumBytes {
		return Frame{}, refuse(PayloadLengthCap)
	}
	if bindings.ExpectedPayloadBytes != nil && declared != *bindings.ExpectedPayloadBytes {
		return Frame{}, refuse(BindingMismatch)
	}
	payloadEnd := headerBytes + int(declared)
	if len(exact) < payloadEnd {
		return Frame{}, refuse(FrameLengthMismatch)
	}
	expectedLength := payloadEnd
	if expectedRole == CompletionRole {
		expectedLength += TrailerBytes
		if len(exact) < expectedLength {
			return Frame{}, refuse(CommitOffset)
		}
	}
	if len(exact) != expectedLength {
		return Frame{}, refuse(FrameLengthMismatch)
	}
	payload := exact[headerBytes:payloadEnd]
	digest := sha256.Sum256(payload)
	if !bytes.Equal(digest[:], exact[digestOffset:digestOffset+32]) {
		return Frame{}, refuse(PayloadDigestMismatch)
	}
	if bindings.ExpectedPayloadDigest != nil && digest != *bindings.ExpectedPayloadDigest {
		return Frame{}, refuse(BindingMismatch)
	}
	frame := Frame{Role: expectedRole, Payload: append([]byte(nil), payload...)}
	if expectedRole == CompletionRole {
		frame.Status = Status(binary.BigEndian.Uint16(exact[112:114]))
		if frame.Status < StatusSucceeded || frame.Status > StatusChildTerminated {
			return Frame{}, refuse(StatusInvalid)
		}
		if frame.Status != StatusSucceeded && !bytes.Equal(payload, []byte("null")) {
			return Frame{}, refuse(NonSuccessPayload)
		}
	}
	if expectedRole != SourceRole {
		if err := ValidateCanonicalJSON(payload); err != nil {
			return Frame{}, refuse(JSONInvalid)
		}
	}
	if expectedRole == CompletionRole {
		trailer := exact[payloadEnd:]
		if !bytes.Equal(trailer[0:8], trailerMagic[:]) {
			return Frame{}, refuse(CommitMagicMismatch)
		}
		if binary.BigEndian.Uint16(trailer[8:10]) != ProtocolVersion {
			return Frame{}, refuse(CommitVersionMismatch)
		}
		if binary.BigEndian.Uint16(trailer[10:12]) != Method {
			return Frame{}, refuse(CommitMethodMismatch)
		}
		if binary.BigEndian.Uint16(trailer[12:14]) != CompletionRole {
			return Frame{}, refuse(CommitRoleMismatch)
		}
		if binary.BigEndian.Uint16(trailer[14:16]) != TrailerBytes {
			return Frame{}, refuse(CommitLengthMismatch)
		}
		if !bytes.Equal(trailer[16:32], exact[16:32]) {
			return Frame{}, refuse(CommitAttemptMismatch)
		}
		committed := sha256.Sum256(exact[:payloadEnd])
		if !bytes.Equal(trailer[32:64], committed[:]) {
			return Frame{}, refuse(CommitDigestMismatch)
		}
	}
	return frame, nil
}

func magicFor(role uint16) [8]byte {
	switch role {
	case SourceRole:
		return sourceMagic
	case InputRole:
		return inputMagic
	default:
		return completionMagic
	}
}

func allZero(value []byte) bool {
	for _, item := range value {
		if item != 0 {
			return false
		}
	}
	return true
}

func refuse(disposition Disposition) error { return &Refusal{Disposition: disposition} }

type jsonMetrics struct {
	nodes, members, elements int
}

// ValidateCanonicalJSON verifies the frozen canonical inline JSON v0 byte profile.
func ValidateCanonicalJSON(exact []byte) error {
	if len(exact) == 0 || len(exact) > PayloadMaximumBytes || bytes.HasPrefix(exact, []byte{0xef, 0xbb, 0xbf}) || !utf8.Valid(exact) {
		return errors.New("json bounds or utf8")
	}
	decoder := json.NewDecoder(bytes.NewReader(exact))
	decoder.UseNumber()
	metrics := jsonMetrics{}
	value, err := decodeJSON(decoder, 1, &metrics)
	if err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return errors.New("json trailing")
	}
	encoded, err := encodeCanonicalJSON(value)
	if err != nil || !bytes.Equal(encoded, exact) {
		return errors.New("json noncanonical")
	}
	return nil
}

func decodeJSON(decoder *json.Decoder, depth int, metrics *jsonMetrics) (any, error) {
	if depth > maximumJSONDepth {
		return nil, errors.New("json depth")
	}
	metrics.nodes++
	if metrics.nodes > maximumJSONNodes {
		return nil, errors.New("json nodes")
	}
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	switch value := token.(type) {
	case json.Delim:
		switch value {
		case '{':
			result := map[string]any{}
			count := 0
			for decoder.More() {
				keyToken, keyErr := decoder.Token()
				key, ok := keyToken.(string)
				if keyErr != nil || !ok || !validKey(key) {
					return nil, errors.New("json key")
				}
				if _, duplicate := result[key]; duplicate {
					return nil, errors.New("json duplicate key")
				}
				count++
				metrics.members++
				if count > maximumJSONObjectMembers || metrics.members > maximumJSONMembers {
					return nil, errors.New("json members")
				}
				child, childErr := decodeJSON(decoder, depth+1, metrics)
				if childErr != nil {
					return nil, childErr
				}
				result[key] = child
			}
			closing, closeErr := decoder.Token()
			if closeErr != nil || closing != json.Delim('}') {
				return nil, errors.New("json object close")
			}
			return result, nil
		case '[':
			result := []any{}
			for decoder.More() {
				metrics.elements++
				if len(result) >= maximumJSONArrayElements || metrics.elements > maximumJSONMembers {
					return nil, errors.New("json elements")
				}
				child, childErr := decodeJSON(decoder, depth+1, metrics)
				if childErr != nil {
					return nil, childErr
				}
				result = append(result, child)
			}
			closing, closeErr := decoder.Token()
			if closeErr != nil || closing != json.Delim(']') {
				return nil, errors.New("json array close")
			}
			return result, nil
		}
	case json.Number:
		text := value.String()
		if strings.ContainsAny(text, ".eE+") || text == "-0" || text == "" {
			return nil, errors.New("json integer")
		}
		parsed, parseErr := strconv.ParseInt(text, 10, 64)
		if parseErr != nil || parsed < -maximumSafeInteger || parsed > maximumSafeInteger {
			return nil, errors.New("json integer")
		}
		return parsed, nil
	case string:
		if !utf8.ValidString(value) || strings.IndexFunc(value, func(r rune) bool { return r >= 0xd800 && r <= 0xdfff }) >= 0 {
			return nil, errors.New("json string")
		}
		return value, nil
	case bool, nil:
		return value, nil
	}
	return nil, errors.New("json type")
}

func validKey(value string) bool {
	if len(value) == 0 || len(value) > maximumJSONKeyBytes {
		return false
	}
	for index, character := range []byte(value) {
		if index == 0 {
			if !isASCIIAlphanumeric(character) {
				return false
			}
		} else if !isASCIIAlphanumeric(character) && character != '_' && character != '.' && character != '-' {
			return false
		}
	}
	return true
}

func isASCIIAlphanumeric(value byte) bool {
	return value >= 'A' && value <= 'Z' || value >= 'a' && value <= 'z' || value >= '0' && value <= '9'
}

func encodeCanonicalJSON(value any) ([]byte, error) {
	var builder strings.Builder
	if err := appendCanonicalJSON(&builder, value); err != nil {
		return nil, err
	}
	return []byte(builder.String()), nil
}

func appendCanonicalJSON(builder *strings.Builder, value any) error {
	switch item := value.(type) {
	case nil:
		builder.WriteString("null")
	case bool:
		builder.WriteString(strconv.FormatBool(item))
	case int64:
		builder.WriteString(strconv.FormatInt(item, 10))
	case string:
		appendCanonicalString(builder, item)
	case []any:
		builder.WriteByte('[')
		for index, child := range item {
			if index > 0 {
				builder.WriteByte(',')
			}
			if err := appendCanonicalJSON(builder, child); err != nil {
				return err
			}
		}
		builder.WriteByte(']')
	case map[string]any:
		keys := make([]string, 0, len(item))
		for key := range item {
			keys = append(keys, key)
		}
		slices.Sort(keys)
		builder.WriteByte('{')
		for index, key := range keys {
			if index > 0 {
				builder.WriteByte(',')
			}
			appendCanonicalString(builder, key)
			builder.WriteByte(':')
			if err := appendCanonicalJSON(builder, item[key]); err != nil {
				return err
			}
		}
		builder.WriteByte('}')
	default:
		return fmt.Errorf("unsupported JSON type %T", value)
	}
	return nil
}

func appendCanonicalString(builder *strings.Builder, value string) {
	builder.WriteByte('"')
	for _, character := range value {
		switch character {
		case '"':
			builder.WriteString("\\\"")
		case '\\':
			builder.WriteString("\\\\")
		default:
			if character <= 0x1f {
				fmt.Fprintf(builder, "\\u00%02x", character)
			} else {
				builder.WriteRune(character)
			}
		}
	}
	builder.WriteByte('"')
}
