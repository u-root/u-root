// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zimage

import (
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
