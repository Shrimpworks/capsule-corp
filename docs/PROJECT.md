# Project Definition

## Summary

Capsule is a local-first trusted execution platform for bounded JavaScript and TypeScript jobs
proposed by AI agents. It separates untrusted planning, trusted human authorization, user-content
custody, and hostile-guest execution so that compromising the agent-facing service does not also
grant approval, content, or launch authority.

The task—not the container, VM, shell, or development environment—is the public abstraction.

## Current status

Capsule is an architecture and buildable scaffold. It does not yet implement the intended security
boundary and must not be used to execute hostile code without another trusted sandbox.

The current JSON Schemas and TypeScript `Job` interfaces are pre-freeze scaffolding. They describe
the repository's current API surface, but they intentionally do not define the target v0 protocol.
Blocking feasibility spikes must determine what the platform can actually enforce before those
contracts are replaced and frozen.

Phase 2A now includes verified passive candidates for the first narrow `JobProposal`, minimum
`ExecutionPlan`, and `PlanRegistration`, with byte-exact internal fixtures and Go/TypeScript
decoded views. They do not replace the scaffold, activate an endpoint, or authorize execution.

Phase 2B now retains a closed 82-rule, 262-case, 368-fixture conformance corpus for raw decoding,
proposal/source/input semantics, exact plan and registration bytes, registration-state transitions,
and the unwired ADR-0024 approval/attempt boundary. The unwired implementation verifies 80
TypeScript proposal targets: 62 raw/schema cases plus all 18 semantic-resolution cases. It also
verifies 177 Go targets: 81 internal-CBOR/wrapper cases, all 40 registration-state cases, 44
passive approval/attempt cases, and 12 fixed-store transition oracles. The approval/attempt work
adds distinct identifiers/references, the fixed internal classification vocabulary, the exact
candidate approval known answer, a bounded fixture-only verifier, and one unwired transactionally
colocated fixed registration/approval/attempt store. The no-guest fake lifecycle now resolves and
recovers only committed `AttemptID` records, revalidates exact plan and copied bindings before fake
prepare, and retains 12 top-level focused lifecycle tests for binding, replay, concurrency, fault,
and startup-recovery behavior. Its lifecycle store remains bounded single-process memory. There is
no consumer, authenticated IPC, production approval, evidence, real backend, runtime, or guest. A
focused unwired TypeScript Task 3C slice now constructs and deterministically encodes the minimum
`ExecutionPlan` from only Task 3B provenance-bearing plan inputs and separately issued trusted role
bindings. The next
focused slice now prepares a defensive exact-byte/complete-role handoff and exercises those values
against the real Go `registrationstate` component through a local-only conformance command. Go
independently predecodes, role-binds, hashes, and retains the 530-byte known answer. This is not a
product IPC implementation: Swift, authenticated cross-process transport, consumers, and endpoints
remain pending. See the
[unwired decoder checkpoint](PHASE_2_UNWIRED_DECODER_CHECKPOINT.md). The latest local Gate C
checkpoint rejects stock Bun 1.3.14, its governed-construction branch, hardened full Deno v2.9.4,
and the tested minimal `deno_core` 0.409.0 construction for the required runtime-authority contract.
No first runtime is selected, `RUNTIME-001` remains unsupported, and the next bounded runtime
experiment must retain the unchanged prohibited-power contract. The exact libkrun block-root
profile also remains unsupported because the smallest `NullFs` removal prevents guest bootstrap.
See the
[P0-0 construction review](../experiments/gate-c-bun-runtime-authority/governed-closure/CONSTRUCTION_REVIEW.md)
and [Deno-family disposition](../experiments/gate-c-deno-runtime-authority/RESULTS.md)
and [parallel-task checkpoint](PHASE_2B_GATE_C_TASK_GROUP_CHECKPOINT.md).

## Problem

Agents frequently need to perform work that is more reliable as code: transform data, validate
configuration, analyze source, generate reports, call an approved API, or run bounded checks.
Executing generated code directly on a user's machine gives it excessive ambient authority.
Generic cloud sandboxes can be opaque, heavyweight, old, or disconnected from local workflows.

Capsule aims to make constrained local execution understandable and testable enough to become the
default place for agent-generated JS/TS tasks.

## Goal

Provide a fast execution boundary in which:

- The agent-facing daemon can propose and plan work but cannot approve or launch it.
- A trusted native Broker owns user presence, file selection, and user-only content.
- A small Execution Supervisor independently validates registered plans and is the only component
  allowed to create a hostile guest.
- Every approval authorizes one exact, immutable plan registration and at most one attempt.
- Every installation has an enrolled trust domain with purpose-separated keys and sequence-ordered
  component identity.
- DIDs can represent Capsule principals externally without becoming the local authorization root.
- Modern runtimes are replaceable, pinned artifacts rather than control-plane dependencies.
- Isolation backends can evolve without changing public job semantics.
- User-defined limits are exact, visible, and never silently expanded or rewritten.
- Results and artifacts have explicit validation and audience policies.
- Every execution produces attributable approval and enforcement evidence with honest limitations.

## Product scope and initial wedge

The platform scope is broader than file processing: it is intended for bounded agent-generated
JS/TS tasks. The first executable slice is deliberately smaller:

- Local macOS control experience
- A pinned JS/TS runtime selected only after P0-0 authority closure
- One-shot, dependency-free execution
- Inline JSON input and bounded JSON output
- Explicit prepare, register, human-readable approval, attempt, and execute phases
- Per-installation identity with offline local authorization
- No network, subprocesses, environment inheritance, native addons, FFI, macros, inspector, or
  dynamic installation
