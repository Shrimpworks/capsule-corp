# C5b19 G3a2 independent review

Review instance: **1 of 3** for the distinct new wire/record increment and canonical recovery reconciliation. Review work status: **PASSED**. Parent full Gate 3: **IN_PROGRESS — TRENDING_GOOD**; Gate 4–6/product admission remain **BLOCKED**. Verdict: **Ready**.

## Findings

**P3 — Reconcile the remaining stale merge instruction.** `docs/CI_PARITY_2026_10_15.md:128–131` still says merging PR #385 is the remaining delivery step. The changed document opening records the merge, and its new closeout at lines 133–144 describes the completed merge and next priority. A reader landing at the hosted-delivery section receives contradictory current instructions. Smallest correction: change the two sentences to past tense/completed delivery, or explicitly label them as a pre-merge historical checkpoint. This was non-blocking documentation drift, not an implementation defect. **Resolved in this same review instance:** root replaced the stale instruction with completed implementation, verification and merge acceptance; reviewer read back the corrected lines 128–131. No behavior changed and no further review instance is needed.

No actionable decoder defect identified. In particular:

- `tests/oracle.go:137–174` checks exact outer sizes/headers before slicing and validates frozen declaration-order bindings, fixed policy, immutable obligation digest, nonzero bounded generation and preparation operation relationship.
- `tests/oracle.go:177–198` matches the passive model's candidate shapes, including terminal absence after unconfirmed identity publication; nonzero confirmed runner identity must equal absence identity.
- `tests/oracle.go:200–255` requires full record/request equality and exact reply tuple/body hash/outcome/record digest. Failed/indeterminate outcomes require zero digest. No outcome creates publication or lifecycle authority.
- Fields are fixed copied arrays; no aliases to caller buffers, dynamic untrusted length, external effect, product import, dependency addition or source-manifest expansion is introduced.

## Plan Review

Target: archive diff `1e3f3bb154d0ef24412002ac47adef22ba3c46a9..c4d3c33ed04c1dff213e2c69064a6113892d9653`, exactly five paths (PLAN, ORIGINS, oracle.go, oracle_test.go, deleted obligation_test.go). Archive became committed during review without source changes; final tree clean and supplied hashes matched. Canonical target: working diff from `ba28f9c0b33c2366ccb892c2cf5bc33dcf792600` over the five named documentation/planning files and new foundation ledgers. Later implementation-handoff.md/progress additions were disclosed by the root.

| Plan requirement | Evidence and assessment |
| --- | --- |
| Pure request397/reply98/record395 substrate | Implemented with exact offsets and big-endian interpretation; independently reconstructed literal bytes match. |
| Fixed independent vectors, negative cases, mutation sensitivity | Literal known answers plus every-byte correspondence substitutions, all short sizes/trailing byte, varied valid frozen fields, matching invalid zero/duplicate bindings, malformed candidates with matching hashes, all outcomes, caller ownership. Three retained mutation logs fail in the named cases. |
| Frozen source manifest | Obligation tests were relocated byte-for-byte into listed oracle_test.go; no coverage removed. Current six-file inventory matches PLAN; native/store/trace files remain absent and are described as absent. |
| Immutable provenance and toolchain | All four implementation and three canonical reference hashes verified; Go1.25.13 module and actual local binary match; no third-party requirements. |
| Scope and authority | Structural decoding only; stateful freshness, predecessor generation, non-overlap, late settlement, full trace oracle, C/Go equality, storage and live custody explicitly deferred. This is faithful to the narrowed acceptance criteria. |
| Recovery reconciliation | Local history confirms exact #384/#385 merge commits. Current roadmap selects G3a2 and preserves historical third-review uncertainty, full Gate3 incompleteness and first-run gate. The stale merge instruction was corrected and read back in this pass (P3 above). |
| Delivery | Final source pin, review retention, push/PR metadata and canonical publication remain root-owned pending actions. They were explicitly pending in the review bootstrap, so are not falsely treated as already completed. |

The plan does not quietly use structural generation validation as state-machine conformance. In particular, acceptance of any nonzero generation through `2^53-1` is intentional at this layer; exact predecessor relation is deferred. The old three-review flow is preserved and this review does not independently certify its claimed third verdict.

## Author-Claim Reconciliation

Preliminary ledger `/tmp/c5b19-review-1-preliminary.md` was written after requirements/source/tests and before archive PLAN/ORIGINS, author handoff or recovery claim inspection. Neutral bootstrap contained scope/prior history, not author rationale. Author packet: `/tmp/c5b19-wire-handoff.md`.

