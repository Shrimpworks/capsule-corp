# Governed `deno_core` C1 controlled-development composition

Date: 2026-08-04

Work item: C1 passive controlled-development composition contract

Status: `PASSED`

Parent governed-runtime status: `IN_PROGRESS — TRENDING_GOOD`

Runtime/profile admission: `BLOCKED`

## Scope and result

C1 fixes the smallest app-facing JSON-in/JSON-out composition for the intended governed
`deno_core` development profile. The executable fixture is
[`controlled-development-profile.json`](../../schemas/conformance/c1-governed-deno-core/controlled-development-profile.json),
9,289 bytes with SHA-256 `d5d75e638a15be6c9f4a3230d17309d085f6ec103a73b64d9e0fd656a5423c9e`,
validated by its JSON Schema, an exact-byte Go decoder, and an independent Node known-answer test.

This is a passive contract. It creates no process, runtime, backend, guest, descriptor, signature,
release, or admission state. It does not wire a Source Validator consumer. Every retained artifact
has `admitted: false`, and the profile refuses activation until C2 composition evidence and a
separate admission decision exist.

## Governed construction evidence

The contract pins the merged handoff from `Shrimpworks/capsule-experiments` PR #1 at merge
`fa03d7043b4f0653081d6c5733d597f49f6efd1c`. That handoff passed exact clean Linux/arm64
construction only; it did not admit or release a runtime.

| Role | SHA-256 |
| --- | --- |
| `rusty_v8` raw static archive | `e964d6b1b3689e91f8cf488d8a9f05764a03434b2e2e8347be5067300d39a7de` |
| `rusty_v8` gzip | `1ae209c9e4ba5803d010d2c79ee4cc0af0126c5a7ebcca211c7e41deaede4cd2` |
| governed `deno_core` binary | `56d3acefd2cc2f5136a0b8143c47131e49a58fbf66382dfd3e84f715ce8e2898` |
| startup snapshot | `4e8965217d5a6675a880326eee6f690bbeec7e7cb243decf2f3e9f453a871a2c` |
| two-file runtime bundle | `0cc08f93e82fcfe68b033e8807975a3bd67068a817da811a87a73aedc3f23937` |
| 22-entry root manifest | `100832dbb37737f29341bc5404df6d4405b8d6b706f274028892801fa88e7de8` |
| root tar | `9c46b45c4d220aedcc47c9ee53e875bc71d31d0b881b51740aaa9b882b5741e6` |
| root gzip | `e847651b35cd425dd8f6fe3bd45d693aff0af244e3a7bd30c629fa125cac62e8` |

The governed Deno head is `9adb0b68b55bca81644827f1e7749a3acb091bed`; the governed `rusty_v8`
head is `80e863ddb942a4aa2b384e794fc23e35b9d2bb15`. Fork source identity and
construction equality are inputs to C1, not substitutes for runtime integrity or admission.

## Exact application contract

1. The only source member is byte-exact `main.mjs`, 0 through 262,144 bytes, under
   `capsule.mjs-source/v0`. Its entrypoint, byte length, and SourceManifest digest bind to the
   exact registered ExecutionPlan fields. There is no transformation.
2. The launcher gives the source the fixed internal identity `capsule:main.mjs`. No other module
   exists. There is no static, dynamic, package, filesystem, network, or fallback loader; any
   module request refuses the attempt.
3. The primary-data input is at most 262,144 bytes of `capsule.canonical-inline-json/v0` and is
   parsed into a fresh value. Exactly one copied parsed value is passed to
   `globalThis.capsuleMain` with `globalThis` as its receiver.
4. Module evaluation must leave `globalThis.capsuleMain` as a function. It is invoked exactly once
   and its synchronous value or Promise is awaited. A missing binding, throw, rejection, or second
   invocation refuses success.
5. The returned value is serialized once with exact `JSON.stringify` UTF-8 semantics. Undefined,
   BigInt, cycles, serialization failure, invalid JSON, or output over 262,144 bytes refuses the
   result. The 262,144-byte value is the protocol ceiling; the exact lower-or-equal
   `ExecutionPlan.outputMaxJsonBytes` is enforced without clamping. The maximum P0-3 physical
   completion frame is 262,368 bytes and remains attempt-bound.
