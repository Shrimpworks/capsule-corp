# Passive typed guest transport v1

Date: 2026-08-10

Work item: C5a passive typed source/input/completion transport

Status: `PASSED`

Scope: one dependency-free byte contract, deterministic repository fixtures, independent Go and
Node verification, and passive state/restoration models only

Evidence: 48 generated frame cases, 13 state/fault cases, 23 restoration cases, exact cap and
cap-plus-one fixtures, and byte-for-byte generator drift checks

Remaining work: C5b controlled transport execution, installed identity/composition, real
completion/absence evidence, runtime/profile admission, and product wiring

Next action: request separate C5b authorization naming the exact runtime/profile successor, owned
disposable guest and host, fixed fixture digest, process names, and evidence destination

Parent owner-only hostile-`.mjs` internal alpha: `IN_PROGRESS — TRENDING_GOOD`

Runtime/profile and product admission: `BLOCKED`

## Defensive result and authority boundary

This contract defensively freezes the bytes and refusal semantics selected by
[R2](../TYPED_GUEST_TRANSPORT_RESEARCH.md). It uses only generated files and local decoders. It
creates no endpoint, descriptor, process, listener, runtime adapter, VMM, guest, store record,
signature, installation state, or product consumer. The implementation is
`internal/protocol/typedguesttransportpassive`; no product package imports it.

The v1 object is a versioned successor to the historical 43-vector P0-3 candidate. It retains the
candidate's 152/160/64-byte shapes, big-endian integers, role and attempt bindings, exact-offset
commit, and narrowed C2A caps. The older external archive is not copied into this repository and
its 1 MiB source cap is explicitly rejected: ADR-0040 and C2A limit the only first-release source
to 262,144 bytes. V1 assigns its own role-specific magic and closed status numbers rather than
claiming identity with the historical unpublished byte strings.

## Exact streams and encodings

Each attempt uses exactly one frame on each of three distinct streams:

| role | magic | direction | header | payload | trailer | physical maximum |
| --- | --- | --- | ---: | ---: | ---: | ---: |
| source `1` | `CPSRC001` | Supervisor → launcher | 152 | 262,144 | 0 | 262,296 |
| input `2` | `CPINP001` | Supervisor → launcher | 152 | 262,144 | 0 | 262,296 |
| completion `3` | `CPCMP001` | launcher → Supervisor | 160 | 262,144 | 64 | 262,368 |

The completion reader retains at most 262,369 bytes while separately counting and draining the
entire stream. Once cap-plus-one is observed, refusal is irreversible even if the retained prefix
is a valid frame. EOF only terminates a stream; it is never a commit.

Source payload encoding is the exact already-registered UTF-8 bytes of the sole member named
`main.mjs`, with no BOM, archive, manifest, path, length prefix, normalization, newline rewrite, or
container inside the payload. The registered `SourceManifest`, source length, and source digest
remain independent Supervisor truth and must be supplied to the verifier; a self-consistent frame
digest cannot replace that binding.

Input and completion payloads use `capsule.canonical-inline-json/v0`: UTF-8 without BOM or
whitespace; literals `null`, `true`, and `false`; safe integers only in shortest base-10 form;
arrays in original order; object keys matching
`^[A-Za-z0-9][A-Za-z0-9_.-]{0,127}$` sorted by unsigned ASCII bytes; no duplicate key; no Unicode
normalization; and only quote, backslash, or lowercase `\u00xx` control escapes. Floats, `-0`,
unsafe integers, second documents, invalid UTF-8, and platform-native stringify variants refuse.

## Header and trailer layouts

All integers are unsigned big-endian. Offsets are zero-based and every unused bit is frozen to
zero.

| offset | bytes | source/input | completion |
| ---: | ---: | --- | --- |
| 0 | 8 | role magic | `CPCMP001` |
| 8 | 2 | protocol version `1` | protocol version `1` |
| 10 | 2 | method `1` | method `1` |
| 12 | 2 | role `1` or `2` | role `3` |
| 14 | 2 | header length `152` | header length `160` |
| 16 | 16 | nonzero `AttemptID` | same |
| 32 | 16 | nonzero distinct-domain `RegistrationID` | same |
| 48 | 32 | exact retained plan SHA-256 | same |
| 80 | 32 | exact selected runtime-profile SHA-256 | same |
| 112 | 8 | payload length | status u16, flags u16 zero, reserved u32 zero |
| 120 | 32/8 | payload SHA-256 | payload length u64 |
| 128 | 32 | — | payload SHA-256 |

The 64-byte trailer is `CPEND001`, protocol u16 `1`, method u16 `1`, role u16 `3`, trailer length
u16 `64`, the same 16-byte `AttemptID`, and SHA-256 of the exact 160-byte completion header plus
declared payload. It begins exactly at `160 + declared payload length`, is written after every
other completion byte, and is the final 64 bytes. Early, missing, truncated, changed, duplicate,
or trailing trailer data refuses.

## Closed completion result

