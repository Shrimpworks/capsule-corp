# Handoff: Independent review of Capsule Gate C P0 spike program

## Session Metadata

- Created: 2026-07-31 22:52:41 America/Toronto
- Project: `/Users/dsteele/repos/capsule-corp`
- Branch: `codex/license-free-spikes`
- Draft PR: https://github.com/dills122/capsule-corp/pull/8
- Review baseline: `7cb974771463720743acf5e05a47bb4d36b2fcc8`
- Session duration: approximately one hour of synthesis and read-only source inspection

### Recent Commits

- `7cb9747 Record feasibility evidence and Gate C readiness`
- `1f9f55b Expand license-free security spike evidence`
- `01c5506 Preserve Gate F recovery spike`
- `ebb963d Preserve Gate E Supervisor topology spike`
- `b632184 Preserve Gate D content custody spike`

## Handoff Chain

- Continues from: none
- Supersedes: none

This brief is for an independent reviewer. It is not authorization to implement, patch libkrun,
run a VMM, change external state, or promote a security claim.

## Current State Summary

Capsule completed five libkrun/HVF implementation-readiness spikes and integrated their evidence.
libkrun remains the preferred native Apple candidate, but no exact development profile is admitted.
The project now proposes five P0 campaigns before libkrun handles user-derived bytes: immutable
block custody, `NullFs` disposition, typed completion, disposable output parsing, and an admissible
installed runtime bundle. The next step is an independent adversarial review of that P0 program,
especially whether its assumptions and pass criteria are capable of proving the claimed properties.
Do not treat the plan below or the prior single-model findings as correct merely because they are
documented.

## Reviewer Mandate

Work read-only. Attempt to disprove the plan rather than improve its presentation.

1. Read `AGENTS.md` and the authoritative architecture/security documents before judging the plan.
2. Inspect the retained experiment results and exact pinned source locations yourself.
3. Separate observed evidence, inference, design proposal, and unsupported claim.
4. Identify missing attacker capabilities, false assumptions, circular dependencies, unsafe
   authority transfer, ambiguous pass conditions, and tests that could pass without proving the
   property.
5. Do not implement fixes. Return a review and a revised evidence plan only.
6. If consulting external sources, prefer official Apple documentation/SDK headers and pinned or
   upstream libkrun/Imago source. Cite exact sources and distinguish current upstream behavior from
   the pinned 1.19.4 candidate.

### Required Review Output

Return these sections:

1. **Executive verdict:** `viable`, `conditionally viable`, or `no-go`, with confidence and the
   narrow reason.
2. **Incorrect or overstated premises:** each with severity, exact file/source evidence, and why it
   matters.
3. **Missing threat cases:** attacker, precondition, action, affected boundary, and expected safe
   result.
4. **Dependency/order defects:** work that is premature, circular, or missing.
5. **Revised P0 tasks:** each with a falsifiable claim, exact pass, exact fail/inconclusive, retained
   evidence, and safe fallback/pivot.
6. **P0-1 recommendation:** compare the existing pathname API, process-local `/dev/fd` over an
   immediately unlinked object, a narrow FD-native libkrun/Imago API, different-owner/component
   storage, and abandoning the native candidate. Do not select an option without evidence.
7. **Questions requiring official-source research or user policy decisions.**
8. **Prohibited claims:** wording the project must not use on current evidence.

## Codebase Understanding

## Architecture Overview

- The agent-facing daemon proposes work but cannot approve, own user-only content, or launch a
  backend.
- The Approval Broker renders Supervisor-registered typed plan data and issues one-use,
  attempt-bound approval.
- Only the Execution Supervisor authorizes and owns hostile guest creation, termination,
  destruction, and reconciliation. A helper may perform only a sealed, enrolled operation.
- User content becomes an immutable content-addressed snapshot before execution; live host paths
  never enter guest or execute-time authority.
