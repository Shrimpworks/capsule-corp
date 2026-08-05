# macOS installation I2A protected-root bootstrap decision and I2B fault plan

Date: 2026-08-04

```text
Work item: I2A protected Supervisor-root/bootstrap-owner architecture and contract
Status: PASSED
Scope: owner selection, authority split, passive request/record contract, descriptor-relative
  creation/open ordering, execution-disabled transition, and the exact I2B implementation/fault
  slices only; this slice created or used no installed Capsule app/profile/key/Keychain item,
  service/process, protected root/lock/store, runtime, backend, or guest
Evidence or reason: current Apple public interfaces support an embedded on-demand XPC service,
  per-user SMAppService LaunchAgent registration, App-Group-named XPC, peer requirements, App
  Sandbox containers, and launch constraints. The selected composition gives object creation to
  the only process with Supervisor-container authority while keeping installation-root signing in
  a separate user-presence-gated role. The closed fault plan preserves every I0/I1/G2 no-create,
  no-guest, and repair-required invariant.
Remaining work: production wrapper review,
  separately authorized test-only key/service/container mutations, installed denial/fault
  evidence, and descriptor-relative fixed-v1 composition remain I2B work. Product-store selection,
  product IPC, update/restore, runtime/backend/guest, and attempt activation remain outside I2A.
Next action: implement I2B2 unsigned installation-only construction; I2B1 request/record fixtures
  and field authority are retained in the passive contract, and I1B is already complete.
Parent status: installed I2 protected-root composition is BLOCKED on the named I2B2-I2B5 evidence;
  macOS installation remains IN_PROGRESS — TRENDING_GOOD.
```

Decision: [Proposed ADR-0038](adr/0038-select-one-shot-coordinator-supervisor-bootstrap.md).

## Defensive scope and invariant

Defensively validate only Capsule's protected Supervisor bootstrap using passive repository
fixtures, then separately authorized exact Apple-signed components and owned disposable per-user
installation state on the named test Mac. Do not access another identity, credential, user,
session, container, service, process, store, runtime, backend, guest, or network target.

The invariant is:

> A separately enrolled, on-demand Trust Coordinator uses the Capsule installation-root key to
> authorize one exact bootstrap. The authenticated Supervisor creates and observes the protected
> root, owner object, and disabled fixed-store genesis inside its own App Sandbox container. The
> Coordinator alone constructs and signs the final bootstrap record over those observations. No
> ordinary startup, caller path, visible app, daemon, updater, replacer, or helper can create or
> replace those objects.

`FakeBackend.CreatesGuest() == false` is mandatory throughout I2. Passing I2 never enables an
attempt and never admits a product store, runtime, backend, or guest.

## Selected composition

The selected installation-only process tree adds one eighth role after I1A's exact inactive
seven-role tree:

```text
Capsule.app (visible Broker/setup UI)
  ├─ SMAppService registers CapsuleSupervisor.app LaunchAgent
  └─ Contents/XPCServices/CapsuleTrustBootstrap.xpc
       one on-demand, separately signed Trust Coordinator process
       installation-root Keychain group; no Approval/Supervisor key or state-root access
                    │ dedicated bootstrap-only App-Group-named XPC
                    ▼
CapsuleSupervisor.app (per-user App-Sandboxed LaunchAgent)
  creates only its fixed private-container objects; owns ADR-0033 flock and store
```

Exact proposed identities for the I2B installed fixture are:

| Item | Identity |
| --- | --- |
| Trust Coordinator signing identifier | `com.capsulecorp.capsule.trust-bootstrap.v1` |
| Coordinator/Supervisor bootstrap App Group | `3DDR84M4JS.com.capsulecorp.capsule.bootstrap.v0` |
| bootstrap-only Mach service | `3DDR84M4JS.com.capsulecorp.capsule.bootstrap.v0.supervisor` |
| bootstrap XPC protocol | `capsule.supervisor-protected-root-bootstrap/xpc-v0` |
| fixed Supervisor state-root entry | `supervisor.state` |
| fixed entries below the state root | `supervisor.owner`, `supervisor.store`, `supervisor.bootstrap-request`, `supervisor.bootstrap-record` |
| fixed transaction entries | `CapsuleSupervisor.bootstrap-pending` below `Application Support`; `supervisor.state.staging.` plus 64 lowercase request-payload-digest hex characters and `supervisor.state.publish-intent` below `CapsuleSupervisor` |

