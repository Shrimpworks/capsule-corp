# Architecture Decision Records

Architecture Decision Records capture decisions that materially affect security, compatibility, or
the long-term structure of Capsule.

Status values:

- `Proposed`
- `Accepted`
- `Superseded by ADR-NNNN`
- `Rejected`

Accepted ADRs are historical records. Replace an outdated decision with a new ADR instead of
silently editing its conclusion.

## Index

- [ADR-0001: Jobs, not computers](0001-jobs-not-computers.md)
- [ADR-0002: External isolation is mandatory](0002-external-isolation.md)
- [ADR-0003: Bun-first, runtime-neutral protocol](0003-bun-first.md)
- [ADR-0004: Guest output is controlled egress](0004-controlled-egress.md)
- [ADR-0005: Go for the initial trusted control plane](0005-go-control-plane.md)
- [ADR-0006: Signed execution plans and per-device identity](0006-signed-plans-and-device-identity.md)
- [ADR-0007: Regular-file capabilities reference immutable snapshots](0007-snapshot-file-capabilities.md)
- [ADR-0008: Apple Container and gVisor are independent backend targets](0008-apple-container-and-gvisor-backends.md)
- [ADR-0009: Resource policy is user-owned and exact](0009-user-owned-exact-resource-limits.md)

Use [the ADR template](TEMPLATE.md) for new decisions.
