# G3a8 preliminary review

Review instance: 2 of 3. Status: IN_PROGRESS — TRENDING_GOOD. Internal read-only review; parent Gate 3 incomplete and child execution BLOCKED.

Scope verified: canonical codex/c5b19-decision-reconciliation at 8e4fd22cc00cd5e961544f6a803974600dc0e751, archive codex/c5b19-benign-teardown at 5d0e53cbccae3d4f587205df70253b776eae63e5. Dirty paths match bootstrap; unrelated .serena/project.yml preserved. No source mutations.

Independent source/test/plan inspection precedes reading author handoffs, prior review report, and mutation scripts/logs. Bootstrap and current plan disclose the previous evidence defect, so this is not blind to its existence; detailed author rationale has not yet been consulted.

## Preliminary ledger

- No actionable implementation defect identified. decisionHistory owns fixed value state, validates before trial mutation, derives acceptance from exact pending bytes plus private settlement, rejects false recorded outcomes atomically, preserves history on represented rejection, latches fatal monotonically, and bounds decisions and cross-call control chronology.
- Tests independently assert projections, exact whole-state preservation, valid wrong-anchor fixtures, late settlement outcomes, duplicate/stale rejection, borrowed-buffer isolation, bounds and uint64 edges. History-only begin after fatal is deliberately not execution admission.
- Passive adapter change distinguishes cleanup context cancellation from deadline expiration. Since only worker cleanup cancels the background-derived method context, and the buffered send precedes cleanup, receiving the result on context.Canceled is safe. Synthetic tests force both select cases ready and preserve caller cancellation/slot semantics.
- Plan scope is appropriately bounded: serial action groups, no full control-state/admission, no malformed diagnostic rejection, no missing-event detection, no native fixture or child. Required future publication/pinning is Brain-owned and remains open.
- Evidence question pending: verify corrected mutant retains exact preparation settlement and reaches the wrong-anchor assertion, rather than failing prerequisite runner begin; verify restored hashes and current provenance. Do not execute mutation tooling.
- Verify proportionate local tests/race; inspect rather than rerun unchanged broad canonical gates.

## Retrieval limitations

Serena initialized. codebase-memory list_projects paginated completely; graph search for adapter symbols paginated 100+58, runPassiveAdmitted traced both directions and snippet inspected. Exact six source/test/model review paths report outside_project, unavailable coverage: graph targets main checkout, not these frozen worktrees. Current exact source reads supply fallback; no completeness claim.
