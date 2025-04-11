// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/net/proxy"

	"github.com/mdlayher/vsock"
	"github.com/u-root/u-root/pkg/netcat"
)

var osConnectors = map[netcat.SocketType]func(string, string) (net.Conn, error){}

type connectFn func(output io.WriteCloser, network, address string) error

func (c *cmd) connectMode(output io.WriteCloser, network, address string, connect connectFn) error {
	if network == "tcp" || network == "udp" {
		err4 := connect(output, network+"4", address)
		if err4 == nil {
			return nil
		}
		err6 := connect(output, network+"6", address)
		if err6 == nil {
			return nil
		}
		return fmt.Errorf("connect mode: %w", errors.Join(err4, err6))
	}

	return connect(output, network, address)
}

func (c *cmd) connect(output io.WriteCloser, network, address string) error {
	if c.config.ConnectionModeOptions.ScanPorts && !c.config.ConnectionModeOptions.ZeroIO {
		return fmt.Errorf("scanning ports is only supported in Zero-I/O mode")
	}

	if c.config.ConnectionModeOptions.ScanPorts {
		return c.scanPorts()
	}

	if c.config.ProtocolOptions.SocketType == netcat.SOCKET_TYPE_UDP_UNIX &&
		c.config.ConnectionModeOptions.SourceHost == "" {
		c.config.ConnectionModeOptions.SourceHost = filepath.Join(os.TempDir(), fmt.Sprintf("netcat.%x.sock", rand.Uint64()))
		defer os.Remove(c.config.ConnectionModeOptions.SourceHost)
	}

	conn, err := c.establishConnection(network, address)
	if err != nil {
		return fmt.Errorf("failed to establish connection: %w", err)
	}
	defer conn.Close()

	log.Printf("connected to %s [%s]", conn.RemoteAddr(), network)

	// Return immediately if Zero-I/O mode is enabled and connection is established
	if c.config.ConnectionModeOptions.ZeroIO {
		return nil
	}

	if c.config.ProtocolOptions.SocketType == netcat.SOCKET_TYPE_UDP ||
		c.config.ProtocolOptions.SocketType == netcat.SOCKET_TYPE_UDP_UNIX {
		return c.transferPackets(output, conn.(net.PacketConn), false)
	}

	var wg sync.WaitGroup

	if !c.config.Misc.ReceiveOnly {
		// prepare command execution on the server
		if c.config.CommandExec.Type != netcat.EXEC_TYPE_NONE {
			if err := c.config.CommandExec.Execute(conn, c.stderr, c.config.Misc.EOL); err != nil {
				return fmt.Errorf("run command: %w", err)
			}

			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			c.writeToRemote(conn)
		}()
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

			return fmt.Errorf("failed to write: %w", err)
		}

		if err = output.Close(); err != nil {
			log.Printf("failed to close output: %v", err)
		}

		break
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
			dialer.LocalAddr, err = net.ResolveTCPAddr(network, net.JoinHostPort(c.config.ConnectionModeOptions.SourceHost, c.config.ConnectionModeOptions.SourcePort))
			if err != nil {
				return nil, fmt.Errorf("failed to resolve source address %w", err)
			}

		case netcat.SOCKET_TYPE_UDP:
			dialer.LocalAddr, err = net.ResolveUDPAddr(network, net.JoinHostPort(c.config.ConnectionModeOptions.SourceHost, c.config.ConnectionModeOptions.SourcePort))
			if err != nil {
				return nil, fmt.Errorf("failed to resolve source address %w", err)
			}

		case netcat.SOCKET_TYPE_UNIX:
			dialer.LocalAddr, err = net.ResolveUnixAddr(network, c.config.ConnectionModeOptions.SourceHost)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve source address %w", err)
			}

		case netcat.SOCKET_TYPE_UDP_UNIX:
			dialer.LocalAddr, err = net.ResolveUnixAddr(network, c.config.ConnectionModeOptions.SourceHost)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve source address %w", err)
			}

		case netcat.SOCKET_TYPE_VSOCK:
			cid, port, err := netcat.SplitVSockAddr(address)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve VSOCK address: %w", err)
			}

			return vsock.Dial(cid, port, nil)

		// unsupported socket types
		default:
			osConn, ok := osConnectors[c.config.ProtocolOptions.SocketType]
			if !ok {
				return nil, fmt.Errorf("socket type %q:%w", c.config.ProtocolOptions.SocketType, os.ErrNotExist)
			}
			return osConn(network, address)
		}
	}

	// Proxy Support
	if c.config.ProxyConfig.Enabled {
		proxyDialer, err := c.proxyDialer(dialer)
		if err != nil {
			return nil, err
		}

		conn, err = proxyDialer.Dial(network, address)
		if err != nil {
			return nil, err
		}
	} else {
		// TLS Support
		if c.config.SSLConfig.Enabled || c.config.SSLConfig.VerifyTrust {
			tlsConfig, err := c.config.SSLConfig.GenerateTLSConfiguration(false)
			if err != nil {
				return nil, err
			}

			conn, err = tls.DialWithDialer(dialer, network, address, tlsConfig)
			if err != nil {
				return nil, err
			}
		} else {
			conn, err = dialer.Dial(network, address)
			if err != nil {
				return nil, err
			}
		}
	}

	if c.config.Timing.Timeout > 0 {
		conn.SetDeadline(time.Now().Add(c.config.Timing.Timeout))
	}

	return conn, nil
}

func (c *cmd) scanPorts() error {
	for {
		if c.config.ConnectionModeOptions.CurrentPort > c.config.ConnectionModeOptions.EndPort {
			return nil
		}

		network, address, err := c.connection()
		if err != nil {
			return fmt.Errorf("failed to parse connection: %w", err)
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
		return nil, fmt.Errorf("invalid proxy URL: %w", err)
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
	if !c.config.Misc.NoShutdown {
		if err := netcat.CloseWrite(conn); err != nil {
			log.Printf("failed to shut down socket for writing: %v", err)
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
