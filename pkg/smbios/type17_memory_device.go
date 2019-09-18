// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"errors"
	"fmt"
	"strings"
)

// Much of this is auto-generated. If adding a new type, see README for instructions.

// MemoryDevice is defined in DSP0134 7.18.
type MemoryDevice struct {
	Table
	PhysicalMemoryArrayHandle         uint16                              // 04h
	MemoryErrorInformationHandle      uint16                              // 06h
	TotalWidth                        uint16                              // 08h
	DataWidth                         uint16                              // 0Ah
	Size                              uint16                              // 0Ch
	FormFactor                        MemoryDeviceFormFactor              // 0Eh
	DeviceSet                         uint8                               // 0Fh
	DeviceLocator                     string                              // 10h
	BankLocator                       string                              // 11h
	Type                              MemoryDeviceType                    // 12h
	TypeDetail                        MemoryDeviceTypeDetail              // 13h
	Speed                             uint16                              // 15h
	Manufacturer                      string                              // 17h
	SerialNumber                      string                              // 18h
	AssetTag                          string                              // 19h
	PartNumber                        string                              // 1Ah
	Attributes                        uint8                               // 1Bh
	ExtendedSize                      uint32                              // 1Ch
	ConfiguredSpeed                   uint16                              // 20h
	MinimumVoltage                    uint16                              // 22h
	MaximumVoltage                    uint16                              // 24h
	ConfiguredVoltage                 uint16                              // 26h
	Technology                        MemoryDeviceTechnology              // 28h
	OperatingModeCapability           MemoryDeviceOperatingModeCapability // 29h
	FirmwareVersion                   string                              // 2Bh
	ModuleManufacturerID              uint16                              // 2Ch
	ModuleProductID                   uint16                              // 2Eh
	SubsystemControllerManufacturerID uint16                              // 30h
	SubsystemControllerProductID      uint16                              // 32h
	NonvolatileSize                   uint64                              // 34h
	VolatileSize                      uint64                              // 3Ch
	CacheSize                         uint64                              // 44h
	LogicalSize                       uint64                              // 4Ch
}

// NewMemoryDevice parses a generic Table into MemoryDevice.
func NewMemoryDevice(t *Table) (*MemoryDevice, error) {
	if t.Type != TableTypeMemoryDevice {
		return nil, fmt.Errorf("invalid table type %d", t.Type)
	}
	if t.Len() < 0x15 {
		return nil, errors.New("required fields missing")
	}
	md := &MemoryDevice{Table: *t}
	_, err := parseStruct(t, 0 /* off */, false /* complete */, md)
	if err != nil {
		return nil, err
	}
	return md, nil
}

// GetSizeBytes returns size of the memory device, in bytes.
func (md *MemoryDevice) GetSizeBytes() uint64 {
	switch md.Size {
	case 0:
		return 0
	case 0x7fff:
		return uint64(md.ExtendedSize&0x7fffffff) * 1024 * 1024
	default:
		mul := uint64(1024 * 1024)
		if md.Size&0x8000 != 0 {
			mul = 1024
		}
		return uint64(md.Size&0x7fff) * mul
	}
}

