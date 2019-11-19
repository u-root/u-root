// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64,!race

package integration

import (
	"testing"

	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest"
)

// TestIO tests the string "UART TEST" is written to the serial port on 0x3f8.
func TestIO(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	// Create the CPIO and start QEMU.
	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		BuildOpts: uroot.Opts{
			Commands: uroot.BusyBoxCmds(
				"github.com/u-root/u-root/cmds/core/io",
			),
		},
		Uinit: "github.com/u-root/u-root/integration/testcmd/io/uinit",
	})
	defer cleanup()

	if err := q.Expect("UART TEST"); err != nil {
		t.Fatal(`expected "UART TEST", got error: `, err)
	}
}
