# Public-key approval to FakeBackend integration

## Result

Work item: backend-independent ApprovalGrant submission through durable attempt creation and the
existing `AttemptID`-only no-guest lifecycle.

Status: `PASSED`

Scope: one repository-local production-shaped public-key verifier, one fixed v0 Supervisor store,
exact retained execution-plan fixtures, one public-only signed ApprovalGrant known answer, the
existing v0-to-v1 migration, and `FakeBackend.CreatesGuest() == false`.

Evidence or reason: the production-shaped verifier now satisfies the method-specific candidate
submission port without accepting candidate-owned nonce, protected key ID, or authorization
identity from the caller. The verifier derives those values from the signed envelope and its
trusted authorization set. The durable Supervisor component independently applies exact
registration, plan, installation, epoch, Supervisor, effective-time, registration-expiry, replay,
and nonce-uniqueness rules before committing a new approval. Existing atomic consume/create and
`AttemptID`-only lifecycle seams are unchanged.

Focused integration tests prove:

- a valid signature for registered plan A cannot authorize separately registered plan B;
- complementary mathematically valid signatures over one canonical payload converge on one
  approval record;
- response loss after approval commit returns the retained approval on replay;
- response loss after atomic consumption and attempt creation returns the same committed
  `AttemptID` after reopen;
- the committed attempt has caused no fake effect before explicit `Drive(AttemptID)`;
- a durable lifecycle record exists while fake call counts are still zero, closing the interval
  between attempt commit and the first fake effect;
- interruption after the fake prepare effect recovers from the same durable intent and effect ID;
- every fake lifecycle operation applies exactly once through destroyed plus fake authoritative
  absence; and
- terminal `AttemptID` replay does not redrive an effect.

Remaining work: authenticated product IPC, installed Broker rendering and live signing, Keychain
and Secure Enclave authorization, protected installed Supervisor state, the selected product store,
real adapter/backend identity and reconciliation, runtime/profile admission, a guest, durable
product completion composition, evidence signing, and receipts.

Next action: connect the same frozen authority ordering only through the separately selected
method-specific authenticated product boundary. Do not pass plan, source, binding, backend, image,
mount, path, resource, or guest bytes at execute time.

Parent status: owner-only hostile-`.mjs` internal alpha remains
`IN_PROGRESS — TRENDING_GOOD`; product admission remains `BLOCKED`.

## Adjacent dependency ordering

The passive ADR-0040 fixed-store stop-policy slice at commit
`54c56c38cef2e33062c8d9f1ff05f312f0a69025` and draft PR #224 is dependency awareness only; no
file or behavior from that branch is integrated here. It is a passive, re-evaluated
observation/admission checker, not complete operational enforcement: it has no persistent trip
latch, selected durable-commit-p95 instrumentation/window/lifetime, `RequestAttempt` wiring, or
response-loss transaction evidence. Full enforcement of ADR-0040's five thresholds remains
`BLOCKED`. Its later guard must apply only to a new attempt, after installation-owner assertion and
fresh full-v2 verification and immediately before the approval-consume/attempt-create mutation. A
replay that finds an already consumed approval must return its existing `AttemptID` before that
guard so response-loss and restart convergence cannot be converted into a new-policy refusal.

There is no conflict with the verifier or `AttemptID`-only lifecycle interfaces. The current
generic `StateStore.commitAttempt` callback, however, exposes only mutable hot authority state; it
does not itself expose a fresh fully verified v2 world, owner assertion, segment/byte inventory, or
timing observations. Safe later wiring therefore needs a v2 transaction-local precondition at the
named insertion point. It must not be a caller-carried permit or an out-of-transaction check that
can race the attempt commit. This slice does not invent that interface.

The passive native XPC v0 prerequisite at commit
`f929dfbc96a9dd4b6e303820a4e97848f6da8497` and draft PR #225 is also dependency awareness only.
It freezes submission, registration, and registered-plan lookup encodings but defines no approval,
attempt, execute, or lifecycle method. A later method-specific adapter can therefore preserve this
slice's ordering: the Broker signs Supervisor-fetched registered data; approval submission carries
the registration identifier and exact signed envelope; attempt creation consumes stored authority;
and lifecycle execution accepts only the Supervisor-issued `AttemptID`. Replacement plan, source,
binding, backend, image, mount, path, resource, or guest bytes remain inadmissible at execute time.

## Frozen submission seam

`approvalattempt.CandidateVerifier.VerifyCandidate` performs bounded exact-envelope framing,
canonical protected-header and payload decoding, public-key signature verification, closed Approval
key-policy validation, signer-to-payload installation/epoch/purpose/audience binding, and grant
interval containment within the authorized key interval. It returns copied exact bytes plus the
signed nonce and verifier-resolved key identity.

