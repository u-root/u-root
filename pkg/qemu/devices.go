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
	// KArgs returns arguments that must be passed to the kernel for this device, or nil.
	KArgs() []string
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

// NewNetwork creates a new QEMU network between QEMU VMs.
//
// The network is closed from the world and only between the QEMU VMs.
func NewNetwork() *Network {
	return &Network{
		port: 1234,
	}
}

// NewVM returns a Device that can be used with a new QEMU VM.
func (n *Network) NewVM() Device {
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
	return networkImpl{args}
}

type networkImpl struct {
	args []string
}

func (n networkImpl) Cmdline() []string {
	return n.args
}

func (n networkImpl) KArgs() []string { return nil }

// ReadOnlyDirectory is a Device that exposes a directory as a /dev/sda1
// readonly vfat partition in the VM.
type ReadOnlyDirectory struct {
	// Dir is the directory to expose as a read-only vfat partition.
	Dir string
}

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

func (ReadOnlyDirectory) KArgs() []string { return nil }

// P9Directory is a Device that exposes a directory as a Plan9 (9p)
// read-write filesystem in the VM.
type P9Directory struct {
	// Dir is the directory to expose as read-write 9p filesystem.
	Dir string

	// Boot: if true, indicates this is the root volume. For this to work,
	// kernel args will need to be added - use KArgs() to get the args. There
	// can only be one boot 9pfs at a time.
	Boot bool

	// Tag is an identifier that is used within the VM when mounting an fs,
	// e.g. 'mount -t 9p my-vol-ident mountpoint'. If not specified, a default
	// tag of 'tmpdir' will be used.
	//
	// Ignored if Boot is true, as the tag in that case is special.
	//
	// For non-boot devices, this is also used as the id linking the `-fsdev`
	// and `-device` args together.
	//
	// Because the tag must be unique for each dir, if multiple non-boot
	// P9Directory's are used, tag may be omitted for no more than one.
	Tag string
}

func (p P9Directory) Cmdline() []string {
	if len(p.Dir) == 0 {
		return nil
	}

	var tag, id string
	if p.Boot {
		tag = "/dev/root"
	} else {
		tag = p.Tag
		if len(tag) == 0 {
			tag = "tmpdir"
		}
	}
	if p.Boot {
		id = "rootdrv"
	} else {
		id = tag
	}

	// Expose the temp directory to QEMU
	return []string{
		// security_model=mapped-file seems to be the best choice. It gives
		// us control over uid/gid/mode seen in the guest, without requiring
		// elevated perms on the host.
		"-fsdev", fmt.Sprintf("local,id=%s,path=%s,security_model=mapped-file", id, p.Dir),
		"-device", fmt.Sprintf("virtio-9p-pci,fsdev=%s,mount_tag=%s", id, tag),
	}
}

func (p P9Directory) KArgs() []string {
	if len(p.Dir) == 0 {
		return nil
	}
	if p.Boot {
		return []string{
			"devtmpfs.mount=1",
			"root=/dev/root",
			"rootfstype=9p",
			"rootflags=trans=virtio,version=9p2000.L",
		}
	}
	return []string{
		//seen as an env var by the init process
		"UROOT_USE_9P=1",
	}
}

// VirtioRandom is a Device that exposes a PCI random number generator to the
// QEMU VM.
type VirtioRandom struct{}

func (VirtioRandom) Cmdline() []string {
	return []string{"-device", "virtio-rng-pci"}
}

func (VirtioRandom) KArgs() []string { return nil }

// ArbitraryArgs is a Device that allows users to add arbitrary arguments to
// the QEMU command line.
type ArbitraryArgs []string

func (aa ArbitraryArgs) Cmdline() []string {
	return aa
}

func (ArbitraryArgs) KArgs() []string { return nil }
