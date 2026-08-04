package sourcevalidatorpassive

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

// R1 is a passive, unwired contract. These identities and codecs do not create
// an XPC endpoint, launch a parser, enroll an artifact, or authorize planning or
// approval. Active resource measurements remain deliberately absent until R4.
const (
	V1ProtocolIdentity               = "capsule.source-validator.protocol/v1"
	V1DaemonMethodIdentity           = "capsule.source-validator.validate-mjs-source.daemon/v1"
	V1BrokerMethodIdentity           = "capsule.source-validator.validate-mjs-source.approval-broker/v1"
	V1DaemonRequestIdentity          = "capsule.source-validator-request.daemon/v1"
	V1BrokerRequestIdentity          = "capsule.source-validator-request.approval-broker/v1"
	V1DaemonResultIdentity           = "capsule.source-validator-result.daemon/v1"
	V1BrokerResultIdentity           = "capsule.source-validator-result.approval-broker/v1"
	V1DaemonProcessProfileIdentity   = "capsule.source-validator.macos-xpc-parser-child.daemon/v1"
	V1BrokerProcessProfileIdentity   = "capsule.source-validator.macos-xpc-parser-child.approval-broker/v1"
	V1DaemonArtifactProfileIdentity  = "capsule.source-validator.artifact-profile.daemon/v1"
	V1BrokerArtifactProfileIdentity  = "capsule.source-validator.artifact-profile.approval-broker/v1"
	V1ResourcePolicyIdentity         = "capsule.source-validator.reactive-footprint-policy/v1"
	V1DaemonServiceIdentity          = "com.capsulecorp.capsule.source-validator.daemon.v1"
	V1BrokerServiceIdentity          = "com.capsulecorp.capsule.source-validator.approval-broker.v1"
	V1DaemonParserIdentity           = "com.capsulecorp.capsule.source-validator-parser.daemon.v1"
	V1BrokerParserIdentity           = "com.capsulecorp.capsule.source-validator-parser.approval-broker.v1"
	V1DaemonContainingBundleIdentity = "com.capsulecorp.capsule.daemon"
	V1BrokerContainingBundleIdentity = "com.capsulecorp.capsule.broker"

	V1ProtocolVersion              uint16 = 1
	V1RequestFrameKind             uint16 = 1
	V1ResultFrameKind              uint16 = 2
	V1ResourcePolicyFrameKind      uint16 = 3
	V1ProcessProfileFrameKind      uint16 = 4
	V1ArtifactProfileFrameKind     uint16 = 5
	V1ConsumerProjectionFrameKind  uint16 = 6
	V1RequestHeaderBytes                  = 216
	V1RequestMaximumBytes                 = V1RequestHeaderBytes + v0candidate.MJSMainSourceMaxBytes
	V1ResultFrameBytes                    = 248
	V1ResourcePolicyFrameBytes            = 256
	V1ProcessProfileFrameBytes            = 256
	V1ArtifactProfileFrameBytes           = 256
	V1ConsumerProjectionFrameBytes        = 192
)

var (
	v1RequestMagic           = [8]byte{'C', 'S', 'V', '1', 'R', 'E', 'Q', '0'}
	v1ResultMagic            = [8]byte{'C', 'S', 'V', '1', 'R', 'E', 'S', '0'}
	v1PolicyMagic            = [8]byte{'C', 'S', 'V', '1', 'P', 'O', 'L', '0'}
	v1ProcessMagic           = [8]byte{'C', 'S', 'V', '1', 'P', 'R', 'C', '0'}
	v1ArtifactMagic          = [8]byte{'C', 'S', 'V', '1', 'A', 'R', 'T', '0'}
	v1ConsumerMagic          = [8]byte{'C', 'S', 'V', '1', 'C', 'O', 'N', '0'}
	v0RequestMagicForRefusal = [8]byte{'C', 'A', 'P', 'M', 'J', 'S', 'R', 'Q'}
	v0ResultMagicForRefusal  = [8]byte{'C', 'A', 'P', 'M', 'J', 'S', 'R', 'S'}
)

type V1ConsumerRole uint16

const (
	V1DaemonRole         V1ConsumerRole = 1
	V1ApprovalBrokerRole V1ConsumerRole = 2
)

func (role V1ConsumerRole) String() string {
	switch role {
	case V1DaemonRole:
		return "daemon"
	case V1ApprovalBrokerRole:
		return "approval-broker"
	default:
		return "unknown"
	}
}

func (role V1ConsumerRole) Other() V1ConsumerRole {
	if role == V1DaemonRole {
		return V1ApprovalBrokerRole
	}
	return V1DaemonRole
}

type V1CorrelationID [16]byte
type V1InstallationID [16]byte
type V1Digest [32]byte

type V1RequestBindings struct {
	CorrelationID         V1CorrelationID
	InstallationID        V1InstallationID
	EpochSequence         uint64
	EpochDigest           V1Digest
	ArtifactProfileDigest V1Digest
	ResourcePolicyDigest  V1Digest
}

// field-authority-object: capsule.source-validator-request v1
type V1ValidationRequest struct {
	ProtocolVersion       uint16
	Role                  V1ConsumerRole
	Method                string
	RequestIdentity       string
	SourceProfile         string
	SourceMediaType       string
	CorrelationID         V1CorrelationID
	InstallationID        V1InstallationID
	EpochSequence         uint64
	EpochDigest           V1Digest
	SourceByteLength      uint32
	SourceDigest          V1Digest
	ArtifactProfileDigest V1Digest
	ResourcePolicyDigest  V1Digest
	sourceBytes           []byte
}

func (request *V1ValidationRequest) SourceBytes() []byte {
	if request == nil {
		return nil
	}
	return bytes.Clone(request.sourceBytes)
}

func NewV1ValidationRequest(role V1ConsumerRole, bindings V1RequestBindings, source []byte) (*V1ValidationRequest, error) {
	if err := validateV1Role(role); err != nil {
		return nil, err
	}
	if err := validateV1Bindings(bindings); err != nil {
		return nil, err
	}
	validated, err := v0candidate.NewMJSMainSource(source)
	if err != nil {
		return nil, reject(ClassificationDomain, "v1-source-profile")
	}
	exact := validated.AuthoritativeBytes()
	digest := sha256.Sum256(exact)
	method, requestIdentity := v1RequestIdentities(role)
	return &V1ValidationRequest{
		ProtocolVersion:       V1ProtocolVersion,
		Role:                  role,
		Method:                method,
		RequestIdentity:       requestIdentity,
		SourceProfile:         SourceProfileIdentity,
		SourceMediaType:       SourceMediaTypeIdentity,
		CorrelationID:         bindings.CorrelationID,
		InstallationID:        bindings.InstallationID,
		EpochSequence:         bindings.EpochSequence,
		EpochDigest:           bindings.EpochDigest,
		SourceByteLength:      uint32(len(exact)), // #nosec G115 -- M1 caps exact bytes at 262,144.
		SourceDigest:          V1Digest(digest),
		ArtifactProfileDigest: bindings.ArtifactProfileDigest,
		ResourcePolicyDigest:  bindings.ResourcePolicyDigest,
		sourceBytes:           exact,
	}, nil
}

