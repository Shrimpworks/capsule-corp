# Workstream and evidence ledger

Date: 2026-08-02

Status: durable coordination index. This ledger records where completed task conclusions were
integrated; it is not independent security evidence, a posture promotion, or a replacement for the
linked experiment results, synthesis documents, ADRs, or conformance fixtures.

The consolidated outcomes, exact known answers, application status, combined verification, and
next dependency boundary for the latest group are recorded in the
[Phase 2B and Gate C parallel-task checkpoint](PHASE_2B_GATE_C_TASK_GROUP_CHECKPOINT.md).

The first product implementation follow-up, including task provenance, executable language
coverage, the corrected PlanRegistration depth fixture, verification, and next interfaces, is
recorded in the [Phase 2 unwired decoder checkpoint](PHASE_2_UNWIRED_DECODER_CHECKPOINT.md).

## Purpose and authority

Capsule used parallel Codex tasks, external read-only reviews, retained experiments, and several
integration pull requests to move from Gate C feasibility work into the Phase 2 contract
foundation. Chat history must not be required to recover the resulting decisions.

The resumed coordination task `019fba00-c1b5-7231-945c-8df30940f145` recovered the stranded Gate C
work, dispatched the bounded adversarial replacement, reconciled the external reviews, integrated
PRs #8 through #11, and performed this final documentation audit. Its durable output is the linked
repository record, not the task transcript.

When records disagree, use this order:

1. the current accepted ADR and current contract or implementation;
2. the latest repository synthesis or reconciliation document;
3. the retained experiment result and selected evidence;
4. a task handoff or transcript; and
5. raw external-review prose.

Task IDs and pull requests provide provenance and recovery pointers. They do not make a claim
authoritative. Historical experiment observations remain valid within their recorded environment,
but later synthesis may narrow their interpretation or replace the proposed next step.

## Gate C task reconciliation

The original coordinator task `019fb58b-04a8-7121-98c9-82d304cf82a5` created five bounded
implementation-readiness tracks and then ended in `systemError`. The task failure produced no
product finding. Every completed track result was recovered into the repository and reconciled in
[Gate C implementation readiness](GATE_C_READINESS_CHECKPOINT.md).

