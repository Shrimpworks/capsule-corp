# Contributing

Capsule is in its architecture and threat-modeling phase. Changes should preserve the distinction
between requested authority and granted authority.

## Before contributing

1. Read the [project definition](docs/PROJECT.md).
2. Read the [architecture](docs/ARCHITECTURE.md).
3. Read the [threat model](docs/security/THREAT_MODEL.md).
4. Check the existing [architecture decisions](docs/adr/README.md).
5. For Go changes, read the [Go engineering standards](docs/GO_ENGINEERING_STANDARDS.md) and the
   [Capsule domain primer](docs/CAPSULE_DOMAIN_PRIMER.md).

## Development workflow

```bash
corepack enable
pnpm install --frozen-lockfile
make fmt
make ci
```

Every behavior change should include tests. Security-boundary changes should also add or update an
adversarial test case once that harness exists.

## Pull requests

Keep pull requests small and describe:

- The capability or component being changed.
- Which security property must continue to hold.
- How the change was verified.
- Any new ambient authority, network path, parser, dependency, or output channel.

Use focused commits and do not mix generated output or unrelated formatting into a behavioral
change. The pull request should pass `make ci` before review.

Do not weaken a default policy merely to make a workload pass. Introduce an explicit, reviewable
capability instead.

## Dependencies

New dependencies require justification. In particular, call out:

- Lifecycle scripts
- Native code or FFI
- Network behavior
- Dynamic loading
- Parsing of untrusted formats
- Transitive dependency size

## Architecture decisions

Material decisions should be recorded in `docs/adr/` using the existing format. Accepted ADRs are
append-only; supersede them with a new ADR rather than rewriting history.
