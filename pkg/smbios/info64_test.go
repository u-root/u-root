// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func defaultType1Table() *Table {
	return &Table{
		Header: Header{
			Type:   TableTypeSystemInfo,
			Length: 27,
			Handle: 0,
		},
		data:    []byte{1, 27, 0, 0, 1, 2, 3, 4, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 1, 2, 3, 0, 5, 6},
		strings: []string{"Manufacturer", "ProductName", "Version", "SerialNumber", "SKUNumber", "Family"},
	}
}

func defaultType2Table() *Table {
	return &Table{
		Header: Header{
			Type:   TableTypeBaseboardInfo,
			Length: 17,
			Handle: 0,
		},
		data: []byte{
			2, 17, 0, 0, 1, 2, 3, 4, 5, 0, 6, 0, 0,
			10, // BoardTypeMotherboardIncludesProcessorMemoryAndIO
			1, 10, 0,
		},
		strings: []string{"Manufacturer", "Product", "Version", "1234-5678", "8765-4321", "Location"},
	}
}

func type1Table(manufacturer string) *Table {
	t1 := defaultType1Table()
	t1.strings[0] = manufacturer
	return t1
}

func type2Table(asset string) *Table {
	t2 := defaultType2Table()
	t2.strings[4] = asset
	return t2
}

func getRawT(t *testing.T, table *Table) []byte {
	result, err := table.MarshalBinary()
	if err != nil {
		t.Fatalf("Error marshalling table: %v", err)
	}
	return result
}

func Test64Marshal(t *testing.T) {
	oldType1 := defaultType1Table()
	oldType1Raw := getRawT(t, oldType1)

	oldType2 := defaultType2Table()
	oldType2Raw := getRawT(t, oldType2)

	newTag := "new-tag"
	newType2 := type2Table(newTag)
	newType2Raw := getRawT(t, newType2)
	baseboardOpt := ReplaceBaseboardInfoMotherboard(nil, nil, nil, nil, &newTag, nil, nil, nil, nil, nil)

	longStr := "a-very-loooooooooooooooooooooooooooooooooooong-string"
	longType1 := type1Table(longStr)
	longType1Raw := getRawT(t, longType1)
	longType1Opt := ReplaceSystemInfo(&longStr, nil, nil, nil, nil, nil, nil, nil)

	shortStr := "-"
	shortType1 := type1Table(shortStr)
	shortType1Raw := getRawT(t, shortType1)
	shortType1Opt := ReplaceSystemInfo(&shortStr, nil, nil, nil, nil, nil, nil, nil)

	tests := []struct {
		name      string
		options   []OverrideOpt
		wantTable []byte
		wantEntry []byte
	}{
		{
			name:      "SameTable",
			wantTable: joinBytesT(t, oldType1Raw, oldType2Raw),
			wantEntry: []byte{95, 83, 77, 51, 95, 180, 24, 2, 1, 1, 0, 0, 167, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			name:      "LongerTable",
			options:   []OverrideOpt{longType1Opt},
			wantTable: joinBytesT(t, longType1Raw, oldType2Raw),
			wantEntry: []byte{95, 83, 77, 51, 95, 139, 24, 2, 1, 1, 0, 0, 208, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			name:      "ShorterTable",
			options:   []OverrideOpt{shortType1Opt},
			wantTable: joinBytesT(t, shortType1Raw, oldType2Raw),
			wantEntry: []byte{95, 83, 77, 51, 95, 191, 24, 2, 1, 1, 0, 0, 156, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			name:      "2 OverrideOpts",
			options:   []OverrideOpt{shortType1Opt, baseboardOpt},
			wantTable: joinBytesT(t, shortType1Raw, newType2Raw),
			wantEntry: []byte{95, 83, 77, 51, 95, 193, 24, 2, 1, 1, 0, 0, 154, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &Info{
				Entry64: defaultEntry64(),
				Tables:  []*Table{defaultType1Table(), defaultType2Table()},
			}

			gotEntry, gotTable, err := info.Marshal(tt.options...)
			if err != nil {
				t.Fatalf("Error marshalling info: %v", err)
			}

			if diff := cmp.Diff(gotEntry, tt.wantEntry); diff != "" {
				t.Errorf("Wrong marshalled entry, diff (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(gotTable, tt.wantTable); diff != "" {
				t.Errorf("Wrong marshalled tables, diff (-want +got):\n%s", diff)
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
