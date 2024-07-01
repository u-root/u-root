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

func (c *cmd) establishConnection(network, address string) (net.Conn, error) {
	var (
		err  error
		conn net.Conn
	)

	dialer := &net.Dialer{
		Timeout: c.config.Timing.Wait,
	}

	if c.config.ConnectionModeOptions.SourceHost != "" {
		switch c.config.ProtocolOptions.SocketType {

		case netcat.SOCKET_TYPE_TCP:
			dialer.LocalAddr, err = net.ResolveTCPAddr(network, fmt.Sprintf("%v:%v", c.config.ConnectionModeOptions.SourceHost, c.config.ConnectionModeOptions.SourcePort))
			if err != nil {
				return nil, fmt.Errorf("connection: failed to resolve source address %v", err)
			}

		case netcat.SOCKET_TYPE_UDP:
			dialer.LocalAddr, err = net.ResolveUDPAddr(network, fmt.Sprintf("%v:%v", c.config.ConnectionModeOptions.SourceHost, c.config.ConnectionModeOptions.SourcePort))
			if err != nil {
				return nil, fmt.Errorf("connection: failed to resolve source address %v", err)
			}

		case netcat.SOCKET_TYPE_UNIX:
			dialer.LocalAddr, err = net.ResolveUnixAddr(network, c.config.ConnectionModeOptions.SourceHost)
			if err != nil {
				return nil, fmt.Errorf("connection: failed to resolve source address %v", err)
			}

		case netcat.SOCKET_TYPE_UDP_UNIX:
			dialer.LocalAddr, err = net.ResolveUnixAddr(network, c.config.ConnectionModeOptions.SourceHost)
			if err != nil {
				return nil, fmt.Errorf("connection: failed to resolve source address %v", err)
			}

		case netcat.SOCKET_TYPE_SCTP:
			sctpAddr, err := sctp.ResolveSCTPAddr(network, address)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve SCTP address: %w", err)
			}

			return sctp.DialSCTP(network, nil, sctpAddr)

		case netcat.SOCKET_TYPE_VSOCK:
			cid, port, err := netcat.SplitVSockAddr(address)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve VSOCK address: %v", err)
			}

			return vsock.Dial(cid, port, nil)

		// unsupported socket types
		case netcat.SOCKET_TYPE_UDP_VSOCK:
			return nil, fmt.Errorf("currently unsupported socket type %q", c.config.ProtocolOptions.SocketType)

		case netcat.SOCKET_TYPE_NONE:
		default:
			return nil, fmt.Errorf("undefined socket type %q", c.config.ProtocolOptions.SocketType)
		}
	}

	// Proxy Support
	if c.config.ProxyConfig.Enabled {
		proxyDialer, err := c.proxyDialer(dialer)
		if err != nil {
			return nil, fmt.Errorf("connection: %v", err)
		}

		conn, err = proxyDialer.Dial(network, address)
		if err != nil {
			return nil, fmt.Errorf("connection: %v", err)
		}
	} else {
		// TLS Support
		if c.config.SSLConfig.Enabled || c.config.SSLConfig.VerifyTrust {
			tlsConfig, err := c.config.SSLConfig.GenerateTLSConfiguration()
			if err != nil {
				return nil, fmt.Errorf("connection: %v", err)
			}

			conn, err = tls.DialWithDialer(dialer, network, address, tlsConfig)
			if err != nil {
				return nil, fmt.Errorf("connection: %v", err)
			}
		} else {
			conn, err = dialer.Dial(network, address)
			if err != nil {
				return nil, fmt.Errorf("connection: %v", err)
			}
		}
	}

	if c.config.Timing.Timeout > 0 {
		conn.SetDeadline(time.Now().Add(c.config.Timing.Timeout))
	}

	return conn, nil
}
