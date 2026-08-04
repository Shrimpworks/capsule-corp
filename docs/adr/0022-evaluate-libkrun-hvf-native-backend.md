# ADR-0022: Evaluate libkrun/HVF as the native Apple backend candidate

- Status: Accepted
- Date: 2026-07-31
- Readiness synthesis: 2026-07-31
- P0 reconciliation: 2026-08-01
- P0-0 stock-runtime result: 2026-08-02
- P0-1 FD-native patch candidate: 2026-08-02
- P0-3 cross-language and console coverage result: 2026-08-03
- Governed fork source reconciliation: 2026-08-03
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

The five retained libkrun patches now have a public governed source record. Their unchanged
aggregate identity is `d19fd0ff159c699acccda2621519de45a09408bf3847b418ac34e02b79e805d5`,
reconstructed from upstream `728df8125077d0db44265f6e997c72b81b65c015` to governed base
`4ea8d1de861ed1c0636fc800b6da8fb71a086aa5`. A bounded console/raw-FD follow-up merged from exact
head `8a2c91943793668f31a1cf7af431933be935bb58` as
`cf0333cdba478cc34a8570a65b38412da7fd3ecc`. It fixed two locally observed console
shutdown/lifecycle defects and raised measured coverage, but ran no guest and received no recorded
independent human or CODEOWNER review. The post-merge baseline branch also no longer satisfies the
verifier's hardcoded `4ea8d1de861ed1c0636fc800b6da8fb71a086aa5` branch pin, so current
source-governance consistency fails closed even though the five patch bytes and aggregate remain
exact. This changes source identity and library evidence only.

## Decision

- libkrun/HVF becomes the **lead native Apple candidate under evaluation** for the next
  implementation and validation slice. This prioritizes evidence collection; it is not backend
  admission, final selection, or a posture promotion.
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
  the fallback. It must separately pass stable attachment identity, frozen-object construction, and
  adversarial end-to-end custody. Final digest/length are computed through the exact retained
  descriptor whose `F_GETFL` access mode is exactly `O_RDONLY`, only after every writable alias/
  mapping is closed and its sole pathname is unlinked. Guest read-only flags and post-stop hashing
  do not make a same-user mutable pathname immutable. First-slice source/input bytes use bounded
  dedicated console ports; later file storage must pass its own immutable custody and parsing gates.
  Writable scratch uses a fresh bounded disposable raw block device.
- The governed P0-1 fallback is a patch candidate on pinned libkrun 1.19.4: its additive fixed-role
  API accepts only an exact finalized raw read-only descriptor, takes duplicate ownership, and
  routes imago directly from `File` without a pathname or format-selection input. Local custody,
  mutation, and fixed owned-guest digest evidence passed. This does not close P0-1C because the
  exact signed/notarized installed App Sandbox and protected-construction corpus remains untested.
- The trusted guest launcher is part of the reviewed runtime bundle, not an invocation of `su` or a
  host UID setter. It remains distinct instead of replacing itself with Bun, retains the completion
  descriptor while dropping and spawning the workload without that authority, verifies complete
  source/input before child start, uses a fixed child FD/argv/environment manifest, waits for exact
  child-tree termination, and commits completion last. Ordinary success requires that typed attempt-
  bound frame plus input, result-validation, integrity, and teardown evidence; runner exit status
  alone is never success.
- Port admission includes the guest-facing virtio control/queue/descriptor implementation and a
  closed host-runner FD manifest, not only application framing. The governed follow-up bounds the
  previously unchecked port-ID and stop-aware transmit cases, prevents zero-progress re-pop during
  shutdown, and makes shutdown return a port to `Inactive`. Its four-file coverage is 37/88
  functions and 298/733 lines, with 2 functions/26 lines in `port.rs` and 14 lines in
  `process_tx.rs` still uncovered. Undocumented negative-FD directionality, shared `O_NONBLOCK`
  status mutation, remaining hostile control/queue/descriptor coverage, and real composition must
  still pass an exact fail-closed corpus. The host continuously drains to cap-plus-one and never
  uses EOF as completion.
- Pinning Bun proves byte identity, not absence of runtime powers. The exact stock Bun 1.3.14 P0-0
  investigation and the follow-up governed-construction review both failed. The latter found a
  40-hand-authored plus 10-generated-output minimum and triggered its broad/unreviewable stop rule.
  Execution requiring the contract remains unsupported; alternate-runtime selection and an
  ADR-0003 superseding decision are required. A contract change requires a separate ADR.
- Host capture retains at most the exact approved prefix per stream while continuously draining.
  Wall/cancel actions remain independent of guest cooperation, and revalidated exact-process
  `SIGKILL` is the required fallback because graceful eventfd shutdown did not pass.
