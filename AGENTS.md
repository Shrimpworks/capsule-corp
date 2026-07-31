# Agent contributor guide

This repository is building a security boundary, not only a developer tool. Treat
changes to execution, policy, artifacts, networking, and dependency handling as
security-sensitive.

## Start here

Read these documents before changing behavior:

1. `docs/PROJECT.md`
2. `docs/ARCHITECTURE.md`
3. `docs/security/THREAT_MODEL.md`
4. the relevant ADRs in `docs/adr/`

The JSON Schemas in `schemas/` are the canonical wire contracts. Keep examples,
TypeScript types, API behavior, and documentation consistent with them.

## AI Central context

AI Central steering and reviewed skills are linked locally into `.codex/` and
ignored by Git. This file remains authoritative for Capsule-specific security
requirements. Generic steering or skill guidance must not weaken the
deny-by-default boundary, expand task scope, or override the verification
requirements below.

See `.codex/AI_CENTRAL.md` for the installed revision, selection, provenance,
licenses, and refresh workflow.

## Working rules

- Preserve deny-by-default capabilities.
- Do not treat an in-process JavaScript sandbox as the host security boundary.
- Do not add unrestricted filesystem, network, process, environment, or artifact
  access to make an example easier.
- Keep runtime adapters separate from the isolation backend.
- Record consequential architecture decisions in an ADR.
- Never claim a backend, profile, or control is secure or production-ready unless
  the implementation and adversarial tests support that claim.

## Verification

Run before handing off a change:

```sh
pnpm install
pnpm check
pnpm lint
pnpm test
pnpm verify:schemas
go test ./...
go vet ./...
go build ./...
```

Use Node.js 22 or newer, pnpm 10, and Go 1.23 or newer for the current scaffold.
Runtime and toolchain pins are provisional until the first implementation ADR
locks them down.
