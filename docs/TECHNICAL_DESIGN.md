# Technical Design

Status: agreed foundation; contracts remain draft until Phase 0 exit criteria pass.

This document maps the intended v0 system across the public protocol, trusted control plane,
identity and approval model, capability brokers, execution backends, controlled egress, and
evidence. The JSON Schemas in `schemas/` remain the canonical wire contracts. Where this document
describes a contract that does not yet exist in `schemas/`, it is a requirement for the Phase 0
contract-freeze work rather than implemented behavior.

## Reference workflow

The first complete Capsule workflow is intentionally narrow:

1. An agent proposes a one-shot TypeScript data transformation.
2. A trusted user supplies inline JSON or selects one regular file.
3. Capsule validates the proposal, snapshots granted inputs, resolves policy and a runtime profile,
   and produces an immutable execution plan.
4. The trusted host presents a human-readable plan and obtains explicit user approval for its exact
   digest.
5. Bun executes in a disposable Linux sandbox with no network, ambient environment, subprocesses,
   native addons, FFI, macros, inspector, or package installation.
6. The guest writes one or more declared JSON, JSONL, CSV, or text outputs.
7. The egress broker validates every guest-controlled channel. The user receives approved content;
   the agent receives metadata by default.
8. Capsule destroys the guest environment and emits a signed execution receipt that records the
   plan, enforcement evidence, outputs, violations, and teardown result.

The task remains the public abstraction. Containers, lightweight VMs, runtimes, and process trees
are backend details.

## Design invariants

- An agent can propose authority but cannot grant it.
- An approval is valid for one immutable execution-plan digest, audience, installation, and
  validity period.
- Hashes establish byte identity; signatures establish control of an enrolled key; policy
  determines whether that key may authorize the requested action.
- The daemon never executes an approved plan that differs from the plan shown to the user.
- Unknown principals, DID methods, keys, algorithms, capabilities, profiles, fields, transitions,
  output types, and backend controls fail closed.
- A guest never receives a live user-selected host path.
- User-defined resource values are not silently clamped or rewritten.
- Every guest-controlled output crosses one egress boundary.
- Every execution reaches teardown, including failure, timeout, and cancellation paths.
- A receipt describes what trusted components claim to have enforced; it does not prove guest-code
  correctness.

## Principals and trust roots

### Agent client

The agent-facing client is an untrusted proposer. It may submit source, requested capabilities,
limits, declared outputs, and references to capabilities already issued by a trusted host. It
cannot issue capabilities, approve execution plans, alter the trust registry, select arbitrary
images, or receive user-only content.

### Trusted user and host

The trusted host is responsible for user interaction, file selection, policy configuration,
capability issuance, plan presentation, approval, and user-audience content delivery. Explicit
approval is required for every v0 execution plan.

### Capsule installation identity

Each installation has an independent device root identity. On supported Apple hardware, its P-256
private key should be generated and retained in the Secure Enclave. Other platforms implement the
same interface with the strongest available non-exportable OS or hardware-backed key store.

The initial public identifier is a P-256 `did:key` DID:

```text
did:key:<multibase-encoded-public-key>
```

`did:key` is resolved locally and deterministically. Capsule v0 does not perform network DID
resolution, load resolver plugins, or accept arbitrary DID methods. The DID identifies a key; it
does not make that key trusted. Enrollment in the local trust registry provides trust.

The root key is used sparingly to authorize distinct operational keys:

- **Approval key:** signs execution and content-access grants; user-presence protected.
- **Receipt key:** signs daemon assertions and execution receipts; noninteractive and narrowly
  scoped.
- **Transport key:** authenticates a local or remote Capsule session when asymmetric transport
  identity is needed.
- **Delegation key:** deferred; may enroll another device or worker under explicit user approval.

Key separation prevents compromise of an always-running receipt signer from automatically becoming
execution-approval authority.

### Trust registry

