# Gate C P0 reconciliation

Date: 2026-07-31

Status: planning decision after independent adversarial review and targeted source research. This
document refines the remaining work in the
[Gate C readiness synthesis](GATE_C_READINESS_CHECKPOINT.md). It records hypotheses to test; it is
not new backend evidence, a frozen libkrun profile, or permission to execute user bytes.

## Purpose

The five Gate C readiness tracks produced a conditionally viable native libkrun/HVF candidate and
identified five follow-up campaigns. An independent read-only review then challenged their scope,
ordering, threat coverage, and claim wording. Targeted research checked the pinned libkrun and Bun
surfaces and the likely no-root macOS packaging shape.

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
| Does stock Bun satisfy the documented no-subprocess/no-FFI contract? | Not evidenced. Bun 1.3.14 exposes subprocess and FFI APIs and has no demonstrated non-bypassable `--no-spawn`/`--no-ffi` profile. Preserve the contract and make runtime-authority closure a P0 gate; change it only through an explicit ADR. | Current stock-Bun profile is unsupported for this claim. |
| Does no-root remain plausible? | Yes, provisionally: test one complete signed/notarized app containing the enrolled per-user Supervisor/runner/runtime topology and an embedded `SMAppService` registration. Do not add privilege silently if it fails. | Packaging hypothesis; clean-host and minimum-OS evidence remain open. |

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
one typed completion frame containing bounded inline JSON
        │
        ▼
host validation → teardown evidence → terminal transcript → Broker-held result
```

The pinned guest kernel and trusted launcher are part of this development profile's TCB. The
completion frame is not attestation against a compromised guest kernel. Host VMM containment must
still treat guest behavior as hostile.

## P0 admission program

### P0-0: runtime-authority closure

Goal: preserve the documented dependency-free Bun profile with no subprocess, FFI, native addon,
inspector, macro, environment-file, or package-install authority.

Required evidence includes direct APIs, aliases, workers, child-process creation, native/FFI
loading, inspector activation, dynamic import/package behavior, environment/config discovery, and
launcher/descriptor inheritance under the exact runtime build. Runtime flags are defense in depth,
not a boundary unless bypass attempts demonstrate otherwise.

Pass: an exact pinned mechanism refuses every prohibited power. Fail: choose a governed runtime
patch or alternate runtime, or explicitly revise the product contract in a new ADR. Do not keep
building around a stock Bun profile that contradicts the plan shown to the user.

### P0-1: immutable runtime-root custody

Goal: prove that a concurrent same-user attacker cannot change or substitute bytes observed through
libkrun's pathname disk API.

First candidate:

1. Create an unguessable file with exclusive creation inside Supervisor-protected storage.
2. Finish and verify the exact bytes using the creator descriptor.
3. Open a distinct genuine read-only descriptor and verify mode, device, inode, type, link count,
   size, and digest.
4. Close and unmap every writable handle, unlink the pathname, and retain only the read-only
   descriptor.
5. Inherit that descriptor directly into the signed runner and attach `/dev/fd/N`.
6. Revalidate descriptor mode/identity before start and test readback from the guest.

The corpus includes pre-creation name substitution, the create-to-read-only window, hard links,
writable mappings, path replacement, inherited writable aliases, runner/Supervisor crashes,
debugger or Mach task-port attachment, explicit container grants, and installed App Sandbox
behavior. `O_CREAT|O_EXCL`, `chmod`, secrecy, read-only guest flags, or post-stop hashing alone do
not pass.

Fallback: one narrowly governed FD-native libkrun API change followed by the same corpus. If both
paths fail, reopen the native/no-root backend decision rather than expanding privilege by default.

### P0-2: `NullFs` disposition

Goal: eliminate the unexpected device or accept only a bounded, understood VMM surface.

Pass by either removing the device and rerunning the full device/cross-job corpus, or by defining
the exact exposed protocol, independently reviewing and fuzzing malformed guest messages, and
proving that no product input can select a host-backed share. Custody and `NullFs` use independent
decisions even if one libkrun patch eventually changes both areas.

### P0-3: typed port transport and completion

Goal: move only exact bounded source/input bytes into the guest and receive exactly one
attempt-bound completion/result without ambient network or vsock authority.

The corpus covers wrong/duplicate/stale attempts, partial and oversized frames, malformed lengths
and JSON, output floods, host-reader stalls/death, launcher crash at every boundary, corrupt root,
missing executable, guest panic/OOM, descriptor leakage, child inheritance, workload direct-open,
subprocess/FFI forgery attempts, and terminal classification when teardown or input integrity is
missing. The unprivileged workload must not own the completion descriptor or port node.

Pass: ordinary success requires a valid frame plus separate input-integrity, bounded-result,
runtime-integrity, runner-lifecycle, and teardown dispositions. Runner exit zero never substitutes
for a missing frame.

### P0-4: complete installed development bundle

Goal: admit the exact bytes and topology that will run the first development slice without root.

Build and exercise one complete app containing the Supervisor, runner, libkrun/libkrunfw,
firmware/kernel/root, Bun, trusted launcher, entitlements, manifests, source/license material, and
embedded per-user service registration. Prepare the signing/notarization pipeline early, but only
the final topology counts for admission.

Pass requires pinned materials, governed patches, manifest/SBOM/provenance/source completeness,
minimum-OS builds, signing and notarization, staple/Gatekeeper assessment, installed-byte readback,
clean-host launch, inherited-descriptor behavior, exact process identity, recovery, and an explicit
macOS support floor. A runner-only ticket is not sufficient.

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

## Primary research references

- [libkrun 1.19.4 public header](https://raw.githubusercontent.com/libkrun/libkrun/v1.19.4/include/libkrun.h)
- [Bun FFI documentation](https://bun.sh/docs/runtime/ffi)
- [Bun child-process documentation](https://bun.sh/docs/runtime/child-process)
- [Apple protected app-container guidance](https://developer.apple.com/documentation/xcode/protecting-local-app-data-using-containers)
- [Apple `SMAppService`](https://developer.apple.com/documentation/servicemanagement/smappservice)
