# Rust engineering standards

This is the code-quality bar for Rust in this repository. It is a companion to `AGENTS.md`, not a
replacement — `AGENTS.md` governs security-sensitive behavior and wins on any conflict.

Scope: Rust currently exists only in `artifacts/mjs-source-validator-v1` and
`artifacts/mjs-source-validator-r2` — spike/experiment artifacts under `AGENTS.md`'s spike-code
rules, retained here as evidence rather than as a wired product component today (see
`docs/EXPERIMENT_ARCHIVE.md`). This document is the bar for that Rust whenever it's touched, and
the floor a future production Rust component inherits if one is ever admitted. It assumes standard
Rust idiom (`rustfmt` defaults, `clippy::pedantic`-adjacent judgment); where this document is
silent, follow the existing style in `artifacts/mjs-source-validator-r2/src/lib.rs`.

## Toolchain and dependencies

- Exact-pin every crate version in `Cargo.toml` (`oxc_parser = "=0.140.0"`, not `"0.140"` or
  `"^0.140.0"`). This Rust surface's entire job is parsing untrusted bytes (`.mjs` source); an
  unreviewed transitive upgrade landing silently is exactly the dependency-handling risk
  `AGENTS.md` treats as security-sensitive.
- Pin the toolchain via `rust-toolchain.toml` (`channel = "1.95.0"`), and build with `--locked
  --offline` (see `reproduce.sh` in both artifact directories) — never drop `--locked` to work
  around a lockfile mismatch; regenerate the lockfile deliberately instead.
- Keep `[profile.release]`'s `panic = "abort"`, `codegen-units = 1`, `lto = true`, and `strip =
  "symbols"` for anything built as a shipped binary artifact.
- Follow `AGENTS.md`'s `docs/ECOSYSTEM_REUSE_AND_ADOPTION.md` checklist before adding a crate, same
  as any other language here.

## Memory and error handling

- No `unsafe`. Zero `unsafe` blocks exist in this codebase today (verified) — a new one needs an
  ADR, not just a justifying comment, given `AGENTS.md`'s stance on reviewing new primitives.
- No `.unwrap()`/`.expect()` on attacker-controlled input. Only use them where the value is
  provably fixed-size or fixed-offset and the surrounding code proves it — `read_u16`'s
  `.expect("fixed offset")` in `lib.rs` is the model: it only ever slices a range the caller already
  length-checked, and the comment says so.
- Prefer `Result`/`Option` combinators over manual bounds checks: `checked_add`, `ok_or`,
  `map_err`, `try_fold`. `SyntaxCounts::increment`'s explicit `checked_add` overflow guard (rather
  than a bare `+= 1`) is the pattern for any counter fed by attacker-controlled repetition.
- Define fixed-size protocol frames as named constants (`REQUEST_HEADER_BYTES`,
  `RESULT_FRAME_BYTES` in `lib.rs`), check an explicit magic/version/kind/length header before
  touching any variable-length region, and reject on any mismatch rather than reading past what the
  header declared.

## Testing

- `#[cfg(test)]` modules driven off the shared `schemas/conformance/` fixture corpus (see `lib.rs`'s
  `corpus_root()`/`fixture()` helpers) rather than hand-written byte literals in the test itself — a
  new wire-format case gets a fixture file, not an inline array, so Go/TypeScript/Rust all exercise
  the same known answer.
- Cover both the accepted path and every rejection path (role mismatch, version mismatch,
  oversized/undersized frame, mutated binding) for anything parsing a request frame — pair a
  known-answer test like `emits_exact_role_specific_v1_known_answers` with a rejection test like
  `rejects_cross_role_v0_cap_plus_one_and_mutated_bindings`.

## Reproducibility

- Keep `SOURCE_DATE_EPOCH=0`, `TZ=UTC`, `LC_ALL=C`, and `--remap-path-prefix` set in any build
  script touching these binaries (see `reproduce.sh` in both artifact directories). A binary that
  isn't byte-reproducible across two independent clean builds fails the artifact's own evidence
  requirement, not just a style preference.
- A build script must independently build two clean copies and byte-compare them (`cmp`/`diff -r`
  in `reproduce.sh`) before anything is copied into `dist/`. Keep this gate on any future Rust build
  script; do not short-circuit it to "save time" locally.

## Precedence

If this document and `AGENTS.md` (or `docs/security/THREAT_MODEL.md`) conflict, `AGENTS.md` and the
threat model win. Raise the conflict as an issue or ADR rather than silently picking one.
