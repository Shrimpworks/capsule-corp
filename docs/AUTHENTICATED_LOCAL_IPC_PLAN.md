# Authenticated local IPC implementation and conformance plan

Status: proposed architecture/conformance plan. It creates no product endpoint, production signer,
backend, runtime, process launch, or guest.

Owner: Execution Supervisor and macOS platform boundary.

Decision: [Proposed ADR-0029](adr/0029-select-authenticated-local-ipc-topology.md).

Protected-root bootstrap refinement:
[Proposed ADR-0038](adr/0038-select-one-shot-coordinator-supervisor-bootstrap.md) selects a
separate installation-only Coordinator/Supervisor service with two closed setup messages. That
service carries no ordinary product authority and becomes a fixed already-enrolled refusal after
genesis; the two services and four calls in this plan remain the disabled ordinary surface.

## Scope and safety boundary

Defensively validate Capsule's daemon/Broker/Supervisor role separation using only passive typed
fixtures, controlled local processes, the no-guest fake lifecycle, and an owned disposable macOS
installed-service harness. Do not access any unrelated system, identity, credential, user data,
workload, runtime, backend, guest, deployment, signing service, or third-party target. Do not add a
production endpoint or privileged helper.

The behavior oracle is the current `registrationstate`, `approvalattempt`, `lifecyclestate`, and
`registeredlifecycle` path. `SupervisorCore`, `RegistrationID`-keyed drive/recovery, daemon approval
forwarding, and execute-time plan replacement are forbidden.

## Current passive evidence

On 2026-08-03 the existing Gate E smoke probe was rerun from this repository without its optional
Apple-container service ping. The controlled probe created only ignored local build products and
called requirement parsing and dynamic validity on its own ad-hoc processes. It created no service,
key, endpoint, backend, runtime, content, or guest.

| Item | Observed value |
| --- | --- |
| macOS | 26.5.2 (25F84), arm64 |
| SDK | 26.5 |
| Swift | 6.3.3 |
| Apple clang | 21.0.0 |
| Go | 1.23.4 darwin/arm64 |
| Swift probe | valid requirement `0`; malformed requirement `22`; same-team setter `0`; self dynamic validity `0` |
| Go/C probe | valid requirement `0`; malformed requirement `22`; same-team setter `0`; self dynamic validity `0`; Go tests passed |

This reconfirms API availability and the feasibility of either native language binding on this
host. It does not authenticate a peer, inspect or use an Apple-issued identity, establish Team or
distribution enrollment, or validate production IPC. The retained Gate B installed-service and
Apple-credentialed results remain the evidence for those distinct observations.

## Dependency graph

```text
accepted ADR-0034 + passed M1 source/manifest fixtures
  -> M2/S1 passive plan-v0 registration/fetch fixture and facade cutover (PASSED)
       -> S1 passive RegisterPlanV0/GetRegisteredPlanV0 records/envelopes (PASSED)
       -> ADR-0044 passive SubmitMainMJSV0 product-adapter record/envelope/flow (PASSED)
       -> remaining SubmitApprovalV0/RequestAttemptV0 passive records/envelopes (BLOCKED)
       -> S3 native authentication/cap harness
            -> S4 single-process ad-hoc composition with fixed store + fake lifecycle
                 -> S5 installed Apple Development identity/session/update matrix
                      -> S6 Developer ID/notarized clean-host evidence
                           -> later consumer activation review

Existing independent blockers that S0-S6 do not close:
  ADR-0019 production wrappers/key authorization
  Supervisor archive/compaction + rollback/backup + installed ADR-0038 protected-root bootstrap
  Broker UI/user-presence/content custody
  runtime/backend/profile admission and evidence composition
```

ADR-0040 makes Source Validator R2-R5B future defense-in-depth rather than a first-release ordering
gate. The passive M2/S1 fixture/facade contract, follow-on two-method passive S1 records/logical
envelopes, and ADR-0044's exact CLI submission method plus aggregate flow contract are `PASSED`.
The other two logical methods and native/authenticated S3 remain `BLOCKED`. S5 requires S4 and
valid Apple Development identities. S6 requires final intended package bytes, Developer ID/
notarization authority, and clean disposable hosts.

## S0: decision review and invariant lock

Retain review sign-off on ADR-0029 before code is wired. Review must confirm:

- one unprivileged per-user Supervisor process and no helper;
- exactly two ordinary role-specific service names and four ordinary calls, separate from
  ADR-0038's setup-only service and two bootstrap messages;
