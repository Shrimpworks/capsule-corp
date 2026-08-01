# Gate C implementation-readiness synthesis

Date: 2026-07-31
P0 plan last reconciled: 2026-08-01

Status: completed cross-track synthesis. This document records spike decisions and coordination
history. It is not a backend posture promotion, a `BackendValidationRecord`, or permission to run
hostile/user workloads with the retained experiment code.

## Bottom line

All five libkrun/HVF implementation-readiness tracks completed and their retained source,
documentation, and selected evidence are now in this repository. Generated `.build/` and `.runs/`
directories remain ignored and disposable.

The combined decision is:

- continue backend-independent registered-plan, fake-backend, custody-ledger, and contract work;
- retain libkrun/HVF as the lead native Apple candidate under evaluation for one eventual
  development profile;
- do not admit the currently built runtime bytes or freeze a libkrun execution profile yet; and
- do not claim `validated-local`, production readiness, absence of a host share, graceful shutdown,
  exact host CPU/memory quotas, or safe output extraction.

The candidate solved Apple Containerization's hidden-helper lifecycle problem, but the follow-ups
found two profile blockers and three release/composition blockers. Immutable block custody is not
enforced through libkrun's pathname disk API, and the block-root path exposes an unexpected
guest-visible `NullFs` virtiofs device. Typed guest completion, a production-safe output parser,
and admissible signed/notarized runtime bytes were also incomplete. A subsequent independent review
and targeted source research narrowed the first inline slice: source/input and bounded inline JSON
may use dedicated virtio-console ports, making filesystem-image parsing a later artifact gate, while
stock Bun's subprocess/FFI surface adds a P0 runtime-authority blocker. Those are actionable
findings, not a failure of the overall architecture.

## Completed tracks

| Track | Task ID | Decision | Retained evidence |
| --- | --- | --- | --- |
| Storage, scratch/output, and egress | `019fb9e7-25ec-78b2-8bf2-6558a2a9f250` | Conditional pass for development-only raw-block staging; immutable same-user host custody remains failed | [`gate-c-libkrun-storage-egress/RESULTS.md`](../experiments/gate-c-libkrun-storage-egress/RESULTS.md) |
| Console, timeout, cancellation, CPU, and memory | `019fb9e7-25ec-78b2-8bf2-659ad6ccdef1` | Conditional pass for bounded console, wall/cancel scheduling, exact forced teardown, and closed resource profiles | [`gate-c-libkrun-console-lifecycle/RESULTS.md`](../experiments/gate-c-libkrun-console-lifecycle/RESULTS.md) |
| Installed lifecycle and crash recovery | `019fb9e7-25ec-78b2-8bf2-65688a86f41f` | Conditional pass for same-machine installed runner-read and recovery mechanics; distribution and authority closure remain open | [`gate-c-libkrun-installed-recovery/RESULTS.md`](../experiments/gate-c-libkrun-installed-recovery/RESULTS.md) |
| Adversarial VMM and cross-job isolation | replacement `019fba17-9659-74f1-ab6f-82bfb72bc991` | Conditional fail for the exact profile because of the unexpected `NullFs` virtiofs device | [`gate-c-libkrun-adversarial/RESULTS.md`](../experiments/gate-c-libkrun-adversarial/RESULTS.md) |
| Packaging, provenance, patches, and supply chain | `019fb9e7-27b1-7c42-9903-8be99f620602` | Conditional feasibility pass; current runtime bytes are no-go for admission | [`gate-c-libkrun-supply-chain/RESULTS.md`](../experiments/gate-c-libkrun-supply-chain/RESULTS.md) |

## Cross-track findings

### Storage and output

- Separate raw root/source/input devices were read-only to the guest. A fixed 12,582,912-byte
  scratch image reached `ENOSPC`, and the spike extractor rejected oversized, sparse, hard-linked,
  and device outputs.
- A same-user host process changed the source backing image while the guest device was live.
  Post-stop hashing detected the mutation too late to prove the guest executed approved bytes.
- `chmod`, path secrecy, advisory locks, App Sandbox on only the runner, and post-stop hashes do not
  close this race. Product use needs an already-open immutable object, OS-enforced component
  storage with equivalent evidence, or a reviewed libkrun API/fork that accepts the required
  immutable handle semantics.
- The shell/`debugfs`/Docker collector is an experiment oracle, not a product parser. It must be
  replaced by a purpose-built, disposable, bounded parser sandbox before filesystem artifact
  release. It does not block the first bounded inline-JSON result path.

### Console, resource limits, and termination

- Per-stream 4,096-byte stdout/stderr prefixes, explicit truncation, and continuous draining stayed
  bounded under output flood, backpressure, reader stalls, concurrency, and controller recovery.
