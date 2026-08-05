# ADR-0039: License Capsule-owned material under Apache-2.0

- Status: Accepted
- Date: 2026-08-05

## Context

Capsule was published without an open-source license while its architecture and security boundary
were being established. That default allowed public inspection but did not grant permission to
use, modify, or redistribute Capsule-owned material. Before accepting outside contributions or
distributing an alpha, the project needs a well-understood license that supports independent
security review, downstream integration, commercial use, and clear contributor patent terms.

This decision applies to Capsule-owned material only. Governed upstream forks, vendored or copied
third-party material, guest/runtime packages, and generated bundles retain their own licenses and
notice or corresponding-source obligations.

## Decision

License Capsule-owned source code, documentation, schemas, and fixtures under the Apache License,
Version 2.0. Place the complete license at repository root, retain the project attribution in
`NOTICE`, identify Apache-2.0 in package metadata, and document the third-party boundary in
`docs/LICENSING.md`.

Section 5 of Apache-2.0 supplies the default inbound contribution terms unless a contributor
explicitly states otherwise. Capsule does not currently require a contributor license agreement
or copyright assignment. Package publication and product admission remain separate decisions;
adding a source license does not admit a runtime, backend, guest, control, or security posture.

## Alternatives considered

### MIT

MIT is short, permissive, and compatible with Capsule's upstream dependencies. It was not selected
because Apache-2.0 adds an explicit contributor patent grant and termination provision that are
valuable for a security and runtime project.

### MPL-2.0

MPL-2.0 would require distributed modifications to covered files to remain available while still
allowing combination with proprietary files. It was not selected because Capsule currently favors
broad integration and upstream adoption over file-level copyleft.

### GPL-3.0 or AGPL-3.0

Strong or network copyleft would preserve broader downstream source availability. It was not
selected because it would add avoidable integration and distribution constraints across Capsule's
SDK, native macOS components, governed runtime, and commercial adopters; AGPL's network condition
also does not align closely with the first local macOS release.

### A custom or source-available license

A custom license could reserve additional commercial rights but would increase legal review,
compatibility uncertainty, and contributor friction. It was rejected in favor of a standard
OSI-approved license.

## Consequences

- Anyone may use, modify, redistribute, and commercially integrate Capsule-owned material under
  Apache-2.0's conditions.
- Downstream users are not required to publish their modifications.
- Contributors provide the license's copyright and patent grants for accepted contributions.
- Redistributions must preserve the license, applicable notices, and modification notices.
- Apache-2.0 grants no general right to use Capsule or Shrimpworks names or marks.
- Every release still needs independent third-party license, notice, corresponding-source, and
  SBOM verification.
- A later proprietary dual-license program may require a new contribution agreement and cannot
  assume ownership of already accepted outside contributions.
