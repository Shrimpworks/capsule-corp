# G3a11 canonical baseline verification

Status: **PASSED** for required canonical baseline checks. Existing unfiltered lint debt #217 remains **BLOCKED**. Parent full Gate3 remains **IN_PROGRESS — TRENDING_GOOD**; child execution, Gates4–6 and product admission remain **BLOCKED**.

## Scope and ownership

Internal verification-only G3a11-VERIFY assignment in `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-write-results`, HEAD/base `d042e53a877e494bd2eb0e325ec084d4b9cc692a` (merged PR408). The question tested was whether the canonical baseline passes its required local checks before bounded supplied write/result correspondence delivery. This worker made no repository edits, branch, commit, push, PR or subagent. Brain owns retention, final metadata checks, independent review, integration and shared PR publication. Concurrent planning and unrelated `.serena/project.yml` changes were preserved.

Defensive authorized scope was ordinary automated checks of this owned local repository. No C5b19 native fixture, child, guest, lifecycle/signal/wait target, live store, installed service, credential or user data was exercised. Existing routine advisory queries were used by the scanner/audit. Archive Go1.25.13 and run/admission gates were not changed. This baseline does not establish correctness or review of the future G3a11 archive increment, hosted acceptance or product admission.

## Verification and environment

Exact Go1.26.11 darwin/arm64 was selected by prepending `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.11.darwin-arm64/bin` to PATH with `GOTOOLCHAIN=local`. Make's scanner recipe explicitly uses `GOTOOLCHAIN=auto`, invoking that same exact Go binary with the unchanged module1.26.11 requirement. Observed Node22.23.3, pnpm10.28.2 and golangci-lint2.12.2 built with Go1.26.2. Local Node differs from the documented22.22.1 patch; no upgrade was attempted.

| Command | Result |
| --- | --- |
| `pnpm install --frozen-lockfile` | exit0 |
| `make ci` | exit0 |
| `pnpm test` | exit0 |
| `golangci-lint run ./...` | exit1; known #217 revive debt, 50 displayed under output cap |
| `git diff --check` | exit0 |

Make's CI covers Go formatting/vet, correctness and changed-line revive lint, pinned govulncheck1.6.0, dependency audit, documentation build, schema/ADR verification, TypeScript checks/lint, changed JSDoc, full Go tests/build, and built workspace/script tests. Both CI and standalone tests report177 workspace tests (164protocol +10SDK +3MCP) and132 script tests, all passing. Go packages passed, with cached package results disclosed in the raw log. No lockfile, go.mod or go.sum diff was present after verification.

Observed scanner output: `No vulnerabilities found.` Dependency audit separately reports four findings: two low and two moderate, below the unchanged high threshold. This is not a claim that dependencies have no vulnerabilities. Scoped Go lint reports zero issues; unfiltered lint shows50 revive findings under its output cap, not a total debt count. The installed binary's help confirms `--max-issues-per-linter` defaults to50 and `--max-same-issues` defaults to3. No check, threshold, exclusion or branch protection was weakened.

## Evidence and limits

Checks started `2026-11-14T03:57:53.250047+00:00` and completed `2026-11-14T03:58:56.129713+00:00`. Exact versions, commands, timestamps, exits and hashes are in `/tmp/c5b19-g3a11-baseline.json`.

- `/tmp/c5b19-g3a11-install.log`: SHA256 `b060e90e5ec372a9a236ba107f313eee32fbd2f55eedcfbfefed7fa8f8d5a3b1`.
- `/tmp/c5b19-g3a11-ci.log`: SHA256 `3eeba71dac97b94e986356d186c985c9a3fd02c717447cf424bc1095c330f507`.
- `/tmp/c5b19-g3a11-pnpm-test.log`: SHA256 `6f92b92b444e2fcf902166ce7ec814b57fe7d3d5b1e2ca6e0b6d0c287fa8e846`.
- `/tmp/c5b19-g3a11-lint-full.log`: SHA256 `88584e1b606d6b3a2ac30421370ffda306685a095e4824a9f25ca95a527056c6`.
- `/tmp/c5b19-g3a11-diff-check.log`: SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.

All five log hashes were verified before handoff; all owned processes completed. Raw `/tmp` logs are ephemeral. Brain should retain this handoff and validation JSON in canonical planning. Future JSON formatting, later documentation edits and final hosted checks are not certified here; final integration metadata validation remains Brain-owned.

Codebase-memory Tier2 index_status confirms the primary checkout index, ready generation2026-10-04T22:18:10Z, old/not watched and not this worktree. Coverage checked all five relied-on paths: AGENTS.md, Makefile and go.mod reported best-effort metadata_match; biome.json was not_tracked; the G3a11 plan had missing freshness. Exact target-worktree reads supplied current configuration and planning evidence. No structural exploration or whole-graph completeness claim was needed. Serena remained initialized in this worker session. Current target-worktree AGENTS, Makefile, go.mod, biome.json and plan were read; the four configuration files have no diff from the previously verified base.

Confidence is high for the exact observed command outcomes. No new failure was found beyond disclosed baseline lint/audit limitations. Recommended next action: Brain retain evidence, complete G3a11 implementation and fresh review, run final metadata checks, and publish through the shared destinations while keeping child execution blocked. There are no unresolved verification questions for this bounded baseline; final implementation/metadata/hosted outcomes remain deferred.
