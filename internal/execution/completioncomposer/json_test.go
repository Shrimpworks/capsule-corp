package completioncomposer

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// TestValidateTypedJSONAcceptsEveryOrdinaryValueShape covers
// validateTypedJSON/decodeJSONValue's happy path for every JSON value kind
// it accepts, including a real array -- decodeJSONArray otherwise has no
// direct test at all in this package.
func TestValidateTypedJSONAcceptsEveryOrdinaryValueShape(t *testing.T) {
	for _, exact := range []string{
		"null", "true", "false", "0", "-1", `"a string"`,
		`[]`, `[1,2,3]`, `{}`, `{"a":1}`,
		`{"a":[1,{"b":true},null],"c":"d"}`,
	} {
		if err := validateTypedJSON([]byte(exact)); err != nil {
			t.Fatalf("validateTypedJSON(%q) = %v", exact, err)
		}
	}
}

func TestValidateTypedJSONCapsAndFramingRefusals(t *testing.T) {
	cases := []struct {
		name  string
		exact []byte
	}{
		{"empty", []byte{}},
		{"oversized", append([]byte{'"'}, bytes.Repeat([]byte{'a'}, MaximumResultJSONBytes)...)},
		{"bom", append([]byte{0xef, 0xbb, 0xbf}, []byte("null")...)},
		{"invalid-utf8", []byte("\"\xff\"")},
		{"trailing-value", []byte("null true")},
		{"top-level-syntax-error", []byte("}")},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if err := validateTypedJSON(testCase.exact); classification(t, err) != ClassificationMalformed {
				t.Fatalf("%s = %v", testCase.name, err)
			}
		})
	}
}

func TestValidateTypedJSONDepthCapExactAndPlusOne(t *testing.T) {
	opening := strings.Repeat("[", maximumJSONDepth)
	closing := strings.Repeat("]", maximumJSONDepth)
	if err := validateTypedJSON([]byte(opening + closing)); err != nil {
		t.Fatalf("exact depth = %v", err)
	}
	openingOver := strings.Repeat("[", maximumJSONDepth+1)
	closingOver := strings.Repeat("]", maximumJSONDepth+1)
	if err := validateTypedJSON([]byte(openingOver + closingOver)); classification(t, err) != ClassificationMalformed {
		t.Fatalf("depth cap+1 = %v", err)
	}
}

func TestValidateTypedJSONNodesCapPlusOneRefuses(t *testing.T) {
	elements := make([]string, maximumJSONNodes) // + 1 for the array itself = cap+1 total nodes
	for index := range elements {
		elements[index] = "0"
	}
	exact := "[" + strings.Join(elements, ",") + "]"
	if err := validateTypedJSON([]byte(exact)); classification(t, err) != ClassificationMalformed {
		t.Fatalf("nodes cap+1 = %v", err)
	}
}

// TestDecodeJSONValueStringCapDirectly covers decodeJSONValue's own
// per-string byte cap by calling it directly rather than through
// validateTypedJSON. maximumJSONStringBytes == MaximumResultJSONBytes, so no
// document validateTypedJSON accepts can ever contain a string long enough
// to trip this check: the string's content plus its two quote bytes always
// exceeds validateTypedJSON's own top-level document-byte cap first (see
// TestValidateTypedJSONCapsAndFramingRefusals's "oversized" case). This
// check is real defense-in-depth, not dead code -- it holds even if a
// future change ever lets the two caps diverge -- so it's worth exercising
// directly instead of leaving uncovered.
func TestDecodeJSONValueStringCapDirectly(t *testing.T) {
	over := append([]byte{'"'}, bytes.Repeat([]byte{'a'}, maximumJSONStringBytes+1)...)
	over = append(over, '"')
	decoder := json.NewDecoder(bytes.NewReader(over))
	decoder.UseNumber()
	metrics := jsonMetrics{}
	if err := decodeJSONValue(decoder, 1, &metrics); classification(t, err) != ClassificationMalformed {
		t.Fatalf("string value cap+1 = %v", err)
	}
}

