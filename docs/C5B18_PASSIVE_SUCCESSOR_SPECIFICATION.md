# C5b18 passive teardown successor specification

Date: 2026-09-18

Status: first review in the human-authorized successor cycle returned **Not ready**
at `0aa6cbc`: an observed service tick could be backdated by a later trigger.
The causal-order correction, local verification and self-review are `PASSED`;
independent review of the corrected head remains pending.
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
- Model transitions exercise no filesystem, process, network, real clock,
  credential, backend, VM or guest. One structural test reads local repository Go
  files to prove no product consumer imports the package.
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
- wall action anchored at `t_start + 1,000 ms`; retain separately the observed
  callback-service tick, never substitute the scheduled anchor for service;
- no graceful-teardown extension in this successor;
- forced absence no later than `t_force + 1,000 ms`;
- absence no later than `t_action + 1,200 ms`;
- `t_action` is the earliest accepted cancellation, wall or fatal-fault anchor;
- repeated triggers never move `t_action` later;
- same-lifetime service, start, signal and absence observations use a nondecreasing
  observed-tick watermark; a delayed cancellation or fault presented without a
  separate current service tick is refused if its tick predates an observation
  already retained, while a wall action anchor may predate its actual service;
- service later than the mandatory wall anchor is a timing violation, even if
  signal and absence follow promptly;
- a pending or late storage response never resets a clock;
- restart invalidates live ticks and supplies no fresh budget.

The model reports only whether retained same-lifetime observations satisfy these
relations. It proves no scheduler, kernel, suspend/wake or platform latency.

## State separation

The model contains three separate projections.

### Durable projection

Candidate durable facts are limited to:

- the immutable obligation and its confirmed preparation outcome;
- confirmed runner-identity publication with exact process identity, when any;
- a confirmed teardown/absence record bound to exact observed process identity,
  including when runner publication failed, and its closed absence-recorded or
  timing-violated disposition
  only after authoritative absence;
- the latest settled operation ID, kind, generation and frozen candidate; and
- explicit failed or indeterminate storage outcome.

Preparation never stores a PID, process identity, signal, exit, absence or cleanup
claim. Runner publication freezes exact current-lifetime process identity. A
teardown/absence record is never a full job terminal join: this model lacks typed
result, lifecycle, cleanup and completion-last proof and therefore never grants
public completion, output or capacity release. Both record dispositions remain
recovery-required.

### Current-lifetime control projection

Volatile facts are limited to:

- one-use creation consumption observed after exact custody; this model has no
  durable publication for that consumption;
- exact live custody observed in the current Supervisor lifetime;
- monotonic stop latch, earliest trigger/anchor, and separate service observation;
- latest accepted observed tick, even when a later trigger does not replace the
  earliest stop anchor;
- first signal request and any uncertain response;
- authoritative absence observation with exact custody identity; and
- same-lifetime clock observations.

Restart clears custody and every live clock/signal/absence fact. Even if creation
was not observed before restart, confirmed preparation is treated conservatively
as may-exist: creation is closed and recovery required. V1
has no operation that adopts a PID or restores custody after restart. Consumed
creation never reopens after natural absence or restart.

### Pending storage projection

At most one bounded storage operation exists. It binds:

- exact immutable installation/epoch/Supervisor/attempt/approval/registration/
  plan/profile/runner/preparation bindings plus a nonzero operation ID;
- a strictly increasing generation;
- one closed kind: preparation, runner identity or terminal join; and
- the exact candidate durable projection frozen when the operation began,
  including process identity for a runner-publication candidate.

No second operation, attempt or generation may bypass it. Settlement must match
the exact bindings, operation ID, generation, kind and frozen candidate. Every settlement,
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
   pending. Stop evaluation never depends on settlement of that write. Wall
   latching requires an actual service tick distinct from its fixed start-derived
   anchor; a late service remains a timing violation. A trigger with a service
   tick older than a retained observation is refused without mutation; the
   earliest action anchor is not an excuse to backdate later signal or absence.
   No trigger may latch or rewrite timing after authoritative absence.
6. Start may be attempted once only after confirmed runner identity, with exact
   current-lifetime custody and no stop latch. `t_start` is captured immediately
   before that first attempt and never reset.
7. Signal may be requested once only after stop is latched and the caller supplies
   the exact current-lifetime custody identity. Zero, stale or substituted identity
   refuses without mutation. A lost/uncertain response prohibits redrive.
8. Authoritative absence may be recorded only with that same exact identity, and
   terminal evidence freezes it even if runner identity publication failed.
   Signal return, EOF, root removal and fixture alarm are not absence.
