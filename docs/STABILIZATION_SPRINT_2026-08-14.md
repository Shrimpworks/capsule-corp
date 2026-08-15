# Stabilization sprint — 2026-08-14

Work item: pre-feature stabilization and cleanup
Status: `IN_PROGRESS — TRENDING_GOOD`
Scope: close the current correctness, test-runtime, dependency, documentation, and CI-governance
issues that can obscure feature work; then return to the C5b7/C5b8 no-run transport sequence.
Evidence or reason: the repository is synchronized at `ebdbf6847c2baa5a049b58f230706d9c079ec69f`,
the local AI Central integration is current, and the open-issue review identified separable bounded
lanes. No lane authorizes a runtime, backend, VM, or guest execution.
Remaining work: deliver and review the active lanes below, land STAB-284 before completing the
STAB-281 aggregate race gate, and update canonical status after GitHub readback.
Next action: finish STAB-288, STAB-284, and STAB-279 in parallel; obtain a genuinely fresh-context
review for the security-sensitive STAB-267-I and STAB-288 changes.
Parent status: owner-only hostile-`.mjs` internal alpha is `IN_PROGRESS — TRENDING_GOOD`;
runtime/profile product admission remains `BLOCKED`.

## Delivery classification and integration

The brain/orchestrator owns this execution index and final reconciliation. Each numbered issue lane
is explicitly classified as an independent user-visible task: it uses its own `codex/*` branch,
worktree, commit, and pull request. Review and merge remain coordinator gates. The lanes must not
modify or consume the user's staged `.gitignore` change on `main`.

| Stable ID | Work item | Dependency | Owner / delivery unit | Acceptance and verification | Status |
| --- | --- | --- | --- | --- | --- |
| STAB-288 | Fix duplicate-instance reconciliation quarantine | fresh-context review | independent task / ready PR #298 at `5baeee4dbc5e7514399fa5e2a7a51ea951acbb3b` | scoped implementation plus local and GitHub gates pass; fresh-context review pending | `IN_PROGRESS — TRENDING_GOOD` |
| STAB-284 | Reduce the 64-segment backup race-test runtime without weakening its boundary | none | merged PR #300 / `3b8d9bf45012d5c96aa7699c8f60ec5708b240e3` | public 63→64 activation retained; race target 595.12→50.34 seconds; full package race and GitHub gates pass; no production/timeout change | `PASSED` |
| STAB-281 | Update `golang.org/x/sys` and the Biome schema only | none | merged PR #297 / `8dd85eeb970e282c326506f4d0c9d9ff8f6f7abb` | metadata-only diff; Go 1.25.13 aggregate race, vulnerability scan, and all GitHub gates pass | `PASSED` |
| STAB-267-D | Decide the external-script CI trust boundary | none | coordinator | decision recorded and closed on #267; implementation scoped in #293 | `PASSED` |
| STAB-267-I | Vendor/review the compiled-artifact verifier in Capsule and fix temporary-clone cleanup | STAB-267-D; fresh-context review | independent task / ready PR #296 | local verifier executes against pinned archive data; local and GitHub checks pass; fresh-context review pending | `IN_PROGRESS — TRENDING_GOOD` |
| STAB-280 | Correct the two bounded factual documentation defects | none | merged PR #295 / `5ad7449e2f5a267fdea2a3a5d65df137a3b05d60` | generated count, docs checks, TypeScript gates, GitHub checks, and body/readback pass | `PASSED` |
| STAB-CODEQL | Synchronize CodeQL action versions and retire duplicate stale Dependabot PRs | none | merged PR #294 / `a0cb9c12a1aa60c741d302346e9daa1a0edd2ca2` | exact v4.37.7 synchronized pin, local and GitHub checks pass; stale PRs #243/#245 retired | `PASSED` |
| STAB-279 | Correct macOS installation-plan trust-boundary documentation and coverage gaps | fresh-context review | independent task / ready PR #299 at `ff60ca4818429738c651145281a1826181043aec` | scoped implementation plus local/GitHub gates and internal adversarial review pass; fresh-context review pending; no canonical signing-state or authority change | `IN_PROGRESS — TRENDING_GOOD` |
| NEXT-C5B7 | Rebuild the immutable runtime root without execution | stabilization merge gate | separately scoped feature task | exact no-run construction evidence and canonical status update | `BLOCKED` |
| NEXT-C5B8 | Implement real controlled-test effects behind test doubles without guest execution | stabilization merge gate; compatible C5b7 inputs | separately scoped feature task | fail-closed test seams, retained effect evidence, and no-run verification | `BLOCKED` |

STAB-267-I is implemented in #296 but remains unmerged until its required fresh-context review.
STAB-281 was rebased after STAB-284 and its aggregate race gate passed before merge. C5b7/C5b8
resume after this sprint's remaining review/merge gates; controlled C5b guest execution remains a
separate exact authorization.

## Integration verification snapshot

The coordinator verified the merged stabilization baseline plus this documentation-only index with:

- `pnpm install --frozen-lockfile`, `pnpm check`, `pnpm lint`, `pnpm test`,
  `pnpm verify:schemas`, and `pnpm verify:adrs` — passed;
- `go test ./...`, `go vet ./...`, and `go build ./...` — passed outside the filesystem sandbox so
  the existing daemon loopback tests could bind locally;
- `GOTOOLCHAIN=go1.25.13 go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...` — no
  vulnerabilities found;
- `golangci-lint run ./...` — blocked only by the same 50 pre-existing `revive` exported-comment
  findings tracked in #217; the delivered STAB-284 changed-code lint reported zero findings;
- STAB-281's post-STAB-284 aggregate `go test ./... -race -count=1` — passed, with the
  registration-state package completing in 557.308 seconds under the unchanged timeout.

Running the vulnerability scanner with the host's unpinned Go 1.26.5 instead reports three standard
library advisories fixed in Go 1.26.6; that host-toolchain limitation does not apply to the exact
repository Go 1.25.13 gate above and is not being hidden or widened into this sprint.

## Coordinator merge gate

Before this sprint is `PASSED`, the coordinator must:

1. collect each independent task's branch, commit, pull request, verification, limitations, and
   deferred work;
2. obtain fresh-context review for security-sensitive behavior changes;
3. confirm required checks and pull-request bodies by GitHub readback;
4. merge or explicitly document the integration destination for every retained commit;
5. rerun proportionate integration verification on the resulting `main`;
6. update this index and any affected canonical plan/status documents without overstating product
   admission.
