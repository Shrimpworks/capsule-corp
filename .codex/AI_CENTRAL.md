# AI Central Integration

This directory contains project-local AI context integrated with the sibling `ai-central`
repository.

## Installed selection

- Install mode: local symlinks ignored by Git
- Shared steering: `javascript-esm-steering.md`
- Skills: curated bundles in `ai-central-skills.json`: core, node, orchestration, documentation,
  delivery, engineering, rust, security-testing, and planning

This uses the shared-checkout pattern. Repository-specific instructions and steering remain
real, tracked files. Reusable skills are installed canonically under `.agents/skills/` with
`.codex/skills/` compatibility links. Shared skills and steering do not become part of the GitHub
repository.

The existing root `AGENTS.md` remains authoritative for repository-specific security guidance.
Files under `steering/` supplement it. If generic AI Central content conflicts with `AGENTS.md`, the
root instructions win.

## Refresh workflow

Create or refresh the links:

```sh
pnpm codex:links
```

The command expects AI Central at `../ai-central/templates` by default. When it lives elsewhere,
set `AI_CENTRAL_HOME` to either the AI Central repository or its `templates` directory.

Preview without writing:

```sh
pnpm codex:links -- --dry-run
```

The setup command preserves real repo-owned files and
delegates skill selection and pruning to AI Central's maintained `install-skill-bundle.sh --sync`.
The selection retains engineering and security-testing workflows and omits product strategy,
framework-specific frontend, infrastructure, and brevity bundles. Source attribution and licenses
remain owned by AI Central because these links are local and ignored by Git.

## Provenance pin

`AGENTS.md` requires that generic steering or skill guidance never weaken this repository's
security posture. Nothing verifies that on its own until a commit is pinned: the linker will link
whatever `SKILL.md`-containing directories it finds under `AI_CENTRAL_HOME`, with no check that the
content matches a reviewed revision.

After reviewing a linked `ai-central` checkout, record its exact commit:

```sh
pnpm codex:links -- --record-pin
```

This writes `.codex/ai-central-pin.json` (a real, tracked file — it travels with the repository and
its changes are reviewable in `git diff`, unlike the symlinks themselves). Every later `pnpm
codex:links` run compares the checkout's current commit against that pin:

- No pin recorded yet: linking proceeds with a warning. This is expected before the first
  `--record-pin`, and is not a substitute for actually running it.
- Checkout commit matches the pin: linking proceeds silently.
- Checkout commit differs from the pin: linking refuses (nonzero exit) until the change is reviewed
  and either accepted (`--record-pin` again) or reverted.
- `AI_CENTRAL_HOME` is not a git checkout: linking proceeds with a warning, since there is no commit
  to compare.

This detects a drifted or substituted `ai-central` source directory; it does not vet AI Central's
own content, and it does not apply retroactively to content already linked before a pin existed.

## Language and code graph tooling

Serena reads `.serena/project.yml`: Go, TypeScript/JavaScript, C, Swift, shell, YAML, and JSON.
Caches, logs, and onboarding memories remain local. Run `serena project index .`,
`serena project health-check .`, and `serena memories check .` after setup. Restart an already
running Serena MCP session after changing its language selection; project activation alone can
reuse its cached configuration. YAML/JSON servers provide symbols but do not support every
reference-search operation.

Codebase Memory is optional contributor tooling. Index this checkout in `full` mode to include
scripts and test sources, check `index_status`, and use `check_index_coverage` before relying on
cited paths. Invalid UTF-8 adversarial fixtures intentionally cannot parse; partial Markdown
coverage and ignored/generated/binary files require source fallback. Neither MCP is a CI or
product dependency.
