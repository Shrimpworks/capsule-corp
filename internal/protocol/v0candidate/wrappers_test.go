package v0candidate

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCandidateMediaTypeFixtures(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		path           string
		validate       func(string) error
		classification Classification
	}{
		{name: "execution plan exact", path: "shared/media-execution-plan-exact.txt", validate: ValidateExecutionPlanMediaType},
		{name: "execution plan case insensitive type", path: "shared/media-execution-plan-uppercase.txt", validate: ValidateExecutionPlanMediaType},
		{name: "execution plan missing version", path: "shared/media-execution-plan-missing-version.txt", validate: ValidateExecutionPlanMediaType, classification: ClassificationMalformed},
		{name: "execution plan additional parameter", path: "shared/media-execution-plan-additional-parameter.txt", validate: ValidateExecutionPlanMediaType, classification: ClassificationMalformed},
		{name: "execution plan unknown version", path: "shared/media-execution-plan-unknown-version.txt", validate: ValidateExecutionPlanMediaType, classification: ClassificationUnsupported},
		{name: "plan registration exact", path: "shared/media-plan-registration-exact.txt", validate: ValidatePlanRegistrationMediaType},
		{name: "plan registration case insensitive type", path: "shared/media-plan-registration-uppercase.txt", validate: ValidatePlanRegistrationMediaType},
		{name: "plan registration missing version", path: "shared/media-plan-registration-missing-version.txt", validate: ValidatePlanRegistrationMediaType, classification: ClassificationMalformed},
		{name: "plan registration additional parameter", path: "shared/media-plan-registration-additional-parameter.txt", validate: ValidatePlanRegistrationMediaType, classification: ClassificationMalformed},
		{name: "plan registration unknown version", path: "shared/media-plan-registration-unknown-version.txt", validate: ValidatePlanRegistrationMediaType, classification: ClassificationUnsupported},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			value := string(readConformanceFixture(t, test.path))
			assertClassification(t, test.validate(value), test.classification)
		})
	}
}

func TestExecutionPlanPredecoderFixtures(t *testing.T) {
	t.Parallel()
	tests := append(
		cborBoundaryFixtures("execution-plan"),
		cborDeterministicFixtures()...,
	)
	tests = append(tests, cborFixture{
		name:           "invalid break byte",
		path:           "execution-plan/invalid-break.cbor",
		classification: ClassificationMalformed,
	})
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := PredecodeExecutionPlanCBOR(readConformanceFixture(t, test.path))
			assertClassification(t, err, test.classification)
		})
	}
}

func TestPlanRegistrationPredecoderFixtures(t *testing.T) {
	t.Parallel()
	tests := append(
		cborBoundaryFixtures("plan-registration"),
		cborDeterministicFixtures()...,
	)
	tests = append(tests, cborFixture{
		name: "ordinary registration",
		path: "plan-registration/ordinary.cbor",
	})
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := PredecodePlanRegistrationCBOR(readConformanceFixture(t, test.path))
			assertClassification(t, err, test.classification)
		})
	}
}

func TestExecutionPlanWrapperRetainsKnownAnswerAndCopiesEveryView(t *testing.T) {
	t.Parallel()
	received := readConformanceFixture(t, "execution-plan/ordinary.cbor")
	expected := bytes.Clone(received)
	decoded, err := DecodeExecutionPlan(received, ordinaryExecutionPlanBindings())
	if err != nil {
		t.Fatalf("decode ordinary ExecutionPlan: %v", err)
	}
	received[0] ^= 0xff

	if !bytes.Equal(decoded.AuthoritativeBytes(), expected) {
		t.Fatal("authoritative bytes changed after caller-buffer mutation")
	}
	if len(decoded.AuthoritativeBytes()) != 530 {
		t.Fatalf("unexpected authoritative byte length %d", len(decoded.AuthoritativeBytes()))
	}
	expectedDigest := mustHex32[ExecutionPlanDigest](t, "627f9524479000dab6f3cee1d70c0428c63285bcadbc2cb3c6e8018b2dea008c")
	if decoded.Digest() != expectedDigest {
		t.Fatalf("unexpected exact-byte plan digest %x", decoded.Digest())
	}

	returnedBytes := decoded.AuthoritativeBytes()
	returnedBytes[0] ^= 0xff
	if !bytes.Equal(decoded.AuthoritativeBytes(), expected) {
		t.Fatal("authoritative bytes changed after returned-buffer mutation")
	}
	view := decoded.View()
	if view.ObjectType != ExecutionPlanObjectType || view.ObjectVersion != CandidateObjectVersion {
		t.Fatalf("unexpected object header: %#v", view)
	}
	if view.SourceEntrypoint != "src/main.ts" || view.WallTimeMS != 5_000 || view.ExpiresAt != 1_785_456_300 {
		t.Fatalf("unexpected decoded ExecutionPlan view: %#v", view)
	}
	view.ProfileReviewAttestationDigests[0][0] ^= 0xff
	if decoded.View().ProfileReviewAttestationDigests[0][0] != 0x66 {
		t.Fatal("decoded view retained caller-mutable review-digest storage")
	}
}

