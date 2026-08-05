# Workstream and evidence ledger

Date: 2026-08-04

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
| Governed `deno_core` C1 passive composition | User-visible implementation task delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | Scoped passive contract `PASSED`; parent governed-runtime work `IN_PROGRESS — TRENDING_GOOD`; runtime/profile admission `BLOCKED`. Exact-byte schema/fixture plus independent Go and Node validators freeze the app JSON-in/JSON-out behavior, desired globals/ops/modules/files, loader absence, logical descriptors, resource references, all eight construction identities, and refusal/stops with every authority effect false. C2 owns numeric descriptors, machine resources, guest execution, and composed evidence. | [C1 contract and limitations](protocol/GOVERNED_DENO_CORE_C1_COMPOSITION.md), `schemas/conformance/c1-governed-deno-core/controlled-development-profile.json`, and independent validators |
| Governed `deno_core` C2A passive execution-profile preparation | User-visible implementation/research task delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` after merged PR #155 | Scoped passive preparation `PASSED`; parent governed-runtime work `IN_PROGRESS — TRENDING_GOOD`; C2B and runtime/profile admission `BLOCKED`; `RUNTIME-001` and `VMM-001` unsupported. Exact-byte schema/fixture plus independent Go and Node validators consume unchanged C1 bytes and freeze plan roles, numeric descriptors, candidate machine/transport/teardown values, known-answer bytes, null final-artifact/runtime-profile/resource blockers, 91 C2B cases, 18 restoration mutations, and all-zero authority effects. No guest, process, runtime, credential, release, or admission was created. | [C2A contract, matrix, blockers, and limitations](protocol/GOVERNED_DENO_CORE_C2A_EXECUTION_PROFILE.md), `schemas/conformance/c2a-governed-deno-core/passive-execution-profile.json`, and independent validators |
| Governed `deno_core` C2B fixed-fixture passive binding | User-visible implementation/documentation task delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | Scoped passive reconciliation `PASSED`; parent governed-runtime work `IN_PROGRESS — TRENDING_GOOD`; C2B composed-profile/guest execution and runtime/profile admission `BLOCKED`; `RUNTIME-001` and `VMM-001` unsupported. The 8,221-byte v1 object preserves exact C1/C2A bytes and binds Deno draft PR #2 head/tree, experiments draft PR #3 head/tree, the v1/v2 evidence lineage, fixed known-answer bytes, six artifacts, 143 recursive field-authority paths, and zero effects. Same-host equality is not independent-builder equality. Both draft PRs are merge dependencies and any identity change requires a new binding. | [C2B passive binding, validation, limitations, and next gate](protocol/GOVERNED_DENO_CORE_C2B_PASSIVE_BINDING.md), `schemas/conformance/c2b-governed-deno-core/passive-binding.json`, mirrored draft-head inputs, and independent Go/TypeScript validators |
| P0-0 TypeScript approved-byte boundary | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` from merged checkpoint `5448943` | Passed only the exact pre-approval byte-ownership question for a strip-only ESM TypeScript subset. Exact Node 22.22.1/Amaro 1.1.5 emission was deterministic for fixed fixtures and mutation/cap/diagnostic cases failed closed. Proposed ADR-0026 binds original and executable source roles before registration. No current contract changed, no owner/runtime was selected, and `RUNTIME-001` remains unsupported. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/typescript-approved-byte-boundary/RESULTS.md), fixed Node and `deno_ast` probes, independent Go verifier, and selected evidence |
| ADR-0034 M1 source-byte/manifest foundation and scanner disposition | `019fca55-554d-7b23-bb2e-6d2c4b0e16cd`, delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | Passive byte/manifest foundation `PASSED`. The handwritten scanner is `NO_GO` and was removed after a retained division-versus-regexp counterexample exposed a live-`import()` false negative. The replacement parser/process work proceeded separately; no runtime/no-loader claim follows. | [`MJS_MODULE_REQUEST_VALIDATOR_HOLD.md`](MJS_MODULE_REQUEST_VALIDATOR_HOLD.md), generated MJS conformance fixtures, SourceManifest CDDL, and field-authority target |
| `.mjs` parser/process boundary | `019fc2de-552d-77a0-aa47-35ac39d02edc` user-visible defensive/local-only research delegation | Scoped parser/process selection `PASSED`; later product Source Validator work is now `BLOCKED` by the separately retained V1 enrollment and V2 resource/confinement stops. Exact Oxc 0.140.0 matched all 33 local parse-only cases and all 28 canonical M1 outcomes. Accepted ADR-0035 selects a future one-shot Source Validator outside daemon/Broker/Supervisor address spaces; ADR acceptance records architecture, not product admission. Product artifact enrollment, supported confinement, parent integration, and runtime no-loader evidence remain. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/mjs-parser-boundary/RESULTS.md), M1 mapping, exact lock/fixtures/classifications, supply inventory, measurements, and [`HANDOFF.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/mjs-parser-boundary/HANDOFF.md) |
| ADR-0035 V0 passive Source Validator contract | `019fc2de-552d-77a0-aa47-35ac39d02edc` user-visible defensive/local-only implementation delegation | Passive byte-contract observation only: exact request/result/engineering-candidate/artifact-profile fixed frames, caps, domain separation, closed status/count relationships, 128 generated cases, zero-effect refusals, and independent unwired Go plus test-only Rust agreement are retained. The generator consumes M1 bytes in place. No Oxc product dependency, parser invocation, child, sandbox, endpoint, artifact enrollment, consumer, Approval/key effect, runtime, backend, or guest exists; accepted ADR-0035 does not change this evidence boundary. | [`MJS_SOURCE_VALIDATOR_PASSIVE_CONTRACT.md`](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_CONTRACT.md), generated conformance corpus, Go passive package, and Rust test-only oracle |
| ADR-0035 V1 Source Validator artifact | `019fc2de-552d-7b23-bb2e-6d2c4b0e16cd` separate user-visible defensive/local-only implementation delegation | Current translation: the bounded unwired V1 artifact slice `PASSED`; product enrollment is `BLOCKED`, not `NO_GO`. Exact Oxc 0.140.0 plus Rust 1.95.0 and pinned `sha2` produces a 1,146,656-byte macOS arm64 executable at SHA-256 `ba2a6b38...cb600`. The 74-dependency lock, registry/source hashes, licenses/notices, CycloneDX SBOM, unsigned in-toto statement, V0/M1 deterministic results, mutations, and same-host two-directory reproduction are retained. The actual V0 profile binds executable/build/assessment but is explicitly not enrolled. Independent builder/clean host, installation-authority signatures, vulnerability owner/SLA, V2 confinement, V3/V4 consumers, V5 expansion, V6 runtime, backend, and guest remain absent. | [`README.md`](../artifacts/mjs-source-validator-v1/README.md), [`ASSESSMENT.md`](../artifacts/mjs-source-validator-v1/ASSESSMENT.md), exact binary/profile, build/source/license/SBOM/provenance/reproduction evidence, Rust tests, and process interruption oracle |
| ADR-0035 V2 Source Validator process profile | `019fcd7b-b0d1-74b1-9eb9-4dc1db12509d` separate user-visible defensive/local-only implementation/research delegation | Exact V2 candidate `BLOCKED`. The strict test-only macOS bootstrap fixes target/profile argv, empty environment, cwd, post-exec FD inventory, CPU/file/FD/child/output/wall limits, group kill/reap, and fault refusal, but `RLIMIT_AS` returns `EINVAL` before V1 `exec`. An explicitly unbounded diagnostic mutation verifies ordinary/exact-maximum frames, partial/duplicate/trailing/oversize/crash/signal/timeout/cancel and clean-later behavior while proving owned out-of-cwd reads, IPv4/Unix socket creation without connect, cwd metadata/empty-file writes, and a 512 MiB mapping remain. Supported App Sandbox child entitlements change the fixed V1 bytes; deprecated custom sandboxing is not used. V1 remains unchanged/not enrolled; no Keychain/Supervisor state, product consumer, runtime, backend, or guest participated. | [`README.md`](../artifacts/mjs-source-validator-v2/README.md), [`OBSERVATIONS.md`](../artifacts/mjs-source-validator-v2/OBSERVATIONS.md), machine-readable candidate profile, Darwin Go tests, native bootstrap/probe source, and entitlement-copy mutation oracle |
| ADR-0035 V2 supported macOS replacement design | `019fc2de-552d-77a0-aa47-35ac39d02edc` user-visible defensive/local-only research/design delegation after merged PR #123 | Scoped research/design `PASSED`; product Source Validator remains `BLOCKED`. Official Apple guidance makes direct App Sandbox inheritance `NO_GO` for this parser because it preserves daemon/Broker static rights. The only plausible supported composition found is a separately App-Sandboxed method-specific XPC launcher owning one fresh parser child, but its topology is unselected, App Sandbox grants a writable private container, and no public unprivileged hard memory cap is usable. The public-SDK footprint setter returned `KERN_NO_ACCESS`; `proc_pid_rusage` monitoring is reactive, not a peak ceiling. R0-R5 fixes the new identity, signing, resource, kill/drain/reap, fault, consumer, and update plan. V0/V1/V2 bytes remain unchanged; no signing identity, service, consumer, runtime, backend, or guest was used. | [`MJS_SOURCE_VALIDATOR_MACOS_PROFILE_REPLACEMENT.md`](MJS_SOURCE_VALIDATOR_MACOS_PROFILE_REPLACEMENT.md), [`README.md`](../artifacts/mjs-source-validator-v2-replacement/README.md), fixed C probe, machine-readable observation, official Apple links, and installed macOS 26.5 SDK header observations |
| ADR-0036 Source Validator R0 architecture/resource decision | `019fc2de-552d-77a0-aa47-35ac39d02edc` user-visible architecture/contract delegation after merged PR #126 | R0 architecture decision `PASSED`; product Source Validator and M2/S1 remain `BLOCKED`. Accepted ADR-0036 selects two role-specific private App-Sandboxed XPC launchers and matching parser/profile families; no shared service/result/cache/container/group/key exists. Each private container is residual scratch authority only, with mandatory cleanup/residue evidence that is not a confidentiality proof. A later evidence-derived reactive physical-footprint watermark replaces the unavailable hard ceiling without a hard-peak/exact-cap or host-availability claim; numeric threshold/cadence/baseline/overshoot/kill values remain unset. Exact stop conditions cover unsupported private reachability, authority/native-loading/network/filesystem escape, orphan/cleanup failure, mixed versions, and unacceptable measured host risk. V0/V1/V2 bytes remain unchanged; no artifact, signing identity, service, parser execution, consumer, runtime, backend, or guest was created or used. | [ADR-0036](adr/0036-select-role-separated-source-validator-launchers.md), [passive v1 boundary](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1.md), reconciled replacement/implementation plans, and canonical architecture/threat/roadmap records |
| ADR-0036 Source Validator R1 passive v1 contract | current user-visible defensive/local-only implementation | Scoped R1 `PASSED`; product Source Validator remains `BLOCKED`. Forty-six generated cases freeze separate daemon/Broker request, result, process-profile, artifact-profile, consumer-projection, and inactive resource-policy frames. Exact/cap-plus-one, trusted-context mismatch, cross-role, cross-version, trailing, and invented-measurement refusals have zero effects. Go and Node independently decode all cases; 14 role-distinct known answers are retained. The inactive policy binds structural frame/concurrency limits only, while all evidence-derived measurements remain unset for R4. No XPC service, child, parser execution, signing, consumer, runtime, backend, or guest exists. | [`MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1.md`](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1.md), generated conformance fixtures, `sourcevalidatorpassive` v1 Go codec/tests, independent Node verifier, and 335-field/26-target field-authority manifest |
| ADR-0036 Source Validator R2 unsigned role bundles | current user-visible defensive/local-only implementation | Scoped R2 `PASSED`; product Source Validator remains `BLOCKED`. Two private role-specific XPC bundle layouts, two native launchers, and two exact Oxc 0.140.0/Rust 1.95.0 parser children build offline and compare byte-for-byte across two clean same-host directories. Four executable digests are role-distinct. Source/lock/license/notice/CycloneDX/unsigned in-toto evidence and `libSystem`-only launcher closure are retained. Each launcher accepts only one `request` data value and validates the fixed role/frame/source/inactive-policy binding. Because R1's only accepted policy is inactive, it refuses without spawn. No Apple identity, signing, installation, enrollment, reachable service, consumer, active resource values, runtime, backend, or guest participated; this is not independent-builder evidence. | [`README.md`](../artifacts/mjs-source-validator-r2/README.md), [`ASSESSMENT.md`](../artifacts/mjs-source-validator-r2/ASSESSMENT.md), exact bundle/parser bytes, construction/build/source/license/SBOM/provenance/reproduction evidence, Rust tests, and benign role-known-answer process verifier |
| P0-1 FD-native immutable runtime-root custody | `019fc4c1-7d40-77b3-a2e9-51d3e2775972`, delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` | PATCH-CANDIDATE: the narrow fixed-role raw-only API passed local attachment/custody, focused sanitizer tests, five deliberate mutations, and four owned unsandboxed guest digest runs with no root-path opens. P0-1C remains open until the exact final signed installed App Sandbox/protected-construction corpus passes. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-root-custody/RESULTS.md), [`FD_NATIVE_PATCH_REVIEW.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-root-custody/FD_NATIVE_PATCH_REVIEW.md), governed patch, and selected evidence |
| P0-2 `NullFs` disposition | Earlier replacement `019fc2e8-445e-7cb2-b4c2-54d84282c3fe`, replacing task `019fc2e6-cf9d-7210-b2f3-f3bf2244e83a`; later prototype merged in PR #30 | `GOVERNED-PATCH`: the smallest deletion failed bootstrap, but the later direct-block-root prototype booted without virtiofs, reran 36 adversarial plus four identity cases without the original failure, and made removal credible. It is not admitted; independent patch review, route closure, P0-1 custody, P0-3 transport, and final signed P0-4 evidence remain. | [`NULLFS_P0_2.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-adversarial/NULLFS_P0_2.md), [`NULLFS_P0_2_DISPOSITION.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-adversarial/NULLFS_P0_2_DISPOSITION.md), governed prototype patch, and compact evidence |
| P0-3 backend-independent framing | Merged in PR #27 | Conditional candidate pass only: 43 byte-exact vectors measured the proposed source/input/result/frame caps and retained binding, role, JSON, commit, drain, stall/death, EOF, runner-exit, and crash dispositions. No transport, launcher, guest, VMM, App Sandbox, Supervisor, approval, or teardown mechanism participated. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-p0-3-protocol-conformance/RESULTS.md), 43-vector manifest, local model, and measurement record |
| P0-3 libkrun console correctness | Merged in PR #28 | At that checkpoint stock could not proceed as-is. Governed patch `584ce48548fe969684fe3c55e57fbf56e7dae40af28c241c24c47b138faf1283` passed 51 local library tests and four regressions but still lacked the later sanitizer/coverage follow-up and all real composition. | [`RESULTS.md`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-libkrun-console-correctness/RESULTS.md), governed patch, verification record, and focused tests |
| P0-3 cross-language/console follow-up | Delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc` on 2026-08-03 | This retained pre-fork checkpoint added independent Node agreement on all 43 vectors, six re-encodings, ten local pipe fault classes, Clippy, AddressSanitizer, repetition, four mutations, and the historical before measurement of 90/728 patched-file lines. The governed-fork row below supersedes it for current coverage. | Updated P0-3 and console `RESULTS.md`, independent verifier/fault harness, cross-language evidence, coverage summary, and mutation patches |
| Governed libkrun console/raw-FD source reconciliation | User-visible read-only reconciliation delegated from `019fc2de-552d-77a0-aa47-35ac39d02edc`; public fork PR #2 | Exact head `8a2c91943793668f31a1cf7af431933be935bb58` merged as `cf0333cdba478cc34a8570a65b38412da7fd3ecc` over governed base `4ea8d1de861ed1c0636fc800b6da8fb71a086aa5`. The unchanged five-patch aggregate is `d19fd0ff159c699acccda2621519de45a09408bf3847b418ac34e02b79e805d5`. The merge fixes queued-backpressure shutdown and inactive-state defects; four-file coverage moves from 13/88 to 37/88 functions and 90/728 to 298/733 lines. The no-guest property/raw-FD, ASan, Clippy, repetition, mutation, reconstruction, and macOS default-init cross-build head checks passed. At 2026-08-03T22:57:43Z Linux-arm64 unit tests remained queued, the overall head state was pending, no merge-commit checks or submitted independent/CODEOWNER review existed, and the advanced baseline branch no longer satisfied the verifier's hardcoded earlier-base pin. No guest, VMM transport, installed product, release, or admission was evidenced. | [Public PR #2](https://github.com/Shrimpworks/libkrun/pull/2), exact [head](https://github.com/Shrimpworks/libkrun/commit/8a2c91943793668f31a1cf7af431933be935bb58), exact [merge](https://github.com/Shrimpworks/libkrun/commit/cf0333cdba478cc34a8570a65b38412da7fd3ecc), and canonical Gate C reconciliation |
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
`a061291fe76d3bb460673adf25a322b0aa6d87d43619503eacaf3889eef4144b`, bundle-manifest SHA-256
`f706e3597958a6f694de7fb7c57f3e66d9cd5cd6a7f99e389de40018923c5c5d`, exact I1A/I1B/I2B1
cross-links, inactive service/entitlement/constraint inputs, and no-create refusal. No identity,
profile, signing, key, service, process, container, protected state, runtime, backend, or guest was
used. Installed I2B remains `BLOCKED` on production wrapper review and separately authorized I2B3
signing/key/App Group/service/container handoff evidence.

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
old-or-new fault/death evidence. F5-F6 remain deferred.

The current closed conformance corpus has 89 rules, 330 cases, and 433 fixtures. The unwired
Go/TypeScript implementation covers the previously recorded 177 Go and 80 TypeScript proposal
targets plus 40 independently verified MJS byte/manifest targets in each language. Twenty-eight
exact source-language adjudication cases remain `pending`. Slice C adds focused Go tests rather
than manifest cases.

## Production CBOR/COSE dependency comparison

The [standalone comparison](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/production-cbor-cose-profile/RESULTS.md) is defensive,
local-only selection evidence. It pins and records `fxamacker/cbor` v2.9.2,
`veraison/go-cose` v1.3.0, and `x448/float16` v0.8.4, replays the retained object and
cross-language known answers, exercises restoration mutations and trusted-key/binding refusal, and
retains bounded fuzz/resource/footprint results. It selects only object-specific fxamacker typed
encode/decode and records a production NO-GO for go-cose.

The later [v0 object-set and wrapper result](V0_CBOR_OBJECT_SET_AND_WRAPPER.md) freezes only
`SourceManifest` v0 as implementation-eligible, retains plan/registration/approval and the
conditional TypeScript family as pre-freeze, and adds exact fxamacker v2.9.2 plus float16 v0.8.4 to
the root module behind one unwired object-specific codec. Existing authority packages do not
import it. The retained handwritten implementation remains the independent oracle; no signing key,
COSE dependency, consumer, IPC, store, runtime, backend, guest, ADR promotion, or product-control
status is added.

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
- passive field-authority verification for 335 fields across 26 selected pre-freeze targets,
  including nested MJS SourceManifest member coverage without claiming recursive coverage for
  future Source Preparer or plan-v1 objects; and
- archive Slice F1 passive projections, limits/known answers, defensive copies, and deterministic
  eligibility selection; F2's owner-asserted all-hot v1-to-v2 migration and read-only full
  verifier; and F3's one immutable complete-cohort segment publication and atomic activation.
  Retained lookup, v2 authority mutation, second activation, backup/orphan cleanup, adapter calls,
  and consumers remain absent; and
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
  arm64 library build passed while its unit job was cancelled; the post-merge branch/verifier
  mismatch, absent independent review, remaining uncovered code, release obligations, and every
  guest/product admission boundary remain explicit.

Current dependency and priority view:

1. **Next archive slice:** F2 migration/full verification, F3 first-segment activation, F4A
   retained lookup/replay/passive-collision routing, and the
   [F4B atomic mutation result](SUPERVISOR_ARCHIVE_F4B_MUTATION_RESULT.md), and
   [F4C bounded-growth result](SUPERVISOR_ARCHIVE_F4C_GROWTH_RESULT.md) are `PASSED` in their exact
   local scopes. Continue with F5 backup/orphan cleanup and reporting; keep F6 production-engine
   work separate.
2. **Source Validator R3 signed composition passed; product blocked:** Accepted ADR-0036 closes R0 with
   two role-specific private launchers, residual private-container scratch plus mandatory cleanup,
   and an evidence-derived reactive watermark with no hard-peak or host-availability claim. Product
   R1 passive v1 contracts/field authority, R2 unsigned role-specific construction, and the exact
   Apple Development R3 signed/installed/refusal/cleanup composition are `PASSED`; product work
   remains `BLOCKED`. Continue with the R4 measured confinement/resource/residue corpus and then
   consumers.
   Independently freeze the signed-object set, narrow `fxamacker/cbor` wrapper responsibilities,
   and maintain exact documentation plus field authority.
3. **Installation plan:** I0, I1A, I1B, I2A, I2B1 passive request/record objects, I2B2 unsigned
   construction, and the post-I1B platform-semantics research are `PASSED` in their exact scopes.
   Continue only after production wrapper review with separately authorized I2B3 signed
   Coordinator/bootstrap, protected-container, and descriptor-relative fixed-v1 evidence. Manual whole-bundle replacement
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
   Its documentation slice is `PASSED`; it changes no credential or product admission. It requires
   one explicit W4-versus-3DDR decision, final bundle/App Group topology, and matching selected-team
   certificates/profiles before R3/I1, and defers Developer ID/notarization/CI release custody to
   separately authorized work.
   A genuinely independent
   Linux/arm64 builder is viable but not currently planned; same-host/GitHub-CI equality remains
   limited and independent-builder equality is deferred.
7. **Next Source Validator path:** accepted ADR-0034's M1 bytes, Accepted ADR-0035's passive V0
   frames, and the bounded V1/V2 evidence checkpoints are retained unchanged. Accepted ADR-0036
   closes R0 with two role-specific private launchers, residual private-container scratch plus
   mandatory cleanup, and an evidence-derived reactive watermark with no hard-peak/availability
   claim. Its exact [R3 execution packet](SOURCE_VALIDATOR_R3_EXECUTION_PACKET.md) records Team 3DDR,
   R2 byte identities/placement, profiles/entitlements, refusal/cleanup, and mutation gates. R1
   passive v1 contracts/field authority, R2 unsigned construction, and R3 exact signed installed
   composition are `PASSED`. Product work proceeds sequentially: R4 confinement/
   resource/residue evidence; R5D daemon consumer; R5B Broker consumer; then M2/S1 checkpoint.
   Threshold/cadence/overshoot values remain unset until R4. No S1 fixture or product endpoint is
   authorized.
8. **Next passive contract path:** accepted ADR-0034 removes Source Preparer/plan-v1 from the
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
