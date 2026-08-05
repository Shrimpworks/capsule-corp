package approvalattempt

import (
	"bytes"
	"context"
	"crypto/elliptic"
	"encoding/hex"
	"math/big"
	"testing"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

const approvalPublicVectorEnvelopeHex = "d2845854a3012603782b6170706c69636174696f6e2f63617073756c652e617070726f76616c2d6772616e742b63626f723b763d30045820aefada573a36d4abac6bbc06b0723eb9737c2ae349e299b83ae1000bfd32401da058eaac017663617073756c652e617070726f76616c2d6772616e74020003501111111111111111111111111111111104582022222222222222222222222222222222222222222222222222222222222222220550333333333333333333333333333333330658204444444444444444444444444444444444444444444444444444444444444444075055555555555555555555555555555555085066666666666666666666666666666666097463617073756c652e706c616e2e617070726f76650a781c63617073756c652e657865637574696f6e2d73757065727669736f720b1a773594000c1a7735952c58403fd2a1fa446dc9fcd53fea6b04034dbad0d3f8e8e3a5557b3a46d9c5993f2a0a288441ea2bcc98104afe10831f408ef703e121617122ddfc8a21c2c2f419c6c6"
const approvalPublicVectorXHex = "984225585d2285c138033d6140e3cef8b91859704e53c313f8b636ba4f967649"
const approvalPublicVectorYHex = "9734144f46fd19a767a545287c4396b97b69dd38faaea8981adc1a4fed9b401e"

func TestProductionShapedVerifierAcceptsPublicKnownAnswerAndComplementaryS(t *testing.T) {
	verifier, authorization, bindings, envelope := approvalPublicVector(t)
	verified, err := verifier.Verify(context.Background(), envelope, bindings)
	if err != nil {
		t.Fatal(err)
	}
	parts, err := frameApprovalEnvelope(envelope)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(verified.EnvelopeBytes(), envelope) ||
		!bytes.Equal(verified.ProtectedHeaderBytes(), envelope[parts.protectedStart:parts.protectedEnd]) ||
		!bytes.Equal(verified.PayloadBytes(), envelope[parts.payloadStart:parts.payloadEnd]) ||
		!bytes.Equal(verified.ProtectedKeyID(), authorization.KeyID[:]) ||
		verified.AuthorizationIdentity() != authorization.AuthorizationIdentity ||
		verified.ResolvedEpochSequence() != authorization.EpochSequence {
		t.Fatal("verified approval lost exact bytes or trusted key authorization")
	}
	if len(envelope) != 391 || parts.payloadEnd-parts.payloadStart != 234 ||
		parts.protectedEnd-parts.protectedStart != ApprovalProductionProtectedCalculatedMaxBytes ||
		len(envelope) > ApprovalProductionEnvelopeCalculatedMaxBytes ||
		parts.payloadEnd-parts.payloadStart > ApprovalProductionPayloadCalculatedMaxBytes {
		t.Fatalf("public vector or production calculated maxima drift: envelope=%d payload=%d protected=%d", len(envelope), parts.payloadEnd-parts.payloadStart, parts.protectedEnd-parts.protectedStart)
	}

	complementary := bytes.Clone(envelope)
	s := new(big.Int).SetBytes(complementary[parts.signatureStart+32 : parts.signatureEnd])
	s.Sub(elliptic.P256().Params().N, s)
	complementaryS := s.FillBytes(make([]byte, 32))
	copy(complementary[parts.signatureStart+32:parts.signatureEnd], complementaryS)
	if _, err := verifier.Verify(context.Background(), complementary, bindings); err != nil {
		t.Fatalf("mathematically valid complementary S refused: %v", err)
	}

	returned := verified.EnvelopeBytes()
	returned[0] ^= 0xff
	if bytes.Equal(returned, verified.EnvelopeBytes()) {
		t.Fatal("verified bytes alias caller memory")
	}
}

func TestProductionShapedVerifierBindsEveryAuthorityFieldAndTime(t *testing.T) {
	verifier, _, bindings, envelope := approvalPublicVector(t)
	tests := []struct {
		name   string
		mutate func(*ApprovalGrantRoleBindings)
		class  Classification
	}{
		{"installation", func(value *ApprovalGrantRoleBindings) { value.InstallationID[0] ^= 0xff }, ClassificationBinding},
		{"epoch-sequence", func(value *ApprovalGrantRoleBindings) { value.EpochSequence++ }, ClassificationBinding},
		{"epoch-digest", func(value *ApprovalGrantRoleBindings) { value.EpochDigest[0] ^= 0xff }, ClassificationBinding},
		{"registration", func(value *ApprovalGrantRoleBindings) { value.RegistrationID[0] ^= 0xff }, ClassificationBinding},
		{"plan", func(value *ApprovalGrantRoleBindings) { value.PlanDigest[0] ^= 0xff }, ClassificationBinding},
		{"supervisor", func(value *ApprovalGrantRoleBindings) { value.SupervisorID[0] ^= 0xff }, ClassificationBinding},
		{"attempt-nonce", func(value *ApprovalGrantRoleBindings) { value.AttemptNonce[0] ^= 0xff }, ClassificationBinding},
		{"protected-key", func(value *ApprovalGrantRoleBindings) { value.ProtectedKeyID[0] ^= 0xff }, ClassificationBinding},
		{"authorization", func(value *ApprovalGrantRoleBindings) { value.AuthorizationIdentity[0] ^= 0xff }, ClassificationBinding},
		{"before-issue", func(value *ApprovalGrantRoleBindings) { value.EffectiveNow = 1_999_999_999 }, ClassificationStale},
		{"grant-expired", func(value *ApprovalGrantRoleBindings) { value.EffectiveNow = 2_000_000_300 }, ClassificationStale},
		{"registration-expired", func(value *ApprovalGrantRoleBindings) { value.RegistrationExpiresAt = value.EffectiveNow }, ClassificationStale},
		{"grant-past-registration", func(value *ApprovalGrantRoleBindings) { value.RegistrationExpiresAt = 2_000_000_299 }, ClassificationStale},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := cloneRoleBindings(bindings)
			test.mutate(&candidate)
			_, err := verifier.Verify(context.Background(), envelope, candidate)
			assertApprovalClassification(t, err, test.class)
		})
	}
}

