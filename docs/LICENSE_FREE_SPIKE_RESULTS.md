# License-free feasibility spike results

Date: 2026-07-31

Status: research synthesis. This records development evidence and current decisions; it does not
promote any control, backend, profile, or component to production-ready.

## Outcome

An Apple Developer Program account is **not blocking the next implementation slice**. Local ad-hoc
signing and public macOS APIs were enough to test the serialization profile, exact live-process and
XPC code requirements, message-derived peer identity, cross-process descriptor custody, direct
VM-backed Containerization, no-network and resource controls, bounded output, controller crashes,
and real `SIGKILL` recovery of the trust-transition stores.

Apple-credentialed and second-wave follow-ups observed installed per-user XPC services,
Team/channel/identifier requirements, provisioned disjoint Keychain groups, protected app
containers, persistent Secure Enclave keys, interactive approval, and crash-safe security-epoch key
transitions. Supervisor notarization completed; Broker and daemon remain in Apple's queue.

The major negative result remains decisive for stock Apple Container and direct Containerization:
neither provides the required durable host-side VM/helper identity or restart enumeration. A later
libkrun/Hypervisor.framework spike avoided that specific failure by making one signed VMM process
the VM lifecycle object. Its isolation, App Sandbox, exact teardown, controller-crash handshake,
resource, and Bun checks conditionally passed. Five implementation-readiness tracks then observed
bounded scratch/output and console mechanics, exact forced teardown, installed crash recovery, an
adversarial corpus, and reproducible-build feasibility. They also found a live mutable-path input
race, an unexpected block-root `NullFs` virtiofs device, ambiguous runner-zero completion,
non-product output parsing, incomplete distribution, and a no-go for the current runtime bytes.
libkrun/HVF remains the lead native candidate under evaluation; OCI plus gVisor remains an
independent comparison and contingency. Neither has a production posture. A later independent P0
review and targeted source research narrowed the first inline slice to port-based source/input/results,
deferred filesystem-image parsing until file artifacts, and added stock-Bun runtime-authority
closure as a pre-user-byte blocker.

## Current gate position

| Gate | License-free result | Current decision | Smallest remaining work |
| --- | --- | --- | --- |
| A: JCS/JWS | RFC 8785 JSON failed because the Swift/Foundation number representation differed and no maintained strict Swift JCS path passed. ES256 itself interoperated. | **Fail the original format.** | Do not revive generic JCS/JWS without new evidence. |
| A2: deterministic CBOR/COSE | A second hardening run retained 90 cases across `ApprovalGrant` and a mutually exclusive transcript shape. Go, Swift, and TypeScript agreed on all 4 accepts and 86 rejects, and every producer verified through every verifier. | **Conditional pass strengthened; no format pivot.** | Freeze object-specific CDDL/bounds, independently review wrappers, run sustained fuzzing, and prove exact-byte retention in product registration/rendering. |
| B: macOS authority separation | Exact signed XPC, installed per-user LaunchAgents, current UID/session checks, service crash reactivation, stale client/service denial, disjoint Keychain groups, Secure Enclave keys, interactive approval, and protected stores passed. Fresh groups/keys per identity-changing security epoch passed 14 modeled and 8 provisioned process-kill checkpoints. | **Conditional pass strengthened; security-epoch groups are the preferred key design.** | Complete Developer ID package/update/session matrices, remaining notarizations, locked-Keychain/restore faults, and signed target-key authorization. |
| C: execution backend | Stock Apple Container and direct Containerization failed lifecycle authority. libkrun/HVF made one signed VMM process the VM; readiness tracks observed fixed raw scratch, bounded console, wall/cancel handling, exact forced teardown, installed reparent/reap, adversarial cases, and build controls. Live backing-file mutation and an unexpected `NullFs` virtiofs device block the exact profile; typed completion/distribution are incomplete and current bytes are no-go. Stock Bun also lacks an evidenced no-subprocess/no-FFI profile. | **Containerization fails; libkrun/HVF remains the lead native candidate under evaluation but is not ready to freeze or admit. gVisor remains independent.** | Before a user-byte adapter, close runtime authority, immutable root custody, `NullFs`, typed port transport/completion with inline JSON, and complete bundle admission. Require disposable output parsing later before file artifacts; then run composition, installed-host, durability, hostile-workload, supply-chain, soak, and exact gVisor campaigns. |
| D: content custody | The original snapshot/descriptor cases passed, followed by a real multi-process SQLite ledger with atomic one-use races, crash/restart reconciliation, bounded pipes, exact idempotent commits, quarantine, tombstones, and guarded GC. | **Strong conditional pass.** | Compose the ledger with distribution-signed XPC/protected storage and the selected backend; add ENOSPC/I/O/corruption/migration/power-loss and signed cross-store saga tests. |
| E: Supervisor topology | Unprivileged per-user service activation and exact authenticated XPC work without a host-root helper. The direct Swift Containerization controller remains useful for development, but Gate C invalidated it as the production backend authority. | **Per-user/no-host-root topology retained; production backend language split reopened.** | Define the narrow local macOS Supervisor versus Linux OCI/gVisor worker boundary after the gVisor run. Keep Go for portable lifecycle/policy code and Swift for necessary macOS APIs. |
| F: trust transition/recovery | The original state model/process-kill corpus passed. A durability follow-up added writer/CAS races, simulated disk-full, WAL/checkpoint damage, corruption, partial/coherent restores, clock failure, atomic replacement, and installer/backend effect crashes; 17 tests including 18 exact-PID kills passed fail-closed. | **Conditional pass strengthened for process/storage ordering.** | Port to the real Supervisor store; add real APFS capacity/I/O, installed-package, Keychain anchor, backend, migration, and VM power-cut campaigns. Coherent rollback still needs an independent anchor. |

