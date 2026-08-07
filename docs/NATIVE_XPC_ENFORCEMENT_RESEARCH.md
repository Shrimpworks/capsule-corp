# Native XPC enforcement research

Date: 2026-08-07

Status: `PASSED`

Parent S3 controlled harness: `BLOCKED` on separate exact authorization

Parent installed authenticated IPC: `BLOCKED`

## Defensive and authorized scope

Research the supported macOS API and ordering needed to defensively validate Capsule's
authenticated local IPC boundary. The work used Apple public documentation and declarations from
the locally installed macOS 26.5 SDK in Xcode 26.6. It created no listener, connection, process,
service, credential, Keychain item, protected root, store, backend, runtime, or guest. It accessed
no unrelated system, identity, credential, or data.

## Question

Which supported native APIs and ordering can the separately authorized S3 harness use to test:

- a peer requirement installed before application message delivery;
- message-derived running-code identity;
- effective-user and audit-session binding;
- bounded copy and flow ownership;
- cancellation, interruption, deadline, and response-loss behavior; and
- process-fault convergence without treating transport state as durable truth?

## Primary sources and environment

Public Apple sources:

- [`xpc_connection_set_peer_code_signing_requirement`](https://developer.apple.com/documentation/xpc/xpc_connection_set_peer_code_signing_requirement(_:_:))
  is the macOS 12 string-requirement API selected for the low-level harness;
  [`xpc_connection_set_peer_requirement`](https://developer.apple.com/documentation/xpc/xpc_connection_set_peer_requirement)
  is the newer macOS 26 composed-requirement form; and
  [`xpc_listener_set_peer_requirement`](https://developer.apple.com/documentation/xpc/xpc_listener_set_peer_requirement)
  define system-enforced peer filtering and listener-side request dropping.
- [`xpc_listener_set_peer_code_signing_requirement`](https://developer.apple.com/documentation/xpc/xpc_listener_set_peer_code_signing_requirement(_:_:))
  exposes the string-requirement form on the newer listener API.
- [`SecCodeCreateWithXPCMessage`](https://developer.apple.com/documentation/security/seccodecreatewithxpcmessage(_:_:_:))
  derives a running-code object from a received XPC message.
- [`SecCodeCheckValidity`](https://developer.apple.com/documentation/security/seccodecheckvalidity(_:_:_:))
  performs dynamic validation and applies an explicit code requirement.
- [`SecCodeCopySigningInformation`](https://developer.apple.com/documentation/security/seccodecopysigninginformation(_:_:_:))
  documents that validity must be checked before relying on returned signing information.
- [`xpc_connection_get_asid`](https://developer.apple.com/documentation/xpc/xpc_connection_get_asid(_:))
  returns the remote peer's audit-session identifier at connection time.
- [`xpc_connection_cancel`](https://developer.apple.com/documentation/xpc/xpc_connection_cancel(_:))
  documents asynchronous, non-preemptive cancellation.
- Apple's [XPC updates](https://developer.apple.com/documentation/updates/xpc) and
  [code-requirement guidance](https://developer.apple.com/documentation/security/applying-code-requirements)
  describe the supported requirement families and the verifier-owned policy model.

Local declaration readback:

| Input | Observed value |
| --- | --- |
| SDK | macOS 26.5 |
| Xcode | 26.6 (`17F113`) |
| Legacy low-level connection API | `xpc_connection_create_mach_service`, peer requirement, EUID, ASID, asynchronous reply, cancellation |
| New listener/session API | listener/session creation from macOS 14.0; string peer requirements from 14.4; composed `xpc_peer_requirement_t` from 26.0 |
| Security API | `SecCodeCreateWithXPCMessage`, `SecCodeCheckValidity`, dynamic signing information/status |

SDK availability is an observation about the installed development environment, not a frozen
minimum-OS claim.

## Documented behavior

### Peer requirement and delivery

The low-level `xpc_connection_t` requirement functions are supported on macOS and state that all
received messages are checked. For a listener connection, requests that fail the requirement are
dropped. When the client expected a reply, its side receives the peer-code-signing-requirement
error instead of an application reply.

The newer `xpc_listener_t` requirement functions likewise require an inactive listener and state
that nonmatching requests are dropped. The macOS 26 `xpc_peer_requirement_t` declaration also says
that peer sessions created from a listener do not inherit the listener requirement; a consumer
using that API must not assume the accepted session independently carries the same constraint.

Only one member of each `xpc_*_set_peer_*_requirement` family may be installed on an object. A
complex policy must therefore be expressed as one reviewed requirement rather than a sequence of
setter calls whose conjunction is assumed.

### Message-derived code identity

`SecCodeCreateWithXPCMessage` creates a `SecCodeRef` for the sender represented by the exact
received message. `SecCodeCheckValidity` then performs dynamic validity checking and can apply the
same explicit requirement used at listener admission. Signing information is not authority until
that validity call succeeds.

Dynamic signing information can expose the signing identifier, Team identifier, unique code
identity/CDHash family, entitlements, and a dynamic status snapshot. The debugged status bit is
sticky after attachment, but Apple warns that the status dictionary is a snapshot; the harness
must record the observation without turning it into a timeless attestation.

### User and session

`xpc_connection_get_euid` and `xpc_connection_get_asid` return the remote EUID and audit-session ID
at connection establishment. They are connection observations, not message-body fields. PID and
service name remain diagnostic values and do not establish code authority.

The public session/listener declarations inspected in the macOS 26.5 SDK do not expose equivalent
EUID/ASID accessors on `xpc_session_t`. That is an absence in the inspected public SDK, not proof
that no other supported API can ever provide the information.

### Cancellation, interruption, and reply state

Low-level XPC cancellation is asynchronous and non-preemptive. It discards unsent messages and
causes awaiting reply handlers to receive an invalid-connection error, but it does not interrupt
an already-running event handler. A remote service exit before reply produces an interrupted-
connection error; failure before send produces an invalid-connection error.

Apple's connection declarations explicitly say transport cannot indicate that a sent message was
processed and that transaction lifetime must be tracked at the protocol layer. XPC supplies no
server-side durable deadline or commit oracle.

## Capsule harness direction

### Selected controlled-harness baseline

Use the low-level C `xpc_connection_t` Mach-service listener for the first S3 harness:

1. it supports a peer code-signing requirement before listener activation;
2. its listener-side documentation defines nonmatching request dropping;
3. accepted peer connections expose public connection-time EUID and ASID observations; and
4. each delivered dictionary can still be revalidated through
   `SecCodeCreateWithXPCMessage` and `SecCodeCheckValidity`.

This is a harness/API selection within Accepted ADR-0029's native-front topology, not a new
component or product implementation. The newer listener/session API remains viable only if a
future minimum-OS decision and supported peer/session-observation contract close all required
fields without private SPI.

### Required order

The controlled harness must implement and retain this order:

```text
construct exact role-specific requirement from trusted fixture
  -> create inactive role-specific Mach-service listener
  -> install exactly one peer code-signing requirement
  -> install event handler and target queue
  -> activate listener
  -> accept peer connection
  -> compare connection-time EUID and ASID with trusted expected values
  -> receive exact XPC dictionary
  -> derive SecCode from that exact message
  -> dynamically validate against the same requirement
  -> inspect validated identifier/Team/CDHash/entitlement/status observations
  -> acquire connection/global flow slot
  -> validate outer key/type/version/width/method/current-state fields
  -> copy bounded application data
  -> invoke only the method-specific test bridge
```

A peer-requirement failure must produce zero application delivery. Any failure through outer
validation must produce zero body copy and zero bridge call. Bridge pointer/length/version faults
terminate the disposable process and leave durable semantics to the controlled restart oracle.

### Deadline and response delivery

Capsule owns the method deadline with a bounded monotonic timer starting after admission. The
timer may cancel or suppress reply delivery, but it cannot preempt a running handler, roll back a
Go call, infer whether state committed, or release a durable mutation slot early. After dispatch,
disconnect, cancellation, timeout, interruption, and client death are response-delivery events
only. Store truth decides replay and recovery.

The harness should use asynchronous replies. The synchronous send API can block indefinitely and
must not be used as the server's deadline mechanism.

## Required controlled evidence

The S3 experiment must retain at least:

| Case family | Required observation |
| --- | --- |
| Wrong code requirement | zero server message delivery/body copy/bridge call; client receives only transport failure |
| Right family, wrong CDHash/role | refusal before body copy and zero bridge call |
| Wrong EUID or audit session | connection refusal; zero body copy and bridge call |
| Message-derived substitution | exact message sender fails the repeated dynamic requirement |
| Debugged or invalid code | fail closed before body copy |
| Missing/extra/wrong-type/cap+1 fields | frozen S3 refusal and zero bridge call |
| Flow exhaustion | fixed `CAPACITY` after outer admission and before body copy |
| Cancel before dispatch | zero bridge call |
| Cancel/deadline after dispatch | bridge result converges independently; reply may be absent |
| Peer or server death | fixed interrupted/invalid transport observation plus protocol-owned convergence |
| Bridge fault | whole disposable server terminates; no ordinary application refusal is invented |

The harness must deliberately test whether listener rejection and client reply behavior match the
header contract on the exact owned host. It must not promote a documented rule into product
evidence without that observation.

## Stop conditions

Stop and revise the topology before product implementation if:

- the listener cannot install the exact requirement before any application delivery;
- supported APIs cannot provide the required EUID/session and exact-message identity observations;
- meeting the contract requires private SPI, a generic dispatcher, or another authority process;
- cancellation can preempt an already-dispatched mutation in a way that breaks store convergence;
  or
- the exact minimum-OS candidate lacks the selected supported APIs.

## Consequence

R1 is `PASSED` for its primary-source and local-declaration research scope. C2 is no longer blocked
on R1 or CL1, but still requires separate authorization naming `Shrimpworks/capsule-experiments`,
the owned Mac, exact S3 fixture digest, disposable process/service names, and the defensive
no-product scope. Installed identity/update evidence and product IPC remain `BLOCKED`.

## Confidence and limitations

Confidence is high in the documented peer-filter, message-derived code, connection-time EUID/ASID,
and non-preemptive cancellation semantics. Confidence is medium in their complete composition
until the controlled harness observes exact delivery, reply, process-fault, and session behavior.
No minimum OS, listener implementation, code requirement bytes, product identity, or admission is
selected by this research.
