// These tests exercise cborscan directly rather than through a consumer's
// fixture corpus. The package previously had no test file of its own: it was
// covered only transitively by v0candidate and bootstrapauthpassive, whose
// object models happen not to contain several of the malformed shapes this
// scanner exists to refuse. Measured with -coverpkg, five refusal branches
// were never taken by any test in the repository — empty payload, unsafe
// negative integer, truncated string, decoded-string budget, and truncated
// argument — and (*Error).Error was never called. Every Reason the package
// can emit now has a named case here.
package cborscan

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"testing"
)

// permissiveProfile is a deliberately loose profile: each test narrows only
// the one limit it is proving, so a refusal cannot come from an unrelated cap.
func permissiveProfile() Profile {
	return Profile{
		MaxBytes:         1024,
		MaxDepth:         8,
		MaxItems:         64,
		MaxMapEntries:    16,
		MaxArrayElements: 16,
		AllowArray:       true,
	}
}

// uint64Head encodes one CBOR head with the 8-byte argument form, which is
// the only form that can carry an argument above the safe-integer bound.
func uint64Head(major byte, argument uint64) []byte {
	head := []byte{major<<5 | 27, 0, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint64(head[1:], argument)
	return head
}

// assertReason runs Predecode and requires it to refuse with exactly want.
func assertReason(t *testing.T, received []byte, profile Profile, want Reason) {
	t.Helper()
	err := Predecode(received, profile)
	if err == nil {
		t.Fatalf("Predecode(% x) accepted the payload; want refusal reason %d", received, want)
	}
	var scanErr *Error
	if !errors.As(err, &scanErr) {
		t.Fatalf("Predecode(% x) returned %T; want *cborscan.Error", received, err)
	}
	if scanErr.Reason != want {
		t.Fatalf("Predecode(% x) refused with reason %d; want %d", received, scanErr.Reason, want)
	}
}

// assertAccepted runs Predecode and requires it to accept.
func assertAccepted(t *testing.T, received []byte, profile Profile) {
	t.Helper()
	if err := Predecode(received, profile); err != nil {
		t.Fatalf("Predecode(% x) refused with %v; want acceptance", received, err)
	}
}

func TestPredecodeRefusesEveryMalformedShape(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		received []byte
		profile  Profile
		reason   Reason
	}{
		// Predecode's own two entry checks, both of which every current
		// consumer performs itself before calling in, so neither was reached.
		{"empty payload", []byte{}, permissiveProfile(), ReasonEmptyPayload},
		{"over raw byte limit", []byte{0x43, 0x61, 0x62, 0x63}, Profile{MaxBytes: 3, MaxDepth: 8, MaxItems: 64}, ReasonRawByteLimit},
		{"trailing data", []byte{0x01, 0x02}, permissiveProfile(), ReasonTrailingData},

		// Head parsing.
		{"truncated item", []byte{0x81}, permissiveProfile(), ReasonTruncatedItem},
		{"truncated argument", []byte{0x19, 0x00}, permissiveProfile(), ReasonTruncatedArgument},
		{"nonpreferred argument", []byte{0x18, 0x17}, permissiveProfile(), ReasonNonpreferredArgument},
		{"indefinite length", []byte{0x9f}, permissiveProfile(), ReasonIndefiniteOrReserved},
		{"reserved additional", []byte{0x1c}, permissiveProfile(), ReasonIndefiniteOrReserved},

		// Integer bounds. The negative case is the one no consumer fixture
		// reaches: -9007199254740992 is one past the safe-integer floor.
		{"unsigned past UInt53", uint64Head(0, maxSafeInteger+1), permissiveProfile(), ReasonUnsafeInteger},
		{"negative past safe floor", uint64Head(1, maxSafeInteger), permissiveProfile(), ReasonUnsafeNegativeInteger},

		// Strings. "truncated string" declares five bytes with two present.
		{"truncated string", []byte{0x45, 0x61, 0x62}, permissiveProfile(), ReasonTruncatedString},
		{"invalid utf8 text", []byte{0x61, 0xff}, permissiveProfile(), ReasonInvalidUTF8},

		// Structural limits, each at the smallest refused value.
		{"depth limit", []byte{0x81, 0x81, 0x01}, narrowed(func(p *Profile) { p.MaxDepth = 2 }), ReasonDepthLimit},
		{"item limit", []byte{0x82, 0x01, 0x02}, narrowed(func(p *Profile) { p.MaxItems = 2 }), ReasonItemLimit},
		{"array element limit", []byte{0x82, 0x01, 0x02}, narrowed(func(p *Profile) { p.MaxArrayElements = 1 }), ReasonArrayElementLimit},
		{"map entry limit", []byte{0xa2, 0x01, 0x01, 0x02, 0x02}, narrowed(func(p *Profile) { p.MaxMapEntries = 1 }), ReasonMapEntryLimit},

		// Canonical map-key ordering, both directions of CompareDeterministic.
		{"duplicate map key", []byte{0xa2, 0x01, 0x01, 0x01, 0x02}, permissiveProfile(), ReasonDuplicateMapKeyOrder},
		{"map keys out of byte order", []byte{0xa2, 0x02, 0x01, 0x01, 0x02}, permissiveProfile(), ReasonDuplicateMapKeyOrder},
		{"longer map key before shorter", []byte{0xa2, 0x62, 0x61, 0x61, 0x01, 0x61, 0x61, 0x02}, permissiveProfile(), ReasonDuplicateMapKeyOrder},
		// Key 24 encodes as two bytes starting 0x18; the empty array encodes
		// as the single byte 0x80. Length-first ordering refuses this pair,
		// but a byte-wise-only comparison would accept it, so this is the map
		// case that actually proves the scanner applies RFC 8949 4.2.1.
		{"longer map key before byte-wise larger shorter key", []byte{0xa2, 0x18, 0x18, 0x01, 0x80, 0x02}, permissiveProfile(), ReasonDuplicateMapKeyOrder},

		// Tags and simple values, refused both when the profile disallows the
		// class outright and when it allows a different member of it.
		{"tag when disallowed", []byte{0xd2, 0x01}, permissiveProfile(), ReasonSemanticTag},
		{"wrong tag number", []byte{0xd2, 0x01}, narrowed(func(p *Profile) { p.AllowTag = true; p.AllowedTag = 17 }), ReasonSemanticTag},
		// Arrays are the third opt-in class. MaxArrayElements cannot express
		// "no arrays": a cap of 0 still admits the empty array, so an object
		// set that carries none has to leave AllowArray unset. Both an empty
		// and a populated array must be refused, and as a type refusal rather
		// than an element-count one.
		{"empty array when disallowed", []byte{0x80}, narrowed(func(p *Profile) { p.AllowArray = false }), ReasonMajorType},
		{"populated array when disallowed", []byte{0x81, 0x01}, narrowed(func(p *Profile) { p.AllowArray = false }), ReasonMajorType},
		{"array under a zero element cap is still refused by AllowArray", []byte{0x80}, narrowed(func(p *Profile) { p.AllowArray = false; p.MaxArrayElements = 0 }), ReasonMajorType},

		{"simple value when disallowed", []byte{0xf4}, permissiveProfile(), ReasonSimpleOrFloat},
		{"true when only false allowed", []byte{0xf5}, narrowed(func(p *Profile) { p.AllowFalse = true }), ReasonSimpleOrFloat},
		{"float when only false allowed", []byte{0xfb, 0x40, 0x09, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18}, narrowed(func(p *Profile) { p.AllowFalse = true }), ReasonSimpleOrFloat},

		// A refusal raised inside a nested item must propagate out of its
		// container rather than being swallowed by the enclosing loop. Each
		// of these nests the same nonpreferred-argument head in a different
		// recursion site: array element, map key, map value, tag payload.
		{"malformed array element", []byte{0x81, 0x18, 0x17}, permissiveProfile(), ReasonNonpreferredArgument},
		{"malformed map key", []byte{0xa1, 0x18, 0x17, 0x01}, permissiveProfile(), ReasonNonpreferredArgument},
		{"malformed map value", []byte{0xa1, 0x01, 0x18, 0x17}, permissiveProfile(), ReasonNonpreferredArgument},
		{"malformed tag payload", []byte{0xd2, 0x18, 0x17}, narrowed(func(p *Profile) { p.AllowTag = true; p.AllowedTag = 18 }), ReasonNonpreferredArgument},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assertReason(t, test.received, test.profile, test.reason)
		})
	}
}

