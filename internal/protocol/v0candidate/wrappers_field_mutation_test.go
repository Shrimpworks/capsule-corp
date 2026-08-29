package v0candidate

import (
	"bytes"
	"testing"
)

// This file closes issue #273's branch-coverage gap on decodeExecutionPlan
// and decodePlanRegistration. Both functions are flat, linear field-by-field
// CBOR decoders: for every field they call reader.requiredKey(N), decode the
// field's value, and (for a few fields) check the decoded value against a
// closed set of accepted values. Every one of those checks is a distinct
// refusal branch that a truncated or corrupted wire field can hit, but only
// the "happy path" plus a handful of branches were previously exercised.
//
// The tests below call decodeExecutionPlan and decodePlanRegistration
// directly (both unexported, same-package white-box tests) rather than
// through DecodeExecutionPlan/DecodePlanRegistration. That is deliberate:
// the public wrappers run PredecodeExecutionPlanCBOR/PredecodePlanRegistrationCBOR
// first, which reject a map whose header entry count does not match its
// actual encoded pairs before the strict decoder ever runs — exactly the
// byte shape a "drop field N" mutation produces. Calling the strict decoders
// directly is the only way to exercise their own per-field refusal branches
// in isolation, including the couple of branches (the PlanRegistration
// object type observed inside decodeExecutionPlan, and vice versa) that the
// public entry points' hasObjectType prefilter makes unreachable from
// outside the package.

// dropTopLevelField removes one top-level map key's key/value pair entirely
// while leaving the map header's own declared entry count untouched. The
// mutated bytes model a field that went missing on the wire: every later
// requiredKey call still expects to see the key it was written for, so the
// decoder must notice the next real key is the wrong one and refuse.
func dropTopLevelField(t *testing.T, source []byte, profile cborProfile, targetKey uint64) []byte {
	t.Helper()
	keyStart, _, valueEnd := topLevelFieldRange(t, source, profile, targetKey)
	result := make([]byte, 0, len(source)-(valueEnd-keyStart))
	result = append(result, source[:keyStart]...)
	result = append(result, source[valueEnd:]...)
	return result
}

// replaceTopLevelFieldValue swaps one top-level map key's encoded value for
// newValue, leaving its key and every other field's bytes untouched.
func replaceTopLevelFieldValue(t *testing.T, source []byte, profile cborProfile, targetKey uint64, newValue []byte) []byte {
	t.Helper()
	_, valueStart, valueEnd := topLevelFieldRange(t, source, profile, targetKey)
	result := make([]byte, 0, len(source)-(valueEnd-valueStart)+len(newValue))
	result = append(result, source[:valueStart]...)
	result = append(result, newValue...)
	result = append(result, source[valueEnd:]...)
	return result
}

// topLevelFieldRange scans a top-level CBOR map for exactly one occurrence
// of targetKey, returning the byte offsets bracketing its key
// (keyStart..valueStart) and its value (valueStart..valueEnd) within source.
// It is the shared byte-location primitive the mutators above build on, and
// mirrors the scan spliceTextFieldValue already performs elsewhere in this
// package's tests, generalized to every field type rather than only text.
func topLevelFieldRange(t *testing.T, source []byte, profile cborProfile, targetKey uint64) (keyStart, valueStart, valueEnd int) {
	t.Helper()
	scanner := cborScanner{bytes: source, profile: profile}
	major, count, err := scanner.readHead()
	if err != nil || major != 5 {
		t.Fatalf("fixture is not a top-level CBOR map: %v", err)
	}
	for index := uint64(0); index < count; index++ {
		ks := scanner.offset
		if err := scanner.scanItem(1); err != nil {
			t.Fatalf("scan key %d: %v", index, err)
		}
		keyReader := cborReader{bytes: source[ks:scanner.offset]}
		key, err := keyReader.unsigned()
		if err != nil {
			t.Fatalf("decode key %d: %v", index, err)
		}
		vs := scanner.offset
		if err := scanner.scanItem(1); err != nil {
			t.Fatalf("scan value for key %d: %v", key, err)
		}
		if key == targetKey {
			return ks, vs, scanner.offset
		}
	}
	t.Fatalf("target key %d not found in fixture", targetKey)
	return 0, 0, 0
}

