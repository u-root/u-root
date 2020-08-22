// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package integration

import (
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/qemu"
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
			"CMDLINE = (cat /proc/cmdline)",
			"SUFFIX = $CMDLINE[-7:]",
			"echo SAW $SUFFIX",
			"kexec -i /testdata/initramfs.cpio -c $CMDLINE' KEXEC=Y' /testdata/kernel",
		},
		QEMUOpts: qemu.Options{
			Timeout: 20 * time.Second,
		},
	})
	defer cleanup()

	if err := q.Expect("SAW KEXEC=Y"); err != nil {
		t.Fatal(err)
	}
}