- Wall timeout and cancellation were scheduled outside the guest. Revalidated exact-process
  `SIGKILL` is the only evidenced non-cooperative teardown path.
- The tested shutdown eventfd did not stop the guest within three seconds. Graceful shutdown is
  bounded best effort, never the sole cleanup mechanism.
- Supported resource vocabulary is narrow: integer vCPU topology and closed RAM profiles. The
  minimal fixture passed at one vCPU/64 MiB; 32 MiB, 48 MiB, and 96 MiB failed in different ways.
  The smallest evidenced Bun profile remains one vCPU/256 MiB.
- CPU percentages, CPU-time ceilings, arbitrary RAM, and exact total host/VMM memory are
  unsupported and must be rejected rather than approximated.

### Completion and lifecycle

- Several malformed roots, guest panics, a missing executable, and an intentional guest crash
  still produced host runner status zero. Runner exit status is lifecycle evidence, not guest
  success evidence.
- Ordinary success requires an exact attempt-bound typed completion record plus expected result,
  input-integrity, applicable bounded validation/parser, and teardown evidence.
- A Developer ID-signed, hardened, App-Sandboxed runner read a sealed bundle-local root without a
  temporary path entitlement. With `AbandonProcessGroup=true`, six runners survived Supervisor
  `SIGKILL`, reparented to PID 1, and were exactly recovered using PID/start/user/path/code identity
  and installed-byte identity. Disabling that launchd property removed the child with the job.
- Corrupt records, live binary replacement, identity-helper failure, a second cooperating
  Supervisor, and durable record-write failure all failed closed in the tested mechanics.
- This is one-host cooperating-process evidence. The notary submission
  `1a67daee-ec4e-4572-ad9a-1a1fa3f63bcf` was still `In Progress`; the tested topology has no accepted
  ticket, staple, Gatekeeper/clean-machine result, logout/login, sleep/wake, reboot, or macOS-floor
  validation. The upload covers only the runner app, not the complete installed topology.

### VMM surface

- The reviewed adversarial report contains 36 VMM cases, four identity cases, 11 findings, one
  limitation, and one failing assertion. Every recorded runner was gone at collection time.
- The failing assertion is real: `krun_set_root_disk_remount` creates a guest-visible virtiofs
  device backed by `NullFs` with `shared_dir: None`. No host directory was configured or mounted,
  but the device and its VMM-side attack surface are present.
- Integration must either accept and independently validate that exact `NullFs` surface or remove
  it in a governed libkrun change and rerun the corpus. Until then the exact profile is not frozen.
- Optional virtiofs, vsock, GPU, sound, duplicate-disk, and missing-disk calls can return success
  before start even when the compiled feature is absent. A closed typed runner surface—not libkrun
  return codes—defines allowed capability.

### Runtime and supply chain

- Unsigned libkrun/libkrunfw outputs matched across two controlled local source directories only
  after locked/offline dependency resolution, path remapping, and an explicit macOS 14 deployment
  target. This is reproducibility feasibility, not an independent two-builder claim.
- The default build was path-dependent. Current runner/firmware metadata declares macOS 26, the
  sysroot and complete kernel/source materials are not fully pinned, the 115-component SBOM is only
  an input inventory, and final signed/notarized-byte identity is incomplete.
- Both Capsule patches require a governed fork and review until they are upstreamed and the exact
  released profile is revalidated. Update, rollback, disable, revocation, and corresponding-source
  workflows are designs, not exercised product mechanisms.
- Pinning Bun 1.3.14 proved that the runtime executed the fixture; it did not prove the documented
  no-subprocess/no-FFI contract. The stock runtime exposes those APIs and no non-bypassable disable
  profile has been evidenced, so the current runtime profile must refuse admission until P0 closes
  or an explicit ADR changes the contract.
- The planning floor remains Apple silicon/macOS 14+ only as a provisional source/platform target.
  It is not a claim about the current package or a validated supported-host range.

## Post-synthesis independent review

An independent read-only adversarial review challenged the P0 premises, threat coverage, and
dependency order. Its viable `/dev/fd/N` direction and conditional-parser finding were accepted;
its use of post-spawn `SCM_RIGHTS` transfer was refined to direct Supervisor-to-runner inheritance
as the smaller first topology. Targeted research then found libkrun 1.19.4's generic multiport
console APIs and the unresolved stock-Bun authority mismatch.

The durable disposition, exact hypotheses, pass/fail branches, first-slice data path, and prohibited
claims are in the [Gate C P0 reconciliation](GATE_C_P0_RECONCILIATION.md). These are research and
planning conclusions, not additions to the completed spike evidence above.

