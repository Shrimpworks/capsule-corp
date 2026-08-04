# ADR-0032: Select an enrolled TypeScript Source Preparer and immutable source store

- Status: Proposed
- Date: 2026-08-03
- Refines if accepted: ADR-0010, ADR-0013, ADR-0018, ADR-0026, ADR-0029, and ADR-0030

## Context

ADR-0026 fixes the ordering rule for the first TypeScript profile: exact Node 22.22.1/Amaro 1.1.5
strip-only emission finishes before plan construction, registration, Broker rendering, approval,
or execution. ADR-0030 defines the original-authoring, executable-JavaScript, profile, normalized-
options, transformation-record, and record-set objects, but deliberately leaves their production
owner and immutable store unresolved.

That owner cannot be chosen as a local implementation convenience. The component invoking a
parser over agent-controlled bytes can be compromised; the component retaining emitted bytes can
substitute or lose them; and any caller-supplied path, executable, option, or post-registration
transform request would break exact approved-byte authority. The owner also has to coexist with
ADR-0029's selected Supervisor topology without turning its native bridge into a generic helper or
freezing the legacy 562-byte `RegisterPlanV0` submitted-binding record.

The merged [authenticated-IPC S1 consistency stop](../AUTHENTICATED_LOCAL_IPC_S1_CONSISTENCY_STOP.md)
selects ADR-0030's versioned atomic-cutover path and makes an accepted owner/store topology its first
dependency. The merged
[field-authority manifest](../../schemas/authority/field-authority-manifest.json) classifies the
current passive approved-byte objects but deliberately omits the illustrative plan v1 and any
Source Preparer method/store fields. This ADR supplies the topology decision, not those later
contracts or classifications.

This ADR selects the smallest separate planning component that can own both transformation and the
role-separated immutable source set. It does not implement that component, accept ADR-0026 or
ADR-0030, freeze `ExecutionPlan` v1 or `RegisterPlanV1`, activate a consumer, admit a runtime, or
establish production process isolation.

The retained
[P0 authority/TCB checkpoint](../TYPESCRIPT_SOURCE_PREPARER_P0_AUTHORITY_REVIEW.md) adversarially
reviewed this topology before passive P1 contracts or bytes are frozen. Its decision is a bounded
**P1 HOLD**: the topology remains the least-dangerous viable TypeScript option, but it is not yet
implementable or evidenceable without the entry criteria and stop conditions recorded below. This
ADR therefore remains Proposed.

## Proposed decision

### One deliberately enrolled planning component

Capsule adds one unprivileged, per-user **TypeScript Source Preparer**. It is a method-specific
planning component, not an execution helper, Content/Approval Broker, parser service, runtime
adapter, or Supervisor extension. A trusted installer registers it as an embedded
`SMAppService.agent`; the daemon cannot install, replace, reconfigure, or enroll it.

One Source Preparer service process contains:

1. a small native C/Objective-C XPC/Security front end that authenticates exact enrolled peers and
   enforces fixed message, connection, deadline, and copied-buffer caps before Go dispatch;
2. a Go core that validates source-set semantics, owns the authoritative immutable source store,
   invokes the fixed transformer worker, constructs the exact ADR-0030 objects, and owns recovery;
   and
3. an exact, one-shot Node 22.22.1 worker child for one source set. The child runs only the
   release-owned strip-only bootstrap and bundled Amaro 1.1.5, receives source through bounded
   pipes, returns emitted bytes or a fixed refusal, and exits. It has no listener, store API,
   operational key, network entitlement, backend interface, content capability, or caller-selected
   file path.

The worker is an implementation detail inside this enrolled component, not a daemon-selected or
generic helper. The daemon talks only to the authenticated Source Preparer service. The Source
Preparer never talks to a backend or guest, and the Supervisor never asks it to transform.

The Source Preparer has no Installation-root, Approval, Supervisor-evidence, content, TUF, or
backend key. It cannot approve, register a plan, consume a grant, create an attempt, launch a guest,
clear trust state, release user-only content, or sign an enforcement claim.

It also has no installation-root or trust-snapshot signing key, Approval private key or keychain
access group, Supervisor-evidence private key, content-unwrapping key or content access group,
TUF-root/update signing key, runtime/backend credential, daemon authority database, or shared app
group. Its package, store, sockets, audit handles, and temporary objects must not be reachable
through a daemon/Broker/Supervisor access group. Absence is a signed-manifest and installed-artifact
property, not an assumption derived from a process name or UID.

