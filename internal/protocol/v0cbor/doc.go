// Package v0cbor provides the unwired fxamacker-backed codec for the exact
// implementation-eligible Capsule v0 CBOR object set.
//
// The set currently contains only SourceManifest v0.  ExecutionPlan,
// PlanRegistration, ApprovalGrant, and the conditional TypeScript family are
// intentionally absent until their separate schema and signed-object gates
// freeze.  This package contains no COSE surface and no generic CBOR API.
package v0cbor
