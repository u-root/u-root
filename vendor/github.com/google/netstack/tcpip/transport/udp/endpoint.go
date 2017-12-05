// Copyright 2016 The Netstack Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package udp

import (
	"sync"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/buffer"
	"github.com/google/netstack/tcpip/header"
	"github.com/google/netstack/tcpip/stack"
	"github.com/google/netstack/waiter"
)

type udpPacket struct {
	udpPacketEntry
	senderAddress tcpip.FullAddress
	data          buffer.VectorisedView
	// views is used as buffer for data when its length is large
	// enough to store a VectorisedView.
	views [8]buffer.View
}

type endpointState int

const (
	stateInitial endpointState = iota
	stateBound
	stateConnected
	stateClosed
)

// endpoint represents a UDP endpoint. This struct serves as the interface
// between users of the endpoint and the protocol implementation; it is legal to
// have concurrent goroutines make calls into the endpoint, they are properly
// synchronized.
type endpoint struct {
	// The following fields are initialized at creation time and do not
	// change throughout the lifetime of the endpoint.
	stack       *stack.Stack
	netProto    tcpip.NetworkProtocolNumber
	waiterQueue *waiter.Queue

	// The following fields are used to manage the receive queue, and are
	// protected by rcvMu.
	rcvMu         sync.Mutex
	rcvReady      bool
	rcvList       udpPacketList
	rcvBufSizeMax int
	rcvBufSize    int
	rcvClosed     bool

	// The following fields are protected by the mu mutex.
	mu         sync.RWMutex
	sndBufSize int
	id         stack.TransportEndpointID
	state      endpointState
	bindNICID  tcpip.NICID
	bindAddr   tcpip.Address
	regNICID   tcpip.NICID
	route      stack.Route
	dstPort    uint16
	v6only     bool

	// effectiveNetProtos contains the network protocols actually in use. In
	// most cases it will only contain "netProto", but in cases like IPv6
	// endpoints with v6only set to false, this could include multiple
	// protocols (e.g., IPv6 and IPv4) or a single different protocol (e.g.,
	// IPv4 when IPv6 endpoint is bound or connected to an IPv4 mapped
	// address).
	effectiveNetProtos []tcpip.NetworkProtocolNumber
}

func newEndpoint(stack *stack.Stack, netProto tcpip.NetworkProtocolNumber, waiterQueue *waiter.Queue) *endpoint {
	// TODO: Use the send buffer size initialized here.
	return &endpoint{
		stack:         stack,
		netProto:      netProto,
		waiterQueue:   waiterQueue,
		rcvBufSizeMax: 32 * 1024,
		sndBufSize:    32 * 1024,
	}
}

// NewConnectedEndpoint creates a new endpoint in the connected state using the
// provided route.
func NewConnectedEndpoint(stack *stack.Stack, r *stack.Route, id stack.TransportEndpointID, waiterQueue *waiter.Queue) (tcpip.Endpoint, *tcpip.Error) {
	ep := newEndpoint(stack, r.NetProto, waiterQueue)

	// Register new endpoint so that packets are routed to it.
	if err := stack.RegisterTransportEndpoint(r.NICID(), []tcpip.NetworkProtocolNumber{r.NetProto}, ProtocolNumber, id, ep); err != nil {
		ep.Close()
		return nil, err
	}

	ep.id = id
	ep.route = r.Clone()
	ep.dstPort = id.RemotePort
	ep.regNICID = r.NICID()

	ep.state = stateConnected

	return ep, nil
}

// Close puts the endpoint in a closed state and frees all resources
// associated with it.
func (e *endpoint) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch e.state {
	case stateBound, stateConnected:
		e.stack.UnregisterTransportEndpoint(e.regNICID, e.effectiveNetProtos, ProtocolNumber, e.id)
	}

	// Close the receive list and drain it.
	e.rcvMu.Lock()
	e.rcvClosed = true
	e.rcvBufSize = 0
	for !e.rcvList.Empty() {
		p := e.rcvList.Front()
		e.rcvList.Remove(p)
	}
	e.rcvMu.Unlock()

	e.route.Release()

	// Update the state.
	e.state = stateClosed
}

