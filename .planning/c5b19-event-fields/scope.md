# G3a5 scope — event field shape and lane ownership

Status: PASSED for this read-only scope assessment. Full Gate3 remains
IN_PROGRESS — TRENDING_GOOD; child execution, Gate4–6 and product admission BLOCKED.

Canonical base fe3d22c84b5d41e0a4edbfacb531563040e16995 (merged PR402), branch
codex/c5b19-event-fields. Archive base3e26b809a40832a3351fdc6b5e0ffc1ef0094792,
existing codex/c5b19-benign-teardown / draft PR38. Parent owns integration.

## Question, method and recommendation

Question: what local kind-specific checks can safely follow G3a4's syntactic
256-byte envelope without pretending to validate complete run evidence?

Recommend one pure decoder/validator for the table-defined local event shapes
below, enforcing forbidden-field zeros, required tuple/identity presence, closed
site/outcome/application values and writer ownership. Keep G3a4 API behavior
unchanged; callers explicitly opt into G3a5. No sidecar lookup, history settlement,
causal pairing, timing arithmetic or syscall/wait-status interpretation.

Defensive scope: malformed synthetic in-memory event bytes in owned archive unit
tests only. No fixture child, guest, process target, store, network, credential,
real clock, signal/wait, service or user data. No dependency/new source manifest
entry. Existing oracle.go and oracle_test.go suffice; PLAN/ORIGINS retain evidence.

Read current Gate2 packet, Gate1 custody explanation, current eventEnvelope/
decodeEventEnvelope/validateLaneEnvelope source, new canonical plan, and inherited
ADR47/reuse context. No implementation or tests executed. Archive source is
unindexed and read directly. Cited canonical coverage was checked: generation
2026-10-04T22:18:10Z, metadata_match/no_recorded_issue, full4-path page, best effort.

## Source rules and notation

Authoritative packet: docs/C5B19_GATE2_IMPLEMENTATION_PACKET.md at canonical base.

- Lines83–94: control owns creation/cancel/start/signal/wait; worker owns storage.
- Lines122–131: positive signed-Darwin-pid_t PID and nonzero identity; identity is
  only an opaque binding, never process authority.
- Lines187–225: three fixed write kinds, generation/candidate/binding constraints,
  reply outcomes1–3 and zero record digest for failed/indeterminate outcomes.
- Lines308–323: lane1 control,2 worker,3 harness; existing envelope fields/codes.
- Lines325–336: names1–29, harness-only29, flags only for store-edge A–D;
  wait-status only for exact terminal wait result; every unmentioned field zero.
- Lines338–350: transport tuple and nonzero confirmed record digest.
- Lines364–375: exact per-kind field table.
- Lines377–393: common sequence/tick/timebase/binding, operation binding for3–8
  and26–27, control-only8, and causal obligations deferred here.

All matrix rows inherit G3a4 magic/version/reserved bytes, known lane/kind/code,
frozen valid binding hash and nonzero matching timebase. These common fields
remain present: Lane,Kind,Sequence,Tick,Numerator,Denominator,BindingDigest.
All **non-common fields not listed in a row MUST be zero**.

Abbreviations:
- T = Operation[16], Generation, RequestDigest[32], CandidateDigest[32].
  Local presence: nonzero Operation, Generation1..2^53-1, nonzero request/candidate
  digests. This denotes a claimed exact operation only; actual equality to raw
  request/candidate bytes remains a sidecar check. The packet explicitly names
  this tuple for3–8; applying the same tuple to26–27 follows “additionally bind
  the exact operation” at378 and “operation tuple” at374.
- I = nonzero ProcessIdentity[32] plus positive PID<=2^31-1.
- R = SyscallResult and Errno, copied signed32 observations; both may be zero.
  Do not impose syscall-specific result/errno ranges or success combinations.
- O = Outcome1..3. Zero denotes no outcome and cannot stand in for a required O.
- D = RecordDigest: nonzero iff O==1; zero for O==2/3.
- A = AuxiliaryTick, carried for the named anchor; zero is not itself proof of a
  missing event. This slice does not establish observed anchor provenance.
- S = Flags1..4 exactly (A/B/C/D store site).
- W = WaitStatus. Zero is allowed: terminal successful exit can encode zero.
  If nonzero, require SyscallResult==int32(PID), a necessary local condition.
  This does NOT establish terminality, interpret Darwin status bits, or prove a
  matching prior wait call. No errno==0 rule is invented.

