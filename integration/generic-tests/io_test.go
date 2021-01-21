// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64,!race

package integration

import (
	"fmt"
	"testing"

	"github.com/u-root/u-root/pkg/vmtest"
)

// TestIO tests the string "UART TEST" is written to the serial port on 0x3f8.
func TestIO(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	testCmds := []string{}
	for _, b := range []byte("UART TEST\r\n") {
		testCmds = append(testCmds, fmt.Sprintf("io outb 0x3f8 %d", b))
	}

	// Create the CPIO and start QEMU.
	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		TestCmds: testCmds,
	})
	defer cleanup()

	if err := q.Expect("UART TEST"); err != nil {
		t.Fatal(`expected "UART TEST", got error: `, err)
	}
}

// TestCMOS runs a series of cmos read and write commands and then checks if the changes to CMOS are reflected.
func TestCMOS(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}
	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Name: "ShellScript",
		TestCmds: []string{
			"io cw 14 1 cr 14 cw 14 0 cr 14",
			"shutdown -h",
		},
	})
	defer cleanup()

	if err := q.Expect("0x01"); err != nil {
		t.Fatal(`expected "0x01", got error: `, err)
	}
	if err := q.Expect("0x00"); err != nil {
		t.Fatal(`expected "0x00", got error: `, err)
	}
}
