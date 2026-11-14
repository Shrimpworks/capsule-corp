# G3a6 author implementation handoff

Status: PASSED for bounded implementation and local verification. Independent
review/publication remain IN_PROGRESS — TRENDING_GOOD. Full Gate3 remains
IN_PROGRESS — TRENDING_GOOD; child execution, Gates4–6 and admission remain BLOCKED.
This is author testimony, not an independent readiness verdict.

## Intent and success criteria

Check exact correspondence between selected kind3/5/7 event bytes and captured
request397/reply98 bytes under frozen bindings and timebase. A mismatched event,
raw frame or peer pair must refuse without mutating anything. Only this bounded
selection is assessed, never a full trace or real publication/settlement.

## Plan-to-implementation traceability

Parent accepted /tmp/c5b19-g3a6-scope.md plus <=3 distinct operation anchors overall.
The pair API, caps, duplicates, request/peer equality, tests-first failure and
mutation/restoration evidence implement that contract. No serialization or source
manifest was invented. All requested owned files are within
`experiments/typed-guest-transport-c5b19-benign-teardown/`.

Archive /tmp/capsule-c5b19-recovery, branch codex/c5b19-benign-teardown,
base/HEAD7aac77de36e33059ef0759b44ca02ab5ad080343. Working-tree review boundary:
only tests/oracle.go, tests/oracle_test.go, PLAN.md, inputs/ORIGINS.json.
Canonical reference e1a966e7c6f18665784984e97b6697cc472cb288 in parent recovery worktree.
No branch, commit, push or PR was created. Parent owns integration destinations:
existing draft archive PR38 and new canonical PR on codex/c5b19-frame-sidecars.
No subagent was spawned. Other workers' edits and unrelated state were preserved.

## Technical approach and flow

`capturedFramePair` carries raw Event/Request/Reply slices only during the call.
`validateCapturedFrame` composes unchanged decodeEventFields and decodeRequest,
checks operation/generation/body-hash/candidate-hash equality, refuses a reply for
kind3, and for kind5/7 composes validateReply with event outcome/record equality.
validateReply already computes the exact confirmed395-byte record digest from
request bytes; this change does not duplicate that codec.

`validateCapturedFrames` first checks bindings, nonzero timebase and <=9 entries.
Each pair must validate; fixed arrays count <=3/category and retain copied event
metadata for bounded O(9^2) pairwise checks. Same-category Operation reuse always
refuses; shared Operation across categories requires identical raw Request bytes;
worker5/control7 peers require identical raw Reply bytes. <=3 distinct Operation
anchors overall models one bounded selection. No resets, storage, retained
references or mutable history are present.

## Changed-component walkthrough

- oracle.go:87 appended lines for pair data, pair validator and collection wrapper.
- oracle_test.go:281 appended lines, seven behavior tests and independent fixtures.
- PLAN.md:appended scope, commands/results, selected mutant evidence, hashes/limits.
- ORIGINS.json:one new g3a6 object. Existing parsed values remain unchanged;
  run_authorized=false and product_admitted=false remain intact.

## Decisions and rejected alternatives

Use in-memory pair structs because the packet does not define a serialized sidecar
container/index. Separate event/capture lists would introduce more mapping state
without proving whole-trace completeness; pair structs make missing members fail
while honestly leaving omitted whole pairs undetectable. Reuse existing strict
parsers/SHA256 and exact frozen Go1.25.13 standard library, consistent with ADR47
and the C5b19 reuse-map row. Do not add dependencies or new primitives.

An error-only result avoids suggesting that correspondence yields a settled
capability. Keep earlier envelope/field/history APIs byte-for-byte unchanged and
opt-in rather than retroactively strengthening their acceptance contract.

## Invariants and boundary conditions

All event kind/lane/forbidden-field checks flow through the existing strict field
layer. Hash exact request body (exclude8 header bytes), exact68-byte candidate,
and exact395-byte record. Exact header/length checks still run. Counts refuse
before unbounded work; arrays have fixed9-event/8-kind sizes and only closed3/5/7
indices reach them. Any missing event/request or required reply refuses. Empty
selection validates context but proves nothing. Same-operation peer request bytes
must match even if supplied events individually validate under different anchors.
No global order, callback, process, syscall, clock, store, network or authority.

## Verification performed and results

Exact binary /Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go,
GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/c5b19-g3a6-go-cache,
cwd the archive experiment directory.

- `test ./... -count=1` compiler-red:exit1, undefined capturedFramePair and validators.
- Same green after implementation/final test additions:exit0 (0.278s).
- `test -count=1 ./...` after restored mutations:exit0 (0.373s).
- `test -race -count=1 ./...`:exit0 (1.819s).
- `vet ./...`, `build ./...`:exit0, no diagnostics.
- `git diff --check`:exit0.

The first attempted red invocation used archive root and failed module discovery;
it was corrected before implementation. That setup error is not counted as TDD
red. Final red log is the undefined-API failure. No gate weakened.

Seven tests cover operation kinds1/2/3, terminal identity/disposition variants,
all outcomes, request-only/worker-only/received-only/empty and matched selections,
nine pairs at cap, all meaningful event tuple/hash comparisons, body-vs-whole hash,
changed valid raw candidate, malformed/short/extra frames, unsupported valid kind8,
wrong lane/context, duplicates (same and changed generation), caps, peer mismatches,
input nonmutation and no implicit pending-write settlement. No real transport or
child was needed or invoked.

