# ADR-0029: Select one native-fronted Go Supervisor process for authenticated local IPC

- Status: Proposed
- Date: 2026-08-03
- First-release source-contract refinement: ADR-0034 on 2026-08-03
- Refines if accepted: ADR-0005, ADR-0011, ADR-0012, ADR-0013, ADR-0018, ADR-0019,
  ADR-0023, ADR-0024, and ADR-0025

## Context

Capsule has durable, unwired Go mechanics for exact plan registration, approval submission,
atomic approval consumption and attempt creation, and an `AttemptID`-keyed no-guest fake lifecycle.
It does not have an authenticated process boundary. Tests currently inject
`AuthenticatedCallContext`; this proves role checks inside the Go components but not the identity of
a live caller.

Retained macOS evidence establishes a narrower platform boundary than a new privileged helper:

- unprivileged per-user launchd services can expose separate Mach services for daemon and Broker
  callers;
- XPC listener requirements can reject wrong-role and stale exact-build peers before message
  delivery;
- `SecCodeCreateWithXPCMessage` can bind a received message to its actual audit-token sender;
- effective UID and audit-session observations compose with exact code requirements;
- both Swift and Go through a small native binding can call the required XPC and Security APIs; and
- the no-host-root installed topology can register and restart an embedded `SMAppService` agent,
  although its ad-hoc evidence is not distribution-signing evidence.

The same evidence also rules out several shortcuts. PID, pathname, process name, EUID, service
name, or a designated requirement alone is not exact enrolled identity. A stable Keychain access
group is not an exact-build or trust-epoch boundary. A separate helper that owns lifecycle
operations would add another enrolled authority, update seam, request protocol, response-loss
case, and recovery owner. The current backend direction has a narrow C ABI and supplies no evidence
that host root or another process is necessary.

This ADR selects the narrowest implementation topology for the four current local authority calls.
It does not activate an endpoint, freeze a public API, implement production signing or approval,
admit a backend/runtime, create a guest, or accept ADR-0019, ADR-0023, ADR-0024, or ADR-0025.

## Proposed decision

### One process and one crash/recovery boundary

The macOS v0 candidate is one unprivileged per-user Execution Supervisor process registered as an
embedded `SMAppService.agent`. A trusted installer or containing application owns registration and
update ceremony; the daemon cannot register, replace, kickstart as authority, or reconfigure the
service.

One enrolled executable contains:

1. a small native C/Objective-C front end for XPC, audit-token/session inspection, Security.framework
   code validation, fixed connection accounting, and application-visible message caps; and
2. the existing Go registration, approval/attempt, and lifecycle core, linked into the same
   executable through a versioned in-process C ABI.

There is no Swift XPC service, backend helper, LaunchDaemon, root process, daemon-to-helper route,
or daemon-to-backend route. Swift remains the preferred Broker language. The Supervisor front end
does not need Swift UI, LocalAuthentication, or the former Swift Containerization API, and the
retained C XPC/Security probes already exercise the required identity primitives. If a future
platform API can be used only from Swift, it may replace the in-process native shim without
creating another process or authority; a separate process remains a new decision.

The Go core is the only owner of the authoritative Supervisor store, exact registered plan bytes,
approval ledger, created attempts, lifecycle records, effect checkpoints, and recovery set. The
native front end owns only live XPC/Security objects, peer observations, bounded copied call
buffers, and connection/deadline accounting. It cannot manufacture a registration, approval,
attempt, lifecycle record, backend descriptor, or success classification.

Any unexpected native fault, Go invariant panic, or bridge contract violation terminates the whole
Supervisor. It is not converted to an ordinary protocol refusal. Process death releases the
installation owner lock; the restarted process validates the same store and completes
`AttemptID`-only recovery before authority-changing methods reopen.

### Enrolled services and peers

The candidate has exactly two role-specific Mach services, recorded verbatim in the active
`InstallationManifest`:

| Service | Enrolled peer | Closed calls |
| --- | --- | --- |
| `com.capsulecorp.capsule.supervisor.daemon.v0` | daemon | `RegisterPlanV0`, `RequestAttemptV0` |
| `com.capsulecorp.capsule.supervisor.broker.v0` | Approval Broker | `GetRegisteredPlanV0`, `SubmitApprovalV0` |

These names are versioned internal candidates, not public discovery names. There is no listener
whose peer requirement is a daemon-or-Broker disjunction and no generic command method. The
Supervisor exposes no IPC operation for `Drive`, `Recover`, backend configuration, trust reset,
quarantine clearing, approval invalidation, archive deletion, evidence signing, key use, or service
administration.

