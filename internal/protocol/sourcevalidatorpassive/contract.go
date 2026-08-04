// Package sourcevalidatorpassive defines the passive, unwired byte contract for
// Proposed ADR-0035 V0. It does not launch a process, invoke a parser, create an
// endpoint, enroll an artifact, mutate authority, or authorize execution.
package sourcevalidatorpassive

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

const (
	ProtocolIdentity         = "capsule.source-validator.protocol/v0"
	MethodIdentity           = "capsule.source-validator.validate-mjs-source/v0"
	RequestIdentity          = "capsule.source-validator-request/v0"
	ResultIdentity           = "capsule.source-validator-result/v0"
	ValidatorProfileIdentity = "capsule.source-validator.one-shot/v0"
	SourceProfileIdentity    = "capsule.mjs-source/v0"
	SourceMediaTypeIdentity  = "application/capsule.javascript-source;v=0;module=esm"
	CorrelationIdentity      = "capsule.source-validator.validation-request-id/v0"

	EngineeringCandidateIdentity = "capsule.source-validator.oxc-engineering-candidate/v0"
	ArtifactProfileIdentity      = "capsule.source-validator.artifact-profile/v0"

	ProtocolVersion        uint16 = 0
	RequestFrameKind       uint16 = 1
	ResultFrameKind        uint16 = 2
	MethodTag              uint16 = 1
	ValidatorProfileTag    uint16 = 1
	SourceProfileTag       uint16 = 1
	SourceMediaTypeTag     uint16 = 1
	CorrelationDomainTag   uint16 = 1
	SourceDigestDomainTag  uint16 = 0x0101
	CandidateDigestDomain  uint16 = 0x0102
	ExecutableDigestDomain uint16 = 0x0103
	BuildDigestDomain      uint16 = 0x0104
	AssessmentDigestDomain uint16 = 0x0105
	ArtifactProfileDomain  uint16 = 0x0106

	RequestHeaderBytes        = 80
	RequestMaximumBytes       = RequestHeaderBytes + v0candidate.MJSMainSourceMaxBytes
	RequestMaximumBodyBytes   = RequestMaximumBytes - 4
	ResultFrameBytes          = 138
	ArtifactProfileFrameBytes = 160
	CandidateFrameBytes       = 292
	ObservationCountMaximum   = v0candidate.MJSMainSourceMaxBytes

	candidateSemanticMode uint16 = 0x000f
	candidateDirectCrates uint16 = 6
	candidatePackages     uint16 = 65
	candidateSourceBytes  uint32 = 24_449_903
	candidateBinaryBytes  uint32 = 1_854_528
)

var (
	requestMagic   = [8]byte{'C', 'A', 'P', 'M', 'J', 'S', 'R', 'Q'}
	resultMagic    = [8]byte{'C', 'A', 'P', 'M', 'J', 'S', 'R', 'S'}
	candidateMagic = [8]byte{'C', 'A', 'P', 'M', 'J', 'S', 'C', 'I'}
	profileMagic   = [8]byte{'C', 'A', 'P', 'M', 'J', 'S', 'A', 'P'}
)

const (
	ClassificationMalformed   = "MALFORMED"
	ClassificationUnsupported = "UNSUPPORTED"
	ClassificationSchema      = "SCHEMA"
	ClassificationDomain      = "DOMAIN"
	ClassificationBinding     = "BINDING"
)

type contractError struct {
	classification string
	code           string
}

func (e *contractError) Error() string { return fmt.Sprintf("%s: %s", e.classification, e.code) }

func ErrorClassification(err error) (string, bool) {
	var classified *contractError
	if errors.As(err, &classified) {
		return classified.classification, true
	}
	return "", false
}

func reject(classification, code string) error {
	return &contractError{classification: classification, code: code}
}

type ValidationRequestID [16]byte
type SourceDigest [32]byte
type EngineeringCandidateDigest [32]byte
type ExecutableDigest [32]byte
type BuildManifestDigest [32]byte
type AssessmentDigest [32]byte
type ArtifactProfileDigest [32]byte

type ParseDisposition uint8

