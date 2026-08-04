package sourcevalidatorpassive

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type conformanceManifest struct {
	Cases []conformanceCase `json:"cases"`
}

type fixtureReference struct {
	Path string `json:"path"`
}

type conformanceCase struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Fixture fixtureReference `json:"fixture"`
	Context struct {
		Kind                 string            `json:"kind"`
		Source               *fixtureReference `json:"source"`
		Request              *fixtureReference `json:"request"`
		ArtifactProfile      *fixtureReference `json:"artifactProfile"`
		EngineeringCandidate *fixtureReference `json:"engineeringCandidate"`
		SourceManifest       *fixtureReference `json:"sourceManifest"`
	} `json:"context"`
	Expected struct {
		Decision       string `json:"decision"`
		Classification string `json:"classification"`
		Owner          string `json:"owner"`
		Effects        *struct {
			State       bool `json:"state"`
			Approval    bool `json:"approval"`
			Key         bool `json:"key"`
			IPCEndpoint bool `json:"ipcEndpoint"`
			Process     bool `json:"process"`
			Runtime     bool `json:"runtime"`
			Backend     bool `json:"backend"`
			Guest       bool `json:"guest"`
		} `json:"effects"`
	} `json:"expected"`
	Implementations struct {
		Go string `json:"go"`
	} `json:"implementations"`
}

