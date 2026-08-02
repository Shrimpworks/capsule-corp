# Execution Supervisor

Status: intended v0 authority and interface; unprivileged per-user macOS topology is feasible. The
native libkrun/HVF adapter is the lead Apple candidate under evaluation, while stock-Bun authority,
immutable runtime-root custody, `NullFs` disposition, typed port transport/completion, complete
installed-bundle admission, and the OCI/gVisor comparison remain gated. Bounded filesystem-image
parsing is a later gate before file artifacts.

Implementation checkpoint: the older `internal/execution.SupervisorCore` remains in-memory scaffold
evidence and is not the oracle for ADR-0024. The current unwired path instead comprises:

- `internal/execution/approvalattempt`: passive typed domains, fixed classifications/states, and a
  retained-vector-only bounded verifier;
- `internal/execution/registrationstate`: a fixed file snapshot colocating exact registrations,
  effective-time high-water, approvals, and immutable created attempts, including atomic
  consume/create and reopen validation; and
- `internal/execution/registeredlifecycle`: an `AttemptID`-keyed, no-guest fake lifecycle that
  revalidates the committed attempt and exact registered plan before fake prepare.

These are implemented local mechanics with retained Go tests, not a deployed Supervisor. The
authority store has no production archive/compaction or multi-process locking, the lifecycle store
is bounded single-process memory, and no consumer, authenticated IPC, macOS code identity,
Keychain/user presence, protected production storage, production cryptography, content, evidence,
runtime/backend authority, or guest is present. ADR-0019 and ADR-0024 remain Proposed.

## Purpose

The Execution Supervisor is the sole component allowed to create a hostile guest. Its narrow job is
to turn a previously registered immutable plan plus a valid one-use approval into at most one
bounded attempt, and to produce an attributable transcript of its enforcement observations.

The Supervisor is not the public Capsule API, general policy engine, content UI, updater, parser
host, or runtime package manager.

## Owned authority

The Supervisor owns:

- independent raw, schema, and hard-safety plan validation;
- exact registered plan bytes and registration sequences;
- approval verification and durable atomic consumption;
- one-attempt lifecycle and terminal classification;
- runtime-integrity/quarantine state;
- exact backend capability matching;
- job-scoped content-handle redemption and staged-digest verification;
- backend creation, observation, termination, destruction, and reconciliation;
- filesystem-safe collection and artifact manifests;
- hash-linked enforcement events and the Supervisor evidence key.

The Supervisor must not own:

- agent-facing protocol parsing or user-facing job prose;
- general policy authoring;
- file selection or original-path metadata;
- user-only artifact delivery;
- rich parsers/previews;
- arbitrary network/DID/TUF resolution;
- runtime/profile building or dependency installation;
- a general update UI;
- responsibilities added only for convenience without an ADR.

## Interface sketch

The exact wire contracts remain pending, but the conceptual interfaces are:

```text
registerPlan(exactCanonicalPlanBytes) -> PlanRegistration
getRegisteredPlan(registrationId) -> exactCanonicalPlanBytes
submitApproval(registrationId, exactSignedApprovalEnvelopeBytes) -> approvalReference
requestAttempt(registrationId, approvalReference) -> attemptReference
cancelAttempt(attemptId) -> acceptedState
getFixedStatus(attemptId) -> boundedStatus
getTranscript(attemptId) -> EnforcementTranscript
reconcile() -> boundedRecoveryReport
```

No execute/attempt operation accepts replacement plan bytes, backend flags, guest paths, images, or
policy overrides.

The Broker has a direct authenticated read/approval/content-handle surface. The daemon receives
only planning/registration/attempt/status/cancellation operations. A privileged launcher, if any,
receives only a sealed typed launch descriptor created by the Supervisor.

These role-specific surfaces are intended and proposed interfaces. Current tests inject an
`AuthenticatedCallContext`; they do not implement or validate authenticated IPC or choose its
process topology.

## Independent validation

Registration repeats strict raw and schema validation even if the daemon claims it already passed.
The Supervisor then applies versioned hard-safety invariants. For v0 it rejects any representation
of:

- network or host IPC authority;
- subprocess names or general process authority;
- environment inheritance or secret lookup;
- native addons, FFI, macros, inspector, package installation, or dependency resolution;
- arbitrary runtime images, registries, kernels, backend flags, mounts, or management sockets;
- host or guest paths supplied by the agent;
- unsupported output/content parsers or future capability unions;
- resource values the exact backend cannot enforce.

Unsupported powers fail as unsupported protocol. They are not accepted merely because policy says
`false`. The runtime profile must also refuse those powers in execution; the current stock Bun
profile has not evidenced non-bypassable subprocess/FFI denial.

## Durable lifecycle

