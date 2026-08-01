# ADR-0022: Evaluate libkrun/HVF as the native Apple backend candidate

- Status: Accepted
- Date: 2026-07-31
- Readiness synthesis: 2026-07-31
- P0 reconciliation: 2026-07-31
- Refines: ADR-0020

## Context

ADR-0020 rejected stock Apple Container and direct Apple Containerization as production backends
because their public lifecycle surfaces did not provide a durable host-side VM/helper identity,
restart enumeration, reopen, or exact force-reap authority. It selected OCI plus gVisor as the next
candidate, not as an already validated backend.

An additional Gate C spike tested libkrun 1.19.4 directly on Apple Hypervisor.framework. Unlike
Containerization, libkrun runs one VM inside the calling VMM process. The signed process therefore
provides a concrete lifecycle object that macOS can identify using PID, process start time,
Security.framework code identity/CDHash, and the expected executable location.

The spike compiled out virtio-net, disabled implicit vsock/TSI, used only a raw read-only block
root, ran a trusted guest launcher as UID/GID 65534 with no capabilities and `no_new_privs`, and
enforced rlimits. A properly bundled Developer ID runner also ran under App Sandbox with only an
exact read-only disk exception; the no-disk-authority control failed. Concurrent cancellation and
three controller-`SIGKILL` checkpoints passed. Digest-pinned Bun 1.3.14 ran native TypeScript, and
five sequential minimal guests completed in 0.13–0.15 seconds.

Two narrow source changes were required: signed-bundle `@rpath` firmware resolution and correct
Linux read-only block-root mount flags. Five follow-up tracks then tested block storage/output,
bounded console and resources, installed recovery, adversarial VMM behavior, and runtime supply
chain. They strengthened the lifecycle evidence but found that the exact profile is not ready to
freeze: libkrun's pathname disk API did not prevent live same-user backing-file mutation, and its
block-root API added a guest-visible `NullFs` virtiofs device even though no host directory was
configured or mounted. Runner exit zero was also ambiguous across corrupt roots, guest panics, and
missing executables. The current runtime bytes failed the development admission checklist.

A subsequent independent adversarial review and targeted source research refined the remaining
work without promoting the candidate. libkrun 1.19.4 exposes generic multiport console APIs, so the
first inline slice can test port-based source/input and typed inline results rather than requiring
user-data disks. A genuine inherited read-only descriptor exposed as `/dev/fd/N` is the first
runtime-root custody candidate. The same review found that stock Bun's subprocess and FFI APIs do
not satisfy the documented prohibited-power contract merely because the binary is pinned. All
three conclusions remain P0 hypotheses until their exact installed corpora pass.

## Decision

- libkrun/HVF becomes the preferred **native Apple backend candidate** for the next
  implementation and validation slice. This is a candidate selection, not a posture promotion.
- The five readiness tracks do not admit a real libkrun adapter for user bytes yet. Before freezing
  one development profile, Capsule must close the exact Bun runtime-authority restrictions,
  immutable runtime-root custody, independent acceptance or removal of the `NullFs` device, typed
  attempt-bound port transport/completion with bounded inline JSON, and an admissible complete
  installed bundle. A bounded disposable filesystem-image parser is mandatory before later file
  artifacts, not before the inline JSON slice.
- OCI plus gVisor remains an independent candidate and contingency. Final selection compares the
  completed exact profiles rather than assuming either backend is secure by reputation.
- Apple Containerization remains development-only; this ADR does not reverse ADR-0020's conclusion
  about its lifecycle authority.
- The native runner uses one signed App-Sandboxed VMM process and one VM per attempt. The Supervisor
  durably records and verifies the exact live process identity before authorizing VM start through
  a fail-closed inherited control pipe.
- V0 compiles only required libkrun device features. It exposes no network, implicit vsock/TSI,
  host-backed virtiofs directory, host socket, arbitrary image format, GPU, or sound device. The
  current block-root path's `NullFs` device is unresolved and must not be described as absent.
- The first custody campaign tests protected creation and a genuine read-only descriptor inherited
  directly into the runner and exposed to libkrun as `/dev/fd/N`; a narrow FD-native API change is
  the fallback. Guest read-only flags and post-stop hashing do not make a same-user mutable
  pathname immutable. First-slice source/input bytes use bounded dedicated console ports; later
  file storage must pass its own immutable custody and parsing gates. Writable scratch uses a fresh
  bounded disposable raw block device.
