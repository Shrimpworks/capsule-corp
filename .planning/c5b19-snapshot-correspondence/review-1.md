# C5b19 G3a9 independent review

Review instance: 1 of 3. Review status: PASSED. Verdict: Ready for the bounded supplied-correspondence increment. Parent full Gate3 remains IN_PROGRESS — TRENDING_GOOD; child execution, Gates4–6 and product admission remain BLOCKED. Integration/publication is still Brain-owned and incomplete.

## Findings

No actionable findings in the frozen implementation or accepted bounded plan. No heavy-pivot gate is needed. Passing this review does not authorize a child/native run, establish physical durability, or complete Gate3.

The source-first preliminary ledger is `/tmp/c5b19-g3a9-review-1-preliminary.md`. It was written before reading implementation-handoff.md, canonical-verification.md, archive PLAN rationale or ORIGINS.g3a9. The required scope.md and plan.md were read first; they contain design rationale and implementation-freeze summaries, so the first pass was not blind to those summaries. Author claims were subsequently reconciled with source and evidence, not used as the readiness verdict.

## Exact scope and environment

- Canonical: `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-snapshot-correspondence`, HEAD/base `f3f2e0755c0c284dd28019225bef0cf8482d3cdd`. Reviewed dirty `docs/CURRENT_WORK_PLAN.md`, `docs/PROJECT_RECOVERY_2026_10_15.md`, `docs/C5B19_BENIGN_TEARDOWN_EXPERIMENT_PLAN.md`, and all six files in `.planning/c5b19-snapshot-correspondence/`.
- Archive: `/tmp/capsule-c5b19-recovery`, branch `codex/c5b19-benign-teardown`, HEAD/base `febb88d6421058f048f75319baae7cf9bd5a0cbb`. Reviewed dirty `experiments/typed-guest-transport-c5b19-benign-teardown/{tests/oracle.go,tests/oracle_test.go,PLAN.md,inputs/ORIGINS.json}`.
- Both branches, HEADs and dirty scopes were verified. No added commits. Unrelated canonical `.serena/project.yml` is excluded and preserved. No repository edits, staging, commits, branches, pushes, PR operations, additional reviewers or subagents were performed.
- Defensive intent: verify the separation of supplied record visibility, represented accepted control settlement and optional worker evidence using only owned inert bytes and ordinary offline Go unit/compiler processes. No mutation script execution, native fixture, child/guest, signal/wait/lifecycle target, live store/clock/network or external target was used.

## Plan Review

The accepted scope is faithfully implemented without expanding the Supervisor, runtime, dependency or serialization boundary. The new `snapshotHistory` is a fixed-size value wrapper, not a store reader or a record-import authority. `oracle.go:900–944` delegates to unchanged G3a8 on trial copies; only successful begin retains a copied request, and only accepted code1 stores reply/outcome against an exact existing request. Code2 preserves prior accepted evidence while retaining inherited fatal/decision effects.

`oracle.go:951–1013` returns supplied facts without mutating history. Explicit absence requires no unused bytes/worker and preserves historical confirmation/pending state. Present records strictly decode and match a retained submission. Worker kind5 and exact request linkage are mandatory when evidence is supplied; accepted and worker replies must agree byte-for-byte. Outcomes2/3 can coexist with visible bytes. `LaterSubmitted` is entry order only and cannot roll back current projection. No full-confirmation, final-store consistency or admission flag is introduced.

This matches packet253–271 and342–359, ADR0047, and normative `internal/execution/teardownpassive/model.go:185–211`: every accepted outcome settles pending state, only confirmed updates durable projection, and late settlement cannot clear control latches. The C5b19 reuse row at `docs/ECOSYSTEM_REUSE_AND_ADOPTION.md:135` stays ADOPT-PLATFORM / BUILD-NARROWLY, test-only. No new dependency or privileged responsibility requires a new adoption or architecture decision here.

Tests cover the ten accepted categories with complete-state equality on successful/refused queries, caller/result value ownership, chronology/bounds and independent valid substitutions. Fixture helpers are composed from existing strict decoders; the semantic mismatch cases explicitly validate the alternative fixture before testing refusal. Strict digest checks intentionally overlap request/record equality; the plan correctly avoids claiming isolated necessity for redundant guards.

