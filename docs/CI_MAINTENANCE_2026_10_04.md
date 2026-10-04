# CI and contributor-tooling maintenance — 2026-10-04

## Scope and failure evidence

Defensive maintenance of this repository's dependency audit and local contributor tooling;
verification uses repository fixtures and local test processes only. No guest, installed service,
credential state, or runtime admission changes are authorized by this slice.

At main `a6ef83a6b43cd4940547b138529ebbb91cff8da6`, nightly CI run
[37205018420](https://github.com/Shrimpworks/capsule-corp/actions/runs/37205018420) failed at
TypeScript's dependency audit. Go, the full nightly Go suite, and the immutable archive job passed.
Run [36580264472](https://github.com/Shrimpworks/capsule-corp/actions/runs/36580264472), the first
failure in the recent streak, has the same failing step. This evidence identifies a deterministic
audit failure; it does not establish an intermittent test failure.

The lockfile retained Ajv's `fast-uri@3.1.6`. Upgrade that transitive package to `3.1.8` within
Ajv's existing range. This closes [authority injection](https://github.com/advisories/GHSA-qw65-cvwx-89v3),
[host confusion](https://github.com/advisories/GHSA-58mr-gqgx-xq4g), and the moderate host-case
normalization advisory reported by pnpm. Keep the existing high-severity audit gate unchanged.
No direct dependency, override, major upgrade, or advisory exception is added.

## Dependency-policy disposition

- Capability/roadmap: existing schema/fixture build validation; reuse-map row “TypeScript
  schemas/SDK/CLI” remains BUILD-NARROWLY. This is maintenance of existing build tooling,
  not product admission.
- Exact candidate: fast-uri 3.1.8; npm integrity bytes and complete graph retained in
  `pnpm-lock.yaml`. Fastify maintains the MIT-licensed package and its advisory response.
- Trust/authority: build/test dependency through Ajv; no new product, key, guest, filesystem,
  process, or network capability. Supported Node/pnpm pins remain unchanged.
- Provenance: package-manager integrity verification, frozen installation and upstream advisory
  readback on 2026-10-04; no independent release signature or attestation claim.
- Compatibility/fault evidence: existing schema negative/restoration/mutation corpus and script
  tests remain the oracle. No new parser primitive or new boundary claim.
- Reproduction/rollback: `pnpm install --frozen-lockfile`; revert this lockfile change to recover
  prior dependency bytes, which will again fail the audit. No offline cache availability claim.
- Upgrade/removal owner: repository maintainers through Dependabot and the existing nightly
  `pnpm audit:dependencies` gate; high findings remain blocking and require a reviewed update.

## Contributor tooling

AI Central selection is pinned to `e08d63e874f39264ed36f31812476f0816b9e49e` and curated in
`.codex/ai-central-skills.json`. The refresh command delegates to the maintained bundle installer
with `--sync`, preserving real files and unrelated links. Biome excludes both canonical
`.agents/skills` and compatibility `.codex/skills` links to keep local sample fixtures outside
repository linting. Repository-owned context remains subject to the normal gates.

Serena's shared project configuration selects Go, TypeScript/JavaScript, C, Swift, shell, YAML,
and JSON. Codebase Memory is indexed in full mode so scripts and test sources are included.
Both remain optional local contributor tools; caches, logs and generated memories are ignored.
See `.codex/AI_CENTRAL.md` for refresh, reindex, coverage limitations and session reload details.

## Verification and limitations

Scoped work status: `PASSED` for local verification. CI/default-branch restoration remains
`IN_PROGRESS — TRENDING_GOOD` until the maintenance PR passes and is merged.

- Frozen pnpm installation, audit (zero advisories), check, lint, full `pnpm test`, schema/ADR
  verification and site build passed. Bootstrap regression verifies curated dry-run delegation,
  no writes, source-pin refusal and CLI invocation through a symlinked temporary path.
- Go 1.25.13 full tests, vet and build passed. Security/correctness lint with an isolated cache
  passed with zero issues. Required unfiltered lint was run and reports 50 existing revive
  documentation findings, the backlog tracked in issue #217; no Go source changes in this slice.
- `go run ...govulncheck@v1.6.0` under the host's Go 1.26.5 reports host standard-library findings.
  Repeating with Go 1.26.6 passed. Building the same pinned scanner under 1.26.6 and running it
  with `GOTOOLCHAIN=go1.25.13` also reports no vulnerabilities. The project pin remains unchanged.
- AI Central's full check passed. Serena health/index and memory-reference check passed.
  Codebase Memory full indexing retained 19,638 nodes and 62,388 edges; coverage readback confirms
  scripts, TypeScript tests and the current teardown model. Partial Markdown and intentionally
  invalid UTF-8 fixture parsing remain disclosed source-fallback cases.
- Removed 24 missing worktree registrations and one clean detached review checkout. Preserved
  the old review checkout with modified instructions and an untracked MCP file. Local instruction
  migration was retained in this change with a named stash as an additional recovery copy.
- Local curated installation contains 70 skills, no broken links; 151 deselected managed links
  were pruned across canonical and compatibility directories. No unrelated real skill was removed.

Final self-review covered pin enforcement, argument-array subprocess invocation, dry-run behavior,
managed-link ownership, ignored cache boundaries and the narrow dependency diff. No product
execution or authority implementation changed.
