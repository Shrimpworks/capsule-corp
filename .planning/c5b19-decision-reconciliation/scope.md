# C5b19 G3a8 pending-write decision reconciliation scope

Status: PASSED for read-only scope derivation. Parent Gate 3 remains IN_PROGRESS — TRENDING_GOOD; child execution, Gates 4–6 and product admission remain BLOCKED.

Assignment: internal read-only G3a8-SCOPE. Defensive intent: validate the offline distinction between exact pending-write settlement and rejection, using owned inert synthetic request/reply/event bytes only. No fixture, child, guest, lifecycle, signal/wait, store, live clock or network operation was performed. No repository files were changed.

## Source and method

Canonical base: `8e4fd22cc00cd5e961544f6a803974600dc0e751` in `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-decision-reconciliation`.
Archive base: `5d0e53cbccae3d4f587205df70253b776eae63e5` in `/tmp/capsule-c5b19-recovery`, branch `codex/c5b19-benign-teardown`. Source below is within `experiments/typed-guest-transport-c5b19-benign-teardown/tests/oracle.go`.

Read AGENTS, relevant starter-document boundaries, ADR-0047, reuse-table C5b19 row, frozen Gate2 packet, normative passive model/types and prior G3a7 planning. Used codebase-memory Tier 2: primary index ready at `2026-10-04T22:18:10Z`, not watched; packet/types/model/ADR coverage metadata_match. SettleWrite/samePending lookup returned two complete results; direct current source checks establish the material claims. Archive and unindexed planning use direct source fallback. This is bounded evidence, not a whole-graph completeness claim.

Grounding:

- Packet lines 242–248: one logical pending operation until its exact predecessor settles; no timing-based implicit settlement.
- Packet lines 342–360: kind8 code1 follows exact tuple/request/candidate validation plus C5b18-equivalent settlement; code2 leaves durable projection unchanged and latches fatal. Captured requests/replies are bounded to three each.
- Packet lines 367–382, 397–410: kind ownership, private control tick watermark, no worker clock authority.
- Packet lines 430–453: full stop/apply ordering, causal/lifecycle assertions and native first-run gates remain later work.
- `internal/execution/teardownpassive/model.go:185–211,233–235`: SettleWrite checks outcome and exact active pending before projection; late confirmation changes its durable fact without changing any control latch. Runner begin separately rejects stop; terminal begin uses absence/control permission. These are not derivable from request/reply bytes alone.
- Existing oracle `writeHistory` lines 260–390: bounded private copied pending request and operation history; accepted reply settlement only; begin explicitly records a submitted candidate and is not permission to submit.
- Existing oracle G3a7 lines 707–792: individual pair validates strict represented correspondence, not recorded decision correctness. Collection additionally forbids duplicate operation IDs and should remain unchanged.
- Reuse table selects ADOPT-PLATFORM / BUILD-NARROWLY for prospective C5b19 lanes, still test-only. This slice composes existing local primitives; no new dependency, transport, packet, privileged responsibility or product admission is warranted.

## Recommended minimal contract

Add a separate bounded pure reconciler around an internally owned **value copy** of writeHistory, immutable frozen bindings/timebase, a copied snapshot, monotonic modeled-fatal flag, decision count, and a private control sequence/tick watermark. Constructor validates bindings and nonzero timebase even before any action. No reset, borrowed mutable history, external pending setter, timeout or record-loading path.

`begin` accepts an existing capturedFramePair with strict kind3 event/request correspondence. Require its sequence strictly greater than the preceding accepted watermark and tick at least the preceding tick. Delegate to a trial copy of existing history.begin, then commit all state together. Preserve its generation, unique-ID, at-most-three-writes, single-pending and candidate rules. Begin is history recording only, explicitly **not control admission**. Do not introduce a fatal/stop-based begin predicate: actual permission requires additional control state not represented here.

`reconcile` accepts one existing receivedAppliedPair, uses existing **individual** G3a7 validation, requires receive sequence strictly after the preceding watermark and receive tick at least the prior tick, and caps accepted represented decisions at three. The pair validator already enforces receive before apply and nondecreasing within-pair tick. Compare recorded application code with independently derived expected code before committing.

| Validated pair relative to private history | Expected code | Successful modeled effect |
| --- | --- | --- |
| Pending exists, entire 397-byte anchor equals private pendingRequest, trial applyReply succeeds | 1 | Commit existing history settlement; confirmed outcome alone updates projection |
| No pending | 2 | Preserve complete history; latch fatal |
| Well-formed anchor/reply belongs to another operation, generation or candidate | 2 | Preserve complete history including current pending; latch fatal |

An exact anchor and correctly validated reply should make trial settlement succeed under current invariants. Derive code1 from both exact byte equality and successful private-pending applyReply, rather than assuming tuple/digest equality proves it. A structurally valid anchor may describe a request never begun in this history: it remains a decode context, not proof of submission.

