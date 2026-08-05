package runtimec2bpassivev2

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestKnownAnswerNoGuestBuildClosure(t *testing.T) {
	binding, err := Decode(readFixture(t, "passive-binding-v2.json"))
	if err != nil {
		t.Fatal(err)
	}
	if binding.SchemaVersion != 2 || binding.ArchiveEvidence.MergeCommit != "50108417ebf1aa45788a4e9a6b4ca6b4448e9972" ||
		binding.ArchiveEvidence.ReviewedHead != "518eea04e1f81d27e61178b7f4ff164b955dea76" {
		t.Fatal("successor or archive identity changed")
	}
	if len(binding.ConstructedArtifacts) != 6 || binding.RuntimeManifest.CanonicalAuthority != nil ||
		binding.HostPreflightHarness.FinalHostRunnerIdentity != nil || binding.HostPreflightHarness.Final ||
		binding.HostPreflightHarness.CallsLibkrun || binding.NextConsumerGate.Implemented ||
		binding.NextConsumerGate.AdmissionEffect || binding.Effects != (Effects{}) {
		t.Fatal("passive binding acquired a consumer, final runner, or authority")
	}
	if !allUnresolvedNil(binding.Unresolved) {
		t.Fatal("required unresolved field acquired invented evidence")
	}
}

func TestCapCapPlusOneAndCrossVersionRefuse(t *testing.T) {
	exact := readFixture(t, "passive-binding-v2.json")
	if len(exact) != BindingBytes {
		t.Fatalf("known-answer length: %d", len(exact))
	}
	if _, err := Decode(append(append([]byte(nil), exact...), ' ')); err == nil || !strings.Contains(err.Error(), "CAP") {
		t.Fatalf("cap+1 was not refused as cap: %v", err)
	}
	if _, err := Decode(readFixture(t, "passive-binding.json")); err == nil {
		t.Fatal("v1 bytes accepted by v2 decoder")
	}
}

func TestUnknownMissingDuplicateAndTrailingRefuse(t *testing.T) {
	exact := readFixture(t, "passive-binding-v2.json")
	var generic map[string]any
	if err := json.Unmarshal(exact, &generic); err != nil {
		t.Fatal(err)
	}
	missing := cloneMap(generic)
	delete(missing, "domain")
	missingBytes, err := json.MarshalIndent(missing, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	missingBytes = append(missingBytes, '\n')
	cases := [][]byte{
		bytes.Replace(exact, []byte("{\n"), []byte("{\n  \"unknown\": true,\n"), 1),
		bytes.Replace(exact, []byte("{\n"), []byte("{\n  \"objectType\": \"capsule.governed-deno-core-c2b-passive-binding\",\n"), 1),
		missingBytes,
		append(append([]byte(nil), exact...), []byte("{}")...),
	}
	for index, candidate := range cases {
		if _, err := Decode(candidate); err == nil {
			t.Fatalf("strict mutation %d accepted", index)
		}
	}
}

func TestSubstitutionNullOrderAndAdmissionMutationsRefuse(t *testing.T) {
	original, err := Decode(readFixture(t, "passive-binding-v2.json"))
	if err != nil {
		t.Fatal(err)
	}
	value := "invented"
	limit := int64(1)
	tests := []struct {
		name   string
		mutate func(*Binding)
	}{
		{"version", func(candidate *Binding) { candidate.SchemaVersion = 1 }},
		{"archive", func(candidate *Binding) { candidate.ArchiveEvidence.MergeCommit = strings.Repeat("0", 40) }},
		{"artifact", func(candidate *Binding) { candidate.ConstructedArtifacts[0].SHA256 = strings.Repeat("0", 64) }},
		{"order", func(candidate *Binding) {
			candidate.ConstructedArtifacts[0], candidate.ConstructedArtifacts[1] = candidate.ConstructedArtifacts[1], candidate.ConstructedArtifacts[0]
		}},
		{"runner-final", func(candidate *Binding) { candidate.HostPreflightHarness.Final = true }},
		{"runner-identity", func(candidate *Binding) { candidate.HostPreflightHarness.FinalHostRunnerIdentity = &value }},
		{"firmware", func(candidate *Binding) { candidate.Unresolved.SeparateFirmwareIdentity = &value }},
		{"cpu", func(candidate *Binding) { candidate.Unresolved.CPUTimeLimitMS = &limit }},
		{"memory", func(candidate *Binding) { candidate.Unresolved.HostVMMMemoryLimitBytes = &limit }},
		{"scratch", func(candidate *Binding) { candidate.Unresolved.ScratchMaximumBytes = &limit }},
		{"guest", func(candidate *Binding) { candidate.Unresolved.GuestEvidence = &value }},
		{"admission-null", func(candidate *Binding) { candidate.Unresolved.RuntimeProfileAdmission = &value }},
		{"admission-effect", func(candidate *Binding) { candidate.Effects.Admission = true }},
		{"consumer", func(candidate *Binding) { candidate.NextConsumerGate.Implemented = true }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := clone(original)
			test.mutate(candidate)
			if err := Validate(candidate); err == nil {
				t.Fatal("mutation accepted")
			}
		})
	}
}

func TestDecodedBindingsOwnDefensiveCopies(t *testing.T) {
	first, err := Decode(readFixture(t, "passive-binding-v2.json"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := Decode(readFixture(t, "passive-binding-v2.json"))
	if err != nil {
		t.Fatal(err)
	}
	first.ConstructedArtifacts[0].Role = "changed"
	first.NextConsumerGate.Requirements[0] = "changed"
	if second.ConstructedArtifacts[0].Role == "changed" || second.NextConsumerGate.Requirements[0] == "changed" {
		t.Fatal("decoded bindings share mutable storage")
	}
}

func TestNoProductConsumerImportsPassiveV2Package(t *testing.T) {
	root := os.DirFS(filepath.Join("..", "..", ".."))
	needle := "internal/protocol/runtimec2bpassivev2"
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
		if filepath.Ext(path) != ".go" || strings.Contains(path, "runtimec2bpassivev2") {
			return nil
		}
		contents, readErr := fs.ReadFile(root, path)
		if readErr != nil {
			return readErr
		}
		if bytes.Contains(contents, []byte(needle)) {
			t.Errorf("passive v2 package has a product consumer: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func allUnresolvedNil(value Unresolved) bool {
	return value.FinalHostRunnerIdentity == nil && value.SeparateFirmwareIdentity == nil &&
		value.ComposedProfileIdentity == nil && value.ComposedProfileDigest == nil &&
		value.CPUTimeLimitMS == nil && value.HostVMMMemoryLimitBytes == nil &&
		value.ScratchMaximumBytes == nil && value.GuestEvidence == nil &&
		value.RuntimeManifestAdmission == nil && value.RuntimeProfileAdmission == nil &&
		value.BackendAdmission == nil
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