Eleven targeted one-at-a-time mutants (tuple, request-hash, candidate-hash, outcome,
record, category, lane, duplicate, caps, anchors, peers) compiled and failed their
intended assertions. Exact original oracle.go restored after each using finally;
restored test/race/vet/build passed. PLAN and ORIGINS retain command/test names,
failure excerpts and hashes. /tmp/c5b19-g3a6-mutations.json additionally records
exact replacement strings; /tmp/c5b19-g3a6-evidence.json records full evidence.
Raw logs/scripts in /tmp are ephemeral; selected results/hashes are retained.

Read-only Python checks proved prior oracle/source test byte prefixes unchanged,
go.mod and model_test.go byte-identical to HEAD, and all historical ORIGINS values
preserved. Full canonical tests are parent/verification-agent owned, not claimed
as executed by this implementation worker.

## Risks, tradeoffs and maintenance costs

This adds a small bounded quadratic collection check (at most9 elements) and a
second decodeEventFields call per accepted pair to keep the pair API error-only.
No independently mutable state or parser format is added. Selected mutants and
unit tests are sensitivity evidence, not exhaustive protocol proof.

## Deviations, deferrals and known gaps

Omitted whole pairs are undetectable; no full capture completeness or delivery
claim. Request anchors used for replies are context, not worker-receipt evidence.
Worker received-request has no named event mapping and remains deferred. Arbitrary
malformed/rejected reply diagnostics lack a frozen encoding and remain unsupported.
Kind8, received/applied ordering, expected projections from actual process facts,
durable snapshots, sequence origins, native semantics, causal/timing validation,
full-case success and live durability/lifecycle claims remain future work.
No child fixture or guest ran; no product admission or run authorization changed.

## Challenge points for reviewer

Inspect whether exact offsets/hash domains and reply/record coupling match Gate2,
whether supplied peers can smuggle conflicting anchors or outcomes, whether count
limits operate independently of inferred causality, and whether docs overclaim
capture or settlement evidence. Verify lower APIs remain unchanged and the
retained red/mutation logs actually show intended failures. No author verdict
should be substituted for the fresh independent review.

## Frozen scope hashes

- `tests/oracle.go`: `5106897e28cfadd491f4080c1c8d2cd3901fc49bbafe9ce8900924e35adac46e`
- `tests/oracle_test.go`: `a44445ff7cf4bbef309205971eb914fead4dcd93ac45b40d499bef66cea42ea4`
- `PLAN.md`: `1099bca9dd78f19cc38148de33f02de1967d0e9c3ff131a965645b88eb27a449`
- `inputs/ORIGINS.json`: `e5de8d7d107ef2fb86950e03b320d72289a6724bad202d5fdf56bbd031f318c3`

Unchanged normative reference hashes:

- `docs/C5B19_GATE2_IMPLEMENTATION_PACKET.md`: `65c885425298ce056f8619c5d0837f97ba54dee91c025fea5621a2d75d5ef4f0`
- `internal/execution/teardownpassive/types.go`: `249b7ef054134f47974bf769c3fb37309b867b1be42ae53c18baf23dc8bacfd0`
- `internal/execution/teardownpassive/model.go`: `1e3a8be8527554d3953d494086c5313ceeca5806b40a7661ab488717f4a94c40`

Log hashes (exact commands and retained excerpts also in PLAN):

- `red`: `51cd1433dc9adb2d07c359001fe36eba1a5395300349112e8191705f54667e7e`
- `green`: `8aa6efdf3a3c9270a2dd274117fdb9f3a1b5d56698c6278841a406615697b79f`
- `restored`: `33bbd8e4efefa78d2115334a95550c4a72afa3f16eb60f820e1c41627f6f3890`
- `race`: `33f3a509f3dc4da9d9da9b3e3f992d22c74c9615ad96ef14f1fa5aa0eaae6157`
- `vet`: `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
- `build`: `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
- `mutant-tuple`: `e9fc8dfa750eaf88d7871519fe221cf137f59375c8306ec41719e241e641b132`
- `mutant-request-hash`: `7d35989982688149b6da3572049de04faae1bd1c9d8216aa07bca49b2c34ee5a`
- `mutant-candidate-hash`: `a9a19d8b9a1c1f93f4488ac8898ac42d02d3ad825d72cf6ea042eee0fc3cb567`
- `mutant-outcome`: `68a70087b352ca6901b958a3d1a27ce4aa8a3f706c6e4cc8bb644964daf29cc1`
- `mutant-record`: `b999141d8e5df4f32549a52ded85030a25f92792ffe33e82f8f1f49704a6959d`
- `mutant-category`: `7fefef9fb36d67687efb0bc25794fd4e001ebcef0374fcb62cb2fc0cda0b7941`
- `mutant-lane`: `afd4a753c7ea049438f67301fd007642d402cd67789c821cdf6a792afd2cd7f5`
- `mutant-duplicate`: `2392adbec59a0dec2b94da1f6555d34b41e917884d8319a1d28bda1ab72bdccd`
- `mutant-caps`: `806f97312fa2b8376872d0678b0878ee5499510e07a15b88b8bd13b608ce097d`
- `mutant-anchors`: `8aeb8493b8f49415a764e0a217612d30e57e073e828b8a5a69c9dc3eb5845475`
- `mutant-peers`: `25aadfe56dca600e481a58e27e977cf5b2d774ace53d370932e1a4af837b0b38`


## Retrieval and authorized environment

Serena initialized for canonical project; archive source is unindexed and was
read directly. Parent Tier2 graph evidence is from2026-10-04 primary checkout,
not watched, best-effort metadata_match for packet/types/model; current exact
reference hashes were checked above. No graph completeness claim. Read relevant
instructions, TDD skill, accepted scope, packet, prior oracle/tests and normative
model context. Authorized environment is only owned inert byte fixtures and
ordinary offline compiler/unit tests. No credentials, external data or system
were targeted. Return control to Brain for review/integration; no owned process
session remains active.
