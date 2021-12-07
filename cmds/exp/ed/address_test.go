// Copyright 2015-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"
)

func TestResolveOffset(t *testing.T) {
	for _, tt := range []struct {
		name          string
		input         string
		wantOffset    int
		wantOffsetCmd int
		want          string
	}{
		{
			name:          "len(ms) == 0",
			input:         "",
			wantOffset:    0,
			wantOffsetCmd: 0,
		},
		{
			name:  "case len(ms[2]) > 0 with err in atoi",
			input: " 11111111111111111111",
			want:  "strconv.Atoi: parsing \"11111111111111111111\": value out of range",
		},
		{
			name:          "case len(ms[2]) > 0",
			input:         " 1",
			wantOffset:    1,
			wantOffsetCmd: 2,
		},
		{
			name:  "case len(ms[3]) > 0 with err in atoi",
			input: " -11111111111111111111",
			want:  "strconv.Atoi: parsing \"11111111111111111111\": value out of range",
		},
		{
			name:          "case len(ms[3]) > 0 case -",
			input:         " -1",
			wantOffset:    -1,
			wantOffsetCmd: 3,
		},
		{
			name:          "case len(ms[3]) > 0 case +",
			input:         " +1",
			wantOffset:    1,
			wantOffsetCmd: 3,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			gotOffset, gotOffsetCmd, err := buffer.ResolveOffset(tt.input)
			buffer = &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}, addr: 2}
			if err != nil {
				if err.Error() != tt.want {
					t.Errorf("ResolveAddr() = %q, want: %q", err.Error(), tt.want)
				}
			} else {
				if gotOffset != tt.wantOffset || gotOffsetCmd != tt.wantOffsetCmd {
					t.Errorf("ResolveAddrs() = %d %d, want: %d %d", gotOffset, gotOffsetCmd, tt.wantOffset, tt.wantOffsetCmd)
				}
			}
		})
	}
}

func TestResolveAddr(t *testing.T) {
	for _, tt := range []struct {
		name       string
		input      string
		wantLine   int
		wantOffset int
		want       string
	}{
		{
			name:       "case single symbol $",
			input:      "$",
			wantLine:   3,
			wantOffset: 1,
		},
		{
			name:       "case single symbol .",
			input:      ".",
			wantLine:   2,
			wantOffset: 1,
		},
		{
			name:       "case number",
			input:      "1",
			wantLine:   0,
			wantOffset: 1,
		},
		{
			name:  "case number with err",
			input: "11111111111111111111",
			want:  "strconv.Atoi: parsing \"11111111111111111111\": value out of range",
		},
		{
			name:       "case offset with case +",
			input:      "+",
			wantLine:   3,
			wantOffset: 1,
		},
		{
			name:       "case offset with case -",
			input:      "-",
			wantLine:   1,
			wantOffset: 1,
		},
		{
			name:  "case mark",
			input: "'t",
			want:  "mark was cleared: t",
		},
		{
			name:       "test7",
			input:      "//",
			wantLine:   2,
			wantOffset: 2,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buffer = &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}, addr: 2, marks: map[byte]int{'t': 5}}
			gotLine, gotOffset, err := buffer.ResolveAddr(tt.input)
			if err != nil {
				if err.Error() != tt.want {
					t.Errorf("ResolveAddr() = %q, want: %q", err.Error(), tt.want)
				}
			} else {
				if gotLine != tt.wantLine || gotOffset != tt.wantOffset {
					t.Errorf("ResolveAddrs() = %d %d, want: %d %d", gotLine, gotOffset, tt.wantLine, tt.wantOffset)
				}
			}
		})
	}
}

func TestResolveAddrs(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input string
		want1 int
		want2 int
		want  string
	}{
		{
			name:  "test1",
			input: ",,",
			want1: 2,
			want2: 1,
		},
		{
			name:  "test2",
			input: ";;",
			want1: 2,
			want2: 1,
		},
		{
			name:  "test3",
			input: "% %",
			want1: 3,
			want2: 2,
		},
		{
			name:  "test4",
			input: "abc",
			want1: 1,
			want2: 0,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buffer = &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}}
			got1, got2, err := buffer.ResolveAddrs(tt.input)
			if err != nil {
				if err.Error() != tt.want {
					t.Errorf("ResolveAddrs() = %q, want: %q", err.Error(), tt.want)
				}
			} else {
				if len(got1) != tt.want1 || got2 != tt.want2 {
					t.Errorf("ResolveAddrs() = %d %d, want: %d %d", len(got1), got2, tt.want1, tt.want2)
				}
			}
		})
	}
}

func TestAddrValue(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input []int
		want  int
		err   error
	}{
		{
			name:  "without err",
			input: []int{2},
			want:  2,
		},
		{
			name:  "len(addr) == 0",
			input: []int{},
			err:   ErrINV,
		},
		{
			name:  "OOB",
			input: []int{5},
			err:   ErrOOB,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buffer = &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}}
			got, err := buffer.AddrValue(tt.input)
			if err != nil {
				if err != tt.err {
					t.Errorf("ResolveAddrs() = %q, want: %q", err.Error(), tt.want)
				}
			} else {
				if got != tt.want {
					t.Errorf("ResolveAddrs() = %d, want: %d", got, tt.want)
				}
			}
		})
	}
}

func TestAddrRange(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input []int
		want  [2]int
		err   error
	}{
		{
			name:  "address out of order",
			input: []int{2, 0},
			err:   fmt.Errorf("address out of order"),
		},
		{
			name:  "case 0",
			input: []int{},
			err:   ErrINV,
		},
		{
			name:  "case 1",
			input: []int{0},
			want:  [2]int{0, 0},
		},
		{
			name:  "OOB",
			input: []int{5},
			err:   ErrOOB,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buffer = &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}}
			got, err := buffer.AddrRange(tt.input)
			if err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("ResolveAddrs() = %q, want: %q", err.Error(), tt.err.Error())
				}
			} else {
				if got != tt.want {
					t.Errorf("ResolveAddrs() = %d, want: %d", got, tt.want)
				}
			}
		})
	}
}

func TestAddrRangeOrLine(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input []int
		want  [2]int
		err   error
	}{
		{
			name:  "len(addrs) > 1",
			input: []int{2, 0},
			err:   fmt.Errorf("address out of order"),
		},
		{
			name:  "else",
			input: []int{5},
			err:   ErrOOB,
		},
		{
			name:  "r[1] = r[0]",
			input: []int{0},
			want:  [2]int{0, 0},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buffer = &FileBuffer{buffer: []string{"0", "1", "2", "3"}, file: []int{0, 1, 2, 3}}
			got, err := buffer.AddrRangeOrLine(tt.input)
			if err != nil {
				if err.Error() != tt.err.Error() {
					t.Errorf("ResolveAddrs() = %q, want: %q", err.Error(), tt.err.Error())
				}
			} else {
				if got != tt.want {
					t.Errorf("ResolveAddrs() = %d, want: %d", got, tt.want)
				}
			}
		})
	}
}