func EncodeV1Request(request *V1ValidationRequest) ([]byte, error) {
	if err := validateV1Request(request); err != nil {
		return nil, err
	}
	frame := make([]byte, V1RequestHeaderBytes+len(request.sourceBytes))
	v1PutLengthMagic(frame, v1RequestMagic)
	v1PutRequestCommon(frame, request)
	copy(frame[V1RequestHeaderBytes:], request.sourceBytes)
	return frame, nil
}

func DecodeV1Request(expectedRole V1ConsumerRole, frame []byte) (*V1ValidationRequest, error) {
	if err := v1ValidateFrame(frame, v1RequestMagic, V1RequestHeaderBytes, V1RequestMaximumBytes, V1RequestFrameKind); err != nil {
		return nil, err
	}
	role := V1ConsumerRole(binary.BigEndian.Uint16(frame[16:18]))
	if err := v1RequireRole(expectedRole, role); err != nil {
		return nil, err
	}
	if err := v1ValidateRequestTags(frame, role); err != nil {
		return nil, err
	}
	request := &V1ValidationRequest{ProtocolVersion: V1ProtocolVersion, Role: role}
	request.Method, request.RequestIdentity = v1RequestIdentities(role)
	request.SourceProfile, request.SourceMediaType = SourceProfileIdentity, SourceMediaTypeIdentity
	copy(request.CorrelationID[:], frame[44:60])
	copy(request.InstallationID[:], frame[60:76])
	request.EpochSequence = binary.BigEndian.Uint64(frame[76:84])
	copy(request.EpochDigest[:], frame[84:116])
	request.SourceByteLength = binary.BigEndian.Uint32(frame[116:120])
	copy(request.SourceDigest[:], frame[120:152])
	copy(request.ArtifactProfileDigest[:], frame[152:184])
	copy(request.ResourcePolicyDigest[:], frame[184:216])
	if int(request.SourceByteLength) != len(frame)-V1RequestHeaderBytes {
		return nil, reject(ClassificationBinding, "v1-source-length")
	}
	request.sourceBytes = bytes.Clone(frame[V1RequestHeaderBytes:])
	if err := validateV1Request(request); err != nil {
		return nil, err
	}
	return request, nil
}

// VerifyV1Request repeats decoding and compares every caller-independent binding
// with the launcher's separately trusted current installation/profile context.
func VerifyV1Request(expectedRole V1ConsumerRole, expected V1RequestBindings, frame []byte) (*V1ValidationRequest, error) {
	if err := validateV1Bindings(expected); err != nil {
		return nil, err
	}
	request, err := DecodeV1Request(expectedRole, frame)
	if err != nil {
		return nil, err
	}
	if request.CorrelationID != expected.CorrelationID || request.InstallationID != expected.InstallationID ||
		request.EpochSequence != expected.EpochSequence || request.EpochDigest != expected.EpochDigest ||
		request.ArtifactProfileDigest != expected.ArtifactProfileDigest || request.ResourcePolicyDigest != expected.ResourcePolicyDigest {
		return nil, reject(ClassificationBinding, "v1-request-current-context")
	}
	return request, nil
}

// field-authority-object: capsule.source-validator-result v1
type V1ValidationResult struct {
	ProtocolVersion            uint16
	Role                       V1ConsumerRole
	Method                     string
	ResultIdentity             string
	CorrelationID              V1CorrelationID
	InstallationID             V1InstallationID
	EpochSequence              uint64
	EpochDigest                V1Digest
	SourceByteLength           uint32
	SourceDigest               V1Digest
	ArtifactProfileDigest      V1Digest
	ResourcePolicyDigest       V1Digest
	ParseDisposition           ParseDisposition
	PolicyDisposition          PolicyDisposition
	Classification             ResultClassification
	CleanupDisposition         V1CleanupDisposition
	RefusalDisposition         V1RefusalDisposition
	StaticImportCount          uint32
	ExportFromCount            uint32
	ImportExpressionCount      uint32
	ImportMetaCount            uint32
	FreeCommonJSReferenceCount uint32
}

type V1CleanupDisposition uint8
type V1RefusalDisposition uint16

const (
	V1CleanupComplete V1CleanupDisposition = 1
	V1CleanupFailed   V1CleanupDisposition = 2

	V1RefusalNone             V1RefusalDisposition = 0
	V1RefusalCancelled        V1RefusalDisposition = 1
	V1RefusalDisconnected     V1RefusalDisposition = 2
	V1RefusalDeadline         V1RefusalDisposition = 3
	V1RefusalWatermark        V1RefusalDisposition = 4
	V1RefusalSampleFailure    V1RefusalDisposition = 5
	V1RefusalMalformedResult  V1RefusalDisposition = 6
	V1RefusalDiagnosticOutput V1RefusalDisposition = 7
	V1RefusalChildCrash       V1RefusalDisposition = 8
	V1RefusalChildSignal      V1RefusalDisposition = 9
	V1RefusalUnexpectedChild  V1RefusalDisposition = 10
	V1RefusalCleanupFailure   V1RefusalDisposition = 11
	V1RefusalMixedUpdate      V1RefusalDisposition = 12
)

const (
	V1ParseValid      = ParseValid
	V1PolicyAllow     = PolicyAllow
	V1ResultNoFinding = ResultNoFinding
)

func NewV1ValidationResult(request *V1ValidationRequest, parse ParseDisposition, policy PolicyDisposition, class ResultClassification, counts [5]uint32) (*V1ValidationResult, error) {
	if err := validateV1Request(request); err != nil {
		return nil, err
	}
	_, identity := v1ResultIdentities(request.Role)
	return &V1ValidationResult{
		ProtocolVersion: V1ProtocolVersion, Role: request.Role, Method: request.Method, ResultIdentity: identity,
		CorrelationID: request.CorrelationID, InstallationID: request.InstallationID,
		EpochSequence: request.EpochSequence, EpochDigest: request.EpochDigest,
		SourceByteLength: request.SourceByteLength, SourceDigest: request.SourceDigest,
		ArtifactProfileDigest: request.ArtifactProfileDigest, ResourcePolicyDigest: request.ResourcePolicyDigest,
		ParseDisposition: parse, PolicyDisposition: policy, Classification: class,
		CleanupDisposition: V1CleanupComplete, RefusalDisposition: V1RefusalNone,
		StaticImportCount: counts[0], ExportFromCount: counts[1], ImportExpressionCount: counts[2],
		ImportMetaCount: counts[3], FreeCommonJSReferenceCount: counts[4],
	}, nil
}

