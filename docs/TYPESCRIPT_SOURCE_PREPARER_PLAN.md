# TypeScript Source Preparer implementation, conformance, and fault plan

Status: future-conditional proposed implementation plan. It retains a design decision only; no
Source Preparer, transformer consumer, product endpoint, store, key, runtime, backend, guest, or
deployment exists. Accepted ADR-0034 removes P0-P7 from the first-release critical path.

Owner: proposed TypeScript Source Preparer planning component.

Decision: [Proposed ADR-0032](adr/0032-select-enrolled-typescript-source-preparer.md).

## Defensive scope and claim boundary

The first release follows ADR-0034's single-file byte-exact `.mjs` plan-v0 contract. Run this plan
only if Capsule later reselects TypeScript; it does not block first-release source/plan, IPC,
runtime packaging, or admission work.

Defensively validate the pre-registration TypeScript approved-byte control using only passive
repository fixtures, exact retained Node 22.22.1/Amaro 1.1.5 benign sources, fixed local process
harnesses, temporary local stores, and deliberately owned installed identities when a later slice
authorizes them. Do not access any unrelated system, identity, credential, user data, runtime,
backend, guest, deployment, signing service, or third-party target.

No slice may claim semantic equivalence, type correctness, Node API stability, production process
isolation, runtime admission, external isolation, authenticated production IPC, or readiness to
execute hostile code. A skipped credentialed, sandbox, update, pressure, or power-loss case is not
a pass.

## Dependency order

```text
P0 Proposed ADR-0032 adversarial authority/TCB review
  -> P0A protected-store, worker-confinement, bootstrap/update, and authority-schema closure
       -> P1 passive method/store objects and exact cross-language fixtures
            -> P2 unwired Go validator and fault-injected immutable temporary store
                 -> P3 fixed one-shot Node worker and governed package evidence
                      -> P4 no-product native authentication/sandbox composition
                           -> P5 Apple Development installed/update/recovery matrix
                                -> P6 ADR-0030 ExecutionPlan/RegisterPlan v1 atomic cutover
                                     -> P7 Developer ID/notarized clean-host activation review
```

P1 is on hold at the retained
[P0 authority/TCB checkpoint](TYPESCRIPT_SOURCE_PREPARER_P0_AUTHORITY_REVIEW.md). P0A must satisfy
every checkpoint entry criterion and stop condition first. P1 may then prepare fixture readers in
Go, TypeScript, and Swift in parallel only after the canonical object inventory and a recursive
field-authority schema/verifier design are reviewable together. P2 and P3 remain unwired and may
proceed independently after P1. P4 composes them. P6 cannot begin until P1-P5 evidence either
closes the relevant gate or the ADR is revised. Authenticated-local-IPC S1 resumes only against the
versioned v1 registration shape selected in P6; the legacy 562-byte `RegisterPlanV0` record is never
frozen or reinterpreted.

## P0 — decision and authority review

Review ADR-0032 with ADR-0010/0011/0013/0018/0026/0029/0030/0031, the merged S1 consistency stop,
and the canonical field-authority manifest. The current manifest covers the passive approved-byte
family but not plan v1 or Source Preparer methods/store state. P0 must not add unverifiable manifest
targets before their canonical objects exist; P1 must add the objects, entries, and verifier support
atomically rather than creating a parallel classification. Confirm:

- a separate unprivileged planning component, not a generic or privileged helper;
- no Broker, Supervisor, updater, daemon, SDK, runtime, or backend transformation owner;
- two role-specific services and eight closed methods only;
- source-set ID is a trusted-resolution input, not plan bytes or execute-time authority;
- Source Preparer owns every stored source/object byte and only the mutable retention state;
- daemon, Supervisor, and Broker originate and validate only their assigned fields;
- no path/descriptor/transform option crosses a caller boundary; and
- plan v1 registration replaces rather than extends the v0 binding record.

