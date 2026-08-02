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
| P0-2 `NullFs` disposition replacement | `019fc2e8-445e-7cb2-b4c2-54d84282c3fe`, replacing task `019fc2e6-cf9d-7210-b2f3-f3bf2244e83a` after an app-layer failure | The smallest removal built but prevented guest bootstrap; no independent feature toggle or retained sanitizer/coverage corpus closes the residual path. The exact block-root profile remains unsupported and barred from user bytes. | [`NULLFS_P0_2.md`](../experiments/gate-c-libkrun-adversarial/NULLFS_P0_2.md), audit script, patch, and selected evidence |

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
- construction-level Bun runtime-authority closure;
- a distinct trusted launcher and exact host/guest descriptor allowlists;
- hostile virtio-console control, queue, descriptor, backpressure, shutdown, and partial-write
  coverage; and
- early topology testing followed by final signed/notarized-byte rebuild and revalidation.

The authoritative remaining P0 and P1 campaigns are in the Gate C synthesis and roadmap. The exact
stock-runtime branch of P0-0 is now closed as a failure; its governed patch/external mechanism
branch remains open. The reviews and experiment found no architecture-killing contradiction and
did not authorize libkrun to handle user bytes.

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

## Integrated pull-request checkpoints

| Integration | Merge commit | What became durable |
| --- | --- | --- |
| PR #8, Gate C feasibility and readiness | `633d249` | Five retained libkrun readiness experiments, Gate C synthesis, ADR-0022 updates, roadmap, and evidence-matrix reconciliation. |
| PR #9, Phase 2A contract foundation | `3a75098` | Passive `JobProposal`, minimum `ExecutionPlan` and `PlanRegistration` candidates, byte-exact fixtures, role-specific Go/TypeScript views, and parallel-review synthesis. |
| PR #10, Phase 2B boundary decisions | `e565e6b` | Proposed ADR-0023 and an implementation-ready breakdown for strict decoding, exact identity, registration, and conformance. |
| PR #11, Phase 2B conformance foundation | `f6de7ec` | Closed manifest, integrity verifier, 37 rules, 105 cases, and 91 unique raw/media/scalar/CBOR fixtures; corrected derived JSON-node and `PlanRegistration` item limits. |
| Local Phase 2B Task 2.3 integration | `4afbdfa` | Proposal/source/input rules, fixed resolver contexts, and known-answer source-manifest/canonical-input fixtures; corpus expanded to 55 rules, 147 cases, and 130 fixtures. |
| Local Gate C P0-2 evidence integration | `d3be865` | Retained the bounded `NullFs` route/removal investigation and kept the exact profile unsupported. |
| Local Gate C P0-0 evidence integration | `004c047` | Rejected stock Bun 1.3.14 for `RUNTIME-001` and retained the source/probe evidence and governed-patch-or-alternate decision. |
| Local Phase 2B Task 2.4 integration | `fd51ac4` | Exact plan/registration/domain/state fixtures and proposed ADR-0023 state addendum; corpus expanded to 67 rules, 206 cases, and 278 fixtures. |

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
  plan/registration/domain/state fixtures.

Next backend-independent work:

1. Unwired strict proposal decoding and semantic planning.
2. Bounded internal wrappers and exact-byte Supervisor registration.
3. The registered-plan/fault-injectable fake-backend vertical slice.
4. Only then, the coordinated public consumer cutover and mixed-`Job` removal.

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
