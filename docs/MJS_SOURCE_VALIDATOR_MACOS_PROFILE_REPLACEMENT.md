# `.mjs` Source Validator supported macOS profile replacement

Date: 2026-08-04

## Work status

```text
Work item: Source Validator R0 supported replacement architecture/resource decision
Status: PASSED
Scope: official public Apple documentation and macOS 26.5 SDK interfaces, the exact retained V0,
  V1, and V2 artifacts, one controlled unprivileged resource probe, and architecture/conformance
  design only
Evidence or reason: Accepted ADR-0036 selects two role-specific private App-Sandboxed XPC launcher
  services, accepts each private container as residual scratch authority under a narrowed no-store
  rule, and replaces the unavailable hard ceiling with an evidence-derived reactive footprint
  watermark and explicit availability limitations.
Remaining work: R1 passive contracts/fixtures and R2 unsigned role-specific construction are
  `PASSED`; continue with separately authorized R3 signing/install; R4 confinement/resource/residue evidence; R5 daemon then Broker
  consumers; and the later M2/S1 checkpoint.
Next action: separately authorize only R3 signing/install and private reachability for the exact R2
  bytes.
Parent status: Product Source Validator and downstream M2/S1 remain BLOCKED. ADR-0035 and ADR-0036
  are Accepted architecture decisions, but no product artifact, endpoint, consumer, runtime,
  backend, or guest is admitted.
```

The direct embedded command-line helper candidate is `NO_GO` in its exact scope: supported App
Sandbox inheritance gives it the launching daemon's or Broker's static sandbox rights, so it does
not create a lower-authority parser boundary. The Source Validator capability is not abandoned.
Its parent workstream remains `BLOCKED`, not `NO_GO`.

## Answer

There is no supported direct-child replacement for the historical V2 profile. Accepted ADR-0036
selects the supported replacement architecture, not a product profile.

Apple supports two materially different compositions:

1. A command-line tool launched directly with `posix_spawn`, `Process`, or fork/exec inherits the
   sandbox capabilities of its parent. The tool must be signed with exactly
   `com.apple.security.app-sandbox=true` and `com.apple.security.inherit=true`. This is suitable
   only when the child may hold the parent's static rights. It is not suitable for an Oxc parser
   launched from the daemon or Approval Broker.
2. An XPC service can have a distinct App Sandbox. ADR-0036 selects two method-specific private XPC
   launchers, one embedded for the daemon and one for the Approval Broker. Each owns its own fresh,
   role-specific parser child and no Capsule key, product store, path, package, loader, backend,
   runtime, or guest authority. No service, result, cache, container, or accepted artifact profile
   is shared across roles.

Private reachability and installed confinement remain unproven. App Sandbox gives each launcher and
child a role-specific writable private container, accepted only as residual scratch authority.
The supported replacement makes no zero-filesystem or confidentiality-through-cleanup claim.
macOS still supplies no usable unprivileged hard per-process address-space or physical-footprint
cap on the observed host. ADR-0036 therefore accepts only a quantified reactive watermark whose
threshold, cadence, baseline, overshoot, and kill latency come from the later signed corpus. It is
not a hard peak ceiling or host-availability guarantee.

## Preserved evidence and ownership

- Passive V0 and the unwired V1 artifact remain `PASSED` in their exact scopes. Their bytes,
  digests, profiles, results, and limitations are unchanged.
- The exact V2 checkpoint remains `BLOCKED`. Its `RLIMIT_AS=256 MiB` `EINVAL`, ambient-authority
  counterevidence, and entitlement mutation are not reinterpreted.
- The daemon invokes one fresh validator over its exact decoded copied bytes before planning.
- The Approval Broker independently invokes a fresh validator over exact copied bytes fetched from
  Supervisor-retained registration state before rendering or any Approval-key operation.
- The Supervisor retains, rehashes, atomically binds, and defensively returns exact source bytes.
  It does not parse or launch the Source Validator.
- No daemon result, launcher result, service lifetime, or parser cache can substitute for the
  Broker's fresh invocation.
- Runtime no-loader enforcement remains a separate V6 admission gate.

## Supported mechanism comparison

