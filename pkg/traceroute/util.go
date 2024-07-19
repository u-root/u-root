// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import (
	"net"
	"strings"
)

type coms struct {
	sendChan chan *Probe
	recvChan chan *Probe
	exitChan chan bool
}

// Given a host name convert it to a 4 byte IP address.
func destAddr(dest, proto string) (net.IP, error) {
	addrs, err := net.LookupHost(dest)
	if err != nil {
		return nil, err
	}

	addr := addrs[0]
	if strings.Contains(proto, "6") {
		addr = addrs[1]
	}

	ipAddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return nil, err
	}

	return ipAddr.IP, nil
}

func findDestinationTTL(printMap map[int]*Probe) int {
	icmp := false
	destttl := 1
	var icmpfinalpb *Probe
	for _, pb := range printMap {
		if destttl < pb.ttl {
			destttl = pb.ttl
		}
		if pb.ttl == 0 {
			// ICMP TCPProbe needs to increase return value by one
			icmpfinalpb = pb
			icmp = true
		}
	}

	if icmp {
		destttl++
		newttl := destttl
		icmpfinalpb.ttl = newttl
	}

	return destttl
}

func getProbesByTLL(printMap map[int]*Probe, ttl int) []*Probe {
	pbs := make([]*Probe, 0)
	for _, pb := range printMap {
		if pb.ttl == ttl {
			pbs = append(pbs, pb)
		}
	}
	return pbs
}
