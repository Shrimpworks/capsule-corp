# Phase 2B and Gate C current maintainer checkpoint

Date: 2026-08-04

Status: repository checkpoint through current main plus the passive V0 contract, bounded V1
artifact, blocked V2 process-profile evidence, and the passed supported-profile replacement design.
This is a status and dependency
index. It does not accept a Proposed ADR, activate a consumer or endpoint, admit a runtime or
backend, authorize user bytes, or authorize a guest.

## How to read the status

Work progress in this checkpoint follows the
[canonical status language](STATUS_LANGUAGE.md). A scoped implementation or experiment may be
`PASSED` while its parent product work remains `IN_PROGRESS` or `BLOCKED`. The evidence classes
below are a separate dimension and must not be read as work-status synonyms.

Keep these evidence classes separate:

| Class | Meaning at this checkpoint |
| --- | --- |
| Selected design | Accepted ADR-0028 selects governed `deno_core` as the first runtime engineering candidate; accepted ADR-0034 freezes the single-file `.mjs` first-release source/plan direction. Proposed ADRs describe reviewable directions only. |
| Implemented local mechanics | Unwired repository code and tests implement only their named passive or no-guest behavior. They are not product authority boundaries. |
| Experiment evidence | A bounded mechanism was observed in its recorded local environment. The result does not transfer automatically to an installed product or composed profile. |
| Governed external source | A named fork commit or merge establishes source identity and maintenance provenance. It does not admit artifacts, releases, runtimes, backends, or profiles. |
| Product admission | No runtime, backend, guest profile, authenticated product IPC, production approval path, or installed security boundary is admitted today. |

## Durable repository state

The current closed conformance corpus has 94 rules, 458 cases, and 561 fixtures. The unwired
Go/TypeScript implementation covers the previously recorded 177 Go and 80 TypeScript proposal
targets plus 40 independently verified MJS byte/manifest targets in each language. Go and a
standalone test-only Rust oracle additionally verify 128 exact passive Source Validator frames.
Twenty-eight exact M1 HOLD result oracles are retained and the unwired V1 artifact reproduces them,
and both are `PASSED` in their exact passive or unwired scopes. Product Source Validator work is
`BLOCKED` on the retained V2 macOS resource/confinement stop. The passive field-authority foundation covers 228 fields across 20
selected pre-freeze targets, including SourceManifest and the four passive Source Validator
objects; it does not classify future Source Preparer or plan-v1 objects.

The no-guest fixed-store lifecycle remains at E5 `local-mechanic`: exact registration,
approval/attempt, lifecycle intent/effect, recovery, and 256-active/4,096-retained capacity oracles
exist, while ownership is still injected in-process and `FakeBackend.CreatesGuest() == false`.

Archive Slice F1 is now implemented as passive `internal/execution/archivestate` types, exact
limits and known-answer digests, defensive copies, and a deterministic complete-cohort eligibility
selector. F1 writes no file, migrates no store, moves no cohort, activates no archive, resolves no
retained authority, deletes nothing, and calls no lifecycle adapter. The passive F2 format
correction now adds scope-separated global/segment indexes, typed hot/archive locations/counts, a
distinct migration-genesis checkpoint, and generated answers without changing that boundary. The
follow-on valid-v1 mapping contradiction is also passively resolved with an explicit
absent/present lifecycle union and independent lifecycle counts. The exact executable
[F2 v1 mapping resolution](SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md) retains the real crash
witness. The [stateful F2 result](SUPERVISOR_ARCHIVE_F2_MIGRATION_RESULT.md) now retains exact v2
known answers and implements only the owner-asserted v1-to-v2 migration, downgrade refusal, and
empty-archive full verifier. It preserves the one-attempt/zero-lifecycle world without recovery or
invention and adds deterministic fault, corruption, capacity, concurrency, and process-death
oracles. It moves no cohort, publishes no segment, routes no retained lookup, and calls no adapter.

Proposed ADR-0033 selects one installer-enrolled pre-created sibling inode plus a lifetime
nonblocking BSD `flock` for cooperating Supervisor ownership. The retained harness covers local
process/descriptor behavior and refusal-before-store ordering only. It is not same-UID containment,
protected storage, rollback protection, or Source Preparer store protection. G1's passive
Go/Darwin bootstrap and G2's owner-required current-v1/no-guest startup composition are both
`PASSED` in their local owned-temporary-root scopes. G3's installed protected-root, signed
bootstrap, and session/update matrix remains `BLOCKED` on the exact Team/profile and composition
dependencies below.