// Read reads data from the endpoint. This method does not block if
// there is no data pending.
func (e *endpoint) Read(addr *tcpip.FullAddress) (buffer.View, *tcpip.Error) {
	e.rcvMu.Lock()

	if e.rcvList.Empty() {
		err := tcpip.ErrWouldBlock
		if e.rcvClosed {
			err = tcpip.ErrClosedForReceive
		}
		e.rcvMu.Unlock()
		return buffer.View{}, err
	}

	p := e.rcvList.Front()
	e.rcvList.Remove(p)
	e.rcvBufSize -= p.data.Size()

	e.rcvMu.Unlock()

	if addr != nil {
		*addr = p.senderAddress
	}

	return p.data.ToView(), nil
}

// prepareForWrite prepares the endpoint for sending data. In particular, it
// binds it if it's still in the initial state. To do so, it must first
// reacquire the mutex in exclusive mode.
//
// Returns true for retry if preparation should be retried.
func (e *endpoint) prepareForWrite(to *tcpip.FullAddress) (retry bool, err *tcpip.Error) {
	switch e.state {
	case stateInitial:
	case stateConnected:
		return false, nil

	case stateBound:
		if to == nil {
			return false, tcpip.ErrDestinationRequired
		}
		return false, nil
	default:
		return false, tcpip.ErrInvalidEndpointState
	}

	e.mu.RUnlock()
	defer e.mu.RLock()

	e.mu.Lock()
	defer e.mu.Unlock()

	// The state changed when we released the shared locked and re-acquired
	// it in exclusive mode. Try again.
	if e.state != stateInitial {
		return true, nil
	}

	// The state is still 'initial', so try to bind the endpoint.
	if err := e.bindLocked(tcpip.FullAddress{}, nil); err != nil {
		return false, err
	}

	return true, nil
}

// Write writes data to the endpoint's peer. This method does not block
// if the data cannot be written.
func (e *endpoint) Write(v buffer.View, to *tcpip.FullAddress) (uintptr, *tcpip.Error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Prepare for write.
	for {
		retry, err := e.prepareForWrite(to)
		if err != nil {
			return 0, err
		}

		if !retry {
			break
		}
	}

	route := &e.route
	dstPort := e.dstPort
	if to != nil {
		// Reject destination address if it goes through a different
		// NIC than the endpoint was bound to.
		nicid := to.NIC
		if e.bindNICID != 0 {
			if nicid != 0 && nicid != e.bindNICID {
				return 0, tcpip.ErrNoRoute
			}

			nicid = e.bindNICID
		}

		toCopy := *to
		to = &toCopy
		netProto, err := e.checkV4Mapped(to, true)
		if err != nil {
			return 0, err
		}

		// Find the enpoint.
		r, err := e.stack.FindRoute(nicid, e.bindAddr, to.Addr, netProto)
		if err != nil {
			return 0, err
		}
		defer r.Release()

		route = &r
		dstPort = to.Port
	}

	sendUDP(route, v, e.id.LocalPort, dstPort)
	return uintptr(len(v)), nil
}

// Peek only returns data from a single datagram, so do nothing here.
func (e *endpoint) Peek([][]byte) (uintptr, *tcpip.Error) {
	return 0, nil
}

// SetSockOpt sets a socket option. Currently not supported.
func (e *endpoint) SetSockOpt(opt interface{}) *tcpip.Error {
	// TODO: Actually implement this.
	switch v := opt.(type) {
	case tcpip.V6OnlyOption:
		// We only recognize this option on v6 endpoints.
		if e.netProto != header.IPv6ProtocolNumber {
			return tcpip.ErrInvalidEndpointState
		}

		e.mu.Lock()
		defer e.mu.Unlock()

		// We only allow this to be set when we're in the initial state.
		if e.state != stateInitial {
			return tcpip.ErrInvalidEndpointState
		}

		e.v6only = v != 0
	}
	return nil
}

