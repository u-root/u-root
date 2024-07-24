// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"bytes"
	"net"
	"testing"

	"github.com/vishvananda/netlink"
)

func TestPrintTCPMetrics(t *testing.T) {
	// Mock data
	resp := []*netlink.InetDiagTCPInfoResp{
		{
			InetDiagMsg: &netlink.Socket{
				ID: netlink.SocketID{
					Source:      net.ParseIP("192.168.1.2"),
					Destination: net.ParseIP("192.168.1.1"),
				},
				Expires: 123,
			},
			TCPInfo: &netlink.TCPInfo{
				Snd_cwnd: 10,
				Rtt:      100,
				Rttvar:   50,
			},
		},
		{
			InetDiagMsg: &netlink.Socket{
				ID: netlink.SocketID{
					Source: net.ParseIP("192.168.1.2"),
				},
				Expires: 123,
			},
			TCPInfo: &netlink.TCPInfo{
				Snd_cwnd: 10,
				Rtt:      100,
				Rttvar:   50,
			},
		},
		{
			InetDiagMsg: &netlink.Socket{
				ID: netlink.SocketID{
					Source:      net.ParseIP("192.168.1.2"),
					Destination: net.ParseIP("123.0.0.4"),
				},
				Expires: 123,
			},
			TCPInfo: &netlink.TCPInfo{
				Snd_cwnd: 10,
				Rtt:      100,
				Rttvar:   50,
			},
		},
	}

	// Mock output buffer
	var out bytes.Buffer
	cmd := cmd{Out: &out}

	// Call the function with the mock data
	cmd.printTCPMetrics(resp, net.ParseIP("192.168.1.1"))

	// Expected output
	expectedOutput := "192.168.1.1 age 123sec cwnd 10 rtt 100 rttvar 50us source 192.168.1.2\n"

	// Check the output
	if out.String() != expectedOutput {
		t.Errorf("expected %q, got %q", expectedOutput, out.String())
	}
}
