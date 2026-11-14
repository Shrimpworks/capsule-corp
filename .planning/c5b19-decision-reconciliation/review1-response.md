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
