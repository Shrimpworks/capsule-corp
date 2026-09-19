# C5b18 passive teardown successor specification

Date: 2026-09-18

Status: review-instance-1 corrections `IN_PROGRESS — TRENDING_GOOD`; fresh-context
review instance 2 of 3 and refreshed full verification pending.
Parent owner-only hostile-`.mjs` internal alpha: `IN_PROGRESS — TRENDING_GOOD`.
Runnable successor, installed lifecycle, guest execution and product admission: `BLOCKED`.
ADR-0047 lifecycle: **Proposed**, not Accepted.

## Objective

Freeze the first passive successor contract required by
[ADR-0047](adr/0047-prepare-teardown-obligation-before-launch.md). The model must
show that a confirmed attempt-bound cleanup obligation permits exact destructive
cleanup to progress independently of later storage work without turning that
obligation into process custody, start authority, signal evidence, absence or
completion.

Success is a deterministic no-effect model with executable failure, restoration
and mutation cases. It does not select a scheduler, storage engine, lock, native
process API, installed service or runnable composition.

## Contract identity and project structure

- Contract: `capsule.c5b18.teardown-successor-model/v1`.
- Record: `capsule.c5b18.teardown-obligation/v1`.
- Go package: `internal/execution/teardownpassive`.
- Model tests remain in that package and exercise no filesystem, process, network,
  clock, credential, backend, VM or guest.
- Canonical specification and status remain in `docs/`; no product package may
  import the passive model.

The new identities do not reinterpret existing C5b11/C5b14B/C5b16 driver,
snapshot or cursor values. Cursor 17 remains historical recovery evidence only.

## Closed bindings

One immutable obligation binds these domain-separated, nonzero values:

| Field | Width | Authority meaning |
| --- | ---: | --- |
| `InstallationID` | 16 bytes | Enrolled local installation only |
| `TrustEpochDigest` | 32 bytes | Exact local trust epoch |
| `SupervisorIdentityDigest` | 32 bytes | Exact expected Supervisor identity |
| `AttemptID` | 16 bytes | Consumed attempt only |
| `ApprovalID` | 16 bytes | Consumed approval only |
| `RegistrationID` | 16 bytes | Registered plan only |
| `PlanDigest` | 32 bytes | Exact plan bytes |
| `ProfileDigest` | 32 bytes | Exact runtime/profile bytes |
| `RunnerDigest` | 32 bytes | Exact runner candidate bytes |
| `PreparationOperationID` | 16 bytes | Supervisor-generated preparation operation |

Same-width values remain distinct Go types and cannot substitute for one another.
The obligation permits only cancellation, wall-deadline and Supervisor-detected
fatal-fault teardown triggers for the bound attempt. It grants no create, start,
observe, result-release or replacement-attempt authority.

## Exact clock policy

The v1 policy is fixed:

- monotonic ticks only within one Supervisor lifetime;
- wall action at `t_start + 1,000 ms`;
- no graceful-teardown extension in this successor;
- forced absence no later than `t_force + 1,000 ms`;
- absence no later than `t_action + 1,200 ms`;
- `t_action` is the earliest accepted cancellation, wall or fatal-fault anchor;
- repeated triggers never move `t_action` later;
- a pending or late storage response never resets a clock;
- restart invalidates live ticks and supplies no fresh budget.

The model reports only whether retained same-lifetime observations satisfy these
relations. It proves no scheduler, kernel, suspend/wake or platform latency.

## State separation

The model contains three separate projections.

### Durable projection

Durable facts are limited to:

- the immutable obligation and its confirmed preparation outcome;
- the monotonic consumed-creation fact after first exact custody observation;
- confirmed runner-identity publication, when any;
- a confirmed terminal record and its closed completed or timing-violated disposition
  only after authoritative absence;
- the latest settled operation ID, kind, generation and frozen candidate; and
- explicit failed or indeterminate storage outcome.

Preparation never stores a PID, process identity, signal, exit, absence or cleanup
claim. Only the completed terminal disposition permits public completion, output
release and capacity release. A confirmed timing-violated record remains unresolved
and recovery-required.

### Current-lifetime control projection

Volatile facts are limited to:

- exact live custody observed in the current Supervisor lifetime;
- monotonic stop latch and earliest trigger/anchor;
- first signal request and any uncertain response;
- authoritative absence observation; and
- same-lifetime clock observations.

Restart clears custody and every live clock/signal/absence fact. It restores the
durable obligation as conservative may-exist state and sets recovery required. V1
has no operation that adopts a PID or restores custody after restart. Consumed
creation never reopens after natural absence or restart.

### Pending storage projection

At most one bounded storage operation exists. It binds:

- a nonzero operation ID distinct from preparation identity;
- a strictly increasing generation;
- one closed kind: preparation, runner identity or terminal join; and
- the exact candidate durable projection frozen when the operation began.

No second operation, attempt or generation may bypass it. Settlement must match
the exact operation ID, generation, kind and frozen candidate. Every settlement,
including failed and indeterminate outcomes, retains those fields plus outcome.
A failed or indeterminate outcome does not silently retry. Late success may retain
its durable fact but cannot clear a stop latch, release start, replace first
observations or release capacity.

## Transition contract

1. Begin preparation before any creation observation.
2. Refused or indeterminate preparation permits no create or start. Reopen may
   classify the complete old/new state only.
