package authorityplane

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

type passiveFixtureReference struct {
	Path       string `json:"path"`
	ByteLength int    `json:"byteLength"`
	SHA256     string `json:"sha256"`
}

type passiveMethodFixture struct {
	ObjectType                            string     `json:"objectType"`
	ObjectVersion                         uint64     `json:"objectVersion"`
	Method                                string     `json:"method"`
	MethodVersion                         uint64     `json:"methodVersion"`
	RoleBindingRecordVersion              uint64     `json:"roleBindingRecordVersion"`
	Service                               string     `json:"service"`
	ExpectedRole                          CallerRole `json:"expectedRole"`
	Audience                              string     `json:"audience"`
	Purpose                               string     `json:"purpose"`
	RequestDataMaxBytes                   uint64     `json:"requestDataMaxBytes"`
	ReplyDataMaxBytes                     uint64     `json:"replyDataMaxBytes"`
	DeadlineMilliseconds                  uint64     `json:"deadlineMilliseconds"`
	ResponseLossDisposition               string     `json:"responseLossDisposition"`
	PeerRequirementBeforeDeliveryRequired bool       `json:"peerRequirementBeforeDeliveryRequired"`
	MethodAuthorityDisposition            string     `json:"methodAuthorityDisposition"`
	RequestIDAuthorityDisposition         string     `json:"requestIdAuthorityDisposition"`
	ServiceMaxConnections                 uint64     `json:"serviceMaxConnections"`
	ConnectionMaxInFlight                 uint64     `json:"connectionMaxInFlight"`
	ProcessMaxAdmittedRequests            uint64     `json:"processMaxAdmittedRequests"`
	ApplicationQueueCapacity              uint64     `json:"applicationQueueCapacity"`
	InFlightRequestDataMaxBytes           uint64     `json:"inFlightRequestDataMaxBytes"`
}

type passiveRegisterRequestFixture struct {
	ObjectType           string                  `json:"objectType"`
	ObjectVersion        uint64                  `json:"objectVersion"`
	FixtureSerialization string                  `json:"fixtureSerialization"`
	ProtocolVersion      uint64                  `json:"protocolVersion"`
	RequestID            string                  `json:"requestId"`
	InstallationID       string                  `json:"installationId"`
	EpochSequence        uint64                  `json:"epochSequence"`
	EpochDigest          string                  `json:"epochDigest"`
	MethodRecord         passiveFixtureReference `json:"methodRecord"`
	Body                 struct {
		ExactPlanBytes      passiveFixtureReference `json:"exactPlanBytes"`
		RoleBindingBytes    passiveFixtureReference `json:"roleBindingBytes"`
		SourceManifestBytes passiveFixtureReference `json:"sourceManifestBytes"`
		SourceBytes         passiveFixtureReference `json:"sourceBytes"`
	} `json:"body"`
	ApplicationDataBytes int `json:"applicationDataBytes"`
}

type passiveRegisterReplyFixture struct {
	ObjectType           string `json:"objectType"`
	ObjectVersion        uint64 `json:"objectVersion"`
	FixtureSerialization string `json:"fixtureSerialization"`
	RequestID            string `json:"requestId"`
	Body                 struct {
		PlanRegistrationBytes passiveFixtureReference `json:"planRegistrationBytes"`
	} `json:"body"`
	ApplicationDataBytes int `json:"applicationDataBytes"`
}

type passiveGetRequestFixture struct {
	ObjectType           string                  `json:"objectType"`
	ObjectVersion        uint64                  `json:"objectVersion"`
	FixtureSerialization string                  `json:"fixtureSerialization"`
	ProtocolVersion      uint64                  `json:"protocolVersion"`
	RequestID            string                  `json:"requestId"`
	InstallationID       string                  `json:"installationId"`
	EpochSequence        uint64                  `json:"epochSequence"`
	EpochDigest          string                  `json:"epochDigest"`
	MethodRecord         passiveFixtureReference `json:"methodRecord"`
	Body                 struct {
		RegistrationID string `json:"registrationId"`
	} `json:"body"`
	ApplicationDataBytes int `json:"applicationDataBytes"`
}

type passiveGetReplyFixture struct {
	ObjectType           string `json:"objectType"`
	ObjectVersion        uint64 `json:"objectVersion"`
	FixtureSerialization string `json:"fixtureSerialization"`
	RequestID            string `json:"requestId"`
	Body                 struct {
		ExactPlanBytes        passiveFixtureReference `json:"exactPlanBytes"`
		RoleBindingBytes      passiveFixtureReference `json:"roleBindingBytes"`
		PlanRegistrationBytes passiveFixtureReference `json:"planRegistrationBytes"`
		SourceManifestBytes   passiveFixtureReference `json:"sourceManifestBytes"`
		SourceBytes           passiveFixtureReference `json:"sourceBytes"`
	} `json:"body"`
	ApplicationDataBytes int `json:"applicationDataBytes"`
}

type passiveManifestFixture struct {
	KnownAnswers map[string]passiveFixtureReference `json:"knownAnswers"`
}

type countingPassiveCore struct {
	registerCalls int
	getCalls      int
	registerReply []byte
	getReply      GetRegisteredPlanV0Reply
	getError      error
}

