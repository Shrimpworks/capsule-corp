package hostrunnermaterialized

import (
	"bytes"
	"crypto/sha256"
	"debug/macho"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"slices"
	"strings"
)

const (
	ProfileMaximumBytes = 16 * 1024
	ProfileBytes        = 10301
	ProfileSHA256       = "198688bacd50aaee4f57b4cd7c56cea6b939c10aa220fbbeba7d315de820d1fd"
	ComposedSHA256      = "e390085caaaba73ebc19f95bc9871305e4f9268c2283d7394133fa4491f4ba82"
)

type artifactKnownAnswer struct {
	path   string
	bytes  int
	sha256 string
}

var artifactKnownAnswers = []artifactKnownAnswer{
	{"libkrun.h", 54658, "dce44d1d70ab770b1089e57646e025281a4137fe5052b9dd8eaefb80c01a1bd8"},
	{"libkrun-abi-audit.c", 2512, "419256ea91de9b5e5323e1f1d6d42afb0a5fa85a8835d0d0404734af0ee92356"},
	{"libkrun.1.dylib", 4393448, "055d9d18dc964fec4aba21948c4a344cb7a51cb48a2c70017484b718eae12f9f"},
	{"capsule-host-runner.c", 7917, "5a5560fa667390253bf504d7c045fcbcc304fa5829b22a8acf1fff00a8e37eb9"},
	{"capsule-host-runner", 100488, "a30e3f7cba5f480b6e164536854749b5e1ba3349f20af6c9c8e5d2590bffe1ad"},
}

var orderedSourceCalls = []string{
	"krun_create_ctx()",
	"krun_set_vm_config(context, 1, 256)",
	"krun_disable_implicit_console(context)",
	"krun_disable_implicit_init(context)",
	"krun_disable_implicit_vsock(context)",
	"krun_add_read_only_raw_root_fd(",
	"krun_set_root_disk_remount(",
	"krun_add_virtio_console_multiport(context)",
	"\"capsule.source\", 5, -1",
	"\"capsule.input\", 6, -1",
	"\"capsule.completion\", -1, 7",
	"krun_set_kernel_console(context, \"hvc0\")",
	"krun_set_workdir(context, \"/\")",
	"krun_set_exec(",
	"write_ready()",
	"require_start_authorization()",
	"krun_start_enter(context)",
}

var forbiddenSourceCalls = []string{
	"krun_set_kernel(", "krun_set_firmware(", "krun_add_disk(", "krun_set_root(",
	"krun_add_virtiofs(", "krun_add_vsock(", "krun_add_net_unixstream(",
	"krun_set_gpu_options(", "krun_set_snd_device(", "open(", "socket(",
	"execve(", "posix_spawn(", "dlopen(",
}

var expectedRunnerKrunImports = []string{
	"_krun_add_console_port_inout",
	"_krun_add_read_only_raw_root_fd",
	"_krun_add_virtio_console_multiport",
	"_krun_create_ctx",
	"_krun_disable_implicit_console",
	"_krun_disable_implicit_init",
	"_krun_disable_implicit_vsock",
	"_krun_set_exec",
	"_krun_set_kernel_console",
	"_krun_set_root_disk_remount",
	"_krun_set_vm_config",
	"_krun_set_workdir",
	"_krun_start_enter",
}

type Profile struct {
	Identity        string          `json:"identity"`
	Status          string          `json:"status"`
	Predecessors    json.RawMessage `json:"predecessors"`
	GovernedSources json.RawMessage `json:"governedSources"`
	AcceptedHeader  json.RawMessage `json:"acceptedHeader"`
	LibkrunBuild    json.RawMessage `json:"libkrunBuild"`
	ABIReview       json.RawMessage `json:"abiReview"`
	BootRole        BootRole        `json:"bootRole"`
	Runner          Runner          `json:"runner"`
	ComposedProfile ComposedProfile `json:"composedProfile"`
	WorkStatus      WorkStatus      `json:"workStatus"`
	NextGuestGate   NextGuestGate   `json:"nextGuestGate"`
	Effects         Effects         `json:"effects"`
}

