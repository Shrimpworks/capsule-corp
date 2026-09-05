// This file pins the exact refusal code predecodeApprovalCBOR emits for every
// malformed shape it can observe. Those codes are contract surface: they reach
// callers through ErrorClassification and appear in conformance fixtures, so a
// consolidation of the scan logic underneath must leave every one of them
// byte-identical. The table was captured against the package's own hand-rolled
// scanner before that scanner was replaced by internal/protocol/cborscan, and
// is the differential baseline for issue #346.
package approvalattempt

import (
	"errors"
	"testing"

	"capsule.local/capsule/internal/protocol/cborscan"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

// payloadProfile and protectedProfile mirror the two profiles the production
// verifier actually uses, at production_verifier.go's decodePayload and
// decodeProtected call sites.
func payloadProfile() approvalPredecodeProfile {
	return approvalPredecodeProfile{
		maxBytes:      ApprovalPayloadRawMaxBytes,
		maxDepth:      2,
		maxItems:      25,
		maxMapEntries: 12,
	}
}

func protectedProfile() approvalPredecodeProfile {
	return approvalPredecodeProfile{
		maxBytes:      ApprovalProtectedRawMaxBytes,
		maxDepth:      2,
		maxItems:      7,
		maxMapEntries: 3,
	}
}

// tightened returns payloadProfile with one bound narrowed, so a row states
// only the constraint it is proving.
func tightened(adjust func(*approvalPredecodeProfile)) approvalPredecodeProfile {
	profile := payloadProfile()
	adjust(&profile)
	return profile
}

// oversizedByteString builds a byte-string head declaring length bytes without
// the body, which is how the "argument exceeds the envelope ceiling" branch is
// reached without allocating that many bytes.
func oversizedByteString(length uint16) []byte {
	return []byte{0x59, byte(length >> 8), byte(length)}
}

func TestPredecodeApprovalCBORRefusalCodes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		value   []byte
		profile approvalPredecodeProfile
		code    string
	}{
		// Entry checks. Empty and over-limit deliberately share one code;
		// preserve that collapse rather than splitting it during any refactor.
		{"empty payload", []byte{}, payloadProfile(), "approval-cbor-byte-limit"},
		{"over byte limit", make([]byte, ApprovalPayloadRawMaxBytes+1), payloadProfile(), "approval-cbor-byte-limit"},
		{"trailing data", []byte{0x01, 0x02}, payloadProfile(), "approval-cbor-trailing"},

		// Structural bounds.
		{"depth limit", []byte{0xa1, 0x01, 0xa1, 0x01, 0x01}, payloadProfile(), "approval-cbor-depth"},
		{"item limit", []byte{0xa2, 0x01, 0x01, 0x02, 0x02}, tightened(func(p *approvalPredecodeProfile) { p.maxItems = 4 }), "approval-cbor-items"},
		{"map entry limit", []byte{0xa2, 0x01, 0x01, 0x02, 0x02}, tightened(func(p *approvalPredecodeProfile) { p.maxMapEntries = 1 }), "approval-map-entries"},

		// Integer bounds.
		{"unsafe unsigned integer", uint64Argument(0, v0candidate.MaxSafeInteger+1), payloadProfile(), "approval-unsafe-integer"},
		{"unsafe negative integer", uint64Argument(1, v0candidate.MaxSafeInteger), payloadProfile(), "approval-unsafe-negative-integer"},

		// Strings. Both the envelope-ceiling branch and the
		// truncation/budget branch report one shared code.
		{"byte string past envelope ceiling", oversizedByteString(ApprovalEnvelopeRawMaxBytes + 1), payloadProfile(), "approval-string-limit"},
		{"truncated byte string", []byte{0x45, 0x61, 0x62}, payloadProfile(), "approval-string-limit"},
		{"invalid utf8 text", []byte{0x61, 0xff}, payloadProfile(), "approval-invalid-utf8"},

		// Canonical map-key ordering.
		{"duplicate map key", []byte{0xa2, 0x01, 0x01, 0x01, 0x02}, payloadProfile(), "approval-map-key-order"},
		{"map keys out of order", []byte{0xa2, 0x02, 0x01, 0x01, 0x02}, payloadProfile(), "approval-map-key-order"},
		{"longer map key first", []byte{0xa2, 0x18, 0x18, 0x01, 0x40, 0x02}, payloadProfile(), "approval-map-key-order"},

		// Major types the approval object set does not contain. Arrays, tags,
		// and simple/float values all report the same type refusal, including
		// an empty array.
		{"empty array", []byte{0x80}, payloadProfile(), "approval-cbor-type"},
		{"populated array", []byte{0x81, 0x01}, payloadProfile(), "approval-cbor-type"},
		{"semantic tag", []byte{0xd2, 0x01}, payloadProfile(), "approval-cbor-type"},
		{"simple value false", []byte{0xf4}, payloadProfile(), "approval-cbor-type"},
		{"float", []byte{0xfb, 0x40, 0x09, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18}, payloadProfile(), "approval-cbor-type"},

		// Head parsing. A head cut short and an argument cut short share one
		// code, unlike cborscan, which distinguishes the two.
		{"truncated item", []byte{0xa1}, payloadProfile(), "approval-cbor-truncated"},
		{"truncated argument", []byte{0x19, 0x00}, payloadProfile(), "approval-cbor-truncated"},
		{"indefinite length", []byte{0xbf}, payloadProfile(), "approval-cbor-indefinite"},
		{"reserved additional", []byte{0x1c}, payloadProfile(), "approval-cbor-indefinite"},
		{"nonpreferred argument", []byte{0x18, 0x17}, payloadProfile(), "approval-cbor-nonpreferred"},

		// The protected-header profile is tighter; prove its bounds bind too.
		{"protected over byte limit", make([]byte, ApprovalProtectedRawMaxBytes+1), protectedProfile(), "approval-cbor-byte-limit"},
		{"protected map entry limit", []byte{0xa4, 0x01, 0x01, 0x02, 0x02, 0x03, 0x03, 0x04, 0x04}, protectedProfile(), "approval-map-entries"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := predecodeApprovalCBOR(test.value, test.profile)
			if err == nil {
				t.Fatalf("predecodeApprovalCBOR(% x) accepted the payload; want code %q", test.value, test.code)
			}
			var contract *contractError
			if !errors.As(err, &contract) {
				t.Fatalf("predecodeApprovalCBOR(% x) returned %T; want *contractError", test.value, err)
			}
			if contract.classification != ClassificationMalformed {
				t.Fatalf("classification = %q; want %q", contract.classification, ClassificationMalformed)
			}
			if contract.code != test.code {
				t.Fatalf("predecodeApprovalCBOR(% x) code = %q; want %q", test.value, contract.code, test.code)
			}
		})
	}
}

