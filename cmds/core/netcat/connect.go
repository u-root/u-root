// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/proxy"

	"github.com/u-root/u-root/pkg/netcat"
)

func (c *cmd) connectMode(output io.Writer, network, address string) error {
	if c.config.ConnectionModeOptions.ScanPorts && !c.config.ConnectionModeOptions.ZeroIO {
		return fmt.Errorf("scanning ports is only supported in Zero-I/O mode")
	}

	if c.config.ConnectionModeOptions.ScanPorts {
		return c.scanPorts()
	}

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
	for {
		if _, err := io.Copy(output, conn); err != nil {
			if errors.Is(err, io.ErrShortWrite) {
				continue
			}

			return fmt.Errorf("failed to write: %v", err)
		}

		break
	}

	wg.Wait()

	return nil
}

func (c *cmd) scanPorts() error {
	for {
		if c.config.ConnectionModeOptions.CurrentPort > c.config.ConnectionModeOptions.EndPort {
			return nil
		}

		network, address, err := c.connection()
		if err != nil {
			return fmt.Errorf("failed to parse connection: %v", err)
		}

		_, err = c.establishConnection(network, address)
		if err != nil {
			log.Printf("connect to %v port %v (%v) failed: %v", c.config.Host, c.config.ConnectionModeOptions.CurrentPort, c.config.ProtocolOptions.SocketType, err)
			c.config.ConnectionModeOptions.CurrentPort++
			continue
		}

		log.Printf("connect to %v port %v (%v) succeeded", c.config.Host, c.config.ConnectionModeOptions.CurrentPort, c.config.ProtocolOptions.SocketType)
		c.config.ConnectionModeOptions.CurrentPort++
	}
}

func (c *cmd) proxyDialer(dialer proxy.Dialer) (proxy.Dialer, error) {
	var proxyAuth string
	if c.config.ProxyConfig.Auth != "" {
		proxyAuth = c.config.ProxyConfig.Auth + "@"
	}

	proxyAddr := fmt.Sprintf("%v://%v%v", c.config.ProxyConfig.Type, proxyAuth, c.config.ProxyConfig.Address)
	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %v", err)
	}

	return proxy.FromURL(proxyURL, dialer)
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