| Workstream | Task provenance | Final disposition | Durable record |
| --- | --- | --- | --- |
| Raw storage, scratch, and controlled egress | `019fb9e7-25ec-78b2-8bf2-6558a2a9f250` | Conditional development-feasibility pass; guest read-only and fixed scratch bounds worked, but same-user host mutation defeated immutable custody and the spike extractor is not a product parser. | [`RESULTS.md`](../experiments/gate-c-libkrun-storage-egress/RESULTS.md) and Gate C synthesis |
| Console, timeout, cancellation, CPU, and memory | `019fb9e7-25ec-78b2-8bf2-659ad6ccdef1` | Conditional pass for bounded prefixes, external wall/cancel scheduling, exact forced teardown, and closed vCPU/RAM profiles. Runner exit zero, graceful-only shutdown, CPU quota, and exact host-memory claims were rejected. | [`RESULTS.md`](../experiments/gate-c-libkrun-console-lifecycle/RESULTS.md) and Gate C synthesis |
| Installed lifecycle and crash recovery | `019fb9e7-25ec-78b2-8bf2-65688a86f41f` | Conditional same-host pass for enrolled launch mechanics and exact recovery. Complete installed authority, accepted distribution, clean-host, session, reboot, and minimum-OS evidence remain open. | [`RESULTS.md`](../experiments/gate-c-libkrun-installed-recovery/RESULTS.md) and Gate C synthesis |
| Adversarial VMM and cross-job isolation | original `019fb9e7-2692-7cb0-a5c2-cebb9378e07f`; completed replacement `019fba17-9659-74f1-ab6f-82bfb72bc991` | The original task failed at the task/app layer. The bounded replacement verified and preserved the existing report. The exact profile conditionally failed because `NullFs` remains guest-visible even without a configured host directory. | [`RESULTS.md`](../experiments/gate-c-libkrun-adversarial/RESULTS.md), [`VALIDATION_RECEIPT.md`](../experiments/gate-c-libkrun-adversarial/VALIDATION_RECEIPT.md), and Gate C synthesis |
| Runtime packaging, provenance, patch governance, and supply chain | `019fb9e7-27b1-7c42-9903-8be99f620602` | Conditional build/release feasibility pass; current bytes are no-go for admission. Controlled same-host equality is not two-builder provenance, and the complete final signed/notarized topology is not available. | [`RESULTS.md`](../experiments/gate-c-libkrun-supply-chain/RESULTS.md), [`HANDOFF.md`](../experiments/gate-c-libkrun-supply-chain/HANDOFF.md), and Gate C synthesis |
| P0-0 exact stock Bun runtime authority | `019fc2e6-cf9f-7482-82b5-240992c79419`, delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | Failed for stock Bun 1.3.14: process, `execve`, FFI/native-loader, inspector, Worker, and descriptor authority remain reachable. `RUNTIME-001` stays unsupported; only a governed patched/external branch or alternate runtime remains. | [`RESULTS.md`](../experiments/gate-c-bun-runtime-authority/RESULTS.md), source inventory, synthetic probes, and selected evidence |
| P0-0 governed Bun construction | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` after merged PR #33 | NO-GO under the explicit broad/unreviewable stop rule. Exact review found a conservative 40-hand-authored plus 10-generated-output minimum before a patch or governed build; a narrow process/exec self-seal cannot independently close Worker or native loading. Alternate-runtime investigation and an ADR-0003 superseding decision are required. | [`CONSTRUCTION_REVIEW.md`](../experiments/gate-c-bun-runtime-authority/governed-closure/CONSTRUCTION_REVIEW.md), reproducible inventory, exact source hashes, and selected evidence |
| P0-0 governed `deno_core` physical omission | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` after PRs #39 and #40; merged in PR #43 | `PHYSICAL-OMISSION-PASS; NO RUNTIME ADMISSION`: a one-file patch reduced the exact built-in registry from 99 ops to three bootstrap-required ops; ASLR-controlled clean builds reproduced the snapshot and final binary, and restoration tests failed closed. TypeScript, independently reconstructible packaging/provenance, full profile admission, and external-isolation composition remain open, so `RUNTIME-001` stays unsupported. | [`RESULTS.md`](../experiments/gate-c-deno-core-physical-omission/RESULTS.md), Phase A inventory, exact patch, fixed probe, and selected evidence |
| P0-0 governed `deno_core` reproducible package | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` after merged PR #43 | BYTE-REPRODUCTION PASS; RUNTIME-SELECTION EVIDENCE NO-GO. A digest-pinned no-apt builder and complete 191-crate offline source bundle reproduced the prior snapshot and binary across two clean same-host containers with complete declared package-byte equality. Independent-builder provenance, archive-corresponding V8 source/notices, a standalone dynamic runtime root, and production TypeScript ownership/wiring remain blockers; `RUNTIME-001` stays unsupported. | [`RESULTS.md`](../experiments/gate-c-deno-core-reproducible-package/RESULTS.md), bundle manifest, CycloneDX SBOM, source/license inventory, unsigned provenance, and selected evidence |
| P0-0 governed `deno_core` V8 source/license closure | `019fc58f-90fd-78c2-88f8-20c7a1f38359`, delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` after merged PR #50 | `SOURCE-LICENSE-CLOSURE-NO-GO`. Exact asset/job, rusty_v8/V8 revisions, 20 gitlinks, Chromium V8 base, four-patch stack, 1,875 archive members, and 726 license/notice candidates are retained. Mutable publisher inputs, missing GN/Ninja link closure, and absent generated notices prevent exact rebuilding and complete notice/SBOM closure. Governed `deno_core` is the intended engineering direction, not an admitted runtime; `RUNTIME-001` stays unsupported. | [`RESULTS.md`](../experiments/gate-c-deno-v8-source-license-closure/RESULTS.md), [`SOURCE_PUBLICATION.md`](../experiments/gate-c-deno-v8-source-license-closure/SOURCE_PUBLICATION.md), exact manifests, archive inventory, and fail-closed checklist |
| P0-0 governed `deno_core` self-contained runtime root | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` after merged PR #50 | STANDALONE DYNAMIC ROOT PASS; NO RUNTIME ADMISSION. Exact Debian snapshot packages produced a reproducible 22-entry root for the unchanged binary/snapshot; empty-environment scratch execution and file-open/mutation evidence used no ambient Bookworm library/config/NSS/locale/timezone/package/cache bytes. Independent-builder, V8 publication, real governed-fork/release, TypeScript, external-isolation, and profile-admission blockers remain; `RUNTIME-001` stays unsupported. | [`RESULTS.md`](../experiments/gate-c-deno-core-runtime-root/RESULTS.md), closed root/package manifests, ELF/file-open/mutation evidence, SBOM, and source/notice mapping |
| P0-0 TypeScript approved-byte boundary | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` from merged checkpoint `5448943` | Passed only the exact pre-approval byte-ownership question for a strip-only ESM TypeScript subset. Exact Node 22.22.1/Amaro 1.1.5 emission was deterministic for fixed fixtures and mutation/cap/diagnostic cases failed closed. Proposed ADR-0026 binds original and executable source roles before registration. No current contract changed, no owner/runtime was selected, and `RUNTIME-001` remains unsupported. | [`RESULTS.md`](../experiments/typescript-approved-byte-boundary/RESULTS.md), fixed Node and `deno_ast` probes, independent Go verifier, and selected evidence |
| P0-1 FD-native immutable runtime-root custody | `019fc4c1-7d40-77b3-a2e9-51d3e2775972`, delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | PATCH-CANDIDATE: the narrow fixed-role raw-only API passed local attachment/custody, focused sanitizer tests, five deliberate mutations, and four owned unsandboxed guest digest runs with no root-path opens. P0-1C remains open until the exact final signed installed App Sandbox/protected-construction corpus passes. | [`RESULTS.md`](../experiments/gate-c-libkrun-root-custody/RESULTS.md), [`FD_NATIVE_PATCH_REVIEW.md`](../experiments/gate-c-libkrun-root-custody/FD_NATIVE_PATCH_REVIEW.md), governed patch, and selected evidence |
| P0-2 `NullFs` disposition | Earlier replacement `019fc2e8-445e-7cb2-b4c2-54d84282c3fe`, replacing task `019fc2e6-cf9d-7210-b2f3-f3bf2244e83a`; later prototype merged in PR #30 | `GOVERNED-PATCH`: the smallest deletion failed bootstrap, but the later direct-block-root prototype booted without virtiofs, reran 36 adversarial plus four identity cases without the original failure, and made removal credible. It is not admitted; independent patch review, route closure, P0-1 custody, P0-3 transport, and final signed P0-4 evidence remain. | [`NULLFS_P0_2.md`](../experiments/gate-c-libkrun-adversarial/NULLFS_P0_2.md), [`NULLFS_P0_2_DISPOSITION.md`](../experiments/gate-c-libkrun-adversarial/NULLFS_P0_2_DISPOSITION.md), governed prototype patch, and compact evidence |
| P0-3 backend-independent framing | Merged in PR #27 | Conditional candidate pass only: 43 byte-exact vectors measured the proposed source/input/result/frame caps and retained binding, role, JSON, commit, drain, stall/death, EOF, runner-exit, and crash dispositions. No transport, launcher, guest, VMM, App Sandbox, Supervisor, approval, or teardown mechanism participated. | [`RESULTS.md`](../experiments/gate-c-p0-3-protocol-conformance/RESULTS.md), 43-vector manifest, local model, and measurement record |
| P0-3 libkrun console correctness | Merged in PR #28 | Stock route cannot proceed as-is. Governed patch `584ce48548fe969684fe3c55e57fbf56e7dae40af28c241c24c47b138faf1283` passed 51 local library tests, including four new regressions, but remains unreviewed and lacks full sanitizer/coverage and real transport/launcher/guest/installed composition. | [`RESULTS.md`](../experiments/gate-c-libkrun-console-correctness/RESULTS.md), governed patch, verification record, and focused tests |
| P0-4A installed development topology | Merged in PR #34 | Conditional topology pass only: 18 roles, 17 installed-entry readbacks, per-user registration/explicit activation, exact ad-hoc IPC identity, refusal cases, and same-session recovery passed without host root or a guest. App Sandbox failed before `main`; valid signing, Team enrollment, notarization/stapling/Gatekeeper, on-demand activation, clean hosts, sessions, and the macOS support floor remain open. | [`RESULTS.md`](../experiments/gate-c-installed-development-topology/RESULTS.md), closed manifests, scripts/tests, and selected installed-run evidence |

