// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test64ParseInfo(t *testing.T) {
	info, err := setupMockData()
	if err != nil {
		t.Errorf("error parsing info data: %v", err)
	}
	if info.Entry32 != nil {
		t.Errorf("false detection of 32-bit SMBIOS table header")
	}
}

func Test64ParseInfoHeaderMalformed(t *testing.T) {
	data, err := os.ReadFile("./testdata/smbios_table.bin")
	if err != nil {
		t.Errorf("error reading mockup smbios tables: %v", err)
	}

	entryData := data[:10]
	data = data[32:]

	_, err = ParseInfo(entryData, data)
	if err == nil {
		t.Errorf("error parsing info data: %v", err)
	}
}

func Test64MajorVersion(t *testing.T) {
	info, err := setupMockData()
	if err != nil {
		t.Errorf("error parsing info data: %v", err)
	}
	if info.MajorVersion() != 3 {
		t.Errorf("major version should be 3 - got %d", info.MajorVersion())
	}
}

func Test64MinorVersion(t *testing.T) {

	info, err := setupMockData()
	if err != nil {
		t.Errorf("error parsing info data: %v", err)
	}
	if info.MinorVersion() != 1 {
		t.Errorf("minor version should be 1 - got %d", info.MinorVersion())
	}
}

func Test64DocRev(t *testing.T) {

	info, err := setupMockData()
	if err != nil {
		t.Errorf("error parsing info data: %v", err)
	}
	if info.DocRev() != 1 {
		t.Errorf("doc revision should be 1 - got %d", info.DocRev())
	}
}

func Test64GetTablesByType(t *testing.T) {

	info, err := setupMockData()
	if err != nil {
		t.Errorf("error parsing info data: %v", err)
	}

	table := info.GetTablesByType(TableTypeBIOSInfo)
	if table == nil {
		t.Errorf("unable to get type")
	}
	if table != nil {
		if table[0].Header.Type != TableTypeBIOSInfo {
			t.Errorf("Wrong type. Got %v but want %v", TableTypeBIOSInfo, table[0].Header.Type)
		}
	}
}

func systemInfoTable() (*Table, error) {
	info, err := setupMockData()
	if err != nil {
		return nil, err
	}

	for _, t := range info.Tables {
		if t.Header.Type == TableTypeSystemInfo {
			return t, nil
		}
	}
	return nil, fmt.Errorf("Unable to find type 1 table")
}

func Test64Marshal(t *testing.T) {
	vanillaTable, err := systemInfoTable()
	if err != nil {
		t.Fatalf("Error setup mock system info table: %v", err)
	}
	vanillaTableRaw, err := vanillaTable.MarshalBinary()
	if err != nil {
		t.Fatalf("Error marshalling vanilla table: %v", err)
	}

	longTable, err := systemInfoTable()
	if err != nil {
		t.Fatalf("Error setup mock system info table: %v", err)
	}
	longTable.strings[0] += "a-very-long-string"

	shortTable, err := systemInfoTable()
	if err != nil {
		t.Fatalf("Error setup mock system info table: %v", err)
	}
	shortTable.strings[0] = "-"

	tests := []struct {
		name          string
		replacedTable *Table
		wantEntry     []byte
	}{
		{
			name:          "SameTable",
			replacedTable: vanillaTable,
			wantEntry:     []byte{95, 83, 77, 51, 95, 45, 24, 3, 1, 1, 1, 0, 248, 12, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:          "LongerTable",
			replacedTable: longTable,
			wantEntry:     []byte{95, 83, 77, 51, 95, 26, 24, 3, 1, 1, 1, 0, 10, 13, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:          "ShorterTable",
			replacedTable: shortTable,
			wantEntry:     []byte{95, 83, 77, 51, 95, 50, 24, 3, 1, 1, 1, 0, 243, 12, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := setupMockData()
			if err != nil {
				t.Fatalf("Error setup mock data: %v", err)
			}
			_, originalTablesRaw, err := setupMockRawData()
			if err != nil {
				t.Fatalf("Error setup mock raw data: %v", err)
			}

			b, err := tt.replacedTable.MarshalBinary()
			if err != nil {
				t.Fatalf("Error marshalling replaced table: %v", err)
			}
			if bytes.Index(originalTablesRaw, vanillaTableRaw) == -1 {
				t.Fatal("table raw should contain original table bytes")
			}
			wantTables := bytes.Replace(originalTablesRaw, vanillaTableRaw, b, -1)

			gotEntry, gotTables, err := info.Marshal(ReplaceTable(TableTypeSystemInfo, tt.replacedTable))
			if err != nil {
				t.Fatalf("Error marshalling info: %v", err)
			}

			if !bytes.Equal(tt.wantEntry, gotEntry) {
				t.Errorf("Wrong marshalled entry, want: %#v\ngot:%#v", tt.wantEntry, gotEntry)
			}
			if !bytes.Equal(wantTables, gotTables) {
				t.Errorf("Wrong marshalled tables, want: %#v\ngot:%#v", wantTables, gotTables)
			}
		})
	}
}

func setupMockData() (*Info, error) {
	data, err := os.ReadFile("./testdata/smbios_table.bin")
	if err != nil {
		return nil, err
	}

	entryData := data[:32]
	data = data[32:]

	info, err := ParseInfo(entryData, data)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func setupMockRawData() ([]byte, []byte, error) {
	data, err := os.ReadFile("./testdata/smbios_table.bin")
	if err != nil {
		return nil, nil, err
	}

	entryData := make([]byte, 24) // Entry64 length should be 24 instead of 32, the test data padded with zeros from 25th to 32nd byte.
	tableData := make([]byte, len(data)-32)
	copy(entryData, data[:24])
	copy(tableData, data[32:])
	return entryData, tableData, nil
}

func FuzzParseInfo(f *testing.F) {
	seeds, err := filepath.Glob("testdata/*.bin")
	if err != nil {
		f.Fatalf("failed to find seed corpora files: %v", err)
	}

	for _, seed := range seeds {
		seedBytes, err := os.ReadFile(seed)
		if err != nil {
			f.Fatalf("failed read seed corpora from files %v: %v", seed, err)
		}

		f.Add(seedBytes)
	}

	f.Fuzz(func(t *testing.T, data []byte) {

		if len(data) < 64 || len(data) > 4096 {
			return
		}

		entryData := data[:32]
		data = data[32:]

		info, err := ParseInfo(entryData, data)
		if err != nil {
			return
		}

		var entry []byte
		if info.Entry32 != nil {
			entry, err = info.Entry32.MarshalBinary()
		} else if info.Entry64 != nil {
			entry, err = info.Entry64.MarshalBinary()

		} else {
			t.Fatalf("expected a SMBIOS 32-Bit or 64-Bit entry point but got none")
		}

		if err != nil {
			t.Fatalf("failed to unmarshal entry data")
		}

		reparsedInfo, err := ParseInfo(entry, data)
		if err != nil {
			t.Fatalf("failed to reparse the SMBIOS info struct")
		}
		if !reflect.DeepEqual(info, reparsedInfo) {
			t.Errorf("expected: %#v\ngot:%#v", info, reparsedInfo)
		}
	})
}
