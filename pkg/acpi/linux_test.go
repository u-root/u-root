// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux

package acpi

import (
	"testing"

	"github.com/hugelgupf/vmtest/guest"
)

// TestLinux just verifies that tables read OK.
// It does not verify content as content varies all
// the time.
func TestLinux(t *testing.T) {
	guest.SkipIfNotInVM(t)

	tab, err := RawTablesFromSys()
	if err != nil {
		t.Fatalf("Got %v, want nil", err)
	}
	for _, tt := range tab {
		t.Logf("table %s", String(tt))
	}
}
