package hostrunnerpassive

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strconv"
	"strings"

	"capsule.local/capsule/internal/execution/hostrunnercontract"
)

const (
	ManifestMaximumBytes = 32 * 1024
	SourceMaximumBytes   = 16 * 1024
	SourceBytes          = 3996
	SourceSHA256         = "7dbc8e5c21acd776fec653665517b6da1ecfe07278da0444d710c251744978b2"
)

var expectedSequence = []string{
	"preflight:validate_one_supervisor_owned_process_per_attempt_id",
	"preflight:validate_exact_argv0_empty_environment_fds_0_through_7_close_from_8_no_execute_time_replacements",
	"preflight:validate_fd4_unlinked_regular_mode_0400_ordonly_fixed_length_and_digest",
	"libkrun:krun_create_ctx",
	"libkrun:krun_set_vm_config_vcpus_1_ram_mib_256",
	"libkrun:krun_disable_implicit_console",
	"libkrun:krun_disable_implicit_init",
	"libkrun:krun_disable_implicit_vsock",
	"libkrun:krun_add_read_only_raw_root_fd_fd4_vda",
	"libkrun:krun_set_root_disk_remount_dev_vda_ext4_ro_nosuid_nodev",
	"libkrun:krun_add_virtio_console_multiport_require_console_id_0",
	"libkrun:krun_add_console_port_inout_capsule_source_fd5_minus1_require_port_id_0",
	"libkrun:krun_add_console_port_inout_capsule_input_fd6_minus1_require_port_id_1",
	"libkrun:krun_add_console_port_inout_capsule_completion_minus1_fd7_require_port_id_2",
	"libkrun:krun_set_kernel_console_hvc0",
	"libkrun:krun_set_workdir_root",
	"libkrun:krun_set_exec_fixed_init_single_argv_three_entry_environment",
	"handshake:emit_ready_then_require_exact_one_byte_G_followed_by_eof",
	"libkrun:krun_start_enter",
}

var expectedForbiddenCalls = []string{
	"krun_set_kernel", "krun_set_firmware", "krun_add_disk", "krun_set_root",
	"krun_add_virtiofs", "krun_add_virtiofs2", "krun_add_virtiofs3", "krun_add_virtiofs4",
	"krun_add_vsock", "krun_add_vsock_port", "krun_add_vsock_port2", "krun_add_net_unixstream",
	"krun_set_gpu_options", "krun_set_snd_device",
}

var expectedForbiddenAuthority = []string{
	"network", "vsock", "tsi", "virtiofs", "NullFs", "writable-block", "additional-block",
	"host-socket", "live-host-path", "kernel-path", "firmware-path", "legacy-disk-path",
	"root-path", "execute-time-config",
}

var expectedBlockers = []string{
	"accepted-libkrun-header-and-current-source-dylib-not-retained-locally",
	"final-runner-abi-build-bytes-digest-and-call-audit-not-created",
	"new-versioned-materialized-profile-not-created",
	"owned-disposable-guest-authorization-not-granted",
}

type Manifest struct {
	Contract           string                    `json:"contract"`
	Status             string                    `json:"status"`
	Source             Source                    `json:"source"`
	Predecessor        Predecessor               `json:"predecessor"`
	AcceptedLibkrun    AcceptedLibkrun           `json:"acceptedLibkrun"`
	LaunchAuthority    LaunchAuthority           `json:"launchAuthority"`
	Preflight          Preflight                 `json:"preflight"`
	Ports              []hostrunnercontract.Port `json:"ports"`
	Devices            []string                  `json:"devices"`
	ForbiddenAuthority []string                  `json:"forbiddenAuthority"`
	Teardown           Teardown                  `json:"teardown"`
	Blockers           []string                  `json:"blockers"`
	Effects            Effects                   `json:"effects"`
}

type Source struct {
	Path                string `json:"path"`
	Language            string `json:"language"`
	Bytes               int    `json:"bytes"`
	SHA256              string `json:"sha256"`
	FinalRunnerArtifact bool   `json:"finalRunnerArtifact"`
	CallsLibkrun        bool   `json:"callsLibkrun"`
}

type Predecessor struct {
	C2BV3Bytes          int    `json:"c2bV3Bytes"`
	C2BV3SHA256         string `json:"c2bV3Sha256"`
	C2BV3ContractSHA256 string `json:"c2bV3ContractSha256"`
}

