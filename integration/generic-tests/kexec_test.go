// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race
// +build !race

package integration

import (
	"os/exec"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest"
)

// TestMountKexec tests that kexec occurs correctly by checking the kernel cmdline.
// This is possible because the generic initramfs ensures that we mount the
// testdata directory containing the initramfs and kernel used in the VM.
func TestMountKexec(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" && vmtest.TestArch() != "arm64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		TestCmds: []string{
			"var CMDLINE = (cat /proc/cmdline)",
			"var SUFFIX = $CMDLINE[-7..]",
			"echo SAW $SUFFIX",
			"kexec -i /testdata/initramfs.cpio -c $CMDLINE' KEXEC=Y' /testdata/kernel",
		},
		QEMUOpts: qemu.Options{
			Timeout: 20 * time.Second,
			Devices: []qemu.Device{
				qemu.ArbitraryArgs{"-m", "8192"},
			},
		},
	})
	defer cleanup()

	if err := q.Expect("SAW KEXEC=Y"); err != nil {
		t.Fatal(err)
	}
}

// TestMountKexecLoad is same as TestMountKexec except it test calling
// kexec_load syscall than file load.
func TestMountKexecLoad(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" && vmtest.TestArch() != "arm64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	gzipP, err := exec.LookPath("gzip")
	if err != nil {
		t.Skipf("no gzip found, skip it as it won't be able to de-compress kernel")
	}

	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		BuildOpts: uroot.Opts{
			ExtraFiles: []string{gzipP},
		},
		TestCmds: []string{
			"var CMDLINE = (cat /proc/cmdline)",
			"var SUFFIX = $CMDLINE[-7..]",
			"echo SAW $SUFFIX",
			"kexec -d -i /testdata/initramfs.cpio --loadsyscall -c $CMDLINE' KEXEC=Y' /testdata/kernel",
		},
		QEMUOpts: qemu.Options{
			Timeout: 20 * time.Second,
			Devices: []qemu.Device{
				qemu.ArbitraryArgs{"-m", "8192"},
			},
		},
	})
	defer cleanup()

	if err := q.Expect("SAW KEXEC=Y"); err != nil {
		t.Fatal(err)
	}
}

// TestMountKexecLoadOnly test kexec loads without a kexec reboot.
func TestMountKexecLoadOnly(t *testing.T) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" && vmtest.TestArch() != "arm64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	gzipP, err := exec.LookPath("gzip")
	if err != nil {
		t.Skipf("no gzip found, skip it as it won't be able to de-compress kernel")
	}

	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		BuildOpts: uroot.Opts{
			ExtraFiles: []string{gzipP},
		},
		TestCmds: []string{
			"var CMDLINE = (cat /proc/cmdline)",
			"echo kexecloadresult ?(kexec -d -l -i /testdata/initramfs.cpio --loadsyscall -c $CMDLINE /testdata/kernel)",
		},
		QEMUOpts: qemu.Options{
			Timeout: 20 * time.Second,
			Devices: []qemu.Device{
				qemu.ArbitraryArgs{"-m", "8192"},
			},
		},
	})
	defer cleanup()

	if err := q.Expect("kexecloadresult $ok"); err != nil {
		t.Fatal(err)
	}
}

// TestMountKexecLoadCustomDTB test kexec_load a Arm64 Image with a user provided dtb.
func TestMountKexecLoadCustomDTB(t *testing.T) {
	if vmtest.TestArch() != "arm64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}
	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		TestCmds: []string{
			"var CMDLINE = (cat /proc/cmdline)",
			"var SUFFIX = $CMDLINE[-7..]",
			"echo SAW $SUFFIX",
			"cp /sys/firmware/fdt /tmp/userfdt",
			"kexec -d --dtb /tmp/userfdt   -i /testdata/initramfs.cpio --loadsyscall -c $CMDLINE' KEXEC=Y' /testdata/kernel",
		},
		QEMUOpts: qemu.Options{
			Timeout: 20 * time.Second,
			Devices: []qemu.Device{
				qemu.ArbitraryArgs{"-m", "8192"},
			},
		},
	})
	defer cleanup()

	if err := q.Expect("SAW KEXEC=Y"); err != nil {
		t.Fatal(err)
	}
}
