# Supervisor archive F2 format blocker resolution

Status: original passive format contradiction resolved; stateful F2 remains absent and is stopped on a follow-on v1 mapping contradiction.

Date opened: 2026-08-04

Date resolved: 2026-08-04

Scope: defensive, repository-local correction of the merged F1 passive archive projections using
only `internal/execution/archivestate`, generated fixtures, and local tests. This work opened no
store, archive file, runtime, backend, adapter, process, service, guest, identity, credential, user
data, or deployment and moved no cohort.

## Original stop decision

PR #78 correctly stopped F2 before any v2 writer, opener, or migration path was added. The merged
F1 model could not represent a valid v1 snapshot with a nonzero currently visible lifecycle
`EffectID`, zero archive descriptors, and a generation-one migration checkpoint without inventing
security-relevant bytes. Existing v1 bytes and selector behavior remained unchanged.

This correction resolves only those three contradictions. It does not implement F2.

## Frozen correction

### 1. Retained-global and segment-derived indexes are different domains

Every `ArchiveIndexesView` now carries one closed digest-bound scope:

- `retained-global` is the complete active-snapshot identity projection reconstructed from every
  decoded hot record and every decoded referenced archive record; and
- `segment-derived` is one segment's archive-only derived projection and is rejected where the
  active snapshot requires `retained-global`.

The scope is included in every per-index digest. Therefore an empty segment-derived index cannot be
substituted for an empty retained-global index even though both contain zero entries.

The active-state count equations are exact for every registration, approval, attempt, lifecycle,
nonce, effect, instance, approval-replay, and attempt-replay family:

```text
hotCounts + archivedCounts = totalCounts
retainedGlobalIndexes.counts = totalCounts
retainedGlobalIndexes.hotLocationCounts = hotCounts
retainedGlobalIndexes.archiveLocationCounts = archivedCounts
sum(referencedDescriptor.counts) = archivedCounts
```

Consequently, migration may serialize a nonzero visible-v1 effect seed in the retained-global
effect index while the descriptor set and every archived count remain zero. The seed is a hot
effect tombstone at migration; it is not counted as an archived cohort or segment entry.

The coverage label remains exactly `visible-v1-seed-plus-all-v2-issued`. Only nonzero effect IDs
still present in v1 lifecycle records are marked `VisibleV1Seed`. Overwritten pre-v2 EffectIDs are
not reconstructed, synthesized, counted, or claimed. The marker, current operation sequence and
operation, issuance generation, attempt binding, and complete effect-index order are digest-bound.

### 2. Every global entry has a discriminated record anchor

`RecordLocation` contains a digest-bound `hot` or `archive` discriminant and the expected
`RecordKind`:

- `hot` requires one positive canonical hot-record ordinal and requires all archive ordinals to be
  zero; and
- `archive` requires positive segment, cohort, and record ordinals and requires the hot ordinal to
  be zero.

No general constructor accepts both arms, a zero arm, an unknown location kind, or an unknown
record kind. Each index validates the exact record anchor below:

| Retained-global entry | Required record anchor |
| --- | --- |
| registration | registration record |
| approval | approval record |
| attempt | attempt record |
| nonce | approval record |
| effect | attempt record |
| instance | lifecycle record |
| approval replay | approval record |
| attempt replay | attempt record |

The entry's identity and full typed bindings still decide record identity; a location never grants
authority by itself. Hot ordinals refer to positions in the same canonical sorted hot collections
whose set digests are retained by the active snapshot. Archive ordinals refer only to a referenced,
fully verified descriptor/segment/cohort/record world. Record kind, location kind, and the selected
ordinal arm are all included in the applicable index digest.

This is the exact projection F2 must reconstruct and compare. It cannot reinterpret a hot ordinal
as a segment ordinal, use an approval location for a registration, or omit hot identities merely
because no archive segment exists.

The field-authority closure is unchanged and explicit: index scope, location kind, record kind,
ordinals, checkpoint kind, counts, and digests are all Supervisor-derived from a fully validated
local snapshot; no caller, daemon, Broker, IPC body, backend, guest, path, or display string supplies
them. The passive constructors validate them and the applicable index/checkpoint digest binds them.
They have no approval-display, user-content, or guest-control role. No existing target in
`schemas/authority/field-authority-manifest.json` changes because these remain internal unwired
archive projections rather than protocol objects; F2 must not expose them as a new public target.

### 3. Migration genesis and archive activation are distinct checkpoints

`MigrationGenesisCheckpoint` is a dedicated generation-one value with digest domain
`capsule.supervisor.archive-migration-genesis.v0`. Its constructor requires and binds:

- store format 2 and migration source format 1;
- result snapshot generation 1 and archive generation 1;
- the exact empty descriptor-set digest;
- every digest and the combined digest of the complete all-hot `retained-global` indexes;
- all four nonzero hot registration/approval/attempt/lifecycle set digests;
- the sorted visible-v1 effect seed, its count, and its digest, exactly matching effect entries
  marked `VisibleV1Seed`;
