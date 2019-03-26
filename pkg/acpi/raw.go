package acpi

import (
	"encoding/binary"
	"fmt"

	"github.com/u-root/u-root/pkg/io"
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

// ReadRaw reads a full table in, given an address.
func ReadRaw(a int64) (Tabler, error) {
	var u io.Uint32
	// Read the table size at a+4
	if err := io.Read(a+4, &u); err != nil {
		return nil, err
	}
	dat := make([]byte, u)
	for i := range dat {
		var d io.Uint8
		if err := io.Read(a+int64(i), &d); err != nil {
			return nil, err
		}
		dat[i] = uint8(d)
	}
	return &Raw{data: dat}, nil
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
