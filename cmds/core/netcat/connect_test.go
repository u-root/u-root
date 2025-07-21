// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/netcat"
)

type closableDiscard struct{}

func (cd *closableDiscard) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (cd *closableDiscard) Close() error {
	return nil
}

func TestConnectMode(t *testing.T) {
	connectErr := errors.New("failed to connect")

	tests := []struct {
		name        string
		network     string
		address     string
		connectFunc func(output io.WriteCloser, network, address string) error
		err         error
	}{
		{
			name:    "TCPv4 success",
			network: "tcp",
			address: "localhost:8080",
			connectFunc: func(output io.WriteCloser, network, address string) error {
				if network == "tcp4" {
					return nil
				}
				return connectErr
			},
			err: nil,
		},
		{
			name:    "TCPv6 success",
			network: "tcp",
			address: "localhost:8080",
			connectFunc: func(output io.WriteCloser, network, address string) error {
				if network == "tcp6" {
					return nil
				}
				return connectErr
			},
			err: nil,
		},
		{
			name:    "UDPv4 success",
			network: "udp",
			address: "localhost:8080",
			connectFunc: func(output io.WriteCloser, network, address string) error {
				if network == "udp4" {
					return nil
				}
				return connectErr
			},
			err: nil,
		},
		{
			name:    "UDPv6 success",
			network: "udp",
			address: "localhost:8080",
			connectFunc: func(output io.WriteCloser, network, address string) error {
				if network == "udp6" {
					return nil
				}
				return connectErr
			},
			err: nil,
		},
		{
			name:    "TCPv4 and TCPv6 failure",
			network: "tcp",
			address: "localhost:8080",
			connectFunc: func(output io.WriteCloser, network, address string) error {
				return connectErr
			},
			err: connectErr,
		},
		{
			name:    "UDPv4 and UDPv6 failure",
			network: "udp",
			address: "localhost:8080",
			connectFunc: func(output io.WriteCloser, network, address string) error {
				return connectErr
			},
			err: connectErr,
		},
		{
			name:    "Other network success",
			network: "unix",
			address: "/tmp/socket",
			connectFunc: func(output io.WriteCloser, network, address string) error {
				if network == "unix" {
					return nil
				}
				return connectErr
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmd{}
			err := cmd.connectMode(&closableDiscard{}, tt.network, tt.address, tt.connectFunc)
			if !errors.Is(err, tt.err) {
				t.Errorf("got %v, want %v", err, tt.err)
			}
		})
	}
}

// Mock for the net.Conn interface
type mockConn struct {
	net.Conn
	read  func(b []byte) (n int, err error)
	write func(b []byte) (n int, err error)
	close func() error
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.read(b)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.write(b)
}

func (m *mockConn) Close() error {
	return m.close()
}

func TestConnect(t *testing.T) {
	response := "World"
	tests := []struct {
		name        string
		address     string
		stdin       string
		stderr      string
		stdout      string
		config      *netcat.Config // Assuming Config is the type of c.config
		expectError bool
	}{
		{
			name:    "zero I/O",
			address: "127.0.0.1:8080",
			config: &netcat.Config{
				ProtocolOptions:       netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_TCP},
				CommandExec:           netcat.Exec{Type: netcat.EXEC_TYPE_NONE},
				ConnectionModeOptions: netcat.ConnectModeOptions{ZeroIO: true},
			},
		},
		{
			name:    "successful connection",
			address: "127.0.0.2:8080",
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_TCP},
				CommandExec:     netcat.Exec{Type: netcat.EXEC_TYPE_NONE},
			},
			stdout: response,
		},
		{
			name:    "successful connection with send only",
			address: "127.0.0.3:8080",
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_TCP},
				CommandExec:     netcat.Exec{Type: netcat.EXEC_TYPE_NONE},
				Misc:            netcat.MiscOptions{SendOnly: true},
			},
			stdout: "",
		},
		{
			name:    "successful connection with receive only",
			address: "127.0.0.4:8080",
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_TCP},
				CommandExec:     netcat.Exec{Type: netcat.EXEC_TYPE_NONE},
				Misc:            netcat.MiscOptions{ReceiveOnly: true},
			},
			stdout: response,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup

			l, err := net.Listen("tcp", tt.address)
			if err != nil {
				t.Fatal(err)
			}

			wg.Add(1)
			time.AfterFunc(500*time.Millisecond, func() {
				defer wg.Done()
				l.Close()
			})

			wg.Add(1)
			go func() {
				defer wg.Done()
				conn, err := l.Accept()
				if err != nil {
					return
				}
				conn.Write([]byte(response))
				conn.Close()
			}()

			stdin := strings.NewReader(tt.stdin)

			c := &cmd{
				stdin:  stdin,
				config: tt.config,
			}

			var output closableBuffer
			err = c.connect(&output, "tcp", "127.0.0.1:8080")
			if err != nil {
				if !tt.expectError {
					return
				}
				t.Errorf("Expected no error, got %v", err)
			}

			if output.String() != tt.stdout {
				t.Errorf("Expected %q, got %q", tt.stdout, output.String())
			}

			wg.Wait()
		})
	}
}

