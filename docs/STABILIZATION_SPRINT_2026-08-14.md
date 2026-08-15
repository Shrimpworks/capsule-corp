# Stabilization sprint — 2026-08-14

Work item: pre-feature stabilization and cleanup
Status: `PASSED`
Scope: close the current correctness, test-runtime, dependency, documentation, and CI-governance
issues that can obscure feature work; then return to the C5b7/C5b8 no-run transport sequence.
Evidence or reason: every stabilization lane below is merged. GitHub readback confirms exact heads
and merge commits for PRs #296, #298, and #299, and the completed retrospective fresh-context
review found no PR-introduced P0-P3 findings. Its verdicts were #296 **Ready with non-blocking
follow-ups**, #298 **Ready**, and #299 **Ready**. No lane authorized a runtime, backend, VM, or
guest execution.
Remaining work: outside this passed sprint, track the pre-existing #296 included-root symlink
hardening opportunity and continue C5b7/C5b8 as separately scoped no-run tasks.
Next action: proceed with C5b7 and C5b8 under their own delivery and review gates, then integrate
them through the C5b9 immutable no-run composite.
Parent status: owner-only hostile-`.mjs` internal alpha is `IN_PROGRESS — TRENDING_GOOD`;
runtime/profile product admission remains `BLOCKED`.

## Delivery classification and integration

The brain/orchestrator owns this execution index and final reconciliation. Each numbered issue lane
is explicitly classified as an independent user-visible task: it uses its own `codex/*` branch,
worktree, commit, and pull request. Review and merge remain coordinator gates. The lanes must not
modify or consume the user's staged `.gitignore` change on `main`.

| Stable ID | Work item | Dependency | Owner / delivery unit | Acceptance and verification | Status |
| --- | --- | --- | --- | --- | --- |
| STAB-288 | Fix duplicate-instance reconciliation quarantine | none | merged PR #298 / head `5baeee4dbc5e7514399fa5e2a7a51ea951acbb3b` / merge `e3b8aa537d60e1c7783377491813d35d2acae356` | scoped implementation and local/GitHub gates pass; fresh-context verdict **Ready**, with no P0-P3 findings | `PASSED` |
| STAB-284 | Reduce the 64-segment backup race-test runtime without weakening its boundary | none | merged PR #300 / `3b8d9bf45012d5c96aa7699c8f60ec5708b240e3` | public 63→64 activation retained; race target 595.12→50.34 seconds; full package race and GitHub gates pass; no production/timeout change | `PASSED` |
| STAB-281 | Update `golang.org/x/sys` and the Biome schema only | none | merged PR #297 / `8dd85eeb970e282c326506f4d0c9d9ff8f6f7abb` | metadata-only diff; Go 1.25.13 aggregate race, vulnerability scan, and all GitHub gates pass | `PASSED` |
| STAB-267-D | Decide the external-script CI trust boundary | none | coordinator | decision recorded and closed on #267; implementation scoped in #293 | `PASSED` |
| STAB-267-I | Vendor/review the compiled-artifact verifier in Capsule and fix temporary-clone cleanup | STAB-267-D | merged PR #296 / head `79c4da365734f8dc9bb8f7d4690ab3058f13773e` / merge `294c3196d667d8017db075fbd489a042bbf1b175` | exact pinned archive and local/GitHub gates pass; fresh-context verdict **Ready with non-blocking follow-ups**, with no PR-introduced P0-P3 findings | `PASSED` |
| STAB-280 | Correct the two bounded factual documentation defects | none | merged PR #295 / `5ad7449e2f5a267fdea2a3a5d65df137a3b05d60` | generated count, docs checks, TypeScript gates, GitHub checks, and body/readback pass | `PASSED` |
| STAB-CODEQL | Synchronize CodeQL action versions and retire duplicate stale Dependabot PRs | none | merged PR #294 / `a0cb9c12a1aa60c741d302346e9daa1a0edd2ca2` | exact v4.37.7 synchronized pin, local and GitHub checks pass; stale PRs #243/#245 retired | `PASSED` |
| STAB-279 | Correct macOS installation-plan trust-boundary documentation and coverage gaps | none | merged PR #299 / head `ff60ca4818429738c651145281a1826181043aec` / merge `4ff08150b767c1c15fb75d8cdf449baf2640af99` | scoped implementation, local/GitHub gates, and internal adversarial review pass; fresh-context verdict **Ready**, with no P0-P3 findings and no canonical signing-state or authority change | `PASSED` |
| NEXT-C5B7 | Rebuild the immutable runtime root without execution | none; stabilization gate passed | separately scoped feature task | exact no-run construction evidence and canonical status update | `IN_PROGRESS — TRENDING_GOOD` |
| NEXT-C5B8 | Implement real controlled-test effects behind test doubles without guest execution | compatible C5b7 inputs; stabilization gate passed | separately scoped feature task | fail-closed test seams, retained effect evidence, and no-run verification | `IN_PROGRESS — TRENDING_GOOD` |

STAB-281 was rebased after STAB-284 and its aggregate race gate passed before merge. PRs #296,
#298, and #299 then merged at the exact identities above. Their fresh-context review completed
after merge, so the verdicts are retrospective and GitHub contains no recorded PR review for those
three merges. The review nevertheless reverified exact heads, diffs, focused and race checks,
the #296 pinned archive, and green GitHub checks without finding a PR-introduced P0-P3 issue.
C5b7/C5b8 may now proceed as separate no-run tasks; controlled C5b guest execution remains a
separate exact authorization.

## Fresh-context review closeout

| Pull request | Exact reviewed head | Exact merge commit | Verdict | Residual or limitation |
| --- | --- | --- | --- | --- |
| #296 | `79c4da365734f8dc9bb8f7d4690ab3058f13773e` | `294c3196d667d8017db075fbd489a042bbf1b175` | **Ready with non-blocking follow-ups** | The pinned historical verifier and Capsule-owned copy check descendant symlinks but do not explicitly `lstat` the two included root entries before traversal. This pre-existing limitation was neither introduced nor worsened by #296; harden it separately. |
| #298 | `5baeee4dbc5e7514399fa5e2a7a51ea951acbb3b` | `e3b8aa537d60e1c7783377491813d35d2acae356` | **Ready** | Evidence remains limited to the local fixed store and FakeBackend; no real backend or installed Supervisor claim follows. |
| #299 | `ff60ca4818429738c651145281a1826181043aec` | `4ff08150b767c1c15fb75d8cdf449baf2640af99` | **Ready** | Active-profile activation/update branches remain intentionally unreachable; any future active profile requires its own decision, fixtures, and installed evidence. |

This review is retained as sprint evidence, not as an approval or a product-admission decision.

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

## Coordinator closeout

The coordinator closeout is `PASSED` for this exact stabilization scope:

1. Every retained stabilization commit has the merged integration destination recorded above.
2. The required fresh-context review completed with no PR-introduced P0-P3 findings and the exact
   verdicts and limitations retained in this index.
3. Required checks, pull-request metadata, exact heads, merge commits, and merged state were read
   back from GitHub.
4. The integration verification snapshot above remains the resulting `main` evidence for the
   sprint; this closeout changes documentation only.
5. The sprint is `PASSED`; the owner-only internal-alpha parent remains
   `IN_PROGRESS — TRENDING_GOOD`, and product/runtime-profile admission remains `BLOCKED`.
