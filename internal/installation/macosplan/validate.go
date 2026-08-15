package macosplan

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"reflect"
)

func ValidateProfile(profile Profile) Result {
	if !reflect.DeepEqual(profile, CanonicalProfile()) {
		return refuse(ClassificationBinding, "profile-not-canonical-i0")
	}
	return accept("canonical-passive-profile")
}

func EvaluateInventory(profile Profile, observed InventoryObservation) Result {
	if result := ValidateProfile(profile); result.Decision != DecisionAccept {
		return result
	}
	if observed.ProfileID != profile.ProfileID {
		return refuse(ClassificationBinding, "inventory-profile-mismatch")
	}
	seen := make(map[string]struct{}, len(observed.Roles))
	for _, candidate := range observed.Roles {
		if candidate.RoleID == "" {
			return refuse(ClassificationMalformed, "component-empty-role-id")
		}
		if _, exists := seen[candidate.RoleID]; exists {
			return refuse(ClassificationMalformed, "component-duplicate-role-id")
		}
		seen[candidate.RoleID] = struct{}{}
	}
	if len(observed.Roles) < len(profile.Roles) {
		return refuse(ClassificationBinding, "component-missing")
	}
	if len(observed.Roles) > len(profile.Roles) {
		return refuse(ClassificationBinding, "component-extra")
	}
	for index, expected := range profile.Roles {
		if !reflect.DeepEqual(observed.Roles[index], expected) {
			return refuse(ClassificationBinding, "component-mixed")
		}
	}
	return accept("component-inventory-exact")
}

func EvaluateActivation(profile Profile, observed InventoryObservation) Result {
	if result := EvaluateInventory(profile, observed); result.Decision != DecisionAccept {
		return result
	}
	// CanonicalProfile is permanently inactive in I0, so the signing refusal is
	// the only reachable outcome below. The bootstrap, service, entitlement, and
	// accept branches remain intentionally dormant until a future active-profile
	// variant adds its own tests before changing this trust-boundary contract.
	if profile.Signing.State != ProfileActive || profile.Signing.TeamIdentifier == "" ||
		profile.Signing.ProvisioningProfileSetID == "" || profile.Signing.CodeDirectoryHashSetID == "" ||
		profile.Signing.EntitlementsDigestSetID == "" {
		return refuse(ClassificationTrustState, "signing-profile-inactive")
	}
	if profile.Bootstrap.State != ProfileActive || profile.Bootstrap.BootstrapOwner == "" ||
		profile.Bootstrap.StoreFormatIdentity == unselectedProductStoreIdentity {
		return refuse(ClassificationTrustState, "bootstrap-profile-inactive")
	}
	for _, service := range profile.Services {
		if service.ActivationState != ProfileActive {
			return refuse(ClassificationTrustState, "service-profile-inactive")
		}
	}
	for _, entitlement := range profile.Entitlements {
		if entitlement.Disposition == EntitlementInactiveUnresolved {
			return refuse(ClassificationTrustState, "entitlement-profile-inactive")
		}
	}
	return accept("profile-active")
}

