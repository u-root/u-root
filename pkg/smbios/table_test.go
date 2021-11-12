// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package smbios

import (
	"io/ioutil"
	"testing"
)

var (
	testbinary = "testdata/satellite_pro_l70_testdata.bin"
)

func TestParseSMBIOS(t *testing.T) {
	data, err := ioutil.ReadFile(testbinary)
	if err != nil {
		t.Error(err)
	}
	datalen := len(data)
	readlen := 0
	for i := 0; datalen > i; i += readlen {
		_, rest, err := ParseTable(data)
		if err != nil {
			t.Log(err)
		}
		readlen = datalen - len(rest)
	}
}
