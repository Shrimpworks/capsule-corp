# NEXT2-SCOPE handoff

Work item: Read-only scoping of the first C5b19 Gate-3 implementation increment.
Status: PASSED.
Scope: Canonical Capsule plan, Gate-2 packet, ADR-0047, C5b18 passive bindings/model, and ecosystem custom-primitive policy. No repository edits, no lifecycle calls, no child fixture, no signal, service, credential, guest, or external target.
Evidence: Focused current-checkout reads and graph coverage below.
Remaining work: Implement and independently review the selected increment; full Gate-3 composition, Gate-4 implemented-byte review/owner first-run authorization, Gate-5 execution, Gate-6 conclusion remain outside this scope.
Next action: Continue and reconcile existing archive draft PR #38 (head 1e3f3bb154d0ef24412002ac47adef22ba3c46a9), preserving its history and review count; do not recreate the already-implemented foundation. Confirm exact offline Go 1.25.13 availability before unit tests.
Parent status: IN_PROGRESS — TRENDING_GOOD. Runnable C5b19 evidence and product admission remain BLOCKED.

## Archive reconciliation update

The parallel archive investigator reports draft PR #38, branch `codex/c5b19-benign-teardown`, head `1e3f3bb154d0ef24412002ac47adef22ba3c46a9`, based on main `d8fd74a158f2227277ccfce5532e9bd8ad3e714e`. Its seven files already implement the cancellation parser, process identity digest and obligation encoding subset suggested below. Continue that delivery unit; the following list describes allowed scope, not instructions to duplicate its source. Archive investigator reports PR body claims third review Ready, PLAN says third pending, and GitHub reviews/comments/checks are empty: reconcile retained evidence before claiming acceptance; preserve the three-pass history.

Reported source list includes `tests/obligation_test.go`, which is not in the Gate-2 frozen list. Preferred bounded correction is consolidate these tests into the already-listed `tests/oracle_test.go` or `tests/model_test.go`, preserving coverage. Do not implicitly amend the reviewed manifest or claim full Gate-3/4 completion. Root must read archive rules and existing code before selecting the concrete patch; this scoped analysis has not reviewed PR38 implementation.

## Decision

The next critical security path is C5b19, not a newly invented runtime or a broad cleanup. The current plan explicitly names the reviewed Gate-2 packet as next security experiment. A pure pre-run parser/model/oracle increment can begin without widening the packet or adding an unlisted source file. Gate-2 lines 446–453 explicitly require failing parser/model/oracle tests first and require simulated process/store effects for every pre-run test. The preceding reviewed conditional owner authorization is recorded at packet lines 15–22 and plan lines 386–390. The old generic stop-rule paragraph at packet lines 480–485 must not be misread as canceling that later explicit recorded authorization.

Recommended smallest useful increment: **C5b19 G3a independent bounded wire/record oracle substrate**. Complete exact independent decoding, digest/binding validation, and frozen request/reply/record correspondence for inert in-memory byte vectors; optionally cancellation/process-identity known answers in the same increment. It is a parser/conformance foundation, not the complete runtime oracle, store implementation, platform mechanism or gate closure.

Use only these already-listed archive files:

- `experiments/typed-guest-transport-c5b19-benign-teardown/PLAN.md`: acceptance, scope, dependency/custom-primitive checklist, exact commands, red/green results, known limitations and handoff.
- `go.mod`: exact `go 1.25.13`, no third-party requirement and no go.sum.
- `inputs/ORIGINS.json`: immutable Capsule Gate-2 reviewed head `e7b68c57d59abe66c68d849473e21fe74f4e31b7`, exact current source reference commit for C5b18, source hashes, test-only provenance. Do not identify this new source as retained Gate-2 bytes.
- `tests/oracle.go`: pure bounded decoders/validators and independently derived expected hashes. Keep it compatible with its future compiled read-only-oracle role; no executable main/live targets are needed for this partial slice.
- `tests/oracle_test.go`: literal known-answer fixtures, malformed-frame/digest/binding negatives, exact bounded checks.
- `tests/model_test.go`: pure simulated request/settlement relation tests if helpful; omit this file rather than create an empty placeholder if not used.
- `scripts/check.sh` only if needed: offline explicit exact-toolchain test command with no fixture invocation; permit only its pre-run unit-test responsibility. Inline vectors avoid new `testdata` files outside frozen manifest.

