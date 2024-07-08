// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package universalpayload

import (
	"strings"
	"unsafe"

	guid "github.com/google/uuid"
	"github.com/u-root/u-root/pkg/boot/kexec"
)

// Efi HOB Handoff Info Table
const EfiHobHandoffTableVersion = 0x0009
const EfiHobAlignment = 0x100000

// Efi HOB Types
const (
	EfiHobTypeHandoff            = 0x0001
	EfiHobTypeMemoryAllocation   = 0x0002
	EfiHobTypeResourceDescriptor = 0x0003
	EfiHobTypeGUIDExtension      = 0x0004
	EfiHobTypeFv                 = 0x0005
	EfiHobTypeCPU                = 0x0006
	EfiHobTypeMemoryPool         = 0x0007
	EfiHobTypeFv2                = 0x0009
	EfiHobTypeLoadPeimUnused     = 0x000A
	EfiHobTypeUefiCapsule        = 0x000B
	EfiHobTypeFv3                = 0x000C
	EfiHobTypeUnused             = 0xFFFE
	EfiHobTypeEndOfHobList       = 0xFFFF
)

// Efi Resource Types
const (
	EfiResourceSystemMemory       = 0x00000000
	EfiResourceMemoryMappedIO     = 0x00000001
	EfiResourceIO                 = 0x00000002
	EfiResourceFirmwareDevice     = 0x00000003
	EfiResourceMemoryMappedIOPort = 0x00000004
	EfiResourceMemoryReserved     = 0x00000005
	EfiResourceIOReserved         = 0x00000006
	EfiResourceMemoryUnaccepted   = 0x00000007
	EfiResourceMaxMemoryType      = 0x00000008
	EfiResourceUnimplemeted       = 0x00000008
)

// Efi Resource Attributes
const (
	EfiResourceAttributePresent               = 0x00000001
	EfiResourceAttributeInitialized           = 0x00000002
	EfiResourceAttributeTested                = 0x00000004
	EfiResourceAttributeReadProtected         = 0x00000080
	EfiResourceAttributeWriteProtected        = 0x00000100
	EfiResourceAttributeExecutionProtected    = 0x00000200
	EfiResourceAttributePersistent            = 0x00800000
	EfiResourceAttributeSingleBitECC          = 0x00000008
	EfiResourceAttributeMultipleBitECC        = 0x00000010
	EfiResourceAttributeECCReserved1          = 0x00000020
	EfiResourceAttributeECCReserved2          = 0x00000040
	EfiResourceAttributeUncacheable           = 0x00000400
	EfiResourceAttributeWriteCombineable      = 0x00000800
	EfiResourceAttributeWriteThroughCacheable = 0x00001000
	EfiResourceAttributeWriteBackCacheable    = 0x00002000
	EfiResourceAttribute16BitIO               = 0x00004000
	EfiResourceAttribute32BitIO               = 0x00008000
	EfiResourceAttribute64BitIO               = 0x00010000
	EfiResourceAttributeUncachedExported      = 0x00020000
	EfiResourceAttributeReadProtectable       = 0x00100000
	EfiResourceAttributeWriteProtectable      = 0x00200000
	EfiResourceAttributeExecutionProtectable  = 0x00400000
	EfiResourceAttributePersistable           = 0x01000000
	EfiResourceAttributeReadOnlyProtected     = 0x00040000
	EfiResourceAttributeReadOnlyProtectable   = 0x00080000
	EfiResourceAttributeEncrypted             = 0x04000000
	EfiResourceAttributeSpecialPurpose        = 0x08000000
	EfiResourceAttributeMoreReliable          = 0x02000000
)

// HandoffHob values
const (
	EfiHobHandoffInfoVersion          = 0x09
	EfiHobHandoffInfoBootMode         = 0x0
	EfiHobHandoffInfoEfiMemoryTop     = 0x4E00000
	EfiHobHandoffInfoEfiMemoryBottom  = 0x800000
	EfiHobHandoffInfoFreeEfiMemoryTop = 0x4E00000
	EfiHobHandoffInfoFreeMemoryBottom = 0xE00000
)

// Type Aliases
type EfiBootMode uint32
type EfiResourceType uint32
type EfiPhysicalAddress uint64
type EfiResourceAttributeType uint32

// Efi HOB Generic Header
type EfiHobGenericHeader struct {
	HobType   uint16
	HobLength uint16
	Reserved  uint32
}

