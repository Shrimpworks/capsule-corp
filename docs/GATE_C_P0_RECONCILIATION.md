# Gate C P0 reconciliation

Date: 2026-08-01
Evidence last reconciled: 2026-08-03

Status: planning decision with retained P0 evidence checkpoints after independent adversarial
review and targeted source research. This document refines the remaining work in the
[Gate C readiness synthesis](GATE_C_READINESS_CHECKPOINT.md). It records hypotheses to test; it is
not new backend evidence, a frozen libkrun profile, or permission to execute user bytes.

## Purpose

The five Gate C readiness tracks produced a conditionally viable native libkrun/HVF candidate and
identified five follow-up campaigns. An independent read-only review then challenged their scope,
ordering, threat coverage, and claim wording. Targeted research checked the pinned libkrun and Bun
surfaces and the likely no-host-root macOS packaging shape.

The review and research agree on the architectural bottom line: continue backend-independent
implementation now, but do not connect user bytes to libkrun until a smaller fail-fast P0 program
closes. Broad backend research is no longer the next step.

## Reconciled decisions

| Question | Decision | Evidence status |
| --- | --- | --- |
| Must Capsule fork libkrun for immutable disks? | The stock `/dev/fd/N` route passed the local identity/guest corpus but retained pathname semantics. A governed raw-only FD-native fallback is now a patch candidate. Prefer direct Supervisor-to-runner inheritance; add descriptor-transfer IPC only if the installed lifecycle requires post-spawn transfer. | FD-native local and owned-guest corpus passed; exact signed installed App Sandbox/protected-construction corpus remains untested. |
| Must source and inline input be block devices? | No for the first slice. Use bounded, attempt-bound virtio-console ports for source/input delivery. Keep the immutable block-custody problem scoped to the trusted runtime root and any later file-artifact storage. | A 43-vector backend-independent framing candidate conditionally passed; real transport, launcher, guest, and installed checks remain. |
| Must completion use runner status or a result disk? | No. Runner exit is lifecycle evidence only. Use one bounded typed completion/result frame on a dedicated port and bind it to the attempt and expected profile. | The backend-independent frame/commit model conditionally passed; it is not implemented in the real transport or launcher. |
| Is ext4 parsing P0 for inline JSON? | No. Carry bounded JSON in the typed result frame. A disposable ext4/raw-image parser becomes a gate for the later file-artifact slice. | Scope decision; no parser posture is promoted. |
| Is `NullFs` coupled to custody? | No. Carry the direct-block-root design as a narrow reviewed governed patch; fall back to complete virtiofs acceptance or reject the profile if that branch fails. | `GOVERNED-PATCH`: prototype removal is credible, not admitted; the current unpatched profile remains unsupported. |
| Does stock Bun satisfy the documented no-subprocess/no-FFI contract? | No. The exact Bun 1.3.14 source/binary investigation observed process, `execve`, FFI/native-loader, inspector, Worker, and descriptor authority despite all relevant stock flags. Preserve the contract; evaluate only a governed patched/external mechanism or alternate runtime, and change the contract only by ADR. | Stock Bun is rejected for this claim; `RUNTIME-001` remains unsupported. |
| Does no-host-root remain plausible? | Yes, provisionally. Continue the per-user Supervisor/runner/runtime topology and embedded `SMAppService`; do not add privilege silently if final packaging fails. | P0-4A conditionally passed 18-role topology enumeration, ad-hoc identity, registration, explicit activation, and same-session recovery. Valid Apple signing, App Sandbox, notarization, Gatekeeper, clean-host, on-demand activation, and supported-floor evidence remain open. |

The multiport API names observed in libkrun 1.19.4 are
`krun_add_virtio_console_multiport` and `krun_add_console_port_inout`. Their presence is not proof
that Capsule's framing, descriptor ownership, guest-node permissions, failure handling, or App
Sandbox profile is safe.

## First development-slice data path

```text
registered and approved source/input bytes
        │ bounded attempt-bound host writers
        ▼
dedicated virtio-console source/input ports
        │
        ▼
pinned kernel + trusted launcher ──fork/drop authority──▶ governed deno_core workload
        │ completion descriptor retained only by launcher
        ▼
one typed completion frame + commit trailer containing bounded inline JSON
        │
        ▼
host validation → teardown evidence → terminal transcript → Broker-held result
```

The pinned guest kernel and trusted launcher are part of this development profile's TCB. The
completion frame is not attestation against a compromised guest kernel. Host VMM containment must
still treat guest behavior as hostile.

## Authority and transport classes

The first slice does not use one generic custody mechanism for every byte class:

| Class | Authoritative producer | First-slice transport/custody | Gate |
| --- | --- | --- | --- |
| Runtime root | Admitted runtime-bundle manifest | Finalized unlinked object retained through a genuine read-only descriptor and attached through the governed raw-only FD-native API | P0-1 |
| Registered source | Exact Supervisor-registered source bytes | Bounded attempt-bound host-writer/guest-reader port | P0-3 |
| Approved inline input | Exact Broker-approved canonical JSON bytes | Separate bounded attempt-bound host-writer/guest-reader port | P0-3 |
| Completion/inline result | Trusted launcher under the admitted guest profile | Fresh launcher-owned guest-writer/host-reader port; unavailable to the workload | P0-3 |
| Temporary working data | Untrusted workload | Guest memory/tmpfs only; never extracted or released | P0-3 profile corpus |
| Future file artifacts | Untrusted workload | Fresh bounded raw storage plus a disposable bounded parser | Deferred file-artifact gate |