## First-release source and authenticated IPC

Accepted ADR-0034 now freezes the first-release source contract as one byte-exact pass-through
`main.mjs` member, one plan-v0 source-manifest role, and no static/dynamic dependency request or
module-loader fallback. The passive source-byte and SourceManifest foundation and bounded Oxc
parser/process selection, passive V0 byte contract, and unwired V1 artifact are `PASSED` in their
exact scopes. The removed handwritten scanner is
`NO_GO` because of the exact [grammar counterexample](MJS_MODULE_REQUEST_VALIDATOR_HOLD.md).
Accepted ADR-0035 has an
exact [passive V0 byte contract](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_CONTRACT.md). The V1 artifact
profile remains unsigned/not enrolled and lacks independent reproduction, V2 confinement, an
endpoint, and a consumer. V2 now retains exact fixed-I/O/fault mechanics plus the stop: `RLIMIT_AS`
returns `EINVAL`, the explicit unbounded diagnostic mutation retains file/socket/write authority,
and supported App Sandbox child entitlements change the fixed V1 bytes. Those V1-V5 gates are
`BLOCKED`; JobProposal
narrowing and all S1/M2 registration/fetch work remain `BLOCKED` on them.

The supported replacement research and R0 architecture slices pass only their exact scopes.
Direct App Sandbox inheritance is `NO_GO`. Accepted ADR-0036 selects two role-specific private
App-Sandboxed XPC launchers and matching fresh parser/profile families with no shared service,
result, cache, container, group, or key. Each writable private container is residual scratch
authority only with mandatory cleanup/residue evidence that is not a confidentiality proof. The
public footprint setter returned `KERN_NO_ACCESS`, so the accepted policy is a later evidence-
derived reactive watermark with no hard-peak/exact-cap or host-availability claim. R1-R5B now own
new passive identities, unsigned construction, separately authorized signing/install, installed
confinement/resource/residue evidence, sequential consumers, and updates. No numeric resource value,
signing identity, service, or consumer exists.

The proposed
`RegisterPlanV0` would atomically submit exact plan bytes, the
complete 562-byte role projection, the exact 87..95-byte canonical source manifest, and 0..262,144
source bytes; `GetRegisteredPlanV0` returns Supervisor-retained defensive copies. Candidate
application-data maxima are 328,337 request bytes and 332,433 fetch-reply bytes, pending generated
fixture verification from the complete field-authority projection. No S1 fixture, facade, bridge,
endpoint, consumer, or source-custody store extension exists yet.

PR #72 retained the Source Preparer as **BLOCKED**. A separately enrolled unprivileged Source
Preparer remains only a conditional later TypeScript design. P1 passive
contracts have not begun. Its entry blockers remain:

- an exact single-member protected store with baseline same-user negative-access evidence;
- exact one-shot Node worker confinement and bounded process-tree death/cleanup evidence;
- sealed installer-owned store genesis and update authority;
- settled source retention, archive, and positive release authority;
- canonical Source Preparer and plan-v1 objects with recursive nested-member field-authority
  classification and every independent validator; and
- closed cancellation, recovery, epoch, rollback, resource, and refusal-side state semantics.

Those TypeScript blockers no longer block the first release. If TypeScript is later selected,
ADR-0030 still requires one atomic plan-v1/RegisterPlanV1 cutover with no dual active v0/v1
acceptance. The 626-byte TypeScript arithmetic remains no layout, cap, or known answer.

## Governed runtime and libkrun source state

