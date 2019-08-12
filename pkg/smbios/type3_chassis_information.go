// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

// Much of this is auto-generated. If adding a new type, see README for instructions.

// ChassisInformation is defined in DSP0134 7.4.
type ChassisInformation struct {
	Table
	Manufacturer                  string                    // 04h
	Type                          ChassisType               // 05h
	Version                       string                    // 06h
	SerialNumber                  string                    // 07h
	AssetTagNumber                string                    // 08h
	BootupState                   ChassisState              // 09h
	PowerSupplyState              ChassisState              // 0Ah
	ThermalState                  ChassisState              // 0Bh
	SecurityStatus                ChassisSecurityStatus     // 0Ch
	OEMInformation                uint32                    // 0Dh
	Height                        uint8                     // 11h
	NumberOfPowerCords            uint8                     // 12h
	ContainedElementCount         uint8                     // 13h
	ContainedElementsRecordLength uint8                     // 14h
	ContainedElements             []ChassisContainedElement `smbios:"-"` // 15h
	SKUNumber                     string                    `smbios:"-"` // 15h + CEC * CERL
}

// ChassisContainedElement is defined in DSP0134 7.4.4.
type ChassisContainedElement struct {
	Type ChassisElementType // 00h
	Min  uint8              // 01h
	Max  uint8              // 02h
}

// NewChassisInformation parses a generic Table into ChassisInformation.
func NewChassisInformation(t *Table) (*ChassisInformation, error) {
	if t.Type != TableTypeChassisInformation {
		return nil, fmt.Errorf("invalid table type %d", t.Type)
	}
	if t.Len() < 0x9 {
		return nil, errors.New("required fields missing")
	}
	si := &ChassisInformation{Table: *t}
	off, err := parseStruct(t, 0 /* off */, false /* complete */, si)
	if err != nil {
		return nil, err
	}
	if si.ContainedElementCount > 0 && si.ContainedElementsRecordLength > 0 {
		if t.Len() < off+int(si.ContainedElementCount)*int(si.ContainedElementsRecordLength) {
			return nil, fmt.Errorf("invalid data length %d %d %d %d", t.Len(), off, si.ContainedElementCount, si.ContainedElementsRecordLength)
		}
		for i := 0; i < int(si.ContainedElementCount); i++ {
			var e ChassisContainedElement
			eb, _ := t.GetBytesAt(off, int(si.ContainedElementsRecordLength))
			if err = binary.Read(bytes.NewReader(eb), binary.LittleEndian, &e); err != nil {
				return nil, err
			}
			si.ContainedElements = append(si.ContainedElements, e)
			off += int(si.ContainedElementsRecordLength)
		}
	}
	if off < t.Len() {
		si.SKUNumber, _ = t.GetStringAt(off)
	}
	return si, nil
}

func (si *ChassisInformation) String() string {
	lockStr := "Not Present"
	if si.Type&0x80 != 0 {
		lockStr = "Present"
	}
	lines := []string{
		si.Header.String(),
		fmt.Sprintf("Manufacturer: %s", si.Manufacturer),
		fmt.Sprintf("Type: %s", si.Type),
		fmt.Sprintf("Lock: %s", lockStr),
		fmt.Sprintf("Version: %s", si.Version),
		fmt.Sprintf("Serial Number: %s", si.SerialNumber),
		fmt.Sprintf("Asset Tag: %s", si.AssetTagNumber),
	}
	if si.Len() >= 9 { // 2.1+
		lines = append(lines,
			fmt.Sprintf("Boot-up State: %s", si.BootupState),
			fmt.Sprintf("Power Supply State: %s", si.PowerSupplyState),
			fmt.Sprintf("Thermal State: %s", si.ThermalState),
			fmt.Sprintf("Security Status: %s", si.SecurityStatus),
		)
	}
	if si.Len() >= 0xd { // 2.3+
		heightStr, numPCStr := "Unspecified", "Unspecified"
		if si.Height != 0 {
			heightStr = fmt.Sprintf("%d U", si.Height)
		}
		if si.NumberOfPowerCords != 0 {
			numPCStr = fmt.Sprintf("%d", si.NumberOfPowerCords)
		}
		lines = append(lines,
			fmt.Sprintf("OEM Information: 0x%08X", si.OEMInformation),
			fmt.Sprintf("Height: %s", heightStr),
			fmt.Sprintf("Number Of Power Cords: %s", numPCStr),
		)
		lines = append(lines,
			fmt.Sprintf("Contained Elements: %d", si.ContainedElementCount),
		)
		for _, e := range si.ContainedElements {
			lines = append(lines,
				fmt.Sprintf("\t%s %d-%d", e.Type, e.Min, e.Max),
			)
		}
	}
	if si.Len() > 0x15+int(si.ContainedElementCount)*int(si.ContainedElementsRecordLength) {
		lines = append(lines,
			fmt.Sprintf("SKU Number: %s", si.SKUNumber),
		)
	}
	return strings.Join(lines, "\n\t")
}

