# G3a12 canonical baseline verification

Status: **PASSED** for required canonical baseline checks. Existing unfiltered lint debt #217 remains **BLOCKED**. Parent full Gate3 remains **IN_PROGRESS — TRENDING_GOOD**; child execution, Gates4–6 and product admission remain **BLOCKED**.

## Scope and ownership

Internal verification-only G3a12-VERIFY assignment in `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-transport-lanes`, HEAD/base `f1ded569af4a6ae20f2a93b7e476825959f97143` (merged PR409). The question tested was whether the canonical baseline passes its required local checks before bounded supplied transport lane consistency delivery. This worker made no repository edits, branch, commit, push, PR or subagent. Brain owns retention, final metadata checks, independent review, integration and shared PR publication. Concurrent planning and unrelated `.serena/project.yml` changes were preserved.

Defensive authorized scope was ordinary automated checks of this owned local repository. No C5b19 native fixture, child, guest, lifecycle/signal/wait target, live store, installed service, credential or user data was exercised. Existing routine advisory queries were used by the scanner/audit. Archive Go1.25.13 and run/admission gates were not changed. This baseline does not establish correctness or review of the future G3a12 archive increment, hosted acceptance or product admission.

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

Checks started `2026-11-14T04:20:45.703726+00:00` and completed `2026-11-14T04:21:41.570383+00:00`. Exact versions, commands, timestamps, exits and hashes are in `/tmp/c5b19-g3a12-baseline.json`.

- `/tmp/c5b19-g3a12-install.log`: SHA256 `0c0fe3efd5cf2e2244cfefdeffef78af189b462bc3abeef549f6a42b44fee0e4`.
- `/tmp/c5b19-g3a12-ci.log`: SHA256 `ace991c83e982dd2b1ff25953a8f3743936be30d250cb45e000dca95c5437794`.
- `/tmp/c5b19-g3a12-pnpm-test.log`: SHA256 `c8097e778068de7689a5aff7cbb24a09b45b6bab5532dfee8768bccc442b9c11`.
- `/tmp/c5b19-g3a12-lint-full.log`: SHA256 `b12416c706c0ba0f921ded07b90baf9c6d078c2d02ae8dcd36461e468f078425`.
- `/tmp/c5b19-g3a12-diff-check.log`: SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.

All five log hashes were verified before handoff; all owned processes completed. Raw `/tmp` logs are ephemeral. Brain should retain this handoff and validation JSON in canonical planning. Future JSON formatting, later documentation edits and final hosted checks are not certified here; final integration metadata validation remains Brain-owned.

Codebase-memory Tier2 index_status confirms the primary checkout index, ready generation2026-10-04T22:18:10Z, old/not watched and not this worktree. Coverage checked all five relied-on paths: AGENTS.md, Makefile and go.mod reported best-effort metadata_match; biome.json was not_tracked; the G3a12 plan had missing freshness. Exact target-worktree reads supplied current configuration and planning evidence. No structural exploration or whole-graph completeness claim was needed. Serena remained initialized in this worker session. Current target-worktree Makefile, go.mod and plan were read. Previously read target-worktree AGENTS and biome.json were verified unchanged; all four configuration files have no diff from the prior verified base d042e53a877e494bd2eb0e325ec084d4b9cc692a or worktree modifications.

Confidence is high for the exact observed command outcomes. No new failure was found beyond disclosed baseline lint/audit limitations. Recommended next action: Brain retain evidence, complete G3a12 implementation and fresh review, run final metadata checks, and publish through the shared destinations while keeping child execution blocked. There are no unresolved verification questions for this bounded baseline; final implementation/metadata/hosted outcomes remain deferred.