type BootRole struct {
	Mode                              string          `json:"mode"`
	LibkrunFeatures                   []string        `json:"libkrunFeatures"`
	Libkrunfw                         json.RawMessage `json:"libkrunfw"`
	Kernel                            json.RawMessage `json:"kernel"`
	SeparateFirmware                  string          `json:"separateFirmware"`
	RunnerCallsKernelOrFirmwareSetter bool            `json:"runnerCallsKernelOrFirmwareSetter"`
}

type Runner struct {
	Source           json.RawMessage `json:"source"`
	Artifact         json.RawMessage `json:"artifact"`
	Authority        Authority       `json:"authority"`
	Preflight        Preflight       `json:"preflight"`
	Ports            []Port          `json:"ports"`
	OrderedCalls     []string        `json:"orderedCalls"`
	ForbiddenImports []string        `json:"forbiddenImports"`
}

type Authority struct {
	Owner                           string `json:"owner"`
	Cardinality                     string `json:"cardinality"`
	AcceptsAttemptIDBytes           bool   `json:"acceptsAttemptIdBytes"`
	AcceptsReplacementConfiguration bool   `json:"acceptsReplacementConfiguration"`
	DaemonRoute                     bool   `json:"daemonRoute"`
	Service                         bool   `json:"service"`
	PrivilegedHelper                bool   `json:"privilegedHelper"`
	TeardownOwner                   string `json:"teardownOwner"`
}

type Preflight struct {
	Argv               []string `json:"argv"`
	Environment        string   `json:"environment"`
	AllowedFDs         []int    `json:"allowedFds"`
	CloseFromInclusive int      `json:"closeFromInclusive"`
	CloseMechanism     string   `json:"closeMechanism"`
	RootFD             int      `json:"rootFd"`
	RootType           string   `json:"rootType"`
	RootBytes          int      `json:"rootBytes"`
	RootSHA256         string   `json:"rootSha256"`
	ReadyByte          string   `json:"readyByte"`
	StartAuthorization string   `json:"startAuthorization"`
}

type Port struct {
	PortID    int    `json:"portId"`
	Name      string `json:"name"`
	InputFD   int    `json:"inputFd"`
	OutputFD  int    `json:"outputFd"`
	GuestNode string `json:"guestNode"`
}

type ComposedProfile struct {
	Runtime         json.RawMessage `json:"runtime"`
	TrustedInit     json.RawMessage `json:"trustedInit"`
	TrustedLauncher json.RawMessage `json:"trustedLauncher"`
	Devices         []string        `json:"devices"`
	Resources       json.RawMessage `json:"resources"`
	Teardown        json.RawMessage `json:"teardown"`
	DigestAlgorithm string          `json:"digestAlgorithm"`
	DigestSHA256    string          `json:"digestSha256"`
	Admitted        bool            `json:"admitted"`
}

type WorkStatus struct {
	MaterializationSlice    string `json:"materializationSlice"`
	FixedOwnedGuest         string `json:"fixedOwnedGuest"`
	RuntimeProfileAdmission string `json:"runtimeProfileAdmission"`
	ParentGovernedRuntime   string `json:"parentGovernedRuntime"`
	Runtime001              string `json:"runtime001"`
	VMM001                  string `json:"vmm001"`
}

type NextGuestGate struct {
	State               string `json:"state"`
	Blocker             string `json:"blocker"`
	GuestAuthorization  bool   `json:"guestAuthorization"`
	ConsumerImplemented bool   `json:"consumerImplemented"`
}

type Effects struct {
	ControlledCompilation bool `json:"controlledCompilation"`
	RunnerExecuted        bool `json:"runnerExecuted"`
	LibkrunLoaded         bool `json:"libkrunLoaded"`
	HVFCalled             bool `json:"hvfCalled"`
	VMCreated             bool `json:"vmCreated"`
	GuestCreated          bool `json:"guestCreated"`
	ArbitraryWorkload     bool `json:"arbitraryWorkload"`
	CredentialUsed        bool `json:"credentialUsed"`
	Signed                bool `json:"signed"`
	Installed             bool `json:"installed"`
	ReleasePublished      bool `json:"releasePublished"`
	ProductWired          bool `json:"productWired"`
	AdmissionChanged      bool `json:"admissionChanged"`
}

type Verified struct {
	profile Profile
}

