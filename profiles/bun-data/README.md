# `bun-data@1`

Status: draft; not executable or approved.

This profile is the first reference environment for bounded Bun/TypeScript data tasks. It will
contain a small, pinned, reviewed set of packages for CSV, JSON, JSONL, schema, date, and text work.

The package versions in `profile.json` are intentionally marked `pending-review`. Before activation:

1. Select exact versions.
2. Review lifecycle scripts, native code, maintainers, and transitive dependencies.
3. Generate and retain an SBOM.
4. Build the immutable runtime image.
5. Sign the image and record its digest.
6. Run profile conformance and adversarial tests.
7. Change the status to `active` only after the image digest is present.

Guest tasks cannot modify the package set or perform runtime installation.
