# G3a12 independent review

Review instance: 1 of 3. Internal read-only G3a12-REVIEW.

## Findings

No actionable findings in the frozen implementation, tests or accepted bounded plan. **PASSED — review completed; verdict Ready for this G3a12 slice.** Parent full Gate3 remains **IN_PROGRESS — TRENDING_GOOD**. Child execution, Gates4–6 and product admission remain **BLOCKED**.

The implementation directly enforces the accepted partial-selection contract. It does not establish complete transport or trace validity. No heavy pivot or additional review instance is recommended for these unchanged source/test bytes.

## Review identity and independence

Canonical worktree: `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-transport-lanes`, HEAD/base `f1ded569af4a6ae20f2a93b7e476825959f97143`.
Archive worktree: `/tmp/capsule-c5b19-recovery`, branch `codex/c5b19-benign-teardown`, HEAD/base `7ec433aa17c623f58a81118bab87f1ada49770f3`.

Reviewed the explicit uncommitted boundary: archive `experiments/typed-guest-transport-c5b19-benign-teardown/{tests/oracle.go,tests/oracle_test.go,PLAN.md,inputs/ORIGINS.json}`, canonical three named status/experiment documents, and `.planning/c5b19-transport-lanes/`. Unrelated `.serena/project.yml` remains excluded and preserved. No new commit, branch, push or PR was created by this reviewer.

Started from `/tmp/c5b19-g3a12-review-bootstrap.md` and read the independent-review skill in reviewer mode. Recorded `/tmp/c5b19-g3a12-review-1-preliminary.md` after inspecting requirements, scope, implementation and tests, before reading the detailed implementation handoff, implementation evidence and canonical verification handoff. The accepted plan itself included a brief verification summary, so the preliminary pass was independent of detailed author rationale/evidence, not unaware that checks were reported. No implementation conversation was supplied.

Defensive authorized scope: owned inert synthetic event/frame fixtures, current source/document reads, and ordinary offline Go1.25.13 compiler/unit processes. No mutation execution, native fixture, child/guest, lifecycle/signal/wait target, live store/clock/network, credentials or user data. No repository edits or subagent dispatch.

Serena initialized. Codebase-memory project discovery and coverage checks identified the old primary index (2026-10-04), missing planning freshness and outside-project current worktree/archive paths. Current exact source reads supplied the evidence. No graph freshness, exhaustive call-graph or repository-wide correctness claim is made.

## Plan Review

The accepted `scope.md` contract is implemented by archive `tests/oracle.go:1207` (`validateTransportLanes`):

- Context/timebase and the 18-event cap precede processing and therefore apply to empty input.
- Each event passes existing `decodeEventFields` (`oracle.go:522`) under the same bindings/timebase; kind closure precedes indexing the nine-element counter array.
- Bounded decoded value storage, three-per-kind and three-distinct-operation caps, and exact kind/operation duplicate refusal match the plan. Generation/digest changes cannot hide a duplicate because those fields are absent from the identity comparison.
- Same-lane sequence equality and both directions of inconsistent tick order are checked against every prior event using unsigned comparisons. No subtraction, sorting, global watermark, cross-lane comparison or numeric operation/generation ordering is introduced.
- No correspondence/history validator is implicitly invoked. Returned information is error-only, and the new function retains no input bytes or state. Prior code/test prefixes are byte-exact.

The ten appended test groups at `tests/oracle_test.go:3452` onward exercise empty/context, each supported kind, malformed/unsupported events, valid local-pair counterexamples, both lanes and sequence directions, cross-kind pair boundaries and interleaving, positive equal ticks/gaps/cross-lane independence, unsigned extrema, caps/dedupe, partial-layer separation and caller preservation. The 18-event deterministic permutation and signed diagnostic values are positive acceptance tests. Existing pair validators are explicitly called before the new guard for inter-pair counterexamples, preventing malformed setup from masquerading as a lane defect.

