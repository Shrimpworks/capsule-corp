# Supervisor archive F2 format blocker

Status: retained local design evidence; no F2 store bytes were implemented.

Date: 2026-08-04

Scope: defensive, repository-local review of the merged F1 passive archive projections against
the v1-to-v2 migration and empty-archive verifier required by Proposed ADR-0031 and
`SUPERVISOR_ARCHIVE_COMPACTION_PLAN.md`. This review opened no runtime, backend, process, service,
guest, identity, credential, user data, or deployment.

## Stop decision

F2 cannot choose one exact closed v2 encoding from the current proposal and merged F1 types
without inventing security-relevant bytes. The conflict is reachable for a valid v1 snapshot with
at least one nonzero currently visible lifecycle `EffectID`, which F2 is explicitly required to
preserve as a limited-coverage seed tombstone.

Implementation stopped before adding a v2 writer, opener, or migration path. Existing v1 bytes
and behavior remain unchanged.

## Retained contradictions

### 1. The only serialized effect index is counted as archived state

The migration must create an empty descriptor set and empty archive-location indexes while
retaining every nonzero effect ID visible in v1. The merged `ActiveStateV2` derives the visible-v1
seed only from `Indexes.Effects`, then requires all index counts to equal `ArchivedCounts`, and
requires `ArchivedCounts` to equal the sum of referenced descriptor counts.

For an empty archive, descriptor counts are all zero. Therefore:

- leaving `Indexes.Effects` empty satisfies descriptor/count validation but loses a required
  visible-v1 seed; and
- inserting a visible-v1 seed preserves the ID but makes `indexes.counts().Effects > 0` while
  descriptor counts remain zero, so `NewActiveStateV2` rejects the candidate.

The retained F1 active-state test exercises only an empty visible-v1 seed. It does not resolve the
nonzero migration case.

### 2. The global reconstruction rule has no hot-record location representation

The plan requires full verification to reconstruct serialized indexes from all decoded hot and
archived records. F1 registration, approval, and attempt index entries each require a positive
three-part `ArchiveLocation`. F2 intentionally has no segment and leaves every authority record
hot, so those global hot projections cannot be serialized with the existing entry types.

Treating those three indexes as archive-only avoids the invalid location, but conflicts with the
plan's all-hot-and-archived reconstruction language and its migration requirement to seed retained
registration/approval/attempt identities. The exact partition between hot-derived uniqueness,
global tombstones, and archive-location indexes is not specified.

### 3. The required generation-one migration has no constructible checkpoint

ADR-0031 requires migration to create v2 at generation one with an empty domain-separated
checkpoint. F1's only `ArchiveCheckpoint` constructor requires both source and result generations
to be positive and requires `result > source`. It therefore cannot construct a checkpoint whose
result snapshot generation is one.

The retained value described by the F1 test as the empty checkpoint uses source generation one and
result generation two. It is an activation-shaped value, not a generation-one migration genesis
checkpoint. Using zero digests in the v2 envelope or silently starting the v2 snapshot at
generation two would choose bytes and semantics not selected by the proposal.

## Existing known answers that remain valid

The contradiction does not invalidate these merged F1 component known answers:

| Projection | SHA-256 known answer |
| --- | --- |
| Empty descriptor set | `a84af9da7e16fadb5aa76f4385558d4bc622ed1ea32ef435899ff02c20e863b3` |
| Empty combined archive index | `2dd78bdddb4e186229d709bdda5b666e4e2d668e5c1216c751be2f4abb46648e` |
| Empty visible-v1 effect seed | `17de5f44f523dab94ca4b215ce7779358146fb094fa6d208e0190cb0ba69e0a1` |
| Activation-shaped empty checkpoint fixture | `69dced13926ca3bfdf7324f7862035480af3b6d26c8b942a36a1ec0db5ee7d54` |

The last value must not be relabeled as the migration genesis answer because its source/result
generations are one/two and its other fixture bindings are specific test values.

## Narrow proposed correction before F2

Revise ADR-0031, the conformance plan, and F1 together in one reviewed format-correction slice:

1. Define separate closed projections for archive body locations and complete visible identity
   tombstones, or define one explicit discriminated hot/archive location type. State exactly which
   registration, approval, attempt, nonce, effect, instance, and replay entries are archive-only,
   hot-derived, or global.
2. Define how hot, archived, and total counts relate to each serialized projection. In particular,
   a visible-v1 effect seed attached to a hot lifecycle must be representable while descriptor and
   archived record counts remain zero.
3. Define a dedicated generation-one migration-genesis checkpoint constructor and digest domain,
   or explicitly authorize a distinct source-generation rule for genesis. Retain a literal known
   answer that binds the real hot set digests, installation/Supervisor/epoch, durable high water,
   empty descriptor/index digests, and visible-v1 seed digest.
4. Add passive tests for a migration-shaped state containing nonempty hot authority sets and a
   nonzero visible-v1 effect seed. Require exact reconstruction, defensive copies, count equality,
   and rejection of seed/index/checkpoint mutation.
5. Freeze the complete closed v2 disk envelope and a byte-exact migration fixture only after those
   projections are selected.

This correction need not authorize segment writing, archive activation, cohort movement, retained
lookup, deletion, a platform owner lock, a production engine, an adapter call, a runtime/backend,
or a guest. Once merged, F2 can implement the explicit three-check owner-lock-asserted migration
and empty-archive opener without guessing.

## Reproduction references

- `docs/adr/0031-checkpoint-closed-supervisor-cohorts.md`, migration steps 3-4 and active-v2 fields
- `docs/SUPERVISOR_ARCHIVE_COMPACTION_PLAN.md`, v2 format, reconstruction, and migration matrix
- `internal/execution/archivestate/projections.go`, `ArchiveCheckpoint` and `ActiveStateV2`
- `internal/execution/archivestate/indexes.go`, archive locations, counts, and visible-v1 seed
- `internal/execution/archivestate/archivestate_test.go`, empty-only active-state fixture
