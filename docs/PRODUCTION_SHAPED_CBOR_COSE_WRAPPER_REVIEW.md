# Production-shaped CBOR/COSE wrapper review

Date: 2026-08-05

```text
Work item: signed-object freeze and independent review of the narrow Go CBOR/COSE wrapper
Status: PASSED
Scope: passive, local-only review and hardening of the I2B1 request/record verifier using checked-in
  public-key fixtures; no signer, live key, ceremony, service, IPC activation, runtime, backend,
  VM, guest, release, or product admission
Evidence or reason: the reviewed Go wrapper and independent CryptoKit-backed Swift verifier agree
  on 95 generated cases: 9 accepts and 86 refusals. The review corrected request-pair binding,
  complete repeated-field binding, and payload-owned replay identity; added exact-byte/digest
  accessors and bounded fuzz targets; and retained direct restoration cases for both object types.
Remaining work: a production caller that establishes authenticated role/profile/CDHash facts before
  body copy, a local installation-key authorization resolver, a durable transaction/replay ledger,
  installed same-byte consumers, live-key signing evidence, and the remaining signed-object family.
Next action: keep this package passive and unwired. Admit a production consumer only with the
  separately authorized installed Coordinator/Supervisor ceremony and its retained fault corpus.
Parent status: ADR-0019 and the complete production signing workstream remain BLOCKED.
```

## Applicable signed-object freeze

This review freezes the applicable signed set for the current production-shaped Go wrapper, not
the future object inventory and not ADR-0019 as a whole.

| Object | Decision | Reason |
| --- | --- | --- |
| `SupervisorBootstrapRequest` v0 | **FROZEN FOR PASSIVE VERIFICATION** | Closed 34-field CDDL, exact field authority, request media/purpose/audience, caps, known answers, and Go/Swift verification exist. |
| `SupervisorBootstrapRecord` v0 | **FROZEN FOR PASSIVE VERIFICATION** | Closed 69-field CDDL, exact request/observation/transition bindings, caps, known answers, and Go/Swift verification exist. |
| `ApprovalGrant` v0 | **SIGNED PRE-FREEZE / EXCLUDED** | Production Approval-key authorization, active role/profile/CDHash path, Swift wrapper, and Broker/Supervisor same-byte path remain open. Its retained verifier stays fixture-only. |
| `ExecutionPlan` v0 and `PlanRegistration` v0 | **PRE-FREEZE / EXCLUDED** | Their authority model and consumer contract remain pre-freeze; registration is not independently signed portable authority. |
| transcript, receipt, trust, update, runtime, TUF, and witness candidates | **NOT IN THIS SET** | No complete production-frozen signed profile and consumer pair exists. |

`SourceManifest` v0 remains eligible for the separate canonical-CBOR wrapper but is not a signed
COSE object. No current exclusion may enter `bootstrapauthpassive` through a generic object switch.

## Review findings and retained corrections

The initial passive implementation had three authorization-relevant defects:

1. Record verification hashed the separately supplied request envelope and request payload without
   proving that the payload was the envelope's embedded signed payload. A valid record could bind a
   split pair assembled by its caller.
2. Record verification compared only a subset of the request fields projected into the record.
   Component profile and other installation, Supervisor, root, owner, store, key, epoch, name, and
   time fields could be omitted from the cross-object check.
3. Replay decisions keyed the complete envelope digest. Mathematically equivalent ES256 signatures
   over one canonical payload could therefore split one semantic transaction, contrary to
   ADR-0019 and ADR-0024.

The retained implementation now caps and frames the supplied request before copying it, proves the
exact embedded-payload pair, verifies that request with the independently supplied trusted public
key, and compares every repeated request field. Replay state keys the canonical payload digest and
nonce. Envelope digest and exact envelope bytes remain integrity/evidence and response-recovery
material; an equivalent signature never replaces the already retained envelope. ADR-0038 and the
signed disposition string now state this payload-owned identity explicitly.

The verifier also retains exact protected, payload, and envelope bytes plus separate payload and
envelope digests behind defensive-copy accessors. No decode-and-re-encode output replaces received
identity.

## Control audit

