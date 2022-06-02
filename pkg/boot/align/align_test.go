// Copyright 2015-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package align

import (
	"testing"
)

func TestAlignUpPageSize(t *testing.T) {
	for _, tt := range []struct {
		name string
		val  uint
		want uint
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
			val:       uint(0x05),
			alignSize: uint(0x05),
			want:      uint(0x05),
		},
		{
			name:      "different alignSize, below",
			val:       uint(0x05),
			alignSize: uint(0x06),
			want:      uint(0x06),
		},
		{
			name:      "different alignSize, next",
			val:       uint(0x07),
			alignSize: uint(0x06),
			want:      uint(0x12),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := AlignUpBySize(tt.val)
			if got != tt.want {
				t.Errorf("AlignUpPageSize(%#02x) = %#02x, want: %#02x", tt.val, got, tt.want)
			}
		})
	}

}