func TestPlanRegistrationWrapperRetainsKnownAnswerAndCopiesCallerBuffer(t *testing.T) {
	t.Parallel()
	received := readConformanceFixture(t, "plan-registration/ordinary.cbor")
	expected := bytes.Clone(received)
	decoded, err := DecodePlanRegistration(received, ordinaryPlanRegistrationBindings())
	if err != nil {
		t.Fatalf("decode ordinary PlanRegistration: %v", err)
	}
	received[0] ^= 0xff

	if !bytes.Equal(decoded.AuthoritativeBytes(), expected) {
		t.Fatal("authoritative registration bytes changed after caller-buffer mutation")
	}
	if len(decoded.AuthoritativeBytes()) != 165 {
		t.Fatalf("unexpected authoritative registration byte length %d", len(decoded.AuthoritativeBytes()))
	}
	if digest := sha256.Sum256(decoded.AuthoritativeBytes()); hex.EncodeToString(digest[:]) != "f3569d37ad6d787c2cdd575ef9ec6c369bbe495157c43110fc9e9d610a277614" {
		t.Fatalf("unexpected registration known-answer digest %x", digest)
	}
	view := decoded.View()
	if view.ObjectType != PlanRegistrationObjectType || view.RegistrationSequence != 1 || view.ExpiresAt != 1_785_456_300 {
		t.Fatalf("unexpected decoded PlanRegistration view: %#v", view)
	}
	returned := decoded.AuthoritativeBytes()
	returned[0] ^= 0xff
	if !bytes.Equal(decoded.AuthoritativeBytes(), expected) {
		t.Fatal("authoritative registration bytes changed after returned-buffer mutation")
	}
}

func TestCrossObjectFixturesFailAtDomainBoundary(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		decode func([]byte) error
		path   string
	}{
		{
			name: "registration at execution plan entry point",
			path: "execution-plan/plan-registration-object.cbor",
			decode: func(value []byte) error {
				_, err := DecodeExecutionPlan(value, ordinaryExecutionPlanBindings())
				return err
			},
		},
		{
			name: "execution plan at registration entry point",
			path: "plan-registration/execution-plan-object.cbor",
			decode: func(value []byte) error {
				_, err := DecodePlanRegistration(value, ordinaryPlanRegistrationBindings())
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assertClassification(t, test.decode(readConformanceFixture(t, test.path)), ClassificationDomain)
		})
	}
}

func TestWrongDomainExecutionPlanFixturesFailClosed(t *testing.T) {
	t.Parallel()
	paths := []string{
		"execution-plan/domain-supervisor-as-installation.cbor",
		"execution-plan/domain-inline-input-as-epoch.cbor",
		"execution-plan/domain-inline-input-as-source.cbor",
		"execution-plan/domain-source-manifest-as-inline-input.cbor",
		"execution-plan/domain-profile-registry-as-runtime-bundle.cbor",
		"execution-plan/domain-runtime-bundle-as-review-attestation.cbor",
		"execution-plan/domain-review-attestation-as-profile-registry.cbor",
		"execution-plan/domain-backend-config-as-validation.cbor",
		"execution-plan/domain-backend-validation-as-config.cbor",
		"execution-plan/domain-policy-as-trust-snapshot.cbor",
		"execution-plan/domain-trust-snapshot-as-policy.cbor",
	}
	for _, path := range paths {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			t.Parallel()
			decoded, err := DecodeExecutionPlan(readConformanceFixture(t, path), ordinaryExecutionPlanBindings())
			if decoded != nil {
				t.Fatal("wrong-domain plan produced a decoded wrapper")
			}
			assertClassification(t, err, ClassificationDomain)
		})
	}
}

