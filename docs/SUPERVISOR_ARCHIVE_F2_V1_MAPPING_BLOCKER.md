# Supervisor archive F2 v1 mapping resolution

Status: passive mapping contradiction resolved; the follow-on stateful F2 slice is `PASSED` in its
exact local fixed-store scope.

Date opened: 2026-08-04

Date resolved: 2026-08-04

Scope: defensive, repository-local review using only the current fixed v1 store, passive
`archivestate` constructors, retained approval/attempt fixtures, and local tests. No v2 file was
written, no store migration or opener was added, no cohort moved, and no lifecycle adapter,
runtime, backend, process, service, guest, identity, credential, user data, or deployment was
accessed.

## Original stop decision

The original F2 review correctly stopped before v2 bytes. Current v1 deliberately permits this
complete durable startup-recovery world:

```text
registrations = 1
approvals     = 1 consumed
attempts      = 1 created
lifecycles    = 0
```

Approval consumption and immutable attempt creation commit before lifecycle establishment.
`validateV1State` accepts this state, `OpenFixedFileStoreV1` reopens it, and
`RecoveryAttemptIDs` returns the `AttemptID` so ordinary startup can establish lifecycle later.
The earlier passive v2 attempt index instead required one lifecycle state per attempt and derived
the lifecycle count from the attempt count. Omitting the attempt lost retained authority identity;
supplying `prepare-pending` invented a lifecycle record; and retaining exact counts `attempts = 1,
lifecycles = 0` failed migration genesis.

## Decision

Use one explicit, digest-bound lifecycle-presence union inside each retained-global attempt index
entry:

```text
attempt lifecycle =
  absent
| present(lifecycle state, typed lifecycle-record location, full lifecycle-record digest)
```

The `absent` arm has no lifecycle state, location, or record digest. It is permitted only for a hot
attempt because a cohort cannot archive until every attempt has a destroyed lifecycle record with
cleanup false after authoritative absence. It retains the immutable attempt and its attempt-record
location/digest, remains active startup work, and grants no lifecycle or cleanup state.

The `present` arm requires a valid closed lifecycle state, a `RecordLifecycle` location, and a
nonzero full-record digest. Its lifecycle location must use the same hot/archive arm as the attempt;
for archive locations it must also use the same segment and cohort. The attempt and lifecycle keep
separate typed locations and full-record digests. A location is never authority without the fully
verified referenced record and cross-links.

Counts are now derived independently:

```text
attempt count   = number of retained attempt entries
lifecycle count = number of attempt entries on the present arm
```

Hot/archive lifecycle location counts come from the present lifecycle-record location, not the
attempt location. All existing equations remain exact:

```text
hotCounts + archivedCounts = totalCounts
retainedGlobalIndexes.counts = totalCounts
retainedGlobalIndexes.hotLocationCounts = hotCounts
retainedGlobalIndexes.archiveLocationCounts = archivedCounts
sum(referencedDescriptor.counts) = archivedCounts
```

Full F2 reconstruction must additionally prove a total one-to-zero-or-one join: every retained
attempt has one attempt entry; every decoded lifecycle record has exactly one matching present arm;
an attempt with no decoded lifecycle uses the absent arm; and no effect or instance entry can bind
an absent lifecycle. Existing v1 authority/lifecycle cross-link validation remains authoritative.

## Alternatives compared

| Alternative | Safety and cost | Decision |
| --- | --- | --- |
| Explicit lifecycle-presence union | Directly represents the v1 crash state; binds a present lifecycle to its own typed location and digest; changes one existing index family and makes counts independent | **Selected as the narrowest coherent representation** |
| Separate retained-global lifecycle index | Also safe if it adds a complete ninth index/digest, exact AttemptID join, independent counts, locations, caps, fixtures, and checkpoint bindings | Not selected because it duplicates the same one-to-zero-or-one relationship and widens every index/checkpoint projection without adding authority or lookup semantics |
| Narrow migration plus authorized pre-migration recovery ceremony | Can make lifecycle present only by committing a real v1 lifecycle-establishment transaction first | Rejected as an F2 prerequisite because it adds a state-changing ceremony, backend-binding decision, capacity/fence handling, and new confirmed/indeterminate boundaries; it would make a supported valid v1 crash state conditionally migratable instead of mapping it exactly |

An operator recovery or repair ceremony may be designed later for its own purpose. It is not part of
migration and cannot be used to conceal a missing representation.

## Preserved invariants

- The committed attempt remains immutable and is never omitted from replay/non-reuse indexes.
- Lifecycle absence is represented, not interpreted as `prepare-pending`, destroyed, cleanup false,
  or any other invented state.
- `Drive`, `Recover`, migration recovery enumeration, and lifecycle establishment remain
  `AttemptID`-only.
- Attempts with missing lifecycle remain active and hot-retained. Only a present, fully verified
  `destroyed` record with cleanup false after authoritative absence can release active capacity or
  make a complete cohort archive-eligible.
