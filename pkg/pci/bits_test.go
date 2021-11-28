// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import "testing"

func TestControlBits(t *testing.T) {
	tests := []struct {
		c Control
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
	tests := []struct {
		c Status
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

func TestBAR(t *testing.T) {
	tests := []struct {
		bar string
		res string
		err string
	}{
		{bar: "0x0000000000001860 0x0000000000001867 0x0000000000040101", res: "Region 0: I/O ports at 1860 [size=8]"},
		{bar: "0x0000000000001814 0x0000000000001817 0x0000000000040101", res: "Region 1: I/O ports at 1814 [size=4]"},
		{bar: "0x0000000000001818 0x000000000000181f 0x0000000000040101", res: "Region 2: I/O ports at 1818 [size=8]"},
		{bar: "0x0000000000001810 0x0000000000001813 0x0000000000040101", res: "Region 3: I/O ports at 1810 [size=4]"},
		{bar: "0x0000000000001840 0x000000000000185f 0x0000000000040101", res: "Region 4: I/O ports at 1840 [size=32]"},
		{bar: "0x00000000f2827000 0x00000000f28277ff 0x0000000000040200", res: "Region 5: Memory at f2827000 (32-bit, non-prefetchable) [size=0x800]"},
		{bar: "0x0000000000000000 0x0000000000000000 0x0000000000000000", res: ""},
		{bar: "z 0x0000000000080000 0x0000000000000000", res: "", err: "strconv.ParseUint: parsing \"z\": invalid syntax"},
		{bar: " 0x0000000000080000 0x0000000000000000", res: "", err: "bar \" 0x0000000000080000 0x0000000000000000\" should have 3 fields"},
		//		{bar: "0x0000000000000000 0x0000000000000000 0x000000000000000f", res: "Region 1: Can't get type from \"0x0000000000000000 0x0000000000000000 0x000000000000000f\""},
		{bar: "0x00000000000c0000 0x00000000000dffff 0x0000000000000212", res: "Region 9: (Disabled)Expansion ROM at 000c0000 (low 1Mbyte) [size=0x20000]"},
		{bar: "0x00000000000c0001 0x00000000000dffff 0x0000000000000212", res: "Region 10: Expansion ROM at 000c0000 (low 1Mbyte) [size=0x20000]"},
		{bar: "0x0000000000080000 0x000000000008ffff 0x0000000000000212", res: "Region 11: Memory at 00080000 (32-bit, low 1Mbyte, non-prefetchable) [size=0x10000]"},
	}
	for i, tt := range tests {
		b, l, a, err := BaseLimType(tt.bar)
		if err != nil && err.Error() != tt.err {
			t.Errorf("BAR %s: got error \n%q, want \n%q", tt.bar, err, tt.err)
		}
		if err == nil && len(tt.err) != 0 {
			t.Errorf("BAR %s: got nil, want \n%q", tt.bar, tt.err)
		}
		bar := BAR{Index: i, Base: b, Lim: l, Attr: a}
		s := bar.String()
		if s != tt.res {
			t.Errorf("BAR %s: got \n%q, want \n%q", tt.bar, s, tt.res)
		}
	}
}
