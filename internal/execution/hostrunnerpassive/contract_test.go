package hostrunnerpassive

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"capsule.local/capsule/internal/execution/hostrunnercontract"
)

func TestRetainedHostRunnerSourceContract(t *testing.T) {
	verified, err := VerifyFixtureFS(fixtureFS())
	if err != nil {
		t.Fatal(err)
	}
	manifest := verified.Manifest()
	if manifest.LaunchAuthority.Cardinality != "exactly-one-runner-process-per-AttemptID" ||
		manifest.LaunchAuthority.ExecuteTimeReplacementValues || manifest.Source.FinalRunnerArtifact ||
		manifest.Source.CallsLibkrun || manifest.Teardown.Owner != "execution-supervisor" {
		t.Fatal("passive host-runner boundary widened")
	}
	source := verified.Source()
	source[0] ^= 0xff
	if bytes.Equal(source, verified.Source()) {
		t.Fatal("source accessor leaked retained storage")
	}
}

func TestSourceSequenceMutationsRefuse(t *testing.T) {
	source, err := fs.ReadFile(fixtureFS(), "capsule-host-runner-contract.c")
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string][]byte{
		"disable-console-removed":     bytes.Replace(source, []byte("krun_disable_implicit_console"), []byte("krun_enable_implicit_console"), 1),
		"port-order-swapped":          bytes.Replace(source, []byte("CAPSULE_LIBKRUN_CALL(12"), []byte("CAPSULE_LIBKRUN_CALL(13"), 1),
		"raw-root-replaced":           bytes.Replace(source, []byte("krun_add_read_only_raw_root_fd_fd4_vda"), []byte("krun_set_root_live_path"), 1),
		"start-before-handshake":      bytes.Replace(source, []byte("CAPSULE_LIBKRUN_CALL(19, krun_start_enter);"), []byte("CAPSULE_LIBKRUN_CALL(18, krun_start_enter);"), 1),
		"execute-replacement-enabled": bytes.Replace(source, []byte("CAPSULE_EXECUTE_TIME_REPLACEMENT_VALUES 0"), []byte("CAPSULE_EXECUTE_TIME_REPLACEMENT_VALUES 1"), 1),
	}
	for name, candidate := range cases {
		t.Run(name, func(t *testing.T) {
			if bytes.Equal(candidate, source) {
				t.Fatal("mutation did not apply")
			}
			if err := validateSourceContract(candidate); err == nil {
				t.Fatal("source mutation accepted")
			}
		})
	}
}

