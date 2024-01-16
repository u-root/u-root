// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chips_test

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/flash/chips"
)

func TestLookup(t *testing.T) {
	tests := []struct {
		name string
		id   chips.ID
		want *chips.Chip
		err  error
	}{
		{
			name: "Valid ID",
			id:   0xbf2541,
			want: &chips.Chips[0], // Assume [0] is the SST25VF016B chip
			err:  nil,
		},
		{
			name: "Non-existent ID",
			id:   0xdeadbeef,
			want: nil,
			err:  os.ErrNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chips.Lookup(tt.id)

			if !errors.Is(err, tt.err) {
				t.Fatalf("Lookup(%06x) got %v, want %v", tt.id, err, tt.err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Lookup(%d) got %v, want %v", tt.id, got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name string
		chip chips.Chip
		want string
	}{
		{
			name: "SST25VF016B",
			chip: chips.Chips[0], // Assume [0] is the SST25VF016B chip
			want: "Vendor:SST Chip:SST25VF016B ID:bf2541 Size:2097152 4BA:false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.chip.String()
			if got != tt.want {
				t.Errorf("String() got %v, want %v", got, tt.want)
			}
		})
	}
}
