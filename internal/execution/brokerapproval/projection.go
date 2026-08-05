// Package brokerapproval freezes an unwired read-only Approval Broker
// projection. It has no UI, key, signer, IPC endpoint, store mutation,
// LocalAuthentication operation, runtime, backend, or guest effect.
package brokerapproval

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"capsule.local/capsule/internal/execution/authorityplane"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	SourceDisplayEncoding = "capsule.bytewise-ascii-escape/v0"
	SourceDisplayMaxBytes = 4 * v0candidate.MJSMainSourceMaxBytes

	SourceContentPolicy      = "exact strict-utf8 main.mjs; no BOM, normalization, rewrite, transform, or second member"
	InlineInputContentPolicy = "current Supervisor projection contains canonical JSON digest and length only; content bytes are not shown"
	InteractionEvidenceState = "contract-only-no-installed-ui-evidence"
)

var fixedWarnings = [...]string{
	"Generated code may encode approved input through allowed output, metadata, state, or timing.",
	"User presence attributes one key operation; it does not prove source comprehension or correct UI logic.",
	"Focus, UI activation, and synthetic input are not approval evidence.",
	"Current Supervisor readback does not contain inline JSON content bytes; this projection shows only their exact digest and length.",
	"Host source validation is not an internal-alpha admission claim; runtime syntax and loader refusal remain separate blocked evidence.",
	"Runtime/profile identifiers and review digests do not admit a runtime, backend, or product path.",
}

type Classification string

const (
	Malformed    Classification = "MALFORMED"
	Unsupported  Classification = "UNSUPPORTED"
	Schema       Classification = "SCHEMA"
	Domain       Classification = "DOMAIN"
	Binding      Classification = "BINDING"
	Stale        Classification = "STALE"
	LocalFailure Classification = "LOCAL_FAILURE"
)

type projectionError struct {
	classification Classification
	code           string
}

func (err *projectionError) Error() string {
	return fmt.Sprintf("%s: %s", err.classification, err.code)
}

func ErrorClassification(err error) (Classification, bool) {
	var target *projectionError
	if errors.As(err, &target) {
		return target.classification, true
	}
	if value, ok := authorityplane.ErrorClassification(err); ok {
		return Classification(value), true
	}
	if value, ok := v0candidate.ErrorClassification(err); ok {
		return Classification(value), true
	}
	return "", false
}

func refused(classification Classification, code string) error {
	return &projectionError{classification: classification, code: code}
}

// RegisteredPlanReadback is satisfied by authorityplane.GetRegisteredPlanV0Reply.
// Its accessors return defensive Supervisor-retained bytes.
type RegisteredPlanReadback interface {
	ExactPlanBytes() []byte
	ResolvedRoleBindingBytes() []byte
	PlanRegistrationBytes() []byte
	SourceManifestBytes() []byte
	SourceBytes() []byte
}

var _ RegisteredPlanReadback = authorityplane.GetRegisteredPlanV0Reply{}

type TrustedSupervisorContext struct {
	InstallationID v0candidate.InstallationID
	EpochSequence  v0candidate.UInt53
	EpochDigest    v0candidate.TrustEpochDigest
	SupervisorID   v0candidate.SupervisorID
}

type SourceProjection struct {
	Profile             string
	MemberMediaType     string
	ManifestMediaType   string
	Entrypoint          string
	FileCount           uint64
	ByteLength          uint64
	ContentDigestHex    string
	ManifestDigestHex   string
	ContentPolicy       string
	DisplayEncoding     string
	EscapedExactContent string
}

type InlineJSONProjection struct {
	Slot              string
	DigestHex         string
	ByteLength        uint64
	ContentPolicy     string
	ContentBytesShown bool
}

type RuntimeProfileProjection struct {
	Alias                         string
	RuntimeBundleDigestHex        string
	ProfileReviewDigestHexes      []string
	ProfileRegistryDigestHex      string
	BackendValidationDigestHex    string
	BackendConfigurationDigestHex string
	TrustSnapshotDigestHex        string
	PolicyDecisionDigestHex       string
}

type LimitsProjection struct {
	SourceBytes        uint64
	InlineInputBytes   uint64
	WallTimeMS         uint64
	WallTimeOrigin     string
	OutputSlot         string
	OutputMaxJSONBytes uint64
}

type ExpiryProjection struct {
	PlanExpiresAt         uint64
	RegistrationExpiresAt uint64
	EffectiveNow          uint64
}

type InteractionContract struct {
	FocusIsApprovalEvidence           bool
	SyntheticInputIsApprovalEvidence  bool
	KeyOperationRequired              bool
	FreshContextPerOperation          bool
	ContextReusePermitted             bool
	CancelInvalidatesContext          bool
	FocusLossInvalidatesContext       bool
	SessionReplacementRequiresRefetch bool
	InstalledUIEvidence               bool
	EvidenceState                     string
}

