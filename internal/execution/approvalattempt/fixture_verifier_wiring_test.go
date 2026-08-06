package approvalattempt

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestFixtureVerifierHasNoNonTestGoConsumer is the durable repository guard
// that keeps the byte-equality FixtureVerifier out of product authority
// wiring. The definition remains available to cross-package conformance tests,
// but any use from a non-test Go file fails the required repository gate.
func TestFixtureVerifierHasNoNonTestGoConsumer(t *testing.T) {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve fixture-verifier guard path")
	}
	repositoryRoot := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "../../.."))
	definitions := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "fixture_verifier.go"))

	err := filepath.WalkDir(repositoryRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if path != repositoryRoot && (entry.Name() == ".git" || entry.Name() == "node_modules") {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") || filepath.Clean(path) == definitions {
			return nil
		}
		parsed, parseErr := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if parseErr != nil {
			return parseErr
		}
		var forbidden string
		ast.Inspect(parsed, func(node ast.Node) bool {
			identifier, isIdentifier := node.(*ast.Ident)
			if isIdentifier && (identifier.Name == "FixtureVerifier" || identifier.Name == "NewFixtureVerifier") {
				forbidden = identifier.Name
				return false
			}
			return forbidden == ""
		})
		if forbidden != "" {
			relative, _ := filepath.Rel(repositoryRoot, path)
			t.Errorf("test-only %s is referenced by non-test Go file %s", forbidden, relative)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
