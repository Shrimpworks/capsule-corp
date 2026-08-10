package typedguesttransportpassive

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type fixtureManifest struct {
	Contract   string            `json:"contract"`
	Version    int               `json:"version"`
	Status     string            `json:"status"`
	Bindings   map[string]string `json:"bindings"`
	Cases      []fixtureCase     `json:"cases"`
	StateCases []struct {
		ID                string `json:"id"`
		Disposition       string `json:"disposition"`
		DurableCompletion bool   `json:"durableCompletion"`
	} `json:"stateCases"`
	RestorationCases []struct {
		ID              string `json:"id"`
		RefusingControl string `json:"refusingControl"`
		Detected        bool   `json:"detected"`
	} `json:"restorationCases"`
	Effects map[string]bool `json:"effects"`
}

type fixtureCase struct {
	ID          string      `json:"id"`
	Decision    string      `json:"decision"`
	Role        uint16      `json:"role"`
	Path        string      `json:"path"`
	Bytes       int         `json:"bytes"`
	SHA256      string      `json:"sha256"`
	Disposition Disposition `json:"disposition"`
}

func TestIndependentFixtureDecoder(t *testing.T) {
	t.Parallel()
	root := filepath.Join("..", "..", "..", "schemas", "conformance", "typed-guest-transport-v1")
	manifestBytes, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest fixtureManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Contract != "capsule.typed-guest-transport" || manifest.Version != 1 || manifest.Status != "passive-conformance-only" {
		t.Fatalf("unexpected contract identity: %+v", manifest)
	}
	bindings := decodeBindings(t, manifest.Bindings)
	for _, test := range manifest.Cases {
		test := test
		t.Run(test.ID, func(t *testing.T) {
			exact, readErr := os.ReadFile(filepath.Join(root, test.Path))
			if readErr != nil {
				t.Fatal(readErr)
			}
			if len(exact) != test.Bytes {
				t.Fatalf("bytes = %d, want %d", len(exact), test.Bytes)
			}
			digest := sha256.Sum256(exact)
			if hex.EncodeToString(digest[:]) != test.SHA256 {
				t.Fatal("fixture digest mismatch")
			}
			_, decodeErr := Decode(exact, test.Role, bindings)
			if test.Decision == "accept" {
				if decodeErr != nil {
					t.Fatalf("accepted fixture refused: %v", decodeErr)
				}
				return
			}
			var refusal *Refusal
			if !errors.As(decodeErr, &refusal) || refusal.Disposition != test.Disposition {
				t.Fatalf("refusal = %v, want %s", decodeErr, test.Disposition)
			}
		})
	}
	if len(manifest.Cases) < 40 {
		t.Fatalf("only %d frame cases", len(manifest.Cases))
	}
	expectedStates := map[string]struct {
		disposition string
		durable     bool
	}{
		"cancel-before-g":                         {"REFUSED_CANCELLED_BEFORE_START", false},
		"cancel-during-source":                    {"REFUSED_INPUT_CANCELLED", false},
		"source-short-write-error":                {"REFUSED_SOURCE_TRANSPORT_FAULT", false},
		"input-zero-progress-deadline":            {"REFUSED_INPUT_STALL", false},
		"completion-reader-death":                 {"REFUSED_COMPLETION_READER_FAULT", false},
		"completion-reset-before-trailer":         {"REFUSED_MISSING_COMMIT", false},
		"valid-frame-before-absence":              {"FRAME_ONLY_AWAIT_TERMINAL_PROOF", false},
		"launcher-death-after-trailer":            {"REFUSED_LIFECYCLE", false},
		"runner-zero-without-frame":               {"REFUSED_MISSING_COMMIT", false},
		"caller-response-loss-after-commit":       {"REPLAY_BYTE_IDENTICAL", true},
		"store-indeterminate":                     {"RECOVERY_REQUIRED_FENCE", false},
		"cancel-concurrent-before-terminal-proof": {"REFUSED_CANCELLED", false},
		"cancel-after-durable-commit":             {"REPLAY_BYTE_IDENTICAL", true},
	}
	if len(manifest.StateCases) != len(expectedStates) {
		t.Fatalf("state case count = %d", len(manifest.StateCases))
	}
	for _, test := range manifest.StateCases {
		expected, ok := expectedStates[test.ID]
		if !ok || expected.disposition != test.Disposition || expected.durable != test.DurableCompletion {
			t.Fatalf("unexpected state case %+v", test)
		}
	}
	expectedRestorations := map[string]string{
		"source-input-endpoint-swap":        "endpoint-object-identity-and-role-magic",
		"completion-source-endpoint-swap":   "endpoint-object-identity-and-role-magic",
		"descriptor-alias":                  "dedicated-open-description-manifest",
		"descriptor-wrong-mode":             "access-mode-manifest",
		"descriptor-cloexec-cleared":        "pre-post-cloexec-canary",
		"descriptor-nonblocking-changed":    "pre-post-status-flag-canary",
		"descriptor-inherited-by-workload":  "closed-child-fd-manifest",
		"wrong-attempt":                     "typed-attempt-binding",
		"wrong-registration":                "typed-registration-binding",
		"plan-digest-substitution":          "retained-plan-digest",
		"runtime-profile-substitution":      "retained-runtime-profile-digest",
		"payload-flood-after-valid-prefix":  "full-drained-byte-count",
		"early-trailer":                     "calculated-final-trailer-offset",
		"eof-as-commit":                     "commit-trailer-required",
		"runner-zero-as-success":            "terminal-proof-join",
		"diagnostic-console-substitution":   "typed-role-magic-and-channel",
		"implicit-console-restored":         "runner-call-and-device-inventory",
		"vsock-restored":                    "runner-call-and-device-inventory",
		"network-device-restored":           "runner-call-and-device-inventory",
		"virtiofs-restored":                 "runner-call-and-device-inventory",
		"cleanup-unresolved":                "terminal-proof-join",
		"response-loss-result-substitution": "immutable-attempt-replay",
		"durable-commit-corrupt":            "store-reopen-full-validation",
	}
	if len(manifest.RestorationCases) != len(expectedRestorations) {
		t.Fatalf("restoration case count = %d", len(manifest.RestorationCases))
	}
	for _, test := range manifest.RestorationCases {
		if expectedRestorations[test.ID] != test.RefusingControl || !test.Detected {
			t.Fatalf("unexpected restoration %+v", test)
		}
	}
	for effect, enabled := range manifest.Effects {
		if enabled {
			t.Fatalf("passive fixture enables %s", effect)
		}
	}
}

