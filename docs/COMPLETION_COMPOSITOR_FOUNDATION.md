# Unwired FakeBackend completion compositor foundation

## Status and scope

Work item: fixed completion/transcript/public-summary compositor foundation

Status: `PASSED`

Parent internal-alpha workstream: `IN_PROGRESS — TRENDING_GOOD`

Product admission: `BLOCKED`

This slice is defensive, local-only, and unwired. Package
`internal/execution/completioncomposer` reads three narrow interfaces representing
Supervisor/store truth: a committed created attempt, its lifecycle record, and a committed-last
completion record. Tests use deterministic repository fixtures and `FakeBackend` identity only.
The package has no writer, endpoint, Broker UI, signing key, receipt, adapter call, runtime,
backend activation, VM, or guest path.

## Terminal proof rule

The compositor emits a terminal projection only when all of these retained facts agree:

1. the attempt, consumed approval, retained registration, exact plan bytes, and role bindings
   cross-link and revalidate;
2. a bounded typed-JSON completion record explicitly says `fake-no-guest-json`, `succeeded`, and
   `committed-last` for the same attempt, registration, plan, and immutable binding digest;
3. the lifecycle is the exact completed six-operation `FakeBackend` sequence with confirmed
   destroy, no first failure, no cleanup requirement, and authoritative absence; and
4. the result fits both the fixed 262,144-byte JSON ceiling and the approved plan's output cap.

A missing completion record or incomplete lifecycle returns `NOT_TERMINAL`. Corrupt retained
truth returns `RECOVERY_REQUIRED`; contradictory identity or lifecycle facts return `BINDING`.
The compositor never manufactures a completion fact. A future Supervisor-owned store and producer
must commit that fact last before wiring this reader.

EOF, process or daemon exit zero, daemon/backend prose, guest diagnostics, paths, artifact names,
and timing do not appear in the completion input, transcript, or public summary and cannot
establish success. `destroyed` proves cleanup only until the separate completion-last fact is also
present and consistent.

## Fixed bounded objects

- `CompletionRecord` owns and validates exact typed JSON: maximum 262,144 bytes, integer-only
  numbers within the JSON safe-integer range, maximum depth 32, fixed node/member/element/key
  bounds, duplicate-key refusal, strict UTF-8/surrogate handling, and trailing-data refusal.
- `TerminalProjection` owns the result bytes and includes only retained typed identities, digests,
  completion disposition, lifecycle sequence/state, teardown, cleanup, and absence truth.
- The deterministic unsigned transcript is at most 4,096 bytes. It carries only fixed fields and
  a result digest/length, not result content.
- The deterministic public summary is at most 256 bytes and contains only `state`, a typed attempt
  identifier, and a transcript digest identifier.

The transcript labels itself `fake-no-guest-local-mechanic` and states that it is neither a signed
attestation nor product-success evidence. The stale receipt schema is not consumed or extended.

## Retained evidence

`internal/execution/completioncomposer/compositor_test.go` covers the frozen known answer,
response-loss restore/reopen and replay convergence, exact cap and max-plus-one, malformed and
untyped JSON, missing/corrupt/contradictory facts, plan-cap refusal, and defensive byte ownership.
`internal/execution/completioncomposer/testdata/known-answer.json` freezes exact result,
transcript, transcript SHA-256, and public-summary bytes.

This uses the standard library only. It follows the fixed bounded envelope `BUILD-NARROWLY`
recommendation in `ECOSYSTEM_REUSE_AND_ADOPTION.md`; it adds no dependency.

## Limitations and next action

The scoped foundation is complete, but the parent remains `IN_PROGRESS — TRENDING_GOOD` and
product admission remains `BLOCKED`. There is no production durable completion object or writer,
guest completion port, integrity/result producer, real process-tree absence proof, authenticated
consumer, signature, receipt, or installed-profile evidence. FakeBackend creates no guest, so this
slice cannot validate guest completion or teardown.

Next work must separately freeze and implement the Supervisor-owned durable completion producer
and store transaction, then connect the compositor only after the completion-last, transport,
integrity, teardown, and authoritative-absence controls have their own retained evidence and
authorization.
