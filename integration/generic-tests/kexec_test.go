// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/hugelgupf/vmtest/qemu"
	"github.com/hugelgupf/vmtest/scriptvm"
	"github.com/hugelgupf/vmtest/testtmp"
	"github.com/u-root/mkuimage/uimage"
)

// TestMountKexec tests that kexec occurs correctly by checking the kernel cmdline.
// This is possible because the generic initramfs ensures that we mount the
// testdata directory containing the initramfs and kernel used in the VM.
func TestMountKexec(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchAMD64, qemu.ArchArm64)

	script := `
		CMDLINE=$(cat /proc/cmdline)
		SUFFIX=${CMDLINE:(-7)}
		echo SAW $SUFFIX
		kexec -l -i /mount/9p/initramfs/initramfs.cpio -c "${CMDLINE} KEXEC=Y" /kernel
		sync
		kexec -e
	`

	initrd := filepath.Join(testtmp.TempDir(t), "initramfs.cpio")
	vm := scriptvm.Start(t, "vm", script,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/cat",
				"github.com/u-root/u-root/cmds/core/sync",
			),
			// Build kexec as a binary command to get accurate GOCOVERDIR
			// integration coverage data (busybox rewrites command code).
			uimage.WithCoveredCommands("github.com/u-root/u-root/cmds/core/kexec"),
			uimage.WithFiles(fmt.Sprintf("%s:kernel", os.Getenv("VMTEST_KERNEL"))),
			uimage.WithCPIOOutput(initrd),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			qemu.ArbitraryArgs("-m", "8192"),
			qemu.WithInitramfs(initrd),
			// Initramfs available at /mount/9p/initramfs/initramfs.cpio.
			qemu.P9Directory(filepath.Dir(initrd), "initramfs"),
		),
	)

	if _, err := vm.Console.ExpectString("SAW KEXEC=Y"); err != nil {
		t.Fatal(err)
	}
	if err := vm.Kill(); err != nil {
		t.Errorf("Kill: %v", err)
	}
	_ = vm.Wait()
}

// TestMountKexecLoad is same as TestMountKexec except it test calling
// kexec_load syscall than file load.
func TestMountKexecLoad(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchAMD64, qemu.ArchArm64)

	gzipP, err := exec.LookPath("gzip")
	if err != nil {
		t.Skipf("no gzip found, skip it as it won't be able to de-compress kernel")
	}

	script := `
		CMDLINE=$(cat /proc/cmdline)
		SUFFIX=${CMDLINE:(-7)}
		echo SAW $SUFFIX
		kexec -l -d -i /mount/9p/initramfs/initramfs.cpio --loadsyscall -c "${CMDLINE} KEXEC=Y" /kernel
		sync
		kexec -e
	`

	initrd := filepath.Join(testtmp.TempDir(t), "initramfs.cpio")
	vm := scriptvm.Start(t, "vm", script,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/cat",
				"github.com/u-root/u-root/cmds/core/sync",
			),
			// Build kexec as a binary command to get accurate GOCOVERDIR
			// integration coverage data (busybox rewrites command code).
			uimage.WithCoveredCommands("github.com/u-root/u-root/cmds/core/kexec"),
			uimage.WithFiles(
				fmt.Sprintf("%s:kernel", os.Getenv("VMTEST_KERNEL")),
				gzipP,
			),
			uimage.WithCPIOOutput(initrd),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			qemu.ArbitraryArgs("-m", "8192"),
			qemu.WithInitramfs(initrd),
			// Initramfs available at /mount/9p/initramfs/initramfs.cpio.
			qemu.P9Directory(filepath.Dir(initrd), "initramfs"),
		),
	)
	if _, err := vm.Console.ExpectString("SAW KEXEC=Y"); err != nil {
		t.Error(err)
	}
	if err := vm.Kill(); err != nil {
		t.Errorf("Kill: %v", err)
	}
	_ = vm.Wait()
}

// TestMountKexecLoadOnly test kexec loads without a kexec reboot.
func TestMountKexecLoadOnly(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchAMD64, qemu.ArchArm64)

	gzipP, err := exec.LookPath("gzip")
	if err != nil {
		t.Skipf("no gzip found, skip it as it won't be able to de-compress kernel")
	}

	script := `
		CMDLINE=$(cat /proc/cmdline)
		kexec -d -l -i /mount/9p/initramfs/initramfs.cpio --loadsyscall -c "${CMDLINE}" /kernel
		echo kexecloadresult $?
	`

	initrd := filepath.Join(testtmp.TempDir(t), "initramfs.cpio")
	vm := scriptvm.Start(t, "vm", script,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/cat",
			),
			// Build kexec as a binary command to get accurate GOCOVERDIR
			// integration coverage data (busybox rewrites command code).
			uimage.WithCoveredCommands("github.com/u-root/u-root/cmds/core/kexec"),
			uimage.WithFiles(
				fmt.Sprintf("%s:kernel", os.Getenv("VMTEST_KERNEL")),
				gzipP,
			),
			uimage.WithCPIOOutput(initrd),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			qemu.ArbitraryArgs("-m", "8192"),
			qemu.WithInitramfs(initrd),
			// Initramfs available at /mount/9p/initramfs/initramfs.cpio.
			qemu.P9Directory(filepath.Dir(initrd), "initramfs"),
		),
	)

	if _, err := vm.Console.ExpectString("kexecloadresult 0"); err != nil {
		t.Error(err)
	}
	if err := vm.Wait(); err != nil {
		t.Errorf("Wait: %v", err)
	}
}

// TestMountKexecLoadCustomDTB test kexec_load a Arm64 Image with a user provided dtb.
func TestMountKexecLoadCustomDTB(t *testing.T) {
	qemu.SkipIfNotArch(t, qemu.ArchArm64)

	script := `
		CMDLINE=$(cat /proc/cmdline)
		SUFFIX=${CMDLINE:(-7)}
		echo SAW $SUFFIX
		cp /sys/firmware/fdt /tmp/userfdt
		kexec -d --dtb /tmp/userfdt -i /mount/9p/initramfs/initramfs.cpio --loadsyscall -c "${CMDLINE} KEXEC=Y" /kernel
	`
	initrd := filepath.Join(testtmp.TempDir(t), "initramfs.cpio")
	vm := scriptvm.Start(t, "vm", script,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/core/cat",
				"github.com/u-root/u-root/cmds/core/cp",
			),
			// Build kexec as a binary command to get accurate GOCOVERDIR
			// integration coverage data (busybox rewrites command code).
			uimage.WithCoveredCommands("github.com/u-root/u-root/cmds/core/kexec"),
			uimage.WithFiles(
				fmt.Sprintf("%s:kernel", os.Getenv("VMTEST_KERNEL")),
			),
			uimage.WithCPIOOutput(initrd),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(time.Minute),
			qemu.ArbitraryArgs("-m", "8192"),
			qemu.WithInitramfs(initrd),
			// Initramfs available at /mount/9p/initramfs/initramfs.cpio.
			qemu.P9Directory(filepath.Dir(initrd), "initramfs"),
		),
	)

	if _, err := vm.Console.ExpectString("SAW KEXEC=Y"); err != nil {
		t.Error(err)
	}
	if err := vm.Kill(); err != nil {
		t.Errorf("Kill: %v", err)
	}
	_ = vm.Wait()
}
