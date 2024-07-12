// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import "net"

type Trace struct {
	Dest     net.IP
	destPort uint16
	src      net.IP
	//srcPort     uint16
	PortOffset   int32
	MaxHops      int
	SendChan     chan<- *Probe
	ReceiveChan  chan<- *Probe
	exitChan     chan<- bool
	debug        bool
	TracesPerHop int
	PacketRate   int
}

func NewTrace(proto string, dAddr [4]byte, sAddr *net.UDPAddr, cc coms, debug bool) *Trace {
	var destPort uint16

	switch proto {
	case "udp4":
		destPort = 33434
	case "tcp4":
		destPort = 443
		// ICMP does not require a port, duh!
		// case "icmp4":
	}

	ret := &Trace{
		Dest:         net.IPv4(dAddr[0], dAddr[1], dAddr[2], dAddr[3]),
		destPort:     destPort,
		src:          sAddr.IP,
		PortOffset:   0,
		MaxHops:      DEFNUMHOPS, // for IPv4 for now
		SendChan:     cc.sendChan,
		ReceiveChan:  cc.recvChan,
		exitChan:     cc.exitChan,
		debug:        debug,
		TracesPerHop: DEFNUMTRACES,
		PacketRate:   1,
	}
	return ret
}
