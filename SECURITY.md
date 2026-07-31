# Security Policy

## Current status

Capsule is pre-release research software. The repository does not yet provide a reviewed sandbox
implementation and must not be relied upon to contain hostile code.

The word "secure" in design documents refers to intended, testable properties—not a current
certification or guarantee.

## Intended security properties

The authoritative backend is expected to enforce:

- No ambient host filesystem access
- No inherited host credentials or environment
- No network unless explicitly granted
- No subprocess, native addon, or FFI authority unless explicitly granted
- Bounded CPU, memory, process count, storage, time, logs, and artifacts
- Fresh execution state and forced teardown
- Controlled, audience-aware egress
- An auditable receipt describing the effective policy and runtime identity

See the [threat model](docs/security/THREAT_MODEL.md) for scope and non-guarantees.

## Reporting a vulnerability

Do not publish suspected sandbox escapes or credential exposures in a public issue.

For a GitHub-hosted repository, use its private vulnerability reporting flow under **Security →
Advisories → Report a vulnerability**. Repository administrators must enable private vulnerability
reporting before a public preview. If that flow is unavailable, do not open a public issue or
disclose exploit details; contact a maintainer through an established private channel first.

Include, where possible:

- Affected commit or version
- Isolation backend and host platform
- Runtime profile and digest
- Minimal reproduction
- Expected and observed capability boundary
- Potential impact

## Supported versions

No released version is currently supported. Security support policy will be defined before the
first public release.