func TestWrongDomainPlanRegistrationFixturesFailClosed(t *testing.T) {
	t.Parallel()
	paths := []string{
		"plan-registration/domain-installation-as-registration.cbor",
		"plan-registration/domain-epoch-as-plan-digest.cbor",
		"plan-registration/domain-plan-as-epoch-digest.cbor",
		"plan-registration/domain-registration-as-supervisor.cbor",
	}
	for _, path := range paths {
		path := path
		t.Run(filepath.Base(path), func(t *testing.T) {
			t.Parallel()
			decoded, err := DecodePlanRegistration(readConformanceFixture(t, path), ordinaryPlanRegistrationBindings())
			if decoded != nil {
				t.Fatal("wrong-domain registration produced a decoded wrapper")
			}
			assertClassification(t, err, ClassificationDomain)
		})
	}
}

func TestObjectShapeAndUInt53FailuresReturnNoWrapper(t *testing.T) {
	t.Parallel()
	plan := readConformanceFixture(t, "execution-plan/ordinary.cbor")
	zeroInstallation := replaceFirst(t, plan, bytes.Repeat([]byte{0x11}, 16), make([]byte, 16))
	decodedPlan, err := DecodeExecutionPlan(zeroInstallation, ordinaryExecutionPlanBindings())
	if decodedPlan != nil {
		t.Fatal("zero installation ID produced a decoded plan")
	}
	assertClassification(t, err, ClassificationSchema)

	registration := readConformanceFixture(t, "plan-registration/ordinary.cbor")
	zeroSequence := replaceFirst(t, registration, []byte{0x04, 0x01, 0x05}, []byte{0x04, 0x00, 0x05})
	decodedRegistration, err := DecodePlanRegistration(zeroSequence, ordinaryPlanRegistrationBindings())
	if decodedRegistration != nil {
		t.Fatal("zero registration sequence produced a decoded registration")
	}
	assertClassification(t, err, ClassificationSchema)

	unsafeUInt53 := []byte{0x1b, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	maximumUInt53 := []byte{0x1b, 0x00, 0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	assertClassification(t, PredecodeExecutionPlanCBOR(maximumUInt53), "")
	assertClassification(t, PredecodePlanRegistrationCBOR(maximumUInt53), "")
	assertClassification(t, PredecodeExecutionPlanCBOR(unsafeUInt53), ClassificationMalformed)
	assertClassification(t, PredecodePlanRegistrationCBOR(unsafeUInt53), ClassificationMalformed)
}

func TestClosedObjectShapeTypeVersionAndWidths(t *testing.T) {
	t.Parallel()
	plan := readConformanceFixture(t, "execution-plan/ordinary.cbor")
	versionOne := replaceFirst(t, plan, []byte{0x02, 0x00, 0x03}, []byte{0x02, 0x01, 0x03})
	decoded, err := DecodeExecutionPlan(versionOne, ordinaryExecutionPlanBindings())
	if decoded != nil {
		t.Fatal("unknown object version produced a decoded plan")
	}
	assertClassification(t, err, ClassificationUnsupported)

	missingLastField := append([]byte{0xb7}, plan[2:len(plan)-7]...)
	decoded, err = DecodeExecutionPlan(missingLastField, ordinaryExecutionPlanBindings())
	if decoded != nil {
		t.Fatal("missing required field produced a decoded plan")
	}
	assertClassification(t, err, ClassificationSchema)

	unknownLastField := bytes.Clone(plan)
	lastField := bytes.LastIndex(unknownLastField, []byte{0x18, 0x18, 0x1a})
	if lastField < 0 {
		t.Fatal("final ExecutionPlan field label not found")
	}
	unknownLastField[lastField+1] = 0x19
	decoded, err = DecodeExecutionPlan(unknownLastField, ordinaryExecutionPlanBindings())
	if decoded != nil {
		t.Fatal("unknown field produced a decoded plan")
	}
	assertClassification(t, err, ClassificationUnsupported)

	shortDigest := bytes.Clone(plan)
	digestHead := bytes.Index(shortDigest, []byte{0x58, 0x20, 0x22, 0x22})
	if digestHead < 0 {
		t.Fatal("ExecutionPlan epoch digest not found")
	}
	shortDigest[digestHead+1] = 0x1f
	shortDigest = append(shortDigest[:digestHead+2+31], shortDigest[digestHead+2+32:]...)
	decoded, err = DecodeExecutionPlan(shortDigest, ordinaryExecutionPlanBindings())
	if decoded != nil {
		t.Fatal("short digest produced a decoded plan")
	}
	assertClassification(t, err, ClassificationSchema)
}

func TestSharedScalarFixtures(t *testing.T) {
	t.Parallel()
	idTests := []struct {
		path           string
		classification Classification
	}{
		{path: "shared/id-nonzero-16.bin"},
		{path: "shared/id-short.bin", classification: ClassificationSchema},
		{path: "shared/id-long.bin", classification: ClassificationSchema},
		{path: "shared/id-zero-16.bin", classification: ClassificationSchema},
	}
	for _, test := range idTests {
		value := readConformanceFixture(t, test.path)
		_, err := NewInstallationID(value)
		if test.classification == "" {
			if err != nil {
				t.Fatalf("%s: unexpected ID error: %v", test.path, err)
			}
		} else {
			assertClassification(t, err, test.classification)
		}
	}

	digestTests := []struct {
		path           string
		classification Classification
	}{
		{path: "shared/digest-nonzero-32.bin"},
		{path: "shared/digest-zero-32.bin"},
		{path: "shared/digest-short.bin", classification: ClassificationSchema},
		{path: "shared/digest-long.bin", classification: ClassificationSchema},
	}
	for _, test := range digestTests {
		_, err := NewSourceManifestDigest(readConformanceFixture(t, test.path))
		assertClassification(t, err, test.classification)
	}
}

func TestSharedDigestRoleContextIsRequired(t *testing.T) {
	t.Parallel()
	plan := readConformanceFixture(t, "execution-plan/ordinary.cbor")
	bindings := ordinaryExecutionPlanBindings()
	if _, err := DecodeExecutionPlan(plan, bindings); err != nil {
		t.Fatalf("matching source-manifest role: %v", err)
	}
	bindings.SourceManifestDigest = SourceManifestDigest(bindings.InlineInputDigest)
	decoded, err := DecodeExecutionPlan(plan, bindings)
	if decoded != nil {
		t.Fatal("wrong digest role context produced a decoded plan")
	}
	assertClassification(t, err, ClassificationDomain)
}

type cborFixture struct {
	name           string
	path           string
	classification Classification
}

func cborBoundaryFixtures(slug string) []cborFixture {
	fixtures := make([]cborFixture, 0, 10)
	for _, dimension := range []string{"raw-bytes", "depth", "items", "map-entries", "array-elements"} {
		fixtures = append(fixtures,
			cborFixture{
				name: dimension + " exact maximum",
				path: "shared/cbor-" + slug + "-" + dimension + "-exact.bin",
			},
			cborFixture{
				name:           dimension + " cap plus one",
				path:           "shared/cbor-" + slug + "-" + dimension + "-over.bin",
				classification: ClassificationMalformed,
			},
		)
	}
	return fixtures
}

func cborDeterministicFixtures() []cborFixture {
	return []cborFixture{
		{name: "canonical map", path: "shared/cbor-profile-canonical-map.bin"},
		{name: "indefinite container", path: "shared/cbor-profile-indefinite-container.bin", classification: ClassificationMalformed},
		{name: "float", path: "shared/cbor-profile-float.bin", classification: ClassificationMalformed},
		{name: "bignum", path: "shared/cbor-profile-bignum.bin", classification: ClassificationMalformed},
		{name: "tag", path: "shared/cbor-profile-tag.bin", classification: ClassificationMalformed},
		{name: "duplicate key", path: "shared/cbor-profile-duplicate-key.bin", classification: ClassificationMalformed},
		{name: "trailing data", path: "shared/cbor-profile-trailing-data.bin", classification: ClassificationMalformed},
		{name: "nonpreferred integer", path: "shared/cbor-profile-nonpreferred-integer.bin", classification: ClassificationMalformed},
		{name: "invalid UTF-8", path: "shared/cbor-profile-invalid-utf8.bin", classification: ClassificationMalformed},
		{name: "map order", path: "shared/cbor-profile-map-order.bin", classification: ClassificationMalformed},
	}
}

func ordinaryExecutionPlanBindings() ExecutionPlanRoleBindings {
	return ExecutionPlanRoleBindings{
		InstallationID:                  repeated16[InstallationID](0x11),
		EpochDigest:                     repeated32[TrustEpochDigest](0x22),
		SourceManifestDigest:            mustHex32NoTest[SourceManifestDigest]("e5e09b2435baedf897526a89c698c0b0531437a69472372ae426f62d801fc171"),
		InlineInputDigest:               mustHex32NoTest[InlineInputDigest]("bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72"),
		RuntimeBundleManifestDigest:     repeated32[RuntimeBundleManifestDigest](0x55),
		ProfileReviewAttestationDigests: []ProfileReviewAttestationDigest{repeated32[ProfileReviewAttestationDigest](0x66), repeated32[ProfileReviewAttestationDigest](0x67)},
		ProfileRegistryEntryDigest:      repeated32[ProfileRegistryEntryDigest](0x77),
		BackendValidationRecordDigest:   repeated32[BackendValidationRecordDigest](0x88),
		BackendConfigurationDigest:      repeated32[BackendConfigurationDigest](0x99),
		TrustSnapshotDigest:             repeated32[TrustSnapshotDigest](0xaa),
		PolicyDecisionDigest:            repeated32[PolicyDecisionDigest](0xbb),
	}
}

func ordinaryPlanRegistrationBindings() PlanRegistrationRoleBindings {
	return PlanRegistrationRoleBindings{
		RegistrationID: repeated16[RegistrationID](0x77),
		PlanDigest:     mustHex32NoTest[ExecutionPlanDigest]("627f9524479000dab6f3cee1d70c0428c63285bcadbc2cb3c6e8018b2dea008c"),
		InstallationID: repeated16[InstallationID](0x11),
		EpochDigest:    repeated32[TrustEpochDigest](0x22),
		SupervisorID:   repeated16[SupervisorID](0x55),
	}
}

func repeated16[T ~[16]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}

func repeated32[T ~[32]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}

func mustHex32[T ~[32]byte](t *testing.T, value string) T {
	t.Helper()
	decoded, err := hex.DecodeString(value)
	if err != nil {
		t.Fatalf("decode test digest: %v", err)
	}
	var result T
	copy(result[:], decoded)
	return result
}

func mustHex32NoTest[T ~[32]byte](value string) T {
	decoded, err := hex.DecodeString(value)
	if err != nil {
		panic(err)
	}
	var result T
	copy(result[:], decoded)
	return result
}

func assertClassification(t *testing.T, err error, expected Classification) {
	t.Helper()
	if expected == "" {
		if err != nil {
			t.Fatalf("unexpected rejection: %v", err)
		}
		return
	}
	if err == nil {
		t.Fatalf("expected %s rejection", expected)
	}
	actual, ok := ErrorClassification(err)
	if !ok {
		t.Fatalf("error has no decode classification: %v", err)
	}
	if actual != expected {
		t.Fatalf("classification = %s, want %s: %v", actual, expected, err)
	}
}

func readConformanceFixture(t *testing.T, path string) []byte {
	t.Helper()
	value, err := os.ReadFile(filepath.Join("..", "..", "..", "schemas", "conformance", "v0", path))
	if err != nil {
		t.Fatalf("read conformance fixture %s: %v", path, err)
	}
	return value
}

func replaceFirst(t *testing.T, source, old, replacement []byte) []byte {
	t.Helper()
	if len(old) != len(replacement) {
		t.Fatal("test replacement must preserve byte length")
	}
	index := bytes.Index(source, old)
	if index < 0 {
		t.Fatalf("test pattern %x not found", old)
	}
	result := bytes.Clone(source)
	copy(result[index:index+len(old)], replacement)
	return result
}

// TestExecutionPlanBoundedTextFieldLengths proves issue #142's coverage gap:
// SourceEntrypoint's 1-256 byte bound and RuntimeProfileAlias's 1-128 byte
// bound (both enforced only by decodeBoundedText) are exercised at 0, 1,
// exact-maximum, and maximum-plus-one, mirroring this package's established
// exact-maximum/cap-plus-one pattern for every other bounded field.
func TestExecutionPlanBoundedTextFieldLengths(t *testing.T) {
	t.Parallel()
	ordinary := readConformanceFixture(t, "execution-plan/ordinary.cbor")

	tests := []struct {
		name           string
		key            uint64
		length         int
		classification Classification
	}{
		{name: "source-entrypoint zero bytes", key: 7, length: 0, classification: ClassificationSchema},
		{name: "source-entrypoint one byte", key: 7, length: 1},
		{name: "source-entrypoint exact maximum (256)", key: 7, length: 256},
		{name: "source-entrypoint maximum plus one (257)", key: 7, length: 257, classification: ClassificationSchema},
		{name: "runtime-profile-alias zero bytes", key: 12, length: 0, classification: ClassificationSchema},
		{name: "runtime-profile-alias one byte", key: 12, length: 1},
		{name: "runtime-profile-alias exact maximum (128)", key: 12, length: 128},
		{name: "runtime-profile-alias maximum plus one (129)", key: 12, length: 129, classification: ClassificationSchema},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			mutated := spliceTextFieldValue(t, ordinary, test.key, strings.Repeat("x", test.length))
			decoded, err := DecodeExecutionPlan(mutated, ordinaryExecutionPlanBindings())
			assertClassification(t, err, test.classification)
			if test.classification == "" {
				if decoded == nil {
					t.Fatal("accepted case produced no decoded plan")
				}
				view := decoded.View()
				var got string
				switch test.key {
				case 7:
					got = view.SourceEntrypoint
				case 12:
					got = view.RuntimeProfileAlias
				}
				if got != strings.Repeat("x", test.length) {
					t.Fatalf("decoded field length = %d, want %d", len(got), test.length)
				}
			}
		})
	}
}

