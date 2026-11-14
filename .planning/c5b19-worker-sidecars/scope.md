# G3a10 supplied worker-received/sent correspondence — scope recommendation

PASSED — internal read-only G3a10-SCOPE derivation. Parent Gate3 remains IN_PROGRESS — TRENDING_GOOD; child execution, Gates4–6 and admission remain BLOCKED. Brain accepted the three-row contract and detailed boundaries before implementation; this report records that frozen scope. This report is the sole retained output of this assignment; no repository edits, tests, mutation execution, native fixtures, child/guest execution, lifecycle effects or external targets occurred.

Question: can an offline, nonmutating selection query distinguish actual supplied raw worker receive/send slots from decoding context, and compare those slots with privately retained submissions, strict kind5 evidence, optional kind7 evidence and any actually accepted reply? Defensive intent is to reject substituted/malformed evidence without manufacturing process authority or settlement. Authorized assets are inert synthetic byte fixtures in the owned archive; no other systems, identities, credentials or data.

## Source and retrieval basis

Canonical base fe9009d5398a883735219e2a3a33c351c79aa54c, branch codex/c5b19-worker-sidecars, /Users/dsteele/.codex/worktrees/project-recovery/capsule-corp. Archive base verified 7646d8faabb670ce8b5ff1b213c414d91deda2d9, branch codex/c5b19-benign-teardown, /tmp/capsule-c5b19-recovery.

- Current .planning/c5b19-worker-sidecars/plan.md and prior .planning/c5b19-snapshot-correspondence/next-scope.md define a bounded correspondence query, not full trace completion.
- docs/C5B19_GATE2_IMPLEMENTATION_PACKET.md:338–359 distinguishes control request/reply slots from separately retained worker received/sent slots; worker reply alone never licenses start/terminal. Lines395–410 make each lane independent and worker clocks observational. Lines429–451 require exact identity and causal control events in the eventual full oracle, forbid total cross-lane timestamp ordering, and keep first-run authorization separate.
- Archive experiments/typed-guest-transport-c5b19-benign-teardown/tests/oracle.go:622–626 explicitly says kind5/7 capturedFramePair.Request is decoding context, not worker receipt. validateCapturedFrame:632–660 closes event kind/fields, request/reply framing and hashes. validateCapturedFrames:664 onward explicitly promises only correspondence, not lane/causal/timing order or omitted-pair detection.
- snapshotHistory:875–1013 already privately owns up to three submitted request arrays and accepted reply arrays/outcomes; begin/reconcile record them transactionally. inspect provides the precedent for nonmutating historical queries, exact accepted-reply comparison and copied value results. No parallel ledger is needed.
- internal/execution/teardownpassive/model.go:185–211 validates outcome and exact pending before settlement; only confirmed outcome changes durable projection. Parent supplied complete SettleWrite outbound trace with four callees; exact source was reread. This query must call no settlement method.
- Applicable canonical/archive AGENTS and bounded starter-document reads preserve nonproduction claims, exclusive Supervisor process authority and explicit gates. ADR0047 (0047-prepare-teardown-obligation-before-launch.md):26–74 preserves unknown-outcome uncertainty and one pending operation without new lifecycle authority. Ecosystem reuse row135 is test-only ADOPT-PLATFORM / BUILD-NARROWLY; no new dependency or primitive is needed.
- Tier2 best-effort graph: capsule-corp ready, generation2026-10-04T22:18:10Z, old primary checkout, not watched. index_status and coverage checks completed; canonical packet/model/ADR/reuse/starter paths reported metadata_match, planning files missing, archive paths outside_project. These are qualified metadata signals, not freshness/completeness proofs. Exact current worktree reads cover relied-on packet, model, ADR, reuse and planning ranges; exact archive reads cover the helpers and private ledger. One mistaken canonical-relative archive test path was absent; corrected archive lookup located existing snapshot tests2169–2630. No test execution occurred. No exhaustive repository claim is made.

## Proposed API and exact selection contract

Append a pure query method to existing snapshotHistory. Names are implementation choices. Each input item contains four distinct members:

1. WorkerReceived []byte: independent raw worker-received slot, required and exactly397 bytes.
2. WorkerSent []byte: independent raw worker-sent slot, absent only when nil, otherwise exactly98 bytes.
3. WorkerReply *capturedFramePair: strict kind5 event plus its decoding context and reply; pointer nil means absent.
4. ControlReceived *capturedFramePair: optional strict kind7 event plus its decoding context and reply; pointer nil means absent.

