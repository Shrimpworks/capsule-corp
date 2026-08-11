package authorityplane

// This file freezes the passive native XPC dictionary contract needed before
// an S3 harness can be built. It imports no XPC API, creates no listener or
// service, authenticates no peer, and dispatches no call. Each message tag is
// only a method-specific cross-check after the service has already selected
// one closed entry point; it is never a generic opcode.

type NativeXPCValueType string

// Passive native-XPC value types are closed to flat dictionary primitives.
const (
	NativeXPCTypeUInt64 NativeXPCValueType = "XPC_TYPE_UINT64"
	NativeXPCTypeData   NativeXPCValueType = "XPC_TYPE_DATA"
	NativeXPCTypeString NativeXPCValueType = "XPC_TYPE_STRING"
)

// Passive native-XPC field names, key counts, method constants, and limits.
const (
	NativeXPCEncoding                         = "xpc-dictionary-v0"
	NativeXPCEncodingVersion                  = uint64(0)
	NativeXPCProtocolVersionKey               = "capsule.protocol-version"
	NativeXPCMethodVersionKey                 = "capsule.method-version"
	NativeXPCMessageTagKey                    = "capsule.message-tag"
	NativeXPCRequestIDKey                     = "capsule.request-id"
	NativeXPCInstallationIDKey                = "capsule.installation-id"
	NativeXPCEpochSequenceKey                 = "capsule.epoch-sequence"
	NativeXPCEpochDigestKey                   = "capsule.epoch-digest"
	NativeXPCAudienceKey                      = "capsule.audience"
	NativeXPCPurposeKey                       = "capsule.purpose"
	NativeXPCStatusKey                        = "capsule.status"
	NativeXPCReasonKey                        = "capsule.reason"
	NativeXPCJobProposalKey                   = "capsule.job-proposal"
	NativeXPCExecutionPlanKey                 = "capsule.execution-plan"
	NativeXPCRoleBindingsKey                  = "capsule.role-bindings"
	NativeXPCSourceManifestKey                = "capsule.source-manifest"
	NativeXPCSourceKey                        = "capsule.source"
	NativeXPCRegistrationIDKey                = "capsule.registration-id"
	NativeXPCPlanRegistrationKey              = "capsule.plan-registration"
	NativeXPCApprovalEnvelopeKey              = "capsule.approval-envelope"
	NativeXPCApprovalIDKey                    = "capsule.approval-id"
	NativeXPCApprovalStateKey                 = "capsule.approval-state"
	NativeXPCAttemptIDKey                     = "capsule.attempt-id"
	NativeXPCAttemptStateKey                  = "capsule.attempt-state"
	NativeXPCRequestCommonKeyCount            = uint64(9)
	NativeXPCReplyCommonKeyCount              = uint64(6)
	NativeXPCSubmitRequestKeyCount            = uint64(10)
	NativeXPCSubmitReplyKeyCount              = uint64(7)
	NativeXPCRegisterRequestKeyCount          = uint64(13)
	NativeXPCRegisterReplyKeyCount            = uint64(7)
	NativeXPCGetRequestKeyCount               = uint64(10)
	NativeXPCGetReplyKeyCount                 = uint64(11)
	NativeXPCSubmitApprovalRequestKeyCount    = uint64(11)
	NativeXPCSubmitApprovalReplyKeyCount      = uint64(8)
	NativeXPCRequestAttemptRequestKeyCount    = uint64(11)
	NativeXPCRequestAttemptReplyKeyCount      = uint64(8)
	NativeXPCRefusalReplyKeyCount             = uint64(6)
	NativeXPCExtraObjectsAllowed              = uint64(0)
	NativeXPCFileDescriptorsAllowed           = uint64(0)
	NativeXPCEndpointsAllowed                 = uint64(0)
	NativeXPCMachRightsAllowed                = uint64(0)
	NativeXPCNestedContainersAllowed          = uint64(0)
	NativeXPCDeadlineBoundaryRule             = "dispatch-only-when-elapsed-milliseconds-is-strictly-less-than-deadline-milliseconds"
	NativeXPCDeadlineExpiredDisposition       = "deadline-expired-before-dispatch"
	NativeXPCDeadlineAfterDispatchDisposition = "response-unknown-store-semantic-result-or-recovery-fence-controls"
	SubmitApprovalV0Method                    = "SubmitApprovalV0"
	SubmitApprovalV0MethodVersion             = uint64(0)
	SubmitApprovalV0Purpose                   = "capsule.ipc.submit-approval.v0"
	SubmitApprovalV0RequestDataMaxBytes       = uint64(528)
	SubmitApprovalV0ReplyDataMaxBytes         = uint64(16)
	SubmitApprovalV0DeadlineMilliseconds      = uint64(5_000)
	SubmitApprovalV0ResponseLossDisposition   = "canonical-payload-and-signer-authorization-replay-returns-same-approval-and-current-state"
	RequestAttemptV0Method                    = "RequestAttemptV0"
	RequestAttemptV0MethodVersion             = uint64(0)
	RequestAttemptV0Purpose                   = "capsule.ipc.request-attempt.v0"
	RequestAttemptV0RequestDataMaxBytes       = uint64(32)
	RequestAttemptV0ReplyDataMaxBytes         = uint64(16)
	RequestAttemptV0DeadlineMilliseconds      = uint64(5_000)
	RequestAttemptV0ResponseLossDisposition   = "registration-and-approval-reference-replay-returns-same-attempt-and-current-state"
)

