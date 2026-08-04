# ADR-0025: Colocate durable attempt lifecycle state with Supervisor authority state

- Status: Proposed
- Date: 2026-08-02
- Slice E1-E5 implementation checkpoint: 2026-08-03
- Refines if accepted: ADR-0011, ADR-0012, ADR-0013, ADR-0023, and ADR-0024

## Context

ADR-0024 and Phase 2B Slice B establish one versioned, fault-injectable Supervisor snapshot that
atomically consumes an approval and creates one immutable `ExecutionAttempt` before any backend
effect. Slice C then resolves that committed attempt by `AttemptID` and drives a no-guest fake
backend, but its lifecycle `MemoryStore` is bounded, single-process, and non-durable. Reusing that
store across component instances tests only in-process restart behavior.

The next boundary must preserve all of the already-retained authority properties while making the
prepare/create/start/observe/stop/destroy checkpoints durable. A second lifecycle store would
introduce two independently committed truths about whether a committed attempt owns cleanup work,
whether an effect may be retried, and whether capacity was released. That ambiguity is not needed
for the first local fake-only slice.

This ADR proposes storage and coordination semantics only. It does not accept ADR-0019 or ADR-0024,
activate a consumer, select authenticated IPC, compose evidence, admit a runtime or real backend,
or create a guest.

## Proposed decision

### One authoritative installation snapshot

Lifecycle records will be colocated with registrations, durable time high-water, approvals, and
immutable created attempts in the same versioned Supervisor installation snapshot and transaction
implementation. The store remains owned only by the Execution Supervisor.

The consume/create transaction remains the sole authority-creation boundary:

```text
usable approval + no attempt
  -- one committed Supervisor transaction -->
consumed approval + one immutable created attempt
```

Only after that commit is confirmed durable may the lifecycle component create a lifecycle record.
Creating the record is a second Supervisor transaction, keyed only by the committed `AttemptID`. It
copies bindings from the store's own attempt/approval/registration records, sets a cleanup
obligation, and records `prepare` as the next operation before any fake effect. It creates no new
approval or attempt authority.

Every lifecycle mutation uses the same snapshot generation, validation, temporary-file sync,
atomic rename, and directory-sync boundary as the authority records. A complete snapshot therefore
cannot contain a lifecycle record without its committed attempt, cannot make an attempt terminal
without its lifecycle disposition, and cannot disagree with another local database about cleanup.

### Exact durable lifecycle record

There is at most one record for an `AttemptID`. The fixed first record is a closed structure with
these fields:

| Group | Required fields and rules |
| --- | --- |
| Key and version | nonzero `AttemptID`; lifecycle record format version; monotonically increasing positive record version; snapshot generation at last commit |
| Copied authority | `ApprovalID`, `AttemptNonce`, `RegistrationID`, registration sequence, plan digest, installation ID, epoch sequence and digest, Supervisor ID, exact approval purpose and audience, approval payload digest, resolved Approval-key authorization identity, attempt `createdAt`, and attempt/approval/registration storage versions |
| Copied plan bindings | source-manifest digest, inline-input digest, runtime-bundle-manifest digest, ordered profile-review-attestation digests, profile-registry-entry digest, backend-validation-record digest, backend-configuration digest, trust-snapshot digest, and policy-decision digest |
| Backend binding | closed adapter kind, adapter protocol version, implementation identity digest, backend configuration digest, backend validation-record digest, and `createsGuest` fixed to `false` for the first slice |
| Effect identity | positive operation sequence, operation, nonzero Supervisor-issued opaque effect ID, and effect status `none`, `intent`, `confirmed`, or `indeterminate`; an effect ID is never caller supplied and is reused for reconciliation or retry of that same logical effect |
| Backend instance | optional closed instance-identity kind plus bounded opaque bytes and digest; absent before confirmed create; immutable after first confirmed create or exact reconciliation |
| Lifecycle | closed state, cleanup-required boolean, last confirmed checkpoint, first fixed failure classification and operation, terminal or unresolved disposition, and automatic-recovery exhaustion flag |
| Recovery | last reconciliation result `none`, `applied`, `not-applied`, `authoritatively-absent`, `present`, `unknown`, or `identity-mismatch`; automatic recovery-attempt count; next-recovery Unix second; recovery-fence reason code |
| Time | lifecycle-opened, last-transition, and optional terminal Unix seconds in `UInt53`; all are no earlier than attempt `createdAt`, never decrease, and never exceed the installation durable time high-water used by the same committed snapshot |
| Integrity | lifecycle-set digest and a closed digest over the immutable binding projection; no free-form backend, caller, path, plan, approval, or guest text |

