// Copyright 2012-2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netcat

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type UdpRemoteConn struct {
	Raddr   *net.UDPAddr
	Conn    *net.UDPConn
	Wg      *sync.WaitGroup
	Once    *sync.Once
	Stderr  io.Writer
	Verbose bool
}

// Implement interface io.ReadWriter for UdpRemoteConn
func (u *UdpRemoteConn) Read(b []byte) (int, error) {
	n, raddr, err := u.Conn.ReadFromUDP(b)
	if err != nil {
		return n, err
	}
	setRaddr := func() {
		u.Raddr = raddr
		if u.Verbose {
			fmt.Fprintln(u.Stderr, "Connected to", raddr)
		}
		u.Wg.Done()
	}
	u.Once.Do(setRaddr)
	return n, nil
}

func (u *UdpRemoteConn) Write(b []byte) (int, error) {
	// we can't answer without raddr, so waiting for incomming request
	u.Wg.Wait()
	return u.Conn.WriteToUDP(b, u.Raddr)
}

// Netcat connection wrapper
type NetcatConnection struct {
	Conn net.Conn
}

// Implement interface io.ReadWriter for UdpRemoteConn
func (nc *NetcatConnection) Read(b []byte) (int, error) {
	return 0, nil
}

func (nc *NetcatConnection) Write(b []byte) (int, error) {
	return 0, nil
}

func (nc *NetcatConnection) Close() error {
	return nc.Conn.Close()
}