func (c *countingPassiveCore) RegisterPlanV0(context.Context, CallContext, RegisterPlanV0Request) ([]byte, error) {
	c.registerCalls++
	if c.registerReply != nil {
		return bytes.Clone(c.registerReply), nil
	}
	return []byte{0xa0}, nil
}

func (c *countingPassiveCore) GetRegisteredPlanV0(context.Context, CallContext, GetRegisteredPlanV0Request) (GetRegisteredPlanV0Reply, error) {
	c.getCalls++
	if c.getError != nil {
		return GetRegisteredPlanV0Reply{}, c.getError
	}
	return c.getReply, nil
}

type sequenceIdentifiers struct {
	mu     sync.Mutex
	values []v0candidate.RegistrationID
}

func (s *sequenceIdentifiers) NewRegistrationID(context.Context) (v0candidate.RegistrationID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.values) == 0 {
		return v0candidate.RegistrationID{}, refused(LocalFailure, "identifier-fixture-exhausted")
	}
	value := s.values[0]
	s.values = s.values[1:]
	return value, nil
}

func TestPassiveMethodRecordsBindVersionsRolesAndCaps(t *testing.T) {
	register := ExpectedRegisterPlanV0MethodRecord()
	fetch := ExpectedGetRegisteredPlanV0MethodRecord()
	if register.MethodVersion != 0 || register.RoleBindingRecordVersion != 0 ||
		fetch.MethodVersion != 0 || fetch.RoleBindingRecordVersion != 0 {
		t.Fatal("method and role-binding records did not cut over atomically at v0")
	}
	if register.Service != DaemonServiceV0 || register.ExpectedRole != Daemon ||
		register.Purpose != RegisterPlanV0Purpose || register.Audience != PassiveIPCAudience {
		t.Fatal("RegisterPlanV0 method binding drift")
	}
	if fetch.Service != BrokerServiceV0 || fetch.ExpectedRole != Broker ||
		fetch.Purpose != GetRegisteredPlanV0Purpose || fetch.Audience != PassiveIPCAudience {
		t.Fatal("GetRegisteredPlanV0 method binding drift")
	}
	if register.RequestDataMaxBytes != 328_337 || register.ReplyDataMaxBytes != 4_096 ||
		fetch.RequestDataMaxBytes != 16 || fetch.ReplyDataMaxBytes != 332_433 {
		t.Fatal("passive method cap drift")
	}
	if register.DeadlineMilliseconds != 5_000 || fetch.DeadlineMilliseconds != 2_000 {
		t.Fatal("passive method deadline drift")
	}
}

