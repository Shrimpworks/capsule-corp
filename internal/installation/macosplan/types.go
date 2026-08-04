package macosplan

import (
	"fmt"

	"capsule.local/capsule/internal/protocol/sourcevalidatorpassive"
)

const (
	ProfileIdentity = "capsule.macos-installation.no-guest/i0"
	ProductName     = "Capsule"
	ApplicationPath = "Capsule.app"

	BrokerRoleID                  = "capsule.role.approval-broker/v0"
	DaemonRoleID                  = "capsule.role.agent-daemon/v0"
	SupervisorRoleID              = "capsule.role.execution-supervisor/v0"
	DaemonValidatorLauncherRoleID = "capsule.role.source-validator-launcher.daemon/v1"
	BrokerValidatorLauncherRoleID = "capsule.role.source-validator-launcher.approval-broker/v1"
	DaemonValidatorParserRoleID   = "capsule.role.source-validator-parser.daemon/v1"
	BrokerValidatorParserRoleID   = "capsule.role.source-validator-parser.approval-broker/v1"

	BrokerBundleIdentity     = sourcevalidatorpassive.V1BrokerContainingBundleIdentity
	DaemonBundleIdentity     = sourcevalidatorpassive.V1DaemonContainingBundleIdentity
	SupervisorBundleIdentity = "com.capsulecorp.capsule.supervisor"

	SupervisorDaemonServiceIdentity = "com.capsulecorp.capsule.supervisor.daemon.v0"
	SupervisorBrokerServiceIdentity = "com.capsulecorp.capsule.supervisor.broker.v0"

	DaemonValidatorServiceIdentity = sourcevalidatorpassive.V1DaemonServiceIdentity
	BrokerValidatorServiceIdentity = sourcevalidatorpassive.V1BrokerServiceIdentity
	DaemonValidatorParserIdentity  = sourcevalidatorpassive.V1DaemonParserIdentity
	BrokerValidatorParserIdentity  = sourcevalidatorpassive.V1BrokerParserIdentity
)

type Decision string

const (
	DecisionAccept Decision = "accept"
	DecisionRefuse Decision = "refuse"
)

type Classification string

const (
	ClassificationNone       Classification = "none"
	ClassificationMalformed  Classification = "malformed"
	ClassificationBinding    Classification = "binding"
	ClassificationTrustState Classification = "trust-state"
	ClassificationRepair     Classification = "repair-required"
)

type Result struct {
	Decision       Decision       `json:"decision"`
	Classification Classification `json:"classification"`
	Reason         string         `json:"reason"`
}

func accept(reason string) Result {
	return Result{Decision: DecisionAccept, Classification: ClassificationNone, Reason: reason}
}

func refuse(classification Classification, reason string) Result {
	return Result{Decision: DecisionRefuse, Classification: classification, Reason: reason}
}

type ProfileState string

const (
	ProfileInactiveRefusing ProfileState = "inactive-refusing"
	ProfileActive           ProfileState = "active"
)

type PackagingClass string

const (
	PackagingVisibleApplication PackagingClass = "visible-application"
	PackagingEmbeddedAgent      PackagingClass = "embedded-per-user-agent"
	PackagingPrivateXPCService  PackagingClass = "private-xpc-service"
	PackagingParserChild        PackagingClass = "embedded-parser-child"
)

type Visibility string

const (
	VisibilityUserVisible Visibility = "user-visible"
	VisibilityEmbedded    Visibility = "embedded-only"
)

// field-authority-object: capsule.macos-installation-role-projection v0
type RoleProjection struct {
	RecordVersion         uint16         `json:"recordVersion"`
	RoleID                string         `json:"roleId"`
	ContainingRoleID      string         `json:"containingRoleId"`
	Packaging             PackagingClass `json:"packaging"`
	Visibility            Visibility     `json:"visibility"`
	BundlePath            string         `json:"bundlePath"`
	ExecutablePath        string         `json:"executablePath"`
	ServiceDescriptorPath string         `json:"serviceDescriptorPath"`
	ResourcePolicyPath    string         `json:"resourcePolicyPath"`
	BundleIdentifier      string         `json:"bundleIdentifier"`
	SigningIdentifier     string         `json:"signingIdentifier"`
	ServiceIdentifier     string         `json:"serviceIdentifier"`
	Required              bool           `json:"required"`
}

