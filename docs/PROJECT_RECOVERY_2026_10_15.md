# Project recovery — 2026-10-15

Work item: RECOVERY-20261015. Status: `PASSED` for initial recovery deliverables and reviewed
maintenance repairs published in [PR #384](https://github.com/Shrimpworks/capsule-corp/pull/384).
Scope: contributor tooling, reproducible baseline, implementation/plan reconciliation and two
bounded maintenance repairs. Base: `fc8af2df91d443f52a7677c2bf9f7dc1d96f9ab1`.
Owner: Brain orchestrator; internal agents supply readiness, architecture, validation and review.
Remaining work outside this delivery: remote CI/merge and the separately scoped roadmap lanes.
Next action: complete PR acceptance, then assign the next bounded roadmap item.
Parent: owner-only internal alpha remains `IN_PROGRESS — TRENDING_GOOD`; installed security
boundary and product/runtime/profile admission remain `BLOCKED`.

This is an initial recovery checkpoint, not a security audit or a declaration that the product
is ready. Evidence labels mean **verified** (direct source/tool/test observation), **inferred**
(engineering judgment from that evidence), or **unresolved** (not established). Historical
experiment summaries do not become newly verified run evidence by being read here.

## Tooling readiness report

The bounded readiness investigation is `PASSED`; optional-tool limitations remain below.
No critical retrieval failure prevents these maintenance repairs. Their own tests and independent
review remain required. No tooling upgrade or provenance-pin reset is needed.

| Tool | Verified observation | Limitation and operating rule |
| --- | --- | --- |
| Serena | Activation and Go `NewServer`, TypeScript `CapsuleClient`, C `set_limit` symbol/reference probes succeeded; shell variable references returned uses | Swift `parseEnvelope` references returned empty despite two source calls; JS cross-file test calls and shell `trap cleanup` were omitted. Use exact source search for these relationships. |
| YAML/JSON LSP | Symbol overviews succeeded | Reference queries returned `-32601`; unsupported, not an absence proof. |
| Serena health/memory | Six memories read; `serena memories check .` passed | Historical health command caught a reference failure then reported success. Per-language probes supersede that broad claim. No blanket memory-freshness guarantee. |
| Codebase Memory | Ready/full index at base HEAD, generation `2026-10-04T22:18:10Z`, 19,631 nodes/62,375 edges; symbol/snippet/trace/coverage calls succeeded | One heuristic `NewServer` self-edge contradicted source. Graph is a navigation aid; check material edges against source. |
| Graph coverage | Exact cited tooling and core source paths returned `metadata_match`; relevant result pages exhausted | Generated/cache/link paths excluded; two intentional invalid-UTF-8 fixtures and partial Markdown need source fallback. Clean coverage is best-effort. Recovery worktree edits require fresh coverage or direct diff reads. |
| AI Central | Both active/default source checkouts clean at reviewed `e08d63e874f39264ed36f31812476f0816b9e49e`; all 143 local links resolve and are ignored; 140 skill links contain `SKILL.md` | Refresh selected absent `javascript-esm-steering.md` and silently reported zero steering links. R4A fixes maintained selection and missing-input failure. |
| Supporting tools | AST-grep MCP found `http.NewServeMux()`; ast-grep 0.45.3, shellcheck 0.11.0, actionlint 1.7.12, golangci-lint 2.12.2 available | Availability is not a completed review. `lychee` absent; no need to install for this scope. Fresh independent review is separate. |
| GitHub | Remote main equals base; open PR/issue and CI reads succeeded through normal CLI authentication | No credentials extracted or changed; no unrelated repository mutation. |

Actual AI Central links use the local `.ai-central` checkout; default refresh uses sibling
`repos/ai-central`. Both identities matched the pin at inspection. Matching today does not prove
future link content: retain the source pin and recheck after refresh. Existing source-owned files,
unrelated links, dry-run nonmutation and pin-mismatch refusal must survive R4A.

## Project recovery assessment

**Purpose — verified.** Capsule aims to let AI agents propose bounded JavaScript tasks while a
human Broker approves exact Supervisor-retained bytes and only the Execution Supervisor owns
hostile guest lifecycle. The public abstraction is a task. [PROJECT](PROJECT.md),
[ARCHITECTURE](ARCHITECTURE.md), [TECHNICAL_DESIGN](TECHNICAL_DESIGN.md),
[threat model](security/THREAT_MODEL.md) and [ADR-0040](adr/0040-freeze-owner-only-internal-alpha-posture.md)
remain the authoritative boundaries. One owner Mac, one exact `main.mjs`, bounded inline JSON,
concurrency one and one fresh guest per attempt define the first alpha.

**Working implementation — verified within bounded source/test inspection.** The repository has
Go/TypeScript scaffolding, strict passive protocol/canonicalization fixtures, immutable source and
plan registration, public-key approval verification, one-use approval/attempt transactions,
FakeBackend lifecycle/completion-last mechanics, fixed-store archive/backup conformance, and a
passive teardown model. These are useful components, not a connected installed product path.
`internal/api/server.go` remains a diagnostic health surface; it is not an authenticated Supervisor
endpoint. The [vertical-flow readiness map](ALPHA_VERTICAL_FLOW_READINESS.md) correctly separates
passive libraries from missing product consumers.

**Incomplete — verified from source and canonical contracts.** Authenticated product IPC,
installed identity/key separation, protected state, live Broker UI/signing, real adapter/recovery,
controlled transport, completion/store/receipt composition and exact profile admission remain
open. [Issue #309](https://github.com/Shrimpworks/capsule-corp/issues/309) records unavailable v19/v27
raw evidence. Historical observations cannot support reproducible release/admission claims until
recovered or replaced by a separately authorized, newly identified experiment.

**Architecture — verified design, incomplete enforcement.** The daemon proposes, the Broker
renders/approves, and the Supervisor independently authorizes and owns execution. No recovery
change introduces a privileged helper or broadens authority. [ADR-0047](adr/0047-prepare-teardown-obligation-before-launch.md)
selects precommitted cleanup obligations and storage-independent stop control after the bounded
C5b16 timing finding. The C5b18 passive model and C5b19 reviewed packet are design/conformance
progress, not proof of runnable or installed teardown.

**Planning drift — verified.** The current plan still asked to integrate Gate 2 although #378 was
already on main, and its bottom checkpoint described providers as missing despite later benign
provider slices. This recovery reconciles those entries without promoting their scope. Setup
text cited Go 1.23 and claimed full `make ci` parity; both conflict with actual pins/workflows.
The C5b19 Gate-4 heading's older “on gate 3” wording is read with its actual requirements:
implemented source/binary review and separate first-run authority are still outstanding.

**What prevents progress — verified/inferred.** A newly failing vulnerability gate, silently
incomplete tooling refresh, stale next-action text, absent historical raw evidence and unclosed
installed composition gates all matter. The inference is to restore reliable feedback first,
then resume the existing narrow security experiment sequence. A rewrite or broad cleanup would
not close these blockers. Four open dependency PRs (#366, #367, #381, #382) and 16 open issues were
observed; their presence is not evidence they are safe to merge or still correctly scoped.

**Reconsideration and unresolved questions.** Confirm whether C5b19 implementation exists in the
owned archive before dispatching a new cross-repository assignment. This local recovery did not
inspect or modify that repository. Revalidate its frozen host/toolchain contract before any future
build/run; the repository Go update does not silently revise that historical packet. Determine
whether v19/v27 evidence is recoverable before proposing a new run. Keep Source Validator and
TypeScript preparation off the internal-alpha critical path as ADR-0040 requires. F6 storage and
external-alpha continuity remain later decisions, not prerequisites invented by this recovery.

## Baseline validation report

Baseline uses Node 22.22.1, pnpm 10.28.2, Go 1.25.13 on macOS 26.6.2/25G83 arm64;
golangci-lint 2.12.2 was built with Go 1.26.2. Exact toolchain selection is explicit because host
Go was 1.26.5. Initial sandbox install DNS/cache-write refusal was environmental; the authorized
normal network/cache retry succeeded with frozen dependencies.

The latest observed [main CI run 38736586664](https://github.com/Shrimpworks/capsule-corp/actions/runs/38736586664)
failed at Go vulnerability scanning. TypeScript, compiled-archive verification and the full
nightly Go suite passed. Recent preceding CI runs also failed; this is preexisting at the base.
No new product exploit was attempted or demonstrated.

Local detailed exit codes and final candidate results are recorded in the validation appendix
below. A passing local test does not substitute for remote CI or an installed security corpus.
Unfiltered Go lint's known exported-contract backlog (#217) must be reported separately from the
existing CI security/correctness and changed-code documentation gates; it is not suppressed here.

## Prioritized recovery roadmap and initial execution plan

Effort is relative: S = bounded maintenance, M = multiple reviewed slices, L = multi-gate program;
these are engineering estimates, not dates. Each child reports scope separately from alpha.

| ID / priority / owner | Outcome and dependencies | Scope / exclusions | Acceptance and verification | Risk / effort / status |
| --- | --- | --- | --- | --- |
| R4A / P0 / internal tooling worker | Reliable pinned AI Central refresh; R0 | Linker, existing tests and integration documentation; no upstream or pin changes | Red/green missing-steering regression; all selected files validated before mutation; dry-run, owned-file preservation, curated sync and pin-refusal tests | Local link regression; S; PASSED |
| R4B / P0 / internal baseline worker | Restore Go vulnerability gate; R2 and upstream advisories | Repository Go declaration and current onboarding; no module dependency or historical experiment rewrite | Same scanner old/new comparison; full Go tests/vet/build, lint accounting, unchanged schema/crypto/restoration fixtures; PR CI | Go release semantics/platform differences; S; PASSED locally; remote CI is a separate acceptance gate |
| R5 / P0 / Brain + fresh reviewer | Deliver recovery evidence and repairs; R4A/R4B | One shared branch/PR, canonical index and onboarding reconciliation | Independent review max three attempts; required checks; commit/push/body/state readback | Avoid overstated health; S; PASSED for review and PR publication |
| NEXT-1 / P1 / next bounded contributor task | Reproduce CI locally; R5 | Issue #265: JSDoc entrypoint first, then explicit CI/Make parity; no permissions/advisory exceptions | Positive/negative exported-doc fixture, unchanged file selection, inventory all local/CI differences before deduplicating builds | Full-audit backlog differs from PR scope; S–M; BLOCKED on next bounded task assignment |
| NEXT-2 / P1 / Brain + archive implementation worker | Resume C5b19 implementation-only lane; R5 and archive-state readback | Exact frozen packet in named archive; no child run, guest, installed service or product activation | Inventory source state; TDD no-effect storage/oracle checks; exact source/build manifest; actual-byte independent review before Gate 4 | Frozen host/toolchain mismatch reopens packet review; M; BLOCKED on archive reconciliation and scoped dispatch |
| NEXT-3 / P1 / evidence owner | Resolve #309 independently of implementation | Exact authorized backup/snapshot evidence recovery; no recreated historical identities | Recovered manifest/raw evidence verified and pinned immutably, or explicit owner decision for a new versioned experiment | Backup availability unknown; M; BLOCKED on owner-supplied evidence location or separate run decision |
| NEXT-4 / P1 / Brain + platform owners | Close installed authority composition in readiness-map order | C3 identity/root, C2 IPC, C6 signing, real adapter/transport, then C7 admission; no wholesale redesign | Exact named environments, independent review, rejection/restoration/fault corpora and retained immutable evidence per gate | Highest authority risk; L; BLOCKED on named platform/evidence/authorization gates |
| NEXT-5 / P2 / maintenance owner | Reduce #217/#216 and storage/protocol debt when touched | Small changed-component batches; no broad storage rewrite | No new lint debt; contract tests and independent review for behavioral changes | Cleanup can obscure boundary changes; M; BLOCKED on prioritization after critical lanes |

The Brain owns integration and the next assignment. All current workers are internal subagents on
`codex/project-recovery-20261015`; none independently commits or publishes. Review attempts are
recorded below. Three failed implementation/review attempts stop the cycle and require a focused
research/alignment decision; this recovery does not reset historical C5b review counts.

Progress measures: required gates reproduced; failing advisory gate removed on exact candidate;
missing-steering regression proven and fixed; stale execution instructions reconciled; reviewed
PR published. Sustained progress is measured at subsequent acceptance checkpoints, not by agent
activity. At each merge, update CURRENT_WORK_PLAN with owner, next runnable scope and evidence;
recheck the existing nightly CI signal, archive claims and remaining dependencies. No new
recurring automation or background execution is established by this document.

## R4B evidence-backed dependency/toolchain decision

Capability: maintain existing Go scaffold and its CI security feedback. Reuse-map rows GO-3
(ADOPT-PLATFORM/STDLIB) and TEST-1 (ADOPT-PINNED/TEST-ONLY) apply. Candidate: Go 1.26.9,
released 2026-10-08; existing govulncheck remains v1.6.0. The [Go release history](https://go.dev/doc/devel/release)
and [GO-2026-6617 advisory](https://pkg.go.dev/vuln/GO-2026-6617), retrieved 2026-10-15, support the
minimum fixed maintained line. Go supports a release until two newer major releases exist; Go
1.25 is outside that window. This is a security-motivated update, not a latest-version upgrade.

The unchanged-base scanner reports GO-2026-6617, 6613, 6612, 6611, 6608, 6607 and 6603 in the
standard library, all fixed by 1.26.9. Candidate validation uses the same scanner against both
versions. Advisory reachability is not a demonstrated exploit. [Go 1.26 release notes](https://go.dev/doc/go1.26)
identify runtime/GC/cgo/compiler changes, so existing boundary and restoration fixtures must
pass; a clean scan alone is insufficient. [Toolchain selection](https://go.dev/doc/toolchain)
explains why `GOTOOLCHAIN=go1.26.9` is necessary for exact local reproduction: the module directive
sets a minimum, not a forced downgrade of a newer host compiler.

Policy disposition: Go project maintenance/security ownership and BSD-style distribution license;
existing product/build/test TCB, no new authority, module or primitive. Existing module graph and
`go.sum` remain unchanged. Go's verified toolchain download is used; no independent signature,
reproducible compiler build or offline cache availability claim. Existing input/resource limits
and refusal/crash/corruption/cancellation semantics remain unchanged and covered by their current
fixtures; no new installed or guest evidence is claimed. Maintainers own subsequent advisory
response, compatibility review and upgrade/removal. Reverting the pin restores old build semantics
but reopens the known failing scan, so it is not a green rollback. Binary-size/platform-wide
comparison and release-grade offline database/toolchain snapshots are unresolved and cannot be
used for product admission. Historical experiment toolchains require their own review.

## Validation appendix and review ledger

Baseline characterization: `PASSED`; baseline all-gates result: `BLOCKED` by old-pin advisories
and preexisting documentation lint debt. R4A/R4B local repair verification: `PASSED` for their
bounded scopes; independent review and PR publication are `PASSED`. Remote CI/merge is separate.

| Command/evidence | Initial base | Repaired tree |
| --- | --- | --- |
| `pnpm install --frozen-lockfile` | Exit 0 after environmental retry | Dependencies/lockfile unchanged |
| `pnpm check`, `pnpm lint`, `pnpm test` | Each exit 0 | Each exit 0 on final integrated tree |
| `pnpm verify:schemas`, `pnpm verify:adrs`, `pnpm site:build` | Each exit 0 | Each exit 0 on final integrated tree |
| `pnpm audit:dependencies` | Exit 0; no advisories | Dependency bytes unchanged |
| `go test ./...` without `-short` | Exit 0; 168.55 s | Exit 0; 167.57 s |
| `go vet ./...`, `go build ./...` | Each exit 0 | Each exit 0 |
| `golangci-lint run ./...` | Exit 1; revive backlog | Exit 1; revive backlog |
| Existing CI `golangci-lint run --disable revive ./...` | Not a substitute for full lint | Exit 0 |
| `go run golang.org/x/vuln/cmd/govulncheck@v1.6.0 ./...` | Auto host scan exit 1 against host Go 1.26.5 | Exact Go 1.26.9 exit 0 |
| Same separately built v1.6.0 scanner with exact `GOTOOLCHAIN` | Go 1.25.13 exit 3; seven stdlib findings | Go 1.26.9 exit 0; none |
| `go mod verify` | Not separately run | Exit 0 |
| Linker regression suite | New tests: 7 pass / 6 fail before fix | 13 pass / 0 fail |
| Real pinned-source linker dry run | Exit 0 but zero selected steering links | Exit 0 and three selected links; no writes |

The 50 reported revive findings are the default per-linter output cap, **not the total backlog
count** or a claim of identical finding lists. No Go source changed and no lint rule was disabled
by this patch. Closure remains owned by maintainers under #217. The original Go scan invocation
used `GOTOOLCHAIN=auto`, so its host findings are not mislabeled as exact-project results; the
separate same-scanner comparison establishes the exact old/new pin result.

Observed Go toolchain identity: `golang.org/toolchain@v0.0.1-go1.26.9.darwin-arm64`, cached module
sum `h1:u1IV5D8HPeURJjWAqzcC8QRXDJqvXlh7FpsA4lGIg1c=`, Go command SHA-256
`318dee99db233501f6c419fbb4616e5a345ed0b831dbd6c7669b24e1ff8e9e2e`.
Unchanged `go.sum` SHA-256:
`fcefebf3982c202f169f7f1e12cf661287a26ff163c8bfad1689fc665583f472`.
Direct modules remain `fxamacker/cbor/v2 v2.9.2` and `x/sys v0.47.0`; indirect `float16 v0.8.4`.
These are local identity observations, not independent compiler provenance.

Local raw logs are under `/tmp/capsule-recovery-logs`; their filenames, commands, exit codes and
hashes are retained in the [validation ledger](../.planning/project-recovery-20261015/validation.json).
Temporary raw logs are not a durable release archive. The retained summaries and ordinary CI are
sufficient for this maintenance claim; no guest or admission evidence depends on them.

Review instance 1 of 3: fresh-context read-only reviewer dispatched against the explicit working
copy at base `fc8af2d`; verdict **Ready**, no actionable findings. Reviewer independently ran
the 13 linker tests and real-source dry run and verified all 32 retained command-log hashes.
Full Go suite evidence was inspected, not independently reexecuted. The preliminary read preceded
author rationale, but no separate preliminary ledger was saved; this process limitation is
retained in the [review report](../.planning/project-recovery-20261015/review.md). No review budget
from an older C5b workstream is reused.

Publication: branch `codex/project-recovery-20261015`, reviewed implementation commit
`22aad310d60009642097c21cf4b822b13bd14069`, intentionally ready PR #384 (`isDraft=false`), base
`main`. Stored body read back as real multiline Markdown; base/head/commit/state verified.
Initial [CI run 38762955023](https://github.com/Shrimpworks/capsule-corp/actions/runs/38762955023)
passed compiled-archive and dependency review; remaining jobs were pending at this retained
snapshot. Latest GitHub checks, rather than this historical snapshot, govern merge. No merge,
guest/installed operation or product-admission decision was performed by recovery. This closing
status record is documentation-only follow-up to the reviewed implementation.
