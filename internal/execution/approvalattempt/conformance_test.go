package approvalattempt

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

const conformanceRoot = "../../../schemas/conformance/v0"

type conformanceManifest struct {
	Cases []struct {
		ID      string           `json:"id"`
		Object  string           `json:"object"`
		Fixture fixtureReference `json:"fixture"`
		Context struct {
			Kind            string            `json:"kind"`
			Operation       fixtureReference  `json:"operation"`
			Before          fixtureReference  `json:"before"`
			Payload         *fixtureReference `json:"payload"`
			ProtectedHeader *fixtureReference `json:"protectedHeader"`
		} `json:"context"`
		Expected struct {
			Decision                   string  `json:"decision"`
			Classification             *string `json:"classification"`
			AuthorityStateChanged      bool    `json:"authorityStateChanged"`
			TimeHighWaterChanged       bool    `json:"timeHighWaterChanged"`
			TrustStateTightened        bool    `json:"trustStateTightened"`
			FakeBackendEffectPermitted bool    `json:"fakeBackendEffectPermitted"`
			StateDelta                 struct {
				After fixtureReference `json:"after"`
			} `json:"stateDelta"`
		} `json:"expected"`
		Implementations map[string]string `json:"implementations"`
	} `json:"cases"`
}

type fixtureReference struct {
	Path       string `json:"path"`
	SHA256     string `json:"sha256"`
	ByteLength int    `json:"byteLength"`
}

type operationFixture struct {
	ContextType     string `json:"contextType"`
	ContextVersion  uint64 `json:"contextVersion"`
	Mode            string `json:"mode"`
	ExpectedDomain  string `json:"expectedDomain"`
	ProvidedDomain  string `json:"providedDomain"`
	ReferenceKind   string `json:"referenceKind"`
	Vector          string `json:"vector"`
	BindingMutation string `json:"bindingMutation"`
	CallerMutation  bool   `json:"callerMutation"`
}

func TestSliceAConformanceManifest(t *testing.T) {
	manifest := readJSON[conformanceManifest](t, fixturePath("manifest.json"))
	caseCount := 0
	acceptCount := 0
	rejectCount := 0
	for _, testCase := range manifest.Cases {
		if testCase.Context.Kind != "approval-attempt-state" {
			continue
		}
		operation := readJSON[operationFixture](t, fixturePath(testCase.Context.Operation.Path))
		if operation.Mode == "store-transition" {
			continue
		}
		caseCount++
		if testCase.Expected.Decision == "accept" {
			acceptCount++
		} else {
			rejectCount++
		}
		if testCase.Implementations["go"] != "verified" {
			t.Fatalf("Slice A case %s is not marked verified for Go", testCase.ID)
		}
		testCase := testCase
		t.Run(testCase.ID, func(t *testing.T) {
			before := readBytes(t, testCase.Context.Before)
			after := readBytes(t, testCase.Expected.StateDelta.After)
			if !bytes.Equal(before, after) {
				t.Fatal("passive case changed retained approval/attempt state")
			}
			if testCase.Expected.AuthorityStateChanged || testCase.Expected.TimeHighWaterChanged ||
				testCase.Expected.TrustStateTightened || testCase.Expected.FakeBackendEffectPermitted {
				t.Fatal("passive case permits an authority, time, trust, or fake-backend effect")
			}

			fixture := readBytes(t, testCase.Fixture)
			var err error
			switch operation.Mode {
			case "identifier":
				err = exerciseIdentifier(operation, fixture)
			case "reference":
				err = exerciseReference(operation, fixture)
			case "classification-vocabulary":
				err = exerciseVocabulary(fixture)
			case "fixture-verifier":
				err = exerciseVerifier(t, operation, fixture, testCase.Context.Payload, testCase.Context.ProtectedHeader)
			default:
				t.Fatalf("unimplemented Slice A operation mode %q", operation.Mode)
			}
			assertDecision(t, testCase.Expected.Decision, testCase.Expected.Classification, err)
		})
	}
	if caseCount != 44 || acceptCount != 10 || rejectCount != 34 {
		t.Fatalf("Slice A matrix = %d cases (%d accept, %d reject), want 44 (10, 34)", caseCount, acceptCount, rejectCount)
	}
}

func exerciseIdentifier(operation operationFixture, fixture []byte) error {
	value, err := NewDomainIdentifier(IdentifierDomain(operation.ProvidedDomain), fixture)
	if err != nil {
		return err
	}
	switch IdentifierDomain(operation.ExpectedDomain) {
	case DomainApprovalID:
		_, err = NewApprovalID(value)
	case DomainAttemptID:
		_, err = NewAttemptID(value)
	case DomainAttemptNonce:
		_, err = NewAttemptNonce(value)
	default:
		return classified(ClassificationDomain, "unknown-expected-domain")
	}
	return err
}

