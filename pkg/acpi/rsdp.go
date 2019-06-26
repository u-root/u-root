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

	"github.com/u-root/u-root/pkg/memio"
)

const (
	// Revision marks lowest ACPI revision we support.
	Revision    = 2
	cSUM1Off    = 8  // Checksum1 offset in packet.
	cSUM2Off    = 32 // Checksum2 offset
	xSDTLenOff  = 20
	xSDTAddrOff = 24
)

var pageMask = uint64(os.Getpagesize() - 1)

// RSDP is the v2 version of the RSDP struct, containing 32 and 64
// bit pointers.
// RSDP don't quite follow the ACPI table standard,
// so some things return empty values. It has nevertheless proven
// useful to have them.
// We just define the RSDP for v2 and later here. It's the only
// one that matters. This whole layout is typical of the overall
// Failure Of Vision that is ACPI. 64-bit micros had existed for 10 years
// when ACPI was defined, and they still wired in 32-bit pointer assumptions,
// and had to backtrack and fix it later. We don't use this struct below,
// it's only worthwhile as documentation. The RSDP has not changed in 20 years.
type RSDP struct {
	sign     [8]byte `Align:"16" Default:"RSDP PTR "`
	v1CSUM   uint8   // This was the checksum, which we are pretty sure is ignored now.
	oemid    [6]byte
	revision uint8  `Default:"2"`
	obase    uint32 // was RSDT, but you're not supposed to use it any more.
	length   uint32
	base     uint64 // XSDT address, the only one you should use
	checksum uint8
	_        [3]uint8
	data     [HeaderLength]byte
}

var (
	defaultRSDP = []byte("RSDP PTR U-ROOT\x02")
	_           = Tabler(&RSDP{})
)

// Marshal fails to marshal an RSDP.
func (r *RSDP) Marshal() ([]byte, error) {
	return nil, fmt.Errorf("Marshal RSDP: not yet")
}

// NewRSDP returns a new and partially initialized RSDP, setting only
// the defaultRSDP values, address, length, and signature.
func NewRSDP(addr uintptr, len uint) []byte {
	var r [HeaderLength]byte
	copy(r[:], defaultRSDP)
	// This is a bit of a cheat. All the fields are 0.
	// So we get a checksum, set up the
	// XSDT fields, get the second checksum.
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

// OEMTableID returns the RSDP OEMTableID
func (r *RSDP) OEMTableID() string {
	return "rsdp?"
}

// Revision returns the RSDP revision, which
// after 2002 should be >= 2
func (r *RSDP) Revision() uint8 {
	return r.revision
}

// OEMRevision returns the table OEMRevision.
func (r *RSDP) OEMRevision() uint32 {
	return 0
}

// CheckSum returns the table CheckSum.
func (r *RSDP) CheckSum() uint8 {
	return uint8(r.checksum)
}

// CreatorID returns the table CreatorID.
func (r *RSDP) CreatorID() uint32 {
	return uint32(0)
}

// CreatorRevision returns the table CreatorRevision.
func (r *RSDP) CreatorRevision() uint32 {
	return 0
}

// Base returns a base address or the [RX]SDT.
// It will preferentially return the XSDT, but if that is
// 0 it will return the RSDT address.
func (r *RSDP) Base() int64 {
	Debug("Base %v data len %d", r, len(r.data))
	b := int64(binary.LittleEndian.Uint32(r.data[16:20]))
	if b != 0 {
		return b
	}
	return int64(binary.LittleEndian.Uint64(r.data[24:32]))
}

func readRSDP(base int64) (*RSDP, error) {
	r := &RSDP{}
	r.base = uint64(base)
	dat := memio.ByteSlice(make([]byte, len(r.data)))
	if err := memio.Read(base, &dat); err != nil {
		return nil, err
	}
	copy(r.data[:], dat)
	return r, nil
}

func getRSDPEFI() (int64, *RSDP, error) {
	file, err := os.Open("/sys/firmware/efi/systab")
	if err != nil {
		return -1, nil, err
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
		base, err := strconv.ParseInt(start, 0, 64)
		if err != nil {
			continue
		}
		rsdp, err := readRSDP(base)
		if err != nil {
			continue
		}
		return base, rsdp, nil
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error while reading EFI systab: %v", err)
	}
	return -1, nil, fmt.Errorf("invalid efi/systab file")
}

func num(n string, i int) (uint64, error) {
	b, err := ioutil.ReadFile(fmt.Sprintf("/sys/firmware/memmap/%d/%s", i, n))
	if err != nil {
		return 0, err
	}
	start, err := strconv.ParseUint(string(b), 0, 64)
	return start, err
}

// getRSDPmem is the option of last choice, it just grovels through
// the e0000-ffff0 area, 16 bytes at a time, trying to find an RSDP.
// These are well-known addresses for 20+ years.
func getRSDPmem() (int64, *RSDP, error) {
	for base := int64(0xe0000); base < 0xffff0; base += 16 {
		var r memio.Uint64
		if err := memio.Read(base, &r); err != nil {
			continue
		}
		if r != 0x2052545020445352 {
			continue
		}
		rsdp, err := readRSDP(base)
		if err != nil {
			return -1, nil, err
		}
		return base, rsdp, nil
	}
	return -1, nil, fmt.Errorf("No ACPI RSDP via /dev/mem")
}

// You can change the getters if you wish for testing.
var getters = []func() (int64, *RSDP, error){getRSDPEFI, getRSDPmem}

// GetRSDP gets an RSDP.
// It is able to use several methods, because there is no consistency
// about how it is done. The base is also returned.
func GetRSDP() (base int64, rsdp *RSDP, err error) {
	for _, f := range getters {
		base, r, err := f()
		if err != nil {
			log.Print(err)
		}
		if err == nil {
			return base, r, nil
		}
	}
	return -1, nil, fmt.Errorf("Can't find an RSDP")
}
