# Gate C P0 reconciliation

Date: 2026-08-01

Status: planning decision after independent adversarial review and targeted source research. This
document refines the remaining work in the
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
| Must Capsule fork libkrun for immutable disks? | Test a genuine inherited read-only descriptor exposed as `/dev/fd/N` first. Prefer direct Supervisor-to-runner inheritance; add descriptor-transfer IPC only if the installed lifecycle requires post-spawn transfer. A narrow FD-native libkrun change is the fallback. | Promising hypothesis; not yet validated under the exact App Sandbox/install profile. |
| Must source and inline input be block devices? | No for the first slice. Use bounded, attempt-bound virtio-console ports for source/input delivery. Keep the immutable block-custody problem scoped to the trusted runtime root and any later file-artifact storage. | The pinned libkrun header exposes generic multiport console APIs; the exact protocol and hostile cases remain untested. |
| Must completion use runner status or a result disk? | No. Runner exit is lifecycle evidence only. Use one bounded typed completion/result frame on a dedicated port and bind it to the attempt and expected profile. | Required design; not implemented. |
| Is ext4 parsing P0 for inline JSON? | No. Carry bounded JSON in the typed result frame. A disposable ext4/raw-image parser becomes a gate for the later file-artifact slice. | Scope decision; no parser posture is promoted. |
| Is `NullFs` coupled to custody? | No. Remove it or independently review and fuzz the exact guest-facing virtiofs parser surface. | Existing exact profile remains failed until this closes. |
| Does stock Bun satisfy the documented no-subprocess/no-FFI contract? | No. The exact Bun 1.3.14 source/binary investigation observed process, `execve`, FFI/native-loader, inspector, Worker, and descriptor authority despite all relevant stock flags. Preserve the contract; evaluate only a governed patched/external mechanism or alternate runtime, and change the contract only by ADR. | Stock Bun is rejected for this claim; `RUNTIME-001` remains unsupported. |
| Does no-host-root remain plausible? | Yes, provisionally: test one complete signed/notarized app containing the enrolled per-user Supervisor/runner/runtime topology and an embedded `SMAppService` registration. Do not add privilege silently if it fails. | Packaging hypothesis; clean-host and minimum-OS evidence remain open. |

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
pinned kernel + trusted launcher ──fork/drop authority──▶ untrusted Bun workload
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
| Runtime root | Admitted runtime-bundle manifest | Finalized unlinked object retained through a genuine read-only descriptor and attached through `/dev/fd/N` | P0-1 |
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

Goal: preserve the documented dependency-free Bun profile with no subprocess, FFI, native addon,
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
building around a stock Bun profile that contradicts the plan shown to the user.

The 2026-08-02 exact-stock investigation failed this hypothesis. Bun 1.3.14 commit
`0d9b296af33f2b851fcbf4df3e9ec89751734ba4` exposed direct and aliased subprocess, `execve`, FFI,
SQLite native loading, workload-started inspector, Worker, and inherited-descriptor authority under
all relevant stock deny flags. Addon, macro, and environment/config mutations showed that those
individual flags were active, but no stock `--no-spawn` or `--no-ffi` closure exists. See the
[retained P0-0 result](../experiments/gate-c-bun-runtime-authority/RESULTS.md). `RUNTIME-001`
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
[construction review](../experiments/gate-c-bun-runtime-authority/governed-closure/CONSTRUCTION_REVIEW.md).

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

Evidence checkpoint (2026-08-02): the bounded replacement investigation found no independent
virtiofs feature toggle and confirmed that `krun_set_root_disk_remount` always adds `NullFs` in the
exact retained block-root route. Removing only that device built but failed before init because the
dummy virtiofs root supplies the bootstrap file and mount points. This falsifies the smallest
removal, not all alternate-bootstrap designs. The existing 47-device-test/one-libkrun-test suites
and Go profile fuzzing do not cover the complete FUSE/queue/worker/overlay path, and the pinned tree
has no relevant fuzz target or retained sanitizer/coverage corpus. P0-2 therefore remains open and
the current profile remains unsupported. See the
[P0-2 investigation](../experiments/gate-c-libkrun-adversarial/NULLFS_P0_2.md).

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

### P0-4: complete installed development bundle

Goal: admit the exact bytes and topology that will run the first development slice without host-
root authority.

P0-4A builds an early signed installed-topology harness containing the Supervisor, runner,
libkrun/libkrunfw, firmware/kernel/root, candidate runtime and launcher, entitlements, manifests,
and embedded per-user service registration. It tests App Sandbox, service, descriptor-manifest,
process-identity, recovery, and minimum-OS assumptions but cannot admit the backend.

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
