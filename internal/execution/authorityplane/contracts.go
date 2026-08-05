// Package authorityplane implements the passive, unwired RegisterPlanV0 and
// GetRegisteredPlanV0 contract. It exposes no service or execution path.
package authorityplane

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	RoleBindingRecordVersion    = byte(0)
	RoleBindingRecordBytes      = 562
	RegisterPlanV0MaxBytes      = v0candidate.ExecutionPlanMaxCBORBytes + RoleBindingRecordBytes + v0candidate.SourceManifestMaxCBORBytes + v0candidate.MJSMainSourceMaxBytes
	GetRegisteredPlanV0MaxBytes = RegisterPlanV0MaxBytes + v0candidate.PlanRegistrationMaxCBORBytes
)

var (
	_ [328_337 - RegisterPlanV0MaxBytes]byte
	_ [RegisterPlanV0MaxBytes - 328_337]byte
	_ [332_433 - GetRegisteredPlanV0MaxBytes]byte
	_ [GetRegisteredPlanV0MaxBytes - 332_433]byte
)

type Classification string

const (
	Malformed      Classification = "MALFORMED"
	Unsupported    Classification = "UNSUPPORTED"
	Schema         Classification = "SCHEMA"
	Binding        Classification = "BINDING"
	Authentication Classification = "AUTHENTICATION"
	Capacity       Classification = "CAPACITY"
	LocalFailure   Classification = "LOCAL_FAILURE"
)

type contractError struct {
	classification Classification
	code           string
}

func (e *contractError) Error() string { return fmt.Sprintf("%s: %s", e.classification, e.code) }
func ErrorClassification(err error) (Classification, bool) {
	var target *contractError
	if errors.As(err, &target) {
		return target.classification, true
	}
	if value, ok := v0candidate.ErrorClassification(err); ok {
		return Classification(value), true
	}
	return "", false
}

// ErrorCode returns the fixed non-caller-controlled reason code for a
// recognized authority-plane refusal.
func ErrorCode(err error) (string, bool) {
	var target *contractError
	if !errors.As(err, &target) {
		return "", false
	}
	return target.code, true
}
func refused(classification Classification, code string) error {
	return &contractError{classification: classification, code: code}
}

// field-authority-object: capsule.register-plan-v0-request v0
// RegisterPlanV0Request is the complete application-data body. Every accessor
// returns a defensive copy; there is no path, URL, descriptor, or replacement source field.
type RegisterPlanV0Request struct {
	plan     []byte
	bindings []byte
	manifest []byte
	source   []byte
}

func NewRegisterPlanV0Request(plan, bindings, manifest, source []byte) (RegisterPlanV0Request, error) {
	if len(plan) == 0 || len(plan) > v0candidate.ExecutionPlanMaxCBORBytes {
		return RegisterPlanV0Request{}, refused(Malformed, "register-plan-bytes")
	}
	if len(bindings) != RoleBindingRecordBytes {
		return RegisterPlanV0Request{}, refused(Malformed, "register-role-bindings")
	}
	if len(manifest) < v0candidate.SourceManifestMinCBORBytes || len(manifest) > v0candidate.SourceManifestMaxCBORBytes {
		return RegisterPlanV0Request{}, refused(Malformed, "register-source-manifest")
	}
	if len(source) > v0candidate.MJSMainSourceMaxBytes {
		return RegisterPlanV0Request{}, refused(Malformed, "register-source-bytes")
	}
	if len(plan)+len(bindings)+len(manifest)+len(source) > RegisterPlanV0MaxBytes {
		return RegisterPlanV0Request{}, refused(Malformed, "register-aggregate-bytes")
	}
	return RegisterPlanV0Request{
		plan: bytes.Clone(plan), bindings: bytes.Clone(bindings),
		manifest: bytes.Clone(manifest), source: bytes.Clone(source),
	}, nil
}

func (r RegisterPlanV0Request) ExactPlanBytes() []byte          { return bytes.Clone(r.plan) }
func (r RegisterPlanV0Request) NominalRoleBindingBytes() []byte { return bytes.Clone(r.bindings) }
func (r RegisterPlanV0Request) SourceManifestBytes() []byte     { return bytes.Clone(r.manifest) }
func (r RegisterPlanV0Request) SourceBytes() []byte             { return bytes.Clone(r.source) }

type GetRegisteredPlanV0Request struct{ RegistrationID v0candidate.RegistrationID }

// field-authority-object: capsule.get-registered-plan-v0-reply v0
type GetRegisteredPlanV0Reply struct {
	plan         []byte
	bindings     []byte
	registration []byte
	manifest     []byte
	source       []byte
}

func (r GetRegisteredPlanV0Reply) ExactPlanBytes() []byte           { return bytes.Clone(r.plan) }
func (r GetRegisteredPlanV0Reply) ResolvedRoleBindingBytes() []byte { return bytes.Clone(r.bindings) }
func (r GetRegisteredPlanV0Reply) PlanRegistrationBytes() []byte    { return bytes.Clone(r.registration) }
func (r GetRegisteredPlanV0Reply) SourceManifestBytes() []byte      { return bytes.Clone(r.manifest) }
func (r GetRegisteredPlanV0Reply) SourceBytes() []byte              { return bytes.Clone(r.source) }

