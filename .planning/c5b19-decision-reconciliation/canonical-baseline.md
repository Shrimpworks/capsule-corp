# G3a8 canonical verification handoff

Status: **PASSED** for required canonical baseline checks. Existing unfiltered lint debt #217 remains **BLOCKED**. Parent full Gate 3 remains **IN_PROGRESS — TRENDING_GOOD**; child execution and admission remain **BLOCKED**.

## Scope and integration ownership

Internal verification-only assignment. Canonical worktree `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-decision-reconciliation`, HEAD/base `8e4fd22cc00cd5e961544f6a803974600dc0e751` (merged PR #405). This worker made no repository edits, branch, commit, push, PR or subagent. Existing `.serena/project.yml` changes and Brain-owned `.planning/c5b19-decision-reconciliation/` were preserved. Brain owns evidence retention, final metadata checks, review and publication.

Defensive authorized scope was the owned repository's standard automated checks. No C5b19 native fixture, child, guest, lifecycle/signal/wait target, live store, installed service, credential or user data was exercised. Archive Go 1.25.13 and run/admission gates were unchanged. This baseline does not establish G3a8 archive implementation correctness, independent review, hosted CI or later admission outcomes.

## Environment and verification

Observed Node 22.23.3, pnpm 10.28.2, exact Go 1.26.11 darwin/arm64, and golangci-lint 2.12.2 (built with Go 1.26.2). The local Node patch differs from documented 22.22.1; no upgrade was attempted. The installed directory `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.11.darwin-arm64/bin` was prepended to PATH with `GOTOOLCHAIN=local`. Make's scanner recipe explicitly sets `GOTOOLCHAIN=auto`, invoking that exact Go binary with unchanged module 1.26.11 requirement.

| Command | Result |
| --- | --- |
| `pnpm install --frozen-lockfile` | exit0; lockfile unchanged |
| `make ci` | exit0 |
| `pnpm test` | exit0 |
| `golangci-lint run ./...` | exit1; known #217 revive debt, 50 displayed findings under output cap |
| `git diff --check` | exit0 |

The current Makefile covers Go formatting/vet, full correctness plus changed-line revive lint, pinned govulncheck 1.6.0, dependency audit, site build, schema/ADR checks, TypeScript check/lint, changed JSDoc, full `go test ./...`, `go build ./...`, and built TS/script tests. Standalone `pnpm test` also passed. Both test invocations report 177 workspace tests (164 protocol + 10 SDK + 3 MCP) and 132 script tests, all passing. Go package tests passed; some were cached as disclosed in the logs. No lockfile, go.mod or go.sum diff was present after verification.

Selected observed output:

```text
0 issues.
0 issues.
No vulnerabilities found.
4 vulnerabilities found
Severity: 2 low | 2 moderate
Checked 165 files in 111ms. No fixes applied.
```

The zero-issue lines are scoped Go lint. Govulncheck reports no vulnerabilities; dependency audit retains two low/two moderate advisories below its unchanged high threshold. This is not a vulnerability-free dependency claim. Full unfiltered lint displays 50 revive findings under its output cap, not a total backlog count. No threshold, exclusion, branch protection or test was changed.

## Logs and retained evidence

Exact commands, versions, UTC timestamps, exit codes and hashes are in `/tmp/c5b19-g3a8-baseline.json`. Started `2026-11-14T01:54:47.019345+00:00`, completed `2026-11-14T01:55:50.733215+00:00` (2026-11-13 America/Toronto). Current advisory data was queried, not pinned as a release snapshot. All four log hashes were verified before handoff.

- `/tmp/c5b19-g3a8-install.log`: SHA-256 `521e0578089590f404060102c99a888a9593782a82e65a60eafdacb65477e196`.
- `/tmp/c5b19-g3a8-ci.log`: SHA-256 `69c67fad16c6dc4d62de7481b10f72de49c574a50836dcd10962a8964e759925`.
- `/tmp/c5b19-g3a8-pnpm-test.log`: SHA-256 `922cf53df26db208da42cebe2dd6fe6cf0e560a2c2cdaaacec501aacbeea4892`.
- `/tmp/c5b19-g3a8-lint-full.log`: SHA-256 `17e0ef3071e79b780b42efa83bfee696c5118cdfcb8def07dc73a876564e9433`.

Raw `/tmp` logs are ephemeral. Brain should retain this handoff and JSON in canonical evidence. All owned command processes completed; no active process session remains.

## Final integration dependencies and limitations

Future retained JSON formatting and documentation written after these checks are not certified here. Brain owns explicit Biome validation of retained JSON and final site/lint/diff checks after reconciliation. No future check is claimed as passed.

Codebase-memory index status was checked: primary checkout generation `2026-10-04T22:18:10Z`, not watched, with older documentation-link metadata. Coverage of five consulted paths reported best-effort metadata_match for AGENTS.md/Makefile/go.mod, not_tracked for biome.json, and missing freshness for the prior verification handoff. All were read directly in the exact target worktree, supplying authoritative fallback. No structural exploration or repository-wide completeness claim was needed.

Confidence is high for these exact observed command outcomes. No new failure or unresolved environment issue was found beyond disclosed baseline limitations. Recommended next action: retain evidence, complete final metadata checks and independent implementation review, and publish the bounded slice while keeping child execution blocked.
