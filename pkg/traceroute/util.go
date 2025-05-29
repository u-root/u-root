// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Coms struct {
	SendChan chan *Probe
	RecvChan chan *Probe
}

// Given a host name convert it to a IP address according to provided ip protocol.
func DestAddr(dest, proto string) (net.IP, error) {
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

func isVirtual(iface net.Interface) bool {
	dev := "/sys/devices/virtual/net/" + iface.Name
	if _, err := os.Stat(dev); err == nil {
		return true
	}
	return false
}

func SrcAddr(proto string) (*net.IP, error) {
	var sAddr net.Addr
	var found bool

	if strings.Contains(proto, "6") {
		// this assumes connection to a well-known DNS server
		conn, err := net.Dial("udp6", "[2001:4860:4860::8844]:53")
		if err != nil {
			ifaces, err := net.Interfaces()
			if err != nil {
				return nil, fmt.Errorf("failed to get interfaces: %v", err)
			}

			for _, i := range ifaces {
				if found {
					break
				}

				addrs, err := i.Addrs()
				if err != nil {
					return nil, err
				}
				for _, addr := range addrs {
					var ip net.IP
					if addr.(*net.IPNet).IP.To4() != nil || isVirtual(i) {
						continue // skip IPv4 addresses
					}

					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
						found = true
					case *net.IPAddr:
						ip = v.IP
						found = true
					}
					if found {
						sAddr = &net.UDPAddr{IP: ip, Port: 0}
						break
					}
				}
			}
			if !found {
				return nil, fmt.Errorf("no suitable IPv6 address found for traceroute")
			}
		} else {
			sAddr = conn.LocalAddr().(*net.UDPAddr)
			conn.Close()
		}
	} else {
		conn, err := net.Dial("udp", "8.8.8.8:53")
		if err != nil {
			return nil, err
		}
		sAddr = conn.LocalAddr().(*net.UDPAddr)
		conn.Close()
	}
	return &sAddr.(*net.UDPAddr).IP, nil
}

func DestTTL(printMap map[int]*Probe) int {
	icmp := false
	destttl := 1
	var icmpFinalPB *Probe
	for _, pb := range printMap {
		if destttl < pb.TTL {
			destttl = pb.TTL
		}
		if pb.TTL == 0 {
			// ICMP TCPProbe needs to increase return value by one
			icmpFinalPB = pb
			icmp = true
		}
	}

	if icmp {
		destttl++
		newTTL := destttl
		icmpFinalPB.TTL = newTTL
	}

	return destttl
}

func GetProbesByTLL(printMap map[int]*Probe, ttl int) []*Probe {
	pbs := make([]*Probe, 0)
	for _, pb := range printMap {
		if pb.TTL == ttl {
			pbs = append(pbs, pb)
		}
	}
	return pbs
}
