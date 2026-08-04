# ADR-0031: Checkpoint closed Supervisor cohorts into immutable retained archives

- Status: Proposed
- Date: 2026-08-03
- Refines if accepted: ADR-0011, ADR-0012, ADR-0019, ADR-0023, ADR-0024, ADR-0025,
  and ADR-0029

## Context

The current unwired Supervisor authority path has reached durable lifecycle Slice E5. One fixed,
versioned v1 snapshot retains exact registrations, approval records, immutable created attempts,
lifecycle records, stable effect and fake-instance identities, set digests, and a durable time high
water. The retained tests establish these local oracles:

- one approval is consumed in the same transaction that creates one immutable `AttemptID` before
  any fake effect;
- lifecycle establishment, drive, recovery, and startup enumeration accept only `AttemptID`;
- a durable intent and stable `EffectID` precede each fake effect;
- active capacity is released only by a durable `destroyed` lifecycle record with cleanup false
  after authoritative absence;
- observed, stopped, destroy-confirmed, unresolved, quarantined, and cleanup-bearing records stay
  active;
- 256 active and 4,096 retained lifecycle records are exact ceilings; and
- retained-capacity exhaustion refuses without eviction, resurrection, or rewrite.

One limitation of that oracle matters here: a lifecycle record retains the current operation's
`EffectID`, and a later operation replaces that field. E5 therefore proves stable identity for an
in-flight/reconciled logical effect, not a durable set of every earlier effect ID. Archive cannot
reconstruct overwritten values. The v2 boundary adds a never-delete effect tombstone set for every
newly issued v2 effect; migration can seed only the nonzero effect IDs still visible in v1 records
and must label earlier fake-only history unavailable. No production global-effect-non-reuse claim
may be inferred across that migration.

That fixed snapshot is deliberately not a continuous-service store. At 4,096 retained lifecycle
records the existing E5 population encodes to 30,321,818 bytes under a 64 MiB v1 read bound. The
store has no archive, production owner lock, coherent backup, rollback-resistant checkpoint,
power-loss evidence, or production engine. Raising the ceiling would postpone rather than resolve
the boundary.

Archive and compaction are security-sensitive because the retained records are not merely history.
They prevent nonce and identifier reuse, make exact approval and attempt replay idempotent, preserve
the consumed-approval/attempt relationship, explain cleanup and failure state, and prevent an old
registration or approval from becoming usable after restart or restore. A missing record is not
proof that an approval was unused, an effect did not happen, an instance was destroyed, or an
identifier may be reused.

This ADR defines the next local conformance boundary. It does not activate a consumer, select a
production database, implement a platform owner lock, authorize deletion of retained authority
history, compose evidence, admit a runtime/backend, or create a guest.

## Proposed decision

### Boundary, authority, and status

Only the Execution Supervisor owns archive selection, creation, validation, activation, lookup,
backup, restore admission, and repair fencing. No daemon, Broker, updater, public API, IPC body, or
backend adapter can select archive records, supply archive paths, request deletion, reset a
tombstone, or claim that capacity was released.

The implementation series uses a minimal fixed-store archive checkpoint only as an unwired
conformance oracle. Slice F1 retains passive projections, exact limits/known answers, defensive
copies, and deterministic eligibility selection. The passive F2 format-correction slice now
resolves PR #78's three contradictions with scope-separated retained-global/segment-derived
indexes, kinded hot/archive record locations, exact hot/archive/total count equations, and a
domain-separated generation-one migration-genesis checkpoint. Generated exact answers and the
before/after contract are retained in the
[F2 format blocker resolution](../SUPERVISOR_ARCHIVE_F2_FORMAT_BLOCKER.md). The follow-on valid-v1
mapping contradiction is also passively resolved: each attempt index entry now carries an explicit
`absent | present(state, lifecycle location, lifecycle record digest)` union, and lifecycle counts
derive only from present arms. The executable
[v1 mapping resolution](../SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md) retains the original
committed-attempt-before-lifecycle witness and exact `attempts = 1, lifecycles = 0` genesis answer.
Stateful F2 implements the explicit owner-asserted v1-to-v2 migration and empty-archive full
verifier. Stateful F3 implements exactly one sealed immutable-segment prepare/verify/publish/
activate transaction with publish-before-reference ordering, atomic generation-two activation,
valid-orphan reporting without deletion, and complete predecessor-or-successor reopen. The
[F3 result](../SUPERVISOR_ARCHIVE_F3_ACTIVATION_RESULT.md) retains its exact known answers and
fault/corruption/substitution/concurrency/owner-loss/process-death limitations. Retained lookup,
v2 authority mutation, multi-segment growth, backup, cleanup policy, and later slices remain
unimplemented. The fixed oracle is not selected as the production engine.
SQLite remains the leading production-engine candidate, but its exact locking, journal/WAL,
checkpoint, sync, backup, migration, corruption, and real power-loss behavior require a separate
evidence-backed selection. Convenience, familiarity, or avoiding a migration is not a selection
criterion.