Before activation, the Supervisor constructs each listener requirement only from the verified
active manifest. The enrolled peer record contains:

- requirement/profile version and exact service name;
- Apple certificate channel and anchor, Team ID, and role-specific signing identifier;
- the accepted CodeDirectory hash set for every permitted architecture/versioned equivalent;
- the exact relevant entitlement digest, Hardened Runtime/library-validation requirements,
  `get-task-allow` absence, and any accepted launch/responsible-process constraint;
- expected effective UID and Aqua audit-session policy;
- installation ID, epoch sequence and digest, distribution identity, and transition state; and
- the fixed local audience and per-method purposes below.

The designated requirement or Team/signing identifier identifies a release family; it never
replaces the active CDHash set. A byte-identical copy is the same code identity. PID is recorded
only for connection accounting/diagnostics, never as authority.

Connection and message admission occurs in this order:

1. install the exact listener peer requirement before listener activation, so a mismatching peer's
   message is not delivered to application code;
2. require XPC-observed EUID and audit-session ID to match the enrolled current login session;
3. for every delivered request, derive `SecCode` from that exact XPC message, perform strict dynamic
   validity against the same active requirement, require the expected Team/channel/identifier,
   accepted CDHash, entitlement digest and runtime constraints, and reject debugged state;
4. acquire the fixed connection/global flow-control slot;
5. inspect the closed XPC dictionary key count, value types, version, fixed-width header fields,
   and application-visible data lengths without decoding the method body;
6. require the request installation ID and epoch sequence/digest to match the already-open
   Supervisor state; and
7. dispatch only the method implied by the role-specific service and closed message tag, then copy
   its bounded body into the method-specific bridge.

The derived call, never a caller-supplied body field, supplies the existing Go context:

| IPC call | Go operation/context |
| --- | --- |
| `RegisterPlanV0` | `Component.RegisterPlan`, daemon, `register-plan` |
| `GetRegisteredPlanV0` | new read-only Broker facade over the retained registration record, Broker, `get-registered-plan` |
| `SubmitApprovalV0` | `ApprovalAttemptComponent.SubmitApproval`, Broker, `submit-approval` |
| `RequestAttemptV0` | `ApprovalAttemptComponent.RequestAttempt`, daemon, `request-attempt` |

The local audience for all four calls is exactly
`capsule.execution-supervisor.local.v0`. Their anti-confusion purposes are exactly
`capsule.ipc.register-plan.v0`, `capsule.ipc.get-registered-plan.v0`,
`capsule.ipc.submit-approval.v0`, and `capsule.ipc.request-attempt.v0`. Service, peer role, message
tag, audience, and purpose must all agree. The signed approval still separately requires
`capsule.plan.approve` for audience `capsule.execution-supervisor`; local channel purpose never
substitutes for signed-object purpose or audience.

An identity-changing update disables both services for authority changes, closes old connections,
and keeps attempts disabled through `pending-verification` and component acceptance. No active
epoch contains an implicit old-or-new peer requirement. A deliberately enrolled multi-CDHash set
is one reviewed epoch state; otherwise stale, mixed, absent, restored, or transition-fenced state
fails closed. Exact XPC denial does not revoke an old Keychain group, so operational keys remain
subject to the separate epoch-scoped group/key transition.

### Closed message set and byte ownership

Every request uses an XPC dictionary with an exact set of fixed ASCII keys, exact value types, and
these common fields: protocol version `0`, one nonzero 16-byte request ID, the 16-byte installation
ID, epoch sequence in `UInt53`, 32-byte epoch digest, the call's fixed numeric message tag, and the
call's exact fixed audience and purpose strings. The service and tag supply role and independently
select the only accepted audience/purpose pair; caller values cannot override them. Unknown or
missing keys, wrong types, unknown tags/versions, extra file descriptors, endpoints, Mach rights,
or nested dictionaries/arrays refuse before a Go call.

The request ID is correlation only. It is echoed in a reply and bounded diagnostics but is not a
registration, approval, attempt, nonce, idempotency, replay, authorization, or store key. In
particular, it does not weaken ADR-0023's rule that every successful `RegisterPlan` call creates a
fresh registration. Because a connection permits only one in-flight call, request ID uniqueness is
required only for that live call. Reuse after completion or reconnect is treated as a fresh
transport call and follows the method's semantic replay rules below; no bounded or durable request
ID history is consulted.

Method bodies and successful replies are exactly:

| Call | Request body owned by caller, copied before Go decode | Successful reply |
| --- | --- | --- |
| `RegisterPlanV0` | exact plan data, 1..65,536 bytes; complete typed role-binding record, exactly 562 bytes; exact canonical single-member source manifest, 87..95 bytes; exact `main.mjs` data, 0..262,144 bytes | exact Supervisor-issued `PlanRegistration`, at most 4,096 bytes |
| `GetRegisteredPlanV0` | nonzero `RegistrationID`, exactly 16 bytes | exact retained plan, at most 65,536 bytes; exact 562-byte binding record; exact retained registration, at most 4,096 bytes; exact 87..95-byte source manifest; exact 0..262,144-byte `main.mjs` data |
| `SubmitApprovalV0` | nonzero `RegistrationID`, 16 bytes; exact candidate envelope, 1..512 bytes | `ApprovalID`, 16 bytes; closed approval-state tag |
| `RequestAttemptV0` | nonzero `RegistrationID`, 16 bytes; nonzero `ApprovalID`, 16 bytes | `AttemptID`, 16 bytes; closed attempt-state tag |

The 562-byte submitted binding record is a bridge-only, fixed-layout v0 projection of the current
`ExecutionPlanRoleBindings`: one version byte; installation ID; epoch, source, inline-input, and
runtime-bundle digests; one review count byte; eight fixed review-digest slots; and profile
registry, backend validation, backend configuration, trust snapshot, and policy digests. Unused
review slots are zero-filled and rejected if nonzero. The review count is at most eight. This
record is not signed, stored as a replacement plan, or accepted from Broker approval. The
daemon's role labels are not trusted merely because it supplied the record. A method-specific Go
facade resolves every submitted identity through the Supervisor's locally trusted installation,
epoch, policy, profile, validation, trust-snapshot, source, and input resolvers and constructs the
trusted `ExecutionPlanRoleBindings` passed to `Component.RegisterPlan`. The Supervisor then
independently decodes the plan and requires every nominal role and exact locally resolved value to
match before state change. The current no-product slice may use only the fixed retained resolver
fixture and must not claim that production resolvers exist. Broker fetch returns the Supervisor-
retained resolved binding projection, not caller labels reinterpreted on read.

Accepted ADR-0034 freezes the first-release plan-v0 source role as one byte-exact pass-through
`main.mjs` member. Registration additionally validates and atomically retains the exact canonical
source manifest and member bytes; Broker fetch returns defensive copies of those same retained
bytes. This does not add a role to the 562-byte plan-binding projection. The source manifest and
member data are method-owned registration inputs, not replacement plan bytes or execute-time
authority. Registration commits none of plan, registration, bindings, manifest, or source unless
all validate and commit together.

The native layer owns each received XPC object and fixed-cap buffer. It copies no body until peer,
session, flow slot, outer shape, installation, and epoch checks pass. Each bridge call is
method-specific; there is no opcode entry point. Calls run synchronously on one dedicated core
queue. Go copies every input before returning and never retains a native pointer. Native code
allocates fixed-cap output buffers and never retains a Go pointer or receives a pointer to Go
memory. Expected refusals return a closed numeric status and fixed reason code. Oversize output,
short write, bridge version disagreement, or pointer/length inconsistency terminates the process as
a local integrity fault.

The exact application-visible aggregate request-data caps are 328,337 bytes for
`RegisterPlanV0`, 16 for `GetRegisteredPlanV0`, 528 for `SubmitApprovalV0`, and 32 for
`RequestAttemptV0`, excluding the fixed common header whose fields have the widths above. The
`GetRegisteredPlanV0` successful reply-data cap is 332,433 bytes. The revised registration/fetch
caps come from the complete ADR-0034 field-authority projection: 65,536 plan bytes, 562 binding
bytes, at most 95 canonical manifest bytes, 262,144 source bytes, and—on fetch—at most 4,096
registration bytes. XPC owns its transport framing; Capsule does not claim a raw Mach-message byte
cap it cannot observe. Passive fixtures must generate and verify the aggregate maxima from the
closed message definitions before implementation; the prose arithmetic alone is not a known
answer. A cap-plus-one data value is rejected before allocation growth or Go decode.

### Flow control, deadlines, cancellation, and partial delivery

Each service accepts at most four live authenticated connections, each connection has at most one
request in flight, and the process has at most eight admitted requests total. One dedicated core
queue serializes store-affecting bridge calls; there is no unbounded application queue. A slot miss
returns `CAPACITY` after authentication and outer-header validation and before body copy/decode or
state change.

