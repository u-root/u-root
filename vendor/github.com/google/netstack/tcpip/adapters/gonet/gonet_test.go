// Copyright 2016 The Netstack Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonet

import (
	"net"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/link/loopback"
	"github.com/google/netstack/tcpip/network/ipv4"
	"github.com/google/netstack/tcpip/network/ipv6"
	"github.com/google/netstack/tcpip/stack"
	"github.com/google/netstack/tcpip/transport/tcp"
	"github.com/google/netstack/tcpip/transport/udp"
	"github.com/google/netstack/waiter"
)

const (
	NICID = 1
)

func TestTimeouts(t *testing.T) {
	nc := NewConn(nil, nil)
	dlfs := []struct {
		name string
		f    func(time.Time) error
	}{
		{"SetDeadline", nc.SetDeadline},
		{"SetReadDeadline", nc.SetReadDeadline},
		{"SetWriteDeadline", nc.SetWriteDeadline},
	}

	for _, dlf := range dlfs {
		if err := dlf.f(time.Time{}); err != nil {
			t.Errorf("got %s(time.Time{}) = %v, want = %v", dlf.name, err, nil)
		}
	}
}

func newLoopbackStack() (*stack.Stack, *tcpip.Error) {
	// Create the stack and add a NIC.
	s := stack.New([]string{ipv4.ProtocolName, ipv6.ProtocolName}, []string{tcp.ProtocolName, udp.ProtocolName})

	if err := s.CreateNIC(NICID, loopback.New()); err != nil {
		return nil, err
	}

	// Add default route.
	s.SetRouteTable([]tcpip.Route{
		// IPv4
		{
			Destination: tcpip.Address(strings.Repeat("\x00", 4)),
			Mask:        tcpip.Address(strings.Repeat("\x00", 4)),
			Gateway:     "",
			NIC:         NICID,
		},

		// IPv6
		{
			Destination: tcpip.Address(strings.Repeat("\x00", 16)),
			Mask:        tcpip.Address(strings.Repeat("\x00", 16)),
			Gateway:     "",
			NIC:         NICID,
		},
	})

	return s, nil
}

type testConnection struct {
	wq *waiter.Queue
	e  *waiter.Entry
	ch chan struct{}
	ep tcpip.Endpoint
}

func connect(s *stack.Stack, addr tcpip.FullAddress) (*testConnection, *tcpip.Error) {
	wq := &waiter.Queue{}
	ep, err := s.NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, wq)

	entry, ch := waiter.NewChannelEntry(nil)
	wq.EventRegister(&entry, waiter.EventOut)

	err = ep.Connect(addr)
	if err == tcpip.ErrConnectStarted {
		<-ch
		err = ep.GetSockOpt(tcpip.ErrorOption{})
	}
	if err != nil {
		return nil, err
	}

	wq.EventUnregister(&entry)
	wq.EventRegister(&entry, waiter.EventIn)

	return &testConnection{wq, &entry, ch, ep}, nil
}

func (c *testConnection) close() {
	c.wq.EventUnregister(c.e)
	c.ep.Close()
}

// TestCloseReader tests that Conn.Close() causes Conn.Read() to unblock.
func TestCloseReader(t *testing.T) {
	s, err := newLoopbackStack()
	if err != nil {
		t.Fatalf("newLoopbackStack() = %v", err)
	}

	addr := tcpip.FullAddress{NICID, tcpip.Address(net.IPv4(169, 254, 10, 1).To4()), 11211}

	s.AddAddress(NICID, ipv4.ProtocolNumber, addr.Addr)

	l, e := NewListener(s, addr, ipv4.ProtocolNumber)
	if e != nil {
		t.Fatalf("NewListener() = %v", e)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		c, err := l.Accept()
		if err != nil {
			t.Fatalf("l.Accept() = %v", err)
		}

		// Give c.Read() a chance to block before closing the connection.
		time.AfterFunc(time.Millisecond*50, func() {
			t.Log("c.Close()")
			c.Close()
			t.Log("c.Close() ok")
		})

		buf := make([]byte, 256)
		t.Log("c.Read()")
		n, err := c.Read(buf)
		got, ok := err.(*net.OpError)
		want := tcpip.ErrConnectionAborted
		if n != 0 || !ok || got.Err.Error() != want.String() {
			t.Errorf("c.Read() = (%d, %v), want (0, OpError(%v))", n, err, want)
		}
		t.Logf("c.Read() = %d, %v", n, err)
	}()
	sender, err := connect(s, addr)
	if err != nil {
		t.Fatalf("connect() = %v", err)
	}

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Errorf("c.Read() didn't unblock")
	}
	sender.close()
}