The retained generated `.build/` and `.runs/` directories are intentionally disposable and ignored.
The repository keeps the source, reproduction instructions, selected evidence, validation receipts,
machine-readable inputs, and conclusions needed to reproduce or challenge the findings. No product
package may import the experiment implementations.

## P0 review and research reconciliation

After the five tracks, the plan received three additional forms of challenge:

- an independent Claude read-only review supplied through an external handoff;
- a separate ChatGPT Pro third-eye review of that plan and reconciliation; and
- a fresh-context issues-only review of the updated plan and pinned upstream source.

The raw review prose is not a project authority and is not required for recovery. Every accepted,
rejected, or narrowed conclusion is retained in
[Gate C P0 reconciliation](GATE_C_P0_RECONCILIATION.md), including:

- direct inherited read-only-descriptor custody as a hypothesis with three separately falsifiable
  subclaims, not a completed immutability mechanism;
- independent `NullFs` disposition;
- bounded attempt-bound console-port source/input and typed completion for the inline slice;
- deferral of the ext4/raw-image parser until file artifacts;
- construction-level governed runtime-authority closure;
- a distinct trusted launcher and exact host/guest descriptor allowlists;
- hostile virtio-console control, queue, descriptor, backpressure, shutdown, and partial-write
  coverage; and
