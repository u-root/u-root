// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package universalpayload supports to load FIT (Flat Image Tree) image.
// FIT is a common Payload image format to faciliate the loading process,
// and defined in UniversalPayload Specification.
// More Details about UniversalPayload Specification, please refer:
// https://github.com/universalpayload/spec
package universalpayload

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"unsafe"

	guid "github.com/google/uuid"
	"github.com/u-root/u-root/pkg/acpi"
	"github.com/u-root/u-root/pkg/boot"
	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/smbios"
)

const (
	UniversalPayloadSerialPortInfoGUID       = "0d197eaa-21be-0944-8e67-a2cd0f61e170"
	UniversalPayloadSerialPortInfoRevision   = 1
	UniversalPayloadSerialPortRegisterStride = 1
	UniversalPayloadSerialPortBaudRate       = 115200
	UniversalPayloadSerialPortRegisterBase   = 0x3f8
)

const (
	UniversalPayloadBaseGUID = "1dc6d403-1327-c54e-a1cc-883be9dc18e5"
)

const (
	UniversalPayloadAcpiTableGUID     = "06959a9f-9755-1545-bab6-8bcde784ba87"
	UniversalPayloadAcpiTableRevision = 1
)

const (
	UniversalPayloadSmbiosTableGUID     = "260d0a59-e506-204d-8a82-59ea1b34982d"
	UniversalPayloadSmbiosTableRevision = 1
)

var (
	kexecMemoryMapFromSysfsMemmap = kexec.MemoryMapFromSysfsMemmap
	acpiGetRSDP                   = acpi.GetRSDP
	smbiosSMBIOSBase              = smbios.SMBIOSBase
)

type UniversalPayloadGenericHeader struct {
	Revision uint8
	Reserved uint8
	Length   uint16
}

type UniversalPayloadSerialPortInfo struct {
	Header         UniversalPayloadGenericHeader
	UseMmio        uint8
	RegisterStride uint8
	BaudRate       uint32
	RegisterBase   EfiPhysicalAddress
}

// Structure member 'Pad' is introduced to match the offset of 'Entry'
// in structure UNIVERSAL_PAYLOAD_BASE which is defined in EDK2 UPL.
type UniversalPayloadBase struct {
	Header UniversalPayloadGenericHeader
	Pad    [4]byte
	Entry  EfiPhysicalAddress
}

type UniversalPayloadAcpiTable struct {
	Header UniversalPayloadGenericHeader
	Rsdp   EfiPhysicalAddress
}

type UniversalPayloadSmbiosTable struct {
	Header           UniversalPayloadGenericHeader
	SmBiosEntryPoint EfiPhysicalAddress
}

// Map GUID string to size of corresponding structure. Use
// this map to simplify the length calculation in function
// constructGUIDHob.
var (
	guidToLength = map[string]uintptr{
		UniversalPayloadSerialPortInfoGUID: unsafe.Sizeof(UniversalPayloadSerialPortInfo{}),
		UniversalPayloadBaseGUID:           unsafe.Sizeof(UniversalPayloadBase{}),
		UniversalPayloadAcpiTableGUID:      unsafe.Sizeof(UniversalPayloadAcpiTable{}),
		UniversalPayloadSmbiosTableGUID:    unsafe.Sizeof(UniversalPayloadSmbiosTable{}),
	}
)

var (
	ErrParseGUIDFail                   = errors.New("failed to parse GUID")
	ErrFailToGetRSDPTable              = errors.New("failed to get RSDP table")
	ErrFailToGetSmbiosTable            = errors.New("failed to get smbios base")
	ErrReadGetMemoryMap                = errors.New("failed to get memory map from sysfs")
	ErrWriteHobBufMemoryMap            = errors.New("failed to write memory map to buffer")
	ErrWriteHobBufSerialPort           = errors.New("failed to append serial port hob to buffer")
	ErrWriteHobBufUniversalPayloadBase = errors.New("failed to append universal payload base to buffer")
	ErrWriteHobBufAcpiTable            = errors.New("failed to append acpi table to buffer")
	ErrWriteHobSmbiosTable             = errors.New("failed to append smbios table to buffer")
	ErrWriteHobEfiCPU                  = errors.New("failed to append CPU HOB to buffer")
	ErrWriteHobBufList                 = errors.New("failed to append HOB list to buffer")
	ErrWriteHobLengthNotMatch          = errors.New("length mismatch when appending")
	ErrKexecLoadFailed                 = errors.New("kexec.Load() failed")
	ErrKexecExecuteFailed              = errors.New("kexec.Execute() failed")
)

