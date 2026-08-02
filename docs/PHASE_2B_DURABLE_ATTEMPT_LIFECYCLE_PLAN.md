# Phase 2B durable attempt lifecycle implementation and conformance plan

Status: design only. Proposed ADR-0025 selects one colocated Supervisor snapshot and transaction
domain for future durable `AttemptID`-keyed fake lifecycle state. No durable lifecycle behavior,
consumer, IPC, evidence, runtime, real backend, or guest is implemented by this document.

Normative proposal:
[ADR-0025](adr/0025-colocate-durable-attempt-lifecycle-state.md).

## Objective and defensive scope

Defensively validate Capsule's durable before-effect lifecycle ordering using only the current
repository, the existing fixed Supervisor store, owned conformance fixtures, and
`registeredlifecycle.FakeBackend` in a local test worktree. Do not access another system, identity,
credential, user data, runtime, backend, or guest. Preserve `AttemptID`-only execution,
consume/create-before-effect ordering, and `FakeBackend.CreatesGuest() == false`.

## Compatibility baseline

The implementation oracle is the merged Slice B/C path, not `execution.SupervisorCore`:

- one fixed versioned snapshot atomically commits consumed approval plus immutable created attempt;
- `ResolveCreated` returns only a committed attempt and copied exact registration/plan bindings;
- `Drive` and `Recover` accept only `AttemptID`;
- two approvals for one registration remain independent;
- exact and concurrent replay do not redrive fake effects;
- recovery ignores approval usability and registration expiry after attempt commit; and
- fake fault outcomes end destroyed, unresolved, or recovery-fenced without a guest.

The durable design may refine active-attempt accounting once lifecycle terminal state is colocated,
but it may not change any authority binding or add an execute-by-registration path.

## Selected architecture

Extend the existing fixed installation snapshot with a closed lifecycle collection and digest.
Every lifecycle transition shares the authority store's validation and atomic replacement
transaction. No second state file, database, journal, or lifecycle writer is introduced.

The first implementation remains an unwired fixed-store checkpoint. A production database,
platform lock, archive, backup, anti-rollback anchor, consumer, and real adapter are later gates.

## Falsifiable state and effect oracles

### Lifecycle establishment oracle

| Boundary | Durable precondition | Adapter calls | Durable/recovery result |
| --- | --- | ---: | --- |
| Before lifecycle record | one committed created attempt and consumed approval; no lifecycle record | 0 | startup recovery candidate |
| Creation commit confirmed abort | same committed authority state | 0 | no lifecycle record; retry may create the one record |
| Creation commit indeterminate | no record or one complete `prepare-pending` record | 0 | process fenced; reopen decides exactly which state exists |
| Creation commit confirmed | one record copied from colocated authority state, cleanup true, prepare pending | 0 | process death/restart resumes prepare by the same `AttemptID` |
| Duplicate/concurrent creation | the first complete record exists | 0 | exact binding returns that record; mismatch is repair-required |

Lifecycle establishment never consumes another approval, creates another attempt, accepts caller
bindings, or calls the adapter.

### Uniform transaction/death oracle

This table applies to each prepare/create/start/observe/stop/destroy effect.

| Boundary | Durable precondition | Permitted adapter calls | Durable postcondition | Reopen/recovery oracle |
| --- | --- | ---: | --- | --- |
| Before intent transaction | exact predecessor version; no other in-flight effect | 0 | predecessor unchanged | resume from predecessor |
| Intent commit confirmed abort, including `FULL`/`IOERR` before rename | predecessor | 0 | predecessor unchanged | resume from predecessor; no reconcile needed |
| Intent commit indeterminate | predecessor or complete intent may be on disk | 0 | current process recovery-fenced | reopen finds exactly predecessor or intent; never call a new effect before that decision |
| Intent confirmed | operation, sequence, stable effect ID, cleanup, and owner permit durable | at most 1 holder of the sealed permit | intent remains until after-effect commit | process death reconciles that exact effect ID |
| Adapter confirmed not applied | intent | 1 | fixed failure/cleanup transition or same-ID retry eligibility commits | never infer application |
| Adapter applied and response returned | intent | 1 | confirmed checkpoint must commit before successor | commit failure fences; reopen reconciles intent |
| Adapter effect applied but response lost/after-effect fault | intent | 1 | intent or explicit indeterminate state; no successor effect | reconcile exact effect ID; never allocate a replacement effect ID |
| After-effect commit confirmed abort or ENOSPC/I/O failure | adapter may have applied effect | 1 | durable intent remains; in-process recovery fence | reopen and reconcile; no later operation runs |
| After-effect commit indeterminate | intent or complete confirmed checkpoint may be on disk | 1 | in-process recovery fence | reopen finds one complete state and reconciles only if intent remains |
| Process death after confirmed checkpoint | successor state durable | 1 | checkpoint and last-confirmed fields agree | resume from successor; do not redrive confirmed effect |
| Corrupt or impossible reopen | unknown | 0 | original bytes retained, open result repair-required | no adapter call or rewrite; explicit repair required |

