package sourcevalidatorpassive

import "encoding/binary"

// This file is the parameterized frame codec shared by this package's v0
// (putBodyLength/validateFrame/putCommon/validateCommon in contract.go) and
// v1 (v1PutLengthMagic/v1ValidateFrame/v1PutCommonFrame/v1ValidateRoleTags
// in contract_v1.go) frame-header scaffolding. Both versions frame a value
// as a 4-byte big-endian declared body length, an 8-byte magic, and a run
// of consecutive big-endian uint16 identity tags starting at a
// version-specific offset — the primitives below are that exact shared
// shape. Each version keeps its own field layout, error classification,
// error code, and check order (they differ enough between v0 and v1 — R1
// succeeds R0's fixed 8-tag header with a role-parameterized 11-tag one,
// and classifies frame-shape refusals differently — that forcing one on the
// other would change a refusal code, not just deduplicate scaffolding).

// putFrameHeader writes the shared header shape both versions' fixed
// contract frames start with: a 4-byte big-endian declared body length
// (len(frame)-4) followed by the 8-byte magic.
func putFrameHeader(frame []byte, magic [8]byte) {
	// Every caller allocates one of this package's fixed, reviewed contract
	// frame sizes (the largest is R1's request header plus the M1 source
	// cap), all well within uint32.
	binary.BigEndian.PutUint32(frame[:4], uint32(len(frame)-4)) //nolint:gosec
	copy(frame[4:12], magic[:])
}

// declaredBodyLength reads the frame's 4-byte big-endian declared body
// length. It performs no bounds checking: every caller has already
// established frame is at least 4 bytes before calling this (v0's
// validateFrame checks explicitly; v1's minimum-length check always runs
// first, and every v1 frame minimum is well above 4).
func declaredBodyLength(frame []byte) uint32 {
	return binary.BigEndian.Uint32(frame[:4])
}

// putSequentialTags writes each of values as a consecutive big-endian
// uint16, starting at offset. Both versions' common-frame encoders
// (putCommon's 8 fixed identity tags at offset 12; v1PutCommonFrame's
// version+kind+role+11 role tags at offset 12; encodeV1Profile's and
// encodeV1Consumer's 6-tag role runs at offset 18) are this exact loop over
// a different values slice and starting offset.
func putSequentialTags(frame []byte, offset int, values []uint16) {
	for index, value := range values {
		binary.BigEndian.PutUint16(frame[offset+index*2:offset+2+index*2], value)
	}
}

// mismatchedTagIndex compares each of want against the consecutive
// big-endian uint16 fields starting at offset, returning the index of the
// first mismatch, or -1 if every field matches. Both versions' common-frame
// tag validators (validateCommon's 8-field comparison; v1ValidateRoleTags'
// role-tag comparison) are this exact loop; each keeps its own
// classification/code decision for the mismatched index.
func mismatchedTagIndex(frame []byte, offset int, want []uint16) int {
	for index, value := range want {
		if binary.BigEndian.Uint16(frame[offset+index*2:offset+2+index*2]) != value {
			return index
		}
	}
	return -1
}
