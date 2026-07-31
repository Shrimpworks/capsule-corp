# Architecture

## System context

```text
Untrusted agent                      Trusted user / host
      │ proposal                         │ file selection,
      │ MCP / SDK                        │ policy and approval
      ▼                                  ▼
                 Trusted Capsule daemon
                 ├── identity and trust registry
                 ├── proposal validation and planning
                 ├── policy and approval verification
                 ├── input/capability broker
                 ├── runtime/profile registry
                 ├── durable job lifecycle
                 ├── artifact/egress broker
                 └── signed receipt generation
                                │ exact approved plan
                                ▼
                       Isolation backend
                       ├── Apple Container VM
                       └── OCI + gVisor sandbox
                                │
                                ▼
                        Disposable runtime
```

The daemon, trusted-host UI, enrolled policy, trust registry, and isolation launcher are trusted.
Agent-generated source, proposals, guest processes, selected input content, dependencies, and all
guest-controlled output are untrusted. Authentication identifies a caller; it does not make request
content trusted or authorized.

## Components

### Proposal, plan, and approval protocol

An agent proposal declares source, runtime profile, references to previously issued capabilities,
requested capabilities, resource limits, expected artifacts, and audience rules. JSON Schemas in
`schemas/` are canonical.

The protocol does not accept arbitrary container flags, host paths from an agent, or arbitrary
runtime images.

The planner resolves the proposal, immutable inputs, runtime profile, user policy, exact limits,
output audiences, and backend requirements into an immutable execution plan. The trusted host shows
a human-readable view derived from that typed plan. The user signs its digest with a purpose-limited
approval key. Any change produces a new plan that requires new approval.

### Identity and trust registry

Each Capsule installation has an independent P-256 device root identity. v0 represents it as a
`did:key` DID and resolves it entirely offline. Enrollment in the local trust registry—not the DID
syntax itself—grants trust.

The root authorizes separate keys for user approval, daemon receipt assertions, and transport. The
registry stores purpose, issuer, validity, revocation, and replacement information. Unknown DID
methods, algorithms, keys, or purposes are rejected. v0 does not fetch DID documents from a network
or support portable multi-device user identity.

### Policy engine

The policy engine converts requested capabilities into an effective plan before approval. It runs
outside the guest and may deny or explicitly narrow a request while preparing the plan. Nothing may
change after approval.

Trusted user policy owns defaults and ceilings. Requests above a ceiling are rejected rather than
silently clamped. Missing values resolve to explicit user defaults in the plan. A backend that
cannot enforce an exact approved value must refuse execution.

### Input and capability broker

The broker issues and redeems opaque capabilities. A capability is granted by a trusted user or
host action and is bound to an audience, operation, expiry, and job or session. An agent may
reference an existing capability but cannot create one by naming an arbitrary host path.

The first broker supports inline JSON and explicitly selected regular files. A regular-file grant
is copied into private staging, hashed, recorded in a signed snapshot manifest, and exposed
read-only to the guest. The grant identifies the immutable snapshot bytes, not the original path.
Complex parsing occurs inside the guest.

Repository snapshots, directories, archives, network destinations, API operations, secrets, and
subprocesses are later capability types.

### Runtime profile registry

A runtime profile is a signed manifest identifying an immutable runtime plus its approved modules,
hardening settings, limits, image digest, provenance, and supported backend controls. Friendly names
resolve to content digests before execution. The agent cannot supply a registry URL or untrusted
image.

Profile signing authenticates its publisher and protects integrity. A separate review attestation
records the policy or security-review verdict; signing alone does not activate a profile.

The first profile is `bun-data@1`. Node and Deno adapters will follow after the job contract is
proven.

### Isolation backend

The isolation backend enforces the hostile-code boundary. The fake backend exercises lifecycle
without executing guest code. Apple Container is the macOS integration candidate and gives each job
a fresh OCI-compatible Linux lightweight VM on supported Apple silicon. OCI plus gVisor is the Linux
reference backend. Firecracker is a later high-assurance option.

Apple Container and gVisor are independent implementations under one contract. Neither is
authoritative until its exact pinned implementation, host configuration, and runtime profile pass
the mandatory attack corpus. Capsule never runs untrusted Bun directly on the host.

Backends implement a conceptual lifecycle:

```text
prepare → create → stage → execute → collect → destroy
```

### Artifact and egress broker

The broker treats stdout, stderr, structured guest results, filenames, metadata, and files as
untrusted data. It enforces counts, byte limits, file types, regular-file requirements, symlink
denial, parsing or schema validation, and audience policy. Structured results use this same
controlled-object path rather than an implicitly trusted side channel.

Default exposure for user data is full delivery to the user and metadata-only delivery to the
agent. Agent-readable content requires a separate signed, short-lived content-access grant.

### Receipt generator

The receipt records hashes and identities for the approved plan, source, runtime profile, policy,
input snapshots, outputs, approval signer, backend, timings, resource usage, violations, recovery,
and teardown. It is signed by the receipt key. It is evidence of what trusted Capsule components
claim to have enforced, not proof that guest computation was correct.

### Durable lifecycle

Job state, approval consumption, identity status, capability metadata, cleanup leases, and receipt
indexes become durable before Capsule executes hostile code. The expected initial metadata store is
SQLite behind repository interfaces. Sensitive snapshot and artifact content remains in bounded,
access-controlled content storage.

On restart, Capsule reconciles nonterminal jobs, prevents approval replay, attempts backend teardown,
and records recovery or teardown failure explicitly.

## Trust boundaries

```text
Trusted                                Untrusted
─────────────────────────────────────  ──────────────────────────────
Trusted user and host UI               Agent proposal and source
Enrolled device and operational keys   Agent-facing MCP client
Capsule daemon and policy              Guest runtime process
Capability issuance UI/CLI             User-selected input content
Signed profile registry metadata       Third-party dependencies
Isolation launcher                     stdout and stderr
Artifact validators                    Structured guest results
                                       Output files and filenames
```

Important internal boundaries remain:

- Agent-facing credentials cannot call trusted-host issuance, approval, identity-administration, or
  user-content operations.
- The approval key cannot be replaced by the noninteractive receipt key.
- The daemon passes only an exact approved plan to a backend.
- Backend management channels, including VM sockets, are never guest capabilities.
- Content delivery authority is separate from execution authority.

## Portability

There are three independent adapter boundaries:

- **Client adapter:** MCP, CLI, SDK, or HTTP
- **Runtime adapter:** Bun, Node, or Deno
- **Isolation backend:** fake lifecycle, Apple Container, OCI/gVisor, or Firecracker
- **Platform key provider:** Apple Secure Enclave/Keychain or another OS-backed implementation

The protocol should remain stable while these implementations evolve. Source programs may use
runtime-specific APIs and are not required to run unchanged on every runtime.

## Execution state

Every job receives fresh guest state. Immutable verified images and dependency artifacts may be
cached. Previously used guest processes, writable filesystems, module caches, temporary
directories, network state, and connection pools must not cross job trust boundaries.

The control lifecycle is:

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

Once guest resources exist, success, failure, timeout, cancellation, collection failure, and daemon
recovery all pass through destruction. A teardown failure is a distinct terminal classification.

## Detailed design

See [Technical Design](TECHNICAL_DESIGN.md) for signed-envelope requirements, plan contents,
resource-policy semantics, input snapshotting, backend constraints, lifecycle recovery, error
taxonomy, and the ordered implementation plan.
