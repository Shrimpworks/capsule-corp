package authorityplane

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

type passiveSubmitMethodFixture struct {
	ObjectType                            string     `json:"objectType"`
	ObjectVersion                         uint64     `json:"objectVersion"`
	Method                                string     `json:"method"`
	MethodVersion                         uint64     `json:"methodVersion"`
	Service                               string     `json:"service"`
	ExpectedRole                          CallerRole `json:"expectedRole"`
	ExpectedSigningIdentifier             string     `json:"expectedSigningIdentifier"`
	Audience                              string     `json:"audience"`
	Purpose                               string     `json:"purpose"`
	RequestMediaType                      string     `json:"requestMediaType"`
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

type passiveSubmitRequestFixture struct {
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
		ExactJobProposalBytes passiveFixtureReference `json:"exactJobProposalBytes"`
	} `json:"body"`
	ApplicationDataBytes int `json:"applicationDataBytes"`
}

type passiveSubmitReplyFixture struct {
	ObjectType           string `json:"objectType"`
	ObjectVersion        uint64 `json:"objectVersion"`
	FixtureSerialization string `json:"fixtureSerialization"`
	RequestID            string `json:"requestId"`
	Body                 struct {
		RegistrationID string `json:"registrationId"`
	} `json:"body"`
	ApplicationDataBytes int `json:"applicationDataBytes"`
}

type controlledSubmissionCore struct {
	mu              sync.Mutex
	calls           int
	proposals       [][]byte
	started         chan struct{}
	release         chan struct{}
	respectDeadline bool
	registrationID  v0candidate.RegistrationID
}

func (c *controlledSubmissionCore) SubmitMainMJSV0(ctx context.Context, proposal []byte) (v0candidate.RegistrationID, error) {
	c.mu.Lock()
	c.calls++
	c.proposals = append(c.proposals, bytes.Clone(proposal))
	c.mu.Unlock()
	if c.started != nil {
		select {
		case c.started <- struct{}{}:
		default:
		}
	}
	if c.release != nil {
		if c.respectDeadline {
			select {
			case <-c.release:
			case <-ctx.Done():
				return v0candidate.RegistrationID{}, ctx.Err()
			}
		} else {
			<-c.release
		}
	}
	return c.registrationID, nil
}

func (c *controlledSubmissionCore) callCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.calls
}

func TestPassiveProductAdapterFreezesOnePrivateXPCSubmission(t *testing.T) {
	method := ExpectedSubmitMainMJSV0MethodRecord()
	if method.Method != SubmitMainMJSV0Method || method.MethodVersion != 0 ||
		method.Service != SubmitMainMJSV0Service || method.ExpectedRole != InternalAlphaCLI ||
		method.ExpectedSigningIdentifier != InternalAlphaCLISigningID ||
		method.Audience != SubmitMainMJSV0Audience || method.Purpose != SubmitMainMJSV0Purpose ||
		method.RequestMediaType != JobProposalV0MediaType {
		t.Fatal("SubmitMainMJSV0 method identity drift")
	}
	if !method.PeerRequirementBeforeDeliveryRequired ||
		method.MethodAuthorityDisposition != PassiveMethodAuthorityDisposition ||
		method.RequestIDAuthorityDisposition != PassiveRequestIDAuthorityDisposition {
		t.Fatal("submission authority ordering drift")
	}
	if method.RequestDataMaxBytes != 2_097_152 || method.ReplyDataMaxBytes != 16 ||
		method.DeadlineMilliseconds != 10_000 || method.ServiceMaxConnections != 4 ||
		method.ConnectionMaxInFlight != 1 || method.ProcessMaxAdmittedRequests != 4 ||
		method.ApplicationQueueCapacity != 0 || method.InFlightRequestDataMaxBytes != 8_388_608 {
		t.Fatal("submission flow contract drift")
	}

	proposal := readFixture(t, "job-proposal/ordinary.json")
	request, err := NewPassiveSubmitMainMJSV0Request(
		0, repeatedID[PassiveRequestID](0x31), repeatedID[v0candidate.InstallationID](0x11), 7,
		repeatedDigest[v0candidate.TrustEpochDigest](0x22), proposal,
	)
	if err != nil {
		t.Fatal(err)
	}
	proposal[0] ^= 0xff
	if bytes.Equal(proposal, request.ProposalBytes()) {
		t.Fatal("submission request retained caller proposal storage")
	}
	returned := request.ProposalBytes()
	returned[0] ^= 0xff
	if bytes.Equal(returned, request.ProposalBytes()) {
		t.Fatal("submission proposal accessor aliases retained bytes")
	}
	if _, err := NewPassiveSubmitMainMJSV0Request(
		0, repeatedID[PassiveRequestID](0x31), repeatedID[v0candidate.InstallationID](0x11), 7,
		repeatedDigest[v0candidate.TrustEpochDigest](0x22), bytes.Repeat([]byte{'x'}, JobProposalV0RawMaxBytes+1),
	); classificationOf(err) != Malformed {
		t.Fatalf("proposal cap plus one classification = %v", err)
	}
}