type NativeXPCMessageTag uint64

// Passive native-XPC method tags are stable cross-checks, not generic opcodes.
const (
	NativeXPCMessageTagInvalid             NativeXPCMessageTag = 0
	NativeXPCMessageTagSubmitMainMJSV0     NativeXPCMessageTag = 1
	NativeXPCMessageTagRegisterPlanV0      NativeXPCMessageTag = 2
	NativeXPCMessageTagGetRegisteredPlanV0 NativeXPCMessageTag = 3
	NativeXPCMessageTagSubmitApprovalV0    NativeXPCMessageTag = 4
	NativeXPCMessageTagRequestAttemptV0    NativeXPCMessageTag = 5
)

// NativeXPCApprovalStateTag is the closed numeric projection returned by SubmitApprovalV0.
type NativeXPCApprovalStateTag uint64

// SubmitApprovalV0 reply-state tags are closed; zero is never a success value.
const (
	NativeXPCApprovalStateInvalid     NativeXPCApprovalStateTag = 0
	NativeXPCApprovalStateUsable      NativeXPCApprovalStateTag = 1
	NativeXPCApprovalStateConsumed    NativeXPCApprovalStateTag = 2
	NativeXPCApprovalStateInvalidated NativeXPCApprovalStateTag = 3
)

// NativeXPCAttemptStateTag is the closed numeric projection returned by RequestAttemptV0.
type NativeXPCAttemptStateTag uint64

// RequestAttemptV0 reply-state tags are closed; zero is never a success value.
const (
	NativeXPCAttemptStateInvalid NativeXPCAttemptStateTag = 0
	NativeXPCAttemptStateCreated NativeXPCAttemptStateTag = 1
)

type NativeXPCStatusTag uint64

const (
	NativeXPCStatusOK               NativeXPCStatusTag = 0
	NativeXPCStatusMalformed        NativeXPCStatusTag = 1
	NativeXPCStatusUnsupported      NativeXPCStatusTag = 2
	NativeXPCStatusSchema           NativeXPCStatusTag = 3
	NativeXPCStatusBinding          NativeXPCStatusTag = 4
	NativeXPCStatusAuthentication   NativeXPCStatusTag = 5
	NativeXPCStatusStale            NativeXPCStatusTag = 6
	NativeXPCStatusReplay           NativeXPCStatusTag = 7
	NativeXPCStatusCapacity         NativeXPCStatusTag = 8
	NativeXPCStatusTrustState       NativeXPCStatusTag = 9
	NativeXPCStatusLocalFailure     NativeXPCStatusTag = 10
	NativeXPCStatusRecoveryRequired NativeXPCStatusTag = 11
	NativeXPCStatusSemantic         NativeXPCStatusTag = 12
	NativeXPCStatusDomain           NativeXPCStatusTag = 13
)

