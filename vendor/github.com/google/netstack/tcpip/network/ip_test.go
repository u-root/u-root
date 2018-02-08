// Copyright 2016 The Netstack Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ip_test

import (
	"testing"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/buffer"
	"github.com/google/netstack/tcpip/header"
	"github.com/google/netstack/tcpip/network/ipv4"
	"github.com/google/netstack/tcpip/network/ipv6"
	"github.com/google/netstack/tcpip/stack"
)

// testObject implements two interfaces: LinkEndpoint and TransportDispatcher.
// The former is used to pretend that it's a link endpoint so that we can
// inspect packets written by the network endpoints. The latter is used to
// pretend that it's the network stack so that it can inspect incoming packets
// that have been handled by the network endpoints.
//
// Packets are checked by comparing their fields/values against the expected
// values stored in the test object itself.
type testObject struct {
	t        *testing.T
	protocol tcpip.TransportProtocolNumber
	contents []byte
	srcAddr  tcpip.Address
	dstAddr  tcpip.Address
	v4       bool
}

// checkValues verifies that the transport protocol, data contents, src & dst
// addresses of a packet match what's expected. If any field doesn't match, the
// test fails.
func (t *testObject) checkValues(protocol tcpip.TransportProtocolNumber, vv *buffer.VectorisedView, srcAddr, dstAddr tcpip.Address) {
	v := vv.ToView()
	if protocol != t.protocol {
		t.t.Errorf("protocol = %v, want %v", protocol, t.protocol)
	}

	if srcAddr != t.srcAddr {
		t.t.Errorf("srcAddr = %v, want %v", srcAddr, t.srcAddr)
	}

	if dstAddr != t.dstAddr {
		t.t.Errorf("dstAddr = %v, want %v", dstAddr, t.dstAddr)
	}

	if len(v) != len(t.contents) {
		t.t.Fatalf("len(payload) = %v, want %v", len(v), len(t.contents))
	}

	for i := range t.contents {
		if t.contents[i] != v[i] {
			t.t.Fatalf("payload[%v] = %v, want %v", i, v[i], t.contents[i])
		}
	}
}

// DeliverTransportPacket is called by network endpoints after parsing incoming
// packets. This is used by the test object to verify that the results of the
// parsing are expected.
func (t *testObject) DeliverTransportPacket(r *stack.Route, protocol tcpip.TransportProtocolNumber, vv *buffer.VectorisedView) {
	t.checkValues(protocol, vv, r.RemoteAddress, r.LocalAddress)
}

// Attach is only implemented to satisfy the LinkEndpoint interface.
func (*testObject) Attach(stack.NetworkDispatcher) {}

// MTU implements stack.LinkEndpoint.MTU. It just returns a constant that
// matches the linux loopback MTU.
func (*testObject) MTU() uint32 {
	return 65536
}

// MaxHeaderLength is only implemented to satisfy the LinkEndpoint interface.
func (*testObject) MaxHeaderLength() uint16 {
	return 0
}

// LinkAddress returns the link address of this endpoint.
func (*testObject) LinkAddress() tcpip.LinkAddress {
	return ""
}

// WritePacket is called by network endpoints after producing a packet and
// writing it to the link endpoint. This is used by the test object to verify
// that the produced packet is as expected.
func (t *testObject) WritePacket(_ *stack.Route, hdr *buffer.Prependable, payload buffer.View, protocol tcpip.NetworkProtocolNumber) *tcpip.Error {
	var prot tcpip.TransportProtocolNumber
	var srcAddr tcpip.Address
	var dstAddr tcpip.Address

	if t.v4 {
		h := header.IPv4(hdr.UsedBytes())
		prot = tcpip.TransportProtocolNumber(h.Protocol())
		srcAddr = h.SourceAddress()
		dstAddr = h.DestinationAddress()

	} else {
		h := header.IPv6(hdr.UsedBytes())
		prot = tcpip.TransportProtocolNumber(h.NextHeader())
		srcAddr = h.SourceAddress()
		dstAddr = h.DestinationAddress()
	}
	var views [1]buffer.View
	vv := payload.ToVectorisedView(views)
	t.checkValues(prot, &vv, srcAddr, dstAddr)
	return nil
}

