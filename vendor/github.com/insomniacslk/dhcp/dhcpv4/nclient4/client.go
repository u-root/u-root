// Copyright 2018 the u-root Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.12

// Package nclient4 is a small, minimum-functionality client for DHCPv4.
//
// It only supports the 4-way DHCPv4 Discover-Offer-Request-Ack handshake as
// well as the Request-Ack renewal process.
package nclient4

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

const (
	defaultTimeout   = 5 * time.Second
	defaultRetries   = 3
	defaultBufferCap = 5
	maxMessageSize   = 1500

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

var (
	// ErrNoResponse is returned when no response packet is received.
	ErrNoResponse = errors.New("no matching response packet received")
)

// pendingCh is a channel associated with a pending TransactionID.
type pendingCh struct {
	// SendAndRead closes done to indicate that it wishes for no more
	// messages for this particular XID.
	done <-chan struct{}

	// ch is used by the receive loop to distribute DHCP messages.
	ch chan<- *dhcpv4.DHCPv4
}

type logger interface {
	Printf(format string, v ...interface{})
	PrintMessage(prefix string, message *dhcpv4.DHCPv4)
}

type emptyLogger struct{}

func (e emptyLogger) Printf(format string, v ...interface{})             {}
func (e emptyLogger) PrintMessage(prefix string, message *dhcpv4.DHCPv4) {}

type shortSummaryLogger struct {
	*log.Logger
}

func (s shortSummaryLogger) Printf(format string, v ...interface{}) {
	s.Logger.Printf(format, v...)
}
func (s shortSummaryLogger) PrintMessage(prefix string, message *dhcpv4.DHCPv4) {
	s.Printf("%s: %s", prefix, message)
}

type debugLogger struct {
	*log.Logger
}

func (d debugLogger) Printf(format string, v ...interface{}) {
	d.Logger.Printf(format, v...)
}
func (d debugLogger) PrintMessage(prefix string, message *dhcpv4.DHCPv4) {
	d.Printf("%s: %s", prefix, message.Summary())
}

// Client is an IPv4 DHCP client.
type Client struct {
	ifaceHWAddr net.HardwareAddr
	conn        net.PacketConn
	timeout     time.Duration
	retry       int
	logger      logger

	// bufferCap is the channel capacity for each TransactionID.
	bufferCap int

	// serverAddr is the UDP address to send all packets to.
	//
	// This may be an actual broadcast address, or a unicast address.
	serverAddr *net.UDPAddr

	// closed is an atomic bool set to 1 when done is closed.
	closed uint32

	// done is closed to unblock the receive loop.
	done chan struct{}

	// wg protects any spawned goroutines, namely the receiveLoop.
	wg sync.WaitGroup

	pendingMu sync.Mutex
	// pending stores the distribution channels for each pending
	// TransactionID. receiveLoop uses this map to determine which channel
	// to send a new DHCP message to.
	pending map[dhcpv4.TransactionID]*pendingCh
}

// New returns a client usable with an unconfigured interface.
func New(iface string, opts ...ClientOpt) (*Client, error) {
	i, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}
	pc, err := NewRawUDPConn(iface, ClientPort)
	if err != nil {
		return nil, err
	}
	return NewWithConn(pc, i.HardwareAddr, opts...)
}

