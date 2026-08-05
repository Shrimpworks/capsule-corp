# Alpha vertical-flow readiness

## Status and scope

The first fake-backend alpha vertical flow is `BLOCKED`. Many of its byte contracts and local
mechanics are `PASSED` in narrower scopes, but there is no product path from submitted `main.mjs`
bytes through validation, planning, registration, approval, attempt creation, lifecycle drive, and
a fixed completion summary.

This document is a dependency and consumer map. It does not select a new endpoint, widen a
component's authority, admit a runtime or backend, or turn passive evidence into a product claim.
The governed guest path remains a separate workstream; the first vertical integration described
here ends at `FakeBackend.CreatesGuest() == false`.

## What exists and what is missing

| Stage | Retained contract or mechanic | Scoped status | Missing product boundary |
| --- | --- | --- | --- |
| Submit one `main.mjs` | [`v0candidate.MJSMainSource` and `SourceManifest`](../internal/protocol/v0candidate/source_manifest.go), exact byte/cap/digest fixtures | `PASSED` for passive bytes and manifests | No selected agent-facing submission adapter consumes these bytes. [`internal/api`](../internal/api/server.go) remains a read-only, unauthenticated diagnostic surface. |
| Validate source | [`sourcevalidatorpassive.ValidationRequest` and `ValidationResult`](../internal/protocol/sourcevalidatorpassive/contract.go), Accepted ADR-0035/0036, role-separated passive and unsigned artifacts | `PASSED` in their exact passive/construction scopes | Signed installed launchers, private reachability, confinement/resource/residue evidence, and independent daemon and Broker consumers remain `BLOCKED`. |
| Build the plan | `v0candidate.ExecutionPlan`, strict wrapper decoding, the pure TypeScript [`job-proposal-resolver`](../packages/protocol/src/job-proposal-resolver.ts) and [`execution-plan-builder`](../packages/protocol/src/execution-plan-builder.ts), candidate proposal/schema fixtures | `PASSED` for the pure unwired resolver/builder, candidate bytes, registration handoff, and decoding | No daemon consumer supplies independently validated source facts and trusted bindings, invokes the builder, or sends its handoff to the Supervisor. Product invocation waits for the daemon-side Source Validator consumer. |
| Register the plan | [`registrationstate.Component.RegisterPlan`](../internal/execution/registrationstate/component.go) and fixed-store fault/reopen mechanics | `PASSED` as unwired local mechanics | No authenticated `RegisterPlanV0` method facade or installed daemon-to-Supervisor XPC consumer reaches it. |
| Fetch and render approval | Proposed ADR-0029 `GetRegisteredPlanV0`; retained Swift installation/status shell | `PASSED` only for the separate passive method and app-construction scopes | The Approval Broker has no product fetch, independent validation, rendering, user-presence, or approval-signing flow. |
| Verify and retain approval | `approvalattempt.ApprovalGrant`, [`FixtureVerifier`](../internal/execution/approvalattempt/fixture_verifier.go), approval/attempt store mechanics | `PASSED` for retained-vector and local store evidence | `FixtureVerifier` is not a production COSE verifier. The authorized-key, LocalAuthentication/Keychain signer, narrow production verifier, and Broker/Supervisor integration remain `BLOCKED`. |
| Create one attempt | `ApprovalAttemptComponent.RequestAttempt` / `Attempt` and immutable `CreatedAttempt` | `PASSED` as unwired local mechanics | No authenticated product call composes the registered plan and production approval into this path. |
| Drive and recover lifecycle | `registeredlifecycle.Component.Drive` / `Recover`, durable lifecycle records, owner coordination, `FakeBackend` | `PASSED` for the no-guest local mechanic | No product consumer invokes the lifecycle from committed attempt truth. A real backend remains deliberately outside this first alpha integration. |
| Compose completion | Pre-freeze receipt schema and intended `ExecutionReceipt`, `EnforcementTranscript`, `ArtifactManifest`, and fixed `AgentExecutionSummary` model | `BLOCKED` | No Go/Swift/TypeScript compositor produces a receipt/transcript or projects a terminal attempt into the fixed agent summary. Schema presence alone is not implementation evidence. |

