# C5b19 G3a8 pending-write/control-state decision reconciliation

Status: PASSED for bounded implementation, local checks and independent review.
Publication is IN_PROGRESS — TRENDING_GOOD. Full Gate3 remains IN_PROGRESS — TRENDING_GOOD;
child execution, Gates 4–6 and product admission remain BLOCKED.

## Objective and delivery boundary

Merge reviewed PR #405 and deliver a bounded inert pending-write/control-state
decision reconciliation increment under frozen Gate2. #405 merged normally at
8e4fd22cc00cd5e961544f6a803974600dc0e751 on 2026-11-14T01:53:15Z
(2026-11-13 America/Toronto), matching reviewed head1198bba.
CI43542813931 and CodeQL43542814060 passed; no bypass or protection change.

Canonical branch codex/c5b19-decision-reconciliation, existing recovery worktree,
base8e4fd22cc00cd5e961544f6a803974600dc0e751. Archive destination: draft PR38,
branch codex/c5b19-benign-teardown, base5d0e53cbccae3d4f587205df70253b776eae63e5.
Preserve unrelated .serena/project.yml. Brain owns integration, commits and PRs.

## Assignments and dependencies

- G3a8-SCOPE: internal read-only agent, derive exact state/decision contract and
  prevent expansion into unsupported full control/timing/native semantics.
- G3a8-VERIFY: internal verification worker, canonical gates/evidence; parallel.
- G3a8-IMPLEMENT: internal worker after scope acceptance, tests first, minimal
  pure state composition and selected guard mutants with exact restoration.
- G3a8-DOCS: Brain, roadmap/evidence reconciliation and immutable archive pin.
- G3a8-REVIEW: fresh-context independent read-only review of frozen source/plan,
  maximum three instances for this distinct slice; no historical budget reset.

Workers share Brain worktrees, preserve all other edits and publish no independent
commits/PRs. Brain collects every handoff and reconciles before delivery.

## Defensive scope and acceptance

Only owned inert synthetic in-memory bytes and ordinary offline compiler/unit
processes. No native fixture, child/guest, lifecycle/signal/wait target, live clock/
store/network, credential or user-data operation. Keep archive Go1.25.13, manifest,
packet/ABI, historical provenance and earlier APIs. No dependency, privileged
responsibility or product import. No actual execution authority from modeled state.

Acceptance: source-grounded bounded state/decision contract; meaningful red before
minimal implementation; positive/negative/atomicity/boundary tests and compiled
mutants with exact restoration; offline test/race/vet/build; canonical required
checks; fresh independent review; exact archive source pin and all commits pushed
to named PR destinations; body/newlines/head/base/draft and hosted gate readback.
Separate represented pending-write reconciliation from full lifecycle/stop/timing
and actual durability/authority claims. Scope acceptance precedes behavior edits.

## Tooling and retrieval

Serena initialized, canonical project/language servers reported active.
Codebase-memory Tier2 primary index ready at2026-10-04T22:18:10Z, not watched.
Packet/types/model/ADR47 coverage best-effort metadata_match; current exact source
reads supply archive/new planning fallback. SettleWrite/samePending lookup complete;
SettleWrite outbound4 and source establish exact pending/outcome checks before
projection, with late settlement preserving control latches. No completeness claim.

## Post-merge stop-the-line correction

IN_PROGRESS — TRENDING_BAD: post-merge CI 43543205245 failed in
`TestPassiveCancellationDeadlineStallAndResponseLossSemantics` at line 453;
CodeQL 43543205196 passed. The pre-change local baseline passed all required
checks. Investigation found that `runPassiveAdmitted` publishes its result then
cancels the method context during worker cleanup; the select may mistake that
cleanup cancellation for deadline expiry. This is an existing adapter race.

Before any G3a8 implementation, internal assignment CI-RACE owns only
`internal/execution/authorityplane/passive_product_adapter.go` and its test file
in this branch/worktree. Defensively reproduce and fix result/cleanup arbitration
using local passive synthetic harnesses; preserve actual deadline disposition,
caller cancellation independence, admission slot lifetime and all safeguards.
No child, guest or external target is authorized. Add tests first, capture red
and green evidence, verify the original scenario, and hand off without commit,
push or PR. Brain explicitly groups this necessary CI correction in the same
canonical PR and independent review; resume G3a8 after targeted verification.

The scope contract in `scope.md` is accepted: value-owned trial reconciliation,
exact private pending matching, represented rejection with monotonic modeled
fatal, atomic invalid-evidence refusal, three decisions and serial control
watermarks. History recording does not grant control admission. Child execution
remains BLOCKED. The exact-anchor and private-settlement checks are redundant
for well-formed differing anchors; mutation evidence must disclose this rather
than claim independent sensitivity for both.

## Stop-the-line targeted verification

PASSED: the new synchronized completed-result regression failed against the
unchanged adapter on attempt0, then passed after the minimal cleanup-cancellation
check. The new regression, caller-cancellation/slot-lifetime test and original
deadline/cancellation test passed100 repetitions. G3a8 resumed only after this
result. Broader post-fix canonical and package race checks follow; no future
outcome is claimed here. Hosted failure log SHA256:
`143b82f60fe539b692e89e5a1b91ae7564ef7ee7d4be210f2f0e1dffc1b4177e`.

Read-only next-slice recommendation is retained in `next-scope.md`: supplied
snapshot correspondence, with visible bytes kept distinct from confirmation.
This is a later scope recommendation, not implemented or authorized child work.

## Post-fix canonical gates

