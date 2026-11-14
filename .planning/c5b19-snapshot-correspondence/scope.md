# C5b19 G3a9 supplied durable-snapshot correspondence scope

Status: PASSED for read-only scope derivation. Parent Gate3 remains IN_PROGRESS — TRENDING_GOOD; child execution, Gates4–6 and product admission remain BLOCKED.

Internal G3a9-SCOPE assignment. Defensive intent: correlate owned inert synthetic record/request/reply/event bytes with a bounded represented history, preserving the difference between record presence, worker outcome and accepted control settlement. No tests, mutation scripts, native fixture, child/guest, live store/clock/network, lifecycle operation or external target were executed. No repository files, branches, commits or PRs were changed. This report is the only retained assignment output; Brain owns integration and publication.

## Sources and method

Canonical base `f3f2e0755c0c284dd28019225bef0cf8482d3cdd`, branch `codex/c5b19-snapshot-correspondence`, root `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`. Archive base `febb88d6421058f048f75319baae7cf9bd5a0cbb`, branch `codex/c5b19-benign-teardown`, root `/tmp/capsule-c5b19-recovery`. Archive source below is `experiments/typed-guest-transport-c5b19-benign-teardown/tests/oracle.go`.

Read applicable contributor boundaries, starter-document status/authority boundaries, ADR0047, C5b19 reuse recommendation, current delivery plan and prior next-scope recommendation. Material grounding:

- Gate2 packet253–271 fixes the395-byte record, exact request/candidate fields and digest relation. The record stores no outcome. Presence after failed/unknown publication never proves confirmation; recovery preserves uncertainty.
- Packet273–280 includes after-rename, after-directory-fsync and dropped-reply faults. A visible later candidate can coexist with an earlier accepted confirmed projection.
- Packet342–359 requires exact control settlement and a received/applied pair for claimed confirmation, and compares those bytes to worker reply and durable snapshot. Worker-only evidence grants no control authority.
- ADR0047 preserves stop/fence state across late replies and keeps persisted obligation/uncertainty distinct from custody, authority and completion.
- `decodeRecord`, `validateRecordForRequest` and `validateReply` (oracle144–251) strictly bind record/request fields and recompute confirmed reply digest; outcome2/3 replies carry zero record digest.
- `writeHistory` (oracle260–390) retains pending and last-settled state, at-most-three unique submitted writes and current claimed confirmed projection; it does not retain every earlier raw request.
- `decisionHistory` (oracle775–875) validates serial kind3 / whole kind7+8 action groups, exact private settlement, atomic refusal and rejection-derived fatal; its prior API and bytes remain unchanged.
- Normative `Model.SettleWrite`185–211 checks exact pending, settles every valid outcome and updates durable projection only for confirmed. Its four direct callees are applyDurableProjection, samePending, settledWrite and validWriteOutcome.

Serena initialized. Codebase-memory Tier2 project list returned all78 projects; symbol search and depth-one trace above were complete and snippet read. Primary-checkout index generation2026-10-04T22:18:10Z is old/not watched, not the reviewed worktree. Coverage of starter docs/packet/model/types/ADR/reuse reported best-effort metadata_match; planning missing and archive outside_project used direct current reads. No whole-graph completeness claim. This slice composes existing standard-library/local primitives under the C5b19 ADOPT-PLATFORM / BUILD-NARROWLY recommendation; no dependency, serialization, manifest, packet, privileged responsibility or product admission changes.

## Minimal accepted-contract candidate

Add a new bounded wrapper (illustrative name `snapshotHistory`) around an internally owned **value copy** of G3a8 `decisionHistory`. Retain at most three value-owned entries, one per successfully recorded begin: exact397-byte request, whether code1 settlement occurred, and its exact98-byte accepted reply/outcome. Do not borrow caller slices, attach caller-created history, load a record into history, or modify prior G3a8 fields/APIs. The constructor validates frozen bindings/timebase by constructing G3a8 normally. No reset or import path.