## Important observations

### Containerization failed; libkrun/HVF reopened a distinct native path

The original stock CLI/API remains unsuitable: a restarted API server reported a still-running
helper as stopped. The lower-level direct library behaved materially better. Controller death
removed its Virtualization helper in the original single case and in a simultaneous two-controller
control, while an unrelated baseline helper remained alive.

The focused identity/recovery experiment found no public host-side VM/helper identity, enumeration,
reopen, reconnect, or force-reap surface. `VZGenericMachineIdentifier` identifies guest virtual
hardware only. Concurrent unrelated helper churn then proved that process-name/PID-delta logic
cannot safely attribute one exact guest. Gate C therefore fails the direct Apple path for
production while retaining it for explicit development posture.

The missing PID cgroup was a small API-surface gap, not a missing guest mechanism. The retained
patch maps one optional configuration field to Containerization's existing `LinuxPids` OCI value.
It remains useful development evidence but cannot solve the missing lifecycle identity.

The alternate OCI harness could not run gVisor because `runsc` was not installed or registered and
the spike deliberately made no daemon changes. Under the existing runc runtime, the same harness
observed durable container IDs/labels across client death plus network-none, non-root/no-new-
privileges, empty capabilities, read-only root, cgroup memory/PID/CPU limits, exact tmpfs capacity,
bounded retained logs, and forced cancellation. Those are control-plane results, not gVisor
isolation evidence.

The libkrun/HVF follow-up does not reuse Containerization's lifecycle model. Hypervisor.framework
runs one VM per process, and the observed signed runner was itself the VMM with no new helper. A
private inherited pipe prevented VM start until the controller had verified and durably recorded
PID, start time, path, Team/signing identifier, and CDHash. `SIGKILL` before that record or before
authorization left no running guest; `SIGKILL` after authorization left a recorded, reparented,
code-valid runner that recovery killed exactly.

The same profile ran inside App Sandbox with only one exact read-only disk exception, while the
no-disk-authority control failed. The guest had no usable network or vsock, a guest-read-only raw root,
UID/GID 65534, no groups/capabilities, `no_new_privs`, and active file/process/fd limits. This is a
conditional candidate result, not a validation record. The follow-up tracks closed several
mechanical questions but proved that guest-read-only is not immutable host custody, runner exit is
not workload success, and the exact block-root surface includes `NullFs`. The full comparison and
P0/P1 campaigns are in the
[Gate C implementation-readiness synthesis](GATE_C_READINESS_CHECKPOINT.md) and its
[P0 reconciliation](GATE_C_P0_RECONCILIATION.md).

