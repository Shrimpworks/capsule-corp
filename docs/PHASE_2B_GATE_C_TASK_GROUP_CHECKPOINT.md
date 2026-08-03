# Phase 2B and Gate C current maintainer checkpoint

Date: 2026-08-03

Status: integrated repository checkpoint through Capsule PR #74 merge
`e930f9dbd877bea0cbd55870060f48c9c7fdd72f`, including PR #75 merge
`6fc31a049c476acf5085071c48d3d5e36f27240f`. This is a status and dependency
index. It does not accept a Proposed ADR, activate a consumer or endpoint, admit a runtime or
backend, authorize user bytes, or authorize a guest.

## How to read the status

Keep these evidence classes separate:

| Class | Meaning at this checkpoint |
| --- | --- |
| Selected design | Accepted ADR-0028 selects governed `deno_core` as the first runtime engineering candidate. Proposed ADRs describe reviewable directions only. |
| Implemented local mechanics | Unwired repository code and tests implement only their named passive or no-guest behavior. They are not product authority boundaries. |
| Experiment evidence | A bounded mechanism was observed in its recorded local environment. The result does not transfer automatically to an installed product or composed profile. |
| Governed external source | A named fork commit or merge establishes source identity and maintenance provenance. It does not admit artifacts, releases, runtimes, backends, or profiles. |
| Product admission | No runtime, backend, guest profile, authenticated product IPC, production approval path, or installed security boundary is admitted today. |

## Durable repository state

The current closed conformance corpus remains 82 rules, 262 cases, and 368 fixtures. The unwired
Go/TypeScript implementation covers the previously recorded 177 Go targets and 80 TypeScript
proposal targets. The passive field-authority foundation covers 164 top-level fields across 15
selected pre-freeze targets; it does not recursively classify future Source Preparer or plan-v1
objects.

The no-guest fixed-store lifecycle remains at E5 `local-mechanic`: exact registration,
approval/attempt, lifecycle intent/effect, recovery, and 256-active/4,096-retained capacity oracles
exist, while ownership is still injected in-process and `FakeBackend.CreatesGuest() == false`.

Archive Slice F1 is now implemented as passive `internal/execution/archivestate` types, exact
limits and known-answer digests, defensive copies, and a deterministic complete-cohort eligibility
selector. F1 writes no file, migrates no store, moves no cohort, activates no archive, resolves no
retained authority, deletes nothing, and calls no lifecycle adapter. F2—the explicit lock-held
fixed-store v1-to-v2 migration and full empty-archive verifier—is the next archive slice.

Proposed ADR-0033 selects one installer-enrolled pre-created sibling inode plus a lifetime
nonblocking BSD `flock` for cooperating Supervisor ownership. The retained harness covers local
process/descriptor behavior and refusal-before-store ordering only. It is not same-UID containment,
protected storage, rollback protection, or Source Preparer store protection. G1—the passive
Go/Darwin bootstrap and opaque-owner port using owned temporary roots—is next; the owner-required
store composition and installed protected-root matrix remain later slices.

## Source preparation and authenticated IPC

PR #72 retained Source Preparer P0 as a bounded **P1 HOLD / NO-GO today**. A separately enrolled
unprivileged Source Preparer remains only a conditional design. P1 passive contracts have not
begun. Its entry blockers include:

- an exact single-member protected store with baseline same-user negative-access evidence;
- exact one-shot Node worker confinement and bounded process-tree death/cleanup evidence;
- sealed installer-owned store genesis and update authority;
- settled source retention, archive, and positive release authority;
- canonical Source Preparer and plan-v1 objects with recursive nested-member field-authority
  classification and every independent validator; and
- closed cancellation, recovery, epoch, rollback, resource, and refusal-side state semantics.

JavaScript-only is an acceptable bounded first-release fallback if those gates do not close. User
direction narrows that future fallback to modern ESM `.mjs` only: no CommonJS, package resolution,
legacy Node module surface, or widening of the governed runtime contract. This is planning
direction, not a frozen media/profile decision. The exact contract and applicable ADRs must be
updated and reviewed before implementation.

