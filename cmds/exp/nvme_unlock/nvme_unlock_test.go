// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestRunEmptyDisk(t *testing.T) {
	err := run("", false, false, false, false)
	if err == nil {
		t.Errorf("run(\"\") = nil, want error")
	}
}