type EntitlementDisposition string

const (
	EntitlementRequiredTrue       EntitlementDisposition = "required-true"
	EntitlementRequiredAbsent     EntitlementDisposition = "required-absent"
	EntitlementInactiveUnresolved EntitlementDisposition = "inactive-unresolved"
)

// field-authority-object: capsule.macos-installation-entitlement-projection v0
type EntitlementProjection struct {
	RecordVersion uint16                 `json:"recordVersion"`
	RoleID        string                 `json:"roleId"`
	Key           string                 `json:"key"`
	Disposition   EntitlementDisposition `json:"disposition"`
	ValueIdentity string                 `json:"valueIdentity"`
	Source        string                 `json:"source"`
}

type ServiceClass string

const (
	ServicePerUserAgent       ServiceClass = "smappservice-agent-candidate"
	ServiceSupervisorListener ServiceClass = "supervisor-mach-listener"
	ServicePrivateXPC         ServiceClass = "role-private-xpc"
)

// field-authority-object: capsule.macos-installation-service-projection v0
type ServiceProjection struct {
	RecordVersion    uint16       `json:"recordVersion"`
	ServiceID        string       `json:"serviceId"`
	ServiceName      string       `json:"serviceName"`
	ServiceClass     ServiceClass `json:"serviceClass"`
	OwnerRoleID      string       `json:"ownerRoleId"`
	PeerRoleID       string       `json:"peerRoleId"`
	ProtocolIdentity string       `json:"protocolIdentity"`
	MethodIdentities string       `json:"methodIdentities"`
	Reachability     string       `json:"reachability"`
	ActivationState  ProfileState `json:"activationState"`
}

// field-authority-object: capsule.macos-installation-signing-projection v0
type SigningProjection struct {
	RecordVersion            uint16       `json:"recordVersion"`
	State                    ProfileState `json:"state"`
	TeamIdentifier           string       `json:"teamIdentifier"`
	DistributionChannel      string       `json:"distributionChannel"`
	ProvisioningProfileSetID string       `json:"provisioningProfileSetId"`
	CodeDirectoryHashSetID   string       `json:"codeDirectoryHashSetId"`
	EntitlementsDigestSetID  string       `json:"entitlementsDigestSetId"`
	Blocker                  string       `json:"blocker"`
}

// field-authority-object: capsule.macos-installation-bootstrap-contract v0
type BootstrapContract struct {
	RecordVersion            uint16       `json:"recordVersion"`
	ContractID               string       `json:"contractId"`
	State                    ProfileState `json:"state"`
	BootstrapOwner           string       `json:"bootstrapOwner"`
	StateRootClass           string       `json:"stateRootClass"`
	OwnerLockEntryName       string       `json:"ownerLockEntryName"`
	OwnerLockMode            uint32       `json:"ownerLockMode"`
	OwnerLockLinkCount       uint64       `json:"ownerLockLinkCount"`
	StoreEntryName           string       `json:"storeEntryName"`
	StoreFormatIdentity      string       `json:"storeFormatIdentity"`
	OwnerMechanismIdentity   string       `json:"ownerMechanismIdentity"`
	OrdinaryStartupOperation string       `json:"ordinaryStartupOperation"`
	InitialAttemptState      string       `json:"initialAttemptState"`
	LocalEvidenceState       string       `json:"localEvidenceState"`
	InstalledEvidenceState   string       `json:"installedEvidenceState"`
}

// field-authority-object: capsule.macos-installation-profile v0
type Profile struct {
	RecordVersion   uint16                  `json:"recordVersion"`
	ProfileID       string                  `json:"profileId"`
	ProductName     string                  `json:"productName"`
	ApplicationPath string                  `json:"applicationPath"`
	VisibleRoleID   string                  `json:"visibleRoleId"`
	Roles           []RoleProjection        `json:"roles"`
	Entitlements    []EntitlementProjection `json:"entitlements"`
	Services        []ServiceProjection     `json:"services"`
	Signing         SigningProjection       `json:"signing"`
	Bootstrap       BootstrapContract       `json:"bootstrap"`
	ExcludedRoles   []string                `json:"excludedRoles"`
}

