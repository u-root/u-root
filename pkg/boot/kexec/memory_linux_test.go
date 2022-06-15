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

	"github.com/google/go-cmp/cmp"
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

func TestParseMemoryMapFromFDT(t *testing.T) {
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
			km := &Memory{}
			err := km.ParseMemoryMapFromFDT(tc.fdt)
			if err != tc.wantErr {
				t.Errorf("km.ParseMemoryMapFromFDT returned error %v, want error %v", err, tc.wantErr)
			}
			checkMemoryMap(t, km.Phys, tc.wantMap)
		})
	}
}

func TestParseMemoryMap(t *testing.T) {
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

	phys, err := internalParseMemoryMap(root)
	if err != nil {
		t.Fatalf("ParseMemoryMap() error: %v", err)
	}
	if !reflect.DeepEqual(phys, want) {
		t.Errorf("ParseMemoryMap() got %v, want %v", phys, want)
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

func TestAvailableRAM(t *testing.T) {
	old := pageMask
	defer func() {
		pageMask = old
	}()
	// suppose we have 4K pages.
	pageMask = 4095

	var mem Memory
	mem.Phys = MemoryMap{
		TypedRange{Range: Range{Start: 0, Size: 8192}, Type: RangeRAM},
		TypedRange{Range: Range{Start: 8192, Size: 8000}, Type: RangeRAM},
		TypedRange{Range: Range{Start: 20480, Size: 1000}, Type: RangeRAM},
		TypedRange{Range: Range{Start: 24576, Size: 1000}, Type: RangeRAM},
		TypedRange{Range: Range{Start: 28672, Size: 1000}, Type: RangeRAM},
	}

	mem.Segments = []Segment{
		{Phys: Range{Start: 40, Size: 50}},
		{Phys: Range{Start: 8000, Size: 200}},
		{Phys: Range{Start: 18000, Size: 1000}},
		{Phys: Range{Start: 24600, Size: 1000}},
		{Phys: Range{Start: 28000, Size: 10000}},
	}

	want := Ranges{
		Range{Start: 0, Size: 40},
		Range{Start: 4096, Size: 8000 - 4096},
		Range{Start: 12288, Size: 8192 + 8000 - 12288},
		Range{Start: 20480, Size: 1000},
		Range{Start: 24576, Size: 24},
	}

	got := mem.AvailableRAM()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AvailableRAM() got %+v, want %+v", got, want)
	}
}

func TestAlignAndMerge(t *testing.T) {
	for _, tt := range []struct {
		name    string
		in      Segments
		want    Segments
		wantErr bool
	}{
		{
			name: "3 buffers in the same page",
			in: Segments{
				NewSegment([]byte("test"), Range{Start: 0, Size: 5}),
				NewSegment([]byte("foo"), Range{Start: 10, Size: 5}),
				NewSegment([]byte("haha"), Range{Start: 15, Size: 4}),
				NewSegment([]byte("hahahahaha"), Range{Start: 0x1000, Size: 10}),
				NewSegment([]byte("hhhhhhhhhh"), Range{Start: 0x2000, Size: 1}),
			},
			want: Segments{
				NewSegment([]byte("test\x00\x00\x00\x00\x00\x00foo\x00\x00haha"), Range{Start: 0, Size: 0x1000}),
				NewSegment([]byte("hahahahaha"), Range{Start: 0x1000, Size: 0x1000}),
				NewSegment([]byte("h"), Range{Start: 0x2000, Size: 0x1000}),
			},
		},
		{
			name: "no buffer",
			in: Segments{
				NewSegment(nil, Range{Start: 0, Size: 0x1000}),
			},
			want: Segments{
				NewSegment(nil, Range{Start: 0, Size: 0x1000}),
			},
		},
		{
			name: "no buffers",
			in:   Segments{},
			want: Segments{},
		},
		{
			name: "perfectly aligned buffer",
			in: Segments{
				NewSegment([]byte("test"), Range{Start: 0, Size: 0x1000}),
			},
			want: Segments{
				NewSegment([]byte("test"), Range{Start: 0, Size: 0x1000}),
			},
		},
		{
			name: "truncate buffer",
			in: Segments{
				NewSegment([]byte("testtest"), Range{Start: 0, Size: 5}),
			},
			want: Segments{
				NewSegment([]byte("testt"), Range{Start: 0, Size: 0x1000}),
			},
		},
		{
			name: "backfill, truncate buffer",
			in: Segments{
				NewSegment([]byte("testtest"), Range{Start: 2, Size: 5}),
			},
			want: Segments{
				NewSegment([]byte("\x00\x00testt"), Range{Start: 0, Size: 0x1000}),
			},
		},
		{
			name: "physical address overlaps, conflicting content",
			in: Segments{
				NewSegment([]byte("testt"), Range{Start: 0, Size: 5}),
				NewSegment([]byte("aabbc"), Range{Start: 3, Size: 4}),
			},
			wantErr: true,
		},
		{
			// This could potentially be solved some day.
			name: "physical address overlaps, same content",
			in: Segments{
				NewSegment([]byte("testt"), Range{Start: 0, Size: 5}),
				NewSegment([]byte("ttaaa"), Range{Start: 3, Size: 4}),
			},
			wantErr: true,
		},
		{
			name: "two segments in the same page with first-buffer-truncation and second-buffer-padding",
			in: Segments{
				NewSegment([]byte("testtest"), Range{Start: 0, Size: 5}),
				NewSegment([]byte("foo"), Range{Start: 10, Size: 5}),
			},
			want: Segments{
				NewSegment([]byte("testt\x00\x00\x00\x00\x00foo"), Range{Start: 0, Size: 0x1000}),
			},
		},
		{
			name: "two segments in the same page with first-buffer-perfect and second-buffer-truncation",
			in: Segments{
				NewSegment([]byte("testt"), Range{Start: 0, Size: 5}),
				NewSegment([]byte("foofoofoo"), Range{Start: 10, Size: 5}),
			},
			want: Segments{
				NewSegment([]byte("testt\x00\x00\x00\x00\x00foofo"), Range{Start: 0, Size: 0x1000}),
			},
		},
		{
			name: "two segments in the same page with first-buffer-padding and second-buffer-perfect",
			in: Segments{
				NewSegment([]byte("tes"), Range{Start: 0, Size: 5}),
				NewSegment([]byte("foofo"), Range{Start: 10, Size: 5}),
			},
			want: Segments{
				NewSegment([]byte("tes\x00\x00\x00\x00\x00\x00\x00foofo"), Range{Start: 0, Size: 0x1000}),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AlignAndMerge(tt.in)
			gotErr := err != nil
			if gotErr != tt.wantErr {
				t.Errorf("AlignAndMerge = %v, want error %t", err, tt.wantErr)
			} else if err != nil {
				return
			}

			if err := got.IsSupersetOf(tt.in); err != nil {
				t.Errorf("AlignAndMerge = %v: %v", got, err)
			}

			gotRanges := got.Phys()
			wantRanges := tt.want.Phys()
			if diff := cmp.Diff(wantRanges, gotRanges); diff != "" {
				t.Errorf("AlignAndMerge physical ranges = (-want, +got):\n%s", diff)
			}
			for i, s := range got {
				b := s.Buf.toSlice()
				if diff := cmp.Diff(tt.want[i].Buf.toSlice(), b); diff != "" {
					t.Errorf("segment %s bytes differ (-want, +got):\n%s", got[i].Phys, diff)
				}
			}
		})
	}
}

func TestFindSpaceIn(t *testing.T) {
	for i, tt := range []struct {
		name  string
		rs    Ranges
		size  uint
		limit Range
		want  Range
		err   error
	}{
		{
			name: "no space above 0x1000",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1000},
			},
			size:  0x10,
			limit: RangeFromInterval(0x1000, MaxAddr),
			err:   ErrNotEnoughSpace{Size: 0x10},
		},
		{
			name: "no space under 0x1000",
			rs: Ranges{
				Range{Start: 0x1000, Size: 0x10},
			},
			size:  0x10,
			limit: RangeFromInterval(0, 0x1000),
			err:   ErrNotEnoughSpace{Size: 0x10},
		},
		{
			name: "disjunct space above 0x1000",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1000},
				Range{Start: 0x1000, Size: 0x10},
			},
			size:  0x10,
			limit: RangeFromInterval(0x1000, MaxAddr),
			want:  Range{Start: 0x1000, Size: 0x10},
		},
		{
			name: "just enough space under 0x1000",
			rs: Ranges{
				Range{Start: 0xFF, Size: 0xf},
				Range{Start: 0xFF0, Size: 0x10},
				Range{Start: 0x1000, Size: 0x10},
			},
			size:  0x10,
			limit: RangeFromInterval(0, 0x1000),
			want:  Range{Start: 0xFF0, Size: 0x10},
		},
		{
			name: "all spaces abvoe 0x1000 and under 0x2000 are too small",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1000},
				Range{Start: 0x1000, Size: 0xf},
				Range{Start: 0x1010, Size: 0xf},
				Range{Start: 0x1f00, Size: 0xf},
				Range{Start: 0x2000, Size: 0x10},
			},
			size:  0x10,
			limit: RangeFromInterval(0x1000, 0x2000),
			err:   ErrNotEnoughSpace{Size: 0x10},
		},
		{
			name: "space is split across 0x1000, with enough space above",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1010},
			},
			size:  0x10,
			limit: RangeFromInterval(0x1000, MaxAddr),
			want:  Range{Start: 0x1000, Size: 0x10},
		},
		{
			name: "space is split across 0x1000, with enough space under",
			rs: Ranges{
				Range{Start: 0xFF0, Size: 0x20},
			},
			size:  0x10,
			limit: RangeFromInterval(0, 0x1000),
			want:  Range{Start: 0xFF0, Size: 0x10},
		},
		{
			name: "space is split across 0x1000 and 0x2000, but not enough space above or below",
			rs: Ranges{
				Range{Start: 0xFF1, Size: 0xf + 0xf},
				Range{Start: 0x1FF1, Size: 0xf + 0xf},
			},
			size:  0x10,
			limit: RangeFromInterval(0x1000, 0x2000),
			err:   ErrNotEnoughSpace{Size: 0x10},
		},
		{
			name: "space is split across 0x1000, with enough space in the next one",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x100f},
				Range{Start: 0x1010, Size: 0x10},
			},
			size:  0x10,
			limit: RangeFromInterval(0x1000, MaxAddr),
			want:  Range{Start: 0x1010, Size: 0x10},
		},
		{
			name:  "no ranges",
			rs:    Ranges{},
			size:  0x10,
			limit: RangeFromInterval(0, MaxAddr),
			err:   ErrNotEnoughSpace{Size: 0x10},
		},
		{
			name:  "no ranges, zero size",
			rs:    Ranges{},
			size:  0,
			limit: RangeFromInterval(0, MaxAddr),
			err:   ErrNotEnoughSpace{Size: 0},
		},
	} {
		t.Run(fmt.Sprintf("test_%d_%s", i, tt.name), func(t *testing.T) {
			got, err := tt.rs.FindSpaceIn(tt.size, tt.limit)
			if !reflect.DeepEqual(got, tt.want) || err != tt.err {
				t.Errorf("%s.FindSpaceIn(%#x, limit = %s) = (%#x, %v), want (%#x, %v)", tt.rs, tt.size, tt.limit, got, err, tt.want, tt.err)
			}
		})
	}
}

