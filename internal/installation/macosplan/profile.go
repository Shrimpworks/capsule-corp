package macosplan

import "capsule.local/capsule/internal/protocol/sourcevalidatorpassive"

const (
	brokerBundlePath     = ApplicationPath
	daemonBundlePath     = ApplicationPath + "/Contents/Library/Helpers/CapsuleDaemon.app"
	supervisorBundlePath = ApplicationPath + "/Contents/Library/Helpers/CapsuleSupervisor.app"

	brokerValidatorBundlePath = ApplicationPath + "/Contents/XPCServices/CapsuleSourceValidatorBroker.xpc"
	daemonValidatorBundlePath = daemonBundlePath + "/Contents/XPCServices/CapsuleSourceValidatorDaemon.xpc"
)

// CanonicalProfile returns a fresh copy of the one closed I0 no-guest profile.
// It is structurally complete but deliberately cannot activate.
func CanonicalProfile() Profile {
	roles := []RoleProjection{
		role(BrokerRoleID, "", PackagingVisibleApplication, VisibilityUserVisible,
			brokerBundlePath, brokerBundlePath+"/Contents/MacOS/Capsule", BrokerBundleIdentity, BrokerBundleIdentity, ""),
		role(DaemonRoleID, BrokerRoleID, PackagingEmbeddedAgent, VisibilityEmbedded,
			daemonBundlePath, daemonBundlePath+"/Contents/MacOS/CapsuleDaemon", DaemonBundleIdentity, DaemonBundleIdentity, DaemonBundleIdentity),
		role(SupervisorRoleID, BrokerRoleID, PackagingEmbeddedAgent, VisibilityEmbedded,
			supervisorBundlePath, supervisorBundlePath+"/Contents/MacOS/CapsuleSupervisor", SupervisorBundleIdentity, SupervisorBundleIdentity, SupervisorBundleIdentity),
		role(DaemonValidatorLauncherRoleID, DaemonRoleID, PackagingPrivateXPCService, VisibilityEmbedded,
			daemonValidatorBundlePath, daemonValidatorBundlePath+"/Contents/MacOS/CapsuleSourceValidatorDaemonLauncher",
			DaemonValidatorServiceIdentity, DaemonValidatorServiceIdentity, DaemonValidatorServiceIdentity),
		role(DaemonValidatorParserRoleID, DaemonValidatorLauncherRoleID, PackagingParserChild, VisibilityEmbedded,
			daemonValidatorBundlePath+"/Contents/Resources/capsule-mjs-source-validator-daemon",
			daemonValidatorBundlePath+"/Contents/Resources/capsule-mjs-source-validator-daemon",
			"", DaemonValidatorParserIdentity, ""),
		role(BrokerValidatorLauncherRoleID, BrokerRoleID, PackagingPrivateXPCService, VisibilityEmbedded,
			brokerValidatorBundlePath, brokerValidatorBundlePath+"/Contents/MacOS/CapsuleSourceValidatorBrokerLauncher",
			BrokerValidatorServiceIdentity, BrokerValidatorServiceIdentity, BrokerValidatorServiceIdentity),
		role(BrokerValidatorParserRoleID, BrokerValidatorLauncherRoleID, PackagingParserChild, VisibilityEmbedded,
			brokerValidatorBundlePath+"/Contents/Resources/capsule-mjs-source-validator-approval-broker",
			brokerValidatorBundlePath+"/Contents/Resources/capsule-mjs-source-validator-approval-broker",
			"", BrokerValidatorParserIdentity, ""),
	}
	roles[1].ServiceDescriptorPath = ApplicationPath + "/Contents/Library/LaunchAgents/com.capsulecorp.capsule.daemon.plist"
	roles[2].ServiceDescriptorPath = ApplicationPath + "/Contents/Library/LaunchAgents/com.capsulecorp.capsule.supervisor.plist"
	roles[3].ResourcePolicyPath = daemonValidatorBundlePath + "/Contents/Resources/resource-policy-inactive.bin"
	roles[5].ResourcePolicyPath = brokerValidatorBundlePath + "/Contents/Resources/resource-policy-inactive.bin"

	return Profile{
		RecordVersion:   0,
		ProfileID:       ProfileIdentity,
		ProductName:     ProductName,
		ApplicationPath: ApplicationPath,
		VisibleRoleID:   BrokerRoleID,
		Roles:           roles,
		Entitlements:    canonicalEntitlements(),
		Services:        canonicalServices(),
		Signing: SigningProjection{
			RecordVersion:            0,
			State:                    ProfileInactiveRefusing,
			TeamIdentifier:           "",
			DistributionChannel:      "unresolved",
			ProvisioningProfileSetID: "",
			CodeDirectoryHashSetID:   "",
			EntitlementsDigestSetID:  "",
			Blocker:                  "matching-team-profiles-and-signed-installed-corpus-unavailable",
		},
		Bootstrap: BootstrapContract{
			RecordVersion:            0,
			ContractID:               "capsule.macos-installation.bootstrap/i0",
			State:                    ProfileInactiveRefusing,
			BootstrapOwner:           "unresolved-containing-application-or-supervisor-refusing",
			StateRootClass:           "supervisor-private-app-sandbox-container",
			OwnerLockEntryName:       "supervisor.owner",
			OwnerLockMode:            0o600,
			OwnerLockLinkCount:       1,
			StoreEntryName:           "supervisor.store",
			StoreFormatIdentity:      "capsule.supervisor.product-store/unselected",
			OwnerMechanismIdentity:   "capsule.supervisor.owner.darwin-openat-flock/v0",
			OrdinaryStartupOperation: "open-without-create",
			InitialAttemptState:      "disabled",
			LocalEvidenceState:       "g2-passed-current-v1-no-guest-local-mechanic",
			InstalledEvidenceState:   "g3-blocked-team-profile-bootstrap-and-descriptor-relative-store-open",
		},
		ExcludedRoles: []string{
			"capsule.role.runner",
			"capsule.role.update-verifier",
			"capsule.role.trust-bootstrap-coordinator",
			"capsule.role.bundle-replacer",
			"capsule.role.source-preparer",
			"capsule.role.runtime",
			"capsule.role.backend",
			"capsule.role.guest",
		},
	}
}

