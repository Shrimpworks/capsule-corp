# G3a10 implementation handoff

PASSED — bounded implementation and local verification only. Independent review
IN_PROGRESS — TRENDING_GOOD (pending fresh reviewer); full Gate3 IN_PROGRESS —
TRENDING_GOOD; child execution, Gates4–6 and product admission BLOCKED. This author
provides no independent readiness verdict. Brain owns publication/integration.

## Neutral freeze identity

Archive /tmp/capsule-c5b19-recovery, branch codex/c5b19-benign-teardown, unchanged
HEAD/base7646d8faabb670ce8b5ff1b213c414d91deda2d9. Exactly four owned files changed
under experiments/typed-guest-transport-c5b19-benign-teardown; no commit/push/PR
created by this worker. Shared draft PR38 is the parent-owned destination.
Canonical basefe9009d5398a883735219e2a3a33c351c79aa54c, branch
codex/c5b19-worker-sidecars. No canonical file was changed by this worker.

- tests/oracle.go: `df113eff8793f07a2156984aebdf69e3dc77deda39dbdc1c80573ec8c6e392ef`
- tests/oracle_test.go: `05197b095ed3ed7da6ff9562b1b6e6da34a2624bf86821c138623a6f12a9448c`
- PLAN.md: `fbccc799aaa324a38348624bde314e61555dcf7bb172c24383a9450a92636085`
- inputs/ORIGINS.json: `5f3a8ee2d1e15eaac750a62511587e62342172270102e2b9da5c80726be7b1b6`

Full four-file patch `/tmp/c5b19-g3a10-freeze.patch`: `59479399ae70aaec90bd3a15c6304151ca651f2859286bb4cf3a65efa4eb8642`.
Evidence `/tmp/c5b19-g3a10-evidence.json`; accepted scope
`/tmp/c5b19-g3a10-scope.md`. ORIGINS.g3a10 retains exact mutation definitions,
commands, assertion locations/counts, log hashes and invalid-construction records.
PLAN.md retains the contract, limits and verification ledger. Existing oracle.go
and oracle_test.go prefixes remain byte-exact; go.mod/model_test.go and all earlier
parsed ORIGINS values, including false run/admission flags, remain unchanged.

## Author explanation: intent and traceability

The prior capturedFramePair Request for kind5/7 was explicitly decoding context,
not evidence that the worker received bytes. This slice adds distinct independent
raw receive/send inputs and checks their correspondence with actual privately
retained submissions and represented reply observations. Packet350–359 requires
worker captures separately from control captures. Packet432–435 prevents worker
clocks from reordering control. ADR47 preserves unknown outcomes and exclusive
Supervisor authority. This implementation supplies neither process effects nor
settlement authority. The accepted scope report freezes optionality and boundaries.

Defensive assets: owned synthetic397/98/256-byte fixtures, value histories and
ordinary offline compiler/unit-test processes in the archive. No native fixture,
child/guest, lifecycle/signal/wait target, live clock/store/network, credential or
user-data access. Ten new test groups and98 appended implementation lines; no
new import, dependency, primitive, packet format, manifest, ABI or product consumer.

## Flow and invariants

`inspectWorkerSidecars` first requires initialized snapshotHistory and ≤3 inputs.
For each item it strictly decodes required WorkerReceived, matches the entire
397-byte request to an existing private submission, and rejects duplicate operation.
Nil alone omits WorkerSent; any nonnil value must have a kind5 pair. Optional kind7
requires that pair. Supported rows are receive only, receive+sent+kind5, and that
row plus kind7. Kind5 and optional kind7 pass unchanged strict capture validators;
their entire request/reply byte arrays must equal independent raw slots. Any
nonzero accepted outcome also requires sent bytes equal the complete retained
accepted reply. Results contain copied fixed value facts and current history;
errors return all zero with no partial collection. The receiver never changes.

Received presence, sent/control presence and worker/accepted outcomes stay
separate. Accepted state never fills omitted worker observations. Confirmed worker
or kind7 observations never apply a reply, settle pending, clear fatal, change a
projection, grant start/terminal or establish physical transport. Old observations
and arbitrary selection ordering are accepted as correspondence only; no current
watermark or lane ordering is consulted. Strict single-event lane/timebase checks
still apply. G3a9 absent/older snapshot semantics and prior APIs remain unchanged.

## Alternatives and costs

Reuse snapshotHistory's three private entries rather than introduce a parallel
ledger/wrapper; append a read-only method rather than alter prior semantics. A
fixed three-entry result makes ownership and all-zero error output explicit;
received-presence is redundant but explicit. Nil-only send omission distinguishes
missing data from malformed allocated-empty data. In-memory selection structs
avoid introducing a new serialized sidecar format or worker-receive event kind.

At most three requests, three membership candidates per request, and three prior
operation comparisons; bounded strict frame parses/SHA256 operations reuse existing
helpers. No work or allocation proportional to unchecked byte lengths: decoders
refuse wrong frame sizes before processing, comparisons use bounded retained bytes,
and result storage is fixed. Caller-supplied slices/pointers are not retained.
Concurrent mutation by the caller remains outside this sequential test-only API.

