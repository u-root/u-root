package acpi

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

func (r *Generic) Len() int {
	return len(r.data)
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

func (r *Generic) OEMRevision() u32 {
	return r.Header.OEMRevision
}

func (r *Generic) CreatorID() u32 {
	return r.Header.CreatorID
}

func (r *Generic) CreatorRevision() u32 {
	return r.Header.CreatorRevision
}

func (r *Generic) Revision() u8 {
	return r.Header.Revision
}

func (r *Generic) CheckSum() u8 {
	return r.Header.CheckSum
}
