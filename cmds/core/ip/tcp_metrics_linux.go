// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"net"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

const tcpMetricsHelp = `Usage:	ip tcp_metrics/tcpmetrics { COMMAND | help }
	ip tcp_metrics show SELECTOR
SELECTOR := [ [ address ] PREFIX ]`

func tcpMetrics(w io.Writer) error {
	if len(arg) == 1 {
		return showTCPMetrics(w, nil)
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

			return showTCPMetrics(w, addr)
		}

		return showTCPMetrics(w, nil)
	case "help":
		fmt.Fprint(w, tcpMetricsHelp)

		return nil
	}

	return usage()
}

func showTCPMetrics(w io.Writer, address net.IP) error {
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

		fmt.Fprintf(w, "%v age %vsec %s source %v\n", v.InetDiagMsg.ID.Destination.String(), v.InetDiagMsg.Expires, tcpInfo, v.InetDiagMsg.ID.Source.String())
	}

	return nil
}
