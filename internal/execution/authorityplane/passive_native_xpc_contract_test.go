package authorityplane

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type nativeXPCFixtureField struct {
	Key             string             `json:"key"`
	ValueType       NativeXPCValueType `json:"valueType"`
	Required        bool               `json:"required"`
	MinDataBytes    uint64             `json:"minDataBytes"`
	MaxDataBytes    uint64             `json:"maxDataBytes"`
	FixedUInt64     *uint64            `json:"fixedUInt64"`
	FixedString     *string            `json:"fixedString"`
	AllowedUInt64   []uint64           `json:"allowedUInt64"`
	ApplicationData bool               `json:"applicationData"`
	NonZeroData     bool               `json:"nonZeroData"`
}

type nativeXPCFixtureEnvelope struct {
	Method                  string                  `json:"method"`
	Direction               string                  `json:"direction"`
	MessageTag              NativeXPCMessageTag     `json:"messageTag"`
	ProtocolVersion         uint64                  `json:"protocolVersion"`
	MethodVersion           uint64                  `json:"methodVersion"`
	ExactKeyCount           uint64                  `json:"exactKeyCount"`
	RequiredKeyCount        uint64                  `json:"requiredKeyCount"`
	OptionalKeyCount        uint64                  `json:"optionalKeyCount"`
	ClosedMap               bool                    `json:"closedMap"`
	ApplicationDataMaxBytes uint64                  `json:"applicationDataMaxBytes"`
	Fields                  []nativeXPCFixtureField `json:"fields"`
}

type nativeXPCFixtureMethod struct {
	Request      nativeXPCFixtureEnvelope `json:"request"`
	SuccessReply nativeXPCFixtureEnvelope `json:"successReply"`
	RefusalReply nativeXPCFixtureEnvelope `json:"refusalReply"`
}

type nativeXPCFixtureCase struct {
	ID                     string                     `json:"id"`
	Method                 string                     `json:"method"`
	Methods                string                     `json:"methods"`
	Direction              *string                    `json:"direction"`
	Mutation               *string                    `json:"mutation"`
	MismatchSet            *string                    `json:"mismatchSet"`
	SelectedFailure        *string                    `json:"selectedFailure"`
	ValidationPrecedence   []string                   `json:"validationPrecedence"`
	DispatchIdentity       *string                    `json:"dispatchIdentity"`
	DeliveryPrecondition   *string                    `json:"deliveryPrecondition"`
	Expected               json.RawMessage            `json:"expected"`
	ReplyDisposition       string                     `json:"replyDisposition"`
	TerminationDisposition string                     `json:"terminationDisposition"`
	NoState                *nativeXPCFixtureZeroState `json:"noState"`
	PostCoreState          *string                    `json:"postCoreState"`
}

type nativeXPCFixtureZeroState struct {
	AuthorityStateChanged bool   `json:"authorityStateChanged"`
	CoreCalls             uint64 `json:"coreCalls"`
	RegistrationsCreated  uint64 `json:"registrationsCreated"`
	ApprovalsConsumed     uint64 `json:"approvalsConsumed"`
	AttemptsCreated       uint64 `json:"attemptsCreated"`
	LifecycleCalls        uint64 `json:"lifecycleCalls"`
	BackendCalls          uint64 `json:"backendCalls"`
}

type nativeXPCFixtureCaseTableEntry struct {
	ID                     string                     `json:"id"`
	Method                 string                     `json:"method"`
	Direction              *string                    `json:"direction"`
	Mutation               *string                    `json:"mutation"`
	Decision               string                     `json:"decision"`
	Classification         *string                    `json:"classification"`
	StatusTag              *uint64                    `json:"statusTag"`
	ReasonTag              *uint64                    `json:"reasonTag"`
	BodyCopied             *bool                      `json:"bodyCopied"`
	CoreCalls              *uint64                    `json:"coreCalls"`
	ReplyDisposition       string                     `json:"replyDisposition"`
	TerminationDisposition string                     `json:"terminationDisposition"`
	NoState                *nativeXPCFixtureZeroState `json:"noState"`
	PostCoreState          *string                    `json:"postCoreState"`
}

type nativeXPCFixtureRefusalReply struct {
	Classification  string `json:"classification"`
	StatusTag       uint64 `json:"statusTag"`
	ReasonTag       uint64 `json:"reasonTag"`
	BodyKeysAllowed uint64 `json:"bodyKeysAllowed"`
	ExactKeyCount   uint64 `json:"exactKeyCount"`
}

type nativeXPCFixtureMethodBinding struct {
	EntryPoint                string              `json:"entryPoint"`
	Service                   string              `json:"service"`
	ExpectedRole              CallerRole          `json:"expectedRole"`
	ExpectedSigningIdentifier *string             `json:"expectedSigningIdentifier"`
	Audience                  string              `json:"audience"`
	Purpose                   string              `json:"purpose"`
	MessageTag                NativeXPCMessageTag `json:"messageTag"`
	MethodVersion             uint64              `json:"methodVersion"`
	DeadlineMilliseconds      uint64              `json:"deadlineMilliseconds"`
}

type nativeXPCFixtureResponseLoss struct {
	Method                    string `json:"method"`
	Disposition               string `json:"disposition"`
	SemanticIdentity          string `json:"semanticIdentity"`
	SameApprovalID            bool   `json:"sameApprovalID"`
	SameAttemptID             bool   `json:"sameAttemptID"`
	CurrentStateReturned      bool   `json:"currentStateReturned"`
	DuplicateAuthorityEffects uint64 `json:"duplicateAuthorityEffects"`
	DuplicateAttempts         uint64 `json:"duplicateAttempts"`
	DuplicateLifecycleEffects uint64 `json:"duplicateLifecycleEffects"`
	RequestIDIsIdempotencyKey bool   `json:"requestIdIsIdempotencyKey"`
}

type nativeXPCFixtureDeadlineCase struct {
	ID                       string                           `json:"id"`
	Method                   string                           `json:"method"`
	DeadlineMilliseconds     uint64                           `json:"deadlineMilliseconds"`
	StartsAt                 string                           `json:"startsAt"`
	ClientExtensionAllowed   bool                             `json:"clientExtensionAllowed"`
	BoundaryRule             string                           `json:"boundaryRule"`
	AfterDispatchDisposition string                           `json:"afterDispatchDisposition"`
	Boundary                 string                           `json:"boundary"`
	ElapsedMilliseconds      uint64                           `json:"elapsedMilliseconds"`
	Expected                 nativeXPCFixtureDeadlineExpected `json:"expected"`
	NoState                  *nativeXPCFixtureZeroState       `json:"noState"`
	PostCoreState            *string                          `json:"postCoreState"`
}

type nativeXPCFixtureDeadlineExpected struct {
	Decision         string  `json:"decision"`
	StatusTag        *uint64 `json:"statusTag"`
	ReasonTag        *uint64 `json:"reasonTag"`
	ReplyDisposition string  `json:"replyDisposition"`
	BodyCopied       bool    `json:"bodyCopied"`
	CoreCalls        uint64  `json:"coreCalls"`
}

type nativeXPCFixtureFutureHarnessOracle struct {
	ID              string `json:"id"`
	Scope           string `json:"scope"`
	Expected        string `json:"expected"`
	CurrentEvidence bool   `json:"currentEvidence"`
	InBandRefusal   bool   `json:"inBandRefusal"`
}

