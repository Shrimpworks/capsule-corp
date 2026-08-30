package hostrunnermaterialized

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"capsule.local/capsule/internal/execution/hostrunnercontract"
)

func TestRetainedMaterializedProfile(t *testing.T) {
	verified, err := VerifyFixtureFS(fixtureFS())
	if err != nil {
		t.Fatal(err)
	}
	profile := verified.Profile()
	if profile.Runner.Authority.Owner != "execution-supervisor" ||
		profile.Runner.Authority.AcceptsReplacementConfiguration ||
		profile.BootRole.RunnerCallsKernelOrFirmwareSetter || profile.ComposedProfile.Admitted ||
		profile.NextGuestGate.GuestAuthorization || profile.Effects.RunnerExecuted ||
		profile.Effects.LibkrunLoaded || profile.Effects.HVFCalled || profile.Effects.GuestCreated {
		t.Fatal("materialized build/static boundary widened")
	}
	profile.Runner.Ports[0].Name = "mutated"
	if verified.Profile().Runner.Ports[0].Name == "mutated" {
		t.Fatal("profile accessor leaked retained slice")
	}
	profile.Predecessors[0] = 'X'
	profile.BootRole.Libkrunfw[0] = 'X'
	profile.Runner.Source[0] = 'X'
	profile.ComposedProfile.Runtime[0] = 'X'
	again := verified.Profile()
	if again.Predecessors[0] == 'X' || again.BootRole.Libkrunfw[0] == 'X' ||
		again.Runner.Source[0] == 'X' || again.ComposedProfile.Runtime[0] == 'X' {
		t.Fatal("profile accessor leaked aliased json.RawMessage bytes")
	}
}

func TestKnownAnswerMutationsRefuse(t *testing.T) {
	base := fixtureFS()
	for _, path := range []string{
		"materialized-profile.json", "libkrun.h", "libkrun-abi-audit.c",
		"libkrun.1.dylib", "capsule-host-runner.c", "capsule-host-runner",
	} {
		t.Run(path, func(t *testing.T) {
			mutated := cloneFS(t, base)
			value := append([]byte(nil), mutated[path].Data...)
			value[len(value)/2] ^= 1
			mutated[path] = &fstest.MapFile{Data: value}
			if _, err := VerifyFixtureFS(mutated); err == nil {
				t.Fatal("one-byte artifact mutation accepted")
			}
		})
	}
}

