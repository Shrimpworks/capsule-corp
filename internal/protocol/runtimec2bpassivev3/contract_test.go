package runtimec2bpassivev3

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestKnownAnswerPassiveSuccessorContract(t *testing.T) {
	binding, err := Decode(readFixture(t, "passive-binding-v3.json"))
	if err != nil {
		t.Fatal(err)
	}
	if binding.SchemaVersion != 3 || binding.Status != "PASSED-passive-successor-contract-only" {
		t.Fatal("passive successor identity changed")
	}
	assertJSONField(t, binding.GovernedInputs, []string{"libkrun", "acceptedCommit"}, "7432eda5a49220976b0167005aa43ee622f9d632")
	assertJSONField(t, binding.ComposedProfile, []string{"contractDigestSha256"}, ContractSHA256)
	assertJSONField(t, binding.NextGuestGate, []string{"state"}, "BLOCKED")
	assertJSONField(t, binding.NextGuestGate, []string{"guestAuthorization"}, false)
	var generic any
	if err := json.Unmarshal(readFixture(t, "passive-binding-v3.json"), &generic); err != nil {
		t.Fatal(err)
	}
	if containsNull(generic) {
		t.Fatal("passive successor contains an authority-bearing null")
	}
	for _, unsupported := range []string{"cpuTimeLimitMs", "hostVmmMemoryLimitBytes", "scratchMaximumBytes"} {
		if bytes.Contains(binding.ResourceContract, []byte(unsupported)) {
			t.Fatalf("unsupported resource field present: %s", unsupported)
		}
	}
}

func TestCapCapPlusOneAndCrossVersionRefuse(t *testing.T) {
	exact := readFixture(t, "passive-binding-v3.json")
	if len(exact) != BindingBytes {
		t.Fatalf("known-answer length: %d", len(exact))
	}
	if _, err := Decode(append(append([]byte(nil), exact...), ' ')); err == nil || !strings.Contains(err.Error(), "CAP") {
		t.Fatalf("cap+1 was not refused as cap: %v", err)
	}
	if _, err := Decode(readFixture(t, "passive-binding-v2.json")); err == nil {
		t.Fatal("v2 bytes accepted by v3 decoder")
	}
}

func TestUnknownMissingDuplicateTrailingAndNullRefuse(t *testing.T) {
	exact := readFixture(t, "passive-binding-v3.json")
	var generic map[string]any
	if err := json.Unmarshal(exact, &generic); err != nil {
		t.Fatal(err)
	}
	missing := cloneMap(generic)
	delete(missing, "domain")
	missingBytes := canonicalJSON(t, missing)
	withNull := cloneMap(generic)
	withNull["unknown"] = nil
	cases := [][]byte{
		bytes.Replace(exact, []byte("{\n"), []byte("{\n  \"unknown\": true,\n"), 1),
		bytes.Replace(exact, []byte("{\n"), []byte("{\n  \"objectType\": \"capsule.governed-deno-core-c2b-passive-binding\",\n"), 1),
		missingBytes,
		canonicalJSON(t, withNull),
		append(append([]byte(nil), exact...), []byte("{}")...),
	}
	for index, candidate := range cases {
		if _, err := Decode(candidate); err == nil {
			t.Fatalf("strict mutation %d accepted", index)
		}
	}
}

func TestSourceRunnerDeviceRuntimeResourceTeardownDigestAndAuthorityMutationsRefuse(t *testing.T) {
	exact := readFixture(t, "passive-binding-v3.json")
	tests := []struct {
		name string
		old  string
		new  string
	}{
		{"libkrun-source", "7432eda5a49220976b0167005aa43ee622f9d632", strings.Repeat("0", 40)},
		{"runner-final", "BLOCKED-final-runner-bytes-and-digest-not-constructed", "PASSED-final-runner-bytes-and-digest-not-constructed"},
		{"implicit-vsock", "explicitly-disabled", "unexpectedly-enabled"},
		{"module-loader", `"moduleLoader": "none"`, `"moduleLoader": "filesystem"`},
		{"string-code-generation", "v8-context-set_allow_generation_from_strings-false-and-verified", "enabled"},
		{"resource", `"vcpus": 1`, `"vcpus": 2`},
		{"teardown", `"forcedSignal": "SIGKILL"`, `"forcedSignal": "SIGTERM"`},
		{"profile-digest", ContractSHA256, strings.Repeat("0", 64)},
		{"guest-authority", `"guestAuthorization": false`, `"guestAuthorization": true `},
		{"admission", `"admission": false`, `"admission": true `},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := bytes.Replace(exact, []byte(test.old), []byte(test.new), 1)
			if bytes.Equal(candidate, exact) {
				t.Fatal("mutation did not apply")
			}
			if len(candidate) != len(exact) {
				candidate = canonicalMutation(t, candidate)
			}
			if _, err := Decode(candidate); err == nil {
				t.Fatal("mutation accepted")
			}
		})
	}
}

func TestDecodedBindingsOwnDefensiveCopies(t *testing.T) {
	first, err := Decode(readFixture(t, "passive-binding-v3.json"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := Decode(readFixture(t, "passive-binding-v3.json"))
	if err != nil {
		t.Fatal(err)
	}
	first.ArtifactRoles[0] = 'X'
	if second.ArtifactRoles[0] == 'X' {
		t.Fatal("decoded bindings share mutable storage")
	}
}

func TestNoProductConsumerImportsPassiveV3Package(t *testing.T) {
	root := os.DirFS(filepath.Join("..", "..", ".."))
	needle := "internal/protocol/runtimec2bpassivev3"
	err := fs.WalkDir(root, ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" || entry.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" || strings.Contains(path, "runtimec2bpassivev3") {
			return nil
		}
		contents, readErr := fs.ReadFile(root, path)
		if readErr != nil {
			return readErr
		}
		if bytes.Contains(contents, []byte(needle)) {
			t.Errorf("passive v3 package has a product consumer: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func assertJSONField(t *testing.T, raw json.RawMessage, path []string, expected any) {
	t.Helper()
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		t.Fatal(err)
	}
	for _, part := range path {
		object, ok := value.(map[string]any)
		if !ok {
			t.Fatalf("non-object before %s", part)
		}
		value = object[part]
	}
	if value != expected {
		t.Fatalf("got %#v want %#v", value, expected)
	}
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("..", "..", "..", "schemas", "conformance", "c2b-governed-deno-core", name)
	value, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func cloneMap(value map[string]any) map[string]any {
	encoded, _ := json.Marshal(value)
	var copy map[string]any
	_ = json.Unmarshal(encoded, &copy)
	return copy
}

func canonicalJSON(t *testing.T, value any) []byte {
	t.Helper()
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return append(encoded, '\n')
}

func canonicalMutation(t *testing.T, candidate []byte) []byte {
	t.Helper()
	var value any
	if err := json.Unmarshal(candidate, &value); err != nil {
		return candidate
	}
	return canonicalJSON(t, value)
}
