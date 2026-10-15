# Preliminary independent review

Review instance: 1 of 3 for the distinct G3a2 wire/record delta and recovery reconciliation. Status: IN_PROGRESS — TRENDING_GOOD.

Written before reading author handoff, archive PLAN/ORIGINS, or recovery reconciliation claims. Neutral bootstrap includes prior-flow history and stated scope; no author rationale has been read. Archive source now committed unchanged at c4d3c33ed04c1dff213e2c69064a6113892d9653 against 1e3f3bb154d0ef24412002ac47adef22ba3c46a9.

Evidence inspected: Gate-2 packet fixed transport/durable schema and source/build closure; teardownpassive types.go bindings and projections, model.go preparation/runner/terminal construction and generation/settlement rules; complete oracle.go and current oracle_test.go; deleted obligation_test.go diff.

Preliminary findings: no actionable implementation defect identified. Fixed offsets independently reconcile to request397/reply98/record395. Strict outer lengths guard slicing. Frozen bindings/policy/obligation checks precede copied operation acceptance. Candidate kinds preserve preparation, runner, terminal with optional unconfirmed runner and matching nonzero absence. Failed/indeterminate replies demand zero durable digest and exact body hash. This is structural correspondence, not stateful confirmation. The preparation tests removed from obligation_test.go were moved to oracle_test.go.

Tests include literal preparation vectors/hashes, all byte substitutions in correspondence, every truncation and trailing byte, varied valid bindings, zero/duplicate matching context, generation bounds, candidate faults with freshly matched hashes, all three outcomes and terminal variants. Literal known answers are not a decoder/encoder roundtrip. The non-preparation helpers restate offsets; a future full model conformance test remains outside this increment.

Pending checks: execute offline tests/vet; reconcile archive manifest hashes and plan limits; inspect recovery diff and retained validation claims; independently check literal vector calculations. No lifecycle, storage, signal/wait, guest, credential, or external-target behavior will be exercised.

Coverage limitation: codebase-memory skill read; list_projects pages 0 and 50 show neither exact review checkout indexed. Focused source reads used for exact target instead of a graph for another checkout. No completeness claim from graph metadata.
