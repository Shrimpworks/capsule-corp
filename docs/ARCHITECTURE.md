# Architecture

## System context

```text
AI client
   │  MCP / SDK / CLI
   ▼
Trusted Capsule daemon
   ├── job validation and policy
   ├── input/capability broker
   ├── runtime/profile registry
   ├── execution lifecycle
   ├── artifact/egress broker
   └── receipt generation
             │
             ▼
      Isolation backend
             │
             ▼
      Disposable runtime
```

The daemon is trusted. Agent-generated source, guest processes, inputs, dependencies, and all
guest-controlled output are untrusted.

## Components

### Job protocol

The protocol declares source, runtime profile, references to granted inputs, requested capabilities,
resource limits, expected artifacts, and audience rules. JSON Schemas in `schemas/` are canonical.

The protocol does not accept arbitrary container flags, host paths from an agent, or arbitrary
runtime images.

### Policy engine

The policy engine converts requested capabilities into an effective policy. It runs outside the
guest and may deny or narrow a request. Future inputs may include user preferences, organization
policy, backend capabilities, and runtime-profile constraints.

### Input and capability broker

The broker issues and redeems opaque capabilities. A capability is granted by a trusted user or
host action and is bound to an audience, operation, expiry, and job or session. An agent may
reference an existing capability but cannot create one by naming an arbitrary host path.

The first broker will support inline JSON and explicitly selected regular files. Repository
snapshots, network destinations, API operations, secrets, and subprocesses are later capability
types.

### Runtime profile registry

A runtime profile identifies an immutable runtime plus its approved modules and hardening settings.
Friendly names resolve to content digests before execution. The agent cannot supply a registry URL
or untrusted image.

The first profile is `bun-data@1`. Node and Deno adapters will follow after the job contract is
proven.

### Isolation backend

The isolation backend enforces the hostile-code boundary. Development backends may use lightweight
local OS primitives and must be labeled non-authoritative. The first serious security target is an
isolated Linux worker using OCI plus gVisor. Firecracker is a later high-assurance option.

Backends implement a conceptual lifecycle:

```text
prepare → create → stage → execute → collect → destroy
```

### Artifact and egress broker

The broker treats stdout, stderr, structured guest results, filenames, and files as untrusted data.
It enforces counts, byte limits, file types, regular-file requirements, symlink denial, parsing or
schema validation, and audience policy.

Default exposure for user data is full delivery to the user and metadata-only delivery to the
agent. Content access requires a separate grant.

### Receipt generator

The receipt records hashes and identities for source, runtime profile, effective policy, inputs,
outputs, timings, resource usage, violations, and exit classification. It is evidence of what
Capsule enforced, not proof that the guest computation was correct.

## Trust boundaries

```text
Trusted                          Untrusted
───────────────────────────────  ──────────────────────────────
Capsule daemon                   AI-generated source
Policy configuration             Guest runtime process
Capability issuance UI/CLI       Input content
Profile registry metadata        Third-party dependencies
Isolation launcher               stdout and stderr
Artifact validators              Structured guest results
                                 Output files and filenames
```

## Portability

There are three independent adapter boundaries:

- **Client adapter:** MCP, CLI, SDK, or HTTP
- **Runtime adapter:** Bun, Node, or Deno
- **Isolation backend:** local sandbox, OCI/gVisor, or Firecracker

The protocol should remain stable while these implementations evolve. Source programs may use
runtime-specific APIs and are not required to run unchanged on every runtime.

## Execution state

Every job receives fresh guest state. Immutable images and dependency artifacts may be cached.
Previously used guest processes, writable filesystems, module caches, temporary directories, and
connection pools must not cross job trust boundaries.