- native authentication before application-body decode;
- method-specific, copy-only in-process bridge ownership;
- Go-only durable authority/lifecycle ownership;
- correlation-only request IDs and existing API idempotency;
- attempt creation before lifecycle effects; and
- startup enumeration and recovery by committed `AttemptID` only.

Exit evidence: accepted or explicitly revised ADR text. A topology or responsibility change stops
this plan and requires an updated Proposed ADR before implementation.

## M1 and S1/M2: passive contract boundary

Accepted ADR-0034 removes TypeScript and the Source Preparer from the first-release critical path.
The first active contract is one byte-exact pass-through `main.mjs` member under the existing
plan-v0 source role. The retained [S1 consistency stop](AUTHENTICATED_LOCAL_IPC_S1_CONSISTENCY_STOP.md)
still governs a conditional later TypeScript plan-v1 cutover; its 626-byte arithmetic remains no
layout, cap, or known answer.

M1 has atomically narrowed the passive proposal/source contract, generated the single-member
87..95-byte canonical source-manifest boundary, and adds recursive source-manifest field-authority
coverage. It replaces incompatible `.js`,
`.cjs`, `.ts`, `.mts`, and `.cts` accepts rather than adding a second active source profile.

ADR-0040 removes the role-specific Source Validator sequence from this first-release contract's
predecessor chain. R2-R5B remain future conditional defense-in-depth work. No validator result is
accepted as plan, registration, or fetch authority.

S1/M2 now adds no-product Go and TypeScript-readable fixtures for the finalized application data,
four role-specific operations, success replies, and fixed refusals. `RegisterPlanV0` atomically
carries exact plan bytes, the complete 562-byte plan-role projection, the exact source manifest,
and exact `main.mjs` bytes. `GetRegisteredPlanV0` returns defensive Supervisor-retained copies of
all four plus the registration. Candidate application-data maxima are 328,337 request bytes and
332,433 fetch-reply bytes; the fixture generator derives and verifies them from closed canonical
definitions before code treats them as caps. Native/Swift transport framing remains a later slice.

The fixture manifest must retain exact plan/registration/approval and complete-binding known
answers; protocol, service, role, tag, audience, purpose, request, installation, epoch, source,
length, deadline, and first-owner fields; exact maxima and cap-plus-one; closed classification/
reason and state/time/trust/core/store/adapter oracles; copy-ownership mutations; response-loss/
idempotency classifications; structural missing/extra/wrong-type and cross-object refusals; and
byte equality across implementations.

The follow-on [passive S1 contract](AUTHENTICATED_LOCAL_IPC_S1_PASSIVE_CONTRACT.md) freezes only
the two facade-backed method records and logical request/reply projections. It retains exact caps,
cap-plus-one, current installation/epoch, fixed refusal/no-state, copy-ownership, and response-loss
oracles with independent Go/Node agreement. At that slice it deliberately left XPC key spellings,
numeric message/status tags, transport encoding, and peer-authentication evidence unset. The later
[S3 native-contract prerequisite](AUTHENTICATED_LOCAL_IPC_S3_NATIVE_CONTRACT.md) now freezes the
first three items for exactly the three passive methods while leaving peer-authentication evidence
unset.

Accepted [ADR-0044](adr/0044-select-private-xpc-internal-alpha-cli-adapter.md) and the
[passive product-adapter result](INTERNAL_ALPHA_PRODUCT_ADAPTER_PASSIVE_CONTRACT.md) select exactly
one developer-signed internal-alpha private-XPC candidate, `SubmitMainMJSV0`, from the CLI to the
daemon. Its typed logical request carries only the exact strict-JSON `JobProposal`; its reply
carries only a Supervisor-issued `RegistrationID`. Generated fixtures and the unexported
in-process harness freeze cap/cap-plus-one, method-derived bindings, correlation-only request IDs,
zero-queue aggregate flow, deadline/cancellation/stall, zero-state refusal, and response-loss
semantics across submission and the existing registration/fetch plane. This is passive evidence:
the XPC encoding, listener, peer check, signing, installed inventory, daemon consumer, and protected
state remain `BLOCKED`, and diagnostic HTTP gains no mutation route.

There is one active plan-v0/source shape, no optional transform field, generic fallback, field
inference, or dual active v0/v1 acceptance. No fixture may import experiment code into product
packages or claim ADR-0019 acceptance.