## Verification and mutation evidence

Tests preceded the query and compiler-red identified missing API names; the minimal
query then passed. Refinements separated valid independently substituted pairs from
malformed evidence, added historical/rejected/mixed-row and input-refusal checks,
and kept prior target lines stable. Exact installed executable:
`/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go`
reports Go1.25.13 darwin/arm64. All commands cwd the experiment with
`GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off GOCACHE=/tmp/c5b19-g3a10-go-cache`.

- red: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./...`, exit1, log `/tmp/c5b19-g3a10-red.log`, SHA256 `b88c066e49a61c1ce659e872fae2f17b411c3b2cd4901c630185c0f2d8f3ad5d`.
- green: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./...`, exit0, log `/tmp/c5b19-g3a10-green.log`, SHA256 `e5bd60ef196d77fdecbef9ce3131b7f9c2fb3bfbace7055c2450957a9cbeb54c`.
- refined-green: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./... -count=1`, exit0, log `/tmp/c5b19-g3a10-refined-green.log`, SHA256 `f5e624fcb40238f144faf7ae9f2bcbaa353cfc68e5e5f31f3c70a757bdc56994`.
- mutation-baseline: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./... -count=1`, exit0, log `/tmp/c5b19-g3a10-mutation-baseline.log`, SHA256 `ca273146ad5426302e32e046919ac582694de81ce83a3dfad7e8d4f9cae53b08`.
- restored: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./... -count=1`, exit0, log `/tmp/c5b19-g3a10-restored.log`, SHA256 `28471f8304d1de94427a5e58e0eca88b9e9303e6c442391ca07ee88997d4b3a2`.
- race: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test -race ./... -count=1`, exit0, log `/tmp/c5b19-g3a10-race.log`, SHA256 `12373c05568dd4c5611f25d202f093d726175770476b3499213dc7dbaa5962a8`.
- vet: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go vet ./...`, exit0, log `/tmp/c5b19-g3a10-vet.log`, SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.
- build: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go build ./...`, exit0, log `/tmp/c5b19-g3a10-build.log`, SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.
- final-test: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test ./... -count=1`, exit0, log `/tmp/c5b19-g3a10-final-test.log`, SHA256 `068b8ae823e4e6f5577b8f1557c50233869c22bc57fb28c32238f5167b108a73`.
- final-race: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test -race ./... -count=1`, exit0, log `/tmp/c5b19-g3a10-final-race.log`, SHA256 `90eea5e2916f45bb0b006dfc760b0365aff1dda3cfee6ff903832238ee8d930b`.
- final-vet: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go vet ./...`, exit0, log `/tmp/c5b19-g3a10-final-vet.log`, SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.
- final-build: `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go build ./...`, exit0, log `/tmp/c5b19-g3a10-final-build.log`, SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.

Fourteen corrected selected mutants compiled and exited1 with17 exact intended
assertion occurrences after valid setup. Each source replacement, selected test
command, assertion and count is retained in ORIGINS/evidence. `finally` restored
exact source bytes after every attempt and campaign stop. Source restoration hash
is the oracle.go freeze hash above. `git diff --check` passed.

The control request/reply checks overlap through validated request digest. Their
disclosed combined bypass demonstrates the relationship, not isolated necessity
of both guards. The cap is likewise reinforced by private three-entry membership
and duplicate refusal; no cap-only mutant was claimed. Invalid constructions are
explicitly separated: deleting control-kind removed its only variable use and
failed compilation; overly broad optionality bypass reached its intended assertion
then panicked in a parent-level malformed-member case. Neither run is counted.
Corrected mutations retain variable use/nil safety and reach exact intended
assertions without compilation errors or panic. No survivor/setup failure is
reported as a kill, and this is not exhaustive mutation proof. Final explicit
refusal-input test appended afterward leaves prior mutant target bytes/line numbers
unchanged; all final checks passed on the complete final tests.

## Review challenge points, limits and deferred work

Reviewer should challenge independent slots versus context, nil/empty member
semantics, retained membership and accepted-reply equality for outcomes2/3,
historical/rejected evidence behavior, optionality paths, whole-zero late failure,
value ownership, and the exact mutation target counts after valid setup. The tests
use existing independent packet fixtures rather than query-generated fixtures.
Review should separately assess what redundant guards and selected mutants cannot
prove; no broad coverage claim follows from fourteen examples.

Supplied correspondence cannot establish authenticity/completeness, omitted frame
absence, actual syscall publication, completed roundtrip, full lane/causal/timing
consistency, storage durability, native absence or execution authority. Partial
receive-only and sent-without-control remain unknown beyond represented facts.
Next recommended bounded slice is supplied kind3→4/kind5→6 write/result
correspondence; it is expressly deferred here along with full trace composition,
native code and all first-run/admission gates. This worker changed only the four
assigned archive files and created no branch, commit, push or PR.

Retrieval: Serena initialized session8ca2ac4b. Tier2 primary graph generation
2026-10-04T22:18:10Z is old/not watched; own coverage reports both archive sources
outside_project. Exact source reads and supplied parent SettleWrite evidence ground
this bounded composition; no graph completeness/freshness claim is made.
