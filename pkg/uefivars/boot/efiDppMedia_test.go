// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package boot

import (
	"strings"
	"testing"
	"unicode/utf16"
)

// func ParseDppMediaHdd(h EfiDevicePathProtocolHdr, b []byte) (*DppMediaHDD, error)
func TestParseDppMediaHdd(t *testing.T) {
	in := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		0x20, 0x21, 0x22, 0x23, 0x02, 0x02,
	}
	hdr := EfiDevicePathProtocolHdr{
		ProtoType:    4,
		ProtoSubType: 1,
		Length:       uint16(len(in) + 4),
	}
	p, err := ParseDppMediaHdd(hdr, in)
	if err != nil {
		t.Fatal(err)
	}
	want := "HD(50462976,GPT,17161514-1918-1b1a-1c1d-1e1f20212223,0xb0a090807060504,0x131211100f0e0d0c)"
	got := p.String()
	if want != got {
		t.Errorf("\nwant %s\n got %s", want, got)
	}
	wantp := "HDD"
	gotp := p.ProtoSubTypeStr()
	if wantp != gotp {
		t.Errorf("want %s got %s", wantp, gotp)
	}
}

// func ParseDppMediaFilePath(h EfiDevicePathProtocolHdr, b []byte) (*DppMediaFilePath, error)
func TestParseDppMediaFilePath(t *testing.T) {
	str := `blah\blah\blah.efi`
	// convert to utf16 ([]uint16)
	u16 := utf16.Encode([]rune(str))
	//...and then to []byte
	var in []byte
	for _, u := range u16 {
		in = append(in, byte(u&0xff), byte(u>>8&0xff))
	}
	hdr := EfiDevicePathProtocolHdr{
		ProtoType:    4,
		ProtoSubType: 4,
		Length:       uint16(len(in) + 4),
	}
	p, err := ParseDppMediaFilePath(hdr, in)
	if err != nil {
		t.Fatal(err)
	}
	want := "File(" + strings.ReplaceAll(str, "\\", "/") + ")"
	got := p.String()
	if want != got {
		t.Errorf("\nwant %s\n got %s", want, got)
	}
	wantp := "FilePath"
	gotp := p.ProtoSubTypeStr()
	if wantp != gotp {
		t.Errorf("want %s got %s", wantp, gotp)
	}
}

// func ParseDppMediaPIWGFV(h EfiDevicePathProtocolHdr, b []byte) (*DppMediaPIWGFV, error)
func TestParseDppMediaPIWGFV(t *testing.T) {
	in := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	hdr := EfiDevicePathProtocolHdr{
		ProtoType:    4,
		ProtoSubType: 7,
		Length:       uint16(len(in) + 4),
	}
	p, err := ParseDppMediaPIWGFV(hdr, in)
	if err != nil {
		t.Fatal(err)
	}
	want := "Fv(03020100-0504-0706-0809-0a0b0c0d0e0f)"
	got := p.String()
	if want != got {
		t.Errorf("\nwant %s\n got %s", want, got)
	}
	wantp := "PIWG Firmware Volume"
	gotp := p.ProtoSubTypeStr()
	if wantp != gotp {
		t.Errorf("want %s got %s", wantp, gotp)
	}
}

// func ParseDppMediaPIWGFF(h EfiDevicePathProtocolHdr, b []byte) (*DppMediaPIWGFF, error)
func TestParseDppMediaPIWGFF(t *testing.T) {
	in := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	hdr := EfiDevicePathProtocolHdr{
		ProtoType:    4,
		ProtoSubType: 6,
		Length:       uint16(len(in) + 4),
	}
	p, err := ParseDppMediaPIWGFF(hdr, in)
	if err != nil {
		t.Fatal(err)
	}
	want := "FvFile(03020100-0504-0706-0809-0a0b0c0d0e0f)"
	got := p.String()
	if want != got {
		t.Errorf("\nwant %s\n got %s", want, got)
	}
	wantp := "PIWG Firmware File"
	gotp := p.ProtoSubTypeStr()
	if wantp != gotp {
		t.Errorf("want %s got %s", wantp, gotp)
	}
}