func exerciseReference(operation operationFixture, fixture []byte) error {
	value, err := NewDomainIdentifier(IdentifierDomain(operation.ProvidedDomain), fixture)
	if err != nil {
		return err
	}
	switch operation.ReferenceKind {
	case "approval-reference":
		identifier, conversionErr := NewApprovalID(value)
		if conversionErr != nil {
			return conversionErr
		}
		_, err = NewApprovalReference(identifier)
	case "attempt-reference":
		identifier, conversionErr := NewAttemptID(value)
		if conversionErr != nil {
			return conversionErr
		}
		_, err = NewAttemptReference(identifier)
	default:
		return classified(ClassificationDomain, "unknown-reference-kind")
	}
	return err
}

func exerciseVocabulary(fixture []byte) error {
	var got []Classification
	if err := json.Unmarshal(fixture, &got); err != nil {
		return classified(ClassificationSchema, "classification-json")
	}
	want := Classifications()
	if !reflect.DeepEqual(got, want) {
		return classified(ClassificationSchema, "classification-vocabulary")
	}
	return nil
}

func exerciseVerifier(
	t *testing.T,
	operation operationFixture,
	envelope []byte,
	payloadReference *fixtureReference,
	protectedReference *fixtureReference,
) error {
	t.Helper()
	payload := readOptionalBytes(t, payloadReference)
	protectedHeader := readOptionalBytes(t, protectedReference)
	view := ordinaryGrantView()
	keyID := []byte("approval-test-key")
	if operation.Vector == "calculated-maximum" {
		view.IssuedAt = v0candidate.UInt53(v0candidate.MaxSafeInteger - 1)
		view.ExpiresAt = v0candidate.UInt53(v0candidate.MaxSafeInteger)
		keyID = bytes.Repeat([]byte{0x6b}, 64)
	}
	switch operation.Vector {
	case "object-type":
		view.ObjectType = "capsule.execution-attempt"
	case "object-version":
		view.ObjectVersion = 1
	case "purpose":
		view.Purpose = "capsule.execution.attest"
	case "audience":
		view.Audience = "capsule.agent-daemon"
	}
	authorizationIdentity := repeatedAuthorizationIdentity(0x99)
	vectorEnvelope := envelope
	if operation.Vector == "envelope-cap-plus-one" {
		vectorEnvelope = []byte{0x01}
	}
	vector := FixtureVector{
		EnvelopeBytes:         bytes.Clone(vectorEnvelope),
		PayloadBytes:          bytes.Clone(payload),
		ProtectedHeaderBytes:  bytes.Clone(protectedHeader),
		ProtectedKeyID:        bytes.Clone(keyID),
		View:                  view,
		ResolvedEpochSequence: 7,
		AuthorizationIdentity: authorizationIdentity,
		SignatureAccepted:     true,
	}
	verifier, err := NewFixtureVerifier([]FixtureVector{vector})
	if err != nil {
		t.Fatalf("new fixture verifier: %v", err)
	}
	if operation.CallerMutation {
		vector.EnvelopeBytes[0] ^= 0xff
		if len(vector.PayloadBytes) > 0 {
			vector.PayloadBytes[0] ^= 0xff
		}
		if len(vector.ProtectedHeaderBytes) > 0 {
			vector.ProtectedHeaderBytes[0] ^= 0xff
		}
		vector.ProtectedKeyID[0] ^= 0xff
	}
	bindings := ordinaryRoleBindings(keyID, authorizationIdentity)
	mutateBindings(&bindings, operation.BindingMutation)
	originalEnvelope := bytes.Clone(envelope)
	verified, err := verifier.Verify(context.Background(), envelope, bindings)
	if err != nil {
		return err
	}
	if operation.CallerMutation {
		envelope[0] ^= 0xff
		returnedEnvelope := verified.EnvelopeBytes()
		returnedPayload := verified.PayloadBytes()
		returnedProtected := verified.ProtectedHeaderBytes()
		returnedKeyID := verified.ProtectedKeyID()
		returnedEnvelope[0] ^= 0xff
		if len(returnedPayload) > 0 {
			returnedPayload[0] ^= 0xff
		}
		if len(returnedProtected) > 0 {
			returnedProtected[0] ^= 0xff
		}
		returnedKeyID[0] ^= 0xff
		if !bytes.Equal(verified.EnvelopeBytes(), originalEnvelope) ||
			!bytes.Equal(verified.PayloadBytes(), payload) ||
			!bytes.Equal(verified.ProtectedHeaderBytes(), protectedHeader) ||
			!bytes.Equal(verified.ProtectedKeyID(), keyID) {
			t.Fatal("verified approval exposed caller-owned byte storage")
		}
	}
	if operation.Vector == "ordinary" {
		if digest := sha256.Sum256(verified.EnvelopeBytes()); hex.EncodeToString(digest[:]) != "fb0a9e7c983f6f3986260dce857edf6b18cba99ee386f9532300dbdc31a5a3bd" {
			t.Fatalf("ordinary known-answer digest changed: %x", digest)
		}
	}
	if operation.Vector == "calculated-maximum" &&
		(len(verified.EnvelopeBytes()) != ApprovalEnvelopeCalculatedMaxBytes ||
			len(verified.PayloadBytes()) != ApprovalPayloadCalculatedMaxBytes ||
			len(verified.ProtectedHeaderBytes()) != ApprovalProtectedCalculatedMaxBytes) {
		t.Fatal("calculated candidate maxima changed")
	}
	return nil
}

