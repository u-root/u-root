// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/netcat"
)

var connectionTests = []struct {
	name   string
	config netcat.NetcatConfig
	args   []string
	stdin  string
}{
	{
		name: "TCP Listen IPv4",
		config: netcat.NetcatConfig{
			ConnectionMode: netcat.CONNECTION_MODE_LISTEN,
			ProtocolOptions: netcat.NetcatProtocolOptions{
				IPType:     netcat.IP_V4,
				SocketType: netcat.SOCKET_TYPE_TCP,
			},
		},
		args: []string{""},
	},
	{
		name: "UDP Listen IPv4",
		config: netcat.NetcatConfig{
			ConnectionMode: netcat.CONNECTION_MODE_LISTEN,
			ProtocolOptions: netcat.NetcatProtocolOptions{
				IPType:     netcat.IP_V4,
				SocketType: netcat.SOCKET_TYPE_UDP,
			},
		},
		args: []string{""},
	},
	{
		name: "TCP Dial IPv4",
		config: netcat.NetcatConfig{
			ConnectionMode: netcat.CONNECTION_MODE_CONNECT,
			ProtocolOptions: netcat.NetcatProtocolOptions{
				IPType:     netcat.IP_V4,
				SocketType: netcat.SOCKET_TYPE_TCP,
			},
		},
		args: []string{""},
	},
	{
		name: "UDP Dial IPv4",
		config: netcat.NetcatConfig{
			ConnectionMode: netcat.CONNECTION_MODE_CONNECT,
			ProtocolOptions: netcat.NetcatProtocolOptions{
				IPType:     netcat.IP_V4,
				SocketType: netcat.SOCKET_TYPE_UDP,
			},
		},
		args: []string{""},
	},
}

func TestConnection(t *testing.T) {
	for _, tt := range connectionTests {
		stdin := strings.NewReader("hello client")
		stdout := bufio.NewWriter(&bytes.Buffer{})
		stderr := bufio.NewWriter(&bytes.Buffer{})
		// cmdIO := io.NewReaderWriter(stdin, stdout, stderr)

		t.Run(tt.name, func(t *testing.T) {
			cmd, err := command(stdin, stdout, stderr, tt.config, tt.args)
			if err != nil {
				t.Fatalf("command() = %v", err)
			}

			con, err := cmd.connection()
			if err != nil {
				t.Fatalf("cmd.connection() = %v, want nil", err)
			}

			if cmd.config.ConnectionMode == netcat.CONNECTION_MODE_LISTEN {
				// receive connection
				go func() {
					if _, err := io.Copy(con, cmd.stdin); err != nil {
						fmt.Fprintln(cmd.stderr, err)
					}
				}()
			} else {
				// Make connection to fake server
			}
		})
	}
}
