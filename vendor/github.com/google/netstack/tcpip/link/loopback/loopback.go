// Copyright 2016 The Netstack Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package loopback provides the implemention of loopback data-link layer
// endpoints. Such endpoints just turn outbound packets into inbound ones.
//
// Loopback endpoints can be used in the networking stack by calling New() to
// create a new endpoint, and then passing it as an argument to
// Stack.CreateNIC().
package loopback

import (
	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/buffer"
	"github.com/google/netstack/tcpip/stack"
)

type endpoint struct {
	dispatcher stack.NetworkDispatcher
}

// New creates a new loopback endpoint. This link-layer endpoint just turns
// outbound packets into inbound packets.
func New() tcpip.LinkEndpointID {
	return stack.RegisterLinkEndpoint(&endpoint{})
}

// Attach implements stack.LinkEndpoint.Attach. It just saves the stack network-
// layer dispatcher for later use when packets need to be dispatched.
func (e *endpoint) Attach(dispatcher stack.NetworkDispatcher) {
	e.dispatcher = dispatcher
}

// MTU implements stack.LinkEndpoint.MTU. It returns a constant that matches the
// linux loopback interface.
func (*endpoint) MTU() uint32 {
	return 65536
}

// MaxHeaderLength implements stack.LinkEndpoint.MaxHeaderLength. Given that the
// loopback interface doesn't have a header, it just returns 0.
func (*endpoint) MaxHeaderLength() uint16 {
	return 0
}

// LinkAddress returns the link address of this endpoint.
func (*endpoint) LinkAddress() tcpip.LinkAddress {
	return ""
}

// WritePacket implements stack.LinkEndpoint.WritePacket. It delivers outbound
// packets to the network-layer dispatcher.
func (e *endpoint) WritePacket(_ *stack.Route, hdr *buffer.Prependable, payload buffer.View, protocol tcpip.NetworkProtocolNumber) *tcpip.Error {
	if len(payload) == 0 {
		// We don't have a payload, so just use the buffer from the
		// header as the full packet.
		v := hdr.View()
		vv := v.ToVectorisedView([1]buffer.View{})
		e.dispatcher.DeliverNetworkPacket(e, "", protocol, &vv)
	} else {
		views := []buffer.View{hdr.View(), payload}
		vv := buffer.NewVectorisedView(len(views[0])+len(views[1]), views)
		e.dispatcher.DeliverNetworkPacket(e, "", protocol, &vv)
	}

	return nil
}
