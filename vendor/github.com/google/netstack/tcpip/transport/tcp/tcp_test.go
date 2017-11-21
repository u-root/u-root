// Copyright 2016 The Netstack Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tcp_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/buffer"
	"github.com/google/netstack/tcpip/checker"
	"github.com/google/netstack/tcpip/header"
	"github.com/google/netstack/tcpip/network/ipv4"
	"github.com/google/netstack/tcpip/seqnum"
	"github.com/google/netstack/tcpip/transport/tcp"
	"github.com/google/netstack/tcpip/transport/tcp/testing/context"
	"github.com/google/netstack/waiter"
)

const (
	// defaultMTU is the MTU, in bytes, used throughout the tests, except
	// where another value is explicitly used. It is chosen to match the MTU
	// of loopback interfaces on linux systems.
	defaultMTU = 65535

	// defaultIPv4MSS is the MSS sent by the network stack in SYN/SYN-ACK for an
	// IPv4 endpoint when the MTU is set to defaultMTU in the test.
	defaultIPv4MSS = defaultMTU - header.IPv4MinimumSize - header.TCPMinimumSize
)

func TestGiveUpConnect(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	var wq waiter.Queue
	ep, err := c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, &wq)
	if err != nil {
		t.Fatalf("NeEndpoint failed: %v", err)
	}

	// Register for notification, then start connection attempt.
	waitEntry, notifyCh := waiter.NewChannelEntry(nil)
	wq.EventRegister(&waitEntry, waiter.EventOut)
	defer wq.EventUnregister(&waitEntry)

	err = ep.Connect(tcpip.FullAddress{Addr: context.TestAddr, Port: context.TestPort})
	if err != tcpip.ErrConnectStarted {
		t.Fatalf("Unexpected return value from Connect: %v", err)
	}

	// Close the connection, wait for completion.
	ep.Close()

	// Wait for ep to become writable.
	<-notifyCh
	err = ep.GetSockOpt(tcpip.ErrorOption{})
}

func TestActiveHandshake(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)
}

func TestNonBlockingClose(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)
	ep := c.EP
	c.EP = nil

	// Close the endpoint and measure how long it takes.
	t0 := time.Now()
	ep.Close()
	if diff := time.Now().Sub(t0); diff > 3*time.Second {
		t.Fatalf("Took too long to close: %v", diff)
	}
}

func TestConnectResetAfterClose(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)
	ep := c.EP
	c.EP = nil

	// Close the endpoint, make sure we get a FIN segment, then acknowledge
	// to complete closure of sender, but don't send our own FIN.
	ep.Close()
	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.TCPFlags(header.TCPFlagAck|header.TCPFlagFin),
		),
	)
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  790,
		AckNum:  c.IRS.Add(1),
		RcvWnd:  30000,
	})

	// Wait for the ep to give up waiting for a FIN, and send a RST.
	time.Sleep(3 * time.Second)
	for {
		b := c.GetPacket()
		tcp := header.TCP(header.IPv4(b).Payload())
		if tcp.Flags() == header.TCPFlagAck|header.TCPFlagFin {
			// This is a retransmit of the FIN, ignore it.
			continue
		}

		checker.IPv4(t, b,
			checker.TCP(
				checker.DstPort(context.TestPort),
				checker.SeqNum(uint32(c.IRS)+1),
				checker.AckNum(790),
				checker.TCPFlags(header.TCPFlagAck|header.TCPFlagRst),
			),
		)
		break
	}
}

func TestSimpleReceive(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	we, ch := waiter.NewChannelEntry(nil)
	c.WQ.EventRegister(&we, waiter.EventIn)
	defer c.WQ.EventUnregister(&we)

	if _, err := c.EP.Read(nil); err != tcpip.ErrWouldBlock {
		t.Fatalf("Unexpected error from Read: %v", err)
	}

	data := []byte{1, 2, 3}
	c.SendPacket(data, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  790,
		AckNum:  c.IRS.Add(1),
		RcvWnd:  30000,
	})

	// Wait for receive to be notified.
	select {
	case <-ch:
	case <-time.After(1 * time.Second):
		t.Fatalf("Timed out waiting for data to arrive")
	}

	// Receive data.
	v, err := c.EP.Read(nil)
	if err != nil {
		t.Fatalf("Unexpected error from Read: %v", err)
	}

	if bytes.Compare(data, v) != 0 {
		t.Fatalf("Data is different: expected %v, got %v", data, v)
	}

	// Check that ACK is received.
	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(uint32(790+len(data))),
			checker.TCPFlags(header.TCPFlagAck),
		),
	)
}

func TestOutOfOrderReceive(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	we, ch := waiter.NewChannelEntry(nil)
	c.WQ.EventRegister(&we, waiter.EventIn)
	defer c.WQ.EventUnregister(&we)

	if _, err := c.EP.Read(nil); err != tcpip.ErrWouldBlock {
		t.Fatalf("Unexpected error from Read: %v", err)
	}

	// Send second half of data first, with seqnum 3 ahead of expected.
	data := []byte{1, 2, 3, 4, 5, 6}
	c.SendPacket(data[3:], &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  793,
		AckNum:  c.IRS.Add(1),
		RcvWnd:  30000,
	})

	// Check that we get an ACK specifying which seqnum is expected.
	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.TCPFlags(header.TCPFlagAck),
		),
	)

	// Wait 200ms and check that no data has been received.
	time.Sleep(200 * time.Millisecond)
	if _, err := c.EP.Read(nil); err != tcpip.ErrWouldBlock {
		t.Fatalf("Unexpected error from Read: %v", err)
	}

	// Send the first 3 bytes now.
	c.SendPacket(data[:3], &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  790,
		AckNum:  c.IRS.Add(1),
		RcvWnd:  30000,
	})

	// Receive data.
	read := make([]byte, 0, 6)
	for len(read) < len(data) {
		v, err := c.EP.Read(nil)
		if err != nil {
			if err == tcpip.ErrWouldBlock {
				// Wait for receive to be notified.
				select {
				case <-ch:
				case <-time.After(5 * time.Second):
					t.Fatalf("Timed out waiting for data to arrive")
				}
				continue
			}
			t.Fatalf("Unexpected error from Read: %v", err)
		}

		read = append(read, v...)
	}

	// Check that we received the data in proper order.
	if bytes.Compare(data, read) != 0 {
		t.Fatalf("Data is different: expected %v, got %v", data, read)
	}

	// Check that the whole data is acknowledged.
	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(uint32(790+len(data))),
			checker.TCPFlags(header.TCPFlagAck),
		),
	)
}

func TestOutOfOrderFlood(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	// Create a new connection with initial window size of 10.
	opt := tcpip.ReceiveBufferSizeOption(10)
	c.CreateConnected(789, 30000, &opt)

	if _, err := c.EP.Read(nil); err != tcpip.ErrWouldBlock {
		t.Fatalf("Unexpected error from Read: %v", err)
	}

	// Send 100 packets before the actual one that is expected.
	data := []byte{1, 2, 3, 4, 5, 6}
	for i := 0; i < 100; i++ {
		c.SendPacket(data[3:], &context.Headers{
			SrcPort: context.TestPort,
			DstPort: c.Port,
			Flags:   header.TCPFlagAck,
			SeqNum:  796,
			AckNum:  c.IRS.Add(1),
			RcvWnd:  30000,
		})

		checker.IPv4(t, c.GetPacket(),
			checker.TCP(
				checker.DstPort(context.TestPort),
				checker.SeqNum(uint32(c.IRS)+1),
				checker.AckNum(790),
				checker.TCPFlags(header.TCPFlagAck),
			),
		)
	}

	// Send packet with seqnum 793. It must be discarded because the
	// out-of-order buffer was filled by the previous packets.
	c.SendPacket(data[3:], &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  793,
		AckNum:  c.IRS.Add(1),
		RcvWnd:  30000,
	})

	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.TCPFlags(header.TCPFlagAck),
		),
	)

	// Now send the expected packet, seqnum 790.
	c.SendPacket(data[:3], &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  790,
		AckNum:  c.IRS.Add(1),
		RcvWnd:  30000,
	})

	// Check that only packet 790 is acknowledged.
	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(793),
			checker.TCPFlags(header.TCPFlagAck),
		),
	)
}