| value | name | payload rule |
| ---: | --- | --- |
| 1 | `succeeded` | canonical typed JSON within both the fixed and approved cap |
| 2 | `workload-failed` | exact UTF-8 `null` |
| 3 | `result-invalid` | exact UTF-8 `null` |
| 4 | `child-terminated` | exact UTF-8 `null` |

No non-success payload carries guest prose, raw errors, host or guest paths, artifact names,
diagnostics, byte counts beyond fixed contract fields, or timing. Unknown values and nonzero flags
or reserved bytes refuse.

## Ordered frame refusal precedence

Independent decoders apply this order:

```text
FRAME_OVERSIZE → HEADER_TRUNCATED → MAGIC → PROTOCOL_VERSION → METHOD → ROLE →
HEADER_LENGTH → FLAGS_RESERVED → IDENTIFIER → BINDING → PAYLOAD_LENGTH_CAP →
FRAME_LENGTH → COMMIT_OFFSET → PAYLOAD_DIGEST → STATUS → NON_SUCCESS_PAYLOAD → JSON →
COMMIT_MAGIC → COMMIT_VERSION → COMMIT_METHOD → COMMIT_ROLE → COMMIT_LENGTH →
COMMIT_ATTEMPT → COMMIT_DIGEST
```

The physical count wins first so cap-plus-one is never reinterpreted as a valid prefix plus an
ignorable byte. `COMMIT_OFFSET` distinguishes a missing or truncated completion trailer after the
declared payload from an ordinary truncated payload. Source/input consumers additionally compare
the declared payload length and recomputed digest to independently retained Supervisor facts.

## Monotonic attempt state

```text
BOUND → ENDPOINTS_READY → RUNNER_READY → INPUT_TRANSFER → LAUNCHER_VALIDATED →
CHILD_RUNNING → RESULT_VALIDATED → TRAILER_WRITTEN → FRAME_OBSERVED →
TERMINAL_PROOF → DURABLE_COMMIT → COMPLETE
```

Reset, cancellation, timeout, stall, reader death, short-write error, wrong endpoint, or process
fault moves the affected machine to terminal refusal and never back to writable or acceptable.
Source and input must both validate before any child exists. The Supervisor starts every drain
before authorization, handles short writes without changing bytes, treats zero progress as a
deadline-bounded stall, closes each writer only after the entire one-frame write, and externally
tears down on cancellation independently of blocked I/O.

The Supervisor creates fresh dedicated `CLOEXEC` endpoints. Host runner FDs 5/6/7 map only to
source/input/completion; launcher FDs 3/4/5 have the corresponding roles. Object identity, access
mode, open-description aliasing, `CLOEXEC`, status flags, inherited manifests, and close-from
boundaries are pre-start canaries. The workload receives neither completion FD 5 nor a reopenable
node. Cap-plus-one output remains fully drained until producer closure or external teardown and
produces no partial authority effect.

The generated state cases freeze cancellation before start and during input/child/commit,
short-write, zero-progress, reader death, reset, launcher/runner death, valid-frame-before-absence,
response loss, indeterminate store outcome, and cancellation concurrent with commit. Only replay
after a completed durable commit returns the same immutable bytes; it never reruns a guest or
accepts replacement output.

## Completion-to-store projection

`FRAME_OBSERVED` is one bounded observation, not durable terminal truth. Publication consumes
independently established authority, source/input, runtime/profile, result, lifecycle, teardown,
authoritative runner and child-tree absence, cleanup-false, and publication facts. The guest frame
cannot assert or override those facts. Only the Supervisor-owned atomic completion transaction may
publish `committed-last`; EOF, exit zero, guest status, or a caller response cannot.

The 23 passive restoration cases cover endpoint substitution, alias/mode/flag/CLOEXEC inheritance,
wrong attempt/registration and digests, valid-prefix flooding, early trailer, EOF/exit substitution,
diagnostic-console substitution, implicit console/vsock/network/virtiofs restoration, unresolved
cleanup, changed response-loss replay, and corrupt durable commit. Each names the refusing control;
no mutation activates a product mechanism.

## Fixtures and verification

The deterministic corpus is in
`schemas/conformance/typed-guest-transport-v1`. The manifest records every file's byte length and
SHA-256, exact offsets, known bindings, state model, endpoint contract, refusal order, effects, and
case disposition. Its SHA-256 is
`79767a34a27bcc32a5f9a479b6a8737f9f5791447fa425ad83455546eadae235`. Exact maximum source,
input, and completion frames are retained, along with the 262,369-byte completion cap-plus-one
case.

```sh
node scripts/generate-typed-guest-transport-fixtures.mjs --check
node scripts/verify-typed-guest-transport-fixtures.mjs
go test ./internal/protocol/typedguesttransportpassive
```

These checks prove passive cross-language agreement only. C5b real console/launcher/VMM behavior,
installed identities, runtime/profile admission, process-tree absence, recovery, durable product
integration, and hostile-source execution remain `BLOCKED` and separately authorized.
