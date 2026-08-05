# Supervisor archive F4B atomic mutation and effect-tombstone result

Status: `PASSED` for the exact fixed-store v2 local-conformance scope described here.

Date: 2026-08-04

Parent status: `IN_PROGRESS — TRENDING_GOOD`. F4C second-segment growth, F5 backup/restore and
cleanup policy, F6 production-engine selection, product admission, and every adapter/runtime/
backend/guest/consumer integration remain open.

## Defensive scope

This slice uses only repository-owned fixed-store implementations, deterministic in-memory
fixtures, owned temporary files, and the existing fake no-guest lifecycle harness. It accessed no
third-party system, identity, credential, user content, production database, runtime, backend,
guest, IPC endpoint, or consumer. It created no caller-directed archive path and performed no
archive deletion or silent repair.

## Exact retained-format correction

The fixed-store v2 hot state now has one independent Supervisor-owned, sorted, append-only
`effectTombstones` collection. Each entry retains the existing `EffectIndexEntry` facts:
`EffectID`, `AttemptID`, operation sequence, operation, issuance snapshot generation,
`VisibleV1Seed`, and the typed attempt-record location. Its separately domain-separated digest is
`capsule.supervisor.effect-tombstone-set.v0` and is exposed as the fifth
`HotSetDigests.EffectTombstones` member.

The first F4B mutation materializes three new top-level fields: `effectTombstones`,
`effectTombstoneSetDigest`, and the exact `migrationGenesisIndexes` needed to preserve the original
checkpoint independently of later hot mutations. The omission compatibility rule is deliberately
closed to the exact pinned F2 `1/1` migration-genesis world and F3 `2/2` first-activation world.
Every later snapshot requires all three fields. Omission, partial presence, digest mismatch, or a
cross-link mismatch is repair-required and never rewrites evidence.

The F2 and F3 wire bytes and known answers remain unchanged. F4B's one-intent materialized known
answers are:

| Answer | SHA-256/domain digest |
| --- | --- |
| Active snapshot file | `ac7809f713af0cc17ba6b223407869c74eeb4738a4030af7490782142288f87d` |
| Hot effect-tombstone set | `690afbf0d8fff394110858f0fd20f8a493854e91b2bdaf7160320ed79898a1e1` |
| Retained-global combined index | `de753444b8792df867c71f132a7fdb07fd0cdcd7ab5b86bb6ea2ae104c68b4ca` |

That fixture has snapshot generation `2` and exactly one v2-issued tombstone.

## Atomic mutation and replay semantics

Registration, approval, attempt, time-high-water, and every lifecycle mutation now use a complete
v2 read/validate/derive/write transaction. The transaction reloads and fully verifies the active
snapshot and every referenced segment, checks the expected generation/checkpoint, applies only a
Supervisor-owned mutation, reconstructs the complete retained-global indexes, repeats collision
validation for registration, approval, attempt, nonce, effect, instance, approval replay, and
attempt replay, computes all hot/index/count digests, and then uses temporary-file sync, atomic
rename, and directory sync. No caller supplies an archive location or retained record.

`BeginEffect` appends exactly one immutable tombstone in the same transaction as the durable
intent. Same-process, concurrent, response-loss, and process-death retries reconstruct and return
the same permit and `EffectID` without calling the identifier source or appending again. A later
operation may replace the lifecycle record's current effect but cannot remove an older tombstone.
Visible-v1 entries are seeded only from effects actually visible at migration and retain
`VisibleV1Seed=true`; missing v1 history is never invented.

Reconstruction of hot state, complete worlds, and immutable archive segments now reads the
independent tombstone source directly. `ResolveEffect` returns `current` plus the live lifecycle
only when the IDs match. Historical lookup returns `superseded-by-current` and the tombstone-bound
facts with no issuing lifecycle, so the current lifecycle is never misrepresented as the older
issuance record.

## Retained invariants and evidence

The implementation preserves lifecycle `absent | present` joins, hot-only `AttemptID` recovery,
owner/session fencing before and after publication, publish-before-reference segment ordering,
defensive copies, exact active `256`/retained `4,096` authority ceilings, the existing archive
caps, no eviction, and corruption refusal without rewrite, fallback scan, repair, or deletion.
Lifecycle record generations are never moved backwards when the independent v2 snapshot
generation starts below a migrated v1 lifecycle counter.

The focused corpus covers deterministic migration/materialization/reopen, confirmed abort,
pre-state and post-rename indeterminate outcomes, response loss, process death, concurrent exact
replay, visible-v1 and v2 historical/current lookup, mutation after first-segment activation,
segment preservation, owner loss, collection omission/partial/digest/cross-link corruption,
restoration, defensive-copy, substitution, capacity, collision, and full reconstruction oracles.
The inherited F2/F3/F4A and lifecycle suites retain the exact 64 MiB snapshot/segment,
64-segment, 262,144-index-entry, active `256`, retained `4,096`, and cap-plus-one refusal cases.
The focused two-package coverage run (`-coverpkg` across `archivestate` and `registrationstate`)
reports 80.4% combined statement coverage; the F4B race selection and twenty-run deterministic
repetition selection pass separately.

## Boundary after F4B

This is a finite, local fixed-store conformance checkpoint. It is not a production database,
continuous-service claim, real power-loss result, anti-rollback mechanism, coherent backup/restore
design, deletion policy, or product-admitted archive. It adds no second segment and does not begin
F4C, F5, or F6. ADR-0031 remains Proposed.
