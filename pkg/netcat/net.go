// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"fmt"
	"net"

	"github.com/u-root/u-root/pkg/ulog"

	// To build the dependencies of this package with TinyGo, we need to include
	// the cpuid package, since tinygo does not support the asm code in the
	// cpuid package. The cpuid package will use the tinygo bridge to get the
	// CPU information. For further information see
	// github.com/u-root/cpuid/cpuid_amd64_tinygo_bridge.go
	_ "github.com/u-root/cpuid"
)

// UDPListener implements net.UDPListener for UDP
type UDPListener struct {
	conn net.Conn
}

// NewUDPListener creates a new UDPListener
func NewUDPListener(network, addr string, _ ulog.Logger) (*UDPListener, error) {
	var conn net.Conn

	switch network {
	case "udp", "udp4", "udp6":
		udpAddr, err := net.ResolveUDPAddr(network, addr)
		if err != nil {
			return nil, err
		}

		conn, err = net.ListenUDP(network, udpAddr)
		if err != nil {
			return nil, err
		}
	case "unixgram":
		unixgramAddr, err := net.ResolveUnixAddr(network, addr)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve Unixgram address: %w", err)
		}

		conn, err = net.ListenUnixgram(network, unixgramAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to listen on Unixgram address: %w", err)
		}
	}
	return &UDPListener{conn: conn}, nil
}

// Accept waits for and returns the next connection to the listener.
func (l *UDPListener) Accept() (net.Conn, error) {
	return l.conn, nil
}

func (l *UDPListener) Close() error {
	return l.conn.Close()
}

func (l *UDPListener) Addr() net.Addr {
	return l.conn.LocalAddr()
}
