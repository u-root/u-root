package acpi

import (
	"encoding/binary"
	"fmt"
)

// Raw is just a table embedded in a []byte.  Operations on Raw are
// useful for unpacking into a more refined table or just figuring out
// how to skip a table you don't care about.
type Raw struct {
	data []byte
}

var _ = Tabler(&Raw{})

func NewRaw(b []byte) (Tabler, error) {
	u := binary.LittleEndian.Uint32(b[LengthOffset : LengthOffset+4])
	return &Raw{data: b[0:u]}, nil
}

func (r *Raw) Len() int {
	return len([]byte(r.data))
}

func (r *Raw) Data() []byte {
	return r.data
}

func (r *Raw) Sig() string {
	return fmt.Sprintf("%s", r.data[:4])
}

func (r *Raw) OEMID() string {
	return fmt.Sprintf("%s", r.data[10:16])
}

func (r *Raw) OEMTableID() string {
	return fmt.Sprintf("%s", r.data[16:24])
}

func (r *Raw) OEMRevision() uint32 {
	u := binary.LittleEndian.Uint32(r.data[LengthOffset : LengthOffset+4])
	return u
}

func (r *Raw) CreatorID() uint32 {
	u := binary.LittleEndian.Uint32(r.data[LengthOffset : LengthOffset+4])
	return u
}

func (r *Raw) VendorID() uint32 {
	u := binary.LittleEndian.Uint32(r.data[LengthOffset : LengthOffset+4])
	return u
}

func (r *Raw) CreatorRevision() uint32 {
	u := binary.LittleEndian.Uint32(r.data[LengthOffset : LengthOffset+4])
	return u
}

func (r *Raw) Revision() uint8 {
	return r.data[8]
}

func (r *Raw) Checksum() uint8 {
	return r.data[9]
}
