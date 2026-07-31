# Testing And Quality Gates

Testing should protect behavior, contracts, and integration boundaries.

## Default Expectations

- Add or update focused tests for behavior changes.
- Cover edge cases for parsing, validation, permissions, persistence, and external integrations.
- Keep test fixtures small and explicit.
- Prefer deterministic tests over timing-sensitive assertions.

## Before Finishing Work

Run the smallest reliable command that validates the changed area:

- Lint and static analysis: `make check`
- Unit tests: `make test`
- Schema contracts: `pnpm verify:schemas`
- Build and typecheck: `make build`
- Complete handoff suite: `make ci`

If a command cannot run locally, document why and what risk remains.

## Quality Gates

- No known failing tests introduced by the change.
- No unrelated formatting churn.
- Public contracts updated when behavior changes.
- Docs updated for setup, command, or workflow changes.
- Security-boundary changes include positive, negative, and adversarial tests where applicable.
