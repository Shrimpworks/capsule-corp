# G3a7 author implementation handoff

Status: PASSED for bounded implementation and local verification. Independent
review/publication remain IN_PROGRESS — TRENDING_GOOD. Parent full Gate3 remains
IN_PROGRESS — TRENDING_GOOD; child execution, Gates4–6 and admission remain BLOCKED.
This is author testimony, not an independent readiness verdict.

## Intent and success criteria

Bind a selected kind7 received record and its exact request/reply sidecars to a
kind8 decision record, checking local ordinal/numeric consistency without claiming
the recorded decision was correct. Missing, substituted or inconsistent supplied
members must refuse. Inputs and prior APIs must remain unchanged.

## Plan-to-implementation traceability

Accepted scope: /tmp/c5b19-g3a7-scope.md and parent canonical
.planning/c5b19-received-applied/plan.md. Pair equality, direct sequence/tick order,
<=3 collection, duplicate operation/ordinal refusal and cross-pair tick consistency
are implemented with tests-first and selected guard/restoration evidence.

Archive /tmp/capsule-c5b19-recovery, branch codex/c5b19-benign-teardown,
base/HEAD14365fb4ba5d1525285406ae90feb1abdcc46269. Explicit working-tree boundary:
experiments/typed-guest-transport-c5b19-benign-teardown/{tests/oracle.go,
tests/oracle_test.go,PLAN.md,inputs/ORIGINS.json} only. Canonical reference is
041b6d4e04df9667e6217f57df26f04eec7d2ddc in the parent recovery worktree.
Brain owns archive draft PR38 and new canonical PR on codex/c5b19-received-applied.
No branch, commit, push, PR or subagent was created by this worker. Others' edits
were preserved; all owned processes finished.

## Technical approach and flow

receivedAppliedPair contains Received capturedFramePair and Applied raw bytes.
validateReceivedApplied composes unchanged validateCapturedFrame, then strict
received kind7 / applied kind8 decoding under the same context. It compares six
represented values: Operation, Generation, RequestDigest, CandidateDigest, Outcome
and RecordDigest. Existing strict decoders recompute raw request/candidate/record
hashes, bind exact fields/lane and enforce application codes1/2. Direct comparisons
require received.Sequence < applied.Sequence and received.Tick <= applied.Tick.
No sequence increment, adjacency rule or inferred origin is used.

validateReceivedAppliedPairs validates context even when empty and caps the list
at3 before bounded work. Fixed6-event metadata holds copied records. Duplicate
Operation IDs refuse even across different generations. Pairwise comparisons
require unique sequence values and nondecreasing tick in sequence order for all
supplied endpoints, independent of list order. Caller bytes are never sorted or
retained. A nil error is only correspondence success, no returned capability.

## Changed-component walkthrough

- oracle.go:69 appended lines for pair shape and two error-only validators.
- oracle_test.go:271 appended lines, seven test functions and independent fixtures.
- PLAN.md:appended scope, commands/results, selected mutations, hashes and limits.
- ORIGINS.json:new g3a7 object only; all prior parsed values and false run/admission
  flags remain unchanged.

## Decisions and rejected alternatives

Reuse G3a5/G3a6 strict parsing and standard-library primitives per the C5b19
BUILD-NARROWLY row and ADR47, rather than adding a parser/serialization format or
new dependency. Keep pure correspondence separate from mutable writeHistory and
control transition interpretation. A pair-struct list avoids inventing a sidecar
index and explicitly admits that omitted whole pairs are not detectable.
A bounded quadratic comparison over at most6 records is simpler than sorting,
mutable maps or a new state object; no performance claim beyond this cap is made.

## Invariants and boundary conditions

Exact received sidecars validate before comparison; lower decoder failures refuse
before bytes can be interpreted. Both events share frozen binding/timebase rules.
Both decision codes and all outcomes are representable; matching code1 does not
prove a valid settle, and code2 does not prove a fatal latch. Same ticks and zero
origins are representable; equal/reversed pair ordinals and tick regressions refuse.
Near-uint64 boundaries use comparison only, avoiding arithmetic overflow. Missing
members refuse; empty selection validates context but provides no evidence.
Earlier source prefixes remain byte-identical and opt-in APIs retain their contract.