func (md *MemoryDevice) String() string {
	ehStr := ""
	switch md.MemoryErrorInformationHandle {
	case 0xffff:
		ehStr = "No Error"
	case 0xfffe:
		ehStr = "Not Provided"
	default:
		ehStr = fmt.Sprintf("0x%04X", md.MemoryErrorInformationHandle)
	}

	bitWidthStr := func(v uint16) string {
		if v == 0 || v == 0xffff {
			return "Unknown"
		}
		return fmt.Sprintf("%d bits", v)
	}

	setStr := ""
	switch md.DeviceSet {
	case 0:
		setStr = "None"
	case 0xff:
		setStr = "Unknown"
	default:
		setStr = fmt.Sprintf("%d", md.DeviceSet)
	}

	speedStr := func(v uint16) string {
		if v == 0 || v == 0xffff {
			return "Unknown"
		}
		return fmt.Sprintf("%d MT/s", v)
	}

	moduleSizeStr := "No Module Installed"
	if md.GetSizeBytes() > 0 {
		moduleSizeStr = kmgt(md.GetSizeBytes())
	}

	lines := []string{
		md.Header.String(),
		fmt.Sprintf("Array Handle: 0x%04X", md.PhysicalMemoryArrayHandle),
		fmt.Sprintf("Error Information Handle: %s", ehStr),
		fmt.Sprintf("Total Width: %s", bitWidthStr(md.TotalWidth)),
		fmt.Sprintf("Data Width: %s", bitWidthStr(md.DataWidth)),
		fmt.Sprintf("Size: %s", moduleSizeStr),
		fmt.Sprintf("Form Factor: %s", md.FormFactor),
		fmt.Sprintf("Set: %s", setStr),
		fmt.Sprintf("Locator: %s", md.DeviceLocator),
		fmt.Sprintf("Bank Locator: %s", md.BankLocator),
		fmt.Sprintf("Type: %s", md.Type),
		fmt.Sprintf("Type Detail: %s", md.TypeDetail),
	}
	if md.Len() > 0x15 {
		lines = append(lines,
			fmt.Sprintf("Speed: %s", speedStr(md.Speed)),
			fmt.Sprintf("Manufacturer: %s", md.Manufacturer),
			fmt.Sprintf("Serial Number: %s", md.SerialNumber),
			fmt.Sprintf("Asset Tag: %s", md.AssetTag),
			fmt.Sprintf("Part Number: %s", md.PartNumber),
		)
	}
	if md.Len() > 0x1b {
		rankStr := "Unknown"
		if md.Attributes&0xf != 0 {
			rankStr = fmt.Sprintf("%d", md.Attributes&0xf)
		}
		lines = append(lines, fmt.Sprintf("Rank: %s", rankStr))
	}
	if md.Len() > 0x1c {
		lines = append(lines, fmt.Sprintf("Configured Memory Speed: %s", speedStr(md.ConfiguredSpeed)))
	}
	if md.Len() > 0x22 {
		voltageStr := func(v uint16) string {
			switch {
			case v == 0:
				return "Unknown"
			case v%100 == 0:
				return fmt.Sprintf("%.1f V", float32(v)/1000.0)
			default:
				return fmt.Sprintf("%g V", float32(v)/1000.0)
			}
		}
		lines = append(lines,
			fmt.Sprintf("Minimum Voltage: %s", voltageStr(md.MinimumVoltage)),
			fmt.Sprintf("Maximum Voltage: %s", voltageStr(md.MaximumVoltage)),
			fmt.Sprintf("Configured Voltage: %s", voltageStr(md.ConfiguredVoltage)),
		)
	}
	if md.Len() > 0x28 {
		manufacturerIDStr := func(id uint16) string {
			if id == 0 {
				return "Unknown"
			}
			return fmt.Sprintf("Bank %d, Hex 0x%02X", (id&0x7F)+1, id>>8)
		}
		productIDStr := func(v uint16) string {
			if v == 0 {
				return "Unknown"
			}
			return fmt.Sprintf("0x%04X", v)
		}
		sizeStr := func(v uint64) string {
			switch v {
			case 0:
				return "None"
			case 0xffffffffffffffff:
				return "Unknown"
			default:
				return kmgt(v)
			}
		}
		lines = append(lines,
			fmt.Sprintf("Memory Technology: %s", md.Technology),
			fmt.Sprintf("Memory Operating Mode Capability: %s", md.OperatingModeCapability),
			fmt.Sprintf("Firmware Version: %s", md.FirmwareVersion),
			fmt.Sprintf("Module Manufacturer ID: %s", manufacturerIDStr(md.ModuleManufacturerID)),
			fmt.Sprintf("Module Product ID: %s", productIDStr(md.ModuleProductID)),
			fmt.Sprintf("Memory Subsystem Controller Manufacturer ID: %s", manufacturerIDStr(md.SubsystemControllerManufacturerID)),
			fmt.Sprintf("Memory Subsystem Controller Product ID: %s", productIDStr(md.SubsystemControllerProductID)),
			fmt.Sprintf("Non-Volatile Size: %s", sizeStr(md.NonvolatileSize)),
			fmt.Sprintf("Volatile Size: %s", sizeStr(md.VolatileSize)),
			fmt.Sprintf("Cache Size: %s", sizeStr(md.CacheSize)),
			fmt.Sprintf("Logical Size: %s", sizeStr(md.LogicalSize)),
		)
	}
	return strings.Join(lines, "\n\t")
}

