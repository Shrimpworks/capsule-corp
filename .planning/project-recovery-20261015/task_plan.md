# Project recovery execution index

Date: 2026-10-15
Work item: RECOVERY-20261015
Status: PASSED for initial recovery deliverables and reviewed repairs published in PR #384
Scope: Recover contributor readiness, implementation reality, reproducible baseline and a prioritized next execution slice in Capsule. No guest, installed service, signing, authority expansion or cross-repository implementation.
Evidence: Clean main at fc8af2d; isolated codex/project-recovery-20261015 worktree created.
Remaining work outside this scope: remote CI/merge and next bounded roadmap assignments.
Next action: Complete PR acceptance, then dispatch the next scoped roadmap item.
Parent: Owner-only internal alpha remains IN_PROGRESS — TRENDING_GOOD; product admission BLOCKED.

Canonical requirements: AGENTS.md, docs/PROJECT.md, docs/CURRENT_WORK_PLAN.md, docs/STATUS_LANGUAGE.md, docs/ALPHA_VERTICAL_FLOW_READINESS.md. This index coordinates recovery; it does not replace product truth.

| ID | Outcome / owner | Dependency | Acceptance / verification | Status |
| --- | --- | --- | --- | --- |
| R0 | Tooling readiness / internal readiness agent | none | Exercise Serena symbols/references, graph coverage/source parity, AI Central provenance and links; identify limitations | PASSED |
| R1 | Recovery assessment / internal rediscovery agent | none | Reconcile implemented paths, canonical plans and evidence; verified/inferred/unresolved labels | PASSED |
| R2 | Baseline / internal baseline agent | none | Required AGENTS.md commands, versions, exit codes, failures and scope limitations | PASSED |
| R3 | Recovery roadmap / Brain | R0–R2 | Prioritized dependency-aware tasks with scope, tests, acceptance, risk, effort and owners | PASSED |
| R4 | Go security pin and AI steering repair / baseline and readiness workers | R3 | Exact old/new scan; red/green steering tests; fresh review max three attempts | PASSED for local implementation and independent review |
| R5 | Integration / Brain | R4 | Reviewed documentation, required checks, commit, push, PR body/state readback | PASSED for review, commit/push and PR body/state readback |

All children are internal subagents sharing this integration destination. They create no branches, commits or PRs. Read-only findings are staged in /tmp before reconciliation. Brain waits for all handoffs. Review failures stop after three attempts and trigger a focused research/alignment decision; no silent cycle reset.

Completion boundary: Initial recovery deliverables and one evidence-selected bounded slice delivered through a PR. Sustained product implementation remains an explicit follow-on workstream; recovery does not imply alpha admission or authorize a gated experiment run.

R0/R1/R2 handoffs collected. Priorities changed on verified CI advisory failure and steering omission: implement R4A/R4B; defer optional JSDoc convenience. Detailed roadmap and evidence: docs/PROJECT_RECOVERY_2026_10_15.md.

Integration destination: https://github.com/Shrimpworks/capsule-corp/pull/384, ready, base main. Implementation commit 22aad310d60009642097c21cf4b822b13bd14069. Remote CI/merge remains a separate gate. All child handoffs collected; no orphan child commits.
