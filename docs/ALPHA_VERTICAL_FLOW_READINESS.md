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

| Checkpoint | Status | Boundary |
| --- | --- | --- |
| First fixed benign owned guest | `PASSED` | One separately authorized [v19 experimental successor](FIRST_OWNED_GUEST_EXECUTION_CHECKPOINT.md) booted and ran the fixed known answer, matched its bounded console proof, exited/reaped normally, and completed unlinked-root teardown. It is not product alpha, accepts no user source, and does not close typed transport or admission. |
| Hostile owner-only internal alpha | `IN_PROGRESS — TRENDING_GOOD` | One named owner Mac, manual Apple Development installation, exact `main.mjs` plus bounded inline JSON, and one fresh guest per attempt. Product admission remains `BLOCKED`. |
| Later external alpha | `BLOCKED` | Developer ID/notarization/update, clean-host/minimum-OS, F6/restore/continuity, and public-distribution evidence are absent. |

## What exists and what is missing

| Stage | Retained contract or mechanic | Scoped status | Missing product boundary |
| --- | --- | --- | --- |
| Submit one `main.mjs` | Exact-one proposal plus [selected passive private-XPC CLI adapter](INTERNAL_ALPHA_PRODUCT_ADAPTER_PASSIVE_CONTRACT.md), bounded logical flow, and the [exact passive native dictionary prerequisite](AUTHENTICATED_LOCAL_IPC_S3_NATIVE_CONTRACT.md) | `PASSED` for the passive/unwired contract | Run the separately authorized native harness, then implement peer authentication, the versioned signed CLI inventory, and the daemon consumer. Diagnostic HTTP remains read-only. |
| Build the plan | Exact-one-`main.mjs` proposal, `ExecutionPlan`, strict wrappers, pure resolver/builder | `PASSED` for passive/unwired mechanics | Connect only the selected method-specific adapter and bind admitted limits plus exact runtime/profile identity. |
| Register source and plan | Atomic passive fixed-store registration/fetch facade | `PASSED` for the unwired oracle | Expose it only through method-specific authenticated product registration/fetch backed by protected production state. |
| Fetch and render approval | Accepted ADR-0043 read-only projection over passive Supervisor readback | `PASSED` for the unwired projection | Authenticated product fetch, native UI, installed spoof/focus/cancel evidence, and Supervisor custody of inline-input content if it must be shown. |
| Sign, verify, and retain approval | ADR-0043 public-key-only strict verifier connected to local consume/create/replay mechanics by one exact-plan public fixture | `PASSED` for the unwired public-key/store integration | Fresh live `LAContext`, Secure Enclave key operation, installed key authorization and same-byte Broker/Supervisor product integration. |
| Create and drive attempt | [Public-key approval to FakeBackend integration](APPROVAL_FAKE_LIFECYCLE_INTEGRATION.md) plus durable attempt/lifecycle mechanics | `PASSED` for the no-guest local integration | Authenticated product call, protected owner/store, real sealed adapter, stable backend identity and reconciliation, and one fresh guest. |
| Compose completion | [Unwired fixed FakeBackend compositor foundation](COMPLETION_COMPOSITOR_FOUNDATION.md) | `PASSED` for the read-only no-guest foundation | No durable completion producer/store transaction, guest port, integrity/result producer, real teardown/absence proof, authenticated consumer, signing, receipt, or installed composition. |

The Oxc Source Validator's passive, artifact, and R3 signed inactive evidence remains retained.
Product R4/R5 is `BLOCKED`, and exact R4-v1 candidates are `NO_GO`. ADR-0040 moves this control to
post-alpha defense-in-depth; the internal alpha instead requires runtime-level syntax/module
refusal and physical absence of host authority.

## Dependency graph