// MemoryDeviceFormFactor is defined in DSP0134 7.18.1.
type MemoryDeviceFormFactor uint8

// MemoryDeviceFormFactor values are defined in DSP0134 7.18.1.
const (
	MemoryDeviceFormFactorOther           MemoryDeviceFormFactor = 0x01 // Other
	MemoryDeviceFormFactorUnknown                                = 0x02 // Unknown
	MemoryDeviceFormFactorSIMM                                   = 0x03 // SIMM
	MemoryDeviceFormFactorSIP                                    = 0x04 // SIP
	MemoryDeviceFormFactorChip                                   = 0x05 // Chip
	MemoryDeviceFormFactorDIP                                    = 0x06 // DIP
	MemoryDeviceFormFactorZIP                                    = 0x07 // ZIP
	MemoryDeviceFormFactorProprietaryCard                        = 0x08 // Proprietary Card
	MemoryDeviceFormFactorDIMM                                   = 0x09 // DIMM
	MemoryDeviceFormFactorTSOP                                   = 0x0a // TSOP
	MemoryDeviceFormFactorRowOfChips                             = 0x0b // Row of chips
	MemoryDeviceFormFactorRIMM                                   = 0x0c // RIMM
	MemoryDeviceFormFactorSODIMM                                 = 0x0d // SODIMM
	MemoryDeviceFormFactorSRIMM                                  = 0x0e // SRIMM
	MemoryDeviceFormFactorFBDIMM                                 = 0x0f // FB-DIMM
)

func (v MemoryDeviceFormFactor) String() string {
	switch v {
	case MemoryDeviceFormFactorOther:
		return "Other"
	case MemoryDeviceFormFactorUnknown:
		return "Unknown"
	case MemoryDeviceFormFactorSIMM:
		return "SIMM"
	case MemoryDeviceFormFactorSIP:
		return "SIP"
	case MemoryDeviceFormFactorChip:
		return "Chip"
	case MemoryDeviceFormFactorDIP:
		return "DIP"
	case MemoryDeviceFormFactorZIP:
		return "ZIP"
	case MemoryDeviceFormFactorProprietaryCard:
		return "Proprietary Card"
	case MemoryDeviceFormFactorDIMM:
		return "DIMM"
	case MemoryDeviceFormFactorTSOP:
		return "TSOP"
	case MemoryDeviceFormFactorRowOfChips:
		return "Row of chips"
	case MemoryDeviceFormFactorRIMM:
		return "RIMM"
	case MemoryDeviceFormFactorSODIMM:
		return "SODIMM"
	case MemoryDeviceFormFactorSRIMM:
		return "SRIMM"
	case MemoryDeviceFormFactorFBDIMM:
		return "FB-DIMM"
	}
	return fmt.Sprintf("%#x", uint8(v))
}

// MemoryDeviceType is defined in DSP0134 7.18.2.
type MemoryDeviceType uint8

