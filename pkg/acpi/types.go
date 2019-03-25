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

// Table is the interface to ACPI tables.
type Tabler interface {
	Len() int
	Base() int64
	Data() []byte
	Sig() string
	OEMID() string
	Revision() uint8
}
