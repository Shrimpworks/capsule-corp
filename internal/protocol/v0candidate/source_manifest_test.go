package v0candidate

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type mjsConformanceManifest struct {
	Cases []mjsConformanceCase `json:"cases"`
}

type mjsConformanceCase struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	WireFormat string `json:"wireFormat"`
	Fixture    struct {
		Path string `json:"path"`
	} `json:"fixture"`
	Context struct {
		Kind   string `json:"kind"`
		Source struct {
			Path string `json:"path"`
		} `json:"source"`
	} `json:"context"`
	Expected struct {
		Decision       string         `json:"decision"`
		Classification Classification `json:"classification"`
	} `json:"expected"`
	Implementations struct {
		Go string `json:"go"`
	} `json:"implementations"`
}

func TestMJSFoundationConformance(t *testing.T) {
	t.Parallel()
	manifestBytes := readConformanceFixture(t, "manifest.json")
	var manifest mjsConformanceManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatalf("decode conformance manifest: %v", err)
	}

	cases := make([]mjsConformanceCase, 0, 40)
	for _, candidate := range manifest.Cases {
		if candidate.Implementations.Go == "verified" &&
			(candidate.Object == "MjsSource" || candidate.Object == "SourceManifest") {
			cases = append(cases, candidate)
		}
	}
	if len(cases) != 40 {
		t.Fatalf("verified Go M1 foundation cases = %d, want 40", len(cases))
	}

	for _, candidate := range cases {
		candidate := candidate
		t.Run(candidate.ID, func(t *testing.T) {
			t.Parallel()
			fixture := readConformanceFixture(t, candidate.Fixture.Path)
			var err error
			switch {
			case candidate.WireFormat == "media-type" && candidate.ID == "mjs-source.closed-identifiers.profile":
				err = ValidateMJSSourceProfile(string(fixture))
			case candidate.WireFormat == "media-type" && candidate.ID == "mjs-source.closed-identifiers.profile-reject":
				err = ValidateMJSSourceProfile(string(fixture))
			case candidate.WireFormat == "media-type" && candidate.Object == "MjsSource":
				err = ValidateMJSSourceMediaType(string(fixture))
			case candidate.WireFormat == "media-type":
				err = ValidateSourceManifestMediaType(string(fixture))
			case candidate.Object == "MjsSource":
				_, err = NewMJSMainSource(fixture)
			case candidate.Object == "SourceManifest":
				sourcePath := candidate.Context.Source.Path
				if sourcePath == "" {
					sourcePath = "mjs-source/ordinary.mjs"
				}
				_, err = DecodeSourceManifest(
					fixture,
					SourceManifestMediaType,
					readConformanceFixture(t, sourcePath),
				)
			default:
				t.Fatalf("unsupported M1 conformance case shape")
			}

			if candidate.Expected.Decision == "accept" {
				assertClassification(t, err, "")
				return
			}
			assertClassification(t, err, candidate.Expected.Classification)
		})
	}
}