- Backends enforce exact approved controls or refuse. A library capability report never promotes
  posture by itself.
- Spike code is non-production and must remain isolated from product packages.

## Critical Files

| File | Purpose | Relevance |
| --- | --- | --- |
| `AGENTS.md` | Repository security and verification rules | Authoritative constraints; do not weaken |
| `docs/PROJECT.md` | Product boundary and milestones | Defines the first functional and validated-local milestones |
| `docs/ARCHITECTURE.md` | Component authority and backend model | Defines Supervisor-only lifecycle ownership |
| `docs/TECHNICAL_DESIGN.md` | Contract, limits, runtime, and evidence design | Establishes exact-or-refused and evidence vocabulary |
| `docs/security/THREAT_MODEL.md` | Attacker capabilities and trust boundaries | Use to challenge same-UID and guest assumptions |
| `docs/GATE_C_READINESS_CHECKPOINT.md` | Integrated Gate C synthesis and current P0/P1 list | Primary artifact under review |
| `docs/adr/0022-evaluate-libkrun-hvf-native-backend.md` | Candidate selection and conditions | Candidate decision, not profile admission |
| `docs/security/CONTROL_EVIDENCE_MATRIX.md` | Claim registry | Shows DATA-002, COMPLETE-001, EGR-002, SUPPLY-001 gaps |
| `experiments/gate-c-libkrun-storage-egress/RESULTS.md` | Storage/egress observations | Contains the live same-user mutation counterexample |
| `experiments/gate-c-libkrun-adversarial/RESULTS.md` | Local adversarial VMM observations | Contains `NullFs` and runner-zero findings |
| `experiments/gate-c-libkrun-console-lifecycle/RESULTS.md` | Console/resource/teardown observations | Establishes narrow resource and forced-kill evidence |
| `experiments/gate-c-libkrun-installed-recovery/RESULTS.md` | Installed runner/recovery evidence | Establishes same-host mechanics and distribution gaps |
| `experiments/gate-c-libkrun-supply-chain/RESULTS.md` | Build/release admission findings | Current bytes are no-go |
| `experiments/gate-c-libkrun-supply-chain/evidence/admission-checklist.json` | Machine-readable admission gaps | Every required row must be honestly dispositioned |
| `experiments/gate-c-libkrun-adversarial/internal/preflight/preflight.go` | Exact-path preflight experiment | Opens and hashes an FD, then documents the later path gap |
| `experiments/gate-c-libkrun-storage-egress/src/runner.c` | Four-disk experiment runner | Calls pathname-only `krun_add_disk` |

### Exact Local Source Observations to Verify Independently

The local pinned source was present at `/private/tmp/capsule-libkrun-v1.19.4` during this handoff.
It is temporary state; verify the commit and patches before relying on it.

- `src/libkrun/src/lib.rs:762-799`: `krun_add_disk` stores a string path in
  `BlockDeviceConfig`; it does not retain the caller's already-open FD.
- `src/libkrun/src/lib.rs:2327-2418`: block-root remount constructs an `FsDeviceConfig` with
  `shared_dir: None`, built-in init/pivot entries, and a `NullFs`-backed virtiofs device.
- `src/libkrun/src/lib.rs:2881-2885`: block devices are constructed later during start.
- `src/devices/src/virtio/block/device.rs:235-263`: block construction opens the stored pathname
  for device identity, then asks Imago to open the pathname again. Stable preflight FD identity does
  not automatically carry through these opens.
- `src/vmm/src/vmm_config/block.rs:28-75`: block configuration is cloneable and path-based.

### Key Patterns Discovered

- Runner exit zero is lifecycle evidence only. Corrupt roots, kernel memory failures, a missing
  executable, and intentional guest failure produced ambiguous zero statuses.
- Guest read-only flags do not make the host backing object immutable.
- `shared_dir: None` narrows the `NullFs` finding but does not make the device absent.
- App Sandbox narrows enrolled component authority; it does not make an unsandboxed same-UID
  attacker unable to mutate an ordinary same-owner pathname.
