# CDDL contract status

These CDDL files define candidate canonical-CBOR contracts for Capsule's internal security
objects. They are separate from the public JSON API and from the repository's current pre-freeze
JSON Schema scaffold.

[`approval-grant-v0.cddl`](approval-grant-v0.cddl) is the first Gate A2-derived contract. It fixes
the tested map labels, field types, sizes, COSE protected headers, empty unprotected map, and raw
ES256 signature shape. It is **not yet a production-frozen contract**: ADR-0019 remains Proposed,
and the final approval/session bindings and identifier formats still need Phase 2 review.

CDDL describes the data model, not every acceptance rule. Implementations must additionally:

- reject non-deterministic encodings, duplicate keys, indefinite lengths, invalid UTF-8, unknown
  tags/headers/fields, floats, unsupported integers, detached payloads, and trailing bytes;
- enforce raw byte, collection, nesting, and allocation bounds before ordinary decoding;
- verify the exact protected header, signature, object semantics, freshness, and local key
  authorization;
- retain the received canonical payload bytes as authoritative instead of replacing them with a
  decode-and-re-encode result.

The byte fixture in [`../fixtures/approval-grant-v0.json`](../fixtures/approval-grant-v0.json) is
verified by `pnpm verify:schemas`. Full envelope and negative vectors remain in the isolated Gate A2
experiment until production wrappers and conformance suites replace them.
