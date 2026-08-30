package hostrunnercontract

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
)

// Port describes one virtio-console port endpoint wired between the host
// runner and the guest: its port ID, name, host input/output file
// descriptors (-1 when unused), and the guest-side device node it appears as.
type Port struct {
	PortID    int    `json:"portId"`
	Name      string `json:"name"`
	InputFD   int    `json:"inputFd"`
	OutputFD  int    `json:"outputFd"`
	GuestNode string `json:"guestNode"`
}

// ExpectedPorts returns the three fixed console ports every stage of the C2B
// host-runner pipeline must agree on: capsule.source, capsule.input, and
// capsule.completion. It returns a fresh slice on every call so callers may
// freely hold or compare it without risk of mutating a shared literal.
func ExpectedPorts() []Port {
	return []Port{
		{0, "capsule.source", 5, -1, "/dev/hvc0"},
		{1, "capsule.input", 6, -1, "/dev/vport0p1"},
		{2, "capsule.completion", -1, 7, "/dev/vport0p2"},
	}
}

// ReadBounded reads path from fsys, refusing to return more than maximum
// bytes. It checks and propagates both the read error and the file's Close
// error rather than discarding either.
func ReadBounded(fsys fs.FS, path string, maximum int) ([]byte, error) {
	file, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}
	value, readErr := io.ReadAll(io.LimitReader(file, int64(maximum)+1))
	closeErr := file.Close()
	if readErr != nil {
		return nil, readErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	if len(value) > maximum {
		return nil, errors.New("bounded read exceeded")
	}
	return value, nil
}

// RequireEOF decodes one more value from decoder and refuses unless decoding
// immediately hits end of input. A genuine decode error (malformed trailing
// bytes, not merely their presence) is wrapped with errCode via %w so its
// underlying cause is preserved; a clean decode of unexpected trailing data
// (no underlying error) returns errCode alone.
func RequireEOF(decoder *json.Decoder, errCode string) error {
	var trailing any
	if err := decoder.Decode(&trailing); errors.Is(err, io.EOF) {
		return nil
	} else if err != nil {
		return fmt.Errorf("%s: %w", errCode, err)
	}
	return errors.New(errCode)
}