`begin` delegates to G3a8 on a complete trial copy, then records the copied exact request in the next private slot only on success. The existing one-pending/generation/unique-ID/kind/terminal/three-write/control-watermark rules still apply. `reconcile` delegates to G3a8 on a complete trial copy; successful code1 associates its exact reply and outcome with the already recorded exact request. Successful code2 may latch G3a8 fatal but **never** creates a submission entry or overwrites accepted evidence. Any refusal commits nothing anywhere, including hidden slots, counts, pending bytes, watermarks and fatal. Begin remains history recording, not permission to submit after fatal or stop.

The read-only query takes:

1. Explicit `present bool` for the supplied observation.
2. Supplied record bytes.
3. Optional captured worker pair (existing `capturedFramePair`, represented by nil when absent).

Return only copied facts such as `Present`, matched `wireOperation`, `LaterSubmitted`, `AcceptedOutcome`, `WorkerPresent`, `WorkerOutcome`, and a copied G3a8 snapshot. Names may follow local style, but their semantics below are mandatory. Do **not** add a `Confirmed`, `FullyCorroborated`, admission or final-case-verdict boolean.

## Query rules

| Supplied case | Required handling |
| --- | --- |
| nil/uninitialized wrapper | Refuse with no mutation |
| present=false, zero record bytes, no worker | Return supplied absence, zero matched/outcome/worker fields and unchanged copied G3a8 history |
| present=false with any record bytes or supplied worker | Refuse ambiguous/unused evidence |
| present=true with malformed/short/extra/unbound record | Refuse using strict existing record decoding |
| present=true, valid record but no exact retained submission match | Refuse; a caller-supplied valid request or rejected reply is not submission evidence |
| present=true, exact retained submission match | Return that copied operation and its retained accepted-outcome evidence, if any |
| worker omitted | WorkerPresent=false, WorkerOutcome=0; missing evidence stays unknown |
| worker supplied | Require strict G3a6 kind5 worker-reply capture, same frozen context, exact copied submitted request, and valid exact reply; then return WorkerPresent=true and its decoded outcome |
| accepted reply exists and worker supplied | Require byte-for-byte equality between accepted reply and supplied worker reply; contradictory but individually well-formed outcomes refuse |

For present records, use `validateRecordForRequest` against the selected private exact request (or equivalent strict operation equality after both decoders). For accepted outcome1, compare SHA256(supplied395-byte record) with the retained accepted reply's record digest. For supplied worker outcome1, compare the same hash with its reply digest. Outcomes2/3 require their existing zero-digest reply convention but **do not require the supplied record to be absent**. A failed/indeterminate write can leave the exact candidate visible.

`AcceptedOutcome` is0 if no successful represented code1 for this entry, otherwise its validated outcome1/2/3. Code2 does not supply accepted outcome. A worker-only outcome1 remains `AcceptedOutcome=0`. A successful code1 outcome1 remains an historical represented control fact even when worker evidence is omitted; omission must never be labeled fully corroborated confirmation. No result claims actual fsync, complete capture or native publication evidence.

`LaterSubmitted` is true exactly when a later request entry was successfully recorded in this selected history. It describes submission order only. It must not be named or interpreted as proven physical supersession, stale final-store state, rollback or corruption. A previously submitted record remains queryable after newer submissions or confirmations. Its own accepted outcome remains historical; it must never replace G3a8's current confirmed projection. Query current and older submitted candidates without inferring unobserved publication or snapshot capture order.

Absent observations remain queryable even with existing accepted confirmation or pending work. The copied history still preserves those facts. A successful absent query means only that the supplied input explicitly represents no record; it does **not** prove a consistent post-history final store, erase confirmation, prove never-published, or resolve uncertainty. Checking final capture timing/completeness and contradiction against physical publication requires a later trace/capture contract. No timestamp/sequence is invented for this query.

Every query, successful or refused, is nonmutating: no accepted-outcome/fatal/projection/pending changes, no slot/count/watermark consumption, no worker evidence retention, and no alias to caller bytes. Repeated queries do not consume the three-action/decision budget. Optional worker sequence/ticks undergo existing strict field validation but do not acquire control-lane ordering or timing authority.

## Required tests

