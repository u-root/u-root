// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package universalpayload supports to load FIT (Flat Image Tree) image.
// FIT is a common Payload image format to facilitate the loading process,
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
	kexecMemoryMapFromIOMem = kexec.MemoryMapFromIOMem
	getAcpiRSDP             = acpi.GetRSDP
	getSMBIOSBase           = smbios.SMBIOSBase
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
	RegisterBase   EFIPhysicalAddress
}

// Structure member 'Pad' is introduced to match the offset of 'Entry'
// in structure UNIVERSAL_PAYLOAD_BASE which is defined in EDK2 UPL.
type UniversalPayloadBase struct {
	Header UniversalPayloadGenericHeader
	Pad    [4]byte
	Entry  EFIPhysicalAddress
}

type UniversalPayloadAcpiTable struct {
	Header UniversalPayloadGenericHeader
	Rsdp   EFIPhysicalAddress
}

type UniversalPayloadSmbiosTable struct {
	Header           UniversalPayloadGenericHeader
	SmBiosEntryPoint EFIPhysicalAddress
}

// Map GUID string to size of corresponding structure. Use
// this map to simplify the length calculation in function
// constructGUIDHOB.
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
	ErrWriteHOBBufMemoryMap            = errors.New("failed to write memory map to buffer")
	ErrWriteHOBBufSerialPort           = errors.New("failed to append serial port hob to buffer")
	ErrWriteHOBBufUniversalPayloadBase = errors.New("failed to append universal payload base to buffer")
	ErrWriteHOBBufAcpiTable            = errors.New("failed to append acpi table to buffer")
	ErrWriteHOBSmbiosTable             = errors.New("failed to append smbios table to buffer")
	ErrWriteHOBEFICPU                  = errors.New("failed to append CPU HOB to buffer")
	ErrWriteHOBBufList                 = errors.New("failed to append HOB list to buffer")
	ErrWriteHOBLengthNotMatch          = errors.New("length mismatch when appending")
	ErrKexecLoadFailed                 = errors.New("kexec.Load() failed")
	ErrKexecExecuteFailed              = errors.New("kexec.Execute() failed")
	ErrMemMapIoMemExecuteFailed        = errors.New("failed to get memory from /proc/iomem")
)

// Create GUID HOB with specified GUID string
func constructGUIDHOB(name string) (*EFIHOBGUIDType, error) {
	length := uint16(unsafe.Sizeof(EFIHOBGUIDType{}) + guidToLength[name])

	id, err := guid.Parse(name)
	if err != nil {
		return nil, errors.Join(ErrParseGUIDFail, err)
	}

	return &EFIHOBGUIDType{
		Header: EFIHOBGenericHeader{
			HOBType:   EFIHOBTypeGUIDExtension,
			HOBLength: EFIHOBLength(length),
		},
		Name: id,
	}, nil
}

