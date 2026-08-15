package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepositoryRootFindsPlantedGoModFromNestedDirectory(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.test/root\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "internal", "installation", "macosplan")
	if err := os.MkdirAll(nested, 0o750); err != nil {
		t.Fatal(err)
	}

	original, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(nested); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(original); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	want, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	got, err := repositoryRoot()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("repositoryRoot() = %q, want %q", got, want)
	}
}