- early topology testing followed by final signed/notarized-byte rebuild and revalidation.

The authoritative remaining P0 and P1 campaigns are in the Gate C synthesis and roadmap. Both the
stock-runtime and governed-construction branches of Bun P0-0 are closed as failures. ADR-0028 now
selects governed `deno_core` as the first engineering candidate and preserves the runtime-neutral
protocol, but no profile is admitted. The reviews and experiments did not authorize libkrun to
handle user bytes.

## Phase 2 parallel-task reconciliation

Three read-only tasks independently reviewed the first passive contract slice. Their full useful
output is reconciled in
[Phase 2A parallel-review synthesis](PHASE_2A_PARALLEL_REVIEW_SYNTHESIS.md).

| Review | Task ID | Retained conclusion |
| --- | --- | --- |
| Contract shape | `019fbd6b-0de7-7201-befc-be983552e5ed` | Three closed, non-interchangeable objects and role-specific scalars are the right boundary; backend-dependent fields and unresolved values must remain deferred. |
| Deprecated mixed-`Job` migration | `019fbd6b-0de5-7003-82b2-40049e07fcb7` | A clean cutover is possible, but dormant direct-execution surfaces, adjacent profile/receipt models, SDK/MCP consumers, examples, tests, and documentation must migrate atomically without a permissive adapter. |
| Fail-closed conformance corpus | `019fbd6b-0de5-7003-82b2-3febecc775aa` | Use one manifest-driven byte-exact corpus with explicit raw, schema/CDDL, semantic, and registration ownership plus no-authority-state-change rejection oracles. |

These reviews were intentionally read-only. Their deliverable was the synthesis and ordered work
plan, not separate code or experiment directories.

## Phase 2B approval/attempt Slice A-C reconciliation

The approval/attempt work was delivered as three merged, backend-independent slices under parent
orchestrator task `019fc2de-552d-77a0-aa47-35ac39d02edc`. It is defensive local conformance work,
not consumer activation or backend admission.

