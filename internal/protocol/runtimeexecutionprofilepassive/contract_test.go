package runtimeexecutionprofilepassive

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestC2AKnownAnswerAndPassiveBoundary(t *testing.T) {
	t.Parallel()
	exact := readContract(t)
	contract, err := Decode(exact)
	if err != nil {
		t.Fatalf("decode C2A contract: %v", err)
	}
	if contract.WorkStatus.C2A != "PASSED-passive-preparation-only" || contract.WorkStatus.C2B != "BLOCKED" {
		t.Fatal("C2A/C2B status boundary changed")
	}
	if contract.Effects != (Effects{}) {
		t.Fatal("passive C2A contract has an authority effect")
	}

	c1Path := filepath.Join("..", "..", "..", "schemas", "conformance", "c1-governed-deno-core", "controlled-development-profile.json")
	c1, err := os.ReadFile(c1Path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(c1)
	if len(c1) != C1Bytes || hex.EncodeToString(sum[:]) != C1SHA256 {
		t.Fatal("C2A no longer consumes exact C1 bytes")
	}
}

func TestC2ASemanticMutationsRefuse(t *testing.T) {
	t.Parallel()
	contract, err := Decode(readContract(t))
	if err != nil {
		t.Fatal(err)
	}
	zero := "00"
	contract.ArtifactClosure.Required[0].SHA256 = &zero
	if err := Validate(contract); err == nil {
		t.Fatal("invented final artifact identity was accepted")
	}

	contract, _ = Decode(readContract(t))
	contract.C1PlanAndProfileBindings.RuntimeProfile.SelectedIdentity = new(string)
	if err := Validate(contract); err == nil {
		t.Fatal("invented composed runtime-profile identity was accepted")
	}

	contract, _ = Decode(readContract(t))
	contract.RunnerDescriptorProfile.Entries[5].FD = 6
	if err := Validate(contract); err == nil {
		t.Fatal("reused descriptor was accepted")
	}

	contract, _ = Decode(readContract(t))
	contract.MachineProfile.CPUTimeLimitMs = new(int)
	if err := Validate(contract); err == nil {
		t.Fatal("unsupported resource field was activated")
	}

	contract, _ = Decode(readContract(t))
	contract.WorkStatus.Runtime001 = "supported"
	if err := Validate(contract); err == nil {
		t.Fatal("passive fixture advanced RUNTIME-001")
	}
}

func TestC2ADecodedValueOwnsDefensiveCopies(t *testing.T) {
	t.Parallel()
	first, err := Decode(readContract(t))
	if err != nil {
		t.Fatal(err)
	}
	first.C2BMatrix[0].Cases[0] = "changed"
	first.RunnerDescriptorProfile.Entries[0].Role = "changed"
	second, err := Decode(readContract(t))
	if err != nil {
		t.Fatal(err)
	}
	if second.C2BMatrix[0].Cases[0] == "changed" || second.RunnerDescriptorProfile.Entries[0].Role == "changed" {
		t.Fatal("decoded contracts share caller-mutated storage")
	}
}

func TestC2ADecodeRejectsWrongLengthAndDigest(t *testing.T) {
	t.Parallel()
	exact := readContract(t)

	if _, err := Decode(exact[:len(exact)-1]); err == nil {
		t.Fatal("truncated contract accepted")
	}

	tampered := append([]byte(nil), exact...)
	tampered[100] ^= 0xFF
	if len(tampered) != len(exact) {
		t.Fatal("test setup changed the byte length")
	}
	if _, err := Decode(tampered); err == nil {
		t.Fatal("digest-mismatched contract accepted")
	}
}

// TestRequireEOFAcceptsOnlyExactlyOneTrailingAbsence covers requireEOF
// directly: Decode only ever calls it with the exact pinned fixture bytes
// (any tampering is already caught by Decode's own length/digest gate before
// requireEOF runs), so its own error branches are otherwise unreachable
// through the public Decode entry point.
func TestRequireEOFAcceptsOnlyExactlyOneTrailingAbsence(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{"nothing-after-eof", "", false},
		{"only-whitespace-after-eof", "   \n", false},
		{"trailing-json-value", "1", true},
		{"trailing-malformed-syntax", "}", true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			decoder := json.NewDecoder(strings.NewReader(testCase.body))
			err := requireEOF(decoder)
			if (err != nil) != testCase.wantErr {
				t.Fatalf("requireEOF(%q) error = %v, want error = %v", testCase.body, err, testCase.wantErr)
			}
		})
	}
}

