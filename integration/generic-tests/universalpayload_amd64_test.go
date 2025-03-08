// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build amd64 && !race

package integration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/mkuimage/uimage"
)

// TestUPLBootAmd64 tests '/bbin/kexec UplFitX64.fit' to boot to UEFI Shell.
func TestUPLBootAmd64(t *testing.T) {
	//os.Setenv("VMTEST_ARCH", "amd64")
	qemu.SkipIfNotArch(t, qemu.ArchAMD64)

	// Check required images, including:
	//   OVMF.fd 		-- UEFI OVMF binary
	//   bzImage 		-- Linuxboot kernel image
	//   UplFitX64.fit	-- UPL Amd64 image with FDT enabled

	var ovmf string
	var bzImage string
	var upl string

	if ovmf = os.Getenv("VMTEST_OVMF"); len(ovmf) == 0 {
		t.Skipf("VMTEST_OVMF not set!!")
	}

	if _, err := os.Stat(ovmf); err != nil {
		t.Skipf("OVMF.fd image is not found: %s\n", ovmf)
	}

	if bzImage = os.Getenv("VMTEST_KERNEL"); len(bzImage) == 0 {
		t.Skipf("VMTEST_KERNEL not set!!")
	}

	if _, err := os.Stat(bzImage); err != nil {
		t.Skipf("Linux kernel image is not found: %s\n", bzImage)
	}

	if upl = os.Getenv("UROOT_TEST_UPLFIT"); len(upl) == 0 {
		t.Skipf("UROOT_TEST_UPLFIT not set!!")
	}

	if _, err := os.Stat(upl); err != nil {
		t.Skipf("UniversalPaylad image is not found: %s\n", upl)
	}

	vm := scriptvm.Start(t, "upl-vm", "",
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/init",
				"github.com/u-root/u-root/cmds/core/kexec",
				"github.com/u-root/u-root/cmds/core/gosh",
			),
			uimage.WithFiles(fmt.Sprintf("%s:/ext/upl", upl)),
			uimage.WithUinitCommand("/bbin/kexec /ext/upl"),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(5*time.Minute),
			qemu.ArbitraryArgs("-machine", "q35"),
			qemu.ArbitraryArgs("-m", "4096"),
			qemu.ArbitraryArgs("-bios", ovmf),
			qemu.ArbitraryArgs("-kernel", bzImage),
		),
	)

	if _, err := vm.Console.Expect(expect.All(
		// Boot target prompted from BDS
		expect.String("[Bds]Booting UEFI Shell"),
		// Last code before booting to UEFI Shell
		expect.String("PROGRESS CODE: V03058001 I0"),
	)); err != nil {
		t.Errorf("VM output did not match expectations: %v", err)
	}

	if err := vm.Kill(); err != nil {
		t.Errorf("Wait for VM process to be killed: %v", err)
	}

	vm.Wait()
}