Each row needs its own digest, bound, ownership, lifecycle, and failure evidence. Passing runtime-
root custody does not prove port integrity; passing port integrity does not prove runtime-root
custody or future artifact safety.

## Attacker tiers for P0

| Tier | Capability | Required treatment |
| --- | --- | --- |
| Baseline same-user attacker | Arbitrary unprivileged process under the same login; ordinary same-UID pathname and directory access; races, links, replacement, and observation; retained descriptors or writable mappings to original user files; malformed IPC; component impersonation and attach attempts | The v0 design must resist it. A new custody object is not custodied while this attacker can acquire or retain writable authority to it. |
| Elevated user-granted attacker | Successful task-port/debug attachment, Full Disk Access, an explicit foreign-container grant, private VFS privilege, or another broad user-authorized capability | Separate elevated-adversary posture. P0 tests shipping denial and fail-closed behavior but makes no general resistance claim unless the exact capability is separately evidenced. |
| Trusted-platform compromise | Malicious root/administrator, SIP bypass, kernel/hypervisor compromise, or an authorized compromised Supervisor | Outside the local containment guarantee. It must not be confused with an ordinary same-user process. |

The baseline may already hold writable authority to the original selected file; that is why the
Broker snapshots bytes. It must not be able to acquire writable authority to the new custody object.
Shipping targets must use hardened runtime without `com.apple.security.get-task-allow` and must test
that attach attempts fail against the exact enrolled components.

## P0 admission program

### P0-0: runtime-authority closure

Goal: preserve the documented dependency-free JS/TS v0 profile with no subprocess, FFI, native addon,
inspector, macro, environment-file, or package-install authority.

Required evidence includes direct APIs, aliases, workers, child-process creation, native/FFI
loading, inspector activation, dynamic import/package behavior, environment/config discovery, and
launcher/descriptor inheritance under the exact runtime build. A finite API corpus alone cannot
prove absence. The mechanism needs a construction-level closure argument: enumerate the pinned
runtime's native-capability registry and syscall/module-loading paths, remove or governed-patch
prohibited entry points, apply external child-exec/native-loading denial where the exact launcher
can do so without breaking the JIT, review the source diff, and mutation-test deliberate restoration
of each denied primitive. Runtime flags are defense in depth unless the exact mechanism makes a
bypass structurally unavailable. Any runtime update invalidates this evidence.

Pass: an exact pinned mechanism refuses every prohibited power. Fail: choose a governed runtime
patch or alternate runtime, or explicitly revise the product contract in a new ADR. Do not keep
building around a runtime profile that contradicts the plan shown to the user.

The 2026-08-02 exact-stock investigation failed this hypothesis. Bun 1.3.14 commit
`0d9b296af33f2b851fcbf4df3e9ec89751734ba4` exposed direct and aliased subprocess, `execve`, FFI,
SQLite native loading, workload-started inspector, Worker, and inherited-descriptor authority under
all relevant stock deny flags. Addon, macro, and environment/config mutations showed that those
individual flags were active, but no stock `--no-spawn` or `--no-ffi` closure exists. See the
[retained P0-0 result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-bun-runtime-authority/RESULTS.md). `RUNTIME-001`
therefore remains unsupported and execution requiring it must refuse. That result left P0-0 open
only for a governed construction-level patch plus an exact external enforcement mechanism; failure
of that branch required an alternate runtime and an ADR-0003 update. This finding did not admit
runtime or backend bytes.

Governed-branch decision (2026-08-02): **NO-GO**. Exact source review found a conservative minimum
surface of 40 hand-authored files plus 10 generated outputs, spanning independent registries,
loader dispatch, native sinks, globals, configuration/resolution, build identity, and restoration
backstops. That triggered the campaign's broad/unreviewable fail-fast rule before a candidate diff
or governed binary existed. A narrow post-initialization process/exec self-seal remains plausible in
isolation, but cannot close Worker or native loading while preserving Bun/JSC lazy threads and JIT.
P0-0 is closed as a Bun NO-GO; alternate-runtime investigation and an ADR-0003 superseding decision
are now required. See the
[construction review](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-bun-runtime-authority/governed-closure/CONSTRUCTION_REVIEW.md).

