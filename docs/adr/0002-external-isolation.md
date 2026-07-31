# ADR-0002: External isolation is mandatory

- Status: Accepted
- Date: 2026-07-30

## Context

Node permissions, Deno permissions, Bun settings, SES, and LavaMoat reduce authority but do not form
a sufficient hostile-code boundary for every supported runtime and workload.

## Decision

All untrusted jobs require an externally enforced isolation backend. Runtime restrictions and
JavaScript confinement are supplemental layers.

Strong security claims will be tied to a named backend tier and adversarial test results. The first
serious target is OCI plus gVisor on Linux. Lightweight local backends may be used for development
but must be labeled non-authoritative.

## Consequences

- The core guarantee does not depend on SES or runtime compatibility.
- Local cross-platform execution requires additional virtualization or platform adapters.
- Startup overhead must be measured at the complete sandbox level.
- Receipts must record the isolation backend and security tier.