func TestFullWindowReceive(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	opt := tcpip.ReceiveBufferSizeOption(10)
	c.CreateConnected(789, 30000, &opt)

	we, ch := waiter.NewChannelEntry(nil)
	c.WQ.EventRegister(&we, waiter.EventIn)
	defer c.WQ.EventUnregister(&we)

	_, err := c.EP.Read(nil)
	if err != tcpip.ErrWouldBlock {
		t.Fatalf("Unexpected error from Read: %v", err)
	}

	// Fill up the window.
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	c.SendPacket(data, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  790,
		AckNum:  c.IRS.Add(1),
		RcvWnd:  30000,
	})

	// Wait for receive to be notified.
	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.Fatalf("Timed out waiting for data to arrive")
	}

	// Check that data is acknowledged, and window goes to zero.
	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(uint32(790+len(data))),
			checker.TCPFlags(header.TCPFlagAck),
			checker.Window(0),
		),
	)

	// Receive data and check it.
	v, err := c.EP.Read(nil)
	if err != nil {
		t.Fatalf("Unexpected error from Read: %v", err)
	}

	if bytes.Compare(data, v) != 0 {
		t.Fatalf("Data is different: expected %v, got %v", data, v)
	}

	// Check that we get an ACK for the newly non-zero window.
	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(uint32(790+len(data))),
			checker.TCPFlags(header.TCPFlagAck),
			checker.Window(10),
		),
	)
}

func TestNoWindowShrinking(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	// Start off with a window size of 10, then shrink it to 5.
	opt := tcpip.ReceiveBufferSizeOption(10)
	c.CreateConnected(789, 30000, &opt)

	opt = 5
	if err := c.EP.SetSockOpt(opt); err != nil {
		t.Fatalf("SetSockOpt failed: %v", err)
	}

	we, ch := waiter.NewChannelEntry(nil)
	c.WQ.EventRegister(&we, waiter.EventIn)
	defer c.WQ.EventUnregister(&we)

	_, err := c.EP.Read(nil)
	if err != tcpip.ErrWouldBlock {
		t.Fatalf("Unexpected error from Read: %v", err)
	}

	// Send 3 bytes, check that the peer acknowledges them.
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	c.SendPacket(data[:3], &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  790,
		AckNum:  c.IRS.Add(1),
		RcvWnd:  30000,
	})

	// Wait for receive to be notified.
	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.Fatalf("Timed out waiting for data to arrive")
	}

	// Check that data is acknowledged, and that window doesn't go to zero
	// just yet because it was previously set to 10. It must go to 7 now.
	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(793),
			checker.TCPFlags(header.TCPFlagAck),
			checker.Window(7),
		),
	)

	// Send 7 more bytes, check that the window fills up.
	c.SendPacket(data[3:], &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  793,
		AckNum:  c.IRS.Add(1),
		RcvWnd:  30000,
	})

	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.Fatalf("Timed out waiting for data to arrive")
	}

	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(uint32(790+len(data))),
			checker.TCPFlags(header.TCPFlagAck),
			checker.Window(0),
		),
	)

	// Receive data and check it.
	read := make([]byte, 0, 10)
	for len(read) < len(data) {
		v, err := c.EP.Read(nil)
		if err != nil {
			t.Fatalf("Unexpected error from Read: %v", err)
		}

		read = append(read, v...)
	}

	if bytes.Compare(data, read) != 0 {
		t.Fatalf("Data is different: expected %v, got %v", data, read)
	}

	// Check that we get an ACK for the newly non-zero window, which is the
	// new size.
	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(uint32(790+len(data))),
			checker.TCPFlags(header.TCPFlagAck),
			checker.Window(5),
		),
	)
}

func TestSimpleSend(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	data := []byte{1, 2, 3}
	view := buffer.NewView(len(data))
	copy(view, data)

	if _, err := c.EP.Write(view, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	// Check that data is received.
	b := c.GetPacket()
	checker.IPv4(t, b,
		checker.PayloadLen(len(data)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)

	if p := b[header.IPv4MinimumSize+header.TCPMinimumSize:]; bytes.Compare(data, p) != 0 {
		t.Fatalf("Data is different: expected %v, got %v", data, p)
	}

	// Acknowledge the data.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  790,
		AckNum:  c.IRS.Add(1 + seqnum.Size(len(data))),
		RcvWnd:  30000,
	})
}

func TestZeroWindowSend(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 0, nil)

	data := []byte{1, 2, 3}
	view := buffer.NewView(len(data))
	copy(view, data)

	_, err := c.EP.Write(view, nil)
	if err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	// Since the window is currently zero, check that no packet is received.
	c.CheckNoPacket("Packet received when window is zero")

	// Open up the window. Data should be received now.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  790,
		AckNum:  c.IRS.Add(1),
		RcvWnd:  30000,
	})

	// Check that data is received.
	b := c.GetPacket()
	checker.IPv4(t, b,
		checker.PayloadLen(len(data)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)

	if p := b[header.IPv4MinimumSize+header.TCPMinimumSize:]; bytes.Compare(data, p) != 0 {
		t.Fatalf("Data is different: expected %v, got %v", data, p)
	}

	// Acknowledge the data.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  790,
		AckNum:  c.IRS.Add(1 + seqnum.Size(len(data))),
		RcvWnd:  30000,
	})
}