func (verified *Verified) Profile() Profile {
	value := verified.profile
	value.Predecessors = bytes.Clone(value.Predecessors)
	value.GovernedSources = bytes.Clone(value.GovernedSources)
	value.AcceptedHeader = bytes.Clone(value.AcceptedHeader)
	value.LibkrunBuild = bytes.Clone(value.LibkrunBuild)
	value.ABIReview = bytes.Clone(value.ABIReview)
	value.BootRole.LibkrunFeatures = slices.Clone(value.BootRole.LibkrunFeatures)
	value.BootRole.Libkrunfw = bytes.Clone(value.BootRole.Libkrunfw)
	value.BootRole.Kernel = bytes.Clone(value.BootRole.Kernel)
	value.Runner.Source = bytes.Clone(value.Runner.Source)
	value.Runner.Artifact = bytes.Clone(value.Runner.Artifact)
	value.Runner.Preflight.Argv = slices.Clone(value.Runner.Preflight.Argv)
	value.Runner.Preflight.AllowedFDs = slices.Clone(value.Runner.Preflight.AllowedFDs)
	value.Runner.Ports = slices.Clone(value.Runner.Ports)
	value.Runner.OrderedCalls = slices.Clone(value.Runner.OrderedCalls)
	value.Runner.ForbiddenImports = slices.Clone(value.Runner.ForbiddenImports)
	value.ComposedProfile.Runtime = bytes.Clone(value.ComposedProfile.Runtime)
	value.ComposedProfile.TrustedInit = bytes.Clone(value.ComposedProfile.TrustedInit)
	value.ComposedProfile.TrustedLauncher = bytes.Clone(value.ComposedProfile.TrustedLauncher)
	value.ComposedProfile.Resources = bytes.Clone(value.ComposedProfile.Resources)
	value.ComposedProfile.Teardown = bytes.Clone(value.ComposedProfile.Teardown)
	value.ComposedProfile.Devices = slices.Clone(value.ComposedProfile.Devices)
	return value
}

func VerifyFixtureFS(fsys fs.FS) (*Verified, error) {
	profileBytes, err := readBounded(fsys, "materialized-profile.json", ProfileMaximumBytes)
	if err != nil {
		return nil, fmt.Errorf("C2B_V4_PROFILE_READ: %w", err)
	}
	if len(profileBytes) != ProfileBytes || digest(profileBytes) != ProfileSHA256 {
		return nil, errors.New("C2B_V4_PROFILE_KNOWN_ANSWER")
	}
	if err := verifyComposedDigest(profileBytes); err != nil {
		return nil, err
	}
	profile, err := decodeProfile(profileBytes)
	if err != nil {
		return nil, err
	}
	for _, answer := range artifactKnownAnswers {
		exact, readErr := readBounded(fsys, answer.path, answer.bytes)
		if readErr != nil {
			return nil, fmt.Errorf("C2B_V4_ARTIFACT_READ_%s: %w", answer.path, readErr)
		}
		if len(exact) != answer.bytes || digest(exact) != answer.sha256 {
			return nil, fmt.Errorf("C2B_V4_ARTIFACT_KNOWN_ANSWER_%s", answer.path)
		}
	}
	if err := verifyRunnerSource(fsys); err != nil {
		return nil, err
	}
	if err := verifyMachOArtifacts(fsys); err != nil {
		return nil, err
	}
	return &Verified{profile: profile}, nil
}

func decodeProfile(exact []byte) (Profile, error) {
	decoder := json.NewDecoder(bytes.NewReader(exact))
	decoder.DisallowUnknownFields()
	var profile Profile
	if err := decoder.Decode(&profile); err != nil {
		return Profile{}, fmt.Errorf("C2B_V4_PROFILE_DECODE: %w", err)
	}
	if err := requireEOF(decoder); err != nil {
		return Profile{}, err
	}
	if err := ValidateProfile(profile); err != nil {
		return Profile{}, err
	}
	return profile, nil
}

