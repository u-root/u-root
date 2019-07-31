package client6

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv6"
)

// Client constants
const (
	DefaultWriteTimeout       = 3 * time.Second // time to wait for write calls
	DefaultReadTimeout        = 3 * time.Second // time to wait for read calls
	DefaultInterfaceUpTimeout = 3 * time.Second // time to wait before a network interface goes up
	MaxUDPReceivedPacketSize  = 8192            // arbitrary size. Theoretically could be up to 65kb
)

// Client implements a DHCPv6 client
type Client struct {
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	LocalAddr     net.Addr
	RemoteAddr    net.Addr
	SimulateRelay bool
	RelayOptions  dhcpv6.Options // These options will be added to relay message if SimulateRelay is true
}

// NewClient returns a Client with default settings
func NewClient() *Client {
	return &Client{
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
	}
}

// Exchange executes a 4-way DHCPv6 request (Solicit, Advertise, Request,
// Reply). The modifiers will be applied to the Solicit and Request packets.
// A common use is to make sure that the Solicit packet has the right options,
// see modifiers.go
func (c *Client) Exchange(ifname string, modifiers ...dhcpv6.Modifier) ([]dhcpv6.DHCPv6, error) {
	conversation := make([]dhcpv6.DHCPv6, 0)
	var err error

	// Solicit
	solicit, advertise, err := c.Solicit(ifname, modifiers...)
	if solicit != nil {
		conversation = append(conversation, solicit)
	}
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, advertise)

	// Decapsulate advertise if it's relayed before passing it to Request
	if advertise.IsRelay() {
		advertiseRelay := advertise.(*dhcpv6.RelayMessage)
		advertise, err = advertiseRelay.GetInnerMessage()
		if err != nil {
			return conversation, err
		}
	}
	request, reply, err := c.Request(ifname, advertise.(*dhcpv6.Message), modifiers...)
	if request != nil {
		conversation = append(conversation, request)
	}
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, reply)
	return conversation, nil
}

