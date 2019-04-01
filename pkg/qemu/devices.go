// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qemu

import (
	"fmt"
	"net"
	"sync/atomic"
)

// Device is a QEMU device to expose to a VM.
type Device interface {
	// Cmdline returns arguments to append to the QEMU command line for this device.
	Cmdline() []string
}

// Network is a Device that can connect multiple QEMU VMs to each other.
//
// Network uses the QEMU socket mechanism to connect multiple VMs with a simple
// TCP socket.
type Network struct {
	port uint16

	// numVMs must be atomically accessed so VMs can be started in parallel
	// in goroutines.
	numVMs uint32
}

func NewNetwork() *Network {
	return &Network{
		port: 1234,
	}
}

// Cmdline implements Device.
func (n *Network) Cmdline() []string {
	if n == nil {
		return nil
	}

	newNum := atomic.AddUint32(&n.numVMs, 1)
	num := newNum - 1

	// MAC for the virtualized NIC.
	//
	// This is from the range of locally administered address ranges.
	mac := net.HardwareAddr{0x0e, 0x00, 0x00, 0x00, 0x00, byte(num)}

	args := []string{"-net", fmt.Sprintf("nic,macaddr=%s", mac)}
	// Note: QEMU in CircleCI seems to in solve cases fail when using just ':1234' format.
	// It fails with "address resolution failed for :1234: Name or service not known"
	// hinting that this is somehow related to DNS resolution. To work around this,
	// we explicitly bind to 127.0.0.1 (IPv6 [::1] is not parsed correctly by QEMU).
	if num != 0 {
		args = append(args, "-net", fmt.Sprintf("socket,connect=127.0.0.1:%d", n.port))
	} else {
		args = append(args, "-net", fmt.Sprintf("socket,listen=127.0.0.1:%d", n.port))
	}
	return args
}

// ReadOnlyDirectory is a Device that exposes a directory as a /dev/sda1
// readonly vfat partition in the VM.
type ReadOnlyDirectory struct {
	// Dir is the directory to expose as a read-only vfat partition.
	Dir string
}

// Cmdline implements Device.
func (rod ReadOnlyDirectory) Cmdline() []string {
	if len(rod.Dir) == 0 {
		return nil
	}

	// Expose the temp directory to QEMU as /dev/sda1
	return []string{
		"-drive", fmt.Sprintf("file=fat:rw:%s,if=none,id=tmpdir", rod.Dir),
		"-device", "ich9-ahci,id=ahci",
		"-device", "ide-drive,drive=tmpdir,bus=ahci.0",
	}
}

// VirtioRandom exposes a PCI random number generator Device to the QEMU VM.
type VirtioRandom struct{}

// Cmdline implements Device.
func (VirtioRandom) Cmdline() []string {
	return []string{"-device", "virtio-rng-pci"}
}

// ArbitraryArgs allows users to add arbitrary arguments to the QEMU command
// line.
type ArbitraryArgs []string

// Cmdline implements Device.
func (aa ArbitraryArgs) Cmdline() []string {
	return aa
}