func EncodeV1Result(result *V1ValidationResult) ([]byte, error) {
	if err := validateV1Result(result); err != nil {
		return nil, err
	}
	frame := make([]byte, V1ResultFrameBytes)
	v1PutLengthMagic(frame, v1ResultMagic)
	v1PutResultCommon(frame, result)
	frame[216], frame[217], frame[218], frame[219] = byte(result.ParseDisposition), byte(result.PolicyDisposition), byte(result.Classification), byte(result.CleanupDisposition)
	binary.BigEndian.PutUint16(frame[220:222], uint16(result.RefusalDisposition))
	counts := [5]uint32{result.StaticImportCount, result.ExportFromCount, result.ImportExpressionCount, result.ImportMetaCount, result.FreeCommonJSReferenceCount}
	for index, count := range counts {
		binary.BigEndian.PutUint32(frame[224+index*4:228+index*4], count)
	}
	return frame, nil
}

func DecodeV1Result(expectedRole V1ConsumerRole, frame []byte) (*V1ValidationResult, error) {
	if err := v1ValidateFrame(frame, v1ResultMagic, V1ResultFrameBytes, V1ResultFrameBytes, V1ResultFrameKind); err != nil {
		return nil, err
	}
	role := V1ConsumerRole(binary.BigEndian.Uint16(frame[16:18]))
	if err := v1RequireRole(expectedRole, role); err != nil {
		return nil, err
	}
	if err := v1ValidateResultTags(frame, role); err != nil {
		return nil, err
	}
	result := &V1ValidationResult{ProtocolVersion: V1ProtocolVersion, Role: role}
	result.Method, result.ResultIdentity = v1ResultIdentities(role)
	copy(result.CorrelationID[:], frame[44:60])
	copy(result.InstallationID[:], frame[60:76])
	result.EpochSequence = binary.BigEndian.Uint64(frame[76:84])
	copy(result.EpochDigest[:], frame[84:116])
	result.SourceByteLength = binary.BigEndian.Uint32(frame[116:120])
	copy(result.SourceDigest[:], frame[120:152])
	copy(result.ArtifactProfileDigest[:], frame[152:184])
	copy(result.ResourcePolicyDigest[:], frame[184:216])
	result.ParseDisposition, result.PolicyDisposition = ParseDisposition(frame[216]), PolicyDisposition(frame[217])
	result.Classification, result.CleanupDisposition = ResultClassification(frame[218]), V1CleanupDisposition(frame[219])
	result.RefusalDisposition = V1RefusalDisposition(binary.BigEndian.Uint16(frame[220:222]))
	result.StaticImportCount = binary.BigEndian.Uint32(frame[224:228])
	result.ExportFromCount = binary.BigEndian.Uint32(frame[228:232])
	result.ImportExpressionCount = binary.BigEndian.Uint32(frame[232:236])
	result.ImportMetaCount = binary.BigEndian.Uint32(frame[236:240])
	result.FreeCommonJSReferenceCount = binary.BigEndian.Uint32(frame[240:244])
	if err := validateV1Result(result); err != nil {
		return nil, err
	}
	return result, nil
}

func VerifyV1Result(request *V1ValidationRequest, frame []byte) (*V1ValidationResult, error) {
	if err := validateV1Request(request); err != nil {
		return nil, err
	}
	result, err := DecodeV1Result(request.Role, frame)
	if err != nil {
		return nil, err
	}
	if result.CorrelationID != request.CorrelationID || result.InstallationID != request.InstallationID ||
		result.EpochSequence != request.EpochSequence || result.EpochDigest != request.EpochDigest ||
		result.SourceByteLength != request.SourceByteLength || result.SourceDigest != request.SourceDigest ||
		result.ArtifactProfileDigest != request.ArtifactProfileDigest || result.ResourcePolicyDigest != request.ResourcePolicyDigest {
		return nil, reject(ClassificationBinding, "v1-result-request-binding")
	}
	return result, nil
}

type V1PolicyActivation uint16

const V1PolicyInactive V1PolicyActivation = 0

// field-authority-object: capsule.source-validator-reactive-footprint-policy v1
type V1ResourcePolicy struct {
	RecordVersion                   uint16
	Role                            V1ConsumerRole
	Activation                      V1PolicyActivation
	ProcessProfileDigest            V1Digest
	ArtifactProfileDigest           V1Digest
	SupportedHostBuildDigest        V1Digest
	DirectChildMaximum              uint16
	DescendantMaximum               uint16
	ActiveRequestMaximum            uint16
	CombinedRoleChildMaximum        uint16
	RequestByteMaximum              uint32
	ResultByteMaximum               uint32
	DiagnosticByteMaximum           uint32
	ConnectionMaximum               uint16
	QueueMaximum                    uint16
	InFlightByteMaximum             uint32
	CPUMillisecondMaximum           uint32
	WallMillisecondMaximum          uint32
	DescriptorMaximum               uint16
	FileMaximum                     uint16
	ProcessMaximum                  uint16
	CoreByteMaximum                 uint32
	StackByteMaximum                uint32
	ObservedFootprintThreshold      uint64
	SamplingIntervalNanoseconds     uint64
	MeasuredCleanBaseline           uint64
	MaximumObservedOvershoot        uint64
	MaximumObservedKillLatencyNanos uint64
	SystemPressureDisposition       uint16
	KillDisposition                 uint16
	DrainDisposition                uint16
	ReapDisposition                 uint16
	RestartDisposition              uint16
	UpdateDisposition               uint16
	StartupDisposition              uint16
	ResidueCleanupDisposition       uint16
}

func NewV1InactiveResourcePolicy(role V1ConsumerRole) *V1ResourcePolicy {
	return &V1ResourcePolicy{
		RecordVersion: V1ProtocolVersion, Role: role, Activation: V1PolicyInactive,
		DirectChildMaximum: 1, DescendantMaximum: 0, ActiveRequestMaximum: 1, CombinedRoleChildMaximum: 2,
		RequestByteMaximum: V1RequestMaximumBytes, ResultByteMaximum: V1ResultFrameBytes,
		DiagnosticByteMaximum: 0, ConnectionMaximum: 1, QueueMaximum: 0,
		InFlightByteMaximum: V1RequestMaximumBytes + V1ResultFrameBytes,
		ProcessMaximum:      1,
	}
}

