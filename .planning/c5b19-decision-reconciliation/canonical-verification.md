# G3a8 post-fix canonical verification

Status: **PASSED** for required post-fix canonical checks. Known unfiltered lint debt #217 remains **BLOCKED**. Parent full Gate 3 remains **IN_PROGRESS — TRENDING_GOOD**; child execution and admission remain **BLOCKED**.

## Scope and exact candidate

Internal verification-only assignment in `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-decision-reconciliation`, HEAD/base `8e4fd22cc00cd5e961544f6a803974600dc0e751`. The CI-RACE worker confirmed the two source files frozen before verification. The following SHA-256 values matched its handoff and were identical at check start and end:

- `internal/execution/authorityplane/passive_product_adapter.go`: `785ac54438cea50838b39d8de8a4cf1b8f9b0dd6b887bfedd51beea5543cc9f4`.
- `internal/execution/authorityplane/passive_product_adapter_test.go`: `4621e65f44dd4fcc1ded6f553d8ba29e2c38245fdd6550a51ef89a09274f9484`.

The defensive authorized target was standard automated checks of this owned local repository. No C5b19 native fixture, child, guest, live store, installed service, or user data was exercised. This worker made no repository edits, branch, commit, push, PR, or subagent. Other workers' documentation/planning edits and `.serena/project.yml` were preserved. The Brain owns final metadata validation, independent review, evidence retention, and publication.

## Verification and environment

Exact Go 1.26.11 darwin/arm64 was selected by prepending `/Users/dsteele/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.11.darwin-arm64/bin` to PATH with `GOTOOLCHAIN=local`. Make's scanner explicitly uses `GOTOOLCHAIN=auto`, invoking the same exact Go binary and unchanged module requirement. Observed versions remain Node 22.23.3, pnpm 10.28.2, golangci-lint 2.12.2 built with Go 1.26.2. Node differs from documented patch 22.22.1; no upgrade was performed.

| Command | Result |
| --- | --- |
| `make ci` | exit 0 |
| `pnpm test` | exit 0 |
| `git diff --check` | exit 0 |

Both CI and standalone tests report 177 workspace tests (164 protocol, 10 SDK, 3 MCP) plus 132 script tests, all passing. Full Go package tests passed; some packages were cached, while `authorityplane` ran and passed in 0.348s. Canonical CI includes formatting/vet, correctness and changed-line revive lint, pinned govulncheck, dependency audit, site/schema/ADR checks, TS check/lint, changed JSDoc, Go tests/build, and built JS/script tests. Govulncheck reports no vulnerabilities. Dependency audit reports two low and two moderate advisories below the unchanged high threshold. No lockfile, go.mod, or go.sum diff was present.

Frozen install and unfiltered lint were not repeated, as assigned. Their retained baseline evidence is `/tmp/c5b19-g3a8-baseline.md` and `/tmp/c5b19-g3a8-baseline.json`: frozen install passed; full unfiltered lint exited 1 with known #217 revive debt (50 displayed under an output cap, not the total backlog). These are prior results, not post-fix reruns. No threshold, exclusion, branch protection, or test was weakened.

## Evidence and limits

Checks started `2026-11-14T02:05:17.132094+00:00` and completed `2026-11-14T02:06:39.395815+00:00`. Exact versions, command timestamps/exits, test counts, source hashes, and log hashes are in `/tmp/c5b19-g3a8-postfix-verification.json`.

- `/tmp/c5b19-g3a8-postfix-ci.log`: SHA-256 `387e67562766619e5da287f3d97e2e9f546f3651e7e4ac962506c8f5cb4d0593`.
- `/tmp/c5b19-g3a8-postfix-pnpm-test.log`: SHA-256 `98c3a8f19b26736ddb09354b361c50495793d608d117afad008c476557871160`.

All owned command processes completed. These ephemeral `/tmp` artifacts require Brain-owned retention. Later retained JSON formatting or documentation changes are not certified; final metadata checks remain Brain-owned. This result does not establish archived G3a8 harness correctness, hosted CI, independent review, or product admission.

Codebase-memory coverage was checked for both source paths, Makefile, and go.mod: all four reported best-effort `metadata_match`, but the primary checkout index generation is old (`2026-10-04T22:18:10Z`) and not watched. Exact canonical worktree reads and source hashes supply the current evidence; no graph-completeness claim is made.

Confidence is high for these exact command outcomes. No new failure was observed. Recommended next action: retain post-fix evidence, finish independent review and final metadata validation, and publish the bounded slice while keeping child execution blocked.