### Exact package and transformer custody

The installed Source Preparer release binds all of the following through the active installation
manifest, trust epoch, and verified local trust snapshot:

- the native front end, Go core, fixed bootstrap, and source-store format;
- the official Node 22.22.1 macOS arm64 distribution archive with SHA-256
  `261da057fb25ff2912dd6abb7842fc915ddf7947a2cb3c8cce90875d2b9bb667`;
- the installed Node executable with SHA-256
  `245e0321af97d3c21dd4e7104457334dfe3c3ba7982d0db75363e354565f8cbb`;
- bundled Amaro exactly `1.1.5` and the official Node source archive with SHA-256
  `87104b07e7acee748bcc5391e1bc69cf3571caa0fdfb8b1d6b5fd3f9599b7849`;
- the exact ADR-0030 transformer-profile and normalized-options bytes; and
- the fixed worker argv, empty inherited environment, fixed working directory, descriptor map,
  sandbox/entitlement set, and bootstrap identity.

The updater/trust verifier may fetch and verify these release artifacts and reduce them into the
bounded trust snapshot. It never invokes the transformer or creates per-job objects. The installer
places the exact Node executable and bootstrap in the protected, non-user-writable component
bundle. Before every worker spawn, the Go owner opens and checks the release-owned files without
following links, verifies regular-file identity, length, digest, active release/epoch, and code
validity, and refuses any mismatch. `process.version` and `process.versions.amaro` must match the
profile inside the worker result frame; that observation supplements rather than replaces installed
artifact verification.

The fixed worker accepts no user argv, environment, module specifier, loader, import hook, config,
package, option, path, or output destination. Standard input and output carry one closed binary
frame; standard error is drained to a fixed 16,384-byte diagnostic sink and is never returned or
persisted as authoritative data. Only fixed refusal codes cross the parent boundary. The worker is
killed on protocol violation, output cap, cancellation, or deadline.

P0 does not yet accept the worker as less trusted than its parent. Until a platform-specific child
profile proves that the worker cannot open the source store or package tree, inherit authority,
spawn or retain descendants, load native code, use the network, attach to peers, or leave live
temporary state, compromise of Node/Amaro is compromise of the complete Source Preparer planning
and store TCB. Merely fixing argv, environment, cwd, and inherited descriptors is insufficient.

The exact source archive is retained as governed reproduction material; it is not parsed or built
on the live preparation path. A different Node, Amaro, platform, architecture, executable,
bootstrap contract, or transformation option is a new reviewed profile and release, not a fallback.

### Closed services and methods

The active installation manifest enrolls two role-specific Mach services:

| Service | Enrolled peer | Closed methods |
| --- | --- | --- |
| `com.capsulecorp.capsule.source-preparer.daemon.v0` | daemon | `PrepareSourceSetV0`, `AbandonPreparedSourceSetV0` |
| `com.capsulecorp.capsule.source-preparer.supervisor.v0` | Execution Supervisor | `ResolveSourceSetForRegistrationV1`, `AbortSourceRegistrationV1`, `RetainRegisteredSourceSetV1`, `ReadApprovalSourceSetV1`, `ReadExecutableSourceSetV1`, `ReleaseRegisteredSourceSetV1` |

There is no Broker listener, generic command, arbitrary object lookup, path lookup, transform-by-
digest method, execute-time transform method, store enumeration, or updater method. The Broker
receives a bounded copy through the Supervisor's registered-plan facade; it never gains direct
source-store authority.

Both listeners use ADR-0029's authentication order: listener requirement before activation,
message-derived `SecCode`, exact Team/channel/signing identifier/CDHash/entitlement and dynamic-
validity checks, expected effective user and audit session, fixed flow slot, closed outer shape,
then installation/epoch/audience/purpose binding before a copied method body reaches Go. The local
audience is exactly `capsule.source-preparer.local.v0`; every method has a distinct fixed purpose.
Caller-supplied role, purpose, executable identity, profile, or options never becomes trusted.

### Input, transformation, and commit

`PrepareSourceSetV0` accepts only an idempotency nonce and a copied, already strictly decoded but
still untrusted proposal source bundle:

