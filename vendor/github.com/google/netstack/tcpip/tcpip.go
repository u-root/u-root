// Copyright 2016 The Netstack Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tcpip provides the interfaces and related types that users of the
// tcpip stack will use in order to create endpoints used to send and receive
// data over the network stack.
//
// The starting point is the creation and configuration of a stack. A stack can
// be created by calling the New() function of the tcpip/stack/stack package;
// configuring a stack involves creating NICs (via calls to Stack.CreateNIC()),
// adding network addresses (via calls to Stack.AddAddress()), and
// setting a route table (via a call to Stack.SetRouteTable()).
//
// Once a stack is configured, endpoints can be created by calling
// Stack.NewEndpoint(). Such endpoints can be used to send/receive data, connect
// to peers, listen for connections, accept connections, etc., depending on the
// transport protocol selected.
package tcpip

import (
	"errors"
	"fmt"

	"github.com/google/netstack/tcpip/buffer"
	"github.com/google/netstack/waiter"
)

// Error represents an error in the netstack error space. Using a special type
// ensures that errors outside of this space are not accidentally introduced.
type Error struct {
	string
}

// String implements fmt.Stringer.String.
func (e *Error) String() string {
	return e.string
}

// Errors that can be returned by the network stack.
var (
	ErrUnknownProtocol       = &Error{"unknown protocol"}
	ErrUnknownNICID          = &Error{"unknown nic id"}
	ErrUnknownProtocolOption = &Error{"unknown option for protocol"}
	ErrDuplicateNICID        = &Error{"duplicate nic id"}
	ErrDuplicateAddress      = &Error{"duplicate address"}
	ErrNoRoute               = &Error{"no route"}
	ErrBadLinkEndpoint       = &Error{"bad link layer endpoint"}
	ErrAlreadyBound          = &Error{"endpoint already bound"}
	ErrInvalidEndpointState  = &Error{"endpoint is in invalid state"}
	ErrAlreadyConnecting     = &Error{"endpoint is already connecting"}
	ErrAlreadyConnected      = &Error{"endpoint is already connected"}
	ErrNoPortAvailable       = &Error{"no ports are available"}
	ErrPortInUse             = &Error{"port is in use"}
	ErrBadLocalAddress       = &Error{"bad local address"}
	ErrClosedForSend         = &Error{"endpoint is closed for send"}
	ErrClosedForReceive      = &Error{"endpoint is closed for receive"}
	ErrWouldBlock            = &Error{"operation would block"}
	ErrConnectionRefused     = &Error{"connection was refused"}
	ErrTimeout               = &Error{"operation timed out"}
	ErrAborted               = &Error{"operation aborted"}
	ErrConnectStarted        = &Error{"connection attempt started"}
	ErrDestinationRequired   = &Error{"destination address is required"}
	ErrNotSupported          = &Error{"operation not supported"}
	ErrQueueSizeNotSupported = &Error{"queue size querying not supported"}
	ErrNotConnected          = &Error{"endpoint not connected"}
	ErrConnectionReset       = &Error{"connection reset by peer"}
	ErrConnectionAborted     = &Error{"connection aborted"}
	ErrNoSuchFile            = &Error{"no such file"}
	ErrInvalidOptionValue    = &Error{"invalid option value specified"}
)

// Errors related to Subnet
var (
	errSubnetLengthMismatch = errors.New("subnet length of address and mask differ")
	errSubnetAddressMasked  = errors.New("subnet address has bits set outside the mask")
)

// Address is a byte slice cast as a string that represents the address of a
// network node. Or, in the case of unix endpoints, it may represent a path.
type Address string

// AddressMask is a bitmask for an address.
type AddressMask string

// Subnet is a subnet defined by its address and mask.
type Subnet struct {
	address Address
	mask    AddressMask
}

// NewSubnet creates a new Subnet, checking that the address and mask are the same length.
func NewSubnet(a Address, m AddressMask) (Subnet, error) {
	if len(a) != len(m) {
		return Subnet{}, errSubnetLengthMismatch
	}
	for i := 0; i < len(a); i++ {
		if a[i]&^m[i] != 0 {
			return Subnet{}, errSubnetAddressMasked
		}
	}
	return Subnet{a, m}, nil
}

// Contains returns true iff the address is of the same length and matches the
// subnet address and mask.
func (s *Subnet) Contains(a Address) bool {
	if len(a) != len(s.address) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i]&s.mask[i] != s.address[i] {
			return false
		}
	}
	return true
}

// ID returns the subnet ID.
func (s *Subnet) ID() Address {
	return s.address
}

