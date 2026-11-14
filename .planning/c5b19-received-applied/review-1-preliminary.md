# G3a7 independent preliminary review

Review instance: 1 of 3.
Status: IN_PROGRESS — TRENDING_GOOD. Parent Gate 3: IN_PROGRESS — TRENDING_GOOD; child execution: BLOCKED.

No actionable code finding identified in the initial requirements/tests/source pass. This is preliminary, before author-handoff/outcome reconciliation and independent test execution.

Scope verified: canonical codex/c5b19-received-applied at 041b6d4e04df9667e6217f57df26f04eec7d2ddc; archive codex/c5b19-benign-teardown at 14365fb4ba5d1525285406ae90feb1abdcc46269. Explicit uncommitted changes are three canonical docs plus .planning/c5b19-received-applied and four archive paths. Unrelated .serena/project.yml excluded. No implementation commits or source writes performed.

Independent reconstruction: Gate2 packet lines 290–382 defines common binding/timebase, control kind7 received and kind8 decision records, exact reply tuple and eventual paired confirmation. Accepted scope deliberately selects well-formed correspondence and a bounded subsequence, with state-transition correctness and full trace/causal/timing composition deferred. The new validator delegates exact sidecar validation, requires kinds7/8, compares six tuple fields, orders pair sequence/tick by direct comparison, then bounds collections to three distinct operations and checks all six supplied control events for unique sequences/nondecreasing ticks. It returns errors only and retains no input or state. The implementation matches that bounded contract.

Tests inspected before implementation: seven added tests cover operation shapes, both decisions/all outcomes, tuple mismatch, kind mismatch, contexts/member lengths, sidecar corruption, local boundary ordering, cap/reuse, cross-role sequence collision, cross-pair tick order in both input orders, and nonmutation/prior API independence. Remaining checks: independently execute exact offline test/race/vet/build; verify baseline prefixes/module/model/provenance/hash claims; inspect mutation and canonical-check evidence; assess doc claims and pending final integration gates.

Blinding qualification: bootstrap and accepted plan disclosed author progress summaries (claimed compiler-red, checks and mutants) before the independent pass, so blinding was imperfect. Author handoff, new archive PLAN outcome additions and ORIGINS additions have not yet been read. No author verdict adopted.

Graph: codebase-memory Tier2, primary index ready at 2026-10-04T22:18:10Z, not watched. SettleWrite search one result, outbound trace four callees without continuation, exact snippet confirms pending/outcome validation before projection and representable failed/indeterminate settlement. Exact-path coverage checked for 16 canonical evidence files: normative sources/docs metadata_match; scope/plan/recovery freshness missing and direct reads used for scope/plan (recovery doc reconciliation pending). Archive unindexed: direct source reads. No completeness claim.

Defensive environment: owned inert synthetic in-memory bytes and ordinary offline compiler/unit tests only. No native fixture/child/guest/lifecycle/signal/wait/live-clock/store/network/credential/user-data operation authorized or performed.
