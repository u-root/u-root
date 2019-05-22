// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"testing"
)

func TestFindFileSystem(t *testing.T) {
	var tests = []struct {
		name string
		err  string
	}{
		{"rootfs", "<nil>"},
		{"bogusfs", "bogusfs not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FindFileSystem(tt.name)
			// There has to be a better way to do this.
			if fmt.Sprintf("%v", err) != tt.err {
				t.Errorf("%s: got %v, want %v", tt.name, err, tt.err)
			}
		})
	}
}