Exit evidence: the retained P0 checkpoint, ADR/threat/control refinements, and an explicit P1 HOLD.
P0A exit evidence is: an exact OS-enforced single-member store container with hostile same-user and
stale-process negative probes; an exact worker child profile or an explicit decision that Node is
full store TCB; a sealed installer-owned genesis/update descriptor; settled archive/release and
cancellation/death semantics; canonical Source Preparer object identifiers; and a field-authority
schema/verifier design that classifies every nested tuple member and permits one required validator
per independent consumer. Any new key, backend route, Broker content access, Supervisor transform
responsibility, generic lookup, third service, or weakened same-user model stops the plan and
requires another ADR revision.

## P1 — passive objects, calculated maxima, and known answers

Add candidate-only CDDL or an equally closed internal binary contract for:

- `PrepareSourceSetV0` request/result and idempotency identity;
- the immutable source-set descriptor and mutable lifecycle record;
- `ResolveSourceSetForRegistrationV1` request/result plus the exact registration-intent binding;
- `AbortSourceRegistrationV1` request/reply;
- retain/read-approval/read-executable/release/abandon requests and fixed replies;
- fixed worker request/result/refusal frames; and
- fixed refusal/fault observations with state, copy, call, and cleanup oracles.

Generate exact maxima from the closed contracts. The generator must prove the 32-file,
262,144-byte file, independent 1,048,576-byte aggregate, 8,192-byte path-sum, record-count, object-
depth/member, and method-specific copy budgets. It must fail repository verification if a field is
added without updating the authority manifest or calculated maximum. Do not use the ADR's
2,359,296-byte reservation allowance as a wire cap; each accepted contract owns its exact encoded
maximum.

Required positive known answers:

1. all-JavaScript pass-through with zero records;
2. one exact retained 391-byte ordinary TypeScript transformation;
3. mixed TypeScript/JavaScript set with ordered records and unchanged paths;
4. exact 32-file/file-size/original-aggregate/emitted-aggregate maxima;
5. exact idempotent prepare replay;
6. plan-v1 source resolution with all roles;
7. registration-intent resolve, abort, and exact replay;
8. registration retain and approval readback;
9. executable-only attempt readback; and
10. final idempotent release refusal with a retained tombstone and unchanged blob references.

A positive blob-reference decrement is deliberately excluded from P1 until the archive/retention
decision selects one exact durable release authority. When selected, it needs a new positive known
answer and crash/replay/capacity oracles before implementation.

Required negative oracles include:

- every missing/extra field, wrong type/domain/version/service/role/purpose/audience;
- zero/33 files, duplicate/unsorted/dot/overlong paths, bad segment, entrypoint mismatch;
- per-file, original aggregate, emitted aggregate, path sum, object, frame, and reply cap-plus-one;
- invalid UTF-8, BOM, LF/CRLF, composed/decomposed Unicode, and trailing-newline distinctions;
- `.tsx`, `.cts`, JSX, enum, namespace, decorator, parameter property, malformed syntax, diagnostic;
- profile, Node, Amaro, source archive, distribution, executable, bootstrap, options, media,
  original, emitted, record, order, map, URL, and diagnostic-count mutation;
- original/executable/profile/options/record cross-domain digest substitution even when widths or
  bytes match;
- unknown or stale source set, nonce/body conflict, expired prepare, wrong installation/epoch,
  wrong/reused registration intent or plan digest, abort after retain, unretained approval read,
  executable read without committed attempt, and daemon release of a registered set;
- caller/native/Go/worker/store/accessor/reply buffer mutation after every copy boundary; and
- old v0 binding presented to v1, v1 binding presented to v0, and a synthetic 626-byte arithmetic
  record with no selected schema.

Exit evidence: independent Go/TypeScript/Swift agreement where each language participates,
fixture-generator idempotence, exact file hashes/counts, recursive nested-member field-authority
coverage with all required independent validators, and no product import or endpoint.

## P2 — unwired Go owner, validators, and temporary store

Implement an internal Go core behind injected interfaces; connect no XPC listener and invoke no
Node binary. A fixed fake transformer returns only fixture outputs or fixed faults. The Go core must:

1. copy and revalidate untrusted source entries before allocating aggregate buffers;
2. reserve transformation memory and durable capacity before calling the transformer;
3. construct every digest and ADR-0030 object itself;
4. write a temporary role-namespaced store with create-exclusive/no-link checks, file sync,
   directory sync, atomic same-volume publication, reopen, and full graph validation;