type nativeXPCFixtureDocument struct {
	ObjectType                string                                   `json:"objectType"`
	ObjectVersion             uint64                                   `json:"objectVersion"`
	Status                    string                                   `json:"status"`
	TransportEncoding         string                                   `json:"transportEncoding"`
	TransportEncodingVersion  uint64                                   `json:"transportEncodingVersion"`
	FixtureSerialization      string                                   `json:"fixtureSerialization"`
	TopLevelObjectType        string                                   `json:"topLevelObjectType"`
	ValueTypes                map[string]string                        `json:"valueTypes"`
	Keys                      map[string]string                        `json:"keys"`
	MessageTags               map[string]uint64                        `json:"messageTags"`
	ApprovalStateTags         map[string]uint64                        `json:"approvalStateTags"`
	AttemptStateTags          map[string]uint64                        `json:"attemptStateTags"`
	StatusTags                map[string]uint64                        `json:"statusTags"`
	ReasonTags                map[string]uint64                        `json:"reasonTags"`
	ClassificationToStatus    map[string]uint64                        `json:"classificationToStatus"`
	StructuralReasonToStatus  map[string]uint64                        `json:"structuralReasonToStatus"`
	CoreRefusalMapping        string                                   `json:"coreRefusalMapping"`
	LocalIntegrityMapping     string                                   `json:"localIntegrityMapping"`
	MethodBindings            map[string]nativeXPCFixtureMethodBinding `json:"methodBindings"`
	SignedObjectBindings      json.RawMessage                          `json:"signedObjectBindings"`
	IdentifierDomains         json.RawMessage                          `json:"identifierDomains"`
	RequestCommonKeyCount     uint64                                   `json:"requestCommonKeyCount"`
	ReplyCommonKeyCount       uint64                                   `json:"replyCommonKeyCount"`
	ExtraObjectsAllowed       uint64                                   `json:"extraObjectsAllowed"`
	FileDescriptorsAllowed    uint64                                   `json:"fileDescriptorsAllowed"`
	EndpointsAllowed          uint64                                   `json:"endpointsAllowed"`
	MachRightsAllowed         uint64                                   `json:"machRightsAllowed"`
	NestedContainersAllowed   uint64                                   `json:"nestedContainersAllowed"`
	MessageTagDisposition     string                                   `json:"messageTagDisposition"`
	ValidationPrecedence      []string                                 `json:"validationPrecedence"`
	RequestIDDisposition      string                                   `json:"requestIdDisposition"`
	CopyDisposition           string                                   `json:"copyDisposition"`
	LocalIntegrityDisposition string                                   `json:"localIntegrityDisposition"`
	Envelopes                 map[string]nativeXPCFixtureMethod        `json:"envelopes"`
	Cases                     []nativeXPCFixtureCase                   `json:"cases"`
	CaseTable                 []nativeXPCFixtureCaseTableEntry         `json:"caseTable"`
	RefusalReplies            []nativeXPCFixtureRefusalReply           `json:"refusalReplies"`
	ResponseLoss              []nativeXPCFixtureResponseLoss           `json:"responseLoss"`
	DeadlineBoundaryRule      string                                   `json:"deadlineBoundaryRule"`
	DeadlineCases             []nativeXPCFixtureDeadlineCase           `json:"deadlineCases"`
	FutureHarnessOracles      []nativeXPCFixtureFutureHarnessOracle    `json:"futureNativeHarnessOracles"`
	PeerAuthentication        any                                      `json:"peerAuthenticationEvidence"`
	ListenerActivated         bool                                     `json:"listenerActivated"`
	ServiceRegistered         bool                                     `json:"serviceRegistered"`
	CasesDigest               string                                   `json:"-"`
}