Authenticated IPC S1 and `RegisterPlanV1` remain blocked behind the Source Preparer decision/evidence
and the coordinated ADR-0030 plan-v1 authority model. There is no retained 562-byte S1 fixture and
the 626-byte arithmetic is not a layout, cap, or known answer. The native/Go bridge, installed
authenticated endpoints, production identities, Approval verification, and consumers remain
unimplemented.

## Governed runtime and libkrun source state

The governed libkrun follow-up merged external fork PR #2 from head
`8a2c91943793668f31a1cf7af431933be935bb58` as
`cf0333cdba478cc34a8570a65b38412da7fd3ecc`. It preserves the five-patch aggregate
`d19fd0ff159c699acccda2621519de45a09408bf3847b418ac34e02b79e805d5`, fixes the observed
queued-backpressure shutdown loop and inactive-state lifecycle defect, and moves the four changed
console files from 13/88 to 37/88 covered functions and 90/728 to 298/733 covered lines. The
remaining uncovered code, branch/verifier pin mismatch, independent review, caller shared-status
hazard, hostile control/queue/descriptor work, real transport, distinct launcher, guest/VMM,
installed composition, release obligations, and final profile reruns remain open. This is governed
fork source and local library evidence only; no guest, backend, or profile is admitted.

Governed `rusty_v8` PR #4 remains external work in progress at exact head
`aa921fa48901bf28774d61248b0187c8b91c55a4`. At the publication check on
2026-08-03T23:57:25Z, its contract jobs had passed and two clean Linux/arm64 `arm64-full-build` jobs
were still in progress. The PR is unmerged and has no accepted handoff, reusable artifact, release,
or admission effect. A successful job would still need exact artifact/evidence review and merge/
governance reconciliation before Capsule reuses it.

## Priority and dependency view

| Priority | Work | Dependency boundary |
| --- | --- | --- |
| Independently actionable now | Archive F2 | Build the explicit fixed-store v2 migration/full verifier from passive F1; no cohort leaves hot state and no archive segment exists. |
| Independently actionable now | Owner-lock G1 | Port only the passive bootstrap/opaque owner boundary under owned temporary roots; do not claim protected storage. |
| Independently actionable now | Source Preparer blockers | Run bounded protected-container and worker-confinement feasibility/design work, close genesis/update and retention authority, and revise the ADR if a stop condition fires. Do not start P1 bytes. |
| Independently actionable now | Documentation and field authority | Keep exact identities, counts, recursive-authority requirements, and refusal boundaries synchronized; do not classify nonexistent P1/plan-v1 fields as implemented. |
| Waiting | Fork-native runtime bundle | Wait for an accepted successful Linux/arm64 `rusty_v8` source/artifact handoff. Do not reuse an in-progress workflow artifact. |
| Later composition | Governed runtime plus libkrun | Requires governed runtime artifacts and explicit authorization for an owned disposable development guest, followed by the exact transport/launcher/root/device/teardown corpus. |
| Credential/environment dependent | Apple Development/provisioning and installed matrices | Current Individual Team `W4QUR9FUL4` has a valid local Apple Development identity, so signed/provisioned experiments can proceed after exact W4 role identifiers, entitlements, and profiles are deliberately created. The three cached profiles belong to historical Team `3DDR84M4JS` and are not reusable. Paid owned clean-host/minimum-OS and final Developer ID/notarized matrices are not planned now and remain deferred activation/distribution evidence, not blockers for F2, G1, or other local mechanics. |
| Environment dependent | Independent Linux/arm64 reconstruction | A genuinely independent builder is viable but not currently planned. Same-host and GitHub-CI evidence remains limited; independent-builder equality stays deferred. |

## Maintainer decisions and resources that can unblock work

- TypeScript remains conditional. If its P0A evidence cannot close without widening authority,
  proceed toward a separately reviewed modern-ESM `.mjs`-only JavaScript fallback contract rather
  than implementing around the stop.
- Current Apple Developer membership is Individual / Team `W4QUR9FUL4`, and read-only
  `security find-identity -v -p codesigning` now reports a valid Apple Development identity for
  that Team. Local signed/provisioned experiments can proceed once exact W4 role identifiers,
  entitlements, and profiles are deliberately created. The three profiles Xcode 26.6 cached through
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
  validate bytes, reviews, controls, releases, or admission. Transfers of Capsule, `rusty_v8`, or
  libkrun are later work and must not disturb their active or queued work.

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