// NewWithConn creates a new DHCP client that sends and receives packets on the
// given interface.
func NewWithConn(conn net.PacketConn, ifaceHWAddr net.HardwareAddr, opts ...ClientOpt) (*Client, error) {
	c := &Client{
		ifaceHWAddr: ifaceHWAddr,
		timeout:     defaultTimeout,
		retry:       defaultRetries,
		serverAddr:  DefaultServers,
		bufferCap:   defaultBufferCap,
		conn:        conn,
		logger:      emptyLogger{},

		done:    make(chan struct{}),
		pending: make(map[dhcpv4.TransactionID]*pendingCh),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.conn == nil {
		return nil, fmt.Errorf("no connection given")
	}
	c.wg.Add(1)
	go c.receiveLoop()
	return c, nil
}

// Close closes the underlying connection.
func (c *Client) Close() error {
	// Make sure not to close done twice.
	if !atomic.CompareAndSwapUint32(&c.closed, 0, 1) {
		return nil
	}

	err := c.conn.Close()

	// Closing c.done sets off a chain reaction:
	//
	// Any SendAndRead unblocks trying to receive more messages, which
	// means rem() gets called.
	//
	// rem() should be unblocking receiveLoop if it is blocked.
	//
	// receiveLoop should then exit gracefully.
	close(c.done)

	// Wait for receiveLoop to stop.
	c.wg.Wait()

	return err
}

func isErrClosing(err error) bool {
	// Unfortunately, the epoll-connection-closed error is internal to the
	// net library.
	return strings.Contains(err.Error(), "use of closed network connection")
}

func (c *Client) receiveLoop() {
	defer c.wg.Done()
	for {
		// TODO: Clients can send a "max packet size" option in their
		// packets, IIRC. Choose a reasonable size and set it.
		b := make([]byte, maxMessageSize)
		n, _, err := c.conn.ReadFrom(b)
		if err != nil {
			if !isErrClosing(err) {
				c.logger.Printf("error reading from UDP connection: %v", err)
			}
			return
		}

		msg, err := dhcpv4.FromBytes(b[:n])
		if err != nil {
			// Not a valid DHCP packet; keep listening.
			continue
		}

		if msg.OpCode != dhcpv4.OpcodeBootReply {
			// Not a response message.
			continue
		}

		// This is a somewhat non-standard check, by the looks
		// of RFC 2131. It should work as long as the DHCP
		// server is spec-compliant for the HWAddr field.
		if c.ifaceHWAddr != nil && !bytes.Equal(c.ifaceHWAddr, msg.ClientHWAddr) {
			// Not for us.
			continue
		}

		c.pendingMu.Lock()
		p, ok := c.pending[msg.TransactionID]
		if ok {
			select {
			case <-p.done:
				close(p.ch)
				delete(c.pending, msg.TransactionID)

			// This send may block.
			case p.ch <- msg:
			}
		}
		c.pendingMu.Unlock()
	}
}

// ClientOpt is a function that configures the Client.
type ClientOpt func(*Client)

// WithTimeout configures the retransmission timeout.
//
// Default is 5 seconds.
func WithTimeout(d time.Duration) ClientOpt {
	return func(c *Client) {
		c.timeout = d
	}
}

// WithSummaryLogger logs one-line DHCPv4 message summarys when sent & received.
func WithSummaryLogger() ClientOpt {
	return func(c *Client) {
		c.logger = shortSummaryLogger{
			Logger: log.New(os.Stderr, "[dhcpv4]", log.LstdFlags),
		}
	}
}

// WithDebugLogger logs multi-line full DHCPv4 messages when sent & received.
func WithDebugLogger() ClientOpt {
	return func(c *Client) {
		c.logger = debugLogger{
			Logger: log.New(os.Stderr, "[dhcpv4]", log.LstdFlags),
		}
	}
}

func withBufferCap(n int) ClientOpt {
	return func(c *Client) {
		c.bufferCap = n
	}
}

// WithRetry configures the number of retransmissions to attempt.
//
// Default is 3.
func WithRetry(r int) ClientOpt {
	return func(c *Client) {
		c.retry = r
	}
}

// WithServerAddr configures the address to send messages to.
func WithServerAddr(n *net.UDPAddr) ClientOpt {
	return func(c *Client) {
		c.serverAddr = n
	}
}

// Matcher matches DHCP packets.
type Matcher func(*dhcpv4.DHCPv4) bool

// IsMessageType returns a matcher that checks for the message type.
//
// If t is MessageTypeNone, all packets are matched.
func IsMessageType(t dhcpv4.MessageType) Matcher {
	return func(p *dhcpv4.DHCPv4) bool {
		return p.MessageType() == t || t == dhcpv4.MessageTypeNone
	}
}

// DiscoverOffer sends a DHCPDiscover message and returns the first valid offer
// received.
func (c *Client) DiscoverOffer(ctx context.Context, modifiers ...dhcpv4.Modifier) (*dhcpv4.DHCPv4, error) {
	// RFC 2131, Section 4.4.1, Table 5 details what a DISCOVER packet should
	// contain.
	discover, err := dhcpv4.NewDiscovery(c.ifaceHWAddr, dhcpv4.PrependModifiers(modifiers,
		dhcpv4.WithOption(dhcpv4.OptMaxMessageSize(maxMessageSize)))...)
	if err != nil {
		return nil, err
	}
	return c.SendAndRead(ctx, c.serverAddr, discover, IsMessageType(dhcpv4.MessageTypeOffer))
}

// Request completes the 4-way Discover-Offer-Request-Ack handshake.
//
// Note that modifiers will be applied *both* to Discover and Request packets.
func (c *Client) Request(ctx context.Context, modifiers ...dhcpv4.Modifier) (offer, ack *dhcpv4.DHCPv4, err error) {
	offer, err = c.DiscoverOffer(ctx, modifiers...)
	if err != nil {
		return nil, nil, err
	}

	// TODO(chrisko): should this be unicast to the server?
	req, err := dhcpv4.NewRequestFromOffer(offer, dhcpv4.PrependModifiers(modifiers,
		dhcpv4.WithOption(dhcpv4.OptMaxMessageSize(maxMessageSize)))...)
	if err != nil {
		return nil, nil, err
	}
	ack, err = c.SendAndRead(ctx, c.serverAddr, req, nil)
	if err != nil {
		return nil, nil, err
	}
	return offer, ack, nil
}

// send sends p to destination and returns a response channel.
//
// Responses will be matched by transaction ID and ClientHWAddr.
//
// The returned lambda function must be called after all desired responses have
// been received in order to return the Transaction ID to the usable pool.
func (c *Client) send(dest *net.UDPAddr, msg *dhcpv4.DHCPv4) (resp <-chan *dhcpv4.DHCPv4, cancel func(), err error) {
	c.pendingMu.Lock()
	if _, ok := c.pending[msg.TransactionID]; ok {
		c.pendingMu.Unlock()
		return nil, nil, fmt.Errorf("transaction ID %s already in use", msg.TransactionID)
	}

	ch := make(chan *dhcpv4.DHCPv4, c.bufferCap)
	done := make(chan struct{})
	c.pending[msg.TransactionID] = &pendingCh{done: done, ch: ch}
	c.pendingMu.Unlock()

	cancel = func() {
		// Why can't we just close ch here?
		//
		// Because receiveLoop may potentially be blocked trying to
		// send on ch. We gotta unblock it first, and then we can take
		// the lock and remove the XID from the pending transaction
		// map.
		close(done)

		c.pendingMu.Lock()
		if p, ok := c.pending[msg.TransactionID]; ok {
			close(p.ch)
			delete(c.pending, msg.TransactionID)
		}
		c.pendingMu.Unlock()
	}

	if _, err := c.conn.WriteTo(msg.ToBytes(), dest); err != nil {
		cancel()
		return nil, nil, fmt.Errorf("error writing packet to connection: %v", err)
	}
	return ch, cancel, nil
}

// This error should never be visible to users.
// It is used only to increase the timeout in retryFn.
var errDeadlineExceeded = errors.New("INTERNAL ERROR: deadline exceeded")

// SendAndRead sends a packet p to a destination dest and waits for the first
// response matching `match` as well as its Transaction ID and ClientHWAddr.
//
// If match is nil, the first packet matching the Transaction ID and
// ClientHWAddr is returned.
func (c *Client) SendAndRead(ctx context.Context, dest *net.UDPAddr, p *dhcpv4.DHCPv4, match Matcher) (*dhcpv4.DHCPv4, error) {
	var response *dhcpv4.DHCPv4
	err := c.retryFn(func(timeout time.Duration) error {
		ch, rem, err := c.send(dest, p)
		if err != nil {
			return err
		}
		c.logger.PrintMessage("sent message", p)
		defer rem()

		for {
			select {
			case <-c.done:
				return ErrNoResponse

			case <-time.After(timeout):
				return errDeadlineExceeded

			case <-ctx.Done():
				return ctx.Err()

			case packet := <-ch:
				if match == nil || match(packet) {
					c.logger.PrintMessage("received message", packet)
					response = packet
					return nil
				}
			}
		}
	})
	if err == errDeadlineExceeded {
		return nil, ErrNoResponse
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Client) retryFn(fn func(timeout time.Duration) error) error {
	timeout := c.timeout

	// Each retry takes the amount of timeout at worst.
	for i := 0; i < c.retry || c.retry < 0; i++ {
		switch err := fn(timeout); err {
		case nil:
			// Got it!
			return nil

		case errDeadlineExceeded:
			// Double timeout, then retry.
			timeout *= 2

		default:
			return err
		}
	}

	return errDeadlineExceeded
}
