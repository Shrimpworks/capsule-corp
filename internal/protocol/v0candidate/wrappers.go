package v0candidate

import (
	"bytes"
	"crypto/sha256"
	"strings"
)

const (
	ExecutionPlanMediaType    = "application/capsule.execution-plan+cbor;v=0"
	PlanRegistrationMediaType = "application/capsule.plan-registration+cbor;v=0"
)

// ExecutionPlanRoleBindings are trusted, role-resolved identities expected in
// one candidate plan. CBOR byte strings do not carry nominal roles, so callers
// must supply this context to reject correct-width cross-domain substitutions.
type ExecutionPlanRoleBindings struct {
	InstallationID                  InstallationID
	EpochDigest                     TrustEpochDigest
	SourceManifestDigest            SourceManifestDigest
	InlineInputDigest               InlineInputDigest
	RuntimeBundleManifestDigest     RuntimeBundleManifestDigest
	ProfileReviewAttestationDigests []ProfileReviewAttestationDigest
	ProfileRegistryEntryDigest      ProfileRegistryEntryDigest
	BackendValidationRecordDigest   BackendValidationRecordDigest
	BackendConfigurationDigest      BackendConfigurationDigest
	TrustSnapshotDigest             TrustSnapshotDigest
	PolicyDecisionDigest            PolicyDecisionDigest
}

// PlanRegistrationRoleBindings are trusted identities expected in one
// Supervisor-issued registration response.
type PlanRegistrationRoleBindings struct {
	RegistrationID RegistrationID
	PlanDigest     ExecutionPlanDigest
	InstallationID InstallationID
	EpochDigest    TrustEpochDigest
	SupervisorID   SupervisorID
}

// DecodedExecutionPlan retains an immutable private copy of the exact received
// authoritative bytes alongside the decoded passive view.
type DecodedExecutionPlan struct {
	authoritativeBytes []byte
	digest             ExecutionPlanDigest
	view               ExecutionPlan
}

// AuthoritativeBytes returns a defensive copy of the exact received bytes.
func (p *DecodedExecutionPlan) AuthoritativeBytes() []byte {
	return bytes.Clone(p.authoritativeBytes)
}

// Digest returns SHA-256 over the exact received authoritative bytes.
func (p *DecodedExecutionPlan) Digest() ExecutionPlanDigest {
	return p.digest
}

// View returns a copy of the decoded passive candidate view.
func (p *DecodedExecutionPlan) View() ExecutionPlan {
	view := p.view
	view.ProfileReviewAttestationDigests = append(
		[]ProfileReviewAttestationDigest(nil),
		p.view.ProfileReviewAttestationDigests...,
	)
	return view
}

// DecodedPlanRegistration retains an immutable private copy of the exact
// received registration bytes alongside the decoded passive view.
type DecodedPlanRegistration struct {
	authoritativeBytes []byte
	view               PlanRegistration
}

// AuthoritativeBytes returns a defensive copy of the exact received bytes.
func (r *DecodedPlanRegistration) AuthoritativeBytes() []byte {
	return bytes.Clone(r.authoritativeBytes)
}

// View returns a copy of the decoded passive candidate view.
func (r *DecodedPlanRegistration) View() PlanRegistration {
	return r.view
}

// ValidateExecutionPlanMediaType validates the object-specific conformance
// media type. The typed local IPC method remains an independent binding.
func ValidateExecutionPlanMediaType(value string) error {
	return validateMediaType(value, "application/capsule.execution-plan+cbor")
}

// ValidatePlanRegistrationMediaType validates the object-specific conformance
// media type. The typed local IPC method remains an independent binding.
func ValidatePlanRegistrationMediaType(value string) error {
	return validateMediaType(value, "application/capsule.plan-registration+cbor")
}

func validateMediaType(value, expectedType string) error {
	separator := strings.IndexByte(value, ';')
	if separator <= 0 || separator != strings.LastIndexByte(value, ';') || separator == len(value)-1 {
		return malformed("media type must contain exactly one version parameter")
	}
	if !strings.EqualFold(value[:separator], expectedType) {
		return malformed("unexpected media type")
	}
	parameter := value[separator+1:]
	if parameter == "v=0" {
		return nil
	}
	if strings.HasPrefix(parameter, "v=") {
		return unsupported("unsupported media-type version")
	}
	return malformed("media type must contain the exact version parameter")
}