type Projection struct {
	RegistrationIDHex         string
	RegistrationSequence      uint64
	PlanDigestHex             string
	InstallationIDHex         string
	EpochSequence             uint64
	EpochDigestHex            string
	SupervisorIDHex           string
	Source                    SourceProjection
	InlineJSON                InlineJSONProjection
	RuntimeProfile            RuntimeProfileProjection
	Limits                    LimitsProjection
	Expiry                    ExpiryProjection
	Warnings                  []string
	Interaction               InteractionContract
	ApprovalEligible          bool
	ApprovalIneligibilityCode string
}

// BuildProjection treats locator only as an untrusted lookup value. It emits a
// view only after the retained registration, exact plan, complete role
// bindings, manifest, source, and trusted Supervisor context agree.
func BuildProjection(locator v0candidate.RegistrationID, readback RegisteredPlanReadback, trusted TrustedSupervisorContext, effectiveNow v0candidate.UInt53) (Projection, error) {
	if locator == (v0candidate.RegistrationID{}) || readback == nil {
		return Projection{}, refused(Malformed, "registration-locator-or-readback")
	}
	if uint64(trusted.EpochSequence) > v0candidate.MaxSafeInteger || uint64(effectiveNow) > v0candidate.MaxSafeInteger {
		return Projection{}, refused(Malformed, "trusted-supervisor-context-range")
	}
	planBytes := readback.ExactPlanBytes()
	bindingBytes := readback.ResolvedRoleBindingBytes()
	registrationBytes := readback.PlanRegistrationBytes()
	manifestBytes := readback.SourceManifestBytes()
	sourceBytes := readback.SourceBytes()
	if len(planBytes)+len(bindingBytes)+len(registrationBytes)+len(manifestBytes)+len(sourceBytes) > authorityplane.GetRegisteredPlanV0MaxBytes {
		return Projection{}, refused(Malformed, "broker-readback-aggregate-cap")
	}
	bindings, err := authorityplane.DecodeRoleBindingsV0(bindingBytes)
	if err != nil {
		return Projection{}, err
	}
	plan, err := v0candidate.DecodeExecutionPlan(planBytes, bindings)
	if err != nil {
		return Projection{}, err
	}
	planView := plan.View()
	if planView.InstallationID != trusted.InstallationID || planView.EpochSequence != trusted.EpochSequence ||
		planView.EpochDigest != trusted.EpochDigest || bindings.InstallationID != trusted.InstallationID || bindings.EpochDigest != trusted.EpochDigest {
		return Projection{}, refused(Binding, "trusted-supervisor-plan-binding")
	}
	registration, err := v0candidate.DecodePlanRegistration(registrationBytes, v0candidate.PlanRegistrationRoleBindings{
		RegistrationID: locator, PlanDigest: plan.Digest(), InstallationID: trusted.InstallationID,
		EpochDigest: trusted.EpochDigest, SupervisorID: trusted.SupervisorID,
	})
	if err != nil {
		return Projection{}, err
	}
	registrationView := registration.View()
	if registrationView.EpochSequence != trusted.EpochSequence || registrationView.ExpiresAt != planView.ExpiresAt {
		return Projection{}, refused(Binding, "registration-plan-state-binding")
	}
	if effectiveNow >= planView.ExpiresAt {
		return Projection{}, refused(Stale, "registration-expired")
	}
	manifest, err := v0candidate.DecodeSourceManifest(manifestBytes, v0candidate.SourceManifestMediaType, sourceBytes)
	if err != nil {
		return Projection{}, err
	}
	manifestView := manifest.View()
	if manifest.Digest() != planView.SourceManifestDigest || manifestView.Entrypoint != v0candidate.MJSMainPath ||
		manifestView.AggregateByteLength != planView.SourceByteLength || len(manifestView.Members) != 1 {
		return Projection{}, refused(Binding, "source-plan-binding")
	}
	escaped, err := EscapeSourceForDisplay(sourceBytes)
	if err != nil {
		return Projection{}, err
	}
	reviews := make([]string, len(planView.ProfileReviewAttestationDigests))
	for index, digest := range planView.ProfileReviewAttestationDigests {
		reviews[index] = hex.EncodeToString(digest[:])
	}
	member := manifestView.Members[0]
	planDigest := plan.Digest()
	manifestDigest := manifest.Digest()
	projection := Projection{
		RegistrationIDHex:    hex.EncodeToString(registrationView.RegistrationID[:]),
		RegistrationSequence: uint64(registrationView.RegistrationSequence), PlanDigestHex: hex.EncodeToString(planDigest[:]),
		InstallationIDHex: hex.EncodeToString(trusted.InstallationID[:]), EpochSequence: uint64(trusted.EpochSequence),
		EpochDigestHex: hex.EncodeToString(trusted.EpochDigest[:]), SupervisorIDHex: hex.EncodeToString(trusted.SupervisorID[:]),
		Source: SourceProjection{
			Profile: v0candidate.MJSSourceProfile, MemberMediaType: v0candidate.MJSSourceMediaType,
			ManifestMediaType: v0candidate.SourceManifestMediaType, Entrypoint: string(manifestView.Entrypoint),
			FileCount: 1, ByteLength: uint64(manifestView.AggregateByteLength),
			ContentDigestHex: hex.EncodeToString(member.ContentDigest[:]), ManifestDigestHex: hex.EncodeToString(manifestDigest[:]),
			ContentPolicy: SourceContentPolicy, DisplayEncoding: SourceDisplayEncoding, EscapedExactContent: escaped,
		},
		InlineJSON: InlineJSONProjection{
			Slot: string(planView.InputSlot), DigestHex: hex.EncodeToString(planView.InlineInputDigest[:]),
			ByteLength: uint64(planView.InlineInputByteLength), ContentPolicy: InlineInputContentPolicy, ContentBytesShown: false,
		},
		RuntimeProfile: RuntimeProfileProjection{
			Alias: planView.RuntimeProfileAlias, RuntimeBundleDigestHex: hex.EncodeToString(planView.RuntimeBundleManifestDigest[:]),
			ProfileReviewDigestHexes: reviews, ProfileRegistryDigestHex: hex.EncodeToString(planView.ProfileRegistryEntryDigest[:]),
			BackendValidationDigestHex:    hex.EncodeToString(planView.BackendValidationRecordDigest[:]),
			BackendConfigurationDigestHex: hex.EncodeToString(planView.BackendConfigurationDigest[:]),
			TrustSnapshotDigestHex:        hex.EncodeToString(planView.TrustSnapshotDigest[:]), PolicyDecisionDigestHex: hex.EncodeToString(planView.PolicyDecisionDigest[:]),
		},
		Limits: LimitsProjection{
			SourceBytes: uint64(planView.SourceByteLength), InlineInputBytes: uint64(planView.InlineInputByteLength),
			WallTimeMS: uint64(planView.WallTimeMS), WallTimeOrigin: string(planView.WallTimeOrigin),
			OutputSlot: string(planView.OutputSlot), OutputMaxJSONBytes: uint64(planView.OutputMaxJSONBytes),
		},
		Expiry:   ExpiryProjection{PlanExpiresAt: uint64(planView.ExpiresAt), RegistrationExpiresAt: uint64(registrationView.ExpiresAt), EffectiveNow: uint64(effectiveNow)},
		Warnings: append([]string(nil), fixedWarnings[:]...),
		Interaction: InteractionContract{
			FocusIsApprovalEvidence: false, SyntheticInputIsApprovalEvidence: false, KeyOperationRequired: true,
			FreshContextPerOperation: true, ContextReusePermitted: false, CancelInvalidatesContext: true,
			FocusLossInvalidatesContext: true, SessionReplacementRequiresRefetch: true,
			InstalledUIEvidence: false, EvidenceState: InteractionEvidenceState,
		},
		ApprovalEligible: false, ApprovalIneligibilityCode: "PASSIVE_CONTRACT_NO_INSTALLED_UI_OR_KEY_OPERATION",
	}
	return projection, nil
}

