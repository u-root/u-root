// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/u-root/u-root/pkg/netcat"
)

func (c *cmd) listenMode(output io.Writer, network, address string) error {
	listener, err := c.setupListener(network, address)
	if err != nil {
		return fmt.Errorf("failed to setup listener: %v", err)
	}

	return c.readFromConnections(output, listener)
}

func (c *cmd) setupListener(network, address string) (net.Listener, error) {
	// If listing mode and Zero-I/O mode are combined the program will block indefinitely
	if c.config.ConnectionModeOptions.ZeroIO {
		for {
			time.Sleep(1 * time.Hour)
		}
	}

	if c.config.Misc.NoDNS {
		return nil, fmt.Errorf("disabling DNS resolution is not supported in listen mode")
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

	// unsupported socket types
	case netcat.SOCKET_TYPE_SCTP, netcat.SOCKET_TYPE_VSOCK, netcat.SOCKET_TYPE_UDP_VSOCK:
		return nil, fmt.Errorf("currently unsupported socket type %q", c.config.ProtocolOptions.SocketType)

	case netcat.SOCKET_TYPE_NONE:
	default:
		return nil, fmt.Errorf("undefined socket type %q", c.config.ProtocolOptions.SocketType)
	}

	return nil, fmt.Errorf("unexpected error")
}

// readFromConnections listens for incoming connections and reads from the first connection that is allowed by the access control list.
func (c *cmd) readFromConnections(output io.Writer, listener net.Listener) error {
	// If keep open is set, the maximum number of connections is set to maxConnections else it is set to 1
	var (
		maxConnections     uint32 = 1
		connectionsHandled uint32
	)

	log.Printf("Listening on %s", listener.Addr().String())
	if c.config.ListenModeOptions.KeepOpen {
		maxConnections = c.config.ListenModeOptions.MaxConnections
	}

	for atomic.LoadUint32(&connectionsHandled) < maxConnections {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			remoteAddr := conn.RemoteAddr().String()

			// Perform a reverse lookup to get the domain names associated with the address
			names, err := net.LookupAddr(remoteAddr)
			if err != nil {
				log.Printf("failed to resolve host address: %v", err)
			}

			if c.config.AccessControl.IsAllowed(append(names, remoteAddr)) {
				atomic.AddUint32(&connectionsHandled, 1)
				// read from the connection
				if _, err := io.Copy(output, conn); err != nil {
					log.Printf("failed to read from connection: %v", err)
				}

			}

			conn.Close()
		}()
	}

	return nil
}
