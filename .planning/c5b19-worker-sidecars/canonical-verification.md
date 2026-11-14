# G3a10 canonical baseline verification

Status: **PASSED** for required canonical baseline checks. Existing unfiltered lint debt #217 remains **BLOCKED**. Parent full Gate3 remains **IN_PROGRESS — TRENDING_GOOD**; child execution, Gates4–6 and product admission remain **BLOCKED**.

## Scope and ownership

Internal verification-only G3a10-VERIFY assignment in `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-worker-sidecars`, HEAD/base `fe9009d5398a883735219e2a3a33c351c79aa54c` (merged PR407). The question tested was whether the canonical baseline passes its required local checks before bounded supplied worker-received/sent sidecar correspondence delivery. This worker made no repository edits, branch, commit, push, PR or subagent. Brain owns retention, final metadata checks, independent review, integration and shared PR publication. Concurrent planning and unrelated `.serena/project.yml` changes were preserved.

Defensive authorized scope was ordinary automated checks of this owned local repository. No C5b19 native fixture, child, guest, lifecycle/signal/wait target, live store, installed service, credential or user data was exercised. Existing routine advisory queries were used by the scanner/audit. Archive Go1.25.13 and run/admission gates were not changed. This baseline does not establish correctness or review of the future G3a10 archive increment, hosted acceptance or product admission.

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

Checks started `2026-11-14T02:58:53.971542+00:00` and completed `2026-11-14T02:59:56.323703+00:00`. Exact versions, commands, timestamps, exits and hashes are in `/tmp/c5b19-g3a10-baseline.json`.

- `/tmp/c5b19-g3a10-install.log`: SHA256 `3b31b43edef06a35b5a169a7ca5c7a567d778cf8a92835ad8c3218f50302bf58`.
- `/tmp/c5b19-g3a10-ci.log`: SHA256 `9d8f61c3e52206490cdeb9877fba990760fe28945fbf85a0944dcfd17f6653fd`.
- `/tmp/c5b19-g3a10-pnpm-test.log`: SHA256 `cded3059139f59fcdd63ef5960d26d76cf7c97ec0e0805174c9d1f195fa3494b`.
- `/tmp/c5b19-g3a10-lint-full.log`: SHA256 `b3b3311d71127c02eb0c0fd9915b2fd8a7b4a108cb87f70b9fe3ac350d0fec4f`.
- `/tmp/c5b19-g3a10-diff-check.log`: SHA256 `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`.

All five log hashes were verified before handoff; all owned processes completed. Raw `/tmp` logs are ephemeral. Brain should retain this handoff and validation JSON in canonical planning. Future JSON formatting, later documentation edits and final hosted checks are not certified here; final integration metadata validation remains Brain-owned.

Codebase-memory Tier2 index_status confirms the primary checkout index, ready generation2026-10-04T22:18:10Z, old/not watched and not this worktree. Coverage checked all six relied-on paths: AGENTS.md, Makefile and go.mod reported best-effort metadata_match; biome.json was not_tracked; both consulted planning reports had missing freshness. Exact target-worktree reads supplied current configuration and planning evidence. No structural exploration or whole-graph completeness claim was needed. Serena initialized. An initial index_status call without its required project argument returned a usage error; the explicit project retry succeeded and supplied the recorded evidence.

Confidence is high for the exact observed command outcomes. No new failure was found beyond disclosed baseline lint/audit limitations. Recommended next action: Brain retain evidence, complete G3a10 implementation and fresh review, run final metadata checks, and publish through the shared destinations while keeping child execution blocked. There are no unresolved verification questions for this bounded baseline; final implementation/metadata/hosted outcomes remain deferred.
