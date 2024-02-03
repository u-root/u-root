// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package launcher

import (
	"testing"
)

func TestIsInitrdSet(t *testing.T) {
	if IsInitrdSet() {
		t.Fatalf("IsInitrdSet() = true, not false")
	}
}

func TestIsValidBootEntry(t *testing.T) {
	// Entries can consist of only letters, numbers, '-', '_', and '.'.
	testStr1 := "valid_entry-12.3"
	testStr2 := "invalid entry 1"
	testStr3 := "invalid@entry#2"

	if !IsValidBootEntry(testStr1) {
		t.Fatalf("IsValidBootEntry('%s') = false, not true", testStr1)
	}

	if IsValidBootEntry(testStr2) {
		t.Fatalf("IsValidBootEntry('%s') = true, not false", testStr2)
	}

	if IsValidBootEntry(testStr3) {
		t.Fatalf("IsValidBootEntry('%s') = true, not false", testStr3)
	}
}