Do not add a generic `source/protocol` package, arbitrary test files, FFI adapter, real durable store or extra generated source to make this slice convenient. C encoding can follow in the already-listed `inputs/protocol.h`/`.c`, but separate actual C execution requires carefully accounted test runners; avoiding C for G3a keeps the first increment narrow. A Go parser alone is not C/Go byte equality.

## Frozen requirements that must survive implementation

Packet lines 187–241 fix request=397 bytes, reply=98, cancellation=57; record lines 253–269 fix 395 bytes. Reject shorter and trailing bytes, incorrect length, magic/version/role/kind/outcome and noncanonical booleans/disposition. Encode/decode integers big-endian and do not serialize native structs or rely on padding.

Bindings[240] use **declaration order**, not grouped-by-size order:
InstallationID[16], TrustEpochDigest[32], SupervisorIdentityDigest[32], AttemptID[16], ApprovalID[16], RegistrationID[16], PlanDigest[32], ProfileDigest[32], RunnerDigest[32], PreparationOperationID[16]. All 16-byte identities are nonzero and pairwise different; likewise the five digests. The policy is exactly 1000,1000,1200. Preparation operation must equal immutable preparation ID; later operation identities must be fresh in stateful simulation. Generation is nonzero, advances exactly once per settled predecessor, and preserves C5b18 maximum `2^53-1`; a stateless frame decoder cannot alone establish freshness.

The obligation SHA-256 domain/preimage is immutable and remains equal across operations; never include candidate, process identity, generation or operation in it. Request hash covers exact request body, excluding the 8-byte header. The record digest hashes all exact 395 record bytes. Candidate has prepared, runner-confirmed, runner identity, absence identity, terminal-confirmed, terminal disposition in that order (68 bytes). Failed/indeterminate reply has all-zero record digest; confirmed reply matches the recomputed exact candidate record.

C5b18 candidate transitions require positive preparation, preparation+runner before starting, and terminal evidence only after authoritative absence. **A terminal candidate can have runner-confirmed=false with zero runner identity and still bind nonzero exact absence identity** when identity publication failed; do not accidentally require runner confirmation for all valid terminal candidates. `TerminalTimingViolated` is retained failure evidence, not successful completion. A parser-valid terminal record/reply cannot prove durable publication, process custody, public completion or capacity release.

Process identity preimage and cancellation fields are fixed at packet lines 121–129 and 228–236. A synthetic PID inside an inert known-answer byte vector is not a real process and grants no signal authority. Invalid zero/out-of-range PID and substituted fixture/attempt/digest values must refuse. Any parser API must receive bounded in-memory bytes, not filesystem paths or process targets.

## Tests and acceptance