const (
	ParseValid                ParseDisposition = 1
	ParseDiagnosticRefusal    ParseDisposition = 2
	SemanticDiagnosticRefusal ParseDisposition = 3
)

type PolicyDisposition uint8

const (
	PolicyAllow        PolicyDisposition = 1
	PolicyDeny         PolicyDisposition = 2
	PolicyNotEvaluated PolicyDisposition = 3
)

type ResultClassification uint8

const (
	ResultNoFinding ResultClassification = iota
	ResultForbiddenSyntax
	ResultFreeCommonJS
	ResultForbiddenSyntaxAndCommonJS
	ResultParserDiagnostic
	ResultSemanticDiagnostic
)

// field-authority-object: capsule.source-validator-request v0
type ValidationRequest struct {
	ProtocolVersion  uint16
	Method           string
	ValidatorProfile string
	SourceProfile    string
	SourceMediaType  string
	CorrelationID    ValidationRequestID
	SourceDigest     SourceDigest
	SourceByteLength uint32
	sourceBytes      []byte
}

func (request *ValidationRequest) SourceBytes() []byte {
	if request == nil {
		return nil
	}
	return bytes.Clone(request.sourceBytes)
}

// field-authority-object: capsule.source-validator-result v0
type ValidationResult struct {
	ProtocolVersion            uint16
	Method                     string
	ValidatorProfile           string
	SourceProfile              string
	SourceMediaType            string
	CorrelationID              ValidationRequestID
	SourceDigest               SourceDigest
	SourceByteLength           uint32
	ArtifactProfileDigest      ArtifactProfileDigest
	ParseDisposition           ParseDisposition
	PolicyDisposition          PolicyDisposition
	Classification             ResultClassification
	StaticImportCount          uint32
	ExportFromCount            uint32
	ImportExpressionCount      uint32
	ImportMetaCount            uint32
	FreeCommonJSReferenceCount uint32
}

// field-authority-object: capsule.source-validator-engineering-candidate v0
type EngineeringCandidate struct {
	RecordVersion        uint16
	Implementation       string
	OxcVersionMajor      uint16
	OxcVersionMinor      uint16
	OxcVersionPatch      uint16
	RustToolchainMajor   uint16
	RustToolchainMinor   uint16
	RustToolchainPatch   uint16
	SemanticMode         uint16
	DirectCrateCount     uint16
	LockedPackageCount   uint16
	CachedSourceBytes    uint32
	ReleaseProbeBytes    uint32
	CargoLockDigest      [32]byte
	OxcAllocatorChecksum [32]byte
	OxcASTChecksum       [32]byte
	OxcASTVisitChecksum  [32]byte
	OxcParserChecksum    [32]byte
	OxcSemanticChecksum  [32]byte
	OxcSpanChecksum      [32]byte
}

// field-authority-object: capsule.source-validator-artifact-profile v0
type EnrolledArtifactProfile struct {
	RecordVersion              uint16
	ProtocolVersion            uint16
	Method                     string
	ValidatorProfile           string
	SourceProfile              string
	SourceMediaType            string
	EngineeringCandidateDigest EngineeringCandidateDigest
	ExecutableDigest           ExecutableDigest
	BuildManifestDigest        BuildManifestDigest
	AssessmentDigest           AssessmentDigest
}

func NewValidationRequest(id ValidationRequestID, source []byte) (*ValidationRequest, error) {
	if id == (ValidationRequestID{}) {
		return nil, reject(ClassificationSchema, "zero-correlation-id")
	}
	validated, err := v0candidate.NewMJSMainSource(source)
	if err != nil {
		return nil, reject(ClassificationDomain, "source-profile")
	}
	exact := validated.AuthoritativeBytes()
	digest := sha256.Sum256(exact)
	// M1 validation above proves len(exact) <= 262,144 before this narrowing.
	sourceByteLength := uint32(len(exact)) // #nosec G115 -- bounded by the canonical M1 cap.
	return &ValidationRequest{
		ProtocolVersion:  ProtocolVersion,
		Method:           MethodIdentity,
		ValidatorProfile: ValidatorProfileIdentity,
		SourceProfile:    SourceProfileIdentity,
		SourceMediaType:  SourceMediaTypeIdentity,
		CorrelationID:    id,
		SourceDigest:     SourceDigest(digest),
		SourceByteLength: sourceByteLength,
		sourceBytes:      exact,
	}, nil
}

