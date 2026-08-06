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
	ID       string          `json:"id"`
	Expected json.RawMessage `json:"expected"`
}

type nativeXPCFixtureRefusalReply struct {
	Classification  string `json:"classification"`
	StatusTag       uint64 `json:"statusTag"`
	ReasonTag       uint64 `json:"reasonTag"`
	BodyKeysAllowed uint64 `json:"bodyKeysAllowed"`
	ExactKeyCount   uint64 `json:"exactKeyCount"`
}

type nativeXPCFixtureMethodBinding struct {
	Service                   string              `json:"service"`
	ExpectedRole              CallerRole          `json:"expectedRole"`
	ExpectedSigningIdentifier *string             `json:"expectedSigningIdentifier"`
	Audience                  string              `json:"audience"`
	Purpose                   string              `json:"purpose"`
	MessageTag                NativeXPCMessageTag `json:"messageTag"`
	MethodVersion             uint64              `json:"methodVersion"`
}

type nativeXPCFixtureResponseLoss struct {
	Method      string `json:"method"`
	Disposition string `json:"disposition"`
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
		StatusTags               map[string]uint64                        `json:"statusTags"`
		ReasonTags               map[string]uint64                        `json:"reasonTags"`
		ClassificationToStatus   map[string]uint64                        `json:"classificationToStatus"`
		StructuralReasonToStatus map[string]uint64                        `json:"structuralReasonToStatus"`
		MethodBindings           map[string]nativeXPCFixtureMethodBinding `json:"methodBindings"`
		ValidationPrecedence     []string                                 `json:"validationPrecedence"`
		Envelopes                map[string]nativeXPCFixtureMethod        `json:"envelopes"`
		Cases                    []nativeXPCFixtureCase                   `json:"cases"`
		RefusalReplies           []nativeXPCFixtureRefusalReply           `json:"refusalReplies"`
		ResponseLoss             []nativeXPCFixtureResponseLoss           `json:"responseLoss"`
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
	}
	if !reflect.DeepEqual(fixture.MessageTags, wantTags) {
		t.Fatalf("native message tags mismatch: got %#v want %#v", fixture.MessageTags, wantTags)
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
		if received.Service != expected.Service || received.ExpectedRole != expected.ExpectedRole || receivedSigningID != expected.ExpectedSigningIdentifier || received.Audience != expected.Audience || received.Purpose != expected.Purpose || received.MessageTag != expected.MessageTag || received.MethodVersion != expected.MethodVersion {
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
		"all.exact-key-and-type-sets", "submit.body-cap-plus-one",
		"register.body-cap-plus-one", "get.body-cap-plus-one", "all.missing-key",
		"all.extra-key", "all.wrong-type", "all.nested-container", "all.zero-request-id",
		"all.unknown-protocol-version", "all.unknown-method-version",
		"register.joint-protocol-method-record-version-mismatch",
		"register.joint-method-record-version-mismatch", "register.embedded-record-version-mismatch",
		"cross-service.SubmitMainMJSV0.tag-RegisterPlanV0",
		"cross-service.SubmitMainMJSV0.tag-GetRegisteredPlanV0",
		"cross-service.RegisterPlanV0.tag-SubmitMainMJSV0",
		"cross-service.RegisterPlanV0.tag-GetRegisteredPlanV0",
		"cross-service.GetRegisteredPlanV0.tag-SubmitMainMJSV0",
		"cross-service.GetRegisteredPlanV0.tag-RegisterPlanV0",
		"all.wrong-service", "all.wrong-role", "all.wrong-audience", "all.wrong-purpose",
		"all.local-audience-replaced-by-signed-object-audience",
		"all.local-purpose-replaced-by-signed-object-purpose",
		"all.request-id-replaced-by-installation-id-bytes",
		"all.installation-id-replaced-by-request-id-bytes",
		"get.registration-id-replaced-by-request-id-bytes", "all.wrong-installation",
		"all.wrong-epoch-digest", "all.epoch-uint53-cap-plus-one",
		"all.extra-file-descriptor", "all.extra-endpoint", "all.extra-mach-send-right",
		"all.success-reply-extra-key", "all.success-reply-request-id-mismatch",
		"all.refusal-reply-extra-body", "all.local-integrity-output-fault",
	}
	gotCaseIDs := make([]string, len(fixture.Cases))
	for index, candidate := range fixture.Cases {
		gotCaseIDs[index] = candidate.ID
		if len(candidate.Expected) == 0 || candidate.Expected[0] == '"' {
			continue
		}
		var expected struct {
			Classification string  `json:"classification"`
			StatusTag      *uint64 `json:"statusTag"`
			ReasonTag      uint64  `json:"reasonTag"`
		}
		if err := json.Unmarshal(candidate.Expected, &expected); err != nil {
			t.Fatalf("decode native case %s: %v", candidate.ID, err)
		}
		if expected.StatusTag != nil && expected.Classification != "" && *expected.StatusTag != fixture.StatusTags[expected.Classification] {
			t.Fatalf("native case %s status/classification mismatch", candidate.ID)
		}
		if expected.ReasonTag > uint64(NativeXPCReasonLocalIntegrityFault) {
			t.Fatalf("native case %s uses unknown reason %d", candidate.ID, expected.ReasonTag)
		}
	}
	if !reflect.DeepEqual(gotCaseIDs, wantCaseIDs) {
		t.Fatalf("native case order/set mismatch:\n got %#v\nwant %#v", gotCaseIDs, wantCaseIDs)
	}
	wantResponseLoss := []nativeXPCFixtureResponseLoss{
		{Method: SubmitMainMJSV0Method, Disposition: SubmitMainMJSV0ResponseLossDisposition},
		{Method: RegisterPlanV0Method, Disposition: RegisterPlanV0ResponseLossDisposition},
		{Method: GetRegisteredPlanV0Method, Disposition: GetRegisteredPlanV0ResponseLossDisposition},
	}
	if !reflect.DeepEqual(fixture.ResponseLoss, wantResponseLoss) {
		t.Fatalf("native response-loss table mismatch: got %#v want %#v", fixture.ResponseLoss, wantResponseLoss)
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
		if receivedField.Key != expectedField.Key || receivedField.ValueType != expectedField.ValueType || receivedField.MinDataBytes != expectedField.MinDataBytes || receivedField.MaxDataBytes != expectedField.MaxDataBytes || receivedField.ApplicationData != expectedField.ApplicationData || receivedField.NonZeroData != expectedField.NonZeroData {
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
}

func TestNativeXPCContractHasNoGenericDispatchOrAuthorityFields(t *testing.T) {
	for _, spec := range []NativeXPCEnvelopeSpec{
		ExpectedNativeXPCSubmitMainMJSV0Request(),
		ExpectedNativeXPCRegisterPlanV0Request(),
		ExpectedNativeXPCGetRegisteredPlanV0Request(),
	} {
		for _, field := range spec.Fields {
			switch field.Key {
			case "capsule.method", "capsule.command", "capsule.opcode", "capsule.role", "capsule.service", "capsule.backend", "capsule.path", "capsule.image", "capsule.mount":
				t.Fatalf("forbidden generic or authority-bearing native field %q", field.Key)
			}
		}
	}
}
