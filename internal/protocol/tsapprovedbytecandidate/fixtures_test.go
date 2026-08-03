package tsapprovedbytecandidate

import (
	"bytes"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"
)

func TestVerifyRetainedFixtureSet(t *testing.T) {
	set, err := VerifyFixtureFS(os.DirFS("../../../schemas/conformance/typescript-approved-byte-v0"))
	if err != nil {
		t.Fatalf("verify retained fixture set: %v", err)
	}
	original := set.Bytes("files/ordinary.ts")
	if len(original) != 391 {
		t.Fatalf("unexpected original byte length: %d", len(original))
	}
	original[0] ^= 0xff
	if bytes.Equal(original, set.Bytes("files/ordinary.ts")) {
		t.Fatal("fixture accessor leaked mutable retained storage")
	}
	bindings := set.Bindings()
	if bindings.OriginalManifest == OriginalManifestDigest(bindings.ExecutableManifest) ||
		bindings.OriginalManifest == OriginalManifestDigest(bindings.TransformationRecordSet) {
		t.Fatal("nominal source roles unexpectedly shared one retained digest")
	}
}

func TestVerifyFixtureSetRejectsByteMutation(t *testing.T) {
	base := os.DirFS("../../../schemas/conformance/typescript-approved-byte-v0")
	mutated := fstest.MapFS{}
	err := fs.WalkDir(base, ".", func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() {
			return walkErr
		}
		value, err := fs.ReadFile(base, path)
		if err != nil {
			return err
		}
		mutated[path] = &fstest.MapFile{Data: value}
		return nil
	})
	if err != nil {
		t.Fatalf("copy fixture FS: %v", err)
	}
	value := append([]byte(nil), mutated["objects/transformation-record.cbor"].Data...)
	value[len(value)-1] ^= 1
	mutated["objects/transformation-record.cbor"] = &fstest.MapFile{Data: value}
	if _, err := VerifyFixtureFS(mutated); err == nil {
		t.Fatal("expected transformation-record mutation to refuse")
	}
}

func TestVerifyFixtureSetBoundsManifestBeforeDecode(t *testing.T) {
	fixture := fstest.MapFS{
		"manifest.json": &fstest.MapFile{Data: bytes.Repeat([]byte{' '}, ManifestMaxBytes+1)},
	}
	if _, err := VerifyFixtureFS(fixture); err == nil {
		t.Fatal("expected oversized manifest to refuse")
	}
}