- V1 authority and lifecycle records remain byte-for-byte and digest-for-digest unchanged by the
  mapping. Failure evidence is neither normalized nor rewritten.
- Migration invokes no adapter and does not establish lifecycle. Overwritten pre-v2 effect IDs
  remain unknowable.

## Passive contract and known answers

`internal/execution/archivestate` now contains only passive constructors and tests for the closed
union. `TestAttemptLifecyclePresenceUnionKeepsExactIndependentCounts` proves absent `1/0` and
present `1/1` attempt/lifecycle counts, digest separation, typed lifecycle anchors, hot/archive arm
agreement, and refusal of an archived absent arm.

`TestFixedStoreV1AttemptWithoutLifecycleHasExactV2Projection` retains the original real v1 witness,
reopens and validates it, maps its attempt to the absent arm, and constructs migration genesis with
exact counts `attempts = 1, lifecycles = 0`. It also proves that lifecycle presence cannot be
asserted using an attempt-record anchor.

The generated fixture
`internal/execution/archivestate/testdata/format-correction-known-answers.json` now retains these
additional exact answers:

| Projection | SHA-256 known answer |
| --- | --- |
| Missing-lifecycle retained-global combined index | `5f77dd10f8cbe8db00c47eed0ee27b8ec81dd2ccee188be9a2069f5641ab7232` |
| Missing-lifecycle migration genesis | `0af76dca782bdf198d5ef80b6b2856fb35ae01539e7d8866198c4f2af643f621` |

Because no v2 writer or stored v2 object ever existed, the corrected attempt-index projection
replaces the earlier passive known answers rather than migrating any bytes.

## Next stateful F2 implementation and conformance plan

The next slice may implement only the explicit fixed-store v1-to-v2 migration and empty-archive
full verifier:

1. Acquire and recheck the installation owner before source open. Bound-read, closed-decode, and
   fully validate v1 without mutation.
2. Reconstruct every all-hot registration, approval, attempt, lifecycle, nonce, visible effect,
   instance, and replay projection from exact v1 records. Join lifecycle by `AttemptID`: absent when
   no lifecycle exists, present only from the decoded lifecycle record and its canonical hot
   ordinal/full-record digest.
3. Independently recompute hot counts from the four decoded record collections and require them to
   equal index counts/location counts. Require absent attempts to be in `RecoveryAttemptIDs`; reject
   lifecycle or effect/instance cross-link omission, duplication, or invention.
4. Construct and fully validate the empty-descriptor active v2 candidate and dedicated migration
   genesis in memory. Prove all v1 authority/lifecycle record bytes and set digests are unchanged.
5. Recheck the owner immediately before the sole temp-file sync/rename/directory-sync commit. No
   adapter, lifecycle establishment, archive segment, or cohort movement is permitted.
6. Recheck ownership, reopen only by top-level version, and run the complete v2 verifier before
   returning success.

Required focused cases include zero attempts, only missing lifecycle, mixed missing/present
lifecycles, all lifecycle states, exact 256 active and 4,096 retained limits, missing lifecycle at
attempt-capacity, wrong/duplicate lifecycle joins, effect/instance on an absent arm, location-arm or
archive-cohort mismatch, count/digest mutation, and defensive-copy/ordering cases. All existing
v0/v1 approval, attempt, lifecycle, recovery, capacity, collision, and corruption tests remain
unchanged.

The fault and downgrade matrix is exact:

- confirmed failure before rename preserves byte-identical valid v1;
- an indeterminate rename/directory sync fences and reopens as exactly valid v1 or exactly valid v2;
- the valid v1 old world may still contain the committed attempt with no lifecycle;
- the valid v2 new world must contain the same attempt on the absent arm with zero lifecycle count;
- v0/v1 openers refuse v2, the v2 opener refuses v0/v1, and neither path rewrites, falls back,
  creates a missing store, or interprets a missing archive collection as empty; and
- ordinary startup recovery may establish lifecycle only after a successful version-specific
  reopen. Migration itself never does so.

## Subsequent stateful result

The logical missing-lifecycle mapping remains closed. The follow-on
[F2 stateful migration result](SUPERVISOR_ARCHIVE_F2_MIGRATION_RESULT.md) implements
`FixedFileStoreV2`, the owner-asserted v1-to-v2 file migration, the closed v2 opener/full verifier,
and its bounded local fault/corruption/capacity/concurrency/process-death corpus. It preserves the
exact valid v1 crash witness as one absent lifecycle arm with independent zero lifecycle count and
runs no recovery adapter.

Archive segment publication/activation, cohort movement, retained lookup, deletion, backup,
production storage, v2 authority mutation, adapter calls, runtime/backend/guest wiring, and every
consumer remain absent. ADR-0031 remains Proposed. This result advances no production, archive
activation, continuous-service, rollback-resistance, runtime, backend, guest, or product-admission
claim.