func EncodeRequest(request *ValidationRequest) ([]byte, error) {
	if err := validateRequestValue(request); err != nil {
		return nil, err
	}
	frame := make([]byte, RequestHeaderBytes+len(request.sourceBytes))
	putBodyLength(frame)
	copy(frame[4:12], requestMagic[:])
	putCommon(frame, RequestFrameKind, request.CorrelationID, request.SourceByteLength, request.SourceDigest)
	copy(frame[RequestHeaderBytes:], request.sourceBytes)
	return frame, nil
}

func DecodeRequest(received []byte) (*ValidationRequest, error) {
	if err := validateFrame(received, requestMagic, RequestHeaderBytes, RequestMaximumBytes); err != nil {
		return nil, err
	}
	if err := validateCommon(received, RequestFrameKind); err != nil {
		return nil, err
	}
	declaredLength := binary.BigEndian.Uint32(received[44:48])
	if declaredLength > v0candidate.MJSMainSourceMaxBytes {
		return nil, reject(ClassificationDomain, "source-cap")
	}
	if int(declaredLength) != len(received)-RequestHeaderBytes {
		return nil, reject(ClassificationBinding, "source-length")
	}
	var id ValidationRequestID
	copy(id[:], received[28:44])
	if id == (ValidationRequestID{}) {
		return nil, reject(ClassificationSchema, "zero-correlation-id")
	}
	source := bytes.Clone(received[RequestHeaderBytes:])
	validated, err := v0candidate.NewMJSMainSource(source)
	if err != nil {
		return nil, reject(ClassificationDomain, "source-profile")
	}
	exact := validated.AuthoritativeBytes()
	digest := sha256.Sum256(exact)
	if !bytes.Equal(digest[:], received[48:80]) {
		return nil, reject(ClassificationBinding, "source-digest")
	}
	return &ValidationRequest{
		ProtocolVersion:  ProtocolVersion,
		Method:           MethodIdentity,
		ValidatorProfile: ValidatorProfileIdentity,
		SourceProfile:    SourceProfileIdentity,
		SourceMediaType:  SourceMediaTypeIdentity,
		CorrelationID:    id,
		SourceDigest:     SourceDigest(digest),
		SourceByteLength: declaredLength,
		sourceBytes:      exact,
	}, nil
}

func NewValidationResult(
	request *ValidationRequest,
	profileDigest ArtifactProfileDigest,
	parse ParseDisposition,
	policy PolicyDisposition,
	classification ResultClassification,
	counts [5]uint32,
) (*ValidationResult, error) {
	if err := validateRequestValue(request); err != nil {
		return nil, err
	}
	result := &ValidationResult{
		ProtocolVersion: ProtocolVersion, Method: MethodIdentity,
		ValidatorProfile: ValidatorProfileIdentity, SourceProfile: SourceProfileIdentity,
		SourceMediaType: SourceMediaTypeIdentity, CorrelationID: request.CorrelationID,
		SourceDigest: request.SourceDigest, SourceByteLength: request.SourceByteLength,
		ArtifactProfileDigest: profileDigest, ParseDisposition: parse,
		PolicyDisposition: policy, Classification: classification,
		StaticImportCount: counts[0], ExportFromCount: counts[1],
		ImportExpressionCount: counts[2], ImportMetaCount: counts[3],
		FreeCommonJSReferenceCount: counts[4],
	}
	if err := validateResultValue(result); err != nil {
		return nil, err
	}
	return result, nil
}

