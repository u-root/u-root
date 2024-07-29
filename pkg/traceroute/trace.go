// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import "net"

type Trace struct {
	DestIP   net.IP
	destPort uint16
	SrcIP    net.IP
	//srcPort     uint16
	PortOffset   int32
	MaxHops      int
	SendChan     chan<- *Probe
	ReceiveChan  chan<- *Probe
	TracesPerHop int
	PacketRate   int
}

func NewTrace(proto string, dAddr net.IP, sAddr net.IP, cc Coms, f *Flags) *Trace {
	var ret *Trace
	var destAddr, srcAddr net.IP
	var dPort uint16

	switch proto {
	case "udp4":
		destAddr = dAddr.To4()
		srcAddr = sAddr.To4()
		dPort = 33434
	case "udp6":
		destAddr = dAddr.To16()
		srcAddr = sAddr.To16()
		dPort = 33434
	case "tcp4":
		destAddr = dAddr.To4()
		srcAddr = sAddr.To4()
		dPort = 443
	case "tcp6":
		destAddr = dAddr.To16()
		srcAddr = sAddr.To16()
		dPort = 443
	case "icmp4":
		destAddr = dAddr.To4()
		srcAddr = sAddr.To4()
		dPort = 0
	case "icmp6":
		destAddr = dAddr.To16()
		srcAddr = sAddr.To16()
		dPort = 0
	}

	ret = &Trace{
		DestIP:       destAddr,
		destPort:     dPort,
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