| Mechanism | Structural effect | Capsule disposition |
| --- | --- | --- |
| App Sandbox on a directly spawned embedded tool | Child inherits the parent's static sandbox rights; the child signature contains only App Sandbox plus inheritance | `NO_GO` for daemon/Broker direct launch because it exposes parent authority to hostile parser memory |
| App-Sandboxed XPC service | Distinct, minimally privileged, `launchd`-managed service sandbox; private to the containing app | Selected as two role-specific private launchers by ADR-0036; caller reachability, lifecycle, and signed installed evidence remain `BLOCKED` |
| App Sandbox without network entitlements | Denies entitlement-gated incoming/outgoing network connections and sensitive resources | Required, but needs installed negative evidence for IPv4/IPv6/DNS/Unix/Mach/vsock-equivalent routes; socket creation alone is not network authority |
| App Sandbox container | Grants the sandboxed component unrestricted read/write access to its private container | Accepted only as role-private residual scratch authority with mandatory cleanup/refusal evidence; no zero-filesystem, confidentiality, or secure-erasure claim |
| Hardened Runtime with no exceptions | Denies JIT, unsigned executable memory, DYLD environment injection, debugging/task-port access, and arbitrary third-party library loading by default | Required defense-in-depth; not a resource, filesystem, network, or syscall sandbox |
| Launch environment constraints | Kernel checks the signed child, parent, and responsible executable requirements before launch | Bind both exact daemon and Broker launch contexts, or the exact XPC launcher context, in every signed profile; signing-dependent |
| Library constraints | Kernel refuses dynamic libraries that do not satisfy the signed constraint | Bind the exact required dynamic closure and supported OS build; absence of `disable-library-validation` alone still permits Apple-signed libraries |
| `posix_spawn` with fixed absolute enrolled path, empty `envp`, `POSIX_SPAWN_CLOEXEC_DEFAULT`, fixed file actions, signal reset, and fresh process group | Fixes launch inputs and prevents unintended file-descriptor inheritance | Required inside a selected launcher; it does not create a different sandbox or apply rlimits atomically |
| Explicit descriptor closure | Only copied request, fixed result, and zero-byte diagnostic pipes survive to parser `exec` | Required and independently inventoried after `exec`; inherited Mach/bootstrap rights require separate sandbox/entitlement review |
| `RLIMIT_CPU` | Sends `SIGXCPU` after the fixed CPU budget | Retain the observed one-second candidate and rerun in the signed profile |
| `RLIMIT_FSIZE` | Prevents nonempty file growth past the limit | Retain zero bytes, but it does not prevent empty files, directories, links, sockets, xattrs, or metadata mutation |
| `RLIMIT_NOFILE`, `RLIMIT_NPROC`, `RLIMIT_CORE`, and `RLIMIT_STACK` | Bound descriptor count, child creation, core output, and stack extension in their documented scopes | Required; none substitutes for filesystem/network denial or total memory control |
| `RLIMIT_DATA` | Limits the `sbrk` data segment only | Reject as a total-memory claim; mappings and other allocator paths are outside that contract |
| `RLIMIT_RSS` | SDK alias of `RLIMIT_AS` | Same blocked interface; it is not a separate resident-memory control |
| `RLIMIT_MEMLOCK` | Limits locked memory | Useful hardening only; it does not limit allocated, resident, or virtual memory |
| `task_set_phys_footprint_limit` from the public SDK | A Mach interface nominally requests a physical-footprint limit | Controlled self-call returned `KERN_NO_ACCESS`; not an available unprivileged product mechanism |
| `proc_pid_rusage` physical-footprint observation | Parent can sample current and lifetime maximum physical footprint for a live child | Monitoring and evidence only. Sampling plus kill is reactive and may overshoot arbitrarily |
| Dispatch memory-pressure notification / system pressure termination | Reports or responds to system-wide pressure | Availability fallback only; not per-job enforcement and not an exact Capsule ceiling |
| Exported jetsam/spawn-memory symbols absent from public SDK declarations | Undocumented/private implementation surface | Rejected. Exported symbols are not supported API authority |
| `sandbox-exec`, `sandbox_init`, or custom Seatbelt profiles | Deprecated or unsupported custom profiles | Rejected and never an admission fallback |
| Endpoint Security or privileged helper | Additional observation or privileged lifecycle surface | Rejected by task scope and unnecessary for the supported candidate topology |

## Required new identities

No existing V1 or V2 identity can be reused. Signing, entitlements, bundle layout, launch/library
constraints, protocol changes, or a launcher executable each change the reviewed object.

The next passive conformance slice must define all of these closed identities:

- `capsule.source-validator.protocol/v1` plus separate daemon and Approval-Broker method, request,
  and result identities from the
  [passive v1 boundary](protocol/MJS_SOURCE_VALIDATOR_PASSIVE_BOUNDARY_V1.md). Version/tag changes
  require new bytes, known answers, decoders, cross-role refusals, and cross-version refusals.
