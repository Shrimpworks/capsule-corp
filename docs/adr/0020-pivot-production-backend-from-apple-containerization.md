# ADR-0020: Pivot the production backend from Apple Containerization

- Status: Accepted
- Date: 2026-07-31
- Supersedes: ADR-0008

## Context

Capsule requires a durable backend identity that survives Supervisor death and supports
authoritative enumeration, inspection, forced destruction, and independent teardown evidence.
Missing process, handle, or response state cannot prove that a hostile guest is gone.

Gate C tested Apple Container 1.0.0 and direct Containerization 0.33.3. The stock service retained
inconsistent lifecycle state after its controller died and lacked required exact controls. Direct
Containerization improved controller-coupled cleanup and exposed useful VM, network, filesystem,
memory, and CPU mechanisms. A narrow experimental patch also enforced `pids.max`.

The lower-level public API nevertheless exposes no durable host-side VM/helper identity and no
supported enumeration, reopen, or reconnect operation after the owning controller exits.
`VZGenericMachineIdentifier` is guest virtual-hardware identity, not a host lifecycle handle.
Concurrent unrelated Virtualization helpers also demonstrated that PID-delta/process-name
attribution cannot safely substitute for a supported identity. Management vsock reachability
remains unresolved.

An OCI control harness running under the already-installed runc runtime demonstrated stable
engine-issued container IDs and labels across client death plus the required configuration shape
for network, privilege, cgroup, storage, output, and cancellation controls. Actual gVisor/runsc and
engine/outer-VM crash behavior remain unvalidated.

## Decision

- Direct Apple Containerization remains an optional development-only macOS backend. It may support
  local integration work but cannot claim authoritative teardown or `validated-local` posture.
- Any ambiguous Apple controller loss remains `unresolved` with execution and ordinary artifact
  release disabled. Capsule will not add private APIs, helper-PID guessing, process-name scans,
  root-only introspection, or a privileged helper to manufacture the missing lifecycle authority.
- OCI plus a pinned gVisor runtime becomes the primary production-backend candidate. It must run on
  a governed disposable Linux worker or managed Linux VM and pass the complete shared attack
  corpus, including engine/containerd and outer-VM failure.
- The fake backend remains the only backend used for contract and lifecycle implementation before
  hostile execution.
- Apple can re-enter production evaluation only if a future supported public API provides durable
  identity and enumerate/reconnect/force-reap semantics, or a verifiable documented destruction
  guarantee satisfies the same evidence contract.

This decision selects the next candidate, not a validated backend. The current runc result is
control-harness evidence only and makes no gVisor isolation claim.

## Consequences

- The portable backend protocol and task contract do not change.
- macOS development can remain native for control-plane, Broker, XPC, Keychain, and fake-backend
  work; authoritative gVisor validation requires Linux infrastructure.
- A Mac-hosted gVisor path adds an outer Linux VM, Docker/containerd or another governed OCI
  engine, a high-authority engine endpoint, image/runtime update policy, and more lifecycle state.
- The Supervisor adapter must own an exact runtime/image allowlist, stable attempt labels and IDs,
  bounded output independent of engine logs, and durable effect intents before engine calls.
- Apple-specific resource and lifecycle experiments remain retained evidence, not production code.

## Evidence

- [`experiments/apple-container-gate-c/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/apple-container-gate-c/RESULTS.md)
- [`experiments/apple-containerization-direct/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/apple-containerization-direct/RESULTS.md)
- [`experiments/gate-c-identity-recovery/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-identity-recovery/RESULTS.md)
- [`experiments/gate-c-gvisor-contingency/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-gvisor-contingency/RESULTS.md)
