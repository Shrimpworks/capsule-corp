// Package hostrunnercontract holds the small pieces of the C2B host-runner
// console-port contract that hostrunnerpassive and hostrunnermaterialized both
// need to state identically.
//
// The two verifier packages check byte-exact known-answer artifacts for two
// distinct, independently governed pipeline stages (the passive source
// contract, C2B v1, versus the materialized build profile, C2B v4 — see
// docs/protocol/GOVERNED_DENO_CORE_C2B_HOST_RUNNER_SOURCE_V1.md and
// docs/protocol/GOVERNED_DENO_CORE_C2B_MATERIALIZED_PROFILE_V4.md). Each
// package's own validation logic and refusal codes remain independently
// frozen and are not touched by this package. What is shared here is narrow
// and physical, not governance state:
//
//   - the Port wire shape and the three fixed console ports every stage of
//     the pipeline must agree on (capsule.source, capsule.input,
//     capsule.completion), since a real port re-wiring must change every
//     verifier in lockstep and previously had nothing to catch a mismatch
//     between the two independently hand-maintained literals; and
//   - the bounded-read and decoder-trailing-data helpers, whose error
//     handling had drifted between the two packages by accident (one silently
//     dropped Close errors and collapsed decode errors to a generic code)
//     rather than by any deliberate per-stage design choice.
//
// This package has no product consumer, exactly like its two callers.
package hostrunnercontract