// setTopLevelMapHeaderCount rewrites just the top-level map header's
// declared entry count, leaving every field byte untouched. Both ordinary
// fixtures this package ships use one of two canonical header encodings: a
// 1-byte immediate count (used by plan-registration/ordinary.cbor, 10
// fields) or a 1-byte-immediate-plus-uint8 count (used by
// execution-plan/ordinary.cbor, 24 fields).
func setTopLevelMapHeaderCount(t *testing.T, source []byte, count byte) []byte {
	t.Helper()
	result := bytes.Clone(source)
	switch {
	case result[0] == 0xb8:
		result[1] = count
	case result[0]&0xe0 == 0xa0 && result[0]&0x1f < 24:
		result[0] = 0xa0 | count
	default:
		t.Fatalf("unrecognized top-level map header encoding %x", result[:2])
	}
	return result
}

// encodeCBORHead encodes a CBOR item head for majorType (0 = unsigned, 2 =
// byte string, 3 = text string, 4 = array, ...) at length, using the same
// preferred/minimal-length encoding rule every one of this package's CBOR
// head-encoding test helpers needs: an immediate count under 24, a
// one-byte-count marker for lengths up to 0xff, or a two-byte-count marker
// for lengths up to 0xffff (sufficient for every count or marker these
// tests need). appendTextValue in wrappers_test.go, and
// encodeUnsignedValue/encodeByteStringValue/encodeArrayHeaderValue below,
// all apply this exact branch to a different major type; centralizing it
// here means the minimal-encoding rule only has to be right in one place.
func encodeCBORHead(majorType byte, length int) []byte {
	base := majorType << 5
	switch {
	case length < 24:
		return []byte{base | byte(length)}
	case length <= 0xff:
		return []byte{base | 24, byte(length)}
	default:
		return []byte{base | 25, byte(length >> 8), byte(length)}
	}
}

// encodeUnsignedValue encodes an unsigned integer using preferred/minimal
// CBOR length encoding, sufficient for the small counts and markers these
// tests need.
func encodeUnsignedValue(value uint64) []byte {
	return encodeCBORHead(0, int(value))
}

// encodeByteStringValue encodes a CBOR major-type-2 byte string using
// preferred/minimal length encoding.
func encodeByteStringValue(value []byte) []byte {
	return append(encodeCBORHead(2, len(value)), value...)
}

// encodeArrayHeaderValue encodes a CBOR major-type-4 array header for count
// elements, without any element bytes.
func encodeArrayHeaderValue(count int) []byte {
	return encodeCBORHead(4, count)
}

var (
	// cborEmptyByteString is a self-contained, well-formed CBOR value (major
	// type 2, zero length) used to corrupt a field that should hold an
	// unsigned integer or byte string of the wrong specific width.
	cborEmptyByteString = []byte{0x40}
	// cborEmptyTextString is a self-contained, well-formed CBOR value (major
	// type 3, zero length) used to corrupt a field that should hold a byte
	// string.
	cborEmptyTextString = []byte{0x60}
	// cborUnsignedZero is a self-contained, well-formed CBOR value (major
	// type 0) used to corrupt a field that should hold a text string.
	cborUnsignedZero = []byte{0x00}
)

// TestDecodeExecutionPlanEveryFieldTruncated proves every one of
// decodeExecutionPlan's 24 requiredKey(N) refusal branches (issue #273) is
// reachable: dropping any one field entirely leaves the next real key
// exactly one position ahead of what the decoder expects next, which its
// requiredKey check must classify as a missing required field.
func TestDecodeExecutionPlanEveryFieldTruncated(t *testing.T) {
	t.Parallel()
	ordinary := readConformanceFixture(t, "execution-plan/ordinary.cbor")

	fields := []struct {
		key  uint64
		name string
	}{
		{1, "object-type"},
		{2, "object-version"},
		{3, "installation-id"},
		{4, "epoch-sequence"},
		{5, "epoch-digest"},
		{6, "source-manifest-digest"},
		{7, "source-entrypoint"},
		{8, "source-byte-length"},
		{9, "input-slot"},
		{10, "inline-input-digest"},
		{11, "inline-input-byte-length"},
		{12, "runtime-profile-alias"},
		{13, "runtime-bundle-manifest-digest"},
		{14, "profile-review-attestation-digests"},
		{15, "profile-registry-entry-digest"},
		{16, "backend-validation-record-digest"},
		{17, "backend-configuration-digest"},
		{18, "trust-snapshot-digest"},
		{19, "policy-decision-digest"},
		{20, "wall-time-ms"},
		{21, "wall-time-origin"},
		{22, "output-slot"},
		{23, "output-max-json-bytes"},
	}
	for _, field := range fields {
		field := field
		t.Run(field.name, func(t *testing.T) {
			t.Parallel()
			mutated := dropTopLevelField(t, ordinary, executionPlanCBORProfile, field.key)
			_, err := decodeExecutionPlan(mutated)
			assertClassification(t, err, ClassificationSchema)
		})
	}
}

