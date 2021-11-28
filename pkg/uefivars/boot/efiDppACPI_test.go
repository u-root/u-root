// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: BSD-3-Clause
//

package boot

import (
	"testing"
)

// func ParseDppAcpiExDevPath(h EfiDevicePathProtocolHdr, b []byte) (*DppAcpiExDevPath, error)
func TestParseDppAcpiExDevPath(t *testing.T) {
	in := []byte{
		0, 0, 0, 1,
		0, 0, 0, 2,
		0, 0, 0, 3,
	}
	// not sure what's supposed to go into these fields, so use made-up values
	in = append(in, cstr("HIDSTR")...)
	in = append(in, cstr("UIDSTR")...)
	in = append(in, cstr("CIDSTR")...)
	hdr := EfiDevicePathProtocolHdr{
		ProtoType:    2,
		ProtoSubType: 2,
		Length:       uint16(len(in) + 4),
	}

	p, err := ParseDppAcpiExDevPath(hdr, in)
	if err != nil {
		t.Error(err)
	} else {
		want := "ACPI_EX(0x00000001,0x00000002,0x00000003,HIDSTR,UIDSTR,CIDSTR)"
		got := p.String()
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
		gotp := p.ProtoSubTypeStr()
		wantp := "Expanded Device Path"
		if gotp != wantp {
			t.Errorf("want %s, got %s", wantp, gotp)
		}
	}
}

func cstr(s string) []byte { return append([]byte(s), 0) }

// func ParseDppAcpiDevPath(h EfiDevicePathProtocolHdr, b []byte) (*DppAcpiDevPath, error)
func TestParseDppAcpiDevPath(t *testing.T) {
	in := []byte{
		0, 0, 0, 1,
		0, 0, 0, 2,
	}
	hdr := EfiDevicePathProtocolHdr{
		ProtoType:    2,
		ProtoSubType: 4,
		Length:       uint16(len(in) + 4),
	}
	p, err := ParseDppAcpiDevPath(hdr, in)
	if err != nil {
		t.Error(err)
	} else {
		want := "ACPI(0x00000001,0x00000002)"
		got := p.String()
		if got != want {
			t.Errorf("want %s, got %s", want, got)
		}
		gotp := p.ProtoSubTypeStr()
		wantp := "NVDIMM"
		if gotp != wantp {
			t.Errorf("want %s, got %s", wantp, gotp)
		}
	}
}
