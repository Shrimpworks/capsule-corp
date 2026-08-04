# M1 ECMAScript module-request validator hold

Status: **HOLD — exact parser/validation boundary unselected and unimplemented**

This record defensively evaluates only fixed repository-owned `.mjs` byte fixtures with a passive
prototype. No fixture was executed as JavaScript, no runtime/backend/guest was launched, and no
network, package, credential, identity, or unrelated data was accessed.

## Question and method

ADR-0034 requires a bounded non-executing validator that distinguishes live ECMAScript module
requests from comments, strings, and ordinary grammar positions. A lightweight token scanner was
prototyped locally and fed exact UTF-8 fixture bytes. The prototype was never retained as product
or passive validation code.

The focused observations were:

| Exact source | Required interpretation | Prototype observation |
| --- | --- | --- |
| `obj.import.meta;` | property access, not `import.meta` | accepted |
| `({ import() {} });` | method named `import`, not dynamic import | rejected as a module request |
| ``const value = `${await import('./evil.mjs')}`;`` | live dynamic import in template interpolation | rejected |
| `eval("import('./evil.mjs')");` | string data at this boundary, not live import syntax | accepted |
| `const of = 9; of / import("evil") / divisor;` | live dynamic import between division operators | **accepted** |

The last case is the decisive false negative. The scanner treated contextual identifier `of` as a
grammar keyword that could precede a regular-expression literal, then consumed
`/ import("evil") /` as an opaque regexp. In this source, those slashes are division operators and
`import("evil")` is live dynamic-import syntax. Correct slash lexical goals require real grammar
context; patching individual preceding tokens would remain an ad hoc parser.

The exact retained counterexample is
[`language-hold-division-regexp-counterexample.mjs`](../schemas/conformance/v0/mjs-source/language-hold-division-regexp-counterexample.mjs):
45 bytes, SHA-256 `f5d7998154870fabd6f99ccc53c88e4f9edbe82bf50d63abd8e5db6138b663a1`.
The conformance manifest retains the other adjudication cases with all language-validator
implementations marked `pending`; fixture-integrity hashing is verified.

## Decision

- Do not retain or ship the prototype scanner.
- Do not add TypeScript or another broad parser dependency as an unreviewed shortcut.
- Do not weaken ADR-0034's no-module-request contract.
- Retain only the independently valid source-byte and deterministic-CBOR `SourceManifest` v0
  foundation.
- Hold JobProposal narrowing, semantic resolution/plan construction for the MJS profile, and all
  downstream M2/S1 activation until a separate decision selects an exact, pinned, governed,
  bounded ECMAScript parser/validation boundary.

`eval("import(...)")` demonstrates a separate property: syntax validation sees string data, while a
later runtime could interpret that string. Pre-approval syntax validation therefore cannot prove
runtime safety. ADR-0034's separately admitted no-loader runtime boundary must still make direct
and dynamically constructed module requests fail closed. This hold neither implements nor
substitutes for that runtime evidence.

## Retained scope and limitation

The passive foundation closes the exact media/profile/path identifiers, strict UTF-8 and leading
BOM handling, 0..262,144-byte identity, byte-preserving digest/length relationships, and the
single-member 87..95-byte SourceManifest. It does not narrow the existing public JobProposal,
construct an MJS ExecutionPlan, implement IPC/custody/approval/staging, or admit any runtime,
backend, or guest.
