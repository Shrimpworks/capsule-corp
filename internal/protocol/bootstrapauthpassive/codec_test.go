package bootstrapauthpassive

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const fixtureRoot = "../../../schemas/conformance/i2b-bootstrap-v0"

type fixtureManifest struct {
	ManifestVersion string `json:"manifestVersion"`
	Maxima          struct {
		Request struct{ Payload, Protected, Envelope int } `json:"request"`
		Record  struct{ Payload, Protected, Envelope int } `json:"record"`
	} `json:"maxima"`
	Effects map[string]int `json:"effects"`
	Cases   []fixtureCase  `json:"cases"`
}

type fixtureCase struct {
	ID              string  `json:"id"`
	Object          string  `json:"object"`
	Fixture         string  `json:"fixture"`
	Expected        string  `json:"expected"`
	Decision        string  `json:"decision"`
	TrustedNow      *uint64 `json:"trustedNow"`
	Replay          string  `json:"replay"`
	SelfExpected    bool    `json:"selfExpected"`
	RequestVariant  string  `json:"requestVariant"`
	RequestEnvelope string  `json:"requestEnvelope"`
	RequestPayload  string  `json:"requestPayload"`
}

func TestPassiveBootstrapConformance(t *testing.T) {
	manifest := readManifest(t)
	if manifest.ManifestVersion != "capsule.i2b-bootstrap-conformance/v0" || len(manifest.Cases) != 95 {
		t.Fatalf("unexpected corpus identity/count: %q/%d", manifest.ManifestVersion, len(manifest.Cases))
	}
	if manifest.Maxima.Request.Payload != RequestPayloadCalculatedMaxBytes || manifest.Maxima.Request.Protected != RequestProtectedCalculatedMaxBytes || manifest.Maxima.Request.Envelope != RequestEnvelopeCalculatedMaxBytes ||
		manifest.Maxima.Record.Payload != RecordPayloadCalculatedMaxBytes || manifest.Maxima.Record.Protected != RecordProtectedCalculatedMaxBytes || manifest.Maxima.Record.Envelope != RecordEnvelopeCalculatedMaxBytes {
		t.Fatalf("calculated maxima drift: %+v", manifest.Maxima)
	}
	for effect, count := range manifest.Effects {
		if count != 0 {
			t.Fatalf("passive effect %s = %d, want zero", effect, count)
		}
	}

	codec, err := NewCodec()
	if err != nil {
		t.Fatal(err)
	}
	ordinaryRequestEnvelope := readFixture(t, "request/ordinary.cose")
	ordinaryRequestPayload := readFixture(t, "request/ordinary.payload.cbor")
	maximumRequestEnvelope := readFixture(t, "request/calculated-maximum.cose")
	maximumRequestPayload := readFixture(t, "request/calculated-maximum.payload.cbor")
	ordinaryRecordEnvelope := readFixture(t, "record/ordinary.cose")
	ordinaryRequest := decodeRequestFixture(t, codec, ordinaryRequestEnvelope)
	ordinaryRecord := decodeRecordFixture(t, codec, ordinaryRecordEnvelope)

	for _, tc := range manifest.Cases {
		t.Run(tc.ID, func(t *testing.T) {
			envelope := readFixture(t, tc.Fixture)
			switch tc.Object {
			case "request":
				expected := ordinaryRequest
				if tc.SelfExpected {
					expected = decodeRequestFixture(t, codec, envelope)
				}
				replay := replayFor(tc.Replay, envelope, expected.Nonce)
				now := uint64(2_000_000_001)
				if tc.TrustedNow != nil {
					now = *tc.TrustedNow
				}
				verified, verifyErr := codec.VerifyRequest(envelope, expected, now, replay)
				assertOutcome(t, tc, Decision(tc.Decision), verifiedDecision(verified), verifyErr)
			case "record":
				requestEnvelope, requestPayload := ordinaryRequestEnvelope, ordinaryRequestPayload
				if tc.RequestVariant == "maximum" {
					requestEnvelope, requestPayload = maximumRequestEnvelope, maximumRequestPayload
				}
				if tc.RequestEnvelope != "" {
					requestEnvelope = readFixture(t, tc.RequestEnvelope)
				}
				if tc.RequestPayload != "" {
					requestPayload = readFixture(t, tc.RequestPayload)
				}
				expected := ordinaryRecord
				if tc.SelfExpected {
					expected = decodeRecordFixture(t, codec, envelope)
				}
				replay := replayFor(tc.Replay, envelope, expected.RequestNonce)
				now := uint64(2_000_000_101)
				if tc.TrustedNow != nil {
					now = *tc.TrustedNow
				}
				verified, verifyErr := codec.VerifyRecord(envelope, RecordBindings{
					Expected: expected, RequestEnvelope: requestEnvelope, RequestPayload: requestPayload, TrustedNow: now, Replay: replay,
				})
				assertOutcome(t, tc, Decision(tc.Decision), verifiedRecordDecision(verified), verifyErr)
			default:
				t.Fatalf("unknown object %q", tc.Object)
			}
		})
	}
}

