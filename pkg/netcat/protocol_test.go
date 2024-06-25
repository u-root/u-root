// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import "testing"

func TestNetcatSocketType_String(t *testing.T) {
	testCases := []struct {
		name       string
		socketType SocketType
		expected   string
	}{
		{"TCP", SOCKET_TYPE_TCP, "tcp"},
		{"UDP", SOCKET_TYPE_UDP, "udp"},
		{"UNIX", SOCKET_TYPE_UNIX, "unix"},
		{"VSOCK", SOCKET_TYPE_VSOCK, "vsock"},
		{"SCTP", SOCKET_TYPE_SCTP, "sctp"},
		{"UDP-VSOCK", SOCKET_TYPE_UDP_VSOCK, "udp-vsock"},
		{"UNIXGRAM", SOCKET_TYPE_UDP_UNIX, "unixgram"},
		{"NONE", SOCKET_TYPE_NONE, "none"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.socketType.String()
			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestNetcatConfig_Network(t *testing.T) {
	testCases := []struct {
		name     string
		config   ProtocolOptions
		expected string
		wantErr  bool
	}{
		{name: "TCP", config: ProtocolOptions{SocketType: SOCKET_TYPE_TCP, IPType: IP_V4_V6}, expected: "tcp", wantErr: false},
		{name: "TCP IPv4", config: ProtocolOptions{SocketType: SOCKET_TYPE_TCP, IPType: IP_V4}, expected: "tcp4", wantErr: false},
		{name: "TCP IPv6", config: ProtocolOptions{SocketType: SOCKET_TYPE_TCP, IPType: IP_V6}, expected: "tcp6", wantErr: false},
		{name: "UDP", config: ProtocolOptions{SocketType: SOCKET_TYPE_UDP, IPType: IP_V4_V6}, expected: "udp", wantErr: false},
		{name: "UDP IPv4", config: ProtocolOptions{SocketType: SOCKET_TYPE_UDP, IPType: IP_V4}, expected: "udp4", wantErr: false},
		{name: "UDP IPv6", config: ProtocolOptions{SocketType: SOCKET_TYPE_UDP, IPType: IP_V6}, expected: "udp6", wantErr: false},
		{name: "UNIX", config: ProtocolOptions{SocketType: SOCKET_TYPE_UNIX}, expected: "unix", wantErr: false},
		{name: "UNIXGRAM", config: ProtocolOptions{SocketType: SOCKET_TYPE_UDP_UNIX}, expected: "unixgram", wantErr: false},
		{name: "VSOCK", config: ProtocolOptions{SocketType: SOCKET_TYPE_VSOCK}, expected: "", wantErr: false},
		{name: "UDP VSOCK", config: ProtocolOptions{SocketType: SOCKET_TYPE_UDP_VSOCK}, expected: "", wantErr: false},
		{name: "SCTP", config: ProtocolOptions{SocketType: SOCKET_TYPE_SCTP}, expected: "", wantErr: false},
		{name: "Invalid", config: ProtocolOptions{SocketType: 999, IPType: 999}, expected: "", wantErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			network, err := tc.config.Network()
			if (err != nil) != tc.wantErr {
				t.Errorf("NetcatConfig.Network() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if network != tc.expected {
				t.Errorf("NetcatConfig.Network() = %v, expected %v", network, tc.expected)
			}
		})
	}
}

func TestParseSocketType(t *testing.T) {
	tests := []struct {
		name         string
		udp          bool
		unix         bool
		vsock        bool
		sctp         bool
		expectedType SocketType
		expectError  bool
	}{
		{name: "TCP only", udp: false, unix: false, vsock: false, sctp: false, expectedType: SOCKET_TYPE_TCP, expectError: false},
		{name: "UDP only", udp: true, unix: false, vsock: false, sctp: false, expectedType: SOCKET_TYPE_UDP, expectError: false},
		{name: "UDP and UNIX", udp: true, unix: true, vsock: false, sctp: false, expectedType: SOCKET_TYPE_UDP_UNIX, expectError: false},
		{name: "UDP and VSOCK", udp: true, unix: false, vsock: true, sctp: false, expectedType: SOCKET_TYPE_UDP_VSOCK, expectError: false},
		{name: "UNIX only", udp: false, unix: true, vsock: false, sctp: false, expectedType: SOCKET_TYPE_UNIX, expectError: false},
		{name: "VSOCK only", udp: false, unix: false, vsock: true, sctp: false, expectedType: SOCKET_TYPE_VSOCK, expectError: false},
		{name: "SCTP only", udp: false, unix: false, vsock: false, sctp: true, expectedType: SOCKET_TYPE_SCTP, expectError: false},
		{name: "Invalid combination", udp: true, unix: true, vsock: true, sctp: true, expectedType: SOCKET_TYPE_NONE, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, err := ParseSocketType(tt.udp, tt.unix, tt.vsock, tt.sctp)
			if (err != nil) != tt.expectError {
				t.Errorf("ParseSocketType() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if gotType != tt.expectedType {
				t.Errorf("ParseSocketType() = %v, want %v", gotType, tt.expectedType)
			}
		})
	}
}
