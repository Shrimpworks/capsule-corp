package brokerapproval

import (
	"bytes"
	"context"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"capsule.local/capsule/internal/execution/authorityplane"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

type exactRoles struct {
	expected v0candidate.ExecutionPlanRoleBindings
}

func (roles exactRoles) ResolvePlanRoles(_ context.Context, received v0candidate.ExecutionPlanRoleBindings) (v0candidate.ExecutionPlanRoleBindings, error) {
	return roles.expected, nil
}

type fixedRegistrationID struct{ value v0candidate.RegistrationID }

func (source fixedRegistrationID) NewRegistrationID(context.Context) (v0candidate.RegistrationID, error) {
	return source.value, nil
}

type retainedReadback struct {
	plan, bindings, registration, manifest, source []byte
}

func (value retainedReadback) ExactPlanBytes() []byte           { return bytes.Clone(value.plan) }
func (value retainedReadback) ResolvedRoleBindingBytes() []byte { return bytes.Clone(value.bindings) }
func (value retainedReadback) PlanRegistrationBytes() []byte    { return bytes.Clone(value.registration) }
func (value retainedReadback) SourceManifestBytes() []byte      { return bytes.Clone(value.manifest) }
func (value retainedReadback) SourceBytes() []byte              { return bytes.Clone(value.source) }

func TestBuildProjectionUsesOnlyBoundSupervisorReadback(t *testing.T) {
	locator, reply, trusted := registeredReadback(t)
	projection, err := BuildProjection(locator, reply, trusted, 1_785_456_000)
	if err != nil {
		t.Fatal(err)
	}
	if projection.RegistrationIDHex != strings.Repeat("77", 16) || projection.RegistrationSequence != 1 ||
		projection.PlanDigestHex != "ef268a0b829adc1ce1307203f4b805f63379954ccf41e8e20a7487b6e5acf241" ||
		projection.InstallationIDHex != strings.Repeat("11", 16) || projection.EpochSequence != 7 ||
		projection.EpochDigestHex != strings.Repeat("22", 32) || projection.SupervisorIDHex != strings.Repeat("55", 16) {
		t.Fatal("trusted registration projection drift")
	}
	if projection.Source.Profile != v0candidate.MJSSourceProfile ||
		projection.Source.MemberMediaType != v0candidate.MJSSourceMediaType ||
		projection.Source.ManifestMediaType != v0candidate.SourceManifestMediaType ||
		projection.Source.Entrypoint != "main.mjs" || projection.Source.FileCount != 1 || projection.Source.ByteLength != 50 ||
		projection.Source.ContentDigestHex != "681f39365de1369ee486fa34e88b993c60df5a835006b65e0d8916df717c31cc" ||
		projection.Source.ManifestDigestHex != "c387c80094027ffbcacb573f44f5f6b4dec4d243bb436b24dd644434feaa1d14" ||
		projection.Source.ContentPolicy != SourceContentPolicy || projection.Source.DisplayEncoding != SourceDisplayEncoding ||
		projection.Source.EscapedExactContent != "export default function (value) { return value; }\\n" {
		t.Fatalf("exact source projection drift: %#v", projection.Source)
	}
	if projection.InlineJSON.Slot != "primary-data" || projection.InlineJSON.ByteLength != 118 ||
		projection.InlineJSON.DigestHex != "bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72" ||
		projection.InlineJSON.ContentPolicy != InlineInputContentPolicy || projection.InlineJSON.ContentBytesShown {
		t.Fatalf("inline JSON projection drift: %#v", projection.InlineJSON)
	}
	if projection.Limits.SourceBytes != 50 || projection.Limits.InlineInputBytes != 118 || projection.Limits.WallTimeMS != 5_000 ||
		projection.Limits.WallTimeOrigin != "requested" || projection.Limits.OutputSlot != "transformed-json" || projection.Limits.OutputMaxJSONBytes != 65_536 {
		t.Fatalf("limit projection drift: %#v", projection.Limits)
	}
	if projection.Expiry.PlanExpiresAt != 1_785_456_300 || projection.Expiry.RegistrationExpiresAt != 1_785_456_300 || projection.Expiry.EffectiveNow != 1_785_456_000 {
		t.Fatalf("expiry projection drift: %#v", projection.Expiry)
	}
	if len(projection.RuntimeProfile.ProfileReviewDigestHexes) != 2 ||
		projection.RuntimeProfile.Alias != "fixture-active@1" ||
		projection.RuntimeProfile.RuntimeBundleDigestHex != strings.Repeat("55", 32) ||
		projection.RuntimeProfile.ProfileReviewDigestHexes[0] != strings.Repeat("66", 32) ||
		projection.RuntimeProfile.ProfileReviewDigestHexes[1] != strings.Repeat("67", 32) ||
		projection.RuntimeProfile.ProfileRegistryDigestHex != strings.Repeat("77", 32) ||
		projection.RuntimeProfile.BackendValidationDigestHex != strings.Repeat("88", 32) ||
		projection.RuntimeProfile.BackendConfigurationDigestHex != strings.Repeat("99", 32) ||
		projection.RuntimeProfile.TrustSnapshotDigestHex != strings.Repeat("aa", 32) ||
		projection.RuntimeProfile.PolicyDecisionDigestHex != strings.Repeat("bb", 32) {
		t.Fatalf("runtime/profile projection drift: %#v", projection.RuntimeProfile)
	}
	if len(projection.Warnings) != len(fixedWarnings) || projection.ApprovalEligible ||
		projection.ApprovalIneligibilityCode != "PASSIVE_CONTRACT_NO_INSTALLED_UI_OR_KEY_OPERATION" ||
		projection.Interaction.FocusIsApprovalEvidence || projection.Interaction.SyntheticInputIsApprovalEvidence ||
		!projection.Interaction.KeyOperationRequired || !projection.Interaction.FreshContextPerOperation ||
		projection.Interaction.ContextReusePermitted || !projection.Interaction.CancelInvalidatesContext ||
		!projection.Interaction.FocusLossInvalidatesContext || !projection.Interaction.SessionReplacementRequiresRefetch ||
		projection.Interaction.InstalledUIEvidence || projection.Interaction.EvidenceState != InteractionEvidenceState {
		t.Fatal("warnings, eligibility, or interaction contract drift")
	}
	for index, warning := range fixedWarnings {
		if projection.Warnings[index] != warning {
			t.Fatal("fixed warning text drift")
		}
	}
}

func TestBuildProjectionRefusesStaleMixedOrMutatedReadback(t *testing.T) {
	locator, reply, trusted := registeredReadback(t)
	tests := []struct {
		name    string
		locator v0candidate.RegistrationID
		mutate  func(*retainedReadback, *TrustedSupervisorContext)
		class   Classification
	}{
		{"mixed-registration", repeatedBrokerID[v0candidate.RegistrationID](0x78), func(*retainedReadback, *TrustedSupervisorContext) {}, Domain},
		{"trusted-installation", locator, func(_ *retainedReadback, value *TrustedSupervisorContext) { value.InstallationID[0] ^= 0xff }, Binding},
		{"trusted-epoch", locator, func(_ *retainedReadback, value *TrustedSupervisorContext) { value.EpochDigest[0] ^= 0xff }, Binding},
		{"mutated-source", locator, func(value *retainedReadback, _ *TrustedSupervisorContext) { value.source[len(value.source)-2] ^= 1 }, Domain},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := retainedReadbackFrom(reply)
			candidateTrusted := trusted
			test.mutate(&candidate, &candidateTrusted)
			_, err := BuildProjection(test.locator, candidate, candidateTrusted, 1_785_456_000)
			assertProjectionClassification(t, err, test.class)
		})
	}
	if _, err := BuildProjection(locator, reply, trusted, 1_785_456_300); err == nil {
		t.Fatal("expiry equality accepted")
	} else {
		assertProjectionClassification(t, err, Stale)
	}
	if _, err := BuildProjection(locator, reply, trusted, v0candidate.UInt53(v0candidate.MaxSafeInteger+1)); err == nil {
		t.Fatal("unsafe effective time accepted")
	} else {
		assertProjectionClassification(t, err, Malformed)
	}
}

