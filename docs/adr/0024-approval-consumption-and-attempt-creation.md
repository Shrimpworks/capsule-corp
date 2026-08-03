# ADR-0024: Bound approval consumption and attempt creation before effects

- Status: Proposed
- Date: 2026-08-02
- Refines if accepted: ADR-0011, ADR-0012, ADR-0013, ADR-0019, and ADR-0023

## Context

ADR-0011 requires the Supervisor to verify and consume one approval and create one
`ExecutionAttempt` in one durable transaction before any backend side effect. The passive
`ApprovalGrant` candidate and Gate A2 fixture establish a narrow signed payload shape, while
ADR-0019 deliberately leaves that profile Proposed until its production-wrapper, cross-language,
fuzzing, and exact-byte acceptance conditions pass.

The current unwired implementations stop on opposite sides of the required boundary:

- `registrationstate` durably resolves exact registered plan bytes by `RegistrationID`, using
  bounded storage, monotonic effective time, trust fencing, and fixed internal classifications;
- `registeredlifecycle` accepts only `RegistrationID` and drives a no-guest fake backend, but it
  creates lifecycle state without an approval or `ExecutionAttempt`; and
- the older `execution.SupervisorCore` models approval and attempt state in memory with one generic
  `OpaqueID` type. It is executable scaffold evidence, not the durable store or the contract for
  the registered-plan slice.

Without an exact seam, implementation would have to guess whether an approval reference is a new
authority object, whether the signed nonce is the attempt ID, which time observation controls
expiry, what a retry after process death returns, and whether lifecycle recovery remains keyed by a
registration that may legitimately receive another later approval.

This ADR resolves only the values required for an unwired, fault-injectable approval ledger and
atomic attempt creation. It does not accept ADR-0019, freeze the Gate A2 fixture as a production
profile, activate an endpoint, or authorize a backend or guest.

## Proposed decision

### Boundary and authority

The Supervisor owns two typed local operations:

```text
submitApproval(registrationId, exactSignedApprovalEnvelopeBytes)
  -> ApprovalReference

requestAttempt(registrationId, ApprovalReference)
  -> AttemptReference
```

`submitApproval` is available only to the authenticated Broker role for the exact local IPC
purpose `submit-approval`. `requestAttempt` is available only to the authenticated daemon role for
the exact local IPC purpose `request-attempt`. These method purposes authenticate
the local calls; they do not replace the signed grant purpose or audience.

Neither operation accepts plan bytes, a `PlanRegistration`, backend configuration, image, mount,
path, content, evidence, or lifecycle override. The daemon cannot submit, clear, invalidate,
garbage-collect, or reset approval records. The Broker cannot call the backend or create an
attempt directly.

An `ApprovalReference` is the typed local projection of one Supervisor-issued `ApprovalID`; it is
not a fourth random value, signed object, portable capability, or public textual identifier:

```text
ApprovalReference = { approvalId: ApprovalID }
AttemptReference  = { attemptId: AttemptID }
```

Possession of a reference does not relax any lookup, binding, time, trust, capacity, or state
check. These references are returned only over authenticated typed local IPC. Their future public
text projections are outside this ADR.

### Identifier and nonce domains

`ApprovalID`, `AttemptID`, and `AttemptNonce` are three distinct semantic domains. Each is exactly
a 16-byte byte string, rejects the all-zero value, and must never be accepted in a field belonging
to another domain even when the bytes have the correct length.

- The Supervisor generates `ApprovalID` when it commits the first usable record for one canonical
  approval payload.
- The Supervisor generates `AttemptID` immediately before the consume-and-create transaction and
  commits it only with that transaction.
- The Broker generates `AttemptNonce` before signing the `ApprovalGrant`. It is signed replay
  material and is copied into the approval and attempt records. It is not an `AttemptID`, is never
  used as a store key, and does not select a backend object.

