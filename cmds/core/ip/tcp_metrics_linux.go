// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

const tcpMetricsHelp = `Usage:	ip tcp_metrics/tcpmetrics { COMMAND | help }
	ip tcp_metrics show SELECTOR
SELECTOR := [ [ address ] PREFIX ]
`

func (cmd *cmd) tcpMetrics() error {
	if !cmd.tokenRemains() {
		return cmd.showTCPMetrics(nil)
	}

	switch cmd.nextToken("show", "help") {
	case "show":
		if cmd.tokenRemains() {
			addr, err := cmd.parseAddress()
			if err != nil {
				return err
			}

			return cmd.showTCPMetrics(addr)
		}

		return cmd.showTCPMetrics(nil)
	case "help":
		fmt.Fprint(cmd.Out, tcpMetricsHelp)

		return nil
	}

	return cmd.usage()
}

func (cmd *cmd) showTCPMetrics(address net.IP) error {
	var (
		resp []*netlink.InetDiagTCPInfoResp
		err  error
	)

	if cmd.Family > 255 || cmd.Family < 0 {
		return fmt.Errorf("invalid protocol family %v", cmd.Family)
	}

	if cmd.Family == netlink.FAMILY_ALL {
		responseIP4, err := netlink.SocketDiagTCPInfo(uint8(netlink.FAMILY_V4))
		if err != nil {
			return fmt.Errorf("failed to get TCP metrics: %w", err)
		}

		responseIP6, err := netlink.SocketDiagTCPInfo(uint8(netlink.FAMILY_V6))
		if err != nil {
			return fmt.Errorf("failed to get TCP metrics: %w", err)
		}

		resp = append(responseIP4, responseIP6...)
	} else {
		resp, err = netlink.SocketDiagTCPInfo(uint8(cmd.Family))
		if err != nil {
			return fmt.Errorf("failed to get TCP metrics: %w", err)
		}
	}

	cmd.printTCPMetrics(resp, address)
	return nil
}

func (cmd *cmd) printTCPMetrics(resp []*netlink.InetDiagTCPInfoResp, address net.IP) {
	for _, v := range resp {
		if v.InetDiagMsg.ID.Destination.IsUnspecified() || v.InetDiagMsg.ID.Source.IsUnspecified() || v.InetDiagMsg.ID.Source == nil || v.InetDiagMsg.ID.Destination == nil {
			continue
		}

		if address != nil && !v.InetDiagMsg.ID.Destination.Equal(address) {
			continue
		}

		var tcpInfo string

		if v.TCPInfo != nil {
			tcpInfo = fmt.Sprintf("cwnd %v rtt %v rttvar %vus", v.TCPInfo.Snd_cwnd, v.TCPInfo.Rtt, v.TCPInfo.Rttvar)
		}

		fmt.Fprintf(cmd.Out, "%v age %vsec %s source %v\n", v.InetDiagMsg.ID.Destination.String(), v.InetDiagMsg.Expires, tcpInfo, v.InetDiagMsg.ID.Source.String())
	}
}
