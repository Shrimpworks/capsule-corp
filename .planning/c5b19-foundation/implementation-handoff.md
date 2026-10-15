# C5b19 G3a2 implementation handoff

Status: PASSED for local implementation/verification; BLOCKED pending fresh independent G3a2 review. Parent full Gate3 IN_PROGRESS — TRENDING_GOOD; Gate4–6 and admission BLOCKED.

Internal worker modified existing lead-owned /tmp/capsule-c5b19-recovery branch codex/c5b19-benign-teardown from PR38 head1e3f3bb154d0ef24412002ac47adef22ba3c46a9. No branch, commit, push or PR created. No additional agents. Source is frozen for review.

## Defensive scope

Pure in-memory synthetic request/reply/record conformance against frozen Gate2 and C5b18 source, in the exact owned archive. No test code subprocess/filesystem/network/process/clock effects, native fixture, guest, child, signals, waits, credentials or services. Commands use normal Go test compilation/process only. No new requirements or go.sum, no downloads. No product import/admission.

## Change

- tests/oracle.go adds request397 and record395 decoders, copied wireOperation, candidate structural validation, exact frozen bindings/policy/obligation verification, request/record correspondence and reply98 tuple/body-hash/outcome/digest correspondence.
- Terminal records accept absent runner identity when publication failed but absence identity is nonzero. Failed/indeterminate replies require zero record digest. Returned outcome is byte correspondence only, never authoritative publication evidence.
- tests/oracle_test.go adds independent literal vectors and hashes, every-byte substitution (890), every short length/trailing byte, structural and frozen-binding refusals (10 zero +20 duplicate pairs), valid field variation, ownership and uncertainty cases.
- Original obligation_test.go deleted after exact tests consolidated into frozen-listed oracle_test.go. Existing model_test.go unchanged (identity tests, not a state-machine claim).
- PLAN.md retains red/green/mutation evidence, source inventory, BUILD-NARROWLY checklist and historical review provenance discrepancy. ORIGINS pins references/current source SHA256s.

## Exact offline verification

Working directory /tmp/capsule-c5b19-recovery/experiments/typed-guest-transport-c5b19-benign-teardown.
Prefix every Go command with GOCACHE=/tmp/c5b19-go-cache GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off and use /Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go.

- test ./... initially compiler-red on missing decodeRequest/decodeRecord/validateRecordForRequest/validateReply. /tmp/c5b19-wire-red.log retained, excerpts in PLAN.
- test ./... -count=1 passed; /tmp/c5b19-wire-green.log.
- test -race ./... -count=1 passed; /tmp/c5b19-wire-race.log.
- vet ./... and build ./... passed with no diagnostics; corresponding logs.
- git diff --check passed.
- Controlled mutant candidate-validation bypass: named structural tests failed, exit1. /tmp/c5b19-wire-mutant-candidate.log.
- Controlled request-body-hash bypass: uncertain reply body-substitution test failed, exit1. /tmp/c5b19-wire-mutant-request-hash.log.
- Controlled uncertain-record-digest bypass: all six candidate classes failed, exit1. /tmp/c5b19-wire-mutant-uncertain-digest.log.
- Exact source restored and hash checked; post-restoration test ./... -count=1 passed, /tmp/c5b19-wire-restored.log. No mutants retained.
- Initial wrong-relative-path command did not write tests; initial default build cache sandbox denial resolved using isolated /tmp cache before the actual compiler-red test. No source behavior preceded actual red run.

## Review boundary and limitations

Fresh review is for distinct G3a2 source/diff, including interactions with G3a; it is not a fourth retry or corroboration of old PR body review3 Ready claim. PLAN previously said pending, external report unavailable; both claims retained honestly. No claim that full oracle or state machine passes.

Stateless validation cannot establish exact predecessor generation, fresh operation IDs beyond preparation ID, non-overlap, late-settlement/control latch behavior, raw trace causality, C/Go equality, storage durability/faults, current process custody, timing, start or release. All deferred. Root owns final tests/review/publication and canonical exact-commit links.

## Frozen file SHA256

- `PLAN.md`: `eb102a9cb173978fb0070f6c8f6b9e8ecef003c4f8949c0a35d1979c6660cbd9`
- `inputs/ORIGINS.json`: `75ba6a28940bb66142b8539294f330dde5f29d2b0db1f1009ce2a5cf14f8b3d1`
- `tests/oracle.go`: `0c370da8d89c792d6b0efcfce42b64772c74e42f4c4bae45f1df4c6398e5ff7b`
- `tests/oracle_test.go`: `c505b8fea9df4059a503b683b916343a95261b2172436eaa0308fa4a108f15fe`
- `tests/model_test.go`: `faa1a4ed113c9c30a7528b9b1096d734b7134173b4cf05f07e3b84c0736e0a2e`
- `go.mod`: `5f40cde18c2d82fe5d57f3a8822377c12e473f1026f174165dfecd126bea3eeb`

Deleted: tests/obligation_test.go (coverage consolidated).
