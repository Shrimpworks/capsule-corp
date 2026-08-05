# macOS installation I2B1 passive bootstrap request/record contract

Date: 2026-08-04

```text
Work item: I2B1 passive Supervisor bootstrap request/record boundary
Status: PASSED
Scope: closed deterministic-CBOR/COSE_Sign1 request and record shapes, nominal Go views, bounded
  passive verification, independently generated Node known answers, independent Go/Swift
  verification, recursive field authority, replay/time/refusal fixtures, and future transaction
  ordering only; no signer or installed ceremony is implemented
Evidence or reason: 71 repository cases contain 6 accepts and 65 refusals. Go and an actual
  CryptoKit-backed Swift verifier agree on every case and verify the independent RFC6979-P256
  Node signatures. Raw/calculate caps, canonical bytes, strict fields, signature shape, signer
  authorization, time, replay, request binding, and defensive copy ownership are enforced.
Remaining work: production wrapper review/fuzzing, authorized signing/key/container/service
  mutations, protected-root creation, owner/store composition, and installed fault evidence.
Next action: separately authorize I2B2 unsigned installation-only construction; do not activate or
  perform the I2B3 ceremony from this passive slice.
Parent status: installed I2B is BLOCKED on the separately authorized I2B2-I2B5 corpus.
```

## Exact identities and authority

The only payload object identities are:

| Object | Object/version | purpose | audience | media type |
| --- | --- | --- | --- | --- |
| request | `capsule.supervisor-bootstrap-request` / `0` | `capsule.installation.bootstrap.request` | `capsule.execution-supervisor.bootstrap` | `application/capsule.supervisor-bootstrap-request+cbor;v=0` |
| record | `capsule.supervisor-bootstrap-record` / `0` | `capsule.installation.bootstrap.record` | `capsule.execution-supervisor` | `application/capsule.supervisor-bootstrap-record+cbor;v=0` |

The on-demand Trust Coordinator is the only originator of both signed payloads and the only future
holder of the installation-root private-key reference. The Supervisor is the first owner of the
root/owner/store observation facts copied into the record; the Coordinator independently validates
that typed observation and constructs the record. The visible app, Broker, daemon, updater,
replacer, Approval grant, attempt authority, runtime admission, and Apple code signature cannot
originate or substitute either object.

The installation-root public key is one exact 77-byte deterministic COSE_Key map:
`{1: 2, 3: -7, -1: 1, -2: x-bstr32, -3: y-bstr32}`. `signerKeyId` is SHA-256 of those exact 77
bytes. `signerAuthorizationIdentity` is:

```text
SHA-256(
  UTF-8("capsule.installation-root-key-authorization/v0") || 0x00 ||
  signerKeyId[32] || installationId[16] || trustEpochSequence-u64be[8] ||
  trustEpochCandidateDigest[32]
)
```

That identity is a locally resolved binding, not a key-discovery mechanism. The protected `kid`
must equal the payload key ID and the locally expected authorization identity must match. No key,
URL, DID resolver, certificate, or dynamic trust material is accepted from an envelope.

The request binds the installation identity; containing-release, I1-profile-candidate, component-
profile, and installation-manifest-candidate digests; epoch sequence `1` and candidate digest;
Supervisor identity and nonroot UID; every fixed root/retained-entry/owner/store profile value;
zero attempt and guest authority; a nonzero 32-byte nonce; `issuedAt <= notBefore < expiresAt` with
an inclusive maximum 300-second issued-to-expiry window; and
`one-shot-exact-envelope-replay-only`.

The record repeats those stable installation/release/I1/component/epoch/Supervisor/key bindings,
binds SHA-256 and length of the exact request payload/envelope plus the nonce, binds the observed
root/owner/store facts and genesis digest, carries its own `issuedAt <= notBefore < expiresAt`
window bounded to 300 seconds and to the retained request expiry, and fixes:

- installation transition `installation-candidate-to-protected-root-disabled`;
- state transition `protected-root-validated-disabled`;
- one-use disposition `commit-once`;
- replay disposition `exact-retained-envelope-response-replay-only`; and
- attempts, runtime, backend, and guest values all `false`.

The I2A phrase “record payload digest” is resolved without a recursive field: SHA-256 of the exact
record payload is derived after encoding and retained beside the exact envelope/file/anchor. It is
not embedded inside the payload whose digest it would define. The signed payload instead fixes the
exact record-retention policy and all authority-bearing bindings.

## Deterministic envelope and exact limits

Both objects use tagged COSE_Sign1 (`tag 18`), an embedded payload, protected `alg=-7`, exact
content type, 32-byte `kid`, an empty unprotected map, empty external AAD, and a 64-byte P1363
`R || S` signature. Detached payloads, DER signatures, any other tag/header, and trailing data
refuse.

Raw caps are checked on caller-owned bytes before a defensive copy. Nested bstr lengths are framed
and capped before copying. The calculated maxima come from independently encoding the complete
closed maximum shapes; no string or collection is padded to fit a chosen cap.

| Object | payload raw / calculated | protected raw / calculated | envelope raw / calculated | depth/items/map/array payload profile |
| --- | ---: | ---: | ---: | --- |
| request | 2,048 / **861** | 256 / **98** | 4,096 / **1,033** | 2 / 69 / 34 / 0 |
| record | 4,096 / **1,527** | 256 / **97** | 6,144 / **1,698** | 2 / 139 / 69 / 0 |