Exact plan bytes remain authoritative in the colocated registration record and are not duplicated
inside every lifecycle record. Before the first lifecycle transaction, the lifecycle component
still independently decodes and hashes those bytes. The store transaction then independently
cross-checks the copied projection against its own registration, consumed approval, and created
attempt. Recovery does not recheck approval usability or registration expiry.

`cleanupRequired` becomes `true` when the lifecycle record commits, before `prepare`. It remains
true through every intent, failure, indeterminate result, backend instance, unresolved state, and
quarantine. Only a committed `destroyed` disposition following an authoritative absence
observation may clear it.

The first store obtains every effect ID from a cryptographic random source immediately before the
intent transaction, rejects zero and collisions across all retained lifecycle records, and commits
the ID only with that intent. An operation retry reuses the committed ID; it never draws another.
`EffectID` is a distinct nominal 16-byte domain and is never accepted as an approval, attempt,
nonce, registration, or backend instance ID. This proves bounded local uniqueness only. Archive
tombstones plus rollback-resistant checkpointing remain production requirements.

### Closed state and checkpoint model

The durable states are:

```text
prepare-pending -> prepare-intent -> prepared
  -> create-intent -> created
  -> start-intent -> started
  -> observe-intent -> observed
  -> stop-intent -> stopped
  -> destroy-intent -> destroy-confirmed
  -> destroyed

any nonterminal state -> unresolved
identity/cross-link mismatch -> quarantined + installation repair-required
```

`destroyed` is the only ordinary terminal lifecycle disposition in the first fake-only slice. It
means cleanup is complete for the fake lifecycle; it is not job success, guest completion, backend
validation, or evidence. `unresolved` and `quarantined` are retained nonterminal dispositions with
cleanup still required.

Every prepare/create/start/observe/stop/destroy operation uses one protocol:

1. Acquire the in-process per-attempt coordinator and re-read the durable record.
2. In one confirmed transaction, require the exact predecessor record version, allocate or reuse
   the operation's effect ID, and commit the operation's `intent` state before calling the adapter.
3. Call the adapter only with the sealed store-returned permit containing `AttemptID`, operation,
   operation sequence, effect ID, immutable binding digest, owner-session ID, and any exact backend
   instance identity.
4. Classify the adapter result as confirmed applied, confirmed not applied, or indeterminate. An
   error or lost response is indeterminate unless the adapter contract proves no effect occurred.
5. Commit the corresponding after-effect checkpoint. A create checkpoint cannot confirm without a
   nonzero exact instance identity. A destroy response does not clear cleanup; reconciliation must
   subsequently report authoritative absence for that exact identity.

The previous confirmed checkpoint and durable intent are the recovery oracle. An absent later
checkpoint never proves the effect did not occur.

### Failure and process-death rules

- A confirmed lifecycle-record creation failure leaves the committed attempt without a lifecycle
  record, calls the adapter zero times, and remains startup work. An indeterminate creation commit
  fences the process; reopen finds either no lifecycle record or one complete `prepare-pending`
  record. Process death after that record commits but before prepare resumes from `prepare-pending`.
- A confirmed failure of the before-effect intent transaction calls the adapter zero times and
  leaves the prior durable state authoritative.
- An indeterminate before-effect transaction calls the adapter zero times, fences all mutation in
  the current process, and requires reopen. Reopen may find the prior state or the intent; it never
  invents a third state.
- A confirmed adapter `not-applied` result may follow the operation's fixed failure path. All other
  adapter errors, after-effect injected faults, timeouts, cancellation, and lost responses leave the
  intent as recovery work.
- If the adapter reports success but the after-effect transaction has a confirmed `FULL`, `IOERR`,
  sync, or other commit failure, the current process enters a recovery fence immediately. No later
  effect is called. Reopen sees the durable intent and reconciles it.