// TestCloseReaderWithForwarder tests that Conn.Close() wakes Conn.Read() when
// using tcp.Forwarder.
func TestCloseReaderWithForwarder(t *testing.T) {
	s, err := newLoopbackStack()
	if err != nil {
		t.Fatalf("newLoopbackStack() = %v", err)
	}

	addr := tcpip.FullAddress{NICID, tcpip.Address(net.IPv4(169, 254, 10, 1).To4()), 11211}
	s.AddAddress(NICID, ipv4.ProtocolNumber, addr.Addr)

	done := make(chan struct{})

	fwd := tcp.NewForwarder(s, 30000, 10, func(r *tcp.ForwarderRequest) {
		defer close(done)

		var wq waiter.Queue
		ep, err := r.CreateEndpoint(&wq)
		if err != nil {
			t.Fatalf("r.CreateEndpoint() = %v", err)
		}
		defer ep.Close()
		r.Complete(false)

		c := NewConn(&wq, ep)

		// Give c.Read() a chance to block before closing the connection.
		time.AfterFunc(time.Millisecond*50, func() {
			t.Log("c.Close()")
			c.Close()
			t.Log("c.Close() ok")
		})

		buf := make([]byte, 256)
		t.Log("c.Read()")
		n, e := c.Read(buf)
		got, ok := e.(*net.OpError)
		want := tcpip.ErrConnectionAborted
		if n != 0 || !ok || got.Err.Error() != want.String() {
			t.Errorf("c.Read() = (%d, %v), want (0, OpError(%v))", n, e, want)
		}
		t.Logf("c.Read() = %d, %v", n, e)
	})
	s.SetTransportProtocolHandler(tcp.ProtocolNumber, fwd.HandlePacket)

	sender, err := connect(s, addr)
	if err != nil {
		t.Fatalf("connect() = %v", err)
	}

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Errorf("c.Read() didn't unblock")
	}
	sender.close()
}

// TestDeadlineChange tests that changing the deadline affects currently blocked reads.
func TestDeadlineChange(t *testing.T) {
	s, err := newLoopbackStack()
	if err != nil {
		t.Fatalf("newLoopbackStack() = %v", err)
	}

	addr := tcpip.FullAddress{NICID, tcpip.Address(net.IPv4(169, 254, 10, 1).To4()), 11211}

	s.AddAddress(NICID, ipv4.ProtocolNumber, addr.Addr)

	l, e := NewListener(s, addr, ipv4.ProtocolNumber)
	if e != nil {
		t.Fatalf("NewListener() = %v", e)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		c, err := l.Accept()
		if err != nil {
			t.Fatalf("l.Accept() = %v", err)
		}

		c.SetDeadline(time.Now().Add(time.Minute))
		// Give c.Read() a chance to block before closing the connection.
		time.AfterFunc(time.Millisecond*50, func() {
			t.Log("c.SetDeadline()")
			c.SetDeadline(time.Now().Add(time.Millisecond * 10))
			t.Log("c.SetDeadline() ok")
		})

		buf := make([]byte, 256)
		t.Log("c.Read()")
		n, err := c.Read(buf)
		got, ok := err.(*net.OpError)
		want := "i/o timeout"
		if n != 0 || !ok || got.Err == nil || got.Err.Error() != want {
			t.Errorf("c.Read() = (%d, %v), want (0, OpError(%s))", n, err, want)
		}
		t.Logf("c.Read() = %d, %v", n, err)
	}()
	sender, err := connect(s, addr)
	if err != nil {
		t.Fatalf("connect() = %v", err)
	}

	select {
	case <-done:
	case <-time.After(time.Millisecond * 500):
		t.Errorf("c.Read() didn't unblock")
	}
	sender.close()
}

