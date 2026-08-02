# Phase 2B approval-consumption and attempt-creation checkpoint

Status: Slice A passive contracts and fixtures, Slice B's unwired fixed durable store, and Slice
C's attempt-keyed no-guest fake lifecycle seam were implemented on 2026-08-02. Slice D reconciles
their documentation and evidence without adding behavior. No consumer, endpoint, real-backend
authority, or guest is implemented.

Normative proposal: [ADR-0024](adr/0024-approval-consumption-and-attempt-creation.md).

## Objective and defensive scope

Defensively validate Capsule's one-use approval and pre-effect attempt-creation control using only
repository fixtures, a fixed fault-injectable local store, the existing exact registration-state
component, and the no-guest `registeredlifecycle.FakeBackend` in owned local worktrees. Do not
access any other system, identity, credential, or data, and preserve Capsule's existing safeguards.

The retained Slice A-C implementation answers one local conformance question: a fixture-verified
candidate approval can become one durable usable ledger record and then one consumed record plus
one `created` `ExecutionAttempt`, atomically and recoverably, before the no-guest fake lifecycle
performs an effect. This is not evidence that a production signer, consumer, process boundary,
backend, runtime, or guest can exercise that path.

## Inputs and fixed test authorities

The conformance implementation uses:

- exact registered plan and `PlanRegistration` bytes retained by `registrationstate`;
- the current candidate `ApprovalGrant` payload and envelope known answer, copied into the v0
  conformance tree without changing its field set;
- an injected verifier that returns exact payload/envelope copies, a decoded view, and one fixed
  fixture-only authorized Approval-key binding;
- a scripted trusted clock whose observations update the registration store's durable time
  high-water mark;
- role-specific deterministic identifier sources for `ApprovalID` and `AttemptID`; and
- a fixed file store with confirmed-abort, indeterminate-commit, corrupt-reopen, and process-reopen
  fault points.

The injected verifier is not production cryptography or Keychain evidence. Production wrapper and
signer-authorization work remains blocked on ADR-0019.

## Evidence interpretation

The checkpoint separates four different things that must not be collapsed:

| Evidence kind | Current meaning |
| --- | --- |
| Implemented local mechanic | Repository Go code implements the named unwired state transition or no-guest fake-lifecycle behavior. |
| Retained test evidence | A named focused test or manifest case exercises that exact local mechanic in the owned test environment. |
| Planned production control | Architecture or a Proposed ADR requires the control, but no production mechanism or consumer exists. |
| Unsupported claim | Current evidence cannot support the claim; callers must not infer activation, production cryptography, authenticated IPC, evidence, runtime/backend enforcement, or guest execution. |

The manifest remains exactly 82 rules, 262 cases, and 368 fixtures. Slice C changes no manifest or
fixture count; its additional evidence is the 12 named top-level focused tests in
`internal/execution/registeredlifecycle/component_test.go`.

## Manifest extension

The closed conformance manifest includes the `approval-state` and `attempt-state` object kinds.
Each such case records:

- exact fixture hashes and context paths;
- first owning layer and one ADR-0024 classification;
- before and after approval/attempt counts, set digests, durable time high-water, trust phase, and
  recovery fence;
- expected `ApprovalID`, `AttemptID`, approval state, attempt state, and idempotent reference when
  accepted;
- `authorityStateChanged`, plus separate `timeHighWaterChanged` and `trustStateTightened` booleans;
- whether any fake-backend effect is permitted; and
- Go implementation status, with Swift/TypeScript marked pending only where their future product
  roles actually participate.

The runner must fail if a listed state fixture is skipped, an unlisted approval/attempt fixture is
added, a rejection creates or consumes authority, or a pre-commit case records a fake effect.

## Exact conformance matrix

### 1. Scalar and domain cases

- nonzero 16-byte `ApprovalID`, `AttemptID`, and `AttemptNonce` accepted in their own domains;
- all-zero and 15/17-byte values rejected;
- each correct-width cross-domain substitution rejected as `DOMAIN`;
- duplicate generated Approval or Attempt ID rejected without state change;
- same nonce plus different payload rejected as `REPLAY`; and
- attempt ID is demonstrably independent of the signed nonce.

### 2. Submission authentication, profile, and binding

