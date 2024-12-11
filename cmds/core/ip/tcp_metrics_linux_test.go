// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

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

func TestTcpMetrics(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		family         int
		cursor         int
		expectedOutput string
		expectError    bool
	}{
		{
			name:        "Show TCP Metrics with valid address",
			args:        []string{"ip", "tcpmetrics", "show", "192.168.1.1"},
			cursor:      1,
			family:      netlink.FAMILY_V4,
			expectError: false,
		},
		{
			name:        "Show TCP Metrics without address",
			args:        []string{"ip", "tcpmetrics", "show"},
			family:      netlink.FAMILY_V4,
			cursor:      1,
			expectError: false,
		},
		{
			name:        "Show IPv6 TCP Metrics without address",
			args:        []string{"ip", "tcpmetrics", "show"},
			family:      netlink.FAMILY_V6,
			cursor:      1,
			expectError: false,
		},
		{
			name:        "Show all TCP Metrics without address",
			args:        []string{"ip", "tcpmetrics", "show"},
			family:      netlink.FAMILY_ALL,
			cursor:      1,
			expectError: false,
		},
		{
			name:           "Help command",
			args:           []string{"ip", "tcpmetrics", "help"},
			family:         netlink.FAMILY_V4,
			expectedOutput: tcpMetricsHelp,
			cursor:         1,
			expectError:    false,
		},
		{
			name:        "Invalid command",
			args:        []string{"ip", "tcpmetrics", "invalid"},
			family:      netlink.FAMILY_V4,
			cursor:      1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			cmd := &cmd{
				Args:   tt.args,
				Out:    &out,
				Cursor: tt.cursor,
				Family: tt.family,
			}

			err := cmd.tcpMetrics()
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}

			if tt.expectedOutput != "" {
				if out.String() != tt.expectedOutput {
					t.Errorf("expected output: %q, got: %q", tt.expectedOutput, out.String())
				}
			}
		})
	}
}
