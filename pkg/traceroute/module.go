// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import "net"

type Module interface {
	SendTraces()
	ReceiveTraces()
}

const (
	DEFNUMHOPS   = 30
	DEFNUMTRACES = 3
)

func NewModule(mod string, dAddr [4]byte, sAddr *net.UDPAddr, cc coms, debug bool) Module {
	var ret Module

	switch mod {
	case "udp":
		ret = &UDP4Trace{
			Dest:         net.IPv4(dAddr[0], dAddr[1], dAddr[2], dAddr[3]),
			destPort:     33434,
			src:          sAddr.IP,
			PortOffset:   0,
			MaxHops:      DEFNUMHOPS, // for IPv4 for now
			SendChan:     cc.sendChan,
			ReceiveChan:  cc.recvChan,
			exitChan:     cc.exitChan,
			debug:        debug,
			TracesPerHop: DEFNUMTRACES,
		}
		return ret
	}
	return nil
}
