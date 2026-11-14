# C5b19 G3a6 independent review

Review instance: 1 of 3.
Review work status: PASSED. Verdict: Ready for the frozen bounded implementation and accepted plan. Parent Gate 3 remains IN_PROGRESS — TRENDING_GOOD; child execution, later gates and product admission remain BLOCKED. Delivery/publication remains Brain-owned and incomplete.

## Findings

No actionable findings. No implementation or plan correction is required by this review.

The archive pair validator at `experiments/typed-guest-transport-c5b19-benign-teardown/tests/oracle.go:634` composes strict event/request/reply decoders, restricts kinds to 3/5/7, compares operation/generation, recomputes the request-body and candidate hashes, and compares validated reply outcome/record digest. Kind3 refuses any reply. Existing exact-length checks and error short-circuiting protect every new byte slice. Confirmed record hashes remain recomputed by the existing 395-byte record codec; failed/indeterminate replies require zero record digest.

The collection validator at `tests/oracle.go:666` refuses invalid context even for an empty selection, more than nine entries, more than three entries per category, duplicate same-category operations, more than three distinct operation anchors, inconsistent same-operation request bytes, and inconsistent worker/control peer replies. Work is bounded by nine inputs with fixed metadata arrays. Success does not return a capability or retain references; neither validator has a history, store or lifecycle effect path. The second event decode is a small bounded cost consistent with the error-only API.

Seven added tests at `tests/oracle_test.go:1246`, `:1272`, `:1296`, `:1353`, `:1393`, `:1443` and `:1467` exercise valid operation/outcome combinations, partial and capped selections, exact comparisons, malformed lengths/headers, unsupported kind8, context/lane rejection, duplicate/cap/peer rules, and nonmutation. Fixtures are constructed from literal packet offsets and prior independent vectors rather than the new validators. The unrelated-history assertion is supplementary; source inspection establishes that no history object is reachable from this API.

## Plan Review

The accepted scope and implementation conform to Gate2 packet transport and raw-evidence requirements (request397/reply98, request body excluding the eight-byte header, candidate68, record395, kind3/5/7 raw slots). The pair-struct alternative and accepted three-distinct-anchor limit are documented in `.planning/c5b19-frame-sidecars/scope.md` and `plan.md`. Missing members refuse; wholly omitted pairs remain undetectable, explicitly and consistently documented.

The small opt-in standard-library composition follows the C5b19 BUILD-NARROWLY reuse row and ADR47's separation of represented storage bytes from control settlement. `internal/execution/teardownpassive/model.go:185` applies a projection only after exact pending-token validation and confirmed outcome; the sidecar code does not claim or perform that step. No dependency, source manifest, serialized container, privileged responsibility or frozen Go1.25.13 archive pin change appears in the diff.

Plan ordering is appropriate: scope, tests first, implementation/mutations, source freeze, independent review, then Brain-owned integration and immutable publication. Canonical documentation accurately leaves full Gate3 in progress and later execution/admission blocked. Archive pin advancement, review acceptance metadata, final formatting/site checks, commits/pushes, PR readback and hosted CI are explicit later delivery gates, not silently credited as complete by this review. No pivot or additional workstream is needed.

## Author-Claim Reconciliation

| Author claim | Evidence inspected | Status and consequence |
| --- | --- | --- |
| Exact local correspondence, no settlement | oracle.go:634 and :666, lower decoders, Gate2 packet and C5b18 SettleWrite | Confirmed; no authority expansion |
| All previous APIs unchanged | Byte-prefix comparison to archive HEAD for oracle.go/oracle_test.go; exact equality for go.mod/model_test.go | Confirmed; earlier acceptance contracts preserved |
| Historical provenance and false run/admission flags preserved | Parsed ORIGINS compared to HEAD; only g3a6 key added; flags both false | Confirmed |
| Tests-first compiler red | Hashed red log reports undefined capturedFramePair/validators | Confirmed as retained evidence; chronology is author-recorded, not independently replayed |
| Eleven compiled guard mutants fail intended assertions | PLAN table, mutation JSON replacement strings, all eleven hashed logs | Confirmed as retained selected sensitivity evidence; not rerun by reviewer |
| Restored source and green test/race/vet/build | Frozen hashes before/after plus independent reruns | Confirmed |
| Canonical baseline passed except known unfiltered lint debt | Retained validation.json/handoff and four hash-verified logs | Confirmed as inspected prior evidence, not independent rerun; existing #217 debt and two low/two moderate audit findings remain |
| Whole-trace, causal/timing, delivery and lifecycle claims excluded | API comments, scope/plan, archive PLAN/ORIGINS and canonical diff | Confirmed; no overclaim found |

