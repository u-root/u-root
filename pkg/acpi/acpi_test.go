// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build linux

package acpi

import (
	"bytes"
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
	var b = &bytes.Buffer{}
	if err := WriteTables(b, ttab); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(b.Bytes(), ttab.Data()) {
		t.Fatalf("Written table does not match")
	}

}
