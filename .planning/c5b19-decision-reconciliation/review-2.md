# Independent engineering review

Review instance: 2 of 3.

Review status: PASSED. Bounded implementation/evidence verdict: **Ready**. Parent delivery and full Gate 3 remain IN_PROGRESS — TRENDING_GOOD. Child execution, Gates 4–6 and product admission remain BLOCKED.

## Findings

No actionable findings remain in the frozen review scope. Review1's P2 mutation-evidence finding is resolved by the corrected targeted mutation and consistent evidence corrections; no implementation or test bytes changed.

The correction preserves the complete exact-match/private-settlement branch and adds an `else if trial.history.state.Pending { expected = 1 }` bypass only for a wrong anchor with pending state. Thus exact preparation settles, runner begin succeeds, and the false acceptance reaches `tests/oracle_test.go:1905`. The retained corrected log has exactly twelve `invalid decision history accepted` assertions at1905 for operation1/2/3, generation1/2/3, kind1/2/3 and candidate1/2/3, with neither setup1887 nor build failure. The mutation ledger and script enforce that exact case set and assertion location. Source restoration matches the unchanged frozen hash, and the restored targeted command's retained log passes.

The original unconditional-bypass log remains preserved and explicitly excluded from wrong-pending sensitivity in PLAN, ORIGINS, revised handoff, review response and canonical roadmap. Its failure at1887 does not count as correction evidence. The historical initial handoff is retained as superseded testimony, not silently treated as current truth. Isolated necessity of overlapping exact-anchor and private-reply guards is not claimed.

## Plan Review

The current code meets the accepted bounded contract: value-owned private write history; validated trial-state transitions; exact pending anchor plus private settlement before code1; code2 preserving every history field and monotonically setting modeled fatal; atomic refusal on wrong decisions or bad evidence; three represented decisions; and strict serial control-sequence/nondecreasing-tick watermarks. Existing begin rules, reply validation and earlier APIs remain unchanged. Only confirmed outcome updates the claimed durable projection; failed and indeterminate outcomes still settle the exact pending operation. Later correct settlement preserves fatal.

The tests exercise independently encoded valid substitutions separately from malformed evidence, exact private-state equality, copied buffers/snapshots, reuse of ordinals after refusal, duplicate/stale rejection, late settlement outcomes, cap and boundary cases. The corrected mutant now satisfies the plan's required intended-assertion sensitivity. Inspection of the twelve remaining original mutation definitions and corresponding assertion locations supports their stated boundaries: wrong decision, missing settlement, rejection history mutation, fatal latch/preservation, decision cap, chronology and watermark updates.

The canonical adapter change is narrowly justified. `runPassiveAdmitted` owns a Background-rooted timeout's explicit cancel function. Worker publication to a one-element buffered channel precedes deferred slot release and context cancellation. A `context.Canceled` result in this branch therefore means the buffered result is available; retrieving it avoids misclassifying cleanup as deadline expiry. Actual `DeadlineExceeded` retains response-unknown behavior. Caller cancellation independence and slot lifetime remain intact. Tests synchronize readiness of result and cleanup and verify exact value/error preservation, caller cancellation and original deadline/stall/response-loss semantics. Randomized select choice is correctly disclosed.

No dependency, serialization, privileged responsibility, product import or admission change is introduced. The plan deliberately defers complete trace processing, interleaved receive/begin/apply groups, malformed rejected-frame diagnostics, omitted-event detection, stop/start/fatal admission, worker/snapshot correspondence, actual durability and lifecycle/timing proof. The next-scope recommendation keeps visible snapshot bytes distinct from confirmation. These limitations are consistent with ADR0047 and the frozen Gate2 packet.

Brain-owned retention, immutable archive pin, commit/push/PR publication, final metadata verification and hosted checks remain explicit delivery steps. This Ready verdict assesses the frozen implementation, evidence and plan; it does not assert those later steps have already occurred or authorize a child run.

## Author-Claim Reconciliation

