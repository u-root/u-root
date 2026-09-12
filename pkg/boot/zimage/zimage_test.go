// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zimage

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

var testData = &ZImage{
	Header: Header{
		Magic:      0x16f2818,
		Start:      0x0,
		End:        0xd5638,
		Endianness: 0x4030201,
		TableMagic: 0x45454545,
		TableAddr:  0x25bc,
	},
	Table: []TableEntry{
		{
			Tag:  0x5a534c4b,
			Data: []uint32{0xd55f5, 0x2b83c},
		},
	},
}

func TestParse(t *testing.T) {
	f, err := os.Open("testdata/zImage")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	z, err := Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(z, testData) {
		t.Errorf(`Parse("testdata/zImage") = %#v; want %#v`, z, testData)
	}
}

func TestKernelSizes(t *testing.T) {
	piggySizeAddr, kernelBSSSize, err := testData.GetKernelSizes()
	if err != nil {
		t.Fatal(err)
	}
	if piggySizeAddr != 0xd55f5 {
		t.Errorf("want piggySizeAddr=0xd55f5, got piggySizeAddr=%#x", piggySizeAddr)
	}
	if kernelBSSSize != 0x2b83c {
		t.Errorf("want kernelBSSSize=0x2b83c, got kernelBSSSize=%#x", kernelBSSSize)
	}
}

// Linux has appended TEXT_OFFSET and MALLOC_SIZE to the kernel size entry
// since this parser was written, so the entry is now six words rather than
// two. Read the first two and ignore the rest, as kexec-tools does.
//
// The values are from a real 6.12.45-bone34 zImage on a BeagleBone Black.
func TestKernelSizesLongEntry(t *testing.T) {
	z := &ZImage{
		Table: []TableEntry{
			{
				Tag: TagKernelSize,
				Data: []uint32{
					0x875534, // __piggy_size_addr - _start
					0x8c9a0,  // _kernel_bss_size
					0x8000,   // TEXT_OFFSET
					0x10000,  // MALLOC_SIZE
					0,
					0x5b1c1ca1,
				},
			},
		},
	}

	piggySizeAddr, kernelBSSSize, err := z.GetKernelSizes()
	if err != nil {
		t.Fatalf("GetKernelSizes() = %v, want nil", err)
	}
	if piggySizeAddr != 0x875534 {
		t.Errorf("want piggySizeAddr=0x875534, got %#x", piggySizeAddr)
	}
	if kernelBSSSize != 0x8c9a0 {
		t.Errorf("want kernelBSSSize=0x8c9a0, got %#x", kernelBSSSize)
	}
}

func TestKernelSizesShortEntry(t *testing.T) {
	z := &ZImage{
		Table: []TableEntry{
			{Tag: TagKernelSize, Data: []uint32{0x875534}},
		},
	}

	if _, _, err := z.GetKernelSizes(); !errors.Is(err, os.ErrInvalid) {
		t.Errorf("GetKernelSizes() on a one-word entry = %v, want %v", err, os.ErrInvalid)
	}
}
