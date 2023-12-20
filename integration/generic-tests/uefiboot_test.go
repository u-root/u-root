// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64 && !race
// +build amd64,!race

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/hugelgupf/vmtest"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/u-root/u-root/pkg/uroot"
)

// TestUefiboot tests uefiboot commmands to boot to uefishell.
func TestUefiBoot(t *testing.T) {
	vmtest.SkipIfNotArch(t, qemu.ArchAMD64)

	payload := "UEFIPAYLOAD.fd"
	src := fmt.Sprintf("/home/circleci/%v", payload)
	if tk := os.Getenv("UROOT_TEST_UEFIPAYLOAD_DIR"); len(tk) > 0 {
		src = filepath.Join(tk, payload)
	}

	if _, err := os.Stat(src); err != nil && os.IsNotExist(err) {
		t.Skipf("UEFI payload image is not found: %s\n Usage: uefiboot <payload>", src)
	}

	testCmds := []string{
		"uefiboot /dev/sda",
	}
	vm := vmtest.StartVMAndRunCmds(t, testCmds,
		vmtest.WithMergedInitramfs(uroot.Opts{Commands: uroot.BusyBoxCmds(
			"github.com/u-root/u-root/cmds/exp/uefiboot",
		)}),
		vmtest.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			qemu.IDEBlockDevice(src),
			qemu.ArbitraryArgs("-machine", "q35"),
			qemu.ArbitraryArgs("-m", "4096"),
		),
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
