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

func (r *Raw) AllData() []byte {
	return r.data
}

func (r *Raw) TableData() []byte {
	return r.data[36:]
}

func (r *Raw) Sig() sig {
	return sig(fmt.Sprintf("%s", r.data[:4]))
}

func (r *Raw) OEMID() oem {
	return oem(fmt.Sprintf("%s", r.data[10:16]))
}

func (r *Raw) OEMTableID() tableid {
	return tableid(fmt.Sprintf("%s", r.data[16:24]))
}

func (r *Raw) OEMRevision() u32 {
	u := binary.LittleEndian.Uint32(r.data[LengthOffset : LengthOffset+4])
	return u32(u)
}

func (r *Raw) CreatorID() u32 {
	u := binary.LittleEndian.Uint32(r.data[LengthOffset : LengthOffset+4])
	return u32(u)
}

func (r *Raw) VendorID() u32 {
	u := binary.LittleEndian.Uint32(r.data[LengthOffset : LengthOffset+4])
	return u32(u)
}

func (r *Raw) CreatorRevision() u32 {
	u := binary.LittleEndian.Uint32(r.data[LengthOffset : LengthOffset+4])
	return u32(u)
}

func (r *Raw) Revision() u8 {
	return u8(r.data[8])
}

func (r *Raw) CheckSum() u8 {
	return u8(r.data[9])
}
