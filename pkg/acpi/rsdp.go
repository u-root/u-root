// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/u-root/u-root/pkg/io"
)

const (
	Revision    = 2 // always
	RSDPLen     = 36
	CSUM1Off    = 8  // Checksum1 offset in packet.
	CSUM2Off    = 32 // Checksum2 offset
	XSDTLenOff  = 20
	XSDTAddrOff = 24
)

var pageMask = uint64(os.Getpagesize() - 1)

// We just define the real one for 2 and later here. It's the only
// one that matters. This whole layout is typical of the overall
// Failure Of Vision that is ACPI. 64-bit micros had existed for 10 years
// when ACPI was defined, and they still wired in 32-bit pointer assumptions,
// and had to backtrack and fix it later. We don't use this struct below,
// it's only worthwhile as documentation. The RSDP has not changed in 20 years.
type RSDP struct {
	sign     [8]byte `Align:"16", Default:"RSDP PTR "`
	v1CSUM   uint8   // This was the checksum, which we are pretty sure is ignored now.
	oemid    [6]byte
	revision uint8  `Default:"2"`
	_        uint32 // was RSDT, but you're not supposed to use it any more.
	length   uint32
	base     uint64 // XSDT address, the only one you should use
	checksum uint8
	_        [3]uint8
	data     [RSDPLen]byte
}

var (
	defaultRSDP = []byte("RSDP PTR U-ROOT\x02")
	_           = Tabler(&RSDP{})
)

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

// Len returns the RSDP length
func (r *RSDP) Len() int {
	return len(r.data)
}

// Base returns the RSDP base address as an int64
// We use int64 for compatibility with Go notions
// of offsets.
func (r *RSDP) Base() int64 {
	return int64(r.base)
}

// Data returns the RSDP as a []byte
func (r *RSDP) AllData() []byte {
	return r.data[:]
}

// TableData returns the RSDP table data as a []byte
func (r *RSDP) TableData() []byte {
	return r.data[36:]
}

// Sig returns the RSDP signature
func (r *RSDP) Sig() sig {
	return sig(r.data[:8])
}

// OEMID returns the RSDP OEMID
func (r *RSDP) OEMID() oem {
	return oem(r.data[9:15])
}

func (r *RSDP) OEMTableID() tableid {
	return "rsdp?"
}

// Revision returns the RSDP revision, which
// after 2002 should be >= 2
func (r *RSDP) Revision() u8 {
	return u8(r.revision)
}

func (r *RSDP) OEMRevision() u32 {
	return u32(0)
}

func (r *RSDP) CheckSum() u8 {
	return u8(r.checksum)
}

func (r *RSDP) CreatorID() u32 {
	return u32(0)
}

func (r *RSDP) VendorID() u32 {
	return u32(0)
}

func (r *RSDP) CreatorRevision() u32 {
	return u32(0)
}

func readRSDP(base int64) (*RSDP, error) {
	r := &RSDP{}
	r.base = uint64(base)
	for i := range r.data {
		var d io.Uint8
		if err := io.Read(base+int64(i), &d); err != nil {
			return nil, err
		}
		r.data[i] = uint8(d)
	}
	return r, nil
}

func getRSDPEFI() (*RSDP, error) {
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
		rsdp, err := strconv.ParseInt(start, 0, 64)
		if err != nil {
			continue
		}
		return readRSDP(rsdp)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error while reading EFI systab: %v", err)
	}
	return nil, fmt.Errorf("invalid efi/systab file")
}

func num(n string, i int) (uint64, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("/sys/firmware/memmap/%d/%s", i, n))
	if err != nil {
		return 0, err
	}
	start, err := strconv.ParseUint(string(b), 0, 64)
	return start, err
}

// get RSDPmem is the option of last choice, it just grovels through
// the e0000-ffff0 area, 16 bytes at a time, trying to find an RSDP.
// These are well-known addresses for 20+ years.
func getRSDPmem() (*RSDP, error) {
	for base := int64(0xe0000); base < 0xffff0; base += 16 {
		var r io.Uint64
		if err := io.Read(base, &r); err != nil {
			continue
		}
		if r != 0x2052545020445352 {
			continue
		}
		return readRSDP(base)
	}
	return nil, fmt.Errorf("No ACPI RSDP via /dev/mem")
}

func GetRSDP() (*RSDP, error) {
	for _, f := range []func() (*RSDP, error){getRSDPEFI, getRSDPmem} {
		r, err := f()
		if err == nil {
			return r, nil
		}
	}
	return nil, fmt.Errorf("Can't find an RSDP")
}
