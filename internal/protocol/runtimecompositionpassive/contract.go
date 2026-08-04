package runtimecompositionpassive

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
)

const (
	ContractIdentity = "capsule.governed-deno-core.controlled-development-composition/c1"
	ContractVersion  = 1
	ContractStatus   = "passive-not-admitted"

	ContractBytes       = 9289
	ContractSHA256      = "d5d75e638a15be6c9f4a3230d17309d085f6ec103a73b64d9e0fd656a5423c9e"
	MaximumSourceBytes  = 262144
	MaximumInputBytes   = 262144
	MaximumResultBytes  = 262144
	MaximumFrameBytes   = 262368
	RuntimeRootEntries  = 22
	RuntimeTarget       = "aarch64-unknown-linux-gnu"
	DenoCoreVersion     = "0.409.0"
	ExperimentsMerge    = "fa03d7043b4f0653081d6c5733d597f49f6efd1c"
	DenoHead            = "9adb0b68b55bca81644827f1e7749a3acb091bed"
	DenoMerge           = "ea18b9dc21ff8ebd19347be7095f47937ee14ec2"
	DenoUpstream        = "14eea3160ae5834476aa3b9d317b8d41d991b982"
	RustyV8Head         = "80e863ddb942a4aa2b384e794fc23e35b9d2bb15"
	RustyV8Merge        = "cbf56de2e1156b1cf1561fdbaea7172a0aa056f4"
	RustyV8Upstream     = "d305e6afa7736f6e298c30ae6646f7709ee9382b"
	SnapshotDigest      = "4e8965217d5a6675a880326eee6f690bbeec7e7cb243decf2f3e9f453a871a2c"
	RootManifestDigest  = "100832dbb37737f29341bc5404df6d4405b8d6b706f274028892801fa88e7de8"
	TransportProfileRef = "capsule.gate-c.p0-3.measured-limits/v0"
)

var expectedArtifacts = map[string]Artifact{
	"rustyV8Raw":     {Role: "governed-rusty-v8-static-archive-raw", SHA256: "e964d6b1b3689e91f8cf488d8a9f05764a03434b2e2e8347be5067300d39a7de"},
	"rustyV8Gzip":    {Role: "governed-rusty-v8-static-archive-gzip", SHA256: "1ae209c9e4ba5803d010d2c79ee4cc0af0126c5a7ebcca211c7e41deaede4cd2"},
	"denoCoreBinary": {Role: "governed-deno-core-binary", SHA256: "56d3acefd2cc2f5136a0b8143c47131e49a58fbf66382dfd3e84f715ce8e2898"},
	"snapshot":       {Role: "governed-deno-core-startup-snapshot", SHA256: SnapshotDigest},
	"twoFileBundle":  {Role: "governed-deno-core-two-file-bundle", SHA256: "0cc08f93e82fcfe68b033e8807975a3bd67068a817da811a87a73aedc3f23937"},
	"rootManifest":   {Role: "governed-runtime-root-manifest", SHA256: RootManifestDigest},
	"rootTar":        {Role: "governed-runtime-root-tar", SHA256: "9c46b45c4d220aedcc47c9ee53e875bc71d31d0b881b51740aaa9b882b5741e6"},
	"rootGzip":       {Role: "governed-runtime-root-gzip", SHA256: "e847651b35cd425dd8f6fe3bd45d693aff0af244e3a7bd30c629fa125cac62e8"},
}

var expectedOps = []string{
	"op_get_ext_import_meta_proto",
	"op_get_extras_binding_object",
	"op_set_captured_bootstrap",
}

var expectedPermittedGlobals = []string{
	"AggregateError", "Array", "ArrayBuffer", "Boolean", "DataView", "Error", "EvalError",
	"Float32Array", "Float64Array", "Infinity", "Int16Array", "Int32Array", "Int8Array", "JSON",
	"Map", "Math", "NaN", "Number", "Object", "Promise", "Proxy", "RangeError", "ReferenceError",
	"Reflect", "RegExp", "Set", "String", "Symbol", "SyntaxError", "TypeError", "URIError",
	"Uint16Array", "Uint32Array", "Uint8Array", "Uint8ClampedArray", "WeakMap", "WeakSet",
	"decodeURI", "decodeURIComponent", "encodeURI", "encodeURIComponent", "globalThis", "isFinite",
	"isNaN", "parseFloat", "parseInt", "undefined",
}