func TestPassiveProductAdapterCrossLanguageKnownAnswer(t *testing.T) {
	var fixture passiveSubmitMethodFixture
	readPassiveJSON(t, "submit-main-mjs-v0.method.json", &fixture)
	expected := ExpectedSubmitMainMJSV0MethodRecord()
	actual := SubmitMainMJSV0MethodRecord{
		Method: fixture.Method, MethodVersion: fixture.MethodVersion,
		Service: fixture.Service, ExpectedRole: fixture.ExpectedRole,
		ExpectedSigningIdentifier: fixture.ExpectedSigningIdentifier,
		Audience:                  fixture.Audience, Purpose: fixture.Purpose, RequestMediaType: fixture.RequestMediaType,
		RequestDataMaxBytes: fixture.RequestDataMaxBytes, ReplyDataMaxBytes: fixture.ReplyDataMaxBytes,
		DeadlineMilliseconds:                  fixture.DeadlineMilliseconds,
		ResponseLossDisposition:               fixture.ResponseLossDisposition,
		PeerRequirementBeforeDeliveryRequired: fixture.PeerRequirementBeforeDeliveryRequired,
		MethodAuthorityDisposition:            fixture.MethodAuthorityDisposition,
		RequestIDAuthorityDisposition:         fixture.RequestIDAuthorityDisposition,
		ServiceMaxConnections:                 fixture.ServiceMaxConnections,
		ConnectionMaxInFlight:                 fixture.ConnectionMaxInFlight,
		ProcessMaxAdmittedRequests:            fixture.ProcessMaxAdmittedRequests,
		ApplicationQueueCapacity:              fixture.ApplicationQueueCapacity,
		InFlightRequestDataMaxBytes:           fixture.InFlightRequestDataMaxBytes,
	}
	if fixture.ObjectType != "capsule.authenticated-local-ipc-submit-main-mjs-v0-method-record" ||
		fixture.ObjectVersion != 0 || actual != expected {
		t.Fatal("Go SubmitMainMJSV0 method record disagrees with generated Node fixture")
	}

	var requestFixture passiveSubmitRequestFixture
	readPassiveJSON(t, "submit-main-mjs-v0.request.json", &requestFixture)
	request, err := NewPassiveSubmitMainMJSV0Request(
		requestFixture.ProtocolVersion,
		decodeFixedHex16[PassiveRequestID](t, requestFixture.RequestID),
		decodeFixedHex16[v0candidate.InstallationID](t, requestFixture.InstallationID),
		v0candidate.UInt53(requestFixture.EpochSequence),
		decodeFixedHex32[v0candidate.TrustEpochDigest](t, requestFixture.EpochDigest),
		readReferencedFixture(t, requestFixture.Body.ExactJobProposalBytes),
	)
	if err != nil || request.ApplicationDataBytes() != requestFixture.ApplicationDataBytes {
		t.Fatalf("Go rejected generated submission request: %v", err)
	}
	assertPassiveReference(t, readPassiveFixture(t, "submit-main-mjs-v0.method.json"), requestFixture.MethodRecord)
	var decoded struct {
		Source struct {
			Files map[string]string `json:"files"`
		} `json:"source"`
	}
	if err := json.Unmarshal(request.ProposalBytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Source.Files["main.mjs"] != string(readReferencedFixture(t, passiveFixtureReference{
		Path: "schemas/conformance/authority-plane-v0/main.mjs", ByteLength: 50,
		SHA256: "681f39365de1369ee486fa34e88b993c60df5a835006b65e0d8916df717c31cc",
	})) {
		t.Fatal("proposal main.mjs differs from registered source fixture")
	}

	var reply passiveSubmitReplyFixture
	readPassiveJSON(t, "submit-main-mjs-v0.reply.json", &reply)
	if reply.ObjectType != "capsule.authenticated-local-ipc-submit-main-mjs-v0-reply" ||
		reply.ObjectVersion != 0 || reply.RequestID != requestFixture.RequestID ||
		reply.ApplicationDataBytes != 16 || decodeFixedHex16[v0candidate.RegistrationID](t, reply.Body.RegistrationID) == (v0candidate.RegistrationID{}) {
		t.Fatal("Go submission reply projection disagrees with generated fixture")
	}
}

func TestPassiveProductAdapterMethodAndHeaderRefusalsReachNoCore(t *testing.T) {
	supervisor := passiveProductSupervisor()
	core := &controlledSubmissionCore{registrationID: repeatedID[v0candidate.RegistrationID](0x77)}
	harness := newPassiveProductHarnessForTest(t, core, supervisor)
	request := passiveSubmitRequest(t, supervisor)

	tests := []struct {
		name   string
		method SubmitMainMJSV0MethodRecord
		mutate func(*PassiveSubmitMainMJSV0Request)
		want   Classification
	}{
		{name: "version", method: mutateSubmitMethod(func(value *SubmitMainMJSV0MethodRecord) { value.MethodVersion++ }), want: Unsupported},
		{name: "service", method: mutateSubmitMethod(func(value *SubmitMainMJSV0MethodRecord) { value.Service = DaemonServiceV0 }), want: Authentication},
		{name: "role", method: mutateSubmitMethod(func(value *SubmitMainMJSV0MethodRecord) { value.ExpectedRole = Daemon }), want: Authentication},
		{name: "purpose", method: mutateSubmitMethod(func(value *SubmitMainMJSV0MethodRecord) { value.Purpose = RegisterPlanV0Purpose }), want: Authentication},
		{name: "audience", method: mutateSubmitMethod(func(value *SubmitMainMJSV0MethodRecord) { value.Audience = PassiveIPCAudience }), want: Authentication},
		{name: "installation", method: ExpectedSubmitMainMJSV0MethodRecord(), mutate: func(value *PassiveSubmitMainMJSV0Request) { value.InstallationID[0] ^= 0xff }, want: Binding},
		{name: "epoch", method: ExpectedSubmitMainMJSV0MethodRecord(), mutate: func(value *PassiveSubmitMainMJSV0Request) { value.EpochDigest[0] ^= 0xff }, want: Binding},
		{name: "request-id", method: ExpectedSubmitMainMJSV0MethodRecord(), mutate: func(value *PassiveSubmitMainMJSV0Request) { value.RequestID = PassiveRequestID{} }, want: Schema},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := request
			candidate.proposalBytes = bytes.Clone(request.proposalBytes)
			if test.mutate != nil {
				test.mutate(&candidate)
			}
			_, err := harness.submitMainMJSV0(context.Background(), repeatedID[PassiveConnectionID](0x41), test.method, candidate)
			if classificationOf(err) != test.want {
				t.Fatalf("classification = %v, want %v (%v)", classificationOf(err), test.want, err)
			}
		})
	}
	if core.callCount() != 0 {
		t.Fatal("pre-dispatch refusals reached submission core")
	}
}

