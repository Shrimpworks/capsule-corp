package authorityplane

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

// This file freezes one passive internal-alpha product-adapter candidate. It
// creates no XPC listener, service registration, peer authentication, CLI,
// daemon consumer, signer, runtime, backend, or guest.
const (
	SubmitMainMJSV0Method        = "SubmitMainMJSV0"
	SubmitMainMJSV0MethodVersion = uint64(0)
	SubmitMainMJSV0Service       = "com.capsulecorp.capsule.daemon.cli.v0"
	SubmitMainMJSV0Audience      = "capsule.daemon.local.v0"
	SubmitMainMJSV0Purpose       = "capsule.ipc.submit-main-mjs.v0"
	InternalAlphaCLISigningID    = "com.capsulecorp.capsule.cli"

	JobProposalV0MediaType                 = "application/capsule.job-proposal+json;v=0"
	JobProposalV0RawMaxBytes               = 2_097_152
	SubmitMainMJSV0RequestDataMaxBytes     = JobProposalV0RawMaxBytes
	SubmitMainMJSV0ReplyDataMaxBytes       = 16
	SubmitMainMJSV0DeadlineMilliseconds    = uint64(10_000)
	SubmitMainMJSV0ResponseLossDisposition = "committed-retry-may-create-fresh-registration"

	DaemonServiceMaxConnections       = uint64(4)
	DaemonConnectionMaxInFlight       = uint64(1)
	DaemonProcessMaxAdmittedRequests  = uint64(4)
	DaemonApplicationQueueCapacity    = uint64(0)
	DaemonInFlightRequestDataMaxBytes = uint64(DaemonProcessMaxAdmittedRequests * SubmitMainMJSV0RequestDataMaxBytes)
)

var (
	_ [2_097_152 - JobProposalV0RawMaxBytes]byte
	_ [JobProposalV0RawMaxBytes - 2_097_152]byte
	_ [8_388_608 - DaemonInFlightRequestDataMaxBytes]byte
	_ [DaemonInFlightRequestDataMaxBytes - 8_388_608]byte
	_ [2_626_696 - SupervisorInFlightRequestDataMaxBytes]byte
	_ [SupervisorInFlightRequestDataMaxBytes - 2_626_696]byte
)

// field-authority-object: capsule.authenticated-local-ipc-submit-main-mjs-v0-method-record v0
type SubmitMainMJSV0MethodRecord struct {
	Method                                string
	MethodVersion                         uint64
	Service                               string
	ExpectedRole                          CallerRole
	ExpectedSigningIdentifier             string
	Audience                              string
	Purpose                               string
	RequestMediaType                      string
	RequestDataMaxBytes                   uint64
	ReplyDataMaxBytes                     uint64
	DeadlineMilliseconds                  uint64
	ResponseLossDisposition               string
	PeerRequirementBeforeDeliveryRequired bool
	MethodAuthorityDisposition            string
	RequestIDAuthorityDisposition         string
	ServiceMaxConnections                 uint64
	ConnectionMaxInFlight                 uint64
	ProcessMaxAdmittedRequests            uint64
	ApplicationQueueCapacity              uint64
	InFlightRequestDataMaxBytes           uint64
}

func ExpectedSubmitMainMJSV0MethodRecord() SubmitMainMJSV0MethodRecord {
	return SubmitMainMJSV0MethodRecord{
		Method: SubmitMainMJSV0Method, MethodVersion: SubmitMainMJSV0MethodVersion,
		Service: SubmitMainMJSV0Service, ExpectedRole: InternalAlphaCLI,
		ExpectedSigningIdentifier: InternalAlphaCLISigningID,
		Audience:                  SubmitMainMJSV0Audience, Purpose: SubmitMainMJSV0Purpose,
		RequestMediaType:                      JobProposalV0MediaType,
		RequestDataMaxBytes:                   SubmitMainMJSV0RequestDataMaxBytes,
		ReplyDataMaxBytes:                     SubmitMainMJSV0ReplyDataMaxBytes,
		DeadlineMilliseconds:                  SubmitMainMJSV0DeadlineMilliseconds,
		ResponseLossDisposition:               SubmitMainMJSV0ResponseLossDisposition,
		PeerRequirementBeforeDeliveryRequired: true,
		MethodAuthorityDisposition:            PassiveMethodAuthorityDisposition,
		RequestIDAuthorityDisposition:         PassiveRequestIDAuthorityDisposition,
		ServiceMaxConnections:                 DaemonServiceMaxConnections,
		ConnectionMaxInFlight:                 DaemonConnectionMaxInFlight,
		ProcessMaxAdmittedRequests:            DaemonProcessMaxAdmittedRequests,
		ApplicationQueueCapacity:              DaemonApplicationQueueCapacity,
		InFlightRequestDataMaxBytes:           DaemonInFlightRequestDataMaxBytes,
	}
}

