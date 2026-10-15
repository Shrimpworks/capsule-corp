# CI parity recovery slice — issue #265

Work item: CI265 implementation/verification. Status: PASSED; PR #385 pending merge.
Scope: contributor Make/CI command parity, unchanged security gates, one emitted TypeScript build
and one test-suite execution, reproducible exported-documentation checks and tool-version clarity.
Parent: initial recovery delivered; owner-only alpha IN_PROGRESS — TRENDING_GOOD, admission BLOCKED.
Base: 9a495bc404a6f2dde67b34e72409a909d2964dd1, merged PR #384 on 2026-10-15 after verified
head 0c879ad and green applicable CI. Branch codex/ci-parity-265 reuses the clean attached worktree.

## Delivery and ownership

All children are internal subagents in this worktree; Brain owns integration, commits and one PR.
- CI265-GATES: read-only gate inventory and edge-case audit; no retained implementation.
- CI265-IMPLEMENT: Makefile, root/workspace scripts, CI workflow, minimal helpers and regression tests.
- CI265-DOCS: Brain owns this record, Development/onboarding and current execution index.
- CI265-REVIEW: fresh-context independent reviewer after implementation/verification; maximum three attempts.

Wait for both child handoffs before final review. Three unsuccessful review attempts stop the cycle
and require focused research/alignment, not a silent reset.

## Acceptance and tests

1. Local CI covers TS dependency audit, site, schema, ADR, type/lint, changed-file JSDoc and tests.
   Explicitly distinguish GitHub-only archive/dependency-review/CodeQL services and nightly scheduling.
2. `pnpm verify:jsdoc` exposes existing verifier; explicit files and no-argument full audit remain
   understandable. Shared changed-file selection fails closed on invalid base, succeeds on empty
   scope without invoking full audit, preserves CI scope, handles filenames/deletions/renames and
   local staged/unstaged/new files. Positive and missing-doc temporary fixtures prove it.
3. Exactly one emitted TS build and one suite execution per TS CI job. Preserve standalone
   build-safe check/test commands; do not mistake `--noEmit` checking for a second emitted build.
4. Remove unused informational coverage rerun without dropping enforced tests or checks.
5. Biome schema and declared version range agree; current 2.5.12 floor/schema already match.
6. Local Go CI uses existing CI security/correctness and changed-line documentation policy;
   full unfiltered audit remains explicit. Pin/document golangci-lint 2.12.2 and verify its version.
7. Run before/after `make ci`, required AGENTS commands, requested recursive TS check/test/lint,
   focused red-green regressions and workflow lint. Review independently and verify PR CI/logs.

## Scope decisions

Issue items 1–6 are mandatory. Item 7 Dependabot grouping is an optional triage preference and is
not selected. No permissions/action-pin, immutable-archive execution-policy, product, guest,
installed service, signing or dependency-version changes. Minimum helper/test/package-script
additions are included because the user explicitly requests tested complete parity; deleting only
the final workflow build would leave package-test rebuilds. No broad CI framework or new dependency.
Historical recovery reports remain historical. Current instructions are updated to the final behavior.

## Progress

- PASSED: #384 merged and merge commit read back; new branch starts at the merge.
- PASSED: CI265-GATES read-only audit collected; exact gate inventory and edge cases reconciled.
- PASSED: baseline make ci exits 2 on known revive debt; prior CI emits nine times across three packages.
- PASSED: initial focused regressions 17/17, real-linter fixtures 3/3 and aggregate make ci.
- PASSED: standalone recursive TypeScript check/test/lint.
- BLOCKED: review instance 1 returned Not ready: shared linter caches preserve another checkout's
  absolute paths and can silently drop new revive findings. Two retained roots reproduce it.
- PASSED: cache correction frozen; red two-root reproduction, 18/18 focused and 4/4 real-linter tests.
  Default cache is per-checkout Git metadata; caller cache parents are namespaced. No cache deletion.
- PASSED: corrected full make ci; 177 workspace tests and 132 script tests; full audit still exposes #217 debt.
- PASSED: independent review instance 2 of 3 returned Ready; no remaining actionable findings.
- PASSED: implementation 80a6efeddda24308ebe595b2abe1825277be1be5 pushed; PR #385 published
  with body/newline and branch/draft readback. CI 38783476966 and CodeQL 38783476875 green.
- PASSED: hosted TypeScript job116589384544 emits once/package (3 total vs baseline9) and no coverage rerun.
- Evidence-only follow-up records hosted results; exact-head CI and ready state readback are the final
  publication checks. No independently reviewed implementation changed after instance2.
- Remaining delivery action: maintainer merge of PR #385; #265 closes on merge.
- Deferred: optional Dependabot grouping; full documentation debt #217; GitHub-only service gates
  remain distinct from local CI; owner-only alpha remains IN_PROGRESS — TRENDING_GOOD/admission BLOCKED.
- Integration destination: https://github.com/Shrimpworks/capsule-corp/pull/385 for all retained commits.
- Next action: finish exact-head readback and ready transition, then report completion to the user.
