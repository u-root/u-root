// +build !windows

// Copyright (c) 2019, Google LLC All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tpmutil

import (
	"fmt"
	"net"
	"testing"
	"time"
)

// mockConn records the number of bytes that are read from it and written to it
// and tracks whether or not it has been closed.
type mockConn struct {
	network string
	path    string
	open    bool
}

// dialMockConn returns a mockConn that holds the given network and path info.
func dialMockConn(network, path string) (net.Conn, error) {
	return &mockConn{
		network: network,
		path:    path,
		open:    true,
	}, nil
}

// Read implements a mock version of Read.
func (mc *mockConn) Read(b []byte) (int, error) {
	// Always read zeros into the bytes for the given length.
	for i := range b {
		b[i] = 0
	}
	return len(b), nil
}

// Write implements a mock version of Write.
func (mc *mockConn) Write(b []byte) (int, error) {
	return len(b), nil
}

// Close implements a mock version of Close.
func (mc *mockConn) Close() error {
	if !mc.open {
		return fmt.Errorf("mockConn is already closed")
	}
	mc.open = false
	return nil
}

// LocalAddr returns nil.
func (mc *mockConn) LocalAddr() net.Addr {
	return nil
}

// RemoteAddr returns nil.
func (mc *mockConn) RemoteAddr() net.Addr {
	return nil
}

// SetDeadline returns nil.
func (mc *mockConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline returns nil.
func (mc *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline returns nil.
func (mc *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func newMockEmulator() *EmulatorReadWriteCloser {
	path := "/dev/null/fake"
	rwc := NewEmulatorReadWriteCloser(path)
	rwc.dialer = dialMockConn
	return rwc
}

var (
	input  = []byte(`input`)
	output = make([]byte, 1)
)

func TestEmulatorReadWriteCloserMultipleReads(t *testing.T) {
	rwc := newMockEmulator()
	n, err := rwc.Write(input)
	if err != nil {
		t.Errorf("failed to write: %v", err)
	}
	if n != len(input) {
		t.Errorf("wrong number of bytes written: got %d, expected %d", n, len(input))
	}

	n, err = rwc.Read(output)
	if err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if n != len(output) {
		t.Errorf("wrong number of bytes read: got %d, expected %d", n, len(output))
	}

	n, err = rwc.Write(input)
	if err != nil {
		t.Errorf("failed to write: %v", err)
	}
	if n != len(input) {
		t.Errorf("wrong number of bytes written: got %d, expected %d", n, len(input))
	}

	n, err = rwc.Read(output)
	if err != nil {
		t.Errorf("failed to read: %v", err)
	}
	if n != len(output) {
		t.Errorf("wrong number of bytes read: got %d, expected %d", n, len(output))
	}
}

func TestEmulatorReadWriteCloserClose(t *testing.T) {
	rwc := newMockEmulator()
	if err := rwc.Close(); err == nil {
		t.Errorf("incorrectly closed a connection that hadn't been opened")
	}

	if _, err := rwc.Write(input); err != nil {
		t.Errorf("failed to write: %v", err)
	}

	if err := rwc.Close(); err != nil {
		t.Errorf("failed to close an open connection: %v", err)
	}

	if err := rwc.Close(); err == nil {
		t.Errorf("incorrectly closed a connection that had already been closed")
	}
}

func TestEmulatorReadWriteCloseReadAfterClose(t *testing.T) {
	rwc := newMockEmulator()
	if _, err := rwc.Write(input); err != nil {
		t.Errorf("failed to write: %v", err)
	}

	if err := rwc.Close(); err != nil {
		t.Errorf("failed to close the connection: %v", err)
	}

	if _, err := rwc.Read(output); err == nil {
		t.Errorf("incorrectly read on a closed connection")
	}
}

func TestEmulatorReadWriteCloserReadBeforeWrite(t *testing.T) {
	rwc := newMockEmulator()
	b := make([]byte, 1)
	if _, err := rwc.Read(b); err == nil {
		t.Errorf("incorrectly read on a connection before writing")
	}
}

func TestEmulatorReadWriteCloserDoubleWrite(t *testing.T) {
	rwc := newMockEmulator()
	if _, err := rwc.Write(input); err != nil {
		t.Errorf("failed to write: %v", err)
	}
	if _, err := rwc.Write(input); err == nil {
		t.Errorf("incorrectly wrote a second time without reading in between")
	}
}

func TestEmulatorReadWriteCloserDialerError(t *testing.T) {
	rwc := newMockEmulator()
	rwc.dialer = func(_, _ string) (net.Conn, error) { return nil, fmt.Errorf("invalid") }

	if _, err := rwc.Write(input); err == nil {
		t.Errorf("incorrectly wrote when the dialer returned an error")
	}
}
