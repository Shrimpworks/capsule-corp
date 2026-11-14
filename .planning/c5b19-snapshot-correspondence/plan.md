# C5b19 G3a9 supplied durable-snapshot correspondence

Status: PASSED for bounded implementation, local checks and independent review.
Publication is IN_PROGRESS — TRENDING_GOOD. Full Gate3 remains IN_PROGRESS — TRENDING_GOOD;
child execution, Gates4–6 and admission remain BLOCKED.

## Objective and delivery boundary

Merge reviewed PR406, then implement only bounded supplied durable-snapshot
correspondence in the owned experiment archive, with tests-first evidence, fresh
independent review and a canonical roadmap PR. PR406 merged normally at
`f3f2e0755c0c284dd28019225bef0cf8482d3cdd` on2026-11-14T02:28:04Z
(2026-11-13 America/Toronto), matching reviewed head
`8ba5c6809fdd5d1d7357cdfa03eb378136920be6`. Exact-head CI43546021065 and
CodeQL43546021007 passed; no bypass or protection change. Post-merge CI43546467978 and CodeQL43546468003 also passed at the merge head.

Canonical branch `codex/c5b19-snapshot-correspondence` in the existing recovery
worktree, base f3f2e0755c0c284dd28019225bef0cf8482d3cdd. Archive branch
`codex/c5b19-benign-teardown`, base febb88d6421058f048f75319baae7cf9bd5a0cbb,
in /tmp/capsule-c5b19-recovery; existing draft PR38 remains integration destination.
Preserve unrelated .serena/project.yml. Brain owns commits/push/PR and integration.

## Internal assignments and ordering

- G3a9-SCOPE: read-only internal explorer; derive source-grounded bounded contract,
  including absent/superseded snapshots and visible-but-unconfirmed candidates.
- G3a9-VERIFY: internal verification worker; canonical baseline and evidence.
- G3a9-IMPLEMENT: internal worker after explicit scope acceptance; own only archive
  oracle.go/oracle_test.go, PLAN.md and ORIGINS.g3a9. Tests first; selected compiled
  mutants must reach their intended assertions; exact restoration required.
- G3a9-DOCS: Brain; canonical roadmap, source pins and collected evidence.
- G3a9-REVIEW: fresh-context read-only reviewer after freeze; max3 instances for this
  distinct bounded slice. Prior budgets and missing historical reports stay intact.

Workers share Brain branches/worktrees, are not independent delivery units, and
must preserve others' edits. No worker branch/commit/push/PR. Scope and baseline
may proceed in parallel; behavior waits for scope acceptance. Brain collects all
handoffs, reconciles findings, verifies and publishes both named destinations.

## Defensive boundary and acceptance

Only owned inert synthetic byte fixtures and ordinary offline unit/compiler
processes. No native fixture, child/guest, lifecycle/signal/wait target, live store/
clock/network, credential or user-data operation. Canonical routine checks may
query their existing advisory sources. Keep archive Go1.25.13 frozen offline,
manifest/packet/ABI, earlier APIs and provenance unchanged. No new dependency or
privileged responsibility. C5b19 BUILD-NARROWLY reuse row remains test-only.

Record presence must never become confirmation. Correlate exact supplied395-byte
records with actually recorded submitted candidates and accepted decision evidence;
failed/indeterminate/unapplied publication can leave visible bytes. Define absent,
stale/superseded and confirmation classifications without inventing unobserved
publication order or claiming full trace/recovery/durability/control admission.
The prior next-scope recommendation guides, but does not replace, current scope.

Acceptance: explicit bounded contract; red/green tests; malformed/substitution/
uncertainty/nonmutation/bounds evidence; selected mutants reach intended assertions
and restore exact source; offline test/race/vet/build; canonical required checks;
fresh review; exact archive pin, committed/pushed PR destinations; body/newlines/
base/head/draft verification and hosted gate readback. Child execution stays gated.

## Retrieval context

Serena session d55b465c initialized. Codebase-memory Tier2 primary index generation
2026-10-04T22:18:10Z, old/not watched, not this worktree. Full project listing paged;
SettleWrite/applyDurableProjection search complete. Packet/model/types/ADR47 report
best-effort metadata_match; prior next-scope missing freshness and archive outside
project use exact source reads. No completeness claim. Current packet253–271
forbids record-presence confirmation;342–359 requires accepted pair evidence and
comparison with worker reply/snapshot. Scope acceptance must preserve these limits.

## Canonical baseline

PASSED: frozen install, exact Go1.26.11 `make ci`, standalone `pnpm test` and diff
check. Each test run reports177 workspace and132 script tests, zero failures.
Govulncheck reports no vulnerabilities; dependency audit retains2low/2moderate
below unchanged high threshold. Unfiltered lint exits1 for known#217 revive debt,
50 displayed under output cap, not total. `canonical-verification.md` and
`validation.json` retain exact commands/log hashes/versions and limitations.
No canonical code/config changes planned; final retained metadata checks remain
Brain-owned. Post-merge CI/CodeQL also passed on f3f2e07. Child remains BLOCKED.

