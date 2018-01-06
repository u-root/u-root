// Package client is a small, minimum-functionality client for DHCPv4.
//
// It only supports the 4-way DHCPv4 Discover-Offer-Request-Ack handshake as
// well as the Request-Ack renewal process.
package client

import (
	"bytes"
	"math/rand"
	"net"
	"time"

	"github.com/u-root/dhcp4"
	"github.com/u-root/dhcp4/opts"
)

const (
	maxMessageSize = 1500
)

// Client is a simple DHCPv4 client.
type Client struct {
	hardwareAddr net.HardwareAddr
	timeout      time.Duration
	conn         *packetSock
}

// New creates a new DHCPv4 client that sends and receives packets on the given
// interface.
func New(link *net.Interface) (*Client, error) {
	conn, err := newPacketSock(link.Index)
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
		c.conn.SetReadTimeout(c.timeout)
		p := make([]byte, maxMessageSize)
		n, _, err := c.conn.Read(p)
		if err != nil {
			return nil, err
		}

		resp := new(dhcp4.Packet)
		if err := resp.UnmarshalBinary(p[:n]); err != nil {
			return nil, err
		}

		if !bytes.Equal(packet.TransactionID[:], resp.TransactionID[:]) {
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
	return c.conn.Write(p)
}

func (c *Client) discoverPacket() *dhcp4.Packet {
	packet := dhcp4.NewPacket(dhcp4.BootRequest)
	rand.Read(packet.TransactionID[:])
	packet.CHAddr = c.hardwareAddr
	packet.Broadcast = true

	packet.Options.Add(dhcp4.OptionDHCPMessageType, opts.DHCPDiscover)
	packet.Options.Add(dhcp4.OptionMaximumDHCPMessageSize, opts.Uint16(maxMessageSize))
	return packet
}

func (c *Client) requestPacket(reply *dhcp4.Packet) *dhcp4.Packet {
	packet := dhcp4.NewPacket(dhcp4.BootRequest)

	packet.CHAddr = c.hardwareAddr
	packet.TransactionID = reply.TransactionID
	packet.CIAddr = reply.CIAddr
	packet.SIAddr = reply.SIAddr
	packet.Broadcast = true

	packet.Options.Add(dhcp4.OptionDHCPMessageType, opts.DHCPRequest)
	packet.Options.Add(dhcp4.OptionMaximumDHCPMessageSize, opts.Uint16(maxMessageSize))
	// Request the offered IP address.
	packet.Options.Add(dhcp4.OptionRequestedIPAddress, opts.IP(reply.YIAddr))

	sid, err := opts.GetServerIdentifier(reply.Options)
	if err == nil {
		packet.Options.Add(dhcp4.OptionServerIdentifier, opts.IP(sid))
	}
	return packet
}
