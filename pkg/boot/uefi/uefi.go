// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uefi

import (
	"bytes"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/boot/kexec"
)

// TODO(chengchieh): Add integration test for uefi package.

const fvEntryImageOffset int64 = 0xA0

var kexecLoad = kexec.Load
var kexecParseMemoryMap = kexec.ParseMemoryMap
var getRSDP = acpi.GetRSDP

// Memory type in Linuxboot.h
// TODO(chengchieh): Merge them with pkg/boot/kexec/memory_linux.go
const (
	MemTypeRAM      = 1
	MemTypeDefault  = 2
	MemTypeAcpi     = 3
	MemTypeNVS      = 4
	MemTypeReserved = 5
)

type memoryMapEntry struct {
	Start   uint64
	End     uint64
	MemType uint32
}

var rangeTypeToPayloadMemType = map[kexec.RangeType]uint32{
	kexec.RangeRAM:      MemTypeRAM,
	kexec.RangeDefault:  MemTypeDefault,
	kexec.RangeACPI:     MemTypeAcpi,
	kexec.RangeNVS:      MemTypeNVS,
	kexec.RangeReserved: MemTypeReserved,
}


// SerialPortConfig defines debug port configuration
// This struct will be used to initialize SERIAL_PORT_INFO
// in payload (UefiPayloadPkg/Include/Guid/SerialPortInfoGuid.h)
type SerialPortConfig struct {
	Type        uint32
	BaseAddr    uint32
	Baud        uint32
	RegWidth    uint32
	InputHertz  uint32
	UartPciAddr uint32
}

// Serial port type for Type in UefiPayloadPkg's SERIAL_PORT_INFO
const (
	SerialPortTypeIO   = 1
	SerialPortTypeMMIO = 2
)

// Current Config Version: 1
const PayloadConfigVersion = 1

type payloadConfig struct {
	AcpiBase            uint64
	AcpiSize            uint64
	SerialConfig        SerialPortConfig
	NumMemoryMapEntries uint32
}

func convertToPayloadMemType(rt kexec.RangeType) uint32 {
	mt, ok := rangeTypeToPayloadMemType[rt]
	if !ok {
		// return reserved if range type is not recognized
		return MemTypeReserved
	}
	return mt
}

// FvImage is a structure for loading a firmware volume
type FvImage struct {
	name         string
	mem          kexec.Memory
	entryAddress uintptr
	ImageBase    uintptr
	SerialConfig SerialPortConfig
}

func checkFvAndGetEntryPoint(name string) (uintptr, error) {
	f, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	// Entry PE image for UEFI Payload is 0xA0 offset from FV head.
	pef := io.NewSectionReader(f, fvEntryImageOffset, math.MaxInt64)
	pf, err := pe.NewFile(pef)
	if err != nil {
		return 0, err
	}
	op64, ok := pf.OptionalHeader.(*pe.OptionalHeader64)
	if !ok {
		return 0, fmt.Errorf("it is not OptionalHeader64")
	}
	return uintptr(fvEntryImageOffset) + uintptr(op64.AddressOfEntryPoint), nil
}

// New loads the file and return FvImage stucture if entry image is found
func New(n string) (*FvImage, error) {
	entry, err := checkFvAndGetEntryPoint(n)
	if err != nil {
		return nil, err
	}
	return &FvImage{name: n, mem: kexec.Memory{}, entryAddress: entry}, nil
}

// Reserved 64kb for passing params
const uefiPayloadConfigSize = 0x10000

// Load loads fimware volume payload and boot the the payload
func (fv *FvImage) Load(verbose bool) error {
	// Install payload
	dat, err := ioutil.ReadFile(fv.name)
	if err != nil {
		return err
	}
	fv.mem.Segments.Insert(kexec.NewSegment(dat, kexec.Range{Start: fv.ImageBase, Size: uint(len(dat))}))

	// Install payload config & its memory map: 64 kb below the image
	// We cannot use the memory above the image base because it may be used by HOB
	var configAddr uintptr = fv.ImageBase - uintptr(uefiPayloadConfigSize)

	// Get MemoryMap
	mm, err := kexecParseMemoryMap()
	if err != nil {
		return err
	}

	// Get Acpi Basc (RSDP)
	rsdp, err := getRSDP()
	if err != nil {
		return err
	}

	pc := payloadConfig{
		AcpiBase:            uint64(rsdp.RSDPAddr()),
		AcpiSize:            uint64(rsdp.Len()),
		SerialConfig:        fv.SerialConfig,
		NumMemoryMapEntries: uint32(len(mm)),
	}

	pcbuf := &bytes.Buffer{}
	if err := binary.Write(pcbuf, binary.LittleEndian, pc); err != nil {
		return err
	}

	// Convert MemoryMap to UefiPayload style
	var mmPayload []memoryMapEntry
	for _, entry := range mm {
		mmPayload = append(mmPayload, memoryMapEntry{
			Start:   uint64(entry.Start),
			End:     uint64(entry.Start) + uint64(entry.Size) - 1,
			MemType: convertToPayloadMemType(entry.Type),
		})
	}

	if err := binary.Write(pcbuf, binary.LittleEndian, mmPayload); err != nil {
		return err
	}

	if len(pcbuf.Bytes()) > uefiPayloadConfigSize {
		return fmt.Errorf("Config/Memmap size is greater than reserved size: %d bytes", len(pcbuf.Bytes()))
	}

	fv.mem.Segments.Insert(kexec.NewSegment(pcbuf.Bytes(), kexec.Range{Start: configAddr, Size: uint(len(pcbuf.Bytes()))}))

	if verbose {
		log.Printf("segments cmdline %v", fv.mem.Segments)
	}

	if err := kexecLoad(fv.ImageBase+uintptr(fv.entryAddress), fv.mem.Segments, 0); err != nil {
		return fmt.Errorf("kexec.Load() error: %v", err)
	}

	return nil
}