## Accepted G3a9 contract

PASSED: scope derivation accepted before behavior edits. `scope.md` defines a new
value-owned wrapper retaining at most3 exact submitted requests and their accepted
code1 replies, composed atomically with unchanged G3a8. Nonmutating explicit-
presence queries strictly match a supplied record to a retained submission and
separately report historical AcceptedOutcome, optional strict kind5 WorkerOutcome,
and order-only LaterSubmitted. Supplied worker reply must match any accepted reply
byte-for-byte. Record presence/worker outcome cannot manufacture control settlement.
No full-confirmation/admission/final-verdict boolean; absence never erases facts or
certifies a consistent final store. Earlier records remain queryable without
physical supersession/rollback claims. All query paths preserve private state.

Mutation evidence must reach intended assertions after valid setup. Overlapping
record/request/reply/digest guards and redundant capacity guards must not receive
isolated sensitivity claims. Test independent expected bytes and meaningful
composed guards; retain failed constructions honestly. Full capture/control/timing/
durability/native evidence remains outside scope. Implementation may now proceed.

## Deferred next-slice recommendation

Read-only G3a9-SCOPE follow-up returned `next-scope.md`: supplied worker-received/
sent sidecar correspondence, preserving independent raw capture roles and partial
stuck/drop evidence. This is a later bounded recommendation, not implemented or
child-authorized work. Full control-trace composition remains separately deferred.

## Implementation freeze

PASSED for bounded implementation and local checks: ten new tests; meaningful
compiler-red before green; twelve compiled selected mutants each reached exact
intended assertions (membership4, worker-request8, other targets1 each); byte-exact
restoration and full restored test/race/vet/build. No unexpected failure, surviving
mutant or invalid construction. Original source/test prefixes, go.mod/model_test
and all prior parsed provenance remain unchanged. Exact definitions, command,
assertion/count and log hashes are retained in archive ORIGINS.g3a9, not only /tmp.
`implementation-handoff.md` freezes four archive paths and their hashes.
Fresh independent review1/3 is next; readiness and hosted acceptance unclaimed.

## Independent acceptance

PASSED: fresh review1/3 returned **Ready**, no actionable findings or heavy pivot.
Reviewer independently passed exact offline Go1.25.13 test/race/vet/build and both
diff checks, verified frozen source/reference/provenance and all20 archive author
logs plus4 canonical baseline logs, and inspected all12 intended mutation assertions.
Mutation scripts and canonical full checks were inspected, not independently rerun.
Report SHA256 `eb877a513f345642468b776f0b9fc6e6ac8ed6b64c1949c45e09dc719d4103bb`.
Brain accepts the report. No implementation/test change follows this freeze; only
acceptance metadata, immutable pins and publication remain. All internal handoffs
collected. No extra review or scope split is needed for these metadata steps.

## Archive publication and integration

PASSED: archive commit `7646d8faabb670ce8b5ff1b213c414d91deda2d9` pushed to
`codex/c5b19-benign-teardown`, existing draft PR38. Archive worktree is clean.
Reviewed oracle/test bytes and historical provenance remain unchanged; acceptance
metadata records the independent report hash. Canonical roadmap pins this exact
archive commit. Canonical branch `codex/c5b19-snapshot-correspondence` carries only
roadmap and retained evidence; unrelated `.serena/project.yml` is excluded. All
internal handoffs are collected; all retained changes have these shared destinations.

Remaining delivery: final metadata checks, canonical commit/push/ready PR, archive
PR description refresh, and exact-head hosted CI/CodeQL readback. These final
publication facts will be recorded in PR bodies without changing reviewed code.
No child fixture ran; full Gate3 and later gates remain incomplete.

## Final metadata verification

PASSED: explicit retained-JSON Biome validation, `pnpm lint` after staging all
retained evidence, `pnpm site:build` and staged diff check. Only this result ledger
was added afterward; no code or configuration changed. Ephemeral log identities:

- `json`: exit0, SHA256 `9db65e5b8464d24822105a525af2b621ca9d3a5e56b109b66ca01893241df743`.
- `lint`: exit0, SHA256 `647347038eb77e708f0fe75e1e09d6616d0f5b0e15c1c4e19cf330f56d009315`.
- `site`: exit0, SHA256 `198149bc8cd29d5959ebe14aaa0721396f6a4b1bc52d5737f767ca97aaa9b464`.
- `diff`: exit0, SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.

Final publication/hosted acceptance and body/base/head/draft readback are recorded
in the published PR bodies. Full Gate3 remains IN_PROGRESS — TRENDING_GOOD; child
execution, Gates4–6 and admission remain BLOCKED. No new run authority is granted.
