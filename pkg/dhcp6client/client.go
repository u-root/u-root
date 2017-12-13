// Package dhcp6client implements a DHCPv6 client as per RFC 3315.
package dhcp6client

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/mdlayher/dhcp6"
	"github.com/mdlayher/dhcp6/dhcp6opts"
	"github.com/mdlayher/eui64"
	"github.com/vishvananda/netlink"
)

// RFC 3315 Section 5.2.
const (
	// ClientPort is the port clients use to listen for DHCP messages.
	ClientPort = 546

	// ServerPort is the port servers and relay agents use to listen for
	// DHCP messages.
	ServerPort = 547
)

var (
	// All DHCP servers and relay agents on the local network segment (RFC
	// 3315, Section 5.1.).
	AllServers = net.ParseIP("ff02::1:2")

	// DefaultServers is the default AllServers IP combined with the
	// ServerPort.
	DefaultServers = &net.UDPAddr{
		IP:   AllServers,
		Port: ServerPort,
	}
)

// Client is a simple DHCPv6 client implementing RFC 3315.
//
//
// Shortest Example:
//
//  c, err := dhcp6client.New(iface, 10*time.Second, 2)
//  ...
//  iana, packet, err := c.RapidSolicit()
//  ...
//  // iana now contains the IP assigned in the IAAddr option.
//
//
// Example selecting which advertising server to request from:
//
//   c, err := dhcp6client.New(iface, 10*time.Second, 2)
//   ...
//   ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
//   defer cancel()
//
//   ads, err := c.Solicit(ctx)
//   ...
//   // Selecting the advertisement of server 3.
//   request, err := dhcp6client.RequestIANAFrom(ads[2])
//   ...
//   iana, packet, err := c.RequestOne(request)
//   ...
//   // iana now contains the IP assigned in the IAAddr option.
type Client struct {
	// The HardwareAddr to send in the request.
	srcMAC net.HardwareAddr

	iface netlink.Link

	// Packet socket to send on.
	conn net.PacketConn

	// Max number of attempts to multicast DHCPv6 solicits.
	// -1 means infinity.
	retry int

	// Timeout for each Solicit try.
	timeout time.Duration
}

