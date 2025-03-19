// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConnectionModeString(t *testing.T) {
	tests := []struct {
		name     string
		mode     ConnectionMode
		expected string
	}{
		{"Connect mode", CONNECTION_MODE_CONNECT, "connect"},
		{"Listen mode", CONNECTION_MODE_LISTEN, "listen"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.expected {
				t.Errorf("ConnectionMode.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAddressConnectMode(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config // Use a pointer to Config
		wantAddr string
		wantErr  bool
	}{
		{
			name: "TCPv4 Connect with valid host",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_CONNECT,
				Host:           "127.0.0.1",
				Port:           8080,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
			},
			wantAddr: "127.0.0.1:8080",
			wantErr:  false,
		},
		{
			name: "TCPv6 Connect with valid host",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_CONNECT,
				Host:           "::1",
				Port:           8080,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
			},
			wantAddr: "[::1]:8080",
			wantErr:  false,
		},
		{
			name: "TCPv4 Connect with no DNS host",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_CONNECT,
				Host:           "127.0.0.1",
				Port:           8080,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
				Misc: MiscOptions{
					NoDNS: true,
				},
			},
			wantAddr: "127.0.0.1:8080",
			wantErr:  false,
		},
		{
			name: "TCPv6 Connect with no DNS host",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_CONNECT,
				Host:           "::1",
				Port:           8080,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
				Misc: MiscOptions{
					NoDNS: true,
				},
			},
			wantAddr: "[::1]:8080",
			wantErr:  false,
		},
		{
			name: "TCPv4 Portscan",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_CONNECT,
				Host:           "127.0.0.1",
				Port:           1234,
				ConnectionModeOptions: ConnectModeOptions{
					ScanPorts:   true,
					CurrentPort: 2345,
					EndPort:     3456,
				},
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
			},
			wantAddr: "127.0.0.1:2345",
			wantErr:  false,
		},
		{
			name: "TCPv6 Portscan",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_CONNECT,
				Host:           "::1",
				Port:           1234,
				ConnectionModeOptions: ConnectModeOptions{
					ScanPorts:   true,
					CurrentPort: 2345,
					EndPort:     3456,
				},
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
			},
			wantAddr: "[::1]:2345",
			wantErr:  false,
		},

		{
			name: "TCP Connect with no DNS host failing",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_CONNECT,
				Host:           "dns.host",
				Port:           8080,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
				Misc: MiscOptions{
					NoDNS: true,
				},
			},
			wantErr: true,
		},
		{
			name: "TCPv4 Connect with valid host and default port",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_CONNECT,
				Host:           "127.0.0.1",
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
			},
			wantAddr: "127.0.0.1:0", // Assuming default port is 0 when not specified
			wantErr:  false,
		},
		{
			name: "TCPv6 Connect with valid host and default port",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_CONNECT,
				Host:           "::1",
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
			},
			wantAddr: "[::1]:0", // Assuming default port is 0 when not specified
			wantErr:  false,
		},

		{
			name: "Invalid Connection Mode",
			config: &Config{
				ConnectionMode: 999, // Assuming 999 is an invalid connection mode
				Host:           "127.0.0.1",
				Port:           8080,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
			},
			wantAddr: "",
			wantErr:  true,
		},
		{
			name: "Empty Host",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_CONNECT,
				Host:           "",
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
			},
			wantAddr: ":8080",
			wantErr:  true,
		},
		{
			name: "UNIX Socket",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_CONNECT,
				Host:           "/tmp/sock",
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_UDP_UNIX,
				},
			},
			wantAddr: "/tmp/sock",
			wantErr:  false,
		},
		{
			name: "SCTPv4",
			config: &Config{
				Host:           "127.0.0.1",
				Port:           8080,
				ConnectionMode: CONNECTION_MODE_CONNECT,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_SCTP,
				},
			},
			wantAddr: "127.0.0.1:8080",
		},
		{
			name: "SCTPv6",
			config: &Config{
				Host:           "::1",
				Port:           8080,
				ConnectionMode: CONNECTION_MODE_CONNECT,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_SCTP,
				},
			},
			wantAddr: "[::1]:8080",
		},
		{
			name: "VSOCK",
			config: &Config{
				Host:           "123",
				Port:           8080,
				ConnectionMode: CONNECTION_MODE_CONNECT,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_VSOCK,
				},
			},
			wantAddr: "123:8080",
		},
		{
			name: "Unsupported Socket Type",
			config: &Config{
				Host:           "/tmp/sock",
				ConnectionMode: CONNECTION_MODE_CONNECT,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_UDP_VSOCK,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := tt.config.Address()
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Errorf("Address() error = %v, wantErr %v", err, tt.wantErr)
			}

			if addr != tt.wantAddr {
				t.Errorf("Address() = %v, want %v", addr, tt.wantAddr)
			}
		})
	}
}

