# Claude project instructions

Follow the repository-wide instructions in `AGENTS.md`. In particular, read the
project definition, architecture, threat model, and applicable ADRs before
changing execution or policy behavior.

Keep machine-specific hooks in ignored `.claude/settings.local.json`.

<!-- codebase-memory-mcp:start -->
## Codebase Memory MCP

For structural codebase exploration, use the installed `codebase-memory` skill.
When Codebase Memory MCP tools are available, call `list_projects` or
`index_status` first; index this checkout with `index_repository` if needed.
Use `search_graph` for symbols, `trace_path` for callers and dependencies,
and `get_code_snippet` for source. Check `check_index_coverage` for cited paths
and scopes, and paginate relevant results.
When Codebase Memory MCP tools are unavailable, or coverage is incomplete,
use `rg`, `rg --files`, and focused source reads; state the limitation and continue.
Record durable decisions in repository documentation or `manage_adr`.
<!-- codebase-memory-mcp:end -->
