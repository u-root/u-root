// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestNotFound(t *testing.T) {
	var err error
	f, err := os.CreateTemp("", "cbmemNotFound")
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
	f, err := os.CreateTemp("", "cbmemAPU2")
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
		os.WriteFile("json", j, 0o666)
	}
	if string(j) != apu2JSON {
		t.Errorf("APU2 JSON: got %s, want %s", j, apu2JSON)
	}
}

func TestAPU2CBMemWrap(t *testing.T) {
	var err error
	f, err := os.CreateTemp("", "cbmemWRAPAPU2")
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
	f, err := os.CreateTemp("", "cbmemWRAPAPU2")
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
	f, err := os.CreateTemp("", "cbmemWRAPAPU2")
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
	memFile, err := os.CreateTemp("", "cbmemAPU2")
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
	f, err := os.CreateTemp("", "cbmemAPU2")
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

func TestCbmem(t *testing.T) {
	for _, tt := range []struct {
		name       string
		mem        string
		version    bool
		verbose    bool
		timestamps bool
		dumpJSON   bool
		list       bool
		console    bool
		want       string
	}{
		{
			name:    "version true",
			version: true,
			want:    "cbmem in Go, including JSON output\n",
		},
		{
			name:    "verbose true",
			verbose: true,
		},
		{
			name:     "dumpJSON true",
			dumpJSON: true,
			want:     "{\n\t\"Memory\": {\n\t\t\"Tag\": 1,\n\t\t\"Size\": 108,\n\t\t\"Maps\": [\n\t\t\t{\n\t\t\t\t\"Start\": 0,\n\t\t\t\t\"Size\": 4096,\n\t\t\t\t\"Mtype\": 16\n\t\t\t},\n\t\t\t{\n\t\t\t\t\"Start\": 4096,\n\t\t\t\t\"Size\": 651264,\n\t\t\t\t\"Mtype\": 1\n\t\t\t},\n\t\t\t{\n\t\t\t\t\"Start\": 786432,\n\t\t\t\t\"Size\": 2012143616,\n\t\t\t\t\"Mtype\": 1\n\t\t\t},\n\t\t\t{\n\t\t\t\t\"Start\": 2012930048,\n\t\t\t\t\"Size\": 335872,\n\t\t\t\t\"Mtype\": 16\n\t\t\t}\n\t\t]\n\t},\n\t\"MemConsole\": {\n\t\t\"Tag\": 23,\n\t\t\"Size\": 131064,\n\t\t\"Address\": 2013130752,\n\t\t\"CSize\": 0,\n\t\t\"Cursor\": 240,\n\t\t\"Data\": \"PCEngines apu2\\r\\ncoreboot build 20170228\\r\\n2032 MB DRAM\\r\\n\\r\\n\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\\u0000\"\n\t},\n\t\"Consoles\": [\n\t\t\"\"\n\t],\n\t\"TimeStampsTable\": {\n\t\t\"Tag\": 0,\n\t\t\"Size\": 0,\n\t\t\"Addr\": 0\n\t},\n\t\"TimeStamps\": null,\n\t\"UART\": [\n\t\t{\n\t\t\t\"Tag\": 15,\n\t\t\t\"Size\": 20,\n\t\t\t\"Type\": 1,\n\t\t\t\"BaseAddr\": 1016,\n\t\t\t\"Baud\": 115200,\n\t\t\t\"RegWidth\": 16\n\t\t}\n\t],\n\t\"MainBoard\": {\n\t\t\"Tag\": 3,\n\t\t\"Size\": 40,\n\t\t\"Vendor\": \"PC Engines\",\n\t\t\"PartNumber\": \"PCEngines apu2\"\n\t},\n\t\"Hwrpb\": {\n\t\t\"Tag\": 0,\n\t\t\"Size\": 0,\n\t\t\"HwrPB\": 0\n\t},\n\t\"CBMemory\": null,\n\t\"BoardID\": {\n\t\t\"Tag\": 37,\n\t\t\"Size\": 16,\n\t\t\"BoardID\": 2012962816\n\t},\n\t\"StringVars\": {\n\t\t\"LB_TAG_BUILD\": \"Tue Feb 28 22:34:13 UTC 2017\",\n\t\t\"LB_TAG_COMPILE_BY\": \"root\",\n\t\t\"LB_TAG_COMPILE_DOMAIN\": \"\",\n\t\t\"LB_TAG_COMPILE_HOST\": \"3aa919ff57dc\",\n\t\t\"LB_TAG_COMPILE_TIME\": \"22:34:13\",\n\t\t\"LB_TAG_EXTRA_VERSION\": \"-4.0.7\",\n\t\t\"LB_TAG_VERSION\": \"8b10004\"\n\t},\n\t\"BootMediaParams\": {\n\t\t\"Tag\": 0,\n\t\t\"Size\": 0,\n\t\t\"FMAPOffset\": 0,\n\t\t\"CBFSOffset\": 0,\n\t\t\"CBFSSize\": 0,\n\t\t\"BootMediaSize\": 0\n\t},\n\t\"VersionTimeStamp\": 38,\n\t\"Unknown\": null,\n\t\"Ignored\": null\n}\n",
		},
		{
			name: "list true",
			list: true,
		},
		{
			name:    "console true",
			console: true,
			want:    "PCEngines apu2\r\ncoreboot build 20170228\r\n2032 MB DRAM\r\n\r\n\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Create(filepath.Join(t.TempDir(), "cbmemAPU2"))
			if err != nil {
				t.Errorf("could not gen file: %v", err)
			}
			defer f.Close()
			if err := genFile(f, t.Logf, apu2); err != nil {
				t.Errorf("could not gen file: %v", err)
			}

			*mem = f.Name()
			version = tt.version
			verbose = tt.verbose
			timestamps = tt.timestamps
			dumpJSON = tt.dumpJSON
			list = tt.list
			console = tt.console

			buf := &bytes.Buffer{}

			got := cbMem(buf)
			if got != nil {
				if got.Error() != tt.want {
					t.Errorf("cbmem() = %q, want: %q", got.Error(), tt.want)
				}
			} else {
				if buf.String() != tt.want {
					t.Errorf("cbmem() = %q, want: %q", buf.String(), tt.want)
				}
			}
		})
	}
}