The visible app may render setup state, call `SMAppService.register`, and relay fixed progress. It
must verify the installed bundle, register the Supervisor, and read back an enabled status before
invoking its private Coordinator. Approval-required or denied status stops without a key, request,
root, or store. The app is not a member of the bootstrap App Group or installation-root
Keychain group, never receives a private-key reference, and cannot supply a path, field override,
record, or generic bytes-to-sign request. The Coordinator and Supervisor are the only bootstrap
App Group members. No file, defaults domain, socket, key, or durable authority state is placed in
that group container. Capsule creates no group-container entry and ignores any platform-owned
metadata there; group membership is retained as a real residual IPC, shared-container, and
potential Keychain-group capability and must pass the negative I2B corpus.

The Coordinator executable remains embedded for later explicit trust ceremonies, but its process
is on-demand and exits after a bounded request. Because it needs the interactive user's Keychain
and LocalAuthentication services, I2B must test `JoinExistingSession=true`; the documented
`false` default is not assumed to work. It is not an `SMAppService`, login item, permanent
agent, privileged helper, root service, network client, updater, replacer, or general signing
oracle. Its private installation-root key is nonexportable where the admitted Security/Secure
Enclave profile supports the selected algorithm, and every request/record signature requires
fresh system-mediated user presence. Loss of that key means a new installation identity, not
silent restoration.

The Supervisor bootstrap listener is a third, setup-only service that refines but does not widen
ADR-0029's two ordinary services and four product calls. It accepts only the two fixed bootstrap
messages below from the exact Coordinator profile. Once the bootstrap record is committed, every
new prepare request returns a fixed already-enrolled refusal without state access or mutation. The
four ordinary daemon/Broker calls remain disabled and outside I2.

## Why the alternatives stop

| Candidate | Decision | Reason |
| --- | --- | --- |
| visible setup app creates the root/lock/store | `NO_GO` | It would need Supervisor-container or shared-state authority and would make a compromised Broker a state replacement path. User-granted foreign-container access is an elevated capability, not an installation design. |
| Supervisor invents its installation root and self-authorizes | `NO_GO` | It collapses installation trust-transition authority into the component that owns execution and makes the bootstrap record self-asserted rather than independently authorized. |
| daemon or Broker owns the installation-root key or signed record | `NO_GO` | Either path violates the existing key and execution-authority split. The visible app is UI/orchestration only. |
| updater, bundle replacer, package, or ordinary installer owns bootstrap | `NO_GO` | Release replacement is not local installation-root authority and must not reset or recreate protected state. |
| permanent privileged helper, root LaunchDaemon, or root service | `NO_GO` | The selected per-user private-container operation needs no host-root primitive; the added authority and recovery surface is unjustified. |
| caller-selected root/store path or normal-start create-if-missing | `NO_GO` | It permits path substitution and plausible empty-state recreation. All names are closed release fields and ordinary startup is open-only. |
| Broker-relayed unsigned request or Apple code identity alone | `NO_GO` | Peer identity authenticates a process, not Capsule installation-root authorization over exact bootstrap bytes. |
| on-demand Trust Coordinator authorizes Supervisor creation | selected | It is the smallest composition that keeps the installation-root private key outside Supervisor while letting only the Supervisor create inside its private container. |