`FakeBackend.CreatesGuest() == false` remains mandatory throughout this boundary. Archive work has
no adapter call and does not change `Drive(AttemptID)`, `Recover(AttemptID)`, or consume/create-
before-effect ordering.

### Archive unit: one closed registration cohort

The indivisible archive unit is a **registration cohort**:

1. one exact retained registration entry, including its wire registration, exact plan bytes,
   resolved role bindings, registration index, sequence, expiry, and storage metadata;
2. every retained approval bound to that registration, including exact envelope, payload, and
   protected-header bytes, authorization identity, payload/envelope digests, nonce, state, and
   consumed-attempt link;
3. every immutable attempt bound through those approvals; and
4. every lifecycle record bound through those attempts, including immutable bindings, its current
   retained effect identity, backend instance identity, failure disposition, cleanup state,
   recovery counters, and timestamps; and
5. every v2 effect tombstone and instance tombstone bound to those attempts.

The fixed v1 record currently retains the latest logical operation `EffectID`, not a general
append-only effect history. The explicit v1-to-v2 migration records that limitation and seeds only
visible IDs. After v2 activation, the same transaction that commits a new effect intent also
inserts its exact effect tombstone before any fake call; later lifecycle operations never remove
it. This tombstone ledger prevents future v2 reuse but is not a transcript, does not prove an
effect occurred, and does not synthesize lost pre-v2 fake history. A later evidence design may add
an independently retained event set.

A cohort is eligible only when one lock-held, fully validated snapshot proves all of the following:

- durable effective time is at or beyond the registration expiry;
- the registration has no approval that is currently usable under that same durable high water;
- every consumed approval names exactly one retained attempt and every retained attempt names
  exactly one consumed approval in the cohort;
- every attempt has exactly one lifecycle record in `destroyed`, with cleanup false, terminal time
  present, and last reconciliation `authoritatively-absent`;
- no cohort record appears in `RecoveryAttemptIDs`;
- no record is unresolved, quarantined, repair-required, transition-fenced, recovery-fenced, or
  automatic-recovery exhausted;
- the cohort contains all records and indexes that refer to its registration, approvals, attempts,
  nonces, effect IDs, and instance identities; and
- no explicit retention hold identifies the cohort as the sole required input to trust repair,
  coherent backup verification, or future evidence composition.

An expired approval may still have durable state `usable`; expiry remains a predicate and archive
does not rewrite it to invalidated. Its exact record and replay projection remain retained. A
registration with no approvals may be archived after expiry. Cohorts are selected in ascending
registration sequence, then registration-ID byte order, so the same valid state and budget produce
the same batch.

No partial cohort may be archived. A registration remains hot if even one bound attempt still owns
cleanup or one cross-link is indeterminate.

### Hot state, retained archive, and tombstone indexes

The v2 conformance snapshot divides capacity, not authority:

- **active capacity** keeps the existing 256 ceilings for unexpired registrations, usable
  approvals, and attempts not durably destroyed with cleanup false;
- **hot-retained capacity** keeps the existing 4,096 ceilings for registrations, approvals,
  attempts, and lifecycle records present in the mutable snapshot; and
- **archive-retained capacity** contains complete immutable closed cohorts plus exact tombstone
  indexes. Archived records remain authoritative retained state, not evicted state.

The active v2 snapshot retains two non-interchangeable index domains. The `retained-global`
projection contains every hot and archived entry below. Each entry has a typed `RecordLocation`:
`hot` plus one positive canonical hot-record ordinal, or `archive` plus positive segment/cohort/
record ordinals. The location also binds its `RecordKind`; unused arms must be zero. A segment's
`segment-derived` projection is archive-only and cannot substitute for the global projection.

