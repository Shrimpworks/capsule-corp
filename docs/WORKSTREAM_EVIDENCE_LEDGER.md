# Workstream and evidence ledger

Date: 2026-08-05

Status: durable coordination index. This ledger records where completed task conclusions were
integrated; it is not independent security evidence, a posture promotion, or a replacement for the
linked experiment results, synthesis documents, ADRs, or conformance fixtures.

Current work progress is summarized in the
[canonical status dashboard](STATUS_LANGUAGE.md#current-workstream-dashboard). Historical result
names in this ledger are provenance, not today's work status. Read each completed slice separately
from its parent workstream, ADR lifecycle, control-evidence state, and product admission.

The consolidated outcomes, exact known answers, application status, combined verification, and
next dependency boundary for the latest group are recorded in the
[Phase 2B and Gate C parallel-task checkpoint](PHASE_2B_GATE_C_TASK_GROUP_CHECKPOINT.md).

## 2026-08-05 merged-state reconciliation

| Slice | Scoped result | Parent boundary |
| --- | --- | --- |
| Passive authority-plane atomic cutover | `PASSED`: exact-one-`main.mjs` plan, bindings, manifest, source, registration, and defensive fetch publish atomically in the fixed-store oracle. | Authenticated product facades, Broker/approval, attempt, runtime, backend, and guest remain `BLOCKED`. |
| Broker rendering and approval conformance v0 | `PASSED`: Accepted ADR-0043 freezes a read-only ASCII-safe projection over bound Supervisor bytes, the exact Secure Enclave/user-presence/no-fallback key contract, and a Capsule-owned public-key-only ApprovalGrant Sign1 verifier with canonical/header/payload/signature/key/role/time/refusal tests. | Inline-input content is not present or shown. Installed UI, Keychain/LocalAuthentication, private key/signing, authenticated IPC, verifier/store wiring, activation, runtime, backend, guest, and product admission remain `BLOCKED`. |
| C2B v3 passive successor | `PASSED`: accepted fork identities and exact runner/libkrunfw/kernel, descriptor/device/resource/teardown semantics are bound in an 18,357-byte zero-effect fixture with 128 field classifications. | Current-source libkrun dylib, final runner, fixed-owned-guest eligibility, and admission remain `BLOCKED`. |
| C2B v4 build/static materialization | `PASSED`: exact accepted header, twice-reproduced current-source unsigned libkrun dylib, independent ABI audit, byte-equal unsigned final runner, and composed digest are retained with no execution. | Fixed-owned-guest eligibility is `BLOCKED` on separate authorization naming the exact v4 digest; admission remains `BLOCKED`. |
| Production-shaped I2B1 CBOR/COSE wrapper review | `PASSED`: 95 Go/Swift checked-in-vector cases close pairing, repeated-field binding, payload-owned replay, caps, mutations, and fuzz targets. | Live signing, caller/key authorization, durable replay, installed consumers, and product admission remain `BLOCKED`. |
| Compiled-artifact archive migration | `PASSED`: 210 completed payload/evidence files are pinned to capsule-experiments commit `0944ffd8cfd01ec23e4ae99138b0931d56804077`; Capsule retains compact conformance metadata and six deterministic I2B2 source inputs. | No product, signing, installation, runtime, or admission status changes. |
| Governed fork promotions | `PASSED`: Deno r3, `rusty_v8` r5, and libkrun r3 accepted lines are locked defaults with historical refs protected. | Governed runtime engineering is `IN_PROGRESS — TRENDING_GOOD`; release and admission remain `BLOCKED`. |
| Source Validator R3/R4 | R3 `PASSED` its exact Apple Development signed, installed, inactive-policy composition. Exact R4-v1 candidates are `NO_GO`; R4-v2 was not executed. | Product R4/R5 remains `BLOCKED` and post-alpha defense-in-depth. |
| Archive F5 | `PASSED`: owner-held coherent backup, complete-copy verification, read-only exact-anchor restore admission, known-orphan cleanup, and offline reporting. | F6, restore activation, production-engine, power-loss/rollback, and external-alpha continuity remain `BLOCKED`. |

The current generated conformance totals are 95 rules, 502 cases, and 624 fixtures. The current
field-authority manifest has 1,159 fields across 91 profiles and 57 targets. These whole-repository
totals supersede older checkpoint counts without rewriting their historical evidence.

## 2026-08-05 internal-alpha architecture audit synthesis

Five defensive read-only audits independently reviewed the complete architecture/critical path,
governed runtime/guest composition, installed macOS and Source Validator topology,
protocol/approval/client flow, and persistence/recovery posture at default-branch commit
`23e48e242765eea219cdd80724b865638fc02200`. Each audit is `PASSED` for its exact decision question;
the internal-alpha parent is `IN_PROGRESS — TRENDING_GOOD` and product admission remains `BLOCKED`.

Accepted [ADR-0040](adr/0040-freeze-owner-only-internal-alpha-posture.md) and the retained
[synthesis](ALPHA_ARCHITECTURE_AND_RELEASE_AUDIT.md) record the reconciled outcome: one exact
`main.mjs`, bounded inline JSON, a native human Broker, one fresh governed guest per `AttemptID`, a
fixed benign guest checkpoint before product alpha, runtime no-loader/host-authority enforcement,
Source Validator deferred as post-alpha defense-in-depth, and a narrowly bounded owner-only fixed-
store exception before F6. No audit created a guest, signing/install action, product key, service,
store consumer, runtime/profile admission, or external mutation.

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
| Raw storage, scratch, and controlled egress | `019fb9e7-25ec-78b2-8bf2-6558a2a9f250` | Conditional development-feasibility pass; guest read-only and fixed scratch bounds worked, but same-user host mutation defeated immutable custody and the spike extractor is not a product parser. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-storage-egress/RESULTS.md) and Gate C synthesis |
| Console, timeout, cancellation, CPU, and memory | `019fb9e7-25ec-78b2-8bf2-659ad6ccdef1` | Conditional pass for bounded prefixes, external wall/cancel scheduling, exact forced teardown, and closed vCPU/RAM profiles. Runner exit zero, graceful-only shutdown, CPU quota, and exact host-memory claims were rejected. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-console-lifecycle/RESULTS.md) and Gate C synthesis |
| Installed lifecycle and crash recovery | `019fb9e7-25ec-78b2-8bf2-65688a86f41f` | Conditional same-host pass for enrolled launch mechanics and exact recovery. Complete installed authority, accepted distribution, clean-host, session, reboot, and minimum-OS evidence remain open. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-installed-recovery/RESULTS.md) and Gate C synthesis |
| Adversarial VMM and cross-job isolation | original `019fb9e7-2692-7cb0-a5c2-cebb9378e07f`; completed replacement `019fba17-9659-74f1-ab6f-82bfb72bc991` | The original task failed at the task/app layer. The bounded replacement verified and preserved the existing report. The exact profile conditionally failed because `NullFs` remains guest-visible even without a configured host directory. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-adversarial/RESULTS.md), [`VALIDATION_RECEIPT.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-adversarial/VALIDATION_RECEIPT.md), and Gate C synthesis |
| Runtime packaging, provenance, patch governance, and supply chain | `019fb9e7-27b1-7c42-9903-8be99f620602` | Conditional build/release feasibility pass; current bytes are no-go for admission. Controlled same-host equality is not two-builder provenance, and the complete final signed/notarized topology is not available. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-supply-chain/RESULTS.md), [`HANDOFF.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-supply-chain/HANDOFF.md), and Gate C synthesis |
| P0-0 exact stock Bun runtime authority | `019fc2e6-cf9f-7482-82b5-240992c79419`, delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | Failed for stock Bun 1.3.14: process, `execve`, FFI/native-loader, inspector, Worker, and descriptor authority remain reachable. `RUNTIME-001` stays unsupported; only a governed patched/external branch or alternate runtime remains. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-bun-runtime-authority/RESULTS.md), source inventory, synthetic probes, and selected evidence |
| P0-0 governed Bun construction | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` after merged PR #33 | NO-GO under the explicit broad/unreviewable stop rule. Exact review found a conservative 40-hand-authored plus 10-generated-output minimum before a patch or governed build; a narrow process/exec self-seal cannot independently close Worker or native loading. Alternate-runtime investigation and an ADR-0003 superseding decision are required. | [`CONSTRUCTION_REVIEW.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-bun-runtime-authority/governed-closure/CONSTRUCTION_REVIEW.md), reproducible inventory, exact source hashes, and selected evidence |
| P0-0 governed `deno_core` physical omission | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` after PRs #39 and #40; merged in PR #43 | Current translation: scoped slice `PASSED`; parent governed-runtime work `IN_PROGRESS — TRENDING_GOOD`; product admission blocked and `RUNTIME-001` unsupported. Historically recorded as `PHYSICAL-OMISSION-PASS; NO RUNTIME ADMISSION`. A one-file patch reduced the exact built-in registry from 99 ops to three bootstrap-required ops; ASLR-controlled clean builds reproduced the snapshot and final binary, and restoration tests failed closed. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-core-physical-omission/RESULTS.md), Phase A inventory, exact patch, fixed probe, and selected evidence |
| P0-0 governed `deno_core` reproducible package | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` after merged PR #43 | Current translation: scoped same-host reproduction `PASSED`; parent governed-runtime work `IN_PROGRESS — TRENDING_GOOD`; runtime admission `BLOCKED`. A digest-pinned no-apt builder and complete 191-crate offline source bundle reproduced the prior snapshot and binary across two clean same-host containers. Independent-builder provenance, archive-corresponding V8 source/notices, a standalone dynamic runtime root, and production source ownership remained outside this passed scope; `RUNTIME-001` stays unsupported. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-core-reproducible-package/RESULTS.md), bundle manifest, CycloneDX SBOM, source/license inventory, unsigned provenance, and selected evidence |
| P0-0 governed `deno_core` V8 source/license closure | `019fc58f-90fd-78c2-88f8-20c7a1f38359`, delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` after merged PR #50 | The exact official prebuilt-asset publication route is `NO_GO`; replacement governed-fork work is `IN_PROGRESS — TRENDING_GOOD`. Exact asset/job, rusty_v8/V8 revisions, 20 gitlinks, Chromium V8 base, four-patch stack, 1,875 archive members, and 726 license/notice candidates are retained. Mutable publisher inputs, missing GN/Ninja link closure, and absent generated notices prevent exact rebuilding and complete notice/SBOM closure from that asset. `RUNTIME-001` stays unsupported. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-v8-source-license-closure/RESULTS.md), [`SOURCE_PUBLICATION.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-v8-source-license-closure/SOURCE_PUBLICATION.md), exact manifests, archive inventory, and fail-closed checklist |
| P0-0 governed `deno_core` self-contained runtime root | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` after merged PR #50 | Current translation: standalone dynamic-root slice `PASSED`; parent governed-runtime work `IN_PROGRESS — TRENDING_GOOD`; runtime admission `BLOCKED`. Exact Debian snapshot packages produced a reproducible 22-entry root for the unchanged binary/snapshot; empty-environment scratch execution and file-open/mutation evidence used no ambient Bookworm library/config/NSS/locale/timezone/package/cache bytes. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-core-runtime-root/RESULTS.md), closed root/package manifests, ELF/file-open/mutation evidence, SBOM, and source/notice mapping |
| P0-0 fork-native governed runtime bundle | `019fc8f3-44d1-77b1-86c6-8a38e2b47211`, delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc`; completed by merged [`Shrimpworks/capsule-experiments` PR #1](https://github.com/Shrimpworks/capsule-experiments/pull/1) | Scoped exact Linux/arm64 construction `PASSED`; parent governed-runtime work `IN_PROGRESS — TRENDING_GOOD`; runtime/profile admission `BLOCKED`. Merge `fa03d7043b4f0653081d6c5733d597f49f6efd1c` reconstructs the governed `rusty_v8`, `deno_core` binary/snapshot/two-file bundle, and 22-entry root identities from empty outputs under the retained method. No guest, release, signature, or admission follows. | Exact merged handoff and retained construction manifests, checksums, SBOM/provenance/source/notice evidence, and offline verifier |
| P0-0 governed `rusty_v8` Linux/arm64 builder | External governed-fork work merged in [`Shrimpworks/rusty_v8` PR #4](https://github.com/Shrimpworks/rusty_v8/pull/4) | Scoped fork construction `PASSED`; parent governed-runtime work `IN_PROGRESS — TRENDING_GOOD`; runtime/profile admission `BLOCKED`. Exact head `80e863ddb942a4aa2b384e794fc23e35b9d2bb15` and merge `cbf56de2e1156b1cf1561fdbaea7172a0aa056f4` retain the Linux/arm64 builder and GN-query correction. Exact-head workflow-dispatch run `30925045754` passed network-disabled build/test/evidence collection and uploaded an unsigned bundle. No governed release, accepted Capsule bundle, independent-builder equality, or admission follows. | [Successful exact-head ARM64 workflow](https://github.com/Shrimpworks/rusty_v8/actions/runs/30925045754), governed fork source, and merged PR history |
| Governed Deno and `rusty_v8` fork acceptance | User-visible defensive governance/integration task after merged Capsule PRs #181/#184, [`Shrimpworks/deno` PR #3](https://github.com/Shrimpworks/deno/pull/3), and [`Shrimpworks/rusty_v8` PR #5](https://github.com/Shrimpworks/rusty_v8/pull/5) | Scoped fork-governance transition `PASSED`; parent governed-runtime work remains `IN_PROGRESS — TRENDING_GOOD`; runtime/profile admission remains `BLOCKED`. Deno accepted/default ref `capsule/accepted-v2.9.4-r3` is locked at verified merge `3fa21d1ae7705ab4bcb4bc98955f25301b20122a`, tree `6060cb0eb4cd3395a4c141f054634968744617d2`, descended from anchor `14eea3160ae5834476aa3b9d317b8d41d991b982` and preceding accepted `4cce46bafccd0df9d1709cf406cd03c05b5daa0b`. `rusty_v8` accepted/default ref `capsule/accepted-v150.2.0-r5` is locked at verified merge `d09221062280ae1675fe26c53c3f43871aae2055`, tree `2632901e6e7e9ac88662756ceb658d4e3e49fceb`, descended from anchor `d305e6afa7736f6e298c30ae6646f7709ee9382b` and preceding accepted `cbf56de2e1156b1cf1561fdbaea7172a0aa056f4`. Both completed review targets and all historical refs remain protected; each `main` remains protected, mutable upstream-integration state. No release, artifact publication, runtime wiring, or admission follows. | Exact merged PR bodies and refs; GitHub default/ref/protection readback; [repository fork ledger and settings](REPOSITORY_SETUP.md#accepted-governed-fork-readback) |
| Governed `deno_core` C1 passive composition | User-visible implementation task delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | Scoped passive contract `PASSED`; parent governed-runtime work `IN_PROGRESS — TRENDING_GOOD`; runtime/profile admission `BLOCKED`. Exact-byte schema/fixture plus independent Go and Node validators freeze the app JSON-in/JSON-out behavior, desired globals/ops/modules/files, loader absence, logical descriptors, resource references, all eight construction identities, and refusal/stops with every authority effect false. C2 owns numeric descriptors, machine resources, guest execution, and composed evidence. | [C1 contract and limitations](protocol/GOVERNED_DENO_CORE_C1_COMPOSITION.md), `schemas/conformance/c1-governed-deno-core/controlled-development-profile.json`, and independent validators |
| Governed `deno_core` C2A passive execution-profile preparation | User-visible implementation/research task delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` after merged PR #155 | Scoped passive preparation `PASSED`; parent governed-runtime work `IN_PROGRESS — TRENDING_GOOD`; C2B and runtime/profile admission `BLOCKED`; `RUNTIME-001` and `VMM-001` unsupported. Exact-byte schema/fixture plus independent Go and Node validators consume unchanged C1 bytes and freeze plan roles, numeric descriptors, candidate machine/transport/teardown values, known-answer bytes, null final-artifact/runtime-profile/resource blockers, 91 C2B cases, 18 restoration mutations, and all-zero authority effects. No guest, process, runtime, credential, release, or admission was created. | [C2A contract, matrix, blockers, and limitations](protocol/GOVERNED_DENO_CORE_C2A_EXECUTION_PROFILE.md), `schemas/conformance/c2a-governed-deno-core/passive-execution-profile.json`, and independent validators |
| Governed `deno_core` C2B fixed-fixture passive binding v1 | User-visible implementation/documentation task delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | Scoped passive reconciliation `PASSED`; parent governed-runtime work `IN_PROGRESS — TRENDING_GOOD`; C2B composed-profile/guest execution and runtime/profile admission `BLOCKED`; `RUNTIME-001` and `VMM-001` unsupported. The immutable 8,221-byte v1 object preserves exact C1/C2A bytes and its historical Deno PR #2 / experiments PR #3 dependency checkpoint, fixed known answers, six artifacts, 143 recursive field-authority paths, and zero effects. Both PRs later merged; v1 bytes were not reinterpreted. | [C2B v1 passive binding](protocol/GOVERNED_DENO_CORE_C2B_PASSIVE_BINDING.md), exact fixture, mirrored historical inputs, and independent Go/TypeScript validators |
| Governed `deno_core` C2B no-guest build-closure passive binding v2 | User-visible defensive implementation/conformance task delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | Scoped no-guest artifact closure `PASSED`; parent governed-runtime work `IN_PROGRESS — TRENDING_GOOD`; composed-profile/guest execution and runtime/profile admission `BLOCKED`; `RUNTIME-001` and `VMM-001` unsupported. The immutable 7,115-byte v2 successor pins capsule-experiments PR #4 merge `50108417...9972`, reviewed head `518eea04...a76`, six libkrun/kernel/init/launcher/root identities, the 8,845-byte unadmitted manifest candidate, and the build-only non-final preflight harness. Final runner, separate firmware, composed profile, CPU-time/host-memory/scratch limits, guest evidence, and admission are null/refusing. Independent Go/TypeScript decoders, schema mutations, defensive copies, cap+1/cross-version checks, 120 recursive field classifications, and no-consumer structure retain zero effects. | [C2B v2 passive binding](protocol/GOVERNED_DENO_CORE_C2B_PASSIVE_BINDING_V2.md), exact schema/fixture, merged archive evidence, and independent validators |
| Governed `deno_core` C2B fixed-owned-guest passive successor v3 | User-visible defensive implementation/conformance task delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | Scoped passive successor contract `PASSED`; parent governed-runtime work `IN_PROGRESS — TRENDING_GOOD`; fixed-owned-guest eligibility and runtime/profile admission `BLOCKED`; `RUNTIME-001` and `VMM-001` unsupported. The immutable 18,357-byte v3 fixture and `8b1ec936...bcc1` composed passive-contract digest preserve C1/C2A/C2B v1/v2, bind current accepted Deno/`rusty_v8`/libkrun commits and trees, resolve the per-attempt runner and libkrunfw/kernel roles, freeze exact FD/port/virtio/runtime/resource/teardown semantics, omit unsupported resource fields, and contain no JSON null. The retained libkrun dylib predates accepted runtime source and the preflight is not a runner; both are historical evidence only. No consumer, process, VM, guest, signing, release, or admission effect exists. | [C2B v3 passive successor](protocol/GOVERNED_DENO_CORE_C2B_PASSIVE_BINDING_V3.md), exact schema/fixture, Accepted ADR-0041, and independent Go/TypeScript validators |
| C2B host-runner passive source contract v1 | User-visible defensive/local-only implementation task delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | Scoped passive/static source boundary `PASSED`; parent governed-runtime work remains `IN_PROGRESS — TRENDING_GOOD`; fixed-owned-guest eligibility and runtime/profile admission remain `BLOCKED`; `RUNTIME-001` and `VMM-001` unsupported. Exact 3,996-byte C17 source-contract bytes and dependency-free Go verifier freeze one Supervisor-owned runner per AttemptID, sealed inherited FDs 0–7, raw-root FD 4, exactly three console ports, explicit implicit-console/init/vsock disable calls, exact fail-closed order, forbidden path/network/vsock/virtiofs authority, and Supervisor teardown/absence. Mutations refuse. Accepted ABI header/current-source dylib inputs are not retained locally; no final runner build, libkrun/HVF call, process, guest, signature, installation, materialized profile, or admission exists. | [Host-runner source contract and blockers](protocol/GOVERNED_DENO_CORE_C2B_HOST_RUNNER_SOURCE_V1.md), `schemas/conformance/c2b-host-runner-source-v1`, and `internal/execution/hostrunnerpassive` |
| C2B materialized host-runner profile v4 | User-visible defensive/local-only implementation task delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | Scoped build/static materialization `PASSED`; parent governed-runtime work remains `IN_PROGRESS — TRENDING_GOOD`; fixed-owned-guest eligibility and runtime/profile admission remain `BLOCKED`; `RUNTIME-001` and `VMM-001` unsupported. Exact accepted libkrun commit `7432eda5...d632`/tree `7671440c...e346`, 54,658-byte header, twice-reproduced 4,393,448-byte unsigned current-source dylib, 2,512-byte ABI audit, 7,917-byte final source, 100,488-byte unsigned runner, and composed digest `e390085c...ba82` are retained. Verification closes FDs 0–7, close-from 8, three ordered ports, disable calls, call/import closure, no replacement authority, libkrunfw-only boot, no separate firmware, and external teardown. No runner/libkrun execution, HVF, VM, guest, signing, installation, release publication, consumer, or admission occurred. | [C2B v4 materialization, exact known answers, and stop boundary](protocol/GOVERNED_DENO_CORE_C2B_MATERIALIZED_PROFILE_V4.md), `schemas/conformance/c2b-host-runner-materialized-v4`, and `internal/execution/hostrunnermaterialized` |
| P0-0 TypeScript approved-byte boundary | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` from merged checkpoint `5448943` | Passed only the exact pre-approval byte-ownership question for a strip-only ESM TypeScript subset. Exact Node 22.22.1/Amaro 1.1.5 emission was deterministic for fixed fixtures and mutation/cap/diagnostic cases failed closed. Proposed ADR-0026 binds original and executable source roles before registration. No current contract changed, no owner/runtime was selected, and `RUNTIME-001` remains unsupported. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/typescript-approved-byte-boundary/RESULTS.md), fixed Node and `deno_ast` probes, independent Go verifier, and selected evidence |
| ADR-0034 M1 source-byte/manifest foundation and scanner disposition | `019fca55-554d-7b23-bb2e-6d2c4b0e16cd`, delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | Passive byte/manifest foundation `PASSED`. The handwritten scanner is `NO_GO` and was removed after a retained division-versus-regexp counterexample exposed a live-`import()` false negative. The replacement parser/process work proceeded separately; no runtime/no-loader claim follows. | [`MJS_MODULE_REQUEST_VALIDATOR_HOLD.md`](MJS_MODULE_REQUEST_VALIDATOR_HOLD.md), generated MJS conformance fixtures, SourceManifest CDDL, and field-authority target |
| `.mjs` parser/process boundary | `019fc2de-552d-77a0-aa47-35ac39d02edc` user-visible defensive/local-only research delegation | Scoped parser/process selection `PASSED`; later product Source Validator work is now `BLOCKED` by the separately retained V1 enrollment and V2 resource/confinement stops. Exact Oxc 0.140.0 matched all 33 local parse-only cases and all 28 canonical M1 outcomes. Accepted ADR-0035 selects a future one-shot Source Validator outside daemon/Broker/Supervisor address spaces; ADR acceptance records architecture, not product admission. Product artifact enrollment, supported confinement, parent integration, and runtime no-loader evidence remain. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/mjs-parser-boundary/RESULTS.md), M1 mapping, exact lock/fixtures/classifications, supply inventory, measurements, and [`HANDOFF.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/mjs-parser-boundary/HANDOFF.md) |
| ADR-0035 V0 passive Source Validator contract | `019fc2de-552d-77a0-aa47-35ac39d02edc` user-visible defensive/local-only implementation delegation | Passive byte-contract observation only: exact request/result/engineering-candidate/artifact-profile fixed frames, caps, domain separation, closed status/count relationships, 128 generated cases, zero-effect refusals, and independent unwired Go plus test-only Rust agreement are retained. The generator consumes M1 bytes in place. No Oxc product dependency, parser invocation, child, sandbox, endpoint, artifact enrollment, consumer, Approval/key effect, runtime, backend, or guest exists; accepted ADR-0035 does not change this evidence boundary. | [`MJS_SOURCE_VALIDATOR_PASSIVE_CONTRACT.md`](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_CONTRACT.md), generated conformance corpus, Go passive package, and Rust test-only oracle |
| ADR-0035 V1 Source Validator artifact | `019fc2de-552d-7b23-bb2e-6d2c4b0e16cd` separate user-visible defensive/local-only implementation delegation | Current translation: the bounded unwired V1 artifact slice `PASSED`; product enrollment is `BLOCKED`, not `NO_GO`. Exact Oxc 0.140.0 plus Rust 1.95.0 and pinned `sha2` produces a 1,146,656-byte macOS arm64 executable at SHA-256 `ba2a6b38...cb600`. The 74-dependency lock, registry/source hashes, licenses/notices, CycloneDX SBOM, unsigned in-toto statement, V0/M1 deterministic results, mutations, and same-host two-directory reproduction are retained. The actual V0 profile binds executable/build/assessment but is explicitly not enrolled. Independent builder/clean host, installation-authority signatures, vulnerability owner/SLA, V2 confinement, V3/V4 consumers, V5 expansion, V6 runtime, backend, and guest remain absent. | [Pinned V1 archive](https://github.com/Shrimpworks/capsule-experiments/tree/0944ffd8cfd01ec23e4ae99138b0931d56804077/experiments/completed-compiled-artifact-payloads/payloads/capsule-corp/artifacts/mjs-source-validator-v1), exact binary/profile, build/source/license/SBOM/provenance/reproduction evidence, Rust tests, and process interruption oracle |
| ADR-0035 V2 Source Validator process profile | `019fcd7b-b0d1-74b1-9eb9-4dc1db12509d` separate user-visible defensive/local-only implementation/research delegation | Exact V2 candidate `BLOCKED`. The strict test-only macOS bootstrap fixes target/profile argv, empty environment, cwd, post-exec FD inventory, CPU/file/FD/child/output/wall limits, group kill/reap, and fault refusal, but `RLIMIT_AS` returns `EINVAL` before V1 `exec`. An explicitly unbounded diagnostic mutation verifies ordinary/exact-maximum frames, partial/duplicate/trailing/oversize/crash/signal/timeout/cancel and clean-later behavior while proving owned out-of-cwd reads, IPv4/Unix socket creation without connect, cwd metadata/empty-file writes, and a 512 MiB mapping remain. Supported App Sandbox child entitlements change the fixed V1 bytes; deprecated custom sandboxing is not used. V1 remains unchanged/not enrolled; no Keychain/Supervisor state, product consumer, runtime, backend, or guest participated. | [Pinned V2 archive](https://github.com/Shrimpworks/capsule-experiments/tree/0944ffd8cfd01ec23e4ae99138b0931d56804077/experiments/completed-compiled-artifact-payloads/payloads/capsule-corp/artifacts/mjs-source-validator-v2), machine-readable candidate profile, Darwin Go tests, native bootstrap/probe source, and entitlement-copy mutation oracle |
| ADR-0035 V2 supported macOS replacement design | `019fc2de-552d-77a0-aa47-35ac39d02edc` user-visible defensive/local-only research/design delegation after merged PR #123 | Scoped research/design `PASSED`; product Source Validator remains `BLOCKED`. Official Apple guidance makes direct App Sandbox inheritance `NO_GO` for this parser because it preserves daemon/Broker static rights. The only plausible supported composition found is a separately App-Sandboxed method-specific XPC launcher owning one fresh parser child, but its topology is unselected, App Sandbox grants a writable private container, and no public unprivileged hard memory cap is usable. The public-SDK footprint setter returned `KERN_NO_ACCESS`; `proc_pid_rusage` monitoring is reactive, not a peak ceiling. R0-R5 fixes the new identity, signing, resource, kill/drain/reap, fault, consumer, and update plan. V0/V1/V2 bytes remain unchanged; no signing identity, service, consumer, runtime, backend, or guest was used. | [`MJS_SOURCE_VALIDATOR_MACOS_PROFILE_REPLACEMENT.md`](MJS_SOURCE_VALIDATOR_MACOS_PROFILE_REPLACEMENT.md), [`README.md`](../artifacts/mjs-source-validator-v2-replacement/README.md), fixed C probe, machine-readable observation, official Apple links, and installed macOS 26.5 SDK header observations |
| ADR-0036 Source Validator R0 architecture/resource decision | `019fc2de-552d-77a0-aa47-35ac39d02edc` user-visible architecture/contract delegation after merged PR #126 | R0 architecture decision `PASSED`; product Source Validator and M2/S1 remain `BLOCKED`. Accepted ADR-0036 selects two role-specific private App-Sandboxed XPC launchers and matching parser/profile families; no shared service/result/cache/container/group/key exists. Each private container is residual scratch authority only, with mandatory cleanup/residue evidence that is not a confidentiality proof. A later evidence-derived reactive physical-footprint watermark replaces the unavailable hard ceiling without a hard-peak/exact-cap or host-availability claim; numeric threshold/cadence/baseline/overshoot/kill values remain unset. Exact stop conditions cover unsupported private reachability, authority/native-loading/network/filesystem escape, orphan/cleanup failure, mixed versions, and unacceptable measured host risk. V0/V1/V2 bytes remain unchanged; no artifact, signing identity, service, parser execution, consumer, runtime, backend, or guest was created or used. | [ADR-0036](adr/0036-select-role-separated-source-validator-launchers.md), [passive v1 boundary](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1.md), reconciled replacement/implementation plans, and canonical architecture/threat/roadmap records |
| ADR-0036 Source Validator R1 passive v1 contract | current user-visible defensive/local-only implementation | Scoped R1 `PASSED`; product Source Validator remains `BLOCKED`. Forty-six generated cases freeze separate daemon/Broker request, result, process-profile, artifact-profile, consumer-projection, and inactive resource-policy frames. Exact/cap-plus-one, trusted-context mismatch, cross-role, cross-version, trailing, and invented-measurement refusals have zero effects. Go and Node independently decode all cases; 14 role-distinct known answers are retained. The inactive policy binds structural frame/concurrency limits only, while all evidence-derived measurements remain unset for R4. No XPC service, child, parser execution, signing, consumer, runtime, backend, or guest exists. | [`MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1.md`](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1.md), generated conformance fixtures, `sourcevalidatorpassive` v1 Go codec/tests, independent Node verifier, and the then-current 335-field/26-target manifest checkpoint; current whole-manifest totals are recorded above |
| ADR-0036 Source Validator R2 unsigned role bundles | current user-visible defensive/local-only implementation | Scoped R2 `PASSED`; product Source Validator remains `BLOCKED`. Two private role-specific XPC bundle layouts, two native launchers, and two exact Oxc 0.140.0/Rust 1.95.0 parser children build offline and compare byte-for-byte across two clean same-host directories. Four executable digests are role-distinct. Source/lock/license/notice/CycloneDX/unsigned in-toto evidence and `libSystem`-only launcher closure are retained. Each launcher accepts only one `request` data value and validates the fixed role/frame/source/inactive-policy binding. Because R1's only accepted policy is inactive, it refuses without spawn. No Apple identity, signing, installation, enrollment, reachable service, consumer, active resource values, runtime, backend, or guest participated; this is not independent-builder evidence. | [Pinned R2 archive](https://github.com/Shrimpworks/capsule-experiments/tree/0944ffd8cfd01ec23e4ae99138b0931d56804077/experiments/completed-compiled-artifact-payloads/payloads/capsule-corp/artifacts/mjs-source-validator-r2), exact bundle/parser bytes, construction/build/source/license/SBOM/provenance/reproduction evidence, Rust tests, and benign role-known-answer process verifier |
| P0-1 FD-native immutable runtime-root custody | `019fc4c1-7d40-77b3-a2e9-51d3e2775972`, delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | PATCH-CANDIDATE: the narrow fixed-role raw-only API passed local attachment/custody, focused sanitizer tests, five deliberate mutations, and four owned unsandboxed guest digest runs with no root-path opens. P0-1C remains open until the exact final signed installed App Sandbox/protected-construction corpus passes. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-root-custody/RESULTS.md), [`FD_NATIVE_PATCH_REVIEW.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-root-custody/FD_NATIVE_PATCH_REVIEW.md), governed patch, and selected evidence |
| P0-2 `NullFs` disposition | Earlier replacement `019fc2e8-445e-7cb2-b4c2-54d84282c3fe`, replacing task `019fc2e6-cf9d-7210-b2f3-f3bf2244e83a`; later prototype merged in PR #30 | `GOVERNED-PATCH`: the smallest deletion failed bootstrap, but the later direct-block-root prototype booted without virtiofs, reran 36 adversarial plus four identity cases without the original failure, and made removal credible. It is not admitted; independent patch review, route closure, P0-1 custody, P0-3 transport, and final signed P0-4 evidence remain. | [`NULLFS_P0_2.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-adversarial/NULLFS_P0_2.md), [`NULLFS_P0_2_DISPOSITION.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-adversarial/NULLFS_P0_2_DISPOSITION.md), governed prototype patch, and compact evidence |
| P0-3 backend-independent framing | Merged in PR #27 | Conditional candidate pass only: 43 byte-exact vectors measured the proposed source/input/result/frame caps and retained binding, role, JSON, commit, drain, stall/death, EOF, runner-exit, and crash dispositions. No transport, launcher, guest, VMM, App Sandbox, Supervisor, approval, or teardown mechanism participated. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-p0-3-protocol-conformance/RESULTS.md), 43-vector manifest, local model, and measurement record |
| P0-3 libkrun console correctness | Merged in PR #28 | At that checkpoint stock could not proceed as-is. Governed patch `584ce48548fe969684fe3c55e57fbf56e7dae40af28c241c24c47b138faf1283` passed 51 local library tests and four regressions but still lacked the later sanitizer/coverage follow-up and all real composition. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-console-correctness/RESULTS.md), governed patch, verification record, and focused tests |
| P0-3 cross-language/console follow-up | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` on 2026-08-03 | This retained pre-fork checkpoint added independent Node agreement on all 43 vectors, six re-encodings, ten local pipe fault classes, Clippy, AddressSanitizer, repetition, four mutations, and the historical before measurement of 90/728 patched-file lines. The governed-fork row below supersedes it for current coverage. | Updated P0-3 and console `RESULTS.md`, independent verifier/fault harness, cross-language evidence, coverage summary, and mutation patches |
| Governed libkrun console/raw-FD source and fork reconciliation | Earlier read-only reconciliation delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc`; later user-visible defensive governance/integration task after public fork PRs [#2](https://github.com/Shrimpworks/libkrun/pull/2) and [#4](https://github.com/Shrimpworks/libkrun/pull/4) | Scoped r3 source/governance transition `PASSED`; parent governed libkrun work remains `IN_PROGRESS — TRENDING_GOOD`; backend/profile admission remains `BLOCKED`. The original patch-queue merge remains locked as `capsule/baseline-v1.19.4-r1` at `4ea8d1de861ed1c0636fc800b6da8fb71a086aa5`; the previous accepted merge remains locked as `capsule/upstream-v1.19.4` at `cf0333cdba478cc34a8570a65b38412da7fd3ecc`. PR #4 separately backported control-port and descriptor-chain validation and merged as `7432eda5a49220976b0167005aa43ee622f9d632`, tree `7671440cfbafa58fe20aebf8d4deb2a843ebe346`, with verified ancestry from upstream `728df8125077d0db44265f6e997c72b81b65c015` and the previous accepted head. Locked `capsule/upstream-v1.19.4-r3` is now the default. The five-patch aggregate remains `d19fd0ff159c699acccda2621519de45a09408bf3847b418ac34e02b79e805d5`; `main` remains protected, mutable integration state. No guest, installed product, release, or admission was evidenced. | Exact merged PR bodies and refs; governed verifier/check readback; [repository fork ledger and settings](REPOSITORY_SETUP.md#accepted-governed-fork-readback); canonical Gate C reconciliation |
| P0-4A installed development topology | Merged in PR #34 | Conditional topology pass only: 18 roles, 17 installed-entry readbacks, per-user registration/explicit activation, exact ad-hoc IPC identity, refusal cases, and same-session recovery passed without host root or a guest. App Sandbox failed before `main`; valid signing, Team enrollment, notarization/stapling/Gatekeeper, on-demand activation, clean hosts, sessions, and the macOS support floor remain open. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-installed-development-topology/RESULTS.md), closed manifests, scripts/tests, and selected installed-run evidence |

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

Proposed ADR-0033 selects a pre-created enrolled sibling object, exact
UID/mode/type/link/device/inode validation, and lifetime nonblocking BSD `flock`. Proposed
ADR-0038 now selects the installation-root-authorized, Supervisor-created enrollment ceremony. The retained
[development-only experiment](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/supervisor-owner-lock-boundary/RESULTS.md) observed
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

The first bounded
[G3 installed identity/session/update checkpoint](https://github.com/Shrimpworks/capsule-experiments/blob/3e9c9cbc3e0314439771151f1fd99c2b3a5a50b9/experiments/supervisor-owner-lock-installed-g3/RESULTS.md)
left installed G3 **BLOCKED before installed build**. Certificate SHA-1
`1638CFBD9250A00B4DBD81AE8FD1C790B42F61E3` is displayed as an Apple Development W4 identity, but
its public subject OU and the TeamIdentifier on an exact harmless signed probe are
`3DDR84M4JS`; every cached provisioning profile is also 3DDR. The separate Developer ID identity
was not used, and there was no ad-hoc fallback. The checkpoint retains exact test-only role,
service, entitlement/profile, state/lock/store, and complete bootstrap-field fixtures plus a pure
v1/v2 exact-tuple refusal model. Its noncredential run passes the real G1/G2 Darwin corpus. It
creates no signed bootstrap or installed/protected result and records the then-unresolved
trusted-installer-to-private-container bootstrap plus descriptor-relative store-open blockers.
I2A/Proposed ADR-0038 now resolves the first blocker in design by selecting a one-shot
Coordinator-authorized/Supervisor-created private root; passive object fixtures and installed I2B
evidence remain absent.
The 2026-08-04 exact-selector follow-up confirms the selector never chose Developer ID or another
fallback: the development certificate itself emits 3DDR. It also rejects the default designated
requirement as W4 evidence because that requirement binds only the misleading common name, not the
Team OU. Apple Membership Details later confirmed `3DDR84M4JS` is the Individual account Team ID
and `W4QUR9FUL4` is a member/display suffix. User-run discovery also reported new valid Apple
Development SHA-1 `80A4969BCD1B3926020888094B9D812A283D3793` with private-key pairing; presence is
not authorization. Both standard local profile caches contain only the same three 3DDR Gate
B/wildcard profiles with nonmatching App IDs, and installed G3/I2B remains blocked before build or
launch on exact role profiles, passive signed-object fixtures, the selected handoff/container
corpus, and descriptor-relative store composition.

I2A is `PASSED` in its exact architecture/contract scope. Proposed ADR-0038 selects an on-demand,
separately signed Trust Coordinator with a Coordinator-only installation-root Keychain group and a
bootstrap-only Coordinator/Supervisor App Group. The Coordinator constructs and installation-root-
signs the bounded request/final record; the Supervisor alone creates the fixed root, owner, and
fixed-v1 no-guest genesis in its private container. The retained I2A plan freezes field authority,
purpose/audience, nonce/time/replay/death behavior, descriptor-relative file/sync/rename ordering,
same-user/stale/debug/session/update/rollback oracles, and I2B1-I2B5. It created no fixture bytes,
signature, key, service, container, file, process, runtime, backend, or guest. Installed I2 remains
`BLOCKED`; product store selection and attempt activation remain outside I2A.

I2B2 is `PASSED` for unsigned installation-only bytes and layout. Two clean repository-local
constructions retain an exact 31-file/eight-role bundle, profile SHA-256
`7c6d410bd99b165a7f882914ca889d8796366d6ba60f0c76d5b30577abc6f5b7`, bundle-manifest SHA-256
`e92f7629774258f1dff68df7882b663479916c5feb4110db5460de3cef0af903`, exact I1A/I1B/I2B1
cross-links, inactive service/entitlement/constraint inputs, and no-create refusal. No identity,
profile, signing, key, service, process, container, protected state, runtime, backend, or guest was
used. The narrow production-shaped wrapper review is `PASSED` for passive checked-in vectors.
Installed I2B remains `BLOCKED` on separately authorized I2B3 signing, caller/key authorization,
App Group/service/container handoff, and descriptor-relative store-open evidence.

Source Preparer P0 remains `BLOCKED` as a separate conditional later feature merged in PR #72 from head
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
locations/counts, a distinct migration-genesis checkpoint, and generated answers. The follow-on
valid-v1 mapping contradiction is now passively resolved with a closed absent/present lifecycle
union and counts derived only from present records. The
[F2 v1 mapping resolution](SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md) and its real-v1 executable
witness are retained. The [stateful F2 result](SUPERVISOR_ARCHIVE_F2_MIGRATION_RESULT.md) now
implements the owner-asserted v1-to-v2 migration, downgrade refusal, and read-only empty-archive
full verifier with exact known answers and local fault/corruption/capacity/concurrency/process-
death evidence. It preserves the absent-lifecycle witness without recovery or invention. The
[stateful F3 result](SUPERVISOR_ARCHIVE_F3_ACTIVATION_RESULT.md) now implements the first sealed
immutable-segment prepare/verify/publish/activate transaction with exact known answers, complete
visible identity/tombstone retention, publish-before-reference ordering, atomic generation-two
activation, passive valid-orphan reporting, and fault/response-loss/corruption/substitution/
concurrency/owner-loss/process-death oracles. It adds no retained lookup, v2 mutation, second
segment, orphan deletion, backup, engine, adapter, or consumer. The design deliberately leaves
referenced-history deletion, implementation/installed validation of the selected owner lock and
power loss, coherent rollback prevention, continuous service, consumers, and guests blocked.
F4A's retained lookup/replay/passive-collision/hot-only-recovery scope is also `PASSED`. The
retained [F4B blocker](SUPERVISOR_ARCHIVE_F4B_MUTATION_BLOCKER.md) records the former effect-history
contradiction and its selected correction. The
[F4B result](SUPERVISOR_ARCHIVE_F4B_MUTATION_RESULT.md) now passes the exact fixed-store atomic
mutation, independent tombstone, durable replay, historical lookup, and fault/collision scope. The
[F4C result](SUPERVISOR_ARCHIVE_F4C_GROWTH_RESULT.md) passes deterministic second/later segment
activation, exact 64-segment acceptance/65 refusal, retained-global reconstruction, and complete
old-or-new fault/death evidence. The
[F5 result](SUPERVISOR_ARCHIVE_F5_BACKUP_RESULT.md) passes owner-held manifest-last coherent backup,
complete-copy verification, read-only exact-anchor restore admission, bounded offline reporting,
and explicit deletion of only a sealed known-unreferenced segment after live and supplied-backup
reference scans. It preserves unknown/corrupt/mixed/cross-installation evidence and activates no
restore. F6 remains deferred.

The current closed conformance corpus has 95 rules, 502 cases, and 624 fixtures. The generated
verifier is canonical for its current per-language distribution. Historical per-slice counts remain
evidence for those slices, not current whole-repository totals.

## Production CBOR/COSE dependency comparison

The [standalone comparison](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/production-cbor-cose-profile/RESULTS.md) is defensive,
local-only selection evidence. It pins and records `fxamacker/cbor` v2.9.2,
`veraison/go-cose` v1.3.0, and `x448/float16` v0.8.4, replays the retained object and
cross-language known answers, exercises restoration mutations and trusted-key/binding refusal, and
retains bounded fuzz/resource/footprint results. It selects only object-specific fxamacker typed
encode/decode and records a production NO-GO for go-cose.

The later [v0 object-set and wrapper result](V0_CBOR_OBJECT_SET_AND_WRAPPER.md) keeps
`SourceManifest` v0 implementation-eligible and freezes the passive I2B1 request/record set while
retaining plan/registration/approval and the conditional TypeScript family as pre-freeze. Exact
fxamacker v2.9.2 plus float16 v0.8.4 remain behind object-specific codecs. The subsequent
[production-shaped wrapper review](PRODUCTION_SHAPED_CBOR_COSE_WRAPPER_REVIEW.md) corrects exact
signed-request pairing, complete repeated-field binding, and payload-owned replay identity; Go and
independent Swift agree on 95 cases and Go retains three fuzz targets. No production signing key,
COSE dependency, authenticated role/profile/CDHash path, local key-authorization registry, durable
replay ledger, consumer, IPC activation, store, runtime, backend, guest, ADR-0019 promotion, or
product-control status is added.

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
| PR #72, Source Preparer P0 authority review | `2e268b01d4174fe90397c00abc5973a3dd785606` | Source head `a12041c36d90815474598f0929c595b32dc68e11`; the conditional later TypeScript feature is currently `BLOCKED` pending protected-store, worker, bootstrap/update, retention/release, recursive field-authority, and lifecycle evidence. No P1 contract exists. |
| PR #73, governed libkrun reconciliation | `f6fcf172af752a425afb29ce62680d0b115f6998` | Source head `5e17ac8cec21320c3693049c53e7575bb9dbc15a`; reconciled external fork PR #2 head `8a2c91943793668f31a1cf7af431933be935bb58` and merge `cf0333cdba478cc34a8570a65b38412da7fd3ecc`, two lifecycle fixes, and 37/88-function plus 298/733-line four-file coverage without guest/backend/profile admission. |
| PR #74, Supervisor owner-lock boundary | `e930f9dbd877bea0cbd55870060f48c9c7fdd72f` | Final reviewed source head `afd148c92f4b9f6f35f2a7d9161502cd1175a107`; proposed ADR-0033 selects enrolled pre-created inode plus lifetime nonblocking BSD `flock` for cooperating Supervisor ownership only. It supplies neither same-UID containment nor Source Preparer storage protection; G1 was the next slice at that merge. |
| PR #75, passive archive F1 | `6fc31a049c476acf5085071c48d3d5e36f27240f` | Source head `20c8d7df1d9ed3eb009e8ce9a0afbd41e03807ef`; implemented passive archive types, limits, known answers, defensive copies, and eligibility only. No file write, migration, activation, lookup consumer, or adapter call exists; F2 is next. |
| PR #86, ecosystem reuse map | `5e202b2` | Recorded primary-source-backed reuse/adoption lanes, including the bounded production CBOR/COSE comparison, without admitting a dependency or implementation. |
| PR #87, passive `.mjs` source manifest | `599d091` | Added the accepted single-member source profile's Go/TypeScript passive bytes, manifest, fixtures, and verification without a parser worker, store, consumer, runtime, or guest. |
| PR #100, `.mjs` parser/process selection | `33193ad` | Selected exact Oxc 0.140.0 and a future disposable-process boundary after the bounded parse-only comparison; no product parser or runtime enforcement was admitted. |
| PR #102, archive F2 v1 mapping stop | `e1ce148` | Retained the exact valid-v1 attempt-without-lifecycle witness and stopped before inventing v2 state. The later passive correction and PR #122 supersede this as the current implementation status. |
| PR #103, owner-lock G3 discovery | `9cca9dd` | Retained the certificate/profile mismatch and installed-build stop. Current translation: the bounded discovery task `PASSED`; installed G3 remains `BLOCKED`, not `NO_GO`. |
| PR #104, production CBOR/COSE comparison | `316a174` | Selected pinned `fxamacker/cbor` conditionally for a future narrow object wrapper and rejected `go-cose` as a product dependency; no root dependency or authority path changed. |
| PR #105, Source Validator passive V0 | `566e323` | Added exact passive Go/Rust request/result/profile frames and known answers without launching a parser or adding a consumer. |
| PR #106, owner-lock G3 W4 diagnostics | `810cbcc` | Confirmed the local display-name suffix does not match the signed TeamIdentifier/profile Team and retained the noncredential G1/G2 regression model. |
| PR #107, experiment archive move | `1a53623` | Moved completed one-time experiment code and evidence to the commit-pinned `Shrimpworks/capsule-experiments` repository while retaining canonical decisions and links here. |
| PR #108, archive F2 passive mapping correction | `74f10cb` | Added the absent/present lifecycle union and independent counts needed to represent the valid v1 crash world without invention. |
| PR #109, Source Validator V1 artifact | `810479c` | Retained the exact unwired Oxc artifact, locked graph, manifests, V0/M1 agreement, and same-host reproduction without enrollment or a product consumer. |
| PR #110, work-status language | `f485340` | Standardized `PASSED`, trended `IN_PROGRESS`, `BLOCKED`, and path-abandonment-only `NO_GO` reporting across current planning. |
| PR #117, organization-link reconciliation | `7afe669` | Updated canonical repository links after the Capsule and idle governed-fork ownership changes without changing source or evidence claims. |
| PR #118, automated dependency guards | `37c038c` | Added Dependabot/toolchain constraints so incompatible or manually governed update classes are not repeatedly proposed. |
| PR #119, bounded bug-hunt fixes | `c5050af` | Corrected low-risk script and protocol defects under the existing contracts. |
| PR #120, grouped bug-fix reconciliation | `946d9ea` | Integrated the reviewed CBOR canonicalization, shutdown fault-injection, and attestation-set comparison corrections from the independent bug-hunt tasks. |
| PR #121, Source Validator V2 process profile | `6641ea1` | The evidence task `PASSED`; exact V2 remains `BLOCKED`. Fixed launch/fault mechanics passed, while unsupported `RLIMIT_AS`, ambient authority in the unbounded diagnostic, and App Sandbox byte-identity changes require a new enrolled artifact and supported confinement design. |
| PR #122, archive F2 stateful migration | `4af9c50` | The exact local F2 slice `PASSED`: owner-asserted v1-to-v2 migration, downgrade refusal, missing-lifecycle preservation, empty-archive full verification, and fault/corruption/capacity/concurrency/death oracles. F3+ and production admission remain outside that scope. |

The ecosystem reuse audit is retained in
[`ECOSYSTEM_REUSE_AND_ADOPTION.md`](ECOSYSTEM_REUSE_AND_ADOPTION.md) at PR #86's integrated
checkpoint. It is a primary-source-backed planning map, not implementation or security evidence.
It records three bounded decision lanes: the completed production CBOR/COSE profile comparison,
SQLite comparison only after the F4-F5 remainder of the logical archive oracle, and the now-blocked Source Validator
V2 replacement-profile boundary. It does not preempt F4, owner-lock G3, or the separate ARM64
`rusty_v8` work.

The 2026-08-04 external macOS installation/update/repair/uninstall review is retained through the
[reconciled installation plan](MACOS_INSTALLATION_AND_DISTRIBUTION_PLAN.md), not as product
evidence. Its one-app/DMG direction and state-machine work are carried forward. Its direct
daemon/Broker parser-child topology is superseded by the supported-profile result: only distinct
App-Sandboxed role-specific launcher candidates proceed. Proposed ADR-0038 now selects the
protected-root bootstrap/Trust Coordinator contract, while its App Group and installed evidence,
ordinary IPC groups, Bundle Replacer, minimum-OS, update, and erasure claims remain explicit
decisions or signed-evidence gates rather than accepted platform facts.

The post-I1B
[Apple-platform semantics research](MACOS_INSTALLATION_PLATFORM_RESEARCH.md) is `PASSED` in its
bounded research scope. It retains the documented `SMAppService` lifecycle, supported private-XPC
and App-Group-named bootstrap candidate, `JoinExistingSession=true` requirement to test,
App-Group residual authority, stale-Coordinator Keychain limitation, bounded
`NSUpdateSecurityPolicy` claim, same-volume authorized-replacement candidate, exact failure
tables, and one-host versus clean-host evidence matrix. It performed no signing, portal,
Keychain, service, installation, replacement, runtime, backend, or guest mutation. Installed I2B,
replacement I4, and distribution I6 remain `BLOCKED` on their named evidence.

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
- passive field-authority verification now covering 1,111 fields across 81 profiles and 50 targets,
  without claiming coverage for future Source Preparer or plan-v1 objects; and
- archive Slice F1 passive projections and eligibility selection; F2 owner-asserted all-hot
  migration/full verification; F3 immutable-segment publication/activation; F4A retained lookup;
  F4B atomic mutation/independent effect tombstones; F4C bounded later growth; and F5 coherent
  backup/read-only restore admission/explicit known-orphan cleanup/offline reporting. Production
  engine, restore activation, referenced-history deletion, adapter calls, and consumers remain
  absent; and
- Proposed ADR-0033's local owner-lock mechanism selection, passive G1 Go/Darwin acquisition, and
  bounded G2 current-v1/no-guest startup composition, plus I2A's passed owner/contract decision;
  signed-object fixtures and installed protected-root evidence remain unimplemented, and installed
  G3/I2B is blocked by exact Team profiles, the protected-container/handoff corpus, and descriptor-
  relative store composition; and
- governed `deno_core` physical omission, same-host package reproduction, the abandoned official
  V8 asset-publication route,
  standalone dynamic-root evidence, and the now-closed fork Linux/arm64 construction blocker.
  Accepted ADR-0028
  selects its engineering order without admitting a profile; the real Deno and `rusty_v8`
  governed baseline branches are merged, but no governed release or admitted artifact exists.
- exact public governed libkrun source identity through merge
  `cf0333cdba478cc34a8570a65b38412da7fd3ecc`, with the unchanged five-patch aggregate, bounded
  console/raw-FD library evidence, two local lifecycle fixes, and improved coverage. The Linux-
  arm64 library build passed while its unit job was cancelled; absent independent review,
  remaining uncovered code, release obligations, and every
  guest/product admission boundary remain explicit.

Current dependency and priority view:

1. **Archive logical oracle complete through F5:** F2 migration/full verification, F3 first-segment activation, F4A
   retained lookup/replay/passive-collision routing, and the
   [F4B atomic mutation result](SUPERVISOR_ARCHIVE_F4B_MUTATION_RESULT.md), and
   [F4C bounded-growth result](SUPERVISOR_ARCHIVE_F4C_GROWTH_RESULT.md), and
   [F5 backup result](SUPERVISOR_ARCHIVE_F5_BACKUP_RESULT.md) are `PASSED` in their exact local
   scopes. ADR-0040 permits only its bounded owner-only internal-alpha exception; keep F6
   production-engine work separate and do not infer restore activation or continuity.
2. **Source Validator R3 signed composition passed; product blocked and post-alpha:** Accepted ADR-0036 closes R0 with
   two role-specific private launchers, residual private-container scratch plus mandatory cleanup,
   and an evidence-derived reactive watermark with no hard-peak or host-availability claim. Product
   R1 passive v1 contracts/field authority, R2 unsigned role-specific construction, and the exact
   Apple Development R3 signed/installed/refusal/cleanup composition are `PASSED`; product work
   remains `BLOCKED`; exact R4-v1 candidates are `NO_GO` and R4-v2 was not executed. ADR-0040 moves
   R4/R5 off the internal-alpha critical path. Resume only with a supportable child-lifetime and
   residual-container contract, then consumers.
   Independently freeze the signed-object set, narrow `fxamacker/cbor` wrapper responsibilities,
   and maintain exact documentation plus field authority.
3. **Installation plan:** I0, I1A, I1B, I2A, I2B1 passive request/record objects, I2B2 unsigned
   construction, and the post-I1B platform-semantics research are `PASSED` in their exact scopes.
   The narrow production-shaped wrapper review is `PASSED`; continue with separately authorized
   I2B3 signed Coordinator/bootstrap, caller/key authorization, protected-container, durable replay,
   and descriptor-relative fixed-v1 evidence. Manual whole-bundle replacement
   remains I4 and `BLOCKED`; automatic TUF/update-replacer, Developer ID distribution, a support-floor matrix,
   and complete uninstall remain I5-I6.
4. **Governed runtime:** the `rusty_v8` fork's bounded Linux/arm64 build is merged and the original
   construction blocker is closed. Next consume only an independently reviewed exact fork artifact
   into the fork-native Capsule bundle/reproduction path; governed release publication,
   independent-builder equality, and runtime/profile admission remain open.
5. **Later composition:** governed runtime plus libkrun requires admitted artifacts, the remaining
   transport/launcher/root/device/teardown work, and explicit authorization for an owned disposable
   development guest. No current task authorizes a guest.
6. **Credential/environment dependent:** Apple Development identity/profile use was deliberately
   authorized only for the exact I1B/R3 installed matrix. Exact G3 discovery
   disproved the display-name inference: SHA-1 `1638...61E3` says W4 in its common name, but its
   subject OU and emitted TeamIdentifier are `3DDR84M4JS`. Apple Membership Details confirms 3DDR
   is the account Team and W4 is a member/display suffix. Apple Development SHA-1
   `80A4...D3793` and three exact explicit-App-ID profiles were selected for I1B/R3; public metadata
   and hashes are retained, while raw profiles and credentials are not. I2B requires its additional
   exact Coordinator/bootstrap identifiers and profiles plus the protected-root signed-record/
   store-open composition. A separate Developer ID Application identity for Team `3DDR84M4JS` is
   later distribution authority requiring explicit authorization and matching-Team package design;
   it does not make Developer ID/notarization work current. Paid owned
   clean-host/minimum-OS coverage is not currently planned
   and remains deferred activation/distribution evidence, not a blocker for F4 local archive work;
   owner-lock installed G3 remains `BLOCKED` until the named credential and design
   blockers close.
   The canonical practical guide is now
   [Apple certificates, credentials, identifiers, entitlements, and Capsule keys](APPLE_CERTIFICATES_CREDENTIALS_AND_KEYS.md).
   Its documentation slice is `PASSED`; it changes no credential or product admission. Team
   `3DDR84M4JS`, the ADR-0037 containing topology, and the three exact I1B/R3 profiles are resolved
   for that completed experiment. Later I2B3 still requires exact Coordinator/bootstrap profiles,
   caller/key authorization, and App Group/service/container handoff; Developer ID/notarization/CI
   release custody remains deferred to separately authorized work.
   A genuinely independent
   Linux/arm64 builder is viable but not currently planned; same-host/GitHub-CI equality remains
   limited and independent-builder equality is deferred.
7. **Later Source Validator path:** accepted ADR-0034's M1 bytes, Accepted ADR-0035's passive V0
   frames, and the bounded V1/V2 evidence checkpoints are retained unchanged. Accepted ADR-0036
   closes R0 with two role-specific private launchers, residual private-container scratch plus
   mandatory cleanup, and an evidence-derived reactive watermark with no hard-peak/availability
   claim. Its exact [R3 execution packet](SOURCE_VALIDATOR_R3_EXECUTION_PACKET.md) records Team 3DDR,
   R2 byte identities/placement, profiles/entitlements, refusal/cleanup, and mutation gates. R1
   passive v1 contracts/field authority, R2 unsigned construction, and R3 exact signed installed
   composition are `PASSED`. Exact R4-v1 candidates are `NO_GO`; R4-v2 is unexecuted. Product work
   remains `BLOCKED` outside the internal-alpha critical path and may resume only with a revised
   supportable lifetime/residue contract. M2/S1 no longer waits on R5 under ADR-0040. No product
   validator endpoint is authorized.
8. **Passive authority-plane cutover passed; product connection blocked:** accepted ADR-0034 removes
   Source Preparer/plan-v1 from the first-release critical path. M1 and the exact-one-`main.mjs`
   proposal feed the atomic `RegisterPlanV0`/fetch fixed-store oracle with generated caps and
   complete current field authority. No authenticated product endpoint, Broker/approval consumer,
   attempt, runtime, backend, or guest is authorized.

TypeScript remains conditional. If later selected, Source Preparer P0A and ADR-0030's atomic
plan-v1/RegisterPlanV1 cutover still apply with no dual active v0/v1 acceptance. CommonJS, package
resolution, legacy Node module surface, and runtime-contract widening remain forbidden.

The governed Deno and `rusty_v8` forks now default to their locked latest accepted lines:
`capsule/accepted-v2.9.4-r3` at `3fa21d1ae7705ab4bcb4bc98955f25301b20122a` and
`capsule/accepted-v150.2.0-r5` at `d09221062280ae1675fe26c53c3f43871aae2055`. Governed libkrun
now defaults to locked `capsule/upstream-v1.19.4-r3` at
`7432eda5a49220976b0167005aa43ee622f9d632`; its original patch-queue merge and superseded accepted
line remain locked. Each fork preserves `main` as protected mutable upstream-integration state, not
Capsule adoption. This is repository governance only, not source/build independence, release,
security-control evidence, profile validation, or product admission. Integration must still pin
exact commits and digests and verify ancestry rather than trust a branch name.

Proposed ADR-0029 selects an IPC topology but does not implement or validate its native bridge,
installed endpoints, peer identities, or production transport. This checkpoint also does not decide
retained archive mutation/growth or production-engine archive/compaction beyond the F1-F4A local oracle,
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
