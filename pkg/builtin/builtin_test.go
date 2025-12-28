// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin

import (
	"testing"
)

func TestBuiltIn(t *testing.T) {
	c := Command("hi", "there")

	if c.Path != "hi" {
		t.Errorf("c.Path: got %v, want hi", c.Path)
	}
	if len(c.Args) != 2 {
		t.Fatalf("c.Args: got %d elements, want 2", len(c.Args))
	}
	if c.Args[0] != "hi" {
		t.Errorf("c.Args[0]: got %v, want hi", c.Args[0])
	}
	if c.Args[1] != "there" {
		t.Errorf("c.Args[1]: got %v, want there", c.Args[1])
	}
}
