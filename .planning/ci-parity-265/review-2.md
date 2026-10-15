# CI265 independent review

Review instance: 2 of 3.
Review status: PASSED. Parent delivery status: IN_PROGRESS — TRENDING_GOOD; publication and hosted verification remain.
Verdict: **Ready** for publication and final hosted verification.

## Findings

No actionable findings remain in the reviewed implementation. The earlier P2 cache-path defect is corrected in the production wrapper and exercised by a failure-sensitive retained regression. No further review instance is required for the frozen implementation.

Target: `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/ci-parity-265`, HEAD/base `9a495bc404a6f2dde67b34e72409a909d2964dd1`, complete explicit tracked working-tree diff and untracked implementation, tests, report, plan and prior-review record. Source remained frozen during review; the lead updated progress/evidence prose. No unrelated changes were found. The review does not certify a subsequently changed commit.

## Plan Review

The plan covers mandatory issue #265 items 1–6. Make now runs dependency audit, site build, schemas, ADRs and changed-file JSDoc alongside the existing type/lint, Go correctness/security/vulnerability, build and test gates. The aggregate orders prerequisites correctly even under parallel Make. Actions uses an explicit build followed by built-output check/test entrypoints, retaining package selectors and root script/conformance checks. The discarded informational coverage rerun is removed. Standalone root check/test remain build-safe. Existing Biome schema and declared dependency floor are both 2.5.12, meeting the issue's permitted-range acceptance without a dependency change.

Go lint validates version 2.12.2 and configuration, runs all-code correctness/security checks, and restricts only revive to changed lines. Invalid bases fail, and untracked Go files require staging rather than being silently omitted. Full lint remains available. The checkout-specific cache is a bounded correction within the plan, not an architecture change. SHA-256 of canonical checkout root selects reusable child cache entries under either the supplied parent or per-checkout Git metadata; no shared cache is deleted.

The helper/workspace/test/doc expansion beyond the issue's older file list is explicitly authorized by the handoff and necessary for tested parity and complete deduplication. Optional grouping, action pins, permissions, archive policy, dependencies and product authority remain unchanged. Current documentation correctly distinguishes whole-file JSDoc from changed-line Go documentation, full audit debt, dirty-tree additions and hosted-only checks.

Publication and actual hosted after-change build-count evidence remain the plan's final delivery actions. The plan's progress/next-action text should be advanced to this review result and completed local aggregate evidence when the lead records this report; this is routine evidence bookkeeping, not an implementation defect.

## Author-Claim Reconciliation

The preliminary ledger `/tmp/capsule-ci265-review2-preliminary.md` was written and sent before opening `/tmp/capsule-ci265-author-2.md`, the author report or review-1. The required neutral bootstrap/plan disclosed the existence of an earlier cache defect, so this was an independent reconstruction rather than blindness to that defect's existence.

| Author claim | Evidence inspected | Status | Review consequence |
| --- | --- | --- | --- |
| Missing local gates now run and aggregate emits once per package | Makefile, workflow, all four manifests, corrected aggregate log | Confirmed | Three emitting compiler invocations plus three no-emit checks observed |
| Existing enforced suite/conformance selectors are retained | Manifest diff, corrected aggregate and standalone logs | Confirmed | 177 workspace tests and 132 script tests in corrected aggregate; no omitted original command found |
| Changed JSDoc preserves CI scope and fails closed | Git helper, verifier, independently executed 18 focused tests | Confirmed | Whole-file packages-wide scope; invalid/empty/dirty/rename/deletion/unusual-name cases covered |
| Local Go gates preserve security/correctness and new documentation | Wrapper and independently executed real-linter tests | Confirmed for reviewed policy and fixtures | Unchanged G402 fails, old documentation is excluded, new committed/staged/unstaged documentation fails |
| Accepted P2 corrected by checkout cache isolation | Wrapper, unit tests, red log, independent 4/4 real-linter run | Confirmed | Both retained roots reject new undocumented exports on two rounds sharing a cache parent |
| Full audits are green | Author does not claim this; full-lint logs retain revive debt | Not claimed | Existing issue #217 remains accurately separated |
| Hosted final evidence is complete | Author explicitly leaves it pending publication | Not claimed | Lead must inspect exact final-commit CI and build counts before declaring issue complete |

