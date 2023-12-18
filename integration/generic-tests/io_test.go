// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64 && !race
// +build amd64,!race

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/u-root/u-root/pkg/uroot"
)

// TestIO tests the string "UART TEST" is written to the serial port on 0x3f8.
func TestIO(t *testing.T) {
	vmtest.SkipIfNotArch(t, qemu.ArchAMD64)

	testCmds := []string{}
	for _, b := range []byte("UART TEST\r\n") {
		testCmds = append(testCmds, fmt.Sprintf("io outb 0x3f8 %d", b))
	}

	vm := vmtest.StartVMAndRunCmds(t, testCmds,
		vmtest.WithMergedInitramfs(uroot.Opts{Commands: uroot.BusyBoxCmds(
			"github.com/u-root/u-root/cmds/core/io",
		)}),
		vmtest.WithQEMUFn(qemu.WithVMTimeout(30*time.Second)),
	)

	if _, err := vm.Console.ExpectString("UART TEST"); err != nil {
		t.Error(`expected "UART TEST", got error: `, err)
	}
	if err := vm.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}

// TestCMOS runs a series of cmos read and write commands and then checks if the changes to CMOS are reflected.
func TestCMOS(t *testing.T) {
	vmtest.SkipIfNotArch(t, qemu.ArchAMD64)

	testCmds := []string{
		"io cw 14 1 cr 14 cw 14 0 cr 14",
		"shutdown -h",
	}
	vm := vmtest.StartVMAndRunCmds(t, testCmds,
		vmtest.WithMergedInitramfs(uroot.Opts{Commands: uroot.BusyBoxCmds(
			"github.com/u-root/u-root/cmds/core/io",
			"github.com/u-root/u-root/cmds/core/shutdown",
		)}),
		vmtest.WithQEMUFn(qemu.WithVMTimeout(30*time.Second)),
	)

	if _, err := vm.Console.ExpectString("0x01"); err != nil {
		t.Error(`expected "0x01", got error: `, err)
	}
	if _, err := vm.Console.ExpectString("0x00"); err != nil {
		t.Error(`expected "0x00", got error: `, err)
	}
	if err := vm.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}
