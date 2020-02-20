// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
)

// Info contains the SMBIOS information.
type Info struct {
	Entry32 *Entry32
	Entry64 *Entry64
	Tables  []*Table
}

// ParseInfo parses SMBIOS information from binary data.
func ParseInfo(entryData, tableData []byte) (*Info, error) {
	info := &Info{}
	var err error
	info.Entry32, info.Entry64, err = ParseEntry(entryData)
	if err != nil {
		return nil, fmt.Errorf("error parsing entry point structure: %v", err)
	}
	for len(tableData) > 0 {
		t, remainder, err := ParseTable(tableData)
		if err != nil && err != errEndOfTable {
			return nil, err
		}
		info.Tables = append(info.Tables, t)
		tableData = remainder
	}
	return info, nil
}

// MajorVersion return major version of the SMBIOS spec.
func (i *Info) MajorVersion() uint8 {
	if i.Entry64 != nil {
		return i.Entry64.SMBIOSMajorVersion
	}
	if i.Entry32 != nil {
		return i.Entry32.SMBIOSMajorVersion
	}
	return 0
}

// MinorVersion return minor version of the SMBIOS spec.
func (i *Info) MinorVersion() uint8 {
	if i.Entry64 != nil {
		return i.Entry64.SMBIOSMinorVersion
	}
	if i.Entry32 != nil {
		return i.Entry32.SMBIOSMinorVersion
	}
	return 0
}

// DocRev return document revision of the SMBIOS spec.
func (i *Info) DocRev() uint8 {
	if i.Entry64 != nil {
		return i.Entry64.SMBIOSDocRev
	}
	return 0
}

// GetTablesByType returns tables of specific type.
func (i *Info) GetTablesByType(tt TableType) []*Table {
	var res []*Table
	for _, t := range i.Tables {
		if t.Type == tt {
			res = append(res, t)
		}
	}
	return res
}

// GetBIOSInformation returns the Bios Information (type 0) table, if present.
func (i *Info) GetBIOSInformation() (*BIOSInformation, error) {
	bt := i.GetTablesByType(TableTypeBIOSInformation)
	if len(bt) == 0 {
		return nil, ErrTableNotFound
	}
	// There can only be one of these.
	return NewBIOSInformation(bt[0])
}

// GetSystemInformation returns the System Information (type 1) table, if present.
func (i *Info) GetSystemInformation() (*SystemInformation, error) {
	bt := i.GetTablesByType(TableTypeSystemInformation)
	if len(bt) == 0 {
		return nil, ErrTableNotFound
	}
	// There can only be one of these.
	return NewSystemInformation(bt[0])
}

// GetBaseboardInformation returns all the Baseboard Information (type 2) tables present.
func (i *Info) GetBaseboardInformation() ([]*BaseboardInformation, error) {
	var res []*BaseboardInformation
	for _, t := range i.GetTablesByType(TableTypeBaseboardInformation) {
		bi, err := NewBaseboardInformation(t)
		if err != nil {
			return nil, err
		}
		res = append(res, bi)
	}
	return res, nil
}

// GetChassisInformation returns all the Chassis Information (type 3) tables present.
func (i *Info) GetChassisInformation() ([]*ChassisInformation, error) {
	var res []*ChassisInformation
	for _, t := range i.GetTablesByType(TableTypeChassisInformation) {
		ci, err := NewChassisInformation(t)
		if err != nil {
			return nil, err
		}
		res = append(res, ci)
	}
	return res, nil
}

// GetProcessorInformation returns all the Processor Information (type 4) tables present.
func (i *Info) GetProcessorInformation() ([]*ProcessorInformation, error) {
	var res []*ProcessorInformation
	for _, t := range i.GetTablesByType(TableTypeProcessorInformation) {
		pi, err := NewProcessorInformation(t)
		if err != nil {
			return nil, err
		}
		res = append(res, pi)
	}
	return res, nil
}

// GetCacheInformation returns all the Cache Information (type 7) tables present.
func (i *Info) GetCacheInformation() ([]*CacheInformation, error) {
	var res []*CacheInformation
	for _, t := range i.GetTablesByType(TableTypeCacheInformation) {
		ci, err := NewCacheInformation(t)
		if err != nil {
			return nil, err
		}
		res = append(res, ci)
	}
	return res, nil
}

// GetMemoryDevices returns all the Memory Device (type 17) tables present.
func (i *Info) GetMemoryDevices() ([]*MemoryDevice, error) {
	var res []*MemoryDevice
	for _, t := range i.GetTablesByType(TableTypeMemoryDevice) {
		ci, err := NewMemoryDevice(t)
		if err != nil {
			return nil, err
		}
		res = append(res, ci)
	}
	return res, nil
}

// GetIPMIDeviceInformation returns all the IPMI Device Information (type 38) tables present.
func (i *Info) GetIPMIDeviceInformation() ([]*IPMIDeviceInformation, error) {
	var res []*IPMIDeviceInformation
	for _, t := range i.GetTablesByType(TableTypeIPMIDeviceInformation) {
		d, err := NewIPMIDeviceInformation(t)
		if err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil
}

// GetTPMDevices returns all the TPM Device (type 43) tables present.
func (i *Info) GetTPMDevices() ([]*TPMDevice, error) {
	var res []*TPMDevice
	for _, t := range i.GetTablesByType(TableTypeTPMDevice) {
		d, err := NewTPMDevice(t)
		if err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil
}