- one through 32 entries;
- unique canonical ASCII logical paths, at most 256 bytes each and 64 bytes per segment;
- exact entrypoint membership;
- exact TypeScript or JavaScript media enum;
- at most 262,144 bytes per file and 1,048,576 aggregate original bytes; and
- no path, descriptor, URL, package, transformer option, profile choice, or emitted-byte authority.

The application-visible copied input budget is therefore exactly 1,048,576 source bytes plus at
most 8,192 logical-path bytes; fixed-width fields and XPC keys are checked separately before either
sum is copied. Cap-plus-one refuses without queueing, partial transformation, clamping, or a
committed source set. An SDK may calculate emitted bytes as a local hint, but v0 does not transmit
that hint. A future comparison method must remain untrusted and the Source Preparer must reproduce
the output itself.

The Go owner revalidates path, order, membership, media, UTF-8/BOM, file, aggregate, and syntax-
eligibility boundaries. JavaScript is copied unchanged. TypeScript enters the one-shot worker. The
parent, not the worker, checks every output length/media/UTF-8 boundary, independently hashes every
original and emitted file, and deterministically encodes the exact profile, options, manifests,
per-file records, and ordered record set. The worker never supplies a digest or successful record.

No successful object is visible until one source-set transaction writes every exact byte, syncs
every file and directory, atomically publishes the set, reopens it without links, and revalidates
the complete role graph. A response contains only a fresh `PreparedSourceSetID`, installation/
epoch/store-format binding, the three future plan digests, both aggregate lengths, and defensive
copies of the exact plan-source projection. The daemon uses that projection to construct a plan;
it cannot use the ID to retrieve bytes or mint another source set.

The idempotency nonce is scoped to daemon identity, installation, epoch, and the digest of the
complete request. Exact replay returns the same committed set; nonce reuse with different bytes
refuses. It exists only to make timeout/response-loss recovery determinate and grants no lookup or
registration authority.

### Immutable source-store topology

The Source Preparer is the only writer. Its protected private store has no shared app group and
contains:

- role-namespaced immutable blobs for exact original files, executable files, transformer profile,
  normalized options, original/executable manifests, records, and record sets;
- one immutable source-set descriptor binding installation, epoch, store/profile version,
  `PreparedSourceSetID`, ordered logical paths, entrypoint, every role/length/digest/media
  relationship, both aggregate lengths, and the proposal-request digest; and
- a bounded mutable lifecycle record whose states are `preparing`, `prepared`,
  `registration-pending`, `registered`, `release-pending`, `released`, or `quarantined`.

Physical deduplication, if implemented, is permitted only within the same nominal role and exact
`(digest, length, media)` tuple. Cross-role equal digests never share a lookup key. Logical path is
manifest data, never a store path or authority. Every store accessor reopens and rehashes the exact
role before returning a fresh copied buffer. No caller receives an mmap, mutable slice, file
descriptor, bookmark, URL, or host path.

“Only writer” and “protected private store” are security claims, not a pathname convention. A sole
role-namespaced store is accepted only as the ownership topology. It is rejected as a same-user
security boundary unless the P1 platform probe identifies one exact OS-enforced protected
container whose single enrolled member is the Source Preparer and proves negative access from the
daemon, Broker, Supervisor, updater, a stale Source Preparer, and an unrelated same-user process.
Mode bits, pathname secrecy, PID checks, code signing without a protected container, and same-Team
identity do not meet this requirement. If the probe cannot establish the boundary, TypeScript
preparation is NO-GO under this ADR. User planning direction accepts a bounded modern ESM
`.mjs`-only JavaScript first-release fallback, with no CommonJS, package resolution, legacy Node
module surface, or widening of the governed runtime contract. This ADR does not freeze that media/
profile decision; the exact contract and applicable ADRs must be updated before implementation.

The store accepts at most 32 prepared-but-unregistered sets, 256 registered nonterminal sets,
4,096 lifecycle/tombstone records, and 536,870,912 exact blob bytes. It never evicts a prepared or
registered set to admit another. Admission reserves the worst-case original plus emitted bytes and
object overhead before worker launch; the final object-specific CBOR maxima must be generated from
the accepted CDDL before implementation. Exhaustion returns `CAPACITY` and changes no authority.