All concurrency tests count calls. A second goroutine, component instance, startup pass, or exact
replay must either observe the first record or wait on the per-attempt coordinator; it cannot hold a
second permit.

### Reconciliation transaction oracle

| Boundary | Adapter calls | Durable/recovery result |
| --- | ---: | --- |
| Recovery-count/eligibility commit confirmed abort | 0 | intent and prior retry schedule unchanged |
| Recovery-count/eligibility commit indeterminate | 0 | process fenced; reopen prior or incremented schedule |
| Eligibility commit confirmed | 0 before commit, then one observation | count and next eligible time already durable |
| Observation response lost or unknown | 1 observation, 0 effects | unresolved intent and consumed recovery slot |
| Observation result commit confirmed abort/ENOSPC/I/O | 1 observation, 0 effects | process fenced; intent remains; reopen may repeat only observation |
| Observation result commit indeterminate | 1 observation, 0 effects | reopen intent/schedule or one complete result; never redrive effect from result ambiguity |
| Result commit confirmed | 1 observation, 0 effects | applied/not-applied/absence/mismatch rule advances exactly once |

### Per-operation pre/post/recovery oracle

| Operation | Required predecessor | Durable intent | Confirmed after-effect state | Recovery decision |
| --- | --- | --- | --- | --- |
| `prepare` | `prepare-pending`, no instance | `prepare-intent`, cleanup true | `prepared` and last confirmed `prepare` | applied advances; not-applied retries same effect ID; unknown unresolved |
| `create` | `prepared`, no instance | `create-intent`, cleanup true | `created` with one exact nonzero instance identity | applied must return same identity; not-applied retries same effect ID; unknown unresolved; mismatch quarantines |
| `start` | `created`, exact instance | `start-intent` naming that instance | `started` | applied advances; exact not-applied may retry same ID; missing/stale/mismatch unresolved or quarantined |
| `observe` | `started`, exact instance | `observe-intent` naming that instance | `observed` | applied advances; exact not-applied may retry same ID; unknown unresolved |
| `stop` | `observed` or cleanup path with exact instance | `stop-intent` naming that instance | `stopped` | applied advances; exact not-applied/present may retry same ID; unknown unresolved |
| `destroy` | cleanup path with exact instance | `destroy-intent` naming that instance | `destroy-confirmed`, cleanup still true | applied then authoritative absence commits `destroyed`; not-applied/present retries same ID; unknown unresolved |

An authoritative absence result before any create could have applied may close cleanup as destroyed.
Absence after a confirmed create/start/observe without a confirmed destroy remains unresolved. A
PID, path, process name, missing handle, empty list, or fake map miss is not by itself authoritative
absence for a future real backend.

### Record consistency oracles

Reopen validation must reject all of these as repair-required:

- lifecycle without a committed created attempt or consumed approval;
- two lifecycle records for one attempt, one approval owning two attempts, or two attempts sharing
  an effect or instance identity;
- changed copied authority, plan-role, backend, or immutable binding digest;
- lifecycle state inconsistent with cleanup, effect status, operation, instance identity, last
  confirmed checkpoint, terminal time, recovery count, or next-retry time;
- record version zero/decrease, unsupported storage version, capacity overflow, set-digest mismatch,
  timestamp rollback, or timestamp above the snapshot high-water; and
- destroyed with cleanup true, unresolved/quarantined with cleanup false, or create-confirmed
  without an instance identity.

No corruption case rewrites, deletes, truncates, recreates, normalizes, or drops the evidence file.

## Startup algorithm and ownership

The startup coordinator has one exact order:

