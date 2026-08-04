package macosplan

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
)

const FixtureDirectory = "internal/installation/macosplan/testdata"

type InventoryFixtureCase struct {
	ID       string               `json:"id"`
	Input    InventoryObservation `json:"input"`
	Expected Result               `json:"expected"`
}

type BootstrapFixtureCase struct {
	ID       string          `json:"id"`
	Input    BootstrapInput  `json:"input"`
	Expected BootstrapOutput `json:"expected"`
}

type UpdateFixtureCase struct {
	ID       string       `json:"id"`
	Input    UpdateInput  `json:"input"`
	Expected UpdateOutput `json:"expected"`
}

type RepairFixtureCase struct {
	ID       string       `json:"id"`
	Input    RepairInput  `json:"input"`
	Expected RepairOutput `json:"expected"`
}

type UninstallFixtureCase struct {
	ID       string          `json:"id"`
	Input    UninstallMode   `json:"input"`
	Expected UninstallOutput `json:"expected"`
}

type FixtureFile struct {
	Path   string `json:"path"`
	Bytes  int    `json:"bytes"`
	SHA256 string `json:"sha256"`
}

type FixtureManifest struct {
	ManifestVersion string        `json:"manifestVersion"`
	ProfileID       string        `json:"profileId"`
	Status          string        `json:"status"`
	Files           []FixtureFile `json:"files"`
}

type GeneratedFixture struct {
	Path  string
	Bytes []byte
}

func GenerateFixtures() ([]GeneratedFixture, error) {
	profile := CanonicalProfile()
	files := []GeneratedFixture{}
	add := func(name string, value any) error {
		encoded, err := canonicalFixtureJSON(value)
		if err != nil {
			return err
		}
		files = append(files, GeneratedFixture{Path: FixtureDirectory + "/" + name, Bytes: encoded})
		return nil
	}

	if err := add("profile.json", profile); err != nil {
		return nil, err
	}
	if err := add("inventory-cases.json", inventoryFixtureCases(profile)); err != nil {
		return nil, err
	}
	if err := add("bootstrap-cases.json", bootstrapFixtureCases()); err != nil {
		return nil, err
	}
	if err := add("update-cases.json", updateFixtureCases(profile)); err != nil {
		return nil, err
	}
	if err := add("repair-cases.json", repairFixtureCases()); err != nil {
		return nil, err
	}
	if err := add("uninstall-cases.json", uninstallFixtureCases()); err != nil {
		return nil, err
	}

	sort.Slice(files, func(left, right int) bool { return files[left].Path < files[right].Path })
	manifest := FixtureManifest{
		ManifestVersion: "capsule.macos-installation-i0-fixtures/v0",
		ProfileID:       ProfileIdentity,
		Status:          "passive-inactive-no-side-effects",
	}
	for _, file := range files {
		digest := sha256.Sum256(file.Bytes)
		manifest.Files = append(manifest.Files, FixtureFile{
			Path:   file.Path,
			Bytes:  len(file.Bytes),
			SHA256: hex.EncodeToString(digest[:]),
		})
	}
	manifestBytes, err := canonicalFixtureJSON(manifest)
	if err != nil {
		return nil, err
	}
	files = append(files, GeneratedFixture{Path: FixtureDirectory + "/manifest.json", Bytes: manifestBytes})
	return files, nil
}

