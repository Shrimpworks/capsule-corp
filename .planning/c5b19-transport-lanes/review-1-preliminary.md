# Preliminary review ledger

Review instance: 1 of 3. Internal read-only G3a12-REVIEW.
Status: IN_PROGRESS — TRENDING_GOOD. Parent Gate3 remains IN_PROGRESS — TRENDING_GOOD; child execution, Gates4–6 and admission BLOCKED.

## Independent first pass

Started from neutral bootstrap only, read the independent-review skill, initialized Serena, and read applicable contributor guides. Read accepted scope/plan, current archive implementation and all appended tests before opening implementation-handoff.md, canonical-verification.md or implementation-evidence.json. Plan contains a short author verification summary, so the first pass was blind to detailed rationale/evidence, not to the presence of reported check results.

Canonical HEAD f1ded569af4a6ae20f2a93b7e476825959f97143 and archive HEAD 7ec433aa17c623f58a81118bab87f1ada49770f3 match bootstrap. Canonical dirty docs and untracked planning are in scope; unrelated .serena/project.yml remains excluded/preserved. Archive has exactly the four named dirty files. Freeze hashes/prefixes are being independently checked alongside this ledger.

## Preliminary concerns / disposition

- No actionable implementation defect identified. validateTransportLanes validates context before empty success, strict-decodes each supplied event, closes kinds3–8, bounds storage and kind/operation counts, rejects duplicate kind+operation, compares unsigned lane-local sequence/tick coordinates pairwise, and retains no mutable state or caller bytes.
- Partial semantics agree with accepted scope. Gaps, equal ticks, cross-lane sequence reuse, isolated results/application, and tuple disagreement across different kinds remain intentionally accepted. This is not a complete-stream or correspondence verdict. Packet complete-stream nonsequential-index refusal remains separate.
- New tests include valid local-pair counterexamples across control/worker, forward/reverse/interleaved order, 18-event deterministic permutations, uint64 boundaries, unsupported/malformed fields, duplicate/cap boundaries, partial-composition positives and nonmutation checks.
- Total/per-kind caps are partly implied by other limits. Their tests can refuse for another legitimate guard; this is acknowledged in scope and is not independent mutation-kill evidence.
- Need verify actual offline compiler/unit checks; exact historical red/green/mutation assertions and restoration; provenance preservation; canonical check limitations; author claims and final unchanged hashes. These are pending evidence checks, not findings.
- Graph belongs to old primary checkout (2026-10-04) and does not index archive/current worktree paths. Coverage queries and focused current-source fallback used; no graph completeness/freshness claim.

## Scope and boundaries

Read canonical starter-document boundary sections, ADR0047, C5b19 reuse row, packet lane/trace requirements and relevant passive model/types. Review is bounded to owned inert synthetic fixture bytes and ordinary offline Go1.25.13 unit/compiler processes; no native fixture, child/guest/lifecycle/signal/wait target, live store/clock/network or mutation execution. Brain owns metadata, immutable pin and publication. No repository edits or review dispatch.