The current governed integration destination is
[`Shrimpworks/libkrun`](https://github.com/Shrimpworks/libkrun), transferred from the historical
`dills122/libkrun` location. It remains a public fork of `libkrun/libkrun`; transferred PRs #1 and
#2 retain their exact history, and `capsule/upstream-v1.19.4` currently points to
`cf0333cdba478cc34a8570a65b38412da7fd3ecc`. This transfer is repository governance only. Consumers
must pin exact reviewed commits and digests and verify ancestry from upstream anchor
`728df8125077d0db44265f6e997c72b81b65c015`, never consume the movable baseline branch name as
authority.

The governed libkrun follow-up merged PR #2 from head
`8a2c91943793668f31a1cf7af431933be935bb58` as
`cf0333cdba478cc34a8570a65b38412da7fd3ecc`. It preserves the five-patch aggregate
`d19fd0ff159c699acccda2621519de45a09408bf3847b418ac34e02b79e805d5`, fixes the observed
queued-backpressure shutdown loop and inactive-state lifecycle defect, and moves the four changed
console files from 13/88 to 37/88 covered functions and 90/728 to 298/733 covered lines. The
remaining uncovered code, branch/verifier pin mismatch, independent review, caller shared-status
hazard, hostile control/queue/descriptor work, real transport, distinct launcher, guest/VMM,
installed composition, release obligations, and final profile reruns remain open. This is governed
fork source and local library evidence only; no guest, backend, or profile is admitted.

Governed `rusty_v8` PR #4 is `PASSED` for its exact bounded fork-construction scope and merged at
head `80e863ddb942a4aa2b384e794fc23e35b9d2bb15` and merge
`cbf56de2e1156b1cf1561fdbaea7172a0aa056f4`. Exact-head workflow-dispatch run
[`30925045754`](https://github.com/Shrimpworks/rusty_v8/actions/runs/30925045754) passed the clean
Linux/arm64 network-disabled build, fixed `get_version` test, corrected GN evidence query,
evidence collection, and unsigned bundle upload. The parent governed-runtime work remains
`IN_PROGRESS — TRENDING_GOOD`: no governed release, accepted Capsule runtime bundle, independent-
builder equality, or runtime/profile admission follows from that fork result.

## Priority and dependency view

| Priority | Work | Dependency boundary |
| --- | --- | --- |
| `PASSED` local conformance slices / `IN_PROGRESS — TRENDING_GOOD` parent workstream | Archive F2/F3/F4A / F4B+ | F2 migration/full verification, F3 first immutable-segment activation, and F4A read-only retained lookup/replay/passive-collision routing pass with exact missing-history preservation, visible tombstones, publish-before-reference ordering, complete old-or-new reopen, hot/archive semantic equality, and zero adapter calls. Continue next with F4B atomic mutation/new tombstones and F4C bounded growth; production archive/store admission remains later. |
| `PASSED` local mechanic / `BLOCKED` installed lane | Owner-lock G2/G3 | G2 passed the exact current-v1/no-guest local composition. G3 is blocked: the displayed W4 certificate actually emits Team `3DDR84M4JS`, no W4 profile is cached locally, and protected-root bootstrap/signed-record/descriptor-relative store composition remains open. |
| `PASSED` prerequisites, replacement research, and R0 decision / product and downstream `BLOCKED` | `.mjs` Source Validator and S1/M2 | The exact Oxc selection, passive frames, unwired V1 artifact, negative V2 checkpoint, and ADR-0036 R0 decision passed their bounded scopes. Direct inherited helpers are rejected; two role-specific private launchers, residual scratch/cleanup, and reactive-resource claim boundaries are selected. Keep product and S1/M2 blocked. |
| `BLOCKED` sequential implementation/evidence | Source Validator R1-R5B | Run passive v1 contracts/field authority; unsigned construction; separately authorized signing/install; confinement/resource/residue corpus; daemon consumer; Broker consumer; then M2/S1. Stop on unsupported private reachability, authority/native-loading/network/filesystem escape, orphan/cleanup, mixed update, or unacceptable measured host risk. |
| `PASSED` bounded unwired conformance / signed parent `BLOCKED` | Production CBOR wrapper preparation | The [v0 object-set/wrapper result](V0_CBOR_OBJECT_SET_AND_WRAPPER.md) freezes only `SourceManifest` v0 as eligible and adds one unwired fxamacker v2.9.2 codec behind Capsule predecode, caps, canonical-byte comparison, binding, and exact-byte ownership. Plan, registration, approval, Swift, same-byte consumers, and production COSE remain blocked; `go-cose` stays test-only. |
| Future conditional | Source Preparer blockers | If TypeScript is reselected, run bounded protected-container and worker-confinement feasibility/design work, close genesis/update and retention authority, and revise the ADR if a stop condition fires. Do not start P1 bytes. |
| Independently actionable now | Documentation and field authority | Keep exact identities, counts, recursive-authority requirements, and refusal boundaries synchronized; do not classify nonexistent P1/plan-v1 fields as implemented. |
| Independently actionable now | Fork-native runtime bundle | Consume only the exact merged Deno and `rusty_v8` governed heads, reconstruct the Linux/arm64 Capsule binary/snapshot/runtime-root bundle from empty outputs and caches, compare it with retained known answers, and retain provenance/SBOM/source/notice evidence. Do not publish, sign, admit, or compose a guest in this slice. |
| Later composition | Governed runtime plus libkrun | Requires governed runtime artifacts and explicit authorization for an owned disposable development guest, followed by the exact transport/launcher/root/device/teardown corpus. |
| Credential/environment dependent | Apple Development/provisioning and installed matrices | Exact G3 discovery disproved the W4 display-name inference: certificate SHA-1 `1638...61E3` has subject OU and signed-byte TeamIdentifier `3DDR84M4JS`, and all cached profiles are also 3DDR. A matching/reissued W4 certificate and exact W4 profiles are required before signed/provisioned work. Paid owned clean-host/minimum-OS and final Developer ID/notarized matrices remain deferred and do not block unrelated local mechanics. |
| Environment dependent | Independent Linux/arm64 reconstruction | A genuinely independent builder is viable but not currently planned. Same-host and GitHub-CI evidence remains limited; independent-builder equality stays deferred. |

## Maintainer decisions and resources that can unblock work

- TypeScript remains conditional and off the first-release critical path. Follow accepted ADR-0034
  for `.mjs`; do not implement around Source Preparer P0A or widen the runtime contract.
- The intended Team remains `W4QUR9FUL4`, but exact G3 readback found that the local certificate
  displayed with that suffix has subject OU and emitted TeamIdentifier `3DDR84M4JS`; it is not W4
  evidence. Local W4 signed/provisioned experiments require a matching/reissued certificate plus
  exact W4 profiles. The three profiles Xcode 26.6 cached through
  Download Manual Profiles all belong to historical Team `3DDR84M4JS` (Gate B Broker, Gate B
  Supervisor, and wildcard); they are not reusable for W4 tests. A separate Developer ID
  Application identity for historical Team `3DDR84M4JS` is later distribution authority only: its
  use requires explicit authorization and matching-Team package design, must not enter W4
  development evidence, and does not mean Developer ID or notarization work is currently planned.
- Owned clean macOS hosts across the support floor and paid clean-host testing are deferred. They
  remain required before the corresponding installed/distribution claim can advance.
- A genuinely independent Linux/arm64 builder would strengthen byte-equality provenance, but is
  not currently scheduled. GitHub-hosted or same-host equality must not be described as independent.
- The idle governed Deno fork has transferred from `dills122/deno` to
  [`Shrimpworks/deno`](https://github.com/Shrimpworks/deno) with its fork relationship, branches,
  merged PR #1, and Actions history intact. This is repository governance only; it does not
  validate bytes, reviews, controls, releases, or admission. The governed libkrun fork has likewise
  transferred from `dills122/libkrun` to
  [`Shrimpworks/libkrun`](https://github.com/Shrimpworks/libkrun) with PR history and exact commit
  identities intact; that is also governance only. Capsule and `rusty_v8` deliberately remain at
  their current owners while their active work continues.

## Historical and exact integration record

The earlier Phase 2B/Gate C task-group observations remain durable in the
[workstream and evidence ledger](WORKSTREAM_EVIDENCE_LEDGER.md), the
[Gate C readiness checkpoint](GATE_C_READINESS_CHECKPOINT.md), and the
[Gate C P0 reconciliation](GATE_C_P0_RECONCILIATION.md). PRs #64-#75 and their exact source/merge
identities are reconciled in the ledger. Historical observations are not silently rewritten when a
later merge narrows their interpretation.

The repository still contains no job-submission endpoint, production Broker approval/signature
flow, authenticated product IPC, protected production authority/source store, admitted runtime,
real execution adapter, hostile guest, or composed production evidence chain.
