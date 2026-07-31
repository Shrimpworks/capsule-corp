# ADR-0018: Platform-specific trusted components use least privilege

- Status: Accepted
- Date: 2026-07-31
- Refines: ADR-0005 and ADR-0008

## Context

Go remains a strong fit for the agent-facing orchestration daemon, but the Broker requires native
macOS UI, LocalAuthentication, Keychain, XPC, code-signing, and file-selection integration. The
Supervisor also needs deep platform/backend integration, but its safest language and privilege
model cannot be chosen before the capability spikes.

## Decision

The agent-facing daemon remains Go. The macOS Trusted Host Broker is preferably Swift/native.

Supervisor language remains deferred while a spike compares native Swift, Go with narrow native
bindings, and a hybrid unprivileged Supervisor with a tiny launcher helper if required. The default
deployment target is unprivileged per-user service/agent. Capsule does not assume a root
LaunchDaemon.

If limited privilege is required, isolate it in the smallest helper that accepts only a sealed typed
launch descriptor and parses no agent proposal, policy, approval UI, content format, or trust
network data.

Another implementation language, including Rust, requires a separate ADR showing a narrow
interface and demonstrated security/operational benefit.

## Consequences

- Native macOS authority controls can be used directly where they matter.
- The repository may become multi-language, so shared cryptographic/protocol fixtures are mandatory.
- Supervisor implementation is not prematurely locked before Apple Container, XPC, Keychain,
  storage, and privilege evidence exists.
- Language memory safety supplements but does not replace the external guest boundary or authority
  design.
- Privileged surface stays smaller than the Supervisor whenever a helper is unavoidable.
