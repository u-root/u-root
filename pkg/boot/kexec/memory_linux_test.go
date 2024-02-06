// Copyright 2018-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kexec

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
				NewSegment([]byte{}, Range{Start: 0, Size: 0x1000}),
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
				if diff := cmp.Diff(tt.want[i].Buf, s.Buf); diff != "" {
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
			err:   ErrNotEnoughSpace,
		},
		{
			name: "no space under 0x1000",
			rs: Ranges{
				Range{Start: 0x1000, Size: 0x10},
			},
			size:  0x10,
			limit: RangeFromInterval(0, 0x1000),
			err:   ErrNotEnoughSpace,
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
			err:   ErrNotEnoughSpace,
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
			err:   ErrNotEnoughSpace,
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
			err:   ErrNotEnoughSpace,
		},
		{
			name:  "no ranges, zero size",
			rs:    Ranges{},
			size:  0,
			limit: RangeFromInterval(0, MaxAddr),
			err:   ErrNotEnoughSpace,
		},
	} {
		t.Run(fmt.Sprintf("test_%d_%s", i, tt.name), func(t *testing.T) {
			got, err := tt.rs.FindSpaceIn(tt.size, tt.limit)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s.FindSpaceIn(%#x, limit = %s) = %#x, want %#x", tt.rs, tt.size, tt.limit, got, tt.want)
			}
			if !errors.Is(err, tt.err) {
				t.Errorf("%s.FindSpaceIn(%#x, limit = %s) = %v, want %v", tt.rs, tt.size, tt.limit, err, tt.err)
			}
		})
	}
}

