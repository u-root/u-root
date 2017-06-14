package dhcp6client

import (
	"fmt"
	"net"
	"time"

	"github.com/mdlayher/dhcp6"
	// "golang.org/x/net/icmp"
	//"golang.org/x/net/ipv6"
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
	Write(packet []byte, mac net.HardwareAddr) error
	// WriteNeighborAd(src, dst net.IP, pb []byte) error
	ReadFrom() ([]byte, error)
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

func (c *Client) Request(mac *net.HardwareAddr) ([]byte, error) {
	solicitPacket, err := c.SendSolicitPacket(mac)
	if err != nil {
		return solicitPacket, err
	}
	y, z := c.GetOffer()
	fmt.Printf("get offer: %v, %v\n", y, z)
	// err = c.connection.Close()
	// if err != nil {
	// 	return false, solicitPacket, err
	// }
	return y, nil
}

func (c *Client) GetOffer() ([]byte, error) {
	pb, err := c.connection.ReadFrom()
	if err != nil {
		return nil, err
	}

	ipv6Hdr := unmarshalIPv6Hdr(pb[:40])

	if ipv6Hdr.NextHeader == 17 { // if next header is UDP 17
		udphdr := unmarshalUdpHdr(pb[40:48])
		if udphdr.Dst == 546 {
			var p dhcp6.Packet
			if err = p.UnmarshalBinary(pb[48:]); err != nil {
				return pb, err
			}
			fmt.Printf("dhcp packet: %v\n", p)
		}
	}

	return pb, nil
}

func (c *Client) SendSolicitPacket(mac *net.HardwareAddr) ([]byte, error) {
	options, err := addSolicitOptions(mac)
	if err != nil {
		return nil, err
	}

	pb, err := newDhcpPacket(dhcp6.MessageTypeSolicit, [3]byte{0, 1, 2}, nil, options)
	if err != nil {
		return nil, err
	}
	return pb, c.connection.Write(pb, *mac)
}

// func (c *Client) SendNeighborAdPacket(src, dst net.IP, icmpMsg *icmp.Message) ([]byte, error) {
// 	flags := []byte{0x40, 0x00, 0x00, 0x00}
// 	targetAddr := make([]byte, len(dst))
// 	copy(targetAddr[:], src)
// 	// m := icmp.Message {
// 	// 	Type: ipv6.ICMPTypeNeighborAdvertisement,
// 	// 	Code: 0,
// 	// 	Body: nil, // &icmp.DefaultMessageBody {
// 	// 		//Data: []byte{},//append(flags, targetAddr...),
// 	// 	//},
// 	// }
//
// 	// psh := icmp.IPv6PseudoHeader(src, dst)
//
// 	// var mtype int
// 	// switch typ := m.Type.(type) {
// 	// case ipv6.ICMPType:
// 	// 	mtype = int(typ)
// 	// default:
// 	// 	return nil, syscall.EINVAL
// 	// }
//
// 	// pb, err := m.Marshal(nil)
// 	// fmt.Printf("sending ad: %v, %v, %v, %v\n", mtype, psh, pb, err)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	pb := []byte{136, 0, 0x73, 0x26}
// 	pb = append(pb, append(flags, targetAddr...)...)
// 	return pb, c.connection.WriteNeighborAd(src, dst, pb)
// }