func TestPassiveBootstrapDefensiveCopies(t *testing.T) {
	codec, err := NewCodec()
	if err != nil {
		t.Fatal(err)
	}
	received := readFixture(t, "request/ordinary.cose")
	expected := decodeRequestFixture(t, codec, received)
	original := bytes.Clone(received)
	verified, err := codec.VerifyRequest(received, expected, 2_000_000_001, ReplayState{Disposition: ReplayFresh})
	if err != nil {
		t.Fatal(err)
	}
	received[0] ^= 0xff
	if !bytes.Equal(verified.ExactEnvelope(), original) {
		t.Fatal("caller mutation changed retained envelope")
	}
	accessor := verified.ExactEnvelope()
	accessor[0] ^= 0xff
	if !bytes.Equal(verified.ExactEnvelope(), original) {
		t.Fatal("accessor mutation changed retained envelope")
	}
	payload := verified.ExactPayload()
	payload[0] ^= 0xff
	if bytes.Equal(payload, verified.ExactPayload()) {
		t.Fatal("payload accessor did not return a defensive copy")
	}
	protected := verified.ExactProtected()
	protected[0] ^= 0xff
	if bytes.Equal(protected, verified.ExactProtected()) {
		t.Fatal("protected accessor did not return a defensive copy")
	}
	if verified.EnvelopeDigest() != sha256.Sum256(original) || verified.PayloadDigest() != sha256.Sum256(verified.ExactPayload()) {
		t.Fatal("request exact-byte digest identity changed")
	}

	recordBytes := readFixture(t, "record/ordinary.cose")
	recordExpected := decodeRecordFixture(t, codec, recordBytes)
	recordVerified, err := codec.VerifyRecord(recordBytes, RecordBindings{Expected: recordExpected, RequestEnvelope: original, RequestPayload: readFixture(t, "request/ordinary.payload.cbor"), TrustedNow: 2_000_000_101, Replay: ReplayState{Disposition: ReplayFresh}})
	if err != nil {
		t.Fatal(err)
	}
	recordOriginal := bytes.Clone(recordBytes)
	recordBytes[0] ^= 0xff
	if !bytes.Equal(recordVerified.ExactEnvelope(), recordOriginal) {
		t.Fatal("record caller mutation changed retained envelope")
	}
	recordPayload := recordVerified.ExactPayload()
	recordPayload[0] ^= 0xff
	if bytes.Equal(recordPayload, recordVerified.ExactPayload()) {
		t.Fatal("record payload accessor did not return a defensive copy")
	}
	recordProtected := recordVerified.ExactProtected()
	recordProtected[0] ^= 0xff
	if bytes.Equal(recordProtected, recordVerified.ExactProtected()) {
		t.Fatal("record protected accessor did not return a defensive copy")
	}
	if recordVerified.EnvelopeDigest() != sha256.Sum256(recordOriginal) || recordVerified.PayloadDigest() != sha256.Sum256(recordVerified.ExactPayload()) {
		t.Fatal("record exact-byte digest identity changed")
	}
}

func TestPassiveBootstrapReplayIdentityUsesPayload(t *testing.T) {
	codec, err := NewCodec()
	if err != nil {
		t.Fatal(err)
	}
	ordinaryRequestEnvelope := readFixture(t, "request/ordinary.cose")
	ordinaryRequest := decodeRequestFixture(t, codec, ordinaryRequestEnvelope)
	ordinaryRequestReplay := replayFor("pending-exact", ordinaryRequestEnvelope, ordinaryRequest.Nonce)
	complementaryRequestEnvelope := readFixture(t, "request/complementary-s.cose")

	request, err := codec.VerifyRequest(complementaryRequestEnvelope, ordinaryRequest, 2_000_000_001, ordinaryRequestReplay)
	if err != nil {
		t.Fatal(err)
	}
	if request.Decision() != DecisionResumeExact || request.PayloadDigest() != ordinaryRequestReplay.PayloadDigest || request.EnvelopeDigest() == ordinaryRequestReplay.EnvelopeDigest {
		t.Fatal("complementary-S request did not retain payload-owned replay identity and distinct envelope evidence")
	}

	ordinaryRecordEnvelope := readFixture(t, "record/ordinary.cose")
	ordinaryRecord := decodeRecordFixture(t, codec, ordinaryRecordEnvelope)
	ordinaryRecordReplay := replayFor("completed-exact", ordinaryRecordEnvelope, ordinaryRecord.RequestNonce)
	complementaryRecordEnvelope := readFixture(t, "record/complementary-s.cose")
	record, err := codec.VerifyRecord(complementaryRecordEnvelope, RecordBindings{
		Expected: ordinaryRecord, RequestEnvelope: ordinaryRequestEnvelope,
		RequestPayload: readFixture(t, "request/ordinary.payload.cbor"), TrustedNow: 2_000_000_101,
		Replay: ordinaryRecordReplay,
	})
	if err != nil {
		t.Fatal(err)
	}
	if record.Decision() != DecisionReturnRetained || record.PayloadDigest() != ordinaryRecordReplay.PayloadDigest || record.EnvelopeDigest() == ordinaryRecordReplay.EnvelopeDigest {
		t.Fatal("complementary-S record did not retain payload-owned replay identity and distinct envelope evidence")
	}
}

