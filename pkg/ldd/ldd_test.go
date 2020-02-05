// Copyright 2017-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build freebsd linux

package ldd

import (
	"fmt"
	"strings"
	"testing"
)

// TestLdd tests Ldd against /bin/date.
// This is just about guaranteed to have
// some output on most linux systems.
func TestLdd(t *testing.T) {
	n, err := Ldd([]string{"/bin/date"})
	if err != nil {
		t.Fatalf("Ldd on /bin/date: want nil, got %v", err)
	}
	t.Logf("TestLdd: /bin/date has deps of")
	for i := range n {
		t.Logf("\t%v", n[i])
	}
}

// lddOne is a helper that runs Ldd on one file. It returns
// the list so we can use elements from the list on another
// test. We do it this way because, unlike /bin/date, there's
// almost nothing else we can assume exists, e.g. /lib/libc.so
// is a different name on almost every *ix* system.
func lddOne(name string) ([]string, error) {
	var libMap = make(map[string]bool)
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
	n, err := lddOne("/bin/date")
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
		n, err := lddOne(f)
		if err != nil {
			t.Error(err)
		}
		t.Logf("%v has deps of %v", f, n)
	}
}
