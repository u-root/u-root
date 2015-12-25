// Copyright 2015 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// By Manoel Vilela <manoel_vilela@engineer.com>

package main

import (
	"testing"
)

// Simple test for ansi command: clear call
func TestAnsiClear(t *testing.T) {
	var test = []string{"clear"}
	if err := ansi(test); err != nil {
		t.Error(err)
	}
}
