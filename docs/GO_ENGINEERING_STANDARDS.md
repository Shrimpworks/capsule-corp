# Go engineering standards

This is the code-quality bar for Go in this repository: naming, structure, error handling,
testing, and lint hygiene. It is a companion to `AGENTS.md`, not a replacement — `AGENTS.md`
governs security-sensitive behavior and wins on any conflict. This document governs everything
`AGENTS.md` does not: the ordinary craft of writing and reviewing Go here.

It assumes and extends, rather than restates, [Effective Go](https://go.dev/doc/effective_go),
the [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments) wiki, and the general
[Google Go Style Guide](https://google.github.io/styleguide/go/). Where this document is silent,
those apply. Where it is specific, the specifics here win for this repository.

## Naming

- Package names are short, lowercase, no underscores, and do not stutter with their contents
  (`registrationstate`, not `registration_state` or `state.RegistrationState`). This repository's
  packages already follow this; keep it.
- Exported identifiers get a doc comment starting with the identifier's name (`// Server is ...`),
  so `godoc` and IDE hovers read correctly. Every exported type, func, const, and var needs one.
- Prefer precise domain names over abbreviation (`Classification`, `AttemptID`, `PlanRegistration`)
  over generic ones (`Type`, `ID`, `Data`). This costs a few extra characters and buys a
  self-documenting diff — worth it in a codebase whose whole point is precise authority boundaries.
- A getter is not `GetX()`; it is `X()`. A single-method interface is usually named for what it
  does, suffixed `-er` (`Verifier`, `Selector`), matching existing packages.

## Function and file shape (clean code)

- Keep functions small and at one level of abstraction. If a function mixes "what" (business
  decision) and "how" (byte-level mechanics), extract the mechanics into a named helper even if it
  has one call site — the name documents intent.
- Guard clauses over nested conditionals. Return/continue early; do not wrap the happy path in an
  `if` that could instead be a precondition check that exits.
- Don't introduce an interface until there are two real concrete implementations (or one concrete
  implementation plus a test fake, which counts). "Accept interfaces, return structs" — see
  `internal/api.NewServer`, which returns `http.Handler` deliberately narrow but is constructed
  concretely.
- Rule of three: tolerate one duplicated block, extract on the third occurrence, not the second.
  Premature extraction produces the wrong abstraction, which is more expensive to unwind later than
  duplication is to tolerate now.
- A file that is hard to summarize in one sentence is a candidate to split. This repository already
  splits by concern within a package (`digests.go`, `errors.go`, `types.go`) — follow that pattern
  rather than growing one file per package.

## Errors

- Wrap with `%w` and add call-site context (`fmt.Errorf("open registration store: %w", err)`).
  Never discard an error silently. If an error truly cannot matter (e.g., a best-effort cleanup),
  say so in a comment at the discard site — don't just drop it.
- Use typed/sentinel errors plus `errors.Is`/`errors.As` for anything a caller needs to branch on.
  The `Classification` + `contractError` pattern in `internal/execution/lifecyclestate/errors.go`
  is the model: a closed, named vocabulary instead of string-matching error text.
- Fail closed. An unrecognized classification, capability, or state transition is an error, not a
  default. This is also an explicit project principle in `docs/PROJECT.md` — it is a code-quality
  rule here because a permissive `default:` case in a switch is where that principle quietly dies.
- Reserve `panic` for programmer errors (invariant violations, impossible states) that should never
  cross an API boundary — never for expected control flow, expected malformed input, or anything a
  caller could reasonably trigger.
- `writeJSON`-style helpers that intentionally swallow a late write error (headers already sent,
  nothing left to do about it) must say so in a comment at the point of discard, not just `return`.
  Silence and "I checked and there's genuinely nothing to do" look identical in a diff; only one of
  them is reviewable.

## Concurrency

- Thread `context.Context` through anything that does I/O or can block; it is the first parameter,
  never stored on a struct.
- No goroutine may outlive the request/operation that started it without an explicit owner and
  shutdown path. If a goroutine is intentionally long-lived, document who stops it and how.
- Run `go test -race` before merge on anything touching shared state, stores, or the lifecycle
  packages. Add it to local habit even though it is not yet wired into `make check` (see
  [Lint and CI wiring](#lint-and-ci-wiring) below).
- Document lock ordering and invariants at the type, not the call site, when a type has one
  (`fixed_store_v2.go`-style code is exactly where an undocumented ordering bug hides).

## Testing

- Table-driven tests are the default shape for anything with more than two input cases. Name each
  case; a failure should be findable by name alone, not by counting table rows.
- Byte-exact protocol and store code (`v0candidate`, `v0cbor`, `registrationstate`) uses
  conformance/known-answer fixtures, not hand-rolled equality assertions. This repository already
  does this at scale (hundreds of fixture cases) — keep growing the corpus rather than reverting to
  ad hoc assertions when a new encode/decode path is added.
- Use `t.Helper()` in every test helper that can fail, so failures point at the caller's line.
- No `time.Sleep`-based synchronization in tests. Use channels, explicit fakes, or condition
  polling with a bounded timeout. A flaky test is worse than no test — it trains reviewers to
  ignore CI failures.
- One behavior per test function. Arrange/Act/Assert order, blank-line separated, is the default
  layout; don't interleave setup and assertions.
- A behavior change needs a test change in the same commit. A security-boundary change additionally
  needs a positive and a negative case — this is already required by `AGENTS.md` and
  `docs/DEVELOPMENT.md`; it is restated here because it is also simply correct test design.

## Documentation

- Every package gets a package doc comment — either at the top of the most central file or in a
  dedicated `doc.go` when the package needs more than a paragraph (see
  `internal/execution/doc.go`). Keep following this convention; it is what lets a new contributor
  orient without reading every file.
- Comments explain *why*, not *what*. `// increment the counter` above `count++` is noise; `// two
  observers may race here, so this is the linearization point` is not.
- When a comment describes project status, use the exact vocabulary in `docs/STATUS_LANGUAGE.md`
  (`PASSED`, `IN_PROGRESS — TRENDING_GOOD/BAD`, `BLOCKED`, `NO_GO`) rather than informal words like
  "done" or "working". This keeps code comments and project docs from silently disagreeing.
- A lint suppression (`//nolint`, a `.golangci.yml` exclusion) must carry an inline rationale, the
  same way the existing gosec exclusions in `.golangci.yml` do. An unexplained suppression is a
  future landmine, not a saved argument.

## Lint and CI wiring

- `gofmt` is mandatory and covers `cmd` and `internal` (`make fmt` / `make check`).
- `.golangci.yml` enables the standard linter set, `gosec`, and the `revive` exported/package
  documentation rules. CI runs `golangci-lint run --disable revive ./...` on every event so
  security and correctness findings always gate across the whole repository. A separate
  `--enable-only revive` gate uses `only-new-issues` on PR/push events; scheduled and manual runs
  have no change baseline and skip only that documentation gate. Existing documentation debt
  remains tracked in issue #217 and is fixed in package-level batches.
- `.github/workflows/ci.yml` pins `golangci-lint` to `v2.12.2` via `golangci-lint-action`; match
  that version locally. To check a branch locally, run the full security/correctness command above
  and `golangci-lint run --enable-only revive --new-from-rev=origin/main ./...` against an
  up-to-date base. `make check` (and therefore `make ci`) still runs the unrestricted
  `golangci-lint run ./...` documentation audit and can fail on the existing #217 backlog even
  when CI passes; it also runs `govulncheck ./...` alongside `gofmt`/`go vet`.
- If `golangci-lint run ./...` reports issues whose file paths don't match anything in your working
  tree, or that vanish after touching unrelated files, suspect a stale cache from a different git
  worktree sharing your `~/.cache/golangci-lint` (or platform equivalent) — run
  `golangci-lint cache clean` before trusting the result. This is a caching quirk of running
  multiple worktrees of the same module, not a lint config problem.
- `govulncheck` needs a newer Go toolchain than `go.mod` currently pins for the project itself;
  run it with `GOTOOLCHAIN=auto` (as `make check` and CI both do) so `go run` can fetch a matching
  toolchain for the tool alone, instead of bumping the project's own pin.

## Security-adjacent Go hygiene

These are code-quality rules, not new security policy — `docs/security/THREAT_MODEL.md` and
`AGENTS.md` remain the authority on the latter. See
`docs/security/SECURE_CODING_STANDARDS.md` for the cross-language version of these principles
shared with the JavaScript/TypeScript and Rust codebases.

- Never build a filesystem path from untrusted input without validation (gosec `G304`). Where a
  path is provably repo-controlled (test fixtures), say so in the surrounding comment, matching the
  existing `.golangci.yml` per-path exclusions — don't blanket-disable the rule.
- Integer conversions that could truncate or overflow (gosec `G115`) need either a bounds check or
  a comment proving the invariant that makes the check unnecessary, the way
  `internal/protocol/v0candidate/cbor.go` documents its offset invariant against the byte-exact
  conformance fixtures. "gosec can't see it" is not sufficient justification by itself; the comment
  is what makes the suppression reviewable.
- Check or explicitly discard-with-reason every `Close()`/`Sync()` error on anything durable
  (files, stores). A silently dropped `Close` error on a write path is a silent data-loss bug.
- Never log secrets, private keys, approval content, or user job content. If a struct might
  contain any of those, do not give it a default `%v`/`%+v`-friendly shape that a future `log.Printf`
  could accidentally dump — implement `String()`/`GoString()` to redact, or keep such fields
  unexported and undumped.
- `os/exec` arguments are never built by string concatenation or shell interpolation from anything
  outside the trusted policy layer. This is the same rule as `AGENTS.md`'s "no ambient authority"
  principle applied at the syscall boundary.
- No hand-rolled cryptographic primitives. Use the standard library or an already-adopted,
  reviewed dependency per `docs/ECOSYSTEM_REUSE_AND_ADOPTION.md`.

## Dependencies

Before adding a Go dependency, follow `AGENTS.md`'s requirement to consult
`docs/ECOSYSTEM_REUSE_AND_ADOPTION.md` and complete its checklist. From a code-quality angle,
also check: does it pull in `cgo`, does it have its own large transitive graph, and does it fit
this repository's near-zero-dependency posture (`go.mod` currently declares exactly two direct
dependencies)? A dependency that fails any of these is a harder sell here than in an average Go
service.

## Precedence

If this document and `AGENTS.md` conflict, `AGENTS.md` wins. Raise the conflict as an issue or ADR
rather than silently picking one.