func EncodeResult(result *ValidationResult) ([]byte, error) {
	if err := validateResultValue(result); err != nil {
		return nil, err
	}
	frame := make([]byte, ResultFrameBytes)
	putBodyLength(frame)
	copy(frame[4:12], resultMagic[:])
	putCommon(frame, ResultFrameKind, result.CorrelationID, result.SourceByteLength, result.SourceDigest)
	binary.BigEndian.PutUint16(frame[80:82], ArtifactProfileDomain)
	copy(frame[82:114], result.ArtifactProfileDigest[:])
	frame[114] = byte(result.ParseDisposition)
	frame[115] = byte(result.PolicyDisposition)
	frame[116] = byte(result.Classification)
	counts := [...]uint32{result.StaticImportCount, result.ExportFromCount, result.ImportExpressionCount, result.ImportMetaCount, result.FreeCommonJSReferenceCount}
	for index, count := range counts {
		binary.BigEndian.PutUint32(frame[118+index*4:122+index*4], count)
	}
	return frame, nil
}

func DecodeResult(received []byte) (*ValidationResult, error) {
	if err := validateFrame(received, resultMagic, ResultFrameBytes, ResultFrameBytes); err != nil {
		return nil, err
	}
	if err := validateCommon(received, ResultFrameKind); err != nil {
		return nil, err
	}
	if binary.BigEndian.Uint16(received[80:82]) != ArtifactProfileDomain {
		return nil, reject(ClassificationDomain, "artifact-profile-digest-domain")
	}
	if received[117] != 0 {
		return nil, reject(ClassificationUnsupported, "result-reserved-byte")
	}
	var id ValidationRequestID
	copy(id[:], received[28:44])
	if id == (ValidationRequestID{}) {
		return nil, reject(ClassificationSchema, "zero-correlation-id")
	}
	var sourceDigest SourceDigest
	copy(sourceDigest[:], received[48:80])
	var profileDigest ArtifactProfileDigest
	copy(profileDigest[:], received[82:114])
	result := &ValidationResult{
		ProtocolVersion: ProtocolVersion, Method: MethodIdentity,
		ValidatorProfile: ValidatorProfileIdentity, SourceProfile: SourceProfileIdentity,
		SourceMediaType: SourceMediaTypeIdentity, CorrelationID: id,
		SourceDigest: sourceDigest, SourceByteLength: binary.BigEndian.Uint32(received[44:48]),
		ArtifactProfileDigest: profileDigest,
		ParseDisposition:      ParseDisposition(received[114]), PolicyDisposition: PolicyDisposition(received[115]),
		Classification:             ResultClassification(received[116]),
		StaticImportCount:          binary.BigEndian.Uint32(received[118:122]),
		ExportFromCount:            binary.BigEndian.Uint32(received[122:126]),
		ImportExpressionCount:      binary.BigEndian.Uint32(received[126:130]),
		ImportMetaCount:            binary.BigEndian.Uint32(received[130:134]),
		FreeCommonJSReferenceCount: binary.BigEndian.Uint32(received[134:138]),
	}
	if err := validateResultValue(result); err != nil {
		return nil, err
	}
	return result, nil
}

func VerifyResult(request *ValidationRequest, artifactProfileBytes, resultBytes []byte) (*ValidationResult, error) {
	if err := validateRequestValue(request); err != nil {
		return nil, err
	}
	profile, err := DecodeArtifactProfile(artifactProfileBytes)
	if err != nil {
		return nil, err
	}
	result, err := DecodeResult(resultBytes)
	if err != nil {
		return nil, err
	}
	profileDigest := ArtifactProfileIdentityDigest(artifactProfileBytes)
	if result.ArtifactProfileDigest != profileDigest || profile.ProtocolVersion != ProtocolVersion ||
		profile.Method != MethodIdentity || profile.ValidatorProfile != ValidatorProfileIdentity ||
		profile.SourceProfile != SourceProfileIdentity || profile.SourceMediaType != SourceMediaTypeIdentity {
		return nil, reject(ClassificationBinding, "artifact-profile")
	}
	if result.CorrelationID != request.CorrelationID || result.SourceDigest != request.SourceDigest ||
		result.SourceByteLength != request.SourceByteLength {
		return nil, reject(ClassificationBinding, "request-result")
	}
	return result, nil
}

