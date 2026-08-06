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

// TestNonCanonicalProfileIsRefusedAtEveryEntryPoint covers ValidateProfile's
// own refuse branch and its delegation through EvaluateInventory and
// EvaluateActivation. Every other test in this file evaluates the exact
// CanonicalProfile(), which never exercises this path.
func TestNonCanonicalProfileIsRefusedAtEveryEntryPoint(t *testing.T) {
	tampered := CanonicalProfile()
	tampered.ProfileID = tampered.ProfileID + "-tampered"
	observed := InventoryObservation{ProfileID: tampered.ProfileID, Roles: append([]RoleProjection(nil), tampered.Roles...)}

	if result := ValidateProfile(tampered); result.Decision != DecisionRefuse || result.Reason != "profile-not-canonical-i0" {
		t.Fatalf("ValidateProfile = %+v", result)
	}
	if result := EvaluateInventory(tampered, observed); result.Decision != DecisionRefuse || result.Reason != "profile-not-canonical-i0" {
		t.Fatalf("EvaluateInventory = %+v", result)
	}
	if result := EvaluateActivation(tampered, observed); result.Decision != DecisionRefuse || result.Reason != "profile-not-canonical-i0" {
		t.Fatalf("EvaluateActivation = %+v", result)
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

// TestEvaluateBootstrapCoversEveryReachableRefusalAndTheReadyTransition covers
// the EvaluateBootstrap branches TestBootstrapKeepsAttemptsDisabledAndRefusesInventedActiveProfile
// does not reach: every event-specific and structural refusal, plus the one
// path that actually lands on BootstrapReady with attempts enabled. (The
// i0-active-bootstrap-profile-unsupported guard blocks nothing here, since
// none of these fixtures assert an active signing/bootstrap-authority
// profile — that posture-guard case is already covered separately.)
func TestEvaluateBootstrapCoversEveryReachableRefusalAndTheReadyTransition(t *testing.T) {
	base := BootstrapInput{
		RecordVersion: 0, ProfileID: ProfileIdentity, InventoryDecision: DecisionAccept,
		SigningState: ProfileInactiveRefusing, BootstrapAuthorityState: ProfileInactiveRefusing,
	}

	wrongVersion := base
	wrongVersion.RecordVersion = 1
	wrongVersion.State, wrongVersion.Event = BootstrapAbsent, BootstrapBegin
	if output := EvaluateBootstrap(wrongVersion); output.State != BootstrapRepairRequired || output.Result.Reason != "bootstrap-input-profile" {
		t.Fatalf("wrong record version = %+v", output)
	}

	wrongProfile := base
	wrongProfile.ProfileID = ProfileIdentity + "-mismatch"
	wrongProfile.State, wrongProfile.Event = BootstrapAbsent, BootstrapBegin
	if output := EvaluateBootstrap(wrongProfile); output.State != BootstrapRepairRequired || output.Result.Reason != "bootstrap-input-profile" {
		t.Fatalf("wrong profile id = %+v", output)
	}

	inventoryRefused := base
	inventoryRefused.State, inventoryRefused.Event = BootstrapAbsent, BootstrapBegin
	inventoryRefused.InventoryDecision = DecisionRefuse
	if output := EvaluateBootstrap(inventoryRefused); output.State != BootstrapRepairRequired || output.Result.Reason != "bootstrap-component-inventory-refused" {
		t.Fatalf("refused inventory = %+v", output)
	}

	invalidTransition := base
	invalidTransition.State, invalidTransition.Event = BootstrapReady, BootstrapBegin
	if output := EvaluateBootstrap(invalidTransition); output.State != BootstrapRepairRequired || output.Result.Reason != "bootstrap-transition-invalid" {
		t.Fatalf("invalid transition = %+v", output)
	}

	authorityInactive := base
	authorityInactive.State, authorityInactive.Event = BootstrapJournaledDisabled, BootstrapRecordIdentity
	if output := EvaluateBootstrap(authorityInactive); output.State != BootstrapRepairRequired || output.Result.Reason != "bootstrap-authority-inactive" {
		t.Fatalf("inactive bootstrap authority = %+v", output)
	}

	signingInactive := base
	signingInactive.State, signingInactive.Event = BootstrapIdentityPreparedDisabled, BootstrapRegisterServices
	if output := EvaluateBootstrap(signingInactive); output.State != BootstrapRepairRequired || output.Result.Reason != "signing-profile-inactive" {
		t.Fatalf("inactive signing for register-services = %+v", output)
	}

	rootNotExact := base
	rootNotExact.State, rootNotExact.Event = BootstrapServicesPreparedDisabled, BootstrapEnrollProtectedRoot
	if output := EvaluateBootstrap(rootNotExact); output.State != BootstrapRepairRequired || output.Result.Reason != "protected-root-owner-store-not-exact" {
		t.Fatalf("protected root not exact = %+v", output)
	}

	epochNotExact := base
	epochNotExact.State, epochNotExact.Event = BootstrapComponentsVerifiedDisabled, BootstrapCommitEpoch
	if output := EvaluateBootstrap(epochNotExact); output.State != BootstrapRepairRequired || output.Result.Reason != "epoch-one-not-exact" {
		t.Fatalf("epoch not exact = %+v", output)
	}

	recoveryNotClean := base
	recoveryNotClean.State, recoveryNotClean.Event = BootstrapEpochCommittedRecovering, BootstrapCompleteRecovery
	if output := EvaluateBootstrap(recoveryNotClean); output.State != BootstrapRepairRequired || output.Result.Reason != "owner-store-recovery-not-clean" {
		t.Fatalf("recovery not clean = %+v", output)
	}

	ready := base
	ready.State, ready.Event = BootstrapEpochCommittedRecovering, BootstrapCompleteRecovery
	ready.RecoveryState = "owner-store-recovery-clean"
	output := EvaluateBootstrap(ready)
	if output.State != BootstrapReady || !output.AttemptsEnabled || output.Result != accept("bootstrap-transition-exact") {
		t.Fatalf("clean recovery should reach ready = %+v", output)
	}
}

// TestEvaluateUpdateCoversInputAndInventoryRefusalsBeforeTheIdentityGate
// covers the three EvaluateUpdate refusals that precede its
// Predecessor/Successor-vs-canonical comparison. Every branch after that
// comparison requires an active signing profile to be reached, which is
// mutually exclusive with matching CanonicalProfile()'s permanently
// inactive I0 posture — those branches are exercised, along with this
// package's own i0-active-update-profile-unsupported guard, by
// TestUpdateIdentityAndInactiveReplacementRefusal instead.
func TestEvaluateUpdateCoversInputAndInventoryRefusalsBeforeTheIdentityGate(t *testing.T) {
	wrongVersion := UpdateInput{RecordVersion: 1}
	if output := EvaluateUpdate(wrongVersion); output.Result.Reason != "update-input-version" || output.Compatibility != "refused" {
		t.Fatalf("wrong record version = %+v", output)
	}

	inventoryRefused := UpdateInput{RecordVersion: 0, ComponentInventoryResult: refuse(ClassificationBinding, "component-mixed")}
	if output := EvaluateUpdate(inventoryRefused); output.Result.Reason != "update-component-inventory-refused" {
		t.Fatalf("refused inventory = %+v", output)
	}

	incomplete := UpdateInput{RecordVersion: 0, ComponentInventoryResult: accept("component-inventory-exact")}
	if output := EvaluateUpdate(incomplete); output.Result.Reason != "update-incomplete-or-mixed" {
		t.Fatalf("incomplete bundle = %+v", output)
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

	wrongVersion := ClassifyRepair(RepairInput{RecordVersion: 1})
	if wrongVersion.Class != RepairRefuseAutomatic || wrongVersion.Result.Reason != "repair-input-version" {
		t.Fatalf("wrong record version = %+v", wrongVersion)
	}

	epochMismatch := ClassifyRepair(RepairInput{
		RecordVersion: 0, InstallationRootAvailable: true, StateRootContinuityProven: true,
		OwnerLockContinuityProven: true, StoreAndHistoryComplete: true, CleanupObligationsResolved: true,
	})
	if epochMismatch.Class != RepairRefuseAutomatic || epochMismatch.Result.Reason != "epoch-mismatch-without-forward-authority" {
		t.Fatalf("epoch mismatch without forward authority = %+v", epochMismatch)
	}

	restoreFiles := ClassifyRepair(RepairInput{
		RecordVersion: 0, InstallationRootAvailable: true, StateRootContinuityProven: true,
		OwnerLockContinuityProven: true, StoreAndHistoryComplete: true, CleanupObligationsResolved: true,
		CurrentEpochExact: true,
	})
	if restoreFiles.Class != RepairRestoreApplicationFiles || restoreFiles.Result != accept("restore-current-release-preserve-state") {
		t.Fatalf("restore application files = %+v", restoreFiles)
	}
}

// TestDigestJSONPanicsOnUnmarshalableInput guards digestJSON's deliberate
// panic instead of silently swallowing an encoding error: every current
// caller passes it a package-controlled struct/slice that always marshals,
// so if that ever stops holding, this should fail loudly, not return a
// digest of the empty string.
func TestDigestJSONPanicsOnUnmarshalableInput(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("digestJSON did not panic on unmarshalable input")
		}
	}()
	digestJSON(make(chan int))
}

// TestResultValidateAcceptsWellFormedTuplesAndRejectsMalformedOnes covers
// Result.Validate, the invariant checker every accept/refuse tuple in this
// package is meant to satisfy: an accept must carry ClassificationNone and a
// reason, a refuse must carry a non-none classification and a reason, and
// every other combination is malformed.
func TestResultValidateAcceptsWellFormedTuplesAndRejectsMalformedOnes(t *testing.T) {
	if err := accept("component-inventory-exact").Validate(); err != nil {
		t.Fatalf("well-formed accept rejected: %v", err)
	}
	if err := refuse(ClassificationBinding, "component-missing").Validate(); err != nil {
		t.Fatalf("well-formed refuse rejected: %v", err)
	}

	malformed := []Result{
		{Decision: DecisionAccept, Classification: ClassificationBinding, Reason: "component-missing"},
		{Decision: DecisionAccept, Classification: ClassificationNone, Reason: ""},
		{Decision: DecisionRefuse, Classification: ClassificationNone, Reason: "component-missing"},
		{Decision: DecisionRefuse, Classification: ClassificationBinding, Reason: ""},
		{Decision: "invented-decision", Classification: ClassificationNone, Reason: "x"},
	}
	for index, result := range malformed {
		if err := result.Validate(); err == nil {
			t.Fatalf("malformed[%d] %+v accepted", index, result)
		}
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
	if output := ClassifyUninstall(UninstallMode("invented-mode")); output.Result.Reason != "uninstall-mode" || len(output.Dispositions) != 0 {
		t.Fatalf("invalid mode = %+v", output)
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
