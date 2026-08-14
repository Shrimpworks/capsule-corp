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

- `GovernedDenoCoreC2bPassiveBinding` (passive conformance evidence only; version 1)
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

## Field-authority classification

Before any target object freezes, every field in that object must appear in a closed,
machine-readable authority manifest. Each entry identifies:

- the object version and exact field path;
- the only role permitted to originate the value;
- the component responsible for structural validation and any trusted resolution;
- whether the field changes authority, narrows authority, selects one predeclared alternative, or
  records non-authoritative data/evidence;
- whether trusted approval UI must render it and from which registered bytes;
- whether it may contain user content or guest-controlled material;
- the exact digest, signature, registration, or durable record that binds it; and
- the required fail-closed behavior for unknown or unsupported values.

The initial classification vocabulary must at least distinguish untrusted proposal data,
approval-visible authority, Supervisor-authoritative state, content capabilities, runtime
authority, backend authority, and evidence-only observations. A display label or identifier does
not acquire authority merely because it is classified.

Strict decoders already reject fields absent from an object's schema. The authority manifest adds
a separate development invariant: adding a schema/CDDL field without deciding who may assert it,
what it can authorize, how it is bound, and whether it is shown during approval fails repository
verification. The manifest should begin with the passive candidate objects and expand only through
the coordinated versioned migrations that add new target fields; it must not extend the deprecated
mixed `Job` model.

The passive implementation now retains 1,203 closed field classifications across 60 selected
targets and 95 profiles; the generated authority manifest is the source of truth for current
counts. It covers the `JobProposal`, `ExecutionPlan`, `PlanRegistration`, and `ApprovalGrant`
candidates;
the passive approval reference, attempt reference, durable approval record, and immutable attempt;
the TypeScript approved-byte object family plus its future-plan source-binding projection; the
single-member MJS `SourceManifest` including its nested path/digest/length fields; the passive
Source Validator v0 request, result, engineering-candidate, and artifact-profile records; the
v1 request, result, resource-policy, process-profile, artifact-profile, and consumer projection;
and the closed `GovernedDenoCoreC2bPassiveBinding` v1 plus its exact C1, C2A, fork-supplement,
build-evidence, artifact, dependency, limitation, status, and next-gate descendants.
Repository verification compares those classifications directly with the current JSON Schema,
numbered CDDL maps, and version-marked Go passive structs. Focused mutations prove rejection of a
missing field (including a nested member field), unknown classification, duplicate path, stale
object version, and classification for a field absent from its canonical target. The durable
approval envelope digest is classified
as evidence-only and never as replay or ledger authority.

Every C2B binding field is supplied only by the repository's passive conformance-fixture generator,
validated by the independent strict Go and TypeScript C2B validators, and retained only as the
versioned repository conformance fixture. Current consumers are those validators and repository
schema/known-answer tests. The only named later eligible consumer is a separately authorized
composed-profile/owned-guest task after both draft dependency PRs merge and every dependency and
artifact identity is reverified. That eligibility does not activate a consumer, admit a runtime or
profile, or grant execution authority. Any changed dependency head, tree, or artifact requires a
new binding/evidence identity; the v1 bytes remain immutable.

The Source Validator result fields are evidence-only observations, and its fixture artifact profile
does not enroll an executable. This coverage remains pre-freeze and unwired. It does not classify
or extend the deprecated mixed
`Job`, activate an endpoint or consumer, classify the illustrative future `ExecutionPlan` v1 as a
current object, resolve the TypeScript transformation owner, admit a runtime/backend, or authorize
execution. Any versioned cutover must update its canonical target and authority classifications in
the same reviewed change.

Accepted ADR-0036 now defines the
[passive Source Validator v1 implementation boundary](MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1.md)
with separate daemon and Approval-Broker request/result/process/artifact-profile families plus an
evidence-derived reactive-resource-policy family. R1 now retains those canonical layouts, known
answers, and classifications while preserving all V0 targets unchanged and rejecting cross-role/
cross-version reuse. Active threshold/cadence/baseline/overshoot/kill-latency values remain unset
until the signed R4 corpus.

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
- `GovernedDenoCoreC2bPassiveBinding`: immutable evidence-only reconciliation of unchanged C1/C2A
  known answers with one exact fixed-fixture fork/build candidate; never a runtime/profile
  admission, plan, registration, approval, descriptor, or execution authorization.

These terms are never interchangeable. `jobId` alone cannot substitute for registration or attempt
identity.

Proposed ADR-0030 refines ADR-0026 with a passive versioned object family and atomic cutover plan
without changing the current candidate: an original-authoring source manifest, an executable
JavaScript source manifest, and an ordered transformation-record set all bind into a future plan v1
before registration. Approval continues to bind the exact plan digest; the runtime may receive only
the executable role and may never transform from an original-only digest after approval. Proposed
ADR-0032 selects a separately enrolled Source Preparer and immutable source store as the
conditional later TypeScript topology. Accepted ADR-0034 removes that path from the first-release
critical path: plan v0 binds one canonical single-member `SourceManifest` and exact pass-through
`main.mjs` bytes. The byte/manifest foundation is passive; JobProposal narrowing and plan
construction remain blocked by the M1 source-language parser-boundary hold. TypeScript remains
unimplemented and would still require the atomic plan-v1
cutover if later selected.

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
| SourceManifest | SHA-256 | Registered/approved plan binds exact source identity; Proposed ADR-0026 later splits authoring and executable roles |
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
- Every target field has a closed authority-manifest entry, and repository verification rejects an
  unclassified field or an unknown classification value.
- Registration/attempt/evidence bindings are explicit.
- Limit fields correspond to proven backend semantics or are rejected for that backend.
- Examples and site copy no longer imply the current mixed `Job` is final.