func TestFindSpace(t *testing.T) {
	for i, tt := range []struct {
		name string
		rs   Ranges
		size uint
		want Range
		err  error
	}{
		{
			name: "just enough space under 0x1000",
			rs: Ranges{
				Range{Start: 0xFF, Size: 0xf},
				Range{Start: 0xFF0, Size: 0x10},
				Range{Start: 0x1000, Size: 0x10},
			},
			size: 0x10,
			want: Range{Start: 0xFF0, Size: 0x10},
		},
		{
			name: "no ranges",
			rs:   Ranges{},
			size: 0x10,
			err:  ErrNotEnoughSpace{Size: 0x10},
		},
		{
			name: "no ranges, zero size",
			rs:   Ranges{},
			size: 0,
			err:  ErrNotEnoughSpace{Size: 0},
		},
	} {
		t.Run(fmt.Sprintf("test_%d_%s", i, tt.name), func(t *testing.T) {
			got, err := tt.rs.FindSpace(tt.size)
			if !reflect.DeepEqual(got, tt.want) || err != tt.err {
				t.Errorf("%s.FindSpace(%#x) = (%#x, %v), want (%#x, %v)", tt.rs, tt.size, got, err, tt.want, tt.err)
			}
		})
	}
}

func TestFindSpaceAbove(t *testing.T) {
	for i, tt := range []struct {
		name string
		rs   Ranges
		size uint
		min  uintptr
		want Range
		err  error
	}{
		{
			name: "no space above 0x1000",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1000},
			},
			size: 0x10,
			min:  0x1000,
			err:  ErrNotEnoughSpace{Size: 0x10},
		},
		{
			name: "disjunct space above 0x1000",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1000},
				Range{Start: 0x1000, Size: 0x10},
			},
			size: 0x10,
			min:  0x1000,
			want: Range{Start: 0x1000, Size: 0x10},
		},
		{
			name: "space is split across 0x1000, with enough space above",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1010},
			},
			size: 0x10,
			min:  0x1000,
			want: Range{Start: 0x1000, Size: 0x10},
		},
		{
			name: "space is split across 0x1000, with enough space in the next one",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x100f},
				Range{Start: 0x1010, Size: 0x10},
			},
			size: 0x10,
			min:  0x1000,
			want: Range{Start: 0x1010, Size: 0x10},
		},
		{
			name: "just enough space under 0x1000",
			rs: Ranges{
				Range{Start: 0xFF, Size: 0xf},
				Range{Start: 0xFF0, Size: 0x10},
				Range{Start: 0x1000, Size: 0x10},
			},
			size: 0x10,
			min:  0,
			want: Range{Start: 0xFF0, Size: 0x10},
		},
		{
			name: "no ranges",
			rs:   Ranges{},
			size: 0x10,
			err:  ErrNotEnoughSpace{Size: 0x10},
		},
		{
			name: "no ranges, zero size",
			rs:   Ranges{},
			size: 0,
			err:  ErrNotEnoughSpace{Size: 0},
		},
	} {
		t.Run(fmt.Sprintf("test_%d_%s", i, tt.name), func(t *testing.T) {
			got, err := tt.rs.FindSpaceAbove(tt.size, tt.min)
			if !reflect.DeepEqual(got, tt.want) || err != tt.err {
				t.Errorf("%s.FindSpaceAbove(%#x, min=%#x) = (%#x, %v), want (%#x, %v)", tt.rs, tt.size, tt.min, got, err, tt.want, tt.err)
			}
		})
	}
}

