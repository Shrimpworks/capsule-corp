# G3a11 independent review

Review instance: 1 of 3.

## Findings

No actionable findings. **Ready** for Brain acceptance and publication of the exact bounded G3a11 implementation and plan reviewed here.

Review status: **PASSED**. Parent delivery and full Gate3 remain **IN_PROGRESS — TRENDING_GOOD**; child execution, Gates4–6 and product admission remain **BLOCKED**. This verdict does not complete publication or admit any runnable control.

## Review identity and independence

Internal read-only assignment. Started from `/tmp/c5b19-g3a11-review-bootstrap.md`, applicable repository instructions and the independent-review skill in reviewer mode. Read canonical requirements, accepted scope/plan, actual diff, tests and dependency source before writing `/tmp/c5b19-g3a11-review-1-preliminary.md`. Only afterward opened the detailed implementation handoff, canonical verification, appended archive PLAN/ORIGINS rationale and implementation evidence. The accepted plan includes brief author verification summaries; the preliminary ledger explicitly recorded those as unverified testimony, not the basis for judgment.

Canonical worktree `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-write-results`, HEAD/base `d042e53a877e494bd2eb0e325ec084d4b9cc692a`. Archive `/tmp/capsule-c5b19-recovery`, branch `codex/c5b19-benign-teardown`, HEAD/base `d698bd9332a000ad3c2b9f23f87d0353dc5964b4`. Reviewed explicit uncommitted changes to four archive files, three canonical documents and untracked `.planning/c5b19-write-results/`. Unrelated `.serena/project.yml` excluded and preserved. No branch, commit, push, PR, repository edit, mutation execution or subagent dispatch occurred.

Archive prefix below is `experiments/typed-guest-transport-c5b19-benign-teardown/`:

- `tests/oracle.go`: SHA256 `4b9012aa8e19b5b35995e5fb48a7e6fd8cd2f858ac7a92ea2c4678e85ce6bc69`.
- `tests/oracle_test.go`: SHA256 `acb23007854fc92f8ddbffa7daea3c1bfeaa172ea85e88644fc7a0ef2f51975c`.
- `PLAN.md`: SHA256 `c0f185ed7548cbb6fa0aabcbb475fd7fd2652498fae4ba7ceb6350a9447b989e`.
- `inputs/ORIGINS.json`: SHA256 `086c03605ff2f240450c5bd94ea8256a59ad8316788a5e2b0ad9edca57a8d968`.
- Four-file patch: SHA256 `e8ea12e51b7185fb025842eab087f2ff9b046a774ac79734e4cccbbaa2d977ed`, matched before and after verification.

## Plan Review

The accepted `.planning/c5b19-write-results/scope.md` contract is implemented by `inspectWriteResults` at archive `tests/oracle.go:1141–1198`. It composes existing strict frame/event/request validation instead of introducing authority or state. The local kind closure precedes indexing `counts[before.Kind]`; six-pair input cap precedes fixed-array indexing. Ignored repeat decode errors are justified by successful validation of exactly those bytes earlier in the sequential call; concurrent caller mutation is outside the declared API contract.

Exact operation/generation/request/candidate identity and outcome/record equality, increasing sequence, nondecreasing tick, duplicate refusal, three unique operations and full request equality across categories match the contract. Output is fixed-size copied values, input order is preserved, and every refusal path returns the zero aggregate. The new function does not call history mutation or canonical `SettleWrite`. Nine appended test groups at `tests/oracle_test.go:3106–3406` cover the material accepted cases, with valid substitution setups that reach the new comparisons. Raw signed boundaries and category-specific 397/98 counts are directly asserted; same-tick and maximum-width ordering do not depend on subtraction.

The choice to expose only `ExactFrameCount` and raw result/errno is faithful to packet lines187–214 and300–381 and to the accepted scope. Neither strict diagnostic validity nor exact represented byte count asserts syscall realizability or transport progress. No source-supported normalized-errno rule is invented. Independent worker receipt, control application, durability and full-stream consistency remain separate obligations.

Canonical CURRENT_WORK_PLAN, recovery continuation and experiment plan consistently retain incomplete Gate3 and blocked child/later gates. ADR0047 authority and C5b19 test-only ADOPT-PLATFORM / BUILD-NARROWLY boundaries are preserved. No dependency, schema, manifest, ABI, privileged responsibility or product behavior changes. The next-scope recommendation identifies sequence reuse and cross-pair tick regression without silently implementing or approving it. Final acceptance metadata, immutable archive commit links and shared PR publication are still Brain-owned delivery work, not missing implementation requirements concealed as completion.

## Author-Claim Reconciliation