func TestEstablishConnection(t *testing.T) {
	addr := "localhost:3000"
	addr6 := "[::1]:3000"

	tests := []struct {
		name        string
		network     string
		address     string
		config      *netcat.Config
		expectError bool
	}{
		{
			name:    "Successful TCP connection",
			network: "tcp",
			address: addr,
			config: &netcat.Config{
				ConnectionModeOptions: netcat.ConnectModeOptions{
					SourceHost: "localhost",
					SourcePort: "8081",
				},
				ProtocolOptions: netcat.ProtocolOptions{
					SocketType: netcat.SOCKET_TYPE_TCP,
				},
				Timing: netcat.TimingOptions{
					Wait:    5 * time.Second,
					Timeout: 5 * time.Second,
				},
			},

			expectError: false,
		},
		{
			name:    "Successful TCPv6 connection",
			network: "tcp6",
			address: addr6,
			config: &netcat.Config{
				ConnectionModeOptions: netcat.ConnectModeOptions{
					SourceHost: "::1",
					SourcePort: "8081",
				},
				ProtocolOptions: netcat.ProtocolOptions{
					SocketType: netcat.SOCKET_TYPE_TCP,
				},
				Timing: netcat.TimingOptions{
					Wait:    5 * time.Second,
					Timeout: 5 * time.Second,
				},
			},

			expectError: false,
		},
		{
			name:    "Unsuccessful TCP connection",
			network: "tcp",
			address: "localhost:3030",
			config: &netcat.Config{
				ConnectionModeOptions: netcat.ConnectModeOptions{
					SourceHost: "localhost",
					SourcePort: "8081",
				},
				ProtocolOptions: netcat.ProtocolOptions{
					SocketType: netcat.SOCKET_TYPE_TCP,
				},
				Timing: netcat.TimingOptions{
					Wait:    5 * time.Second,
					Timeout: 5 * time.Second,
				},
			},

			expectError: true,
		},
		{
			name:    "Successful UDP connection",
			network: "udp",
			address: addr,
			config: &netcat.Config{
				ConnectionModeOptions: netcat.ConnectModeOptions{
					SourceHost: "localhost",
					SourcePort: "8081",
				},
				ProtocolOptions: netcat.ProtocolOptions{
					SocketType: netcat.SOCKET_TYPE_UDP,
				},
				Timing: netcat.TimingOptions{
					Wait:    5 * time.Second,
					Timeout: 5 * time.Second,
				},
			},
			expectError: false,
		},
		{
			name:    "Successful UDPv6 connection",
			network: "udp6",
			address: addr6,
			config: &netcat.Config{
				ConnectionModeOptions: netcat.ConnectModeOptions{
					SourceHost: "::1",
					SourcePort: "8081",
				},
				ProtocolOptions: netcat.ProtocolOptions{
					SocketType: netcat.SOCKET_TYPE_UDP,
				},
				Timing: netcat.TimingOptions{
					Wait:    5 * time.Second,
					Timeout: 5 * time.Second,
				},
			},
			expectError: false,
		},
		{
			name:    "unimplemented socket",
			network: "unix",
			address: "localhost:3077",
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{
					SocketType: netcat.SOCKET_TYPE_VSOCK,
				},
			},
			expectError: true,
		},
		{
			name:    "none socket",
			network: "unix",
			address: "localhost:3077",
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{
					SocketType: netcat.SOCKET_TYPE_NONE,
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var listenAddr string
			var wg sync.WaitGroup

			// github does not enable IPv6 for docker containers
			if tt.network == "tcp6" || tt.network == "udp6" {
				_, disable_ipv6 := os.LookupEnv("NETCAT_CONNECT_TEST_DISABLE_IPV6")
				if disable_ipv6 {
					log.Printf("skipping %s", tt.name)
					return
				}
			}

			switch tt.network {
			case "tcp", "tcp6":
				if tt.network == "tcp" {
					listenAddr = addr
				} else {
					listenAddr = addr6
				}
				l, err := net.Listen(tt.network, listenAddr)
				if err != nil {
					t.Fatal(err)
				}

				wg.Add(1)
				time.AfterFunc(500*time.Millisecond, func() {
					defer wg.Done()
					l.Close()
				})

				wg.Add(1)
				go func() {
					defer wg.Done()
					conn, err := l.Accept()
					if err != nil {
						return
					}

					defer conn.Close()
				}()
			case "udp", "udp6":
				if tt.network == "udp" {
					listenAddr = addr
				} else {
					listenAddr = addr6
				}
				l, err := net.ListenPacket(tt.network, listenAddr)
				if err != nil {
					t.Fatal(err)
				}

				time.AfterFunc(500*time.Millisecond, func() {
					l.Close()
				})
			}

			c := &cmd{
				config: tt.config,
			}

			conn, err := c.establishConnection(tt.network, tt.address)

			if tt.expectError && err != nil {
				return
			}

			if conn == nil {
				t.Errorf("Expected a connection, got nil")
			}
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			wg.Wait()
		})
	}
}

