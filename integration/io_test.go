// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64

package integration

import (
	"testing"
)

// TestIO tests the string "UART TEST" is written to the serial port on 0x3f8.
func TestIO(t *testing.T) {
	// Create the CPIO and start QEMU.
	tmpDir, q := testWithQEMU(t, options{
		uinitName: "io",
	})
	defer cleanup(t, tmpDir, q)

	if err := q.Expect("UART TEST"); err != nil {
		t.Fatal(`expected "UART TEST", got error: `, err)
	}
}