// TestClonePanicsOnUnmarshalableRawMessage guards clone's deliberate panic:
// every product-controlled contract round-trips through Marshal/Unmarshal
// cleanly, so if that ever stops holding this should fail loudly rather
// than silently return a truncated copy.
func TestClonePanicsOnUnmarshalableRawMessage(t *testing.T) {
	t.Parallel()
	contract, err := Decode(readContract(t))
	if err != nil {
		t.Fatal(err)
	}
	contract.Scope = json.RawMessage("{not-valid-json")
	defer func() {
		if recover() == nil {
			t.Fatal("clone did not panic on an invalid RawMessage field")
		}
	}()
	clone(contract)
}

// TestC2AValidateRejectsEveryDistinctFieldMismatch extends
// TestC2ASemanticMutationsRefuse's coverage to every remaining independent
// check in Validate, one mutation per check so a future accidental
// loosening of any single comparison is caught.
func TestC2AValidateRejectsEveryDistinctFieldMismatch(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name   string
		mutate func(*Contract)
	}{
		{"contract-identity", func(c *Contract) { c.Contract = "wrong" }},
		{"c1-binding", func(c *Contract) { c.C1Binding.Consumption = "wrong" }},
		{"plan-binding-role-list", func(c *Contract) {
			c.C1PlanAndProfileBindings.ExecutionPlan.RequiredReferenceRoles = append(
				c.C1PlanAndProfileBindings.ExecutionPlan.RequiredReferenceRoles, "extra")
		}},
		{"artifact-role-count", func(c *Contract) {
			c.ArtifactClosure.Required = c.ArtifactClosure.Required[:len(c.ArtifactClosure.Required)-1]
		}},
		{"known-answer-digest", func(c *Contract) { c.KnownAnswer.Source.SHA256 = "0" }},
		{"known-answer-struct-field", func(c *Contract) { c.KnownAnswer.Source.Path = "wrong.mjs" }},
		{"runner-descriptor-range", func(c *Contract) { c.RunnerDescriptorProfile.NumericRange.Maximum = 99 }},
		{"runner-port-calls", func(c *Contract) { c.RunnerDescriptorProfile.PortCalls[0].InputFD = 999 }},
		{"guest-descriptor-profile", func(c *Contract) { c.GuestDescriptorProfile.CloseFromInclusive = 99 }},
		{"guest-descriptor-entry", func(c *Contract) { c.GuestDescriptorProfile.PreChild[0].Role = "wrong" }},
		{"transport-profile", func(c *Contract) { c.TransportProfile.HeaderVersion = 99 }},
		{"teardown-profile", func(c *Contract) { c.TeardownProfile.WallActionAtMs = 99 }},
		{"c2b-matrix-length", func(c *Contract) { c.C2BMatrix = c.C2BMatrix[:len(c.C2BMatrix)-1] }},
		{"c2b-matrix-group-count", func(c *Contract) { c.C2BMatrix[0].Cases = append(c.C2BMatrix[0].Cases, "extra") }},
		{"restoration-length", func(c *Contract) {
			c.RestorationMutations = c.RestorationMutations[:len(c.RestorationMutations)-1]
		}},
		{"restoration-too-short", func(c *Contract) { c.RestorationMutations[0] = "MUT" }},
		{"restoration-wrong-prefix", func(c *Contract) { c.RestorationMutations[0] = "XXX-001-something" }},
		{"work-status", func(c *Contract) { c.WorkStatus.C2A = "wrong" }},
		{"effects", func(c *Contract) { c.Effects.Process = true }},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			contract, err := Decode(readContract(t))
			if err != nil {
				t.Fatal(err)
			}
			testCase.mutate(contract)
			if err := Validate(contract); err == nil {
				t.Fatalf("%s was accepted", testCase.name)
			}
		})
	}
}

func readContract(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join("..", "..", "..", "schemas", "conformance", "c2a-governed-deno-core", "passive-execution-profile.json")
	exact, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return exact
}
