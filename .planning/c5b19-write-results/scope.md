# G3a11 supplied write/result correspondence — proposed frozen contract

PASSED — internal read-only G3a11-SCOPE derivation; awaiting Brain acceptance before implementation. Parent Gate3 IN_PROGRESS — TRENDING_GOOD; child execution, Gates4–6 and admission BLOCKED. Defensive scope: owned inert archive byte fixtures only. No tests, native/child/guest/lifecycle/live-store/network execution or repository edits occurred. Sole output: this /tmp report; Brain owns integration into canonical planning and existing archive draft PR38. No branch/commit/push/PR created.

Canonical worktree branch codex/c5b19-write-results, base d042e53a877e494bd2eb0e325ec084d4b9cc692a. Archive base d698bd9332a000ad3c2b9f23f87d0353dc5964b4 in /tmp/capsule-c5b19-recovery. Preserve existing APIs, source/test prefixes and provenance.

## Critical semantic resolution

Packet187–214 requires one whole397-byte request or98-byte reply write; short/zero/error/excess is not transport progress or a retry. Packet300–321 stores both result and errno as signed i32_be. The current strict event decoder explicitly does not interpret syscall semantics; archive PLAN543–544 says signed result/errno observations are preserved without guessed success rules. No cited source specifies resetting errno before a successful call or zero-normalizing the captured errno afterward.

Therefore freeze these three distinct concepts:

1. Diagnostic validity: any signed32 result and signed32 errno is permitted by this new layer if strict event shape, identity and pair chronology validate. Negative, zero, excess and seemingly contradictory values remain uninterpreted diagnostic observations; acceptance never asserts they could arise from a conforming syscall.
2. ExactFrameCount: true iff result equals397 for kind3→4 or98 for kind5→6, independent of errno. Preserve both raw values verbatim. This is only a represented count equality, not successful syscall evidence.
3. Full successful write/transport progress: no such verdict is exposed in this slice. In particular do not introduce FullWrite, Success or TransportPassed booleans, infer success from errno==0, or refuse a diagnostic solely because errno is signed/unrecognized. An exact count with nonzero errno remains exact-count diagnostic evidence. Full native syscall capture semantics require a later explicitly normalized producer contract/evidence.

This closes the unresolved errno issue without inventing packet semantics, a numeric whitelist or platform calls. If Brain wants a differently named conservative derived conjunction, that requires explicit scope amendment before coding; ExactFrameCount plus raw Errno is the recommended minimal contract.

## Pure pair API

Append an in-memory selection type containing a capturedFramePair (pre-write event plus raw decoding frames) and a required result-event byte slice. Prefer a stateless pure function taking selections, frozen obligationBindings and numer/denom, composing existing validateCapturedFrame/decodeEventFields/decodeRequest. No new private ledger, history mutation or serialized format is needed.

Valid pair kinds are exclusively3→4 and5→6. A selected item with either member nil, empty, malformed or wrong-kind refuses. The kind3 capture requires its exact397 request and empty Reply; kind5 capture requires its exact397 request context and exact98 reply. Kind5 Request remains decoding context, not independent worker receipt; only G3a10 makes the latter distinction.

Validate both events strictly under identical frozen bindings/timebase. Result operationID, generation, request digest and candidate digest must equal the pre-write event. For kind5→6, outcome and record digest must also equal. Pair lanes must be control1 for3/4 and worker2 for5/6 (already enforced by strict helpers, retain explicit closed pair-kind check).

Result.Sequence > Pre.Sequence and Result.Tick >= Pre.Tick, using direct comparisons without subtraction/overflow. Same tick is valid; sequence gaps are valid. These are local pair observations only. Do not enforce adjacency, compare worker/control clocks, consult current history watermark, or claim chronology relative to omitted/intervening events.

## Collection, context and ownership