Supervisor-owned method deadlines start at admission: two seconds for `GetRegisteredPlanV0` and
five seconds for each stateful call. Clients cannot extend them. Disconnect or cancellation before
bridge dispatch makes no Go call. After dispatch, client cancellation is only response
cancellation: it cannot decide whether a durable operation committed. The Supervisor completes the
current determinate store operation or enters the existing recovery fence for an indeterminate
outcome. A deadline never rewrites an expired registration/approval, rolls back a consumed grant,
or clears cleanup work.

XPC does not deliver a partial dictionary as an accepted request. Connection interruption before
dispatch has no state effect. Interruption after dispatch follows the same response-loss rules as
process death. No method interprets EOF, connection close, or absence of a reply as transaction
abort or success.

### Idempotency and response-loss recovery

The existing durable semantics remain the oracle:

- `RegisterPlanV0` is deliberately non-idempotent. If a reply is lost after commit, the daemon may
  submit the same exact plan and bindings in a new call; that successful call creates a new
  registration. The unreachable earlier registration remains retained and expires normally. The
  transport does not invent a caller-controlled deduplication key.
- `GetRegisteredPlanV0` is a read-only exact-byte lookup and may be repeated by registration ID.
- `SubmitApprovalV0` is idempotent by verified canonical approval payload and resolved signer
  authorization identity. Exact/equivalent replay returns the same `ApprovalID` and current state.
- `RequestAttemptV0` is idempotent by `(RegistrationID, ApprovalReference)`. Replay or concurrency
  after atomic consume/create returns the same committed `AttemptID` and never drives another
  attempt.

An indeterminate store commit returns or recovers as `RECOVERY_REQUIRED`; mutation stays fenced
until reopen validation. No XPC cache or request-ID table overrides store truth.

### Startup and recovery boundary

Startup follows ADR-0025 before normal call dispatch:

1. acquire the one installation-scoped nonblocking owner lock;
2. open the authoritative store without creation and validate version, bounds, digests,
   cross-links, installation, Supervisor identity, epoch, and trust state;
3. enumerate the store's sorted, duplicate-free committed recovery `AttemptID` set;
4. call only `Recover(AttemptID)` for each entry; and
5. keep all four calls fail-closed while recovery, quarantine, repair-required state, exhausted
   reconciliation, or unresolved cleanup prevents attempt enablement.

Proposed ADR-0033 supplies the exact owner-lock decision for step 1: an installer-enrolled
pre-created sibling object, exact descriptor identity checks, and lifetime nonblocking BSD
`flock`. The lock is acquired before store access and a busy contender performs no core or adapter
work. That design has local temporary-process evidence only; this ADR's installed identity/session/
update gates still require owner-required store/startup composition and the protected-state-root
matrix. Passive G1's internal Go/Darwin acquisition does not close those gates.

Recovery never calls `RegisterPlan`, re-submits approval bytes, asks whether an approval is still
usable, or rechecks registration/approval expiry for an already committed attempt. It never accepts
a registration ID, approval reference, plan bytes, role bindings, backend flags, image, mount,
path, or caller-selected recovery action. A missing lifecycle record remains work for the committed
`AttemptID`; a missing process or response is not proof of absence.

### Refusal classes and no-state rule

Wrong OS peer requirements are dropped by XPC and receive no application reply. Delivered calls use
only the current fixed internal classifications:

| First refusal boundary | Classification |
| --- | --- |
| wrong role/service, EUID/session, message-derived code identity, debug state, or local purpose/audience | `AUTHENTICATION` |
| malformed XPC shape/type/width, cap-plus-one body, zero identifier, or submitted-record length mismatch | `MALFORMED` or `SCHEMA` as owned by the current decoder |
| unknown protocol/message version or unavailable tag | `UNSUPPORTED` |
| installation, epoch, registration, plan, Supervisor, approval-reference, or role-binding mismatch | `BINDING` |
| expired registration/approval or non-idempotent replay condition | `STALE` or `REPLAY` as owned by ADR-0023/0024 |
| connection/in-flight/store ceiling | `CAPACITY` |
| transition fence, quarantine, repair-required, attempts disabled | `TRUST_STATE` |
| trusted local prerequisite or determinate bridge/store failure | `LOCAL_FAILURE` |
| indeterminate commit/cross-record outcome | `RECOVERY_REQUIRED` |

Every refusal before the Go call has zero durable effect. Existing stateful refusals may advance
only the durable time high-water or tighten a fail-closed trust/recovery state as already specified;
they create/widen no authority, consume no approval, create no attempt, and invoke no lifecycle or
backend effect. Diagnostics contain request ID, fixed role/service/method/classification codes,
epoch sequence, and enrolled build identifiers only—never plan, approval, content, path, guest, or
arbitrary caller text.

## Rejected alternatives

