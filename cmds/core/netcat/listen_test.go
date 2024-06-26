// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"net"
	"reflect"
	"testing"

	"github.com/u-root/u-root/pkg/netcat"
)

// Define a struct for your test cases
type setupListenerTestCase struct {
	name     string
	config   *netcat.Config // Assuming Config is the type of c.config
	network  string
	address  string
	wantErr  bool
	wantType interface{}
}

// Example test cases
func TestSetupListener(t *testing.T) {
	tests := []setupListenerTestCase{
		{
			name: "TCP without SSL",
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_TCP},
			},
			network:  "tcp",
			address:  "127.0.0.1:0",
			wantErr:  false,
			wantType: &net.TCPListener{},
		},
		{
			name: "Unsupported Socket Type",
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_SCTP},
			},
			network: "sctp",
			address: "127.0.0.1:0",
			wantErr: true,
		},
		{
			name: "UDP without SSL",
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_UDP},
				SSLConfig:       netcat.SSLOptions{Enabled: false},
			},
			network:  "udp",
			address:  "127.0.0.1:0",
			wantErr:  false,
			wantType: &netcat.UDPListener{},
		},
		{
			name: "NoDNS set",
			config: &netcat.Config{
				Misc:            netcat.MiscOptions{NoDNS: true},
				ProtocolOptions: netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_TCP},
				SSLConfig:       netcat.SSLOptions{Enabled: false},
			},
			wantErr: true,
		},
		{
			name: "SourceHost set ",
			config: &netcat.Config{
				ConnectionModeOptions: netcat.ConnectModeOptions{SourceHost: "192.168.1.1"},
				ProtocolOptions:       netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_TCP},
			},
			wantErr: true,
		},
		{
			name: "SourcePort set",
			config: &netcat.Config{
				ConnectionModeOptions: netcat.ConnectModeOptions{SourcePort: "8080"},
				ProtocolOptions:       netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_TCP},
			},
			wantErr: true,
		},
		{
			name: "Unsupported Socket Type None",
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_NONE},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &cmd{config: tc.config}
			listener, err := c.setupListener(tc.network, tc.address)
			if (err != nil) != tc.wantErr {
				t.Errorf("setupListener() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if err == nil {
				defer listener.Close()
				if reflect.TypeOf(listener) != reflect.TypeOf(tc.wantType) {
					t.Errorf("Expected listener type %v, got %v", reflect.TypeOf(tc.wantType), reflect.TypeOf(listener))
				}
			}
		})
	}
}