These are in-memory selections only, not new serialized formats or event kinds. Required WorkerReceived is never inferred from kind3, kind5.Request, kind7.Request, successful writes or accepted history. Nil WorkerSent means absence; nonnil zero-length WorkerSent is present but malformed and must refuse; a nonnil pair pointer is present and must validate, even if its slices are empty. There is no zero-sidecar anchor item.

| Received | Sent | Kind5 | Kind7 | Disposition |
|---|---|---|---|---|
| present | absent | absent | absent | Valid partial received-only evidence |
| present | present | present | absent | Valid partial sent-without-control evidence |
| present | present | present | present | Valid supplied correspondence only |
| absent | any | any | any | Refuse selected item |
| present | present | absent | any | Refuse incomplete sent/event pair |
| present | absent | present | any | Refuse incomplete sent/event pair |
| present | absent | absent | present | Refuse unsupported control-only row |

All valid rows require:

- Strict decode of WorkerReceived using privately retained frozen bindings; exact byte equality against one actual snapshotHistory entry. A well-formed but unsubmitted request is refused. Do not accept a caller-provided history snapshot or request context as membership.
- For sent rows, strict validateCapturedFrame(WorkerReply), explicit Kind==5, and exact WorkerReply.Request==WorkerReceived plus WorkerReply.Reply==WorkerSent. Independently supplied sent bytes cannot be silently replaced with pair.Reply.
- If ControlReceived exists, strict validateCapturedFrame, explicit Kind==7, exact ControlReceived.Request==WorkerReceived and ControlReceived.Reply==WorkerSent. Its kind7 bytes establish only represented receipt; do not invent kind8/application or consult the current history watermark.
- If the retained entry outcome is nonzero, a supplied WorkerSent must equal the entire retained accepted reply array. Here accepted means ApplicationResult1; its outcome may be confirmed1, failed2 or indeterminate3. Do not compare only confirmed digests or only outcome. A historical rejected ApplicationResult2 does not replace this accepted evidence. Receive-only rows do not infer a worker send from accepted evidence.
- If no accepted reply exists, valid worker outcomes1/2/3 remain observable but do not settle pending, clear fatal, or update any projection. A supplied kind7 never implies application.

## Collection, chronology and ownership

- Nil or uninitialized snapshotHistory refuses every query, including empty selection. An initialized history accepts empty selection with zero selected observations: no evidence or completeness claim. Initialized zero-submission history still refuses every nonempty selection because membership cannot hold.
- Maximum three selected items, irrespective of row shape. At most one item per operation/submission; duplicate identical rows, duplicate rows with different optionality and same-operation substitutions all refuse. Input order is arbitrary and result order follows selection order.
- This slice imposes no new same-lane sequence ordering, sequence uniqueness, gap coverage, causal ordering or cross-lane clock relation, following G3a6. Strict per-event field/lane/timebase validation remains mandatory. Equal ticks, sequence gaps, arbitrary selection order, old event sequences and worker ticks on either side of control ticks are not chronology errors here. This is deliberate partial selection acceptance, not evidence that a complete trace could pass. No current control watermark comparison is allowed in this historical query. Full stream consistency remains deferred.
- All successful outputs contain copied value facts only: selected wireOperation, explicit worker-sent/control-received presence, worker outcome (zero if absent), retained accepted outcome (zero if absent), and a copied decision-history snapshot as useful. A received-presence field is optional because every returned item has received evidence. No confirmed/complete-roundtrip/final-verdict boolean. No returned slices or pointers may alias caller bytes or private arrays; a newly allocated result slice of value-only items is acceptable.
- A refusal returns an error and no partial result (nil/zero result collection), including when a later selected item fails after an earlier valid one. Neither success nor error mutates the receiver, supplied bytes, accepted entries, pending state, fatal latch, control watermark or projection. Mutation of caller inputs or returned values after a call cannot change history or previously returned facts.

## Tests-first acceptance matrix

