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

// IPMIDeviceInfo is defined in DSP0134 7.39.
type IPMIDeviceInfo struct {
	Table
	InterfaceType                    BMCInterfaceType // 04h
	IPMISpecificationRevision        uint8            // 05h
	I2CSlaveAddress                  uint8            // 06h
	NVStorageDeviceAddress           uint8            // 07h
	BaseAddress                      uint64           // 08h
	BaseAddressModifierInterruptInfo uint8            // 10h
	InterruptNumber                  uint8            // 11h
}

// ParseIPMIDeviceInfo parses a generic Table into IPMIDeviceInfo.
func ParseIPMIDeviceInfo(t *Table) (*IPMIDeviceInfo, error) {
	if t.Type != TableTypeIPMIDeviceInfo {
		return nil, fmt.Errorf("invalid table type %d", t.Type)
	}
	if t.Len() < 0x12 {
		return nil, errors.New("required fields missing")
	}
	di := &IPMIDeviceInfo{Table: *t}
	if _, err := parseStruct(t, 0 /* off */, false /* complete */, di); err != nil {
		return nil, err
	}
	return di, nil
}

func (di *IPMIDeviceInfo) String() string {
	nvs := "Not Present"
	if di.NVStorageDeviceAddress != 0xff {
		nvs = fmt.Sprintf("%d", di.NVStorageDeviceAddress)
	}

	baType := "Memory-mapped"
	if di.BaseAddress&1 != 0 {
		baType = "I/O"
	}
	ba := (di.BaseAddress & 0xfffffffffffffffe) | uint64((di.BaseAddressModifierInterruptInfo>>4)&1)

	lines := []string{
		di.Header.String(),
		fmt.Sprintf("Interface Type: %s", di.InterfaceType),
		fmt.Sprintf("Specification Version: %d.%d", di.IPMISpecificationRevision>>4, di.IPMISpecificationRevision&0xf),
		fmt.Sprintf("I2C Slave Address: 0x%02x", di.I2CSlaveAddress>>1),
		fmt.Sprintf("NV Storage Device: %s", nvs),
		fmt.Sprintf("Base Address: 0x%016X (%s)", ba, baType),
	}
	if di.InterfaceType != BMCInterfaceTypeSSIFSMBusSystemInterface {
		rss := ""
		switch (di.BaseAddressModifierInterruptInfo >> 6) & 3 {
		case 0:
			rss = "Successive Byte Boundaries"
		case 1:
			rss = "32-bit Boundaries"
		case 2:
			rss = "16-bit Boundaries"
		case 3:
			rss = outOfSpec
		}
		lines = append(lines, fmt.Sprintf("Register Spacing: %s", rss))
		if di.BaseAddressModifierInterruptInfo&(1<<3) != 0 {
			if di.BaseAddressModifierInterruptInfo&(1<<1) != 0 {
				lines = append(lines, "Interrupt Polarity: Active High")
			} else {
				lines = append(lines, "Interrupt Polarity: Active Low")
			}
			if di.BaseAddressModifierInterruptInfo&(1<<0) != 0 {
				lines = append(lines, "Interrupt Trigger Mode: Level")
			} else {
				lines = append(lines, "Interrupt Trigger Mode: Edge")
			}
		}
	}
	if di.InterruptNumber != 0 {
		lines = append(lines, fmt.Sprintf("Interrupt Number: %d", di.InterruptNumber))
	}
	return strings.Join(lines, "\n\t")
}

// BMCInterfaceType is defined in DSP0134 7.39.1.
type BMCInterfaceType uint8

// BMCInterfaceType values are defined in DSP0134 7.39.1.
const (
	BMCInterfaceTypeUnknown                           BMCInterfaceType = 0x00 // Unknown
	BMCInterfaceTypeKCSKeyboardControllerStyle        BMCInterfaceType = 0x01 // KCS: Keyboard Controller Style
	BMCInterfaceTypeSMICServerManagementInterfaceChip BMCInterfaceType = 0x02 // SMIC: Server Management Interface Chip
	BMCInterfaceTypeBTBlockTransfer                   BMCInterfaceType = 0x03 // BT: Block Transfer
	BMCInterfaceTypeSSIFSMBusSystemInterface          BMCInterfaceType = 0x04 // SSIF: SMBus System Interface
)

func (v BMCInterfaceType) String() string {
	names := map[BMCInterfaceType]string{
		BMCInterfaceTypeUnknown:                           "Unknown",
		BMCInterfaceTypeKCSKeyboardControllerStyle:        "KCS (Keyboard Control Style)",
		BMCInterfaceTypeSMICServerManagementInterfaceChip: "SMIC (Server Management Interface Chip)",
		BMCInterfaceTypeBTBlockTransfer:                   "BT (Block Transfer)",
		BMCInterfaceTypeSSIFSMBusSystemInterface:          "SSIF (SMBus System Interface)",
	}
	if name, ok := names[v]; ok {
		return name
	}
	return fmt.Sprintf("%#x", uint8(v))
}
