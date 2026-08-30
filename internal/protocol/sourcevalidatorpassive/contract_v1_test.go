package sourcevalidatorpassive

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

func TestV1KnownAnswersRemainRoleDistinct(t *testing.T) {
	t.Parallel()
	answers := []struct {
		path   string
		length int
		digest string
	}{
		{"mjs-source-validator-v1/daemon/resource-policy-inactive.bin", 256, "c198dac71f3b5c2d2e8cca34fc3e9c01ff7b8093ef1a881d8160a34800ff1098"},
		{"mjs-source-validator-v1/approval-broker/resource-policy-inactive.bin", 256, "b0ce8504190b5fe9b0a0296c22340a6439ab453cb32f32c19ddb6e594698568d"},
		{"mjs-source-validator-v1/daemon/process-profile.bin", 256, "272ed9a2bda282a455c042070824bad317c5756cf90c9b62cd417eb27e3ab1e8"},
		{"mjs-source-validator-v1/approval-broker/process-profile.bin", 256, "99e48f3de3ab673f131019db5aa39d2689cc4b1264e11a204643f3af95319383"},
		{"mjs-source-validator-v1/daemon/artifact-profile.bin", 256, "6a35faadeaafb2161af879922f4055dd515335654e8bb27d98060efe5b974208"},
		{"mjs-source-validator-v1/approval-broker/artifact-profile.bin", 256, "fec134344e26cfd3cdfdb85e1398d2bc915c47e106b24a94ffaf4f159549a256"},
		{"mjs-source-validator-v1/daemon/consumer-projection.bin", 192, "0785bf0c4dbedf0a0c4f07eb27c529b74638f05e95ac9021d8362135aad80cb2"},
		{"mjs-source-validator-v1/approval-broker/consumer-projection.bin", 192, "ecc2baa63c5b74a5f290d352f01cd8fc11c9a04cd7a10b6bff7664bf59d9e599"},
		{"mjs-source-validator-v1/daemon/request-ordinary.bin", 273, "7e388687da643527d37118bf7d2bc6c77a1bbdabab9b7c23e0c05e90993ba932"},
		{"mjs-source-validator-v1/approval-broker/request-ordinary.bin", 273, "a19eea21d15a7e8b414de140f563f7acdaf64427be052419313a79ace9351824"},
		{"mjs-source-validator-v1/daemon/request-exact-maximum.bin", V1RequestMaximumBytes, "ce4ab493a95bfe4dfd7f0c61bea27fc14a8451dc4f2043a6aacac0c97ecd5c20"},
		{"mjs-source-validator-v1/approval-broker/request-exact-maximum.bin", V1RequestMaximumBytes, "ee9b8d75228d9faa5df5ddb525c93f2bdcdae07c7322b39ec3c1585e32a26dd5"},
		{"mjs-source-validator-v1/daemon/result-ordinary.bin", 248, "afbc80dd7256bd2b86ae40fcef4ff5c44527f80625bad974b32eb27be4d181cb"},
		{"mjs-source-validator-v1/approval-broker/result-ordinary.bin", 248, "4ee78a722f8da48d1529aad31a12aa9481e4a291665d4276a9e1bdfa1f9c98ee"},
	}
	for _, answer := range answers {
		fixture := readFixture(t, answer.path)
		if len(fixture) != answer.length {
			t.Fatalf("%s length = %d, want %d", answer.path, len(fixture), answer.length)
		}
		digest := sha256.Sum256(fixture)
		if got := hex.EncodeToString(digest[:]); got != answer.digest {
			t.Fatalf("%s digest = %s, want %s", answer.path, got, answer.digest)
		}
	}
}

