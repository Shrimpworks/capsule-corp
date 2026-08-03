# Architecture Decision Records

Architecture Decision Records capture decisions that materially affect security, compatibility, or
the long-term structure of Capsule.

Status values:

- `Proposed`
- `Accepted`
- `Superseded by ADR-NNNN`
- `Rejected`

Accepted ADRs are historical records. Replace an outdated decision with a new ADR instead of
silently rewriting its conclusion. A later ADR may refine an accepted decision without superseding
its core conclusion.

## Index

- [ADR-0001: Jobs, not computers](0001-jobs-not-computers.md)
- [ADR-0002: External isolation is mandatory](0002-external-isolation.md)
- [ADR-0003: Bun-first, runtime-neutral protocol](0003-bun-first.md) — first-runtime ordering
  superseded by ADR-0028; runtime-neutral protocol retained
- [ADR-0004: Guest output is controlled egress](0004-controlled-egress.md)
- [ADR-0005: Go for the initial trusted control plane](0005-go-control-plane.md) — refined by
  ADR-0018
- [ADR-0006: Signed execution plans and per-device identity](0006-signed-plans-and-device-identity.md)
  — superseded by ADR-0010, ADR-0011, and ADR-0014
- [ADR-0007: Regular-file capabilities reference immutable snapshots](0007-snapshot-file-capabilities.md)
- [ADR-0008: Apple Container and gVisor are independent backend targets](0008-apple-container-and-gvisor-backends.md)
  — superseded by ADR-0020
- [ADR-0009: Resource policy is user-owned and exact](0009-user-owned-exact-resource-limits.md)
  — refined by ADR-0011
- [ADR-0010: Separate planning, trusted-host approval/content, and execution authority](0010-separate-planning-approval-content-and-execution.md)
- [ADR-0011: Execute only Supervisor-registered immutable plans](0011-supervisor-registered-plans-and-one-use-attempts.md)
- [ADR-0012: Installation Trust Domain, signed manifests, and trust epochs](0012-installation-trust-domain-and-epochs.md)
- [ADR-0013: Runtime integrity uses authenticated IPC and point-in-time assessments](0013-point-in-time-runtime-integrity.md)
- [ADR-0014: TUF anchors external trust; DIDs are optional identifiers](0014-tuf-trust-and-optional-dids.md)
- [ADR-0015: Supervisor transcripts and composed execution receipts](0015-supervisor-transcripts-and-composed-receipts.md)
- [ADR-0016: Runtime bundle, review, activation, and backend validation are separate](0016-separate-runtime-profile-evidence.md)
- [ADR-0017: V0 omits unsupported authority and uses logical resource slots](0017-narrow-v0-proposal-and-logical-slots.md)
- [ADR-0018: Platform-specific trusted components use least privilege](0018-platform-specific-trusted-components.md)
- [ADR-0019: Use bounded deterministic CBOR and object-specific COSE profiles](0019-bounded-deterministic-cbor-and-cose.md)
- [ADR-0020: Pivot the production backend from Apple Containerization](0020-pivot-production-backend-from-apple-containerization.md)
  — refined by ADR-0022
- [ADR-0021: Scope operational Keychain groups to security epochs](0021-security-epoch-keychain-groups.md)
- [ADR-0022: Evaluate libkrun/HVF as the native Apple backend candidate](0022-evaluate-libkrun-hvf-native-backend.md)
- [ADR-0023: Bound protocol decoding and registration semantics](0023-bound-protocol-decoding-and-registration.md)
  — Proposed
- [ADR-0024: Bound approval consumption and attempt creation before effects](0024-approval-consumption-and-attempt-creation.md)
  — Proposed
- [ADR-0025: Colocate durable attempt lifecycle state with Supervisor authority state](0025-colocate-durable-attempt-lifecycle-state.md)
  — Proposed
- [ADR-0026: Bind strip-only TypeScript emission before plan registration](0026-bind-pre-approval-typescript-erasure.md)
  — Proposed
- [ADR-0027: Retire the SupervisorCore in-memory scaffold](0027-retire-supervisorcore-scaffold.md)
- [ADR-0028: Select governed deno_core as the first runtime candidate](0028-select-governed-deno-core-first.md)
- [ADR-0029: Select one native-fronted Go Supervisor process for authenticated local IPC](0029-select-authenticated-local-ipc-topology.md)
  — Proposed
- [ADR-0030: Define the passive TypeScript approved-byte migration boundary](0030-define-typescript-approved-byte-migration-boundary.md)
- [ADR-0031: Checkpoint closed Supervisor cohorts into immutable retained archives](0031-checkpoint-closed-supervisor-cohorts.md)
  — Proposed

Use [the ADR template](TEMPLATE.md) for new decisions.
