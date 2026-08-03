package lifecyclestate

import (
	"bytes"
	"reflect"

	"capsule.local/capsule/internal/execution/approvalattempt"
	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	MaxProfileReviewAttestationDigests = 8
	MaxBackendInstanceIdentityBytes    = 64
	FakeBackendProtocolVersion         = v0candidate.PositiveUInt53(1)
	SliceBRegistrationStorageVersion   = uint64(0)
	SliceBApprovalStorageVersion       = uint64(0)
	SliceBAttemptStorageVersion        = uint64(0)
)

// BackendBindingView is the closed fake-only backend projection retained by
// one lifecycle record. CreatesGuest must remain false in this slice.
type BackendBindingView struct {
	Kind                          BackendKind
	ProtocolVersion               v0candidate.PositiveUInt53
	ImplementationIdentityDigest  BackendImplementationDigest
	BackendConfigurationDigest    v0candidate.BackendConfigurationDigest
	BackendValidationRecordDigest v0candidate.BackendValidationRecordDigest
	CreatesGuest                  bool
}

// BackendBinding is an immutable, validated backend projection.
type BackendBinding struct {
	view BackendBindingView
}

func NewBackendBinding(view BackendBindingView) (BackendBinding, error) {
	if view.Kind != BackendFakeNoGuest {
		return BackendBinding{}, classified(ClassificationUnsupported, "backend-kind")
	}
	if view.ProtocolVersion != FakeBackendProtocolVersion {
		return BackendBinding{}, classified(ClassificationUnsupported, "backend-protocol-version")
	}
	if zeroDigest(view.ImplementationIdentityDigest) ||
		zeroDigest(view.BackendConfigurationDigest) ||
		zeroDigest(view.BackendValidationRecordDigest) {
		return BackendBinding{}, classified(ClassificationSchema, "backend-digest-zero")
	}
	if view.CreatesGuest {
		return BackendBinding{}, classified(ClassificationUnsupported, "backend-creates-guest")
	}
	return BackendBinding{view: view}, nil
}

func (binding BackendBinding) View() BackendBindingView { return binding.view }

func (binding BackendBinding) valid() bool {
	validated, err := NewBackendBinding(binding.view)
	return err == nil && validated.view == binding.view
}

// BackendInstanceIdentityView is the serializable projection of a bounded
// opaque instance identity and its deterministic domain-separated digest.
type BackendInstanceIdentityView struct {
	Kind   BackendInstanceKind
	Value  []byte
	Digest BackendInstanceDigest
}

// BackendInstanceIdentity owns a defensive copy of bounded opaque bytes.
type BackendInstanceIdentity struct {
	kind   BackendInstanceKind
	value  []byte
	digest BackendInstanceDigest
}

func NewBackendInstanceIdentity(
	kind BackendInstanceKind,
	value []byte,
) (BackendInstanceIdentity, error) {
	if kind != BackendInstanceFake {
		return BackendInstanceIdentity{}, classified(ClassificationUnsupported, "backend-instance-kind")
	}
	if len(value) == 0 || len(value) > MaxBackendInstanceIdentityBytes {
		return BackendInstanceIdentity{}, classified(ClassificationSchema, "backend-instance-length")
	}
	copied := bytes.Clone(value)
	encoder := newDigestEncoder("capsule.lifecycle.backend-instance.v1")
	encoder.text(string(kind))
	encoder.bytes(copied)
	digest := digestSum[BackendInstanceDigest](encoder)
	return BackendInstanceIdentity{kind: kind, value: copied, digest: digest}, nil
}

func RestoreBackendInstanceIdentity(view BackendInstanceIdentityView) (BackendInstanceIdentity, error) {
	identity, err := NewBackendInstanceIdentity(view.Kind, view.Value)
	if err != nil {
		return BackendInstanceIdentity{}, err
	}
	if view.Digest != identity.digest {
		return BackendInstanceIdentity{}, classified(ClassificationBinding, "backend-instance-digest")
	}
	return identity, nil
}

func (identity BackendInstanceIdentity) Present() bool {
	return identity.kind != "" || len(identity.value) != 0 || identity.digest != (BackendInstanceDigest{})
}

func (identity BackendInstanceIdentity) Kind() BackendInstanceKind { return identity.kind }

func (identity BackendInstanceIdentity) Value() []byte { return bytes.Clone(identity.value) }

func (identity BackendInstanceIdentity) Digest() BackendInstanceDigest { return identity.digest }

func (identity BackendInstanceIdentity) View() BackendInstanceIdentityView {
	return BackendInstanceIdentityView{
		Kind: identity.kind, Value: bytes.Clone(identity.value), Digest: identity.digest,
	}
}

