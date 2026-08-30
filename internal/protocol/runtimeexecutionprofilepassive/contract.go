package runtimeexecutionprofilepassive

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"capsule.local/capsule/internal/protocol/strictjson"
)

const (
	ContractIdentity = "capsule.governed-deno-core.execution-profile-preparation/c2a"
	ContractVersion  = 1
	ContractStatus   = "passive-preparation-refusing"

	ContractBytes  = 26850
	ContractSHA256 = "d4ce88888186266f5d251e6246c889b1fd46d7746bb0ba56bcc4b3ce4675992f"
	C1Bytes        = 9289
	C1SHA256       = "d5d75e638a15be6c9f4a3230d17309d085f6ec103a73b64d9e0fd656a5423c9e"
)

var expectedArtifactRoles = []string{
	"host-runner", "governed-libkrun-dylib", "libkrunfw-dylib", "firmware",
	"guest-kernel", "trusted-init", "trusted-launcher", "raw-runtime-root",
	"governed-runtime-bundle-manifest",
}

var expectedRunnerRoles = []string{
	"runner-stdin-null", "runner-stdout", "runner-stderr", "record-before-start-control",
	"runtime-root", "registered-source-port-input", "approved-inline-input-port-input",
	"completion-result-port-output",
}

var expectedRunnerModes = []string{
	"O_RDONLY", "O_WRONLY", "O_WRONLY", "O_RDONLY", "O_RDONLY", "O_RDONLY", "O_RDONLY", "O_WRONLY",
}

var expectedGuestRoles = []string{
	"launcher-stdin-null", "launcher-stdout-null", "launcher-stderr-null",
	"registered-source", "approved-inline-input", "completion-result",
}

var expectedGuestModes = []string{"O_RDONLY", "O_WRONLY", "O_WRONLY", "O_RDONLY", "O_RDONLY", "O_WRONLY"}

var expectedMatrixGroups = map[string]int{
	"preflight-and-identities":       5,
	"runner-descriptors":             10,
	"root-custody-and-device":        8,
	"runtime-surface":                7,
	"no-loader-v6":                   8,
	"file-open-and-syscall-traces":   8,
	"real-p0-3-transport":            12,
	"launcher-and-child-tree":        8,
	"death-cancel-and-teardown":      12,
	"device-network-and-restoration": 6,
	"evidence-and-classification":    7,
}

type Contract struct {
	Contract                 string                   `json:"contract"`
	Version                  int                      `json:"version"`
	Status                   string                   `json:"status"`
	Scope                    json.RawMessage          `json:"scope"`
	C1Binding                C1Binding                `json:"c1Binding"`
	C1PlanAndProfileBindings C1PlanAndProfileBindings `json:"c1PlanAndProfileBindings"`
	EvidencePins             json.RawMessage          `json:"evidencePins"`
	GovernedSource           json.RawMessage          `json:"governedSource"`
	ArtifactClosure          ArtifactClosure          `json:"artifactClosure"`
	KnownAnswer              KnownAnswer              `json:"knownAnswer"`
	RunnerDescriptorProfile  RunnerDescriptorProfile  `json:"runnerDescriptorProfile"`
	GuestDescriptorProfile   GuestDescriptorProfile   `json:"guestDescriptorProfile"`
	MachineProfile           MachineProfile           `json:"machineProfile"`
	TransportProfile         TransportProfile         `json:"transportProfile"`
	TeardownProfile          TeardownProfile          `json:"teardownProfile"`
	C2BMatrix                []MatrixGroup            `json:"c2bMatrix"`
	RestorationMutations     []string                 `json:"restorationMutations"`
	EvidenceRetention        []string                 `json:"evidenceRetention"`
	Blockers                 []string                 `json:"blockers"`
	RefusalCodes             []string                 `json:"refusalCodes"`
	StopConditions           []string                 `json:"stopConditions"`
	WorkStatus               WorkStatus               `json:"workStatus"`
	Effects                  Effects                  `json:"effects"`
}

