# Work status language

Capsule reports work progress with four statuses. These labels answer one question only: **what is
the state of this exact work item?** They do not describe an ADR's lifecycle, a control's evidence
level, or the readiness of a larger parent system.

## The four statuses

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
| Single-file `.mjs` byte and `SourceManifest` foundation | `PASSED` | Exact passive bytes, caps, canonical manifests, and Go/TypeScript fixtures meet the M1 scope. |
| `.mjs` parser/process selection | `PASSED` | The bounded parse-only comparison selected exact Oxc 0.140.0 and the disposable-process topology for further implementation. |
| Product Source Validator and downstream M2/S1 consumers | `BLOCKED` | V0, the unwired V1 artifact, the negative V2 checkpoint, replacement research, and accepted R0 architecture each retain their bounded status. ADR-0036 selects two role-specific private App-Sandboxed launchers, accepts their containers only as residual scratch with mandatory cleanup, and replaces the unavailable hard ceiling with a later evidence-derived reactive watermark that makes no hard-peak or host-availability claim. Resume sequentially with passive v1 contracts/field authority, unsigned construction, separately authorized signing/install, confinement/resource/residue corpus, daemon consumer, Broker consumer, then M2/S1. V0/V1/V2 bytes never change. |
| No-guest Supervisor lifecycle E1-E5 and owner-lock G2 local mechanic | `PASSED` | The exact unwired fixed-store/FakeBackend/local-owner scope passes; it creates no guest and is not installed-product evidence. |
| Installed owner-lock G3 | `BLOCKED` | Resume when a certificate/profile set with matching Team identity and a selected protected-root bootstrap, signed record, and descriptor-relative store-open composition are available. |
| Archive fixed-store F2 migration, F3 first-segment activation, and F4A read-only lookup | `PASSED` | F2's owner-asserted migration/full verifier, F3's sealed first-segment transaction, and F4A's retained-global lookup/replay/passive-collision/hot-only-recovery routing pass their exact local scopes, including missing-history preservation, publish-before-reference ordering, atomic old-or-new reopen, complete visible tombstones, hot/archive semantic equality, and retained fault/corruption/substitution/concurrency/restoration oracles. |
| Archive fixed-store F4B atomic mutation | `BLOCKED` | F4A reconstructs and resolves only the lifecycle record's current effect, while ADR-0031 requires every earlier v2 effect tombstone to remain after later operations replace that field. ADR-0031's effect-tombstone-source-of-truth correction now records the required decision: an independent, same-transaction, Supervisor-only effect-tombstone set that reconstruction and `ResolveEffect` must read from directly. F4B remains `BLOCKED` on implementing and verifying that correction; see the [exact blocker](SUPERVISOR_ARCHIVE_F4B_MUTATION_BLOCKER.md). F4C bounded growth and F5-F6 backup/engine work remain deferred. |
| Governed `deno_core` C1 passive composition | `PASSED` | Exact governed construction identities and the intended `.mjs` JSON-in/JSON-out, runtime-surface, logical-descriptor, resource-reference, and refusal contracts pass as zero-effect fixtures; no runtime or guest was created. |
| Governed `deno_core` runtime engineering | `IN_PROGRESS — TRENDING_GOOD` | Physical omission, same-host reproduction, standalone root, fork-native Linux/arm64 reconstruction, and C1 passive composition passed; C2 owns the first authorized composed-profile execution evidence. |
| Runtime/profile admission | `BLOCKED` | It needs governed release review, exact C2 external-isolation/descriptor/resource/transport evidence, the remaining P0 controls, and final installed/profile evidence. |
| Governed libkrun source and local library hardening | `IN_PROGRESS — TRENDING_GOOD` | The governed fork and several bounded library/custody/transport slices passed; independent review, real guest transport, launcher, installed composition, and final profile reruns remain. |
| macOS installation and distribution | `IN_PROGRESS — TRENDING_GOOD` | I0's exact inactive one-app/seven-role profile, generated fixtures, field authority, and pure refusal/state classifiers are `PASSED`. I1-I4 developer-signed setup/protected-root/IPC/manual-replacement work remain; automatic updates, Developer ID distribution, a support-floor matrix, and complete uninstall are later gates. |
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
