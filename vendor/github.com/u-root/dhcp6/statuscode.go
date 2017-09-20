package dhcp6

import (
	"encoding/binary"
	"errors"
)

var (
	// errInvalidStatusCode is returned when a byte slice does not contain
	// enough bytes to parse a valid StatusCode value.
	errInvalidStatusCode = errors.New("not enough bytes for valid StatusCode")
)

// StatusCode represents a Status Code, as defined in RFC 3315, Section 5.4.
// DHCP clients and servers can use status codes to communicate successes
// or failures, and provide additional information using a message to describe
// specific failures.
type StatusCode struct {
	// Code specifies the Status value stored within this StatusCode, such as
	// StatusSuccess, StatusUnspecFail, etc.
	Code Status

	// Message specifies a human-readable message within this StatusCode, useful
	// for providing information about successes or failures.
	Message string
}

// NewStatusCode creates a new StatusCode from an input Status value and a
// string message.
func NewStatusCode(code Status, message string) *StatusCode {
	return &StatusCode{
		Code:    code,
		Message: message,
	}
}

// MarshalBinary allocates a byte slice containing the data from a StatusCode.
func (s *StatusCode) MarshalBinary() ([]byte, error) {
	// 2 bytes: status code
	// N bytes: message
	b := make([]byte, 2+len(s.Message))

	binary.BigEndian.PutUint16(b[0:2], uint16(s.Code))
	copy(b[2:], []byte(s.Message))

	return b, nil
}

// UnmarshalBinary unmarshals a raw byte slice into a StatusCode.
//
// If the byte slice does not contain enough data to form a valid StatusCode,
// errInvalidStatusCode is returned.
func (s *StatusCode) UnmarshalBinary(b []byte) error {
	// Too short to contain valid StatusCode
	if len(b) < 2 {
		return errInvalidStatusCode
	}

	s.Code = Status(binary.BigEndian.Uint16(b[0:2]))
	s.Message = string(b[2:])

	return nil
}