The first fixed store enforces uniqueness of every nonzero ID within its role and uniqueness of
`AttemptNonce` across every retained approval record in the installation, including records from
older epochs. It never derives one domain from another. Cryptographic randomness quality, durable
non-reuse after a future archive/compaction, and rollback-resistant uniqueness remain production
admission requirements; the unwired slice demonstrates collision rejection, not those stronger
properties.

### Exact signed and resolved bindings

The submitted candidate grant is usable only when its verified canonical payload and the current
Supervisor state satisfy all of the following at the same installation:

| Binding | Exact rule |
| --- | --- |
| object | `objectType == "capsule.approval-grant"` and `objectVersion == 0` |
| registration | grant `registrationId` equals the requested and stored registration ID |
| plan | grant `planDigest` equals SHA-256 of the stored exact plan bytes and the stored registration digest |
| installation | grant, registration, plan, and active state have the same nonzero `installationId` |
| epoch | grant `epochDigest` equals the registration and active epoch digest; the registration and plan epoch sequence and digest must also equal the active state |
| Supervisor | grant `supervisorId` equals the registration and active nonzero Supervisor ID |
| purpose | exact UTF-8 string `capsule.plan.approve` |
| audience | exact UTF-8 string `capsule.execution-supervisor` |
| signer | for new admission, protected `kid` selects a locally authorized active-epoch Approval key for this object and purpose; a DID, key ID, or valid signature alone grants no authority |
| nonce | nonzero `AttemptNonce`, unused by every retained approval record in this installation |

The current candidate payload binds an epoch digest, not an epoch sequence. The Supervisor binds
the digest to its one active sequence and copies that resolved sequence into stored approval and
attempt records. Adding an epoch-sequence field to the signed profile is neither required nor
authorized by this ADR.

The verifier must return the exact received envelope bytes, the exact embedded canonical payload
bytes, the decoded candidate view, and the locally resolved Approval-key authorization identity.
The ledger copies them before commit. It hashes payload bytes for semantic replay identity and
envelope bytes only for evidence/integrity; it never uses signature bytes or an envelope digest as
approval identity. Equivalent low-S/high-S signatures over the same payload therefore address the
same approval record.

Production `KeyAuthorization` identity and the production COSE wrapper remain blocked on the
relevant ADR-0019 acceptance conditions. The unwired implementation must inject a fixed verifier
port and fixture-only authorized-key binding; it must not import the Gate A2 experiment or claim
Keychain, user-presence, or production signature validation.

The first unwired wrapper applies these inclusive raw-copy/predecoder budgets before ordinary
allocation:

| Candidate bytes | Maximum |
| --- | ---: |
| complete tagged COSE_Sign1 envelope | 512 |
| embedded canonical ApprovalGrant payload | 256 |
| encoded protected-header map | 128 |

The current closed CDDL has calculated canonical maxima of 431, 242, and 116 bytes respectively,
including a 64-byte `kid` and maximum-width `UInt53` timestamps. The larger power-of-two budgets
bound rejection before the exact candidate decoder enforces those calculated maxima. Cap-plus-one
inputs at each raw boundary reject without ledger change. These are provisional unwired candidate
budgets, not acceptance of ADR-0019 or permission to add fields until the budget is filled.

### Effective time and expiry

All timestamps are nonnegative integer Unix seconds in `UInt53`. The Supervisor uses the same
durable effective-time rule as registration state:

```text
effectiveNow = max(trustedClockObservation, durableTimeHighWater)
```

The high-water mark is persisted before a submission or attempt request can commit. A confirmed
high-water write failure creates no approval and consumes no approval. An indeterminate high-water
write enters recovery-required handling; callers cannot infer that the earlier time remains
authoritative. Clock rollback cannot make an approval or registration usable again, while a
forward jump may expire it early and fails safe.

At approval submission, all of these rules must hold:

```text
issuedAt <= effectiveNow < expiresAt
issuedAt < expiresAt
expiresAt - issuedAt <= 300
expiresAt <= registration.expiresAt
```

At consume-and-create commit, the Supervisor obtains and persists a new trusted-clock observation
and requires again:

