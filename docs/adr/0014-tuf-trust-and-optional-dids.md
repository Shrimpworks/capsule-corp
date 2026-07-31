# ADR-0014: TUF anchors external trust; DIDs are optional identifiers

- Status: Accepted
- Date: 2026-07-31
- Supersedes: the normative per-installation `did:key` root decision in ADR-0006

## Context

DID Core provides interoperable identifiers and verification methods, but method resolution does
not grant Capsule authority. `did:key` is useful offline but cannot update or deactivate, making it
an incomplete long-lived recovery/revocation foundation. Release/profile distribution also needs
delegation, snapshot consistency, expiration, and rollback defenses.

## Decision

The normative local identity is an opaque installation ID plus locally authorized public keys.
DIDs remain first-class optional representations for organizations, reviewers, workers,
verification methods, and exported receipts. Local v0 approval/execution performs no network DID
resolution, arbitrary method loading, resolver plugin execution, or remote context retrieval.

Pinned TUF root metadata anchors release/profile distribution. A network-capable updater verifies
external metadata and produces a compact signed local `TrustSnapshot`; the Supervisor consumes that
snapshot without live network/TUF resolution.

TUF carries Capsule-defined revocation/disable objects. Capsule policy defines their semantics.

## Consequences

- DID interoperability remains available without making method behavior the authorization root.
- `did:key` can render an operational key but local status/replacement remains authoritative.
- URLs, TLS, DIDs, and service `valid` responses cannot replace pinned roots.
- Production update/trust infrastructure is substantial but stays outside the live execution TCB.
- Offline policy must explicitly define freshness, grace, downgrade, and refusal behavior.
