# ADR-0013: Runtime integrity uses authenticated IPC and point-in-time assessments

- Status: Accepted
- Date: 2026-07-31

## Context

Code signatures at installation time do not authenticate each live IPC peer, detect partial update,
or prove continuous runtime integrity. PID, path, and process name are weak identities. Endpoint
Security monitoring is deployment-sensitive and cannot be assumed in v0.

## Decision

Trusted macOS IPC uses OS-enforced XPC peer code requirements where available, plus protocol
binding to installation and epoch. Connection/preflight checks record expected signing identifier,
team, effective user/session, exact active build/code-directory identity, relevant entitlements,
dynamic validity, and debug state where supported.

The Supervisor creates one bounded `RuntimeIntegrityAssessment` before grant consumption and
backend creation. Without an independent Guardian, evidence is labeled point-in-time. A short
internal permission to proceed is not exported as attestation.

Integrity state can degrade, quarantine, require repair, or become compromised; the daemon cannot
clear it. A post-consumption integrity failure burns approval and prevents ordinary artifact
release/success.

## Consequences

- Local peer authentication and exact build/epoch checks are part of the security boundary.
- Valid signatures remain identity/integrity observations, not proof of correct logic.
- Continuous-monitoring or platform-attestation labels are prohibited until an implemented
  mechanism supports them.
- Endpoint Security remains an optional notify-only research component and never authorizes jobs.