func TestPredecodeApprovalCBORAcceptsProfileShapes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		value   []byte
		profile approvalPredecodeProfile
	}{
		{"empty map", []byte{0xa0}, payloadProfile()},
		{"single entry map", []byte{0xa1, 0x01, 0x01}, payloadProfile()},
		{"canonically ordered map", []byte{0xa2, 0x01, 0x01, 0x02, 0x02}, payloadProfile()},
		{"byte string", []byte{0x43, 0x61, 0x62, 0x63}, payloadProfile()},
		{"text string", []byte{0x63, 0x61, 0x62, 0x63}, payloadProfile()},
		{"unsigned at UInt53 ceiling", uint64Argument(0, v0candidate.MaxSafeInteger), payloadProfile()},
		{"negative at safe floor", uint64Argument(1, v0candidate.MaxSafeInteger-1), payloadProfile()},
		{"map entries at ceiling", []byte{0xa1, 0x01, 0x01}, tightened(func(p *approvalPredecodeProfile) { p.maxMapEntries = 1 })},
		{"bytes at ceiling", []byte{0x43, 0x61, 0x62, 0x63}, tightened(func(p *approvalPredecodeProfile) { p.maxBytes = 4 })},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if err := predecodeApprovalCBOR(test.value, test.profile); err != nil {
				t.Fatalf("predecodeApprovalCBOR(% x) refused with %v; want acceptance", test.value, err)
			}
		})
	}
}

// TestApprovalPredecodeRefusalCodeCoversEveryReason fails if cborscan grows a
// Reason this package has no code for. Without it, a new refusal would fall
// through to approvalUnmappedPredecodeCode at runtime instead of failing the
// build or a test — the same silent-default problem the shared scanner's own
// reasonNames table guards against.
func TestApprovalPredecodeRefusalCodeCoversEveryReason(t *testing.T) {
	t.Parallel()
	for reason := cborscan.ReasonEmptyPayload; reason <= cborscan.ReasonNonpreferredArgument; reason++ {
		code := approvalPredecodeRefusalCode(&cborscan.Error{Reason: reason})
		if code == approvalUnmappedPredecodeCode {
			t.Errorf("cborscan reason %q (%d) has no approval refusal code", reason, int(reason))
		}
	}

	// A non-cborscan error must still classify rather than panic.
	if got := approvalPredecodeRefusalCode(errors.New("unrelated")); got != approvalUnmappedPredecodeCode {
		t.Errorf("approvalPredecodeRefusalCode(non-cborscan error) = %q; want %q", got, approvalUnmappedPredecodeCode)
	}
}

// TestApprovalProfileRefusesEveryOptionalMajorType pins the three cborscan
// opt-in switches as closed for this package. Setting any of them would widen
// the approval object set silently, since the refusal code for all three is
// the same "approval-cbor-type".
func TestApprovalProfileRefusesEveryOptionalMajorType(t *testing.T) {
	t.Parallel()
	shared := payloadProfile().sharedProfile()
	if shared.AllowArray {
		t.Error("approval profile sets AllowArray; the approval object set carries no array")
	}
	if shared.AllowTag {
		t.Error("approval profile sets AllowTag; the approval payload and protected header carry no tag")
	}
	if shared.AllowFalse {
		t.Error("approval profile sets AllowFalse; the approval object set carries no simple value")
	}
}

// uint64Argument encodes one CBOR head using the 8-byte argument form.
func uint64Argument(major byte, argument uint64) []byte {
	head := []byte{major<<5 | 27, 0, 0, 0, 0, 0, 0, 0, 0}
	for index := 0; index < 8; index++ {
		head[8-index] = byte(argument >> (8 * index))
	}
	return head
}