- An indeterminate after-effect commit similarly fences mutation. Reopen finds either the intent or
  the complete confirmed checkpoint.
- Process death before the intent commit has no effect. Process death after intent and before a
  confirmed after-effect checkpoint is reconciled by the same effect ID. Process death after the
  checkpoint resumes only from its recorded successor.
- Store corruption, an unsupported version, a digest/cross-link failure, or an impossible state
  returns `repair-required`. The original file is not rewritten, deleted, truncated, or silently
  recreated. The failed validation itself is the durable fence until an explicit repair procedure
  preserves and replaces the evidence.
- A backend instance identity mismatch, stale reuse, or observation that could name another
  attempt enters `quarantined` plus installation `repair-required`; Capsule does not signal, stop,
  destroy, or adopt the mismatched object.
- Missing backend state is authoritative absence only when the adapter proves absence for the exact
  effect or instance identity. Missing process/handle/path state by itself is `unknown`.
- Unknown reconciliation records `unresolved`, retains cleanup, and never becomes ordinary success
  or capacity-releasing terminal state.

### Reconciliation and bounded automatic retry

Reconciliation is an observation-only adapter operation over the exact stored effect and instance
identity. It returns a closed result and performs no create/start/stop/destroy action itself.

Recovery commits the attempt count and next eligible time before calling reconciliation. A
confirmed failure of that pre-observation transaction makes zero adapter calls. An indeterminate
transaction fences until reopen. If the observation returns but its result cannot be durably
committed, the process fences with the original effect intent still authoritative; reopening may
repeat only the observation, not the effect. An indeterminate result commit reopens to either the
old intent/recovery count or the complete reconciliation result.

For an intent, `applied` commits the missing after-effect checkpoint; `not-applied` permits one
retry with the same effect ID; and `unknown` retains unresolved cleanup. A confirmed created or
later state whose exact instance is unexpectedly absent is unresolved, not destroyed. Only a
confirmed destroy plus authoritative absence reaches `destroyed`.

Automatic reconciliation is installation-global and bounded per record across restarts. The first
fixed policy permits three automatic observations: the first immediately, the second no earlier
than one durable effective-time second later, and the third no earlier than five additional
seconds later. The count and next eligible time commit before each observation. After a third
unknown result, `automaticRecoveryExhausted` is set and repeated startup performs no further adapter
call without an explicit Supervisor-owned repair action. A process death or lost response after the
pre-observation commit safely consumes that bounded slot. Trusted-clock failure blocks new forward
effects but does not block destructive cleanup or reconciliation; those may retain the prior
high-water timestamp rather than move time backward.

### Startup and single-owner coordination

Startup occurs in this order:

1. Acquire one nonblocking installation-scoped exclusive owner lock before opening state. On macOS,
   Proposed ADR-0033 selects a pre-created sibling object opened relative to the enrolled protected
   state root with `O_NOFOLLOW|O_CLOEXEC`, validated by exact UID/mode/type/link/device/inode, and
   held by `flock(LOCK_EX|LOCK_NB)` for the entire process lifetime. Normal startup never creates a
   missing object. Because the store itself is atomically replaced, the stable sibling inode—not
   the snapshot inode—is locked. A second Supervisor process refuses with a fixed busy result
   before store read-modify-write, recovery, archive, or adapter work. Process death closes the last
   inherited description and releases the lock. The local semantics and passive G1 Go/Darwin
   acquisition now pass under owned temporary roots; owner-required store/startup composition and
   the installed protected-root matrix remain unimplemented.
2. Open without automatic creation, enforce file type/ownership/permission policy, bound the read,
   decode the exact supported version, and validate snapshot, set digests, capacities, immutable
   bindings, cross-links, state transitions, and time rules.
3. Refuse unsupported version or corruption as repair-required without rewriting the file.
4. Establish an owner-session ID for sealed in-process effect permits. Per-attempt coordination
   serializes `Drive`, `Recover`, startup recovery, and later cancellation within that sole owner.
5. Enumerate a sorted, duplicate-free recovery set by joining committed created attempts with
   lifecycle records. Include only attempts with no lifecycle record or with nonterminal,
   cleanup-bearing, unresolved, or quarantined state. A lifecycle record without its committed
   attempt is corruption.
