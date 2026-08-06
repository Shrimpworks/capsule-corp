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
| Current contract/field-authority corpus | `PASSED` | Generated verification closes 95 rules, 502 cases, 624 fixtures, and 1,159 classified fields across 91 profiles and 57 targets. These are current whole-repository totals; older slice totals remain historical evidence only. |
| `.mjs` parser/process selection | `PASSED` | The bounded parse-only comparison selected exact Oxc 0.140.0 and the disposable-process topology for further implementation. |
| Product Source Validator | `BLOCKED` | R1/R2 and signed inactive R3 passed their exact scopes. Exact R4-v1 candidates are `NO_GO`; R4-v2 was not executed. ADR-0040 moves host AST validation off the internal-alpha critical path while preserving ADR-0035/0036 as later defense-in-depth work. Resume only with a new internally consistent lifetime/residue contract; the accepted private container is residual scratch authority. |
| Owner-only hostile-`.mjs` internal alpha | `IN_PROGRESS — TRENDING_GOOD` | The scope and ordered gates are frozen by ADR-0040. Passive atomic source custody plus ADR-0043 Broker projection/public-key verification are `PASSED`; authenticated product IPC, installed Broker UI/signing/key authorization, protected installed state, real adapter/recovery, runtime/profile admission, and completion composition remain `BLOCKED`. See the [readiness map](ALPHA_VERTICAL_FLOW_READINESS.md). |
| Passive authority-plane atomic cutover | `PASSED` | The fixed-store oracle atomically retains one exact `main.mjs`, its manifest, resolved bindings, plan, and registration, then fetches defensive copies by `RegistrationID`. It creates no endpoint, approval, attempt, process, runtime, backend, or guest and admits no product path. |
| Authenticated-local-IPC S3 native-contract prerequisite | `PASSED` | Exact XPC dictionary keys/types, numeric method/status/reason tags, versions, key counts, caps, and cross-language refusal fixtures are frozen for `SubmitMainMJSV0`, `RegisterPlanV0`, and `GetRegisteredPlanV0`. The native harness, listener, peer authentication, installed profile, consumers, and product IPC remain `BLOCKED`. |
| Passive Broker rendering and approval verification v0 | `PASSED` | Accepted ADR-0043 freezes an ASCII-safe read-only projection over exact Supervisor-retained plan/bindings/registration/manifest/source bytes, the Secure Enclave/user-presence key contract, and a strict Capsule-owned public-key-only ApprovalGrant COSE verifier. Inline-input content is not present or shown. No UI, Keychain, LocalAuthentication, private key, signer, IPC, activation, runtime, backend, or guest exists; the installed parent remains `BLOCKED`. |
| Public-key approval to FakeBackend authority integration | `PASSED` | One public-only signed grant bound to the exact retained plan now traverses the production-shaped verifier, durable replay/consume/create store, and `AttemptID`-only no-guest lifecycle. Plan A/B substitution refuses; equivalent signatures, commit response loss, reopen, lifecycle interruption, and terminal replay converge without duplicate fake effects. No live signer, UI, IPC, protected product state, real adapter/backend, runtime, guest, or product admission exists. See the [integration result](APPROVAL_FAKE_LIFECYCLE_INTEGRATION.md). |
| No-guest Supervisor lifecycle E1-E5 and owner-lock G2 local mechanic | `PASSED` | The exact unwired fixed-store/FakeBackend/local-owner scope passes; it creates no guest and is not installed-product evidence. |
| Passive durable completion-last transaction | `PASSED` | The version-0 fixed-file/FakeBackend oracle atomically retains one immutable completion, bounded transcript, and fixed summary only after retained terminal lifecycle and fake authoritative absence. EOF/exit zero, stale replay, forged or mixed state, response loss, restart, and unresolved teardown cases fail closed or converge. Real runner identity, process-tree absence, product-store integration, signing, and admission remain `BLOCKED`. |
| Production-shaped I2B1 CBOR/COSE wrapper review | `PASSED` | The exact passive request/record pairing, repeated-field binding, payload-owned replay, caps, digests, mutations, fuzz targets, and independent CryptoKit check pass against checked-in public-key vectors. It provides no live signing, caller/key authorization, durable replay ledger, installed same-byte consumer, or product admission. |
| Installed owner-lock G3/I2B | `BLOCKED` | Proposed ADR-0038 selects the bootstrap owner and signed-record/store-open contract. Proposed ADR-0045 now selects a separate versioned Supervisor-authority-epoch candidate; its passive decision and inert experiment packet are `PASSED`, not OS evidence. The stable Supervisor identity is legacy residue. Resume only after separately authorized Apple Development identity-separation mutations pass, then separately authorize the Keychain/service/root corpus. |
| Archive fixed-store F2-F5 | `PASSED` | F2 migration/full verification, F3 first-segment activation, F4A retained lookup, F4B atomic authority/lifecycle mutation, F4C bounded growth, and F5 coherent backup/read-only restore admission/explicit known-orphan cleanup/offline reporting pass their exact owner-held local fixed-store scopes. See the [F5 result](SUPERVISOR_ARCHIVE_F5_BACKUP_RESULT.md). |
| Archive F6 and external-alpha storage | `BLOCKED` | ADR-0040 permits only a bounded owner-only fixed-store exception with no restore/continuity claim. F6 production-engine selection, installed APFS/power-loss and restore activation, retention/deletion policy, and external-alpha storage remain blocked. |
| Completed compiled-artifact archive migration | `PASSED` | Bulky completed V1/V2/R2, I1A, R3, and I2B2 payload/evidence files are pinned to exact capsule-experiments archive commit `0944ffd8cfd01ec23e4ae99138b0931d56804077`; Capsule retains compact conformance metadata and required deterministic I2B2 source inputs. Product, signing, installation, runtime, and admission status are unchanged. |
| Governed `deno_core` C1 passive composition | `PASSED` | Exact governed construction identities and the intended `.mjs` JSON-in/JSON-out, runtime-surface, logical-descriptor, resource-reference, and refusal contracts pass as zero-effect fixtures; no runtime or guest was created. |
| Governed `deno_core` C2A passive execution-profile preparation | `PASSED` | The unchanged C1 binding, exact numeric descriptor map, bounded candidate machine/transport/teardown values, known answer, artifact blockers, and complete C2B/restoration matrix pass as zero-effect refusing fixtures. C2B remains separately authorized and blocked. |
| Governed `deno_core` C2B passive bindings, source contract, and v4 materialization | `PASSED` | Immutable v1/v2/v3 and the 3,996-byte source contract preserve their exact scopes. V4 retains the accepted header, twice-reproduced current-source unsigned libkrun dylib, independent ABI audit, byte-equal unsigned final runner, and composed digest. Its verifier freezes FDs 0–7, close-from 8, three ports, explicit implicit-device disable calls, exact call/import order, no replacement authority, and external teardown. No runner/libkrun execution, HVF, VM, guest, signing, install, consumer, or admission occurred in v4's immutable scope. |
| First fixed benign owned guest v19 | `PASSED` | One separately authorized experimental successor booted the exact fixed governed-`deno_core` known answer in one owned-disposable libkrun/HVF guest, matched the full bounded console digest proof, exited normally, was reaped without force-kill, and completed unlinked-root teardown. It used immutable in-root inputs and a diagnostic console proof, not arbitrary source or the final typed transport. See the [checkpoint](FIRST_OWNED_GUEST_EXECUTION_CHECKPOINT.md). |
| Fixed denial-test v20 no-launch materialization | `PASSED` | Independent A/B construction reproduced the exact network-disabled root, guest binary, signed runner, profile, and controller; strict signature, digest, profile, source/control/sink, and no-network-call checks passed. No runner, libkrun/HVF process, or guest was launched. Dynamic v20 execution is separately `BLOCKED` on fresh exact one-use authorization. |
| Governed dependency fork acceptance/default transition | `PASSED` | Deno `capsule/accepted-v2.9.4-r3`, `rusty_v8` `capsule/accepted-v150.2.0-r5`, and libkrun `capsule/upstream-v1.19.4-r3` are locked at their verified merged heads and are their fork defaults. Historical refs remain protected and each `main` remains protected mutable integration state. This changes repository governance only. |
| Governed `deno_core` runtime engineering | `IN_PROGRESS — TRENDING_GOOD` | Physical omission, same-host reproduction, standalone root, fork-native Linux/arm64 reconstruction, governed-fork promotions, C1, passive C2A, C2B v1-v4, the exact v19 fixed-owned-guest experiment, and v20 no-launch materialization passed in their scopes. V20 dynamic execution remains separately authorization-bound; typed transport, hostile-source composition, installed evidence, and admission remain. |
| Runtime/profile admission | `BLOCKED` | It needs governed release review, exact C2 external-isolation/descriptor/resource/transport evidence, the remaining P0 controls, and final installed/profile evidence. |
| Governed libkrun source and local library hardening | `IN_PROGRESS — TRENDING_GOOD` | The governed r3 backport, immutable default transition, bounded library/custody/transport slices, and one exact fixed benign guest composition passed. The final typed transport/launcher profile, denial and hostile corpora, independent product-admission review, installed composition, and final profile reruns remain. |
| macOS installation and distribution | `IN_PROGRESS — TRENDING_GOOD` | Proposed ADR-0045 closes the design choice with versioned Supervisor/Coordinator identities, private container, LaunchAgent, bootstrap group/service, role-specific Keychain groups, descriptor, and state-engine binding. The decision and inert matrix are `PASSED`; installed I2B-I4 remains `BLOCKED` pending exact Apple Development identity-separation and later key/service/root evidence. Automatic updates, Developer ID distribution, a support-floor matrix, and complete uninstall are later gates. |
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
