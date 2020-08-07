// Copyright 2018 the u-root Authors and Andrea Barberio. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.12

package nclient6

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/hugelgupf/socketpair"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/dhcpv6/server6"
	"github.com/stretchr/testify/require"
)

type handler struct {
	mu       sync.Mutex
	received []*dhcpv6.Message

	// Each received packet can have more than one response (in theory,
	// from different servers sending different Advertise, for example).
	responses [][]*dhcpv6.Message
}

func (h *handler) handle(conn net.PacketConn, peer net.Addr, msg dhcpv6.DHCPv6) {
	h.mu.Lock()
	defer h.mu.Unlock()

	m := msg.(*dhcpv6.Message)

	h.received = append(h.received, m)

	if len(h.responses) > 0 {
		resps := h.responses[0]
		// What should we send in response?
		for _, resp := range resps {
			if _, err := conn.WriteTo(resp.ToBytes(), peer); err != nil {
				panic(err)
			}
		}
		h.responses = h.responses[1:]
	}
}

func serveAndClient(ctx context.Context, responses [][]*dhcpv6.Message, opt ...ClientOpt) (*Client, net.PacketConn) {
	// Fake connection between client and server. No raw sockets, no port
	// weirdness.
	clientRawConn, serverRawConn, err := socketpair.PacketSocketPair()
	if err != nil {
		panic(err)
	}

	o := []ClientOpt{WithRetry(1), WithTimeout(2 * time.Second)}
	o = append(o, opt...)
	mc, err := NewWithConn(clientRawConn, net.HardwareAddr{0xa, 0xb, 0xc, 0xd, 0xe, 0xf}, o...)
	if err != nil {
		panic(err)
	}

	h := &handler{
		responses: responses,
	}
	s, err := server6.NewServer("", nil, h.handle, server6.WithConn(serverRawConn))
	if err != nil {
		panic(err)
	}
	go func() {
		if err := s.Serve(); err != nil {
			panic(err)
		}
	}()

	return mc, serverRawConn
}

func ComparePacket(got *dhcpv6.Message, want *dhcpv6.Message) error {
	if got == nil && got == want {
		return nil
	}
	if (want == nil || got == nil) && (got != want) {
		return fmt.Errorf("packet got %v, want %v", got, want)
	}
	if !bytes.Equal(got.ToBytes(), want.ToBytes()) {
		return fmt.Errorf("packet got %v, want %v", got, want)
	}
	return nil
}

func pktsExpected(got []*dhcpv6.Message, want []*dhcpv6.Message) error {
	if len(got) != len(want) {
		return fmt.Errorf("got %d packets, want %d packets", len(got), len(want))
	}

	for i := range got {
		if err := ComparePacket(got[i], want[i]); err != nil {
			return err
		}
	}
	return nil
}

func newPacket(xid dhcpv6.TransactionID) *dhcpv6.Message {
	p, err := dhcpv6.NewMessage()
	if err != nil {
		panic(fmt.Sprintf("newpacket: %v", err))
	}
	p.TransactionID = xid
	return p
}

func withBufferCap(n int) ClientOpt {
	return func(c *Client) {
		c.bufferCap = n
	}
}

