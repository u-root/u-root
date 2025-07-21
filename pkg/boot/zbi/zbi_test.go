// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package zbi

import (
	"fmt"
	"reflect"
	"testing"
)

var testData = map[string]*Image{
	"chain-load-test.zbi": {
		Header: NewContainerHeader(203560),
		BootItems: []BootItem{
			{
				Header:         NewDefaultHeader(ZBITypeKernelX64, 203528),
				PayloadAddress: 0x40,
			},
		},
		Bootable: true,
	},
	"x86-boot-shim-tests.zbi": {
		Header: NewContainerHeader(264008),
		BootItems: []BootItem{
			{
				Header:         NewDefaultHeader(ZBITypeKernelX64, 106520),
				PayloadAddress: 0x40,
			},
			{
				Header:         NewDefaultHeader(ZBITypeKernelX64, 101248),
				PayloadAddress: 0x1a078,
			},
			{
				Header:         NewDefaultHeader(ZBITypeKernelX64, 56144),
				PayloadAddress: 0x32c18,
			},
		},
		Bootable: true,
	},
	"zbi-chain-load-hello-world-test.zbi": {
		Header: NewContainerHeader(228064),
		BootItems: []BootItem{
			{
				Header:         NewDefaultHeader(ZBITypeKernelX64, 203528),
				PayloadAddress: 0x40,
			},
			{
				Header: Header{
					Type:      ZBITypeStorageKernel,
					Length:    24467,
					Extra:     0xdb90,
					Flags:     0x30001,
					Reserved0: 0,
					Reserved1: 0,
					Magic:     ItemMagic,
					CRC32:     0x9d0920d0,
				},
				PayloadAddress: 0x31b68,
			},
		},
		Bootable: true,
	},
}

func TestLoad(t *testing.T) {
	for filename, expected := range testData {
		actual, err := Load(fmt.Sprintf("testdata/%s", filename))
		if err != nil {
			t.Fatalf("Failed to load %s: %s", filename, err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf(`Parse("testdata/%s") = %#v; want %#v`, filename, actual, expected)
		}
	}
}

func NewDefaultHeader(itemType ZBIType, length uint32) Header {
	return Header{
		Type:      itemType,
		Length:    length,
		Extra:     0,
		Flags:     VersionFlag,
		Reserved0: 0,
		Reserved1: 0,
		Magic:     ItemMagic,
		CRC32:     NoCRC32Flag,
	}
}