type NativeXPCReasonTag uint64

const (
	NativeXPCReasonNone                NativeXPCReasonTag = 0
	NativeXPCReasonKeySet              NativeXPCReasonTag = 1
	NativeXPCReasonValueType           NativeXPCReasonTag = 2
	NativeXPCReasonDataWidth           NativeXPCReasonTag = 3
	NativeXPCReasonDataCap             NativeXPCReasonTag = 4
	NativeXPCReasonZeroIdentifier      NativeXPCReasonTag = 5
	NativeXPCReasonEpochSequence       NativeXPCReasonTag = 6
	NativeXPCReasonProtocolVersion     NativeXPCReasonTag = 7
	NativeXPCReasonMethodVersion       NativeXPCReasonTag = 8
	NativeXPCReasonMessageTag          NativeXPCReasonTag = 9
	NativeXPCReasonMethodBinding       NativeXPCReasonTag = 10
	NativeXPCReasonCurrentState        NativeXPCReasonTag = 11
	NativeXPCReasonCapacity            NativeXPCReasonTag = 12
	NativeXPCReasonCoreRefusal         NativeXPCReasonTag = 13
	NativeXPCReasonLocalIntegrityFault NativeXPCReasonTag = 14
)

// NativeXPCFieldSpec describes one exact top-level XPC dictionary value. A
// zero data minimum is meaningful; fixed-width fields have equal minima and
// maxima. FixedUInt64 and FixedString are method-derived checks.
type NativeXPCFieldSpec struct {
	Key             string
	ValueType       NativeXPCValueType
	Required        bool
	MinDataBytes    uint64
	MaxDataBytes    uint64
	FixedUInt64     uint64
	HasFixedUInt64  bool
	FixedString     string
	HasFixedString  bool
	AllowedUInt64   []uint64
	ApplicationData bool
	NonZeroData     bool
}

type NativeXPCEnvelopeSpec struct {
	Method                  string
	Direction               string
	MessageTag              NativeXPCMessageTag
	ProtocolVersion         uint64
	MethodVersion           uint64
	ExactKeyCount           uint64
	RequiredKeyCount        uint64
	OptionalKeyCount        uint64
	ClosedMap               bool
	ApplicationDataMaxBytes uint64
	Fields                  []NativeXPCFieldSpec
}

type NativeXPCMethodBinding struct {
	Method                    string
	EntryPoint                string
	Service                   string
	ExpectedRole              CallerRole
	ExpectedSigningIdentifier string
	Audience                  string
	Purpose                   string
	MessageTag                NativeXPCMessageTag
	MethodVersion             uint64
	DeadlineMilliseconds      uint64
}

type NativeXPCErrorMapping struct {
	Classification Classification
	StatusTag      NativeXPCStatusTag
}

type NativeXPCRefusalMapping struct {
	ReasonTag NativeXPCReasonTag
	StatusTag NativeXPCStatusTag
}

func NativeXPCErrorMappings() []NativeXPCErrorMapping {
	return []NativeXPCErrorMapping{
		{Malformed, NativeXPCStatusMalformed},
		{Unsupported, NativeXPCStatusUnsupported},
		{Schema, NativeXPCStatusSchema},
		{Semantic, NativeXPCStatusSemantic},
		{Domain, NativeXPCStatusDomain},
		{Binding, NativeXPCStatusBinding},
		{Authentication, NativeXPCStatusAuthentication},
		{Stale, NativeXPCStatusStale},
		{Replay, NativeXPCStatusReplay},
		{Capacity, NativeXPCStatusCapacity},
		{TrustState, NativeXPCStatusTrustState},
		{LocalFailure, NativeXPCStatusLocalFailure},
		{RecoveryRequired, NativeXPCStatusRecoveryRequired},
	}
}

