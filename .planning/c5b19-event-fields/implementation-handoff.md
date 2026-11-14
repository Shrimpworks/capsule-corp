# G3a5 implementation handoff and author explanation

Status: PASSED for bounded implementation and local verification; independent
review IN_PROGRESS — TRENDING_GOOD, pending instance 1 of 3. Full Gate 3 remains
IN_PROGRESS — TRENDING_GOOD. Child execution, Gates 4–6 and admission BLOCKED.

## Intent and success criteria

Add an explicit opt-in field/lane validator after G3a4 syntax, using the accepted
29-kind table-defined subset. It must reject wrong writers, every forbidden
non-common field, missing operation/identity members, invalid generation/PID/
site boundaries, reply outcome/digest mismatch and a nonzero wait status without
a matching PID return. It must preserve scalar zero observations and leave
syscall, paired-frame, causal, timing and lifecycle evidence unclaimed.

Defensive scope: owned archive synthetic in-memory bytes only. No child fixture,
guest, process target, signal/wait, filesystem/store, network, credential, real
clock, service, user data or goroutine effects in new code/tests. Normal offline
compiler/unit-test tools were the only executed processes for this assignment.

## Frozen review target and integration destination

Archive /tmp/capsule-c5b19-recovery, branch codex/c5b19-benign-teardown,
HEAD/base 3e26b809a40832a3351fdc6b5e0ffc1ef0094792. Explicit working diff:
experiments/typed-guest-transport-c5b19-benign-teardown/{tests/oracle.go,
tests/oracle_test.go,PLAN.md,inputs/ORIGINS.json}. No other archive changes.
Earlier source is an unchanged prefix of both Go files; go.mod/model_test.go
are byte-equal to baseline. All previous ORIGINS values are unchanged; only g3a5
is appended. No branch, commit, push or PR was created. Root owns existing draft
archive PR38 plus the new canonical PR on codex/c5b19-event-fields.

Canonical reference fe3d22c84b5d41e0a4edbfacb531563040e16995. Acceptance is
/tmp/c5b19-g3a5-scope.md and canonical .planning/c5b19-event-fields/plan.md.
Read archive AGENTS and current Gate2 field/lane, process identity and harness
ambiguity sections directly, with prior ecosystem/ADR context inherited. Canonical
packet/model/type hashes are in ORIGINS. Inherited best-effort graph coverage
reports metadata_match for those canonical paths at generation
2026-10-04T22:18:10Z; archive is unindexed and exact source was read directly.
No new structural graph or coverage claim is made.

## Plan traceability and implementation flow

1. Added failing tests before decodeEventFields existed. Missing-symbol compiler
red is retained, followed by green after the implementation.
2. decodeEventFields accepts raw bytes and invokes unchanged decodeEventEnvelope,
so direct struct construction cannot bypass framing/binding/timebase checks.
3. A local 30-slot table visibly names every kind 1–29 and its writer/allowed
field groups. Only trace-fault kind 28 permits every known lane. Non-common
presence bits are checked against the row; all forbidden nonzero fields fail.
4. Required tuple/identity rules and bounded scalar relations follow. Operation
IDs and request/candidate digests must be nonzero; generation is 1..2^53-1.
Identity requires a nonzero digest and positive signed-range PID. Spawn result
permits PID zero without pretending to know syscall success. Outcome/application
presence, transport digest coupling, store site and wait-return relation close
the local checks. Every refusal returns zero; success returns the copied envelope.
5. PLAN retains exact boundary/dependency decision, commands, excerpts and hashes;
ORIGINS appends references, current implementation hashes and log hashes. Both
execution/admission flags stay false. Independent review and publication pending.

## Component responsibilities and choices

Only tests/oracle.go adds behavior: 130 lines, one raw-byte wrapper and a fixed
mask vocabulary. tests/oracle_test.go adds 342 lines of independent byte-offset
fixtures and boundary/negative cases. PLAN/ORIGINS append evidence; all earlier
behavior/provenance remains untouched.

A visible fixed table is chosen over scattered per-kind ad hoc checks so a reviewer
can audit all 29 rows together. The presence mask separates forbidden zeros from
required nonzero tuples, avoiding the false rule that every named scalar must be
nonzero. Reusing G3a4 avoids a second frame decoder. No new lane walker, production
import, schema, dependency or manifest entry is needed. Standard-library-only
Go 1.25.13 and existing bytes/binary/error primitives suffice.