- `capsule.source-validator.macos-xpc-parser-child.daemon/v1` and
  `capsule.source-validator.macos-xpc-parser-child.approval-broker/v1`.
- `capsule.source-validator.artifact-profile.daemon/v1` and
  `capsule.source-validator.artifact-profile.approval-broker/v1`, each binding one complete signed
  role-specific artifact set. Neither is a relabeling of the V1 executable or V0 profile.
- `capsule.source-validator.reactive-footprint-policy/v1`, whose active numeric evidence fields
  remain unset until R4.

The v1 artifact profile must bind, at minimum:

1. exact Oxc engineering-candidate, Cargo lock, toolchain, target, build-flag, source, SBOM, notice,
   provenance, assessment, vulnerability-owner, SLA, and removal/upgrade identities;
2. exact signed parser Mach-O bytes, size, SHA-256, CodeDirectory hashes, signing identifier, Team
   identifier, designated requirement, architecture, minimum OS, SDK, UUID, and linked-library
   closure;
3. exact signed role-specific launcher/XPC-service executable and bundle bytes, `Info.plist`,
   placement, service identifier, launch mode, containing-bundle role, and private caller route;
4. exact parent/responsible/self launch-constraint and library-constraint bytes and digests;
5. exact entitlements for each executable, including absence of network, app-group, Keychain-group,
   user-selected-file, JIT, unsigned-memory, DYLD, library-validation-disable, debugger, and
   `get-task-allow` exceptions;
6. Hardened Runtime, App Sandbox, notarization, staple, Gatekeeper, installation-manifest, trust-
   epoch, and active-profile bindings; and
7. exact process-profile values, request/result caps, deadlines, evidence-derived reactive resource
   policy, cleanup policy, and supported host/build matrix.

Each result profile digest identifies one complete enrolled role-specific profile. An installation
may not combine roles or take a parser from one profile, launcher from another, or parent constraint
from a third.

## Typed request, result, and process semantics

The replacement retains Capsule's fixed binary operation, not a generic XPC, JSON, Codable, or
RPC bus.

- Parent-to-launcher transport carries one closed role-specific method/version field and exactly one `xpc_data`
  value containing a single fixed Source Validator request frame. No path, filename, environment,
  package, loader, option, profile selector, diagnostic string, or authority-bearing identifier is
  accepted.
- The semantic source cap remains 262,144 bytes. If the v1 framing preserves the V0 layout, the
  minimum/maximum request sizes remain 80 and 262,224 bytes; that must be generated and independently
  verified rather than copied from this prose.
- Launcher-to-parser transport is exactly stdin/stdout/stderr after fixed `dup2` actions: one
  request, one result, and no diagnostic bytes. A valid result remains exactly 138 bytes only if
  the generated v1 layout proves it.
- The launcher recomputes request length/digest before spawn and result framing, role,
  installation/epoch, correlation, source length/digest, role-specific profile/resource-policy
  digest, status/count relationships, and trailing-data absence after child exit. Parents repeat
  those checks over defensive copies.
- The child receives empty `envp`, fixed `argv`, a fixed cwd chosen by the launcher, reset signal
  state, a fresh process group, no inherited FD above 2, and no second request.

Success requires all of the following before one launcher reply exists: exact result bytes are
continuously drained through cap plus one; stderr stays empty; the child exits zero; no partial,
duplicate, trailing, or extra frame exists; every binding matches; and the complete child process
group is gone. EOF alone is never success.

On deadline, cancellation, parent disconnect, output cap, stderr byte, malformed result, crash,
signal, watermark/sample failure, unexpected child, or residue, the launcher sends `SIGKILL` to the
child process group, continues a bounded drain, waits for the direct child, establishes the required
group absence, discards all output, cleans its container, and emits only a fixed refusal. It does
not retry automatically. A later explicit request creates a new child. XPC interruption or
launcher death is a parent refusal; `launchd` restart cannot turn the failed call into success and
startup cleanup must pass before a new request.

The daemon and Broker each enforce their own aggregate concurrency, in-flight byte, queue, and
deadline caps. Each launcher admits one active child/request; the signed corpus tests the combined
two-role maximum. A shared launcher process, parser cache, accepted prior result, or service retry
is forbidden.

## Memory and resource decision

The historical V2 contract required an exact memory ceiling or refusal. No reviewed supported
macOS mechanism met that rule:

- `RLIMIT_AS` and its `RLIMIT_RSS` alias returned `EINVAL` when lowered on the observed macOS 26.5.2
  arm64 host.
- `RLIMIT_DATA`, `RLIMIT_STACK`, and `RLIMIT_MEMLOCK` govern narrower resources and cannot be
  described as total memory limits.
- The public-SDK `task_set_phys_footprint_limit` self-call returned `KERN_NO_ACCESS` in the retained
  unprivileged probe.
- `proc_pid_rusage` can observe physical footprint, but a compromised child can allocate and touch
  memory between observations. A watermark-kill policy has an exact observed threshold and sampling
  interval but no exact peak ceiling.
- System memory-pressure and XPC lifecycle policies are OS availability behavior, not a user-owned
  per-validation resource limit.

ADR-0036 selects the second honest path: a quantified reactive watermark with measured overshoot,
combined role concurrency, host-pressure, and availability limitations. Every hard/exact/peak
memory claim is removed. A single compromised parse may still cause collateral host pressure
before the monitor acts, so host availability is not bounded by the watermark.

No threshold, sample interval, baseline, overshoot, or kill-latency value is selected here. R1 must
not invent an active policy. R4 derives those values from a separately authorized signed corpus on
every supported host; profile review must accept the measured risk before consumers activate.

## Exact next implementation and conformance plan

### R0 — architecture and policy decision

- `PASSED` through ADR-0036: two parent-private role-specific services, no shared bus/result/cache;
  private container accepted only as residual scratch authority with mandatory cleanup; and an
  evidence-derived reactive footprint policy with no hard-peak/availability claim.
- Private-XPC reachability, new bytes, signing, confinement, values, cleanup, and consumers remain
  later gates. Every ADR-0036 stop condition applies.

### R1 — passive v1 contracts and fixtures

- Generate protocol/profile v1 request, result, candidate, bundle-manifest, and artifact-profile
  records with exact known answers and field-authority classifications.
- Preserve every V0 fixture and require cross-version, cross-profile, and V1-relabel refusal.
- Generate separate role-specific frames/profiles and bind the selected two-launcher topology
  without granting either consumer access to the other's service, result, container, key, or store.
- Freeze copied ownership, cleanup/refusal fields, mixed-update refusal, and resource-policy field
  widths/domains. Do not invent an active threshold, cadence, baseline, overshoot, or kill latency.

### R2 — unsigned role-specific launcher/parser construction only

Status: `PASSED` for exact unsigned construction. Two role-specific bundle layouts, native
launchers, and parser children rebuild offline and compare byte-for-byte across two clean same-host
directories. The retained source/lock/license/notice/SBOM/provenance/static-closure evidence is not
independent-builder evidence. The launchers accept only the fixed role-specific `request` data
field and validate frame/source/policy bindings. Because the exact R1 policy is inactive, they
refuse without parser spawn; active descriptor/resource/monitor/drain/kill/reap/cleanup behavior
remains R4 after R3 installation evidence. No product consumer or Apple credential was used.

### R3 — separately authorized signing and installation

- Use only explicitly authorized Apple signing identities and profiles. Sign the launcher/service
  and parser as new bytes; never mutate V1 or call the result V1.
- Verify exact entitlements, Hardened Runtime exceptions, launch/library constraints, nested-code
  sealing, notarization/staple/Gatekeeper, active Team/signing identifiers, and clean/minimum-OS
  hosts.
- Verify each XPC service's actual distinct sandbox/container, responsible process, role-local
  private reach, no cross-role reach, and `JoinExistingSession=false` behavior.
- Follow the exact [R3 execution packet](SOURCE_VALIDATOR_R3_EXECUTION_PACKET.md). Team identity is
  reconciled to `3DDR84M4JS`. R3 later `PASSED` its exact signed, installed, inactive-policy scope;
  that result did not activate a parser or product consumer. Exact R4-v1 candidates are `NO_GO`,
  R4-v2 is unexecuted, and ADR-0040 moves the validator to post-alpha defense-in-depth.

### R4 — confinement, reactive-resource, and residue corpus

- Repeat every V2 request/result/fault/kill/restart oracle.
- Add out-of-bundle/container and authority-store reads; container write/metadata/link/socket state;
  network connect/listen/DNS/IPv4/IPv6/loopback/Unix/Mach routes; Keychain; process/fork/spawn/signal/
  attach; dynamic library/framework loading; environment/cache; inherited FD/Mach/bootstrap ports;
  CPU/stack/file/FD/process/core limits; and parent/launcher/child death at every boundary.