// DecodeExecutionPlan copies, predecodes, strictly decodes, and role-binds one
// passive ExecutionPlan candidate. It creates no registration or authority
// state.
func DecodeExecutionPlan(received []byte, bindings ExecutionPlanRoleBindings) (*DecodedExecutionPlan, error) {
	if len(received) > ExecutionPlanMaxCBORBytes {
		return nil, malformed("CBOR raw-byte limit exceeded")
	}
	authoritative := bytes.Clone(received)
	if hasObjectType(authoritative, PlanRegistrationObjectType) {
		if err := PredecodePlanRegistrationCBOR(authoritative); err != nil {
			return nil, err
		}
		if _, err := decodePlanRegistration(authoritative); err != nil {
			return nil, err
		}
		return nil, domain("PlanRegistration supplied at ExecutionPlan entry point")
	}
	if err := PredecodeExecutionPlanCBOR(authoritative); err != nil {
		return nil, err
	}
	view, err := decodeExecutionPlan(authoritative)
	if err != nil {
		return nil, err
	}
	if err := bindExecutionPlan(view, bindings); err != nil {
		return nil, err
	}
	digest := sha256.Sum256(authoritative)
	return &DecodedExecutionPlan{
		authoritativeBytes: authoritative,
		digest:             ExecutionPlanDigest(digest),
		view:               view,
	}, nil
}

// DecodePlanRegistration copies, predecodes, strictly decodes, and role-binds
// one passive Supervisor-issued PlanRegistration candidate. It creates no
// registration or authority state.
func DecodePlanRegistration(received []byte, bindings PlanRegistrationRoleBindings) (*DecodedPlanRegistration, error) {
	if len(received) > PlanRegistrationMaxCBORBytes {
		return nil, malformed("CBOR raw-byte limit exceeded")
	}
	authoritative := bytes.Clone(received)
	if hasObjectType(authoritative, ExecutionPlanObjectType) {
		if err := PredecodeExecutionPlanCBOR(authoritative); err != nil {
			return nil, err
		}
		if _, err := decodeExecutionPlan(authoritative); err != nil {
			return nil, err
		}
		return nil, domain("ExecutionPlan supplied at PlanRegistration entry point")
	}
	if err := PredecodePlanRegistrationCBOR(authoritative); err != nil {
		return nil, err
	}
	view, err := decodePlanRegistration(authoritative)
	if err != nil {
		return nil, err
	}
	if err := bindPlanRegistration(view, bindings); err != nil {
		return nil, err
	}
	return &DecodedPlanRegistration{authoritativeBytes: authoritative, view: view}, nil
}

// hasObjectType recognizes only the other closed candidate discriminator. It
// performs no allocation and grants no acceptance; the full profile still
// predecodes every non-cross-object payload before ordinary decoding.
func hasObjectType(authoritative []byte, expected string) bool {
	reader := cborReader{bytes: authoritative}
	if _, err := reader.mapLength(); err != nil {
		return false
	}
	key, err := reader.unsigned()
	if err != nil || key != 1 {
		return false
	}
	major, length, err := reader.head()
	if err != nil || major != 3 || length != uint64(len(expected)) || length > uint64(len(authoritative)-reader.offset) {
		return false
	}
	value := authoritative[reader.offset : reader.offset+int(length)]
	for index := range value {
		if value[index] != expected[index] {
			return false
		}
	}
	return true
}

