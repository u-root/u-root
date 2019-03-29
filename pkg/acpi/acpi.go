// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package acpi began life as a relatively simple set of functions.
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
// by DOS, 16 bit addresses, 4 bit segments, 20 bits, and the incredible failure
// of vision in 1999 of building a bunch of tables with 32 bit pointers in them.
// The x86 world is nothing if not consistent: short-sighted decisions for 40+ years).
// The two-byte EBDA pointer (16 bits) is shifted left 4 bits to get a 20-bit address
// which points to *the range of memory* containing, maybe, the RSDP.
// The RSDP has to live in the low 20 bits of address space (remember 20 bits, right?)
// (Yes, I know about the EFI tables which make this restriction less of a problem,
//  but for now all the firmware I've seen adheres to the "RSDP in e or f segment" rule.)
// The RSDP has two pointers in it, one or both containing a value: 32 bit pointer or 64 bit
// pointer to the RSDT or XSDT. (see "failure of vision" above.)
// RSDT has 32-bit pointers, XSDT 64-bit pointers.
// The spec recommends you ignore the 32-bit pointer variants but as of 2019,
// I still see use of them; this package will read and write the variants but the
// internal structs are defined with 64-bit pointers (oh, if only the ACPICA code
// had gotten this right, but that code is pretty bad too).
// What's really fun: when you generate the [RX]SDT, you need to generate pointers
// too. It's best to keep those tables in a physically contiguous area, so we do;
// not all systems follow this rule, which is asking for trouble.
// Finally, not all tables are generated in a consistent way, e.g. the IBFT is not like
// most other tables containing variable length data, because the people who created it
// are Klever. IBFT has its own heap, and the elements in the heap follow rules unlike
// all other tables. Nice!
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

// addUnMarshaler is intended to be called by init functions
// in this package. It adds an UnMarshaler for a given
// ACPI signature.
func addUnMarshaler(n string, f func(Tabler) (Tabler, error)) {
	if _, ok := unmarshalers[sig(n)]; ok {
		log.Fatalf("Can't add %s; already in use", n)
	}
	unmarshalers[sig(n)] = f
}

// GetHeader extracts a Header from a Tabler and returns a reference to it.
func GetHeader(t Tabler) *Header {
	return &Header{
		Sig:             sig(t.Sig()),
		Length:          t.Len(),
		Revision:        t.Revision(),
		CheckSum:        t.CheckSum(),
		OEMID:           oem(t.OEMID()),
		OEMTableID:      tableid(t.OEMTableID()),
		OEMRevision:     t.OEMRevision(),
		CreatorID:       t.CreatorID(),
		CreatorRevision: t.CreatorRevision(),
	}
}

// Flags takes 0 or more flags and produces a uint8 value.
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
			return 0, fmt.Errorf("%s is not a valid value: only 0 or 1 are valid", f)
		}
		bit++
	}
	return i, nil
}

// w writes 0 or more values to a bytes.Buffer, in LittleEndian order.
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
	default:
		return fmt.Errorf("Invalid bit length for uw: %d, only supoprt 8, 16, 32, 64", bits)
	}
	return nil
}

// Marshal marshals a Tabler into a byte slice.
// Once marshaling is done, it inserts the length into
// the standard place at LengthOffset, and then generates and inserts
// a checksum at CSUMOffset.
func Marshal(t Tabler) ([]byte, error) {
	Debug("Marshal %T", t)
	b, err := t.Marshal()
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

// UnMarshal unmarshals a single table and returns a Tabler.
// If the table is one of the many we don't care about we
// just return a Raw table, which can be easily written out
// again if needed. If it has an UnMarshal registered we use
// that instead once the Raw is unmarshaled.
func UnMarshal(a int64) (Tabler, error) {
	r, err := ReadRaw(a)
	if err != nil {
		return nil, err
	}
	Debug("Raw table: %q", r)
	if m, ok := unmarshalers[sig(r.Sig())]; ok {
		return m(r)
	}

	return r, nil
}

// UnMarshalSDT unmarshals an SDT.
// It's pretty much impossible for the RSDP to point to
// anything else so we mainly do the unmarshal and type assertion.
func UnMarshalSDT(r *RSDP) (*SDT, error) {
	s, err := UnMarshal(r.Base())
	if err != nil {
		return nil, err
	}
	Debug("SDT: %q", s)
	return s.(*SDT), nil
}

// UnMarshalAll takes an SDT and unmarshals all the tables
// using UnMarshal. It returns a []Tabler. In most cases,
// the tables will be Raw, but in a few cases they might be
// further converted.
func UnMarshalAll(s *SDT) ([]Tabler, error) {
	var tab []Tabler
	for _, a := range s.Tables {
		t, err := UnMarshal(a)
		if err != nil {
			return nil, err
		}
		tab = append(tab, t)
	}

	return tab, nil
}
