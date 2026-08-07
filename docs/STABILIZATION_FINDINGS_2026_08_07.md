# Stabilization findings — 2026-08-07

Work item: post-checkpoint research, alignment, and targeted bug hunt

Status: `PASSED`

Parent owner-only hostile-`.mjs` internal alpha: `IN_PROGRESS — TRENDING_GOOD`

Product admission and installed security boundary: `BLOCKED`

## Scope

This stabilization pass began from fresh `main` commit
`d1972b04be1de38c07f0fcc471d33e796e61d507`. It was defensive and local-only. It reviewed the
passive authority, approval, registration, completion, and archive packages; accepted and aligned
the ADR-0029 decision; completed native XPC and typed guest-transport research; and ran repository
verification. It created no listener, service, credential, signature, protected root, runtime,
adapter, VMM, or guest and changed no product authority.

The review used correctness, readability, architecture, security, and performance as separate
axes. A finding was retained only when a local command reproduced it or exact source established
it. No speculative code change was made.

## Validated finding and fix

### S1 — test-only verifier guard escaped the current checkout

Classification: correctness and test isolation

Status: `PASSED`

`TestFixtureVerifierHasNoNonTestGoConsumer` recursively scanned hidden local tooling state beneath
the repository. With locally retained `.claude/worktrees`, it parsed those worktrees' own
`fixture_verifier.go` definitions and reported them as product consumers in the current checkout.
The test therefore depended on unrelated local worktree inventory and failed before evaluating the
intended security invariant.

The guard now skips the directory-name classes Go package discovery excludes: dot-prefixed and
underscore-prefixed directories, `testdata`, `vendor`, and `node_modules`. A table test freezes the
tool/worktree exclusions and preserves ordinary `internal` package traversal. The guard still
parses every current-module non-test Go file it reaches and still refuses any use of
`FixtureVerifier` or `NewFixtureVerifier` outside the definition and tests. `.golangci.yml` now
also excludes `.claude` and `.codex` local state with an explicit rationale, so lint cannot analyze
complete nested worktrees as though they belonged to the current checkout.

Reproduction before the fix:

```text
go test -race ./internal/execution/approvalattempt
--- FAIL: TestFixtureVerifierHasNoNonTestGoConsumer
test-only FixtureVerifier is referenced by non-test Go file
.claude/worktrees/.../internal/execution/approvalattempt/fixture_verifier.go
```

Verification after the fix:

- `go test -race ./internal/execution/approvalattempt` — `PASSED`
- 20 consecutive ordinary runs of `authorityplane`, `approvalattempt`, `completioncomposer`, and
  `archivestate` — `PASSED`
- focused `go vet` over those four packages — `PASSED`

## Reviewed areas with no validated defect

- ADR-0029/S3 topology and refusal ownership: the accepted decision, three existing passive
  methods, two later approval/attempt methods, and ADR-0044's separate CLI method now have one
  consistent dependency order and claim boundary.
- Passive native-XPC contract: caps, deadlines, copy-before-parse ordering, correlation-only
  request IDs, and response-loss ownership remain separated from platform enforcement.
- Approval and attempt verification: exact public-key authorization, identifier domains, replay,
  defensive byte ownership, and passive fixture/product-verifier separation showed no additional
  validated defect.
- Completion composer/store: strict JSON, immutable first completion, response-loss convergence,
  atomic creation/publication, terminal lifecycle joins, defensive copies, and fixed public summary
  showed no additional validated defect.
- Archive-state decomposition: the six-stage `validateCandidateCohort` refactor preserves the
  original refusal order, forward/reverse approval-attempt links, collision sets, lifecycle
  disposition checks, exact-record digests, and final encoded-byte underflow check. Its focused
  mutation suite and repeated package runs pass.

This means only that the inspected passive/local mechanics and tests supplied no second reproduced
finding. It is not an independent security review, installed-state result, real transport result,
or admission evidence.

## Test-scalability limitation

Status: `BLOCKED`

