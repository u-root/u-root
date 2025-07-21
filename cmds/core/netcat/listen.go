// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/mdlayher/vsock"
	"github.com/u-root/u-root/pkg/netcat"
)

var osListeners = map[netcat.SocketType]func(string, string) (net.Listener, error){}

type listenFn func(output io.WriteCloser, network, address string) error

func (c *cmd) listenMode(output io.WriteCloser, network, address string, listen listenFn) error {
	if network == "tcp" || network == "udp" {
		err4 := listen(output, network+"4", address)
		if err4 == nil {
			return nil
		}
		err6 := listen(output, network+"6", address)
		if err6 == nil {
			return nil
		}
		return fmt.Errorf("listen mode: %w", errors.Join(err4, err6))
	}

	return listen(output, network, address)
}

func (c *cmd) listen(output io.WriteCloser, network, address string) error {
	listener, err := c.setupListener(network, address)
	if err != nil {
		return fmt.Errorf("failed to setup listener: %w", err)
	}

	return c.listenForConnections(output, listener, 0)
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
	case netcat.SOCKET_TYPE_TCP:
		if c.config.SSLConfig.Enabled {
			tlsConfig, err := c.config.SSLConfig.GenerateTLSConfiguration(true)
			if err != nil {
				return nil, fmt.Errorf("failed generating TLS configuration: %w", err)
			}

			return tls.Listen(network, address, tlsConfig)

		}
		fallthrough

	case netcat.SOCKET_TYPE_UNIX:
		return net.Listen(network, address)

	case netcat.SOCKET_TYPE_UDP, netcat.SOCKET_TYPE_UDP_UNIX:
		return netcat.NewUDPListener(network, address, c.config.Output.Logger)

	case netcat.SOCKET_TYPE_VSOCK:
		cid, port, err := netcat.SplitVSockAddr(address)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve VSOCK address: %w", err)
		}

		return vsock.ListenContextID(cid, port, nil)

	default:
		l, ok := osListeners[c.config.ProtocolOptions.SocketType]
		if !ok {
			return nil, fmt.Errorf("currently unsupported socket type %q", c.config.ProtocolOptions.SocketType)
		}
		return l(network, address)
	}
}

// acceptAllowed accepts a connection from listener such that the peer address
// is permitted by the allow list. acceptAllowed drops connections from
// forbidden peers transparently (apart from counting test attempts during unit
// testing).
func (c *cmd) acceptAllowed(listener net.Listener, testLimit uint32, testID *uint32) (net.Conn, error) {
	for testLimit == 0 || *testID < testLimit {
		conn, err := listener.Accept()

		if testLimit > 0 {
			(*testID)++
		}

		if err != nil {
			return nil, err
		}

		if c.config.ProtocolOptions.SocketType == netcat.SOCKET_TYPE_TCP &&
			!c.config.AccessControl.IsAllowed(parseRemoteAddr(c.config.ProtocolOptions.SocketType, conn.RemoteAddr().String())) {
			conn.Close()
			continue
		}
		return conn, nil
	}
	return nil, errors.New("testLimit exhausted")
}

// fullCopy copies src to dst with io.Copy, transparently retrying on
// io.ErrShortWrite.
func fullCopy(dst io.Writer, src io.Reader) error {
	var err error
	for {
		_, err = io.Copy(dst, src)
		if err != nil && errors.Is(err, io.ErrShortWrite) {
			continue
		}
		break
	}
	return err
}

// connections holds all the active connections of a listener.
type connections struct {
	capacity    uint32
	used        uint32
	connections map[uint32]net.Conn
	mutex       sync.Mutex
	isAvailable *sync.Cond
	isEmpty     *sync.Cond
}

func newConnections(capacity uint32) *connections {
	conns := connections{
		capacity:    capacity,
		connections: make(map[uint32]net.Conn),
		mutex:       sync.Mutex{},
	}
	conns.isAvailable = sync.NewCond(&conns.mutex)
	conns.isEmpty = sync.NewCond(&conns.mutex)
	return &conns
}

// add adds a new connection. If the number of concurrent connections has
// reached the maximum (i.e., capacity has been exhausted) before entry to the
// function, add blocks until another goroutine calls delete.
func (c *connections) add(id uint32, conn net.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for c.used == c.capacity {
		c.isAvailable.Wait()
	}
	c.used++
	c.connections[id] = conn
}

func (c *connections) delete(id uint32) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.used == c.capacity {
		c.isAvailable.Signal()
	}
	c.connections[id].Close()
	delete(c.connections, id)
	c.used--
	if c.used == 0 {
		c.isEmpty.Signal()
	}
}

func (c *connections) drain() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for c.used > 0 {
		c.isEmpty.Wait()
	}
}

// broadcast sends a message to all connections in the connections object (except the sender) and the output writer.
func (c *connections) broadcast(output io.Writer, senderID uint32, message string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, err := io.WriteString(output, message); err != nil && !errors.Is(err, io.ErrShortWrite) {
		log.Printf("failed to write to output: %v", err)
	}

	for id, conn := range c.connections {
		if id == senderID {
			continue
		}

		if _, err := io.WriteString(conn, message); err != nil && !errors.Is(err, io.ErrShortWrite) {
			log.Printf("failed to write to connection %v: %v", id, err)
			break
		}
	}
}

