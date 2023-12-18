// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package network provides net device configurators for use with the Go qemu
// API.
package network

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/hugelgupf/vmtest/qemu"
)

// InterVM is a Device that can connect multiple QEMU VMs to each other.
//
// InterVM uses the QEMU socket mechanism to connect multiple VMs with a simple
// TCP socket.
type InterVM struct {
	port uint16

	// numVMs must be atomically accessed so VMs can be started in parallel
	// in goroutines.
	numVMs uint32
}

// NewInterVM creates a new QEMU network between QEMU VMs.
//
// The network is closed from the world and only between the QEMU VMs.
func NewInterVM() *InterVM {
	return &InterVM{
		port: 1234,
	}
}

// Opt returns additional QEMU command-line parameters based on the net
// device ID.
type Opt func(netdev string, id *qemu.IDAllocator) []string

// WithPCAP captures network traffic and saves it to outputFile.
func WithPCAP(outputFile string) Opt {
	return func(netdev string, id *qemu.IDAllocator) []string {
		return []string{
			"-object",
			fmt.Sprintf("filter-dump,id=%s,netdev=%s,file=%s", id.ID("filter"), netdev, outputFile),
		}
	}
}

// NewVM returns a Device that can be used with a new QEMU VM.
func (n *InterVM) NewVM(nopts ...Opt) qemu.Fn {
	if n == nil {
		return nil
	}

	newNum := atomic.AddUint32(&n.numVMs, 1)
	num := newNum - 1

	// MAC for the virtualized NIC.
	//
	// This is from the range of locally administered address ranges.
	mac := net.HardwareAddr{0x0e, 0x00, 0x00, 0x00, 0x00, byte(num)}
	return func(alloc *qemu.IDAllocator, opts *qemu.Options) error {
		devID := alloc.ID("vm")

		args := []string{"-device", fmt.Sprintf("e1000,netdev=%s,mac=%s", devID, mac)}
		// Note: QEMU in CircleCI seems to in solve cases fail when using just ':1234' format.
		// It fails with "address resolution failed for :1234: Name or service not known"
		// hinting that this is somehow related to DNS resolution. To work around this,
		// we explicitly bind to 127.0.0.1 (IPv6 [::1] is not parsed correctly by QEMU).
		if num != 0 {
			args = append(args, "-netdev", fmt.Sprintf("socket,id=%s,connect=127.0.0.1:%d", devID, n.port))
		} else {
			args = append(args, "-netdev", fmt.Sprintf("socket,id=%s,listen=127.0.0.1:%d", devID, n.port))
		}

		for _, opt := range nopts {
			args = append(args, opt(devID, alloc)...)
		}
		opts.AppendQEMU(args...)
		return nil
	}
}

// IPv4HostNetwork provides QEMU user-mode networking to the host.
//
// Net must be an IPv4 network.
// It uses the e1000 NIC.
func IPv4HostNetwork(net *net.IPNet, nopts ...Opt) qemu.Fn {
	return func(alloc *qemu.IDAllocator, opts *qemu.Options) error {
		if net.IP.To4() == nil {
			return fmt.Errorf("HostNetwork must be configured with an IPv4 address")
		}

		netdevID := alloc.ID("netdev")
		args := []string{
			"-device", fmt.Sprintf("e1000,netdev=%s", netdevID),
			"-netdev", fmt.Sprintf("user,id=%s,net=%s,dhcpstart=%s,ipv6=off", netdevID, net, nthIP(net, 8)),
		}

		for _, opt := range nopts {
			args = append(args, opt(netdevID, alloc)...)
		}
		opts.AppendQEMU(args...)
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