### Final preimplementation source review

A fresh-context issues-only review of the reconciled plan and pinned upstream source completed on
2026-08-01. It preserved the decision to begin backend-independent Phase 2–4 work and found four
additional blockers before a real libkrun adapter receives user bytes:

- pinned virtio-console control handling uses guest-supplied port IDs to index VMM port/queue state
  without an evidenced bounds check, so P0-3 now includes the hostile control/queue/descriptor path;
- the retained experiment launcher replaces itself with the workload, while the planned trusted
  launcher must remain distinct, withhold completion authority, and commit only after exact child-
  tree termination;
- the host VMM runner needs an exact role-specific descriptor allowlist because App Sandbox does not
  revoke ambient inherited descriptors; and
- pinned console TX backpressure/shutdown and partial-error behavior cannot supply the required
  bounded failure semantics without an exact fail-closed corpus or a governed patch.

The same review corrected evidence sequencing: an early installed harness may test topology, but
final admission rebuilds selected mechanisms and reruns affected P0 gates on final signed/notarized
bytes. A P0 record is explicitly `development-admitted`, never `validated-local`. Bun denial now
requires construction-level source/mechanism closure rather than a finite API corpus alone.

Suspected impossibilities that source did not support were rejected: pinned implementation can
construct one-directional ports when the unused FD is negative, raw block data uses positional I/O,
both block opens use the configured `/dev/fd/N` string, and the ext4 parser remains correctly
deferred until file artifacts. Directionality is still undocumented public behavior and therefore
requires a pinned canary or governed API.

## Contract consequences

The following backend-independent vocabulary is ready to implement and test with a fake backend:

- execute only by Supervisor-issued registration ID; never accept replacement plan or backend
  bytes at execute time;
- bind immutable runtime/profile/validation references before approval and resolve mutable aliases
  before registration;
- use a closed capability/admission result; unsupported controls refuse execution;
- resolve exact limits before approval and never clamp or substitute them;
- distinguish configured guest RAM from host/VMM-memory accounting, and vCPU topology from CPU
  quotas;
- distinguish physical scratch-image bytes, artifact count, per-artifact logical bytes, total
  logical bytes, console prefix bytes, bounded port-frame bytes, and parser limits;
- keep guest completion, runner lifecycle, input integrity, output validation, teardown, and
  terminal classification as separate evidence;
- preserve `cleanup-required`, `unresolved`, `integrity-failed`, `unsafe-output`, and quarantine
  states across recovery; and
- require `CreatesGuest() == false` for the fake backend and a hard fence between fake tests and
  any VMM launch path.

The following libkrun-specific values are not frozen: runtime-root descriptor/path mechanics,
console-port framing and guest permissions, arbitrary resource values, a claim that virtiofs is
absent, graceful-only cancellation, runner-zero success, stock-Bun authority restrictions,
temporary App Sandbox path exceptions, the current runtime byte manifest, and a macOS 14 package
claim.

## Remaining spike campaigns

### P0 — before a libkrun development adapter handles user bytes

1. **Runtime-authority closure:** prove that the exact Bun/launcher profile refuses subprocess,
   FFI, native-addon, inspector, macro, environment-file, and package-install paths. If stock Bun
   cannot do so, choose a governed patch/alternate runtime or explicitly revise the contract in a
   new ADR.
2. **Immutable runtime-root custody:** test protected exclusive creation, a distinct genuine
   read-only descriptor, closure of every writable alias/mapping, unlink, final digest through that
   exact descriptor, direct runner inheritance, and `/dev/fd/N` under the exact installed App
   Sandbox profile. Separately pass stable attachment identity, frozen-object construction, and
   adversarial end-to-end custody, including FD reuse/shared-state, same-user pre-custody, pathname,
   hard-link, mapping, debugger/task-port attempt, crash, and recovery cases. A narrow FD-native
   libkrun change is the fallback; failure of both rejects libkrun for v0.
3. **`NullFs` disposition:** independently remove or validate the unexpected device and rerun the
   exact device/cross-job corpus. Accepted residual surface requires the complete guest-reachable
   virtiofs/FUSE/queue/worker path, sanitizer fuzzing, retained coverage/limitations, and zero
   unresolved high-severity findings—not a claim that fuzzing proved the surface bug-free. Do not
   couple this admission decision to custody merely because one fork could touch both mechanisms.