The retained-global indexes are sorted and duplicate-free:

| Index | Exact retained projection |
| --- | --- |
| registration | `RegistrationID`, sequence, plan digest, expiry, registration-record location, full-record digest |
| approval | `ApprovalID`, registration, payload digest, authorization identity, nonce, state, consumed `AttemptID`, expiry, approval-record location, full-record digest |
| attempt | `AttemptID`, approval, registration, created time, explicit absent/present lifecycle union, attempt-record location, full-record digest; the present arm also binds lifecycle disposition, lifecycle-record location, and full-record digest |
| nonce | `AttemptNonce` to approval payload digest, `ApprovalID`, and approval-record location |
| effect | every v2-issued nonzero `EffectID` to `AttemptID`, operation sequence, operation, issuance snapshot generation, attempt-record location, and explicitly labeled visible-v1 seed membership |
| instance | instance identity digest to `AttemptID` and lifecycle-record location |
| approval replay | canonical payload digest plus exact-payload digest, authorization identity, `ApprovalID`, state, and approval-record location |
| attempt replay | `(RegistrationID, ApprovalID)` to the one committed `AttemptID`, state, and attempt-record location |

The exact payload bytes remain in the immutable segment. Approval replay follows the existing
oracle: a verified canonical payload and the same retained authorization identity return the same
approval and current fixed state; a digest collision with different exact payload bytes is
`REPLAY`, not a second record. Attempt replay returns the same attempt. Archived registrations are
expired and cannot create new approval or attempt authority, but read-only exact lookup may load
their retained bytes from the segment.

Every new registration, approval, attempt, nonce, effect ID, and instance identity checks both hot
records and archive tombstones in the same logical transaction. A v2 effect intent inserts its
tombstone in that transaction and refuses if the ID exists anywhere in the visible seed/v2 set.
Archive location is never
authority by itself; the referenced full record and every digest/cross-link must verify.

For every count family, `hot + archived = total`; retained-global index counts equal total; entries
on the hot location arm equal hot; entries on the archive arm equal archived; and summed
referenced-descriptor counts equal archived. A visible-v1 effect attached to a hot lifecycle is
therefore counted in hot and total effects while archive descriptors and archived effects remain
zero. It never inflates an archive cohort or descriptor count.

### Closed archive formats and bounded checkpoint oracle

The conformance oracle introduces an explicit v2 active snapshot and immutable archive-segment v0.
Both use closed JSON only because the current fixed-store oracle does; this does not select JSON for
a production store or a signed protocol object.

The active v2 snapshot adds:

- store format version `2` and migration source version `1`;
- positive snapshot and archive generations;
- durable time high water, installation, Supervisor, and epoch identity;
- the complete hot registration/approval/attempt/lifecycle sets and their existing set digests;
- a sorted archive-segment descriptor set and descriptor-set digest;
- the scope-tagged retained-global indexes above and one digest per index;
- effect-tombstone coverage `visible-v1-seed-plus-all-v2-issued`, including the migration seed
  count and digest;
- a combined archive-index digest;
- kinded prior and current checkpoint references; and
- exact hot, archived, and total counts per record and tombstone family.

Each immutable archive segment contains:

- segment format/version and positive ordinal;
- installation, Supervisor, epoch, source snapshot generation, archive generation, and durable time
  high water;
- prior archive-checkpoint digest;
- a sorted complete cohort set;
- separate registration, approval, attempt, lifecycle, nonce, effect, instance, and replay set
  digests;
- exact counts and encoded byte length; and
- a domain-separated digest over every preceding field and the canonical segment bytes with the
  digest field omitted.

The migration-genesis checkpoint is a distinct generation-one representation domain separated as
`capsule.supervisor.archive-migration-genesis.v0`. It requires store/source versions 2/1, result
snapshot and archive generations 1/1, the empty descriptor-set digest, every complete all-hot
retained-global index digest, all hot set digests, exact visible-v1 seed count/digest, hot counts,
installation/Supervisor/epoch identity, and durable time high water. It has no previous checkpoint,
new segment, or archive location. Its retained all-present generated known answer is
`657c86aff68354e535369bdecf9e8c23bbfb1457e8ada57576bc1a52666a5fc9`. The valid
missing-lifecycle fixture separately retains genesis
`0af76dca782bdf198d5ef80b6b2856fb35ae01539e7d8866198c4f2af643f621` with exact counts
`attempts = 1, lifecycles = 0`.