- A fixed, low-bandwidth agent summary; full output remains user-only by default
- Development posture until the exact backend configuration passes its retained attack corpus

Regular-file snapshots and JSONL, text, and CSV artifacts follow only after the inline JSON slice
exercises the authority boundaries successfully.

## Primary users

- Developers building agents that need a safe computation surface
- Desktop AI users who want bounded local execution
- Teams introducing policy around generated code
- Tool authors who need a vendor-neutral execution contract

## Non-goals for v0

- General remote shell access
- Long-lived development environments
- Browser automation
- Docker-in-Docker
- Arbitrary languages or package installation
- Network access, subprocesses, secrets, or environment inheritance
- Arbitrary directories, repositories, archives, devices, sockets, or special files
- Background services and public preview URLs
- Multi-tenant hosted scheduling
- Portable multi-device identity or automated recovery
- Proof that guest code is correct or aligned with user intent
- Proof that the local kernel, Secure Enclave, hypervisor, or correctly signed program logic is
  uncompromised
- Proof that permitted outputs, metadata, or timing cannot encode granted input data

## Principles

1. **No ambient authority.** Every external effect requires an explicit grant.
2. **Separate authorities.** Planning, approval/content custody, and execution enforcement do not
   share one compromise boundary.
3. **Fail closed.** Unknown capabilities, identities, profiles, transitions, and backend controls
   are rejected.
4. **External isolation is mandatory.** Language permissions are supplemental hardening.
5. **Approval is exact and attempt-bound.** The user approves typed registered plan bytes; a grant
   can produce at most one attempt.
6. **Data authority uses handles.** Agents cannot turn paths, URLs, names, or identifiers into
   authority.
7. **Identity is not authorization.** Installation policy gives enrolled keys narrow purposes; a
   DID is an optional interoperable identifier.
8. **Trust changes are explicit.** Component identities and policy are bound into signed,
   sequence-ordered trust epochs with crash-safe update and repair rules.
9. **Egress is a capability.** Logs, structured results, filenames, metrics, and artifacts are
   untrusted and observable channels.
10. **Limits belong to the user.** Defaults and ceilings come from trusted policy and are enforced
    exactly or the attempt is refused.
11. **Evidence is attributable, not magical.** Receipts compose signed claims and retained test
    evidence; they are not independent platform attestation.
12. **Security claims are testable.** No backend, profile, or component posture advances without
    exact mechanisms, adversarial tests, and retained evidence.

## Agreed v0 direction

- The agent-facing Go daemon performs strict proposal validation and planning only.
- A native macOS Trusted Host Broker contains logically separate Approval and Content Broker
  interfaces and has no agent-facing endpoint.
- The Execution Supervisor is the sole hostile-guest launch authority and independently enforces
  non-overridable hard-safety rules.
- Execute operations accept a Supervisor-issued registration identifier, never replacement plan
  bytes.
- Approval binds the plan digest, registration, installation, trust epoch, expected Supervisor,
  attempt nonce, purpose, audience, and expiry.
- The normative local identity is a random installation ID plus locally authorized public keys.
  DIDs are first-class optional representations for interoperability and exported evidence.
- External release and profile trust uses pinned TUF roots. Live execution consumes a compact,
  verified local trust snapshot and performs no network trust lookup.
- Direct Apple Containerization is retained only as a macOS development backend after failing the
  durable identity/recovery gate. A follow-up libkrun/Hypervisor.framework spike conditionally
  passed exact-process lifecycle, isolation, App Sandbox, controller-crash, and Bun checks, making
  it the lead native Apple candidate under evaluation. Its readiness tracks and the subsequent P0
  review found that runtime-authority closure, immutable runtime-root custody, the block-root
  `NullFs` surface, typed port transport/completion, and complete installed-bundle admission must
  close before one exact development profile freezes. Source/input and bounded inline JSON can use
  attempt-bound virtio-console ports in the first slice; an ext4/raw-image parser is deferred until
  file artifacts. OCI plus gVisor remains an independent unvalidated candidate and contingency
  until both exact profiles run the shared attack corpus.
- Go remains the daemon language; Swift is preferred for the native Broker. Supervisor language and
  privilege remain gated by the backend adapter; the current macOS service evidence supports an
  unprivileged per-user Swift component with no host-root helper.

## Success criteria

The first functional milestone succeeds when a client can submit a dependency-free inline JSON
job, obtain a Supervisor registration, receive explicit user-presence approval from the Broker,
consume that grant for one attempt, run Bun in a disposable development backend, release only a
fixed agent summary, record backend-specific teardown evidence, and compose a receipt from Broker
approval evidence and a Supervisor enforcement transcript.

The first validated-local milestone additionally requires the exact backend, runtime bundle,
component identities, and host configuration to survive the documented adversarial corpus. An
unresolved teardown or integrity state fails closed and cannot be reported as ordinary success.

## Near-term method

The first feasibility program is complete. Capsule now proceeds on two bounded lanes: implement
backend-independent contracts, registration, approval consumption, fake-backend recovery, evidence,
and inline JSON; and close only the remaining fail-fast P0 gates before connecting user bytes to a
real libkrun adapter. Prototype code may be discarded. Fixtures, measurements, limitations, and the
resulting ADR decisions remain durable project evidence.

See [Feasibility Spikes](FEASIBILITY_SPIKES.md), [Technical Design](TECHNICAL_DESIGN.md), and the
[Roadmap](ROADMAP.md). The exact branch point is recorded in the
[Gate C P0 reconciliation](GATE_C_P0_RECONCILIATION.md).