## Invariants and ambiguity limits

- G3a4 remains syntactic-only and accepts its old literal; callers opt into G3a5.
- Kind 29 supports payload-free containment only. Optional harness identity has
  no fixed slot mapping here; payload-bearing variants are unsupported by this
  subset, not declared forbidden by the whole frozen packet.
- Kind 8 supports complete represented replies for both decisions 1 and 2.
  Arbitrary malformed-reply diagnostic mapping is not invented.
- Transport kinds 5–8 require record digest iff confirmed; store kind 27 always
  requires zero record digest, even for confirmed outcome.
- Zero wait status is allowed. Nonzero status requires returned target PID but
  proves neither Darwin terminality nor a preceding wait or real absence.
- Signed result/errno fields and zero tick/anchor values are observations. No
  syscall success relation, missing-observation inference or time math is added.
- Successful single-record shape validation cannot settle history or grant
  start, signal, completion, release, custody, absence or a case-pass verdict.

## Verification performed

All Go commands used:
GOCACHE=/tmp/c5b19-g3a5-go-cache GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go

- test ./... -count=1 initially exit 1: undefined decodeEventFields.
- test ./... -count=1 after implementation exit 0 (0.261s).
- test -race ./... -count=1 exit 0 (1.795s).
- vet ./... and build ./... exit 0, no diagnostics.
- Restored test ./... -count=1 exit 0 (0.236s).
- git diff --check exit 0.

Tests transcribe literal byte offsets independently from implementation masks,
using the existing independently calculated common binding digest. All 29 kinds
are tested against all three lanes, including all 31 supported combinations.
Every forbidden non-common field is injected independently. Missing tuples/
identities, generation and PID bounds, represented reply outcomes and both
kind 8 decisions, invalid record coupling, all store sites, store outcome with
zero digest, matching/nonmatching wait return, zero wait status, signed spawn/
errno observations, zero scalar observations, envelope composition, zero error
results and caller ownership are covered. Absence shape acceptance explicitly
requires no prior wait in this local API and does not claim causal evidence.

Eight source mutants compiled and failed intended assertions: lane ownership,
forbidden mask, tuple presence, identity presence, kind 27 record digest,
reply outcome/digest coupling, store site bounds, wait status/return relation.
Source restored exactly after each in a finally block; final full tests passed.
/tmp/c5b19-g3a5-mutations.json records test names, exit codes, failing assertions,
log hashes and restoration hash. All logs /tmp/c5b19-g3a5-*.log are ephemeral;
essential excerpts/results/hashes are retained in PLAN/ORIGINS. This is selected
sensitivity, not exhaustive mutation or concurrent-safety proof.

## Costs, risks, deferrals and reviewer challenges

The fixed masks/table add small auditable parser complexity; no dynamic protocol
or runtime component is introduced. The tables in tests are an independent raw-
offset transcription, not reused production masks. Review should compare both
against the frozen packet, especially kind 27 vs transport replies and kind 24
vs kind 25 status handling. Check signed conversions inherited from G3a4, every
required member, generation/PID boundary and zero-result error behavior.

No intended deviation from accepted scope. Optional harness identity mapping,
malformed-reply diagnostics, syscall conventions, native terminal wait status,
raw sidecars, pairing, complete-capture, causal/timing composition, native/store/
build closure and all execution remain deferred. This exact subset must not be
presented as the full event oracle. Historical review counts are unchanged.

## Frozen SHA256

- `tests/oracle.go`: `d82c587990d8516e89b442d1d89707fed7a2114a6a906c9b76a624dfd39fd114`
- `tests/oracle_test.go`: `fb643536635b8e73fb2c2d52aeae782b7350162613b09577a01cbc50f34ad410`
- `PLAN.md`: `c64725ced2c26f9f43ef34a1cd4db120cf7a5f4e1075d309ffea0f03c6b67d28`
- `inputs/ORIGINS.json`: `8059c0f5f0a9d5cec04838b687d96f8bb278fb07ab6cf5c80d0c9202872ad944`

Recommended next action: root retains this handoff, freezes the four files,
starts fresh review 1 of 3, then owns fixes and exact publication after acceptance.
