# Supervisor archive, compaction, and replay-retention conformance plan

Status: proposed and local-only. Slice F1 implements passive archive projections, exact
limits/domain-separated known answers, defensive copies, and the pure closed-cohort eligibility
selector. The passive F2 format-correction slice resolved PR #78's three contradictions with
scope-separated global/segment indexes, typed hot/archive locations and count equations, and a
distinct migration-genesis checkpoint. Generated answers and before/after semantics are retained in
[the F2 blocker resolution](SUPERVISOR_ARCHIVE_F2_FORMAT_BLOCKER.md). Proposed ADR-0031 still
defines the unimplemented immutable retained archive
and selects a minimal fixed-store checkpoint only as the conformance oracle. The follow-on valid-v1
mapping contradiction is now passively resolved with an explicit absent/present lifecycle union on
attempt entries and independent lifecycle counts. The exact real-v1 witness, alternative analysis,
decision, and fault plan are retained in the
[F2 v1 mapping resolution](SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md). The
[stateful F2 result](SUPERVISOR_ARCHIVE_F2_MIGRATION_RESULT.md) now implements the owner-asserted
v1-to-v2 migration, downgrade refusal, and read-only empty-archive full verifier with exact known
answers and local fault/corruption/capacity/concurrency/process-death evidence. It implements no
retained lookup, consumer, IPC, runtime, backend, service, identity, credential, user data,
deployment, or guest. The follow-on
[stateful F3 result](SUPERVISOR_ARCHIVE_F3_ACTIVATION_RESULT.md) now implements exactly one sealed
immutable-segment prepare/verify/publish/activate transaction. It retains the complete cohort and
all visible identities/tombstones, publishes before reference, atomically installs the generation-
two checkpoint, reports valid orphans without deleting them, and fully reopens either the complete
predecessor or successor under the owner assertion. The follow-on
[F4A result](SUPERVISOR_ARCHIVE_F4A_LOOKUP_RESULT.md) now implements only read-only retained-global
lookup, replay, passive collision routing, and hot-only `AttemptID` recovery. The retained
[F4B blocker](SUPERVISOR_ARCHIVE_F4B_MUTATION_BLOCKER.md) records the former effect-history gap;
the [F4B result](SUPERVISOR_ARCHIVE_F4B_MUTATION_RESULT.md) now passes atomic authority/lifecycle
mutation and the independent append-only effect source in the fixed-store scope. The
[F4C result](SUPERVISOR_ARCHIVE_F4C_GROWTH_RESULT.md) now passes deterministic second/later
immutable-segment activation and exact segment 64/65 bounded growth. The
[F5 result](SUPERVISOR_ARCHIVE_F5_BACKUP_RESULT.md) now passes owner-held coherent backup,
read-only restore admission, explicit known-orphan handling, and fixed-shape offline reporting.
F6 and production admission remain outside the passed scope.

Normative proposal:
[ADR-0031](adr/0031-checkpoint-closed-supervisor-cohorts.md).

## Objective and defensive scope

Defensively validate Capsule's no-resurrection, replay-retention, cleanup, and archive-activation
ordering using only the repository's current v1 fixed store, authority/lifecycle types, retained
fixtures, compact generated populations, and local tests. Do not access any runtime, backend,
guest, process/service, identity, credential, user data, deployment, or other system.

Preserve throughout:

- atomic approval consumption plus immutable attempt creation before any effect;
- `Drive(AttemptID)`, `Recover(AttemptID)`, and startup recovery by `AttemptID` only;
- no lifecycle-active, cleanup-bearing, unresolved, quarantined, repair-required, or recovery-
  fenced record enters an archive;
- no automatic eviction or referenced-history deletion;
- no approval, registration, nonce, attempt, effect, or instance resurrection or reuse;
- confirmed versus indeterminate commit distinctions; and
- `FakeBackend.CreatesGuest() == false` with zero adapter calls from archive code.

## Implementation oracle and non-goals

The stateful implementation oracle is current Slice B/E5:

- `registrationstate.FixedFileStoreV1` closed-decodes and validates the complete authority and
  lifecycle snapshot;
- registration, approval, attempt, and lifecycle set digests and cross-links are recomputed on
  reopen;
- the fixed v1 store is bounded at 64 MiB with 256 active and 4,096 retained lifecycle records;
- only `destroyed` plus cleanup false after authoritative absence releases active capacity;
- indeterminate rename/directory-sync outcomes fence the open store until reopen; and
- the only lifecycle adapter is the no-guest fake.

Slice F1 plus its passive format correction add `internal/execution/archivestate` values and tests.
They implement closed generation/ordinal/digest/index/descriptor/checkpoint/plan projections,
scope- and kind-separated record locations, exact candidate limits and generated known answers,
defensive-copy behavior, and deterministic complete-cohort selection with `RecoveryAttemptIDs`
exclusion. Stateful F2 uses those projections to migrate a fully validated v1 file into one all-hot
v2 file and fully verify it on reopen. It does not write a segment, release hot capacity, route
lookup, mutate v2 authority, or invoke the fake.

E5 retains only the current lifecycle record's `EffectID`; later operations replace that field.
Slice F2 records this as a migration limitation and seeds tombstones from every nonzero ID still
visible in v1. F4B must add v2 mutation behavior so `BeginEffect` commits
the new effect tombstone in the same transaction as the intent. Tests must not imply that
overwritten pre-v2 fake IDs were reconstructed or that a transcript exists.

This work does not:

- accept ADR-0019, ADR-0024, ADR-0025, ADR-0029, or ADR-0031;
- select SQLite or another production engine;
- implement the selected ADR-0033 macOS owner lock or protected production storage;
- authorize archive-history deletion, secure deletion, online backup, or restore activation;
- establish rollback prevention or identifier non-reuse across a coherent older restore;
- compose enforcement evidence or receipts; or
- advance any runtime/backend/guest or continuous-service claim.

