// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crabefi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/smbios"
)

const (
	// Keep the handoff away from the payload's usual 1 MiB load address.
	// This is a placement policy, not a coreboot ABI requirement.
	arenaMin         = 128 << 20
	arenaSize        = 1 << 20
	handoffOffset    = 0x1000   // one page of entry code; no caller stack needed
	mappedMemorySize = 64 << 30 // CrabEFI's initial identity map

	// Offsets relative to the coreboot table, not the start of the arena.
	smbiosEntryOffset = 0x1000
	smbiosDataOffset  = 0x1100
	pageSize          = 4096

	// CrabEFI's coreboot table parser uses a fixed-capacity memory map.
	maxMemoryRegions = 64

	// coreboot/src/commonlib/include/commonlib/coreboot_tables.h and
	// coreboot/src/commonlib/bsd/include/commonlib/bsd/cbmem_id.h.
	tagMemory      = 0x01
	tagFramebuffer = 0x12
	tagCBMEMEntry  = 0x31
	tagRSDP        = 0x43
	cbmemIDSMBIOS  = 0x534d4254 // "SMBT"

	memoryRAM      = 1
	memoryReserved = 2
	memoryACPI     = 3
	memoryNVS      = 4
)

// These structures are serialized with encoding/binary, without Go padding.
type tableHeader struct {
	Signature      [4]byte
	HeaderBytes    uint32
	HeaderChecksum uint32
	TableBytes     uint32
	TableChecksum  uint32
	TableEntries   uint32
}
type recordHeader struct{ Tag, Size uint32 }
type memoryRange struct {
	Base, Size uint64
	Kind       uint32
}
type cbmemEntry struct {
	Address  uint64
	Size, ID uint32
}
type colorChannel struct{ Position, Size uint8 }
type framebufferInfo struct {
	Address                    uint64
	Width, Height, Stride      uint32
	BitsPerPixel               uint8
	Red, Green, Blue, Reserved colorChannel
}

// checksum is coreboot's little-endian Internet checksum (ipchksum).
func checksum(data []byte) uint16 {
	var sum uint32
	for i, b := range data {
		sum += uint32(b) << (8 * (i % 2))
	}
	for sum>>16 != 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return ^uint16(sum)
}

func byteSum(data []byte) uint8 {
	var sum uint8
	for _, b := range data {
		sum += b
	}
	return sum
}

func relocateSMBIOS(entry, data []byte, address uint64) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty SMBIOS table")
	}
	e32, e64, err := smbios.ParseEntry(entry)
	if err != nil {
		return nil, err
	}
	if e64 != nil {
		if len(entry) != int(e64.Length) {
			return nil, fmt.Errorf("unsupported SMBIOS 3 entry length")
		}
		if uint64(len(data)) > uint64(e64.StructMaxSize) {
			return nil, fmt.Errorf("SMBIOS table exceeds declared size")
		}
		e64.StructMaxSize = uint32(len(data))
		e64.StructTableAddr = address
		return e64.MarshalBinary()
	}
	if len(entry) != int(e32.Length) {
		return nil, fmt.Errorf("unsupported SMBIOS 2 entry length")
	}
	if address > math.MaxUint32 || len(data) != int(e32.StructTableLength) {
		return nil, fmt.Errorf("invalid SMBIOS 2 address or table length")
	}
	e32.StructTableAddr = uint32(address)
	return e32.MarshalBinary()
}

func makeHandoff(memory kexec.MemoryMap, base, rsdp uint64, fb framebufferInfo, entry, data []byte) ([]byte, error) {
	if len(data) > arenaSize-handoffOffset-smbiosDataOffset {
		return nil, fmt.Errorf("SMBIOS exceeds handoff arena")
	}
	ep, err := relocateSMBIOS(entry, data, base+smbiosDataOffset)
	if err != nil {
		return nil, err
	}
	var ranges []memoryRange
	for _, r := range memory {
		kind := uint32(memoryReserved)
		switch r.Type {
		case kexec.RangeRAM:
			kind = memoryRAM
		case kexec.RangeACPI:
			kind = memoryACPI
		case kexec.RangeNVS:
			kind = memoryNVS
		}
		ranges = append(ranges, memoryRange{uint64(r.Start), uint64(r.Size), kind})
	}
	var table bytes.Buffer
	for _, record := range []struct {
		tag  uint32
		data any
	}{
		{tagMemory, ranges}, {tagRSDP, rsdp}, {tagFramebuffer, fb},
		{tagCBMEMEntry, cbmemEntry{base + smbiosEntryOffset, uint32(smbiosDataOffset + len(data) - smbiosEntryOffset), cbmemIDSMBIOS}},
	} {
		body, err := binary.Append(nil, binary.LittleEndian, record.data)
		if err != nil {
			return nil, err
		}
		size := (len(body) + 8 + 3) &^ 3
		if err := binary.Write(&table, binary.LittleEndian, recordHeader{record.tag, uint32(size)}); err != nil {
			return nil, err
		}
		table.Write(body)
		table.Write(make([]byte, size-8-len(body)))
	}
	header := tableHeader{
		Signature: [4]byte{'L', 'B', 'I', 'O'}, HeaderBytes: 24,
		TableBytes: uint32(table.Len()), TableChecksum: uint32(checksum(table.Bytes())),
		TableEntries: 4,
	}
	h, err := binary.Append(nil, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}
	header.HeaderChecksum = uint32(checksum(h))
	h, err = binary.Append(nil, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}
	if len(h)+table.Len() > smbiosEntryOffset || smbiosEntryOffset+len(ep) > smbiosDataOffset {
		return nil, fmt.Errorf("handoff records overlap SMBIOS")
	}
	result := make([]byte, smbiosDataOffset+len(data))
	copy(result, h)
	copy(result[len(h):], table.Bytes())
	copy(result[smbiosEntryOffset:], ep)
	copy(result[smbiosDataOffset:], data)
	return result, nil
}