func DefaultEngineeringCandidate() EngineeringCandidate {
	return EngineeringCandidate{
		RecordVersion: 0, Implementation: EngineeringCandidateIdentity,
		OxcVersionMajor: 0, OxcVersionMinor: 140, OxcVersionPatch: 0,
		RustToolchainMajor: 1, RustToolchainMinor: 95, RustToolchainPatch: 0,
		SemanticMode: candidateSemanticMode, DirectCrateCount: candidateDirectCrates,
		LockedPackageCount: candidatePackages, CachedSourceBytes: candidateSourceBytes,
		ReleaseProbeBytes:    candidateBinaryBytes,
		CargoLockDigest:      mustHex32("505669a07338603876bc96c242f8d5af386d3a13139e70110a8b52f39bae69ac"),
		OxcAllocatorChecksum: mustHex32("0f8245ba555b465d3577732d5f9d9306babb0aaa7b80e97a2ce21f74fae442a3"),
		OxcASTChecksum:       mustHex32("3305400b90fff2a30b272b58fe6080d25369407b2ac37c4ac652996a9677efe0"),
		OxcASTVisitChecksum:  mustHex32("4640e6d0de2e0f6c820d1444a468d070c710111df76ce90a1694ac386641e133"),
		OxcParserChecksum:    mustHex32("8abd68f81349d37ea79f1d99d2370e15f282cc9fbe66e8544d072595744ab38e"),
		OxcSemanticChecksum:  mustHex32("5967f96881e1694d10b453311fa681b4df0f38760628e1de613b046566cd8c8e"),
		OxcSpanChecksum:      mustHex32("e83fa0a0fe6e5e2f5abb173a64afe8db711bb612acbc002c663fd13b08a8cbf3"),
	}
}

func EncodeEngineeringCandidate(candidate EngineeringCandidate) ([]byte, error) {
	if candidate != DefaultEngineeringCandidate() {
		return nil, reject(ClassificationUnsupported, "engineering-candidate-identity")
	}
	frame := make([]byte, CandidateFrameBytes)
	putBodyLength(frame)
	copy(frame[4:12], candidateMagic[:])
	values := []uint16{candidate.RecordVersion, 1, candidate.OxcVersionMajor, candidate.OxcVersionMinor, candidate.OxcVersionPatch, candidate.RustToolchainMajor, candidate.RustToolchainMinor, candidate.RustToolchainPatch, candidate.SemanticMode, candidate.DirectCrateCount, candidate.LockedPackageCount, 0}
	for index, value := range values {
		binary.BigEndian.PutUint16(frame[12+index*2:14+index*2], value)
	}
	binary.BigEndian.PutUint32(frame[36:40], candidate.CachedSourceBytes)
	binary.BigEndian.PutUint32(frame[40:44], candidate.ReleaseProbeBytes)
	copy(frame[44:76], candidate.CargoLockDigest[:])
	checksums := [...][32]byte{candidate.OxcAllocatorChecksum, candidate.OxcASTChecksum, candidate.OxcASTVisitChecksum, candidate.OxcParserChecksum, candidate.OxcSemanticChecksum, candidate.OxcSpanChecksum}
	for index, checksum := range checksums {
		offset := 76 + index*36
		binary.BigEndian.PutUint16(frame[offset:offset+2], uint16(index+1)) // #nosec G115 -- six fixed entries.
		copy(frame[offset+4:offset+36], checksum[:])
	}
	return frame, nil
}

func DecodeEngineeringCandidate(received []byte) (EngineeringCandidate, error) {
	if err := validateFrame(received, candidateMagic, CandidateFrameBytes, CandidateFrameBytes); err != nil {
		return EngineeringCandidate{}, err
	}
	want, _ := EncodeEngineeringCandidate(DefaultEngineeringCandidate())
	if !bytes.Equal(received, want) {
		return EngineeringCandidate{}, reject(ClassificationUnsupported, "engineering-candidate-identity")
	}
	return DefaultEngineeringCandidate(), nil
}

func EngineeringCandidateIdentityDigest(candidateBytes []byte) EngineeringCandidateDigest {
	return EngineeringCandidateDigest(domainHash("capsule.source-validator.engineering-candidate/v0", candidateBytes))
}

