# G3a7 received/applied correspondence scope

Status: PASSED for read-only scope derivation. Implementation/review remain pending.
Parent Gate 3: IN_PROGRESS — TRENDING_GOOD. Child execution and admission: BLOCKED.
Internal research assignment; no repository changes or publication.

## Question and evidence

Defensively validate correspondence between owned inert kind7 received-event bytes,
exact request/reply sidecars and a kind8 application-decision event. Scope is the
owned local archive and canonical checkout; no native fixture, child, guest,
lifecycle target, clock, live store/network, credentials or user data.

Canonical base041b6d4e04df9667e6217f57df26f04eec7d2ddc; archive
base14365fb4ba5d1525285406ae90feb1abdcc46269. Gate2 packet lines342–360 defines
kind7 as the complete frame received and kind8 as the recorded control validation
decision; lines367–382 bind exact operation fields and require a received/applied
pair for claimed confirmation. Lines290–327 require one clock, nondecreasing ticks
and lane-local sequence. Existing G3a5 supports well-formed represented kind8 for
both decisions; G3a6 validates exact kind7 sidecars. No new serialized container,
source manifest, dependency, native path or toolchain change is justified.

## Accepted minimum recommendation

Add an error-only pure pair validator, e.g. an in-memory value containing a
`capturedFramePair` for Received and raw Applied bytes. Keep earlier APIs unchanged.
No reference is retained, no input mutated, no history or durable state is touched.

1. Require the received event to be exactly kind7. Run existing
   validateCapturedFrame on its event/request/reply to validate exact frozen
   bindings, timebase, field/lane shape, request397/reply98 bytes and recomputed
   request-body/candidate/confirmed-record hashes. Kind3 or5 is not a received
   control event even if its own sidecar validates.
2. Decode Applied with decodeEventFields under the same frozen context and require
   kind8 (therefore control lane). ApplicationResult is1 or2, as already enforced.
3. Require equality of Operation, Generation, RequestDigest, CandidateDigest,
   Outcome and RecordDigest between received and applied. Common binding/timebase
   equality follows from both decoders using the same frozen context. All other
   forbidden fields retain the lower decoder's zero rules.
4. Require received.Sequence < applied.Sequence and received.Tick <= applied.Tick.
   Same-tick pairs are valid when sequence establishes receive-before-decision.
   Do not require consecutive sequences: other control events may intervene.
   Compare directly, without sequence+1 arithmetic or overflow. Zero origins/ticks
   remain representable; this API does not prove capture origin or missing events.
5. Return only correspondence success/failure. A matching decision byte does not
   establish that code1's state transition was valid or that code2 latched fatal.

Both decision1/2 and all reply outcomes1/2/3 are representable. An accepted failed
or indeterminate reply can settle a pending operation without durable projection;
SettleWrite confirms this distinction. A rejected but well-formed frame can be
represented here, but the history/fault reason is not validated. Arbitrary malformed
or truncated rejected-reply diagnostics remain unsupported because the packet has
not frozen their event mapping. Do not call these unsupported encodings universally
invalid under the whole packet.

## Small bounded collection

A list of pair structs is the smallest sufficient wrapper: 0..3 pairs for one
attempt, no duplicate Operation ID (including changed generation), and no missing
received/applied/request/reply member. Empty validates bindings and nonzero timebase
but means no represented correspondence. Omitted entire pairs remain undetectable.
A received-only capture stays supported by earlier G3a6; this API intentionally
validates only selections claiming a corresponding decision.

Across all represented control events (at most6), require unique lane sequence
values and nondecreasing Tick in sequence order. Input list order need not equal
trace order; pairwise comparisons over six fixed slots suffice. This is consistency
of the supplied control subsequence, not completeness, adjacency, operation ordering,
state-machine causality or deadline analysis. No sorting/mutation of caller bytes.
The <=3 pair limit already bounds distinct request anchors; no new frame-table
format or per-call history is necessary.

## Test and mutation acceptance

Use literal-offset fixtures, preserving old vectors/APIs. Positive tests cover all
operation shapes1/2/3, every reply outcome with both application codes, same and
later ticks, sequence gaps, valid shuffled pair-list order, max3 pairs, empty valid
context, and near-uint64 sequence/tick boundaries without overflow.

Negative tests cover each required equality independently; received wrong kind3/5
and applied wrong kind7/other kinds; wrong lanes, forbidden fields, bindings and
timebase; application0/unknown; nil/short/extra event/frame members; malformed or
substituted request/reply; missing one member; duplicate operation with same or
changed generation; fourth pair; equal/reversed pair sequence; regressed pair tick;
cross-pair duplicate sequence and tick regression even when each pair is locally
ordered. Test whole-pair omission only as documented no-evidence, not detectable
failure. Assert input nonmutation and preserve old field/sidecar acceptance APIs.

Meaningful compiled guard mutants should remove tuple/digest/outcome/record equality,
received/applied kind restriction, pair sequence/tick checks, duplicate/cap checks,
and cross-pair sequence/tick consistency. Record targeted failing assertions and
restore exact source bytes. Exact offline Go1.25.13 test/race/vet/build, canonical
required checks and fresh independent review remain delivery gates.

## Limits and next boundary

This establishes matching supplied receive/decision records and local control
subsequence consistency only. It does not validate whether ApplicationResult is
correct, perform writeHistory.applyReply, prove a C5b18 settle transition, fatal
latch, durable snapshot, expected projection from process observations, full
capture completeness, lane origin, request delivery, worker receipt, stop/apply
causality, deadlines, native syscall facts, storage durability or lifecycle authority.
Reversed same-tick stop/apply order cannot be checked without the stop event and
control model; it remains a later composition obligation, not silently passed.
A successful code1 pair is not a claimed confirmation verdict. Next bounded work
should connect represented accepted/rejected decisions to pending-write/control
state under an explicit scope before full causal/timing and snapshot composition.

## Method, confidence and handoff

Read current AGENTS in both repositories, relevant starter-document boundaries,
ADR47 and C5b19 BUILD-NARROWLY reuse row, packet, prior scope/plan, current
G3a5/G3a6 source and normative types/SettleWrite. Tier2 index ready at
2026-10-04T22:18:10Z on primary checkout, not watched. Exact coverage checked for
all cited canonical files: normative/docs metadata_match; prior planning freshness
missing, read directly. SettleWrite search returned one result; outbound trace had
four callees and no continuation; snippet confirmed exact pending/outcome check
before projection. Archive is unindexed and current source was read directly.
No graph completeness claim. No tests ran because this assignment is read-only.

Confidence is high for the exact local correspondence and selected sequence/tick
constraints. Future transition correctness and malformed diagnostics are explicit
limitations, not guessed rules. Only /tmp/c5b19-g3a7-scope.md was written. Parent
owns acceptance, canonical retention, implementation dispatch and all publication.
No branch, commit, PR, source edit or subagent was created by this assignment.