func TestSort(t *testing.T) {
	for _, tt := range []struct {
		in   Ranges
		want Ranges
	}{
		{
			in: Ranges{
				Range{Start: 2, Size: 5},
				Range{Start: 1, Size: 5},
			},
			want: Ranges{
				Range{Start: 1, Size: 5},
				Range{Start: 2, Size: 5},
			},
		},
		{
			in: Ranges{
				Range{Start: 1, Size: 5},
				Range{Start: 1, Size: 6},
			},
			want: Ranges{
				Range{Start: 1, Size: 6},
				Range{Start: 1, Size: 5},
			},
		},
	} {
		var deepCopy Ranges
		deepCopy = append(deepCopy, tt.in...)
		tt.in.Sort()
		if !reflect.DeepEqual(tt.in, tt.want) {
			t.Errorf("%v.Sort() = %v, want\n%v", deepCopy, tt.in, tt.want)
		}
	}
}

func TestIntersection(t *testing.T) {
	for i, tt := range []struct {
		r             Range
		r2            Range
		wantOverlap   bool
		wantIntersect *Range
	}{
		{
			r:             Range{Start: 0, Size: 50},
			r2:            Range{Start: 49, Size: 1},
			wantOverlap:   true,
			wantIntersect: &Range{Start: 49, Size: 1},
		},
		{
			r:             Range{Start: 0, Size: 50},
			r2:            Range{Start: 50, Size: 1},
			wantOverlap:   false,
			wantIntersect: nil,
		},
		{
			r:             Range{Start: 49, Size: 1},
			r2:            Range{Start: 0, Size: 50},
			wantOverlap:   true,
			wantIntersect: &Range{Start: 49, Size: 1},
		},
		{
			r:             Range{Start: 50, Size: 1},
			r2:            Range{Start: 0, Size: 50},
			wantOverlap:   false,
			wantIntersect: nil,
		},
		{
			r:             Range{Start: 0, Size: 50},
			r2:            Range{Start: 10, Size: 1},
			wantOverlap:   true,
			wantIntersect: &Range{Start: 10, Size: 1},
		},
		{
			r:             Range{Start: 10, Size: 1},
			r2:            Range{Start: 0, Size: 50},
			wantOverlap:   true,
			wantIntersect: &Range{Start: 10, Size: 1},
		},
	} {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			if got := tt.r.Overlaps(tt.r2); got != tt.wantOverlap {
				t.Errorf("%s.Overlaps(%s) = %v, want %v", tt.r, tt.r2, got, tt.wantOverlap)
			}
			if got := tt.r.Intersect(tt.r2); !reflect.DeepEqual(got, tt.wantIntersect) {
				t.Errorf("%s.Intersect(%s) = %v, want %v", tt.r, tt.r2, got, tt.wantIntersect)
			}
		})
	}
}

