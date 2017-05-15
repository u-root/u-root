package dhcp6

import (
	"encoding/binary"
	"io"
	"time"
)

// IANA represents an Identity Association for Non-temporary Addresses, as
// defined in RFC 3315, Section 22.4.
//
// Multiple IANAs may be present in a single DHCP request.
type IANA struct {
	// IAID specifies a DHCP identity association identifier.  The IAID
	// is a unique, client-generated identifier.
	IAID [4]byte

	// T1 specifies how long a DHCP client will wait to contact this server,
	// to extend the lifetimes of the addresses assigned to this IANA
	// by this server.
	T1 time.Duration

	// T2 specifies how long a DHCP client will wait to contact any server,
	// to extend the lifetimes of the addresses assigned to this IANA
	// by this server.
	T2 time.Duration

	// Options specifies a map of DHCP options specific to this IANA.
	// Its methods can be used to retrieve data from an incoming IANA, or send
	// data with an outgoing IANA.
	Options Options
}

// NewIANA creates a new IANA from an IAID, T1 and T2 durations, and an
// Options map.  If an Options map is not specified, a new one will be
// allocated.
func NewIANA(iaid [4]byte, t1 time.Duration, t2 time.Duration, options Options) *IANA {
	if options == nil {
		options = make(Options)
	}

	return &IANA{
		IAID:    iaid,
		T1:      t1,
		T2:      t2,
		Options: options,
	}
}

// MarshalBinary allocates a byte slice containing the data from a IANA.
func (i *IANA) MarshalBinary() ([]byte, error) {
	// 4 bytes: IAID
	// 4 bytes: T1
	// 4 bytes: T2
	// N bytes: options slice byte count
	opts := i.Options.enumerate()
	b := make([]byte, 12+opts.count())

	copy(b[0:4], i.IAID[:])
	binary.BigEndian.PutUint32(b[4:8], uint32(i.T1/time.Second))
	binary.BigEndian.PutUint32(b[8:12], uint32(i.T2/time.Second))
	opts.write(b[12:])

	return b, nil
}

// UnmarshalBinary unmarshals a raw byte slice into a IANA.
//
// If the byte slice does not contain enough data to form a valid IANA,
// io.ErrUnexpectedEOF is returned.
func (i *IANA) UnmarshalBinary(b []byte) error {
	// IANA must contain at least an IAID, T1, and T2.
	if len(b) < 12 {
		return io.ErrUnexpectedEOF
	}

	iaid := [4]byte{}
	copy(iaid[:], b[0:4])
	i.IAID = iaid

	i.T1 = time.Duration(binary.BigEndian.Uint32(b[4:8])) * time.Second
	i.T2 = time.Duration(binary.BigEndian.Uint32(b[8:12])) * time.Second

	options, err := parseOptions(b[12:])
	if err != nil {
		return err
	}
	i.Options = options

	return nil
}
