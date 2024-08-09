// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd || linux

package ldd

import (
	"strings"
	"testing"
)

// TestLdd tests Ldd against /bin/date.
// This is just about guaranteed to have
// some output on most linux systems.
func TestLdd(t *testing.T) {
	n, err := List("/bin/date")
	if err != nil {
		t.Fatalf("Ldd on /bin/date: want nil, got %v", err)
	}
	t.Logf("TestLdd: /bin/date has deps of")
	for i := range n {
		t.Logf("\t%v", n[i])
	}
}

// TestLddList tests that the LddList is the
// same as the info returned by Ldd
func TestLddList(t *testing.T) {
	n, err := List("/bin/date")
	if err != nil {
		t.Fatal(err)
	}

	// Find the first name in the array that contains "lib"
	// Test 'em all
	for _, f := range n {
		if !strings.Contains(f, "lib") {
			continue
		}
		t.Logf("Test %v", f)
		n, err := List(f)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v has deps of %v", f, n)
	}
}
