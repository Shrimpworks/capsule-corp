# Runtime Integrity

Status: intended point-in-time macOS design; signed mechanism and installed per-user lifecycle
feasibility is observed, while distribution, cross-session, and product-store validation remains
pending.

## Claim boundary

Capsule intends to claim only that, at recorded checkpoints, macOS reported participating
processes satisfied enrolled code identities and dynamic validity requirements, authenticated local
channels enforced expected peer requirements, components bound the same trust epoch, and no checked
disqualifying condition was present.

This does not prove the kernel, hypervisor, Secure Enclave, signing infrastructure, or correctly
signed program logic is uncompromised.

## Peer authentication

Trusted local IPC uses OS-enforced XPC peer code requirements where available. A trusted connection
checks:

- expected signing/team identity;
- expected distribution/certificate channel;
- component signing identifier;
- purpose-specific entitlement or code requirement where useful;
- expected effective user and session;
- exact enrolled build/epoch at the protocol layer;
- connection/start identity rather than PID/path/name alone.

Each component accepts only the peers and operations required by its authority. The Broker does not
expose an agent endpoint; the backend-control endpoint is Supervisor-only.

On the tested macOS 26 host, listener peer requirements and `SecCodeCreateWithXPCMessage` composed
in both directions with exact Team/channel/identifier/code-directory-hash predicates. Three
unprivileged per-user LaunchAgents then exercised four purpose-specific channels, effective UID and
audit-session agreement, stale client/service denial, and activation of a fresh Supervisor after
`SIGKILL`. This remains Apple Development spike evidence: distribution packaging, cross-user/fast-
user switching, logout/login, in-flight retry, and a real durable epoch store remain pending.

Keychain access-group membership must not be inferred from an exact peer result. Gate B observed a
stale exact-hash-denied build use its stable group key; the follow-up security-epoch group/key
transition is the separate key-use boundary.

## Dynamic validation

At connection and attempt preflight, record/check as supported:

- dynamic code validity;
- Hardened Runtime expectations;
- prohibition of debugging for validated posture;
- active code-directory hash;
- signing identifier and team ID;
- relevant entitlements;
- effective user/session and process-start identity;
- active trust epoch.

Code signing establishes reported identity/integrity state, not logic correctness.

## RuntimeIntegrityAssessment

The Supervisor records one bounded assessment per attempt. It includes:

- installation manifest and epoch identity;
- Supervisor self-check result;
- daemon and Broker peer-check results;
- Approval-key authorization and status;
- local trust-snapshot identity/freshness;
- runtime bundle, review, registry, and backend-validation identities;
- exact backend probe/capability result;
- grant-ledger and recovery health;
- active degraded/quarantine/repair flags;
- checkpoint time and evidence mode.

The assessment supports an internal short-lived permission to proceed from preflight to backend
creation. This internal permission is not exported as platform attestation.

## State machine

```text
uninitialized
  → verifying
  → verified
  → authoritative-ready
  → executing

failure/degradation:
  → degraded
  → quarantined
  → repair-required
  → compromised
```

- `degraded`: optional evidence or trust freshness does not meet requested posture.
- `quarantined`: component, debug, signature, backend, epoch, or other material mismatch.
- `repair-required`: interrupted update or inconsistent enrolled state.
- `compromised`: strong evidence of unauthorized modification; execution disabled.

The daemon cannot clear or downgrade these states.

## Evidence modes

Receipts use only implemented modes:

- `startup-only`
- `preflight-point-in-time`
- `periodically-revalidated`
- `continuously-monitored`
- `platform-attested` (future only)

Without an independent Runtime Guardian, v0 is at most point-in-time. Periodic renewal is meaningful
only when it rechecks independent signals. A timer that merely refreshes an in-memory flag is not
new integrity evidence.

## Failure during an attempt

- Before grant consumption: reject without consuming when state remains determinate.
- After consumption and before guest start: terminal `runtime-integrity-failed`; grant remains
  consumed.
- While running: stop new attempts, terminate/destroy the guest, quarantine outputs, and sign only a
  failure transcript.
- After collection but before destruction/release: quarantine artifacts and report failed or
  indeterminate integrity, never success-with-warning.
- When teardown state is unknown: record unresolved/teardown failure and keep cleanup responsibility
  durable.

## Optional Runtime Guardian

An Endpoint Security Guardian is a gated, non-blocking research item. It may observe component exec
identity, code-signature invalidation, trace/debug attempts, unexpected children, relevant signals,
and protected-state changes.

It requires entitlement, deployment, root/admin, user-approval, Full Disk Access, event-coverage,
performance, and false-positive analysis. It begins notify-only and never authorizes an attempt.

## Required tests

- unsigned and ad-hoc peer;
- same-team wrong signing identifier;
- expected binary with stale epoch;
- copied/replaced binary and mismatched code-directory hash;
- debugged or dynamically invalid component;
- wrong effective user/session;
- replayed PID/path/name assumptions;
- trust-snapshot downgrade/staleness;
- integrity failure before/after consumption, during execution, after collection, and during
  teardown;
- daemon attempt to clear quarantine or fabricate a passing assessment.

Results must identify the exact macOS/SDK/API availability and may not generalize beyond tested
platform ranges.