## Exact conformance constants

The fixed oracle introduces these exact candidate constants:

```go
const (
    SupervisorStoreFormatV2       = uint64(2)
    SupervisorArchiveFormatV0     = uint64(0)
    MaxSupervisorStateV2Bytes     = int64(64 << 20)
    MaxSupervisorArchiveBytes     = int64(64 << 20)
    MaxArchiveCohortsPerSegment   = 256
    MaxArchiveApprovalsPerSegment = 4_096
    MaxArchiveAttemptsPerSegment  = 4_096
    MaxArchiveLifecyclesPerSegment = 4_096
    MaxReferencedArchiveSegments  = 64
    MaxArchiveIndexEntries        = 262_144
)
```

Names may change during implementation, but values and cap-plus-one oracles must not drift without
updating ADR-0031 and the retained fixtures. Byte caps apply before ordinary decode or allocation.
Count and byte limits are inclusive; the first exceeded dimension returns `CAPACITY` without state
change.

## Passive formats and exact digests

### V2 active snapshot

The v2 envelope is closed and contains exactly:

- `storeFormatVersion = 2` and `migrationSourceVersion = 1`;
- positive `snapshotGeneration` and `archiveGeneration`;
- the existing installation state and complete hot authority collections;
- the existing lifecycle collection and lifecycle-set digest;
- sorted archive descriptors and descriptor-set digest;
- one scope-tagged `retained-global` registration, approval, attempt, nonce, effect, instance,
  approval-replay, and attempt-replay projection plus one digest per index; every entry binds a
  typed record location whose `hot` and `archive` arms are mutually exclusive; every attempt also
  binds an explicit lifecycle `absent` arm or a `present` arm with lifecycle state, separately typed
  lifecycle-record location, and full-record digest;
- an effect-tombstone coverage tag fixed to `visible-v1-seed-plus-all-v2-issued`, seed count, and
  seed digest;
- one combined retained-global index digest;
- prior/current kinded checkpoint references; and
- exact hot/archive/total counts for every collection/index.

All arrays are non-null even when empty. Unknown, missing, duplicate, trailing, nonpreferred, or
wrong-version data refuses. `segment-derived` index digests are domain-separated from
`retained-global` digests and are valid only inside one archive segment. The exact descriptor,
index, visible-v1 seed, migration-genesis, and activation-checkpoint answers are generated into the
retained fixture named by the blocker resolution.

### Archive segment v0

The segment is closed and contains exactly the metadata, cohort collections, set digests, counts,
encoded length, and segment digest specified by ADR-0031. Cohorts sort by registration sequence and
then registration ID. Within a cohort:

- approvals sort by `ApprovalID`;
- attempts sort by `AttemptID`;
- lifecycle records sort by `AttemptID`; and
- derived nonce/effect/instance/replay projections sort by raw key bytes and then their complete
  disambiguating tuple.

Digest construction uses length-prefixed fields and domain-separated SHA-256. JSON concatenation,
map iteration order, path strings, file names, or signature bytes never define semantic identity.
The implementation adds byte-exact empty, one-cohort, multi-attempt-cohort, and maximum-boundary
known answers.

### Index reconstruction oracle

Full verification discards serialized retained-global indexes in memory, reconstructs them from all
decoded hot and archived records, and requires exact equality with the serialized sorted
projections and their scope-separated digests. Every reconstructed entry binds the expected record
kind plus exactly one location arm: a positive ordinal in the canonical sorted hot record set, or
positive segment/cohort/record ordinals in a referenced archive. The reconstruction must detect:

- one identifier at two locations;
- one nonce bound to two payloads;
- one approval payload replay key bound to two approval IDs;
- one `(RegistrationID, ApprovalID)` bound to two attempts;
- one effect ID bound to another attempt/operation sequence;
- one instance digest bound to another attempt;
- a tombstone without a full record;
- a full record without every required tombstone; and
- hot/archive overlap or a cross-cohort reference.

Attempt/lifecycle reconstruction is a total one-to-zero-or-one join. Every attempt is indexed.
Only a decoded lifecycle record may create a present arm; it must use the same hot/archive world as
its attempt and the same segment/cohort when archived. An absent arm carries no state, location, or
digest, remains hot active startup work, and cannot have effect or instance entries. Lifecycle
counts and location counts derive only from present arms.

For every family it also requires:

```text
hotCounts + archivedCounts = totalCounts
retainedGlobalIndexes.counts = totalCounts
retainedGlobalIndexes.hotLocationCounts = hotCounts
retainedGlobalIndexes.archiveLocationCounts = archivedCounts
sum(referencedDescriptor.counts) = archivedCounts
```

Thus the migration visible-v1 effect seed is a hot retained-global effect entry and does not
increase archive descriptors, segment-derived entries, or archived counts.

A hash collision does not merge records. Exact bytes and complete typed bindings decide whether a
case is idempotent replay or `REPLAY`/repair-required.

For every v2 `BeginEffect`, the candidate effect ID is checked against the complete visible-v1
seed plus v2 hot/archive effect tombstones, and the tombstone commits atomically with the intent.
Confirmed/indeterminate intent behavior remains the E5 oracle: no adapter call occurs before that
transaction, and reopen finds neither ID/intent or both ID/intent. A tombstone without its issuing
intent/archived lifecycle cross-link is repair-required.

### Checkpoint kinds

Migration creates `MigrationGenesisCheckpoint` under
`capsule.supervisor.archive-migration-genesis.v0`. It is the only constructor that accepts result
snapshot generation one. It requires archive generation one, no previous checkpoint, no segment,
the empty descriptor digest, all-hot retained-global indexes, exact hot set digests/counts, the
sorted visible-v1 seed/count/digest, and installation/Supervisor/epoch/time bindings.

