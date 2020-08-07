package async

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	promise "github.com/fanliao/go-promise"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/client4"
)

// Default ports
const (
	DefaultServerPort = 67
	DefaultClientPort = 68
)

// Client implements an asynchronous DHCPv4 client
// It doesn't use the broadcast socket! Which means it should be used only when
// the network is already established.
// https://github.com/insomniacslk/dhcp/issues/143
type Client struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	LocalAddr    net.Addr
	RemoteAddr   net.Addr
	IgnoreErrors bool

	connection   *net.UDPConn
	cancel       context.CancelFunc
	stopping     *sync.WaitGroup
	receiveQueue chan *dhcpv4.DHCPv4
	sendQueue    chan *dhcpv4.DHCPv4
	packetsLock  sync.Mutex
	packets      map[dhcpv4.TransactionID]*promise.Promise
	errors       chan error
}

// NewClient creates an asynchronous client
func NewClient() *Client {
	return &Client{
		ReadTimeout:  client4.DefaultReadTimeout,
		WriteTimeout: client4.DefaultWriteTimeout,
	}
}

// Open starts the client. The requests made with Send function call are first
// put to the buffered channel and dispatched in FIFO order. BufferSize
// indicates the number of packets that can be waiting to be send before
// blocking the caller exectution.
func (c *Client) Open(bufferSize int) error {
	var (
		addr *net.UDPAddr
		ok   bool
		err  error
	)

	if addr, ok = c.LocalAddr.(*net.UDPAddr); !ok {
		return fmt.Errorf("Invalid local address: %v not a net.UDPAddr", c.LocalAddr)
	}

	// prepare the socket to listen on for replies
	c.connection, err = net.ListenUDP("udp4", addr)
	if err != nil {
		return err
	}
	c.stopping = new(sync.WaitGroup)
	c.sendQueue = make(chan *dhcpv4.DHCPv4, bufferSize)
	c.receiveQueue = make(chan *dhcpv4.DHCPv4, bufferSize)
	c.packets = make(map[dhcpv4.TransactionID]*promise.Promise)
	c.packetsLock = sync.Mutex{}
	c.errors = make(chan error)

	var ctx context.Context
	ctx, c.cancel = context.WithCancel(context.Background())
	go c.receiverLoop(ctx)
	go c.senderLoop(ctx)

	return nil
}

// Close stops the client
func (c *Client) Close() {
	// Wait for sender and receiver loops
	c.stopping.Add(2)
	c.cancel()
	c.stopping.Wait()

	close(c.sendQueue)
	close(c.receiveQueue)
	close(c.errors)

	c.connection.Close()
}

// Errors returns a channel where runtime errors are posted
func (c *Client) Errors() <-chan error {
	return c.errors
}

func (c *Client) addError(err error) {
	if !c.IgnoreErrors {
		c.errors <- err
	}
}

func (c *Client) receiverLoop(ctx context.Context) {
	defer func() { c.stopping.Done() }()
	for {
		select {
		case <-ctx.Done():
			return
		case packet := <-c.receiveQueue:
			c.receive(packet)
		}
	}
}

func (c *Client) senderLoop(ctx context.Context) {
	defer func() { c.stopping.Done() }()
	for {
		select {
		case <-ctx.Done():
			return
		case packet := <-c.sendQueue:
			c.send(packet)
		}
	}
}

func (c *Client) send(packet *dhcpv4.DHCPv4) {
	c.packetsLock.Lock()
	p := c.packets[packet.TransactionID]
	c.packetsLock.Unlock()

	raddr, err := c.remoteAddr()
	if err != nil {
		_ = p.Reject(err)
		return
	}

	if err := c.connection.SetWriteDeadline(time.Now().Add(c.WriteTimeout)); err != nil {
		log.Printf("Warning: cannot set write deadline: %v", err)
		return
	}
	_, err = c.connection.WriteTo(packet.ToBytes(), raddr)
	if err != nil {
		_ = p.Reject(err)
		log.Printf("Warning: cannot write to %s: %v", raddr, err)
		return
	}

	c.receiveQueue <- packet
}

func (c *Client) receive(_ *dhcpv4.DHCPv4) {
	var (
		oobdata  = []byte{}
		received *dhcpv4.DHCPv4
	)

	if err := c.connection.SetReadDeadline(time.Now().Add(c.ReadTimeout)); err != nil {
		log.Printf("Warning: cannot set write deadline: %v", err)
		return
	}
	for {
		buffer := make([]byte, client4.MaxUDPReceivedPacketSize)
		n, _, _, _, err := c.connection.ReadMsgUDP(buffer, oobdata)
		if err != nil {
			if err, ok := err.(net.Error); !ok || !err.Timeout() {
				c.addError(fmt.Errorf("Error receiving the message: %s", err))
			}
			return
		}
		received, err = dhcpv4.FromBytes(buffer[:n])
		if err == nil {
			break
		}
	}

	c.packetsLock.Lock()
	if p, ok := c.packets[received.TransactionID]; ok {
		delete(c.packets, received.TransactionID)
		_ = p.Resolve(received)
	}
	c.packetsLock.Unlock()
}

func (c *Client) remoteAddr() (*net.UDPAddr, error) {
	if c.RemoteAddr == nil {
		return &net.UDPAddr{IP: net.IPv4bcast, Port: DefaultServerPort}, nil
	}

	if addr, ok := c.RemoteAddr.(*net.UDPAddr); ok {
		return addr, nil
	}
	return nil, fmt.Errorf("Invalid remote address: %v not a net.UDPAddr", c.RemoteAddr)
}

// Send inserts a message to the queue to be sent asynchronously.
// Returns a future which resolves to response and error.
func (c *Client) Send(message *dhcpv4.DHCPv4, modifiers ...dhcpv4.Modifier) *promise.Future {
	for _, mod := range modifiers {
		mod(message)
	}

	p := promise.NewPromise()
	c.packetsLock.Lock()
	c.packets[message.TransactionID] = p
	c.packetsLock.Unlock()
	c.sendQueue <- message
	return p.Future
}
