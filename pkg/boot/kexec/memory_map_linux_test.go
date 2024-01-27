// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"fmt"
	"os"
	"path"
	"reflect"
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

func TestMemoryMapFromEFI(t *testing.T) {
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

	phys, err := memoryMapFromEFI(root)
	if err != nil {
		t.Fatalf("MemoryMapFromEFI() error: %v", err)
	}
	if !reflect.DeepEqual(phys, want) {
		t.Errorf("MemoryMapFromEFI() got %v, want %v", phys, want)
	}
}

func TestAsPayloadParam(t *testing.T) {
	var mem Memory
	mem.Phys = MemoryMap{
		TypedRange{Range: Range{Start: 0, Size: 50}, Type: RangeRAM},
		TypedRange{Range: Range{Start: 100, Size: 50}, Type: RangeACPI},
		TypedRange{Range: Range{Start: 200, Size: 50}, Type: RangeNVS},
		TypedRange{Range: Range{Start: 300, Size: 50}, Type: RangeReserved},
		TypedRange{Range: Range{Start: 400, Size: 50}, Type: RangeRAM},
	}
	want := PayloadMemoryMapParam{
		{Start: 0, End: 49, Type: PayloadTypeRAM},
		{Start: 100, End: 149, Type: PayloadTypeACPI},
		{Start: 200, End: 249, Type: PayloadTypeNVS},
		{Start: 300, End: 349, Type: PayloadTypeReserved},
		{Start: 400, End: 449, Type: PayloadTypeRAM},
	}
	mm := mem.Phys.AsPayloadParam()
	if !reflect.DeepEqual(mm, want) {
		t.Errorf("MemoryMap.AsPayloadParam() got %v, want %v", mm, want)
	}
}

func TestMemoryMapInsert(t *testing.T) {
	for i, tt := range []struct {
		m    MemoryMap
		r    TypedRange
		want MemoryMap
	}{
		{
			// r is entirely within m's one range.
			m: MemoryMap{
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
			m: MemoryMap{
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
			m: MemoryMap{
				TypedRange{Range: Range{Start: 0x100, Size: 0x50}, Type: RangeRAM},
			},
			r: TypedRange{Range: Range{Start: 0x100, Size: 0x100}, Type: RangeReserved},
			want: MemoryMap{
				TypedRange{Range: Range{Start: 0x100, Size: 0x100}, Type: RangeReserved},
			},
		},
		{
			// r is the first range in the map.
			m: MemoryMap{},
			r: TypedRange{Range: Range{Start: 0x100, Size: 0x100}, Type: RangeReserved},
			want: MemoryMap{
				TypedRange{Range: Range{Start: 0x100, Size: 0x100}, Type: RangeReserved},
			},
		},
	} {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			// Make a copy for the Errorf print.
			m := tt.m
			tt.m.Insert(tt.r)

			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("\n%v.Insert(%s) =\n%v, want\n%v", m, tt.r, tt.m, tt.want)
			}
		})
	}
}