func EvaluateBootstrap(input BootstrapInput) BootstrapOutput {
	output := BootstrapOutput{RecordVersion: 0, State: input.State, AttemptsEnabled: false}
	if input.RecordVersion != 0 || input.ProfileID != ProfileIdentity {
		output.State = BootstrapRepairRequired
		output.Result = refuse(ClassificationMalformed, "bootstrap-input-profile")
		return output
	}
	if input.AttemptsEnabled && input.State != BootstrapReady {
		output.State = BootstrapRepairRequired
		output.Result = refuse(ClassificationRepair, "attempts-enabled-before-ready")
		return output
	}
	if input.InventoryDecision == DecisionRefuse {
		output.State = BootstrapRepairRequired
		output.Result = refuse(ClassificationRepair, "bootstrap-component-inventory-refused")
		return output
	}
	if input.SigningState == ProfileActive || input.BootstrapAuthorityState == ProfileActive {
		output.State = BootstrapRepairRequired
		output.Result = refuse(ClassificationTrustState, "i0-active-bootstrap-profile-unsupported")
		return output
	}

	next, ok := bootstrapTransitions[input.State][input.Event]
	if !ok {
		output.State = BootstrapRepairRequired
		output.Result = refuse(ClassificationRepair, "bootstrap-transition-invalid")
		return output
	}
	if input.Event == BootstrapRecordIdentity && input.BootstrapAuthorityState != ProfileActive {
		output.State = BootstrapRepairRequired
		output.Result = refuse(ClassificationTrustState, "bootstrap-authority-inactive")
		return output
	}
	if input.Event == BootstrapRegisterServices && input.SigningState != ProfileActive {
		output.State = BootstrapRepairRequired
		output.Result = refuse(ClassificationTrustState, "signing-profile-inactive")
		return output
	}
	if input.Event == BootstrapEnrollProtectedRoot &&
		(input.ProtectedRootState != "enrolled-exact" || input.OwnerLockState != "enrolled-exact" || input.StoreState != "open-without-create-exact") {
		output.State = BootstrapRepairRequired
		output.Result = refuse(ClassificationRepair, "protected-root-owner-store-not-exact")
		return output
	}
	if input.Event == BootstrapCommitEpoch && input.EpochState != "epoch-one-prepared-exact" {
		output.State = BootstrapRepairRequired
		output.Result = refuse(ClassificationRepair, "epoch-one-not-exact")
		return output
	}
	if input.Event == BootstrapCompleteRecovery && input.RecoveryState != "owner-store-recovery-clean" {
		output.State = BootstrapRepairRequired
		output.Result = refuse(ClassificationRepair, "owner-store-recovery-not-clean")
		return output
	}

	output.State = next
	output.AttemptsEnabled = next == BootstrapReady
	output.Result = accept("bootstrap-transition-exact")
	return output
}

var bootstrapTransitions = map[BootstrapState]map[BootstrapEvent]BootstrapState{
	BootstrapAbsent: {
		BootstrapBegin: BootstrapJournaledDisabled,
	},
	BootstrapJournaledDisabled: {
		BootstrapRecordIdentity: BootstrapIdentityPreparedDisabled,
	},
	BootstrapIdentityPreparedDisabled: {
		BootstrapRegisterServices: BootstrapServicesPreparedDisabled,
	},
	BootstrapServicesPreparedDisabled: {
		BootstrapEnrollProtectedRoot: BootstrapRootPreparedDisabled,
	},
	BootstrapRootPreparedDisabled: {
		BootstrapVerifyComponents: BootstrapComponentsVerifiedDisabled,
	},
	BootstrapComponentsVerifiedDisabled: {
		BootstrapCommitEpoch: BootstrapEpochCommittedRecovering,
	},
	BootstrapEpochCommittedRecovering: {
		BootstrapCompleteRecovery: BootstrapReady,
	},
}

func UpdateIdentityForProfile(profile Profile) UpdateIdentity {
	return UpdateIdentity{
		RecordVersion:           0,
		ProfileID:               profile.ProfileID,
		RoleSetDigest:           digestJSON(profile.Roles),
		ServiceSetDigest:        digestJSON(profile.Services),
		EntitlementsSetDigest:   digestJSON(profile.Entitlements),
		BootstrapContractID:     profile.Bootstrap.ContractID,
		StoreFormatIdentity:     profile.Bootstrap.StoreFormatIdentity,
		SupervisorIPCGeneration: 0,
		ValidatorGeneration:     1,
		SigningState:            profile.Signing.State,
		ReleaseIdentity:         "",
	}
}