The local trust registry stores enrolled device and operational identities, verification methods,
authorized purposes, validity periods, status, issuer, and replacement relationships. Status is at
least `active`, `suspended`, `revoked`, or `replaced`.

`did:key` cannot update or deactivate itself. Capsule therefore owns application-level operational
key rotation and revocation. A lost or replaced device root becomes a new DID and requires trusted
reenrollment. Multi-device identity portability is deferred; v0 identities are per installation.

## Signed objects and verification

Security-relevant handoffs use typed, versioned, signed envelopes. Planned object types include:

- `RuntimeProfileManifest`
- `InputSnapshotManifest`
- `ExecutionPlan`
- `KeyAuthorization`
- `ApprovalGrant`
- `ArtifactManifest`
- `ContentAccessGrant`
- `ExecutionReceipt`

Each signed envelope binds:

- schema and object type
- issuer and verification-method identifiers
- intended audience and Capsule installation
- subject digest and job identifier where applicable
- capability or signing purpose
- nonce
- issue and expiry times
- signature algorithm and key identifier
- signature over domain-separated canonical bytes

The contract-freeze work must select one well-reviewed canonicalization and envelope standard that
has maintained Go and TypeScript implementations. Capsule will not invent a cryptographic
primitive. SHA-256 is the initial content digest. P-256 is the initial device-signing algorithm
because it can be hardware-backed by the Apple Secure Enclave.

Verification order is fail closed:

1. Enforce request and envelope byte limits.
2. Parse against a strict versioned schema.
3. Resolve only an allowed verification method without network access.
4. Confirm enrollment, purpose, status, audience, installation, nonce, and validity period.
5. Verify canonical subject bytes, digest, and signature.
6. Evaluate Capsule policy for the requested state transition.
7. Atomically consume single-use grants before execution begins.

Transport security and object security are complementary. Local clients use an OS-protected local
transport and peer identity. A future remote transport uses mutually authenticated encryption.
Signed envelopes remain independently verifiable after transport.

## Prepare, approve, execute

```text
untrusted proposal
      │
      ▼
strict validation ──deny──► policy-denied result
      │
      ▼
resolve profile + inputs + user policy + backend controls
      │
      ▼
immutable ExecutionPlan + digest + human-readable summary
      │
      ▼
explicit user approval over exact plan digest
      │
      ▼
single-use ApprovalGrant
      │
      ▼
reverify + atomically consume grant + execute
```

### Proposal

A proposal contains requested source, inputs, capabilities, limits, result contracts, and runtime
selection. It is untrusted even when it arrives over an authenticated client connection.

### Execution plan

The planner resolves all aliases and implicit values before approval. An execution plan includes:

- canonical source manifest and digest
- immutable input snapshot identities and digests
- runtime profile name, version, image digest, and signing identity
- exact effective capabilities
- exact resource and output limits
- declared result and artifact contracts
- user and agent exposure policy
- backend requirements and minimum security tier
- installation, audience, nonce, and expiry

No policy or backend default may be introduced after the plan is signed. If a backend cannot
enforce the plan exactly, execution is refused.

### Human-readable approval

The trusted host displays a stable summary of source identity, inputs, runtime, capabilities,
limits, outputs, audiences, and backend tier. The user never approves an opaque digest alone.
Approval UI and serialized plan are generated from the same typed model.

## Policy and user-owned limits

The user owns resource policy. A user policy defines explicit defaults and ceilings. An agent may
request values at or below those ceilings. Missing values resolve to user defaults during planning.
Requests above a ceiling are rejected instead of silently clamped.

The execution plan contains the final exact values shown to and approved by the user. The daemon
does not modify them after approval. The backend either enforces them or rejects the job.

For v0:

- network requests are denied
- subprocess requests are denied
- environment inheritance is denied
- native addons, FFI, macros, inspector access, and package installation are denied
- limits exist for request bytes, source files and total source bytes, inline JSON bytes, input
  bytes, wall and CPU time, memory, PIDs, scratch bytes, log bytes, artifact count, per-artifact
  bytes, and total artifact bytes

