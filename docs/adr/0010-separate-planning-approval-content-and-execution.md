# ADR-0010: Separate planning, trusted-host approval/content, and execution authority

- Status: Accepted
- Date: 2026-07-31
- Supersedes: the single-daemon authority portions of ADR-0006

## Context

The original design placed planning, approval verification, content brokerage, backend launch,
grant state, and receipt signing around one daemon. A compromise of the agent-facing parser/service
could therefore approach every authority needed to run hostile code or disclose user content.

Separate processes improve containment only when operating-system identity, key, storage, and IPC
controls make them separate authorities.

## Decision

Capsule uses three local authorities:

- the Go daemon authenticates agents and constructs proposed immutable plans;
- a native Trusted Host Broker owns user-presence approval and user-content custody;
- an Execution Supervisor independently validates registered plans and alone creates hostile guests.

The daemon has no Approval/Supervisor private key, user-only content access, backend launch path,
grant-reset authority, or quarantine-reset authority. The Broker has no agent-facing endpoint or
backend launch authority. The Supervisor has no public agent parser, file-picker UX, rich parser, or
network trust resolution.

Approval and Content Broker interfaces may share one native v0 process but remain logically,
cryptographically, and persistently separable.

## Consequences

- Daemon compromise can propose malicious work and cause denial of service but cannot by itself
  approve, launch, disclose Broker content, or forge Supervisor enforcement evidence.
- Local IPC, Keychain access groups, protected storage, exact code identity, and entitlements become
  mandatory feasibility evidence; a diagram or Unix mode bits are insufficient.
- More components and cross-store recovery add deployment and state-machine complexity.
- A later high-assurance deployment can split Approval and Content Brokers without changing the
  public job abstraction.
