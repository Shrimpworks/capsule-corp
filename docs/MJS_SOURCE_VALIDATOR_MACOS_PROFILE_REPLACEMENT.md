# `.mjs` Source Validator supported macOS profile replacement

Date: 2026-08-04

## Work status

```text
Work item: supported replacement design for the blocked Source Validator V2 macOS process profile
Status: PASSED
Scope: official public Apple documentation and macOS 26.5 SDK interfaces, the exact retained V0,
  V1, and V2 artifacts, one controlled unprivileged resource probe, and architecture/conformance
  design only
Evidence or reason: the review identifies the supported privilege-separation composition, the
  controls it can and cannot enforce, the required new identities and bindings, and the exact
  blockers that prevent activation.
Remaining work: construct and sign no artifact yet; select the launcher topology, obtain signing-
  authorized evidence, resolve the writable-container and hard-memory postures, then run the
  complete replacement corpus.
Next action: an architecture owner must decide whether a distinct unprivileged validator-launcher
  role and a quantified reactive memory posture are acceptable. If both are accepted, a separately
  authorized signing task may build the newly identified artifact/profile and run the installed
  corpus. Otherwise select a stronger disposable parsing boundary in a new ADR.
Parent status: Product Source Validator V1-V5 and downstream M2/S1 remain BLOCKED. Proposed
  ADR-0035 remains Proposed and no product artifact, endpoint, consumer, runtime, backend, or guest
  is admitted.
```

The direct embedded command-line helper candidate is `NO_GO` in its exact scope: supported App
Sandbox inheritance gives it the launching daemon's or Broker's static sandbox rights, so it does
not create a lower-authority parser boundary. The Source Validator capability is not abandoned.
Its parent workstream remains `BLOCKED`, not `NO_GO`.

## Answer

There is no supported direct-child replacement that currently satisfies Capsule's V2 acceptance
rule.

Apple supports two materially different compositions:

1. A command-line tool launched directly with `posix_spawn`, `Process`, or fork/exec inherits the
   sandbox capabilities of its parent. The tool must be signed with exactly
   `com.apple.security.app-sandbox=true` and `com.apple.security.inherit=true`. This is suitable
   only when the child may hold the parent's static rights. It is not suitable for an Oxc parser
   launched from the daemon or Approval Broker.
2. An XPC service, login item, or background helper app can have a distinct App Sandbox. Of those,
   only a method-specific XPC launcher plus a fresh parser child plausibly preserves copied binary
   input/output and child-owned kill/reap semantics without a privileged helper. The XPC launcher
   would own no Capsule key, store, path, package, loader, backend, runtime, or guest authority and
   would spawn one fresh Oxc child for each call. The child would inherit the launcher's minimal
   sandbox rather than either authority parent's sandbox.

The second composition is a design candidate, not a selected profile. It adds a trusted launcher
role or two parent-private launcher instances, while a private XPC service is managed by `launchd`
and is available only to its containing app. That installed topology has not been authorized or
proven for both the daemon and Broker. App Sandbox also gives the sandboxed service a private
container with unrestricted read/write access, and the inherited parser child receives that
static access. No supported custom profile removes it. Finally, macOS supplies no usable
unprivileged hard per-process address-space or physical-footprint cap on the observed host.

Those are release blockers under the current exact-memory, no-store, one-shot acceptance rule.
Bounded input, output, concurrency, CPU, and wall time plus reactive parent monitoring materially
reduce exposure, but do not cap the peak memory a compromised child can consume between samples or
before termination. System memory pressure is a host policy, not a Capsule limit. The current v0
contract therefore may not call that composition memory-bounded or admit it.

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
| App-Sandboxed XPC service | Distinct, minimally privileged, `launchd`-managed service sandbox; private to the containing app | Only plausible supported launcher boundary; topology, caller reachability, lifecycle, and signed evidence remain `BLOCKED` |
| App Sandbox without network entitlements | Denies entitlement-gated incoming/outgoing network connections and sensitive resources | Required, but needs installed negative evidence for IPv4/IPv6/DNS/Unix/Mach/vsock-equivalent routes; socket creation alone is not network authority |
| App Sandbox container | Grants the sandboxed component unrestricted read/write access to its private container | Counterevidence to an absolute no-file/no-store claim; no supported zero-container or zero-durable-write mode was found |
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

