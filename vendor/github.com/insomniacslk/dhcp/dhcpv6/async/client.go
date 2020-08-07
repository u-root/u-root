package async

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	promise "github.com/fanliao/go-promise"
	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/dhcpv6/client6"
)

// Client implements an asynchronous DHCPv6 client
type Client struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	LocalAddr    net.Addr
	RemoteAddr   net.Addr
	IgnoreErrors bool

	connection   *net.UDPConn
	cancel       context.CancelFunc
	stopping     *sync.WaitGroup
	receiveQueue chan dhcpv6.DHCPv6
	sendQueue    chan dhcpv6.DHCPv6
	packetsLock  sync.Mutex
	packets      map[dhcpv6.TransactionID]*promise.Promise
	errors       chan error
}

// NewClient creates an asynchronous client
func NewClient() *Client {
	return &Client{
		ReadTimeout:  client6.DefaultReadTimeout,
		WriteTimeout: client6.DefaultWriteTimeout,
	}
}

// OpenForInterface starts the client on the specified interface, replacing
// client LocalAddr with a link-local address of the given interface and
// standard DHCP port (546).
func (c *Client) OpenForInterface(ifname string, bufferSize int) error {
	addr, err := dhcpv6.GetLinkLocalAddr(ifname)
	if err != nil {
		return err
	}
	c.LocalAddr = &net.UDPAddr{IP: addr, Port: dhcpv6.DefaultClientPort, Zone: ifname}
	return c.Open(bufferSize)
}

// Open starts the client
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
	c.connection, err = net.ListenUDP("udp6", addr)
	if err != nil {
		return err
	}
	c.stopping = new(sync.WaitGroup)
	c.sendQueue = make(chan dhcpv6.DHCPv6, bufferSize)
	c.receiveQueue = make(chan dhcpv6.DHCPv6, bufferSize)
	c.packets = make(map[dhcpv6.TransactionID]*promise.Promise)
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

func (c *Client) send(packet dhcpv6.DHCPv6) {
	transactionID, err := dhcpv6.GetTransactionID(packet)
	if err != nil {
		c.addError(fmt.Errorf("Warning: This should never happen, there is no transaction ID on %s", packet))
		return
	}
	c.packetsLock.Lock()
	p := c.packets[transactionID]
	c.packetsLock.Unlock()

	raddr, err := c.remoteAddr()
	if err != nil {
		_ = p.Reject(err)
		log.Printf("Warning: cannot get remote address :%v", err)
		return
	}

	if err := c.connection.SetWriteDeadline(time.Now().Add(c.WriteTimeout)); err != nil {
		_ = p.Reject(err)
		log.Printf("Warning: cannot set write deadline :%v", err)
		return
	}
	_, err = c.connection.WriteTo(packet.ToBytes(), raddr)
	if err != nil {
		_ = p.Reject(err)
		log.Printf("Warning: cannot write to %s :%v", raddr, err)
		return
	}

	c.receiveQueue <- packet
}

func (c *Client) receive(_ dhcpv6.DHCPv6) {
	var (
		oobdata  = []byte{}
		received dhcpv6.DHCPv6
	)

	if err := c.connection.SetReadDeadline(time.Now().Add(c.ReadTimeout)); err != nil {
		log.Printf("Warning: cannot set read deadline :%v", err)
	}
	for {
		buffer := make([]byte, client6.MaxUDPReceivedPacketSize)
		n, _, _, _, err := c.connection.ReadMsgUDP(buffer, oobdata)
		if err != nil {
			if err, ok := err.(net.Error); !ok || !err.Timeout() {
				c.addError(fmt.Errorf("Error receiving the message: %s", err))
			}
			return
		}
		received, err = dhcpv6.FromBytes(buffer[:n])
		if err != nil {
			// skip non-DHCP packets
			continue
		}
		break
	}

	transactionID, err := dhcpv6.GetTransactionID(received)
	if err != nil {
		c.addError(fmt.Errorf("Unable to get a transactionID for %s: %s", received, err))
		return
	}

	c.packetsLock.Lock()
	if p, ok := c.packets[transactionID]; ok {
		delete(c.packets, transactionID)
		_ = p.Resolve(received)
	}
	c.packetsLock.Unlock()
}

func (c *Client) remoteAddr() (*net.UDPAddr, error) {
	if c.RemoteAddr == nil {
		return &net.UDPAddr{IP: dhcpv6.AllDHCPRelayAgentsAndServers, Port: dhcpv6.DefaultServerPort}, nil
	}

	if addr, ok := c.RemoteAddr.(*net.UDPAddr); ok {
		return addr, nil
	}
	return nil, fmt.Errorf("Invalid remote address: %v not a net.UDPAddr", c.RemoteAddr)
}

// Send inserts a message to the queue to be sent asynchronously.
// Returns a future which resolves to response and error.
func (c *Client) Send(message dhcpv6.DHCPv6, modifiers ...dhcpv6.Modifier) *promise.Future {
	for _, mod := range modifiers {
		mod(message)
	}

	transactionID, err := dhcpv6.GetTransactionID(message)
	if err != nil {
		return promise.Wrap(err)
	}

	p := promise.NewPromise()
	c.packetsLock.Lock()
	c.packets[transactionID] = p
	c.packetsLock.Unlock()
	c.sendQueue <- message
	return p.Future
}