### Separate Swift XPC front end plus Go Supervisor service

Rejected. It adds an authenticated protocol and response-loss boundary between two processes, a
second enrolled update/epoch component, and ambiguity over store/lifecycle ownership. No required
API or privilege currently justifies that additional authority boundary.

### Rewrite the existing authority/lifecycle core in Swift

Rejected for this slice. Native Swift can implement the platform calls, but rewriting the durable
Go oracle would duplicate the attempt-keyed fault/recovery semantics without reducing process or
privilege count.

### Go process with no native adapter

Rejected as an inaccurate boundary description. The supported XPC, Security, audit-token, and
ServiceManagement primitives are native platform APIs. A small in-process native adapter is
explicitly part of the TCB even when called through cgo.

### Entitlement-bearing or privileged launcher helper

Rejected for v0. No host-root operation is evidenced, and a helper capable of create/start/stop/
destroy/reconcile is neither tiny nor stateless. The lead backend's C ABI and direct process
identity needs can remain in the Supervisor process. Any later exception requires a new ADR and a
sealed descriptor protocol; it may not receive plan, approval, content, paths, images, mounts,
flags, or public messages.

### One multiplexed service or generic NSXPC/Codable/JSON command envelope

Rejected. A broad listener requirement and generic decode surface create role/type confusion and
allocation before method-specific caps. Two services and four closed messages are smaller and map
directly to the current authority APIs.

### Unix-domain socket authenticated by path/mode/PID/EUID

Rejected. Same-user file modes, socket path, PID, and EUID do not establish enrolled component
identity, exact build, or current audit session.

### Durable request-ID deduplication for registration

Rejected. It would change ADR-0023's fresh-registration semantics and add another replay/retention
ledger. Approval and attempt operations already carry the semantic idempotency identities required
for safe response-loss recovery.

## Consequences and blockers

- One executable, owner lock, store, and restart path define the Supervisor authority boundary.
- Native code performs authentication before application-body decode; Go remains the only durable
  authority/lifecycle owner.
- The daemon can supply only exact plan bytes plus complete trusted-role bindings and can later
  request an attempt only by registration and approval reference.
- The Broker fetches the exact Supervisor-retained typed plan and submits approval directly; the
  daemon never forwards approval bytes or receives the Approval private key.
- No added process or helper receives backend, key, content, or recovery authority.
- Registration response loss may leave an expired retained registration; solving that operational
  cost requires changing ADR-0023 rather than hiding a deduplication key in transport.

This ADR remains Proposed until the implementation/conformance plan retains passive bridge
fixtures and fault oracles, an ad-hoc no-product two-service harness proves the ordering and
no-state refusals, and an Apple-signed installed package proves the exact identity/session/update
matrix. Consumer activation additionally remains blocked on production Supervisor
archive/compaction, implementation and installed evidence for the selected multi-process owner
lock, production approval verification and
key authorization, protected storage, update/repair integration, and all existing runtime/backend/
content/evidence gates.

[Proposed ADR-0031](0031-checkpoint-closed-supervisor-cohorts.md) defines the archive/compaction and
replay-retention semantics, but implements no storage behavior and deliberately does not select a
production engine. Its archive operations depend on this ADR's still-unimplemented lifetime owner
lock and expose no new IPC call.

Ad-hoc local evidence may support API availability, exact-CDHash denial, message shape/caps,
purpose/audience/epoch refusal, process death, replay, flow control, and bridge ownership. It cannot
support Team/channel enrollment, Developer ID/notarized distribution, provisioning-backed
entitlements/Keychain groups, App Sandbox/protected-container claims, wrong-user/session and
fast-user-switch behavior, installed update races, clean-host/minimum-OS support, or authenticated
production IPC.

## Implementation and conformance plan

The dependency-ordered slices, exact fixtures, fault matrix, evidence boundaries, and activation
gates are in [the authenticated local IPC implementation and conformance plan](../AUTHENTICATED_LOCAL_IPC_PLAN.md).

## Evidence

- [Gate B macOS authority-separation results](../../experiments/macos-authority-separation/RESULTS.md)
- [Gate B installed per-user services](../../experiments/gate-b-installed-services/RESULTS.md)
- [Gate E Supervisor topology](../../experiments/gate-e-supervisor-topology/README.md)
- [P0-4A installed development topology](../../experiments/gate-c-installed-development-topology/RESULTS.md)
- [Execution Supervisor current boundary](../EXECUTION_SUPERVISOR.md)
- [Installation trust](../security/INSTALLATION_TRUST.md)
- [Update and recovery](../UPDATE_AND_RECOVERY.md)
