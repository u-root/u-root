// Copyright 2012-2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

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
func (c *cmd) listenForConnections(output io.WriteCloser, listener net.Listener) error {
	var (
		connectionsHandled uint32
		wg                 sync.WaitGroup
		once               sync.Once
	)

	connections := NewConnections()

	maxConnections := c.config.ListenModeOptions.MaxConnections

	for atomic.LoadUint32(&connectionsHandled) < maxConnections {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("wait for connection: %v", err)
			continue
		}

		if c.config.ProtocolOptions.SocketType == netcat.SOCKET_TYPE_TCP {
			if !c.config.AccessControl.IsAllowed(parseRemoteAddr(c.config.ProtocolOptions.SocketType, conn.RemoteAddr().String())) {
				defer conn.Close()
				break
			}
		}

		go once.Do(func() {
			if _, err := io.Copy(conn, c.stdin); err != nil {
				log.Printf("write to connection: %v", err)
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
					log.Printf("read from connection: %v", err)
				}
			} else {
				// without broker mode, read from the connection and write to the output
				for {
					connections.mutex.Lock()
					connection := connections.Connections[id]
					connections.mutex.Unlock()

					if _, err = io.Copy(output, connection); err != nil {
						if errors.Is(err, io.ErrShortWrite) {
							continue
						}

						log.Printf("output: %v", err)
					}

					break
				}
			}
		}(connections, connectionID)
	}

	wg.Wait() // Wait for all connections to finish

	return nil
}
