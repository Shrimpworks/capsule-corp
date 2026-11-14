# G3a6 canonical verification handoff

Status: **PASSED** for required canonical baseline CI checks. Existing unfiltered lint debt #217 remains **BLOCKED**. Parent full C5b19 Gate 3 remains **IN_PROGRESS — TRENDING_GOOD**; child execution and later gates remain **BLOCKED**.

## Scope and integration ownership

Internal verification-only assignment. Canonical worktree `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-frame-sidecars`, HEAD/base `e1a966e7c6f18665784984e97b6697cc472cb288` (merged PR #403). The worker made no source, metadata or dependency edits, branch, commit, push or PR. Existing `.serena/project.yml` changes and Brain-owned `.planning/c5b19-frame-sidecars/` were preserved. Brain owns final integration, evidence retention, review and publishing.

Defensive authorized scope was the owned repository's standard automated checks. No C5b19 native fixture, child, guest, signal/wait target, live store, installed service or user data was exercised. Archive frozen Go 1.25.13 and run/admission gates were unchanged. This baseline does not establish implementation correctness, independent review, hosted CI or later admission outcomes for G3a6.

## Environment and verification

Observed Node 22.23.3, pnpm 10.28.2, exact Go 1.26.11 darwin/arm64 and golangci-lint 2.12.2 (built with Go 1.26.2). The local Node patch differs from documented 22.22.1; no upgrade was attempted. The exact installed directory `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.11.darwin-arm64/bin` was prepended to PATH with `GOTOOLCHAIN=local`. Make's scanner recipe explicitly sets `GOTOOLCHAIN=auto`, while invoking that exact Go binary and the unchanged module 1.26.11 requirement.

| Command | Result |
| --- | --- |
| `pnpm install --frozen-lockfile` | exit 0, lockfile unchanged |
| `make ci` | exit 0 |
| `pnpm test` | exit 0 |
| `golangci-lint run ./...` | exit 1, known #217 revive backlog; 50 displayed findings under output cap |
| `git diff --check` | exit 0 |

The current Makefile runs Go formatting/vet, full correctness plus changed-line revive lint, pinned govulncheck 1.6.0, dependency audit, site build, schema/ADR checks, TypeScript check/lint, changed JSDoc, full `go test ./...`, `go build ./...`, and built TS/script tests. Standalone `pnpm test` also passed. Both test invocations report 177 workspace tests (164 protocol + 10 SDK + 3 MCP) and 132 script tests, with no failures. Go package tests passed; some were cached, as shown in the retained log.

Selected observed output:

```text
0 issues.
0 issues.
No vulnerabilities found.
4 vulnerabilities found
Severity: 2 low | 2 moderate
Checked 163 files in 106ms. No fixes applied.
```

The zero-issue lines are scoped Go lint. Govulncheck found no vulnerabilities. The dependency audit retained two low and two moderate advisories below its unchanged high threshold, so this is not a claim that dependencies are vulnerability-free. Unfiltered lint displayed 50 revive findings under its output cap; that is not a total debt count. No gate threshold, exclusion or branch protection changed.

## Logs and retained evidence

Exact commands, versions, UTC timestamps, exit codes and log hashes are in `/tmp/c5b19-g3a6-baseline.json`. Commands ran 2026-11-14 01:17–01:18 UTC (2026-11-13 America/Toronto). The current advisory database was queried; its contents were not pinned as a release snapshot.

- `/tmp/c5b19-g3a6-install.log`: SHA-256 `6fc261a7f154e871facf2248a6dc160978c1824cb744a2a8067f999d9a3ca589`.
- `/tmp/c5b19-g3a6-ci.log`: SHA-256 `8be51156fe3c7654fc3aa969d736ea59822ac26c907a1570452ca9382925b41f`.
- `/tmp/c5b19-g3a6-pnpm-test.log`: SHA-256 `d6e7c1f4501b432bcc390766c7b90ba7835ccfc941c1d2dac666ed222a941baa`.
- `/tmp/c5b19-g3a6-lint-full.log`: SHA-256 `44829f3e16725109ca6748363d744259b4d8e6d93304201888a5071c824bc519`.

Raw `/tmp` logs are ephemeral. Brain should retain this handoff and JSON in canonical evidence. No process session remains running.

## Final integration dependencies and limitations

This baseline cannot certify formatting of future retained JSON or documentation written after these checks. Brain owns explicit Biome checking of retained JSON and final site/lint/diff checks after metadata reconciliation. No future check is claimed as passed here.

Codebase-memory index status was checked: primary checkout generation `2026-10-04T22:18:10Z`, not watched, with older documentation-link metadata. Coverage for Makefile/go.mod reported best-effort `metadata_match`; biome.json was `not_tracked` and the prior verification report had missing freshness. All consulted files were read directly in the exact target worktree, providing the authoritative fallback. No structural code query or repository-wide completeness claim was needed.

Confidence is high for these exact local command outcomes. The new G3a6 archive implementation, final canonical prose, review acceptance, publication and hosted results remain separate Brain-owned gates. Recommended next action: retain the evidence, finish metadata checks, complete independent review and publish the bounded slice while keeping child execution blocked.