// field-authority-object: capsule.authenticated-local-ipc-submit-main-mjs-v0-request v0
type PassiveSubmitMainMJSV0Request struct {
	ProtocolVersion uint64
	RequestID       PassiveRequestID
	InstallationID  v0candidate.InstallationID
	EpochSequence   v0candidate.UInt53
	EpochDigest     v0candidate.TrustEpochDigest
	proposalBytes   []byte
}

func NewPassiveSubmitMainMJSV0Request(
	protocolVersion uint64,
	requestID PassiveRequestID,
	installationID v0candidate.InstallationID,
	epochSequence v0candidate.UInt53,
	epochDigest v0candidate.TrustEpochDigest,
	proposalBytes []byte,
) (PassiveSubmitMainMJSV0Request, error) {
	if len(proposalBytes) == 0 || len(proposalBytes) > JobProposalV0RawMaxBytes {
		return PassiveSubmitMainMJSV0Request{}, refused(Malformed, "submit-main-mjs-proposal-bytes")
	}
	return PassiveSubmitMainMJSV0Request{
		ProtocolVersion: protocolVersion, RequestID: requestID,
		InstallationID: installationID, EpochSequence: epochSequence, EpochDigest: epochDigest,
		proposalBytes: bytes.Clone(proposalBytes),
	}, nil
}

func (r PassiveSubmitMainMJSV0Request) ProposalBytes() []byte     { return bytes.Clone(r.proposalBytes) }
func (r PassiveSubmitMainMJSV0Request) ApplicationDataBytes() int { return len(r.proposalBytes) }

// field-authority-object: capsule.authenticated-local-ipc-submit-main-mjs-v0-reply v0
type PassiveSubmitMainMJSV0Reply struct {
	RequestID      PassiveRequestID
	RegistrationID v0candidate.RegistrationID
}

func (r PassiveSubmitMainMJSV0Reply) ApplicationDataBytes() int { return len(r.RegistrationID) }

type passiveSubmissionCore interface {
	SubmitMainMJSV0(context.Context, []byte) (v0candidate.RegistrationID, error)
}

type PassiveConnectionID [16]byte

type passiveProcess string

const (
	passiveDaemonProcess     passiveProcess = "daemon"
	passiveSupervisorProcess passiveProcess = "supervisor"
)

type passiveFlowLimits struct {
	serviceMaxConnections       uint64
	connectionMaxInFlight       uint64
	processMaxAdmittedRequests  uint64
	applicationQueueCapacity    uint64
	inFlightRequestDataMaxBytes uint64
}

type passiveAdmissionPool struct {
	mu           sync.Mutex
	processes    map[passiveProcess]uint64
	processBytes map[passiveProcess]uint64
	connections  map[string]uint64
	services     map[string]map[PassiveConnectionID]struct{}
}

func newPassiveAdmissionPool() *passiveAdmissionPool {
	return &passiveAdmissionPool{
		processes:    make(map[passiveProcess]uint64),
		processBytes: make(map[passiveProcess]uint64),
		connections:  make(map[string]uint64),
		services:     make(map[string]map[PassiveConnectionID]struct{}),
	}
}

func (p *passiveAdmissionPool) acquire(
	process passiveProcess,
	service string,
	connection PassiveConnectionID,
	requestBytes uint64,
	limits passiveFlowLimits,
) (func(), error) {
	if p == nil || connection == (PassiveConnectionID{}) {
		return nil, refused(Schema, "ipc-connection-id")
	}
	if limits.applicationQueueCapacity != 0 {
		return nil, refused(LocalFailure, "ipc-queue-contract")
	}
	key := service + ":" + fmt.Sprintf("%x", connection)
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.connections[key] >= limits.connectionMaxInFlight {
		return nil, refused(Capacity, "ipc-connection-in-flight")
	}
	activeConnections := p.services[service]
	if activeConnections == nil {
		activeConnections = make(map[PassiveConnectionID]struct{})
		p.services[service] = activeConnections
	}
	if _, active := activeConnections[connection]; !active && uint64(len(activeConnections)) >= limits.serviceMaxConnections {
		return nil, refused(Capacity, "ipc-service-connections")
	}
	if p.processes[process] >= limits.processMaxAdmittedRequests {
		return nil, refused(Capacity, "ipc-process-in-flight")
	}
	if requestBytes > limits.inFlightRequestDataMaxBytes-p.processBytes[process] {
		return nil, refused(Capacity, "ipc-process-in-flight-bytes")
	}
	p.connections[key]++
	activeConnections[connection] = struct{}{}
	p.processes[process]++
	p.processBytes[process] += requestBytes
	var once sync.Once
	return func() {
		once.Do(func() {
			p.mu.Lock()
			defer p.mu.Unlock()
			p.connections[key]--
			if p.connections[key] == 0 {
				delete(p.connections, key)
				delete(activeConnections, connection)
			}
			p.processes[process]--
			p.processBytes[process] -= requestBytes
		})
	}, nil
}