type C1Binding struct {
	Contract    string `json:"contract"`
	Fixture     string `json:"fixture"`
	Bytes       int    `json:"bytes"`
	SHA256      string `json:"sha256"`
	Consumption string `json:"consumption"`
	Mutation    string `json:"mutation"`
}

type C1PlanAndProfileBindings struct {
	ExecutionPlan  ExecutionPlanBinding  `json:"executionPlan"`
	RuntimeProfile RuntimeProfileBinding `json:"runtimeProfile"`
}

type ExecutionPlanBinding struct {
	ObjectType             string   `json:"objectType"`
	ObjectVersion          int      `json:"objectVersion"`
	MediaType              string   `json:"mediaType"`
	SourceProfile          string   `json:"sourceProfile"`
	SourceBinding          string   `json:"sourceBinding"`
	InputEncoding          string   `json:"inputEncoding"`
	InputBinding           string   `json:"inputBinding"`
	OutputLimitBinding     string   `json:"outputLimitBinding"`
	RequiredReferenceRoles []string `json:"requiredReferenceRoles"`
}

type RuntimeProfileBinding struct {
	SelectedIdentity         *string `json:"selectedIdentity"`
	SelectedDigest           *string `json:"selectedDigest"`
	MachineProfileIdentity   string  `json:"machineProfileIdentity"`
	TransportProfileIdentity string  `json:"transportProfileIdentity"`
	State                    string  `json:"state"`
}

type ArtifactClosure struct {
	Activation                    string     `json:"activation"`
	Required                      []Artifact `json:"required"`
	HistoricalEvidenceNotReusable []string   `json:"historicalEvidenceNotReusable"`
}

type Artifact struct {
	Role         string  `json:"role"`
	SelectedPath *string `json:"selectedPath"`
	SHA256       *string `json:"sha256"`
	State        string  `json:"state"`
	Admitted     bool    `json:"admitted"`
}

type KnownAnswer struct {
	Source             SourceKnownAnswer     `json:"source"`
	CanonicalInput     BytesKnownAnswer      `json:"canonicalInput"`
	ExpectedCompletion CompletionKnownAnswer `json:"expectedCompletion"`
}

type SourceKnownAnswer struct {
	Path                 string `json:"path"`
	UTF8                 string `json:"utf8"`
	Bytes                int    `json:"bytes"`
	SHA256               string `json:"sha256"`
	SourceManifestBytes  int    `json:"sourceManifestBytes"`
	SourceManifestSHA256 string `json:"sourceManifestSha256"`
}

type BytesKnownAnswer struct {
	UTF8   string `json:"utf8"`
	Bytes  int    `json:"bytes"`
	SHA256 string `json:"sha256"`
}

type CompletionKnownAnswer struct {
	TerminalStatus  string `json:"terminalStatus"`
	JSONUTF8        string `json:"jsonUtf8"`
	JSONBytes       int    `json:"jsonBytes"`
	JSONSHA256      string `json:"jsonSha256"`
	Commit          string `json:"commit"`
	RunnerExitAlone string `json:"runnerExitAlone"`
}

type RunnerDescriptorProfile struct {
	Identity                  string            `json:"identity"`
	NumericRange              NumericRange      `json:"numericRange"`
	Entries                   []DescriptorEntry `json:"entries"`
	PortCalls                 []PortCall        `json:"portCalls"`
	InheritanceSequence       []string          `json:"inheritanceSequence"`
	ForbiddenExtras           []string          `json:"forbiddenExtras"`
	RuntimeCreatedDescriptors string            `json:"runtimeCreatedDescriptors"`
}

type NumericRange struct {
	Minimum            int `json:"minimum"`
	Maximum            int `json:"maximum"`
	CloseFromInclusive int `json:"closeFromInclusive"`
}

type DescriptorEntry struct {
	FD          int    `json:"fd"`
	Role        string `json:"role"`
	AccessMode  string `json:"accessMode"`
	Producer    string `json:"producer"`
	Consumer    string `json:"consumer"`
	Direction   string `json:"direction"`
	Inheritance string `json:"inheritance"`
	Cloexec     string `json:"cloexec"`
	Lifetime    string `json:"lifetime"`
}

