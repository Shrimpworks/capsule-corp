# G3a5 canonical verification handoff

Status: PASSED for required canonical baseline CI checks. Existing unfiltered lint debt #217 remains BLOCKED. Parent full C5b19 Gate 3 remains IN_PROGRESS — TRENDING_GOOD; child execution and later gates remain BLOCKED.

## Scope and integration ownership

Internal verification-only assignment. Canonical worktree `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-event-fields`, HEAD/base `fe3d22c84b5d41e0a4edbfacb531563040e16995` (merged PR402). No source, metadata or dependency edits, branch, commit, push or PR were made by this worker. Existing `.serena/project.yml` modifications and root-owned `.planning/c5b19-event-fields/` were preserved. Root owns integration, review and publishing.

Defensive scope was the owned repository's standard canonical automated checks. No C5b19 native fixture, child, guest, signal/wait target, live store or installed service was run. The archive's frozen Go1.25.13 packet was unchanged. This assignment does not establish any Gate3 implementation, independent review, hosted CI or later admission outcome.

## Environment and commands

Observed Node22.23.3, pnpm10.28.2, exact Go1.26.11 darwin/arm64, golangci-lint2.12.2. Node's observed patch differs from documented22.22.1 but satisfies package engines>=22.22.0; no toolchain upgrade was attempted.

All commands used the canonical worktree. The exact directory `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.11.darwin-arm64/bin` was prepended to PATH; `GOTOOLCHAIN=local` was set. Make's scanner recipe explicitly overrides this to `GOTOOLCHAIN=auto`, while using that exact Go binary and the unchanged module1.26.11 requirement.

| Command | Result |
| --- | --- |
| `pnpm install --frozen-lockfile` | exit0; locked dependencies unchanged |
| `make ci` | exit0 |
| `pnpm test` | exit0 |
| `golangci-lint run ./...` | exit1; known revive backlog,50 displayed findings under output cap |
| `git diff --check` | exit0 |

The Makefile covers required Go formatting/vet, full correctness plus changed-line revive lint, pinned govulncheck1.6.0, dependency audit, site build, schema/ADR checks, TypeScript check/lint, changed JSDoc, full `go test ./...`, `go build ./...`, and built TS/script tests. Standalone `pnpm test` was also run. Results include177 workspace tests (164protocol+10SDK+3MCP) and132 script tests.

Selected observed output:

```text
0 issues.
0 issues.
No vulnerabilities found.
4 vulnerabilities found
Severity: 2 low | 2 moderate
Checked 162 files in 101ms. No fixes applied.
```

The two zero-issue lines are scoped Go lint. The vulnerability-free line is govulncheck. The4 lower-severity findings are the dependency audit, whose existing high threshold passed; this is not a claim that dependencies have no vulnerabilities. Full unfiltered lint reported50revive findings; this is capped output, not a complete debt count. No gate threshold or exclusion changed.

## Logs and retained evidence

Commands, versions, exit codes and SHA256 hashes are in `/tmp/c5b19-g3a5-baseline.json`.

- Install: `/tmp/c5b19-g3a5-install.log`, SHA256 `8bb7491be1262e14400b996750fb575332df0d2ed37ba8d5a681491b16c45956`.
- CI: `/tmp/c5b19-g3a5-ci.log`, SHA256 `6b2a178e7a9f48b74a2f89d78ee57ad59bbb6b723d82f55fdcfba9a3005fd963`.
- Standalone test: `/tmp/c5b19-g3a5-pnpm-test.log`, SHA256 `5a3d9cdd4034c105727b8da3c348e9579c1c76bd925666c62ca3b0c72305e94d`.
- Full lint: `/tmp/c5b19-g3a5-lint-full.log`, SHA256 `f2e2396820d2e5b7e74cc8d3234b9c3e643e0d1dbf251e4786729bb78a0075e0`.

Raw /tmp logs are ephemeral; root should retain this handoff and JSON in its canonical evidence. Current advisory DB was used on2026-11-13 America/Toronto (UTC2026-11-14), not a pinned release DB snapshot.

## Final metadata lint integration dependency

At baseline completion the new planning directory contained only untracked `plan.md`, so this baseline cannot certify future retained JSON formatting. Root explicitly took ownership of formatting and explicit Biome validation after retaining final JSON, to avoid a circular report update. Root must run `pnpm exec biome check <exact retained JSON paths>` or lint after files are tracked, then record that distinct final integration result in the PR body. No claim of that future check is made here.

## Retrieval and limitations

The consulted current Makefile/configs were read directly in the exact target worktree. Codebase-memory primary index generation2026-10-04T22:18:10Z reported best-effort metadata_match for go.mod/Makefile/.gitignore, and biome.json not_tracked; direct source reads supply the authoritative fallback. No structural code query or repository-wide completeness claim was needed.

All owned checks completed; no process session remains. Recommended next action: retain these results, finish the final metadata lint check after staging/explicit selection, obtain independent implementation review, and publish with hosted checks while preserving child execution gates.
