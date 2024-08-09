// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package acpi

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

// TestTabWrite verifies that we can write a table and
// read it back the same. For fun, we use the DSDT, which is
// the important one. We don't keep tables here as data since we don't know
// which ones we can copy, so we use what's it in sys or skip the test.
func TestTabWrite(t *testing.T) {
	tabs, err := RawTablesFromSys()
	if err != nil || len(tabs) == 0 {
		t.Logf("Skipping test, no ACPI tables in /sys")
		return
	}

	var ttab Table
	for _, tab := range tabs {
		if tab.Sig() == "DSDT" {
			ttab = tab
			break
		}
	}
	if ttab == nil {
		ttab = tabs[0]
	}
	b := &bytes.Buffer{}
	if err := WriteTables(b, ttab); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(b.Bytes(), ttab.Data()) {
		t.Fatalf("Written table does not match")
	}
}

// TestAddr tests the decoding of those stupid "use this 64-bit if set else 32 things"
// that permeate ACPI
func TestAddr(t *testing.T) {
	tests := []struct {
		n   string
		dat []byte
		a64 int64
		a32 int64
		val int64
		err error
	}{
		{n: "zero length data", dat: []byte{}, a64: 5, a32: 1, val: -1, err: fmt.Errorf("no 64-bit address at 5, no 32-bit address at 1, in 0-byte slice")},
		{n: "32 bits at 1, no 64-bit", dat: []byte{1, 2, 3, 4, 5}, a64: 5, a32: 1, val: 84148994, err: nil},
		{n: "64 bits at 5", dat: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}, a64: 5, a32: 1, val: 940138559942690566, err: nil},
	}
	Debug = t.Logf
	for _, tt := range tests {
		v, err := getaddr(tt.dat, tt.a64, tt.a32)
		if v != tt.val {
			t.Errorf("Test %s: got %d, want %d", tt.n, v, tt.val)
		}
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("Test %s: got %v, want %v", tt.n, err, tt.err)
		}
	}
}
