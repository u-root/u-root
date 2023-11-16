// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
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