func TestSendAndReadUntil(t *testing.T) {
	for _, tt := range []struct {
		desc   string
		send   *dhcpv6.Message
		server []*dhcpv6.Message

		// If want is nil, we assume server contains what is wanted.
		want    *dhcpv6.Message
		wantErr error
	}{
		{
			desc: "two response packets",
			send: newPacket([3]byte{0x33, 0x33, 0x33}),
			server: []*dhcpv6.Message{
				newPacket([3]byte{0x33, 0x33, 0x33}),
				newPacket([3]byte{0x33, 0x33, 0x33}),
			},
			want: newPacket([3]byte{0x33, 0x33, 0x33}),
		},
		{
			desc: "one response packet",
			send: newPacket([3]byte{0x33, 0x33, 0x33}),
			server: []*dhcpv6.Message{
				newPacket([3]byte{0x33, 0x33, 0x33}),
			},
			want: newPacket([3]byte{0x33, 0x33, 0x33}),
		},
		{
			desc: "one response packet, one invalid XID",
			send: newPacket([3]byte{0x33, 0x33, 0x33}),
			server: []*dhcpv6.Message{
				newPacket([3]byte{0x77, 0x33, 0x33}),
				newPacket([3]byte{0x33, 0x33, 0x33}),
			},
			want: newPacket([3]byte{0x33, 0x33, 0x33}),
		},
		{
			desc: "discard wrong XID",
			send: newPacket([3]byte{0x33, 0x33, 0x33}),
			server: []*dhcpv6.Message{
				newPacket([3]byte{0, 0, 0}),
			},
			want:    nil,
			wantErr: ErrNoResponse,
		},
		{
			desc:    "no response, timeout",
			send:    newPacket([3]byte{0x33, 0x33, 0x33}),
			wantErr: ErrNoResponse,
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			// Both server and client only get 2 seconds.
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			mc, _ := serveAndClient(ctx, [][]*dhcpv6.Message{tt.server},
				// Use an unbuffered channel to make sure we
				// have no deadlocks.
				withBufferCap(0))
			defer mc.Close()

			rcvd, err := mc.SendAndRead(context.Background(), AllDHCPServers, tt.send, nil)
			if err != tt.wantErr {
				t.Error(err)
			}

			if err := ComparePacket(rcvd, tt.want); err != nil {
				t.Errorf("got unexpected packets: %v", err)
			}
		})
	}
}

func TestSimpleSendAndReadDiscardGarbage(t *testing.T) {
	pkt := newPacket([3]byte{0x33, 0x33, 0x33})

	responses := []*dhcpv6.Message{
		newPacket([3]byte{0x33, 0x33, 0x33}),
	}

	// Both the server and client only get 2 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mc, udpConn := serveAndClient(ctx, [][]*dhcpv6.Message{responses})
	defer mc.Close()

	// Too short for valid DHCPv4 packet.
	_, err := udpConn.WriteTo([]byte{0x01}, nil)
	require.NoError(t, err)

	rcvd, err := mc.SendAndRead(context.Background(), AllDHCPServers, pkt, nil)
	if err != nil {
		t.Error(err)
	}

	if err := ComparePacket(rcvd, responses[0]); err != nil {
		t.Errorf("got unexpected packets: %v", err)
	}
}

func TestMultipleSendAndReadOne(t *testing.T) {
	for _, tt := range []struct {
		desc    string
		send    []*dhcpv6.Message
		server  [][]*dhcpv6.Message
		wantErr []error
	}{
		{
			desc: "two requests, two responses",
			send: []*dhcpv6.Message{
				newPacket([3]byte{0x33, 0x33, 0x33}),
				newPacket([3]byte{0x44, 0x44, 0x44}),
			},
			server: [][]*dhcpv6.Message{
				[]*dhcpv6.Message{ // Response for first packet.
					newPacket([3]byte{0x33, 0x33, 0x33}),
				},
				[]*dhcpv6.Message{ // Response for second packet.
					newPacket([3]byte{0x44, 0x44, 0x44}),
				},
			},
			wantErr: []error{
				nil,
				nil,
			},
		},
	} {
		// Both server and client only get 2 seconds.
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		mc, _ := serveAndClient(ctx, tt.server)
		defer mc.conn.Close()

		for i, send := range tt.send {
			rcvd, err := mc.SendAndRead(context.Background(), AllDHCPServers, send, nil)

			if wantErr := tt.wantErr[i]; err != wantErr {
				t.Errorf("SendAndReadOne(%v): got %v, want %v", send, err, wantErr)
			}
			if err := pktsExpected([]*dhcpv6.Message{rcvd}, tt.server[i]); err != nil {
				t.Errorf("got unexpected packets: %v", err)
			}
		}
	}
}
