// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ldd

import (
	"fmt"
	"testing"
)

// lddOne is a helper that runs Ldd on one file. It returns
// the list so we can use elements from the list on another
// test. We do it this way because, unlike /bin/date, there's
// almost nothing else we can assume exists, e.g. /lib/libc.so
// is a different name on almost every *ix* system.
func lddOne(name string) ([]string, error) {
	libMap := make(map[string]bool)
	n, err := Ldd([]string{name})
	if err != nil {
		return nil, fmt.Errorf("Ldd on %v: want nil, got %v", name, err)
	}
	l, err := List([]string{name})
	if err != nil {
		return nil, fmt.Errorf("LddList on %v: want nil, got %v", name, err)
	}
	if len(n) != len(l) {
		return nil, fmt.Errorf("%v: Len of Ldd(%v) and LddList(%v): want same, got different", name, len(n), len(l))
	}
	for i := range n {
		libMap[n[i].FullName] = true
	}
	for i := range n {
		if !libMap[l[i]] {
			return nil, fmt.Errorf("%v: %v was in LddList but not in Ldd", name, l[i])
		}
	}
	return l, nil
}

// TestLddList tests that the LddList is the
// same as the info returned by Ldd
func TestLddList(t *testing.T) {
	n, err := lddOne("/etc/resolv.conf")
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range n {
		t.Logf("Test %v", f)
		n, err := lddOne(f)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v has deps of %v", f, n)
	}
}