// TestDecodeExecutionPlanEveryFieldCorrupted proves every one of
// decodeExecutionPlan's per-field decode/value refusal branches not already
// exercised by known-answer or boundary fixtures elsewhere in this package
// is reachable: replacing one field's value with a wrong-type-but-otherwise
// well-formed CBOR item (or, for the three closed-vocabulary fields, a
// same-type wrong value) must be refused as a schema violation.
func TestDecodeExecutionPlanEveryFieldCorrupted(t *testing.T) {
	t.Parallel()
	ordinary := readConformanceFixture(t, "execution-plan/ordinary.cbor")

	tests := []struct {
		key   uint64
		name  string
		value []byte
	}{
		{2, "object-version wrong type", cborEmptyByteString},
		{3, "installation-id wrong type", cborEmptyTextString},
		{4, "epoch-sequence wrong type", cborEmptyByteString},
		{6, "source-manifest-digest wrong type", cborEmptyTextString},
		{8, "source-byte-length wrong type", cborEmptyByteString},
		{9, "input-slot wrong value", appendTextValue(nil, "not-primary-data")},
		{10, "inline-input-digest wrong type", cborEmptyTextString},
		{11, "inline-input-byte-length wrong type", cborEmptyByteString},
		{13, "runtime-bundle-manifest-digest wrong type", cborEmptyTextString},
		{15, "profile-registry-entry-digest wrong type", cborEmptyTextString},
		{16, "backend-validation-record-digest wrong type", cborEmptyTextString},
		{17, "backend-configuration-digest wrong type", cborEmptyTextString},
		{18, "trust-snapshot-digest wrong type", cborEmptyTextString},
		{19, "policy-decision-digest wrong type", cborEmptyTextString},
		{20, "wall-time-ms wrong type", cborEmptyByteString},
		{21, "wall-time-origin wrong value", appendTextValue(nil, "not-a-real-origin")},
		{22, "output-slot wrong value", appendTextValue(nil, "not-transformed-json")},
		{23, "output-max-json-bytes wrong type", cborEmptyByteString},
		{24, "expires-at wrong type", cborEmptyByteString},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mutated := replaceTopLevelFieldValue(t, ordinary, executionPlanCBORProfile, test.key, test.value)
			_, err := decodeExecutionPlan(mutated)
			assertClassification(t, err, ClassificationSchema)
		})
	}
}

// TestDecodeExecutionPlanReviewAttestationDigestsRefusals covers the one
// array-shaped ExecutionPlan field's refusal branches individually: a
// wrong-type value, too few elements, too many elements, and one corrupted
// element.
func TestDecodeExecutionPlanReviewAttestationDigestsRefusals(t *testing.T) {
	t.Parallel()
	ordinary := readConformanceFixture(t, "execution-plan/ordinary.cbor")

	nineDigests := encodeArrayHeaderValue(9)
	for i := 0; i < 9; i++ {
		nineDigests = append(nineDigests, encodeByteStringValue(bytes.Repeat([]byte{byte(i)}, 32))...)
	}
	oneCorruptElement := encodeArrayHeaderValue(2)
	oneCorruptElement = append(oneCorruptElement, cborEmptyTextString...)
	oneCorruptElement = append(oneCorruptElement, encodeByteStringValue(bytes.Repeat([]byte{0x01}, 32))...)

	tests := []struct {
		name  string
		value []byte
	}{
		{"wrong type", cborEmptyByteString},
		{"zero elements", encodeArrayHeaderValue(0)},
		{"nine elements", nineDigests},
		{"first element wrong type", oneCorruptElement},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mutated := replaceTopLevelFieldValue(t, ordinary, executionPlanCBORProfile, 14, test.value)
			_, err := decodeExecutionPlan(mutated)
			assertClassification(t, err, ClassificationSchema)
		})
	}
}

