# NEXT2-ARCHIVE reconciliation — 2026-10-15

Status: `PASSED` for read-only archive-state reconciliation. C5b19 Gate-3 parent: `IN_PROGRESS — TRENDING_GOOD`; complete implementation, Gate-4 actual-byte/first-run admission, Gate-5 runnable evidence and product admission remain `BLOCKED`.

## Scope and method

Internal read-only assignment for defensive recovery of the exact owned `Shrimpworks/capsule-experiments` repository and the `experiments/typed-guest-transport-c5b19-benign-teardown` source. Used Keychain-backed `env -u GH_TOKEN -u GITHUB_TOKEN gh` REST/PR readbacks. No clone, repository mutation, fixture/test execution, child/guest/service run, or credential/key/identity change. `/Users/dsteele/repos/capsule-experiments` was absent per parent inspection. Archive has no graph index: immutable recursive Git tree/contents APIs provide source fallback; both inspected recursive trees reported `truncated:false`.

## Exact refs and integration destination

- Default branch: `main`, immutable readback `d8fd74a158f2227277ccfce5532e9bd8ad3e714e`, 2026-09-06T18:18:51Z; merge of PR #37 (`test: retain C5b16 publication timing and fault evidence`). The default tree has no C5b19 path.
- Existing C5b19 draft: [archive PR #38](https://github.com/Shrimpworks/capsule-experiments/pull/38), `test: begin C5b19 build-only teardown oracle`, OPEN, isDraft:true; branch `codex/c5b19-benign-teardown`, base main at the exact ref above, head `1e3f3bb154d0ef24412002ac47adef22ba3c46a9`.
- PR author dills122; created 2026-09-19T21:16:29Z; updated 2026-09-19T21:33:28Z. Head committed 2026-09-19T21:28:31Z. Readback mergeStateStatus:CLEAN, reviews:[], reviewDecision:"", statusCheckRollup:[], issue comments empty. CLEAN alone is not review/evidence completion.
- Existing implementation commits: `38b16f615c014d0f0e6cd53d5e6aa0477ebf1f75` (cancellation); `7e064c9f4cd1eb4fff937a8c8e081595dc3ddc4b` (identity); `cd9e20a145f584f7e4ee66684c2cb39621112621` (obligation); `1ef2dcabeeef6591e5f45ba4a0d0a5a6e6661fcb` (mutation sensitivity); `1e3f3bb154d0ef24412002ac47adef22ba3c46a9` (hardcoded cancellation refusal).
- Pinned C5b16 `0efd03def6bc333a92c8b9809bc56b7b3ce9ea80` is available through an untruncated recursive tree, including `experiments/typed-guest-transport-c5b16-timing-fault/{PLAN.md,go.mod,inputs/ORIGINS.json}`. It is the second parent of default merge d8fd74a. Preserve as immutable reference.
- Other branches read back: codex/c5b14-storage-core, codex/c5b15-runner-reconciliation, codex/c5b-s1-fixed-runner-no-run-successor, main. No second C5b19 branch found; PR list through all 38 entries found no later C5b19 PR. No `.github/` workflow tree on default head.

## Actual source present

C5b19 subtree hash `4ba695ab9caccea8d00c0d21cf379450ba1a728d` at exact PR head contains exactly seven files:

1. `PLAN.md`
2. `go.mod` — module `capsule.experiments/c5b19`, Go1.25.13, no requirements
3. `inputs/ORIGINS.json`
4. `tests/model_test.go`
5. `tests/obligation_test.go`
6. `tests/oracle.go`
7. `tests/oracle_test.go`

The implementation source imports bytes, crypto/sha256, encoding/binary and errors only. It implements a pure 57-byte cancellation validator, a domain-separated process-identity digest and the immutable 295-byte obligation preimage. It has no process/filesystem/network/time/transport authority and no executable entry point. This is not the full state machine, full oracle, storage or native mechanism.

The rest of the Gate-2 source/build manifest is absent: fixture/supervisor/harness C, protocol and clock C/headers, bridge header and Go bridge, Go store and tests, build/check/run scripts, executable manifest and raw run evidence. The existing extra obligation_test.go is a test source split beyond the literal intended Gate-2 file list and must be transparently inventoried by the lead's scope reconciliation; no native/linking bytes were inferred from these file names.

ORIGINS pins archive base d8fd74a, Gate2 reviewed e7b68c57d59abe66c68d849473e21fe74f4e31b7, integration a6ef83a6b43cd4940547b138529ebbb91cff8da6, and C5b16 ref above; Go1.25.13, run_authorized:false, product_admitted:false. Do not silently replace this historical packet pin with product Go1.26.9.

## Test and review evidence bounds

Static test inventory: cancellation literal known answer; mutated magic/version/role/length/each ID/kind; short/trailing/empty; matching changed valid bindings and stale-frame refusal; matching invalid zero/duplicate bindings. Identity tests contain known-answer digest, each observation varied, zero values and PID overflow. Obligation tests contain known-answer digest, all ten field offsets independently varied, zero and duplicate refusals.

PR body reports Go1.25.13 `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./...`, `git diff --check` passing. These are author reports, not executed by this investigation; no hosted checks were present to corroborate.

Retained PLAN records first three test-before-implementation failures as author-observed task history, explicitly not retained raw logs. It records controlled no-child mutants for disabled cancellation binding checks, hardcoded plan binding and hardcoded original cancellation tuple, each reportedly caught then restored. No raw mutation logs found in subtree.

Review history must not reset: instance1/3 at cd9e20a returned Not ready (two P2 test-sensitivity gaps); instance2/3 at1ef2dca returned Not ready (hardcoded cancellation tuple survived); corrected head1e3f3bb. PR body says instance3/3 Ready/no actionable findings on1e3f3bb; retained PLAN still says instance3 pending. No GitHub review/comment artifact corroborates the claimed third verdict. This is provenance/status drift requiring honest canonical reconciliation or locating retained reviewer evidence, not license to invent or retry a fourth pass under the old budget.

## Archive conventions

Only root AGENTS.md appears in exact head tree. It requires defensive local scope, non-production boundary, exact question/versions/fixtures/environment/authority; distinguish observations/inferences; retain positive/negative/fault/mutation evidence; verification/limitations/unresolved questions; canonical decision link back in capsule-corp. Never import archive code into product packages. README places disposable implementation here, with product retaining immutable decision/evidence links. No root build convention/workflow was established beyond the per-experiment Go commands and frozen packet.

## Recommendation

Continue from existing PR38/head after source custody setup, not a replacement initial implementation. Root owns canonical recovery roadmap PR and explicitly names archive PR38 as retained implementation destination if it elects to continue that branch. Reconcile #384/#385 closure and this pre-existing work first. Bound the next child to non-executing parser/model/oracle tests only until actual implementation/run gates close; do not infer complete Gate3 or create a child fixture from the present green claims. Preserve the existing three-pass review history and require an explicitly scoped new review budget only for distinct new work or owner-authorized extension, not a disguised reset.

Confidence: high for exact refs/tree/source inventory and absence in these complete trees; medium for historical tests/review claims because current retained evidence is author-written prose with a contradictory PLAN status. No assertion of runtime correctness, timing, independent reproduction, branch-protection sufficiency or run authorization.

## Verification reads

`gh api repos/Shrimpworks/capsule-experiments`; `gh api repos/Shrimpworks/capsule-experiments/commits/main`; `gh pr list --repo Shrimpworks/capsule-experiments --state all --limit 100`; `gh pr view 38 --repo Shrimpworks/capsule-experiments --json ...`; `gh api --paginate repos/Shrimpworks/capsule-experiments/branches?per_page=100`; immutable `git/trees/<ref>?recursive=1` and `contents/<path>?ref=1e3f3bb...` with raw Accept; `gh api repos/Shrimpworks/capsule-experiments/issues/38/comments`. All ran outside sandbox with inherited token env vars unset; no secrets output.