// MemoryDeviceType values are defined in DSP0134 7.18.2.
const (
	MemoryDeviceTypeOther                    MemoryDeviceType = 0x01 // Other
	MemoryDeviceTypeUnknown                                   = 0x02 // Unknown
	MemoryDeviceTypeDRAM                                      = 0x03 // DRAM
	MemoryDeviceTypeEDRAM                                     = 0x04 // EDRAM
	MemoryDeviceTypeVRAM                                      = 0x05 // VRAM
	MemoryDeviceTypeSRAM                                      = 0x06 // SRAM
	MemoryDeviceTypeRAM                                       = 0x07 // RAM
	MemoryDeviceTypeROM                                       = 0x08 // ROM
	MemoryDeviceTypeFlash                                     = 0x09 // Flash
	MemoryDeviceTypeEEPROM                                    = 0x0a // EEPROM
	MemoryDeviceTypeFEPROM                                    = 0x0b // FEPROM
	MemoryDeviceTypeEPROM                                     = 0x0c // EPROM
	MemoryDeviceTypeCDRAM                                     = 0x0d // CDRAM
	MemoryDeviceType3DRAM                                     = 0x0e // 3DRAM
	MemoryDeviceTypeSDRAM                                     = 0x0f // SDRAM
	MemoryDeviceTypeSGRAM                                     = 0x10 // SGRAM
	MemoryDeviceTypeRDRAM                                     = 0x11 // RDRAM
	MemoryDeviceTypeDDR                                       = 0x12 // DDR
	MemoryDeviceTypeDDR2                                      = 0x13 // DDR2
	MemoryDeviceTypeDDR2FBDIMM                                = 0x14 // DDR2 FB-DIMM
	MemoryDeviceTypeDDR3                                      = 0x18 // DDR3
	MemoryDeviceTypeFBD2                                      = 0x19 // FBD2
	MemoryDeviceTypeDDR4                                      = 0x1a // DDR4
	MemoryDeviceTypeLPDDR                                     = 0x1b // LPDDR
	MemoryDeviceTypeLPDDR2                                    = 0x1c // LPDDR2
	MemoryDeviceTypeLPDDR3                                    = 0x1d // LPDDR3
	MemoryDeviceTypeLPDDR4                                    = 0x1e // LPDDR4
	MemoryDeviceTypeLogicalNonvolatileDevice                  = 0x1f // Logical non-volatile device
)

func (v MemoryDeviceType) String() string {
	switch v {
	case MemoryDeviceTypeOther:
		return "Other"
	case MemoryDeviceTypeUnknown:
		return "Unknown"
	case MemoryDeviceTypeDRAM:
		return "DRAM"
	case MemoryDeviceTypeEDRAM:
		return "EDRAM"
	case MemoryDeviceTypeVRAM:
		return "VRAM"
	case MemoryDeviceTypeSRAM:
		return "SRAM"
	case MemoryDeviceTypeRAM:
		return "RAM"
	case MemoryDeviceTypeROM:
		return "ROM"
	case MemoryDeviceTypeFlash:
		return "Flash"
	case MemoryDeviceTypeEEPROM:
		return "EEPROM"
	case MemoryDeviceTypeFEPROM:
		return "FEPROM"
	case MemoryDeviceTypeEPROM:
		return "EPROM"
	case MemoryDeviceTypeCDRAM:
		return "CDRAM"
	case MemoryDeviceType3DRAM:
		return "3DRAM"
	case MemoryDeviceTypeSDRAM:
		return "SDRAM"
	case MemoryDeviceTypeSGRAM:
		return "SGRAM"
	case MemoryDeviceTypeRDRAM:
		return "RDRAM"
	case MemoryDeviceTypeDDR:
		return "DDR"
	case MemoryDeviceTypeDDR2:
		return "DDR2"
	case MemoryDeviceTypeDDR2FBDIMM:
		return "DDR2 FB-DIMM"
	case MemoryDeviceTypeDDR3:
		return "DDR3"
	case MemoryDeviceTypeFBD2:
		return "FBD2"
	case MemoryDeviceTypeDDR4:
		return "DDR4"
	case MemoryDeviceTypeLPDDR:
		return "LPDDR"
	case MemoryDeviceTypeLPDDR2:
		return "LPDDR2"
	case MemoryDeviceTypeLPDDR3:
		return "LPDDR3"
	case MemoryDeviceTypeLPDDR4:
		return "LPDDR4"
	case MemoryDeviceTypeLogicalNonvolatileDevice:
		return "Logical non-volatile device"
	}
	return fmt.Sprintf("%#x", uint8(v))
}

