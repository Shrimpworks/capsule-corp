# GitHub Repository Setup

This checklist covers settings that cannot be enforced by files in the repository. Complete it
before a public preview and review it whenever ownership or release policy changes.

## Repository identity

- Choose the permanent GitHub organization and repository name.
- Replace the provisional `capsule.local/capsule` Go module path before publishing Go modules.
- Add the repository description, topics, and project URL.
- Apache-2.0 is selected for Capsule-owned material. Preserve `LICENSE`, `NOTICE`, the documented
  third-party boundary, and release-specific corresponding-source and notice obligations.
- Adopt a code of conduct and private enforcement contact before accepting external contributions.
- Keep the canonical `capsule-corp` repository's default branch named `main`.

## Branch protection

Protect `main` with:

- Pull requests required for changes.
- The `Go` and `TypeScript` CI jobs required to pass.
- Review conversations required to be resolved.
- Force pushes and branch deletion disabled.
- Administrator bypass limited and auditable.

Add `.github/CODEOWNERS` after the responsible GitHub team or maintainers are known. Require review
from those owners for `schemas/`, `internal/`, `profiles/`, `docs/TECHNICAL_DESIGN.md`,
`docs/security/`, and `docs/adr/`.

Runtime, guest-kernel, VMM, launcher, security-profile, and security-critical native patch changes
require at least one independent approval from a designated security owner; the author cannot
satisfy that approval. If no independent qualified reviewer is available, experiments and evidence
may record the limitation, but the changed bytes cannot enter an admitted runtime bundle or
support a new validation claim. Preserve exact reviewed patch/source digests and rerun deliberate
capability-restoration tests after any affected rebase or rebuild.

## Governed upstream forks

Do not reuse the canonical repository's `main` convention blindly for owned dependency forks.
Configure each fork according to its recorded role:

- A fork primarily carrying Capsule-governed bytes uses its latest accepted governed line as the
  default branch. Its upstream-oriented `main` remains integration state and does not imply Capsule
  adoption.
- Preserve the upstream anchor, original governed patch-queue merge, and every superseded accepted
  head under explicit versioned refs. Lock immutable refs with administrator enforcement; disable
  force-pushes and deletion.
- Start each governed update from the preceding accepted commit on a fresh versioned target branch.
  Require a pull request, keep it draft through evidence and human review, then lock the accepted
  head before switching the fork default.
- Never merge an upstream development branch into a pinned governed line. Backport logical commits
  separately, resolve version-specific conflicts without widening behavior, and rerun the governed
  source, unit, mutation, sanitizer, coverage, architecture, and authorized guest gates that apply.
- Upstream contribution preparation uses an explicitly named upstream baseline or integration
  branch. It cannot change Capsule product state without a later governed backport.
- Every automation or agent must create pull requests with explicit repository, base, and head
  arguments and read those fields back with the stored body and draft state.

When an owned fork has only one qualified maintainer, protect its mutable default and governed
review targets with required pull requests, green required checks, resolved conversations,
administrator enforcement, no force-push, no deletion, and explicit evidence/settings readback.
Set the required approving-review count to zero, require-approval-of-most-recent-push to false, and
required CODEOWNER review to false; the maintainer must perform and record a self-review. Enable
external approval enforcement only when a second qualified maintainer is available. This GitHub
merge setting does not waive the independent security approval required above for admitted runtime
bytes or new validation claims, and it does not waive DCO or upstream-submission authorization.

Record the active default, immutable ref names and commits, candidate target, protections, and
verification commands in the fork's governed metadata. A branch name alone is not an immutable
identity.

### Accepted governed fork readback

The 2026-08-05 post-merge acceptance transition independently verified each merge's parents,
ancestry, and tree before changing GitHub settings. These are the current accepted defaults:

| Fork | Immutable upstream source | Locked accepted default | Accepted tree | Upstream-integration `main` |
| --- | --- | --- | --- | --- |
| [`Shrimpworks/deno`](https://github.com/Shrimpworks/deno) | `capsule/anchor-v2.9.4` / `14eea3160ae5834476aa3b9d317b8d41d991b982` | `capsule/accepted-v2.9.4-r3` / `3fa21d1ae7705ab4bcb4bc98955f25301b20122a` | `6060cb0eb4cd3395a4c141f054634968744617d2` | `98bcf8375eb9f9a8fa7d83fb2f7885ef38244219` |
| [`Shrimpworks/rusty_v8`](https://github.com/Shrimpworks/rusty_v8) | `capsule/anchor-v150.2.0` / `d305e6afa7736f6e298c30ae6646f7709ee9382b` | `capsule/accepted-v150.2.0-r5` / `d09221062280ae1675fe26c53c3f43871aae2055` | `2632901e6e7e9ac88662756ceb658d4e3e49fceb` | `42d6a1adc3b6ab97eff922638c1e6c3c847a8800` |
| [`Shrimpworks/libkrun`](https://github.com/Shrimpworks/libkrun) | upstream v1.19.4 / `728df8125077d0db44265f6e997c72b81b65c015` | `capsule/upstream-v1.19.4-r3` / `7432eda5a49220976b0167005aa43ee622f9d632` | `7671440cfbafa58fe20aebf8d4deb2a843ebe346` | `1622c9f46751fc8341555f44d2000589e5d28360` |

The Deno and `rusty_v8` completed r3/r5 review targets are locked at the same accepted commits.
The libkrun original patch-queue merge remains locked as `capsule/baseline-v1.19.4-r1` at
`4ea8d1de861ed1c0636fc800b6da8fb71a086aa5`; its superseded accepted line remains locked as
`capsule/upstream-v1.19.4` at `cf0333cdba478cc34a8570a65b38412da7fd3ecc`. The earlier Deno and
`rusty_v8` anchor, reviewed-head, accepted, and recovery refs likewise remain protected and
unchanged.

Settings readback showed every accepted default and completed review target locked with
administrator enforcement, required pull requests, resolved conversations, no force-push, and no
deletion. All three `main` branches retain the same protections but remain unlocked integration
state; this transition added the previously missing protection to libkrun `main`. No required
status-check contexts were configured on these branches at readback. Required approvals are zero,
most-recent-push approval is false, and CODEOWNER approval is false while there is one qualified
maintainer. This is repository governance only: it creates no release and changes no dependency,
runtime, backend, profile, control-evidence, or product-admission state.

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
5. Confirm Apache-2.0 metadata, `LICENSE`, `NOTICE`, contribution expectations, and all
   third-party notices are current.
6. Confirm local AI Central symlinks are ignored and `pnpm codex:links -- --dry-run` documents how
   contributors recreate them.
