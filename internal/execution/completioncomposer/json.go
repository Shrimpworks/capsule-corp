package completioncomposer

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	maximumJSONDepth         = 32
	maximumJSONNodes         = 8_193
	maximumJSONMembers       = 4_096
	maximumJSONObjectMembers = 256
	maximumJSONArrayElements = 256
	maximumJSONKeyBytes      = 128
	maximumJSONStringBytes   = MaximumResultJSONBytes
)

type jsonMetrics struct {
	nodes    int
	members  int
	elements int
}

func validateTypedJSON(exact []byte) error {
	if len(exact) == 0 {
		return classified(ClassificationMalformed, "completion-json-empty")
	}
	if len(exact) > MaximumResultJSONBytes {
		return classified(ClassificationMalformed, "completion-result-cap")
	}
	if bytes.HasPrefix(exact, []byte{0xef, 0xbb, 0xbf}) {
		return classified(ClassificationMalformed, "completion-json-bom")
	}
	if !utf8.Valid(exact) || !validEscapedSurrogates(exact) {
		return classified(ClassificationMalformed, "completion-json-utf8")
	}
	decoder := json.NewDecoder(bytes.NewReader(exact))
	decoder.UseNumber()
	metrics := jsonMetrics{}
	if err := decodeJSONValue(decoder, 1, &metrics); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return classified(ClassificationMalformed, "completion-json-trailing")
	}
	return nil
}

func decodeJSONValue(decoder *json.Decoder, depth int, metrics *jsonMetrics) error {
	if depth > maximumJSONDepth {
		return classified(ClassificationMalformed, "completion-json-depth")
	}
	metrics.nodes++
	if metrics.nodes > maximumJSONNodes {
		return classified(ClassificationMalformed, "completion-json-nodes")
	}
	token, err := decoder.Token()
	if err != nil {
		return classified(ClassificationMalformed, "completion-json-syntax")
	}
	switch value := token.(type) {
	case json.Delim:
		switch value {
		case '{':
			return decodeJSONObject(decoder, depth, metrics)
		case '[':
			return decodeJSONArray(decoder, depth, metrics)
		default:
			return classified(ClassificationMalformed, "completion-json-syntax")
		}
	case string:
		if len([]byte(value)) > maximumJSONStringBytes {
			return classified(ClassificationMalformed, "completion-json-string")
		}
		return nil
	case json.Number:
		return validateJSONInteger(string(value))
	case bool, nil:
		return nil
	default:
		return classified(ClassificationMalformed, "completion-json-type")
	}
}

func decodeJSONObject(decoder *json.Decoder, depth int, metrics *jsonMetrics) error {
	keys := make(map[string]struct{})
	count := 0
	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			return classified(ClassificationMalformed, "completion-json-syntax")
		}
		key, ok := keyToken.(string)
		if !ok || len([]byte(key)) > maximumJSONKeyBytes {
			return classified(ClassificationMalformed, "completion-json-key")
		}
		if _, exists := keys[key]; exists {
			return classified(ClassificationMalformed, "completion-json-duplicate-key")
		}
		keys[key] = struct{}{}
		count++
		metrics.members++
		if count > maximumJSONObjectMembers || metrics.members > maximumJSONMembers {
			return classified(ClassificationMalformed, "completion-json-members")
		}
		if err := decodeJSONValue(decoder, depth+1, metrics); err != nil {
			return err
		}
	}
	closing, err := decoder.Token()
	if err != nil || closing != json.Delim('}') {
		return classified(ClassificationMalformed, "completion-json-syntax")
	}
	return nil
}

func decodeJSONArray(decoder *json.Decoder, depth int, metrics *jsonMetrics) error {
	count := 0
	for decoder.More() {
		count++
		metrics.elements++
		if count > maximumJSONArrayElements || metrics.elements > maximumJSONMembers {
			return classified(ClassificationMalformed, "completion-json-elements")
		}
		if err := decodeJSONValue(decoder, depth+1, metrics); err != nil {
			return err
		}
	}
	closing, err := decoder.Token()
	if err != nil || closing != json.Delim(']') {
		return classified(ClassificationMalformed, "completion-json-syntax")
	}
	return nil
}

func validateJSONInteger(value string) error {
	if strings.ContainsAny(value, ".eE+") || value == "-0" || value == "" {
		return classified(ClassificationMalformed, "completion-json-integer")
	}
	if value[0] == '-' {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed < -int64(v0candidate.MaxSafeInteger) {
			return classified(ClassificationMalformed, "completion-json-integer")
		}
		return nil
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil || parsed > v0candidate.MaxSafeInteger {
		return classified(ClassificationMalformed, "completion-json-integer")
	}
	return nil
}

func validEscapedSurrogates(value []byte) bool {
	inString := false
	for index := 0; index < len(value); index++ {
		switch value[index] {
		case '"':
			inString = !inString
		case '\\':
			if !inString || index+1 >= len(value) {
				continue
			}
			index++
			if value[index] != 'u' || index+4 >= len(value) {
				continue
			}
			code, ok := parseHex4(value[index+1 : index+5])
			if !ok {
				continue
			}
			index += 4
			if code >= 0xdc00 && code <= 0xdfff {
				return false
			}
			if code >= 0xd800 && code <= 0xdbff {
				if index+6 >= len(value) || value[index+1] != '\\' || value[index+2] != 'u' {
					return false
				}
				low, lowOK := parseHex4(value[index+3 : index+7])
				if !lowOK || low < 0xdc00 || low > 0xdfff {
					return false
				}
				index += 6
			}
		}
	}
	return true
}

func parseHex4(value []byte) (uint16, bool) {
	if len(value) != 4 {
		return 0, false
	}
	var result uint16
	for _, character := range value {
		result <<= 4
		switch {
		case character >= '0' && character <= '9':
			result += uint16(character - '0')
		case character >= 'a' && character <= 'f':
			result += uint16(character-'a') + 10
		case character >= 'A' && character <= 'F':
			result += uint16(character-'A') + 10
		default:
			return 0, false
		}
	}
	return result, true
}
