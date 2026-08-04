# `.mjs` module-request validation boundary — independent defensive review

Date: 2026-08-03

Status: retained independent design review. No validator, parser dependency, or ADR is added or
selected by this review. It answers one question for the active M1 task under Accepted ADR-0034:
whether a bounded validator can reliably reject imports and prohibited module features in a
single-file `main.mjs` without executing source or embedding a broad JavaScript
parser/runtime dependency.

Reviewer: Claude, independent read-only defensive review at the request of the Capsule
orchestrator (codex).

## Defensive scope and method

Read-only review. No file in this repository was modified to produce it, no untrusted workload was
executed, and no external system was accessed. The review read, in order: `AGENTS.md`,
[ADR-0034](adr/0034-freeze-mjs-first-release-contract.md),
[`TYPESCRIPT_APPROVED_BYTE_CUTOVER_PLAN.md`](TYPESCRIPT_APPROVED_BYTE_CUTOVER_PLAN.md),
[`security/THREAT_MODEL.md`](security/THREAT_MODEL.md),
[`protocol/OBJECT_MODEL.md`](protocol/OBJECT_MODEL.md),
[ADR-0028](adr/0028-select-governed-deno-core-first.md),
[`GOVERNED_DENO_CORE_WORK_PLAN.md`](GOVERNED_DENO_CORE_WORK_PLAN.md), `PROJECT.md`/`ARCHITECTURE.md`
(daemon-language sections), `go.mod`, and the `internal/`/`packages/`/`cmd/` trees, grepped for any
existing validator/scanner/tokenizer code.

**Grounding fact:** no module-request validator exists in the repository yet. `grep` for
`import.meta`, `module-request`, `denylist`, `require(`, `CommonJS`, `scanner`, `tokenizer`,
`lexer` across `internal/` and `packages/` turns up nothing but Node/Vitest test-fixture
boilerplate and unrelated CBOR/owner-lock code. ADR-0034 itself lists this as an open blocker: *"A
bounded non-executing module-request validator, Supervisor source-store projection, independent
Broker validation/rendering, source transport framing, and runtime no-loader evidence remain
unimplemented."* This review is therefore prospective design analysis, not a code audit.

Second grounding fact: the component that would run this validator — *"The agent-facing Go daemon
performs strict proposal validation and planning only"* (`PROJECT.md`) — is written in Go, and
`go.mod` currently declares exactly one dependency, `golang.org/x/sys`. No pinned ECMAScript parser
exists anywhere in this repo's Go code today. The only ECMAScript-aware tooling present is the
`typescript` npm devDependency used by `packages/protocol`, `packages/sdk`, and `packages/mcp-server`
— not something the Go daemon can call in-process today.

## Decision

