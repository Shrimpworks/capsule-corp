# C5b19 G3a3 preliminary independent review

Review instance: 1 of 3. Status: IN_PROGRESS — TRENDING_GOOD.

This record precedes reading the author handoff, archive PLAN/ORIGINS, author logs and canonical outcome claims. Bootstrap scope and scope.md acceptance recommendation were available. Read-only review; no source edits or lifecycle fixtures.

Frozen archive HEAD 2274a6c7c751ded537541438f473e81ccd293457, branch codex/c5b19-benign-teardown, four changed paths match bootstrap. Source SHA256s match oracle.go fe36b0504a4a69623b1f64e791487e35ae78869c8e1524ff131d9790d3180dc2 and model_test.go b0f3b9d063c225eb887cecaf16f4058f0ed07ae8f14c3d92c49acacabb0a931e. Canonical HEAD 080e9203c6c96d1494ea3a52e8855816fec3fa38 has the three declared modified docs and untracked planning directory.

No actionable implementation defect identified on the first pass. Exact source and tests were checked against Gate-2 fixed transport/durable/trace requirements and teardownpassive/types.go and model.go, particularly beginWrite, SettleWrite, BeginRunnerIdentityWrite and BeginTerminalWrite. New history enforces initialization, frozen bindings, a copied sole pending request, first/next generation, consumed IDs and kinds, confirmed preparation, terminal closure, preserved confirmed runner prefix and attempted runner identity. It mutates only after validation and exposes only value-owned snapshots. applyReply explicitly represents simulated control application. No lifecycle permission methods, imports or dependency extensions occur.

Test review: positive outcome/disposition matrix covers skipped, confirmed, failed and indeterminate runner publication; separate tests cover failed preparation, pending versus worker-arrival, stale replies and every reply byte, generation/order/freshness, state-preserving rejections, terminal continuity, initialization and copying, varied valid bindings. Helpers build replies independently of new history. A dedicated pending runner → otherwise-valid terminal replacement test avoids accidental rejection by unrelated guards.

Residual checks before final verdict: independently run unit/race/vet/build checks; reconcile red/green/mutant/restoration claims and metadata against actual logs; inspect recovery integration and status boundaries. Initial unit test command started, result pending when this record was written.

Graph: codebase-memory skill applied. list_projects paginated 57 projects; no index for exact archive. Canonical capsule-corp generation 2026-10-04T22:18:10Z reports metadata_match/no_recorded_issue for packet/types/model and scope with no more pages; best-effort only. Focused exact source reads used for archive and canonical worktree. git diff confirms canonical packet/types/model unchanged between primary fc8af2df and target080e920. No graph completeness claim.

Full Gate3 remains IN_PROGRESS — TRENDING_GOOD; Gate4+ BLOCKED; child execution gated.