func ValidateProfile(profile Profile) error {
	if profile.Identity != "capsule.governed-deno-core.fixed-owned-guest-materialized/v4" ||
		profile.Status != "PASSED" {
		return errors.New("C2B_V4_IDENTITY")
	}
	if profile.BootRole.Mode != "non-efi" ||
		!slices.Equal(profile.BootRole.LibkrunFeatures, []string{"blk"}) ||
		profile.BootRole.SeparateFirmware != "inapplicable-no-path-authority" ||
		profile.BootRole.RunnerCallsKernelOrFirmwareSetter {
		return errors.New("C2B_V4_BOOT_ROLE")
	}
	wantAuthority := Authority{
		Owner: "execution-supervisor", Cardinality: "exactly-one-runner-per-committed-AttemptID",
		TeardownOwner: "execution-supervisor",
	}
	if profile.Runner.Authority != wantAuthority {
		return errors.New("C2B_V4_RUNNER_AUTHORITY")
	}
	wantFDs := []int{0, 1, 2, 3, 4, 5, 6, 7}
	if !slices.Equal(profile.Runner.Preflight.Argv, []string{"capsule-host-runner"}) ||
		profile.Runner.Preflight.Environment != "empty-after-removing-only-__CF_USER_TEXT_ENCODING" ||
		!slices.Equal(profile.Runner.Preflight.AllowedFDs, wantFDs) ||
		profile.Runner.Preflight.CloseFromInclusive != 8 || profile.Runner.Preflight.RootFD != 4 ||
		profile.Runner.Preflight.RootBytes != 134217728 ||
		profile.Runner.Preflight.RootSHA256 != "390a4786a20d45f1c691ec8c203f84f5e9d372a30e98f867cc8309a144ca6798" ||
		profile.Runner.Preflight.ReadyByte != "R" ||
		profile.Runner.Preflight.StartAuthorization != "exact-one-byte-G-followed-by-eof" {
		return errors.New("C2B_V4_PREFLIGHT")
	}
	wantPorts := []Port{
		{0, "capsule.source", 5, -1, "/dev/hvc0"},
		{1, "capsule.input", 6, -1, "/dev/vport0p1"},
		{2, "capsule.completion", -1, 7, "/dev/vport0p2"},
	}
	if !slices.Equal(profile.Runner.Ports, wantPorts) || len(profile.Runner.OrderedCalls) != 16 {
		return errors.New("C2B_V4_PORT_CALL_CONTRACT")
	}
	if profile.ComposedProfile.DigestSHA256 != ComposedSHA256 || profile.ComposedProfile.Admitted ||
		profile.WorkStatus.MaterializationSlice != "PASSED" ||
		profile.WorkStatus.FixedOwnedGuest != "BLOCKED" ||
		profile.WorkStatus.RuntimeProfileAdmission != "BLOCKED" ||
		profile.WorkStatus.ParentGovernedRuntime != "IN_PROGRESS — TRENDING_GOOD" ||
		profile.WorkStatus.Runtime001 != "unsupported" || profile.WorkStatus.VMM001 != "unsupported" {
		return errors.New("C2B_V4_STATUS_OR_ADMISSION")
	}
	if profile.NextGuestGate.State != "BLOCKED" || profile.NextGuestGate.GuestAuthorization ||
		profile.NextGuestGate.ConsumerImplemented ||
		profile.Effects != (Effects{ControlledCompilation: true}) {
		return errors.New("C2B_V4_EFFECT_OR_GUEST_GATE")
	}
	return nil
}

func verifyComposedDigest(exact []byte) error {
	needle := []byte(`"digestSha256": "` + ComposedSHA256 + `"`)
	if bytes.Count(exact, needle) != 1 {
		return errors.New("C2B_V4_COMPOSED_DIGEST_FIELD")
	}
	zero := []byte(`"digestSha256": "` + strings.Repeat("0", 64) + `"`)
	normalized := bytes.Replace(exact, needle, zero, 1)
	if digest(normalized) != ComposedSHA256 {
		return errors.New("C2B_V4_COMPOSED_DIGEST")
	}
	return nil
}