Deno-family decision (2026-08-02): **DENO-FAMILY-NO-GO**. Under the exact full-Deno v2.9.4
profile, the initial static graph bypassed read/import denial, blob Workers executed, SIGUSR1
activated a loopback inspector, Node compatibility remained constructed, and runtime-managed web
storage APIs remained available outside ordinary read/write permissions (with state confined to
disposable tmpfs). A smaller Capsule-owned `deno_core` 0.409.0 prototype had no module
loader or ambient extensions and used V8 `--jitless`, but `JsRuntime` still physically registered
99 built-in core ops before middleware disabled 96. `deno_core` also has no TypeScript pipeline;
preserving exact approved-byte semantics requires a separately pinned pre-approval transformation
and coordinated plan/schema/ADR binding. Neither construction is selected, no Proposed ADR
supersedes ADR-0003, and `RUNTIME-001` continues to refuse. The governed physical-omission follow-up
then reduced the registry to three bootstrap ops and reproduced its snapshot/binary. A later
packaging experiment replaced the local-only builder with a digest-pinned no-apt recipe and complete
offline Cargo source bundle, reproducing the same bytes in two clean same-host containers. The
later self-contained-root experiment closed the standalone dynamic-root item with an exact
package-derived 22-entry root, explicit cache-free loader invocation, and file-open/mutation
evidence. The stronger admission-evidence gate still fails on independent-builder provenance,
exact archive-corresponding V8 source/notices, a fork-native Linux/arm64 builder/release, and the
remaining admission corpus. The governed fork branches now exist, but the first fork-native
integration check stopped before building because the `rusty_v8` fork supports only Linux/amd64.
See
the [retained Deno-family result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-runtime-authority/RESULTS.md),
[physical-omission result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-core-physical-omission/RESULTS.md), and
[package result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-core-reproducible-package/RESULTS.md), and
[self-contained-root result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-core-runtime-root/RESULTS.md).

Exact V8 closure follow-up (2026-08-02): **SOURCE-LICENSE-CLOSURE-NO-GO**. The official
Linux/arm64 asset is now bound to its successful release job, exact `rusty_v8` commit, 20 recursive
gitlinks, exact Deno V8 commit, Chromium V8 base, and four-patch stack. The 1,875 archive members and
726 source-tree license/notice candidates were inventoried. The release still omits immutable
runner/package/action resolutions, effective GN/Ninja closure, generated notices, and a
corresponding-source bundle, so an exact rebuild was not possible and PR #50's CycloneDX
composition remains incomplete. Accepted ADR-0028 selects governed `deno_core` as the first
engineering candidate and supersedes ADR-0003's Bun-first ordering only; it does not admit a
profile or change `RUNTIME-001`. See
the [retained closure result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-v8-source-license-closure/RESULTS.md).

Fork-native integration follow-up (2026-08-03): **LINUX/ARM64 CONSTRUCTION BLOCKED; NO BUILD OR
RUNTIME ADMISSION**. Exact public refs, merge parents, upstream ancestry, Deno three-op/fixture
oracles, the `rusty_v8` 20-gitlink source lock, and existing tool locks passed independent
inspection at governed Deno head `9adb0b68...91bed` and governed `rusty_v8` follow-up head
`a43ee748...33cf`. The latter's only profile is Linux/amd64:
`x86_64-unknown-linux-gnu`, including its builder image, LLVM/bindgen, sysroot, GN/Ninja, output,
collection, and provenance paths. The experiment therefore performed no prefetch or build and did
not substitute amd64. The smallest next fork change is a digest-pinned Linux/arm64 sibling profile
with network-disabled compilation/test/evidence collection. `RUNTIME-001` remains unsupported.
See the [retained blocker](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-fork-native-deno-runtime-bundle/RESULTS.md).

Current external-work note (2026-08-03): governed `rusty_v8` PR #4 now carries that follow-up at
exact unmerged head `aa921fa48901bf28774d61248b0187c8b91c55a4`. Its contract jobs pass while clean
Linux/arm64 full-build work remains in progress. This is not a durable Capsule integration, accepted
handoff, reusable artifact, release, or admission result. The retained blocked experiment above is
unchanged until exact successful evidence is reviewed and merged.

Governed `deno_core` follow-up (2026-08-02): **PHYSICAL-OMISSION-PASS; NO RUNTIME ADMISSION**.
The exact governed patched construction reduced the built-in registry from 99 ops to the three
bootstrap-required ops with a one-file physical-omission patch; runtime/symbol inspection observed
only those three, fixed
restoration mutations failed closed, and ASLR-controlled clean builds reproduced the snapshot and
binary. This closes only the pre-registration/final-link question. TypeScript approved-byte
semantics, independently reconstructible packaging/provenance, complete restoration/backstop
review, external isolation composition, and runtime-profile admission remain open. The later
ADR-0028 ordering decision does not promote this narrow result or `RUNTIME-001`. See the
[retained physical-omission result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-core-physical-omission/RESULTS.md).

Approved-byte follow-up (2026-08-03): **BOUNDARY-PASS; NO RUNTIME ADMISSION**. Exact Node
22.22.1/Amaro 1.1.5 `stripTypeScriptTypes` in strip-only mode emitted byte-identical fixed outputs
across repeated processes, refused transform-requiring syntax and unknown options, preserved
Unicode/line endings without normalization, and detected source/output/toolchain/options
mutations. Proposed ADR-0026 therefore places transformation before executable-source manifest and
plan construction and binds original, emitted, options, transformer, diagnostics-zero, and explicit
source-map absence into the future plan. Current schemas/types remain unchanged, production
ownership/topology and cross-platform provenance remain open, no `deno_core` module loader or
runtime was wired, and `RUNTIME-001` still refuses. See the
[retained result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/typescript-approved-byte-boundary/RESULTS.md).

### P0-1: immutable runtime-root custody

Goal: prove that a concurrent same-user attacker cannot change or substitute bytes observed through
libkrun's pathname disk API.

