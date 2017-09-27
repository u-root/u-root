package dhcp6

import (
	"encoding/binary"
	"io"
)

// This option is used by a client to identify the vendor that
// manufactured the hardware on which the client is running.  The
// information contained in the data area of this option is contained in
// one or more opaque fields that identify details of the hardware
// configuration.
type VendorClass struct {
	// EnterpriseNumber specifies an IANA-assigned vendor Private Enterprise
	// Number.
	EnterpriseNumber uint32

	// The vendor-class-data is composed of a series of separate items, each
	// of which describes some characteristic of the client's hardware
	// configuration.  Examples of vendor-class-data instances might include
	// the version of the operating system the client is running or the
	// amount of memory installed on the client.
	VendorClassData Data
}

// MarshalBinary allocates a byte slice containing the data from a VendorClass.
func (vc *VendorClass) MarshalBinary() ([]byte, error) {
	// 4 bytes: EnterpriseNumber
	// N bytes: VendorClassData
	vcb, err := vc.VendorClassData.MarshalBinary()
	if err != nil {
		return nil, err
	}

	b := make([]byte, 4, 4+len(vcb))
	binary.BigEndian.PutUint32(b, vc.EnterpriseNumber)
	b = append(b, vcb...)

	return b, nil
}

// UnmarshalBinary unmarshals a raw byte slice into a VendorClass.
//
// If the byte slice is less than 4 bytes in length,
// or if VendorClassData is malformed, io.ErrUnexpectedEOF is
// returned.
func (vc *VendorClass) UnmarshalBinary(b []byte) error {
	if len(b) < 4 {
		return io.ErrUnexpectedEOF
	}

	// Extract EnterpriseNumber
	vc.EnterpriseNumber = binary.BigEndian.Uint32(b[:4])

	return vc.VendorClassData.UnmarshalBinary(b[4:])
}