It accepts no registration state or caller-supplied candidate projection. The durable
`registrationstate.ApprovalAttemptComponent` continues to resolve the requested registration and
stored plan itself, advance effective time, reject stale or mismatched state, enforce nonce
uniqueness, and own payload replay. This split preserves the frozen rule that an equivalent
signature replay verifies cryptographically but does not re-admit or resurrect the already retained
record.

## Public-only fixture

The integration known answer binds the SHA-256 digest
`ef268a0b829adc1ce1307203f4b805f63379954ccf41e8e20a7487b6e5acf241` of
`schemas/conformance/v0/execution-plan/ordinary.cbor`, registration `33…33`, installation `11…11`,
epoch digest `22…22`, Supervisor `55…55`, and nonce `66…66`.

The fixture was generated once during an explicitly authorized local test-only signing operation.
The P-256 private key existed only in process memory and was neither printed nor written. The
repository retains only public X/Y coordinates and the signed envelope in
`internal/execution/registeredlifecycle/production_approval_integration_test.go`. It is not a live
key, Keychain item, Secure Enclave key, Apple identity, installation authorization, or product
credential.

## Claim boundary

This is an unwired local conformance integration, not a product approval or execution path. It
creates no endpoint, process, runtime, adapter, VM, guest, output, completion, signature-bearing
evidence, or user-visible success. Fake prepare/create/start/observe/stop/destroy are simulator
state only; fake absence proves no real process-tree teardown. The fixed stores remain local
fault-injection oracles and provide no installed owner, same-UID containment, power-loss, restore,
rollback, continuity, or production durability claim.

## Reproduction

```sh
GOCACHE=/tmp/capsule-go-cache go test \
  ./internal/execution/approvalattempt \
  ./internal/execution/registrationstate \
  ./internal/execution/registeredlifecycle
```

Repository handoff verification remains the complete command set required by `AGENTS.md`.

## Read-only audit disposition

| Requested case | Disposition |
| --- | --- |
| Zero fake effects before atomic consume/create and through lifecycle-record commit | Implemented here by `TestProductionVerifiedApprovalCommitAndRestartConvergeBeforeFakeEffects`, including a test-store wrapper that returns only after `EnsureLifecycle` has durably committed and has observed zero fake calls. |
| Keep `FixtureVerifier` test-only | Implemented here by `TestFixtureVerifierHasNoNonTestGoConsumer`, which parses every non-test Go file and fails if the fixture type or constructor is consumed outside its definition. The integration injects only `ProductionShapedVerifier`. |
| Strict predecode at the real submission integration point | Implemented here by `TestProductionIntegrationRefusesNonProfileSign1BeforeDurableState`: untagged and DER Sign1 forms refuse as `MALFORMED`, a non-equivalent raw-signature mutation refuses as `BINDING`, durable state remains byte-identical, and no attempt or fake effect appears. |
| Low-S/high-S replay after consumption | Implemented here: the low-S form is submitted, its high-S complement converges, the approval is consumed into one attempt, and both forms then return the same consumed record rather than a fresh usable approval. |
| Lost submit reply, exact/equivalent retry, and consumed replay | Implemented here across injected post-commit response loss and store reopen. |
| Nonce reuse with changed payload and old-epoch/binding refusal | Already covered at the durable semantic seam by `TestApprovalAttemptExpiryRollbackNonceAndIdentifierRules` and `TestApprovalSubmissionBindingMatrix`; ADR-0043 production verifier tests independently reject epoch/key authorization mutations. This integration retains one authorized public signed payload and no private signer, so it does not invent a second valid signature solely to duplicate those store tests. |
| Plan, registration, or source replacement | Plan A versus B and registration mismatch are implemented here. Source replacement is structurally absent from approval and execute calls: the stored plan carries its registered source binding, and `Drive` accepts only `AttemptID`; no replacement-source parameter exists. |
| Confirmed versus indeterminate submit/attempt commit outcomes | Already covered by `TestApprovalAttemptFaultAndProcessDeathMatrix`; this integration adds production-verifier response loss after confirmed approval and attempt commits. The verifier seam does not change the store's confirmed/indeterminate transaction mechanics. |
| Confirmed versus indeterminate `Complete` outcomes | Already covered independently by `completioncomposer.TestCrashEdgesResponseLossRestartAndStaleReplay`. Durable completion remains unwired to this no-guest slice, so claiming an integrated `Complete` path here would exceed scope. |
| ADR-0024 versus ADR-0043 signer-policy authority | Resolved by an explicit cumulative-authority cross-reference in ADR-0024: ADR-0024 owns durable admission bindings; accepted ADR-0043 owns the additional Team/role/access-group/key-policy scope. |
| ADR-0040 stop policy and passive native XPC v0 | Dispositioned under adjacent dependency ordering above. Both remain unmerged dependency awareness; neither branch is copied here. |
