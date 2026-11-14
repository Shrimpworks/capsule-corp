# G3a4 independent implementation and plan review

Review instance: 1 of 3. Verdict: **Ready** for the exact bounded working-tree candidate.

Scoped review status: **PASSED**. Bounded implementation/local validation: **PASSED**. Publication and immutable source-pin closure remain the parent delivery step, **IN_PROGRESS — TRENDING_GOOD**. Full Gate 3 and parent owner-only internal alpha remain **IN_PROGRESS — TRENDING_GOOD**. Gates 4–6, child execution and product/runtime admission remain **BLOCKED**. This verdict does not reopen or reset earlier G3a/G3a2/G3a3 review budgets.

## Findings

No actionable findings in the reviewed delta. No source fixes or plan pivot are required.

The review found no path from envelope acceptance to write-history settlement or lifecycle authority. The 256-byte decoder checks length before slicing, exact header/code ranges, zero reserved tail, valid frozen bindings, SHA-256 of precisely the 240 binding bytes, and matching nonzero timebase before returning owned fixed arrays. The lane checker refuses malformed sizes and excess counts before iteration, checks sequence arithmetic before adding indices, and applies only a private control tick watermark. These observations support this narrow conformance prerequisite, not complete event validity.

## Review identity and method

Internal, read-only independent review in a fresh agent with no implementation conversation inherited. Used the independent-review and codebase-memory skills. Started with `/tmp/c5b19-g3a4-review-bootstrap.md`; read scope/acceptance, normative packet, tests and implementation before author explanation. Saved `/tmp/c5b19-g3a4-review-1-preliminary.md` before reading `/tmp/c5b19-g3a4-implementation.md`, archive outcome additions, or canonical outcome records. Scope/plan progress summaries were necessarily visible but not treated as verification evidence.

Canonical root: `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-event-conformance`, HEAD/base `2aefb15ee45018bb1045e711c77a0c7bc9f4c54d`. Reviewed named active Go pin/docs, roadmap changes, and `.planning/c5b19-events/` evidence. Excluded pre-existing `.serena/project.yml` tool migration.

Archive root: `/tmp/capsule-c5b19-recovery`, branch `codex/c5b19-benign-teardown`, HEAD/base `6575898a32ff4a3d24ca4c757f42b1807798a372`. Exact four-file diff under `experiments/typed-guest-transport-c5b19-benign-teardown/`: `tests/oracle.go`, `tests/oracle_test.go`, `PLAN.md`, `inputs/ORIGINS.json`. Archive status remained exactly those four files. Old implementation/test bytes remain exact prefixes; historical ORIGINS fields are unchanged.

Verified frozen SHA-256:

| Path | SHA-256 |
| --- | --- |
| tests/oracle.go | 4d7010b67089bd8d25e3c6f35101e9652c66e7147bcb5210e2103b3314f34f4b |
| tests/oracle_test.go | d134a318eb49335c3113e35fb42928ebc427022612413b7925c5fbe514cdc234 |
| PLAN.md | 5ff9ecd39a49909408059ad8f813cf93e84611c8e4b0c9d39ad2a647d0611b32 |
| inputs/ORIGINS.json | 838f9628eb886a914eaa84143d7bb2fc54d3278ac6d6c9c5038d382662447845 |

Defensive question: does the owned archive's pure parser reject malformed/substituted synthetic envelopes and malformed finite lanes within the accepted boundary? Exact targets were inert in-memory byte fixtures under ordinary Go unit-test/compiler tools. No experiment child, guest, lifecycle target, signal/wait operation, live store, real clock, network target, identity, credential or user content was exercised by this review. No source edits, staging, commits, pushes, PR changes, or subagents were created. Only the requested two report files were written; ordinary compilation used the prescribed temporary Go cache.

## Plan Review

The bounded acceptance is coherent with the staged full Gate-3 plan. Canonical `docs/C5B19_GATE2_IMPLEMENTATION_PACKET.md:309` fixes slot layout/caps; its kind/common-field and lane-isolation sections supply the exact common syntax. `.planning/c5b19-events/scope.md` expressly scopes out per-kind required/forbidden fields, lane ownership, captured-frame correspondence, causality, timing, complete capture, and live evidence. The implementation names and comments preserve that distinction (`tests/oracle.go:389`, `decodeEventEnvelope`, `validateLaneEnvelope`). A deliberately semantically invalid kind/lane combination in the known-answer test makes the limit observable instead of claiming full-case success.

