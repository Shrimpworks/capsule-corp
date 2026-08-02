package registrationstate

import (
	"encoding/binary"

	"capsule.local/capsule/internal/protocol/v0candidate"
)

func encodePlanRegistration(view v0candidate.PlanRegistration) []byte {
	result := make([]byte, 0, 192)
	result = appendCBORArgument(result, 5, 10)
	result = appendUnsigned(result, 1)
	result = appendText(result, view.ObjectType)
	result = appendUnsigned(result, 2)
	result = appendUnsigned(result, uint64(view.ObjectVersion))
	result = appendUnsigned(result, 3)
	result = appendByteString(result, view.RegistrationID[:])
	result = appendUnsigned(result, 4)
	result = appendUnsigned(result, uint64(view.RegistrationSequence))
	result = appendUnsigned(result, 5)
	result = appendByteString(result, view.PlanDigest[:])
	result = appendUnsigned(result, 6)
	result = appendByteString(result, view.InstallationID[:])
	result = appendUnsigned(result, 7)
	result = appendUnsigned(result, uint64(view.EpochSequence))
	result = appendUnsigned(result, 8)
	result = appendByteString(result, view.EpochDigest[:])
	result = appendUnsigned(result, 9)
	result = appendByteString(result, view.SupervisorID[:])
	result = appendUnsigned(result, 10)
	result = appendUnsigned(result, uint64(view.ExpiresAt))
	return result
}

func appendText(destination []byte, value string) []byte {
	destination = appendCBORArgument(destination, 3, uint64(len(value)))
	return append(destination, value...)
}

func appendByteString(destination []byte, value []byte) []byte {
	destination = appendCBORArgument(destination, 2, uint64(len(value)))
	return append(destination, value...)
}

func appendUnsigned(destination []byte, value uint64) []byte {
	return appendCBORArgument(destination, 0, value)
}

func appendCBORArgument(destination []byte, major byte, value uint64) []byte {
	switch {
	case value < 24:
		return append(destination, major<<5|byte(value))
	case value <= 0xff:
		return append(destination, major<<5|24, byte(value))
	case value <= 0xffff:
		destination = append(destination, major<<5|25, 0, 0)
		binary.BigEndian.PutUint16(destination[len(destination)-2:], uint16(value))
		return destination
	case value <= 0xffffffff:
		destination = append(destination, major<<5|26, 0, 0, 0, 0)
		binary.BigEndian.PutUint32(destination[len(destination)-4:], uint32(value))
		return destination
	default:
		destination = append(destination, major<<5|27, 0, 0, 0, 0, 0, 0, 0, 0)
		binary.BigEndian.PutUint64(destination[len(destination)-8:], value)
		return destination
	}
}
