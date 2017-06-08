package dhcp6client

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"

	// "github.com/mdlayher/dhcp6"
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
func (pc *packetSock) Write(pb []byte) error {
	// Define linke layer
	lladdr := unix.SockaddrLinklayer{
		Ifindex:  pc.ifindex,
		Protocol: swap16(unix.ETH_P_IPV6),
		Halen:    uint8(len(bcastMAC)),
	}
	copy(lladdr.Addr[:], bcastMAC)

	flowLabel := rand.Int() & 0xfffff

	h := ipv6.Header{
		Version:      ipv6.Version,
		TrafficClass: 0,
		FlowLabel:    flowLabel,
		PayloadLen:   udpHdrLen + len(pb),
		NextHeader:   unix.IPPROTO_UDP,
		HopLimit:     1,
		// TODO: src ip harded coded for now
		Src: net.ParseIP("fe80::baae:edff:fe79:6191"),
		Dst: net.ParseIP("FF02::1:2"),
	}

	pkt := make([]byte, ipv6HdrLen+udpHdrLen+len(pb))
	ipv6hdr := marshalIPv6Hdr(h)
	udphdr := fillUDPHdr(len(pb))

	// Wrap up packet
	copy(pkt[0:ipv6HdrLen], ipv6hdr)
	copy(pkt[ipv6HdrLen:ipv6HdrLen+udpHdrLen], udphdr)
	copy(pkt[ipv6HdrLen+udpHdrLen:len(pkt)], pb)

	udpChksum(ipv6hdr, udphdr, pb, pkt[ipv6HdrLen+6:ipv6HdrLen+8])

	// Send out request from link layer
	return unix.Sendto(pc.fd, pkt, 0, &lladdr)
}

func (pc *packetSock) ReadFrom() ([]byte, error) {
	fmt.Printf("starts reading\n")
	pb := make([]byte, 500)
	n, _, err := unix.Recvfrom(pc.fd, pb, 0)
	if err != nil {
		return nil, err
	}
	return pb[:n], nil
	//packet := dhcp6.Packet{}
	//UnmarshalBinary(&packet, pb)

	//ipv6hdr := marshalIPv6Hdr(pb)
	//fmt.Printf("ip header: %v\n", ipv6hdr)
	//fmt.Printf("read from server: %v, %v, %v\n", n, pb, err)
}

func (pc *packetSock) Close() error {
	return unix.Close(pc.fd)
}

func swap16(x uint16) uint16 {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], x)
	return binary.LittleEndian.Uint16(b[:])
}

func unmarshalIPv6Hdr(hb []byte) ipv6.Header {
	return ipv6.Header{
		Version:      ipv6.Version,
		TrafficClass: int((hb[0]&0x0f)<<4 | (hb[1]&0xf0)>>4),
		FlowLabel:    int((hb[1]&0x0f)<<16 | hb[2]<<8 | hb[3]),
		PayloadLen:   int(hb[4]<<8 | hb[5]),
		NextHeader:   int(hb[6]),
		HopLimit:     int(hb[7]),
		Src:          hb[8:24],
		Dst:          hb[24:40],
	}
}

func marshalIPv6Hdr(h ipv6.Header) []byte {
	ipv6hdr := make([]byte, ipv6HdrLen)
	// ver + first half byte of traffic class
	ipv6hdr[0] = byte(h.Version<<4 | (h.TrafficClass >> 4))
	// second half byte of traffic class + first half byte of flow label
	ipv6hdr[1] = byte(((h.TrafficClass & 0x0f) << 4) | (h.FlowLabel >> 16))
	// flow label
	ipv6hdr[2] = byte(h.FlowLabel & 0x0ff00 >> 8)
	ipv6hdr[3] = byte(h.FlowLabel & 0x000ff)
	// payload length
	binary.BigEndian.PutUint16(ipv6hdr[4:6], uint16(h.PayloadLen))
	// next header
	ipv6hdr[6] = byte(h.NextHeader)
	// hop limit
	ipv6hdr[7] = byte(h.HopLimit)
	// src
	copy(ipv6hdr[8:24], h.Src)
	// dst
	copy(ipv6hdr[24:40], h.Dst)

	return ipv6hdr
}

