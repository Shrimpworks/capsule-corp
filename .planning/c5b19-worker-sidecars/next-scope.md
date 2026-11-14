# Next bounded Gate3 slice: supplied write/result correspondence

PASSED — internal read-only G3a10-NEXT recommendation. Full Gate3 remains IN_PROGRESS — TRENDING_GOOD; child execution, Gates4–6 and admission remain BLOCKED. No implementation, tests, native fixtures, repository edits or external access occurred. Brain owns integration of this sole retained output.

Recommend exactly one slice: **supplied request-write and worker-reply-write result pairing**, using inert owned archive bytes only.

Gap: G3a6 closes kind3/5/7 frame correspondence; G3a8/G3a9 retain submissions and accepted replies; accepted G3a10 distinguishes independent worker receive/send slots. These contracts do not pair kind3 with kind4 request-write-result or kind5 with kind6 worker-reply-write-result. Packet338–342 separates pre-write events from returned byte-count/errno evidence;361–381 requires closed fields and pairing;189–196 requires one whole-frame write and refuses short/zero/error/excess rather than retrying.

Acceptance boundary for a separately approved consuming scope:

- Pure supplied-pair validation using existing strict helpers: kind3 captured397 request plus kind4 result, or kind5 captured98 reply plus kind6 result. No new serialization, event kind, state ledger or dependency.
- Match operation/generation/request/candidate identity and correct same lane; kind5/6 also match outcome/record digest. Result sequence must follow its paired event, tick must not regress; equal ticks allowed. No cross-lane ordering. This new local pairing claim does not strengthen G3a10 retroactively.
- At most three pairs/category and three distinct operations overall; reject duplicate category/operation. Arbitrary selection order; empty selection is no evidence. Omitted whole pairs remain unknown; explicitly selected incomplete/mismatched pairs refuse. Neither category requires the other.
- Preserve raw result/errno facts. Any full-write observation requires exact397/98 length and the frozen successful-result representation; short/zero/excess/error never becomes transport progress, retry or settlement. Pin errno representation before coding; invent no numeric errno whitelist. A failed represented syscall can remain valid diagnostic evidence without a passing transport verdict.
- Keep G3a10 receive/send/accepted evidence separately evaluated. A write-result cannot manufacture worker/control receipt, kind8 application or durable confirmation. Preserve history/input nonmutation, owned returned facts and atomic zero-result errors.

Tests: strict malformed/wrong kind/lane/fields; substituted tuple/digest; kind5/6 outcome mismatch; missing counterpart; reversed/equal sequence and regressed ticks; valid same-tick pairing; exact versus short/zero/excess/error result; duplicate/cap boundaries; arbitrary selection order and partial categories; cross-lane tick independence; ownership/nonmutation. Retain tests-first red/green, compiled mutants reaching intended assertions, exact restoration, ordinary offline checks and fresh review. No native syscall/pipe/child/guest/clock/lifecycle fixture is needed.

Limits: no authenticity/completeness, whole-lane coverage, global chronology, physical transfer, reader receipt, actual syscall success, durable publication, full round trip, control settlement, start/stop permission, timing verdict, native absence or admission. Full stream/control/timing composition remains deferred.

Priority: this closes a finite named transport-evidence gap immediately adjacent to the current correspondence chain without architecture expansion.

Basis: canonical packet182–211,300–381,395–410,429–451 at fe9009d5398a883735219e2a3a33c351c79aa54c; archive oracle strict event/capture and private snapshotHistory contracts at7646d8faabb670ce8b5ff1b213c414d91deda2d9; accepted /tmp/c5b19-g3a10-scope.md. Reused prior qualified Tier2 graph evidence (old2026-10-04 primary not watched; archive outside_project) and bounded exact source reads; no broad audit/completeness claim. Confidence high for this gap; errno representation and final API require freezing in the consuming scope.