1. acquire the installation exclusive-open port; the production macOS target is the pre-created,
   no-symlink, installation-owned mode-0600 sibling lock file and lifetime nonblocking exclusive
   descriptor lock proposed by ADR-0025, while the first slice injects an in-process substitute;
2. open the existing store without create-on-missing fallback;
3. validate top-level version, bounds, set digests, cross-links, lifecycle invariants, and migration
   completion;
4. establish one owner-session token and keep the installation lock for the full process lifetime;
5. persist a trusted-clock high-water observation when available;
6. enumerate a sorted unique join of committed created attempts whose lifecycle is missing or not
   durably destroyed;
7. call `Recover(AttemptID)` once for each eligible record in the initial pass;
8. retain the startup fence if any record is repair-required, quarantined, unresolved, or automatic
   recovery exhausted; and
9. only after a clean pass report internal readiness. This first slice has no consumer to enable.

Automatic reconciliation eligibility is durable: attempt 1 is immediate, attempt 2 is no earlier
than high-water plus one second, and attempt 3 no earlier than high-water plus five more seconds.
The count and next time commit before the observation. After three unknowns, startup and timers make
no further adapter call until a future explicit repair API resets the schedule. That repair API is
not part of the first slice.

A trusted-clock failure prevents new forward prepare/create/start/observe work. It does not prevent
stop/destroy/reconcile cleanup, which may reuse the current high-water timestamp. Time never moves
backward and expiry is never re-evaluated for a committed attempt.

The first slice uses an injected exclusive-open/coordinator port and proves duplicate-owner refusal
within the local harness. Consumer activation requires a selected platform implementation and
real multi-process/process-death tests; the in-process port is not that evidence.

## Store and adapter failure matrix

The retained fault harness must inject each row at every applicable operation.

| Fault | Required result |
| --- | --- |
| before intent mutation | predecessor unchanged; zero adapter calls |
| intent validation/capacity failure | no intent, no eviction, zero adapter calls |
| intent temp create/write/sync/close/rename failure | confirmed abort when knowable; zero adapter calls |
| intent post-rename directory-sync failure | recovery-fenced; reopen predecessor or intent; zero adapter calls before reopen |
| fake before-effect fault | confirmed not-applied; fixed failure and cleanup oracle |
| fake after-effect fault | effect indeterminate; intent retained; restart reconciliation |
| process death immediately after each fake effect | intent retained; new component/store instance resolves exact effect once |
| after-effect temp/write/sync/rename failure | recovery-fenced; intent remains; restart reconciliation |
| after-effect post-rename directory-sync failure | reopen intent or confirmed checkpoint; no effect replay when checkpoint exists |
| lifecycle-record creation confirmed/indeterminate failure | no fake effect; reopen missing record or one complete prepare-pending record |
| reconciliation eligibility/result commit failure | observation-only recovery fence; never redrive an effect from reconciliation ambiguity |
| create success with lost instance response | reconcile by effect ID; adopt only the exact returned identity |
| stale/wrong instance identity | no stop/destroy against it; quarantined and installation repair-required |
| backend reports exact not-applied | same logical effect ID may be retried; no new effect identity |
| backend reports unknown | unresolved, cleanup retained, durable bounded backoff |
| backend state missing after confirmed create | unresolved unless exact adapter proof establishes the relevant destroy effect and absence |
| destroy response without absence proof | `destroy-confirmed`, cleanup retained |
| repeated startup before next retry time | no adapter call and no state mutation |
| third unknown reconciliation | automatic recovery exhausted; later startup performs no automatic call |
| lifecycle active/retained cap plus one | `CAPACITY`; no eviction or partial record |
| unsupported version/corruption/truncation/trailing data | repair-required; file unchanged |
| explicit migration before rename failure | v0 remains valid |
| explicit migration post-rename indeterminate | reopen exactly v0 or v1; never merge states |

## Internal APIs and likely files

The names below are implementation targets, not public contracts.

### Passive record and adapter types

Add `internal/execution/lifecyclestate` with closed types independent of a runtime or real backend:

```go
type Operation string
type RecordState string
type EffectStatus string
type ReconcileStatus string
type RecordVersion uint64
type EffectID [16]byte
type BackendInstanceIdentity struct { Kind string; Value []byte }
type BackendBinding struct { /* closed digests and CreatesGuest */ }
type Record struct { /* ADR-0025 fields */ }
type EffectPermit struct { /* sealed owner/record/effect binding */ }
type EffectResult struct { /* applied/not-applied/indeterminate + identity */ }
type ReconcileResult struct { /* closed status + exact identities */ }
```