// Bits returns the number of ones (network bits) and zeros (host bits) in the
// subnet mask.
func (s *Subnet) Bits() (ones int, zeros int) {
	for _, b := range []byte(s.mask) {
		for i := uint(0); i < 8; i++ {
			if b&(1<<i) == 0 {
				zeros++
			} else {
				ones++
			}
		}
	}
	return
}

// Prefix returns the number of bits before the first host bit.
func (s *Subnet) Prefix() int {
	for i, b := range []byte(s.mask) {
		for j := 7; j >= 0; j-- {
			if b&(1<<uint(j)) == 0 {
				return i*8 + 7 - j
			}
		}
	}
	return len(s.mask) * 8
}

// NICID is a number that uniquely identifies a NIC.
type NICID int32

// ShutdownFlags represents flags that can be passed to the Shutdown() method
// of the Endpoint interface.
type ShutdownFlags int

// Values of the flags that can be passed to the Shutdown() method. They can
// be OR'ed together.
const (
	ShutdownRead ShutdownFlags = 1 << iota
	ShutdownWrite
)

// FullAddress represents a full transport node address, as required by the
// Connect() and Bind() methods.
type FullAddress struct {
	// NIC is the ID of the NIC this address refers to.
	//
	// This may not be used by all endpoint types.
	NIC NICID

	// Addr is the network address.
	Addr Address

	// Port is the transport port.
	//
	// This may not be used by all endpoint types.
	Port uint16
}

// Endpoint is the interface implemented by transport protocols (e.g., tcp, udp)
// that exposes functionality like read, write, connect, etc. to users of the
// networking stack.
type Endpoint interface {
	// Close puts the endpoint in a closed state and frees all resources
	// associated with it.
	Close()

	// Read reads data from the endpoint and optionally returns the sender.
	// This method does not block if there is no data pending.
	// It will also either return an error or data, never both.
	Read(*FullAddress) (buffer.View, *Error)

	// Write writes data to the endpoint's peer, or the provided address if
	// one is specified. This method does not block if the data cannot be
	// written.
	//
	// Note that unlike io.Writer.Write, it is not an error for Write to
	// perform a partial write.
	Write(buffer.View, *FullAddress) (uintptr, *Error)

	// Peek reads data without consuming it from the endpoint.
	//
	// This method does not block if there is no data pending.
	Peek([][]byte) (uintptr, *Error)

	// Connect connects the endpoint to its peer. Specifying a NIC is
	// optional.
	//
	// There are three classes of return values:
	//	nil -- the attempt to connect succeeded.
	//	ErrConnectStarted -- the connect attempt started but hasn't
	//		completed yet. In this case, the actual result will
	//		become available via GetSockOpt(ErrorOption) when
	//		the endpoint becomes writable. (This mimics the
	//		connect(2) syscall behavior.)
	//	Anything else -- the attempt to connect failed.
	Connect(address FullAddress) *Error

	// Shutdown closes the read and/or write end of the endpoint connection
	// to its peer.
	Shutdown(flags ShutdownFlags) *Error

	// Listen puts the endpoint in "listen" mode, which allows it to accept
	// new connections.
	Listen(backlog int) *Error

	// Accept returns a new endpoint if a peer has established a connection
	// to an endpoint previously set to listen mode. This method does not
	// block if no new connections are available.
	//
	// The returned Queue is the wait queue for the newly created endpoint.
	Accept() (Endpoint, *waiter.Queue, *Error)

	// Bind binds the endpoint to a specific local address and port.
	// Specifying a NIC is optional.
	//
	// An optional commit function will be executed atomically with respect
	// to binding the endpoint. If this returns an error, the bind will not
	// occur and the error will be propagated back to the caller.
	Bind(address FullAddress, commit func() *Error) *Error

	// GetLocalAddress returns the address to which the endpoint is bound.
	GetLocalAddress() (FullAddress, *Error)

	// GetRemoteAddress returns the address to which the endpoint is
	// connected.
	GetRemoteAddress() (FullAddress, *Error)

	// Readiness returns the current readiness of the endpoint. For example,
	// if waiter.EventIn is set, the endpoint is immediately readable.
	Readiness(mask waiter.EventMask) waiter.EventMask

	// SetSockOpt sets a socket option. opt should be one of the *Option types.
	SetSockOpt(opt interface{}) *Error

	// GetSockOpt gets a socket option. opt should be a pointer to one of the
	// *Option types.
	GetSockOpt(opt interface{}) *Error
}

// ErrorOption is used in GetSockOpt to specify that the last error reported by
// the endpoint should be cleared and returned.
type ErrorOption struct{}