func TestPredecodeAcceptsProfileAdmittedShapes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		received []byte
		profile  Profile
	}{
		{"unsigned at UInt53 ceiling", uint64Head(0, maxSafeInteger), permissiveProfile()},
		{"negative at safe floor", uint64Head(1, maxSafeInteger-1), permissiveProfile()},
		{"empty array", []byte{0x80}, permissiveProfile()},
		{"empty map", []byte{0xa0}, permissiveProfile()},
		{"canonically ordered map", []byte{0xa2, 0x01, 0x01, 0x02, 0x02}, permissiveProfile()},
		{"shorter map key before longer", []byte{0xa2, 0x61, 0x61, 0x01, 0x62, 0x61, 0x61, 0x02}, permissiveProfile()},
		{"text string", []byte{0x63, 0x61, 0x62, 0x63}, permissiveProfile()},
		// The accepted halves of the tag and simple-value branches. No
		// consumer profile reaches these from both sides, so without them the
		// AllowTag/AllowFalse switches are only ever proven to refuse.
		{"allowed tag", []byte{0xd2, 0x01}, narrowed(func(p *Profile) { p.AllowTag = true; p.AllowedTag = 18 })},
		{"allowed false", []byte{0xf4}, narrowed(func(p *Profile) { p.AllowFalse = true })},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assertAccepted(t, test.received, test.profile)
		})
	}
}