// acceptSingle is the main transfer routine for listen mode if neither
// keep-open nor broker mode is enabled. The logic here reflects connect mode:
// after accepting a single connection, acceptSingle copies stdin to socket,
// and copies socket to output. In both directions, acceptSingle propagates
// EOF. acceptSingle returns when transfers have completed in both directions.
func (c *cmd) acceptSingle(output io.WriteCloser, listener net.Listener) error {
	conn, err := c.acceptAllowed(listener, 0, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	stdinToConnErrChan := make(chan error)
	go func() {
		innerErr := fullCopy(conn, c.stdin)
		if innerErr == nil && !c.config.Misc.NoShutdown {
			innerErr = netcat.CloseWrite(conn)
		}
		stdinToConnErrChan <- innerErr
	}()

	err = fullCopy(output, conn)
	if err == nil {
		err = output.Close()
	}

	stdinToConnErr := <-stdinToConnErrChan

	if stdinToConnErr != nil {
		return stdinToConnErr
	}
	return err
}

// acceptForever is the main transfer routine for listen mode if keep-open or
// broker mode is enabled. The function never returns (unless "testLimit" is
// set to a positive value, for unit testing).
//
// acceptForever keeps accepting new connections. Each connection is handled
// (i.e., read) in a separate goroutine, concurrently with the others. Whenever
// the maximum number of simultaneous connections is reached, acceptForever
// doesn't proceed with another connection until one of the existent
// connections completes (for some definition of "completes"; see below).
//
// In broker mode:
//
// (1) acceptForever does not read stdin.
//
// (2) For each connection, in separation:
//
// (2.1) acceptForever keeps reading another line of data, and writes it out to
// every different connection, plus "output". With chat mode enabled,
// acceptForever prefixes each line written with the identifier of the source
// connection (i.e., where the line has been read from). These prefixes are
// 1-based, not 0-based.
//
// (2.2) acceptForever completes (closes and deletes) the connection when EOF
// is read from the connection.
//
// (3) acceptForever never propagates EOF to any connection or to output, as
// connections accepted in the future can always send data to current
// connections and to output.
//
// In keep-open mode:
//
// (1) For each connection, in separation:
//
// (1.1) acceptForever copies all data from the connection to output.
//
// (1.2) For the first connection accepted, acceptForever copies all data from
// stdin, and propagates EOF, to the connection.
//
// (1.3) For connections accepted after the first one, acceptForever sends an
// EOF at once.
//
// (1.4) acceptForever completes (closes and deletes) the connection when
// transfers have completed in both directions.
//
// (2) acceptForever never propagates EOF to output, as connections accepted in
// the future can always send data to output.
func (c *cmd) acceptForever(output io.WriteCloser, listener net.Listener,
	testLimit uint32,
) error {
	var testID uint32
	var connID uint32

	conns := newConnections(c.config.ListenModeOptions.MaxConnections)
	for testLimit == 0 || testID < testLimit {
		conn, err := c.acceptAllowed(listener, testLimit, &testID)
		if err != nil {
			log.Printf("wait for connection: %v", err)
			continue
		}

		conns.add(connID, conn)

		go func(myConnID uint32, myConn net.Conn) {
			defer conns.delete(myConnID)

			var myConnErr error

			if c.config.ListenModeOptions.BrokerMode {
				scanner := bufio.NewScanner(myConn)
				for scanner.Scan() {
					line := scanner.Text()

					var formattedLine string
					if c.config.ListenModeOptions.ChatMode {
						formattedLine = fmt.Sprintf("user<%d>: %s\n", myConnID+1, line)
					} else {
						formattedLine = fmt.Sprintf("%s\n", line)
					}

					conns.broadcast(output, myConnID, formattedLine)
				}

				myConnErr = scanner.Err()
				if myConnErr != nil {
					log.Printf("failed to scan connection: %v", myConnErr)
				}

				return
			}

			stdinToConnErrChan := make(chan error)
			go func() {
				var innerErr error
				if myConnID == 0 {
					innerErr = fullCopy(myConn, c.stdin)
				}
				if innerErr == nil && !c.config.Misc.NoShutdown {
					innerErr = netcat.CloseWrite(myConn)
				}
				stdinToConnErrChan <- innerErr
			}()

			myConnErr = fullCopy(output, myConn)
			if myConnErr != nil {
				log.Printf("failed to copy connection to output: %v", myConnErr)
			}

			myConnErr = <-stdinToConnErrChan
			if myConnErr != nil {
				log.Printf("failed to send on connection: %v", myConnErr)
			}
		}(connID, conn)

		connID++
	}

	conns.drain()
	return nil
}

// acceptSingleUDP is the transfer routine for listen mode with UDP.
func (c *cmd) acceptSingleUDP(output io.Writer, listener net.Listener) error {
	conn, err := listener.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	return c.transferPackets(output, conn.(net.PacketConn), true)
}

// listenForConnections listens for incoming connections on a specified listener and reads data from these.
// Arguments:
//   - output: The io.Writer object to which the function writes the data read from the connections.
//   - listener: The net.Listener object on which the function listens for incoming connections. This listener should already be initialized
//     and listening on the desired port.
//   - testLimit: in keep-open or broker mode, stop accepting connections after
//     this many connections have been accepted
func (c *cmd) listenForConnections(output io.WriteCloser, listener net.Listener, testLimit uint32) error {
	if c.config.ListenModeOptions.KeepOpen || c.config.ListenModeOptions.BrokerMode {
		return c.acceptForever(output, listener, testLimit)
	}
	if c.config.ProtocolOptions.SocketType == netcat.SOCKET_TYPE_UDP ||
		c.config.ProtocolOptions.SocketType == netcat.SOCKET_TYPE_UDP_UNIX {
		return c.acceptSingleUDP(output, listener)
	}
	return c.acceptSingle(output, listener)
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