// EncodeRoleBindingsV0 emits the fixed bridge-only record. Unused review slots are zero.
func EncodeRoleBindingsV0(value v0candidate.ExecutionPlanRoleBindings) ([]byte, error) {
	if len(value.ProfileReviewAttestationDigests) > 8 {
		return nil, refused(Schema, "role-binding-review-count")
	}
	result := make([]byte, RoleBindingRecordBytes)
	result[0] = RoleBindingRecordVersion
	offset := 1
	copy(result[offset:], value.InstallationID[:])
	offset += 16
	for _, digest := range [][]byte{value.EpochDigest[:], value.SourceManifestDigest[:], value.InlineInputDigest[:], value.RuntimeBundleManifestDigest[:]} {
		copy(result[offset:], digest)
		offset += 32
	}
	var reviewCount byte
	for range value.ProfileReviewAttestationDigests {
		reviewCount++
	}
	result[offset] = reviewCount
	offset++
	for _, digest := range value.ProfileReviewAttestationDigests {
		copy(result[offset:], digest[:])
		offset += 32
	}
	offset = 1 + 16 + 4*32 + 1 + 8*32
	for _, digest := range [][]byte{value.ProfileRegistryEntryDigest[:], value.BackendValidationRecordDigest[:], value.BackendConfigurationDigest[:], value.TrustSnapshotDigest[:], value.PolicyDecisionDigest[:]} {
		copy(result[offset:], digest)
		offset += 32
	}
	if offset != RoleBindingRecordBytes {
		panic("role-binding layout drift")
	}
	return result, nil
}

func DecodeRoleBindingsV0(received []byte) (v0candidate.ExecutionPlanRoleBindings, error) {
	if len(received) != RoleBindingRecordBytes {
		return v0candidate.ExecutionPlanRoleBindings{}, refused(Malformed, "role-binding-length")
	}
	if received[0] != RoleBindingRecordVersion {
		return v0candidate.ExecutionPlanRoleBindings{}, refused(Unsupported, "role-binding-version")
	}
	offset := 1
	installation, err := v0candidate.NewInstallationID(received[offset : offset+16])
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, err
	}
	offset += 16
	readDigest := func() []byte { value := received[offset : offset+32]; offset += 32; return value }
	epoch, err := v0candidate.NewTrustEpochDigest(readDigest())
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, err
	}
	source, err := v0candidate.NewSourceManifestDigest(readDigest())
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, err
	}
	input, err := v0candidate.NewInlineInputDigest(readDigest())
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, err
	}
	runtime, err := v0candidate.NewRuntimeBundleManifestDigest(readDigest())
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, err
	}
	count := int(received[offset])
	offset++
	if count > 8 {
		return v0candidate.ExecutionPlanRoleBindings{}, refused(Schema, "role-binding-review-count")
	}
	reviews := make([]v0candidate.ProfileReviewAttestationDigest, 0, count)
	for index := 0; index < 8; index++ {
		raw := readDigest()
		if index < count {
			value, err := v0candidate.NewProfileReviewAttestationDigest(raw)
			if err != nil {
				return v0candidate.ExecutionPlanRoleBindings{}, err
			}
			reviews = append(reviews, value)
		} else if !bytes.Equal(raw, make([]byte, 32)) {
			return v0candidate.ExecutionPlanRoleBindings{}, refused(Schema, "role-binding-unused-review-slot")
		}
	}
	profile, err := v0candidate.NewProfileRegistryEntryDigest(readDigest())
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, err
	}
	backendValidation, err := v0candidate.NewBackendValidationRecordDigest(readDigest())
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, err
	}
	backendConfiguration, err := v0candidate.NewBackendConfigurationDigest(readDigest())
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, err
	}
	trust, err := v0candidate.NewTrustSnapshotDigest(readDigest())
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, err
	}
	policy, err := v0candidate.NewPolicyDecisionDigest(readDigest())
	if err != nil {
		return v0candidate.ExecutionPlanRoleBindings{}, err
	}
	if offset != len(received) {
		return v0candidate.ExecutionPlanRoleBindings{}, refused(Malformed, "role-binding-trailing-data")
	}
	return v0candidate.ExecutionPlanRoleBindings{InstallationID: installation, EpochDigest: epoch, SourceManifestDigest: source, InlineInputDigest: input, RuntimeBundleManifestDigest: runtime, ProfileReviewAttestationDigests: reviews, ProfileRegistryEntryDigest: profile, BackendValidationRecordDigest: backendValidation, BackendConfigurationDigest: backendConfiguration, TrustSnapshotDigest: trust, PolicyDecisionDigest: policy}, nil
}

func appendCBORArgument(destination []byte, major byte, value uint64) []byte {
	switch {
	case value < 24:
		return append(destination, major<<5|byte(value))
	case value <= 0xff:
		return append(destination, major<<5|24, byte(value))
	case value <= 0xffff:
		return binary.BigEndian.AppendUint16(append(destination, major<<5|25), uint16(value))
	case value <= 0xffffffff:
		return binary.BigEndian.AppendUint32(append(destination, major<<5|26), uint32(value))
	default:
		return binary.BigEndian.AppendUint64(append(destination, major<<5|27), value)
	}
}
func appendUnsigned(destination []byte, value uint64) []byte {
	return appendCBORArgument(destination, 0, value)
}
func appendText(destination []byte, value string) []byte {
	return append(appendCBORArgument(destination, 3, uint64(len(value))), value...)
}
func appendBytes(destination, value []byte) []byte {
	return append(appendCBORArgument(destination, 2, uint64(len(value))), value...)
}
