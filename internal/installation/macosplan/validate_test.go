package macosplan

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestCanonicalProfileClosesOneVisibleNoGuestRoleTree(t *testing.T) {
	profile := CanonicalProfile()
	if result := ValidateProfile(profile); result.Decision != DecisionAccept {
		t.Fatalf("canonical profile refused: %+v", result)
	}
	if profile.VisibleRoleID != BrokerRoleID || profile.ApplicationPath != "Capsule.app" {
		t.Fatalf("visible application = %q %q", profile.VisibleRoleID, profile.ApplicationPath)
	}
	if len(profile.Roles) != 7 {
		t.Fatalf("role count = %d, want 7", len(profile.Roles))
	}
	visible := 0
	for _, role := range profile.Roles {
		if role.Visibility == VisibilityUserVisible {
			visible++
		}
		if !role.Required {
			t.Fatalf("role %q is optional", role.RoleID)
		}
	}
	if visible != 1 {
		t.Fatalf("visible role count = %d, want 1", visible)
	}
	for _, excluded := range profile.ExcludedRoles {
		for _, role := range profile.Roles {
			if excluded == role.RoleID {
				t.Fatalf("excluded role %q is present", excluded)
			}
		}
	}
	if profile.Signing.State != ProfileInactiveRefusing || profile.Signing.TeamIdentifier != "" ||
		profile.Signing.ProvisioningProfileSetID != "" {
		t.Fatalf("signing profile invented active identity: %+v", profile.Signing)
	}
	if profile.Bootstrap.State != ProfileInactiveRefusing || profile.Bootstrap.BootstrapOwner == "" {
		t.Fatalf("bootstrap profile does not retain its refusing owner decision: %+v", profile.Bootstrap)
	}
}

func TestCanonicalProfileComposesADRServiceAndValidatorIdentities(t *testing.T) {
	profile := CanonicalProfile()
	want := []string{
		DaemonBundleIdentity,
		SupervisorBundleIdentity,
		SupervisorDaemonServiceIdentity,
		SupervisorBrokerServiceIdentity,
		DaemonValidatorServiceIdentity,
		BrokerValidatorServiceIdentity,
	}
	if len(profile.Services) != len(want) {
		t.Fatalf("service count = %d, want %d", len(profile.Services), len(want))
	}
	for index, identity := range want {
		if profile.Services[index].ServiceName != identity {
			t.Fatalf("service[%d] = %q, want %q", index, profile.Services[index].ServiceName, identity)
		}
		if profile.Services[index].ActivationState != ProfileInactiveRefusing {
			t.Fatalf("service %q unexpectedly active", identity)
		}
	}
	if profile.Roles[4].SigningIdentifier != DaemonValidatorParserIdentity ||
		profile.Roles[6].SigningIdentifier != BrokerValidatorParserIdentity {
		t.Fatal("role-specific parser identities are not closed")
	}
	if profile.Roles[1].ServiceDescriptorPath != "Capsule.app/Contents/Library/LaunchAgents/com.capsulecorp.capsule.daemon.plist" ||
		profile.Roles[2].ServiceDescriptorPath != "Capsule.app/Contents/Library/LaunchAgents/com.capsulecorp.capsule.supervisor.plist" {
		t.Fatal("per-user service descriptor paths are not closed")
	}
	if profile.Roles[3].ResourcePolicyPath == "" || profile.Roles[5].ResourcePolicyPath == "" {
		t.Fatal("R2 inactive resource-policy paths are not closed")
	}
}

func TestCanonicalProfileReturnsDefensiveCopies(t *testing.T) {
	left := CanonicalProfile()
	left.Roles[0].RoleID = "mutated"
	left.Entitlements[0].Key = "mutated"
	left.Services[0].MethodIdentities = "mutated"
	left.ExcludedRoles[0] = "mutated"
	right := CanonicalProfile()
	if result := ValidateProfile(right); result.Decision != DecisionAccept {
		t.Fatalf("fresh profile changed through returned slices: %+v", result)
	}
}

func TestInventoryRefusesMissingMixedExtraBeforeInactiveSigning(t *testing.T) {
	profile := CanonicalProfile()
	tests := inventoryFixtureCases(profile)
	wantReasons := []string{
		"signing-profile-inactive",
		"component-missing",
		"component-extra",
		"component-mixed",
		"component-mixed-or-duplicate",
		"component-mixed-or-duplicate",
		"inventory-profile-mismatch",
	}
	for index, test := range tests {
		result := EvaluateActivation(profile, test.Input)
		if result.Decision != DecisionRefuse || result.Reason != wantReasons[index] {
			t.Fatalf("%s = %+v, want reason %q", test.ID, result, wantReasons[index])
		}
	}
}

