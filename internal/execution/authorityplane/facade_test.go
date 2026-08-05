package authorityplane

import (
	"bytes"
	"context"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

type exactRoleResolver struct {
	expected v0candidate.ExecutionPlanRoleBindings
}

func (r exactRoleResolver) ResolvePlanRoles(_ context.Context, received v0candidate.ExecutionPlanRoleBindings) (v0candidate.ExecutionPlanRoleBindings, error) {
	if !EqualRoleBindings(received, r.expected) {
		return v0candidate.ExecutionPlanRoleBindings{}, refused(Binding, "nominal-role-mismatch")
	}
	return r.expected, nil
}

type fixedIdentifier struct{ value v0candidate.RegistrationID }

func (i fixedIdentifier) NewRegistrationID(context.Context) (v0candidate.RegistrationID, error) {
	return i.value, nil
}

func TestPassiveAtomicRegistrationAndBrokerReadback(t *testing.T) {
	bindings := ordinaryBindings(t)
	bindingBytes, err := EncodeRoleBindingsV0(bindings)
	if err != nil {
		t.Fatal(err)
	}
	if len(bindingBytes) != 562 {
		t.Fatalf("binding bytes = %d", len(bindingBytes))
	}
	if expected := readAuthorityFixture(t, "role-bindings.bin"); !bytes.Equal(bindingBytes, expected) {
		t.Fatal("Go role-binding encoding disagrees with generated cross-language known answer")
	}
	plan := readFixture(t, "execution-plan/ordinary.cbor")
	manifest := readFixture(t, "source-manifest/ordinary.cbor")
	source := []byte("export default function (value) { return value; }\n")
	request, err := NewRegisterPlanV0Request(plan, bindingBytes, manifest, source)
	if err != nil {
		t.Fatal(err)
	}
	registrationID := repeatedID[v0candidate.RegistrationID](0x7a)
	supervisor := SupervisorContext{InstallationID: bindings.InstallationID, EpochSequence: 7, EpochDigest: bindings.EpochDigest, SupervisorID: repeatedID[v0candidate.SupervisorID](0x55)}
	store := &FixedStore{}
	facade, err := NewFacade(store, exactRoleResolver{bindings}, fixedIdentifier{registrationID}, supervisor)
	if err != nil {
		t.Fatal(err)
	}

	registration, err := facade.RegisterPlanV0(context.Background(), CallContext{Authenticated: true, Role: Daemon, Purpose: RegisterPlanV0Purpose}, request)
	if err != nil {
		t.Fatal(err)
	}
	if store.count() != 1 {
		t.Fatalf("registration count = %d", store.count())
	}
	reply, err := facade.GetRegisteredPlanV0(context.Background(), CallContext{Authenticated: true, Role: Broker, Purpose: GetRegisteredPlanV0Purpose}, GetRegisteredPlanV0Request{RegistrationID: registrationID})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(reply.ExactPlanBytes(), plan) || !bytes.Equal(reply.ResolvedRoleBindingBytes(), bindingBytes) || !bytes.Equal(reply.PlanRegistrationBytes(), registration) || !bytes.Equal(reply.SourceManifestBytes(), manifest) || !bytes.Equal(reply.SourceBytes(), source) {
		t.Fatal("Broker readback differs from atomic retained bytes")
	}

	plan[0] ^= 0xff
	manifest[0] ^= 0xff
	source[0] ^= 0xff
	bindingBytes[0] ^= 0xff
	registration[0] ^= 0xff
	again, err := facade.GetRegisteredPlanV0(context.Background(), CallContext{Authenticated: true, Role: Broker, Purpose: GetRegisteredPlanV0Purpose}, GetRegisteredPlanV0Request{RegistrationID: registrationID})
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(again.ExactPlanBytes(), plan) || bytes.Equal(again.SourceBytes(), source) {
		t.Fatal("caller mutation changed retained custody")
	}
	returned := again.SourceBytes()
	returned[0] ^= 0xff
	third, err := facade.GetRegisteredPlanV0(context.Background(), CallContext{Authenticated: true, Role: Broker, Purpose: GetRegisteredPlanV0Purpose}, GetRegisteredPlanV0Request{RegistrationID: registrationID})
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(returned, third.SourceBytes()) {
		t.Fatal("reply mutation changed retained custody")
	}
}

func TestRegistrationCommitFailureLeavesNoSplitCustody(t *testing.T) {
	bindings := ordinaryBindings(t)
	bindingBytes, _ := EncodeRoleBindingsV0(bindings)
	request, err := NewRegisterPlanV0Request(readFixture(t, "execution-plan/ordinary.cbor"), bindingBytes, readFixture(t, "source-manifest/ordinary.cbor"), []byte("export default function (value) { return value; }\n"))
	if err != nil {
		t.Fatal(err)
	}
	store := &FixedStore{}
	store.setCommitFailureForTest(true)
	facade, err := NewFacade(store, exactRoleResolver{bindings}, fixedIdentifier{repeatedID[v0candidate.RegistrationID](0x7a)}, SupervisorContext{InstallationID: bindings.InstallationID, EpochSequence: 7, EpochDigest: bindings.EpochDigest, SupervisorID: repeatedID[v0candidate.SupervisorID](0x55)})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := facade.RegisterPlanV0(context.Background(), CallContext{Authenticated: true, Role: Daemon, Purpose: RegisterPlanV0Purpose}, request); err == nil {
		t.Fatal("commit failure accepted")
	}
	if store.count() != 0 {
		t.Fatal("failed transaction retained partial custody")
	}
}

func TestSourceBindingFailureLeavesNoCustody(t *testing.T) {
	bindings := ordinaryBindings(t)
	bindingBytes, err := EncodeRoleBindingsV0(bindings)
	if err != nil {
		t.Fatal(err)
	}
	request, err := NewRegisterPlanV0Request(
		readFixture(t, "execution-plan/ordinary.cbor"),
		bindingBytes,
		readFixture(t, "source-manifest/ordinary.cbor"),
		[]byte("export default function (value) { return value + 1; }\n"),
	)
	if err != nil {
		t.Fatal(err)
	}
	store := &FixedStore{}
	facade, err := NewFacade(
		store,
		exactRoleResolver{bindings},
		fixedIdentifier{repeatedID[v0candidate.RegistrationID](0x7a)},
		SupervisorContext{
			InstallationID: bindings.InstallationID,
			EpochSequence:  7,
			EpochDigest:    bindings.EpochDigest,
			SupervisorID:   repeatedID[v0candidate.SupervisorID](0x55),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := facade.RegisterPlanV0(
		context.Background(),
		CallContext{Authenticated: true, Role: Daemon, Purpose: RegisterPlanV0Purpose},
		request,
	); err == nil {
		t.Fatal("source replacement accepted")
	}
	if store.count() != 0 {
		t.Fatal("source-binding refusal retained partial custody")
	}
}

func TestMethodAuthenticationRefusesBeforeApplicationData(t *testing.T) {
	bindings := ordinaryBindings(t)
	store := &FixedStore{}
	facade, err := NewFacade(
		store,
		exactRoleResolver{bindings},
		fixedIdentifier{repeatedID[v0candidate.RegistrationID](0x7a)},
		SupervisorContext{
			InstallationID: bindings.InstallationID,
			EpochSequence:  7,
			EpochDigest:    bindings.EpochDigest,
			SupervisorID:   repeatedID[v0candidate.SupervisorID](0x55),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	_, registerErr := facade.RegisterPlanV0(context.Background(), CallContext{}, RegisterPlanV0Request{})
	if classification, ok := ErrorClassification(registerErr); !ok || classification != Authentication {
		t.Fatalf("unauthenticated registration classification = %q, %v", classification, registerErr)
	}
	_, fetchErr := facade.GetRegisteredPlanV0(
		context.Background(),
		CallContext{Authenticated: true, Role: Daemon, Purpose: GetRegisteredPlanV0Purpose},
		GetRegisteredPlanV0Request{},
	)
	if classification, ok := ErrorClassification(fetchErr); !ok || classification != Authentication {
		t.Fatalf("wrong-role fetch classification = %q, %v", classification, fetchErr)
	}
	if store.count() != 0 {
		t.Fatal("authentication refusal changed custody")
	}
}

func TestClosedRoleRecordRefuses(t *testing.T) {
	bindings := ordinaryBindings(t)
	encoded, err := EncodeRoleBindingsV0(bindings)
	if err != nil {
		t.Fatal(err)
	}
	mutated := bytes.Clone(encoded)
	mutated[0] = 1
	if _, err := DecodeRoleBindingsV0(mutated); err == nil {
		t.Fatal("unknown version accepted")
	}
	mutated = bytes.Clone(encoded)
	mutated[145] = 9
	if _, err := DecodeRoleBindingsV0(mutated); err == nil {
		t.Fatal("review count above eight accepted")
	}
	mutated = bytes.Clone(encoded)
	mutated[1+16+4*32+1+2*32] = 1
	if _, err := DecodeRoleBindingsV0(mutated); err == nil {
		t.Fatal("nonzero unused review slot accepted")
	}
	if RegisterPlanV0MaxBytes != 328337 || GetRegisteredPlanV0MaxBytes != 332433 {
		t.Fatal("generated cap drift")
	}
}

func ordinaryBindings(t *testing.T) v0candidate.ExecutionPlanRoleBindings {
	t.Helper()
	return v0candidate.ExecutionPlanRoleBindings{
		InstallationID: repeatedID[v0candidate.InstallationID](0x11), EpochDigest: repeatedDigest[v0candidate.TrustEpochDigest](0x22),
		SourceManifestDigest: hexDigest[v0candidate.SourceManifestDigest](t, "c387c80094027ffbcacb573f44f5f6b4dec4d243bb436b24dd644434feaa1d14"), InlineInputDigest: hexDigest[v0candidate.InlineInputDigest](t, "bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72"),
		RuntimeBundleManifestDigest: repeatedDigest[v0candidate.RuntimeBundleManifestDigest](0x55), ProfileReviewAttestationDigests: []v0candidate.ProfileReviewAttestationDigest{repeatedDigest[v0candidate.ProfileReviewAttestationDigest](0x66), repeatedDigest[v0candidate.ProfileReviewAttestationDigest](0x67)},
		ProfileRegistryEntryDigest: repeatedDigest[v0candidate.ProfileRegistryEntryDigest](0x77), BackendValidationRecordDigest: repeatedDigest[v0candidate.BackendValidationRecordDigest](0x88), BackendConfigurationDigest: repeatedDigest[v0candidate.BackendConfigurationDigest](0x99), TrustSnapshotDigest: repeatedDigest[v0candidate.TrustSnapshotDigest](0xaa), PolicyDecisionDigest: repeatedDigest[v0candidate.PolicyDecisionDigest](0xbb),
	}
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	value, err := os.ReadFile(filepath.Join("..", "..", "..", "schemas", "conformance", "v0", name))
	if err != nil {
		t.Fatal(err)
	}
	return value
}
func readAuthorityFixture(t *testing.T, name string) []byte {
	t.Helper()
	value, err := os.ReadFile(filepath.Join("..", "..", "..", "schemas", "conformance", "authority-plane-v0", name))
	if err != nil {
		t.Fatal(err)
	}
	return value
}
func repeatedID[T ~[16]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}
func repeatedDigest[T ~[32]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}
func hexDigest[T ~[32]byte](t *testing.T, value string) T {
	t.Helper()
	raw, err := hex.DecodeString(value)
	if err != nil {
		t.Fatal(err)
	}
	var result T
	copy(result[:], raw)
	return result
}