```text
issuedAt <= effectiveNow < approval.expiresAt
effectiveNow < registration.expiresAt
```

Equality with either expiry is expired. A grant may be issued at the current second. Submission
does not extend or rewrite either signed timestamp, and attempt creation does not extend the
registration or approval lifetime. `createdAt` in the attempt record is the final effective time
used by the transaction.

### Approval submission, durable states, and replay

`submitted` is an operation phase, not durable authority. The Supervisor bounds, copies,
profile-verifies, and role-resolves the envelope, computes the canonical payload digest, and first
checks for an existing record with that digest. An existing coherent record returns its reference
and current state without re-admission or mutation only after the envelope verifies against the
same signer authorization identity retained by that record. A new payload must pass every current
binding, time, nonce, trust, and capacity check before a ledger record can commit. Untrusted or
rejected bytes are not retained in the authority ledger.

The first ledger has these durable approval states:

| State | Meaning and permitted transition |
| --- | --- |
| `usable` | verified and durably committed; may transition once to `consumed` or `invalidated` |
| `consumed` | terminal; atomically names exactly one committed `AttemptID`; this is the approval being burned even if every later lifecycle operation fails |
| `invalidated` | terminal without an attempt, only from a trust/installation fence or explicit Supervisor-owned repair transition; cannot become usable |

Expiry is a failed usability predicate, not a state transition. An expired record remains retained
and cannot be resurrected. `rejected` is a fixed operation outcome with no approval record. A
bounded diagnostic or later transcript event may record a content-free rejection classification,
but it is not replay authority and is outside this ledger.

Canonical payload digest is a unique ledger key. Resubmitting any valid envelope with the same
canonical payload, including a mathematically equivalent signature, returns the already committed
`ApprovalReference` and its current fixed state without mutation. It never creates a second
approval or makes an expired, consumed, or invalidated approval usable. A retained nonce collision
with a different payload is `REPLAY` and creates no record.

This idempotent lookup is required for a lost response after a successful submission commit. An
indeterminate commit outcome returns `RECOVERY_REQUIRED`, fences approval/attempt mutation, and
requires store reopen and validation. Recovery either finds the one committed record and returns
its reference or establishes that no record committed. It never guesses from the missing response.

### Atomic consume and attempt creation

Before the transaction, the Supervisor may perform only bounded, non-backend preflight:

1. authenticate the caller and validate the nonzero typed references;
2. persist the final trusted-time observation;
3. resolve the exact registration and revalidate its retained plan digest and active bindings; and
4. obtain the point-in-time `RuntimeIntegrityAssessment` required by ADR-0013.

Integrity assessment failure before the transaction leaves a usable approval unconsumed. It may
tighten quarantine or repair state, but cannot create authority. No fake-backend operation,
lifecycle-store write, cleanup lease, content redemption, or guest effect occurs in preflight.

One durable store transaction then performs all of the following against freshly read state:

1. require stable trust, attempts enabled, and no recovery fence;
2. require the registration and approval to exist and repeat every registration/plan/
   installation/epoch/Supervisor/purpose/audience/time binding;
3. require the approval to be `usable` and the approval reference to name that record;
4. enforce attempt capacity and uniqueness of the generated nonzero `AttemptID`;
5. insert one `ExecutionAttempt` in `created` state with no backend handle and no cleanup
   obligation; and
6. transition the approval to `consumed` with that same `AttemptID`.

The attempt copies at least these immutable fields: `AttemptID`, `ApprovalID`, `AttemptNonce`,
`RegistrationID`, registration sequence, plan digest, installation ID, epoch sequence and digest,
Supervisor ID, exact approval purpose and audience, approval payload digest, resolved Approval-key
authorization identity, `createdAt`, and store-format version.

The transaction commits both changes or neither. A confirmed abort leaves the approval usable and
creates no attempt. An indeterminate result returns `RECOVERY_REQUIRED` and fences mutation until
the store is reopened and validated. Recovery may find either the pre-transaction usable record or
the post-transaction consumed record plus its attempt; no state allows consumed-without-attempt or
attempt-without-consumption.

