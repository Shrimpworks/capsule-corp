package runtimec2bpassive

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"capsule.local/capsule/internal/protocol/strictjson"
)

func TestKnownAnswerCrossLinksAllPassiveInputs(t *testing.T) {
	binding, err := Decode(read(t, "passive-binding.json"), readInputs(t))
	if err != nil {
		t.Fatalf("decode passive binding: %v", err)
	}
	if binding.WorkStatus.PassiveReconciliation != "PASSED" ||
		binding.WorkStatus.ParentGovernedRuntime != "IN_PROGRESS — TRENDING_GOOD" ||
		binding.WorkStatus.C2BComposedProfileGuestExecution != "BLOCKED" ||
		binding.WorkStatus.RuntimeProfileAdmission != "BLOCKED" ||
		binding.WorkStatus.Runtime001 != "unsupported" || binding.WorkStatus.VMM001 != "unsupported" {
		t.Fatal("status or admission boundary changed")
	}
	if binding.BuildEvidence.IndependentBuilder || binding.NextConsumerGate.Implemented ||
		binding.NextConsumerGate.AdmissionEffect || binding.Effects != (Effects{}) {
		t.Fatal("passive binding acquired authority or overstated build evidence")
	}
	if len(binding.Artifacts) != 6 || len(binding.DependencyMergePolicy.Dependencies) != 2 {
		t.Fatal("closed artifact or dependency inventory changed")
	}
}

func TestSemanticMutationClassesRefuse(t *testing.T) {
	original, err := decodeClosed(read(t, "passive-binding.json"))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name   string
		mutate func(*Binding)
	}{
		{"wrong-version", func(value *Binding) { value.SchemaVersion = 2 }},
		{"wrong-domain", func(value *Binding) { value.Domain = "capsule.other/v1" }},
		{"wrong-length", func(value *Binding) { value.FixedFixture.Source.Bytes++ }},
		{"wrong-digest", func(value *Binding) { value.FixedFixture.Source.SHA256 = strings.Repeat("0", 64) }},
		{"cross-link", func(value *Binding) { value.BuildEvidence.ForkSupplementIdentity = "capsule.other" }},
		{"substitution", func(value *Binding) { value.ForkSupplement.Head = strings.Repeat("0", 40) }},
		{"order", func(value *Binding) { value.Artifacts[0], value.Artifacts[1] = value.Artifacts[1], value.Artifacts[0] }},
		{"cap", func(value *Binding) { value.FixedFixture.Source.MaximumBytes++ }},
		{"admission-state", func(value *Binding) { value.Effects.Admission = true }},
		{"dependency-state", func(value *Binding) { value.DependencyMergePolicy.State = "PASSED" }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := strictjson.Clone(original)
			test.mutate(candidate)
			if err := Validate(candidate); err == nil {
				t.Fatal("mutation accepted")
			}
		})
	}
}

func TestClosedDecoderRejectsUnknownDuplicateMissingAndTrailing(t *testing.T) {
	exact := read(t, "passive-binding.json")
	var generic map[string]any
	if err := json.Unmarshal(exact, &generic); err != nil {
		t.Fatal(err)
	}
	generic["unknown"] = true
	unknown, err := json.Marshal(generic)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := decodeClosed(unknown); err == nil {
		t.Fatal("unknown field accepted")
	}
	delete(generic, "unknown")
	delete(generic, "domain")
	missing, err := json.Marshal(generic)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := decodeClosed(missing)
	if err != nil {
		t.Fatal(err)
	}
	if err := Validate(decoded); err == nil {
		t.Fatal("missing field accepted")
	}
	duplicate := bytes.Replace(exact, []byte("{\n"), []byte("{\n  \"objectType\": \"capsule.governed-deno-core-c2b-passive-binding\",\n"), 1)
	if _, err := decodeClosed(duplicate); err == nil || !strings.Contains(err.Error(), "DUPLICATE") {
		t.Fatalf("duplicate field was not classified: %v", err)
	}
	if _, err := decodeClosed(append(append([]byte(nil), exact...), []byte("{}")...)); err == nil {
		t.Fatal("trailing value accepted")
	}
}

func TestInputLengthDigestAndSelfDigestRefuse(t *testing.T) {
	exact := read(t, "passive-binding.json")
	inputs := readInputs(t)
	mutations := []struct {
		name   string
		mutate func(*Inputs)
	}{
		{"c1-length", func(value *Inputs) { value.C1 = value.C1[:len(value.C1)-1] }},
		{"c2a-digest", func(value *Inputs) { value.C2A = mutateByte(value.C2A) }},
		{"fork-digest", func(value *Inputs) { value.ForkSupplement = mutateByte(value.ForkSupplement) }},
		{"evidence-digest", func(value *Inputs) { value.BuildEvidence = mutateByte(value.BuildEvidence) }},
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			candidate := cloneInputs(inputs)
			mutation.mutate(&candidate)
			if _, err := Decode(exact, candidate); err == nil {
				t.Fatal("input mutation accepted")
			}
		})
	}
	if err := validateEvidenceSelfDigest(inputs.BuildEvidence); err != nil {
		t.Fatalf("evidence self-digest: %v", err)
	}
}

func TestDecodedBindingsOwnDefensiveCopies(t *testing.T) {
	inputs := readInputs(t)
	first, err := Decode(read(t, "passive-binding.json"), inputs)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Decode(read(t, "passive-binding.json"), inputs)
	if err != nil {
		t.Fatal(err)
	}
	first.Artifacts[0].Role = "changed"
	first.DependencyMergePolicy.Dependencies[0].Role = "changed"
	first.NextConsumerGate.Requirements[0] = "changed"
	if second.Artifacts[0].Role == "changed" || second.DependencyMergePolicy.Dependencies[0].Role == "changed" || second.NextConsumerGate.Requirements[0] == "changed" {
		t.Fatal("decoded bindings share mutable storage")
	}
}

func readInputs(t *testing.T) Inputs {
	t.Helper()
	return Inputs{
		C1:             readFromRoot(t, "schemas", "conformance", "c1-governed-deno-core", "controlled-development-profile.json"),
		C2A:            readFromRoot(t, "schemas", "conformance", "c2a-governed-deno-core", "passive-execution-profile.json"),
		ForkSupplement: read(t, "inputs", "fork-supplement.json"),
		BuildEvidence:  read(t, "inputs", "build-evidence.json"),
	}
}

func read(t *testing.T, parts ...string) []byte {
	t.Helper()
	base := []string{"schemas", "conformance", "c2b-governed-deno-core"}
	return readFromRoot(t, append(base, parts...)...)
}

func readFromRoot(t *testing.T, parts ...string) []byte {
	t.Helper()
	path := filepath.Join(append([]string{"..", "..", ".."}, parts...)...)
	value, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func mutateByte(value []byte) []byte {
	copy := append([]byte(nil), value...)
	copy[len(copy)/2] ^= 1
	return copy
}

func cloneInputs(inputs Inputs) Inputs {
	return Inputs{
		C1:             append([]byte(nil), inputs.C1...),
		C2A:            append([]byte(nil), inputs.C2A...),
		ForkSupplement: append([]byte(nil), inputs.ForkSupplement...),
		BuildEvidence:  append([]byte(nil), inputs.BuildEvidence...),
	}
}
