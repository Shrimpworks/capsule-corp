# C5b19 Gate-3 continuation after PR #387

Status: PASSED for bounded implementation, local verification and independent
review. Brain owns final PR delivery. Parent full Gate 3 remains IN_PROGRESS — TRENDING_GOOD;
Gate 4–6, child execution and product admission remain BLOCKED.

## Objective and delivery boundary

Merge reviewed PR #387 and deliver the next bounded inert event-conformance
increment from the unchanged Gate-2 packet. PR #387 merged normally at
2aefb15ee45018bb1045e711c77a0c7bc9f4c54d (2026-11-13 America/Toronto;
GitHub timestamp 2026-11-14T00:23:03Z), matching reviewed head4ff7edda.
CI38820204405 and CodeQL38820204303 passed. No bypass or protection change.

Canonical branch: codex/c5b19-event-conformance in the existing recovery
worktree. Archive integration destination: existing draft PR38 on
codex/c5b19-benign-teardown, base6575898a32ff4a3d24ca4c757f42b1807798a372.
The previous temporary archive checkout was absent and was restored from that
published branch. Preserve unrelated .serena/project.yml changes outside delivery.

## Assignments and dependencies

- G3a4-SCOPE — internal read-only subagent: derive a bounded synthetic event
  acceptance envelope from the exact packet; return scope/source evidence.
- G3a4-IMPLEMENT — internal worker after scope acceptance: test first, implement
  only approved frozen-manifest files, retain red/green and restored mutations.
- G3a4-VERIFY — internal verification worker: canonical required checks, existing
  lint debt and dated failures reported separately; no dependency changes.
- G3a4-DOCS — Brain: canonical status, source pins, retained evidence and PRs.
- G3a4-REVIEW — fresh-context independent subagent after source freeze, at most
  three instances for this distinct delta. Preserve historical review counts.

Brain collects all handoffs, reconciles findings and owns commits, pushes and PRs.
No worker branches, commits or publishes independently. Scope analysis and baseline
checks can run independently; implementation and review are sequential.

## Defensive scope and gates

Defensively validate rejection of malformed synthetic event bytes in the owned
archive's local in-memory fixtures. No child fixture/guest, lifecycle target,
signal/wait, actual store, clock, network, credential or user-content operation in
new code/tests. Normal compiler/unit-test tools are permitted. No source-manifest,
ABI, native host profile, dependency or authority expansion. Exact cached Go1.25.13
runs offline; canonical tooling remains separately pinned. Full causal/timing and
implemented-byte/host review remain separate gates; first child run requires owner
authorization after Gate 4.

## Acceptance and verification

Scope must state what syntactic acceptance proves and what it does not. Retain
compiler/behavior-red evidence, positive/negative vectors, meaningful targeted
mutant failures and byte-exact restoration. Run archive test/race/vet/build and
diff checks, canonical required gates, hash/source/provenance reconciliation,
and independent review. Publish all retained commits; read back PR body/newlines,
head/base/draft and hosted results. New canonical PR is intentionally ready only
after review; archive PR38 remains draft for unfinished full Gate 3.

## Progress

- PASSED: #387 normal pinned-head merge and immutable archive restoration.
- IN_PROGRESS — TRENDING_GOOD: scope assessment and canonical verification.

## Dated canonical prerequisite

IN_PROGRESS — TRENDING_BAD: post-merge CI43533822264 failed the unchanged
vulnerability gate on Go1.26.9: GO-2026-6859 (net/http, fixed1.26.11) and
GO-2026-6809 (crypto/x509, fixed1.26.10). CodeQL43533822342 passed.
The verifier reproduced the scanner failure locally. The smallest supported
same-minor fix is Go1.26.11; official [release history](https://go.dev/doc/devel/release#go1.26.11)
records its 2026-11-10 security release, and the
[x509 advisory](https://pkg.go.dev/vuln/GO-2026-6809) confirms the older affected range.

G3a4-VERIFY also owns the minimal canonical active toolchain pin/docs update,
with verified toolchain download and full checks. Preserve historical evidence,
scanner enforcement, and frozen archive Go1.25.13. No product behavior or dependency
addition is intended. This prerequisite joins the same independent review and PR;
no unsupported host-default upgrade to Go1.27 is used.

## Accepted G3a4 boundary

PASSED: scope handoff collected and accepted. The [scope report](scope.md) defines
exact 256-byte envelope offsets, computed Bindings[240] hash, expected nonzero
timebase, known code ranges, reserved zeros and owned decoded values. Finite lane
checks enforce explicit expected lane/first sequence, nonempty whole frames,
32768/1024/64 caps, overflow-safe sequence continuity and private control tick
watermark. Sequence origin is a trusted synthetic-fixture parameter because the
packet does not choose zero versus one; no capture completeness is inferred.

Per-kind payload/lane ownership, sidecars, causal and timing validity remain
expressly deferred. Syntactic envelope success cannot settle history or grant
execution/completion authority. No ABI convention is invented. The implementation
worker owns only existing oracle.go/oracle_test.go and PLAN/ORIGINS archive paths.
IN_PROGRESS — TRENDING_GOOD: implementation dispatched, with test-first evidence,
five targeted mutants and restoration; fresh review follows a frozen candidate.

- PASSED: canonical verification handoff collected; exact Go1.26.11 make ci and
  standalone pnpm test pass, including177 workspace/132 script tests. Brain checked
  seven command-log hashes and five active-file hashes. Scoped lint passes;
  unfiltered #217 revive debt remains BLOCKED (50-output cap, not backlog total).
  Dependency audit has2 low/2 moderate findings below its unchanged high threshold.
  Local Node22.23.3 differs from documented22.22.1; hosted workflow uses its pin.

- PASSED: archive implementation and verification handoff reported; four-file
  candidate frozen. Brain independently verified4 source/12 log/3 reference hashes.
  Exact offline test/race/vet/build and six intended mutation failures/restoration
  passed. Known-answer/code/binding/length/cap/sequence/watermark cases are covered.
- IN_PROGRESS — TRENDING_GOOD: fresh review instance1/3 against the explicit
  working diffs and canonical patch prerequisite. Source stays frozen during review.

## Accepted result and integration destinations

- PASSED: fresh review instance1/3 **Ready**, no actionable findings. Report and
  preliminary ledger retained; no correction loop required. Source/tests stayed
  unchanged after freeze; only acceptance metadata was appended before commit.
- PASSED: archive commit `3e26b809a40832a3351fdc6b5e0ffc1ef0094792` pushed to existing
  draft PR38. Historical source/provenance preserved, execution/admission false.
- PASSED: every internal handoff collected; no retained child work lacks a
  documented destination. Canonical branch codex/c5b19-event-conformance receives
  the exact pin, Go patch, roadmap and evidence in one new PR. Its final URL/head/
  hosted checks are recorded in the PR body after publication.
- Deferred: per-kind fields/lane ownership, paired frame sidecars, causal/timing
  oracle, native/store/build closure and all live evidence. Sequence writer origin
  still needs a later explicit choice; parameterized framing is not capture proof.
  Full Gate3 IN_PROGRESS — TRENDING_GOOD; child execution/Gate4–6/product admission BLOCKED.
- Next bounded scope after acceptance: kind-specific field/lane conformance, then
  sidecar correspondence and causal/timing composition under the same no-child gate.