func TestNativeXPCContractMatchesGeneratedCrossLanguageFixture(t *testing.T) {
	bytes, err := os.ReadFile(filepath.Join("..", "..", "..", "schemas", "conformance", "authenticated-local-ipc-v0", "native-xpc-v0.contract.json"))
	if err != nil {
		t.Fatal(err)
	}
	fixture, err := decodeNativeXPCFixture(bytes)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateHardenedNativeXPCFixture(fixture); err != nil {
		t.Fatal(err)
	}
	if fixture.ObjectType != "capsule.authenticated-local-ipc-native-xpc-contract" || fixture.ObjectVersion != 0 || fixture.Status != "passive-unwired-no-listener" {
		t.Fatalf("unexpected native fixture identity: %#v", fixture)
	}
	if fixture.TransportEncoding != NativeXPCEncoding || fixture.TransportEncodingVersion != NativeXPCEncodingVersion {
		t.Fatalf("native encoding mismatch: %q v%d", fixture.TransportEncoding, fixture.TransportEncodingVersion)
	}
	if fixture.PeerAuthentication != nil || fixture.ListenerActivated || fixture.ServiceRegistered {
		t.Fatal("passive native contract must not claim peer authentication, listener activation, or service registration")
	}
	wantPrecedence := []string{
		"protocol-version", "method-version",
		"service-entry-point-role-message-tag-audience-purpose",
		"installation-epoch-current-state", "application-data-copy",
		"embedded-record-version-and-core-validation",
	}
	if !reflect.DeepEqual(fixture.ValidationPrecedence, wantPrecedence) {
		t.Fatalf("native validation precedence mismatch: got %#v want %#v", fixture.ValidationPrecedence, wantPrecedence)
	}
	wantTags := map[string]uint64{
		"invalid":                 0,
		SubmitMainMJSV0Method:     uint64(NativeXPCMessageTagSubmitMainMJSV0),
		RegisterPlanV0Method:      uint64(NativeXPCMessageTagRegisterPlanV0),
		GetRegisteredPlanV0Method: uint64(NativeXPCMessageTagGetRegisteredPlanV0),
		SubmitApprovalV0Method:    uint64(NativeXPCMessageTagSubmitApprovalV0),
		RequestAttemptV0Method:    uint64(NativeXPCMessageTagRequestAttemptV0),
	}
	if !reflect.DeepEqual(fixture.MessageTags, wantTags) {
		t.Fatalf("native message tags mismatch: got %#v want %#v", fixture.MessageTags, wantTags)
	}
	wantApprovalStates := map[string]uint64{
		"invalid": uint64(NativeXPCApprovalStateInvalid), "usable": uint64(NativeXPCApprovalStateUsable),
		"consumed": uint64(NativeXPCApprovalStateConsumed), "invalidated": uint64(NativeXPCApprovalStateInvalidated),
	}
	wantAttemptStates := map[string]uint64{
		"invalid": uint64(NativeXPCAttemptStateInvalid), "created": uint64(NativeXPCAttemptStateCreated),
	}
	if !reflect.DeepEqual(fixture.ApprovalStateTags, wantApprovalStates) || !reflect.DeepEqual(fixture.AttemptStateTags, wantAttemptStates) {
		t.Fatalf("native state tags mismatch: approval=%#v attempt=%#v", fixture.ApprovalStateTags, fixture.AttemptStateTags)
	}
	wantStatuses := map[string]uint64{
		"OK": uint64(NativeXPCStatusOK), "MALFORMED": uint64(NativeXPCStatusMalformed),
		"UNSUPPORTED": uint64(NativeXPCStatusUnsupported), "SCHEMA": uint64(NativeXPCStatusSchema),
		"BINDING": uint64(NativeXPCStatusBinding), "AUTHENTICATION": uint64(NativeXPCStatusAuthentication),
		"STALE": uint64(NativeXPCStatusStale), "REPLAY": uint64(NativeXPCStatusReplay),
		"CAPACITY": uint64(NativeXPCStatusCapacity), "TRUST_STATE": uint64(NativeXPCStatusTrustState),
		"LOCAL_FAILURE": uint64(NativeXPCStatusLocalFailure), "RECOVERY_REQUIRED": uint64(NativeXPCStatusRecoveryRequired),
		"SEMANTIC": uint64(NativeXPCStatusSemantic), "DOMAIN": uint64(NativeXPCStatusDomain),
	}
	if !reflect.DeepEqual(fixture.StatusTags, wantStatuses) {
		t.Fatalf("native status tags mismatch: got %#v want %#v", fixture.StatusTags, wantStatuses)
	}
	for _, mapping := range NativeXPCErrorMappings() {
		if got := fixture.ClassificationToStatus[string(mapping.Classification)]; got != uint64(mapping.StatusTag) {
			t.Fatalf("classification %s maps to %d, want %d", mapping.Classification, got, mapping.StatusTag)
		}
	}
	if fixture.ReasonTags["none"] != uint64(NativeXPCReasonNone) || fixture.ReasonTags["localIntegrityFault"] != uint64(NativeXPCReasonLocalIntegrityFault) {
		t.Fatalf("native reason tags mismatch: %#v", fixture.ReasonTags)
	}
	reasonNames := map[NativeXPCReasonTag]string{
		NativeXPCReasonKeySet: "keySet", NativeXPCReasonValueType: "valueType",
		NativeXPCReasonDataWidth: "dataWidth", NativeXPCReasonDataCap: "dataCap",
		NativeXPCReasonZeroIdentifier: "zeroIdentifier", NativeXPCReasonEpochSequence: "epochSequence",
		NativeXPCReasonProtocolVersion: "protocolVersion", NativeXPCReasonMethodVersion: "methodVersion",
		NativeXPCReasonMessageTag: "messageTag", NativeXPCReasonMethodBinding: "methodBinding",
		NativeXPCReasonCurrentState: "currentState", NativeXPCReasonCapacity: "capacity",
	}
	for _, mapping := range NativeXPCStructuralRefusalMappings() {
		name := reasonNames[mapping.ReasonTag]
		if fixture.ReasonTags[name] != uint64(mapping.ReasonTag) || fixture.StructuralReasonToStatus[name] != uint64(mapping.StatusTag) {
			t.Fatalf("native structural reason mapping mismatch for %s", name)
		}
	}
	for _, expected := range ExpectedNativeXPCMethodBindings() {
		received, ok := fixture.MethodBindings[expected.Method]
		if !ok {
			t.Fatalf("missing native method binding %s", expected.Method)
		}
		receivedSigningID := ""
		if received.ExpectedSigningIdentifier != nil {
			receivedSigningID = *received.ExpectedSigningIdentifier
		}
		if received.EntryPoint != expected.EntryPoint || received.Service != expected.Service || received.ExpectedRole != expected.ExpectedRole || receivedSigningID != expected.ExpectedSigningIdentifier || received.Audience != expected.Audience || received.Purpose != expected.Purpose || received.MessageTag != expected.MessageTag || received.MethodVersion != expected.MethodVersion || received.DeadlineMilliseconds != expected.DeadlineMilliseconds {
			t.Fatalf("native method binding mismatch for %s: got %#v want %#v", expected.Method, received, expected)
		}
	}

	goMethods := map[string][3]NativeXPCEnvelopeSpec{
		SubmitMainMJSV0Method: {
			ExpectedNativeXPCSubmitMainMJSV0Request(),
			ExpectedNativeXPCSubmitMainMJSV0Reply(),
			ExpectedNativeXPCSubmitMainMJSV0RefusalReply(),
		},
		RegisterPlanV0Method: {
			ExpectedNativeXPCRegisterPlanV0Request(),
			ExpectedNativeXPCRegisterPlanV0Reply(),
			ExpectedNativeXPCRegisterPlanV0RefusalReply(),
		},
		GetRegisteredPlanV0Method: {
			ExpectedNativeXPCGetRegisteredPlanV0Request(),
			ExpectedNativeXPCGetRegisteredPlanV0Reply(),
			ExpectedNativeXPCGetRegisteredPlanV0RefusalReply(),
		},
		SubmitApprovalV0Method: {
			ExpectedNativeXPCSubmitApprovalV0Request(),
			ExpectedNativeXPCSubmitApprovalV0Reply(),
			ExpectedNativeXPCSubmitApprovalV0RefusalReply(),
		},
		RequestAttemptV0Method: {
			ExpectedNativeXPCRequestAttemptV0Request(),
			ExpectedNativeXPCRequestAttemptV0Reply(),
			ExpectedNativeXPCRequestAttemptV0RefusalReply(),
		},
	}
	for method, goEnvelopes := range goMethods {
		fixtureMethod, ok := fixture.Envelopes[method]
		if !ok {
			t.Fatalf("missing fixture method %s", method)
		}
		fixtureEnvelopes := [3]nativeXPCFixtureEnvelope{fixtureMethod.Request, fixtureMethod.SuccessReply, fixtureMethod.RefusalReply}
		for index := range goEnvelopes {
			assertNativeXPCEnvelope(t, goEnvelopes[index], fixtureEnvelopes[index])
		}
	}
	wantCaseIDs := []string{
		"all.exact-key-and-type-sets", "submit-approval.request.exact-maximum",
		"request-attempt.request.exact-maximum", "submit.body-cap-plus-one",
		"register.body-cap-plus-one", "get.body-cap-plus-one",
		"submit-approval.body-cap-plus-one", "request-attempt.approval-id-width-plus-one",
		"all.missing-key",
		"all.extra-key", "all.wrong-type", "all.nested-container", "all.zero-request-id",
		"all.unknown-protocol-version", "all.unknown-method-version",
		"register.joint-protocol-method-record-version-mismatch",
		"register.joint-method-record-version-mismatch", "register.embedded-record-version-mismatch",
		"submit-approval.joint-protocol-method-record-version-mismatch",
		"submit-approval.joint-method-record-version-mismatch",
		"submit-approval.embedded-record-version-mismatch",
	}
	orderedMethods := []string{
		SubmitMainMJSV0Method, RegisterPlanV0Method, GetRegisteredPlanV0Method,
		SubmitApprovalV0Method, RequestAttemptV0Method,
	}
	for _, selected := range orderedMethods {
		for _, received := range orderedMethods {
			if selected != received {
				wantCaseIDs = append(wantCaseIDs, "cross-service."+selected+".tag-"+received)
			}
		}
	}
	wantCaseIDs = append(wantCaseIDs,
		"all.wrong-service", "all.wrong-role", "all.wrong-session", "all.wrong-audience", "all.wrong-purpose",
		"all.local-audience-replaced-by-signed-object-audience",
		"all.local-purpose-replaced-by-signed-object-purpose",
		"submit-approval.signed-audience-replaced-by-local-channel-audience",
		"submit-approval.signed-purpose-replaced-by-local-channel-purpose",
		"all.request-id-replaced-by-installation-id-bytes",
		"all.installation-id-replaced-by-request-id-bytes",
		"get.registration-id-replaced-by-request-id-bytes",
		"submit-approval.registration-id-replaced-by-request-id-bytes",
		"request-attempt.registration-id-replaced-by-approval-id-bytes",
		"request-attempt.approval-id-replaced-by-registration-id-bytes",
		"all.wrong-installation",
		"all.wrong-epoch-digest", "all.epoch-uint53-cap-plus-one",
		"all.extra-file-descriptor", "all.extra-endpoint", "all.extra-mach-send-right",
		"all.success-reply-extra-key", "all.success-reply-request-id-mismatch",
		"all.refusal-reply-extra-body", "all.local-integrity-output-fault",
		"submit-approval.cancel-before-dispatch", "submit-approval.cancel-after-dispatch",
		"request-attempt.cancel-before-dispatch", "request-attempt.cancel-after-dispatch",
	)
	gotCaseIDs := make([]string, len(fixture.Cases))
	derivedCaseTable := make([]nativeXPCFixtureCaseTableEntry, len(fixture.Cases))
	for index, candidate := range fixture.Cases {
		gotCaseIDs[index] = candidate.ID
		method := candidate.Method
		if method == "" {
			method = candidate.Methods
		}
		derived := nativeXPCFixtureCaseTableEntry{
			ID: candidate.ID, Method: method, Direction: candidate.Direction, Mutation: candidate.Mutation,
			ReplyDisposition: candidate.ReplyDisposition, TerminationDisposition: candidate.TerminationDisposition,
			NoState: candidate.NoState, PostCoreState: candidate.PostCoreState,
		}
		if len(candidate.Expected) == 0 {
			t.Fatalf("native case %s has no expected result", candidate.ID)
		}
		if candidate.Expected[0] == '"' {
			if err := json.Unmarshal(candidate.Expected, &derived.Decision); err != nil {
				t.Fatalf("decode native case %s decision: %v", candidate.ID, err)
			}
			derivedCaseTable[index] = derived
			continue
		}
		var expected struct {
			Decision       string  `json:"decision"`
			Classification *string `json:"classification"`
			StatusTag      *uint64 `json:"statusTag"`
			ReasonTag      *uint64 `json:"reasonTag"`
			BodyCopied     *bool   `json:"bodyCopied"`
			CoreCalls      *uint64 `json:"coreCalls"`
		}
		if err := json.Unmarshal(candidate.Expected, &expected); err != nil {
			t.Fatalf("decode native case %s: %v", candidate.ID, err)
		}
		if expected.StatusTag != nil && expected.Classification != nil && *expected.StatusTag != fixture.StatusTags[*expected.Classification] {
			t.Fatalf("native case %s status/classification mismatch", candidate.ID)
		}
		if expected.ReasonTag != nil && *expected.ReasonTag > uint64(NativeXPCReasonLocalIntegrityFault) {
			t.Fatalf("native case %s uses unknown reason %d", candidate.ID, *expected.ReasonTag)
		}
		derived.Decision = expected.Decision
		derived.Classification = expected.Classification
		derived.StatusTag = expected.StatusTag
		derived.ReasonTag = expected.ReasonTag
		derived.BodyCopied = expected.BodyCopied
		derived.CoreCalls = expected.CoreCalls
		derivedCaseTable[index] = derived
	}
	if !reflect.DeepEqual(gotCaseIDs, wantCaseIDs) {
		t.Fatalf("native case order/set mismatch:\n got %#v\nwant %#v", gotCaseIDs, wantCaseIDs)
	}
	if !reflect.DeepEqual(fixture.CaseTable, derivedCaseTable) {
		t.Fatalf("native case table does not exactly project the complete ordered case set:\n got %#v\nwant %#v", fixture.CaseTable, derivedCaseTable)
	}
	wantResponseLoss := []nativeXPCFixtureResponseLoss{
		{Method: SubmitMainMJSV0Method, Disposition: SubmitMainMJSV0ResponseLossDisposition},
		{Method: RegisterPlanV0Method, Disposition: RegisterPlanV0ResponseLossDisposition},
		{Method: GetRegisteredPlanV0Method, Disposition: GetRegisteredPlanV0ResponseLossDisposition},
		{
			Method: SubmitApprovalV0Method, Disposition: SubmitApprovalV0ResponseLossDisposition,
			SemanticIdentity: "canonical-approval-payload+resolved-signer-authorization-identity",
			SameApprovalID:   true, CurrentStateReturned: true, DuplicateAuthorityEffects: 0,
			RequestIDIsIdempotencyKey: false,
		},
		{
			Method: RequestAttemptV0Method, Disposition: RequestAttemptV0ResponseLossDisposition,
			SemanticIdentity: "registration-id+approval-reference",
			SameAttemptID:    true, CurrentStateReturned: true, DuplicateAttempts: 0,
			DuplicateLifecycleEffects: 0, RequestIDIsIdempotencyKey: false,
		},
	}
	if !reflect.DeepEqual(fixture.ResponseLoss, wantResponseLoss) {
		t.Fatalf("native response-loss table mismatch: got %#v want %#v", fixture.ResponseLoss, wantResponseLoss)
	}
	wantDeadlines := expectedNativeXPCDeadlineCases()
	if !reflect.DeepEqual(fixture.DeadlineCases, wantDeadlines) {
		t.Fatalf("native deadline table mismatch: got %#v want %#v", fixture.DeadlineCases, wantDeadlines)
	}
	wantFutureOracles := []nativeXPCFixtureFutureHarnessOracle{{
		ID: "os-peer-requirement-mismatch", Scope: "future-external-native-harness-only",
		Expected: "no-message-delivery-and-no-application-reply", CurrentEvidence: false,
		InBandRefusal: false,
	}}
	if !reflect.DeepEqual(fixture.FutureHarnessOracles, wantFutureOracles) {
		t.Fatalf("native future-harness oracle mismatch: got %#v want %#v", fixture.FutureHarnessOracles, wantFutureOracles)
	}
	if len(fixture.RefusalReplies) != len(NativeXPCErrorMappings()) {
		t.Fatalf("native refusal reply count: got %d want %d", len(fixture.RefusalReplies), len(NativeXPCErrorMappings()))
	}
	for _, reply := range fixture.RefusalReplies {
		if reply.StatusTag != fixture.ClassificationToStatus[reply.Classification] || reply.ReasonTag != uint64(NativeXPCReasonCoreRefusal) || reply.BodyKeysAllowed != 0 || reply.ExactKeyCount != NativeXPCRefusalReplyKeyCount {
			t.Fatalf("invalid native refusal reply %#v", reply)
		}
	}
}

