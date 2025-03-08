// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build arm64 && !race

package integration

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/u-root/gobusybox/src/pkg/golang"
	"github.com/u-root/mkuimage/uimage"
)

// TestUPLBootArm64 tests '/bbin/kexec UplFitARM64.fit' boot to UEFI Shell.
func TestUPLBootArm64(t *testing.T) {
	// CAUTION: Before running this VM test, environment variable must be set:
	//     "VMTEST_ARCH=arm64"
	// Since "github.com/hugelgupf/vmtest/qemu" package parses .vmtest.yaml
	// with $VMTEST_ARCH to get $VMTEST_QEMU.

	qemu.SkipIfNotArch(t, qemu.ArchArm64)

	// Check required images, including:
	//   QEMU_EFI.fd		-- UEFI OVMF binary
	//   Image 				-- Linuxboot kernel Image
	//   UplFitARM64.fit	-- UPL Arm64 image with FDT enabled

	var ovmf string
	var image string
	var upl string

	if img := os.Getenv("VMTEST_OVMF"); len(img) == 0 {
		t.Skipf("VMTEST_OVMF not set!!")
	} else {
		ovmf = img
	}

	if _, err := os.Stat(ovmf); err != nil && os.IsNotExist(err) {
		t.Skipf("OVMF.fd image is not found: %s\n", ovmf)
	}

	if img := os.Getenv("VMTEST_KERNEL"); len(img) == 0 {
		t.Skipf("VMTEST_KERNEL not set!!")
	} else {
		image = img
	}

	if _, err := os.Stat(image); err != nil && os.IsNotExist(err) {
		t.Skipf("Linux kernel image image is not found: %s\n", image)
	}

	if img := os.Getenv("UROOT_TEST_UPLFIT"); len(img) == 0 {
		t.Skipf("UROOT_TEST_UPLFIT not set!!")
	} else {
		upl = img
	}

	if _, err := os.Stat(upl); err != nil && os.IsNotExist(err) {
		t.Skipf("UniversalPaylad image is not found: %s\n", upl)
	}

	vm := scriptvm.Start(t, "upl-vm", "",
		scriptvm.WithUimage(
			uimage.WithEnv(golang.WithGOARCH("arm64")),
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
			qemu.ArbitraryArgs("-machine", "virt,gic-version=3"),
			qemu.ArbitraryArgs("-m", "4096"),
			qemu.ArbitraryArgs("-bios", ovmf),
			qemu.ArbitraryArgs("-kernel", image),
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
		fmt.Printf("Wait for VM process to be killed: %v\n", err)
		t.Errorf("Wait for VM process to be killed: %v", err)
	}

	vm.Wait()
}