## Dependency graph

```text
selected client adapter
  -> exact main.mjs bytes and SourceManifest
  -> daemon-role Source Validator consumer
  -> daemon invocation of the passed pure plan builder
  -> authenticated RegisterPlanV0 facade
  -> Broker GetRegisteredPlanV0 + independent Source Validator consumer
  -> Broker render + user-presence approval signing
  -> Supervisor production approval verification
  -> RequestAttemptV0
  -> Drive/Recover against FakeBackend
  -> receipt/transcript composition
  -> fixed AgentExecutionSummary
```

The order is intentional. Registration does not precede source validation; approval does not render
daemon-supplied display text; execution never accepts replacement plan or backend bytes; and the
completion summary is derived from retained Supervisor truth rather than caller or backend claims.

## Corrections that future slices must preserve

1. **Do not silently turn the diagnostic HTTP server into the job authority boundary.** The public
   client adapter may eventually be MCP, CLI, SDK, or HTTP, but its authentication, aggregate
   connection/concurrency/byte/deadline envelope, cancellation, overload, and forwarding behavior
   must be selected and tested before a mutation endpoint is activated.
2. **Registration is deliberately fresh, not idempotent.** Every successful `RegisterPlanV0` core
   call creates a fresh registration. A lost response may leave an unreachable registration that
   expires normally; transport request IDs do not add a deduplication ledger.
3. **`FixtureVerifier` cannot verify newly signed approvals.** It performs exact lookup against
   retained vectors. A test may use those vectors, but production approval requires the separately
   reviewed signer, authorized-key resolution, strict envelope parser, and verifier.
4. **Do not restore the rejected inherited-helper Source Validator topology.** Product consumers use
   the two role-specific private App-Sandboxed launcher design from Accepted ADR-0036 after its
   signed installed and confinement evidence passes.
5. **Schemas are not consumers.** Existing receipt or summary schemas describe pre-freeze shape;
   they do not prove that evidence is composed, signed, stored, disclosed correctly, or projected
   from terminal lifecycle truth.
6. **Do not mislabel replay semantics.** `SubmitApprovalV0` converges by verified canonical payload
   and resolved signer identity, and `RequestAttemptV0` converges by
   `(RegistrationID, ApprovalReference)`; `RegisterPlanV0` does not.

## Ordered alpha closure

1. Complete the Source Validator signed-install, confinement/resource/residue, daemon-consumer, and
   Broker-consumer slices without changing V0/V1/V2 retained bytes.
2. Hold the M2/S1 checkpoint: freeze the complete plan construction and method-specific
   `RegisterPlanV0` / `GetRegisteredPlanV0` fixtures, caps, field authority, refusal behavior, and
   response-loss oracles.
3. Implement the native-fronted authenticated XPC method facades so the daemon can register and the
   Broker can fetch only Supervisor-owned typed bytes.
4. Implement read-only Broker rendering before adding production approval key use. Then add the
   LocalAuthentication/Keychain signer and narrow Supervisor approval verifier with exact
   authorization, replay, and user-presence evidence.
5. Connect committed approval to attempt creation and drive the existing durable lifecycle against
   `FakeBackend`, preserving AttemptID-only authority and startup recovery from store truth.
6. Freeze and implement evidence/receipt composition and the fixed agent summary. The fake vertical
   integration passes only when response loss, restart, refusal, unresolved, quarantine, and
   terminal teardown cases produce the expected bounded public result without guest-controlled
   diagnostic leakage.
7. Keep the governed `deno_core`/libkrun C2 path separate until runtime/profile admission. Replacing
   `FakeBackend` is not required to prove the first consumer-to-Supervisor alpha integration.

## Honest alpha exit condition

The fake-backend alpha vertical flow becomes `PASSED` only when one selected client adapter can
submit the fixed `.mjs` proposal and the installed product components complete the sequence above
through a fixed summary, with authenticated role-specific IPC, real Broker approval authority,
durable response-loss/restart recovery, and no guest creation. Passing that flow will still not
admit a runtime, backend, hostile-code execution profile, updater, or production distribution.