func assertNativeXPCEnvelope(t *testing.T, expected NativeXPCEnvelopeSpec, received nativeXPCFixtureEnvelope) {
	t.Helper()
	if received.Method != expected.Method || received.Direction != expected.Direction || received.MessageTag != expected.MessageTag || received.ProtocolVersion != expected.ProtocolVersion || received.MethodVersion != expected.MethodVersion || received.ExactKeyCount != expected.ExactKeyCount || received.RequiredKeyCount != expected.RequiredKeyCount || received.OptionalKeyCount != expected.OptionalKeyCount || received.ClosedMap != expected.ClosedMap || received.ApplicationDataMaxBytes != expected.ApplicationDataMaxBytes {
		t.Fatalf("native envelope mismatch: got %#v want %#v", received, expected)
	}
	if uint64(len(expected.Fields)) != expected.ExactKeyCount || len(received.Fields) != len(expected.Fields) {
		t.Fatalf("native envelope field-count mismatch for %s %s", expected.Method, expected.Direction)
	}
	seen := make(map[string]struct{}, len(expected.Fields))
	var derivedApplicationDataMaxBytes uint64
	for index, expectedField := range expected.Fields {
		receivedField := received.Fields[index]
		if _, exists := seen[expectedField.Key]; exists {
			t.Fatalf("duplicate native key %q", expectedField.Key)
		}
		for _, value := range []byte(expectedField.Key) {
			if value < 0x21 || value > 0x7e {
				t.Fatalf("native key is not visible ASCII: %q", expectedField.Key)
			}
		}
		seen[expectedField.Key] = struct{}{}
		if expectedField.ApplicationData {
			derivedApplicationDataMaxBytes += expectedField.MaxDataBytes
		}
		if receivedField.Key != expectedField.Key || receivedField.ValueType != expectedField.ValueType || receivedField.Required != expectedField.Required || receivedField.MinDataBytes != expectedField.MinDataBytes || receivedField.MaxDataBytes != expectedField.MaxDataBytes || !reflect.DeepEqual(receivedField.AllowedUInt64, expectedField.AllowedUInt64) || receivedField.ApplicationData != expectedField.ApplicationData || receivedField.NonZeroData != expectedField.NonZeroData {
			t.Fatalf("native field mismatch for %s %s: got %#v want %#v", expected.Method, expected.Direction, receivedField, expectedField)
		}
		if expectedField.HasFixedUInt64 {
			if receivedField.FixedUInt64 == nil || *receivedField.FixedUInt64 != expectedField.FixedUInt64 {
				t.Fatalf("native fixed uint mismatch for %s", expectedField.Key)
			}
		} else if receivedField.FixedUInt64 != nil {
			t.Fatalf("unexpected native fixed uint for %s", expectedField.Key)
		}
		if expectedField.HasFixedString {
			if receivedField.FixedString == nil || *receivedField.FixedString != expectedField.FixedString {
				t.Fatalf("native fixed string mismatch for %s", expectedField.Key)
			}
		} else if receivedField.FixedString != nil {
			t.Fatalf("unexpected native fixed string for %s", expectedField.Key)
		}
	}
	if derivedApplicationDataMaxBytes != expected.ApplicationDataMaxBytes {
		t.Fatalf("native envelope application-data aggregate mismatch for %s %s: got %d want %d", expected.Method, expected.Direction, expected.ApplicationDataMaxBytes, derivedApplicationDataMaxBytes)
	}
}

