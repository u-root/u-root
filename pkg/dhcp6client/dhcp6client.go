package dhcp6client

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/header"
	"github.com/mdlayher/dhcp6"
)

const (
	srcPort = 546
	dstPort = 547
)

// All DHCP servers and relay agents on the local network segment (RFC 3315)
// IPv6 Multicast (RFC 2464)
// insert the low 32 Bits of the multicast IPv6 Address into the Ethernet Address (RFC 7042 2.3.1.)
var multicastMAC = net.HardwareAddr([]byte{0x33, 0x33, 0x00, 0x01, 0x00, 0x02})

type Client struct {
	// The HardwareAddr to send in the request.
	srcMAC net.HardwareAddr

	// Packet socket to send on.
	connection *packetSock

	// Timeout
	timeout time.Duration
}

func New(haddr net.HardwareAddr, packetSock *packetSock, t time.Duration) *Client {
	return &Client{
		srcMAC:     haddr,
		connection: packetSock,
		timeout:    t,
	}
}

func (c *Client) Solicit() ([]*dhcp6.IAAddr, *dhcp6.Packet, error) {
	solicitPacket, err := newSolicitPacket(c.srcMAC)
	if err != nil {
		return nil, nil, fmt.Errorf("new solicit packet: %v", err)
	}

	var packet *dhcp6.Packet
	for {
		if err := c.SendPacket(solicitPacket); err != nil {
			return nil, nil, fmt.Errorf("send solicit packet(%v) = err %v", solicitPacket, err)
		}

		packet, err = c.ReadReply()
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				log.Printf("%v: resending DHCPv6 Solicit Message...", err)
				continue
			}
			return nil, nil, fmt.Errorf("error while reading reply: %v", err)
		} else {
			break
		}
	}

	iana, containsIANA, err := packet.Options.IANA()
	if err != nil {
		return nil, packet, fmt.Errorf("error: reply does not contain valid IANA: %v", err)
	}
	if !containsIANA {
		return nil, packet, fmt.Errorf("error: reply does not contain IANA")
	}

	iaAddrs, containsIAAddr, err := iana[0].Options.IAAddr()
	if err != nil {
		return nil, packet, fmt.Errorf("error: reply does not contain valid IAAddr: %v", err)
	}
	if !containsIAAddr {
		return nil, packet, fmt.Errorf("error: reply does not contain IAAddr")
	}

	return iaAddrs, packet, nil
}

func (c *Client) SendPacket(p *dhcp6.Packet) error {
	pkt, err := ipv6UDPPacket(p, c.srcMAC)
	if err != nil {
		return err
	}

	return c.connection.WriteTo(pkt, multicastMAC)
}

type buffer struct {
	buf []byte
}

func (p *buffer) consume(size int) []byte {
	consumed := p.buf[:size]
	p.buf = p.buf[size:]
	return consumed
}

func (p *buffer) remaining() []byte {
	return p.consume(len(p.buf))
}

func (c *Client) ReadReply() (*dhcp6.Packet, error) {
	var err error
	for i := 0; i < 5; i++ {
		// set a random timeout within range [c.timeout, 2*c.timeout)
		// so as to prevent swamping a DHCP server.
		c.connection.SetReadTimeout(c.timeout + time.Duration(rand.Intn(int(c.timeout))))
		pb := make([]byte, 1500)
		var n int
		n, err = c.connection.ReadFrom(pb)
		if err != nil {
			// Return error if read time of socket is due,
			// so that request can be resent instantly
			if err, ok := err.(net.Error); ok && err.Timeout() {
				return nil, err
			}
			continue
		}

		p := &buffer{pb[:n]}
		ipv6 := header.IPv6(p.consume(header.IPv6MinimumSize))

		if ipv6.DestinationAddress() == tcpip.Address(mac2ipv6(c.srcMAC)) ||
			ipv6.NextHeader() == uint8(header.UDPProtocolNumber) {
			udp := header.UDP(p.consume(header.UDPMinimumSize))

			if udp.DestinationPort() == srcPort {
				dhcp6p := &dhcp6.Packet{}
				if err = dhcp6p.UnmarshalBinary(p.remaining()); err != nil {
					continue
				}
				return dhcp6p, nil
			}
		}
	}
	return nil, fmt.Errorf("failed to get ipv6 address after five attempts: %v", err)
}