An activation checkpoint is separately domain separated as
`capsule.supervisor.archive-checkpoint.v0` and binds the kinded previous checkpoint reference, new segment
digest, complete descriptor-set digest, complete archive-index digest, all hot set digests, source
and resulting snapshot generations, archive generation, installation/Supervisor/epoch identity,
and durable time high water. It requires a new segment and a strictly increasing snapshot
generation. Checkpoint references bind `migration-genesis` or `activation` kind so one
representation cannot be interpreted as the other. A checkpoint detects accidental omission or mix-and-match inside the visible
archive set; it is not an anti-rollback anchor or signature.

The exact unwired fixed-store caps are:

| Dimension | Inclusive maximum |
| --- | ---: |
| active v2 snapshot bytes | 64 MiB |
| one archive segment bytes | 64 MiB |
| cohorts selected per activation | 256 |
| registrations per segment | 256 |
| approvals per segment | 4,096 |
| attempts per segment | 4,096 |
| lifecycle records per segment | 4,096 |
| referenced archive segments | 64 |
| entries in any complete archive tombstone index | 262,144 |

Every limit is checked before a record leaves the hot snapshot. The first reached limit returns
`CAPACITY` with no archive activation, hot-record removal, eviction, or external effect. These are
conformance bounds, not production sizing or a performance claim. At 64 full 64 MiB segments the
fixed oracle may retain roughly 4 GiB and require an expensive full verification; tests use compact
generated populations rather than allocating the maximum. Indefinite continuous service is not a
claim of this design.

### Archive creation, verification, and activation

One archive activation occurs under the installation owner lock and one in-process archive
coordinator:

1. Reopen and fully validate the existing v2 active snapshot and every referenced segment. Refuse
   while recovery, repair, quarantine, transition, or ownership is unresolved.
2. Deterministically select a complete eligible cohort batch within every count and byte budget.
   Selection creates no mutation.
3. Build the complete candidate segment and candidate post-activation hot state/indexes in memory.
   Recompute every record, set, descriptor, index, and checkpoint digest from source records.
4. Revalidate the candidate as an independent closed world: no missing/duplicate cross-link, no
   active record selected, no hot/archive overlap, no tombstone omission, no identifier collision,
   and no capacity overflow.
5. Write the segment to a same-volume mode-0600 temporary regular file; sync and close it; reopen
   without symlink following; bound-read; closed-decode; recompute all digests/index projections;
   and require byte-for-byte equality with the prepared candidate.
6. Rename the verified segment to a digest-addressed final name and sync the archive directory.
   At this point it is immutable but inactive.
7. Build the complete v2 active snapshot that references the final segment, installs the new
   tombstone indexes/checkpoint, and removes the exact cohort records from hot sets. Recompute and
   validate every hot and archive digest/cross-link.
8. Commit the active snapshot through temp-file sync, atomic rename, and directory sync.
   Activation is the active-snapshot rename, never segment existence alone.
9. Reopen the active snapshot and all referenced segments and run full verification before
   returning success.

The segment is published and directory-synced before an active snapshot may reference it. The hot
records are removed only in the same complete active-snapshot commit that installs their segment
descriptor and tombstones. There is therefore no committed state in which an archived cohort is
neither hot nor reachable.

An unreferenced prepared/final segment grants no authority and is ignored by normal lookup. It may
be removed only by the narrow orphan rule below.

### Crash, power-loss, and indeterminate outcomes

Confirmed failures before segment rename preserve the old active snapshot and create no referenced
archive. Confirmed failures after segment publication but before active-snapshot rename leave the
old hot state authoritative and one unreferenced segment. They do not release capacity.

An indeterminate segment-directory sync or active-snapshot rename/directory-sync immediately
fences all mutation. Reopen must establish one of two complete worlds:

- old active snapshot, cohort still hot, and any new segment unreferenced; or
- new active snapshot, cohort absent from hot sets, referenced segment and every tombstone present.

Any third state is `repair-required`. A referenced missing/corrupt segment never falls back to the
old hot snapshot or an earlier checkpoint because that could resurrect approval or lose cleanup
history. Hot/archive duplication under a checkpoint that references both is also corruption, even
if bytes match.