// GetSockOpt implements tcpip.Endpoint.GetSockOpt.
func (e *endpoint) GetSockOpt(opt interface{}) *tcpip.Error {
	switch o := opt.(type) {
	case tcpip.ErrorOption:
		return nil

	case *tcpip.SendBufferSizeOption:
		e.mu.Lock()
		*o = tcpip.SendBufferSizeOption(e.sndBufSize)
		e.mu.Unlock()
		return nil

	case *tcpip.ReceiveBufferSizeOption:
		e.rcvMu.Lock()
		*o = tcpip.ReceiveBufferSizeOption(e.rcvBufSizeMax)
		e.rcvMu.Unlock()
		return nil

	case *tcpip.V6OnlyOption:
		// We only recognize this option on v6 endpoints.
		if e.netProto != header.IPv6ProtocolNumber {
			return tcpip.ErrUnknownProtocolOption
		}

		e.mu.Lock()
		v := e.v6only
		e.mu.Unlock()

		*o = 0
		if v {
			*o = 1
		}
		return nil

	case *tcpip.ReceiveQueueSizeOption:
		e.rcvMu.Lock()
		if e.rcvList.Empty() {
			*o = 0
		} else {
			p := e.rcvList.Front()
			*o = tcpip.ReceiveQueueSizeOption(p.data.Size())
		}
		e.rcvMu.Unlock()
		return nil
	}

	return tcpip.ErrUnknownProtocolOption
}

// sendUDP sends a UDP segment via the provided network endpoint and under the
// provided identity.
func sendUDP(r *stack.Route, data buffer.View, localPort, remotePort uint16) *tcpip.Error {
	// Allocate a buffer for the UDP header.
	hdr := buffer.NewPrependable(header.UDPMinimumSize + int(r.MaxHeaderLength()))

	// Initialize the header.
	udp := header.UDP(hdr.Prepend(header.UDPMinimumSize))

	length := uint16(hdr.UsedLength())
	xsum := r.PseudoHeaderChecksum(ProtocolNumber)
	if data != nil {
		length += uint16(len(data))
		xsum = header.Checksum(data, xsum)
	}

	udp.Encode(&header.UDPFields{
		SrcPort: localPort,
		DstPort: remotePort,
		Length:  length,
	})

	udp.SetChecksum(^udp.CalculateChecksum(xsum, length))

	return r.WritePacket(&hdr, data, ProtocolNumber)
}

func (e *endpoint) checkV4Mapped(addr *tcpip.FullAddress, allowMismatch bool) (tcpip.NetworkProtocolNumber, *tcpip.Error) {
	netProto := e.netProto
	if header.IsV4MappedAddress(addr.Addr) {
		// Fail if using a v4 mapped address on a v6only endpoint.
		if e.v6only {
			return 0, tcpip.ErrNoRoute
		}

		netProto = header.IPv4ProtocolNumber
		addr.Addr = addr.Addr[header.IPv6AddressSize-header.IPv4AddressSize:]
		if addr.Addr == "\x00\x00\x00\x00" {
			addr.Addr = ""
		}
	}

	// Fail if we're bound to an address length different from the one we're
	// checking.
	if l := len(e.id.LocalAddress); !allowMismatch && l != 0 && l != len(addr.Addr) {
		return 0, tcpip.ErrInvalidEndpointState
	}

	return netProto, nil
}

