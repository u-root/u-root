// Copyright 2016-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
	"os"
)

// Modifier modifies the SMBIOS data
type Modifier struct {
	Info
	memFile   *os.File
	entryAddr int64
	tableAddr int64
}

func getMemFile() (*os.File, error) {
	memFile, err := os.OpenFile("/dev/mem", os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open /dev/mem: %w", err)
	}
	return memFile, nil
}

func getEntries(smbiosBase func() (int64, int64, error), memFile *os.File) (*Entry32, *Entry64, int64, error) {
	var err error
	var entryAddr, sz int64
	entryAddr, sz, err = smbiosBase()
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to find SMBIOS base: %w", err)
	}
	entryData := make([]byte, sz)
	if _, err := memFile.ReadAt(entryData, entryAddr); err != nil {
		return nil, nil, 0, fmt.Errorf("failed to read entryData at address: 0x%x, error:%w", entryAddr, err)
	}

	e32, e64, err := ParseEntry(entryData)
	return e32, e64, entryAddr, err
}

// OverrideOpt is a function return overridden tables given another tables the marshaler
type OverrideOpt func(t []*Table) ([]*Table, error)

// ReplaceSystemInfo returns override options of the SystemInfo table with the given values
func ReplaceSystemInfo(manufacturer, productName, version, serialNumber, sku, family *string, uuid *UUID, wakeupType *WakeupType) OverrideOpt {
	return func(tables []*Table) ([]*Table, error) {
		var result []*Table

		for _, t := range tables {
			if t.Type != TableTypeSystemInfo {
				result = append(result, t)
				continue
			}
			// replace it
			si, err := ParseSystemInfo(t)
			if err != nil {
				return nil, fmt.Errorf("failed to parse system info: %w", err)
			}

			if manufacturer != nil {
				si.Manufacturer = *manufacturer
			}
			if productName != nil {
				si.ProductName = *productName
			}
			if version != nil {
				si.Version = *version
			}
			if serialNumber != nil {
				si.SerialNumber = *serialNumber
			}
			if sku != nil {
				si.SKUNumber = *sku
			}
			if family != nil {
				si.Family = *family
			}
			if uuid != nil {
				si.UUID = *uuid
			}
			if wakeupType != nil {
				si.WakeupType = *wakeupType
			}

			sit, err := si.toTable()
			if err != nil {
				return nil, fmt.Errorf("failed to convert system info to table: %w", err)
			}
			result = append(result, sit)
		}
		return result, nil
	}
}

// ReplaceBaseboardInfoMotherboard returns override options that only overrides table with Type = BaseboardInfo and BoardType = BoardTypeMotherboardIncludesProcessorMemoryAndIO
func ReplaceBaseboardInfoMotherboard(manufacturer, product, version, serialNumber, assetTag, locationInChassis *string, boardFeatures *BoardFeatures, chassisHandle *uint16, boardType *BoardType, containedObjectHandles *[]uint16) OverrideOpt {
	return func(tables []*Table) ([]*Table, error) {
		var result []*Table
		for _, t := range tables {
			if t.Type != TableTypeBaseboardInfo {
				result = append(result, t)
				continue
			}

			bi, err := ParseBaseboardInfo(t)
			if err != nil {
				return nil, fmt.Errorf("failed to parse baseboard info: %w", err)
			}
			if bi.BoardType != BoardTypeMotherboardIncludesProcessorMemoryAndIO {
				result = append(result, t)
				continue
			}

			// replace it
			if manufacturer != nil {
				bi.Manufacturer = *manufacturer
			}
			if product != nil {
				bi.Product = *product
			}
			if version != nil {
				bi.Version = *version
			}
			if serialNumber != nil {
				bi.SerialNumber = *serialNumber
			}
			if assetTag != nil {
				bi.AssetTag = *assetTag
			}
			if locationInChassis != nil {
				bi.LocationInChassis = *locationInChassis
			}
			if boardFeatures != nil {
				bi.BoardFeatures = *boardFeatures
			}
			if chassisHandle != nil {
				bi.ChassisHandle = *chassisHandle
			}
			if boardType != nil {
				bi.BoardType = *boardType
			}
			if containedObjectHandles != nil {
				bi.NumberOfContainedObjectHandles = uint8(len(*containedObjectHandles))
				bi.ContainedObjectHandles = *containedObjectHandles
			}
			biT, err := bi.toTable()
			if err != nil {
				return nil, fmt.Errorf("failed to convert baseboard info to table: %w", err)
			}
			result = append(result, biT)
		}
		return result, nil
	}
}

// Modify modifies SMBIOS tables in system memory given override options
func (m *Modifier) Modify(opts ...OverrideOpt) error {
	entry, tables, err := m.Info.Marshal(opts...)
	if err != nil {
		return fmt.Errorf("failed to marshal info: %w", err)
	}

	if _, err := m.memFile.WriteAt(entry, m.entryAddr); err != nil {
		return fmt.Errorf("failed to write entry data at address: 0x%x, error:%w", m.entryAddr, err)
	}
	if _, err = m.memFile.WriteAt(tables, m.tableAddr); err != nil {
		return fmt.Errorf("failed to write table data at address: 0x%x, error:%w", m.tableAddr, err)
	}
	return nil
}

// CloseMemFile closes Modifier memory file
func (m *Modifier) CloseMemFile() error {
	return m.memFile.Close()
}

// NewModifier returns a Modifier and initialize all necessary fields
func NewModifier() (*Modifier, error) {
	return newModifier(getMemFile, SMBIOSBase)
}

func newModifier(getMemFile func() (*os.File, error), smbiosBase func() (int64, int64, error)) (*Modifier, error) {
	var err error
	m := &Modifier{}
	m.memFile, err = getMemFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get mem file: %w", err)
	}
	m.Entry32, m.Entry64, m.entryAddr, err = getEntries(smbiosBase, m.memFile)
	if err != nil {
		return nil, err
	}

	var entryData, tableData []byte
	if m.Entry32 != nil {
		m.tableAddr = int64(m.Entry32.StructTableAddr)
		tableData = make([]byte, m.Entry32.StructTableLength)
		entryData, err = m.Entry32.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Entry32: %w", err)
		}
	}
	if m.Entry64 != nil {
		m.tableAddr = int64(m.Entry64.StructTableAddr)
		tableData = make([]byte, m.Entry64.StructMaxSize)
		entryData, err = m.Entry64.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal Entry64: %w", err)
		}
	}
	if _, err := m.memFile.ReadAt(tableData, m.tableAddr); err != nil {
		return nil, fmt.Errorf("failed to ReadAt table from address: 0x%x, error:%w", m.tableAddr, err)
	}

	// load data
	info, err := ParseInfo(entryData, tableData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse info: %w", err)
	}
	m.Info = *info
	return m, nil
}