func (identity BackendInstanceIdentity) valid() bool {
	if !identity.Present() {
		return true
	}
	restored, err := RestoreBackendInstanceIdentity(identity.View())
	return err == nil && restored.digest == identity.digest
}

// ImmutableBindingsView contains every copied Slice B authority and plan
// binding required by ADR-0025. Profile review digests retain their order.
type ImmutableBindingsView struct {
	AttemptID             approvalattempt.AttemptID
	ApprovalID            approvalattempt.ApprovalID
	AttemptNonce          approvalattempt.AttemptNonce
	RegistrationID        v0candidate.RegistrationID
	RegistrationSequence  v0candidate.PositiveUInt53
	PlanDigest            v0candidate.ExecutionPlanDigest
	InstallationID        v0candidate.InstallationID
	EpochSequence         v0candidate.UInt53
	EpochDigest           v0candidate.TrustEpochDigest
	SupervisorID          v0candidate.SupervisorID
	ApprovalPurpose       string
	ApprovalAudience      string
	ApprovalPayloadDigest approvalattempt.ApprovalPayloadDigest
	AuthorizationIdentity approvalattempt.ApprovalKeyAuthorizationIdentity
	AttemptCreatedAt      v0candidate.UInt53

	AttemptStorageVersion      uint64
	ApprovalStorageVersion     uint64
	RegistrationStorageVersion uint64

	SourceManifestDigest            v0candidate.SourceManifestDigest
	InlineInputDigest               v0candidate.InlineInputDigest
	RuntimeBundleManifestDigest     v0candidate.RuntimeBundleManifestDigest
	ProfileReviewAttestationDigests []v0candidate.ProfileReviewAttestationDigest
	ProfileRegistryEntryDigest      v0candidate.ProfileRegistryEntryDigest
	BackendValidationRecordDigest   v0candidate.BackendValidationRecordDigest
	BackendConfigurationDigest      v0candidate.BackendConfigurationDigest
	TrustSnapshotDigest             v0candidate.TrustSnapshotDigest
	PolicyDecisionDigest            v0candidate.PolicyDecisionDigest

	Backend BackendBinding
}

// ImmutableBindings owns the copied projection and its deterministic digest.
// Its only slice is private and every projection returns another copy.
type ImmutableBindings struct {
	view   ImmutableBindingsView
	digest ImmutableBindingDigest
}

func NewImmutableBindings(view ImmutableBindingsView) (ImmutableBindings, error) {
	view.ProfileReviewAttestationDigests = append(
		[]v0candidate.ProfileReviewAttestationDigest(nil),
		view.ProfileReviewAttestationDigests...,
	)
	if err := validateImmutableBindingsView(view); err != nil {
		return ImmutableBindings{}, err
	}
	binding := ImmutableBindings{view: view}
	binding.digest = digestImmutableBindings(view)
	return binding, nil
}

func RestoreImmutableBindings(
	view ImmutableBindingsView,
	digest ImmutableBindingDigest,
) (ImmutableBindings, error) {
	binding, err := NewImmutableBindings(view)
	if err != nil {
		return ImmutableBindings{}, err
	}
	if binding.digest != digest {
		return ImmutableBindings{}, classified(ClassificationBinding, "immutable-binding-digest")
	}
	return binding, nil
}

func (binding ImmutableBindings) View() ImmutableBindingsView {
	view := binding.view
	view.ProfileReviewAttestationDigests = append(
		[]v0candidate.ProfileReviewAttestationDigest(nil),
		binding.view.ProfileReviewAttestationDigests...,
	)
	return view
}

func (binding ImmutableBindings) Digest() ImmutableBindingDigest { return binding.digest }

func (binding ImmutableBindings) Equal(other ImmutableBindings) bool {
	return binding.digest == other.digest && reflect.DeepEqual(binding.View(), other.View())
}

func (binding ImmutableBindings) valid() bool {
	restored, err := RestoreImmutableBindings(binding.View(), binding.digest)
	return err == nil && restored.Equal(binding)
}