The fixed oracle can inject confirmed pre-rename and indeterminate post-rename outcomes. Those tests
model ordering; they are not real power-loss evidence. Consumer or production activation requires
the selected engine and filesystem to pass abrupt process death and real APFS/power-interruption
tests at each write, sync, rename, directory-sync, checkpoint, and reopen boundary.

### Compaction and deletion eligibility

For this ADR, **compaction means only moving a complete eligible cohort from the mutable hot
snapshot into a verified immutable segment while retaining the full cohort and exact tombstones**.
It does not mean erasing security history or merging unknown states.

No referenced archive segment, full cohort, registration, approval, attempt, nonce, effect,
instance, replay tombstone, checkpoint, or sole explanatory failure record is eligible for
deletion under this proposal. Capacity pressure never changes that rule.

Only these bytes may be removed automatically, under the owner lock and after a full reference
scan:

- a same-volume temporary file whose name and inode were created by the current failed operation;
- a digest-addressed segment that is not referenced by the current active snapshot or any coherent
  backup manifest and whose digest exactly matches a locally recomputed, fully decoded candidate;
  and
- a superseded temporary backup/index file that was never activated.

Unknown, malformed, wrongly owned, linked, symlinked, or digest-mismatched files are quarantined,
not deleted. Secure deletion is not claimed on copy-on-write filesystems, SSDs, snapshots, or
backups.

Deleting referenced archive history requires a later ADR that defines retention policy, receipt
and trust dependencies, exact replay behavior after body retirement, independently protected
identifier generations or equivalent non-reuse evidence, backup/witness interaction, and honest
secure-deletion limitations. Until then, reaching the segment/index cap refuses new authority.

### Backup, restore, rollback, and identifier non-reuse

Ordinary file copy is not a coherent backup. A coherent backup operation must:

1. hold the installation owner lock and disable mutation;
2. fully verify the active snapshot and every referenced segment;
3. copy the active snapshot, every referenced segment, and a closed backup manifest into a new
   destination without following links;
4. bind installation, Supervisor, epoch, active/archive generations, durable time high water,
   checkpoint head, exact file names, lengths, and digests;
5. sync all files and directories; and
6. independently reopen and fully verify the copied set before calling it complete.

No online fuzzy backup is selected. The first conformance implementation may create and verify a
backup set but may not activate restore.

Restore is offline and owner-lock held. It first verifies the complete candidate set. If an
independently protected latest checkpoint exists, the restored checkpoint, generations, epoch, and
time high water must equal that anchor. A lower, missing, or incomparable checkpoint refuses.

Without such an anchor, a cryptographically and structurally valid backup may still be a coherent
older world. It may be inspected offline but cannot silently enable attempts. It enters a fixed
`rollback-uncertain`/`repair-required` state and requires a separately authorized forward repair
and trust-epoch transition. This ADR does not define that repair ceremony.

Within the visible archive history, hot records plus tombstone indexes prevent reuse of every
registration, approval, attempt, nonce, retained effect, and instance identity. A coherent rollback
can hide later random identifiers and permit accidental reuse against the invisible history.
Therefore this boundary does not claim durable global non-reuse, rollback prevention, or monotonic
history without an independent checkpoint/non-rollbackable anchor. Randomness alone does not close
that requirement.

### Migration and version refusal

The fixed oracle uses one explicit offline, lock-held `v1 -> v2` migration:

1. fully validate v1 under its current 64 MiB bound;
2. require no recovery/repair/quarantine/transition fence;
3. construct v2 with an empty descriptor set, zero archived counts, generation one, one complete
   all-hot `retained-global` projection reconstructed from every v1 registration/approval/attempt/
   lifecycle record, and one dedicated migration-genesis checkpoint; represent lifecycle as absent
   or present strictly from the decoded v1 lifecycle collection and derive lifecycle counts
   independently of attempt counts; seed registration/approval/
   attempt/nonce/replay identities from all v1 hot records and seed the effect tombstone set from
   only the nonzero effect IDs still visible in v1 lifecycle records;
4. prove authority and lifecycle sets are byte-for-byte and digest-identical to v1 and record the
   exact visible-v1 effect seed count/digest and limited coverage label;
5. require retained-global counts and location arms to equal hot/total counts while every segment-
   derived/descriptor/archived count remains zero, then sync/rename/sync once; and
6. reopen and fully validate v2, including recomputation of the migration-genesis checkpoint.