**A bounded/conservative substring- or regex-style scanner is not defensible and should not be
built.** The counterexamples below are not adversarial-only; several (regex literals containing
import-shaped text, `obj.import.meta`, `{ import() {} }`, comments/strings with lookalike text) are
things ordinary, non-malicious AI-generated code will plausibly produce, meaning a substring/regex
approach fails on its own stated purpose as a *usability* gate (ADR-0034: "this precheck is a
usability and contract check") before it even reaches the security question.

**The alternative:** a real tokenizer with full string/template/comment/regex-division correctness,
plus grammar-position awareness for the `import` keyword (statement/expression position vs.
property-name/method-name position), implemented via an **existing, pinned, governed ECMAScript
parser** — not a hand-rolled lexer inside the Go daemon — run out-of-process or otherwise
resource-bounded against fully adversarial input, and checked at AST-node granularity
(`ImportDeclaration`, `ExportNamedDeclaration` with a `source`, `ImportExpression`,
`MetaProperty`). Given the Go daemon's near-zero dependency footprint today, adopting any such
parser is itself a consequential architecture decision under `AGENTS.md`'s own rules and should get
its own ADR rather than being folded quietly into M1/M2 of ADR-0034's implementation plan.

Independent of which parser is chosen, this validator's completeness must not stand in for the real
security property: **the runtime module loader must unconditionally refuse every static or
dynamically constructed module request, with zero filesystem/network/package/import-map/ambient
fallback** — exactly as ADR-0034 and the Threat Model already state, and exactly the thing both
documents currently list as unimplemented and unproven. That runtime-loader evidence, not this
validator, is the actual blocker for any claim that Capsule can safely admit single-file `.mjs`
sources.

## Concrete counterexamples

### 1. Static import syntax

**False refusal** — substring denylist over raw text:

```js
const s = "please use import { foo } from 'bar' in your code";
```

Zero real import syntax; the text is string-literal *data*. Any scanner that greps for
`import ... from` without knowing it's inside a string wrongly refuses this.

**False refusal** — comments containing lookalike text:

```js
// TODO: import { logger } from './log.mjs' once we support it
export function run() {}
```

Legitimate ESM module, zero real imports; a substring scanner over raw bytes refuses it because the
banned text sits inside a `//` comment.

### 2. Dynamic `import()` disguised in strings / indirect via `eval`/`Function`

**False accept:**

```js
eval("import('./evil.mjs')");
new Function("p", "return import(p)")("./evil.mjs");
```

The literal `import(` token never appears outside a string in the source a token-scanner sees. A
scanner limited to "reject `ImportExpression`/`ImportDeclaration` syntax" necessarily misses this;
catching it requires *also* banning `eval` and the `Function` constructor as identifiers/call
targets, a different (CommonJS/dynamic-code) concern layered on top of the import-syntax check.

Note: `globalThis['im'+'port']` does *not* actually work in real engines — `import` is pure syntax,
not a property reachable off `globalThis`, so string-concatenation tricks against `globalThis`
cannot summon it. Worth stating explicitly so the validator design isn't over-built against a
non-threat while under-built against the real one (`eval`/`Function` string payloads).

### 3. `import.meta`

**False accept** — comments between tokens are legal, so an exact-substring check for
`"import.meta"` misses:

```js
const url = import/**/.meta.url;
```

Valid JS; `import`, `.`, `meta` are three separate tokens with a comment between the first two.
`.includes("import.meta")` is `false`, yet this is a live `MetaProperty` access.

**False refusal — the sharpest counterexample in this review** — `import` is only reserved as an
*IdentifierReference*; it remains a perfectly legal `IdentifierName` in *property-name position*:

```js
const obj = { import: { meta: 42 } };
console.log(obj.import.meta); // 42, no ESM import anywhere
```

The token stream literally contains `import` `.` `meta` — the same three consecutive tokens as real
`import.meta` — but here `import` is preceded by a `.` (it's `obj`'s property, not a primary
expression). A rule that fires on "`import` token immediately followed by `.meta`" without checking
whether `import` is itself preceded by a `.` false-refuses this. Distinguishing the two cases is a
*grammar-position* fact, not a lexical one — it needs at least one token of backward context and
knowledge of expression- vs. property-name position, which is parser territory, not pure lexing.

A companion false refusal on the `(` side of the same ambiguity — `import` is also legal as an
object/class method name:

```js
const cache = {
  import(path) {           // method named "import", not a call
    return this._store[path];
  },
};
```

Token-for-token this starts identically to `import(path)` used as a real dynamic-import call. A
rule that fires on "`import` immediately followed by `(`" without knowing it's in
`PropertyDefinition`/`ClassElement` (method-name) position false-refuses this.

### 4. Comments hiding/faking syntax

JS block comments do **not** nest — `/*` inside a comment has no special meaning; the comment ends
at the *first* `*/`. This produces two classic scanner bugs, both direct consequences of the
ECMA-262 lexical grammar for `MultiLineComment`, not edge cases:

**False accept via a greedy regex spanning a later `*/` inside a string:**

```js
/* start
*/ import x from './evil.mjs';
const s = "trailing */ decoy";
```

A scanner that strips comments with a greedy, dot-all regex (rather than non-greedy and
string-aware) matches from the *first* `/*` to the *last* `*/` in the whole file — which here sits
inside the string literal — and swallows the real `import` statement into what it thinks is "all
comment."

**False accept via line-comment detection that ignores string state:**

```js
const s = "config: //"; import x from './evil.mjs';
```

A naive "strip from `//` to end of line" pass that doesn't track string state sees the `//` inside
the string literal and discards the real `import` statement that follows on the same line.

Both bugs share a root cause: comments and strings cannot be recognized by two independent regex
passes over raw text — the scanner must interleave state in a single left-to-right pass.

### 5. String literals: escapes and line continuations

**False refusal from bad string-boundary detection:**

```js
const s = "an escaped quote \" then says import x from 'y'";
```

If the scanner locates a string's end by finding the next `"` without accounting for `\"` as an
escape, it misreads the trailing text as bare source.

**Same failure mode from line continuations:**

```js
const s = "abc\
import x from 'y'";
```

Backslash-newline is a valid `LineContinuation`, producing no line break in the string's value. A
scanner that assumes double-quoted strings never span a raw newline terminates the string early and
misreads the remainder as live code.

### 6. Template literals — the false-accept sharpest edge

**Lexical trap** — non-nesting-aware boundary matching:

```js
const x = `outer ${`inner ${1 + 1} template`} end`;
```

A regex treating backtick-to-backtick as one opaque string terminates the outer template at the
first inner backtick.

**Semantic trap — the important one.** Unlike plain strings, the `${...}` portions of a template
literal are **live, executable expression syntax**, not inert text. Any rule that treats "content
between matching backticks" as opaque string data — the exact rule that's correct for `'...'`/
`"..."` — is a **false accept** for templates:

```js
const x = `${ await import('./evil.mjs') }`;
```

A real `ImportExpression` sits inside `${...}` and is never inspected if templates are treated as
opaque strings. Nesting compounds this, since a `${...}` can contain a full nested template
literal recursively — real JS lexers (V8's, Acorn's) maintain an explicit template-brace stack for
exactly this reason; a `}` sometimes closes a substitution and sometimes closes an ordinary block.

### 7. Regex-literal vs. division-operator ambiguity

The ECMAScript spec delegates disambiguation of a lone `/` to the *syntactic grammar*
(`InputElementRegExp` vs. `InputElementDiv`), not to the lexer alone. A pure tokenizer without
parser-level context can only approximate this with a "was the previous significant token
value-producing" heuristic:

```js
a = b
/hi/g.exec(c).map(d)
```

Because `b` is value-producing and no ASI rule forces a semicolon here, this parses as **division**
(`a = b / hi / g.exec(c).map(d)`), not an assignment followed by a regex-literal statement. If a
scanner's heuristic wrongly decides a span *is* a regex literal, it skips everything up to the next
`/` as opaque regex body — and if real code (potentially containing a genuine `import` token) sits
in that span, it is silently swallowed, producing a false accept.

A concrete false-refusal example for the substring-denylist approach specifically:

```js
// perfectly ordinary lint/codemod helper — zero real ESM imports
const importStatementPattern = /import\s+.*\s+from\s+'.*'/;
```

The regex's source text contains the literal substrings `import` and `from`; a substring/regex
denylist over raw bytes flags this even though it's a data pattern, never an actual import.

### 8. Unicode escapes in identifiers / normalization

An escaped `import` keyword (`import`) is not usable as the actual `import` keyword —
ECMAScript's grammar disallows Unicode escapes inside `ReservedWord` tokens specifically to prevent
this obfuscation (V8/Node reports a syntax error). This claim is stated with moderate confidence —
it is a well-established, spec-driven restriction, but was not re-verified against a running engine
in this review (see Limitations). So the specific "escape the keyword" attack does not work against
real engines — good news for a validator's threat model, but only if the validator itself doesn't
try to defend against it with a broader, riskier mechanism.

The riskier mechanism is exactly what a naive implementer reaches for: **blindly decoding all
`\uXXXX`/`\u{...}` escapes across the whole raw source before running any substring check.** That
conflates identifier-position escapes with escapes that are legitimately part of string-literal
*content*:

```js
const label = "we import nothing here";
```

Fully ordinary, harmless string content. A "decode escapes everywhere, then substring-match"
scanner turns this into the literal text `we import nothing here` and false-refuses it.

Separately, ADR-0034 establishes that Capsule performs **no Unicode normalization** on `main.mjs`
bytes. If the validator's own internal identifier comparisons route through a Unicode-normalizing
string-equality routine (a common library default), it would be inconsistent with the contract's
own byte-identity model.

### 9. CommonJS and Node-specific surfaces

`require`, `module.exports`, `exports`, `__dirname`, `__filename`, `process`, `Buffer`, `global` are
ordinary, syntactically valid free-variable references in ESM grammar — `require('fs')` parses as
an unremarkable function-call expression referencing an undeclared identifier, with no syntax-level
distinction from any other undefined-function call. A syntax scanner cannot reject it as "forbidden
grammar" without being wrong about what grammar is. ADR-0034's own "Closed module-loading policy"
already reflects this: it lists exactly four refusal categories — static `import`,
`export ... from`, dynamic `import()`, and `import.meta` — and says nothing about banning these
identifiers as syntax.

**Concrete false refusal if these names are denylisted as syntax anyway:**

```js
function require(dependencyName) {
  return this._localPolyfills[dependencyName];
}
```

A fully valid, harmless local helper. Banning the *token* `require` breaks it for no security
benefit.

The correct enforcement point for this whole category is the **runtime global-object surface**, not
the pre-approval scanner: per the Threat Model and ADR-0034, the governed `deno_core` profile is
supposed to expose no `require`, no CommonJS shim, no `process`/`Buffer`/`global`, and no
filesystem/network/package/import-map/fallback loader at all. If true, `require('fs')` in
`main.mjs` simply throws `ReferenceError: require is not defined` at execution — a runtime failure,
not a pre-flight syntax violation. `node:`-scheme specifiers are moot for the same reason: they are
refused by being an `ImportDeclaration`/`ImportExpression` regardless of specifier text.

## Smallest defensible validation design

1. **A single left-to-right, stateful tokenizer** correctly implementing ECMA-262's lexical grammar
   interaction with syntactic context: string literals (full escape table, line continuations),
   template literals with a nesting stack for `${...}` (substitutions must be re-entered into the
   tokenizer as live code), comments (non-nesting, string/template-aware), and regex-vs-division
   disambiguation. This is materially more than "a lexer with some regexes" — it's the scanner
   front-end of a real parser, and the class of component with the highest historical defect rate
   for hand-rolled implementations.
2. **Grammar-position awareness on top of tokens**, not just token adjacency, specifically for the
   `import` keyword: property-name/method-name position (inert) vs. expression/statement position
   (live — `ImportDeclaration`, `ImportExpression`, `MetaProperty`). This means checking real AST
   node types, i.e. a full parse, not a partial one.
3. **Do not hand-roll this inside the Go daemon.** ADR-0028 treats even the *execution-time* JS
   engine with extreme rigor — pinned upstream commits, a hard cap on hand-authored reviewable
   files, reproducible builds, SBOM/provenance — precisely because parser/runtime code is uniquely
   dangerous. A bespoke ECMAScript lexer parsing fully adversarial AI-generated bytes (up to
   262,144 bytes per ADR-0034's own cap) sits in the same risk category.
4. **Dependency options**, given the Go daemon's near-zero dependency graph today (`go.mod` has
   exactly `golang.org/x/sys`):
   - A maintained, pinned, actively-audited ECMAScript parser used strictly in "parse and classify"
     mode (never executed, never used to transform `B` — the byte-exact pass-through invariant
     stays intact). Needs the same fork-pin-govern treatment ADR-0028 already established for
     `deno_core`/`rusty_v8`, not an ad hoc `go get`.
   - Reusing the same governed V8 already committed to for execution (via a narrow, out-of-process,
     resource-bounded "parse and report node kinds" call) removes parser-differential risk between
     what the validator approved and what the runtime actually runs.
   - Reusing the already-pinned `typescript` devDependency via a narrow, bounded CLI is possible but
     is a bigger architectural change than it sounds — it means the Go control plane depends on
     spawning a Node subprocess per proposal, and needs its own ADR under `AGENTS.md`'s
     dependency/language-addition rule.
   - Whichever parser is chosen, run it **sandboxed and resource-bounded** — the validator now
     parses fully hostile untrusted bytes, so parser bugs (stack overflow, memory exhaustion,
     ReDoS-equivalent inputs) are a realistic DoS surface if linked in-process without isolation.
5. **Treat this validator's verdict as a usability/early-refusal gate only, never the security
   boundary** — exactly what ADR-0034 already says ("This precheck is a usability and contract
   check, not the runtime security boundary") and what the Threat Model reiterates. The real
   security property needed is the governed `deno_core` module loader unconditionally refusing every
   static or dynamically constructed module request, with no fallback wired in at all — a property
   both documents currently list as unimplemented and unproven.

## Confidence and limitations

Moderate-to-high confidence on the JavaScript-semantics claims (grounded in the stable, well-
established ECMA-262 specification), moderate confidence on the repo-specific architectural
conclusions (a synthesis of this repo's own stated principles, not external opinion).

- No runtime access: none of the snippets above were executed against Node/V8/`deno_core` in this
  review; they rely on knowledge of the ECMA-262 specification and well-known JS lexer literature,
  not empirical verification.
- No validator code exists yet to review directly — this is a design-level assessment against
  ADR-0034's stated intent, not a review of an implementation with real bugs.
- Docs are explicitly marked pre-freeze/aspirational in places (`THREAT_MODEL.md`: "the current
  scaffold does not implement or satisfy these properties"); some details (the governed `deno_core`
  fork's actual patch set / exact runtime global-object surface) are asserted intent, not
  independently verified fact in this review.
- `experiments/` and any git-ignored AI-Central history were not exhaustively searched for prior art
  on this exact question.
- ASI (automatic semicolon insertion) interaction with the regex/division ambiguity was not
  exhaustively enumerated; only a citation-grade example is included.
