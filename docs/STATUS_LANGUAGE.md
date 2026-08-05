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
| Current contract/field-authority corpus | `PASSED` | Generated verification closes 95 rules, 502 cases, 624 fixtures, and 1,111 classified fields across 81 profiles and 50 targets. These are current whole-repository totals; older slice totals remain historical evidence only. |
| `.mjs` parser/process selection | `PASSED` | The bounded parse-only comparison selected exact Oxc 0.140.0 and the disposable-process topology for further implementation. |
| Product Source Validator | `BLOCKED` | R1/R2 and signed inactive R3 passed their exact scopes. Exact R4-v1 candidates are `NO_GO`; R4-v2 was not executed. ADR-0040 moves host AST validation off the internal-alpha critical path while preserving ADR-0035/0036 as later defense-in-depth work. Resume only with a new internally consistent lifetime/residue contract; the accepted private container is residual scratch authority. |
| Owner-only hostile-`.mjs` internal alpha | `IN_PROGRESS — TRENDING_GOOD` | The scope and ordered gates are frozen by ADR-0040. Passive atomic source custody is `PASSED`; authenticated product IPC, Broker render/sign, production approval verification, protected installed state, real adapter/recovery, runtime/profile admission, and completion composition remain `BLOCKED`. See the [readiness map](ALPHA_VERTICAL_FLOW_READINESS.md). |
| Passive authority-plane atomic cutover | `PASSED` | The fixed-store oracle atomically retains one exact `main.mjs`, its manifest, resolved bindings, plan, and registration, then fetches defensive copies by `RegistrationID`. It creates no endpoint, approval, attempt, process, runtime, backend, or guest and admits no product path. |
| Fake-backend authority integration | `BLOCKED` | Passive source/plan bytes, atomic custody/fetch, approval/attempt, lifecycle, and FakeBackend mechanics exist, but no selected product adapter, authenticated facades, Broker flow, production verifier, or compositor connects them. It no longer waits on Source Validator R4/R5. |
| No-guest Supervisor lifecycle E1-E5 and owner-lock G2 local mechanic | `PASSED` | The exact unwired fixed-store/FakeBackend/local-owner scope passes; it creates no guest and is not installed-product evidence. |
| Production-shaped I2B1 CBOR/COSE wrapper review | `PASSED` | The exact passive request/record pairing, repeated-field binding, payload-owned replay, caps, digests, mutations, fuzz targets, and independent CryptoKit check pass against checked-in public-key vectors. It provides no live signing, caller/key authorization, durable replay ledger, installed same-byte consumer, or product admission. |
| Installed owner-lock G3/I2B | `BLOCKED` | Proposed ADR-0038 selects the bootstrap owner and signed-record/store-open contract. I1B, I2B1 passive objects, the narrow production-shaped wrapper review, and I2B2 unsigned construction are `PASSED` in their exact scopes; resume with separately authorized I2B3 exact Team-3DDR signing, caller/key authorization, App Group/service/container handoff, and descriptor-relative fault evidence. |
| Archive fixed-store F2-F5 | `PASSED` | F2 migration/full verification, F3 first-segment activation, F4A retained lookup, F4B atomic authority/lifecycle mutation, F4C bounded growth, and F5 coherent backup/read-only restore admission/explicit known-orphan cleanup/offline reporting pass their exact owner-held local fixed-store scopes. See the [F5 result](SUPERVISOR_ARCHIVE_F5_BACKUP_RESULT.md). |
| Archive F6 and external-alpha storage | `BLOCKED` | ADR-0040 permits only a bounded owner-only fixed-store exception with no restore/continuity claim. F6 production-engine selection, installed APFS/power-loss and restore activation, retention/deletion policy, and external-alpha storage remain blocked. |
| Completed compiled-artifact archive migration | `PASSED` | Bulky completed V1/V2/R2, I1A, R3, and I2B2 payload/evidence files are pinned to exact capsule-experiments archive commit `0944ffd8cfd01ec23e4ae99138b0931d56804077`; Capsule retains compact conformance metadata and required deterministic I2B2 source inputs. Product, signing, installation, runtime, and admission status are unchanged. |
| Governed `deno_core` C1 passive composition | `PASSED` | Exact governed construction identities and the intended `.mjs` JSON-in/JSON-out, runtime-surface, logical-descriptor, resource-reference, and refusal contracts pass as zero-effect fixtures; no runtime or guest was created. |
| Governed `deno_core` C2A passive execution-profile preparation | `PASSED` | The unchanged C1 binding, exact numeric descriptor map, bounded candidate machine/transport/teardown values, known answer, artifact blockers, and complete C2B/restoration matrix pass as zero-effect refusing fixtures. C2B remains separately authorized and blocked. |
| Governed `deno_core` C2B passive bindings | `PASSED` | Immutable v1/v2 preserve their historical checkpoints. The 18,357-byte v3 successor binds current accepted Deno/`rusty_v8`/libkrun heads, resolves runner/libkrunfw/kernel roles, fixes exact FD/port/device/runtime/resource/teardown semantics, and retains a 128-field zero-effect contract. The retained dylib predates accepted libkrun source and the preflight is not a runner; current-source dylib, final runner, guest evidence, and admission remain `BLOCKED`. |
| Governed dependency fork acceptance/default transition | `PASSED` | Deno `capsule/accepted-v2.9.4-r3`, `rusty_v8` `capsule/accepted-v150.2.0-r5`, and libkrun `capsule/upstream-v1.19.4-r3` are locked at their verified merged heads and are their fork defaults. Historical refs remain protected and each `main` remains protected mutable integration state. This changes repository governance only. |
| Governed `deno_core` runtime engineering | `IN_PROGRESS — TRENDING_GOOD` | Physical omission, same-host reproduction, standalone root, fork-native Linux/arm64 reconstruction, governed-fork promotions, C1, passive C2A, and C2B v1-v3 passed in their exact scopes. A current-source libkrun dylib and final runner must close before a separately authorized fixed-owned-guest experiment. |
| Runtime/profile admission | `BLOCKED` | It needs governed release review, exact C2 external-isolation/descriptor/resource/transport evidence, the remaining P0 controls, and final installed/profile evidence. |
| Governed libkrun source and local library hardening | `IN_PROGRESS — TRENDING_GOOD` | The governed r3 backport, immutable default transition, and several bounded library/custody/transport slices passed; independent product-admission review, real guest transport, launcher, installed composition, and final profile reruns remain. |
| macOS installation and distribution | `IN_PROGRESS — TRENDING_GOOD` | I0's exact inactive one-app/seven-role profile, I1A's unsigned construction, I1B's exact development-signed execution-disabled composition, I2A's protected-root decision, I2B1 passive objects, I2B2 unsigned eight-role construction, and the post-I1B platform research are `PASSED` in their exact scopes. Installed I2B-I4 protected-root/IPC/manual-replacement evidence remains `BLOCKED`; automatic updates, Developer ID distribution, a support-floor matrix, and complete uninstall are later gates. |
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