// NativeXPCStructuralRefusalMappings fixes the outer-dictionary error result.
// NativeXPCReasonCoreRefusal instead preserves the independently classified
// core status, and NativeXPCReasonLocalIntegrityFault terminates without a
// reply; neither has one default status in this table.
func NativeXPCStructuralRefusalMappings() []NativeXPCRefusalMapping {
	return []NativeXPCRefusalMapping{
		{NativeXPCReasonKeySet, NativeXPCStatusMalformed},
		{NativeXPCReasonValueType, NativeXPCStatusMalformed},
		{NativeXPCReasonDataWidth, NativeXPCStatusSchema},
		{NativeXPCReasonDataCap, NativeXPCStatusMalformed},
		{NativeXPCReasonZeroIdentifier, NativeXPCStatusSchema},
		{NativeXPCReasonEpochSequence, NativeXPCStatusSchema},
		{NativeXPCReasonProtocolVersion, NativeXPCStatusUnsupported},
		{NativeXPCReasonMethodVersion, NativeXPCStatusUnsupported},
		{NativeXPCReasonMessageTag, NativeXPCStatusUnsupported},
		{NativeXPCReasonMethodBinding, NativeXPCStatusAuthentication},
		{NativeXPCReasonCurrentState, NativeXPCStatusBinding},
		{NativeXPCReasonCapacity, NativeXPCStatusCapacity},
	}
}

func ExpectedNativeXPCMethodBindings() []NativeXPCMethodBinding {
	return []NativeXPCMethodBinding{
		{
			Method: SubmitMainMJSV0Method, EntryPoint: SubmitMainMJSV0Method, Service: SubmitMainMJSV0Service,
			ExpectedRole: InternalAlphaCLI, ExpectedSigningIdentifier: InternalAlphaCLISigningID,
			Audience: SubmitMainMJSV0Audience, Purpose: SubmitMainMJSV0Purpose,
			MessageTag: NativeXPCMessageTagSubmitMainMJSV0, MethodVersion: SubmitMainMJSV0MethodVersion,
			DeadlineMilliseconds: SubmitMainMJSV0DeadlineMilliseconds,
		},
		{
			Method: RegisterPlanV0Method, EntryPoint: RegisterPlanV0Method, Service: DaemonServiceV0, ExpectedRole: Daemon,
			Audience: PassiveIPCAudience, Purpose: RegisterPlanV0Purpose,
			MessageTag: NativeXPCMessageTagRegisterPlanV0, MethodVersion: RegisterPlanV0MethodVersion,
			DeadlineMilliseconds: RegisterPlanV0DeadlineMilliseconds,
		},
		{
			Method: GetRegisteredPlanV0Method, EntryPoint: GetRegisteredPlanV0Method, Service: BrokerServiceV0, ExpectedRole: Broker,
			Audience: PassiveIPCAudience, Purpose: GetRegisteredPlanV0Purpose,
			MessageTag: NativeXPCMessageTagGetRegisteredPlanV0, MethodVersion: GetRegisteredPlanV0MethodVersion,
			DeadlineMilliseconds: GetRegisteredPlanV0DeadlineMilliseconds,
		},
		{
			Method: SubmitApprovalV0Method, EntryPoint: SubmitApprovalV0Method, Service: BrokerServiceV0, ExpectedRole: Broker,
			Audience: PassiveIPCAudience, Purpose: SubmitApprovalV0Purpose,
			MessageTag: NativeXPCMessageTagSubmitApprovalV0, MethodVersion: SubmitApprovalV0MethodVersion,
			DeadlineMilliseconds: SubmitApprovalV0DeadlineMilliseconds,
		},
		{
			Method: RequestAttemptV0Method, EntryPoint: RequestAttemptV0Method, Service: DaemonServiceV0, ExpectedRole: Daemon,
			Audience: PassiveIPCAudience, Purpose: RequestAttemptV0Purpose,
			MessageTag: NativeXPCMessageTagRequestAttemptV0, MethodVersion: RequestAttemptV0MethodVersion,
			DeadlineMilliseconds: RequestAttemptV0DeadlineMilliseconds,
		},
	}
}

