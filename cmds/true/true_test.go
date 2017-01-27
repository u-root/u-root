// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os/exec"
	"testing"
)

// Ensure 0 is returned.
func TestTrue(t *testing.T) {
	out, err := exec.Command("go", "run", "true.go").CombinedOutput()
	if err != nil || len(out) != 0 {
		t.Fatal("Expected no output and no error")
	}
}
