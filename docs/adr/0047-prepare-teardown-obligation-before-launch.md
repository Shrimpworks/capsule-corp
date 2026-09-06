# ADR-0047: Prepare an attempt-bound teardown obligation before launch

- Status: Proposed
- Date: 2026-09-06
- Refines if accepted: ADR-0011, ADR-0025, ADR-0041, ADR-0042, and ADR-0046

## Context

[C5b16](../C5B_TIMING_FAULT_CHECKPOINT.md) observed controlled 1,200/1,600-ms
publication delays exceeding the 1,200-ms total teardown bound, while post-gate
cleanup stayed within 1,000 ms. Publication refusal left the benign child live
without a Supervisor signal or native absence proof. The fixture self-alarm is
containment only. This result does not diagnose a product vulnerability or prove
installed behavior.

The current source also performs durable fencing and cursor lookup before the
teardown gate, then durable checkpoints before reap and absence reconciliation.
Moving one clock or one write cannot close that whole dependency chain.
[C5b17](../C5B_TEARDOWN_DEADLINE_DESIGN.md) records exact sources, anchors, alternatives,
failure cases and the next validation gate.

## Proposed decision

For a separately versioned successor, propose committing a closed cleanup
obligation with the consumed attempt before any process creation. The obligation
permits only termination and reconciliation of that attempt's exactly owned child;
it is neither a new execution grant nor a claim that termination has occurred.
No signal may precede confirmed durable preparation. Unknown commit outcome
permits no spawn or start.

The record must bind the installation and trust epoch, Supervisor identity,
AttemptID, consumed approval and registration, exact plan/profile/runner identities,
a Supervisor-generated operation identity, permitted teardown triggers, and exact
clock policy. Actual PID and process custody cannot be fabricated before creation.
Their later observation stays distinct from the precommitted obligation. Existing
record-before-start identity publication remains required before releasing start.

After creation, one Supervisor-owned control path must service cancellation,
deadlines, identity checks, signaling and authoritative reap/absence without
waiting for store I/O, store locks, completion/drain joins, or a pending storage
worker. A scheduling split inside the Supervisor is proposed; no separate helper,
daemon-to-backend route, additional authority or concurrent lifecycle owner is
selected. Exact scheduling, synchronization and platform custody mechanisms must
be specified and validated before runnable implementation.

At most one bounded storage operation may be outstanding. Its late response cannot
release start after cancellation, clear a stop/fence latch, restore authority,
overwrite newer terminal state, or release capacity. The owner retains its lock;
a new operation/attempt cannot bypass the pending worker. The control path may
continue only the already prepared destructive cleanup using live exact custody.
A storage error that suggests lost process identity or custody does not permit a
signal. State-store uncertainty alone never manufactures or destroys a separately
established process-custody fact.

The signal-attempt latch is monotonic and private to the current Supervisor
lifetime. A lost response never permits redriving the signal. Persisted preparation
is not a persisted assertion of signal issuance, exit, absence or cleanup. Outcome
publication follows independently observed facts; uncertainty remains explicit.
Completion and capacity release still require the full durable terminal join.

A restart restores obligation and uncertainty, not PID custody, fresh authority,
a new clock budget or permission to signal a reused PID. Without separately proven
restart/process-tree custody, mark recovery required, admit no replacement attempt,
and report no timely-absence or completion claim. Supervisor death and host
suspend/power failure remain unclosed installed-lifecycle gates.

The existing exact C1/C2A/C2B/C5a contracts, C5b11 driver, C5b14B/C5b16 storage formats,
24-provider ABI and accepted ADR lifecycles remain unchanged. A successor must use
new format/profile/driver identities; it must not reinterpret cursor 17 as proof
that the prepared obligation was discharged. No schema, API, dependency, service,
process or product consumer is introduced by this proposed ADR.

## Alternatives

| Candidate | Disposition and reason |
| --- | --- |
| Keep synchronous per-effect publication; move the timer earlier | Insufficient: observes a miss but cannot interrupt the storage wait. |
| Give each write a timeout | Insufficient without a proved cancellation/settlement mechanism; a timed-out write may still commit and hold shared locks. |
| Change storage engine or improve normal latency | May improve latency; does not alone prove behavior for stalled or indeterminate publication. |
| Signal first and record intent later | Rejected: removes durable-before-effect protection. |
| Guest self-alarm, daemon watchdog, or independent cleanup helper | Rejected for this slice: lacks required Supervisor authority/custody or expands the architecture. |
| Precommit the obligation and isolate the Supervisor control path | Proposed for passive validation; preserves prior durable intent while removing later publication from the stop dependency chain. |

## Consequences and acceptance gate

The proposed ordering can remove one source of unbounded waiting; it is not a
real-time guarantee. Identity checks, scheduler delay, kernel signaling, process
reap and descendant absence still require exact platform evidence.

Before accepting this ADR, require a reviewed passive successor contract/model
that distinguishes preparation, actual custody, stop latch, signal uncertainty,
physical absence and durable completion; specifies pending-write settlement;
reconciles every changed driver/cursor dependency; and passes the C5b17 failure
matrix. Require explicit maintainer acceptance of that concrete decision.

Only after that gate may a separately scoped benign fixture implement the selected
mechanism. It must demonstrate deadline independence under controlled blocked
publication and retain mutation sensitivity. Installed/guest/product promotion
requires its own evidence and authorization. No architecture remedy is activated
by proposing or merging this document while its status is Proposed.
