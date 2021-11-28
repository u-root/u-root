// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		in  string
		out time.Duration
		err error
	}{
		{"", time.Duration(0), errDuration},
		{"xyz", time.Duration(0), errDuration},
		{"-2.5", time.Duration(0), errDuration},
		{"-2.5s", time.Duration(0), errDuration},
		{"2.5", time.Duration(2500 * time.Millisecond), nil},
		{"2.5s", time.Duration(2500 * time.Millisecond), nil},
		{"300m", time.Duration(300 * time.Minute), nil},
		{"2h45m", time.Duration(2*time.Hour + 45*time.Minute), nil},
	}

	// Table-driven testing
	for _, tt := range tests {
		out, err := parseDuration(tt.in)
		if out != tt.out || err != tt.err {
			t.Errorf("parseDuration(%#v) = %v, %v; want %v, %v",
				tt.in, out, err, tt.out, tt.err)
		}
	}
}