The retained red cache log specifically fails at round 1, second checkout with two `0 issues` reports. The corrected integration test asserts both nonzero exit and the expected `Added` documentation diagnostic; an unrelated process failure cannot satisfy it. The correction is in the wrapper rather than only in fixture environment setup.

## Verification Performed

From the target worktree, using `login:false`, Node v22.22.1 via `PATH=/Users/dsteele/.local/share/fnm/node-versions/v22.22.1/installation/bin:$PATH`, and `GOTOOLCHAIN=go1.26.9` where applicable:

- PASSED: `node --test scripts/ci-parity.test.mjs scripts/verify-exported-jsdoc.test.mjs` — 18/18.
- PASSED: `node --test scripts/ci-parity-go.integration.mjs` — 4/4, with authorized normal host cache access. Includes two retained checkouts/repeated runs, untracked refusal then staged detection, new documentation states and unchanged G402.
- PASSED: `pnpm lint` — 159 files, no fixes.
- PASSED: `actionlint`.
- PASSED: `git diff --check`.

Inspected rather than redundantly executed:

- `/tmp/capsule-ci265-logs/corrected-make-ci.log`: complete corrected aggregate, reported exit 0 by lead; audit/site/ADR/JSDoc commands present, full Go results, exactly three emitted compiler calls and three noEmit checks, 177 workspace and 132 script tests, conformance tail completed.
- `/tmp/capsule-ci265-logs/standalone-ts.log`: standalone recursive check/test/lint completed successfully.
- `/tmp/capsule-ci265-logs/before-typescript-ci.log`: historical hosted baseline has nine emitted compiler calls and three noEmit checks, consistent with the report.
- `/tmp/capsule-ci265-cache-red.log`: retained pre-correction two-root false-green reproduction.
- `/tmp/capsule-ci265-logs/corrected-full-lint.log`: full wrapper still returns failure on known revive documentation findings. This is not a passing full-audit claim.

## Open Questions And Residual Risks

- Hosted after-change CI results and exact emitted build counts remain unverified until publication; local source/tests cannot prove hosted service execution. This Ready verdict authorizes proceeding to that planned verification, not declaring the parent issue complete.
- Full documentation debt remains tracked in #217. Network advisory state, GitHub event baselines, service gates and scheduling can differ from a local run; documentation states those limits.
- Cache isolation retains per-checkout entries and incurs initial warmup for each checkout. Tested failure handling covers a non-directory cache parent; no broad filesystem-permission guarantee is inferred.
- Code graph metadata belongs to the primary checkout and marked four new files missing. Direct exact-worktree source reads supplied authoritative review evidence; no negative structural claim relies on incomplete graph coverage.
- Defensive scope stayed within contributor checks and disposable local Git/Go fixtures. No guest, credential, identity, unrelated service/data or user cache deletion was involved. No source edits, staging, commits, publication or additional reviewer instances were performed by this reviewer.

## Verdict

**Ready.** The complete frozen implementation and plan have no remaining actionable review findings. The accepted P2 is corrected and independently verified. This is a contributor-tooling review, not a product security/admission claim.

## Recommended Next Actions

1. Retain this report and update plan/report progress with the completed local evidence.
2. Commit and publish the intended PR, then confirm applicable hosted checks on the exact head and inspect TypeScript logs for one emitted build per workspace and one enforced suite execution.
3. Record hosted evidence and remaining #217/hosted-only limitations before closing the parent delivery. Do not start another review instance unless implementation materially changes.