func TestEstablishConnectionUnix(t *testing.T) {
	socketPath := "/tmp/uroot_test_unix"
	sourcePath := "/tmp/uroot_test_unix_source"

	tests := []struct {
		name        string
		network     string
		address     string
		remove      bool
		config      *netcat.Config
		expectError bool
	}{
		{
			name:    "Successful Unix connection",
			network: "unix",
			address: socketPath,
			config: &netcat.Config{
				ConnectionModeOptions: netcat.ConnectModeOptions{
					SourceHost: sourcePath,
				},
				ProtocolOptions: netcat.ProtocolOptions{
					SocketType: netcat.SOCKET_TYPE_UNIX,
				},
				Timing: netcat.TimingOptions{
					Wait:    5 * time.Second,
					Timeout: 5 * time.Second,
				},
			},

			expectError: false,
		},
		{
			name:    "Successful Unix connection (unnamed client socket)",
			network: "unix",
			address: socketPath,
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{
					SocketType: netcat.SOCKET_TYPE_UNIX,
				},
				Timing: netcat.TimingOptions{
					Wait:    5 * time.Second,
					Timeout: 5 * time.Second,
				},
			},

			expectError: false,
		},
		{
			name:    "Successful UDP Unix connection",
			network: "unixgram",
			address: socketPath,
			config: &netcat.Config{
				ConnectionModeOptions: netcat.ConnectModeOptions{
					SourceHost: sourcePath,
				},
				ProtocolOptions: netcat.ProtocolOptions{
					SocketType: netcat.SOCKET_TYPE_UDP_UNIX,
				},
				Timing: netcat.TimingOptions{
					Wait:    5 * time.Second,
					Timeout: 5 * time.Second,
				},
			},

			expectError: false,
		},
		{
			name:    "Successful UDP Unix connection  (temporary client socket)",
			network: "unixgram",
			address: socketPath,
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{
					SocketType: netcat.SOCKET_TYPE_UDP_UNIX,
				},
				Timing: netcat.TimingOptions{
					Wait:    5 * time.Second,
					Timeout: 5 * time.Second,
				},
			},

			expectError: false,
		},
		{
			name:    "Unsuccessful Unix connection",
			network: "unix",
			address: "/tmp/not_available",
			config: &netcat.Config{
				ConnectionModeOptions: netcat.ConnectModeOptions{
					SourceHost: sourcePath,
				},
				ProtocolOptions: netcat.ProtocolOptions{
					SocketType: netcat.SOCKET_TYPE_UNIX,
				},
				Timing: netcat.TimingOptions{
					Wait:    5 * time.Second,
					Timeout: 5 * time.Second,
				},
			},

			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := filepath.Glob("/tmp/uroot_test_unix*")
			if err != nil {
				t.Fatalf("Failed to find files: %v", err)
			}

			for _, file := range files {
				if err := os.Remove(file); err != nil {
					t.Fatalf("Failed to remove file %s: %v", file, err)
				}
			}

			var wg sync.WaitGroup

			defer os.Remove(socketPath)

			switch tt.network {
			case "unix":
				unixL, err := net.Listen(tt.network, socketPath)
				if err != nil {
					t.Fatal(err)
				}

				time.AfterFunc(500*time.Millisecond, func() {
					unixL.Close()
				})

				wg.Add(1)
				go func() {
					defer wg.Done()
					conn, err := unixL.Accept()
					if err != nil {
						return
					}
					defer conn.Close()
				}()
			case "unixgram":
				unixL, err := net.ListenPacket(tt.network, socketPath)
				if err != nil {
					t.Fatal(err)
				}

				time.AfterFunc(500*time.Millisecond, func() {
					unixL.Close()
				})
			}

			c := &cmd{
				config: tt.config,
			}

			conn, err := c.establishConnection(tt.network, tt.address)

			if tt.expectError && err != nil {
				return
			}

			if conn == nil {
				t.Errorf("Expected a connection, got nil")
			}
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if tt.remove {
				if err := os.Remove(tt.address); err != nil {
					t.Errorf("Failed to remove file: %v", err)
				}
			}

			wg.Wait()
		})
	}
}

