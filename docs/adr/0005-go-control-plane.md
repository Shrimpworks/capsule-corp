# ADR-0005: Go for the initial trusted control plane

- Status: Accepted
- Date: 2026-07-30
- Refined by: ADR-0018

## Context

Capsule needs a trusted daemon for protocol validation, policy evaluation, capability coordination,
job lifecycle management, isolation-backend orchestration, controlled egress, and receipt
generation.

Both Go and Rust are viable memory-safe implementation languages. The initial control plane is
primarily an orchestration service rather than a custom kernel-facing sandbox. It benefits from a
small dependency graph, a strong standard library, direct process and HTTP support, simple
deployment, and mature OCI and gVisor integration.

The implementation language does not replace the externally enforced isolation boundary required
by ADR-0002.

## Decision

The initial trusted daemon and isolation orchestration layer are implemented in Go. TypeScript
remains the language for client adapters and the first guest workloads.

Security-critical helpers may be implemented in another memory-safe language, including Rust, only
when they have a narrow interface, run across an explicit process or privilege boundary where
appropriate, and provide a documented security or operational benefit. A consequential addition of
another implementation language requires an ADR.

## Consequences

- The initial control plane can use Go's standard HTTP, context, process, and structured-logging
  facilities with few dependencies.
- The repository avoids a mixed-language trusted computing base during the first vertical slice.
- Low-level privileged helpers are not required to share the daemon's implementation language.
- Security claims continue to depend on effective policy, external isolation, teardown, controlled
  egress, and adversarial tests—not on Go alone.
- The decision can be revisited if implementation evidence shows a specific component needs
  stronger low-level control or a smaller privileged footprint.