func TestPassiveConformanceCorpus(t *testing.T) {
	t.Parallel()
	manifest := readManifest(t)
	cases := make([]conformanceCase, 0, 128)
	for _, candidate := range manifest.Cases {
		if candidate.Expected.Owner == "source-validator-passive-contract" && candidate.Implementations.Go == "verified" {
			cases = append(cases, candidate)
		}
	}
	if len(cases) != 128 {
		t.Fatalf("verified Go passive Source Validator cases = %d, want 128", len(cases))
	}

	for _, candidate := range cases {
		candidate := candidate
		t.Run(candidate.ID, func(t *testing.T) {
			t.Parallel()
			fixture := readFixture(t, candidate.Fixture.Path)
			var err error
			switch candidate.Object {
			case "SourceValidatorEngineeringCandidate":
				_, err = DecodeEngineeringCandidate(fixture)
			case "SourceValidatorArtifactProfile":
				_, err = DecodeArtifactProfile(fixture)
			case "SourceValidatorRequest":
				_, err = DecodeRequest(fixture)
			case "SourceValidatorResult":
				if candidate.Context.Request == nil || candidate.Context.ArtifactProfile == nil {
					t.Fatal("result case lacks request or artifact-profile context")
				}
				request, requestErr := DecodeRequest(readFixture(t, candidate.Context.Request.Path))
				if requestErr != nil {
					t.Fatalf("decode contextual request: %v", requestErr)
				}
				_, err = VerifyResult(request, readFixture(t, candidate.Context.ArtifactProfile.Path), fixture)
			default:
				t.Fatalf("unexpected passive object %q", candidate.Object)
			}

			if candidate.Expected.Decision == "accept" {
				if err != nil {
					t.Fatalf("accepted case refused: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("rejected case accepted")
			}
			classification, ok := ErrorClassification(err)
			if !ok || classification != candidate.Expected.Classification {
				t.Fatalf("classification = %q, %v; want %q (%v)", classification, ok, candidate.Expected.Classification, err)
			}
			if candidate.Expected.Effects == nil || candidate.Expected.Effects.State || candidate.Expected.Effects.Approval ||
				candidate.Expected.Effects.Key || candidate.Expected.Effects.IPCEndpoint || candidate.Expected.Effects.Process ||
				candidate.Expected.Effects.Runtime || candidate.Expected.Effects.Backend || candidate.Expected.Effects.Guest {
				t.Fatal("rejection does not assert the complete zero-effect oracle")
			}
		})
	}
}

func TestKnownAnswersAndCalculatedMaxima(t *testing.T) {
	t.Parallel()
	answers := []struct {
		path   string
		length int
		digest string
	}{
		{"mjs-source-validator/engineering-candidate.bin", 292, "2c39757c40198074f1b1dd6e0ed37fb6c75c1c699c0090e2aa4b8ae88cecc9af"},
		{"mjs-source-validator/artifact-profile.bin", 160, "2075ee498ce4b3d81843d57c5289f0056092aa5e3b575d885018c7348419fc8b"},
		{"mjs-source-validator/request-minimum.bin", 80, "4e40f3057a7d4fe7814b74806d794c91e30b5f79f8b55d12c0d76c0c177fc4f7"},
		{"mjs-source-validator/request-ordinary.bin", 137, "5ea0960c1b8200f483ee29fa2756b15c69a62aadcb1dc675c244569b295795fc"},
		{"mjs-source-validator/request-exact-maximum.bin", RequestMaximumBytes, "24db6d0599ba05dead0e7eac1f1454262506c8c187b5542694bb6bb12e6a0571"},
		{"mjs-source-validator/request-cap-plus-one.bin", RequestMaximumBytes + 1, "4add49e80567b4c097bbad0f989863fe028b868304b18d4ef40c9349a888428b"},
		{"mjs-source-validator/result-ordinary.bin", ResultFrameBytes, "5a428178447f70607367f0c052f674ce74315a60376cf15826de06562ea1ca25"},
	}
	for _, answer := range answers {
		value := readFixture(t, answer.path)
		if len(value) != answer.length {
			t.Fatalf("%s length = %d, want %d", answer.path, len(value), answer.length)
		}
		digest := sha256.Sum256(value)
		if hex.EncodeToString(digest[:]) != answer.digest {
			t.Fatalf("%s digest = %x, want %s", answer.path, digest, answer.digest)
		}
	}

	candidateBytes := readFixture(t, "mjs-source-validator/engineering-candidate.bin")
	candidateDigest := EngineeringCandidateIdentityDigest(candidateBytes)
	if got := hex.EncodeToString(candidateDigest[:]); got != "11a08cace3ddc4dde925bd93ea3de8ea313f7f5766839deca802078df038d0c6" {
		t.Fatalf("engineering candidate identity digest = %s", got)
	}
	profileBytes := readFixture(t, "mjs-source-validator/artifact-profile.bin")
	profileDigest := ArtifactProfileIdentityDigest(profileBytes)
	if got := hex.EncodeToString(profileDigest[:]); got != "0fe1523741ac3cedc32e062778836d73f4eb21f422b3587d78ad81df8c180908" {
		t.Fatalf("artifact profile identity digest = %s", got)
	}
}

func TestIndependentGoEncodersAndDefensiveCopies(t *testing.T) {
	t.Parallel()
	candidateBytes, err := EncodeEngineeringCandidate(DefaultEngineeringCandidate())
	if err != nil {
		t.Fatalf("encode candidate: %v", err)
	}
	if !bytes.Equal(candidateBytes, readFixture(t, "mjs-source-validator/engineering-candidate.bin")) {
		t.Fatal("Go engineering-candidate encoding disagrees with generated bytes")
	}

	var executable ExecutableDigest
	var build BuildManifestDigest
	var assessment AssessmentDigest
	for index := range executable {
		executable[index], build[index], assessment[index] = 0x11, 0x22, 0x33
	}
	profile, err := NewEnrolledArtifactProfile(executable, build, assessment)
	if err != nil {
		t.Fatalf("construct profile: %v", err)
	}
	profileBytes, err := EncodeArtifactProfile(profile)
	if err != nil {
		t.Fatalf("encode profile: %v", err)
	}
	if !bytes.Equal(profileBytes, readFixture(t, "mjs-source-validator/artifact-profile.bin")) {
		t.Fatal("Go artifact-profile encoding disagrees with generated bytes")
	}

	callerSource := readFixture(t, "mjs-source/ordinary.mjs")
	var requestID ValidationRequestID
	for index := range requestID {
		requestID[index] = byte(index + 1)
	}
	request, err := NewValidationRequest(requestID, callerSource)
	if err != nil {
		t.Fatalf("construct request: %v", err)
	}
	requestBytes, err := EncodeRequest(request)
	if err != nil {
		t.Fatalf("encode request: %v", err)
	}
	if !bytes.Equal(requestBytes, readFixture(t, "mjs-source-validator/request-ordinary.bin")) {
		t.Fatal("Go request encoding disagrees with generated bytes")
	}
	callerSource[0] ^= 0xff
	returned := request.SourceBytes()
	returned[0] ^= 0xff
	if !bytes.Equal(request.SourceBytes(), readFixture(t, "mjs-source/ordinary.mjs")) {
		t.Fatal("request retained caller or accessor storage")
	}

	result, err := NewValidationResult(request, ArtifactProfileIdentityDigest(profileBytes), ParseValid, PolicyAllow, ResultNoFinding, [5]uint32{})
	if err != nil {
		t.Fatalf("construct result: %v", err)
	}
	resultBytes, err := EncodeResult(result)
	if err != nil {
		t.Fatalf("encode result: %v", err)
	}
	if !bytes.Equal(resultBytes, readFixture(t, "mjs-source-validator/result-ordinary.bin")) {
		t.Fatal("Go result encoding disagrees with generated bytes")
	}
	resultBytes[0] ^= 0xff
	verified, err := VerifyResult(request, profileBytes, readFixture(t, "mjs-source-validator/result-ordinary.bin"))
	if err != nil || verified.PolicyDisposition != PolicyAllow {
		t.Fatalf("verify independent result after caller mutation: %v", err)
	}
}

func readManifest(t *testing.T) conformanceManifest {
	t.Helper()
	bytes := readFixture(t, "manifest.json")
	var manifest conformanceManifest
	if err := json.Unmarshal(bytes, &manifest); err != nil {
		t.Fatalf("decode manifest: %v", err)
	}
	return manifest
}

func readFixture(t *testing.T, path string) []byte {
	t.Helper()
	root := filepath.Join("..", "..", "..", "schemas", "conformance", "v0")
	value, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return value
}