func TestMinusRange(t *testing.T) {
	for i, tt := range []struct {
		r    Range
		r2   Range
		want []Range
	}{
		{
			// r2 contained completely within r.
			r:  Range{Start: 0x100, Size: 0x200},
			r2: Range{Start: 0x150, Size: 0x50},
			want: []Range{
				{Start: 0x100, Size: 0x50},
				{Start: 0x1a0, Size: 0x160},
			},
		},
		{
			// r contained completely within r2.
			r:    Range{Start: 0x100, Size: 0x50},
			r2:   Range{Start: 0x90, Size: 0x100},
			want: nil,
		},
		{
			r:    Range{Start: 0x100, Size: 0x50},
			r2:   Range{Start: 0x100, Size: 0x100},
			want: nil,
		},
		{
			r:    Range{Start: 0x100, Size: 0x50},
			r2:   Range{Start: 0xf0, Size: 0x60},
			want: nil,
		},
		{
			// Overlaps to the right.
			r:  Range{Start: 0x100, Size: 0x100},
			r2: Range{Start: 0x150, Size: 0x100},
			want: []Range{
				{Start: 0x100, Size: 0x50},
			},
		},
		{
			// Overlaps to the left.
			r:  Range{Start: 0x100, Size: 0x100},
			r2: Range{Start: 0x50, Size: 0x100},
			want: []Range{
				{Start: 0x150, Size: 0xb0},
			},
		},
		{
			// Doesn't overlap at all.
			r:  Range{Start: 0x100, Size: 0x100},
			r2: Range{Start: 0x200, Size: 0x100},
			want: []Range{
				{Start: 0x100, Size: 0x100},
			},
		},
	} {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			if got := tt.r.Minus(tt.r2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s minus %s = %v, want %v", tt.r, tt.r2, got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	for i, tt := range []struct {
		r    Range
		p    uintptr
		want bool
	}{
		{
			r:    Range{Start: 0, Size: 50},
			p:    50,
			want: false,
		},
		{
			r:    Range{Start: 0, Size: 50},
			p:    49,
			want: true,
		},
		{
			r:    Range{Start: 50, Size: 50},
			p:    49,
			want: false,
		},
	} {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			if got := tt.r.Contains(tt.p); got != tt.want {
				t.Errorf("%s.Contains(%#x) = %v, want %v", tt.r, tt.p, got, tt.want)
			}
		})
	}
}

