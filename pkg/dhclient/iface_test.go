// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhclient

import "testing"

// TestSimple tests the cases where we should have at least one.
func TestInterfaces(t *testing.T) {
	l, err := Interfaces(".")
	if err != nil {
		t.Fatalf("Checking \".\": got %v, want nil", err)
	}
	if len(l) == 0 {
		t.Fatalf("Checking \".\": got no elements, want at least one")
	}
	// Grab the first one, wrap it in ^$, and we should only get one back.
	one := "^" + l[0].Attrs().Name + "$"
	l, err = Interfaces(one)
	if err != nil {
		t.Fatalf("Checking %s: got %v, want nil", one, err)
	}
	if len(l) != 1 {
		t.Errorf("Matching %s: got %d elements(%v), expect 1", one, len(l), l)
	}
}
