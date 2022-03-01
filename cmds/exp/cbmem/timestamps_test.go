// Copyright 2016-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestReadTimeStamps(t *testing.T) {
	tmpDir := t.TempDir()
	for _, tt := range []struct {
		name string
		addr uint32
		want string
	}{
		{
			name: "addr = 0",
			addr: 0,
			want: "no time stamps",
		},
		{
			name: "addr = 0xffffffff, err",
			addr: 0xffffffff,
			want: "creating TSHeader offsetReader @ 0xffffffff: cbmem tables can only be in 32-bit space and (0xffffffff-0x10000000f is outside it",
		},
		{
			name: "addr = 0x000004ff, no err",
			addr: 0x000004ff,
			want: "26728810564608",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := CBmem{}
			c.TimeStampsTable.Addr = tt.addr

			// Creating file
			file, err := os.Create(filepath.Join(tmpDir, "file"))
			if err != nil {
				t.Errorf("Failed to create file: %v", err)
			}
			defer file.Close()

			if err := genFile(file, t.Logf, apu2); err != nil {
				t.Errorf("could not gen file: %v", err)
			}

			got, err := c.readTimeStamps(file)
			if err != nil {
				if err.Error() != tt.want {
					t.Errorf("readTimeStamps() = %q, want: %q", err.Error(), tt.want)
				}
			} else {
				if fmt.Sprint(got.BaseTime) != tt.want {
					t.Errorf("readTimeStamps() = '%d', want: %q", got.BaseTime, tt.want)
				}
			}
		})
	}
}
