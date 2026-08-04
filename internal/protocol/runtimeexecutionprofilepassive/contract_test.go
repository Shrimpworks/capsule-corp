package runtimeexecutionprofilepassive

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
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

func readContract(t *testing.T) []byte {
	t.Helper()
	path := filepath.Join("..", "..", "..", "schemas", "conformance", "c2a-governed-deno-core", "passive-execution-profile.json")
	exact, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return exact
}
