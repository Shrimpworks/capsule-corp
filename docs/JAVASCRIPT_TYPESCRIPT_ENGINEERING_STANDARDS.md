# JavaScript/TypeScript engineering standards

This is the code-quality bar for JavaScript and TypeScript in this repository: naming, structure,
error handling, testing, and lint hygiene. It is a companion to `AGENTS.md`, not a replacement —
`AGENTS.md` governs security-sensitive behavior and wins on any conflict. This document governs
everything `AGENTS.md` does not: the ordinary craft of writing and reviewing JS/TS here. It is the
JS/TS counterpart to `docs/GO_ENGINEERING_STANDARDS.md`; where the two overlap in spirit
(closed-world validation, fail-closed defaults, no hand-rolled crypto), they should stay in sync.

It assumes standard TypeScript/ECMAScript idiom. Where this document is silent, follow the
existing style in `packages/protocol/src`, the most mature package in the workspace. Where it is
specific, the specifics here win for this repository.

## Naming

- Filenames are kebab-case and match the primary export
  (`decoded-job-proposal-candidate.ts`, `execution-plan-builder.ts`). A file named `internal-*.ts`
  is package-private: never re-export it through a package's `index.ts` barrel. (This convention
  was broken once — see the tracked fix for `internal-contract-candidates.ts` — keep it followed
  going forward.)
- Prefer precise domain names (`JobProposalSchemaRefusalCode`, `TrustedExecutionPlanDigestRole`)
  over generic ones (`ErrorCode`, `Role`). This codebase's whole point is precise authority
  boundaries; a self-documenting name is worth the extra characters.
- A value that means something specific — a validated path, a positive integer, a 32-byte digest
  bound to one role — gets its own nominal branded type (`declare const xBrand: unique symbol`),
  not a bare `string`/`number`/`Uint8Array`. See `SourcePath`, `PositiveSafeInteger`, and
  `CandidateDigest<Role>` in `job-proposal.ts` and `internal-contract-candidates.ts`. A new value
  with the same "this string isn't just any string" property should follow the same pattern rather
  than staying primitive.

## Function and file shape

- One concern per file. The decode/resolve/build pipeline is split across
  `job-proposal-decoder.ts`, `job-proposal-resolver.ts`, and `execution-plan-builder.ts` on
  purpose — keep splitting by pipeline stage rather than growing one file per package.
  `internal-cbor-primitives.ts` isolates the one shared low-level encoder so every higher-level
  encoder doesn't reimplement CBOR head/length logic, and it stays unexported from `index.ts` on
  purpose — do not "simplify" by inlining it back into its callers or exporting it publicly.
- Small, single-purpose validation helpers (`requireObject`, `requireString`, `requireClosedFields`
  in `job-proposal-decoder.ts`) factored out and reused across every decode path, rather than
  repeating an inline `typeof`/`Array.isArray` check at each call site.
- Prefer a discriminated result object over a thrown exception for expected outcomes (see Errors
  below); reserve small internal exception classes (`SchemaFailure`, `RawDecodeFailure`,
  `ManifestDecodeFailure`) strictly as unwind mechanics inside one module, caught before the
  module's public boundary.

## Errors and refusals

- No thrown exception crosses a public API boundary for expected malformed input. Return a
  discriminated `{ ok: true, value } | { ok: false, refusal }` result instead — this is the pattern
  throughout `packages/protocol` (`JobProposalDecodeResult`, `ExecutionPlanConstructionResult`,
  `JobProposalResolutionResult`, and so on). Reserve `throw`/`TypeError` for programmer-error and
  invariant violations at construction time (a WeakSet provenance check failing, a closed-shape
  input check failing on a hand-built object) that a correct caller should never trigger.
- Every refusal carries a closed `classification` plus a closed `code`, both exhaustive TypeScript
  unions — never a free-text message a caller might string-match. A new refusal reason gets a new
  named union member, not a reused generic code.
- Closed-world by default: `requireClosedFields`'s "reject any key not in the admitted set" is the
  model. A validator that silently ignores an unrecognized field or falls through an unhandled
  union member to an implicit accept is a bug, not a shortcut.

## Provenance and immutability

This pattern is the strongest thing in this codebase and is not optional style — keep applying it.

- A value produced through validated construction is `Object.freeze`d before it's returned, and any
  array or typed-array member gets its own defensive copy (`Uint8Array.from`, `Array.from`) rather
  than aliasing the caller's buffer. See `retainRoleBytes`/`retainExactBytes` in
  `execution-plan-builder.ts` and the equivalent copies in `job-proposal-resolver.ts`.
- When a value's *origin* matters — only the decoder may produce a "decoded candidate," only the
  resolver may produce "resolved plan inputs" — brand it structurally **and** register its object
  identity in a module-local `WeakSet`/`WeakMap`, then have every downstream consumer check
  `isRetained*`/`isTrusted*` before trusting it. See `decoded-job-proposal-candidate.ts`,
  `resolved-job-proposal-provenance.ts`, and the binding checks throughout
  `execution-plan-builder.ts`. A same-shaped object built by hand must still be rejected — shape can
  be forged, `WeakSet` membership cannot. Do not skip the provenance check "because the shape
  already proves it."

## Testing

- `node:test`, run against built output (`node --test dist/*.test.js`), not source directly — this
  catches ESM/type-emit issues that testing straight from `.ts` source would miss.
- Byte-exact/conformance-fixture testing for anything that encodes or decodes a wire format.
  Fixtures live under `schemas/conformance/`, are produced by `scripts/generate-*-fixtures.mjs`, and
  are independently re-verified by a paired `scripts/verify-*.mjs`. New wire-format code needs a
  fixture generator plus an independent verifier, not just inline assertions — see the
  `mjs-source-validator` and `contract-conformance` generator/verifier pairs.
