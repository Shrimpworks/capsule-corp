# Threat Model

Status: draft

This document describes the intended security boundary. The current scaffold does not yet implement
or satisfy these properties.

## Security objective

Contain an AI-generated JS/TS task so that it can exercise only authority explicitly granted by a
trusted user, host application, or policy, while bounding resource consumption and controlling all
data returned from the guest.

Capsule must also ensure that the plan executed by a backend is the exact immutable plan approved
by the user, that approval cannot be replayed or substituted across devices or jobs, and that every
security-relevant handoff has verifiable identity, integrity, purpose, audience, and freshness.

## Assets

- Host filesystem and user data not granted to the job
- Host and AI-client credentials
- Control-plane process and configuration
- Device root, approval, receipt, and transport private keys
- Local DID trust registry, revocation state, and policy identities
- Execution plans, approval grants, nonces, and approval-consumption state
- Other jobs and their inputs, outputs, and state
- Immutable input snapshots and their retention metadata
- Runtime profile registry and signing trust roots
- Integrity of runtime-profile review attestations
- Integrity of the effective plan, artifact manifest, and execution receipt
- Host availability and bounded resource use
- Confidentiality of guest output not approved for the agent

## Adversaries and untrusted inputs

Capsule assumes the following may be malicious:

- AI-generated source code
- Prompts and model output
- Submitted source bundles
- User-selected input content
- Third-party dependencies
- Runtime parsers and libraries
- Signed envelopes from unenrolled, revoked, expired, or wrong-purpose keys
- DID strings, verification-method identifiers, nonces, timestamps, signatures, and claims supplied
  by clients
- Guest stdout, stderr, filenames, structured results, and artifacts
- A client attempting to request more authority than the user intended
- A client attempting to substitute, replay, race, or reuse an approval
- Malicious OCI images, manifests, archives, registries, or profile metadata before trust
  verification

The trusted user, local operating-system administrator, trusted-host application, Capsule daemon,
enrolled policy, and authoritative backend operator are trusted in the initial model. Defending
against a compromised host kernel, Secure Enclave, hypervisor, or authorized backend operator is out
of scope for local execution.

The agent-facing client is not trusted to perform user actions merely because it is authenticated.
An enrolled signature authenticates control of a key; it does not prove that the signed source is
safe, that a claim is true, or that the signer is authorized for every purpose.

## Trust boundaries

### Client to daemon

Trusted-host and agent-facing credentials have separate authority. An agent-facing client may
propose jobs and inspect policy-approved metadata but cannot issue capabilities, approve plans,
administer identity, or read user-only content.

Every client request remains untrusted after authentication. A syntactically valid or correctly
signed request is not automatically authorized.

### Device identity and signed handoffs

Each installation has a per-device root identity enrolled in a local trust registry. v0 accepts
only offline P-256 `did:key` resolution. The registry assigns purpose, validity, status, and
replacement state to root and operational keys.

Approval, receipt, and transport signing purposes are separated. Unknown DID methods, resolver
plugins, algorithms, keys, purposes, audiences, installations, or envelope versions fail closed.
Job authorization never depends on live network DID resolution.

Signed objects bind their type, version, subject digest, issuer, purpose, audience, installation,
nonce, and validity period. Verification precedes state transition. A single-use grant is consumed
atomically before execution.

### Proposal to execution

The daemon resolves all source, input, profile, policy, limit, output, audience, and backend values
into one immutable execution plan. The trusted host displays a human-readable summary derived from
that model. The user approves its exact digest.

No value may be defaulted, clamped, substituted, or widened after approval. A changed plan requires
new approval. A backend that cannot enforce the exact plan refuses the job.

### Capability issuance

Only trusted user or host actions may turn a host resource into a capability. An agent may reference
an issued capability but may not mint one from a path, URL, environment variable, or identifier it
invented.

A regular-file capability references a private immutable snapshot. The broker safely opens a
user-selected regular file, copies and hashes its bytes, detects material mutation, records a
snapshot manifest, and stages the snapshot read-only. The original host path is not part of the
agent or guest contract.

### Daemon to isolation backend

The daemon passes a fully resolved effective policy to the launcher. Backend configuration is
generated from trusted code and must not contain guest-controlled shell interpolation, mount flags,
image references, or seccomp rules.