```text
registered
  → awaiting-approval
  → preflighting
  → approval-consumed / attempt-created
  → preparing
  → created
  → staged
  → running
  → terminating
  → collecting
  → destroying
  → reconciled
  → terminal
```

Pre-consumption rejection may leave a registration approvable if its state is still determinate.
Once the approval is consumed, every result belongs to that attempt and another attempt requires a
new grant.

Every path after backend creation reaches destroy/reconcile. Terminal classifications distinguish
guest failure, policy/integrity denial, cancellation, timeout, resource exhaustion, backend
failure, egress rejection, teardown failure, unresolved backend, and success.

The current local evidence stops at a narrower boundary. `ApprovalAttemptComponent` durably
creates one immutable `created` attempt and exposes `ResolveCreated(AttemptID)` plus startup
enumeration. `registeredlifecycle.Drive` and `Recover` accept only that `AttemptID`, revalidate the
attempt/consumed-approval/registration/plan bindings, and key lifecycle records, fake instances,
faults, and recovery by the attempt. Separately approved attempts for one registration therefore
remain distinct. The lifecycle `MemoryStore` is non-durable, and the fake lifecycle has no job
success result and creates no guest.

## Side-effect ordering

1. Persist registration before issuing the registration ID. The unwired fixed store exercises
   this locally.
2. Persist grant and validate its bindings. The fixed fixture verifier and authority store
   exercise only the candidate local profile.
3. Atomically consume grant and create attempt before backend side effects. Slice B and Slice C
   tests exercise this ordering against the no-guest fake lifecycle.
4. Persist a cleanup lease before or transactionally with backend creation intent.
5. Persist the durable backend handle immediately after creation and reconcile ambiguous results.
6. Verify staged bytes before start.
7. Persist collection manifest before content release.
8. Persist terminal transcript and teardown classification before ordinary success is visible.

Steps 4 through 8 remain planned product controls. The current lifecycle memory store does not
provide their required durability. In the intended system, a crash between steps produces an
explicit recovery state; missing later state never proves that a grant was unused or a backend was
absent.

## Backend lifecycle

The Supervisor-facing backend contract is equivalent to:

```text
probe → prepare → create → stage → start → wait/inspect
      → terminate → collect → destroy → reconcile
```

Each operation accepts bounded typed data, is idempotent or has a documented reconciliation key,
and returns a stable machine-readable result. Backend handles are never guest or agent authority.

Direct Apple Containerization cannot currently supply the required durable handle after controller
death and is development-only. The libkrun/HVF candidate instead uses one signed VMM process per
attempt. The Supervisor persists and verifies PID, start time, live code identity/CDHash, expected
path, and attempt state before a private inherited pipe authorizes VM start. EOF before
authorization makes the runner exit without starting a VM. Recovery still treats a mismatched
tuple as unresolved, never as authority to kill an arbitrary PID. The OCI/gVisor candidate must
persist an engine-issued identity and exact attempt labels before create, then prove
list/inspect/kill/delete reconciliation across Supervisor, engine/containerd, and outer-VM failure.
A missing response or local process is never authoritative absence.

The VMM's exit status is never authoritative guest completion. The tested runner returned zero for
corrupt roots, guest kernel failures, and missing executables. The first-slice backend uses bounded
dedicated virtio-console ports for source/input and exactly one fixed-cap typed attempt-bound
completion frame containing inline JSON. Before implementation, separate exact source, canonical-
input, completion-frame, and JSON-payload caps plus per-channel role, version, attempt,
registration/plan, profile, length, digest, terminal-status, and commit-trailer semantics are frozen.
Each channel fails rather than resizing; the host continuously drains cap-plus-one and never uses
EOF as completion; the commit trailer is written only after the complete payload. A backend reports
runner lifecycle, guest-reported completion, input integrity, applicable result validation/parser
disposition, and teardown as separate evidence.
Any missing required element blocks ordinary success. The trusted launcher, not the unprivileged
workload, retains completion authority. It remains a distinct process instead of replacing itself
with Bun, verifies complete source/input before child start, gives the child a fixed argv/
environment/cwd/FD manifest without the completion endpoint or node, caps the child result, waits
for exact child-tree termination, and commits completion last. A compromised guest kernel remains
outside what this record attests.

The host runner starts from a fail-closed role-specific descriptor allowlist. After exec it may hold
only the finalized read-only runtime root, dedicated directional port endpoints, and indispensable
runtime/control descriptors. Unexpected database, key, XPC, log, temporary, writable, or inherited
descriptors reject start. Port validation includes hostile virtio control IDs, queues and descriptor
chains plus pinned upstream backpressure, partial-error, directionality, and shared-status behavior;
application framing alone does not admit the device.

`BackendCapabilityReport` identifies exact mechanisms and unsupported controls. A plan proceeds
only when all its required controls match. Capability discovery does not self-certify validation;
the local trust snapshot separately identifies accepted `BackendValidationRecord` digests and
enforces each record's explicit verdict/posture ceiling. A `development-admitted` record never
authorizes `validated-local`.

