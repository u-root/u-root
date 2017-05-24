package dhcp6client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/mdlayher/dhcp6"
	"github.com/vishvananda/netlink"
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

func (c *Client) Request() (bool, []byte, error) {
	solicitPacket, err := c.SendSolicitPacket()
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

func (c *Client) SendSolicitPacket() ([]byte, error) {
	var id = [4]byte{'r', 'o', 'o', 't'}
	options := make(dhcp6.Options)

	if err := options.Add(dhcp6.OptionIATA, dhcp6.NewIATA(id, nil)); err != nil {
		return nil, err
	}
	duid := dhcp6.NewDUIDLL(6, []byte{0xb8, 0xae, 0xed, 0x79, 0x61, 0x91})
	db, err := duid.MarshalBinary()
	if err != nil {
		return nil, err
	}
	addRaw(options, dhcp6.OptionClientID, db)

	solicitPacket := NewPacket(dhcp6.MessageTypeSolicit, [3]byte{0, 1, 2}, nil, options)
	padToMinSize(&solicitPacket)
	return solicitPacket, c.connection.Write(solicitPacket)
}

func (c *Client) GetAdvertisePacket() {
	c.connection.ReadFrom()
}

var padder [272]byte

func padToMinSize(p *[]byte) {
	if n := len(*p); n < 272 {
		*p = append(*p, padder[:272-n]...)
	}
}

func Dhclient6(ifname string) error {
	timeout := 15 * time.Second

	n, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/net/%s/address", ifname))
	if err != nil {
		return fmt.Errorf("cannot get mac for %v: %v\n", ifname, err)
	}

	// This is truly amazing but /sys appends newlines to all this data.
	n = bytes.TrimSpace(n)
	mac, err := net.ParseMAC(string(n))
	if err != nil {
		return fmt.Errorf("mac error: %v\n", err)
	}

	iface, err := netlink.LinkByName(ifname)
	if err != nil {
		return fmt.Errorf("%s: netlink.LinkByName failed: %v\n", ifname, err)
	}

	if err := netlink.LinkSetUp(iface); err != nil {
		return fmt.Errorf("%v: %v can't make it up: %v\n", ifname, iface, err)
	}

	conn, err := NewPacketSock(iface.Attrs().Index)
	if err != nil {
		return fmt.Errorf("client conection generation: %v\n", err)
	}

	client, err := New(mac, conn, timeout)
	if err != nil {
		return fmt.Errorf("error: %v\n", err)
	}

	success, packet, err := client.Request()
	if err != nil {
		return fmt.Errorf("result: %v, %v, %v\n", success, packet, err)
	}
	fmt.Printf("result: %v, %v, %v\n", success, packet, err)
	return nil
}
