package dhcp6

import (
	"encoding/binary"
	"io"
)

// The Authentication option carries authentication information to
// authenticate the identity and contents of DHCP messages. The use of
// the Authentication option is described in section 21.
type Authentication struct {
	// The authentication protocol used in this authentication option
	Protocol byte

	// The algorithm used in the authentication protocol
	Algorithm byte

	// The replay detection method used in this authentication option
	RDM byte

	// The replay detection information for the RDM
	ReplayDetection uint64

	// The authentication information,
	// as specified by the protocol and
	// algorithm used in this authentication
	// option
	AuthenticationInformation []byte
}

// MarshalBinary allocates a byte slice containing the data from a Authentication.
func (a *Authentication) MarshalBinary() ([]byte, error) {
	// 1 byte:  Protocol
	// 1 byte:  Algorithm
	// 1 byte:  RDM
	// 8 bytes: ReplayDetection
	// N bytes: AuthenticationInformation (can have 0 len byte)
	b := make([]byte, 11+len(a.AuthenticationInformation))
	_ = append(b[:0], a.Protocol, a.Algorithm, a.RDM)
	binary.BigEndian.PutUint64(b[3:11], a.ReplayDetection)
	copy(b[11:], a.AuthenticationInformation)

	return b, nil
}

// UnmarshalBinary unmarshals a raw byte slice into a Authentication.
// If the byte slice does not contain enough data to form a valid
// Authentication, io.ErrUnexpectedEOF is returned.
func (a *Authentication) UnmarshalBinary(b []byte) error {
	// Too short to be valid Authentication
	if len(b) < 11 {
		return io.ErrUnexpectedEOF
	}

	a.Protocol = b[0]
	a.Algorithm = b[1]
	a.RDM = b[2]
	a.ReplayDetection = binary.BigEndian.Uint64(b[3:])
	a.AuthenticationInformation = make([]byte, len(b[11:]))
	copy(a.AuthenticationInformation, b[11:])

	return nil
}
