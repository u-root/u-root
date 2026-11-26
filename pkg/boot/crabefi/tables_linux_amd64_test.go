// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !tinygo

package crabefi

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"slices"
	"testing"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/smbios"
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

func TestHandoffWithoutFramebuffer(t *testing.T) {
	f := testFirmware(t, 3)
	for _, tt := range []struct {
		fb   *framebufferInfo
		want []uint32
	}{
		{f.framebuffer, []uint32{tagMemory, tagRSDP, tagFramebuffer, tagCBMEMEntry}},
		{nil, []uint32{tagMemory, tagRSDP, tagCBMEMEntry}},
	} {
		data, err := makeHandoff(f.memory, 0x8004000, f.rsdp, tt.fb, f.smbiosEntry, f.smbiosData)
		if err != nil {
			t.Fatal(err)
		}
		var header tableHeader
		if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &header); err != nil {
			t.Fatal(err)
		}
		if checksum(data[:24]) != 0 || checksum(data[24:24+header.TableBytes]) != uint16(header.TableChecksum) {
			t.Fatalf("makeHandoff(framebuffer %v) produced invalid checksums", tt.fb != nil)
		}
		var tags []uint32
		r := bytes.NewReader(data[24 : 24+header.TableBytes])
		for r.Len() != 0 {
			var record recordHeader
			if err := binary.Read(r, binary.LittleEndian, &record); err != nil {
				t.Fatal(err)
			}
			tags = append(tags, record.Tag)
			if _, err := r.Seek(int64(record.Size-8), io.SeekCurrent); err != nil {
				t.Fatal(err)
			}
		}
		if header.TableEntries != uint32(len(tt.want)) || !slices.Equal(tags, tt.want) {
			t.Errorf("makeHandoff(framebuffer %v) records = %d, %#x; want %#x", tt.fb != nil, header.TableEntries, tags, tt.want)
		}
	}
}

func TestSMBIOS2TrimmedTable(t *testing.T) {
	f := testFirmware(t, 2)
	// Linux trims the exported table to its end-of-table structure but
	// exports the firmware's entry point unchanged.
	trimmed := f.smbiosData[:4]
	ep, err := relocateSMBIOS(f.smbiosEntry, trimmed, 0x8000000)
	if err != nil {
		t.Fatalf("relocateSMBIOS(trimmed SMBIOS 2) = %v; want nil", err)
	}
	e32, _, err := smbios.ParseEntry(ep)
	if err != nil {
		t.Fatalf("relocated entry point is invalid: %v", err)
	}
	if e32.StructTableLength != uint16(len(trimmed)) || e32.StructTableAddr != 0x8000000 {
		t.Errorf("relocated SMBIOS 2 table = %d bytes at %#x; want %d bytes at 0x8000000", e32.StructTableLength, e32.StructTableAddr, len(trimmed))
	}
}