func EvaluateUpdate(input UpdateInput) UpdateOutput {
	output := UpdateOutput{RecordVersion: 0, Compatibility: "refused", AttemptsState: "disabled"}
	if input.RecordVersion != 0 {
		output.Result = refuse(ClassificationMalformed, "update-input-version")
		return output
	}
	if input.ComponentInventoryResult.Decision != DecisionAccept {
		output.Result = refuse(ClassificationRepair, "update-component-inventory-refused")
		return output
	}
	if !input.CompleteBundleObserved || !input.BothValidatorTuplesExact || !input.ServicesStopped {
		output.Result = refuse(ClassificationRepair, "update-incomplete-or-mixed")
		return output
	}
	if input.Predecessor.ProfileID != input.Successor.ProfileID ||
		input.Predecessor.RoleSetDigest != input.Successor.RoleSetDigest ||
		input.Predecessor.ServiceSetDigest != input.Successor.ServiceSetDigest ||
		input.Predecessor.EntitlementsSetDigest != input.Successor.EntitlementsSetDigest ||
		input.Predecessor.BootstrapContractID != input.Successor.BootstrapContractID ||
		input.Predecessor.SupervisorIPCGeneration != input.Successor.SupervisorIPCGeneration ||
		input.Predecessor.ValidatorGeneration != input.Successor.ValidatorGeneration {
		output.Result = refuse(ClassificationRepair, "update-identity-mixed")
		return output
	}
	canonical := UpdateIdentityForProfile(CanonicalProfile())
	if !reflect.DeepEqual(input.Predecessor, canonical) || !reflect.DeepEqual(input.Successor, canonical) {
		output.Result = refuse(ClassificationTrustState, "i0-active-update-profile-unsupported")
		return output
	}
	if input.Predecessor.SigningState != ProfileActive || input.Successor.SigningState != ProfileActive ||
		input.Predecessor.ReleaseIdentity == "" || input.Successor.ReleaseIdentity == "" {
		output.Result = refuse(ClassificationTrustState, "update-signing-profile-inactive")
		return output
	}
	if !input.StateRootIdentityPreserved || !input.OwnerLockIdentityPreserved {
		output.Compatibility = "repair-required"
		output.Result = refuse(ClassificationRepair, "protected-root-or-owner-lock-not-preserved")
		return output
	}
	if input.Predecessor.StoreFormatIdentity != input.Successor.StoreFormatIdentity {
		if !input.StoreMigrationAuthorized {
			output.Compatibility = "migration-required"
			output.Result = refuse(ClassificationRepair, "store-migration-not-authorized")
			return output
		}
		output.Compatibility = "authorized-store-migration"
		output.Result = accept("manual-whole-bundle-replacement-with-migration")
		return output
	}
	output.Compatibility = "manual-whole-bundle-replacement-compatible"
	output.Result = accept("manual-whole-bundle-replacement-compatible")
	return output
}

func ClassifyRepair(input RepairInput) RepairOutput {
	output := RepairOutput{RecordVersion: 0, AttemptsState: "disabled"}
	if input.RecordVersion != 0 {
		output.Class = RepairRefuseAutomatic
		output.Result = refuse(ClassificationMalformed, "repair-input-version")
		return output
	}
	if !input.InstallationRootAvailable {
		output.Class = RepairNewInstallationRequired
		output.Result = refuse(ClassificationRepair, "installation-root-unavailable")
		return output
	}
	if !input.StoreAndHistoryComplete || !input.CleanupObligationsResolved {
		output.Class = RepairRefuseAutomatic
		output.Result = refuse(ClassificationRepair, "history-or-cleanup-unresolved")
		return output
	}
	if !input.StateRootContinuityProven || !input.OwnerLockContinuityProven {
		output.Class = RepairProtectedStateRequired
		output.Result = refuse(ClassificationRepair, "protected-state-continuity-unproven")
		return output
	}
	if !input.CurrentEpochExact {
		if input.ForwardTransitionAuthorized {
			output.Class = RepairForwardTrustTransition
			output.Result = accept("forward-trust-transition-authorized")
			return output
		}
		output.Class = RepairRefuseAutomatic
		output.Result = refuse(ClassificationRepair, "epoch-mismatch-without-forward-authority")
		return output
	}
	if !input.ApplicationFilesExact {
		output.Class = RepairRestoreApplicationFiles
		output.Result = accept("restore-current-release-preserve-state")
		return output
	}
	output.Class = RepairRestoreApplicationFiles
	output.Result = accept("application-and-state-already-exact")
	return output
}