// EscapeSourceForDisplay produces an ASCII-only, byte-exact representation.
// It never emits untrusted control, bidi, normalization-sensitive, or markup
// characters. Every escaped output can be reversed to the exact source bytes.
func EscapeSourceForDisplay(source []byte) (string, error) {
	validated, err := v0candidate.NewMJSMainSource(source)
	if err != nil {
		return "", err
	}
	exact := validated.AuthoritativeBytes()
	var result bytes.Buffer
	result.Grow(len(exact))
	hexDigits := "0123456789abcdef"
	for _, value := range exact {
		switch {
		case value == '\\':
			result.WriteString("\\\\")
		case value == '\n':
			result.WriteString("\\n")
		case value == '\r':
			result.WriteString("\\r")
		case value == '\t':
			result.WriteString("\\t")
		case value >= 0x20 && value <= 0x7e && value != '<' && value != '>' && value != '&' &&
			value != '"' && value != '\'' && value != '`':
			result.WriteByte(value)
		default:
			result.WriteString("\\x")
			result.WriteByte(hexDigits[value>>4])
			result.WriteByte(hexDigits[value&0x0f])
		}
		if result.Len() > SourceDisplayMaxBytes {
			return "", refused(LocalFailure, "escaped-source-cap")
		}
	}
	return result.String(), nil
}