// Construct Serial Port HOB
func constructSerialPortHOB() *UniversalPayloadSerialPortInfo {
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

// Construct universal payload base HOB
func constructUniversalPayloadBase(addr uint64) *UniversalPayloadBase {
	return &UniversalPayloadBase{
		Header: UniversalPayloadGenericHeader{
			Revision: 0,
			Length:   uint16(unsafe.Sizeof(UniversalPayloadBase{})),
		},
		Entry: EFIPhysicalAddress(addr),
	}
}

// Construct UniversalPayloadAcpiTable HOB
func constructRSDPTable() (*UniversalPayloadAcpiTable, error) {
	rsdp, err := getAcpiRSDP()
	if err != nil {
		return nil, errors.Join(ErrFailToGetRSDPTable, err)
	}

	return &UniversalPayloadAcpiTable{
		Header: UniversalPayloadGenericHeader{
			Revision: UniversalPayloadAcpiTableRevision,
			Length:   uint16(unsafe.Sizeof(UniversalPayloadAcpiTable{})),
		},
		Rsdp: EFIPhysicalAddress(rsdp.RSDPAddr()),
	}, nil
}

// Construct UniversalPayloadSmbiosTable HOB
func constructSmbiosTable() (*UniversalPayloadSmbiosTable, error) {
	smbiosTableBase, _, err := getSMBIOSBase()
	if err != nil {
		return nil, errors.Join(ErrFailToGetSmbiosTable, err)
	}

	return &UniversalPayloadSmbiosTable{
		Header: UniversalPayloadGenericHeader{
			Revision: UniversalPayloadSmbiosTableRevision,
			Length:   uint16(unsafe.Sizeof(UniversalPayloadSmbiosTable{})),
		},
		SmBiosEntryPoint: EFIPhysicalAddress(smbiosTableBase),
	}, nil
}

// Construct system memory resource HOB
func appendMemMapHOB(buf *bytes.Buffer, hobLen *uint64, memMap kexec.MemoryMap) error {
	prev := buf.Len()
	memHOB, length := hobFromMemMap(memMap)
	if err := binary.Write(buf, binary.LittleEndian, memHOB); err != nil {
		return errors.Join(ErrWriteHOBBufMemoryMap, err)
	}

	if err := alignHOBLength(length, buf.Len()-prev, buf); err != nil {
		return errors.Join(ErrWriteHOBLengthNotMatch, err)
	}

	*hobLen += length

	return nil
}

// Construct serial port HOB
func appendSerialPortHOB(buf *bytes.Buffer, hobLen *uint64) error {
	serialPortInfo := constructSerialPortHOB()
	serialGUIDHOB, err := constructGUIDHOB(UniversalPayloadSerialPortInfoGUID)
	if err != nil {
		return err
	}

	length := uint64(unsafe.Sizeof(EFIHOBGUIDType{}) + unsafe.Sizeof(UniversalPayloadSerialPortInfo{}))
	prev := buf.Len()

	if err := binary.Write(buf, binary.LittleEndian, serialGUIDHOB); err != nil {
		return errors.Join(ErrWriteHOBBufSerialPort, err)
	}

	if err := binary.Write(buf, binary.LittleEndian, serialPortInfo); err != nil {
		return errors.Join(ErrWriteHOBBufSerialPort, err)
	}

	if err := alignHOBLength(length, buf.Len()-prev, buf); err != nil {
		return errors.Join(ErrWriteHOBLengthNotMatch, err)
	}

	*hobLen += length

	return nil
}

func appendUniversalPayloadBase(buf *bytes.Buffer, hobLen *uint64, load uint64) error {
	uplBase := constructUniversalPayloadBase(load)
	uplBaseGUIDHOB, err := constructGUIDHOB(UniversalPayloadBaseGUID)
	if err != nil {
		return err
	}

	length := uint64(unsafe.Sizeof(EFIHOBGUIDType{}) + unsafe.Sizeof(UniversalPayloadBase{}))
	prev := buf.Len()

	if err := binary.Write(buf, binary.LittleEndian, uplBaseGUIDHOB); err != nil {
		return errors.Join(ErrWriteHOBBufUniversalPayloadBase, err)
	}

	if err := binary.Write(buf, binary.LittleEndian, uplBase); err != nil {
		return errors.Join(ErrWriteHOBBufUniversalPayloadBase, err)
	}

	if err := alignHOBLength(length, buf.Len()-prev, buf); err != nil {
		return fmt.Errorf("%w, func = appendUniversalPayloadBase()", ErrWriteHOBLengthNotMatch)
	}

	*hobLen += length

	return nil
}

func appendAcpiTableHOB(buf *bytes.Buffer, hobLen *uint64) error {
	// Construct universal payload ACPI (RSDP) table HOB
	rsdpTable, err := constructRSDPTable()
	if err != nil {
		return err
	}

	rsdpTableGUIDHOB, err := constructGUIDHOB(UniversalPayloadAcpiTableGUID)
	if err != nil {
		return err
	}

	length := uint64(unsafe.Sizeof(EFIHOBGUIDType{}) + unsafe.Sizeof(UniversalPayloadAcpiTable{}))
	prev := buf.Len()

	if err := binary.Write(buf, binary.LittleEndian, rsdpTableGUIDHOB); err != nil {
		return errors.Join(ErrWriteHOBBufAcpiTable, err)
	}

	if err := binary.Write(buf, binary.LittleEndian, rsdpTable); err != nil {
		return errors.Join(ErrWriteHOBBufAcpiTable, err)
	}

	if err := alignHOBLength(length, buf.Len()-prev, buf); err != nil {
		return fmt.Errorf("%w, func = appendAcpiTableHOB()", ErrWriteHOBLengthNotMatch)
	}

	*hobLen += length

	return nil
}

func appendSmbiosTableHOB(buf *bytes.Buffer, hobLen *uint64) error {
	// Construct SMBIOS HOB
	smbiosTable, err := constructSmbiosTable()
	if err != nil {
		return err
	}

	smbiosTableGUIDHOB, err := constructGUIDHOB(UniversalPayloadSmbiosTableGUID)
	if err != nil {
		return err
	}

	length := uint64(unsafe.Sizeof(EFIHOBGUIDType{}) + unsafe.Sizeof(UniversalPayloadSmbiosTable{}))
	prev := buf.Len()

	if err := binary.Write(buf, binary.LittleEndian, smbiosTableGUIDHOB); err != nil {
		return errors.Join(ErrWriteHOBSmbiosTable, err)
	}

	if err := binary.Write(buf, binary.LittleEndian, smbiosTable); err != nil {
		return errors.Join(ErrWriteHOBSmbiosTable, err)
	}

	if err := alignHOBLength(length, buf.Len()-prev, buf); err != nil {
		return fmt.Errorf("%w, func = appendSmbiosTableHOB()", ErrWriteHOBLengthNotMatch)
	}

	*hobLen += length

	return nil
}

func appendEFICPUHOB(buf *bytes.Buffer, hobLen *uint64) error {
	cpuHOB, err := hobCreateEFIHOBCPU()
	if err != nil {
		return err
	}

	length := uint64(unsafe.Sizeof(EFIHOBCPU{}))
	prev := buf.Len()
	if err := binary.Write(buf, binary.LittleEndian, cpuHOB); err != nil {
		return errors.Join(ErrWriteHOBEFICPU, err)
	}

	if err := alignHOBLength(length, buf.Len()-prev, buf); err != nil {
		return fmt.Errorf("%w, func = appendEFICPUHOB()", ErrWriteHOBLengthNotMatch)
	}

	*hobLen += length

	return nil
}

func constructHOBList(dst *bytes.Buffer, src *bytes.Buffer, hobLen *uint64) error {
	handoffHOB := hobCreateEFIHOBHandoffInfoTable(*hobLen)
	if err := binary.Write(dst, binary.LittleEndian, handoffHOB); err != nil {
		return errors.Join(ErrWriteHOBBufList, err)
	}

	if err := binary.Write(dst, binary.LittleEndian, src.Bytes()); err != nil {
		return errors.Join(ErrWriteHOBBufList, err)
	}

	hobEndHeader := hobCreateEndHOB()
	prev := dst.Len()
	length := uint64(unsafe.Sizeof(EFIHOBGenericHeader{}))

	if err := binary.Write(dst, binary.LittleEndian, hobEndHeader); err != nil {
		return errors.Join(ErrWriteHOBBufList, err)
	}

	if length != (uint64)(dst.Len()-prev) {
		return fmt.Errorf("%w, func = constructHOBList()", ErrWriteHOBLengthNotMatch)
	}

	*hobLen += length

	return nil
}

func prepareBootEnv(hobAddr uint64, entry uint64, mem *kexec.Memory) error {
	stackBuffer := make([]byte, tmpStackSize)

	s := kexec.NewSegment(stackBuffer, kexec.Range{
		Start: uintptr(hobAddr + tmpHobSize),
		Size:  tmpStackSize,
	})
	mem.Segments.Insert(s)

	var trampoline []uint8
	trampoline = constructTrampoline(trampoline, hobAddr, entry)

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, trampoline)

	s = kexec.NewSegment(buf.Bytes(), kexec.Range{
		Start: uintptr(hobAddr + tmpEntryOffset),
		Size:  uint(buf.Len()),
	})

	mem.Segments.Insert(s)

	return nil
}