func TestScaledWindowConnect(t *testing.T) {
	// This test ensures that window scaling is used when the peer
	// does advertise it and connection is established with Connect().
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	// Set the window size greater than the maximum non-scaled window.
	opt := tcpip.ReceiveBufferSizeOption(65535 * 3)
	c.CreateConnectedWithRawOptions(789, 30000, &opt, []byte{
		header.TCPOptionWS, 3, 0, header.TCPOptionNOP,
	})

	data := []byte{1, 2, 3}
	view := buffer.NewView(len(data))
	copy(view, data)

	if _, err := c.EP.Write(view, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	// Check that data is received, and that advertised window is 0xbfff,
	// that is, that it is scaled.
	b := c.GetPacket()
	checker.IPv4(t, b,
		checker.PayloadLen(len(data)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.Window(0xbfff),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)
}

func TestNonScaledWindowConnect(t *testing.T) {
	// This test ensures that window scaling is not used when the peer
	// doesn't advertise it and connection is established with Connect().
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	// Set the window size greater than the maximum non-scaled window.
	opt := tcpip.ReceiveBufferSizeOption(65535 * 3)
	c.CreateConnected(789, 30000, &opt)

	data := []byte{1, 2, 3}
	view := buffer.NewView(len(data))
	copy(view, data)

	if _, err := c.EP.Write(view, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	// Check that data is received, and that advertised window is 0xffff,
	// that is, that it's not scaled.
	b := c.GetPacket()
	checker.IPv4(t, b,
		checker.PayloadLen(len(data)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.Window(0xffff),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)
}

func TestScaledWindowAccept(t *testing.T) {
	// This test ensures that window scaling is used when the peer
	// does advertise it and connection is established with Accept().
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	// Create EP and start listening.
	wq := &waiter.Queue{}
	ep, err := c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, wq)
	if err != nil {
		t.Fatalf("NewEndpoint failed: %v", err)
	}
	defer ep.Close()

	// Set the window size greater than the maximum non-scaled window.
	if err := ep.SetSockOpt(tcpip.ReceiveBufferSizeOption(65535 * 3)); err != nil {
		t.Fatalf("SetSockOpt failed failed: %v", err)
	}

	if err := ep.Bind(tcpip.FullAddress{Port: context.StackPort}, nil); err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if err := ep.Listen(10); err != nil {
		t.Fatalf("Listen failed: %v", err)
	}

	// Do 3-way handshake.
	c.PassiveConnectWithOptions(100, 2, header.TCPSynOptions{MSS: defaultIPv4MSS})

	// Try to accept the connection.
	we, ch := waiter.NewChannelEntry(nil)
	wq.EventRegister(&we, waiter.EventIn)
	defer wq.EventUnregister(&we)

	c.EP, _, err = ep.Accept()
	if err == tcpip.ErrWouldBlock {
		// Wait for connection to be established.
		select {
		case <-ch:
			c.EP, _, err = ep.Accept()
			if err != nil {
				t.Fatalf("Accept failed: %v", err)
			}

		case <-time.After(1 * time.Second):
			t.Fatalf("Timed out waiting for accept")
		}
	}

	data := []byte{1, 2, 3}
	view := buffer.NewView(len(data))
	copy(view, data)

	if _, err := c.EP.Write(view, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	// Check that data is received, and that advertised window is 0xbfff,
	// that is, that it is scaled.
	b := c.GetPacket()
	checker.IPv4(t, b,
		checker.PayloadLen(len(data)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.Window(0xbfff),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)
}

func TestNonScaledWindowAccept(t *testing.T) {
	// This test ensures that window scaling is not used when the peer
	// doesn't advertise it and connection is established with Accept().
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	// Create EP and start listening.
	wq := &waiter.Queue{}
	ep, err := c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, wq)
	if err != nil {
		t.Fatalf("NewEndpoint failed: %v", err)
	}
	defer ep.Close()

	// Set the window size greater than the maximum non-scaled window.
	if err := ep.SetSockOpt(tcpip.ReceiveBufferSizeOption(65535 * 3)); err != nil {
		t.Fatalf("SetSockOpt failed failed: %v", err)
	}

	if err := ep.Bind(tcpip.FullAddress{Port: context.StackPort}, nil); err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if err := ep.Listen(10); err != nil {
		t.Fatalf("Listen failed: %v", err)
	}

	// Do 3-way handshake.
	c.PassiveConnect(100, 2, header.TCPSynOptions{MSS: defaultIPv4MSS})

	// Try to accept the connection.
	we, ch := waiter.NewChannelEntry(nil)
	wq.EventRegister(&we, waiter.EventIn)
	defer wq.EventUnregister(&we)

	c.EP, _, err = ep.Accept()
	if err == tcpip.ErrWouldBlock {
		// Wait for connection to be established.
		select {
		case <-ch:
			c.EP, _, err = ep.Accept()
			if err != nil {
				t.Fatalf("Accept failed: %v", err)
			}

		case <-time.After(1 * time.Second):
			t.Fatalf("Timed out waiting for accept")
		}
	}

	data := []byte{1, 2, 3}
	view := buffer.NewView(len(data))
	copy(view, data)

	if _, err := c.EP.Write(view, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	// Check that data is received, and that advertised window is 0xffff,
	// that is, that it's not scaled.
	b := c.GetPacket()
	checker.IPv4(t, b,
		checker.PayloadLen(len(data)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.Window(0xffff),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)
}

func TestZeroScaledWindowReceive(t *testing.T) {
	// This test ensures that the endpoint sends a non-zero window size
	// advertisement when the scaled window transitions from 0 to non-zero,
	// but the actual window (not scaled) hasn't gotten to zero.
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	// Set the window size such that a window scale of 4 will be used.
	const wnd = 65535 * 10
	const ws = uint32(4)
	opt := tcpip.ReceiveBufferSizeOption(wnd)
	c.CreateConnectedWithRawOptions(789, 30000, &opt, []byte{
		header.TCPOptionWS, 3, 0, header.TCPOptionNOP,
	})

	// Write chunks of 50000 bytes.
	remain := wnd
	sent := 0
	data := make([]byte, 50000)
	for remain > len(data) {
		c.SendPacket(data, &context.Headers{
			SrcPort: context.TestPort,
			DstPort: c.Port,
			Flags:   header.TCPFlagAck,
			SeqNum:  seqnum.Value(790 + sent),
			AckNum:  c.IRS.Add(1),
			RcvWnd:  30000,
		})
		sent += len(data)
		remain -= len(data)
		checker.IPv4(t, c.GetPacket(),
			checker.PayloadLen(header.TCPMinimumSize),
			checker.TCP(
				checker.DstPort(context.TestPort),
				checker.SeqNum(uint32(c.IRS)+1),
				checker.AckNum(uint32(790+sent)),
				checker.Window(uint16(remain>>ws)),
				checker.TCPFlags(header.TCPFlagAck),
			),
		)
	}

	// Make the window non-zero, but the scaled window zero.
	if remain >= 16 {
		data = data[:remain-15]
		c.SendPacket(data, &context.Headers{
			SrcPort: context.TestPort,
			DstPort: c.Port,
			Flags:   header.TCPFlagAck,
			SeqNum:  seqnum.Value(790 + sent),
			AckNum:  c.IRS.Add(1),
			RcvWnd:  30000,
		})
		sent += len(data)
		remain -= len(data)
		checker.IPv4(t, c.GetPacket(),
			checker.PayloadLen(header.TCPMinimumSize),
			checker.TCP(
				checker.DstPort(context.TestPort),
				checker.SeqNum(uint32(c.IRS)+1),
				checker.AckNum(uint32(790+sent)),
				checker.Window(0),
				checker.TCPFlags(header.TCPFlagAck),
			),
		)
	}

	// Read some data. An ack should be sent in response to that.
	v, err := c.EP.Read(nil)
	if err != nil {
		t.Fatalf("Unexpected error from Read: %v", err)
	}

	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(uint32(790+sent)),
			checker.Window(uint16(len(v)>>ws)),
			checker.TCPFlags(header.TCPFlagAck),
		),
	)
}

func testBrokenUpWrite(t *testing.T, c *context.Context, maxPayload int) {
	payloadMultiplier := 10
	dataLen := payloadMultiplier * maxPayload
	data := make([]byte, dataLen)
	for i := range data {
		data[i] = byte(i)
	}

	view := buffer.NewView(len(data))
	copy(view, data)

	if _, err := c.EP.Write(view, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	// Check that data is received in chunks.
	bytesReceived := 0
	numPackets := 0
	for bytesReceived != dataLen {
		b := c.GetPacket()
		numPackets++
		tcp := header.TCP(header.IPv4(b).Payload())
		payloadLen := len(tcp.Payload())
		checker.IPv4(t, b,
			checker.TCP(
				checker.DstPort(context.TestPort),
				checker.SeqNum(uint32(c.IRS)+1+uint32(bytesReceived)),
				checker.AckNum(790),
				checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
			),
		)

		pdata := data[bytesReceived : bytesReceived+payloadLen]
		if p := tcp.Payload(); bytes.Compare(pdata, p) != 0 {
			t.Fatalf("Data is different: expected %v, got %v", pdata, p)
		}
		bytesReceived += payloadLen
		var options []byte
		if c.TimeStampEnabled {
			// If timestamp option is enabled, echo back the timestamp and increment
			// the TSEcr value included in the packet and send that back as the TSVal.
			parsedOpts := tcp.ParsedOptions()
			tsOpt := header.EncodeTSOption(parsedOpts.TSEcr+1, parsedOpts.TSVal)
			options = append(options, tsOpt[:]...)
		}
		// Acknowledge the data.
		c.SendPacket(nil, &context.Headers{
			SrcPort: context.TestPort,
			DstPort: c.Port,
			Flags:   header.TCPFlagAck,
			SeqNum:  790,
			AckNum:  c.IRS.Add(1 + seqnum.Size(bytesReceived)),
			RcvWnd:  30000,
			TCPOpts: options,
		})
	}
	if numPackets == 1 {
		t.Fatalf("expected write to be broken up into multiple packets, but got 1 packet")
	}
}

func TestSendGreaterThanMTU(t *testing.T) {
	const maxPayload = 100
	c := context.New(t, uint32(header.TCPMinimumSize+header.IPv4MinimumSize+maxPayload))
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)
	testBrokenUpWrite(t, c, maxPayload)
}

func TestActiveSendMSSLessThanMTU(t *testing.T) {
	const maxPayload = 100
	c := context.New(t, 65535)
	defer c.Cleanup()

	c.CreateConnectedWithRawOptions(789, 30000, nil, []byte{
		header.TCPOptionMSS, 4, byte(maxPayload / 256), byte(maxPayload % 256),
	})
	testBrokenUpWrite(t, c, maxPayload)
}

func TestPassiveSendMSSLessThanMTU(t *testing.T) {
	const maxPayload = 100
	const mtu = 1200
	c := context.New(t, mtu)
	defer c.Cleanup()

	// Create EP and start listening.
	wq := &waiter.Queue{}
	ep, err := c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, wq)
	if err != nil {
		t.Fatalf("NewEndpoint failed: %v", err)
	}
	defer ep.Close()

	// Set the buffer size to a deterministic size so that we can check the
	// window scaling option.
	const rcvBufferSize = 0x20000
	const wndScale = 2
	if err := ep.SetSockOpt(tcpip.ReceiveBufferSizeOption(rcvBufferSize)); err != nil {
		t.Fatalf("SetSockOpt failed failed: %v", err)
	}

	if err := ep.Bind(tcpip.FullAddress{Port: context.StackPort}, nil); err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if err := ep.Listen(10); err != nil {
		t.Fatalf("Listen failed: %v", err)
	}

	// Do 3-way handshake.
	c.PassiveConnect(maxPayload, wndScale, header.TCPSynOptions{MSS: mtu - header.IPv4MinimumSize - header.TCPMinimumSize})

	// Try to accept the connection.
	we, ch := waiter.NewChannelEntry(nil)
	wq.EventRegister(&we, waiter.EventIn)
	defer wq.EventUnregister(&we)

	c.EP, _, err = ep.Accept()
	if err == tcpip.ErrWouldBlock {
		// Wait for connection to be established.
		select {
		case <-ch:
			c.EP, _, err = ep.Accept()
			if err != nil {
				t.Fatalf("Accept failed: %v", err)
			}

		case <-time.After(1 * time.Second):
			t.Fatalf("Timed out waiting for accept")
		}
	}

	// Check that data gets properly segmented.
	testBrokenUpWrite(t, c, maxPayload)
}

func TestSynCookiePassiveSendMSSLessThanMTU(t *testing.T) {
	const maxPayload = 536
	const mtu = 2000
	c := context.New(t, mtu)
	defer c.Cleanup()

	// Set the SynRcvd threshold to zero to force a syn cookie based accept
	// to happen.
	saved := tcp.SynRcvdCountThreshold
	defer func() {
		tcp.SynRcvdCountThreshold = saved
	}()
	tcp.SynRcvdCountThreshold = 0

	// Create EP and start listening.
	wq := &waiter.Queue{}
	ep, err := c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, wq)
	if err != nil {
		t.Fatalf("NewEndpoint failed: %v", err)
	}
	defer ep.Close()

	if err := ep.Bind(tcpip.FullAddress{Port: context.StackPort}, nil); err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if err := ep.Listen(10); err != nil {
		t.Fatalf("Listen failed: %v", err)
	}

	// Do 3-way handshake.
	c.PassiveConnect(maxPayload, -1, header.TCPSynOptions{MSS: mtu - header.IPv4MinimumSize - header.TCPMinimumSize})

	// Try to accept the connection.
	we, ch := waiter.NewChannelEntry(nil)
	wq.EventRegister(&we, waiter.EventIn)
	defer wq.EventUnregister(&we)

	c.EP, _, err = ep.Accept()
	if err == tcpip.ErrWouldBlock {
		// Wait for connection to be established.
		select {
		case <-ch:
			c.EP, _, err = ep.Accept()
			if err != nil {
				t.Fatalf("Accept failed: %v", err)
			}

		case <-time.After(1 * time.Second):
			t.Fatalf("Timed out waiting for accept")
		}
	}

	// Check that data gets properly segmented.
	testBrokenUpWrite(t, c, maxPayload)
}

func TestForwarderSendMSSLessThanMTU(t *testing.T) {
	const maxPayload = 100
	const mtu = 1200
	c := context.New(t, mtu)
	defer c.Cleanup()

	s := c.Stack()
	ch := make(chan *tcpip.Error, 1)
	f := tcp.NewForwarder(s, 65536, 10, func(r *tcp.ForwarderRequest) {
		var err *tcpip.Error
		c.EP, err = r.CreateEndpoint(&c.WQ)
		ch <- err
	})
	s.SetTransportProtocolHandler(tcp.ProtocolNumber, f.HandlePacket)

	// Do 3-way handshake.
	c.PassiveConnect(maxPayload, 1, header.TCPSynOptions{MSS: mtu - header.IPv4MinimumSize - header.TCPMinimumSize})

	// Wait for connection to be available.
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("Error creating endpoint: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("Timed out waiting for connection")
	}

	// Check that data gets properly segmented.
	testBrokenUpWrite(t, c, maxPayload)
}

func TestSynOptionsOnActiveConnect(t *testing.T) {
	const mtu = 1400
	c := context.New(t, mtu)
	defer c.Cleanup()

	// Create TCP endpoint.
	var err *tcpip.Error
	c.EP, err = c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, &c.WQ)
	if err != nil {
		t.Fatalf("NewEndpoint failed: %v", err)
	}

	// Set the buffer size to a deterministic size so that we can check the
	// window scaling option.
	const rcvBufferSize = 0x20000
	const wndScale = 2
	if err := c.EP.SetSockOpt(tcpip.ReceiveBufferSizeOption(rcvBufferSize)); err != nil {
		t.Fatalf("SetSockOpt failed failed: %v", err)
	}

	// Start connection attempt.
	we, ch := waiter.NewChannelEntry(nil)
	c.WQ.EventRegister(&we, waiter.EventOut)
	defer c.WQ.EventUnregister(&we)

	err = c.EP.Connect(tcpip.FullAddress{Addr: context.TestAddr, Port: context.TestPort})
	if err != tcpip.ErrConnectStarted {
		t.Fatalf("Unexpected return value from Connect: %v", err)
	}

	// Receive SYN packet.
	b := c.GetPacket()

	checker.IPv4(t, b,
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.TCPFlags(header.TCPFlagSyn),
			checker.TCPSynOptions(header.TCPSynOptions{MSS: mtu - header.IPv4MinimumSize - header.TCPMinimumSize, WS: wndScale}),
		),
	)

	tcp := header.TCP(header.IPv4(b).Payload())
	c.IRS = seqnum.Value(tcp.SequenceNumber())

	// Wait for retransmit.
	time.Sleep(1 * time.Second)
	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.TCPFlags(header.TCPFlagSyn),
			checker.SrcPort(tcp.SourcePort()),
			checker.SeqNum(tcp.SequenceNumber()),
			checker.TCPSynOptions(header.TCPSynOptions{MSS: mtu - header.IPv4MinimumSize - header.TCPMinimumSize, WS: wndScale}),
		),
	)

	// Send SYN-ACK.
	iss := seqnum.Value(789)
	c.SendPacket(nil, &context.Headers{
		SrcPort: tcp.DestinationPort(),
		DstPort: tcp.SourcePort(),
		Flags:   header.TCPFlagSyn | header.TCPFlagAck,
		SeqNum:  iss,
		AckNum:  c.IRS.Add(1),
		RcvWnd:  30000,
	})

	// Receive ACK packet.
	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.TCPFlags(header.TCPFlagAck),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(uint32(iss)+1),
		),
	)

	// Wait for connection to be established.
	select {
	case <-ch:
		err = c.EP.GetSockOpt(tcpip.ErrorOption{})
		if err != nil {
			t.Fatalf("Unexpected error when connecting: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("Timed out waiting for connection")
	}
}

func TestCloseListener(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	// Create listener.
	var wq waiter.Queue
	ep, err := c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, &wq)
	if err != nil {
		t.Fatalf("NewEndpoint failed: %v", err)
	}

	if err := ep.Bind(tcpip.FullAddress{}, nil); err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if err := ep.Listen(10); err != nil {
		t.Fatalf("Listen failed: %v", err)
	}

	// Close the listener and measure how long it takes.
	t0 := time.Now()
	ep.Close()
	if diff := time.Now().Sub(t0); diff > 3*time.Second {
		t.Fatalf("Took too long to close: %v", diff)
	}
}