A `prepared` set expires after ten minutes of trusted effective time and may be abandoned only
through the daemon method. A `registration-pending` or `registered` set has no time-based or
pressure eviction. Only an exact Supervisor registration-intent abort can return a pending set to
`prepared`; only an authenticated, idempotent Supervisor release after its durable registration/
approval/attempt/archive and source-retention policy permits release can decrement registered
references. Supervisor cohort eligibility under Proposed ADR-0031 does not by itself authorize
source deletion: any explicit reproduction/evidence retention hold must already be discharged or
transferred to a separately reviewed immutable source archive. Unknown or inconsistent references
enter `quarantined`; startup never recreates an empty store or treats absence as release authority.

P1 may freeze refusal, quarantine, tombstone, and reconciliation objects, but it must not freeze a
positive blob-reference decrement until the archive/retention rule selects exactly one durable
authority for release. Store genesis and migration likewise require a sealed, installer-owned,
installation/epoch/store-format-bound descriptor. Neither an ordinary first open, a daemon request,
nor a mutable updater request may authorize store creation, repair, migration, rollback, or
quarantine clearance.

### Registration, approval, and execution integration

`RegisterPlanV1` remains to be defined by ADR-0030's atomic cutover and the field-authority
manifest. It carries a submitted `PreparedSourceSetID` as a trusted-resolution input outside the
plan bytes. It does not extend or reinterpret `RegisterPlanV0` or its legacy 562-byte record.

Before calling the Source Preparer, the Supervisor durably records a fresh
`SourceRegistrationIntentID`, submitted `PreparedSourceSetID`, exact plan digest, installation, and
epoch. `ResolveSourceSetForRegistrationV1` atomically changes `prepared` to
`registration-pending`, binds that exact intent and plan digest, and returns bounded defensive
copies of all original, executable, profile, options, manifest, record, and record-set bytes. Exact
replay of the intent is a repeatable read; another intent or plan refuses. The Supervisor's
independent Go decoders then verify canonical bytes, nominal digest domains,
path/order/member/entrypoint equality, both aggregate lengths, every transformed-file record, every
JavaScript pass-through, profile/options/media/diagnostic/map/URL disposition, and the plan's three
digest roles. It does not rerun Node and does not claim semantic equivalence.

On validation refusal before registration commit, the Supervisor calls the idempotent
`AbortSourceRegistrationV1` with the exact retained intent; only that pending binding may return to
`prepared`. On success, the Supervisor durably commits the exact plan, registration, complete role
bindings, intent, and internal source-set reference in `source-retention-pending`; it then calls the
idempotent `RetainRegisteredSourceSetV1`, which binds the registration and changes the set to
`registered`. Approval fetch and attempt creation remain disabled until retention is durably
acknowledged. Startup resolves every durable intent from Supervisor store truth: committed
registration retries retain, while an uncommitted intent retries abort. Response loss never creates
a second plan registration, changes plan bytes, or leaves a pending set subject to expiry.

For approval, the Supervisor calls `ReadApprovalSourceSetV1` from its retained registration state
and returns the exact plan plus bounded copied source objects to the Broker. The Broker's independent
Swift decoders repeat the complete relationship checks and render the original/executable roles
and profile before signing. It trusts neither daemon prose nor a Source Preparer display string.

For each ordered logical path, approval rendering must distinguish original and executable media,
length, and digest; label pass-through versus strip-only transformation; show the exact profile,
Node/Amaro/bootstrap identity, normalized options, record-set identity, and diagnostics/source-map/
URL dispositions; and offer bounded inspection of both copied byte roles. It must warn that the
Source Preparer selected the executable JavaScript and that independent validators prove internal
consistency, not faithful TypeScript erasure or semantic equivalence. A Source Preparer-generated
summary is never approval text.

For execution, `ReadExecutableSourceSetV1` accepts only the Supervisor-authenticated registration
and committed attempt identities derived from retained state. It returns only the executable
manifest and executable file bytes. The Supervisor rehashes them against the registered plan before
staging through the backend's typed source transport. Original bytes, profile, options, records,
worker, `PreparedSourceSetID`, and every host path stay out of the guest and runtime. Attempt and
execute entry points remain registration-ID-only; no caller can request another transform or
provide replacement bytes.

### Deadlines, cancellation, flow control, and diagnostics

The Source Preparer admits at most four authenticated connections, one in-flight call per
connection, two concurrent transformations, and eight concurrent read/validation calls globally.
Each admitted transformation reserves 4,456,448 bytes: two independent 1,048,576-byte source
aggregates plus a 2,359,296-byte bounded framing/object/staging allowance. The global transformation
reservation is therefore 8,912,896 bytes. There is no unbounded queue; saturation refuses before
body copy or worker launch.

