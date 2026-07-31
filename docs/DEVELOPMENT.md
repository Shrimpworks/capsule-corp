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

- `@capsule-corp/protocol` contains TypeScript views of the current canonical, pre-freeze JSON
  scaffold.
- `@capsule-corp/sdk` is a thin client for the trusted daemon.
- `@capsule-corp/mcp-server` will translate MCP tool calls into SDK requests.

Current JSON Schemas remain canonical for current tests, but they are not the intended v0 object
model. Do not extend the mixed `Job` capability union. The replacement must update schemas,
TypeScript/Go/Swift views, examples, SDK behavior, and shared fixtures together after the blocking
spikes. See [Schema status](../schemas/README.md).

## Feasibility spike workflow

Blocking macOS, backend, cryptographic, content-handle, privilege, and update prototypes may be
one-off and are not required to resemble final product structure.

If spike code is retained in the repository:

- place it under `experiments/` with a development-only README;
- state its hypothesis, exact environment, privileges/entitlements, and pass/fail criteria;
- include reproducible positive, negative, misuse, and failure-injection commands;
- prevent production packages from importing it;
- retain fixtures/evidence and record the ADR consequence;
- name the condition under which the prototype is removed or replaced.

A successful spike can conclude that a backend/control is unsuitable. See
[Feasibility Spikes](FEASIBILITY_SPIKES.md).

## Local daemon

The scaffolded daemon exposes only informational endpoints:

```bash
go run ./cmd/capsuled
curl http://127.0.0.1:7777/healthz
curl http://127.0.0.1:7777/v1/version
curl http://127.0.0.1:7777/v1/runtimes
```

It does not execute guest code yet.

## Backend development

Ordinary Go and TypeScript development does not require Linux. The planned backend workflow is:

- Native unit and contract tests on macOS and CI
- Registered-plan and recovery tests against a fake backend before hostile guest execution
- Disposable Apple Container capability probes on exact supported Apple silicon/macOS versions
- OCI plus gVisor tests on Linux CI or a dedicated Linux worker

Apple Container may create the Linux lightweight VM for an individual job; contributors do not
normally manage a permanent Linux VM. Its required no-network, resource, storage, management-
channel, orphan-recovery, and teardown semantics are not assumed until the spike proves them. No
development backend may execute untrusted Bun directly on the host. Until an exact backend
configuration passes the documented attack corpus, receipts and UI must label its isolation posture
`development`.

## Security-sensitive changes

Before changing identity, signatures, approval, capability resolution, path handling, process
spawning, output collection, profile resolution, persistence, recovery, trusted IPC, updater/trust,
or backend configuration:

1. State the security property being preserved.
2. Add positive and negative tests.
3. Add an attack-corpus case where appropriate.
4. Update the threat model or an ADR if the trust boundary changes.