Only after the commit returns durably successful may the internal registered lifecycle perform a
fake side effect. Every post-commit failure belongs to the created attempt and never restores the
approval.

### Request replay, concurrency, and process death

- Repeating the exact `(RegistrationID, ApprovalReference)` after consumption returns the same
  `AttemptReference` and current fixed attempt state without creating or driving another attempt.
- Using the reference with another registration is `BINDING`, even after consumption, and reveals
  no other record fields.
- Concurrent exact requests serialize at the transaction boundary. One creates the attempt; all
  others return the same committed attempt. There are never two successively generated attempts
  for one approval.
- Process death before a confirmed commit leaves no attempt and the approval remains usable.
- Process death after the atomic commit leaves one consumed approval and one `created` attempt.
  Startup recovery enumerates that attempt and may drive it; it does not require new approval.
- Process death after any later fake-backend effect follows the lifecycle cleanup/reconciliation
  rules and still cannot restore the approval.
- Corrupt, cross-linked, partially present, or unverifiable approval/attempt records enter
  `repair-required`; they are never deleted, rewritten into a usable grant, or reported as an
  ordinary rejection.

### Bounded first store and retention

The unwired fixed store uses these inclusive per-installation ceilings:

| Dimension | Maximum |
| --- | ---: |
| usable approvals | 256 |
| retained approval records | 4,096 |
| nonterminal attempts | 256 |
| retained attempt records | 4,096 |

Expired approvals do not consume the usable-approval allowance but do consume retained capacity.
`consumed` and `invalidated` approvals remain retained. `created`, active, terminal, unresolved,
and cleanup-bearing attempts remain retained. The store performs no automatic eviction,
compaction, archival, identifier/nonce tombstone deletion, or capacity-driven state transition.

Submission at approval capacity returns `CAPACITY` without a record. Attempt request at attempt
capacity leaves the approval usable. Capacity checks and their corresponding insert/update occur
in the same transaction. Increasing a ceiling, deleting an explanatory record, or treating expiry
as permission to forget replay state is not a production retention design.

These limits are sufficient only for unwired conformance and the no-guest fake lifecycle. Consumer
activation remains blocked on a reviewed Supervisor-owned archive/compaction and replay-retention
design that preserves approval, attempt, nonce, evidence, cleanup, trust, and rollback invariants.

[Proposed ADR-0031](0031-checkpoint-closed-supervisor-cohorts.md) now supplies that design boundary
without implementing or accepting it. It retains complete closed registration cohorts and exact
replay/non-reuse tombstones in finite immutable segments, forbids referenced-history deletion, and
keeps production engine, owner-lock, power-loss, backup/rollback, and continuous-service evidence
blocked.

### Fixed internal classifications

This slice uses exactly these content-free internal classes; they are not public protocol error
codes or execution results:

| Classification | First-owned examples |
| --- | --- |
| `AUTHENTICATION` | wrong/unauthenticated Broker or daemon method caller |
| `MALFORMED` | bounded envelope/predecoder failure |
| `UNSUPPORTED` | wrong media/object version or unsupported COSE feature |
| `SCHEMA` | closed candidate field/type/width failure |
| `DOMAIN` | correct-width value used in the wrong semantic role |
| `BINDING` | registration, plan, installation, epoch, Supervisor, purpose, audience, or reference mismatch |
| `STALE` | not-yet-effective or expired grant/registration |
| `REPLAY` | nonce reused by another payload or a non-idempotent replay conflict |
| `CAPACITY` | usable/retained approval or open/retained attempt ceiling |
| `TRUST_STATE` | transition-fenced, quarantined, repair-required, or attempts disabled |
| `LOCAL_FAILURE` | confirmed local clock, identifier, verifier-port, read, or transaction failure |
| `RECOVERY_REQUIRED` | commit outcome indeterminate or durable cross-record validation cannot decide safely |

Fake lifecycle failures remain `LIFECYCLE_FAILURE` and `CLEANUP_UNRESOLVED` after attempt creation;
they are not approval-ledger classifications.

