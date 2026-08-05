# Passive authority-plane cutover result

Date: 2026-08-05

Status: `PASSED` for the exact passive/unwired fixture and facade scope. Product activation remains
`BLOCKED`.

## Result

The legacy broad JobProposal source shape is replaced in v0 by exactly one `main.mjs` field. The
decoder refuses another entrypoint, another member, or a historical JS/TS/CJS source map. Semantic
resolution retains strict UTF-8 bytes without normalization, rejects a leading BOM and bytes above
262,144, and derives the canonical one-member SourceManifest.

`internal/execution/authorityplane/` defines one closed `RegisterPlanV0Request`, the exact 562-byte
role projection, and one `GetRegisteredPlanV0Reply`. Registration resolves nominal roles before
decoding the plan, binds the plan to the exact SourceManifest and source bytes, and publishes plan,
resolved bindings, registration, manifest, and source in one cloned fixed-store transaction. A
failed commit publishes none. Broker lookup accepts only a Supervisor-issued RegistrationID and
returns defensive copies of those retained bytes; it has no execute-time replacement argument.

The aggregate limits are fixture-generated from closed field maxima:

- RegisterPlanV0: 65,536 + 562 + 95 + 262,144 = 328,337 bytes.
- GetRegisteredPlanV0 reply: 328,337 + 4,096 = 332,433 bytes.

`schemas/conformance/authority-plane-v0/` retains the cross-language known answers. The independent
Node verifier and Go encoder/decoder agree on the exact role projection, plan, SourceManifest, and
source bytes. The field-authority manifest classifies both method projections and the exact
`main.mjs` proposal field in the same change.

## Boundaries

This result adds no service name, XPC listener, authenticated endpoint, key use, approval signing,
runtime, backend, VM, guest, or execute path. The fixed store is a passive in-process transaction
oracle, not the production durable Supervisor store. Historical registrationstate objects and old
fixtures remain repository history and passive test material; the new facade exposes no route that
accepts their source-less shape as equivalent authority.

ADR-0040 removes Source Validator R5B from this contract's predecessor chain. A future Broker
consumer may independently validate and render only the exact fetched Supervisor-owned bytes; it
must not add a validator result to registration or revive broad source authority.

The next passive boundary is retained in
[the S1 logical-envelope result](AUTHENTICATED_LOCAL_IPC_S1_PASSIVE_CONTRACT.md). It versions the
two method records with the existing v0 role record and adds no native transport or authentication
claim.
