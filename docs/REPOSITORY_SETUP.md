# GitHub Repository Setup

This checklist covers settings that cannot be enforced by files in the repository. Complete it
before a public preview and review it whenever ownership or release policy changes.

## Repository identity

- Choose the permanent GitHub organization and repository name.
- Replace the provisional `capsule.local/capsule` Go module path before publishing Go modules.
- Add the repository description, topics, and project URL.
- Choose a license before accepting external contributions or distributing releases.
- Adopt a code of conduct and private enforcement contact before accepting external contributions.
- Keep the default branch named `main`.

## Branch protection

Protect `main` with:

- Pull requests required for changes.
- The `Go` and `TypeScript` CI jobs required to pass.
- Review conversations required to be resolved.
- Force pushes and branch deletion disabled.
- Administrator bypass limited and auditable.

Add `.github/CODEOWNERS` after the responsible GitHub team or maintainers are known. Require review
from those owners for `schemas/`, `internal/`, `profiles/`, `docs/security/`, and `docs/adr/`.

## Security settings

Enable:

- Private vulnerability reporting.
- Dependabot alerts and security updates.
- Secret scanning and push protection.
- GitHub Actions with read-only default token permissions.

Review allowed GitHub Actions and pin third-party actions to full-length commit SHAs as the workflow
supply-chain policy is finalized.

Do not describe Capsule as a security boundary in repository metadata or releases until an
authoritative backend has passed the documented adversarial corpus.

## Merge and release policy

- Prefer squash merging for a readable early history.
- Delete merged branches automatically.
- Do not publish packages, images, or binaries until provenance, signing, and versioning policies
  are defined.
- Do not activate a runtime profile until its immutable image, digest, dependency review, SBOM,
  signature, conformance tests, and adversarial tests are complete.

## Public repository review

Before changing repository visibility:

1. Run `make ci` from a clean checkout.
2. Confirm no credentials, local artifacts, `.env` files, signing material, or private assessment
   details are present in Git history.
3. Verify `SECURITY.md` points to an enabled private reporting path.
4. Confirm the README, support policy, threat model, and security status accurately describe the
   current implementation.
5. Confirm the selected license and contribution expectations.
6. Confirm local AI Central symlinks are ignored and `pnpm codex:links -- --dry-run` documents how
   contributors recreate them.
