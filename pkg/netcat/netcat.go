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

// Network functionality
type UdpRemoteConn struct {
	raddr   *net.UDPAddr
	conn    *net.UDPConn
	wg      *sync.WaitGroup
	once    *sync.Once
	stderr  io.Writer
	verbose bool
}

func (u *UdpRemoteConn) Read(b []byte) (int, error) {
	n, raddr, err := u.conn.ReadFromUDP(b)
	if err != nil {
		return n, err
	}
	setRaddr := func() {
		u.raddr = raddr
		if u.verbose {
			fmt.Fprintln(u.stderr, "Connected to", raddr)
		}
		u.wg.Done()
	}
	u.once.Do(setRaddr)
	return n, nil
}

func (u *UdpRemoteConn) Write(b []byte) (int, error) {
	// we can't answer without raddr, so waiting for incomming request
	u.wg.Wait()
	return u.conn.WriteToUDP(b, u.raddr)
}
