# C5b19 G3a10 supplied worker-received/sent sidecar correspondence

Status: IN_PROGRESS — TRENDING_GOOD. Full Gate 3 remains IN_PROGRESS — TRENDING_GOOD;
child execution, Gates 4–6 and product admission remain BLOCKED.

## Objective and delivery boundary

Merge reviewed PR #407, then deliver only bounded supplied worker-received/sent
sidecar correspondence with tests-first evidence, fresh independent review and a
canonical roadmap PR. PR #407 merged normally at
`fe9009d5398a883735219e2a3a33c351c79aa54c` on 2026-11-14T02:55:46Z
(2026-11-13 America/Toronto), matching reviewed head
`f61e2a0dfdf7cb89b35747d1dc598c59ee5d8070`. Exact-head CI 43548608076 and
CodeQL 43548608135 passed; no bypass or protection change. Post-merge CI 43548933421 and CodeQL 43548933564 passed.

Canonical branch `codex/c5b19-worker-sidecars` in the recovery worktree, base
fe9009d5398a883735219e2a3a33c351c79aa54c. Archive branch
`codex/c5b19-benign-teardown`, base 7646d8faabb670ce8b5ff1b213c414d91deda2d9,
in /tmp/capsule-c5b19-recovery; existing draft PR #38 is the integration destination.
Brain owns commits, pushes, PRs and integration. Preserve unrelated .serena/project.yml.

## Internal assignments and dependencies

- G3a10-SCOPE: read-only internal explorer; derive source-grounded bounded input/
  optional-member/collection contract and tests, preserving partial evidence.
- G3a10-VERIFY: internal verification worker; canonical baseline and evidence.
- G3a10-IMPLEMENT: internal worker after explicit scope acceptance; own only archive
  oracle.go/oracle_test.go, PLAN.md and new ORIGINS.g3a10. Tests first, minimal pure
  composition, selected compiled mutants reaching intended assertions, exact restoration.
- G3a10-DOCS: Brain; canonical roadmap, immutable pins and collected handoffs.
- G3a10-REVIEW: fresh-context read-only reviewer after freeze; max 3 instances for
  this distinct bounded slice. No reset of historical budgets/evidence gaps.

Workers share Brain branches/worktrees, preserve others' changes and create no
branch, commit, push or PR. Scope and baseline may run in parallel; implementation
waits for accepted scope. Brain reconciles all handoffs and publishes shared results.

## Defensive boundary and acceptance

Only owned inert synthetic byte fixtures and ordinary offline unit/compiler
processes. No native fixture, child/guest, lifecycle/signal/wait target, live store/
clock/network, credential or user data. Canonical routine checks may use their
existing advisory sources. Keep archive Go 1.25.13 frozen offline, all prior APIs/
source/test prefixes/provenance, manifest/packet/ABI unchanged. No new dependency,
serialization, privileged responsibility or product import. C5b19 reuse row remains
test-only ADOPT-PLATFORM / BUILD-NARROWLY.

The current kind 5/7 Request is decoding context, not independent worker receipt.
Define separate supplied worker-received 397-byte and worker-sent 98-byte slots,
matched to actually retained submitted requests, strict kind 5 worker reply and any
supplied kind 7 control-received reply. Preserve received-without-sent and sent-
without-control evidence as partial, not completed transport. No worker-request-
received event kind exists: do not invent one. Worker clocks cannot reorder control.
No supplied sidecar can settle history, promote snapshot visibility to confirmation,
or establish capture authenticity/completeness, syscall success or execution authority.

Acceptance: bounded explicit optional-member/cap/duplicate contract; red/green,
malformed/substitution/partial/nonmutation/boundary tests; compiled mutants reach
intended assertions after valid setup and restore exact source; offline test/race/
vet/build; canonical gates; fresh review; exact archive pin; committed/pushed PR
destinations; real-newline body and base/head/draft readback; exact-head hosted checks.
Full control/causal/timing/native and first-run gates remain deferred.

## Retrieval context

Serena session 448553d9 initialized. Codebase-memory Tier 2 index ready at
2026-10-04T22:18:10Z, primary checkout, old/not watched. SettleWrite lookup,
four-callee outbound trace and snippet complete: exact pending/outcome validation
precedes projection. Packet/model/types/ADR47 report best-effort metadata_match;
prior next-scope missing freshness and archive outside_project use exact reads.
Current packet 350–359 requires separate worker frame slots; 432–435 forbids total
cross-lane timestamp order. Archive oracle 622–626 explicitly separates decode
context from worker receipt. No completeness claim. Scope must ground further choices.

