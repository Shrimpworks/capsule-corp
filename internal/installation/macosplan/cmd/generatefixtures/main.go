// Command generatefixtures regenerates the passive macOS installation I0
// conformance fixtures under macosplan.GenerateFixtures and writes them to
// their repository-relative paths. It is a local development/CI generator,
// not a daemon, service, or installed component; it holds no bootstrap,
// signing, or protected-root authority.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"capsule.local/capsule/internal/installation/macosplan"
)

func main() {
	root, err := repositoryRoot()
	if err != nil {
		fatal(err)
	}
	fixtures, err := macosplan.GenerateFixtures()
	if err != nil {
		fatal(err)
	}
	for _, fixture := range fixtures {
		path := filepath.Join(root, filepath.FromSlash(fixture.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			fatal(err)
		}
		if err := os.WriteFile(path, fixture.Bytes, 0o600); err != nil {
			fatal(err)
		}
	}
	fmt.Printf("generated %d passive macOS installation I0 fixtures\n", len(fixtures))
}

func repositoryRoot() (string, error) {
	directory, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(directory, "go.mod")); err == nil {
			return directory, nil
		}
		parent := filepath.Dir(directory)
		if parent == directory {
			return "", fmt.Errorf("repository root containing go.mod not found")
		}
		directory = parent
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