Those numbers remain illustrative until P1 derives request, response, CBOR object, staging,
diagnostic, per-connection, global reservation, and tombstone maxima from canonical bytes. Every
admitted connection and read must reserve its complete worst case before copying. No partial copy,
worker spawn, staging object, lifecycle transition, or nonce binding may survive a cap-plus-one,
reservation failure, or saturation refusal.

`PrepareSourceSetV0` has a ten-second end-to-end deadline, including a five-second worker deadline
and a five-second validation/commit budget. Read, retain, abandon, and release methods have a
two-second deadline. These are fail-closed protocol ceilings, not performance claims. Cancellation
before commit terminates the worker and discards verified staging; cancellation or response loss
after commit returns through idempotent replay. A timed-out caller never authorizes the service to
publish a second set.

Cancellation is response control after an authoritative transition: it cannot undo a committed
publish, registration intent, retain, release, or quarantine transition. The caller must recover
the exact result through the same idempotency or intent binding. Before a transition, cancellation
must kill the exact worker process tree, drain bounded stdio, close every mapping and descriptor,
unlink only transaction-owned staging, release reservations, and leave no authority state.

Replies and logs contain only opaque IDs, fixed component/build identifiers, byte/count metrics,
and closed codes such as `MALFORMED`, `UNSUPPORTED`, `BINDING`, `CAPACITY`, `TRUST_STATE`,
`TRANSFORM_REFUSED`, `TRANSFORM_TIMEOUT`, `LOCAL_FAILURE`, and `RECOVERY_REQUIRED`. They never
contain source, emitted bytes, logical or host paths, Node stderr, arbitrary diagnostics, approval
data, content handles, or guest strings.

### Crash, recovery, cleanup, update, and reproduction

The service holds one nonblocking installation-scoped owner lock. Startup opens the store without
creation, validates version/bounds/digests/cross-links/installation/epoch, removes only proven
unpublished staging owned by an exact durable preparation record, and otherwise quarantines
unknown state. A durable worker lease is committed before spawn and then bound to the observed
PID/start/code identity. After parent death, recovery terminates only that exact still-live worker,
never a PID/path match alone, and resumes or aborts from store truth. A crash before publication
has no visible set; a crash after publication returns the same set on idempotent replay.

The spawn protocol must not leave an unbound live-child interval: the parent reserves and records a
non-authoritative spawn nonce before launch, binds the returned process handle and immutable code
identity before source is sent, and refuses preparation until recovery can prove the exact child is
dead. PID alone is never identity. All pipes are nonblocking and drained under fixed caps so a
stderr or stdout flood cannot deadlock cleanup.

An identity, profile, Node, bootstrap, entitlement, or store-format change is a prepared trust-epoch
transition. Preparation and new registration stop; active attempts drain or are explicitly
reconciled; old connections close; store migration uses lock-held validate/write/sync/rename/sync/
reopen ordering; and all peers accept one exact new epoch before calls reopen. Old binaries refuse
new stores. Before epoch commit, only an explicitly reversible prepared transition may restore the
old world; afterward repair is forward-only. Existing registered sets remain exact executable-byte
authority across a compatible migration and need no re-transformation.

Old-epoch services and connections lose store access before the new epoch accepts calls. A stale
process may not retain an open store descriptor or mapping across transition. Incompatible
store/profile changes require drain or an explicitly reviewed dual-reader migration; silent
reinterpretation and rollback to an older writer are forbidden.

Every retained set contains enough exact bytes and identities for offline reproduction: original
and emitted files, all ADR-0030 objects, profile/options, Source Preparer release/bootstrap identity,
Node source and distribution target identities, and the request digest. A separate bounded offline
verifier may rerun the exact installed worker with no network and compare emitted/object bytes. It
does not modify the live store or establish semantic equivalence, toolchain stability, or an
execution approval. If retention policy later permits blob release, the tombstone must say that
offline reproduction is unavailable unless an independently validated immutable reproduction
archive was committed first; this ADR does not select that archive format or owner.

## Alternatives considered

### Embed Node or a transformer in the daemon