func EncodeV1ResourcePolicy(policy *V1ResourcePolicy) ([]byte, error) {
	if err := validateV1ResourcePolicy(policy); err != nil {
		return nil, err
	}
	frame := make([]byte, V1ResourcePolicyFrameBytes)
	v1PutLengthMagic(frame, v1PolicyMagic)
	binary.BigEndian.PutUint16(frame[12:14], V1ProtocolVersion)
	binary.BigEndian.PutUint16(frame[14:16], V1ResourcePolicyFrameKind)
	binary.BigEndian.PutUint16(frame[16:18], uint16(policy.Role))
	binary.BigEndian.PutUint16(frame[18:20], uint16(policy.Activation))
	binary.BigEndian.PutUint16(frame[20:22], 1)
	copy(frame[32:64], policy.ProcessProfileDigest[:])
	copy(frame[64:96], policy.ArtifactProfileDigest[:])
	copy(frame[96:128], policy.SupportedHostBuildDigest[:])
	binary.BigEndian.PutUint16(frame[128:130], policy.DirectChildMaximum)
	binary.BigEndian.PutUint16(frame[130:132], policy.DescendantMaximum)
	binary.BigEndian.PutUint16(frame[132:134], policy.ActiveRequestMaximum)
	binary.BigEndian.PutUint16(frame[134:136], policy.CombinedRoleChildMaximum)
	binary.BigEndian.PutUint32(frame[136:140], policy.RequestByteMaximum)
	binary.BigEndian.PutUint32(frame[140:144], policy.ResultByteMaximum)
	binary.BigEndian.PutUint32(frame[144:148], policy.DiagnosticByteMaximum)
	binary.BigEndian.PutUint16(frame[148:150], policy.ConnectionMaximum)
	binary.BigEndian.PutUint16(frame[150:152], policy.QueueMaximum)
	binary.BigEndian.PutUint32(frame[152:156], policy.InFlightByteMaximum)
	binary.BigEndian.PutUint32(frame[156:160], policy.CPUMillisecondMaximum)
	binary.BigEndian.PutUint32(frame[160:164], policy.WallMillisecondMaximum)
	binary.BigEndian.PutUint16(frame[164:166], policy.DescriptorMaximum)
	binary.BigEndian.PutUint16(frame[166:168], policy.FileMaximum)
	binary.BigEndian.PutUint16(frame[168:170], policy.ProcessMaximum)
	binary.BigEndian.PutUint32(frame[170:174], policy.CoreByteMaximum)
	binary.BigEndian.PutUint32(frame[174:178], policy.StackByteMaximum)
	binary.BigEndian.PutUint64(frame[178:186], policy.ObservedFootprintThreshold)
	binary.BigEndian.PutUint64(frame[186:194], policy.SamplingIntervalNanoseconds)
	binary.BigEndian.PutUint64(frame[194:202], policy.MeasuredCleanBaseline)
	binary.BigEndian.PutUint64(frame[202:210], policy.MaximumObservedOvershoot)
	binary.BigEndian.PutUint64(frame[210:218], policy.MaximumObservedKillLatencyNanos)
	dispositions := []uint16{policy.SystemPressureDisposition, policy.KillDisposition, policy.DrainDisposition, policy.ReapDisposition, policy.RestartDisposition, policy.UpdateDisposition, policy.StartupDisposition, policy.ResidueCleanupDisposition}
	for index, value := range dispositions {
		binary.BigEndian.PutUint16(frame[218+index*2:220+index*2], value)
	}
	return frame, nil
}

func DecodeV1ResourcePolicy(expectedRole V1ConsumerRole, frame []byte) (*V1ResourcePolicy, error) {
	if err := v1ValidateFrame(frame, v1PolicyMagic, V1ResourcePolicyFrameBytes, V1ResourcePolicyFrameBytes, V1ResourcePolicyFrameKind); err != nil {
		return nil, err
	}
	role := V1ConsumerRole(binary.BigEndian.Uint16(frame[16:18]))
	if err := v1RequireRole(expectedRole, role); err != nil {
		return nil, err
	}
	if binary.BigEndian.Uint16(frame[20:22]) != 1 || !allZero(frame[22:32]) || !allZero(frame[234:]) {
		return nil, reject(ClassificationSchema, "v1-policy-reserved")
	}
	policy := &V1ResourcePolicy{RecordVersion: V1ProtocolVersion, Role: role, Activation: V1PolicyActivation(binary.BigEndian.Uint16(frame[18:20]))}
	copy(policy.ProcessProfileDigest[:], frame[32:64])
	copy(policy.ArtifactProfileDigest[:], frame[64:96])
	copy(policy.SupportedHostBuildDigest[:], frame[96:128])
	policy.DirectChildMaximum = binary.BigEndian.Uint16(frame[128:130])
	policy.DescendantMaximum = binary.BigEndian.Uint16(frame[130:132])
	policy.ActiveRequestMaximum = binary.BigEndian.Uint16(frame[132:134])
	policy.CombinedRoleChildMaximum = binary.BigEndian.Uint16(frame[134:136])
	policy.RequestByteMaximum = binary.BigEndian.Uint32(frame[136:140])
	policy.ResultByteMaximum = binary.BigEndian.Uint32(frame[140:144])
	policy.DiagnosticByteMaximum = binary.BigEndian.Uint32(frame[144:148])
	policy.ConnectionMaximum = binary.BigEndian.Uint16(frame[148:150])
	policy.QueueMaximum = binary.BigEndian.Uint16(frame[150:152])
	policy.InFlightByteMaximum = binary.BigEndian.Uint32(frame[152:156])
	policy.CPUMillisecondMaximum = binary.BigEndian.Uint32(frame[156:160])
	policy.WallMillisecondMaximum = binary.BigEndian.Uint32(frame[160:164])
	policy.DescriptorMaximum = binary.BigEndian.Uint16(frame[164:166])
	policy.FileMaximum = binary.BigEndian.Uint16(frame[166:168])
	policy.ProcessMaximum = binary.BigEndian.Uint16(frame[168:170])
	policy.CoreByteMaximum = binary.BigEndian.Uint32(frame[170:174])
	policy.StackByteMaximum = binary.BigEndian.Uint32(frame[174:178])
	policy.ObservedFootprintThreshold = binary.BigEndian.Uint64(frame[178:186])
	policy.SamplingIntervalNanoseconds = binary.BigEndian.Uint64(frame[186:194])
	policy.MeasuredCleanBaseline = binary.BigEndian.Uint64(frame[194:202])
	policy.MaximumObservedOvershoot = binary.BigEndian.Uint64(frame[202:210])
	policy.MaximumObservedKillLatencyNanos = binary.BigEndian.Uint64(frame[210:218])
	values := []*uint16{&policy.SystemPressureDisposition, &policy.KillDisposition, &policy.DrainDisposition, &policy.ReapDisposition, &policy.RestartDisposition, &policy.UpdateDisposition, &policy.StartupDisposition, &policy.ResidueCleanupDisposition}
	for index, value := range values {
		*value = binary.BigEndian.Uint16(frame[218+index*2 : 220+index*2])
	}
	if err := validateV1ResourcePolicy(policy); err != nil {
		return nil, err
	}
	return policy, nil
}