type InventoryObservation struct {
	ProfileID string           `json:"profileId"`
	Roles     []RoleProjection `json:"roles"`
}

type BootstrapState string

const (
	BootstrapAbsent                     BootstrapState = "absent"
	BootstrapJournaledDisabled          BootstrapState = "journaled-disabled"
	BootstrapIdentityPreparedDisabled   BootstrapState = "identity-prepared-disabled"
	BootstrapServicesPreparedDisabled   BootstrapState = "services-prepared-disabled"
	BootstrapRootPreparedDisabled       BootstrapState = "protected-root-prepared-disabled"
	BootstrapComponentsVerifiedDisabled BootstrapState = "components-verified-disabled"
	BootstrapEpochCommittedRecovering   BootstrapState = "epoch-committed-recovering"
	BootstrapReady                      BootstrapState = "ready"
	BootstrapRepairRequired             BootstrapState = "repair-required"
)

type BootstrapEvent string

const (
	BootstrapBegin               BootstrapEvent = "begin"
	BootstrapRecordIdentity      BootstrapEvent = "record-installation-identity"
	BootstrapRegisterServices    BootstrapEvent = "register-services"
	BootstrapEnrollProtectedRoot BootstrapEvent = "enroll-protected-root"
	BootstrapVerifyComponents    BootstrapEvent = "verify-components"
	BootstrapCommitEpoch         BootstrapEvent = "commit-epoch"
	BootstrapCompleteRecovery    BootstrapEvent = "complete-owner-store-recovery"
)

// field-authority-object: capsule.macos-installation-bootstrap-input v0
type BootstrapInput struct {
	RecordVersion           uint16         `json:"recordVersion"`
	ProfileID               string         `json:"profileId"`
	State                   BootstrapState `json:"state"`
	Event                   BootstrapEvent `json:"event"`
	InventoryDecision       Decision       `json:"inventoryDecision"`
	SigningState            ProfileState   `json:"signingState"`
	BootstrapAuthorityState ProfileState   `json:"bootstrapAuthorityState"`
	ProtectedRootState      string         `json:"protectedRootState"`
	OwnerLockState          string         `json:"ownerLockState"`
	StoreState              string         `json:"storeState"`
	EpochState              string         `json:"epochState"`
	RecoveryState           string         `json:"recoveryState"`
	AttemptsEnabled         bool           `json:"attemptsEnabled"`
}

// field-authority-object: capsule.macos-installation-bootstrap-output v0
type BootstrapOutput struct {
	RecordVersion   uint16         `json:"recordVersion"`
	State           BootstrapState `json:"state"`
	AttemptsEnabled bool           `json:"attemptsEnabled"`
	Result          Result         `json:"result"`
}

// field-authority-object: capsule.macos-installation-update-identity v0
type UpdateIdentity struct {
	RecordVersion           uint16       `json:"recordVersion"`
	ProfileID               string       `json:"profileId"`
	RoleSetDigest           string       `json:"roleSetDigest"`
	ServiceSetDigest        string       `json:"serviceSetDigest"`
	EntitlementsSetDigest   string       `json:"entitlementsSetDigest"`
	BootstrapContractID     string       `json:"bootstrapContractId"`
	StoreFormatIdentity     string       `json:"storeFormatIdentity"`
	SupervisorIPCGeneration uint16       `json:"supervisorIpcGeneration"`
	ValidatorGeneration     uint16       `json:"validatorGeneration"`
	SigningState            ProfileState `json:"signingState"`
	ReleaseIdentity         string       `json:"releaseIdentity"`
}