The backend receives only an approved plan. Runtime images, kernels, init environments, management
channels, OCI configuration, and backend binaries are pinned trusted inputs. VM sockets, container
engine sockets, and launcher control channels are not guest capabilities.

### Guest to host

The guest is hostile. Its syscalls, filesystem access, process creation, IPC, and network access are
controlled externally. Runtime permission systems are supplemental only.

### Guest output to user or agent

All guest-controlled output crosses an egress broker. Content delivery is separate from metadata
delivery and follows an explicit audience policy.

Execution approval and content-access approval are distinct signed grants. Structured results and
logs do not bypass artifact policy.

## Mandatory security properties

### Identity, approval, and integrity

- A hash is never treated as proof of origin without an independently trusted expected value or
  signature.
- Device root keys are non-exportable where supported and operational signing purposes are
  separated.
- Approval requires explicit user action for every v0 execution plan.
- Approval binds the exact plan digest, installation, audience, job, purpose, nonce, and expiry.
- Unknown, revoked, suspended, replaced, expired, wrong-purpose, wrong-audience, and
  wrong-installation keys or grants are rejected.
- Single-use approvals and content grants are consumed atomically and cannot be replayed after
  daemon restart.
- Human-readable approval views and signed plans derive from the same typed model.
- Network DID resolution, arbitrary DID methods, dynamic resolver plugins, and unpinned contexts
  are not used in v0 authorization.
- Signatures authenticate enrolled assertions; policy still validates the requested action and
  resource.

### Isolation

- The guest cannot read or write arbitrary host paths.
- Inputs are immutable snapshots staged with the exact granted read-only access.
- The guest never receives the original host path or a live user-file mount.
- The runtime root is immutable during a job.
- Scratch and output storage are isolated and size-limited.
- Host environment variables, open descriptors, sockets, and credentials are not inherited.
- Guest processes cannot signal, inspect, or attach to host or other-job processes.
- Guest state is destroyed after completion or cancellation.
- A teardown failure is reported distinctly and is never treated as successful cleanup.
- Capsule does not execute untrusted Bun directly on the host.

### Network and IPC

- Network is externally denied by default.
- DNS, loopback, link-local, metadata services, and Unix sockets are included in the denial.
- v0 backend configuration provisions no usable network interface rather than depending on broken
  connectivity.
- Future network access must use either an isolated network policy or a broker that revalidates
  destinations and redirects.
- Container-engine, agent, SSH, credential, display, and host-service sockets are never inherited.
- VM management and vsock channels are minimized, authenticated, and tested as hostile guest attack
  surfaces.

### Resource control

- Wall time and CPU time are bounded independently.
- Memory, PIDs, temporary storage, output count, and output bytes are bounded externally.
- Output flooding cannot exhaust daemon memory.
- Cancellation terminates the complete guest process tree.
- Trusted user policy supplies explicit defaults and ceilings.
- Requests above user ceilings are denied rather than silently clamped.
- The exact approved values are enforced or the backend refuses execution.
- Request, source, inline input, snapshot input, log, per-artifact, and total artifact bytes are
  independently bounded.

### Runtime hardening

- Bun automatic installation and `.env` loading are disabled.
- Native addons, FFI, dynamic package installation, macros, subprocesses, and inspector access are
  disabled unless a future profile explicitly grants them.
- Profiles are selected from a trusted registry and resolved to immutable digests.
- Active profiles have pinned images, dependency manifests, SBOMs, publisher signatures, and
  separate review attestations.
- Profile signing authenticates origin and integrity but does not replace policy review.

### Egress

- stdout and stderr are capped and not automatically exposed in full to the agent.
- Only declared artifacts from the dedicated output volume are considered.
- v0 artifacts must be regular files with exact paths, allowed types, and bounded sizes.
- Symlinks, hard-link tricks, device files, sockets, FIFOs, sparse-file abuse, and archives are
  rejected in v0.
- Agent content access requires an audience grant separate from user delivery.
- Structured results use the same controlled validation and exposure path.
- Artifact and log metadata cannot disclose original host paths or unbounded guest strings.

### Durability and recovery

- Job, approval-consumption, trust-registry, capability, cleanup, and receipt state is durable before
  hostile execution is enabled.
