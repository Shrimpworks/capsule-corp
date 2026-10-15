# CI265 independent review

Review instance: 1 of 3.
Work status: BLOCKED on a reproducible local Go documentation-gate false pass.
Verdict: **Not ready**.

## Findings

### P2 — Isolate golangci-lint cache by checkout before filtering new findings

Evidence: `scripts/lint-go.mjs:26–32` inherits a shared `GOLANGCI_LINT_CACHE` (or the normal global cache) for both scans. With pinned golangci-lint 2.12.2, cached issue positions retain another checkout's absolute source paths. The new `--new-from-rev` stage silently filters those findings out when running an identical module from another checkout.

Decisive reproduction: `/tmp/capsule-ci265-two-roots.mjs` creates two retained disposable repositories under `/tmp/capsule-ci265-two-roots-pOgDwR/`, both with identical `example.invalid/capsule_review_two_roots` module/source content, documented baseline, and a committed `func Added() {}` without documentation. Both use fresh writable shared Go/linter caches. The first wrapper run correctly exits 1 with `source.go:4:1` revive failure. The second wrapper run exits 0 and reports `0 issues` twice, with no warning. Both fixture directories remain present, so deletion, unavailable source, sandbox cache permissions and cleanup are not necessary for the false pass. Log: `/tmp/capsule-ci265-two-roots.log`.

Control: rerunning the second checkout with the SAME Go cache and a distinct `GOLANGCI_LINT_CACHE` correctly exits 1 on `Added`. Log: `/tmp/capsule-ci265-two-roots-isolated-lint.log`. Tool: `/opt/homebrew/bin/golangci-lint`, version 2.12.2, built with go1.26.2; runtime GOTOOLCHAIN=go1.26.9, Node v22.22.1.

Impact: contributors switching between clones/worktrees can get a green `make ci` for newly undocumented Go exports that clean hosted CI rejects. The old unfiltered local lint did not have this changed-line filtering failure, so the new wrapper exposes it. The unchanged-security negative test still fails; no claim of a general security-linter bypass is established.

Smallest credible correction: make the wrapper choose a linter cache scoped to the canonical checkout root (including any caller-selected cache root), or a disposable private linter cache per invocation, without clearing the user's shared cache. Add regression coverage that runs identical source through two retained roots sharing a parent cache and requires both undocumented changes to fail, plus repeat-run integration coverage. Merely isolating caches inside the tests would conceal the real wrapper defect. Update retained evidence after remediation.

The existing real-linter integration suite also fails on repeat runs because it creates/deletes equivalent modules with shared cache state (`scripts/ci-parity-go.integration.mjs:14–34`); this is the same root issue, not a separate finding.

## Plan Review

The plan faithfully addresses issue #265 mandatory items 1–6 and intentionally leaves optional Dependabot grouping out. The expanded helper/workspace-script/test/doc scope is justified by complete build deduplication and tested parity. No dependency, action pin, permissions, guest or authority changes were found.

The JSDoc scope/invalid-base/dirty-tree design, full security scan, standalone build-safe commands and serialized Make ordering are coherent. Exactly one emitted build per workspace package is visible in the actual aggregate log. All original package test selectors and root conformance commands remain present. No heavy architectural pivot is required to correct the finding.

Acceptance requiring reproducible changed-line Go enforcement is currently unmet. Publication and hosted job/build-count confirmation remain subsequent delivery steps; the plan/report correctly do not declare the overall issue complete.

## Author-Claim Reconciliation

| Author claim | Evidence inspected | Status | Review consequence |
| --- | --- | --- | --- |
| Local aggregate adds audit/site/ADR/JSDoc and emits once per package | Makefile, manifests, workflow, `/tmp/capsule-ci265-logs/final-make-ci.log` | Confirmed | No issue; three emitted and three noEmit invocations observed |
| Package tests/conformance checks are retained | Root/workspace manifest diff and standalone/aggregate logs | Confirmed | No dropped enforced suite found |
| Changed JSDoc handles empty, invalid, renamed, deleted and dirty files | Helper, verifier and focused 17-test run | Confirmed | Whole changed-file enforcement preserved |
| Go new-documentation enforcement is proven by pinned-linter tests | Prior 3/3 log, fresh-cache 3/3, repeated-cache 1/3, two retained roots | Contradicted as a general reproducibility claim | Cache relocation causes a real false pass; finding above |
| Unchanged correctness/security findings remain enforced | Whole-code wrapper stage and G402 integration case | Confirmed for tested case | No evidence to claim broader security coverage than tested |
| All required full audits pass | Author explicitly says unfiltered lint retains debt | Not claimed | Known #217 debt accurately distinguished |
| Hosted proof is complete | Author explicitly says pending publication | Not claimed | Must still collect hosted CI/build-count readback |

