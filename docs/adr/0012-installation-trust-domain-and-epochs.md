# ADR-0012: Installation Trust Domain, signed manifests, and trust epochs

- Status: Accepted
- Date: 2026-07-31

## Context

Purpose-separated keys do not by themselves detect mixed component versions, restored policy,
stale peers, partial updates, or changed entitlements. Capsule needs one enrolled view of the local
installation and explicit transition/recovery semantics.

## Decision

Each installation has a random normative `installationId`, hardware-backed installation root where
supported, operational-key authorizations, signed `InstallationManifest`, and a sequence-ordered
trust-epoch chain.

The manifest binds expected component code identities and relevant entitlements, operational keys,
policy/profile and external-trust checkpoints, storage formats, prior epoch digest, and transition
reason. Trusted IPC, plans, registrations, approvals, attempts, and receipts bind the active epoch.

Component-changing transitions use a prepared, explicitly authorized, crash-safe update ceremony.
Partial state enters `repair-required`; the daemon cannot clear it.

## Consequences

- Partial updates and stale/mismatched enrolled components fail closed.
- Installation-root use remains rare and unavailable to the daemon.
- Trust state and update/recovery become security-critical persisted data.
- Epochs are sequence-ordered but not inherently rollback-proof. Coherent rollback requires a
  non-rollbackable or independent checkpoint to detect reliably.
- Lost root/OS replacement creates a new installation identity and invalidates pending approvals.