func ExpectedNativeXPCSubmitMainMJSV0Request() NativeXPCEnvelopeSpec {
	return nativeXPCRequestSpec(
		SubmitMainMJSV0Method,
		NativeXPCMessageTagSubmitMainMJSV0,
		SubmitMainMJSV0MethodVersion,
		SubmitMainMJSV0Audience,
		SubmitMainMJSV0Purpose,
		SubmitMainMJSV0RequestDataMaxBytes,
		NativeXPCSubmitRequestKeyCount,
		[]NativeXPCFieldSpec{nativeXPCDataField(NativeXPCJobProposalKey, 1, JobProposalV0RawMaxBytes, true, false)},
	)
}

func ExpectedNativeXPCRegisterPlanV0Request() NativeXPCEnvelopeSpec {
	return nativeXPCRequestSpec(
		RegisterPlanV0Method,
		NativeXPCMessageTagRegisterPlanV0,
		RegisterPlanV0MethodVersion,
		PassiveIPCAudience,
		RegisterPlanV0Purpose,
		RegisterPlanV0RequestDataMaxBytes,
		NativeXPCRegisterRequestKeyCount,
		[]NativeXPCFieldSpec{
			nativeXPCDataField(NativeXPCExecutionPlanKey, 1, 65_536, true, false),
			nativeXPCDataField(NativeXPCRoleBindingsKey, RoleBindingRecordBytes, RoleBindingRecordBytes, true, false),
			nativeXPCDataField(NativeXPCSourceManifestKey, 87, 95, true, false),
			nativeXPCDataField(NativeXPCSourceKey, 0, 262_144, true, false),
		},
	)
}

func ExpectedNativeXPCGetRegisteredPlanV0Request() NativeXPCEnvelopeSpec {
	return nativeXPCRequestSpec(
		GetRegisteredPlanV0Method,
		NativeXPCMessageTagGetRegisteredPlanV0,
		GetRegisteredPlanV0MethodVersion,
		PassiveIPCAudience,
		GetRegisteredPlanV0Purpose,
		GetRegisteredPlanV0RequestDataMaxBytes,
		NativeXPCGetRequestKeyCount,
		[]NativeXPCFieldSpec{nativeXPCDataField(NativeXPCRegistrationIDKey, 16, 16, true, true)},
	)
}

// ExpectedNativeXPCSubmitApprovalV0Request returns the passive request dictionary specification.
func ExpectedNativeXPCSubmitApprovalV0Request() NativeXPCEnvelopeSpec {
	return nativeXPCRequestSpec(
		SubmitApprovalV0Method,
		NativeXPCMessageTagSubmitApprovalV0,
		SubmitApprovalV0MethodVersion,
		PassiveIPCAudience,
		SubmitApprovalV0Purpose,
		SubmitApprovalV0RequestDataMaxBytes,
		NativeXPCSubmitApprovalRequestKeyCount,
		[]NativeXPCFieldSpec{
			nativeXPCDataField(NativeXPCRegistrationIDKey, 16, 16, true, true),
			nativeXPCDataField(NativeXPCApprovalEnvelopeKey, 1, 512, true, false),
		},
	)
}