func TestPassiveCrossLanguageKnownAnswers(t *testing.T) {
	var manifest passiveManifestFixture
	if err := json.Unmarshal(readPassiveFixture(t, "manifest.json"), &manifest); err != nil {
		t.Fatal(err)
	}
	for path, known := range manifest.KnownAnswers {
		bytes := readPassiveFixture(t, path)
		assertPassiveReference(t, bytes, known)
	}

	var registerMethod passiveMethodFixture
	readPassiveJSON(t, "register-plan-v0.method.json", &registerMethod)
	expectedRegister := ExpectedRegisterPlanV0MethodRecord()
	if registerMethod.ObjectType != "capsule.authenticated-local-ipc-register-plan-v0-method-record" ||
		registerMethod.ObjectVersion != 0 ||
		(RegisterPlanV0MethodRecord{
			Method: registerMethod.Method, MethodVersion: registerMethod.MethodVersion,
			RoleBindingRecordVersion: registerMethod.RoleBindingRecordVersion,
			Service:                  registerMethod.Service, ExpectedRole: registerMethod.ExpectedRole,
			Audience: registerMethod.Audience, Purpose: registerMethod.Purpose,
			RequestDataMaxBytes:                   registerMethod.RequestDataMaxBytes,
			ReplyDataMaxBytes:                     registerMethod.ReplyDataMaxBytes,
			DeadlineMilliseconds:                  registerMethod.DeadlineMilliseconds,
			ResponseLossDisposition:               registerMethod.ResponseLossDisposition,
			PeerRequirementBeforeDeliveryRequired: registerMethod.PeerRequirementBeforeDeliveryRequired,
			MethodAuthorityDisposition:            registerMethod.MethodAuthorityDisposition,
			RequestIDAuthorityDisposition:         registerMethod.RequestIDAuthorityDisposition,
			ServiceMaxConnections:                 registerMethod.ServiceMaxConnections,
			ConnectionMaxInFlight:                 registerMethod.ConnectionMaxInFlight,
			ProcessMaxAdmittedRequests:            registerMethod.ProcessMaxAdmittedRequests,
			ApplicationQueueCapacity:              registerMethod.ApplicationQueueCapacity,
			InFlightRequestDataMaxBytes:           registerMethod.InFlightRequestDataMaxBytes,
		}) != expectedRegister {
		t.Fatal("Go RegisterPlanV0 method record disagrees with generated TypeScript fixture")
	}

	var getMethod passiveMethodFixture
	readPassiveJSON(t, "get-registered-plan-v0.method.json", &getMethod)
	expectedGet := ExpectedGetRegisteredPlanV0MethodRecord()
	if getMethod.ObjectType != "capsule.authenticated-local-ipc-get-registered-plan-v0-method-record" ||
		getMethod.ObjectVersion != 0 ||
		(GetRegisteredPlanV0MethodRecord{
			Method: getMethod.Method, MethodVersion: getMethod.MethodVersion,
			RoleBindingRecordVersion: getMethod.RoleBindingRecordVersion,
			Service:                  getMethod.Service, ExpectedRole: getMethod.ExpectedRole,
			Audience: getMethod.Audience, Purpose: getMethod.Purpose,
			RequestDataMaxBytes:                   getMethod.RequestDataMaxBytes,
			ReplyDataMaxBytes:                     getMethod.ReplyDataMaxBytes,
			DeadlineMilliseconds:                  getMethod.DeadlineMilliseconds,
			ResponseLossDisposition:               getMethod.ResponseLossDisposition,
			PeerRequirementBeforeDeliveryRequired: getMethod.PeerRequirementBeforeDeliveryRequired,
			MethodAuthorityDisposition:            getMethod.MethodAuthorityDisposition,
			RequestIDAuthorityDisposition:         getMethod.RequestIDAuthorityDisposition,
			ServiceMaxConnections:                 getMethod.ServiceMaxConnections,
			ConnectionMaxInFlight:                 getMethod.ConnectionMaxInFlight,
			ProcessMaxAdmittedRequests:            getMethod.ProcessMaxAdmittedRequests,
			ApplicationQueueCapacity:              getMethod.ApplicationQueueCapacity,
			InFlightRequestDataMaxBytes:           getMethod.InFlightRequestDataMaxBytes,
		}) != expectedGet {
		t.Fatal("Go GetRegisteredPlanV0 method record disagrees with generated TypeScript fixture")
	}

	var requestFixture passiveRegisterRequestFixture
	readPassiveJSON(t, "register-plan-v0.request.json", &requestFixture)
	request, err := NewPassiveRegisterPlanV0Request(
		requestFixture.ProtocolVersion,
		decodeFixedHex16[PassiveRequestID](t, requestFixture.RequestID),
		decodeFixedHex16[v0candidate.InstallationID](t, requestFixture.InstallationID),
		v0candidate.UInt53(requestFixture.EpochSequence),
		decodeFixedHex32[v0candidate.TrustEpochDigest](t, requestFixture.EpochDigest),
		readReferencedFixture(t, requestFixture.Body.ExactPlanBytes),
		readReferencedFixture(t, requestFixture.Body.RoleBindingBytes),
		readReferencedFixture(t, requestFixture.Body.SourceManifestBytes),
		readReferencedFixture(t, requestFixture.Body.SourceBytes),
	)
	if err != nil {
		t.Fatalf("Go rejected TypeScript-generated ordinary request: %v", err)
	}
	if requestFixture.ObjectType != "capsule.authenticated-local-ipc-register-plan-v0-request" ||
		requestFixture.ObjectVersion != 0 || requestFixture.FixtureSerialization != "exact-json-not-xpc-framing" ||
		request.ApplicationDataBytes() != requestFixture.ApplicationDataBytes {
		t.Fatal("Go ordinary request projection disagrees with generated fixture")
	}
	assertPassiveReference(t, readPassiveFixture(t, "register-plan-v0.method.json"), requestFixture.MethodRecord)

	var registerReply passiveRegisterReplyFixture
	readPassiveJSON(t, "register-plan-v0.reply.json", &registerReply)
	registration := readReferencedFixture(t, registerReply.Body.PlanRegistrationBytes)
	if registerReply.ObjectType != "capsule.authenticated-local-ipc-register-plan-v0-reply" ||
		registerReply.ObjectVersion != 0 || registerReply.FixtureSerialization != "exact-json-not-xpc-framing" ||
		registerReply.RequestID != requestFixture.RequestID ||
		registerReply.ApplicationDataBytes != len(registration) {
		t.Fatal("Go RegisterPlanV0 reply projection disagrees with generated fixture")
	}

	var getRequestFixture passiveGetRequestFixture
	readPassiveJSON(t, "get-registered-plan-v0.request.json", &getRequestFixture)
	getRequest := PassiveGetRegisteredPlanV0Request{
		ProtocolVersion: getRequestFixture.ProtocolVersion,
		RequestID:       decodeFixedHex16[PassiveRequestID](t, getRequestFixture.RequestID),
		InstallationID:  decodeFixedHex16[v0candidate.InstallationID](t, getRequestFixture.InstallationID),
		EpochSequence:   v0candidate.UInt53(getRequestFixture.EpochSequence),
		EpochDigest:     decodeFixedHex32[v0candidate.TrustEpochDigest](t, getRequestFixture.EpochDigest),
		RegistrationID:  decodeFixedHex16[v0candidate.RegistrationID](t, getRequestFixture.Body.RegistrationID),
	}
	if getRequestFixture.ObjectType != "capsule.authenticated-local-ipc-get-registered-plan-v0-request" ||
		getRequestFixture.ObjectVersion != 0 ||
		getRequestFixture.FixtureSerialization != "exact-json-not-xpc-framing" ||
		getRequest.ApplicationDataBytes() != getRequestFixture.ApplicationDataBytes ||
		getRequest.InstallationID != request.InstallationID ||
		getRequest.EpochSequence != request.EpochSequence || getRequest.EpochDigest != request.EpochDigest {
		t.Fatal("Go GetRegisteredPlanV0 request projection disagrees with generated fixture")
	}
	assertPassiveReference(t, readPassiveFixture(t, "get-registered-plan-v0.method.json"), getRequestFixture.MethodRecord)

	var getReply passiveGetReplyFixture
	readPassiveJSON(t, "get-registered-plan-v0.reply.json", &getReply)
	getReplyBytes := 0
	for _, reference := range []passiveFixtureReference{
		getReply.Body.ExactPlanBytes,
		getReply.Body.RoleBindingBytes,
		getReply.Body.PlanRegistrationBytes,
		getReply.Body.SourceManifestBytes,
		getReply.Body.SourceBytes,
	} {
		getReplyBytes += len(readReferencedFixture(t, reference))
	}
	if getReply.ObjectType != "capsule.authenticated-local-ipc-get-registered-plan-v0-reply" ||
		getReply.ObjectVersion != 0 || getReply.FixtureSerialization != "exact-json-not-xpc-framing" ||
		getReply.RequestID != getRequestFixture.RequestID || getReply.ApplicationDataBytes != getReplyBytes {
		t.Fatal("Go GetRegisteredPlanV0 reply projection disagrees with generated fixture")
	}
}

