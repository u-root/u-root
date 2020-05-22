// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreboot

import (
	"fmt"
	"io"

	"github.com/u-root/u-root/pkg/acpi"
)

type madtType uint8

const (
	minMADT madtType = 8 // 32-bit address and flags
)
const (
	entryLAPIC         madtType = 0
	entryIOAPIC        madtType = 1
	entryISOR          madtType = 2 // pronounced: eyesore, for reasons.
	entryBOGUS         madtType = 3
	entryNMI           madtType = 4
	entryLAPICOverRide madtType = 5
	entryLX2APIC       madtType = 9
	entryLX2NMI        madtType = 10

	sizeEntryLAPIC         uint8 = 8
	sizeEntryIOAPIC        uint8 = 12
	sizeEntryISOR          uint8 = 10 // pronounced: eyesore, for reasons.
	sizeEntryNMI           uint8 = 6
	sizeEntryLAPICOverRide uint8 = 12
	sizeEntryLX2APIC       uint8 = 16
	sizeEntryLX2NMI        uint8 = 12
)

var (
	// The MADT signature is APIC.
	// Another failure of vision; they never thought there'd
	// be more than one APIC I guess.
	sigMADT = "APIC"
	info    = map[madtType]struct {
		name string
		size uint8
	}{
		entryLAPIC:         {"LAPIC", sizeEntryLAPIC},
		entryIOAPIC:        {"IOAPIC", sizeEntryIOAPIC},
		entryISOR:          {"ISOR", sizeEntryISOR},
		entryNMI:           {"NMI", sizeEntryNMI},
		entryLAPICOverRide: {"LAPICOverRide", sizeEntryLAPICOverRide},
		entryLX2APIC:       {"LX2APIC", sizeEntryLX2APIC},
		entryLX2NMI:        {"X2NMI", sizeEntryLX2NMI},
	}
)

type subtable struct {
	t madtType
	b []byte
}

// MADT defines the madt struct.
type MADT struct {
	acpi.Table
	// base is the base address of the MADT struct in physical memory.
	base      int64
	subtables []subtable
}

var _ Corebooter = &MADT{}

func init() {
	corebooters["APIC"] = NewMADT
}

// rv reverses a slice
func rv(b []byte) []byte {
	for i := 0; i < len(b)/2; i++ {
		j := len(b) - i - 1
		b[i], b[j] = b[j], b[i]
	}
	return b
}

// NewMADT returns a MADT, from a Table, iff the signature is correct.
func NewMADT(t acpi.Table) (Corebooter, error) {
	if t.Sig() != sigMADT {
		return nil, fmt.Errorf("%s is not a MADT", t.Sig())
	}
	b := t.TableData()
	Debug("NEWMADT Subtables: %d bytes", len(b))
	if len(b) < int(minMADT) {
		return nil, fmt.Errorf("Table is only %d bytes, need %d", len(t.TableData()), minMADT)
	}
	b = b[8:]
	Debug("NEWMADT Subtables: after removing 8 bytes of fluff, %d bytes", len(b))
	m := &MADT{Table: t, base: t.Address()}
	for len(b) > 0 {
		Debug("NewMADT loop: %d bytes\n", len(b))
		if len(b) < 2 {
			return nil, fmt.Errorf("subtable is too short: have %d bytes, need 2", len(b))
		}
		l := b[1]
		typ := madtType(b[0])
		Debug("subtable: %d bytes, l is %#x, t is %#x", len(b), l, t)
		i, ok := info[typ]
		if ! ok {
			Debug("Bad type: %#x, offset %#x, skipping %d bytes", typ, len(t.TableData())-len(b), l)
			b = b[l:]
			continue
		}
		if l != i.size {
			return nil, fmt.Errorf("Type %s has size %d should be %d", i.name, l, i.size)
		}
		if len(b) < int(l) {
			return nil, fmt.Errorf("Type %s has buf length %d should be %d", i.name, len(b), i.size)
		}
		m.subtables = append(m.subtables, subtable{t: typ, b: b[:l]})
		b = b[l:]
	}
	Debug("NEWMADT Subtables: %d bytes tab %d bytes", len(b), len(m.TableData()))
	return m, nil
}

// CoreBoot generates the coreboot C code for a MADT.
// MADT consists of TLV entries: 1-byte type, 1-byte length, value.
func (m *MADT) Coreboot(w io.Writer) error {
	s := fmt.Sprintf(`unsigned long acpi_fill_madt(unsigned long current)
{
	/* create all subtables for processors */
`)
	// OK this is a bit nasty. Right now coreboot is trying to fill out
	// the lapic address and flags. So we ignore those in the table for now.
	// If it's getting them wrong, we will have to go back and
	// fill them in.
	//lapic(%#x, %#x);\n", b[:4], b[4:8])
	// The coreboot serialization code for ACPI uses the C compiler to drive serialization.
	// It uses this pattern:
	// int acpi_create_madt_lapic(acpi_madt_lapic_t *lapic, u8 cpu, u8 apic);
	// The first arg is a uintptr cast to a struct pointer of the table type.
	// The function stores the arguments into the packed struct, the idea being that
	// the compiler generates code to store into the struct that exactly matches
	// the ACPI standard.
	for _, tab := range m.subtables {
		b := tab.b
		switch tab.t {
		case entryLAPIC:
			// TODO: Check b[4] to see if the LAPIC is enabled? When would it not be, if it's in the table?
			s += fmt.Sprintf("\tcurrent += acpi_create_madt_lapic((acpi_madt_lapic_t *)current, %#x, %#x);\n", b[2], b[3])
		case entryIOAPIC:
			s += fmt.Sprintf("\tcurrent += acpi_create_madt_ioapic((acpi_madt_ioapic_t *)current, %#x, %#x, %#x);\n", b[2], rv(b[4:8]), rv(b[8:12]))
		case entryISOR:
			s += fmt.Sprintf("\tcurrent += acpi_create_madt_irqoverride((acpi_madt_irqoverride_t *)current, %#x, %#x, %#x, %#x);\n", b[2], b[3], rv(b[4:8]), rv(b[8:10]))
		case entryNMI:
			s += fmt.Sprintf("\tcurrent += acpi_create_madt_lapic_nmi((acpi_madt_lapic_nmi_t *)current, %#x, %#x, %#x);\n", b[2], rv(b[3:5]), b[5])
		case entryLAPICOverRide:
			// coreboot appears not to support this.
			Debug("coreboot does not support Local APIC Address Override, skipping")
			//s += fmt.Sprintf("current += acpi_create_madt_irqoverride((acpi_madt_irqoverride_t *), %#x, %#x);\n", rv(b[2:4]), b[4:12])
		case entryLX2APIC:
			// coreboot does not support the flags parameter ...
			s += fmt.Sprintf("\tcurrent += acpi_create_madt_lx2apic((acpi_madt_lx2apic_t *) current,  %#x, %#x);\n", rv(b[12:16]), rv(b[4:8]))
		case entryLX2NMI:
			s += fmt.Sprintf("\tcurrent += acpi_create_madt_lx2apic_nmi((acpi_madt_lx2apic_nmi_t *) current,  %#x, %#x, %#x);\n", rv(b[4:8]), rv(b[2:4]), b[8:9])
		}
	}
	s += "\treturn current;\n}\n"
	_, err := w.Write([]byte(s))
	return err
}
