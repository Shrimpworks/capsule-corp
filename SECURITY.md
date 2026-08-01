# Security Policy

## Current status

Capsule is pre-release research software. It does not yet provide a reviewed sandbox,
authorization boundary, runtime-integrity system, or authoritative receipt and must not be relied
upon to contain hostile code.

Security statements in design documents are intended falsifiable properties—not current
certification, implementation, or guarantees. Current schemas/types are pre-freeze scaffolding.

## Intended authorization properties

- Agent-facing planning is separate from user approval/content custody and backend launch.
- Only the Execution Supervisor may create a hostile guest.
- Execute uses a Supervisor-issued plan registration ID and never replacement plan bytes.
- The Broker renders exact registered typed data and requires fresh user presence per v0 plan.
- Approval is purpose-, audience-, installation-, epoch-, registration-, Supervisor-, attempt-, and
  expiry-bound and consumed atomically once.
- The Supervisor independently rejects unsupported v0 powers and backend capability mismatches.
- The daemon cannot use Approval/Supervisor keys, read user-only content, launch a backend, reset
  grants, or clear quarantine/trust state.

## Intended runtime-integrity properties

- Trusted local IPC authenticates exact enrolled component code identity, effective user/session,
  relevant entitlements, and common trust epoch using supported OS mechanisms.
- Partial update, stale build, debug/dynamic validity failure, or inconsistent enrolled state fails
  closed or enters repair/quarantine.
- Point-in-time assessments are not described as continuous monitoring or platform attestation.
- Trust epochs are sequence-ordered, not called rollback-proof without an independent anchor.

## Intended isolation and content properties

- No ambient host filesystem, credentials, environment, sockets, process, network, native, package,
  or backend authority.
- Inputs are immutable content identities and scoped handles rather than live host paths.
- Exact user-owned resource values are enforced or the attempt is refused.
- Every post-create path terminates/destroys/reconciles the guest; unknown teardown is not success.
- Bounded typed validation gates inline results; filesystem safety and disposable bounded parsing
  additionally gate later file artifacts.
- User content stays Broker-owned; the agent receives a fixed minimized summary by default.
- Rich parsing stays out of the daemon and Supervisor.

## Intended evidence properties

A user receipt at any validated posture requires both:

- Broker-attributable plan approval; and
- Supervisor-attributable enforcement/teardown transcript.

The receipt binds the same plan registration, attempt, installation, and trust epoch. These are
signed claims, not proof that the host, user interface, signer logic, kernel, or guest computation
was correct.

See the [threat model](docs/security/THREAT_MODEL.md),
[component compromise matrix](docs/security/COMPONENT_COMPROMISE_MATRIX.md), and
[control evidence matrix](docs/security/CONTROL_EVIDENCE_MATRIX.md).

## Reporting a vulnerability

Do not publish suspected sandbox escapes, authorization bypasses, key exposures, content-boundary
bypasses, runtime-integrity failures, or false authoritative evidence in a public issue.

Use the repository's private vulnerability reporting flow under **Security → Advisories → Report a
vulnerability**. If unavailable, contact a maintainer through an established private channel before
sharing exploit details.

Include where possible:

- affected commit/version and whether code is product, spike, test, or documentation;
- component identity and role: daemon, Broker, Supervisor, updater, Guardian, backend, or guest;
- host OS/hardware, code-signing identity, entitlements, and active trust epoch;
- isolation backend, exact configuration, runtime bundle/profile, and validation record;
- plan registration, attempt, trust-snapshot, and evidence identifiers without private content;
- minimal reproduction and attacker prerequisites;
- expected/observed authorization, runtime-integrity, isolation, content, or evidence boundary;
- cleanup/teardown and potential impact.

## Supported versions

No released version is currently supported. Security support policy will be defined before the
first public release.
