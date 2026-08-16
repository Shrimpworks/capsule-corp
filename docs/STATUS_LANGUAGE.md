# Work status language

Capsule reports work progress with five statuses. These labels answer one question only: **what is
the state of this exact work item?** They do not describe an ADR's lifecycle, a control's evidence
level, or the readiness of a larger parent system.

## The five statuses

| Status | Use it when | Required detail |
| --- | --- | --- |
| `PASSED` | The exact scoped item met every declared acceptance condition. | Name the scope that passed and any larger parent item that remains incomplete. |
| `IN_PROGRESS — TRENDING_GOOD` | Work is active and the latest evidence is closing risk or blockers. | Name the evidence gained, remaining work, and next action. |
| `IN_PROGRESS — TRENDING_BAD` | Work is active, but the latest evidence found new difficulty, expanded the required surface, or produced repeated failures. The path is still being pursued. | Name the adverse evidence, why the path remains open, and the next decision point. |
| `BLOCKED` | The path is still intended, but work cannot proceed until a named dependency, decision, artifact, credential, or contract is available. | Name the blocker, the owner or source of the unblock, and the action that resumes work. |
| `NO_GO` | The exact candidate or path has been abandoned. No further work is planned unless a new explicit decision reopens it. | Name the rejected candidate, why it was abandoned, and the selected replacement or fallback. |

`IN_PROGRESS` without a trend is incomplete reporting. `NO_GO` is never shorthand for “not done,”
“not admitted,” “more testing is required,” “the next slice is missing,” or “currently blocked.”

## Keep status, evidence, and decisions separate

Capsule uses other vocabularies for different questions. Do not combine them into a compound work
status.

| Dimension | Example | Meaning |
| --- | --- | --- |
| Work status | `IN_PROGRESS — TRENDING_GOOD` | Whether the exact work item is complete, active, blocked, or abandoned. |
| ADR lifecycle | `Proposed`, `Accepted`, `Superseded` | Whether an architecture decision is under review, selected, or replaced. |
| Control/evidence state | `unsupported`, `local-mechanic`, `spike-observed`, `proposed` | What the retained evidence supports for one security control. |
| Product admission | admitted or not admitted | Whether exact composed product bytes/profile may handle the stated workload. |

For example, the governed `deno_core` physical-omission slice is `PASSED`. The larger governed
runtime workstream is `IN_PROGRESS — TRENDING_GOOD`. Product runtime admission is still blocked,
and `RUNTIME-001` remains `unsupported`. None of those facts contradicts another.

## Parent and child work

Always report the smallest item being judged and its parent separately when their status differs.

```text
Slice: deno_core physical built-in omission
Status: PASSED
Evidence: the governed build registered and linked only the three required bootstrap ops.

Parent: governed runtime construction and admission
Status: IN_PROGRESS — TRENDING_GOOD
Remaining work: governed release review, C2 external-isolation composition, and the installed
profile corpus.
Control evidence: RUNTIME-001 unsupported
```

A completed research spike can therefore be `PASSED` even when the product work it informs is
`BLOCKED` or `IN_PROGRESS`. Conversely, a spike may pass because it produced a decisive rejection:
the spike itself is `PASSED`, while the exact rejected candidate is `NO_GO`.

## Required reporting format

Every task plan, handoff, roadmap checkpoint, and current-status summary must include:

```text
Work item:
Status: PASSED | IN_PROGRESS — TRENDING_GOOD | IN_PROGRESS — TRENDING_BAD | BLOCKED | NO_GO
Scope:
Evidence or reason:
Remaining work:
Blocker and owner:     # required for BLOCKED
Replacement/fallback: # required for NO_GO
Next action:
Parent status:        # required when the parent differs
```

For `PASSED`, “remaining work” means work outside the passed scope, not an undisclosed failure of
that scope. For `BLOCKED`, the next action must say what makes the item runnable again. For
`NO_GO`, the replacement may be “none” only when the capability itself has been deliberately
removed from scope.

