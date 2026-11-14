# Gate-3 documentation reset independent review

Review instance: 1 of 3.

Work item: fresh-context, read-only review of the owner-requested Gate-3 delivery reset.
Status: PASSED for this review. Verdict: Ready.
Parent documentation delivery: IN_PROGRESS — TRENDING_GOOD pending Brain's retention/publication work. Gate-3 implementation: IN_PROGRESS — TRENDING_BAD at the reset; owner implementation authorization remains PASSED. Child execution, Gate 4 onward and product admission remain BLOCKED on their named implementation/review/owner-authorization dependencies.

## Findings

No actionable findings. The reset accurately distinguishes file presence, local validators, integrated implementation and actual-run evidence. No changed authority, ABI, dependency or native behavior was introduced by the reviewed documentation.

Reviewed boundary: working-tree target on `codex/c5b19-gate3-reset`, HEAD/base `a9e1d32afbb1fa1ef9c94175321c9c3bbb3d3361`, in `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`. In scope: new `docs/C5B19_GATE3_COMPLETION_PLAN.md`, four tracked documentation updates, and supporting `.planning/c5b19-gate3-reset/{plan,inventory,requirements}.md`. Unrelated `.serena/project.yml` remains excluded. Archive tree: `a9c457cb3b89a4ecd8c38684fd3d804d2089214c`, `experiments/typed-guest-transport-c5b19-benign-teardown`, read only. No reviewer branch, commit or PR was created.

The neutral bootstrap was read first. Preliminary assessment was recorded in `/tmp/c5b19-reset-review-1-preliminary.md` before reading `/tmp/c5b19-reset-author.md`. This was a blind first pass with respect to the author's explanation.

## Plan Review

The packet's twenty-path manifest independently resolves to six present and fourteen absent paths, with no unexpected present path. The new completion plan explicitly rejects treating those counts as a completion percentage. Its grouping accurately assigns native and store/bridge implementation as missing, existing tests/oracle as partial, and `go.mod` as complete only for its declaration role.

Source inspection supports the functional inventory. At the immutable archive pin, `tests/oracle.go:517` implements local shapes for all 29 event kinds while excluding full lifecycle semantics and some diagnostics. `tests/oracle.go:1201` implements only partial transport selections. `tests/oracle_test.go:3453` deliberately accepts empty/lone-kind8 selections, and `:3674` deliberately accepts lane-consistent selections whose separate pair/frame correspondence fails. Thus the reset correctly identifies composition and required-event completeness as new obligations rather than rebranding the helpers as a full verifier. Package `oracle` has no executable entry point; packaging is explicitly assigned to D4.

The four deliverables are coherent and finite at the planning level:

- D1 (`docs/C5B19_GATE3_COMPLETION_PLAN.md:80`) owns independent complete-attempt judgments, modeled transitions and timing, finite scenarios/mutations, and expected capacity/readiness/EINTR verdicts. Scenario-appropriate completeness permits legitimate preparation refusal without inventing child observations.
- D2 (`:115`) owns the frozen wire, real store implementation and sole-worker bridge, with simulated pre-run store effects and explicit uncertainty.
- D3 (`:139`) exercises actual control code through simulated syscall effects, preserves custody and stop ownership, closes timing/trace production, and consumes D1 scenarios rather than treating parser-cap tests as control-loop evidence.
- D4 (`:168`) closes fixture/harness/build/provenance and the full review packet, including executable composition, permitted negative tests, exact inventories and gated entry refusal.

The completion rule at `:212` requires every frozen obligation implemented and supported by permitted pre-run evidence, with blocking implementation findings closed; it does not move unfinished implementation into Gate 4. The requirements audit preserves detailed packet obligations that the shorter canonical plan groups by outcome. The packet remains controlling, so compact bullet summaries do not remove such details as worker pre-fault flush, descriptor premises or exact wait semantics.

Roadmap edits supersede stale helper scheduling while preserving historical evidence. The former same-operation ordering recommendation is absorbed into D1. The explicit four-outcome plan avoids both an arbitrary PR-count promise and an open-ended helper milestone chain. Historical audit GOOD trends are acknowledged as prior reports; the new canonical reset status is BAD for the broader unfinished implementation, while D1's reusable foundation remains a distinct scope.

Authority remains unchanged: all pre-run process/store effects are simulated; actual native custody/timing and durable publication are not credited. Ordered Gate 4 independent actual-byte review and separate owner first-run authorization remain required. New manifest/ABI/dependency/authority/host mismatches reopen applicable review and authorization (`:220`). Real suspend/restart, descendants, guest and installed-product admission remain excluded.