type AcceptedLibkrun struct {
	Repository          string `json:"repository"`
	Upstream            string `json:"upstream"`
	AcceptedCommit      string `json:"acceptedCommit"`
	AcceptedTree        string `json:"acceptedTree"`
	HeaderAndDylibState string `json:"headerAndDylibState"`
}

type LaunchAuthority struct {
	Owner                        string `json:"owner"`
	Cardinality                  string `json:"cardinality"`
	Descriptor                   string `json:"descriptor"`
	RunnerAcceptsAttemptIDBytes  bool   `json:"runnerAcceptsAttemptIDBytes"`
	ExecuteTimeReplacementValues bool   `json:"executeTimeReplacementValues"`
	DaemonRoute                  bool   `json:"daemonRoute"`
	Service                      bool   `json:"service"`
	PrivilegedHelper             bool   `json:"privilegedHelper"`
}

type Preflight struct {
	Argv               string   `json:"argv"`
	Environment        string   `json:"environment"`
	HostFDs            []HostFD `json:"hostFds"`
	CloseFromInclusive int      `json:"closeFromInclusive"`
	Root               string   `json:"root"`
	Failure            string   `json:"failure"`
}

type HostFD struct {
	FD         int    `json:"fd"`
	Role       string `json:"role"`
	AccessMode string `json:"accessMode"`
}

type Teardown struct {
	Owner                      string `json:"owner"`
	RunnerExit                 string `json:"runnerExit"`
	IdentityRevalidation       string `json:"identityRevalidation"`
	ForcedSignal               string `json:"forcedSignal"`
	ForcedAbsenceDeadlineMs    int    `json:"forcedAbsenceDeadlineMs"`
	MaximumActionToAbsenceMs   int    `json:"maximumActionToAbsenceMs"`
	IdentityMismatch           string `json:"identityMismatch"`
	AuthoritativeAbsenceNeeded bool   `json:"authoritativeAbsenceRequired"`
}

type Effects struct {
	Process      bool `json:"process"`
	Libkrun      bool `json:"libkrun"`
	HVF          bool `json:"hvf"`
	VM           bool `json:"vm"`
	Guest        bool `json:"guest"`
	Signing      bool `json:"signing"`
	Installation bool `json:"installation"`
	Admission    bool `json:"admission"`
}

type Verified struct {
	manifest Manifest
	source   []byte
}

func (verified *Verified) Manifest() Manifest {
	return cloneManifest(verified.manifest)
}

func (verified *Verified) Source() []byte {
	return bytes.Clone(verified.source)
}

func VerifyFixtureFS(fsys fs.FS) (*Verified, error) {
	manifestBytes, err := hostrunnercontract.ReadBounded(fsys, "manifest.json", ManifestMaximumBytes)
	if err != nil {
		return nil, fmt.Errorf("C2B_RUNNER_MANIFEST_READ: %w", err)
	}
	manifest, err := decodeManifest(manifestBytes)
	if err != nil {
		return nil, err
	}
	source, err := hostrunnercontract.ReadBounded(fsys, manifest.Source.Path, SourceMaximumBytes)
	if err != nil {
		return nil, fmt.Errorf("C2B_RUNNER_SOURCE_READ: %w", err)
	}
	if len(source) != SourceBytes || digest(source) != SourceSHA256 ||
		manifest.Source.Bytes != SourceBytes || manifest.Source.SHA256 != SourceSHA256 {
		return nil, errors.New("C2B_RUNNER_SOURCE_KNOWN_ANSWER")
	}
	if err := validateSourceContract(source); err != nil {
		return nil, err
	}
	return &Verified{manifest: cloneManifest(manifest), source: bytes.Clone(source)}, nil
}

