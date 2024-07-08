// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/u-root/u-root/pkg/netcat"
)

func (c *cmd) listenMode(output io.Writer, network, address string) error {
	listener, err := c.setupListener(network, address)
	if err != nil {
		return fmt.Errorf("failed to setup listener: %v", err)
	}

	return c.listenForConnections(output, listener)
}

// Connections holds all the active connections of a listener.
type Connections struct {
	Connections map[uint32]net.Conn
	mutex       sync.Mutex
}

func NewConnections() *Connections {
	return &Connections{
		Connections: make(map[uint32]net.Conn),
		mutex:       sync.Mutex{},
	}
}

func (c *Connections) Add(id uint32, conn net.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Connections[id] = conn
}

func (c *Connections) Delete(id uint32) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.Connections, id)
}

// Broadcast sends a message to all connections in the Connections object (except the sender) and the output writer.
func (c *Connections) Broadcast(output io.Writer, senderID uint32, message string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, err := io.WriteString(output, message); err != nil && !errors.Is(err, io.ErrShortWrite) {
		log.Printf("failed to write to output: %v", err)
	}

	for id, conn := range c.Connections {
		if id == senderID {
			continue
		}

		if _, err := io.WriteString(conn, message); err != nil && !errors.Is(err, io.ErrShortWrite) {
			log.Printf("failed to write to connection %v: %v", id, err)
			break
		}
	}
}

// listenForConnections listens for incoming connections on a specified listener and reads data from these.
// The function reads data from the connections and writes it to the output writer.
// The first connection to be accepted is used to write data to from stdin.
// If keep open is set, the maximum number of connections is set to maxConnections else it is set to 1.
// In broker mode, the function reads from all connections and broadcasts the messages to all other connections.
// In chat mode, the function prepends the user id to the message before broadcasting.
// Arguments:
//   - output: The io.Writer object to which the function writes the data read from the connections.
//   - listener: The net.Listener object on which the function listens for incoming connections. This listener should already be initialized
//     and listening on the desired port.
func (c *cmd) listenForConnections(output io.Writer, listener net.Listener) error {
	var (
		connectionsHandled uint32
		wg                 sync.WaitGroup
		once               sync.Once
	)

	connections := NewConnections()

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

		go once.Do(func() {
			if _, err := io.Copy(conn, c.stdin); err != nil {
				log.Printf("failed to write to connection: %v", err)
			}
		})

		atomic.AddUint32(&connectionsHandled, 1)
		connectionID := atomic.LoadUint32(&connectionsHandled)
		connections.Add(connectionID, conn)

		wg.Add(1)
		go func(connections *Connections, id uint32) {
			defer wg.Done()
			defer conn.Close()
			defer func() {
				connections.Delete(id)
			}()

			// broadcast messages to all connections in broker mode
			if c.config.ListenModeOptions.BrokerMode {
				scanner := bufio.NewScanner(conn)
				for scanner.Scan() {
					line := scanner.Text()

					var formattedLine string
					// if chat-mode is enabled, prepend the user id to the message
					if c.config.ListenModeOptions.ChatMode {
						formattedLine = fmt.Sprintf("user<%d>: %s\n", id, line)
					} else {
						formattedLine = fmt.Sprintf("%s\n", line)
					}

					// broadcast the message to all connections except itseld
					connections.Broadcast(output, id, formattedLine)
				}

				if err := scanner.Err(); err != nil {
					log.Printf("failed to read from connection: %v", err)
				}
			} else {
				// without broker mode, read from the connection and write to the output
				for {
					if _, err = io.Copy(output, connections.Connections[id]); err != nil {
						if errors.Is(err, io.ErrShortWrite) {
							continue
						}

						log.Printf("failed to write to output: %v", err)
					}

					break
				}
			}
		}(connections, connectionID)
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
