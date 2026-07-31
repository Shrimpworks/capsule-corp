# `bun-data@1`

Status: draft; not executable or approved.

This profile is the first reference environment for bounded Bun/TypeScript data tasks. It will
contain a small, pinned, reviewed set of packages for CSV, JSON, JSONL, schema, date, and text work.

The package versions in `profile.json` are intentionally marked `pending-review`. Before activation:

1. Select exact versions.
2. Review lifecycle scripts, native code, maintainers, and transitive dependencies.
3. Generate and retain an SBOM.
4. Build the immutable runtime image.
5. Sign the profile and image manifests and record their digests.
6. Record a separate review attestation; a publisher signature alone does not approve the profile.
7. Pin backend-specific kernel and init identities where applicable.
8. Run profile conformance and adversarial tests.
9. Change the status to `active` only after immutable identities and retained evidence are present.

Guest tasks cannot modify the package set or perform runtime installation.

Profile defaults and maximums do not silently override user policy. During planning, Capsule
resolves missing values from trusted user defaults, rejects requests above user ceilings, and asks
the backend to enforce the exact approved plan.
