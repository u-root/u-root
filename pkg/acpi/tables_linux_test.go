// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"reflect"
	"testing"

	"github.com/hugelgupf/vmtest/guest"
)

// TestTable tests whether any method for getting a table works.
// If it succeeds, it is called again with the method it returns
// to verify at least that much is the same.
func TestTable(t *testing.T) {
	guest.SkipIfNotInVM(t)

	m, tg, err := GetTable()
	if err != nil {
		t.Skip("no table to get")
	}

	f, ok := Methods[m]
	if !ok {
		t.Fatalf("method type returned from GetTable: got not found in Methods, expect to find it")
	}
	tt, err := f()
	if err != nil {
		t.Fatalf("Getting table via %q: got %v, want nil", m, err)
	}
	if !reflect.DeepEqual(tg, tt) {
		t.Fatalf("Getting table via GetTable and %q: got different(%s, %s), want same", m, tg, tt)
	}
}