### The authority and custody design composes locally

The live XPC result closes an important feasibility uncertainty without overstating ad-hoc signing.
The OS applied an exact code requirement on the listener before message delivery. Accepted
messages were independently tied to their audit-token sender through `SecCodeCreateWithXPCMessage`,
then carried an already-open read-only descriptor. The stale client did not reach the protocol;
the exact authenticated client with a malformed operation reached it and was rejected there.

This supports the planned split:

```text
code requirement + message-derived sender  -> component authority
typed installation/epoch/attempt binding   -> protocol authority
one-use Broker ledger                      -> content authority
read-only FD copied into attempt storage   -> byte custody
```

The Apple-credentialed follow-up proved the same composition with Team/channel/role identities,
installed per-user services, provisioned component groups, and protected stores. It also showed
that exact IPC identity does not make a stable Keychain group exact-build or epoch scoped. The
security-epoch transition follow-up then proved disjoint new groups/keys, create-if-absent
fingerprints, old-key retirement, rollback/forward-repair boundaries, and fail-closed process-death
recovery.

The content path now also has a real multi-process ledger. Exactly one process won each repeated
redemption race; consumed authority never resurrected; ambiguous output became quarantined; and GC
waited for terminal attempt state plus the retention horizon.

### The update model now survived real process death

The first Gate F model raised an exception after closing SQLite, which tested ordering but not an
unclean process. The retained child harness now pauses after the committed checkpoint and is killed
with uncatchable `SIGKILL`; neither SQLite connection is closed. A fresh process opens the WAL-backed
stores and reconciles the result. Every incomplete update remained execution-disabled, the final
re-enable checkpoint recovered stable, consumed grants stayed consumed, backend-create ambiguity
kept cleanup responsibility, and completed external release was observed rather than rolled back.

The durability follow-up additionally covered concurrent writers, deterministic SQLite-full,
WAL/checkpoint interruption and damage, corruption, partial restore, clock rollback/unavailability,
atomic file replacement, and fake installer/backend effects. This remains process/filesystem-
ordering evidence rather than an actual power-loss claim. A coherent restore of every local file
still passed, confirming the need for an independently protected checkpoint when rollback detection
is required.

## Apple-credentialed follow-up status

Completed locally:

- Apple Development and Developer ID Team/channel/role/exact-build XPC matrices;
- development-provisioned Broker/Supervisor Keychain groups and daemon/sibling denial;
- persistent Secure Enclave evidence and user-presence Approval keys;
- distinct signed sandbox containers and cross-role denial; and
- an explicit stale same-team build attack confirming stable-group replacement-key use;
- installed per-user launchd/XPC activation, stale replacement denial, and crash reactivation;
- a security-epoch group/key transition with modeled and provisioned crash recovery; and
- three Developer ID exports uploaded through an authenticated notarization profile.

Still open:

- final Apple notarization acceptance, stapling, and Gatekeeper assessment;
- complete runner/Supervisor/helper/LaunchAgent packaging and clean-host distribution (the latest
  runner-only submission remained `In Progress` at the evidence cutoff);
- installed Developer ID update/replay/logout/login and in-flight retry cases;
- protected-container override, Full Disk Access, fast-user-switching, MDM, and migration;
- installer/team/profile/container/key transition behavior; and
- locked-Keychain, real disk/I/O/APFS restore, and power-cut cases.

## Next implementation slice

1. Expand the proposed CBOR/COSE ADR/CDDL from `ApprovalGrant` to every registered or signed v0
   object, preserving exact received payload bytes. The 90-case two-object hardening corpus is now
   the seed for independent review, sustained fuzzing, and runtime/CDDL conformance.
