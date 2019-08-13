// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

// These types are used in emitting binary from a JSON.
// You can serialize from JSON into structs with these types, and then
// they can be used in a type switch to serialize out different ways.
// In one case, the sheap, they serialize into the head and the heap.
type (
	// These are serialized as:
	sig      string
	oem      string
	tableid  string
	ipaddr   string // 16 byte IP4 or IP6 addr
	sockaddr string // IP addr as above and 2 byte port
	flag     string // with other flags in the struct as one u32
	mac      string // Ethernet MAC to get converted to six uint8's
	bdf      string // u32 pci bus/dev/function
	sheap    string // string placed into the heap, with u16 len and offset in header
	u8       string // 1 byte unsigned
	u16      string // 2 byte unsigned
	u32      string // 4 byte unsigned
	u64      string // 8 byte unsigned
)

// Tabler is the interface to ACPI tables, be they
// held in memory as a byte slice, header and byte slice,
// or more complex struct.
type Tabler interface {
	Sig() string
	Len() uint32
	Revision() uint8
	CheckSum() uint8
	OEMID() string
	OEMTableID() string
	OEMRevision() uint32
	CreatorID() uint32
	CreatorRevision() uint32
	AllData() []byte
	TableData() []byte
	Marshal() ([]byte, error)
}

// Header is the standard header for all ACPI tables, except the
// ones that don't use it. (That's a joke. So is ACPI.)
// We use types that we hope are easy to read; they in turn
// make writing marshal code with type switches very convenient.
type Header struct {
	Sig             sig
	Length          uint32
	Revision        uint8
	CheckSum        uint8
	OEMID           oem
	OEMTableID      tableid
	OEMRevision     uint32
	CreatorID       uint32
	CreatorRevision uint32
}
