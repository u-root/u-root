package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/u-root/u-root/pkg/netcat"
)

func (c *cmd) connectMode(output io.Writer, network, address string) error {
	var (
		err  error
		conn net.Conn
	)

	dialer := &net.Dialer{
		Timeout: c.config.Timing.Wait,
	}

	switch c.config.ProtocolOptions.SocketType {

	case netcat.SOCKET_TYPE_TCP:
		dialer.LocalAddr, err = net.ResolveTCPAddr(c.config.ConnectionModeOptions.SourceHost, c.config.ConnectionModeOptions.SourcePort)
		if err != nil {
			return fmt.Errorf("connection: failed to resolve source address %v", err)
		}

	case netcat.SOCKET_TYPE_UDP:
		dialer.LocalAddr, err = net.ResolveUDPAddr(c.config.ConnectionModeOptions.SourceHost, c.config.ConnectionModeOptions.SourcePort)
		if err != nil {
			return fmt.Errorf("connection: failed to resolve source address %v", err)
		}

	case netcat.SOCKET_TYPE_UNIX:
		dialer.LocalAddr, err = net.ResolveUnixAddr(c.config.ConnectionModeOptions.SourceHost, c.config.ConnectionModeOptions.SourcePort)
		if err != nil {
			return fmt.Errorf("connection: failed to resolve source address %v", err)
		}

	case netcat.SOCKET_TYPE_UDP_UNIX:
		dialer.LocalAddr, err = net.ResolveUnixAddr(c.config.ConnectionModeOptions.SourceHost, c.config.ConnectionModeOptions.SourcePort)
		if err != nil {
			return fmt.Errorf("connection: failed to resolve source address %v", err)
		}

	// unsupported socket types
	case netcat.SOCKET_TYPE_SCTP, netcat.SOCKET_TYPE_VSOCK, netcat.SOCKET_TYPE_UDP_VSOCK:
		return fmt.Errorf("currently unsupported socket type %q", c.config.ProtocolOptions.SocketType)

	case netcat.SOCKET_TYPE_NONE:
	default:
		return fmt.Errorf("undefined socket type %q", c.config.ProtocolOptions.SocketType)
	}

	// TLS Support
	if c.config.SSLConfig.Enabled || c.config.SSLConfig.VerifyTrust {
		tlsConfig, err := c.generateTLSConfiguration()
		if err != nil {
			return fmt.Errorf("connection: %v", err)
		}

		conn, err = tls.DialWithDialer(dialer, network, address, tlsConfig)
		if err != nil {
			return fmt.Errorf("connection: %v", err)
		}
	} else {
		conn, err = dialer.Dial(network, address)
		if err != nil {
			return fmt.Errorf("connection: %v", err)
		}
	}

	if c.config.Timing.Timeout > 0 {
		conn.SetDeadline(time.Now().Add(c.config.Timing.Timeout))
	}

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

func (c *cmd) writeToRemote(conn io.Writer) {
	eolReader := netcat.NewEOLReader(c.stdin, c.config.Misc.EOL)

	// If the delay is set, read the input line by line in time intervals of the delay duration
	if c.config.Timing.Delay > 0 {
		scanner := bufio.NewScanner(eolReader)

		for scanner.Scan() {
			if _, err := conn.Write([]byte(scanner.Text() + "\n")); err != nil {
				netcat.FLogf(c.config, c.stderr, "run copy: %v", err)
			}

			time.Sleep(c.config.Timing.Delay)
		}

		if err := scanner.Err(); err != nil {
			netcat.FLogf(c.config, c.stderr, "run copy: %v", err)
		}
	} else {
		if _, err := io.Copy(conn, eolReader); err != nil {
			netcat.FLogf(c.config, c.stderr, "run copy: %v", err)
		}
	}

	// do not shutdown the connection if the no-shutdown flag is set
	if c.config.Misc.NoShutdown {
		for {
			time.Sleep(1 * time.Hour)
		}
	}
}
