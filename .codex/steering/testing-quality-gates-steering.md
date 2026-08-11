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

## Durable Experiment Evidence

Important evidence-bearing work includes security-boundary, hostile-guest, platform, signing,
durability, fault-injection, performance-gate, and product-admission experiments. Its raw evidence
must not exist only in `/tmp`, `/private/tmp`, an agent worktree, a disposable guest, or a
short-lived CI artifact.

Before execution, name the evidence owner, staging location, expected manifest, and durable archive
destination. During the run, preserve the exact harness source, environment and immutable input
identities, commands, reports, receipts, captures, cleanup evidence, file sizes, and SHA-256 hashes
needed for independent readback. Exclude credentials, private keys, raw secrets, and unrelated
user or host data.

Before task completion or destructive cleanup:

1. copy the complete evidence packet to an owner-controlled non-ephemeral location;
2. verify its manifest and hashes independently;
3. publish disposable harnesses and raw evidence to
   `Shrimpworks/capsule-experiments` through a reviewed immutable commit;
4. read the remote commit and tree back and rerun the archive verifier; and
5. record the exact archive commit and narrowly scoped conclusion in Capsule's canonical docs.

Do not delete the source workspace or describe evidence as durable, retained, release-grade, or
admission-grade until the remote archive readback succeeds. If backup, publication, or verification
fails, preserve every available copy, report the evidence-retention step as `BLOCKED`, and keep any
product or admission claim blocked. Follow the complete repository/archive boundary and publishing
workflow in [Experiment archive](../../docs/EXPERIMENT_ARCHIVE.md).