func validateImmutableBindingsView(view ImmutableBindingsView) error {
	if view.AttemptID == (approvalattempt.AttemptID{}) ||
		view.ApprovalID == (approvalattempt.ApprovalID{}) ||
		view.AttemptNonce == (approvalattempt.AttemptNonce{}) ||
		view.RegistrationID == (v0candidate.RegistrationID{}) ||
		view.InstallationID == (v0candidate.InstallationID{}) ||
		view.SupervisorID == (v0candidate.SupervisorID{}) {
		return classified(ClassificationSchema, "immutable-identifier-zero")
	}
	if view.RegistrationSequence == 0 ||
		uint64(view.RegistrationSequence) > v0candidate.MaxSafeInteger ||
		uint64(view.EpochSequence) > v0candidate.MaxSafeInteger ||
		uint64(view.AttemptCreatedAt) > v0candidate.MaxSafeInteger {
		return classified(ClassificationSchema, "immutable-uint53")
	}
	if view.ApprovalPurpose != approvalattempt.ApprovalGrantPurpose ||
		view.ApprovalAudience != approvalattempt.ApprovalGrantAudience {
		return classified(ClassificationBinding, "approval-purpose-audience")
	}
	if view.AttemptStorageVersion != SliceBAttemptStorageVersion ||
		view.ApprovalStorageVersion != SliceBApprovalStorageVersion ||
		view.RegistrationStorageVersion != SliceBRegistrationStorageVersion {
		return classified(ClassificationUnsupported, "slice-b-storage-version")
	}
	if zeroDigest(view.PlanDigest) || zeroDigest(view.EpochDigest) ||
		zeroDigest(view.ApprovalPayloadDigest) || zeroDigest(view.AuthorizationIdentity) ||
		zeroDigest(view.SourceManifestDigest) || zeroDigest(view.InlineInputDigest) ||
		zeroDigest(view.RuntimeBundleManifestDigest) ||
		zeroDigest(view.ProfileRegistryEntryDigest) ||
		zeroDigest(view.BackendValidationRecordDigest) ||
		zeroDigest(view.BackendConfigurationDigest) ||
		zeroDigest(view.TrustSnapshotDigest) || zeroDigest(view.PolicyDecisionDigest) {
		return classified(ClassificationSchema, "immutable-digest-zero")
	}
	if len(view.ProfileReviewAttestationDigests) == 0 ||
		len(view.ProfileReviewAttestationDigests) > MaxProfileReviewAttestationDigests {
		return classified(ClassificationCapacity, "profile-review-digest-count")
	}
	for _, digest := range view.ProfileReviewAttestationDigests {
		if zeroDigest(digest) {
			return classified(ClassificationSchema, "profile-review-digest-zero")
		}
	}
	if !view.Backend.valid() {
		return classified(ClassificationUnsupported, "backend-binding")
	}
	backend := view.Backend.View()
	if backend.BackendConfigurationDigest != view.BackendConfigurationDigest ||
		backend.BackendValidationRecordDigest != view.BackendValidationRecordDigest {
		return classified(ClassificationBinding, "backend-plan-cross-link")
	}
	return nil
}

func digestImmutableBindings(view ImmutableBindingsView) ImmutableBindingDigest {
	encoder := newDigestEncoder("capsule.lifecycle.immutable-bindings.v1")
	encoder.bytes(view.AttemptID[:])
	encoder.bytes(view.ApprovalID[:])
	encoder.bytes(view.AttemptNonce[:])
	encoder.bytes(view.RegistrationID[:])
	encoder.uint64(uint64(view.RegistrationSequence))
	encoder.bytes(view.PlanDigest[:])
	encoder.bytes(view.InstallationID[:])
	encoder.uint64(uint64(view.EpochSequence))
	encoder.bytes(view.EpochDigest[:])
	encoder.bytes(view.SupervisorID[:])
	encoder.text(view.ApprovalPurpose)
	encoder.text(view.ApprovalAudience)
	encoder.bytes(view.ApprovalPayloadDigest[:])
	encoder.bytes(view.AuthorizationIdentity[:])
	encoder.uint64(uint64(view.AttemptCreatedAt))
	encoder.uint64(view.AttemptStorageVersion)
	encoder.uint64(view.ApprovalStorageVersion)
	encoder.uint64(view.RegistrationStorageVersion)
	encoder.bytes(view.SourceManifestDigest[:])
	encoder.bytes(view.InlineInputDigest[:])
	encoder.bytes(view.RuntimeBundleManifestDigest[:])
	encoder.uint64(uint64(len(view.ProfileReviewAttestationDigests)))
	for _, digest := range view.ProfileReviewAttestationDigests {
		encoder.bytes(digest[:])
	}
	encoder.bytes(view.ProfileRegistryEntryDigest[:])
	encoder.bytes(view.BackendValidationRecordDigest[:])
	encoder.bytes(view.BackendConfigurationDigest[:])
	encoder.bytes(view.TrustSnapshotDigest[:])
	encoder.bytes(view.PolicyDecisionDigest[:])
	backend := view.Backend.View()
	encoder.text(string(backend.Kind))
	encoder.uint64(uint64(backend.ProtocolVersion))
	encoder.bytes(backend.ImplementationIdentityDigest[:])
	encoder.bytes(backend.BackendConfigurationDigest[:])
	encoder.bytes(backend.BackendValidationRecordDigest[:])
	encoder.boolean(backend.CreatesGuest)
	return digestSum[ImmutableBindingDigest](encoder)
}