// MemoryDeviceTypeDetail is defined in DSP0134 7.18.3.
type MemoryDeviceTypeDetail uint16

// MemoryDeviceTypeDetail fields are defined in DSP0134 7.18.3
const (
	MemoryDeviceTypeDetailOther                  = (1 << 1)  // Other
	MemoryDeviceTypeDetailUnknown                = (1 << 2)  // Unknown
	MemoryDeviceTypeDetailFastpaged              = (1 << 3)  // Fast-paged
	MemoryDeviceTypeDetailStaticColumn           = (1 << 4)  // Static column
	MemoryDeviceTypeDetailPseudostatic           = (1 << 5)  // Pseudo-static
	MemoryDeviceTypeDetailRAMBUS                 = (1 << 6)  // RAMBUS
	MemoryDeviceTypeDetailSynchronous            = (1 << 7)  // Synchronous
	MemoryDeviceTypeDetailCMOS                   = (1 << 8)  // CMOS
	MemoryDeviceTypeDetailEDO                    = (1 << 9)  // EDO
	MemoryDeviceTypeDetailWindowDRAM             = (1 << 10) // Window DRAM
	MemoryDeviceTypeDetailCacheDRAM              = (1 << 11) // Cache DRAM
	MemoryDeviceTypeDetailNonvolatile            = (1 << 12) // Non-volatile
	MemoryDeviceTypeDetailRegisteredBuffered     = (1 << 13) // Registered (Buffered)
	MemoryDeviceTypeDetailUnbufferedUnregistered = (1 << 14) // Unbuffered (Unregistered)
	MemoryDeviceTypeDetailLRDIMM                 = (1 << 15) // LRDIMM
)

func (v MemoryDeviceTypeDetail) String() string {
	if v&0xfffe == 0 {
		return "None"
	}
	var lines []string
	if v&MemoryDeviceTypeDetailOther != 0 {
		lines = append(lines, "Other")
	}
	if v&MemoryDeviceTypeDetailUnknown != 0 {
		lines = append(lines, "Unknown")
	}
	if v&MemoryDeviceTypeDetailFastpaged != 0 {
		lines = append(lines, "Fast-paged")
	}
	if v&MemoryDeviceTypeDetailStaticColumn != 0 {
		lines = append(lines, "Static column")
	}
	if v&MemoryDeviceTypeDetailPseudostatic != 0 {
		lines = append(lines, "Pseudo-static")
	}
	if v&MemoryDeviceTypeDetailRAMBUS != 0 {
		lines = append(lines, "RAMBUS")
	}
	if v&MemoryDeviceTypeDetailSynchronous != 0 {
		lines = append(lines, "Synchronous")
	}
	if v&MemoryDeviceTypeDetailCMOS != 0 {
		lines = append(lines, "CMOS")
	}
	if v&MemoryDeviceTypeDetailEDO != 0 {
		lines = append(lines, "EDO")
	}
	if v&MemoryDeviceTypeDetailWindowDRAM != 0 {
		lines = append(lines, "Window DRAM")
	}
	if v&MemoryDeviceTypeDetailCacheDRAM != 0 {
		lines = append(lines, "Cache DRAM")
	}
	if v&MemoryDeviceTypeDetailNonvolatile != 0 {
		lines = append(lines, "Non-volatile")
	}
	if v&MemoryDeviceTypeDetailRegisteredBuffered != 0 {
		lines = append(lines, "Registered (Buffered)")
	}
	if v&MemoryDeviceTypeDetailUnbufferedUnregistered != 0 {
		lines = append(lines, "Unbuffered (Unregistered)")
	}
	if v&MemoryDeviceTypeDetailLRDIMM != 0 {
		lines = append(lines, "LRDIMM")
	}
	return strings.Join(lines, " ")
}