func TestReceiveOnResetConnection(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	// Send RST segment.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagRst,
		SeqNum:  790,
		RcvWnd:  30000,
	})

	// Try to read.
	we, ch := waiter.NewChannelEntry(nil)
	c.WQ.EventRegister(&we, waiter.EventIn)
	defer c.WQ.EventUnregister(&we)

loop:
	for {
		switch _, err := c.EP.Read(nil); err {
		case nil:
			t.Fatalf("Unexpected success.")
		case tcpip.ErrWouldBlock:
			select {
			case <-ch:
			case <-time.After(1 * time.Second):
				t.Fatalf("Timed out waiting for reset to arrive")
			}
		case tcpip.ErrConnectionReset:
			break loop
		default:
			t.Fatalf("Unexpected error: want %v, got %v", tcpip.ErrConnectionReset, err)
		}
	}
}

func TestSendOnResetConnection(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	// Send RST segment.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagRst,
		SeqNum:  790,
		RcvWnd:  30000,
	})

	// Wait for the RST to be received.
	time.Sleep(1 * time.Second)

	// Try to write.
	view := buffer.NewView(10)
	_, err := c.EP.Write(view, nil)
	if err != tcpip.ErrConnectionReset {
		t.Fatalf("Unexpected error from Write: want %v, got %v", tcpip.ErrConnectionReset, err)
	}
}