// spliceTextFieldValue replaces the CBOR-encoded value of one top-level
// ExecutionPlan map key with a freshly encoded text string of newValue's
// exact length, using the production predecode scanner to locate the
// existing value's byte range so the replacement works regardless of the
// field's current encoded length. Every other field, and every byte outside
// the target value, is left untouched.
func spliceTextFieldValue(t *testing.T, source []byte, targetKey uint64, newValue string) []byte {
	t.Helper()
	scanner := cborScanner{bytes: source, profile: executionPlanCBORProfile}
	major, count, err := scanner.readHead()
	if err != nil || major != 5 {
		t.Fatalf("fixture is not a top-level CBOR map: %v", err)
	}
	for index := uint64(0); index < count; index++ {
		keyStart := scanner.offset
		if err := scanner.scanItem(1); err != nil {
			t.Fatalf("scan key %d: %v", index, err)
		}
		keyReader := cborReader{bytes: source[keyStart:scanner.offset]}
		key, err := keyReader.unsigned()
		if err != nil {
			t.Fatalf("decode key %d: %v", index, err)
		}
		valueStart := scanner.offset
		if err := scanner.scanItem(1); err != nil {
			t.Fatalf("scan value for key %d: %v", key, err)
		}
		if key != targetKey {
			continue
		}
		result := make([]byte, 0, len(source)-(scanner.offset-valueStart)+len(newValue)+4)
		result = append(result, source[:valueStart]...)
		result = appendTextValue(result, newValue)
		result = append(result, source[scanner.offset:]...)
		return result
	}
	t.Fatalf("target key %d not found in fixture", targetKey)
	return nil
}

// appendTextValue encodes a CBOR major-type-3 text string using the same
// minimal/preferred length encoding the predecode scanner requires.
func appendTextValue(destination []byte, value string) []byte {
	length := len(value)
	switch {
	case length < 24:
		destination = append(destination, 0x60|byte(length))
	case length <= 0xff:
		destination = append(destination, 0x78, byte(length))
	default:
		destination = append(destination, 0x79, byte(length>>8), byte(length))
	}
	return append(destination, value...)
}
