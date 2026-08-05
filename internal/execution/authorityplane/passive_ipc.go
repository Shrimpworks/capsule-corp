package authorityplane

import (
	"bytes"
	"context"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

// This file freezes only the passive application contract between a future
// native front end and the existing unwired Go facade. It defines no XPC key
// spelling, numeric message tag, transport encoding, endpoint, or peer-auth
// evidence. Those values remain owned by a later native-transport slice.
const (
	PassiveIPCProtocolVersion = uint64(0)
	// #nosec G101 -- fixed local audience discriminator, not a credential.
	PassiveIPCAudience = "capsule.execution-supervisor.local.v0"

	DaemonServiceV0 = "com.capsulecorp.capsule.supervisor.daemon.v0"
	BrokerServiceV0 = "com.capsulecorp.capsule.supervisor.broker.v0"

	RegisterPlanV0Method             = "RegisterPlanV0"
	GetRegisteredPlanV0Method        = "GetRegisteredPlanV0"
	RegisterPlanV0MethodVersion      = uint64(0)
	GetRegisteredPlanV0MethodVersion = uint64(0)

	RegisterPlanV0RequestDataMaxBytes      = RegisterPlanV0MaxBytes
	RegisterPlanV0ReplyDataMaxBytes        = v0candidate.PlanRegistrationMaxCBORBytes
	GetRegisteredPlanV0RequestDataMaxBytes = 16
	GetRegisteredPlanV0ReplyDataMaxBytes   = GetRegisteredPlanV0MaxBytes

	RegisterPlanV0DeadlineMilliseconds      = uint64(5_000)
	GetRegisteredPlanV0DeadlineMilliseconds = uint64(2_000)

	RegisterPlanV0ResponseLossDisposition      = "committed-retry-creates-fresh-registration"
	GetRegisteredPlanV0ResponseLossDisposition = "repeatable-read-by-registration-id"
)

var (
	_ [328_337 - RegisterPlanV0RequestDataMaxBytes]byte
	_ [RegisterPlanV0RequestDataMaxBytes - 328_337]byte
	_ [4_096 - RegisterPlanV0ReplyDataMaxBytes]byte
	_ [RegisterPlanV0ReplyDataMaxBytes - 4_096]byte
	_ [16 - GetRegisteredPlanV0RequestDataMaxBytes]byte
	_ [GetRegisteredPlanV0RequestDataMaxBytes - 16]byte
	_ [332_433 - GetRegisteredPlanV0ReplyDataMaxBytes]byte
	_ [GetRegisteredPlanV0ReplyDataMaxBytes - 332_433]byte
)

type PassiveRequestID [16]byte

// field-authority-object: capsule.authenticated-local-ipc-register-plan-v0-method-record v0
// RegisterPlanV0MethodRecord binds the method and its v0 role record as one
// passive contract. The future native layer must derive this record from its
// role-specific service and method dispatch, never from application body data.
type RegisterPlanV0MethodRecord struct {
	Method                   string
	MethodVersion            uint64
	RoleBindingRecordVersion uint64
	Service                  string
	ExpectedRole             CallerRole
	Audience                 string
	Purpose                  string
	RequestDataMaxBytes      uint64
	ReplyDataMaxBytes        uint64
	DeadlineMilliseconds     uint64
	ResponseLossDisposition  string
}

// field-authority-object: capsule.authenticated-local-ipc-get-registered-plan-v0-method-record v0
type GetRegisteredPlanV0MethodRecord struct {
	Method                   string
	MethodVersion            uint64
	RoleBindingRecordVersion uint64
	Service                  string
	ExpectedRole             CallerRole
	Audience                 string
	Purpose                  string
	RequestDataMaxBytes      uint64
	ReplyDataMaxBytes        uint64
	DeadlineMilliseconds     uint64
	ResponseLossDisposition  string
}

func ExpectedRegisterPlanV0MethodRecord() RegisterPlanV0MethodRecord {
	return RegisterPlanV0MethodRecord{
		Method: RegisterPlanV0Method, MethodVersion: RegisterPlanV0MethodVersion,
		RoleBindingRecordVersion: uint64(RoleBindingRecordVersion),
		Service:                  DaemonServiceV0, ExpectedRole: Daemon, Audience: PassiveIPCAudience,
		Purpose: RegisterPlanV0Purpose, RequestDataMaxBytes: RegisterPlanV0RequestDataMaxBytes,
		ReplyDataMaxBytes:       RegisterPlanV0ReplyDataMaxBytes,
		DeadlineMilliseconds:    RegisterPlanV0DeadlineMilliseconds,
		ResponseLossDisposition: RegisterPlanV0ResponseLossDisposition,
	}
}

func ExpectedGetRegisteredPlanV0MethodRecord() GetRegisteredPlanV0MethodRecord {
	return GetRegisteredPlanV0MethodRecord{
		Method: GetRegisteredPlanV0Method, MethodVersion: GetRegisteredPlanV0MethodVersion,
		RoleBindingRecordVersion: uint64(RoleBindingRecordVersion),
		Service:                  BrokerServiceV0, ExpectedRole: Broker, Audience: PassiveIPCAudience,
		Purpose: GetRegisteredPlanV0Purpose, RequestDataMaxBytes: GetRegisteredPlanV0RequestDataMaxBytes,
		ReplyDataMaxBytes:       GetRegisteredPlanV0ReplyDataMaxBytes,
		DeadlineMilliseconds:    GetRegisteredPlanV0DeadlineMilliseconds,
		ResponseLossDisposition: GetRegisteredPlanV0ResponseLossDisposition,
	}
}

// field-authority-object: capsule.authenticated-local-ipc-register-plan-v0-request v0
type PassiveRegisterPlanV0Request struct {
	ProtocolVersion     uint64
	RequestID           PassiveRequestID
	InstallationID      v0candidate.InstallationID
	EpochSequence       v0candidate.UInt53
	EpochDigest         v0candidate.TrustEpochDigest
	exactPlanBytes      []byte
	roleBindingBytes    []byte
	sourceManifestBytes []byte
	sourceBytes         []byte
}

func NewPassiveRegisterPlanV0Request(
	protocolVersion uint64,
	requestID PassiveRequestID,
	installationID v0candidate.InstallationID,
	epochSequence v0candidate.UInt53,
	epochDigest v0candidate.TrustEpochDigest,
	plan, bindings, manifest, source []byte,
) (PassiveRegisterPlanV0Request, error) {
	body, err := NewRegisterPlanV0Request(plan, bindings, manifest, source)
	if err != nil {
		return PassiveRegisterPlanV0Request{}, err
	}
	return PassiveRegisterPlanV0Request{
		ProtocolVersion: protocolVersion, RequestID: requestID,
		InstallationID: installationID, EpochSequence: epochSequence, EpochDigest: epochDigest,
		exactPlanBytes: body.ExactPlanBytes(), roleBindingBytes: body.NominalRoleBindingBytes(),
		sourceManifestBytes: body.SourceManifestBytes(), sourceBytes: body.SourceBytes(),
	}, nil
}

func (r PassiveRegisterPlanV0Request) ExactPlanBytes() []byte { return bytes.Clone(r.exactPlanBytes) }
func (r PassiveRegisterPlanV0Request) RoleBindingBytes() []byte {
	return bytes.Clone(r.roleBindingBytes)
}
func (r PassiveRegisterPlanV0Request) SourceManifestBytes() []byte {
	return bytes.Clone(r.sourceManifestBytes)
}
func (r PassiveRegisterPlanV0Request) SourceBytes() []byte { return bytes.Clone(r.sourceBytes) }
func (r PassiveRegisterPlanV0Request) ApplicationDataBytes() int {
	return len(r.exactPlanBytes) + len(r.roleBindingBytes) + len(r.sourceManifestBytes) + len(r.sourceBytes)
}

func (r PassiveRegisterPlanV0Request) coreRequest() (RegisterPlanV0Request, error) {
	return NewRegisterPlanV0Request(r.exactPlanBytes, r.roleBindingBytes, r.sourceManifestBytes, r.sourceBytes)
}

// field-authority-object: capsule.authenticated-local-ipc-get-registered-plan-v0-request v0
type PassiveGetRegisteredPlanV0Request struct {
	ProtocolVersion uint64
	RequestID       PassiveRequestID
	InstallationID  v0candidate.InstallationID
	EpochSequence   v0candidate.UInt53
	EpochDigest     v0candidate.TrustEpochDigest
	RegistrationID  v0candidate.RegistrationID
}

func (r PassiveGetRegisteredPlanV0Request) ApplicationDataBytes() int { return len(r.RegistrationID) }

// field-authority-object: capsule.authenticated-local-ipc-register-plan-v0-reply v0
type PassiveRegisterPlanV0Reply struct {
	RequestID             PassiveRequestID
	planRegistrationBytes []byte
}

func (r PassiveRegisterPlanV0Reply) PlanRegistrationBytes() []byte {
	return bytes.Clone(r.planRegistrationBytes)
}
func (r PassiveRegisterPlanV0Reply) ApplicationDataBytes() int {
	return len(r.planRegistrationBytes)
}

// field-authority-object: capsule.authenticated-local-ipc-get-registered-plan-v0-reply v0
type PassiveGetRegisteredPlanV0Reply struct {
	RequestID             PassiveRequestID
	exactPlanBytes        []byte
	roleBindingBytes      []byte
	planRegistrationBytes []byte
	sourceManifestBytes   []byte
	sourceBytes           []byte
}

func (r PassiveGetRegisteredPlanV0Reply) ExactPlanBytes() []byte {
	return bytes.Clone(r.exactPlanBytes)
}
func (r PassiveGetRegisteredPlanV0Reply) RoleBindingBytes() []byte {
	return bytes.Clone(r.roleBindingBytes)
}
func (r PassiveGetRegisteredPlanV0Reply) PlanRegistrationBytes() []byte {
	return bytes.Clone(r.planRegistrationBytes)
}
func (r PassiveGetRegisteredPlanV0Reply) SourceManifestBytes() []byte {
	return bytes.Clone(r.sourceManifestBytes)
}
func (r PassiveGetRegisteredPlanV0Reply) SourceBytes() []byte { return bytes.Clone(r.sourceBytes) }
func (r PassiveGetRegisteredPlanV0Reply) ApplicationDataBytes() int {
	return len(r.exactPlanBytes) + len(r.roleBindingBytes) + len(r.planRegistrationBytes) + len(r.sourceManifestBytes) + len(r.sourceBytes)
}

// field-authority-object: capsule.authenticated-local-ipc-refusal-reply v0
// PassiveRefusalReply is a logical fixed-code oracle, not an XPC or C ABI
// serialization. Numeric status values remain intentionally unselected.
type PassiveRefusalReply struct {
	RequestID      PassiveRequestID
	Classification Classification
	Reason         string
}

type passiveAuthorityCore interface {
	RegisterPlanV0(context.Context, CallContext, RegisterPlanV0Request) ([]byte, error)
	GetRegisteredPlanV0(context.Context, CallContext, GetRegisteredPlanV0Request) (GetRegisteredPlanV0Reply, error)
}

var _ passiveAuthorityCore = (*Facade)(nil)

// passiveIPCBoundary exists only for in-package conformance. Keeping it
// unexported prevents passive fixture admission from becoming a product
// authentication bypass or callable transport surface.
type passiveIPCBoundary struct {
	core       passiveAuthorityCore
	supervisor SupervisorContext
}

func newPassiveIPCBoundary(facade passiveAuthorityCore, supervisor SupervisorContext) (*passiveIPCBoundary, error) {
	if facade == nil {
		return nil, refused(LocalFailure, "passive-ipc-facade")
	}
	return &passiveIPCBoundary{core: facade, supervisor: supervisor}, nil
}

func (b *passiveIPCBoundary) registerPlanV0(
	ctx context.Context,
	method RegisterPlanV0MethodRecord,
	request PassiveRegisterPlanV0Request,
) (PassiveRegisterPlanV0Reply, error) {
	if err := validateRegisterMethod(method); err != nil {
		return PassiveRegisterPlanV0Reply{}, err
	}
	if err := b.validateHeader(request.ProtocolVersion, request.RequestID, request.InstallationID, request.EpochSequence, request.EpochDigest); err != nil {
		return PassiveRegisterPlanV0Reply{}, err
	}
	body, err := request.coreRequest()
	if err != nil {
		return PassiveRegisterPlanV0Reply{}, err
	}
	// Authenticated is a synthetic facade precondition in this unexported
	// conformance adapter. It is not evidence of a peer-authenticated channel.
	registration, err := b.core.RegisterPlanV0(ctx, CallContext{Authenticated: true, Role: Daemon, Purpose: RegisterPlanV0Purpose}, body)
	if err != nil {
		return PassiveRegisterPlanV0Reply{}, err
	}
	if len(registration) == 0 {
		panic("register-plan-reply-shape")
	}
	if len(registration) > RegisterPlanV0ReplyDataMaxBytes {
		panic("register-plan-reply-cap")
	}
	return PassiveRegisterPlanV0Reply{RequestID: request.RequestID, planRegistrationBytes: bytes.Clone(registration)}, nil
}

func (b *passiveIPCBoundary) getRegisteredPlanV0(
	ctx context.Context,
	method GetRegisteredPlanV0MethodRecord,
	request PassiveGetRegisteredPlanV0Request,
) (PassiveGetRegisteredPlanV0Reply, error) {
	if err := validateGetMethod(method); err != nil {
		return PassiveGetRegisteredPlanV0Reply{}, err
	}
	if err := b.validateHeader(request.ProtocolVersion, request.RequestID, request.InstallationID, request.EpochSequence, request.EpochDigest); err != nil {
		return PassiveGetRegisteredPlanV0Reply{}, err
	}
	if request.RegistrationID == (v0candidate.RegistrationID{}) {
		return PassiveGetRegisteredPlanV0Reply{}, refused(Schema, "registration-id")
	}
	// Authenticated is a synthetic facade precondition in this unexported
	// conformance adapter. It is not evidence of a peer-authenticated channel.
	reply, err := b.core.GetRegisteredPlanV0(ctx, CallContext{Authenticated: true, Role: Broker, Purpose: GetRegisteredPlanV0Purpose}, GetRegisteredPlanV0Request{RegistrationID: request.RegistrationID})
	if err != nil {
		return PassiveGetRegisteredPlanV0Reply{}, err
	}
	return newPassiveGetRegisteredPlanV0Reply(request.RequestID, reply), nil
}

func (b *passiveIPCBoundary) validateHeader(
	version uint64,
	requestID PassiveRequestID,
	installationID v0candidate.InstallationID,
	epochSequence v0candidate.UInt53,
	epochDigest v0candidate.TrustEpochDigest,
) error {
	if version != PassiveIPCProtocolVersion {
		return refused(Unsupported, "ipc-protocol-version")
	}
	if requestID == (PassiveRequestID{}) {
		return refused(Schema, "ipc-request-id")
	}
	if uint64(epochSequence) > v0candidate.MaxSafeInteger {
		return refused(Schema, "ipc-epoch-sequence")
	}
	if b == nil || b.core == nil {
		return refused(LocalFailure, "passive-ipc-core")
	}
	if installationID != b.supervisor.InstallationID || epochSequence != b.supervisor.EpochSequence || epochDigest != b.supervisor.EpochDigest {
		return refused(Binding, "ipc-current-supervisor-state")
	}
	return nil
}

func validateRegisterMethod(received RegisterPlanV0MethodRecord) error {
	expected := ExpectedRegisterPlanV0MethodRecord()
	if received.MethodVersion != expected.MethodVersion || received.RoleBindingRecordVersion != expected.RoleBindingRecordVersion {
		return refused(Unsupported, "register-plan-method-record-version")
	}
	if received.Method != expected.Method {
		return refused(Unsupported, "register-plan-method")
	}
	if received != expected {
		return refused(Authentication, "register-plan-method-binding")
	}
	return nil
}

func validateGetMethod(received GetRegisteredPlanV0MethodRecord) error {
	expected := ExpectedGetRegisteredPlanV0MethodRecord()
	if received.MethodVersion != expected.MethodVersion || received.RoleBindingRecordVersion != expected.RoleBindingRecordVersion {
		return refused(Unsupported, "get-registered-plan-method-record-version")
	}
	if received.Method != expected.Method {
		return refused(Unsupported, "get-registered-plan-method")
	}
	if received != expected {
		return refused(Authentication, "get-registered-plan-method-binding")
	}
	return nil
}

func newPassiveGetRegisteredPlanV0Reply(requestID PassiveRequestID, reply GetRegisteredPlanV0Reply) PassiveGetRegisteredPlanV0Reply {
	plan := reply.ExactPlanBytes()
	bindings := reply.ResolvedRoleBindingBytes()
	registration := reply.PlanRegistrationBytes()
	manifest := reply.SourceManifestBytes()
	source := reply.SourceBytes()
	if len(plan) == 0 || len(plan) > v0candidate.ExecutionPlanMaxCBORBytes ||
		len(bindings) != RoleBindingRecordBytes ||
		len(registration) == 0 || len(registration) > v0candidate.PlanRegistrationMaxCBORBytes ||
		len(manifest) < v0candidate.SourceManifestMinCBORBytes || len(manifest) > v0candidate.SourceManifestMaxCBORBytes ||
		len(source) > v0candidate.MJSMainSourceMaxBytes {
		panic("get-registered-plan-reply-shape")
	}
	if len(plan)+len(bindings)+len(registration)+len(manifest)+len(source) > GetRegisteredPlanV0ReplyDataMaxBytes {
		panic("get-registered-plan-reply-cap")
	}
	return PassiveGetRegisteredPlanV0Reply{
		RequestID:      requestID,
		exactPlanBytes: bytes.Clone(plan), roleBindingBytes: bytes.Clone(bindings),
		planRegistrationBytes: bytes.Clone(registration), sourceManifestBytes: bytes.Clone(manifest),
		sourceBytes: bytes.Clone(source),
	}
}

func passiveRefusal(requestID PassiveRequestID, err error) (PassiveRefusalReply, bool) {
	classification, ok := ErrorClassification(err)
	if !ok {
		return PassiveRefusalReply{}, false
	}
	reason := "core-refusal"
	if code, recognized := ErrorCode(err); recognized {
		reason = code
	}
	return PassiveRefusalReply{RequestID: requestID, Classification: classification, Reason: reason}, true
}