## Content access

The Supervisor necessarily receives transient access to exact source/input bytes and produced
output while staging/collecting. The security property is scoped capability, not impossible access.

Handles bind installation, epoch, registration, attempt, content identity, direction, operation,
byte limit, expiry, and redemption state. The daemon cannot redeem them. The Supervisor never gets
an original host path or ambient access to the Broker store. In the first inline slice the
Supervisor writes the redeemed exact bytes through bounded attempt-bound ports rather than making
user data a guest block device. Later file staging requires its separate custody/parser gates.

## Evidence transcript

The hash-linked event stream records bounded security-relevant transitions such as:

- plan registration and validation version;
- approval verification and consumption;
- runtime-integrity assessment identity;
- profile/trust/backend identities;
- content-handle redemption and staged digests;
- backend operation results and resource events;
- cancellation/timeout/integrity failure;
- artifact filesystem gate and manifest digest;
- terminate/destroy/reconcile observations;
- terminal classification.

The terminal `EnforcementTranscript` is signed for `capsule.execution.attest`. It is an attributable
Supervisor claim, not independent proof that each event was true or that the host was uncompromised.

## Language and privilege decision

Installed-service evidence supports an unprivileged per-user macOS component with purpose-specific
authenticated XPC and no host-root helper. Host-root execution, a separate-owner host service, and
a privileged host helper are prohibited for v0. If the lead candidate cannot close custody without
one, it fails the v0 profile; any later exception requires a new ADR. Swift remains the preferred
implementation where direct macOS Security, XPC, LocalAuthentication, or lifecycle APIs materially
reduce bindings.

Gate C invalidated direct Swift Containerization as the production backend authority. The
libkrun/HVF candidate has a narrow C ABI and direct macOS process/code-identity needs; the next
slice compares a Go-owned lifecycle adapter with narrow native Security bindings against a small
native launcher surface. The possible Linux OCI/gVisor worker remains a separate boundary. Go
continues to own portable policy/state-machine/backend-contract code. A post-v0 privileged helper
may be considered only under the separate ADR above and must never receive public parsing,
approval, content, or general engine authority.

## Current checkpoint and next boundary

Retained focused evidence includes:

- `TestApprovalAttemptDurableIdempotentAtomicHappyPath` and
  `TestApprovalAttemptConcurrentExactRequestsConverge` for canonical-payload idempotency, exact
  replay/concurrency, and the consumed-approval/created-attempt cross-link;
- `TestApprovalSubmissionTimeBoundaries`, `TestApprovalAttemptFaultAndProcessDeathMatrix`,
  `TestApprovalAttemptReopenRejectsCorruptionAndCrossLinks`, and
  `TestApprovalAttemptHighWaterAndAttemptCapacity` for effective time, capacity, confirmed versus
  indeterminate faults, recovery fencing, reopen, and corruption refusal; and
- the 12 top-level `registeredlifecycle` tests for `AttemptID`-only drive/recovery, complete binding
  revalidation before fake prepare, distinct attempts per registration, startup enumeration,
  exact replay, all fake before/after-effect faults, post-effect restart recovery, and the explicit
  `FakeBackend.CreatesGuest() == false` invariant.

[Proposed ADR-0025](adr/0025-colocate-durable-attempt-lifecycle-state.md) and the
[Phase 2B durable lifecycle plan](PHASE_2B_DURABLE_ATTEMPT_LIFECYCLE_PLAN.md) select the narrow next
implementation boundary: add `AttemptID`-keyed lifecycle records and before/after-effect checkpoints
to the same versioned Supervisor snapshot and transaction domain as the authority records. The
proposal defines startup ownership, reconciliation, repair, capacity, and migration oracles without
creating a registration-keyed execute path or weakening consume/create-before-effect ordering.

This is not implemented durable lifecycle behavior. `registeredlifecycle.MemoryStore` remains the
current non-durable implementation. Archive/compaction and replay retention, real multi-process
locking, rollback/backup, authenticated IPC, production COSE/Swift/Keychain/user-presence signing
and verification, production backend reconciliation, consumers, content, evidence composition,
runtime/backend admission, and public cutover remain blocked. No guest exists in this checkpoint.

## Target acceptance tests

- daemon cannot call backend/launcher directly;
- execute with replacement bytes has no interface;
- daemon-approved unsupported power is rejected independently;
- plan A approval cannot run plan B;
- one grant cannot create two concurrent or sequential attempts;
- crash at every durable/side-effect boundary preserves grant and cleanup invariants;
- forged/stale/cross-job content handle fails;
- all post-create failures reach destroy/reconcile;
- daemon cannot forge success without a valid terminal transcript;
- evidence key cannot sign an Approval object under the accepted profile.