| Slice | Durable implementation checkpoint | Retained evidence and limitation |
| --- | --- | --- |
| A: passive contracts | `internal/execution/approvalattempt` defines distinct `ApprovalID`, `AttemptID`, and `AttemptNonce` domains, typed references, three approval states, one created-attempt state, the fixed 12-class vocabulary, defensive byte ownership, and the retained-vector-only verifier. | `TestSliceAConformanceManifest` covers 44 cases (10 accept, 34 reject), the exact 375-byte known answer, 512/256/128 raw budgets, and 431/242/116 calculated maxima. It is not production COSE, signer authorization, or cryptography. |
| B: fixed authority store | `registrationstate.FixedFileStore` colocates registration, effective-time high-water, approval, and immutable created-attempt state; payload-digest idempotency, retained nonce uniqueness, atomic consume/create, exact replay/concurrency, capacity, recovery fencing, reopen validation, and corruption refusal are implemented locally. | Twelve manifest state-transition cases plus focused component tests retain the exact behavior. The fixed no-eviction file-snapshot store has no production archive/compaction, multi-process locking, backup/restore, or rollback-resistant uniqueness and remains unwired. |
| C: attempt-keyed fake lifecycle | `ApprovalAttemptComponent` supplies `ResolveCreated` and created-attempt enumeration; `registeredlifecycle` accepts only `AttemptID`, revalidates every copied binding and exact plan before fake prepare, keeps attempts for one registration distinct, and keys replay/fault/recovery state by attempt. | Twelve top-level focused tests cover binding refusal, exact/concurrent replay, startup enumeration, all fake fault moments, and post-effect restart recovery. `FakeBackend.CreatesGuest() == false`; lifecycle `MemoryStore` is bounded single-process memory and non-durable. |

The closed conformance manifest remains 82 rules, 262 cases, and 368 fixtures after Slice C. Its Go
manifest-backed coverage remains 177 targets: 81 internal-CBOR/wrapper cases, 40 registration-state
cases, 44 passive approval/attempt cases, and 12 Slice B state transitions. Slice C adds focused Go
tests rather than manifest cases.

## Integrated pull-request checkpoints

