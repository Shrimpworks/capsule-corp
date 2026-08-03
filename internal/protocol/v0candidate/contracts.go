// Package v0candidate contains strict, unwired wrappers and decoded views of
// passive Phase 2 contract candidates. It does not persist registrations,
// activate a Supervisor implementation, or authorize execution. Exact
// received deterministic-CBOR bytes remain authoritative.
package v0candidate

import "fmt"

// MaxSafeInteger is the largest integer both this decoder and the
// TypeScript side of the protocol can represent exactly (2^53-1),
// matching TypeScript's Number.isSafeInteger ceiling.
const MaxSafeInteger uint64 = 9_007_199_254_740_991

// ObjectVersion is the closed candidate object-version discriminator.
type ObjectVersion uint64

// CandidateObjectVersion is the only object version this decoder accepts; any other value is unsupported protocol.
const CandidateObjectVersion ObjectVersion = 0

const (
	ExecutionPlanObjectType    = "capsule.execution-plan"
	PlanRegistrationObjectType = "capsule.plan-registration"
)

// InstallationID, RegistrationID, and SupervisorID are distinct nominal
// 16-byte identity domains. Each is validated nonzero by its New*
// constructor and is never accepted in a field belonging to another role,
// even though all three share the same underlying width (ADR-0023).
type InstallationID [16]byte
type RegistrationID [16]byte
type SupervisorID [16]byte

// The digest types below are distinct nominal 32-byte SHA-256 domains, one
// per role. Unlike the ID types above, every 32-byte bit pattern is
// structurally accepted, including all-zero: existence and byte equality
// are binding checks owned by the role-specific resolver, not this decoder
// (see newDigest).
type ExecutionPlanDigest [32]byte
type TrustEpochDigest [32]byte
type SourceManifestDigest [32]byte
type InlineInputDigest [32]byte
type RuntimeBundleManifestDigest [32]byte
type ProfileReviewAttestationDigest [32]byte
type ProfileRegistryEntryDigest [32]byte
type BackendValidationRecordDigest [32]byte
type BackendConfigurationDigest [32]byte
type TrustSnapshotDigest [32]byte
type PolicyDecisionDigest [32]byte

// UInt53 is an unsigned integer bounded to the safe-integer range (0 through MaxSafeInteger).
type UInt53 uint64

// PositiveUInt53 is a UInt53 additionally constrained to be nonzero.
type PositiveUInt53 uint64

// InputSlot is the closed set of logical input slot identifiers (ADR-0017); v0 defines exactly one.
type InputSlot string

const PrimaryDataInputSlot InputSlot = "primary-data"

// OutputSlot is the closed set of logical output slot identifiers (ADR-0017); v0 defines exactly one.
type OutputSlot string

const TransformedJSONOutputSlot OutputSlot = "transformed-json"

// WallTimeOrigin records whether a resolved wall-time limit came from the caller's request or the trusted default.
type WallTimeOrigin string

const (
	WallTimeRequested      WallTimeOrigin = "requested"
	WallTimeTrustedDefault WallTimeOrigin = "trusted-default"
)

// ExecutionPlan is the decoded view of the minimum passive candidate. The
// candidate intentionally omits unresolved resource and transport fields and
// therefore cannot authorize execution.
type ExecutionPlan struct {
	ObjectType                      string
	ObjectVersion                   ObjectVersion
	InstallationID                  InstallationID
	EpochSequence                   UInt53
	EpochDigest                     TrustEpochDigest
	SourceManifestDigest            SourceManifestDigest
	SourceEntrypoint                string
	SourceByteLength                UInt53
	InputSlot                       InputSlot
	InlineInputDigest               InlineInputDigest
	InlineInputByteLength           UInt53
	RuntimeProfileAlias             string
	RuntimeBundleManifestDigest     RuntimeBundleManifestDigest
	ProfileReviewAttestationDigests []ProfileReviewAttestationDigest
	ProfileRegistryEntryDigest      ProfileRegistryEntryDigest
	BackendValidationRecordDigest   BackendValidationRecordDigest
	BackendConfigurationDigest      BackendConfigurationDigest
	TrustSnapshotDigest             TrustSnapshotDigest
	PolicyDecisionDigest            PolicyDecisionDigest
	WallTimeMS                      PositiveUInt53
	WallTimeOrigin                  WallTimeOrigin
	OutputSlot                      OutputSlot
	// OutputMaxJSONBytes is decoded as a structural positive integer only.
	// Unlike WallTimeMS, it is not yet checked against a trusted policy
	// ceiling anywhere in this package or the resolver — a deliberately
	// deferred gap per ADR-0023 and PHASE_2A_CONTRACT_FOUNDATION.md, not
	// an oversight. See packages/protocol/src/job-proposal.ts's
	// InlineJsonOutputProposal for the matching TS-side note.
	OutputMaxJSONBytes PositiveUInt53
	ExpiresAt          UInt53
}