// field-authority-object: capsule.source-validator-process-profile v1
type V1ProcessProfile struct {
	RecordVersion         uint16
	Role                  V1ConsumerRole
	ProcessProfile        string
	ContainingBundle      string
	PrivateService        string
	ParserSigningIdentity string
	Method                string
	ResourcePolicyDigest  V1Digest
	ExecutableDigest      V1Digest
	BuildManifestDigest   V1Digest
	SignatureDigest       V1Digest
	EntitlementsDigest    V1Digest
	ConstraintsDigest     V1Digest
	AssessmentDigest      V1Digest
}

// field-authority-object: capsule.source-validator-artifact-profile v1
type V1ArtifactProfile struct {
	RecordVersion        uint16
	Role                 V1ConsumerRole
	ArtifactProfile      string
	ProcessProfileDigest V1Digest
	ResourcePolicyDigest V1Digest
	ExecutableDigest     V1Digest
	BuildManifestDigest  V1Digest
	SignatureDigest      V1Digest
	EntitlementsDigest   V1Digest
	ConstraintsDigest    V1Digest
}

// field-authority-object: capsule.source-validator-consumer-projection v1
type V1ConsumerProjection struct {
	RecordVersion         uint16
	Role                  V1ConsumerRole
	ContainingBundle      string
	Method                string
	InstallationID        V1InstallationID
	EpochSequence         uint64
	EpochDigest           V1Digest
	ArtifactProfileDigest V1Digest
	ResourcePolicyDigest  V1Digest
	ProcessProfileDigest  V1Digest
}

func NewFixtureOnlyV1ProcessProfile(role V1ConsumerRole, policyDigest [32]byte) *V1ProcessProfile {
	profile, bundle, service, parser, method := v1RoleProfileIdentities(role)
	value := &V1ProcessProfile{RecordVersion: V1ProtocolVersion, Role: role, ProcessProfile: profile, ContainingBundle: bundle, PrivateService: service, ParserSigningIdentity: parser, Method: method, ResourcePolicyDigest: V1Digest(policyDigest)}
	fillFixtureDigests([]*V1Digest{&value.ExecutableDigest, &value.BuildManifestDigest, &value.SignatureDigest, &value.EntitlementsDigest, &value.ConstraintsDigest, &value.AssessmentDigest}, role, 0x60)
	return value
}

func NewFixtureOnlyV1ArtifactProfile(role V1ConsumerRole, processDigest, policyDigest [32]byte) *V1ArtifactProfile {
	identity := V1DaemonArtifactProfileIdentity
	if role == V1ApprovalBrokerRole {
		identity = V1BrokerArtifactProfileIdentity
	}
	value := &V1ArtifactProfile{RecordVersion: V1ProtocolVersion, Role: role, ArtifactProfile: identity, ProcessProfileDigest: V1Digest(processDigest), ResourcePolicyDigest: V1Digest(policyDigest)}
	fillFixtureDigests([]*V1Digest{&value.ExecutableDigest, &value.BuildManifestDigest, &value.SignatureDigest, &value.EntitlementsDigest, &value.ConstraintsDigest}, role, 0x70)
	return value
}

func NewFixtureOnlyV1ConsumerProjection(role V1ConsumerRole, artifactDigest, policyDigest [32]byte) *V1ConsumerProjection {
	_, bundle, _, _, method := v1RoleProfileIdentities(role)
	value := &V1ConsumerProjection{RecordVersion: V1ProtocolVersion, Role: role, ContainingBundle: bundle, Method: method, EpochSequence: 7, ArtifactProfileDigest: V1Digest(artifactDigest), ResourcePolicyDigest: V1Digest(policyDigest)}
	value.InstallationID[0], value.EpochDigest[0], value.ProcessProfileDigest[0] = 0x11, 0x22, 0x80+v1FixtureRoleByte(role)
	return value
}

func EncodeV1ProcessProfile(value *V1ProcessProfile) ([]byte, error) {
	return encodeV1Profile(value, V1ProcessProfileFrameKind, v1ProcessMagic)
}
func EncodeV1ArtifactProfile(value *V1ArtifactProfile) ([]byte, error) {
	return encodeV1Profile(value, V1ArtifactProfileFrameKind, v1ArtifactMagic)
}
func EncodeV1ConsumerProjection(value *V1ConsumerProjection) ([]byte, error) {
	return encodeV1Consumer(value)
}

func DecodeV1ProcessProfile(role V1ConsumerRole, frame []byte) (*V1ProcessProfile, error) {
	return decodeV1ProcessProfile(role, frame)
}
func DecodeV1ArtifactProfile(role V1ConsumerRole, frame []byte) (*V1ArtifactProfile, error) {
	return decodeV1ArtifactProfile(role, frame)
}
func DecodeV1ConsumerProjection(role V1ConsumerRole, frame []byte) (*V1ConsumerProjection, error) {
	return decodeV1Consumer(role, frame)
}

func encodeV1Profile(value any, kind uint16, magic [8]byte) ([]byte, error) {
	frame := make([]byte, V1ProcessProfileFrameBytes)
	v1PutLengthMagic(frame, magic)
	var role V1ConsumerRole
	var digests []V1Digest
	switch typed := value.(type) {
	case *V1ProcessProfile:
		if err := validateV1ProcessProfile(typed); err != nil {
			return nil, err
		}
		role = typed.Role
		digests = []V1Digest{typed.ResourcePolicyDigest, typed.ExecutableDigest, typed.BuildManifestDigest, typed.SignatureDigest, typed.EntitlementsDigest, typed.ConstraintsDigest, typed.AssessmentDigest}
	case *V1ArtifactProfile:
		if err := validateV1ArtifactProfile(typed); err != nil {
			return nil, err
		}
		role = typed.Role
		digests = []V1Digest{typed.ProcessProfileDigest, typed.ResourcePolicyDigest, typed.ExecutableDigest, typed.BuildManifestDigest, typed.SignatureDigest, typed.EntitlementsDigest, typed.ConstraintsDigest}
	default:
		return nil, reject(ClassificationSchema, "v1-profile-type")
	}
	binary.BigEndian.PutUint16(frame[12:14], V1ProtocolVersion)
	binary.BigEndian.PutUint16(frame[14:16], kind)
	binary.BigEndian.PutUint16(frame[16:18], uint16(role))
	for index := 0; index < 6; index++ {
		binary.BigEndian.PutUint16(frame[18+index*2:20+index*2], v1RoleTag(role, uint16(index+1)))
	}
	for index, digest := range digests {
		copy(frame[32+index*32:64+index*32], digest[:])
	}
	return frame, nil
}

func decodeV1ProcessProfile(expected V1ConsumerRole, frame []byte) (*V1ProcessProfile, error) {
	if err := v1ValidateProfileFrame(expected, frame, v1ProcessMagic, V1ProcessProfileFrameKind); err != nil {
		return nil, err
	}
	profile, bundle, service, parser, method := v1RoleProfileIdentities(expected)
	value := &V1ProcessProfile{RecordVersion: 1, Role: expected, ProcessProfile: profile, ContainingBundle: bundle, PrivateService: service, ParserSigningIdentity: parser, Method: method}
	digests := []*V1Digest{&value.ResourcePolicyDigest, &value.ExecutableDigest, &value.BuildManifestDigest, &value.SignatureDigest, &value.EntitlementsDigest, &value.ConstraintsDigest, &value.AssessmentDigest}
	for index, digest := range digests {
		copy(digest[:], frame[32+index*32:64+index*32])
	}
	if err := validateV1ProcessProfile(value); err != nil {
		return nil, err
	}
	return value, nil
}