This package performs no I/O and imports no experiment, runtime, backend, content, IPC, or evidence
package.

### Colocated store port

Extend `internal/execution/registrationstate` without importing `registeredlifecycle`:

```go
type DurableLifecycleStore interface {
    EnsureLifecycle(context.Context, AttemptID, lifecyclestate.BackendBinding) (lifecyclestate.Record, bool, error)
    ReadLifecycle(context.Context, AttemptID) (lifecyclestate.Record, error)
    BeginEffect(context.Context, AttemptID, RecordVersion, lifecyclestate.Operation) (lifecyclestate.EffectPermit, error)
    ConfirmEffect(context.Context, lifecyclestate.EffectPermit, lifecyclestate.EffectResult) (lifecyclestate.Record, error)
    RecordIndeterminate(context.Context, lifecyclestate.EffectPermit, Classification) (lifecyclestate.Record, error)
    RecordReconciliation(context.Context, AttemptID, RecordVersion, lifecyclestate.ReconcileResult) (lifecyclestate.Record, error)
    RecoveryAttemptIDs(context.Context) ([]AttemptReference, error)
}
```

Likely files:

- `internal/execution/registrationstate/types.go`: lifecycle store port and snapshot projection;
- `internal/execution/registrationstate/fixed_store.go`: v1 envelope, lifecycle collection/digest,
  explicit migration, transaction fault points, and reopen validation;
- `internal/execution/registrationstate/approval_attempt_validation.go`: cross-record and active
  capacity validation; and
- new focused `lifecycle_store_test.go` and migration/corruption tests.

The immutable `ExecutionAttempt` type remains unchanged. Active-attempt counting becomes a joined
predicate: created attempt plus no destroyed lifecycle record.

### Lifecycle driver

Replace `MemoryStore` use in `internal/execution/registeredlifecycle/component.go` with the durable
port. Add an injected per-attempt coordinator and owner-session provider. Keep `Drive`, `Recover`,
and startup entry points `AttemptID`-only.

Refine `fake_backend.go` so every operation consumes the sealed effect identity and returns a
closed applied/not-applied/indeterminate result. Reconciliation observes the exact effect/instance
identity and has no external effect. The fake caches one result per effect ID and applies that
logical effect at most once even when the same ID is retried. Keep the constructor hard-fenced on
`CreatesGuest() == false`. Remove `memory_store.go` only after all migrated tests pass; do not retain
a permissive fallback.

## Small implementation slices

### Slice E1: passive lifecycle contract

- Add the closed runtime-neutral record, effect, instance, result, and validation types.
- Add byte/collection ceilings and defensive-copy tests.
- Preserve every copied Slice B binding and prove wrong-domain/cross-link rejection.

Acceptance: passive tests only; no store write, adapter call, consumer, or guest.

### Slice E2: fixed snapshot v1 and explicit migration

- Add lifecycle collection/set digest and v1 validation.
- Add explicit lock-held v0-to-v1 migration with pre/post-rename fault oracles.
- Refuse auto-create, downgrade, unsupported versions, corruption, and partial lifecycle records.

Acceptance: reopen/migration/corruption tests pass; lifecycle still cannot call the fake backend.

### Slice E3: transactional lifecycle checkpoints

- Add ensure/read/begin/confirm/indeterminate/reconcile/recovery-set transactions.
- Reserve lifecycle capacity when attempt creation commits or prove from matched ceilings that it
  cannot fail after an attempt exists.
- Inject confirmed and indeterminate persistence failures at every transaction boundary.

Acceptance: no transaction test calls a backend; every intent and after-effect oracle is exact.

### Slice E4: fake-only durable driver and startup coordinator

- Convert the fake backend to stable effect and instance identities.
- Replace `MemoryStore` with the v1 fixed store.
- Serialize all `Drive`/`Recover`/startup calls per attempt under one injected owner session.
- Implement the six operation protocols, exact reconciliation, durable backoff, and terminal/
  unresolved dispositions.

Acceptance: all migrated Slice C tests plus the new death/fault matrix pass with a newly reopened
store and newly constructed component. `FakeBackend.CreatesGuest() == false` is asserted in every
constructor path. No consumer, IPC, content, evidence, runtime, real adapter, process, or guest is
introduced.

