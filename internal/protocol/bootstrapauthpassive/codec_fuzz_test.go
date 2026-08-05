package bootstrapauthpassive

import (
	"bytes"
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"
)

func FuzzVerifyRequest(f *testing.F) {
	ordinary := fuzzFixture(f, "request/ordinary.cose")
	maximum := fuzzFixture(f, "request/calculated-maximum.cose")
	f.Add(ordinary)
	f.Add(maximum)
	f.Add(fuzzFixture(f, "request/request-nonempty-external-aad.cose"))
	f.Add(make([]byte, RequestEnvelopeRawMaxBytes+1))

	codec, err := NewCodec()
	if err != nil {
		f.Fatal(err)
	}
	expected := fuzzRequest(f, codec, ordinary)
	f.Fuzz(func(t *testing.T, received []byte) {
		verified, verifyErr := codec.VerifyRequest(received, expected, 2_000_000_001, ReplayState{Disposition: ReplayFresh})
		if verifyErr != nil {
			return
		}
		if verified.Decision() != DecisionAdmitOnce || verified.View() != expected ||
			!bytes.Equal(verified.ExactEnvelope(), received) ||
			verified.EnvelopeDigest() != sha256.Sum256(received) ||
			verified.PayloadDigest() != sha256.Sum256(verified.ExactPayload()) {
			t.Fatal("accepted request did not preserve its exact authorized byte identity")
		}
	})
}

func FuzzVerifyRecord(f *testing.F) {
	ordinaryRequestEnvelope := fuzzFixture(f, "request/ordinary.cose")
	ordinaryRequestPayload := fuzzFixture(f, "request/ordinary.payload.cbor")
	ordinaryRecord := fuzzFixture(f, "record/ordinary.cose")
	f.Add(ordinaryRecord)
	f.Add(fuzzFixture(f, "record/calculated-maximum.cose"))
	f.Add(fuzzFixture(f, "record/record-nonempty-external-aad.cose"))
	f.Add(make([]byte, RecordEnvelopeRawMaxBytes+1))

	codec, err := NewCodec()
	if err != nil {
		f.Fatal(err)
	}
	expected := fuzzRecord(f, codec, ordinaryRecord)
	f.Fuzz(func(t *testing.T, received []byte) {
		verified, verifyErr := codec.VerifyRecord(received, RecordBindings{
			Expected: expected, RequestEnvelope: ordinaryRequestEnvelope,
			RequestPayload: ordinaryRequestPayload, TrustedNow: 2_000_000_101,
			Replay: ReplayState{Disposition: ReplayFresh},
		})
		if verifyErr != nil {
			return
		}
		if verified.Decision() != DecisionCommitOnce || verified.View() != expected ||
			!bytes.Equal(verified.ExactEnvelope(), received) ||
			verified.EnvelopeDigest() != sha256.Sum256(received) ||
			verified.PayloadDigest() != sha256.Sum256(verified.ExactPayload()) {
			t.Fatal("accepted record did not preserve its exact authorized byte identity")
		}
	})
}

func FuzzVerifyRecordRequestBinding(f *testing.F) {
	requestEnvelope := fuzzFixture(f, "request/ordinary.cose")
	requestPayload := fuzzFixture(f, "request/ordinary.payload.cbor")
	recordEnvelope := fuzzFixture(f, "record/ordinary.cose")
	f.Add(requestEnvelope, requestPayload)
	f.Add(fuzzFixture(f, "request/complementary-s.cose"), requestPayload)
	f.Add(requestEnvelope, fuzzFixture(f, "request/calculated-maximum.payload.cbor"))

	codec, err := NewCodec()
	if err != nil {
		f.Fatal(err)
	}
	expected := fuzzRecord(f, codec, recordEnvelope)
	f.Fuzz(func(t *testing.T, retainedEnvelope, retainedPayload []byte) {
		_, verifyErr := codec.VerifyRecord(recordEnvelope, RecordBindings{
			Expected: expected, RequestEnvelope: retainedEnvelope,
			RequestPayload: retainedPayload, TrustedNow: 2_000_000_101,
			Replay: ReplayState{Disposition: ReplayFresh},
		})
		if verifyErr == nil && (!bytes.Equal(retainedEnvelope, requestEnvelope) || !bytes.Equal(retainedPayload, requestPayload)) {
			t.Fatal("record accepted a substituted retained request pair")
		}
	})
}

func fuzzFixture(f *testing.F, path string) []byte {
	f.Helper()
	value, err := os.ReadFile(filepath.Join(fixtureRoot, path))
	if err != nil {
		f.Fatal(err)
	}
	return value
}

func fuzzRequest(f *testing.F, codec *Codec, envelope []byte) Request {
	f.Helper()
	parts, err := frameEnvelope(envelope, RequestEnvelopeRawMaxBytes, RequestProtectedRawMaxBytes, RequestPayloadRawMaxBytes)
	if err != nil {
		f.Fatal(err)
	}
	view, err := codec.decodeRequest(envelope[parts.payloadStart:parts.payloadEnd])
	if err != nil {
		f.Fatal(err)
	}
	return view
}

func fuzzRecord(f *testing.F, codec *Codec, envelope []byte) Record {
	f.Helper()
	parts, err := frameEnvelope(envelope, RecordEnvelopeRawMaxBytes, RecordProtectedRawMaxBytes, RecordPayloadRawMaxBytes)
	if err != nil {
		f.Fatal(err)
	}
	view, err := codec.decodeRecord(envelope[parts.payloadStart:parts.payloadEnd])
	if err != nil {
		f.Fatal(err)
	}
	return view
}