func TestAdjacent(t *testing.T) {
	for i, tt := range []struct {
		r1   Range
		r2   Range
		want bool
	}{
		{
			r1:   Range{Start: 0, Size: 50},
			r2:   Range{Start: 50, Size: 50},
			want: true,
		},
		{
			r1:   Range{Start: 0, Size: 40},
			r2:   Range{Start: 41, Size: 50},
			want: false,
		},
		{
			r1:   Range{Start: 10, Size: 40},
			r2:   Range{Start: 0, Size: 10},
			want: true,
		},
		{
			r1:   Range{Start: 10, Size: 39},
			r2:   Range{Start: 40, Size: 50},
			want: false,
		},
	} {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			got1 := tt.r1.Adjacent(tt.r2)
			got2 := tt.r2.Adjacent(tt.r1)
			if got1 != tt.want {
				t.Errorf("%s.Adjacent(%s) = %v, want %v", tt.r1, tt.r2, got1, tt.want)
			}
			if got2 != tt.want {
				t.Errorf("%s.Adjacent(%s) = %v, want %v", tt.r2, tt.r1, got2, tt.want)
			}
		})
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

func TestSegmentsInsert(t *testing.T) {
	for i, tt := range []struct {
		segs Segments
		s    Segment
		want Segments
	}{
		{
			segs: Segments{
				Segment{Phys: Range{Start: 0x2000, Size: 0x20}},
				Segment{Phys: Range{Start: 0x4000, Size: 0x20}},
			},
			s: Segment{Phys: Range{Start: 0x3000, Size: 0x20}},
			want: Segments{
				Segment{Phys: Range{Start: 0x2000, Size: 0x20}},
				Segment{Phys: Range{Start: 0x3000, Size: 0x20}},
				Segment{Phys: Range{Start: 0x4000, Size: 0x20}},
			},
		},
		{
			segs: Segments{},
			s:    Segment{Phys: Range{Start: 0x3000, Size: 0x20}},
			want: Segments{
				Segment{Phys: Range{Start: 0x3000, Size: 0x20}},
			},
		},
	} {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			before := tt.segs
			tt.segs.Insert(tt.s)

			if !reflect.DeepEqual(tt.segs, tt.want) {
				t.Errorf("\n%v.Insert(%v) = \n%v, want \n%v", before, tt.s, tt.segs, tt.want)
			}
		})
	}
}