// Create GUID Hob with specified GUID string
func constructGUIDHob(name string) (*EfiHobGUIDType, error) {
	length := uint16(unsafe.Sizeof(EfiHobGUIDType{}) + guidToLength[name])

	id, err := guid.Parse(name)
	if err != nil {
		return nil, errors.Join(ErrParseGUIDFail, err)
	}

	return &EfiHobGUIDType{
		Header: EfiHobGenericHeader{
			HobType:   EfiHobTypeGUIDExtension,
			HobLength: length,
		},
		Name: id,
	}, nil
}

func constructSerialPortHob() *UniversalPayloadSerialPortInfo {
	// Construct Serial Port Hob
	return &UniversalPayloadSerialPortInfo{
		Header: UniversalPayloadGenericHeader{
			Revision: UniversalPayloadSerialPortInfoRevision,
			Length:   uint16(unsafe.Sizeof(UniversalPayloadSerialPortInfo{})),
		},
		UseMmio:        0,
		RegisterStride: UniversalPayloadSerialPortRegisterStride,
		BaudRate:       UniversalPayloadSerialPortBaudRate,
		RegisterBase:   UniversalPayloadSerialPortRegisterBase,
	}
}

func constructUnivesralPayloadBase(addr uint64) *UniversalPayloadBase {
	return &UniversalPayloadBase{
		Header: UniversalPayloadGenericHeader{
			Revision: 0,
			Length:   uint16(unsafe.Sizeof(UniversalPayloadBase{})),
		},
		Entry: EfiPhysicalAddress(addr),
	}
}

func constructRSDPTable() (*UniversalPayloadAcpiTable, error) {
	rsdp, err := acpiGetRSDP()
	if err != nil {
		return nil, errors.Join(ErrFailToGetRSDPTable, err)
	}

	return &UniversalPayloadAcpiTable{
		Header: UniversalPayloadGenericHeader{
			Revision: UniversalPayloadAcpiTableRevision,
			Length:   uint16(unsafe.Sizeof(UniversalPayloadAcpiTable{})),
		},
		Rsdp: EfiPhysicalAddress(rsdp.RSDPAddr()),
	}, nil
}

func constructSmbiosTable() (*UniversalPayloadSmbiosTable, error) {
	smbiosTableBase, _, err := smbiosSMBIOSBase()
	if err != nil {
		return nil, errors.Join(ErrFailToGetSmbiosTable, err)
	}

	return &UniversalPayloadSmbiosTable{
		Header: UniversalPayloadGenericHeader{
			Revision: UniversalPayloadSmbiosTableRevision,
			Length:   uint16(unsafe.Sizeof(UniversalPayloadSmbiosTable{})),
		},
		SmBiosEntryPoint: EfiPhysicalAddress(smbiosTableBase),
	}, nil
}

func appendMemMapHob(buf *bytes.Buffer, hobLen *uint64) error {
	// Construct system memory resource Hob
	memMap, err := kexecMemoryMapFromSysfsMemmap()
	if err != nil {
		return errors.Join(ErrReadGetMemoryMap, err)
	}
	prev := buf.Len()
	memHob, length := hobFromMemMap(memMap)
	if err := binary.Write(buf, binary.LittleEndian, memHob); err != nil {
		return errors.Join(ErrWriteHobBufMemoryMap, err)
	}

	if err := alignHOBLength(length, buf.Len()-prev, buf); err != nil {
		return errors.Join(ErrWriteHobLengthNotMatch, err)
	}

	*hobLen += length

	return nil
}

