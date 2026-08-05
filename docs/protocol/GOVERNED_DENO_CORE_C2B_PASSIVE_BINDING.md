# Governed `deno_core` C2B fixed-fixture passive binding

Date: 2026-08-04

Work item: C2B fixed-fixture build-candidate passive reconciliation

Status: `PASSED`

Scope: immutable, versioned, zero-effect binding of unchanged C1/C2A bytes to one governed
fixed-fixture fork supplement, build-evidence record, and closed artifact inventory

Parent governed-runtime status: `IN_PROGRESS — TRENDING_GOOD`

C2B composed profile/guest execution and runtime/profile admission: `BLOCKED`

Control evidence: `RUNTIME-001` and `VMM-001` remain `unsupported`

## Result and object boundary

The canonical object is
`capsule.governed-deno-core-c2b-passive-binding/c1-c2a-v1`, object type
`capsule.governed-deno-core-c2b-passive-binding`, schema version 1. Its exact fixture is
[`passive-binding.json`](../../schemas/conformance/c2b-governed-deno-core/passive-binding.json),
8,221 bytes at SHA-256
`3540d5224bdc81edbceafa1f0f17ac119904a70feab604957ab349dd116961a6`.

This is a new passive supplement. It neither replaces nor reinterprets the 9,289-byte C1 fixture at
`d5d75e638a15be6c9f4a3230d17309d085f6ec103a73b64d9e0fd656a5423c9e` or the 26,850-byte C2A
fixture at `d4ce88888186266f5d251e6246c889b1fd46d7746bb0ba56bcc4b3ce4675992f`.
Both predecessors are consumed by exact bytes, read-only, with mutation forbidden.

The binding has no registration, approval, attempt, runtime, profile, adapter, backend, VM, guest,
credential, signing, release, or admission effect. It does not define a general runtime interface,
loader, module path, source/input route, IPC method, descriptor, or admission flag.

## Immutable dependency inputs

The two draft-head inputs are mirrored byte-for-byte under
[`inputs/`](../../schemas/conformance/c2b-governed-deno-core/inputs/) so Go and TypeScript
validation is offline and does not trust a movable branch or live network response.

| Input | Exact identity | Exact bytes |
| --- | --- | --- |
| Fork supplement | `capsule.governed-deno-core.c2b-fixed-fixture/c1-c2a-v1`; Deno PR #2 head `29b71f06c2df5ab06721ccbb7bc744fb8104356e`, tree `172e57551fe5a6683f11c886a81f9634023a5514`, governed base `ea18b9dc21ff8ebd19347be7095f47937ee14ec2` | 1,895 bytes; SHA-256 `41350bcfc854338ded5e62f77475daf86486351356104dbbf647a8f8b5f11946` |
| Build evidence v2 | `capsule.c2b-fixed-fixture.runtime-build-evidence/c1-c2a-v2`; experiments PR #3 head `e016386ce6260dbca9f451cc07986fae24dfb334`, tree `cc03aca8841fa4d2f57eab3365faf6715e4eac34` | 6,026 bytes; raw SHA-256 `7de51e88ccc7cf35ee04b822d028f6c2184b72eaa833352fe9fc656b4a7baa13`; normalized self-digest `732301bf8553b0c59b3fe0e4f2b9e070dcc3a1b478e742dc13bd438873b7e488` |

Historical build-evidence v1 identity
`capsule.c2b-fixed-fixture.runtime-build-evidence/c1-c2a-v1` and self-digest
`6a673b88dc99e8939bc46ec88fb4f869caf7a9ff5909aa445e62afc5a3a83f87` remain immutable. The v2
record succeeds it because of the exact governed-fork formatter-policy commit; it does not edit or
erase v1.

Both upstream pull requests are explicit merge dependencies:

- [`Shrimpworks/deno` PR #2](https://github.com/Shrimpworks/deno/pull/2);
- [`Shrimpworks/capsule-experiments` PR #3](https://github.com/Shrimpworks/capsule-experiments/pull/3).

## Fixed bytes and artifact inventory

The exact source, canonical input, and completion bytes were read from the governed fork head and
independently matched to the retained C2A known answers. No transcription or newly invented
workload is accepted.

| Role | Bytes | SHA-256 |
| --- | ---: | --- |
| `main.mjs` source, including its trailing LF | 103 | `c8e940feb89b342de2d5e6bd13c413226676de9a539fce34c4107516e635b475` |
| deterministic-CBOR single-member SourceManifest | 89 | `712b1bd9739e4f6b0b027600207cbb08fb21b159a57bd34a15cf0ff8f32661b0` |
| canonical inline JSON | 36 | `9de0c909cfb111bd99c3b0b5f7a10972894270c2867022a71b6b6f3c0cd1af6e` |
| expected completion JSON | 35 | `bb7234ee486b0fbccc2091859ec93499e6a14ea7d6e091cdef60a0e2a6e8371c` |

The closed artifact array is ordered and has exactly six entries:

| Role | Bytes | SHA-256 | Additional count |
| --- | ---: | --- | ---: |
| Deno source archive | 32,352,414 | `7073152cccd4df42d5081ecec5c8ab36f8d6914039faa806060656d55a9e4cf3` | — |
| governed binary | 68,496,520 | `e781a90236cdf1272a9a16189c6be033164fa25a5aa9e52376ef998982ec0a77` | — |
| startup snapshot | 699,988 | `4e8965217d5a6675a880326eee6f690bbeec7e7cb243decf2f3e9f453a871a2c` | — |
| deterministic two-file bundle | 20,981,992 | `ad908b8289c86f25c3413713fa3e60c4c8bb91fec0d52763e870d7a186865ee6` | — |
| bundle manifest | 244 | `e3572e934b20a3fdfdd4b3abe2261d1348cd61bbc3523058fc2f5baac2724442` | — |
| connected acquisition bundle | 70,134,953 | `1e96e49a516e4cf6a9ec79acae9a9eb3d0ee52b332695fa11476a97e1e50d1d4` | 189 packages |

Same-host byte equality is retained exactly as observed. It is not independent-builder equality.
The fixed-fixture evidence did not construct the runtime root or a composed profile and did not run
a guest.

## Validation and field authority

The JSON Schema is closed recursively. Independent Go and TypeScript decoders bind the exact
binding, C1, C2A, fork-supplement, and build-evidence byte lengths and digests; reject duplicate,
unknown, missing, trailing, wrong-version, wrong-domain, wrong-length, digest, cross-link,
substitution, order, cap, dependency-state, and admission-state mutations; rederive the evidence
self-digest; cross-check the C2A/fork known answers; and return independently owned decoded values.

The field-authority manifest classifies all 143 top-level and nested fields. The passive fixture
generator supplies them; the schema plus Go/TypeScript validators validate them; the Capsule
repository conformance fixture retains them. Current consumers are only those validators. The only
later eligible consumer is a separately authorized composed-profile/owned-guest task after the
gate below. Caller path, argv, environment, arbitrary bytes, runtime selection, and admission never
become authority.

## Exact next consumer gate

This passive reconciliation is `PASSED` in its local scope. The next task remains `BLOCKED` until:

1. both dependency pull requests merge;
2. their exact merged heads and trees are reverified;
3. every bound artifact identity is reverified;
4. all C2A composed-profile blockers are closed; and
5. a separate task receives explicit authorization for the named owned disposable guest.

Only then may that separately authorized task consume this passive binding while composing the
fixed profile and running the already-fixed benign fixture. This task does not implement that
consumer. Any dependency head, tree, source, fixture, evidence, or artifact change requires a new
binding and evidence identity; in-place mutation is forbidden.

Even after that gate, runtime/profile admission remains a separate decision. `RUNTIME-001` and
`VMM-001` remain unsupported until their exact mechanisms and retained composed evidence close.
