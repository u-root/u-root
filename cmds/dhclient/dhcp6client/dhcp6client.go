package dhcp6client

import (
	"fmt"
	"net"
	"time"

	"github.com/google/netstack/tcpip/header"
	"github.com/mdlayher/dhcp6"
)

const (
	ipv6HdrLen = 40
	udpHdrLen  = 8

	srcPort     = 546
	dstPort     = 547
	protocolUDP = 17
)

type Client struct {
	hardwareAddr  net.HardwareAddr //The HardwareAddr to send in the request.
	ignoreServers []net.IP         //List of Servers to Ignore requests from.
	timeout       time.Duration    //Time before we timeout.
	broadcast     bool             //Set the Bcast flag in BOOTP Flags
	connection    connection       //The Connection Method to use
}

/*
* Abstracts the type of underlying socket used
 */

func New(haddr net.HardwareAddr, conn connection, timeout time.Duration) (*Client, error) {
	c := Client{
		broadcast: true,
	}

	c.hardwareAddr = haddr
	c.connection = conn
	c.timeout = timeout
	return &c, nil
}

func (c *Client) Request(mac *net.HardwareAddr) (*dhcp6.Packet, error) {
	solicitPacket, err := newSolicitPacket(mac)
	if err != nil {
		return nil, fmt.Errorf("Request Error:\nnew solicit packet: %v\nerr: %v", solicitPacket, err)
	}

	if err = c.SendSolicitPacket(solicitPacket, mac); err != nil {
		return nil, fmt.Errorf("Request Error:\nsend solicit packet: %v\nerr: %v", solicitPacket, err)
	}

	advertisePacket, err := c.GetOffer()
	if err != nil {
		return nil, fmt.Errorf("Request Error:\nadvertise packet: %v\nerr: %v", advertisePacket, err)
	}
	return advertisePacket, nil
}

func (c *Client) SendSolicitPacket(p *dhcp6.Packet, mac *net.HardwareAddr) error {
	return c.connection.Write(p, mac)
}

func (c *Client) GetOffer() (*dhcp6.Packet, error) {
	var err error
	for i := 0; i < 5; i++ { // five attempts
		var pb []byte
		pb, err = c.connection.ReadFrom()
		if err != nil {
			continue
		}

		ipv6 := header.IPv6(pb[:header.IPv6MinimumSize])
		pb = pb[header.IPv6MinimumSize:]

		if ipv6.NextHeader() == uint8(header.UDPProtocolNumber) {
			udp := header.UDP(pb[:header.UDPMinimumSize])
			pb = pb[header.UDPMinimumSize:]

			if udp.DestinationPort() == srcPort {
				dhcp6p := &dhcp6.Packet{}
				if err = dhcp6p.UnmarshalBinary(pb); err != nil {
					continue
				}
				return dhcp6p, nil
			}
		}
	}
	return nil, fmt.Errorf("failed to get ipv6 address after five attempts: %v", err)
}