func TestFinImmediately(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	// Shutdown immediately, check that we get a FIN.
	if err := c.EP.Shutdown(tcpip.ShutdownWrite); err != nil {
		t.Fatalf("Unexpected error from Shutdown: %v", err)
	}

	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.TCPFlags(header.TCPFlagAck|header.TCPFlagFin),
		),
	)

	// Ack and send FIN as well.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck | header.TCPFlagFin,
		SeqNum:  790,
		AckNum:  c.IRS.Add(2),
		RcvWnd:  30000,
	})

	// Check that the stack acks the FIN.
	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+2),
			checker.AckNum(791),
			checker.TCPFlags(header.TCPFlagAck),
		),
	)
}

func TestFinRetransmit(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	// Shutdown immediately, check that we get a FIN.
	if err := c.EP.Shutdown(tcpip.ShutdownWrite); err != nil {
		t.Fatalf("Unexpected error from Shutdown: %v", err)
	}

	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.TCPFlags(header.TCPFlagAck|header.TCPFlagFin),
		),
	)

	// Don't acknowledge yet. We should get a retransmit of the FIN.
	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.TCPFlags(header.TCPFlagAck|header.TCPFlagFin),
		),
	)

	// Ack and send FIN as well.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck | header.TCPFlagFin,
		SeqNum:  790,
		AckNum:  c.IRS.Add(2),
		RcvWnd:  30000,
	})

	// Check that the stack acks the FIN.
	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+2),
			checker.AckNum(791),
			checker.TCPFlags(header.TCPFlagAck),
		),
	)
}

func TestFinWithNoPendingData(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	// Write something out, and have it acknowledged.
	view := buffer.NewView(10)
	if _, err := c.EP.Write(view, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	next := uint32(c.IRS) + 1
	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(len(view)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(790),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)
	next += uint32(len(view))

	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  790,
		AckNum:  seqnum.Value(next),
		RcvWnd:  30000,
	})

	// Shutdown, check that we get a FIN.
	if err := c.EP.Shutdown(tcpip.ShutdownWrite); err != nil {
		t.Fatalf("Unexpected error from Shutdown: %v", err)
	}

	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(790),
			checker.TCPFlags(header.TCPFlagAck|header.TCPFlagFin),
		),
	)
	next++

	// Ack and send FIN as well.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck | header.TCPFlagFin,
		SeqNum:  790,
		AckNum:  seqnum.Value(next),
		RcvWnd:  30000,
	})

	// Check that the stack acks the FIN.
	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(791),
			checker.TCPFlags(header.TCPFlagAck),
		),
	)
}

func DisabledTestFinWithPendingDataCwndFull(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	// Write something out but don't ACK it yet.
	view := buffer.NewView(10)
	if _, err := c.EP.Write(view, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	next := uint32(c.IRS) + 1
	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(len(view)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(790),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)
	next += uint32(len(view))

	// Shutdown the connection, check that the FIN segment isn't sent
	// because the congestion window doesn't allow it. Wait until a
	// retransmit is received.
	if err := c.EP.Shutdown(tcpip.ShutdownWrite); err != nil {
		t.Fatalf("Unexpected error from Shutdown: %v", err)
	}

	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(len(view)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next-uint32(len(view))),
			checker.AckNum(790),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)

	// Send the ACK that will allow the FIN to be sent as well.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  790,
		AckNum:  seqnum.Value(next),
		RcvWnd:  30000,
	})

	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(790),
			checker.TCPFlags(header.TCPFlagAck|header.TCPFlagFin),
		),
	)
	next++

	// Send a FIN that acknowledges everything. Get an ACK back.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck | header.TCPFlagFin,
		SeqNum:  790,
		AckNum:  seqnum.Value(next),
		RcvWnd:  30000,
	})

	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(791),
			checker.TCPFlags(header.TCPFlagAck),
		),
	)
}

