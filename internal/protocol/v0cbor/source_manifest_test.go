package v0cbor

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"capsule.local/capsule/internal/protocol/v0candidate"
	"github.com/fxamacker/cbor/v2"
)

type conformanceManifest struct {
	Cases []conformanceCase `json:"cases"`
}

type conformanceCase struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	WireFormat string `json:"wireFormat"`
	Fixture    struct {
		Path string `json:"path"`
	} `json:"fixture"`
	Context struct {
		Source struct {
			Path string `json:"path"`
		} `json:"source"`
	} `json:"context"`
	Expected struct {
		Decision       string                `json:"decision"`
		Classification FailureClassification `json:"classification"`
	} `json:"expected"`
	Implementations struct {
		Go         string `json:"go"`
		TypeScript string `json:"typescript"`
	} `json:"implementations"`
}

func TestSourceManifestCompleteApplicableCorpus(t *testing.T) {
	t.Parallel()
	codec := mustCodec(t)
	var manifest conformanceManifest
	if err := json.Unmarshal(readFixture(t, "manifest.json"), &manifest); err != nil {
		t.Fatalf("decode conformance manifest: %v", err)
	}

	cases := make([]conformanceCase, 0, 19)
	for _, candidate := range manifest.Cases {
		if candidate.Object == "SourceManifest" &&
			candidate.Implementations.Go == "verified" &&
			candidate.Implementations.TypeScript == "verified" {
			cases = append(cases, candidate)
		}
	}
	if len(cases) != 19 {
		t.Fatalf("applicable Go/TypeScript SourceManifest cases = %d, want 19", len(cases))
	}

	for _, candidate := range cases {
		candidate := candidate
		t.Run(candidate.ID, func(t *testing.T) {
			t.Parallel()
			fixture := readFixture(t, candidate.Fixture.Path)
			var err error
			if candidate.WireFormat == "media-type" {
				err = v0candidate.ValidateSourceManifestMediaType(string(fixture))
			} else {
				sourcePath := candidate.Context.Source.Path
				if sourcePath == "" {
					sourcePath = "mjs-source/ordinary.mjs"
				}
				_, err = codec.Decode(fixture, v0candidate.SourceManifestMediaType, readFixture(t, sourcePath))
			}
			assertClassification(t, err, candidate.Expected.Decision, candidate.Expected.Classification)
		})
	}
}

func TestSourceManifestCrossImplementationKnownAnswers(t *testing.T) {
	t.Parallel()
	codec := mustCodec(t)
	answers := []struct {
		name   string
		length int
		digest string
	}{
		{name: "minimum", length: 87, digest: "b29a880488898dbac54c5890ec56b97a8aca461cbf6585af60fa3a56bc9e9044"},
		{name: "ordinary", length: 89, digest: "052dce0c353e1efeb70f93405a8757ef6fa4d29f91a4d6bcaa67d00c45abc0d6"},
		{name: "exact-maximum", length: 95, digest: "21b9b521ed0ab6973be563b14f1eae4591097423e17d4d00ba0025c3fe4298d8"},
	}
	for _, answer := range answers {
		answer := answer
		t.Run(answer.name, func(t *testing.T) {
			t.Parallel()
			source := readFixture(t, "mjs-source/"+answer.name+".mjs")
			expected := readFixture(t, "mjs-source-manifest/"+answer.name+".cbor")
			encoded, err := codec.Encode(source)
			if err != nil {
				t.Fatalf("encode: %v", err)
			}
			if !bytes.Equal(encoded, expected) {
				t.Fatal("fxamacker encoding disagrees with retained Go/TypeScript known answer")
			}
			if len(encoded) != answer.length {
				t.Fatalf("encoded length = %d, want %d", len(encoded), answer.length)
			}
			digest := sha256.Sum256(encoded)
			if hex.EncodeToString(digest[:]) != answer.digest {
				t.Fatalf("encoded digest = %x, want %s", digest, answer.digest)
			}
			legacySource, legacyErr := v0candidate.NewMJSMainSource(source)
			if legacyErr != nil {
				t.Fatalf("handwritten source oracle: %v", legacyErr)
			}
			legacy, legacyErr := v0candidate.EncodeSourceManifest(legacySource)
			if legacyErr != nil || !bytes.Equal(encoded, legacy) {
				t.Fatalf("handwritten encoding oracle disagreement: %v", legacyErr)
			}
		})
	}
}

