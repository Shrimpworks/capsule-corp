# Development

## Prerequisites

- Go 1.23 or newer, as declared by `go.mod`
- Node.js 22.22.1, as declared by `.node-version`
- pnpm 10.28.2, as declared by `package.json`
- Bun 1.3.14 for runtime-profile experiments

The runtime version used to execute guest jobs is independent of the Node.js version used for
repository tooling.

Toolchain pins are part of the repository contract. Update their declarations, CI configuration,
and documentation together.

## Setup

```bash
corepack enable
pnpm install --frozen-lockfile
```

## Common commands

```bash
make fmt
make check
make test
make build
make ci
```

`make ci` runs the complete verification suite used for handoff. Run `make fmt` separately when
formatting changes are needed.

## Local AI context

Create or refresh the ignored local AI Central links:

```bash
pnpm codex:links
```

The command expects `../ai-central/templates` by default. Set `AI_CENTRAL_HOME` to the AI Central
repository or its `templates` directory when it lives elsewhere. Use
`pnpm codex:links -- --dry-run` to preview changes.

Repository-specific `AGENTS.md` and steering remain tracked. Shared steering and the full reviewed
AI Central skill catalog—including Caveman and Hallmark—remain local and are ignored by Git.

## Packages

- `@capsule-corp/protocol` contains TypeScript views of the canonical JSON protocol.
- `@capsule-corp/sdk` is a thin client for the trusted daemon.
- `@capsule-corp/mcp-server` will translate MCP tool calls into SDK requests.

JSON Schemas remain canonical. TypeScript and Go representations must be tested against schema
fixtures once code generation or validation is introduced.

## Local daemon

The scaffolded daemon exposes only informational endpoints:

```bash
go run ./cmd/capsuled
curl http://127.0.0.1:7777/healthz
curl http://127.0.0.1:7777/v1/version
curl http://127.0.0.1:7777/v1/runtimes
```

It does not execute guest code yet.

## Security-sensitive changes

Before changing capability resolution, path handling, process spawning, output collection, profile
resolution, or backend configuration:

1. State the security property being preserved.
2. Add positive and negative tests.
3. Add an attack-corpus case where appropriate.
4. Update the threat model or an ADR if the trust boundary changes.