Every maximum is inclusive. The exact calculated-maximum objects accept. Calculated maximum plus
one refuses at the calculated gate. Each raw cap plus one refuses before copy, state, signature,
or authority work. The predecoder additionally rejects nonpreferred integers/lengths, indefinite
values, duplicate or non-deterministically ordered keys, invalid UTF-8, unsafe integers, floats,
unsupported simple values/tags, unknown fields, wrong depth/item/map/array shapes, and trailing
bytes. The typed fxamacker v2.9.2 decode runs only after these checks on the one private copy,
re-encodes the closed private wire struct, and requires byte equality. No generic CBOR/COSE API is
exported. `go-cose` remains absent and production-NO_GO.

## Replay, time, and refusal semantics

- Fresh request: admit once only while `notBefore <= trustedNow < expiresAt`.
- Pending exact envelope digest plus exact nonce: resume the same transaction.
- Pending different envelope/payload/key/nonce or a reused nonce: refuse replay/substitution.
- Completed request, including an exact replay: fixed already-enrolled refusal before state.
- Fresh/pending record: validate its own active time window and commit once against the exact
  retained request bindings.
- Completed exact record after response loss: return the byte-identical retained envelope.
- Completed different record/nonce/envelope: refuse; never sign, create, or commit again.

Signature bytes are retained for byte-identical response recovery but never define semantic or
replay identity. Payload bytes, payload/envelope digests, the authorization identity, installation,
epoch, Supervisor, and nonce own the decision.

## Future installed transaction ordering — not implemented

The future I2B3/I2B4 implementation must preserve this exact order:

1. Authenticate the exact Coordinator peer and cap the message before body copy; validate the
   signed request, locally authorized key identity, purpose/audience/install/release/I1/epoch/
   Supervisor bindings, time, nonce, and fresh-or-exact-pending disposition.
2. From Supervisor-obtained container and retained directory descriptors only, commit the exact
   pending request journal. Create and fully validate only the fixed staging root, owner, retained
   request, and conformance fixed-v1/no-guest store; hold the enrolled owner lock. These staging
   objects remain attempts-disabled and are not yet an enrolled protected root.
3. Return the bounded typed observation on the authenticated live connection. The Coordinator
   constructs/signs the record; the Supervisor validates its signature, authorization identity,
   exact request digests, observations, transition, retention, and zero-admission fields **before**
   publishing or enrolling the protected root.
4. Exclusively publish the validated staging root from Supervisor-controlled descriptors, retain
   the exact record file, commit the create-only epoch anchor once, and byte-readback the file and
   anchor. No existing nonidentical object may be replaced or adopted.
5. Reopen the protected root/record/owner/store without creation, reacquire and recheck the owner,
   perform fixed-v1/no-guest recovery, and return only `protected-root-validated-disabled` when
   every check is clean. Otherwise retain evidence and enter `repair-required`.

This ordering refines I2A by keeping the descriptor-created world staged until the signed record is
validated; “create protected root” means the exclusive publication/enrollment edge, not permission
for the Coordinator or caller to create filesystem objects.

## Failure matrix and zero-effect oracle

| Failure | Required future result |
| --- | --- |
| wrong domain/purpose/audience/install/epoch/Supervisor/release/I1/root/lock/store/key | refuse before creation/publication; zero enrollment |
| zero nonce, reused nonce, different exact bytes for the replay tuple | replay refusal; retained request/record unchanged |
| future, stale, expired, or over-300-second first admission | time refusal; zero filesystem/Keychain/store authority |
| noncanonical/duplicate/unknown/trailing CBOR, cross-object payload, wrong signature shape | bounded structural refusal before semantic state |
| payload/envelope/request-digest substitution or unknown signer authorization | signature/binding refusal; no publish/anchor/store work |
| death before journal commit | no bootstrap effect |
| death after journal but before signed-record validation | exact pending replay only; never a new request/root/store/key/epoch |
| signed-record response loss after commit | byte-identical retained envelope; no second signature or enrollment |
| record valid but publish/anchor/readback/reopen indeterminate | reopen exact predecessor/new world; otherwise `repair-required` |
| root/lock/store/record missing, corrupt, linked, replaced, relocated, or mismatched | `repair-required`; no recreate, normalize, delete, repair, or attempt activation |

All current fixtures retain zero state mutation, zero Keychain/key operation, zero IPC endpoint,
zero process/service, and zero runtime/backend/guest operation. The Swift verifier uses only a
checked-in public P-256 fixture key and CryptoKit verification; it performs no signing, key
generation, Keychain, Secure Enclave, or LocalAuthentication operation.

## Corpus and limitations

`scripts/generate-i2b-bootstrap-conformance.mjs` independently encodes deterministic CBOR and
RFC6979 P-256 signatures. Go verifies the closed fxamacker-backed passive boundary and defensive
ownership. `scripts/verify-i2b-bootstrap-fixtures.swift` is a real independent strict CBOR,
Sig_structure, CryptoKit ES256, binding, time, and replay verifier. Agreement is 71/71 cases: six
accepts and 65 refusals.

This does not admit a production signer or wrapper, perform user presence, select a production key
implementation, access signing identities, create a Keychain item/container/root/lock/store,
activate an IPC/service/process, select a product store, enable an attempt, or admit a runtime,
backend, Apple profile, code signature, or guest. Installed I2B remains `BLOCKED` on the separately
authorized I2B2-I2B5 work.