func (c *Client) sendReceive(ifname string, packet dhcpv6.DHCPv6, expectedType dhcpv6.MessageType) (dhcpv6.DHCPv6, error) {
	if packet == nil {
		return nil, fmt.Errorf("Packet to send cannot be nil")
	}
	// if no LocalAddr is specified, get the interface's link-local address
	var laddr net.UDPAddr
	if c.LocalAddr == nil {
		llAddr, err := dhcpv6.GetLinkLocalAddr(ifname)
		if err != nil {
			return nil, err
		}
		laddr = net.UDPAddr{IP: llAddr, Port: dhcpv6.DefaultClientPort, Zone: ifname}
	} else {
		if addr, ok := c.LocalAddr.(*net.UDPAddr); ok {
			laddr = *addr
		} else {
			return nil, fmt.Errorf("Invalid local address: not a net.UDPAddr: %v", c.LocalAddr)
		}
	}
	if c.SimulateRelay {
		var err error
		packet, err = dhcpv6.EncapsulateRelay(packet, dhcpv6.MessageTypeRelayForward, net.IPv6zero, laddr.IP)
		if err != nil {
			return nil, err
		}
		// Add Relay Options to ecapsulated Packet
		for _, opt := range c.RelayOptions {
			packet.UpdateOption(opt)
		}
	}
	if expectedType == dhcpv6.MessageTypeNone {
		// infer the expected type from the packet being sent
		if packet.Type() == dhcpv6.MessageTypeSolicit {
			expectedType = dhcpv6.MessageTypeAdvertise
		} else if packet.Type() == dhcpv6.MessageTypeRequest {
			expectedType = dhcpv6.MessageTypeReply
		} else if packet.Type() == dhcpv6.MessageTypeRelayForward {
			expectedType = dhcpv6.MessageTypeRelayReply
		} else if packet.Type() == dhcpv6.MessageTypeLeaseQuery {
			expectedType = dhcpv6.MessageTypeLeaseQueryReply
		} // and probably more
	}

	// if no RemoteAddr is specified, use AllDHCPRelayAgentsAndServers
	var raddr net.UDPAddr
	if c.RemoteAddr == nil {
		raddr = net.UDPAddr{IP: dhcpv6.AllDHCPRelayAgentsAndServers, Port: dhcpv6.DefaultServerPort}
	} else {
		if addr, ok := c.RemoteAddr.(*net.UDPAddr); ok {
			raddr = *addr
		} else {
			return nil, fmt.Errorf("Invalid remote address: not a net.UDPAddr: %v", c.RemoteAddr)
		}
	}

	// prepare the socket to listen on for replies
	conn, err := net.ListenUDP("udp6", &laddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	// wait for the listener to be ready, fail if it takes too much time
	deadline := time.Now().Add(time.Second)
	for {
		if now := time.Now(); now.After(deadline) {
			return nil, errors.New("Timed out waiting for listener to be ready")
		}
		if conn.LocalAddr() != nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// send the packet out
	if err := conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout)); err != nil {
		return nil, err
	}
	_, err = conn.WriteTo(packet.ToBytes(), &raddr)
	if err != nil {
		return nil, err
	}

	// wait for a reply
	oobdata := []byte{} // ignoring oob data
	if err := conn.SetReadDeadline(time.Now().Add(c.ReadTimeout)); err != nil {
		return nil, err
	}
	var (
		adv       dhcpv6.DHCPv6
		isMessage bool
	)
	defer conn.Close()
	msg, ok := packet.(*dhcpv6.Message)
	if ok {
		isMessage = true
	}
	for {
		buf := make([]byte, MaxUDPReceivedPacketSize)
		n, _, _, _, err := conn.ReadMsgUDP(buf, oobdata)
		if err != nil {
			return nil, err
		}
		adv, err = dhcpv6.FromBytes(buf[:n])
		if err != nil {
			// skip non-DHCP packets
			continue
		}
		if recvMsg, ok := adv.(*dhcpv6.Message); ok && isMessage {
			// if a regular message, check the transaction ID first
			// XXX should this unpack relay messages and check the XID of the
			// inner packet too?
			if msg.TransactionID != recvMsg.TransactionID {
				// different XID, we don't want this packet for sure
				continue
			}
		}
		if expectedType == dhcpv6.MessageTypeNone {
			// just take whatever arrived
			break
		} else if adv.Type() == expectedType {
			break
		}
	}
	return adv, nil
}

// Solicit sends a Solicit, returns the Solicit, an Advertise (if not nil), and
// an error if any. The modifiers will be applied to the Solicit before sending
// it, see modifiers.go
func (c *Client) Solicit(ifname string, modifiers ...dhcpv6.Modifier) (dhcpv6.DHCPv6, dhcpv6.DHCPv6, error) {
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return nil, nil, err
	}
	solicit, err := dhcpv6.NewSolicit(iface.HardwareAddr)
	if err != nil {
		return nil, nil, err
	}
	for _, mod := range modifiers {
		mod(solicit)
	}
	advertise, err := c.sendReceive(ifname, solicit, dhcpv6.MessageTypeNone)
	return solicit, advertise, err
}

// Request sends a Request built from an Advertise. It returns the Request, a
// Reply (if not nil), and an error if any. The modifiers will be applied to
// the Request before sending it, see modifiers.go
func (c *Client) Request(ifname string, advertise *dhcpv6.Message, modifiers ...dhcpv6.Modifier) (dhcpv6.DHCPv6, dhcpv6.DHCPv6, error) {
	request, err := dhcpv6.NewRequestFromAdvertise(advertise)
	if err != nil {
		return nil, nil, err
	}
	for _, mod := range modifiers {
		mod(request)
	}
	reply, err := c.sendReceive(ifname, request, dhcpv6.MessageTypeNone)
	return request, reply, err
}