func TestAuthorityAndCallMutationsRefuse(t *testing.T) {
	base := fixtureFS()
	cases := map[string]func(fstest.MapFS){
		"two-runners": func(candidate fstest.MapFS) {
			mutateText(t, candidate, "materialized-profile.json", "exactly-one-runner-per-committed-AttemptID", "two-runners")
		},
		"guest-authorized": func(candidate fstest.MapFS) {
			mutateText(t, candidate, "materialized-profile.json", `"guestAuthorization": false`, `"guestAuthorization": true`)
		},
		"fd-nine": func(candidate fstest.MapFS) {
			mutateText(t, candidate, "materialized-profile.json", `"allowedFds": [0, 1, 2, 3, 4, 5, 6, 7]`, `"allowedFds": [0, 1, 2, 3, 4, 5, 6, 7, 8]`)
		},
		"firmware-call": func(candidate fstest.MapFS) {
			mutateText(t, candidate, "capsule-host-runner.c", "krun_set_kernel_console(context, \"hvc0\")", "krun_set_firmware(context, \"/tmp/fw\")")
		},
		"port-swap": func(candidate fstest.MapFS) {
			mutateText(t, candidate, "capsule-host-runner.c", "\"capsule.source\", 5, -1", "\"capsule.source\", 6, -1")
		},
		"start-before-handshake": func(candidate fstest.MapFS) {
			mutateText(t, candidate, "capsule-host-runner.c", "require_start_authorization();", "krun_start_enter(context);")
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			candidate := cloneFS(t, base)
			mutate(candidate)
			if _, err := VerifyFixtureFS(candidate); err == nil {
				t.Fatal("authority/call mutation accepted")
			}
		})
	}

	profileBytes, err := fs.ReadFile(base, "materialized-profile.json")
	if err != nil {
		t.Fatal(err)
	}
	profile, err := decodeProfile(profileBytes)
	if err != nil {
		t.Fatal(err)
	}
	fieldCases := map[string]func(*Profile){
		// validateIdentity
		"identity-mismatch": func(value *Profile) { value.Identity = "other" },
		"status-mismatch":   func(value *Profile) { value.Status = "other" },

		// validateBootRole
		"boot-role-mode":              func(value *Profile) { value.BootRole.Mode = "efi" },
		"boot-role-features":          func(value *Profile) { value.BootRole.LibkrunFeatures = []string{"other"} },
		"boot-role-separate-firmware": func(value *Profile) { value.BootRole.SeparateFirmware = "other" },
		"boot-role-calls-setter":      func(value *Profile) { value.BootRole.RunnerCallsKernelOrFirmwareSetter = true },

		// validateRunnerAuthority
		"authority-owner":               func(value *Profile) { value.Runner.Authority.Owner = "other" },
		"authority-cardinality":         func(value *Profile) { value.Runner.Authority.Cardinality = "other" },
		"authority-accepts-attempt-id":  func(value *Profile) { value.Runner.Authority.AcceptsAttemptIDBytes = true },
		"authority-accepts-replacement": func(value *Profile) { value.Runner.Authority.AcceptsReplacementConfiguration = true },
		"authority-daemon-route":        func(value *Profile) { value.Runner.Authority.DaemonRoute = true },
		"authority-service":             func(value *Profile) { value.Runner.Authority.Service = true },
		"authority-privileged-helper":   func(value *Profile) { value.Runner.Authority.PrivilegedHelper = true },
		"authority-teardown-owner":      func(value *Profile) { value.Runner.Authority.TeardownOwner = "other" },

		// validatePreflight
		"preflight-argv":                func(value *Profile) { value.Runner.Preflight.Argv = []string{"other"} },
		"preflight-environment":         func(value *Profile) { value.Runner.Preflight.Environment = "other" },
		"preflight-allowed-fds":         func(value *Profile) { value.Runner.Preflight.AllowedFDs = []int{0, 1, 2, 3, 4, 5, 6} },
		"preflight-close-from":          func(value *Profile) { value.Runner.Preflight.CloseFromInclusive = 9 },
		"preflight-root-fd":             func(value *Profile) { value.Runner.Preflight.RootFD = 5 },
		"preflight-root-bytes":          func(value *Profile) { value.Runner.Preflight.RootBytes = 1 },
		"preflight-root-sha256":         func(value *Profile) { value.Runner.Preflight.RootSHA256 = "x" },
		"preflight-ready-byte":          func(value *Profile) { value.Runner.Preflight.ReadyByte = "X" },
		"preflight-start-authorization": func(value *Profile) { value.Runner.Preflight.StartAuthorization = "other" },

		// validatePortsAndCalls
		"port-id-mismatch":         func(value *Profile) { value.Runner.Ports[0].PortID = 9 },
		"port-name-mismatch":       func(value *Profile) { value.Runner.Ports[0].Name = "other" },
		"port-input-fd-mismatch":   func(value *Profile) { value.Runner.Ports[0].InputFD = 99 },
		"port-output-fd-mismatch":  func(value *Profile) { value.Runner.Ports[2].OutputFD = 99 },
		"port-guest-node-mismatch": func(value *Profile) { value.Runner.Ports[0].GuestNode = "/dev/other" },
		"fourth-port": func(value *Profile) {
			value.Runner.Ports = append(value.Runner.Ports, hostrunnercontract.Port{PortID: 3})
		},
		"ordered-calls-length": func(value *Profile) {
			value.Runner.OrderedCalls = value.Runner.OrderedCalls[:len(value.Runner.OrderedCalls)-1]
		},

		// validateComposedProfileAndStatus
		"composed-digest-mismatch":      func(value *Profile) { value.ComposedProfile.DigestSHA256 = "x" },
		"composed-admitted":             func(value *Profile) { value.ComposedProfile.Admitted = true },
		"work-status-materialization":   func(value *Profile) { value.WorkStatus.MaterializationSlice = "other" },
		"work-status-fixed-owned-guest": func(value *Profile) { value.WorkStatus.FixedOwnedGuest = "other" },
		"work-status-runtime-admission": func(value *Profile) { value.WorkStatus.RuntimeProfileAdmission = "other" },
		"work-status-parent-runtime":    func(value *Profile) { value.WorkStatus.ParentGovernedRuntime = "other" },
		"work-status-runtime-001":       func(value *Profile) { value.WorkStatus.Runtime001 = "other" },
		"work-status-vmm-001":           func(value *Profile) { value.WorkStatus.VMM001 = "other" },

		// validateEffectsAndGuestGate
		"guest-gate-state":                func(value *Profile) { value.NextGuestGate.State = "other" },
		"guest-gate-authorization":        func(value *Profile) { value.NextGuestGate.GuestAuthorization = true },
		"guest-gate-consumer-implemented": func(value *Profile) { value.NextGuestGate.ConsumerImplemented = true },
		"effect-controlled-compilation":   func(value *Profile) { value.Effects.ControlledCompilation = false },
		"effect-runner-executed":          func(value *Profile) { value.Effects.RunnerExecuted = true },
		"effect-libkrun-loaded":           func(value *Profile) { value.Effects.LibkrunLoaded = true },
		"effect-hvf-called":               func(value *Profile) { value.Effects.HVFCalled = true },
		"effect-vm-created":               func(value *Profile) { value.Effects.VMCreated = true },
		"effect-guest-created":            func(value *Profile) { value.Effects.GuestCreated = true },
		"effect-arbitrary-workload":       func(value *Profile) { value.Effects.ArbitraryWorkload = true },
		"effect-credential-used":          func(value *Profile) { value.Effects.CredentialUsed = true },
		"effect-signed":                   func(value *Profile) { value.Effects.Signed = true },
		"effect-installed":                func(value *Profile) { value.Effects.Installed = true },
		"effect-release-published":        func(value *Profile) { value.Effects.ReleasePublished = true },
		"effect-product-wired":            func(value *Profile) { value.Effects.ProductWired = true },
		"effect-admission-changed":        func(value *Profile) { value.Effects.AdmissionChanged = true },
	}
	for name, mutate := range fieldCases {
		t.Run(name, func(t *testing.T) {
			candidate := cloneProfileForTest(profile)
			mutate(&candidate)
			if err := ValidateProfile(candidate); err == nil {
				t.Fatal("profile field mutation accepted")
			}
		})
	}
}

