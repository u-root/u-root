// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import "encoding/binary"

const (
	Revision    = 2 // always
	RSDPLen     = 36
	CSUM1Off    = 8  // Checksum1 offset in packet.
	CSUM2Off    = 32 // Checksum2 offset
	XSDTLenOff  = 20
	XSDTAddrOff = 24
)

// We just define the real one for 2 and later here. It's the only
// one that matters. This whole layout is typical of the overall
// Failure Of Vision that is ACPI. 64-bit micros had existed for 10 years
// when ACPI was defined, and they still wired in 32-bit pointer assumptions,
// and had to backtrack and fix it later. We don't use this struct below,
// it's only worthwhile as documentation. The RSDP has not changed in 20 years.
type RSDP struct {
	Signature [8]byte `Align:"16", Default:"RSDP PTR "`
	V1CSUM    uint8   // This was the checksum, which we are pretty sure is ignored now.
	OEMID     [6]byte
	Revision  uint8  `Default:"2"`
	_         uint32 // was RSDT, but you're not supposed to use it any more.
	Length    uint32
	Address   uint64 // XSDT address, the only one you should use
	Checksum  uint8
	_         [3]uint8
}

var defaultRSDP = []byte("RSDP PTR U-ROOT\x02")

func NewRSDP(addr uintptr, len uint) []byte {
	var r [RSDPLen]byte
	copy(r[:], defaultRSDP)
	// This is a bit of a cheat. All the fields are 0.
	// So we get a checksum, set up the
	// XSDT fields, get the second checksum.
	r[CSUM1Off] = gencsum(r[:])
	binary.LittleEndian.PutUint32(r[XSDTLenOff:], uint32(len))
	binary.LittleEndian.PutUint64(r[XSDTAddrOff:], uint64(addr))
	r[CSUM2Off] = gencsum(r[:])
	return r[:]
}
