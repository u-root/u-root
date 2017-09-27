package dhcp6

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/header"
)

const (
	srcPort = 546
	dstPort = 547
)

var (
	// All DHCP servers and relay agents on the local network segment (RFC 3315)
	// IPv6 Multicast (RFC 2464)
	// insert the low 32 Bits of the multicast IPv6 Address into the Ethernet Address (RFC 7042 2.3.1.)
	multicastMAC = net.HardwareAddr([]byte{0x33, 0x33, 0x00, 0x01, 0x00, 0x02})
)

type Client struct {
	// The HardwareAddr to send in the request.
	srcMAC net.HardwareAddr

	// Packet socket to send on.
	connection *packetSock

	// Timeout
	timeout time.Duration

	// Max number of attempts to send DHCPv6 solicit to server.
	// A valid DHCPv6 reply is supposed to be received by client before timeout.
	// -1 means infinity.
	retry int
}

func New(haddr net.HardwareAddr, packetSock *packetSock, t time.Duration, r int) *Client {
	return &Client{
		srcMAC:     haddr,
		connection: packetSock,
		timeout:    t,
		retry:      r,
	}
}

func (c *Client) Solicit() ([]*IAAddr, *Packet, error) {
	solicitPacket, err := newSolicitPacket(c.srcMAC)
	if err != nil {
		return nil, nil, fmt.Errorf("new solicit packet: %v", err)
	}

	var packet *Packet
	for i := 0; i < c.retry || c.retry < 0; i++ { // each retry takes the amount of timeout at worst.
		if err := c.SendPacket(solicitPacket); err != nil {
			return nil, nil, fmt.Errorf("send solicit packet(%v) = err %v", solicitPacket, err)
		}

		packet, err = c.ReadReply()
		if err != nil {
			log.Printf("%v\nResending DHCPv6 Solicit Message...", err)
			continue
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

func (c *Client) SendPacket(p *Packet) error {
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

// If the client fails to receive a valid DHCPv6 packet via socket,
// it keeps listening until the time is out.
func (c *Client) ReadReply() (*Packet, error) {
	start := time.Now()

	for {
		remainingTime := c.timeout - time.Since(start)
		if remainingTime <= 0 {
			return nil, fmt.Errorf("waiting for response timed out")
		}

		c.connection.SetReadTimeout(remainingTime)
		pb := make([]byte, 1500)
		var n int
		n, err := c.connection.ReadFrom(pb)
		if err != nil {
			continue
		}

		p := &buffer{pb[:n]}
		ipv6 := header.IPv6(p.consume(header.IPv6MinimumSize))

		if ipv6.DestinationAddress() == tcpip.Address(mac2ipv6(c.srcMAC)) ||
			ipv6.NextHeader() == uint8(header.UDPProtocolNumber) {
			udp := header.UDP(p.consume(header.UDPMinimumSize))

			if udp.DestinationPort() == srcPort {
				dhcp6p := &Packet{}
				if err := dhcp6p.UnmarshalBinary(p.remaining()); err != nil {
					// Not a valid DHCPv6 reply; keep listening.
					continue
				}
				return dhcp6p, nil
			}
		}
	}
}