4. **Typed port transport and completion:** use dedicated bounded attempt-bound virtio-console
   ports for source/input and one typed completion frame containing bounded inline JSON. Freeze
   separate source, canonical-input, physical-frame, and JSON-payload caps plus per-channel role,
   version, attempt, plan/profile, length/digest, and completion commit-trailer semantics before
   implementation. Prove descriptor/node isolation, framing without EOF dependence, cap-plus-one
   draining, partial writes, backpressure, reader death/stall, invalid/swapped descriptors, bounded
   shutdown/forced teardown, forgery resistance, launcher failure, stale/duplicate attempts, and
   separate integrity/teardown dispositions without network/vsock. Pinned upstream has known
   unchecked guest port IDs, non-stop-aware transmit waiting, undocumented directional behavior,
   shared-status mutation, and partial-then-error assumptions; patch or reject the transport if the
   exact VMM/control/queue/descriptor and stream corpus cannot fail closed. Freeze a distinct
   launcher that verifies inputs before child start, withholds completion authority, uses a closed
   child manifest, waits for child-tree termination, and commits last. Start the host runner with an
   exact descriptor allowlist and reject every unexpected inherited authority.
5. **Admissible complete development bundle:** pin all build inputs, govern patches, produce the
   complete manifest/SBOM/provenance/source bundle, build the intended minimum-OS bytes, and
   sign/notarize/staple/read back the complete Supervisor/runner/runtime/per-user-service topology
   on clean hosts. Use an early installed harness for topology assumptions, then rebuild and rerun
   affected P0 gates on final bytes after the mechanisms are selected. A runner-only ticket or
   early harness does not pass.

The disposable bounded ext4/raw-image parser is deferred to the regular-file/file-artifact slice;
it remains mandatory before any filesystem artifact is released.

Host-root execution, a separate-owner host service, and a privileged host helper are prohibited for
v0. If no-host-root custody cannot close, a future exception requires a new ADR; the spike may not
silently widen privilege or the supported attacker tier.

### P1 — before `validated-local` or production claims

1. **End-to-end composition and fault injection:** compose registration, approval, custody,
   staging/port transfer, start, completion, applicable parsing, release, cleanup, and receipt;
   kill processes or inject storage failure at every durable/side-effect edge.
2. **Real installed lifecycle matrix:** test sleep/wake, logout/login, fast user switching, locked
   Keychain, reboot recovery with explicit human coordination, update/replacement, quarantine,
   first launch, multiple supported macOS versions, and clean machines.
3. **Real storage/durability pressure:** test APFS ENOSPC, I/O errors, fsync/power interruption,
   corruption, migration/restore, concurrent attempts, and long-run cross-job remanence.
4. **Hostile runtime/workload corpus:** exercise the exact Bun/launcher surface, disabled process/
   FFI/native-addon/inspector/package powers, kernel and memory pressure, malformed inputs, and
   repeated concurrent hostile guests.
5. **Independent OCI/gVisor comparison:** run checksum-pinned `runsc`/shim on a disposable Linux
   worker through engine, containerd, cgroup, network, storage, log, Bun, crash, and outer-VM
   failure cases. The existing runc harness makes no gVisor isolation claim.
6. **Release/update response drill:** use independent builders, complete source/license publication,
   vulnerability review, TUF-backed activation, emergency disable/revocation, partial update,
   rollback to an explicitly supported prior bundle, and repair-required recovery.
7. **Soak and quantitative budgets:** measure startup, thermal/CPU/RSS/disk/I/O pressure, descriptor
   and process leakage, cleanup latency, concurrency, and repeated controller/VMM failure.

### P2 — only if stronger claims require them

- external/non-rollbackable witnessing for coherent local rollback detection;
- Endpoint Security Guardian feasibility without granting approval or content authority;
- Full Disk Access, MDM, restore/migration, and enterprise deployment matrices;
- Intel Mac support through a separately selected backend; and
- rich document/media/archive parsing in separate disposable parser sandboxes.

No finite spike list can prove that every possible pitfall is absent. These campaigns close the
known claim-critical gaps and add composition, mutation, fault-injection, independent-backend, and
long-run testing specifically to discover unknown interactions.

## Coordination failure history

The original coordinator task `019fb58b-04a8-7121-98c9-82d304cf82a5` ended in `systemError` after
creating the five tracks. The original adversarial task
`019fb9e7-2692-7cb0-a5c2-cebb9378e07f` then ended in `systemError` twice. The replacement task was
scoped explicitly to local, bounded VMM configuration validation and completed normally as
`019fba17-9659-74f1-ab6f-82bfb72bc991`.

Those failures were Codex task/app lifecycle failures, not evidence loss and not a product security
finding. Keeping every retained result in the repository, with task IDs and explicit integration
state here, removes chat history as a project dependency.

## Verification record

Each isolated task recorded its own experiment and repository verification. The integrated branch
must also pass the repository-required Node 22/pnpm/Go suite before handoff; the resulting commit
and draft PR are the authoritative integration checkpoint.
