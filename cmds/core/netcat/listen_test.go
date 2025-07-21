// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"reflect"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/u-root/u-root/pkg/netcat"
)

func TestListenMode(t *testing.T) {
	listenErr := errors.New("failed to listen")

	tests := []struct {
		name       string
		network    string
		address    string
		listenFunc func(output io.WriteCloser, network, address string) error
		err        error
	}{
		{
			name:    "TCPv4 success",
			network: "tcp",
			address: "localhost:8080",
			listenFunc: func(output io.WriteCloser, network, address string) error {
				if network == "tcp4" {
					return nil
				}
				return listenErr
			},
			err: nil,
		},
		{
			name:    "TCPv6 success",
			network: "tcp",
			address: "localhost:8080",
			listenFunc: func(output io.WriteCloser, network, address string) error {
				if network == "tcp6" {
					return nil
				}
				return listenErr
			},
			err: nil,
		},
		{
			name:    "UDPv4 success",
			network: "udp",
			address: "localhost:8080",
			listenFunc: func(output io.WriteCloser, network, address string) error {
				if network == "udp4" {
					return nil
				}
				return listenErr
			},
			err: nil,
		},
		{
			name:    "UDPv6 success",
			network: "udp",
			address: "localhost:8080",
			listenFunc: func(output io.WriteCloser, network, address string) error {
				if network == "udp6" {
					return nil
				}
				return listenErr
			},
			err: nil,
		},
		{
			name:    "TCPv4 and TCPv6 failure",
			network: "tcp",
			address: "localhost:8080",
			listenFunc: func(output io.WriteCloser, network, address string) error {
				return listenErr
			},
			err: listenErr,
		},
		{
			name:    "UDPv4 and UDPv6 failure",
			network: "udp",
			address: "localhost:8080",
			listenFunc: func(output io.WriteCloser, network, address string) error {
				return listenErr
			},
			err: listenErr,
		},
		{
			name:    "Other network success",
			network: "unix",
			address: "/tmp/socket",
			listenFunc: func(output io.WriteCloser, network, address string) error {
				if network == "unix" {
					return nil
				}
				return listenErr
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cmd{}
			err := cmd.listenMode(&closableDiscard{}, tt.network, tt.address, tt.listenFunc)
			if !errors.Is(err, tt.err) {
				t.Errorf("got %v, want %v", err, tt.err)
			}
		})
	}
}

// Define a struct for your test cases
type setupListenerTestCase struct {
	name     string
	config   *netcat.Config // Assuming Config is the type of c.config
	network  string
	address  string
	wantErr  bool
	wantType any
}

// Example test cases
func TestSetupListener(t *testing.T) {
	tests := []setupListenerTestCase{
		{
			name: "TCP without SSL",
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_TCP},
			},
			network:  "tcp",
			address:  "127.0.0.1:0",
			wantErr:  false,
			wantType: &net.TCPListener{},
		},
		{
			name: "Unsupported Socket Type",
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_UDP_VSOCK},
			},
			network: "sctp",
			address: "127.0.0.1:0",
			wantErr: true,
		},
		{
			name: "UDP without SSL",
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_UDP},
				SSLConfig:       netcat.SSLOptions{Enabled: false},
			},
			network:  "udp",
			address:  "127.0.0.1:0",
			wantErr:  false,
			wantType: &netcat.UDPListener{},
		},
		{
			name: "NoDNS set",
			config: &netcat.Config{
				Misc:            netcat.MiscOptions{NoDNS: true},
				ProtocolOptions: netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_TCP},
				SSLConfig:       netcat.SSLOptions{Enabled: false},
			},
			wantErr: true,
		},
		{
			name: "SourceHost set ",
			config: &netcat.Config{
				ConnectionModeOptions: netcat.ConnectModeOptions{SourceHost: "192.168.1.1"},
				ProtocolOptions:       netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_TCP},
			},
			wantErr: true,
		},
		{
			name: "SourcePort set",
			config: &netcat.Config{
				ConnectionModeOptions: netcat.ConnectModeOptions{SourcePort: "8080"},
				ProtocolOptions:       netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_TCP},
			},
			wantErr: true,
		},
		{
			name: "Unsupported Socket Type None",
			config: &netcat.Config{
				ProtocolOptions: netcat.ProtocolOptions{SocketType: netcat.SOCKET_TYPE_NONE},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := &cmd{config: tc.config}
			listener, err := c.setupListener(tc.network, tc.address)
			if (err != nil) != tc.wantErr {
				t.Errorf("setupListener() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if err == nil {
				defer listener.Close()
				if reflect.TypeOf(listener) != reflect.TypeOf(tc.wantType) {
					t.Errorf("Expected listener type %v, got %v", reflect.TypeOf(tc.wantType), reflect.TypeOf(listener))
				}
			}
		})
	}
}