2. Implement the Go authoritative state machines and ports around fake Broker/backend adapters,
   including the Gate F transition fence, component-acceptance barrier, cleanup obligations, and
   idempotent release intent.
   **Started:** the no-guest `SupervisorCore` covers exact registration, approval binding and
   atomic in-memory one-use consumption, transition fencing/acceptance, create-intent cleanup, and
   authoritative-absence/ambiguous-outcome reconciliation. Durable storage, release intent,
   production crypto/IPC, and real process-crash testing remain.
3. Before a real libkrun adapter handles user bytes, close stock-Bun runtime authority, immutable
   runtime-root custody, the independent `NullFs` accept/remove decision, typed port transport and
   completion with bounded inline JSON, and the complete installed-bundle admission checklist.
   Defer the bounded disposable filesystem-image parser only until file artifacts. Retain the
   observed durable-record-before-start, console, forced-teardown, raw-scratch, and installed-
   recovery mechanics as spike inputs only.
4. Run the retained OCI harness with a checksum-pinned `runsc`/shim on a disposable Linux worker,
   including engine/containerd/outer-VM death and Bun compatibility, then compare exact profiles.
5. Port the proven custody ledger and Gate F ordering into durable product stores behind fake
   adapters before composing distribution-signed XPC and a real backend.

## Evidence index

- Gate A: [`../experiments/gate-a-signing-canonicalization/README.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-a-signing-canonicalization/README.md)
- Gate A2: [`../experiments/gate-a2-cbor-cose/README.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-a2-cbor-cose/README.md)
- Gate A2 hardening: [`../experiments/gate-a2-profile-hardening/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-a2-profile-hardening/RESULTS.md)
- Gate B: [`../experiments/macos-authority-separation/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/macos-authority-separation/RESULTS.md)
- Gate B installed services: [`../experiments/gate-b-installed-services/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-b-installed-services/RESULTS.md)
- Gate B key transition: [`../experiments/gate-b-key-rotation/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-b-key-rotation/RESULTS.md)
- Gate C stock API: [`../experiments/apple-container-gate-c/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/apple-container-gate-c/RESULTS.md)
- Gate C direct backend: [`../experiments/apple-containerization-direct/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/apple-containerization-direct/RESULTS.md)
- Gate C identity/recovery: [`../experiments/gate-c-identity-recovery/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-identity-recovery/RESULTS.md)
- Gate C OCI/gVisor contingency: [`../experiments/gate-c-gvisor-contingency/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-gvisor-contingency/RESULTS.md)
- Gate C libkrun/HVF follow-up: [`../experiments/gate-c-libkrun-hvf/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-hvf/RESULTS.md)
- Gate C libkrun storage/egress: [`../experiments/gate-c-libkrun-storage-egress/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-storage-egress/RESULTS.md)
- Gate C libkrun console/lifecycle: [`../experiments/gate-c-libkrun-console-lifecycle/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-console-lifecycle/RESULTS.md)
- Gate C libkrun installed recovery: [`../experiments/gate-c-libkrun-installed-recovery/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-installed-recovery/RESULTS.md)
- Gate C libkrun adversarial: [`../experiments/gate-c-libkrun-adversarial/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-adversarial/RESULTS.md)
- Gate C libkrun supply chain: [`../experiments/gate-c-libkrun-supply-chain/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-supply-chain/RESULTS.md)
- Gate C P0 review/research reconciliation: [`GATE_C_P0_RECONCILIATION.md`](GATE_C_P0_RECONCILIATION.md)
- Gate D: [`../experiments/gate-d-content-custody/README.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-d-content-custody/README.md)
- Gate D ledger: [`../experiments/gate-d-custody-ledger/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-d-custody-ledger/RESULTS.md)
- Gate E: [`../experiments/gate-e-supervisor-topology/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-e-supervisor-topology/RESULTS.md)
- Gate F: [`../experiments/gate-f-trust-transition/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-f-trust-transition/RESULTS.md)
- Gate F durability: [`../experiments/gate-f-durability/RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-f-durability/RESULTS.md)
