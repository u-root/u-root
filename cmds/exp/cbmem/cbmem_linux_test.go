// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"
)

func genFile(f *os.File, s []seg) error {
	for _, r := range s {
		if _, err := f.WriteAt(r.dat, r.off); err != nil {
			return err
		}
	}
	return nil
}

func TestAPU2(t *testing.T) {
	var err error
	if memFile, err = ioutil.TempFile("", "cbmemAPU2"); err != nil {
		t.Fatal(err)
	}
	if err := genFile(memFile, apu2); err != nil {
		t.Fatal(err)
	}
	var c *CBmem
	debug = t.Logf
	for _, addr := range []int64{0, 0xf0000} {
		t.Logf("Check %#08x", addr)
		if c, err = parseCBtable(addr, 0x10000); err == nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("Reading coreboot table: %v", err)
	}
	var b = &bytes.Buffer{}
	DumpMem(c, b)
	t.Logf("%s", b.String())
	apu2Mem := `               Name    Start     Size
       LB_MEM_TABLE 00000000 00001000
         LB_MEM_RAM 00001000 0009f000
         LB_MEM_RAM 000c0000 77eee000
       LB_MEM_TABLE 77fae000 00052000
`
	if b.String() != apu2Mem {
		t.Errorf("APU2 DumpMem: got \n%s\n, want \n%s\n", hex.Dump(b.Bytes()), hex.Dump([]byte(apu2Mem)))
	}

}

func TestIO(t *testing.T) {
	var (
		b = [8]byte{1}
		i uint64
	)
	if err := readOneSize(bytes.NewReader(b[:1]), &i, 1, 23); err == nil {
		t.Errorf("readOne on too small buffer: got nil, want err")
	}
	if err := readOneSize(bytes.NewReader(b[:]), &i, 1, 23); err == nil {
		t.Fatalf("readOne on too small buffer: got %v, want nil", err)
	}
	if i != i {
		t.Fatalf("readOne value: got %d, want 1", i)
	}
}
