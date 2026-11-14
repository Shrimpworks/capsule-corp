# G3a11 implementation handoff

PASSED — bounded implementation and local verification only. Independent review
IN_PROGRESS — TRENDING_GOOD (pending); parent full Gate3 IN_PROGRESS — TRENDING_GOOD.
Child execution, Gates4–6 and product admission BLOCKED. No author readiness verdict.

## Neutral freeze identity

Archive /tmp/capsule-c5b19-recovery, branch codex/c5b19-benign-teardown, unchanged
HEAD/base d698bd9332a000ad3c2b9f23f87d0353dc5964b4. Exactly four assigned files changed
under experiments/typed-guest-transport-c5b19-benign-teardown. Canonical base
 d042e53a877e494bd2eb0e325ec084d4b9cc692a, branch codex/c5b19-write-results.
No canonical change, branch, commit, push or PR was made by this worker. Brain owns
shared archive draft PR38 and canonical PR publication.

- tests/oracle.go: `4b9012aa8e19b5b35995e5fb48a7e6fd8cd2f858ac7a92ea2c4678e85ce6bc69`
- tests/oracle_test.go: `acb23007854fc92f8ddbffa7daea3c1bfeaa172ea85e88644fc7a0ef2f51975c`
- PLAN.md: `c0f185ed7548cbb6fa0aabcbb475fd7fd2652498fae4ba7ceb6350a9447b989e`
- inputs/ORIGINS.json: `086c03605ff2f240450c5bd94ea8256a59ad8316788a5e2b0ad9edca57a8d968`

Neutral four-file patch `/tmp/c5b19-g3a11-freeze.patch`, SHA256 `e8ea12e51b7185fb025842eab087f2ff9b046a774ac79734e4cccbbaa2d977ed`.
Evidence `/tmp/c5b19-g3a11-evidence.json`; accepted scope
`/tmp/c5b19-g3a11-scope.md`; mutation driver `/tmp/c5b19-g3a11-mutate.py`.
ORIGINS.g3a11 retains all exact selected replacements/commands/assertions/counts/log
hashes. Prior oracle/test byte prefixes, go.mod/model_test.go and all historical
parsed ORIGINS values, including false run/admission flags, remain unchanged.

## 1. Intent

Close the bounded gap between represented pre-write captured frames and their
returned result events: kind3→4 request write and kind5→6 worker reply write.
This is diagnostic byte correspondence and local pair ordering only, not physical
transport completion or successful-syscall classification.

## 2. Requirements and source traceability

Packet187–214 freezes397-byte request and98-byte reply writes, with no retries for
short/zero/error/excess. Packet300–381 defines signed result/errno slots and strict
kind fields; packet432–435 forbids worker clocks reordering control. Existing
strict helpers already validate frame identity, event fields, lanes and context,
while preserving raw signed diagnostics. Neither packet nor current archive
PLAN543–544 defines errno reset/normalization on success. Brain explicitly accepted
ExactFrameCount as numerical result equality independent of errno before tests.
ADR47 authority and unknown-state boundaries remain unchanged. The canonical plan,
accepted scope and next-scope identify this as a separate bounded slice.

## 3. Defensive scope and authorization

Only owned inert synthetic request/reply/event bytes and ordinary offline Go1.25.13
compiler/unit processes in this archive. No native fixture, child/guest, lifecycle,
signal/wait target, live clock/store/network, credentials or user data. Brain
accepted scope in the G3a11 assignment before edits/tests. Tests-first skill reused.
Serena session8ca2ac4b remains initialized. Tier2 primary graph generation
2026-10-04T22:18:10Z is old/unwatched; own index_status/coverage reports packet
metadata_match and archive outside_project. Exact current source reads supply
archive evidence; no graph freshness/completeness claim.

## 4. Implemented data flow

`inspectWriteResults` first validates frozen bindings, nonzero timebase and cap6.
Each captured pre-event and result must pass existing strict validators and the
closed3→4/5→6 pairing. Compare operation/generation/request/candidate and
outcome/recorddigest; require increasing sequence and nondecreasing tick. Count
at most3/category and3 unique operations, reject duplicate category/operation,
and require exact full397 request equality across categories for a shared ID.
Decode the validated request into copied operation facts, retain raw signed result
and errno, derive only category-specific ExactFrameCount, and append to fixed
six-entry value output. Every failure returns the completely zero result.

## 5. Decisions and alternatives

A stateless function composes existing validators; no history wrapper or redundant
submission ledger. Fixed value output expresses cap/ownership and avoids partial
prefix publication. Diagnostic acceptance preserves every signed i32 result/errno
rather than guessing a platform errno whitelist or normalization. ExactFrameCount
names only represented equality; a FullWrite/Success/TransportPassed flag would
invent unsupported producer semantics and is deliberately absent. Local direct
sequence/tick comparisons permit gaps/equality of ticks without overflow. No
inter-pair ordering or sequence uniqueness is added: full lane consistency is the
next separately bounded recommendation.

## 6. Invariants

