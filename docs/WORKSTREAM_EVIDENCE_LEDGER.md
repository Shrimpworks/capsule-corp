# Workstream and evidence ledger

Date: 2026-08-03

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
| P0-0 fork-native governed runtime bundle | `019fc8f3-44d1-77b1-86c6-8a38e2b47211`, delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` after the Deno/rusty_v8 fork PRs merged | LINUX/ARM64 CONSTRUCTION BLOCKED; NO BUILD OR RUNTIME ADMISSION. Exact fork refs/ancestry, Deno three-op/fixture oracles, and the rusty_v8 20-gitlink/source/tool locks passed, but the only merged rusty_v8 builder is Linux/amd64. The experiment stopped before prefetch/build and did not substitute amd64. | [`RESULTS.md`](../experiments/gate-c-fork-native-deno-runtime-bundle/RESULTS.md), closed intended input contract, exact ref/lock evidence, and offline verifier |
| P0-0 TypeScript approved-byte boundary | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` from merged checkpoint `5448943` | Passed only the exact pre-approval byte-ownership question for a strip-only ESM TypeScript subset. Exact Node 22.22.1/Amaro 1.1.5 emission was deterministic for fixed fixtures and mutation/cap/diagnostic cases failed closed. Proposed ADR-0026 binds original and executable source roles before registration. No current contract changed, no owner/runtime was selected, and `RUNTIME-001` remains unsupported. | [`RESULTS.md`](../experiments/typescript-approved-byte-boundary/RESULTS.md), fixed Node and `deno_ast` probes, independent Go verifier, and selected evidence |
| ADR-0034 M1 source-byte/manifest foundation and language-validator hold | `019fca55-554d-7b23-bb2e-6d2c4b0e16cd`, delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | Passive foundation only: exact MJS identifiers, strict UTF-8/BOM and byte caps, deterministic single-member SourceManifest, and Go/TypeScript known answers pass. A retained division-versus-regexp counterexample proved the ad hoc scanner had a live-`import()` false negative, so it was removed and JobProposal/plan/M2/S1 work is on HOLD pending a separate parser-boundary decision. No runtime/no-loader claim follows. | [`MJS_MODULE_REQUEST_VALIDATOR_HOLD.md`](MJS_MODULE_REQUEST_VALIDATOR_HOLD.md), generated MJS conformance fixtures, SourceManifest CDDL, and field-authority target |
| `.mjs` parser/process boundary | `019fc2de-552d-77a0-aa47-35ac39d02edc` user-visible defensive/local-only research delegation | `PARSER-CANDIDATE-GO; PRODUCT-NO-GO`. Exact Oxc 0.140.0 with parser diagnostics, semantic early-error checking, and unresolved free-CommonJS-reference accounting matched all 33 local parse-only cases and all 28 canonical merged M1 HOLD outcomes in place (22 grammar/ordinary plus six CommonJS-reference cases). Exact deno_ast/SWC and tree-sitter modes recovered from required parse errors, and V8 compile-only could not observe dynamic import/`import.meta`. Proposed ADR-0035 selects a future one-shot stateless Source Validator outside daemon/Broker/Supervisor address spaces. No product validator, process sandbox, runtime no-loader control, runtime, backend, or guest exists. | [`RESULTS.md`](../experiments/mjs-parser-boundary/RESULTS.md), M1 HOLD mapping, exact lock/fixtures/classifications, supply inventory, measurements, and [`HANDOFF.md`](../experiments/mjs-parser-boundary/HANDOFF.md) |
| P0-1 FD-native immutable runtime-root custody | `019fc4c1-7d40-77b3-a2e9-51d3e2775972`, delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | PATCH-CANDIDATE: the narrow fixed-role raw-only API passed local attachment/custody, focused sanitizer tests, five deliberate mutations, and four owned unsandboxed guest digest runs with no root-path opens. P0-1C remains open until the exact final signed installed App Sandbox/protected-construction corpus passes. | [`RESULTS.md`](../experiments/gate-c-libkrun-root-custody/RESULTS.md), [`FD_NATIVE_PATCH_REVIEW.md`](../experiments/gate-c-libkrun-root-custody/FD_NATIVE_PATCH_REVIEW.md), governed patch, and selected evidence |
| P0-2 `NullFs` disposition | Earlier replacement `019fc2e8-445e-7cb2-b4c2-54d84282c3fe`, replacing task `019fc2e6-cf9d-7210-b2f3-f3bf2244e83a`; later prototype merged in PR #30 | `GOVERNED-PATCH`: the smallest deletion failed bootstrap, but the later direct-block-root prototype booted without virtiofs, reran 36 adversarial plus four identity cases without the original failure, and made removal credible. It is not admitted; independent patch review, route closure, P0-1 custody, P0-3 transport, and final signed P0-4 evidence remain. | [`NULLFS_P0_2.md`](../experiments/gate-c-libkrun-adversarial/NULLFS_P0_2.md), [`NULLFS_P0_2_DISPOSITION.md`](../experiments/gate-c-libkrun-adversarial/NULLFS_P0_2_DISPOSITION.md), governed prototype patch, and compact evidence |
| P0-3 backend-independent framing | Merged in PR #27 | Conditional candidate pass only: 43 byte-exact vectors measured the proposed source/input/result/frame caps and retained binding, role, JSON, commit, drain, stall/death, EOF, runner-exit, and crash dispositions. No transport, launcher, guest, VMM, App Sandbox, Supervisor, approval, or teardown mechanism participated. | [`RESULTS.md`](../experiments/gate-c-p0-3-protocol-conformance/RESULTS.md), 43-vector manifest, local model, and measurement record |
| P0-3 libkrun console correctness | Merged in PR #28 | At that checkpoint stock could not proceed as-is. Governed patch `584ce48548fe969684fe3c55e57fbf56e7dae40af28c241c24c47b138faf1283` passed 51 local library tests and four regressions but still lacked the later sanitizer/coverage follow-up and all real composition. | [`RESULTS.md`](../experiments/gate-c-libkrun-console-correctness/RESULTS.md), governed patch, verification record, and focused tests |
| P0-3 cross-language/console follow-up | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` on 2026-08-03 | This retained pre-fork checkpoint added independent Node agreement on all 43 vectors, six re-encodings, ten local pipe fault classes, Clippy, AddressSanitizer, repetition, four mutations, and the historical before measurement of 90/728 patched-file lines. The governed-fork row below supersedes it for current coverage. | Updated P0-3 and console `RESULTS.md`, independent verifier/fault harness, cross-language evidence, coverage summary, and mutation patches |
| Governed libkrun console/raw-FD source reconciliation | User-visible read-only reconciliation delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc`; public fork PR #2 | Exact head `8a2c91943793668f31a1cf7af431933be935bb58` merged as `cf0333cdba478cc34a8570a65b38412da7fd3ecc` over governed base `4ea8d1de861ed1c0636fc800b6da8fb71a086aa5`. The unchanged five-patch aggregate is `d19fd0ff159c699acccda2621519de45a09408bf3847b418ac34e02b79e805d5`. The merge fixes queued-backpressure shutdown and inactive-state defects; four-file coverage moves from 13/88 to 37/88 functions and 90/728 to 298/733 lines. The no-guest property/raw-FD, ASan, Clippy, repetition, mutation, reconstruction, and macOS default-init cross-build head checks passed. At 2026-08-03T22:57:43Z Linux-arm64 unit tests remained queued, the overall head state was pending, no merge-commit checks or submitted independent/CODEOWNER review existed, and the advanced baseline branch no longer satisfied the verifier's hardcoded earlier-base pin. No guest, VMM transport, installed product, release, or admission was evidenced. | [Public PR #2](https://github.com/dills122/libkrun/pull/2), exact [head](https://github.com/dills122/libkrun/commit/8a2c91943793668f31a1cf7af431933be935bb58), exact [merge](https://github.com/dills122/libkrun/commit/cf0333cdba478cc34a8570a65b38412da7fd3ecc), and canonical Gate C reconciliation |
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