9. Teardown/absence-record publication may begin only after authoritative absence
   and freezes one closed disposition. Neither `absence-recorded` nor
   `timing-violated` permits public completion, output or capacity release. A
   started attempt missing its mandatory wall action at the bound is timing-violated.
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
| Started attempt missing mandatory wall action | late absence is timing-violated; no public release |
| Wall callback serviced after its fixed anchor | retain actual service tick; late action is timing-violated |
| Runner publication failed before exact absence | teardown record retains observed process identity without forging runner success |
| Signal response uncertain | second signal refused; exact custody may still reconcile absence |
| Absence observed; terminal write fails | local absence retained; public completion/output/capacity withheld |
| Natural absence before stop | creation remains consumed; no replacement custody or late trigger |
| Timing bound violated, later absence observed | timing failure remains unresolved; terminal evidence may persist but completion/output/capacity stay withheld |
| Restart, lost reaper, PID reuse or mismatch | recovery required; exact identity mismatch mutates nothing; no custody adoption or signal |
| Repeated cancel/wall/fatal race | earliest anchor retained; one stop latch and at most one signal |
| Earlier cancellation delivered after observed wall service | Backdated trigger refused; signal/absence before that service refused without mutation |
| Later trigger does not replace earliest anchor | Latest observed tick still fences signal/absence timestamps |

## Test strategy

Tests use a pure deterministic model and fixed byte values. They must cover every
matrix row plus exact/capacity/domain validation and defensive copies.

Guard scenarios corresponding to future mutation targets:

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
11. settling a mutated frozen candidate or losing its exact settled operation fields; and
12. backdating a trigger, signal or absence behind any retained same-lifetime observation.

The named Go tests assert these guards on the unmodified model. They are not an
executed mutant campaign; actual mutation sensitivity remains unverified and is
required before a mutation-proof claim.
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

Observed on the review-3 correction worktree on 2026-09-18:

- `pnpm install --frozen-lockfile`, `pnpm check`, `pnpm lint`, `pnpm test`,
  `pnpm verify:schemas`, `pnpm verify:adrs`, `pnpm audit:dependencies` and
  `pnpm site:build` passed;
- `go test ./...`, `go vet ./...`, `go build ./...`, package race/coverage
  (`86.3%` statements), formatting, CI security/correctness lint, new-code
  `revive` lint and pinned `govulncheck@v1.6.0` under the declared Go 1.25.13
  toolchain passed; and
- unrestricted `golangci-lint run ./...` reports only the 50 pre-existing
  exported-comment findings tracked in issue #217. No C5b18 finding remains.

For the causal-order correction, the new late-wall regression failed before the
fix and passes after it. Refreshed `pnpm install --frozen-lockfile`, `pnpm check`,
`pnpm lint`, `pnpm test`, schema/ADR verification, dependency audit and static-site
build passed. Go test/vet/build, package race/coverage (87.5%), blocking lint,
package `revive` and pinned `govulncheck@v1.6.0` under Go 1.25.13 passed. Full
`golangci-lint run ./...` still reports only the 50 pre-existing issue-#217
`revive` findings. An independent verdict on the corrected head is not yet claimed.

Review instance 3 returned **Not ready** at `cbaf84d`; its two corrections,
refreshed local verification and CI pass at `567c733`. A human-authorized fresh
cycle began at `0aa6cbc`; its first reviewer returned **Not ready** on clock
causality. The current correction rejects observations older than the latest
accepted service tick, including when an earlier stop anchor remains selected.
No independent verdict is claimed for the correction head. The
host's alternate Go 1.26.5 is vulnerable to three standard-library
advisories; it is not the declared build toolchain. Do not use it for this candidate.

## Independent review loop

Review instance 1 of 3 inspected exact head `26fa5bc` and returned **Not ready**.
The implementation task accepts all four findings:

1. natural absence reopened creation for one consumed attempt;
2. timing violation could publish success and release output/capacity;
3. signal and absence did not compare exact process identity; and
4. pending/settled storage projections omitted required candidate/operation fields.

Commits `dacb95f`, `bffc71e`, `39cfe7b`, and `b459609` correct those findings;
`5307e81` additionally prevents a post-absence trigger from reclassifying the first
terminal observation. These corrections materially changed authority and terminal
state. Review instance 2 examined `3dc70e6` and found absence-only public release,
cross-attempt settlement-token substitution, an unpersisted creation fact called
durable, and an unsupported executed-mutation claim. A second blind review of the
same head also found a missed wall-action timing violation and a documented record
identity absent from model state. All findings are accepted for correction.

Review instance 3 inspected `cbaf84d` and found that a scheduled wall anchor
could masquerade as actual late callback service, and terminal evidence after
failed runner publication omitted exact observed absence identity. Both findings
are accepted; no review verdict is claimed for the subsequent correction head.

The first review in the explicitly human-authorized successor cycle inspected
`0aa6cbc` and returned **Not ready**. It reproduced wall service at tick 1,100,
then accepted cancellation at 250, signal at 260 and absence at 300, falsely
classifying timing as satisfied. The finding is accepted. The correction uses a
same-lifetime observation watermark and causally ordered repeated-trigger tests;
its independent review is pending.

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
- every failure/restoration row and twelve guard scenarios execute and pass;
- no product consumer imports the package;
- ADR-0047 remains Proposed with the concrete packet linked for maintainer review;
- canonical status documents distinguish scoped `PASSED` from blocked runnable,
  installed, guest and product work; and
- full required verification and an independent verdict on the corrected head
  before merge.
