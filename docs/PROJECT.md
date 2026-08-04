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

Current work uses the [canonical status language](STATUS_LANGUAGE.md). In short: the passive
single-file `.mjs` foundation, bounded Oxc parser/process selection, no-guest Supervisor lifecycle,
and G2 local owner composition are `PASSED` in their exact scopes. Archive F2's fixed-store
v1-to-v2 migration/full verifier, F3's first immutable-segment activation, and F4A's read-only
retained lookup/replay/uniqueness routing are also `PASSED` in their narrow local-conformance
scopes. Product
Source Validator work, installed owner-lock G3, and runtime/profile admission are `BLOCKED` on the
named dependencies below. Governed `deno_core` and governed libkrun remain `IN_PROGRESS —
TRENDING_GOOD`; the next local archive work starts at F4B. None of those blocked or incomplete parent
items is `NO_GO`. Product admission and control-evidence maturity remain separate from work status.

The concise current dependency and claim checkpoint is
[Phase 2B and Gate C current maintainer checkpoint](PHASE_2B_GATE_C_TASK_GROUP_CHECKPOINT.md).
This section summarizes durable repository state; the linked checkpoint separates selected design,
implemented local mechanics, experiment evidence, governed external source, and product admission.

The reviewed [macOS installation and distribution plan](MACOS_INSTALLATION_AND_DISTRIBUTION_PLAN.md)
now records the intended user-facing shape as one Swift `Capsule.app` distributed in a DMG with
embedded, narrowly enrolled per-user components. That direction is `IN_PROGRESS — TRENDING_GOOD`,
not installed-product evidence. A developer-signed MVP deliberately starts with manual whole-app
replacement; automatic TUF updates, a custom Bundle Replacer, Developer ID/notarized distribution,
minimum-OS support, and complete uninstall semantics remain later work. The protected Supervisor-
container bootstrap, App Group/private-XPC residual authority, two role-specific Source Validator
launchers, and any Trust Coordinator or Bundle Replacer authority require separate decisions and
signed installed evidence.

The current JSON Schemas and TypeScript `Job` interfaces are pre-freeze scaffolding. They describe
the repository's current API surface, but they intentionally do not define the target v0 protocol.
Blocking feasibility spikes must determine what the platform can actually enforce before those
contracts are replaced and frozen.

Phase 2A now includes verified passive candidates for the first narrow `JobProposal`, minimum
`ExecutionPlan`, and `PlanRegistration`, with byte-exact internal fixtures and Go/TypeScript
decoded views. They do not replace the scaffold, activate an endpoint, or authorize execution.