## Related-systems review reconciliation

A later public-source comparison considered Deno's permission broker, HARP, Redan, Qubes qrexec,
XDG portals, Confidential Containers/Kata, reproducible-build systems, DSSE/in-toto, capability
systems, and agent-sandbox products. The bounded comparison is retained in
[Related systems and design influences](RELATED_SYSTEMS.md). It is an input to design review, not
implementation evidence, a competitor-completeness claim, a novelty opinion, or authority to relax
Capsule's boundary.

The review reinforced existing decisions: approval cannot add authority absent from the registered
plan, indeterminate or failed resolution must not be treated as absence, content should cross trust
boundaries through narrow handles, the guest is hostile, receipts describe observed execution rather
than prove correctness, and network remains out of the first slice. It did not justify importing
full Deno permission brokering, dynamic authority negotiation, URL fetches, Unix-socket trust, or
product-specific policy formats.

One concrete planning refinement was retained: before schema freeze, every authority-bearing field
must have a machine-readable classification naming its origin role, validator, authority effect,
approval visibility, content status, binding, and unknown-field behavior. Unclassified fields fail
verification. This tightens the coordinated object-model migration without changing the current
priority order.

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
| C: attempt-keyed fake lifecycle | `ApprovalAttemptComponent` supplies `ResolveCreated` and created-attempt enumeration; `registeredlifecycle` accepts only `AttemptID`, revalidates every copied binding and exact plan before fake prepare, keeps attempts for one registration distinct, and keys replay/fault/recovery state by attempt. | The original twelve top-level focused tests cover binding refusal, exact/concurrent replay, startup enumeration, all fake fault moments, and post-effect restart recovery. E4/E5 now preserve those oracles through the colocated v1 store and stable fake effect/instance IDs. `FakeBackend.CreatesGuest() == false`; no real adapter or guest exists. |