func TestAddressListenMode(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		wantAddr string
		wantErr  bool
	}{
		{
			name: "TCPv4 Listen with valid host",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_LISTEN,
				Host:           "127.0.0.1",
				Port:           8080,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
			},
			wantAddr: "127.0.0.1:8080",
			wantErr:  false,
		},
		{
			name: "TCPv6 Listen with valid host",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_LISTEN,
				Host:           "::1",
				Port:           8080,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
			},
			wantAddr: "[::1]:8080",
			wantErr:  false,
		},
		{
			name: "TCPv4 Listen with valid host and default port",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_LISTEN,
				Host:           "127.0.0.1",
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
			},
			wantAddr: "127.0.0.1:0", // Assuming default port is 0 when not specified
			wantErr:  false,
		},
		{
			name: "TCPv6 Listen with valid host and default port",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_LISTEN,
				Host:           "::1",
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
				},
			},
			wantAddr: "[::1]:0", // Assuming default port is 0 when not specified
			wantErr:  false,
		},
		{
			name: "Empty Host with IPV4_V6",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_LISTEN,
				Port:           DEFAULT_PORT,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
					IPType:     IP_V4_V6,
				},
			},
			wantAddr: fmt.Sprintf("%s:%d", DEFAULT_IPV4_ADDRESS, DEFAULT_PORT),
			wantErr:  false,
		},
		{
			name: "Empty Host with IPV4",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_LISTEN,
				Port:           DEFAULT_PORT,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
					IPType:     IP_V4,
				},
			},
			wantAddr: fmt.Sprintf("%s:%d", DEFAULT_IPV4_ADDRESS, DEFAULT_PORT),
			wantErr:  false,
		},
		{
			name: "Empty Host with IPV6",
			config: &Config{
				Port:           DEFAULT_PORT,
				ConnectionMode: CONNECTION_MODE_LISTEN,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_TCP,
					IPType:     IP_V6,
				},
			},
			wantAddr: fmt.Sprintf("[%s]:%d", DEFAULT_IPV6_ADDRESS, DEFAULT_PORT),
			wantErr:  false,
		},
		{
			name: "Empty Host with UNIX",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_LISTEN,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_UDP_UNIX,
				},
			},
			wantAddr: DEFAULT_UNIX_SOCKET,
			wantErr:  false,
		},
		{
			name: "Explicit pathname for UNIX domain socket (stream)",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_LISTEN,
				Host:           "/tmp/stream.sock",
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_UNIX,
				},
			},
			wantAddr: "/tmp/stream.sock",
			wantErr:  false,
		},
		{
			name: "Explicit pathname for UNIX domain socket (datagram)",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_LISTEN,
				Host:           "/tmp/datagram.sock",
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_UDP_UNIX,
				},
			},
			wantAddr: "/tmp/datagram.sock",
			wantErr:  false,
		},
		{
			name: "Unsupported Socket Type",
			config: &Config{
				ConnectionMode: CONNECTION_MODE_LISTEN,
				ProtocolOptions: ProtocolOptions{
					SocketType: SOCKET_TYPE_VSOCK,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := tt.config.Address()
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Errorf("Address() error = %v, wantErr %v", err, tt.wantErr)
			}

			if addr != tt.wantAddr {
				t.Errorf("Address() = %v, want %v", addr, tt.wantAddr)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	defaultConfig := DefaultConfig()

	if defaultConfig.Port != DEFAULT_PORT {
		t.Errorf("expected Port to be %d, got %d", DEFAULT_PORT, defaultConfig.Port)
	}

	if defaultConfig.ConnectionMode != DEFAULT_CONNECTION_MODE {
		t.Errorf("expected ConnectionMode to be %d, got %d", DEFAULT_CONNECTION_MODE, defaultConfig.ConnectionMode)
	}

	if len(defaultConfig.ConnectionModeOptions.LooseSourceRouterPoints) != 0 {
		t.Errorf("expected LooseSourceRouterPoints to be empty, got %v", defaultConfig.ConnectionModeOptions.LooseSourceRouterPoints)
	}

	if defaultConfig.ConnectionModeOptions.SourcePort != DEFAULT_SOURCE_PORT {
		t.Errorf("expected SourcePort to be %v, got %v", DEFAULT_SOURCE_PORT, defaultConfig.ConnectionModeOptions.SourcePort)
	}

	if defaultConfig.ListenModeOptions.MaxConnections != DEFAULT_CONNECTION_MAX {
		t.Errorf("expected MaxConnections to be %d, got %d", DEFAULT_CONNECTION_MAX, defaultConfig.ListenModeOptions.MaxConnections)
	}

	if defaultConfig.ProtocolOptions.IPType != DEFAULT_IP_TYPE {
		t.Errorf("expected IPType to be %d, got %d", DEFAULT_IP_TYPE, defaultConfig.ProtocolOptions.IPType)
	}

	if defaultConfig.ProtocolOptions.SocketType != SOCKET_TYPE_TCP {
		t.Errorf("expected SocketType to be %d, got %d", SOCKET_TYPE_TCP, defaultConfig.ProtocolOptions.SocketType)
	}

	if cmp.Diff(defaultConfig.SSLConfig.Ciphers, []string{}) != "" {
		t.Errorf("expected Ciphers to be %s, got %s", DEFAULT_SSL_SUITE_STR, defaultConfig.SSLConfig.Ciphers)
	}

	if defaultConfig.ProxyConfig.Type != PROXY_TYPE_NONE {
		t.Errorf("expected Proxy Type to be %d, got %d", PROXY_TYPE_NONE, defaultConfig.ProxyConfig.Type)
	}

	if len(defaultConfig.AccessControl.ConnectionList) != 0 {
		t.Errorf("expected ConnectionList to be empty, got %v", defaultConfig.AccessControl.ConnectionList)
	}

	if defaultConfig.CommandExec.Type != EXEC_TYPE_NONE {
		t.Errorf("expected CommandExec Type to be %d, got %d", EXEC_TYPE_NONE, defaultConfig.CommandExec.Type)
	}

	if cmp.Diff(defaultConfig.Misc.EOL, DEFAULT_LF) != "" {
		t.Errorf("expected EOL to be %s, got %s", DEFAULT_LF, defaultConfig.Misc.EOL)
	}

	if defaultConfig.Timing.Wait != DEFAULT_WAIT {
		t.Errorf("expected Wait to be %d, got %d", DEFAULT_WAIT, defaultConfig.Timing.Wait)
	}
}