func TestNativeXPCContractHasNoGenericDispatchOrAuthorityFields(t *testing.T) {
	for _, spec := range []NativeXPCEnvelopeSpec{
		ExpectedNativeXPCSubmitMainMJSV0Request(),
		ExpectedNativeXPCRegisterPlanV0Request(),
		ExpectedNativeXPCGetRegisteredPlanV0Request(),
		ExpectedNativeXPCSubmitApprovalV0Request(),
		ExpectedNativeXPCRequestAttemptV0Request(),
	} {
		for _, field := range spec.Fields {
			switch field.Key {
			case "capsule.method", "capsule.command", "capsule.opcode", "capsule.role", "capsule.service", "capsule.backend", "capsule.path", "capsule.image", "capsule.mount":
				t.Fatalf("forbidden generic or authority-bearing native field %q", field.Key)
			}
		}
	}
}

func TestNativeXPCHardenedVerifierRejectsBoundedMutations(t *testing.T) {
	path := filepath.Join("..", "..", "..", "schemas", "conformance", "authenticated-local-ipc-v0", "native-xpc-v0.contract.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var baseline map[string]any
	if err := json.Unmarshal(data, &baseline); err != nil {
		t.Fatal(err)
	}
	mutations := []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"missing-dictionary-field", func(root map[string]any) {
			fields := nativeXPCMutationFields(root, SubmitApprovalV0Method, "request")
			nativeXPCMutationEnvelope(root, SubmitApprovalV0Method, "request")["fields"] = fields[:len(fields)-1]
		}},
		{"extra-dictionary-field", func(root map[string]any) {
			fields := nativeXPCMutationFields(root, RequestAttemptV0Method, "request")
			nativeXPCMutationEnvelope(root, RequestAttemptV0Method, "request")["fields"] = append(fields, map[string]any{"key": "capsule.extra", "valueType": "XPC_TYPE_DATA"})
		}},
		{"changed-type-or-width", func(root map[string]any) {
			fields := nativeXPCMutationFields(root, SubmitApprovalV0Method, "request")
			fields[len(fields)-1].(map[string]any)["maxDataBytes"] = float64(511)
		}},
		{"absent-required-no-state", func(root map[string]any) {
			delete(nativeXPCMutationCase(root, "all.missing-key"), "noState")
		}},
		{"modified-cancellation-expected-result", func(root map[string]any) {
			nativeXPCMutationCase(root, "submit-approval.cancel-after-dispatch")["expected"].(map[string]any)["decision"] = "no-core-call-no-application-reply"
		}},
		{"missing-refusal-case", func(root map[string]any) {
			cases := root["cases"].([]any)
			for index, candidate := range cases {
				if candidate.(map[string]any)["id"] == "all.wrong-role" {
					root["cases"] = append(cases[:index], cases[index+1:]...)
					return
				}
			}
		}},
		{"response-loss-table-drift", func(root map[string]any) {
			root["responseLoss"].([]any)[3].(map[string]any)["sameApprovalID"] = false
		}},
		{"exact-at-boundary-inversion", func(root map[string]any) {
			for _, candidate := range root["deadlineCases"].([]any) {
				entry := candidate.(map[string]any)
				if entry["id"] == "SubmitApprovalV0.deadline.exactly-at" {
					entry["expected"].(map[string]any)["decision"] = "dispatch-core-before-deadline"
					return
				}
			}
		}},
	}
	for _, mutation := range mutations {
		t.Run(mutation.name, func(t *testing.T) {
			clonedBytes, err := json.Marshal(baseline)
			if err != nil {
				t.Fatal(err)
			}
			var changed map[string]any
			if err := json.Unmarshal(clonedBytes, &changed); err != nil {
				t.Fatal(err)
			}
			mutation.mutate(changed)
			changedBytes, err := json.Marshal(changed)
			if err != nil {
				t.Fatal(err)
			}
			fixture, decodeErr := decodeNativeXPCFixture(changedBytes)
			if decodeErr == nil {
				decodeErr = validateHardenedNativeXPCFixture(fixture)
			}
			if decodeErr == nil {
				t.Fatal("hardened native XPC verifier accepted mutation")
			}
		})
	}
}

func TestNativeXPCOraclesAreCompleteAndClosed(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "..", "schemas", "conformance", "authenticated-local-ipc-v0", "oracles.json"))
	if err != nil {
		t.Fatal(err)
	}
	var oracle map[string]json.RawMessage
	if err := json.Unmarshal(data, &oracle); err != nil {
		t.Fatal(err)
	}
	wantKeys := []string{"objectType", "objectVersion", "maxima", "refusals", "copyOwnership", "responseLoss", "flowControl", "cancellationAndDeadline", "sourceCustody", "zeroEffects"}
	if len(oracle) != len(wantKeys) {
		t.Fatalf("oracle top-level key count: got %d want %d", len(oracle), len(wantKeys))
	}
	for _, key := range wantKeys {
		if _, ok := oracle[key]; !ok {
			t.Fatalf("oracle missing top-level key %q", key)
		}
	}
	wantDigests := map[string]string{
		"maxima":                  "75d701adfa713008b0b21de1a754e3999448443d3d4794fdeac0dc0dab7af11e",
		"refusals":                "314390ececfbef518bce7cf84ce0966f68905b1fb4837ee100426150bf1ed07d",
		"responseLoss":            "ab064f9f00898f9164a03b7839f4bdd9582ee05cf3a45c12fd26596b6859c297",
		"flowControl":             "3db47429084a43b0e265400c2b53feac1d509b0179a6089bd7860449a194488e",
		"cancellationAndDeadline": "0b00d6b05e125c53f57bc8cdaacfa7c17bfec0e3fed96463f633ebcb58fe47af",
	}
	for section, want := range wantDigests {
		var compact bytes.Buffer
		if err := json.Compact(&compact, oracle[section]); err != nil {
			t.Fatalf("compact oracle %s: %v", section, err)
		}
		if got := fmt.Sprintf("%x", sha256.Sum256(compact.Bytes())); got != want {
			t.Fatalf("complete oracle %s drift: got %s want %s", section, got, want)
		}
	}
}

func nativeXPCMutationEnvelope(root map[string]any, method, direction string) map[string]any {
	return root["envelopes"].(map[string]any)[method].(map[string]any)[direction].(map[string]any)
}

func nativeXPCMutationFields(root map[string]any, method, direction string) []any {
	return nativeXPCMutationEnvelope(root, method, direction)["fields"].([]any)
}

func nativeXPCMutationCase(root map[string]any, id string) map[string]any {
	for _, candidate := range root["cases"].([]any) {
		entry := candidate.(map[string]any)
		if entry["id"] == id {
			return entry
		}
	}
	panic("missing mutation case " + id)
}

