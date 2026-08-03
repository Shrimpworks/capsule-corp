package archivestate

import (
	"crypto/sha256"
	"encoding/binary"
	"hash"
)

const digestBytes = 32

type ArchiveRecordDigest [digestBytes]byte
type ArchiveCohortDigest [digestBytes]byte
type ArchivePlanDigest [digestBytes]byte
type ArchiveSegmentDigest [digestBytes]byte
type ArchiveCheckpointDigest [digestBytes]byte
type ArchiveIndexDigest [digestBytes]byte
type ArchiveDescriptorSetDigest [digestBytes]byte
type ArchiveCombinedIndexDigest [digestBytes]byte
type ArchiveEffectSeedDigest [digestBytes]byte

type digestEncoder struct{ hash hash.Hash }

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

func zeroDigest[T ~[digestBytes]byte](value T) bool { return value == (T{}) }

// RecordKind prevents exact bytes from crossing registration, approval,
// attempt, and lifecycle digest domains.
type RecordKind string

const (
	RecordRegistration RecordKind = "registration"
	RecordApproval     RecordKind = "approval"
	RecordAttempt      RecordKind = "attempt"
	RecordLifecycle    RecordKind = "lifecycle"
)

func validRecordKind(kind RecordKind) bool {
	switch kind {
	case RecordRegistration, RecordApproval, RecordAttempt, RecordLifecycle:
		return true
	default:
		return false
	}
}

// DigestRecord hashes copied exact record bytes in a closed semantic domain.
func DigestRecord(kind RecordKind, exactBytes []byte) (ArchiveRecordDigest, error) {
	if !validRecordKind(kind) {
		return ArchiveRecordDigest{}, classified(ClassificationDomain, "record-kind")
	}
	if len(exactBytes) == 0 || int64(len(exactBytes)) > MaxSupervisorArchiveBytes {
		return ArchiveRecordDigest{}, classified(ClassificationSchema, "record-bytes")
	}
	encoder := newDigestEncoder("capsule.supervisor.archive-record.v0")
	encoder.text(string(kind))
	encoder.bytes(exactBytes)
	return digestSum[ArchiveRecordDigest](encoder), nil
}