## Phase 2B durable lifecycle Slice E1-E5 reconciliation

The local durable-lifecycle checkpoint implements Proposed ADR-0025 without accepting it or
activating a consumer:

| Slice | Retained checkpoint | Evidence and limitation |
| --- | --- | --- |
| E1-E3 | Closed passive record/effect types, explicit v0-to-v1 migration/open validation, and colocated lifecycle transactions | Exact binding, collision, migration, corruption, fault, and recovery-set tests; no adapter call through E3 |
| E4 | `registeredlifecycle` drives only the no-guest fake through stable store-issued effect permits and exact fake instance identities | PR #59 retains before/after-effect, death, reopen, reconciliation, backoff, owner-session, and fake-no-guest oracles; the owner/coordinator is injected in-process |
| E5 | Joined v1 capacity and repeated-startup checkpoint | `TestFixedStoreV1ExactActiveCapacityReleasesOnlyAfterDurableDestroy` proves 256 active and destroyed-only release; `TestFixedStoreV1ExactRetainedLifecycleCapacityNeverEvicts` proves 4,096 retained and cap-plus-one no rewrite/eviction with a 30,321,818-byte population; repeated/concurrent startup tests prove terminal omission, unresolved/exhausted retention, no fourth observation, and owner/coordinator mismatch refusal |

This is an unwired `local-mechanic` result. Passive archive F1 types and selection now exist, but
archive file/store behavior, installed protected-root/bootstrap evidence, backup/rollback,
production reconciliation, consumers, content, evidence, runtime, backend, and guest remain absent.

## Supervisor owner-lock design checkpoint

Proposed ADR-0033 selects a pre-created installer-enrolled sibling object, exact
UID/mode/type/link/device/inode validation, and lifetime nonblocking BSD `flock`. The retained
[development-only experiment](../experiments/supervisor-owner-lock-boundary/RESULTS.md) observed
duplicate-process refusal before store/recovery markers, last-description/process-death release,
dup/fork/exec/`CLOEXEC` behavior, POSIX record-lock any-close risk, macOS 26 OFD behavior,
`O_EXLOCK` interoperability, retained-directory `openat`, and rename/unlink/replacement limits.

