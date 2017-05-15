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
	var p dhcp6.Packet
	for i := 0; i < 5; i++ { // five attempts
		pb, err := c.connection.ReadFrom()
		if err != nil {
			continue
		}

		ipv6Hdr := unmarshalIPv6Hdr(pb[:40])

		if ipv6Hdr.NextHeader == protocolUDP { // if next header is UDP
			udphdr := unmarshalUdpHdr(pb[40:48])
			if udphdr.Dst == srcPort {
				if err = p.UnmarshalBinary(pb[48:]); err != nil {
					continue
				}
				return &p, nil
			}
		}
	}
	return nil, fmt.Errorf("failed to get ipv6 address after five attempts: %v", err)
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
