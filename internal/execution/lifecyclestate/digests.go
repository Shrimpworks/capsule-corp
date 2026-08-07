package lifecyclestate

import (
	"crypto/sha256"
	"encoding/binary"
	"hash"
)

const digestBytes = 32

// ImmutableBindingDigest identifies one complete immutable lifecycle binding.
type ImmutableBindingDigest [digestBytes]byte

// BackendImplementationDigest identifies the exact passive backend implementation.
type BackendImplementationDigest [digestBytes]byte

// BackendInstanceDigest identifies one bounded opaque backend instance value.
type BackendInstanceDigest [digestBytes]byte

func zeroDigest[T ~[digestBytes]byte](value T) bool { return value == (T{}) }

type digestEncoder struct {
	hash hash.Hash
}

func newDigestEncoder(domain string) *digestEncoder {
	encoder := &digestEncoder{hash: sha256.New()}
	encoder.bytes([]byte(domain))
	return encoder
}

func (encoder *digestEncoder) bytes(value []byte) {
	var length [8]byte
	binary.BigEndian.PutUint64(length[:], uint64(len(value)))
	_, _ = encoder.hash.Write(length[:])
	_, _ = encoder.hash.Write(value)
}

func (encoder *digestEncoder) text(value string) { encoder.bytes([]byte(value)) }

func (encoder *digestEncoder) uint64(value uint64) {
	var encoded [8]byte
	binary.BigEndian.PutUint64(encoded[:], value)
	_, _ = encoder.hash.Write(encoded[:])
}

func (encoder *digestEncoder) boolean(value bool) {
	if value {
		encoder.bytes([]byte{1})
		return
	}
	encoder.bytes([]byte{0})
}

func digestSum[T ~[digestBytes]byte](encoder *digestEncoder) T {
	var result T
	copy(result[:], encoder.hash.Sum(nil))
	return result
}
