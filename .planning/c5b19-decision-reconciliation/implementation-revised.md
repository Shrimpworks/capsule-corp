# G3a8 implementation author handoff

Status: PASSED for bounded implementation/local verification. Independent review
IN_PROGRESS — TRENDING_GOOD, pending instance 2 of 3 after accepted review1 P2 correction. Full Gate 3 IN_PROGRESS —
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
both remain. Original combined bypass failed setup and is excluded from mismatch sensitivity.
Corrected targeted bypass preserves exact settlement and fails all12 intended
assertions. Isolated necessity of overlapping guards is not claimed.

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
- Initial13 mutants compiled and failed assertions; review1 found the combined
  bypass fails wrong-anchor setup. It is excluded from mismatch sensitivity.
  Twelve other results stand; corrected targeted bypass supplies the thirteenth
  credited guard mutation. All12 intended assertion1905 failures, exact finally
  restoration and targeted green. No construction failures.
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
PLAN.md 76c234feccaf3cd3477caf5754183dec0b7607c64d882428c3f8120a63142817
inputs/ORIGINS.json 6916cfcd64657496b06aa850ff003e3a23bd9ac8f2432b8d9d778b87e8b102dd
```


# G3a8 review1 response

PASSED: accepted P2 evidence correction and targeted restoration check.
Fresh review2/3 IN_PROGRESS — TRENDING_GOOD, pending. Full Gate3 remains
IN_PROGRESS — TRENDING_GOOD; child execution/product admission BLOCKED.

The original combined mutant failed setup, not the intended wrong-pending
boundary. This claim is corrected and the original evidence explicitly excluded
from that conclusion; original artifacts are preserved. The other12 mutant
results stand. Corrected else-if-pending mutation preserves exact preparation
settlement. All12 valid mismatch cases fail assertion1905, none at setup1887.
Exact source finally restored, same targeted test passes. Only PLAN/ORIGINS.g3a8
metadata changed; source/test frozen hashes unchanged. No commit/push/PR or
child/native/live-clock/store/network effect. Details follow.


### Review1 P2 evidence correction

Review1 returned **Not ready** for the claimed wrong-pending mutation sensitivity,
with no reconciler implementation defect found. Accepted: the original mutant
replaced the entire match/settlement block with expected=1. Exact preparation
never settled, so all12 wrong-anchor cases failed during setup1887. Its original
script/log/hash and initial author/evidence artifacts remain preserved; it is
excluded from the wrong-pending claim. Fresh review2/3 is pending.

The corrected temporary mutant preserves exact preparation settlement and adds
only an else-if-pending wrong-anchor acceptance path. Restored implementation does
not contain this bypass:

```go
if trial.history.state.Pending && bytes.Equal(pair.Received.Request, trial.history.pendingRequest[:]) {
		if err := trial.history.applyReply(pair.Received.Reply); err == nil {
			expected = 1
		}
	} else if trial.history.state.Pending {
		expected = 1
	}
```

Using the exact Go1.25.13 executable/offline environment/cache above:
`go test ./... -run '^TestDecisionHistoryMismatchedAnchorsAndLateSettlement$' -count=1`
compiled and exited1. All12 subtests failed intended assertion
`oracle_test.go:1905: invalid decision history accepted`: operation1/2/3,
generation1/2/3, kind1/2/3, candidate1/2/3. No setup1887 failure occurred.
The script requires the exact12 case set, twelve target assertion occurrences,
no setup1887/build failure and nonzero exit. Source was restored byte-exact in
finally; the same targeted command then exited0. Source/test hashes remain
unchanged, so prior full test/race/vet/build results still apply; no redundant full
rerun. No child/native or live clock/store/network effects.

New evidence SHA256:

```text
mutant-targeted-wrong-pending 274b47e802ba77c1500ab91a663198e8ff89087b60a036de1f8c25be41bb7657
targeted-restored f1a20070b176b6c88a35e0d304d81a962bfe7892c2fbc0fb562292d01f9e92bb
targeted-script 9ff5defb4aade822fa7a39f2e591d3419fec2f6af27890b2777fc66ce43f7661
targeted-ledger 6046bfd42347f68e372fe739920824db34355d0f248a3a14b3df0e4b0048283a
review1 eb85cc2db8e12b84ce504bd3fe86a349c9ab25bc9863eff36157e5ebbbfe6ee0
```

Reproduction script /tmp/c5b19-g3a8-targeted-mutate.py and assertion ledger
/tmp/c5b19-g3a8-targeted-mutation.json complement this retained exact mutation,
command, result and hash record. Child execution/product admission remain BLOCKED.

New freeze:

```text
tests/oracle.go f9cc9258a05f61eb3e2abee65c004aebfd77c4305e7f396bd87536d5f29ee65b
tests/oracle_test.go 493df8d0d94187dc975a6fa763b69eebae33dbff6bb56ac254189aa23d73f36d
PLAN.md 76c234feccaf3cd3477caf5754183dec0b7607c64d882428c3f8120a63142817
inputs/ORIGINS.json 6916cfcd64657496b06aa850ff003e3a23bd9ac8f2432b8d9d778b87e8b102dd
```
