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
	"github.com/u-root/u-root/pkg/align"
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
	getSMBIOSBase           = smbios.SMBIOSBase
	getSMBIOS3HdrSize       = smbios.SMBIOS3HeaderSize
	getAcpiRsdpData         = archGetAcpiRsdpData
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
	ErrComponentsSizeOverflow          = errors.New("reserved components size overflow")
)

var (
	debug      = func(string, ...interface{}) {}
	warningMsg []error
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

func appendUniversalPayloadBase(buf *bytes.Buffer, hobLen *uint64, loadAddr uint64) error {
	uplBase := constructUniversalPayloadBase(loadAddr)
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

func checkComponentsSize(appendSize uint) error {
	componentsSize = componentsSize + appendSize

	debug("Current components size:%X vs. reserved size:%X\n", componentsSize, sizeForComponents)
	if componentsSize > uint(sizeForComponents) {
		return fmt.Errorf("components size check failure:%w", ErrComponentsSizeOverflow)
	}

	return nil
}

func prepareBootEnv(loadAddr uint64, entry uint64, mem *kexec.Memory) error {
	stackSize := pageSize
	stackBuffer := make([]byte, stackSize)

	// Check whether reserved components size is overflowed.
	if err := checkComponentsSize(stackSize); err != nil {
		return err
	}
	s := kexec.NewSegment(stackBuffer, kexec.Range{
		Start: uintptr(loadAddr + tmpStackOffset),
		Size:  stackSize,
	})
	mem.Segments.Insert(s)

	// Next step, trampoline code will be placed.
	trampolineOffset = tmpStackOffset + uint64(stackSize)

	var trampoline []uint8
	trampoline = constructTrampoline(trampoline, loadAddr, entry)

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, trampoline)

	// Check whether reserved components size is overflowed.
	if err := checkComponentsSize(uint(buf.Len())); err != nil {
		return err
	}
	s = kexec.NewSegment(buf.Bytes(), kexec.Range{
		Start: uintptr(loadAddr + trampolineOffset),
		Size:  uint(buf.Len()),
	})

	mem.Segments.Insert(s)

	return nil
}

func prepareHob(buf *bytes.Buffer, length *uint64, loadAddr uint64, mem *kexec.Memory) error {
	if err := appendMemMapHOB(buf, length, mem.Phys); err != nil {
		return err
	}

	if err := appendSerialPortHOB(buf, length); err != nil {
		return err
	}

	if err := appendUniversalPayloadBase(buf, length, loadAddr); err != nil {
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
	rsdpBase, rsdpData, err := getAcpiRsdpData()
	if err != nil {
		debug("universalpayload: failed to get RSDP table data (%v)\n", err)
		return err
	}

	// rsdpBase indicates whether we need to copy RSDP table data to specified
	// location. If rsdpBase equals to zero, then we need to copy data to
	// specified address, otherwise, we will use rsdpBase directly.
	if rsdpBase == 0 {
		// Check whether reserved components size is overflowed.
		if err := checkComponentsSize(align.UpPage(uint(len(rsdpData)))); err != nil {
			return err
		}
		s := kexec.NewSegment(rsdpData, kexec.Range{
			Start: uintptr(loadAddr + rsdpTableOffset),
			Size:  uint(len(rsdpData)),
		})

		mem.Segments.Insert(s)

		rsdpBase = loadAddr + rsdpTableOffset
	}

	// Next step, Handoff Blocks will be placed
	tmpHobOffset = rsdpTableOffset + uint64(align.UpPage(uint64(len(rsdpData))))

	hobBuf := &bytes.Buffer{}
	hobListBuf := &bytes.Buffer{}
	var hobLen uint64

	if err := prepareHob(hobBuf, &hobLen, fdtLoad.Load, mem); err != nil {
		debug("universalpayload: failed to construct HoBs (%v)\n", err)
		return err
	}

	if err := constructHOBList(hobListBuf, hobBuf, &hobLen); err != nil {
		debug("universalpayload: failed to construct HoBList (%v)\n", err)
		return err
	}

	// Check whether reserved components size is overflowed.
	if err := checkComponentsSize(align.UpPage(uint(hobListBuf.Len()))); err != nil {
		return err
	}
	s := kexec.NewSegment(hobListBuf.Bytes(), kexec.Range{
		Start: uintptr(loadAddr + tmpHobOffset),
		Size:  uint(hobListBuf.Len()),
	})

	mem.Segments.Insert(s)

	// Next step, FDT DTB info will be placed
	fdtDtbOffset = tmpHobOffset + uint64(align.UpPage(uint64(hobListBuf.Len())))

	dtBuf := &bytes.Buffer{}

	err = buildDeviceTreeInfo(dtBuf, mem, loadAddr, rsdpBase)
	if err != nil {
		debug("universalpayload: failed to build FDT (%v)\n", err)
		return err
	}

	// Check whether reserved components size is overflowed.
	if err := checkComponentsSize(align.UpPage(uint(dtBuf.Len()))); err != nil {
		return err
	}
	s = kexec.NewSegment(dtBuf.Bytes(), kexec.Range{
		Start: uintptr(loadAddr + fdtDtbOffset),
		Size:  uint(dtBuf.Len()),
	})
	mem.Segments.Insert(s)

	// Next step, temporary stack for trampoline code will be placed
	tmpStackOffset = fdtDtbOffset + uint64(align.UpPage(uint64(dtBuf.Len())))

	return nil
}

func prepareFdtData(fdt *FdtLoad, data []byte, addr uint64, mem *kexec.Memory) error {
	if err := relocateFdtData(addr+uplImageOffset, fdt, data); err != nil {
		debug("universalpayload: failed to relocate FIT image (%v)\n", err)
		return err
	}

	s := kexec.NewSegment(data, kexec.Range{
		Start: uintptr(fdt.Load),
		Size:  uint(len(data)),
	})

	mem.Segments.Insert(s)

	// Next step, ACPI RSDP table content will be placed
	rsdpTableOffset = uplImageOffset + uint64(align.UpPage(uint64(len(data))))

	return nil
}

func loadKexecMemWithHOBs(fdt *FdtLoad, data []byte, mem *kexec.Memory) (uintptr, error) {
	mmRanges := mem.Phys.RAM()

	// Reserved 1MB additional space which is used to place Device Tree info, Handoff Blocks,
	// temporary stack and trampoline code.
	rangeLen := len(data) + int(sizeForComponents)

	// Try to find available Space to locate FIT image and HOB, stack and trampoline code,
	// Device Tree information, and ACPI DATA.
	// 2MB alignment will be easy for target OS/Bootloader to construct page table.
	// The layout of this Space will be placed as following:
	//
	//  |------------------------|  <-- Memory Region top
	//  |     TRAMPOLINE CODE    |
	//  |------------------------|
	//  |      TEMP STACK        |
	//  |------------------------|
	//  |    Device Tree Info    |
	//  |------------------------|
	//  |  BOOTLOADER PARAMETER  |
	//  |------------------------|
	//  |       ACPI DATA        |
	//  |------------------------|
	//  |       FIT IMAGE        |
	//  |------------------------|  <-- Memory Region bottom
	//
	kernelRange, err := mmRanges.FindSpace(uint(rangeLen), kexec.WithAlignment(uplImageAlignment))
	if err != nil {
		debug("universalpayload: failed to find 2MB aligned space (%v)\n", err)
		return 0, err
	}

	loadAddr := uint64(kernelRange.Start)

	if err = prepareFdtData(fdt, data, loadAddr, mem); err != nil {
		debug("universalpayload: failed to prepare FDT data (%v)\n", err)
		return 0, err
	}

	if err = prepareBootloaderParameter(fdt, loadAddr, mem); err != nil {
		debug("universalpayload: failed to prepare boot parameters (%v)\n", err)
		return 0, err
	}

	if err = prepareBootEnv(loadAddr, fdt.EntryStart, mem); err != nil {
		return 0, err
	}

	return (uintptr)(loadAddr + uint64(trampolineOffset)), nil
}

func Load(name string, dbg func(string, ...interface{})) (error, error) {
	if dbg != nil {
		debug = dbg
	}

	debug("universalpayload: Try to get FDT information from:%s\n", name)
	fdtLoad, err := GetFdtInfo(name)
	if err != nil {
		debug("universalpayload: Failed to get FDT information (%v)\n", err)
		return err, errors.Join(warningMsg...)
	}

	debug("universalpayload: Try to fetch file content\n")
	data, err := os.ReadFile(name)
	if err != nil {
		debug("universalpayload: Failed to fetch file content (%v)\n", err)
		return fmt.Errorf("%w: file: %s, err: %w", ErrFailToReadFdtFile, name, err), errors.Join(warningMsg...)
	}

	// Prepare memory.
	debug("universalpayload: Try to get Memory Map from IOMem\n")
	ioMem, err := kexecMemoryMapFromIOMem()
	if err != nil {
		debug("universalpayload: Failed to get Memory Map from IOMem\n")
		return fmt.Errorf("%w: err: %w", ErrMemMapIoMemExecuteFailed, err), errors.Join(warningMsg...)
	}

	mem := kexec.Memory{
		Phys: ioMem,
	}

	// Prepare boot environment, including HoB, stack, bootloader parameter.
	debug("universalpayload: Try to prepare required stuffs\n")
	entry, err := loadKexecMemWithHOBs(fdtLoad, data, &mem)
	if err != nil {
		debug("universalpayload: Failed to prepare parameters with error (%v)\n", err)
		return err, errors.Join(warningMsg...)
	}

	debug("universalpayload: Entry:%x, Segments:%v\n", entry, mem.Segments)
	if err := kexec.Load(entry, mem.Segments, 0); err != nil {
		debug("universalpayload: Failed to load segments with error (%v)\n", err)
		return errors.Join(ErrKexecLoadFailed, err), errors.Join(warningMsg...)
	}

	debug("universalpayload: boot trampoline code at:%x\n", entry)

	return nil, errors.Join(warningMsg...)
}

func Exec() error {
	if err := boot.Execute(); err != nil {
		return errors.Join(ErrKexecExecuteFailed, err)
	}

	return nil
}