- only authenticated Broker plus `submit-approval` accepted;
- unauthenticated, daemon, updater, and wrong-purpose callers rejected as `AUTHENTICATION`;
- malformed/noncandidate envelope classifications remain owned by the injected bounded verifier;
- wrong object/version, registration, plan digest, installation, epoch digest, active epoch
  sequence resolution, Supervisor, signed purpose, audience, and signer authorization each reject
  at their exact class;
- registered plan mutation or stored digest mismatch rejects without an approval record;
- approval for registration A submitted with registration B rejects;
- exact accepted envelope and payload bytes are copied before caller/verifier buffer mutation; and
- envelope/signature digest never substitutes for canonical payload replay identity.

The bounded-verifier fixtures also accept the calculated candidate maxima and reject raw envelope
512+1, payload 256+1, and protected-header 128+1 before ledger allocation or state change.

### 3. Issue, expiry, and effective time

- `issuedAt == effectiveNow` accepted;
- `issuedAt == effectiveNow + 1` rejected as not yet effective;
- `expiresAt == effectiveNow` rejected and `expiresAt == effectiveNow + 1` accepted;
- `issuedAt == expiresAt` rejected;
- exact 300-second signed lifetime accepted and 301 seconds rejected;
- `expiresAt == registration.expiresAt` accepted and plus one rejected;
- final attempt-time equality with approval or registration expiry rejects without consumption;
- clock rollback cannot resurrect an expired approval;
- later persisted high-water wins over a lower final observation; and
- confirmed and indeterminate high-water write faults have distinct no-authority/recovery oracles.

### 4. Submission state and replay

- first valid submission creates one `usable` record and returns its reference;
- identical payload and identical envelope returns the same reference without mutation;
- identical payload under an equivalent signature returns the same reference and retains the first
  exact envelope;
- concurrent duplicate submissions create one record and return one reference;
- retry after response loss returns the committed reference;
- duplicate submission of expired, consumed, or invalidated state returns the same state without
  resurrection; and
- rejected bytes create no replay tombstone or authority ledger record.

### 5. Atomic consume and create

- only authenticated daemon plus `request-attempt` accepted; unauthenticated, Broker, and updater
  callers reject as `AUTHENTICATION` before lookup or state change;
- stable, usable exact bindings commit one consumed approval and one `created` attempt;
- the attempt copies every immutable binding listed by ADR-0024 exactly;
- attempt capacity, identifier failure/collision, stale time, binding mismatch, trust fence, and
  integrity preflight failure create no attempt and leave the approval usable;
- confirmed transaction abort changes neither approval nor attempt set;
- indeterminate commit returns `RECOVERY_REQUIRED` and no fake effect is permitted;
- state validation rejects consumed-without-attempt, attempt-without-consumption, cross-linked IDs,
  mismatched nonce/reference/bindings, and duplicate attempt ownership; and
- no lifecycle `prepare` call is observed before the atomic store commit checkpoint.

### 6. Replay and concurrency after consumption

- exact replay returns the same `AttemptReference` without a second attempt;
- cross-registration replay is `BINDING` and does not reveal or mutate the attempt;
- two simultaneous exact requests return one shared attempt identity;
- sequential replay after fake lifecycle completion does not redrive the lifecycle; and
- an invalidated approval never creates an attempt.

### 7. Process-death and reopen matrix

Inject interruption and reopen at these exact boundaries:

| Boundary | Required recovered outcome |
| --- | --- |
| before approval commit | no approval; retry may submit |
| after approval rename, before directory durability result | recovery-fenced until reopen decides; never create a second record |
| after durable usable approval, before response | same reference returned on retry |
| before consume/create commit | approval usable; no attempt |
| after consume/create rename, before directory durability result | recovery-fenced until reopen decides pre- or post-state |
| after durable consume/create, before lifecycle call | approval consumed; one `created` attempt enumerated for recovery |
| after each fake effect | approval remains consumed; existing lifecycle destroy/reconcile oracle applies |

Reopen validation must recompute registration and approval payload digests, validate every copied
binding and cross-reference, enforce store versions/capacities, and either load one coherent state
or enter repair-required. Missing later state is never evidence that an earlier commit did not
happen.

### 8. Capacity and retention