type PortCall struct {
	PortID    int    `json:"portId"`
	Name      string `json:"name"`
	InputFD   int    `json:"inputFd"`
	OutputFD  int    `json:"outputFd"`
	GuestNode string `json:"guestNode"`
}

type GuestDescriptorProfile struct {
	Identity                 string                 `json:"identity"`
	PreChild                 []GuestDescriptorEntry `json:"preChild"`
	CloseFromInclusive       int                    `json:"closeFromInclusive"`
	WorkloadCompletionAccess string                 `json:"workloadCompletionAccess"`
	ChildManifest            json.RawMessage        `json:"childManifest"`
	State                    string                 `json:"state"`
}

type GuestDescriptorEntry struct {
	FD         int    `json:"fd"`
	Role       string `json:"role"`
	AccessMode string `json:"accessMode"`
}

type MachineProfile struct {
	Identity                string         `json:"identity"`
	Architecture            string         `json:"architecture"`
	Backend                 string         `json:"backend"`
	VCPUs                   int            `json:"vcpus"`
	GuestRAMMiB             int            `json:"guestRamMiB"`
	WallTimeMs              int            `json:"wallTimeMs"`
	CPUTimeLimitMs          *int           `json:"cpuTimeLimitMs"`
	CPUTimePolicy           string         `json:"cpuTimePolicy"`
	HostVMMMemoryLimitBytes *int64         `json:"hostVmmMemoryLimitBytes"`
	HostVMMMemoryPolicy     string         `json:"hostVmmMemoryPolicy"`
	Concurrency             int            `json:"concurrency"`
	Scratch                 ScratchProfile `json:"scratch"`
	UnknownOrDifferentValue string         `json:"unknownOrDifferentValue"`
}

type ScratchProfile struct {
	Mode         string `json:"mode"`
	MaximumBytes *int64 `json:"maximumBytes"`
	State        string `json:"state"`
}

type TransportProfile struct {
	Identity            string              `json:"identity"`
	HeaderVersion       int                 `json:"headerVersion"`
	Source              InputTransport      `json:"source"`
	CanonicalInput      InputTransport      `json:"canonicalInput"`
	Completion          CompletionTransport `json:"completion"`
	RunnerDiagnostics   RunnerDiagnostics   `json:"runnerDiagnostics"`
	LimitBinding        string              `json:"limitBinding"`
	OversizeDisposition string              `json:"oversizeDisposition"`
}

type InputTransport struct {
	PayloadMaximumBytes  int `json:"payloadMaximumBytes"`
	HeaderBytes          int `json:"headerBytes"`
	PhysicalMaximumBytes int `json:"physicalMaximumBytes"`
	HostRunnerFD         int `json:"hostRunnerFd"`
	GuestLauncherFD      int `json:"guestLauncherFd"`
}

type CompletionTransport struct {
	JSONMaximumBytes     int  `json:"jsonMaximumBytes"`
	HeaderBytes          int  `json:"headerBytes"`
	TrailerBytes         int  `json:"trailerBytes"`
	PhysicalMaximumBytes int  `json:"physicalMaximumBytes"`
	RetainMaximumBytes   int  `json:"retainMaximumBytes"`
	HostRunnerFD         int  `json:"hostRunnerFd"`
	GuestLauncherFD      int  `json:"guestLauncherFd"`
	EOFIsCommit          bool `json:"eofIsCommit"`
}

type RunnerDiagnostics struct {
	StdoutRetainedPrefixBytes        int  `json:"stdoutRetainedPrefixBytes"`
	StderrRetainedPrefixBytes        int  `json:"stderrRetainedPrefixBytes"`
	ReadBufferBytes                  int  `json:"readBufferBytes"`
	MaximumCaptureFileBytesPerStream int  `json:"maximumCaptureFileBytesPerStream"`
	ContinuousDrain                  bool `json:"continuousDrain"`
}