This gate has three independent claims. `/dev/fd/N` can pass attachment identity without passing
frozen-object construction.

#### P0-1A: stable attachment identity

Prove that both libkrun internal opens duplicate the exact originally read-only descriptor and never
fall back to the old image pathname. Test descriptor lifetime across VM start, wrong and reused
descriptor numbers, close-on-exec/fork behavior, both internal consumers, shared-open-description
offset/status behavior, and every positional-I/O assumption.

Direct Supervisor-to-runner inheritance is the first topology. Add authenticated descriptor-
transfer IPC only if the installed lifecycle proves that post-spawn transfer is required.

#### P0-1B: frozen-object construction

The first candidate sequence is:

1. Create an unguessable file with exclusive creation inside Supervisor-protected storage.
2. Populate the object from the admitted root bytes and finish every write through the creator
   descriptor.
3. While the protected pathname still exists, open a distinct descriptor whose
   `fcntl(F_GETFL) & O_ACCMODE` is exactly `O_RDONLY`; separately verify file permission mode,
   device, inode, regular-file type, link count of one, and expected length.
4. Close and unmap every writable descriptor, alias, and mapping owned by the construction path.
5. Unlink the only pathname and verify through the retained descriptor that its access mode is
   still exactly `O_RDONLY`, its device/inode/type/permission mode/length are unchanged, and link
   count is now zero.
6. Only after finalization, compute the complete digest and length through that exact retained
   read-only descriptor and compare them with the admitted runtime-root manifest.
7. Bind role, digest, length, device/inode observation, descriptor-transfer mechanism, runtime
   profile, and attempt into the durable launch record.
8. Inherit that descriptor into the signed runner, attach `/dev/fd/N`, and compare the guest-observed
   digest with the finalized custody digest.

Hashing before finalization, hashing a different pathname/object, or relying only on a post-stop
hash does not close custody.

#### P0-1C: adversarial end-to-end custody

The corpus includes pre-creation name substitution, the create-to-read-only window, hard links,
writable mappings, path replacement, inherited writable aliases, runner/Supervisor crashes,
debugger or Mach task-port attachment attempts, explicit container grants, descriptor close/reuse,
shared-offset interference, guest read-before-finalization, recovery pathname reacquisition, and
installed App Sandbox behavior. `O_CREAT|O_EXCL`, `chmod`, secrecy, read-only guest flags, or
post-stop hashing alone do not pass.

Pass only if P0-1A through P0-1C all pass: the admitted manifest digest equals the exact finalized
descriptor digest, the guest-observed digest equals it, no baseline same-user attack can alter it,
and recovery never reconstructs object authority from a pathname.

Fallback: one narrowly governed FD-native libkrun API change followed by the same corpus. If both
paths fail, reject libkrun for the v0 profile. Host-root execution, a separate-owner host service,
or another privileged host helper is a hard v0 non-option; considering one later requires a new ADR
comparing a higher macOS floor, memory-backed storage, a different backend, narrower same-user
claims, and the complete new installation/update/recovery boundary.

Evidence checkpoint (2026-08-02): the fallback is a **PATCH-CANDIDATE**, not a P0-1 pass. Against
libkrun commit `728df8125077d0db44265f6e997c72b81b65c015`, patch SHA-256
`48cdbc307b3fa1209fa0ec68fc3f817634af312983d68f0de259db86c0b43333` adds only the fixed
`runtime-root:vda:raw:read-only` API. It takes an owned `F_DUPFD_CLOEXEC` duplicate, validates exact
finalized descriptor identity, and constructs raw imago storage directly from an owned `File`.
The controlled C/Rust/local custody and five-mutation corpus passed; four owned unsandboxed HVF
guest runs matched finalized host, guest `/dev/vda`, and post-stop digests with zero root-path
opens. The host has no valid signing identity, so the ad-hoc App Sandbox bundle aborted before
`main`; protected construction, exact installed descriptor manifest, task-port/grant denial, and
final signed/notarized bytes remain mandatory. See the
[retained result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-root-custody/RESULTS.md).

### P0-2: `NullFs` disposition

Goal: eliminate the unexpected device or accept only a bounded, understood VMM surface.

Pass by either removing the device and rerunning the full device/cross-job corpus, or by formally
accepting an exact residual surface. Acceptance pins the libkrun commit and build features, proves
`shared_dir=None`, makes every host-directory configuration route unreachable through Capsule,
reviews the complete guest-reachable virtiofs transport/queue/FUSE decoder/worker/overlay/`NullFs`
path, and runs malformed opcode, length, descriptor, queue, concurrency, cancellation, and resource-
exhaustion fuzzing under applicable sanitizers. Retain coverage, corpus, crashes/fixes, zero
unresolved high-severity findings, and residual limitations. Describe the surface as accepted, not
absent or proven bug-free. Custody and `NullFs` remain independent decisions even if one libkrun
patch eventually changes both areas.

