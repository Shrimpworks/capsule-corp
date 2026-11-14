# G3a4 implementation handoff and author explanation

Status: PASSED for bounded implementation/local verification; independent review
IN_PROGRESS — TRENDING_GOOD, pending instance1/3. Parent full Gate3 remains
IN_PROGRESS — TRENDING_GOOD. Child execution, Gate4–6 and product admission BLOCKED.

## Intent and success criteria

Add pure256-byte raw event-envelope decoding and finite per-lane framing as the
next prerequisite for the future C5b19 oracle. This rejects malformed common
fields, wrong binding/timebase, partial/over-cap lanes, wrong sequence and control
tick regression. It does not decide event semantic validity, trace/case success,
lifecycle permissions or timing evidence.

Defensive authorized scope: owned archive checkout /tmp/capsule-c5b19-recovery,
synthetic in-memory bytes under ordinary unit-test tools only. No experiment
fixture/child, guest, lifecycle target, process probe, signal/wait, real store,
network, credential, service, user data or clock effect in new code/tests.

## Frozen review boundary and integration destination

Archive branch codex/c5b19-benign-teardown, HEAD/base
6575898a32ff4a3d24ca4c757f42b1807798a372. Explicit working diff only:
experiments/typed-guest-transport-c5b19-benign-teardown/{tests/oracle.go,
tests/oracle_test.go,PLAN.md,inputs/ORIGINS.json}. All previous source bytes are
unchanged prefixes; model_test.go/go.mod match baseline exactly. Every historical
ORIGINS field is equal to its baseline value; only g3a4 was appended. Source is now
frozen for root's fresh review. No other archive changes, branch, commit, push or
PR was created. Destination remains existing draft archive PR38; root owns all
publication and canonical immutable pins in its codex/c5b19-event-conformance PR.

Canonical packet reference2aefb15ee45018bb1045e711c77a0c7bc9f4c54d is unchanged.
Exact accepted scope: /tmp/c5b19-g3a4-scope.md, with Gate2 lines309–327,377–380,
396–408,427–435 and446–453. Ancestor graph evidence reports best-effort
metadata_match for canonical packet/types/model at generation2026-10-04T22:18:10Z;
archive is unindexed and source was read directly. No structural graph absence
claim is made. Serena manual read and primary Capsule project already active.

## Plan traceability and technical flow

1. Added failing parser/lane tests before any new implementation; retained compiler
red for missing eventEnvelope/decodeEventEnvelope/validateLaneEnvelope.
2. decodeEventEnvelope checks exact length/header, enums, reserved bytes,
validated frozen bindings and SHA256(exact Bindings[240]), nonzero exact expected
timebase; it then returns a fully value-owned struct decoding every field. All
errors return the zero struct. Signed32 fields explicitly convert from big-endian.
3. validateLaneEnvelope checks valid expected lane, nonempty whole-slot length and
its fixed count cap before looping, rejects sequence arithmetic overflow, and
checks every decoded slot's lane/expected index and nondecreasing control tick.
No sorting, retention, shared watermark or writeHistory connection exists.
4. PLAN records scope/policy/TDD/mutations/results. ORIGINS appends current source,
canonical-reference and evidence-log hashes; execution/admission remain false.
5. All local checks and six sensitivity mutants pass their intended outcomes;
restored clean source suite passes. Fresh review and publication remain root work.

## Design choices and invariants

Use existing bytes/binary/SHA256/errors and validated obligation preimage rather
than introducing a frame library or new dependency. The event digest uses only
preimage[30:270], excluding obligation policy/domain. Test literal was separately
assembled with Python struct/hashlib; changed-binding tests compute raw240-byte
hashes independently. No C/native source or ABI changes.

Gate2 does not specify first lane sequence0vs1. The checker therefore takes an
explicit trusted firstSequence, never learns it from untrusted bytes. This is a
parameterized prerequisite, not a new writer convention or complete-capture claim.
Only the control watermark is enforced; worker/harness observations cannot advance
it. Non-common fields and known codes are copied, not semantically approved.

