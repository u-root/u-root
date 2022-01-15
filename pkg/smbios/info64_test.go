// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"io/ioutil"
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
	data, err := ioutil.ReadFile("./testdata/smbios_table.bin")
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
	data, err := ioutil.ReadFile("./testdata/smbios_table.bin")
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