func decodeExecutionPlan(authoritative []byte) (ExecutionPlan, error) {
	reader := cborReader{bytes: authoritative}
	entries, err := reader.mapLength()
	if err != nil {
		return ExecutionPlan{}, err
	}
	var view ExecutionPlan
	if err := reader.requiredKey(1); err != nil {
		return ExecutionPlan{}, err
	}
	view.ObjectType, err = reader.textString()
	if err != nil {
		return ExecutionPlan{}, err
	}
	if view.ObjectType != ExecutionPlanObjectType {
		if view.ObjectType == PlanRegistrationObjectType {
			return ExecutionPlan{}, domain("PlanRegistration supplied at ExecutionPlan entry point")
		}
		return ExecutionPlan{}, unsupported("unknown ExecutionPlan object type")
	}
	if entries < 24 {
		return ExecutionPlan{}, schema("ExecutionPlan is missing required fields")
	}
	if entries > 24 {
		return ExecutionPlan{}, unsupported("ExecutionPlan contains unknown fields")
	}

	if err := reader.requiredKey(2); err != nil {
		return ExecutionPlan{}, err
	}
	version, err := reader.unsigned()
	if err != nil {
		return ExecutionPlan{}, err
	}
	if version != uint64(CandidateObjectVersion) {
		return ExecutionPlan{}, unsupported("unsupported ExecutionPlan object version")
	}
	view.ObjectVersion = ObjectVersion(version)

	if err := reader.requiredKey(3); err != nil {
		return ExecutionPlan{}, err
	}
	if view.InstallationID, err = decodeInstallationID(&reader); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(4); err != nil {
		return ExecutionPlan{}, err
	}
	if view.EpochSequence, err = decodeUInt53(&reader); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(5); err != nil {
		return ExecutionPlan{}, err
	}
	if view.EpochDigest, err = decodeDigest(&reader, NewTrustEpochDigest); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(6); err != nil {
		return ExecutionPlan{}, err
	}
	if view.SourceManifestDigest, err = decodeDigest(&reader, NewSourceManifestDigest); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(7); err != nil {
		return ExecutionPlan{}, err
	}
	if view.SourceEntrypoint, err = decodeBoundedText(&reader, 1, 256); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(8); err != nil {
		return ExecutionPlan{}, err
	}
	if view.SourceByteLength, err = decodeUInt53(&reader); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(9); err != nil {
		return ExecutionPlan{}, err
	}
	inputSlot, err := reader.textString()
	if err != nil || inputSlot != string(PrimaryDataInputSlot) {
		return ExecutionPlan{}, schema("ExecutionPlan input slot must be primary-data")
	}
	view.InputSlot = PrimaryDataInputSlot
	if err := reader.requiredKey(10); err != nil {
		return ExecutionPlan{}, err
	}
	if view.InlineInputDigest, err = decodeDigest(&reader, NewInlineInputDigest); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(11); err != nil {
		return ExecutionPlan{}, err
	}
	if view.InlineInputByteLength, err = decodeUInt53(&reader); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(12); err != nil {
		return ExecutionPlan{}, err
	}
	if view.RuntimeProfileAlias, err = decodeBoundedText(&reader, 1, 128); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(13); err != nil {
		return ExecutionPlan{}, err
	}
	if view.RuntimeBundleManifestDigest, err = decodeDigest(&reader, NewRuntimeBundleManifestDigest); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(14); err != nil {
		return ExecutionPlan{}, err
	}
	reviewCount, err := reader.arrayLength()
	if err != nil || reviewCount < 1 || reviewCount > 8 {
		return ExecutionPlan{}, schema("ExecutionPlan requires one to eight review-attestation digests")
	}
	view.ProfileReviewAttestationDigests = make([]ProfileReviewAttestationDigest, reviewCount)
	for index := range view.ProfileReviewAttestationDigests {
		view.ProfileReviewAttestationDigests[index], err = decodeDigest(&reader, NewProfileReviewAttestationDigest)
		if err != nil {
			return ExecutionPlan{}, err
		}
	}
	if err := reader.requiredKey(15); err != nil {
		return ExecutionPlan{}, err
	}
	if view.ProfileRegistryEntryDigest, err = decodeDigest(&reader, NewProfileRegistryEntryDigest); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(16); err != nil {
		return ExecutionPlan{}, err
	}
	if view.BackendValidationRecordDigest, err = decodeDigest(&reader, NewBackendValidationRecordDigest); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(17); err != nil {
		return ExecutionPlan{}, err
	}
	if view.BackendConfigurationDigest, err = decodeDigest(&reader, NewBackendConfigurationDigest); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(18); err != nil {
		return ExecutionPlan{}, err
	}
	if view.TrustSnapshotDigest, err = decodeDigest(&reader, NewTrustSnapshotDigest); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(19); err != nil {
		return ExecutionPlan{}, err
	}
	if view.PolicyDecisionDigest, err = decodeDigest(&reader, NewPolicyDecisionDigest); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(20); err != nil {
		return ExecutionPlan{}, err
	}
	if view.WallTimeMS, err = decodePositiveUInt53(&reader); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(21); err != nil {
		return ExecutionPlan{}, err
	}
	origin, err := reader.textString()
	if err != nil || (origin != string(WallTimeRequested) && origin != string(WallTimeTrustedDefault)) {
		return ExecutionPlan{}, schema("ExecutionPlan wall-time origin is invalid")
	}
	view.WallTimeOrigin = WallTimeOrigin(origin)
	if err := reader.requiredKey(22); err != nil {
		return ExecutionPlan{}, err
	}
	outputSlot, err := reader.textString()
	if err != nil || outputSlot != string(TransformedJSONOutputSlot) {
		return ExecutionPlan{}, schema("ExecutionPlan output slot must be transformed-json")
	}
	view.OutputSlot = TransformedJSONOutputSlot
	if err := reader.requiredKey(23); err != nil {
		return ExecutionPlan{}, err
	}
	if view.OutputMaxJSONBytes, err = decodePositiveUInt53(&reader); err != nil {
		return ExecutionPlan{}, err
	}
	if err := reader.requiredKey(24); err != nil {
		return ExecutionPlan{}, err
	}
	if view.ExpiresAt, err = decodeUInt53(&reader); err != nil {
		return ExecutionPlan{}, err
	}
	return view, nil
}