Archive activation creates `ArchiveCheckpoint` under
`capsule.supervisor.archive-checkpoint.v0`. It requires a kinded genesis-or-activation predecessor,
one nonzero new segment digest, descriptor/index/hot-set digests, and a result snapshot generation
strictly greater than its source. Active-state references bind checkpoint kind as well as digest;
wrong-kind substitution refuses. The generated fixture retains literal answers for both domains.

## Eligibility oracle

Implement one pure selector over a fully validated snapshot:

```go
func SelectClosedRegistrationCohorts(
    state ValidatedV1OrV2State,
    limits ArchiveLimits,
) (ArchivePlan, error)
```

The returned plan is sealed, defensively copied, and binds the source snapshot generation, archive
generation, checkpoint head, owner session, every selected record digest, and candidate limits.

Focused selector cases cover:

| Cohort state | Result |
| --- | --- |
| expired registration, no approvals | eligible |
| registration at exact expiry high water | eligible if every other condition holds |
| registration one second before expiry | not selected |
| expired durable-`usable` approval, no attempt | eligible; state is not rewritten |
| consumed approval plus destroyed/absent attempt | eligible |
| invalidated approval, no attempt | eligible |
| any usable unexpired approval | not selected |
| missing approval/attempt/lifecycle half | repair-required, not merely ineligible |
| observed/stopped/destroy-confirmed lifecycle | not selected |
| unresolved/quarantined/exhausted/recovery-fenced lifecycle | not selected |
| destroyed with cleanup true or non-authoritative absence | invalid/repair-required |
| cohort present in `RecoveryAttemptIDs` | not selected and test fails if selector disagrees |
| trust transition, quarantine, repair, or recovery fence | archive operation refuses globally |
| one registration with two independently approved destroyed attempts | one indivisible eligible cohort |
| one active and one destroyed attempt for a registration | entire cohort remains hot |
| candidate count or byte cap reached | deterministic prefix only; no partial cohort |

Selection makes no write and no fake call. Repeating selection against the same generation and
limits produces identical cohort identities and order.

## Archive creation and activation state machine

The archive coordinator uses these passive states:

```text
planned -> prepared -> verified -> segment-published -> activated -> reopen-verified
          \-> refused/recovery-fenced/repair-required
```

Only store-issued sealed values cross steps. A stale source generation, checkpoint head, owner
session, cohort digest, or limit rejects before a write.

### Creation/activation oracle

| Boundary | Active snapshot | Segment | Required result |
| --- | --- | --- | --- |
| before prepare | old, cohort hot | absent | no mutation |
| temp create/write/sync/close failure | old, cohort hot | absent/temp only | confirmed abort |
| prepared reopen/verification failure | old, cohort hot | temp/quarantined | no activation |
| segment rename failure | old, cohort hot | absent/temp | confirmed abort |
| segment directory-sync indeterminate | old, cohort hot | absent or unreferenced final | recovery fence; reopen |
| segment published | old, cohort hot | unreferenced final | no capacity release |
| active temp/write/sync/rename confirmed failure | old, cohort hot | unreferenced final | no capacity release |
| active rename/directory-sync indeterminate | old or new | unreferenced or referenced | recovery fence; reopen exactly one complete world |
| activation committed | new, cohort absent/hot indexes updated | referenced | full reopen verification required |
| referenced segment missing/corrupt | new | missing/corrupt | repair-required; never restore old hot state |
| active contains cohort and references it | invalid | referenced | repair-required even for byte-identical duplicate |

Every after-rename injected result is treated as indeterminate even when the test harness can see
the new file. The open process performs no further archive or authority mutation until a newly
constructed store reopens and validates.

### Adapter-call oracle

All archive/migration/backup/offline-verifier tests install a `FakeBackend` call counter and assert
zero prepare/create/start/observe/stop/destroy/reconcile calls. The archive packages must not import
`registeredlifecycle`; a test-only hook observes that no lifecycle entry point was invoked.

## Hot/archive lookup and replay oracle — F4A complete

After one committed archive activation, rerun the existing registration/approval/attempt tests
against both a hot cohort and an archived cohort:

- exact registered-plan lookup returns the original retained bytes and bindings from the segment;
- cap-plus-one or corrupt archived bytes refuse without index fallback;
- exact/equivalent approval envelope replay returns the original `ApprovalID` and archived fixed
  state;
- same payload digest with different exact payload is `REPLAY`;
- same nonce with a different payload is `REPLAY`;
- exact `(RegistrationID, ApprovalID)` attempt replay returns the original `AttemptID`;
- using an archived approval with another registration is `BINDING`;
- archived expiry never becomes usable after a lower clock observation;
- generated registration/approval/attempt/effect/instance collisions with any tombstone refuse;
- an archived `AttemptID` never appears in startup recovery or calls `Recover`; and
- one later, separately approved attempt cannot be created because cohort eligibility required the
  registration to be expired before archive.

The state oracle records before/after hot set digests, archive checkpoint/index digests, counts,
references, time high water, and complete file bytes for every refusal.

## Capacity and bounded-growth oracle — F4C complete

Tests use compact generated records and encoded-size estimators to prove:

- exact active and hot-retained E5 ceilings remain unchanged before archive;
- archiving one eligible cohort releases hot-retained counts only, never active cleanup counts;
- exact 256 cohorts and 64 MiB segment bytes accept; cap plus one refuses before publication;
- exact 64 segment descriptors accept; segment 65 refuses with all cohorts still hot;
- exact 262,144 entries in each index accept; cap plus one refuses without dropping an older
  tombstone;
- active v2 snapshot 64 MiB exact accepts and cap plus one refuses before rename;
- an orphan segment does not consume logical archive capacity until referenced, but physical disk
  exhaustion is a confirmed local failure and never triggers deletion of referenced state; and