## Author-Claim Reconciliation

| Author claim | Evidence inspected | Status | Consequence |
| --- | --- | --- | --- |
| Six paths present; integrated implementation incomplete | Frozen packet manifest, immutable `git ls-tree`, oracle source and representative tests | Confirmed | Inventory supports the reset; no file-count completion claim |
| Four deliverables cover remaining packet work | Completion plan D1–D4, requirements audit and Gate-2 packet | Confirmed | Bounded acceptance structure is suitable |
| Existing helper scheduling is superseded without erasing history | Diff of current work plan, experiment plan and recovery report | Confirmed | New schedule is identifiable; provenance retained |
| No new run or product authority | Plan simulated-effects and Gate-4 clauses; unchanged ordered gates; ADR-0047 | Confirmed | Child execution remains BLOCKED |
| Packaging may expose a packet mismatch | `package oracle`, missing scripts/native paths, D4/reopening clauses | Confirmed | Explicit future decision gate rather than hidden design permission |
| Canonical checks passed except known full-lint debt | Supplied verification JSON, matching SHA-256 log hashes and log excerpts | Confirmed as supplied execution evidence | Not reviewer-executed tests; preserve lint limitation in delivery |
| No historical test/mutation recertification | Reset wording and archive provenance limitations | Confirmed | This review makes no historical rerun or custody claim |

## Verification Performed

Reviewer performed only read-only Git/source/metadata checks and wrote the two assigned `/tmp` review reports:

- `git status --short`, `git branch --show-current`, `git rev-parse HEAD`, scoped `git diff`, and `git diff --check`: target verified; diff check passed.
- `git ls-tree -r --name-only a9c457cb3b89a4ecd8c38684fd3d804d2089214c -- experiments/typed-guest-transport-c5b19-benign-teardown`, immutable `git show`, and focused source reads: manifest 20, present 6, missing 14, extra 0. Archive status clean.
- Read-only Python manifest comparison and local-link existence checks on the new plan/supporting planning files: passed; no missing local targets. Anchor rendering and remote-link availability were not independently tested.
- Codebase-memory index and path coverage checks: primary checkout index dated 2026-10-04; reset/planning documents missing from graph, archive paths outside project. Symbol query returned no archive symbols. Current source/immutable Git fallback supplies evidence; negative graph output was not used to infer absence or complete call-graph closure.
- SHA-256 verification of each log listed in `/tmp/c5b19-reset-verification.json`: all matched. Supplied commands: frozen `pnpm install`, `make ci CI_BASE=a9e1d32afbb1fa1ef9c94175321c9c3bbb3d3361`, `pnpm test`, and `git diff --check` exit 0; unfiltered `golangci-lint run ./...` exit 1 with 50 displayed revive issues. CI log reports no Go vulnerabilities and four package audit advisories. The orchestrator identifies the lint findings as pre-existing #217 debt; no Go source changes are in this reset.

No experiment test, native executable, fixture, child, live store, guest or external target was run. No upstream platform-contract revalidation was attempted; this reviews consistency with the frozen packet, not that packet's original design or future actual-byte admission.

## Open Questions And Residual Risks

The reset intentionally leaves implementation decisions to D1–D4, particularly serialized evidence composition and executable/test-injection packaging. Those may require packet review if the fixed manifest cannot accommodate the result. This is already an explicit reopening condition and does not block the documentation reset.

Full native custody, actual syscall/clock behavior, generated/linked-runtime closure, host compatibility and real-run mutation sensitivity remain unproven. Historical G3a third-review provenance remains unresolved and is not recertified. Existing full-lint debt is outside this documentation change.

Brain's publication note says the stacked reset PR will merge only after a separately authorized merge of #410, then retarget/recheck against main because CodeQL's pull-request trigger is main-only. That is appropriate publication sequencing, not permission to merge either PR. Publication metadata additions and eventual exact-commit remote checks remain Brain's responsibility.

Confidence: high for manifest accuracy, documented helper limits, packet coverage and unchanged authority; this is not exhaustive code correctness review or future implementation validation.

## Verdict

Ready for the documentation reset's retention and publication, with no actionable review findings. Review status PASSED does not complete Gate-3 implementation or authorize any experiment run.

## Recommended Next Actions

Brain should retain this report with the planning artifacts, publish the scoped documentation change with explicit verification limits and stack/merge order, and report its immutable commit/PR. Future implementation should proceed against D1's complete oracle/corpus exit and the frozen packet, preserving all Gate-4/first-run blockers. No further review instance is needed merely to consume the remaining budget; material changes require reassessment within the existing limit.