func prepareHob(buf *bytes.Buffer, length *uint64, addr uint64, mem *kexec.Memory) error {
	if err := appendMemMapHOB(buf, length, mem.Phys); err != nil {
		return err
	}

	if err := appendSerialPortHOB(buf, length); err != nil {
		return err
	}

	if err := appendUniversalPayloadBase(buf, length, addr); err != nil {
		return err
	}

	if err := appendAcpiTableHOB(buf, length); err != nil {
		return err
	}

	if err := appendSmbiosTableHOB(buf, length); err != nil {
		return err
	}

	if err := appendEFICPUHOB(buf, length); err != nil {
		return err
	}

	return nil
}

func prepareBootloaderParameter(fdtLoad *FdtLoad, loadAddr uint64, mem *kexec.Memory) error {
	hobBuf := &bytes.Buffer{}
	hobListBuf := &bytes.Buffer{}
	var hobLen uint64

	if err := prepareHob(hobBuf, &hobLen, fdtLoad.Load, mem); err != nil {
		return err
	}

	if err := constructHOBList(hobListBuf, hobBuf, &hobLen); err != nil {
		return err
	}

	s := kexec.NewSegment(hobListBuf.Bytes(), kexec.Range{
		Start: uintptr(loadAddr),
		Size:  uint(hobListBuf.Len()),
	})

	mem.Segments.Insert(s)

	return nil
}

