// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestParseMemoryMap(t *testing.T) {
	var mem Memory
	root, err := ioutil.TempDir("", "memmap")
	if err != nil {
		t.Fatalf("Cannot create test dir: %v", err)
	}
	defer os.RemoveAll(root)

	old := memoryMapRoot
	memoryMapRoot = root
	defer func() { memoryMapRoot = old }()

	create := func(dir string, start, end uintptr, typ RangeType) error {
		p := path.Join(root, dir)
		if err := os.Mkdir(p, 0755); err != nil {
			return err
		}
		if err := ioutil.WriteFile(path.Join(p, "start"), []byte(fmt.Sprintf("%#x\n", start)), 0655); err != nil {
			return err
		}
		if err := ioutil.WriteFile(path.Join(p, "end"), []byte(fmt.Sprintf("%#x\n", end)), 0655); err != nil {
			return err
		}
		return ioutil.WriteFile(path.Join(p, "type"), append([]byte(typ), '\n'), 0655)
	}

	if err := create("0", 0, 50, RangeRAM); err != nil {
		t.Fatal(err)
	}
	if err := create("1", 100, 150, RangeACPI); err != nil {
		t.Fatal(err)
	}
	if err := create("2", 200, 250, RangeNVS); err != nil {
		t.Fatal(err)
	}
	if err := create("3", 300, 350, RangeReserved); err != nil {
		t.Fatal(err)
	}

	want := []TypedAddressRange{
		{Range: Range{Start: 0, Size: 50}, Type: RangeRAM},
		{Range: Range{Start: 100, Size: 50}, Type: RangeACPI},
		{Range: Range{Start: 200, Size: 50}, Type: RangeNVS},
		{Range: Range{Start: 300, Size: 50}, Type: RangeReserved},
	}

	if err := mem.ParseMemoryMap(); err != nil {
		t.Fatalf("ParseMemoryMap() error: %v", err)
	}
	if !reflect.DeepEqual(mem.Phys, want) {
		t.Errorf("ParseMemoryMap() got %v, want %v", mem.Phys, want)
	}
}

func TestAvailableRAM(t *testing.T) {
	old := pageMask
	defer func() {
		pageMask = old
	}()
	// suppose we have 4K pages.
	pageMask = 4095

	var mem Memory
	mem.Phys = []TypedAddressRange{
		TypedAddressRange{Range: Range{Start: 0, Size: 8192}, Type: RangeRAM},
		TypedAddressRange{Range: Range{Start: 8192, Size: 8000}, Type: RangeRAM},
		TypedAddressRange{Range: Range{Start: 20480, Size: 1000}, Type: RangeRAM},
		TypedAddressRange{Range: Range{Start: 24576, Size: 1000}, Type: RangeRAM},
		TypedAddressRange{Range: Range{Start: 28672, Size: 1000}, Type: RangeRAM},
	}

	mem.Segments = []Segment{
		Segment{Phys: Range{Start: 40, Size: 50}},
		Segment{Phys: Range{Start: 8000, Size: 200}},
		Segment{Phys: Range{Start: 18000, Size: 1000}},
		Segment{Phys: Range{Start: 24600, Size: 1000}},
		Segment{Phys: Range{Start: 28000, Size: 10000}},
	}

	want := []TypedAddressRange{
		TypedAddressRange{Range: Range{Start: 0, Size: 40}, Type: RangeRAM},
		TypedAddressRange{Range: Range{Start: 4096, Size: 8000 - 4096}, Type: RangeRAM},
		TypedAddressRange{Range: Range{Start: 12288, Size: 8192 + 8000 - 12288}, Type: RangeRAM},
		TypedAddressRange{Range: Range{Start: 20480, Size: 1000}, Type: RangeRAM},
		TypedAddressRange{Range: Range{Start: 24576, Size: 24}, Type: RangeRAM},
	}

	got := mem.availableRAM()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("availableRAM() got %+v, want %+v", got, want)
	}
}