var expectedForbiddenGlobals = []string{
	"Atomics", "BigInt", "BigInt64Array", "BigUint64Array", "BroadcastChannel", "Date", "Deno",
	"EventSource", "FinalizationRegistry", "Function", "Intl", "SharedArrayBuffer", "Temporal",
	"TextDecoder", "TextEncoder", "URL", "URLPattern", "WeakRef", "WebAssembly", "WebSocket",
	"Worker", "__bootstrap", "clearInterval", "clearTimeout", "console", "crypto", "eval", "fetch",
	"localStorage", "navigator", "performance", "process", "queueMicrotask", "sessionStorage",
	"setInterval", "setTimeout",
}

var expectedSuccessRequirements = []string{
	"source-binding-valid",
	"input-binding-valid",
	"runtime-integrity-valid",
	"main-module-evaluated",
	"capsuleMain-called-once",
	"result-json-valid-and-bounded",
	"completion-frame-committed",
	"runner-lifecycle-terminal",
	"teardown-authoritative",
}

var expectedDescriptorRoles = []DescriptorRole{
	{Role: "runtime-root", Owner: "supervisor-to-host-runner", Direction: "read-only-host-custody", WorkloadAccess: "none"},
	{Role: "registered-source", Owner: "host-writer-to-guest-launcher-reader", Direction: "host-to-guest", WorkloadAccess: "indirect-input-only"},
	{Role: "approved-inline-input", Owner: "host-writer-to-guest-launcher-reader", Direction: "host-to-guest", WorkloadAccess: "indirect-input-only"},
	{Role: "completion-result", Owner: "guest-launcher-writer-to-host-reader", Direction: "guest-to-host", WorkloadAccess: "none"},
}

var expectedPlanReferenceRoles = []string{
	"runtimeBundleDigest",
	"profileRegistryDigest",
	"backendValidationDigest",
	"backendConfigurationDigest",
}

var expectedRefusalCodes = []string{
	"C1_CONTRACT_MISMATCH",
	"C1_ARTIFACT_MISMATCH",
	"C1_SOURCE_MISMATCH",
	"C1_INPUT_MISMATCH",
	"C1_GLOBAL_SURFACE_MISMATCH",
	"C1_OP_SURFACE_MISMATCH",
	"C1_LOADER_PRESENT",
	"C1_MAIN_BINDING_INVALID",
	"C1_RESULT_INVALID",
	"C1_DESCRIPTOR_PROFILE_UNSELECTED",
	"C1_RESOURCE_PROFILE_UNSELECTED",
	"C1_RUNTIME_NOT_ADMITTED",
}

var expectedStopConditions = []string{
	"Any governed fork commit or artifact digest differs.",
	"Any built-in op, extension, module, loader, or permitted global is added.",
	"Any forbidden global remains reachable or string code generation remains usable.",
	"Source or input bytes come from outside committed registration state.",
	"The workload owns the completion endpoint or can forge the commit trailer.",
	"A live host path, filesystem API, network, subprocess, FFI, inspector, Worker, WebAssembly, or native loader is reachable.",
	"C2 cannot close numeric descriptor ownership without ambient authority.",
	"C2 cannot select exact vCPU, RAM, wall-time, scratch, console, and teardown mechanisms from retained composed evidence.",
	"P0-1 root custody, P0-2 device closure, P0-3 transport, or P0-4 installed composition is incomplete.",
	"A runtime/profile admission record for these exact composed bytes is absent, expired, or mismatched.",
}

type Contract struct {
	Contract       string              `json:"contract"`
	Version        int                 `json:"version"`
	Status         string              `json:"status"`
	Evidence       Evidence            `json:"evidence"`
	Forks          Forks               `json:"forks"`
	Artifacts      map[string]Artifact `json:"artifacts"`
	Application    Application         `json:"application"`
	RuntimeSurface RuntimeSurface      `json:"runtimeSurface"`
	Descriptors    Descriptors         `json:"descriptors"`
	Resources      Resources           `json:"resources"`
	RefusalCodes   []string            `json:"refusalCodes"`
	StopConditions []string            `json:"stopConditions"`
	Effects        Effects             `json:"effects"`
}

type Evidence struct {
	Repository  string `json:"repository"`
	PullRequest int    `json:"pullRequest"`
	MergeCommit string `json:"mergeCommit"`
	HandoffPath string `json:"handoffPath"`
}

type Forks struct {
	Deno    Fork `json:"deno"`
	RustyV8 Fork `json:"rustyV8"`
}

type Fork struct {
	Repository string `json:"repository"`
	Head       string `json:"head"`
	Merge      string `json:"merge"`
	Upstream   string `json:"upstream"`
}

