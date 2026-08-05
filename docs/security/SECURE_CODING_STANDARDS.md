# Secure coding standards

This is the cross-language engineering practice that keeps code from being the reason a documented
mitigation fails in reality. It is a companion to `AGENTS.md` and `docs/security/THREAT_MODEL.md`,
not a replacement for either — they are policy and adversarial analysis (the *what* and *why*); this
document is engineering discipline (the *how*), the same relationship the per-language standards
documents have to `AGENTS.md`.

Every language in this repository touches the security boundary somewhere: Go is the trusted
control plane, TypeScript owns the wire protocol and client SDK, Rust builds the shipped Source
Validator XPC parser binaries, and Swift builds the shipped Broker/status application shell. None of
them get a pass on secure coding practice because their language isn't the "main" one.

## How this fits with everything else

- `AGENTS.md` — binding project-wide security rules. Wins on any conflict with anything below.
- `docs/security/THREAT_MODEL.md` — adversaries, trust boundaries, attack surface, mandatory
  security properties: the *what could go wrong*.
- This document — cross-language coding practice: the *how do we keep it from going wrong here*.
- Per-language detail and worked examples:
  - `docs/GO_ENGINEERING_STANDARDS.md` § Security-adjacent Go hygiene
  - `docs/JAVASCRIPT_TYPESCRIPT_ENGINEERING_STANDARDS.md` § Security-adjacent JS/TS hygiene
  - `docs/RUST_ENGINEERING_STANDARDS.md`
  - Swift has no dedicated engineering-standards document yet — see Known gaps below. The
    cross-language principles here still apply to
    the [archived I1A `CapsuleStatusApp.swift`](https://github.com/Shrimpworks/capsule-experiments/blob/0944ffd8cfd01ec23e4ae99138b0931d56804077/experiments/completed-compiled-artifact-payloads/payloads/capsule-corp/artifacts/macos-i1a-unsigned-app-shell/Sources/CapsuleStatusApp.swift)
    and any future Broker Swift code.

## Cross-language principles

1. **Closed-world validation, not open-world.** Reject unknown fields, variants, and types by
   default; every accept path is an explicit allowlist. Go's exhaustive `Classification` switches
   with no permissive `default:` case, TypeScript's `requireClosedFields`, and Rust's exact
   frame-length/magic checks are the same rule in three languages. A `default: allow` anywhere near
   a trust boundary is a bug, not a style choice.
2. **Bound everything that touches attacker-controlled bytes before parsing it, not after.** Max
   length, max depth, max element count, checked at the earliest possible point —
   `strict-job-proposal-json.ts`'s `JOB_PROPOSAL_RAW_LIMITS`, the Rust Source Validator's fixed
   frame sizes, and Go's byte-exact conformance ceilings are the same control in three languages. An
   unbounded parser over untrusted input is a denial-of-service bug even with zero memory-safety
   issues.
3. **No ambient authority.** A component gets only the filesystem/network/process/environment
   access its role requires, enforced in code, not by convention: the `os/exec` policy layer in Go,
   "never pass live host paths into a guest," and the Rust validator binaries taking stdin bytes
   only with no argv-controlled paths are all the same principle. Never widen this "to make an
   example easier" — `AGENTS.md` says so explicitly, and it is a code-review blocker, not a
   discussion point.
4. **Fail closed on ambiguity.** An error, a malformed input, or an unrecognized state is a
   refusal, never a default-accept. Same rule in every language's error-handling section.
5. **No hand-rolled cryptography, in any language.** `node:crypto`, Go's standard `crypto/*`
   package, or an already-adopted, reviewed dependency per `docs/ECOSYSTEM_REUSE_AND_ADOPTION.md`.
   This includes "just a hash for non-security purposes" — if the digest later gates authority or
   identity, treat computing it as security code from the start.
6. **Defensive copies across trust boundaries.** Once a byte buffer crosses from caller-owned to
   retained, copy it — never alias memory the other side could still mutate. Applies identically to
   Go slices, TypeScript `Uint8Array`s, and Rust `&[u8]` → `Vec<u8>` conversions.
7. **Provenance over shape.** When only one code path may produce a "trusted" value, prove *origin*
   — a `WeakSet` identity check in TypeScript, a sealed/branded type in Go, a typed builder in
   Rust/Swift — not just shape. A hand-constructed value with the right fields must still be
   rejected; shape can be forged, provenance cannot. This is the strongest pattern already present
   in this codebase (see `packages/protocol`'s provenance `WeakSet`s) — keep applying it as new
   trust-sensitive types appear in any language.
8. **Secrets and user content never reach a default log or print.** Redact by construction — no
   default `%v`/`Debug`/`console.log`-friendly dump of any type that might carry a key, approval
   content, or job content — rather than trusting every call site to remember not to log it.
9. **Reproducible, pinned builds for anything shipped as a binary.** Exact-pinned dependency
   versions, a pinned toolchain, `SOURCE_DATE_EPOCH`, and an actual two-build byte-identity check
   before anything reaches `dist/`. This already happens for the Rust and Swift artifacts
   (`reproduce.sh` in each); hold Go and TypeScript build output to the same bar once either ships a
   binary artifact.
10. **A lint or security-analysis suppression always carries an inline rationale at the suppression
    site.** No blanket-disabled rule anywhere, in any language's lint configuration.

## Dependency supply chain

- Every new dependency, in every language, goes through `docs/ECOSYSTEM_REUSE_AND_ADOPTION.md`'s
  checklist before it's added. This is already required by `AGENTS.md`; it is restated here because
  it is the single highest-leverage security control available to a project this size — a small
  footprint with few dependencies is a property worth actively protecting, not a side effect.
- A dependency that carries a lifecycle or build script (npm `postinstall`/`preinstall`, a Cargo
  build script, a Go `go generate` hook pulled from a dependency) gets explicit review before
  adoption — it runs with the same authority as the build itself.
- Prefer exact pins for anything that ends up inside a shipped artifact (`=x.y.z` in `Cargo.toml`,
  exact versions where the ecosystem supports it). Range pins are acceptable for tooling-only
  dependencies, but the lockfile (`pnpm-lock.yaml`, `Cargo.lock`, `go.sum`) is the actual authority
  and must be committed and enforced (`--frozen-lockfile`, `--locked`).
- Build and CI tooling is part of this surface too, not a separate category: pin GitHub Actions to
  a commit rather than a moving tag, and pin one-off tool invocations (like `govulncheck`) to a
  version instead of `@latest`. An unpinned build-time tool has the same blast radius as an
  unpinned runtime dependency — it just runs earlier.
- `docs/security/THREAT_MODEL.md`'s attack-surface enumeration does not yet have its own "build and
  dependency supply chain" subsection distinct from runtime content/egress threats; this section is
  interim coverage for that gap. Folding a proper attacker story into the threat model itself is a
  reasonable follow-up, but it should go through the same review the rest of that document gets, not
  land silently as a side effect of writing this one.

## Known gaps (tracked, not hidden)

Naming these explicitly is the point of writing them down — a steering document that hides its own
blind spots is worse than none.

- Swift has no dedicated engineering-standards document. The cross-language principles above apply,
  but nothing yet captures Swift-specific hygiene (force-unwrap avoidance, `Codable` boundary
  validation, memory ownership at the AppKit/XPC boundary) the way the Go/JS/Rust documents do for
  their languages.
- Biome (JS/TS) currently runs only its `"recommended"` rule preset — there is no JS/TS equivalent
  of `gosec`'s enforcement, and no static security scanner (semgrep, `eslint-plugin-security`, or
  similar) wired into CI for this language, where Go has `gosec` via `golangci-lint`.
- This document and the per-language ones describe practice observed in the codebase as of this
  writing. As the codebase grows, re-verify the concrete examples cited here still exist and still
  demonstrate the rule — a stale worked example is worse than none, because it reads as endorsed
  when it may no longer be current.

## Precedence

`AGENTS.md` wins over `docs/security/THREAT_MODEL.md`, which wins over this document, which wins
over the per-language engineering-standards documents, in that order, on any conflict. Raise a
conflict as an issue or ADR rather than silently picking one.