## Brain scope acceptance

PASSED: accepted the scope agent's three-row contract before implementation.
Initialized context is required, including empty selection. Select at most three
unique retained operations. Each entry requires an independent 397-byte worker
receipt; optional 98-byte sent reply and strict kind 5 capture appear together.
Optional strict kind 7 requires that pair. Thus supported rows are received only,
received plus sent plus kind 5, and those three plus kind 7. All supplied frame
bytes must agree with the actual retained request, and a supplied sent reply must
agree with any retained accepted reply. Omitted evidence stays unknown. Reject
invalid rows, malformed or unsubmitted bytes, duplicates and cap + 1 atomically.
Queries return owned value facts with explicit presence and separate accepted and
worker outcomes; they neither mutate history nor return authority/completeness.
Reuse G3a6's explicit no chronology/causal/timing claim; historical queries do not
consult the current control watermark. Worker authorized to begin tests-first.
Written scope report will supplement this contract; any conflict must be reconciled.

PASSED: canonical baseline handoff collected in canonical-verification.md and
validation.json. Required checks pass; known #217 unfiltered lint debt remains
BLOCKED. Canonical source/config remains unchanged. Exact current roadmap reads
supplement the old primary graph; recovery/plan freshness is missing, other cited
roadmap files report only best-effort metadata match. No completeness claim.

PASSED: written scope handoff collected as scope.md and reconciled with the
accepted contract; no conflict. WorkerSent is absent only when nil; nonnil empty
sent/capture values are malformed. Results preserve selection order and accepted
outcome separately from worker presence/outcome, with an all-zero result on refusal.

PASSED: next-scope handoff collected in next-scope.md. Recommended next bounded
increment is supplied request-write and worker-reply-write result correspondence
(kind 3→4 and kind 5→6). It is deferred: errno representation and the complete
consuming contract must be frozen before a later implementation. No next-slice
implementation or child execution is authorized by this recommendation.

PASSED: implementation handoff and evidence collected and frozen. Ten new test
groups and 98 appended implementation lines; offline test/race/vet/build passed.
Fourteen corrected compiled mutations reached 17 exact intended assertions, with
byte-exact restoration. Two invalid constructions are disclosed and excluded.
The four-file archive patch is 59479399ae70aaec90bd3a15c6304151ca651f2859286bb4cf3a65efa4eb8642.
Source/test hashes and preserved-prefix/provenance checks are retained in the
handoff and implementation-evidence.json. Fresh independent review instance 1 of 3
is next; source/test bytes are frozen until findings are reconciled.

PASSED: independent review instance 1 of 3 returned Ready with no actionable
findings. Brain accepts the report and its limits; source/test bytes remain frozen.
Report SHA256: e7ac4de9d08d2956db21730e47d15aeefc1b830995277d4e5ff88134a1af7620. All internal scope, baseline, implementation, next-scope
and review handoffs are collected. Final metadata, archive pin and publication
remain IN_PROGRESS — TRENDING_GOOD. Child execution remains BLOCKED.

PASSED: archive commit `d698bd9332a000ad3c2b9f23f87d0353dc5964b4` pushed to the existing shared draft PR #38
branch. Reviewed source/test hashes match; acceptance metadata records review
without changing prior provenance. Canonical roadmap pins this immutable archive
head. All retained changes have shared PR integration destinations; final canonical
metadata checks and publication remain IN_PROGRESS — TRENDING_GOOD.

PASSED: final explicit JSON check, staged-file pnpm lint, documentation site build
and staged diff check. Final hosted acceptance will be recorded in the PR body.
- json: SHA256 `c943f80d97ee27c0306b6fa62ad39ec0e2f32d37d6e2c4c94958bb9b75b08254`.
- lint: SHA256 `4ba38fa88c192bd0845959003eee6c0aedb11ffa7adf8a2e755491391f5c798d`.
- site: SHA256 `198149bc8cd29d5959ebe14aaa0721396f6a4b1bc52d5737f767ca97aaa9b464`.
- diff: SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.

No source/config changed after the canonical baseline. Review acceptance metadata
and immutable pins do not alter the reviewed source/test bytes. Known #217 debt
and audit limits remain disclosed; child execution and admission remain BLOCKED.