## Current workstream dashboard

This is the concise interpretation of the current planning documents. Detailed evidence and exact
limitations remain in the linked plans, ledger, ADRs, and control matrix.

| Work item | Status | What that means now |
| --- | --- | --- |
| Internal-alpha architecture and release audit | `PASSED` | Five independent audits converged on Accepted ADR-0040: one owner Mac, exact `main.mjs`, bounded inline JSON, human Broker approval, one fresh governed guest per attempt, and explicit internal/external-alpha separation. The product parent remains `IN_PROGRESS — TRENDING_GOOD` and admission remains `BLOCKED`. |
| Single-file `.mjs` byte and `SourceManifest` foundation | `PASSED` | Exact passive bytes, caps, canonical manifests, and Go/TypeScript fixtures meet the M1 scope. |
| Current contract/field-authority corpus | `PASSED` | Generated verification closes 95 rules, 502 cases, 624 fixtures, and 1,203 classified fields across 95 profiles and 60 targets. These are current whole-repository totals; older slice totals remain historical evidence only. |
| `.mjs` parser/process selection | `PASSED` | The bounded parse-only comparison selected exact Oxc 0.140.0 and the disposable-process topology for further implementation. |
| Product Source Validator | `BLOCKED` | R1/R2 and signed inactive R3 passed their exact scopes. Exact R4-v1 candidates are `NO_GO`; R4-v2 was not executed. ADR-0040 moves host AST validation off the internal-alpha critical path while preserving ADR-0035/0036 as later defense-in-depth work. Resume only with a new internally consistent lifetime/residue contract; the accepted private container is residual scratch authority. |
| Owner-only hostile-`.mjs` internal alpha | `IN_PROGRESS — TRENDING_GOOD` | The scope and ordered gates are frozen by ADR-0040. Passive atomic source custody, ADR-0043 Broker projection/public-key verification, and the fixed-file/FakeBackend completion-last transaction are `PASSED`; authenticated product IPC, installed Broker UI/signing/key authorization, protected installed state, real adapter/recovery, runtime/profile admission, and product completion/store/receipt wiring remain `BLOCKED`. See the [readiness map](ALPHA_VERTICAL_FLOW_READINESS.md) and [current work plan](CURRENT_WORK_PLAN.md). |
| Passive proposal, plan, and authority-plane atomic cutover | `PASSED` | TypeScript `decodeJobProposal`, `resolveJobProposal`, and `constructExecutionPlan` implement the unwired proposal-to-plan mechanics. The Go fixed-store `Facade.RegisterPlanV0`/`Facade.GetRegisteredPlanV0` oracle atomically retains one exact `main.mjs`, its manifest, resolved bindings, plan, and registration, then fetches defensive copies by `RegistrationID`. These symbols are passive library/in-process mechanics: they create no authenticated endpoint, protected product state, approval, attempt, process, runtime, backend, or guest and admit no product path. |
| Authenticated-local-IPC S3 native contract and C2b0 harness | `PASSED` | Exact XPC dictionaries, tags, versions, caps, refusals, replay/response-loss, and strict deadline cases remain frozen. The immutable [C2b0 archive](https://github.com/Shrimpworks/capsule-experiments/tree/3d7bd46352506bf6018286749c2c85a3e2f683df/experiments/authenticated-local-ipc-s3-native-xpc-c2b0-ce7220e523bc43ba-c7ae502b0742bab1e) adds a gated, reproducible, unsigned native harness and independent mutation verification. Nothing was signed, launched, registered, or delivered; native OS enforcement, installed profile, consumers, and product IPC remain `BLOCKED`. |
| Authenticated-local-IPC C4 passive evidence claim | `PASSED` | PR #248 is the canonical predecessor. The focused follow-up closes ordered 4,999/5,000/5,001-ms equality semantics, complete closed dictionaries/maps, every ordered case field, required `noState`, all 20 collisions, refusal/cancellation/replay/response-loss oracles, and bounded Go/Node mutation proofs. Installed and product status remain `BLOCKED`. |
| CL4 IPC refusal-matrix audit | `PASSED` | The independent read-only audit retains historical disposition `AMEND` and found no runtime authority bypass. Its focused evidence-hardening findings are now closed; no ADR change follows. |
| Authenticated-local-IPC ADR-0029 S0 decision review | `PASSED` | Accepted ADR-0029 retains one native-fronted Go Supervisor, two role-specific Supervisor services, four ordinary methods, authentication-before-copy, method-specific bridges, Go-only durable authority, correlation-only request IDs, and `AttemptID`-only recovery. The three-method S3 harness and the two-method C4 contract are separate lanes; platform enforcement and product activation remain `BLOCKED`. |
| Native XPC enforcement research | `PASSED` | Apple primary sources and macOS 26.5 SDK declarations support a controlled low-level `xpc_connection_t` harness with a peer requirement installed before listener activation, exact-message `SecCode` validation, connection-time EUID/ASID checks, non-preemptive cancellation, and protocol-owned deadlines/recovery. The harness, installed identities, services, consumers, and product IPC remain `BLOCKED`. |
| Typed guest transport research, C5b0-C5b9 no-run construction, and compatibility preflight | `PASSED` | The passive R2/C5a design, [C5b0 packet](https://github.com/Shrimpworks/capsule-experiments/tree/b357d0c0fb29100c180494e67cebd7809aabe3c5/experiments/typed-guest-transport-c5b0-v19-successor), [C5b1 construction](https://github.com/Shrimpworks/capsule-experiments/tree/db08ebf277432e06d6cba3b7f7338e3bd4a61252/experiments/typed-guest-transport-c5b1-executable-successor), [C5b2 input closure](https://github.com/Shrimpworks/capsule-experiments/tree/5a2f835e8c9df8279237f940f5af757e119593bd/experiments/typed-guest-transport-c5b2-governed-input-closure), [C5b7 runtime root](https://github.com/Shrimpworks/capsule-experiments/tree/78485fb91a31733c568fe43e5fa295474e5956e1/experiments/typed-guest-transport-c5b7-deterministic-runtime-root), [C5b8 root-binding successor](https://github.com/Shrimpworks/capsule-experiments/tree/b0819d76883eb86cbbc03b2b7033fe55bedbf713/experiments/typed-guest-transport-c5b8-c5b7-root-binding-successor), [C5b9 immutable composite](https://github.com/Shrimpworks/capsule-experiments/tree/3965e6b5cc87d476da7f431d7ed8a5758011a1b8/experiments/typed-guest-transport-c5b9-immutable-no-run-composite), and [compatibility preflight](https://github.com/Shrimpworks/capsule-experiments/tree/7fc3af9c46895b340c3118a96cb50abb26b1d977/experiments/typed-guest-transport-c5b-controlled-harness-preflight) remain exact. The later slices retain exact `rusty_v8`, controller, libkrunfw, descriptive-adapter, governed fixed-fixture Deno, 100,663,296-byte root, sealed test-double effect sequencing, root-compatible object, six-role composite identities, and the exact four-way incompatibility result with closed static/mutation evidence. No retained artifact was loaded or executed. Direct provider-only C5b9 composition is `NO_GO`; a versioned fixed-runner successor, controlled execution, preferred-form libkrunfw/kernel source compliance, installed composition, and admission remain `BLOCKED`. The extracted kernel is evidence-only and separate firmware is inapplicable under ADR-0041. |
| Post-checkpoint stabilization review | `PASSED` | Targeted authority/approval/completion/archive review fixed one current-checkout worktree-isolation defect, completed the `lifecyclestate` exported-contract documentation batch, and retained a race-suite scalability limitation. Required ordinary tests/build/vet and current-change lint pass; full legacy exported-comment lint remains `BLOCKED` on issue #217. |
| Passive Broker rendering and approval verification v0 | `PASSED` | Accepted ADR-0043 freezes an ASCII-safe read-only projection over exact Supervisor-retained plan/bindings/registration/manifest/source bytes, the Secure Enclave/user-presence key contract, and a strict Capsule-owned public-key-only ApprovalGrant COSE verifier. Inline-input content is not present or shown. No UI, Keychain, LocalAuthentication, private key, signer, IPC, activation, runtime, backend, or guest exists; the installed parent remains `BLOCKED`. |
| R3 Approval Broker live-signing research | `PASSED` | The [canonical passive brief](BROKER_LIVE_SIGNING_EVIDENCE_BRIEF.md) maps the Apple-supported fresh-context/key-use mechanism, preserves Supervisor-owned durable `SubmitApprovalV0` authority, and defines experiment-only accessibility/raw-signature candidates and fault oracles. It remains research rather than installed evidence. |
| C6b1a/b Broker/Supervisor evidence construction | `PASSED` | The immutable [unsigned Broker harness](https://github.com/Shrimpworks/capsule-experiments/tree/4a2447d4bd0e03132dc616e608031ca313630cdd/experiments/broker-live-signing-c6b1) retains deterministic public fixtures, public-only checks, entitlement requests, and no-credential interaction tests. The immutable [test-only Supervisor seam](https://github.com/Shrimpworks/capsule-experiments/tree/067fe2beb40361bb714507cab1331004e0a656fa/experiments/broker-live-signing-c6b1-supervisor-seam) retains six commit/replay/response-loss/reopen/concurrency rows with no Broker durable authority. No installed identity, Keychain, LocalAuthentication, private key, installation, authenticated listener, product store/consumer, runtime, backend, VM, or guest participated. |
| C6b1c Broker identity/profile and signed-artifact readback | `PASSED` | The immutable [no-install result](https://github.com/Shrimpworks/capsule-experiments/tree/82d1a799f70482856aaa6030f612d701b39cec67/experiments/broker-live-signing-c6b1c-signed-artifact-readback) retains exact development-profile metadata, signed Broker bytes, strict signature/designated requirement, Team/bundle/CDHash, hardened runtime, and the closed App Sandbox plus one Approval Keychain-group entitlement. The app was not installed or launched, the raw profile was not embedded, and no Keychain/LocalAuthentication/service/product/runtime/guest operation occurred. C6b1d and product wiring remain `BLOCKED`. |
| Public-key approval to FakeBackend authority integration | `PASSED` | One public-only signed grant bound to the exact retained plan now traverses the production-shaped verifier, durable replay/consume/create store, and `AttemptID`-only no-guest lifecycle. Plan A/B substitution refuses; equivalent signatures, commit response loss, reopen, lifecycle interruption, and terminal replay converge without duplicate fake effects. No live signer, UI, IPC, protected product state, real adapter/backend, runtime, guest, or product admission exists. See the [integration result](APPROVAL_FAKE_LIFECYCLE_INTEGRATION.md). |
| No-guest Supervisor lifecycle E1-E5 and owner-lock G2 local mechanic | `PASSED` | The exact unwired fixed-store/FakeBackend/local-owner scope passes; it creates no guest and is not installed-product evidence. |
| Passive durable completion-last transaction | `PASSED` | The version-0 fixed-file/FakeBackend oracle atomically retains one immutable completion, bounded transcript, and fixed summary only after retained terminal lifecycle and fake authoritative absence. EOF/exit zero, stale replay, forged or mixed state, response loss, restart, and unresolved teardown cases fail closed or converge. Real runner identity, process-tree absence, product-store integration, signing, and admission remain `BLOCKED`. |
| Production-shaped I2B1 CBOR/COSE wrapper review | `PASSED` | The exact passive request/record pairing, repeated-field binding, payload-owned replay, caps, digests, mutations, fuzz targets, and independent CryptoKit check pass against checked-in public-key vectors. It provides no live signing, caller/key authorization, durable replay ledger, installed same-byte consumer, or product admission. |
| Installed owner-lock G3/I2B | `BLOCKED` | Proposed ADR-0038 selects the bootstrap owner and signed-record/store-open contract. Proposed ADR-0045 selects a versioned Supervisor-authority-epoch candidate, and the immutable [C3a E0 archive](https://github.com/Shrimpworks/capsule-experiments/tree/dee784d40684100f8315720fab9a5cd3399f492b/experiments/macos-installation-i2b3-supervisor-authority-epoch-e0) passed exact unsigned reproducible construction. The exact legacy profile is restored outside Git. The [C3b correction](https://github.com/Shrimpworks/capsule-experiments/tree/3671a6eb23357ff28de4562dd60e8f68173034ae/experiments/macos-installation-i2b3-supervisor-authority-epoch-e1-app-group-preflight) makes only Developer-portal registration of the frozen macOS-style App Group `NO_GO`; the identity remains intended. The [C3b signed-profile preflight](https://github.com/Shrimpworks/capsule-experiments/tree/ee00ae2abbce64ae6458b82d0b53d904ee39aeb6/experiments/macos-installation-i2b3-supervisor-authority-epoch-e1-signed-profile-preflight) is `PASSED` for exact no-launch profile/signature readback. The separately authorized E1 container matrix remains `BLOCKED`; ADR-0038/0045 remain Proposed. After C3b, C3c must adopt, amend, reject, or supersede both decisions and freeze the descriptor/state-engine binding before C6a may begin. |
| Archive fixed-store F2-F5 | `PASSED` | F2 migration/full verification, F3 first-segment activation, F4A retained lookup, F4B atomic authority/lifecycle mutation, F4C bounded growth, and F5 coherent backup/read-only restore admission/explicit known-orphan cleanup/offline reporting pass their exact owner-held local fixed-store scopes. See the [F5 result](SUPERVISOR_ARCHIVE_F5_BACKUP_RESULT.md). |
| Archive F6 and external-alpha storage | `BLOCKED` | ADR-0040 permits only a bounded owner-only fixed-store exception with no restore/continuity claim. F6 production-engine selection, installed APFS/power-loss and restore activation, retention/deletion policy, and external-alpha storage remain blocked. |
| Completed compiled-artifact archive migration | `PASSED` | Bulky completed V1/V2/R2, I1A, R3, and I2B2 payload/evidence files are pinned to exact capsule-experiments archive commit `0944ffd8cfd01ec23e4ae99138b0931d56804077`; Capsule retains compact conformance metadata and required deterministic I2B2 source inputs. Product, signing, installation, runtime, and admission status are unchanged. |
| Governed `deno_core` C1 passive composition | `PASSED` | Exact governed construction identities and the intended `.mjs` JSON-in/JSON-out, runtime-surface, logical-descriptor, resource-reference, and refusal contracts pass as zero-effect fixtures; no runtime or guest was created. |
| Governed `deno_core` C2A passive execution-profile preparation | `PASSED` | The unchanged C1 binding, exact numeric descriptor map, bounded candidate machine/transport/teardown values, known answer, artifact blockers, and complete C2B/restoration matrix pass as zero-effect refusing fixtures. C2B remains separately authorized and blocked. |
| C5a passive typed guest transport v1 | `PASSED` | Exact source/input/completion layouts, caps and cap-plus-one, big-endian bindings, four closed statuses, canonical JSON, ordered refusals, monotonic state/fault semantics, endpoint custody, completion-last projection, 48 generated frame cases, 13 state cases, and 23 restoration cases agree in Go and Node. No endpoint, runtime, backend, guest, or product consumer exists; controlled C5b execution and admission remain `BLOCKED`. |
| Governed `deno_core` C2B passive bindings, source contract, and v4 materialization | `PASSED` | Immutable v1/v2/v3 and the 3,996-byte source contract preserve their exact scopes. V4 retains the accepted header, twice-reproduced current-source unsigned libkrun dylib, independent ABI audit, byte-equal unsigned final runner, and composed digest. Its verifier freezes FDs 0–7, close-from 8, three ports, explicit implicit-device disable calls, exact call/import order, no replacement authority, and external teardown. No runner/libkrun execution, HVF, VM, guest, signing, install, consumer, or admission occurred in v4's immutable scope. |
| First fixed benign owned guest v19 | `PASSED` | One separately authorized experimental successor booted the exact fixed governed-`deno_core` known answer in one owned-disposable libkrun/HVF guest, matched the full bounded console digest proof, exited normally, was reaped without force-kill, and completed unlinked-root teardown. It used immutable in-root inputs and a diagnostic console proof, not arbitrary source or the final typed transport. Canonical identities remain, but the unpublished raw v10-v27 archive is unavailable and cannot support durable release/admission evidence. See the [checkpoint](FIRST_OWNED_GUEST_EXECUTION_CHECKPOINT.md). |
| Fixed denial-test v20 no-launch materialization | `PASSED` | Independent A/B construction reproduced the exact network-disabled root, guest binary, signed runner, profile, and controller; strict signature, digest, profile, source/control/sink, and no-network-call checks passed. The later exact v20 runtime attempt is `BLOCKED`: its runner exited 125 before readiness, no start authorization or guest launch occurred, and the exact pre-ready stage is unknown because early stderr was not persisted. |
| Fixed denial-test v21 diagnostic materialization and runtime stop | Materialization `PASSED`: A/B runner/profile/controller reproduction, strict signatures, independent composed-digest calculation, C17 warning-as-error builds, Go tests/vet, and false/blocked authority assertions passed. Exact runtime attempt `BLOCKED`: ready EOF arrived before `R`, no start authorization or guest launch occurred, and the controller returned before drained stderr or authoritative waitpid evidence persisted. |
| Fixed denial-test v22 convergence materialization and runtime stop | Materialization `PASSED`: A/B runner/profile/controller reproduction, strict signatures, independent composed-digest calculation, Go tests/vet, clean controller builds, and authority assertions passed for bounded non-`R` convergence. Exact runtime attempt `BLOCKED`: retained stderr identified `preflight-root-sha256`, with authoritative waitpid/teardown/canary evidence and no authorization, libkrun configuration, HVF call, or guest launch. Host-only hashes match the expected root; the inherited-FD/hash-state cause remains unresolved. |
| Fixed denial-test v23 root-digest localization | `PASSED` | Materialization and one exact invocation proved identical staged-path, unlinked-open-FD, runner-computed, and actual-root digests. The runner refused before authorization/libkrun/HVF/guest activity because its embedded expected byte array was malformed from zero-based byte 18 through 31. The hostile denial controls therefore remain unexecuted and `BLOCKED`. |
| Fixed denial-test v24 corrected execution | `PASSED` for corrected preflight, governed-runtime known answer, and early denial controls; complete corpus `BLOCKED` | One exact authorized guest proved non-root/no-new-privileges/zero capabilities, sealed descriptors, root-write denial, absent host paths, mount denial, and root-regain denial. It then failed generically in the vsock-check family before later controls. No connect/send occurred; normal reap, teardown, canary, and network-zero evidence passed. |
| Fixed denial-test v25 runtime candidate | `NO_GO` | Pre-launch semantic review showed `AF_VSOCK` socket creation alone is not usable transport/network authority and can succeed when `VM_SOCKETS_GET_LOCAL_CID` fails. V25 tested the wrong property, so it was neither authorized nor launched. |
| Fixed denial-test v26 consolidated execution | Failure-localization objective `PASSED`; complete corpus `BLOCKED` | One exact authorized guest passed known-answer, active local-CID vsock-unavailable, and raw-block write-open denial, then reported expected down `dummy0` as a non-loopback interface. This is an over-strict probe-policy mismatch, not network-access evidence. No connect/send or traffic occurred; normal reap/teardown/canary evidence passed. Do not rerun v26. |
| Exact fixed denial-test v27 execution | `PASSED` | One exact authorized owned guest completed all 30 fixed markers with exact completion/console proofs: known answer; non-root/no-new-privileges/zero capabilities; sealed descriptors; root/host-path/mount/root-regain denials; active local-CID vsock unavailable without connect/send; raw-block denial; passive network/route denial; absent virtiofs; empty environment; normal reap; unlinked-root teardown; unchanged canary; and zero network/credential authority or traffic. Canonical identities remain, but the unpublished raw reports/receipts/harness/manifest are unavailable; durable release/admission evidence remains `BLOCKED`. |
| Governed dependency fork acceptance/default transition | `PASSED` | Deno `capsule/accepted-v2.9.4-r3`, `rusty_v8` `capsule/accepted-v150.2.0-r5`, and libkrun `capsule/upstream-v1.19.4-r3` are locked at their verified merged heads and are their fork defaults. Historical refs remain protected and each `main` remains protected mutable integration state. This changes repository governance only. |
| Governed `deno_core` runtime engineering | `IN_PROGRESS — TRENDING_GOOD` | Physical omission, same-host reproduction, standalone root, fork-native Linux/arm64 reconstruction, governed-fork promotions, C1, passive C2A, C2B v1-v4, v19 fixed-owned-guest execution, v20-v27 materializations, v23 root-digest localization, v24 early denials, v26 consolidated localization, exact v27 fixed denial execution, and passive C5a transport conformance passed in their scopes. V25 is `NO_GO`; controlled real transport, installed lifecycle/recovery, broader platform matrices, hostile-source composition, product wiring, and admission remain. |
| Runtime/profile admission | `BLOCKED` | It needs governed release review, exact C2 external-isolation/descriptor/resource/transport evidence, the remaining P0 controls, and final installed/profile evidence. |
| Governed libkrun source and local library hardening | `IN_PROGRESS — TRENDING_GOOD` | The governed r3 backport, immutable default transition, bounded library/custody/transport slices, exact fixed benign composition, v24 early denials, v26 active local-CID/raw-block denials, exact v27 fixed denial execution, and passive C5a transport conformance passed. Controlled transport/launcher execution, broader kernel/platform/lifecycle corpora, independent product-admission review, installed composition, and final profile reruns remain. |
| macOS installation and distribution | `IN_PROGRESS — TRENDING_GOOD` | Proposed ADR-0045 closes the design choice with versioned Supervisor/Coordinator identities, private container, LaunchAgent, macOS-style bootstrap group/service, role-specific Keychain groups, descriptor, and state-engine binding. Its inert matrix, C3a reproducible unsigned E0 construction, bounded no-mutation preflights, and exact no-launch profile/signature readback are `PASSED`; installed I2B-I4 remains `BLOCKED` pending exact Apple Development identity-separation and later key/service/root evidence. Automatic updates, Developer ID distribution, a support-floor matrix, and complete uninstall are later gates. |
| TypeScript Source Preparer | `BLOCKED` | It is a conditional later feature, not a first-release dependency. Resume only if TypeScript returns to scope and its protected-store/worker/update/retention authority blockers are closed. |

Current `NO_GO` examples are deliberately narrow: stock Bun 1.3.14, the abandoned broad governed
Bun construction, hardened full Deno as the product runtime, the removed handwritten ECMAScript
scanner, and `go-cose` v1.3.0 as a product dependency. Those exact paths are not being pursued;
their replacements are governed `deno_core`, Oxc in a disposable validator, and a narrow
Capsule-owned COSE_Sign1 wrapper using standard crypto.

## Historical records

Do not rewrite an experiment's quoted historical title or erase the decision it recorded. When an
old title used ambiguous language such as `PHYSICAL-OMISSION-PASS; NO RUNTIME ADMISSION`, preserve
it as historical provenance only and add the current translation:

- scoped slice: `PASSED`;
- parent workstream: `IN_PROGRESS — TRENDING_GOOD` or `BLOCKED`;
- product admission/control evidence: stated separately.

This preserves auditability without making maintainers decode old compound labels to understand
today's plan.
