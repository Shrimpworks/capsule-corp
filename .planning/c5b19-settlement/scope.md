# G3a3 read-only scope recommendation

Work item: Define the next bounded pure C5b19 Gate-3 increment after G3a2.
Status: PASSED for scoping only.
Scope: Existing archive head `2274a6c7c751ded537541438f473e81ccd293457`, frozen Gate-2 packet, C5b18 bindings/pending/settlement model. Read-only source and graph inspection; no network, runtime fixture, subprocess test, signal, service, guest, store, credential or user data access.
Evidence: Exact source at the archive checkout and canonical worktree; material packet/model references below.
Remaining work: Brain selects/implements/tests/reviews/publishes the bounded increment. Full control model, raw-trace oracle and actual implementation/run gates remain separate.
Next action: Assign G3a3 as **pure one-attempt write-history and exact settlement conformance**, using existing tests/oracle.go and tests/model_test.go only for code/tests.
Parent status: IN_PROGRESS — TRENDING_GOOD. Gate4 first-run authorization and Gates5–6 remain BLOCKED.

## Why this increment

G3a2 already decodes request397/record395, validates reply98 correspondence, and checks structural candidates. Its wireOperation comment explicitly disclaims pending capability. decodeWireFields explicitly defers freshness/exact predecessor generation/custody. Stateless checks therefore permit individually valid messages whose sequence is impossible: overlapping writes, reused operation IDs, skipped generation, replayed applied reply, and contradictory durable candidate history. Closing these history-only relations is the smallest useful extension before implementing the much larger 29-event causal/timing oracle.

This is an independent in-memory projection of the **write lane only**. It is not the C5b18 full Model, a control state machine, a store, or a live Supervisor. It must never return MayCreate/MayStart/MaySignal/MayRelease decisions. A valid history remains insufficient evidence of preparation fsync, custody, authoritative absence, start permission, timely stop, public completion or capacity release.

## Suggested interface and owned files

Keep unexported package-local APIs in the already-frozen-listed `tests/oracle.go`; add tests in the existing `tests/model_test.go`, reusing literal G3a2 wire fixtures/helpers where suitable. Existing root planning/provenance records are lead-owned. No new source files, modules, C wrappers, fixture runners or testdata files needed.

Suggested names (implementation may choose equivalent narrow names):

- `newWriteHistory(frozen obligationBindings) (*writeHistory, error)`: validate and copy immutable bindings. No loader, restore/import or caller-supplied generation.
- `(*writeHistory).begin(request []byte) error`: structurally validate through existing decodeRequest, then accept one exact request into private owned storage. This records a submitted candidate, not permission to submit a real operation.
- `(*writeHistory).applyReply(reply []byte) error`: validate using existing validateReply against the **privately retained pending request**, then record one simulated control application. The caller must not supply replacement request/candidate/record bytes here.
- `(*writeHistory).snapshot() writeHistorySnapshot`: return value-owned pending operation, last settled operation+outcome, and confirmed projection/history fields. An explicit Pending boolean and LastSettledPresent are needed; do not make zero-byte outcomes mean failed or confirmed. Do not expose mutable slices/maps/pointers into internal state.

Private state can be strictly bounded: immutable bindings; one `[397]byte` pending request plus copied `wireOperation`; last settled tuple/outcome; `[68]byte` confirmed projection initially zero; maximum three consumed operation IDs; per-kind consumed flags or a small explicit phase. Fixed arrays fit this one-attempt protocol; an unbounded map/queue is unnecessary. No reset, reopen, retry, timeout, restore or capacity-release operation in this slice.

## Acceptance invariants