func verifyRunnerSource(fsys fs.FS) error {
	source, err := fs.ReadFile(fsys, "capsule-host-runner.c")
	if err != nil {
		return err
	}
	offset := 0
	for _, token := range orderedSourceCalls {
		index := bytes.Index(source[offset:], []byte(token))
		if index < 0 {
			return fmt.Errorf("C2B_V4_SOURCE_ORDER_%s", token)
		}
		offset += index + len(token)
	}
	for _, forbidden := range forbiddenSourceCalls {
		if bytes.Contains(source, []byte(forbidden)) {
			return fmt.Errorf("C2B_V4_SOURCE_FORBIDDEN_%s", forbidden)
		}
	}
	for _, required := range []string{
		"CAPSULE_ABI_POINTER_SLOTS 4096", "CAPSULE_HOST_CLOSE_FROM_INCLUSIVE 8",
		"status[4].st_nlink != 0", "status[4].st_mode & 07777", "pread(CAPSULE_ROOT_FD",
		"unsetenv(\"__CF_USER_TEXT_ENCODING\")", "environ[0] != NULL",
	} {
		if !bytes.Contains(source, []byte(required)) {
			return fmt.Errorf("C2B_V4_SOURCE_REQUIRED_%s", required)
		}
	}
	return nil
}

func verifyMachOArtifacts(fsys fs.FS) error {
	runnerBytes, err := fs.ReadFile(fsys, "capsule-host-runner")
	if err != nil {
		return err
	}
	runner, err := macho.NewFile(bytes.NewReader(runnerBytes))
	if err != nil {
		return fmt.Errorf("C2B_V4_RUNNER_MACHO: %w", err)
	}
	defer func() { _ = runner.Close() }()
	if runner.Cpu != macho.CpuArm64 {
		return errors.New("C2B_V4_RUNNER_ARCH")
	}
	libraries, err := runner.ImportedLibraries()
	if err != nil || !slices.Equal(libraries, []string{"@rpath/libkrun.1.dylib", "/usr/lib/libSystem.B.dylib"}) {
		return errors.New("C2B_V4_RUNNER_LIBRARIES")
	}
	imports, err := runner.ImportedSymbols()
	if err != nil {
		return fmt.Errorf("C2B_V4_RUNNER_IMPORTS: %w", err)
	}
	var krunImports []string
	for _, symbol := range imports {
		if strings.HasPrefix(symbol, "_krun_") {
			krunImports = append(krunImports, symbol)
		}
	}
	slices.Sort(krunImports)
	if !slices.Equal(krunImports, expectedRunnerKrunImports) || hasCodeSignature(runner) {
		return errors.New("C2B_V4_RUNNER_IMPORT_OR_SIGNATURE")
	}
	hasRpath := false
	for _, load := range runner.Loads {
		if rpath, ok := load.(*macho.Rpath); ok && rpath.Path == "@executable_path" {
			hasRpath = true
		}
	}
	if !hasRpath {
		return errors.New("C2B_V4_RUNNER_RPATH")
	}

	libBytes, err := fs.ReadFile(fsys, "libkrun.1.dylib")
	if err != nil {
		return err
	}
	lib, err := macho.NewFile(bytes.NewReader(libBytes))
	if err != nil {
		return fmt.Errorf("C2B_V4_LIBKRUN_MACHO: %w", err)
	}
	defer func() { _ = lib.Close() }()
	if lib.Cpu != macho.CpuArm64 || hasCodeSignature(lib) || lib.Symtab == nil {
		return errors.New("C2B_V4_LIBKRUN_ARCH_SIGNATURE_SYMTAB")
	}
	exports := make(map[string]bool)
	for _, symbol := range lib.Symtab.Syms {
		exports[symbol.Name] = true
	}
	for _, symbol := range expectedRunnerKrunImports {
		if !exports[symbol] {
			return fmt.Errorf("C2B_V4_LIBKRUN_MISSING_EXPORT_%s", symbol)
		}
	}
	return nil
}

func hasCodeSignature(file *macho.File) bool {
	const lcCodeSignature = 0x1d
	for _, load := range file.Loads {
		raw := load.Raw()
		if len(raw) >= 4 && binary.LittleEndian.Uint32(raw[:4]) == lcCodeSignature {
			return true
		}
	}
	return false
}

func readBounded(fsys fs.FS, path string, maximum int) ([]byte, error) {
	file, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	value, err := io.ReadAll(io.LimitReader(file, int64(maximum)+1))
	if err != nil {
		return nil, err
	}
	if len(value) > maximum {
		return nil, errors.New("maximum exceeded")
	}
	return value, nil
}

func digest(value []byte) string {
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}

func requireEOF(decoder *json.Decoder) error {
	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("C2B_V4_PROFILE_TRAILING")
	}
	return nil
}