1. TDD: commit/retain red evidence for named missing-behavior tests before implementation; do not just assert roundtrip against the same implementation. Literal independent byte strings and fixed expected SHA-256 values are needed.
2. Positive known answers for each selected frame, obligation/process preimage and stored record; documented offsets and exact bytes. Cross-operation obligation stability, request/record hash relationships and copied inputs survive caller mutation.
3. Negative mutation table for every header field; short/trailing data; every immutable binding zero/substitution/duplicate; fixed policy and mask changes; generation 0/overflow; unknown enums/booleans; illegal candidate combinations; exact predecessor/operation mismatch; nonzero failure digest; confirmed mismatch; request-body/candidate/record substitution. Byte mutation is permitted to produce another syntactically valid value only if contextual binding validation rejects it; separate parser structure from expected-binding policy.
4. Pure settlement tests: no overlapping operation; late confirmed reply after stop never opens start, clears stop or releases capacity; failed/indeterminate outcome remains uncertainty; terminal record never supplies public completion. If these would require duplicating the entire control state machine now, explicitly defer them and name G3a **wire/record** only. Do not claim model completion from field validation.
5. All tests run offline with already-available exact Go 1.25.13, `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off` (or invoke its absolute binary and validate version first), no added modules and no fixture compilation/execution. Go test's normal unit-test process is not the forbidden benign lifecycle fixture; test code itself must not spawn subprocesses.
6. Independent fresh-context engineering/security review compares bytes/tests directly with the exact packet and C5b18. Every blocking finding resolved within configured review limit. Retain archive PR/commit, checks, red/green and review results; brain records exact commit-pinned conclusion in Capsule roadmap and its delivery PR.

G3a PASSED means this explicitly listed parser/conformance subset and its evidence close. Full G3 remains IN_PROGRESS — TRENDING_GOOD while active; G4–6 remain BLOCKED. Do not use full `oracle passes` language: the raw trace/cross-lane causal/timing oracle remains future work.

## Dependencies and policy checklist

Applicable ecosystem row: `C5b19 prospective Supervisor control/storage lane exchange` (ECOSYSTEM_REUSE_AND_ADOPTION.md line 135): ADOPT-PLATFORM / BUILD-NARROWLY; no generic bus or queue. G3a uses only the BUILD-NARROWLY inert wire part and Go standard SHA-256/big-endian support. Native pipe, poll, fcntl, clocks, waiters, storage lock/fsync/rename and lifecycle behavior are deferred.

Complete the policy block in PLAN.md for this custom primitive before implementing it, even though no package is added. Name trust class test/evidence; fixed byte caps; authority solely parsing in-memory synthetic bytes; no process, network, user filesystem or credential authority; exact Go/toolchain provenance/license inventory and offline reproduction; mutation cases; maintainer upgrade/removal ownership. No dependency admission follows.

Exact archive Go pin is **1.25.13**, not Capsule's newly upgraded root 1.26.9. The packet requires pinned local tools and no network fetch. If it is absent, do not silently auto-download or substitute 1.26.9. Continue roadmap/packet planning and report the toolchain prerequisite, or have the brain take an explicitly reviewed packet-revision slice. Current host/SDK changes also cannot silently become the run profile; Gate4 must recheck its exact host and reopen packet review on mismatch. Pure unit tests do not close or refresh the host gate.

Storage code is possible later within `source/store/store.go` and `_test.go` but **real filesystem fault injection is not a pure no-effect slice**. Pre-run tests here must use simulated store effects; no owner lock, actual fsync/rename or delayed worker is necessary for G3a. Full store behavior still owes A–D fault outcomes, ownership and publication proof. Trace parser/oracle implementation is also permitted in listed files but 29 event kinds, lane caps, causal pairings, separate writer sequences and timing rules make it materially larger than G3a.

## Source/coverage and confidence

Graph project `capsule-corp` lists primary checkout, full generation 2026-10-04T22:18:10Z. Coverage of Gate2, experiment plan, ADR0047, ecosystem map, teardownpassive types/model and scope reported metadata_match/no_recorded_issue, best effort only. Exact current managed-worktree source was read; no changes between indexed root commit fc8af2d and current head for types/model. Initial broad types symbol search returned page1 only and was not used for an exhaustive structural claim; targeted Struct search returned all7 results and DurableProjection snippet confirmed six fields. Current plan and recovery report freshness require direct reads (parent covers roadmap).

Confidence high on allowed no-effect scope and exact wire requirements, moderate on minimal archive packaging pending archive investigator's repository rules/workflow findings. No claim about remote archive contents or local Go1.25.13 availability; those were not examined to avoid duplication. No code tests run in this read-only assignment. No branch, commit, or PR created.