- no full-cap case evicts, merges, rewrites, or changes a retained authority state.

Observed result: `PASSED` in the exact fixed-store local-conformance scope. The focused F4C oracle
activates 64 referenced one-cohort segments and refuses segment 65 without rewriting or evicting
the remaining hot cohort or any referenced bytes. The inherited exact inclusive and cap-plus-one
cohort, record, index, segment-byte, and active-byte cases remain passing. Exact answers, fault/
death/race results, and limitations are retained in the
[F4C result](SUPERVISOR_ARCHIVE_F4C_GROWTH_RESULT.md).

The plan makes no throughput, latency, disk-life, or indefinite-service claim. A separate
quantitative budget and production-engine campaign is required before consumer activation.

## Migration and downgrade matrix

The explicit v1-to-v2 migrator checks the owner lock three times: before source open, immediately
before commit, and before returned reopen. It accepts only a complete valid v1 file and constructs
an empty descriptor/segment world plus a complete all-hot retained-global index and the dedicated
generation-one migration-genesis checkpoint.

Required cases:

- v2 opener refuses v0/v1 without rewrite;
- v0/v1 openers refuse v2 without rewrite;
- migration refuses missing file, symlink, non-regular file, wrong owner/mode, missing lock,
  unsupported version, corruption, trailing/duplicate/unknown data, recovery fence, quarantine,
  repair state, and active trust transition;
- pre-rename failure leaves byte-identical valid v1;
- post-rename indeterminate reopens as exactly valid v1 or valid v2;
- v2 authority/lifecycle sets equal the v1 source byte-for-byte and digest-for-digest;
- every v1 registration/approval/attempt/lifecycle identity reconstructs at its typed canonical hot
  ordinal, retained-global indexes/counts equal hot/total counts, and archived/descriptor/segment-
  derived counts remain zero;
- a valid v1 attempt without lifecycle reconstructs on the absent arm, remains in
  `RecoveryAttemptIDs`, contributes one attempt and zero lifecycles, and receives no lifecycle
  state, record location, effect, or instance;
- v2 records the exact count/digest of nonzero effect IDs visible at migration and never claims
  coverage for overwritten pre-v2 fake IDs;
- migration recomputes the genesis checkpoint under its distinct domain and refuses an activation-
  kind reference, wrong hot digest, wrong seed, wrong index scope, or archive location;
- repeated migration refuses v2 and does not reset archive generation/checkpoints; and
- unknown segment/index versions refuse the entire v2 set without ignoring the unknown object.

No automatic migration, fallback, downgrade, or create-on-missing path exists.

## Corruption, repair, and quarantine corpus

For each active/segment/index/checkpoint field, mutate both the field and any obvious local digest
to prove cross-link validation is not merely one checksum. Required families include:

- unsupported/missing versions and generations;
- truncated, trailing, duplicate-name, unknown-field, oversized, wrong-type, and invalid JSON;
- unsafe permissions, symlink, hard-link/replace race, wrong inode, non-regular file, and file-name
  digest mismatch;
- installation, Supervisor, epoch, time-high-water, snapshot-generation, archive-generation,
  ordinal, and prior-checkpoint mismatch;
- record, cohort, set, descriptor, index, combined-index, segment, and checkpoint digest mismatch;
- count/length mismatch and sort-order violation;
- hot/archive overlap or missing cohort;
- split cohort and cross-segment approval/attempt/lifecycle link;
- duplicate registration/approval/attempt/nonce/effect/instance/replay key;
- index points to wrong segment/unit/record digest;
- selected cohort that was not expired/terminal/authoritatively absent; and
- referenced segment replaced by an older otherwise valid segment.

Every case returns `repair-required`, sets attempts disabled in the opened in-memory state where a
state can be established, preserves every on-disk byte, and makes zero adapter calls. Tests compare
file inventories, lengths, digests, and contents before and after refusal.

## Orphan and deletion tests

The only deletion implementation in scope is removal of a known unreferenced artifact created by a
failed local archive/backup operation.

Tests require all of these before removal:

- owner lock held;
- file resides under the configured archive root and was opened without following links;
- active snapshot and every coherent backup manifest do not reference its exact digest/name;
- closed decode, size bound, and digest-addressed name all agree; and
- no unknown link/file-type/ownership/mode condition exists.

Referenced segments, malformed unknown files, digest-mismatched files, tombstones, checkpoint
history, and full cohorts are never deleted. A deletion fault changes no active reference.
Filesystem free-space recovery or secure erase is not claimed.

## Coherent backup and restore-verification plan

### Backup creation

The local conformance backup API accepts only a test-owned destination capability created by the
harness, not an arbitrary caller path. Under the owner lock it copies one fully verified set into a
new empty destination, writes the closed backup manifest last, syncs, and reopens the copy.

Fault injection covers before/after each file copy, file sync, segment directory sync, manifest
rename, root directory sync, and verification read. A partial destination is never a complete
backup and never changes the live store.

### Backup verification

Verification checks the exact manifest inventory, lengths, file digests, active/archive
checkpoint, generations, installation/Supervisor/epoch, time high water, indexes, records, and
cross-links. Extra files are reported with a fixed classification and do not silently join the
backup.

### Restore

This plan implements no live-store replacement. It tests a read-only restore-admission decision:

- exact equality with an injected independent latest-checkpoint fixture may be classified
  `eligible-for-future-restore`;
- older, missing, future, or incomparable checkpoints are `rollback-uncertain`/repair-required;
- no-anchor mode is always rollback-uncertain even when structural verification passes; and
- no decision enables attempts, changes an epoch, deletes the current store, or copies candidate
  bytes over it.

An actual repair/restore ceremony is a later ADR and implementation.

## Offline verifier

