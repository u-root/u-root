// Copyright 2016-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"encoding"
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

// OverrideOpt is a function overriding the marshaler
type OverrideOpt func(over map[TableType]encoding.BinaryMarshaler)

// ReplaceTable returns func replacing the marshaler given table type
func ReplaceTable(typ TableType, table encoding.BinaryMarshaler) OverrideOpt {
	return func(over map[TableType]encoding.BinaryMarshaler) {
		over[typ] = table
	}
}

// ModifySystemInfo modifies the SystemInfo table in system memory
func (m *Modifier) ModifySystemInfo(manufacturer, productName, version, serialNumber string) error {
	sysInfo, err := m.Info.GetSystemInfo()
	if err != nil {
		return fmt.Errorf("failed to get system info: %w", err)
	}
	sysInfo.Manufacturer = manufacturer
	sysInfo.ProductName = productName
	sysInfo.Version = version
	sysInfo.SerialNumber = serialNumber

	entry, tables, err := m.Info.Marshal(ReplaceTable(TableTypeSystemInfo, sysInfo))
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