func TestPassiveFlowControllerRefusesConnectionServiceProcessAndBytesWithoutQueue(t *testing.T) {
	pool := newPassiveAdmissionPool()
	limits := flowLimitsForSubmit(ExpectedSubmitMainMJSV0MethodRecord())
	releases := make([]func(), 0, 4)
	for index := byte(1); index <= 4; index++ {
		release, err := pool.acquire(passiveDaemonProcess, SubmitMainMJSV0Service, repeatedID[PassiveConnectionID](index), JobProposalV0RawMaxBytes, limits)
		if err != nil {
			t.Fatalf("admit connection %d: %v", index, err)
		}
		releases = append(releases, release)
	}
	if _, err := pool.acquire(passiveDaemonProcess, SubmitMainMJSV0Service, repeatedID[PassiveConnectionID](1), 1, limits); errorCodeOf(err) != "ipc-connection-in-flight" {
		t.Fatalf("same-connection overload = %v", err)
	}
	if _, err := pool.acquire(passiveDaemonProcess, SubmitMainMJSV0Service, repeatedID[PassiveConnectionID](5), 1, limits); errorCodeOf(err) != "ipc-service-connections" {
		t.Fatalf("service cap plus one = %v", err)
	}
	for _, release := range releases {
		release()
	}

	supervisorLimits := flowLimitsForSupervisor(
		SupervisorServiceMaxConnections, SupervisorConnectionMaxInFlight,
		SupervisorProcessMaxAdmittedRequests, SupervisorApplicationQueueCapacity,
		SupervisorInFlightRequestDataMaxBytes,
	)
	registerReleases := make([]func(), 0, 4)
	fetchReleases := make([]func(), 0, 4)
	for index := byte(1); index <= 4; index++ {
		release, err := pool.acquire(passiveSupervisorProcess, DaemonServiceV0, repeatedID[PassiveConnectionID](index), RegisterPlanV0RequestDataMaxBytes, supervisorLimits)
		if err != nil {
			t.Fatal(err)
		}
		registerReleases = append(registerReleases, release)
		release, err = pool.acquire(passiveSupervisorProcess, BrokerServiceV0, repeatedID[PassiveConnectionID](index+8), GetRegisteredPlanV0RequestDataMaxBytes, supervisorLimits)
		if err != nil {
			t.Fatal(err)
		}
		fetchReleases = append(fetchReleases, release)
	}
	if _, err := pool.acquire(passiveSupervisorProcess, BrokerServiceV0, repeatedID[PassiveConnectionID](0x20), 16, supervisorLimits); errorCodeOf(err) != "ipc-service-connections" {
		t.Fatalf("supervisor ninth request = %v", err)
	}
	for _, release := range append(registerReleases, fetchReleases...) {
		release()
	}

	byteLimits := supervisorLimits
	byteLimits.inFlightRequestDataMaxBytes =
		RegisterPlanV0RequestDataMaxBytes + GetRegisteredPlanV0RequestDataMaxBytes
	registerRelease, err := pool.acquire(
		passiveSupervisorProcess, DaemonServiceV0, repeatedID[PassiveConnectionID](0x31),
		RegisterPlanV0RequestDataMaxBytes, byteLimits,
	)
	if err != nil {
		t.Fatalf("admit mixed register bytes: %v", err)
	}
	getRelease, err := pool.acquire(
		passiveSupervisorProcess, BrokerServiceV0, repeatedID[PassiveConnectionID](0x32),
		GetRegisteredPlanV0RequestDataMaxBytes, byteLimits,
	)
	if err != nil {
		t.Fatalf("admit mixed get bytes: %v", err)
	}
	if _, err := pool.acquire(passiveSupervisorProcess, BrokerServiceV0, repeatedID[PassiveConnectionID](0x33), 1, byteLimits); errorCodeOf(err) != "ipc-process-in-flight-bytes" {
		t.Fatalf("mixed aggregate byte cap plus one = %v", err)
	}
	getRelease()
	postRelease, err := pool.acquire(
		passiveSupervisorProcess, BrokerServiceV0, repeatedID[PassiveConnectionID](0x34), 1, byteLimits,
	)
	if err != nil {
		t.Fatalf("post-release byte re-admission: %v", err)
	}
	postRelease()
	registerRelease()
}

