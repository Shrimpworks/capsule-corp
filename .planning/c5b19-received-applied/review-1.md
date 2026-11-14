# G3a7 independent received/applied correspondence review

Review instance: 1 of 3.

## Findings

No actionable findings. No correctness, security-boundary, compatibility, or accepted-scope defect was identified in the frozen implementation or tests.

Review status: **PASSED** for this bounded received/applied correspondence slice. Parent Gate 3 remains **IN_PROGRESS — TRENDING_GOOD**; child execution, Gates 4–6 and product admission remain **BLOCKED**. This report does not certify those parent gates.

## Reviewed target and independence

Internal read-only fresh-context review using independent-review reviewer mode. Canonical root `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-received-applied`, HEAD/base `041b6d4e04df9667e6217f57df26f04eec7d2ddc`. Archive root `/tmp/capsule-c5b19-recovery`, branch `codex/c5b19-benign-teardown`, HEAD/base `14365fb4ba5d1525285406ae90feb1abdcc46269`. The target is the explicit uncommitted diff, not new commits.

Archive scope is `experiments/typed-guest-transport-c5b19-benign-teardown/{tests/oracle.go,tests/oracle_test.go,PLAN.md,inputs/ORIGINS.json}`. Canonical scope is three documentation changes (`docs/CURRENT_WORK_PLAN.md`, `docs/PROJECT_RECOVERY_2026_10_15.md`, `docs/C5B19_BENIGN_TEARDOWN_EXPERIMENT_PLAN.md`) plus the five files in `.planning/c5b19-received-applied/`. Unrelated `.serena/project.yml` was excluded and preserved. Start/end Git state matched this boundary.

Requirements, accepted scope, tests and source were reviewed before the author handoff and new archive outcome additions. Preliminary findings were written to `/tmp/c5b19-g3a7-review-1-preliminary.md` first. Blinding was imperfect because the bootstrap/accepted plan already disclosed author progress summaries; no author verdict was adopted. The separate author handoff was subsequently checked against source and logs.

Defensive question: do supplied owned inert kind7 receive bytes and exact request/reply sidecars correspond to a kind8 decision record, with consistent selected control sequence/tick observations? Authorized environment was the two named local checkouts, synthetic in-memory fixtures and ordinary offline compiler/unit-test processes. No native fixture, child/guest, lifecycle target, signal/wait, live clock/store/network, credential or user-data operation occurred.

## Plan Review

The implementation satisfies the accepted minimum scope without extending Gate2's frozen format or authority model:

- `tests/oracle.go:720` validates received sidecars through the unchanged strict validator, decodes both events under the same frozen binding/timebase context, requires kinds7/8, compares Operation, Generation, RequestDigest, CandidateDigest, Outcome and RecordDigest, and requires strictly increasing pair sequence plus nondecreasing tick.
- `tests/oracle.go:745` checks context even for an empty selection, caps input at three pairs before array writes, refuses repeated operations across generations, and compares all at-most-six copied control records for unique sequences and nondecreasing ticks by sequence. Both comparison directions support arbitrary pair-list order. Direct comparisons avoid uint64 increment/elapsed-time overflow.
- Lower `decodeEventFields` and `validateCapturedFrame` contracts supply lane ownership, forbidden-field zeros, decision1/2, exact raw-frame sizes/bindings and request/candidate/record hash checks. The new functions are error-only and cannot mutate history or grant authority. Repeated decoding is bounded and consistent with the existing wrapper style.
- Seven added tests (`tests/oracle_test.go:1520` onward) exercise operation shapes1/2/3, all outcomes1/2/3 with both decisions, independent tuple mismatch offsets, wrong kinds/lanes/context, omitted/short/extra members, raw-frame substitution, same/later ticks, gaps/zero origins/uint64 boundaries, collection limits and operation reuse, all endpoint collision roles, both pair-list orders, cross-pair tick regression and nonmutation/prior API independence.
- The state model supports the scope's distinction between accepted failed/indeterminate settlement and confirmed durable projection: `internal/execution/teardownpassive/model.go:185` checks outcome and exact pending token before settlement, then returns without projection for non-confirmed outcomes.

No behavior-scope deviation was found. The plan deliberately defers decision correctness, rejected-frame fatal latching, full capture/causality/timing, pending/control-state composition and durable snapshots; those deferrals are consistently stated in source and changed docs. This is an appropriate incremental prerequisite, not completion of the whole packet oracle. No new dependency, format, native component, privileged responsibility or product import is introduced.

Final acceptance metadata, immutable archive pin, retained evidence formatting/site checks, pushes and PR publication are explicitly Brain-owned and were pending at freeze. A Ready implementation verdict does not turn them into completed delivery.

## Author-Claim Reconciliation

| Author claim | Evidence inspected | Status | Review consequence |
| --- | --- | --- | --- |
| Exact kind7/kind8 tuple and sidecar correspondence | New validators and unchanged lower decoders; new tests | Confirmed | Matches accepted local contract |
| Three-pair cap; distinct operations/sequences; monotonic ticks independent of list order | Fixed six-event array, bounded nested comparisons, collision/order fixtures | Confirmed | No unchecked growth or sorting/mutation |
| Both decisions/all outcomes representable, with no settlement or authority | Positive matrix, error-only call path, normative SettleWrite | Confirmed | Claimed security boundary remains narrow |
| Earlier APIs/source prefixes and provenance preserved | HEAD-prefix byte checks; exact module/model equality; parsed ORIGINS with only g3a7 removed | Confirmed | Additive change preserves accepted prior slices |
| Tests-first compiler-red then green/restored checks | Hash-verified red/green/restored logs; independent current checks | Confirmed as supplied historical evidence plus independent current success | Red is missing newly introduced APIs, not unrelated failure; historical timing itself not independently reproduced |
| Fourteen compiled mutants fail intended assertions | Exact replacement JSON and all fourteen hash-verified failure logs | Confirmed by inspection, not rerun | Selected guard sensitivity only; excluded cross-tick construction log is a vet/build failure and is not counted |
| Canonical baseline checks pass, with disclosed lint/audit limits | validation.json, canonical-verification.md, all four hash-verified logs and Makefile | Confirmed by inspection, not rerun | Existing unfiltered revive debt remains BLOCKED; no clean full-lint or vulnerability-free dependency claim |
| Exact source/reference/log hashes | Independent SHA-256 computations | Confirmed | Frozen target and inspected evidence agree |
| Prior hosted PR/CI outcomes and no protection bypass | Local plan/docs testimony only; no network check permitted | Unverified independently in this review | Brain retains hosted integration-readback responsibility; not needed to establish local validator correctness |

