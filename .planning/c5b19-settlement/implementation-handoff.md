# G3a3 implementation handoff

Status: PASSED for local implementation verification; BLOCKED pending fresh G3a3 independent review. Parent full Gate3 IN_PROGRESS — TRENDING_GOOD; Gate4+ and product admission BLOCKED.

Source frozen in shared lead-owned /tmp/capsule-c5b19-recovery branch codex/c5b19-benign-teardown, base2274a6c7c751ded537541438f473e81ccd293457. No branch/commit/push/PR actions or subagents. Scope report /tmp/c5b19-g3a3-scope.md implemented in full. Only four owned files changed: PLAN.md, inputs/ORIGINS.json, tests/oracle.go and tests/model_test.go under experiments/typed-guest-transport-c5b19-benign-teardown. Existing oracle_test.go unchanged; its literal fixtures/helpers are reused.

## Defensive scope and implementation

Controlled inert in-memory repository fixtures only. No test subprocess, filesystem, network, clock, goroutine, child/guest, signal/wait, service or credential effects. Normal Go compilation/test process only. No added dependencies, go.sum or source manifest paths. No control or permission stubs.

Constructor validates and copies immutable bindings. A fixed three-ID/three-kind history retains one privately copied397-byte pending request; applyReply accepts only reply bytes and validates against that retained request. The value snapshot distinguishes pending and last-settled presence, exact operation/outcome and confirmed68-byte projection. All refusal paths validate before mutation.

All nine accepted invariants covered: immutable bindings/owned copies; single pending through explicit application; next generation and globally fresh IDs; prep first and confirmed before optional runner/terminal, kinds maxonce and terminal closes; exact settlement/replay refusal; uncertainty consumes settlement but never promotes candidate; terminal preserves confirmed runner prefix; terminal absence matches attempted runner even after failed publication, while direct prep→terminal stays valid; no hidden control state. Nil/uninitialized receivers refuse operations. No importer/reset/retry/timeout/record-loading APIs.

applyReply models simulated accepted control application, NOT worker arrival or actual fsync confirmation. Positive byte history is not evidence of custody, timing, start/completion/capacity authority.

## Tests and TDD

Custom primitive BUILD-NARROWLY checklist was written in PLAN before behavior. Compiler-red /tmp/c5b19-g3a3-red.log shows undefined writeHistory/newWriteHistory before implementation. New tests cover24 runner/terminal/disposition combinations,2 failed/indeterminate preparations, valid varied bindings/process/operation IDs, ownership, zero/nil/invalid construction, sequencing/freshness/kind retries/overlap/fourth requests, exact settlement/byte mutations/candidate and cross-attempt substitutions/replays, and continuity after all outcomes. Rejection helpers compare entire private history; queued replies and repeated snapshots do not settle anything.

## Exact offline commands/results

Cwd /tmp/capsule-c5b19-recovery/experiments/typed-guest-transport-c5b19-benign-teardown.
Use /Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go with GOCACHE=/tmp/c5b19-g3a3-go-cache GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off.

- go test ./... -count=1: PASSED,0.276s.
- go test -race ./... -count=1: PASSED,1.312s.
- go vet ./...: PASSED,no diagnostics.
- go build ./...: PASSED,no diagnostics.
- git diff --check: PASSED.
- Six separately applied/restored mutants compiled and failed via intended assertions (exit1): pending overlap, next-generation equality, fresh-ID guard, uncertain projection promotion, attempted-runner terminal identity, exact pending settlement validation.
- Exact source SHA256 restored; post-restoration go test ./... -count=1 PASSED,0.201s.

Logs /tmp/c5b19-g3a3-{red,green,race,vet,build,restored}.log and /tmp/c5b19-g3a3-mutant-{overlap,generation,fresh-id,uncertainty,terminal-identity,settlement-binding}.log. PLAN retains command/results, red and mutant excerpts and every log SHA256 rather than relying only on ephemeral files. ORIGINS preserves historical g3a2 source hashes at reviewedc4d3c33 and adds separate current g3a3 hashes, canonical080e920 reference, log hashes and review-pending status.

## Review and limits

This is a distinct G3a3 review (max3), preserving all G3a/G3a2 prior review history. No old verdict is rewritten. New review should inspect full archive delta and current canonical plan. None of this closes full model/trace/native/store/build/run gates. Actual trace pairing, stop/start races, timing, creation/wait custody, storage faults and real publication evidence remain deferred. Root owns final review/commit/publication and immutable canonical links.

## Frozen changed-file SHA256

- `PLAN.md`: `9ca93e5349e71c7be3fa346e980e7e88897d162f159be51533d19ca86e71fdaa`
- `inputs/ORIGINS.json`: `4649172116ce6275c7f4cf7b78e51c96b93de5bd22032d2f6fe68ec32f68c3fc`
- `tests/oracle.go`: `fe36b0504a4a69623b1f64e791487e35ae78869c8e1524ff131d9790d3180dc2`
- `tests/model_test.go`: `b0f3b9d063c225eb887cecaf16f4058f0ed07ae8f14c3d92c49acacabb0a931e`