func appendSerialPortHob(buf *bytes.Buffer, hobLen *uint64) error {
	// Construct serial port Hob
	serialPortInfo := constructSerialPortHob()
	serialGUIDHob, err := constructGUIDHob(UniversalPayloadSerialPortInfoGUID)
	if err != nil {
		return err
	}

	length := uint64(unsafe.Sizeof(EfiHobGUIDType{}) + unsafe.Sizeof(UniversalPayloadSerialPortInfo{}))
	prev := buf.Len()

	if err := binary.Write(buf, binary.LittleEndian, serialGUIDHob); err != nil {
		return errors.Join(ErrWriteHobBufSerialPort, err)
	}

	if err := binary.Write(buf, binary.LittleEndian, serialPortInfo); err != nil {
		return errors.Join(ErrWriteHobBufSerialPort, err)
	}

	if err := alignHOBLength(length, buf.Len()-prev, buf); err != nil {
		return errors.Join(ErrWriteHobLengthNotMatch, err)
	}

	*hobLen += length

	return nil
}

func appendUniversalPayloadBase(buf *bytes.Buffer, hobLen *uint64, load uint64) error {
	// Construct universal payload base Hob
	uplBase := constructUnivesralPayloadBase(load)
	uplBaseGUIDHob, err := constructGUIDHob(UniversalPayloadBaseGUID)
	if err != nil {
		return err
	}

	length := uint64(unsafe.Sizeof(EfiHobGUIDType{}) + unsafe.Sizeof(UniversalPayloadBase{}))
	prev := buf.Len()

	if err := binary.Write(buf, binary.LittleEndian, uplBaseGUIDHob); err != nil {
		return errors.Join(ErrWriteHobBufUniversalPayloadBase, err)
	}

	if err := binary.Write(buf, binary.LittleEndian, uplBase); err != nil {
		return errors.Join(ErrWriteHobBufUniversalPayloadBase, err)
	}

	if err := alignHOBLength(length, buf.Len()-prev, buf); err != nil {
		return fmt.Errorf("%w, func = appendUniversalPayloadBase()", ErrWriteHobLengthNotMatch)
	}

	*hobLen += length

	return nil
}

func appendAcpiTableHob(buf *bytes.Buffer, hobLen *uint64) error {
	// Construct universal payload ACPI (RSDP) table Hob
	rsdpTable, err := constructRSDPTable()
	if err != nil {
		return err
	}

	rsdpTableGUIDHob, err := constructGUIDHob(UniversalPayloadAcpiTableGUID)
	if err != nil {
		return err
	}

	length := uint64(unsafe.Sizeof(EfiHobGUIDType{}) + unsafe.Sizeof(UniversalPayloadAcpiTable{}))
	prev := buf.Len()

	if err := binary.Write(buf, binary.LittleEndian, rsdpTableGUIDHob); err != nil {
		return errors.Join(ErrWriteHobBufAcpiTable, err)
	}

	if err := binary.Write(buf, binary.LittleEndian, rsdpTable); err != nil {
		return errors.Join(ErrWriteHobBufAcpiTable, err)
	}

	if err := alignHOBLength(length, buf.Len()-prev, buf); err != nil {
		return fmt.Errorf("%w, func = appendAcpiTableHOB()", ErrWriteHobLengthNotMatch)
	}

	*hobLen += length

	return nil
}

func appendSmbiosTableHob(buf *bytes.Buffer, hobLen *uint64) error {
	// Construct SMBIOS Hob
	smbiosTable, err := constructSmbiosTable()
	if err != nil {
		return err
	}

	smbiosTableGUIDHob, err := constructGUIDHob(UniversalPayloadSmbiosTableGUID)
	if err != nil {
		return err
	}

	length := uint64(unsafe.Sizeof(EfiHobGUIDType{}) + unsafe.Sizeof(UniversalPayloadSmbiosTable{}))
	prev := buf.Len()

	if err := binary.Write(buf, binary.LittleEndian, smbiosTableGUIDHob); err != nil {
		return errors.Join(ErrWriteHobSmbiosTable, err)
	}

	if err := binary.Write(buf, binary.LittleEndian, smbiosTable); err != nil {
		return errors.Join(ErrWriteHobSmbiosTable, err)
	}

	if err := alignHOBLength(length, buf.Len()-prev, buf); err != nil {
		return fmt.Errorf("%w, func = appendSmbiosTableHOB()", ErrWriteHobLengthNotMatch)
	}

	*hobLen += length

	return nil
}

