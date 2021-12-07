// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package abi

import (
	"testing"

	"golang.org/x/sys/unix"
)

// From the cros docs
// Bits	 Content	 Notes
// 63-57	 Unused
// 56	 Successful Boot Flag
// Set to 1 the first time the system has successfully booted from this partition (see the File System/Autoupdate design document for the definition of success).
// 55-52	 Tries Remaining	Number of times to attempt booting this partition. Used only when the Successful Boot Flag is 0.
// 51-48	 Priority	4-bit number: 15 = highest, 1 = lowest, 0 = not bootable.
// 47-0	 Reserved by EFI Spec
func TestCrosGPT(t *testing.T) {
	tests := []struct {
		n string
		f FlagSet
		v uint64
		o string
	}{
		{
			n: "Basic ChromeOS GPT flags",
			f: FlagSet{
				&BitFlag{Name: "Suc", Value: 1 << 56},
				&Field{Name: "Tries", BitMask: 15 << 52, Shift: 52},
				&Field{Name: "prio", BitMask: 15 << 48, Shift: 48},
			},
			v: 1<<56 | 8<<52 | 3<<48,
			o: "Suc|Tries=0x8|prio=0x3",
		},
		{
			n: "Basic ChromeOS GPT flags with standard GPT value",
			f: FlagSet{
				&BitFlag{Name: "Suc", Value: 1 << 56},
				&Field{Name: "Tries", BitMask: 15 << 52, Shift: 52},
				&Field{Name: "prio", BitMask: 15 << 48, Shift: 48},
			},
			v: 1<<56 | 8<<52 | 3<<48 | 1,
			o: "Suc|Tries=0x8|prio=0x3|0x1",
		},
		{
			n: "Simple system call",
			f: FlagSet{
				&Value{Name: "write", Value: unix.SYS_WRITE},
			},
			v: unix.SYS_WRITE,
			o: "write",
		},
	}
	for _, tc := range tests {
		t.Run(tc.n, func(t *testing.T) {
			s := tc.f.Parse(tc.v)
			if s != tc.o {
				t.Fatalf("Got %s, want %s", s, tc.o)
			}
		})
	}
}
