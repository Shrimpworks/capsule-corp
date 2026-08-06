package hostrunnermaterialized

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
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