func decodeV1ArtifactProfile(expected V1ConsumerRole, frame []byte) (*V1ArtifactProfile, error) {
	if err := v1ValidateProfileFrame(expected, frame, v1ArtifactMagic, V1ArtifactProfileFrameKind); err != nil {
		return nil, err
	}
	identity := V1DaemonArtifactProfileIdentity
	if expected == V1ApprovalBrokerRole {
		identity = V1BrokerArtifactProfileIdentity
	}
	value := &V1ArtifactProfile{RecordVersion: 1, Role: expected, ArtifactProfile: identity}
	digests := []*V1Digest{&value.ProcessProfileDigest, &value.ResourcePolicyDigest, &value.ExecutableDigest, &value.BuildManifestDigest, &value.SignatureDigest, &value.EntitlementsDigest, &value.ConstraintsDigest}
	for index, digest := range digests {
		copy(digest[:], frame[32+index*32:64+index*32])
	}
	if err := validateV1ArtifactProfile(value); err != nil {
		return nil, err
	}
	return value, nil
}

func encodeV1Consumer(value *V1ConsumerProjection) ([]byte, error) {
	if err := validateV1Consumer(value); err != nil {
		return nil, err
	}
	frame := make([]byte, V1ConsumerProjectionFrameBytes)
	v1PutLengthMagic(frame, v1ConsumerMagic)
	binary.BigEndian.PutUint16(frame[12:14], 1)
	binary.BigEndian.PutUint16(frame[14:16], V1ConsumerProjectionFrameKind)
	binary.BigEndian.PutUint16(frame[16:18], uint16(value.Role))
	for index := 0; index < 6; index++ {
		binary.BigEndian.PutUint16(frame[18+index*2:20+index*2], v1RoleTag(value.Role, uint16(index+1)))
	}
	copy(frame[32:48], value.InstallationID[:])
	binary.BigEndian.PutUint64(frame[48:56], value.EpochSequence)
	copy(frame[56:88], value.EpochDigest[:])
	copy(frame[88:120], value.ArtifactProfileDigest[:])
	copy(frame[120:152], value.ResourcePolicyDigest[:])
	copy(frame[152:184], value.ProcessProfileDigest[:])
	return frame, nil
}

func decodeV1Consumer(expected V1ConsumerRole, frame []byte) (*V1ConsumerProjection, error) {
	if err := v1ValidateFrame(frame, v1ConsumerMagic, V1ConsumerProjectionFrameBytes, V1ConsumerProjectionFrameBytes, V1ConsumerProjectionFrameKind); err != nil {
		return nil, err
	}
	role := V1ConsumerRole(binary.BigEndian.Uint16(frame[16:18]))
	if err := v1RequireRole(expected, role); err != nil {
		return nil, err
	}
	if err := v1ValidateRoleTags(frame, role, 6); err != nil {
		return nil, err
	}
	if !allZero(frame[30:32]) || !allZero(frame[184:]) {
		return nil, reject(ClassificationSchema, "v1-consumer-reserved")
	}
	_, bundle, _, _, method := v1RoleProfileIdentities(role)
	value := &V1ConsumerProjection{RecordVersion: 1, Role: role, ContainingBundle: bundle, Method: method}
	copy(value.InstallationID[:], frame[32:48])
	value.EpochSequence = binary.BigEndian.Uint64(frame[48:56])
	copy(value.EpochDigest[:], frame[56:88])
	copy(value.ArtifactProfileDigest[:], frame[88:120])
	copy(value.ResourcePolicyDigest[:], frame[120:152])
	copy(value.ProcessProfileDigest[:], frame[152:184])
	if err := validateV1Consumer(value); err != nil {
		return nil, err
	}
	return value, nil
}

func v1PutRequestCommon(frame []byte, request *V1ValidationRequest) {
	binary.BigEndian.PutUint16(frame[12:14], 1)
	binary.BigEndian.PutUint16(frame[14:16], V1RequestFrameKind)
	binary.BigEndian.PutUint16(frame[16:18], uint16(request.Role))
	for index := 0; index < 11; index++ {
		binary.BigEndian.PutUint16(frame[18+index*2:20+index*2], v1RoleTag(request.Role, uint16(index+1)))
	}
	copy(frame[44:60], request.CorrelationID[:])
	copy(frame[60:76], request.InstallationID[:])
	binary.BigEndian.PutUint64(frame[76:84], request.EpochSequence)
	copy(frame[84:116], request.EpochDigest[:])
	binary.BigEndian.PutUint32(frame[116:120], request.SourceByteLength)
	copy(frame[120:152], request.SourceDigest[:])
	copy(frame[152:184], request.ArtifactProfileDigest[:])
	copy(frame[184:216], request.ResourcePolicyDigest[:])
}

func v1PutResultCommon(frame []byte, result *V1ValidationResult) {
	binary.BigEndian.PutUint16(frame[12:14], 1)
	binary.BigEndian.PutUint16(frame[14:16], V1ResultFrameKind)
	binary.BigEndian.PutUint16(frame[16:18], uint16(result.Role))
	for index := 0; index < 11; index++ {
		binary.BigEndian.PutUint16(frame[18+index*2:20+index*2], v1RoleTag(result.Role, uint16(index+1)))
	}
	copy(frame[44:60], result.CorrelationID[:])
	copy(frame[60:76], result.InstallationID[:])
	binary.BigEndian.PutUint64(frame[76:84], result.EpochSequence)
	copy(frame[84:116], result.EpochDigest[:])
	binary.BigEndian.PutUint32(frame[116:120], result.SourceByteLength)
	copy(frame[120:152], result.SourceDigest[:])
	copy(frame[152:184], result.ArtifactProfileDigest[:])
	copy(frame[184:216], result.ResourcePolicyDigest[:])
}