## Exact 29-kind supported matrix

| Kind/name | Lane | Allowed/required non-common fields | Source |
| --- | --- | --- | --- |
|1 endpoint-created|control1|none|366; control endpoint owner83–94|
|2 setup-check|control1|A|366|
|3 write-request|control1|T|338–340,367|
|4 request-write-result|control1|T,R|339–340,367|
|5 worker-reply|worker2|T,O,D|340–342,367|
|6 worker-reply-write-result|worker2|T,R,O,D|340–342,367|
|7 control-reply-received|control1|T,O,D|343,367|
|8 control-reply-applied|control1|T,O,D; ApplicationResult1 or2|344–350,367,379|
|9 spawn-call|control1|none|83–94,368|
|10 spawn-result|control1|PID (zero or <=2^31-1),R; ProcessIdentity zero|368; successful/failed PID semantics deferred|
|11 identity-frozen|control1|I|368;122–131|
|12 cancel-ingress|control1|none|369;227–234|
|13 start-admission|control1|none|370,383|
|14 start-token-call|control1|none|370,383–384|
|15 start-token-result|control1|R|370,384–385|
|16 start-post-service|control1|A|370,385–386|
|17 wall-service|control1|A|371|
|18 fatal-service|control1|none|371|
|19 stop-latch|control1|A|371,389–390|
|20 signal-latch|control1|I|372|
|21 signal-call|control1|I|372|
|22 signal-result|control1|I,R|372|
|23 wait-call|control1|I|373|
|24 wait-result|control1|I,R,W|373,391–392|
|25 absence|control1|I; WaitStatus zero (belongs only to24)|335–336,373|
|26 store-edge|worker2|T,S; Outcome/RecordDigest zero|334–336,374,378|
|27 store-outcome|worker2|T,O; RecordDigest zero even if confirmed|374,378|
|28 trace-fault|any1/2/3|none|375 “lane-local trace fault”|
|29 harness-containment|harness3|none in supported table-defined subset|334,375; optional identity ambiguity below|

“None” never omits common Tick/binding/timebase. Required integer observations R,
A and W can legitimately be zero; do not implement “every listed field nonzero.”
The frozen packet has no separate presence bit for scalar zero. This local check
cannot establish missing/regressed anchor observations or successful syscall
facts; those need paired/control-time validation. Existing control watermark is
available only through validateLaneEnvelope; this single-record check must not
claim that it has run that function.

## Explicit ambiguities and conservative decisions

1. **Harness identity variant:** lines147–150 permit optional harness-child
   identity binding in its own trace but do not map that binding to kind29 fields.
   The later table375 names no non-common payload. Support the payload-free29
   subset and refuse payload-bearing29 as **unsupported by this local subset**,
   not proof that the whole frozen packet forbids optional identity capture.
   Exact mapping requires later packet clarification before that variant is used.
   This does not block the requested table-defined bounded acceptance; do not
   claim complete harness-trace implementation or enable a run.
2. **Store outcome vs reply:**27 names only operation tuple+outcome; no record
   digest is named. Thus RecordDigest remains zero even for confirmed27, unlike
   transport5–8. Store outcome never substitutes for control settlement.
3. **Rejected reply diagnostics:** support the complete well-formed candidate
   reply local shape for both accepted/rejected kind8 decisions (O1..3,D relation).
   Rejected malformed/unknown-code/truncated frames have no defined diagnostic
   slot mapping here. Do not manufacture Outcome0 or stripped tuple as “known
   valid” kind8. Such diagnostics remain unsupported pending explicit mapping;
   absence of their support blocks full oracle completion, not this subset.
4. **Spawn result representation:** posix_spawn is selected, but the trace does
   not freeze how its return error is mirrored into generic Errno. Allow the
   named PID/R fields, require ProcessIdentity zero, and range-check a nonzero
   PID. Do not infer PID>0 iff SyscallResult==0, or equate result with PID.
   Failed spawn must never acquire identity/custody; that causal check is later.
5. **Wait status:** do not use guessed Darwin WIFEXITED/WIFSIGNALED encodings or
   require W!=0. Nonzero status requires the target PID return, but positive PID
   return plus status0 is merely a local shape, not authoritative absence.
