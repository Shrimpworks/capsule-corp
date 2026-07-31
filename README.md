# Capsule

Capsule is an experimental, capability-controlled JavaScript and TypeScript task runtime for
AI agents.

[Project site](https://dills122.github.io/capsule-corp/) ·
[Technical design](docs/TECHNICAL_DESIGN.md) ·
[Architecture](docs/ARCHITECTURE.md) ·
[Threat model](docs/security/THREAT_MODEL.md)

An agent proposes source code, inputs, capabilities, resource limits, and expected outputs.
Capsule applies host policy, runs the task inside a disposable externally enforced sandbox, and
returns only policy-approved results and artifacts.

> [!WARNING]
> Capsule is an early scaffold. It is not yet a security boundary and must not be used to execute
> untrusted code outside an additional trusted sandbox.

## Core rule

> The model proposes work. The trusted host resolves the exact authority. The user approves that
> immutable plan. Capsule controls execution, egress, and the signed receipt.

## Current direction

- Runtime-neutral job protocol with a Bun-first implementation.
- Prepare, explicit human-readable approval, and execution bound to one immutable plan digest.
- Per-installation P-256 device identity with offline DID verification and purpose-separated keys.
- Local-first control experience with Apple Container on macOS and OCI plus gVisor on Linux as
  independent backend candidates.
- Immutable regular-file snapshots instead of live host paths or bind mounts.
- User-owned resource defaults and ceilings with no silent clamping.
- No ambient filesystem, network, environment, subprocess, native-addon, or package-installation
  authority.
- Separate input/capability and artifact/egress brokers.
- Disposable jobs with resource limits and execution receipts.
- MCP, CLI, and SDK adapters over the same protocol.

The initial vertical slice is intentionally narrow: a one-shot Bun/TypeScript task operating on
explicitly granted inputs with no network and returning validated artifacts. The architecture is
intended to support broader bounded JS/TS tasks later.

## Repository layout

```text
cmd/capsuled/          Go daemon entrypoint
internal/              Trusted control-plane implementation
packages/protocol/     TypeScript protocol types
packages/sdk/          TypeScript client SDK
packages/mcp-server/   MCP adapter scaffold
profiles/              Immutable runtime profile declarations
schemas/               Canonical JSON Schemas
docs/                  Project, architecture, security, ADRs, and roadmap
```

## Technology stack

- Go for the trusted daemon, policy, lifecycle, and isolation orchestration.
- TypeScript for the protocol package, SDK, CLI, MCP adapter, and initial guest programs.
- JSON Schema Draft 2020-12 for canonical wire contracts.
- SHA-256 content identities and signed, replay-resistant trust envelopes.
- P-256 per-device identity with Secure Enclave support on compatible Apple hardware.
- Bun as the first guest runtime, with Node and Deno deferred.
- Apple Container for on-demand macOS lightweight-VM execution and OCI plus gVisor as the Linux
  reference backend.
- pnpm, TypeScript, Biome, Go tooling, Make, and GitHub Actions for development and CI.
- Local AI Central links for reviewed skills and shared steering under `.codex/`.

Go is the initial control-plane language, not the security boundary. External isolation remains
mandatory. See [ADR-0005](docs/adr/0005-go-control-plane.md).

## Status

The repository currently contains specifications and buildable scaffolding. It does not yet start
guest runtimes.

Planned implementation order:

1. Freeze plan, identity, approval, capability, egress, receipt, error, and lifecycle contracts.
2. Build per-device trust, exact policy resolution, and immutable plan generation.
3. Prove the durable lifecycle and signed receipts with a fake backend.
4. Build a minimal Bun runner on the Apple Container development backend.
5. Complete the regular-file snapshot-to-artifact workflow through CLI and MCP.
6. Validate Apple Container and OCI plus gVisor independently against the attack corpus.
7. Add Node as the portability proof and Deno later.

See [Documentation](docs/README.md), [Project](docs/PROJECT.md),
[Technical design](docs/TECHNICAL_DESIGN.md), [Architecture](docs/ARCHITECTURE.md),
[Security model](docs/security/THREAT_MODEL.md), and [Roadmap](docs/ROADMAP.md).

## Development

Prerequisites:

- Go 1.23 or newer
- Node.js 22.22.1 for repository tooling
- pnpm 10.28.2
- Bun 1.3.14 for future runtime-profile work

```bash
pnpm install
make ci
```

See [Development](docs/DEVELOPMENT.md) for details.

Before publishing or transferring the repository, follow the
[GitHub repository setup checklist](docs/REPOSITORY_SETUP.md).

AI-agent context is documented in [.codex/AI_CENTRAL.md](.codex/AI_CENTRAL.md). The root
`AGENTS.md` remains authoritative for Capsule-specific security guidance. Run `pnpm codex:links` to
create the ignored local links, including Caveman and Hallmark.

## Support and security

Use [SUPPORT.md](SUPPORT.md) for development questions and bug-reporting guidance. Suspected
vulnerabilities must be reported privately according to [SECURITY.md](SECURITY.md).

## License

No license has been selected. Until a license is added, the repository is not offered under an
open-source license.