func NewEnrolledArtifactProfile(executable ExecutableDigest, build BuildManifestDigest, assessment AssessmentDigest) (EnrolledArtifactProfile, error) {
	if executable == (ExecutableDigest{}) || build == (BuildManifestDigest{}) || assessment == (AssessmentDigest{}) {
		return EnrolledArtifactProfile{}, reject(ClassificationSchema, "zero-artifact-profile-digest")
	}
	candidateBytes, _ := EncodeEngineeringCandidate(DefaultEngineeringCandidate())
	return EnrolledArtifactProfile{
		RecordVersion: 0, ProtocolVersion: ProtocolVersion, Method: MethodIdentity,
		ValidatorProfile: ValidatorProfileIdentity, SourceProfile: SourceProfileIdentity,
		SourceMediaType:            SourceMediaTypeIdentity,
		EngineeringCandidateDigest: EngineeringCandidateIdentityDigest(candidateBytes),
		ExecutableDigest:           executable, BuildManifestDigest: build, AssessmentDigest: assessment,
	}, nil
}

func EncodeArtifactProfile(profile EnrolledArtifactProfile) ([]byte, error) {
	if err := validateArtifactProfile(profile); err != nil {
		return nil, err
	}
	frame := make([]byte, ArtifactProfileFrameBytes)
	putBodyLength(frame)
	copy(frame[4:12], profileMagic[:])
	for index, value := range []uint16{profile.RecordVersion, ProtocolVersion, MethodTag, ValidatorProfileTag, SourceProfileTag, SourceMediaTypeTag} {
		binary.BigEndian.PutUint16(frame[12+index*2:14+index*2], value)
	}
	binary.BigEndian.PutUint16(frame[24:26], CandidateDigestDomain)
	copy(frame[26:58], profile.EngineeringCandidateDigest[:])
	binary.BigEndian.PutUint16(frame[58:60], ExecutableDigestDomain)
	copy(frame[60:92], profile.ExecutableDigest[:])
	binary.BigEndian.PutUint16(frame[92:94], BuildDigestDomain)
	copy(frame[94:126], profile.BuildManifestDigest[:])
	binary.BigEndian.PutUint16(frame[126:128], AssessmentDigestDomain)
	copy(frame[128:160], profile.AssessmentDigest[:])
	return frame, nil
}

func DecodeArtifactProfile(received []byte) (EnrolledArtifactProfile, error) {
	if err := validateFrame(received, profileMagic, ArtifactProfileFrameBytes, ArtifactProfileFrameBytes); err != nil {
		return EnrolledArtifactProfile{}, err
	}
	if binary.BigEndian.Uint16(received[12:14]) != 0 || binary.BigEndian.Uint16(received[14:16]) != ProtocolVersion ||
		binary.BigEndian.Uint16(received[16:18]) != MethodTag || binary.BigEndian.Uint16(received[18:20]) != ValidatorProfileTag ||
		binary.BigEndian.Uint16(received[20:22]) != SourceProfileTag || binary.BigEndian.Uint16(received[22:24]) != SourceMediaTypeTag {
		return EnrolledArtifactProfile{}, reject(ClassificationUnsupported, "artifact-profile-version-or-profile")
	}
	if binary.BigEndian.Uint16(received[24:26]) != CandidateDigestDomain || binary.BigEndian.Uint16(received[58:60]) != ExecutableDigestDomain ||
		binary.BigEndian.Uint16(received[92:94]) != BuildDigestDomain || binary.BigEndian.Uint16(received[126:128]) != AssessmentDigestDomain {
		return EnrolledArtifactProfile{}, reject(ClassificationDomain, "artifact-profile-digest-domain")
	}
	profile := EnrolledArtifactProfile{
		RecordVersion: 0, ProtocolVersion: ProtocolVersion, Method: MethodIdentity,
		ValidatorProfile: ValidatorProfileIdentity, SourceProfile: SourceProfileIdentity,
		SourceMediaType: SourceMediaTypeIdentity,
	}
	copy(profile.EngineeringCandidateDigest[:], received[26:58])
	copy(profile.ExecutableDigest[:], received[60:92])
	copy(profile.BuildManifestDigest[:], received[94:126])
	copy(profile.AssessmentDigest[:], received[128:160])
	if err := validateArtifactProfile(profile); err != nil {
		return EnrolledArtifactProfile{}, err
	}
	return profile, nil
}

