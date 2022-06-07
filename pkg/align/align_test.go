// Copyright 2015-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package align

import (
	"testing"
)

func TestUp(t *testing.T) {
	for _, tt := range []struct {
		name      string
		val       uint
		alignSize uint
		want      uint
	}{
		{
			name:      "below",
			val:       uint(0x02),
			alignSize: uint(0x04),
			want:      uint(0x04),
		},
		{
			name:      "already aligned",
			val:       uint(0x04),
			alignSize: uint(0x04),
			want:      uint(0x04),
		},
		{
			name:      "next",
			val:       uint(0x05),
			alignSize: uint(0x04),
			want:      uint(0x08),
		},
		{
			name:      "different alignSize, already aligned",
			val:       uint(0x08),
			alignSize: uint(0x08),
			want:      uint(0x08),
		},
		{
			name:      "different alignSize, below",
			val:       uint(0x07),
			alignSize: uint(0x08),
			want:      uint(0x08),
		},
		{
			name:      "different alignSize, next",
			val:       uint(0x09),
			alignSize: uint(0x08),
			want:      uint(0x10),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Up(tt.val, tt.alignSize)
			if got != tt.want {
				t.Errorf("Up(%#02x, %#02x) = %#02x, want: %#02x", tt.val, tt.alignSize, got, tt.want)
			}

			pageSize = tt.alignSize
			got = UpPage(tt.val)
			if got != tt.want {
				t.Errorf("UpPage(%#02x) = %#02x, want: %#02x", tt.val, got, tt.want)
			}
		})
	}
}

func TestDown(t *testing.T) {
	for _, tt := range []struct {
		name      string
		val       uint
		alignSize uint
		want      uint
	}{
		{
			name:      "below",
			val:       uint(0x02),
			alignSize: uint(0x04),
			want:      uint(0x00),
		},
		{
			name:      "already aligned",
			val:       uint(0x04),
			alignSize: uint(0x04),
			want:      uint(0x04),
		},
		{
			name:      "next",
			val:       uint(0x05),
			alignSize: uint(0x04),
			want:      uint(0x04),
		},
		{
			name:      "different alignSize, already aligned",
			val:       uint(0x08),
			alignSize: uint(0x08),
			want:      uint(0x08),
		},
		{
			name:      "different alignSize, below",
			val:       uint(0x07),
			alignSize: uint(0x08),
			want:      uint(0x00),
		},
		{
			name:      "different alignSize, next",
			val:       uint(0x09),
			alignSize: uint(0x08),
			want:      uint(0x08),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := Down(tt.val, tt.alignSize)
			if got != tt.want {
				t.Errorf("Down(%#02x, %#02x) = %#02x, want: %#02x", tt.val, tt.alignSize, got, tt.want)
			}

			pageSize = tt.alignSize
			got = DownPage(tt.val)
			if got != tt.want {
				t.Errorf("DownPage(%#02x) = %#02x, want: %#02x", tt.val, got, tt.want)
			}
		})
	}
}
