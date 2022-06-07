// Copyright 2015-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package align

import (
	"testing"
)

func TestAlignUpPageSize(t *testing.T) {
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
			name:      "equal",
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
			pageSize = tt.alignSize
			got := AlignUpBySize(tt.val, tt.alignSize)
			if got != tt.want {
				t.Errorf("AlignUpBySize(%#02x, %#02x) = %#02x, want: %#02x", tt.val, tt.alignSize, got, tt.want)
			}
		})
	}

}
