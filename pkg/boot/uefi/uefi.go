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
	"log"
	"math"
	"os"

	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/smbios"
)

var (
	kexecLoad                     = kexec.Load
	kexecMemoryMapFromSysfsMemmap = kexec.MemoryMapFromSysfsMemmap
	getRSDP                       = acpi.GetRSDP
	getSMBIOSBase                 = smbios.SMBIOSBase
)

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
	Version             uint64
	ACPIBase            uint64
	ACPISize            uint64
	SMBIOSBase          uint64
	SMBIOSSize          uint64
	SerialConfig        SerialPortConfig
	NumMemoryMapEntries uint32
}

// FVImage is a structure for loading a firmware volume
type FVImage struct {
	name         string
	mem          kexec.Memory
	entryAddress uintptr
	ImageBase    uintptr
	SerialConfig SerialPortConfig
}

func checkFVAndGetEntryPoint(name string) (uintptr, error) {
	// Parse entire firmware volume to find SEC's PE32. Since payload will
	// be only few MBs, it should be fine to load entire image for parsing.
	dat, err := os.ReadFile(name)
	if err != nil {
		return 0, err
	}

	secEntry, err := findSecurityCorePEEntry(dat)
	if err != nil {
		return 0, err
	}

	f, err := os.Open(name)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	pef := io.NewSectionReader(f, int64(secEntry), math.MaxInt64)
	pf, err := pe.NewFile(pef)
	if err != nil {
		return 0, err
	}
	op64, ok := pf.OptionalHeader.(*pe.OptionalHeader64)
	if !ok {
		return 0, fmt.Errorf("it is not OptionalHeader64")
	}
	return uintptr(secEntry) + uintptr(op64.AddressOfEntryPoint), nil
}

// New loads the file and return FVImage stucture if entry image is found
func New(n string) (*FVImage, error) {
	entry, err := checkFVAndGetEntryPoint(n)
	if err != nil {
		return nil, err
	}
	return &FVImage{name: n, mem: kexec.Memory{}, entryAddress: entry}, nil
}

// Reserved 64kb for passing params
const uefiPayloadConfigSize = 0x10000

// Load loads fimware volume payload and boot the the payload
func (fv *FVImage) Load(verbose bool) error {
	// Install payload
	dat, err := os.ReadFile(fv.name)
	if err != nil {
		return err
	}
	fv.mem.Segments.Insert(kexec.NewSegment(dat, kexec.Range{Start: fv.ImageBase, Size: uint(len(dat))}))

	// Install payload config & its memory map: 64 kb below the image
	// We cannot use the memory above the image base because it may be used by HOB
	configAddr := fv.ImageBase - uintptr(uefiPayloadConfigSize)

	// Get MemoryMap
	mm, err := kexecMemoryMapFromSysfsMemmap()
	if err != nil {
		return err
	}

	// Get Acpi Basc (RSDP)
	rsdp, err := getRSDP()
	if err != nil {
		return err
	}

	smbiosBase, smbiosSize, err := getSMBIOSBase()
	if err != nil {
		return err
	}

	pc := payloadConfig{
		Version:             PayloadConfigVersion,
		ACPIBase:            uint64(rsdp.RSDPAddr()),
		ACPISize:            uint64(rsdp.Len()),
		SMBIOSBase:          uint64(smbiosBase),
		SMBIOSSize:          uint64(smbiosSize),
		SerialConfig:        fv.SerialConfig,
		NumMemoryMapEntries: uint32(len(mm)),
	}

	pcbuf := &bytes.Buffer{}
	if err := binary.Write(pcbuf, binary.LittleEndian, pc); err != nil {
		return err
	}

	if err := binary.Write(pcbuf, binary.LittleEndian, mm.ToUEFIPayloadMemoryMap()); err != nil {
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
		return fmt.Errorf("kexec.Load() error: %w", err)
	}

	return nil
}