func TestPassiveRequestAndReplyExactCapsAndCapPlusOne(t *testing.T) {
	request, err := NewPassiveRegisterPlanV0Request(
		0,
		repeatedID[PassiveRequestID](0x41),
		repeatedID[v0candidate.InstallationID](0x11),
		7,
		repeatedDigest[v0candidate.TrustEpochDigest](0x22),
		bytes.Repeat([]byte{0x01}, v0candidate.ExecutionPlanMaxCBORBytes),
		bytes.Repeat([]byte{0x02}, RoleBindingRecordBytes),
		bytes.Repeat([]byte{0x03}, v0candidate.SourceManifestMaxCBORBytes),
		bytes.Repeat([]byte{'x'}, v0candidate.MJSMainSourceMaxBytes),
	)
	if err != nil {
		t.Fatalf("exact RegisterPlanV0 application cap refused: %v", err)
	}
	if request.ApplicationDataBytes() != RegisterPlanV0RequestDataMaxBytes {
		t.Fatalf("request application bytes = %d", request.ApplicationDataBytes())
	}
	if _, err := NewPassiveRegisterPlanV0Request(
		0,
		repeatedID[PassiveRequestID](0x41),
		repeatedID[v0candidate.InstallationID](0x11),
		7,
		repeatedDigest[v0candidate.TrustEpochDigest](0x22),
		bytes.Repeat([]byte{0x01}, v0candidate.ExecutionPlanMaxCBORBytes),
		bytes.Repeat([]byte{0x02}, RoleBindingRecordBytes),
		bytes.Repeat([]byte{0x03}, v0candidate.SourceManifestMaxCBORBytes),
		bytes.Repeat([]byte{'x'}, v0candidate.MJSMainSourceMaxBytes+1),
	); err == nil {
		t.Fatal("RegisterPlanV0 source cap plus one accepted")
	}

	maximumReply := GetRegisteredPlanV0Reply{
		plan:         bytes.Repeat([]byte{0x01}, v0candidate.ExecutionPlanMaxCBORBytes),
		bindings:     bytes.Repeat([]byte{0x02}, RoleBindingRecordBytes),
		registration: bytes.Repeat([]byte{0x03}, v0candidate.PlanRegistrationMaxCBORBytes),
		manifest:     bytes.Repeat([]byte{0x04}, v0candidate.SourceManifestMaxCBORBytes),
		source:       bytes.Repeat([]byte{'x'}, v0candidate.MJSMainSourceMaxBytes),
	}
	passiveReply := newPassiveGetRegisteredPlanV0Reply(repeatedID[PassiveRequestID](0x42), maximumReply)
	if passiveReply.ApplicationDataBytes() != GetRegisteredPlanV0ReplyDataMaxBytes {
		t.Fatalf("fetch reply application bytes = %d", passiveReply.ApplicationDataBytes())
	}
	maximumReply.source = bytes.Repeat([]byte{'x'}, v0candidate.MJSMainSourceMaxBytes+1)
	assertPassiveIntegrityPanic(t, "get-registered-plan-reply-shape", func() {
		newPassiveGetRegisteredPlanV0Reply(repeatedID[PassiveRequestID](0x42), maximumReply)
	})

	bindings := ordinaryBindings(t)
	ordinaryRequest, err := NewPassiveRegisterPlanV0Request(
		0,
		repeatedID[PassiveRequestID](0x41),
		bindings.InstallationID,
		7,
		bindings.EpochDigest,
		readFixture(t, "execution-plan/ordinary.cbor"),
		mustRoleBindings(t, bindings),
		readFixture(t, "source-manifest/ordinary.cbor"),
		[]byte("export default function (value) { return value; }\n"),
	)
	if err != nil {
		t.Fatal(err)
	}
	supervisor := SupervisorContext{
		InstallationID: bindings.InstallationID,
		EpochSequence:  7,
		EpochDigest:    bindings.EpochDigest,
		SupervisorID:   repeatedID[v0candidate.SupervisorID](0x55),
	}
	core := &countingPassiveCore{registerReply: bytes.Repeat([]byte{0xa1}, RegisterPlanV0ReplyDataMaxBytes)}
	boundary, err := newPassiveIPCBoundary(core, supervisor)
	if err != nil {
		t.Fatal(err)
	}
	registerReply, err := boundary.registerPlanV0(
		context.Background(), ExpectedRegisterPlanV0MethodRecord(), ordinaryRequest,
	)
	if err != nil || registerReply.ApplicationDataBytes() != RegisterPlanV0ReplyDataMaxBytes {
		t.Fatalf("exact RegisterPlanV0 reply cap refused: bytes=%d err=%v", registerReply.ApplicationDataBytes(), err)
	}
	core.registerReply = bytes.Repeat([]byte{0xa1}, RegisterPlanV0ReplyDataMaxBytes+1)
	assertPassiveIntegrityPanic(t, "register-plan-reply-cap", func() {
		_, _ = boundary.registerPlanV0(
			context.Background(), ExpectedRegisterPlanV0MethodRecord(), ordinaryRequest,
		)
	})
	if core.registerCalls != 2 {
		t.Fatal("reply cap oracle did not execute after the authority core")
	}
}