func TestPacketConnTransfer(t *testing.T) {
	s, e := newLoopbackStack()
	if e != nil {
		t.Fatalf("newLoopbackStack() = %v", e)
	}

	ip1 := tcpip.Address(net.IPv4(169, 254, 10, 1).To4())
	addr1 := tcpip.FullAddress{NICID, ip1, 11211}
	s.AddAddress(NICID, ipv4.ProtocolNumber, ip1)
	ip2 := tcpip.Address(net.IPv4(169, 254, 10, 2).To4())
	addr2 := tcpip.FullAddress{NICID, ip2, 11311}
	s.AddAddress(NICID, ipv4.ProtocolNumber, ip2)

	c1, err := NewPacketConn(s, addr1, ipv4.ProtocolNumber)
	if err != nil {
		t.Fatal("NewPacketConn(port 4):", err)
	}
	c2, err := NewPacketConn(s, addr2, ipv4.ProtocolNumber)
	if err != nil {
		t.Fatal("NewPacketConn(port 5):", err)
	}

	c1.SetDeadline(time.Now().Add(time.Second))
	c2.SetDeadline(time.Now().Add(time.Second))

	sent := "abc123"
	sendAddr := fullToUDPAddr(addr2)
	if n, err := c1.WriteTo([]byte(sent), sendAddr); err != nil || n != len(sent) {
		t.Errorf("got c1.WriteTo(%q, %v) = %d, %v, want = %d, %v", sent, sendAddr, n, err, len(sent), nil)
	}
	recv := make([]byte, len(sent))
	n, recvAddr, err := c2.ReadFrom(recv)
	if err != nil || n != len(recv) {
		t.Errorf("got c2.ReadFrom() = %d, %v, want = %d, %v", n, err, len(recv), nil)
	}

	if recv := string(recv); recv != sent {
		t.Errorf("got recv = %q, want = %q", recv, sent)
	}

	if want := fullToUDPAddr(addr1); !reflect.DeepEqual(recvAddr, want) {
		t.Errorf("got recvAddr = %v, want = %v", recvAddr, want)
	}

	if err := c1.Close(); err != nil {
		t.Error("c1.Close():", err)
	}
	if err := c2.Close(); err != nil {
		t.Error("c2.Close():", err)
	}
}

func TestTCPConnTransfer(t *testing.T) {
	s, e := newLoopbackStack()
	if e != nil {
		t.Fatalf("newLoopbackStack() = %v", e)
	}

	ip := tcpip.Address(net.IPv4(169, 254, 10, 1).To4())
	addr := tcpip.FullAddress{NICID, ip, 11211}
	s.AddAddress(NICID, ipv4.ProtocolNumber, ip)

	l, err := NewListener(s, addr, ipv4.ProtocolNumber)
	if err != nil {
		t.Fatal("NewListener:", err)
	}
	defer func() {
		if err := l.Close(); err != nil {
			t.Error("l.Close():", err)
		}
	}()

	c1, err := DialTCP(s, addr, ipv4.ProtocolNumber)
	if err != nil {
		t.Fatal("DialTCP:", err)
	}
	defer func() {
		if err := c1.Close(); err != nil {
			t.Error("c1.Close():", err)
		}
	}()

	c2, err := l.Accept()
	if err != nil {
		t.Fatal("l.Accept:", err)
	}
	defer func() {
		if err := c2.Close(); err != nil {
			t.Error("c2.Close():", err)
		}
	}()

	c1.SetDeadline(time.Now().Add(time.Second))
	c2.SetDeadline(time.Now().Add(time.Second))

	const sent = "abc123"

	tests := []struct {
		name string
		c1   net.Conn
		c2   net.Conn
	}{
		{"connected to accepted", c1, c2},
		{"accepted to connected", c2, c1},
	}

	for _, test := range tests {
		if n, err := test.c1.Write([]byte(sent)); err != nil || n != len(sent) {
			t.Errorf("%s: got test.c1.Write(%q) = %d, %v, want = %d, %v", test.name, sent, n, err, len(sent), nil)
			continue
		}

		recv := make([]byte, len(sent))
		n, err := test.c2.Read(recv)
		if err != nil || n != len(recv) {
			t.Errorf("%s: got test.c2.Read() = %d, %v, want = %d, %v", test.name, n, err, len(recv), nil)
			continue
		}

		if recv := string(recv); recv != sent {
			t.Errorf("%s: got recv = %q, want = %q", test.name, recv, sent)
		}
	}
}

func TestTCPDialError(t *testing.T) {
	s, e := newLoopbackStack()
	if e != nil {
		t.Fatalf("newLoopbackStack() = %v", e)
	}

	ip := tcpip.Address(net.IPv4(169, 254, 10, 1).To4())
	addr := tcpip.FullAddress{NICID, ip, 11211}

	_, err := DialTCP(s, addr, ipv4.ProtocolNumber)
	got, ok := err.(*net.OpError)
	want := tcpip.ErrNoRoute
	if !ok || got.Err.Error() != want.String() {
		t.Errorf("Got DialTCP() = %v, want = %v", err, tcpip.ErrNoRoute)
	}
}