6. Call idempotent `Recover(AttemptID)` for each eligible record. It never accepts a registration
   ID, approval reference, plan bytes, backend flags, image, mount, path, or replacement binding.
7. Keep attempts disabled until the initial recovery pass finishes. Any repair-required,
   quarantined, exhausted, or still-unresolved record leaves the recovery fence closed.

Repeated startup uses the durable state/version, effect ID, recovery count, and next eligible time.
It cannot create another lifecycle owner or a differently identified logical effect for the same
operation.

The fixed first implementation may exercise only an in-process coordinator and injected
exclusive-open port; it must remain unwired and may not claim the macOS multi-process mechanism
until the selected Supervisor language/platform implementation passes duplicate-process, crash,
replacement, session, update, and hostile-file tests.

### Capacity, retention, and archive boundary

The fixed snapshot permits at most 4,096 retained lifecycle records and at most 256 attempts whose
lifecycle is missing, nonterminal, cleanup-bearing, unresolved, or quarantined. These ceilings
match the existing retained/open attempt ceilings. Lifecycle creation for a committed attempt must
therefore have already-reserved capacity; post-effect `FULL` may not strand recovery identity.

An immutable attempt becomes terminal for active-capacity accounting only when its colocated
lifecycle record is durably `destroyed` with cleanup false. The attempt, consumed approval, nonce,
registration relationship, lifecycle record, effect IDs, instance identity, and explanatory
failure remain retained. The fixed store performs no eviction, archive, compaction, tombstone
deletion, capacity-driven transition, or secure deletion.

Consumer or continuous-service activation is blocked on a separate accepted archive/compaction
design. That design must preserve replay and uniqueness tombstones, cleanup and unresolved work,
sole explanatory evidence, cross-record digests, and a transactionally checkable active/archive
boundary. Raising a ceiling is not archival.

[Proposed ADR-0031](0031-checkpoint-closed-supervisor-cohorts.md) now defines that boundary but is
not accepted or implemented. It permits only complete expired registration cohorts whose attempts
are durably destroyed after authoritative absence, retains full records and exact tombstones, and
forbids referenced-history deletion. Its finite fixed-store checkpoint is a conformance oracle,
not the production engine or continuous-service mechanism required here.

### Versioning, migration, backup, and rollback

Adding lifecycle state requires a new top-level snapshot version. Version 0 never treats a missing
lifecycle collection as an empty production collection. The first unwired migration is an explicit,
offline, lock-held `v0 -> v1` transaction: fully validate v0, add an empty lifecycle set and digest,
record the source version, sync/rename/sync once, then reopen and validate v1. It is not an
automatic fallback. Failure before rename preserves v0; an indeterminate rename requires reopen;
old binaries refuse v1 rather than downgrade it.

Ordinary file-copy backup is not a coherent backup protocol. A future backup must take a verified
snapshot under the owner lock and bind its snapshot generation and installation/epoch checkpoint.
Restore cannot silently reset attempts, approvals, nonces, effect IDs, cleanup, or repair state.

Retained identifiers prevent reuse only within the visible installation history. A coherent
restore can replay an older otherwise valid world, including its random-ID namespace. Production
activation therefore requires either an independently protected latest checkpoint/non-rollbackable
anchor or an explicit posture that detects but does not claim to prevent coherent rollback. Archive
must retain nonce, ID, effect-ID, and sequence tombstones. The fixed unwired store supplies none of
these stronger backup or rollback guarantees.

### Internal interface boundary

The lifecycle remains internal and `AttemptID`-only:

```go
type DurableLifecycleStore interface {
    EnsureLifecycle(context.Context, AttemptID, BackendBinding) (Record, bool, error)
    ReadLifecycle(context.Context, AttemptID) (Record, error)
    BeginEffect(context.Context, AttemptID, RecordVersion, Operation) (EffectPermit, error)
    ConfirmEffect(context.Context, EffectPermit, EffectResult) (Record, error)
    RecordIndeterminate(context.Context, EffectPermit, Classification) (Record, error)
    RecordReconciliation(context.Context, AttemptID, RecordVersion, ReconcileResult) (Record, error)
    RecoveryAttemptIDs(context.Context) ([]AttemptID, error)
}

func (component *Component) Drive(context.Context, AttemptID) (Snapshot, error)
func (component *Component) Recover(context.Context, AttemptID) (Snapshot, error)
func (component *Component) RecoverCreatedAttempts(context.Context) ([]Snapshot, error)
```