func TestFindSpace(t *testing.T) {
	for i, tt := range []struct {
		name string
		rs   Ranges
		opts []FindOptioner
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
			err:  ErrNotEnoughSpace,
		},
		{
			name: "no ranges, zero size",
			rs:   Ranges{},
			size: 0,
			err:  ErrNotEnoughSpace,
		},
		{
			name: "no space above 0x1000",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1000},
			},
			size: 0x10,
			opts: []FindOptioner{WithMinimumAddr(0x1000)},
			err:  ErrNotEnoughSpace,
		},
		{
			name: "disjunct space above 0x1000",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1000},
				Range{Start: 0x1000, Size: 0x10},
			},
			size: 0x10,
			opts: []FindOptioner{WithMinimumAddr(0x1000)},
			want: Range{Start: 0x1000, Size: 0x10},
		},
		{
			name: "space is split across 0x1000, with enough space above",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1010},
			},
			size: 0x10,
			opts: []FindOptioner{WithMinimumAddr(0x1000)},
			want: Range{Start: 0x1000, Size: 0x10},
		},
		{
			name: "space is split across 0x1000, with enough space in the next one",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x100f},
				Range{Start: 0x1010, Size: 0x10},
			},
			size: 0x10,
			opts: []FindOptioner{WithMinimumAddr(0x1000)},
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
			want: Range{Start: 0xFF0, Size: 0x10},
		},
		{
			name: "no space under 0x1000",
			rs: Ranges{
				Range{Start: 0x1000, Size: 0x10},
			},
			size: 0x10,
			opts: []FindOptioner{WithinRange(RangeFromInterval(0, 0x1000))},
			err:  ErrNotEnoughSpace,
		},
		{
			name: "disjunct space above 0x1000",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1000},
				Range{Start: 0x1000, Size: 0x10},
			},
			size: 0x10,
			opts: []FindOptioner{WithinRange(RangeFromInterval(0x1000, MaxAddr))},
			want: Range{Start: 0x1000, Size: 0x10},
		},
		{
			name: "just enough space under 0x1000",
			rs: Ranges{
				Range{Start: 0xFF, Size: 0xf},
				Range{Start: 0xFF0, Size: 0x10},
				Range{Start: 0x1000, Size: 0x10},
			},
			size: 0x10,
			opts: []FindOptioner{WithinRange(RangeFromInterval(0, 0x1000))},
			want: Range{Start: 0xFF0, Size: 0x10},
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
			size: 0x10,
			opts: []FindOptioner{WithinRange(RangeFromInterval(0x1000, 0x2000))},
			err:  ErrNotEnoughSpace,
		},
		{
			name: "space is split across 0x1000, with enough space above",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1010},
			},
			size: 0x10,
			opts: []FindOptioner{WithinRange(RangeFromInterval(0x1000, MaxAddr))},
			want: Range{Start: 0x1000, Size: 0x10},
		},
		{
			name: "space is split across 0x1000, with enough space under",
			rs: Ranges{
				Range{Start: 0xFF0, Size: 0x20},
			},
			size: 0x10,
			opts: []FindOptioner{WithinRange(RangeFromInterval(0, 0x1000))},
			want: Range{Start: 0xFF0, Size: 0x10},
		},
		{
			name: "space is split across 0x1000 and 0x2000, but not enough space above or below",
			rs: Ranges{
				Range{Start: 0xFF1, Size: 0xf + 0xf},
				Range{Start: 0x1FF1, Size: 0xf + 0xf},
			},
			size: 0x10,
			opts: []FindOptioner{WithinRange(RangeFromInterval(0x1000, 0x2000))},
			err:  ErrNotEnoughSpace,
		},
		{
			name: "space is split across 0x1000, with enough space in the next one",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x100f},
				Range{Start: 0x1010, Size: 0x10},
			},
			size: 0x10,
			opts: []FindOptioner{WithinRange(RangeFromInterval(0x1000, MaxAddr))},
			want: Range{Start: 0x1010, Size: 0x10},
		},
		{
			name: "alignment with limit",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1000},
				Range{Start: 0x1010, Size: 0x10},
				Range{Start: 0x2000, Size: 0x10},
			},
			size: 0x10,
			opts: []FindOptioner{WithinRange(RangeFromInterval(0x500, MaxAddr)), WithStartAlignment(0x1000)},
			want: Range{Start: 0x2000, Size: 0x10},
		},
		{
			name: "alignment with limit",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1000},
				Range{Start: 0x1010, Size: 0x1010},
				Range{Start: 0x3000, Size: 0x10},
			},
			size: 0x10,
			opts: []FindOptioner{WithinRange(RangeFromInterval(0x500, MaxAddr)), WithStartAlignment(0x1000)},
			want: Range{Start: 0x2000, Size: 0x10},
		},
		{
			name: "alignment with limit",
			rs: Ranges{
				Range{Start: 0x0, Size: 0x1000},
				Range{Start: 0x1010, Size: 0x1010},
				Range{Start: 0x3000, Size: 0x1000},
			},
			size: 0x10,
			opts: []FindOptioner{WithinRange(RangeFromInterval(0x500, MaxAddr)), WithAlignment(0x1000)},
			want: Range{Start: 0x3000, Size: 0x1000},
		},
	} {
		t.Run(fmt.Sprintf("test_%d_%s", i, tt.name), func(t *testing.T) {
			got, err := tt.rs.FindSpace(tt.size, tt.opts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s.FindSpace(%#x) = %#x, want %#x", tt.rs, tt.size, got, tt.want)
			}
			if !errors.Is(err, tt.err) {
				t.Errorf("%s.FindSpace(%#x) = %v, want %v", tt.rs, tt.size, err, tt.err)
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
			err:  ErrNotEnoughSpace,
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
			err:  ErrNotEnoughSpace,
		},
		{
			name: "no ranges, zero size",
			rs:   Ranges{},
			size: 0,
			err:  ErrNotEnoughSpace,
		},
	} {
		t.Run(fmt.Sprintf("test_%d_%s", i, tt.name), func(t *testing.T) {
			got, err := tt.rs.FindSpaceAbove(tt.size, tt.min)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s.FindSpaceAbove(%#x, min=%#x) = %#x, want %#x", tt.rs, tt.size, tt.min, got, tt.want)
			}
			if !errors.Is(err, tt.err) {
				t.Errorf("%s.FindSpaceAbove(%#x, min=%#x) = %v, want %v", tt.rs, tt.size, tt.min, err, tt.err)
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

func TestRanges(t *testing.T) {
	for _, tt := range []struct {
		start uintptr
		end   uintptr
		conv  func(uintptr, uintptr) Range
		want  Range
	}{
		{
			start: 0,
			end:   0x1000,
			conv:  RangeFromInterval,
			want:  Range{Start: 0, Size: 0x1000},
		},
		{
			start: 0,
			end:   0xfff,
			conv:  RangeFromInclusiveInterval,
			want:  Range{Start: 0, Size: 0x1000},
		},
		{
			start: 0,
			end:   0,
			conv:  RangeFromInterval,
			want:  Range{Start: 0, Size: 0},
		},
		{
			start: 0,
			end:   0,
			conv:  RangeFromInclusiveInterval,
			want:  Range{Start: 0, Size: 0x1},
		},
	} {
		got := tt.conv(tt.start, tt.end)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Range(%#x, %#x) = %v, want %v", tt.start, tt.end, got, tt.want)
		}
	}
}

func TestAlign(t *testing.T) {
	for _, tt := range []struct {
		r         Range
		alignSize uint
		want      Range
	}{
		{
			r:         Range{Start: 0x10, Size: 0x10},
			alignSize: 0x1000,
			want:      Range{Start: 0, Size: 0x1000},
		},
		{
			r:         Range{Start: 0x10, Size: 0},
			alignSize: 0x1000,
			want:      Range{Start: 0, Size: 0},
		},
		{
			r:         Range{Start: 0, Size: 0},
			alignSize: 0x1000,
			want:      Range{Start: 0, Size: 0},
		},
		{
			r:         Range{Start: 0, Size: 0x10},
			alignSize: 0x1000,
			want:      Range{Start: 0, Size: 0x1000},
		},
		{
			r:         Range{Start: 0x10, Size: 0x10},
			alignSize: 0,
			want:      Range{Start: 0x10, Size: 0x10},
		},
		{
			r:         Range{Start: 0x10, Size: 0x10},
			alignSize: 1,
			want:      Range{Start: 0x10, Size: 0x10},
		},
	} {
		got := tt.r.Align(tt.alignSize)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%v.Align(%#x) = %v, want %v", tt.r, tt.alignSize, got, tt.want)
		}
	}
}

func TestAlignPage(t *testing.T) {
	for _, tt := range []struct {
		r    Range
		want Range
	}{
		{
			r:    Range{Start: 0x10, Size: 0x10},
			want: Range{Start: 0, Size: 0x1000},
		},
		{
			r:    Range{Start: 0x10, Size: 0},
			want: Range{Start: 0, Size: 0},
		},
		{
			r:    Range{Start: 0, Size: 0},
			want: Range{Start: 0, Size: 0},
		},
		{
			r:    Range{Start: 0, Size: 0x10},
			want: Range{Start: 0, Size: 0x1000},
		},
	} {
		got := tt.r.AlignPage()
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%v.AlignPage() = %v, want %v", tt.r, got, tt.want)
		}
	}
}
