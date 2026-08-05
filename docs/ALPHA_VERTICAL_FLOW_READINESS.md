# Alpha vertical-flow readiness

## Status and scope

The owner-only hostile-`.mjs` internal alpha is `IN_PROGRESS — TRENDING_GOOD`; product admission is
`BLOCKED`. Accepted [ADR-0040](adr/0040-freeze-owner-only-internal-alpha-posture.md) narrows the
target to one exact `main.mjs`, bounded inline JSON, native human approval, one fresh governed
guest per `AttemptID`, concurrency one, and a deliberately limited fixed-store posture.

The earlier fake-backend vertical flow remains useful integration work, but it is not a separate
release milestone and no longer waits on the host Source Validator. The fixed benign owned-guest
checkpoint is a separate engineering experiment, not product alpha. See the full
[architecture and release audit](ALPHA_ARCHITECTURE_AND_RELEASE_AUDIT.md).

## What exists and what is missing

| Stage | Retained contract or mechanic | Scoped status | Missing product boundary |
| --- | --- | --- | --- |
| Submit one `main.mjs` | `MJSMainSource`, canonical `SourceManifest`, exact byte/cap/digest fixtures | `PASSED` for passive bytes | No selected authenticated bounded CLI-to-daemon submission adapter. Diagnostic HTTP remains read-only. |
| Build the plan | `ExecutionPlan`, strict wrappers, pure unwired resolver/builder candidates | `PASSED` for passive/unwired mechanics | Replace the broad legacy proposal with exactly one `main.mjs`; bind only supported limits and exact runtime/profile identity. |
| Register source and plan | Local registration/store mechanics | `PASSED` only for their current unwired fields | Atomically retain exact plan, complete bindings, manifest, and source bytes; expose method-specific authenticated registration/fetch. |
| Fetch and render approval | ADR-0029 method design and native app scaffolding | `BLOCKED` | Broker must fetch Supervisor-owned bytes, render fixed typed facts, and treat an opaque RegistrationID only as an untrusted locator. |
| Sign and retain approval | Approval fixtures plus local consume/create/replay mechanics | `PASSED` only for fixture/store evidence | Fresh `LAContext`, explicit user-presence key policy, no fallback, strict production COSE/key authorization, and installed Broker/Supervisor integration. |
| Create and drive attempt | Durable attempt/lifecycle/FakeBackend mechanics | `PASSED` for no-guest local mechanics | Authenticated product call, protected owner/store, real sealed adapter, stable backend identity, restart reconciliation, and one fresh guest. |
| Compose completion | Intended transcript/receipt/fixed-summary model | `BLOCKED` | Supervisor-derived compositor joining completion-last, integrity, lifecycle, teardown, and authoritative process-tree absence. |

The Oxc Source Validator's passive, artifact, and R3 signed inactive evidence remains retained.
Product R4/R5 is `BLOCKED`, and exact R4-v1 candidates are `NO_GO`. ADR-0040 moves this control to
post-alpha defense-in-depth; the internal alpha instead requires runtime-level syntax/module
refusal and physical absence of host authority.

## Dependency graph

```text
governance promotion hardening
  -> successor runnable profile closure
  -> fixed benign owned guest checkpoint

bounded authenticated CLI adapter
  -> exact main.mjs + SourceManifest
  -> narrow plan construction
  -> atomic plan/bindings/manifest/source registration
  -> Broker fetches Supervisor-owned bytes
  -> fixed rendering + fresh user-presence approval
  -> production approval verification
  -> one-use AttemptID creation
  -> sealed real adapter + fresh governed guest
  -> completion-last + teardown/absence composition
  -> fixed AgentExecutionSummary
```

The authority path may be integrated through `FakeBackend.CreatesGuest() == false` before the real
adapter is connected. The fixed guest may be tested earlier only with the sealed known answer and
separate authorization; neither path alone is product alpha.

## Corrections every slice must preserve

1. Do not turn the diagnostic HTTP server into a mutation or authority surface.
2. Registration remains deliberately fresh; request IDs are correlation, not generic authority
   deduplication.
3. `FixtureVerifier` is test-only and cannot verify newly signed production approvals.
4. The legacy multi-file JavaScript/TypeScript proposal and stale receipt schema are not alpha
   consumers.
5. Approval renders Supervisor-retained typed bytes, never daemon-supplied prose.
6. Execute-time calls accept only Supervisor-issued identifiers, never replacement bytes, paths,
   images, mounts, resources, or backend flags.
7. A fresh disposable guest is created for every attempt. Exit zero or EOF never proves success.
8. The fixed-store exception has no restore, rollback, continuity, secure-deletion, or hostile
   same-UID host-process claim and stops at ADR-0040's thresholds.
9. `--jitless` does not disable `eval`/`Function`; loader absence and string-code-generation denial
   are independent controls with restoration tests.

## Ordered closure

1. Require governed admission checks, generic governed-branch CI filters, administrator
   enforcement, and a no-rewrite bad-promotion runbook before consuming another promoted head.
2. Freeze the successor runnable composition and run one fixed benign owned guest.
3. Freeze the narrow proposal, atomic custody/fetch object, selected adapter envelope, and complete
   field authority without preserving broad legacy acceptance.
4. Implement authenticated role-specific IPC, protected Supervisor state, Broker rendering,
   production signing/verification, and the bounded fixed-store policy.
5. Connect approval to attempt/lifecycle first against FakeBackend, then through the admitted sealed
   real adapter.
6. Implement the completion/transcript/fixed-summary compositor.
7. Run the minimum hostile source/authority/transport/root/lifecycle/restoration corpus in the exact
   signed-installed profile.

## Honest exit condition

The internal alpha becomes `PASSED` only when the exact Supervisor-retained source shown to and
approved by the owner is the source executed once in a fresh governed guest; response loss and
restart converge; unsupported authority is physically absent; and every terminal result contains
independent teardown/absence truth. Developer ID distribution, notarization, automatic update,
multi-host support, restore, and production storage remain external-alpha gates.
