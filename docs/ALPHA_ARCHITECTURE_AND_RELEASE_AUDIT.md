# Internal-alpha architecture and release audit

## Status and decision

Work item: full architecture, runtime, macOS, protocol/approval, and persistence audit synthesis  
Status: `PASSED`  
Scope: defensive read-only audits of the current repository and official primary sources, followed
by canonical decision reconciliation in this repository  
Parent: owner-only hostile-`.mjs` internal alpha  
Parent status: `IN_PROGRESS — TRENDING_GOOD`  
Product admission: `BLOCKED`

The five audits agree that Capsule's narrow internal-alpha hypothesis is viable. Accepted
[ADR-0040](adr/0040-freeze-owner-only-internal-alpha-posture.md) resolves their two material design
conflicts: the host Source Validator becomes post-alpha defense-in-depth, and the fixed store may
serve only a tightly bounded, disposable owner-only alpha before F6.

No audit executed a guest, signed or installed new code, activated a store, used a product key, or
admitted a runtime/profile. The audits passed because they answered their scoped questions; the
product remains blocked on implementation and composed evidence.

## Frozen internal-alpha shape

| Property | Internal-alpha decision |
| --- | --- |
| Host | One named owner-controlled Apple-silicon Mac; manual Apple Development installation |
| Workload | Exactly one byte-exact UTF-8 `main.mjs`, at most 262,144 bytes |
| I/O | Bounded inline JSON input/output; no file artifacts |
| Approval | Native Broker, fresh human-presence operation, one-use attempt-bound grant |
| Execution | One fresh disposable Linux/arm64 guest per `AttemptID`; concurrency one |
| Runtime | Exact governed `deno_core`/V8 profile with no loader or ambient host operations |
| Backend | Exact governed libkrun/HVF composition with sealed descriptors and fixed devices |
| Store | Bounded fail-closed fixed-store exception; no restore or continuity claim |
| Deferred | TypeScript, host AST validation, updater, public distribution, support matrix, F6 |

The first fixed benign guest is an earlier engineering checkpoint and is not the product alpha.
It uses one immutable known answer and may not accept arbitrary or user source.

## What is going well

- Exact `.mjs` bytes and `SourceManifest` fixtures, canonical CBOR scaffolding, plan codecs, local
  registration, approval/attempt, response-loss, lifecycle, owner-lock, and archive F2-F5 mechanics
  all have meaningful passed evidence in their stated scopes.
- Governed Deno, rusty_v8, and libkrun accepted lines are pinned; Linux/arm64 runtime construction,
  passive C1/C2A/C2B bindings, and no-guest artifact closure materially reduce composition risk.
- libkrun's governed control-port/descriptor-chain validation is merged and guest-regressed.
- The signed inactive Source Validator R3 composition and macOS installation I1B establish useful
  role/signing/topology evidence even though active parser and installed authority remain blocked.
- The architecture consistently keeps the daemon away from Approval keys, user-only content,
  backend launch, and authoritative evidence.

## Blocking implementation gaps

### P0 — before the first fixed benign guest

1. Freeze a successor composed profile from exact accepted fork commits and built byte identities.
2. Close the final runner, libkrun/libkrunfw/kernel role, init, launcher, immutable root, numeric
   descriptor and port map, exact device inventory, and supported resource fields.
3. Explicitly disable implicit vsock/TSI, implicit console/init where applicable, every loader and
   host operation, and V8 string code generation; restoration-test each control.
4. Prove result-frame commit-last semantics, process-tree teardown, and authoritative absence. EOF
   and exit zero are never success.
5. Run the one fixed known answer plus root, device, descriptor, transport, cap, crash, timeout, and
   teardown mutations in an owned disposable guest.

### P0 — before the hostile `.mjs` internal alpha

1. Replace—not extend—the legacy broad proposal with the single-`main.mjs` contract.
2. Atomically retain plan, bindings, manifest, and exact source; Broker fetches only those retained
   bytes and execution accepts identifiers only.
3. Implement authenticated method-specific local IPC and keep diagnostic HTTP read-only.
4. Implement fixed Broker rendering, fresh user-presence signing, strict authorized-key COSE
   verification, one-use attempt creation, and protected installed Supervisor state.
5. Activate the bounded fixed-store exception with exact refusal thresholds and no restore path.
6. Implement the real adapter, restart reconciliation, transcript/receipt compositor, fixed public
   summary, and the composed hostile/refusal/restoration corpus.

## Risks and unknowns