func TestPassiveAdmissionRefusesBeforeCore(t *testing.T) {
	bindings := ordinaryBindings(t)
	body, err := NewPassiveRegisterPlanV0Request(
		0,
		repeatedID[PassiveRequestID](0x41),
		bindings.InstallationID,
		7,
		bindings.EpochDigest,
		readFixture(t, "execution-plan/ordinary.cbor"),
		mustRoleBindings(t, bindings),
		readFixture(t, "source-manifest/ordinary.cbor"),
		[]byte("export default function (value) { return value; }\n"),
	)
	if err != nil {
		t.Fatal(err)
	}
	supervisor := SupervisorContext{
		InstallationID: bindings.InstallationID,
		EpochSequence:  7,
		EpochDigest:    bindings.EpochDigest,
		SupervisorID:   repeatedID[v0candidate.SupervisorID](0x55),
	}

	tests := []struct {
		name           string
		method         RegisterPlanV0MethodRecord
		mutateRequest  func(*PassiveRegisterPlanV0Request)
		classification Classification
		reason         string
	}{
		{name: "method-version", method: mutateRegisterMethod(func(value *RegisterPlanV0MethodRecord) { value.MethodVersion = 1 }), classification: Unsupported, reason: "register-plan-method-record-version"},
		{name: "record-version", method: mutateRegisterMethod(func(value *RegisterPlanV0MethodRecord) { value.RoleBindingRecordVersion = 1 }), classification: Unsupported, reason: "register-plan-method-record-version"},
		{name: "method", method: mutateRegisterMethod(func(value *RegisterPlanV0MethodRecord) { value.Method = GetRegisteredPlanV0Method }), classification: Unsupported, reason: "register-plan-method"},
		{name: "service", method: mutateRegisterMethod(func(value *RegisterPlanV0MethodRecord) { value.Service = BrokerServiceV0 }), classification: Authentication, reason: "register-plan-method-binding"},
		{name: "role", method: mutateRegisterMethod(func(value *RegisterPlanV0MethodRecord) { value.ExpectedRole = Broker }), classification: Authentication, reason: "register-plan-method-binding"},
		{name: "audience", method: mutateRegisterMethod(func(value *RegisterPlanV0MethodRecord) { value.Audience = "capsule.execution-supervisor" }), classification: Authentication, reason: "register-plan-method-binding"},
		{name: "purpose", method: mutateRegisterMethod(func(value *RegisterPlanV0MethodRecord) { value.Purpose = GetRegisteredPlanV0Purpose }), classification: Authentication, reason: "register-plan-method-binding"},
		{name: "protocol-version", method: ExpectedRegisterPlanV0MethodRecord(), mutateRequest: func(value *PassiveRegisterPlanV0Request) { value.ProtocolVersion = 1 }, classification: Unsupported, reason: "ipc-protocol-version"},
		{name: "zero-request-id", method: ExpectedRegisterPlanV0MethodRecord(), mutateRequest: func(value *PassiveRegisterPlanV0Request) { value.RequestID = PassiveRequestID{} }, classification: Schema, reason: "ipc-request-id"},
		{name: "installation", method: ExpectedRegisterPlanV0MethodRecord(), mutateRequest: func(value *PassiveRegisterPlanV0Request) { value.InstallationID[0] ^= 0xff }, classification: Binding, reason: "ipc-current-supervisor-state"},
		{name: "epoch-sequence", method: ExpectedRegisterPlanV0MethodRecord(), mutateRequest: func(value *PassiveRegisterPlanV0Request) { value.EpochSequence++ }, classification: Binding, reason: "ipc-current-supervisor-state"},
		{name: "epoch-digest", method: ExpectedRegisterPlanV0MethodRecord(), mutateRequest: func(value *PassiveRegisterPlanV0Request) { value.EpochDigest[0] ^= 0xff }, classification: Binding, reason: "ipc-current-supervisor-state"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			core := &countingPassiveCore{}
			boundary, err := newPassiveIPCBoundary(core, supervisor)
			if err != nil {
				t.Fatal(err)
			}
			candidate := clonePassiveRegisterRequest(body)
			if test.mutateRequest != nil {
				test.mutateRequest(&candidate)
			}
			_, err = boundary.registerPlanV0(context.Background(), test.method, candidate)
			if classification, ok := ErrorClassification(err); !ok || classification != test.classification {
				t.Fatalf("classification = %q, want %q (%v)", classification, test.classification, err)
			}
			if reason, ok := ErrorCode(err); !ok || reason != test.reason {
				t.Fatalf("reason = %q, want %q (%v)", reason, test.reason, err)
			}
			if core.registerCalls != 0 || core.getCalls != 0 {
				t.Fatal("passive admission refusal reached authority core")
			}
			refusal, ok := passiveRefusal(candidate.RequestID, err)
			if !ok || refusal.Classification != test.classification || refusal.Reason != test.reason {
				t.Fatalf("fixed refusal = %#v, recognized %v", refusal, ok)
			}
		})
	}
}

