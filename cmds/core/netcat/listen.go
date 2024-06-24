package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/u-root/u-root/pkg/netcat"
)

func (c *cmd) listenMode(output io.Writer, network, address string) error {
	var (
		err      error
		listener net.Listener
	)

	// If listing mode and Zero-I/O mode are combined the program will block indefinitely
	if c.config.ConnectionModeOptions.ZeroIO {
		for {
			time.Sleep(1 * time.Hour)
		}
	}

	if c.config.Misc.NoDNS {
		return fmt.Errorf("listen: disabling DNS resolution is not supported in listen mode")
	}

	if c.config.ConnectionModeOptions.SourceHost != "" && c.config.ConnectionModeOptions.SourcePort != "" {
		return fmt.Errorf("listen: source host/port cannot be set in listen mode")
	}

	switch c.config.ProtocolOptions.SocketType {
	case netcat.SOCKET_TYPE_TCP, netcat.SOCKET_TYPE_UNIX:
		if c.config.SSLConfig.Enabled || c.config.SSLConfig.VerifyTrust {
			tlsConfig, err := c.generateTLSConfiguration()
			if err != nil {
				return fmt.Errorf("connection: %v", err)
			}

			listener, err = tls.Listen(network, address, tlsConfig)
			if err != nil {
				return fmt.Errorf("connection: %v", err)
			}

		} else {
			listener, err = net.Listen(network, address)
			if err != nil {
				return err
			}
		}

	case netcat.SOCKET_TYPE_UDP, netcat.SOCKET_TYPE_UDP_UNIX:
		listener, err = netcat.NewUDPListener(network, address, c.config.Output.Verbose)
		if err != nil {
			return err
		}

	// unsupported socket types
	case netcat.SOCKET_TYPE_SCTP, netcat.SOCKET_TYPE_VSOCK, netcat.SOCKET_TYPE_UDP_VSOCK:
		return fmt.Errorf("currently unsupported socket type %q", c.config.ProtocolOptions.SocketType)

	case netcat.SOCKET_TYPE_NONE:
	default:
		return fmt.Errorf("undefined socket type %q", c.config.ProtocolOptions.SocketType)
	}

	return c.readFromConnections(output, listener)
}

// readFromConnections listens for incoming connections and reads from the first connection that is allowed by the access control list.
func (c *cmd) readFromConnections(output io.Writer, listener net.Listener) error {
	// If keep open is set, the maximum number of connections is set to maxConnections else it is set to 1
	var (
		maxConnections     uint32 = 1
		connectionsHandled uint32
	)

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
				netcat.Logf(c.config, "failed to resolve address: %v", err)
			}

			if c.config.AccessControl.IsAllowed(append(names, remoteAddr)) {
				atomic.AddUint32(&connectionsHandled, 1)
				// read from the connection
				if _, err := io.Copy(output, conn); err != nil {
					netcat.Logf(c.config, "run dump: %v", err)
				}

			}

			conn.Close()
		}()
	}

	return nil
}