6. **Ticks/anchors:** no generic nonzero-tick rule was frozen by G3a4. Keep scalar
   values without converting them to evidence presence. “Missing or regressed
   observation fails” at366 requires later observation/anchor pairing and bounds.
7. No first-sequence convention is selected; explicit trusted G3a4 parameter stays.

## Suggested acceptance API

Add `decodeEventFields(raw []byte, frozen obligationBindings, numer, denom uint32)
(eventEnvelope,error)` (equivalent narrow name permitted). It first invokes
existing decodeEventEnvelope, then validates the exact matrix. Error returns
zero event. It accepts raw bytes, so direct construction of eventEnvelope cannot
skip common validation. A package-local pure helper may perform field checks.

Use explicit field masks/booleans or a short switch plus a common forbidden-field
check; keep the matrix visibly auditable. Do not alter decodeEventEnvelope's
syntactic-only contract or wire/history APIs. A new lane walker is unnecessary:
future composition can invoke existing validateLaneEnvelope and G3a5 decoding.
No successful return means a valid whole trace, confirmed durable state,
absence, timely stop, authorization, or capacity release.

## Tests and mutations

- Tests first; retain actual compiler-red then green output. Positive literal
  fixture builder independently sets documented byte offsets, not implementation
  mask tables. All29 kinds and every allowed lane (28 on all3) must pass. Include
  both kind8 decisions, all reply outcomes, all4 site codes, store27 confirmed
  with zero record digest, signed syscall observations, and zero wait status.
- For each kind and each forbidden non-common field, independently inject
  nonzero bytes and require refusal. Include kind25 WaitStatus, store27 record
  digest, kind26 Outcome, PID/identity on setup/start/cancel, and ApplicationResult
  everywhere except8. Use all three lanes for every kind to establish exact
  ownership, rather than testing only one chosen wrong lane.
- Missing required T/I members each fail independently; generation0/2^53-1/
  2^53 boundaries; PID0/max/max+1 whereI required; zero/unknown O/Application/site;
  confirmed missing record digest and failed/indeterminate nonzero record digest.
  Test accepted/rejected kind8 alike. For10, positive/zero PID shapes and out-of-
  signed-range refusal, without inventing success relations.
- Wait24 zero status with matching PID succeeds; nonzero status with matching PID
  succeeds as shape only; nonzero status with result0/-1/wrongPID fails. An errno
  field alone is retained observation, not an inferred syscall contract.
- Envelope failures still fail through the new wrapper; all failure values zero.
  Inputs/returned arrays remain copied, no source mutation. Demonstrate that a
  locally valid25 without preceding24 does not claim causal evidence; no actual
  history settlement occurs. Existing broad G3a4 syntactic tests stay unchanged.
- Controlled restore-one-at-a-time mutations: bypass lane ownership; bypass one
  forbidden-field group; omit required tuple/identity presence; wrongly allow
  record digest on27 (or wrongly require one when confirmed27); bypass reply
  O/D coupling; bypass store-site range; omit nonzero wait-status/PID condition.
  Target otherwise-valid fixtures so each mutant fails for its intended guard.
  Restore exact source bytes and rerun suite; record command/assertion/log hashes.
- Exact offline Go1.25.13 test/race/vet/build, diff check. No fixture/exec/realclock/
  filesystem/network/store behavior in tests. Fresh independent review of source
  and plan, max3 instances for this new slice; prior budgets remain historical.

## Dependency policy, confidence and handoff

Applicable reuse row: C5b19 prospective Supervisor control/storage lane exchange,
BUILD-NARROWLY fixed bounded parser. Existing stdlib only; exact pin unchanged;
input256 bytes and fixed29-kind decision tree; no runtime authority or allocation
proportional to untrusted size; validation before return; no product import.
Capsule maintainer owns upgrade/removal. Retain scope/ambiguities, tests/mutations,
review and archive immutable pin in canonical docs.

Confidence high for table-defined field and lane matrix. The optional29 identity
mapping, malformed-reply diagnostics, syscall result conventions and native
wait-status semantics are explicitly not resolved. These keep full Gate3/4/run
closure blocked. This research does not independently verify future code/tests.
Sole artifact /tmp/c5b19-g3a5-scope.md; no source edit, branch, commit or PR.
