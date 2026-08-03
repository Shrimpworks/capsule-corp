# Authenticated local IPC implementation and conformance plan

Status: proposed architecture/conformance plan. It creates no product endpoint, production signer,
backend, runtime, process launch, or guest.

Owner: Execution Supervisor and macOS platform boundary.

Decision: [Proposed ADR-0029](adr/0029-select-authenticated-local-ipc-topology.md).

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
S0 decision review
  -> S1 passive message/bridge fixtures
       -> S2 Go method-specific facade and Broker fetch projection
       -> S3 native authentication/cap harness
            -> S4 single-process ad-hoc composition with fixed store + fake lifecycle
                 -> S5 installed Apple Development identity/session/update matrix
                      -> S6 Developer ID/notarized clean-host evidence
                           -> later consumer activation review

Existing independent blockers that S0-S6 do not close:
  ADR-0019 production wrappers/key authorization
  Supervisor archive/compaction + rollback/backup + owner lock
  Broker UI/user-presence/content custody
  runtime/backend/profile admission and evidence composition
```

S1 fixtures may proceed in parallel by language after S0. S2 and S3 may proceed in parallel after
the shared fixture contract. S5 requires S4 and valid Apple Development identities. S6 requires
final intended package bytes, Developer ID/notarization authority, and clean disposable hosts.

## S0: decision review and invariant lock

Retain review sign-off on ADR-0029 before code is wired. Review must confirm:

- one unprivileged per-user Supervisor process and no helper;
- two role-specific service names and exactly four calls;
- native authentication before application-body decode;
- method-specific, copy-only in-process bridge ownership;
- Go-only durable authority/lifecycle ownership;
- correlation-only request IDs and existing API idempotency;
- attempt creation before lifecycle effects; and
- startup enumeration and recovery by committed `AttemptID` only.

Exit evidence: accepted or explicitly revised ADR text. A topology or responsibility change stops
this plan and requires an updated Proposed ADR before implementation.

## S1: passive contracts and fixed fixtures

Architecture decision: **follow ADR-0030's versioned atomic cutover**. The merged TypeScript
approved-byte contract requires three distinct future plan source roles, while ADR-0029's proposed
562-byte `RegisterPlanV0` record contains one source-manifest role and is v0-only. `RegisterPlanV0`
and 562 bytes remain historical/current-plan-v0 design only and will not be frozen into an S1
corpus. The observed 626-byte arithmetic is not an approved layout, cap, or known answer. See the
bounded [S1 consistency stop](AUTHENTICATED_LOCAL_IPC_S1_CONSISTENCY_STOP.md).

S1 remains blocked until this dependency order completes:

1. an accepted transformation-owner and immutable source-store topology;
2. finalized `ExecutionPlan` v1 and its complete nominal role model;
3. integration of the separately developed canonical field-authority manifest;
4. a newly versioned registration method and binding record with reviewed field order, exact caps,
   cap-plus-one behavior, and cross-language known answers; and
5. explicit review of the registration fetch, approval submission, and attempt-request projections,
   including method version changes wherever their typed shape changes.

Only then may S1 add no-product native, Go, and Swift-readable fixtures for the finalized common
header, four role-specific operations, success replies, and fixed refusals. The replacement fixture
manifest must retain exact plan/registration/approval and complete-binding known answers; protocol,
service, role, tag, audience, purpose, request, installation, epoch, length, deadline, and first-owner
fields; exact maxima and cap-plus-one; closed classification/reason and state/time/trust/core/store/
adapter oracles; copy-ownership mutations; response-loss/idempotency classifications; structural
missing/extra/wrong-type and cross-object refusals; and byte equality across implementations.

There is no dual active v0/v1 acceptance, optional transformation role, generic fallback, or field
inference. No fixture may import experiment code into product packages or claim ADR-0019 acceptance.

## S2: Go facade and store projection

Status: **blocked** pending completion of the ADR-0030 dependency chain and shared newly versioned
S1 bytes.
S2 must not freeze `RegisterPlanV0`, 562 bytes, or the 626-byte arithmetic, and it must not define
field authority inside the Go facade.

Implement an internal, unwired facade with one Go entry point per method. The facade does not
accept role or purpose from request bytes; each entry point constructs the existing fixed
`AuthenticatedCallContext` internally.

Work:

1. keep the finalized registration, approval-submission, and attempt-request entry points as thin
   method-specific calls to the existing components after method-owned admission;
2. add the intended read-only Broker fetch facade only after its typed projection and method version
   are explicitly reviewed, returning defensive copies of the retained exact plan, complete role
   bindings, and wire registration only after Broker call context, active installation/epoch/trust,
   registration binding, and expiry checks;
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

Status: **blocked** pending completion of the ADR-0030 dependency chain and shared newly versioned
S1 bytes.
Native parsing must not create a de facto field layout, cap, known answer, or method-version
decision.

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
mode-0600 owner-lock fixture. Do not add the executable to a product package or service manifest.

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

This slice is blocked until valid Apple Development identities and provisioning profiles for the
exact daemon, Broker, and Supervisor fixtures are deliberately made available. Do not inspect or
use unrelated identities.

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
- duplicate Supervisor startup and owner-lock contention; and
- exact installed protected store and epoch-scoped operational-key access denial, without exercising
  an Approval private-key operation unless separately authorized by the approval plan.

Exit evidence: Apple Development installed-topology result only. It does not claim Developer ID,
notarization, Gatekeeper, clean-host, shipping, or production authentication.

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
./experiments/gate-e-supervisor-topology/run-smoke.sh
```

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