// TestStructuralLimitsAcceptTheirExactCeiling pairs each refusal above with
// the largest value the same profile must still accept, so an off-by-one in
// either direction fails rather than only an overshoot.
func TestStructuralLimitsAcceptTheirExactCeiling(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		received []byte
		profile  Profile
	}{
		{"depth at ceiling", []byte{0x81, 0x01}, narrowed(func(p *Profile) { p.MaxDepth = 2 })},
		{"items at ceiling", []byte{0x81, 0x01}, narrowed(func(p *Profile) { p.MaxItems = 2 })},
		{"array elements at ceiling", []byte{0x81, 0x01}, narrowed(func(p *Profile) { p.MaxArrayElements = 1 })},
		{"map entries at ceiling", []byte{0xa1, 0x01, 0x01}, narrowed(func(p *Profile) { p.MaxMapEntries = 1 })},
		{"raw bytes at ceiling", []byte{0x43, 0x61, 0x62, 0x63}, Profile{MaxBytes: 4, MaxDepth: 8, MaxItems: 64}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assertAccepted(t, test.received, test.profile)
		})
	}
}

// TestScannerRefusesDecodedStringBudgetOverrun covers the one refusal that
// Predecode cannot reach: it rejects any payload larger than Profile.MaxBytes
// up front, so the running decoded-string total can never exceed that same
// bound through the entry point. Driving Scanner directly is the only way to
// prove the check fires. Note the arguments here stay at or below MaxBytes;
// issue #346 covers the separate case where an argument exceeds MaxBytes.
func TestScannerRefusesDecodedStringBudgetOverrun(t *testing.T) {
	t.Parallel()
	profile := permissiveProfile()
	profile.MaxBytes = 4

	// array(2) of two 3-byte strings: the first fits the budget, the second
	// pushes the running total past it.
	value := []byte{0x82, 0x43, 0x61, 0x62, 0x63, 0x43, 0x64, 0x65, 0x66}
	scanner := Scanner{Value: value, Profile: profile}

	err := scanner.Item(1)
	var scanErr *Error
	if !errors.As(err, &scanErr) {
		t.Fatalf("Item returned %v (%T); want *cborscan.Error", err, err)
	}
	if scanErr.Reason != ReasonDecodedStringLimit {
		t.Fatalf("Item refused with reason %d; want ReasonDecodedStringLimit (%d)", scanErr.Reason, ReasonDecodedStringLimit)
	}
}

func TestCompareDeterministicOrdersByLengthThenBytes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		left  []byte
		right []byte
		want  int
	}{
		{"shorter sorts first", []byte{0x01}, []byte{0x01, 0x02}, -1},
		{"longer sorts last", []byte{0x01, 0x02}, []byte{0x01}, 1},
		{"equal", []byte{0x01, 0x02}, []byte{0x01, 0x02}, 0},
		{"same length lower byte first", []byte{0x01, 0x01}, []byte{0x01, 0x02}, -1},
		{"same length higher byte last", []byte{0x01, 0x02}, []byte{0x01, 0x01}, 1},
		// RFC 8949 4.2.1 orders by encoded length before bytes. The two cases
		// above cannot tell the two rules apart, because a common prefix makes
		// bytes.Compare agree with length. These invert that: the longer
		// operand is byte-wise smaller, so dropping the length rule flips the
		// sign rather than leaving it unchanged.
		{"longer still sorts last when byte-wise smaller", []byte{0x18, 0x18}, []byte{0x80}, 1},
		{"shorter still sorts first when byte-wise larger", []byte{0x80}, []byte{0x18, 0x18}, -1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := CompareDeterministic(test.left, test.right); got != test.want {
				t.Fatalf("CompareDeterministic(% x, % x) = %d; want %d", test.left, test.right, got, test.want)
			}
		})
	}
}

