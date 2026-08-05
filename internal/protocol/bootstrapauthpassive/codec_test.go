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
	ID             string  `json:"id"`
	Object         string  `json:"object"`
	Fixture        string  `json:"fixture"`
	Expected       string  `json:"expected"`
	Decision       string  `json:"decision"`
	TrustedNow     *uint64 `json:"trustedNow"`
	Replay         string  `json:"replay"`
	SelfExpected   bool    `json:"selfExpected"`
	RequestVariant string  `json:"requestVariant"`
}

func TestPassiveBootstrapConformance(t *testing.T) {
	manifest := readManifest(t)
	if manifest.ManifestVersion != "capsule.i2b-bootstrap-conformance/v0" || len(manifest.Cases) != 71 {
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
	digest := sha256.Sum256(envelope)
	switch kind {
	case "fresh":
		return ReplayState{Disposition: ReplayFresh}
	case "pending-exact":
		return ReplayState{Disposition: ReplayPending, EnvelopeDigest: digest, Nonce: nonce}
	case "completed-exact":
		return ReplayState{Disposition: ReplayCompleted, EnvelopeDigest: digest, Nonce: nonce}
	case "pending-other":
		digest[0] ^= 1
		return ReplayState{Disposition: ReplayPending, EnvelopeDigest: digest, Nonce: nonce}
	case "completed-other":
		digest[0] ^= 1
		return ReplayState{Disposition: ReplayCompleted, EnvelopeDigest: digest, Nonce: nonce}
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
