// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package integration

import (
	"os"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/uroot"
	"github.com/u-root/u-root/pkg/vmtest"
)

// TestMountKexec runs an init which mounts a filesystem and kexecs a kernel.
func RunTestMountKexec(t *testing.T, initramfs string) {
	// TODO: support arm
	if vmtest.TestArch() != "amd64" {
		t.Skipf("test not supported on %s", vmtest.TestArch())
	}

	if len(initramfs) == 0 {
		f, err := vmtest.CreateTestInitramfs(
			uroot.Opts{}, "github.com/u-root/u-root/integration/testcmd/kexec/uinit", "")
		if err != nil {
			t.Errorf("failed to create test initramfs: %v", err)
		}
		defer os.Remove(f)
		initramfs = f
	}

	// Create the CPIO and start QEMU.
	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		Initramfs: initramfs,
		QEMUOpts: qemu.Options{
			Timeout: 30 * time.Second,
		},
	})
	defer cleanup()

	if err := q.Expect("KEXECCOUNTER=0"); err != nil {
		t.Fatal(err)
	}
	if err := q.Expect("KEXECCOUNTER=1"); err != nil {
		t.Fatal(err)
	}
}
