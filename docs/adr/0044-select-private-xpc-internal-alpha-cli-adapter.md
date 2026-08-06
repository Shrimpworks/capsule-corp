# ADR-0044: Select one private-XPC internal-alpha CLI submission adapter

- Status: Accepted
- Date: 2026-08-05
- Refines: ADR-0029 and ADR-0040
- Decision scope: passive internal-alpha CLI-to-daemon submission topology and logical contract

## Context

ADR-0040 selects a bounded Apple Development-signed CLI for the owner-only internal alpha and
forbids turning the diagnostic HTTP server into mutation authority. The repository already retains
one exact-one-`main.mjs` `JobProposal`, plan construction mechanics, and the passive atomic
`RegisterPlanV0`/`GetRegisteredPlanV0` authority plane. It did not select the authenticated
CLI-to-daemon adapter that submits those proposal bytes or freeze aggregate flow-control behavior
across that submission and the two existing Supervisor methods.

The choice must preserve ADR-0029's ordering: an OS peer requirement is installed before message
delivery; message-derived service, role, purpose, audience, installation, and epoch select
authority; application bytes cannot choose a method; request IDs are correlation only; and a
missing reply never proves abort. A generic HTTP mutation endpoint, Unix-domain socket identified
by path/mode/PID/EUID, stdin protocol, shared command bus, or multiplexed daemon/Supervisor listener
would weaken those properties.

This decision does not revise ADR-0037's exact passive I0 role inventory. The CLI identity and
service require a later versioned internal-alpha installation profile; they are not silently added
to I0. This decision also does not activate an endpoint, register a service, sign code, use a key,
or connect a product consumer.

## Decision

### One adapter and one method

Select exactly one private-XPC logical adapter from the internal-alpha CLI to the daemon:

| Property | Exact value |
| --- | --- |
| CLI signing identifier candidate | `com.capsulecorp.capsule.cli` |
| daemon service | `com.capsulecorp.capsule.daemon.cli.v0` |
| method | `SubmitMainMJSV0` |
| method version | `0` |
| expected peer role | `internal-alpha-cli` |
| audience | `capsule.daemon.local.v0` |
| purpose | `capsule.ipc.submit-main-mjs.v0` |
| request media type | `application/capsule.job-proposal+json;v=0` |
| request application-data cap | 2,097,152 bytes |
| successful reply | nonzero 16-byte Supervisor-issued `RegistrationID` only |
| method deadline | 10,000 milliseconds from admission |

The future daemon listener must install the exact active peer requirement before activation and
repeat message-derived code identity, current EUID/session, installation, and epoch checks before
body copy or decode. Exact Team/channel, CDHash, entitlement, session, and update values belong to
the later signed installed profile. The passive contract records the requirement but supplies no
peer-authentication evidence.

`SubmitMainMJSV0` carries one exact strict-JSON `JobProposal` document. The current closed decoder
must still reject legacy multifile JavaScript, CommonJS, TypeScript, paths, packages, network,
process, environment, native, backend, image, mount, and execute-time replacement authority. The
adapter does not accept raw plan bytes, role bindings, source manifests, approval bytes, backend
flags, or a caller-selected `RegistrationID`.

After strict decode, semantic resolution, and plan construction, the daemon may invoke only
ADR-0029's method-specific `RegisterPlanV0`. Registration atomically retains the exact plan,
resolved 562-byte role projection, canonical `SourceManifest`, exact `main.mjs`, and registration.
The CLI reply exposes only that Supervisor-issued registration locator. The Broker continues to
fetch Supervisor-owned typed bytes through `GetRegisteredPlanV0`; the CLI does not supply approval
display data.

### Closed flow and delivery behavior

The CLI/daemon service permits at most four live authenticated connections, one in-flight request
per connection, four admitted requests in the daemon, zero queued application requests, and
8,388,608 aggregate in-flight request-data bytes. A missed slot returns `CAPACITY` after peer and
outer-header validation but before proposal body copy/decode or downstream state change.

ADR-0029's Supervisor services retain at most four live authenticated connections each, one
in-flight request per connection, eight admitted requests total, zero queued application requests,
and 2,626,696 aggregate in-flight request-data bytes. `RegisterPlanV0` keeps its 5,000-millisecond
deadline; `GetRegisteredPlanV0` keeps 2,000 milliseconds.

Cancellation before dispatch makes no daemon planner or Supervisor call. Cancellation after
dispatch cancels only response delivery: it cannot decide whether registration committed.
Method-owned deadlines also never infer abort. A stalled downstream call retains its bounded slot
until it returns or the later process-level integrity/restart boundary resolves it; no replacement
work is queued behind it.

Request IDs remain nonzero 16-byte correlation values. They are not registration, idempotency,
replay, authorization, approval, attempt, nonce, or store keys. After response loss, retrying the
same submission is a fresh transport call and may create another registration if the first call
committed. `RegisterPlanV0` therefore remains deliberately fresh; `GetRegisteredPlanV0` remains a
repeatable exact-byte read.

### Diagnostic HTTP remains read-only

`internal/api` and `cmd/capsuled` remain diagnostic-only. No proposal, registration, approval,
attempt, cancellation, execution, or other mutation route may be added there for the internal
alpha. A future public client adapter needs another decision and evidence; it cannot reuse this
private service as a generic command bus.

## Rejected alternatives

- **Diagnostic HTTP submission:** rejected because it makes the read-only diagnostic server a
  mutation/authority surface and lacks the selected OS peer-requirement-before-delivery topology.
- **Unix-domain socket or stdin:** rejected because path/mode/PID/EUID or parent process placement
  does not establish exact enrolled code identity, current session, installation, or epoch.
- **Generic or multiplexed command service:** rejected because application bytes could select
  authority before method-specific caps and role checks.
- **CLI direct to Supervisor:** rejected because it bypasses the daemon's strict public-proposal
  decoding, trusted resolution, and plan-construction role and widens the Supervisor parser TCB.
- **Multiple CLI/SDK/MCP adapters in the alpha:** rejected for this checkpoint. Exactly one client
  adapter keeps aggregate scheduling, installed identity, and response-loss behavior closed.

## Consequences and blockers

The passive adapter and logical envelope can be implemented and verified against repository
fixtures and the unwired authority-plane facade. That slice provides no XPC key spelling, numeric
tag/status encoding, endpoint, service registration, Apple signing, installed identity/session,
protected state, Broker UI/key use, runtime, backend, guest, or product admission evidence.

The later installed profile must deliberately version ADR-0037's inventory, enroll the exact CLI
and daemon peer records, prove the authentication-before-delivery and update/session matrix, and
retain no-state overload/deadline/cancellation/response-loss evidence. If supported private-XPC
reachability requires a different component placement, shared App Group, helper, or listener
topology, this exact candidate stops and a new ADR is required.

The later passive
[S3 native-contract prerequisite](../AUTHENTICATED_LOCAL_IPC_S3_NATIVE_CONTRACT.md) freezes the
exact dictionary keys/types, numeric method/status/reason tags, versions, key counts, and passive
refusal mappings for this method plus `RegisterPlanV0` and `GetRegisteredPlanV0`. It does not
change this decision's topology and supplies no listener, peer-authentication, signing, installed,
or product evidence.
