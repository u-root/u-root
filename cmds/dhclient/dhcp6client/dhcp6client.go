package dhcp6client

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mdlayher/dhcp6"
	// "github.com/d2g/dhcp4"
)

type Client struct {
	hardwareAddr  net.HardwareAddr //The HardwareAddr to send in the request.
	ignoreServers []net.IP         //List of Servers to Ignore requests from.
	timeout       time.Duration    //Time before we timeout.
	broadcast     bool             //Set the Bcast flag in BOOTP Flags
	connection    connection       //The Connection Method to use
}

/*
*  * Abstracts the type of underlying socket used
*   */
type connection interface {
	Close() error
	Write(packet []byte) error
	ReadFrom()
	// SetReadTimeout(t time.Duration) error
}

func New(haddr net.HardwareAddr, conn connection, timeout time.Duration) (*Client, error) {
	c := Client{
		broadcast: true,
	}

	c.hardwareAddr = haddr
	c.connection = conn
	c.timeout = timeout
	return &c, nil
}

func NewPacket(messageType dhcp6.MessageType, txID [3]byte, addr *net.UDPAddr, options dhcp6.Options) []byte {
	packet := &dhcp6.Packet{
		MessageType:   messageType,
		TransactionID: txID,
		Options:       options,
	}

	fmt.Printf("solicitPacket: %v\n", packet)

	pb, err := packet.MarshalBinary()
	if err != nil {
		log.Printf("packet %v marshal to binary err: %v\n", txID, err)
		return nil
	}
	return pb
}

func (c *Client) Request(mac *net.HardwareAddr) (bool, []byte, error) {
	solicitPacket, err := c.SendSolicitPacket(mac)
	if err != nil {
		return false, solicitPacket, err
	}
	c.GetAdvertisePacket()
	err = c.connection.Close()
	if err != nil {
		return false, solicitPacket, err
	}
	return true, solicitPacket, nil
}

func (c *Client) SendSolicitPacket(mac *net.HardwareAddr) ([]byte, error) {
	// make options: iata
	var id = [4]byte{'r', 'o', 'o', 't'}
	options := make(dhcp6.Options)
	if err := options.Add(dhcp6.OptionIATA, dhcp6.NewIATA(id, nil)); err != nil {
		return nil, err
	}
	// make options: duid with mac address
	duid := dhcp6.NewDUIDLL(6, *mac)
	db, err := duid.MarshalBinary()
	if err != nil {
		return nil, err
	}
	addRaw(options, dhcp6.OptionClientID, db)

	pb := NewPacket(dhcp6.MessageTypeSolicit, [3]byte{0, 1, 2}, nil, options)
	return pb, c.connection.Write(pb)
}

func (c *Client) GetAdvertisePacket() {
	c.connection.ReadFrom()
}

func (c *Client) PrintConn() {
	fmt.Printf("print connection: %v\n", c.connection)
}
