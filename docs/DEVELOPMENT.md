# Development

For a sequenced fresh-machine walkthrough (including Apple identity/credential setup and a
machine-loss recovery drill), see
[Development machine setup and recovery](DEV_MACHINE_SETUP_AND_RECOVERY.md). This page remains the
canonical reference for exact toolchain pins and common commands.

## Prerequisites

- Go 1.23 or newer, as declared by `go.mod`
- `golangci-lint` v2.12.2, matching the pin in `.github/workflows/ci.yml`, for `make check`
- Node.js 22.22.1, as declared by `.node-version`
- pnpm 10.28.2, as declared by `package.json`
- Bun 1.3.14 for runtime-profile experiments; it is not an admitted workload profile

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
  scaffold plus passive Phase 2A candidate views that no product endpoint consumes.
- `@capsule-corp/sdk` is a thin client for the trusted daemon.
- `@capsule-corp/mcp-server` will translate MCP tool calls into SDK requests.

Current JSON Schemas remain canonical for current tests, but they are not the intended v0 object
model. Do not extend the mixed `Job` capability union. The replacement must update schemas,
TypeScript/Go/Swift views, examples, SDK behavior, and shared fixtures together after the blocking
spikes. Phase 2A candidate schemas and Go/TypeScript decoded views now establish the first narrow
shape without activating it; see [Phase 2A contract foundation](PHASE_2A_CONTRACT_FOUNDATION.md)
and [Schema status](../schemas/README.md).

## Feasibility spike workflow

Blocking macOS, backend, cryptographic, content-handle, privilege, and update prototypes may be
one-off and are not required to resemble final product structure.

If spike code or one-time evidence is retained:

- place it in the separate
  [`Shrimpworks/capsule-experiments`](https://github.com/Shrimpworks/capsule-experiments)
  archive with a development-only README;
- state its hypothesis, exact environment, privileges/entitlements, and pass/fail criteria;
- include reproducible positive, negative, misuse, and failure-injection commands;
- prevent Capsule product packages from importing it;
- retain the canonical decision and exact commit-pinned evidence link in this repository;
- keep production conformance fixtures in `schemas/` or the owning product package;
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
- Disposable direct Apple Containerization development runs on exact supported Apple silicon/macOS
  versions
- Native libkrun/Hypervisor.framework development runs only after an exact pinned, signed,
  App-Sandboxed Apple-silicon profile closes its P0 readiness blockers
- OCI plus gVisor tests on Linux CI or a dedicated Linux worker

Apple Containerization may create a lightweight Linux VM for a development job, but Gate C failed
it as a production authority because a restarted Supervisor cannot durably identify or enumerate
the exact VM/helper through supported public APIs. Controller loss therefore remains unresolved;
the backend cannot produce ordinary success or release artifacts after ambiguity. No development
backend may execute an untrusted workload runtime directly on the host. The libkrun/HVF follow-up
is the lead native candidate under evaluation after conditionally passing its first isolation and
lifecycle corpus,
but its follow-up tracks found unresolved immutable block custody and a guest-visible `NullFs`
device.
The P0 reconciliation rejected stock and governed Bun for the first runtime, selected governed
`deno_core` as the engineering candidate without admitting it, narrows custody to the runtime root
for the first slice, and proposes bounded console ports for source/input and typed
inline JSON completion. Safe filesystem-image output parsing remains required before file artifacts,
while complete release-byte admission and composed validation remain open. The retained spike
runners are not a product development backend. Until one exact candidate profile passes the
complete documented corpus, receipts and UI must label every executable backend posture
`development`. See [Gate C P0 reconciliation](GATE_C_P0_RECONCILIATION.md).

## Security-sensitive changes

Before changing identity, signatures, approval, capability resolution, path handling, process
spawning, output collection, profile resolution, persistence, recovery, trusted IPC, updater/trust,
or backend configuration:

1. State the security property being preserved.
2. Add positive and negative tests.
3. Add an attack-corpus case where appropriate.
4. Update the threat model or an ADR if the trust boundary changes.