Explicit trusted `firstSequence` is reasonable for this prerequisite: the packet does not select zero versus one. No untrusted first record establishes its own expected origin. Overflow checks cover the full finite range before addition; accepting one record at MaxUint64 and rejecting subsequent wrap matches the plan. Control-only tick ordering and acceptance of independent worker/harness observations match the agreed boundary; neither supplies a cross-lane temporal order or timing validity.

Tests precede implementation inspection in this review and cover independent literal bytes, every decoded field, signed/high unsigned values, invalid common fields/bindings, every short frame length, reserved bytes, closed-code combinations, all three exact capacities and cap+1, per-position lane/sequence/corruption failures, swaps, overflow, equal/regressed ticks, and invocation independence. The fixed value-type return also makes aliasing impossible by construction. Selected guard mutants supply relevant sensitivity evidence, not exhaustive proof.

ADR-0047, the C5b18 passive specification and architecture/threat-model role separation remain intact. No privileged helper or Supervisor responsibility is introduced. The ecosystem BUILD-NARROWLY/standard-library choice is appropriate for fixed bytes; PLAN records trust class, bounds, authority limits, restoration checks and maintainer ownership. No new dependency or source-manifest entry exists.

The canonical Go1.26.9 to Go1.26.11 prerequisite is a five-file active-pin/documentation patch. The unchanged CI workflow reads go.mod; unchanged Makefile retains scanner v1.6.0 and scoped lint policy. Historical and frozen archive Go1.25.13 records remain untouched. Existing scanner failure and final full checks justify the bounded patch without a literal-mirroring test.

Publication ordering is sound: review frozen bytes, then parent publishes archive source, records immutable commit/review links canonically, and verifies PR bodies/state/checks. Current pending pins and hosted acceptance are explicitly future delivery work; this source verdict does not claim their completion.

## Author-Claim Reconciliation

| Author claim | Evidence inspected | Status | Review consequence |
| --- | --- | --- | --- |
| Exact envelope decoding, all fields copied, zero-result refusal | Gate2 layout; oracle.go decodeEventEnvelope; known-answer/malformed tests; independent test run | Confirmed | Meets scoped parser contract |
| Bindings-only hash, existing frozen-binding validation | encodeObligationPreimage at oracle.go:78, preimage[30:270] reuse; literal digest and independent raw request binding tests | Confirmed | No full-obligation substitution accepted |
| Finite lane caps, explicit origin, wrap refusal, private watermark | validateLaneEnvelope; cap/sequence/tick tests; independent test/race runs | Confirmed | Meets finite framing scope; no capture completeness claim |
| Existing history unchanged; no new effects/dependency | Exact prefix comparison of source/tests; imports; four-file status; unchanged module/model hashes | Confirmed | No authority or integration expansion |
| Compiler-red before implementation | Hashed red log with missing symbols; PLAN chronology | Confirmed as retained evidence | Review did not witness historical ordering; log supports stated TDD sequence |
| Six compiled mutants fail intended assertions; exact restoration | Six raw logs, mutations.json, ORIGINS log hashes, current source hash, independent clean suite | Confirmed as retained sensitivity evidence | Did not replay mutations or independently witness every restoration; current restoration verified |
| Exact offline archive test/race/vet/build pass | Independent reviewer commands below plus matching worker logs | Confirmed | Local conformance accepted |
| Canonical scanner initially rejects Go1.26.9; final1.26.11 passes | Seven verified command-log hashes, baseline/final scanner excerpts, final test logs, five active file hashes, exact local binary version/hash | Confirmed from retained local evidence | Patch prerequisite supported; current advisory DB remains dated evidence |
| Canonical full make ci and pnpm test pass | Verified final logs, Makefile/workflow mapping; 177 workspace and132 script tests | Confirmed from retained local evidence | Full canonical suite was not rerun by reviewer |
| Existing revive debt and dependency advisories remain | Final full-lint log capped at50 revive findings; final CI audit2low/2moderate | Confirmed | Neither is concealed as a clean global backlog/audit |
| Historical hosted checks/merge and later hosted acceptance | Canonical statements/local Git baseline; no remote API consulted | Unverified independently | Outside this local source review; parent owns hosted readback and final delivery |
| No full semantics/live execution/product admission | Source comments, explicit scope and roadmap statuses | Confirmed | Full Gate3 stays unfinished and execution gates stay closed |

