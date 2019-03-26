// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The ACPI package began life as a relatively simple set of functions.
// At some point we wanted to be able to extend the ACPI table over a
// kexec and that's when it all went nasty.
// A few issues crop up.
// In theory, one should be able to put tables anywhere. In practice,
// even in the EFI age, this is not simple. The root pointer
// for tables, RSDP, Root Services Data Pointer, is in memory protected by
// the southbridge/ICH/PCH with lockdown bits (disable writes once, can
// not be reset until power on/reset). Writes to it are transparently discarded.
// (This makes sense, as crap DOS programs probably tried to write that area all
// the time. Everything -- everything! -- about x86 and BIOS is explainable
// by DOS, 16 bits addresses, 4 bit segments, 20 bits, and the incredible failure
// of vision in 1999 of building a bunch of tables with 32 bit pointers in them.
// The x86 world is nothing if not consistent: short-sighted decisions for 40+ years).
// The two-byte EBDA pointer (16 bits) is shifted left 4 bits to get a 20-bit address
// which points to *the range of memory* containing, maybe, the RSDP.
// The RSDP has to live in the low 20 bits of address space (remember 20 bits, right?)
// The RSDP has two pointers in, one of both containing a value: 32 bit pointer or 64 bit
// pointer to the RSDT or XSDT. RSDT has 32-bit pointers, XSDT 64-bit pointers.
// The spec recommends you ignore the 32-bit pointer variants but as of 2019,
// I still see use of them; this package will read and write the variants but the
// internal structs are defined with 64-bit pointers (oh, if only the ACPICA code
// had gotten this right, but that code is pretty bad too).
// What's really fun: when you generate the [RX]SDT, you need to generate pointers
// too. It's best to keep those tables in a physically contiguous area, so we do;
// not all systems follow this rule, which is asking for trouble.
// Finally, not all tables are generated in a consistent way, e.g. the IBFT is not like
// most other tables containing variable length data, because the people who created it
// are Klever. IBFT has its own heap.
// If you really look at ACPI closely, you can kind of see that it's optimized for NASM,
// which is part of what makes it so unpleasant. But, it's what we have.
package acpi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
)

type ACPIWriter interface {
	Marshal() ([]byte, error)
}

const (
	// LengthOffset is the offset of the table length
	LengthOffset = 4
	// CSUMOffset is the offset of the single byte checksum in *most* ACPI tables
	CSUMOffset = 9
	// MinTableLength is the minimum length: 4 byte tag, 4 byte length, 1 byte revision, 1 byte checksum,
	MinTableLength = 10
)

var (
	// Debug implements fmt.Sprintf and can be used for debug printing
	Debug        = func(string, ...interface{}) {}
	unmarshalers = map[sig]func(Tabler) (Tabler, error){}
)

func addUnMarshaler(n string, f func(Tabler) (Tabler, error)) {
	if _, ok := unmarshalers[sig(n)]; ok {
		log.Fatalf("Can't add %s; already in use", n)
	}
	unmarshalers[sig(n)] = f
}

// This is the standard header for all ACPI tables, except the
// ones that don't use it.
// We use types that we hope are easy to read; they in turn
// make writing marshal code with type switches very convenient.
type Header struct {
	sig             sig
	length          u32
	revision        u8
	checkSum        u8
	oemID           oem
	oemTableID      tableid
	oemRevision     u32
	creatorID       u32
	creatorRevision u32
}

func GetHeader(t Tabler) *Header {
	return &Header{
		sig:             t.Sig(),
		length:          u32(t.Len()),
		revision:        t.Revision(),
		checkSum:        t.CheckSum(),
		oemID:           t.OEMID(),
		oemTableID:      t.OEMTableID(),
		oemRevision:     t.OEMRevision(),
		creatorID:       t.CreatorID(),
		creatorRevision: t.CreatorRevision(),
	}
}

// Table is a basic ACPI table. All tables consist of the
// Header and then a bunch of table-specific data.
// Because we do not care about most tables, we have
// this "generic" type which we can Unmarshal into
// and Marshal from.
type Table struct {
	Header
	Data []byte
}

// Flags takes 0 or more flags and produces a uint32 value.
// Each argument represents one bit position, with the first flag
// being bit 0. For each flag with a value of "1", the bit for that
// flag will be set. For each flag with a value of "0", that bit in
// the flag will be cleared.
// This is mainly for consistency but it would allow
// us in future to pass in an initial value and set or clear bits in it
// depending on flags.
// Current allowed values are "0" and "1". In future, if the flags are
// not contiguous, we can allow "ignore" in future to ignore a flag.
func flags(s ...flag) (uint8, error) {
	var i, bit uint8
	for _, f := range s {
		switch f {
		case "1":
			i |= 1 << bit
		case "0":
			i &= ^(1 << bit)
		default:
			return 0, fmt.Errorf("%s is not a valid value: only 0 or 1 are valied", f)
		}
		bit++
	}
	return i, nil
}

// w writes 0 or more values to a bytes.Buffer, in LittleEndian order,
func w(b *bytes.Buffer, val ...interface{}) {
	for _, v := range val {
		binary.Write(b, binary.LittleEndian, v)
		Debug("\t %T %v b is %d bytes", v, v, b.Len())
	}
	Debug("w: done: b is %d bytes", b.Len())
}

// uw writes strings as unsigned words to a bytes.Buffer.
// Currently it only supports 16, 32, and 64 bit writes.
func uw(b *bytes.Buffer, s string, bits int) error {
	// convenience case: if they don't set it, it comes in as "",
	// take that to mean 0.
	var v uint64
	if s != "" {
		var err error
		if v, err = strconv.ParseUint(string(s), 0, bits); err != nil {
			return err
		}
	}
	switch bits {
	case 8:
		w(b, uint8(v))
	case 16:
		w(b, uint16(v))
	case 32:
		w(b, uint32(v))
	case 64:
		w(b, uint64(v))
	}
	return nil

}

// Marshal marshals a single ACPI table into a byte slice.
func Marshal(i ACPIWriter) ([]byte, error) {

	Debug("Marshall %T", i)
	b, err := i.Marshal()
	if err != nil {
		return nil, err
	}

	if len(b) < MinTableLength {
		return nil, fmt.Errorf("%v is too short to contain a table", b)
	}

	binary.LittleEndian.PutUint32(b[LengthOffset:], uint32(len(b)))
	c := gencsum(b)
	Debug("CSUM is %#x", c)
	b[CSUMOffset] = c

	return b, nil
}

// UnMarshall unmarshals a single table.
// If the table is one of the many we don't care about we
// just return a raw table, which can be easily written out
// again if needed. If it has an UnMarshal registered we use
// that instead.
func UnMarshal(b []byte) (Tabler, error) {
	r, err := NewRaw(b)
	if err != nil {
		return nil, err
	}
	if m, ok := unmarshalers[r.Sig()]; ok {
		return m(r)
	}

	return r, nil
}

// UnMarshalSDT unmarshals an SDT.
// It's pretty much impossible for the RSDP to point to
// anything else so we mainly do the unmarshal and check the sig.
func UnMarshallSDT(r *RSDP) (*SDT, error) {
	// suck in the raw table, then marshal it.
	// There should have been a marshaler registered.
	raw, err := ReadRaw(r.Base())
	if err != nil {
		return nil, err
	}
	s, err := UnMarshal(raw.AllData())
	if err != nil {
		return nil, err
	}
	return s.(*SDT), nil
}

func UnMarshalAll(s *SDT) (*[]Table, error) {
	return nil, nil
}
