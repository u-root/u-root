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

// TPMDevice is defined in DSP0134 7.44.
type TPMDevice struct {
	Table
	VendorID         TPMDeviceVendorID        `smbios:"-,skip=4"` // 04h
	MajorSpecVersion uint8                    // 08h
	MinorSpecVersion uint8                    // 09h
	FirmwareVersion1 uint32                   // 0Ah
	FirmwareVersion2 uint32                   // 0Eh
	Description      string                   // 12h
	Characteristics  TPMDeviceCharacteristics // 13h
	OEMDefined       uint32                   // 1Bh
}

// NewTPMDevice parses a generic Table into TPMDevice.
func NewTPMDevice(t *Table) (*TPMDevice, error) {
	return newTPMDevice(parseStruct, t)
}

func newTPMDevice(parseFn parseStructure, t *Table) (*TPMDevice, error) {
	if t.Type != TableTypeTPMDevice {
		return nil, fmt.Errorf("invalid table type %d", t.Type)
	}
	if t.Len() < 0x1f {
		return nil, errors.New("required fields missing")
	}
	di := &TPMDevice{Table: *t}
	if _, err := parseFn(t, 0 /* off */, false /* complete */, di); err != nil {
		return nil, err
	}
	vid, _ := di.GetBytesAt(4, 4)
	copy(di.VendorID[:], vid)
	return di, nil
}

func (di *TPMDevice) String() string {
	lines := []string{
		di.Header.String(),
		fmt.Sprintf("Vendor ID: %s", di.VendorID),
		fmt.Sprintf("Specification Version: %d.%d", di.MajorSpecVersion, di.MinorSpecVersion),
	}
	switch di.MajorSpecVersion {
	case 1:
		lines = append(lines, fmt.Sprintf("Firmware Revision: %d.%d",
			(di.FirmwareVersion1>>16)&0xff, (di.FirmwareVersion1>>24)&0xff),
		)
	case 2:
		lines = append(lines, fmt.Sprintf("Firmware Revision: %d.%d",
			(di.FirmwareVersion1>>16)&0xffff, di.FirmwareVersion1&0xff),
		)
	}
	lines = append(lines,
		fmt.Sprintf("Description: %s", di.Description),
		fmt.Sprintf("Characteristics:\n%s", di.Characteristics),
		fmt.Sprintf("OEM-specific Info: 0x%08X", di.OEMDefined),
	)
	return strings.Join(lines, "\n\t")
}

// TPMDeviceVendorID is defined in TCG Vendor ID Registry.
type TPMDeviceVendorID [4]byte

func (vid TPMDeviceVendorID) String() string {
	// DSP0134 specifies Vendor ID field as 4 BYTEs, not a DWORD, and gives an example value.
	// But Infineon ignores it and puts their VID in LE byte order so it ends up backwards.
	if vid[0] == 0 && vid[1] == 'X' && vid[2] == 'F' && vid[3] == 'I' {
		return "IFX"
	}
	s := ""
	for i := 0; i < 4 && vid[i] != 0; i++ {
		s += string(vid[i])
	}
	return s
}

// TPMDeviceCharacteristics is defined in DSP0134 7.44.1.
type TPMDeviceCharacteristics uint8

// TPMDeviceCharacteristics fields are defined in DSP0134 x.x.x
const (
	TPMDeviceCharacteristicsNotSupported                                 TPMDeviceCharacteristics = 1 << 2 // TPM Device Characteristics are not supported.
	TPMDeviceCharacteristicsFamilyConfigurableViaFirmwareUpdate          TPMDeviceCharacteristics = 1 << 3 // Family configurable via firmware update.
	TPMDeviceCharacteristicsFamilyConfigurableViaPlatformSoftwareSupport TPMDeviceCharacteristics = 1 << 4 // Family configurable via platform software support.
	TPMDeviceCharacteristicsFamilyConfigurableViaOEMProprietaryMechanism TPMDeviceCharacteristics = 1 << 5 // Family configurable via OEM proprietary mechanism.
)

func (v TPMDeviceCharacteristics) String() string {
	if v&TPMDeviceCharacteristicsNotSupported != 0 {
		return "\t\tTPM Device characteristics not supported"
	}
	var lines []string
	if v&TPMDeviceCharacteristicsFamilyConfigurableViaFirmwareUpdate != 0 {
		lines = append(lines, "Family configurable via firmware update")
	}
	if v&TPMDeviceCharacteristicsFamilyConfigurableViaPlatformSoftwareSupport != 0 {
		lines = append(lines, "Family configurable via platform software support")
	}
	if v&TPMDeviceCharacteristicsFamilyConfigurableViaOEMProprietaryMechanism != 0 {
		lines = append(lines, "Family configurable via OEM proprietary mechanism")
	}
	return "\t\t" + strings.Join(lines, "\n\t\t")
}