## Verification performed and results

Exact binary /Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go,
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/c5b19-g3a7-go-cache,
cwd the archive experiment directory.

- `test -count=1 ./...` compiler-red:exit1, undefined pair/validator API.
- Green after implementation/final collision-fixture refinement:exit0,0.568s.
- Restored `test -count=1 ./...`:exit0,0.320s.
- `test -race -count=1 ./...`:exit0,1.858s.
- `vet ./...`, `build ./...`:exit0, no diagnostics.
- `git diff --check`:exit0.

Seven tests cover all operation shapes and outcome/decision combinations; exact
value matching, wrong kind/lane/context/forbidden fields; malformed/short/extra
members and substituted raw bytes; zeros/gaps/same ticks/uint64 limits; cap,
operation reuse including changed generation, every endpoint-role sequence
collision, cross-pair tick regression and reversed input-list order; input
nonmutation and old API preservation. Source inspection establishes no reachable
history, store or lifecycle mutation; no synthetic history tautology is needed.

Fourteen final compiled mutants each failed intended assertions: tuple, request
hash, candidate hash, outcome, record, sidecar validation, received kind, applied
kind, pair sequence, pair tick, operation reuse, cap, sequence reuse and cross ticks.
Exact source restored after each, and restored checks passed. The cap mutant also
expands its metadata array to avoid counting an out-of-bounds panic as sensitivity.
The first cross-tick mutant construction used redundant false expressions rejected
by Go vet; that failure is excluded and the corrected construction compiled and
failed the intended test. No implementation fix was needed for that construction.
Its separate log/hash is retained in PLAN/evidence JSON. The other13 successful
mutation results were retained without unnecessary reruns.

PLAN/ORIGINS retain commands, test names, representative failures and hashes.
/tmp/c5b19-g3a7-mutations.json retains exact replacements; /tmp/c5b19-g3a7-evidence.json
retains source/reference/log metadata. Temporary raw logs are ephemeral and are not
release-grade evidence. Source-prefix and exact module/model byte comparisons and
historical ORIGINS parsed comparison passed. Canonical CI is parent-owned, not
claimed as independently executed by this worker.

## Risks, tradeoffs and maintenance costs

The pair uses extra verified event decodes for a small error-only interface; the
wrapper uses at most6 fixed metadata values and bounded pairwise comparisons.
Selected guards are sensitivity evidence, not exhaustive protocol proof. Future
composition must still validate what the decision should have been.

## Deviations, deferrals and known gaps

No behavior scope deviation. Wholly omitted pairs, full capture completeness,
origin/contiguity, worker receipt, actual pending/control state, code1 settle,
code2 fatal latch, stop/apply causality, durable snapshots, expected process
projections, deadline math/native effects/durability/lifecycle claims remain
outside this slice. Arbitrary malformed/truncated rejected-reply mappings remain
unsupported, not redefined as universally invalid. Child execution stays gated.

## Challenge points for reviewer

Check exact field equality and raw-sidecar validation; common context rejection;
local pair and cross-pair ordinal/numeric invariants including zero/max boundaries;
independently constructed collision fixtures; whether docs improperly promote
represented decisions to correctness or confirmation; preservation of lower APIs,
source/provenance, and whether mutation logs are compiled assertion failures.
Review both the accepted plan and source, without adopting this author testimony.

## Frozen scope hashes

- `tests/oracle.go`: `b9cf1d6c311da52062eb86aec40f2e1cce1f03767bfbb5f6686a23d4746672bb`.
- `tests/oracle_test.go`: `2ae581d0ff919974d73319fd953580f878435ae4693a65186e276a4d081fabf6`.
- `PLAN.md`: `5c187764aedf50f1bc6a65f961b2f4d7d0ed6b43a6bf05d2bf3fbac7f74861c7`.
- `inputs/ORIGINS.json`: `c161dd22e021cc2f11467a60fa171d1317029982ed1bf714b4cf9965a18da5ab`.

