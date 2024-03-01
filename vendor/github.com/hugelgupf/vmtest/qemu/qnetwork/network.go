// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package qnetwork provides net device configurators for use with the Go qemu
// API.
package qnetwork

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/hugelgupf/vmtest/qemu"
)

// NIC is a QEMU NIC device string.
//
// Valid values for your QEMU can be found with	`qemu-system-<arch> -device
// help` in the Network devices section.
type NIC string

// A subset of QEMU NIC devices.
const (
	NICE1000     NIC = "e1000"
	NICVirtioNet NIC = "virtio-net"
)

// NetDevice is a definition of a NIC exposed to the guest & a backend to
// service that NIC.
type NetDevice[B Backend] struct {
	Device    Device
	Backend   B
	ExtraArgs []string
}

// NetDevModifier is a modifier of NetDevice.
type NetDevModifier[B Backend] func(netdevID string, alloc *qemu.IDAllocator, opts *qemu.Options, nd *NetDevice[B]) error

// Cmdline returns the full cmdline to add to QEMU args.
func (nd *NetDevice[B]) Cmdline(id string) []string {
	return append([]string{
		"-device", nd.Device.DevArgs(id),
		"-netdev", nd.Backend.NetDev(id),
	}, nd.ExtraArgs...)
}

// New adds a new NIC & network.
func New[B Backend](mods ...NetDevModifier[B]) qemu.Fn {
	return func(alloc *qemu.IDAllocator, opts *qemu.Options) error {
		netdevID := alloc.ID("netdev")

		nd := &NetDevice[B]{
			Device: Device{
				NIC: NICE1000,
				// Default MAC for the virtualized NIC.
				//
				// This is from the range of locally administered address ranges.
				MAC: net.HardwareAddr{0xe, 0, 0, 0, 0, 1},
			},
		}
		for _, mod := range mods {
			if mod != nil {
				if err := mod(netdevID, alloc, opts, nd); err != nil {
					return err
				}
			}
		}
		if err := nd.Backend.Validate(); err != nil {
			return err
		}
		opts.AppendQEMU(nd.Cmdline(netdevID)...)
		return nil
	}
}

// Device defines the device emulated by QEMU to the guest.
type Device struct {
	NIC  NIC
	MAC  net.HardwareAddr
	Args []string
}

// DeviceModifier is a function that modifies Device.
type DeviceModifier func(*Device) error

// WithNIC sets the NIC device exposed to the guest.
func WithNIC(n NIC) DeviceModifier {
	return func(d *Device) error {
		d.NIC = n
		return nil
	}
}

// WithMAC sets the MAC address exposed to the guest.
func WithMAC(mac net.HardwareAddr) DeviceModifier {
	if mac == nil {
		return nil
	}
	return func(d *Device) error {
		d.MAC = mac
		return nil
	}
}

// DevArgs returns the arg to "-device".
func (d *Device) DevArgs(id string) string {
	s := append([]string{string(d.NIC), "netdev=" + id, fmt.Sprintf("mac=%s", d.MAC)}, d.Args...)
	return strings.Join(s, ",")
}

// WithDevice adds device modifiers to NetDevice.
func WithDevice[B Backend](mods ...DeviceModifier) NetDevModifier[B] {
	return func(netdevID string, alloc *qemu.IDAllocator, opts *qemu.Options, nd *NetDevice[B]) error {
		for _, mod := range mods {
			if mod != nil {
				if err := mod(&nd.Device); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// InterVM is a Device that can connect multiple QEMU VMs to each other.
//
// InterVM uses the QEMU socket mechanism to connect multiple VMs with a simple
// unix domain socket.
type InterVM struct {
	socket string
	err    error

	// numVMs must be atomically accessed so VMs can be started in parallel
	// in goroutines.
	numVMs uint32

	wg sync.WaitGroup
}

// NewInterVM creates a new QEMU network between QEMU VMs.
//
// The network is closed from the world and only between the QEMU VMs.
func NewInterVM() *InterVM {
	// Avoid returning an error here if unnecessary.
	dir, err := os.MkdirTemp("", "intervm-")
	return &InterVM{
		err:    err,
		socket: filepath.Join(dir, "intervm.socket"),
	}
}

// NewVM returns a Device that can be used with a new QEMU VM.
func (n *InterVM) NewVM(mods ...NetDevModifier[SocketBackend]) qemu.Fn {
	if n == nil {
		return nil
	}
	if n.err != nil {
		return func(alloc *qemu.IDAllocator, opts *qemu.Options) error {
			return n.err
		}
	}

	newNum := atomic.AddUint32(&n.numVMs, 1)
	num := newNum - 1
	n.wg.Add(1)

	fn := []qemu.Fn{
		New[SocketBackend](
			append([]NetDevModifier[SocketBackend]{
				WithDevice[SocketBackend](WithMAC(net.HardwareAddr{0xe, 0, 0, 0, 0, byte(num)})),
				WithSocket(IsServer(num == 0), WithUnixSocket(n.socket)),
			}, mods...)...,
		),
	}
	if num == 0 {
		// When the server VM exits, wait until all clients
		// close, then delete the socket file and directory.
		fn = append(fn, qemu.WithTask(func(ctx context.Context, notif *qemu.Notifications) error {
			n.wg.Wait()
			return os.RemoveAll(filepath.Dir(n.socket))
		}))
	}

	// When each VM exits, call Done.
	fn = append(fn, qemu.WithTask(qemu.Cleanup(func() error {
		n.wg.Done()
		return nil
	})))
	return qemu.All(fn...)
}

// WithPCAP captures network traffic and saves it to outputFile.
func WithPCAP[B Backend](outputFile string) NetDevModifier[B] {
	return func(netdevID string, alloc *qemu.IDAllocator, opts *qemu.Options, nd *NetDevice[B]) error {
		nd.ExtraArgs = append(nd.ExtraArgs,
			"-object",
			fmt.Sprintf("filter-dump,id=%s,netdev=%s,file=%s", alloc.ID("filter"), netdevID, outputFile),
		)
		return nil
	}
}