- Resource admission uses only closed evidenced vCPU/RAM profiles. For the current Bun fixture the
  smallest observed profile is one vCPU/256 MiB. CPU percentage/time, arbitrary RAM, and exact total
  host/VMM memory are unsupported.
- Apple silicon and macOS 14+ are a provisional source/platform target only. The complete package
  has not passed a macOS 14 clean-host floor, and the current runner/firmware metadata declares
  macOS 26. macOS 14 code-signature protection for sandbox app data containers must be tested
  directly; Capsule does not substitute macOS 15 app-group-container behavior or a broad shared app
  group for that evidence. Intel support requires a separately selected and validated backend.
- Host-root execution, a separate-owner host service, and a privileged host helper are prohibited
  for the v0 milestone. If no-host-root custody fails, libkrun fails that profile. Any later
  exception requires a new ADR covering authorization, installation, update, recovery, compromise
  radius, supported macOS floor, and comparison with memory-backed storage or another backend.
- Capsule now maintains a public governed five-patch source line and exact follow-up merge, but this
  is not a release. The branch/ref/verifier invariant must first be repaired without rewriting the
  retained patch bytes. Independent human/CODEOWNER review, a complete mapping of the follow-up
  deltas, exact source publication, SBOM/provenance, advisories, update/removal operations, and
  LGPL/GPL source-compliance handling for libkrunfw and its kernel remain required unless
  equivalent changes are accepted upstream.
- An early signed installed harness may test App Sandbox, service, descriptor, identity, and minimum-
  OS assumptions, but it cannot admit the backend. After P0 selects mechanisms and patches, Capsule
  rebuilds the complete final app and reruns every affected gate on those signed/notarized bytes;
  affected byte, entitlement, topology, or OS-floor changes invalidate earlier evidence.

No backend may claim `validated-local` until the exact distributed runner, libraries, firmware,
runtime disks, entitlements, host version, limits, and recovery configuration pass the full corpus
and are bound by an accepted `BackendValidationRecord` with that verdict. P0 may produce only a
`development-admitted` verdict with an explicit posture ceiling, limitations, expiry, and
invalidation triggers; it cannot be relabeled for `validated-local`.

## Consequences

- The immediate product path can remain native on Apple silicon without adding a Linux VM, OCI
  engine endpoint, or host-root helper solely to obtain the first real backend.
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

- [`experiments/gate-c-libkrun-hvf/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-hvf/RESULTS.md)
- [`experiments/gate-c-libkrun-hvf/README.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-hvf/README.md)
- [`experiments/gate-c-identity-recovery/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-identity-recovery/RESULTS.md)
- [`experiments/gate-c-gvisor-contingency/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-gvisor-contingency/RESULTS.md)
- [`experiments/gate-c-libkrun-storage-egress/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-storage-egress/RESULTS.md)
- [`experiments/gate-c-libkrun-console-lifecycle/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-console-lifecycle/RESULTS.md)
- [`experiments/gate-c-libkrun-installed-recovery/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-installed-recovery/RESULTS.md)
- [`experiments/gate-c-libkrun-adversarial/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-adversarial/RESULTS.md)
- [`experiments/gate-c-libkrun-adversarial/NULLFS_P0_2.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-adversarial/NULLFS_P0_2.md)
- [`experiments/gate-c-libkrun-supply-chain/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-supply-chain/RESULTS.md)
- [`experiments/gate-c-bun-runtime-authority/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-bun-runtime-authority/RESULTS.md)
- [`experiments/gate-c-libkrun-root-custody/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-root-custody/RESULTS.md)
- [`experiments/gate-c-p0-3-protocol-conformance/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-p0-3-protocol-conformance/RESULTS.md)
- [`experiments/gate-c-libkrun-console-correctness/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-console-correctness/RESULTS.md)
- [`Shrimpworks/libkrun` PR #2](https://github.com/Shrimpworks/libkrun/pull/2), exact
  [source head](https://github.com/Shrimpworks/libkrun/commit/8a2c91943793668f31a1cf7af431933be935bb58),
  and exact [merge commit](https://github.com/Shrimpworks/libkrun/commit/cf0333cdba478cc34a8570a65b38412da7fd3ecc)
- [Gate C implementation-readiness synthesis](../GATE_C_READINESS_CHECKPOINT.md)
- [Gate C P0 reconciliation](../GATE_C_P0_RECONCILIATION.md)

Repository-governance update: the retained PR links above identify the historical location used by
the evidence review. The current governed integration destination is
[`Shrimpworks/libkrun`](https://github.com/Shrimpworks/libkrun), transferred with PR history and
exact source/merge identities unchanged. This does not change this ADR's evidence or admission
state. Consumers must pin exact commits/digests and verify ancestry from upstream
`728df8125077d0db44265f6e997c72b81b65c015` instead of trusting the movable
`capsule/upstream-v1.19.4` branch.