The concrete store builds the immutable record from its own colocated authority records. The
binding, effect, and instance types are runtime-neutral, but the first implementation accepts only
the closed fake binding and rejects any adapter for which `CreatesGuest()` is true. A later real
adapter requires a separate proposal and validation without changing execute authority.

### Current implementation checkpoint

Slices E1 through E5 now implement the passive types, explicit fixed-store v1 migration/open
validation, the complete `DurableLifecycleStore` transaction port above, and the unwired
FakeBackend-only driver. Retained tests cover stable effect/instance identity, every fake effect
boundary, reconciliation, concurrent/repeated startup, owner/coordinator mismatch, three-observation
automatic-recovery exhaustion, exact 256 active and 4,096 retained lifecycle ceilings, and
cap-plus-one refusal without eviction. The exact retained population encodes to 30,321,818 bytes
under a v1-only 64 MiB raw read bound; the v0 bound remains 16 MiB.

Joined v1 validation releases active capacity only for a durable `destroyed` record with cleanup
false and authoritative absence. `observed`, `stopped`, `destroy-confirmed`, `unresolved`, and
`quarantined` remain active and retained. The coordinator and owner session are still injected
in-process mechanics and do not implement the proposed macOS lock. The status remains Proposed,
and no production or guest-facing claim advances.

Proposed ADR-0033 now selects the exact BSD `flock` mechanism and startup/bootstrap contract but
does not change that implementation checkpoint. Owner-session IDs remain only sealed in-process
permit bindings and never substitute for the held OS descriptor.

## Alternatives considered

### Keep a separately durable lifecycle store

Rejected for this slice. A saga could coordinate it, but would require a durable outbox/inbox,
operation IDs, compensating rules, startup ownership, capacity reservation, and proof for every
authority-store/lifecycle-store commit ordering. That is strictly wider than adding one collection
to the existing Supervisor transaction and creates split-store states with no current consumer need.

### Store lifecycle only in an append-only journal beside the snapshot

Rejected. The journal and snapshot would still need an exact atomic checkpoint/compaction protocol,
power-loss ordering, replay validation, and corruption rules. It does not remove split authority.

### Replace the fixed snapshot with SQLite now

Rejected for the first fake-only slice. SQLite may be appropriate for a production Supervisor
store, but journal mode, sync policy, checkpointing, locking, backup, migration, corruption, and
power-loss evidence are separate activation decisions. Changing engines is not required to define
or test the lifecycle transaction protocol.

### Derive lifecycle from backend enumeration after restart

Rejected. Backend absence does not prove that no effect occurred or that cleanup completed, and a
backend cannot recreate the approval/attempt/plan authority record. Durable intent must precede
effects.

### Mutate the immutable `ExecutionAttempt` into the lifecycle state machine

Rejected. The committed created attempt is the permanent output of the consume/create authority
transaction. A separate colocated lifecycle record preserves that immutable fact while allowing
versioned effect and cleanup transitions.

## Consequences and activation blockers

- One transaction domain supplies the narrowest unambiguous authority/lifecycle/capacity view.
- A durable intent plus stable effect identity makes every lost response and process-death point
  recoverable without treating absence as proof or issuing a differently identified effect.
- The fixed snapshot remains an unwired conformance checkpoint, not a production database.
- Archive/compaction and continuous retention, implementation and installed validation of the
  selected multi-process lock, backup/restore and
  rollback-resistant uniqueness, a production backend's exact reconciliation contract, and any
  consumer/public cutover each remain explicit blockers.
- The first implementation remains fake-only and must continue to assert
  `FakeBackend.CreatesGuest() == false` at construction and in retained tests.

## Conformance plan

The exact small-slice, state/effect/death, migration, startup, and fault-injection work is defined in
[the Phase 2B durable attempt lifecycle plan](../PHASE_2B_DURABLE_ATTEMPT_LIFECYCLE_PLAN.md).