Add a read-only package API first; a command-line wrapper is optional and must not precede the API:

```go
type VerificationMode string

const VerificationFull VerificationMode = "full"

type ArchiveVerificationReport struct {
    Version             uint64
    Classification      Classification
    StoreFormatVersion  uint64
    ArchiveFormatVersion uint64
    InstallationID      InstallationID
    SupervisorID        SupervisorID
    EpochSequence       UInt53
    EpochDigest         TrustEpochDigest
    SnapshotGeneration  PositiveUInt53
    ArchiveGeneration   PositiveUInt53
    TimeHighWater       UInt53
    CheckpointDigest    ArchiveCheckpointDigest
    ActiveCounts        ArchiveCounts
    ArchivedCounts      ArchiveCounts
    SegmentCount        uint16
}

func VerifyArchiveSet(context.Context, StoreRoot, VerificationFull)
    (ArchiveVerificationReport, error)
```

`StoreRoot` is a trusted local capability constructed by installation code or the test harness,
not a public string. The report is fixed shape and contains no filenames, paths, labels, plan or
approval bytes, user content, guest text, or arbitrary decoded error. Full verification is
deterministic and read-only; repeated runs return identical reports for identical bytes.

## Exact internal APIs

Add passive types under `internal/execution/archivestate` so format/digest validation imports no
store, lifecycle driver, backend, IPC, content, evidence, or experiment package.

```go
type ArchiveGeneration v0candidate.PositiveUInt53
type ArchiveSegmentDigest [32]byte
type ArchiveCheckpointDigest [32]byte
type ArchiveIndexDigest [32]byte
type ArchiveDescriptorSetDigest [32]byte
type ArchiveCheckpointKind string
type ArchiveIndexScope string
type RecordLocation struct { /* kinded hot ordinal or archive ordinals plus RecordKind */ }

type ArchiveLimits struct {
    MaxCohorts uint16
    MaxBytes   uint64
}

type ArchivePlan struct { /* sealed immutable source/candidate projection */ }
type PreparedArchive struct { /* sealed segment bytes + candidate v2 projection */ }
type VerifiedArchive struct { /* sealed reopened bytes/digests/owner projection */ }
type MigrationGenesisCheckpoint struct { /* generation-one v1-to-v2 checkpoint */ }
type ArchiveCheckpoint struct { /* defensive fixed fields */ }
type ArchiveVerificationReport struct { /* bounded fields above */ }
```

The planned mutation and backup surface remains partitioned behind the internal archive store
boundary; F4A does not implement or call these methods:

```go
type SupervisorArchiveStore interface {
    PlanArchive(context.Context, ArchiveLimits) (archivestate.ArchivePlan, error)
    PrepareArchive(context.Context, archivestate.ArchivePlan) (archivestate.PreparedArchive, error)
    VerifyPrepared(context.Context, archivestate.PreparedArchive) (archivestate.VerifiedArchive, error)
    Activate(context.Context, archivestate.VerifiedArchive) error
    VerifyArchiveSet(context.Context, StoreRoot, VerificationMode)
        (archivestate.ArchiveVerificationReport, error)
    CreateCoherentBackup(context.Context, StoreRoot, BackupRoot) (BackupManifest, error)
    VerifyCoherentBackup(context.Context, BackupRoot) (BackupVerificationReport, error)
}
```

F5 retains the equivalent concrete internal boundary in `registrationstate`: `StoreRoot` and
`BackupRoot` are opaque trusted local capabilities; `CreateCoherentBackup` requires the live
`ArchiveOwner`; `VerifyCoherentBackup` and `VerifyArchiveSet(..., VerificationFull)` are read-only;
`VerifyRestoreAdmission` requires the owner and compares only with an injected
`IndependentLatestCheckpoint`; and `InventoryArchiveArtifacts` returns sealed
`KnownUnreferencedOrphan` values. `DeleteKnownUnreferencedOrphan` accepts only one of those sealed
values and repeats the live and supplied-backup reference scan. None is a public/IPC path API.

F4A extends `internal/execution/registrationstate` with these read-only
`FixedFileStoreV2` methods and closed, location-independent response projections:

```go
ResolveRegistration(context.Context, RegistrationID) (RetainedRegistration, error)
ResolveApproval(context.Context, ApprovalID) (ApprovalRecord, error)
ResolveAttempt(context.Context, AttemptID) (RetainedAttempt, error)
ResolveNonce(context.Context, AttemptNonce) (RetainedNonce, error)
ResolveEffect(context.Context, EffectID) (RetainedEffect, error)
ResolveInstance(context.Context, BackendInstanceDigest) (RetainedInstance, error)
ResolveApprovalReplay(context.Context, ApprovalPayloadDigest, []byte,
    ApprovalKeyAuthorizationIdentity) (ApprovalReference, ApprovalState, error)
ResolveAttemptReplay(context.Context, RegistrationID, ApprovalID)
    (AttemptReference, AttemptState, error)
QueryRetainedIdentityCollisions(context.Context, RetainedIdentityCandidates)
    (RetainedIdentityCollisions, error)
CheckRetainedIdentityAvailability(context.Context, RetainedIdentityCandidates) error
RecoveryAttemptIDs(context.Context) ([]AttemptID, error)
```

Every call freshly full-verifies the active snapshot and referenced immutable segment set, starts
from the `retained-global` index, follows only its record-kind-bound typed location, and rechecks the
full record, digest, cross-links, and applicable replay/uniqueness tombstones. Results exclude the
location, so exact hot and archive records have identical semantics. Missing or corrupt locations
never scan another record collection, fall back to hot history, or try another segment.

