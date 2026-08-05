# Governed `deno_core` C2B no-guest build-closure passive binding v2

Date: 2026-08-05

Work item: immutable canonical successor for reviewed C2B no-guest artifact closure

Scoped status: `PASSED`

Parent governed-runtime status: `IN_PROGRESS — TRENDING_GOOD`

C2B composed-profile/guest execution and runtime/profile admission: `BLOCKED`

Control evidence: `RUNTIME-001` and `VMM-001` remain `unsupported`

## Result and immutable lineage

The canonical successor is
`capsule.governed-deno-core-c2b-passive-binding/c1-c2a-no-guest-build-closure-v2`, schema version 2.
Its exact fixture is
[`passive-binding-v2.json`](../../schemas/conformance/c2b-governed-deno-core/passive-binding-v2.json),
7,115 bytes at SHA-256
`c59f7fdd27834dd7be5a05a3c44a973d6ffa99869b9b99e2531045926827190a`.

V2 consumes C1, C2A, and [C2B v1](GOVERNED_DENO_CORE_C2B_PASSIVE_BINDING.md) by exact immutable
byte length and digest. It does not edit, supersede in place, reinterpret, or make runnable the
9,289-byte C1, 26,850-byte C2A, or 8,221-byte C2B v1 objects. V1 permanently records its historical
draft-head/merge-dependency checkpoint. V2 records later reviewed evidence under a new identity.

## Pinned archive evidence

The only new evidence source is merged
[`Shrimpworks/capsule-experiments` PR #4](https://github.com/Shrimpworks/capsule-experiments/pull/4):

- merge `50108417ebf1aa45788a4e9a6b4ca6b4448e9972`;
- reviewed head `518eea04e1f81d27e61178b7f4ff164b955dea76`;
- shared reviewed/merge tree `cfe9865998bdd1ec556a0424ed30abdf51bb89d9`;
- predecessor merge `22b9eb2e92d17398e2844ad122e6c28faaf3a678`;
- retained root `experiments/gate-c-c2b-no-guest-artifact-closure`.

Independent passive readback verified the merge parent/tree, the archive verifier, both manifest
self-digests, and every manifest-bound evidence digest. The retained runtime-bundle candidate is
8,845 bytes at SHA-256
`d37c9311cf21e87cf693594ebb6bbf6c29bcb50d13c3f8a5e8334a0f02d30607`; it remains a build-only,
unsigned, unadmitted candidate and is not canonical runnable authority.

## Bound constructed identities

| Role | Mode | Bytes | SHA-256 |
| --- | --- | ---: | --- |
| governed libkrun dylib | `0755` | 4,426,736 | `f8e05177ce57a6f773f86d6755a29fe3f2bab92140dfe8caa33663a28584ae52` |
| libkrunfw dylib | `0755` | 24,339,104 | `0b14f4b8005dafd33c38df5935b9efdb6381c724224b3967ba1cecbecf10b6e9` |
| guest kernel | `0644` | 24,117,248 | `b50a4165215d5d897ab3614606a2105756cf8f2b2510cbceda9dc06057a5622d` |
| trusted init | `0755` | 930,144 | `4f4f2c8bc037c3226b183ad0d6daf35395c49467dfe5786d10a33290adf585cd` |
| trusted launcher | `0755` | 995,920 | `fd255394a26affadb1226d3f724494e76fc89785a5cced027a7bb9859d7da32d` |
| raw ext4 runtime root | `0644` | 134,217,728 | `390a4786a20d45f1c691ec8c203f84f5e9d372a30e98f867cc8309a144ca6798` |

The host preflight harness is separately bound as 34,864 build-only bytes at SHA-256
`4d56480bec2a34cd2d23f83f199e9eb62951ac66e53319ae45991af7b4000922`. Its disassembly has no
libkrun call. It is explicitly not a final host runner, has no final-runner identity, and grants no
launch/configuration/teardown authority.

## Exact nulls and refusal boundary

The following remain explicitly null and must refuse any invented value:

- final host-runner identity;
- separate runnable firmware identity, unless a later canonical decision proves the role
  inapplicable;
- composed-profile identity and digest;
- supported CPU-time, host-VMM-memory, and scratch-maximum mechanisms/values;
- guest-observed device, transport, root, child, trace, lifecycle, teardown, and restoration
  evidence; and
- runtime-manifest, runtime-profile, and backend admission.

All process, runtime, profile, consumer, adapter, backend, VM, guest, credential, signing,
publication, release, and admission effects are false. No product package imports the Go passive
decoder, and the TypeScript surface has no product call site.

## Validation and next dependency

Independent Go and TypeScript decoders bind the exact known answer, reject cap-plus-one,
cross-version, unknown, missing, duplicate, trailing, substitution, reordered-artifact,
null-closure, consumer, and admission mutations, and return defensive copies. The field-authority
manifest adds 120 recursive v2 classifications with version-specific closed consumers. Schema
tests independently reject invented runner, limit, guest, and admission state.

The scoped no-guest build-closure slice is `PASSED`. The parent remains `BLOCKED` on the final host
runner, separate-firmware-role resolution, supported resource mechanisms and values, accepted
composed-profile identity/digest, and one separately authorized owned disposable guest run. That
future authorization must name the exact v2/profile/artifact identities and retained evidence
location. It is the only dependency allowed to produce the guest-only device, transport, root,
child, trace, lifecycle, teardown, and restoration evidence; this passive successor does not.

The immutable [v3 successor](GOVERNED_DENO_CORE_C2B_PASSIVE_BINDING_V3.md) later resolves the
firmware, runner, device, runtime, resource-field, teardown, and passive composed-digest semantics
without changing these v2 bytes. V3 still authorizes no guest and keeps missing current-source
libkrun and final-runner artifacts as typed blockers.