func TestManifestAuthorityMutationsRefuse(t *testing.T) {
	exact, err := fs.ReadFile(fixtureFS(), "manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	manifest, err := decodeManifest(exact)
	if err != nil {
		t.Fatal(err)
	}
	cases := map[string]func(*Manifest){
		// validateIdentity
		"contract-mismatch":        func(value *Manifest) { value.Contract = "other" },
		"status-mismatch":          func(value *Manifest) { value.Status = "other" },
		"source-path-mismatch":     func(value *Manifest) { value.Source.Path = "other.c" },
		"source-language-mismatch": func(value *Manifest) { value.Source.Language = "other" },
		"source-bytes-mismatch":    func(value *Manifest) { value.Source.Bytes = 1 },
		"source-sha256-mismatch":   func(value *Manifest) { value.Source.SHA256 = "deadbeef" },
		"source-calls-libkrun":     func(value *Manifest) { value.Source.CallsLibkrun = true },
		"artifact-claim":           func(value *Manifest) { value.Source.FinalRunnerArtifact = true },

		// validatePredecessor
		"predecessor-bytes-mismatch":           func(value *Manifest) { value.Predecessor.C2BV3Bytes = 1 },
		"predecessor-sha256-mismatch":          func(value *Manifest) { value.Predecessor.C2BV3SHA256 = "x" },
		"predecessor-contract-sha256-mismatch": func(value *Manifest) { value.Predecessor.C2BV3ContractSHA256 = "x" },

		// validateAcceptedLibkrun
		"libkrun-repository-mismatch":      func(value *Manifest) { value.AcceptedLibkrun.Repository = "other/other" },
		"libkrun-upstream-mismatch":        func(value *Manifest) { value.AcceptedLibkrun.Upstream = "x" },
		"libkrun-accepted-commit-mismatch": func(value *Manifest) { value.AcceptedLibkrun.AcceptedCommit = "x" },
		"libkrun-accepted-tree-mismatch":   func(value *Manifest) { value.AcceptedLibkrun.AcceptedTree = "x" },
		"libkrun-header-state-mismatch":    func(value *Manifest) { value.AcceptedLibkrun.HeaderAndDylibState = "x" },

		// validateLaunchAuthority
		"two-runners":                  func(value *Manifest) { value.LaunchAuthority.Cardinality = "two-runners-per-AttemptID" },
		"replacement-values":           func(value *Manifest) { value.LaunchAuthority.ExecuteTimeReplacementValues = true },
		"daemon-route":                 func(value *Manifest) { value.LaunchAuthority.DaemonRoute = true },
		"launch-authority-owner":       func(value *Manifest) { value.LaunchAuthority.Owner = "other" },
		"launch-authority-descriptor":  func(value *Manifest) { value.LaunchAuthority.Descriptor = "other" },
		"launch-authority-accepts-id":  func(value *Manifest) { value.LaunchAuthority.RunnerAcceptsAttemptIDBytes = true },
		"launch-authority-service":     func(value *Manifest) { value.LaunchAuthority.Service = true },
		"launch-authority-priv-helper": func(value *Manifest) { value.LaunchAuthority.PrivilegedHelper = true },

		// validatePreflight
		"fd-eight":                   func(value *Manifest) { value.Preflight.CloseFromInclusive = 9 },
		"writable-root":              func(value *Manifest) { value.Preflight.HostFDs[4].AccessMode = "O_RDWR" },
		"preflight-argv-mismatch":    func(value *Manifest) { value.Preflight.Argv = "other" },
		"preflight-env-mismatch":     func(value *Manifest) { value.Preflight.Environment = "other" },
		"preflight-hostfd-role":      func(value *Manifest) { value.Preflight.HostFDs[0].Role = "other" },
		"preflight-hostfd-fd":        func(value *Manifest) { value.Preflight.HostFDs[0].FD = 99 },
		"preflight-root-mismatch":    func(value *Manifest) { value.Preflight.Root = "other" },
		"preflight-failure-mismatch": func(value *Manifest) { value.Preflight.Failure = "other" },

		// validatePortsAndDevices
		"fourth-port":              func(value *Manifest) { value.Ports = append(value.Ports, hostrunnercontract.Port{PortID: 3}) },
		"port-id-mismatch":         func(value *Manifest) { value.Ports[0].PortID = 9 },
		"port-name-mismatch":       func(value *Manifest) { value.Ports[0].Name = "other" },
		"port-input-fd-mismatch":   func(value *Manifest) { value.Ports[0].InputFD = 99 },
		"port-output-fd-mismatch":  func(value *Manifest) { value.Ports[2].OutputFD = 99 },
		"port-guest-node-mismatch": func(value *Manifest) { value.Ports[0].GuestNode = "/dev/other" },
		"network-device":           func(value *Manifest) { value.Devices = append(value.Devices, "network") },
		"devices-content-mismatch": func(value *Manifest) { value.Devices[0] = "other" },

		// validateForbiddenAuthority
		"forbidden-authority-mismatch": func(value *Manifest) { value.ForbiddenAuthority[0] = "other" },

		// validateTeardown
		"runner-teardown":                  func(value *Manifest) { value.Teardown.Owner = "runner" },
		"teardown-runner-exit":             func(value *Manifest) { value.Teardown.RunnerExit = "other" },
		"teardown-identity-revalidation":   func(value *Manifest) { value.Teardown.IdentityRevalidation = "other" },
		"teardown-forced-signal":           func(value *Manifest) { value.Teardown.ForcedSignal = "other" },
		"teardown-forced-absence-deadline": func(value *Manifest) { value.Teardown.ForcedAbsenceDeadlineMs = 1 },
		"teardown-max-action-to-absence":   func(value *Manifest) { value.Teardown.MaximumActionToAbsenceMs = 1 },
		"teardown-identity-mismatch-field": func(value *Manifest) { value.Teardown.IdentityMismatch = "other" },
		"teardown-authoritative-absence":   func(value *Manifest) { value.Teardown.AuthoritativeAbsenceNeeded = false },

		// validateBlockersAndEffects
		"blockers-mismatch":   func(value *Manifest) { value.Blockers[0] = "other" },
		"guest-effect":        func(value *Manifest) { value.Effects.Guest = true },
		"effect-process":      func(value *Manifest) { value.Effects.Process = true },
		"effect-libkrun":      func(value *Manifest) { value.Effects.Libkrun = true },
		"effect-hvf":          func(value *Manifest) { value.Effects.HVF = true },
		"effect-vm":           func(value *Manifest) { value.Effects.VM = true },
		"effect-signing":      func(value *Manifest) { value.Effects.Signing = true },
		"effect-installation": func(value *Manifest) { value.Effects.Installation = true },
		"effect-admission":    func(value *Manifest) { value.Effects.Admission = true },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			candidate := cloneManifest(manifest)
			mutate(&candidate)
			if err := ValidateManifest(candidate); err == nil {
				t.Fatal("manifest mutation accepted")
			}
		})
	}
}

func TestFixtureBoundsAndKnownAnswerMutationsRefuse(t *testing.T) {
	base := fixtureFS()
	mutated := fstest.MapFS{}
	err := fs.WalkDir(base, ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}
		value, readErr := fs.ReadFile(base, path)
		if readErr != nil {
			return readErr
		}
		mutated[path] = &fstest.MapFile{Data: value}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	source := append([]byte(nil), mutated["capsule-host-runner-contract.c"].Data...)
	source[len(source)-2] ^= 1
	mutated["capsule-host-runner-contract.c"] = &fstest.MapFile{Data: source}
	if _, err := VerifyFixtureFS(mutated); err == nil {
		t.Fatal("one-byte source mutation accepted")
	}
	oversized := fstest.MapFS{"manifest.json": &fstest.MapFile{Data: bytes.Repeat([]byte{' '}, ManifestMaximumBytes+1)}}
	if _, err := VerifyFixtureFS(oversized); err == nil || !strings.Contains(err.Error(), "READ") {
		t.Fatalf("oversized manifest not refused at bounded read: %v", err)
	}
}

func TestNoProductConsumerImportsPassiveHostRunner(t *testing.T) {
	root := os.DirFS(filepath.Join("..", "..", ".."))
	needle := "internal/execution/hostrunnerpassive"
	err := fs.WalkDir(root, ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" || entry.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" || strings.Contains(path, "hostrunnerpassive") {
			return nil
		}
		contents, readErr := fs.ReadFile(root, path)
		if readErr != nil {
			return readErr
		}
		if bytes.Contains(contents, []byte(needle)) {
			t.Errorf("passive host-runner verifier has product consumer: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func fixtureFS() fs.FS {
	return os.DirFS("../../../schemas/conformance/c2b-host-runner-source-v1")
}
