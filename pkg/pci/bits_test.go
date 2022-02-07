// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import "testing"

func TestControlBits(t *testing.T) {
	for _, tt := range []struct {
		name    string
		control Control
		want    string
	}{
		{
			name:    "test 0x0",
			control: 0x0,
			want:    "I/O- Memory- DMA- Special- MemWINV- VGASnoop- ParErr- Stepping- SERR- FastB2B- DisInt-",
		},
		{
			name:    "test 0x001",
			control: 0x001,
			want:    "I/O+ Memory- DMA- Special- MemWINV- VGASnoop- ParErr- Stepping- SERR- FastB2B- DisInt-",
		},
		{
			name:    "test 0x003",
			control: 0x003,
			want:    "I/O+ Memory+ DMA- Special- MemWINV- VGASnoop- ParErr- Stepping- SERR- FastB2B- DisInt-",
		},
		{
			name:    "test 0x555",
			control: 0x555,
			want:    "I/O+ Memory- DMA+ Special- MemWINV+ VGASnoop- ParErr+ Stepping- SERR+ FastB2B- DisInt+",
		},
		{
			name:    "test 0xaaa",
			control: 0xaaa,
			want:    "I/O- Memory+ DMA- Special+ MemWINV- VGASnoop+ ParErr- Stepping+ SERR- FastB2B+ DisInt-",
		},
		{
			name:    "test 0xfff",
			control: 0xfff,
			want:    "I/O+ Memory+ DMA+ Special+ MemWINV+ VGASnoop+ ParErr+ Stepping+ SERR+ FastB2B+ DisInt+",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.control.String()
			if got != tt.want {
				t.Errorf("Control bits for '%#x': got: %q, want: %q", tt.control, got, tt.want)
			}
		})
	}
}

func TestStatusBits(t *testing.T) {
	for _, tt := range []struct {
		name   string
		status Status
		want   string
	}{
		{
			name:   "test 0x0",
			status: 0x0,
			want:   "INTx- Cap- 66MHz- UDF- FastB2b- ParErr- DEVSEL- DEVSEL=fast <MABORT- >SERR- <PERR-",
		},
		{
			name:   "test 0x600",
			status: 0x600,
			want:   "INTx- Cap- 66MHz- UDF- FastB2b- ParErr- DEVSEL- DEVSEL=reserved <MABORT- >SERR- <PERR-",
		},
		{
			name:   "test 0x400",
			status: 0x400,
			want:   "INTx- Cap- 66MHz- UDF- FastB2b- ParErr- DEVSEL- DEVSEL=slow <MABORT- >SERR- <PERR-",
		},
		{
			name:   "test 0x200",
			status: 0x200,
			want:   "INTx- Cap- 66MHz- UDF- FastB2b- ParErr- DEVSEL- DEVSEL=medium <MABORT- >SERR- <PERR-",
		},
		{
			name:   "test 0xffff",
			status: 0xffff,
			want:   "INTx+ Cap+ 66MHz+ UDF+ FastB2b+ ParErr+ DEVSEL+ DEVSEL=reserved <MABORT+ >SERR+ <PERR+",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.String()
			if got != tt.want {
				t.Errorf("Control bits for '%#x': got: %q, want: %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestBAR(t *testing.T) {
	for i, tt := range []struct {
		name string
		bar  string
		want string
		err  string
	}{
		{
			name: "test1",
			bar:  "0x0000000000001860 0x0000000000001867 0x0000000000040101",
			want: "Region 0: I/O ports at 1860 [size=8]",
		},
		{
			name: "test2",
			bar:  "0x0000000000001814 0x0000000000001817 0x0000000000040101",
			want: "Region 1: I/O ports at 1814 [size=4]",
		},
		{
			name: "test3",
			bar:  "0x0000000000001818 0x000000000000181f 0x0000000000040101",
			want: "Region 2: I/O ports at 1818 [size=8]",
		},
		{
			name: "test4",
			bar:  "0x0000000000001810 0x0000000000001813 0x0000000000040101",
			want: "Region 3: I/O ports at 1810 [size=4]",
		},
		{
			name: "test5",
			bar:  "0x0000000000001840 0x000000000000185f 0x0000000000040101",
			want: "Region 4: I/O ports at 1840 [size=32]",
		},
		{
			name: "test6",
			bar:  "0x00000000f2827000 0x00000000f28277ff 0x0000000000040200",
			want: "Region 5: Memory at f2827000 (32-bit, non-prefetchable) [size=0x800]",
		},
		{
			name: "test7",
			bar:  "0x0000000000000000 0x0000000000000000 0x0000000000000000",
			want: "",
		},
		{
			name: "test8",
			bar:  "z 0x0000000000080000 0x0000000000000000",
			want: "",
			err:  "strconv.ParseUint: parsing \"z\": invalid syntax",
		},
		{
			name: "test9",
			bar:  " 0x0000000000080000 0x0000000000000000",
			want: "",
			err:  "bar \" 0x0000000000080000 0x0000000000000000\" should have 3 fields",
		},
		{
			name: "test10",
			bar:  "0x00000000000c0000 0x00000000000dffff 0x0000000000000212",
			want: "Region 9: (Disabled)Expansion ROM at 000c0000 (low 1Mbyte) [size=0x20000]",
		},
		{
			name: "test11",
			bar:  "0x00000000000c0001 0x00000000000dffff 0x0000000000000212",
			want: "Region 10: Expansion ROM at 000c0000 (low 1Mbyte) [size=0x20000]",
		},
		{
			name: "test12",
			bar:  "0x0000000000080000 0x000000000008ffff 0x0000000000000212",
			want: "Region 11: Memory at 00080000 (32-bit, low 1Mbyte, non-prefetchable) [size=0x10000]",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			base, limit, attribute, err := BaseLimType(tt.bar)
			if err != nil && err.Error() != tt.err {
				t.Errorf("BAR %s: got error: %q, want %q", tt.bar, err, tt.err)
			}
			if err == nil && len(tt.err) != 0 {
				t.Errorf("BAR %s: got: 'nil', want %q", tt.bar, tt.err)
			}
			bar := BAR{Index: i, Base: base, Lim: limit, Attr: attribute}
			got := bar.String()
			if got != tt.want {
				t.Errorf("BAR %s: got: %q, want: %q", tt.bar, got, tt.want)
			}
		})
	}
}
