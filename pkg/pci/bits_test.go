// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import "testing"

func TestControlBits(t *testing.T) {
	var tests = []struct {
		c PCIControl
		w string
	}{
		{c: 0, w: "I/O- Memory- DMA- Special- MemWINV- VGASnoop- ParErr- Stepping- SERR- FastB2B- DisInt-"},
		{c: 0x001, w: "I/O+ Memory- DMA- Special- MemWINV- VGASnoop- ParErr- Stepping- SERR- FastB2B- DisInt-"},
		{c: 0x003, w: "I/O+ Memory+ DMA- Special- MemWINV- VGASnoop- ParErr- Stepping- SERR- FastB2B- DisInt-"},
		{c: 0x555, w: "I/O+ Memory- DMA+ Special- MemWINV+ VGASnoop- ParErr+ Stepping- SERR+ FastB2B- DisInt+"},
		{c: 0xaaa, w: "I/O- Memory+ DMA- Special+ MemWINV- VGASnoop+ ParErr- Stepping+ SERR- FastB2B+ DisInt-"},
		{c: 0xfff, w: "I/O+ Memory+ DMA+ Special+ MemWINV+ VGASnoop+ ParErr+ Stepping+ SERR+ FastB2B+ DisInt+"},
	}
	for _, tt := range tests {
		s := tt.c.String()
		if s != tt.w {
			t.Errorf("Control bits for %#x: got \n%q\n, want \n%q", tt.c, s, tt.w)
		}
	}

}

func TestStatusBits(t *testing.T) {
	var tests = []struct {
		c PCIStatus
		w string
	}{
		{c: 0, w: "INTx- Cap- 66MHz- UDF- FastB2b- ParErr- DEVSEL- DEVSEL=fast <MABORT- >SERR- <PERR-"},
		{c: 0x600, w: "INTx- Cap- 66MHz- UDF- FastB2b- ParErr- DEVSEL- DEVSEL=reserved <MABORT- >SERR- <PERR-"},
		{c: 0x400, w: "INTx- Cap- 66MHz- UDF- FastB2b- ParErr- DEVSEL- DEVSEL=slow <MABORT- >SERR- <PERR-"},
		{c: 0x200, w: "INTx- Cap- 66MHz- UDF- FastB2b- ParErr- DEVSEL- DEVSEL=medium <MABORT- >SERR- <PERR-"},
		{c: 0xffff, w: "INTx+ Cap+ 66MHz+ UDF+ FastB2b+ ParErr+ DEVSEL+ DEVSEL=reserved <MABORT+ >SERR+ <PERR+"},
	}
	for _, tt := range tests {
		s := tt.c.String()
		if s != tt.w {
			t.Errorf("Control bits for %#x: got \n%q, want \n%q", tt.c, s, tt.w)
		}
	}

}