Earlier evidence checkpoint (2026-08-02): the bounded replacement investigation found no independent
virtiofs feature toggle and confirmed that `krun_set_root_disk_remount` always adds `NullFs` in the
exact retained block-root route. Removing only that device built but failed before init because the
dummy virtiofs root supplies the bootstrap file and mount points. This falsifies the smallest
removal, not all alternate-bootstrap designs. The existing 47-device-test/one-libkrun-test suites
and Go profile fuzzing do not cover the complete FUSE/queue/worker/overlay path, and the pinned tree
has no relevant fuzz target or retained sanitizer/coverage corpus. See the
[P0-2 investigation](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-adversarial/NULLFS_P0_2.md).

Later disposition checkpoint (2026-08-02): **`GOVERNED-PATCH`; removal credible but not admitted**.
A direct-block-root prototype placed the trusted init inside the immutable root, booted
`/dev/vda` directly, and remounted `/` `ro,nosuid,nodev`. It exposed only balloon, RNG, console,
and block devices, no virtiofs device or mount, denied the network probe, exposed no usable vsock,
and reran 36 adversarial plus four identity cases without the original `NullFs` failure. The three
failures were expected ad-hoc-signing identity limitations. This does not close P0-2: the patch
needs independent source/API review, route-closure mutations, final P0-1 custody, P0-3 transport,
and the complete signed/notarized P0-4 corpus. The current unpatched profile remains unsupported.
See the [retained disposition](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-adversarial/NULLFS_P0_2_DISPOSITION.md).

### P0-3: typed port transport and completion

Goal: move only exact bounded source/input bytes into the guest and receive exactly one
attempt-bound completion/result without ambient network or vsock authority.

Before implementation, freeze independent exact maxima for source bytes, canonical input bytes,
the completion physical frame, and the inline-result JSON payload; every channel must fail rather
than resize dynamically. Each source/input envelope binds version, role, attempt, registration/plan,
runtime profile, length, and digest. The completion protocol additionally binds terminal status and
ends with an unambiguous fixed commit trailer written only after the complete payload. The host
continuously drains up to cap-plus-one and never uses stream EOF as proof of a complete frame. This
does not choose byte layouts or cap values without measurement.

Pinned libkrun 1.19.4 source establishes known hazards that this gate must resolve rather than
rediscover accidentally. Its console output wait polls only the output descriptor, without the stop
event, while port shutdown joins the transmit thread before waking the stop event. Successful short
writes can advance, but partial progress followed by error is not accounted exactly and a partial-
then-zero result can reach a complete-write assertion. The API also duplicates supplied descriptors,
sets `O_NONBLOCK` on the shared open-file description, and treats a negative descriptor as an absent
direction only in pinned implementation—not in its public contract. Pass only with a governed
directional API or a pinned canary for that exact behavior, dedicated endpoints with recorded pre/
post flags, and proof that an unused direction has no host source/sink or cross-role effect.

The guest-facing transport implementation is itself hostile-input VMM surface. The pinned console
control handler indexes ports and queues from guest-supplied IDs without a demonstrated bounds check.
Audit and sanitizer/coverage-fuzz control IDs/events, descriptor chains, queue sizes and indices,
port open/close/reset ordering, concurrency, cancellation, and resource exhaustion. Retain coverage,
crashes/fixes, zero unresolved high-severity findings, and limitations. Pass only if continuous
draining, reader stall/death, cap-plus-one closure, partial-then-error/zero, `SIGPIPE`, invalid/
closed/swapped descriptors, shutdown, and external exact forced teardown remain bounded and fail
closed. Otherwise use a governed bounds-checked, stop-aware, partial-write-correct port change or
reject this transport for v0.

Before guest implementation, freeze the trusted launcher/runtime-adapter protocol. The launcher
must remain a distinct process rather than replacing itself with Bun: fully receive and verify both
envelopes before child start; bind immutable source/input objects or descriptors to fixed argv,
empty/fixed environment, cwd, runtime identity, and child FD manifest; withhold the completion port
and node from the child; collect only fixed bounded result/stdout/stderr endpoints; wait for exact
child-tree termination; and emit the commit trailer last. Test swap, truncation, `/proc` access,
launcher/child crash, descriptor inheritance, child/subprocess forgery, and output before final
input/source verification.

The host runner also starts from a closed role-specific FD manifest because a VMM compromise gains
every descriptor it inherits. After exec, allow only the finalized runtime-root descriptor,
directional port endpoints, and indispensable runtime/control descriptors; prove access modes,
flags, ownership, and roles and fail on every unexpected database, key, XPC, log, temporary, or
writable descriptor. Enumerate the actual set in clean launch, crash, recovery, and update cases.

The corpus covers wrong/duplicate/stale attempts, partial and oversized frames, malformed lengths
and JSON, early/missing/duplicate commit trailers, trailing data, output floods, host-reader stalls/
death, launcher crash before and after the commit trailer, corrupt root, missing executable, guest
panic/OOM, descriptor leakage, child inheritance, workload direct-open, subprocess/FFI forgery
attempts, and terminal classification when teardown or input integrity is missing. The unprivileged
workload must not own the completion descriptor or port node.

Pass: ordinary success requires a valid frame plus separate input-integrity, bounded-result,
runtime-integrity, runner-lifecycle, and teardown dispositions. Runner exit zero never substitutes
for a missing frame.