| Author claim | Evidence inspected | Status and consequence |
| --- | --- | --- |
| Pure, bounded pair correspondence with exact identity and owned atomic output | New API and existing decode/validate helpers; all nine new test groups; reviewer offline test/race/vet/build | Confirmed for supplied inert bytes and sequential use. |
| Raw signed diagnostics; exact count independent of errno | `TestWriteResultSignedDiagnostics`, signed decode fields, category lengths and two selected count mutants | Confirmed; no syscall-success inference. |
| Prior APIs, source/test prefixes and provenance preserved | Git base comparison and SHA256 checks; parsed old ORIGINS entries compared with new document | Confirmed: source, test and PLAN prefixes preserved; all18 old top-level provenance entries unchanged. |
| Tests-first compiler-red then green | Hash-verified red and green logs | Confirmed as retained historical evidence: red names missing new types/function, green passes. Reviewer did not recreate this history. |
| Eighteen compiled mutants reach27 intended assertion occurrences | All18 hash-verified raw logs, exact assertion/count checks, replacements and driver inspection | Confirmed as inspected historical evidence. No build failure/panic/malformed-setup output; no mutation run by reviewer. Scope is selected sensitivity only. |
| Byte-exact mutation restoration | Driver `finally`, restoration JSON digest, current source hash and subsequent historical restored/race/vet/build logs | Consistent and confirmed as retained evidence; current hash independently verified. |
| Canonical baseline checks pass with disclosed debt | Five hash-verified logs, validation.json and canonical-verification.md | Confirmed as recorded worker evidence. Full lint exits1 with50 displayed revive issues; audit findings and toolchain differences are disclosed. Reviewer did not rerun canonical gates or advisory/network operations. |
| Hosted merge/check metadata | Canonical plan/document references only | Not independently queried in this offline review; no implementation conclusion depends on live hosted state. Brain owns final readback. |

The retained canonical implementation-evidence JSON matches every ORIGINS.g3a11 field. The mutation-driver hash matches provenance. Cap guards overlap logically; no isolated redundant-guard necessity or exhaustive mutation claim is made.

## Verification Performed

Reviewer-executed commands used experiment cwd `/tmp/capsule-c5b19-recovery/experiments/typed-guest-transport-c5b19-benign-teardown` and exact binary `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go`, with `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/c5b19-g3a11-review-1-cache`:

| Command | Reviewer result |
| --- | --- |
| `go version` using the exact binary | go1.25.13 darwin/arm64 |
| `go test -count=1 ./...` | exit0, `ok capsule.experiments/c5b19/tests 0.310s` |
| `go test -race -count=1 ./...` | exit0, `ok capsule.experiments/c5b19/tests 1.940s` |
| `go vet ./...` | exit0 |
| `go build ./...` | exit0 |
| `git diff --check` in archive and canonical worktrees | exit0 both |
| Read-only Python hash, prefix, provenance and evidence checks | Passed; frozen patch unchanged after checks |

Reviewer race/vet/build logs: `/tmp/c5b19-g3a11-review-1-{race,vet,build}.log`. Unit result is in the tool transcript. Historical logs were read and checked, never rerun as mutations. Canonical checks are worker-executed evidence: frozen install, make ci and pnpm test exit0; unfiltered golangci-lint exit1 with disclosed baseline debt; diff check exit0. All five canonical log hashes match validation.json. All seven archive verification-log hashes and all18 mutant-log hashes match ORIGINS.

Ordinary offline compiler/unit processes operated only on owned inert synthetic fixtures. No native fixture, child/guest/lifecycle/signal/wait target, live store/clock/network or external identity/data was exercised.

## Open Questions And Residual Risks

High confidence in this finite supplied-byte correspondence contract and its bounded ownership/error behavior. No claim of capture authenticity/completeness, actual syscall result, physical transfer, independent receipt, control settlement, durability, whole-lane chronology, timing bounds, native absence, execution authority or product admission follows. A caller can supply individually valid pairs with contradictory inter-pair lane indices; this is explicitly excluded and covered by the deferred next-scope recommendation. Raw contradictory result/errno combinations remain diagnostics, not evidence of realizable native behavior.

Serena initialized. codebase-memory index_status, symbol search, outbound trace and snippet were used for normative SettleWrite orientation; exact-path coverage was checked for all relied-on code/document scopes. The graph is the old2026-10-04 primary checkout, archive paths are outside_project and new planning/recovery freshness is missing. Current target-worktree reads supplied relevant evidence. No fresh whole-graph or exhaustive structural claim is made.

Raw /tmp evidence is ephemeral; retention and immutable archive/canonical linking remain Brain's publishing responsibility. Final hosted check state and PR metadata are outside this offline review. No heavy pivot or unresolved implementation question identified.

## Verdict

**Ready** — bounded G3a11 code and plan only, review instance1 of3. Review **PASSED**; parent full Gate3 **IN_PROGRESS — TRENDING_GOOD**; child execution, Gates4–6 and admission **BLOCKED**.

## Recommended Next Actions

Brain should retain this report and preliminary ledger, record acceptance against the frozen source/test bytes, pin the archive evidence to the accepted immutable commit, and finish the existing archive draft PR38 plus canonical PR with final metadata/readback and exact-head checks. Preserve all existing execution/admission blockers. No additional review instance or implementation work is requested by this report.