// Connect connects the endpoint to its peer. Specifying a NIC is optional.
func (e *endpoint) Connect(addr tcpip.FullAddress) *tcpip.Error {
	if addr.Port == 0 {
		// We don't support connecting to port zero.
		return tcpip.ErrInvalidEndpointState
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	nicid := addr.NIC
	localPort := uint16(0)
	switch e.state {
	case stateInitial:
	case stateBound, stateConnected:
		localPort = e.id.LocalPort
		if e.bindNICID == 0 {
			break
		}

		if nicid != 0 && nicid != e.bindNICID {
			return tcpip.ErrInvalidEndpointState
		}

		nicid = e.bindNICID
	default:
		return tcpip.ErrInvalidEndpointState
	}

	netProto, err := e.checkV4Mapped(&addr, false)
	if err != nil {
		return err
	}

	// Find a route to the desired destination.
	r, err := e.stack.FindRoute(nicid, e.bindAddr, addr.Addr, netProto)
	if err != nil {
		return err
	}
	defer r.Release()

	id := stack.TransportEndpointID{
		LocalAddress:  r.LocalAddress,
		LocalPort:     localPort,
		RemotePort:    addr.Port,
		RemoteAddress: addr.Addr,
	}

	// Even if we're connected, this endpoint can still be used to send
	// packets on a different network protocol, so we register both even if
	// v6only is set to false and this is an ipv6 endpoint.
	netProtos := []tcpip.NetworkProtocolNumber{netProto}
	if e.netProto == header.IPv6ProtocolNumber && !e.v6only {
		netProtos = []tcpip.NetworkProtocolNumber{
			header.IPv4ProtocolNumber,
			header.IPv6ProtocolNumber,
		}
	}

	id, err = e.registerWithStack(nicid, netProtos, id)
	if err != nil {
		return err
	}

	// Remove the old registration.
	if e.id.LocalPort != 0 {
		e.stack.UnregisterTransportEndpoint(e.regNICID, e.effectiveNetProtos, ProtocolNumber, e.id)
	}

	e.id = id
	e.route = r.Clone()
	e.dstPort = addr.Port
	e.regNICID = nicid
	e.effectiveNetProtos = netProtos

	e.state = stateConnected

	e.rcvMu.Lock()
	e.rcvReady = true
	e.rcvMu.Unlock()

	return nil
}

// ConnectEndpoint is not supported.
func (*endpoint) ConnectEndpoint(tcpip.Endpoint) *tcpip.Error {
	return tcpip.ErrInvalidEndpointState
}

// Shutdown closes the read and/or write end of the endpoint connection
// to its peer.
func (e *endpoint) Shutdown(flags tcpip.ShutdownFlags) *tcpip.Error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.state != stateConnected {
		return tcpip.ErrNotConnected
	}

	if flags&tcpip.ShutdownRead != 0 {
		e.rcvMu.Lock()
		wasClosed := e.rcvClosed
		e.rcvClosed = true
		e.rcvMu.Unlock()

		if !wasClosed {
			e.waiterQueue.Notify(waiter.EventIn)
		}
	}

	return nil
}

// Listen is not supported by UDP, it just fails.
func (*endpoint) Listen(int) *tcpip.Error {
	return tcpip.ErrNotSupported
}

// Accept is not supported by UDP, it just fails.
func (*endpoint) Accept() (tcpip.Endpoint, *waiter.Queue, *tcpip.Error) {
	return nil, nil, tcpip.ErrNotSupported
}

func (e *endpoint) registerWithStack(nicid tcpip.NICID, netProtos []tcpip.NetworkProtocolNumber, id stack.TransportEndpointID) (stack.TransportEndpointID, *tcpip.Error) {
	if id.LocalPort != 0 {
		// The endpoint already has a local port, just attempt to
		// register it.
		err := e.stack.RegisterTransportEndpoint(nicid, netProtos, ProtocolNumber, id, e)
		return id, err
	}

	// We need to find a port for the endpoint.
	_, err := e.stack.PickEphemeralPort(func(p uint16) (bool, *tcpip.Error) {
		id.LocalPort = p
		err := e.stack.RegisterTransportEndpoint(nicid, netProtos, ProtocolNumber, id, e)
		switch err {
		case nil:
			return true, nil
		case tcpip.ErrPortInUse:
			return false, nil
		default:
			return false, err
		}
	})

	return id, err
}