The store constructs all archive paths from its trusted root and segment digest. No lookup accepts
plan bytes, approval envelope bytes, backend bindings, caller-selected record bytes, or arbitrary
filesystem paths. F4A's collision result is passive and reserves nothing; F4B must repeat the
complete uniqueness check inside its future atomic mutation. No caller can request hot versus
archive resolution.

## Small implementation slices

### Slice F1: passive archive types and known answers — complete

- Add closed format/index/checkpoint types and domain-separated digest functions.
- Add empty, one-cohort, multi-attempt, exact-bound, cap-plus-one, wrong-domain, aliasing, sort, and
  collision tests.
- Add the pure eligibility selector and prove `RecoveryAttemptIDs` agreement.

Acceptance: passive tests only. No file write, migration, authority mutation, lifecycle entry
point, adapter call, or deletion.

Retained result: merged in PR #75 from exact head
`20c8d7df1d9ed3eb009e8ce9a0afbd41e03807ef` as
`6fc31a049c476acf5085071c48d3d5e36f27240f`. The focused race run reported 86.0% statement
coverage for the passive package. This is a local-mechanic checkpoint, not archive activation.

### Passive F2 format correction — complete

The first implementation review stopped before choosing bytes because merged F1 cannot represent
a required nonzero visible-v1 effect seed with zero archive descriptors, and its activation-shaped
checkpoint cannot produce the required generation-one migration genesis. The retained
[F2 blocker resolution](SUPERVISOR_ARCHIVE_F2_FORMAT_BLOCKER.md) now freezes:

- `retained-global` versus `segment-derived` index domains;
- typed, record-kind-bound `hot` versus `archive` location arms for all global entries;
- exact hot/archive/total/index/descriptor count equations;
- a dedicated migration-genesis constructor/domain/kind distinct from activation; and
- generated digests/counts plus wrong-kind/domain/order/count/cap/collision/copy/mutation cases.

Acceptance: passive types, digests, fixtures, and tests only. No v2 store bytes or lifecycle
behavior changed.

### Passive F2 v1 mapping correction — complete

- Add the closed absent/present lifecycle union to retained attempt entries.
- Derive lifecycle total and location counts only from present arms.
- Bind a present lifecycle to its own typed record location/full-record digest and require the same
  hot/archive world as its attempt.
- Retain generated all-present and missing-lifecycle genesis answers plus the original real-v1
  executable witness.

Acceptance: passive types, digests, fixtures, and tests only. No v2 store bytes or lifecycle
behavior changed.

### Slice F2: explicit fixed-store v2 migration and full verifier — complete

- Add v2 closed open/validation with empty archive known answers.
- Add explicit offline lock-asserted v1-to-v2 migration and downgrade refusal.
- Add full index reconstruction and the corruption corpus for an empty archive.

Acceptance: no cohort leaves hot state and no archive segment exists. V1 behavior and E5 tests stay
unchanged.

Observed result: `PASSED` in this exact scope. The implementation and retained known answers are
recorded in the [F2 stateful migration result](SUPERVISOR_ARCHIVE_F2_MIGRATION_RESULT.md). Owner
checks occur at entry, immediately before commit, and before reopen; confirmed pre-rename failures
preserve byte-identical v1, while post-rename/directory-sync uncertainty returns indeterminate and
fresh reopen observes one complete version. Full reconstruction, downgrade, corruption,
exact/cap+1, concurrency, process-death, and defensive-copy oracles pass with zero adapter calls.

Entry condition was satisfied: the passive contract represents the valid v1 crash state without invention.
See the [F2 v1 mapping resolution](SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md),
`TestAttemptLifecyclePresenceUnionKeepsExactIndependentCounts`, and
`TestFixedStoreV1AttemptWithoutLifecycleHasExactV2Projection`. Stateful work must follow the exact
join, no-adapter, downgrade, and confirmed/indeterminate matrix retained there.

### Slice F3: one immutable segment prepare/verify/activate transaction — complete

- Implement sealed plan/prepared/verified values and one deterministic cohort batch.
- Publish the segment before active reference; add every pre/post-rename fault and reopen oracle.
- Implement zero deletion beyond current-operation temporary cleanup.

Acceptance: one compact generated cohort archives and reopens; all fault cases preserve old or new
complete worlds; zero fake calls.

Observed result: `PASSED` in this exact scope. The implementation, transaction identities, known
answers, owner checks, corruption/substitution/orphan/concurrency/fault/response-loss/process-death
oracles, and limitations are recorded in the
[F3 stateful activation result](SUPERVISOR_ARCHIVE_F3_ACTIVATION_RESULT.md). No retained lookup,
v2 authority mutation, second segment, orphan deletion, backup, production engine, consumer, IPC,
adapter, runtime, backend, or guest was added.

### Slice F4A: read-only retained lookup, replay, and uniqueness routing — complete

- Route registration lookup and approval/attempt replay through hot plus archive indexes.
- Add object-specific registration, approval, attempt, nonce, effect, and instance lookup through
  typed hot/archive locations with identical closed results.
- Add passive candidate collision queries for registration, approval, attempt, nonce, effect,
  instance, approval-replay, and attempt-replay identities.
- Prove archived attempts never enter startup recovery.

Acceptance: every call freshly full-verifies active and referenced immutable bytes; missing,
corrupt, substituted, stale, wrong-kind, or wrong-index-domain locations return the existing
fail-closed classification without fallback or rewrite. Existing Slice B/E5 tests plus archived
equivalents pass. `FakeBackend.CreatesGuest()` remains false and no consumer exists.

Observed result: `PASSED` in this exact read-only scope. The exact APIs, routing semantics, focused
race/corruption/reopen/cap/copy/restoration evidence, and limitations are retained in the
[F4A result](SUPERVISOR_ARCHIVE_F4A_LOOKUP_RESULT.md).

### Slice F4B: atomic v2 authority/lifecycle mutation and effect tombstones — complete