func TestFinWithPendingData(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	// Write something out, and acknowledge it to get cwnd to 2.
	view := buffer.NewView(10)
	if _, err := c.EP.Write(view, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	next := uint32(c.IRS) + 1
	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(len(view)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(790),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)
	next += uint32(len(view))

	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  790,
		AckNum:  seqnum.Value(next),
		RcvWnd:  30000,
	})

	// Write new data, but don't acknowledge it.
	if _, err := c.EP.Write(view, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(len(view)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(790),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)
	next += uint32(len(view))

	// Shutdown the connection, check that we do get a FIN.
	if err := c.EP.Shutdown(tcpip.ShutdownWrite); err != nil {
		t.Fatalf("Unexpected error from Shutdown: %v", err)
	}

	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(790),
			checker.TCPFlags(header.TCPFlagAck|header.TCPFlagFin),
		),
	)
	next++

	// Send a FIN that acknowledges everything. Get an ACK back.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck | header.TCPFlagFin,
		SeqNum:  790,
		AckNum:  seqnum.Value(next),
		RcvWnd:  30000,
	})

	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(791),
			checker.TCPFlags(header.TCPFlagAck),
		),
	)
}

func TestFinWithPartialAck(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	// Write something out, and acknowledge it to get cwnd to 2. Also send
	// FIN from the test side.
	view := buffer.NewView(10)
	if _, err := c.EP.Write(view, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	next := uint32(c.IRS) + 1
	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(len(view)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(790),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)
	next += uint32(len(view))

	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck | header.TCPFlagFin,
		SeqNum:  790,
		AckNum:  seqnum.Value(next),
		RcvWnd:  30000,
	})

	// Check that we get an ACK for the fin.
	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(791),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)

	// Write new data, but don't acknowledge it.
	if _, err := c.EP.Write(view, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(len(view)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(791),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)
	next += uint32(len(view))

	// Shutdown the connection, check that we do get a FIN.
	if err := c.EP.Shutdown(tcpip.ShutdownWrite); err != nil {
		t.Fatalf("Unexpected error from Shutdown: %v", err)
	}

	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(next),
			checker.AckNum(791),
			checker.TCPFlags(header.TCPFlagAck|header.TCPFlagFin),
		),
	)
	next++

	// Send an ACK for the data, but not for the FIN yet.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  791,
		AckNum:  seqnum.Value(next - 1),
		RcvWnd:  30000,
	})

	// Check that we don't get a retransmit of the FIN.
	c.CheckNoPacketTimeout("FIN retransmitted when data was ack'd", 100*time.Millisecond)

	// Ack the FIN.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck | header.TCPFlagFin,
		SeqNum:  791,
		AckNum:  seqnum.Value(next),
		RcvWnd:  30000,
	})
}

