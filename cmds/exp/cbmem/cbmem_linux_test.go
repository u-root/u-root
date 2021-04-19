// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func genFile(f *os.File, p func(string, ...interface{}), s []seg) error {
	// Extend the test file to the full 4G, to match hardware.
	if _, err := f.WriteAt([]byte{1}[:], 0xffffffff); err != nil {
		return err
	}
	for _, r := range s {
		p("Write %d bytes at %#x", len(r.dat), r.off)
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
	if err := genFile(memFile, t.Logf, apu2); err != nil {
		t.Fatal(err)
	}
	var c *CBmem
	debug = t.Logf
	var found bool
	for _, addr := range []int64{0, 0xf0000} {
		t.Logf("Check %#08x", addr)
		if c, found, err = parseCBtable(addr, 0x10000); err == nil {
			break
		}
	}
	if !found {
		t.Fatalf("Looking for coreboot table: got false, want true")
	}
	if err != nil {
		t.Fatalf("Reading coreboot table: got %v, want nil", err)
	}
	var b = &bytes.Buffer{}
	DumpMem(memFile, c, false, b)
	t.Logf("%s", b.String())
	apu2Mem := `               Name    Start     Size
       LB_MEM_TABLE 00000000 00001000
         LB_MEM_RAM 00001000 0009f000
         LB_MEM_RAM 000c0000 77eee000
       LB_MEM_TABLE 77fae000 00052000
`
	o := b.String()
	if o != apu2Mem {
		t.Errorf("APU2 DumpMem: got \n%s\n, want \n%s\n", hex.Dump(b.Bytes()), hex.Dump([]byte(apu2Mem)))
	}
	b.Reset()
	DumpMem(memFile, c, true, b)
	if b.Len() == len(apu2Mem) {
		t.Errorf("APU2 DumpMem: got %d bytes output, want more", b.Len())
	}
	// See if JSON even works. TODO: compare output
	j, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}
	ioutil.WriteFile("json", j, 0666)
	if string(j) != apu2JSON {
		t.Errorf("APU2 JSON: got %s, want %s", j, apu2JSON)
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

func TestOffsetReader(t *testing.T) {
	debug = t.Logf
	memFile, err := ioutil.TempFile("", "cbmemAPU2")
	if err != nil {
		t.Fatal(err)
	}
	if err := genFile(memFile, t.Logf, apu2); err != nil {
		t.Fatal(err)
	}
	o, err := newOffsetReader(memFile, 0x77fdf040, 1)
	if err != nil {
		t.Fatalf("newOffsetReader: got %v, want nil", err)
	}
	var b [9]byte
	for _, i := range []int64{0x8000000000, -1, 0x77fdf03f} {
		_, err := o.ReadAt(b[:], i)
		if err == nil {
			t.Errorf("Reading newOffsetReader at %#x: got nil, want err", i)
		}
	}
	for _, i := range []int64{0x77fdf040} {
		n, err := o.ReadAt(b[:], i)
		if err != io.EOF {
			t.Errorf("Reading newOffsetReader at %#x: got %v, want io.EOF", i, err)
		}
		if n != 1 {
			t.Errorf("Reading newOffsetReader at %#x: got %d bytes, want nil", i, n)
		}
	}

	// Now find the LBIO at 0x77fae000
	if o, err = newOffsetReader(memFile, 0x77fae000, 8); err != nil {
		t.Fatalf("newOffsetReader: got %v, want nil", err)
	}
	for _, i := range []int64{0x77fae000} {
		_, err := o.ReadAt(b[:], i)
		if err == nil {
			t.Errorf("Reading newOffsetReader at %#x: got nil, want err", i)
		}
		if string(b[:4]) != "LBIO" {
			t.Errorf("Reading newOffsetReader at %#x: got %q, want LBIO", 0x77fae000, string(b[:4]))
		}
	}
}
