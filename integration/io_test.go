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
	// TODO: support arm
	if TestArch() != "amd64" {
		t.Skipf("test not supported on %s", TestArch())
	}

	// Create the CPIO and start QEMU.
	q, cleanup := QEMUTest(t, &Options{
		Cmds: []string{
			"github.com/u-root/u-root/integration/testcmd/io/uinit",
			"github.com/u-root/u-root/cmds/init",
			"github.com/u-root/u-root/cmds/io",
		},
	})
	defer cleanup()

	if err := q.Expect("UART TEST"); err != nil {
		t.Fatal(`expected "UART TEST", got error: `, err)
	}
}