func ArtifactProfileIdentityDigest(profileBytes []byte) ArtifactProfileDigest {
	return ArtifactProfileDigest(domainHash("capsule.source-validator.artifact-profile/v0", profileBytes))
}

func validateRequestValue(request *ValidationRequest) error {
	if request == nil || request.ProtocolVersion != ProtocolVersion || request.Method != MethodIdentity ||
		request.ValidatorProfile != ValidatorProfileIdentity || request.SourceProfile != SourceProfileIdentity ||
		request.SourceMediaType != SourceMediaTypeIdentity {
		return reject(ClassificationUnsupported, "request-profile")
	}
	if request.CorrelationID == (ValidationRequestID{}) {
		return reject(ClassificationSchema, "zero-correlation-id")
	}
	if len(request.sourceBytes) > v0candidate.MJSMainSourceMaxBytes {
		return reject(ClassificationBinding, "source-length")
	}
	// The preceding inclusive M1 cap proves this narrowing is exact.
	sourceByteLength := uint32(len(request.sourceBytes)) // #nosec G115 -- bounded by 262,144.
	if sourceByteLength != request.SourceByteLength {
		return reject(ClassificationBinding, "source-length")
	}
	digest := sha256.Sum256(request.sourceBytes)
	if SourceDigest(digest) != request.SourceDigest {
		return reject(ClassificationBinding, "source-digest")
	}
	if _, err := v0candidate.NewMJSMainSource(request.sourceBytes); err != nil {
		return reject(ClassificationDomain, "source-profile")
	}
	return nil
}

func validateResultValue(result *ValidationResult) error {
	if result == nil || result.ProtocolVersion != ProtocolVersion || result.Method != MethodIdentity ||
		result.ValidatorProfile != ValidatorProfileIdentity || result.SourceProfile != SourceProfileIdentity ||
		result.SourceMediaType != SourceMediaTypeIdentity {
		return reject(ClassificationUnsupported, "result-profile")
	}
	if result.CorrelationID == (ValidationRequestID{}) || result.SourceByteLength > v0candidate.MJSMainSourceMaxBytes ||
		result.ArtifactProfileDigest == (ArtifactProfileDigest{}) {
		return reject(ClassificationSchema, "result-identity-or-length")
	}
	counts := [...]uint32{result.StaticImportCount, result.ExportFromCount, result.ImportExpressionCount, result.ImportMetaCount, result.FreeCommonJSReferenceCount}
	for _, count := range counts {
		if count > ObservationCountMaximum || count > result.SourceByteLength {
			return reject(ClassificationDomain, "result-count")
		}
	}
	syntax := counts[0]+counts[1]+counts[2]+counts[3] > 0
	commonjs := counts[4] > 0
	switch result.ParseDisposition {
	case ParseValid:
		if result.PolicyDisposition == PolicyAllow && !syntax && !commonjs && result.Classification == ResultNoFinding {
			return nil
		}
		if result.PolicyDisposition != PolicyDeny {
			return reject(ClassificationDomain, "policy-disposition")
		}
		if syntax && !commonjs && result.Classification == ResultForbiddenSyntax {
			return nil
		}
		if !syntax && commonjs && result.Classification == ResultFreeCommonJS {
			return nil
		}
		if syntax && commonjs && result.Classification == ResultForbiddenSyntaxAndCommonJS {
			return nil
		}
	case ParseDiagnosticRefusal:
		if result.PolicyDisposition == PolicyNotEvaluated && result.Classification == ResultParserDiagnostic && !syntax && !commonjs {
			return nil
		}
	case SemanticDiagnosticRefusal:
		if result.PolicyDisposition == PolicyNotEvaluated && result.Classification == ResultSemanticDiagnostic && !syntax && !commonjs {
			return nil
		}
	}
	return reject(ClassificationDomain, "result-status-classification")
}