Failure before rename preserves v1. An indeterminate rename requires reopen. Old binaries refuse
v2; the v2 opener refuses v0/v1 and never interprets a missing archive collection as empty. Unknown
active, segment, index, or tombstone versions are `repair-required` and never skipped, normalized,
or automatically migrated. A segment format migration is a separate explicit archive rewrite and
is not authorized by this ADR.

### Repair and quarantine

These conditions enter installation `repair-required`, disable attempts, retain original bytes,
and make zero adapter calls:

- unsupported version; malformed, truncated, trailing, duplicate, or unknown data;
- file type, owner, mode, link, or no-follow violation;
- missing referenced segment or checkpoint;
- record, set, descriptor, index, checkpoint, length, or file-name digest mismatch;
- hot/archive overlap or omission;
- duplicate identifier, nonce, payload replay key, effect ID, or instance identity;
- cross-registration cohort split or cross-link mismatch;
- generation, ordinal, durable-time, installation, Supervisor, or epoch mismatch; or
- a referenced record that was not archive-eligible at activation.

The Supervisor never repairs by dropping a segment, rebuilding indexes from whichever records are
easiest to parse, restoring an older active snapshot, or recreating a missing store. Offline
verification may produce a bounded report and preserve a byte-for-byte forensic copy. Actual
repair requires a later authenticated procedure that preserves or explicitly replaces trust,
grant, attempt, cleanup, replay, and evidence history.

### Owner-lock dependency

Every migration, archive preparation, activation, orphan removal, coherent backup, restore check,
and offline repair operation requires the same lifetime installation owner lock as ordinary
Supervisor state mutation. Per-attempt coordinators are insufficient because an archive cohort may
contain many attempts and updates global indexes.

The current injected owner/session mechanics are not that lock. Proposed ADR-0033 selects the
pre-created enrolled sibling plus BSD `flock` mechanism after duplicate-process, process-death,
replacement, and hostile-file tests in an owned temporary harness. Archive behavior remains limited
to an unwired single-process harness until G2 composes the passive Go/Darwin owner with the store
and the installed protected-root/session/update/reboot matrix passes. A second owner refuses before
reading a mutable candidate, creating a segment, changing an index, or calling any adapter.

### Offline verification

The repository will define a read-only full verifier that accepts only a trusted installation-
configured store root, never an agent/Broker/backend path. It performs no mutation and emits a
bounded fixed-shape report containing format versions, installation/Supervisor/epoch identity,
generations, time high water, checkpoint head, file/record counts, digest verdicts, and one fixed
failure classification. It emits no plan, approval, path, guest, user-content, or arbitrary stored
text.

Full verification closed-decodes and recomputes all active and segment records, set/index/checkpoint
digests, counts, sort order, eligibility facts, and cross-links. A faster startup mode is not
selected by this ADR. Any later incremental verification must be justified by an authenticated or
otherwise protected checkpoint and exact invalidation rules; cached success alone is not proof.

### Internal API boundary

The following names are implementation targets, not public or IPC contracts:

```go
type ArchiveLimits struct {
    MaxCohorts uint16 // fixed to 256 in the conformance oracle
    MaxBytes   uint64 // fixed to 64 MiB in the conformance oracle
}

type ArchivePlan interface { sealedArchivePlan() }
type PreparedArchive interface { sealedPreparedArchive() }
type VerifiedArchive interface { sealedVerifiedArchive() }

type SupervisorArchiveStore interface {
    PlanArchive(context.Context, ArchiveLimits) (ArchivePlan, error)
    PrepareArchive(context.Context, ArchivePlan) (PreparedArchive, error)
    VerifyPreparedArchive(context.Context, PreparedArchive) (VerifiedArchive, error)
    ActivateArchive(context.Context, VerifiedArchive) (ArchiveCheckpoint, error)
    VerifyArchiveSet(context.Context, VerificationMode) (ArchiveVerificationReport, error)
    CreateCoherentBackup(context.Context, BackupDestination) (BackupManifest, error)
    VerifyCoherentBackup(context.Context, BackupManifest) (ArchiveVerificationReport, error)
}
```

`ArchivePlan`, `PreparedArchive`, and `VerifiedArchive` are sealed store-issued values bound to the
source snapshot generation, archive generation, checkpoint head, selected cohort digests, and owner
session. They cannot contain caller-supplied records or paths. Activation rejects a stale plan or a
different owner session.

