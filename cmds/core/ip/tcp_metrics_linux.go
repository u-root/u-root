// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

const tcpMetricsHelp = `Usage:	ip tcp_metrics/tcpmetrics { COMMAND | help }
	ip tcp_metrics show SELECTOR
SELECTOR := [ [ address ] PREFIX ]`

func (cmd cmd) tcpMetrics() error {
	if len(arg) == 1 {
		return cmd.showTCPMetrics(nil)
	}

	expectedValues = []string{"show", "help"}
	switch arg[1] {
	case "show":
		cursor++
		if len(arg) > 2 {
			addr, err := parseAddress()
			if err != nil {
				return err
			}

			return cmd.showTCPMetrics(addr)
		}

		return cmd.showTCPMetrics(nil)
	case "help":
		fmt.Fprint(cmd.out, tcpMetricsHelp)

		return nil
	}

	return usage()
}

func (cmd cmd) showTCPMetrics(address net.IP) error {
	var family uint8 = unix.AF_INET
	if family == netlink.FAMILY_V6 {
		family = unix.AF_INET6
	}

	resp, err := netlink.SocketDiagTCPInfo(family)
	if err != nil {
		return err
	}

	for _, v := range resp {
		if v.InetDiagMsg.ID.Destination.IsUnspecified() {
			continue
		}

		if address != nil && !v.InetDiagMsg.ID.Destination.Equal(address) {
			continue
		}

		var tcpInfo string

		if v.TCPInfo != nil {
			tcpInfo = fmt.Sprintf("cwnd %v rtt %v rttvar %vus", v.TCPInfo.Snd_cwnd, v.TCPInfo.Rtt, v.TCPInfo.Rttvar)
		}

		fmt.Fprintf(cmd.out, "%v age %vsec %s source %v\n", v.InetDiagMsg.ID.Destination.String(), v.InetDiagMsg.Expires, tcpInfo, v.InetDiagMsg.ID.Source.String())
	}

	return nil
}