| Risk or unknown | Current disposition |
| --- | --- |
| Host validator child survives launcher death | R4-v1 candidates `NO_GO`; R4-v2 untested. Validator is post-alpha until a supportable contract exists. |
| Guest runtime can resolve imports or regain a host operation | Alpha stop condition; must be physically absent and restoration-tested. |
| Implicit vsock/TSI remains enabled without a network device | Alpha stop condition; explicit disablement and device readback required. |
| Separate firmware identity is unclear | Do not invent one. Prove whether pinned libkrunfw is the sole boot-kernel carrier. |
| Exact host/VMM memory and CPU-time caps are unsupported | Omit them from approved authority; retain only enforceable vCPU, guest RAM, wall, concurrency, and byte caps. |
| Same-user path/store replacement | Installed protected-root and descriptor-relative evidence still required; hostile same-UID host processes remain out of internal-alpha scope. |
| Backend effect during Supervisor crash | Real adapter must enumerate and reconcile a stable guest identity; ambiguous state stays unresolved and disables output/capacity release. |
| Store corruption or rollback | No restore. Corruption retires the installation; F6 is required when continuity or restore enters scope. |
| New governed fork promotion bypasses CI | Before the next promotion, require stable aggregate checks, generic branch filters, admin enforcement, and the no-rewrite runbook. |

## Minimum meaningful hostile corpus

The existing P0-3 vectors and C2A restoration mutations remain immutable subcorpora. The composed
profile adds uniquely identified cases covering:

- static/export-from/dynamic/computed imports, `import.meta`, `eval`, `Function`, WebAssembly,
  inspector/worker/native-op discovery, and built-in restoration;
- filesystem, environment, process/exec, FFI, network, Unix socket, vsock/TSI, and inherited-FD
  discovery;
- source/input port swaps, wrong attempt/profile, truncated/cap-plus-one frames, backpressure,
  partial/zero progress, forged/early/missing completion, replay, and completion-FD access;
- writable-root aliases, path substitution, FD reuse/shared-offset mutation, restored NullFs or
  virtiofs, extra devices/mounts, and every artifact-digest mismatch;
- infinite loop, memory pressure, crash/panic, parent/reader death, cancellation, PID/identity
  mismatch, absence deadline, and restart recovery; and
- JavaScript throw/rejection/non-JSON results, malformed input, exact/max-plus-one result, and
  diagnostic flood/control strings.

The corpus is a minimum boundary suite, not an exhaustiveness claim. Every case must retain exact
plan/profile/runtime/root/attempt identities, transport observations, teardown/absence evidence,
and the expected terminal classification.

## External-alpha boundary

External alpha remains `BLOCKED` and deliberately later. It adds Developer ID, Hardened Runtime
distribution review, notarization/stapling/Gatekeeper, clean-host and minimum-OS/session matrices,
whole-bundle replacement/repair/uninstall, stale-signer and key rotation, F6/APFS/power-loss/
restore/retention evidence, governed release ceremony, support/privacy posture, and broader soak.
Automatic update is optional only if a safe manual replacement path closes rollback and stale-
service behavior.

## Evidence classification and sources

Repository fixtures and experiments support only their declared scopes. Official primary sources
used to check platform and format semantics include Apple's
[App Sandbox inheritance guidance](https://developer.apple.com/library/archive/documentation/Miscellaneous/Reference/EntitlementKeyReference/Chapters/EnablingAppSandbox.html),
[XPC documentation](https://developer.apple.com/documentation/xpc), and
[Hypervisor entitlement](https://developer.apple.com/documentation/BundleResources/Entitlements/com.apple.security.hypervisor);
[RFC 8949](https://www.rfc-editor.org/rfc/rfc8949.html),
[RFC 9052](https://www.rfc-editor.org/rfc/rfc9052.html), and
[RFC 9053](https://www.rfc-editor.org/rfc/rfc9053.html); SQLite's
[locking](https://www.sqlite.org/lockingv3.html), [atomic commit](https://sqlite.org/atomiccommit.html),
and [backup](https://www.sqlite.org/backup.html) documentation; the
[Virtio 1.3 specification](https://docs.oasis-open.org/virtio/virtio/v1.3/virtio-v1.3.html);
and official [libkrun](https://github.com/libkrun/libkrun),
[deno_core runtime options](https://docs.rs/deno_core/latest/deno_core/struct.RuntimeOptions.html),
and [V8 jitless](https://v8.dev/blog/jitless) material.

Confidence is high in the ordering, contradiction resolution, and current implementation gaps;
medium in the final firmware/scratch and installed host behavior until exact composed tests run.