Authority operations use internal hot/archive lookup and uniqueness ports; they do not expose an
archive operation to callers:

```go
type RetainedAuthorityLookup interface {
    ResolveRegistration(context.Context, RegistrationID) (RetainedRegistration, error)
    ResolveApprovalReplay(context.Context, ApprovalPayloadDigest, []byte,
        ApprovalKeyAuthorizationIdentity) (ApprovalReference, ApprovalState, error)
    ResolveAttemptReplay(context.Context, RegistrationID, ApprovalID)
        (AttemptReference, AttemptState, error)
    CheckIdentifierSet(context.Context, CandidateIdentifiers) error
}
```

`Drive`, `Recover`, and startup recovery remain `AttemptID`-only and enumerate only hot active work.
Archived cohorts are terminal and never sent to the lifecycle adapter.

## Engine comparison and selection boundary

### Minimal fixed-store checkpoint

Strengths for the next conformance slice:

- preserves the exact E5 full-state validation and before/after-rename fault oracle;
- makes segment-publication-before-activation ordering visible and testable;
- supports byte-exact fixtures and deliberate corruption without a database-specific tool; and
- ensures the logical archive protocol is specified before an engine can hide it behind
  transactions.

Security and operational costs:

- every active mutation rewrites a bounded full snapshot;
- archive activation spans immutable segment publication and a second active-snapshot commit;
- startup/full verification scales with all retained segments;
- the index itself grows and eventually reaches a hard cap;
- there is no production multi-process lock, online backup, page-level recovery, encryption, or
  real power-loss evidence; and
- a 64-segment fixed oracle is explicitly finite and unsuitable for continuous service.

### SQLite production candidate

SQLite could place hot records, closed cohorts, tombstones, indexes, generations, and checkpoint
metadata in one transactional database. Foreign keys, uniqueness constraints, `WITHOUT ROWID`
tables, and the backup API may reduce custom split-file state. WAL or rollback-journal behavior can
also avoid rewriting the complete active population for each ordinary transaction.

Those are hypotheses, not selection evidence. A production choice must pin and test:

- SQLite version/build options, page size, locking mode, `journal_mode`, `synchronous`, macOS full-
  fsync behavior, WAL checkpoint policy, and directory/file protection;
- `BUSY`, `LOCKED`, `FULL`, `IOERR`, `CORRUPT`, partial page/WAL, stale reader/writer, checkpoint,
  and process/power-loss outcomes;
- schema migration, downgrade refusal, backup/restore, integrity checking, archive export/offline
  verification, and bounded growth/vacuum behavior; and
- the same cohort eligibility, replay, uniqueness, cleanup, repair, rollback, and no-resurrection
  oracles in this ADR.

SQLite does not by itself provide coherent-rollback detection, identifier non-reuse across an
older restored world, secure deletion on APFS/SSD/snapshots, authenticated checkpoints, or correct
archive semantics. Choosing it merely because it offers transactions would leave the security
decision incomplete.

### Decision

Use the minimal fixed-store checkpoint for passive and fault-injectable archive conformance only.
Keep the archive model and lookup semantics engine-neutral. Defer production engine selection until
the fixed oracle exists and SQLite (or another named candidate) runs the same logical corpus plus
real locking, backup, corruption, APFS, and power-loss tests. No consumer may activate on the fixed
archive oracle.

## Consequences

- Closed terminal cohorts can leave the mutable retained sets without becoming absent: the full
  records and exact replay/non-reuse tombstones remain authoritative and verifiable.
- Active and hot-retained capacity become distinct from total retained history. Archive activation
  can release hot-retained capacity but never active cleanup capacity.
- Exact replay and identifier uniqueness consult hot and archived state, preventing compaction from
  resurrecting approvals or reissuing identities within visible history.
- The two-file publication boundary is explicit: segment first, active reference second, reopen on
  every indeterminate outcome.
- No referenced security history is deletable under this proposal. Storage is hard-bounded and
  eventually refuses, so indefinite continuous service remains unsupported.
- Coherent backup is a verified set operation; restore without an independent latest checkpoint is
  rollback-uncertain and cannot silently enable attempts.
- The fixed checkpoint remains an unwired local mechanic. It advances no production, evidence,
  backend, runtime, guest, or continuous-service claim.

## Alternatives considered

### Put lifecycle disposition directly on every attempt without an absence arm