3. Confirmed preparation permits one current-lifetime exact custody observation.
   That observation consumes creation authority permanently for the attempt;
   natural absence cannot authorize replacement custody.
4. Runner identity publication begins after custody and holds start closed.
5. Cancellation, wall or fatal fault may latch stop while any post-create write is
   pending. Stop evaluation never depends on settlement of that write. No trigger
   may latch or rewrite timing after authoritative absence.
6. Start may be attempted once only after confirmed runner identity, with exact
   current-lifetime custody and no stop latch. `t_start` is captured immediately
   before that first attempt and never reset.
7. Signal may be requested once only after stop is latched and the caller supplies
   the exact current-lifetime custody identity. Zero, stale or substituted identity
   refuses without mutation. A lost/uncertain response prohibits redrive.
8. Authoritative absence may be recorded only with that same exact identity.
   Signal return, EOF, root removal and fixture alarm are not absence.
9. Terminal publication may begin only after absence and freezes one closed
   disposition. `completed` permits durable completion, output and capacity release.
   `timing-violated` records failure evidence, retains recovery required and permits
   none of those releases.
10. Restart preserves durable preparation and storage uncertainty, clears live
    authority, marks recovery required and permits no PID adoption, signal,
    start, replacement attempt or completion.

Every refusal leaves the prior model state unchanged.

## Failure and restoration acceptance matrix

| Case | Required model result |
| --- | --- |
| Preparation pending, failed or indeterminate | create/start refused; no retry identity invented |
| Preparation confirmed; crash before creation | obligation retained; recovery required; no custody/signal/completion |
| Runner publication pending/refused | start closed; stop still serviceable with exact custody |
| Stop while a post-create write is pending | stop and first signal request progress without write settlement |
| Late runner-publication success after stop | durable identity may settle; start stays closed; latch unchanged |
| Storage failure with exact live custody | cleanup remains available; completion/capacity stay withheld |
| Signal response uncertain | second signal refused; exact custody may still reconcile absence |
| Absence observed; terminal write fails | local absence retained; public completion/output/capacity withheld |
| Natural absence before stop | creation remains consumed; no replacement custody or late trigger |
| Timing bound violated, later absence observed | timing failure remains unresolved; terminal evidence may persist but completion/output/capacity stay withheld |
| Restart, lost reaper, PID reuse or mismatch | recovery required; exact identity mismatch mutates nothing; no custody adoption or signal |
| Repeated cancel/wall/fatal race | earliest anchor retained; one stop latch and at most one signal |

## Test strategy

Tests use a pure deterministic model and fixed byte values. They must cover every
matrix row plus exact/capacity/domain validation and defensive copies.

Required mutation sensitivity:

1. restoring a storage-settlement dependency to stop;
2. permitting start after a stop latch;
3. permitting a repeated signal after uncertain response;
4. forgetting confirmed preparation on restart;
5. resetting the earliest deadline on a later trigger;
6. adopting PID/process custody after restart; and
7. beginning or committing terminal state before authoritative absence;
8. reopening creation after natural absence;
9. accepting a zero, stale or substituted signal/absence identity;
10. releasing completion, output or capacity after a timing violation; and
11. settling a mutated frozen candidate or losing its exact settled operation fields.

Each mutation must fail a named assertion rather than time out or fail to compile.
The independent decision oracle derives `mayCreate`, `mayStart`, `maySignal`,
`mayPublishTerminal`, `mayRelease` and `recoveryRequired` only from the public
snapshot; it does not call the model's transition guards.

## Commands

Focused test:

```sh
go test ./internal/execution/teardownpassive
```

Repository verification follows `AGENTS.md`: pnpm install/check/lint/test/schema/ADR
verification, Go test/vet/build/lint and pinned `govulncheck`.

The corrected head requires refreshed full results before review instance 2 begins.

## Independent review loop

Review instance 1 of 3 inspected exact head `26fa5bc` and returned **Not ready**.
The implementation task accepts all four findings:

1. natural absence reopened creation for one consumed attempt;
2. timing violation could publish success and release output/capacity;
3. signal and absence did not compare exact process identity; and
4. pending/settled storage projections omitted required candidate/operation fields.

Commits `dacb95f`, `bffc71e`, `39cfe7b`, and `b459609` correct those findings;
`5307e81` additionally prevents a post-absence trigger from reclassifying the first
terminal observation. These corrections materially change authority and terminal
state, so review instance 2 of 3 is mandatory after refreshed verification. No
review verdict is claimed for the corrected head yet.

## Boundaries

Always preserve Supervisor-only lifecycle ownership, durable-before-effect
ordering, one-use attempt binding, no-redrive uncertainty and completion-last
release.

Any new dependency, store engine, helper, service, product consumer, schema/API,
process effect, installed evidence, guest execution or ADR acceptance requires a
separate reviewed task and authorization.

Never treat this model as platform evidence, runtime/profile admission, a timing
guarantee, restored PID custody, or proof that cleanup occurred.

## Success criteria

- exact v1 bindings, clock policy and three state projections compile as a passive
  no-effect model;
- every failure/restoration row and all eleven mutations are executable and pass;
- no product consumer imports the package;
- ADR-0047 remains Proposed with the concrete packet linked for maintainer review;
- canonical status documents distinguish scoped `PASSED` from blocked runnable,
  installed, guest and product work; and
- full required verification and independent review pass before merge.