Evidence checkpoint (2026-08-03): **conditional pass for a falsifiable backend-independent
candidate; P0-3 remains open**. Independent Go and Node models verify the same 43 byte-exact vectors,
and Node independently encodes six accepted known answers. The local process-pipe corpus adds
partial writes, zero progress/stall, reader death, peer-close `EPIPE`, backpressure, cancellation,
runner death before/after commit, three-way role confusion, and EOF/clean-exit refusal. It measured
1,048,576 source bytes, 262,144 canonical-input bytes, 262,144 inline-result JSON bytes, and a
262,368-byte physical completion frame; exact boundaries passed and cap-plus-one was fully drained
and refused. No virtio-console, launcher, runtime, guest, VMM, App Sandbox, Supervisor, approval,
or product teardown mechanism participated. See the
[retained protocol result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-p0-3-protocol-conformance/RESULTS.md).

The sibling console review found that stock libkrun cannot proceed as-is. Governed patch SHA-256
`584ce48548fe969684fe3c55e57fbf56e7dae40af28c241c24c47b138faf1283` passed 51 local library
tests, all 51 under AddressSanitizer, warning-denying Clippy with the known deprecated-call
allowance, 25 shutdown repetitions, and four caught restoration mutations. Exact coverage was only
90/728 lines (12.362637%) across the four patched files, with `port.rs` and `process_tx.rs` at zero.
That retained local result is the before measurement, not the current governed-fork coverage state.
See the [retained console result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-console-correctness/RESULTS.md).

Governed-fork reconciliation (public state read 2026-08-03T22:57:43Z):

