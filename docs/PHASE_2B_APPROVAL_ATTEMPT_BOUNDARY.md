# Phase 2B approval-consumption and attempt-creation conformance plan

Status: Slice A passive contracts and fixtures implemented on 2026-08-02. No approval ledger,
attempt store, consumer, endpoint, lifecycle integration, backend authority, or guest is
implemented.

Normative proposal: [ADR-0024](adr/0024-approval-consumption-and-attempt-creation.md).

## Objective and defensive scope

Defensively validate Capsule's one-use approval and pre-effect attempt-creation control using only
repository fixtures, a fixed fault-injectable local store, the existing exact registration-state
component, and the no-guest `registeredlifecycle.FakeBackend` in owned local worktrees. Do not
access any other system, identity, credential, or data, and preserve Capsule's existing safeguards.

The remaining implementation slices should answer one question: can a verified candidate approval
become one durable usable ledger record and then one consumed record plus one `created`
`ExecutionAttempt`, atomically and recoverably, before the fake lifecycle performs any effect?

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

## Manifest extension

Extend the closed conformance manifest with object kinds `approval-state` and `attempt-state`.
Each case records:

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

Slice B remains wholly pending. `registrationstate`, `registeredlifecycle`, and the older in-memory
`SupervisorCore` remain behaviorally unchanged and do not import the Slice A package.

### Slice B: fixed approval/attempt store

Implement one bounded, versioned, fault-injectable store extending or transactionally colocated
with registration state. Add submission idempotency, nonce/payload indexes, effective-time
handling, consume-and-create, reopen validation, and capacity tests.

Acceptance: the matrix through process-death and capacity passes with no fake-backend calls.

### Slice C: attempt-keyed fake lifecycle integration

Replace the current registration-keyed lifecycle entry points/store with the ADR-0024
`AttemptResolver`, `Drive(AttemptID)`, and `Recover(AttemptID)` seam. Call it only after a confirmed
consume/create commit.

Acceptance: Plan A cannot realize Plan B, one approval cannot create or drive two attempts, two new
approvals for one registration remain distinct, and every injected post-effect failure is
destroyed, unresolved, or recovery-fenced.

### Slice D: documentation checkpoint only

Reconcile `EXECUTION_SUPERVISOR.md`, the phase checkpoint, and control/evidence matrix with observed
test evidence. Do not claim durable product behavior, production approval, or backend validation.

Consumer wiring, Broker UI/signing, evidence composition, content, and real backend work remain
separate tasks after their gates pass.

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

This plan exits when the unwired store and fake lifecycle pass the exact matrix and the repository
records only the demonstrated claim: atomic one-use state mechanics against a no-guest fake
backend.

It does not exit ADR-0019 or authorize consumers. Remaining blockers include production COSE and
Swift wrappers, Approval-key authorization and Keychain/user-presence evidence, reviewed durable
archive/compaction and rollback behavior, public identifiers/endpoints, evidence profiles, content
custody, and every real-backend/runtime P0 gate.