| Material author claim | Evidence inspected | Status / consequence |
| --- | --- | --- |
| Fixed layouts and exact correspondence | Gate-2 fixed transport/durable sections; passive types/model; full source; independent Python struct/hashlib reconstruction | Confirmed for bounded byte structure. |
| No lifecycle/storage/clock/network effect in tests | Complete oracle.go/oracle_test.go/model_test.go imports and bodies | Confirmed for reviewed source. Go compilation/test processes are normal permitted tooling. |
| 890 byte substitutions, ten zero and twenty duplicate binding pairs | Test loops and fixtures | Confirmed; not exhaustive implementation-mutant proof. |
| Existing obligation tests consolidated unchanged | Exact old file tail compared as contiguous text in oracle_test.go | Confirmed. |
| Red, green, race, vet, build and restored test results | All named wire logs; fresh test/race/vet/build; restored source SHA | Confirmed current behavior; historical red/green ordering remains author's chronology, supported by retained missing-symbol log. |
| Three targeted mutants caught | Three corresponding failure logs plus test/source inspection | Confirmed log contents and plausible sensitivity; reviewer did not mutate shared source or independently rerun mutants. |
| Canonical install/CI/tests pass; full lint remains capped debt | validation.json hashes matched all four logs; test totals 164+10+3=177 and 132 script tests; lint log ends 50 revive findings | Confirmed retained logs, not independently rerun canonical suite. |
| Merge timestamps, canceled suite rerun, unchanged protections, current remote issue/nightly states | Canonical readback narrative and supplied facts; local merge history | Local merge identities confirmed; fresh remote metadata unverified in this pass because gh network calls failed. No contrary evidence. |
| Old review3 Ready claim unresolved | Historical PLAN/PR discrepancy retained in new PLAN, inventory and canonical checkpoint | Confirmed honest uncertainty; no retrospective certification. |

## Verification Performed

Archive working directory: `/tmp/capsule-c5b19-recovery/experiments/typed-guest-transport-c5b19-benign-teardown`.

Exact Go binary: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go`; version readback `go1.25.13 darwin/arm64`. Each Go check used `GOCACHE=/tmp/c5b19-review-go-cache GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off`:

- `go test -count=1 ./tests`: exit 0, 0.254s.
- `go test -race -count=1 ./...`: exit 0, 1.314s.
- `go vet ./...`: exit 0, no diagnostics.
- `go build ./...`: exit 0, no diagnostics.
- Archive `git diff --check 1e3f3bb154d0ef24412002ac47adef22ba3c46a9` and canonical `git diff --check`: exit 0.
- Python standard-library read-only verifier: all implementation/reference/log hashes, independent packet-derived request/reply/record bytes and exact obligation-test relocation passed.
- `gh pr view 385 --repo Shrimpworks/capsule-corp --json state,mergedAt,mergeCommit` and `gh issue view 265 --repo Shrimpworks/capsule-corp --json state,closedAt`, with inherited token env unset: failed to connect to api.github.com. No credential material inspected or emitted.

An initial login-shell setup emitted a parse warning while the first test still ran successfully; subsequent commands used a non-login shell. One file-read tool hook stopped an invocation, then allowed its retry; no source mutation was attempted. Reports are the only reviewer-written artifacts; Go cache writes are ordinary build output.

## Open Questions And Residual Risks

- Exact review checkouts lack graph indexes. Read codebase-memory skill and both list_projects pages; used focused actual source instead of another checkout's graph. No exhaustive graph claim made.
- Unit vectors cover bounded structure and correspondence only. They do not establish durable publication, actual process identity/custody, timing, C/Go ABI equality, cancellation ingestion, replay protection, full trace causality, start or release authority.
- Raw wire/canonical logs in /tmp are ephemeral. PLAN retains named command/results/excerpts and distinguishes them from durable runtime evidence. Root should retain this review and the final exact commit link before delivery.
- Independent remote verification remains limited as above. This review neither verifies branch protections nor certifies hosted runtime evidence.
- Existing G3a third review report remains unavailable; new passing checks cannot recreate historical review provenance.

## Verdict

**Ready** for this exact G3a2 increment and canonical reconciliation, subject to the ordinary root-owned publication completion already pending. No blocking implementation finding or heavy pivot. This verdict does not close full Gate3, admit Gate4, authorize a first child run, or reset/certify the historical review flow.

## Recommended Next Actions

1. P3 wording correction is verified complete; no further action on that finding.
2. Retain this report, exact archive commit and verification boundary in the canonical delivery and existing archive PR38; verify published body/draft state and final pins.
3. Continue only the next explicitly scoped implementation slice; keep lifecycle, storage, host and first-run gates closed until their own evidence and authorization are satisfied.
