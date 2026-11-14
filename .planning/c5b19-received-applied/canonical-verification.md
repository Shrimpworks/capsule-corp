# G3a7 canonical verification handoff

Status: **PASSED** for required canonical baseline CI checks. Existing unfiltered lint debt #217 remains **BLOCKED**. Parent full C5b19 Gate 3 remains **IN_PROGRESS — TRENDING_GOOD**; child execution and later gates remain **BLOCKED**.

## Scope and integration ownership

Internal verification-only assignment. Canonical worktree `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-received-applied`, HEAD/base `041b6d4e04df9667e6217f57df26f04eec7d2ddc` (merged PR #404). This worker made no repository edits, branch, commit, push, PR or subagent. Existing `.serena/project.yml` changes and Brain-owned `.planning/c5b19-received-applied/` were preserved. Brain owns evidence retention, final integration, review and publication.

Defensive authorized scope was the owned repository's standard automated checks. No C5b19 native fixture, child, guest, lifecycle/signal/wait target, live store, installed service, credential or user data was exercised. Archive Go 1.25.13 and run/admission gates remained unchanged. This baseline does not establish correctness or independent review of G3a7 archive behavior, hosted CI or later admission outcomes.

## Environment and verification

Observed Node 22.23.3, pnpm 10.28.2, exact Go 1.26.11 darwin/arm64, and golangci-lint 2.12.2 (built with Go 1.26.2). The local Node patch differs from documented 22.22.1; no upgrade was attempted. `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.11.darwin-arm64/bin` was prepended to PATH with `GOTOOLCHAIN=local`. Make's scanner recipe explicitly sets `GOTOOLCHAIN=auto`, invoking that exact Go binary with unchanged module 1.26.11 requirement.

| Command | Result |
| --- | --- |
| `pnpm install --frozen-lockfile` | exit 0; lockfile unchanged |
| `make ci` | exit 0 |
| `pnpm test` | exit 0 |
| `golangci-lint run ./...` | exit 1; known #217 revive backlog, 50 displayed findings under output cap |
| `git diff --check` | exit 0 |

The current Makefile covers Go formatting/vet, full correctness plus changed-line revive lint, pinned govulncheck 1.6.0, dependency audit, site build, schema/ADR checks, TypeScript check/lint, changed JSDoc, full `go test ./...`, `go build ./...`, and built TS/script tests. Standalone `pnpm test` also passed. Each test invocation reports 177 workspace tests (164 protocol + 10 SDK + 3 MCP) and 132 script tests, with zero failures. Go package tests passed; some were cached as shown in the log.

Selected observed output:

```text
0 issues.
0 issues.
No vulnerabilities found.
4 vulnerabilities found
Severity: 2 low | 2 moderate
Checked 164 files in 79ms. No fixes applied.
```

The zero-issue lines are scoped Go lint; govulncheck found no vulnerabilities. Dependency audit retained two low/two moderate advisories below the unchanged high threshold, not a vulnerability-free dependency claim. Full unfiltered lint displayed 50 revive findings under its output cap; this is not a total backlog count. No gate threshold, exclusion or branch protection changed.

## Logs and retained evidence

Exact commands, versions, UTC timestamps, exit codes and log hashes are in `/tmp/c5b19-g3a7-baseline.json`. Commands ran 2026-11-14 01:36–01:37 UTC (2026-11-13 America/Toronto). Current advisory data was queried, not pinned as a release snapshot. All four log hashes were checked before this handoff.

- `/tmp/c5b19-g3a7-install.log`: SHA-256 `cedd53b4756cc62087fae8555bc34b90dd0493a9a0601c51db934cb61f153718`.
- `/tmp/c5b19-g3a7-ci.log`: SHA-256 `cd6c7fdede62915d46f6ab0f767c4397f2eacf4fffc4f0e5a8cf273a16e93543`.
- `/tmp/c5b19-g3a7-pnpm-test.log`: SHA-256 `be60b3c9d720402fd870d76b4cbfe99408cef6c9c4459d9ea5a4d15973057817`.
- `/tmp/c5b19-g3a7-lint-full.log`: SHA-256 `bd76c25c0d4e5f0f3fa9236a46d6f0bffadd5e7c8d1eb6afcb469654784eaa26`.

Raw `/tmp` logs are ephemeral. Brain should retain this handoff and JSON in canonical evidence. All owned command processes completed; no process session remains active.

## Final integration dependencies and limitations

This baseline cannot certify future retained JSON formatting or documentation written after these checks. Brain owns explicit Biome validation of retained JSON and final site/lint/diff checks after metadata reconciliation. No future result is claimed here.

Codebase-memory index status was checked: primary checkout generation `2026-10-04T22:18:10Z`, not watched, with older documentation-link metadata. Exact-path coverage reported best-effort metadata_match for AGENTS.md/Makefile/go.mod, not_tracked for biome.json, and missing freshness for the prior verification handoff. All consulted files were read directly in the exact worktree, providing authoritative fallback. No structural exploration or repository-wide completeness claim was needed.

Confidence is high for these exact local command outcomes. Recommended next action: retain evidence, finish final metadata gates, obtain independent implementation review, and publish the bounded slice while keeping child execution blocked. No new failure or unresolved environment issue was found beyond the disclosed baseline limitations.