The selection depends on installed evidence. Apple documents that `SMAppService.register()`
immediately bootstraps a per-user LaunchAgent and registers it on later logins, that sandboxed apps
receive a system-associated container, that App Groups support Mach/XPC between sandboxed apps but
also grant a shared container and potential Keychain group, and that XPC peer requirements can
drop messages from a nonmatching listener peer. Apple also documents private embedded XPC
services and launch/spawn constraints. Those facts support the candidate; they do not prove the
exact Capsule bundle/profile/container result. See Apple's
[`SMAppService.register()` documentation](https://developer.apple.com/documentation/servicemanagement/smappservice/register%28%29),
[App Sandbox file-access documentation](https://developer.apple.com/documentation/security/accessing-files-from-the-macos-app-sandbox),
[App Groups entitlement](https://developer.apple.com/documentation/bundleresources/entitlements/com.apple.security.application-groups),
[XPC code-signing-requirement API](https://developer.apple.com/documentation/foundation/nsxpcconnection/setcodesigningrequirement%28_%3A%29),
and [launch-constraint guide](https://developer.apple.com/documentation/security/applying-launch-environment-and-library-constraints).

The exact distinction between documented substrate and installed evidence, including the
App-Group residual authority and stale-Coordinator Keychain limit, is retained in the
[post-I1B platform research](MACOS_INSTALLATION_PLATFORM_RESEARCH.md).

## Closed bootstrap messages

Both signed objects use ADR-0019's object-specific deterministic-CBOR and COSE_Sign1 profile:
tagged embedded payload, protected ES256/content-type/key-ID headers, empty unprotected map, no
external AAD, no detached payload, and exact canonical-on-wire comparison. I2A selects the object
fields and bounds, not a production wrapper implementation or dependency admission.

### `SupervisorBootstrapRequestV0`

- object type: `capsule.supervisor-bootstrap-request`
- object version: `0`
- media type: `application/capsule.supervisor-bootstrap-request+cbor;v=0`
- signing purpose: `capsule.installation.bootstrap.request`
- audience: `capsule.execution-supervisor.bootstrap`
- payload maximum: 2,048 bytes; complete COSE envelope maximum: 4,096 bytes

The request contains exactly:

| Field | Width/bound | First authority |
| --- | --- | --- |
| object type/version/purpose/audience | fixed values above | release object profile; never a request override |
| installation ID | nonzero 16 bytes | Coordinator random source |
| installation-root public COSE key and key ID | P-256/ES256-only, key bytes at most 256; key ID is SHA-256 of exact public-key bytes | Coordinator key ceremony |
| Supervisor ID and expected effective UID | nonzero 16 bytes; nonroot UInt32 | exact I1B profile plus observed current login |
| Coordinator/Supervisor component-profile digest | 32 bytes | signed installed profile generated from I1B readback |
| installation-manifest candidate digest | 32 bytes | Coordinator over closed candidate bytes |
| trust epoch sequence and candidate digest | sequence exactly `1`; digest 32 bytes | Coordinator installation genesis |
| state-root class and closed root/entry names | fixed I2 values; each ASCII name 1..255 bytes with no slash, backslash, dot, control, or Unicode | release/bootstrap profile |
| owner/store mechanism and store-format identities | fixed ASCII, each at most 128 bytes | ADR-0033 plus selected I2B conformance-store profile |
| attempts-enabled and backend flags | exactly `false`; `FakeBackendCreatesGuest` exactly `false` | I2 hard boundary |
| request nonce | nonzero 32 bytes | Coordinator random source |
| issued-at and expires-at | UInt53 Unix seconds; `issuedAt < expiresAt`, window at most 300 seconds | Coordinator trusted-time observation |

The Coordinator constructs the payload and envelope. The Broker cannot submit a pre-encoded
payload, choose a field, or request an arbitrary signature. The Coordinator durably retains the
exact envelope before delivery because a Secure Enclave ECDSA retry is not assumed to reproduce
signature bytes.

### `SupervisorBootstrapObservationV0`

This is a bounded authenticated-XPC reply, not a signed portable authority object. It contains the
request payload and envelope digests; exact observed root, owner, store, and retained-request
facts; initial store-byte digest; fixed component-profile digest; and one closed result. It has no
path, key, file descriptor, endpoint, free-form diagnostic, or caller field and is at most 4,096
application-data bytes. The Coordinator accepts it only from the exact live Supervisor peer on the
same request connection, independently checks every release/request-owned value, and constructs
the record payload itself. It never signs Supervisor-supplied bytes wholesale.

### `SupervisorBootstrapRecordV0`

- object type: `capsule.supervisor-bootstrap-record`
- object version: `0`
- media type: `application/capsule.supervisor-bootstrap-record+cbor;v=0`
- signing purpose: `capsule.installation.bootstrap.record`
- audience: `capsule.execution-supervisor`
- payload maximum: 4,096 bytes; complete COSE envelope maximum: 6,144 bytes

The record contains exactly:

| Field | Binding and authority |
| --- | --- |
| request bindings | every stable request field above plus SHA-256 of exact request payload and exact retained request envelope |
| installation and epoch | installation ID, root public key/key ID, manifest candidate digest, epoch sequence exactly `1`, epoch candidate digest, component-profile digest, attempts disabled |
| Supervisor/container binding | Supervisor ID, expected UID, Supervisor-private App Sandbox class, fixed root name, state-root device/inode, type directory, and mode exactly `0700` |
| ADR-0033 owner binding | fixed `supervisor.owner`, same device as root, observed inode, regular-file type, expected UID, mode `0600`, link count one, and owner mechanism identity |
| fixed-store genesis binding | fixed `supervisor.store`, regular-file type, expected UID, mode `0600`, link count one, same device as root, selected I2B conformance-store format, exact initial byte length/digest, and explicit non-product/no-guest label; the mutable store inode is deliberately not enrolled |
| retained bootstrap bytes | fixed request/record entry names, request envelope length/digest, derived record-payload digest retained outside the payload, and exact record-envelope retention policy |
| transition | `protected-root-validated-disabled`, finalization Unix second, no attempt activation, no runtime/backend/guest admission |

The Coordinator exclusively encodes and signs the record using the installation-root private key.
The Supervisor verifies the exact canonical payload, root-key authorization, purpose, audience,
request digest, installation/epoch/component bindings, and every locally re-observed filesystem
fact before retention. The Supervisor may never replace the record with a self-signed projection.

Apple code signing authenticates the exact Coordinator and Supervisor processes, their
entitlements, and the installed handoff. It does **not** sign or authorize the request/record.
Capsule's installation-root signature authorizes those objects. The initial local trust genesis is
the conjunction of explicit user presence, the exact Apple-enrolled Coordinator/Supervisor
channel, create-only Coordinator/Supervisor Keychain records, and the installation-root signature;
none is a substitute for another.

## Nonce, replay, time, death, and retained-byte rules

- A request is first-admissible only while the Coordinator and Supervisor exact profiles are
  active, `issuedAt <= trustedNow < expiresAt`, no bootstrap anchor exists, and the fixed state
  root and pending journal are absent.
- Time limits exposure before durable admission; they are not the replay defense. Once the exact
  request envelope digest is durably journaled, expiry does not force abandonment or creation of a
  replacement request. Only exact replay may finish that transaction.
- Installation ID, root-key ID, Supervisor ID, epoch `1`, component-profile digest, and nonce form
  the semantic replay tuple. Same tuple plus different exact payload/envelope is a fixed replay
  refusal. Any different tuple while a pending or completed transaction exists is repair-required.
- Concurrent exact requests converge on the one journal/root transaction. Concurrent different
  requests both fail except for the one that durably won exclusive journal creation.
- The Coordinator retains its nonexportable root-key reference, exact request envelope, exact
  record envelope once signed, their payload/envelope digests, and the terminal transaction state
  in its Coordinator-only Keychain group. It exposes no private material.
- The Supervisor retains the exact request and record envelopes as mode-`0400`, single-link regular
  files in the protected root. It also retains the exact record envelope in one create-only,
  nonsynchronizable, device-only item in the epoch-1 Supervisor Keychain group and requires exact
  readback. The bootstrap App Group container retains no authority bytes after response.
- Ordinary startup reads the fixed Supervisor Keychain anchor first, verifies the record signature,
  then opens the fixed private-container hierarchy and requires the root copy to be byte-identical.
  It never discovers an installation ID, key, root, or path from a caller or shared container.
- Signature bytes are retained for exact response replay but are never semantic identity. Payload
  bytes and the closed replay tuple own meaning.
- Death before the pending journal commit has no durable bootstrap effect. Death after it permits
  only exact replay from the Coordinator. Death after root publication but before record commit
  leaves execution disabled and resumes only from the retained request/observation. Death after
  record commit or response loss returns the same retained envelope. No retry creates a second
  root, lock, store, key, installation ID, or epoch.

If the Coordinator's terminal ledger says an installation existed, loss of the Supervisor anchor,
root, owner, store, or record never re-enters first-install setup. Only a future explicit abandon-
identity/new-installation or forward-repair contract may proceed. No such contract is selected in
I2A.

## Descriptor-relative creation and open sequence

The only pathname consumed is the Supervisor's own container URL returned by the platform API in
the Supervisor process. It is never sent over IPC or retained as authority. The Supervisor opens
that container directory once with `O_DIRECTORY|O_NOFOLLOW|O_CLOEXEC`, then traverses and creates
only closed components relative to retained directory descriptors. It retains the fixed private
parent descriptor and later reopens the root name on every held-owner check so a post-open root
rename/replacement fences the process.

### One-time authorized creation

1. Authenticate the Coordinator peer before body copy; check exact Team `3DDR84M4JS`, signing ID,
   active CDHash/profile/entitlements, EUID, audit session, no debug state, bootstrap group,
   protocol, message, purpose, audience, and caps. Verify the signed request and user-presence
   result before filesystem creation.
2. Resolve the Supervisor container through the platform API. Open its container, `Library`, and
   `Application Support` descriptors no-follow. Require each existing ancestor to be a directory
   owned by the expected UID, not group/other-writable, on the same expected volume, and opened
   `O_DIRECTORY|O_NOFOLLOW|O_CLOEXEC`; never chmod or normalize a platform-owned ancestor.
3. Under `Application Support`, write the closed request-bound journal as
   `CapsuleSupervisor.bootstrap-pending.tmp`, a mode-`0600`, single-link regular file; `fsync`,
   close/reopen and byte-compare it, publish it with
   `renameatx_np(..., RENAME_EXCL)` to `CapsuleSupervisor.bootstrap-pending`, then `fsync` the
   `Application Support` directory. No other request may replace or adopt it.
4. Only after that journal commit, exclusively create fixed `CapsuleSupervisor` mode `0700` and
   the fixed digest-derived staging root mode `0700`; validate directory type, expected UID, exact
   mode, same device, and descriptor flags. An exact replay may resume only the journal-bound
   parent/staging shape; an unexpected name or fact is repair-required. Under the staging-root
   descriptor, exclusively create `supervisor.owner` with
   `O_CREAT|O_EXCL|O_NOFOLLOW|O_CLOEXEC`, mode `0600`; validate same device, regular type, UID,
   exact mode/link count, `fsync`, close, reopen, and revalidate.
5. Acquire `flock(LOCK_EX|LOCK_NB)` on that exact owner descriptor and retain it. From this point,
   bootstrap and ordinary composition share the ADR-0033 lifetime owner invariant.
6. Create the exact fixed-v1 no-guest genesis bytes through fixed `supervisor.store.tmp`, a
   mode-`0600` regular file; `fsync`, close, reopen/bound-read/full-validate/byte-compare, publish
   with `renameatx_np(..., RENAME_EXCL)` to `supervisor.store`, and `fsync` the staging root. Create
   and verify `supervisor.bootstrap-request` through fixed
   `supervisor.bootstrap-request.tmp` the same way, change its final mode to `0400`, sync it, and
   sync the root. Every temporary and final entry must be regular, expected-UID, single-link, and
   on the root device; no temporary name may survive the next checkpoint.
7. Revalidate every entry and the held lock, sync the staging root, create/sync the exact
   mode-`0600` `supervisor.state.publish-intent` in the private parent, and sync that parent. Publish
   the staging root with `renameatx_np(..., RENAME_EXCL)` to fixed `supervisor.state`, sync the
   private parent, reopen the fixed root no-follow, and require its device/inode to equal the held
   descriptor. Reopen the fixed owner name and require the same device/inode as the held lock
   before returning the observation.
8. Receive only the Coordinator-signed record for that observation. Verify it, write and reopen the
   mode-`0400` `supervisor.bootstrap-record` through fixed
   `supervisor.bootstrap-record.tmp`, sync/reopen/verify, exclusive no-replace rename, and
   root-sync; add the exact
   envelope to the create-only Supervisor epoch Keychain item, and read it back byte-for-byte.
   A pre-existing item succeeds only when every attribute and byte is exact.
9. Remove the now-complete publish-intent and pending-journal entries only after record-file and
   Keychain readback, then sync the private parent and `Application Support` directories. Retain
   the immutable request/record copies inside the root. Their removal cannot authorize recreation:
   the Coordinator terminal ledger and Supervisor anchor now prove prior completion.
10. Without releasing the owner descriptor, derive the closed ADR-0033 enrollment from the record,
    issue one fresh owner-session ID, and perform the descriptor-relative store open/recovery below.
    The resulting state is `protected-root-validated-disabled`, never ready for attempts.

Confirmed failure before a rename preserves the last complete predecessor. An indeterminate rename
or directory sync requires reopen. With a pending journal, exact replay may resume only the exact
staging or final-root world; after publish intent, absence of both is repair-required. Once the
Coordinator terminal ledger or Supervisor anchor exists, a missing root is always repair-required.
Unknown, symlinked,
linked, extra, wrongly owned, wrongly permissioned, or substituted entries are retained for bounded
reporting and never normalized or deleted by ordinary startup.

### Ordinary descriptor-relative open

1. Read and verify the exact fixed Supervisor Keychain bootstrap anchor; absence is not permission
   to create. Derive only typed values and closed names.
2. Resolve the Supervisor's own container with the platform API and open the fixed hierarchy and
   `supervisor.state` descriptor-relative with `O_DIRECTORY|O_NOFOLLOW|O_CLOEXEC`. Validate root
   UID, mode `0700`, directory type, device/inode enrollment, and fixed record byte equality.
3. Open `supervisor.owner` relative to the root with
   `O_RDONLY|O_NOFOLLOW|O_CLOEXEC`; validate regular type, UID, mode `0600`, link count one, device,
   inode, and descriptor flags; acquire nonblocking exclusive `flock`; then reopen both the root
   from its parent and the owner from the root and require exact identity.
4. Only after the held-owner recheck, open `supervisor.store` relative to the retained root with
   `O_RDWR|O_NOFOLLOW|O_CLOEXEC` and no create/truncate/append flag. Validate regular type, expected
   UID, mode `0600`, link count one, same device as root, bounded length, supported selected
   fixed-v1 conformance format, installation/Supervisor/epoch bindings, set digests, and cross-links.
   The store inode is not compared with bootstrap enrollment because atomic commits replace it.
5. Bind the store and coordinator to the one owner-session ID, perform sorted `AttemptID` recovery
   through `FakeBackend.CreatesGuest() == false`, and keep attempts disabled. A failed held-root or
   held-owner check permanently fences the process before another read or mutation.
6. Shutdown listeners/queues, lifecycle work, store handles, owner descriptor, root descriptor,
   and private-parent/container descriptors in that order.

No path, file descriptor, endpoint, root identity, entry name, store bytes, or owner-session ID is
accepted from the Broker, daemon, Coordinator request, environment, argv, defaults, App Group
container, or ordinary IPC caller after enrollment.

## Transition from I1 to protected-root state

I2 has one forward transition and no readiness edge:

```text
I1B signed inventory exact; attempts disabled; no Supervisor root
  -> Coordinator request retained; attempts disabled
  -> Supervisor root/owner/fixed-v1 genesis published; attempts disabled
  -> installation-root-signed record + epoch-1 anchor retained; attempts disabled
  -> owner/store full reopen and no-guest recovery clean
  -> protected-root-validated-disabled
```

`protected-root-validated-disabled` proves only the exact installed I2B root/record/owner/fixed-v1
composition on the tested host. It does not make the I0 `ready` state reachable. Product store,
Approval/Supervisor operational-key activation, four-call authenticated IPC, component acceptance,
runtime/profile/backend admission, and an explicit final attempt-enable transaction remain later
gates. Any missing or uncertain prerequisite is `repair-required` or `BLOCKED`, not success with a
warning.

## I2B installed denial and fault oracles

Every case records request/record payload and envelope digests, Coordinator and Supervisor Keychain
item metadata/digests, private-parent/root inventory and device/inode facts, lock state, store bytes
and digest, owner-session presence, recovery/fake call counts, service state, and response. Tests
use only owned disposable fixtures and never an existing Capsule or unrelated container/key.

| Threat/fault | Required oracle |
| --- | --- |
| baseline same-user opens or mutates Supervisor container | OS denial; zero request/record/store/fake effect. A platform profile that permits it is `NO_GO` for that exact profile. |
| same-user substitutes root path, symlink, hard link, rename, or replacement before open | root/record identity refusal; repair-required; zero lock/store/recovery/fake work |
| root entry changes after open | next held-owner check detects parent/root device-inode split; permanent fence before store mutation |
| lock link/rename/replacement | ADR-0033 identity refusal or post-open fence; no second enrolled owner/store world |
| store symlink, nonregular file, hard link, wrong UID/mode/device, rename, or replacement at open | repair-required after owner acquisition; original bytes unchanged; zero recovery/fake calls |
| stale same-Team Coordinator/Supervisor binary | spawn/peer/CDHash/profile/entitlement/epoch refusal; old epoch Keychain group cannot read current anchor/key; zero bootstrap/store effect. Any stale binary with undetected current-root mutation is a blocker. |
| unsigned, wrong-Team, wrong-role, wrong group/service/session/EUID peer | XPC drop or fixed authentication refusal before body copy and all state |
| debugger/task-port or `get-task-allow` | shipping profile launch/peer refusal and attempts disabled. A successful task port is elevated posture; mutation must still become repair-required, but no general containment claim is made. |
| explicit user grant to another app's container | classify the actor as elevated user-granted; never request the grant in product; any mutation is detected and fenced, not claimed prevented |
| wrong purpose/audience/install/epoch/Supervisor/key ID/nonce/time | fixed binding, authentication, stale, or replay refusal; zero filesystem/Keychain/store/fake effect |
| exact request replay before journal | one journal wins; peers converge on one observation |
| exact replay after observation or record response loss | byte-identical retained observation/record result; no second signature, root, lock, store, installation, or epoch |
| same semantic tuple with different payload/envelope | replay refusal; retained bytes unchanged |
| Coordinator death before/after request or record signing | before retained envelope: no delivery; after retention: exact replay only; fresh user presence is required if the signing context died before record signature |
| Supervisor death before journal, during each file sync/rename, after root publish, before/after record file, Keychain add/readback, or response | pre-commit, exact resumable incomplete transaction, complete record, or repair-required according to the sequence above; never automatic recreation |
| launchd restart/backoff or SMAppService response loss | registration status is re-read; Supervisor reopens truth; no setup completion inferred from service status |
| screen lock, Keychain lock, logout/login, fast-user switch, reboot, sleep/wake | unavailable key/user presence blocks signing; ordinary Supervisor reopen binds the same EUID/audit session/root/record; wrong session has zero effect |
| manual whole-bundle replacement with exact I2 tuple | root/owner/record identity preserved, stale services stopped, current exact profile accepted, attempts still disabled |
| mixed or changed Team, CDHash, entitlements, group, service, constraint, record profile, root/lock, store format, or epoch | repair-required; no old-or-new compatibility window and no automatic migration |
| missing/corrupt Keychain anchor, request/record copy, root, owner, store, or cross-link after completion | repair-required; no create, rewrite, delete, fallback, epoch reset, or recovery/fake call |
| coherent old record/store/root restoration | rollback-uncertain and attempts disabled unless a separately protected latest checkpoint proves equality; I2 makes no rollback-prevention claim |
| any I2 construction | `FakeBackend.CreatesGuest() == false`; zero runtime/backend/guest creation and no ordinary success classification |

## Exact next I2B slices

### I2B1: passive objects and field authority

- Status: `PASSED` for the request/record scope in the
  [passive I2B1 contract](protocol/MACOS_INSTALLATION_I2B_PASSIVE_BOOTSTRAP_CONTRACT.md). The
  observation remains a future authenticated-XPC method shape and was deliberately not added to
  this request/record-only slice.
- Retains closed request and record types/CDDL, exact bounds, media/purpose/audience profiles,
  generated known answers, cap-plus-one/cross-object/replay/time fixtures, and recursive
  field-authority entries.
- Use deterministic test keys only. Do not call Security, Keychain, LocalAuthentication,
  ServiceManagement, XPC, filesystem creation, or a process.
- Exit evidence: independent Go and real Swift fixtures agree on all 71 cases; ADR-0019 remains
  Proposed and no production wrapper is claimed.

### I2B2: unsigned installation-only bundle construction

- Extend the I1 tree with the exact Trust Coordinator XPC bundle, bootstrap-only service
  descriptor, App Group/Keychain/launch-constraint projections, and fixed no-create state names.
- Keep all profiles inactive, Coordinator request signing unavailable, Supervisor placeholder
  nonexecuting, and attempts disabled.
- Exit: deterministic two-directory construction/readback rejects missing, extra, mixed, stale, or
  widened role/service/entitlement/constraint bytes without launch.

### I2B3: separately authorized Apple Development bootstrap handoff

- Requires completed I1B, exact Team-`3DDR84M4JS` profiles, explicit authorization for fresh
  test-only Coordinator/Supervisor Keychain groups and keys, App Group, `SMAppService`, installed
  container, and local process mutations.
- Run the two-message peer-authenticated ceremony with a fresh disposable installation ID and
  nonexportable test root where supported. No existing personal Capsule key/state may be used.
- Exit: exact signed-byte/profile/session/container readback plus request/record agreement; no
  product key, product store, ordinary IPC call, runtime, backend, or guest.

### I2B4: descriptor-relative fixed-v1/no-guest composition and fault matrix

- Replace G2's trusted absolute test path with the retained private-parent/root descriptor
  capability, root-name continuity check, descriptor-relative owner and store open, same-session
  G2 recovery, and ordered shutdown.
- Run every fault/death/substitution/response-loss case above using the current fixed-v1 oracle and
  `FakeBackend.CreatesGuest() == false`.
- Exit: `protected-root-validated-disabled` on the exact host/profile only. The fixed store remains
  a conformance oracle and no product consumer activates.

### I2B5: session and manual-replacement continuity checkpoint

- Repeat restart/backoff, screen/Keychain lock, logout/login, fast-user switch where locally owned,
  reboot, stale binary, exact manual whole-bundle replacement, mixed update, missing/corrupt state,
  explicit foreign-container grant, and rollback-uncertain cases.
- Exit: retained old-or-exact-new or repair-required results, exact root/owner/record continuity,
  no Capsule-created bootstrap App Group container entries, no old-epoch Keychain cross-use, and
  attempts still disabled.

## Boundaries left outside I2A/I2B

This decision does not select the production Supervisor store, production signed-object wrapper,
complete signed installed corpus, archive F4B+, product four-call IPC, Approval/Supervisor evidence
key activation, update verifier, bundle replacer, backup/restore, repair mutation, coherent rollback
anchor, attempt activation, runtime, backend, content, evidence, or guest. Product-store selection
must consume this root/owner/open contract rather than change it. Any store-format transition is a
separately authorized forward epoch and never a normal-start create fallback.