// TestDecodeExecutionPlanObjectTypeRefusals covers decodeExecutionPlan's
// object-type field refusal branches: a wrong-type value, an unrecognized
// object type string, and the closed cross-object PlanRegistration domain
// check. The last of these is unreachable through the public
// DecodeExecutionPlan entry point, since its hasObjectType prefilter routes
// a real PlanRegistration payload away before decodeExecutionPlan ever runs;
// calling decodeExecutionPlan directly is the only way to exercise
// decodeExecutionPlan's own defense-in-depth copy of that check.
func TestDecodeExecutionPlanObjectTypeRefusals(t *testing.T) {
	t.Parallel()
	ordinary := readConformanceFixture(t, "execution-plan/ordinary.cbor")

	t.Run("wrong type", func(t *testing.T) {
		t.Parallel()
		mutated := replaceTopLevelFieldValue(t, ordinary, executionPlanCBORProfile, 1, cborUnsignedZero)
		_, err := decodeExecutionPlan(mutated)
		assertClassification(t, err, ClassificationSchema)
	})
	t.Run("unrecognized object type", func(t *testing.T) {
		t.Parallel()
		mutated := replaceTopLevelFieldValue(t, ordinary, executionPlanCBORProfile, 1, appendTextValue(nil, "capsule.something-else"))
		_, err := decodeExecutionPlan(mutated)
		assertClassification(t, err, ClassificationUnsupported)
	})
	t.Run("plan registration object type", func(t *testing.T) {
		t.Parallel()
		mutated := replaceTopLevelFieldValue(t, ordinary, executionPlanCBORProfile, 1, appendTextValue(nil, PlanRegistrationObjectType))
		_, err := decodeExecutionPlan(mutated)
		assertClassification(t, err, ClassificationDomain)
	})
}

// TestDecodeExecutionPlanObjectShapeRefusals covers decodeExecutionPlan's
// remaining structural refusal branches: a non-map top-level value and a
// map header declaring more than the closed 24-field count.
func TestDecodeExecutionPlanObjectShapeRefusals(t *testing.T) {
	t.Parallel()

	t.Run("not a map", func(t *testing.T) {
		t.Parallel()
		_, err := decodeExecutionPlan([]byte{0x00})
		assertClassification(t, err, ClassificationSchema)
	})
	t.Run("more than 24 fields", func(t *testing.T) {
		t.Parallel()
		ordinary := readConformanceFixture(t, "execution-plan/ordinary.cbor")
		mutated := setTopLevelMapHeaderCount(t, ordinary, 25)
		_, err := decodeExecutionPlan(mutated)
		assertClassification(t, err, ClassificationUnsupported)
	})
}

// TestDecodePlanRegistrationEveryFieldTruncated is decodePlanRegistration's
// analogue of TestDecodeExecutionPlanEveryFieldTruncated: dropping any one
// of its 10 fields must be refused because the next real key no longer
// matches what the decoder expects next.
func TestDecodePlanRegistrationEveryFieldTruncated(t *testing.T) {
	t.Parallel()
	ordinary := readConformanceFixture(t, "plan-registration/ordinary.cbor")

	fields := []struct {
		key  uint64
		name string
	}{
		{1, "object-type"},
		{2, "object-version"},
		{3, "registration-id"},
		{4, "registration-sequence"},
		{5, "plan-digest"},
		{6, "installation-id"},
		{7, "epoch-sequence"},
		{8, "epoch-digest"},
		{9, "supervisor-id"},
		{10, "expires-at"},
	}
	for _, field := range fields {
		field := field
		t.Run(field.name, func(t *testing.T) {
			t.Parallel()
			mutated := dropTopLevelField(t, ordinary, planRegistrationCBORProfile, field.key)
			_, err := decodePlanRegistration(mutated)
			assertClassification(t, err, ClassificationSchema)
		})
	}
}

// TestDecodePlanRegistrationEveryFieldCorrupted is decodePlanRegistration's
// analogue of TestDecodeExecutionPlanEveryFieldCorrupted, covering the
// per-field decode refusal branches not already exercised elsewhere in this
// package (RegistrationSequence's zero-value refusal is already covered by
// TestObjectShapeAndUInt53FailuresReturnNoWrapper).
func TestDecodePlanRegistrationEveryFieldCorrupted(t *testing.T) {
	t.Parallel()
	ordinary := readConformanceFixture(t, "plan-registration/ordinary.cbor")

	tests := []struct {
		key   uint64
		name  string
		value []byte
	}{
		{2, "object-version wrong type", cborEmptyByteString},
		{3, "registration-id wrong type", cborEmptyTextString},
		{5, "plan-digest wrong type", cborEmptyTextString},
		{6, "installation-id wrong type", cborEmptyTextString},
		{7, "epoch-sequence wrong type", cborEmptyByteString},
		{8, "epoch-digest wrong type", cborEmptyTextString},
		{9, "supervisor-id wrong type", cborEmptyTextString},
		{10, "expires-at wrong type", cborEmptyByteString},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mutated := replaceTopLevelFieldValue(t, ordinary, planRegistrationCBORProfile, test.key, test.value)
			_, err := decodePlanRegistration(mutated)
			assertClassification(t, err, ClassificationSchema)
		})
	}
}