func TestSigStructureUsesEmptyExternalAAD(t *testing.T) {
	codec, err := NewCodec()
	if err != nil {
		t.Fatal(err)
	}
	envelope := readFixture(t, "request/ordinary.cose")
	parts, err := frameEnvelope(envelope, RequestEnvelopeRawMaxBytes, RequestProtectedRawMaxBytes, RequestPayloadRawMaxBytes)
	if err != nil {
		t.Fatal(err)
	}
	protected := envelope[parts.protectedStart:parts.protectedEnd]
	payload := envelope[parts.payloadStart:parts.payloadEnd]
	structure := sigStructure(protected, payload)
	protectedStart := 2 + len("Signature1")
	protectedHeaderBytes := appendByteString(nil, protected)
	externalAADOffset := protectedStart + len(protectedHeaderBytes)
	if structure[externalAADOffset] != 0x40 {
		t.Fatalf("Sig_structure external_aad encoding = 0x%02x, want empty bstr 0x40", structure[externalAADOffset])
	}
	if err := verifySignature(decodeRequestFixture(t, codec, envelope).InstallationRootPublicKey, protected, payload, envelope[parts.signatureStart:parts.signatureEnd]); err != nil {
		t.Fatal(err)
	}
}

func decodeRequestFixture(t *testing.T, codec *Codec, envelope []byte) Request {
	t.Helper()
	parts, err := frameEnvelope(envelope, RequestEnvelopeRawMaxBytes, RequestProtectedRawMaxBytes, RequestPayloadRawMaxBytes)
	if err != nil {
		t.Fatal(err)
	}
	view, err := codec.decodeRequest(envelope[parts.payloadStart:parts.payloadEnd])
	if err != nil {
		t.Fatal(err)
	}
	return view
}
func decodeRecordFixture(t *testing.T, codec *Codec, envelope []byte) Record {
	t.Helper()
	parts, err := frameEnvelope(envelope, RecordEnvelopeRawMaxBytes, RecordProtectedRawMaxBytes, RecordPayloadRawMaxBytes)
	if err != nil {
		t.Fatal(err)
	}
	view, err := codec.decodeRecord(envelope[parts.payloadStart:parts.payloadEnd])
	if err != nil {
		t.Fatal(err)
	}
	return view
}
func replayFor(kind string, envelope []byte, nonce Nonce) ReplayState {
	if kind == "fresh" {
		return ReplayState{Disposition: ReplayFresh}
	}
	parts, err := frameEnvelope(envelope, RecordEnvelopeRawMaxBytes, RecordProtectedRawMaxBytes, RecordPayloadRawMaxBytes)
	if err != nil {
		panic(err)
	}
	payloadDigest := sha256.Sum256(envelope[parts.payloadStart:parts.payloadEnd])
	envelopeDigest := sha256.Sum256(envelope)
	switch kind {
	case "pending-exact":
		return ReplayState{Disposition: ReplayPending, PayloadDigest: payloadDigest, EnvelopeDigest: envelopeDigest, Nonce: nonce}
	case "completed-exact":
		return ReplayState{Disposition: ReplayCompleted, PayloadDigest: payloadDigest, EnvelopeDigest: envelopeDigest, Nonce: nonce}
	case "pending-other":
		payloadDigest[0] ^= 1
		return ReplayState{Disposition: ReplayPending, PayloadDigest: payloadDigest, EnvelopeDigest: envelopeDigest, Nonce: nonce}
	case "completed-other":
		payloadDigest[0] ^= 1
		return ReplayState{Disposition: ReplayCompleted, PayloadDigest: payloadDigest, EnvelopeDigest: envelopeDigest, Nonce: nonce}
	default:
		panic("unknown replay fixture")
	}
}
func assertOutcome(t *testing.T, tc fixtureCase, want, got Decision, err error) {
	t.Helper()
	if tc.Expected == "ACCEPT" {
		if err != nil || got != want {
			t.Fatalf("got decision/error %q/%v, want %q/nil", got, err, want)
		}
		return
	}
	if err == nil {
		t.Fatalf("refusal case returned %q without error", got)
	}
	if _, ok := err.(*VerificationError); !ok {
		t.Fatalf("refusal did not use bounded classification: %T %v", err, err)
	}
}
func verifiedDecision(v *VerifiedRequest) Decision {
	if v == nil {
		return ""
	}
	return v.Decision()
}
func verifiedRecordDecision(v *VerifiedRecord) Decision {
	if v == nil {
		return ""
	}
	return v.Decision()
}
func readManifest(t *testing.T) fixtureManifest {
	t.Helper()
	var v fixtureManifest
	if err := json.Unmarshal(readFixture(t, "manifest.json"), &v); err != nil {
		t.Fatal(err)
	}
	return v
}
func readFixture(t *testing.T, path string) []byte {
	t.Helper()
	value, err := os.ReadFile(filepath.Join(fixtureRoot, path))
	if err != nil {
		t.Fatal(err)
	}
	return value
}