func validateV1Request(request *V1ValidationRequest) error {
	if request == nil {
		return reject(ClassificationSchema, "v1-request-nil")
	}
	if err := validateV1Role(request.Role); err != nil {
		return err
	}
	method, identity := v1RequestIdentities(request.Role)
	if request.ProtocolVersion != 1 {
		return reject(ClassificationUnsupported, "v1-request-version")
	}
	if request.Method != method || request.RequestIdentity != identity || request.SourceProfile != SourceProfileIdentity || request.SourceMediaType != SourceMediaTypeIdentity {
		return reject(ClassificationDomain, "v1-request-family")
	}
	bindings := V1RequestBindings{request.CorrelationID, request.InstallationID, request.EpochSequence, request.EpochDigest, request.ArtifactProfileDigest, request.ResourcePolicyDigest}
	if err := validateV1Bindings(bindings); err != nil {
		return err
	}
	if len(request.sourceBytes) > v0candidate.MJSMainSourceMaxBytes || int(request.SourceByteLength) != len(request.sourceBytes) {
		return reject(ClassificationBinding, "v1-source-length")
	}
	digest := sha256.Sum256(request.sourceBytes)
	if V1Digest(digest) != request.SourceDigest {
		return reject(ClassificationBinding, "v1-source-digest")
	}
	if _, err := v0candidate.NewMJSMainSource(request.sourceBytes); err != nil {
		return reject(ClassificationDomain, "v1-source-profile")
	}
	return nil
}

func validateV1Result(result *V1ValidationResult) error {
	if result == nil {
		return reject(ClassificationSchema, "v1-result-nil")
	}
	if err := validateV1Role(result.Role); err != nil {
		return err
	}
	method, identity := v1ResultIdentities(result.Role)
	if result.ProtocolVersion != 1 {
		return reject(ClassificationUnsupported, "v1-result-version")
	}
	if result.Method != method || result.ResultIdentity != identity {
		return reject(ClassificationDomain, "v1-result-family")
	}
	if result.CorrelationID == (V1CorrelationID{}) || result.InstallationID == (V1InstallationID{}) || result.EpochSequence == 0 || result.EpochDigest == (V1Digest{}) || result.SourceDigest == (V1Digest{}) || result.ArtifactProfileDigest == (V1Digest{}) || result.ResourcePolicyDigest == (V1Digest{}) {
		return reject(ClassificationSchema, "v1-result-zero-binding")
	}
	if result.SourceByteLength > v0candidate.MJSMainSourceMaxBytes {
		return reject(ClassificationDomain, "v1-result-source-cap")
	}
	if result.ParseDisposition != ParseValid || result.PolicyDisposition != PolicyAllow || result.Classification != ResultNoFinding {
		return reject(ClassificationSchema, "v1-result-disposition")
	}
	if result.CleanupDisposition != V1CleanupComplete || result.RefusalDisposition != V1RefusalNone {
		return reject(ClassificationSchema, "v1-result-cleanup")
	}
	if result.StaticImportCount|result.ExportFromCount|result.ImportExpressionCount|result.ImportMetaCount|result.FreeCommonJSReferenceCount != 0 {
		return reject(ClassificationSchema, "v1-result-count-on-allow")
	}
	return nil
}

func validateV1ResourcePolicy(policy *V1ResourcePolicy) error {
	if policy == nil {
		return reject(ClassificationSchema, "v1-policy-nil")
	}
	if err := validateV1Role(policy.Role); err != nil {
		return err
	}
	if policy.RecordVersion != 1 {
		return reject(ClassificationUnsupported, "v1-policy-version")
	}
	if policy.Activation != V1PolicyInactive {
		return reject(ClassificationUnsupported, "v1-policy-active-unmeasured")
	}
	if policy.ProcessProfileDigest != (V1Digest{}) || policy.ArtifactProfileDigest != (V1Digest{}) || policy.SupportedHostBuildDigest != (V1Digest{}) || policy.ObservedFootprintThreshold != 0 || policy.SamplingIntervalNanoseconds != 0 || policy.MeasuredCleanBaseline != 0 || policy.MaximumObservedOvershoot != 0 || policy.MaximumObservedKillLatencyNanos != 0 || policy.CPUMillisecondMaximum != 0 || policy.WallMillisecondMaximum != 0 || policy.DescriptorMaximum != 0 || policy.FileMaximum != 0 || policy.CoreByteMaximum != 0 || policy.StackByteMaximum != 0 || policy.SystemPressureDisposition != 0 || policy.KillDisposition != 0 || policy.DrainDisposition != 0 || policy.ReapDisposition != 0 || policy.RestartDisposition != 0 || policy.UpdateDisposition != 0 || policy.StartupDisposition != 0 || policy.ResidueCleanupDisposition != 0 {
		return reject(ClassificationSchema, "v1-policy-inactive-measurement")
	}
	if policy.DirectChildMaximum != 1 || policy.DescendantMaximum != 0 || policy.ActiveRequestMaximum != 1 || policy.CombinedRoleChildMaximum != 2 || policy.RequestByteMaximum != V1RequestMaximumBytes || policy.ResultByteMaximum != V1ResultFrameBytes || policy.DiagnosticByteMaximum != 0 || policy.ConnectionMaximum != 1 || policy.QueueMaximum != 0 || policy.InFlightByteMaximum != V1RequestMaximumBytes+V1ResultFrameBytes || policy.ProcessMaximum != 1 {
		return reject(ClassificationBinding, "v1-policy-fixed-bound")
	}
	return nil
}

func validateV1ProcessProfile(value *V1ProcessProfile) error {
	if value == nil {
		return reject(ClassificationSchema, "v1-process-nil")
	}
	if err := validateV1Role(value.Role); err != nil {
		return err
	}
	profile, bundle, service, parser, method := v1RoleProfileIdentities(value.Role)
	if value.RecordVersion != 1 {
		return reject(ClassificationUnsupported, "v1-process-version")
	}
	if value.ProcessProfile != profile || value.ContainingBundle != bundle || value.PrivateService != service || value.ParserSigningIdentity != parser || value.Method != method {
		return reject(ClassificationDomain, "v1-process-family")
	}
	return requireNonzeroDigests("v1-process-digest", value.ResourcePolicyDigest, value.ExecutableDigest, value.BuildManifestDigest, value.SignatureDigest, value.EntitlementsDigest, value.ConstraintsDigest, value.AssessmentDigest)
}
func validateV1ArtifactProfile(value *V1ArtifactProfile) error {
	if value == nil {
		return reject(ClassificationSchema, "v1-artifact-nil")
	}
	if err := validateV1Role(value.Role); err != nil {
		return err
	}
	identity := V1DaemonArtifactProfileIdentity
	if value.Role == V1ApprovalBrokerRole {
		identity = V1BrokerArtifactProfileIdentity
	}
	if value.RecordVersion != 1 {
		return reject(ClassificationUnsupported, "v1-artifact-version")
	}
	if value.ArtifactProfile != identity {
		return reject(ClassificationDomain, "v1-artifact-family")
	}
	return requireNonzeroDigests("v1-artifact-digest", value.ProcessProfileDigest, value.ResourcePolicyDigest, value.ExecutableDigest, value.BuildManifestDigest, value.SignatureDigest, value.EntitlementsDigest, value.ConstraintsDigest)
}
func validateV1Consumer(value *V1ConsumerProjection) error {
	if value == nil {
		return reject(ClassificationSchema, "v1-consumer-nil")
	}
	if err := validateV1Role(value.Role); err != nil {
		return err
	}
	_, bundle, _, _, method := v1RoleProfileIdentities(value.Role)
	if value.RecordVersion != 1 {
		return reject(ClassificationUnsupported, "v1-consumer-version")
	}
	if value.ContainingBundle != bundle || value.Method != method {
		return reject(ClassificationDomain, "v1-consumer-family")
	}
	if value.InstallationID == (V1InstallationID{}) || value.EpochSequence == 0 || value.EpochDigest == (V1Digest{}) {
		return reject(ClassificationSchema, "v1-consumer-zero-binding")
	}
	return requireNonzeroDigests("v1-consumer-digest", value.ArtifactProfileDigest, value.ResourcePolicyDigest, value.ProcessProfileDigest)
}