- Generated `.build/` and `.runs/` evidence is ignored. Tracked results must retain hashes,
  environment, observations, limitations, and reproducible commands.

## Work Completed

### Tasks Finished

- [x] Integrated all five Gate C readiness experiment packages into the main workspace.
- [x] Reconciled ADR-0022, roadmap, architecture, threat model, and evidence matrix.
- [x] Opened draft PR 8 and verified the complete repository test/lint/build suite.
- [x] Mapped the initial P0 dependency order and inspected the exact path-open/`NullFs` source.
- [x] Ran one fresh-context single-model adversarial review of the proposed P0 plan.

## Files Modified

No tracked file was modified after commit `7cb9747`. This handoff is an untracked file under
`.claude/handoffs/`; do not add it to the branch unless the user explicitly wants it in the draft PR.

## Decisions Made

| Decision | Options considered | Rationale |
| --- | --- | --- |
| Independently review before implementation | Immediately build all P0 spikes; review the dependency/evidence plan first | These are security claims; a poorly specified spike can pass without proving its property |
| Start eventual implementation with custody | NullFs, completion, parser, or release first | Unsafe custody invalidates every real user-derived disk attachment downstream |
| Package last | Begin notarization now; bind the exact outputs of P0-1 through P0-4 | Earlier release bytes are already known to change and currently fail admission |
| Preserve hard negative results | Clamp/reinterpret failures; record conditional fail/no-go | Spike success means a trustworthy decision, including a backend pivot |

## Proposed P0 Review Artifact

### Dependency Order

1. Immutable custody precedes real user-derived block attachment.
2. Custody and `NullFs` should share one governed libkrun fork/upstream decision.
3. Typed completion defines the terminal evidence the output gate consumes.
4. The parser mechanism consumes the exact custody and completion contracts.
5. The signed/notarized bundle binds the exact outputs of the first four campaigns.

### P0-1: Immutable Block Custody

Claim to test: within the declared same-UID, no-root, no-special-entitlement attacker model, a
Supervisor-issued raw object can be attached so no attacker-writable descriptor, mapping, alias, or
backing-store path can change bytes observed by the guest during the attempt.

Candidate mechanisms:

- current pathname API as the known negative control;
- process-local `/dev/fd` referring to an immediately unlinked object;
- a narrow FD-native libkrun/Imago experiment API;
- different-owner or OS-enforced component storage, with its privilege cost made explicit; and
- stopping the native backend candidate.

Required corpus: creation/finalization, read-only handoff, descriptor/mapping inventory, exact
inode/size/digest, symlink/rename/replacement/hard-link/truncate/write attempts, pre-existing writable
alias, inherited descriptor leakage, controller/runner death, cleanup, and the full VMM live-read
case. Passing identity without mutation denial is failure. A modified library proves feasibility
only; it does not admit the shipped profile.

### P0-2: `NullFs` Disposition

Compare:

- explicit acceptance of the exact `shared_dir: None` device, with the host-backed share API absent
  from the closed runner and negative guest/configuration/operation probes; and
- a fixed initramfs/firmware path that supplies init and pivot mount points without constructing a
  virtiofs device.

Rerun source audit, binary import/feature audit, device and mount inventory, boot, malformed
configuration, guest interaction, concurrency, and teardown. Pass only with a closed device profile
and retained fail-closed criteria; otherwise use the other branch or stop the candidate.

### P0-3: Attempt-Bound Completion

Test a fixed-size, commit-last record on a separate bounded raw completion device controlled by the
trusted runtime launcher—not workload stdout and not VMM exit. Define the runtime/kernel/workload
trust boundary explicitly. Bind version, unpredictable attempt capability, registration/plan,
profile/runtime, workload outcome, input/unmount state, and output identity. A checksum detects
tearing but is not authentication; specify how workload forgery is prevented and what a compromised
guest kernel can still fake.