## Verification Performed

Independent commands used cwd `/tmp/capsule-c5b19-recovery/experiments/typed-guest-transport-c5b19-benign-teardown`, binary `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go`, and exact environment `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/c5b19-g3a7-review-cache`:

| Command | Observed result |
| --- | --- |
| `go version` using the exact binary | `go version go1.25.13 darwin/arm64` |
| `test -count=1 ./...` | exit0, `ok capsule.experiments/c5b19/tests 0.359s` |
| `test -race -count=1 ./...` | exit0, `ok capsule.experiments/c5b19/tests 1.796s` |
| `vet ./...` | exit0, no diagnostics |
| `build ./...` | exit0, no diagnostics |
| `git diff --check` in archive and canonical checkouts | both exit0 |

Independent Python read-only comparisons confirmed:

- Previous HEAD bytes are exact prefixes of oracle.go (26,635 prior bytes), oracle_test.go (57,719) and PLAN.md (52,908); go.mod (45) and model_test.go (14,636) are byte-identical to HEAD.
- Historical parsed ORIGINS values are identical after removing the new g3a7 object, preserving prior run/admission flags.
- Four module/source hashes, three normative canonical hashes, all 20 archive log hashes, the separately excluded mutation-construction log, and four canonical log hashes match the supplied evidence.
- `/tmp/c5b19-g3a7-implementation.md` exactly matches retained `implementation-handoff.md`.
- Normative packet/types/model bytes in the reviewed canonical worktree equal those in the primary checkout behind the graph.

Frozen archive SHA-256 values independently confirmed again after checks:

| File | SHA-256 |
| --- | --- |
| tests/oracle.go | `b9cf1d6c311da52062eb86aec40f2e1cce1f03767bfbb5f6686a23d4746672bb` |
| tests/oracle_test.go | `2ae581d0ff919974d73319fd953580f878435ae4693a65186e276a4d081fabf6` |
| PLAN.md | `5c187764aedf50f1bc6a65f961b2f4d7d0ed6b43a6bf05d2bf3fbac7f74861c7` |
| inputs/ORIGINS.json | `c161dd22e021cc2f11467a60fa171d1317029982ed1bf714b4cf9965a18da5ab` |

Canonical baseline logs were inspected rather than rerun. They support frozen install, make ci, standalone pnpm test, 177 workspace plus132 script tests, and govulncheck with no vulnerabilities reported. Full unfiltered lint exits1 with50 displayed revive findings under its output cap, not a total debt count. Dependency audit reports two low/two moderate findings below the unchanged high threshold. Node22.23.3 differs from the documented22.22.1; archive Go1.25.13 is distinct from canonical Go1.26.11. These disclosed limitations are not changed by this review.

Retrieval used codebase-memory Tier2: primary index ready at2026-10-04T22:18:10Z, not watched. SettleWrite search returned one symbol; outbound trace had four callees and no continuation; exact snippet was read. Coverage checked21 canonical evidence paths across two calls: normative files and stable docs report best-effort metadata_match, new planning files and recovery doc freshness missing. Current direct reads/diffs supplied those gaps; source equality bound the normative snippet to the reviewed worktree. Archive is unindexed and was reviewed directly. No graph completeness claim is made.

## Open Questions And Residual Risks

No unresolved question blocks this bounded implementation. Future work must prove that represented decisions match pending/control state, rejected frames latch fatal, stop/apply ordering is correct, expected projections match observations, and worker/snapshot/full trace evidence agrees. Whole-pair omission and capture origin/completeness cannot be detected by this API. Malformed rejected-frame diagnostic mapping is unsupported pending a frozen contract. Success does not prove native syscall facts, clock/deadline guarantees, durable publication, lifecycle authority or admission.

Mutation evidence is selected and inspection-based; no mutations were reapplied because this review is read-only. Raw logs in /tmp are ephemeral and are not release-grade provenance. Brain must retain appropriate evidence and exact commit links during integration. Existing canonical lint debt, dependency audit findings and toolchain patch discrepancy remain as disclosed baseline limitations.

## Verdict

**Ready** for the exact frozen G3a7 received/applied correspondence implementation and bounded plan. Confidence is high within the selected inert in-memory contract; no broader security-control or runnable-mechanism claim is supported.

## Recommended Next Actions

Brain should retain this report, record acceptance metadata, pin the reviewed archive source to its eventual immutable commit, complete final retained-document/JSON/site checks, and publish through the already named archive/canonical PR destinations with required readbacks. Keep full Gate3 **IN_PROGRESS — TRENDING_GOOD** and child execution/admission **BLOCKED**. No additional review instance or workstream was created.

This reviewer wrote only the two authorized /tmp reports plus ordinary compiler cache artifacts. No source fix, staging, branch, commit, push, PR edit, subagent or native execution occurred. All owned command sessions completed. Handoff returns to Brain for acceptance and delivery.