func decodeManifest(exact []byte) (Manifest, error) {
	decoder := json.NewDecoder(bytes.NewReader(exact))
	decoder.DisallowUnknownFields()
	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return Manifest{}, fmt.Errorf("C2B_RUNNER_MANIFEST_DECODE: %w", err)
	}
	if err := hostrunnercontract.RequireEOF(decoder, "C2B_RUNNER_MANIFEST_TRAILING"); err != nil {
		return Manifest{}, err
	}
	if err := ValidateManifest(manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

// ValidateManifest checks every field of a decoded Manifest against the
// frozen passive host-runner known answer. It delegates to one validator per
// manifest section so each refusal reason stays independently testable and
// each function stays small; the sections and their refusal codes are
// unchanged from the original flat check.
func ValidateManifest(manifest Manifest) error {
	for _, validate := range []func(Manifest) error{
		validateIdentity,
		validatePredecessor,
		validateAcceptedLibkrun,
		validateLaunchAuthority,
		validatePreflight,
		validatePortsAndDevices,
		validateForbiddenAuthority,
		validateTeardown,
		validateBlockersAndEffects,
	} {
		if err := validate(manifest); err != nil {
			return err
		}
	}
	return nil
}

func validateIdentity(manifest Manifest) error {
	if manifest.Contract != "capsule.c2b.host-runner-source-contract/v1" ||
		manifest.Status != "PASSED-passive-static-source-contract-only" ||
		manifest.Source != (Source{"capsule-host-runner-contract.c", "C17-passive-source-contract", SourceBytes, SourceSHA256, false, false}) {
		return errors.New("C2B_RUNNER_IDENTITY")
	}
	return nil
}

func validatePredecessor(manifest Manifest) error {
	if manifest.Predecessor != (Predecessor{18357, "d72327bba369484a56db7d543a32e8bbd4eac403230ac65d63709ac3ba3bbdfb", "8b1ec936a7b56370716d28557125e46866dea8f21a149704a01f251a0dddbcc1"}) {
		return errors.New("C2B_RUNNER_PREDECESSOR")
	}
	return nil
}

func validateAcceptedLibkrun(manifest Manifest) error {
	if manifest.AcceptedLibkrun != (AcceptedLibkrun{"Shrimpworks/libkrun", "728df8125077d0db44265f6e997c72b81b65c015", "7432eda5a49220976b0167005aa43ee622f9d632", "7671440cfbafa58fe20aebf8d4deb2a843ebe346", "BLOCKED-not-retained-locally-for-final-build"}) {
		return errors.New("C2B_RUNNER_LIBKRUN_IDENTITY")
	}
	return nil
}

func validateLaunchAuthority(manifest Manifest) error {
	if manifest.LaunchAuthority != (LaunchAuthority{"execution-supervisor", "exactly-one-runner-process-per-AttemptID", "sealed-supervisor-inherited-process-and-fd-manifest-only", false, false, false, false, false}) {
		return errors.New("C2B_RUNNER_LAUNCH_AUTHORITY")
	}
	return nil
}

func validatePreflight(manifest Manifest) error {
	expectedFDs := []HostFD{
		{0, "runner-stdin-null", "O_RDONLY"}, {1, "runner-stdout", "O_WRONLY"},
		{2, "runner-stderr", "O_WRONLY"}, {3, "record-before-start-control", "O_RDONLY"},
		{4, "runtime-root", "O_RDONLY"}, {5, "registered-source-port-input", "O_RDONLY"},
		{6, "approved-inline-input-port-input", "O_RDONLY"}, {7, "completion-result-port-output", "O_WRONLY"},
	}
	if manifest.Preflight.Argv != "exact-argv0-only" || manifest.Preflight.Environment != "empty-after-macos-text-encoding-removal" ||
		!slices.Equal(manifest.Preflight.HostFDs, expectedFDs) || manifest.Preflight.CloseFromInclusive != 8 ||
		manifest.Preflight.Root != "fd4-unlinked-regular-mode-0400-ordonly-fixed-length-and-digest" ||
		manifest.Preflight.Failure != "refuse-before-next-step-and-before-krun_start_enter" {
		return errors.New("C2B_RUNNER_PREFLIGHT")
	}
	return nil
}

func validatePortsAndDevices(manifest Manifest) error {
	if !slices.Equal(manifest.Ports, hostrunnercontract.ExpectedPorts()) ||
		!slices.Equal(manifest.Devices, []string{"balloon", "rng", "console-multiport-with-three-fixed-ports", "block-root-vda-read-only"}) {
		return errors.New("C2B_RUNNER_PORT_DEVICE_CONTRACT")
	}
	return nil
}

func validateForbiddenAuthority(manifest Manifest) error {
	if !slices.Equal(manifest.ForbiddenAuthority, expectedForbiddenAuthority) {
		return errors.New("C2B_RUNNER_FORBIDDEN_AUTHORITY")
	}
	return nil
}

func validateTeardown(manifest Manifest) error {
	if manifest.Teardown != (Teardown{"execution-supervisor", "lifecycle-only-never-workload-success", "pid-start-uid-gid-executable-path-code-identity", "SIGKILL", 1000, 1200, "unresolved-no-signal-no-success", true}) {
		return errors.New("C2B_RUNNER_TEARDOWN")
	}
	return nil
}

func validateBlockersAndEffects(manifest Manifest) error {
	if !slices.Equal(manifest.Blockers, expectedBlockers) || manifest.Effects != (Effects{}) {
		return errors.New("C2B_RUNNER_EFFECT_OR_BLOCKER")
	}
	return nil
}

func validateSourceContract(source []byte) error {
	requiredFragments := []string{
		"#define CAPSULE_HOST_RUNNER_SOURCE_CONTRACT_VERSION 1",
		"#define CAPSULE_RUNNERS_PER_ATTEMPT_ID 1",
		"#define CAPSULE_EXECUTE_TIME_REPLACEMENT_VALUES 0",
		"#define CAPSULE_HOST_CLOSE_FROM_INCLUSIVE 8",
		"#define CAPSULE_ROOT_FD 4",
		"#define CAPSULE_CONSOLE_PORT_COUNT 3",
		`{0, "runner-stdin-null", "O_RDONLY"}`,
		`{1, "runner-stdout", "O_WRONLY"}`,
		`{2, "runner-stderr", "O_WRONLY"}`,
		`{3, "record-before-start-control", "O_RDONLY"}`,
		`{4, "runtime-root", "O_RDONLY"}`,
		`{5, "registered-source-port-input", "O_RDONLY"}`,
		`{6, "approved-inline-input-port-input", "O_RDONLY"}`,
		`{7, "completion-result-port-output", "O_WRONLY"}`,
		`{0, "capsule.source", 5, -1, "/dev/hvc0"}`,
		`{1, "capsule.input", 6, -1, "/dev/vport0p1"}`,
		`{2, "capsule.completion", -1, 7, "/dev/vport0p2"}`,
		`"console-multiport-with-three-fixed-ports"`,
		`"block-root-vda-read-only"`,
	}
	for _, fragment := range requiredFragments {
		if bytes.Count(source, []byte(fragment)) != 1 {
			return fmt.Errorf("C2B_RUNNER_SOURCE_STATIC_CONTRACT: %s", fragment)
		}
	}
	lines := strings.Split(string(source), "\n")
	sequence := make([]string, 0, len(expectedSequence))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		var kind string
		switch {
		case strings.HasPrefix(trimmed, "CAPSULE_PREFLIGHT("):
			kind = "preflight"
		case strings.HasPrefix(trimmed, "CAPSULE_LIBKRUN_CALL("):
			kind = "libkrun"
		case strings.HasPrefix(trimmed, "CAPSULE_HANDSHAKE("):
			kind = "handshake"
		default:
			continue
		}
		if strings.HasPrefix(trimmed, "#define") {
			continue
		}
		open := strings.IndexByte(trimmed, '(')
		close := strings.LastIndex(trimmed, ");")
		if open < 0 || close < 0 || close <= open+1 {
			return errors.New("C2B_RUNNER_SOURCE_SYNTAX")
		}
		parts := strings.Split(trimmed[open+1:close], ",")
		if len(parts) != 2 {
			return errors.New("C2B_RUNNER_SOURCE_SYNTAX")
		}
		order, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil || order != len(sequence)+1 {
			return errors.New("C2B_RUNNER_SOURCE_ORDER")
		}
		sequence = append(sequence, kind+":"+strings.TrimSpace(parts[1]))
	}
	if !slices.Equal(sequence, expectedSequence) {
		return errors.New("C2B_RUNNER_SOURCE_SEQUENCE")
	}
	for _, forbidden := range expectedForbiddenCalls {
		if !bytes.Contains(source, []byte(`"`+forbidden+`"`)) {
			return fmt.Errorf("C2B_RUNNER_FORBIDDEN_CALL_SET: %s", forbidden)
		}
	}
	return nil
}

func digest(value []byte) string {
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}

func cloneManifest(manifest Manifest) Manifest {
	manifest.Preflight.HostFDs = slices.Clone(manifest.Preflight.HostFDs)
	manifest.Ports = slices.Clone(manifest.Ports)
	manifest.Devices = slices.Clone(manifest.Devices)
	manifest.ForbiddenAuthority = slices.Clone(manifest.ForbiddenAuthority)
	manifest.Blockers = slices.Clone(manifest.Blockers)
	return manifest
}