The next-scope recommendation remains a recommendation for supplied worker-received/sent slots; it does not silently implement capture authenticity/completeness or authorize execution. Canonical roadmap deltas keep full Gate3/later gates open and retain historical G3a8 status. Commit-pinning, review retention, final metadata/hosted checks and both PR destinations remain delivery tasks for Brain, not completed by this reviewer. The existing archive draft PR38 is the documented integration destination; its live state was not queried.

## Author-Claim Reconciliation

| Author claim | Evidence inspected | Status / consequence |
| --- | --- | --- |
| Exactly four archive paths, prior implementations unchanged | Git status/diff; HEAD byte-prefix comparison of oracle.go, oracle_test.go and PLAN.md; exact go.mod/model_test.go comparison | Confirmed. 138 implementation and472 test lines appended; no prior byte rewrites. |
| Whole-value atomic transitions and read-only query | oracle.go:875–1013, existing decoder/history composition, ownership/chronology tests | Confirmed. No borrowed slice or imported history path. |
| Visibility, accepted outcome and worker evidence remain separate | Tests at2244–2418 and2470–2490; inspect result assignments | Confirmed. Unknown/missing worker never creates settlement. |
| Historical candidates and absence do not rewrite current history | Historical-candidate/absence tests and complete-state query helpers | Confirmed within represented history only; no physical ordering claim. |
| All previous provenance/run/admission values preserved | Parsed current ORIGINS minus g3a9 equals parsed HEAD ORIGINS | Confirmed. |
| Ten new tests and twelve compiled mutants reach intended assertions | Full new test source; exact mutation definitions/script; all12 raw logs, hash and assertion-count checks | Confirmed as retained execution evidence. Mutants were inspected, not rerun. |
| Exact restoration and test/race/vet/build passed | Source hashes, mutation script finally restoration, restored log; independent reviewer commands | Confirmed for final restored bytes. Script restores each mutant byte-exactly and runs a full restored suite after all mutants; it does not run a restored suite after each individual mutant. |
| Canonical baseline checks passed with disclosed debt | validation.json, canonical-verification.md and four matching raw log hashes | Confirmed as inspected worker evidence, not reviewer-executed canonical checks. Unfiltered lint exits1;2 low/2 moderate audit findings remain. |
| PR406 hosted CI/CodeQL passed without bypass | Local plan/roadmap assertion only | Not independently verified; network/hosted verification is outside this local review. Brain must retain/read back hosted evidence before delivery. |

## Verification Performed

Reviewer executed the following from `/tmp/capsule-c5b19-recovery/experiments/typed-guest-transport-c5b19-benign-teardown`, with exact `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go`, `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/c5b19-g3a9-review1-go-cache`:

| Command | Result | Reviewer log SHA256 |
| --- | --- | --- |
| `go test -count=1 ./...` | exit0;0.406s | `f2aec33b001c30d604c5ff8c068c5052d4169967e443d2f57cbd2abe3d0ef80d` |
| `go test -race -count=1 ./...` | exit0;1.909s | `3cc2d07f2050de42845416a794b5bcdcacc2ee0e2cba9a8ef69f5df20e010007` |
| `go vet ./...` | exit0; no diagnostics | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| `go build ./...` | exit0; no diagnostics | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |

Logs are `/tmp/c5b19-g3a9-review1-{test,race,vet,build}.log`. Both worktrees' `git diff --check` exited0. All author source/reference hashes and all20 archive author log hashes matched ORIGINS.g3a9. All four canonical baseline log hashes matched validation.json. Canonical CI and standalone test logs show177 workspace and132 script tests with zero failures; scanner says no vulnerabilities, audit separately reports2 low/2 moderate, unfiltered lint shows50 capped revive findings. No broad canonical rerun was needed for documentation-only changes; final metadata checks remain pending.

Inspected author compiler-red log shows undefined new snapshot API before implementation; green/refined/baseline/restored logs pass. The mutation script hash is `9b0f3574314feb25d9cf0cf59e70f4fa5a7b04780f8cbdff9ef66e8c66901726`. It bounds replacements to the new wrapper suffix, asserts exactly one replacement site, executes focused tests, rejects build failures/panics, checks exact assertion counts, restores in finally, checks byte equality, then runs the full restored suite. The initially guessed plural script filename did not exist; the actual singular `/tmp/c5b19-g3a9-mutate.py` was located and inspected, never executed.

