// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !race

package integration

import (
	"os"
	"testing"
	"time"

	expect "github.com/google/goexpect"
	"github.com/u-root/u-root/pkg/qemu"
	"github.com/u-root/u-root/pkg/vmtest"
)

func TestESXi(t *testing.T) {
	img := os.Getenv("UROOT_ESXI_IMAGE")
	if _, err := os.Stat(img); err != nil && os.IsNotExist(err) {
		t.Skip("ESXi disk image is not present. Please set UROOT_ESXI_IMAGE= to an installed ESXi disk image.")
	}

	q, cleanup := vmtest.QEMUTest(t, &vmtest.Options{
		TestCmds: []string{
			`esxiboot -d="/dev/sda" --append="vmkBootVerbose=TRUE vmbLog=TRUE debugLogToSerial=1 logPort=com1"`,
		},
		QEMUOpts: qemu.Options{
			Devices: []qemu.Device{
				qemu.IDEBlockDevice{File: img},
				// If at some point we get virtio-net working
				// again in ESXi, you may need to set num CPUs
				// to 4.
				qemu.ArbitraryArgs{"-smp", "2"},

				// QEMU is not good enough to do ESXi without KVM. You'll get kernel panic.
				qemu.ArbitraryArgs{"-enable-kvm"},

				// IvyBridge is the lowest-common-denominator.
				// ESXi 7.0 drops support for SandyBridge
				// afaict.
				qemu.ArbitraryArgs{"-cpu", "IvyBridge"},

				// Min ESXi requirement is 4G of memory, but in
				// ESXi 7.0 some plugins fail to load at 4G
				// under memory pressure, and we never get to
				// "Boot Successful"
				qemu.ArbitraryArgs{"-m", "8192"},
			},
		},
	})
	defer cleanup()

	if out, err := q.ExpectBatch([]expect.Batcher{
		// Last Linux print.
		&expect.BExp{R: "kexec_core: Starting new kernel"},
		// First b.b00 ESXi first-stage kernel print
		&expect.BExp{R: "Serial port initialized..."},
		// ~thirdish k.b00 ESXi second-stage kernel print
		&expect.BExp{R: "Parsing command line boot options"},
		// When we can be confident it's done.
		&expect.BExp{R: "Boot Successful"},
	}, 2*time.Minute); err != nil {
		t.Fatalf("VM output did not match expectations: %v", out)
	}
}
