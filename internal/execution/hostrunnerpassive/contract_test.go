package hostrunnerpassive

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
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
		"two-runners":        func(value *Manifest) { value.LaunchAuthority.Cardinality = "two-runners-per-AttemptID" },
		"replacement-values": func(value *Manifest) { value.LaunchAuthority.ExecuteTimeReplacementValues = true },
		"daemon-route":       func(value *Manifest) { value.LaunchAuthority.DaemonRoute = true },
		"fd-eight":           func(value *Manifest) { value.Preflight.CloseFromInclusive = 9 },
		"writable-root":      func(value *Manifest) { value.Preflight.HostFDs[4].AccessMode = "O_RDWR" },
		"fourth-port":        func(value *Manifest) { value.Ports = append(value.Ports, Port{PortID: 3}) },
		"network-device":     func(value *Manifest) { value.Devices = append(value.Devices, "network") },
		"runner-teardown":    func(value *Manifest) { value.Teardown.Owner = "runner" },
		"guest-effect":       func(value *Manifest) { value.Effects.Guest = true },
		"artifact-claim":     func(value *Manifest) { value.Source.FinalRunnerArtifact = true },
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
