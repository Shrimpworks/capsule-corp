# Roadmap

The roadmap is ordered by risk reduction rather than feature count.

## Phase 0: specification

- Freeze the draft job, profile, and receipt schemas.
- Define capability issuance and audience-aware egress protocols.
- Complete the threat model and explicit security claims.
- Define backend security tiers and the adversarial corpus.
- Record architecture decisions.

## Phase 1: minimal Bun runner

- Trusted daemon lifecycle
- Bun runtime adapter
- Digest-pinned runtime profile
- Inline source and JSON input
- Denied network and environment inheritance
- Hard resource and output limits
- Structured receipt

## Phase 2: reference vertical slice

- User-issued regular-file capabilities
- Read-only input staging
- Immutable `bun-data@1` package set
- Dedicated scratch and output volumes
- Artifact validation and audience rules
- CLI and MCP adapters
- End-to-end data transformation demonstration

## Phase 3: authoritative isolation

- OCI backend
- gVisor execution
- Attack-corpus automation
- Cold and warm performance benchmarks
- Runtime image signing and verification
- Published security claims and limitations

## Phase 4: portability proof

- Node adapter
- Runtime conformance suite
- Runtime discovery and version policy
- Runtime-neutral fixtures

## Phase 5: broader task capabilities

- Custom immutable dependency preparation
- Deno adapter
- Directory and repository snapshots
- Approved network and API brokers
- Approved subprocesses and patch artifacts
- Temporary services
- Hosted scheduling and Firecracker backend

## Deferred decisions

- Public license
- Permanent Go module path
- Public repository and release model
- Supported desktop virtualization strategy
- Hosted offering and multi-tenancy