func (e *endpoint) bindLocked(addr tcpip.FullAddress, commit func() *tcpip.Error) *tcpip.Error {
	// Don't allow binding once endpoint is not in the initial state
	// anymore.
	if e.state != stateInitial {
		return tcpip.ErrInvalidEndpointState
	}

	netProto, err := e.checkV4Mapped(&addr, false)
	if err != nil {
		return err
	}

	// Expand netProtos to include v4 and v6 if the caller is binding to a
	// wildcard (empty) address, and this is an IPv6 endpoint with v6only
	// set to false.
	netProtos := []tcpip.NetworkProtocolNumber{netProto}
	if netProto == header.IPv6ProtocolNumber && !e.v6only && addr.Addr == "" {
		netProtos = []tcpip.NetworkProtocolNumber{
			header.IPv6ProtocolNumber,
			header.IPv4ProtocolNumber,
		}
	}

	if len(addr.Addr) != 0 {
		// A local address was specified, verify that it's valid.
		if e.stack.CheckLocalAddress(addr.NIC, addr.Addr) == 0 {
			return tcpip.ErrBadLocalAddress
		}
	}

	id := stack.TransportEndpointID{
		LocalPort:    addr.Port,
		LocalAddress: addr.Addr,
	}
	id, err = e.registerWithStack(addr.NIC, netProtos, id)
	if err != nil {
		return err
	}
	if commit != nil {
		if err := commit(); err != nil {
			// Unregister, the commit failed.
			e.stack.UnregisterTransportEndpoint(addr.NIC, netProtos, ProtocolNumber, id)
			return err
		}
	}

	e.id = id
	e.regNICID = addr.NIC
	e.effectiveNetProtos = netProtos

	// Mark endpoint as bound.
	e.state = stateBound

	e.rcvMu.Lock()
	e.rcvReady = true
	e.rcvMu.Unlock()

	return nil
}

// Bind binds the endpoint to a specific local address and port.
// Specifying a NIC is optional.
func (e *endpoint) Bind(addr tcpip.FullAddress, commit func() *tcpip.Error) *tcpip.Error {
	e.mu.Lock()
	defer e.mu.Unlock()

	err := e.bindLocked(addr, commit)
	if err != nil {
		return err
	}

	e.bindNICID = addr.NIC
	e.bindAddr = addr.Addr

	return nil
}

// GetLocalAddress returns the address to which the endpoint is bound.
func (e *endpoint) GetLocalAddress() (tcpip.FullAddress, *tcpip.Error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return tcpip.FullAddress{
		NIC:  e.regNICID,
		Addr: e.id.LocalAddress,
		Port: e.id.LocalPort,
	}, nil
}

// GetRemoteAddress returns the address to which the endpoint is connected.
func (e *endpoint) GetRemoteAddress() (tcpip.FullAddress, *tcpip.Error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.state != stateConnected {
		return tcpip.FullAddress{}, tcpip.ErrNotConnected
	}

	return tcpip.FullAddress{
		NIC:  e.regNICID,
		Addr: e.id.RemoteAddress,
		Port: e.id.RemotePort,
	}, nil
}

// Readiness returns the current readiness of the endpoint. For example, if
// waiter.EventIn is set, the endpoint is immediately readable.
func (e *endpoint) Readiness(mask waiter.EventMask) waiter.EventMask {
	// The endpoint is always writable.
	result := waiter.EventOut & mask

	// Determine if the endpoint is readable if requested.
	if (mask & waiter.EventIn) != 0 {
		e.rcvMu.Lock()
		if !e.rcvList.Empty() || e.rcvClosed {
			result |= waiter.EventIn
		}
		e.rcvMu.Unlock()
	}

	return result
}

// HandlePacket is called by the stack when new packets arrive to this transport
// endpoint.
func (e *endpoint) HandlePacket(r *stack.Route, id stack.TransportEndpointID, vv *buffer.VectorisedView) {
	// Get the header then trim it from the view.
	hdr := header.UDP(vv.First())
	if int(hdr.Length()) > vv.Size() {
		// Malformed packet.
		return
	}

	vv.TrimFront(header.UDPMinimumSize)

	e.rcvMu.Lock()

	// Drop the packet if our buffer is currently full.
	if !e.rcvReady || e.rcvClosed || e.rcvBufSize >= e.rcvBufSizeMax {
		e.rcvMu.Unlock()
		return
	}

	wasEmpty := e.rcvBufSize == 0

	// Push new packet into receive list and increment the buffer size.
	pkt := &udpPacket{
		senderAddress: tcpip.FullAddress{
			NIC:  r.NICID(),
			Addr: id.RemoteAddress,
			Port: hdr.SourcePort(),
		},
	}
	pkt.data = vv.Clone(pkt.views[:])
	e.rcvList.PushBack(pkt)
	e.rcvBufSize += vv.Size()

	e.rcvMu.Unlock()

	// Notify any waiters that there's data to be read now.
	if wasEmpty {
		e.waiterQueue.Notify(waiter.EventIn)
	}
}
