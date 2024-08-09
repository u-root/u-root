// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !race

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

func TestESXi(t *testing.T) {
	img := os.Getenv("UROOT_ESXI_IMAGE")
	if _, err := os.Stat(img); err != nil && os.IsNotExist(err) {
		t.Skip("ESXi disk image is not present. Please set UROOT_ESXI_IMAGE= to an installed ESXi disk image.")
	}

	script := `esxiboot -d="/dev/sda" --append="vmkBootVerbose=TRUE vmbLog=TRUE debugLogToSerial=1 logPort=com1"`
	vm := scriptvm.Start(t, "vm", script,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/exp/esxiboot",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			qemu.IDEBlockDevice(img),
			qemu.VirtioRandom(),
			// If at some point we get virtio-net working
			// again in ESXi, you may need to set num CPUs
			// to 4.
			qemu.ArbitraryArgs("-smp", "2"),

			// QEMU is not good enough to do ESXi without KVM. You'll get kernel panic.
			qemu.ArbitraryArgs("-enable-kvm"),

			// IvyBridge is the lowest-common-denominator.
			// ESXi 7.0 drops support for SandyBridge
			// afaict.
			qemu.ArbitraryArgs("-cpu", "IvyBridge"),

			// Min ESXi requirement is 4G of memory, but in
			// ESXi 7.0 some plugins fail to load at 4G
			// under memory pressure, and we never get to
			// "Boot Successful"
			qemu.ArbitraryArgs("-m", "8192"),
		),
	)

	if _, err := vm.Console.Expect(expect.All(
		// First b.b00 ESXi first-stage kernel print
		expect.String("Serial port initialized..."),
		// ~thirdish k.b00 ESXi second-stage kernel print
		expect.String("Parsing command line boot options"),
		// When we can be confident it's done.
		expect.String("Boot Successful"),
	)); err != nil {
		t.Errorf("VM output did not match expectations: %v", err)
	}

	if err := vm.Kill(); err != nil {
		t.Errorf("Wait: %v", err)
	}
	_ = vm.Wait()
}

func TestESXiNVMe(t *testing.T) {
	img := os.Getenv("UROOT_ESXI_IMAGE")
	if _, err := os.Stat(img); err != nil && os.IsNotExist(err) {
		t.Skip("ESXi disk image is not present. Please set UROOT_ESXI_IMAGE= to an installed ESXi disk image.")
	}

	script := `esxiboot -d="/dev/nvme0n1" --append="vmkBootVerbose=TRUE vmbLog=TRUE debugLogToSerial=1 logPort=com1"`
	vm := scriptvm.Start(t, "vm", script,
		scriptvm.WithUimage(
			uimage.WithBusyboxCommands(
				"github.com/u-root/u-root/cmds/exp/esxiboot",
			),
		),
		scriptvm.WithQEMUFn(
			qemu.WithVMTimeout(2*time.Minute),
			qemu.VirtioRandom(),
			qemu.ArbitraryArgs("-device", "nvme,drive=NVME1,serial=nvme-1"),
			qemu.ArbitraryArgs("-drive", fmt.Sprintf("file=%s,if=none,id=NVME1", img)),
			// If at some point we get virtio-net working
			// again in ESXi, you may need to set num CPUs
			// to 4.
			qemu.ArbitraryArgs("-smp", "2"),

			// QEMU is not good enough to do ESXi without KVM. You'll get kernel panic.
			qemu.ArbitraryArgs("-enable-kvm"),

			// IvyBridge is the lowest-common-denominator.
			// ESXi 7.0 drops support for SandyBridge
			// afaict.
			qemu.ArbitraryArgs("-cpu", "IvyBridge"),

			// Min ESXi requirement is 4G of memory, but in
			// ESXi 7.0 some plugins fail to load at 4G
			// under memory pressure, and we never get to
			// "Boot Successful"
			qemu.ArbitraryArgs("-m", "8192"),
		),
	)

	if _, err := vm.Console.Expect(expect.All(
		// First b.b00 ESXi first-stage kernel print
		expect.String("Serial port initialized..."),
		// ~thirdish k.b00 ESXi second-stage kernel print
		expect.String("Parsing command line boot options"),
		// When we can be confident it's done.
		expect.String("Boot Successful"),
	)); err != nil {
		t.Errorf("VM output did not match expectations: %v", err)
	}

	if err := vm.Kill(); err != nil {
		t.Errorf("Wait: %v", err)
	}
	_ = vm.Wait()
}
