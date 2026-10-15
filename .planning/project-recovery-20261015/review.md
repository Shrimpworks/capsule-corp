# Independent recovery review

Review instance: 1 of 3.
Work item: bounded review of RECOVERY-20261015.
Status: PASSED for the bounded independent review. Integration/publication remains IN_PROGRESS — TRENDING_GOOD; product admission remains BLOCKED.
Scope: working-tree diff on codex/project-recovery-20261015 at base/HEAD fc8af2df91d443f52a7677c2bf9f7dc1d96f9ab1, including recovery report and .planning/project-recovery-20261015. Defensive local contributor-tooling and documentation review only. No product, guest, service, credential or cross-repository mutation.

## Findings

No actionable implementation finding identified. The steering source validation precedes installer invocation and directory/link mutation; all three required sources must resolve to files. Existing owned-file preservation, unrelated links, mismatch refusal and dry-run semantics remain intact. The Go change only raises the module minimum; dependencies, scanner version, schemas, historical packets and CI controls remain unchanged.

Preliminary inspection before opening the author explanation covered Git identity/scope, full linker diff and surrounding flow, tests, canonical status rules, current plan, recovery report, Makefile and CI. Preliminary concerns were final evidence completeness and stale planning status; these were communicated to the parent as finalization items, and the updated plan/appendix address them. Process limitation: the preliminary observations were not saved to a separate ledger before the author explanation was opened. The implementation was inspected first, but the record of that separation is less auditable than the requested workflow.

## Plan Review

The report has a bounded contributor-readiness outcome, explicit ownership, dependencies, acceptance criteria, risks and effort. Immediate maintenance repairs precede broad CI parity and installed authority work. The next C5b19 lane requires archive-state reconciliation and actual-byte review; the historical packet and separate first-run gate remain unchanged. Missing v19/v27 evidence is unresolved rather than reconstructed or promoted into an admission claim. Product admission remains BLOCKED; the product parent remains IN_PROGRESS — TRENDING_GOOD.

The dependency decision maps the toolchain to GO-3/TEST-1, records exact local identity, limits provenance claims, names maintenance ownership and states rollback consequences. Lack of release-grade offline compiler/advisory evidence and broader platform evidence is explicitly outside this maintenance admission scope.

## Author-Claim Reconciliation

| Claim | Evidence | Status / consequence |
| --- | --- | --- |
| All selected steering sources validated before refresh writes | findSteeringLinks and main; missing/directory normal and dry-run tests | Confirmed |
| Owned files and unrelated links preserved | ensureSymlink and refresh fixture | Confirmed |
| Current reviewed source supplies selected files | clean source checkout HEAD e08d63e874f39264ed36f31812476f0816b9e49e; file metadata; real-source dry run | Confirmed |
| Same scanner distinguishes old/new Go | candidate-results.json, scanner-go12513.log, scanner-go1269.log | Confirmed from execution records; seven old-pin stdlib findings and none on candidate |
| Full candidate compatibility gates pass | postpin-results.json and postpin-go-test.log | Confirmed from records for full tests, vet/build, scan and module verification |
| Full lint is clean | Author does not claim this; lint-comparison.log and appendix distinguish capped 50 reports | Correctly limited; legacy debt remains |
| Dependency and product authority unchanged | tracked diff and go.sum unchanged | Confirmed |
| Go 1.26.9 is the maintained fixed line | official Go release history and GO-2026-6617 | Confirmed; https://go.dev/doc/devel/release and https://pkg.go.dev/vuln/GO-2026-6617 |

## Verification Performed

- Verified worktree branch, base/HEAD and dirty path scope.
- Independently executed Node 22.22.1 `node --test scripts/setup-codex-links.test.mjs`: 13 passed, zero failed.
- Independently executed `git diff --check`: passed.
- Independently ran pinned-source `AI_CENTRAL_HOME=/Users/dsteele/repos/ai-central node scripts/setup-codex-links.mjs --dry-run`: exit 0, three selected steering links reported. Source HEAD matches recorded pin and status was clean.
- Inspected exact old/new scanner result records, toolchain identity, full post-pin Go results, base JS results and all six final JS result records; verified all 32 retained command-log hashes against local raw logs. Did not rerun the full Go suite; those results are author/worker execution records rather than independent reexecution.
- Read official release history on 2026-10-15: Go 1.26.9 released 2026-10-08; Go 1.27 is current, making 1.25 unsupported under the two-newer-major policy.

## Open Questions And Residual Risks

Remote PR CI/publication is not complete and is owned by the integration task. The existing unfiltered revive debt remains; capped output is not a total or identical-findings claim. Toolchain behavior was exercised locally on macOS arm64; this is not installed-product, compiler reproducibility, guest, transport or runtime admission evidence. Pin checking remains checkout-HEAD checking rather than an immutable-byte trust boundary; that behavior predates this patch. Synthetic installer tests do not prove arbitrary upstream installer behavior.

## Verdict

**Ready.** No actionable findings remain in the bounded implementation or recovery plan. Final retained JS results all exit 0, validation.json is populated, and its command-log hashes match the inspected local outputs. README now correctly distinguishes the Go module minimum from exact reproduction. This verdict supports PR publication/review; it does not certify remote CI, merge, installed enforcement or product admission.

## Recommended Next Actions

Integration owner should record this review disposition and publish through the existing branch/PR workflow. Required remote CI must pass before merge; keep issue #217 and deferred installed/admission gates separately scoped.