Retained mutation assertions checked against current source and raw logs:

| Mutation | Intended source assertion | Count |
| --- | --- | --- |
| membership-and-accepted-digest | oracle_test.go:2398 invalid correspondence accepted |4|
| worker-kind | oracle_test.go:2432 invalid correspondence accepted |1|
| worker-request | oracle_test.go:2591 invalid correspondence accepted |8|
| worker-accepted-reply | oracle_test.go:2300 invalid correspondence accepted |1|
| accepted-worker-separation | oracle_test.go:2260 worker/record promoted settlement |1|
| code2-evidence | oracle_test.go:2484 code2 overwrote evidence |1|
| later-submitted | oracle_test.go:2324 wrong submission order |1|
| query-fatal-mutation | oracle_test.go:2482 private history mutated |1|
| begin-request-bookkeeping | oracle_test.go:2252 expected valid query returns snapshot error |1|
| accepted-outcome-bookkeeping | oracle_test.go:2288 accepted/worker conflated |1|
| reconcile-commit | oracle_test.go:2288 accepted/worker conflated |1|
| absence-unused-evidence | oracle_test.go:2227 invalid correspondence accepted |1|

Membership's combined bypass is accurately disclosed: it removes membership selection and accepted-confirmed digest comparison, preserving valid history setup so all four independently valid unsubmitted records reach the intended refusal assertion. Worker-request uses eight valid failed/indeterminate foreign requests, avoiding an earlier confirmed-digest failure. Worker equality uses failed versus indeterminate zero-digest replies. None is credited merely for a nonzero command exit. Bookkeeping's expected-valid query failure is the intended positive behavior check, not an unrelated setup failure.

Frozen archive hashes matched author handoff:

- oracle.go: `d3287284d42f5280fd4d2f80afeb898ac652e909697fdf140cefef17a0e77484`.
- oracle_test.go: `3d63c8885e15d4c0ea5854272ffaab6b618882cae11e73ce21ddd1cff387c1b9`.
- PLAN.md: `c500bb6ce464255e475b35fa205d13997c53b92a37ffde7d0168700b2d56898c`.
- ORIGINS.json: `c41f2e64ee9dc9da27f0c12c1a510ff803cfe34ad653632d88768bc934f6dd18`.

Serena initialized; codebase-memory Tier2 index/status, search, trace, snippet and coverage checks performed for all relied-on repository paths. The primary index is old, unwatched and not these worktrees; archive paths are excluded and new planning paths have missing freshness. Exact current source/docs and diffs therefore supplied evidence. SettleWrite's four-callee trace and snippet agree with current model source. No graph completeness claim is made. One accumulated-search hook blocked an ancillary read once and permitted the retry; no task action remained blocked.

## Open Questions And Residual Risks

No unresolved question blocks this accepted bounded contract. Confidence is high for the supplied-byte behavior exercised here, not for native execution or physical persistence. Full capture completeness/authenticity, worker receipt proof, snapshot timing/final-store contradiction, interleaved control groups, raw process-derived projection, control admission, fsync/rename/custody, deadlines and restart recovery remain deferred by design. Existing canonical lint debt and audit findings are baseline limitations, not newly introduced defects. Raw /tmp logs are ephemeral; retained definitions, assertion identities and hashes support auditability but do not replace release-grade capture evidence.

## Verdict

Ready for Brain to retain and publish this bounded G3a9 change through the documented shared destinations. This review is PASSED; full Gate3 and delivery are not complete. No architecture pivot, scope split or additional review instance is requested. Material post-review implementation changes would require Brain to reassess review needs within the unchanged three-instance budget.

## Recommended Next Actions

Brain should retain this preliminary/final review, update pending-review status, complete immutable archive commit pinning and canonical metadata checks, and publish/read back both shared PR destinations with their intended draft state and hosted gates. Preserve every stated limitation and keep child execution, later gates and admission BLOCKED. Do not treat the next-scope recommendation as authorization to implement or execute it.