type TeardownProfile struct {
	WallActionAtMs           int      `json:"wallActionAtMs"`
	CancellationAction       string   `json:"cancellationAction"`
	BestEffortGraceMs        int      `json:"bestEffortGraceMs"`
	BestEffortGraceClaim     string   `json:"bestEffortGraceClaim"`
	ForcedSignal             string   `json:"forcedSignal"`
	ForcedAbsenceDeadlineMs  int      `json:"forcedAbsenceDeadlineMs"`
	MaximumActionToAbsenceMs int      `json:"maximumActionToAbsenceMs"`
	IdentityBeforeSignal     []string `json:"identityBeforeSignal"`
	IdentityMismatch         string   `json:"identityMismatch"`
	ChildTree                string   `json:"childTree"`
	RunnerExit               string   `json:"runnerExit"`
	OrdinarySuccessRequires  []string `json:"ordinarySuccessRequires"`
}

type MatrixGroup struct {
	Group string   `json:"group"`
	Cases []string `json:"cases"`
}

type WorkStatus struct {
	C2A                     string `json:"c2a"`
	ParentGovernedRuntime   string `json:"parentGovernedRuntime"`
	C2B                     string `json:"c2b"`
	RuntimeProfileAdmission string `json:"runtimeProfileAdmission"`
	Runtime001              string `json:"runtime001"`
	VMM001                  string `json:"vmm001"`
}

type Effects struct {
	Process    bool `json:"process"`
	Runtime    bool `json:"runtime"`
	Adapter    bool `json:"adapter"`
	Backend    bool `json:"backend"`
	Guest      bool `json:"guest"`
	Credential bool `json:"credential"`
	Signing    bool `json:"signing"`
	Release    bool `json:"release"`
	Admission  bool `json:"admission"`
	Consumer   bool `json:"consumer"`
}

func Decode(exact []byte) (*Contract, error) {
	if len(exact) != ContractBytes {
		return nil, fmt.Errorf("C2A_CONTRACT_MISMATCH: byte length %d", len(exact))
	}
	digest := sha256.Sum256(exact)
	if hex.EncodeToString(digest[:]) != ContractSHA256 {
		return nil, errors.New("C2A_CONTRACT_MISMATCH: sha256")
	}
	decoder := json.NewDecoder(bytes.NewReader(exact))
	decoder.DisallowUnknownFields()
	var contract Contract
	if err := decoder.Decode(&contract); err != nil {
		return nil, fmt.Errorf("C2A_CONTRACT_MISMATCH: decode: %w", err)
	}
	if err := strictjson.RequireEOF(decoder, "C2A_CONTRACT_MISMATCH"); err != nil {
		return nil, err
	}
	if err := Validate(&contract); err != nil {
		return nil, err
	}
	return strictjson.Clone(&contract), nil
}