Unchanged normative references:
- `docs/C5B19_GATE2_IMPLEMENTATION_PACKET.md`: `65c885425298ce056f8619c5d0837f97ba54dee91c025fea5621a2d75d5ef4f0`.
- `internal/execution/teardownpassive/types.go`: `249b7ef054134f47974bf769c3fb37309b867b1be42ae53c18baf23dc8bacfd0`.
- `internal/execution/teardownpassive/model.go`: `1e3a8be8527554d3953d494086c5313ceeca5806b40a7661ab488717f4a94c40`.

Log hashes (paths /tmp/c5b19-g3a7-NAME.log):
- `red`: `7548ce788344528f2ccc290c70ac00090c145bb6f522d7a5ba1d49de06bc820c`.
- `green`: `7b209b31e212d5292cdc71945f403acdbb11a59eb53e811ef1b2c3645aed8354`.
- `restored`: `61ac622f9401eb2ce3bc459d6d0ced2aa686e092447a2251e5e4e4966b0fabb9`.
- `race`: `e5f4e343c4a2730b5a6261772f941b55d1519ac2ce8606fe937b03724a7b34d3`.
- `vet`: `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.
- `build`: `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.
- `mutant-tuple`: `151909881a2b0a4b191f03a3a74662f76c26b3c089960148ed1eea36719268eb`.
- `mutant-request-hash`: `691df4370522d978d57d73eb273809ff55cd0389b0ec4c57dfaf837bbc96b237`.
- `mutant-candidate-hash`: `05ed16dec1dfdb41f36a25afc9c22c3352c4f04c85c65068a6af0005760bab9f`.
- `mutant-outcome`: `44ad1fc9eaca1794acc6db235b52df01b4a39736a16499945633200c95ed8aa4`.
- `mutant-record`: `2f16ba624ede818338a09623e47a35edfb27d31c00ba8b79a0f41121cf4bf9bc`.
- `mutant-sidecar`: `4d7eb5c4674d25ceae53e2d54aab4b45e2621c3d182ff45bb471785ecc437714`.
- `mutant-received-kind`: `4ea5874906c18ec27cc5a19fef3a7228b11236ef732737b5fa9e94bbbf232204`.
- `mutant-applied-kind`: `e6094046435f15f698346cbfe72337a24f5e9758e725132d2299762ca1203c16`.
- `mutant-pair-sequence`: `bd2c0b90f72519455a73cb377febd687f39bf810a1df7559355c5cc8a0b7b281`.
- `mutant-pair-tick`: `4c579c7434a20225478702f1371c0c51670ca459924a84750ccb453f3c81c701`.
- `mutant-operation-reuse`: `c0b670138a96308841169fdb6f80968e8319b07628c6d18c2fe9986421f782fa`.
- `mutant-cap`: `d283da4eaad3294e6b0b78472a884ac22f40c9a7186811fd0a00ed2f74910c9f`.
- `mutant-sequence-reuse`: `85d4e11ab329fbf08f8869c1ae30ae390b128b85ef71ef9a2f9581448bb7593d`.
- `mutant-cross-ticks`: `def79d9ed16e0443820acdd03b36d06dc886b8fdf65801eefd7ff1db47e7db0c`.

## Retrieval and authorized environment

Serena initialized; exact archive source is unindexed and was read directly. Parent
supplied Tier2 graph generation2026-10-04 primary-checkout evidence, not watched:
packet/types/model best-effort metadata_match, SettleWrite exact pending/outcome
before applyDurableProjection, four-callee outbound trace complete. Current reference
hashes agree. No graph completeness claim. Relevant repository boundaries, packet,
ADR47, TDD skill, scope, prior APIs and normative model context were consulted.
Only owned inert bytes and ordinary offline unit-test processes were exercised;
no child/native fixture/guest, lifecycle target, live clock/store/network,
credential or user-data operation. Return to Brain for independent review and
publication; no owned process session remains active.