- Pair every accepted fixture with several rejected mutations. `scripts/verify-schemas.mjs`'s
  `invalidProposals`/`invalidC2A` arrays are the model: each valid example is accompanied by
  deliberate mutations (wrong version, wrong slot, oversized value, fractional number where an
  integer is required) that must be rejected. A new schema or decoder needs both directions.

## Documentation

- Every exported type and function gets a `/** ... */` doc comment. Where behavior follows from an
  ADR or design doc, cite it inline — see `execution-plan-builder.ts`'s ADR-0019/ADR-0023
  references and the `profileReviewAttestationDigestsMatch` comment explaining *why*
  order-independent comparison is correct, not just restating what the code does.
- Use the exact vocabulary from `docs/STATUS_LANGUAGE.md` (`PASSED`, `IN_PROGRESS —
  TRENDING_GOOD/BAD`, `BLOCKED`, `NO_GO`) in comments describing project status — same rule as the
  Go standards document, for the same reason: it keeps code comments and project docs from quietly
  disagreeing.

## Lint and CI wiring

- `biome check .` (`pnpm lint`) is mandatory. It currently runs only the `"recommended"` rule
  preset in `biome.json` — tightening this is tracked separately. Until it's tightened, treat
  avoiding `any`, non-null assertions, and unawaited promises as unwritten team convention this
  codebase already follows throughout `packages/protocol`, `packages/sdk`, and `packages/mcp-server`
  — do not introduce the first violation.
- `tsc --noEmit` (`pnpm check`) against the shared `tsconfig.base.json`: `strict`,
  `noUncheckedIndexedAccess`, `verbatimModuleSyntax`, and `forceConsistentCasingInFileNames` all
  stay on. Don't loosen a package's own `tsconfig.json` to work around a type error — fix the type,
  or if the check is genuinely wrong for a specific case, say why in a comment at the suppression
  site.
- `pnpm verify:schemas`, `pnpm verify:adrs`, and the conformance-fixture generators run in CI
  (`.github/workflows/ci.yml`). A new schema or protocol object needs its own entry in that chain,
  not just a passing `tsc`.

## Security-adjacent JS/TS hygiene

These are code-quality rules, not new security policy — `docs/security/THREAT_MODEL.md`,
`AGENTS.md`, and `docs/security/SECURE_CODING_STANDARDS.md` remain the authority on the latter.

- Never build a filesystem or module path from untrusted input. `scripts/verify-field-authority-
  manifest.mjs`'s `resolveDefinitionPath` (`path.resolve` plus an explicit `startsWith(root)` check)
  is the pattern for any script that reads a caller-supplied relative path.
- Bound every parser before it recurses or accumulates, not after: max depth, max node count, max
  collection size, max string/byte length. See `strict-job-proposal-json.ts`'s
  `JOB_PROPOSAL_RAW_LIMITS` and the per-field checks enforced during parsing, not post-hoc. An
  unbounded parser over untrusted bytes is a denial-of-service surface even with zero memory-safety
  bugs.
- No `eval`, `new Function(...)`, or dynamic `import()`/`require()` of a non-literal specifier
  anywhere in product code. (A fixture containing the literal text `"eval(...)"` as a rejected-input
  test case is fine — that's inert data, never executed.)
- No hand-rolled cryptography. `node:crypto`'s `createHash("sha256")` is the only primitive this
  codebase uses today; if a design ever needs signatures, follow ADR-0019's COSE_Sign1 direction
  through an adopted, reviewed dependency — never a bespoke scheme.
- Defensive-copy anything crossing a trust boundary (see Provenance above). Never return or accept
  a live reference to a `Uint8Array` the other side could mutate after the fact.
- Treat every parsed JSON value as `unknown` until a validator narrows it — the `isRecord` plus
  `assertHealthResponse`/`assertVersionInfo`/`assertRuntimesResponse` pattern in
  `packages/sdk/src/index.ts` is the model for validating a network/daemon response before trusting
  its shape. Never narrow with a bare `as` on external input.
- No `as any` and no non-null assertion (`!`) as a substitute for an actual runtime check. Both are
  currently absent from this codebase — keep it that way. `noUncheckedIndexedAccess` is on
  specifically so an out-of-bounds access is a type error, not a silent `undefined`; don't work
  around it with `!`.

## Dependencies

Before adding one, follow `AGENTS.md`'s requirement to consult
`docs/ECOSYSTEM_REUSE_AND_ADOPTION.md` and complete its checklist, same as for Go.

- Root `devDependencies` are exactly `@biomejs/biome`, `@types/node`, `ajv`, `ajv-formats`, and
  `typescript` today — a near-zero-dependency posture matching the Go side's own "two direct
  dependencies" discipline. A new runtime dependency inside a published `packages/*` package is a
  harder sell than a root tooling `devDependency`.
- A dependency that ships a lifecycle script (`postinstall`, `preinstall`) needs explicit review
  before adoption — it runs with the same authority as the build itself. See
  `docs/security/SECURE_CODING_STANDARDS.md`'s dependency supply-chain section.
- Pin exact versions for anything that ends up inside a published package's runtime surface;
  `pnpm-lock.yaml` is the actual authority and CI installs with `--frozen-lockfile` — never commit a
  change that would make a fresh `--frozen-lockfile` install fail.

## Precedence

If this document and `AGENTS.md` (or `docs/security/THREAT_MODEL.md`) conflict, `AGENTS.md` and the
threat model win. Raise the conflict as an issue or ADR rather than silently picking one.