func TestAlignPhys(t *testing.T) {
	for _, test := range []struct {
		name      string
		seg, want Segment
	}{
		{
			name: "aligned",
			seg: Segment{
				Buf:  Range{Start: 0x1000, Size: 0x1000},
				Phys: Range{Start: 0x2000, Size: 0x1000},
			},
			want: Segment{
				Buf:  Range{Start: 0x1000, Size: 0x1000},
				Phys: Range{Start: 0x2000, Size: 0x1000},
			},
		},
		{
			name: "unaligned",
			seg: Segment{
				Buf:  Range{Start: 0x1011, Size: 0x1022},
				Phys: Range{Start: 0x2011, Size: 0x1022},
			},
			want: Segment{
				Buf:  Range{Start: 0x1000, Size: 0x1033},
				Phys: Range{Start: 0x2000, Size: 0x2000},
			},
		},
		{
			name: "empty_buf",
			seg: Segment{
				Buf:  Range{Start: 0x1011, Size: 0},
				Phys: Range{Start: 0x2011, Size: 0},
			},
			want: Segment{
				Buf:  Range{Start: 0x1000, Size: 0},
				Phys: Range{Start: 0x2000, Size: 0x1000},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := AlignPhys(test.seg)
			if got != test.want {
				t.Errorf("AlignPhys() got %v, want %v", got, test.want)
			}
		})
	}
}

func TestTryMerge(t *testing.T) {
	for _, test := range []struct {
		name   string
		phys   Range
		merged bool
		want   Range
	}{
		{
			name:   "disjunct",
			phys:   Range{Start: 100, Size: 150},
			merged: false,
		},
		{
			name:   "superset",
			phys:   Range{Start: 0, Size: 80},
			merged: true,
			want:   Range{Start: 0, Size: 100},
		},
		{
			name:   "superset",
			phys:   Range{Start: 10, Size: 80},
			merged: true,
			want:   Range{Start: 0, Size: 100},
		},
		{
			name:   "superset",
			phys:   Range{Start: 10, Size: 90},
			merged: true,
			want:   Range{Start: 0, Size: 100},
		},
		{
			name:   "superset",
			phys:   Range{Start: 0, Size: 100},
			merged: true,
			want:   Range{Start: 0, Size: 100},
		},
		{
			name:   "overlap",
			phys:   Range{Start: 0, Size: 150},
			merged: true,
			want:   Range{Start: 0, Size: 150},
		},
		{
			name:   "overlap",
			phys:   Range{Start: 50, Size: 100},
			merged: true,
			want:   Range{Start: 0, Size: 150},
		},
		{
			name:   "overlap",
			phys:   Range{Start: 99, Size: 51},
			merged: true,
			want:   Range{Start: 0, Size: 150},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			a := NewSegment([]byte("aaaa"), Range{Start: 0, Size: 100})
			b := NewSegment([]byte("bbbb"), test.phys)

			merged := a.tryMerge(b)
			if merged != test.merged {
				t.Fatalf("tryMerge() got %v, want %v", merged, test.merged)
			}
			if !merged {
				return
			}
			if a.Phys != test.want {
				t.Fatalf("Wrong merge result: got %+v, want %+v", a.Phys, test.want)
			}

			got := a.Buf.toSlice()
			want := []byte("aaaabbbb")
			if !bytes.Equal(got, want) {
				t.Errorf("Wrong buf: got %s, want %s", got, want)
			}
		})
	}
}

func TestDedup(t *testing.T) {
	s := []Segment{
		NewSegment([]byte("test"), Range{Start: 0, Size: 100}),
		NewSegment([]byte("test"), Range{Start: 100, Size: 100}),
		NewSegment([]byte("test"), Range{Start: 200, Size: 100}),
		NewSegment([]byte("test"), Range{Start: 250, Size: 50}),
		NewSegment([]byte("test"), Range{Start: 300, Size: 100}),
		NewSegment([]byte("test"), Range{Start: 350, Size: 100}),
	}
	want := []Range{
		Range{Start: 0, Size: 100},
		Range{Start: 100, Size: 100},
		Range{Start: 200, Size: 100},
		Range{Start: 300, Size: 150},
	}

	got := Dedup(s)
	for i := range got {
		if got[i].Phys != want[i] {
			t.Errorf("Dedup() got %v, want %v", got[i].Phys, want[i])
		}
	}

}