func DisabledTestExponentialIncreaseDuringSlowStart(t *testing.T) {
	maxPayload := 10
	c := context.New(t, uint32(header.TCPMinimumSize+header.IPv4MinimumSize+maxPayload))
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	const iterations = 7
	data := buffer.NewView(maxPayload * (1 << (iterations + 1)))
	for i := range data {
		data[i] = byte(i)
	}

	// Write all the data in one shot. Packets will only be written at the
	// MTU size though.
	if _, err := c.EP.Write(data, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	expected := 1
	bytesRead := 0
	for i := 0; i < iterations; i++ {
		// Read all packets expected on this iteration. Don't
		// acknowledge any of them just yet, so that we can measure the
		// congestion window.
		for j := 0; j < expected; j++ {
			c.ReceiveAndCheckPacket(data, bytesRead, maxPayload)
			bytesRead += maxPayload
		}

		// Check we don't receive any more packets on this iteration.
		// The timeout can't be too high or we'll trigger a timeout.
		c.CheckNoPacketTimeout("More packets received than expected for this cwnd.", 50*time.Millisecond)

		// Acknowledge all the data received so far.
		c.SendAck(790, bytesRead)

		// Double the number of expected packets for the next iteration.
		expected *= 2
	}
}

func DisabledTestCongestionAvoidance(t *testing.T) {
	maxPayload := 10
	c := context.New(t, uint32(header.TCPMinimumSize+header.IPv4MinimumSize+maxPayload))
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	const iterations = 7
	data := buffer.NewView(2 * maxPayload * (1 << (iterations + 1)))
	for i := range data {
		data[i] = byte(i)
	}

	// Write all the data in one shot. Packets will only be written at the
	// MTU size though.
	if _, err := c.EP.Write(data, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	// Do slow start for a few iterations.
	expected := 1
	bytesRead := 0
	for i := 0; i < iterations; i++ {
		expected = 1 << uint(i)
		if i > 0 {
			// Acknowledge all the data received so far if not on
			// first iteration.
			c.SendAck(790, bytesRead)
		}

		// Read all packets expected on this iteration. Don't
		// acknowledge any of them just yet, so that we can measure the
		// congestion window.
		for j := 0; j < expected; j++ {
			c.ReceiveAndCheckPacket(data, bytesRead, maxPayload)
			bytesRead += maxPayload
		}

		// Check we don't receive any more packets on this iteration.
		// The timeout can't be too high or we'll trigger a timeout.
		c.CheckNoPacketTimeout("More packets received than expected for this cwnd.", 50*time.Millisecond)
	}

	// Don't acknowledge the first packet of the last packet train. Let's
	// wait for them to time out, which will trigger a restart of slow
	// start, and initialization of ssthresh to cwnd/2.
	rtxOffset := bytesRead - maxPayload*expected
	c.ReceiveAndCheckPacket(data, rtxOffset, maxPayload)

	// Acknowledge all the data received so far.
	c.SendAck(790, bytesRead)

	// This part is tricky: when the timeout happened, we had "expected"
	// packets pending, cwnd reset to 1, and ssthresh set to expected/2.
	// By acknowledging "expected" packets, the slow-start part will
	// increase cwnd to expected/2 (which "consumes" expected/2-1 of the
	// acknowledgements), then the congestion avoidance part will consume
	// an extra expected/2 acks to take cwnd to expected/2 + 1. One ack
	// remains in the "ack count" (which will cause cwnd to be incremented
	// once it reaches cwnd acks).
	//
	// So we're straight into congestion avoidance with cwnd set to
	// expected/2 + 1.
	//
	// Check that packets trains of cwnd packets are sent, and that cwnd is
	// incremented by 1 after we acknowledge each packet.
	expected = expected/2 + 1
	for i := 0; i < iterations; i++ {
		// Read all packets expected on this iteration. Don't
		// acknowledge any of them just yet, so that we can measure the
		// congestion window.
		for j := 0; j < expected; j++ {
			c.ReceiveAndCheckPacket(data, bytesRead, maxPayload)
			bytesRead += maxPayload
		}

		// Check we don't receive any more packets on this iteration.
		// The timeout can't be too high or we'll trigger a timeout.
		c.CheckNoPacketTimeout("More packets received than expected for this cwnd.", 50*time.Millisecond)

		// Acknowledge all the data received so far.
		c.SendAck(790, bytesRead)

		// In cogestion avoidance, the packets trains increase by 1 in
		// each iteration.
		expected++
	}
}

func DisabledTestFastRecovery(t *testing.T) {
	maxPayload := 10
	c := context.New(t, uint32(header.TCPMinimumSize+header.IPv4MinimumSize+maxPayload))
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	const iterations = 7
	data := buffer.NewView(2 * maxPayload * (1 << (iterations + 1)))
	for i := range data {
		data[i] = byte(i)
	}

	// Write all the data in one shot. Packets will only be written at the
	// MTU size though.
	if _, err := c.EP.Write(data, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	// Do slow start for a few iterations.
	expected := 1
	bytesRead := 0
	for i := 0; i < iterations; i++ {
		expected = 1 << uint(i)
		if i > 0 {
			// Acknowledge all the data received so far if not on
			// first iteration.
			c.SendAck(790, bytesRead)
		}

		// Read all packets expected on this iteration. Don't
		// acknowledge any of them just yet, so that we can measure the
		// congestion window.
		for j := 0; j < expected; j++ {
			c.ReceiveAndCheckPacket(data, bytesRead, maxPayload)
			bytesRead += maxPayload
		}

		// Check we don't receive any more packets on this iteration.
		// The timeout can't be too high or we'll trigger a timeout.
		c.CheckNoPacketTimeout("More packets received than expected for this cwnd.", 50*time.Millisecond)
	}

	// Send 10 duplicate acks. This should force an immediate retransmit of
	// the pending packet, and inflation of cwnd to expected/2+7.
	rtxOffset := bytesRead - maxPayload*expected
	for i := 0; i < 10; i++ {
		c.SendAck(790, rtxOffset)
	}

	// Receive the retransmitted packet.
	c.ReceiveAndCheckPacket(data, rtxOffset, maxPayload)

	// Acknowledge half of the pending data.
	rtxOffset = bytesRead - expected*maxPayload/2
	c.SendAck(790, rtxOffset)

	// Receive the retransmit due to partial ack.
	c.ReceiveAndCheckPacket(data, rtxOffset, maxPayload)

	// This part is tricky: when the retransmit happened, we had "expected"
	// packets pending, cwnd reset to expected/2, and ssthresh set to
	// expected/2. By acknowledging expected/2 packets, 7 new packets are
	// allowed to be sent immediately.
	for j := 0; j < 7; j++ {
		c.ReceiveAndCheckPacket(data, bytesRead, maxPayload)
		bytesRead += maxPayload
	}

	c.CheckNoPacketTimeout("More packets received than expected for this cwnd.", 50*time.Millisecond)

	// Acknowledge all pending data.
	c.SendAck(790, bytesRead)

	// Now the inflation is removed, so cwnd is expected/2. But since we've
	// received expected+7 packets since cwnd changed, it must now be set
	// expected/2 + 2, given that floor((expected+7)/(expected/2)) == 2.
	expected = expected/2 + 2
	for i := 0; i < iterations; i++ {
		// Read all packets expected on this iteration. Don't
		// acknowledge any of them just yet, so that we can measure the
		// congestion window.
		for j := 0; j < expected; j++ {
			c.ReceiveAndCheckPacket(data, bytesRead, maxPayload)
			bytesRead += maxPayload
		}

		// Check we don't receive any more packets on this iteration.
		// The timeout can't be too high or we'll trigger a timeout.
		c.CheckNoPacketTimeout("More packets received than expected for this cwnd.", 50*time.Millisecond)

		// Acknowledge all the data received so far.
		c.SendAck(790, bytesRead)

		// In cogestion avoidance, the packets trains increase by 1 in
		// each iteration.
		expected++
	}
}

func DisabledTestRetransmit(t *testing.T) {
	maxPayload := 10
	c := context.New(t, uint32(header.TCPMinimumSize+header.IPv4MinimumSize+maxPayload))
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	const iterations = 7
	data := buffer.NewView(maxPayload * (1 << (iterations + 1)))
	for i := range data {
		data[i] = byte(i)
	}

	// Write all the data in two shots. Packets will only be written at the
	// MTU size though.
	if _, err := c.EP.Write(data[:len(data)/2], nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}
	if _, err := c.EP.Write(data[len(data)/2:], nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	// Do slow start for a few iterations.
	expected := 1
	bytesRead := 0
	for i := 0; i < iterations; i++ {
		expected = 1 << uint(i)
		if i > 0 {
			// Acknowledge all the data received so far if not on
			// first iteration.
			c.SendAck(790, bytesRead)
		}

		// Read all packets expected on this iteration. Don't
		// acknowledge any of them just yet, so that we can measure the
		// congestion window.
		for j := 0; j < expected; j++ {
			c.ReceiveAndCheckPacket(data, bytesRead, maxPayload)
			bytesRead += maxPayload
		}

		// Check we don't receive any more packets on this iteration.
		// The timeout can't be too high or we'll trigger a timeout.
		c.CheckNoPacketTimeout("More packets received than expected for this cwnd.", 50*time.Millisecond)
	}

	// Wait for a timeout and retransmit.
	rtxOffset := bytesRead - maxPayload*expected
	c.ReceiveAndCheckPacket(data, rtxOffset, maxPayload)

	// Acknowledge half of the pending data.
	rtxOffset = bytesRead - expected*maxPayload/2
	c.SendAck(790, rtxOffset)

	// Receive the remaining data, making sure that acknowledge data is not
	// retransmitted.
	for offset := rtxOffset; offset < len(data); offset += maxPayload {
		c.ReceiveAndCheckPacket(data, offset, maxPayload)
		c.SendAck(790, offset+maxPayload)
	}

	c.CheckNoPacketTimeout("More packets received than expected for this cwnd.", 50*time.Millisecond)
}

func TestUpdateListenBacklog(t *testing.T) {
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	// Create listener.
	var wq waiter.Queue
	ep, err := c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, &wq)
	if err != nil {
		t.Fatalf("NewEndpoint failed: %v", err)
	}

	if err := ep.Bind(tcpip.FullAddress{}, nil); err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if err := ep.Listen(10); err != nil {
		t.Fatalf("Listen failed: %v", err)
	}

	// Update the backlog with another Listen() on the same endpoint.
	if err := ep.Listen(20); err != nil {
		t.Fatalf("Listen failed to update backlog: %v", err)
	}

	ep.Close()
}

func scaledSendWindow(t *testing.T, scale uint8) {
	// This test ensures that the endpoint is using the right scaling by
	// sending a buffer that is larger than the window size, and ensuring
	// that the endpoint doesn't send more than allowed.
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	maxPayload := defaultMTU - header.IPv4MinimumSize - header.TCPMinimumSize
	c.CreateConnectedWithRawOptions(789, 0, nil, []byte{
		header.TCPOptionMSS, 4, byte(maxPayload / 256), byte(maxPayload % 256),
		header.TCPOptionWS, 3, scale, header.TCPOptionNOP,
	})

	// Open up the window with a scaled value.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck,
		SeqNum:  790,
		AckNum:  c.IRS.Add(1),
		RcvWnd:  1,
	})

	// Send some data. Check that it's capped by the window size.
	view := buffer.NewView(65535)
	if _, err := c.EP.Write(view, nil); err != nil {
		t.Fatalf("Unexpected error from Write: %v", err)
	}

	// Check that only data that fits in the scaled window is sent.
	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen((1<<scale)+header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.TCPFlagsMatch(header.TCPFlagAck, ^uint8(header.TCPFlagPsh)),
		),
	)

	// Reset the connection to free resources.
	c.SendPacket(nil, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagRst,
		SeqNum:  790,
	})
}

func TestScaledSendWindow(t *testing.T) {
	for scale := uint8(0); scale <= 14; scale++ {
		scaledSendWindow(t, scale)
	}
}

