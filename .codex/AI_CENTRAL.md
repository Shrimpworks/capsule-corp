# AI Central Integration

This directory contains project-local AI context integrated with the sibling `ai-central`
repository.

## Installed selection

- Install mode: local symlinks ignored by Git
- Shared steering: `javascript-esm-steering.md`
- Skills: all reviewed skills currently exposed by AI Central

This matches the Wap Labs and Reef pattern. Repository-specific instructions and steering remain
real, tracked files. Reusable skills and shared steering are recreated locally and do not become
part of the GitHub repository.

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

The setup command preserves real repo-owned files, repairs stale or broken managed symlinks, and
links the full reviewed AI Central catalog. That includes Caveman, Cavecrew, Caveman Stats, and
Hallmark. Source attribution and licenses remain owned by AI Central because these links are local
and ignored by Git.
