# `.mjs` Source Validator passive implementation boundary v1

Status: accepted R0 architecture boundary; passive R1 bytes and fixtures are not yet implemented.
This document defines the closed role/version/ownership surface that the next passive slice must
encode. It does not assign production resource measurements, create an endpoint, run a parser,
build or sign a binary, use an entitlement or credential, or admit a product control.

Decision: [ADR-0036](../adr/0036-select-role-separated-source-validator-launchers.md).

## Preserved historical identities

`capsule.source-validator.protocol/v0`, its request/result/candidate/artifact-profile fixtures, the
exact unwired V1 parser artifact, and the V2 direct-child checkpoint remain immutable historical
evidence. A v1 decoder must refuse every such frame/artifact/profile. The passive R1 generator adds
new targets; it does not edit, regenerate, resign, or relabel old bytes.

## Closed v1 nominal identity families

| Role | Daemon family | Approval Broker family |
| --- | --- | --- |
| Protocol | `capsule.source-validator.protocol/v1` | `capsule.source-validator.protocol/v1` |
| Containing bundle | `com.capsulecorp.capsule.daemon` | `com.capsulecorp.capsule.broker` |
| Method | `capsule.source-validator.validate-mjs-source.daemon/v1` | `capsule.source-validator.validate-mjs-source.approval-broker/v1` |
| Request | `capsule.source-validator-request.daemon/v1` | `capsule.source-validator-request.approval-broker/v1` |
| Result | `capsule.source-validator-result.daemon/v1` | `capsule.source-validator-result.approval-broker/v1` |
| Process profile | `capsule.source-validator.macos-xpc-parser-child.daemon/v1` | `capsule.source-validator.macos-xpc-parser-child.approval-broker/v1` |
| Artifact profile | `capsule.source-validator.artifact-profile.daemon/v1` | `capsule.source-validator.artifact-profile.approval-broker/v1` |
| Private service | `com.capsulecorp.capsule.source-validator.daemon.v1` | `com.capsulecorp.capsule.source-validator.approval-broker.v1` |
| Parser signing identity | `com.capsulecorp.capsule.source-validator-parser.daemon.v1` | `com.capsulecorp.capsule.source-validator-parser.approval-broker.v1` |

Both families bind `capsule.source-validator.reactive-footprint-policy/v1`, the existing
`capsule.mjs-source/v0` source profile, and the existing v0 `.mjs` media type. Reusing the source
profile means the semantic `main.mjs` bytes remain unchanged; it does not make the validator wire
or artifact profile backward compatible.

## Role-separated request ownership

Each parent constructs one fixed request in a fresh parent-owned buffer. Before XPC send it must
bind:

- exact protocol and role-specific method/request identities;
- a nonzero in-memory correlation value scoped only to the current role/service call;
- installation ID plus active trust-epoch sequence and digest;
- exact source byte length, source SHA-256, and one copied `main.mjs` byte string at the existing
  262,144-byte semantic maximum;
- exact matching role-specific artifact-profile digest; and
- exact accepted reactive-resource-policy digest.

The correlation value is not a registration, approval, attempt, nonce, replay, cache, store, or
authority identity. No request accepts a caller path, filename, environment, option, package,
loader, profile selector, diagnostic, service name, endpoint, descriptor, Mach right, key,
registration, Approval, backend, runtime, or guest value.

XPC must synchronously copy exactly one bounded data value before the parent buffer can be released.
The role-specific launcher owns its copy, predecodes before allocation growth, recomputes every
fixed identity/length/digest/profile/policy binding, and copies the request into one fixed child
pipe. The parser never borrows or retains parent or XPC memory.

## Role-separated result ownership

The parser writes exactly one fixed result frame. The launcher continuously drains through
result-cap-plus-one and treats the first diagnostic byte as failure. A candidate result contains
only:

- the matching protocol, role-specific method/result, correlation, installation, and epoch values;
- recomputed source length and digest;
- exact matching role-specific artifact-profile and reactive-resource-policy digests;
- closed parse, policy, and classification dispositions; and
- the five existing bounded grammar/CommonJS counts.

The launcher derives and verifies the allowed disposition from the fixed statuses/counts, requires
zero child exit, drains all pipes, reaps the direct child, establishes the required process-group
and cleanup state, and only then copies one result into a new XPC reply. The parent owns the reply
copy and repeats every decode, binding, and derivation check. It never accepts a launcher-owned
pointer, child prose, source echo, path, arbitrary string, profile selector, or result from the
other role.

The daemon result can authorize only continuation to plan construction inside the daemon's current
call. The Broker result can authorize only continuation to fixed rendering and the separately
controlled Approval-key operation inside the Broker's current call. Neither result is persisted or
reused. The Supervisor never receives either result.

## Resource-policy boundary

The R1 passive contract must define a fixed, versioned, domain-separated resource-policy record
whose active signed instance binds at least:

- role-specific process/artifact profile and supported host/build matrix;
- one direct child, no descendants, one active request per launcher, and the combined two-role
  concurrency assumption;
- request, result, diagnostic, connection, in-flight-byte, queue, CPU, wall, descriptor, file,
  process, core, and stack bounds;
- observed physical-footprint threshold and fixed sampling interval;
- measured clean baseline, maximum observed overshoot, maximum kill latency, and system-pressure
  disposition; and
- bounded kill, drain, reap, restart, update, startup, and residue-cleanup dispositions.

R1 may freeze field widths, domains, refusal rules, and inactive/missing-policy cases. It must not
invent an active threshold, sampling interval, baseline, overshoot, or kill-latency value. Those
values and the first active record/known answer are outputs of the later signed R4 corpus and
require profile review. Zero, a prose default, a V2 value, or an operating-system pressure policy
may not stand in for evidence.

## Required passive R1 targets

The next slice adds, for each role, closed request, result, process-profile, bundle/profile, and
consumer-projection targets plus the shared versioned resource-policy target. The canonical
field-authority manifest must classify every top-level and nested field in the same change. At a
minimum the classifications must distinguish:

- parent-origin copied source bytes from parent-computed evidence-only digest/length;
- installation/epoch/profile/policy values that narrow acceptance but grant no new authority;
- launcher-recomputed evidence from child-reported observations;
- role/service/method discriminators that select one predeclared operation;
- correlation-only values with no replay/store authority;
- fixed refusal/cleanup/resource observations from any planning or Approval authority effect; and
- fixture-only or unmeasured policy state from an active evidence-derived signed policy.

The generator must derive frame sizes and aggregate caps from canonical layouts, retain exact
boundary/cap-plus-one vectors, and emit separate known answers for both roles. Independent decoders
must reject cross-role, v0-as-v1, unknown/newer/older version, wrong installation/epoch/profile/
policy, partial/duplicate/trailing, reserved-field, and mixed-update cases. A shared result fixture
or a role field that can be rewritten after decode is a design failure.

The existing 228-field/20-target manifest and all V0 fixture counts remain unchanged until R1 lands.
This document is planning authority only and must not be listed as implemented field coverage.

## Passive stop conditions

R1 stops before implementation if the closed layout cannot express distinct role/service/profile
identity, copied request/reply ownership, active-epoch and resource-policy binding, strict
cross-role/cross-version refusal, or evidence-derived resource activation without a generic bus or
optional extension field. It also stops if official Apple public documentation/SDK evidence shows
the exact two private-service packaging cannot provide the required role-local reachability.