// New returns a new DHCPv6 client based on the given parameters.
func New(iface netlink.Link, timeout time.Duration, retry int) (*Client, error) {
	haddr := iface.Attrs().HardwareAddr
	ip, err := eui64.ParseMAC(net.ParseIP("fe80::"), haddr)
	if err != nil {
		return nil, err
	}

	// If this doesn't work like this, we may have to SO_BINDTODEVICE to
	// the interface. If that doesn't work... then we're back to raw
	// sockets. Hope not.
	c, err := net.ListenUDP("udp6", &net.UDPAddr{
		IP:   ip,
		Port: ClientPort,
		Zone: iface.Attrs().Name,
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		iface:   iface,
		srcMAC:  haddr,
		conn:    c,
		timeout: timeout,
		retry:   retry,
	}, nil
}

// RapidSolicit solicits one non-temporary address assignment by multicasting a
// DHCPv6 solicitation message with the rapid commit option.
//
// RapidSolicit returns the first valid, suitable response by any remote server.
func (c *Client) RapidSolicit() (*dhcp6opts.IANA, *dhcp6.Packet, error) {
	solicit, err := NewRapidSolicit(c.srcMAC)
	if err != nil {
		return nil, nil, err
	}
	return c.RequestOne(solicit)
}

// RequestOne multicasts the `request` and returns the first matching IANA and
// its associated Packet returned by any server.
func (c *Client) RequestOne(request *dhcp6.Packet) (*dhcp6opts.IANA, *dhcp6.Packet, error) {
	ianas, pkt, err := c.Request(request)
	if err != nil {
		return nil, nil, err
	}
	if len(ianas) != 1 {
		return nil, nil, fmt.Errorf("got %d IANAs, expected 1", len(ianas))
	}
	return ianas[0], pkt, nil
}

// Solicit multicasts a Solicit message and collects all Advertise responses
// received before c.timeout expires.
//
// Solicit blocks until either:
// - `ctx` is canceled; or
// - we have exhausted all configured retries and timeouts.
func (c *Client) Solicit(ctx context.Context) ([]*dhcp6.Packet, error) {
	solicit, err := NewSolicitPacket(c.srcMAC)
	if err != nil {
		return nil, err
	}

	out, errCh := c.SendAndRead(ctx, DefaultServers, solicit)

	var ads []*dhcp6.Packet
	// resps is closed by SendAndRead when done.
	for r := range out {
		if r.Packet.MessageType == dhcp6.MessageTypeAdvertise {
			ads = append(ads, r.Packet)
		}
	}

	if err, ok := <-errCh; ok && err != nil {
		return nil, err
	}
	return ads, nil
}

// This name smells.
type manyErrs []string

func newManyErrs() *manyErrs {
	return new(manyErrs)
}

func (e *manyErrs) add(err error) {
	*e = append(*e, err.Error())
}

func (e manyErrs) Error() string {
	return strings.Join([]string(e), "; ")
}

// Request requests non-temporary address assignments by multicasting the given
// message.
//
// This request message may be any DHCPv6 request message type; e.g. a
// Solicit with the Rapid Commit option or a Rebind message.
func (c *Client) Request(request *dhcp6.Packet) ([]*dhcp6opts.IANA, *dhcp6.Packet, error) {
	errs := newManyErrs()

	// These are the IANAs we are looking for in responses.
	reqIANAs, err := dhcp6opts.GetIANA(request.Options)
	if err != nil {
		return nil, nil, fmt.Errorf("request packet contains no IANAs: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out, errCh := c.SendAndRead(ctx, DefaultServers, request)

	for packet := range out {
		if ianas, err := SuitableReply(reqIANAs, packet.Packet); err != nil {
			errs.add(err)
		} else {
			// Guess we found our IANAs! The context will cancel
			// all our problems.
			return ianas, packet.Packet, nil
		}
	}

	// Check if an error occured.
	if err, ok := <-errCh; ok && err != nil {
		errs.add(err)
	}

	errs.add(fmt.Errorf("no suitable responses"))
	return nil, nil, errs
}

type ClientPacket struct {
	Interface netlink.Link
	Packet    *dhcp6.Packet
}

type ClientError struct {
	Interface netlink.Link
	Err       error
}

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

// SendAndRead multicasts a DHCPv6 packet and launches a goroutine to read
// response packets. Those response packets will be sent on the channel
// returned. The sender will close both goroutines when it stops reading
// packets, for example when the context is canceled.
//
// Callers must cancel ctx when they have received the packet they are looking
// for. Otherwise, the spawned goroutine will keep reading until it times out.
// More importantly, if you send another packet, the spawned goroutine may read
// the response faster than the one launched for the other packet.
//
// See Client.Solicit for an example use of SendAndRead.
//
// Callers sending a packet on one interface should use this. Callers intending
// to send packets on many interface at the same time, should look at using
// ParallelSendAndRead instead.
//
// Example Usage:
//
//   func sendRequest(someRequest *Packet...) (*Packet, error) {
//     ctx, cancel := context.WithCancel(context.Background())
//     defer cancel()
//
//     out, errCh := c.SendAndRead(ctx, DefaultServers, someRequest)
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
func (c *Client) SendAndRead(ctx context.Context, dest *net.UDPAddr, p *dhcp6.Packet) (<-chan *ClientPacket, <-chan *ClientError) {
	out := make(chan *ClientPacket, 10)
	errOut := make(chan *ClientError, 1)
	go c.ParallelSendAndRead(ctx, dest, p, out, errOut)
	return out, errOut
}

// ParallelSendAndRead sends the given packet `dest` to `to` and reads
// responses on the UDP connection. Valid responses are sent to `out`; `out` is
// closed by SendAndRead when it returns.
//
// ParallelSendAndRead blocks reading response packets until either:
// - `ctx` is canceled; or
// - we have exhausted all configured retries and timeouts.
//
// Any valid DHCPv6 packet received with the correct Transaction ID is sent on
// `out`.
//
// SendAndRead retries sending the packet and receiving responses according to
// the configured number of c.retry, using a response timeout of c.timeout.
//
// TODO(hugelgupf): SendAndRead should follow RFC 3315 Section 14 for
// retransmission behavior. Also conform to Section 15 for what kind of
// messages must be discarded.
func (c *Client) ParallelSendAndRead(ctx context.Context, dest *net.UDPAddr, p *dhcp6.Packet, out chan<- *ClientPacket, errCh chan<- *ClientError) {
	defer close(errCh)

	// This ensures that
	// - we send at most one error on errCh; and
	// - we don't forget to send err on errCh in the many return statements
	//   of sendAndRead.
	if err := c.sendAndRead(ctx, dest, p, out); err != nil {
		errCh <- err
	}
}

func (c *Client) sendAndRead(ctx context.Context, dest *net.UDPAddr, p *dhcp6.Packet, out chan<- *ClientPacket) *ClientError {
	defer close(out)

	pkt, err := p.MarshalBinary()
	if err != nil {
		return c.newClientErr(err)
	}

	return c.retryFn(func() *ClientError {
		if _, err := c.conn.WriteTo(pkt, dest); err != nil {
			return c.newClientErr(fmt.Errorf("error writing packet to connection: %v", err))
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
				return c.newClientErr(timeoutCtx.Err())
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
			if oerr, ok := err.(*net.OpError); ok && oerr.Timeout() {
				// Continue to check ctx.Done() above and
				// return the appropriate error.
				continue
			} else if err != nil {
				return c.newClientErr(fmt.Errorf("error reading from UDP connection: %v", err))
			}

			pkt := &dhcp6.Packet{}
			if err := pkt.UnmarshalBinary(b[:n]); err != nil {
				// Not a valid DHCPv6 reply; keep listening.
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
				return c.newClientErr(ctx.Err())
			case out <- clientPkt:
			}
		}
	})
}

// SuitableReply validates whether a pkt is a valid Reply message as defined by
// RFC 3315, Section 18.1.8.
//
// It returns all valid IANAs corresponding to requested IANAs.
func SuitableReply(reqIANAs []*dhcp6opts.IANA, pkt *dhcp6.Packet) ([]*dhcp6opts.IANA, error) {
	// RFC 3315, Section 18.1.8.
	// A suitable Reply packet must have:
	//
	// - non-negative status code (or no status), and
	// - an IANA with IAID matching one of the ones we used in our request, and
	// -- a non-negative status code (or no status) in the matching IANA, and
	// -- a non-zero number of IAAddrs in the matching IANA.
	if pkt.MessageType != dhcp6.MessageTypeReply {
		return nil, fmt.Errorf("got DHCP message of type %s, wanted %s", pkt.MessageType, dhcp6.MessageTypeReply)
	}

	if status, err := dhcp6opts.GetStatusCode(pkt.Options); err == nil && status.Code != dhcp6.StatusSuccess {
		return nil, fmt.Errorf("packet has status %s: %s", status.Code, status.Message)
	}

	ianas, err := dhcp6opts.GetIANA(pkt.Options)
	if err != nil {
		return nil, fmt.Errorf("successful packet had problem with IANA: %v", err)
	}

	var returned []*dhcp6opts.IANA
	for _, iana := range ianas {
		for _, reqIANA := range reqIANAs {
			if iana.IAID != reqIANA.IAID {
				continue
			}

			if status, err := dhcp6opts.GetStatusCode(iana.Options); err == nil && status.Code != dhcp6.StatusSuccess {
				continue
			}

			iaAddrs, err := dhcp6opts.GetIAAddr(iana.Options)
			if err != nil || len(iaAddrs) == 0 {
				continue
			}

			returned = append(returned, iana)
		}
	}

	return returned, nil
}

func (c *Client) retryFn(fn func() *ClientError) *ClientError {
	// Each retry takes the amount of timeout at worst.
	for i := 0; i < c.retry || c.retry < 0; i++ {
		switch err := fn(); {
		case err == nil:
			// Got it!
			return nil

		case err.Err == context.DeadlineExceeded:
			// Just retry.
			// TODO(hugelgupf): Sleep here for some random amount of time.

		default:
			return err
		}
	}

	return c.newClientErr(context.DeadlineExceeded)
}
