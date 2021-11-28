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

func TestNotFound(t *testing.T) {
	var err error
	f, err := ioutil.TempFile("", "cbmemNotFound")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(make([]byte, 0x100000)); err != nil {
		t.Fatalf("Writing empty file: got %v, want nil", err)
	}
	var found bool
	for _, addr := range []int64{0, 0xf0000} {
		t.Logf("Check %#08x", addr)
		if _, found, err = parseCBtable(f, addr, 0x10000); err == nil {
			break
		}
	}
	if err != nil {
		t.Errorf("Scanning empty file: got %v, want nil", err)
	}
	if found {
		t.Fatalf("Found a coreboot table in empty file: got nil, want err")
	}
}

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
	f, err := ioutil.TempFile("", "cbmemAPU2")
	if err != nil {
		t.Fatal(err)
	}
	if err := genFile(f, t.Logf, apu2); err != nil {
		t.Fatal(err)
	}
	var c *CBmem
	var found bool
	for _, addr := range []int64{0, 0xf0000} {
		t.Logf("Check %#08x", addr)
		if c, found, err = parseCBtable(f, addr, 0x10000); err == nil {
			break
		}
	}
	if !found {
		t.Fatalf("Looking for coreboot table: got false, want true")
	}
	if err != nil {
		t.Fatalf("Reading coreboot table: got %v, want nil", err)
	}
	b := &bytes.Buffer{}
	DumpMem(f, c, false, b)
	t.Logf("%s", b.String())
	o := b.String()
	if o != apu2Mem {
		t.Errorf("APU2 DumpMem: got \n%s\n, want \n%s\n", hex.Dump(b.Bytes()), hex.Dump([]byte(apu2Mem)))
	}
	b.Reset()
	DumpMem(f, c, true, b)
	t.Logf("2nd dump string is %s", b.String())
	if b.Len() == len(apu2Mem) {
		t.Errorf("APU2 DumpMem: got %d bytes output, want more", b.Len())
	}
	// See if JSON even works. TODO: compare output
	j, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}
	// You can use this to generate new test data. It's a timesaver.
	if false {
		ioutil.WriteFile("json", j, 0o666)
	}
	if string(j) != apu2JSON {
		t.Errorf("APU2 JSON: got %s, want %s", j, apu2JSON)
	}
}

func TestAPU2CBMemWrap(t *testing.T) {
	var err error
	f, err := ioutil.TempFile("", "cbmemWRAPAPU2")
	if err != nil {
		t.Fatal(err)
	}
	// Need to patch this a bit. First, add a patch so that the cursor has wrapped.
	p := apu2
	if true {
		p = append(p, []seg{
			{
				// The buffer size will be 4, and the cursor will be 2 -> wrap.
				off: 0x77fdf000, dat: []byte{
					0x04 /*'ø'*/, 0x00 /*'ÿ'*/, 0x00 /*'\x01'*/, 0x00 /*'\x00'*/, 0x02 /*'\x02' */, 0x00 /*'\x00'*/, 0x00 /*'\x00'*/, 0x80, /*'\x80'*/
				},
			},
		}...)
	}
	if err := genFile(f, t.Logf, p); err != nil {
		t.Fatal(err)
	}
	var c *CBmem
	var found bool
	if c, found, err = parseCBtable(f, 0, 0x10000); err != nil {
		t.Fatalf("reading CB table: got %v, want nil", err)
	}
	if !found {
		t.Fatalf("Looking for coreboot table: got false, want true")
	}
	want := "EnPC"
	got := c.MemConsole.Data
	if got != want {
		t.Fatalf("Console data: got %q, want %q", got, want)
	}
}

func TestAPU2CBBadCursor(t *testing.T) {
	var err error
	f, err := ioutil.TempFile("", "cbmemWRAPAPU2")
	if err != nil {
		t.Fatal(err)
	}
	// Need to patch this a bit. First, add a patch so that the cursor has wrapped.
	p := apu2
	p = append(p, []seg{
		{
			// The buffer size will be 4, and the cursor will be 2 -> wrap.
			off: 0x77fdf000, dat: []byte{
				/*0x77fdf000*/ 0x04 /*'ø'*/, 0x00 /*'ÿ'*/, 0x00 /*'\x01'*/, 0x00 /*'\x00'*/, 0x0a /*'\x0a'*/, 0x00 /*'\x00'*/, 0x00 /*'\x00'*/, 0x80, /*'\x80'*/
			},
		},
	}...)

	if err := genFile(f, t.Logf, p); err != nil {
		t.Fatal(err)
	}
	var c *CBmem
	var found bool
	if c, found, err = parseCBtable(f, 0, 0x10000); err != nil {
		t.Fatalf("reading CB table: got %v, want nil", err)
	}
	if !found {
		t.Fatalf("Looking for coreboot table: got false, want true")
	}
	want := "PCEn"
	got := c.MemConsole.Data
	if got != want {
		t.Fatalf("Console data: got %q, want %q", got, want)
	}
}

func TestAPU2CBBadPtr(t *testing.T) {
	var err error
	f, err := ioutil.TempFile("", "cbmemWRAPAPU2")
	if err != nil {
		t.Fatal(err)
	}
	// Need to patch this a bit. First, add a patch so that the cursor has wrapped.
	p := apu2
	p = append(p, []seg{
		{
			off: 0x77fae170, dat: []byte{
				0xff, 0xff, 0xff, 0xff,
			},
		},
	}...)
	if err := genFile(f, t.Logf, p); err != nil {
		t.Fatal(err)
	}
	if _, _, err = parseCBtable(f, 0, 0x10000); err == nil {
		t.Fatalf("reading CB table: got nil, want err")
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
	if err := readOneSize(bytes.NewReader(b[:]), &i, 0, 8); err != nil {
		t.Fatalf("readOne: got %v, want nil", err)
	}
	if i != 1 {
		t.Fatalf("readOne value: got %d, want 1", i)
	}
}

func TestOffsetReader(t *testing.T) {
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
	if err := memFile.Close(); err != nil {
		t.Fatalf("Closing %s: got %v, want nil", memFile.Name(), err)
	}
	if _, err := newOffsetReader(memFile, 0x77fdf040, 1); err == nil {
		t.Fatalf("newOffsetReader: got nil, want err")
	}
}

func TestTimeStampsAPU2(t *testing.T) {
	f, err := ioutil.TempFile("", "cbmemAPU2")
	if err != nil {
		t.Fatal(err)
	}
	if err := genFile(f, t.Logf, apu2); err != nil {
		t.Fatal(err)
	}
	var c *CBmem
	var found bool
	for _, addr := range []int64{0, 0xf0000} {
		t.Logf("Check %#08x", addr)
		if c, found, err = parseCBtable(f, addr, 0x10000); err == nil {
			break
		}
	}
	if !found {
		t.Fatalf("Looking for coreboot table: got false, want true")
	}
	if err != nil {
		t.Fatalf("Reading coreboot table: got %v, want nil", err)
	}
	if c.TimeStampsTable.Addr != 0 {
		t.Fatalf("TimeStampsTable: got %#x, want 0", c.TimeStampsTable)
	}
}
