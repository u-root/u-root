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

// func ParseDppHwPci(h EfiDevicePathProtocolHdr, b []byte) (*DppHwPci, error)
func TestParseDppHwPci(t *testing.T) {
	in := []byte{3, 4}
	hdr := EfiDevicePathProtocolHdr{
		ProtoType:    1,
		ProtoSubType: 6,
		Length:       uint16(len(in) + 4),
	}

	pci, err := ParseDppHwPci(hdr, in)
	if err != nil {
		t.Fatal(err)
	}
	want := "PCI(0x3,0x4)"
	got := pci.String()
	if want != got {
		t.Errorf("want %s got %s", want, got)
	}
	wantp := "BMC"
	gotp := pci.ProtoSubTypeStr()
	if wantp != gotp {
		t.Errorf("want %s got %s", wantp, gotp)
	}
}
