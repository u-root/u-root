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

// func ParseDppMsgATAPI(h EfiDevicePathProtocolHdr, b []byte) (*DppMsgATAPI, error)
func TestParseDppMsgATAPI(t *testing.T) {
	in := []byte{0, 1, 2, 3}
	hdr := EfiDevicePathProtocolHdr{
		ProtoType:    3,
		ProtoSubType: 1,
		Length:       uint16(len(in) + 4),
	}
	p, err := ParseDppMsgATAPI(hdr, in)
	if err != nil {
		t.Fatal(err)
	}
	want := "ATAPI(pri=true,master=false,lun=770)"
	got := p.String()
	if want != got {
		t.Errorf("\nwant %s\n got %s", want, got)
	}
	wantp := "ATAPI"
	gotp := p.ProtoSubTypeStr()
	if wantp != gotp {
		t.Errorf("want %s got %s", wantp, gotp)
	}
}

// func ParseDppMsgMAC(h EfiDevicePathProtocolHdr, b []byte) (*DppMsgMAC, error)
func TestParseDppMsgMAC(t *testing.T) {
	in := []byte{
		0x00, 0x26, 0xfd, 0x00, 0x26, 0xfd, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x06,
	}
	hdr := EfiDevicePathProtocolHdr{
		ProtoType:    3,
		ProtoSubType: 11,
		Length:       uint16(len(in) + 4),
	}
	p, err := ParseDppMsgMAC(hdr, in)
	if err != nil {
		t.Fatal(err)
	}
	want := "MAC(mac=0026fd0026fd0000000000000000000000000000000000000000000000000000,iftype=0x6)"
	got := p.String()
	if want != got {
		t.Errorf("\nwant %s\n got %s", want, got)
	}
	wantp := "MAC"
	gotp := p.ProtoSubTypeStr()
	if wantp != gotp {
		t.Errorf("want %s got %s", wantp, gotp)
	}
}
