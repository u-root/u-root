// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package acpi can find and parse the RSDP pointer and struct.
package acpi

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/boot/ebda"
	"github.com/u-root/u-root/pkg/memio"
)

const (
	// checksum1 offset in RSDP struct.
	cSUM1Off = 8

	// checksum2 offset in RSDP struct.
	cSUM2Off    = 32
	xSDTLenOff  = 20
	xSDTAddrOff = 24

	// headerLength is a common header length for (almost)
	// all ACPI tables.
	headerLength = 36
)

var (
	defaultRSDP = []byte("RSDP PTR U-ROOT\x02")
)

// gencsum generates a uint8 checksum of a []uint8
func gencsum(b []uint8) uint8 {
	var csum uint8
	for _, bb := range b {
		csum += bb
	}
	return ^csum + 1
}

// RSDP is the v2 version of the ACPI RSDP struct.
type RSDP struct {
	// base is the base address of the RSDP struct in physical memory.
	base uint64

	data [headerLength]byte
}

// NewRSDP returns a new and partially initialized RSDP, setting only
// the defaultRSDP values, address, length, and signature.
func NewRSDP(addr uintptr, len uint) []byte {
	var r [headerLength]byte
	copy(r[:], defaultRSDP)

	// This is a bit of a cheat. All the fields are 0.  So we get a
	// checksum, set up the XSDT fields, get the second checksum.
	r[cSUM1Off] = gencsum(r[:])
	binary.LittleEndian.PutUint32(r[xSDTLenOff:], uint32(len))
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
func (r *RSDP) RSDPAddr() uint64 {
	return r.base
}

// SDTAddr returns a base address or the [RX]SDT.
//
// It will preferentially return the XSDT, but if that is
// 0 it will return the RSDT address.
func (r *RSDP) SDTAddr() uint64 {
	b := uint64(binary.LittleEndian.Uint32(r.data[16:20]))
	if b != 0 {
		return b
	}
	return uint64(binary.LittleEndian.Uint64(r.data[24:32]))
}

func readRSDP(base uint64) (*RSDP, error) {
	r := &RSDP{}
	r.base = base

	dat := memio.ByteSlice(make([]byte, len(r.data)))
	if err := memio.Read(int64(base), &dat); err != nil {
		return nil, err
	}
	copy(r.data[:], dat)
	return r, nil
}

// GetRSDPEFI finds the RSDP in the EFI System Table.
func GetRSDPEFI() (*RSDP, error) {
	file, err := os.Open("/sys/firmware/efi/systab")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	const (
		acpi20 = "ACPI20="
		acpi   = "ACPI="
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		start := ""
		if strings.HasPrefix(line, acpi20) {
			start = strings.TrimPrefix(line, acpi20)
		}
		if strings.HasPrefix(line, acpi) {
			start = strings.TrimPrefix(line, acpi)
		}
		if start == "" {
			continue
		}
		base, err := strconv.ParseUint(start, 0, 64)
		if err != nil {
			continue
		}
		rsdp, err := readRSDP(base)
		if err != nil {
			continue
		}
		return rsdp, nil
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error while reading EFI systab: %v", err)
	}
	return nil, fmt.Errorf("invalid /sys/firmware/efi/systab file")
}

// GetRSDPEBDA finds the RSDP in the EBDA.
func GetRSDPEBDA() (*RSDP, error) {
	f, err := os.OpenFile("/dev/mem", os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	e, err := ebda.ReadEBDA(f)
	if err != nil {
		return nil, err
	}

	return getRSDPMem(uint64(e.BaseOffset), uint64(e.BaseOffset+e.Length))
}

func getRSDPMem(start, end uint64) (*RSDP, error) {
	for base := start; base < end; base += 16 {
		var r memio.Uint64
		if err := memio.Read(int64(base), &r); err != nil {
			continue
		}
		if r != 0x2052545020445352 {
			continue
		}
		rsdp, err := readRSDP(base)
		if err != nil {
			return nil, err
		}
		return rsdp, nil
	}
	return nil, fmt.Errorf("could not find ACPI RSDP via /dev/mem from %#08x to %#08x", start, end)
}

// GetRSDPMem is the option of last choice, it just grovels through
// the e0000-ffff0 area, 16 bytes at a time, trying to find an RSDP.
// These are well-known addresses for 20+ years.
func GetRSDPMem() (*RSDP, error) {
	return getRSDPMem(0xe0000, 0xffff0)
}

// You can change the getters if you wish for testing.
var getters = []func() (*RSDP, error){GetRSDPEFI, GetRSDPEBDA, GetRSDPMem}

// GetRSDP finds the RSDP pointer and struct in memory.
//
// It is able to use several methods, because there is no consistency
// about how it is done.
func GetRSDP() (*RSDP, error) {
	for _, f := range getters {
		r, err := f()
		if err != nil {
			log.Print(err)
		}
		if err == nil {
			return r, nil
		}
	}
	return nil, fmt.Errorf("cannot find an RSDP")
}