func ClassifyUninstall(mode UninstallMode) UninstallOutput {
	output := UninstallOutput{RecordVersion: 0, Mode: mode, AttemptsState: "disabled"}
	if mode != UninstallApplicationOnly && mode != UninstallLocalWhereSafe && mode != UninstallAbandonIdentity {
		output.Result = refuse(ClassificationMalformed, "uninstall-mode")
		return output
	}

	remove := func(dataClass, reason string) UninstallDisposition {
		return UninstallDisposition{RecordVersion: 0, DataClass: dataClass, Disposition: DispositionRemove, Reason: reason}
	}
	retain := func(dataClass, reason string) UninstallDisposition {
		return UninstallDisposition{RecordVersion: 0, DataClass: dataClass, Disposition: DispositionRetain, Reason: reason}
	}
	deferRemoval := func(dataClass, reason string) UninstallDisposition {
		return UninstallDisposition{RecordVersion: 0, DataClass: dataClass, Disposition: DispositionDefer, Reason: reason}
	}

	output.Dispositions = []UninstallDisposition{
		remove("application-bundle", "selected-local-application-removal"),
		remove("service-registrations", "stop-and-unregister-per-user-services"),
		remove("non-authoritative-daemon-cache", "contains-no-installation-authority"),
		deferRemoval("source-validator-container-residue", "cleanup-must-first-prove-empty-inventory"),
		retain("supervisor-state-root", "authoritative-state-and-recovery-history"),
		retain("owner-lock-object", "installation-continuity-identity"),
		retain("grant-attempt-and-cleanup-history", "replay-and-reconciliation-required"),
		retain("archive-and-nonreuse-history", "referenced-history-deletion-not-authorized"),
		retain("quarantine-and-repair-state", "cannot-clear-by-uninstall"),
		{RecordVersion: 0, DataClass: "exported-receipts-external-witnesses-and-backups", Disposition: DispositionOutsideScope, Reason: "not-local-erasure-authority"},
	}

	switch mode {
	case UninstallApplicationOnly:
		output.Dispositions = append(output.Dispositions,
			retain("installation-and-operational-keys", "preserve-reinstall-and-recovery-path"),
			retain("broker-user-content", "application-only-removal-does-not-delete-content"),
		)
	case UninstallLocalWhereSafe:
		output.Dispositions = append(output.Dispositions,
			deferRemoval("installation-and-operational-keys", "remove-only-after-authority-and-cleanup-safety-check"),
			deferRemoval("broker-user-content", "remove-only-when-no-retention-or-release-obligation-remains"),
		)
	case UninstallAbandonIdentity:
		output.Dispositions = append(output.Dispositions,
			UninstallDisposition{RecordVersion: 0, DataClass: "installation-identity", Disposition: DispositionMarkAbandoned, Reason: "make-local-installation-unusable-with-erasure-limit-record"},
			deferRemoval("installation-and-operational-keys", "retire-only-after-abandonment-record-and-cleanup-safety-check"),
			deferRemoval("broker-user-content", "remove-only-when-no-retention-or-release-obligation-remains"),
		)
	}
	output.Result = accept("passive-uninstall-classification")
	return output
}

func digestJSON(value any) string {
	// #279 defers error propagation because it would alter UpdateIdentityForProfile's API; current inputs are package-controlled and marshalable.
	encoded, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:])
}