The next conformance slice, if architecture unblocks it, must define all of these closed identities:

- `capsule.source-validator.protocol/v1` and
  `capsule.source-validator.validate-mjs-source/v1` for the revised enrolled profile. The fixed
  request and result layouts may retain the V0 lengths, but version/tag changes require new bytes,
  known answers, decoders, and cross-version refusal cases.
- `capsule.source-validator.macos-xpc-parser-child/v1` for the supported launcher/child process
  composition.
- `capsule.source-validator.artifact-profile/v1` for one complete signed artifact set. It is a
  new record type, not a relabeling of the V1 executable or V0 profile.

The v1 artifact profile must bind, at minimum:

1. exact Oxc engineering-candidate, Cargo lock, toolchain, target, build-flag, source, SBOM, notice,
   provenance, assessment, vulnerability-owner, SLA, and removal/upgrade identities;
2. exact signed parser Mach-O bytes, size, SHA-256, CodeDirectory hashes, signing identifier, Team
   identifier, designated requirement, architecture, minimum OS, SDK, UUID, and linked-library
   closure;
3. exact signed launcher/XPC-service executable and bundle bytes, `Info.plist`, placement, service
   identifier, launch mode, and containing-bundle role;
4. exact parent/responsible/self launch-constraint and library-constraint bytes and digests;
5. exact entitlements for each executable, including absence of network, app-group, Keychain-group,
   user-selected-file, JIT, unsigned-memory, DYLD, library-validation-disable, debugger, and
   `get-task-allow` exceptions;
6. Hardened Runtime, App Sandbox, notarization, staple, Gatekeeper, installation-manifest, trust-
   epoch, and active-profile bindings; and
7. exact process-profile values, request/result caps, deadlines, resource policy, and supported
   host/build matrix.

The result profile digest continues to identify one complete enrolled profile. An installation
may not combine a parser from one profile, launcher from another, or parent constraint from a third.

## Typed request, result, and process semantics

The replacement retains Capsule's fixed binary operation, not a generic XPC, JSON, Codable, or
RPC bus.

- Parent-to-launcher transport carries one closed method/version field and exactly one `xpc_data`
  value containing a single fixed Source Validator request frame. No path, filename, environment,
  package, loader, option, profile selector, diagnostic string, or authority-bearing identifier is
  accepted.
- The semantic source cap remains 262,144 bytes. If the v1 framing preserves the V0 layout, the
  minimum/maximum request sizes remain 80 and 262,224 bytes; that must be generated and independently
  verified rather than copied from this prose.
- Launcher-to-parser transport is exactly stdin/stdout/stderr after fixed `dup2` actions: one
  request, one result, and no diagnostic bytes. A valid result remains exactly 138 bytes only if
  the generated v1 layout proves it.
- The launcher recomputes request length/digest before spawn and result framing, correlation,
  source length/digest, profile digest, status/count relationships, and trailing-data absence after
  child exit. Parents repeat those checks over defensive copies.
- The child receives empty `envp`, fixed `argv`, a fixed cwd chosen by the launcher, reset signal
  state, a fresh process group, no inherited FD above 2, and no second request.

Success requires all of the following before one launcher reply exists: exact result bytes are
continuously drained through cap plus one; stderr stays empty; the child exits zero; no partial,
duplicate, trailing, or extra frame exists; every binding matches; and the complete child process
group is gone. EOF alone is never success.

On deadline, cancellation, parent disconnect, output cap, stderr byte, malformed result, crash,
signal, or unexpected child, the launcher sends `SIGKILL` to the child process group, continues a
bounded drain, waits for the direct child, discards all output, and emits only a fixed refusal. It
does not retry automatically. A later explicit request creates a new child. XPC interruption or
launcher death is a parent refusal; `launchd` restart cannot turn the failed call into success.

The daemon and Broker each enforce their own aggregate concurrency, in-flight byte, queue, and
deadline caps. A shared launcher process, parser cache, accepted prior result, or service retry is
forbidden.

## Memory and resource decision

The current contract requires an exact memory ceiling or refusal. No reviewed supported macOS
mechanism meets that rule:

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

Consequently, bounded parent input/output, one-child concurrency, one-second CPU, two-second wall,
reactive physical-footprint monitoring, kill/reap, and system pressure are **not sufficient for v0
under the current acceptance rule**. The failure mode is bounded in time and authority but not in
peak host-memory impact. A single compromised parse may cause collateral host pressure before the
monitor acts.