```text
governance promotions PASSED
  -> C2B v3 passive successor PASSED
  -> C2B v4 build/static materialization PASSED
  -> fixed benign owned guest v19 PASSED
  -> fixed denial-test v20 no-launch materialization PASSED
  -> fixed denial-test v20 execution BLOCKED pre-ready with exact stage unknown
  -> fixed denial-test v21 diagnostic materialization PASSED
  -> fixed denial-test v21 execution BLOCKED on ready-EOF evidence convergence
  -> fixed denial-test v22 convergence materialization PASSED
  -> fixed denial-test v22 execution BLOCKED at preflight-root-sha256
  -> fixed denial-test v23 hash-diagnostic materialization PASSED
  -> fixed denial-test v23 root-digest localization PASSED / hostile execution BLOCKED
  -> fixed denial-test v24 corrected materialization PASSED
  -> fixed denial-test v24 early controls PASSED / full corpus BLOCKED at vsock
  -> fixed denial-test v25 vsock-diagnostic materialization PASSED
  -> fixed denial-test v25 runtime candidate NO_GO before launch
  -> fixed denial-test v26 consolidated materialization PASSED
  -> fixed denial-test v26 localization PASSED / full corpus BLOCKED at passive inventory
  -> fixed denial-test v27 passive-network correction PASSED
  -> fixed denial-test v27 execution BLOCKED on fresh exact authorization

passive bounded CLI adapter + native dictionary prerequisite PASSED
  -> native authenticated CLI adapter
  -> exact main.mjs + SourceManifest
  -> narrow plan construction
  -> passive atomic custody/fetch PASSED
  -> authenticated product registration/fetch
  -> passive fixed rendering + public-key verification PASSED
  -> public-key grant + one-use AttemptID + no-guest FakeBackend PASSED
  -> authenticated Broker fetch + installed fixed rendering
  -> fresh user-presence signing + installed verification
  -> one-use AttemptID creation
  -> sealed real adapter + fresh governed guest
  -> completion-last + teardown/absence composition
  -> fixed AgentExecutionSummary
```

The authority path may be integrated through `FakeBackend.CreatesGuest() == false` before the real
adapter is connected. The fixed guest was tested separately with a sealed known answer and exact
one-use authorization; neither path alone is product alpha.

## Corrections every slice must preserve

1. Do not turn the diagnostic HTTP server into a mutation or authority surface.
2. Registration remains deliberately fresh; request IDs are correlation, not generic authority
   deduplication.
3. `FixtureVerifier` is test-only. ADR-0043's cryptographic verifier is production-shaped but
   unwired and cannot authorize a product approval path.
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

1. Preserve the passed governed fork promotions, C2B v3/v4 passive/static successors, and exact
   v19 fixed-owned-guest checkpoint without reclassifying its diagnostic transport as product
   evidence.
2. Preserve the passed exact-one-`main.mjs` proposal, selected passive CLI adapter, and atomic
   custody/fetch oracle without restoring broad legacy acceptance.
3. Preserve ADR-0043's passive projection/verifier and implement the frozen authenticated
   role-specific IPC, protected Supervisor state, native Broker UI, live signing, installed
   verification, and the bounded fixed-store policy.
4. Preserve the passed public-key approval/attempt/FakeBackend integration and connect the same
   ordering only through the admitted sealed real adapter.
5. Implement the completion/transcript/fixed-summary compositor.
6. Connect the unwired completion/transcript/fixed-summary compositor only after implementing the
   Supervisor-owned durable completion producer and the real integrity, teardown, and absence
   facts it must consume.
7. Run the minimum hostile source/authority/transport/root/lifecycle/restoration corpus in the exact
   signed-installed profile.

## Honest exit condition

The internal alpha becomes `PASSED` only when the exact Supervisor-retained source shown to and
approved by the owner is the source executed once in a fresh governed guest; response loss and
restart converge; unsupported authority is physically absent; and every terminal result contains
independent teardown/absence truth. Developer ID distribution, notarization, automatic update,
multi-host support, restore, and production storage remain external-alpha gates.
