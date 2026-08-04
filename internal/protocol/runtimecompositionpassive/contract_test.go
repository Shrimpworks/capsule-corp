package runtimecompositionpassive

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestC1KnownAnswerAndPassiveBoundary(t *testing.T) {
	t.Parallel()
	exact := readContract(t)
	contract, err := Decode(exact)
	if err != nil {
		t.Fatalf("decode C1 contract: %v", err)
	}
	if contract.Resources.MachineProfileRef != nil || contract.Resources.Activation != "refuse-until-c2-and-admission" {
		t.Fatal("C1 contract activated an unselected machine profile")
	}
	if contract.Effects.Guest || contract.Effects.Runtime || contract.Effects.Admission {
		t.Fatal("passive C1 contract has an authority effect")
	}
}

func TestC1ExactBytesAndSemanticMutationsRefuse(t *testing.T) {
	t.Parallel()
	exact := readContract(t)
	mutated := append([]byte(nil), exact...)
	mutated[len(mutated)-2] ^= 1
	if _, err := Decode(mutated); err == nil {
		t.Fatal("one-byte contract mutation was accepted")
	}

	contract, err := Decode(exact)
	if err != nil {
		t.Fatal(err)
	}
	contract.Artifacts["snapshot"] = Artifact{Role: "governed-deno-core-startup-snapshot", SHA256: "00", Admitted: false}
	if err := Validate(contract); err == nil {
		t.Fatal("snapshot mutation was accepted")
	}

	contract, err = Decode(exact)
	if err != nil {
		t.Fatal(err)
	}
	contract.RuntimeSurface.ModuleLoader = "filesystem"
	if err := Validate(contract); err == nil {
		t.Fatal("loader restoration was accepted")
	}

	contract, err = Decode(exact)
	if err != nil {
		t.Fatal(err)
	}
	contract.Resources.MachineProfileRef = new(string)
	if err := Validate(contract); err == nil {
		t.Fatal("invented machine profile was accepted")
	}
}

// TestC1ValidateRefusesEveryMutationClass extends the hand-picked mutation cases in
// TestC1ExactBytesAndSemanticMutationsRefuse with one case per remaining Validate() refusal
// branch, so a future regression that silently loosens any single check fails this suite instead
// of only passing because a different branch happened to catch the same input.
func TestC1ValidateRefusesEveryMutationClass(t *testing.T) {
	t.Parallel()
	exact := readContract(t)

	zeroCommit := strings.Repeat("0", 40)
	cases := []struct {
		name   string
		mutate func(*Contract)
	}{
		{"evidence-repository", func(c *Contract) { c.Evidence.Repository = "Shrimpworks/other" }},
		{"evidence-merge-commit", func(c *Contract) { c.Evidence.MergeCommit = zeroCommit }},
		{"fork-deno-head", func(c *Contract) { c.Forks.Deno.Head = zeroCommit }},
		{"fork-rustyv8-upstream", func(c *Contract) { c.Forks.RustyV8.Upstream = zeroCommit }},
		{"effects-process-true", func(c *Contract) { c.Effects.Process = true }},
		{"effects-runtime-true", func(c *Contract) { c.Effects.Runtime = true }},
		{"effects-guest-true", func(c *Contract) { c.Effects.Guest = true }},
		{"effects-admission-true", func(c *Contract) { c.Effects.Admission = true }},
		{"effects-signing-true", func(c *Contract) { c.Effects.Signing = true }},
		{"effects-release-true", func(c *Contract) { c.Effects.Release = true }},
		{"refusal-codes-truncated", func(c *Contract) {
			c.RefusalCodes = c.RefusalCodes[:len(c.RefusalCodes)-1]
		}},
		{"stop-conditions-truncated", func(c *Contract) {
			c.StopConditions = c.StopConditions[:len(c.StopConditions)-1]
		}},
		{"runtime-surface-globals-mismatch", func(c *Contract) {
			// A permitted global can never legally overlap a forbidden one, but the closest
			// reachable mutation is widening the permitted list, which the exact-match check
			// (and, if it did not exist, the overlap check) must both refuse.
			c.RuntimeSurface.PermittedGlobals = append(
				append([]string(nil), c.RuntimeSurface.PermittedGlobals...), "eval")
		}},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			contract, err := Decode(exact)
			if err != nil {
				t.Fatal(err)
			}
			testCase.mutate(contract)
			if err := Validate(contract); err == nil {
				t.Fatalf("%s mutation was accepted", testCase.name)
			}
		})
	}
}

// TestC1DecoderRejectsUnknownFieldsAndTrailingData exercises Decode's JSON-decoding defenses
// (DisallowUnknownFields, requireEOF) directly. Decode's fixed byte-length/SHA-256 gate runs
// before those defenses and rejects any mutated payload on its own, so a mutation that adds an
// unknown field or trailing bytes can never reach them through the public Decode entry point.
// decodeOnly replicates only Decode's JSON-parsing stage to prove those defenses still fire on
// their own.
func TestC1DecoderRejectsUnknownFieldsAndTrailingData(t *testing.T) {
	t.Parallel()
	exact := readContract(t)

	var asObject map[string]any
	if err := json.Unmarshal(exact, &asObject); err != nil {
		t.Fatal(err)
	}
	asObject["unexpectedField"] = true
	withExtraField, err := json.Marshal(asObject)
	if err != nil {
		t.Fatal(err)
	}
	if err := decodeOnly(withExtraField); err == nil {
		t.Fatal("unknown top-level field was accepted")
	}

	withTrailingData := append(append([]byte(nil), exact...), []byte("\n{}")...)
	if err := decodeOnly(withTrailingData); err == nil {
		t.Fatal("trailing data after the JSON object was accepted")
	}
}

func decodeOnly(payload []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	var contract Contract
	if err := decoder.Decode(&contract); err != nil {
		return err
	}
	return requireEOF(decoder)
}

func TestC1DecodedValueOwnsDefensiveCopies(t *testing.T) {
	t.Parallel()
	exact := readContract(t)
	contract, err := Decode(exact)
	if err != nil {
		t.Fatal(err)
	}
	contract.RuntimeSurface.BuiltinOps[0] = "op_print"
	contract.Artifacts["snapshot"] = Artifact{}

	second, err := Decode(exact)
	if err != nil {
		t.Fatal(err)
	}
	if second.RuntimeSurface.BuiltinOps[0] != expectedOps[0] || second.Artifacts["snapshot"].SHA256 != SnapshotDigest {
		t.Fatal("decoded contracts share caller-mutated storage")
	}
}

func readContract(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join("..", "..", "..", "schemas", "conformance", "c1-governed-deno-core", "controlled-development-profile.json")
	exact, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return exact
}