func decodePlanRegistration(authoritative []byte) (PlanRegistration, error) {
	reader := cborReader{bytes: authoritative}
	entries, err := reader.mapLength()
	if err != nil {
		return PlanRegistration{}, err
	}
	var view PlanRegistration
	if err := reader.requiredKey(1); err != nil {
		return PlanRegistration{}, err
	}
	view.ObjectType, err = reader.textString()
	if err != nil {
		return PlanRegistration{}, err
	}
	if view.ObjectType != PlanRegistrationObjectType {
		if view.ObjectType == ExecutionPlanObjectType {
			return PlanRegistration{}, domain("ExecutionPlan supplied at PlanRegistration entry point")
		}
		return PlanRegistration{}, unsupported("unknown PlanRegistration object type")
	}
	if entries < 10 {
		return PlanRegistration{}, schema("PlanRegistration is missing required fields")
	}
	if entries > 10 {
		return PlanRegistration{}, unsupported("PlanRegistration contains unknown fields")
	}
	if err := reader.requiredKey(2); err != nil {
		return PlanRegistration{}, err
	}
	version, err := reader.unsigned()
	if err != nil {
		return PlanRegistration{}, err
	}
	if version != uint64(CandidateObjectVersion) {
		return PlanRegistration{}, unsupported("unsupported PlanRegistration object version")
	}
	view.ObjectVersion = ObjectVersion(version)
	if err := reader.requiredKey(3); err != nil {
		return PlanRegistration{}, err
	}
	registrationIDBytes, err := reader.byteString()
	if err != nil {
		return PlanRegistration{}, err
	}
	if view.RegistrationID, err = NewRegistrationID(registrationIDBytes); err != nil {
		return PlanRegistration{}, schema("registration ID must be a nonzero 16-byte string")
	}
	if err := reader.requiredKey(4); err != nil {
		return PlanRegistration{}, err
	}
	if view.RegistrationSequence, err = decodePositiveUInt53(&reader); err != nil {
		return PlanRegistration{}, err
	}
	if err := reader.requiredKey(5); err != nil {
		return PlanRegistration{}, err
	}
	if view.PlanDigest, err = decodeDigest(&reader, NewExecutionPlanDigest); err != nil {
		return PlanRegistration{}, err
	}
	if err := reader.requiredKey(6); err != nil {
		return PlanRegistration{}, err
	}
	if view.InstallationID, err = decodeInstallationID(&reader); err != nil {
		return PlanRegistration{}, err
	}
	if err := reader.requiredKey(7); err != nil {
		return PlanRegistration{}, err
	}
	if view.EpochSequence, err = decodeUInt53(&reader); err != nil {
		return PlanRegistration{}, err
	}
	if err := reader.requiredKey(8); err != nil {
		return PlanRegistration{}, err
	}
	if view.EpochDigest, err = decodeDigest(&reader, NewTrustEpochDigest); err != nil {
		return PlanRegistration{}, err
	}
	if err := reader.requiredKey(9); err != nil {
		return PlanRegistration{}, err
	}
	supervisorIDBytes, err := reader.byteString()
	if err != nil {
		return PlanRegistration{}, err
	}
	if view.SupervisorID, err = NewSupervisorID(supervisorIDBytes); err != nil {
		return PlanRegistration{}, schema("Supervisor ID must be a nonzero 16-byte string")
	}
	if err := reader.requiredKey(10); err != nil {
		return PlanRegistration{}, err
	}
	if view.ExpiresAt, err = decodeUInt53(&reader); err != nil {
		return PlanRegistration{}, err
	}
	return view, nil
}