- Validate frozen context and nonzero numer/denom even for empty selections. Valid context plus nil/empty selection returns zero Count/no evidence. Invalid context always refuses. No receiver initialization issue exists for the recommended stateless API.
- Maximum6 pairs total, at most3 per pre-write kind, and at most3 distinct operationIDs overall. Duplicate same-category/same-operation refuses even if byte-identical or with a different result. Cross-category same-operation is allowed once/category only when the entire raw397 request bytes agree; equality covers full binding, operation, generation, kind, candidate and clock policy, not merely hash/ID. Different categories for different operations need not pair up. Entire omitted pairs remain unknown.
- Input order arbitrary; return facts in input order. Cross-category byte consistency does not prove control send caused worker receipt. This slice does not enforce same-lane inter-pair ordering or sequence uniqueness beyond each selected pair; those are full-stream consistency checks. Document this limit explicitly rather than implicitly certifying all supplied events as one valid complete lane.
- Return fixed bounded value-owned facts (for example [6] entries plus Count), containing category/operation, raw signed result/errno, ExactFrameCount, and copied sequence/tick facts if useful. No slices/pointers may alias input or private state. Return the completely zero result on every error, including an invalid later pair. Do not return a validated prefix.
- No mutation of caller bytes, frozen bindings, any existing history or durable projection. Subsequent caller-input/result mutations cannot alter prior returned facts or other calls. No begin/reconcile/SettleWrite invocation.

## Tests and targeted mutation evidence

Tests-first matrix:

1. Valid request and all three reply outcomes, each isolated and mixed; same-tick pair, arbitrary input order, sequence gaps and independent cross-lane tick order. Prove neither category requires the other.
2. Empty/nil selection with valid context; invalid frozen bindings and zero timebase even when empty. Selected missing/empty members refuse rather than becoming omitted pairs.
3. Strict malformed/truncated/overlong events and frames; unsupported pre/result kinds, swapped lane, unknown/forbidden fields, binding/timebase mismatch. Reuse strict decoders without weakening prior contracts.
4. Independently valid result tuple substitutions (ID/generation/request/candidate), worker outcome/recorddigest mismatch and valid same-operation cross-category request substitution. Build complete valid setups to reach the new equality checks; malformed earlier frames are not substitute evidence.
5. Sequence equal/reversed and tick regressed refuse; equal ticks and nonadjacent increasing sequences pass; boundaries zero/max u64 use direct comparisons. No newly implied inter-pair chronology.
6. Raw result/errno Cartesian boundary samples: MinInt32, -2, -1,0,1, frameLength-1, exact frameLength, frameLength+1, MaxInt32 for result; MinInt32,-1,0,1,MaxInt32 for errno. Preserve every signed value on otherwise valid pairs. ExactFrameCount depends only on correct category length; exact count with nonzero/negative errno remains true, all other counts false. A98 request count is false and397 reply count false.
7. Six pairs/three operations with both categories; fourth unique operation refusal; fourth pair in one category; seven total; duplicate same-category operation, including changed result; same operation across categories with identical request succeeds and changed request refuses. Late invalid pair yields completely zero result.
8. Input preservation and no aliases: mutate input arrays after success, verify copied facts stable; mutate returned facts then revalidate original independent fixture; all errors leave bytes unchanged. Existing history APIs are not involved and remain unaffected.

Selected compiled mutants should remove pair-kind closure, one unmasked tuple/outcome comparison, local sequence or tick comparison, cross-category full-request equality, duplicate/unique-operation enforcement, category-specific exact-count calculation, or atomic zero-result behavior. Include an errno-dependent ExactFrameCount mutation to establish that no ungrounded zero-errno invariant slipped in. Retain valid setup and intended assertion evidence, restoration digest, red/green and standard offline checks followed by fresh review. Some cap checks are logically redundant with duplicate/unique/category checks; report survivors honestly, never weaken safeguards to force independent kills. Compilation/setup failures are not useful kills.

## Evidence, confidence and limits

Sources: canonical .planning/c5b19-write-results/plan.md and worker-sidecars/next-scope.md; packet187–214,300–381,395–410,429–451; archive tests/oracle.go decodeEventFields510–617, captured-frame validators622 onward and current inspectWorkerSidecars; archive PLAN543–544. Prior ADR0047/authority and reuse-row boundaries remain unchanged. The installed codebase-memory skill and inherited Tier2 trace context were reused; index_status confirms old2026-10-04 generation at unwatched primary checkout. Coverage: packet metadata_match, both planning files missing freshness, archive outside_project. Exact current-worktree reads cover relied-on planning/packet/helper ranges; no graph freshness, absence or exhaustive claim.

Confidence high for source-defined byte pairing and signed diagnostic preservation; no source establishes normalized errno semantics, so successful-syscall classification is deliberately excluded. No capture authenticity/completeness, real syscall outcome, whole-lane chronology, physical transport, independent worker receipt, control application, durability, full roundtrip, timing, stop/start, process absence, execution authority or product admission is established. First-run authorization remains separate and BLOCKED.
