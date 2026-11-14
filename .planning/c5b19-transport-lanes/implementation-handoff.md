# G3a12 implementation handoff

PASSED — bounded implementation/local verification. Independent review IN_PROGRESS —
TRENDING_GOOD (pending); parent full Gate3 IN_PROGRESS — TRENDING_GOOD. Child execution,
Gates4–6 and product admission BLOCKED. Author provides no readiness verdict.

## Neutral freeze identity

Archive /tmp/capsule-c5b19-recovery, branch codex/c5b19-benign-teardown, unchanged
HEAD/base7ec433aa17c623f58a81118bab87f1ada49770f3. Exactly four assigned files changed
under experiments/typed-guest-transport-c5b19-benign-teardown. Canonical branch
codex/c5b19-transport-lanes at f1ded569af4a6ae20f2a93b7e476825959f97143. No canonical
file, branch, commit, push or PR changed by this worker. Brain owns shared draft
archive PR38 and canonical publication.

- tests/oracle.go: `22e91fd0c95e778f9e8e68480bcb586fb4fe64d4d3c0ebb9b0c56b887b7ac0be`
- tests/oracle_test.go: `8c10cdab499b5e0e9260848548adae83fc788b9f1c26defae054dd06c9cf5b2c`
- PLAN.md: `5d1b65fbd1c292ea0fe74782670ed8993e334db9ec2b020702471b735d6dc645`
- inputs/ORIGINS.json: `8192daa4bb347b469a9f070dc5a8ebcc2c4581a658e71dc0dc71d64e373971a1`

Neutral patch `/tmp/c5b19-g3a12-freeze.patch`: `b1523a16e6d47b8538218a966604cdc6849167c0da8241a59301cfc52cd8c2fb`.
Evidence `/tmp/c5b19-g3a12-evidence.json`; accepted scope
`/tmp/c5b19-g3a12-scope.md`; mutation driver `/tmp/c5b19-g3a12-mutate.py`.
Previous source/test prefixes, go.mod/model_test.go, parsed historical ORIGINS
values and false run/admission flags remain exact. Mutation definitions, commands,
assertion locations/counts and log hashes are retained in ORIGINS.g3a12 and evidence.

## 1. Intent

Reject sequence reuse and tick regression between supplied strict transport events
in the same lane, which can remain invisible to valid individual write/result or
received/applied pairs. This closes only a partial-selection lane contradiction;
it does not produce a complete trace or correspondence verdict.

## 2. Requirements and traceability

Packet300–335 defines lane-local sequence slots and closed event kinds;395–410
separates writer arrays and control watermark;429–435 forbids total cross-lane
ordering. Existing G3a11 explicitly validates local pairs without inter-pair
chronology. G3a7 validates received/applied pairs separately. Brain accepted the
error-only kind3–8 contract before tests and the written scope reconciled without
conflict. Prior ADR47 exclusive process authority and test-only reuse-row boundaries
remain unchanged. Canonical plan and prior next-scope identify this finite gap.

## 3. Defensive scope and authorization

Owned inert synthetic frame/event arrays and ordinary offline Go1.25.13 compiler/
unit tests only. No native fixture, child/guest, lifecycle/signal/wait target,
live store/clock/network, credential or user-data access. No dependency, primitive,
serialization, manifest, ABI, privilege or product consumer. Existing TDD skill
instructions apply. Serena session8ca2ac4b remains initialized. Own index_status
and coverage show old2026-10-04 primary graph, not watched; packet metadata_match
and archive outside_project. Exact current source reads supply bounded evidence;
no graph freshness/completeness claim.

## 4. Implemented flow

`validateTransportLanes` validates bindings, nonzero timebase and max18 before
processing; each input strictly decodes under one context and must be kind3–8.
It caps3 per kind, counts distinct operationIDs, and rejects duplicate kind/operation
regardless of payload differences. For each prior same-lane event it refuses equal
sequence or a tick inconsistent with sequence order using direct unsigned
comparisons. Different lanes bypass chronology comparisons. Fixed local arrays
retain decoded values during the call only. Success returns nil, refusal an error;
no result object, validated prefix or mutable history is exposed.

## 5. Decisions and alternatives

A separate error-only validator avoids a new facts wrapper or rewriting prior
APIs. Raw event selection keeps frame/tuple/pair/history obligations explicit;
calling those validators implicitly would overstate this slice. Bounded pairwise
comparison is simpler than sorting and leaves caller outer order and slice identity
unchanged. Gaps and nonzero origins are valid because this is partial selection;
requiring adjacency would confuse it with the eventual complete-lane verifier.
Independent lane clocks prohibit a global sequence/timestamp watermark. Direct
comparison avoids subtraction or signed-cast overflow.

## 6. Invariants

Valid empty context plus empty selection is no evidence; invalid context refuses
even empty. Each event is strictly shaped with proper lane; no harness/nontransport
kind admitted. Same-lane sequences unique, ticks nondecreasing in sequence order,
regardless of input order; equal ticks allowed. Other lanes may reuse sequence and
have arbitrary relative ticks. Duplicate identity is exactly kind+operationID,
not generation/digest. Different kinds sharing an operation may have inconsistent
tuple fields here; separate correspondence is still mandatory. Lone application
or result events are allowed partial selections. No implicit pair/order/application
proof, state update, settlement or authority. Inputs/context and outer slice order/
inner identity never mutate; no bytes/state retained across calls.

