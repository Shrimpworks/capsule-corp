# C5b19 G3a9 independent preliminary review

Review instance: 1 of 3. Status: IN_PROGRESS — TRENDING_GOOD. Parent full Gate3 remains IN_PROGRESS — TRENDING_GOOD; child execution and later gates/admission remain BLOCKED.

## Findings first

No actionable implementation defect identified in the source-first pass. This is preliminary, not a readiness verdict. Remaining verification concerns: retained mutation failures must reach the named behavioral assertions rather than setup; prior source/test/provenance must remain preserved; exact source hashes and offline checks must corroborate the frozen target.

## Exact scope and independence

Canonical /Users/dsteele/.codex/worktrees/project-recovery/capsule-corp branch codex/c5b19-snapshot-correspondence HEAD/base f3f2e0755c0c284dd28019225bef0cf8482d3cdd. Archive /tmp/capsule-c5b19-recovery branch codex/c5b19-benign-teardown HEAD/base febb88d6421058f048f75319baae7cf9bd5a0cbb. Both verified by git. Review uncommitted archive oracle.go/oracle_test.go, PLAN.md, inputs/ORIGINS.json under experiments/typed-guest-transport-c5b19-benign-teardown and canonical docs/CURRENT_WORK_PLAN.md, docs/PROJECT_RECOVERY_2026_10_15.md, docs/C5B19_BENIGN_TEARDOWN_EXPERIMENT_PLAN.md plus .planning/c5b19-snapshot-correspondence/. Excluded unrelated .serena/project.yml preserved. No commits beyond bases. No repository edits, mutation scripts, native fixture, child/guest, store/clock/network or external target execution.

The bootstrap, scope.md and plan.md were read before author explanation, as requested. Plan.md itself contains implementation-freeze summaries and scope.md includes design rationale; those summaries were visible in the first pass and are testimony, not independent evidence. The implementation-handoff.md, canonical-verification.md, archive PLAN rationale and ORIGINS.g3a9 have not yet been read. This ledger was recorded before them.

## Source-first observations

- oracle.go 875–1013 adds a bounded value-owned wrapper with three copied request/reply entries. Both transition methods use a whole-state trial and commit only after delegated success and wrapper bookkeeping.
- inspect strictly decodes a supplied record, matches an actually retained request, reports accepted outcome separately from worker outcome and retains historical entry order without asserting physical snapshot timing.
- Absence permits empty record/no worker only and returns unchanged copied history. Query results and fields are values. Record/worker presence does not settle pending state or erase fatal.
- Exact worker request equality and accepted reply equality close substitution and contradictory outcome cases; existing strict decoders constrain lengths before fixed-size conversions. Explicit digest checks overlap existing byte correspondence and are not independently necessary guards.
- Ten new test functions exercise nil/context/absence, pending worker uncertainty, all accepted outcomes, historical candidates, unsubmitted valid anchors, malformed evidence, code2 preservation, copy/nonmutation, worker-request substitution and bounds/chronology. Tests use public workflow-generated state rather than corrupting private internals.
- Normative Model.SettleWrite preserves confirmed projection on failed/indeterminate outcomes and updates only on confirmation, consistent with this wrapper. ADR0047 and packet253–271/342–359 keep observed bytes, accepted settlement, custody and authority distinct.

## Retrieval and next checks

Serena initialized. Used codebase-memory skill Tier2 index_status, search_graph, trace_path, get_code_snippet and check_index_coverage. Graph is primary checkout generation 2026-10-04T22:18:10Z, unwatched; archive scope excluded. Current exact source is authoritative. Graph trace of Snapshot returned no callees/callers despite snippet reporting callers, so no graph completeness claim is made. Direct reads covered all new implementation/tests and relevant existing decoding/history/model sources.

Next: inspect author claim/evidence ledger, verify source/provenance/log hashes and mutation assertions without executing mutation scripts, run exact Go1.25.13 offline test/race/vet/build in a private /tmp cache, and review canonical documentation deltas. No broader scope or heavy pivot currently indicated.
