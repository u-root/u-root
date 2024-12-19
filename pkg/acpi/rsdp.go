// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"encoding/binary"
	"fmt"

	"github.com/u-root/u-root/pkg/memio"
)

var defaultRSDP = []byte("RSDP PTR U-ROOT\x02")

// RSDP is the v2 version of the ACPI RSDP struct.
type RSDP struct {
	// base is the base address of the RSDP struct in physical memory.
	base int64

	data [headerLength]byte
}

// NewRSDP returns a new and partially initialized RSDP, setting only
// the defaultRSDP values, address, length, and signature.
func NewRSDP(addr uintptr, length uint) []byte {
	var r [headerLength]byte
	copy(r[:], defaultRSDP)

	// This is a bit of a cheat. All the fields are 0.  So we get a
	// checksum, set up the XSDT fields, get the second checksum.
	r[cSUM1Off] = gencsum(r[:])
	binary.LittleEndian.PutUint32(r[xSDTLenOff:], uint32(length))
	binary.LittleEndian.PutUint64(r[xSDTAddrOff:], uint64(addr))
	r[cSUM2Off] = gencsum(r[:])
	return r[:]
}

// Len returns the RSDP length
func (r *RSDP) Len() uint32 {
	return uint32(len(r.data))
}

// AllData returns the RSDP as a []byte
func (r *RSDP) AllData() []byte {
	return r.data[:]
}

// TableData returns the RSDP table data as a []byte
func (r *RSDP) TableData() []byte {
	return r.data[36:]
}

// Sig returns the RSDP signature
func (r *RSDP) Sig() string {
	return string(r.data[:8])
}

// OEMID returns the RSDP OEMID
func (r *RSDP) OEMID() string {
	return string(r.data[9:15])
}

// RSDPAddr returns the physical base address of the RSDP.
func (r *RSDP) RSDPAddr() int64 {
	return r.base
}

// SDTAddr returns a base address or the [RX]SDT.
//
// It will preferentially return the XSDT, but if that is
// 0 it will return the RSDT address.
func (r *RSDP) SDTAddr() int64 {
	b := uint64(binary.LittleEndian.Uint32(r.data[16:20]))
	if b != 0 {
		return int64(b)
	}
	return int64(binary.LittleEndian.Uint64(r.data[24:32]))
}

func readRSDP(base int64) (*RSDP, error) {
	r := &RSDP{}
	r.base = base

	dat := memio.ByteSlice(make([]byte, len(r.data)))
	if err := memio.Read(int64(base), &dat); err != nil {
		return nil, err
	}
	copy(r.data[:], dat)
	return r, nil
}

// GetRSDP finds the RSDP pointer and struct. The rsdpgetters must be defined
// in rsdp_$(GOOS).go, since, e.g.,OSX, BSD, and Linux have some intersections
// but some unique aspects too, and Plan 9 has nothing in common with any of them.
//
// It is able to use several methods, because there is no consistency
// about how it is done.
func GetRSDP() (*RSDP, error) {
	for _, f := range rsdpgetters {
		r, err := f()
		if err == nil {
			return r, nil
		}
	}
	return nil, fmt.Errorf("cannot find an RSDP")
}