func TestBootstrapKeepsAttemptsDisabledAndRefusesInventedActiveProfile(t *testing.T) {
	begin := BootstrapInput{RecordVersion: 0, ProfileID: ProfileIdentity, State: BootstrapAbsent, Event: BootstrapBegin, InventoryDecision: DecisionAccept}
	if output := EvaluateBootstrap(begin); output.State != BootstrapJournaledDisabled || output.AttemptsEnabled || output.Result.Decision != DecisionAccept {
		t.Fatalf("begin = %+v", output)
	}

	early := begin
	early.State = BootstrapServicesPreparedDisabled
	early.Event = BootstrapEnrollProtectedRoot
	early.AttemptsEnabled = true
	if output := EvaluateBootstrap(early); output.State != BootstrapRepairRequired || output.Result.Reason != "attempts-enabled-before-ready" {
		t.Fatalf("early attempts = %+v", output)
	}

	ready := BootstrapInput{
		RecordVersion: 0, ProfileID: ProfileIdentity, State: BootstrapEpochCommittedRecovering,
		Event: BootstrapCompleteRecovery, InventoryDecision: DecisionAccept, SigningState: ProfileActive,
		BootstrapAuthorityState: ProfileActive, RecoveryState: "owner-store-recovery-clean",
	}
	if output := EvaluateBootstrap(ready); output.State != BootstrapRepairRequired || output.AttemptsEnabled || output.Result.Reason != "i0-active-bootstrap-profile-unsupported" {
		t.Fatalf("invented active profile = %+v", output)
	}
}

func TestUpdateIdentityAndInactiveReplacementRefusal(t *testing.T) {
	profile := CanonicalProfile()
	identity := UpdateIdentityForProfile(profile)
	if identity.RoleSetDigest == "" || identity.ServiceSetDigest == "" || identity.EntitlementsSetDigest == "" {
		t.Fatal("update compatibility identities are incomplete")
	}
	input := UpdateInput{
		RecordVersion: 0, Predecessor: identity, Successor: identity, CompleteBundleObserved: true,
		StateRootIdentityPreserved: true, OwnerLockIdentityPreserved: true, BothValidatorTuplesExact: true,
		ServicesStopped: true, ComponentInventoryResult: accept("component-inventory-exact"),
	}
	if output := EvaluateUpdate(input); output.Result.Reason != "update-signing-profile-inactive" || output.AttemptsState != "disabled" {
		t.Fatalf("inactive update = %+v", output)
	}
	input.Successor.ServiceSetDigest = "mismatch"
	if output := EvaluateUpdate(input); output.Result.Reason != "update-identity-mixed" {
		t.Fatalf("mixed update = %+v", output)
	}
	input.Successor = identity
	input.Predecessor.SigningState = ProfileActive
	input.Successor.SigningState = ProfileActive
	input.Predecessor.ReleaseIdentity = "fixture-release"
	input.Successor.ReleaseIdentity = "fixture-release"
	if output := EvaluateUpdate(input); output.Result.Reason != "i0-active-update-profile-unsupported" {
		t.Fatalf("invented active update = %+v", output)
	}
}

func TestRepairClassificationsNeverRewriteSecurityHistory(t *testing.T) {
	for _, fixture := range repairFixtureCases() {
		output := ClassifyRepair(fixture.Input)
		if !reflect.DeepEqual(output, fixture.Expected) || output.AttemptsState != "disabled" {
			t.Fatalf("%s = %+v", fixture.ID, output)
		}
	}
	refusal := ClassifyRepair(RepairInput{
		RecordVersion: 0, InstallationRootAvailable: true, StateRootContinuityProven: true,
		OwnerLockContinuityProven: true, CurrentEpochExact: true,
	})
	if refusal.Class != RepairRefuseAutomatic || refusal.Result.Reason != "history-or-cleanup-unresolved" {
		t.Fatalf("history refusal = %+v", refusal)
	}
}

func TestUninstallKeepsRequiredHistoryAndExternalEvidenceOutOfScope(t *testing.T) {
	for _, mode := range []UninstallMode{UninstallApplicationOnly, UninstallLocalWhereSafe, UninstallAbandonIdentity} {
		output := ClassifyUninstall(mode)
		if output.Result.Decision != DecisionAccept || output.AttemptsState != "disabled" {
			t.Fatalf("%s = %+v", mode, output)
		}
		dispositions := make(map[string]DataDisposition, len(output.Dispositions))
		for _, disposition := range output.Dispositions {
			dispositions[disposition.DataClass] = disposition.Disposition
		}
		if dispositions["application-bundle"] != DispositionRemove ||
			dispositions["grant-attempt-and-cleanup-history"] != DispositionRetain ||
			dispositions["archive-and-nonreuse-history"] != DispositionRetain ||
			dispositions["exported-receipts-external-witnesses-and-backups"] != DispositionOutsideScope {
			t.Fatalf("%s dispositions = %+v", mode, dispositions)
		}
	}
}

func TestGeneratedFixturesAreCurrentAndDigestBound(t *testing.T) {
	root := repositoryRootForTest(t)
	generated, err := GenerateFixtures()
	if err != nil {
		t.Fatal(err)
	}
	for _, fixture := range generated {
		path := filepath.Join(root, filepath.FromSlash(fixture.Path))
		stored, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", fixture.Path, err)
		}
		if !reflect.DeepEqual(stored, fixture.Bytes) {
			t.Fatalf("generated fixture is stale: %s", fixture.Path)
		}
	}

	manifestPath := filepath.Join(root, FixtureDirectory, "manifest.json")
	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var manifest FixtureManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.ProfileID != ProfileIdentity || manifest.Status != "passive-inactive-no-side-effects" {
		t.Fatalf("fixture manifest identity = %+v", manifest)
	}
	for _, file := range manifest.Files {
		contents, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(file.Path)))
		if err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256(contents)
		if len(contents) != file.Bytes || hex.EncodeToString(digest[:]) != file.SHA256 {
			t.Fatalf("fixture manifest mismatch: %s", file.Path)
		}
	}
}

func repositoryRootForTest(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller unavailable")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "../../.."))
}