type EfiHobHandoffInfoTable struct {
	Header              EfiHobGenericHeader
	Version             uint32
	BootMode            EfiBootMode
	EfiMemoryTop        EfiPhysicalAddress
	EfiMemoryBottom     EfiPhysicalAddress
	EfiFreeMemoryTop    EfiPhysicalAddress
	EfiFreeMemoryBottom EfiPhysicalAddress
	EfiEndOfHobList     EfiPhysicalAddress
}

// Efi HOB Resource Descriptor
type EfiHobResourceDescriptor struct {
	Header            EfiHobGenericHeader
	Owner             [16]byte // EfiGUID (assuming a byte array for simplicity)
	ResourceType      EfiResourceType
	ResourceAttribute EfiResourceAttributeType
	PhysicalStart     EfiPhysicalAddress
	ResourceLength    uint64
}

// Efi HOB Firmware Volume
type EfiHobFirmwareVolume struct {
	Header      EfiHobGenericHeader
	BaseAddress EfiPhysicalAddress
	Length      uint64
}

// Efi HOB Firmware Volume
type EfiHobGUIDType struct {
	Header EfiHobGenericHeader
	Name   guid.UUID
}

// Efi HOB CPU
type EfiHobCPU struct {
	Header            EfiHobGenericHeader
	SizeOfMemorySpace uint8
	SizeOfIOSpace     uint8
	Reserved          [6]byte
}

type EfiMemoryMapHob []EfiHobResourceDescriptor

// Translate System Map with "System RAM" type to Resource code Hobs.
func hobFromMemMap(memMap kexec.MemoryMap) (EfiMemoryMapHob, uint64) {
	var memMapHob EfiMemoryMapHob
	var length uint64

	for _, entry := range memMap {

		memType := strings.TrimSpace(string(entry.Type))

		if memType == "System RAM" {
			memMapHob = append(memMapHob, EfiHobResourceDescriptor{
				Header: EfiHobGenericHeader{
					HobType:   EfiHobTypeResourceDescriptor,
					HobLength: uint16(unsafe.Sizeof(EfiHobResourceDescriptor{})),
				},
				ResourceType: EfiResourceSystemMemory,
				ResourceAttribute: EfiResourceAttributePresent |
					EfiResourceAttributeInitialized |
					EfiResourceAttributeTested |
					EfiResourceAttributeUncacheable |
					EfiResourceAttributeWriteCombineable |
					EfiResourceAttributeWriteThroughCacheable |
					EfiResourceAttributeWriteBackCacheable,
				PhysicalStart:  EfiPhysicalAddress(entry.Start),
				ResourceLength: uint64(entry.Size),
			})
		}

		length += uint64(unsafe.Sizeof(EfiHobResourceDescriptor{}))
	}

	return memMapHob, length
}

func hobCreateEndHob() EfiHobGenericHeader {
	return EfiHobGenericHeader{
		HobType:   EfiHobTypeEndOfHobList,
		HobLength: uint16(unsafe.Sizeof(EfiHobGenericHeader{})),
		Reserved:  0,
	}
}

// Create handoff information Hob with specified length.
// Handoff info Hob should be created after all Hobs in Hob list
// have been created, since length of Hob list should be provided
// as input parameter.
func hobCreateEfiHobHandoffInfoTable(length uint64) EfiHobHandoffInfoTable {
	length += uint64(unsafe.Sizeof(EfiHobHandoffInfoTable{}))

	return EfiHobHandoffInfoTable{
		Header: EfiHobGenericHeader{
			HobType:   EfiHobTypeHandoff,
			HobLength: uint16(unsafe.Sizeof(EfiHobHandoffInfoTable{})),
		},
		Version:             EfiHobHandoffInfoVersion,
		BootMode:            EfiHobHandoffInfoBootMode,
		EfiMemoryTop:        EfiHobHandoffInfoEfiMemoryTop,
		EfiMemoryBottom:     EfiHobHandoffInfoEfiMemoryBottom,
		EfiFreeMemoryTop:    EfiHobHandoffInfoFreeEfiMemoryTop,
		EfiFreeMemoryBottom: EfiPhysicalAddress(EfiHobHandoffInfoFreeMemoryBottom + length + 8),
		EfiEndOfHobList:     EfiPhysicalAddress(EfiHobHandoffInfoFreeMemoryBottom + length),
	}
}
