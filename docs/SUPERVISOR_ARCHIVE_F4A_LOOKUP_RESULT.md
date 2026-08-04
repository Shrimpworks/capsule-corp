# Supervisor archive F4A retained lookup and uniqueness-routing result

Status: `PASSED` for the defensive repository-local read-only retained lookup, replay, passive
collision, and recovery-routing scope described here.

Date: 2026-08-04

This result uses only repository-owned fixed-store v2 fixtures, immutable F3 segment files, owned
temporary paths, injected local corruption/read faults, and local tests. It performs no v2 write,
authority/lifecycle mutation, `BeginEffect`, tombstone commit, second-segment activation, archive
cleanup, backup, product IPC, consumer, adapter call, runtime, backend, guest, service, external
access, identity, credential, user data, or deployment operation.

## Exact read-only APIs

`registrationstate.FixedFileStoreV2` now provides:

- `ResolveRegistration(context.Context, RegistrationID) (RetainedRegistration, error)`;
- `ResolveApproval(context.Context, ApprovalID) (ApprovalRecord, error)`;
- `ResolveAttempt(context.Context, AttemptID) (RetainedAttempt, error)`;
- `ResolveNonce(context.Context, AttemptNonce) (RetainedNonce, error)`;
- `ResolveEffect(context.Context, EffectID) (RetainedEffect, error)`;
- `ResolveInstance(context.Context, BackendInstanceDigest) (RetainedInstance, error)`;
- `ResolveApprovalReplay(context.Context, ApprovalPayloadDigest, []byte,
  ApprovalKeyAuthorizationIdentity) (ApprovalReference, ApprovalState, error)`;
- `ResolveAttemptReplay(context.Context, RegistrationID, ApprovalID)
  (AttemptReference, AttemptState, error)`;
- `QueryRetainedIdentityCollisions(context.Context, RetainedIdentityCandidates)
  (RetainedIdentityCollisions, error)`;
- `CheckRetainedIdentityAvailability(context.Context, RetainedIdentityCandidates) error`; and
- `RecoveryAttemptIDs(context.Context) ([]AttemptID, error)`.

The returned registration, approval, attempt/lifecycle, nonce, effect, and instance projections
contain no hot/archive location. Their byte-bearing fields are defensive copies. Exact hot and
archive fixtures therefore return semantically identical projections even though their internal
record-location arms differ.

## Routing and classification semantics

Every API invocation:

1. checks cancellation and the trusted active-store file shape;
2. bound-reads and closed-decodes the active snapshot;
3. fully reconstructs and validates the complete retained-global indexes, descriptors,
   checkpoints, counts, hot sets, and every referenced immutable segment;
4. selects the requested object only from the applicable sorted retained-global index;
5. follows only the index's record-kind-bound `hot` or `archive` location;
6. rechecks the selected full-record digest, nominal identity, registration/approval/attempt/
   lifecycle cross-links, and applicable nonce/effect/instance/replay tombstone; and
7. returns a location-independent defensive projection.

No lookup scans hot records or another segment after an index miss or location failure. Missing
nominal/replay identities retain the existing `BINDING` classification. An approval semantic-digest
collision with different exact payload bytes remains `REPLAY`; matching bytes under a different
authorization identity remain `BINDING`. A retained collision makes
`CheckRetainedIdentityAvailability` return `REPLAY`. Missing, corrupt, substituted, stale,
wrong-kind, wrong-index-domain, cap-exceeding, or cross-linked retained storage returns
`ErrStoreRepairRequired`, preserves every byte, and never falls back.

`RecoveryAttemptIDs` iterates only retained-global attempt entries on the hot location arm. It
preserves the existing committed-attempt-without-lifecycle recovery case and never returns an
archived terminal attempt. It accepts and returns only `AttemptID`; it does not invoke recovery or
an adapter.

## Retained verification result

Focused tests prove:

- hot and archived registration, approval, attempt/lifecycle, nonce, effect, instance, approval
  replay, and attempt replay lookups return identical closed semantics;
- archived replay returns the original `ApprovalID`/`AttemptID` and fixed retained state;
- all eight registration/approval/attempt/nonce/effect/instance/approval-replay/attempt-replay
  collision domains remain unavailable, while novel values remain available in this passive
  query;
- missing identities return `BINDING`, payload collision returns `REPLAY`, and authorization
  substitution returns `BINDING`;
- stale ordinal, wrong record kind, wrong active-index domain, and omitted replay tombstone refuse
  through full verification without rewrite;
- missing, cap-plus-one, and internally valid substituted segment bytes refuse without hot or
  alternate-segment fallback;
- deliberate restoration of the exact immutable segment restores the same read result;
- caller mutation of returned plan, binding, payload, and envelope bytes cannot affect later
  reads;
- fresh reopen and 32 concurrent lookup/recovery readers return the same result under the race
  detector and change no active or segment byte; and
- cancellation before a read changes no state and an ordinary retry succeeds.

Read response loss has no commit interpretation in F4A because every operation is read-only and
reserves nothing. Cancellation/retry, missing/read-fault, reopen, and concurrent-read coverage are
the applicable oracles.

## Boundary after F4A

F4A is a finite fixed-store conformance result, not a product store, continuous-service mechanism,
or product-admitted archive. Its collision API is passive and cannot authorize a later write by
itself.

F4B must add atomic v2 authority/lifecycle mutation, repeat retained-global uniqueness checks in
the same transaction, and commit each new v2 effect tombstone with its `BeginEffect` intent. Its
first implementation review is now
[`BLOCKED`](SUPERVISOR_ARCHIVE_F4B_MUTATION_BLOCKER.md): this F4A scope reconstructs and resolves an
effect only when it equals the lifecycle record's single current effect, so it cannot retain an
older v2 effect after a later operation replaces that field without a passive format/lookup
correction. F4C must separately add second-segment activation and the exact 64-segment/index/byte
bounded-growth oracles. F4A implements neither slice. Backup, cleanup/deletion, production-engine selection,
installed owner/protected-root evidence, APFS/real-power-loss evidence, restore/anti-rollback,
consumer/IPC integration, runtime/backend/guest security, and production admission remain open.
ADR-0031 remains Proposed.