func v1ValidateFrame(frame []byte, magic [8]byte, minimum, maximum int, kind uint16) error {
	if len(frame) >= 12 && (bytes.Equal(frame[4:12], v0RequestMagicForRefusal[:]) || bytes.Equal(frame[4:12], v0ResultMagicForRefusal[:])) {
		return reject(ClassificationUnsupported, "v0-as-v1")
	}
	if len(frame) < minimum {
		return reject(ClassificationMalformed, "v1-frame-truncated")
	}
	if len(frame) > maximum {
		return reject(ClassificationMalformed, "v1-frame-cap")
	}
	if !bytes.Equal(frame[4:12], magic[:]) {
		return reject(ClassificationMalformed, "v1-frame-magic")
	}
	if int(binary.BigEndian.Uint32(frame[:4])) != len(frame)-4 {
		return reject(ClassificationMalformed, "v1-frame-length")
	}
	if binary.BigEndian.Uint16(frame[12:14]) != 1 {
		return reject(ClassificationUnsupported, "v1-frame-version")
	}
	if binary.BigEndian.Uint16(frame[14:16]) != kind {
		return reject(ClassificationDomain, "v1-frame-kind")
	}
	return nil
}
func v1ValidateRequestTags(frame []byte, role V1ConsumerRole) error {
	if err := v1ValidateRoleTags(frame, role, 11); err != nil {
		return err
	}
	if !allZero(frame[40:44]) {
		return reject(ClassificationSchema, "v1-request-reserved")
	}
	return nil
}
func v1ValidateResultTags(frame []byte, role V1ConsumerRole) error {
	if err := v1ValidateRoleTags(frame, role, 11); err != nil {
		return err
	}
	if !allZero(frame[40:44]) || !allZero(frame[222:224]) || !allZero(frame[244:248]) {
		return reject(ClassificationSchema, "v1-result-reserved")
	}
	return nil
}
func v1ValidateRoleTags(frame []byte, role V1ConsumerRole, count int) error {
	for index := 0; index < count; index++ {
		if binary.BigEndian.Uint16(frame[18+index*2:20+index*2]) != v1RoleTag(role, uint16(index+1)) {
			return reject(ClassificationDomain, "v1-role-family-tag")
		}
	}
	return nil
}
func v1ValidateProfileFrame(role V1ConsumerRole, frame []byte, magic [8]byte, kind uint16) error {
	if err := v1ValidateFrame(frame, magic, V1ProcessProfileFrameBytes, V1ProcessProfileFrameBytes, kind); err != nil {
		return err
	}
	actual := V1ConsumerRole(binary.BigEndian.Uint16(frame[16:18]))
	if err := v1RequireRole(role, actual); err != nil {
		return err
	}
	if err := v1ValidateRoleTags(frame, role, 6); err != nil {
		return err
	}
	if !allZero(frame[30:32]) {
		return reject(ClassificationSchema, "v1-profile-reserved")
	}
	return nil
}
func v1PutLengthMagic(frame []byte, magic [8]byte) {
	// Every caller allocates one of the closed frame sizes, all of which are
	// bounded by V1RequestMaximumBytes and therefore fit in uint32.
	binary.BigEndian.PutUint32(frame[:4], uint32(len(frame)-4)) //nolint:gosec
	copy(frame[4:12], magic[:])
}
func v1RoleTag(role V1ConsumerRole, field uint16) uint16 { return uint16(role)*0x100 + field }
func validateV1Role(role V1ConsumerRole) error {
	if role != V1DaemonRole && role != V1ApprovalBrokerRole {
		return reject(ClassificationDomain, "v1-consumer-role")
	}
	return nil
}
func v1RequireRole(expected, actual V1ConsumerRole) error {
	if err := validateV1Role(expected); err != nil {
		return err
	}
	if actual != expected {
		return reject(ClassificationDomain, "v1-cross-role")
	}
	return nil
}
func validateV1Bindings(value V1RequestBindings) error {
	if value.CorrelationID == (V1CorrelationID{}) || value.InstallationID == (V1InstallationID{}) || value.EpochSequence == 0 || value.EpochDigest == (V1Digest{}) || value.ArtifactProfileDigest == (V1Digest{}) || value.ResourcePolicyDigest == (V1Digest{}) {
		return reject(ClassificationSchema, "v1-zero-binding")
	}
	return nil
}
func requireNonzeroDigests(code string, values ...V1Digest) error {
	for _, value := range values {
		if value == (V1Digest{}) {
			return reject(ClassificationSchema, code)
		}
	}
	return nil
}
func allZero(value []byte) bool {
	for _, candidate := range value {
		if candidate != 0 {
			return false
		}
	}
	return true
}
func fillFixtureDigests(values []*V1Digest, role V1ConsumerRole, base byte) {
	roleByte := v1FixtureRoleByte(role)
	for index, value := range values {
		for offset := range value {
			value[offset] = base + byte(index) + roleByte
		}
	}
}
func v1FixtureRoleByte(role V1ConsumerRole) byte {
	switch role {
	case V1DaemonRole:
		return 1
	case V1ApprovalBrokerRole:
		return 2
	default:
		return 0
	}
}
func v1RequestIdentities(role V1ConsumerRole) (string, string) {
	if role == V1DaemonRole {
		return V1DaemonMethodIdentity, V1DaemonRequestIdentity
	}
	return V1BrokerMethodIdentity, V1BrokerRequestIdentity
}
func v1ResultIdentities(role V1ConsumerRole) (string, string) {
	if role == V1DaemonRole {
		return V1DaemonMethodIdentity, V1DaemonResultIdentity
	}
	return V1BrokerMethodIdentity, V1BrokerResultIdentity
}
func v1RoleProfileIdentities(role V1ConsumerRole) (string, string, string, string, string) {
	if role == V1DaemonRole {
		return V1DaemonProcessProfileIdentity, V1DaemonContainingBundleIdentity, V1DaemonServiceIdentity, V1DaemonParserIdentity, V1DaemonMethodIdentity
	}
	return V1BrokerProcessProfileIdentity, V1BrokerContainingBundleIdentity, V1BrokerServiceIdentity, V1BrokerParserIdentity, V1BrokerMethodIdentity
}