func TestProductionShapedVerifierRefusesNonProfileSign1Shapes(t *testing.T) {
	verifier, _, bindings, envelope := approvalPublicVector(t)
	parts, err := frameApprovalEnvelope(envelope)
	if err != nil {
		t.Fatal(err)
	}
	protected := bytes.Clone(envelope[parts.protectedStart:parts.protectedEnd])
	payload := bytes.Clone(envelope[parts.payloadStart:parts.payloadEnd])
	signature := bytes.Clone(envelope[parts.signatureStart:parts.signatureEnd])

	wrongOrder := append([]byte{0xa3}, protected[49:]...)
	wrongOrder = append(wrongOrder, protected[1:49]...)
	wrongAlgorithm := bytes.Clone(protected)
	wrongAlgorithm[2] = 0x25
	wrongContentType := bytes.Clone(protected)
	wrongContentType[6] ^= 1
	unknownProtected := bytes.Clone(protected)
	unknownProtected[0] = 0xa4
	unknownProtected = append(unknownProtected, 0x05, 0x00)
	tests := []struct {
		name     string
		envelope []byte
		class    Classification
	}{
		{"untagged", envelope[1:], ClassificationMalformed},
		{"wrong-array-length", append([]byte(nil), append([]byte{0xd2, 0x83}, envelope[2:]...)...), ClassificationMalformed},
		{"nonempty-unprotected", approvalEnvelope(protected, []byte{0xa1, 0x01, 0x00}, payload, signature), ClassificationMalformed},
		{"detached-payload", approvalEnvelope(protected, []byte{0xa0}, nil, signature), ClassificationMalformed},
		{"short-raw-signature", approvalEnvelope(protected, []byte{0xa0}, payload, signature[:63]), ClassificationMalformed},
		{"noncanonical-protected-order", approvalEnvelope(wrongOrder, []byte{0xa0}, payload, signature), ClassificationMalformed},
		{"wrong-algorithm", approvalEnvelope(wrongAlgorithm, []byte{0xa0}, payload, signature), ClassificationUnsupported},
		{"wrong-content-type", approvalEnvelope(wrongContentType, []byte{0xa0}, payload, signature), ClassificationUnsupported},
		{"unknown-protected-field", approvalEnvelope(unknownProtected, []byte{0xa0}, payload, signature), ClassificationSchema},
		{"tampered-payload", func() []byte { value := bytes.Clone(envelope); value[parts.payloadEnd-1] ^= 1; return value }(), ClassificationBinding},
		{"calculated-maximum-plus-one", approvalEnvelope(protected, []byte{0xa0}, bytes.Repeat([]byte{0}, ApprovalProductionPayloadCalculatedMaxBytes+1), signature), ClassificationSchema},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := verifier.Verify(context.Background(), test.envelope, bindings)
			assertApprovalClassification(t, err, test.class)
		})
	}
}

