// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64,!race

package integration

import (
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest"
)

// TestIO tests the string "UART TEST" is written to the serial port on 0x3f8.
func RunTestIO(t *testing.T, initramfs string) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	if len(initramfs) == 0 {
		f, err := vmtest.CreateTestInitramfs(
			uroot.Opts{}, "github.com/u-root/u-root/integration/testcmd/io/uinit", "")
		if err != nil {
			t.Errorf("failed to create test initramfs: %v", err)
		}
		defer os.Remove(f)
		initramfs = f
	}

	// Create the CPIO and start QEMU.
	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Initramfs: initramfs,
	})
	defer cleanup()

	if err := q.Expect("UART TEST"); err != nil {
		t.Fatal(`expected "UART TEST", got error: `, err)
	}
}