## Verification Performed

Archive working directory:
`/tmp/capsule-c5b19-recovery/experiments/typed-guest-transport-c5b19-benign-teardown`.

Every reviewer Go command used the exact prefix:

```sh
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/c5b19-g3a4-review-cache /Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go
```

| Suffix/command | Reviewer result |
| --- | --- |
| test -count=1 ./... | exit0, ok capsule.experiments/c5b19/tests 0.317s |
| test -race -count=1 ./... | exit0, ok capsule.experiments/c5b19/tests 1.770s |
| vet ./... | exit0, no diagnostics |
| build ./... | exit0, no diagnostics |
| git diff --check (both roots) | exit0 |

Read-only Python/hash checks verified all four archive implementation hashes, all12 archive log hashes, all3 canonical reference hashes against the exact Git commit, all7 canonical command-log hashes, and all5 active pin/doc hashes. Compared historical ORIGINS values to HEAD and verified prior oracle source/tests are exact prefixes. The retained author explanation equals the /tmp handoff byte-for-byte.

Read actual red and mutation logs, checking each expected assertion and absence of compiler-failure output in mutants. Inspected baseline/final scanner text and final test/audit/lint summaries. Canonical cached binary independently prints `go version go1.26.11 darwin/arm64`; bin/go SHA-256 `20e07ffb4f08faaf125e413db024b6e988d3ec2799dcb3568a720d14be0ad133` and LICENSE SHA-256 `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad` match retained provenance.

Codebase-memory list_projects identifies the primary capsule-corp graph at `/Users/dsteele/repos/capsule-corp`, generation2026-10-04T22:18:10Z. Task-directed symbol search returned no archive symbols, so no fabricated trace/snippet result was used. The archive is unindexed; exact source/diff reads supply its evidence. Coverage checks covered consulted canonical paths: most tracked docs/configs returned best-effort metadata_match; recovery/planning records returned missing. All exact target-worktree files were directly read. The initially guessed ADR filename was corrected from the actual listing and coverage checked on `docs/adr/0047-prepare-teardown-obligation-before-launch.md`. No absence/exhaustive claim is based on the graph.

One repository search hook temporarily rejected a combined text-search command; it reset and the read-only reads succeeded afterward. This did not affect evidence or require broader access.

## Open Questions And Residual Risks

No unresolved question blocks this bounded source candidate. Future implementation must choose the writer origin, enforce all29-kind payload/lane rules, bind sidecars and received/applied events, compose causal/timing checks, and prove real custody/absence under separately authorized gates. Envelope success alone remains insufficient.

Reviewer reran archive checks; canonical full-suite and mutation chronology rely on inspected retained worker evidence rather than repeated execution. Raw /tmp logs are ephemeral; retained PLAN excerpts/hashes and canonical validation records are the durable evidence intended for publication. Node22.23.3 differs from documented22.22.1; full lint debt #217 and2low/2moderate dependency advisories remain disclosed. No hosted checks, platform behavior, native layout, process/store synchronization, or live clock guarantees were independently verified.

## Verdict

**Ready** for the exact G3a4 bounded envelope/lane-framing implementation and canonical Go1.26.11 prerequisite. The implementation and its plan are supported by the inspected source, tests, independent offline checks and reconciled evidence. No actionable findings remain and no heavy pivot is needed.

## Recommended Next Actions

Parent should retain this report and preliminary ledger, reconcile the no-findings verdict, publish the frozen archive bytes through existing draft PR38, then retain exact immutable archive/review references and current completion status in the canonical PR. Verify PR bodies/newlines, repository/base/head/draft identities and hosted checks before delivery completion. Preserve unrelated tool changes and all remaining execution/admission gates. These are already planned delivery steps, not new review findings. No further review instance was started.