Neither category requires the other; omitted whole pairs remain unknown. A selected
missing/malformed member refuses. Empty selection validates context and yields no
evidence. Raw signed diagnostics are neither normalized nor treated as a realizable
syscall claim; exact count with negative/nonzero errno remains exact-count data.
Kind5 Request is decoding context, not independent worker receipt. No begin,
reconcile, applyReply or settlement call exists; history/durable projection cannot
change. Caller bytes and frozen value context never mutate, outputs contain no
slices/pointers, input order is preserved, and late failure leaks no partial facts.
All prior APIs/source/test prefixes/provenance are preserved.

## 7. Verification and experimental evidence

Eight test groups were appended first; compiler-red named the missing new pair,
result and query APIs. The minimal89-line implementation then passed. A ninth test
group refined full copied-operation fact checks before mutation baseline.335 new
test lines cover all accepted scope categories, signed boundary samples,
well-formed substitutions, strict malformed context/member cases, local chronology
zero/max/equal/gap/regression, duplicates/caps/cross-category identity, arbitrary
selection order, deliberately repeated inter-pair indices, raw facts and ownership.

Exact Go binary `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go`
reports Go1.25.13 darwin/arm64. Cwd the experiment, environment
`GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/c5b19-g3a11-go-cache`.

- red: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./...`, exit1; log `/tmp/c5b19-g3a11-red.log`, SHA256 `b185a38514ca77801877dd96673a34c6d86da09205bae559be7514768f2f8ab5`.
- green: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./...`, exit0; log `/tmp/c5b19-g3a11-green.log`, SHA256 `0b3a4e17c1e8ac8c460a69bf285084d9bc088314f13d1e74a74245e6b6937960`.
- mutation-baseline: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./... -count=1`, exit0; log `/tmp/c5b19-g3a11-mutation-baseline.log`, SHA256 `b26a014b4ef8d8fb0aa9c6c050a405f349404d17cfa6d3745d83d5594ade6765`.
- restored: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./... -count=1`, exit0; log `/tmp/c5b19-g3a11-restored.log`, SHA256 `adec114e8281fcda240a9dab64468b0f35a0917f8ef1f49c8035b0d2e0b314af`.
- race: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test -race ./... -count=1`, exit0; log `/tmp/c5b19-g3a11-race.log`, SHA256 `517e2bf60ac2bf212c1861377f4f20bee5e6966fd67356af5dfa568f0bf6e160`.
- vet: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go vet ./...`, exit0; log `/tmp/c5b19-g3a11-vet.log`, SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.
- build: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go build ./...`, exit0; log `/tmp/c5b19-g3a11-build.log`, SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.

Eighteen selected mutants compiled and exited1 at27 exact intended assertion
occurrences after valid setup. Result-kind closure; six identity/outcome/digest
fields; sequence/tick guards; duplicate/unique/cross-category request guards;
category-specific count; errno-independent count; raw errno/result preservation;
empty-context validation; and all-zero late-error behavior are covered. All exact
replacements/commands/assertions/counts and hashes are in ORIGINS/evidence. No
construction failure, panic, setup failure or survivor occurred. `finally` restored
exact source bytes after every attempt and final standard test/race/vet/build passed.
`git diff --check` passed. This is selected sensitivity, not exhaustive proof.

## 8. Cost and operational impact

At most6 input pairs, constant-size256/397/98 validation and hashes, six-slot result,
fixed event/count arrays and bounded pairwise scans. Wrong-sized frames refuse
before unbounded processing. No new import/dependency, primitive, serialized
format, manifest/ABI, privilege, product consumer or live effect. Standard-library
parsing/comparison reuse follows C5b19 test-only ADOPT-PLATFORM / BUILD-NARROWLY.
Concurrent caller mutation remains outside the sequential test-only API contract.

## 9. Risks, limits and deferred work

Numerical ExactFrameCount is not transport progress: actual producer errno capture
semantics and real syscall success remain unspecified/unvalidated here. Raw
contradictory diagnostics are preserved, not certified realizable. No capture
truth/completeness, physical transfer, independent worker receipt, control
application, durability, full roundtrip, inter-pair/full-stream chronology, timing,
start/stop, native absence or authority claim. Cap guards overlap with unique and
duplicate checks; no isolated redundant-cap mutant necessity is claimed. Next
transport-lane consistency across selected kinds3–8 is deferred, along with full
trace/control/native and first-run gates. Child execution stays BLOCKED.

## 10. Independent review challenge points

Challenge raw signed result/errno fidelity and the deliberately errno-independent
count label; strict closed pair kinds and identity comparisons; pair-local direct
chronology versus omitted/interleaved events; same-operation cross-category exact
request equality; arbitrary selection order; caps and duplicates; all-zero late
failure; and copied-value/input invariants. Inspect each selected mutation's exact
assertion after valid setup, distinguishing finite sensitivity from exhaustive
proof or isolated redundant-guard necessity. Confirm no interpretation strengthens
prior G3a10 correspondence into full lane chronology or physical transport evidence.
Brain owns independent review, final acceptance, immutable pinning and publication.
