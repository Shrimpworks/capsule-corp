# G3a9 author implementation handoff

Status: PASSED for bounded implementation and local verification. Independent
review IN_PROGRESS — TRENDING_GOOD, pending instance1 of3. Full Gate3 remains
IN_PROGRESS — TRENDING_GOOD; child execution, Gates4–6 and admission remain BLOCKED.
This is author testimony, not an independent readiness verdict.

## Intent and exact scope

Defensively classify supplied durable-record/request/reply correspondence using
owned inert synthetic byte fixtures. A visible record, accepted control outcome
and optional worker outcome remain distinct facts. No final-store consistency,
physical supersession, observed fsync or execution authority is inferred.

Archive /tmp/capsule-c5b19-recovery, branch codex/c5b19-benign-teardown,
HEAD/base febb88d6421058f048f75319baae7cf9bd5a0cbb. Only four assigned paths changed:
experiments/typed-guest-transport-c5b19-benign-teardown/{tests/oracle.go,
tests/oracle_test.go,PLAN.md,inputs/ORIGINS.json}. Canonical reference is
f3f2e0755c0c284dd28019225bef0cf8482d3cdd in the parent recovery worktree and branch
codex/c5b19-snapshot-correspondence. Brain owns existing draft archive PR38 and the
new canonical PR. No worker branch, commit, push, PR or subagent was created.
All other changes were preserved; no owned process remains active.

## Plan and traceability

Brain explicitly accepted /tmp/c5b19-g3a9-scope.md before edits or tests. The
canonical scope and delivery plan retain that contract. Frozen Gate2 packet253–271
forbids record-presence confirmation;342–359 requires exact control/worker/snapshot
correspondence. Existing G3a8 trial state is composed unchanged. Every requested
input/result distinction is implemented; full causal/capture/control semantics
remain deferred as accepted. No behavior scope deviation was needed.

## Technical approach

`snapshotHistory` owns value-copied `decisionHistory`, three fixed request/reply/
outcome entries, initialization and count. Constructor delegates frozen-context
validation. `begin` executes G3a8 on trial state, copies the successfully recorded
request and commits once. `reconcile` similarly delegates to trial G3a8; code1
stores exact reply/outcome only against an existing equal request, while code2
commits only the unchanged evidence plus modeled G3a8 fatal/decision effects.
Refusals discard the entire trial.

`inspect` is read-only. Absence requires zero record bytes and no worker, returns
zero match fields plus copied current history, and erases nothing. Present bytes
strictly decode and must match a privately retained submission. Accepted outcome1
matches the supplied record hash to retained reply digest. Optional worker is
strict kind5, exact private request, valid reply and, when accepted evidence exists,
byte-identical accepted reply. Its outcome1 digest also matches supplied bytes.
Outcomes2/3 remain valid with a present candidate; their zero reply digests never
promote visibility to confirmation. Missing worker stays unknown. `LaterSubmitted`
means only a later retained submission exists, including pending/failed/unknown;
older candidates keep their own outcome and cannot rewrite the current projection.
There is no full-confirmation, admission or case-verdict boolean.

## Changed-component walkthrough

- oracle.go:138 appended lines for fixed evidence slots, wrapper transitions and
  read-only copied correspondence result.
- oracle_test.go:472 appended lines; ten test functions and small test helpers.
- PLAN.md:scope, verification commands, intended-assertion mutation table, hashes
  and limitations appended; prior prefix remains exact.
- ORIGINS.json:new g3a9 object with source/reference/log hashes, exact selected
  mutation replacements/commands/locations/counts and pending review status.
  Every historical parsed value and all run/admission flags remain unchanged.

## Decisions and alternatives

Retain fixed private slots rather than attach a borrowed G3a8 history or trust a
caller-provided request at query time; G3a8 itself retains only pending/last-settled
information, insufficient for earlier submitted candidates. Compose existing
strict decoders rather than add serialization or a storage path. Return separate
outcomes instead of an ambiguous confirmation boolean. Preserve explicit absence
as supplied observation, and order-only LaterSubmitted instead of inventing a
snapshot timestamp or physical replacement claim. No dependency/new primitive;
C5b19 ADOPT-PLATFORM / BUILD-NARROWLY stays test-only.

## Invariants and boundary cases

Whole-value ownership preserves atomicity and no-alias behavior. Earlier APIs and
source bytes remain exact; inherited one-pending/generation/unique-ID/kind/terminal,
control watermark and three-decision rules still apply. Successful/rejected queries
never consume slots, decisions or watermarks. A rejected anchor is not submission;
code2 cannot overwrite accepted evidence; a late settlement and all queries retain
fatal. Newer visible uncertain candidates coexist with earlier confirmed projection;
older or absent observations neither roll it back nor imply physical consistency.

