// Copyright 2016 The Netstack Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ports provides PortManager that manages allocating, reserving and releasing ports.
package ports

import (
	"math"
	"math/rand"
	"sync"

	"github.com/google/netstack/tcpip"
)

const (
	// firstEphemeral is the first ephemeral port.
	firstEphemeral uint16 = 16000

	anyIPAddress = tcpip.Address("")
)

type portDescriptor struct {
	network   tcpip.NetworkProtocolNumber
	transport tcpip.TransportProtocolNumber
	port      uint16
}

// PortManager manages allocating, reserving and releasing ports.
type PortManager struct {
	mu             sync.RWMutex
	allocatedPorts map[portDescriptor]bindAddresses
}

// bindAddresses is a set of IP addresses.
type bindAddresses map[tcpip.Address]struct{}

// isAvailable checks whether an IP address is available to bind to.
func (b bindAddresses) isAvailable(addr tcpip.Address) bool {
	if addr == anyIPAddress {
		return len(b) == 0
	}

	// If all addresses for this portDescriptor are already bound, no
	// address is available.
	if _, ok := b[anyIPAddress]; ok {
		return false
	}

	if _, ok := b[addr]; ok {
		return false
	}
	return true
}

// NewPortManager creates new PortManager.
func NewPortManager() *PortManager {
	return &PortManager{allocatedPorts: make(map[portDescriptor]bindAddresses)}
}

// PickEphemeralPort randomly chooses a starting point and iterates over all
// possible ephemeral ports, allowing the caller to decide whether a given port
// is suitable for its needs, and stopping when a port is found or an error
// occurs.
func (s *PortManager) PickEphemeralPort(testPort func(p uint16) (bool, *tcpip.Error)) (port uint16, err *tcpip.Error) {
	count := uint16(math.MaxUint16 - firstEphemeral + 1)
	offset := uint16(rand.Int31n(int32(count)))

	for i := uint16(0); i < count; i++ {
		port = firstEphemeral + (offset+i)%count
		ok, err := testPort(port)
		if err != nil {
			return 0, err
		}

		if ok {
			return port, nil
		}
	}

	return 0, tcpip.ErrNoPortAvailable
}

// ReservePort marks a port/IP combination as reserved so that it cannot be
// reserved by another endpoint. If port is zero, ReservePort will search for
// an unreserved ephemeral port and reserve it, returning its value in the
// "port" return value.
func (s *PortManager) ReservePort(network []tcpip.NetworkProtocolNumber, transport tcpip.TransportProtocolNumber, addr tcpip.Address, port uint16) (reservedPort uint16, err *tcpip.Error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If a port is specified, just try to reserve it for all network
	// protocols.
	if port != 0 {
		if !s.reserveSpecificPort(network, transport, addr, port) {
			return 0, tcpip.ErrPortInUse
		}
		return port, nil
	}

	// A port wasn't specified, so try to find one.
	return s.PickEphemeralPort(func(p uint16) (bool, *tcpip.Error) {
		return s.reserveSpecificPort(network, transport, addr, p), nil
	})
}

// reserveSpecificPort tries to reserve the given port on all given protocols.
func (s *PortManager) reserveSpecificPort(network []tcpip.NetworkProtocolNumber, transport tcpip.TransportProtocolNumber, addr tcpip.Address, port uint16) bool {
	// Check that the port is available on all network protocols.
	desc := portDescriptor{0, transport, port}
	for _, n := range network {
		desc.network = n
		if addrs, ok := s.allocatedPorts[desc]; ok {
			if !addrs.isAvailable(addr) {
				return false
			}
		}
	}

	// Reserve port on all network protocols.
	for _, n := range network {
		desc.network = n
		m, ok := s.allocatedPorts[desc]
		if !ok {
			m = make(bindAddresses)
			s.allocatedPorts[desc] = m
		}
		m[addr] = struct{}{}
	}

	return true
}

// ReleasePort releases the reservation on a port/IP combination so that it can
// be reserved by other endpoints.
func (s *PortManager) ReleasePort(network []tcpip.NetworkProtocolNumber, transport tcpip.TransportProtocolNumber, addr tcpip.Address, port uint16) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, n := range network {
		desc := portDescriptor{n, transport, port}
		m := s.allocatedPorts[desc]
		delete(m, addr)
		if len(m) == 0 {
			delete(s.allocatedPorts, desc)
		}
	}
}