type passiveTransportError struct{ disposition string }

func (e *passiveTransportError) Error() string { return e.disposition }

func passiveTransportDisposition(err error) (string, bool) {
	value, ok := err.(*passiveTransportError)
	if !ok {
		return "", false
	}
	return value.disposition, true
}

type passiveProductAdapterHarness struct {
	submission       passiveSubmissionCore
	authority        *passiveIPCBoundary
	pool             *passiveAdmissionPool
	deadlineDuration func(uint64) time.Duration
}

func newPassiveProductAdapterHarness(
	submission passiveSubmissionCore,
	authority *passiveIPCBoundary,
) (*passiveProductAdapterHarness, error) {
	if submission == nil || authority == nil {
		return nil, refused(LocalFailure, "passive-product-adapter-core")
	}
	return &passiveProductAdapterHarness{
		submission: submission, authority: authority,
		pool: newPassiveAdmissionPool(), deadlineDuration: passiveMethodDeadline,
	}, nil
}

// passiveMethodDeadline converts a method record's DeadlineMilliseconds into
// the time.Duration passed to runPassiveAdmitted. It computes the value
// directly instead of matching a hardcoded set of known deadlines, so an
// unwired or future method's deadline is never rejected with a panic. Method
// records reach this function only after validateSubmitMethod,
// validateRegisterMethod, or validateGetMethod has matched every field of
// the received record against its fixed expected record, so
// DeadlineMilliseconds is always one of this package's non-zero constants
// on every path wired today; a value of zero would still convert cleanly to
// an already-expired deadline (context.WithTimeout treats a non-positive
// duration as immediately Done), which is a fail-closed outcome, not an
// unbounded or fail-open one.
func passiveMethodDeadline(milliseconds uint64) time.Duration {
	return time.Duration(milliseconds) * time.Millisecond // #nosec G115 -- every wired caller supplies a validated method-record deadline constant far below the millisecond range that would overflow int64 nanoseconds.
}

func (h *passiveProductAdapterHarness) submitMainMJSV0(
	caller context.Context,
	connection PassiveConnectionID,
	method SubmitMainMJSV0MethodRecord,
	request PassiveSubmitMainMJSV0Request,
) (PassiveSubmitMainMJSV0Reply, error) {
	if err := validateSubmitMethod(method); err != nil {
		return PassiveSubmitMainMJSV0Reply{}, err
	}
	if err := h.authority.validateHeader(request.ProtocolVersion, request.RequestID, request.InstallationID, request.EpochSequence, request.EpochDigest); err != nil {
		return PassiveSubmitMainMJSV0Reply{}, err
	}
	requestBytes := uint64(request.ApplicationDataBytes()) // #nosec G115 -- constructor bounds the nonnegative byte-slice length.
	return runPassiveAdmitted(caller, h.pool, passiveDaemonProcess, method.Service, connection,
		requestBytes, flowLimitsForSubmit(method), h.deadlineDuration(method.DeadlineMilliseconds),
		func(coreContext context.Context) (PassiveSubmitMainMJSV0Reply, error) {
			registrationID, err := h.submission.SubmitMainMJSV0(coreContext, request.ProposalBytes())
			if err != nil {
				return PassiveSubmitMainMJSV0Reply{}, err
			}
			if registrationID == (v0candidate.RegistrationID{}) {
				panic("submit-main-mjs-reply-shape")
			}
			return PassiveSubmitMainMJSV0Reply{RequestID: request.RequestID, RegistrationID: registrationID}, nil
		})
}

func (h *passiveProductAdapterHarness) registerPlanV0(
	caller context.Context,
	connection PassiveConnectionID,
	method RegisterPlanV0MethodRecord,
	request PassiveRegisterPlanV0Request,
) (PassiveRegisterPlanV0Reply, error) {
	if err := validateRegisterMethod(method); err != nil {
		return PassiveRegisterPlanV0Reply{}, err
	}
	if err := h.authority.validateHeader(request.ProtocolVersion, request.RequestID, request.InstallationID, request.EpochSequence, request.EpochDigest); err != nil {
		return PassiveRegisterPlanV0Reply{}, err
	}
	requestBytes := uint64(request.ApplicationDataBytes()) // #nosec G115 -- request constructor caps every nonnegative slice length and their sum.
	return runPassiveAdmitted(caller, h.pool, passiveSupervisorProcess, method.Service, connection,
		requestBytes, flowLimitsForSupervisor(method.ServiceMaxConnections, method.ConnectionMaxInFlight, method.ProcessMaxAdmittedRequests, method.ApplicationQueueCapacity, method.InFlightRequestDataMaxBytes),
		h.deadlineDuration(method.DeadlineMilliseconds),
		func(coreContext context.Context) (PassiveRegisterPlanV0Reply, error) {
			return h.authority.registerPlanV0(coreContext, method, request)
		})
}