## Verification and observed results

Exact offline Go1.25.13 binary with GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off and
GOCACHE=/tmp/c5b19-g3a9-go-cache, cwd archive experiment:

- Tests-first `go test ./...`: expected compiler-red, undefined new API.
- Initial implementation green; fixture coverage and focused subtest-label
  refinements green. No unexpected failing behavior or malformed-fixture correction.
- Mutation baseline `go test ./... -count=1`: passed.
- Twelve compiled selected mutants: each exited1 and reached its exact intended
  assertion count (membership4, worker-request8, all other targets1). No survivor,
  panic, setup-only credit or construction failure. Exact source restored in finally.
- Restored `go test ./... -count=1`: passed.
- `go test -race -count=1 ./...`: passed,2.684s.
- `go vet ./...`, `go build ./...`: exit0, no diagnostics.
- Prior oracle/test/PLAN prefix comparison, exact go.mod/model_test.go identity,
  all historical parsed ORIGINS values and `git diff --check`: passed.

Tests cover explicit absence/context; pending worker-only outcomes; accepted
confirmed/failed/indeterminate and worker contradictions; current/older candidates
across three-write/direct-terminal histories; independent valid unsubmitted and
foreign worker fixtures; strict malformed kind/lane/binding/length/digest evidence;
all395 record-byte substitutions; code2/late-fatal preservation; complete state
nonmutation on query/refusal; caller and return-value ownership; inherited bounds.

Exact selected mutation definitions, commands, intended assertion strings/counts
and hashes are retained in ORIGINS.g3a9.selected_mutations, summarized in PLAN.
The membership test uses a disclosed combined membership/accepted-digest bypass:
those guards mathematically overlap for canonical inputs. It preserves valid setup
and reaches all four unsubmitted-record substitutions. Worker-kind uses validkind7
with complete reply; worker-request uses eight valid uncertain foreign requests;
worker/accepted equality uses failed versus indeterminate zero-digest replies.
No isolated digest/capacity necessity or exhaustive mutation proof is claimed.

Raw /tmp/c5b19-g3a9-* logs and script are ephemeral. /tmp/c5b19-g3a9-evidence.json
and mutations.json summarize exact command execution, hashes and failure excerpts;
retained PLAN/ORIGINS carry the durable evidence identities/definitions. Canonical
verification and hosted checks are Brain-owned, not claimed as executed here.

## Costs, limits and unresolved questions

Constant-size value copies and at most three repeated strict record/request decodes
are deliberately simple. No I/O or unbounded processing. Full snapshot timing,
final-store contradictions, complete traces/worker-received mapping, G3a8 interleaved
actions, malformed diagnostic mapping, raw process-derived projections, actual
control admission, fsync/rename/custody, restart recovery, deadline/lifecycle and
native execution remain outside this slice. No unresolved requirement blocks this
accepted bounded contract; confidence is high within inert correspondence only.

## Retrieval and authorized environment

Serena initialized (session67ee4989). Codebase-memory Tier2 primary index ready at
2026-10-04T22:18:10Z, not watched and not this archive. Parent supplied complete
SettleWrite four-callee trace and exact projection/pending semantics. This worker
checked index status and coverage: starter docs/packet/ADR best-effort metadata_match,
archive source/test outside_project. Direct exact archive and canonical reads
supply fallback; no whole-graph completeness claim. Applicable archive AGENTS,
TDD skill, starter authority boundaries, ADR47, frozen packet references and scope
were consulted. Only ordinary offline compiler/unit processes and inert bytes were
used; no child/native fixture/guest, lifecycle, live clock/store/network, credential
or user-data operation occurred.

## Reviewer challenges and freeze

Challenge the distinctions between presence, accepted outcome and worker evidence;
exact retained membership, byte-identical worker/accepted reply linkage, uncertainty
and historical candidate semantics, whole-state query nonmutation, inherited bounds,
and whether each mutation reaches its credited assertion. Verify both plan and
implementation without adopting this author explanation as proof.

Source and evidence FROZEN; no further edits without Brain direction. Final hashes:

- `tests/oracle.go`: `d3287284d42f5280fd4d2f80afeb898ac652e909697fdf140cefef17a0e77484`.
- `tests/oracle_test.go`: `3d63c8885e15d4c0ea5854272ffaab6b618882cae11e73ce21ddd1cff387c1b9`.
- `PLAN.md`: `c500bb6ce464255e475b35fa205d13997c53b92a37ffde7d0168700b2d56898c`.
- `inputs/ORIGINS.json`: `c41f2e64ee9dc9da27f0c12c1a510ff803cfe34ad653632d88768bc934f6dd18`.
