# C5b19 G3a3 independent engineering review

Review instance: 1 of 3. Review status: PASSED. Verdict: Ready.
Parent full Gate3: IN_PROGRESS — TRENDING_GOOD. Gate4–6 and product admission: BLOCKED. Child execution remains separately gated.

## Findings

No actionable findings remain. No corrective source change is requested.

This review covers necessary synthetic write-history consistency, not real custody, authoritative absence, durable publication, timing, start/stop, public completion or capacity evidence. Final publication pins were explicitly pending at dispatch and are not declared complete here.

## Scope and independence

Applied the independent-review and codebase-memory skills. Recorded `/tmp/c5b19-g3a3-review-1-preliminary.md` after inspecting requirements, tests and implementation, before reading author rationale, archive PLAN/ORIGINS, logs or canonical outcome claims. The neutral bootstrap and acceptance recommendation were initially available; the implementation conversation was not.

Archive: `/tmp/capsule-c5b19-recovery`, branch `codex/c5b19-benign-teardown`, HEAD/base `2274a6c7c751ded537541438f473e81ccd293457`. Reviewed the explicit working diff in `experiments/typed-guest-transport-c5b19-benign-teardown/{tests/oracle.go,tests/model_test.go,PLAN.md,inputs/ORIGINS.json}`. No other archive changes appeared. Frozen hashes match:

- oracle.go: `fe36b0504a4a69623b1f64e791487e35ae78869c8e1524ff131d9790d3180dc2`
- model_test.go: `b0f3b9d063c225eb887cecaf16f4058f0ed07ae8f14c3d92c49acacabb0a931e`
- PLAN.md: `9ca93e5349e71c7be3fa346e980e7e88897d162f159be51533d19ca86e71fdaa`
- inputs/ORIGINS.json: `4649172116ce6275c7f4cf7b78e51c96b93de5bd22032d2f6fe68ec32f68c3fc`

Canonical: `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-settlement-conformance`, HEAD/base `080e9203c6c96d1494ea3a52e8855816fec3fa38`. Reviewed changes to docs/CURRENT_WORK_PLAN.md, docs/PROJECT_RECOVERY_2026_10_15.md and docs/C5B19_BENIGN_TEARDOWN_EXPERIMENT_PLAN.md, plus new `.planning/c5b19-settlement/{scope.md,plan.md,validation.json,implementation-handoff.md}`. The copied author handoff was inspected after the preliminary record.

No branch, commit, push, PR operation, source mutation or subagent was created. Reviewer writes are these two reports plus ordinary Go cache output. Fixtures inspected and exercised are inert in-memory values. No child fixture, guest, signal/wait target, live store, service, credential or user data was accessed.

## Plan Review

| Acceptance | Implementation and test evidence | Result |
| --- | --- | --- |
| Frozen valid bindings, owned copies, nil/zero refusal | oracle.go newWriteHistory/begin/applyReply; model_test.go TestWriteHistoryOwnershipAndInitialization and TestWriteHistoryRequestBindingAndRejectionOwnership | PASSED |
| Sole pending through explicit application | begin pending guard and private frame copy; TestWriteHistoryPendingIsNotWorkerArrival and otherwise-valid terminal replacement refusal | PASSED |
| Exact next generation and fresh IDs | begin nextGeneration/usedIDs guards; TestWriteHistoryGenerationFreshnessAndOrdering, including valid reused-ID terminals after every outcome | PASSED |
| Preparation first/confirmed; each kind once; terminal closes | begin attempted flags and prepared projection guard; positive matrix and ordering tests | PASSED |
| Exact settlement and preservation on refusal | applyReply validates privately retained request before mutation; TestWriteHistoryExactSettlementAndReplays; rejection helpers compare whole private state | PASSED |
| Confirmed-only projection | applyReply outcome guard; settleAndCheck and all outcome combinations | PASSED |
| Preserve confirmed runner facts | terminal candidate prefix comparison; TestWriteHistoryTerminalContinuity uses individually valid contradictory candidates | PASSED |
| Attempted identity survives uncertain publication | runnerAttempt and terminal absence comparison; all runner outcomes plus direct preparation-to-terminal path | PASSED |
| No hidden control state or lifecycle authority | fixed private arrays, value snapshots, narrow write API; no reset/import/retry/timeout/control methods | PASSED |

Read canonical Gate2 fixed-transport, durable-format and control-reply-application requirements independently, alongside normative teardownpassive/types.go and model.go. The implementation matches their write-lane relations while explicitly deferring custody/absence guards and real received/applied trace pairing. Decoder validation closes the kind range before array indexing. Three-write bounds make generation overflow unreachable; the decoder still enforces the safe-integer limit. A runner candidate extends the only reachable confirmed preparation state. Terminal continuity prevents inventing or erasing confirmed runner facts. Failed and indeterminate settlements consume history without promoting their candidate.

