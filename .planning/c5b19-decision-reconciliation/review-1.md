# Independent engineering review

Review instance: 1 of 3.

Review status: PASSED — scoped review completed. Delivery status: IN_PROGRESS — TRENDING_BAD pending one evidence correction. Full Gate 3 remains IN_PROGRESS — TRENDING_GOOD; child execution, Gates 4–6 and product admission remain BLOCKED.

## Findings

### P2 — Combined mutation does not establish the claimed mismatched-anchor sensitivity

Evidence: archive `experiments/typed-guest-transport-c5b19-benign-teardown/PLAN.md:1055–1059` says the combined bypass establishes mismatch sensitivity; `inputs/ORIGINS.json:275` repeats that claim. The author handoff's alternatives section makes the same claim. `/tmp/c5b19-g3a8-mutate.py:6` replaces the whole pending/equality/private-settlement block with `expected = 1`, removing normal settlement as well as both match checks.

On an exact preparation response, that mutant reports code1 but leaves preparation pending. Consequently all twelve `TestDecisionHistoryMismatchedAnchorsAndLateSettlement` subtests fail at `tests/oracle_test.go:1887` while beginning the runner, before the substituted anchor is constructed at1889–1900 or the intended false-accept assertion executes at1905. The retained mutant log explicitly shows all twelve failures at1887. Its SHA256 matches the retained value `9004068ee065b3a0f41db271ef93319e3f6a9112b26c35f60c9acfdc2ba33c4f`. The initial excerpt in PLAN correctly names the ordinary settlement test, but that does not substantiate the stronger narrative claim.

The mutation also triggers the no-pending false-accept assertion at1863. That is useful evidence for a different case; it does not show sensitivity to a wrong anchor while an actual different write is pending. The unmutated tests do exercise that case correctly, and no reconciler implementation defect was found. The defect is in the required adversarial evidence and its stated conclusion: scope.md requires each selected mutant to fail its intended behavioral assertion.

Smallest correction: preserve normal exact preparation settlement in a targeted combined-bypass mutant, bypass the exact-anchor/private-match protection for the wrong-anchor case, and require the specific wrong-pending false-accept assertion to fail. Retain that mutation definition and assertion evidence, restore exact source bytes, and refresh the affected hashes/excerpts/claims. Keep disclosure that isolated exact-anchor removal has overlapping protection. This is a bounded evidence correction, not a heavy pivot. Brain acknowledged `Accept` before this report; no fixes occurred during this review.

No other actionable findings identified in the frozen scope.

## Plan Review

The accepted plan adequately bounds this increment: private value-owned history; three represented decisions; strict serial-group sequence and tick watermarks; full-state atomic refusal; represented code2 history preservation plus a monotonic modeled fatal latch; and late settlement that preserves fatal. Individual G3a7 validation is the appropriate composition point because the collection API deliberately disallows repeated operation anchors needed for stale/duplicate rejection cases.

The code follows that plan. Existing history begin/settle rules are unchanged. Confirmed alone advances the projection, while failed/indeterminate outcomes clear the exact pending operation. Constructor and transition context validation, copied arrays, trial state, and value snapshots make the ownership/atomicity argument straightforward. Tests cover well-formed substitutions separately from malformed evidence and exercise zero/gaps/equal ticks/maximum values without arithmetic wrap.

The CI race correction is explicitly grouped before edits and is narrow. `runPassiveAdmitted` creates a background-rooted timeout and retains its only explicit cancel function. The worker publishes into the buffered result channel before deferred slot release and cancellation. Therefore the added `context.Canceled` branch can consume the published result safely. `DeadlineExceeded` retains the existing response-unknown disposition. Caller cancellation and simultaneous-ready arbitration remain unchanged. The new synchronization test and slot-lifetime regression support those distinctions.

The documentation correctly leaves full stop/start admission, receive/begin/apply interleaving, rejected malformed diagnostics, omitted-event detection, worker/snapshot correspondence, deadlines, actual durability, native execution and installed authority outside this slice. The next-scope recommendation keeps record visibility distinct from confirmation. No new dependency or privileged responsibility appears. The sole acceptance gap found is the targeted mutant claim above.

Publication/pinning, final metadata checks and hosted acceptance remain Brain-owned pending delivery steps; they have not been misclassified as completed in this frozen candidate. A review verdict does not close those steps or child gates.

## Author-Claim Reconciliation