Rejected. The public daemon would gain the Node parser/runtime and immutable-store TCB, and a Go-to-
Node process shortcut would violate the daemon-to-helper prohibition. Linking an unsupported
`libnode` construction would no longer be the exact selected Node executable and would require a
different profile and evidence. Daemon compromise must not be able to fabricate a trusted source-
set record without the separately enrolled owner.

### SDK emission plus trusted validation

Rejected as the authority owner. SDK output is agent-controlled. Digest, media, path, and record
validation cannot prove that emitted bytes are the exact result of the selected Node operation.
Hints may improve ergonomics later only if the Source Preparer independently transforms and either
ignores or byte-compares them.

### Broker-owned preprocessing or storage

Rejected. Adding Node and agent-source parsing to the process holding Approval and user-content
authority increases the consequence of parser compromise and violates the Broker's narrow role.
The Broker independently reads and validates copied registered objects but owns neither emission
nor their source store.

### Supervisor-native or Go-owned transformation

Rejected. The Supervisor must remain the smallest execution authority and must not gain a public
source parser, Node worker, or pre-plan transform API. It resolves and validates exact stored roles
at registration but never emits them. ADR-0029's native bridge remains method-specific and contains
no transform method.

### Updater, generic parsing service, or runtime ownership

Rejected. The updater distributes verified release inputs but has no per-job authority. A generic
parsing service creates an open-ended request surface. Runtime or backend transformation occurs
after registration/approval and would make executed bytes depend on unapproved behavior.

### Store original bytes in the daemon and emitted bytes in the Supervisor

Rejected. Split custody creates partial commits and leaves no single owner able to return the exact
closed source set. It also makes Supervisor storage a source-ingest path and lets daemon state loss
change approval rendering. One role-namespaced store plus independent digest validation has a
smaller authority topology.

## Consequences, trust boundary, and blockers

- A compromised daemon can submit malicious original source and cause denial of service, but it
  cannot itself mint a Source Preparer source set, select transformer options, or substitute
  emitted bytes after registration.
- Compromise of the native front, Go core, store implementation, or—until separately proven
  confined—Node/Amaro worker can choose mutually consistent emitted JavaScript, profile/options,
  records, digests, and prepared-set claims before plan construction; retain or disclose
  agent-supplied proposal source; replay within any missing binding; exhaust bounded availability;
  or destroy its own unregistered state. It still cannot approve, register, launch, use operational
  keys, clear quarantine, read Broker-only content, or replace bytes after registration without a
  digest mismatch. The entire component is planning and approval-understanding TCB.
- Supervisor and Broker validators establish exact byte/digest/media/path/record consistency; they
  do not independently implement TypeScript erasure, prove toolchain behavior, or establish
  semantic equivalence. Deterministic independent reproduction is therefore not required for exact-
  byte approval, but approval must expose this trusted-not-proven consequence. If faithful erasure
  despite Source Preparer compromise becomes a required property, this topology is NO-GO and a new
  ADR is required.
- Source confidentiality is not added: proposal source is agent-facing data. This decision grants
  no access to Broker user-only inputs or outputs.
- No live host path, transformer capability, original TypeScript, or post-registration option
  reaches the runtime or backend.

P1 remains on hold until every entry criterion in the P0 checkpoint is satisfied without weakening
the same-user adversary or moving responsibility into another component. In particular, P1 needs
positive protected-container and child-confinement probe results, a closed genesis/update channel,
canonical Source Preparer objects and recursive field-authority verification, exact calculated
object/XPC/resource maxima, settled cancellation/death/cleanup semantics, and a selected retention
rule before positive release is frozen. Any failed platform boundary, required daemon/Broker/
Supervisor/runtime transformation, generic helper, live-path requirement, or unbounded/native-
loading worker design is a stop condition and requires JavaScript-only admission or a new ADR. The
JavaScript-only path is the separately reviewed modern ESM `.mjs`-only fallback above, not
permission to add CommonJS, package resolution, legacy Node module APIs, or wider runtime authority.

After those gates, acceptance remains blocked on a fault-injected fixed-store implementation,
governed Source Preparer/Node packaging, ad-hoc then Apple-signed installed identity and sandbox
evidence, update/rollback tests, independent Go/Swift validators, and ADR-0030's atomic plan/
registration/approval/lifecycle cutover. The exact implementation and falsifiable evidence plan is
[TypeScript Source Preparer implementation, conformance, and fault plan](../TYPESCRIPT_SOURCE_PREPARER_PLAN.md).
