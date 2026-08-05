# CDDL contract status

These CDDL files define candidate canonical-CBOR contracts for Capsule's internal security
objects. They are separate from the public JSON API and from the repository's current pre-freeze
JSON Schema scaffold.

[`approval-grant-v0.cddl`](approval-grant-v0.cddl) is the first Gate A2-derived contract. It fixes
the tested map labels, field types, sizes, COSE protected headers, empty unprotected map, and raw
ES256 signature shape. It is **not yet a production-frozen contract**: ADR-0019 remains Proposed,
and the final approval/session bindings and identifier formats still need Phase 2 review.

[`supervisor-bootstrap-v0.cddl`](supervisor-bootstrap-v0.cddl) freezes the passive I2B1
`SupervisorBootstrapRequestV0` and `SupervisorBootstrapRecordV0` payloads and their exact
Sign1-only envelopes. Its independently generated corpus has real Go/Swift verification, but the
objects are unwired and no production signer, key, Keychain, IPC, protected-root, or store path is
implemented.

[`candidates/`](candidates/) contains the passive Phase 2A `ExecutionPlan` and `PlanRegistration`
payload candidates plus their shared scalar definitions. Their byte-exact fixtures are under
[`../fixtures/`](../fixtures/). The corresponding Go and TypeScript types are decoded views only;
there is no product CBOR codec or registration endpoint in this slice. The minimum execution plan
omits unresolved resource and backend-transport values and cannot authorize execution.

Candidate installation, registration, and Supervisor IDs are distinct semantic domains and reject
the all-zero 16-byte value. CDDL expresses their byte width; object-specific semantic validators
own the nonzero and domain-binding rules. Digest zero-value policy remains an explicit freeze
decision.

CDDL describes the data model, not every acceptance rule. Implementations must additionally:

- reject non-deterministic encodings, duplicate keys, indefinite lengths, invalid UTF-8, unknown
  tags/headers/fields, floats, unsupported integers, detached payloads, and trailing bytes;
- enforce raw byte, collection, nesting, and allocation bounds before ordinary decoding;
- verify the exact protected header, signature, object semantics, freshness, and local key
  authorization;
- retain the received canonical payload bytes as authoritative instead of replacing them with a
  decode-and-re-encode result.

The approval, plan, and registration byte fixtures in [`../fixtures/`](../fixtures/) are verified
by `pnpm verify:schemas`. Full production decoder, envelope, and negative canonical-CBOR vectors
remain future work; the Gate A2 experiment continues to retain the complete approval corpus until
production wrappers and conformance suites replace it.