- If reactive memory is accepted, measure supported-host Oxc baselines and adversarial allocation,
  exact sample cadence, observed threshold, maximum overshoot, kill latency, concurrent maximum-
  input calls, system pressure, and clean later calls. Never relabel the result a hard cap.
- Test mandatory empty-inventory cleanup after every success/refusal, parser crash, launcher death/
  restart, update, and startup. Cleanup is not confidentiality or secure-erasure evidence.

### R5D — daemon consumer

- Implement only the daemon-facing consumer after R1-R4 pass. It uses one fresh role-specific
  service/child and fixed pre-plan refusal. It cannot reach or accept the Broker service/result.

### R5B — Approval Broker consumer and updates

- After R5D passes, implement the independent Broker-facing consumer over Supervisor-fetched exact
  bytes, with one fresh role-specific service/child and zero Approval-key operations on refusal.
- It cannot accept daemon results, profile digests, service state, or caches.
- Run mixed old/new consumer, launcher, parser, entitlement, constraint, OS, and trust-epoch cases.
  Unknown, missing, revoked, rolled-back, partially updated, or profile-mixed artifacts refuse and
  may require quarantine/repair.
- Then complete V5 grammar evidence and hold the M2/S1 checkpoint. V6 runtime no-loader evidence
  remains independently mandatory.

## Retained local observation

The smallest new probe is under
[`artifacts/mjs-source-validator-v2-replacement/`](../artifacts/mjs-source-validator-v2-replacement/README.md).
It compiles one unprivileged child that calls the public-SDK
`task_set_phys_footprint_limit(mach_task_self(), 128, ...)`. On macOS 26.5.2 arm64 the call returned
kernel result 8, `KERN_NO_ACCESS`, before allocation. It used no signing identity, JavaScript,
network, store, key, service, runtime, backend, or guest. No binary is retained.

## Official primary references

Reviewed 2026-08-04:

- Apple, [Embedding a command-line tool in a sandboxed app](https://developer.apple.com/documentation/xcode/embedding-a-helper-tool-in-a-sandboxed-app)
- Apple, [Discovering and diagnosing App Sandbox violations](https://developer.apple.com/documentation/security/discovering-and-diagnosing-app-sandbox-violations)
- Apple, [Creating XPC Services](https://developer.apple.com/library/archive/documentation/MacOSX/Conceptual/BPSystemStartup/Chapters/CreatingXPCServices.html)
- Apple, [Configuring the macOS App Sandbox](https://developer.apple.com/documentation/xcode/configuring-the-macos-app-sandbox)
- Apple, [Accessing files from the macOS App Sandbox](https://developer.apple.com/documentation/security/accessing-files-from-the-macos-app-sandbox)
- Apple, [Configuring the Hardened Runtime](https://developer.apple.com/documentation/xcode/configuring-the-hardened-runtime)
- Apple, [Applying launch environment and library constraints](https://developer.apple.com/documentation/security/applying-launch-environment-and-library-constraints)
- Apple, [Defining launch environment and library constraints](https://developer.apple.com/documentation/security/defining-launch-environment-and-library-constraints)
- Apple, [`setrlimit(2)`](https://developer.apple.com/library/archive/documentation/System/Conceptual/ManPages_iPhoneOS/man2/setrlimit.2.html)
- Apple, [`posix_spawn(2)`](https://developer.apple.com/library/archive/documentation/System/Conceptual/ManPages_iPhoneOS/man2/posix_spawn.2.html)

The installed macOS 26.5 SDK public headers additionally supplied the exact `spawn.h`,
`sys/spawn.h`, `sys/resource.h`, `libproc.h`, and Mach task declarations. Symbols found only in the
dynamic export table but absent from the public SDK declarations were treated as private and
rejected.

## Claim boundary

This document records the accepted R0 design and bounded R1/R2 evidence but selects no admitted
product profile. R2 creates unsigned bundle layouts and parser bytes only. It does not authorize an
installed/reachable XPC service, login item, background helper, app group, signing operation, user identity, daemon/Broker
consumer, Keychain access, runtime, backend, or guest. It does not prove private-XPC reachability,
App Sandbox denial, Hardened Runtime/library constraints, clean-host compatibility, reactive
memory behavior, cleanup, or product admission. V0/V1/V2 evidence remains scoped and unchanged;
product Source Validator remains `BLOCKED`.