- 256 usable approvals accepted; the 257th rejected without eviction;
- expired approvals stop counting as usable but remain in retained count and nonce replay index;
- 4,096 retained approvals accepted; cap plus one rejected;
- 256 nonterminal attempts accepted across independently approved registrations/retries; cap plus
  one leaves its approval usable;
- 4,096 retained attempts accepted; cap plus one rejected;
- terminal, consumed, invalidated, unresolved, cleanup-bearing, and sole explanatory records are
  never evicted; and
- restart does not reset counts, IDs, nonce uniqueness, or set digests.

Large exact-capacity cases may use compact generated contexts plus repository-verified hashes; they
must not add thousands of redundant opaque fixtures.

### 9. `registeredlifecycle` seam

- `Drive` and `Recover` accept only nonzero `AttemptID`;
- no method accepts approval bytes/reference, replacement plan bytes, backend flags, or paths;
- `AttemptResolver.ResolveCreated` returns copied exact plan bytes and immutable attempt bindings;
- mutation after resolution cannot alter stored state;
- wrong attempt/registration/plan/approval binding rejects before fake `prepare`;
- two separately approved attempts for one registration have distinct lifecycle records;
- recovery by Attempt ID survives registration and approval expiry; and
- startup enumeration drives a durable `created` attempt exactly once and reconciles later states.

## Implementation slices

### Slice A: passive contracts and fixtures

Add role-specific Go types, fixed classifications, state/context fixture schema, and the exact
candidate approval known answer. Do not implement COSE or reuse experiment packages.

Acceptance: all byte/state oracles are reviewable, fixture integrity passes, and no product package
or consumer is wired.

Observed checkpoint: `internal/execution/approvalattempt` now provides the three distinct nonzero
16-byte identifier/nonce domains, typed local references, the exact 12-class internal vocabulary,
passive approval/attempt states, and an injected fixture-only verifier. The verifier recognizes only
retained vectors, returns defensive copies, and applies the inclusive 512/256/128 raw ceilings plus
the calculated 431/242/116 closed-candidate maxima. It is not a production COSE or Approval-key
authorization implementation.

The manifest grew from 67 rules, 206 cases, and 278 fixtures to 78 rules, 250 cases, and 350
fixtures. The 44 new Go cases contain 10 accept and 34 reject oracles for identifiers, references,
the vocabulary, the exact 375-byte candidate known answer, complete signed/resolved bindings,
copied byte ownership, calculated maxima, and cap-plus-one refusal. Every Slice A context retains
the same empty approval/attempt state, sets `authorityStateChanged`, `timeHighWaterChanged`, and
`trustStateTightened` false, and forbids a fake-backend effect.

Observed fixture identities include:

- envelope: SHA-256 `fb0a9e7c983f6f3986260dce857edf6b18cba99ee386f9532300dbdc31a5a3bd`, 375 bytes;
- payload: SHA-256 `8ed203acb49409cf2c787bcb04e5e40aaed7139e8bc5b599bd53a49fb3c0e6ea`, 234 bytes; and
- protected header: SHA-256 `b79d430399eb9d3f3690735f03a021a80a24f1ea76821303cf90fd010033ecbf`, 68 bytes.

Before Slice B, `registrationstate`, `registeredlifecycle`, and the older in-memory
`SupervisorCore` remained behaviorally unchanged and did not import the Slice A package.

### Slice B: fixed approval/attempt store

Implement one bounded, versioned, fault-injectable store extending or transactionally colocated
with registration state. Add submission idempotency, nonce/payload indexes, effective-time
handling, consume-and-create, reopen validation, and capacity tests.

Acceptance: the matrix through process-death and capacity passes with no fake-backend calls.

Observed checkpoint: `registrationstate.FixedFileStore` now transactionally colocates registrations,
the shared durable effective-time high-water, approval records, and immutable created-attempt
records in one versioned snapshot. `ApprovalAttemptComponent` exposes only authenticated
`SubmitApproval`, authenticated `RequestAttempt`, defensive record inspection, and created-attempt
enumeration. It injects the Slice A candidate verifier, role-specific identifier sources, a
point-in-time integrity assessor, and process-boundary checkpoints. It does not call
`registeredlifecycle`, `FakeBackend`, or another effect.