func TestPassiveGetAdmissionRefusesBeforeCore(t *testing.T) {
	bindings := ordinaryBindings(t)
	supervisor := SupervisorContext{
		InstallationID: bindings.InstallationID,
		EpochSequence:  7,
		EpochDigest:    bindings.EpochDigest,
		SupervisorID:   repeatedID[v0candidate.SupervisorID](0x55),
	}
	request := PassiveGetRegisteredPlanV0Request{
		ProtocolVersion: 0,
		RequestID:       repeatedID[PassiveRequestID](0x51),
		InstallationID:  bindings.InstallationID,
		EpochSequence:   7,
		EpochDigest:     bindings.EpochDigest,
		RegistrationID:  repeatedID[v0candidate.RegistrationID](0x77),
	}
	tests := []struct {
		name           string
		method         GetRegisteredPlanV0MethodRecord
		mutateRequest  func(*PassiveGetRegisteredPlanV0Request)
		classification Classification
		reason         string
	}{
		{name: "method-version", method: mutateGetMethod(func(value *GetRegisteredPlanV0MethodRecord) { value.MethodVersion = 1 }), classification: Unsupported, reason: "get-registered-plan-method-record-version"},
		{name: "record-version", method: mutateGetMethod(func(value *GetRegisteredPlanV0MethodRecord) { value.RoleBindingRecordVersion = 1 }), classification: Unsupported, reason: "get-registered-plan-method-record-version"},
		{name: "method", method: mutateGetMethod(func(value *GetRegisteredPlanV0MethodRecord) { value.Method = RegisterPlanV0Method }), classification: Unsupported, reason: "get-registered-plan-method"},
		{name: "service", method: mutateGetMethod(func(value *GetRegisteredPlanV0MethodRecord) { value.Service = DaemonServiceV0 }), classification: Authentication, reason: "get-registered-plan-method-binding"},
		{name: "role", method: mutateGetMethod(func(value *GetRegisteredPlanV0MethodRecord) { value.ExpectedRole = Daemon }), classification: Authentication, reason: "get-registered-plan-method-binding"},
		{name: "audience", method: mutateGetMethod(func(value *GetRegisteredPlanV0MethodRecord) { value.Audience = "capsule.execution-supervisor" }), classification: Authentication, reason: "get-registered-plan-method-binding"},
		{name: "purpose", method: mutateGetMethod(func(value *GetRegisteredPlanV0MethodRecord) { value.Purpose = RegisterPlanV0Purpose }), classification: Authentication, reason: "get-registered-plan-method-binding"},
		{name: "protocol-version", method: ExpectedGetRegisteredPlanV0MethodRecord(), mutateRequest: func(value *PassiveGetRegisteredPlanV0Request) { value.ProtocolVersion = 1 }, classification: Unsupported, reason: "ipc-protocol-version"},
		{name: "zero-request-id", method: ExpectedGetRegisteredPlanV0MethodRecord(), mutateRequest: func(value *PassiveGetRegisteredPlanV0Request) { value.RequestID = PassiveRequestID{} }, classification: Schema, reason: "ipc-request-id"},
		{name: "installation", method: ExpectedGetRegisteredPlanV0MethodRecord(), mutateRequest: func(value *PassiveGetRegisteredPlanV0Request) { value.InstallationID[0] ^= 0xff }, classification: Binding, reason: "ipc-current-supervisor-state"},
		{name: "epoch-sequence", method: ExpectedGetRegisteredPlanV0MethodRecord(), mutateRequest: func(value *PassiveGetRegisteredPlanV0Request) { value.EpochSequence++ }, classification: Binding, reason: "ipc-current-supervisor-state"},
		{name: "epoch-digest", method: ExpectedGetRegisteredPlanV0MethodRecord(), mutateRequest: func(value *PassiveGetRegisteredPlanV0Request) { value.EpochDigest[0] ^= 0xff }, classification: Binding, reason: "ipc-current-supervisor-state"},
		{name: "zero-registration-id", method: ExpectedGetRegisteredPlanV0MethodRecord(), mutateRequest: func(value *PassiveGetRegisteredPlanV0Request) { value.RegistrationID = v0candidate.RegistrationID{} }, classification: Schema, reason: "registration-id"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			core := &countingPassiveCore{getError: refused(LocalFailure, "unexpected-test-fetch")}
			boundary, err := newPassiveIPCBoundary(core, supervisor)
			if err != nil {
				t.Fatal(err)
			}
			candidate := request
			if test.mutateRequest != nil {
				test.mutateRequest(&candidate)
			}
			_, err = boundary.getRegisteredPlanV0(context.Background(), test.method, candidate)
			if classification, ok := ErrorClassification(err); !ok || classification != test.classification {
				t.Fatalf("classification = %q, want %q (%v)", classification, test.classification, err)
			}
			if reason, ok := ErrorCode(err); !ok || reason != test.reason {
				t.Fatalf("reason = %q, want %q (%v)", reason, test.reason, err)
			}
			if core.registerCalls != 0 || core.getCalls != 0 {
				t.Fatal("passive fetch admission refusal reached authority core")
			}
			refusal, ok := passiveRefusal(candidate.RequestID, err)
			if !ok || refusal.Classification != test.classification || refusal.Reason != test.reason {
				t.Fatalf("fixed refusal = %#v, recognized %v", refusal, ok)
			}
		})
	}
}

