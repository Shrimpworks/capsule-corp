# Supervisor archive F2 v1 mapping blocker

Status: F2 stopped before v2 bytes; passive format correction is insufficient for every valid v1 state.

Date opened: 2026-08-04

Scope: defensive, repository-local review using only the current fixed v1 store,
passive `archivestate` constructors, retained approval/attempt fixtures, and local
tests. No v2 file was written, no store migration was added, no cohort moved, no
archive segment existed, and no lifecycle adapter, runtime, backend, process,
service, guest, identity, credential, user data, or deployment was accessed.

## Stop decision

ADR-0031 Slice F2 cannot implement a total explicit v1-to-v2 migration from the
merged contracts without inventing security-relevant state.

The current v1 contract deliberately permits this complete, valid, durable
startup-recovery world:

```text
registrations = 1
approvals     = 1 consumed
attempts      = 1 created
lifecycles    = 0
```

Approval consumption and attempt creation commit before lifecycle establishment.
`validateV1State` accepts the state, `OpenFixedFileStoreV1` reopens it, and
`RecoveryAttemptIDs` returns the attempt so startup can call
`EnsureLifecycle(AttemptID, ...)`. This is required crash behavior rather than
corruption.

The corrected passive v2 projection cannot represent that world:

- every `AttemptIndexEntry` requires a valid, nonempty `LifecycleState`;
- there is no `missing`, `not-established`, or equivalent lifecycle disposition;
- `ArchiveIndexes.counts()` sets `Lifecycles` to the number of attempt-index
  entries rather than the number of decoded lifecycle records; and
- migration genesis requires the retained-global index counts to equal exact
  hot counts.

For the valid v1 witness, omitting the attempt loses a retained identity and
violates the required complete all-hot index. Including it requires choosing a
lifecycle state absent from v1. Choosing `prepare-pending` would synthesize a
lifecycle record and cleanup state that were not durably committed. Including
one attempt while retaining the exact zero lifecycle count then fails the
checkpoint count equation. None is a permissible migration.

The executable witness is
`TestFixedStoreV1AttemptWithoutLifecycleBlocksV2Projection`. It constructs the
state through the real approval/attempt transaction, explicitly migrates v0 to
v1, reopens and validates v1, proves the attempt index rejects the absent
lifecycle disposition, and proves that an invented disposition still cannot
produce migration genesis with exact `attempts = 1`, `lifecycles = 0` counts.

Focused repetition:

```sh
go test ./internal/execution/registrationstate \
  -run TestFixedStoreV1AttemptWithoutLifecycleBlocksV2Projection -count=20
```

Result on the retained branch: PASS. This is evidence of the contradiction,
not v2 implementation evidence.

## Why refusal is not an implicit fix

The F2 contract says the explicit migrator accepts a complete valid v1 source,
reconstructs every v1 registration/approval/attempt/lifecycle identity, and
preserves the authority/lifecycle sets byte-for-byte and digest-for-digest. The
witness is a complete valid v1 source and is exactly the crash boundary the v1
startup oracle was designed to retain.

Adding an undocumented precondition that every attempt already has a lifecycle
would make migration unavailable for a supported valid v1 state and would not
be the specified total v1-to-v2 migration. Running recovery inside migration is
also forbidden: it would mutate the authority world, require a backend binding
and lifecycle driver, and violate F2's zero-adapter/no-lifecycle-entry boundary.

## Required contract decision

F2 remains blocked until a reviewed passive correction chooses one exact
representation and regenerates its known answers. Plausible design directions
are inputs to that review, not decisions made here:

1. add an explicit attempt lifecycle-presence/disposition union and make exact
   lifecycle counts independent of attempt counts;
2. add a separate retained-global lifecycle index, with attempt entries binding
   an optional lifecycle location under closed cross-link rules; or
3. version the migration contract to document and justify a narrower accepted
   v1 state set plus an independently authorized pre-migration recovery
   ceremony.

Any correction must preserve the v1 crash oracle, exact hot/archive/total count
equations, typed locations, complete reconstruction, generated digests, replay
and non-reuse state, and the rule that migration invokes no adapter. It must
also specify how attempts with and without lifecycle records affect active and
hot-retained capacity.

## Unchanged conclusions

The earlier passive correction remains valid for the contradictions it solved:
retained-global versus segment-derived index domains, typed hot/archive record
locations, visible-v1 effect seeding, and a distinct migration-genesis
checkpoint. Overwritten pre-v2 `EffectID` values remain unreconstructable.

The owner capability can structurally guard the required pre-open,
immediately-pre-commit, and pre-return checks. The atomic predecessor/successor
formats can be distinguished by their closed top-level versions after an
indeterminate rename. Those checks do not resolve the missing-lifecycle mapping.

No `FixedFileStoreV2`, v1-to-v2 writer, v2 opener, migration fault path,
segment preparation/activation, retained lookup, deletion, backup, production
engine, runtime/backend/guest wiring, or consumer was added. ADR-0031 remains
Proposed; archive status remains passive, unwired, and local-mechanic only.
