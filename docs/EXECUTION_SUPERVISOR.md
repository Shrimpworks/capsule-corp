# Execution Supervisor

Status: intended v0 authority and interface; implementation language/privilege is gated by the
macOS feasibility spikes.

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
submitApproval(registrationId, ApprovalGrant) -> approvalReference
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
`false`.

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

## Side-effect ordering

1. Persist registration before issuing the registration ID.
2. Persist grant and validate its bindings.
3. Atomically consume grant and create attempt before backend side effects.
4. Persist a cleanup lease before or transactionally with backend creation intent.
5. Persist the durable backend handle immediately after creation and reconcile ambiguous results.
6. Verify staged bytes before start.
7. Persist collection manifest before content release.
8. Persist terminal transcript and teardown classification before ordinary success is visible.

A process crash between steps produces an explicit recovery state. The system never infers that a
grant was unused or a backend absent solely because a later record is missing.

## Backend lifecycle

The Supervisor-facing backend contract is equivalent to:

```text
probe → prepare → create → stage → start → wait/inspect
      → terminate → collect → destroy → reconcile
```

Each operation accepts bounded typed data, is idempotent or has a documented reconciliation key,
and returns a stable machine-readable result. Backend handles are never guest or agent authority.

`BackendCapabilityReport` identifies exact mechanisms and unsupported controls. A plan proceeds
only when all its required controls match. Capability discovery does not self-certify validation;
the local trust snapshot separately identifies accepted `BackendValidationRecord` digests.

## Content access

The Supervisor necessarily receives transient access to exact source/input bytes and produced
output while staging/collecting. The security property is scoped capability, not impossible access.

Handles bind installation, epoch, registration, attempt, content identity, direction, operation,
byte limit, expiry, and redemption state. The daemon cannot redeem them. The Supervisor never gets
an original host path or ambient access to the Broker store.

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

Gate E compares native Swift, Go plus narrow native bindings, and a hybrid design. Selection
criteria are platform API coverage, privilege, parsing/IPC TCB, memory/runtime footprint, update
complexity, recovery, testability, and developer cost.

The default is an unprivileged per-user component. A privileged helper is introduced only if a
proven required backend operation cannot be performed safely otherwise. That helper accepts only a
sealed launch descriptor and owns no policy, approval, content, or public parser.

## Acceptance tests

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