func TestPassiveResponseLossAndCopyOwnership(t *testing.T) {
	bindings := ordinaryBindings(t)
	plan := readFixture(t, "execution-plan/ordinary.cbor")
	roleBytes := mustRoleBindings(t, bindings)
	manifest := readFixture(t, "source-manifest/ordinary.cbor")
	source := []byte("export default function (value) { return value; }\n")
	request, err := NewPassiveRegisterPlanV0Request(
		0,
		repeatedID[PassiveRequestID](0x41),
		bindings.InstallationID,
		7,
		bindings.EpochDigest,
		plan,
		roleBytes,
		manifest,
		source,
	)
	if err != nil {
		t.Fatal(err)
	}
	plan[0] ^= 0xff
	roleBytes[0] ^= 0xff
	manifest[0] ^= 0xff
	source[0] ^= 0xff
	if bytes.Equal(request.ExactPlanBytes(), plan) || bytes.Equal(request.SourceBytes(), source) {
		t.Fatal("passive request retained caller storage")
	}
	returnedSource := request.SourceBytes()
	returnedSource[0] ^= 0xff
	if bytes.Equal(returnedSource, request.SourceBytes()) {
		t.Fatal("passive request accessor aliases retained storage")
	}

	ids := &sequenceIdentifiers{values: []v0candidate.RegistrationID{
		repeatedID[v0candidate.RegistrationID](0x77),
		repeatedID[v0candidate.RegistrationID](0x78),
	}}
	store := &FixedStore{}
	supervisor := SupervisorContext{
		InstallationID: bindings.InstallationID,
		EpochSequence:  7,
		EpochDigest:    bindings.EpochDigest,
		SupervisorID:   repeatedID[v0candidate.SupervisorID](0x55),
	}
	facade, err := NewFacade(store, exactRoleResolver{bindings}, ids, supervisor)
	if err != nil {
		t.Fatal(err)
	}
	boundary, err := newPassiveIPCBoundary(facade, supervisor)
	if err != nil {
		t.Fatal(err)
	}

	first, err := boundary.registerPlanV0(context.Background(), ExpectedRegisterPlanV0MethodRecord(), request)
	if err != nil {
		t.Fatal(err)
	}
	// Simulated response loss: discard first reply, then make a fresh transport
	// call with the same application data and a new correlation-only request ID.
	retry := clonePassiveRegisterRequest(request)
	retry.RequestID = repeatedID[PassiveRequestID](0x42)
	second, err := boundary.registerPlanV0(context.Background(), ExpectedRegisterPlanV0MethodRecord(), retry)
	if err != nil {
		t.Fatal(err)
	}
	if store.count() != 2 || bytes.Equal(first.PlanRegistrationBytes(), second.PlanRegistrationBytes()) {
		t.Fatal("lost registration reply retry did not retain two fresh registrations")
	}
	if !bytes.Equal(first.PlanRegistrationBytes(), readFixture(t, "plan-registration/ordinary.cbor")) {
		t.Fatal("passive RegisterPlanV0 reply differs from cross-language known answer")
	}

	fetchRequest := PassiveGetRegisteredPlanV0Request{
		ProtocolVersion: 0,
		RequestID:       repeatedID[PassiveRequestID](0x51),
		InstallationID:  bindings.InstallationID,
		EpochSequence:   7,
		EpochDigest:     bindings.EpochDigest,
		RegistrationID:  repeatedID[v0candidate.RegistrationID](0x77),
	}
	fetched, err := boundary.getRegisteredPlanV0(context.Background(), ExpectedGetRegisteredPlanV0MethodRecord(), fetchRequest)
	if err != nil {
		t.Fatal(err)
	}
	secondFetchRequest := fetchRequest
	secondFetchRequest.RequestID = repeatedID[PassiveRequestID](0x53)
	secondFetchRequest.RegistrationID = repeatedID[v0candidate.RegistrationID](0x78)
	secondFetched, err := boundary.getRegisteredPlanV0(
		context.Background(), ExpectedGetRegisteredPlanV0MethodRecord(), secondFetchRequest,
	)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(secondFetched.PlanRegistrationBytes(), second.PlanRegistrationBytes()) ||
		!bytes.Equal(secondFetched.ExactPlanBytes(), fetched.ExactPlanBytes()) ||
		!bytes.Equal(secondFetched.SourceBytes(), fetched.SourceBytes()) {
		t.Fatal("fresh registration retry was not separately readable with the same retained body")
	}
	fetchRequest.RequestID = repeatedID[PassiveRequestID](0x52)
	repeated, err := boundary.getRegisteredPlanV0(context.Background(), ExpectedGetRegisteredPlanV0MethodRecord(), fetchRequest)
	if err != nil {
		t.Fatal(err)
	}
	if store.count() != 2 || !bytes.Equal(fetched.ExactPlanBytes(), repeated.ExactPlanBytes()) ||
		!bytes.Equal(fetched.RoleBindingBytes(), repeated.RoleBindingBytes()) ||
		!bytes.Equal(fetched.PlanRegistrationBytes(), repeated.PlanRegistrationBytes()) ||
		!bytes.Equal(fetched.SourceManifestBytes(), repeated.SourceManifestBytes()) ||
		!bytes.Equal(fetched.SourceBytes(), repeated.SourceBytes()) {
		t.Fatal("lost fetch reply was not a byte-exact repeatable read")
	}
	mutated := fetched.SourceBytes()
	mutated[0] ^= 0xff
	if bytes.Equal(mutated, fetched.SourceBytes()) {
		t.Fatal("passive fetch reply accessor aliases retained reply storage")
	}
}