| Claim | Evidence inspected | Status | Review consequence |
| --- | --- | --- | --- |
| Bounded atomic exact-pending reconciliation | oracle.go decisionHistory, unchanged writeHistory, new tests and scope | Confirmed | No code finding |
| Code2 preserves history; late settlement preserves fatal | reconcile and whole-state/projection assertions | Confirmed | Matches accepted contract |
| Corrected mutant reaches twelve intended wrong-anchor assertions | Exact script/ledger/log; line1905 in current test source | Confirmed from inspected execution evidence | Review1 P2 closed |
| Original mutant excluded; overlapping protection disclosed | PLAN, ORIGINS, current handoffs and roadmap diffs | Confirmed | No overstated isolated-guard claim |
| Implementation/test bytes frozen | Recomputed source hashes, initial and revised handoffs | Confirmed | Earlier full-check evidence remains applicable |
| Historical source/provenance preserved | Oracle/test/PLAN prefix comparisons against HEAD; every prior parsed ORIGINS value | Confirmed | No historical rewrite |
| Cleanup race regression fails before fix and passes afterward | Red log, source ordering, independent targeted race run | Confirmed | Narrow correction supported |
| Full archive test/race/vet/build and canonical make ci/pnpm test passed | Retained reports and verified log hashes; archive race independently rerun | Inspected except explicit reviewer commands below | No claim of independently rerunning all checks |
| Child/native/full-control/admission remains blocked | Code comments, plan, roadmap and exclusions | Confirmed | No authority expansion |

## Verification Performed

Authorized defensive environment: inert offline archive fixtures and local synthetic passive adapter harnesses only. No mutation script, native fixture, child, guest or external target was executed. No repository source, plan, branch, commit or PR was modified. Only reviewer reports and private temporary build caches were written.

Ground truth:

- Canonical root `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-decision-reconciliation`, HEAD/base `8e4fd22cc00cd5e961544f6a803974600dc0e751`.
- Archive root `/tmp/capsule-c5b19-recovery`, branch `codex/c5b19-benign-teardown`, HEAD/base `5d0e53cbccae3d4f587205df70253b776eae63e5`.
- Exact uncommitted bootstrap paths reviewed, including untracked canonical planning. Unrelated `.serena/project.yml` preserved. No new committed diff.

Independently executed with `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off`:

1. Archive experiment cwd, `GOCACHE=/tmp/c5b19-review-2-archive-cache`, exact `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.13.darwin-arm64/bin/go test -race -count=1 ./...` — PASSED, exit0, package1.880s.
2. Canonical cwd, `GOCACHE=/tmp/c5b19-review-2-canonical-cache`, exact `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.11.darwin-arm64/bin/go test -race -count=10 ./internal/execution/authorityplane -run 'TestPassive(CompletedResultSurvivesMethodContextCleanup|CallerCancellationRetainsMethodContextAndSlot|CancellationDeadlineStallAndResponseLossSemantics)$'` — PASSED, exit0, package1.522s.
3. Both worktrees `git diff --check` — PASSED.
4. Read-only Python comparisons — PASSED: all four ORIGINS implementation hashes, three canonical-reference hashes, all22 archive log hashes, corrected script/ledger/review1 hashes, exact twelve-case/assertion evidence, current restored source, historical source/test/PLAN prefixes and parsed ORIGINS values, both adapter hashes and both canonical post-fix log hashes. Rechecked the ten source/evidence hashes below at report creation.

Inspected, not independently executed in this review: archive original red/green/restored/vet/build, original and corrected mutant executions and targeted restoration; review1's full checks; canonical adapter red and hundred-repeat green/package race logs; canonical make ci and pnpm test verification reports/logs. No unchanged broad suite was unnecessarily repeated. Prior unfiltered lint debt#217 and dependency audit limitations remain disclosed.