func TestSourceManifestDefensiveOwnership(t *testing.T) {
	t.Parallel()
	codec := mustCodec(t)
	source := readFixture(t, "mjs-source/ordinary.mjs")
	manifest := readFixture(t, "mjs-source-manifest/ordinary.cbor")
	expectedSource := bytes.Clone(source)
	expectedManifest := bytes.Clone(manifest)

	decoded, err := codec.Decode(manifest, v0candidate.SourceManifestMediaType, source)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	source[0] ^= 0xff
	manifest[0] ^= 0xff
	returnedSource := decoded.SourceBytes()
	returnedManifest := decoded.ExactBytes()
	returnedSource[0] ^= 0xff
	returnedManifest[0] ^= 0xff
	if !bytes.Equal(decoded.SourceBytes(), expectedSource) {
		t.Fatal("source bytes retained caller- or accessor-mutable storage")
	}
	if !bytes.Equal(decoded.ExactBytes(), expectedManifest) {
		t.Fatal("manifest bytes retained caller- or accessor-mutable storage")
	}
	if decoded.View().Member.LogicalPath != v0candidate.MJSMainPath {
		t.Fatal("typed projection changed")
	}
	if decoded.Digest() != v0candidate.SourceManifestDigest(sha256.Sum256(expectedManifest)) {
		t.Fatal("digest is not over exact received bytes")
	}
}

func TestSourceManifestRestorationCorpus(t *testing.T) {
	t.Parallel()
	codec := mustCodec(t)
	source := readFixture(t, "mjs-source/ordinary.mjs")
	ordinary := readFixture(t, "mjs-source-manifest/ordinary.cbor")

	nonpreferred := replaceOnce(t, ordinary, []byte{0x02, 0x00, 0x03}, []byte{0x02, 0x18, 0x00, 0x03})
	invalidUTF8 := bytes.Clone(ordinary)
	invalidUTF8[bytes.Index(invalidUTF8, []byte(v0candidate.SourceManifestObjectType))] = 0xff
	duplicate := append(bytes.Clone(ordinary), 0x05, 0x18, 0x39)
	duplicate[0] = 0xa6
	unknown := replaceOnce(t, ordinary, []byte{0x05, 0x18, 0x39}, []byte{0x06, 0x18, 0x39})
	depth := replaceOnce(t, ordinary, []byte{0x05, 0x18, 0x39}, []byte{0x05, 0x81, 0x81, 0x81, 0x81, 0x00})
	collection := replaceOnce(t, ordinary, []byte{0x04, 0x81, 0x83, 0x68}, []byte{0x04, 0x84, 0x80, 0x80, 0x80, 0x80})
	tagged := append([]byte{0xc0}, ordinary...)
	indefinite := append(bytes.Clone(ordinary), 0xff)
	indefinite[0] = 0xbf
	floating := replaceOnce(t, ordinary, []byte{0x05, 0x18, 0x39}, []byte{0x05, 0xf9, 0x00, 0x00})
	trailing := append(bytes.Clone(ordinary), 0x00)
	wrongType := replaceOnce(t, ordinary, []byte{0x02, 0x00, 0x03}, []byte{0x02, 0x60, 0x03})

	tests := []struct {
		name           string
		manifest       []byte
		classification FailureClassification
	}{
		{name: "duplicate field", manifest: duplicate, classification: FailureMalformed},
		{name: "retained out of order field", manifest: readFixture(t, "mjs-source-manifest/reject-map-order.cbor"), classification: FailureMalformed},
		{name: "nonpreferred integer", manifest: nonpreferred, classification: FailureMalformed},
		{name: "invalid UTF-8", manifest: invalidUTF8, classification: FailureMalformed},
		{name: "depth cap", manifest: depth, classification: FailureMalformed},
		{name: "collection cap", manifest: collection, classification: FailureMalformed},
		{name: "semantic tag", manifest: tagged, classification: FailureMalformed},
		{name: "indefinite map", manifest: indefinite, classification: FailureMalformed},
		{name: "floating point", manifest: floating, classification: FailureMalformed},
		{name: "trailing item", manifest: trailing, classification: FailureMalformed},
		{name: "unknown field", manifest: unknown, classification: FailureUnsupported},
		{name: "wrong typed field", manifest: wrongType, classification: FailureSchema},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			_, err := codec.Decode(test.manifest, v0candidate.SourceManifestMediaType, source)
			assertClassification(t, err, "reject", test.classification)
		})
	}
}

