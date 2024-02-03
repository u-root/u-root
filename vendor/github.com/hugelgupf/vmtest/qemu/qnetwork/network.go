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

// DeviceOptions are network device options.
//
// These are options for the `-device <nic>,...` command-line arg.
type DeviceOptions struct {
	// NIC is the NIC device that QEMU emulates.
	NIC NIC

	// MAC is the MAC address assigned to this interface in the guest.
	MAC net.HardwareAddr
}

// SetNIC sets the NIC.
func (d *DeviceOptions) setNIC(nic NIC) {
	d.NIC = nic
}

// SetMAC sets the device's MAC.
func (d *DeviceOptions) setMAC(mac net.HardwareAddr) {
	d.MAC = mac
}

// DeviceOptioner is an interface for setting DeviceOptions members.
//
// It exists so DeviceOptions can be extensible through generics, but the
// WithNIC/WithMAC/WithPCAP functions can be the same across all Options
// structs.
type DeviceOptioner interface {
	*UserOptions | *DeviceOptions

	setNIC(NIC)
	setMAC(net.HardwareAddr)
}

// Opt is a configurer useed with either *DeviceOptions or *UserOptions.
type Opt[DO DeviceOptioner] func(netdev string, id *qemu.IDAllocator, qopts *qemu.Options, opts DO) error

// WithPCAP captures network traffic and saves it to outputFile.
func WithPCAP[DO DeviceOptioner](outputFile string) Opt[DO] {
	return func(netdev string, id *qemu.IDAllocator, qopts *qemu.Options, opts DO) error {
		qopts.AppendQEMU(
			"-object",
			fmt.Sprintf("filter-dump,id=%s,netdev=%s,file=%s", id.ID("filter"), netdev, outputFile),
		)
		return nil
	}
}

// WithNIC changes the default NIC device QEMU emulates from e1000 to the given value.
func WithNIC[DO DeviceOptioner](nic NIC) Opt[DO] {
	return func(netdev string, id *qemu.IDAllocator, qopts *qemu.Options, opts DO) error {
		opts.setNIC(nic)
		return nil
	}
}

// WithMAC assigns a MAC address to the guest interface.
func WithMAC[DO DeviceOptioner](mac net.HardwareAddr) Opt[DO] {
	return func(netdev string, id *qemu.IDAllocator, qops *qemu.Options, opts DO) error {
		if mac != nil {
			opts.setMAC(mac)
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
func (n *InterVM) NewVM(nopts ...Opt[*DeviceOptions]) qemu.Fn {
	if n == nil {
		return nil
	}

	newNum := atomic.AddUint32(&n.numVMs, 1)
	num := newNum - 1

	n.wg.Add(1)
	return func(alloc *qemu.IDAllocator, qopts *qemu.Options) error {
		if n.err != nil {
			return n.err
		}
		devID := alloc.ID("vm")

		opts := DeviceOptions{
			// Default NIC.
			NIC: NICE1000,

			// MAC for the virtualized NIC.
			//
			// This is from the range of locally administered address ranges.
			MAC: net.HardwareAddr{0xe, 0, 0, 0, 0, byte(num)},
		}
		for _, opt := range nopts {
			if err := opt(devID, alloc, qopts, &opts); err != nil {
				return err
			}
		}
		args := []string{"-device", fmt.Sprintf("%s,netdev=%s,mac=%s", opts.NIC, devID, opts.MAC)}

		if num != 0 {
			args = append(args, "-netdev", fmt.Sprintf("stream,id=%s,server=false,addr.type=unix,addr.path=%s", devID, n.socket))
		} else {
			args = append(args, "-netdev", fmt.Sprintf("stream,id=%s,server=true,addr.type=unix,addr.path=%s", devID, n.socket))

			// When the server VM exits, wait until all clients
			// close, then delete the socket file and directory.
			qopts.Tasks = append(qopts.Tasks, func(ctx context.Context, notif *qemu.Notifications) error {
				n.wg.Wait()
				return os.RemoveAll(filepath.Dir(n.socket))
			})
		}

		// When each VM exits, call Done.
		qopts.Tasks = append(qopts.Tasks, qemu.Cleanup(func() error {
			n.wg.Done()
			return nil
		}))
		qopts.AppendQEMU(args...)
		return nil
	}
}

// UserOptions are options for a QEMU "user" network.
type UserOptions struct {
	DeviceOptions

	Args []string
}

// WithUserArg adds more comma-separated args to a `-netdev user,arg0,arg1,...`
// invocation.
func WithUserArg(arg ...string) Opt[*UserOptions] {
	return func(netdev string, id *qemu.IDAllocator, qopts *qemu.Options, opts *UserOptions) error {
		opts.Args = append(opts.Args, arg...)
		return nil
	}
}

// IPv4HostNetwork provides QEMU user-mode networking to the host.
//
// Net must be an IPv4 network.
//
// Default NIC is e1000, with a MAC address of 0e:00:00:00:00:01.
func IPv4HostNetwork(cidr string, nopts ...Opt[*UserOptions]) qemu.Fn {
	return func(alloc *qemu.IDAllocator, qopts *qemu.Options) error {
		// TODO: use IP
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return err
		}
		if ipnet.IP.To4() == nil {
			return fmt.Errorf("HostNetwork must be configured with an IPv4 address")
		}

		netdevID := alloc.ID("netdev")
		opts := UserOptions{
			DeviceOptions: DeviceOptions{
				// Default NIC.
				NIC: NICE1000,

				// MAC for the virtualized NIC.
				//
				// This is from the range of locally administered address ranges.
				MAC: net.HardwareAddr{0xe, 0, 0, 0, 0, 1},
			},
		}

		for _, opt := range nopts {
			if err := opt(netdevID, alloc, qopts, &opts); err != nil {
				return err
			}
		}

		var extraUserArgs string
		if len(opts.Args) > 0 {
			extraUserArgs = "," + strings.Join(opts.Args, ",")
		}
		qopts.AppendQEMU(
			"-device", fmt.Sprintf("%s,netdev=%s,mac=%s", opts.NIC, netdevID, opts.MAC),
			"-netdev", fmt.Sprintf("user,id=%s,net=%s,dhcpstart=%s,ipv6=off%s", netdevID, ipnet, nthIP(ipnet, 8), extraUserArgs),
		)
		return nil
	}
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func nthIP(nt *net.IPNet, n int) net.IP {
	ip := make(net.IP, net.IPv4len)
	copy(ip, nt.IP.To4())
	for i := 0; i < n; i++ {
		inc(ip)
	}
	if !nt.Contains(ip) {
		return nil
	}
	return ip
}