func inventoryFixtureCases(profile Profile) []InventoryFixtureCase {
	exact := InventoryObservation{ProfileID: profile.ProfileID, Roles: append([]RoleProjection(nil), profile.Roles...)}
	missing := InventoryObservation{ProfileID: profile.ProfileID, Roles: append([]RoleProjection(nil), profile.Roles[:len(profile.Roles)-1]...)}
	extra := InventoryObservation{ProfileID: profile.ProfileID, Roles: append([]RoleProjection(nil), profile.Roles...)}
	extra.Roles = append(extra.Roles, role(
		"capsule.role.runner", BrokerRoleID, PackagingEmbeddedAgent, VisibilityEmbedded,
		ApplicationPath+"/Contents/Library/Helpers/CapsuleRunner.app",
		ApplicationPath+"/Contents/Library/Helpers/CapsuleRunner.app/Contents/MacOS/CapsuleRunner",
		"com.capsulecorp.capsule.runner", "com.capsulecorp.capsule.runner", "com.capsulecorp.capsule.runner",
	))
	mixed := InventoryObservation{ProfileID: profile.ProfileID, Roles: append([]RoleProjection(nil), profile.Roles...)}
	mixed.Roles[3].ServiceIdentifier = BrokerValidatorServiceIdentity
	duplicate := InventoryObservation{ProfileID: profile.ProfileID, Roles: append([]RoleProjection(nil), profile.Roles...)}
	duplicate.Roles[4].RoleID = duplicate.Roles[3].RoleID
	emptyRoleID := InventoryObservation{ProfileID: profile.ProfileID, Roles: append([]RoleProjection(nil), profile.Roles...)}
	emptyRoleID.Roles[5].RoleID = ""
	profileMismatch := InventoryObservation{ProfileID: profile.ProfileID + "-mismatch", Roles: append([]RoleProjection(nil), profile.Roles...)}

	return []InventoryFixtureCase{
		{ID: "inventory-exact-activation-refuses-inactive-signing", Input: exact, Expected: EvaluateActivation(profile, exact)},
		{ID: "inventory-missing-refuses", Input: missing, Expected: EvaluateActivation(profile, missing)},
		{ID: "inventory-extra-refuses", Input: extra, Expected: EvaluateActivation(profile, extra)},
		{ID: "inventory-mixed-refuses", Input: mixed, Expected: EvaluateActivation(profile, mixed)},
		{ID: "inventory-duplicate-role-id-refuses", Input: duplicate, Expected: EvaluateActivation(profile, duplicate)},
		{ID: "inventory-empty-role-id-refuses", Input: emptyRoleID, Expected: EvaluateActivation(profile, emptyRoleID)},
		{ID: "inventory-profile-id-mismatch-refuses", Input: profileMismatch, Expected: EvaluateActivation(profile, profileMismatch)},
	}
}

func bootstrapFixtureCases() []BootstrapFixtureCase {
	inputs := []struct {
		id    string
		input BootstrapInput
	}{
		{
			id: "begin-creates-disabled-journal-state",
			input: BootstrapInput{RecordVersion: 0, ProfileID: ProfileIdentity, State: BootstrapAbsent, Event: BootstrapBegin,
				InventoryDecision: DecisionAccept, SigningState: ProfileInactiveRefusing, BootstrapAuthorityState: ProfileInactiveRefusing},
		},
		{
			id: "identity-step-refuses-inactive-bootstrap-authority",
			input: BootstrapInput{RecordVersion: 0, ProfileID: ProfileIdentity, State: BootstrapJournaledDisabled, Event: BootstrapRecordIdentity,
				InventoryDecision: DecisionAccept, SigningState: ProfileInactiveRefusing, BootstrapAuthorityState: ProfileInactiveRefusing},
		},
		{
			id: "attempts-enabled-before-ready-enters-repair",
			input: BootstrapInput{RecordVersion: 0, ProfileID: ProfileIdentity, State: BootstrapServicesPreparedDisabled, Event: BootstrapEnrollProtectedRoot,
				InventoryDecision: DecisionAccept, SigningState: ProfileActive, BootstrapAuthorityState: ProfileActive, AttemptsEnabled: true},
		},
		{
			id: "clean-owner-store-recovery-is-only-ready-transition",
			input: BootstrapInput{RecordVersion: 0, ProfileID: ProfileIdentity, State: BootstrapEpochCommittedRecovering, Event: BootstrapCompleteRecovery,
				InventoryDecision: DecisionAccept, SigningState: ProfileActive, BootstrapAuthorityState: ProfileActive,
				ProtectedRootState: "enrolled-exact", OwnerLockState: "enrolled-exact", StoreState: "open-without-create-exact",
				EpochState: "epoch-one-prepared-exact", RecoveryState: "owner-store-recovery-clean"},
		},
	}
	cases := make([]BootstrapFixtureCase, 0, len(inputs))
	for _, candidate := range inputs {
		cases = append(cases, BootstrapFixtureCase{ID: candidate.id, Input: candidate.input, Expected: EvaluateBootstrap(candidate.input)})
	}
	return cases
}

