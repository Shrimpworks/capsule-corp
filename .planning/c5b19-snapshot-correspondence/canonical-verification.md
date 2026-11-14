# G3a9 canonical baseline verification

Status: **PASSED** for required canonical baseline checks. Existing unfiltered lint debt #217 remains **BLOCKED**. Parent full Gate3 remains **IN_PROGRESS — TRENDING_GOOD**; child execution and product admission remain **BLOCKED**.

## Scope and ownership

Internal verification-only G3a9-VERIFY assignment in `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-snapshot-correspondence`, HEAD/base `f3f2e0755c0c284dd28019225bef0cf8482d3cdd` (merged PR406). The question tested was whether the canonical baseline passes its required local checks before the bounded snapshot-correspondence delivery. This worker made no repository edits, branch, commit, push, PR or subagent. Brain owns retention, final metadata checks, independent review, integration and shared PR publication. Root's concurrent documentation/planning changes and unrelated `.serena/project.yml` were preserved.

Defensive authorized scope was ordinary automated checks of this owned local repository. No C5b19 native fixture, child, guest, lifecycle/signal/wait target, live store, installed service, credential or user data was exercised. Existing routine advisory queries were used by the scanner/audit. Archive Go1.25.13 and run/admission gates were not changed. This baseline does not establish correctness or review of the future G3a9 archive increment, hosted acceptance or product admission.

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

Observed scanner output: `No vulnerabilities found.` Dependency audit separately reports four findings: two low and two moderate, below the unchanged high threshold. This is not a claim that dependencies have no vulnerabilities. Scoped Go lint reports zero issues; unfiltered lint shows50 revive findings under its output cap, not a total debt count. No check, threshold, exclusion or branch protection was weakened.

## Evidence and limits

Checks started `2026-11-14T02:30:10.824090+00:00` and completed `2026-11-14T02:31:15.149642+00:00`. Exact versions, commands, timestamps, exits and hashes are in `/tmp/c5b19-g3a9-baseline.json`.

- `/tmp/c5b19-g3a9-install.log`: SHA256 `665010577978fedf41608da33b6ee3b6892ab208ff54a6df4f7b394a6de6806d`.
- `/tmp/c5b19-g3a9-ci.log`: SHA256 `0b991975c98449f0741bae51141f79327bdbe977733326d96e2b88370671f435`.
- `/tmp/c5b19-g3a9-pnpm-test.log`: SHA256 `65c1faae21a458d513d71d5b202f43e9f5b1f955c3e8b6edc1a038ce224c4712`.
- `/tmp/c5b19-g3a9-lint-full.log`: SHA256 `b12416c706c0ba0f921ded07b90baf9c6d078c2d02ae8dcd36461e468f078425`.

All four log hashes were verified before handoff; all owned processes completed. Raw `/tmp` logs are ephemeral. Brain should retain this handoff and validation JSON in canonical planning. Future JSON formatting, later documentation edits and final hosted checks are not certified here; final integration metadata validation remains Brain-owned.

Codebase-memory Tier2 project listing was complete (78 projects), confirming the primary checkout index; coverage was checked for AGENTS.md, Makefile, go.mod, biome.json and both consulted planning reports. Generation2026-10-04T22:18:10Z is old/not watched and not this worktree. AGENTS/Makefile/go.mod reported best-effort metadata_match; biome was not_tracked, planning freshness missing. Exact target-worktree reads supplied current configuration and planning evidence. No structural exploration or whole-graph completeness claim was needed. Serena initialized. One ancillary literal log search was temporarily blocked by a tool hook's accumulated-search counter; direct log inspection then established counts and limitations without affecting checks.

Confidence is high for the exact observed command outcomes. No new failure was found beyond disclosed baseline lint/audit limitations. Recommended next action: Brain retain evidence, complete G3a9 implementation and fresh review, run final metadata checks, and publish through the shared destinations while keeping child execution blocked.
