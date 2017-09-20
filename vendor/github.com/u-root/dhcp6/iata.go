package dhcp6

import (
	"io"
)

// IATA represents an Identity Association for Temporary Addresses, as
// defined in RFC 3315, Section 22.5.
//
// Multiple IATAs may be present in a single DHCP request.
type IATA struct {
	// IAID specifies a DHCP identity association identifier.  The IAID
	// is a unique, client-generated identifier.
	IAID [4]byte

	// Options specifies a map of DHCP options specific to this IATA.
	// Its methods can be used to retrieve data from an incoming IATA, or send
	// data with an outgoing IATA.
	Options Options
}

// NewIATA creates a new IATA from an IAID and an Options map.  If an Options
// map is not specified, a new one will be allocated.
func NewIATA(iaid [4]byte, options Options) *IATA {
	if options == nil {
		options = make(Options)
	}

	return &IATA{
		IAID:    iaid,
		Options: options,
	}
}

// MarshalBinary allocates a byte slice containing the data from a IATA.
func (i *IATA) MarshalBinary() ([]byte, error) {
	// 4 bytes: IAID
	// N bytes: options slice byte count
	opts := i.Options.enumerate()
	b := make([]byte, 4+opts.count())

	copy(b[0:4], i.IAID[:])
	opts.write(b[4:])

	return b, nil
}

// UnmarshalBinary unmarshals a raw byte slice into a IATA.
//
// If the byte slice does not contain enough data to form a valid IATA,
// io.ErrUnexpectedEOF is returned.
func (i *IATA) UnmarshalBinary(b []byte) error {
	// IATA must contain at least an IAID.
	if len(b) < 4 {
		return io.ErrUnexpectedEOF
	}

	iaid := [4]byte{}
	copy(iaid[:], b[0:4])
	i.IAID = iaid

	options, err := parseOptions(b[4:])
	if err != nil {
		return err
	}
	i.Options = options

	return nil
}