// ExpectedNativeXPCRequestAttemptV0Request returns the passive request dictionary specification.
func ExpectedNativeXPCRequestAttemptV0Request() NativeXPCEnvelopeSpec {
	return nativeXPCRequestSpec(
		RequestAttemptV0Method,
		NativeXPCMessageTagRequestAttemptV0,
		RequestAttemptV0MethodVersion,
		PassiveIPCAudience,
		RequestAttemptV0Purpose,
		RequestAttemptV0RequestDataMaxBytes,
		NativeXPCRequestAttemptRequestKeyCount,
		[]NativeXPCFieldSpec{
			nativeXPCDataField(NativeXPCRegistrationIDKey, 16, 16, true, true),
			nativeXPCDataField(NativeXPCApprovalIDKey, 16, 16, true, true),
		},
	)
}

func ExpectedNativeXPCSubmitMainMJSV0Reply() NativeXPCEnvelopeSpec {
	return nativeXPCSuccessReplySpec(
		SubmitMainMJSV0Method,
		NativeXPCMessageTagSubmitMainMJSV0,
		SubmitMainMJSV0MethodVersion,
		SubmitMainMJSV0ReplyDataMaxBytes,
		NativeXPCSubmitReplyKeyCount,
		[]NativeXPCFieldSpec{nativeXPCDataField(NativeXPCRegistrationIDKey, 16, 16, true, true)},
	)
}

func ExpectedNativeXPCRegisterPlanV0Reply() NativeXPCEnvelopeSpec {
	return nativeXPCSuccessReplySpec(
		RegisterPlanV0Method,
		NativeXPCMessageTagRegisterPlanV0,
		RegisterPlanV0MethodVersion,
		RegisterPlanV0ReplyDataMaxBytes,
		NativeXPCRegisterReplyKeyCount,
		[]NativeXPCFieldSpec{nativeXPCDataField(NativeXPCPlanRegistrationKey, 1, 4_096, true, false)},
	)
}

func ExpectedNativeXPCGetRegisteredPlanV0Reply() NativeXPCEnvelopeSpec {
	return nativeXPCSuccessReplySpec(
		GetRegisteredPlanV0Method,
		NativeXPCMessageTagGetRegisteredPlanV0,
		GetRegisteredPlanV0MethodVersion,
		GetRegisteredPlanV0ReplyDataMaxBytes,
		NativeXPCGetReplyKeyCount,
		[]NativeXPCFieldSpec{
			nativeXPCDataField(NativeXPCExecutionPlanKey, 1, 65_536, true, false),
			nativeXPCDataField(NativeXPCRoleBindingsKey, RoleBindingRecordBytes, RoleBindingRecordBytes, true, false),
			nativeXPCDataField(NativeXPCPlanRegistrationKey, 1, 4_096, true, false),
			nativeXPCDataField(NativeXPCSourceManifestKey, 87, 95, true, false),
			nativeXPCDataField(NativeXPCSourceKey, 0, 262_144, true, false),
		},
	)
}

// ExpectedNativeXPCSubmitApprovalV0Reply returns the passive success-reply dictionary specification.
func ExpectedNativeXPCSubmitApprovalV0Reply() NativeXPCEnvelopeSpec {
	return nativeXPCSuccessReplySpec(
		SubmitApprovalV0Method,
		NativeXPCMessageTagSubmitApprovalV0,
		SubmitApprovalV0MethodVersion,
		SubmitApprovalV0ReplyDataMaxBytes,
		NativeXPCSubmitApprovalReplyKeyCount,
		[]NativeXPCFieldSpec{
			nativeXPCDataField(NativeXPCApprovalIDKey, 16, 16, true, true),
			nativeXPCAllowedUInt64Field(NativeXPCApprovalStateKey, []uint64{
				uint64(NativeXPCApprovalStateUsable),
				uint64(NativeXPCApprovalStateConsumed),
				uint64(NativeXPCApprovalStateInvalidated),
			}),
		},
	)
}

