package bsdp

import (
	"errors"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/client4"
)

// Client represents a BSDP client that can perform BSDP exchanges via the
// broadcast address.
type Client struct {
	client4.Client
}

// NewClient constructs a new client with default read and write timeouts from
// dhcpv4.Client.
func NewClient() *Client {
	return &Client{Client: client4.Client{}}
}

// Exchange runs a full BSDP exchange (Inform[list], Ack, Inform[select],
// Ack). Returns a list of DHCPv4 structures representing the exchange.
func (c *Client) Exchange(ifname string) ([]*Packet, error) {
	conversation := make([]*Packet, 0)

	// Get our file descriptor for the broadcast socket.
	sendFd, err := client4.MakeBroadcastSocket(ifname)
	if err != nil {
		return conversation, err
	}
	recvFd, err := client4.MakeListeningSocket(ifname)
	if err != nil {
		return conversation, err
	}

	// INFORM[LIST]
	informList, err := NewInformListForInterface(ifname, dhcpv4.ClientPort)
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, informList)

	// ACK[LIST]
	ackForList, err := c.Client.SendReceive(sendFd, recvFd, informList.v4(), dhcpv4.MessageTypeAck)
	if err != nil {
		return conversation, err
	}

	// Rewrite vendor-specific option for pretty printing.
	conversation = append(conversation, PacketFor(ackForList))

	// Parse boot images sent back by server
	bootImages, err := ParseBootImageListFromAck(ackForList)
	if err != nil {
		return conversation, err
	}
	if len(bootImages) == 0 {
		return conversation, errors.New("got no BootImages from server")
	}

	// INFORM[SELECT]
	informSelect, err := InformSelectForAck(PacketFor(ackForList), dhcpv4.ClientPort, bootImages[0])
	if err != nil {
		return conversation, err
	}
	conversation = append(conversation, informSelect)

	// ACK[SELECT]
	ackForSelect, err := c.Client.SendReceive(sendFd, recvFd, informSelect.v4(), dhcpv4.MessageTypeAck)
	if err != nil {
		return conversation, err
	}
	return append(conversation, PacketFor(ackForSelect)), nil
}