Current raw caps are32768/1024/64 events for control/worker/harness. The wrapper
checks size before iteration; decoder scratch is bounded. No arbitrary input cap
or hidden mutable tracker is introduced. Existing write-history code is untouched.

## Verification and sensitivity

Exact tool for every run:
GOCACHE=/tmp/c5b19-g3a4-go-cache GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go

- test ./... -count=1 initially compiler-red with missing symbols.
- test ./... -count=1 after implementation exit0 (0.260s).
- test -race ./... -count=1 exit0 (1.786s).
- vet ./... and build ./... exit0, no diagnostics.
- Restored test ./... -count=1 exit0 (0.306s).
- git diff --check exit0.

The initial shell wrapper attempted to assign zsh's reserved status variable after
the red Go run; only its exit-report wrapper failed. Red Go output was already
captured and inspected. Subsequent commands use result_code. Tests themselves
contain no shell/subprocess/file/network/clock work.

Coverage: all decoded fields including signed/high unsigned values, every short
frame length, trailing byte, all reserved bytes, known codes1044combinations,
unknown code boundaries, individually valid/zero changes to all10bindings,
representative duplicate ID/digest with matching raw binding hash, stale or whole-
obligation hash substitution, invalid timebase, ownership, all3 lane exact caps
and cap+1, three origins, 3-slot per-position gaps/duplicates/reorders/mixed lanes/
malformed frames, max-u64 wrap and single-max, equal/regressed ticks, lane independence.

Six one-at-a-time mutants all compiled and failed their intended assertion:
binding comparison, reserved-byte refusal, exact sequence comparison, capacity
cap+1, control tick comparison and signed field offset. /tmp/c5b19-g3a4-mutations.json
records test command/name, output hash and exact restored source hash. No mutant
is retained. Original-source SHA256 confirmed after each restoration. Evidence
is selected sensitivity, not exhaustive proof. Logs at /tmp/c5b19-g3a4-*.log;
essential commands/results/excerpts and hashes retained in PLAN/ORIGINS.

## Risks, costs, deviations and deferred work

No intended scope deviation. No new dependency, go.sum, frozen source file or
product import. The decoder adds103 lines and tests288, intentionally keeping
kind-specific/causal semantics out. Passing syntax may include semantically wrong
lane-kind/payload combinations or trace-fault kind28; comments and literal test
make that limit explicit. A future semantic layer MUST reject those as required.

No measured real capture completeness, received/applied pairing, fsync, custody,
authoritative absence, live clocks/bounds, native/C conformance or start/stop/
completion/release evidence. The future writer's sequence origin and full oracle
integration remain deferred. Race-test success is not a concurrent API guarantee.
Raw log files are temporary; retained excerpts/hashes survive publication. The
canonical Go1.26.11 patch is separate parent-owned prerequisite; archive1.25.13
remains exact and no runnable host profile is admitted.

## Reviewer challenge points

Verify offset arithmetic and signed conversions, Bindings-only hash versus full
obligation hash, zero-result refusal, explicit trusted sequence origin and overflow,
all exact caps, the control-only watermark and absence of cross-lane state. Verify
no API or status promotes syntax into semantic or lifecycle evidence. Inspect
literal independence and whether selected guard mutants fail for intended reasons.
Read the neutral bootstrap/scope and source first; this packet is author testimony.

## Frozen SHA256

- `tests/oracle.go`: `4d7010b67089bd8d25e3c6f35101e9652c66e7147bcb5210e2103b3314f34f4b`
- `tests/oracle_test.go`: `d134a318eb49335c3113e35fb42928ebc427022612413b7925c5fbe514cdc234`
- `PLAN.md`: `5ff9ecd39a49909408059ad8f813cf93e84611c8e4b0c9d39ad2a647d0611b32`
- `inputs/ORIGINS.json`: `838f9628eb886a914eaa84143d7bb2fc54d3278ac6d6c9c5038d382662447845`

Recommended next action: root retains this handoff, starts fresh review1/3,
reconciles findings, then commits and publishes exact reviewed bytes.
