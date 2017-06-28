package dhcp6client

import (
	"encoding/binary"
	"math/rand"
	"net"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/header"
	"github.com/mdlayher/dhcp6"
	"golang.org/x/sys/unix"
)

type packetSock struct {
	fd      int
	ifindex int
}

// All DHCP servers and relay agents on the local network segment (RFC 3315)
// IPv6 Multicast (RFC 2464)
// insert the low 32 Bits of the multicast IPv6 Address into the Ethernet Address (RFC 7042 2.3.1.)
var multicastMAC = []byte{0x33, 0x33, 0x00, 0x01, 0x00, 0x02}

/*
* connection
 */
type connection interface {
	Close() error
	Write(p *dhcp6.Packet, mac *net.HardwareAddr) error
	// WriteNeighborAd(src, dst net.IP, pb []byte) error
	ReadFrom() ([]byte, error)
	// SetReadTimeout(t time.Duration) error
}

/*
* PacketSock that implements connection
 */
// NewPackSock creates a new socket that sends and receives packets.
func NewPacketSock(ifindex int) (*packetSock, error) {
	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_DGRAM, int(swap16(unix.ETH_P_IPV6)))
	if err != nil {
		return nil, err
	}
	addr := unix.SockaddrLinklayer{
		Ifindex:  ifindex,
		Protocol: swap16(unix.ETH_P_IPV6),
	}
	if err = unix.Bind(fd, &addr); err != nil {
		return nil, err
	}
	return &packetSock{
		fd:      fd,
		ifindex: ifindex,
	}, nil
}

// Write dhcpv6 requests.
func (pc *packetSock) Write(p *dhcp6.Packet, mac *net.HardwareAddr) error {
	// Define linke layer
	lladdr := unix.SockaddrLinklayer{
		Ifindex:  pc.ifindex,
		Protocol: swap16(unix.ETH_P_IPV6),
		Halen:    uint8(len(multicastMAC)),
	}
	copy(lladdr.Addr[:], multicastMAC)

	flowLabel := rand.Uint32()

	pb, err := p.MarshalBinary()
	if err != nil {
		return err
	}

	length := uint16(header.UDPMinimumSize + len(pb))
	ipv6fields := &header.IPv6Fields{
		FlowLabel:     flowLabel,
		PayloadLength: length,
		NextHeader:    uint8(header.UDPProtocolNumber),
		HopLimit:      1,
		SrcAddr:       tcpip.Address(mac2ipv6(mac)),
		DstAddr:       tcpip.Address(net.ParseIP("FF02::1:2").To16()),
	}
	ipv6header := header.IPv6(make([]byte, header.IPv6MinimumSize))
	ipv6header.Encode(ipv6fields)

	udphdr := header.UDP(make([]byte, header.UDPMinimumSize))
	udphdr.Encode(&header.UDPFields{
		SrcPort: srcPort,
		DstPort: dstPort,
		Length:  length,
	})

	xsum := header.Checksum(pb, header.PseudoHeaderChecksum(ipv6header.TransportProtocol(), ipv6fields.SrcAddr, ipv6fields.DstAddr))
	udphdr.SetChecksum(^udphdr.CalculateChecksum(xsum, length))

	pkt := append([]byte(ipv6header), []byte(udphdr)...)
	pkt = append(pkt, pb...)

	// Send out request from link layer
	return unix.Sendto(pc.fd, pkt, 0, &lladdr)
}

// Read from replies from server.
func (pc *packetSock) ReadFrom() ([]byte, error) {
	pb := make([]byte, 1500)
	n, _, err := unix.Recvfrom(pc.fd, pb, 0)
	if err != nil {
		return nil, err
	}
	return pb[:n], nil
}

// Close socket.
func (pc *packetSock) Close() error {
	return unix.Close(pc.fd)
}

/**
* Utils
 */
func swap16(x uint16) uint16 {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], x)
	return binary.LittleEndian.Uint16(b[:])
}

func mac2ipv6(mac *net.HardwareAddr) []byte {
	v6addr := make([]byte, 6)
	copy(v6addr, *mac)
	v6addr = append(v6addr[:3], append([]byte{0xff, 0xfe}, v6addr[3:]...)...)
	if v6addr[0]&0x02 == 0x02 {
		v6addr[0] -= 0x02
	} else {
		v6addr[0] += 0x02
	}
	v6addr = append([]byte{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, v6addr...)
	return v6addr
}
