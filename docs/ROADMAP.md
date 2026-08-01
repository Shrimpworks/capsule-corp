# Roadmap

The roadmap is ordered by uncertainty and risk reduction rather than feature count. Disposable
spikes may be built outside the final product shape. Their retained evidence—not prototype code
quality—is the deliverable.

## Phase 0: architecture and claim baseline

- Align project, architecture, technical design, threat model, security policy, and ADRs on the
  daemon/Broker/Supervisor authority split.
- Define the component-compromise and control-to-evidence matrices.
- Inventory the target protocol objects and mark current schemas/types as pre-freeze scaffolding.
- State evidence and runtime-integrity claims without implying platform attestation.

Exit evidence:

- No document says the daemon can approve, launch, own user content, and author authoritative
  execution evidence.
- DIDs are consistently described as identifiers, not authorization roots.
- Apple Container and all runtime profiles are development posture.
- Current versus intended behavior is obvious.

## Phase 1: blocking feasibility spikes

Status: initial decision spikes and all five Gate C implementation-readiness tracks completed for
the currently available host/account environment. Gate C produced a required backend pivot and a
conditional native candidate, plus explicit blockers for the exact native profile. Session, MDM,
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
the exact profile is blocked by mutable-path custody and an unexpected `NullFs` virtiofs device. Guest
completion, installed distribution, and release-byte admission also remain open. The post-track P0
reconciliation adds stock-Bun authority closure, proposes bounded console ports for source/input
and inline results, and defers filesystem-image output parsing until file artifacts. OCI plus
gVisor remains independent; only its surrounding OCI/runc harness has run, so gVisor itself is
unvalidated.

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
byte-exact fixtures, and Go/TypeScript decoded views. They are not frozen or activated. The next
lane first closes the strict-decoder, conformance-corpus, and registration-semantics decisions
needed by the registered-plan/fake-backend lifecycle. Public cutover waits for that vertical slice,
an atomic consumer migration, and removal of the dormant direct-execution scaffold. See the
[Phase 2A parallel-review synthesis](PHASE_2A_PARALLEL_REVIEW_SYNTHESIS.md).

- Replace the mixed `Job` schema with narrow `JobProposal` semantics.
- Add schemas for plan, registration, approval, attempt, trust snapshot, integrity assessment,
  transcript, artifact manifest, agent summary, and composed receipt.
- Add closed CDDL and byte-exact fixtures for every canonically registered or signed internal
  security object; do not generalize from the first `ApprovalGrant` candidate.
- Define semantic source-path canonicalization and logical input/output slots.
- Freeze strict raw decoding, canonical bytes, digest, signature, type/domain separation, and safe
  numeric rules using retained cross-language fixtures.
- Define stable error, violation, posture, lifecycle, and recovery records.

Exit evidence:

- JSON Schema, Go, Swift, and TypeScript agree on all applicable fixtures.
- Duplicate keys, unknown fields/versions, unsafe numbers, unsupported powers, wrong types, and
  cross-protocol substitutions fail closed.
- The backend contract requests only controls supported or explicitly rejected by candidates.

## Phase 3: registered-plan and fake-backend lifecycle

- Implement daemon plan generation and Supervisor plan registration.
- Implement direct Broker fetch/render/user-presence approval.
- Implement a locally seeded, signed development `TrustSnapshot`; production TUF service remains
  later work.
- Implement a durable atomic grant ledger and one-attempt semantics.
- Build a fault-injectable fake backend that never runs guest code.
- Implement multi-store saga/reconciliation states and crash injection at every side-effect edge.
- Produce a Supervisor enforcement transcript and composed receipt without guest content.

Exit evidence:

- The daemon cannot execute unregistered or replacement bytes.
- Plan A approval cannot execute plan B, and one grant cannot create two attempts.
- The daemon cannot forge ordinary terminal success.
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

## Phase 5: dependency-free Bun development execution

- Close the five reconciled P0 Gate C campaigns: stock-Bun runtime-authority closure, immutable
  runtime-root custody, independent `NullFs` disposition, typed port transport/completion with
  bounded inline JSON, and an admissible complete installed development bundle.
- Before transport implementation, freeze separate exact source, canonical-input, completion-frame,
  and JSON-payload caps plus per-channel role/binding, length/digest, terminal-status, and commit-
  trailer semantics; continuously drain cap-plus-one and fail instead of resizing, depending on
  EOF, or inferring success from runner exit.
- Patch or close the pinned virtio-console control/queue/descriptor and transmit hazards; define a
  distinct trusted launcher with a fixed child manifest and a host runner with an exact descriptor
  allowlist before any real-backend implementation.
- Build one exact package-free runtime bundle.
- Execute one JSON-in/JSON-out job through the libkrun/HVF candidate in explicit development
  posture, including its durable-record-before-start lifecycle, inherited read-only root custody,
  and bounded console-port data path.
- Deny network and all ambient host resources using the proven mechanisms.
- Bind source, input, runtime, backend, controls, output, integrity evidence, and teardown into the
  attempt transcript.

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
  attestation, registry, and backend validation records; exercise disable, revocation, update,
  repair, and explicit rollback.
- Test runtime-integrity failure, cancellation, restart, orphan, parser failure, teardown ambiguity,
  repeated concurrency, cross-job state, and long-run cleanup paths.
- Measure startup, thermal/CPU/RSS, storage, I/O, descriptor/process leakage, and cleanup behavior.

Exit evidence:

- Only exact pinned configurations supported by retained evidence may use `validated-local` or a
  stronger isolation posture.
- Published claims include limitations and never collapse posture dimensions into one unsupported
  label.

## Phase 8: production trust repository and updates

- Operate root, targets, snapshot, timestamp, and delegated TUF roles.
- Publish releases, runtime bundles, review records, validation records, and Capsule-defined
  revocation/disable objects.
- Produce compact signed local trust snapshots outside the live Supervisor path.
- Support offline bundles and pinned self-hosted repositories.
- Implement explicit crash-safe install, update, repair, and key-replacement ceremonies.

## Phase 9: optional Guardian and external witness

- Evaluate Endpoint Security entitlement and deployment requirements.
- Begin with notify-only observations; never make Guardian approval authority.
- Evaluate privacy-reviewed installation/receipt checkpoint witnessing and transparency monitoring.
- Send no job content, filenames, or per-execution identifiers by default.

## Deferred

- Node and Deno portability
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