// ExpectedNativeXPCRequestAttemptV0Reply returns the passive success-reply dictionary specification.
func ExpectedNativeXPCRequestAttemptV0Reply() NativeXPCEnvelopeSpec {
	return nativeXPCSuccessReplySpec(
		RequestAttemptV0Method,
		NativeXPCMessageTagRequestAttemptV0,
		RequestAttemptV0MethodVersion,
		RequestAttemptV0ReplyDataMaxBytes,
		NativeXPCRequestAttemptReplyKeyCount,
		[]NativeXPCFieldSpec{
			nativeXPCDataField(NativeXPCAttemptIDKey, 16, 16, true, true),
			nativeXPCAllowedUInt64Field(NativeXPCAttemptStateKey, []uint64{uint64(NativeXPCAttemptStateCreated)}),
		},
	)
}

func ExpectedNativeXPCSubmitMainMJSV0RefusalReply() NativeXPCEnvelopeSpec {
	return nativeXPCRefusalReplySpec(SubmitMainMJSV0Method, NativeXPCMessageTagSubmitMainMJSV0, SubmitMainMJSV0MethodVersion)
}

func ExpectedNativeXPCRegisterPlanV0RefusalReply() NativeXPCEnvelopeSpec {
	return nativeXPCRefusalReplySpec(RegisterPlanV0Method, NativeXPCMessageTagRegisterPlanV0, RegisterPlanV0MethodVersion)
}

func ExpectedNativeXPCGetRegisteredPlanV0RefusalReply() NativeXPCEnvelopeSpec {
	return nativeXPCRefusalReplySpec(GetRegisteredPlanV0Method, NativeXPCMessageTagGetRegisteredPlanV0, GetRegisteredPlanV0MethodVersion)
}

// ExpectedNativeXPCSubmitApprovalV0RefusalReply returns the body-free refusal specification.
func ExpectedNativeXPCSubmitApprovalV0RefusalReply() NativeXPCEnvelopeSpec {
	return nativeXPCRefusalReplySpec(SubmitApprovalV0Method, NativeXPCMessageTagSubmitApprovalV0, SubmitApprovalV0MethodVersion)
}

// ExpectedNativeXPCRequestAttemptV0RefusalReply returns the body-free refusal specification.
func ExpectedNativeXPCRequestAttemptV0RefusalReply() NativeXPCEnvelopeSpec {
	return nativeXPCRefusalReplySpec(RequestAttemptV0Method, NativeXPCMessageTagRequestAttemptV0, RequestAttemptV0MethodVersion)
}

func nativeXPCRefusalReplySpec(method string, tag NativeXPCMessageTag, version uint64) NativeXPCEnvelopeSpec {
	return NativeXPCEnvelopeSpec{
		Method: method, Direction: "refusal-reply", MessageTag: tag,
		ProtocolVersion: PassiveIPCProtocolVersion, MethodVersion: version,
		ExactKeyCount: NativeXPCRefusalReplyKeyCount, RequiredKeyCount: NativeXPCRefusalReplyKeyCount, ClosedMap: true,
		Fields: nativeXPCReplyHeader(tag, version, false),
	}
}