func TestPassiveProductAdapterRoutesOnlyMethodSpecificSupervisorCalls(t *testing.T) {
	supervisor := passiveProductSupervisor()
	core := &countingPassiveCore{
		registerReply: readFixture(t, "plan-registration/ordinary.cbor"),
		getReply: GetRegisteredPlanV0Reply{
			plan:         readFixture(t, "execution-plan/ordinary.cbor"),
			bindings:     readPassiveFixture(t, "../authority-plane-v0/role-bindings.bin"),
			registration: readFixture(t, "plan-registration/ordinary.cbor"),
			manifest:     readPassiveFixture(t, "../authority-plane-v0/source-manifest.cbor"),
			source:       readPassiveFixture(t, "../authority-plane-v0/main.mjs"),
		},
	}
	authority, err := newPassiveIPCBoundary(core, supervisor)
	if err != nil {
		t.Fatal(err)
	}
	harness, err := newPassiveProductAdapterHarness(
		&controlledSubmissionCore{registrationID: repeatedID[v0candidate.RegistrationID](0x77)},
		authority,
	)
	if err != nil {
		t.Fatal(err)
	}

	register, err := NewPassiveRegisterPlanV0Request(
		0, repeatedID[PassiveRequestID](0x31), supervisor.InstallationID,
		supervisor.EpochSequence, supervisor.EpochDigest,
		readFixture(t, "execution-plan/ordinary.cbor"),
		readPassiveFixture(t, "../authority-plane-v0/role-bindings.bin"),
		readPassiveFixture(t, "../authority-plane-v0/source-manifest.cbor"),
		readPassiveFixture(t, "../authority-plane-v0/main.mjs"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := harness.registerPlanV0(
		context.Background(), repeatedID[PassiveConnectionID](0x41),
		ExpectedRegisterPlanV0MethodRecord(), register,
	); err != nil {
		t.Fatal(err)
	}

	fetch := PassiveGetRegisteredPlanV0Request{
		ProtocolVersion: 0,
		RequestID:       repeatedID[PassiveRequestID](0x32),
		InstallationID:  supervisor.InstallationID,
		EpochSequence:   supervisor.EpochSequence,
		EpochDigest:     supervisor.EpochDigest,
		RegistrationID:  repeatedID[v0candidate.RegistrationID](0x77),
	}
	if _, err := harness.getRegisteredPlanV0(
		context.Background(), repeatedID[PassiveConnectionID](0x42),
		ExpectedGetRegisteredPlanV0MethodRecord(), fetch,
	); err != nil {
		t.Fatal(err)
	}
	if core.registerCalls != 1 || core.getCalls != 1 {
		t.Fatalf("method-specific core calls = register %d, fetch %d", core.registerCalls, core.getCalls)
	}
}

func TestPassiveCancellationDeadlineStallAndResponseLossSemantics(t *testing.T) {
	supervisor := passiveProductSupervisor()
	request := passiveSubmitRequest(t, supervisor)
	connection := repeatedID[PassiveConnectionID](0x41)

	preCancelled, cancel := context.WithCancel(context.Background())
	cancel()
	preCore := &controlledSubmissionCore{registrationID: repeatedID[v0candidate.RegistrationID](0x77)}
	preHarness := newPassiveProductHarnessForTest(t, preCore, supervisor)
	_, err := preHarness.submitMainMJSV0(preCancelled, connection, ExpectedSubmitMainMJSV0MethodRecord(), request)
	if dispositionOf(err) != "caller-cancelled-before-dispatch" || preCore.callCount() != 0 {
		t.Fatalf("pre-dispatch cancellation = %v, calls %d", err, preCore.callCount())
	}

	release := make(chan struct{})
	stalled := &controlledSubmissionCore{
		started: make(chan struct{}, 1), release: release,
		registrationID: repeatedID[v0candidate.RegistrationID](0x77),
	}
	stallHarness := newPassiveProductHarnessForTest(t, stalled, supervisor)
	stallHarness.deadlineDuration = func(uint64) time.Duration { return 15 * time.Millisecond }
	_, err = stallHarness.submitMainMJSV0(context.Background(), connection, ExpectedSubmitMainMJSV0MethodRecord(), request)
	if dispositionOf(err) != "method-deadline-after-dispatch-response-unknown" || stalled.callCount() != 1 {
		t.Fatalf("downstream stall = %v, calls %d", err, stalled.callCount())
	}
	_, err = stallHarness.submitMainMJSV0(context.Background(), connection, ExpectedSubmitMainMJSV0MethodRecord(), request)
	if errorCodeOf(err) != "ipc-connection-in-flight" {
		t.Fatalf("stalled call did not retain bounded slot: %v", err)
	}
	close(release)
	waitForPassiveSlotRelease(t, stallHarness.pool, connection)

	responseRelease := make(chan struct{})
	responseCore := &controlledSubmissionCore{
		started: make(chan struct{}, 1), release: responseRelease,
		registrationID: repeatedID[v0candidate.RegistrationID](0x77),
	}
	responseHarness := newPassiveProductHarnessForTest(t, responseCore, supervisor)
	caller, cancelCaller := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, callErr := responseHarness.submitMainMJSV0(caller, connection, ExpectedSubmitMainMJSV0MethodRecord(), request)
		result <- callErr
	}()
	<-responseCore.started
	cancelCaller()
	if err := <-result; dispositionOf(err) != "caller-cancelled-after-dispatch-response-unknown" {
		t.Fatalf("post-dispatch cancellation = %v", err)
	}
	close(responseRelease)
	waitForPassiveSlotRelease(t, responseHarness.pool, connection)

	// Response loss never turns request ID into deduplication authority. A new
	// transport call reaches the core again; the downstream registration
	// semantics decide whether a fresh registration was committed.
	retryRequest := request
	retryRequest.RequestID = repeatedID[PassiveRequestID](0x32)
	if _, err := responseHarness.submitMainMJSV0(context.Background(), repeatedID[PassiveConnectionID](0x42), ExpectedSubmitMainMJSV0MethodRecord(), retryRequest); err != nil {
		t.Fatal(err)
	}
	if responseCore.callCount() != 2 {
		t.Fatalf("response-loss retry calls = %d", responseCore.callCount())
	}
}

func passiveProductSupervisor() SupervisorContext {
	return SupervisorContext{
		InstallationID: repeatedID[v0candidate.InstallationID](0x11),
		EpochSequence:  7,
		EpochDigest:    repeatedDigest[v0candidate.TrustEpochDigest](0x22),
		SupervisorID:   repeatedID[v0candidate.SupervisorID](0x55),
	}
}

func passiveSubmitRequest(t *testing.T, supervisor SupervisorContext) PassiveSubmitMainMJSV0Request {
	t.Helper()
	request, err := NewPassiveSubmitMainMJSV0Request(
		0, repeatedID[PassiveRequestID](0x31), supervisor.InstallationID,
		supervisor.EpochSequence, supervisor.EpochDigest, readFixture(t, "job-proposal/ordinary.json"),
	)
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func newPassiveProductHarnessForTest(t *testing.T, submission passiveSubmissionCore, supervisor SupervisorContext) *passiveProductAdapterHarness {
	t.Helper()
	authority, err := newPassiveIPCBoundary(&countingPassiveCore{}, supervisor)
	if err != nil {
		t.Fatal(err)
	}
	harness, err := newPassiveProductAdapterHarness(submission, authority)
	if err != nil {
		t.Fatal(err)
	}
	return harness
}

func mutateSubmitMethod(mutate func(*SubmitMainMJSV0MethodRecord)) SubmitMainMJSV0MethodRecord {
	value := ExpectedSubmitMainMJSV0MethodRecord()
	mutate(&value)
	return value
}

func classificationOf(err error) Classification {
	classification, _ := ErrorClassification(err)
	return classification
}

func errorCodeOf(err error) string {
	code, _ := ErrorCode(err)
	return code
}

func dispositionOf(err error) string {
	disposition, _ := passiveTransportDisposition(err)
	return disposition
}

func waitForPassiveSlotRelease(t *testing.T, pool *passiveAdmissionPool, connection PassiveConnectionID) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	limits := flowLimitsForSubmit(ExpectedSubmitMainMJSV0MethodRecord())
	for time.Now().Before(deadline) {
		release, err := pool.acquire(passiveDaemonProcess, SubmitMainMJSV0Service, connection, 1, limits)
		if err == nil {
			release()
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("passive admission slot did not release")
}
