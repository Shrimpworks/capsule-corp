package runtimecompositionpassive

import (
	"os"
	"path/filepath"
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
