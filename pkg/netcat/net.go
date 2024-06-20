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

type UDPRemoteConn struct {
	Raddr         *net.UDPAddr
	AccessControl NetcatAccessControlOptions
	Conn          *net.UDPConn
	Wg            *sync.WaitGroup
	Once          *sync.Once
	Stderr        io.Writer
	Verbose       bool
}

// NewUDPRemoteConn creates a new UDPRemoteConn object
// 1. Resolve UDP address from network and address
// 2. Get `UDPCon` from `ListenUDP`
// 3. Create a `sync.WaitGroup` with delta `wgDelta`
// 4. Return a new `UDPRemoteConn` object
func NewUDPRemoteConn(network string, address string, stderr io.Writer, accessControl NetcatAccessControlOptions, verbose bool) (*UDPRemoteConn, error) {
	udpAddr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP(network, udpAddr)
	if err != nil {
		return nil, err
	}

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	return &UDPRemoteConn{
		Conn:          conn,
		AccessControl: accessControl,
		Wg:            waitGroup,
		Once:          &sync.Once{},
		Stderr:        stderr,
		Verbose:       verbose,
	}, nil
}

// Implement interface io.ReadWriter for UDPRemoteConn
func (u *UDPRemoteConn) Read(b []byte) (int, error) {
	n, raddr, err := u.Conn.ReadFromUDP(b)
	if err != nil {
		return n, err
	}

	// return if host is not allowed
	if !u.AccessControl.IsAllowed(raddr.IP.String()) {
		return 0, nil
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

// Implement interface io.ReadWriter for UDPRemoteConn
func (u *UDPRemoteConn) Write(b []byte) (int, error) {
	// we can't answer without raddr, so waiting for incomming request
	u.Wg.Wait()
	return u.Conn.WriteToUDP(b, u.Raddr)
}

type UnixRemoteConn struct {
	Raddr         *net.UnixAddr
	AccessControl NetcatAccessControlOptions
	Conn          *net.UnixConn
	Wg            *sync.WaitGroup
	Once          *sync.Once
	Stderr        io.Writer
	Verbose       bool
}

// NewUnixgramRemoteConn creates a new UnixRemoteConn object
// 1. Resolve Unix address from network and address
// 2. Get `UnixCon` from `ListenUnixgram`
// 3. Create a `sync.WaitGroup` with delta `wgDelta`
// 4. Return a new `UnixRemoteConn` object
func NewUnixgramRemoteConn(network string, address string, stderr io.Writer, accessControl NetcatAccessControlOptions, verbose bool) (*UnixRemoteConn, error) {
	addr, err := net.ResolveUnixAddr("unixgram", address)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve Unix address: %v", err)
	}

	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on Unix address: %v", err)
	}

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	return &UnixRemoteConn{
		Conn:          conn,
		AccessControl: accessControl,
		Wg:            waitGroup,
		Once:          &sync.Once{},
		Stderr:        stderr,
		Verbose:       verbose,
	}, nil
}

// Implement interface io.ReadWriter for UDPRemoteConn
func (u *UnixRemoteConn) Read(b []byte) (int, error) {
	n, raddr, err := u.Conn.ReadFromUnix(b)
	if err != nil {
		return n, err
	}

	// return if host is not allowed
	if !u.AccessControl.IsAllowed(raddr.Name) {
		return 0, nil
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

// Implement interface io.ReadWriter for UnixRemoteConn
func (u *UnixRemoteConn) Write(b []byte) (int, error) {
	u.Wg.Wait()
	return u.Conn.Write(b)
}

// Netcat connection wrapper
type NetcatConnection struct {
	Conn net.Conn
}

func (nc *NetcatConnection) Read(b []byte) (int, error) {
	return 0, nil
}

func (nc *NetcatConnection) Write(b []byte) (int, error) {
	return 0, nil
}

func (nc *NetcatConnection) Close() error {
	return nc.Conn.Close()
}