// TestDecodeRejectsOversizedInputBeforeCloning proves the cap-before-copy
// ordering: an oversized received buffer must be rejected by its raw length
// before Decode clones it. A regression that clones first (as this function
// once did) allocates bytes proportional to the attacker-controlled input
// size before the length check ever runs.
func TestDecodeRejectsOversizedInputBeforeCloning(t *testing.T) {
	codec := mustCodec(t)
	source := readFixture(t, "mjs-source/ordinary.mjs")
	oversized := make([]byte, 4<<20) // 4 MiB: far beyond SourceManifestMaxCBORBytes (95).

	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	_, err := codec.Decode(oversized, v0candidate.SourceManifestMediaType, source)
	runtime.ReadMemStats(&after)

	assertClassification(t, err, "reject", FailureMalformed)
	const allocationBudget = 64 * 1024 // bytes; far below the 4 MiB input, comfortably above legitimate small work.
	if delta := after.TotalAlloc - before.TotalAlloc; delta > allocationBudget {
		t.Fatalf("oversized input allocated %d bytes before length rejection (budget %d); cap-before-copy regressed", delta, allocationBudget)
	}
}

func TestCanonicalComparisonSurvivesPredecoderRestoration(t *testing.T) {
	t.Parallel()
	codec := mustCodec(t)
	ordinary := readFixture(t, "mjs-source-manifest/ordinary.cbor")
	nonpreferred := replaceOnce(t, ordinary, []byte{0x02, 0x00, 0x03}, []byte{0x02, 0x18, 0x00, 0x03})

	// The typed library accepts this valid but nonpreferred representation.  The
	// Capsule predecoder rejects it; this direct mode probe additionally proves
	// the exact-byte comparison would reject it if that predecoder rule regressed.
	var wire sourceManifestWire
	if err := codec.decode.Unmarshal(nonpreferred, &wire); err != nil {
		t.Fatalf("typed decoder no longer demonstrates the restoration hazard: %v", err)
	}
	canonical, err := codec.encode.Marshal(wire)
	if err != nil {
		t.Fatalf("canonical re-encode: %v", err)
	}
	if bytes.Equal(canonical, nonpreferred) || !bytes.Equal(canonical, ordinary) {
		t.Fatal("exact canonical comparison no longer distinguishes the restored representation")
	}
	if err := cbor.Unmarshal(nonpreferred, &wire); err != nil {
		t.Fatalf("generic library decoder no longer accepts the restoration probe: %v", err)
	}
	if _, err := codec.Decode(nonpreferred, v0candidate.SourceManifestMediaType, readFixture(t, "mjs-source/ordinary.mjs")); err == nil {
		t.Fatal("object-specific wrapper accepted restored nonpreferred encoding")
	}
}

