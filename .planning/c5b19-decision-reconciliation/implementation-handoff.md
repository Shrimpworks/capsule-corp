# G3a8 implementation author handoff

Status: PASSED for bounded implementation/local verification. Independent review
IN_PROGRESS — TRENDING_GOOD, pending instance 1 of 3. Full Gate 3 IN_PROGRESS —
TRENDING_GOOD; child execution, Gates 4–6 and product admission BLOCKED. This is
an author packet, not an independent Ready verdict.

## 1. Intent and authorized scope

Defensively reconcile pending-write and represented application decisions using
owned inert synthetic request/reply/event fixtures in ordinary offline Go tests.
Only four assigned archive files changed. No native fixture, child, guest,
lifecycle/signal/wait target, live clock/store/network, credential or user-data
operation. No dependency, product import, packet/ABI/manifest/toolchain change.

## 2. Traceability and grounding

Accepted scope: /tmp/c5b19-g3a8-scope.md. Canonical reference
8e4fd22cc00cd5e961544f6a803974600dc0e751; archive baseline
5d0e53cbccae3d4f587205df70253b776eae63e5. Gate2 packet kind8 exact settlement/code2
fatal rule and passive model SettleWrite/samePending late projection behavior
underlie the contract. Prior APIs stay byte-identical. Codebase memory parent
Tier2 snapshot was old/unwatched; exact source reads resolve current authority.
Archive remains unindexed with direct source fallback. No graph completeness claim.

## 3. Data flow and implementation

newDecisionHistory validates immutable context and owns value-copied writeHistory.
Strict kind3 begin and individual G3a7 received/applied pairs enter copied trial
state. Shared control watermark enforces increasing sequence/nondecreasing ticks.
Code1 derives from exact 397-byte pending equality plus successful private reply
settlement. Code2 preserves every history byte and latches fatal. Wrong recorded
decisions/invalid evidence refuse the entire trial. Commit is one value assignment;
returned snapshots contain values only. Three represented decisions bound cost.

## 4. Files and ownership

Append-only oracle.go behavior and oracle_test.go tests, appended PLAN evidence,
and new ORIGINS.g3a8 metadata. go.mod/model_test.go unchanged; earlier source/test
prefixes and all earlier parsed ORIGINS values verified exact. Worker made no
branch, commit, push or PR; shared archive branch codex/c5b19-benign-teardown and
existing draft PR38 are the Brain-owned integration destination.

## 5. Decisions and alternatives

Use individual G3a7 validation, since collection duplicate-operation rejection
would prevent valid stale/duplicate code2 evidence. Whole-pair processing means
serial selected groups only; begin interleaving receive/apply is unsupported,
without declaring that full trace invalid. Fatal is rejection-derived only;
history-only begin after fatal is explicitly allowed, without actual admission.
No invented stop/start/timeout policy. Redundant exact bytes and reply hash checks
both remain. A combined bypass mutant demonstrates mismatch sensitivity, not
isolated necessity of equality in the presence of private reply validation.

## 6. Invariants and test matrix

Exact one-pending history/generation/unique ID/three-write/terminal guards compose
unchanged. Confirmed alone changes projection; all outcomes settle exact pending.
Rejection never changes hidden IDs/attempts/pending/generation/projection; late
settlement never clears fatal. Wrong decisions and chronology errors consume no
budget/watermark. Tests independently encode events and compare expected candidate
bytes, include same-ordinal retries after refusal, null/uninitialized context,
malformed members, valid mismatched anchors, stale/duplicate rejection, zero/gaps/
equal ticks/max uint64, rejection cap, one-pending and terminal refusal, caller
buffer/snapshot independence, input nonmutation, and history-only begin after fatal.

## 7. Verification and honest failures

Exact installed Go1.25.13 binary under GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off
GOCACHE=/tmp/c5b19-g3a8-go-cache, cwd archive experiment:
- tests-first go test ./... failed as expected with undefined decisionHistory API.
- initial implemented test run exposed an invalid clock-policy mismatch fixture;
  policy is frozen. Corrected it to malformed-evidence coverage, retained failure.
- go test ./... green.
- 13 selected mutants: go test ./... -run TestDecisionHistory -count=1;
  all compiled and failed behavioral assertions; no mutation construction failures.
- exact source restored; final go test ./..., go test -race ./..., go vet ./...,
  go build ./... all exit0.
- git diff --check exit0; historical prefix/identity/provenance checks passed.

Raw logs /tmp/c5b19-g3a8-NAME.log; hashes retained in PLAN/ORIGINS. Mutation results
and retained assertion excerpts /tmp/c5b19-g3a8-mutations.json; reproduction script
/tmp/c5b19-g3a8-mutate.py uses finally restoration. No native child execution.

## 8. Cost and dependency footprint

Constant-size value copies; fixed 256/397/98-byte validators; at most three begins
and three represented decisions. Standard library only, no I/O/unbounded loop or
new primitive. The C5b19 ADOPT-PLATFORM / BUILD-NARROWLY recommendation does not
admit this archive into products.

## 9. Limits and deferred work

No full control model, actual submission admission, stop/apply decision causality,
omitted-event detection, worker/durable snapshot correspondence, rejected malformed
diagnostic encoding, deadlines, real fsync/custody, native effect, child execution,
run permission or product admission. Full Gate3 and later gates remain open.
No inference from a successful pure selection to execution authority.

## 10. Review challenges and freeze

Review should challenge wrong-decision atomicity, entire history preservation,
post-rejection late settlement/fatal persistence, chronology boundaries and serial
selection limits, plus independent fixture correctness. Brain owns remaining
independent review, canonical docs/pins, source commit and publication.

SOURCE FROZEN: worker will make no further edits without Brain review/fix direction.
Four owned paths and SHA256 follow:

```text
tests/oracle.go f9cc9258a05f61eb3e2abee65c004aebfd77c4305e7f396bd87536d5f29ee65b
tests/oracle_test.go 493df8d0d94187dc975a6fa763b69eebae33dbff6bb56ac254189aa23d73f36d
PLAN.md 886544e141cb78d1bddb75f8eed0a45f7a85b9c709d6db78481a9ab5f822bea0
inputs/ORIGINS.json 90e2e7e10475dd38b9898ff949628687307fa3ab2e982de477b5f429ffbc843b
```
