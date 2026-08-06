# Passive internal-alpha product-adapter contract result

Date: 2026-08-05

Work item: developer-signed internal-alpha CLI submission adapter and registration/fetch flow
contract.

Status: `PASSED` for the exact passive, unwired logical-contract and in-process flow-harness scope.
The installed authenticated product-adapter parent remains `BLOCKED`.

Scope: repository fixtures, an unexported in-process harness, the unwired authority-plane facade,
and the diagnostic HTTP route table. No endpoint, service registration, Apple signature, key,
approval, process launch, runtime, backend, VM, guest, or product consumer was created.

Evidence or reason: Accepted [ADR-0044](adr/0044-select-private-xpc-internal-alpha-cli-adapter.md)
selects exactly one private-XPC logical adapter: `SubmitMainMJSV0` from
`com.capsulecorp.capsule.cli` to `com.capsulecorp.capsule.daemon.cli.v0`. Generated Node fixtures
and Go tests agree on the request/reply known answers, exact caps, method-derived role/purpose/
audience/install/epoch bindings, correlation-only request IDs, copy ownership, exact source
lineage into the existing atomic registration/fetch plane, capacity, concurrency, zero queue,
deadline, cancellation, downstream stall, refusal, and response-loss dispositions.

The submission request carries 1..2,097,152 exact `JobProposal` JSON bytes and the reply carries
only a nonzero 16-byte Supervisor-issued `RegistrationID`. The retained known answer proves the
proposal's decoded `main.mjs` bytes equal the registered source fixture. Existing registration
continues to commit plan, resolved bindings, manifest, source, and registration atomically; fetch
returns defensive retained copies. No legacy multifile/TypeScript acceptance or execute-time byte
route was added.

Flow limits are frozen as follows:

- CLI/daemon: four authenticated connections, one in flight per connection, four admitted total,
  zero application queue, 8,388,608 in-flight request bytes, 10,000-millisecond deadline.
- Supervisor: four authenticated connections per role service, one in flight per connection,
  eight admitted total, zero application queue, 2,626,696 in-flight request bytes;
  `RegisterPlanV0` remains 5,000 milliseconds and `GetRegisteredPlanV0` remains 2,000.

Pre-dispatch cancellation and every pre-core refusal have zero authority/state/effect. After
dispatch, cancellation, deadline, interruption, or missing reply never proves commit or abort.
Stalled work retains its bounded slot and queues nothing. Submission/registration retry may create
a fresh registration; fetch remains repeatable by `RegistrationID`.

The diagnostic HTTP server still exposes only `GET /healthz`, `GET /v1/version`, and
`GET /v1/runtimes`; focused tests reject candidate submission/mutation paths.

Remaining work: the one-time native XPC harness, actual peer requirement and message-derived
identity checks, Apple Development signing/profile enrollment, versioned installed
CLI inventory, daemon planner consumer, protected Supervisor state, Broker rendering/signing,
production approval verification, attempt/lifecycle integration, runtime/backend admission, and
completion composition.

Blocker and owner: installed authenticated product use is blocked on the macOS installation/
Supervisor/Broker owners supplying the versioned signed profile and S3-S5 evidence; runtime and
backend owners separately own product admission.

Next action: run a separately authorized one-time native passive/ad-hoc transport harness in
`Shrimpworks/capsule-experiments` from the now-frozen exact dictionary records without activating a
product service.

Parent status: owner-only hostile-`.mjs` internal alpha remains
`IN_PROGRESS — TRENDING_GOOD`; product admission remains `BLOCKED`.

Successor prerequisite: the later
[S3 native contract](AUTHENTICATED_LOCAL_IPC_S3_NATIVE_CONTRACT.md) is `PASSED` for exact passive
XPC dictionary keys/types, numeric method/status/reason tags, versions, key counts, and refusal
fixtures for all three methods. The one-time native harness, peer-authentication evidence, installed
profile, consumers, and product admission remain `BLOCKED`.