func TestIsSupersetOf(t *testing.T) {
	for _, tt := range []struct {
		r    Range
		r2   Range
		want bool
	}{
		{
			r:    Range{Start: 0, Size: 0x1000},
			r2:   Range{Start: 1, Size: 0x1000 - 2},
			want: true,
		},
		{
			r:    Range{Start: 0, Size: 0x1000},
			r2:   Range{Start: 1, Size: 0x1000 - 1},
			want: true,
		},
		{
			r:    Range{Start: 0, Size: 0x1000},
			r2:   Range{Start: 1, Size: 0x1000},
			want: false,
		},
		{
			r:    Range{Start: 0, Size: 0x1000},
			r2:   Range{Start: 0, Size: 0},
			want: true,
		},
		{
			// Hmm... this feels wrong. "IsSupersetOf" may be the
			// wrong name, or the implementation should recognize
			// that any 0-size range is inside any other range.
			r:    Range{Start: 0, Size: 0x1000},
			r2:   Range{Start: 0x1001, Size: 0},
			want: false,
		},
	} {
		got := tt.r.IsSupersetOf(tt.r2)
		if got != tt.want {
			t.Errorf("%s.IsSupersetOf(%s) = %t, want %t", tt.r, tt.r2, got, tt.want)
		}
	}
}
