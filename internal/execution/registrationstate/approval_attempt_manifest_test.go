package registrationstate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type sliceBManifest struct {
	Cases []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Context struct {
			Kind      string          `json:"kind"`
			Operation manifestFixture `json:"operation"`
			Before    manifestFixture `json:"before"`
		} `json:"context"`
		Expected struct {
			Decision                   string  `json:"decision"`
			Classification             *string `json:"classification"`
			AuthorityStateChanged      bool    `json:"authorityStateChanged"`
			FakeBackendEffectPermitted bool    `json:"fakeBackendEffectPermitted"`
			StateDelta                 struct {
				After manifestFixture `json:"after"`
			} `json:"stateDelta"`
		} `json:"expected"`
		Implementations map[string]string `json:"implementations"`
	} `json:"cases"`
}

type manifestFixture struct {
	Path       string `json:"path"`
	SHA256     string `json:"sha256"`
	ByteLength int    `json:"byteLength"`
}

type storeOperationOracle struct {
	ContextType    string `json:"contextType"`
	ContextVersion uint64 `json:"contextVersion"`
	Mode           string `json:"mode"`
	Method         string `json:"method"`
	Scenario       string `json:"scenario"`
	Recovery       bool   `json:"recovery"`
}

type storeStateOracle struct {
	RecoveryFence      bool `json:"recoveryFence"`
	ApprovalPopulation struct {
		UsableCount   int    `json:"usableCount"`
		RetainedCount int    `json:"retainedCount"`
		SetDigest     string `json:"setDigest"`
	} `json:"approvalPopulation"`
	AttemptPopulation struct {
		NonterminalCount int    `json:"nonterminalCount"`
		RetainedCount    int    `json:"retainedCount"`
		SetDigest        string `json:"setDigest"`
	} `json:"attemptPopulation"`
}

func TestSliceBDurableStoreManifestOracles(t *testing.T) {
	manifest := readManifestFixture[sliceBManifest](t, manifestFixture{
		Path: "manifest.json",
	})
	cases := 0
	accepts := 0
	rejects := 0
	seen := make(map[string]struct{})
	for _, testCase := range manifest.Cases {
		if testCase.Context.Kind != "approval-attempt-state" {
			continue
		}
		operation := readManifestFixture[storeOperationOracle](t, testCase.Context.Operation)
		if operation.Mode != "store-transition" {
			continue
		}
		cases++
		seen[operation.Scenario] = struct{}{}
		if testCase.Expected.Decision == "accept" {
			accepts++
		} else {
			rejects++
		}
		if testCase.Implementations["go"] != "verified" ||
			testCase.Expected.FakeBackendEffectPermitted ||
			(testCase.Object != "approval-state" && testCase.Object != "attempt-state") ||
			operation.ContextType != "capsule.conformance.approval-attempt-operation" ||
			operation.ContextVersion != 0 {
			t.Fatalf("invalid Slice B manifest case %s", testCase.ID)
		}
		before := readManifestFixture[storeStateOracle](t, testCase.Context.Before)
		after := readManifestFixture[storeStateOracle](t, testCase.Expected.StateDelta.After)
		if testCase.Expected.Decision == "reject" && testCase.Expected.Classification == nil {
			t.Fatalf("Slice B rejection %s lacks a classification", testCase.ID)
		}
		if operation.Scenario == "atomic-commit" {
			if before.ApprovalPopulation.UsableCount != 1 || after.ApprovalPopulation.UsableCount != 0 ||
				after.ApprovalPopulation.RetainedCount != 1 || after.AttemptPopulation.NonterminalCount != 1 ||
				after.AttemptPopulation.RetainedCount != 1 || !testCase.Expected.AuthorityStateChanged {
				t.Fatal("atomic consume/create manifest oracle is split or incomplete")
			}
		}
		if operation.Scenario == "confirmed-abort" &&
			(before != after || testCase.Expected.AuthorityStateChanged) {
			t.Fatal("confirmed-abort manifest oracle changed authority")
		}
		if operation.Scenario == "nonterminal-cap-plus-one" &&
			(after.AttemptPopulation.NonterminalCount != MaxNonterminalAttempts ||
				after.ApprovalPopulation.UsableCount != 1) {
			t.Fatal("attempt capacity oracle does not preserve its usable approval")
		}
		if operation.Recovery && operation.Scenario != "reopen-post-state" && !after.RecoveryFence {
			t.Fatalf("recovery case %s lacks its fence", testCase.ID)
		}
	}
	if cases != 12 || accepts != 6 || rejects != 6 {
		t.Fatalf("Slice B manifest = %d cases (%d accept, %d reject), want 12 (6, 6)", cases, accepts, rejects)
	}
	for _, scenario := range []string{
		"first-usable", "exact-replay", "nonce-replay", "atomic-commit", "confirmed-abort",
		"indeterminate-pre-state", "indeterminate-post-state", "reopen-post-state",
		"usable-exact-256", "usable-cap-plus-one", "nonterminal-cap-plus-one",
	} {
		if _, ok := seen[scenario]; !ok {
			t.Fatalf("Slice B manifest omitted scenario %q", scenario)
		}
	}
}

func readManifestFixture[T any](t *testing.T, fixture manifestFixture) T {
	t.Helper()
	path := filepath.Join(conformanceRoot, fixture.Path)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", fixture.Path, err)
	}
	if fixture.ByteLength != 0 && len(data) != fixture.ByteLength {
		t.Fatalf("%s length = %d, want %d", fixture.Path, len(data), fixture.ByteLength)
	}
	if fixture.SHA256 != "" {
		digest := sha256.Sum256(data)
		if hex.EncodeToString(digest[:]) != fixture.SHA256 {
			t.Fatalf("%s digest mismatch", fixture.Path)
		}
	}
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("decode %s: %v", fixture.Path, err)
	}
	return result
}