func Validate(c *Contract) error {
	if c == nil || c.Contract != ContractIdentity || c.Version != ContractVersion || c.Status != ContractStatus {
		return errors.New("C2A_CONTRACT_MISMATCH: identity, version, or status")
	}
	if c.C1Binding.Contract != "capsule.governed-deno-core.controlled-development-composition/c1" ||
		c.C1Binding.Fixture != "schemas/conformance/c1-governed-deno-core/controlled-development-profile.json" ||
		c.C1Binding.Bytes != C1Bytes || c.C1Binding.SHA256 != C1SHA256 ||
		c.C1Binding.Consumption != "exact-byte-read-only" || c.C1Binding.Mutation != "forbidden" {
		return errors.New("C2A_C1_MISMATCH")
	}
	expectedPlanRoles := []string{"runtimeBundleDigest", "profileRegistryDigest", "backendValidationDigest", "backendConfigurationDigest"}
	plan := c.C1PlanAndProfileBindings.ExecutionPlan
	profile := c.C1PlanAndProfileBindings.RuntimeProfile
	if plan.ObjectType != "capsule.execution-plan" || plan.ObjectVersion != 0 ||
		plan.MediaType != "application/capsule.execution-plan+cbor;v=0" ||
		plan.SourceProfile != "capsule.mjs-source/v0" ||
		plan.SourceBinding != "sourceManifestDigest+sourceByteLength+sourceEntrypoint" ||
		plan.InputEncoding != "capsule.canonical-inline-json/v0" ||
		plan.InputBinding != "inlineInputDigest+inlineInputByteLength" ||
		plan.OutputLimitBinding != "outputMaxJsonBytes" ||
		!slices.Equal(plan.RequiredReferenceRoles, expectedPlanRoles) ||
		profile.SelectedIdentity != nil || profile.SelectedDigest != nil ||
		profile.MachineProfileIdentity != "capsule.c2a.libkrun-hvf-controlled-development/v1" ||
		profile.TransportProfileIdentity != "capsule.c2a.c1-narrowed-p0-3/v1" ||
		profile.State != "missing-composed-runtime-profile-identity-refuse" {
		return errors.New("C2A_C1_MISMATCH: plan or runtime-profile binding")
	}
	if len(c.ArtifactClosure.Required) != len(expectedArtifactRoles) {
		return errors.New("C2A_ARTIFACT_IDENTITY_MISSING: role count")
	}
	for i, artifact := range c.ArtifactClosure.Required {
		if artifact.Role != expectedArtifactRoles[i] || artifact.SHA256 != nil || artifact.Admitted {
			return fmt.Errorf("C2A_ARTIFACT_IDENTITY_MISSING: %s", artifact.Role)
		}
	}
	if err := validateKnownAnswer(c.KnownAnswer); err != nil {
		return err
	}
	if c.RunnerDescriptorProfile.NumericRange != (NumericRange{Minimum: 0, Maximum: 7, CloseFromInclusive: 8}) || len(c.RunnerDescriptorProfile.Entries) != 8 {
		return errors.New("C2A_DESCRIPTOR_MANIFEST_MISMATCH: runner range")
	}
	for i, entry := range c.RunnerDescriptorProfile.Entries {
		if entry.FD != i || entry.Role != expectedRunnerRoles[i] || entry.AccessMode != expectedRunnerModes[i] {
			return fmt.Errorf("C2A_DESCRIPTOR_MANIFEST_MISMATCH: runner fd %d", i)
		}
	}
	expectedPorts := []PortCall{
		{PortID: 0, Name: "capsule.source", InputFD: 5, OutputFD: -1, GuestNode: "/dev/hvc0"},
		{PortID: 1, Name: "capsule.input", InputFD: 6, OutputFD: -1, GuestNode: "/dev/vport0p1"},
		{PortID: 2, Name: "capsule.completion", InputFD: -1, OutputFD: 7, GuestNode: "/dev/vport0p2"},
	}
	if !slices.Equal(c.RunnerDescriptorProfile.PortCalls, expectedPorts) {
		return errors.New("C2A_DESCRIPTOR_MANIFEST_MISMATCH: ports")
	}
	if len(c.GuestDescriptorProfile.PreChild) != 6 || c.GuestDescriptorProfile.CloseFromInclusive != 6 || c.GuestDescriptorProfile.WorkloadCompletionAccess != "none" || string(c.GuestDescriptorProfile.ChildManifest) != "null" {
		return errors.New("C2A_DESCRIPTOR_MANIFEST_MISMATCH: guest")
	}
	for i, entry := range c.GuestDescriptorProfile.PreChild {
		if entry.FD != i || entry.Role != expectedGuestRoles[i] || entry.AccessMode != expectedGuestModes[i] {
			return fmt.Errorf("C2A_DESCRIPTOR_MANIFEST_MISMATCH: guest fd %d", i)
		}
	}
	if err := validateMachine(c.MachineProfile); err != nil {
		return err
	}
	if err := validateTransport(c.TransportProfile); err != nil {
		return err
	}
	if c.TeardownProfile.WallActionAtMs != 1000 || c.TeardownProfile.BestEffortGraceMs != 200 || c.TeardownProfile.ForcedSignal != "SIGKILL" || c.TeardownProfile.ForcedAbsenceDeadlineMs != 1000 || c.TeardownProfile.MaximumActionToAbsenceMs != 1200 || c.TeardownProfile.RunnerExit != "lifecycle-only-never-workload-success" {
		return errors.New("C2A_TEARDOWN_UNSUPPORTED")
	}
	if len(c.C2BMatrix) != len(expectedMatrixGroups) {
		return errors.New("C2A_C2B_MATRIX_MISMATCH")
	}
	seen := make(map[string]bool, len(c.C2BMatrix))
	for _, group := range c.C2BMatrix {
		count, ok := expectedMatrixGroups[group.Group]
		if !ok || len(group.Cases) != count || seen[group.Group] {
			return fmt.Errorf("C2A_C2B_MATRIX_MISMATCH: %s", group.Group)
		}
		seen[group.Group] = true
	}
	if len(c.RestorationMutations) != 18 {
		return errors.New("C2A_RESTORATION_MATRIX_MISMATCH")
	}
	for i, mutation := range c.RestorationMutations {
		if len(mutation) < 7 || mutation[:7] != fmt.Sprintf("MUT-%03d", i+1) {
			return fmt.Errorf("C2A_RESTORATION_MATRIX_MISMATCH: %d", i+1)
		}
	}
	if c.WorkStatus != (WorkStatus{C2A: "PASSED-passive-preparation-only", ParentGovernedRuntime: "IN_PROGRESS-TRENDING_GOOD", C2B: "BLOCKED", RuntimeProfileAdmission: "BLOCKED", Runtime001: "unsupported", VMM001: "unsupported"}) {
		return errors.New("C2A_STATUS_MISMATCH")
	}
	if c.Effects != (Effects{}) {
		return errors.New("C2A_RUNTIME_NOT_ADMITTED: authority effect")
	}
	return nil
}