Every rejection fixture asserts `authorityStateChanged: false`: no approval is created or made
usable, no approval is consumed or resurrected, and no attempt is created. Persisting a later time
high-water mark or tightening quarantine/repair state is permitted and must be asserted separately;
neither change grants authority. `RECOVERY_REQUIRED` may only fence or tighten state.

### Exact `registeredlifecycle` integration seam

The daemon-facing request remains `requestAttempt(RegistrationID, ApprovalReference)`. The
internal lifecycle does not receive approval bytes or references and does not create the attempt.
After atomic commit it is driven by the Supervisor-issued attempt identity:

```go
type AttemptResolver interface {
    ResolveCreated(context.Context, AttemptID) (CreatedAttempt, error)
}

func (component *Component) Drive(context.Context, AttemptID) (Snapshot, error)
func (component *Component) Recover(context.Context, AttemptID) (Snapshot, error)
```

`CreatedAttempt` supplies the immutable attempt projection and Supervisor-retained exact
registration record. `registeredlifecycle` independently recomputes and revalidates the exact plan
bytes and requires every copied attempt binding to match before its first lifecycle transition.

Accordingly, the current `Execute(context.Context, RegistrationID)` and
`Recover(context.Context, RegistrationID)` interface is replaced rather than overloaded. Its
memory store and snapshots become keyed by `AttemptID` and also retain `RegistrationID`. This is
necessary because one registration may receive a later, separately approved attempt; a
registration-keyed lifecycle store would conflate those attempts. The component remains internal,
unwired, and restricted to `FakeBackend.CreatesGuest() == false`.

Startup recovery enumerates nonterminal attempts from the Supervisor store and calls
`Recover(AttemptID)`. It never re-resolves approval usability or registration expiry to decide
whether an already created attempt needs cleanup.

[Proposed ADR-0025](0025-colocate-durable-attempt-lifecycle-state.md) defines the follow-up storage,
effect-checkpoint, and startup-coordination design. It preserves this authority seam and proposes
colocating lifecycle records in the same versioned Supervisor snapshot; it does not accept this ADR
or implement durable lifecycle behavior.

## Conformance requirements

The focused manifest and fault-injection plan is
[Phase 2B approval/attempt boundary](../PHASE_2B_APPROVAL_ATTEMPT_BOUNDARY.md). No implementation may
wire `registeredlifecycle` to approval state before those cases provide exact state oracles.

## Consequences

- The signed nonce, durable approval identity, and execution-attempt identity cannot be confused.
- Submission and attempt requests are idempotent across lost responses without creating a second
  authority record or attempt.
- One transaction establishes the consumed-approval/created-attempt invariant before fake effects.
- Attempt-keyed lifecycle recovery supports more than one independently approved attempt for one
  immutable registration without accepting replacement plan bytes.
- Fixed no-eviction limits make the first store testable but deliberately unsuitable for consumer
  activation or production retention claims.
- The older in-memory `SupervisorCore` remains scaffold evidence; its generic identifiers, `time.Time`
  semantics, duplicate-approval error, and in-memory store are not the implementation oracle for
  this slice.

## Explicit blockers and non-goals

This ADR does not implement or claim:

- Broker UI, user comprehension, LocalAuthentication, Keychain/Secure Enclave signing, or a
  production Approval-key authorization profile;
- acceptance of ADR-0019, production COSE wrappers, cross-language approval decoders, fuzzing, or
  final signed-object byte ceilings;
- public/daemon SDK or MCP endpoints, public identifier text forms, content access, evidence
  signing, receipts, a real backend, or a guest;
- archive/compaction, backup/restore, coherent-rollback detection, multi-process store locking, or
  non-rollbackable identifier/nonce uniqueness; or
- a claim that the candidate `ApprovalGrant` fixture is a frozen, secure, or production-ready
  profile.

Evidence missing for any blocker must be supplied by its existing gate or a named follow-up
experiment. It must not be inferred from the unwired ledger.
