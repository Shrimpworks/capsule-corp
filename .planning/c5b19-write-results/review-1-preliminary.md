# G3a11 preliminary independent review

Review instance: 1 of 3.
Status: IN_PROGRESS — TRENDING_GOOD.

Recorded before opening implementation-handoff.md, canonical-verification.md, detailed implementation evidence, or appended archive PLAN/ORIGINS rationale. The neutral bootstrap, scope.md and accepted plan.md were read; plan.md itself includes summary author verification claims, which remain unverified testimony.

## Scope reconstruction
Canonical branch codex/c5b19-write-results at d042e53a877e494bd2eb0e325ec084d4b9cc692a; archive branch codex/c5b19-benign-teardown at d698bd9332a000ad3c2b9f23f87d0353dc5964b4. Dirty paths match bootstrap; canonical .serena/project.yml is unrelated and excluded. Four archive paths are modified; source and tests append only. No edits or execution of native/child/guest/lifecycle targets.

## Preliminary findings ledger
- No actionable implementation defect identified in the 89-line appended API or nine test groups.
- Verify prior strict decoder guarantees: closed kind/lane, exact request/reply/context, signed diagnostics and owned wireOperation fields justify ignored repeat decode errors and bounded indexing. Source inspection supports these guarantees.
- Collection limits and duplicate/cross-category exact-request checks appear correct; no inter-pair chronology or receipt/transport success is claimed.
- Tests independently construct raw event slots, validate substitutions before refusal, cover both category counts and signed boundaries, local ordering, caps, partial selection, atomic refusal, ownership. Some cap checks overlap logically; inspect mutation evidence for intended assertion reachability rather than relying on exit status.
- ExactFrameCount independent of errno follows accepted scope and packet omission of errno normalization. It is deliberately no syscall-success verdict.
- Canonical docs accurately keep Gate3 incomplete and child/later gates blocked; archive immutable pin and acceptance/publication metadata remain Brain delivery work.

## Remaining verification
Run exact Go1.25.13 offline test/race/vet/build with private cache; check frozen hashes, prior source/test prefixes/provenance and diff hygiene. Inspect historical red/green/mutation/canonical logs and distinguish those from reviewer execution. Read author rationale only after this ledger.

## Tool limitation
Serena initialized. codebase-memory index is the old 2026-10-04 primary checkout, not current canonical/archive roots. Coverage reports archive paths outside_project; current source reads are authoritative. Graph SettleWrite outbound lookup/snippet supplies structural orientation only; current normative source confirms no relationship from the new API to settlement.