// TestDecodePlanRegistrationIdentifierFieldsRejectZeroContent covers
// RegistrationID's and SupervisorID's nonzero-16-byte checks, which
// decodePlanRegistration performs inline (unlike InstallationID, which
// delegates to the shared decodeInstallationID helper).
func TestDecodePlanRegistrationIdentifierFieldsRejectZeroContent(t *testing.T) {
	t.Parallel()
	ordinary := readConformanceFixture(t, "plan-registration/ordinary.cbor")
	zero16 := encodeByteStringValue(make([]byte, 16))

	tests := []struct {
		key  uint64
		name string
	}{
		{3, "registration-id zero content"},
		{9, "supervisor-id zero content"},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mutated := replaceTopLevelFieldValue(t, ordinary, planRegistrationCBORProfile, test.key, zero16)
			_, err := decodePlanRegistration(mutated)
			assertClassification(t, err, ClassificationSchema)
		})
	}
}

// TestDecodePlanRegistrationObjectTypeRefusals is decodePlanRegistration's
// analogue of TestDecodeExecutionPlanObjectTypeRefusals.
func TestDecodePlanRegistrationObjectTypeRefusals(t *testing.T) {
	t.Parallel()
	ordinary := readConformanceFixture(t, "plan-registration/ordinary.cbor")

	t.Run("wrong type", func(t *testing.T) {
		t.Parallel()
		mutated := replaceTopLevelFieldValue(t, ordinary, planRegistrationCBORProfile, 1, cborUnsignedZero)
		_, err := decodePlanRegistration(mutated)
		assertClassification(t, err, ClassificationSchema)
	})
	t.Run("unrecognized object type", func(t *testing.T) {
		t.Parallel()
		mutated := replaceTopLevelFieldValue(t, ordinary, planRegistrationCBORProfile, 1, appendTextValue(nil, "capsule.something-else"))
		_, err := decodePlanRegistration(mutated)
		assertClassification(t, err, ClassificationUnsupported)
	})
	t.Run("execution plan object type", func(t *testing.T) {
		t.Parallel()
		mutated := replaceTopLevelFieldValue(t, ordinary, planRegistrationCBORProfile, 1, appendTextValue(nil, ExecutionPlanObjectType))
		_, err := decodePlanRegistration(mutated)
		assertClassification(t, err, ClassificationDomain)
	})
}

// TestDecodePlanRegistrationObjectShapeRefusals covers
// decodePlanRegistration's remaining structural refusal branches: a
// non-map top-level value, a map header declaring fewer than the closed
// 10-field count, and one declaring more.
func TestDecodePlanRegistrationObjectShapeRefusals(t *testing.T) {
	t.Parallel()

	t.Run("not a map", func(t *testing.T) {
		t.Parallel()
		_, err := decodePlanRegistration([]byte{0x00})
		assertClassification(t, err, ClassificationSchema)
	})
	t.Run("fewer than 10 fields", func(t *testing.T) {
		t.Parallel()
		ordinary := readConformanceFixture(t, "plan-registration/ordinary.cbor")
		mutated := setTopLevelMapHeaderCount(t, ordinary, 9)
		_, err := decodePlanRegistration(mutated)
		assertClassification(t, err, ClassificationSchema)
	})
	t.Run("more than 10 fields", func(t *testing.T) {
		t.Parallel()
		ordinary := readConformanceFixture(t, "plan-registration/ordinary.cbor")
		mutated := setTopLevelMapHeaderCount(t, ordinary, 11)
		_, err := decodePlanRegistration(mutated)
		assertClassification(t, err, ClassificationUnsupported)
	})
}

// TestDecodePlanRegistrationObjectVersionMismatchIsUnsupported covers the
// "correct type, wrong value" branch of PlanRegistration's version check.
// decodeExecutionPlan's equivalent branch is already covered by
// TestClosedObjectShapeTypeVersionAndWidths's "unknown object version" case.
func TestDecodePlanRegistrationObjectVersionMismatchIsUnsupported(t *testing.T) {
	t.Parallel()
	ordinary := readConformanceFixture(t, "plan-registration/ordinary.cbor")
	mutated := replaceTopLevelFieldValue(t, ordinary, planRegistrationCBORProfile, 2, encodeUnsignedValue(1))
	_, err := decodePlanRegistration(mutated)
	assertClassification(t, err, ClassificationUnsupported)
}
