# Phase 2 internal contract candidates

These files are passive deterministic-CBOR candidates. They do not activate a daemon, Broker,
Supervisor, SDK, or MCP wire protocol, and they do not change ADR-0019 from Proposed.

[`common-v0.cddl`](common-v0.cddl) defines shared scalar widths. Candidate object files are reviewed
together with those definitions. Product decoders must eventually enforce deterministic encoding,
raw allocation bounds, exact object wrappers, domain-specific identifiers/digests, and semantic
bindings in addition to the CDDL shape.

`PlanRegistration` is currently an authenticated local Supervisor response, not a portable signed
authority object. Its candidate payload contains no plan bytes and cannot be submitted by a caller
to mint registration authority. Exact authenticated-IPC and optional exported-evidence treatment
remain later decisions.

[`typescript-approved-byte-v0.cddl`](typescript-approved-byte-v0.cddl) defines the separate passive
ADR-0030 original/executable/profile/options/record object family.
[`execution-plan-v1-typescript-approved-byte.cddl`](execution-plan-v1-typescript-approved-byte.cddl)
is an illustrative full-shape atomic migration candidate only. Neither changes the current plan v0
wrapper or permits dual plan-version acceptance.