func TestEscapeSourceForDisplayNeutralizesControlsBidiAndMaximum(t *testing.T) {
	source := []byte("a\\\n\r\t\x00\x1b\u202eb<>&\"'`")
	escaped, err := EscapeSourceForDisplay(source)
	if err != nil {
		t.Fatal(err)
	}
	if escaped != "a\\\\\\n\\r\\t\\x00\\x1b\\xe2\\x80\\xaeb\\x3c\\x3e\\x26\\x22\\x27\\x60" {
		t.Fatalf("escaped source = %q", escaped)
	}
	for _, value := range []byte(escaped) {
		if value < 0x20 || value > 0x7e {
			t.Fatalf("display retained unsafe byte 0x%02x", value)
		}
	}
	maximum := bytes.Repeat([]byte{0x7f}, v0candidate.MJSMainSourceMaxBytes)
	escaped, err = EscapeSourceForDisplay(maximum)
	if err != nil || len(escaped) != SourceDisplayMaxBytes {
		t.Fatalf("exact source maximum: bytes=%d err=%v", len(escaped), err)
	}
	if _, err := EscapeSourceForDisplay(append(maximum, 0x7f)); err == nil {
		t.Fatal("source maximum plus one accepted")
	}
}

func registeredReadback(t *testing.T) (v0candidate.RegistrationID, authorityplane.GetRegisteredPlanV0Reply, TrustedSupervisorContext) {
	t.Helper()
	bindings := brokerBindings(t)
	bindingBytes, err := authorityplane.EncodeRoleBindingsV0(bindings)
	if err != nil {
		t.Fatal(err)
	}
	request, err := authorityplane.NewRegisterPlanV0Request(
		brokerFixture(t, "execution-plan/ordinary.cbor"), bindingBytes,
		brokerFixture(t, "source-manifest/ordinary.cbor"), []byte("export default function (value) { return value; }\n"),
	)
	if err != nil {
		t.Fatal(err)
	}
	locator := repeatedBrokerID[v0candidate.RegistrationID](0x77)
	trusted := TrustedSupervisorContext{
		InstallationID: bindings.InstallationID, EpochSequence: 7, EpochDigest: bindings.EpochDigest,
		SupervisorID: repeatedBrokerID[v0candidate.SupervisorID](0x55),
	}
	facade, err := authorityplane.NewFacade(
		&authorityplane.FixedStore{}, exactRoles{bindings}, fixedRegistrationID{locator},
		authorityplane.SupervisorContext(trusted),
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := facade.RegisterPlanV0(context.Background(), authorityplane.CallContext{Authenticated: true, Role: authorityplane.Daemon, Purpose: authorityplane.RegisterPlanV0Purpose}, request); err != nil {
		t.Fatal(err)
	}
	reply, err := facade.GetRegisteredPlanV0(context.Background(), authorityplane.CallContext{Authenticated: true, Role: authorityplane.Broker, Purpose: authorityplane.GetRegisteredPlanV0Purpose}, authorityplane.GetRegisteredPlanV0Request{RegistrationID: locator})
	if err != nil {
		t.Fatal(err)
	}
	return locator, reply, trusted
}

func retainedReadbackFrom(value authorityplane.GetRegisteredPlanV0Reply) retainedReadback {
	return retainedReadback{
		plan: value.ExactPlanBytes(), bindings: value.ResolvedRoleBindingBytes(), registration: value.PlanRegistrationBytes(),
		manifest: value.SourceManifestBytes(), source: value.SourceBytes(),
	}
}

func brokerBindings(t *testing.T) v0candidate.ExecutionPlanRoleBindings {
	t.Helper()
	return v0candidate.ExecutionPlanRoleBindings{
		InstallationID: repeatedBrokerID[v0candidate.InstallationID](0x11), EpochDigest: repeatedBrokerDigest[v0candidate.TrustEpochDigest](0x22),
		SourceManifestDigest:            brokerHexDigest[v0candidate.SourceManifestDigest](t, "c387c80094027ffbcacb573f44f5f6b4dec4d243bb436b24dd644434feaa1d14"),
		InlineInputDigest:               brokerHexDigest[v0candidate.InlineInputDigest](t, "bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72"),
		RuntimeBundleManifestDigest:     repeatedBrokerDigest[v0candidate.RuntimeBundleManifestDigest](0x55),
		ProfileReviewAttestationDigests: []v0candidate.ProfileReviewAttestationDigest{repeatedBrokerDigest[v0candidate.ProfileReviewAttestationDigest](0x66), repeatedBrokerDigest[v0candidate.ProfileReviewAttestationDigest](0x67)},
		ProfileRegistryEntryDigest:      repeatedBrokerDigest[v0candidate.ProfileRegistryEntryDigest](0x77), BackendValidationRecordDigest: repeatedBrokerDigest[v0candidate.BackendValidationRecordDigest](0x88),
		BackendConfigurationDigest: repeatedBrokerDigest[v0candidate.BackendConfigurationDigest](0x99), TrustSnapshotDigest: repeatedBrokerDigest[v0candidate.TrustSnapshotDigest](0xaa),
		PolicyDecisionDigest: repeatedBrokerDigest[v0candidate.PolicyDecisionDigest](0xbb),
	}
}

func brokerFixture(t *testing.T, path string) []byte {
	t.Helper()
	value, err := os.ReadFile(filepath.Join("..", "..", "..", "schemas", "conformance", "v0", path))
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func brokerHexDigest[T ~[32]byte](t *testing.T, value string) T {
	t.Helper()
	raw, err := hex.DecodeString(value)
	if err != nil {
		t.Fatal(err)
	}
	var result T
	copy(result[:], raw)
	return result
}

func repeatedBrokerID[T ~[16]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}

func repeatedBrokerDigest[T ~[32]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}

func assertProjectionClassification(t *testing.T, err error, expected Classification) {
	t.Helper()
	classification, ok := ErrorClassification(err)
	if !ok || classification != expected {
		t.Fatalf("classification = %q, error = %v, want %q", classification, err, expected)
	}
}