func TestIPv4Send(t *testing.T) {
	o := testObject{t: t, v4: true}
	proto := ipv4.NewProtocol()
	ep, err := proto.NewEndpoint(1, "\x0a\x00\x00\x01", nil, nil, &o)
	if err != nil {
		t.Fatalf("NewEndpoint failed: %v", err)
	}

	// Allocate and initialize the payload view.
	payload := buffer.NewView(100)
	for i := 0; i < len(payload); i++ {
		payload[i] = uint8(i)
	}

	// Allocate the header buffer.
	hdr := buffer.NewPrependable(int(ep.MaxHeaderLength()))

	// Issue the write.
	o.protocol = 123
	o.srcAddr = "\x0a\x00\x00\x01"
	o.dstAddr = "\x0a\x00\x00\x02"
	o.contents = payload

	r := stack.Route{
		RemoteAddress: o.dstAddr,
		LocalAddress:  o.srcAddr,
	}
	if err := ep.WritePacket(&r, &hdr, payload, 123); err != nil {
		t.Fatalf("WritePacket failed: %v", err)
	}
}

func TestIPv4Receive(t *testing.T) {
	o := testObject{t: t, v4: true}
	proto := ipv4.NewProtocol()
	ep, err := proto.NewEndpoint(1, "\x0a\x00\x00\x01", nil, &o, nil)
	if err != nil {
		t.Fatalf("NewEndpoint failed: %v", err)
	}

	totalLen := header.IPv4MinimumSize + 30
	view := buffer.NewView(totalLen)
	ip := header.IPv4(view)
	ip.Encode(&header.IPv4Fields{
		IHL:         header.IPv4MinimumSize,
		TotalLength: uint16(totalLen),
		TTL:         20,
		Protocol:    10,
		SrcAddr:     "\x0a\x00\x00\x02",
		DstAddr:     "\x0a\x00\x00\x01",
	})

	// Make payload be non-zero.
	for i := header.IPv4MinimumSize; i < totalLen; i++ {
		view[i] = uint8(i)
	}

	// Give packet to ipv4 endpoint, dispatcher will validate that it's ok.
	o.protocol = 10
	o.srcAddr = "\x0a\x00\x00\x02"
	o.dstAddr = "\x0a\x00\x00\x01"
	o.contents = view[header.IPv4MinimumSize:totalLen]

	r := stack.Route{
		LocalAddress:  o.dstAddr,
		RemoteAddress: o.srcAddr,
	}
	var views [1]buffer.View
	vv := view.ToVectorisedView(views)
	ep.HandlePacket(&r, &vv)
}

