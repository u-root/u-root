// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mtd

import (
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	// If there's no such device then don't bother with the
	// test.
	if _, err := os.Stat(DevName); err != nil {
		t.Skip("No device to test")
	}
	m, err := NewDev(DevName)
	if err != nil {
		t.Fatal(err)
	}
	if err := m.Close(); err != nil {
		t.Fatal(err)
	}
}