Passive G1 advances the exact internal Go/Darwin acquisition to `local-mechanic`: closed
enrollment and entry-name validation, retained-root `openat`, exact descriptor flags and identity,
nonblocking independent-process contention, opaque owner/session/close behavior, fault injection,
inheritance, process-death, replacement, descriptor-reuse, and refusal-before-downstream markers
are covered under owned temporary roots. The existing migration assertion is also exercised with
an actual held owner. Bounded G2 now adds the owner-required current-v1 opener, same-session
store/coordinator, exact ownership-before-store-before-sorted-recovery composition, permanent
post-open fence, ordered close, response-loss recovery, and child-process death/reopen. No signed
bootstrap, Apple-signed protected state root, wrong-user/session/update/reboot result, production
engine, archive behavior, consumer, runtime, real backend, or guest exists. The advisory lock does
not contain a same-UID process that can mutate its parent directory.

Source Preparer P0 remains a separate bounded NO-GO/HOLD merged in PR #72 from head
`a12041c36d90815474598f0929c595b32dc68e11` as
`2e268b01d4174fe90397c00abc5973a3dd785606`: no single-member source-store container, exact worker
confinement evidence, closed store-genesis/update authority, settled release/archive authority, or
recursive nested-member field-authority design is proven. P1 passive contracts have not begun.
ADR-0033 process exclusion does not supply source-store
confidentiality, integrity, or protected-container membership. Useful composition would require a
separate protected root, component identity, lock object, and owner session after those properties
are independently proven; Supervisor ownership material is not shared.

