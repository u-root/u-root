// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vmtest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/testtmp"
	"github.com/u-root/u-root/pkg/uroot"
)

// RunCmdsInVM starts a VM and runs each command provided in testCmds in a
// shell in the VM. If any command fails, the test fails.
//
// The VM can be configured with o. The kernel can be provided via o or
// VMTEST_KERNEL env var. Guest architecture can be set with VMTEST_ARCH.
//
// Underneath, this generates an Elvish script with these commands. The script
// is shared with the VM and run from a special init.
//
//   - TODO: timeouts for individual individual commands.
//   - TODO: It should check their exit status. Hahaha.
func RunCmdsInVM(t *testing.T, testCmds []string, o ...Opt) {
	vm := StartVMAndRunCmds(t, testCmds, o...)

	if _, err := vm.Console.ExpectString("TESTS PASSED MARKER"); err != nil {
		t.Errorf("Waiting for 'TESTS PASSED MARKER' signal: %v", err)
	}

	if err := vm.Wait(); err != nil {
		t.Errorf("VM exited with %v", err)
	}
}

// StartVMAndRunCmds starts a VM and runs each command provided in testCmds in
// a shell in the VM. If the commands return, the VM will be shutdown.
//
// The VM can be configured with o.
//
// Underneath, this generates an Elvish script with these commands. The script
// is shared with the VM and run from a special init.
func StartVMAndRunCmds(t *testing.T, testCmds []string, o ...Opt) *qemu.VM {
	SkipWithoutQEMU(t)

	sharedDir := testtmp.TempDir(t)

	// Generate Elvish shell script of test commands in o.SharedDir.
	if len(testCmds) > 0 {
		testFile := filepath.Join(sharedDir, "test.elv")
		if err := os.WriteFile(testFile, []byte(strings.Join(testCmds, "\n")), 0o777); err != nil {
			t.Fatal(err)
		}
	}

	initramfs := uroot.Opts{
		Commands: uroot.BusyBoxCmds(
			"github.com/u-root/u-root/cmds/core/init",
			"github.com/u-root/u-root/cmds/core/elvish",
			"github.com/hugelgupf/vmtest/vminit/shelluinit",
		),
		InitCmd:  "init",
		UinitCmd: "shelluinit",
		TempDir:  testtmp.TempDir(t),
	}
	return StartVM(t, append([]Opt{
		WithQEMUFn(qemu.P9Directory(sharedDir, "shelltest")),
		WithMergedInitramfs(initramfs),
		CollectKernelCoverage(),
	}, o...)...)
}