func TestReceivedSegmentQueuing(t *testing.T) {
	// This test sends 200 segments containing a few bytes each to an
	// endpoint and checks that they're all received and acknowledged by
	// the endpoint, that is, that none of the segments are dropped by
	// internal queues.
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	// Send 200 segments.
	data := []byte{1, 2, 3}
	for i := 0; i < 200; i++ {
		c.SendPacket(data, &context.Headers{
			SrcPort: context.TestPort,
			DstPort: c.Port,
			Flags:   header.TCPFlagAck,
			SeqNum:  seqnum.Value(790 + i*len(data)),
			AckNum:  c.IRS.Add(1),
			RcvWnd:  30000,
		})
	}

	// Receive ACKs for all segments.
	last := seqnum.Value(790 + 200*len(data))
	for {
		b := c.GetPacket()
		checker.IPv4(t, b,
			checker.TCP(
				checker.DstPort(context.TestPort),
				checker.SeqNum(uint32(c.IRS)+1),
				checker.TCPFlags(header.TCPFlagAck),
			),
		)
		tcp := header.TCP(header.IPv4(b).Payload())
		ack := seqnum.Value(tcp.AckNumber())
		if ack == last {
			break
		}

		if last.LessThan(ack) {
			t.Fatalf("Acknowledge (%v) beyond the expected (%v)", ack, last)
		}
	}
}

func TestReadAfterClosedState(t *testing.T) {
	// This test ensures that calling Read() or Peek() after the endpoint
	// has transitioned to closedState still works if there is pending
	// data. To transition to stateClosed without calling Close(), we must
	// shutdown the send path and the peer must send its own FIN.
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	c.CreateConnected(789, 30000, nil)

	we, ch := waiter.NewChannelEntry(nil)
	c.WQ.EventRegister(&we, waiter.EventIn)
	defer c.WQ.EventUnregister(&we)

	if _, err := c.EP.Read(nil); err != tcpip.ErrWouldBlock {
		t.Fatalf("Unexpected error from Read: %v", err)
	}

	// Shutdown immediately for write, check that we get a FIN.
	if err := c.EP.Shutdown(tcpip.ShutdownWrite); err != nil {
		t.Fatalf("Unexpected error from Shutdown: %v", err)
	}

	checker.IPv4(t, c.GetPacket(),
		checker.PayloadLen(header.TCPMinimumSize),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+1),
			checker.AckNum(790),
			checker.TCPFlags(header.TCPFlagAck|header.TCPFlagFin),
		),
	)

	// Send some data and acknowledge the FIN.
	data := []byte{1, 2, 3}
	c.SendPacket(data, &context.Headers{
		SrcPort: context.TestPort,
		DstPort: c.Port,
		Flags:   header.TCPFlagAck | header.TCPFlagFin,
		SeqNum:  790,
		AckNum:  c.IRS.Add(2),
		RcvWnd:  30000,
	})

	// Check that ACK is received.
	checker.IPv4(t, c.GetPacket(),
		checker.TCP(
			checker.DstPort(context.TestPort),
			checker.SeqNum(uint32(c.IRS)+2),
			checker.AckNum(uint32(791+len(data))),
			checker.TCPFlags(header.TCPFlagAck),
		),
	)

	// Give the stack the chance to transition to closed state.
	time.Sleep(1 * time.Second)

	// Wait for receive to be notified.
	select {
	case <-ch:
	case <-time.After(1 * time.Second):
		t.Fatalf("Timed out waiting for data to arrive")
	}

	// Check that peek works.
	peekBuf := make([]byte, 10)
	n, err := c.EP.Peek([][]byte{peekBuf})
	if err != nil {
		t.Fatalf("Unexpected error from Peek: %v", err)
	}

	peekBuf = peekBuf[:n]
	if bytes.Compare(data, peekBuf) != 0 {
		t.Fatalf("Data is different: expected %v, got %v", data, peekBuf)
	}

	// Receive data.
	v, err := c.EP.Read(nil)
	if err != nil {
		t.Fatalf("Unexpected error from Read: %v", err)
	}

	if bytes.Compare(data, v) != 0 {
		t.Fatalf("Data is different: expected %v, got %v", data, v)
	}

	// Now that we drained the queue, check that functions fail with the
	// right error code.
	if _, err := c.EP.Read(nil); err != tcpip.ErrClosedForReceive {
		t.Fatalf("Unexpected return from Read: got %v, want %v", err, tcpip.ErrClosedForReceive)
	}

	if _, err := c.EP.Peek([][]byte{peekBuf}); err != tcpip.ErrClosedForReceive {
		t.Fatalf("Unexpected return from Peek: got %v, want %v", err, tcpip.ErrClosedForReceive)
	}
}

func TestReusePort(t *testing.T) {
	// This test ensures that ports are immediately available for reuse
	// after Close on the endpoints using them returns.
	c := context.New(t, defaultMTU)
	defer c.Cleanup()

	// First case, just an endpoint that was bound.
	var err *tcpip.Error
	c.EP, err = c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, &waiter.Queue{})
	if err != nil {
		t.Fatalf("NewEndpoint failed; %v", err)
	}
	if err := c.EP.Bind(tcpip.FullAddress{Port: context.StackPort}, nil); err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	c.EP.Close()
	c.EP, err = c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, &waiter.Queue{})
	if err != nil {
		t.Fatalf("NewEndpoint failed; %v", err)
	}
	if err := c.EP.Bind(tcpip.FullAddress{Port: context.StackPort}, nil); err != nil {
		t.Fatalf("Bind failed: %v", err)
	}
	c.EP.Close()

	// Second case, an endpoint that was bound and is connecting..
	c.EP, err = c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, &waiter.Queue{})
	if err != nil {
		t.Fatalf("NewEndpoint failed; %v", err)
	}
	if err := c.EP.Bind(tcpip.FullAddress{Port: context.StackPort}, nil); err != nil {
		t.Fatalf("Bind failed: %v", err)
	}
	err = c.EP.Connect(tcpip.FullAddress{Addr: context.TestAddr, Port: context.TestPort})
	if err != tcpip.ErrConnectStarted {
		t.Fatalf("Unexpected return value from Connect: %v", err)
	}
	c.EP.Close()

	c.EP, err = c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, &waiter.Queue{})
	if err != nil {
		t.Fatalf("NewEndpoint failed; %v", err)
	}
	if err := c.EP.Bind(tcpip.FullAddress{Port: context.StackPort}, nil); err != nil {
		t.Fatalf("Bind failed: %v", err)
	}
	c.EP.Close()

	// Third case, an endpoint that was bound and is listening.
	c.EP, err = c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, &waiter.Queue{})
	if err != nil {
		t.Fatalf("NewEndpoint failed; %v", err)
	}
	if err := c.EP.Bind(tcpip.FullAddress{Port: context.StackPort}, nil); err != nil {
		t.Fatalf("Bind failed: %v", err)
	}
	err = c.EP.Listen(10)
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	c.EP.Close()

	c.EP, err = c.Stack().NewEndpoint(tcp.ProtocolNumber, ipv4.ProtocolNumber, &waiter.Queue{})
	if err != nil {
		t.Fatalf("NewEndpoint failed; %v", err)
	}
	if err := c.EP.Bind(tcpip.FullAddress{Port: context.StackPort}, nil); err != nil {
		t.Fatalf("Bind failed: %v", err)
	}
	err = c.EP.Listen(10)
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
}