- State transitions are atomic, monotonic, and idempotent.
- After restart, nonterminal jobs are reconciled and their backends are reaped or explicitly
  reported as unresolved.
- Missing backend state is not assumed to prove destruction.
- Snapshot and artifact retention follows bounded access-controlled policy.

## Non-guarantees

Capsule does not prove:

- That guest code performs the requested task correctly
- That a signed assertion or verifiable credential contains a true claim
- That an approved output does not contain copied or encoded input data
- That secret-pattern redaction can identify every sensitive value
- That a supported runtime or host kernel contains no unknown vulnerability
- That Apple Container, gVisor, a hypervisor, a Secure Enclave, or a cryptographic implementation
  contains no unknown vulnerability
- That a signature protects against compromise of its private key or an already trusted signer
- That a receipt proves behavior independently of the trusted component that signed it
- That source written for one runtime behaves identically on another
- That inputs supplied through an AI client's normal attachment flow remain outside model context

## Abuse cases and required tests

The authoritative backend must cover at least:

| Category | Cases |
| --- | --- |
| Filesystem | traversal, absolute paths, symlinks, hard links, `/proc`, home and credential reads |
| Process | fork bomb, worker creation, signals, inspector activation, orphan processes |
| Network | TCP, UDP, DNS, loopback, IPv6, Unix sockets, metadata endpoints |
| Runtime | native addons, FFI, Wasm abuse, dynamic imports, Bun auto-install, `.env` loading |
| Resources | busy loop, heap OOM, native allocation, PID exhaustion, disk fill, output flood |
| Artifacts | undeclared paths, oversized output, sparse files, devices, FIFOs, archive bombs |
| Identity | unknown DID methods, key substitution, wrong purpose, revocation, rotation, root replacement |
| Approval | plan mutation, replay, stale nonce, wrong audience, wrong installation, expiry, double execution |
| Handoff | digest mismatch, signature mismatch, type confusion, canonicalization differences, partial writes |
| Isolation | cross-job state, cached writable data, inherited descriptors, environment leakage, management channels |
| Recovery | cancellation races, daemon crash, orphan backend, replay after restart, teardown failure |
| Protocol | unknown fields, duplicate identifiers, invalid limits, unsupported capabilities |

## Security tiers

- **Development:** fake backend or unvalidated local sandbox; useful for development, not an
  authoritative claim.
- **Authoritative local:** a pinned Apple Container or Linux/gVisor configuration that has passed
  the mandatory local attack corpus.
- **Hosted hardened:** gVisor or stronger boundary plus signed profiles, tenant isolation,
  monitoring, and hosted operational controls.
- **High assurance:** future microVM backend with dedicated tenant boundaries.

Every receipt and user-facing status should identify the backend tier used.

## Open questions

- Canonical signed-envelope representation and maintained Go and TypeScript implementations
- Exact approval user-presence and short-lived session policy
- Device-root recovery and trusted replacement ceremony
- Windows local-isolation strategy
- Safe handling of user-approved network access
- Retention and encryption policy for receipts and staged artifacts
- Quantitative performance and resource budgets

## Severity calibration

### Critical

- Execution of an unapproved or substituted plan through approval bypass, signature confusion, or
  replay
- Guest escape to arbitrary host code execution or unrestricted host filesystem and credential
  access
- Agent-facing authority to issue file capabilities, approve plans, or read user-only content
- Cross-job or cross-user compromise in a hosted authoritative backend

### High

- Reliable network or host-service access when the approved plan denies network
- Acceptance of a revoked, wrong-purpose, wrong-installation, or wrong-audience approval
- Live host-file mounts or path races that expose content beyond the selected snapshot
- Resource-control bypass capable of materially affecting host availability
- Teardown failure hidden as success, leaving hostile execution active

### Medium

- Receipt, metric, or artifact metadata integrity failures that materially weaken auditability
  without expanding guest authority
- Bounded disclosure of guest output to the wrong audience
- Denial of service confined to one local job and its configured resource budget
- Profile provenance or review-status confusion that does not activate an untrusted profile

### Low

- Non-sensitive diagnostic inaccuracies with no effect on authorization, isolation, egress, or
  cleanup
- Development-only availability failures that cannot execute guest code or alter production
  evidence
- Documentation or metadata defects that do not overstate an implemented security tier