1. **Immutable binding:** constructor rejects all invalid frozen bindings already covered by obligation rules; a zero/uninitialized tracker must refuse begin/apply rather than accepting a caller-supplied first binding. Later requests always match the constructor-frozen set. Caller mutation after constructor/begin/snapshot cannot change pending identity or candidate.
2. **Single pending:** any accepted begin, including preparation, occupies the sole pending slot until an exact valid reply is explicitly applied. Merely constructing or validating a worker reply outside history does not clear it. A queued, executing or dropped-response operation all remain pending in this projection; no timeout implicitly settles it.
3. **Generation/identity:** initial accepted request has generation1; each later accepted request equals last settled generation+1 and remains within C5b18's `2^53-1` bound. Operation ID is nonzero and never reused, including after failed/indeterminate outcomes. Preparation uses frozen Preparation ID. Rejection must not consume an ID or advance a generation. Same request cannot be retried/reopened after settlement.
4. **Bounded write ordering:** preparation is first and at most once; only its confirmed application allows later writes in the projection. One runner-identity write is optional, at most once, and precedes terminal; terminal is at most once and closes further write submission regardless of its outcome. Thus valid histories are preparation alone, preparation→runner, preparation→terminal, preparation→runner→terminal, subject to outcomes. No fourth operation or second preparation/runner/terminal. This is the projection of required C5b18 ordering, not sufficient proof of its custody/absence guards.
5. **Exact settlement:** with no pending request, apply fails. Matching tuple alone is insufficient: outcome, request-body hash, candidate/record digest and fixed bindings must match the original frozen request through existing validation. Altered request, candidate, operation, generation, kind, binding or reply bytes fail and preserve the entire history. Duplicate applied reply and stale prior reply against a new pending operation fail. Only after successful validation clear pending and record the exact consumed candidate/outcome once.
6. **Confirmed projection:** confirmed outcome applies exactly pending.Candidate. Failed and indeterminate outcomes retain last settlement/outcome and advance history generation but leave confirmed projection unchanged. No record-loading API exists; record presence/hash alone cannot promote confirmation. A terminal confirmed candidate is teardown evidence only, never public completion.
7. **Candidate continuity:** preparation is the existing canonical prepared-only candidate. The runner candidate extends confirmed preparation with its nonzero runner identity. A terminal candidate preserves **exactly** whether runner identity was confirmed and its exact confirmed identity, then adds nonzero absence identity and terminal disposition. After failed/indeterminate runner settlement it must preserve runner-confirmed=false and zero runner identity, even though the failed pending candidate itself had runner fields. After confirmed runner settlement it must retain that runner and matching absence identity. Neither history rewind nor invented/erased confirmed runner facts is permitted.
8. **No false process claim:** if a runner request occurred, its frozen candidate describes the only modeled process identity, even if publication failed; a later terminal's absence identity must equal that prior attempted runner identity. This is byte-history consistency, not proof of custody. If the runner write was never attempted (early cancellation/absence), terminal can be valid with no confirmed runner and a nonzero synthetic absence identity; proving that identity against a real creation/wait observation belongs to the future control/trace oracle.
9. **No hidden control state:** late reply only changes these write facts. Do not introduce a pretend stop latch or claim late replies cannot reopen start from an always-false stub. Actual stop/start races, timing and authority need the future explicit control/event model and remain deferred.

Invariant8 is a necessary consequence of C5b18 BeginRunnerIdentityWrite freezing ProcessIdentity and BeginTerminalWrite freezing exact observed Absence.Identity, but the minimal history has no actual observation constructor. Document that assumption carefully and test both attempted-runner and skipped-runner paths. If the lead chooses a narrower envelope ledger without invariant4/7/8, explicitly rename its acceptance accordingly and retain those missing relations as deferred; do not call that the complete write-history conformance slice.

## TDD and mutation plan

Add failing behavior tests before history implementation; retain the red command/output excerpt in PLAN.md rather than only chat. Existing parser tests remain intact.

Positive tables:

- Preparation confirmed; runner confirmed; terminal confirmed for both absence-recorded and timing-violated dispositions. Check every pending/settled/confirmed snapshot and exact gen1/2/3 progression.
- Preparation confirmed then terminal directly (no runner request) with runner fields zero.
- Preparation confirmed; runner failed and indeterminate in separate cases; terminal with zero runner fields and the matching attempted runner identity as absence. Terminal may settle confirmed, failed or indeterminate; only confirmed changes projection.
- Failed/indeterminate preparation: last settlement retained, no confirmed projection and no further begin.
- A simulated indefinitely pending operation remains pending through repeated snapshots and standalone stateless reply validation; no actual sleep, clock, goroutine or store delay.
- Multiple independently changed valid bindings/operation IDs/process identities (not only one hardcoded fixture); independent request/reply bytes computed using existing test helpers, not the new history implementation.

Negative tables:

- Zero/uninitialized tracker begin/apply; invalid constructor bindings; runner or terminal first; starting generation0/2/max; skip/repeat/decrease generation; reused operation ID after each outcome; same-kind retry with fresh ID and correct next generation; fourth request; runner after terminal.
- Overlapping preparation/runner/terminal and rejected begin while pending; follow each rejection with a valid apply to show pending wasn't lost.
- Apply without pending; duplicate apply; previous reply during a new pending request; wrong operation/generation/kind/request hash/candidate/digest/outcome; cross-attempt reply; malformed and trailing/short bytes. Check state byte-for-byte unchanged after every rejection.
- Failed/indeterminate runner candidate must not enter confirmed projection; terminal cannot invent a confirmed runner, erase a confirmed runner, substitute runner or absence identity, or resurrect a failed candidate as confirmed history.
- Caller mutations to constructor input, accepted raw request and returned snapshots must not affect later matching reply/settlement; copied arrays are preferable to retention of slice aliases.