- exact hot counts, which must equal retained-global index counts;
- installation, Supervisor, epoch sequence/digest, and durable time high water; and
- no archive location, previous checkpoint, or new segment.

`ArchiveCheckpoint` remains activation-shaped under
`capsule.supervisor.archive-checkpoint.v0`. It requires a nonzero, explicitly kinded previous
checkpoint reference, a new segment digest, descriptor/index/hot-set digests, positive generations,
and `resultSnapshotGeneration > sourceSnapshotGeneration`.

Checkpoint references are themselves kinded as `migration-genesis` or `activation`. The active v2
projection permits a genesis head only with no previous checkpoint, no descriptors, zero archived
counts, and archive generation one. An activation head requires a nonempty referenced archive world
and a nonempty genesis-or-activation predecessor. Wrong-kind substitution fails before F2 storage
interpretation.

## Generated exact known answers

The canonical source is
`internal/execution/archivestate/testdata/format-correction-known-answers.json`. It is regenerated,
not hand-edited, with:

```sh
go test ./internal/execution/archivestate \
  -run TestArchiveFormatCorrectionKnownAnswers \
  -update-archive-format-known-answers
```

The retained answers are:

| Projection | SHA-256 known answer |
| --- | --- |
| Empty retained-global registration index | `5e08ca56a013091d27547c6ee6430b94aabb54f6c113bcf616206275a46aed47` |
| Empty retained-global combined index | `78e817b6a07989095010743601a017e43e3b660ea78ad0231e01d900227e207c` |
| Empty segment-derived registration index | `5d965cfe716f2e877f214ac6d384d86ac41a4a24ec5061edd1b2b27bc68382db` |
| Empty segment-derived combined index | `a2ebab040a1fcc665092285392179b7733b7ac327476dafe6e95ee5b6ba95ffb` |
| Empty descriptor set | `a84af9da7e16fadb5aa76f4385558d4bc622ed1ea32ef435899ff02c20e863b3` |
| Empty visible-v1 effect seed | `17de5f44f523dab94ca4b215ce7779358146fb094fa6d208e0190cb0ba69e0a1` |
| One-entry visible-v1 effect seed | `56fba94c52b81ecf559b143e8771e8ff0e36759567fcec071faf8d9d153f0ffa` |
| One-entry all-archive retained-global combined index | `a7a1c329949fb2895e33680a77bc8871dea1e0e2370d676674d2c9f139e88c9d` |
| Migration-shaped all-hot retained-global combined index | `6335a7bed7fab57e286286aee0cf42b04cdaa666ed8a299fe5a7b8fbbae9aaf7` |
| Generation-one migration genesis checkpoint | `3a71b2da5a03570a746cdd535f73dbb159c108469ed9cc727e538a0c53f74704` |
| Activation checkpoint fixture | `b260cffd0b10aaf791f9c1db4a86d901c24ba05a01af7fe2486dec0048d82dea` |

The migration-shaped fixture has one hot and total entry in every count family, including one
visible-v1 effect, and zero archived entries in every family. Its descriptor count is zero. The
generated file also retains the corrected one- and multi-cohort segment answers because index-scope
domain separation changes their segment digests.

## Falsification coverage

Focused passive tests now cover wrong location kind, wrong record kind, mixed hot/archive arms,
wrong index scope, wrong checkpoint kind, order, count and location-count mismatch, cap plus one,
identifier collision, seed/index/checkpoint mutation, and caller/accessor defensive-copy behavior.
They also prove that location-arm mutation changes the applicable retained index digests and that
genesis and activation checkpoint domains and kinds differ.

All existing F1 selector cases remain unchanged, including deterministic cohort ordering,
multi-attempt indivisibility, `RecoveryAttemptIDs` agreement, and every no-resurrection/replay/
`AttemptID`-only lifecycle invariant outside this passive package.

## Remaining boundary

The original three contradictions are resolved. A later implementation review found a separate
contradiction: v1 validly retains a committed attempt before lifecycle establishment, while the
corrected v2 attempt index requires a lifecycle disposition and derives lifecycle count from
attempt count. The exact executable witness and stop decision are retained in the
[F2 v1 mapping blocker](SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md). F2 therefore may not yet add
the explicit owner-lock-asserted fixed-store v1-to-v2 migration, downgrade refusal, closed v2
opener, or empty-archive full verifier.

This correction does not implement `FixedFileStoreV2`, v1-to-v2 file migration, archive segment
writing or activation, cohort movement, retained lookup, deletion, production storage, owner-lock
G2, adapter calls, a runtime/backend/process/service/guest, or any consumer. Overwritten pre-v2
EffectIDs remain unknowable. F2 found another format contradiction and stopped without choosing
new bytes, as this document required.