// SendBufferSizeOption is used by SetSockOpt/GetSockOpt to specify the send
// buffer size option.
type SendBufferSizeOption int

// ReceiveBufferSizeOption is used by SetSockOpt/GetSockOpt to specify the
// receive buffer size option.
type ReceiveBufferSizeOption int

// SendQueueSizeOption is used in GetSockOpt to specify that the number of
// unread bytes in the output buffer should be returned.
type SendQueueSizeOption int

// ReceiveQueueSizeOption is used in GetSockOpt to specify that the number of
// unread bytes in the input buffer should be returned.
type ReceiveQueueSizeOption int

// V6OnlyOption is used by SetSockOpt/GetSockOpt to specify whether an IPv6
// socket is to be restricted to sending and receiving IPv6 packets only.
type V6OnlyOption int

// NoDelayOption is used by SetSockOpt/GetSockOpt to specify if data should be
// sent out immediately by the transport protocol. For TCP, it determines if the
// Nagle algorithm is on or off.
type NoDelayOption int

// ReuseAddressOption is used by SetSockOpt/GetSockOpt to specify whether Bind()
// should allow reuse of local address.
type ReuseAddressOption int

// PasscredOption is used by SetSockOpt/GetSockOpt to specify whether
// SCM_CREDENTIALS socket control messages are enabled.
//
// Only supported on Unix sockets.
type PasscredOption int

// Route is a row in the routing table. It specifies through which NIC (and
// gateway) sets of packets should be routed. A row is considered viable if the
// masked target address matches the destination adddress in the row.
type Route struct {
	// Destination is the address that must be matched against the masked
	// target address to check if this row is viable.
	Destination Address

	// Mask specifies which bits of the Destination and the target address
	// must match for this row to be viable.
	Mask Address

	// Gateway is the gateway to be used if this row is viable.
	Gateway Address

	// NIC is the id of the nic to be used if this row is viable.
	NIC NICID
}

// Match determines if r is viable for the given destination address.
func (r *Route) Match(addr Address) bool {
	if len(addr) != len(r.Destination) {
		return false
	}

	for i := 0; i < len(r.Destination); i++ {
		if (addr[i] & r.Mask[i]) != r.Destination[i] {
			return false
		}
	}

	return true
}

// LinkEndpointID represents a data link layer endpoint.
type LinkEndpointID uint64

// TransportProtocolNumber is the number of a transport protocol.
type TransportProtocolNumber uint32

// NetworkProtocolNumber is the number of a network protocol.
type NetworkProtocolNumber uint32

// Stats holds statistics about the networking stack.
type Stats struct {
	// UnkownProtocolRcvdPackets is the number of packets received by the
	// stack that were for an unknown or unsupported protocol.
	UnknownProtocolRcvdPackets uint64

	// UnknownNetworkEndpointRcvdPackets is the number of packets received
	// by the stack that were for a supported network protocol, but whose
	// destination address didn't having a matching endpoint.
	UnknownNetworkEndpointRcvdPackets uint64

	// MalformedRcvPackets is the number of packets received by the stack
	// that were deemed malformed.
	MalformedRcvdPackets uint64

	// DroppedPackets is the number of packets dropped due to full queues.
	DroppedPackets uint64
}

// String implements the fmt.Stringer interface.
func (a Address) String() string {
	switch len(a) {
	case 4:
		return fmt.Sprintf("%d.%d.%d.%d", int(a[0]), int(a[1]), int(a[2]), int(a[3]))
	default:
		return fmt.Sprintf("%x", []byte(a))
	}
}

// To4 converts the IPv4 address to a 4-byte representation.
// If the address is not an IPv4 address, To4 returns "".
func (a Address) To4() Address {
	const (
		ipv4len = 4
		ipv6len = 16
	)
	if len(a) == ipv4len {
		return a
	}
	if len(a) == ipv6len &&
		isZeros(a[0:10]) &&
		a[10] == 0xff &&
		a[11] == 0xff {
		return a[12:16]
	}
	return ""
}

// isZeros reports whether a is all zeros.
func isZeros(a Address) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != 0 {
			return false
		}
	}
	return true
}

// LinkAddress is a byte slice cast as a string that represents a link address.
// It is typically a 6-byte MAC address.
type LinkAddress string

// String implements the fmt.Stringer interface.
func (a LinkAddress) String() string {
	switch len(a) {
	case 6:
		return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", a[0], a[1], a[2], a[3], a[4], a[5])
	default:
		return fmt.Sprintf("%x", []byte(a))
	}
}
