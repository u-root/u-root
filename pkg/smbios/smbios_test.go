// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
	"testing"

	"github.com/u-root/u-root/pkg/memio"
)

func TestSMBIOSBaseEFI(t *testing.T) {
	systabPath = "testdata/smbios3_systab"
	base, _, err := SMBIOSBase()
	if err != nil {
		t.Fatal(err)
	}

	var want int64 = 0x12345678

	if base != want {
		t.Errorf("SMBIOSBase(): 0x%x, want 0x%x", base, want)
	}
}

func TestSMBIOSBaseLegacy(t *testing.T) {
	tmpBuf = []byte{0, '_', 'M', 'S', '_', 0, 0, '_', 'S', 'M', '_', 0, 0, 0, 0, 0}
	systabPath = "testdata/systab_NOT_FOUND"
	defer func(old func(base int64, uintn memio.UintN) error) { memioRead = old }(memioRead)
	memioRead = func(base int64, uintn memio.UintN) error {
		dat, ok := uintn.(*memio.ByteSlice)
		if !ok {
			return fmt.Errorf("not supported")
		}
		bufLen := len(tmpBuf)
		for i := int64(0); i < dat.Size(); i++ {
			(*dat)[i] = tmpBuf[(base+i)%int64(bufLen)]
		}
		return nil
	}

	base, _, err := SMBIOSBase()
	if err != nil {
		t.Fatal(err)
	}

	var want int64 = 0xf0007

	if base != want {
		t.Errorf("SMBIOSBase(): 0x%x, want 0x%x", base, want)
	}
}

func TestSMBIOSBaseNotFound(t *testing.T) {
	systabPath = "testdata/systab_NOT_FOUND"
	defer func(old func(base int64, uintn memio.UintN) error) { memioRead = old }(memioRead)
	memioRead = func(base int64, uintn memio.UintN) error {
		return nil
	}

	_, _, err := SMBIOSBase()

	want := "could not find _SM_ or _SM3_ via /dev/mem from 0x000f0000 to 0x00100000"

	if err.Error() != want {
		t.Errorf("SMBIOSBase(): '%v', want %s", err, want)
	}
}
