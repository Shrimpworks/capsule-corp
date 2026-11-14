# G3a4 canonical verification and Go security prerequisite

Status: PASSED for the scoped security patch and canonical required CI gates. Existing unfiltered lint debt #217 remains BLOCKED. Parent C5b19 Gate 3 remains IN_PROGRESS — TRENDING_GOOD; child execution and later gates remain BLOCKED.

## Scope and integration

Internal worker assignment; shared canonical worktree `/Users/dsteele/.codex/worktrees/project-recovery/capsule-corp`, branch `codex/c5b19-event-conformance`, base `2aefb15ee45018bb1045e711c77a0c7bc9f4c54d`. Root owns commit, PR and independent review. No branch, commit, push or PR was created by this worker.

After the baseline detected newly listed standard-library advisories, the parent explicitly assigned the minimal patch prerequisite. Exactly five files were edited: `go.mod`, `AGENTS.md`, `README.md`, `docs/DEVELOPMENT.md`, and `docs/DEV_MACHINE_SETUP_AND_RECOVERY.md`. Their active Go requirement moved from 1.26.9 to 1.26.11. CI already reads `go-version-file: go.mod`; there was no hardcoded workflow pin or `scripts/dev-toolchain.mjs` to change. Existing tests did not hardcode the old version. No new test mirroring the literal pin was added: the retained scanner failure and full-suite pass provide the behavior evidence.

Historical recovery/toolchain evidence and C5b19's frozen offline Go 1.25.13 packet were not rewritten. Existing `.serena/project.yml` changes were preserved and excluded from this worker's ownership. No product source, dependency, lockfile, scanner version, branch protection, runtime authority or lifecycle fixture changed.

## Baseline, correction, and exact commands

Normal Go toolchain verification downloaded `golang.org/toolchain@v0.0.1-go1.26.11.darwin-arm64` using `GOTOOLCHAIN=go1.26.11 go version`; it printed `go version go1.26.11 darwin/arm64`. All suite commands put the exact cached toolchain `bin` first on PATH and set `GOTOOLCHAIN=local`. Make's scanner command explicitly sets `GOTOOLCHAIN=auto`, using that exact Go binary and module directive. The host default was Go1.27.2, which was deliberately not used for the acceptance run. Node22.23.3 meets package engines >=22.22.0; pnpm10.28.2 and golangci-lint2.12.2 match configured pins. Node's observed patch differs from documented22.22.1 and was not changed.

| Command | Original exact Go1.26.9 | Final exact Go1.26.11 |
| --- | --- | --- |
| `pnpm install --frozen-lockfile` | exit0; already up to date | dependency bytes unchanged; not repeated |
| `make ci` | exit2, scanner stopped check | exit0 |
| `pnpm test` | exit0 | exit0 |
| `golangci-lint run ./...` | exit1,50 displayed revive findings | exit1,50 displayed revive findings |
| `git diff --check` | not separately needed before edit | exit0 |

`make ci` covers the AGENTS commands: gofmt, `go vet ./...`, full correctness plus changed-line revive lint, `go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...`, dependency audit, site build, schema/ADR checks, `pnpm check`, `pnpm lint`, changed JSDoc, `go test ./...`, `go build ./...`, and built TS/script tests. Standalone `pnpm test` was also run as requested. Final tests counted177 workspace tests (164+10+3) and132 script tests. Full Go tests passed, including registrationstate. No C5b19 experiment child fixture, guest, native lifecycle target, live store or service was run by this assignment.

Retained baseline scanner excerpt:

```text
GO-2026-6859: HTTP/1 server crash due to MaxBytesReader header race in net/http
Found in: net/http@go1.26.9
Fixed in: net/http@go1.26.11
GO-2026-6809: Distrust compromised intermediate certificates in crypto/x509
Found in: crypto/x509@go1.26.9
Fixed in: crypto/x509@go1.26.10
Your code is affected by 2 vulnerabilities from the Go standard library.
exit status 3
make: *** [check] Error 1
```

Final scanner: `No vulnerabilities found.` Final package audit retained its existing high-severity threshold and reported `4 vulnerabilities found`, `Severity: 2 low | 2 moderate`; its exit0 is not a vulnerability-free dependency claim. The unfiltered linter's50 findings are an output cap, not a full backlog count; all displayed findings were revive. Scoped CI lint remains green.

Parent verified [Go release history](https://go.dev/doc/devel/release) and [GO-2026-6809](https://pkg.go.dev/vuln/GO-2026-6809). The web tool could not retrieve GO-2026-6859 directly; its minimum fix version is evidenced by the hosted and local scanner output. The worker directly observed scanner minima and final behavior; it did not independently browse those pages. Go1.26.11 is the smallest listed patch satisfying both observed minima. This is a prerequisite repair, not a new toolchain family or a C5b19 packet change.

## Dependency-policy/provenance record

Applicable ecosystem rows: GO-3 (ADOPT-PLATFORM/STDLIB) and TEST-1 (ADOPT-PINNED/TEST-ONLY). This patch uses Go-maintained standard-library fixes, adds zero external modules, preserves role boundaries and the v1.6.0 scanner. Toolchain trust is build/test and existing Go product standard-library code, with no new authority. Source distribution license is BSD-3-Clause; Capsule maintainers own future upgrade/removal decisions. Fail-closed scanner rejection was preserved; no exclusion, advisory suppression or threshold relaxation was added. Reverting the five literals reproduces the old pin, while the retained baseline log records its current rejection. No new primitive or protocol restoration mutant applies to a patch-version literal update; existing full conformance and CI tests exercise the changed compiler/stdlib.

- Toolchain module zip hash: `h1:cGxzEKi4mhhQFd6Cz/2NSL4Y9nfIBCirzU1yn2rnOY8=`.
- Toolchain LICENSE SHA256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`.
- Toolchain bin/go SHA256: `20e07ffb4f08faaf125e413db024b6e988d3ec2799dcb3568a720d14be0ad133`.
- Current vulnerability DB was used on 2026-11-13 America/Toronto (UTC date2026-11-14); this is ordinary developer verification, not an immutable release DB snapshot.

## Evidence and limits

Combined machine-readable command results and SHA256 log evidence are `/tmp/c5b19-g3a4-baseline.json`; original results remain `/tmp/c5b19-g3a4-baseline-original.json`, final results `/tmp/c5b19-g3a4-final.json`. Raw logs are named there. Root should retain this report and JSON in the canonical PR; raw /tmp logs are ephemeral. This worker makes no hosted-CI, independent-review, full Gate3, product-admission, timing or live teardown claim.

Graph coverage was checked for consulted tracked configs and lint script against indexed primary capsule-corp generation2026-10-04T22:18:10Z: metadata matched configs best-effort; `.node-version` was reported not_tracked and newer `scripts/lint-go.mjs` missing in that primary snapshot. Exact target-worktree file reads, not graph completeness, establish the commands above. Structural code exploration was unnecessary for the five literal edits.

Recommended next action: include the five-file prerequisite in the same root delivery and fresh independent review; preserve historical evidence, run hosted checks, and retain all later execution gates.
