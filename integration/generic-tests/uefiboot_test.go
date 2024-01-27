// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64 && !race
// +build amd64,!race

package integration

import (
	"os"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/hugelgupf/vmtest"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/u-root/pkg/uroot"
)

// TestUefiboot tests uefiboot commmands to boot to uefishell.
func TestUEFIBoot(t *testing.T) {
	vmtest.SkipIfNotArch(t, qemu.ArchAMD64)

	var payload string
	if tk := os.Getenv("UROOT_TEST_UEFIPAYLOAD"); len(tk) == 0 {
		t.Skipf("UROOT_TEST_UEFIPAYLOAD not set to payload")
	} else {
		payload = tk
	}

	if _, err := os.Stat(payload); err != nil && os.IsNotExist(err) {
		t.Skipf("UEFI payload image is not found: %s\n Usage: uefiboot <payload>", payload)
	}

	// Separate loading and executing new kernel so GOCOVERDIR data is collected for uefiboot command.
	script := `
		uefiboot -e=false /dev/sda
		sync
		kexec -e
	`
	vm := vmtest.StartVMAndRunCmds(t, script,
		vmtest.WithBusyboxCommands(
			"github.com/u-root/u-root/cmds/core/kexec",
			"github.com/u-root/u-root/cmds/core/sync",
		),
		// Since busybox mode rewrites commands, build uefiboot
		// straight up as a binary to get integration test coverage.
		vmtest.WithMergedInitramfs(uroot.Opts{
			Commands: uroot.BinaryCmds("github.com/u-root/u-root/cmds/exp/uefiboot"),
		}),
		vmtest.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			qemu.IDEBlockDevice(payload),
			qemu.ArbitraryArgs("-machine", "q35"),
			qemu.ArbitraryArgs("-m", "4096"),
		),
		// Build uefiboot (and all other initramfs commands) with coverage enabled.
		vmtest.WithGoBuildOpts(&golang.BuildOpts{
			ExtraArgs: []string{"-cover", "-coverpkg=github.com/u-root/u-root/...", "-covermode=atomic"},
		}),
	)

	// Edk2 debug mode will print PROGRESS CODE. We will want to make sure
	// payload is booting to uefi shell correctly.
	if _, err := vm.Console.Expect(expect.All(
		// Finish booting.
		expect.String("PROGRESS CODE: V02070003"),
		// Last code before booting to UEFI Shell
		expect.String("PROGRESS CODE: V03058001"),
	)); err != nil {
		t.Errorf("VM output did not match expectations: %v", err)
	}

	if err := vm.Kill(); err != nil {
		t.Errorf("Wait: %v", err)
	}
	_ = vm.Wait()

}