func decodeNativeXPCFixture(data []byte) (nativeXPCFixtureDocument, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nativeXPCFixtureDocument{}, err
	}
	wantTopLevel := []string{
		"objectType", "objectVersion", "status", "transportEncoding", "transportEncodingVersion",
		"fixtureSerialization", "topLevelObjectType", "valueTypes", "keys", "messageTags",
		"approvalStateTags", "attemptStateTags", "statusTags", "reasonTags", "classificationToStatus",
		"structuralReasonToStatus", "coreRefusalMapping", "localIntegrityMapping", "methodBindings",
		"signedObjectBindings", "identifierDomains", "requestCommonKeyCount", "replyCommonKeyCount",
		"extraObjectsAllowed", "fileDescriptorsAllowed", "endpointsAllowed", "machRightsAllowed",
		"nestedContainersAllowed", "messageTagDisposition", "validationPrecedence", "requestIdDisposition",
		"copyDisposition", "localIntegrityDisposition", "envelopes", "cases", "caseTable",
		"refusalReplies", "responseLoss", "deadlineBoundaryRule", "deadlineCases",
		"futureNativeHarnessOracles", "peerAuthenticationEvidence", "listenerActivated", "serviceRegistered",
	}
	if len(raw) != len(wantTopLevel) {
		return nativeXPCFixtureDocument{}, fmt.Errorf("native XPC top-level key count: got %d want %d", len(raw), len(wantTopLevel))
	}
	for _, key := range wantTopLevel {
		if _, ok := raw[key]; !ok {
			return nativeXPCFixtureDocument{}, fmt.Errorf("native XPC missing top-level key %q", key)
		}
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var fixture nativeXPCFixtureDocument
	if err := decoder.Decode(&fixture); err != nil {
		return nativeXPCFixtureDocument{}, err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return nativeXPCFixtureDocument{}, fmt.Errorf("native XPC trailing JSON: %v", err)
	}
	var canonicalCases any
	if err := json.Unmarshal(raw["cases"], &canonicalCases); err != nil {
		return nativeXPCFixtureDocument{}, err
	}
	canonicalCaseBytes, err := json.Marshal(canonicalCases)
	if err != nil {
		return nativeXPCFixtureDocument{}, err
	}
	fixture.CasesDigest = fmt.Sprintf("%x", sha256.Sum256(canonicalCaseBytes))
	return fixture, nil
}

func nativeXPCCanonicalDigest(raw json.RawMessage) (string, error) {
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return "", err
	}
	canonical, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(canonical)), nil
}

