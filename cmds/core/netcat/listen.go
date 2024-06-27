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
	"sync"
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

// readFromConnections listens for incoming connections on a specified listener and reads data from these.
// If keep open is set, the maximum number of connections is set to maxConnections else it is set to 1.
// Arguments:
//   - listener: The net.Listener object on which the function listens for incoming connections. This listener should already be initialized
//     and listening on the desired port.
//   - acl: An AccessControlList object that contains the rules for which connections are allowed to communicate with this service.
//     The function uses this list to determine if an incoming connection should be accepted or rejected based on the source address.
func (c *cmd) readFromConnections(output io.Writer, listener net.Listener) error {
	var (
		connectionsHandled uint32
		wg                 sync.WaitGroup
	)

	maxConnections := c.config.ListenModeOptions.MaxConnections

	for {
		if atomic.LoadUint32(&connectionsHandled) >= maxConnections {
			break // Stop accepting new connections if max is reached
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		if !c.config.AccessControl.IsAllowed(parseRemoteAddr(c.config.ProtocolOptions.SocketType, conn.RemoteAddr().String())) {
			defer conn.Close()
			break
		}

		atomic.AddUint32(&connectionsHandled, 1)
		wg.Add(1)

		go func(conn net.Conn) {
			defer wg.Done()
			defer conn.Close()

			// Read from the connection.
			if _, err := io.Copy(output, conn); err != nil {
				log.Printf("failed to read from connection: %v", err)
			}
		}(conn)
	}

	wg.Wait() // Wait for all connections to finish

	return nil
}

// parseRemoteAddr parses the remote address of a connection and returns a list of possible addresses.
// For UNIX sockets, the returned address is the path to the socket file.
// For TCP and UDP sockets, the remote addresses are combinations of IP address and port and any domain name.
func parseRemoteAddr(socketType netcat.SocketType, remoteAddr string) []string {
	addresses := []string{remoteAddr}
	switch socketType {
	case netcat.SOCKET_TYPE_TCP, netcat.SOCKET_TYPE_UDP:
		// Strip the port from the remoteAddr, if error occurs, skip this step
		host, _, err := net.SplitHostPort(remoteAddr)
		if err == nil {
			addresses = append(addresses, host)
			// If the address is not in the format host:port, use the original remoteAddr as the host
		} else {
			host = remoteAddr
		}

		// Perform a reverse lookup to get the domain names associated with the host.
		names, err := net.LookupAddr(host)
		if err != nil {
			log.Printf("failed to resolve host address: %v", err)
		}

		return append(addresses, names...)
	case netcat.SOCKET_TYPE_NONE:
		log.Printf("socket type not set, using remote address as is")
	case netcat.SOCKET_TYPE_SCTP, netcat.SOCKET_TYPE_VSOCK, netcat.SOCKET_TYPE_UDP_VSOCK:
		log.Printf("unsupported socket type %q", socketType)
	default:
	}

	return addresses
}