The plan orders scope, pure tests, bounded implementation, verification/mutations, review and publication appropriately. Six recorded mutations exceed the suggested four without expanding scope. No manifest, dependency, privileged responsibility or authority expansion occurs. Appended provenance preserves historical G3a/G3a2 records and limitations. Canonical documentation separates the bounded child increment from the unfinished parent.

## Author-Claim Reconciliation

| Author claim | Evidence inspected | Status | Review consequence |
| --- | --- | --- | --- |
| Nine bounded invariants implemented | Source/requirements/tests and independent full tests | Confirmed | Supports readiness |
| Private pending bytes and value snapshots | Value types, begin copy, applyReply has no replacement-request argument, ownership tests | Confirmed | Caller mutation cannot replace accepted bytes |
| No dependency, manifest, live-effect or permission expansion | Four-file archive diff, module, all test source/imports | Confirmed | Inert authorized boundary preserved |
| 24 runner/terminal/disposition cases and two uncertain preparations | 4×3×2 matrix and two preparation subcases | Confirmed | Positive/uncertain paths covered |
| Compiler-red before behavior | Hash-matching missing-symbol red log and retained PLAN excerpt | Corroborated; exact chronology not independently replayed | No contradictory evidence; not reviewer-observed author chronology |
| Six guard mutants failed and restored | Matching log hashes and intended assertion sites; restored source hashes; independent tests pass | Confirmed as retained sensitivity evidence; not rerun | Bounded sensitivity, not exhaustive mutation proof |
| Offline test/race/vet/build pass | Author logs and independent reruns | Confirmed | No local verification blocker |
| Canonical install/CI/test pass; unfiltered lint debt | validation.json hashes, CI output, 164+10+3 workspace tests,132 script tests, lint output | Confirmed as retained author results; not reviewer reruns | 50-output-capped revive debt stays separate from scoped CI |
| PR386 merged at080e920 | Local commit identity/message/date | Confirmed locally | Correct continuation baseline |
| Hosted success/protection/no-bypass details | Canonical recorded identifiers; no network query | Unverified independently | Orchestrator owns hosted readback |
| Full Gate3/run/admission incomplete | PLAN/ORIGINS and canonical diff | Confirmed | Ready verdict does not open later gates |

Recomputed all four implementation hashes, three canonical reference hashes, twelve archive log hashes and four canonical validation-log hashes: every value matched. Historical G3a2 fields are unchanged in the diff; this review does not reset or recertify them.

## Verification Performed

Working directory: `/tmp/capsule-c5b19-recovery/experiments/typed-guest-transport-c5b19-benign-teardown`.
Every reviewer Go command used this exact prefix:

```sh
GOCACHE=/tmp/c5b19-g3a3-review-cache GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off /Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go
```

- `test -count=1 ./...`: exit0, `ok capsule.experiments/c5b19/tests 0.297s`.
- `test -race -count=1 ./...`: exit0, `ok capsule.experiments/c5b19/tests 1.410s`.
- `vet ./...`: exit0, no diagnostics.
- `build ./...`: exit0, no diagnostics.
- `git diff --check` in both worktrees: exit0, no diagnostics.
- Read-only SHA256 checks of source, references and retained logs: all matched.
- Git branch/HEAD/status checks: declared frozen scope matched, no extra source/dependency/product files.

Graph limitation: list_projects was fully paginated (57 projects); exact archive checkout is unindexed. Canonical capsule-corp coverage generation2026-10-04T22:18:10Z reports metadata_match/no_recorded_issue for packet/types/model and teardownpassive scope, without further pages. This is best-effort, not completeness proof. Focused exact-checkout reads supplied authoritative evidence. Git diff confirms those canonical requirements unchanged between primary fc8af2df and target080e920. No archive index or source was mutated.

## Open Questions And Residual Risks

No implementation question blocks the bounded slice. Synthetic confirmation remains a claimed byte-level write fact. Actual arrival/application pairing, store uncertainty, custody/absence, stop/start races, time bounds, native/C conformance and full build/run closure remain deferred. Race-test success does not establish concurrent API safety or live control behavior; this is a sequential inert projection.

Raw logs are ephemeral; PLAN retains excerpts, commands, outcomes and hashes. Final commits, immutable archive links, review retention, PR metadata and hosted checks remain publication work. Broader canonical checks were inspected as author logs rather than rerun because this review's executable scope was the offline archive suite.

## Verdict

**Ready.** No actionable source, test or plan defect was found in the frozen G3a3 delta and canonical reconciliation. G3a3 review status is PASSED; full Gate3 remains IN_PROGRESS — TRENDING_GOOD and Gate4+ remains BLOCKED.

## Recommended Next Actions

The orchestrator should retain this report/hash, reconcile review-pending metadata, commit/publish the exact archive candidate, replace pending publication references with immutable pins, and perform required PR/body/draft/head/check readback. Preserve historical review records and deferred control/runtime gates. These are planned delivery steps, not implementation fixes or a new review workstream.
