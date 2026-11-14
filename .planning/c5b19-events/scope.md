# G3a4 scope: raw event envelope and lane framing

Status: PASSED for this read-only scope assessment. Parent full Gate 3 remains
IN_PROGRESS — TRENDING_GOOD. Gates 4–6, child execution and product admission remain BLOCKED.

## Recommendation and rationale

Deliver a pure **event-envelope decoder and bounded lane-framing checker** as G3a4.
G3a3 already establishes synthetic write-history relations. The next full oracle
must consume the frozen 256-byte event slots; a small independently tested decoder
and lane structural validator are useful prerequisites before the much larger
29-kind field-use, paired-frame, causality and timing oracle. Do not combine those
later concerns into this slice or call envelope validity an event/case pass.

Canonical baseline: 2aefb15ee45018bb1045e711c77a0c7bc9f4c54d, merged PR #387.
Archive baseline: 6575898a32ff4a3d24ca4c757f42b1807798a372, existing draft PR #38,
branch codex/c5b19-benign-teardown. Only existing frozen-manifest source files
`tests/oracle.go` and `tests/oracle_test.go` need change, plus lead/worker-agreed
PLAN.md and inputs/ORIGINS.json provenance. No new source file or dependency.

Defensive control: reject malformed or substituted synthetic trace envelopes,
nonsequential lane records and over-budget lane inputs. Exact authorized target:
owned archive checkout's inert in-memory byte arrays under Go unit tests. No
child fixture, guest, lifecycle target, signals/waits, store, network, credentials,
real clock, service or user data. Ordinary compiler/unit-test processes are not
execution of the gated experiment fixture.

## Source contract and exact layout

Gate2 packet lines309–327 gives capacities, slot encoding, closed codes and zero
reserved bytes; lines329–337 gives the closed1–29 kind set. Lines377–380 requires
lane-local sequence, tick/timebase and binding; lines396–408 isolates lane state;
lines427–435 rejects absent/truncated/overflowed evidence and forbids cross-lane
clock total-order assumptions. Lines446–453 requires test-first inert unit tests
before separate Gate4/first-run authorization. Lines53–64 freezes source files.
Lines291–296 requires nonzero timebase and rejects mixed clock/regression.

Offsets below are mechanical sums of the packet's explicitly encoded fields;
all multi-byte integers are big endian. Signed fields use signed32 decoding.

| Offset span (half open) | Field |
| --- | --- |
| 0:4 | C5EV magic |
| 4:5 | version1 |
| 5:6 | lane1 control /2 worker /3 harness |
| 6:8 | kind u16, closed1–29 |
| 8:16 | lane sequence u64 |
| 16:24 | tick u64 |
| 24:32 | auxiliary tick u64 |
| 32:36,36:40 | timebase numerator/denominator u32 |
| 40:72 | SHA256 exact Bindings[240] |
| 72:88 | operation ID[16] |
| 88:96 | generation u64 |
| 96:128 | request digest[32] |
| 128:160 | process identity[32] |
| 160:164 | PID u32 |
| 164:168,168:172 | syscall result/errno i32 |
| 172:176,176:180 | flags/wait status u32 |
| 180:181 | outcome0–3 |
| 181:182 | application result0–2 |
| 182:214 | durable record digest[32] |
| 214:246 | candidate digest[32] |
| 246:256 | exactly10 zero reserved bytes |

Use existing validated `encodeObligationPreimage(frozen)` to obtain the exact
Bindings[240] at preimage[30:270], and SHA256 those bytes, **not** the complete
obligation preimage or request. Independently constructed test bytes must verify
that distinction and vary each valid frozen binding.

## API and acceptance boundary

Suggested equivalent narrow unexported APIs:

- `decodeEventEnvelope(raw []byte, frozen obligationBindings, numer, denom uint32)
  (eventEnvelope, error)` returns a fully value-owned decoded record. Validate
  fixed raw length before slicing, magic/version, known lane/kind/code ranges,
  all reserved bytes, valid frozen binding set, exact computed binding digest,
  and exact nonzero frozen timebase. Decode every field at its specified offset.
  All failure results should return the zero event plus error. No slice aliases.
- `validateLaneEnvelope(raw []byte, expectedLane byte, firstSequence uint64,
  frozen obligationBindings, numer, denom uint32) error` validates nonempty,
  whole-slot, bounded input before iterating and calls the decoder. Every slot
  must have the expected lane and exactly firstSequence+index, with overflow
  rejected rather than wrapping. Caps:32768 control,1024 worker,64 harness.
  Control ticks must be nondecreasing; equal ticks are valid and input sequence
  order must be preserved. Do not compare or sort timestamps across lanes.

The packet does **not** explicitly freeze whether the first sequence is0 or1.
Do not invent that ABI choice. Require expected-first as an explicit trusted
synthetic-fixture parameter; never infer it from the first untrusted record.
Tests must cover both origins and checked arithmetic near MaxUint64. This
parameter establishes only continuity relative to the supplied expected origin,
not proof of a complete real capture. Selecting the future writer's origin and
its integration into the complete oracle remains deferred. No new packet
revision is necessary to implement this expressly parameterized prerequisite.