// mockWriter captures writes for verification
type mockWriter struct {
	bytes.Buffer
}

func (mw *mockWriter) Write(p []byte) (n int, err error) {
	return mw.Buffer.Write(p)
}

func TestWriteToRemote(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		delay      time.Duration
		eol        []byte
		expected   string
		expectHang bool
	}{
		{
			name:     "No delay",
			input:    "hello\nworld",
			eol:      []byte("\n"),
			delay:    0,
			expected: "hello\nworld",
		},
		{
			name:     "With CRLF",
			input:    "hello\nworld\n",
			eol:      []byte("\r\n"),
			delay:    0,
			expected: "hello\r\nworld\r\n",
		},
		{
			name:     "With delay",
			input:    "hello\nworld\n",
			eol:      []byte("\n"),
			delay:    10 * time.Millisecond,
			expected: "hello\nworld\n",
		},
		{
			name:       "long delay",
			input:      "hello\nworld",
			delay:      500 * time.Millisecond,
			expectHang: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConn := &mockWriter{}
			cmd := &cmd{
				stdin: strings.NewReader(tt.input),
				config: &netcat.Config{
					Misc: netcat.MiscOptions{
						EOL: tt.eol,
					},
					Timing: netcat.TimingOptions{
						Delay: tt.delay,
					},
				},
			}

			done := make(chan bool)
			go func() {
				cmd.writeToRemote(mockConn)
				done <- true
			}()

			select {
			case <-done:
				if tt.expectHang {
					t.Errorf("Expected writeToRemote to hang, but it did not")
				}
			case <-time.After(100 * time.Millisecond):
				if !tt.expectHang {
					t.Errorf("writeToRemote took too long to complete")
				}
			}

			if !tt.expectHang && mockConn.String() != tt.expected {
				t.Errorf("Expected output %q, got %q", tt.expected, mockConn.String())
			}
		})
	}
}

func TestScanWithCustomEOL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		eol      string
		expected []string
	}{
		{
			name:     "Single custom EOL",
			input:    "Hello, world!#EOL#This is a test",
			eol:      "#EOL#",
			expected: []string{"Hello, world!", "This is a test"},
		},
		{
			name:     "Multiple custom EOL",
			input:    "Line 1#EOL#Line 2#EOL#Line 3",
			eol:      "#EOL#",
			expected: []string{"Line 1", "Line 2", "Line 3"},
		},
		{
			name:     "No custom EOL",
			input:    "No custom EOL here",
			eol:      "#EOL#",
			expected: []string{"No custom EOL here"},
		},
		{
			name:     "Custom EOL at the end",
			input:    "Ends with EOL#EOL#",
			eol:      "#EOL#",
			expected: []string{"Ends with EOL"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(bytes.NewBufferString(tt.input))
			scanner.Split(ScanWithCustomEOL(tt.eol))

			var got []string
			for scanner.Scan() {
				got = append(got, scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				t.Errorf("Scanner encountered an error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Got %v, want %v", got, tt.expected)
			}
		})
	}
}

// mockDialer is a mock implementation of proxy.Dialer for testing purposes.
type mockDialer struct{}

func (m *mockDialer) Dial(network, addr string) (conn net.Conn, err error) {
	// Mock implementation
	return nil, nil
}

func TestProxyDialer(t *testing.T) {
	// Setup
	mockDial := &mockDialer{}
	testCmd := &cmd{
		config: &netcat.Config{
			ProxyConfig: netcat.ProxyOptions{
				Auth:    "user:pass",
				Type:    netcat.PROXY_TYPE_SOCKS5,
				Address: "proxy.example.com:8080",
			},
		},
	}

	// Execute
	dialer, err := testCmd.proxyDialer(mockDial)
	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if dialer == nil {
		t.Error("Expected dialer to be not nil")
	}
}