// ChassisType is defined in DSP0134 7.4.1.
type ChassisType uint8

// ChassisType values are defined in DSP0134 7.4.1.
const (
	ChassisTypeOther               ChassisType = 0x01 // Other
	ChassisTypeUnknown                         = 0x02 // Unknown
	ChassisTypeDesktop                         = 0x03 // Desktop
	ChassisTypeLowProfileDesktop               = 0x04 // Low Profile Desktop
	ChassisTypePizzaBox                        = 0x05 // Pizza Box
	ChassisTypeMiniTower                       = 0x06 // Mini Tower
	ChassisTypeTower                           = 0x07 // Tower
	ChassisTypePortable                        = 0x08 // Portable
	ChassisTypeLaptop                          = 0x09 // Laptop
	ChassisTypeNotebook                        = 0x0a // Notebook
	ChassisTypeHandHeld                        = 0x0b // Hand Held
	ChassisTypeDockingStation                  = 0x0c // Docking Station
	ChassisTypeAllInOne                        = 0x0d // All in One
	ChassisTypeSubNotebook                     = 0x0e // Sub Notebook
	ChassisTypeSpacesaving                     = 0x0f // Space-saving
	ChassisTypeLunchBox                        = 0x10 // Lunch Box
	ChassisTypeMainServerChassis               = 0x11 // Main Server Chassis
	ChassisTypeExpansionChassis                = 0x12 // Expansion Chassis
	ChassisTypeSubChassis                      = 0x13 // SubChassis
	ChassisTypeBusExpansionChassis             = 0x14 // Bus Expansion Chassis
	ChassisTypePeripheralChassis               = 0x15 // Peripheral Chassis
	ChassisTypeRAIDChassis                     = 0x16 // RAID Chassis
	ChassisTypeRackMountChassis                = 0x17 // Rack Mount Chassis
	ChassisTypeSealedcasePC                    = 0x18 // Sealed-case PC
	ChassisTypeMultisystemChassis              = 0x19 // Multi-system chassis
	ChassisTypeCompactPCI                      = 0x1a // Compact PCI
	ChassisTypeAdvancedTCA                     = 0x1b // Advanced TCA
	ChassisTypeBlade                           = 0x1c // Blade
	ChassisTypeBladeChassis                    = 0x1d // Blade Chassis
	ChassisTypeTablet                          = 0x1e // Tablet
	ChassisTypeConvertible                     = 0x1f // Convertible
	ChassisTypeDetachable                      = 0x20 // Detachable
	ChassisTypeIoTGateway                      = 0x21 // IoT Gateway
	ChassisTypeEmbeddedPC                      = 0x22 // Embedded PC
	ChassisTypeMiniPC                          = 0x23 // Mini PC
	ChassisTypeStickPC                         = 0x24 // Stick PC
)

