// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package smbios

import "testing"

var validTableHeaderRaw = []byte{0x0, 0xFF, 0xBE, 0xEF}

func TestParseTableHeader(t *testing.T) {
	var h Header
	if err := h.Parse(validTableHeaderRaw); err != nil {
		t.Error(err)
	}
}

func TestHeaderToString(t *testing.T) {
	h := Header{
		Type:   0x1,
		Length: 0xFF,
		Handle: 0xBEEF,
	}

	expect := `Handle 0xBEEF, DMI type 1, 255 bytes
System Information`

	gotString := h.String()

	if gotString != expect {
		t.Errorf("Header.String(): %v, want %v", gotString, expect)
	}
}