| Integration | Merge commit | What became durable |
| --- | --- | --- |
| PR #8, Gate C feasibility and readiness | `633d249` | Five retained libkrun readiness experiments, Gate C synthesis, ADR-0022 updates, roadmap, and evidence-matrix reconciliation. |
| PR #9, Phase 2A contract foundation | `3a75098` | Passive `JobProposal`, minimum `ExecutionPlan` and `PlanRegistration` candidates, byte-exact fixtures, role-specific Go/TypeScript views, and parallel-review synthesis. |
| PR #10, Phase 2B boundary decisions | `e565e6b` | Proposed ADR-0023 and an implementation-ready breakdown for strict decoding, exact identity, registration, and conformance. |
| PR #11, Phase 2B conformance foundation | `f6de7ec` | Closed manifest, integrity verifier, 37 rules, 105 cases, and 91 unique raw/media/scalar/CBOR fixtures; corrected derived JSON-node and `PlanRegistration` item limits. |
| PR #27, P0-3 protocol candidate | `35b8da7` | Retained the 43-vector backend-independent framing candidate and measurement evidence; no backend or guest. |
| PR #28, console correctness | `7ba37d3` | Retained the governed console patch decision, SHA-256, and 51-test local library result. |
| PR #29, P0-1 initial custody | `b4307af` | Retained `/dev/fd/N` attachment/frozen-object evidence and the open signed installed P0-1C boundary. |
| PR #30, P0-2 direct root | `7fb4146` | Selected `GOVERNED-PATCH`; direct-block-root prototype removed the `NullFs` device in the bounded rerun without admitting the profile. |
| PR #31, site copy | `1614b4f` | Plain-language site rewrite only; no Gate C or Phase 2 security status changed. |
| Local Phase 2B Task 2.3 integration | `4afbdfa` | Proposal/source/input rules, fixed resolver contexts, and known-answer source-manifest/canonical-input fixtures; corpus expanded to 55 rules, 147 cases, and 130 fixtures. |
| Local Gate C P0-2 evidence integration | `d3be865` | Retained the bounded `NullFs` route/removal investigation and kept the exact profile unsupported. |
| Local Gate C P0-0 evidence integration | `004c047` | Rejected stock Bun 1.3.14 for `RUNTIME-001` and retained the source/probe evidence and governed-patch-or-alternate decision. |
| Local Phase 2B Task 2.4 integration | `fd51ac4` | Exact plan/registration/domain/state fixtures and proposed ADR-0023 state addendum; corpus expanded to 67 rules, 206 cases, and 278 fixtures. |
| PR #32, approval/attempt Slice A | `5891c1d` | Passive typed domains/references/states/classifications, exact approval known answer, bounded retained-vector verifier, defensive copies, and 44 manifest cases; corpus expanded to 78 rules, 250 cases, and 350 fixtures. |
| PR #33, Bun source toolchain | `b804889` | Established the governed Bun source-build inventory/toolchain; no runtime admission. |
| PR #34, P0-4A topology | `6ddc576` | Retained the conditional 18-role no-host-root topology result and exact signing/notarization/minimum-OS gaps. |
| PR #35, governed Bun no-go | `170e988` | Closed the Bun construction branch as NO-GO under its broad/unreviewable threshold. |
| PR #36, approval/attempt Slice B | `2adbd5f` | Colocated fixed registration/approval/attempt snapshot, durable effective-time high-water, canonical-payload idempotency, atomic consume/create, capacity, fault/reopen/corruption handling, and 12 state-transition cases; corpus expanded to 82 rules, 262 cases, and 368 fixtures. |
| PR #37, approval/attempt Slice C | `9955d6b` | `AttemptResolver`/created-attempt enumeration, `AttemptID`-only fake lifecycle and recovery, binding revalidation, distinct attempts per registration, startup enumeration, fault/replay recovery, and explicit no-guest invariant. Manifest counts did not change. |
| PR #38, Slice D checkpoint | `6704eb5` | Reconciled the unwired approval/attempt Slices A-C status and remaining boundaries. |
| PR #39, Deno-family no-go | `7db7e39` | Rejected full Deno and the tested unpatched/middleware `deno_core` construction; `RUNTIME-001` remained unsupported. |
| PR #40, durable lifecycle design | `c88b6e3` | Proposed ADR-0025 and the E1-E5 fake-only durable lifecycle plan; implemented no durable behavior. |
| PR #41, FD-native custody | `e447ed8` | Promoted the narrow raw-only FD-native API to `PATCH-CANDIDATE` after local, mutation, and four owned unsandboxed guest runs; P0-1C remained open. |
| PR #42, passive lifecycle E1 | `504e44c` | Implemented runtime-neutral passive lifecycle types and focused tests only; E2-E5 remained unwired. |
| PR #43, physical omission | `5448943` | Recorded `PHYSICAL-OMISSION-PASS; NO RUNTIME ADMISSION` for the governed `deno_core` construction; `RUNTIME-001` remained unsupported. |
| PR #44, Gate C/Phase 2B reconciliation | `3b5fd40` | Reconciled P0-1 through P0-4 evidence and recorded E1 as the then-current durable lifecycle boundary. |
| PR #45, fail-closed/schema fixes | `37fee54` | Replaced panic paths with fail-closed errors and tightened schema bounds without promoting a product claim. |
| PR #46, TypeScript approved bytes | `990d602` | Recorded `BOUNDARY-PASS; NO RUNTIME ADMISSION` and Proposed ADR-0026 for exact strip-only transformation before registration and approval. |
| PR #47, durable lifecycle E2 | `e815179` | Added explicit fixed-store v1 lifecycle-set validation and lock-asserted v0-to-v1 migration with fault/reopen/downgrade refusal. |
| PR #48, output ceiling gap | `0b1adee` | Recorded the deferred output `maxBytes` ceiling as a known contract gap. |
| PR #49, SupervisorCore retirement | `6cdd1bb` | Removed the unused legacy in-memory Supervisor scaffold under accepted ADR-0027; the approval/attempt/lifecycle split is now the sole implementation path. |
| PR #50, governed runtime packaging | `cb987e1` | Reproduced the governed `deno_core` binary/snapshot from pinned inputs on one host, retained package/SBOM/provenance evidence, and kept selection evidence at NO-GO. |
| PRs #51-#53, quality hardening | `64cf70e`, `711e290`, `0fcea59` | Added protocol/API documentation and fail-closed SDK response validation; no runtime/backend admission changed. |
| PR #54, durable lifecycle E3 | `4200033` | Added the unwired fixed-store lifecycle transaction port, stable effect permits, exact fault/reopen/reconciliation behavior, and no adapter calls. |
| PR #55, CI hardening | `a5cb64a` | Added security scanning, deeper linting, ADR-index validation, and coverage gates. |
| PR #56, V8 source/license closure | `abfdaa5` | Recorded `SOURCE-LICENSE-CLOSURE-NO-GO`: exact official source/patch identity was mapped, but mutable publisher inputs, missing build closure, and absent generated notices block governed reuse. |
| PR #57, self-contained runtime root | `96383a8` | Recorded `STANDALONE DYNAMIC ROOT PASS; NO RUNTIME ADMISSION` for the exact 22-entry package-derived root with no ambient library/config fallback. |