6. The workload never owns the completion endpoint. Success requires the launcher to commit the
   validated completion frame and the Supervisor to observe terminal runner lifecycle plus
   authoritative teardown. Runner exit or EOF alone is never success.

`capsule.canonical-inline-json/v0` is ADR-0023's restricted canonical value: null, booleans, safe
integers, strings, arrays, and objects only; no floats or out-of-range integers. Objects use the
bounded ASCII key grammar and unsigned-ASCII key order, arrays preserve order, strings preserve
decoded Unicode scalars without normalization, and the exact encoding contains no BOM or
insignificant whitespace. The registered canonical byte length and SHA-256 must bind the input
before it is copied into the invocation.

`capsuleMain` is a source-created post-evaluation binding. The fixture's `permittedGlobals` list is
the exact desired pre-source runtime global surface; its `forbiddenGlobals` list names ambient
surfaces that must remain unavailable. The lists, `--jitless` and `--random-seed=42`, the empty
extension registry, and the exact three built-in ops are all equality-checked. C1 deliberately
states the desired surface more narrowly than the construction handoff proved. C2 must reproduce
and adversarially verify it; a missing, extra, renamed, or re-enabled global/op/extension refuses
composition.

The workload has no readable or writable file. The 22-entry root is runtime custody identified by
the root-manifest digest above, not workload filesystem authority.

## Descriptor and resource roles

| Logical role | Owner and direction | Workload access |
| --- | --- | --- |
| `runtime-root` | Supervisor to host runner, read-only custody | none |
| `registered-source` | host writer to guest-launcher reader | indirect input only |
| `approved-inline-input` | host writer to guest-launcher reader | copied argument only |
| `completion-result` | guest launcher to host reader | none |

The handoff's construction probe observed descriptors 0, 1, and 2, but that is not a composed
profile manifest. C1 therefore leaves numeric assignment as `c2-required-unselected`; inventing
numbers here would create unsupported authority. C2 must select and prove a closed numeric
descriptor manifest, directionality, ownership, EOF behavior, shared-open-description behavior,
and teardown.

The transport reference is `capsule.gate-c.p0-3.measured-limits/v0`. Plan-owned wall time, console,
vCPU, guest RAM, and scratch references remain roles resolved before approval under ADR-0009. The
machine-profile reference is intentionally null: no exact governed runner/kernel/init/libkrun
composition or enforceable machine values have passed yet. Missing or unresolved values refuse;
they are never silently defaulted or clamped after approval.

## Refusal and stop boundary

The fixture contains the exhaustive C1 refusal-code vocabulary and stop conditions. Identity,
artifact, source, input, binding, serialization, globals/ops, loader, descriptor, resource, or
admission mismatch fails closed before activation. Work stops rather than extending C1 if it would
require inventing authority, changing RUNTIME-001, wiring a Source Validator consumer, admitting a
runtime/backend/profile, launching a guest, signing, publishing, or acquiring credentials.

## Evidence and limitations

- JSON Schema validation fixes structure and constants.
- The Go package fixes the complete file length and SHA-256, rejects unknown fields and trailing
  data, validates semantic equality, and returns defensive copies.
- The independent Node test fixes the same known answer and rechecks the security-relevant
  application, runtime, descriptor, resource, and zero-effect properties.
- Neither validator parses or executes JavaScript. Neither creates descriptors or demonstrates
  isolation. Passing C1 proves only the passive contract bytes and their internal consistency.
- The external handoff is pinned to an exact merge and artifact identities, but its retained
  limitations still apply: no independent-builder equality, governed release publication,
  admission, guest execution, or installed profile evidence.

## Next boundary: C2

C2 is a separately authorized controlled composition experiment in an owned disposable guest. It
must consume this exact C1 fixture; bind exact governed runner, snapshot, root, kernel/init,
libkrun, and launcher identities; choose the closed numeric descriptor and machine-resource
manifests; and reproduce the JSON-in/JSON-out, global/op/module/file, P0-3 transport, fault,
teardown, and restoration-negative corpus. P0-1, P0-2, and P0-4 remain separate required evidence.
Even a passing C2 experiment would not admit the runtime or profile; that remains a later explicit
admission decision after every required control row closes.