Serena initialized. Codebase-memory Tier2: list_projects fully paginated; adapter symbol search fully paginated100+58, runPassiveAdmitted both-direction trace and snippet inspected. Coverage requested for relied-on source/test/model paths, canonical documents/planning and archive/evidence paths. Main-checkout generation2026-10-04T22:18:10Z is not the frozen review worktree: exact worktree/archive paths are outside_project, planning and recovery document are missing, and other main-checkout paths report best-effort metadata_match. Current direct source/diff reads provide fallback. No exhaustive graph/completeness claim.

The preliminary report `/tmp/c5b19-g3a8-review-2-preliminary.md` was written before detailed author handoffs, prior review report and mutation script/log inspection. Bootstrap and required plan disclosed the prior finding, preventing complete blinding to its existence; that limitation is recorded.

### Exact frozen source and evidence SHA256

| Absolute path | SHA256 |
| --- | --- |
| `/tmp/capsule-c5b19-recovery/experiments/typed-guest-transport-c5b19-benign-teardown/tests/oracle.go` | `f9cc9258a05f61eb3e2abee65c004aebfd77c4305e7f396bd87536d5f29ee65b` |
| `/tmp/capsule-c5b19-recovery/experiments/typed-guest-transport-c5b19-benign-teardown/tests/oracle_test.go` | `493df8d0d94187dc975a6fa763b69eebae33dbff6bb56ac254189aa23d73f36d` |
| `/tmp/capsule-c5b19-recovery/experiments/typed-guest-transport-c5b19-benign-teardown/PLAN.md` | `76c234feccaf3cd3477caf5754183dec0b7607c64d882428c3f8120a63142817` |
| `/tmp/capsule-c5b19-recovery/experiments/typed-guest-transport-c5b19-benign-teardown/inputs/ORIGINS.json` | `6916cfcd64657496b06aa850ff003e3a23bd9ac8f2432b8d9d778b87e8b102dd` |
| `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp/internal/execution/authorityplane/passive_product_adapter.go` | `785ac54438cea50838b39d8de8a4cf1b8f9b0dd6b887bfedd51beea5543cc9f4` |
| `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp/internal/execution/authorityplane/passive_product_adapter_test.go` | `4621e65f44dd4fcc1ded6f553d8ba29e2c38245fdd6550a51ef89a09274f9484` |
| `/tmp/c5b19-g3a8-targeted-mutate.py` | `9ff5defb4aade822fa7a39f2e591d3419fec2f6af27890b2777fc66ce43f7661` |
| `/tmp/c5b19-g3a8-targeted-mutation.json` | `6046bfd42347f68e372fe739920824db34355d0f248a3a14b3df0e4b0048283a` |
| `/tmp/c5b19-g3a8-mutant-targeted-wrong-pending.log` | `274b47e802ba77c1500ab91a663198e8ff89087b60a036de1f8c25be41bb7657` |
| `/tmp/c5b19-g3a8-targeted-restored.log` | `f1a20070b176b6c88a35e0d304d81a962bfe7892c2fbc0fb562292d01f9e92bb` |

## Open Questions And Residual Risks

No unresolved implementation or architecture question blocks this bounded candidate. Local checks ran on macOS; hosted Linux acceptance and final immutable publication remain Brain-owned. The corrected mutation establishes sensitivity to accepting a well-formed wrong anchor while another write is pending, not independent necessity of each overlapping guard. The synthetic selector regression controls readiness, not branch selection. Full Gate3, child execution and product admission are unchanged.

## Verdict

**Ready** for Brain's bounded G3a8 plus canonical adapter-correction integration. The prior P2 evidence defect is resolved, implementation/test source remains frozen, and no new actionable findings were identified. Review instance2 of3 is complete; no new reviewer, fix, split or heavy pivot was initiated.

## Recommended Next Actions

Brain should retain this report and preliminary ledger, update acceptance metadata and immutable archive links, verify those final metadata changes, then finish the named shared PR destinations and exact-head hosted checks. Preserve the historical failed evidence and current limitations. Child execution and product admission remain BLOCKED.