func validateHardenedNativeXPCFixture(fixture nativeXPCFixtureDocument) error {
	if fixture.ObjectType != "capsule.authenticated-local-ipc-native-xpc-contract" || fixture.ObjectVersion != 0 || fixture.Status != "passive-unwired-no-listener" || fixture.TransportEncoding != NativeXPCEncoding || fixture.TransportEncodingVersion != NativeXPCEncodingVersion {
		return fmt.Errorf("native XPC identity or encoding mismatch")
	}
	if fixture.FixtureSerialization != "exact-json-description-of-xpc-dictionaries" || fixture.TopLevelObjectType != "XPC_TYPE_DICTIONARY" {
		return fmt.Errorf("native XPC dictionary serialization mismatch")
	}
	if !reflect.DeepEqual(fixture.ValueTypes, map[string]string{"uint64": string(NativeXPCTypeUInt64), "data": string(NativeXPCTypeData), "string": string(NativeXPCTypeString)}) {
		return fmt.Errorf("native XPC complete value-type map mismatch")
	}
	wantKeys := map[string]string{
		"protocolVersion": NativeXPCProtocolVersionKey, "methodVersion": NativeXPCMethodVersionKey,
		"messageTag": NativeXPCMessageTagKey, "requestId": NativeXPCRequestIDKey,
		"installationId": NativeXPCInstallationIDKey, "epochSequence": NativeXPCEpochSequenceKey,
		"epochDigest": NativeXPCEpochDigestKey, "audience": NativeXPCAudienceKey,
		"purpose": NativeXPCPurposeKey, "status": NativeXPCStatusKey, "reason": NativeXPCReasonKey,
		"jobProposal": NativeXPCJobProposalKey, "executionPlan": NativeXPCExecutionPlanKey,
		"roleBindings": NativeXPCRoleBindingsKey, "sourceManifest": NativeXPCSourceManifestKey,
		"source": NativeXPCSourceKey, "registrationId": NativeXPCRegistrationIDKey,
		"planRegistration": NativeXPCPlanRegistrationKey, "approvalEnvelope": NativeXPCApprovalEnvelopeKey,
		"approvalId": NativeXPCApprovalIDKey, "approvalState": NativeXPCApprovalStateKey,
		"attemptId": NativeXPCAttemptIDKey, "attemptState": NativeXPCAttemptStateKey,
	}
	if !reflect.DeepEqual(fixture.Keys, wantKeys) {
		return fmt.Errorf("native XPC complete key map mismatch")
	}
	wantMessageTags := map[string]uint64{
		"invalid": 0, SubmitMainMJSV0Method: 1, RegisterPlanV0Method: 2,
		GetRegisteredPlanV0Method: 3, SubmitApprovalV0Method: 4, RequestAttemptV0Method: 5,
	}
	wantApprovalStates := map[string]uint64{"invalid": 0, "usable": 1, "consumed": 2, "invalidated": 3}
	wantAttemptStates := map[string]uint64{"invalid": 0, "created": 1}
	wantStatuses := map[string]uint64{
		"OK": 0, "MALFORMED": 1, "UNSUPPORTED": 2, "SCHEMA": 3, "BINDING": 4,
		"AUTHENTICATION": 5, "STALE": 6, "REPLAY": 7, "CAPACITY": 8, "TRUST_STATE": 9,
		"LOCAL_FAILURE": 10, "RECOVERY_REQUIRED": 11, "SEMANTIC": 12, "DOMAIN": 13,
	}
	wantReasons := map[string]uint64{
		"none": 0, "keySet": 1, "valueType": 2, "dataWidth": 3, "dataCap": 4,
		"zeroIdentifier": 5, "epochSequence": 6, "protocolVersion": 7, "methodVersion": 8,
		"messageTag": 9, "methodBinding": 10, "currentState": 11, "capacity": 12,
		"coreRefusal": 13, "localIntegrityFault": 14,
	}
	wantClassification := map[string]uint64{
		"MALFORMED": 1, "UNSUPPORTED": 2, "SCHEMA": 3, "BINDING": 4, "AUTHENTICATION": 5,
		"STALE": 6, "REPLAY": 7, "CAPACITY": 8, "TRUST_STATE": 9, "LOCAL_FAILURE": 10,
		"RECOVERY_REQUIRED": 11, "SEMANTIC": 12, "DOMAIN": 13,
	}
	wantStructural := map[string]uint64{
		"keySet": 1, "valueType": 1, "dataWidth": 3, "dataCap": 1, "zeroIdentifier": 3,
		"epochSequence": 3, "protocolVersion": 2, "methodVersion": 2, "messageTag": 2,
		"methodBinding": 5, "currentState": 4, "capacity": 8,
	}
	if !reflect.DeepEqual(fixture.MessageTags, wantMessageTags) || !reflect.DeepEqual(fixture.ApprovalStateTags, wantApprovalStates) || !reflect.DeepEqual(fixture.AttemptStateTags, wantAttemptStates) || !reflect.DeepEqual(fixture.StatusTags, wantStatuses) || !reflect.DeepEqual(fixture.ReasonTags, wantReasons) || !reflect.DeepEqual(fixture.ClassificationToStatus, wantClassification) || !reflect.DeepEqual(fixture.StructuralReasonToStatus, wantStructural) {
		return fmt.Errorf("native XPC complete numeric tag or state table mismatch")
	}
	if len(fixture.MethodBindings) != 5 {
		return fmt.Errorf("native XPC method-binding count: got %d want 5", len(fixture.MethodBindings))
	}
	for _, expected := range ExpectedNativeXPCMethodBindings() {
		received, ok := fixture.MethodBindings[expected.Method]
		if !ok {
			return fmt.Errorf("native XPC missing method binding %s", expected.Method)
		}
		receivedSigningID := ""
		if received.ExpectedSigningIdentifier != nil {
			receivedSigningID = *received.ExpectedSigningIdentifier
		}
		if received.EntryPoint != expected.EntryPoint || received.Service != expected.Service || received.ExpectedRole != expected.ExpectedRole || receivedSigningID != expected.ExpectedSigningIdentifier || received.Audience != expected.Audience || received.Purpose != expected.Purpose || received.MessageTag != expected.MessageTag || received.MethodVersion != expected.MethodVersion || received.DeadlineMilliseconds != expected.DeadlineMilliseconds {
			return fmt.Errorf("native XPC complete method binding mismatch for %s", expected.Method)
		}
	}
	if digest, err := nativeXPCCanonicalDigest(fixture.SignedObjectBindings); err != nil || digest != "c3941c3514bf19d094e24246472f321ceeb699e06a4804ac39cdfbbe15057f1d" {
		return fmt.Errorf("native XPC signed-object binding map mismatch: digest=%s err=%v", digest, err)
	}
	if digest, err := nativeXPCCanonicalDigest(fixture.IdentifierDomains); err != nil || digest != "45a196eae0f0871efae0d427c466a79c89b179cec24b40c4ae9a223c2d1f963d" {
		return fmt.Errorf("native XPC identifier-domain map mismatch: digest=%s err=%v", digest, err)
	}
	if fixture.RequestCommonKeyCount != NativeXPCRequestCommonKeyCount || fixture.ReplyCommonKeyCount != NativeXPCReplyCommonKeyCount || fixture.ExtraObjectsAllowed != 0 || fixture.FileDescriptorsAllowed != 0 || fixture.EndpointsAllowed != 0 || fixture.MachRightsAllowed != 0 || fixture.NestedContainersAllowed != 0 {
		return fmt.Errorf("native XPC common counts or closed-map controls mismatch")
	}
	if fixture.MessageTagDisposition != "method-specific-cross-check-not-dispatch-opcode" || fixture.RequestIDDisposition != "correlation-only" || fixture.CopyDisposition != "body-copy-only-after-peer-flow-shape-current-state-binding" || fixture.LocalIntegrityDisposition != "oversize-output-short-write-pointer-length-or-bridge-version-fault-terminates-process-without-reply" {
		return fmt.Errorf("native XPC dispatch, copy, request, or integrity disposition mismatch")
	}
	if fixture.CoreRefusalMapping != "classification-selects-status;reason=coreRefusal" || fixture.LocalIntegrityMapping != "terminate-without-reply;reason-tag-is-fixture-diagnostic-only" {
		return fmt.Errorf("native XPC refusal or integrity mapping mismatch")
	}
	if fixture.PeerAuthentication != nil || fixture.ListenerActivated || fixture.ServiceRegistered {
		return fmt.Errorf("native XPC passive scope drift")
	}
	wantMethods := map[string][3]NativeXPCEnvelopeSpec{
		SubmitMainMJSV0Method:     {ExpectedNativeXPCSubmitMainMJSV0Request(), ExpectedNativeXPCSubmitMainMJSV0Reply(), ExpectedNativeXPCSubmitMainMJSV0RefusalReply()},
		RegisterPlanV0Method:      {ExpectedNativeXPCRegisterPlanV0Request(), ExpectedNativeXPCRegisterPlanV0Reply(), ExpectedNativeXPCRegisterPlanV0RefusalReply()},
		GetRegisteredPlanV0Method: {ExpectedNativeXPCGetRegisteredPlanV0Request(), ExpectedNativeXPCGetRegisteredPlanV0Reply(), ExpectedNativeXPCGetRegisteredPlanV0RefusalReply()},
		SubmitApprovalV0Method:    {ExpectedNativeXPCSubmitApprovalV0Request(), ExpectedNativeXPCSubmitApprovalV0Reply(), ExpectedNativeXPCSubmitApprovalV0RefusalReply()},
		RequestAttemptV0Method:    {ExpectedNativeXPCRequestAttemptV0Request(), ExpectedNativeXPCRequestAttemptV0Reply(), ExpectedNativeXPCRequestAttemptV0RefusalReply()},
	}
	if len(fixture.Envelopes) != len(wantMethods) {
		return fmt.Errorf("native XPC envelope method count: got %d want %d", len(fixture.Envelopes), len(wantMethods))
	}
	for method, specs := range wantMethods {
		got, ok := fixture.Envelopes[method]
		if !ok {
			return fmt.Errorf("native XPC missing envelope method %s", method)
		}
		want := [3]nativeXPCFixtureEnvelope{
			nativeXPCFixtureEnvelopeFromSpec(specs[0]),
			nativeXPCFixtureEnvelopeFromSpec(specs[1]),
			nativeXPCFixtureEnvelopeFromSpec(specs[2]),
		}
		if received := [3]nativeXPCFixtureEnvelope{got.Request, got.SuccessReply, got.RefusalReply}; !reflect.DeepEqual(received, want) {
			return fmt.Errorf("native XPC complete dictionary mismatch for %s", method)
		}
	}
	if fixture.CasesDigest != "9ac6845baf35651aab057989264ab7fb17305751d3101df38d26b2334b8ef68e" {
		return fmt.Errorf("native XPC complete ordered-case digest mismatch: got %s", fixture.CasesDigest)
	}
	wantIDs := expectedNativeXPCCaseIDs()
	if len(fixture.Cases) != len(wantIDs) || len(fixture.CaseTable) != len(wantIDs) {
		return fmt.Errorf("native XPC case/refusal count mismatch")
	}
	zero := nativeXPCFixtureZeroState{}
	for index, candidate := range fixture.Cases {
		if candidate.ID != wantIDs[index] {
			return fmt.Errorf("native XPC case %d: got %s want %s", index, candidate.ID, wantIDs[index])
		}
		var expected map[string]json.RawMessage
		if len(candidate.Expected) == 0 {
			return fmt.Errorf("native XPC case %s has no expected result", candidate.ID)
		}
		if candidate.ReplyDisposition == "" || candidate.TerminationDisposition == "" {
			return fmt.Errorf("native XPC case %s lacks reply or termination disposition", candidate.ID)
		}
		if candidate.Expected[0] != '"' {
			if err := json.Unmarshal(candidate.Expected, &expected); err != nil {
				return fmt.Errorf("native XPC case %s expected: %w", candidate.ID, err)
			}
			var decision string
			if rawDecision, ok := expected["decision"]; ok {
				if err := json.Unmarshal(rawDecision, &decision); err != nil {
					return err
				}
			}
			requiresNoState := decision == "reject-before-body-copy" || decision == "no-core-call-no-application-reply" || decision == "accept-outer-shape-only"
			if requiresNoState && (candidate.NoState == nil || !reflect.DeepEqual(*candidate.NoState, zero)) {
				return fmt.Errorf("native XPC case %s missing exact required noState", candidate.ID)
			}
		}
	}
	refusalOrder := []NativeXPCErrorMapping{
		{Malformed, NativeXPCStatusMalformed}, {Unsupported, NativeXPCStatusUnsupported},
		{Schema, NativeXPCStatusSchema}, {Binding, NativeXPCStatusBinding},
		{Authentication, NativeXPCStatusAuthentication}, {Stale, NativeXPCStatusStale},
		{Replay, NativeXPCStatusReplay}, {Capacity, NativeXPCStatusCapacity},
		{TrustState, NativeXPCStatusTrustState}, {LocalFailure, NativeXPCStatusLocalFailure},
		{RecoveryRequired, NativeXPCStatusRecoveryRequired}, {Semantic, NativeXPCStatusSemantic},
		{Domain, NativeXPCStatusDomain},
	}
	wantRefusals := make([]nativeXPCFixtureRefusalReply, 0, len(refusalOrder))
	for _, mapping := range refusalOrder {
		wantRefusals = append(wantRefusals, nativeXPCFixtureRefusalReply{
			Classification: string(mapping.Classification), StatusTag: uint64(mapping.StatusTag),
			ReasonTag: uint64(NativeXPCReasonCoreRefusal), BodyKeysAllowed: 0,
			ExactKeyCount: NativeXPCRefusalReplyKeyCount,
		})
	}
	if !reflect.DeepEqual(fixture.RefusalReplies, wantRefusals) {
		return fmt.Errorf("native XPC complete refusal-reply table mismatch")
	}
	if !reflect.DeepEqual(fixture.ResponseLoss, expectedNativeXPCResponseLoss()) {
		return fmt.Errorf("native XPC exact five-entry response-loss table mismatch")
	}
	if fixture.DeadlineBoundaryRule != NativeXPCDeadlineBoundaryRule || !reflect.DeepEqual(fixture.DeadlineCases, expectedNativeXPCDeadlineCases()) {
		return fmt.Errorf("native XPC exact deadline-boundary table mismatch")
	}
	return nil
}

