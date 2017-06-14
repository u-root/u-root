package dhcp6client

import (
	"encoding/binary"
	"math/rand"
	"net"

	"github.com/mdlayher/dhcp6"
	"golang.org/x/net/ipv6"
	"golang.org/x/sys/unix"
)

const (
	ipv6HdrLen = 40
	udpHdrLen  = 8

	srcPort = 546
	dstPort = 547
)

type packetSock struct {
	fd      int
	ifindex int
}

var bcastMAC = []byte{0x33, 0x33, 0x00, 0x01, 0x00, 0x02}

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

// Write dhcpv6 requests
func (pc *packetSock) Write(p *dhcp6.Packet, mac net.HardwareAddr) error {
	// Define linke layer
	lladdr := unix.SockaddrLinklayer{
		Ifindex:  pc.ifindex,
		Protocol: swap16(unix.ETH_P_IPV6),
		Halen:    uint8(len(bcastMAC)),
	}
	copy(lladdr.Addr[:], bcastMAC)

	flowLabel := rand.Int() & 0xfffff
	// src := mac2ipv6(mac)

	pb, err := p.MarshalBinary()
	if err != nil {
		return err
	}

	h1 := &ipv6.Header{
		Version:      ipv6.Version,
		TrafficClass: 0,
		FlowLabel:    flowLabel,
		PayloadLen:   udpHdrLen + len(pb),
		NextHeader:   unix.IPPROTO_UDP,
		HopLimit:     1,
		// TODO: get rid of hard-coded addr
		Src: net.ParseIP("fe80::baae:edff:fe79:6191"),
		// wide-dhcpv6-client addr: net.ParseIP("fe80::179a:1422:c923:2727"),
		Dst: net.ParseIP("FF02::1:2"),
	}

	h2 := &Udphdr{
		Src:    srcPort,
		Dst:    dstPort,
		Length: uint16(udpHdrLen + len(pb)),
	}

	pkt, err := marshalPacket(h1, h2, pb)
	if err != nil {
		return err
	}

	// Send out request from link layer
	return unix.Sendto(pc.fd, pkt, 0, &lladdr)
}

// Write icmpv6 neighbor advertisements
//func (pc *packetSock) WriteNeighborAd(src, dst net.IP, pb []byte) error {
//	mac := ipv62mac([]byte(dst))
//	fmt.Printf("addr: %v, %x \n", []byte(dst), mac)
//	// Define linke layer
//	lladdr := unix.SockaddrLinklayer{
//		Ifindex:  pc.ifindex,
//		Protocol: swap16(unix.ETH_P_IPV6),
//		Halen:    uint8(len(mac)),
//	}
//	copy(lladdr.Addr[:], mac)
//
//	flowLabel := rand.Int() & 0xfffff
//
//	h := &ipv6.Header{
//		Version:      ipv6.Version,
//		TrafficClass: 0,
//		FlowLabel:    flowLabel,
//		PayloadLen:   24,
//		NextHeader:   unix.IPPROTO_ICMPV6,
//		HopLimit:     255,
//		// TODO: src ip harded coded for now
//		Src: src,
//		Dst: dst,
//	}
//
//	pkt := make([]byte, ipv6HdrLen+len(pb))
//	ipv6hdr := &marshalIPv6Hdr(h)
//
//	// Wrap up packet
//	copy(pkt[0:ipv6HdrLen], ipv6hdr)
//	copy(pkt[ipv6HdrLen:], pb)
//
//	// Send out request from link layer
//	return unix.Sendto(pc.fd, pkt, 0, &lladdr)
//}

func (pc *packetSock) ReadFrom() ([]byte, error) {
	pb := make([]byte, 1500)
	n, _, err := unix.Recvfrom(pc.fd, pb, 0)
	if err != nil {
		return nil, err
	}
	return pb[:n], nil
}

func (pc *packetSock) Close() error {
	return unix.Close(pc.fd)
}

func swap16(x uint16) uint16 {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], x)
	return binary.LittleEndian.Uint16(b[:])
}