func ordinaryGrantView() ApprovalGrant {
	nonceValue, _ := NewDomainIdentifier(DomainAttemptNonce, bytes.Repeat([]byte{0x66}, 16))
	nonce, _ := NewAttemptNonce(nonceValue)
	return ApprovalGrant{
		ObjectType:     ApprovalGrantObjectType,
		ObjectVersion:  0,
		InstallationID: repeated16[v0candidate.InstallationID](0x11),
		EpochDigest:    repeated32[v0candidate.TrustEpochDigest](0x22),
		RegistrationID: repeated16[v0candidate.RegistrationID](0x33),
		PlanDigest:     repeated32[v0candidate.ExecutionPlanDigest](0x44),
		SupervisorID:   repeated16[v0candidate.SupervisorID](0x55),
		AttemptNonce:   nonce,
		Purpose:        ApprovalGrantPurpose,
		Audience:       ApprovalGrantAudience,
		IssuedAt:       1_785_456_000,
		ExpiresAt:      1_785_456_300,
	}
}

func ordinaryRoleBindings(
	keyID []byte,
	authorizationIdentity ApprovalKeyAuthorizationIdentity,
) ApprovalGrantRoleBindings {
	view := ordinaryGrantView()
	return ApprovalGrantRoleBindings{
		InstallationID:        view.InstallationID,
		EpochSequence:         7,
		EpochDigest:           view.EpochDigest,
		RegistrationID:        view.RegistrationID,
		PlanDigest:            view.PlanDigest,
		SupervisorID:          view.SupervisorID,
		AttemptNonce:          view.AttemptNonce,
		ProtectedKeyID:        bytes.Clone(keyID),
		AuthorizationIdentity: authorizationIdentity,
	}
}

func mutateBindings(bindings *ApprovalGrantRoleBindings, mutation string) {
	switch mutation {
	case "", "none":
	case "installation":
		bindings.InstallationID[0] ^= 0xff
	case "epoch-sequence":
		bindings.EpochSequence++
	case "epoch-digest":
		bindings.EpochDigest[0] ^= 0xff
	case "registration":
		bindings.RegistrationID[0] ^= 0xff
	case "plan-digest":
		bindings.PlanDigest[0] ^= 0xff
	case "supervisor":
		bindings.SupervisorID[0] ^= 0xff
	case "attempt-nonce":
		bindings.AttemptNonce[0] ^= 0xff
	case "protected-key-id":
		bindings.ProtectedKeyID[0] ^= 0xff
	case "authorization-identity":
		bindings.AuthorizationIdentity[0] ^= 0xff
	default:
		panic("unknown binding mutation: " + mutation)
	}
}

func repeatedAuthorizationIdentity(value byte) ApprovalKeyAuthorizationIdentity {
	identity, err := NewApprovalKeyAuthorizationIdentity(bytes.Repeat([]byte{value}, 32))
	if err != nil {
		panic(err)
	}
	return identity
}

func repeated16[T ~[16]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}

func repeated32[T ~[32]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}

func assertDecision(t *testing.T, decision string, expected *string, err error) {
	t.Helper()
	if decision == "accept" {
		if err != nil {
			t.Fatalf("unexpected rejection: %v", err)
		}
		return
	}
	if err == nil {
		t.Fatal("expected rejection")
	}
	classification, ok := ErrorClassification(err)
	if !ok || expected == nil || string(classification) != *expected {
		t.Fatalf("classification = %q (known %t), want %v: %v", classification, ok, expected, err)
	}
}

func readOptionalBytes(t *testing.T, reference *fixtureReference) []byte {
	t.Helper()
	if reference == nil {
		return nil
	}
	return readBytes(t, *reference)
}

func readBytes(t *testing.T, reference fixtureReference) []byte {
	t.Helper()
	value, err := os.ReadFile(fixturePath(reference.Path))
	if err != nil {
		t.Fatalf("read %s: %v", reference.Path, err)
	}
	if len(value) != reference.ByteLength {
		t.Fatalf("%s length = %d, want %d", reference.Path, len(value), reference.ByteLength)
	}
	digest := sha256.Sum256(value)
	if hex.EncodeToString(digest[:]) != reference.SHA256 {
		t.Fatalf("%s digest mismatch", reference.Path)
	}
	return value
}

func readJSON[T any](t *testing.T, path string) T {
	t.Helper()
	value, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var result T
	if err := json.Unmarshal(value, &result); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return result
}

func fixturePath(path string) string { return filepath.Join(conformanceRoot, path) }
