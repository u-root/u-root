// Copyright 2018 the u-root Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dhcp4client is a small, minimum-functionality client for DHCPv4.
//
// It only supports the 4-way DHCPv4 Discover-Offer-Request-Ack handshake as
// well as the Request-Ack renewal process.
package dhcp4client

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/u-root/dhcp4"
	"github.com/u-root/dhcp4/dhcp4opts"
	"github.com/vishvananda/netlink"
)

const (
	maxMessageSize = 1500

	// ClientPort is the port that DHCP clients listen on.
	ClientPort = 68

	// ServerPort is the port that DHCP servers and relay agents listen on.
	ServerPort = 67
)

var (
	// DefaultServers is the address of all link-local DHCP servers and
	// relay agents.
	DefaultServers = &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: ServerPort,
	}
)

// Client is an IPv4 DHCP client.
type Client struct {
	iface   netlink.Link
	conn    net.PacketConn
	timeout time.Duration
	retry   int
}

// New creates a new DHCP client that sends and receives packets on the given
// interface.
func New(iface netlink.Link, opts ...ClientOpt) (*Client, error) {
	c := &Client{
		iface:   iface,
		timeout: 10 * time.Second,
		retry:   3,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.conn == nil {
		var err error
		c.conn, err = NewPacketUDPConn(iface.Attrs().Name, ClientPort)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

// ClientOpt is a function that configures the Client.
type ClientOpt func(*Client) error

// WithTimeout configures the retransmission timeout.
//
// Default is 10 seconds.
//
// TODO(hugelgupf): Check RFC for retransmission behavior.
func WithTimeout(d time.Duration) ClientOpt {
	return func(c *Client) error {
		c.timeout = d
		return nil
	}
}

// WithRetry configures the number of retransmissions to attempt.
//
// Default is 3.
//
// TODO(hugelgupf): Check RFC for retransmission behavior.
func WithRetry(r int) ClientOpt {
	return func(c *Client) error {
		c.retry = r
		return nil
	}
}

// WithConn configures the packet connection to use.
func WithConn(conn net.PacketConn) ClientOpt {
	return func(c *Client) error {
		c.conn = conn
		return nil
	}
}

// DiscoverOffer sends a DHCPDiscover message and returns the first valid offer
// received.
func (c *Client) DiscoverOffer() (*dhcp4.Packet, error) {
	ctx, cancel := context.WithCancel(context.Background())
	wg, out, errCh := c.SimpleSendAndRead(ctx, DefaultServers, c.DiscoverPacket())
	defer func() {
		// Explicitly cancel first, then wait.
		cancel()
		wg.Wait()
	}()

	for packet := range out {
		msgType := dhcp4opts.GetDHCPMessageType(packet.Packet.Options)
		if msgType == dhcp4opts.DHCPOffer {
			// Deferred cancel will cancel the goroutine.
			return packet.Packet, nil
		}
	}

	if err, ok := <-errCh; ok && err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("didn't get a packet")
}

// Request completes the 4-way Discover-Offer-Request-Ack handshake.
func (c *Client) Request() (*dhcp4.Packet, error) {
	offer, err := c.DiscoverOffer()
	if err != nil {
		return nil, err
	}

	return c.SendAndReadOne(c.RequestPacket(offer))
}

// Renew sends a renewal request packet and waits for the corresponding response.
func (c *Client) Renew(ack *dhcp4.Packet) (*dhcp4.Packet, error) {
	return c.SendAndReadOne(c.RequestPacket(ack))
}

// Close closes the client connection.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// SendAndReadOne sends one packet and returns the first response returned by
// any server.
func (c *Client) SendAndReadOne(packet *dhcp4.Packet) (*dhcp4.Packet, error) {
	ctx, cancel := context.WithCancel(context.Background())
	wg, out, errCh := c.SimpleSendAndRead(ctx, DefaultServers, packet)
	defer func() {
		// Explicitly cancel first, then wait.
		cancel()
		wg.Wait()
	}()

	if response, ok := <-out; ok {
		// We're just gonna take the first packet.
		return response.Packet, nil
	}
	if err, ok := <-errCh; ok && err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("no packet received")
}

// DiscoverPacket returns a valid Discover packet for this client.
//
// TODO: Look at RFC and confirm.
func (c *Client) DiscoverPacket() *dhcp4.Packet {
	packet := dhcp4.NewPacket(dhcp4.BootRequest)
	rand.Read(packet.TransactionID[:])
	packet.CHAddr = c.iface.Attrs().HardwareAddr
	packet.Broadcast = true

	packet.Options.Add(dhcp4.OptionDHCPMessageType, dhcp4opts.DHCPDiscover)
	packet.Options.Add(dhcp4.OptionMaximumDHCPMessageSize, dhcp4opts.Uint16(maxMessageSize))
	return packet
}

// RequestPacket returns a valid DHCPRequest packet for the given offer.
//
// TODO: Look at RFC and confirm.
func (c *Client) RequestPacket(offer *dhcp4.Packet) *dhcp4.Packet {
	packet := dhcp4.NewPacket(dhcp4.BootRequest)

	packet.CHAddr = c.iface.Attrs().HardwareAddr
	packet.TransactionID = offer.TransactionID
	packet.CIAddr = offer.CIAddr
	packet.SIAddr = offer.SIAddr
	packet.Broadcast = true

	packet.Options.Add(dhcp4.OptionDHCPMessageType, dhcp4opts.DHCPRequest)
	packet.Options.Add(dhcp4.OptionMaximumDHCPMessageSize, dhcp4opts.Uint16(maxMessageSize))
	// Request the offered IP address.
	packet.Options.Add(dhcp4.OptionRequestedIPAddress, dhcp4opts.IP(offer.YIAddr))

	sid := dhcp4opts.GetServerIdentifier(offer.Options)
	if sid != nil {
		packet.Options.Add(dhcp4.OptionServerIdentifier, dhcp4opts.IP(sid))
	}
	return packet
}

// ClientPacket is a DHCP packet and the interface it corresponds to.
type ClientPacket struct {
	Interface netlink.Link
	Packet    *dhcp4.Packet
}

// ClientError is an error that occured on the associated interface.
type ClientError struct {
	Interface netlink.Link
	Err       error
}

// Error implements error.
func (ce *ClientError) Error() string {
	if ce.Interface != nil {
		return fmt.Sprintf("error on %q: %v", ce.Interface.Attrs().Name, ce.Err)
	}
	return fmt.Sprintf("error without interface: %v", ce.Err)
}

func (c *Client) newClientErr(err error) *ClientError {
	if err == nil {
		return nil
	}
	return &ClientError{
		Interface: c.iface,
		Err:       err,
	}
}

// SimpleSendAndRead broadcasts a DHCP packet and launches a goroutine to read
// response packets. Those response packets will be sent on the channel
// returned.
//
// Callers must cancel ctx when they have received the packet they are looking
// for. Otherwise, the spawned goroutine will keep reading until it times out.
// More importantly, if you send another packet, the spawned goroutine may read
// the response faster than the one launched for the other packet.
//
// Callers sending a packet on one interface should use this. Callers intending
// to send packets on many interface at the same time should look at using
// SendAndRead instead.
//
// Example Usage:
//
//   func sendRequest(someRequest *Packet...) (*Packet, error) {
//     ctx, cancel := context.WithCancel(context.Background())
//     defer cancel()
//
//     out, errCh := c.SimpleSendAndRead(ctx, DefaultServers, someRequest)
//
//     for response := range out {
//       if response == What You Want {
//         // Context cancelation will stop the reading goroutine.
//         return response, ...
//       }
//     }
//
//     if err, ok := <-errCh; ok && err != nil {
//       return nil, err
//     }
//     return nil, fmt.Errorf("got no valid responses")
//   }
//
// TODO(hugelgupf): since the client only has one connection, maybe it should
// just have one dedicated goroutine for reading from the UDP socket, and use a
// request and response queue.
func (c *Client) SimpleSendAndRead(ctx context.Context, dest *net.UDPAddr, p *dhcp4.Packet) (*sync.WaitGroup, <-chan *ClientPacket, <-chan *ClientError) {
	out := make(chan *ClientPacket, 10)
	errOut := make(chan *ClientError, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		c.SendAndRead(ctx, dest, p, out, errOut)
		close(out)
		close(errOut)
		wg.Done()
	}()
	return &wg, out, errOut
}

// SendAndRead sends the given packet `p` to `dest` and reads responses on the
// UDP connection. Any valid DHCP reply packet is sent to `out`.
//
// SendAndRead blocks reading response packets until either:
// - `ctx` is canceled; or
// - we have exhausted all configured retries and timeouts.
//
// SendAndRead retries sending the packet and receiving responses according to
// the configured number of c.retry, using a response timeout of c.timeout.
//
// TODO(hugelgupf): Make this a little state machine of packet types. See RFC
// 2131, Section 4.4, Figure 5.
func (c *Client) SendAndRead(ctx context.Context, dest *net.UDPAddr, p *dhcp4.Packet, out chan<- *ClientPacket, errCh chan<- *ClientError) {
	// This ensures that
	// - we send at most one error on errCh; and
	// - we don't forget to send err on errCh in the many return statements
	//   of sendAndRead.
	if err := c.sendAndRead(ctx, dest, p, out); err != nil {
		errCh <- err
	}
}

func (c *Client) sendAndRead(ctx context.Context, dest *net.UDPAddr, p *dhcp4.Packet, out chan<- *ClientPacket) *ClientError {
	pkt, err := p.MarshalBinary()
	if err != nil {
		return c.newClientErr(err)
	}

	return c.newClientErr(c.retryFn(func() error {
		if _, err := c.conn.WriteTo(pkt, dest); err != nil {
			return fmt.Errorf("error writing packet to connection: %v", err)
		}

		var numPackets int
		timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
		defer cancel()
		for {
			select {
			case <-timeoutCtx.Done():
				if numPackets > 0 {
					return nil
				}

				// No packets received. Sadness.
				return timeoutCtx.Err()
			default:
			}

			// Since a context can be canceled not just because of
			// a deadline, we must check the context every once in
			// a while. Use what is (hopefully) a small part of the
			// context deadline rather than the context's deadline.
			c.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

			// TODO: Clients can send a "max packet size" option in
			// their packets, IIRC. Choose a reasonable size and
			// set it.
			b := make([]byte, 1500)
			n, _, err := c.conn.ReadFrom(b)
			if oerr, ok := err.(net.Error); ok && oerr.Timeout() {
				// Continue to check ctx.Done() above and
				// return the appropriate error.
				continue
			} else if err != nil {
				return fmt.Errorf("error reading from UDP connection: %v", err)
			}

			pkt := &dhcp4.Packet{}
			if err := pkt.UnmarshalBinary(b[:n]); err != nil {
				// Not a valid DHCP reply; keep listening.
				continue
			}

			if pkt.TransactionID != p.TransactionID {
				// Not the right response packet.
				continue
			}

			numPackets++

			clientPkt := &ClientPacket{
				Packet:    pkt,
				Interface: c.iface,
			}

			// Make sure that sending the response has priority.
			select {
			case out <- clientPkt:
				continue
			default:
			}

			// We deliberately only check the parent context here.
			// c.timeout should only apply to reading from the
			// conn, not sending on out.
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- clientPkt:
			}
		}
	}))
}

func (c *Client) retryFn(fn func() error) error {
	// Each retry takes the amount of timeout at worst.
	for i := 0; i < c.retry || c.retry < 0; i++ {
		switch err := fn(); err {
		case nil:
			// Got it!
			return nil

		case context.DeadlineExceeded:
			// Just retry.
			// TODO(hugelgupf): Sleep here for some random amount of time.

		default:
			return err
		}
	}

	return context.DeadlineExceeded
}