## 7. Verification

Tests were appended first; compiler-red named missing validateTransportLanes.
Written-scope refinements preceded the53-line implementation, which passed green.
Ten groups/341 appended test lines cover all supported kinds, partial/empty/context,
strict malformed/unsupported domain, locally valid write/result and received/applied
contradictions, cross-kind boundary reuse, control/worker forward/reverse/interleaved
regression,24 four-event permutations, full18 deterministic permutation, independent
clocks, equal ticks/gaps, zero/max widths, caps/duplicates, positive layer separation
and outer order/inner identity/byte preservation on success and late failure.

Exact executable `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go`
reports Go1.25.13 darwin/arm64. Cwd the experiment; environment
`GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/c5b19-g3a12-go-cache`.

- red: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./...`, exit1; log `/tmp/c5b19-g3a12-red.log`, SHA256 `1f6176723bf4ead02bb4db7407f79d5b3b9d507b1b3ed9da06314a7097af2d45`.
- green: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./...`, exit0; log `/tmp/c5b19-g3a12-green.log`, SHA256 `e5bd60ef196d77fdecbef9ce3131b7f9c2fb3bfbace7055c2450957a9cbeb54c`.
- mutation-baseline: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./... -count=1`, exit0; log `/tmp/c5b19-g3a12-mutation-baseline.log`, SHA256 `bb187a81052ac031864853d38f58f81b9ff6114da1ecc06e8aad39c63f4bf2e7`.
- restored: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./... -count=1`, exit0; log `/tmp/c5b19-g3a12-restored.log`, SHA256 `7c7510df158490436c19b19d491c83634c5d6ca1eaf0bb290bc90322dedd080f`.
- race: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test -race ./... -count=1`, exit0; log `/tmp/c5b19-g3a12-race.log`, SHA256 `f80ae054612f80c4ffcd0f8d3ea26685cda6d5ba4f9d5d124b9afb755993b853`.
- vet: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go vet ./...`, exit0; log `/tmp/c5b19-g3a12-vet.log`, SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.
- build: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go build ./...`, exit0; log `/tmp/c5b19-g3a12-build.log`, SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.
- final-test: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./... -count=1`, exit0; log `/tmp/c5b19-g3a12-final-test.log`, SHA256 `8aa6efdf3a3c9270a2dd274117fdb9f3a1b5d56698c6278841a406615697b79f`.
- final-race: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test -race ./... -count=1`, exit0; log `/tmp/c5b19-g3a12-final-race.log`, SHA256 `ae26aa7403c36406625080b2a24629e34629b9b868f6508d8747c71435235395`.
- final-vet: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go vet ./...`, exit0; log `/tmp/c5b19-g3a12-final-vet.log`, SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.
- final-build: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go build ./...`, exit0; log `/tmp/c5b19-g3a12-final-build.log`, SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.

Twelve selected mutants compiled and exited1 with17 exact intended assertions after
valid setup. They cover safe supported-kind lower-bound removal, duplicate identity,
same-lane sequence equality, both input-order tick branches, lane filtering, allowed
equal ticks/gaps, unique operations, unintended tuple correspondence, empty context
and caller-slice reordering. No invalid construction, panic, setup failure or survivor.
Each exact replacement/command/assertion/count/log hash is retained. Finally restored
original source after each attempt; final offline full/race/vet/build all passed.
The final full18 permutation/raw-diagnostic positive was appended after the campaign
without changing target bytes or lines, followed by final checks. `git diff --check`
passed. Total/per-kind caps are logically reinforced by operation/duplicate caps;
no independent redundant-cap necessity or exhaustive mutation coverage claimed.

## 8. Costs and operational impact

Maximum18 events, fixed18 decoded records and9 counters, at most153 prior-event
comparisons plus constant-size strict256-byte parsing/hashing. Wrong frame lengths
refuse via existing decoder. No input sorting, caller mutation, heap-sized result,
new I/O or process authority. Existing frame/history APIs remain untouched. As with
other passive selectors, concurrent caller mutation is outside this sequential
fixture API contract. Product code imports nothing from the archive.

## 9. Risks, limits and deferred work

Success cannot prove completeness, contiguous whole-lane indices, capture truth,
raw-frame or operation-tuple correspondence, pair presence/semantic order,
control-history decisions, physical effects, syscall success, durability, full
roundtrip, timing, start/stop, absence or admission. Raw signed diagnostics remain
uninterpreted. Removing only safe lower kind bound is selected mutant evidence;
no unsafe high-kind indexing bypass is run or claimed. Next same-operation control
ordering is explicitly deferred, as are full trace/native and first-run gates.
Child execution remains BLOCKED.

## 10. Independent review challenge points

Challenge the distinction between partial lane consistency and full trace validity;
operation dedupe versus tuple agreement; both sequence/tick comparison directions;
no cross-lane constraints; correct zero/max behavior; positive gaps/equal ticks;
strict kind bounds before counter indexing; separate composition and lone events;
outer/inner slice preservation; locally valid setups for each inter-pair defect;
and exact mutation assertion counts rather than generic nonzero exits. Assess
limits of finite mutations and overlapping caps separately from tested behavior.
Brain owns fresh review, acceptance, immutable pinning, commits and publication.