func TestRetainedPayloadBindingIsIndependentOfFrameSelfDigest(t *testing.T) {
	t.Parallel()
	root := filepath.Join("..", "..", "..", "schemas", "conformance", "typed-guest-transport-v1")
	manifestBytes, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest fixtureManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatal(err)
	}
	bindings := decodeBindings(t, manifest.Bindings)
	var source fixtureCase
	for _, candidate := range manifest.Cases {
		if candidate.ID == "accept-source-ordinary" {
			source = candidate
		}
	}
	exact, err := os.ReadFile(filepath.Join(root, source.Path))
	if err != nil {
		t.Fatal(err)
	}
	payloadBytes := uint64(len(exact) - InputHeaderBytes)
	payloadDigest := sha256.Sum256(exact[InputHeaderBytes:])
	bindings.ExpectedPayloadBytes = &payloadBytes
	bindings.ExpectedPayloadDigest = &payloadDigest
	if _, err := Decode(exact, SourceRole, bindings); err != nil {
		t.Fatal(err)
	}
	wrongDigest := payloadDigest
	wrongDigest[0] ^= 1
	bindings.ExpectedPayloadDigest = &wrongDigest
	_, err = Decode(exact, SourceRole, bindings)
	var refusal *Refusal
	if !errors.As(err, &refusal) || refusal.Disposition != BindingMismatch {
		t.Fatalf("expected retained payload binding refusal, got %v", err)
	}
}

func TestCanonicalJSONRules(t *testing.T) {
	t.Parallel()
	for _, accepted := range []string{
		"null", "true", "false", "0", "-1", `"text"`, `[1,"x",null]`, `{"a":1,"b":2}`,
	} {
		if err := ValidateCanonicalJSON([]byte(accepted)); err != nil {
			t.Errorf("accepted %q: %v", accepted, err)
		}
	}
	for _, rejected := range []string{
		"", " 0", "0 ", "-0", "1.0", "1e2", "9007199254740992", `{"b":2,"a":1}`,
		`{"a":1,"a":2}`, `{"bad key":1}`, `"\/"`, "nullnull",
	} {
		if err := ValidateCanonicalJSON([]byte(rejected)); err == nil {
			t.Errorf("accepted noncanonical %q", rejected)
		}
	}
}

func decodeBindings(t *testing.T, fields map[string]string) Bindings {
	t.Helper()
	var result Bindings
	decode := func(name string, destination []byte) {
		exact, err := hex.DecodeString(fields[name])
		if err != nil || len(exact) != len(destination) {
			t.Fatalf("invalid %s", name)
		}
		copy(destination, exact)
	}
	decode("attemptId", result.AttemptID[:])
	decode("registrationId", result.RegistrationID[:])
	decode("planDigest", result.PlanDigest[:])
	decode("runtimeProfileDigest", result.RuntimeProfileDigest[:])
	return result
}