The policy decision records requested values, resolved values, denial codes, and the trusted policy
identity that supplied each default or ceiling.

## Input capability model

### Inline JSON

Inline JSON is validated for syntax, nesting, item counts, and total encoded bytes before planning.
It is canonicalized and hashed into the execution plan.

### Regular-file snapshots

“Arbitrary file input” means opaque content from an explicitly selected regular file. It does not
mean an arbitrary host path or filesystem object.

The trusted input broker:

1. Receives a file selection from a trusted user action.
2. Opens the selected object without following links using platform-appropriate safe-open
   primitives.
3. Verifies it is a regular file and satisfies user-defined byte limits.
4. Copies from the opened descriptor into private staging while hashing the copied bytes.
5. Detects material mutation during the copy and rejects an unstable selection.
6. Records a signed snapshot manifest and issues an opaque, short-lived, job-bound capability.
7. Stages the snapshot read-only for the guest without exposing the original host path.

The snapshot bytes—not the mutable source path—are the granted resource. Complex file parsing occurs
inside the guest. The trusted daemon does not extract archives or invoke rich document, image,
spreadsheet, or media parsers.

v0 rejects directory capabilities, symlinks, devices, sockets, FIFOs, host-side archive extraction,
and other special filesystem objects. Additional resource kinds require new explicit capability
contracts.

## Runtime profiles and source

Source bundles use a deterministic file manifest with path, byte length, and content digest for
each file plus an overall digest. Paths are relative, normalized, bounded, and collision checked.
The source contract includes both per-file and total-bundle limits.

Runtime profiles are immutable signed manifests. An active profile identifies:

- runtime and exact version
- OCI image digest
- kernel and init identity where backend-specific
- curated dependency versions and review status
- SBOM and provenance references
- hardening settings
- default and maximum limits
- supported backend controls
- profile signer and independent review attestation

Signing authenticates the publisher and protects manifest integrity. A separate review attestation
records the security-review verdict, following the same separation as code signing and
notarization. A profile remains `draft` until its immutable artifacts and required tests exist.

## Isolation backends

All backends implement the same lifecycle:

```text
prepare → create → stage → execute → collect → destroy
```

### Fake backend

The fake backend exercises planning, state transitions, cancellation, persistence, egress, and
receipt generation without executing guest code. It is a development backend.

### Apple Container backend

Apple Container is the initial macOS integration candidate. It runs OCI-compatible Linux containers
inside per-container lightweight virtual machines on Apple silicon. Capsule uses one fresh VM per
job and destroys the VM after collection.

The backend must pin the Apple Container implementation, Linux kernel, init environment, runtime
image, and OCI manifest. It must provision no network interface, expose no host bind mounts or host
sockets, use a read-only root, provide dedicated input/scratch/output filesystems, and enforce
resource limits externally.

Apple Container remains a development-tier backend until the exact implementation and host version
pass the mandatory adversarial corpus. Platform support is initially limited to supported Apple
silicon and macOS versions.

### OCI plus gVisor backend

OCI plus gVisor is the Linux reference isolation target. It runs each job in a separate gVisor
sandbox with externally configured namespaces, filesystem exposure, cgroups, networking, and
teardown.

It remains non-authoritative until its exact runtime, profile, host configuration, and attack-corpus
results support the claim. Linux CI or dedicated workers exercise this backend; developers do not
need to perform ordinary repository work on Linux.

### Backend independence

Apple Container and gVisor are separate enforcement implementations under one backend contract.
Capsule does not depend on stacking gVisor inside an Apple lightweight VM. A backend tier is
evidence-based and recorded in every plan and receipt.

Capsule does not provide a backend that executes untrusted Bun directly on the host.

## Job lifecycle