Required cases: success, workload nonzero, missing executable, corrupt root, panic/OOM,
timeout/forced kill, output flood, child-forged stdout, wrong attempt/profile/runtime, replay,
duplicate/partial/torn record, mutation after record, missing flush/unmount, and runner-zero
ambiguity. Output sealing and input unmount must precede the final record commit. Host validation
occurs only after exact runner stop.

### P0-4: Disposable Output Parser

Predeclare comparison criteria for a native App-Sandboxed parser receiving exact handles versus a
disposable no-network parser VM: exact input authority, ambient filesystem/network/process denial,
resource enforcement, lifecycle identity/recovery, compromise consequence, bounded output channel,
and teardown. The fallback for either design is release nothing.

For the selected experiment, accept only a closed raw/ext4 profile and fixed
the fixed guest slot "result/data.json". Exercise unknown ext4 features, filesystem errors, undeclared entries,
links and special inodes, sparse/overlap/out-of-range extents, xattrs/ACLs, owners/modes/times,
inode/entry/logical/allocated/stream limits, and mutated superblocks/journals/extents/directories/
checksums/feature flags. Inject parser crash, timeout, OOM, handle substitution, controller death,
partial copy/fsync/rename, storage pressure, and cleanup ambiguity. Do not claim power-loss atomicity
before the separate APFS campaign.

Also challenge whether general ext4 artifact parsing is truly P0 for the first inline JSON-only
development milestone, or whether a smaller bounded result channel can safely defer it to the
regular-file artifact slice.

### P0-5: Admissible Installed Bundle

Only after P0-1 through P0-4, pin and close the sysroot, Cargo/Rust/Xcode/LLD, libkrun,
libkrunfw/kernel, Bun, launcher/parser, patches, configs, entitlements, deployment target, licenses,
and build order. Produce complete SBOM, corresponding-source bundle, provenance, and independent
builder comparison. Build/sign/notarize/staple/read back the complete per-user topology, not only
the runner app. Exercise update, disable, explicit rollback, partial install, and repair-required
ordering.

Clean-host quarantine, Gatekeeper, first-launch, minimum-OS, newer-OS, session, recovery, uninstall,
and repair cases require coordinated external hosts. No automated logout or reboot is authorized.
Passing the admission checklist permits consideration of one development profile only.

## Prior Single-Model Findings to Challenge

These are inputs to review, not accepted conclusions:

1. Exact inode/object identity does not prove immutability if a writable FD, mapping, alias, or
   backing-store path remains.
2. An experiment-only libkrun patch proves feasibility, not the shipped profile.
3. Launcher, runner, and parser authority boundaries were underspecified relative to
   Supervisor-only lifecycle ownership.
4. Registration-ID-only execution and rejection of execute-time plan/backend/image/path overrides
   need explicit negative evidence.
5. A nonce/checksum completion record needs a defined observation boundary, workload-forgery
   defense, commit-last ordering, and torn-write semantics.
6. `shared_dir: None` plus no observed mount is insufficient to accept `NullFs` without explicit
   guest-reachable negative tests.
7. Parser comparison needs predeclared isolation criteria and must not overclaim power-loss atomicity.
8. Every task needs exact pass, fail/inconclusive, retained evidence, and fallback thresholds.

## Pending Work

## Immediate Next Steps

1. Perform the independent read-only review using the required output format above.
2. Return the review to the user; do not edit the repository or begin P0 implementation.
3. Reconcile each finding as contract misread, valid/actionable, valid trade-off, or noise before
   approving an implementation plan.

### Blockers/Open Questions

- [ ] What exact same-UID attacker powers are in scope: pre-existing writable descriptors/mappings,
  unsandboxed filesystem access, process inspection, debugging, Full Disk Access, or file-ID APIs?