The combined race-detector command reached Go's default ten-minute package timeout in
`TestCoherentBackupAcceptsExact64SegmentWorldAndPreservesCapPlusOneRefusal`. It reported no data
race before timeout. The same exact 64-segment/cap-plus-one test passes without race instrumentation
in about 52 seconds on this host, and the required ordinary suite remains the release gate.

This is retained as a race-suite scalability limitation, not hidden by increasing a timeout and not
misreported as a product race. Follow-up should profile the repeated growing-state decode/activation
path and choose deliberately between a faster semantically equivalent boundary fixture, a split
race corpus, or an explicit longer race job. That work must preserve the exact 64-segment and
cap-plus-one coverage.

## Exported-contract documentation batch

Status: `PASSED`

Full repository lint exposed the already-tracked issue #217 documentation backlog. This branch
completed one authority-bearing package batch for `internal/execution/lifecyclestate`: the closed
enums, digest domains, identifiers, immutable bindings, passive permits/results, record versions,
constructors, and defensive projections now state their purpose and caller-visible authority
limits. No identifier, signature, field, return type, validation branch, digest input, or persisted
byte changed.

The remaining legacy packages stay tracked under issue #217. CI deliberately gates only issues
introduced by the current diff until those package-at-a-time batches finish. This branch's
CI-equivalent new-issue lint passes with zero findings; the exact full-repository lint command was
also run and still fails on that known backlog. It is not represented as fully remediated.

## Research and alignment outcomes

- CL1 `PASSED`: ADR-0029 is Accepted with one unprivileged per-user native-fronted Go Supervisor,
  two role-specific services/four ordinary methods, authentication before body copy, method-specific
  bridges, Go-only durable authority, correlation-only request IDs, and `AttemptID`-only recovery.
- R1 `PASSED`: the controlled harness baseline is low-level `xpc_connection_t`, peer requirement
  before activation, exact-message `SecCode` validation, connection-time EUID/ASID checks,
  protocol-owned deadlines, and store-owned replay. C2 execution still requires exact separate
  authorization.
- R2 `PASSED`: C5a now has a passive three-stream state machine, C2A's narrowed 262,144-byte
  first-slice payload caps, launcher-only completion authority, cap-plus-one continuous drain,
  cancellation/reset/response-loss taxonomy, durable terminal-proof join, and restoration matrix.
  C5a implementation and C5b execution remain incomplete.

## Verification

The final branch verification is recorded in the pull request. At this checkpoint:

- pinned Node 22.22.1/pnpm 10.28.2 install, build/check, Biome lint, tests, schemas, conformance,
  field-authority, and ADR-index gates — `PASSED`;
- `go test ./...`, `go vet ./...`, and `go build ./...` — `PASSED`;
- `go run golang.org/x/vuln/cmd/govulncheck@latest ./...` — `PASSED`, with zero called
  vulnerabilities;
- `golangci-lint run --new-from-rev=d1972b04be1de38c07f0fcc471d33e796e61d507 ./...` —
  `PASSED`, zero issues;
- `golangci-lint run ./...` — `BLOCKED` on the known issue #217 legacy exported-comment backlog;
  the completed `lifecyclestate` batch no longer appears;
- the combined race-detector package run — `BLOCKED` only by the ten-minute 64-segment test timeout
  described above; targeted `authorityplane` and `approvalattempt` race runs passed, and no race was
  reported before the long run timed out.

## Next sprint starting line

The clean dependency-ordered implementation choices are:

1. C4: freeze passive `SubmitApprovalV0` and `RequestAttemptV0` without activating a listener or
   signer.
2. C5a: freeze the typed transport contract from the R2 state machine without running libkrun or a
   guest.
3. Continue issue #217 one package at a time, with `registrationstate` as the next authority-bearing
   documentation batch.
4. C2, C1, or C3 only after the owner supplies their already-documented exact external workspace,
   host, fixture, or Apple identity authorization.

Installed composition, Broker signing, real transport execution, and product admission remain
later gates; this stabilization pass does not pull them forward.