5. return defensive copies only;
6. implement nonce replay, prepared-only expiry, intent resolve/abort, retain/read/release,
   quarantine, and capacity exactly;
7. open the store without automatic creation and hold one owner lock; and
8. expose no arbitrary lookup, enumeration, mutable reference, path, descriptor, transform profile,
   or caller-selected cleanup operation.

Inject faults before and after reservation, staging creation, every write/sync/rename/directory
sync, descriptor/lifecycle commit, publication, reopen, response, retain, read, release, reference
decrement, tombstone, quarantine, and owner-lock acquisition. For every point assert exact store
tree/digests, state, reference counts, reservation release, fake-transform call count, reply or
interruption, and restart result.

Specific recovery oracles:

- before publication: no visible ID and only provably owned staging is removable;
- after publication/before reply: exact nonce replay returns the one committed ID;
- indeterminate rename or cross-record result: `RECOVERY_REQUIRED`, no second set, no empty store;
- resolve response loss: one pending intent and repeatable copied readback;
- validation failure/restart before registration commit: exact intent abort returns only that set
  to prepared state;
- retain response loss or restart after registration commit: one registered reference and
  repeatable acknowledgement;
- release response loss: one tombstone and no reference underflow;
- release with an ADR-0031 cohort still hot, an evidence/reproduction retention hold, or no
  separately validated reproduction archive refuses without decrementing blobs;
- corruption, unknown entry, link, wrong owner/mode, role mismatch, count overflow, or version
  mismatch: quarantine/repair-required without rewrite;
- capacity exhaustion: existing prepared/registered sets remain readable, no eviction; and
- duplicate owner: second process refuses before mutation.

Exit evidence: focused Go tests, race/repetition tests, mutation tests that restore each missing
check, and an explicit statement that the fake transformer proves no Node behavior.

## P3 — fixed worker and governed packaging

Build a no-product one-shot worker from the exact Node 22.22.1 macOS arm64 distribution. The fixed
bootstrap calls only `node:module.stripTypeScriptTypes` with mode `strip`, no source URL, and no
source map. It accepts one closed input frame and emits one closed output/refusal frame.

Package and verify:

- source/distribution/executable/Amaro/bootstrap hashes and complete notice/license inventory;
- fixed argv, empty environment, working directory, inherited descriptor allowlist, and no loader,
  config, package, network, debug, inspector, addon, or user path input;
- an exact OS child profile that denies the source store and package tree after executable open,
  network, process discovery/attach, descendant persistence, dynamic/native loading, arbitrary
  filesystem access, and every descriptor except bounded stdin/stdout/stderr and the process handle;
- exact regular-file/no-link/non-user-writable installed artifact custody;
- worker exit, signal, timeout, pipe truncation, extra frame, trailing bytes, stdout/stderr flood,
  parent cancellation, parent death, pre-bind death, and process-handle/start/code-identity orphan
  reaping without PID-only action;
- two-worker/8,912,896-byte reservation saturation and a third-call `CAPACITY` refusal; and
- all retained TypeScript semantic/cap/mutation fixtures, 100 repeated same-process-set runs, and
  25 fresh-process runs.

The parent must ignore worker-supplied digests, counts, media, and records and calculate them itself.
Change the worker to lie about each one and prove the parent either ignores it or refuses. Exercise
parser crash and hang only with fixed benign fault hooks or retained non-hostile fixtures.

Exit evidence: reproducible exact package inventory, SBOM/notices, worker fixture hashes, timing and
peak-memory distributions at exact maxima, and explicit limits. This slice does not establish App
Sandbox, signed peer identity, independent rebuilding, or production source/license closure unless
the retained evidence actually does so. If the child profile cannot prove store and package
separation, Node/Amaro remains full Source Preparer/store TCB; if that consequence is unacceptable,
the TypeScript path stops.

## P4 — no-product native authentication and composition