func TestReadHeadReportsConsumedOffset(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		value        []byte
		wantMajor    byte
		wantArgument uint64
		wantNext     int
	}{
		{"immediate argument", []byte{0x0a}, 0, 10, 1},
		{"one byte argument", []byte{0x18, 0x18}, 0, 24, 2},
		{"two byte argument", []byte{0x19, 0x01, 0x00}, 0, 256, 3},
		{"four byte argument", []byte{0x1a, 0x00, 0x01, 0x00, 0x00}, 0, 65536, 5},
		{"eight byte argument", uint64Head(0, 1<<32), 0, 1 << 32, 9},
		{"major type is the top three bits", []byte{0xa1}, 5, 1, 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			major, argument, next, err := ReadHead(test.value, 0)
			if err != nil {
				t.Fatalf("ReadHead(% x) returned %v; want success", test.value, err)
			}
			if major != test.wantMajor || argument != test.wantArgument || next != test.wantNext {
				t.Fatalf("ReadHead(% x) = (%d, %d, %d); want (%d, %d, %d)",
					test.value, major, argument, next, test.wantMajor, test.wantArgument, test.wantNext)
			}
		})
	}
}

// TestErrorMessageNamesItsReason pins the refusal string. It was one constant
// shared by every Reason, so a log line or a test failure holding only the
// error could not say which check fired.
func TestErrorMessageNamesItsReason(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		received []byte
		want     string
	}{
		{"empty payload", []byte{}, "cborscan: predecode refused: empty-payload"},
		{"trailing data", []byte{0x01, 0x02}, "cborscan: predecode refused: trailing-data"},
		{"invalid utf8", []byte{0x61, 0xff}, "cborscan: predecode refused: invalid-utf8"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := Predecode(test.received, permissiveProfile())
			if err == nil {
				t.Fatalf("Predecode(% x) accepted the payload; want refusal", test.received)
			}
			if got := err.Error(); got != test.want {
				t.Fatalf("Error() = %q; want %q", got, test.want)
			}
		})
	}
}

// TestReasonNamesCoverEveryReason fails when a Reason is added to the const
// block without a matching entry in reasonNames, so a new refusal cannot
// silently format as "unknown(N)".
func TestReasonNamesCoverEveryReason(t *testing.T) {
	t.Parallel()
	// ReasonNonpreferredArgument is the last value in the iota block; if a
	// Reason is appended after it, this bound moves with it.
	for reason := ReasonEmptyPayload; reason <= ReasonNonpreferredArgument; reason++ {
		name := reason.String()
		if name == "" || name == fmt.Sprintf("unknown(%d)", int(reason)) {
			t.Errorf("Reason(%d) has no name in reasonNames", int(reason))
		}
	}
	if got := Reason(-1).String(); got != "unknown(-1)" {
		t.Errorf("Reason(-1).String() = %q; want %q", got, "unknown(-1)")
	}
	if got := Reason(len(reasonNames)).String(); got != fmt.Sprintf("unknown(%d)", len(reasonNames)) {
		t.Errorf("out-of-range Reason formatted as %q; want an unknown(N) form", got)
	}
}

// TestConsumerProfilesRoundTrip guards the two shapes the real consumers rely
// on, so a change made to satisfy the cases above cannot quietly break them:
// bootstrapauthpassive wraps a COSE_Sign1 tag around an array, and
// v0candidate refuses tags entirely.
func TestConsumerProfilesRoundTrip(t *testing.T) {
	t.Parallel()

	// tag(18) wrapping array(2) of two byte strings, the COSE_Sign1 shape.
	signed := []byte{0xd2, 0x82, 0x41, 0x01, 0x41, 0x02}
	tagging := permissiveProfile()
	tagging.AllowTag = true
	tagging.AllowedTag = 18
	assertAccepted(t, signed, tagging)

	// The same bytes under a profile that admits no tag at all.
	assertReason(t, signed, permissiveProfile(), ReasonSemanticTag)

	// A scanned payload must consume exactly its own bytes.
	scanner := Scanner{Value: signed, Profile: tagging}
	if err := scanner.Item(1); err != nil {
		t.Fatalf("Item returned %v; want success", err)
	}
	if scanner.Offset != len(signed) {
		t.Fatalf("Offset = %d after scanning %d bytes; want full consumption", scanner.Offset, len(signed))
	}
	if !bytes.Equal(scanner.Value, signed) {
		t.Fatal("Scanner mutated the payload it was given")
	}
}

// narrowed returns permissiveProfile with one limit tightened, so each table
// row states only the constraint it is proving.
func narrowed(adjust func(*Profile)) Profile {
	profile := permissiveProfile()
	adjust(&profile)
	return profile
}