// MockListener is a simple mock for net.Listener
type MockListener struct {
	AcceptFunc func() (net.Conn, error)
	CloseFunc  func() error
}

func (m *MockListener) Accept() (net.Conn, error) {
	return m.AcceptFunc()
}

func (m *MockListener) Close() error {
	return m.CloseFunc()
}

func (m *MockListener) Addr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
}

// MockConn is a simple mock for net.Conn
type MockConn struct {
	Output               bytes.Buffer
	ReadFunc             func([]byte) (int, error)
	CloseFunc            func() error
	LocalAddrFunc        func() net.Addr
	RemoteAddrFunc       func() net.Addr
	SetDeadlineFunc      func(t int64) error
	SetReadDeadlineFunc  func(t int64) error
	SetWriteDeadlineFunc func(t int64) error
}

func (m *MockConn) Close() error {
	return m.CloseFunc()
}

func (m *MockConn) RemoteAddr() net.Addr {
	return m.RemoteAddrFunc()
}

func (m *MockConn) LocalAddr() net.Addr {
	return m.LocalAddrFunc()
}

func (m *MockConn) SetDeadline(t time.Time) error {
	return m.SetDeadlineFunc(t.Unix())
}

func (m *MockConn) SetReadDeadline(t time.Time) error {
	return m.SetReadDeadlineFunc(t.Unix())
}

func (m *MockConn) SetWriteDeadline(t time.Time) error {
	return m.SetWriteDeadlineFunc(t.Unix())
}

func (m *MockConn) Read(b []byte) (int, error) {
	return m.ReadFunc(b)
}

func (m *MockConn) Write(b []byte) (int, error) {
	return m.Output.Write(b)
}