func updateFixtureCases(profile Profile) []UpdateFixtureCase {
	identity := UpdateIdentityForProfile(profile)
	exactInventory := accept("component-inventory-exact")
	exact := UpdateInput{
		RecordVersion: 0, Predecessor: identity, Successor: identity, CompleteBundleObserved: true,
		StateRootIdentityPreserved: true, OwnerLockIdentityPreserved: true, BothValidatorTuplesExact: true,
		ServicesStopped: true, ComponentInventoryResult: exactInventory,
	}
	mixed := exact
	mixed.Successor.RoleSetDigest = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	incomplete := exact
	incomplete.CompleteBundleObserved = false

	return []UpdateFixtureCase{
		{ID: "exact-passive-identity-refuses-inactive-signing", Input: exact, Expected: EvaluateUpdate(exact)},
		{ID: "mixed-role-set-refuses", Input: mixed, Expected: EvaluateUpdate(mixed)},
		{ID: "incomplete-whole-bundle-refuses", Input: incomplete, Expected: EvaluateUpdate(incomplete)},
	}
}

func repairFixtureCases() []RepairFixtureCase {
	candidates := []struct {
		id    string
		input RepairInput
	}{
		{"restore-application-files", RepairInput{RecordVersion: 0, InstallationRootAvailable: true, StateRootContinuityProven: true, OwnerLockContinuityProven: true, StoreAndHistoryComplete: true, CleanupObligationsResolved: true, CurrentEpochExact: true}},
		{"authorized-forward-transition", RepairInput{RecordVersion: 0, ApplicationFilesExact: true, InstallationRootAvailable: true, StateRootContinuityProven: true, OwnerLockContinuityProven: true, StoreAndHistoryComplete: true, CleanupObligationsResolved: true, ForwardTransitionAuthorized: true}},
		{"protected-state-repair", RepairInput{RecordVersion: 0, ApplicationFilesExact: true, InstallationRootAvailable: true, StoreAndHistoryComplete: true, CleanupObligationsResolved: true, CurrentEpochExact: true}},
		{"new-installation", RepairInput{RecordVersion: 0}},
		{"refuse-history-loss", RepairInput{RecordVersion: 0, ApplicationFilesExact: true, InstallationRootAvailable: true, StateRootContinuityProven: true, OwnerLockContinuityProven: true, CurrentEpochExact: true}},
	}
	cases := make([]RepairFixtureCase, 0, len(candidates))
	for _, candidate := range candidates {
		cases = append(cases, RepairFixtureCase{ID: candidate.id, Input: candidate.input, Expected: ClassifyRepair(candidate.input)})
	}
	return cases
}

func uninstallFixtureCases() []UninstallFixtureCase {
	modes := []UninstallMode{UninstallApplicationOnly, UninstallLocalWhereSafe, UninstallAbandonIdentity}
	cases := make([]UninstallFixtureCase, 0, len(modes))
	for _, mode := range modes {
		cases = append(cases, UninstallFixtureCase{ID: string(mode), Input: mode, Expected: ClassifyUninstall(mode)})
	}
	return cases
}

func canonicalFixtureJSON(value any) ([]byte, error) {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode fixture: %w", err)
	}
	return append(encoded, '\n'), nil
}