Packet lines300–335 and395–435 preserve independent writer lanes and reject total cross-lane timestamp ordering. Its eventual complete-stream sequential-index requirement is not weakened by this separate partial selector. The partial gap, lone-event and inconsistent-tuple positives are therefore deliberate, not missing correspondence validation. ADR0047 and the test-only ADOPT-PLATFORM / BUILD-NARROWLY row remain unchanged; no product authority or runnable mechanism is added.

The plan sequences scope acceptance, tests-first evidence, implementation, mutation sensitivity, offline checks and fresh review appropriately. Canonical baseline checks and known limitations are retained separately. Immutable archive pinning, review acceptance metadata, canonical publication and exact-head hosted checks remain Brain-owned delivery work; they are pending rather than falsely complete. The next-slice control ordering recommendation is explicitly deferred and requires its own producer-edge scope freeze. It is not accepted as implemented by this review.

## Author-Claim Reconciliation

| Material claim | Evidence inspected | Status and consequence |
| --- | --- | --- |
| One stateless bounded partial lane validator; prior APIs unchanged | Entire appended function, relevant strict decoder and canonical passive model/types; source prefix comparison | Confirmed. Correct separation of partial consistency from effects and authority. |
| Local pairs can pass while a combined lane selection contradicts | `TestTransportLanesLocallyValidWritePairs`, `LocallyValidReceivePairs`, `CrossKindPairBoundary`; independent full test run | Confirmed. New guard addresses an actual finite gap without requiring complete pairs. |
| Input-order independence, unsigned boundaries, independent lanes, equal ticks and gaps | Direct pairwise logic plus boundary/permutation tests | Confirmed. No overflow arithmetic or global clock state. |
| Ten new test groups, 53 appended implementation lines, prior source/test prefixes preserved | Diff and current source; byte-prefix checks | Confirmed. No unrelated implementation changes. |
| Tests-first red/green | Hash-verified retained red/green logs: undefined validator then successful package test | Supported by inspected historical evidence. Reviewer did not recreate or independently observe authoring chronology. |
| Twelve compiled mutants reach seventeen intended assertions | All twelve raw logs and exact replacement/test definitions, assertion counts and hashes; mutation driver read only | Confirmed as retained evidence. Every named assertion was present; no build failure/panic/setup refusal was substituted for a kill. Some tests additionally hit outer boundary assertions; seventeen is the named intended count, not all failure lines. |
| Exact restoration and final offline checks | Current source hashes, retained restoration/final logs and driver finally-restoration logic; reviewer full/race/vet/build | Confirmed current restored bytes and final outcomes. Per-attempt restoration is historical driver/log evidence, not a reviewer-executed campaign. |
| Historical provenance and non-admission flags preserved | Parsed old ORIGINS versus current; only new top-level `g3a12` key; unchanged go.mod/model_test.go hashes | Confirmed. Existing historical values and false gates were not rewritten. |
| Canonical required baseline passed, with disclosed lint/audit limits | Five baseline log hashes, raw summaries, validation.json and handoff | Confirmed as inspected baseline evidence. Unfiltered revive lint still exits1; two low/two moderate dependency findings remain. No clean-all-lint or zero-dependency-findings claim. |
| Full trace, causal control ordering, timing/native execution and admission remain incomplete | Source comments, archive plan, canonical diffs, accepted plan and next-scope | Confirmed. No completion or execution-authority inflation. |

## Verification Performed

Reviewer-executed commands ran in `/tmp/capsule-c5b19-recovery/experiments/typed-guest-transport-c5b19-benign-teardown` using exact executable `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go`, with `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/c5b19-review1-go-cache`.

| Command | Observed outcome | Reviewer log / SHA256 |
| --- | --- | --- |
| `go version` | go1.25.13 darwin/arm64 | tool output |
| `go test -count=1 ./...` | exit0, package passed | `/tmp/c5b19-review1-test.log`; `4da1980650b2b132618fa6bcd2fee57ccf7d314f48579f7549986b9c6bfd994d` |
| `go test -race -count=1 ./...` | exit0, package passed | `/tmp/c5b19-review1-race.log`; `ef6263f68b76368c62b9de20054087360877763e8a3dabe62f15f5a5ffb0d5c3` |
| `go vet ./...` | exit0 | `/tmp/c5b19-review1-vet.log`; `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `go build ./...` | exit0 | `/tmp/c5b19-review1-build.log`; `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `git diff --check` in both worktrees | exit0 each | tool output |