type Artifact struct {
	Role     string `json:"role"`
	SHA256   string `json:"sha256"`
	Admitted bool   `json:"admitted"`
}

type Application struct {
	Source     Source     `json:"source"`
	Input      Input      `json:"input"`
	Main       Main       `json:"main"`
	Completion Completion `json:"completion"`
}

type Source struct {
	Profile        string `json:"profile"`
	MediaType      string `json:"mediaType"`
	FileCount      int    `json:"fileCount"`
	Path           string `json:"path"`
	Entrypoint     string `json:"entrypoint"`
	MinimumBytes   int    `json:"minimumBytes"`
	MaximumBytes   int    `json:"maximumBytes"`
	Binding        string `json:"binding"`
	Transformation string `json:"transformation"`
}

type Input struct {
	Slot         string `json:"slot"`
	MediaType    string `json:"mediaType"`
	Encoding     string `json:"encoding"`
	MaximumBytes int    `json:"maximumBytes"`
	Binding      string `json:"binding"`
	Injection    string `json:"injection"`
}

type Main struct {
	Binding      string   `json:"binding"`
	Receiver     string   `json:"receiver"`
	RequiredType string   `json:"requiredType"`
	Invocations  int      `json:"invocations"`
	AwaitResult  bool     `json:"awaitResult"`
	Arguments    []string `json:"arguments"`
}

type Completion struct {
	Slot                      string   `json:"slot"`
	MediaType                 string   `json:"mediaType"`
	Serialization             string   `json:"serialization"`
	MaximumJSONBytes          int      `json:"maximumJsonBytes"`
	MaximumPhysicalFrameBytes int      `json:"maximumPhysicalFrameBytes"`
	LimitBinding              string   `json:"limitBinding"`
	AttemptBinding            string   `json:"attemptBinding"`
	WorkloadOwnsEndpoint      bool     `json:"workloadOwnsEndpoint"`
	SuccessRequires           []string `json:"successRequires"`
}

type RuntimeSurface struct {
	Target           string        `json:"target"`
	DenoCore         string        `json:"denoCore"`
	SnapshotSHA256   string        `json:"snapshotSha256"`
	V8Flags          []string      `json:"v8Flags"`
	Extensions       []string      `json:"extensions"`
	BuiltinOps       []string      `json:"builtinOps"`
	ModuleLoader     string        `json:"moduleLoader"`
	Modules          []Module      `json:"modules"`
	PermittedGlobals []string      `json:"permittedGlobals"`
	ForbiddenGlobals []string      `json:"forbiddenGlobals"`
	WorkloadFiles    WorkloadFiles `json:"workloadFiles"`
}

type Module struct {
	Identity string `json:"identity"`
	Role     string `json:"role"`
	Source   string `json:"source"`
}

type WorkloadFiles struct {
	Readable                  []string `json:"readable"`
	Writable                  []string `json:"writable"`
	RuntimeRootManifestSHA256 string   `json:"runtimeRootManifestSha256"`
	RuntimeRootEntries        int      `json:"runtimeRootEntries"`
}

type Descriptors struct {
	ConstructionProbeObserved []int            `json:"constructionProbeObserved"`
	LogicalRoles              []DescriptorRole `json:"logicalRoles"`
	NumericAssignment         string           `json:"numericAssignment"`
	WorkloadCompletionAccess  string           `json:"workloadCompletionAccess"`
}

type DescriptorRole struct {
	Role           string `json:"role"`
	Owner          string `json:"owner"`
	Direction      string `json:"direction"`
	WorkloadAccess string `json:"workloadAccess"`
}

type Resources struct {
	TransportProfileRef string   `json:"transportProfileRef"`
	PlanReferenceRoles  []string `json:"planReferenceRoles"`
	MachineProfileRef   *string  `json:"machineProfileRef"`
	Activation          string   `json:"activation"`
}

type Effects struct {
	Process   bool `json:"process"`
	Runtime   bool `json:"runtime"`
	Backend   bool `json:"backend"`
	Guest     bool `json:"guest"`
	Admission bool `json:"admission"`
	Signing   bool `json:"signing"`
	Release   bool `json:"release"`
}

