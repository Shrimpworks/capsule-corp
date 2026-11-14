# G3a10 preliminary independent review

Review instance: 1 of 3. Status: IN_PROGRESS — TRENDING_GOOD.

Read-only internal assignment; source-first ledger recorded before opening implementation-handoff.md, canonical-verification.md, implementation-evidence.json, or archive PLAN/ORIGINS rationale. The required scope/plan included high-level author verification summaries; those are not treated as verified evidence. No implementation conversation was supplied.

Scope verified: canonical codex/c5b19-worker-sidecars at fe9009d5398a883735219e2a3a33c351c79aa54c, archive codex/c5b19-benign-teardown at 7646d8faabb670ce8b5ff1b213c414d91deda2d9. Dirty paths match bootstrap; unrelated .serena/project.yml excluded. Four-file archive patch appends 98 source lines and 433 test lines plus metadata. Canonical edits are three roadmap documents and untracked planning.

Read: both AGENTS, bounded starter-document scope, ADR0047, reuse row135, packet338–359/395–451, normative passive model/types, scope.md and plan.md, complete new source/test diff, existing frame validators and snapshotHistory retention/query context. Serena initialized. Codebase-memory primary index is old and for another checkout; search for teardownpassive history/snapshot symbols returned no result. Coverage reports canonical metadata_match, planning missing and archive outside_project. Exact current source reads are authoritative; graph absence is not an absence claim.

## Preliminary findings ledger

No actionable defect identified yet.

1. New query requires initialized history and <=3 selections; it strictly decodes the independent received bytes and matches exact private retained request bytes. Duplicate operation check is appropriate because existing history admits each operation once.
2. Presence combinations enforce exactly receive-only, receive+send+kind5, and those plus kind7. Nonnil empty send remains present and fails byte comparison with strict valid reply. Events are independently validated and assigned explicit kind roles.
3. Every supplied send must equal the worker capture reply, and accepted history when present. Optional control request/reply must equal independent slots. Full captured-frame validation precedes rereading event fields.
4. Returned structs contain arrays/scalars only; source makes no receiver/input assignment and every error returns a zero struct. Pending/fatal/accepted state is copied, never promoted by observation.
5. Tests cover all16 optional-member masks, three outcomes, received-only settled evidence, substitutions with structurally valid pairs, malformed lengths/lane/reserved data, historical selections, duplicates/cap, late failure, ownership, and rejection preserving prior accepted outcome.
6. Follow-up verification questions: confirm prior source/test bytes and all prior ORIGINS members are unchanged; frozen hashes match; retained mutant failures reach named assertions rather than compile/setup failure; ordinary offline test/race/vet/build pass. Verify tests import no lifecycle/network authority before running.
7. Potential coverage nuance, not finding: cap guard is redundant with three private entries plus unique operations; do not infer an independent cap mutation kill. Query's lack of chronology is explicit scope, not omission of full-stream validation.
8. Publication/immutable archive pin and full parent Gate3 remain incomplete; review readiness can apply only to this frozen bounded slice, not native execution or product admission.

Defensive scope: owned inert synthetic byte fixtures and ordinary offline compiler/unit processes only; no native fixture, child/guest, signal/wait target, live store/clock/network or external identity/data. No repository edits, mutation execution, commit, push, PR or further reviewer dispatch.