func (v ChassisType) String() string {
	switch v & 0x7f {
	case ChassisTypeOther:
		return "Other"
	case ChassisTypeUnknown:
		return "Unknown"
	case ChassisTypeDesktop:
		return "Desktop"
	case ChassisTypeLowProfileDesktop:
		return "Low Profile Desktop"
	case ChassisTypePizzaBox:
		return "Pizza Box"
	case ChassisTypeMiniTower:
		return "Mini Tower"
	case ChassisTypeTower:
		return "Tower"
	case ChassisTypePortable:
		return "Portable"
	case ChassisTypeLaptop:
		return "Laptop"
	case ChassisTypeNotebook:
		return "Notebook"
	case ChassisTypeHandHeld:
		return "Hand Held"
	case ChassisTypeDockingStation:
		return "Docking Station"
	case ChassisTypeAllInOne:
		return "All In One"
	case ChassisTypeSubNotebook:
		return "Sub Notebook"
	case ChassisTypeSpacesaving:
		return "Space-saving"
	case ChassisTypeLunchBox:
		return "Lunch Box"
	case ChassisTypeMainServerChassis:
		return "Main Server Chassis"
	case ChassisTypeExpansionChassis:
		return "Expansion Chassis"
	case ChassisTypeSubChassis:
		return "Sub Chassis"
	case ChassisTypeBusExpansionChassis:
		return "Bus Expansion Chassis"
	case ChassisTypePeripheralChassis:
		return "Peripheral Chassis"
	case ChassisTypeRAIDChassis:
		return "RAID Chassis"
	case ChassisTypeRackMountChassis:
		return "Rack Mount Chassis"
	case ChassisTypeSealedcasePC:
		return "Sealed-case PC"
	case ChassisTypeMultisystemChassis:
		return "Multi-system"
	case ChassisTypeCompactPCI:
		return "CompactPCI"
	case ChassisTypeAdvancedTCA:
		return "AdvancedTCA"
	case ChassisTypeBlade:
		return "Blade"
	case ChassisTypeBladeChassis:
		return "Blade Chassis"
	case ChassisTypeTablet:
		return "Tablet"
	case ChassisTypeConvertible:
		return "Convertible"
	case ChassisTypeDetachable:
		return "Detachable"
	case ChassisTypeIoTGateway:
		return "IoT Gateway"
	case ChassisTypeEmbeddedPC:
		return "Embedded PC"
	case ChassisTypeMiniPC:
		return "Mini PC"
	case ChassisTypeStickPC:
		return "Stick PC"
	}
	return fmt.Sprintf("%d", v)
}

// ChassisState is defined in DSP0134 7.4.2.
type ChassisState uint8

// ChassisState values are defined in DSP0134 7.4.2.
const (
	ChassisStateOther          ChassisState = 0x01 // Other
	ChassisStateUnknown                     = 0x02 // Unknown
	ChassisStateSafe                        = 0x03 // Safe
	ChassisStateWarning                     = 0x04 // Warning
	ChassisStateCritical                    = 0x05 // Critical
	ChassisStateNonrecoverable              = 0x06 // Non-recoverable
)

func (v ChassisState) String() string {
	switch v {
	case ChassisStateOther:
		return "Other"
	case ChassisStateUnknown:
		return "Unknown"
	case ChassisStateSafe:
		return "Safe"
	case ChassisStateWarning:
		return "Warning"
	case ChassisStateCritical:
		return "Critical"
	case ChassisStateNonrecoverable:
		return "Non-recoverable"
	}
	return fmt.Sprintf("%d", v)
}

// ChassisSecurityStatus is defined in DSP0134 7.4.3.
type ChassisSecurityStatus uint8

// ChassisSecurityStatus values are defined in DSP0134 7.4.3.
const (
	ChassisSecurityStatusOther                      ChassisSecurityStatus = 0x01 // Other
	ChassisSecurityStatusUnknown                                          = 0x02 // Unknown
	ChassisSecurityStatusNone                                             = 0x03 // None
	ChassisSecurityStatusExternalInterfaceLockedOut                       = 0x04 // External interface locked out
	ChassisSecurityStatusExternalInterfaceEnabled                         = 0x05 // External interface enabled
)

func (v ChassisSecurityStatus) String() string {
	switch v {
	case ChassisSecurityStatusOther:
		return "Other"
	case ChassisSecurityStatusUnknown:
		return "Unknown"
	case ChassisSecurityStatusNone:
		return "None"
	case ChassisSecurityStatusExternalInterfaceLockedOut:
		return "External Interface Locked Out"
	case ChassisSecurityStatusExternalInterfaceEnabled:
		return "External Interface Enabled"
	}
	return fmt.Sprintf("%d", v)
}

// ChassisElementType is defined in DSP0134 7.4.4.
type ChassisElementType uint8

func (v ChassisElementType) String() string {
	if v&0x80 != 0 {
		return fmt.Sprintf("%s", TableType(v&0x7f))
	}
	return fmt.Sprintf("%s", BoardType(v&0x7f))
}