// Decode rejects unknown JSON keys (DisallowUnknownFields) and trailing bytes (requireEOF), but
// encoding/json has no equivalent guard for *omitted* keys: a missing field silently decodes to
// its Go zero value. That is safe only because every field's zero value already coincides with
// its expected/safe constant here (empty slices, false Effects, zero byte counts) — Validate then
// rejects the zero value the same way it would reject any other wrong value. This is an
// invariant of the current field set, not a structural guarantee. A future field whose
// safe/expected value is non-zero (for example a boolean that must be true to prove a property)
// must not rely on this coincidence; it needs its own explicit presence check, or Decode should
// gain a general required-key check mirroring governed-deno-core-c1-composition.schema.json's
// `required` lists.
func Decode(exact []byte) (*Contract, error) {
	if len(exact) != ContractBytes {
		return nil, fmt.Errorf("C1_CONTRACT_MISMATCH: byte length %d", len(exact))
	}
	digest := sha256.Sum256(exact)
	if hex.EncodeToString(digest[:]) != ContractSHA256 {
		return nil, errors.New("C1_CONTRACT_MISMATCH: sha256")
	}
	decoder := json.NewDecoder(bytes.NewReader(exact))
	decoder.DisallowUnknownFields()
	var contract Contract
	if err := decoder.Decode(&contract); err != nil {
		return nil, fmt.Errorf("C1_CONTRACT_MISMATCH: decode: %w", err)
	}
	if err := requireEOF(decoder); err != nil {
		return nil, err
	}
	if err := Validate(&contract); err != nil {
		return nil, err
	}
	return clone(&contract), nil
}

