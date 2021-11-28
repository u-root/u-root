// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package upath

import (
	"testing"
)

func TestUrootPath(t *testing.T) {
	tests := []struct {
		name      string
		urootRoot string
		out       string
	}{
		{"ubin/cat", "", "/ubin/cat"},
		{"ubin/cat", "/", "/ubin/cat"},
		{"ubin/cat", "usr/local", "/usr/local/ubin/cat"},
	}

	for _, tt := range tests {
		root = tt.urootRoot
		o := UrootPath(tt.name)
		if o != tt.out {
			t.Errorf("%v: got %v, want %v", tt, o, tt.out)
		}
	}
}