Accepted ADR-0034 freezes the first-release source contract to one dependency-free, byte-exact
`main.mjs`. The merged M1 passive byte/manifest foundation is canonical. A bounded parse-only
comparison now supports [Accepted ADR-0035](adr/0035-select-disposable-mjs-source-validator.md):
exact Oxc 0.140.0 parser/AST/semantic checking is the engineering candidate inside a future
one-shot disposable Source Validator, invoked independently before planning and before approval.
The experiment is retained in the commit-pinned
[`mjs-parser-boundary` archive](https://github.com/Shrimpworks/capsule-experiments/tree/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/mjs-parser-boundary).
The passive V0 slice freezes fixed request/result/candidate/artifact-profile frames and independent
Go/Rust test oracles without launching or invoking a parser. The
bounded V1 follow-up now retains one unwired macOS arm64 artifact at exact digest
`ba2a6b38be6b4eea8c067887cf80988756e2f4a551d128bf2dabdaf7f2ecb600`, its Rust 1.95.0 locked
offline graph, source/license/notice/SBOM/provenance manifests, V0 result agreement, and same-host
two-directory byte reproduction. Its identity-free linker ad-hoc signature is recorded, but it is
not installation-signed or enrolled: independent-builder, vulnerability-owner, V2 sandbox/resource,
V3/V4 consumer, and V6 runtime evidence remain blockers. No product endpoint, runtime no-loader
boundary, runtime, backend, or guest is implemented or admitted.

The retained V2 macOS checkpoint now makes that process-profile blocker exact. A strict test-only
bootstrap fixes copied I/O, argv, empty environment, cwd, descriptor closure, CPU/file/FD/child/
wall limits, kill/reap, and fault refusal, but `RLIMIT_AS` returns `EINVAL` before the artifact
executes. An explicitly unbounded diagnostic mutation proves file reads, IPv4/Unix socket creation,
cwd metadata writes, and a 512 MiB mapping remain possible. Apple's supported embedded-tool App
Sandbox entitlement shape changes the fixed V1 bytes, while deprecated custom sandbox profiles are
not acceptable evidence. V2 and the product work are therefore `BLOCKED`, not `NO_GO`. V0/V1/V2
bytes and status remain unchanged; the supported replacement proceeds only through R1-R5B under
ADR-0036.

The supported-profile replacement research and R0 architecture-decision slices are `PASSED` in
their exact scopes and keep the product parent blocked. Direct App Sandbox helper inheritance is
`NO_GO` because it preserves daemon/Broker static rights. [Accepted ADR-0036](adr/0036-select-role-separated-source-validator-launchers.md)
selects two role-specific private App-Sandboxed XPC launchers, one per consumer, with distinct
services/parser identities/profiles and no shared bus, result, cache, container, group, or key.
Each private container is accepted only as residual scratch authority: no persistent Capsule
product state, source/diagnostic log, cache, or reusable result is permitted, and cleanup/residue
testing is mandatory after every request, crash, restart, update, and startup without becoming a
confidentiality proof. The unavailable hard ceiling is replaced with a later evidence-derived
reactive physical-footprint watermark, one direct child per launcher request, bounded combined
two-role concurrency, and kill/drain/reap; it is not a hard peak/exact cap or host-availability
guarantee. Threshold, cadence, baseline, and overshoot remain unset until the separately authorized
signed corpus. Resume with the [new passive v1 boundary](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1.md)
and sequential R1-R5B plan. No user signing identity was used or authorized.

Phase 2B now retains a closed 94-rule, 458-case, 561-fixture conformance corpus for raw decoding,
proposal/source/input semantics, exact plan and registration bytes, registration-state transitions,
and the unwired ADR-0024 approval/attempt boundary. The unwired implementation verifies 80
TypeScript proposal targets: 62 raw/schema cases plus all 18 semantic-resolution cases. It also
verifies those 80 proposal targets plus 40 MJS byte/manifest targets in TypeScript. Go verifies
345 targets: the previously retained 177 plus the same 40 MJS byte/manifest targets and 128
passive Source Validator frame targets; a standalone test-only Rust oracle verifies those same
128 frames and is not a product dependency. The earlier
177 comprise 81 internal-CBOR/wrapper cases, all 40 registration-state cases, 44
passive approval/attempt cases, and 12 fixed-store transition oracles. The approval/attempt work
adds distinct identifiers/references, the fixed internal classification vocabulary, the exact
candidate approval known answer, a bounded fixture-only verifier, and one unwired transactionally
colocated fixed registration/approval/attempt store. The no-guest fake lifecycle now resolves and
recovers only committed `AttemptID` records, revalidates exact plan and copied bindings before fake
prepare, and retains the original 12 top-level focused lifecycle tests for binding, replay,
concurrency, fault, and startup-recovery behavior. The E4/E5 local checkpoint adds focused durable
tests and now drives that no-guest fake through
the colocated v1 fixed snapshot and retains exact effect identities, restart reconciliation,
256-active/4,096-retained lifecycle ceilings, and repeated-startup/exhaustion evidence. Proposed
ADR-0033 now selects an enrolled
pre-created sibling object plus lifetime nonblocking BSD `flock` as the later macOS owner
mechanism after a bounded local process/file experiment. Passive G1 adds the internal opaque
Go/Darwin acquisition capability. G2 now composes that capability with an owner-required current
v1 opener, the same owner-session ID in the store and per-attempt coordinator, sorted `AttemptID`
recovery, permanent post-open ownership fencing, and store-before-owner shutdown ordering. The
retained race/fault/process tests use owned temporary roots and the no-guest fake only. This is not
product startup, a signed bootstrap record, protected installed storage, v2/archive composition,
or production evidence. The first bounded G3 discovery found a certificate/profile identity
mismatch before installed build. Installed G3 is **BLOCKED**, not abandoned:
the authorized Apple Development certificate's display name ends in `W4QUR9FUL4`, but its public
subject OU and an exact signed-byte `TeamIdentifier` are `3DDR84M4JS`; all three cached profiles
are also 3DDR. The retained G3 fixture fixes test-only role/state/bootstrap fields and reruns the
noncredential G1/G2 corpus, but no W4 app, service, protected root, signed per-installation record,
session, or update case ran. A matching W4 certificate/profile set, selected protected-root
bootstrap ceremony, signed record, and descriptor-relative store open remain blockers. There is no
consumer, authenticated IPC, production approval, evidence, real backend, runtime, or guest.
[Proposed ADR-0025](adr/0025-colocate-durable-attempt-lifecycle-state.md) selects a colocated
lifecycle record/effect-checkpoint extension to the same Supervisor snapshot, with a separate
fake-only implementation plan. Slices E1 through E5 now provide the passive lifecycle contract,
explicit fixed-store v0-to-v1 migration and validation, durable ensure/read/intent/result/
indeterminate/reconciliation/recovery-set transactions, the FakeBackend-only driver, and the
capacity/repeated-startup evidence checkpoint. ADR-0025 remains Proposed and the implementation
status remains unwired `local-mechanic` only.
[Proposed ADR-0031](adr/0031-checkpoint-closed-supervisor-cohorts.md) defines the exact
archive/compaction boundary: complete expired registration cohorts move from the mutable snapshot
to immutable retained segments only after every bound attempt is durably destroyed with
authoritative absence, while full records and replay/non-reuse tombstones remain retained. It
selects a finite fixed-store checkpoint only as a local conformance oracle, not a production
engine, and authorizes no referenced-history deletion. Slice F1 now implements passive archive
types, known-answer digests, defensive copies, and the pure complete-cohort eligibility selector.
It writes no file, migrates no store, moves no cohort, activates no archive, and routes no consumer.
The first F2 review stopped before implementation because merged F1 could not represent a required
nonzero visible-v1 effect seed with zero archive descriptors, reconstruct global hot indexes, or
construct the specified generation-one migration checkpoint. The passive
[F2 format blocker resolution](SUPERVISOR_ARCHIVE_F2_FORMAT_BLOCKER.md) now freezes separate global/
segment index domains, typed hot/archive record locations and counts, a distinct migration-genesis
checkpoint, and generated known answers. It writes no store bytes and changes no v1 behavior. The
follow-on valid-v1 missing-lifecycle contradiction is now also passively resolved with a closed
absent/present lifecycle union on each attempt entry and lifecycle counts derived only from present
records. The executable [F2 v1 mapping resolution](SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md)
retains the real committed-attempt-before-lifecycle witness and exact `attempts = 1, lifecycles =
0` genesis answer. Stateful F2 now implements the owner-asserted closed v1-to-v2 migration and
read-only full verifier as a finite local oracle. It preserves that witness without recovery or
invention, reconstructs all-hot retained-global indexes and migration genesis, and retains exact
known answers plus pre/post-rename, corruption, capacity, concurrency, and subprocess-death tests.
See the [F2 stateful migration result](SUPERVISOR_ARCHIVE_F2_MIGRATION_RESULT.md). It moves no
cohort, creates no segment, routes no retained lookup, calls no adapter, and is not a production
store or durability claim. Stateful F3 now adds exactly one sealed segment prepare/verify/publish/
activate transaction. It publishes the fully reopened digest-addressed immutable segment before
atomically installing its generation-two active reference, preserves every full cohort record and
visible identity/tombstone, reports valid unreferenced artifacts without deleting them, and retains
pre/post-publication, pre/post-activation, response-loss, corruption/substitution, concurrency,
owner-loss, and subprocess-death oracles. See the
[F3 stateful activation result](SUPERVISOR_ARCHIVE_F3_ACTIVATION_RESULT.md). It adds no retained
lookup or v2 authority mutation, second segment, orphan deletion, backup, adapter, IPC, runtime,
backend, guest, production engine, or real power-loss claim.
Read-only F4A now adds fresh full-verification lookup of retained registration, approval, attempt,
nonce, effect, instance, and replay identities through typed hot/archive retained-global locations,
plus passive collision queries and hot-only `AttemptID` recovery. See the
[F4A retained lookup result](SUPERVISOR_ARCHIVE_F4A_LOOKUP_RESULT.md). It adds no v2 write,
authority/lifecycle mutation, new effect tombstone, second segment, adapter, consumer, runtime,
backend, or guest. F4B mutation and F4C bounded growth remain open.
[Proposed ADR-0029](adr/0029-select-authenticated-local-ipc-topology.md) now selects one
unprivileged per-user Supervisor process with a small native XPC/Security front end and the existing
Go authority/lifecycle core in-process. It defines two role-specific Mach services and four closed
calls, but no bridge fixture, installed authenticated endpoint, product consumer, or production
identity evidence is implemented.
A focused unwired TypeScript Task 3C slice constructs and deterministically encodes the minimum
`ExecutionPlan` from only Task 3B provenance-bearing plan inputs and separately issued trusted role
bindings. A later completed focused slice prepares a defensive exact-byte/complete-role handoff and
exercises those values
against the real Go `registrationstate` component through a local-only conformance command. Go
independently predecodes, role-binds, hashes, and retains the 530-byte known answer. This is not a
product IPC implementation: the native-to-Go bridge, authenticated cross-process transport,
consumers, and endpoints remain pending. See the
[unwired decoder checkpoint](PHASE_2_UNWIRED_DECODER_CHECKPOINT.md). The latest local Gate C
checkpoint rejects stock Bun 1.3.14, its governed-construction branch, hardened full Deno v2.9.4,
and the tested minimal `deno_core` 0.409.0 construction for the required runtime-authority contract.
The later governed `deno_core` patched construction passed only the narrow physical-omission
question. Its package follow-up reproduced the exact snapshot and binary on one owned host from a
digest-pinned no-apt
builder and complete locked Cargo source bundle. The later self-contained-root follow-up closed its
standalone dynamic Bookworm-root blocker with a 22-entry package-derived root and no ambient
library/config fallback. The subsequent exact V8 closure trace proved the official asset,
publisher job, source gitlinks, V8 base, and patch stack. It rejected that exact official-asset
publication route as `NO_GO`: mutable publisher inputs, missing GN/Ninja link metadata, and absent
generated notices prevent an independent rebuild and complete notice closure. The replacement
governed-fork path is `IN_PROGRESS — TRENDING_GOOD`; independent-builder provenance, governed
release construction, and admission remain open. Accepted
[ADR-0028](adr/0028-select-governed-deno-core-first.md) selects governed `deno_core` as the first
runtime engineering candidate after the hard Bun pivot and records the real `Shrimpworks/deno` and
`Shrimpworks/rusty_v8` forks. Their first governed branches are merged at exact commits. The first
fork-native integration check stopped before construction because the merged `rusty_v8` fork had
only a Linux/amd64 builder. Governed `rusty_v8` PR #4 is now merged at head
`80e863ddb942a4aa2b384e794fc23e35b9d2bb15` and merge
`cbf56de2e1156b1cf1561fdbaea7172a0aa056f4`. Its clean Linux/arm64 build, fixed `get_version`
test, corrected GN evidence query, complete network-disabled build, and unsigned bundle upload all
passed at the merged head. The fork and its retained branch, PR, and Actions history have transferred
from `dills122/rusty_v8` to `Shrimpworks/rusty_v8`. This closes the fork's bounded ARM64 construction
blocker; independent-builder equality, evidence review, governed release publication, and runtime
admission remain open. No fork release or new admitted runtime artifact exists.
The decision supersedes ADR-0003's Bun-first
ordering only; it does not admit a
runtime, and `RUNTIME-001` remains unsupported. The
libkrun direct-block-root prototype made `NullFs` removal credible and selected `GOVERNED-PATCH`,
but the current and prototype profiles remain unsupported until final governed installed bytes
close P0-1 through P0-4. P0-3 now retains independent Go/Node verification of the 43 framing
vectors and a local process-pipe fault corpus. The public governed libkrun follow-up merged from
exact head `8a2c91943793668f31a1cf7af431933be935bb58` as
`cf0333cdba478cc34a8570a65b38412da7fd3ecc`. It retains the unchanged five-patch aggregate
`d19fd0ff159c699acccda2621519de45a09408bf3847b418ac34e02b79e805d5`, adds bounded console and
raw-FD library corpora, fixes two locally observed shutdown/lifecycle defects, and raises the four
console files from 13/88 to 37/88 covered functions and from 90/728 to 298/733 covered lines.
This is fork-source governance, local library, and CI build evidence only. Measured uncovered code,
the post-merge governed-branch/verifier pin mismatch, independent human/CODEOWNER review, the real
transport, launcher, guest/VMM, installed composition, signing/notarization,
distribution/source/SBOM obligations, and final profile reruns remain open.
P0-4A conditionally passed the no-host-root topology only; Gatekeeper, clean-host, and minimum-OS
admission also remain open as deferred activation/distribution evidence. They do not block current
local F2, G1, or documentation mechanics.
Accepted [ADR-0034](adr/0034-freeze-mjs-first-release-contract.md) now freezes the first release as
one byte-exact pass-through `main.mjs` member under the existing plan-v0 source role, with no
static/dynamic dependency request, CommonJS, package resolution, legacy Node module surface, or
module-loader fallback. The exact source-byte and deterministic-CBOR SourceManifest foundation is
implemented passively, but the module-request/source-language validator is on the retained
[grammar-counterexample hold](MJS_MODULE_REQUEST_VALIDATOR_HOLD.md). JobProposal narrowing,
plan construction, IPC registration/fetch, Supervisor source custody, Broker rendering, runtime
no-loader evidence, and every consumer remain unimplemented. The bounded TypeScript approved-byte follow-up passed only the pre-approval
byte-ownership question: exact Node 22.22.1/Amaro 1.1.5 strip-only emission was deterministic for the fixed
fixtures and Proposed ADR-0026 binds original and emitted roles before registration. That
experiment did not make the later ADR-0028 selection, choose a production transformer owner,
change current contracts, or admit a runtime. Proposed ADR-0032 selects a separate enrolled Source
Preparer with a one-shot exact Node worker and a role-namespaced immutable source store. That
conditional later feature is **BLOCKED** pending protected-store, worker-confinement,
genesis/update, retention/release, recursive field-authority, and lifecycle evidence. P1 passive
contracts have not begun. No component, store, endpoint, consumer, installed identity, or cutover
exists. ADR-0034
removes that conditional TypeScript path from the first-release critical path; ADR-0030/0032
remain future-conditional and still require an atomic plan-v1 cutover if resumed.
See the
[P0-0 construction review](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-bun-runtime-authority/governed-closure/CONSTRUCTION_REVIEW.md)
and [Deno-family disposition](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-runtime-authority/RESULTS.md)
and [governed package result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-core-reproducible-package/RESULTS.md)
and [V8 source/license closure](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-deno-v8-source-license-closure/RESULTS.md)
and [fork-native Linux/arm64 blocker](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/gate-c-fork-native-deno-runtime-bundle/RESULTS.md)
and [TypeScript approved-byte result](https://github.com/Shrimpworks/capsule-experiments/blob/0d8233b55f153b27a901a9ec45a3834208e3aa86/experiments/typescript-approved-byte-boundary/RESULTS.md)
and [governed runtime work plan](GOVERNED_DENO_CORE_WORK_PLAN.md)
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
- One byte-exact pass-through `main.mjs`; TypeScript remains conditional later work
- A pinned JavaScript runtime selected only after runtime/profile authority closure
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
- Go remains the daemon language and Swift/native remains preferred for the Broker. Proposed
  ADR-0029 selects an unprivileged per-user Supervisor that keeps the existing Go
  authority/lifecycle core and links a small native C/Objective-C XPC/Security front end in the
  same process. It adds no Swift Supervisor service, host-root process, or privileged helper.

## Success criteria

The first functional milestone succeeds when a client can submit a dependency-free inline JSON
job, obtain a Supervisor registration, receive explicit user-presence approval from the Broker,
consume that grant for one attempt, run an explicitly admitted dependency-free runtime in a
disposable development backend, release only a fixed agent summary, record backend-specific
teardown evidence, and compose a receipt from Broker approval evidence and a Supervisor enforcement
transcript. Accepted ADR-0028 makes governed `deno_core` the first engineering candidate, but this
milestone remains blocked until a separate profile-admission decision accepts exact runtime bytes.

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
[Gate C P0 reconciliation](GATE_C_P0_RECONCILIATION.md). Public precedents and the boundary between
reusable lessons and Capsule evidence are recorded in
[Related systems and design influences](RELATED_SYSTEMS.md). The canonical
[ecosystem reuse and adoption map](ECOSYSTEM_REUSE_AND_ADOPTION.md) ties platform, library,
standard, tool, and governed-fork choices to exact roadmap slices. Future work consults that map
and completes its dependency-policy checklist before adding a dependency or inventing a primitive;
the map does not admit a product component or advance a security claim.
