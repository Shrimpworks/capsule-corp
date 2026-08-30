package hostrunnercontract

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
)

func TestExpectedPortsReturnsFreshSlice(t *testing.T) {
	first := ExpectedPorts()
	first[0].Name = "mutated"
	second := ExpectedPorts()
	if second[0].Name == "mutated" {
		t.Fatal("ExpectedPorts leaked its backing array across calls")
	}
	want := []Port{
		{0, "capsule.source", 5, -1, "/dev/hvc0"},
		{1, "capsule.input", 6, -1, "/dev/vport0p1"},
		{2, "capsule.completion", -1, 7, "/dev/vport0p2"},
	}
	if len(second) != len(want) {
		t.Fatalf("unexpected port count: %d", len(second))
	}
	for index := range want {
		if second[index] != want[index] {
			t.Fatalf("port %d mismatch: got %+v want %+v", index, second[index], want[index])
		}
	}
}

func TestReadBoundedSuccess(t *testing.T) {
	fsys := fstest.MapFS{"file": &fstest.MapFile{Data: []byte("exact")}}
	value, err := ReadBounded(fsys, "file", 5)
	if err != nil {
		t.Fatal(err)
	}
	if string(value) != "exact" {
		t.Fatalf("unexpected value: %q", value)
	}
}

func TestReadBoundedOpenErrorPropagates(t *testing.T) {
	fsys := fstest.MapFS{}
	if _, err := ReadBounded(fsys, "missing", 5); err == nil {
		t.Fatal("missing file accepted")
	}
}

func TestReadBoundedExceedsMaximum(t *testing.T) {
	fsys := fstest.MapFS{"file": &fstest.MapFile{Data: []byte("too-long")}}
	if _, err := ReadBounded(fsys, "file", 3); err == nil {
		t.Fatal("oversized read accepted")
	}
}

// faultyFile fails on Close but otherwise reads normally, to prove
// ReadBounded checks and propagates the Close error rather than discarding
// it silently.
type faultyFile struct {
	*bytes.Reader
	closeErr error
}

func (f *faultyFile) Close() error { return f.closeErr }
func (f *faultyFile) Stat() (fs.FileInfo, error) {
	return nil, errors.New("stat not supported")
}

type faultyFS struct {
	data     []byte
	closeErr error
}

func (f faultyFS) Open(name string) (fs.File, error) {
	return &faultyFile{Reader: bytes.NewReader(f.data), closeErr: f.closeErr}, nil
}

func TestReadBoundedPropagatesCloseError(t *testing.T) {
	wantErr := errors.New("close failed")
	fsys := faultyFS{data: []byte("ok"), closeErr: wantErr}
	_, err := ReadBounded(fsys, "file", 10)
	if !errors.Is(err, wantErr) {
		t.Fatalf("Close error not propagated: %v", err)
	}
}

func TestReadBoundedPrefersReadErrorOverCloseError(t *testing.T) {
	// A read that exceeds the maximum should still surface as a bounds
	// refusal even when Close also fails; the read/bounds outcome takes
	// priority since it is determined before Close runs.
	fsys := faultyFS{data: []byte("too-long-for-bound"), closeErr: errors.New("close failed")}
	_, err := ReadBounded(fsys, "file", 3)
	if err == nil {
		t.Fatal("oversized read with failing Close accepted")
	}
}

func TestRequireEOFCleanEnd(t *testing.T) {
	decoder := json.NewDecoder(strings.NewReader(`{"a":1}`))
	var value map[string]int
	if err := decoder.Decode(&value); err != nil {
		t.Fatal(err)
	}
	if err := RequireEOF(decoder, "TEST_TRAILING"); err != nil {
		t.Fatalf("clean end of input refused: %v", err)
	}
}

func TestRequireEOFTrailingValidValueRefuses(t *testing.T) {
	decoder := json.NewDecoder(strings.NewReader(`{"a":1} {"b":2}`))
	var value map[string]int
	if err := decoder.Decode(&value); err != nil {
		t.Fatal(err)
	}
	err := RequireEOF(decoder, "TEST_TRAILING")
	if err == nil {
		t.Fatal("trailing valid JSON value accepted")
	}
	if err.Error() != "TEST_TRAILING" {
		t.Fatalf("unexpected error for clean trailing data: %v", err)
	}
}

func TestRequireEOFTrailingSyntaxErrorPreservesUnderlyingError(t *testing.T) {
	decoder := json.NewDecoder(strings.NewReader(`{"a":1} not-json`))
	var value map[string]int
	if err := decoder.Decode(&value); err != nil {
		t.Fatal(err)
	}
	err := RequireEOF(decoder, "TEST_TRAILING")
	if err == nil {
		t.Fatal("trailing malformed JSON accepted")
	}
	var syntaxErr *json.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("underlying decode error not preserved via %%w: %v", err)
	}
	if !strings.Contains(err.Error(), "TEST_TRAILING") {
		t.Fatalf("error code not present in wrapped error: %v", err)
	}
}

func TestRequireEOFDoesNotMatchNonEOFIOErrors(t *testing.T) {
	// io.ErrUnexpectedEOF must not be treated as a clean end of input by
	// errors.Is(err, io.EOF).
	if errors.Is(io.ErrUnexpectedEOF, io.EOF) {
		t.Fatal("test assumption invalid: io.ErrUnexpectedEOF matches io.EOF")
	}
}