// cloneProfileForTest deep-copies the slice fields TestAuthorityAndCallMutationsRefuse
// mutates so each subtest's candidate does not alias another subtest's backing array.
func cloneProfileForTest(profile Profile) Profile {
	clone := profile
	clone.BootRole.LibkrunFeatures = slices.Clone(profile.BootRole.LibkrunFeatures)
	clone.Runner.Preflight.Argv = slices.Clone(profile.Runner.Preflight.Argv)
	clone.Runner.Preflight.AllowedFDs = slices.Clone(profile.Runner.Preflight.AllowedFDs)
	clone.Runner.Ports = slices.Clone(profile.Runner.Ports)
	clone.Runner.OrderedCalls = slices.Clone(profile.Runner.OrderedCalls)
	clone.Runner.ForbiddenImports = slices.Clone(profile.Runner.ForbiddenImports)
	clone.ComposedProfile.Devices = slices.Clone(profile.ComposedProfile.Devices)
	return clone
}

func TestBoundsRefuse(t *testing.T) {
	oversized := fstest.MapFS{
		"materialized-profile.json": &fstest.MapFile{Data: bytes.Repeat([]byte{' '}, ProfileMaximumBytes+1)},
	}
	if _, err := VerifyFixtureFS(oversized); err == nil || !strings.Contains(err.Error(), "READ") {
		t.Fatalf("oversized profile not refused at bounded read: %v", err)
	}
}

func TestNoProductConsumerImportsMaterializedVerifier(t *testing.T) {
	root := os.DirFS(filepath.Join("..", "..", ".."))
	needle := "internal/execution/hostrunnermaterialized"
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
		if filepath.Ext(path) != ".go" || strings.Contains(path, "hostrunnermaterialized") {
			return nil
		}
		contents, readErr := fs.ReadFile(root, path)
		if readErr != nil {
			return readErr
		}
		if bytes.Contains(contents, []byte(needle)) {
			t.Errorf("materialized verifier has product consumer: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func cloneFS(t *testing.T, source fs.FS) fstest.MapFS {
	t.Helper()
	result := fstest.MapFS{}
	err := fs.WalkDir(source, ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}
		value, readErr := fs.ReadFile(source, path)
		if readErr != nil {
			return readErr
		}
		result[path] = &fstest.MapFile{Data: value}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func mutateText(t *testing.T, candidate fstest.MapFS, path, before, after string) {
	t.Helper()
	value := candidate[path].Data
	mutated := bytes.Replace(value, []byte(before), []byte(after), 1)
	if bytes.Equal(mutated, value) {
		t.Fatalf("mutation did not apply: %s", path)
	}
	candidate[path] = &fstest.MapFile{Data: mutated}
}

func fixtureFS() fs.FS {
	return os.DirFS("../../../schemas/conformance/c2b-host-runner-materialized-v4")
}
