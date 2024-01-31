// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/dt"
)

func checkMemoryMap(t *testing.T, got, want MemoryMap) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("got memory map length %d, want memory map length %d", len(got), len(want))
	}
	for idx, r := range got {
		if r.Type != want[idx].Type {
			t.Errorf("got memory at index %d type %v, want type  %v", idx, r.Type, want[idx].Type)
		}
		if r.Range.Start != want[idx].Start || r.Size != want[idx].Size {
			t.Errorf("got memory at index %d range %v, want range %v", idx, r.Range, want[idx].Range)
		}
	}
}

func TestMemoryMapFromFDT(t *testing.T) {
	for _, tc := range []struct {
		name    string
		fdt     *dt.FDT
		wantMap MemoryMap
		wantErr error
	}{
		{
			"empty",
			&dt.FDT{RootNode: &dt.Node{Name: "/"}},
			MemoryMap{},
			nil,
		},
		{
			"add system memory ok",
			&dt.FDT{
				RootNode: &dt.Node{
					Name: "/",
					Children: []*dt.Node{
						{
							Name: "test memory",
							Properties: []dt.Property{
								{"device_type", append([]byte("memory"), 0)},
								{"reg", []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
							},
						},
						{
							Name: "test memory 2",
							Properties: []dt.Property{
								{"device_type", append([]byte("memory"), 0)},
								{"reg", []byte{0x0, 0x01, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
							},
						},
						{
							Name: "test memory 3",
							Properties: []dt.Property{
								{"device_type", append([]byte("memory"), 0)},
								{"reg", []byte{0x0, 0x03, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x02, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
							},
						},
					},
				},
			},
			MemoryMap{
				TypedRange{Range{Start: uintptr(0x0), Size: 0xffffffffffff}, "System RAM"},
				TypedRange{Range{Start: uintptr(0x1000000000000), Size: 0x1ffffffffffff}, "System RAM"},
				TypedRange{Range{Start: uintptr(0x3000000000000), Size: 0x2ffffffffffff}, "System RAM"},
			},
			nil,
		},
		{
			"add system memory, and reserved memory ok",
			&dt.FDT{
				RootNode: &dt.Node{
					Name: "/",
					Children: []*dt.Node{
						{
							Name: "test memory",
							Properties: []dt.Property{
								{"device_type", append([]byte("memory"), 0)},
								{"reg", []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
							},
						},
						{
							Name: "test memory 2",
							Properties: []dt.Property{
								{"device_type", append([]byte("memory"), 0)},
								{"reg", []byte{0x0, 0x01, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
							},
						},
						{
							Name: "test memory 3",
							Properties: []dt.Property{
								{"device_type", append([]byte("memory"), 0)},
								{"reg", []byte{0x0, 0x03, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x02, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
							},
						},
						{
							Name: "reserved-memory",
							Properties: []dt.Property{
								{"reg", []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff}},
							},
							Children: []*dt.Node{
								{
									Name: "reserved mem child node",
									Properties: []dt.Property{
										{
											"reg", []byte{0x0, 0x03, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff},
										},
									},
								},
							},
						},
					},
				},
			},
			MemoryMap{
				TypedRange{Range{Start: uintptr(0x0), Size: 0xffffffff}, "Reserved"},
				TypedRange{Range{Start: uintptr(0xffffffff), Size: 0xffff00000000}, "System RAM"}, // carve out reserved portion from "reserved-memory".
				TypedRange{Range{Start: uintptr(0x1000000000000), Size: 0x1ffffffffffff}, "System RAM"},
				TypedRange{Range{Start: uintptr(0x3000000000000), Size: 0xffffffffff}, "Reserved"},
				TypedRange{Range{Start: uintptr(0x300ffffffffff), Size: 0x2ff0000000000}, "System RAM"}, // Carve out reserved portion from "reserved mem child node".
			},
			nil,
		},
		{
			"add system memory, reserved memory, and reserved entries ok",
			&dt.FDT{
				ReserveEntries: []dt.ReserveEntry{
					{
						Address: uint64(0x1000000000000),
						Size:    uint64(0xffff),
					},
				},
				RootNode: &dt.Node{
					Name: "/",
					Children: []*dt.Node{
						{
							Name: "test memory",
							Properties: []dt.Property{
								{"device_type", append([]byte("memory"), 0)},
								{"reg", []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
							},
						},
						{
							Name: "test memory 2",
							Properties: []dt.Property{
								{"device_type", append([]byte("memory"), 0)},
								{"reg", []byte{0x0, 0x01, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
							},
						},
						{
							Name: "test memory 3",
							Properties: []dt.Property{
								{"device_type", append([]byte("memory"), 0)},
								{"reg", []byte{0x0, 0x03, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x02, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
							},
						},
						{
							Name: "reserved-memory",
							Properties: []dt.Property{
								{"reg", []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff, 0xff, 0xff}},
							},
							Children: []*dt.Node{
								{
									Name: "reserved mem child node",
									Properties: []dt.Property{
										{
											"reg", []byte{0x0, 0x03, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff},
										},
									},
								},
							},
						},
					},
				},
			},
			MemoryMap{
				TypedRange{Range{Start: uintptr(0x0), Size: 0xffffffff}, "Reserved"},
				TypedRange{Range{Start: uintptr(0xffffffff), Size: 0xffff00000000}, "System RAM"}, // carve out reserved portion from "reserved-memory".
				TypedRange{Range{Start: uintptr(0x1000000000000), Size: 0xffff}, "Reserved"},
				TypedRange{Range{Start: uintptr(0x100000000ffff), Size: 0x1ffffffff0000}, "System RAM"}, // carve out reserve entry.
				TypedRange{Range{Start: uintptr(0x3000000000000), Size: 0xffffffffff}, "Reserved"},
				TypedRange{Range{Start: uintptr(0x300ffffffffff), Size: 0x2ff0000000000}, "System RAM"}, // Carve out reserved portion from "reserved mem child node".
			},
			nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mm, err := MemoryMapFromFDT(tc.fdt)
			if err != tc.wantErr {
				t.Errorf("MemoryMapFromFDT returned error %v, want error %v", err, tc.wantErr)
			}
			checkMemoryMap(t, mm, tc.wantMap)
		})
	}
}

func TestMemoryMapFromSysfsMemmap(t *testing.T) {
	root := t.TempDir()

	create := func(dir string, start, end uintptr, typ RangeType) error {
		p := path.Join(root, dir)
		if err := os.Mkdir(p, 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path.Join(p, "start"), []byte(fmt.Sprintf("%#x\n", start)), 0o655); err != nil {
			return err
		}
		if err := os.WriteFile(path.Join(p, "end"), []byte(fmt.Sprintf("%#x\n", end)), 0o655); err != nil {
			return err
		}
		return os.WriteFile(path.Join(p, "type"), append([]byte(typ), '\n'), 0o655)
	}

	if err := create("0", 0, 49, RangeRAM); err != nil {
		t.Fatal(err)
	}
	if err := create("1", 100, 149, RangeACPI); err != nil {
		t.Fatal(err)
	}
	if err := create("2", 200, 249, RangeNVS); err != nil {
		t.Fatal(err)
	}
	if err := create("3", 300, 349, RangeReserved); err != nil {
		t.Fatal(err)
	}

	want := MemoryMap{
		{Range: Range{Start: 0, Size: 50}, Type: RangeRAM},
		{Range: Range{Start: 100, Size: 50}, Type: RangeACPI},
		{Range: Range{Start: 200, Size: 50}, Type: RangeNVS},
		{Range: Range{Start: 300, Size: 50}, Type: RangeReserved},
	}

	phys, err := memoryMapFromSysfsMemmap(root)
	if err != nil {
		t.Fatalf("MemoryMapFromSysfsMemmap() error: %v", err)
	}
	if !reflect.DeepEqual(phys, want) {
		t.Errorf("MemoryMapFromSysfsMemmap() got %v, want %v", phys, want)
	}
}

func TestToUEFIPayloadMemoryMap(t *testing.T) {
	mm := MemoryMap{
		TypedRange{Range: Range{Start: 0, Size: 50}, Type: RangeRAM},
		TypedRange{Range: Range{Start: 100, Size: 50}, Type: RangeACPI},
		TypedRange{Range: Range{Start: 200, Size: 50}, Type: RangeNVS},
		TypedRange{Range: Range{Start: 300, Size: 50}, Type: RangeReserved},
		TypedRange{Range: Range{Start: 400, Size: 50}, Type: RangeRAM},
	}
	want := UEFIPayloadMemoryMap{
		{Start: 0, End: 49, Type: UEFIPayloadTypeRAM},
		{Start: 100, End: 149, Type: UEFIPayloadTypeACPI},
		{Start: 200, End: 249, Type: UEFIPayloadTypeNVS},
		{Start: 300, End: 349, Type: UEFIPayloadTypeReserved},
		{Start: 400, End: 449, Type: UEFIPayloadTypeRAM},
	}
	uefiMM := mm.ToUEFIPayloadMemoryMap()
	if !reflect.DeepEqual(uefiMM, want) {
		t.Errorf("ToUEFIPayloadMemoryMap() got %v, want %v", uefiMM, want)
	}
}

func TestMemoryMapInsert(t *testing.T) {
	for i, tt := range []struct {
		mm   MemoryMap
		r    TypedRange
		want MemoryMap
	}{
		{
			// r is entirely within m's one range.
			mm: MemoryMap{
				TypedRange{Range: Range{Start: 0, Size: 0x2000}, Type: RangeRAM},
			},
			r: TypedRange{Range: Range{Start: 0x100, Size: 0x100}, Type: RangeReserved},
			want: MemoryMap{
				TypedRange{Range: Range{Start: 0, Size: 0x100}, Type: RangeRAM},
				TypedRange{Range: Range{Start: 0x100, Size: 0x100}, Type: RangeReserved},
				TypedRange{Range: Range{Start: 0x200, Size: 0x2000 - 0x200}, Type: RangeRAM},
			},
		},
		{
			// r sits across three RAM ranges.
			mm: MemoryMap{
				TypedRange{Range: Range{Start: 0, Size: 0x150}, Type: RangeRAM},
				TypedRange{Range: Range{Start: 0x150, Size: 0x50}, Type: RangeRAM},
				TypedRange{Range: Range{Start: 0x1a0, Size: 0x100}, Type: RangeRAM},
			},
			r: TypedRange{Range: Range{Start: 0x100, Size: 0x100}, Type: RangeReserved},
			want: MemoryMap{
				TypedRange{Range: Range{Start: 0, Size: 0x100}, Type: RangeRAM},
				TypedRange{Range: Range{Start: 0x100, Size: 0x100}, Type: RangeReserved},
				TypedRange{Range: Range{Start: 0x200, Size: 0xa0}, Type: RangeRAM},
			},
		},
		{
			// r is a superset of the ranges in m.
			mm: MemoryMap{
				TypedRange{Range: Range{Start: 0x100, Size: 0x50}, Type: RangeRAM},
			},
			r: TypedRange{Range: Range{Start: 0x100, Size: 0x100}, Type: RangeReserved},
			want: MemoryMap{
				TypedRange{Range: Range{Start: 0x100, Size: 0x100}, Type: RangeReserved},
			},
		},
		{
			// r is the first range in the map.
			mm: MemoryMap{},
			r:  TypedRange{Range: Range{Start: 0x100, Size: 0x100}, Type: RangeReserved},
			want: MemoryMap{
				TypedRange{Range: Range{Start: 0x100, Size: 0x100}, Type: RangeReserved},
			},
		},
	} {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			// Make a copy for the Errorf print.
			mm := tt.mm
			tt.mm.Insert(tt.r)

			if !reflect.DeepEqual(tt.mm, tt.want) {
				t.Errorf("\n%v.Insert(%s) =\n%v, want\n%v", mm, tt.r, tt.mm, tt.want)
			}
		})
	}
}

func TestMemoryMapFromIOMem(t *testing.T) {
	f := `10000000-101fffff : reserved
10201000-10202fff : reserved
14000000-1effffff : System RAM
  14154000-14154fff : reserved
  141c0000-14bcffff : reserved
  14c10000-1636ffff : Kernel code
  16370000-1686ffff : reserved
  16870000-1734ffff : Kernel data
  17350000-17377fff : reserved`
	mm, err := memoryMapFromIOMem(strings.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}

	want := MemoryMap{
		TypedRange{Range: RangeFromInterval(0x10000000, 0x101fffff+1), Type: RangeReserved},
		TypedRange{Range: RangeFromInterval(0x10201000, 0x10202fff+1), Type: RangeReserved},
		TypedRange{Range: RangeFromInterval(0x14000000, 0x14154000), Type: RangeRAM},
		TypedRange{Range: RangeFromInterval(0x14154000, 0x14154fff+1), Type: RangeReserved},
		TypedRange{Range: RangeFromInterval(0x14155000, 0x141c0000), Type: RangeRAM},
		TypedRange{Range: RangeFromInterval(0x141c0000, 0x14bcffff+1), Type: RangeReserved},
		TypedRange{Range: RangeFromInterval(0x14bd0000, 0x14c10000), Type: RangeRAM},
		TypedRange{Range: RangeFromInterval(0x14c10000, 0x1636ffff+1), Type: RangeType("Kernel code")},
		TypedRange{Range: RangeFromInterval(0x16370000, 0x1686ffff+1), Type: RangeReserved},
		TypedRange{Range: RangeFromInterval(0x16870000, 0x1734ffff+1), Type: RangeType("Kernel data")},
		TypedRange{Range: RangeFromInterval(0x17350000, 0x17377fff+1), Type: RangeReserved},
		TypedRange{Range: RangeFromInterval(0x17378000, 0x1effffff+1), Type: RangeRAM},
	}
	if !reflect.DeepEqual(mm, want) {
		t.Errorf("Not equal, got %v", mm)
	}

	ignored := `00000000-00000000 : reserved
10000000-101fffff
10201000 : reserved
: System RAM
  141GGGGG-14154fff : reserved
  141c0000-14GGGGGG : reserved`
	mm2, err := memoryMapFromIOMem(strings.NewReader(ignored))
	if err != nil {
		t.Fatal(err)
	}

	if want := MemoryMap(nil); !reflect.DeepEqual(mm2, want) {
		t.Errorf("Memory maps not equal, got %v, want %v", mm2, want)
	}
}

func TestMemoryMapFromMemblock(t *testing.T) {
	memory := `  0: 0x0000004000000000..0x00000040113fffff
   1: 0x0000004011400000..0x00000040123fffff
   2: 0x0000004012400000..0x00000040dfffffff
   3: 0x0000004400000000..0x00000044dfffffff`
	reserved := `  0: 0x0000004000000000..0x00000040113fffff
   1: 0x0000004012400000..0x00000040125fffff
   2: 0x0000004012800000..0x00000040137fffff`
	mm, err := memoryMapFromMemblock(strings.NewReader(memory), strings.NewReader(reserved))
	if err != nil {
		t.Fatal(err)
	}

	want := MemoryMap{
		TypedRange{Range: RangeFromInterval(0x4000000000, 0x4011400000), Type: RangeReserved},
		TypedRange{Range: RangeFromInterval(0x4011400000, 0x4012400000), Type: RangeRAM},
		TypedRange{Range: RangeFromInterval(0x4012400000, 0x4012600000), Type: RangeReserved},
		TypedRange{Range: RangeFromInterval(0x4012600000, 0x4012800000), Type: RangeRAM},
		TypedRange{Range: RangeFromInterval(0x4012800000, 0x4013800000), Type: RangeReserved},
		TypedRange{Range: RangeFromInterval(0x4013800000, 0x40e0000000), Type: RangeRAM},
		TypedRange{Range: RangeFromInterval(0x4400000000, 0x44e0000000), Type: RangeRAM},
	}
	if !reflect.DeepEqual(mm, want) {
		t.Errorf("Not equal, got %v", mm)
	}

	memIgnored := `  0: 0x0000000000000000..0x0000000000000000
   0: 0x0000004000000000..
   1: 0x0000004011400000..0x00000040GGGGGGGG
   0x0000004012400000..0x00000040dfffffff
   2: 0x000000401GGGGGGG..0x00000040dfffffff
   3: 0x00000044000000000x00000044dfffffff`
	reservedIgnored := `  0: 0x0000004000000000..
   1: 0x0000004011400000..0x00000040GGGGGGGG
   0x0000004012400000..0x00000040dfffffff
   2: 0x000000401GGGGGGG..0x00000040dfffffff
   3: 0x00000044000000000x00000044dfffffff`
	mm2, err := memoryMapFromMemblock(strings.NewReader(memIgnored), strings.NewReader(reservedIgnored))
	if err != nil {
		t.Fatal(err)
	}

	if want := MemoryMap(nil); !reflect.DeepEqual(mm2, want) {
		t.Errorf("Memory maps not equal, got %v, want %v", mm2, want)
	}
}

func TestMemoryMapMerge(t *testing.T) {
	mm := MemoryMap{
		TypedRange{Range: Range{Start: 0, Size: 50}, Type: RangeRAM},
		TypedRange{Range: Range{Start: 50, Size: 20}, Type: RangeRAM},
		TypedRange{Range: Range{Start: 70, Size: 40}, Type: RangeRAM},
		TypedRange{Range: Range{Start: 111, Size: 50}, Type: RangeRAM},
		TypedRange{Range: Range{Start: 121, Size: 50}, Type: RangeRAM},
		TypedRange{Range: Range{Start: 400, Size: 50}, Type: RangeReserved},
		TypedRange{Range: Range{Start: 500, Size: 50}, Type: RangeReserved},
		TypedRange{Range: Range{Start: 500, Size: 20}, Type: RangeReserved},
		TypedRange{Range: Range{Start: 600, Size: 20}, Type: RangeReserved},
		TypedRange{Range: Range{Start: 600, Size: 50}, Type: RangeReserved},
	}

	want := MemoryMap{
		TypedRange{Range: Range{Start: 0, Size: 110}, Type: RangeRAM},
		TypedRange{Range: Range{Start: 111, Size: 60}, Type: RangeRAM},
		TypedRange{Range: Range{Start: 400, Size: 50}, Type: RangeReserved},
		TypedRange{Range: Range{Start: 500, Size: 50}, Type: RangeReserved},
		TypedRange{Range: Range{Start: 600, Size: 50}, Type: RangeReserved},
	}

	mm.mergeAdjacent()
	if !reflect.DeepEqual(mm, want) {
		t.Errorf("Merge() got %v, want %v", mm, want)
	}
}
