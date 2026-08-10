package authorityplane

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type nativeXPCFixtureField struct {
	Key             string             `json:"key"`
	ValueType       NativeXPCValueType `json:"valueType"`
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
	ApplicationDataMaxBytes uint64                  `json:"applicationDataMaxBytes"`
	Fields                  []nativeXPCFixtureField `json:"fields"`
}

type nativeXPCFixtureMethod struct {
	Request      nativeXPCFixtureEnvelope `json:"request"`
	SuccessReply nativeXPCFixtureEnvelope `json:"successReply"`
	RefusalReply nativeXPCFixtureEnvelope `json:"refusalReply"`
}

type nativeXPCFixtureCase struct {
	ID        string          `json:"id"`
	Method    string          `json:"method"`
	Methods   string          `json:"methods"`
	Direction *string         `json:"direction"`
	Mutation  *string         `json:"mutation"`
	Expected  json.RawMessage `json:"expected"`
}

type nativeXPCFixtureCaseTableEntry struct {
	ID             string  `json:"id"`
	Method         string  `json:"method"`
	Direction      *string `json:"direction"`
	Mutation       *string `json:"mutation"`
	Decision       string  `json:"decision"`
	Classification *string `json:"classification"`
	StatusTag      *uint64 `json:"statusTag"`
	ReasonTag      *uint64 `json:"reasonTag"`
	BodyCopied     *bool   `json:"bodyCopied"`
	CoreCalls      *uint64 `json:"coreCalls"`
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
	Method                   string `json:"method"`
	DeadlineMilliseconds     uint64 `json:"deadlineMilliseconds"`
	StartsAt                 string `json:"startsAt"`
	ClientExtensionAllowed   bool   `json:"clientExtensionAllowed"`
	AfterDispatchDisposition string `json:"afterDispatchDisposition"`
}

type nativeXPCFixtureFutureHarnessOracle struct {
	ID              string `json:"id"`
	Scope           string `json:"scope"`
	Expected        string `json:"expected"`
	CurrentEvidence bool   `json:"currentEvidence"`
	InBandRefusal   bool   `json:"inBandRefusal"`
}

func TestNativeXPCContractMatchesGeneratedCrossLanguageFixture(t *testing.T) {
	bytes, err := os.ReadFile(filepath.Join("..", "..", "..", "schemas", "conformance", "authenticated-local-ipc-v0", "native-xpc-v0.contract.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		ObjectType               string                                   `json:"objectType"`
		ObjectVersion            uint64                                   `json:"objectVersion"`
		Status                   string                                   `json:"status"`
		TransportEncoding        string                                   `json:"transportEncoding"`
		TransportEncodingVersion uint64                                   `json:"transportEncodingVersion"`
		MessageTags              map[string]uint64                        `json:"messageTags"`
		ApprovalStateTags        map[string]uint64                        `json:"approvalStateTags"`
		AttemptStateTags         map[string]uint64                        `json:"attemptStateTags"`
		StatusTags               map[string]uint64                        `json:"statusTags"`
		ReasonTags               map[string]uint64                        `json:"reasonTags"`
		ClassificationToStatus   map[string]uint64                        `json:"classificationToStatus"`
		StructuralReasonToStatus map[string]uint64                        `json:"structuralReasonToStatus"`
		MethodBindings           map[string]nativeXPCFixtureMethodBinding `json:"methodBindings"`
		ValidationPrecedence     []string                                 `json:"validationPrecedence"`
		Envelopes                map[string]nativeXPCFixtureMethod        `json:"envelopes"`
		Cases                    []nativeXPCFixtureCase                   `json:"cases"`
		CaseTable                []nativeXPCFixtureCaseTableEntry         `json:"caseTable"`
		RefusalReplies           []nativeXPCFixtureRefusalReply           `json:"refusalReplies"`
		ResponseLoss             []nativeXPCFixtureResponseLoss           `json:"responseLoss"`
		DeadlineCases            []nativeXPCFixtureDeadlineCase           `json:"deadlineCases"`
		FutureHarnessOracles     []nativeXPCFixtureFutureHarnessOracle    `json:"futureNativeHarnessOracles"`
		PeerAuthentication       any                                      `json:"peerAuthenticationEvidence"`
		ListenerActivated        bool                                     `json:"listenerActivated"`
		ServiceRegistered        bool                                     `json:"serviceRegistered"`
	}
	if err := json.Unmarshal(bytes, &fixture); err != nil {
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
	wantDeadlines := []nativeXPCFixtureDeadlineCase{
		{
			Method: SubmitApprovalV0Method, DeadlineMilliseconds: SubmitApprovalV0DeadlineMilliseconds,
			StartsAt: "admission", ClientExtensionAllowed: false,
			AfterDispatchDisposition: "response-unknown-store-semantic-result-or-recovery-fence-controls",
		},
		{
			Method: RequestAttemptV0Method, DeadlineMilliseconds: RequestAttemptV0DeadlineMilliseconds,
			StartsAt: "admission", ClientExtensionAllowed: false,
			AfterDispatchDisposition: "response-unknown-store-semantic-result-or-recovery-fence-controls",
		},
	}
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
	if received.Method != expected.Method || received.Direction != expected.Direction || received.MessageTag != expected.MessageTag || received.ProtocolVersion != expected.ProtocolVersion || received.MethodVersion != expected.MethodVersion || received.ExactKeyCount != expected.ExactKeyCount || received.ApplicationDataMaxBytes != expected.ApplicationDataMaxBytes {
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
		if receivedField.Key != expectedField.Key || receivedField.ValueType != expectedField.ValueType || receivedField.MinDataBytes != expectedField.MinDataBytes || receivedField.MaxDataBytes != expectedField.MaxDataBytes || !reflect.DeepEqual(receivedField.AllowedUInt64, expectedField.AllowedUInt64) || receivedField.ApplicationData != expectedField.ApplicationData || receivedField.NonZeroData != expectedField.NonZeroData {
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
