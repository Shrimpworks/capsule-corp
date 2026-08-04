# Supervisor archive F3 immutable-segment activation result

Status: `PASSED` for the defensive repository-local first-segment prepare, verify, publish, activate,
and full-reopen scope described here.

Date: 2026-08-04

This result uses only repository-owned fixed-store v2 fixtures, owned temporary files and
directories, injected local faults, and bounded subprocess-death tests. It invokes no lifecycle
adapter, product IPC, runtime, backend, guest, service, identity, credential, user data, deployment,
or third-party target. It adds no dependency and makes no SQLite or production-engine selection.

## Exact predecessor and successor

The only predecessor accepted by F3 is the fully verified F2 migration-genesis world:

- active store format/source versions `2/1`;
- snapshot/archive generations `1/1`;
- one `migration-genesis` checkpoint head;
- no referenced segment; and
- complete all-hot retained-global indexes reconstructed from the exact v1 records.

The only F3 successor is the closed first-activation world:

- active store format/source versions remain `2/1`;
- snapshot/archive generations are exactly `2/2`;
- archive-segment format/ordinal are exactly `0/1`;
- the segment is stored as `segment-<semantic-segment-digest>.json` under the store-derived archive
  sibling, with a closed hex semantic digest and exact encoded length;
- the prior checkpoint is the F2 migration genesis and the current checkpoint kind is
  `activation`;
- every selected registration, approval, attempt, lifecycle, nonce, visible-v1 effect, instance,
  approval-replay, and attempt-replay identity remains in the full segment and retained-global
  indexes on typed archive-location arms; and
- the exact selected records leave the hot collections only in the complete active-snapshot rename
  that installs the descriptor, tombstones, counts, and checkpoint.

F3 accepts one activation and one immutable segment only. A later activation is stale/refused in
this slice; multi-segment growth and v2 mutation remain F4.

## Transaction and owner boundary

`FixedFileStoreV2` now exposes sealed internal plan, prepared, and verified values. They bind the
source snapshot generation, archive generation, checkpoint head, selected cohort digests, limits,
and owner session. Plan selection and preparation perform no filesystem mutation. Verification
closed-decodes the candidate segment, reconstructs its full records and segment-derived indexes,
reconstructs the complete successor retained-global indexes, revalidates the original migration
genesis, and requires deterministic byte equality.

Activation performs this exact ordering:

1. check the single installation owner and revalidate the F2 predecessor;
2. create a mode-`0600` same-directory segment temporary file, write/sync/close/reopen it, bound-read
   and closed-decode it, and require byte-for-byte equality with the sealed candidate;
3. change the verified temporary segment to mode `0400`, sync it, owner-check immediately before
   publication, rename it to its digest-derived final name, and sync the archive directory;
4. write/sync/close a mode-`0600` active-snapshot temporary file, owner-check immediately before
   activation, atomically rename it over the active v2 file, and sync the active directory; and
5. owner-check before a version-specific full reopen of the active snapshot and referenced segment.

Existence of an unreferenced segment never releases hot capacity or grants authority. A valid
unreferenced final segment or interrupted current-operation temporary file is detected and reported
as an orphan; F3 deletes neither. Confirmed failures clean up only the exact current-operation
temporary path. No referenced segment or retained history is deletable.

## Retained known answers

`TestFixedStoreV2ArchivePrepareVerifyActivateKnownAnswer` retains these exact literals for the
compact one-cohort fixture with one consumed approval, one immutable attempt, and one destroyed
lifecycle that retains its visible effect and instance identities:

| Projection | SHA-256 or domain-separated digest |
| --- | --- |
| Activated v2 active file | `dfcaf10b8c3b747a50bf553e73a637574f009874b62615680a5c9813aae28c45` |
| Archive-segment file | `2e1b883ade8b5e7349eea6a842b8f0f8320386e58ba29bfbfb9ff3397303abcd` |
| Semantic archive segment | `e0532c64181191d67d6d8b1f5dd4423277f2dcfa228bc2176e8580880c39fbdc` |
| Activation checkpoint | `1dddc0d54bd5591f9fa586cafbd01543023abb556215220d966794aaab386a92` |
| Retained-global combined index | `0a4bbc13e6deea3f05deca26871e69aa59b7190c8e4b7eba8f0d4efddcd3da3e` |

The successor has zero hot entries and exactly one archived entry in every identity family:
registration, approval, attempt, lifecycle, nonce, effect, instance, approval replay, and attempt
replay. Accessors return defensive copies, and independent preparation under different owned roots
produces identical segment bytes and identities.

## Fault and corruption result

The retained tests cover:

- every segment and active temporary create/write/sync boundary, pre-publication, post-publication,
  segment-directory sync, pre-activation, post-activation, active-directory sync, and response loss
  after successful reopen;
- process death before publication, after publication, before activation, and after activation;
- owner loss before publication, immediately before activation, and before reopen;
- eight concurrent contenders, with exactly one complete successor and seven stale refusals;
- truncated, trailing, duplicate-name, unknown-field, unsupported-version, count, cap-plus-one,
  duplicate/order, cross-link, cohort-digest, wrong-index-domain, effect-tombstone omission, and
  semantic-digest corruption;
- missing referenced segment and substitution with a different internally valid segment; and
- valid-orphan detection without deletion or authority effect.

Every post-rename injected failure is indeterminate. Fresh reopen observes either the byte-identical
generation-one predecessor with the cohort hot and any published segment unreferenced, or the
complete generation-two successor with the cohort archived and referenced. Missing or corrupt
referenced bytes return repair-required and never fall back to hot history. Missing lifecycle state
is preserved as absence and refuses F3 selection; it is never invented or normalized. All merged
F2 known answers and tests remain unchanged.

## Boundary after F3

This is a finite fixed-store conformance oracle, not a production store, durability profile, or
product-admitted archive. F4 retained lookup/replay/uniqueness routing, v2 mutations and additional
segments; F5 backup/orphan-cleanup/offline-report policy; and F6 production-engine comparison remain
open. F3 makes no real power-loss, APFS, multi-process owner, installed protected-root, backup,
restore, rollback-resistance, continuous-service, secure-deletion, runtime/backend/guest security,
or production-readiness claim. ADR-0031 remains Proposed.
