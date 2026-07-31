# Repository Scope And Priorities

This repository builds Capsule, a local-first, capability-controlled JavaScript and TypeScript task
runtime for AI agents.

Primary deliverables:

- a trusted Go control plane for validation, policy, lifecycle, egress, and receipts
- runtime-neutral JSON contracts with TypeScript SDK, CLI, and MCP adapters
- disposable Bun execution through externally enforced isolation backends

Core priorities:

- deny-by-default authority and fail-closed behavior
- externally enforced isolation with explicit backend security tiers
- stable typed contracts between modules
- controlled, audience-aware output and reproducible evidence

## Active Boundaries

- `schemas/` owns the canonical wire contracts.
- `internal/` and `cmd/capsuled/` own the trusted Go control plane.
- `packages/` owns client-facing protocol views and adapters.
- `profiles/` owns reviewed runtime-profile declarations, not mutable guest configuration.

## Safe Refactor Boundaries

Do not refactor these without explicit instruction:

- capability issuance, redemption, scope, expiry, or audience semantics
- isolation-backend configuration or security-tier claims
- artifact validation, output exposure, or receipt semantics
- canonical schemas, public API routes, and runtime-profile identity

Safe default changes:

- feature-scoped improvements
- fail-closed endpoint hardening and validation
- focused test additions
- typing improvements that remain consistent with canonical schemas
