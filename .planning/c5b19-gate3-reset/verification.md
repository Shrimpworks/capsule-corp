# Gate-3 reset verification

Status: `PASSED` for canonical CI-equivalent checks and standalone tests. Existing
unfiltered lint debt #217 remains `BLOCKED`; no check or threshold changed.
Scope: documentation-only reset on `codex/c5b19-gate3-reset`, base
`a9e1d32afbb1fa1ef9c94175321c9c3bbb3d3361`. No C5b19 archive implementation or
execution. Child execution, Gates 4–6 and product admission remain `BLOCKED`.

| Command | Observed result |
| --- | --- |
| `pnpm install --frozen-lockfile` | exit 0 |
| `make ci CI_BASE=a9e1d32afbb1fa1ef9c94175321c9c3bbb3d3361` | exit 0 |
| `pnpm test` | exit 0 |
| `golangci-lint run ./...` | exit 1; existing revive debt #217, 50 displayed under output cap |
| `git diff --check` | exit 0 |
| Local Markdown target check across five changed canonical documents | All file targets exist; no external URL availability claim |

Go 1.26.11 darwin/arm64 selected via its exact local toolchain PATH and
`GOTOOLCHAIN=local`; make's scanner explicitly uses auto with that same binary.
Node 22.23.3, pnpm 10.28.2, golangci-lint 2.12.2 built with Go 1.26.2.
Node's patch differs from documented 22.22.1; no toolchain change in this reset.
The archive module remains frozen at Go 1.25.13 and was not executed.

Make covers formatting, vet, scoped correctness/changed-line lint, pinned
vulnerability scanner, dependency audit, site build, schemas/ADRs, TypeScript
checks/lint, changed JSDoc, Go tests/build and built workspace/script tests.
Observed 177 workspace tests (164 protocol, 10 SDK, 3 MCP) and 132 script tests
passing in both CI and standalone runs; Go package results include cached results.
Govulncheck reports no vulnerabilities. Dependency audit reports two low and two
moderate findings, below the unchanged high threshold. This is not a zero-advisory
claim. The full lint output cap is not a count of all outstanding debt.

[Machine-readable evidence](verification.json) records exact commands, exits,
UTC timestamps and SHA256 hashes of ephemeral `/tmp` logs. Checks ran
2026-11-14 05:02–05:03 UTC. No code/dependency changes were made; no new tests
were warranted for this reversible planning edit. Immutable source audits establish
the map; historical archive tests/mutations are not rerun or recertified here.
Archive tree was clean at `a9c457cb3b89a4ecd8c38684fd3d804d2089214c` after checks.
Unrelated `.serena/project.yml` remains excluded. Independent review and publication
acceptance are recorded separately in the delivery plan.

## Hosted metadata correction

Initial CI run 43559875578 failed TypeScript lint solely because verification.json,
retained after the local baseline checks, used expanded JSON arrays rather than the
repository formatter's output. Applied Biome formatting and asserted parsed JSON
identity with the reviewed commit. No evidence values, roadmap or behavior changed.
Final `pnpm lint` and `git diff --check` passed after this correction; subsequent
hosted results are recorded in the PR acceptance readback. This metadata-only fix
does not require another independent review instance.
