# G3a4 independent preliminary ledger

Review instance: 1 of 3.
Status: IN_PROGRESS — TRENDING_GOOD. Full Gate 3 remains IN_PROGRESS — TRENDING_GOOD; Gates 4–6, child execution and product admission BLOCKED.

This ledger was saved before reading /tmp/c5b19-g3a4-implementation.md, archive PLAN/ORIGINS outcome additions, or canonical verification/outcome records. The neutral scope/plan necessarily contained brief progress summaries; those were not used as proof of checks.

## Scope reconstructed

Canonical /Users/dsteele/.codex/worktrees/project-recovery/capsule-corp is codex/c5b19-event-conformance at 2aefb15ee45018bb1045e711c77a0c7bc9f4c54d. Archive /tmp/capsule-c5b19-recovery is codex/c5b19-benign-teardown at 6575898a32ff4a3d24ca4c757f42b1807798a372. Status agrees with bootstrap; unrelated .serena/project.yml excluded. Archive has exactly four modified paths. Both oracle hashes match frozen bootstrap.

Read canonical Gate2 raw event contract, scoped acceptance and delivery plan, archive tests before implementation, reused binding encoder, and canonical five-file Go patch/doc plus roadmap diff. New 103-line implementation is pure bounded envelope/lane parsing. Each field offset matches the packet; error returns are zero-valued; fixed arrays own values; length checks precede slices; binding digest uses exactly preimage[30:270]; lane count cap precedes iteration; checked sequence arithmetic excludes wrapping; only control ticks have a local watermark. No semantic event/trace claim or writeHistory linkage introduced.

## Preliminary findings

No actionable defect identified from the independent source/test/diff inspection. Per-kind fields/lane ownership, causal/timing composition, capture completeness and origin choice are expressly excluded, not accidental omissions. Relaxed worker/harness tick ordering is explicitly accepted by the scoped plan and remains distinct from the future full timing oracle.

## Questions to close with evidence

1. Confirm test-first compiler-red, selected mutation sensitivity and exact restoration from retained logs and hashes.
2. Confirm canonical Go1.26.11 checks and unchanged scanner enforcement, with pre-existing lint debt kept separate.
3. Reconcile archive PLAN/ORIGINS provenance and final publication obligations; pending commits/pins/hosted acceptance are planned later work, not this source verdict.
4. Finish independent offline test/race/vet/build/diff checks; initial exact Go1.25.13 test is running.

## Coverage limitation

Codebase-memory list_projects identifies capsule-corp at the main checkout, not this review worktree. search_graph for new/reused archive symbols returned no matches; archive is unindexed, so exact source/diff reads supply evidence. Canonical coverage generation2026-10-04T22:18:10Z returned best-effort metadata_match for most paths but missing for current recovery/planning files; these were read directly. The first requested ADR filename was incorrect and resolved from actual directory listing; use 0047-prepare-teardown-obligation-before-launch.md. No graph absence proves anything about the archive/current diff.
