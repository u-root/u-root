// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package coreboot

import (
	"bytes"
	"testing"

	"github.com/u-root/u-root/pkg/acpi"
)

func TestMADT(t *testing.T) {
	// This is a 36 byte table which is about as small as it can get.
	var rawAPIC = []byte{
		0x00, 0x50, 0x49, 0x43, 0x24, 0x00, 0x00, 0x00, 0x01, 0x8f, 0x50, 0x54, 0x4c, 0x54, 0x44, 0x20,
		0x09, 0x20, 0x41, 0x50, 0x49, 0x43, 0x20, 0x20, 0x00, 0x00, 0x04, 0x06, 0x20, 0x4c, 0x54, 0x50,
		0x00, 0x00, 0x00, 0x00,
		// Now we have entries.
	}

	tab, err := acpi.NewRaw(rawAPIC)
	if err != nil {
		t.Fatalf("acpi.NewRaw: got %v, want nil", err)
	}
	// Test geting a MADT with a bad signature
	if _, err = NewMADT(tab[0]); err == nil {
		t.Errorf("MADT with incorrect signature: got nil, want err")
	}
	// Test geting a MADT that is too short.
	rawAPIC[0] = 'A'
	if _, err = acpi.NewRaw(rawAPIC); err != nil {
		t.Fatalf("acpi.NewRaw: got %v, want nil", err)
	}
	if _, err = NewMADT(tab[0]); err == nil {
		t.Errorf("MADT that is too short: got nil, want err")
	}

	// Test a good minimal APIC table
	rawAPIC = append(rawAPIC, 0xde, 0xad, 0xbe, 0xef, 1, 0, 0, 0)
	rawAPIC[4] += 8
	if tab, err = acpi.NewRaw(rawAPIC); err != nil {
		t.Fatalf("acpi.NewRaw: got %v, want nil", err)
	}
	m, err := NewMADT(tab[0])
	if err != nil {
		t.Errorf("New good MADT: got %v, want nil", err)
	}
	var b bytes.Buffer
	if err := m.Coreboot(&b); err != nil {
		t.Errorf("CoreBoot: got %v, want nil", err)
	}
	// Now add too short subtable entry
	// i.e. we have a good madt with no subentries; let's add a bogus entry
	// that is too short.
	Debug = t.Logf
	tooshort := append(rawAPIC, 1)
	tooshort[4]++
	if tab, err = acpi.NewRaw(tooshort); err != nil {
		t.Fatalf("acpi.NewRaw: got %v, want nil", err)
	}
	if _, err = NewMADT(tab[0]); err == nil {
		t.Fatalf("Too short MADT subtable: got nil, want err")
	}
	// now try everything.
	t.Logf("rawapic len %d", len(rawAPIC))
	rawAPIC = append(rawAPIC, byte(entryLAPIC), sizeEntryLAPIC, 1, 2, 3, 0, 0, 0)
	rawAPIC = append(rawAPIC, byte(entryIOAPIC), sizeEntryIOAPIC, 32, 0, 4, 5, 6, 7, 9, 10, 11, 12)
	rawAPIC = append(rawAPIC, byte(entryISOR), sizeEntryISOR, 22, 23, 1, 2, 3, 4, 0, 1)
	rawAPIC = append(rawAPIC, byte(entryNMI), sizeEntryNMI, 0xfe, 0xfd, 0xfc, 0xfb)

	t.Logf("rawapic len %d", len(rawAPIC))
	rawAPIC[4] = uint8(len(rawAPIC))
	if tab, err = acpi.NewRaw(rawAPIC); err != nil {
		t.Fatalf("acpi.NewRaw: got %v, want nil", err)
	}
	t.Logf("Table %dbytes\n", tab[0].TableData())
	Debug = t.Logf
	m, err = NewMADT(tab[0])
	if err != nil {
		t.Errorf("New good MADT: got %v, want nil", err)
	}
	w := &bytes.Buffer{}
	if err := m.Coreboot(w); err != nil {
		t.Error(err)
	}
	good := `unsigned long acpi_fill_madt(unsigned long current)
{
	/* create all subtables for processors */
	current += acpi_create_madt_lapic((acpi_madt_lapic_t *)current, 0x1, 0x2);
	current += acpi_create_madt_ioapic((acpi_madt_ioapic_t *)current, 0x20, 0x07060504, 0x0c0b0a09);
	current += acpi_create_madt_irqoverride((acpi_madt_irqoverride_t *)current, 0x16, 0x17, 0x04030201, 0x0100);
	current += acpi_create_madt_lapic_nmi((acpi_madt_lapic_nmi_t *)current, 0xfe, 0xfcfd, 0xfb);
	return current;
}
`
	out := w.String()
	if len(out) != len(good) {
		t.Errorf("coreboot: got %d byte string, want %d bytes", len(out), len(good))
	}
	if w.String() != good {
		t.Errorf("Coreboot: got %v, want %v", w.String(), good)
		t.Errorf("Coreboot: got %q, want %q", w.String(), good)
	}
	// now shorten it by 1 to simulate a truncated table
	rawAPIC = rawAPIC[:len(rawAPIC)-1]
	rawAPIC[4] = uint8(len(rawAPIC))
	if tab, err = acpi.NewRaw(rawAPIC); err != nil {
		t.Errorf("acpi.NewRaw: got %v, want nil", err)
	}
	if _, err = NewMADT(tab[0]); err == nil {
		t.Errorf("Short MADT subtable: got nil, want err")
	}

}