func appendEfiCPUHob(buf *bytes.Buffer, hobLen *uint64) error {
	cpuHob, err := hobCreateEfiHobCPU()
	if err != nil {
		return err
	}

	length := uint64(unsafe.Sizeof(EfiHobCPU{}))
	prev := buf.Len()
	if err := binary.Write(buf, binary.LittleEndian, cpuHob); err != nil {
		return errors.Join(ErrWriteHobEfiCPU, err)
	}

	if err := alignHOBLength(length, buf.Len()-prev, buf); err != nil {
		return fmt.Errorf("%w, func = appendEFICPUHOB()", ErrWriteHobLengthNotMatch)
	}

	*hobLen += length

	return nil
}

func constructHobList(dst *bytes.Buffer, src *bytes.Buffer, hobLen *uint64) error {
	handoffHob := hobCreateEfiHobHandoffInfoTable(*hobLen)
	if err := binary.Write(dst, binary.LittleEndian, handoffHob); err != nil {
		return errors.Join(ErrWriteHobBufList, err)
	}

	if err := binary.Write(dst, binary.LittleEndian, src.Bytes()); err != nil {
		return errors.Join(ErrWriteHobBufList, err)
	}

	hobEndHeader := hobCreateEndHob()
	prev := dst.Len()
	length := uint64(unsafe.Sizeof(EfiHobGenericHeader{}))

	if err := binary.Write(dst, binary.LittleEndian, hobEndHeader); err != nil {
		return errors.Join(ErrWriteHobBufList, err)
	}

	if length != (uint64)(dst.Len()-prev) {
		return fmt.Errorf("%w, func = constructHOBList()", ErrWriteHobLengthNotMatch)
	}

	*hobLen += length

	return nil
}

func loadKexecMemWithHobs(fdtLoad *FdtLoad, fdtData []byte) (kexec.Memory, error) {
	//Step 1, Prepare memory
	mem := kexec.Memory{}

	//Step 2, Insert tianocore raw binary
	mem.Segments.Insert(kexec.NewSegment(fdtData, kexec.Range{Start: uintptr(fdtLoad.Load), Size: uint(len(fdtData))}))

	// Step 3, Prepare HobList
	// TODO: remove hardcode HoB Address here
	hobAddr := fdtLoad.Load - 0x100000
	hobBuf := &bytes.Buffer{}
	hobListBuf := &bytes.Buffer{}
	var hobLen uint64

	if err := appendMemMapHob(hobBuf, &hobLen); err != nil {
		return mem, err
	}

	if err := appendSerialPortHob(hobBuf, &hobLen); err != nil {
		return mem, err
	}

	if err := appendUniversalPayloadBase(hobBuf, &hobLen, fdtLoad.Load); err != nil {
		return mem, err
	}

	if err := appendAcpiTableHob(hobBuf, &hobLen); err != nil {
		return mem, err
	}

	if err := appendSmbiosTableHob(hobBuf, &hobLen); err != nil {
		return mem, err
	}

	if err := appendEfiCPUHob(hobBuf, &hobLen); err != nil {
		return mem, err
	}

	if err := constructHobList(hobListBuf, hobBuf, &hobLen); err != nil {
		return mem, err
	}

	mem.Segments.Insert(kexec.NewSegment(hobListBuf.Bytes(), kexec.Range{Start: uintptr(hobAddr), Size: uint(hobLen)}))
	return mem, nil
}

func Load(name string) error {
	fdtLoad, err := GetFdtInfo(name)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(name)
	if err != nil {
		return fmt.Errorf("%w: file: %s, err: %w", ErrFailToReadFdtFile, name, err)
	}

	// Prepare memory
	mem, err := loadKexecMemWithHobs(fdtLoad, data)

	if err != nil {
		return err
	}

	if err := kexec.Load(uintptr(fdtLoad.EntryStart), mem.Segments, 0); err != nil {
		return errors.Join(ErrKexecLoadFailed, err)
	}

	if err := boot.Execute(); err != nil {
		return errors.Join(ErrKexecExecuteFailed, err)
	}

	return nil
}
