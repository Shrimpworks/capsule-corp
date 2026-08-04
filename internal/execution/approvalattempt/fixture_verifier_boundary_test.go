package approvalattempt

import (
	"bytes"
	"context"
	"testing"
)

// TestFixtureVerifierRawAndCalculatedByteGates directly exercises
// FixtureVerifier.Verify and FixtureVerifier.VerifyCandidate at the true
// 512/256/128 raw byte caps and the separate, tighter 431/242/116
// calculated closed-candidate maxima, asserting the exact classification
// and code returned for each -- not only the classification enum. This
// complements the JSON conformance corpus's *-raw-maximum and
// *-calculated-maximum-plus-one cases (schemas/conformance/v0/manifest.json,
// dispatched by TestSliceAConformanceManifest): those prove the same
// boundary through the generic manifest-driven harness, which only asserts
// the classification enum and never calls VerifyCandidate. This test
// additionally proves the exact rejection code and exercises
// VerifyCandidate directly.
//
// Claim boundary: FixtureVerifier is a byte-equality fixture lookup, not a
// COSE/CBOR parser -- payload/protected-header bytes here are trusted
// fixture metadata, not derived by decoding the envelope. These tests prove
// FixtureVerifier's own gate ordering and classification, not general
// COSE/CDDL validity or production-parser behavior.
func TestFixtureVerifierRawAndCalculatedByteGates(t *testing.T) {
	keyID := []byte("approval-test-key")
	authorizationIdentity := repeatedAuthorizationIdentity(0x99)
	view := ordinaryGrantView()
	bindings := ordinaryRoleBindings(keyID, authorizationIdentity)

	newVector := func(envelopeBytes, payloadBytes, protectedBytes int) FixtureVector {
		return FixtureVector{
			EnvelopeBytes:         bytes.Repeat([]byte{0x00}, envelopeBytes),
			PayloadBytes:          bytes.Repeat([]byte{0x00}, payloadBytes),
			ProtectedHeaderBytes:  bytes.Repeat([]byte{0x00}, protectedBytes),
			ProtectedKeyID:        keyID,
			View:                  view,
			ResolvedEpochSequence: 7,
			AuthorizationIdentity: authorizationIdentity,
			SignatureAccepted:     true,
		}
	}

	// A genuinely acceptable vector, well under every raw and calculated
	// bound, registered alongside each rejected vector to prove the
	// rejection leaves the verifier's other fixture untouched.
	accept := newVector(24, 24, 24)

	cases := []struct {
		name    string
		vector  FixtureVector
		wantErr string
	}{
		{"envelope-cap-plus-one", newVector(513, 1, 1), "MALFORMED: envelope-raw-byte-limit"},
		{"payload-cap-plus-one", newVector(8, 257, 1), "MALFORMED: payload-raw-byte-limit"},
		{"protected-cap-plus-one", newVector(8, 1, 129), "MALFORMED: protected-raw-byte-limit"},
		{"envelope-calculated-maximum-plus-one", newVector(432, 1, 1), "SCHEMA: calculated-candidate-maximum"},
		{"envelope-raw-maximum", newVector(512, 1, 1), "SCHEMA: calculated-candidate-maximum"},
		{"payload-calculated-maximum-plus-one", newVector(8, 243, 1), "SCHEMA: calculated-candidate-maximum"},
		{"payload-raw-maximum", newVector(8, 256, 1), "SCHEMA: calculated-candidate-maximum"},
		{"protected-calculated-maximum-plus-one", newVector(8, 1, 117), "SCHEMA: calculated-candidate-maximum"},
		{"protected-raw-maximum", newVector(8, 1, 128), "SCHEMA: calculated-candidate-maximum"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			verifier, err := NewFixtureVerifier([]FixtureVector{testCase.vector, accept})
			if err != nil {
				t.Fatalf("new fixture verifier: %v", err)
			}

			verified, verifyErr := verifier.Verify(context.Background(), testCase.vector.EnvelopeBytes, bindings)
			if verifyErr == nil || verifyErr.Error() != testCase.wantErr {
				t.Fatalf("Verify error = %v, want %q", verifyErr, testCase.wantErr)
			}
			if verified != nil {
				t.Fatal("Verify returned a non-nil approval on rejection")
			}

			candidateVerified, candidateErr := verifier.VerifyCandidate(context.Background(), testCase.vector.EnvelopeBytes)
			if candidateErr == nil || candidateErr.Error() != testCase.wantErr {
				t.Fatalf("VerifyCandidate error = %v, want %q", candidateErr, testCase.wantErr)
			}
			if candidateVerified != nil {
				t.Fatal("VerifyCandidate returned a non-nil approval on rejection")
			}

			// Zero authority/state change: the same verifier, still holding the
			// genuinely acceptable vector registered alongside the rejected
			// one, must accept it exactly as it would have if the rejected
			// call above had never happened.
			acceptVerified, acceptErr := verifier.Verify(context.Background(), accept.EnvelopeBytes, bindings)
			if acceptErr != nil {
				t.Fatalf("Verify of the unrelated accepted vector failed after rejection: %v", acceptErr)
			}
			if acceptVerified == nil {
				t.Fatal("Verify of the unrelated accepted vector returned no approval")
			}
			acceptCandidateVerified, acceptCandidateErr := verifier.VerifyCandidate(context.Background(), accept.EnvelopeBytes)
			if acceptCandidateErr != nil {
				t.Fatalf("VerifyCandidate of the unrelated accepted vector failed after rejection: %v", acceptCandidateErr)
			}
			if acceptCandidateVerified == nil {
				t.Fatal("VerifyCandidate of the unrelated accepted vector returned no approval")
			}
		})
	}
}
