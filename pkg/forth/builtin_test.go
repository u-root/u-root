// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package forth

import (
	"bytes"
	"testing"
)

func TestBuiltIn(t *testing.T) {
	c := Command("forth", "2", "2", "+")
	Debug = t.Logf
	if err := c.Run(); err != nil {
		t.Fatalf("2 2 +: got %v, want nil", err)
	}
	s := c.Stdout.(*bytes.Buffer).String()
	t.Logf("output %v", s)
	if s != "4/1" {
		t.Errorf("eval: got %s, want 4/1", s)
	}
}