// Do not set testLimit higher than the number of mockConns or else the function will run indefinitely waiting for more connections
func TestListenForConnections(t *testing.T) {
	partialData := "partial data"
	tests := []struct {
		name           string
		mockConns      []*MockConn
		config         *netcat.Config
		testLimit      uint32
		expectedOutput []string
		expectError    bool
	}{
		{
			name: "Successful read",
			mockConns: []*MockConn{
				{
					ReadFunc: func(b []byte) (int, error) {
						return 0, io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.1.1"), Port: 8081}
					},
				},
			},
			config:         &netcat.Config{},
			expectedOutput: []string{""},
			expectError:    false,
		},
		{
			name:           "Read with partial data",
			expectedOutput: []string{partialData},
			mockConns: []*MockConn{
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, partialData)
						return len(partialData), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.2.1"), Port: 8081}
					},
				},
			},
			config:      &netcat.Config{},
			expectError: false,
		},
		{
			name: "Read with data from multiple connections",
			expectedOutput: []string{
				"part 1part 2part 3",
				"part 1part 3part 2",
				"part 2part 1part 3",
				"part 2part 3part 1",
				"part 3part 1part 2",
				"part 3part 2part 1",
			},
			mockConns: []*MockConn{
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 1")
						return len("part 1"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.3.1"), Port: 8081}
					},
				},
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 2")
						return len("part 2"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.3.2"), Port: 8081}
					},
				},
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 3")
						return len("part 3"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.3.3"), Port: 8081}
					},
				},
			},
			config: &netcat.Config{
				ListenModeOptions: netcat.ListenModeOptions{
					KeepOpen:       true,
					MaxConnections: netcat.DEFAULT_CONNECTION_MAX,
				},
			},
			testLimit:   3,
			expectError: false,
		},
		{
			name: "Read with data from multiple connections in chat mode",
			expectedOutput: []string{
				"user<1>: part 1\nuser<2>: part 2\nuser<3>: part 3\n",
				"user<1>: part 1\nuser<3>: part 3\nuser<2>: part 2\n",
				"user<2>: part 2\nuser<1>: part 1\nuser<3>: part 3\n",
				"user<2>: part 2\nuser<3>: part 3\nuser<1>: part 1\n",
				"user<3>: part 3\nuser<1>: part 1\nuser<2>: part 2\n",
				"user<3>: part 3\nuser<2>: part 2\nuser<1>: part 1\n",
			},
			mockConns: []*MockConn{
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 1\n")
						return len("part 1\n"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.3.1"), Port: 8081}
					},
				},
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 2\n")
						return len("part 2\n"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.3.2"), Port: 8081}
					},
				},
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 3\n")
						return len("part 3\n"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.3.3"), Port: 8081}
					},
				},
			},
			config: &netcat.Config{
				ListenModeOptions: netcat.ListenModeOptions{
					KeepOpen:       true,
					MaxConnections: netcat.DEFAULT_CONNECTION_MAX,
					ChatMode:       true,
					BrokerMode:     true,
				},
			},
			testLimit:   3,
			expectError: false,
		},
		{
			name: "Read with data from multiple connections in broker mode",
			expectedOutput: []string{
				"part 1\npart 2\n",
				"part 2\npart 1\n",
			},
			mockConns: []*MockConn{
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 1\n")
						return len("part 1\n"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.3.1"), Port: 8081}
					},
				},
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 2\n")
						return len("part 2\n"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.3.2"), Port: 8081}
					},
				},
			},
			config: &netcat.Config{
				ListenModeOptions: netcat.ListenModeOptions{
					KeepOpen:       true,
					MaxConnections: netcat.DEFAULT_CONNECTION_MAX,
					ChatMode:       false,
					BrokerMode:     true,
				},
			},
			testLimit:   2,
			expectError: false,
		},
		{
			name: "Disallow a host",
			expectedOutput: []string{
				"part 1part 2",
				"part 2part 1",
			},
			mockConns: []*MockConn{
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 1")
						return len("part 1"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.4.1"), Port: 8081}
					},
				},
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 2")
						return len("part 2"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.4.2"), Port: 8081}
					},
				},
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 3")
						return len("part 3"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.4.3"), Port: 8081}
					},
				},
			},
			config: &netcat.Config{
				ListenModeOptions: netcat.ListenModeOptions{
					KeepOpen:       true,
					MaxConnections: netcat.DEFAULT_CONNECTION_MAX,
				},
				AccessControl: netcat.AccessControlOptions{
					ConnectionList: map[string]bool{
						"127.0.4.3": false,
					},
				},
			},
			testLimit:   3,
			expectError: false,
		},
		{
			name:           "Disallow all localhosts",
			expectedOutput: []string{""},
			mockConns: []*MockConn{
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 1")
						return len("part 1"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8081}
					},
				},
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 2")
						return len("part 2"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.0.2"), Port: 8081}
					},
				},
				{
					ReadFunc: func(b []byte) (int, error) {
						copy(b, "part 3")
						return len("part 3"), io.EOF
					},
					CloseFunc: func() error {
						return nil
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.0.3"), Port: 8081}
					},
				},
			},
			config: &netcat.Config{
				ListenModeOptions: netcat.ListenModeOptions{
					KeepOpen:       true,
					MaxConnections: netcat.DEFAULT_CONNECTION_MAX,
				},
				AccessControl: netcat.AccessControlOptions{
					ConnectionList: map[string]bool{
						"127.0.0.1": false,
						"127.0.0.2": false,
						"127.0.0.3": false,
					},
				},
			},
			testLimit:   3,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentConnIndex := 0

			// MockListener that returns the MockConns in the order they are defined in the test case
			mockListener := &MockListener{
				AcceptFunc: func() (net.Conn, error) {
					if currentConnIndex < len(tt.mockConns) {
						conn := tt.mockConns[currentConnIndex]
						currentConnIndex++
						return conn, nil
					}

					return nil, errors.New("no more MockConns available")
				},
			}

			output := &closableBuffer{}
			cmd := &cmd{
				stdin:  &bytes.Buffer{},
				config: tt.config,
			}

			err := cmd.listenForConnections(netcat.NewConcurrentWriteCloser(output), mockListener, tt.testLimit)
			if (err != nil) != tt.expectError {
				t.Fatalf("Expected error: %v, got: %v", tt.expectError, err != nil)
			}

			// The output may appear in different order as the io.Copy from the connections are executed concurrently
			// So we check if the output is a permutation of the connection reads
			matchFound := slices.Contains(tt.expectedOutput, output.String())

			if !matchFound {
				t.Errorf("Expected output:\n'%v', got:\n'%v'", tt.expectedOutput, output.String())
			}
		})
	}
}

