# `bun-data@1`

Status: draft; not executable or approved.

This pre-freeze profile scaffold explores a later bounded Bun/TypeScript data environment with a
small, pinned, reviewed package set. It is **not** the first executable v0 bundle under the current
plan. Capsule must first prove a dependency-free JSON-in/JSON-out bundle and the profile evidence
model before activating curated packages.

The package versions in `profile.json` are intentionally marked `pending-review`. The object shape
itself may change during protocol/profile freeze. Before any later activation:

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