Suggested four controlled local mutants, applied one at a time and restored exactly:

1. Bypass pending guard: overlap test must fail.
2. Bypass generation equality or fresh-ID guard: use dedicated matching otherwise-valid history test so structural parsing is not the accidental rejection reason.
3. Apply candidate for failed/indeterminate outcome: projection test must fail.
4. Omit terminal continuation identity check (especially failed runner identity→wrong terminal absence): dedicated state-consistent structural vector must catch it.

Also ensure a guard mutant replacing pending request at apply cannot pass stale/candidate-swap test. Selected mutation tests are sensitivity evidence only, not exhaustive proof.

Offline commands: exact existing `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go`, with isolated GOCACHE, `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off`; test, race, vet, build, diff check and restored test. Test code itself must not spawn a process or touch filesystem/network/clocks. Normal Go tool compilation/unit-test execution is not lifecycle-fixture execution. Fresh-context review is for this new G3a3 delta plus integration, max3, preserving all earlier G3a/G3a2 review records.

## Precise sources

- Gate-2 fixed transport/identity lines187–250: only three kinds; exact request/reply schemas and hashes; all C5b18 binding/generation/candidate rules; logical pending covers queued, executing and unconsumed reply; exact predecessor settlement and terminal-after-absence.
- Gate-2 durable format lines253–287: immutable candidate, confirmed/failed/unknown distinction; record presence alone never confirmation; no delayed/lost reply creates another operation/generation or changes candidate.
- Gate-2 trace/oracle lines338–360 and446–453: worker reply is not control application; kind8 applies only exact validated frozen operation; pre-run tests use simulated process/store effects.
- C5b18 types.go207–237: candidate and exact PendingWrite/SettledWrite value shapes.
- C5b18 model.go67–82: preparation identity/one shot;104–125 runner one shot and frozen identity;130–153 terminal candidate/one shot;155–181 single pending, used IDs and LastSettled.Generation+1;185–211 exact pending settlement, uncertain outcome preserving durable projection, terminal recovery-required;233–245 exact whole-token equality and retained settlement.
- C5b18 model.go Decision and control.go lifecycle guards remain outside this projection; do not infer their behavior from wire success.
- Archive tests/oracle.go123–260 provides all needed stateless decoding/correspondence primitives at head2274a6c. It currently imports only bytes/sha256/binary/errors; preserving those imports is feasible.

## Dependency/custom-primitive decision and boundaries

Applicable row: `C5b19 prospective Supervisor control/storage lane exchange` in ECOSYSTEM_REUSE_AND_ADOPTION.md. Recommendation remains ADOPT-PLATFORM / BUILD-NARROWLY; this slice uses only BUILD-NARROWLY fixed bounded in-memory history and already used standard-library primitives. No native platform API, dependency, go.sum, transport bus, queue, thread or runtime module is added.

Checklist record in PLAN.md should identify: test/evidence trust class; Go1.25.13 and standard-library-only graph/licensing/provenance inherited from ORIGINS; fixed one attempt/three operations/one397-byte pending request/98-byte reply/68-byte projection caps; no authority beyond synthetic-byte validation; no process/network/credential/user filesystem/storage effects; atomic validation-before-state-mutation and explicit uncertainty; exact offline tool reproduction; positive/negative/restoration/mutation evidence; Capsule maintainer upgrade/removal ownership. Exact packet/toolchain/host/ABI remains unchanged. No existing product implementation or ADR responsibility changes.

Unresolved boundary to state before coding: history.applyReply is **a simulated accepted control application**, not merely worker reply arrival; absent actual trace/event pairing the resulting confirmation remains an offline claimed write fact, not real durable evidence. Do not use a name such as verifyExecution or completeAttempt. No new authorization is needed for this purely inert subset of the already-authorized packet; child execution stays separately gated.

## Confidence and verification limits

High confidence in the proposed necessary write-history relations. No independent review of future implementation or proof of live storage/lifecycle/timing behavior. No tests run in this read-only assignment.

Codebase-memory skill already read and applied. list_projects confirmed Capsule project; exact-path and teardownpassive scope coverage reported metadata_match/no_recorded_issue at full generation2026-10-04T22:18:10Z, best-effort only. Targeted six-symbol search was fully paginated (6returned/no more); SettleWrite snippet read. Current canonical types/model/packet were also read directly. Archive source is not established as indexed; direct exact-head reads used and git status was clean, HEAD2274a6c. No branch, commit or PR created. Sole retained file is this /tmp report for parent reconciliation.