1. Context: nil/zero receiver with empty and nonempty selection; initialized empty history empty success/nonempty refusal; nil and allocated-empty selections. Verify no evidence on empty success.
2. Every valid row, all worker outcomes1/2/3, pending and settled histories, and historical entries after later submissions. Received-only accepted history reports accepted outcome while sent presence/outcome remains absent/zero. Sent+kind5 without control remains partial; kind7 without kind8 never settles.
3. Enumerate all unsupported optional combinations, including nil versus empty sent, nonnil-empty capture pointers, absent required receive, malformed sent in an otherwise present row, and forbidden control-only row.
4. Boundaries: exact397/98/256 frames; zero/truncated/overlong requests/replies/events; invalid headers, frozen bindings, request/candidate digests, outcome/recorddigest relation, unused fields, role/lane/timebase, wrong event kinds (including kind7 in worker slot and kind5 in control slot). Reuse strict validator helpers; do not weaken prior validators.
5. Exact membership: independently well-formed unsubmitted receive; request substitution with same or different op/gen; a decoding-context-only request is insufficient. Anchor substitution tests must remain otherwise valid enough to reach intended equality.
6. Independent sent comparison: WorkerSent differs from a separately valid WorkerReply.Reply for the same request; use valid failed2 versus indeterminate3 replies with zero recorddigests to isolate equality. Test worker context differs from independent receive. Test kind7 request and reply independently differ while each pair itself validates.
7. Accepted linkage: accepted outcomes1/2/3 require exact retained reply equality for supplied worker send; valid contradictory2/3 isolates accepted-equality guard. After valid rejected ApplicationResult2, a matching prior accepted reply stays queryable; rejection does not overwrite it. Historical entries still resolve after current pending advances.
8. Collection: three distinct retained operations in arbitrary order, mixed valid rows, four inputs refusal, repeated identical operation, same operation with changed optionality, and late invalid item yields zero results. Duplicate and cap checks count selections, not only sent/control-present rows.
9. Explicit lack of chronology: consistent byte selections with scrambled input order, old sequences relative to current watermark, same ticks, gaps, independent worker/control tick order remain correspondence-valid; no inferred lane ordering is asserted. Strict event lane validation still refuses incorrect lanes.
10. Nonmutation/ownership: compare whole private history before/after every success/refusal; pending remains pending after confirmed worker observation, fatal remains latched, accepted evidence remains unchanged. Compare input bytes before/after. Change caller byte arrays and copied returned facts after success, then requery retained history with independent original fixture bytes. Verify no aliases and no partially returned evidence on failure.

## Targeted mutant proposals and limitations

Mutants are future controlled edits for implementation verification, not executed here. Select compilable behavior mutations that reach intended named assertions after valid setup: remove exact retained membership; ignore independent received versus context; ignore independent sent versus kind5 reply; allow wrong kind5/kind7 role; bypass optional dependency guard; ignore control request/reply correspondence; skip retained accepted-reply equality; permit duplicate operations; leak partial collection on later error; return aliases or mutate history during query if such implementation paths exist. Use otherwise-valid fixtures and distinguish guard reachability from malformed-parser failure.

Some checks compose redundantly: validated kind5 identity, exact request matching and membership can mask one another; the cap is also structurally implied by three private entries plus duplicate refusal. A cap-only mutant may survive because all fourth items already fail membership/duplicate checks. Do not claim independent mutation coverage for unreachable/redundant guards, and do not weaken other safeguards to force a kill. Record the intended guard, concrete assertion, observed failure and restoration digest. Compilation failure, fixture setup failure or an unrelated earlier assertion is not a useful mutation kill. Test helper changes must remain within new append-only tests; prior source/test prefixes and provenance stay exact.

## Exclusions and disposition

No capture authenticity, omitted-frame detection, capture completeness, successful syscall publication, physical transport completion, durable publication, full round trip, full lane ordering/coverage, full control-trace causality, elapsed-time bound, stop/start permission, native absence, final verdict, execution authority or product admission follows from this query. A worker-received event kind does not exist and must not be invented. Entirely omitted selections remain unknown. Existing G3a9 absent/older/pending/failed/indeterminate snapshot semantics remain unchanged.

Confidence: high for the bounded contract grounded in existing helpers and private retained entries; eventual full-stream composition and native behavior are expressly unresolved. Recorded decision: Brain accepted this exact narrow contract and boundary details; freeze it and proceed with tests-first implementation in the existing shared archive worktree with ordinary offline unit/compiler checks and fresh independent review. No source change is needed merely to finalize this scope report.