1. Constructor/nil/uninitialized validation; explicit empty-history absence; present malformed/empty record; absent with bytes or worker refuses. Assert complete private-state equality on every query path.
2. Exact supplied record for pending preparation: no worker, worker confirmed, failed and indeterminate variants remain accepted-outcome0. Visible bytes and worker outcome1 never create code1 settlement or alter projection.
3. Code1 accepted confirmed/failed/indeterminate outcomes: return exact outcome and record correspondence; failed/indeterminate with present bytes remain representable and never promote the projection. Include valid three-write and direct preparation-to-terminal histories.
4. Optional strict worker capture: correct kind5; kind7/kind3 substitution; mismatched exact request, reply, operation/generation/candidate/outcome/record digest; malformed binding/length/lane fields. Use independently well-formed wrong-request/outcome fixtures for semantic mismatch tests.
5. Same submitted request with accepted failed versus supplied worker indeterminate (and reverse), or accepted confirmed versus a valid failed worker reply, refuses exact-reply contradiction. Absence of worker remains unknown, not contradiction.
6. Exact current and earlier submitted records before/after later pending, failed/indeterminate and confirmed writes. Assert `LaterSubmitted` based only on entry order and preserve each entry's historical outcome. A newer visible uncertain candidate must coexist with earlier current confirmed projection without changing it; an older record must not roll it back.
7. Valid but never-submitted record substitutions: operation, generation, kind/candidate and other bound fields where constructible. Rejected well-formed code2 anchors must not become eligible submissions. Avoid tests that merely trip malformed framing before the membership assertion.
8. Successful duplicate/stale code2 never overwrites a prior accepted reply/outcome; wrong-pending code2 preserves pending entry with accepted-outcome0. Later exact settlement preserves fatal. Snapshot queries must not clear the latch.
9. Caller mutation after begin/reconcile and after query, returned value mutation, and repeated query success/refusal cannot change internal history/evidence. Compare complete private state, not projection alone.
10. Invalid begin/reconcile consumes neither entry nor evidence slot; same correct action afterward succeeds. Existing chronology, unique-ID/generation, one-pending and three-write/decision bounds remain. No imported or fabricated history path.

## Mutation evidence and restoration

Select only guards with reachable intended assertions. Useful targets: recorded-submission membership; optional worker kind/request linkage; worker/accepted reply equality; accepted-versus-worker outcome separation; code2 evidence preservation; LaterSubmitted classification; no projection/history mutation; and successful-trial evidence bookkeeping. Preserve valid preparation/setup when mutating later matching logic. For each mutant record the exact replacement, compiling command, intended assertion/subtest and observed failure location, then restore source bytes exactly and rerun the relevant restored test. Earlier G3a8 showed that a setup failure is not evidence for the intended later assertion.

Strict record/request equality, private reply validation and explicit confirmed digest equality overlap mathematically for valid canonical inputs. Removing one redundant comparison may survive because the other checks already force the same395 bytes/hash. Do not claim isolated necessity, fabricate hash collisions or corrupt private state merely to make a mutant fail. A carefully bounded combined bypass may demonstrate the behavioral boundary only if it preserves valid setup and reaches that exact assertion; otherwise disclose the overlap and cover it through exact independent expected-byte/digest tests plus composed checks. Similarly, the new wrapper's capacity is already bounded by G3a8 begin rules; do not claim independent sensitivity for a redundant capacity guard.

## Limits and handoff

This query classifies supplied correspondence, not an observed filesystem snapshot or full run verdict. Missing worker evidence, absent/current/older snapshot capture timing, complete worker-received/sent selections, complete trace causality, interleaved G3a8 action groups, malformed-reply diagnostics, process-observation-derived projection, stop/start admission, deadlines, real fsync/rename/custody, restart recovery and native execution stay deferred. No record loading, unknown-to-confirmed promotion, run/admission flag change or product import.

Confidence high for this bounded contract; no unresolved question requires a broader architecture choice. Brain must explicitly accept before behavior edits, retain this report, collect TDD/mutation/restoration evidence and fresh independent review, update canonical status/pins, and publish the shared destinations. This assignment ran no tests and makes no implementation-readiness claim.