func validateKnownAnswer(k KnownAnswer) error {
	checks := []struct {
		value  string
		bytes  int
		digest string
	}{
		{k.Source.UTF8, k.Source.Bytes, k.Source.SHA256},
		{k.CanonicalInput.UTF8, k.CanonicalInput.Bytes, k.CanonicalInput.SHA256},
		{k.ExpectedCompletion.JSONUTF8, k.ExpectedCompletion.JSONBytes, k.ExpectedCompletion.JSONSHA256},
	}
	for _, check := range checks {
		sum := sha256.Sum256([]byte(check.value))
		if len([]byte(check.value)) != check.bytes || hex.EncodeToString(sum[:]) != check.digest {
			return errors.New("C2A_KNOWN_ANSWER_MISMATCH")
		}
	}
	if k.Source.Path != "main.mjs" || k.Source.SourceManifestBytes != 89 || k.Source.SourceManifestSHA256 != "712b1bd9739e4f6b0b027600207cbb08fb21b159a57bd34a15cf0ff8f32661b0" || k.ExpectedCompletion.TerminalStatus != "succeeded" || k.ExpectedCompletion.Commit != "required-last" || k.ExpectedCompletion.RunnerExitAlone != "never-success" {
		return errors.New("C2A_KNOWN_ANSWER_MISMATCH")
	}
	return nil
}

func validateMachine(m MachineProfile) error {
	if m.Architecture != "aarch64" || m.VCPUs != 1 || m.GuestRAMMiB != 256 || m.WallTimeMs != 1000 || m.Concurrency != 1 || m.CPUTimeLimitMs != nil || m.HostVMMMemoryLimitBytes != nil || m.Scratch.MaximumBytes != nil || m.CPUTimePolicy != "unsupported-refuse" || m.HostVMMMemoryPolicy != "unsupported-refuse" || m.Scratch.State != "inactive-refuse-missing-exact-enforcement-evidence" || m.UnknownOrDifferentValue != "refuse-no-default-no-clamp" {
		return errors.New("C2A_RESOURCE_UNSUPPORTED")
	}
	return nil
}

func validateTransport(t TransportProfile) error {
	if t.HeaderVersion != 1 || t.Source != (InputTransport{262144, 152, 262296, 5, 3}) || t.CanonicalInput != (InputTransport{262144, 152, 262296, 6, 4}) || t.Completion != (CompletionTransport{262144, 160, 64, 262368, 262369, 7, 5, false}) || t.RunnerDiagnostics != (RunnerDiagnostics{4096, 4096, 32768, 4194, true}) || t.OversizeDisposition != "drain-all-retain-cap-plus-one-refuse" {
		return errors.New("C2A_TRANSPORT_MISMATCH")
	}
	return nil
}
