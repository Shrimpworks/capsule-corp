package v0candidate

import (
	"bytes"
	"testing"
)

func TestCandidateIdentifiersRejectInvalidValuesAndCopyInput(t *testing.T) {
	source := bytes.Repeat([]byte{0x11}, 16)
	installationID, err := NewInstallationID(source)
	if err != nil {
		t.Fatalf("create installation ID: %v", err)
	}
	source[0] = 0xff

	if installationID[0] != 0x11 {
		t.Fatal("installation ID retained caller-owned storage")
	}
	if _, err := NewRegistrationID(make([]byte, 15)); err == nil {
		t.Fatal("expected a short registration ID to be rejected")
	}
	if _, err := NewSupervisorID(make([]byte, 16)); err == nil {
		t.Fatal("expected an all-zero Supervisor ID to be rejected")
	}
}

func TestCandidateScalarsRejectOutOfRangeValues(t *testing.T) {
	if _, err := NewUInt53(MaxSafeInteger + 1); err == nil {
		t.Fatal("expected an unsafe integer to be rejected")
	}
	if _, err := NewPositiveUInt53(0); err == nil {
		t.Fatal("expected zero to be rejected for a positive integer")
	}
	if _, err := NewExecutionPlanDigest(make([]byte, 31)); err == nil {
		t.Fatal("expected a short execution-plan digest to be rejected")
	}
}

func TestPlanRegistrationKeepsIdentityAndDigestDomainsDistinct(t *testing.T) {
	registrationID, err := NewRegistrationID(bytes.Repeat([]byte{0x77}, 16))
	if err != nil {
		t.Fatalf("create registration ID: %v", err)
	}
	installationID, err := NewInstallationID(bytes.Repeat([]byte{0x11}, 16))
	if err != nil {
		t.Fatalf("create installation ID: %v", err)
	}
	supervisorID, err := NewSupervisorID(bytes.Repeat([]byte{0x55}, 16))
	if err != nil {
		t.Fatalf("create Supervisor ID: %v", err)
	}
	planDigest, err := NewExecutionPlanDigest(bytes.Repeat([]byte{0x88}, 32))
	if err != nil {
		t.Fatalf("create execution-plan digest: %v", err)
	}
	epochDigest, err := NewTrustEpochDigest(bytes.Repeat([]byte{0x22}, 32))
	if err != nil {
		t.Fatalf("create trust-epoch digest: %v", err)
	}

	registration := PlanRegistration{
		ObjectType:           PlanRegistrationObjectType,
		ObjectVersion:        CandidateObjectVersion,
		RegistrationID:       registrationID,
		RegistrationSequence: 1,
		PlanDigest:           planDigest,
		InstallationID:       installationID,
		EpochSequence:        7,
		EpochDigest:          epochDigest,
		SupervisorID:         supervisorID,
		ExpiresAt:            1_785_456_300,
	}

	if registration.ObjectType != "capsule.plan-registration" {
		t.Fatalf("unexpected object type %q", registration.ObjectType)
	}
	if registration.PlanDigest == (ExecutionPlanDigest{}) {
		t.Fatal("expected a nonzero plan digest")
	}
}
