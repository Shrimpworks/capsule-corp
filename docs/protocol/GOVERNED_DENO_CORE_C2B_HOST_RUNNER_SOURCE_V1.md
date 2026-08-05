# Governed `deno_core` C2B host-runner source contract v1

```text
Work item: C2B final host-runner passive source/contract boundary
Status: PASSED
Scope: dependency-free static C17 source contract and local Go mutation verifier only; no ABI
  build, final runner artifact, libkrun/HVF call, process, VM, guest, signature, installation, or
  admission effect
Evidence or reason: exact retained source bytes freeze Supervisor ownership, one runner per
  AttemptID, sealed inherited process/FD manifest only, FDs 0 through 7, raw-root FD 4, exactly
  three console ports, explicit implicit-console/init/vsock disable calls, closed call order,
  forbidden authority, fail-closed preflight, and Supervisor-owned teardown/absence. Mutations of
  cardinality, replacement values, FDs, root access, ports, devices, call order, teardown owner,
  artifact claim, or authority effects refuse.
Remaining work: retain the accepted libkrun header and current-source dylib locally; implement and
  review exact ABI calls; build final App-Sandboxed runner bytes; retain digest/call audit and a new
  materialized profile; obtain separate owned-disposable-guest authorization.
Next action: construct the final ABI implementation only after exact accepted header/current-source
  dylib inputs are retained locally and separately reviewed.
Parent status: C2B fixed-owned-guest eligibility and runtime/profile admission are BLOCKED.
```

This slice follows immutable C2B v3 and Accepted ADR-0041. It does not alter v3 bytes, fill v3's
final-runner artifact role, or reinterpret the earlier build-only preflight. Its only durable output
is a no-effect source contract plus a verifier that later runner work must preserve.

## Placement and language

The retained verifier lives in `internal/execution/hostrunnerpassive`, beside execution authority
mechanics but with no product consumer. The canonical source contract lives at
`schemas/conformance/c2b-host-runner-source-v1/capsule-host-runner-contract.c`.

C17 remains the narrow placement because the future role is one native App-Sandboxed VMM process
calling libkrun's C ABI. This slice deliberately does not copy or guess that ABI. Task scope barred
external accepted-fork retrieval, and the exact accepted header/current-source dylib are not
retained in this repository. Therefore the source records function names and semantic arguments,
not compilable prototypes or runnable calls. This preserves local-only scope and fails closed
instead of manufacturing a final artifact.

Source known answer:

- bytes: `3,996`;
- SHA-256: `7dbc8e5c21acd776fec653665517b6da1ecfe07278da0444d710c251744978b2`;
- predecessor C2B v3 fixture SHA-256:
  `d72327bba369484a56db7d543a32e8bbd4eac403230ac65d63709ac3ba3bbdfb`;
- predecessor passive-contract SHA-256:
  `8b1ec936a7b56370716d28557125e46866dea8f21a149704a01f251a0dddbcc1`.

## Authority boundary

Attempt identity and process cardinality remain Supervisor state. Runner receives no AttemptID,
plan, profile, path, backend flag, image, mount, or serialized configuration bytes. The sealed
Supervisor descriptor is the fixed inherited process and FD manifest, not an execute-time policy
object. Only the Supervisor may create one runner for one committed AttemptID, authorize start,
drain streams, interpret completion, cancel, signal, and prove authoritative absence.

Runner source contract owns only:

1. exact argv, empty environment, FDs 0 through 7, close-from 8, access-mode, and raw-root checks;
2. exact fixed libkrun configuration sequence;
3. ready then exact one-byte `G` plus EOF start handshake; and
4. entry into `krun_start_enter` only after every earlier step returns its exact expected result.

Any unexpected value refuses before the next step. Runner exit remains lifecycle evidence only.

## Static call and device contract

Ordered source obligations are:

1. validate one Supervisor-owned process per AttemptID;
2. validate exact argv/environment/FD manifest and no execute-time replacement values;
3. validate unlinked regular read-only raw-root FD 4, mode `0400`, fixed length, and fixed digest;
4. create context and set one vCPU/256 MiB;
5. explicitly disable implicit console, init, and vsock;
6. attach FD 4 as read-only raw `/dev/vda` and select fixed ext4 `ro,nosuid,nodev` remount;
7. add one multiport console and exactly ports `capsule.source`, `capsule.input`, and
   `capsule.completion` with IDs 0, 1, and 2;
8. select `hvc0`, `/`, and fixed init path/argv/environment;
9. complete the start handshake; and
10. call `krun_start_enter`.

Allowed device inventory remains balloon, RNG, one three-port console, and one read-only raw-root
block device. Network, vsock, TSI, virtiofs/`NullFs`, writable or additional block, host sockets,
live host paths, and kernel/firmware/legacy-disk/root path calls remain forbidden.

## Teardown and blockers

Teardown remains external. Supervisor revalidates PID, start time, UID/GID, executable path, and
code identity before `SIGKILL`, then requires authoritative absence within 1,000 ms and at most
1,200 ms from initial action. Identity mismatch is unresolved: no signal and no success.

Accepted libkrun identity stays upstream `728df812...c015`, accepted commit
`7432eda5...d632`, tree `7671440c...e346`. Those source identities do not supply a locally retained
header, current-source dylib, final ABI runner bytes, review, signature, or admission. A later slice
must create a new versioned materialization; it must not mutate this source contract or C2B v3.

## Immutable successor

The later [v4 materialized profile](GOVERNED_DENO_CORE_C2B_MATERIALIZED_PROFILE_V4.md) is `PASSED`
for its exact build/static scope. It retains the accepted header, twice-reproduced current-source
unsigned libkrun dylib, independent ABI audit, byte-equal unsigned final runner, and composed
digest. It preserves this source contract without changing its 3,996 bytes. Neither runner nor
libkrun was executed or loaded; fixed-owned-guest eligibility and admission remain `BLOCKED`.