func role(
	roleID string,
	containingRoleID string,
	packaging PackagingClass,
	visibility Visibility,
	bundlePath string,
	executablePath string,
	bundleIdentifier string,
	signingIdentifier string,
	serviceIdentifier string,
) RoleProjection {
	return RoleProjection{
		RecordVersion:     0,
		RoleID:            roleID,
		ContainingRoleID:  containingRoleID,
		Packaging:         packaging,
		Visibility:        visibility,
		BundlePath:        bundlePath,
		ExecutablePath:    executablePath,
		BundleIdentifier:  bundleIdentifier,
		SigningIdentifier: signingIdentifier,
		ServiceIdentifier: serviceIdentifier,
		Required:          true,
	}
}

func canonicalServices() []ServiceProjection {
	return []ServiceProjection{
		service("capsule.service.daemon-agent/i0", DaemonBundleIdentity, ServicePerUserAgent, DaemonRoleID, BrokerRoleID,
			"none", "", "containing-application-registration-only"),
		service("capsule.service.supervisor-agent/i0", SupervisorBundleIdentity, ServicePerUserAgent, SupervisorRoleID, BrokerRoleID,
			"none", "", "containing-application-registration-only"),
		service("capsule.service.supervisor.daemon/v0", SupervisorDaemonServiceIdentity, ServiceSupervisorListener, SupervisorRoleID, DaemonRoleID,
			"capsule.execution-supervisor.local.v0", "RegisterPlanV0,RequestAttemptV0", "pairwise-transport-unresolved-refusing"),
		service("capsule.service.supervisor.broker/v0", SupervisorBrokerServiceIdentity, ServiceSupervisorListener, SupervisorRoleID, BrokerRoleID,
			"capsule.execution-supervisor.local.v0", "GetRegisteredPlanV0,SubmitApprovalV0", "pairwise-transport-unresolved-refusing"),
		service("capsule.service.source-validator.daemon/v1", sourcevalidatorpassive.V1DaemonServiceIdentity, ServicePrivateXPC,
			DaemonValidatorLauncherRoleID, DaemonRoleID, sourcevalidatorpassive.V1ProtocolIdentity,
			sourcevalidatorpassive.V1DaemonMethodIdentity, "role-private-reachability-unproven-refusing"),
		service("capsule.service.source-validator.approval-broker/v1", sourcevalidatorpassive.V1BrokerServiceIdentity, ServicePrivateXPC,
			BrokerValidatorLauncherRoleID, BrokerRoleID, sourcevalidatorpassive.V1ProtocolIdentity,
			sourcevalidatorpassive.V1BrokerMethodIdentity, "role-private-reachability-unproven-refusing"),
	}
}

func service(
	serviceID string,
	serviceName string,
	class ServiceClass,
	ownerRoleID string,
	peerRoleID string,
	protocolIdentity string,
	methodIdentities string,
	reachability string,
) ServiceProjection {
	return ServiceProjection{
		RecordVersion:    0,
		ServiceID:        serviceID,
		ServiceName:      serviceName,
		ServiceClass:     class,
		OwnerRoleID:      ownerRoleID,
		PeerRoleID:       peerRoleID,
		ProtocolIdentity: protocolIdentity,
		MethodIdentities: methodIdentities,
		Reachability:     reachability,
		ActivationState:  ProfileInactiveRefusing,
	}
}