### Slice E5: capacity, repeated-startup, and documentation checkpoint

- Prove exact 256 active and 4,096 retained ceilings, capacity release only after durable destroyed,
  no eviction, and no retry after exhaustion.
- Prove duplicate owner refusal through the injected port and retain the real multi-process blocker.
- Reconcile canonical documentation and the evidence matrix as `local-mechanic` only after code and
  tests exist.

Acceptance: no production or guest-facing claim is promoted.

## Required focused tests

Retain the current 12 top-level `registeredlifecycle` behaviors and add focused tests with exact
oracles equivalent to:

- `TestDurableLifecycleRecordCopiesCommittedAttemptBindings`;
- `TestLifecycleCreationCommitAndProcessDeathMatrix`;
- `TestLifecycleIntentCommitsBeforeEveryFakeEffect`;
- `TestIntentConfirmedAbortCallsNoEffect`;
- `TestIntentAndConfirmationIndeterminateReopenMatrix`;
- `TestLostEffectResponseReconcilesSameEffectID`;
- `TestConfirmedEffectCommitFailureFencesUntilReopen`;
- `TestProcessDeathAtEveryEffectBoundaryUsesOneOwnerAndEffect`;
- `TestCreateRecoveryRequiresExactInstanceIdentity`;
- `TestMissingStaleAndUnknownBackendStateFailClosed`;
- `TestDestroyRequiresAuthoritativeAbsenceBeforeCleanupClears`;
- `TestRepeatedStartupHonorsDurableRetryBudgetAndBackoff`;
- `TestReconciliationCommitFailureRepeatsObservationNotEffect`;
- `TestStartupEnumeratesOnlyCommittedCreatedNonterminalAttempts`;
- `TestRecoveryIgnoresApprovalAndRegistrationExpiry`;
- `TestLifecycleCapacityNeverEvictsTerminalOrCleanupState`;
- `TestLifecycleV0ToV1MigrationAndDowngradeRefusal`;
- `TestLifecycleCorruptionEntersRepairRequiredWithoutRewrite`; and
- `TestDuplicateOwnerRefusesBeforeStoreMutationOrFakeEffect`.

Every test records before/after authority, lifecycle-set digest, record version, operation/effect
identity, cleanup flag, recovery fence, adapter call count, and reopen result. Cap-plus-one cases use
compact generated state rather than thousands of opaque fixtures.

## Activation blockers and dependencies

The fake-only fixed-store slice does not satisfy any of these:

1. **Archive/compaction:** accepted active/archive transaction, replay tombstones, retention,
   secure-deletion limitations, power-loss evidence, and continuous-service capacity.
2. **Multi-process and platform ownership:** selected Supervisor topology, real installation lock,
   process crash, session/reboot/update behavior, and duplicate-owner refusal across processes.
3. **Rollback, backup, and uniqueness:** coherent backup/restore protocol, independent checkpoint or
   explicit rollback limitation, and durable non-reuse for registration/approval/attempt/nonce/
   effect identities across archive and restore.
4. **Production backend reconciliation:** exact adapter-specific effect/instance identity,
   enumeration, authoritative absence, stale-object defense, idempotency, and retained adversarial
   evidence. No real adapter is authorized by this plan.
5. **Public cutover:** accepted ADR-0019/0024 dependencies as applicable, production Approval
   verification, authenticated IPC, consumer ownership, evidence composition, daemon aggregate
   envelope, and atomic removal of dormant direct-execution paths.

No blocker may be inferred complete from fixed-file or fake-backend evidence.

## Verification for each retained implementation slice

Use Node.js 22.22.1 or newer, pnpm 10, and Go 1.23 or newer:

```sh
fnm exec --using=22.22.1 -- pnpm install --frozen-lockfile
fnm exec --using=22.22.1 -- pnpm check
fnm exec --using=22.22.1 -- pnpm lint
fnm exec --using=22.22.1 -- pnpm test
fnm exec --using=22.22.1 -- pnpm verify:schemas
fnm exec --using=22.22.1 -- pnpm format:check
go test ./...
go vet ./...
go build ./...
git diff --check
```

Also run the focused Slice B/C and new lifecycle store, migration, fault, capacity, and repeated-
startup tests. Markdown links and ADR status/index entries must be checked before handoff.
