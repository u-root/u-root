// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import "net"

type Trace struct {
	DestIP net.IP

	// See --port in traceroute(8):
	// For UDP tracing, specifies the destination port base traceroute will use (the destination port number
	// will be incremented by each probe).
	// For ICMP tracing, specifies the initial ICMP sequence value (incremented by each probe too).
	// For TCP and others specifies just the (constant) destination port to connect.
	DestPort uint16

	SrcIP        net.IP
	PortOffset   int32
	MaxHops      int
	SendChan     chan<- *Probe
	ReceiveChan  chan<- *Probe
	TracesPerHop int
	PacketRate   int
	ICMPSeqStart uint16
}

func NewTrace(proto string, dAddr net.IP, sAddr net.IP, cc Coms, f *Flags) *Trace {
	var (
		ret               *Trace
		destAddr, srcAddr net.IP
		dPort             uint16
	)

	switch proto {
	case "udp4":
		destAddr = dAddr.To4()
		srcAddr = sAddr.To4()
		dPort = UDPDEFPORT
	case "udp6":
		destAddr = dAddr.To16()
		srcAddr = sAddr.To16()
		dPort = UDPDEFPORT
	case "tcp4":
		destAddr = dAddr.To4()
		srcAddr = sAddr.To4()
		dPort = TCPDEFPORT
	case "tcp6":
		destAddr = dAddr.To16()
		srcAddr = sAddr.To16()
		dPort = TCPDEFPORT
	case "icmp4":
		destAddr = dAddr.To4()
		srcAddr = sAddr.To4()
		dPort = 0
	case "icmp6":
		destAddr = dAddr.To16()
		srcAddr = sAddr.To16()
		dPort = 0
	}

	// update only when the user specifies a port and it is not already 0 (icmp)
	if f.DestPortSeq != 0 {
		dPort = uint16(f.DestPortSeq)
	}

	ret = &Trace{
		DestIP:       destAddr,
		DestPort:     dPort,
		SrcIP:        srcAddr,
		PortOffset:   0,
		MaxHops:      DEFNUMHOPS,
		SendChan:     cc.SendChan,
		ReceiveChan:  cc.RecvChan,
		TracesPerHop: DEFNUMTRACES,
		PacketRate:   1,
	}

	return ret
}