// MemoryDeviceTechnology is defined in DSP0134 7.18.6.
type MemoryDeviceTechnology uint8

// MemoryDeviceTechnology values are defined in DSP0134 7.18.6.
const (
	MemoryDeviceTechnologyOther                 MemoryDeviceTechnology = 0x01 // Other
	MemoryDeviceTechnologyUnknown                                      = 0x02 // Unknown
	MemoryDeviceTechnologyDRAM                                         = 0x03 // DRAM
	MemoryDeviceTechnologyNVDIMMN                                      = 0x04 // NVDIMM-N
	MemoryDeviceTechnologyNVDIMMF                                      = 0x05 // NVDIMM-F
	MemoryDeviceTechnologyNVDIMMP                                      = 0x06 // NVDIMM-P
	MemoryDeviceTechnologyIntelPersistentMemory                        = 0x07 // Intel persistent memory
)

func (v MemoryDeviceTechnology) String() string {
	switch v {
	case MemoryDeviceTechnologyOther:
		return "Other"
	case MemoryDeviceTechnologyUnknown:
		return "Unknown"
	case MemoryDeviceTechnologyDRAM:
		return "DRAM"
	case MemoryDeviceTechnologyNVDIMMN:
		return "NVDIMM-N"
	case MemoryDeviceTechnologyNVDIMMF:
		return "NVDIMM-F"
	case MemoryDeviceTechnologyNVDIMMP:
		return "NVDIMM-P"
	case MemoryDeviceTechnologyIntelPersistentMemory:
		return "Intel persistent memory"
	}
	return fmt.Sprintf("%#x", uint8(v))
}

// MemoryDeviceOperatingModeCapability is defined in DSP0134 7.18.7.
type MemoryDeviceOperatingModeCapability uint16

// MemoryDeviceOperatingModeCapability fields are defined in DSP0134 x.x.x
const (
	MemoryDeviceOperatingModeCapabilityOther                           = (1 << 1) // Other
	MemoryDeviceOperatingModeCapabilityUnknown                         = (1 << 2) // Unknown
	MemoryDeviceOperatingModeCapabilityVolatileMemory                  = (1 << 3) // Volatile memory
	MemoryDeviceOperatingModeCapabilityByteaccessiblePersistentMemory  = (1 << 4) // Byte-accessible persistent memory
	MemoryDeviceOperatingModeCapabilityBlockaccessiblePersistentMemory = (1 << 5) // Block-accessible persistent memory
)

func (v MemoryDeviceOperatingModeCapability) String() string {
	if v&0xfffe == 0 {
		return "None"
	}
	var lines []string
	if v&MemoryDeviceOperatingModeCapabilityOther != 0 {
		lines = append(lines, "Other")
	}
	if v&MemoryDeviceOperatingModeCapabilityUnknown != 0 {
		lines = append(lines, "Unknown")
	}
	if v&MemoryDeviceOperatingModeCapabilityVolatileMemory != 0 {
		lines = append(lines, "Volatile memory")
	}
	if v&MemoryDeviceOperatingModeCapabilityByteaccessiblePersistentMemory != 0 {
		lines = append(lines, "Byte-accessible persistent memory")
	}
	if v&MemoryDeviceOperatingModeCapabilityBlockaccessiblePersistentMemory != 0 {
		lines = append(lines, "Block-accessible persistent memory")
	}
	return strings.Join(lines, " ")
}