func TestMJSFoundationKnownAnswersAndDefensiveCopies(t *testing.T) {
	t.Parallel()
	answers := []struct {
		path   string
		length int
		digest string
	}{
		{"mjs-source/ordinary.mjs", 57, "2412c601eb0977eac9a5d4d30f459f48ab787c3f80551e080c38103543d89d75"},
		{"mjs-source-manifest/ordinary.cbor", 89, "052dce0c353e1efeb70f93405a8757ef6fa4d29f91a4d6bcaa67d00c45abc0d6"},
		{"mjs-source-manifest/minimum.cbor", 87, "b29a880488898dbac54c5890ec56b97a8aca461cbf6585af60fa3a56bc9e9044"},
		{"mjs-source-manifest/exact-maximum.cbor", 95, "21b9b521ed0ab6973be563b14f1eae4591097423e17d4d00ba0025c3fe4298d8"},
	}
	for _, answer := range answers {
		value := readConformanceFixture(t, answer.path)
		if len(value) != answer.length {
			t.Fatalf("%s length = %d, want %d", answer.path, len(value), answer.length)
		}
		digest := sha256.Sum256(value)
		if hex.EncodeToString(digest[:]) != answer.digest {
			t.Fatalf("%s digest = %x, want %s", answer.path, digest, answer.digest)
		}
	}
	for _, name := range []string{"minimum", "ordinary", "exact-maximum"} {
		source, err := NewMJSMainSource(readConformanceFixture(t, "mjs-source/"+name+".mjs"))
		if err != nil {
			t.Fatalf("retain %s source: %v", name, err)
		}
		encoded, err := EncodeSourceManifest(source)
		if err != nil {
			t.Fatalf("encode %s manifest: %v", name, err)
		}
		if !bytes.Equal(encoded, readConformanceFixture(t, "mjs-source-manifest/"+name+".cbor")) {
			t.Fatalf("Go %s SourceManifest encoding disagrees with the generated known answer", name)
		}
	}

	callerSource := readConformanceFixture(t, "mjs-source/ordinary.mjs")
	source, err := NewMJSMainSource(callerSource)
	if err != nil {
		t.Fatalf("retain source: %v", err)
	}
	encoded, err := EncodeSourceManifest(source)
	if err != nil {
		t.Fatalf("encode manifest: %v", err)
	}
	expectedManifest := readConformanceFixture(t, "mjs-source-manifest/ordinary.cbor")
	if !bytes.Equal(encoded, expectedManifest) {
		t.Fatal("Go SourceManifest encoding disagrees with the generated known answer")
	}

	callerManifest := bytes.Clone(expectedManifest)
	decoded, err := DecodeSourceManifest(callerManifest, SourceManifestMediaType, callerSource)
	if err != nil {
		t.Fatalf("decode manifest: %v", err)
	}
	callerSource[0] ^= 0xff
	callerManifest[0] ^= 0xff
	returnedSource := decoded.SourceBytes()
	returnedManifest := decoded.AuthoritativeBytes()
	returnedSource[0] ^= 0xff
	returnedManifest[0] ^= 0xff
	if !bytes.Equal(decoded.SourceBytes(), readConformanceFixture(t, "mjs-source/ordinary.mjs")) {
		t.Fatal("source defensive copy changed")
	}
	if !bytes.Equal(decoded.AuthoritativeBytes(), expectedManifest) {
		t.Fatal("manifest defensive copy changed")
	}
	view := decoded.View()
	view.Members[0].LogicalPath = "changed.mjs"
	if decoded.View().Members[0].LogicalPath != MJSMainPath {
		t.Fatal("manifest view retained caller-mutable member storage")
	}
}

// TestDecodeSourceManifestRejectsOversizedInputBeforeCloning proves the
// cap-before-copy ordering: an oversized received buffer must be rejected by
// its raw length before DecodeSourceManifest clones it. A regression that
// clones first (as this function once did) allocates bytes proportional to
// the attacker-controlled input size before the length check ever runs.
func TestDecodeSourceManifestRejectsOversizedInputBeforeCloning(t *testing.T) {
	source := readConformanceFixture(t, "mjs-source/ordinary.mjs")
	oversized := make([]byte, 4<<20) // 4 MiB: far beyond SourceManifestMaxCBORBytes (95).

	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	_, err := DecodeSourceManifest(oversized, SourceManifestMediaType, source)
	runtime.ReadMemStats(&after)

	if err == nil {
		t.Fatal("expected rejection of an oversized SourceManifest")
	}
	if classification, ok := ErrorClassification(err); !ok || classification != ClassificationMalformed {
		t.Fatalf("classification = %v (recognized %v), want MALFORMED: %v", classification, ok, err)
	}
	const allocationBudget = 64 * 1024 // bytes; far below the 4 MiB input, comfortably above legitimate small work.
	if delta := after.TotalAlloc - before.TotalAlloc; delta > allocationBudget {
		t.Fatalf("oversized input allocated %d bytes before length rejection (budget %d); cap-before-copy regressed", delta, allocationBudget)
	}
}

func TestMJSBlockerEvidenceExactBytes(t *testing.T) {
	t.Parallel()
	path := filepath.Join("..", "..", "..", "schemas", "conformance", "v0", "mjs-source", "language-hold-division-regexp-counterexample.mjs")
	value, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read blocker fixture: %v", err)
	}
	expected := []byte("const of = 9; of / import(\"evil\") / divisor;\n")
	if !bytes.Equal(value, expected) {
		t.Fatalf("blocker fixture bytes changed: %x", value)
	}
}
