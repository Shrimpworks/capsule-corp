# Owner-only internal-alpha fixed-store passive observation checker

Work item: ADR-0040 fixed-store threshold observation/admission checker
Status: `PASSED`
Scope: passive, owner-held, read-only, re-evaluated checker over the local fixed-store v2 oracle
Evidence or reason: focused exact-boundary, cap-plus-one, combined-precedence, corruption, segment-bearing, orphan, restart, no-rewrite, unchanged-approval, and race tests plus every `AGENTS.md` repository gate pass
Remaining work: no persistent installation trip latch exists; product p95 source/window/lifetime, authenticated `RequestAttempt` wiring, and full five-threshold product enforcement remain `BLOCKED`
Next action: a future separately authorized product-consumer slice must resolve trip/reset semantics, select p95 instrumentation, and place an admitted guard immediately before attempt creation
Parent status: owner-only hostile-`.mjs` internal alpha remains `IN_PROGRESS — TRENDING_GOOD`; product admission remains `BLOCKED`

## Defensive and authorized scope

This slice defensively validates passive comparisons for ADR-0040's fixed-store thresholds using only
`internal/execution/registrationstate` fixtures, fake authority records, and owned temporary roots.
It accesses no runtime, backend, guest, product endpoint, service, Apple identity, credential,
installed state, restore path, or user data.

The fixed store remains labeled `owner-only-disposable-alpha-oracle`. This slice does not select it
as a production engine and makes no continuity, rollback, restore, secure-deletion, APFS, power-
loss, multi-process, or external-alpha claim.

## Decision and insertion point

No new ADR is needed. Accepted ADR-0040 already selects the exact exception and thresholds, and it
explicitly refines ADR-0031 only for the owner-only disposable internal alpha. The implementation
adds no Supervisor responsibility: it evaluates already Supervisor-owned, fully verified store
facts and returns no permit, mutation, or authority-bearing value.

The exact future insertion point is immediately before the attempt-creation transaction, after the
installation owner is confirmed held and after the existing fixed-store v2 full verifier has
reconstructed the active snapshot and every referenced segment. The passive API is not wired to
`RequestAttempt`, IPC, a daemon, Broker, adapter, runtime, backend, or guest in this slice.

Startup uses a separate owner-held full verification to issue a sealed duration observation. A
later pre-attempt check requires that observation, repeats full verification, requires the same
installation/Supervisor/epoch and owner session, and checks a sealed durable-commit p95
observation. Missing timing evidence, owner loss, a mismatched session, repair/quarantine/transition
state, disabled attempts, a known unreferenced orphan, unknown extra data, missing data, or corrupt
data refuses as `unknown-state` without rewrite.

These timing observations are re-evaluated inputs, not durable installation state. The checker
persists no first-trip latch: startup verification may be measured again in a later owner session,
and a later caller may supply a newly measured p95 observation. ADR-0040 does not yet select the
trip/reset semantics required for a product consumer, so this slice makes no sticky-stop claim.

## Exact threshold semantics

ADR-0040 selects stops at the stated attempt, active-byte, and segment capacities, while it names
timing observations only when they exceed their stated durations. The checker therefore compares
count/size/segment observations with `>=` and timing observations with `>`. It preserves observed
values and never clamps them:

| Dimension | Last admitted observation | First refusal |
| --- | ---: | ---: |
| cumulative retained attempts, hot plus archived | 127 | 128 |
| active encoded store bytes | 8 MiB minus one byte | 8 MiB |
| referenced archive segments | 15 | 16 |
| startup full-verification duration | exactly 2 seconds | greater than 2 seconds |
| durable-commit p95 | exactly 250 milliseconds | greater than 250 milliseconds |

The first matching refusal is reported in ADR-0040 order: attempts, active bytes, segments,
startup verification, then durable-commit p95. Integrity or unknown-state refusal precedes every
numeric threshold.

This slice compares the p95 threshold with a sealed caller-supplied observed value but deliberately
does not select or authenticate the future product instrumentation's source, sample window,
aggregation lifetime, or durable storage. Those mechanics belong to the later product consumer and
quantitative campaign; inventing them here would exceed the passive fixed-store scope.

## Retained tests

`internal_alpha_fixed_store_policy_test.go` retains:

- every exact boundary and cap-plus-one outcome;
- exact two-second and 250-millisecond acceptance plus one-nanosecond refusal;
- combined unknown state taking precedence over every over-threshold numeric observation;
- missing timing, corrupt active bytes, and known-unreferenced-orphan refusal;
- a referenced-segment world proving active encoded bytes equal only the primary v2 file length and
  exclude segment bytes;
- two startup/reopen cycles over identical bytes;
- unchanged active/archive inventory across every live policy check; and
- unchanged retained approval, attempt count, and `AttemptsDisabled` authority state after read-only
  evaluation.

## Limitations and deferred work

- No product attempt call consumes this oracle yet; wiring requires the future authenticated
  owner-only alpha consumer and its own integration/fault corpus.
- No persistent installation trip latch is implemented. Startup and p95 observations can be
  remeasured; the checker evaluates the observation supplied to the current call.
- Durable-commit source instrumentation, sampling window, aggregation lifetime, and persistence are
  not implemented or selected here.
- The startup duration is local monotonic elapsed-time evidence, not a service-level guarantee or a
  durable historical trip record.
- The fixed store remains finite and disposable. Threshold or integrity refusal may permanently
  stop the installation under ADR-0040.
- F6 production-engine selection, installed protected state, real power-loss evidence, restore,
  continuity, automatic update, external alpha, full five-threshold product enforcement, and
  product admission remain `BLOCKED`.

## Verification

Focused verification completed:

```sh
GOCACHE=/private/tmp/capsule-go-cache go test ./internal/execution/registrationstate -run InternalAlpha -count=1 -v
```

Full repository verification completed:

```sh
fnm exec --using=22.22.1 -- pnpm install --frozen-lockfile
fnm exec --using=22.22.1 -- pnpm check
fnm exec --using=22.22.1 -- pnpm lint
fnm exec --using=22.22.1 -- pnpm test
fnm exec --using=22.22.1 -- pnpm verify:schemas
fnm exec --using=22.22.1 -- pnpm verify:adrs
fnm exec --using=22.22.1 -- pnpm format:check
GOCACHE=/private/tmp/capsule-go-cache go test -race ./internal/execution/registrationstate -run InternalAlpha -count=1
GOCACHE=/private/tmp/capsule-go-cache go test ./internal/execution/registrationstate -count=1
GOCACHE=/private/tmp/capsule-go-cache go test ./...
GOCACHE=/private/tmp/capsule-go-cache go vet ./...
GOCACHE=/private/tmp/capsule-go-cache go build ./...
GOCACHE=/private/tmp/capsule-go-cache-a11a-policy GOLANGCI_LINT_CACHE=/private/tmp/capsule-golangci-cache-a11a-policy golangci-lint run ./...
GOCACHE=/private/tmp/capsule-go-cache-a11a-policy GOMODCACHE=/private/tmp/capsule-go-mod-cache-a11a-policy go run golang.org/x/vuln/cmd/govulncheck@latest ./...
git diff --check
```

`golangci-lint` reported `0 issues`. `govulncheck` found zero reachable vulnerabilities.