func TestSourceManifestDeterministicPropertyCorpus(t *testing.T) {
	codec := mustCodec(t)
	state := uint64(0x6a09e667f3bcc909)
	for index := 0; index < 10_000; index++ {
		state ^= state << 13
		state ^= state >> 7
		state ^= state << 17
		length := int(state % 1025)
		source := make([]byte, length)
		for offset := range source {
			state = state*6364136223846793005 + 1
			source[offset] = byte(0x20 + state%0x5f)
		}
		encoded, err := codec.Encode(source)
		if err != nil {
			t.Fatalf("property case %d encode: %v", index, err)
		}
		decoded, err := codec.Decode(encoded, v0candidate.SourceManifestMediaType, source)
		if err != nil {
			t.Fatalf("property case %d decode: %v", index, err)
		}
		if !bytes.Equal(decoded.ExactBytes(), encoded) || !bytes.Equal(decoded.SourceBytes(), source) {
			t.Fatalf("property case %d lost exact bytes", index)
		}
		legacySource, legacyErr := v0candidate.NewMJSMainSource(source)
		if legacyErr != nil {
			t.Fatalf("property case %d handwritten source: %v", index, legacyErr)
		}
		legacy, legacyErr := v0candidate.EncodeSourceManifest(legacySource)
		if legacyErr != nil || !bytes.Equal(legacy, encoded) {
			t.Fatalf("property case %d handwritten oracle disagreement: %v", index, legacyErr)
		}
	}
}

func FuzzSourceManifestCodec(f *testing.F) {
	for _, name := range []string{"minimum", "ordinary", "exact-maximum"} {
		f.Add(
			readFixture(f, "mjs-source-manifest/"+name+".cbor"),
			readFixture(f, "mjs-source/"+name+".mjs"),
		)
	}
	f.Add(readFixture(f, "mjs-source-manifest/cap-plus-one.cbor"), readFixture(f, "mjs-source/exact-maximum.mjs"))

	codec, err := NewSourceManifestCodec()
	if err != nil {
		f.Fatalf("construct codec: %v", err)
	}
	f.Fuzz(func(t *testing.T, manifest, source []byte) {
		if len(manifest) > v0candidate.SourceManifestMaxCBORBytes+1 || len(source) > v0candidate.MJSMainSourceMaxBytes+1 {
			return
		}
		decoded, gotErr := codec.Decode(manifest, v0candidate.SourceManifestMediaType, source)
		oracle, oracleErr := v0candidate.DecodeSourceManifest(manifest, v0candidate.SourceManifestMediaType, source)
		if (gotErr == nil) != (oracleErr == nil) {
			t.Fatalf("wrapper/oracle acceptance disagreement: wrapper=%v oracle=%v", gotErr, oracleErr)
		}
		if gotErr == nil {
			if !bytes.Equal(decoded.ExactBytes(), oracle.AuthoritativeBytes()) ||
				!bytes.Equal(decoded.SourceBytes(), oracle.SourceBytes()) ||
				decoded.Digest() != oracle.Digest() {
				t.Fatal("accepted wrapper/oracle exact-byte disagreement")
			}
		}
	})
}

type testingTB interface {
	Helper()
	Fatalf(string, ...any)
}

func mustCodec(t testingTB) *SourceManifestCodec {
	t.Helper()
	codec, err := NewSourceManifestCodec()
	if err != nil {
		t.Fatalf("construct SourceManifest codec: %v", err)
	}
	return codec
}

func readFixture(t testingTB, relative string) []byte {
	t.Helper()
	path := filepath.Join("..", "..", "..", "schemas", "conformance", "v0", relative)
	value, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", relative, err)
	}
	return value
}

func assertClassification(t *testing.T, err error, decision string, expected FailureClassification) {
	t.Helper()
	if decision == "accept" {
		if err != nil {
			t.Fatalf("unexpected rejection: %v", err)
		}
		return
	}
	if err == nil {
		t.Fatalf("expected %s rejection", expected)
	}
	actual, ok := ErrorClassification(err)
	if !ok || actual != expected {
		t.Fatalf("classification = %q (recognized %v), want %q: %v", actual, ok, expected, err)
	}
}

func replaceOnce(t *testing.T, value, old, replacement []byte) []byte {
	t.Helper()
	index := bytes.Index(value, old)
	if index < 0 || bytes.Contains(value[index+len(old):], old) {
		t.Fatalf("mutation anchor %x is missing or ambiguous", old)
	}
	result := make([]byte, 0, len(value)-len(old)+len(replacement))
	result = append(result, value[:index]...)
	result = append(result, replacement...)
	result = append(result, value[index+len(old):]...)
	return result
}
