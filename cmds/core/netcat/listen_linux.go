// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/ishidawataru/sctp"
	"github.com/mdlayher/vsock"
	"github.com/u-root/u-root/pkg/netcat"
)

// setupListener initializes a network listener based on the configuration provided in the cmd struct.
// It supports various protocols and configurations, including TCP, UNIX, UDP, and their secure versions with TLS.
//
// Arguments:
//   - network: A string representing the network type (e.g., "tcp", "unix").
//   - address: A string representing the address to listen on. The format is "host:port" for TCP and UDP,
//     or a file path for UNIX sockets.
//
// Returns:
//   - net.Listener: An interface representing the initialized network listener. This can be a standard Go net.Listener
//     or a custom listener for specific protocols like UDP.
func (c *cmd) setupListener(network, address string) (net.Listener, error) {
	// If listing mode and Zero-I/O mode are combined the program will block indefinitely
	if c.config.ConnectionModeOptions.ZeroIO {
		for {
			time.Sleep(1 * time.Hour)
		}
	}

	if c.config.ConnectionModeOptions.SourceHost != "" && c.config.ConnectionModeOptions.SourcePort != "" {
		return nil, fmt.Errorf("source host/port cannot be set in listen mode")
	}

	switch c.config.ProtocolOptions.SocketType {
	case netcat.SOCKET_TYPE_TCP, netcat.SOCKET_TYPE_UNIX:
		if c.config.SSLConfig.Enabled || c.config.SSLConfig.VerifyTrust {
			tlsConfig, err := c.config.SSLConfig.GenerateTLSConfiguration()
			if err != nil {
				return nil, fmt.Errorf("failed generating TLS configuration: %v", err)
			}

			return tls.Listen(network, address, tlsConfig)

		} else {
			return net.Listen(network, address)
		}

	case netcat.SOCKET_TYPE_UDP, netcat.SOCKET_TYPE_UDP_UNIX:
		return netcat.NewUDPListener(network, address, c.config.Output.Logger)

	case netcat.SOCKET_TYPE_SCTP:
		sctpAddr, err := sctp.ResolveSCTPAddr(network, address)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve SCTP address: %w", err)
		}

		return sctp.ListenSCTP(network, sctpAddr)
	case netcat.SOCKET_TYPE_VSOCK:
		cid, port, err := netcat.SplitVSockAddr(address)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve VSOCK address: %v", err)
		}

		return vsock.ListenContextID(cid, port, nil)

	// unsupported socket types
	case netcat.SOCKET_TYPE_UDP_VSOCK:
		return nil, fmt.Errorf("currently unsupported socket type %q", c.config.ProtocolOptions.SocketType)

	case netcat.SOCKET_TYPE_NONE:
	default:
		return nil, fmt.Errorf("undefined socket type %q", c.config.ProtocolOptions.SocketType)
	}

	return nil, fmt.Errorf("unexpected error")
}