func Validate(contract *Contract) error {
	if contract == nil || contract.Contract != ContractIdentity || contract.Version != ContractVersion || contract.Status != ContractStatus {
		return errors.New("C1_CONTRACT_MISMATCH: identity")
	}
	if contract.Evidence.Repository != "Shrimpworks/capsule-experiments" || contract.Evidence.PullRequest != 1 || contract.Evidence.MergeCommit != ExperimentsMerge || contract.Evidence.HandoffPath != "experiments/gate-c-fork-native-deno-runtime-bundle/HANDOFF.md" {
		return errors.New("C1_CONTRACT_MISMATCH: evidence")
	}
	if contract.Forks.Deno != (Fork{Repository: "Shrimpworks/deno", Head: DenoHead, Merge: DenoMerge, Upstream: DenoUpstream}) || contract.Forks.RustyV8 != (Fork{Repository: "Shrimpworks/rusty_v8", Head: RustyV8Head, Merge: RustyV8Merge, Upstream: RustyV8Upstream}) {
		return errors.New("C1_ARTIFACT_MISMATCH: fork")
	}
	if len(contract.Artifacts) != len(expectedArtifacts) {
		return errors.New("C1_ARTIFACT_MISMATCH: count")
	}
	for name, expected := range expectedArtifacts {
		actual, ok := contract.Artifacts[name]
		if !ok || actual.Role != expected.Role || actual.SHA256 != expected.SHA256 || actual.Admitted {
			return fmt.Errorf("C1_ARTIFACT_MISMATCH: %s", name)
		}
	}
	if contract.Application.Source.Profile != "capsule.mjs-source/v0" || contract.Application.Source.MediaType != "application/capsule.javascript-source;v=0;module=esm" || contract.Application.Source.Path != "main.mjs" || contract.Application.Source.Entrypoint != "main.mjs" || contract.Application.Source.FileCount != 1 || contract.Application.Source.MinimumBytes != 0 || contract.Application.Source.MaximumBytes != MaximumSourceBytes || contract.Application.Source.Binding != "ExecutionPlan.sourceManifestDigest+sourceByteLength+sourceEntrypoint" || contract.Application.Source.Transformation != "none-byte-exact" {
		return errors.New("C1_SOURCE_MISMATCH")
	}
	if contract.Application.Input.Slot != "primary-data" || contract.Application.Input.MediaType != "application/json" || contract.Application.Input.Encoding != "capsule.canonical-inline-json/v0" || contract.Application.Input.MaximumBytes != MaximumInputBytes || contract.Application.Input.Binding != "ExecutionPlan.inlineInputDigest+inlineInputByteLength" || contract.Application.Input.Injection != "single-copied-argument" {
		return errors.New("C1_INPUT_MISMATCH")
	}
	if contract.Application.Main.Binding != "globalThis.capsuleMain" || contract.Application.Main.Receiver != "globalThis" || contract.Application.Main.RequiredType != "function" || contract.Application.Main.Invocations != 1 || !contract.Application.Main.AwaitResult || !equalStrings(contract.Application.Main.Arguments, []string{"parsed-canonical-inline-json"}) {
		return errors.New("C1_MAIN_BINDING_INVALID")
	}
	if contract.Application.Completion.Slot != "transformed-json" || contract.Application.Completion.MediaType != "application/json" || contract.Application.Completion.Serialization != "JSON.stringify-exact-utf8" || contract.Application.Completion.MaximumJSONBytes != MaximumResultBytes || contract.Application.Completion.MaximumPhysicalFrameBytes != MaximumFrameBytes || contract.Application.Completion.LimitBinding != "ExecutionPlan.outputMaxJsonBytes" || contract.Application.Completion.AttemptBinding != "P0-3-attempt-bound-completion-frame" || contract.Application.Completion.WorkloadOwnsEndpoint || !equalStrings(contract.Application.Completion.SuccessRequires, expectedSuccessRequirements) {
		return errors.New("C1_RESULT_INVALID")
	}
	if contract.RuntimeSurface.Target != RuntimeTarget || contract.RuntimeSurface.DenoCore != DenoCoreVersion || contract.RuntimeSurface.SnapshotSHA256 != SnapshotDigest || contract.RuntimeSurface.ModuleLoader != "none" || len(contract.RuntimeSurface.Extensions) != 0 || !equalStrings(contract.RuntimeSurface.BuiltinOps, expectedOps) {
		return errors.New("C1_OP_SURFACE_MISMATCH")
	}
	if !equalStrings(contract.RuntimeSurface.V8Flags, []string{"--jitless", "--random-seed=42"}) || len(contract.RuntimeSurface.Modules) != 1 || contract.RuntimeSurface.Modules[0] != (Module{Identity: "capsule:main.mjs", Role: "main", Source: "registered-main.mjs"}) {
		return errors.New("C1_LOADER_PRESENT")
	}
	if !equalStrings(contract.RuntimeSurface.PermittedGlobals, expectedPermittedGlobals) || !equalStrings(contract.RuntimeSurface.ForbiddenGlobals, expectedForbiddenGlobals) || overlaps(contract.RuntimeSurface.PermittedGlobals, contract.RuntimeSurface.ForbiddenGlobals) {
		return errors.New("C1_GLOBAL_SURFACE_MISMATCH")
	}
	files := contract.RuntimeSurface.WorkloadFiles
	if len(files.Readable) != 0 || len(files.Writable) != 0 || files.RuntimeRootManifestSHA256 != RootManifestDigest || files.RuntimeRootEntries != RuntimeRootEntries {
		return errors.New("C1_ARTIFACT_MISMATCH: workload-files")
	}
	if !equalInts(contract.Descriptors.ConstructionProbeObserved, []int{0, 1, 2}) || !equalDescriptorRoles(contract.Descriptors.LogicalRoles, expectedDescriptorRoles) || contract.Descriptors.NumericAssignment != "c2-required-unselected" || contract.Descriptors.WorkloadCompletionAccess != "none" {
		return errors.New("C1_DESCRIPTOR_PROFILE_UNSELECTED")
	}
	if contract.Resources.TransportProfileRef != TransportProfileRef || !equalStrings(contract.Resources.PlanReferenceRoles, expectedPlanReferenceRoles) || contract.Resources.MachineProfileRef != nil || contract.Resources.Activation != "refuse-until-c2-and-admission" {
		return errors.New("C1_RESOURCE_PROFILE_UNSELECTED")
	}
	if !equalStrings(contract.RefusalCodes, expectedRefusalCodes) || !equalStrings(contract.StopConditions, expectedStopConditions) {
		return errors.New("C1_CONTRACT_MISMATCH: refusal-boundary")
	}
	if contract.Effects.Process || contract.Effects.Runtime || contract.Effects.Backend || contract.Effects.Guest || contract.Effects.Admission || contract.Effects.Signing || contract.Effects.Release {
		return errors.New("C1_RUNTIME_NOT_ADMITTED")
	}
	return nil
}

func requireEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("C1_CONTRACT_MISMATCH: trailing data")
	}
	return nil
}

func clone(contract *Contract) *Contract {
	exact, err := json.Marshal(contract)
	if err != nil {
		panic(err)
	}
	var copied Contract
	if err := json.Unmarshal(exact, &copied); err != nil {
		panic(err)
	}
	return &copied
}

func equalStrings(left, right []string) bool {
	return slices.Equal(left, right)
}

func equalInts(left, right []int) bool {
	return slices.Equal(left, right)
}

func equalDescriptorRoles(left, right []DescriptorRole) bool {
	return slices.Equal(left, right)
}

func overlaps(left, right []string) bool {
	seen := make(map[string]struct{}, len(left))
	for _, value := range left {
		seen[value] = struct{}{}
	}
	for _, value := range right {
		if _, ok := seen[value]; ok {
			return true
		}
	}
	return false
}