func TestWriteFromListenerToConnection(t *testing.T) {
	tests := []struct {
		name           string
		config         *netcat.Config
		expectedOutput string
		expectError    bool
		mockConns      []*MockConn // Ensure this is defined and includes at least one mock connection
	}{
		{
			name:           "Successful write",
			config:         &netcat.Config{},
			expectedOutput: "abc",
			expectError:    false,
			mockConns: []*MockConn{
				{
					ReadFunc: func(b []byte) (int, error) {
						return 0, io.EOF
					},
					RemoteAddrFunc: func() net.Addr {
						return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8081}
					},
					CloseFunc: func() error {
						return nil
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentConnIndex := 0

			mockListener := &MockListener{
				AcceptFunc: func() (net.Conn, error) {
					if currentConnIndex < len(tt.mockConns) {
						conn := tt.mockConns[currentConnIndex]
						currentConnIndex++
						return conn, nil
					}
					return nil, errors.New("no more MockConns available")
				},
			}

			output := &closableBuffer{}
			cmd := &cmd{
				stdin:  bytes.NewBufferString(tt.expectedOutput),
				config: tt.config,
			}

			err := cmd.listenForConnections(netcat.NewConcurrentWriteCloser(output), mockListener, 0)
			if (err != nil) != tt.expectError {
				t.Fatalf("Expected error: %v, got: %v", tt.expectError, err != nil)
			}

			// Sleep for a while to allow the goroutine to write to the connection
			// This is necessary because the write is done in a goroutine which cannot be waited on
			time.Sleep(1 * time.Second)

			for _, conn := range tt.mockConns {
				if conn.Output.String() != tt.expectedOutput {
					t.Errorf("Expected output 'abc', got '%v'", conn.Output.String())
				}
			}
		})
	}
}

func TestBroadcastMessage(t *testing.T) {
	// Setup
	connections := newConnections(2)

	senderConn, receiverConn := net.Pipe()
	defer senderConn.Close()
	defer receiverConn.Close()

	// Add connections
	connections.add(1, senderConn)
	connections.add(2, receiverConn)

	// Prepare a buffer to capture the broadcast output for the receiver
	var (
		wg             sync.WaitGroup
		outputBuffer   bytes.Buffer
		senderBuffer   bytes.Buffer
		receiverBuffer bytes.Buffer
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(&receiverBuffer, receiverConn)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(&senderBuffer, senderConn)
	}()

	message := "Broadcasted Message!"
	connections.broadcast(netcat.NewConcurrentWriter(&outputBuffer), 2, message)
	senderConn.Close()
	receiverConn.Close()

	wg.Wait()

	if receiverBuffer.String() != message {
		t.Errorf("Expected message to receiver to be %q, got %q", message, receiverBuffer.String())
	}

	if outputBuffer.String() != message {
		t.Errorf("Expected output to be %q, got %q", message, outputBuffer.String())
	}

	if senderBuffer.String() != "" {
		t.Errorf("Expected sender buffer to be empty, got %q", senderBuffer.String())
	}
}

func TestParseRemoteAddr(t *testing.T) {
	tests := []struct {
		name        string
		remoteAddr  string
		socketType  netcat.SocketType
		wantAddress []string
	}{
		{
			socketType:  netcat.SOCKET_TYPE_TCP,
			name:        "IP and Port",
			remoteAddr:  "127.0.0.1:80",
			wantAddress: []string{"127.0.0.1:80", "127.0.0.1", "localhost"},
		},
		{
			socketType:  netcat.SOCKET_TYPE_TCP,
			name:        "IP",
			remoteAddr:  "127.0.0.1",
			wantAddress: []string{"127.0.0.1", "localhost"},
		},
		{
			socketType:  netcat.SOCKET_TYPE_TCP,
			name:        "IPv6 and Port",
			remoteAddr:  "[::1]:80",
			wantAddress: []string{"[::1]:80", "::1"},
		},
		{
			socketType:  netcat.SOCKET_TYPE_TCP,
			name:        "IPv6",
			remoteAddr:  "::1",
			wantAddress: []string{"::1"},
		},
		{
			socketType:  netcat.SOCKET_TYPE_UNIX,
			name:        "Unix Socket",
			remoteAddr:  "/tmp/socket",
			wantAddress: []string{"/tmp/socket"},
		},
		{
			socketType:  netcat.SOCKET_TYPE_NONE,
			name:        "None Socket",
			remoteAddr:  "/tmp/socket",
			wantAddress: []string{"/tmp/socket"},
		},
		{
			socketType:  netcat.SOCKET_TYPE_SCTP,
			name:        "unsupported socket",
			remoteAddr:  "/tmp/socket",
			wantAddress: []string{"/tmp/socket"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddress := parseRemoteAddr(tt.socketType, tt.remoteAddr)
			if !isSubset(t, gotAddress, tt.wantAddress) {
				t.Errorf("parseRemoteAddr(%v, %v) = %v, want a subset of %v", tt.socketType, tt.remoteAddr, gotAddress, tt.wantAddress)
			}
		})
	}
}

func isSubset(t *testing.T, gotAddress, wantAddress []string) bool {
	for _, want := range wantAddress {
		found := slices.Contains(gotAddress, want)
		if !found {
			return false
		}
	}

	return true
}