The archive design is retained as
[Proposed ADR-0031](adr/0031-checkpoint-closed-supervisor-cohorts.md) and the
[Supervisor archive/compaction conformance plan](SUPERVISOR_ARCHIVE_COMPACTION_PLAN.md). It selects
complete closed registration cohorts, immutable full-record segments, exact replay/non-reuse
tombstone indexes, publish-before-activate fault ordering, coherent backup verification, and
read-only offline verification. Slice F1 now implements passive archive projections, exact
limits/known answers, defensive copies, and deterministic complete-cohort eligibility only. It
writes no file, migrates no store, activates no archive, and routes no retained lookup. The passive
F2 format correction now freezes scope-separated global/segment indexes, typed hot/archive
locations/counts, a distinct migration-genesis checkpoint, and generated answers. The stateful F2
review stopped before v2 bytes because the corrected projection cannot represent a valid v1
committed attempt with no lifecycle record without inventing state or violating exact counts. The
[F2 v1 mapping blocker](SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md) and its executable witness are
retained; migration/full verification now waits for another passive contract decision. The design
deliberately leaves
referenced-history deletion, implementation/installed validation of the selected owner lock and
power loss, coherent rollback prevention, continuous service, consumers, and guests blocked.

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
| PR #58, governed runtime work plan | `37482fb` | Recorded the real-fork governance and release/admission sequence for governed `deno_core` without admitting a profile. |
| PR #59, durable lifecycle E4 | `68b75fd` | Replaced the active memory driver with the colocated v1 store, added stable fake effect/instance identity and injected startup coordination, and retained the death/fault/reopen matrix without a consumer or guest. |
| PR #60, durable lifecycle E5 | `d88661b` | Retained exact active/retained lifecycle capacities, destroyed-only release, repeated startup, and bounded recovery-exhaustion evidence at `local-mechanic` status. |
| PR #61, P0-3 cross-language/console | `91c4aa4` | Added independent Node verification of 43 vectors and local process-fault evidence while retaining the real transport/guest and console-coverage gaps. |
| PR #62, fork-native runtime bundle | `4f1edd7` | Stopped the first governed Linux/arm64 construction before build because the merged `rusty_v8` publication contract supported only Linux/amd64. |
| PR #63, libkrun patch provenance | `ada09cd` | Restored the exact governed base-patch bytes and made their canonical hashes unconditional verifier inputs without changing patch semantics. |
| External `dills122/libkrun` PR #2, governed console/raw-FD follow-up | `cf0333cdba478cc34a8570a65b38412da7fd3ecc` | Merged exact follow-up head `8a2c91943793668f31a1cf7af431933be935bb58`, preserving the five-patch aggregate while adding bounded library tests, two source fixes, measured coverage, and CI build evidence. This is an external fork source checkpoint, not Capsule product admission. |
| PR #64, security overview | `030668857a5c0fbba76f5a81b905ae03aed63d28` | Source head `df0a48e9518db67356d8645da3f58f83f6692662`; added a plain-language explanation of Capsule's intended security design and current claim boundary. |
| PR #65, authenticated IPC topology | `758bed7eb0f85124b4477c37f8893401f8c2c037` | Source head `b0acd66ae8fd06f2de8522723469ca5ac9fb01de`; proposed one unprivileged per-user Supervisor with a native XPC/Security front end, the existing Go core, two role-specific services, and four closed calls; no product IPC was implemented. |
| PR #66, passive TypeScript approved-byte contract | `4be6d0e5059f179a66b6b439de76d6f8ed4e2f08` | Source head `d2f6969122d840e045fc63bb61d6525b95f8da22`; added nine known answers, 14 refusal mutations, and the illustrative plan-v1 source roles without selecting an owner, activating a consumer, or admitting a runtime. |
| PR #67, related-systems/current-plan alignment | `5dee741a1aeb465dc5eb77917089d98cf1ad4564` | Source head `ad6c89c9f429fb93540ec711d5c821a3fbf835cf`; retained the field-authority requirement and reconciled then-current runtime, lifecycle, and IPC ordering without adding evidence. |
| PR #68, authenticated IPC S1 consistency stop | `77e39fdb83cee22243c711507318327777dcb607` | Source head `c3ad12c6c0aa40684268b7f5154d5f87a8183b85`; selected the versioned atomic-cutover dependency and generated zero S1 fixtures/known answers rather than freezing the v0-only 562-byte record or 626-byte arithmetic. |
| PR #69, passive field authority | `ba743febfe77a8a96cf368f3513fbf5795effb3c` | Source head `0b7e5a911e66a960677e66cbafc3ab349ccdc11d`; added fail-closed verification for 164 top-level fields across 15 passive targets. It does not recursively classify future Source Preparer or plan-v1 objects. |
| PR #70, Supervisor archive design | `66f7fc1120eeb0bc7f4b15e955ca07988f500ddb` | Source head `2217f297fdec08d48890418899cc25ed5f176c66`; proposed ADR-0031 and the F1-F6 plan, with no archive implementation in that merge. |
| PR #71, Source Preparer topology | `c68dfb1535b6763ad7c89d5f401fa9002f225b26` | Source head `5edc7fd90025c918291b5967ae0f06297b72540e`; proposed the separate enrolled owner/store topology and P0-P7 plan without implementing it. |
| PR #72, Source Preparer P0 authority review | `2e268b01d4174fe90397c00abc5973a3dd785606` | Source head `a12041c36d90815474598f0929c595b32dc68e11`; retained P1 HOLD/NO-GO pending protected-store, worker, bootstrap/update, retention/release, recursive field-authority, and lifecycle evidence. No P1 contract exists. |
| PR #73, governed libkrun reconciliation | `f6fcf172af752a425afb29ce62680d0b115f6998` | Source head `5e17ac8cec21320c3693049c53e7575bb9dbc15a`; reconciled external fork PR #2 head `8a2c91943793668f31a1cf7af431933be935bb58` and merge `cf0333cdba478cc34a8570a65b38412da7fd3ecc`, two lifecycle fixes, and 37/88-function plus 298/733-line four-file coverage without guest/backend/profile admission. |
| PR #74, Supervisor owner-lock boundary | `e930f9dbd877bea0cbd55870060f48c9c7fdd72f` | Final reviewed source head `afd148c92f4b9f6f35f2a7d9161502cd1175a107`; proposed ADR-0033 selects enrolled pre-created inode plus lifetime nonblocking BSD `flock` for cooperating Supervisor ownership only. It supplies neither same-UID containment nor Source Preparer storage protection; G1 was the next slice at that merge. |
| PR #75, passive archive F1 | `6fc31a049c476acf5085071c48d3d5e36f27240f` | Source head `20c8d7df1d9ed3eb009e8ce9a0afbd41e03807ef`; implemented passive archive types, limits, known answers, defensive copies, and eligibility only. No file write, migration, activation, lookup consumer, or adapter call exists; F2 is next. |