func validateArtifactProfile(profile EnrolledArtifactProfile) error {
	if profile.RecordVersion != 0 || profile.ProtocolVersion != ProtocolVersion || profile.Method != MethodIdentity ||
		profile.ValidatorProfile != ValidatorProfileIdentity || profile.SourceProfile != SourceProfileIdentity ||
		profile.SourceMediaType != SourceMediaTypeIdentity {
		return reject(ClassificationUnsupported, "artifact-profile")
	}
	candidateBytes, _ := EncodeEngineeringCandidate(DefaultEngineeringCandidate())
	if profile.EngineeringCandidateDigest != EngineeringCandidateIdentityDigest(candidateBytes) {
		return reject(ClassificationBinding, "engineering-candidate")
	}
	if profile.ExecutableDigest == (ExecutableDigest{}) || profile.BuildManifestDigest == (BuildManifestDigest{}) || profile.AssessmentDigest == (AssessmentDigest{}) {
		return reject(ClassificationSchema, "zero-artifact-profile-digest")
	}
	return nil
}

func putCommon(frame []byte, kind uint16, id ValidationRequestID, length uint32, digest SourceDigest) {
	for index, value := range []uint16{ProtocolVersion, kind, MethodTag, ValidatorProfileTag, SourceProfileTag, SourceMediaTypeTag, CorrelationDomainTag, SourceDigestDomainTag} {
		binary.BigEndian.PutUint16(frame[12+index*2:14+index*2], value)
	}
	copy(frame[28:44], id[:])
	binary.BigEndian.PutUint32(frame[44:48], length)
	copy(frame[48:80], digest[:])
}

func validateCommon(frame []byte, kind uint16) error {
	values := []uint16{ProtocolVersion, kind, MethodTag, ValidatorProfileTag, SourceProfileTag, SourceMediaTypeTag, CorrelationDomainTag, SourceDigestDomainTag}
	for index, want := range values {
		got := binary.BigEndian.Uint16(frame[12+index*2 : 14+index*2])
		if got == want {
			continue
		}
		classification := ClassificationUnsupported
		if index == 6 || index == 7 {
			classification = ClassificationDomain
		}
		return reject(classification, fmt.Sprintf("common-field-%d", index))
	}
	return nil
}

func validateFrame(frame []byte, magic [8]byte, minimum, maximum int) error {
	if len(frame) < 4 {
		return reject(ClassificationMalformed, "truncated-length")
	}
	declared := uint64(binary.BigEndian.Uint32(frame[:4])) + 4
	// All callers supply fixed nonnegative contract constants no larger than 262,224.
	maximum64 := uint64(maximum) // #nosec G115 -- fixed, reviewed frame maximum.
	minimum64 := uint64(minimum) // #nosec G115 -- fixed, reviewed frame minimum.
	if declared > maximum64 {
		return reject(ClassificationDomain, "frame-cap")
	}
	// A frame that reached this point is at most maximum bytes, so the host length is safe to widen.
	frameLength := uint64(len(frame)) // #nosec G115 -- len(frame) is nonnegative and frame-capped.
	if declared < minimum64 || declared != frameLength {
		return reject(ClassificationMalformed, "frame-length-or-trailing")
	}
	if len(frame) < 12 || !bytes.Equal(frame[4:12], magic[:]) {
		return reject(ClassificationDomain, "frame-domain")
	}
	return nil
}

func putBodyLength(frame []byte) {
	// Encoders allocate only the fixed contract frames, whose largest body is 262,220 bytes.
	bodyLength := uint32(len(frame) - 4) // #nosec G115 -- construction is contract-capped.
	binary.BigEndian.PutUint32(frame[:4], bodyLength)
}

func domainHash(domain string, value []byte) [32]byte {
	hash := sha256.New()
	hash.Write([]byte(domain))
	hash.Write([]byte{0})
	hash.Write(value)
	var result [32]byte
	copy(result[:], hash.Sum(nil))
	return result
}

func mustHex32(value string) [32]byte {
	var result [32]byte
	for index := range result {
		_, err := fmt.Sscanf(value[index*2:index*2+2], "%02x", &result[index])
		if err != nil {
			panic(err)
		}
	}
	return result
}