| Author claim | Evidence inspected | Status | Consequence |
| --- | --- | --- | --- |
| Value-owned bounded atomic reconciliation | oracle.go775–875; unchanged writeHistory; nine new tests | Confirmed | No code finding |
| Rejection preserves history, latches fatal; late settlement preserves latch | reconcile; whole-state and projection assertions | Confirmed | Matches accepted bounded contract |
| Thirteen compiled mutants fail assertions | Mutation script, all raw logs and hash ledger | Confirmed as a broad statement | One failure is not at its claimed targeted boundary |
| Combined bypass establishes wrong-anchor sensitivity | Script and all mismatch-subtest failures at1887 | Contradicted | P2 finding; evidence correction required |
| Existing source/provenance preserved | Byte-prefix comparisons, unchanged go.mod/model_test.go, historical parsed ORIGINS equality | Confirmed | No prior contract rewrite |
| Source hashes match frozen handoffs | Four archive implementation hashes, three canonical reference hashes, both adapter hashes | Confirmed | Review targets exact claimed bytes |
| Adapter regression red before fix and repeated green | Red log, source reasoning, independent repetition/race checks | Confirmed | Narrow cleanup/deadline correction supported |
| Canonical make ci and pnpm test passed | Retained report/JSON and matching log hashes | Inspected evidence, not independently rerun | No claim of reviewer execution or hosted acceptance |
| Full control, native behavior and admission are deferred | Scope, plan, implementation comments and roadmap diffs | Confirmed | Child remains BLOCKED |

## Verification Performed

Defensive environment: owned offline inert archive bytes and ordinary local synthetic passive adapter harnesses only. No mutation script was executed by the reviewer, no native fixture/child/guest/external target was run, and no repository source was edited.

Canonical root: `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-decision-reconciliation`, base/HEAD `8e4fd22cc00cd5e961544f6a803974600dc0e751`.

Archive root: `/tmp/capsule-c5b19-recovery`, branch `codex/c5b19-benign-teardown`, base/HEAD `5d0e53cbccae3d4f587205df70253b776eae63e5`. Scope is the exact uncommitted changes listed by the bootstrap, including canonical planning files. Unrelated `.serena/project.yml` preserved.

Commands used `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off` and reviewer-specific caches. Archive commands ran at `experiments/typed-guest-transport-c5b19-benign-teardown` using `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go`, cache `/tmp/c5b19-review-1-archive-cache`:

- `go test -count=1 ./...` — PASSED, 0.595s.
- `go test -race -count=1 ./...` — PASSED, 2.285s.
- `go vet ./...` — PASSED, no diagnostic.
- `go build ./...` — PASSED, no diagnostic.

Canonical commands used `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.11.darwin-arm64/bin/go`, cache `/tmp/c5b19-review-1-canonical-cache`:

- `go test -race -count=1 ./internal/execution/authorityplane` — PASSED, 2.045s.
- `go test ./internal/execution/authorityplane -run 'TestPassive(CompletedResultSurvivesMethodContextCleanup|CallerCancellationRetainsMethodContextAndSlot|CancellationDeadlineStallAndResponseLossSemantics)$' -count=100` — PASSED, 2.327s.

Read-only Python comparisons passed: archive oracle/test/PLAN prefix equality against HEAD; unchanged go.mod/model_test.go; equality of every historical parsed ORIGINS value; all declared G3a8 implementation/reference hashes; all20 archive log hashes; both canonical adapter freeze hashes; and the two retained post-fix canonical log hashes. `git diff --check` passed in both worktrees. One ancillary literal search used the wrong root and reported a missing file; it was corrected with the absolute archive path and did not affect verification.

Serena initialized. Codebase-memory Tier2: list_projects paginated fully; search found runPassiveAdmitted, SettleWrite and samePending; SettleWrite outbound trace and snippet inspected; path coverage requested for relied-on canonical code/docs/planning and archive files. Main checkout index generation2026-10-04T22:18:10Z is not the reviewed worktree; metadata_match is best-effort, planning is missing, archive outside_project. Exact source/diff reads supply fallback. No whole-graph or exhaustive-impact claim.

Preliminary ledger `/tmp/c5b19-g3a8-review-1-preliminary.md` preceded detailed author handoffs, archive evidence rationale and verification reports. Required plan/scope contained short author-result summaries, so complete blinding to those summaries was impossible; this is disclosed in that ledger.

## Open Questions And Residual Risks

No unresolved architecture question or heavy-pivot gate identified. Evidence correction and final delivery remain. The adapter test synchronizes readiness, not randomized select choice; its observed original failure and repeated/race green runs support the regression without claiming deterministic branch selection. Reviewer tests ran on macOS; hosted Linux acceptance remains unverified here. Existing unfiltered lint debt#217 and dependency-audit limitations are carried by the canonical verification report, not resolved by this review.

This review does not authorize child execution, close full Gate3, certify actual persistence, or establish product admission.

## Verdict

**Not ready** for the frozen delivery's evidence claim. Implementation and tests appear sound within their bounded contracts, but the specific required mismatched-anchor mutation conclusion is unsupported by the executed mutant. Correct that narrow evidence gap before claiming completed review readiness.

## Recommended Next Actions

Brain should correct the mutant/evidence as described, preserve this instance1 report, refresh affected metadata, and determine the bounded follow-up review under the existing three-instance limit. No new workstream or architectural pivot is needed. Then finish immutable archive linkage, final metadata verification, shared PR publication and exact-head hosted checks. Retain child execution and admission as BLOCKED.
