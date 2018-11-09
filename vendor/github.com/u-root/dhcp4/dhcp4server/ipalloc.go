// Copyright 2018 the u-root Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhcp4server

import (
	"encoding/binary"
	"fmt"
	"net"
)

func nextIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}

type ipAllocator struct {
	// subnet is the range of IP addresses that can be allocated by this
	// DHCP server.
	subnet *net.IPNet

	// allocated is the set of IP addresses currently allocated to a
	// client.
	allocated map[uint32]struct{}
}

func newIPAllocator(subnet *net.IPNet) *ipAllocator {
	return &ipAllocator{
		subnet:    subnet,
		allocated: make(map[uint32]struct{}),
	}
}

func ipToUint32(ip net.IP) uint32 {
	return binary.LittleEndian.Uint32(ip.To4())
}

func (ia *ipAllocator) usable(ip net.IP) bool {
	_, ok := ia.allocated[ipToUint32(ip)]
	return !ok
}

func (ia *ipAllocator) grab(ip net.IP) bool {
	if !ia.subnet.Contains(ip) || !ia.usable(ip) {
		return false
	}
	ia.allocated[ipToUint32(ip)] = struct{}{}
	return true
}

func (ia *ipAllocator) alloc() net.IP {
	// Make a copy so we can modify it.
	try := make([]byte, len(ia.subnet.IP))
	copy(try, ia.subnet.IP)

	// Just try em all.
	for ia.subnet.Contains(try) {
		if ia.usable(try) {
			ia.allocated[ipToUint32(try)] = struct{}{}
			return try
		}
		nextIP(try)
	}
	return nil
}

func (ia *ipAllocator) free(ip net.IP) error {
	if ia.usable(ip) {
		return fmt.Errorf("cannot free unallocated IP %v", ip)
	}
	delete(ia.allocated, ipToUint32(ip))
	return nil
}