The retained [F4B blocker](SUPERVISOR_ARCHIVE_F4B_MUTATION_BLOCKER.md) records the former
current-effect contradiction and ADR-0031's independent-source correction. The resumed
[F4B result](SUPERVISOR_ARCHIVE_F4B_MUTATION_RESULT.md) implements that correction and atomic v2
mutation in the exact fixed-store local-conformance scope.

- Add ordinary v2 registration, approval, attempt, and lifecycle mutations without routing any
  mutation from caller-selected archive bytes or locations.
- Check every new registration, approval, attempt, nonce, effect, instance, approval-replay, and
  attempt-replay identity against the complete retained-global indexes in the same transaction.
- Commit every newly issued v2 effect tombstone atomically with its `BeginEffect` intent and retain
  it when later lifecycle operations replace the current effect field.

Acceptance: existing Slice B/E5 mutation, replay, response-loss, fault, and recovery tests pass on
v2 with archived collision equivalents. No second segment, new tombstone-only commit, consumer,
adapter call, runtime, backend, or guest is added.

Observed result: `PASSED` in this exact scope. F4C and F5 are separately complete; F6 and product
admission remain open.

### Slice F4C: second-segment activation and bounded growth — complete

- Generalize the sealed F3 activation transaction to a second and later segment within the exact
  64-segment bound.
- Add 64-segment, 262,144-entry-per-index, active-byte, segment-byte, and cap-plus-one no-eviction
  refusal oracles.
- Preserve publish-before-reference ordering and complete old-or-new reopen for every activation.

Acceptance: a second eligible cohort activates into a second referenced segment; every exact cap
accepts and cap plus one refuses with all hot and referenced history unchanged. No cleanup,
deletion, backup, production engine, consumer, adapter, runtime, backend, or guest is added.

Observed result: `PASSED` in this exact scope. The
[F4C result](SUPERVISOR_ARCHIVE_F4C_GROWTH_RESULT.md) retains exact second-segment known answers,
deterministic later-generation/checkpoint ancestry, full retained-global lookup/replay/recovery,
historical/current effect resolution across a mutation between activations, all thirteen activation
fault points, four process-death points, concurrent retry convergence, second-segment mutation and
substitution refusals, exact segment 64 acceptance, and segment 65 no-rewrite refusal. F5 is
separately complete; F6 and product admission remain `BLOCKED`.

### Slice F5: coherent backup, orphan handling, and offline verification — complete

- Implement verified backup creation and read-only restore admission.
- Implement the narrow known-unreferenced orphan deletion rule.
- Add deterministic fixed-shape offline verification report and the complete corruption/fault
  matrix.

Acceptance: no restore activation or referenced-history deletion exists. No rollback-resistance or
secure-deletion claim advances.

Observed result: `PASSED` in the exact owner-held fixed-store local-conformance scope. Backup
manifest v0 copies `state.json` plus every referenced digest-addressed segment and binds exact file
digests/lengths, installation/Supervisor/epoch, current and previous checkpoints, migration
genesis, snapshot/archive ancestry, visible-v1 seed, independent effect-tombstone digest, and
hot/archived/total counts. `manifest.json` is written last and the complete set is independently
reopened. Read-only restore admission is eligible only for exact equality with an injected
independent latest-checkpoint fixture; no-anchor, older, future, and incomparable candidates remain
rollback-uncertain. Offline inventory preserves unknown/corrupt/mixed/cross-installation evidence,
and deletion requires explicit selection of a sealed fully verified known-unreferenced segment plus
a repeated live/supplied-backup reference scan. Exact faults, death/race/cap evidence, known answer,
and limitations are retained in the [F5 result](SUPERVISOR_ARCHIVE_F5_BACKUP_RESULT.md). F6 is not
started.

### Slice F6: production-engine experiment and decision

Accepted ADR-0040 permits a narrower exception before F6 only for the owner-only disposable
internal alpha. That consumer must enforce one owner/installation/Supervisor, local APFS,
concurrency one, full verification before every attempt, no restore, installation retirement on
live-store loss/corruption, and the ADR's exact 128-attempt/8-MiB/16-segment/latency stop thresholds.
This is not production-store selection and makes no continuity, rollback, hostile-same-UID,
secure-deletion, or restore claim. F6 remains mandatory before external alpha or any such claim.
The [passive stop-policy result](INTERNAL_ALPHA_FIXED_STORE_STOP_POLICY.md) implements only the
owner-held read-only admission oracle for those exact thresholds. It remains unwired from attempt
creation and selects no p95 sampling window, product consumer, or engine.

- Instantiate the same logical schema/oracles in one exact SQLite build and configuration, or name
  another candidate with equivalent evidence scope.
- Run `BUSY`/`LOCKED`/`FULL`/`IOERR`/`CORRUPT`, WAL/journal/checkpoint, multi-process, backup,
  migration, APFS, process-death, and real power-interruption campaigns.
- Measure bounded startup, mutation, archive, lookup, backup, and verification resource costs.
- Record a separate ADR selecting or rejecting the engine and exact configuration.

Acceptance: evidence and ADR only. Do not replace the fixed oracle or activate a consumer in the
same slice.

## Required focused tests

Retain all current Slice B/E5 tests and add exact equivalents of:

