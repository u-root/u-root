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

func TestScript(t *testing.T) {
	testCmds := []string{
		"echo HELLO WORLD",
		"shutdown -h",
	}
	vm := vmtest.StartVMAndRunCmds(t, testCmds,
		vmtest.WithMergedInitramfs(uroot.Opts{Commands: uroot.BusyBoxCmds(
			"github.com/u-root/u-root/cmds/core/shutdown",
		)}),
		vmtest.WithQEMUFn(qemu.WithVMTimeout(30*time.Second)),
	)

	if _, err := vm.Console.ExpectString("HELLO WORLD"); err != nil {
		t.Errorf("Want HELLO WORLD: %v", err)
	}
	if err := vm.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}