func prepareFdtData(fdt *FdtLoad, data []byte, addr uint64, mem *kexec.Memory) error {
	if err := relocateFdtData(addr, fdt, data); err != nil {
		return err
	}

	s := kexec.NewSegment(data, kexec.Range{
		Start: uintptr(fdt.Load),
		Size:  uint(len(data)),
	})

	mem.Segments.Insert(s)
	return nil
}

func loadKexecMemWithHOBs(fdt *FdtLoad, data []byte, mem *kexec.Memory) (uintptr, error) {
	mmRanges := mem.Phys.RAM()

	fitOffset := tmpHobSize + tmpStackSize + trampolineSize
	rangeLen := fitOffset + len(data)

	// Try to find available Space to locate FIT image and HOB, stack and trampoline code.
	// 2MB alignment will be easy for target OS/Bootloader to construct page table.
	// The layout of this Space will be placed as following:
	//
	//  |------------------------|
	//  |       FIT IMAGE        |
	//  |------------------------|
	//  |     TRAMPOLINE CODE    |
	//  |------------------------|
	//  |      TEMP STACK        |
	//  |------------------------|
	//  |  BOOTLOADER PARAMETER  |
	//  |------------------------|
	//
	kernelRange, err := mmRanges.FindSpace(uint(rangeLen), kexec.WithAlignment(0x200000))
	if err != nil {
		return 0, err
	}

	targetAddr := uint64(kernelRange.Start)
	fitImgAddr := targetAddr + uint64(fitOffset)

	if err = prepareFdtData(fdt, data, fitImgAddr, mem); err != nil {
		return 0, err
	}

	if err = prepareBootloaderParameter(fdt, targetAddr, mem); err != nil {
		return 0, err
	}

	if err = prepareBootEnv(targetAddr, fdt.EntryStart, mem); err != nil {
		return 0, err
	}

	return (uintptr)(targetAddr + tmpEntryOffset), nil
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

	// Prepare memory.
	ioMem, err := kexecMemoryMapFromIOMem()
	if err != nil {
		return fmt.Errorf("%w: err: %w", ErrMemMapIoMemExecuteFailed, err)
	}

	mem := kexec.Memory{
		Phys: ioMem,
	}

	// Prepare boot environment, including HoB, stack, bootloader parameter.
	entry, err := loadKexecMemWithHOBs(fdtLoad, data, &mem)
	if err != nil {
		return err
	}

	if err := kexec.Load(entry, mem.Segments, 0); err != nil {
		return errors.Join(ErrKexecLoadFailed, err)
	}

	if err := boot.Execute(); err != nil {
		return errors.Join(ErrKexecExecuteFailed, err)
	}

	return nil
}
