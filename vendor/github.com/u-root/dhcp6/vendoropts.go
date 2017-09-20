package dhcp6

import (
	"encoding/binary"
	"io"
)

// A VendorOpts is used by clients and servers to exchange
// VendorOpts information.
type VendorOpts struct {
	// EnterpriseNumber specifies an IANA-assigned vendor Private Enterprise
	// Number.
	EnterpriseNumber uint32

	// An opaque object of option-len octets,
	// interpreted by vendor-specific code on the
	// clients and servers
	Options Options
}

// MarshalBinary allocates a byte slice containing the data from a VendorOpts.
func (v *VendorOpts) MarshalBinary() ([]byte, error) {
	// 4 bytes: EnterpriseNumber
	// N bytes: options slice byte count
	opts := v.Options.enumerate()
	b := make([]byte, 4+opts.count())
	binary.BigEndian.PutUint32(b, v.EnterpriseNumber)
	opts.write(b[4:])

	return b, nil
}

// UnmarshalBinary unmarshals a raw byte slice into a VendorOpts.
// If the byte slice does not contain enough data to form a valid
// VendorOpts, io.ErrUnexpectedEOF is returned.
// If option-data are invalid, then ErrInvalidPacket is returned.
func (v *VendorOpts) UnmarshalBinary(b []byte) error {
	// Too short to be valid VendorOpts
	if len(b) < 4 {
		return io.ErrUnexpectedEOF
	}

	options, err := parseOptions(b[4:])
	if err != nil {
		// Invalid options means an invalid RelayMessage
		return ErrInvalidPacket
	}

	v.EnterpriseNumber = binary.BigEndian.Uint32(b[:4])
	v.Options = options

	return nil
}