func nativeXPCFixtureEnvelopeFromSpec(spec NativeXPCEnvelopeSpec) nativeXPCFixtureEnvelope {
	fields := make([]nativeXPCFixtureField, len(spec.Fields))
	for index, field := range spec.Fields {
		fields[index] = nativeXPCFixtureField{
			Key: field.Key, ValueType: field.ValueType, Required: field.Required, MinDataBytes: field.MinDataBytes,
			MaxDataBytes: field.MaxDataBytes, AllowedUInt64: append([]uint64(nil), field.AllowedUInt64...),
			ApplicationData: field.ApplicationData, NonZeroData: field.NonZeroData,
		}
		if field.HasFixedUInt64 {
			value := field.FixedUInt64
			fields[index].FixedUInt64 = &value
		}
		if field.HasFixedString {
			value := field.FixedString
			fields[index].FixedString = &value
		}
	}
	return nativeXPCFixtureEnvelope{
		Method: spec.Method, Direction: spec.Direction, MessageTag: spec.MessageTag,
		ProtocolVersion: spec.ProtocolVersion, MethodVersion: spec.MethodVersion,
		ExactKeyCount: spec.ExactKeyCount, RequiredKeyCount: spec.RequiredKeyCount,
		OptionalKeyCount: spec.OptionalKeyCount, ClosedMap: spec.ClosedMap,
		ApplicationDataMaxBytes: spec.ApplicationDataMaxBytes,
		Fields:                  fields,
	}
}

func expectedNativeXPCResponseLoss() []nativeXPCFixtureResponseLoss {
	return []nativeXPCFixtureResponseLoss{
		{Method: SubmitMainMJSV0Method, Disposition: SubmitMainMJSV0ResponseLossDisposition},
		{Method: RegisterPlanV0Method, Disposition: RegisterPlanV0ResponseLossDisposition},
		{Method: GetRegisteredPlanV0Method, Disposition: GetRegisteredPlanV0ResponseLossDisposition},
		{Method: SubmitApprovalV0Method, Disposition: SubmitApprovalV0ResponseLossDisposition, SemanticIdentity: "canonical-approval-payload+resolved-signer-authorization-identity", SameApprovalID: true, CurrentStateReturned: true, DuplicateAuthorityEffects: 0, RequestIDIsIdempotencyKey: false},
		{Method: RequestAttemptV0Method, Disposition: RequestAttemptV0ResponseLossDisposition, SemanticIdentity: "registration-id+approval-reference", SameAttemptID: true, CurrentStateReturned: true, DuplicateAttempts: 0, DuplicateLifecycleEffects: 0, RequestIDIsIdempotencyKey: false},
	}
}

func expectedNativeXPCDeadlineCases() []nativeXPCFixtureDeadlineCase {
	rule := NativeXPCDeadlineBoundaryRule
	after := NativeXPCDeadlineAfterDispatchDisposition
	postCore := "store-semantic-result-or-recovery-fence-controls"
	zero := nativeXPCFixtureZeroState{}
	var cases []nativeXPCFixtureDeadlineCase
	for _, method := range []string{SubmitApprovalV0Method, RequestAttemptV0Method} {
		cases = append(cases,
			nativeXPCFixtureDeadlineCase{ID: method + ".deadline.immediately-before", Method: method, DeadlineMilliseconds: 5_000, StartsAt: "admission", BoundaryRule: rule, AfterDispatchDisposition: after, Boundary: "immediately-before", ElapsedMilliseconds: 4_999, Expected: nativeXPCFixtureDeadlineExpected{Decision: "dispatch-core-before-deadline", ReplyDisposition: "store-semantic-reply-unless-deadline-cancels-delivery", BodyCopied: true, CoreCalls: 1}, PostCoreState: &postCore},
			nativeXPCFixtureDeadlineCase{ID: method + ".deadline.exactly-at", Method: method, DeadlineMilliseconds: 5_000, StartsAt: "admission", BoundaryRule: rule, AfterDispatchDisposition: after, Boundary: "exactly-at", ElapsedMilliseconds: 5_000, Expected: nativeXPCFixtureDeadlineExpected{Decision: NativeXPCDeadlineExpiredDisposition, ReplyDisposition: "no-application-reply"}, NoState: &zero},
			nativeXPCFixtureDeadlineCase{ID: method + ".deadline.immediately-after", Method: method, DeadlineMilliseconds: 5_000, StartsAt: "admission", BoundaryRule: rule, AfterDispatchDisposition: after, Boundary: "immediately-after", ElapsedMilliseconds: 5_001, Expected: nativeXPCFixtureDeadlineExpected{Decision: NativeXPCDeadlineExpiredDisposition, ReplyDisposition: "no-application-reply"}, NoState: &zero},
		)
	}
	return cases
}

func expectedNativeXPCCaseIDs() []string {
	ids := []string{
		"all.exact-key-and-type-sets", "submit-approval.request.exact-maximum", "request-attempt.request.exact-maximum",
		"submit.body-cap-plus-one", "register.body-cap-plus-one", "get.body-cap-plus-one",
		"submit-approval.body-cap-plus-one", "request-attempt.approval-id-width-plus-one", "all.missing-key",
		"all.extra-key", "all.wrong-type", "all.nested-container", "all.zero-request-id",
		"all.unknown-protocol-version", "all.unknown-method-version",
		"register.joint-protocol-method-record-version-mismatch", "register.joint-method-record-version-mismatch",
		"register.embedded-record-version-mismatch", "submit-approval.joint-protocol-method-record-version-mismatch",
		"submit-approval.joint-method-record-version-mismatch", "submit-approval.embedded-record-version-mismatch",
	}
	methods := []string{SubmitMainMJSV0Method, RegisterPlanV0Method, GetRegisteredPlanV0Method, SubmitApprovalV0Method, RequestAttemptV0Method}
	for _, selected := range methods {
		for _, received := range methods {
			if selected != received {
				ids = append(ids, "cross-service."+selected+".tag-"+received)
			}
		}
	}
	return append(ids,
		"all.wrong-service", "all.wrong-role", "all.wrong-session", "all.wrong-audience", "all.wrong-purpose",
		"all.local-audience-replaced-by-signed-object-audience", "all.local-purpose-replaced-by-signed-object-purpose",
		"submit-approval.signed-audience-replaced-by-local-channel-audience", "submit-approval.signed-purpose-replaced-by-local-channel-purpose",
		"all.request-id-replaced-by-installation-id-bytes", "all.installation-id-replaced-by-request-id-bytes",
		"get.registration-id-replaced-by-request-id-bytes", "submit-approval.registration-id-replaced-by-request-id-bytes",
		"request-attempt.registration-id-replaced-by-approval-id-bytes", "request-attempt.approval-id-replaced-by-registration-id-bytes",
		"all.wrong-installation", "all.wrong-epoch-digest", "all.epoch-uint53-cap-plus-one",
		"all.extra-file-descriptor", "all.extra-endpoint", "all.extra-mach-send-right",
		"all.success-reply-extra-key", "all.success-reply-request-id-mismatch", "all.refusal-reply-extra-body",
		"all.local-integrity-output-fault", "submit-approval.cancel-before-dispatch", "submit-approval.cancel-after-dispatch",
		"request-attempt.cancel-before-dispatch", "request-attempt.cancel-after-dispatch",
	)
}