The store keys semantic replay by SHA-256 of exact canonical payload bytes, retains the first exact
envelope, rejects nonce reuse by another payload, and copies every byte and binding before commit.
One transaction changes `usable` to `consumed`, binds the generated `AttemptID`, and inserts the one
immutable `created` attempt. Exact and concurrent retries return the already committed references.
Confirmed aborts retain the pre-state. An injected post-rename/pre-directory-durability result
returns `RECOVERY_REQUIRED` and fences mutation until reopen validates one coherent snapshot.
Reopen rejects digest, version, capacity, missing-half, duplicate-owner, and cross-link corruption.

The compact manifest now retains 82 rules, 262 cases, and 368 fixtures. Twelve new Go-verified
Slice B state-transition cases cover first submission, payload/nonce replay, atomic consume/create,
exact request replay, confirmed abort, both indeterminate recovery outcomes, reopen, and exact/cap-
plus-one active ceilings. Focused Go tests additionally exercise defensive copies, concurrent lost-
response retries, time equality/lifetime/rollback, caller and signed bindings, trust/integrity
fences, identifier collisions, the exact 256 usable-approval and nonterminal-attempt ceilings, and
4,096 non-evicting retained approvals across reopen without creating thousands of opaque fixtures.
The 4,096 retained-attempt ceiling is enforced by the store validator; reaching it through normal
operations remains impossible before Slice C supplies terminal lifecycle disposition, because
every Slice B-created attempt is intentionally nonterminal.

This is unwired conformance evidence only. It does not accept ADR-0019 or ADR-0024, validate
production cryptography or signer authorization, provide multi-process locking/archive/compaction,
or authorize a lifecycle, backend, runtime, content handle, or guest.

### Slice C: attempt-keyed fake lifecycle integration

Replace the current registration-keyed lifecycle entry points/store with the ADR-0024
`AttemptResolver`, `Drive(AttemptID)`, and `Recover(AttemptID)` seam. Call it only after a confirmed
consume/create commit.

Acceptance: Plan A cannot realize Plan B, one approval cannot create or drive two attempts, two new
approvals for one registration remain distinct, and every injected post-effect failure is
destroyed, unresolved, or recovery-fenced.

Observed checkpoint: `registrationstate.ApprovalAttemptComponent` now implements the narrow
`ResolveCreated(AttemptID)` seam. It validates the complete durable store, consumed-approval/
created-attempt cross-link, exact retained registration record, and role-bound plan bytes, and
returns defensive copies without signed approval bytes. `registeredlifecycle.Drive` and `Recover`
accept only a nonzero `AttemptID`; lifecycle records, fake instances, fault keys, checkpoint hooks,
idempotency, and recovery are all attempt-keyed while retaining the bound registration and complete
immutable plan/approval/trust/policy/runtime bindings.

Startup recovery enumerates the committed Slice B `created` attempts once per recovery pass and
drives a missing lifecycle record or reconciles an existing one. Exact replay returns the retained
lifecycle record without repeating an effect, and separately approved attempts for one registration
remain independent. The migrated local corpus covers every fake before/after-effect fault and every
post-effect interruption checkpoint; outcomes remain destroyed, explicitly unresolved, or
recovery-fenced. `FakeBackend.CreatesGuest()` remains hard-coded false. This is unwired fake
lifecycle evidence only: the lifecycle store remains the existing bounded single-process memory
store, and no consumer, production approval wrapper, content path, real backend, or guest is added.

### Slice D: documentation checkpoint only

Reconciles `EXECUTION_SUPERVISOR.md`, the phase checkpoint, workstream ledger, and control/evidence
matrix with observed test evidence. It does not claim durable product behavior, production
approval, or backend validation.

Observed checkpoint: Slice D changes documentation only. The evidence map is:

