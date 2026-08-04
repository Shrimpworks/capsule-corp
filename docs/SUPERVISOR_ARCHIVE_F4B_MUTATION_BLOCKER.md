# Supervisor archive F4B effect-tombstone mutation blocker

Status: `BLOCKED` for Slice F4B atomic fixed-store v2 authority/lifecycle mutation
*implementation*. The required contract decision below is now recorded in ADR-0031; F4B remains
blocked only on building and verifying it.

Date: 2026-08-04

Parent: Supervisor archive/compaction conformance workstream.

Parent status: `IN_PROGRESS — TRENDING_GOOD`. F1 through F4A remain `PASSED` in their exact
scopes; this result identifies the contract correction required before F4B can safely begin.

## Defensive scope and stop decision

This review used only repository-owned fixed-store v2 types, deterministic local fixtures, closed
in-memory projections, and an owned temporary-file no-rewrite witness. It made no adapter call and
accessed no runtime, backend, guest, IPC, consumer, service, identity, credential, content,
production database, archive deletion, signing operation, or unrelated system or data.

The F4B task carried an explicit stop condition: if F4A's interfaces could not support atomic F4B
without contradicting ADR-0031, retain the exact blocker instead of inventing semantics. That stop
condition is met. No v2 mutation behavior was added.

## Exact contradiction

ADR-0031 requires both of these facts:

1. every newly issued v2 `EffectID` becomes an append-only tombstone in the same transaction as
   its `BeginEffect` intent; and
2. a later lifecycle operation may replace the lifecycle record's single current `EffectID`, but
   must not remove any earlier v2 effect tombstone.

F4A implements a different, narrower invariant that is correct only for the visible-v1/current-
effect scope it passed:

- `reconstructV2Indexes` and `reconstructV2IndexesForWorld` derive at most one effect entry from
  each lifecycle record's current `EffectID`;
- full reopen requires the serialized retained-global effect index to equal that reconstruction;
- archive-segment verification similarly derives segment effects from current lifecycle records;
  and
- `ResolveEffect` requires the selected tombstone's `EffectID`, operation sequence, operation, and
  issuance generation to equal the current lifecycle record before returning the lifecycle as the
  issuing record.

After two v2 lifecycle operations for one attempt, the required state is therefore not
representable under F4A's full-verification and lookup semantics:

```text
required retained effect ledger:  effect-1, effect-2
current lifecycle effect field:              effect-2
F4A full reconstruction:                     effect-2
```

Keeping `effect-1` makes full reopen report a reconstructed-index mismatch. Dropping `effect-1`
violates ADR-0031's non-reuse requirement. Returning the current lifecycle as the issuing record
for `effect-1` violates F4A's exact lookup cross-link. A tombstone-only commit, caller-selected
record/location, alternate scan, fallback, or silent lookup-semantic change is not authorized.

## Retained executable witness

`TestFixedStoreV2F4BHistoricalEffectTombstoneCannotSatisfyF4AReconstruction` constructs the exact
closed two-effect ledger for one attempt: a visible-v1 seed followed by a v2-issued current effect.
The passive `archivestate.ArchiveIndexes` type accepts both sorted entries, proving that the
identifier projection itself is expressible. The F4A fixed-store full verifier then rejects the
same retained ledger because reconstruction can recover only the current effect. The test also
proves refusal preserves the exact active bytes.

This is a representation/verification blocker, not a random-ID collision, capacity issue,
implementation inconvenience, or candidate `NO_GO` decision.

## Required contract decision

The owner of ADR-0031 and the passive archive format must select and retain one exact effect-history
model before F4B resumes. At minimum that decision must define:

- whether every v2 issuance gains a separately digestible immutable effect-tombstone record, or
  whether an index-only tombstone is authoritative;
- how full hot/archive reconstruction proves historical issuance metadata after the lifecycle
  current field changes;
- what `ResolveEffect` returns for a historical tombstone without falsely returning the current
  lifecycle as its issuing record;
- how effect records or index-only tombstones obtain typed hot/archive locations and counts;
- how visible-v1 seeds remain distinguishable from all v2-issued effects;
- how complete-cohort segment construction and verification carry every historical effect; and
- which checkpoint fields and known answers change, with downgrade/version refusal and
  corruption/cap-plus-one tests.

A passive ADR/format/F4A interface-correction slice must land first. It must preserve all passed F2
and F3 byte answers unless the selected versioned correction explicitly replaces them, and it must
retain a compatibility/refusal rule rather than interpreting an omitted effect-history structure
as empty.

## Contract decision recorded

ADR-0031's ["Effect-tombstone source of truth is independent of the lifecycle
record"](adr/0031-checkpoint-closed-supervisor-cohorts.md#effect-tombstone-source-of-truth-is-independent-of-the-lifecycle-record)
now answers each requirement above:

- every v2 issuance gains an entry in an independent, separately digestible effect-tombstone set —
  not an index derived by re-reading lifecycle records, and not a tombstone embedded inside the
  lifecycle record itself;
- full hot/archive reconstruction reads that effect-tombstone set directly, never the current
  snapshot of lifecycle records, so historical issuance metadata survives every later operation on
  the same attempt;
- `ResolveEffect` returns a `superseded-by-current` classification plus the tombstone's own bound
  facts for a historical effect, and only attaches the live lifecycle record as issuing record when
  the tombstone's `EffectID` still matches the attempt's current effect;
- effect-tombstone entries keep the existing `EffectIndexEntry` shape (`archivestate.go`) and its
  existing typed hot/archive `RecordLocation`/count accounting, unchanged;
- the `VisibleV1Seed` discriminant already on each entry continues to distinguish migration-seeded
  v1 effects from every v2-issued effect;
- cohort/segment carriage of the effect index is unchanged, since the correction only changes the
  index's *source*, not its shape or its participation in cohort/segment digests; and
- the only new checkpoint field is one additional `HotSetDigests` member for the effect-tombstone
  set; every other pinned F2/F3 checkpoint field, digest encoding, and known-answer byte is
  unaffected, so no downgrade/version-refusal or corruption/cap-plus-one behavior changes shape —
  the new field is exercised by the same existing tests, extended to cover it.

This decision changes no authority: `EffectID` values remain Supervisor-internal, never
caller-supplied, and only the Supervisor's own committed effect-intent transaction may append to
the effect-tombstone set.

## Resume condition and deferred work

F4B implementation resumes now that this passive contract correction is recorded in ADR-0031. The
implementing slice must add the effect-tombstone hot-state collection, correct
`reconstructV2Indexes`/`reconstructV2IndexesForWorld`/archive-segment verification/`ResolveEffect`
to read from it, add the `HotSetDigests` member, and regenerate the affected F4A/F4B known-answer
fixtures — while still preserving all passed F2 and F3 byte answers untouched. The resumed slice
must still perform same-transaction retained-global checks for registration, approval, attempt,
nonce, effect, instance, approval-replay, and attempt-replay identities; preserve missing lifecycle
as absence; preserve `AttemptID`-only recovery, owner/session fencing, exact capacities,
publish-before-reference ordering, and fail-closed no-rewrite behavior; and add the requested
fault, response-loss, reopen, concurrency, collision, and restoration oracles.

F4C second-segment/bounded-growth work, F5 backup/orphan/offline-report policy, and F6 production-
engine selection remain deferred. This blocker adds no authority, effect, archive, deletion,
adapter, consumer, runtime, backend, guest, or production claim. ADR-0031 remains Proposed.