func TestV1PassiveConformanceCorpus(t *testing.T) {
	t.Parallel()
	manifest := readManifest(t)
	cases := make([]conformanceCase, 0, 46)
	for _, candidate := range manifest.Cases {
		if candidate.Expected.Owner == "source-validator-passive-v1-contract" && candidate.Implementations.Go == "verified" {
			cases = append(cases, candidate)
		}
	}
	if len(cases) != 46 {
		t.Fatalf("verified Go passive v1 cases = %d, want 46", len(cases))
	}

	for _, candidate := range cases {
		candidate := candidate
		t.Run(candidate.ID, func(t *testing.T) {
			t.Parallel()
			role := V1DaemonRole
			if strings.Contains(candidate.ID, ".approval-broker.") {
				role = V1ApprovalBrokerRole
			}
			fixture := readFixture(t, candidate.Fixture.Path)
			var err error
			switch candidate.Object {
			case "SourceValidatorV1ResourcePolicy":
				_, err = DecodeV1ResourcePolicy(role, fixture)
			case "SourceValidatorV1ProcessProfile":
				_, err = DecodeV1ProcessProfile(role, fixture)
			case "SourceValidatorV1ArtifactProfile":
				_, err = DecodeV1ArtifactProfile(role, fixture)
			case "SourceValidatorV1ConsumerProjection":
				_, err = DecodeV1ConsumerProjection(role, fixture)
			case "SourceValidatorV1Request":
				_, err = VerifyV1Request(role, fixtureV1Bindings(t, role), fixture)
			case "SourceValidatorV1Result":
				if candidate.Context.Request == nil {
					t.Fatal("result case lacks contextual request")
				}
				request, requestErr := DecodeV1Request(role, readFixture(t, candidate.Context.Request.Path))
				if requestErr != nil {
					t.Fatalf("decode contextual request: %v", requestErr)
				}
				_, err = VerifyV1Result(request, fixture)
			default:
				t.Fatalf("unexpected passive v1 object %q", candidate.Object)
			}

			if candidate.Expected.Decision == "accept" {
				if err != nil {
					t.Fatalf("accepted case refused: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("rejected case accepted")
			}
			if got := errClassification(err); got != candidate.Expected.Classification {
				t.Fatalf("classification = %q; want %q (%v)", got, candidate.Expected.Classification, err)
			}
			if candidate.Expected.Effects == nil || candidate.Expected.Effects.State || candidate.Expected.Effects.Approval ||
				candidate.Expected.Effects.Key || candidate.Expected.Effects.IPCEndpoint || candidate.Expected.Effects.Process ||
				candidate.Expected.Effects.Runtime || candidate.Expected.Effects.Backend || candidate.Expected.Effects.Guest {
				t.Fatal("rejection does not assert the complete zero-effect oracle")
			}
		})
	}
}

func fixtureV1Bindings(t *testing.T, role V1ConsumerRole) V1RequestBindings {
	t.Helper()
	name := role.String()
	artifactDigest := sha256.Sum256(readFixture(t, "mjs-source-validator-v1/"+name+"/artifact-profile.bin"))
	policyDigest := sha256.Sum256(readFixture(t, "mjs-source-validator-v1/"+name+"/resource-policy-inactive.bin"))
	bindings := V1RequestBindings{
		EpochSequence:         7,
		ArtifactProfileDigest: V1Digest(artifactDigest),
		ResourcePolicyDigest:  V1Digest(policyDigest),
	}
	bindings.CorrelationID[0] = byte(0x50 + role)
	bindings.InstallationID[0] = 0x11
	for index := range bindings.EpochDigest {
		bindings.EpochDigest[index] = 0x22
	}
	return bindings
}

func TestV1RoleSeparatedRequestRoundTripAndCopies(t *testing.T) {
	t.Parallel()
	for _, role := range []V1ConsumerRole{V1DaemonRole, V1ApprovalBrokerRole} {
		role := role
		t.Run(role.String(), func(t *testing.T) {
			t.Parallel()
			bindings := testV1Bindings(role)
			source := []byte("export const answer = 42;\n")
			request, err := NewV1ValidationRequest(role, bindings, source)
			if err != nil {
				t.Fatalf("construct request: %v", err)
			}
			encoded, err := EncodeV1Request(request)
			if err != nil {
				t.Fatalf("encode request: %v", err)
			}
			if len(encoded) != V1RequestHeaderBytes+len(source) {
				t.Fatalf("request bytes = %d", len(encoded))
			}
			decoded, err := DecodeV1Request(role, encoded)
			if err != nil {
				t.Fatalf("decode request: %v", err)
			}
			if decoded.Role != role || !bytes.Equal(decoded.SourceBytes(), source) {
				t.Fatal("decoded request lost its role or exact source")
			}

			source[0] ^= 0xff
			encoded[0] ^= 0xff
			returned := decoded.SourceBytes()
			returned[0] ^= 0xff
			if !bytes.Equal(decoded.SourceBytes(), []byte("export const answer = 42;\n")) {
				t.Fatal("request retained caller, wire, or accessor storage")
			}
		})
	}
}

func TestV1RequestExactCapAndCrossFamilyRefusal(t *testing.T) {
	t.Parallel()
	maximum := bytes.Repeat([]byte{'x'}, v0candidate.MJSMainSourceMaxBytes)
	request, err := NewV1ValidationRequest(V1DaemonRole, testV1Bindings(V1DaemonRole), maximum)
	if err != nil {
		t.Fatalf("exact maximum refused: %v", err)
	}
	encoded, err := EncodeV1Request(request)
	if err != nil || len(encoded) != V1RequestMaximumBytes {
		t.Fatalf("exact maximum encoding = %d, %v", len(encoded), err)
	}
	if _, err := DecodeV1Request(V1ApprovalBrokerRole, encoded); errClassification(err) != ClassificationDomain {
		t.Fatalf("cross-role classification = %q (%v)", errClassification(err), err)
	}
	if _, err := DecodeV1Request(V1DaemonRole, append(encoded, 0)); errClassification(err) != ClassificationMalformed {
		t.Fatalf("cap+1 classification = %q (%v)", errClassification(err), err)
	}
	if _, err := DecodeV1Request(V1DaemonRole, readV0RequestFixture(t)); errClassification(err) != ClassificationUnsupported {
		t.Fatalf("v0-as-v1 classification = %q (%v)", errClassification(err), err)
	}
}

func TestV1InactiveResourcePolicyRefusesInventedMeasurements(t *testing.T) {
	t.Parallel()
	for _, role := range []V1ConsumerRole{V1DaemonRole, V1ApprovalBrokerRole} {
		policy := NewV1InactiveResourcePolicy(role)
		encoded, err := EncodeV1ResourcePolicy(policy)
		if err != nil {
			t.Fatalf("encode inactive policy: %v", err)
		}
		if len(encoded) != V1ResourcePolicyFrameBytes {
			t.Fatalf("resource policy bytes = %d", len(encoded))
		}
		decoded, err := DecodeV1ResourcePolicy(role, encoded)
		if err != nil || decoded.Activation != V1PolicyInactive {
			t.Fatalf("decode inactive policy: %v", err)
		}

		policy.ObservedFootprintThreshold = 1
		if _, err := EncodeV1ResourcePolicy(policy); errClassification(err) != ClassificationSchema {
			t.Fatalf("invented measurement classification = %q (%v)", errClassification(err), err)
		}
	}
}

func TestV1ProfileAndConsumerFamiliesAreRoleClosed(t *testing.T) {
	t.Parallel()
	for _, role := range []V1ConsumerRole{V1DaemonRole, V1ApprovalBrokerRole} {
		set := testV1ProfileSet(t, role)
		processBytes, err := EncodeV1ProcessProfile(set.Process)
		if err != nil || len(processBytes) != V1ProcessProfileFrameBytes {
			t.Fatalf("process profile: %d, %v", len(processBytes), err)
		}
		artifactBytes, err := EncodeV1ArtifactProfile(set.Artifact)
		if err != nil || len(artifactBytes) != V1ArtifactProfileFrameBytes {
			t.Fatalf("artifact profile: %d, %v", len(artifactBytes), err)
		}
		consumerBytes, err := EncodeV1ConsumerProjection(set.Consumer)
		if err != nil || len(consumerBytes) != V1ConsumerProjectionFrameBytes {
			t.Fatalf("consumer projection: %d, %v", len(consumerBytes), err)
		}

		other := role.Other()
		if _, err := DecodeV1ProcessProfile(other, processBytes); errClassification(err) != ClassificationDomain {
			t.Fatalf("cross-role process profile classification = %q", errClassification(err))
		}
		if _, err := DecodeV1ArtifactProfile(other, artifactBytes); errClassification(err) != ClassificationDomain {
			t.Fatalf("cross-role artifact profile classification = %q", errClassification(err))
		}
		if _, err := DecodeV1ConsumerProjection(other, consumerBytes); errClassification(err) != ClassificationDomain {
			t.Fatalf("cross-role consumer classification = %q", errClassification(err))
		}
	}
}

func TestV1ResultBindsExactRequestProfilesAndCleanup(t *testing.T) {
	t.Parallel()
	for _, role := range []V1ConsumerRole{V1DaemonRole, V1ApprovalBrokerRole} {
		bindings := testV1Bindings(role)
		request, err := NewV1ValidationRequest(role, bindings, []byte("export default 1;\n"))
		if err != nil {
			t.Fatal(err)
		}
		result, err := NewV1ValidationResult(request, V1ParseValid, V1PolicyAllow, V1ResultNoFinding, [5]uint32{})
		if err != nil {
			t.Fatal(err)
		}
		encoded, err := EncodeV1Result(result)
		if err != nil || len(encoded) != V1ResultFrameBytes {
			t.Fatalf("result encoding: %d, %v", len(encoded), err)
		}
		verified, err := VerifyV1Result(request, encoded)
		if err != nil || verified.CleanupDisposition != V1CleanupComplete || verified.RefusalDisposition != V1RefusalNone {
			t.Fatalf("verify result: %v", err)
		}
		if _, err := VerifyV1Result(request, append(encoded, 0)); errClassification(err) != ClassificationMalformed {
			t.Fatalf("trailing result classification = %q (%v)", errClassification(err), err)
		}
		if _, err := DecodeV1Result(role.Other(), encoded); errClassification(err) != ClassificationDomain {
			t.Fatalf("cross-role result classification = %q (%v)", errClassification(err), err)
		}
	}
}

type v1TestProfileSet struct {
	Process  *V1ProcessProfile
	Artifact *V1ArtifactProfile
	Consumer *V1ConsumerProjection
}

func testV1Bindings(role V1ConsumerRole) V1RequestBindings {
	var bindings V1RequestBindings
	bindings.InstallationID[0] = 0x11
	bindings.EpochSequence = 7
	bindings.EpochDigest[0] = 0x22
	bindings.ArtifactProfileDigest[0] = byte(0x30 + role)
	bindings.ResourcePolicyDigest[0] = byte(0x40 + role)
	bindings.CorrelationID[0] = byte(0x50 + role)
	return bindings
}

func testV1ProfileSet(t *testing.T, role V1ConsumerRole) v1TestProfileSet {
	t.Helper()
	policy := NewV1InactiveResourcePolicy(role)
	policyBytes, err := EncodeV1ResourcePolicy(policy)
	if err != nil {
		t.Fatal(err)
	}
	policyDigest := sha256.Sum256(policyBytes)
	process := NewFixtureOnlyV1ProcessProfile(role, policyDigest)
	processBytes, err := EncodeV1ProcessProfile(process)
	if err != nil {
		t.Fatal(err)
	}
	processDigest := sha256.Sum256(processBytes)
	artifact := NewFixtureOnlyV1ArtifactProfile(role, processDigest, policyDigest)
	artifactBytes, err := EncodeV1ArtifactProfile(artifact)
	if err != nil {
		t.Fatal(err)
	}
	artifactDigest := sha256.Sum256(artifactBytes)
	consumer := NewFixtureOnlyV1ConsumerProjection(role, artifactDigest, policyDigest)
	return v1TestProfileSet{Process: process, Artifact: artifact, Consumer: consumer}
}

func errClassification(err error) string {
	value, _ := ErrorClassification(err)
	return value
}

func readV0RequestFixture(t *testing.T) []byte {
	t.Helper()
	request, err := NewValidationRequest(ValidationRequestID{1}, []byte("export default 1;\n"))
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	return encoded
}