func (h *passiveProductAdapterHarness) getRegisteredPlanV0(
	caller context.Context,
	connection PassiveConnectionID,
	method GetRegisteredPlanV0MethodRecord,
	request PassiveGetRegisteredPlanV0Request,
) (PassiveGetRegisteredPlanV0Reply, error) {
	if err := validateGetMethod(method); err != nil {
		return PassiveGetRegisteredPlanV0Reply{}, err
	}
	if err := h.authority.validateHeader(request.ProtocolVersion, request.RequestID, request.InstallationID, request.EpochSequence, request.EpochDigest); err != nil {
		return PassiveGetRegisteredPlanV0Reply{}, err
	}
	requestBytes := uint64(request.ApplicationDataBytes()) // #nosec G115 -- RegistrationID length is the fixed nonnegative 16-byte width.
	return runPassiveAdmitted(caller, h.pool, passiveSupervisorProcess, method.Service, connection,
		requestBytes, flowLimitsForSupervisor(method.ServiceMaxConnections, method.ConnectionMaxInFlight, method.ProcessMaxAdmittedRequests, method.ApplicationQueueCapacity, method.InFlightRequestDataMaxBytes),
		h.deadlineDuration(method.DeadlineMilliseconds),
		func(coreContext context.Context) (PassiveGetRegisteredPlanV0Reply, error) {
			return h.authority.getRegisteredPlanV0(coreContext, method, request)
		})
}

type passiveCallResult[T any] struct {
	value T
	err   error
}

func runPassiveAdmitted[T any](
	caller context.Context,
	pool *passiveAdmissionPool,
	process passiveProcess,
	service string,
	connection PassiveConnectionID,
	requestBytes uint64,
	limits passiveFlowLimits,
	deadline time.Duration,
	call func(context.Context) (T, error),
) (T, error) {
	var zero T
	release, err := pool.acquire(process, service, connection, requestBytes, limits)
	if err != nil {
		return zero, err
	}
	if err := caller.Err(); err != nil {
		release()
		return zero, &passiveTransportError{disposition: "caller-cancelled-before-dispatch"}
	}
	methodContext, cancelMethod := context.WithTimeout(context.Background(), deadline)
	result := make(chan passiveCallResult[T], 1)
	go func() {
		defer cancelMethod()
		defer release()
		value, callErr := call(methodContext)
		result <- passiveCallResult[T]{value: value, err: callErr}
	}()
	select {
	case completed := <-result:
		return completed.value, completed.err
	case <-caller.Done():
		return zero, &passiveTransportError{disposition: "caller-cancelled-after-dispatch-response-unknown"}
	case <-methodContext.Done():
		return zero, &passiveTransportError{disposition: "method-deadline-after-dispatch-response-unknown"}
	}
}

func validateSubmitMethod(received SubmitMainMJSV0MethodRecord) error {
	expected := ExpectedSubmitMainMJSV0MethodRecord()
	if received.MethodVersion != expected.MethodVersion {
		return refused(Unsupported, "submit-main-mjs-method-version")
	}
	if received.Method != expected.Method {
		return refused(Unsupported, "submit-main-mjs-method")
	}
	if received != expected {
		return refused(Authentication, "submit-main-mjs-method-binding")
	}
	return nil
}

func flowLimitsForSubmit(method SubmitMainMJSV0MethodRecord) passiveFlowLimits {
	return passiveFlowLimits{
		serviceMaxConnections:       method.ServiceMaxConnections,
		connectionMaxInFlight:       method.ConnectionMaxInFlight,
		processMaxAdmittedRequests:  method.ProcessMaxAdmittedRequests,
		applicationQueueCapacity:    method.ApplicationQueueCapacity,
		inFlightRequestDataMaxBytes: method.InFlightRequestDataMaxBytes,
	}
}

func flowLimitsForSupervisor(service, connection, process, queue, bytes uint64) passiveFlowLimits {
	return passiveFlowLimits{
		serviceMaxConnections:       service,
		connectionMaxInFlight:       connection,
		processMaxAdmittedRequests:  process,
		applicationQueueCapacity:    queue,
		inFlightRequestDataMaxBytes: bytes,
	}
}
