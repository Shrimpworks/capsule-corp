# Protocol Object Model

Status: target inventory for Phase 2 contract freeze. Phase 2A now has verified passive candidates,
but no target object in this document is an activated or frozen production wire contract unless a
later accepted schema and ADR say so.

## Current scaffold warning

`schemas/job.schema.json`, `schemas/receipt.schema.json`, `schemas/profile.schema.json`, and the
current TypeScript interfaces remain canonical for the repository's current scaffold and tests.
They are not the intended final v0 contracts.

Known mismatches include:

- `Job` combines untrusted request and effective authority;
- agent-controlled guest mount/output paths appear in the public shape;
- network, subprocess, environment, FFI, native addon, and package-install powers are represented
  even though v0 must omit them;
- runtime bundle identity, review, activation, and backend validation are conflated;
- the receipt is daemon-centric rather than composed from Broker/Supervisor evidence;
- job, registration, attempt, artifact custody, and agent-summary semantics are not separated.

Do not extend the current unions to add more future authority. Phase 2 replaces them after the
blocking spikes establish the exact feasible vocabulary.

Phase 2A adds a closed passive `JobProposal` candidate under
[`schemas/candidates/`](../../schemas/candidates/) and minimum `ExecutionPlan` and
`PlanRegistration` candidates under [`schemas/cddl/candidates/`](../../schemas/cddl/candidates/).
They have examples, byte-exact fixtures, and Go/TypeScript decoded views, but no daemon, SDK, MCP,
Broker, Supervisor, or backend integration. The minimum plan intentionally omits unresolved
resource and transport values and cannot authorize execution.

Gate A2 also produced the candidate CDDL contract for `ApprovalGrant` v0 under
[`schemas/cddl/`](../../schemas/cddl/). These candidates fix only their tested shapes; they do not
freeze final field sets or promote experiment code to production. ADR-0019 records the remaining
canonical-wrapper acceptance conditions.

## Object inventory

### Untrusted/public request objects

- `JobProposal`
- `SourceBundleProposal`
- `RequestedLimits`
- `RequestedOutputContract`
- `CancellationRequest`
- `AgentExecutionSummary`

### Content-addressed internal objects

- `SourceManifest`
- `InputSnapshotManifest`
- `PolicyDecision`
- `BackendConfiguration`
- `ArtifactManifest`

### Trust and profile objects

- `RuntimeBundleManifest`
- `ProfileReviewAttestation`
- `ProfileRegistryEntry`
- `BackendCapabilityReport`
- `BackendValidationRecord`
- `TrustMetadataCheckpoint`
- `TrustSnapshot`
- `KeyAuthorization`
- `InstallationManifest`
- `TrustEpochRecord`
- `PreparedUpdate`

### Plan and authorization objects

- `ExecutionPlan`
- `PlanRegistration`
- `ApprovalGrant`
- `ContentAccessGrant` (deferred until the content path exists)
- `ContentHandle`

### Execution and evidence objects

- `ExecutionAttempt`
- `RuntimeIntegrityAssessment`
- `SupervisorEvent`
- `EnforcementTranscript`
- `ExecutionReceipt`

### Shared infrastructure

- `SignedEnvelope`
- `ErrorCode`
- `ViolationRecord`
- `LifecycleTransition`
- posture dimension records

## Semantics

- `JobProposal`: untrusted desired work; never effective authority.
- `ExecutionPlan`: exact immutable resolved authority, content identities, policy, limits, profile,
  backend requirements, installation, and epoch.
- `PlanRegistration`: Supervisor's durable acknowledgement of exact plan bytes.
- `ApprovalGrant`: Broker-attributable one-use authorization for one registration/attempt.
- `ExecutionAttempt`: one realization after atomic grant consumption.
- `EnforcementTranscript`: Supervisor-attributable bounded claim about enforcement/lifecycle events.
- `ExecutionReceipt`: user-facing composition of plan, approval, transcript, artifact, and posture
  evidence.
- `AgentExecutionSummary`: separately minimized fixed response; not a redacted full receipt.

These terms are never interchangeable. `jobId` alone cannot substitute for registration or attempt
identity.

## V0 proposal shape

The final frozen schema is pending. The passive first-slice candidate is narrower than this target
inventory; v0 proposals contain only:

- API/object version;
- source bundle/entrypoint under semantic canonical path rules;
- one trusted runtime-profile selector resolved locally;
- inline JSON or opaque input slot reference;
- requested values for the limit dimensions proven enforceable;
- one or more narrowly declared logical output slots within v0 bounds;
- bounded non-authoritative agent labels.

V0 proposals do not represent:

- network destinations or generic network grants;
- subprocess names or shell commands;
- environment variables or secret names;
- native addons, FFI, macros, inspector, or package installation;
- arbitrary runtime image/registry/kernel/init references;
- host paths or agent-selected guest paths;
- backend flags, mounts, sockets, or security profiles;
- arbitrary JSON Schema machinery;
- directories, repositories, archives, devices, sockets, or special files.

Unknown or future powers fail as unsupported protocol. A boolean set to `false` is not a reason to
keep an unsupported authority in the v0 union.

## Logical resource slots

Agent-facing objects use bounded logical identifiers:

```text
input slot:  primary-data
output slot: transformed-json
```

Capsule assigns backend paths internally, for example:

```text
/capsule/source/main.ts
/capsule/inputs/0/data
/capsule/outputs/0/data
```

Guest paths never carry authority and do not appear in approval as agent-selected fields.

## Signed versus content-addressed

| Object | Identity mechanism | Reason |
| --- | --- | --- |
| SourceManifest | SHA-256 | Registered/approved plan binds exact source identity |
| InputSnapshotManifest | SHA-256; optional Broker signature across independent boundary | Exact bytes are primary; signature attributes Broker claim |
| PolicyDecision | SHA-256 | Plan binds exact policy result/source |
| ExecutionPlan | SHA-256 canonical bytes | Supervisor registration and Broker approval bind it |
| RuntimeBundleManifest | Publisher signature | Distributed artifact origin/integrity |
| ProfileReviewAttestation | Reviewer signature | Independent verdict for exact bundle |
| ProfileRegistryEntry | Local policy identity/digest | Mutable activation must remain separate from bundle |
| BackendValidationRecord | Validation authority signature | Explicit verdict/posture ceiling for exact backend/bundle/host/configuration claims, limitations, expiry, invalidation triggers, and evidence; `development-admitted` never implies `validated-local` |
| KeyAuthorization | Installation/trust-authority signature | Grants purpose to a public key |
| InstallationManifest/epoch | Installation-root-authorized signature | Local component/trust continuity |
| ApprovalGrant | Approval-key signature | Human-authorization claim |
| ArtifactManifest | SHA-256, bound by transcript | Exact collected outputs |
| EnforcementTranscript | Supervisor evidence signature | Enforcement/lifecycle claim |
| ExecutionReceipt | Optional packaging signature plus embedded evidence | Packaging integrity cannot replace embedded claims |
| TUF metadata | TUF role signatures | Distribution, delegation, snapshot, freshness/rollback defenses |
| TrustSnapshot | Local trust-verifier signature | Bounded accepted external/local trust state |
| Witness checkpoint | Witness signature | Independent historical checkpoint only |

## Strict decoding

Every externally supplied JSON object is bounded before normal schema processing:

- raw bytes and accepted media type;
- strict UTF-8 and one document only;
- duplicate-key rejection;
- nesting, collection count, string length, and total decoded limits;
- no trailing data;
- exactly interoperable numeric range or canonical decimal-string counters;
- unknown fields and unsupported object versions rejected.

Signed payloads additionally use canonical-on-wire rules, protected-header allowlists, explicit
type/version/purpose, and mutually exclusive object schemas.

## Agent summary

The default summary is deliberately smaller than the receipt:

```json
{
  "state": "completed",
  "attemptId": "attempt_...",
  "receiptId": "receipt_..."
}
```

Exact artifact names/sizes, guest strings, paths, timing, metrics, and detailed violations are not
included by default. This reduces channel capacity but does not prove confidentiality or eliminate
state/timing leakage.

## Schema-freeze acceptance

- Go, Swift, TypeScript, and JSON Schema agree on shared fixtures where each participates.
- Canonical bytes and ES256 signatures match retained known-answer vectors.
- Invalid Unicode/numbers, duplicate keys, unknown fields/versions, wrong object types, and
  cross-protocol signatures fail.
- Proposal schemas cannot express unsupported authority.
- Registration/attempt/evidence bindings are explicit.
- Limit fields correspond to proven backend semantics or are rejected for that backend.
- Examples and site copy no longer imply the current mixed `Job` is final.
