// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import "net"

type Trace struct {
	destIP   net.IP
	destPort uint16
	srcIP    net.IP
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

func NewTrace(proto string, dAddr net.IP, sAddr net.IP, cc coms, f *Flags) *Trace {
	var ret *Trace

	switch proto {
	case "udp4":
		ret = &Trace{
			destIP:       dAddr.To4(),
			destPort:     33434,
			srcIP:        sAddr.To4(),
			PortOffset:   0,
			MaxHops:      DEFNUMHOPS, // for IPv4 for now
			SendChan:     cc.sendChan,
			ReceiveChan:  cc.recvChan,
			exitChan:     cc.exitChan,
			debug:        f.Debug,
			TracesPerHop: DEFNUMTRACES,
			PacketRate:   1,
		}
	case "tcp4":
		ret = &Trace{
			destIP:       dAddr.To4(),
			destPort:     443,
			srcIP:        sAddr.To4(),
			PortOffset:   0,
			MaxHops:      DEFNUMHOPS, // for IPv4 for now
			SendChan:     cc.sendChan,
			ReceiveChan:  cc.recvChan,
			exitChan:     cc.exitChan,
			debug:        f.Debug,
			TracesPerHop: DEFNUMTRACES,
			PacketRate:   1,
		}
	case "icmp4":
		ret = &Trace{
			destIP:       dAddr.To4(),
			destPort:     0,
			srcIP:        sAddr.To4(),
			PortOffset:   0,
			MaxHops:      DEFNUMHOPS, // for IPv4 for now
			SendChan:     cc.sendChan,
			ReceiveChan:  cc.recvChan,
			exitChan:     cc.exitChan,
			debug:        f.Debug,
			TracesPerHop: DEFNUMTRACES,
			PacketRate:   1,
		}
	case "icmp6":
		ret = &Trace{
			destIP:       dAddr,
			destPort:     0,
			srcIP:        sAddr,
			PortOffset:   0,
			MaxHops:      DEFNUMHOPS, // for IPv4 for now
			SendChan:     cc.sendChan,
			ReceiveChan:  cc.recvChan,
			exitChan:     cc.exitChan,
			debug:        f.Debug,
			TracesPerHop: DEFNUMTRACES,
			PacketRate:   1,
		}
	}

	return ret
}