For temporal validation, enforcing only the explicitly private control watermark
is the minimum useful boundary. Worker/harness ticks remain independent
observations; they never advance control state. Numerator/denominator matching
prevents changing the declared conversion within any checked lane. No tick-to-ns
math, clock provenance, timing bounds, tick/auxiliary anchor meanings, or arbitrary
nonzero tick assertion is added here.

**Deliberately deferred:** per-kind allowed/required payload fields; lane-to-kind
ownership (e.g. worker-only26/27, harness-only29, control-only8); operation/request/
candidate correspondence and sidecars; syscall/PID/wait-status semantics;
received/applied pairs; stop/start/absence/terminal causality; trace-fault failure
classification; overflow latch provenance; complete-lane/end markers; cross-lane
merge and timing; all native/store/control/build/run work. This decoder should
not silently pretend to implement those relations. A structurally decoded kind28
trace-fault or wrong-lane payload is still invalid case evidence for a future
semantic checker. Name APIs and comments to make this explicit. No API returns
MayCreate/MayStart/MaySignal/MayRelease or a complete-trace pass verdict.

No mutable tracker is needed: validate finite byte slices without retaining or
sorting them. Existing writeHistory remains unchanged and unconnected; successful
envelope parsing must not settle any write or alter its projection.

## Test-first acceptance and sensitivity

1. Add failing tests before new decoder/checker symbols; retain the compiler-red
   command/output excerpt and hashes in archive evidence.
2. Literal known-answer256-byte fixture, independently encoded from the table,
   checks every decoded field (including negative syscall result/errno and high
   unsigned values). Do not use the new decoder as its own test encoder. Inputs
   and returned value mutations must not affect each other.
3. Reject nil/short/trailing raw frames, every header/code/range/reserved mutation,
   invalid frozen bindings, stale digest, full-obligation hash substituted for
   binding hash, zero/mismatched timebase. Change each valid binding with matching
   digest and confirm acceptance; stale original must fail.
4. Lane tests cover all3 lanes, wrong/unknown expected lane, empty input, partial
   final slot, exact cap, cap+1, every duplicate/gap/reorder position, explicit
   origins0/1 and high origins, and overflow including max-to-zero wrapping.
   Validate unequal tick values from another lane cannot affect control checks.
   Control tick regression fails, equal ticks remain in original sequence.
5. Boundary tests deliberately show envelope-only scope: arbitrary non-common
   payload is copied faithfully and not granted semantic validity; no full-run
   success verdict or history mutation exists. Tests must not imply a synthetic
   full-cap stream is a valid causal trace or trace-fault event is a passing case.
6. Controlled one-at-a-time source mutants should remove/alter (a) binding digest
   comparison, (b) reserved-byte refusal, (c) exact sequence check, (d) lane cap,
   (e) control tick regression check. At least these targeted tests must fail for
   their intended reason. A field-offset/endian mutant is valuable if small.
   Restore exact clean source and rerun the full suite afterward. Retain command,
   failing assertion, source restoration hash and log hash. This is selected
   sensitivity evidence only, never exhaustive proof.
7. Exact pinned offline Go1.25.13 test/race/vet/build and diff check pass. No test
   invokes fixture, os/exec, real clock, filesystem, network, goroutine or live
   store. No dependency/go.sum/source-manifest changes.
8. Fresh-context independent implementation+plan review, at most3 instances for
   distinct G3a4; preserve G3a/G3a2/G3a3 historical counts and claims.

## Dependency policy and delivery

Ecosystem row at docs/ECOSYSTEM_REUSE_AND_ADOPTION.md:135 selects
ADOPT-PLATFORM / BUILD-NARROWLY. Only BUILD-NARROWLY fixed byte decoding applies
here; no platform API is exercised. Record test-only evidence trust class,
Go1.25.13/standard-library-only provenance, explicit per-lane memory/iteration
caps, no authority effects, atomic refusal, offline reproduction, malformed/
substitution/overflow/restoration tests, and Capsule maintainer upgrade/removal
ownership. Archive stays the implementation destination; canonical docs retain
scope, exact commit pins, verification, review and deferred gates. No product
imports or ADR responsibility change. No new external research is needed for
this already-frozen byte contract.

## Method, confidence and limits

Read exact archive source/tests/PLAN/module and canonical Gate2, C5b18, ADR47 and
reuse row; checked archive status (clean). No implementation or tests executed.
Inherited parent graph evidence and independently checked cited canonical paths:
project capsule-corp generation2026-10-04T22:18:10Z, metadata_match/no_recorded_issue,
fully returned5-path coverage page, best-effort only. Archive is not indexed and
was read directly. Canonical source is authoritative over graph.

High confidence in layout/caps and explicit separation from semantic evidence;
sequence origin is genuinely unspecified in the cited packet and deliberately
not resolved by this API. If lead wants per-kind semantic validation in G3a4,
that is a larger acceptance boundary requiring a new exact29-kind matrix rather
than ad hoc inferred guards. Recommended next semantic slice after this one:
closed field-use/lane ownership with raw transport sidecar correspondence, then
causal/timing composition.

Sole retained artifact: /tmp/c5b19-g3a4-scope.md. No branch, commit or PR created.
Integration destination: lead's canonical .planning continuation and archive PR38.