| Control | Reviewed result |
| --- | --- |
| Cap before copy | Raw envelope caps and nested protected/payload bstr caps run in `frameEnvelope` before `bytes.Clone`. Record-bound request envelope/payload are independently capped before their defensive copies. Calculated closed-shape maxima run before copy as well. |
| Deterministic canonical-on-wire bytes | Capsule predecode rejects nonpreferred and unordered forms; private fxamacker typed decode is followed by deterministic re-encoding and byte equality. Received bytes remain authoritative. |
| Closed fields and domains | Duplicate, unknown, out-of-order, wrong-type, wrong-object, wrong-version, wrong-purpose, wrong-audience, cross-object, noncanonical, trailing, and cap-plus-one cases refuse with bounded Capsule classifications. |
| Protected headers | The only accepted protected map is canonical labels `1`, `3`, `4` with ES256 `-7`, the exact object media type, and a 32-byte `kid`. Wrong algorithm/content type, unknown labels, and noncanonical order refuse for each object. |
| COSE framing | Only tagged embedded-payload COSE_Sign1 is accepted. The unprotected map must be exactly empty; detached payloads and trailing bytes refuse. |
| Sig_structure | Capsule constructs `["Signature1", protected, h'', payload]`; valid fixtures signed with nonempty external AAD refuse. |
| ES256 signature | Exactly 64-byte IEEE P1363 `R || S` is accepted. DER and other lengths refuse. Valid ordinary and complementary-S forms verify, but signature bytes never own replay identity. |
| Key and authorization binding | Protected `kid`, payload `signerKeyId`, exact 77-byte P-256 COSE_Key, installation ID, epoch, and domain-separated authorization identity must match the caller-supplied expected object. Signature verification uses that expected public key, not an envelope-discovered key. |
| Role, purpose, audience, install, epoch, profile, digest, nonce, and time | The closed expected request/record values and all record-to-request projections are compared before acceptance. Nonzero nonce, fixed windows, active trusted time, and replay disposition are enforced. |
| Payload replay identity | SHA-256 of exact canonical payload plus nonce owns request/record replay decisions. Complementary-S envelopes resolve to the same payload identity while preserving distinct envelope evidence. |
| Restoration and mutation | The 95-case shared corpus covers parser/profile restoration, exact caps, cross-object and binding substitutions, signature/AAD forms, time, replay, and trusted bindings. Three Go fuzz targets exercise request bytes, record bytes, and the separately supplied request pair. |

The role/profile statement is deliberately conditional. The wrapper validates the exact role-bound
values passed by its caller; it does not authenticate a process or establish a code-signing role.
ADR-0038 requires the future IPC boundary to authenticate the Coordinator, validate exact
profile/CDHash/install facts, and cap the message before body copy. No such service or active
profile exists here. Likewise, the authorization-identity calculation is verified, but the future
local registry that establishes whether that key is authorized for this role, purpose,
installation, and epoch is not implemented.

## Library and Capsule responsibilities

`github.com/fxamacker/cbor/v2` v2.9.2 remains conditionally adopted only for deterministic encoding
of private fixed wire structs and typed field decoding after Capsule predecode. It does not own raw
caps, canonical-on-wire acceptance, duplicate/unknown refusal policy, object domains, COSE framing,
protected-header policy, `Sig_structure`, key authorization, trusted-context comparison, time,
replay, cross-object binding, or byte ownership.

Go's standard cryptographic library supplies SHA-256 and P-256 verification. Capsule constructs
the exact `Sig_structure`, parses raw `R || S`, supplies the trusted key, and owns all authorization
decisions. `veraison/go-cose` remains absent from the product graph and `NO_GO` for product use; no
new evidence justified changing that decision.

## Cross-language comparison and limitations

The Node generator independently emits deterministic CBOR and RFC6979 P-256 fixture signatures.
The Go verifier and an independent strict Swift parser/CryptoKit verifier agree on all 95 cases.
The Swift path compares the retained request envelope with its embedded payload, verifies both
request and record signatures, checks the complete repeated-field projection, and applies the same
payload-owned replay model. It uses only the checked-in public fixture key during verification; no
live private key, Keychain, Secure Enclave, LocalAuthentication, or signing ceremony participates.

This evidence supports only the passive repository wrapper and its exact fixtures. It is not a
production wrapper admission, a continuous-service claim, an installed role/profile result, or
evidence for Approval, execution, runtime, backend, VM, or guest authority.

## Verification

The final tree passed the repository-required Node 22.22.1/pnpm 10 and Go gates:

- `pnpm install`, `pnpm check`, `pnpm lint`, credential-free `pnpm test`,
  `pnpm verify:schemas`, and `pnpm verify:adrs`;
- `go test ./...`, `go vet ./...`, `go build ./...`, and `golangci-lint run ./...`;
- `govulncheck@latest ./...`: no called-symbol or imported-package vulnerability;
- focused bootstrap race and coverage runs: pass, 80.9% statement coverage;
- independent Swift: 9 accepts and 86 refusals; and
- three 10-second Go fuzz runs: 149,099 request, 59,747 record, and 594,298 request-pair
  executions, with no failed invariant.

The full Node test was run with commit signing disabled for its disposable temporary Git
repositories. The ordinary user GPG agent and signing key were not accessed.
