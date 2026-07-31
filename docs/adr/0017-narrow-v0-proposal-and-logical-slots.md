# ADR-0017: V0 omits unsupported authority and uses logical resource slots

- Status: Accepted
- Date: 2026-07-31

## Context

The scaffold `Job` object models network, subprocess, environment, native, package, host capability,
and guest path shapes that v0 will not support. Representing future authority and always setting it
to false makes parsers, approvals, and policy unnecessarily complex and leaves accidental widening
paths.

## Decision

The target `JobProposal` includes only dependency-free source, a trusted profile selector, inline
JSON or opaque input slot, proven limit dimensions, bounded logical output slots, and bounded
non-authoritative labels.

V0 omits network, process, environment, native/FFI, package, arbitrary image/backend, path, special
filesystem, and agent-supplied general schema powers. Unknown/future power is unsupported protocol.

Agents use logical slot identifiers. Capsule assigns all guest paths internally and binds exact
slot/content mappings in the plan.

## Consequences

- The approval surface and Supervisor hard-safety validator are smaller.
- Agent strings cannot become host/guest path authority.
- Adding a future capability requires a new explicit contract, threat review, tests, and ADR where
  consequential.
- Current schemas and TypeScript types remain pre-freeze scaffold until replaced after feasibility
  evidence.