func bindExecutionPlan(view ExecutionPlan, bindings ExecutionPlanRoleBindings) error {
	if view.InstallationID != bindings.InstallationID ||
		view.EpochDigest != bindings.EpochDigest ||
		view.SourceManifestDigest != bindings.SourceManifestDigest ||
		view.InlineInputDigest != bindings.InlineInputDigest ||
		view.RuntimeBundleManifestDigest != bindings.RuntimeBundleManifestDigest ||
		view.ProfileRegistryEntryDigest != bindings.ProfileRegistryEntryDigest ||
		view.BackendValidationRecordDigest != bindings.BackendValidationRecordDigest ||
		view.BackendConfigurationDigest != bindings.BackendConfigurationDigest ||
		view.TrustSnapshotDigest != bindings.TrustSnapshotDigest ||
		view.PolicyDecisionDigest != bindings.PolicyDecisionDigest ||
		!equalReviewDigests(view.ProfileReviewAttestationDigests, bindings.ProfileReviewAttestationDigests) {
		return domain("ExecutionPlan role binding mismatch")
	}
	return nil
}

func bindPlanRegistration(view PlanRegistration, bindings PlanRegistrationRoleBindings) error {
	if view.RegistrationID != bindings.RegistrationID ||
		view.PlanDigest != bindings.PlanDigest ||
		view.InstallationID != bindings.InstallationID ||
		view.EpochDigest != bindings.EpochDigest ||
		view.SupervisorID != bindings.SupervisorID {
		return domain("PlanRegistration role binding mismatch")
	}
	return nil
}

func equalReviewDigests(left, right []ProfileReviewAttestationDigest) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func decodeInstallationID(reader *cborReader) (InstallationID, error) {
	value, err := reader.byteString()
	if err != nil {
		return InstallationID{}, err
	}
	identifier, err := NewInstallationID(value)
	if err != nil {
		return InstallationID{}, schema("installation ID must be a nonzero 16-byte string")
	}
	return identifier, nil
}

func decodeUInt53(reader *cborReader) (UInt53, error) {
	value, err := reader.unsigned()
	if err != nil {
		return 0, err
	}
	return UInt53(value), nil
}

func decodePositiveUInt53(reader *cborReader) (PositiveUInt53, error) {
	value, err := reader.unsigned()
	if err != nil || value == 0 {
		return 0, schema("field must be a positive UInt53")
	}
	return PositiveUInt53(value), nil
}

func decodeBoundedText(reader *cborReader, minimum, maximum int) (string, error) {
	value, err := reader.textString()
	if err != nil {
		return "", err
	}
	if len([]byte(value)) < minimum || len([]byte(value)) > maximum {
		return "", schema("text field is outside its byte-length bounds")
	}
	return value, nil
}

func decodeDigest[T ~[32]byte](reader *cborReader, constructor func([]byte) (T, error)) (T, error) {
	value, err := reader.byteString()
	if err != nil {
		return T{}, err
	}
	digest, err := constructor(value)
	if err != nil {
		return T{}, schema("digest must be a 32-byte string")
	}
	return digest, nil
}