- The trusted guest launcher is part of the reviewed runtime bundle, not an invocation of `su` or a
  host UID setter. It retains the completion descriptor while dropping the workload without that
  authority. Ordinary success requires a typed attempt-bound completion/result frame plus input,
  result-validation, integrity, and teardown evidence; runner exit status alone is never success.
- Pinning Bun proves byte identity, not absence of runtime powers. Until the exact profile refuses
  subprocess, FFI, native-addon, inspector, macro, environment-file, and package-install paths,
  execution requiring that contract is unsupported. A contract change requires a separate ADR.
- Host capture retains at most the exact approved prefix per stream while continuously draining.
  Wall/cancel actions remain independent of guest cooperation, and revalidated exact-process
  `SIGKILL` is the required fallback because graceful eventfd shutdown did not pass.
- Resource admission uses only closed evidenced vCPU/RAM profiles. For the current Bun fixture the
  smallest observed profile is one vCPU/256 MiB. CPU percentage/time, arbitrary RAM, and exact total
  host/VMM memory are unsupported.
- Apple silicon and macOS 14+ are a provisional source/platform target only. The complete package
  has not passed a macOS 14 clean-host floor, and the current runner/firmware metadata declares
  macOS 26. Intel support requires a separately selected and validated backend.
- Capsule must upstream the two patches or maintain a governed fork with exact source publication,
  SBOM/provenance, advisories, and LGPL/GPL source-compliance handling for libkrunfw and its kernel.

No backend may claim `validated-local` until the exact distributed runner, libraries, firmware,
runtime disks, entitlements, host version, limits, and recovery configuration pass the full corpus
and are bound by an accepted `BackendValidationRecord`.

## Consequences

- The immediate product path can remain native on Apple silicon without adding a Linux VM, OCI
  engine endpoint, or root helper solely to obtain the first real backend.
- The backend adapter must preserve the record-before-start handshake and compare PID/start/code
  identity on recovery. PID, path, name, or process scans alone remain non-authoritative.
- App Sandbox materially narrows VMM-compromise impact, but the VMM and Hypervisor boundary still
  require malicious-guest testing and security-update ownership.
- Runtime-root custody, bounded port framing, and console capture become explicit TCB components
  with their own ownership, quota, and crash-recovery designs. Block-image extraction joins that
  TCB only when file artifacts are introduced and must run in a disposable bounded parser sandbox.
- Installed recovery depends on an exact enrolled launch profile. The experiment's
  `AbandonProcessGroup=true` allowed exact reparent/reap recovery; changing it changed orphan
  behavior. The complete Supervisor/helper/LaunchAgent package, notarization, Gatekeeper, session,
  reboot, and clean-host matrix remains required.
- The current libkrun/libkrunfw/runner bytes are no-go for a development validation record. A
  controlled path-remapped build demonstrated reproducibility feasibility, not independent
  two-builder provenance or release admission.
- gVisor work is reprioritized from the only production path to an independent comparison and
  fallback; its existing runc control result still makes no gVisor isolation claim.

## Evidence

- [`experiments/gate-c-libkrun-hvf/RESULTS.md`](../../experiments/gate-c-libkrun-hvf/RESULTS.md)
- [`experiments/gate-c-libkrun-hvf/README.md`](../../experiments/gate-c-libkrun-hvf/README.md)
- [`experiments/gate-c-identity-recovery/RESULTS.md`](../../experiments/gate-c-identity-recovery/RESULTS.md)
- [`experiments/gate-c-gvisor-contingency/RESULTS.md`](../../experiments/gate-c-gvisor-contingency/RESULTS.md)
- [`experiments/gate-c-libkrun-storage-egress/RESULTS.md`](../../experiments/gate-c-libkrun-storage-egress/RESULTS.md)
- [`experiments/gate-c-libkrun-console-lifecycle/RESULTS.md`](../../experiments/gate-c-libkrun-console-lifecycle/RESULTS.md)
- [`experiments/gate-c-libkrun-installed-recovery/RESULTS.md`](../../experiments/gate-c-libkrun-installed-recovery/RESULTS.md)
- [`experiments/gate-c-libkrun-adversarial/RESULTS.md`](../../experiments/gate-c-libkrun-adversarial/RESULTS.md)
- [`experiments/gate-c-libkrun-supply-chain/RESULTS.md`](../../experiments/gate-c-libkrun-supply-chain/RESULTS.md)
- [Gate C implementation-readiness synthesis](../GATE_C_READINESS_CHECKPOINT.md)
- [Gate C P0 reconciliation](../GATE_C_P0_RECONCILIATION.md)
