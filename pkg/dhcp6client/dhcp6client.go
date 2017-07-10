package dhcp6client

import (
	"fmt"
	"net"

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

	// Max number of attempts to receive a valid DHCPv6 reply from server.
	attempts int
}

func New(haddr net.HardwareAddr, packetSock *packetSock, n int) *Client {
	return &Client{
		srcMAC:     haddr,
		connection: packetSock,
		attempts:   n,
	}
}

func (c *Client) Solicit() ([]*dhcp6.IAAddr, *dhcp6.Packet, error) {
	solicitPacket, err := newSolicitPacket(c.srcMAC)
	if err != nil {
		return nil, nil, fmt.Errorf("Request Error:\nnew solicit packet: %v\nerr: %v", solicitPacket, err)
	}

	if err := c.SendPacket(solicitPacket); err != nil {
		return nil, nil, fmt.Errorf("Request Error:\nsend solicit packet: %v\nerr: %v", solicitPacket, err)
	}

	packet, err := c.ReadReply()
	if err != nil {
		return nil, nil, fmt.Errorf("Request Error:\nadvertise packet: %v\nerr: %v", packet, err)
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
		return nil, packet, fmt.Errorf("error: reply does not contain valid Iaaddr: %v", err)
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
	for i := 0; i < c.attempts; i++ { // five attempts
		pb := make([]byte, 1500)
		var n int
		n, err = c.connection.ReadFrom(pb)
		if err != nil {
			continue
		}

		p := &buffer{pb[:n]}
		ipv6 := header.IPv6(p.consume(header.IPv6MinimumSize))

		if ipv6.NextHeader() == uint8(header.UDPProtocolNumber) {
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