A retained dependency-aware implementation plan for the `RegisterPlanV0`/`GetRegisteredPlanV0` S1/M2
slice — field ownership, request/reply derivation methodology, cap-calculation methodology, native/
Go copy-ownership tests, the full refusal and response-loss/replay test matrices, field-authority
manifest additions, and an explicit list of every value blocked on M1 — is in
[M2/S1 conformance slice plan](M2_S1_CONFORMANCE_SLICE_PLAN.md).

The Source Validator launchers do not add an ordinary Supervisor call or change the two ordinary
Supervisor services and four methods in ADR-0029. ADR-0038's separate setup-only service also does
not widen this product surface. A source-validation result is never submitted to the Supervisor
and cannot become a registration field, cache, or substitute for Supervisor byte/manifest custody.

## S2: Go facade and store projection

Status: `PASSED` for the passive RegisterPlanV0/GetRegisteredPlanV0 facade and fixed-store
transaction oracle. Authenticated service activation, the other two methods, and native transport
remain `BLOCKED`. Conditional TypeScript, plan v1, and Source Preparer evidence are not
first-release dependencies. S2 does not treat the 626-byte TypeScript arithmetic as a record.

Implement an internal, unwired facade with one Go entry point per method. The facade does not
accept role or purpose from request bytes; each entry point constructs the existing fixed
`AuthenticatedCallContext` internally.

Work:

1. keep the finalized registration, approval-submission, and attempt-request entry points as thin
   method-specific calls to the existing components after method-owned admission;
2. add the intended read-only Broker fetch facade only after its typed projection and method version
   are explicitly reviewed, returning defensive copies of the retained exact plan, complete role
   bindings, wire registration, canonical source manifest, and exact `main.mjs` bytes only after
   Broker call context, active installation/epoch/trust, registration binding, source revalidation,
   and expiry checks;
3. treat the finalized complete-role record as submitted nominal identities, resolve every role
   through fixed injected local resolvers, and construct the trusted complete-role binding passed to
   registration; the first slice uses only retained fixed resolver fixtures and does not claim
   production policy/profile/trust/content resolution;
4. expose no whole-store snapshot, approval bytes, created-attempt resolver, lifecycle driver,
   backend handle, key operation, or repair method;
5. keep approval and attempt replay semantics in the authoritative store, not a transport cache;
   and
6. make every returned byte slice non-aliasing and every refusal carry only the fixed class/reason.

Focused tests:

- wrong facade/service/role/purpose cannot reach a component;
- Broker fetch of unknown, other-installation/epoch, expired, transition-fenced, corrupt, or
  role-binding-mismatched registration refuses without mutation;
- caller mutation after dispatch and accessor/reply mutation after return cannot change retained
  bytes;
- missing, unknown, wrong-domain, or locally mismatched submitted role identities refuse before
  registration state change;
- concurrent exact approval submissions converge on one `ApprovalID`;
- concurrent exact attempt requests converge on one `AttemptID`; and
- registration retry remains fresh-registration behavior.

Exit evidence: focused Go tests plus the existing registration, approval/attempt, lifecycle, v1
store, and conformance tests. No XPC or product consumer exists.

## S3: native authentication and cap harness

Status: `BLOCKED` on a separately authorized local harness. The exact passive native XPC
key/type/version/tag/status/reason contract for `SubmitMainMJSV0`, `RegisterPlanV0`, and
`GetRegisteredPlanV0` is `PASSED` as the
[S3 native-contract prerequisite](AUTHENTICATED_LOCAL_IPC_S3_NATIVE_CONTRACT.md). Source Validator
R2-R5B is not a first-release dependency under ADR-0040. Native parsing must consume that generated
contract and the completed logical Register/fetch fixtures; it must not reinterpret fixture JSON as
raw XPC framing or invent a conflicting method/record version.

Build a strictly local, no-product C/Objective-C XPC harness from fixed identities/messages. It may
use only ad-hoc signing unless the explicit S5 credentialed environment is active. It creates no
operational key, protected production store, backend, runtime, content, or guest.

The listener must install its peer requirement before activation and then prove, in order:

1. wrong exact peer receives no delivered request;
2. delivered message yields `SecCode` through `SecCodeCreateWithXPCMessage` and passes the same
   strict active requirement;
3. EUID and audit-session match the current enrolled harness session;
4. flow slot and outer fixed header are checked before body copy;
5. installation/epoch/service/tag/audience/purpose agree before bridge dispatch; and
6. only then is the bounded body copied to a test double that records call count and bytes.

Ad-hoc identity cases: exact CDHash, stale different CDHash, same identifier impostor, unsigned,
copied byte-identical executable, debug-entitled/debugged where locally possible, PID/name/path
substitution, and fresh process instance after service death. State explicitly that Team/channel
and distribution claims are not tested ad hoc.