func TestDecodeJSONObjectKeyAndMemberCaps(t *testing.T) {
	longKey := `"` + strings.Repeat("k", maximumJSONKeyBytes+1) + `":0`
	if err := validateTypedJSON([]byte("{" + longKey + "}")); classification(t, err) != ClassificationMalformed {
		t.Fatalf("long key = %v", err)
	}

	members := make([]string, maximumJSONObjectMembers+1)
	for index := range members {
		members[index] = `"k` + string(rune('a'+index%26)) + string(rune('A'+index/26)) + `":0`
	}
	if err := validateTypedJSON([]byte("{" + strings.Join(members, ",") + "}")); classification(t, err) != ClassificationMalformed {
		t.Fatalf("object member cap+1 = %v", err)
	}

	if err := validateTypedJSON([]byte(`{"a":1,"a":2}`)); classification(t, err) != ClassificationMalformed {
		t.Fatalf("duplicate key = %v", err)
	}

	if err := validateTypedJSON([]byte(`{"a":1`)); classification(t, err) != ClassificationMalformed {
		t.Fatalf("truncated object = %v", err)
	}
}

func TestDecodeJSONArrayElementCapAndTruncation(t *testing.T) {
	elements := make([]string, maximumJSONArrayElements+1)
	for index := range elements {
		elements[index] = "0"
	}
	if err := validateTypedJSON([]byte("[" + strings.Join(elements, ",") + "]")); classification(t, err) != ClassificationMalformed {
		t.Fatalf("array element cap+1 = %v", err)
	}

	if err := validateTypedJSON([]byte(`[1,2`)); classification(t, err) != ClassificationMalformed {
		t.Fatalf("truncated array = %v", err)
	}
}

func TestValidateJSONIntegerBoundsAndSyntax(t *testing.T) {
	cases := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"exact-max-safe", "9007199254740991", false},
		{"exact-min-safe", "-9007199254740991", false},
		{"zero", "0", false},
		{"positive-over-cap", "9007199254740992", true},
		{"negative-over-cap", "-9007199254740992", true},
		{"positive-overflow-uint64", "99999999999999999999999", true},
		{"negative-overflow-int64", "-99999999999999999999999", true},
		{"negative-zero", "-0", true},
		{"empty", "", true},
		{"decimal-point", "1.5", true},
		{"exponent-lower", "1e2", true},
		{"exponent-upper", "1E2", true},
		{"explicit-plus", "+1", true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateJSONInteger(testCase.value)
			if testCase.wantErr && classification(t, err) != ClassificationMalformed {
				t.Fatalf("validateJSONInteger(%q) = %v, want malformed", testCase.value, err)
			}
			if !testCase.wantErr && err != nil {
				t.Fatalf("validateJSONInteger(%q) = %v, want nil", testCase.value, err)
			}
		})
	}
}

func TestValidEscapedSurrogatesAcceptsAndRejects(t *testing.T) {
	// Padded so the high-surrogate branch's own length guard
	// (index+6>=len(value)) doesn't short-circuit before reaching the
	// escape-shape and low-surrogate-value checks these cases target.
	pad := strings.Repeat("A", 10)
	cases := []struct {
		name  string
		value string
		want  bool
	}{
		{"no-escapes", `"plain"`, true},
		{"valid-surrogate-pair", `"😀"`, true},
		{"lone-high-surrogate-then-close-quote", `"\ud83d"`, false},
		{"lone-low-surrogate", `"\ude00"`, false},
		{"high-surrogate-then-non-escape-with-padding", `"\ud83dX` + pad + `"`, false},
		{"high-surrogate-then-non-u-escape-with-padding", `"\ud83d\n` + pad + `"`, false},
		{"high-surrogate-then-valid-escape-out-of-range-low-value", `"\ud83d\u0041` + pad + `"`, false},
		{"valid-surrogate-pair-with-padding", `"😀` + pad + `"`, true},
		{"trailing-backslash", "\"\\", true},
		{"short-unicode-escape", `"\u12"`, true},
		{"ordinary-escape", `"\n"`, true},
		{"backslash-not-in-string", "x\\y", true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := validEscapedSurrogates([]byte(testCase.value)); got != testCase.want {
				t.Fatalf("validEscapedSurrogates(%q) = %v, want %v", testCase.value, got, testCase.want)
			}
		})
	}
}

func TestParseHex4(t *testing.T) {
	cases := []struct {
		name   string
		value  string
		want   uint16
		wantOK bool
	}{
		{"lowercase", "d83d", 0xd83d, true},
		{"uppercase", "D83D", 0xd83d, true},
		{"digits", "0123", 0x0123, true},
		{"wrong-length", "abc", 0, false},
		{"invalid-character", "abcg", 0, false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got, ok := parseHex4([]byte(testCase.value))
			if got != testCase.want || ok != testCase.wantOK {
				t.Fatalf("parseHex4(%q) = (%v, %v), want (%v, %v)", testCase.value, got, ok, testCase.want, testCase.wantOK)
			}
		})
	}
}