- `TestArchiveSelectorAcceptsOnlyCompleteClosedRegistrationCohorts`;
- `TestArchiveSelectorNeverSplitsMultiAttemptRegistration`;
- `TestArchiveSelectorAgreesWithRecoveryAttemptIDs`;
- `TestArchivePassiveKnownAnswersAndDefensiveCopies`;
- `TestArchiveFormatCorrectionKnownAnswers`;
- `TestArchiveRecordLocationsAreKindedDiscriminatedAndBoundIntoIndexes`;
- `TestAttemptLifecyclePresenceUnionKeepsExactIndependentCounts`;
- `TestFixedStoreV1AttemptWithoutLifecycleHasExactV2Projection`;
- `TestFixedStoreV1ToV2MigrationAndDowngradeRefusal`;
- `TestFixedStoreV2ReconstructsEveryArchiveIndex`;
- `TestArchiveSegmentPublishesBeforeActiveReference`;
- `TestArchiveActivationFaultAndProcessDeathMatrix`;
- `TestArchiveActivationIndeterminateReopensExactlyOldOrNewWorld`;
- `TestReferencedMissingOrCorruptSegmentNeverFallsBack`;
- `TestArchivedApprovalReplayReturnsOriginalApprovalAndState`;
- `TestArchivedAttemptReplayReturnsOriginalAttempt`;
- `TestArchivedNoncePayloadEffectAndInstanceCollisionsRefuse`;
- `TestFixedStoreV2RetainedLookupsHotAndArchiveAreSemanticallyIdentical`;
- `TestFixedStoreV2RetainedCollisionAndReplayQueries`;
- `TestFixedStoreV2RetainedLookupMissingAndCorruptLocationsFailClosed`;
- `TestFixedStoreV2RetainedLookupReadFaultCapSubstitutionAndRestoration`;
- `TestFixedStoreV2RetainedLookupConcurrencyCancellationAndNoRewrite`;
- `TestFixedStoreV2RecoveryAttemptIDsPreserveHotAbsenceAndExcludeArchive`;
- `TestArchivedAttemptNeverEntersStartupRecovery`;
- `TestArchiveReleasesHotRetainedButNeverActiveCleanupCapacity`;
- `TestArchiveCapsNeverEvictOrDeleteReferencedHistory`;
- `TestArchiveCorruptionEntersRepairRequiredWithoutRewrite`;
- `TestKnownUnreferencedOrphanDeletionRequiresFullReferenceScan`;
- `TestCoherentBackupIncludesAndVerifiesExactArchiveSet`;
- `TestRestoreWithoutIndependentCheckpointIsRollbackUncertain`;
- `TestOfflineArchiveVerificationIsReadOnlyAndDeterministic`;
- `TestDuplicateOwnerRefusesBeforeArchiveWriteOrFakeEffect`; and
- `TestArchivePathMakesZeroAdapterCallsAndFakeStillCreatesNoGuest`.

Every stateful test records before/after active bytes, file inventory and digests, hot and archive
set/index/checkpoint digests, generations, time high water, counts, selected cohort identities,
recovery set, tombstone lookups, authority states, and adapter call counts.

## Production-engine comparison campaign

The SQLite candidate must implement the exact ADR-0031 logical schema rather than replacing it with
an engine-specific notion of terminal rows. At minimum compare:

| Dimension | Fixed checkpoint oracle | SQLite candidate evidence required |
| --- | --- | --- |
| atomic authority/lifecycle transaction | full-file rename | one exact transaction and configured durability |
| archive activation | segment publish then v2 reference | internal closed-cohort move and/or exact export checkpoint |
| uniqueness/replay | reconstructed sorted indexes | constraints plus full cross-table verifier |
| indeterminate result | injected post-rename fence/reopen | process/power loss at commit/WAL/checkpoint/fsync edges |
| owner | injected single process | real nonblocking process ownership and SQLite lock interaction |
| corruption | byte mutation and closed reopen | page/WAL/header/index/table corruption and integrity/refusal behavior |
| backup | lock-held exact file set | pinned backup API procedure plus archive/checkpoint binding |
| rollback | unsupported without anchor | same limitation; database is not an anchor |
| deletion | referenced history forbidden | same until a later retention ADR |
| growth | 64 segments/index caps then refuse | exact page/WAL/index/vacuum caps and deterministic refusal |
| offline verification | full segment/state scan | independent full logical export/scan, not only `integrity_check` |

The production ADR must explain why its engine/configuration reduces a measured risk or satisfies a
required operational property. Developer convenience is not sufficient.

## Activation blockers and limitations

Even after F1-F5, all of these remain blockers:

1. the ADR-0033 Go/Darwin multi-process owner-lock port and installed storage protection;
2. production-engine selection and exact APFS/power-loss evidence;
3. coherent restore/repair ceremony and an independent latest checkpoint or explicit permanently
   disabled rollback-uncertain posture;
4. reviewed retention/evidence/export policy before referenced-history deletion;
5. quantitative service/disk/verification budgets and an operational response to hard capacity;
6. production approval verification/key authorization, authenticated IPC, update/trust
   integration, evidence composition, consumers, and public migration; and
7. every existing runtime, backend, content, and guest admission gate.

The archive oracle bounds active and total storage by refusing. It does not provide indefinite
continuous service. Exact non-reuse holds only across visible retained history. Secure deletion,
anti-rollback, external witnessing, and production readiness remain unimplemented.

## Verification for each retained implementation slice

Use Node.js 22.22.1 or newer, pnpm 10, and Go 1.23 or newer:

```sh
fnm exec --using=22.22.1 -- pnpm install --frozen-lockfile
fnm exec --using=22.22.1 -- pnpm check
fnm exec --using=22.22.1 -- pnpm lint
fnm exec --using=22.22.1 -- pnpm test
fnm exec --using=22.22.1 -- pnpm verify:schemas
fnm exec --using=22.22.1 -- pnpm verify:adrs
fnm exec --using=22.22.1 -- pnpm format:check
go test ./...
go vet ./...
go build ./...
golangci-lint run ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
git diff --check
```

Also run the focused archive packages and all existing approval/attempt, fixed-store v1,
lifecycle-store, and `registeredlifecycle` tests. A documentation-only design change must at least
run the complete repository verification required by `AGENTS.md`; unavailable tools or network-
dependent checks are reported exactly and never converted into a pass.
