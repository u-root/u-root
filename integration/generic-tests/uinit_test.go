// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package integration

import (
	"testing"
	"time"

	"github.com/hugelgupf/vmtest"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/u-root/u-root/pkg/uroot"
)

// TestHelloWorld runs an init which prints the string "HELLO WORLD" and exits.
func TestHelloWorld(t *testing.T) {
	vm := vmtest.StartVM(t,
		vmtest.WithMergedInitramfs(uroot.Opts{
			InitCmd:  "init",
			UinitCmd: "uinit",
			Commands: uroot.BusyBoxCmds(
				"github.com/u-root/u-root/integration/testcmd/helloworld/uinit",
				"github.com/u-root/u-root/cmds/core/init",
			),
			TempDir: t.TempDir(),
		}),
		vmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
		),
		vmtest.CollectKernelCoverage(),
	)

	if _, err := vm.Console.ExpectString("HELLO WORLD"); err != nil {
		t.Error(`expected "HELLO WORLD", got error: `, err)
	}
	if err := vm.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}

// TestHelloWorldNegative runs an init which does not print the string "HELLO WORLD".
func TestHelloWorldNegative(t *testing.T) {
	vm := vmtest.StartVM(t,
		vmtest.WithMergedInitramfs(uroot.Opts{
			InitCmd:  "init",
			UinitCmd: "uinit",
			Commands: uroot.BusyBoxCmds(
				"github.com/u-root/u-root/integration/testcmd/helloworld/uinit",
				"github.com/u-root/u-root/cmds/core/init",
			),
			TempDir: t.TempDir(),
		}),
		vmtest.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
		),
		vmtest.CollectKernelCoverage(),
	)

	if _, err := vm.Console.ExpectString("GOODBYE WORLD"); err == nil {
		t.Error(`expected error, but matched "GOODBYE WORLD"`)
	}
	if err := vm.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}