The ecosystem reuse audit is retained in
[`ECOSYSTEM_REUSE_AND_ADOPTION.md`](ECOSYSTEM_REUSE_AND_ADOPTION.md) on its delivery branch until
its draft PR merges; its final merge commit must be added to this table during integration. It is
a primary-source-backed planning map, not implementation or security evidence. It records three
bounded decision lanes: production CBOR/COSE profiling, SQLite comparison after F2/G2, and the
already-active `.mjs` parser boundary. It explicitly does not preempt M1, G2, F2, or the separate
ARM64 `rusty_v8` work.

The merge commits, not former draft-PR state or chat handoffs, are the integration checkpoints.
PR #75 merged before PR #74 finalized, so main's first-parent order is `... f6fcf17 -> 6fc31a0 ->
e930f9d` even though the PR numbers are 73, 75, and 74. The exact source heads above preserve review
identity without treating an unmerged head as durable Capsule evidence.

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
- durable-lifecycle Slices E1-E5: passive runtime-neutral types, explicit fixed-store v1 migration
  and validation, colocated transactions, the no-guest fake driver, exact 256-active/4,096-retained
  ceilings, destroyed-only capacity release, and repeated-startup/exhaustion evidence. No consumer,
  real adapter, runtime, backend process, or guest is present; and
- passive field-authority verification for 172 fields across 16 selected pre-freeze targets,
  including nested MJS SourceManifest member coverage without claiming recursive coverage for
  future Source Preparer or plan-v1 objects; and
- archive Slice F1 passive projections, limits/known answers, defensive copies, and deterministic
  eligibility selection. No file/store/archive behavior or consumer exists; and
- Proposed ADR-0033's local owner-lock mechanism selection, passive G1 Go/Darwin acquisition, and
  bounded G2 current-v1/no-guest startup composition, while signed bootstrap and installed
  protected-root evidence remain unimplemented; and
- governed `deno_core` physical omission, same-host package reproduction, exact V8 closure NO-GO,
  standalone dynamic-root evidence, and the fork-native Linux/arm64 blocker. Accepted ADR-0028
  selects its engineering order without admitting a profile; the real Deno and `rusty_v8`
  governed baseline branches are merged, but no governed arm64 release or admitted artifact exists.
- exact public governed libkrun source identity through merge
  `cf0333cdba478cc34a8570a65b38412da7fd3ecc`, with the unchanged five-patch aggregate, bounded
  console/raw-FD library evidence, two local lifecycle fixes, and improved coverage. The queued
  Linux-arm64 unit job, post-merge branch/verifier mismatch, absent independent review, remaining
  uncovered code, release obligations, and every guest/product admission boundary remain explicit.

Current dependency and priority view:

1. **Blocked:** archive F2 migration/full verification waits for a reviewed representation of the
   valid v1 committed-attempt-without-lifecycle world; do not choose v2 bytes or silently restrict
   migration before that passive decision.
2. **Independently actionable now:** retain bounded owner-lock G2 current-v1/no-guest composition
   while planning its separate installed G3 matrix; follow accepted ADR-0034 for `.mjs` M1/S1
   passive contract work; and maintain exact documentation plus recursive field-authority design.
3. **Waiting:** the fork-native runtime bundle waits for an accepted successful Linux/arm64
   `rusty_v8` source/artifact handoff. External PR #4 is open at exact head
   `aa921fa48901bf28774d61248b0187c8b91c55a4`; passing contract jobs and in-progress full builds are
   not durable Capsule evidence and no workflow artifact may be reused yet.
4. **Later composition:** governed runtime plus libkrun requires admitted artifacts, the remaining
   transport/launcher/root/device/teardown work, and explicit authorization for an owned disposable
   development guest. No current task authorizes a guest.
5. **Credential/environment dependent:** Apple Development identities and provisioning profiles
   must be deliberately authorized before the existing installed matrices run. Current Individual
   membership is Team `W4QUR9FUL4`, and read-only discovery reports a valid Apple Development
   identity for that Team. Local signed/provisioned experiments can proceed once exact W4 role
   identifiers, entitlements, and profiles are deliberately created. All three Xcode 26.6-cached
   profiles belong to historical Team `3DDR84M4JS` and are not reusable for W4 tests. A separate
   Developer ID Application identity for historical Team `3DDR84M4JS` is later distribution
   authority requiring explicit authorization and matching-Team package design; it is not W4
   development evidence and does not make Developer ID/notarization work current. Paid owned
   clean-host/minimum-OS coverage is not currently planned
   and remains deferred activation/distribution evidence, not a blocker for F2 migration/full
   verification or owner-lock G3 planning.
   A genuinely independent
   Linux/arm64 builder is viable but not currently planned; same-host/GitHub-CI equality remains
   limited and independent-builder equality is deferred.
5. **Next passive contract path:** accepted ADR-0034 removes Source Preparer/plan-v1 from the
   first-release critical path. M1 narrows the proposal/source/manifest to one byte-exact
   `main.mjs`; S1/M2 then generates the revised `RegisterPlanV0`/fetch projections and exact caps
   from complete field authority. No S1 fixtures exist and no product endpoint is authorized.

TypeScript remains conditional. If later selected, Source Preparer P0A and ADR-0030's atomic
plan-v1/RegisterPlanV1 cutover still apply with no dual active v0/v1 acceptance. CommonJS, package
resolution, legacy Node module surface, and runtime-contract widening remain forbidden.

The idle governed Deno fork has transferred from `dills122/deno` to
[`Shrimpworks/deno`](https://github.com/Shrimpworks/deno) with its fork relationship, default and
`capsule/upstream-v2.9.4` branches, merged PR #1, and Actions history intact. This is repository-
governance state, not source review, build provenance, security evidence, release admission, or
profile validation. The governed libkrun fork has also transferred from `dills122/libkrun` to
[`Shrimpworks/libkrun`](https://github.com/Shrimpworks/libkrun), still a public fork of
`libkrun/libkrun`. PRs #1 and #2 retain their history; `capsule/upstream-v1.19.4` points to exact
merge `cf0333cdba478cc34a8570a65b38412da7fd3ecc`, with verified ancestry from upstream anchor
`728df8125077d0db44265f6e997c72b81b65c015`. This is likewise governance only. Integration must pin
exact commits/digests and verify ancestry rather than consume the movable baseline branch.
`capsule-corp` and `rusty_v8` deliberately remain at their current owners while active work
continues.

Proposed ADR-0029 selects an IPC topology but does not implement or validate its native bridge,
installed endpoints, peer identities, or production transport. This checkpoint also does not decide
stateful archive migration/activation or production-engine archive/compaction beyond passive F1,
production COSE/Keychain/user-presence signing, consumer
ownership, evidence composition, or public cutover. The authority/lifecycle snapshot lacks real
multi-process locking and rollback-resistant identifier/nonce/effect uniqueness. The fixed snapshot
is durable for controlled local tests, but ownership remains in-process and no production
persistence claim follows. Content, evidence, runtime, real backend, and guest remain absent from
the unwired path.

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
5. record the merge commit that integrated the result; and
6. for a dependency or new custom primitive, cite the matching reuse-map row and retain the
   completed dependency-policy checklist with the consuming task or decision evidence.

No task is complete for coordination purposes if its only useful output exists in chat history.