Protocol cases: wrong role/service/session/install/epoch/purpose/audience, request ID reuse after
completion/reconnect, same reused request ID with a different body, malformed/cap-plus-one, extra
value, partial/interrupted XPC connection, response loss, four-connection and eight-global
saturation, per-connection second in-flight request, and two/five-second deadline boundaries.
Request ID reuse is a fresh transport call and never changes core replay semantics.

Exit evidence: retained source, environment, exact fixture digest, observed results, and a table
separating OS enforcement, protocol enforcement, harness mechanics, inference, and untested claims.

## S4: single-process composition and fault matrix

Compose the native front end and Go facade in one disposable Supervisor executable against a
temporary fixed v1 store and `FakeBackend.CreatesGuest() == false`. Use a sibling no-symlink
mode-0600 owner-lock fixture implementing Proposed ADR-0033's enrolled object, explicit BSD
`flock`, and lock-before-store contract. Do not add the executable to a product package or service
manifest.

Bridge rules under test:

- one function per method, no opcode;
- synchronous dedicated core queue;
- Go copies input before return and retains no native pointer;
- native caller supplies fixed-cap outputs and retains no Go pointer;
- expected errors return fixed numeric class/reason only; and
- invariant panic, short/oversize bridge output, pointer/length disagreement, or ABI-version
  mismatch kills the disposable Supervisor and leaves restart recovery to store truth.

Inject death/fault at each boundary:

1. before/after peer admission;
2. before/after flow-slot acquisition;
3. before/after body copy and bridge dispatch;
4. before/after trusted-time persistence;
5. before/after registration commit and reply;
6. before/after approval commit and reply;
7. before/after atomic consume/create and reply;
8. before/after lifecycle-record creation, each fake effect intent/effect/checkpoint, and reply; and
9. before/during startup owner lock, store open, recovery enumeration, and each `Recover(AttemptID)`.

For every refusal/fault, assert exact store bytes/digests and counts, time-high-water/trust-fence
exceptions, approval state, attempt count, lifecycle/effect count, adapter call count, reply or
interruption, and restart result. No refusal may create/widen authority. Once consumed, an approval
never becomes usable. No fake effect occurs before committed attempt creation. A missing response
never proves abort.

Response-loss oracles:

- lost registration reply followed by retry creates another registration and retains the first;
- lost fetch reply is a repeatable read;
- lost approval reply followed by exact resubmission returns the same approval;
- lost attempt reply followed by exact replay returns the same attempt; and
- startup recovery enumerates committed attempts and calls only `Recover(AttemptID)`, even after
  registration/approval expiry.

Exit evidence: all faults converge to pre-commit, post-commit, or explicit recovery-required state;
no `RegistrationID`-keyed lifecycle path or `SupervisorCore` symbol exists.

## S5: Apple Development installed matrix

This slice requires deliberately authorized Apple Development identities and provisioning profiles
for the exact daemon, Broker, and Supervisor fixtures. Do not inspect or use unrelated identities.

Current non-secret resource discovery: Apple Membership Details confirms Team `3DDR84M4JS` and
shows that `W4QUR9FUL4` in the Apple Development common name is a member/display suffix, not the
Team ID. A new Apple Development identity SHA-1 `80A4...D3793` is locally present but not
authorized for use. Local signed/provisioned experiments still require exact role identifiers,
entitlements, and Team-3DDR profiles. Xcode 26.6 cached three 3DDR Gate B/wildcard profiles, but
their App IDs/entitlements do not match these roles and they are not reusable. A
separate Developer ID Application identity for Team `3DDR84M4JS` is later distribution
authority requiring explicit authorization and matching-Team package design. It must not be used
for current development work and does not make Developer ID or notarization work current.

Package the same single Supervisor executable as an embedded per-user `SMAppService.agent` and
rerun S3/S4 with:

- Team/channel/role identifier, Hardened Runtime, absent `get-task-allow`, relevant entitlement
  digest, and exact active CDHash set in both directions;
- distinct daemon/Broker service names and listener requirements;
- message-derived sender identity on every request/reply;
- correct and wrong EUID/session, simultaneous/fast-user-switch sessions where owned fixtures are
  available, screen lock, logout/login, launchd restart/backoff, and reconnect;
- old daemon/new Supervisor, new daemon/old Supervisor, stale Broker, mixed CDHash set, changed
  entitlement, prepared update, pending verification, epoch finalization, component acceptance,
  crash at swap/acceptance boundaries, and repair-required refusal;