func canonicalEntitlements() []EntitlementProjection {
	var values []EntitlementProjection
	for _, roleID := range []string{
		BrokerRoleID,
		DaemonRoleID,
		SupervisorRoleID,
		DaemonValidatorLauncherRoleID,
		DaemonValidatorParserRoleID,
		BrokerValidatorLauncherRoleID,
		BrokerValidatorParserRoleID,
	} {
		values = append(values,
			entitlement(roleID, "com.apple.security.app-sandbox", EntitlementRequiredTrue, "true", "app-sandbox-baseline"),
			entitlement(roleID, "com.apple.security.get-task-allow", EntitlementRequiredAbsent, "absent", "hardened-runtime-baseline"),
			entitlement(roleID, "com.apple.security.cs.disable-library-validation", EntitlementRequiredAbsent, "absent", "closed-native-loading"),
			entitlement(roleID, "com.apple.security.cs.allow-jit", EntitlementRequiredAbsent, "absent", "closed-native-loading"),
			entitlement(roleID, "com.apple.security.cs.allow-unsigned-executable-memory", EntitlementRequiredAbsent, "absent", "closed-native-loading"),
			entitlement(roleID, "com.apple.security.hypervisor", EntitlementRequiredAbsent, "absent", "no-guest-profile"),
		)
	}

	values = append(values,
		entitlement(BrokerRoleID, "com.apple.security.files.user-selected.read-only", EntitlementRequiredAbsent, "absent", "no-guest-inline-json-profile"),
		entitlement(BrokerRoleID, "keychain-access-groups", EntitlementInactiveUnresolved, "epoch-scoped-broker-group-unset", "signing-dependent"),
		entitlement(BrokerRoleID, "com.apple.security.application-groups", EntitlementInactiveUnresolved, "pairwise-broker-supervisor-group-or-private-service-unselected", "ipc-topology-dependent"),

		entitlement(DaemonRoleID, "keychain-access-groups", EntitlementRequiredAbsent, "absent", "daemon-key-exclusion"),
		entitlement(DaemonRoleID, "com.apple.security.files.user-selected.read-only", EntitlementRequiredAbsent, "absent", "broker-only-file-selection"),
		entitlement(DaemonRoleID, "com.apple.security.application-groups", EntitlementInactiveUnresolved, "pairwise-daemon-supervisor-group-or-private-service-unselected", "ipc-topology-dependent"),

		entitlement(SupervisorRoleID, "keychain-access-groups", EntitlementInactiveUnresolved, "epoch-scoped-supervisor-evidence-group-unset", "signing-dependent"),
		entitlement(SupervisorRoleID, "com.apple.security.files.user-selected.read-only", EntitlementRequiredAbsent, "absent", "broker-only-file-selection"),
		entitlement(SupervisorRoleID, "com.apple.security.application-groups", EntitlementInactiveUnresolved, "two-pairwise-groups-or-private-services-unselected", "ipc-topology-dependent"),
	)

	for _, roleID := range []string{
		DaemonValidatorLauncherRoleID,
		DaemonValidatorParserRoleID,
		BrokerValidatorLauncherRoleID,
		BrokerValidatorParserRoleID,
	} {
		values = append(values,
			entitlement(roleID, "keychain-access-groups", EntitlementRequiredAbsent, "absent", "validator-no-key-authority"),
			entitlement(roleID, "com.apple.security.application-groups", EntitlementRequiredAbsent, "absent", "validator-role-separation"),
			entitlement(roleID, "com.apple.security.files.user-selected.read-only", EntitlementRequiredAbsent, "absent", "validator-copied-source-only"),
			entitlement(roleID, "com.apple.security.network.client", EntitlementRequiredAbsent, "absent", "validator-no-network"),
			entitlement(roleID, "com.apple.security.network.server", EntitlementRequiredAbsent, "absent", "validator-no-network"),
		)
	}

	for _, roleID := range []string{DaemonValidatorParserRoleID, BrokerValidatorParserRoleID} {
		values = append(values,
			entitlement(roleID, "com.apple.security.inherit", EntitlementInactiveUnresolved, "signed-parser-child-shape-unset", "r3-r4-signing-and-confinement-dependent"),
		)
	}
	return values
}

func entitlement(
	roleID string,
	key string,
	disposition EntitlementDisposition,
	valueIdentity string,
	source string,
) EntitlementProjection {
	return EntitlementProjection{
		RecordVersion: 0,
		RoleID:        roleID,
		Key:           key,
		Disposition:   disposition,
		ValueIdentity: valueIdentity,
		Source:        source,
	}
}