PASSED: exact Go1.26.11 `make ci`, standalone `pnpm test` and diff check. Both test
runs report177 workspace and132 script tests, zero failures. The two CI-RACE
source hashes matched the worker freeze and remained unchanged through checks.
Frozen install and unfiltered lint are baseline evidence, not repeated; known
#217 output-capped revive debt and two low/two moderate dependency findings stay
open. Retained `canonical-verification.md` and `validation.json` distinguish this
post-fix result from the earlier baseline. Final metadata remains Brain-owned.

## G3a8 implementation freeze

PASSED for local implementation checks: nine new tests, exact offline Go1.25.13
unit/race/vet/build, thirteen compiled mutant assertion failures, byte-exact source
restoration and prior source/provenance preservation. The initial clock-policy
mismatch fixture was malformed under frozen binding and moved to malformed-input
coverage; its failure is disclosed rather than counted as valid mismatch evidence.
`implementation-handoff.md` retains the author explanation and exact four-file
hashes. Fresh independent review1/3 is next, including the canonical CI-RACE fix
and this plan. Review readiness and hosted acceptance remain unclaimed.

## Independent review instance1 and accepted correction

Review1/3 completed, verdict **Not ready**, one accepted P2 evidence finding.
The combined pending/private-settlement mutant replaced the entire settlement
block with `expected=1`, preventing exact preparation settlement. All12 valid
mismatched-anchor cases therefore failed runner-begin setup at line1887, before
the false-acceptance assertion at line1905. Its compiled failures do not support
the claimed wrong-pending mismatch sensitivity. No other actionable finding;
reviewer independently passed archive test/race/vet/build, canonical authorityplane
race and targeted100-run checks, provenance/hash and both diff checks.

Author response: **Accept**. Correct only mutation/evidence: retain exact-match
settlement, bypass matching/private validation for wrong-anchor pending, require
all12 cases to fail at the intended assertion, restore exact source, and run
restored checks. Preserve old failure evidence as excluded from this claim.
Production/oracle/test bytes stay frozen. This is bounded evidence repair, no
architecture or plan pivot. Fresh review2/3 will assess the correction and
unchanged implementation; no historical budget reset. Publication remains
BLOCKED until evidence and review close. Child execution remains BLOCKED.

## Corrected mutation freeze

PASSED: corrected mismatch-only bypass preserves exact preparation settlement;
all12 valid substitution/outcome subtests fail the intended false-acceptance
assertion at1905, with no setup failure at1887. The compiled mutant exits1;
byte-exact restoration then passes the same targeted test. Both original source
and test hashes are unchanged. Twelve prior targeted mutants remain valid; the
corrected mutant replaces the original for wrong-pending sensitivity. The original
log/script/hashes and initial testimony remain retained, explicitly excluded from
that claim. `review1-response.md` and `implementation-revised.md` record exact
mutation, observed assertion counts, hashes and limits. Fresh review2/3 is next.

## Independent acceptance

PASSED: fresh review2/3 returned **Ready**, no actionable findings. Review1's
accepted P2 is closed by the intended12-assertion mutation evidence and exact
restoration. Reviewer independently ran full archive race tests and three
canonical adapter regressions under race ten times; inspected the mutation
execution evidence, prior full checks and provenance, and verified all relevant
hashes and both diff checks. Report SHA256:
`8a0fd66f9acdfcc623660f65d138dd56430cc2e4ff85747707f442db7e3728a8`.
No behavior/test source changed after review1 or review2. Brain accepts the report
and retains both review instances and the correction response. No third review
or architecture pivot is needed for acceptance-only metadata/publication.

## Archive publication and final integration

PASSED: archive commit `febb88d6421058f048f75319baae7cf9bd5a0cbb` pushed to
`codex/c5b19-benign-teardown`, existing draft PR38. Archive worktree is clean.
Reviewed oracle/test bytes and all prior provenance remain unchanged; acceptance
metadata records review2 report hash. Canonical roadmap pins that exact archive
commit. Canonical branch `codex/c5b19-decision-reconciliation` carries the passive
adapter correction, roadmap, accepted scope, both reviews and verification.
Unrelated `.serena/project.yml` remains excluded. All internal handoffs collected;
no retained worker commit lacks the shared PR destination. No child fixture ran.

Remaining delivery: final metadata checks, canonical commit/push/ready PR, archive
PR body refresh and exact-head hosted CI/CodeQL readback. These publication facts
will be recorded in the PR bodies without changing reviewed implementation.

## Final metadata verification

PASSED: explicit Biome checks for both retained JSON files, `pnpm lint` after
staging all retained evidence, `pnpm site:build` and staged diff check. Canonical
source hashes still match independent review. Logs are ephemeral /tmp evidence;
exact hashes below retain their identities. Only this result ledger and a prose
completion were added afterward; no implementation or configuration changed.

- `json`: exit0, SHA256 `c943f80d97ee27c0306b6fa62ad39ec0e2f32d37d6e2c4c94958bb9b75b08254`.
- `lint`: exit0, SHA256 `29b6db709fadbe3352ccb01b64175b8376dab97613d1d4276e3aec312e4b9b90`.
- `site`: exit0, SHA256 `198149bc8cd29d5959ebe14aaa0721396f6a4b1bc52d5737f767ca97aaa9b464`.
- `diff`: exit0, SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.

Final hosted checks and PR readback are recorded in the published PR body. Full
Gate3 remains IN_PROGRESS — TRENDING_GOOD; child execution, Gates4–6 and admission
remain BLOCKED. This bounded implementation/review completion grants no run authority.