## Verification Performed

From `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`:

- `PATH=/Users/dsteele/.local/share/fnm/node-versions/v22.22.1/installation/bin:$PATH node --test scripts/ci-parity.test.mjs scripts/verify-exported-jsdoc.test.mjs`: PASSED, 17/17.
- `PATH=... pnpm lint`: PASSED, 159 files, no fixes.
- `actionlint .github/workflows/ci.yml`: PASSED.
- `git diff --check`: PASSED.
- `PATH=... GOTOOLCHAIN=go1.26.9 node --test scripts/ci-parity-go.integration.mjs`: failed 2/3; committed/dirty documentation case and staged-new-file case received exit 0. Repeated with authorized escalated normal-cache access: same failure, `/tmp/capsule-ci265-review-go-unsandboxed.log`.
- Same integration suite with BOTH fresh writable temporary GOCACHE and GOLANGCI_LINT_CACHE: PASSED 3/3; immediate repeat with same caches: failed 2/3. Logs `/tmp/capsule-ci265-review-go-fresh-cache.log` and `/tmp/capsule-ci265-review-go-reused-cache.log`. Debug copy only adds logging and fixes repository-root lookup: `/tmp/capsule-ci265-debug-go.integration.mjs`.
- Two-retained-root reproduction and isolated-linter-cache control: results above. Harness `/tmp/capsule-ci265-two-roots.mjs`.

Full `make ci` was not duplicated. Inspected lead's aggregate log and standalone `pnpm --recursive check && pnpm test && pnpm lint` log; their successful execution is consistent with source and build-count claims. A clean changed-Go scope in the real branch does not refute the cache-negative reproduction. The prior author 3/3 integration log is genuine but insufficient for repeated execution.

## Open Questions And Residual Risks

- Review is of explicit working-tree changes, not a final commit: branch `codex/ci-parity-265`, HEAD/base `9a495bc404a6f2dde67b34e72409a909d2964dd1`. Tracked and all untracked implementation/doc/plan paths match the bootstrap. Only lead evidence/doc append activity was allowed during review.
- Preliminary independent ledger was written before reading author rationale: `/tmp/capsule-ci265-review-preliminary.md`. No author implementation history was inherited.
- Code graph covered the older primary checkout, missing new files; direct worktree reads supplied evidence. No repository file edits, staging, commits, pushes or additional reviewer instances were performed.
- Authorized defensive scope was the contributor CI tooling and repository-owned disposable local fixtures only. No guest, network target, credential, identity or unrelated data was accessed. Retained `/tmp` fixtures/logs support the finding.
- Initial sandbox runs also revealed unwritable/missing normal Go cache entries; this confound was eliminated by the fresh writable-cache reproduction and escalated rerun. Do not misclassify that incidental cache-permission issue as the finding.
- Full documentation backlog and GitHub-only service checks remain outside local passing claims. Hosted publication/readback remains pending.

## Verdict

**Not ready.** Correct checkout-sensitive lint-cache handling and retain a regression before relying on the new local CI parity gate. The rest of the reviewed implementation has no actionable findings.

## Recommended Next Actions

1. Accept or dispute the cache finding with the retained two-root evidence; implement the smallest bounded cache correction in the existing slice.
2. Re-run focused/integration checks, including two retained roots and repeated runs, and update canonical evidence without claiming full-backlog completion.
3. Use fresh-context review instance 2 of 3 for the material changed behavior, then publish and verify hosted CI and emitted build counts. No heavy pivot or new workstream is needed.