The merge commits, not the former draft-PR state recorded in task responses, are the integration
checkpoints. The repository was clean at `f6de7ec` except for the user-owned untracked `.claude/`
directory, which is not project evidence and is intentionally excluded from Git.

## Current boundary and remaining work

Completed and retained:

- all five Gate C readiness tracks and their P0 reconciliation;
- Phase 2A passive candidate contracts and three independent reviews;
- Phase 2B exact proposed decoder/registration decisions;
- conformance Tasks 2.1 through 2.4: manifest integrity, foundational byte fixtures, fixed proposal
  resolver contexts, source/canonical-inline-input known answers, and exact
  plan/registration/domain/state fixtures; and
- approval/attempt Slices A-C: passive typed contracts and fixture verifier, the unwired colocated
  fixed authority store with atomic consume/create, and the `AttemptID`-keyed no-guest fake
  lifecycle seam; and
- durable-lifecycle Slices E1-E3: passive runtime-neutral types, explicit fixed-store v1 migration
  and validation, and unwired ensure/read/intent/result/indeterminate/reconciliation/recovery-set
  transactions. No E3 path calls an adapter, consumer, runtime, backend, or guest; and
- governed `deno_core` physical omission, same-host package reproduction, exact V8 closure NO-GO,
  and standalone dynamic-root evidence. Accepted ADR-0028 selects its engineering order without
  admitting a profile; the real Deno and `rusty_v8` forks exist with no governed branches yet.

Next backend-independent work:

1. Continue [Proposed ADR-0025](adr/0025-colocate-durable-attempt-lifecycle-state.md) with Slice E4
   in the [fake-only durable lifecycle plan](PHASE_2B_DURABLE_ATTEMPT_LIFECYCLE_PLAN.md): migrate
   `registeredlifecycle` from `MemoryStore` to the E3 durable transaction port, add stable fake
   effect/instance identities and startup coordination, and rerun the death/fault/reopen matrix.
   Slice E5 remains separate and no product or guest claim advances.
2. In parallel, bootstrap governed branches and reviewable draft PRs in `dills122/deno` and
   `dills122/rusty_v8` from the exact ADR-0028 upstream commits. Follow the
   [governed runtime work plan](GOVERNED_DENO_CORE_WORK_PLAN.md); do not substitute local copies or
   experiment patches for real fork commits/releases.
3. Separately design reviewed Supervisor archive/compaction and replay retention. The fixed
   no-eviction authority store is not a continuous-service store.
4. After the Supervisor language/privilege topology is selected, design authenticated typed IPC
   and production Approval verification/authorization; do not promote the conformance JSON or
   injected caller context into a product transport.
5. Only after those decisions, consumers, evidence composition, and the daemon aggregate service
   envelope exist may a coordinated public cutover and mixed-`Job` removal be considered.

This checkpoint does not decide IPC topology, authority-store archive/compaction, production
COSE/Swift/Keychain/user-presence signing, consumer ownership, evidence composition, or public
cutover. The authority/lifecycle snapshot lacks real multi-process locking and rollback-resistant
identifier/nonce/effect uniqueness, while the active driver remains single-process and
non-durable. Content, evidence, runtime, backend, and guest remain absent from the unwired path.

In parallel, libkrun remains barred from user bytes until the five reconciled P0 campaigns close:
runtime-authority closure, immutable runtime-root custody, `NullFs` disposition, typed port
transport/completion, and an admissible complete installed development bundle. File-artifact
parsing and stronger validation campaigns remain later gates as documented in the roadmap.

## Maintenance rule

For every future parallel task or external review that changes the project direction:

1. record its task ID or source, bounded scope, decision, and durable evidence location here or in
   a linked phase-specific ledger;
2. place reproducible observations in a retained experiment or conformance fixture before relying
   on the task handoff;
3. reconcile design conclusions into the current ADR, synthesis, evidence matrix, and roadmap;
4. mark superseded recommendations explicitly rather than deleting historical observations; and
5. record the merge commit that integrated the result.

No task is complete for coordination purposes if its only useful output exists in chat history.