| Slice | Implemented local mechanic | Retained evidence |
| --- | --- | --- |
| A | Distinct nonzero `ApprovalID`, `AttemptID`, and `AttemptNonce` domains; typed references; fixed states/classifications; defensive byte ownership; retained-vector-only verifier with raw and calculated maxima. | `TestSliceAConformanceManifest`, `TestDomainIdentifierCopiesInputAndReferencesRejectZero`, and `TestClassificationProjectionIsDefensive`; 44 manifest-backed Go cases, including the 375-byte known answer and 431/242/116 calculated maxima. |
| B | One versioned fixed snapshot colocates registration, effective-time high-water, approval, and created-attempt state; canonical-payload replay identity, retained nonce uniqueness, exact request replay, and atomic consume/create precede effects. | `TestApprovalAttemptDurableIdempotentAtomicHappyPath`, `TestApprovalAttemptConcurrentExactRequestsConverge`, `TestApprovalSubmissionTimeBoundaries`, `TestApprovalAttemptFaultAndProcessDeathMatrix`, `TestApprovalAttemptReopenRejectsCorruptionAndCrossLinks`, `TestApprovalAttemptHighWaterAndAttemptCapacity`, `TestApprovalCapacityAndNonEvictingRetainedState`, and `TestSliceBDurableStoreManifestOracles`. |
| C | `AttemptResolver` and `CreatedAttemptEnumerator` expose committed created attempts; `Drive`, `Recover`, lifecycle state, fake instances, and faults are keyed only by `AttemptID`; bindings and exact plan bytes are revalidated before fake prepare. | All 12 top-level `registeredlifecycle` tests, including `TestDriveRejectsMissingMutatedAndCrossLinkedAttemptsBeforePrepare`, `TestConcurrentAndSequentialDriveNeverRedriveEffects`, `TestTwoApprovalsForOneRegistrationHaveIndependentAttemptLifecycles`, `TestFakeFaultMatrixEndsDestroyedOrUnresolved`, `TestPostSideEffectInterruptionRecoversByAttemptAcrossComponentRestart`, and `TestStartupEnumerationDrivesCreatedAttemptOnceAndIgnoresExpiry`. `TestDriveResolvesCommittedSliceBAttemptAndExposesNoSuccessResult` also asserts `FakeBackend.CreatesGuest() == false`. |

The Slice B fixed authority store is durable only within this unwired file-snapshot test design.
The Slice C `MemoryStore` is bounded, single-process, and non-durable; reusing it across component
instances exercises in-process restart logic, not process-death or power-loss persistence.

Consumer wiring, Broker UI/signing, evidence composition, content, and real backend work remain
separate tasks after their gates pass.

## Next backend-independent boundary

[Proposed ADR-0025](adr/0025-colocate-durable-attempt-lifecycle-state.md) and the separate
[durable lifecycle implementation plan](PHASE_2B_DURABLE_ATTEMPT_LIFECYCLE_PLAN.md) now select the
narrow design for the next implementation slice: lifecycle records, effect intents, and terminal/
unresolved dispositions are colocated in the existing versioned Supervisor snapshot and
transaction domain. Startup and recovery remain `AttemptID`-only, resolve only committed created
attempts, and do not recheck approval usability or registration expiry.

That proposal is design only. `MemoryStore` remains the implemented Slice C store until the
fake-only slices and fault matrix land. Authenticated IPC, archive/compaction and replay retention,
production signing/verification, real multi-process locking, rollback/backup, production backend
reconciliation, consumers, evidence composition, and public cutover remain separate blockers.

## Verification for each retained slice

```sh
fnm exec --using=22.22.1 -- pnpm install --frozen-lockfile
fnm exec --using=22.22.1 -- pnpm check
fnm exec --using=22.22.1 -- pnpm lint
fnm exec --using=22.22.1 -- pnpm test
fnm exec --using=22.22.1 -- pnpm verify:schemas
go test ./...
go vet ./...
go build ./...
git diff --check
```

The repository-required commands remain authoritative if this list becomes stale.

## Exit and blockers

This checkpoint exits with the repository recording only the demonstrated claim: atomic one-use
authority-state mechanics and `AttemptID`-keyed recovery against a no-guest fake backend.

It does not accept ADR-0019 or ADR-0024 and does not authorize consumers. The fixed authority store
has no eviction, production archive/compaction, multi-process locking, backup/restore design, or
rollback-resistant identifier/nonce uniqueness. The lifecycle `MemoryStore` is single-process and
non-durable. Production COSE and Swift wrappers, Approval-key authorization,
Keychain/LocalAuthentication user-presence evidence, authenticated IPC, public identifiers and
endpoints, consumers, content custody, evidence composition/signing, runtime/backend enforcement,
and every real-backend/runtime P0 gate remain unresolved. No guest is created.
