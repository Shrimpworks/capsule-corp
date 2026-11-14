# Preliminary independent review

Review instance: 1 of 3.
Status: IN_PROGRESS — TRENDING_GOOD. Child execution remains BLOCKED.

Internal read-only review of canonical worktree `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp` on `codex/c5b19-decision-reconciliation` at `8e4fd22cc00cd5e961544f6a803974600dc0e751` and archive `/tmp/capsule-c5b19-recovery` on `codex/c5b19-benign-teardown` at `5d0e53cbccae3d4f587205df70253b776eae63e5`. Both branches/HEADs and scoped dirty files match bootstrap. Unrelated `.serena/project.yml` excluded and preserved.

This ledger precedes reading `implementation-handoff.md`, `ci-race-handoff.md`, archive PLAN/ORIGINS rationale, and retained verification reports. The required plan/scope contained brief implementation results and race explanation; therefore complete blinding to those summaries was impossible, although detailed author explanations remain unread.

## Preliminary findings ledger

- No actionable implementation defect identified in the initial pass.
- Archive `oracle.go:775–875` composes value-owned writeHistory and existing individual pair validation. Exact pending bytes plus private reply validation derive settlement; false recorded decisions refuse before committing. Code2 preserves all history and latches fatal; later settlement retains the latch. Cross-call sequence/tick guards and three-decision cap match scope. Need independent test execution and retained mutation evidence reconciliation.
- Tests at `oracle_test.go:1772–2169` assert whole private-state equality on refusal, independently expected byte projections, late confirmed/failed/indeterminate settlement, valid substitutions, duplicate/stale decisions, boundaries and ownership. Existing helper implementations and validation support a meaningful independent fixture path. Need verify old evidence/source preservation and exact provenance hashes.
- Canonical `passive_product_adapter.go:378–385` treats `context.Canceled` as cleanup only. Current call construction uses a background-rooted timeout, retains cancelMethod privately, and publishes buffered result before deferred cancellation; the blocking result receive is therefore justified. DeadlineExceeded retains existing response-unknown semantics. Tests synchronize ready cases and exercise caller cancellation without releasing admission early. Need execute race/repetition checks.
- Scope explicitly excludes full control admission, external stop, serial-group interleaving, malformed-reply diagnostics, durability, omitted-event detection and native execution. No widened authority or new dependency observed. Existing history begin after fatal is intentionally recording only, not submission admission.
- Delivery is not yet complete: archive pin/publication, canonical metadata, final hosted acceptance and child gates remain separate obligations. Assess final documentation against these limitations.

## Method and limitations

Serena initialized. Codebase-memory Tier2 used list_projects (both pages), symbol search, SettleWrite outbound trace (four callees), snippet, and coverage for four relevant canonical Go paths plus two archive paths. Index is main checkout, generation2026-10-04T22:18:10Z; canonical metadata_match is best-effort, archive outside_project. Exact reviewed worktree source/diffs supply authoritative fallback, without completeness claim.

Reviewed applicable AGENTS, scope/plan, starter-document boundaries, ADR0047, C5b19 reuse row, Gate2 packet, relevant model/types, archive prior history/pair helpers, changes and tests. Only report files are writable review output; no source mutation or native fixture execution authorized.
