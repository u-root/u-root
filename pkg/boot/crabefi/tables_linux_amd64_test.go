// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package crabefi

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/u-root/u-root/pkg/boot/kexec"
)

func TestCorebootChecksum(t *testing.T) {
	for _, tt := range []struct {
		data []byte
		want uint16
	}{
		{nil, 0xffff},
		{[]byte{1}, 0xfffe},
		{[]byte{0x12, 0x34}, 0xcbed},
		{[]byte{1, 2, 3}, 0xfdfb},
		{[]byte{0xff, 0xff, 1}, 0xfffe},
	} {
		if got := checksum(tt.data); got != tt.want {
			t.Errorf("checksum(%x) = %#x; want %#x", tt.data, got, tt.want)
		}
	}
}

func TestSMBIOSRelocationBounds(t *testing.T) {
	for _, version := range []int{2, 3} {
		f := testFirmware(t, version)
		for _, data := range [][]byte{nil, append(append([]byte(nil), f.smbiosData...), 0)} {
			if _, err := relocateSMBIOS(f.smbiosEntry, data, 0x8000000); err == nil {
				t.Errorf("relocateSMBIOS(version %d, %d bytes) = nil; want error", version, len(data))
			}
		}
	}
	f := testFirmware(t, 2)
	if _, err := relocateSMBIOS(f.smbiosEntry, f.smbiosData, uint64(math.MaxUint32)+1); err == nil {
		t.Error("relocateSMBIOS(SMBIOS 2 above 4 GiB) = nil; want error")
	}
}

func TestHandoffMemoryTypes(t *testing.T) {
	f := testFirmware(t, 3)
	memory := kexec.MemoryMap{
		{Range: kexec.Range{Start: 0x1000, Size: 4096}, Type: kexec.RangeRAM},
		{Range: kexec.Range{Start: 0x2000, Size: 4096}, Type: kexec.RangeReserved},
		{Range: kexec.Range{Start: 0x3000, Size: 4096}, Type: kexec.RangeACPI},
		{Range: kexec.Range{Start: 0x4000, Size: 4096}, Type: kexec.RangeNVS},
		{Range: kexec.Range{Start: 0x5000, Size: 4096}, Type: kexec.RangeType("unknown")},
	}
	data, err := makeHandoff(memory, 0x8004000, f.rsdp, f.framebuffer, f.smbiosEntry, f.smbiosData)
	if err != nil {
		t.Fatal(err)
	}
	r := bytes.NewReader(data[24:])
	var record recordHeader
	if err := binary.Read(r, binary.LittleEndian, &record); err != nil {
		t.Fatal(err)
	}
	if record.Tag != 1 || record.Size != 8+20*uint32(len(memory)) {
		t.Fatalf("memory record = %+v; want coreboot memory tag and packed 20-byte ranges", record)
	}
	for i, want := range []uint32{1, 2, 3, 4, 2} {
		var region memoryRange
		if err := binary.Read(r, binary.LittleEndian, &region); err != nil {
			t.Fatal(err)
		}
		if region.Kind != want || region.Base != uint64(memory[i].Start) || region.Size != uint64(memory[i].Size) {
			t.Errorf("memory record %d = %+v; want %#x, 4096, type %d", i, region, memory[i].Start, want)
		}
	}
	// Too many records to fit before the fixed SMBIOS entry offset must
	// fail rather than silently overwrite another handoff structure.
	if _, err := makeHandoff(make(kexec.MemoryMap, 256), 0x8004000, f.rsdp, f.framebuffer, f.smbiosEntry, f.smbiosData); err == nil {
		t.Error("makeHandoff(overlapping records) = nil; want error")
	}
}