func fillUDPHdr(payloadLen int) []byte {
	udphdr := make([]byte, udpHdrLen)
	//src port
	binary.BigEndian.PutUint16(udphdr[0:2], srcPort)
	// dest port
	binary.BigEndian.PutUint16(udphdr[2:4], dstPort)
	// length
	binary.BigEndian.PutUint16(udphdr[4:6], uint16(udpHdrLen+payloadLen))
	return udphdr
}

func udpChksum(ipv6hdr []byte, udphdr []byte, udpbody []byte, csum []byte) {
	// psuedoheader = srcip + dstip + udpcode + udplen
	psuedoHdr := append(ipv6hdr[8:], []byte{0x00, 0x11 /*udpcode=17*/}...)
	psuedoHdr = append(psuedoHdr, udphdr[4:6]...)

	// udp header + data (excluding checksum)
	udpData := append(udphdr, udpbody...)

	// calculate 16-bit sum
	sumPsuedoHdr := sixteenBitSum(psuedoHdr)
	sumUdpData := sixteenBitSum(udpData)
	sumTotal := sumPsuedoHdr + sumUdpData
	sumTotal = (sumTotal>>16 + sumTotal&0xffff)
	sumTotal = sumTotal + (sumTotal >> 16)

	// one's complement
	sumTotal = ^sumTotal

	csum[0] = uint8(sumTotal & 0xff)
	csum[1] = uint8(sumTotal >> 8)
}

func sixteenBitSum(p []byte) uint32 {
	cklen := len(p)
	s := uint32(0)
	for i := 0; i < (cklen - 1); i += 2 {
		s += uint32(p[i+1])<<8 | uint32(p[i])
	}
	if cklen&1 == 1 {
		s += uint32(p[cklen-1])
	}
	s = (s >> 16) + (s & 0xffff)
	s = s + (s >> 16)

	return s
}

//func MarshalBinary(req *dhcp6.Request) ([]byte, error) {
//	r := *req
//	opts := enumerate(r.Options)
//	addrbyte := []byte(r.RemoteAddr)
//	b := make([]byte, 6+opts.count()+len(addrbyte))
//	b[0] = byte(r.MessageType)
//	copy(b[1:4], r.TransactionID[:])
//	opts.write(b[4 : 4+opts.count()])
//	copy(b[4+opts.count():], addrbyte[:])
//
//	return b, nil
//}

//func UnmarshalBinary(p *dhcp6.Packet, b []byte) error {
//	// Packet must contain at least a message type and transaction ID
//	if len(b) < 4 {
//		return dhcp6.ErrInvalidPacket
//	}
//	p.MessageType = dhcp6.MessageType(b[0])
//	txID := [3]byte{}
//	copy(txID[:], b[1:4])
//	p.TransactionID = txID
//
//	options, err := parseOptions(b[4:])
//	if err != nil {
//		// Invalid options means an invalid packet
//		return dhcp6.ErrInvalidPacket
//	}
//	p.Options = options
//	return nil
//}

//func parseOptions(b []byte) (dhcp6.Options, error) {
//	var length int
//	options := make(dhcp6.Options)
//	buf := bytes.NewBuffer(b)
//	for buf.Len() > 3 {
//		// 2 bytes: option code
//		o := option{}
//		code := dhcp6.OptionCode(binary.BigEndian.Uint16(buf.Next(2)))
//		// If code is 0, bytes are empty after this point
//		if code == 0 {
//			return options, nil
//		}
//
//		o.Code = code
//		// 2 bytes: option length
//		length = int(binary.BigEndian.Uint16(buf.Next(2)))
//
//		// If length indicated is zero, skip to next iteration
//		if length == 0 {
//			continue
//		}
//
//		// N bytes: option data
//		o.Data = buf.Next(length)
//		// Set slice's max for option's data
//		o.Data = o.Data[:len(o.Data):len(o.Data)]
//
//		// If option data has less bytes than indicated by length,
//		// return an error
//		if len(o.Data) < length {
//			return nil, errors.New("invalid options data")
//		}
//
//		addRaw(options, o.Code, o.Data)
//	}
//	// Report error for any trailing bytes
//	if buf.Len() != 0 {
//		return nil, errors.New("invalid options data")
//	}
//	return options, nil
//}