The cap mutation loosens total, category and unique-anchor caps together. Its failure establishes sensitivity to the combined cap rules, not independent necessity of each redundant cap. This agrees with the selected-guard evidence claim and is not an actionable defect.

## Verification Performed

Exact canonical checkout `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-frame-sidecars`, HEAD/base `e1a966e7c6f18665784984e97b6697cc472cb288`. Exact archive checkout `/tmp/capsule-c5b19-recovery`, branch `codex/c5b19-benign-teardown`, HEAD/base `7aac77de36e33059ef0759b44ca02ab5ad080343`. Reviewed the explicit working-tree diff; no implementation commits exist yet. Unrelated canonical `.serena/project.yml` remained untouched. Archive changed paths are exactly the four authorized files.

Independent commands ran from the archive experiment directory using `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go` and `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/c5b19-g3a6-review-cache`:

| Command | Observed outcome |
| --- | --- |
| `go test -count=1 ./...` | exit0, package passed 0.367s |
| `go test -race -count=1 ./...` | exit0, package passed 1.830s |
| `go vet ./...` | exit0, no diagnostics |
| `go build ./...` | exit0, no diagnostics |
| `git diff --check`, both worktrees | exit0 |

Read-only Python checks independently verified source/test/PLAN append-only preservation, exact go.mod/model_test.go equality, historical parsed ORIGINS equality, only the added g3a6 key, four implementation/module hashes, all 17 implementation log hashes, and four canonical baseline log hashes. Each mutant log contains assertion failure rather than build failure. Normative packet/types/model hashes match the handoff. All four frozen archive file hashes remained identical after checks:

- oracle.go: `5106897e28cfadd491f4080c1c8d2cd3901fc49bbafe9ce8900924e35adac46e`
- oracle_test.go: `a44445ff7cf4bbef309205971eb914fead4dcd93ac45b40d499bef66cea42ea4`
- PLAN.md: `1099bca9dd78f19cc38148de33f02de1967d0e9c3ff131a965645b88eb27a449`
- ORIGINS.json: `e5de8d7d107ef2fb86950e03b320d72289a6724bad202d5fdf56bbd031f318c3`

## Open Questions And Residual Risks

No unresolved question blocks this exact slice. Omitted pairs, capture origin, request delivery, worker receipt, kind8 application, received/applied pairing, durable snapshots, lane sequencing, native syscall semantics, causal/timing order, full-case success and actual durability/lifecycle authority remain outside it. Future composition must enforce those obligations before confirmation or execution claims.

Independence: neutral bootstrap and source/requirements were inspected before author explanation and archive outcome additions. Preliminary ledger was saved at `/tmp/c5b19-g3a6-review-1-preliminary.md` first. The accepted plan contains short implementation outcomes, so the preliminary pass was not perfectly blinded; no author verdict was adopted.

Retrieval: Tier2 primary index generation `2026-10-04T22:18:10Z`, not watched. Exact-path coverage checked; canonical packet/types/model and most normative docs reported metadata_match, planning/recovery paths missing, archive outside indexed project. Current exact source/diff reads supplied authoritative fallback. SettleWrite graph search, four-callee outbound trace and snippet were checked against the target checkout. No graph completeness claim. Relevant starter-document boundaries were read; this is a bounded diff review, not a repository-wide audit.

Authorized defensive environment: owned inert synthetic bytes and ordinary offline compiler/unit-test processes only. No native fixture, child/guest, lifecycle/signal/wait target, live clock/store/network, credential or user-data operation was performed. No source mutation, staging, commit, push, PR edit, branch creation or subagent dispatch. Only the two permitted review reports and ordinary compiler cache artifacts were written.

## Verdict

**Ready** for the frozen G3a6 captured-frame correspondence implementation and accepted bounded plan. Review work PASSED; parent Gate3 remains IN_PROGRESS — TRENDING_GOOD. This verdict neither completes publication nor authorizes child execution or product admission.

## Recommended Next Actions

Brain should retain this report and preliminary ledger, record acceptance, then complete the already-planned final metadata checks and immutable archive/canonical publication. Preserve reviewed code/test hashes and explicit limitations. Archive integration destination remains existing draft PR38; canonical destination is the Brain-owned PR from codex/c5b19-frame-sidecars. No further review instance is required absent a material implementation or scope change.
