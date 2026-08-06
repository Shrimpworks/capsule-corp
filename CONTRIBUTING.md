# Contributing

Capsule is in its architecture and threat-modeling phase. Changes should preserve the distinction
between requested authority and granted authority.

## Before contributing

1. Read the [project definition](docs/PROJECT.md).
2. Read the [architecture](docs/ARCHITECTURE.md).
3. Read the [threat model](docs/security/THREAT_MODEL.md).
4. Check the existing [architecture decisions](docs/adr/README.md).
5. Read the [secure coding standards](docs/security/SECURE_CODING_STANDARDS.md) — cross-language
   practice that applies regardless of which language you're changing.
6. For Go changes, read the [Go engineering standards](docs/GO_ENGINEERING_STANDARDS.md) and the
   [Capsule domain primer](docs/CAPSULE_DOMAIN_PRIMER.md). For JavaScript/TypeScript changes, read
   the [JS/TS engineering standards](docs/JAVASCRIPT_TYPESCRIPT_ENGINEERING_STANDARDS.md). For Rust
   changes, read the [Rust engineering standards](docs/RUST_ENGINEERING_STANDARDS.md).
7. Read the [licensing policy](docs/LICENSING.md). Do not submit third-party or generated material
   without compatible provenance, licensing, and retained notices.

## Development workflow

```bash
corepack enable
pnpm install --frozen-lockfile
make fmt
make ci
```

Every behavior change should include tests. Security-boundary changes should also add or update an
adversarial test case once that harness exists.

## Opening issues

Use the narrowest issue form:

- **Bug report** for reproducible non-security incorrect behavior.
- **Capability proposal** for a new authority, runtime, backend, or workflow that still needs review.
- **Architecture or ADR change** for consequential trust-boundary, topology, or authority decisions.
- **Scoped implementation task** for accepted work with exact acceptance and verification criteria.
- **Research or experiment** for one falsifiable, defensively authorized question and retained evidence.
- **Blocked work** for an intended path waiting on a named unblock owner or source.

Suspected vulnerabilities must use
[private vulnerability reporting](https://github.com/Shrimpworks/capsule-corp/security/advisories/new),
not a public issue.

Issue labels use independent dimensions. Apply one primary work-type label (`bug`, `proposal`,
`implementation`, `research`, or `architecture`), the narrowest `component:...` label, and a
`status:...` label only when tracking the canonical status from
[`docs/STATUS_LANGUAGE.md`](docs/STATUS_LANGUAGE.md). `security-boundary` identifies public design,
hardening, or verification work; it is not an invitation to disclose a vulnerability publicly.
Use `proposal` while a change still needs a decision and `enhancement` only after its direction has
been accepted for implementation. See the complete [issue label taxonomy](.github/LABELS.md).

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

## Contribution license

Unless explicitly stated otherwise, contributions intentionally submitted for inclusion in
Capsule are licensed under Apache-2.0 under Section 5 of the project license. Submitting a change
does not transfer copyright ownership. No contributor license agreement or copyright assignment
is currently required.