Read-only Python verification independently checked the frozen patch and all four frozen file hashes before and after review, canonical reference hashes, implementation hashes, prior source/test prefixes, historical parsed ORIGINS preservation, matching retained g3a12 evidence, all eleven historical verification log hashes, all twelve mutant log hashes and exact intended assertion occurrences, mutation driver hash, and all five canonical baseline log hashes. All passed. Current imports/package inventory and focused source inspection support the ordinary offline unit-process scope.

Frozen archive patch SHA256: `b1523a16e6d47b8538218a966604cdc6849167c0da8241a59301cfc52cd8c2fb`.

| Frozen file (archive experiment relative) | SHA256 |
| --- | --- |
| tests/oracle.go | `22e91fd0c95e778f9e8e68480bcb586fb4fe64d4d3c0ebb9b0c56b887b7ac0be` |
| tests/oracle_test.go | `8c10cdab499b5e0e9260848548adae83fc788b9f1c26defae054dd06c9cf5b2c` |
| PLAN.md | `5d1b65fbd1c292ea0fe74782670ed8993e334db9ec2b020702471b735d6dc645` |
| inputs/ORIGINS.json | `8192daa4bb347b469a9f070dc5a8ebcc2c4581a658e71dc0dc71d64e373971a1` |

Historical mutation assertion counts independently read from logs: supported-kind1, duplicate-kind/operation1, same-lane-sequence4, reverse-tick2, forward-tick2, cross-lane-filter1, strict-tick1, unique-operation1, adjacency1, correspondence1, empty-timebase1, caller-reordering1; total17. No mutations were executed by this reviewer.

Canonical `pnpm install --frozen-lockfile`, `make ci`, `pnpm test`, full lint and baseline diff checks were inspected, not rerun. Their Go1.26.11/Node22.23.3/pnpm10.28.2 environment is distinct from the frozen archive Go1.25.13. Full lint debt #217 and the four below-threshold audit findings are pre-existing disclosed limitations, not new G3a12 findings. Hosted results were not queried during this offline review.

## Open Questions And Residual Risks

No unresolved implementation question blocks this bounded slice. Finite tests and twelve selected mutations are not exhaustive proof. Total/per-kind caps overlap with the unique-operation/duplicate limits; no independent redundant-cap necessity is claimed. The supported-kind mutant safely removes only the lower bound. Caller concurrent mutation is outside this sequential fixture API contract; race-pass does not establish that contract for shared concurrently modified inputs.

Successful selection can still omit events, contain lone result/application events, or contain different strict tuple fields across kinds for one operation. It proves neither capture truth/completeness nor raw-frame identity, semantic role order, transport success, decision settlement, timing, process absence or authority. Those limitations are intentional and explicitly documented. Child execution remains BLOCKED.

Temporary raw logs are ephemeral. Brain owns durable review/evidence retention, acceptance-only metadata changes, immutable archive commit pin and publication. This verdict covers the source/test freeze above; later behavior changes require re-evaluation within the existing three-instance limit.

## Verdict

**Ready** for the frozen bounded G3a12 supplied transport lane consistency implementation and plan. **PASSED** for this review. Parent Gate3 **IN_PROGRESS — TRENDING_GOOD**; child execution, Gates4–6 and product admission **BLOCKED**. This is not Gate4/full implementation acceptance or first-run authorization.

## Recommended Next Actions

Brain should retain this report and preliminary ledger, record acceptance, publish the unchanged reviewed source/test bytes through existing draft archive PR38, pin the resulting archive commit in canonical documentation, complete final metadata checks, and publish the canonical PR with verified body/base/head/draft state and exact-head hosted checks. Preserve the disclosed lint/audit limits and blocked execution gates. No fixes, new workstreams or heavy pivot are requested.