// field-authority-object: capsule.macos-installation-update-input v0
type UpdateInput struct {
	RecordVersion              uint16         `json:"recordVersion"`
	Predecessor                UpdateIdentity `json:"predecessor"`
	Successor                  UpdateIdentity `json:"successor"`
	CompleteBundleObserved     bool           `json:"completeBundleObserved"`
	StateRootIdentityPreserved bool           `json:"stateRootIdentityPreserved"`
	OwnerLockIdentityPreserved bool           `json:"ownerLockIdentityPreserved"`
	StoreMigrationAuthorized   bool           `json:"storeMigrationAuthorized"`
	BothValidatorTuplesExact   bool           `json:"bothValidatorTuplesExact"`
	ServicesStopped            bool           `json:"servicesStopped"`
	ComponentInventoryResult   Result         `json:"componentInventoryResult"`
}

// field-authority-object: capsule.macos-installation-update-output v0
type UpdateOutput struct {
	RecordVersion uint16 `json:"recordVersion"`
	Compatibility string `json:"compatibility"`
	AttemptsState string `json:"attemptsState"`
	Result        Result `json:"result"`
}

type RepairClass string

const (
	RepairRestoreApplicationFiles RepairClass = "restore-current-release-preserve-state"
	RepairForwardTrustTransition  RepairClass = "authorized-forward-trust-transition"
	RepairProtectedStateRequired  RepairClass = "protected-state-repair-required"
	RepairNewInstallationRequired RepairClass = "new-installation-required"
	RepairRefuseAutomatic         RepairClass = "refuse-automatic-repair"
)

// field-authority-object: capsule.macos-installation-repair-input v0
type RepairInput struct {
	RecordVersion               uint16 `json:"recordVersion"`
	ApplicationFilesExact       bool   `json:"applicationFilesExact"`
	InstallationRootAvailable   bool   `json:"installationRootAvailable"`
	StateRootContinuityProven   bool   `json:"stateRootContinuityProven"`
	OwnerLockContinuityProven   bool   `json:"ownerLockContinuityProven"`
	StoreAndHistoryComplete     bool   `json:"storeAndHistoryComplete"`
	CleanupObligationsResolved  bool   `json:"cleanupObligationsResolved"`
	CurrentEpochExact           bool   `json:"currentEpochExact"`
	ForwardTransitionAuthorized bool   `json:"forwardTransitionAuthorized"`
}

// field-authority-object: capsule.macos-installation-repair-output v0
type RepairOutput struct {
	RecordVersion uint16      `json:"recordVersion"`
	Class         RepairClass `json:"class"`
	AttemptsState string      `json:"attemptsState"`
	Result        Result      `json:"result"`
}

type UninstallMode string

const (
	UninstallApplicationOnly UninstallMode = "remove-application-preserve-state"
	UninstallLocalWhereSafe  UninstallMode = "remove-local-where-safe"
	UninstallAbandonIdentity UninstallMode = "abandon-installation-identity"
)

type DataDisposition string

const (
	DispositionRemove        DataDisposition = "remove"
	DispositionRetain        DataDisposition = "retain"
	DispositionDefer         DataDisposition = "defer"
	DispositionMarkAbandoned DataDisposition = "mark-abandoned"
	DispositionOutsideScope  DataDisposition = "outside-local-scope"
)

// field-authority-object: capsule.macos-installation-uninstall-disposition v0
type UninstallDisposition struct {
	RecordVersion uint16          `json:"recordVersion"`
	DataClass     string          `json:"dataClass"`
	Disposition   DataDisposition `json:"disposition"`
	Reason        string          `json:"reason"`
}

// field-authority-object: capsule.macos-installation-uninstall-output v0
type UninstallOutput struct {
	RecordVersion uint16                 `json:"recordVersion"`
	Mode          UninstallMode          `json:"mode"`
	AttemptsState string                 `json:"attemptsState"`
	Dispositions  []UninstallDisposition `json:"dispositions"`
	Result        Result                 `json:"result"`
}

func (result Result) Validate() error {
	if result.Decision == DecisionAccept && result.Classification == ClassificationNone && result.Reason != "" {
		return nil
	}
	if result.Decision == DecisionRefuse && result.Classification != ClassificationNone && result.Reason != "" {
		return nil
	}
	return fmt.Errorf("invalid result tuple")
}