// PlanRegistration is a decoded view of the Supervisor-issued registration
// payload returned over authenticated typed local IPC. It contains no plan
// bytes and is not independently signed portable authority.
type PlanRegistration struct {
	ObjectType           string
	ObjectVersion        ObjectVersion
	RegistrationID       RegistrationID
	RegistrationSequence PositiveUInt53
	PlanDigest           ExecutionPlanDigest
	InstallationID       InstallationID
	EpochSequence        UInt53
	EpochDigest          TrustEpochDigest
	SupervisorID         SupervisorID
	ExpiresAt            UInt53
}

// The New* constructors below all validate and copy raw bytes into their
// role-specific nominal type via newIdentifier (16-byte, nonzero) or
// newDigest (32-byte, all-zero-accepting); see those functions for the
// exact validation each performs.

func NewInstallationID(value []byte) (InstallationID, error) {
	return newIdentifier[InstallationID](value, "installation")
}

func NewRegistrationID(value []byte) (RegistrationID, error) {
	return newIdentifier[RegistrationID](value, "registration")
}

func NewSupervisorID(value []byte) (SupervisorID, error) {
	return newIdentifier[SupervisorID](value, "Supervisor")
}

func NewExecutionPlanDigest(value []byte) (ExecutionPlanDigest, error) {
	return newDigest[ExecutionPlanDigest](value, "execution-plan")
}

func NewTrustEpochDigest(value []byte) (TrustEpochDigest, error) {
	return newDigest[TrustEpochDigest](value, "trust-epoch")
}

func NewSourceManifestDigest(value []byte) (SourceManifestDigest, error) {
	return newDigest[SourceManifestDigest](value, "source-manifest")
}

func NewInlineInputDigest(value []byte) (InlineInputDigest, error) {
	return newDigest[InlineInputDigest](value, "inline-input")
}

func NewRuntimeBundleManifestDigest(value []byte) (RuntimeBundleManifestDigest, error) {
	return newDigest[RuntimeBundleManifestDigest](value, "runtime-bundle-manifest")
}

func NewProfileReviewAttestationDigest(value []byte) (ProfileReviewAttestationDigest, error) {
	return newDigest[ProfileReviewAttestationDigest](value, "profile-review-attestation")
}

func NewProfileRegistryEntryDigest(value []byte) (ProfileRegistryEntryDigest, error) {
	return newDigest[ProfileRegistryEntryDigest](value, "profile-registry-entry")
}

func NewBackendValidationRecordDigest(value []byte) (BackendValidationRecordDigest, error) {
	return newDigest[BackendValidationRecordDigest](value, "backend-validation-record")
}

func NewBackendConfigurationDigest(value []byte) (BackendConfigurationDigest, error) {
	return newDigest[BackendConfigurationDigest](value, "backend-configuration")
}

func NewTrustSnapshotDigest(value []byte) (TrustSnapshotDigest, error) {
	return newDigest[TrustSnapshotDigest](value, "trust-snapshot")
}

func NewPolicyDecisionDigest(value []byte) (PolicyDecisionDigest, error) {
	return newDigest[PolicyDecisionDigest](value, "policy-decision")
}

// NewUInt53 validates value is within the safe-integer range.
func NewUInt53(value uint64) (UInt53, error) {
	if value > MaxSafeInteger {
		return 0, schema(fmt.Sprintf("value must not exceed %d", MaxSafeInteger))
	}
	return UInt53(value), nil
}

// NewPositiveUInt53 validates value is nonzero and within the safe-integer range.
func NewPositiveUInt53(value uint64) (PositiveUInt53, error) {
	if value == 0 || value > MaxSafeInteger {
		return 0, schema(fmt.Sprintf("value must be between 1 and %d", MaxSafeInteger))
	}
	return PositiveUInt53(value), nil
}

func newIdentifier[T ~[16]byte](value []byte, role string) (T, error) {
	var result T
	if len(value) != len(result) {
		return result, schema(fmt.Sprintf("%s ID must contain exactly 16 bytes", role))
	}
	copy(result[:], value)
	if result == (T{}) {
		return T{}, schema(fmt.Sprintf("%s ID must be nonzero", role))
	}
	return result, nil
}

// newDigest accepts every 32-byte bit pattern structurally, including all
// zeroes, unlike newIdentifier. ADR-0023 rejects special-casing the all-zero
// SHA-256 value here: width and nominal role are structural, but existence
// and byte equality are binding checks owned by the role-specific resolver.
// An unresolved all-zero digest therefore fails binding like any other
// unknown digest, rather than being rejected at decode time.
func newDigest[T ~[32]byte](value []byte, role string) (T, error) {
	var result T
	if len(value) != len(result) {
		return result, schema(fmt.Sprintf("%s SHA-256 digest must contain exactly 32 bytes", role))
	}
	copy(result[:], value)
	return result, nil
}