Repository-governance update: the historical public read below occurred at `dills122/libkrun`.
The current governed integration destination is
[`Shrimpworks/libkrun`](https://github.com/Shrimpworks/libkrun), still a public fork of
`libkrun/libkrun`. PRs #1 and #2 and their exact identities transferred unchanged;
`capsule/upstream-v1.19.4` now points to merge
`cf0333cdba478cc34a8570a65b38412da7fd3ecc`, whose verified ancestry reaches exact upstream anchor
`728df8125077d0db44265f6e997c72b81b65c015`. This changes repository ownership only. Any reuse must
pin exact commits/digests and verify ancestry rather than trusting that movable branch name.

- [`dills122/libkrun` PR #2](https://github.com/dills122/libkrun/pull/2) merged exact source head
  `8a2c91943793668f31a1cf7af431933be935bb58` as merge commit
  `cf0333cdba478cc34a8570a65b38412da7fd3ecc`. The merge parents are governed base
  `4ea8d1de861ed1c0636fc800b6da8fb71a086aa5` and that exact head. At read time the source branch
  still resolved to the head and `capsule/upstream-v1.19.4` resolved to the merge commit.
- The five retained patch files and aggregate patch-set identity
  `d19fd0ff159c699acccda2621519de45a09408bf3847b418ac34e02b79e805d5` are unchanged. Their
  reconstruction still begins at upstream `728df8125077d0db44265f6e997c72b81b65c015` and matches
  governed base `4ea8d1de861ed1c0636fc800b6da8fb71a086aa5`; the follow-up fixes and tests are source changes
  layered above that base, not rewritten experiment patch bytes.
- The bounded property corpus exposed and fixed two local defects: shutdown racing a queued,
  backpressured zero-progress descriptor could re-pop and retry it, and a port whose workers had
  been joined remained logically active. The merged source checks stop before requeue, parks
  zero-progress work until signaled, and moves the port to `Inactive` as shutdown takes its worker
  handles.
- Four-file console coverage moved from 13/88 functions (14.772727%) and 90/728 lines
  (12.362637%) to 37/88 functions (42.045455%) and 298/733 lines (40.654843%). `port.rs` moved from
  0/17 functions and 0/137 lines to 15/17 functions (88.235294%) and 111/137 lines (81.021898%);
  `process_tx.rs` moved from 0/4 functions and 0/91 lines to 4/4 functions (100%) and 82/96 lines
  (85.416667%). The measured remaining gaps are 2 functions/26 lines in `port.rs` and 0
  functions/14 lines in `process_tx.rs`.
- The governed no-guest wrapper reports 51 retained default tests, 53 `blk` tests including the two
  existing raw-FD tests, four bounded console/property tests, two raw-FD library-boundary tests,
  warning-denying Clippy with the retained deprecation allowance, AddressSanitizer over the
  51+4+2 macOS-arm64 sets, 25/25 shutdown and 25/25 queued-backpressure repetitions, five caught
  raw-FD mutations, four caught console restoration mutations, formatting/source/header contracts,
  exact patch reconstruction, and reverse dry-run. The non-`blk` raw-FD target compiled with zero
  gated tests; no guest was run.
- GitHub reported 16 head check runs: 13 successes, two intentionally skipped guest-running
  integration jobs, and the self-hosted `Unit tests (Linux aarch64)` job still queued with no steps
  or conclusion. The overall head check state was therefore pending. The macOS cross-compilation
  job, including the default Linux init build, succeeded; both push and pull-request governed
  no-guest jobs succeeded. GitHub reported no checks for the merge commit itself.
- Those successes are pre-merge head evidence, not a successful post-merge baseline rerun. At the
  exact merge, `capsule/upstream-v1.19.4` resolves to `cf0333cdba478cc34a8570a65b38412da7fd3ecc`,
  while `verify-patch-queue.sh` still requires that branch to equal
  `4ea8d1de861ed1c0636fc800b6da8fb71a086aa5`. Read-only inspection therefore predicts—and the
  script deliberately defines—a fail-closed `governed baseline branch moved` result before patch
  reconstruction on the current branch tip. The unchanged local five-patch aggregate still
  reproduces exactly; branch-policy consistency does not.
- GitHub reported no submitted review, no requested reviewer/team, and the exact CODEOWNERS rule
  names `@dills122`, who also authored and merged the PR. Independent human review and an evidenced
  CODEOWNER review therefore remain open; merge and green checks do not satisfy that governance
  requirement.

This closes the stale zero-coverage and archived-only descriptions, not P0-3. Measured uncovered
code, the governed branch/verifier invariant, caller-flag closure, control/queue/descriptor
fuzzing, real virtio transport, the distinct launcher and child/runner FD manifests,
guest/VMM/forced-teardown behavior, installed App Sandbox, and final profile evidence remain open.
Fork source identity, local library tests, CI build evidence, earlier owned-guest custody evidence,
and product/backend admission are separate classes.

### P0-4: complete installed development bundle

Goal: admit the exact bytes and topology that will run the first development slice without host-
root authority.

P0-4A builds an early signed installed-topology harness containing the Supervisor, runner,
libkrun/libkrunfw, firmware/kernel/root, candidate runtime and launcher, entitlements, manifests,
and embedded per-user service registration. It tests App Sandbox, service, descriptor-manifest,
process-identity, recovery, and minimum-OS assumptions but cannot admit the backend.

Evidence checkpoint (2026-08-02): **conditional P0-4A topology pass; no admission**. The retained
harness enumerated 18 roles and verified 17 non-self-referential installed entries, kept
`backendAdmitted=false`, registered and explicitly activated an embedded per-user `SMAppService`,
enforced exact ad-hoc CDHash identity in both IPC directions, refused stale/missing/mixed/unexpected
components and FD 8, and recovered after same-session Supervisor death. It created no guest and did
not load libkrun. With zero valid signing identities, the App-Sandboxed runner failed before
`main`, Gatekeeper rejected the app, pure on-demand activation remained open, and the exact assembled
candidate inherited a macOS 26.0 floor from libkrunfw. No notarized, stapled, clean-host, supported-
floor, session, or P0-0-through-P0-3-composed evidence exists. See the
[retained P0-4A result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-installed-development-topology/RESULTS.md).

After P0-0 through P0-3 select exact mechanisms and patches, P0-4B rebuilds the complete final app
with source/license material and reruns every affected P0 gate on the final signed/notarized bytes.
Any runtime, kernel, libkrun patch, entitlement, helper topology, descriptor manifest, or minimum-OS
change invalidates dependent earlier evidence. Only P0-4B can admit the development profile.

Pass requires pinned materials, governed patches, manifest/SBOM/provenance/source completeness,
minimum-OS builds, signing and notarization, staple/Gatekeeper assessment, installed-byte readback,
clean-host launch, inherited-descriptor behavior, exact process identity, recovery, and an explicit
macOS support floor. The resulting `BackendValidationRecord` verdict is at most
`development-admitted`; it binds exact final bytes/configuration, passed claims, known limitations,
evidence digests, expiry and invalidation triggers and cannot authorize `validated-local`. A runner-
only ticket is not sufficient.

## Exact next work and credential boundary

The next non-credential work is:

1. repair the governed baseline/ref/verifier invariant without rewriting the retained five patch
   bytes, then independently review the exact merge and full FD-native, direct-block-root, console,
   and follow-up source/API composition; explicitly map the follow-up deltas and close the measured
   uncovered lines/functions with bounded tests and mutations;
2. carry the independently reproduced 43 P0-3 vectors into the selected host/launcher languages,
   complete hostile control/queue/descriptor coverage, and implement the distinct launcher, child
   manifest, and exact runner FD manifest before any real-guest composition;
3. wait for an exact successful and reviewed handoff from unmerged governed `rusty_v8` PR #4; only
   then reconstruct the merged governed Deno candidate,
   deterministic V8 publication/link/notice evidence, and runtime-root packaging,
   complete restoration/backstop review, assign and wire the retained approved-byte TypeScript
   candidate, and finish external-isolation/profile composition without weakening `RUNTIME-001`;
   and
4. finish the final role, descriptor, runtime/kernel/init, source/license, SBOM/provenance, and
   minimum-OS build inputs so the exact package can be rebuilt before signing.

The exact sequencing and fork exit criteria are recorded in the
[governed deno_core work plan](GOVERNED_DENO_CORE_WORK_PLAN.md).

The signing-dependent blockers are a valid Apple Development or Developer ID signing identity for
the hardened App Sandbox runner, protected-container custody test, and Team-enrolled installed IPC;
Developer ID signing/notarization credentials for complete-bundle notarization, stapling, and
Gatekeeper assessment; and installed-byte, clean-host, session/recovery, and support-floor
validation on those exact signed bytes. P0-4B and final P0-1/P0-2/P0-3 installed reruns cannot pass
without that evidence. Those identities/profiles must be deliberately authorized. Exact G3
discovery found that certificate SHA-1 `1638...61E3` is displayed with a W4 suffix, but its X.509
subject OU and emitted code-signing TeamIdentifier are `3DDR84M4JS`; it is not W4 evidence. The
three Xcode 26.6-cached profiles all belong to historical Team `3DDR84M4JS` and are likewise not
reusable for W4 tests. Local W4 work requires a matching/reissued certificate and exact profiles,
not only role identifiers and entitlements. A separate
Developer ID Application identity for historical Team `3DDR84M4JS` remains later distribution
authority requiring explicit authorization and matching-Team package design; it is not W4
development evidence and does not make Developer ID or notarization work current. Paid owned clean-host/
minimum-OS validation is not currently planned and remains deferred activation/distribution
evidence rather than a blocker for local mechanics. A
genuinely independent Linux/arm64 builder is viable but not currently planned; same-host and
GitHub-CI equality therefore remain limited and independent-builder equality is deferred.

## Work that starts before P0 closes

Backend-independent Phase 2 through Phase 4 work may proceed against a hard-fenced fake backend:

- object-specific contracts and fixtures;
- execute-by-registration-ID and immutable registered bytes;
- one-use approval consumption and attempt creation;
- fault-injectable lifecycle, durable cleanup obligations, reconciliation, and quarantine;
- composed Broker/Supervisor evidence without guest content; and
- bounded inline JSON ownership, user-only result storage, and fixed agent summaries.

The fake backend must report `CreatesGuest() == false` and must have no link or configuration path
that can launch a VMM. The real libkrun adapter begins only after every P0 gate above passes.

## Deferred campaigns

The first inline slice does not require a filesystem artifact parser. Before regular-file or file-
artifact support, Capsule must add the bounded disposable ext4/raw-image parser campaign. Before
`validated-local` or production claims it must additionally run end-to-end real-backend fault
injection, installed session/reboot matrices, APFS/power-loss pressure, the complete hostile-runtime
corpus, a real checksum-pinned gVisor comparison, release/update/revocation drills, and quantitative
soak/resource-leak testing.

## Claim boundary

This reconciliation supports only a plan of work. It does not establish that `/dev/fd/N` is safe
under App Sandbox, that the port protocol is isolated, that Bun powers are disabled, that
`NullFs` is harmless, that a complete package is admissible, or that libkrun is validated or
production-ready. Failure of a P0 mechanism produces an explicit pivot decision; it is not a
reason to reinterpret missing evidence as success.

In particular, do not claim:

- that `/dev/fd/N` alone provides immutable custody;
- that a guest-read-only virtio block device is host-immutable;
- that opening `/dev/fd/N` from an `O_RDWR` descriptor independently downgrades its authority;
- that completion proves the approved code ran correctly or the guest kernel was uncompromised;
- that a profile is “same-UID safe” without naming the attacker tier; or
- that source/input port evidence covers runtime-root custody or future artifact parsing; or
- that a virtio-console stream provides message boundaries, completion, or safe teardown by itself.

## Primary research references

- [libkrun 1.19.4 public header](https://raw.githubusercontent.com/libkrun/libkrun/v1.19.4/include/libkrun.h)
- [libkrun 1.19.4 console control path](https://github.com/libkrun/libkrun/blob/v1.19.4/src/devices/src/virtio/console/device.rs)
- [libkrun 1.19.4 console transmit path](https://github.com/libkrun/libkrun/blob/v1.19.4/src/devices/src/virtio/console/process_tx.rs)
- [libkrun 1.19.4 console port I/O](https://github.com/libkrun/libkrun/blob/v1.19.4/src/devices/src/virtio/console/port_io.rs)
- [libkrun 1.19.4 console shutdown ordering](https://github.com/libkrun/libkrun/blob/v1.19.4/src/devices/src/virtio/console/port.rs)
- [Apple XNU `/dev/fd` implementation](https://github.com/apple-oss-distributions/xnu/blob/main/bsd/miscfs/devfs/devfs_fdesc_support.c)
- [Apple XNU descriptor duplication](https://github.com/apple-oss-distributions/xnu/blob/main/bsd/kern/kern_descrip.c)
- [POSIX `dup` open-file-description semantics](https://pubs.opengroup.org/onlinepubs/9799919799/functions/dup.html)
- [POSIX nonblocking and partial-write semantics](https://pubs.opengroup.org/onlinepubs/9690949599/functions/write.html)
- [Virtio 1.3 console specification](https://docs.oasis-open.org/virtio/virtio/v1.3/virtio-v1.3.html)
- [Bun FFI documentation](https://bun.sh/docs/runtime/ffi)
- [Bun child-process documentation](https://bun.sh/docs/runtime/child-process)
- [Apple protected app-container guidance](https://developer.apple.com/documentation/xcode/protecting-local-app-data-using-containers)
- [Apple App Sandbox file access and macOS 14 container association](https://developer.apple.com/documentation/security/accessing-files-from-the-macos-app-sandbox)
- [Apple macOS 15 app-group container protection](https://developer.apple.com/documentation/macos-release-notes/macos-15-release-notes)
- [Apple notarization guidance for `get-task-allow`](https://developer.apple.com/documentation/security/resolving-common-notarization-issues)
- [Apple `SMAppService`](https://developer.apple.com/documentation/servicemanagement/smappservice)
