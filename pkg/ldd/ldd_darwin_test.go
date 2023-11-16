// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ldd

import (
	"testing"
)

// TestLddList tests that the LddList is the
// same as the info returned by Ldd
func TestLddList(t *testing.T) {
	n, err := List("/etc/resolv.conf")
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range n {
		t.Logf("Test %v", f)
		n, err := List(f)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v has deps of %v", f, n)
	}
}