func TestProductionShapedVerifierFreezesApprovalKeyPolicy(t *testing.T) {
	_, authorization, _, _ := approvalPublicVector(t)
	if authorization.TeamID != ApprovalKeyTeamID || authorization.RoleID != ApprovalBrokerRoleID ||
		authorization.AccessGroup != ApprovalKeyAccessGroupPrefix+"7" ||
		authorization.Protection != ApprovalKeyProtection || authorization.AccessControl != ApprovalKeyAccessControl ||
		authorization.ContextPolicy != ApprovalKeyContextPolicy || authorization.SoftwareFallback {
		t.Fatal("approval key policy projection drift")
	}
	tests := []struct {
		name   string
		mutate func(*ApprovalKeyAuthorization)
	}{
		{"team", func(value *ApprovalKeyAuthorization) { value.TeamID = "OTHER" }},
		{"role", func(value *ApprovalKeyAuthorization) { value.RoleID = "other" }},
		{"access-group", func(value *ApprovalKeyAuthorization) { value.AccessGroup += ".wrong" }},
		{"software-fallback", func(value *ApprovalKeyAuthorization) { value.SoftwareFallback = true }},
		{"protection", func(value *ApprovalKeyAuthorization) { value.Protection = "software-p256" }},
		{"access-control", func(value *ApprovalKeyAuthorization) { value.AccessControl = "userPresence" }},
		{"context-policy", func(value *ApprovalKeyAuthorization) { value.ContextPolicy = "reused" }},
		{"status", func(value *ApprovalKeyAuthorization) { value.Status = "revoked" }},
		{"unsafe-validity", func(value *ApprovalKeyAuthorization) {
			value.ExpiresAt = v0candidate.UInt53(v0candidate.MaxSafeInteger + 1)
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			candidate := authorization
			test.mutate(&candidate)
			if _, err := NewProductionShapedVerifier([]ApprovalKeyAuthorization{candidate}); err == nil {
				t.Fatal("invalid approval key policy accepted")
			}
		})
	}
	if _, err := NewProductionShapedVerifier(nil); err == nil {
		t.Fatal("empty approval key authorization set accepted")
	}
	over := make([]ApprovalKeyAuthorization, ApprovalKeyAuthorizationMaxEntries+1)
	for index := range over {
		over[index] = authorization
	}
	if _, err := NewProductionShapedVerifier(over); err == nil {
		t.Fatal("approval key authorization cap plus one accepted")
	}
}

func approvalPublicVector(t *testing.T) (*ProductionShapedVerifier, ApprovalKeyAuthorization, ApprovalGrantRoleBindings, []byte) {
	t.Helper()
	decode := func(value string) []byte {
		raw, err := hex.DecodeString(value)
		if err != nil {
			t.Fatal(err)
		}
		return raw
	}
	var publicKey ApprovalP256PublicKey
	copy(publicKey.X[:], decode(approvalPublicVectorXHex))
	copy(publicKey.Y[:], decode(approvalPublicVectorYHex))
	installation := repeatedApprovalID[v0candidate.InstallationID](0x11)
	epochDigest := repeatedApprovalDigest[v0candidate.TrustEpochDigest](0x22)
	authorization, err := NewPassiveApprovalKeyAuthorization(publicKey, installation, 7, epochDigest, 1_999_999_900, 2_000_000_400)
	if err != nil {
		t.Fatal(err)
	}
	verifier, err := NewProductionShapedVerifier([]ApprovalKeyAuthorization{authorization})
	if err != nil {
		t.Fatal(err)
	}
	nonce, err := NewAttemptNonce(DomainIdentifier{domain: DomainAttemptNonce, bytes: repeatedApprovalID[[16]byte](0x66)})
	if err != nil {
		t.Fatal(err)
	}
	bindings := ApprovalGrantRoleBindings{
		InstallationID: installation, EpochSequence: 7, EpochDigest: epochDigest,
		RegistrationID: repeatedApprovalID[v0candidate.RegistrationID](0x33),
		PlanDigest:     repeatedApprovalDigest[v0candidate.ExecutionPlanDigest](0x44),
		SupervisorID:   repeatedApprovalID[v0candidate.SupervisorID](0x55), AttemptNonce: nonce,
		ProtectedKeyID: authorization.KeyID[:], AuthorizationIdentity: authorization.AuthorizationIdentity,
		EffectiveNow: 2_000_000_100, RegistrationExpiresAt: 2_000_000_350,
	}
	return verifier, authorization, bindings, decode(approvalPublicVectorEnvelopeHex)
}

func approvalEnvelope(protected, unprotected, payload, signature []byte) []byte {
	result := []byte{0xd2, 0x84}
	result = appendApprovalByteString(result, protected)
	result = append(result, unprotected...)
	result = appendApprovalByteString(result, payload)
	return appendApprovalByteString(result, signature)
}

func appendApprovalByteString(target, value []byte) []byte {
	switch {
	case len(value) < 24:
		target = append(target, 0x40|byte(len(value)))
	case len(value) < 256:
		target = append(target, 0x58, byte(len(value)))
	default:
		target = append(target, 0x59, byte(len(value)>>8), byte(len(value)))
	}
	return append(target, value...)
}

func assertApprovalClassification(t *testing.T, err error, expected Classification) {
	t.Helper()
	classification, ok := ErrorClassification(err)
	if !ok || classification != expected {
		t.Fatalf("classification = %q, error = %v, want %q", classification, err, expected)
	}
}

func repeatedApprovalID[T ~[16]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}

func repeatedApprovalDigest[T ~[32]byte](value byte) T {
	var result T
	for index := range result {
		result[index] = value
	}
	return result
}
