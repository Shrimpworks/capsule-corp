# Roadmap

The roadmap is ordered by uncertainty and risk reduction rather than feature count. Disposable
spikes may be built outside the final product shape. Their retained evidence—not prototype code
quality—is the deliverable.

Live per-item status — Archive F2-F6, owner-lock G2/G3, ADR-0034-0036 Source Validator R0-R5,
governed `deno_core`, macOS installation (including I1B), and the rest — is tracked in the
[current workstream dashboard](STATUS_LANGUAGE.md#current-workstream-dashboard), the single
current-status source; this roadmap defines phase order and exit evidence and does not restate live
item status. The
[Phase 2B and Gate C current maintainer checkpoint](PHASE_2B_GATE_C_TASK_GROUP_CHECKPOINT.md) adds
the separate evidence-class breakdown (selected design vs. implemented mechanics vs. experiment
evidence vs. product admission) for the latest task group.

Installation packaging now has its own staged
[macOS installation and distribution plan](MACOS_INSTALLATION_AND_DISTRIBUTION_PLAN.md). The current
direction is one visible Swift application in a DMG with embedded per-user services and no
permanent privileged helper. I0 passive role/bundle/bootstrap contracts and I1-I4 developer-signed
composition/manual-replacement work precede any automatic updater. TUF, a mechanical Bundle
Replacer, Developer ID/notarized distribution, a frozen minimum OS, and complete uninstall are
later I5-I6 gates and do not block the no-guest development MVP.

Every phase must consult the
[ecosystem reuse and adoption map](ECOSYSTEM_REUSE_AND_ADOPTION.md) before adding a dependency or
custom primitive. The task records the matching row and closed recommendation and completes the
dependency-policy checklist. `ADOPT-PLATFORM`, `ADOPT-PINNED`, or `GOVERN/FORK` is a planning
direction, not product admission or evidence closure.

## Internal-alpha critical path

Accepted [ADR-0040](adr/0040-freeze-owner-only-internal-alpha-posture.md) and the
[full audit synthesis](ALPHA_ARCHITECTURE_AND_RELEASE_AUDIT.md) now define the release path:

1. preserve the passed governed fork promotions, passive C2B v3 successor, and v4 build/static
   materialization;
2. preserve the passed separately authorized
   [v19 fixed benign owned guest](FIRST_OWNED_GUEST_EXECUTION_CHECKPOINT.md) as a controlled
   experiment without treating its diagnostic console proof as the final typed transport;
3. preserve the passed v20 no-launch denial-test materialization and its fail-closed pre-ready
   runtime refusal, v21's fail-closed ready-EOF result, and v22's exact
   `preflight-root-sha256` refusal; preserve v23's confirmed malformed embedded-digest diagnosis,
   v24's passed corrected preflight/known-answer/early denials, and its vsock-family stop; then use
   semantic review to preserve v25 as `NO_GO`, retain v26's passed consolidated localization and
   passive-interface-policy stop, then use the passed v27 down/unbacked-`dummy0` and route correction
   plus its passed exact fixed-denial execution without treating it as product admission;
4. connect one bounded authenticated CLI adapter to the passed exact single-`main.mjs` proposal and
   passive atomic plan/bindings/manifest/source custody, then implement Broker fetch/render/approval;
5. compose protected installed Supervisor state, the bounded fixed-store alpha exception, real
   adapter/recovery, and completion evidence; and
6. run the minimum hostile `.mjs` and restoration corpus with one fresh guest per attempt.

Host Source Validator R4/R5 is later defense-in-depth rather than an internal-alpha gate. Exact
R4-v1 candidates are `NO_GO`; R4-v2 is unexecuted and `BLOCKED`. F6 is also deferred only for the
strict owner-only disposable fixed-store posture. TypeScript, automatic update, Developer ID
distribution, clean-host/minimum-OS coverage, restore, and production storage are external-alpha
work. The fixed guest is not product alpha. V20-v23 refused before readiness without launching a
guest; v24 then launched, passed the known answer and early denials, and stopped in the vsock family.
V25 was abandoned before launch because it tested socket creation rather than usable vsock
capability. V26 then passed the active vsock/raw-block controls and stopped on an over-strict
interface-name policy. V27 corrected that policy and passed the complete exact fixed denial corpus
in one authorized owned guest. The hostile owner-only internal alpha is still the next product
checkpoint, and external alpha remains a later distinct gate.

## Phase 0: architecture and claim baseline

- Align project, architecture, technical design, threat model, security policy, and ADRs on the
  daemon/Broker/Supervisor authority split.
- Define the component-compromise and control-to-evidence matrices.
- Inventory the target protocol objects and mark current schemas/types as pre-freeze scaffolding.
- State evidence and runtime-integrity claims without implying platform attestation.
- State cross-domain microarchitectural leakage and elevated approval-automation capabilities as
  explicit posture limitations until exact platform evidence supports anything stronger.

Exit evidence:

- No document says the daemon can approve, launch, own user content, and author authoritative
  execution evidence.
- DIDs are consistently described as identifiers, not authorization roots.
- Apple Container and all runtime profiles are development posture.
- Current versus intended behavior is obvious.

## Phase 1: blocking feasibility spikes

Status: initial decision spikes and all five Gate C implementation-readiness tracks completed for
the currently available host/account environment. Gate C produced a required backend pivot and a
conditional native candidate, plus explicit blockers for the exact native profile. Runtime P0-0
has rejected the stock/governed Bun and tested full-Deno/minimal-`deno_core` constructions.
Accepted ADR-0028 now selects governed `deno_core` as the first runtime engineering candidate;
`RUNTIME-001` remains unsupported and no runtime profile is admitted. The follow-ups
passed physical omission, same-host byte reproduction, and the exact standalone dynamic-root
question; the TypeScript follow-up passed a strip-only pre-approval byte-binding question. These
results establish bounded construction evidence but do not admit that runtime. The real
`Shrimpworks/deno` and `Shrimpworks/rusty_v8` forks contain merged governed branches, but no governed
release exists. The first fork-native Linux/arm64 construction stopped before building because the
merged `rusty_v8` publication contract supported only Linux/amd64. Governed `rusty_v8` PR #4 now
is merged at exact head `80e863ddb942a4aa2b384e794fc23e35b9d2bb15` and merge
`cbf56de2e1156b1cf1561fdbaea7172a0aa056f4`. Its clean ARM64 build, fixed test, corrected GN
evidence query, network-disabled full build, and unsigned bundle upload passed. The fork has
transferred to `Shrimpworks/rusty_v8` with its branch, PR, and Actions history intact. Independent-
builder equality, evidence review, governed release publication, and admission remain open. The
merged experiments handoff reconstructs the eight exact Linux/arm64 runtime/root identities, and
the passive C1 fixture freezes their intended JSON-in/JSON-out composition without creating a
guest. Passive C2A and immutable C2B v1 retain the refusing profile and historical build candidate
without execution. The Deno and experiments dependencies later merged. The distinct C2B v2
successor pins reviewed capsule-experiments PR #4 and closes six no-guest build identities plus an
unadmitted manifest candidate; it keeps the host runner non-final and every composed-profile,
resource, guest-evidence, and admission field null. Same-host equality is not independent-builder
equality. The v3 passive successor binds current accepted governed source heads/trees, preserves
the stale libkrun dylib as evidence only, resolves boot/runner/device/runtime/resource/teardown
semantics, and supplies an immutable passive-contract digest without adding a consumer. The v4
successor retains exact accepted header/current-source libkrun and final-runner bytes, independently
audits their ABI, and supplies a new immutable composed digest without execution or a consumer. A
later separately authorized [v19 experimental successor](FIRST_OWNED_GUEST_EXECUTION_CHECKPOINT.md)
passed one fixed benign owned guest without mutating v4 or closing the intended typed transport.
C2 remains the separately authorized composed execution boundary. P0-1 is a
`PATCH-CANDIDATE`, P0-2
selected `GOVERNED-PATCH` without admission, and P0-3 has a backend-independent candidate plus an
exact public governed libkrun source merge. That merge fixed two local console lifecycle defects,
added bounded console/raw-FD library tests, and materially improved measured coverage, but retained
uncovered code, a post-merge governed-branch/verifier mismatch, a pending Linux-arm64 unit job,
absent independent review, and no real transport, launcher, guest/VMM, or installed composition.
P0-4A passed only the conditional no-host-root topology. Session, MDM,
power-loss, independent-builder, clean-host, and Linux-worker cases remain later validation work
rather than reasons to delay backend-independent contract implementation.

Run bounded, disposable prototypes in parallel where practical:

- Go/Swift/TypeScript canonical signing interoperability: record the RFC 8785/JWS failure and test
  the bounded deterministic-CBOR/COSE fallback.
- macOS XPC peer requirements, dynamic validation, Keychain/access-group isolation, and protected
  storage separation.
- Apple Container and direct Containerization network, filesystem, resource, management-channel,
  orphan-recovery, identity, and teardown controls; a libkrun/HVF native follow-up; and an
  OCI/gVisor contingency harness.
- Broker-to-Supervisor immutable content-handle transfer without daemon content access or live user
  mounts.
- Supervisor language, per-user versus privileged process model, and need for any tiny launcher.
- Prepared-update, epoch-finalization, interruption, and repair state transitions.

Exit evidence:

- Every spike has recorded setup, exact platform/tool versions, adversarial cases, results,
  limitations, and a decision.
- Daemon key/content/backend access attempts fail in the macOS prototype.
- Stale same-team access to replacement operational keys is either denied by a retained mechanism
  or explicitly blocks the affected update/rotation posture.
- Exact enforceable backend controls and unsupported limits are known.
- Cryptographic implementations agree or a documented alternative is selected.
- Resulting ADRs and contract vocabulary are updated before schema freeze.

Recorded outcome: bounded CBOR/COSE, macOS authority separation, release-key transitions, installed
per-user services, content custody, and trust-transition ordering passed conditionally. Both stock
Apple Container and direct Containerization failed the production lifecycle gate. libkrun/HVF
remains the lead native candidate under evaluation, and its readiness tracks passed mechanics, but
the exact unpatched profile is blocked by mutable-path custody and an unexpected `NullFs` virtiofs
device. Retained follow-ups made FD-native custody and direct-block-root removal credible without
passing their final installed corpora. Guest transport/completion, installed distribution, and
release-byte admission remain open. The post-track P0 reconciliation proposes bounded console
ports for source/input and inline results, and defers filesystem-image output parsing until file
artifacts. OCI plus gVisor remains independent; only its surrounding OCI/runc harness has run, so
gVisor itself is unvalidated.

Gate C now permits freezing backend-independent identifiers, exact-or-refused limits, typed
admission and terminal classifications, and a fake backend that creates no guest. It does not
permit freezing libkrun paths/devices, arbitrary CPU or memory semantics, runner-exit success, the
current runtime manifest, or a stronger posture. See the
[Gate C implementation-readiness synthesis](GATE_C_READINESS_CHECKPOINT.md) and
[Gate C P0 reconciliation](GATE_C_P0_RECONCILIATION.md).

See [Feasibility Spikes](FEASIBILITY_SPIKES.md).

## Phase 2: contract and cryptographic freeze

Status: Phase 2A has implemented a passive, backend-independent foundation: a deliberately narrow
`JobProposal` candidate, minimum `ExecutionPlan` and `PlanRegistration` CDDL candidates,
byte-exact fixtures, and Go/TypeScript decoded views. Phase 2B now provides the closed conformance
manifest, integrity runner, and 82-rule/262-case/368-fixture corpus, including proposal/source/input,
exact plan/registration-state, passive approval/attempt Slice A, and unwired fixed-store Slice B
oracles. Slice C adds an `AttemptID`-keyed no-guest fake-lifecycle seam and 12 top-level focused Go
tests without changing the manifest counts. These slices are not frozen or activated;
the unwired slices now implement strict TypeScript raw/schema proposal decoding, all 18 retained
TypeScript semantic-resolution cases, exact TypeScript minimum-plan construction/encoding, strict
Go exact-byte plan/registration wrappers, all 40 retained Go registration-state cases, 44 passive
Go approval/attempt contract cases, and 12 fixed-store transition cases. Swift remains pending.
The fixed store now colocates registrations, approval/attempt authority, durable time high-water,
and lifecycle/effect checkpoints; it remains unwired. The fake lifecycle accepts only committed
created attempts through `AttemptResolver`, drives and recovers by `AttemptID`, and hard-codes
`FakeBackend.CreatesGuest() == false`.
[Proposed ADR-0025](adr/0025-colocate-durable-attempt-lifecycle-state.md) selects colocated lifecycle
records and effect checkpoints in that same versioned Supervisor transaction domain. Slices E1
through E5 now retain passive contracts, explicit fixed-store v1 migration/open validation,
durable lifecycle transactions, the FakeBackend-only driver, exact 256-active/4,096-retained
capacity behavior, and concurrent/repeated startup plus recovery-exhaustion evidence. The only
adapter used is the closed no-guest fake. A
focused local-only conformance handoff also carries copied TypeScript `ConstructedExecutionPlan`
bytes and complete role bindings into the real Go `registrationstate` component. Neither path is a
product-language/IPC seam or public consumer. E5 supplies no archive mechanic. Archive F1 plus its
passive format correction now supply
scope-separated global/segment indexes, typed hot/archive locations/counts, a distinct
migration-genesis checkpoint, generated answers, defensive copies, and eligibility selection only.
The follow-on valid-v1 mapping contradiction is now passively resolved with a lifecycle
absent/present union and independent attempt/lifecycle counts. The
[F2 v1 mapping resolution](SUPERVISOR_ARCHIVE_F2_V1_MAPPING_BLOCKER.md) retains the exact crash
witness and fault plan. The [stateful F2 result](SUPERVISOR_ARCHIVE_F2_MIGRATION_RESULT.md) now
implements the owner-asserted migration, downgrade refusal, and empty-archive full verifier with
exact known answers and local fault/corruption/capacity evidence, while preserving that witness
without recovery or invention.
Proposed ADR-0033 selects a pre-created
enrolled sibling object plus lifetime nonblocking BSD `flock` after one owned local
descriptor/process corpus. Passive owner-lock G1 supplies the internal Go/Darwin acquisition. G2
now supplies the owner-required current v1/no-guest startup composition, same-session
store/coordinator, sorted recovery, post-open fencing, and ordered close under owned temporary
roots. Proposed ADR-0038 now selects the one-shot Coordinator-authorized/Supervisor-created
bootstrap and its signed-object contract, but passive fixtures and the installed protected-state-
root/session/update matrix remain unimplemented. The first bounded G3 discovery stopped before installed build: certificate SHA-1
`1638CFBD9250A00B4DBD81AE8FD1C790B42F61E3` has display suffix `W4QUR9FUL4` but X.509 subject OU
and signed-byte TeamIdentifier `3DDR84M4JS`, and every cached profile is also 3DDR. Apple Membership
Details later confirmed 3DDR is the account Team and W4 is a member/display suffix. The cached
profiles still have the wrong App IDs. The exact test-only identifiers/bootstrap fields and
noncredential mismatch/update model are retained, but they do not advance installed or
protected-root evidence.
Proposed ADR-0029 selects the authenticated local IPC process/language topology and four-call
surface. Passive logical fixtures and the exact native dictionary prerequisite now cover three
methods without activating transport. Both remaining methods, `SubmitApprovalV0` and
`RequestAttemptV0`, remain `BLOCKED`; this slice assigns neither a native tag nor a deadline to
either method. The native harness, installed endpoints, production identities, consumers, and
platform evidence remain unimplemented. Production approval signing/verification, archive
implementation and production-engine selection, evidence composition, consumers, and atomic public
migration remain separate decisions.

Proposed ADR-0031 defines the reviewed archive/compaction semantics and conformance plan. Complete
expired registration cohorts may move to immutable retained segments
only after all bound attempts are durably destroyed with authoritative absence. Full records and
exact replay/non-reuse tombstones remain retained; referenced deletion is forbidden. A finite
fixed-store v2 checkpoint is selected only as the local oracle. Slice F1 and both passive F2
corrections remain the contract foundation. Stateful F2 writes the all-hot, empty-archive v2
migration successor and fully verifies it. Stateful F3 publishes and activates exactly one
immutable complete-cohort segment under the owner assertion and fully verifies the generation-two
successor. F4A adds read-only retained-global lookup/replay/passive-collision routing and excludes
archived attempts from hot recovery. The retained F4B blocker records the former effect-history
contradiction; F4B now passes atomic mutation and independent append-only effect tombstones in the
exact fixed-store scope. F4C now passes bounded second/later immutable-segment growth through the
exact 64-segment ceiling in its [retained result](SUPERVISOR_ARCHIVE_F4C_GROWTH_RESULT.md). F5
now passes manifest-last coherent backup, complete-copy verification, read-only exact-anchor restore
admission, bounded offline reporting, and explicit known-unreferenced orphan removal in the same
fixed-store oracle. F6 production-engine selection is the next archive slice.
Installed validation of the selected owner lock and power loss, coherent
restore/anti-rollback, continuous service, and all
consumer behavior remain open.

Accepted ADR-0034 freezes the first release as one byte-exact pass-through `main.mjs` member under
the existing plan-v0 source role, with no static/dynamic dependency request or module-loader
fallback. Its passive source-byte/SourceManifest foundation is retained, and an exact
division-versus-regexp counterexample continues to bar the removed scanner. The Source Validator's
passive V0 frames are exact, and V1 retains one unwired exact artifact plus supply-chain evidence,
but the product validator remains future defense-in-depth under ADR-0040. The passive S1/M2
registration/fetch fixture and facade cutover is now `PASSED` without waiting for R3-R5B. R2 retains exact unsigned role-specific bundle/parser construction with inactive
predecode/refusal and no spawn. The exact historical V2 macOS checkpoint is `BLOCKED`: its strict
bootstrap cannot lower `RLIMIT_AS`, the explicit unbounded diagnostic mutation retains ambient
file/socket/write authority, and Apple's supported App Sandbox child entitlement shape changes the
fixed V1 bytes. Resume only with a newly reviewed/enrolled artifact and supported exact memory/
confinement design; deprecated custom sandboxing is not a fallback.
The supported replacement research and R0 architecture slices pass their exact questions while the
product remains blocked. Direct inherited helpers are `NO_GO`. ADR-0036 selects two role-specific
private App-Sandboxed XPC launchers and matching fresh parser children, with no shared service,
result, cache, container, group, key, or accepted profile. Each private container is residual
scratch authority only, with mandatory cleanup/residue testing that is not a confidentiality
proof. The public footprint setter returned `KERN_NO_ACCESS`, so the accepted contract is a later
evidence-derived reactive footprint watermark with one direct child per launcher request, bounded
combined two-role concurrency, fixed sampling, and kill/drain/reap—not a hard peak/exact cap or
host-availability guarantee. R1 now passes with those measurements explicitly inactive and zero.
R3 also passes its exact Apple Development signed, installed, inactive-policy composition. Exact
R4-v1 candidates are `NO_GO`; R4-v2 was not executed, product R4/R5 remains `BLOCKED`, and R4 alone
may derive threshold/cadence/baseline/overshoot/kill-latency values.
Supervisor custody, Broker rendering, and runtime no-loader evidence remain unimplemented.
Accepted ADR-0035 selects exact Oxc 0.140.0 as the engineering candidate for a separate
one-shot Source Validator after a bounded parse-only comparison. V0 observes the fixed typed
protocol and cross-language passive conformance. V1 observes exact build bytes, graph, licenses,
SBOM, unsigned provenance, same-host reproduction, and V0/M1 behavior without enrollment. Later
slices must prove independent reproduction and signed enrollment, disposable OS profile,
independent role-private daemon/Broker invocation, grammar/mutation corpus, resource/residue
evidence, and fault recovery. R1 passive contracts/fixtures, R2 unsigned construction, and R3
signed installed inactive-policy composition are now retained. R2's role-specific bundles and
parser children reproduce offline across two clean same-host directories; their exact inactive
policies predecode and refuse without spawning. R4-v2 confinement/resource/residue, R5D daemon
consumer, and R5B Broker consumer are future conditional work; they are not predecessors
of the now-passed passive M2/S1 fixture/facade checkpoint. Unsupported private-XPC
reachability, widened authority/native loading/network/filesystem escape, orphan/cleanup failure,
mixed-version acceptance, or unacceptable measured host risk stops the exact candidate. No product
validator or runtime enforcement exists.
Proposed ADR-0032's TypeScript Source Preparer and immutable source-store topology is `BLOCKED` as
a conditional later feature. Its protected-store, worker, genesis/update, retention/release,
recursive field-authority, lifecycle, installed evidence, and plan-v1 cutover do not block
first-release plan/IPC/runtime work.
The dormant `SupervisorCore`
scaffold was removed in PR #49 under ADR-0027. See the
[Phase 2A parallel-review synthesis](PHASE_2A_PARALLEL_REVIEW_SYNTHESIS.md) and proposed
[Phase 2B boundary decisions](PHASE_2B_BOUNDARY_DECISIONS.md).

- Replace the mixed `Job` schema with narrow `JobProposal` semantics.
- Add schemas for plan, registration, approval, attempt, trust snapshot, integrity assessment,
  transcript, artifact manifest, agent summary, and composed receipt.
- Add closed CDDL and byte-exact fixtures for every canonically registered or signed internal
  security object; do not generalize from the first `ApprovalGrant` candidate.
- Add a machine-readable field-authority manifest for each target object. Every field identifies
  its origin role, validator/resolver, authority effect, approval visibility, content/guest-control
  status, and cryptographic or durable binding; verification rejects unclassified fields.
- Define semantic source-path canonicalization and logical input/output slots.
- Freeze strict raw decoding, canonical bytes, digest, signature, type/domain separation, and safe
  numeric rules using retained cross-language fixtures.
- Define stable error, violation, posture, lifecycle, and recovery records.
- Retain the completed ADR-0031 effect-history correction,
  [F4B result](SUPERVISOR_ARCHIVE_F4B_MUTATION_RESULT.md), and
  [F4C result](SUPERVISOR_ARCHIVE_F4C_GROWTH_RESULT.md), and
  [F5 result](SUPERVISOR_ARCHIVE_F5_BACKUP_RESULT.md), retaining full closed cohorts and exact
  tombstone indexes without referenced deletion, then compare a pinned production-engine candidate against the same logical,
  corruption, locking, backup, APFS, and power-loss corpus.
- Retain the completed bounded production CBOR/COSE dependency comparison: it selects pinned
  `fxamacker/cbor` only for future object-specific typed encode/decode while keeping Capsule
  predecode/caps/canonical/binding/replay controls, and rejects `go-cose` as a product envelope
  dependency. The narrow I2B1 production-shaped Go/Swift same-byte wrapper review is `PASSED` for
  checked-in public-key vectors only. Before broader root-module admission, freeze each remaining
  object set and close caller/key authorization, live signing, durable replay, installed consumers,
  and independent review.
- After F2 and G2 close their active logical/ownership work, run the single bounded SQLite
  production-engine comparison from the reuse map. Do not turn the fixed snapshot into an
  unbounded product store or select a driver before the complete graph and C/native boundary are
  measured.

Exit evidence:

- JSON Schema, Go, Swift, and TypeScript agree on all applicable fixtures.
- Duplicate keys, unknown fields/versions, unsafe numbers, unsupported powers, wrong types, and
  cross-protocol substitutions fail closed.
- Every frozen field has a closed authority classification, and adding a known schema/CDDL field
  without that classification fails repository verification.
- The backend contract requests only controls supported or explicitly rejected by candidates.

## Phase 3: registered-plan and fake-backend lifecycle

The ordered consumer-to-Supervisor dependency and the distinction between retained mechanics and a
working alpha path are summarized in the
[alpha vertical-flow readiness map](ALPHA_VERTICAL_FLOW_READINESS.md). In particular, the current
diagnostic HTTP server is not already the selected submission boundary, `RegisterPlanV0` remains
fresh-registration behavior, and `FixtureVerifier` remains a retained-vector test oracle rather
than a production approval verifier.

- Implement daemon plan generation and Supervisor plan registration.
- Before activating the candidate agent endpoint, define and enforce the daemon's aggregate
  connection, concurrency, in-flight-byte, queue, deadline, cancellation, downstream-stall, and
  overload envelope; retain maximum-size concurrency and slow/partial-client evidence for
  `DAEMON-001`.
- Implement direct Broker fetch/render/user-presence approval.
- Prove that UI activation or synthetic input alone cannot sign a grant without the configured
  LocalAuthentication/Keychain user-presence operation. Exercise Accessibility, overlay/focus, and
  stale-session cases under their explicit elevated-adversary limitations for `UI-001`.
- Implement a locally seeded, signed development `TrustSnapshot`; production TUF service remains
  later work.
- Implement a durable atomic grant ledger and one-attempt semantics.
- Build a fault-injectable fake backend that never runs guest code.
- Implement multi-store saga/reconciliation states and crash injection at every side-effect edge.
- Produce a Supervisor enforcement transcript and composed receipt without guest content.

Exit evidence:

- The daemon cannot execute unregistered or replacement bytes.
- Plan A approval cannot execute plan B, and one grant cannot create two attempts.
- Approval cannot add or widen a destination, audience, limit, capability, content reference,
  runtime/profile choice, or backend requirement; a selectable option must already be a closed,
  Supervisor-validated part of the registered plan.
- The daemon cannot forge ordinary terminal success.
- Authenticated-client overload stays within the configured daemon envelope, sheds work without
  authority change, and does not create unbounded queues or diagnostics.
- Approval evidence distinguishes the key-gated user-presence operation from UI clicks, focus, and
  unverifiable claims about comprehension or elevated automation resistance.
- Every post-create path reaches explicit destroy, unresolved, or quarantine state.

## Phase 4: inline JSON and content separation

- Implement bounded inline JSON content ownership.
- Implement fixed logical slots and Supervisor staging verification.
- Add a bounded JSON output gate and user-only storage.
- Return a fixed-shape agent summary with no guest-controlled strings, names, sizes, timings, or
  violation detail by default.

Exit evidence:

- Daemon/MCP credentials cannot retrieve user-only content.
- Agent-observable fields stay within the documented channel budget.
- Rich parsing is absent from the daemon and Supervisor.
- Keep the proposed `.mjs` parser outside daemon, Broker, and Supervisor address spaces; validator
  failure refuses, and runtime no-loader admission remains independently required.

## Phase 5: dependency-free runtime development execution

- Close the five reconciled P0 Gate C campaigns. Follow accepted ADR-0028 and bootstrap the real
  `Shrimpworks/deno` and `Shrimpworks/rusty_v8` governed branches from their exact retained upstream
  commits. Preserve the governed `deno_core` `PHYSICAL-OMISSION-PASS` as construction evidence only
  while closing packaging/provenance,
  restoration/backstop, ADR-0034's byte-exact `.mjs` source custody and no-loader evidence,
  Accepted ADR-0035/0036's disposable validator R1-R5B product-evidence gates,
  external-isolation, and profile-admission gaps before admitting a runtime. Carry the FD-native
  `PATCH-CANDIDATE` and direct-block-root `GOVERNED-PATCH` through independent review, closed
  APIs/routes, mutation tests, and composed final-profile reruns. Keep `RUNTIME-001` unsupported
  throughout.
- Before transport implementation, freeze separate exact source, canonical-input, completion-frame,
  and JSON-payload caps plus per-channel role/binding, length/digest, terminal-status, and commit-
  trailer semantics; continuously drain cap-plus-one and fail instead of resizing, depending on
  EOF, or inferring success from runner exit.
- Treat C2B v1, v2, v3, and v4 as immutable inputs. V3 binds the accepted governed heads and resolves
  runner/libkrunfw/kernel, descriptor, device, resource, and teardown semantics without rewriting
  earlier evidence. V4 binds the exact accepted header, current-source unsigned libkrun dylib,
  independently reviewed ABI, unsigned final runner, and new composed digest without execution.
  Reverify every archive, manifest, and artifact identity before use. Any fixed-owned-guest task
  requires separate authorization naming the v4 digest; neither passive consumption nor that
  experiment admits a runtime or profile.
- Patch or close the pinned virtio-console control/queue/descriptor and transmit hazards; define a
  distinct trusted launcher with a fixed child manifest and a host runner with an exact descriptor
  allowlist before any real-backend implementation.
- Carry the independently reproduced Go/Node P0-3 vectors into the selected host/launcher languages,
  close the remaining measured console gaps (2 functions/26 lines in `port.rs` and 14 lines in
  `process_tx.rs`), and run the exact frames through the governed directional transport, distinct
  launcher, closed child/runner descriptor manifests, guest, and forced-teardown lifecycle. The
  library/property/raw-FD, process-pipe, sanitizer, static-analysis, repetition, mutation, and CI
  build results do not cover those checks.
- Bind an exact guest-kernel image, configuration, boot/module/debug policy, provenance, and
  launcher restriction set into the runtime profile. Minimize unused facilities where supported
  and retain deliberate restoration tests, while continuing to require containment of a fully
  hostile guest kernel under `KERNEL-001`.
- Build one exact package-free runtime bundle.
- Execute one JSON-in/JSON-out job through the libkrun/HVF candidate in explicit development
  posture, including its durable-record-before-start lifecycle, inherited read-only root custody,
  and bounded console-port data path.
- Deny network and all ambient host resources using the proven mechanisms.
- Bind source, input, runtime, backend, controls, output, integrity evidence, and teardown into the
  attempt transcript.

The exact non-credential work first repairs the governed libkrun branch/verifier invariant without
rewriting the five retained patch bytes, then obtains independent human/CODEOWNER source/API review
of the exact merge, closes its measured uncovered code and hostile control/queue/descriptor surface,
carries forward the cross-language P0-3 fixtures, resolves runtime/TypeScript and build-provenance
decisions, and completes the final role, descriptor, source, SBOM, and minimum-OS build inputs. The
credential-dependent rerun requires a valid Apple Development or
Developer ID signing identity for the hardened App Sandbox runner, protected-container
construction, and Team-enrolled process identity; Developer ID signing/notarization credentials
for complete-bundle notarization, stapling, and Gatekeeper assessment; and installed-byte,
clean-host, session, and support-floor validation of those exact signed bytes. Those identities and
profiles must be deliberately authorized. The first exact G3 discovery contradicted the earlier
display-name inference: certificate SHA-1 `1638CFBD9250A00B4DBD81AE8FD1C790B42F61E3` is labeled
`Apple Development: Dylan Steele (W4QUR9FUL4)`, but its subject OU and emitted code-signing
TeamIdentifier are `3DDR84M4JS`. Apple Membership Details now confirms 3DDR is the account Team;
W4 is the common-name/member display suffix. A new Apple Development identity SHA-1
`80A4...D3793` is locally present but not authorized for use. All three Xcode 26.6-cached profiles
share Team 3DDR but have the wrong Gate B/wildcard App IDs and are not reusable. Local
signed/provisioned experiments require exact 3DDR role profiles, identifiers, and entitlements. A
separate Developer ID Application identity for Team `3DDR84M4JS` is later distribution authority
requiring explicit authorization and matching-Team package design; it does not make Developer ID/
notarization work current. Paid owned clean-host/minimum-OS
coverage is not currently planned and remains deferred activation/distribution evidence, not a
blocker for local F2/G1 mechanics. A genuinely independent
Linux/arm64 builder is viable but not currently planned, so same-host and GitHub-CI equality remain
limited and independent-builder equality is deferred. A real governed runtime/libkrun composition
also remains later work requiring explicit authorization for an owned disposable development
guest.

Exit evidence:

- Concurrent same-user mutation cannot change approved bytes observed by the guest, and the exact
  accepted device surface has a retained disposition and corpus. Runtime-root custody separately
  passes stable attachment identity, frozen-object construction, and adversarial end-to-end
  custody; `/dev/fd/N` alone does not pass.
- Guest success does not depend on VMM exit status, and no inline result is released without a
  valid attempt-bound frame, bounded JSON validation, and complete terminal evidence.
- The exact runtime refuses every prohibited subprocess, FFI, native-addon, inspector, macro,
  environment-file, and package-install path through a construction-level closure argument, source
  review, deliberate capability-restoration mutations, and the accepted adversarial corpus.
- Accepted ADR-0028 records the governed `deno_core` engineering selection and supersedes
  ADR-0003's Bun-first implementation choice; the runtime-neutral protocol remains intact. This
  exit item does not itself admit the runtime.
- The exact guest-kernel and launcher profile matches its reviewed manifest, exposes no
  undocumented kernel authority, and records remaining facilities and completion-trust limitations.
- The exact configuration passes the minimum development attack suite.
- Unsupported controls refuse execution rather than being silently approximated.
- The backend remains clearly labeled development until its exact profile validation record passes.
- No v0 profile depends on host-root execution, a separate-owner host service, or a privileged host
  helper; failure to close no-host-root custody rejects libkrun for v0 rather than silently
  expanding the boundary.
- An early installed harness informs topology, but final admission rebuilds the selected mechanisms
  and reruns every affected P0 gate on the exact signed/notarized bytes.

## Phase 6: regular-file snapshot vertical slice

- Implement native file selection and immutable regular-file data-fork snapshots.
- Transfer job-scoped handles to the Supervisor without daemon content access.
- Add bounded scratch/output storage and a disposable bounded filesystem-image parser before
  JSON/JSONL/text artifacts, then CSV.
- Add audience-controlled release and separate content grants.
- Complete CLI and MCP adapters over the same daemon protocol.

Exit evidence:

- No agent or guest contract contains an original host path.
- Link, mutation, special-file, oversized-input, formula/terminal/bidi, malformed-output, and
  audience-bypass tests pass.

## Phase 7: authoritative validation

- Compose the native libkrun/HVF candidate's independently observed storage, console, timeout,
  installed-recovery, hostile-guest, completion, runtime-authority, parser, and release controls in
  the applicable exact profiles.
- Inject process death, ENOSPC/I/O failure, corruption, and partial completion at every durable and
  external side-effect edge; add real APFS/power-interruption, sleep/wake, logout/login, reboot,
  fast-user-switch, locked-Keychain, update, and clean-host cases.
- Run the complete shared attack corpus against the exact Linux worker, engine, OCI, cgroup,
  `runsc`/shim, network, storage, output, and recovery configuration for independent comparison.
- Retain Apple Containerization only as a separately labeled development backend and regression
  target unless a future supported durable lifecycle API reopens its gate.
- Build runtime SBOM, two-builder provenance, corresponding-source/license publication, review
  attestation, registry, and backend validation records; require independent non-author approval
  for runtime, guest-kernel, VMM, and security-critical native patches; exercise deliberate
  capability-restoration changes, disable, revocation, update, repair, and explicit rollback.
- Define the exact platform and co-residency policy for each validated profile: hardware model,
  OS/hypervisor build, applicable vendor mitigation state, allowed concurrent guests/jobs, and
  explicit speculative/shared-cache residual risk. Do not claim microarchitectural
  noninterference from platform identity or configuration alone; track this as `PLATFORM-001`.
- Test runtime-integrity failure, cancellation, restart, orphan, parser failure, teardown ambiguity,
  repeated concurrency, cross-job state, and long-run cleanup paths.
- Measure startup, thermal/CPU/RSS, storage, I/O, descriptor/process leakage, and cleanup behavior.

Exit evidence:

- Only exact pinned configurations supported by retained evidence may use `validated-local` or a
  stronger isolation posture.
- Validation records and receipts bind the applicable host/co-residency policy and preserve the
  microarchitectural side-channel non-guarantee unless a separately reviewed exact campaign
  supports a narrower claim.
- Published claims include limitations and never collapse posture dimensions into one unsupported
  label.

## Phase 8: production trust repository and updates

- Complete I0-I4 from the macOS installation plan first: closed bundle/role contracts, protected
  Supervisor bootstrap, pairwise authenticated IPC, Source Validator launchers, and manual
  whole-bundle replacement with service re-registration and retained-state recovery.
- Operate root, targets, snapshot, timestamp, and delegated TUF roles.
- Publish releases, runtime bundles, review records, validation records, and Capsule-defined
  revocation/disable objects.
- Produce compact signed local trust snapshots outside the live Supervisor path.
- Support offline bundles and pinned self-hosted repositories.
- Implement explicit crash-safe install, update, repair, and key-replacement ceremonies.
- Keep Proposed ADR-0038's narrow bootstrap Trust Coordinator and any future update/repair methods
  separate from a reviewed Bundle Replacer ADR; neither replacement authority nor wider
  Coordinator power is implied by choosing a DMG or a TUF metadata format.

## Phase 9: optional Guardian and external witness

- Evaluate Endpoint Security entitlement and deployment requirements.
- Begin with notify-only observations; never make Guardian approval authority.
- Evaluate privacy-reviewed installation/receipt checkpoint witnessing and transparency monitoring.
- Send no job content, filenames, or per-execution identifiers by default.

## Deferred

- Node portability/contingency and additional Deno profiles
- network and API brokers
- secrets brokerage
- directory and repository snapshots
- rich document/media parsing sandboxes
- automatic background update delegation
- portable multi-device identity and recovery
- externally resolved DID methods and general Verifiable Credentials
- Windows
- hosted multi-tenancy
- Firecracker
- platform attestation
- public transparency or blockchain anchoring
