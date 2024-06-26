// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/u-root/u-root/pkg/netcat"
)

func (c *cmd) connectMode(output io.Writer, network, address string) error {
	conn, err := c.establishConnection(network, address)
	if err != nil {
		return fmt.Errorf("failed to establish connection: %v", err)
	}
	log.Printf("Connection to %s [%s] succeeded", address, network)

	// Return immediately if Zero-I/O mode is enabled and connection is established
	if c.config.ConnectionModeOptions.ZeroIO {
		return nil
	}

	var wg sync.WaitGroup

	if !c.config.Misc.ReceiveOnly {
		wg.Add(1)

		go func() {
			defer wg.Done()
			c.writeToRemote(conn)
		}()

		// prepare command execution on the server
		if c.config.CommandExec.Type != netcat.EXEC_TYPE_NONE {
			if err := c.config.CommandExec.Execute(conn, io.MultiWriter(conn, output), c.stderr, c.config.Misc.EOL); err != nil {
				return fmt.Errorf("run command: %v", err)
			}
		}
	}

	// in send-only mode ignore incoming data
	if c.config.Misc.SendOnly {
		return nil
	}

	// read from the connection
	if _, err := io.Copy(output, conn); err != nil {
		return fmt.Errorf("run dump: %v", err)
	}

	wg.Wait()

	return nil
}

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

		// unsupported socket types
		case netcat.SOCKET_TYPE_SCTP, netcat.SOCKET_TYPE_VSOCK, netcat.SOCKET_TYPE_UDP_VSOCK:
			return nil, fmt.Errorf("currently unsupported socket type %q", c.config.ProtocolOptions.SocketType)

		case netcat.SOCKET_TYPE_NONE:
		default:
			return nil, fmt.Errorf("undefined socket type %q", c.config.ProtocolOptions.SocketType)
		}
	}

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

	if c.config.Timing.Timeout > 0 {
		conn.SetDeadline(time.Now().Add(c.config.Timing.Timeout))
	}

	return conn, nil
}

func (c *cmd) writeToRemote(conn io.Writer) {
	eolReader := netcat.NewEOLReader(c.stdin, c.config.Misc.EOL)

	// If the delay is set, read the input line by line in time intervals of the delay duration
	if c.config.Timing.Delay > 0 {
		scanner := bufio.NewScanner(eolReader)
		scanner.Split(ScanWithCustomEOL(string(c.config.Misc.EOL)))

		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Printf("failed reading input: %v", err)
			return
		}
		for i, line := range lines {
			// Determine if this is the last line
			isLastLine := (i == len(lines)-1)

			// Write the line
			if _, err := conn.Write([]byte(line + string(c.config.Misc.EOL))); err != nil {
				log.Printf("failed writing to host: %v", err)
			}

			// Apply the delay if configured
			if !isLastLine { // Avoid delay after the last line
				time.Sleep(c.config.Timing.Delay)
			}
		}
	} else {
		if _, err := io.Copy(conn, eolReader); err != nil {
			log.Printf("failed writing to host: %v", err)
		}
	}

	// do not shutdown the connection if the no-shutdown flag is set
	if c.config.Misc.NoShutdown {
		for {
			time.Sleep(1 * time.Hour)
		}
	}
}

// Custom split function to use with bufio.Scanner
func ScanWithCustomEOL(eol string) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := bytes.Index(data, []byte(eol)); i >= 0 {
			return i + len(eol), data[0:i], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	}
}