func nativeXPCRequestSpec(
	method string,
	tag NativeXPCMessageTag,
	methodVersion uint64,
	audience string,
	purpose string,
	applicationDataMaxBytes uint64,
	exactKeyCount uint64,
	body []NativeXPCFieldSpec,
) NativeXPCEnvelopeSpec {
	fields := []NativeXPCFieldSpec{
		nativeXPCFixedUInt64Field(NativeXPCProtocolVersionKey, PassiveIPCProtocolVersion),
		nativeXPCFixedUInt64Field(NativeXPCMethodVersionKey, methodVersion),
		nativeXPCFixedUInt64Field(NativeXPCMessageTagKey, uint64(tag)),
		nativeXPCDataField(NativeXPCRequestIDKey, 16, 16, false, true),
		nativeXPCDataField(NativeXPCInstallationIDKey, 16, 16, false, false),
		{Key: NativeXPCEpochSequenceKey, ValueType: NativeXPCTypeUInt64, Required: true},
		nativeXPCDataField(NativeXPCEpochDigestKey, 32, 32, false, false),
		nativeXPCFixedStringField(NativeXPCAudienceKey, audience),
		nativeXPCFixedStringField(NativeXPCPurposeKey, purpose),
	}
	fields = append(fields, body...)
	return NativeXPCEnvelopeSpec{
		Method: method, Direction: "request", MessageTag: tag,
		ProtocolVersion: PassiveIPCProtocolVersion, MethodVersion: methodVersion,
		ExactKeyCount: exactKeyCount, RequiredKeyCount: exactKeyCount, ClosedMap: true, ApplicationDataMaxBytes: applicationDataMaxBytes,
		Fields: fields,
	}
}

func nativeXPCSuccessReplySpec(
	method string,
	tag NativeXPCMessageTag,
	methodVersion uint64,
	applicationDataMaxBytes uint64,
	exactKeyCount uint64,
	body []NativeXPCFieldSpec,
) NativeXPCEnvelopeSpec {
	fields := nativeXPCReplyHeader(tag, methodVersion, true)
	fields = append(fields, body...)
	return NativeXPCEnvelopeSpec{
		Method: method, Direction: "success-reply", MessageTag: tag,
		ProtocolVersion: PassiveIPCProtocolVersion, MethodVersion: methodVersion,
		ExactKeyCount: exactKeyCount, RequiredKeyCount: exactKeyCount, ClosedMap: true, ApplicationDataMaxBytes: applicationDataMaxBytes,
		Fields: fields,
	}
}

func nativeXPCReplyHeader(tag NativeXPCMessageTag, methodVersion uint64, success bool) []NativeXPCFieldSpec {
	status := uint64(NativeXPCStatusOK)
	reason := uint64(NativeXPCReasonNone)
	return []NativeXPCFieldSpec{
		nativeXPCFixedUInt64Field(NativeXPCProtocolVersionKey, PassiveIPCProtocolVersion),
		nativeXPCFixedUInt64Field(NativeXPCMethodVersionKey, methodVersion),
		nativeXPCFixedUInt64Field(NativeXPCMessageTagKey, uint64(tag)),
		nativeXPCDataField(NativeXPCRequestIDKey, 16, 16, false, true),
		{Key: NativeXPCStatusKey, ValueType: NativeXPCTypeUInt64, Required: true, FixedUInt64: status, HasFixedUInt64: success},
		{Key: NativeXPCReasonKey, ValueType: NativeXPCTypeUInt64, Required: true, FixedUInt64: reason, HasFixedUInt64: success},
	}
}

func nativeXPCDataField(key string, minimum, maximum uint64, applicationData, nonzero bool) NativeXPCFieldSpec {
	return NativeXPCFieldSpec{Key: key, ValueType: NativeXPCTypeData, Required: true, MinDataBytes: minimum, MaxDataBytes: maximum, ApplicationData: applicationData, NonZeroData: nonzero}
}

func nativeXPCAllowedUInt64Field(key string, allowed []uint64) NativeXPCFieldSpec {
	return NativeXPCFieldSpec{Key: key, ValueType: NativeXPCTypeUInt64, Required: true, AllowedUInt64: append([]uint64(nil), allowed...)}
}

func nativeXPCFixedUInt64Field(key string, value uint64) NativeXPCFieldSpec {
	return NativeXPCFieldSpec{Key: key, ValueType: NativeXPCTypeUInt64, Required: true, FixedUInt64: value, HasFixedUInt64: true}
}

func nativeXPCFixedStringField(key, value string) NativeXPCFieldSpec {
	return NativeXPCFieldSpec{Key: key, ValueType: NativeXPCTypeString, Required: true, FixedString: value, HasFixedString: true}
}