Successful code1 for failed/indeterminate outcome still clears pending and advances LastSettled, but leaves Confirmed unchanged. Code2 is a successful representation of rejection, not a validator error. It leaves **all** history fields unchanged, including hidden used IDs, attempted kinds, runner attempt and copied pending bytes. It advances decision count/watermark and sets fatal true. Fatal never clears, including when a later exact pending reply is accepted.

Malformed/unbound/mismatched event-frame evidence, chronology failure, count overflow or wrong recorded application code returns an error with **no mutation anywhere**, including fatal, counts and watermark. In particular, a recorded code2 for exact settleable pending is false rejection and must refuse; code1 for no/mismatching pending is false acceptance and must refuse. Invalid evidence is not represented rejection and must not itself simulate fatal.

## Ordering, late replies and limits

- Copy the complete reconciler to trial state; validate and compute effects there; replace original only on success. No caller slices or mutable external history are retained. Return value snapshots.
- At most three begins (existing history bound) and three represented receive/apply decisions, including code2. Rejected decisions do not consume operation IDs or generations.
- A later duplicate/stale reply with fresh valid event ordinals may correctly produce code2. Replaying the same event ordinals is invalid evidence and errors atomically. Do not call the G3a7 collection validator here: its independent duplicate-operation rule would preclude the intentional duplicate/stale rejection case.
- Strictly increasing supplied control sequence and nondecreasing tick span kind3, kind7 and kind8. First sequence/tick may be zero; gaps and equal ticks are valid. Compare values without incrementing sequence/tick, so uint64 boundaries cannot wrap. Worker ticks are irrelevant. No elapsed deadline inference.
- A whole-pair API supports serial supplied action groups: receive/apply are processed together. It cannot reconcile a begin interleaved between a pair's receive and apply. Document this unsupported selection ordering; do not infer the corresponding full trace is invalid. No trace completeness or adjacency claim.
- Fatal is only the packet's monotonic rejection-derived modeled latch. It supplies no permission or cancellation behavior. A later valid reply for still-pending bytes settles its write fact and preserves fatal. External stop is unrepresented; therefore this slice neither rejects a late reply solely for stop nor establishes stop/start/signal behavior.
- Existing history-only begin rules may accept later inert history after fatal; that acceptance is not admission to actually submit/create/start/publish. Enforcing real post-fatal/post-stop admission is deferred rather than guessed.
- Malformed diagnostic code2 events remain unsupported by strict G3a7 framing. A corrupt/unbound raw reply cannot be dressed as a valid rejection in this slice. Proper diagnostic mapping needs its own scope.
- No full oracle, omitted-event detection, worker receipt, durable snapshot/readback, actual fsync, full control state, stop/apply causal decision, native effect, child execution or run/admission flag changes.

## Required test matrix and mutation sensitivity

1. Valid three-write history with independently asserted snapshots/projections; each of confirmed, failed and indeterminate settlement effects; direct preparation-to-terminal history remains representable under existing rules.
2. False code1 with no pending and wrong pending anchor; false code2 for exact pending. Assert entire private state unchanged, then reuse the same correct ordinals to prove refusal did not consume watermark/count.
3. Correct no-pending code2; correct wrong-pending code2 preserves exact pending; late valid pending code1 settles and preserves fatal. Failed/indeterminate late settlements also preserve fatal and confirmed projection.
4. Wrong ID, generation, kind/candidate and request-anchor substitutions must use independently well-formed matching request/reply/event fixtures where possible. They produce expected2, rather than merely tripping an earlier malformed-byte validator.
5. Structurally invalid input, wrong outcome/record digest, wrong pair member, malformed application code and nil/uninitialized instance all refuse without state change.
6. Duplicate reply after settlement and stale reply while another write is pending: fresh ordinals code2 succeeds; false code1 refuses; identical-ordinal replay errors. Reused ID/generation on begin still refuses via history.
7. Cross-call control sequence/tick errors covering begin→receive, apply→begin, apply→receive; duplicate ordinals, decreasing ticks, accepted same ticks, zero origin, gaps and max-value boundaries.
8. Fourth decision refuses atomically, including rejection-only histories; begin count/kind/terminal and one-pending bounds retained. Invalid attempts do not consume budget.
9. Caller mutates request/reply/event buffers after acceptance and snapshot values after retrieval: internal history/fatal/watermark unaffected. Whole-state equality checks cover hidden fields, not only exported projection.
10. Explicit history-only begin-after-fatal demonstration, named as absence of control admission; no fabricated stop-state test or claim.

Selected compiling mutants should remove exact pending anchor equality, decision equality, history commit, rejection history preservation, fatal latch/preservation, cap, cross-call sequence and cross-call tick guards. Each must fail its intended behavioral assertion; restore exact source bytes afterward. Tests must verify independently expected byte projections, not only echo implementation output.

## Handoff

Confidence high for bounded selection/settlement contract; full control admission and interleaved trace processing are deliberately outside scope. No unresolved requirement blocks this bounded implementation if those limits are retained. Brain owns implementation, tests-first red/green evidence, independent review, canonical roadmap/evidence update, archive commit/PR destination and canonical PR. Verification here consisted solely of graph/source/document reads; no tests or execution were run by this assignment.