- [ ] Is the fixed guest kernel/launcher trusted against only an unprivileged hostile workload, or
  must completion remain meaningful after guest-kernel compromise?
- [ ] Is the no-root/no-separate-UID topology an absolute v0 constraint or a preference that may be
  reopened by an ADR if immutable custody otherwise fails?
- [ ] Is ext4 artifact extraction genuinely required before the first inline JSON-only libkrun
  development slice?
- [ ] Which clean macOS 14 Apple-silicon and representative newer hosts are available for the final
  package-floor campaign?

### Deferred Items

- Product implementation, libkrun patches, VMM runs, and release builds are deferred until review.
- Reboot, logout/login, clean-host distribution, and power-loss work requires explicit coordination.
- P1 composition, real APFS, hostile Bun, real gVisor, release drills, and soak follow P0 decisions.

## Context for Resuming Agent

## Important Context

The user explicitly requested an independent review. Do not optimize for agreement with the existing
synthesis. A negative decision is useful. The most dangerous failure mode is a test that reports
green while only proving stable metadata, ordinary process behavior, or absence of an observed
mount—not immutable bytes, a closed VMM surface, trustworthy completion, or safe release.

The repository builds a security boundary. Do not introduce a daemon-to-runner/helper path, accept
execute-time plan/image/path bytes, add live host paths, silently approximate limits, parse rich
formats in the daemon/Supervisor, or treat experiment code as product code.

## Assumptions Made

- The review baseline and working tree are clean; verify with Git before relying on this.
- The local `/private/tmp` libkrun tree matches pinned 1.19.4 plus the two recorded Capsule patches;
  independently verify it because temporary state can change.
- The relevant attacker is same UID without root, special entitlements, or authorization to debug
  enrolled hardened components. Whether that is complete is itself a review question.
- No current result promotes libkrun beyond candidate/development research posture.

## Potential Gotchas

- `/dev/fd/N` may reopen or duplicate a descriptor with semantics that differ across macOS/Linux;
  do not infer behavior from another platform.
- Closing the original writer does not help if a duplicate writable descriptor or writable mapping
  remains.
- An attempt nonce is not automatically secret or authentic inside a guest.
- A fixed trusted guest launcher does not create platform attestation; record only what it observed.
- A signed/notarized app is attributable distribution evidence, not safe-logic evidence.
- Notarizing the runner alone does not cover the Supervisor/helper/LaunchAgent/install topology.
- The existing actual-gVisor claim is unvalidated; its retained run used runc only.

## Environment State

### Tools/Services Used

- Git branch `codex/license-free-spikes`, tracking `origin/codex/license-free-spikes`.
- Draft GitHub PR 8 is open.
- Node 22.22.1, pnpm 10.28.2, and the configured Go toolchain passed the integration suite at
  commit `7cb9747`.
- Local source inspection used the temporary pinned libkrun tree described above.

### Active Processes

No VMM, parser, test server, or build process was intentionally left running by this planning task.

### Environment Variables

No environment variables or secret values are required for the independent read-only review.

## Related Resources

- `docs/GATE_C_READINESS_CHECKPOINT.md`
- `docs/FEASIBILITY_SPIKES.md`
- `docs/ROADMAP.md`
- `docs/adr/0022-evaluate-libkrun-hvf-native-backend.md`
- `docs/security/CONTROL_EVIDENCE_MATRIX.md`
- `experiments/gate-c-libkrun-storage-egress/RESULTS.md`
- `experiments/gate-c-libkrun-adversarial/RESULTS.md`
- `experiments/gate-c-libkrun-console-lifecycle/RESULTS.md`
- `experiments/gate-c-libkrun-installed-recovery/RESULTS.md`
- `experiments/gate-c-libkrun-supply-chain/RESULTS.md`

## Security Reminder

Do not include credentials, signing material, private keys, tokens, notarization profiles, or raw
user content in the review. Cite identifiers and hashes already present in tracked evidence only
when necessary.
