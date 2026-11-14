# G3a6 captured-frame correspondence scope

Status: PASSED for read-only scope derivation. Implementation/review remain pending;
parent Gate 3 is IN_PROGRESS — TRENDING_GOOD. Child execution and product admission
remain BLOCKED. Internal research assignment; no repository edits or publication.

## Question, evidence and recommendation

Defensively validate equality between owned inert event-byte fixtures and exact
captured request/reply bytes. All work is offline in the two named repositories;
no child, guest, process-control target, clock, live store, network or user data.
Canonical base e1a966e7c6f18665784984e97b6697cc472cb288; archive base
7aac77de36e33059ef0759b44ca02ab5ad080343.

Gate2 packet lines199–215 freezes request397/reply98 and request-body hash;
lines342–360 requires copied sidecars and independent hashing; line367 explicitly
pairs event kinds3/5/7 with raw slots. Its serialized sidecar container/index is
NOT defined. Implement a pure in-memory opt-in correspondence API; do not invent a
new persistent wire envelope, extend the frozen manifest, or alter previous APIs.
The C5b19 reuse-map row recommends BUILD-NARROWLY, using current standard-library
fixed parsing/SHA-256, not a new dependency. ADR47 and C5b18 SettleWrite keep actual
control application/settlement distinct from worker bytes.

## Minimum pair contract

A pair validator takes raw event bytes, exact raw request anchor bytes, optional
raw reply bytes, frozen bindings and timebase. It first invokes decodeEventFields
and decodeRequest. Only kinds3/5/7 are supported; every other kind refuses. An
error returns no usable decoded output. A copied event result or error-only API
is sufficient; no mutable history or caller-owned input references are retained.

| Event | Lane | Raw capture and equality |
| --- | --- | --- |
|3 write-request|control1|Request anchor is the capture: exact397, decodeRequest valid; operation and generation equal decoded request; RequestDigest=SHA256(request[8:397]); CandidateDigest=SHA256(request[305:373]); reply input absent.|
|5 worker-reply|worker2|Exact98 reply plus exact397 request anchor; validateReply succeeds against that anchor; event operation/generation/request/candidate equal request; event Outcome=reply[65], RecordDigest=reply[66:98].|
|7 control-reply-received|control1|Same byte checks as5; lane/kind remain control-specific. This is represented well-formed reply correspondence, not acceptance or settlement.|

validateReply already recomputes confirmed record SHA256 from the exact request
fields reordered into395 bytes; failed/indeterminate reply digest mustzero. Reuse
it rather than introducing a second hash/record codec. Event request hash excludes
the8-byte header; candidate hash includes exactly68 bytes. Exact fixed headers
still have to validate. Do not compare only caller-supplied digests.

## Bounded collection, if included in this slice

A small wrapper is justified to enforce the packet's at-most-three captures per
supported category and reject missing/orphan/duplicate supplied evidence. Keep
its contract explicitly about the supplied selection, not full trace completeness.

One suitable representation is separate bounded lists of raw selected events and
in-memory tagged captures, plus at most3 unique raw request anchors. A capture tag
names only the supported source role (3 control-request,5 worker-sent-reply,
7 control-received-reply); it is an API discriminator, not a serialized protocol
field. Pair by category and decoded exact operation identity/generation; enforce
one-to-one matching. Each category has0..3 events and captures. Reject duplicate
operations within a category, missing captures, extra captures, unsupported tags,
and invalid/truncated/extra-byte entries before reporting success. Request anchors
must be structurally valid and unambiguous (one exact request per operation).
A pair-struct list may instead be used if it explicitly refuses missing event or
frame members; document that omitted entire pairs cannot be detected.

For replies, the request anchor provides decoding/hash context only. It must not
be presented as worker receipt, delivery or a captured kind3 event unless that
separate supported pair is supplied. Do not invent an event number for worker
received-request capture; that packet category has no named event and is deferred.

Where both worker5 and control7 pairs for the same operation are supplied, require
identical reply bytes. Where a request3 pair and a reply anchor for that operation
are both supplied, require identical request bytes. A mismatched raw frame cannot
be repaired by matching digest fields or a different request supplied per call.
Do not require every request3 to have a reply5 or7, or every worker5 to have7:
short writes, drops, pending work and incomplete evidence selections exist. Do not
infer success, missing-event fault, causal order or capture completeness from the
absence of an entire pair in this local API. A full oracle will assess those facts.

## Tests and mutations

Use independently assembled literal-offset fixtures, all three operation kinds,
all reply outcomes, and both reply writer categories. Positive cases include
request-only, worker-only reply with request context, received-only reply with
request context, matching supplied peers, and up to3 captures per category.

Negative cases: wrong kind/lane/common binding/timebase; wrong event operation,
generation, request hash, candidate hash, outcome or record hash; changed raw
request field (including candidate) with stale event digests; wrong reply role,
kind, request hash, record digest or length; nil/short/long event/request/reply;
nonconfirmed record substitution; missing/orphan/duplicate supplied members;
4th capture; inconsistent request anchors and disagreeing worker/control replies.
Test whole-request hashing instead of body-only hashing fails, and no input or
prior state mutates on refusal. Preserve earlier syntax/fields/history APIs.

Selected guard mutations should demonstrably fail specific tests, then restore
source exactly: tuple comparison, request hash, candidate hash, outcome/record
comparison, category/lane, duplicate/cap and peer equality (as applicable to chosen
API). Compiler-red, exact offline Go1.25.13 test/race/vet/build and fresh independent
review remain mandatory. No child execution is needed or authorized.

## Deliberate limits and next boundary

This slice does not decode arbitrary malformed/rejected reply diagnostics; those
encodings remain unspecified. It does not apply kind8, settle writeHistory, assert
received/applied pairing, compare a durable snapshot, derive expected projections
from actual process observations, validate lane sequencing/capture origin, native
syscall result, causal/timing order, full-case success, storage durability or
lifecycle authority. Existing G3a4 lane framing and G3a5 field validity are separate
layers. Empty supported collections may mean no represented correspondences,
never a valid complete trace. A future composition slice must close event8 and
received/applied ordering plus model transitions before confirmation claims.

## Method, confidence and retention

Read exact current packet sections, archive oracle.go decoders/history/field layer,
C5b18 types/SettleWrite, ADR47, reuse row, repository instructions and relevant
starter-document boundaries. No tests were run because this assignment changes no
behavior. Exact HEADs were read back. Graph index status was ready at generation
2026-10-04T22:18:10Z on the primary checkout, not watched; cited canonical paths
reported best-effort metadata_match. SettleWrite outbound trace returned4 callees,
no continuation; source snippet and current worktree reads confirmed its role.
Archive source is unindexed and read directly. No graph completeness claim.

Confidence is high for pair correspondence and the stated count cap, with the
collection API a narrow implementation choice rather than a frozen wire format.
Parent should retain this scope in canonical planning before implementation,
select the smallest representation, and pass its exact acceptance contract to the
fresh reviewer. No branch, commit or PR created by this assignment. Only this
/tmp report was written; integration destination is the parent's canonical PR.