Compose the native XPC/Security front end, Go owner, temporary store, and exact worker in a local
disposable harness. Do not add a product service manifest. Ad-hoc identities may establish API and
exact-CDHash behavior only.

Prove authentication and allocation order for both service roles: listener requirement, message-
derived `SecCode`, EUID/session, connection/flow slot, fixed outer shape and data lengths,
installation/epoch/audience/purpose, copied bridge, then Go. Wrong identity or saturation must reach
neither body copy nor store/worker. Test exact/stale CDHash, same-identifier impostor, unsigned,
byte-identical copy, debug state, PID/name/path substitution, reconnect, service death, four-
connection/global saturation, per-connection second call, and every two/five/ten-second boundary.

Compose the full fault matrix at:

1. peer admission and flow reservation;
2. body copy and bridge dispatch;
3. worker artifact validation and spawn lease;
4. worker input/output frames and termination;
5. every P2 store publication/recovery boundary;
6. response copy and response loss;
7. Supervisor intent commit/resolve/validation/abort or registration commit/retain acknowledgement;
8. Broker approval read and copied reply;
9. attempt commit and executable read; and
10. release/tombstone/startup cleanup.

The Supervisor and Broker validators in this slice are independent no-product test doubles using
the eventual Go and Swift decoder contracts. Deliberately omit one relationship check from each,
run the corresponding mutation, and prove the evidence fails until restored.

Exit evidence distinguishes OS enforcement, protocol enforcement, Go/store mechanics, Node worker
behavior, inference, and untested claims. Ad-hoc evidence cannot support Team/channel enrollment,
provisioned entitlements, protected containers, wrong-session switching, or production IPC.

## P5 — Apple Development installed and update matrix

This slice requires deliberately authorized Apple Development identities and provisioning profiles
for the exact daemon, Source Preparer, Supervisor, and Broker fixtures. Do not inspect or use
unrelated identities.

Current non-secret resource discovery: the intended Team is `W4QUR9FUL4`, but exact G3 readback
found that the certificate displayed with that suffix has subject OU and emitted TeamIdentifier
`3DDR84M4JS`; it is not W4 evidence. Local signed/provisioned experiments require a matching W4
certificate plus exact role identifiers, entitlements, and profiles. Xcode 26.6 cached three
profiles, all for historical Team `3DDR84M4JS` (Gate B Broker, Gate B Supervisor, and wildcard);
they are not reusable for W4 tests
and do not include a Source Preparer role. A separate Developer ID Application identity for
historical Team `3DDR84M4JS` is later distribution authority requiring explicit authorization and
matching-Team package design. It must not be used as W4 development evidence and does not make
Developer ID or notarization work current.

Install the Source Preparer as its proposed unprivileged `SMAppService.agent` and repeat P4 with:

- exact Team/channel/role identifier, CDHash set, entitlements, Hardened Runtime, absent debug
  authority, message-derived identity, EUID, and Aqua session;
- protected private store with no shared app group; same-user read/write/link/replace/attach probes;
- worker denial of network, arbitrary filesystem, environment, addon, inspector, debug, and service
  access while fixed pipe transformation succeeds;
- logout/login, fast-user-switch where owned, lock/unlock, sleep/wake, launchd restart/backoff,
  memory/disk/process pressure, duplicate service start, owner-lock contention, and orphan cleanup;
- old/new daemon, Source Preparer, Supervisor, Broker, Node/bootstrap, entitlement, store format,
  plan version, and epoch combinations;
- crash at every prepared-update swap/migration/acceptance/finalization boundary;
- reversible pre-commit rollback, forward-only post-commit repair, old-binary refusal, partial
  restore, corrupt backup, full store, and indeterminate rename; and
- registered-set survival without re-transformation, approval refusal while retention is pending,
  and executable-only read after a committed attempt.

Exit evidence is limited to the exact Apple Development package/host matrix. Power interruption,
Developer ID, notarization, Gatekeeper, clean-host, minimum-OS, and production identity remain
separate until actually tested.

## P6 — atomic plan/registration/approval/lifecycle cutover

Only after P0-P5 review, execute ADR-0030's coordinated authority cutover in one reviewed release:

- accept final object CDDL/media/profile/store contracts and their Go/TypeScript/Swift views;
- replace `ExecutionPlan` v0 with v1 original/executable/record-set roles and both aggregates;
- define `RegisterPlanV1` and a newly calculated complete binding record including the internal
  Source Preparer resolution reference and registration intent; never freeze or reinterpret 562
  or assume 626 bytes;
- implement Supervisor source resolution, complete validation, pending-retention registration,
  Broker copied rendering, attempt projection, executable-only staging, lifecycle binding-set,
  transcript, receipt, and archive release;
- replace all downstream known answers together; and
- migrate one locked store/version with validate/write/sync/rename/sync/reopen, old-binary refusal,
  and no dual active plan or mixed source-reference shape.

Authenticated-local-IPC S1 then generates its passive native/Go/Swift fixtures from this v1 shape.
`RequestAttempt`, `GetRegisteredPlan`, and `SubmitApproval` are versioned wherever their projection
changes. No optional transform fields, generic envelope, compatibility retry, or v0/v1 authority
fallback is allowed.

Exit evidence: full exact-byte cross-language corpus, complete field-authority manifest, source-
store readback, response-loss/fault matrix, Broker complete rendering, execution retrieval of only
the emitted role, offline migration tests, and removal of old active acceptance.

## P7 — distribution and activation review

Repeat the installed matrix with final intended Developer ID/notarized/stapled package bytes on
owned clean disposable hosts at every supported macOS floor. Include install, update, reboot,
logout/login, sleep/wake, pressure, locked Keychain, restore, repair, stale peer, stale source store,
and network-offline reproduction. Read back exact identities, entitlements, service registration,
protected storage, store checkpoint, trust epoch, TUF target/profile identities, and stale-component
denials. Before any blob-release behavior activates, separately select and validate either a
bounded retention policy that explicitly ends offline reproduction or an immutable reproduction-
archive format/owner with crash, corruption, capacity, restore, and confidentiality tests.

Activation additionally requires accepted ADR-0019/0026/0029/0030/0032 as applicable, production
approval verification, Supervisor archive/compaction, content/runtime/backend/profile admission,
evidence composition, and every independent blocker. Completing this plan alone does not admit a
runtime or hostile workload.

## Cross-slice verification

Every retained slice runs the repository-required verification and its focused local harness:

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

Also validate relative Markdown links, generated maxima, fixture idempotence, exact tool versions,
store tree ownership/modes, and absence of product imports/endpoints before P6. Record every skipped
platform or credentialed case explicitly.

## Falsifiable security claims

| Claim under test | Required falsifier/oracle |
| --- | --- |
| daemon cannot mint trusted emission | fabricated set ID/object/profile or daemon-direct worker call refuses before registration state |
| no path becomes authority | any request/reply/store accessor accepting or returning a host path, URL, fd, bookmark, or mmap fails conformance |
| worker cannot choose authority | changed argv/env/loader/options/digest/record/output destination refuses or is ignored and parent-recomputed |
| protected source-store ownership | daemon/Broker/Supervisor/updater/stale-Preparer/unrelated same-user opens, links, replaces, maps, renames, or retains no handle to the store; unknown mutation quarantines |
| worker authority is bounded or explicit | child-profile probes deny store/package/network/process/native-loading authority; otherwise evidence labels Node/Amaro full store TCB |
| independent registration validation | each digest/media/path/order/member/pass-through/record mutation refuses before registration commit |
| independent approval validation | same mutations refuse before Broker signature operation |
| registration binds retention | approval/attempt remains disabled until idempotent source retain is committed |
| execution receives emitted bytes only | original/profile/options/record/transform request/path at read/stage yields zero runtime/backend calls |
| response loss does not duplicate authority | exact prepare/retain/release replay converges; nonce conflict refuses |
| cancellation/deadline is bounded | worker and reservation terminate; no hidden second set; post-commit replay returns one set |
| recovery does not invent state | corrupt/missing/mixed store quarantines; no empty recreation or time-based registered eviction |
| update has one active world | every old/new component/profile/store/plan mix refuses until one epoch is accepted |