func mustRoleBindings(t *testing.T, bindings v0candidate.ExecutionPlanRoleBindings) []byte {
	t.Helper()
	encoded, err := EncodeRoleBindingsV0(bindings)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}

func mutateRegisterMethod(mutate func(*RegisterPlanV0MethodRecord)) RegisterPlanV0MethodRecord {
	value := ExpectedRegisterPlanV0MethodRecord()
	mutate(&value)
	return value
}

func mutateGetMethod(mutate func(*GetRegisteredPlanV0MethodRecord)) GetRegisteredPlanV0MethodRecord {
	value := ExpectedGetRegisteredPlanV0MethodRecord()
	mutate(&value)
	return value
}

func clonePassiveRegisterRequest(value PassiveRegisterPlanV0Request) PassiveRegisterPlanV0Request {
	value.exactPlanBytes = bytes.Clone(value.exactPlanBytes)
	value.roleBindingBytes = bytes.Clone(value.roleBindingBytes)
	value.sourceManifestBytes = bytes.Clone(value.sourceManifestBytes)
	value.sourceBytes = bytes.Clone(value.sourceBytes)
	return value
}

func readPassiveJSON(t *testing.T, name string, target any) {
	t.Helper()
	value := readPassiveFixture(t, name)
	decoder := json.NewDecoder(bytes.NewReader(value))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		t.Fatalf("decode passive fixture %s: %v", name, err)
	}
}

func readPassiveFixture(t *testing.T, name string) []byte {
	t.Helper()
	value, err := os.ReadFile(filepath.Join("..", "..", "..", "schemas", "conformance", "authenticated-local-ipc-v0", name))
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func readReferencedFixture(t *testing.T, reference passiveFixtureReference) []byte {
	t.Helper()
	value, err := os.ReadFile(filepath.Join("..", "..", "..", filepath.FromSlash(reference.Path)))
	if err != nil {
		t.Fatal(err)
	}
	assertPassiveReference(t, value, reference)
	return value
}

func assertPassiveReference(t *testing.T, value []byte, reference passiveFixtureReference) {
	t.Helper()
	if len(value) != reference.ByteLength {
		t.Fatalf("%s bytes = %d, want %d", reference.Path, len(value), reference.ByteLength)
	}
	digest := sha256.Sum256(value)
	if hex.EncodeToString(digest[:]) != reference.SHA256 {
		t.Fatalf("%s digest mismatch", reference.Path)
	}
}

func assertPassiveIntegrityPanic(t *testing.T, expected string, action func()) {
	t.Helper()
	defer func() {
		if recovered := recover(); recovered != expected {
			t.Fatalf("integrity panic = %v, want %q", recovered, expected)
		}
	}()
	action()
	t.Fatal("expected local-integrity panic")
}

func decodeFixedHex16[T ~[16]byte](t *testing.T, value string) T {
	t.Helper()
	raw, err := hex.DecodeString(value)
	if err != nil {
		t.Fatal(err)
	}
	var result T
	if len(raw) != 16 {
		t.Fatalf("fixed hex width = %d, want 16", len(raw))
	}
	copy(result[:], raw)
	return result
}

func decodeFixedHex32[T ~[32]byte](t *testing.T, value string) T {
	t.Helper()
	raw, err := hex.DecodeString(value)
	if err != nil {
		t.Fatal(err)
	}
	var result T
	if len(raw) != 32 {
		t.Fatalf("fixed hex width = %d, want 32", len(raw))
	}
	copy(result[:], raw)
	return result
}