func TestIPv4FragmentationReceive(t *testing.T) {
	o := testObject{t: t, v4: true}
	proto := ipv4.NewProtocol()
	ep, err := proto.NewEndpoint(1, "\x0a\x00\x00\x01", nil, &o, nil)
	if err != nil {
		t.Fatalf("NewEndpoint failed: %v", err)
	}

	totalLen := header.IPv4MinimumSize + 24

	frag1 := buffer.NewView(totalLen)
	ip1 := header.IPv4(frag1)
	ip1.Encode(&header.IPv4Fields{
		IHL:            header.IPv4MinimumSize,
		TotalLength:    uint16(totalLen),
		TTL:            20,
		Protocol:       10,
		FragmentOffset: 0,
		Flags:          header.IPv4FlagMoreFragments,
		SrcAddr:        "\x0a\x00\x00\x02",
		DstAddr:        "\x0a\x00\x00\x01",
	})
	// Make payload be non-zero.
	for i := header.IPv4MinimumSize; i < totalLen; i++ {
		frag1[i] = uint8(i)
	}

	frag2 := buffer.NewView(totalLen)
	ip2 := header.IPv4(frag2)
	ip2.Encode(&header.IPv4Fields{
		IHL:            header.IPv4MinimumSize,
		TotalLength:    uint16(totalLen),
		TTL:            20,
		Protocol:       10,
		FragmentOffset: 24,
		SrcAddr:        "\x0a\x00\x00\x02",
		DstAddr:        "\x0a\x00\x00\x01",
	})
	// Make payload be non-zero.
	for i := header.IPv4MinimumSize; i < totalLen; i++ {
		frag2[i] = uint8(i)
	}

	// Give packet to ipv4 endpoint, dispatcher will validate that it's ok.
	o.protocol = 10
	o.srcAddr = "\x0a\x00\x00\x02"
	o.dstAddr = "\x0a\x00\x00\x01"
	o.contents = append(frag1[header.IPv4MinimumSize:totalLen], frag2[header.IPv4MinimumSize:totalLen]...)

	r := stack.Route{
		LocalAddress:  o.dstAddr,
		RemoteAddress: o.srcAddr,
	}

	// Send first segment.
	var views1 [1]buffer.View
	vv1 := frag1.ToVectorisedView(views1)
	ep.HandlePacket(&r, &vv1)

	// Send second segment.
	var views2 [1]buffer.View
	vv2 := frag2.ToVectorisedView(views2)
	ep.HandlePacket(&r, &vv2)
}

func TestIPv6Send(t *testing.T) {
	o := testObject{t: t}
	proto := ipv6.NewProtocol()
	ep, err := proto.NewEndpoint(1, "\x0a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01", nil, nil, &o)
	if err != nil {
		t.Fatalf("NewEndpoint failed: %v", err)
	}

	// Allocate and initialize the payload view.
	payload := buffer.NewView(100)
	for i := 0; i < len(payload); i++ {
		payload[i] = uint8(i)
	}

	// Allocate the header buffer.
	hdr := buffer.NewPrependable(int(ep.MaxHeaderLength()))

	// Issue the write.
	o.protocol = 123
	o.srcAddr = "\x0a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01"
	o.dstAddr = "\x0a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02"
	o.contents = payload

	r := stack.Route{
		RemoteAddress: o.dstAddr,
		LocalAddress:  o.srcAddr,
	}
	if err := ep.WritePacket(&r, &hdr, payload, 123); err != nil {
		t.Fatalf("WritePacket failed: %v", err)
	}
}

func TestIPv6Receive(t *testing.T) {
	o := testObject{t: t}
	proto := ipv6.NewProtocol()
	ep, err := proto.NewEndpoint(1, "\x0a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01", nil, &o, nil)
	if err != nil {
		t.Fatalf("NewEndpoint failed: %v", err)
	}

	totalLen := header.IPv6MinimumSize + 30
	view := buffer.NewView(totalLen)
	ip := header.IPv6(view)
	ip.Encode(&header.IPv6Fields{
		PayloadLength: uint16(totalLen - header.IPv6MinimumSize),
		NextHeader:    10,
		HopLimit:      20,
		SrcAddr:       "\x0a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02",
		DstAddr:       "\x0a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01",
	})

	// Make payload be non-zero.
	for i := header.IPv6MinimumSize; i < totalLen; i++ {
		view[i] = uint8(i)
	}

	// Give packet to ipv6 endpoint, dispatcher will validate that it's ok.
	o.protocol = 10
	o.srcAddr = "\x0a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02"
	o.dstAddr = "\x0a\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01"
	o.contents = view[header.IPv6MinimumSize:totalLen]

	r := stack.Route{
		LocalAddress:  o.dstAddr,
		RemoteAddress: o.srcAddr,
	}
	var views [1]buffer.View
	vv := view.ToVectorisedView(views)
	ep.HandlePacket(&r, &vv)
}
