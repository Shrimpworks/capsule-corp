# G3a5 independent review

Review instance: 1 of 3 for this distinct slice. Verdict: **Ready** for the bounded local field/lane implementation and its plan. Review status: **PASSED**. Parent full Gate 3 remains **IN_PROGRESS — TRENDING_GOOD**; child execution, Gates 4–6 and product admission remain **BLOCKED**. Publication is still owned by the Brain and was not performed by this reviewer.

## Findings

No actionable findings in the frozen working-tree scope. No heavy pivot or additional review instance is requested.

Reviewed repositories and exact bases:

- Canonical `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, `codex/c5b19-event-fields`, HEAD/base `fe3d22c84b5d41e0a4edbfacb531563040e16995`: three named canonical documentation changes and untracked `.planning/c5b19-event-fields/` evidence. The unrelated `.serena/project.yml` change was excluded and preserved.
- Archive `/tmp/capsule-c5b19-recovery`, `codex/c5b19-benign-teardown`, HEAD/base `3e26b809a40832a3351fdc6b5e0ffc1ef0094792`: only `experiments/typed-guest-transport-c5b19-benign-teardown/{tests/oracle.go,tests/oracle_test.go,PLAN.md,inputs/ORIGINS.json}` changed.

The source table at archive `tests/oracle.go:535` agrees with every row of the accepted scope matrix and Gate2 packet field table. The wrapper at `tests/oracle.go:527` invokes the unchanged envelope decoder before indexing the closed kind range. Its forbidden-field mask covers all 14 non-common fields, required tuples/identities have the accepted presence and range checks, and the remaining conditions preserve zero scalar observations while enforcing local reply/store/wait relations. All refusal paths return the zero event. No lifecycle, history-settlement, file, clock or network behavior is added.

Tests at `tests/oracle_test.go:964`, `:988`, `:1013`, `:1059`, `:1099`, `:1128` and `:1172` cover all 29 kinds and three lanes, each forbidden non-common field, independent required-member omissions, accepted and rejected scalar boundaries, represented reply outcomes/decisions, store27 without a record digest, wait-result/PID relations, zero observations, common envelope refusal, zero error outputs, copying, and explicit opt-in preservation. The literal-offset test matrix does not reuse implementation masks.

## Plan Review

The acceptance plan is coherent and appropriately bounded: scope clarification precedes implementation; tests precede the new wrapper; exact source freeze precedes independent review; final publication and immutable archive links follow review. It preserves the prior envelope contract and frozen manifest, standard-library-only Go1.25.13 archive pin, ownership split, and deny-by-default execution gates. ADR47 and the C5b19 reuse-map row are consistent with a pure fixed parser, with no additional authority or dependency.

The planned 29-kind support is explicitly the table-defined local subset. Optional harness identity and arbitrary malformed-reply diagnostic mappings remain unsupported, and are not incorrectly declared universally invalid under the packet. Scalar zero observations are representable; syscall interpretation, terminal status, sidecars, causal ordering, timing and complete-case validity remain deferred. Those are honest parent dependencies, not failures of this slice. Full publication acceptance is not claimed prematurely in the current canonical prose.

The reviewer read acceptance scope/plan, packet, tests and source before the author explanation. `/tmp/c5b19-g3a5-review-1-preliminary.md` was saved before reading the separate author explanation or archive outcome additions. The acceptance `plan.md` itself contained short outcome summaries, so perfect blinding to all author outcomes was impossible; those summaries were not treated as independent verification.

## Author-Claim Reconciliation

| Author claim | Evidence inspected | Status | Review consequence |
| --- | --- | --- | --- |
| Opt-in 29-kind matrix, unchanged prior APIs | `oracle.go:527–619`, packet field table, accepted matrix; byte-prefix comparison against archive HEAD | Confirmed | Additive bounded layer; old decoder/history behavior retained |
| Independent offsets, all kinds/lanes, forbidden fields and boundaries | New seven test functions and fixture matrix; independent offline full test/race runs | Confirmed | Direct evidence for requested local acceptance/refusal |
| No dependencies, manifest growth, product import or live effects | Exact archive dirty paths, unchanged go.mod/model_test.go, imports and new source/tests | Confirmed | No expanded runtime or authority surface |
| Compiler-red before implementation, green afterward | Matching SHA-256 and diagnostic contents of red/green logs; undefined decodeEventFields at test line938 | Confirmed as retained evidence | Supports test-first record; historical sequence was not independently witnessed |
| Eight intended guard mutants fail and source restored | Mutation JSON and all eight matching log hashes/assertions; final source hashes and independent passing tests | Confirmed as retained evidence | Selected sensitivity established by inspected logs; mutations were not independently rerun |
| Prior provenance/source preserved | Parsed ORIGINS comparison excluding only appended g3a5; exact old Go prefixes; unchanged module/model | Confirmed | No retroactive reinterpretation of earlier slices |
| Canonical baseline checks pass with disclosed lint/audit limits | Four canonical command-log hashes, Makefile recipes, observed test counts/results and full lint output | Confirmed as retained evidence | No false all-clean claim; independent canonical full CI rerun was not needed for documentation-only scope |
| No full trace/lifecycle/admission claim | New source comments, archive PLAN/ORIGINS, canonical recovery/plan text | Confirmed | Narrow readiness only |
| Final immutable pins, PR publication and hosted acceptance | Plan identifies these as future Brain delivery | Unverified / pending by design | Complete before reporting integrated delivery; not an implementation defect |

## Verification Performed

Independently executed in the exact archive experiment directory, with `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/c5b19-g3a5-review-cache` and binary `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go`:

- `test -count=1 ./...`: exit0, package tests passed (0.512s).
- `test -race -count=1 ./...`: exit0, package tests passed (2.087s).
- `vet ./...`: exit0, no diagnostics.
- `build ./...`: exit0, no diagnostics.
- `git diff --check` in both worktrees: exit0.

Independent read-only Python comparisons confirmed prior Go source prefixes, unchanged go.mod/model_test.go, all prior ORIGINS values, all four implementation hashes, three canonical reference hashes, all fourteen archive evidence-log hashes and all four canonical command-log hashes. Relevant diagnostic lines and test counts were inspected. Canonical evidence reports 177 workspace tests plus132 script tests, no govulncheck findings, two low plus two moderate dependency advisories below the unchanged threshold, and exit1 from unfiltered lint with50 output-capped revive findings. This review did not rerun canonical package installation, CI, scanner or unfiltered lint.

Frozen archive SHA-256 values match dispatch:

- oracle.go: `d82c587990d8516e89b442d1d89707fed7a2114a6a906c9b76a624dfd39fd114`
- oracle_test.go: `fb643536635b8e73fb2c2d52aeae782b7350162613b09577a01cbc50f34ad410`
- PLAN.md: `c64725ced2c26f9f43ef34a1cd4db120cf7a5f4e1075d309ffea0f03c6b67d28`
- ORIGINS.json: `8059c0f5f0a9d5cec04838b687d96f8bb278fb07ab6cf5c80d0c9202872ad944`

Serena was initialized and activated for the target canonical worktree. Codebase-memory index status and cited-path coverage were checked, generation `2026-10-04T22:18:10Z` at the primary checkout. Canonical metadata is best-effort only; new planning/recovery paths have missing freshness and were read exactly. Structural model queries produced no matches, so exact current source supplied the fallback. Archive source is unindexed and was reviewed directly. No graph-based exhaustive or completeness claim is made.

## Open Questions And Residual Risks

Confidence is high for the accepted fixed local matrix, with selected boundary/guard evidence rather than exhaustive protocol proof. This API can accept an isolated absence shape or trace-fault shape; that is intentional and grants no evidence of successful cleanup or a passing case. A future full oracle must supply sidecar equality, capture completeness, causal/timing composition and native semantics, including the two unsupported diagnostic/identity encodings before using them.

Canonical Node22.23.3 differs from documented22.22.1; known unfiltered lint debt and lower-severity advisories remain disclosed, unchanged baseline limitations. Mutation source variants were not independently reconstructed; this review confirms retained failing assertions/hashes and current restoration, not witnessed historical mutation operations. Final metadata additions after this review require the parent's final formatting/diff check and publication readback.

Defensive authorized scope was owned inert synthetic in-memory byte fixtures and ordinary offline compiler/unit-test processes. No native fixture, guest, lifecycle target, signal/wait, live clock/store/network, credential or user-data operation was invoked. No source fixes, repository commits, branch changes, pushes, PR edits, subagents or broader workstreams were created. Only the requested temporary review artifacts were written.

## Verdict

**Ready** for the frozen G3a5 bounded implementation and plan. Review slice **PASSED**; integrated delivery remains **IN_PROGRESS — TRENDING_GOOD** until the Brain completes publication. This verdict does not close full Gate3, authorize a run, or admit a product control.

## Recommended Next Actions

Brain should retain this report, reconcile review status in canonical/archive evidence, commit/push the reviewed bytes through the named integration destinations, pin the resulting immutable archive commit in canonical documentation, complete final metadata checks, and read back PR body/newlines/base/head/draft state and hosted results. No behavioral fix is requested. Preserve all existing run/admission blockers and historical review budgets. This is instance1/3 only; no further review instance was started.
