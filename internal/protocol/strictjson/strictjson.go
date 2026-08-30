// Package strictjson is the shared "decode trailing-JSON check" and
// "marshal-then-unmarshal deep clone" helper backing the runtimec2bpassive
// family (v1/v2/v3), runtimecompositionpassive, and
// runtimeexecutionprofilepassive contract codecs. Every caller keeps its own
// error-code prefix and its own field-authority Go type; this package only
// supplies the mechanical, type-independent operation each one repeated.
package strictjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// RequireEOF decodes one more JSON value from decoder and reports whether
// the input contains trailing data after the caller's own top-level decode.
// It expects decoder's input to be otherwise exhausted (io.EOF); any other
// outcome — a further decode error, or a further value decoding cleanly — is
// a refusal, wrapped (or, for a genuine trailing value, formatted) with the
// caller-owned prefix so each contract family keeps its own error-code
// vocabulary.
func RequireEOF(decoder *json.Decoder, prefix string) error {
	var trailing any
	if err := decoder.Decode(&trailing); errors.Is(err, io.EOF) {
		return nil
	} else if err != nil {
		return fmt.Errorf("%s: %w", prefix, err)
	}
	return fmt.Errorf("%s: extra JSON value", prefix)
}

// Clone returns a deep, alias-free copy of value via a marshal/unmarshal
// round trip. It panics if value fails to encode or the encoded bytes fail
// to decode back into T; both indicate the caller's own type failed to
// round-trip through its declared JSON shape, not a reachable runtime
// condition for the contract types this package backs.
func Clone[T any](value *T) *T {
	exact, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	var copied T
	if err := json.Unmarshal(exact, &copied); err != nil {
		panic(err)
	}
	return &copied
}