Rejected. V1 deliberately commits an immutable attempt before lifecycle establishment. Requiring a
lifecycle state would either omit retained attempt identity or invent `prepare-pending` and cleanup
state. The selected closed absent/present union preserves the attempt and counts only decoded
lifecycle records.

### Add a separate retained-global lifecycle index

Not selected. It can be safe with a ninth index/digest, exact AttemptID join, typed locations,
independent counts, caps, checkpoint bindings, and regenerated fixtures. The selected union binds
the same one-to-zero-or-one relationship and lifecycle record anchor within the existing attempt
index, so it is the narrower passive correction.

### Restrict migration to attempts with lifecycle after a recovery ceremony

Rejected as the F2 mapping. A real v1 lifecycle-establishment transaction could make a particular
store satisfy that precondition, but requiring it adds backend-binding, capacity, authorization,
and confirmed/indeterminate ceremony boundaries and makes a supported crash state only
conditionally migratable. Migration instead maps absence exactly and invokes no adapter or
lifecycle mutation.

### Delete terminal records and keep only set digests

Rejected. A digest cannot answer exact payload replay, return the original approval/attempt
identity, prevent nonce/effect/instance reuse, validate a cross-link, or explain a terminal
failure. It also cannot distinguish omission from an intentionally empty set without retained
structure.

### Archive each attempt independently

Rejected. Registration expiry and replay semantics span every approval and attempt for one
registration. Attempt-sized units would duplicate or split exact plan/approval state and create a
cross-hot/archive saga for later approvals. A closed registration cohort has no remaining authority
edge back into active work.

### Treat expiry as permission to forget approvals or registrations

Rejected. Expiry prevents new use; it does not erase payload replay, nonce uniqueness, consumed
state, attempt identity, or explanatory history. Clock or restore rollback must not make forgotten
authority usable.

### Put archive segments beside v1 without a new active format

Rejected. V1 has no authenticated/validated reference, tombstone index, generation, or checkpoint
for archived state. Interpreting missing v1 records through an optional sidecar would make absence
ambiguous and permit old binaries to ignore archive history.

### Select SQLite immediately

Rejected for this boundary, not rejected as a production candidate. Engine migration before the
archive semantics and two-world fault oracle exist would conflate logical correctness with engine
behavior and still leave power loss, backup, rollback, deletion, and non-reuse unresolved.

### Append forever with no hard cap

Rejected. It hides availability and disk-exhaustion behavior. The fixed oracle has exact count and
byte caps and returns `CAPACITY` without authority change when retained history reaches them.

## Conformance plan

The exact passive types, migration, archive activation, replay, fault, backup, and offline-
verification slices are defined in
[the Supervisor archive/compaction conformance plan](../SUPERVISOR_ARCHIVE_COMPACTION_PLAN.md).
The retained
[valid-v1 mapping resolution](../SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md) selects the passive
missing-lifecycle representation and rejects a narrower state-changing migration ceremony. The
[stateful F2 migration/full verifier](../SUPERVISOR_ARCHIVE_F2_MIGRATION_RESULT.md) is now `PASSED`
in its exact local fixed-store scope. The
[stateful F3 first-segment activation](../SUPERVISOR_ARCHIVE_F3_ACTIVATION_RESULT.md) is also
`PASSED` in its exact local fixed-store scope. This ADR remains Proposed: F3 does not implement
retained lookup, v2 authority mutation, a second segment, backup/orphan cleanup, a production
engine, or a product consumer.

## Acceptance blockers

This ADR remains Proposed and consumer activation remains blocked until all applicable work exists:

1. the fixed-store conformance slices and byte-exact/corruption/fault corpus in the linked plan;
2. the ADR-0033 Go/Darwin owner-lock port and installed protected-root hostile-file/session/update
   evidence;
3. a separately selected production engine with real APFS/power-loss, locking, backup, migration,
   corruption, and quantitative performance evidence;
4. an independently protected latest checkpoint/non-rollbackable anchor, or an explicit product
   posture that keeps rollback-uncertain restore disabled;
5. reviewed retention, evidence/export, and secure-deletion policy before any referenced archive
   body or tombstone can be removed; and
6. production approval verification, authenticated IPC, trust/update integration, evidence
   composition, consumers, and all existing runtime/backend/content gates.

No blocker may be inferred complete from the fixed-file or no-guest fake evidence.
