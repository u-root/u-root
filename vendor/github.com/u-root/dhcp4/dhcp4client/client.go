// Package dhcp4client is a small, minimum-functionality client for DHCPv4.
//
// It only supports the 4-way DHCPv4 Discover-Offer-Request-Ack handshake as
// well as the Request-Ack renewal process.
package dhcp4client

import (
	"math/rand"
	"net"
	"time"

	"github.com/u-root/dhcp4"
	"github.com/u-root/dhcp4/dhcp4opts"
)

const (
	maxMessageSize = 1500

	// ClientPort is the port that DHCP clients listen on.
	ClientPort = 68

	// ServerPort is the port that DHCP servers and relay agents listen on.
	ServerPort = 67
)

var (
	// AllDHCPServers is the address of all link-local DHCP servers and
	// relay agents.
	AllDHCPServers = &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: ServerPort,
	}
)

// Client is a simple IPv4 DHCP client.
type Client struct {
	hardwareAddr net.HardwareAddr
	conn         net.PacketConn
	timeout      time.Duration
}

// New creates a new DHCP client that sends and receives packets on the given
// interface.
func New(link *net.Interface) (*Client, error) {
	conn, err := NewIPv4UDPConn(link.Name, ClientPort)
	if err != nil {
		return nil, err
	}

	return &Client{
		hardwareAddr: link.HardwareAddr,
		timeout:      time.Second * 10,
		conn:         conn,
	}, nil
}

// Request completes the 4-way Discover-Offer-Request-Ack handshake.
func (c *Client) Request() (*dhcp4.Packet, error) {
	discover := c.discoverPacket()
	offer, err := c.SendAndRead(discover)
	if err != nil {
		return nil, err
	}

	return c.SendAndRead(c.requestPacket(offer))
}

// Renew sends a renewal request packet and waits for the corresponding response.
func (c *Client) Renew(ack *dhcp4.Packet) (*dhcp4.Packet, error) {
	return c.SendAndRead(c.requestPacket(ack))
}

// Close closes the client connection.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// SendAndRead sends a packet and waits for the corresponding response packet;
// matching the transaction ID.
//
// TODO: Make this a little state machine of packet types. See RFC 2131,
// Section 4.4, Figure 5.
func (c *Client) SendAndRead(packet *dhcp4.Packet) (*dhcp4.Packet, error) {
	if err := c.SendPacket(packet); err != nil {
		return nil, err
	}

	for {
		// TODO: We should have an overall timeout for the SendAndRead
		// operation, not one for just the socket read.
		c.conn.SetReadDeadline(time.Now().Add(c.timeout))
		p := make([]byte, maxMessageSize)
		n, _, err := c.conn.ReadFrom(p)
		if err != nil {
			return nil, err
		}

		resp := &dhcp4.Packet{}
		if err := resp.UnmarshalBinary(p[:n]); err != nil {
			return nil, err
		}

		if packet.TransactionID != resp.TransactionID {
			continue
		}
		return resp, nil
	}

}

// SendPacket broadcasts a DHCPv4 packet.
func (c *Client) SendPacket(packet *dhcp4.Packet) error {
	p, err := packet.MarshalBinary()
	if err != nil {
		return err
	}
	_, err = c.conn.WriteTo(p, AllDHCPServers)
	return err
}

func (c *Client) discoverPacket() *dhcp4.Packet {
	packet := dhcp4.NewPacket(dhcp4.BootRequest)
	rand.Read(packet.TransactionID[:])
	packet.CHAddr = c.hardwareAddr
	packet.Broadcast = true

	packet.Options.Add(dhcp4.OptionDHCPMessageType, dhcp4opts.DHCPDiscover)
	packet.Options.Add(dhcp4.OptionMaximumDHCPMessageSize, dhcp4opts.Uint16(maxMessageSize))
	return packet
}

func (c *Client) requestPacket(reply *dhcp4.Packet) *dhcp4.Packet {
	packet := dhcp4.NewPacket(dhcp4.BootRequest)

	packet.CHAddr = c.hardwareAddr
	packet.TransactionID = reply.TransactionID
	packet.CIAddr = reply.CIAddr
	packet.SIAddr = reply.SIAddr
	packet.Broadcast = true

	packet.Options.Add(dhcp4.OptionDHCPMessageType, dhcp4opts.DHCPRequest)
	packet.Options.Add(dhcp4.OptionMaximumDHCPMessageSize, dhcp4opts.Uint16(maxMessageSize))
	// Request the offered IP address.
	packet.Options.Add(dhcp4.OptionRequestedIPAddress, dhcp4opts.IP(reply.YIAddr))

	sid, err := dhcp4opts.GetServerIdentifier(reply.Options)
	if err == nil {
		packet.Options.Add(dhcp4.OptionServerIdentifier, dhcp4opts.IP(sid))
	}
	return packet
}
