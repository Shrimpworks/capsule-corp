# C5b19 G3a5 kind-specific field and lane validation

Status: PASSED for bounded implementation, local verification and independent review.
Brain owns final PR delivery. Full Gate 3 remains IN_PROGRESS — TRENDING_GOOD;
child execution, Gate 4–6 and product admission remain BLOCKED.

## Objective and delivery boundary

Merge reviewed PR #402 and deliver the next bounded inert event-conformance slice:
kind-specific allowed/required fields and lane ownership under the frozen Gate2
packet. PR #402 merged normally at fe3d22c84b5d41e0a4edbfacb531563040e16995 on
2026-11-13 America/Toronto (2026-11-14T00:50:43Z), matching reviewed head b381f683.
CI43535826940 and CodeQL43535826875 passed; no bypass or protection change.

Canonical branch codex/c5b19-event-fields in the existing recovery worktree.
Archive destination remains draft PR38, branch codex/c5b19-benign-teardown,
base3e26b809a40832a3351fdc6b5e0ffc1ef0094792. Preserve unrelated .serena/project.yml
changes outside delivery. Brain owns integration, commits, pushes and PRs.

## Assignments and dependencies

- G3a5-SCOPE — internal read-only agent: derive exact 29-kind field/lane matrix,
  separate local validation from sidecar/causal/timing checks; flag ambiguity.
- G3a5-IMPLEMENT — internal worker after acceptance: tests first, minimal pure
  implementation in existing frozen-manifest files, controlled mutation evidence.
- G3a5-VERIFY — internal verification worker: canonical required checks and exact
  archive integration evidence; preserve known lint/audit limitations.
- G3a5-DOCS — Brain: canonical reconciliation, exact source pins and PR evidence.
- G3a5-REVIEW — fresh-context independent review of frozen implementation and
  plan, maximum three instances for this distinct slice; no historical reset.

Scope analysis and canonical baseline may run independently. Implementation waits
for scope acceptance; review waits for source freeze. No worker creates branches,
commits or publishes independently. Brain collects every handoff before delivery.

## Defensive environment and acceptance

Defensively validate kind-specific field/lane rejection using owned archive inert
in-memory synthetic bytes only. No child fixture, guest, lifecycle target,
signal/wait, real clock/store/network/credentials/user-data effects in code/tests.
Ordinary compiler and unit-test tools are authorized. Exact offline Go1.25.13
archive pin, frozen source manifest, ABI and host packet remain unchanged.
No dependencies, new privileged responsibilities or product imports.

Acceptance: source-grounded explicit matrix; failing tests before behavior;
positive/negative/boundary tests and meaningful restored guard mutants; exact
archive test/race/vet/build and canonical required checks; independent review;
immutable archive pin, all retained commits pushed and linked to PRs; stored body/
newlines/head/base/draft and hosted-check readback. Envelope or local field validity
must not be presented as paired-frame, causal, timing, lifecycle or execution proof.

## Accepted local-validation boundary

PASSED: [scope matrix](scope.md) accepted for all 29 table-defined event kinds.
Implement an opt-in raw-byte decoder that composes the existing envelope decoder
with forbidden-zero fields, mandatory operation/identity presence, lane ownership,
reply outcome/digest coupling, store sites and the necessary wait-result/PID relation.
Existing G3a4 syntactic acceptance stays unchanged. Scalar zero is not treated as
missing evidence; zero wait status is permitted. No syscall semantics are guessed.

Support only the payload-free harness29 subset; the packet's optional harness-child
identity encoding has no exact slot mapping here. Complete represented replies are
supported for both kind8 decisions; arbitrary malformed-reply diagnostics remain
unsupported. These limits do not declare those variants invalid under the whole
packet and do not close the full oracle. Sidecars, pairing, causal/timing evidence
and all live execution remain deferred. Store27 outcome has no record digest.

IN_PROGRESS — TRENDING_GOOD: implementation follows the exact matrix in existing
oracle.go/oracle_test.go plus archive PLAN/ORIGINS; tests first, controlled mutations
and exact restoration. Source freeze precedes independent review.

PASSED: canonical frozen install, make ci and standalone pnpm test on Go1.26.11;
177 workspace and132 script tests. Brain verified four command-log hashes and
explicitly formatted/checked the retained JSON with Biome; its parsed values are
identical to the worker handoff. Final metadata additions will be checked again.
Known #217 unfiltered revive debt remains BLOCKED; low/moderate audit findings and
local Node patch difference remain disclosed in the verification report.

PASSED: post-merge main CI43536982093 and CodeQL43536982271 completed successfully
on fe3d22c84b5d41e0a4edbfacb531563040e16995. These are push results, not a new
nightly result. Branch protections and run/admission gates are unchanged.

PASSED: implementation handoff collected; four archive files frozen. Exact offline
Go1.25.13 test/race/vet/build passed. Eight controlled mutants failed intended
assertions and original bytes were restored. Prior source prefixes, unchanged
module/model tests and historical provenance are preserved.
IN_PROGRESS — TRENDING_GOOD: fresh independent review instance1/3 of the explicit
working-tree delta and canonical reconciliation. No source changes during review.

## Accepted result and integration destinations

- PASSED: fresh independent review instance 1/3 **Ready**, no actionable findings.
  Full report and preliminary ledger retained. No correction loop required.
- PASSED: archive commit `7aac77de36e33059ef0759b44ca02ab5ad080343` pushed to existing
  draft PR38. Reviewed Go/test bytes unchanged after freeze; only acceptance
  metadata appended. Prior APIs/source/provenance and false run/admission flags preserved.
- PASSED: all internal handoffs collected. Canonical branch codex/c5b19-event-fields
  receives roadmap, matrix, exact archive pin and evidence in one new PR; final URL,
  head, formatting checks and hosted acceptance are recorded in that PR body.
- Deferred: optional harness-identity and malformed-reply diagnostic mappings;
  native syscall/terminal-status semantics; raw sidecars, pairing, causal/timing
  composition, build/host closure and all execution/admission gates.
- Next bounded scope after acceptance: captured raw-frame sidecar correspondence,
  before causal/timing composition. Do not infer run authorization or full Gate3
  completion from this local field/lane layer.
