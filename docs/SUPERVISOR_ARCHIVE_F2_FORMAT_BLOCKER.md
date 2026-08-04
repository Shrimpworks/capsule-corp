# Supervisor archive F2 format blocker resolution

Status: original passive format contradiction resolved; follow-on v1 mapping contradiction also
resolved passively; stateful F2 remains absent.

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
| One-entry all-archive retained-global combined index | `0dcb51ad18810d42148164b30df3fe9b5867e1b33d6890b549b09d9ceb5db7f6` |
| Migration-shaped all-hot retained-global combined index | `e31f70203da7f7909b0df43bb10e5219b4427c723a5236d10ec3de1ed0ac97ff` |
| Generation-one all-present migration genesis checkpoint | `657c86aff68354e535369bdecf9e8c23bbfb1457e8ada57576bc1a52666a5fc9` |
| Missing-lifecycle retained-global combined index | `5f77dd10f8cbe8db00c47eed0ee27b8ec81dd2ccee188be9a2069f5641ab7232` |
| Missing-lifecycle migration genesis checkpoint | `0af76dca782bdf198d5ef80b6b2856fb35ae01539e7d8866198c4f2af643f621` |
| Activation checkpoint fixture | `32d8dd2b25d32ee2e23d5c134f37fbd421291052ad90c97c4829fcdd9f94d56d` |

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

The original three contradictions and the later valid-v1 missing-lifecycle contradiction are
resolved at the passive contract layer. The
[F2 v1 mapping resolution](SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md) selects an explicit
attempt-lifecycle absence/presence union with independently derived lifecycle counts and retains
the exact real-v1 witness. F2 may next add only the explicit owner-lock-asserted fixed-store
v1-to-v2 migration, downgrade refusal, closed v2 opener, and empty-archive full verifier under the
recorded conformance/fault plan.

This correction does not implement `FixedFileStoreV2`, v1-to-v2 file migration, archive segment
writing or activation, cohort movement, retained lookup, deletion, production storage, owner-lock
G2, adapter calls, a runtime/backend/process/service/guest, or any consumer. Overwritten pre-v2
EffectIDs remain unknowable. No v2 bytes were chosen or written by either passive correction.