There are only two honest ways to resume:

1. select and retain a supported hard memory mechanism on every supported host; or
2. use a new ADR to change the Source Validator resource contract explicitly to a quantified
   reactive watermark with measured overshoot, concurrency, host-pressure, and availability
   limitations. That policy must remove every “hard/exact memory cap” claim, state that host
   availability is not bounded by the watermark, and be accepted before implementation.

Option 2 is not selected here. Until one option is accepted, memory remains a release blocker.

## Exact next implementation and conformance plan

### R0 — architecture and policy decision

- Decide whether one dedicated unprivileged validator launcher service or two parent-private
  instances fit ADR-0035 and ADR-0029 without an app group, Keychain group, temporary Mach lookup
  exception, or generic bus.
- Decide whether the launcher's empty writable App Sandbox container is acceptable residual
  scratch authority. If not, the App Sandbox process candidate cannot satisfy the no-store rule.
- Decide the hard-memory mechanism or explicitly approve the quantified reactive-policy ADR.
- Stop if the result adds privileged authority, shares parent/Broker stores or keys, or cannot
  close native loading, file, network, and child-process authority to the accepted residual scope.

### R1 — passive v1 profile contract

- Generate protocol/profile v1 request, result, candidate, bundle-manifest, and artifact-profile
  records with exact known answers and field-authority classifications.
- Preserve every V0 fixture and require cross-version, cross-profile, and V1-relabel refusal.
- Bind both allowed consumer roles and the selected launcher topology without granting either
  consumer access to the other's service, container, key, or store.

### R2 — unsigned construction only

- Rebuild Oxc from the exact reviewed lock with the v1 codec and no formatter, transformer,
  resolver, runtime, package, or loader.
- Build the smallest native launcher that only predecodes the fixed frame, applies the fixed spawn
  descriptor/rlimits, monitors, drains, kills, waits, verifies, and returns a fixed result.
- Retain no product consumer. Run offline source/SBOM/notice/provenance and two-builder review before
  any enrollment request.

### R3 — separately authorized signed installed experiment

- Use only explicitly authorized Apple signing identities and profiles. Sign the launcher/service
  and parser as new bytes; never mutate V1 or call the result V1.
- Verify exact entitlements, Hardened Runtime exceptions, launch/library constraints, nested-code
  sealing, notarization/staple/Gatekeeper, active Team/signing identifiers, and clean/minimum-OS
  hosts.
- Verify the XPC service's actual distinct sandbox, container, responsible process, private reach,
  and `JoinExistingSession=false` behavior for both consumer roles.

### R4 — confinement/resource corpus

- Repeat every V2 request/result/fault/kill/restart oracle.
- Add out-of-bundle/container and authority-store reads; container write/metadata/link/socket state;
  network connect/listen/DNS/IPv4/IPv6/loopback/Unix/Mach routes; Keychain; process/fork/spawn/signal/
  attach; dynamic library/framework loading; environment/cache; inherited FD/Mach/bootstrap ports;
  CPU/stack/file/FD/process/core limits; and parent/launcher/child death at every boundary.
- If reactive memory is accepted, measure supported-host Oxc baselines and adversarial allocation,
  exact sample cadence, observed threshold, maximum overshoot, kill latency, concurrent maximum-
  input calls, system pressure, and clean later calls. Never relabel the result a hard cap.

### R5 — independent consumers and updates

- Implement V3 daemon and V4 Broker integrations only after R0-R4 pass, using separate fresh child
  invocations and zero key/state effects on every refusal.
- Run mixed old/new consumer, launcher, parser, entitlement, constraint, OS, and trust-epoch cases.
  Unknown, missing, revoked, rolled-back, partially updated, or profile-mixed artifacts refuse and
  may require quarantine/repair.
- Then complete V5 grammar evidence. V6 runtime no-loader evidence remains independently mandatory.

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

This document selects no product profile. It does not authorize an XPC service, login item,
background helper, app group, signing operation, user identity, parser rebuild, daemon/Broker
consumer, Keychain access, runtime, backend, or guest. It does not prove App Sandbox denial,
Hardened Runtime/library constraints, clean-host compatibility, exact memory behavior, or product
admission. V0/V1 evidence remains scoped and V2 remains blocked.