Control state and terminal result are separate concepts:

```text
submitted
  → validating
  → awaiting-approval
  → prepared
  → staging
  → running
  → collecting
  → destroying
  → terminal
```

Pre-execution validation or policy denial may transition directly to a terminal result. Once guest
resources exist, every path passes through `destroying`.

Terminal classifications include:

- `succeeded`
- `guest-failed`
- `policy-denied`
- `approval-expired`
- `cancelled`
- `timed-out`
- `resource-exhausted`
- `backend-failed`
- `egress-rejected`
- `teardown-failed`

State transitions are atomic, monotonic, and idempotent. Cancellation terminates the complete guest
environment, not only the first process. Teardown failure is visible and cannot be reported as
ordinary success.

## Persistence and recovery

The lifecycle and trust stores are defined behind repository interfaces. A fake implementation may
start in memory, but Capsule must use durable local metadata before executing hostile code.

The initial durable store is expected to be SQLite for:

- jobs and state transitions
- plan digests and approval consumption
- trusted identities and key authorizations
- capability metadata and expiry
- runtime-profile metadata
- artifact descriptors
- receipt indexes
- cleanup leases and backend handles

Input snapshots, artifact content, and detailed receipts live in access-controlled content storage,
not as unbounded database blobs.

On restart, Capsule finds nonterminal jobs, prevents grant replay, reconciles backend state, attempts
teardown, and emits an explicit recovery or teardown result. Recovery never assumes an absent
backend handle means successful destruction.

## Controlled egress and receipts

stdout, stderr, structured results, artifact names, metadata, and files are untrusted. The broker
applies count, byte, type, parsing, path, regular-file, link, and audience policy before exposing
anything.

Structured results use the same controlled-object path as artifacts rather than receiving an
implicitly trusted channel. Logs are bounded artifacts with explicit truncation and exposure
metadata.

Default v0 exposure:

- user: approved full content
- agent: metadata only

Agent-readable content requires a separate, signed, short-lived `ContentAccessGrant`. An execution
approval does not imply content-access approval.

The execution receipt records at least:

- proposal, plan, source, profile, input, output, and policy digests
- approval issuer and verification method
- backend identity and security tier
- exact effective capabilities and limits
- timestamps, exit classification, metrics, violations, and truncation
- artifact descriptors and manifest root
- cancellation and recovery events
- teardown outcome
- receipt signer and signature

An optional hash-linked transparency log may be added later for historical tamper evidence. v0 does
not require a blockchain or distributed consensus.

## Client and API separation

The daemon distinguishes trusted-host and agent-facing credentials and endpoints.

Trusted-host operations include:

- file selection and capability issuance
- policy administration
- identity enrollment, rotation, and revocation
- plan approval
- user-audience content access

Agent-facing operations include:

- list profiles
- prepare a proposal
- submit an already approved plan
- inspect status and policy-approved metadata
- cancel work it is authorized to manage
- request, but not grant, content access
- read execution receipts without user-only content

Local production communication should use an OS-protected local transport such as a Unix-domain
socket with peer checks. Loopback HTTP remains development-only unless it gains equivalent
authentication and authorization.

## Error and violation taxonomy

Stable machine-readable codes distinguish:

- malformed or unsupported protocol
- authentication failure
- untrusted, revoked, or wrong-purpose signing key
- stale, replayed, expired, or wrong-audience approval
- policy denial
- profile unavailable or digest mismatch
- backend capability mismatch
- guest failure
- resource-limit violation
- artifact or output validation failure
- cancellation
- recovery or teardown failure

Human-readable messages are bounded and must not echo secrets, arbitrary guest output, or host
paths.

## Verification strategy

Phase 0 requires shared positive and negative fixtures across JSON Schema, Go, and TypeScript. It
also requires deterministic digest and signature test vectors.

Every backend candidate runs the same attack corpus:

- host path and credential reads
- traversal, symlink, hard-link, device, FIFO, sparse-file, and archive attacks
- environment, inherited descriptor, socket, and metadata-service access
- TCP, UDP, DNS, loopback, IPv6, Unix socket, and backend management-channel access
- workers, signals, inspectors, subprocesses, orphans, and process exhaustion
- CPU, memory, disk, PID, log, and output exhaustion
- Bun auto-install, `.env`, native addon, FFI, macro, and dynamic-import abuse
- undeclared, oversized, malformed, or audience-violating output
- cancellation, daemon crash, restart recovery, and cross-job state leakage
- signature substitution, unknown algorithms, replay, stale approval, wrong audience, wrong
  installation, and revoked identity

No backend or profile is called authoritative until the required tests pass for its exact pinned
identity.

## Ordered implementation plan

### Foundation A: contract freeze

- Add schemas for execution plans, approvals, capability snapshots, signed envelopes, errors, and
  lifecycle state.
- Tighten receipt capability and limit types.
- Define canonicalization, digest, signature, nonce, and replay rules.
- Add shared Go and TypeScript fixtures and test vectors.

### Foundation B: trust and planning

- Implement the local trust registry and platform key interface.
- Implement offline P-256 `did:key` resolution.
- Implement user-policy resolution with no silent clamping.
- Implement proposal validation and immutable execution-plan generation.
- Implement human-readable plan summaries from the typed model.

### Foundation C: lifecycle

- Implement the durable job state machine and fake backend.
- Implement approval consumption, cancellation, recovery, teardown, and signed receipts.
- Exercise the complete lifecycle without guest execution.

### Reference execution

- Implement the Bun runtime adapter and draft profile verification.
- Implement the Apple Container development backend.
- Execute inline source and inline JSON with all optional authority denied.
- Add controlled logs and one validated structured artifact.

### File-to-artifact vertical slice

- Implement trusted regular-file selection and snapshot capabilities.
- Stage immutable read-only input plus dedicated scratch and output storage.
- Complete CLI and MCP flows over the same daemon contracts.

### Authoritative validation

- Implement the OCI plus gVisor backend.
- Run the shared attack corpus and performance benchmarks against both backend candidates.
- Publish only the security claims supported by pinned artifacts and retained evidence.

## Deferred work

- Network and API brokers
- directories, repository snapshots, archives, and patch capabilities
- custom dependency preparation
- Node and Deno runtime adapters
- portable multi-device user identity
- externally resolved DID methods and general Verifiable Credentials
- hosted scheduling and multi-tenancy
- Firecracker
- public transparency or blockchain anchoring

## Remaining Phase 0 decisions

- Canonical signed-envelope representation and maintained Go/TypeScript libraries
- Exact user-presence requirements and approval-session policy
- Retention and encryption policy for input snapshots, artifacts, and receipts
- Device-root recovery and trusted replacement ceremony
- Windows local-isolation backend
- Quantitative startup, memory, and throughput targets

## External design references

- [W3C Decentralized Identifiers v1.0](https://www.w3.org/TR/did-core/)
- [W3C Verifiable Credentials Data Model v2.0](https://www.w3.org/TR/vc-data-model/)
- [`did:key` Method](https://w3c-ccg.github.io/did-key-spec/)
- [Apple app code signing](https://support.apple.com/guide/security-pdf/app-code-signing-process-sec3ad8e6e53/web)
- [Apple Secure Enclave key protection](https://developer.apple.com/documentation/Security/protecting-keys-with-the-secure-enclave)
- [Apple Container](https://github.com/apple/container)
- [Apple Containerization](https://github.com/apple/containerization)
- [gVisor security model](https://gvisor.dev/docs/architecture_guide/security/)
- [Ethereum transactions](https://ethereum.org/developers/docs/transactions)
- [Ethereum Merkle Patricia tries](https://ethereum.org/developers/docs/data-structures-and-encoding/patricia-merkle-trie/)
