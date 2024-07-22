// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import (
	"fmt"
	"log"
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

	var addr string
	for _, a := range addrs {
		if strings.Contains(a, ":") && strings.Contains(proto, "6") {
			addr = a
			break
		} else if strings.Contains(a, ".") && strings.Contains(proto, "4") {
			addr = a
			break
		}
	}

	if len(addr) < 1 {
		return nil, fmt.Errorf("no valid ip address for proto: %s", proto)
	}

	ipAddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return nil, err
	}

	return ipAddr.IP, nil
}

func srcAddr(proto string) (net.IP, error) {
	var sAddr net.Addr
	if strings.Contains(proto, "6") {
		conn, err := net.Dial("udp6", "[2001:4860:4860::8844]:53")
		if err != nil {
			log.Fatal(err)
		}
		sAddr = conn.LocalAddr().(*net.UDPAddr)
		conn.Close()
	} else {
		conn, err := net.Dial("udp", "8.8.8.8:53")
		if err != nil {
			log.Fatal(err)
		}
		sAddr = conn.LocalAddr().(*net.UDPAddr)
		conn.Close()
	}
	return sAddr.(*net.UDPAddr).IP, nil
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