- duplicate Supervisor startup and owner-lock contention;
- lock-object preservation across update plus enrolled device/inode mismatch refusal; and
- exact installed protected store and epoch-scoped operational-key access denial, without exercising
  an Approval private-key operation unless separately authorized by the approval plan.

Exit evidence: Apple Development installed-topology result only. It does not claim Developer ID,
notarization, Gatekeeper, clean-host, shipping, or production authentication.

A retained read-only provisioning and installed-test plan for this slice — proposed bundle
identifiers, required App IDs/profiles, minimum and prohibited entitlements per role, XPC/Mach
service identities, and an installation/rollback test matrix separating what is testable on the
current owned Mac from what requires a second or clean host — is in
[Apple Development provisioning plan](APPLE_DEVELOPMENT_PROVISIONING_PLAN.md).

## S6: Developer ID and clean-host admission evidence

This slice is blocked on the final intended package bytes, Developer ID identities, successful
notarization/stapling, and owned clean disposable hosts for every supported macOS floor.

Repeat the full negative matrix after install, update, rollback/forward repair, reboot, logout/
login, sleep/wake, pressure, locked Keychain, restore, and power interruption. Read back exact
Team/channel/role/CDHash/entitlements, service registration, protected storage, epoch-specific
groups/keys, Gatekeeper/quarantine/translocation state, owner lock, store checkpoint, and stale-peer
denials. An updater/package change invalidates the dependent observation.

Exit evidence may support only the exact tested distribution/platform matrix. Production consumer
activation still requires accepted production protocol/crypto/store/retention decisions and every
independent content/runtime/backend/evidence gate.

## Conformance matrix

| Threat/fault | Required oracle |
| --- | --- |
| wrong peer/role/service | OS drop or `AUTHENTICATION`; zero core calls/state/effects |
| wrong EUID/session | `AUTHENTICATION`; zero body copy/core calls/state/effects |
| wrong installation/epoch | `BINDING` or transition `TRUST_STATE`; zero authority/effects |
| wrong purpose/audience/tag | `AUTHENTICATION` or `UNSUPPORTED`; zero core calls/state/effects |
| malformed/short/cap+1/extra right | `MALFORMED`/`SCHEMA`; bounded allocation, zero state/effects |
| flow saturation | `CAPACITY`; no queue/body copy/state/effects |
| registration replay | fresh registration per successful core call; exact plan retained twice |
| approval replay/concurrency | one approval ID, current fixed state, no resurrection |
| attempt replay/concurrency | one consumed approval and one attempt ID |
| response loss | method-specific replay behavior from ADR-0029; no inferred abort |
| death before commit | pre-state or explicit indeterminate fence; never partial authority |
| death after commit | committed registration/approval/attempt remains authoritative |
| death after fake effect | durable intent/checkpoint drives same-effect reconciliation |
| stale/update race | services fenced until one active epoch and all acceptances agree |
| recovery after expiry | committed `AttemptID` recovered without authority revalidation |
| duplicate Supervisor | second owner refuses before store mutation or adapter call |
| corruption/version mismatch | repair-required; original state not recreated or rewritten |

## Verification at each retained slice

Run the repository-required verification from `AGENTS.md` plus focused native/Go tests. Record exact
tool versions and distinguish a skipped credentialed test from a pass.

```sh
pnpm install
pnpm check
pnpm lint
pnpm test
pnpm verify:schemas
pnpm verify:adrs
go test ./...
go vet ./...
go build ./...
golangci-lint run ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

The historical passive native smoke harness is archived at commit
`0d8233b55f153b27a901a9ec45a3834208e3aa86` under
[`experiments/gate-e-supervisor-topology/run-smoke.sh`](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-e-supervisor-topology/run-smoke.sh).
Replay it only from an explicitly checked-out archive worktree; it is not part of Capsule's normal
verification suite.

Also validate every relative Markdown link in the ADR, this plan, and any evidence report. S3-S6
must retain their own exact local harness commands; no credentialed or installed test is implied by
the passive smoke command above.

## Claim boundary

The current retained Gate B/E/P0-4A evidence plus the passive smoke probe supports selecting an
implementable process/language topology. It does not establish authenticated production IPC.
Until S5/S6 and all activation blockers pass, documents and status must say `proposed`,
`development-only`, `ad-hoc exact-CDHash observed`, or `Apple Development installed observation`
as applicable—never `production authenticated`, `shipping validated`, `secure`, `attested`, or
`production-ready`.
