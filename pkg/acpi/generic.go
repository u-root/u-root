package acpi

import (
	"encoding/binary"
)

// Generic is the generic ACPI table, with a header and data
type Generic struct {
	Header
	data []byte
}

var _ = Tabler(&Generic{})

func NewGeneric(b []byte) (Tabler, error) {
	t, err := NewRaw(b)
	if err != nil {
		return nil, err
	}
	return &Generic{Header: *GetHeader(t), data: t.AllData()}, nil
}

// Marshal marshals Generic tables. The main use of this function
// is when you want to tweak the header a bit; you can convert a Raw
// table to a Generic table, do what you wish, and write it out.
func (g *Generic) Marshal() ([]byte, error) {
	// Marshal the header, as it may be changed.
	h, err := g.Header.Marshal()
	if err != nil {
		return nil, err
	}
	// Append only the table data.
	h = append(h, g.TableData()...)
	binary.LittleEndian.PutUint32(h[LengthOffset:], uint32(len(h)))
	h[CSUMOffset] = 0
	c := gencsum(h)
	Debug("CSUM is %#x", c)
	h[CSUMOffset] = c
	return h, nil
}
func (r *Generic) Len() uint32 {
	return uint32(len(r.data))
}

func (r *Generic) AllData() []byte {
	return r.data
}

func (r *Generic) TableData() []byte {
	return r.data[36:]
}

func (r *Generic) Sig() sig {
	return r.Header.Sig
}

func (r *Generic) OEMID() oem {
	return r.Header.OEMID
}

func (r *Generic) OEMTableID() tableid {
	return r.Header.OEMTableID
}

func (r *Generic) OEMRevision() uint32 {
	return r.Header.OEMRevision
}

func (r *Generic) CreatorID() uint32 {
	return r.Header.CreatorID
}

func (r *Generic) CreatorRevision() uint32 {
	return r.Header.CreatorRevision
}

func (r *Generic) Revision() uint8 {
	return r.Header.Revision
}

func (r *Generic) CheckSum() uint8 {
	return r.Header.CheckSum
}
