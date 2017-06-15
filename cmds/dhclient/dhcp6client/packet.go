package dhcp6client

import (
	"encoding/binary"
	"net"

	"github.com/mdlayher/dhcp6"
	"golang.org/x/net/ipv6"
)

type Udphdr struct {
	Src    uint16
	Dst    uint16
	Length uint16
	Csum   uint16
}

func marshalPacket(h1 *ipv6.Header, h2 *Udphdr, pb []byte) ([]byte, error) {
	ipv6hdr := marshalIPv6Hdr(h1)
	udphdr := marshalUdpHdr(h2)
	// add checksum to udp header
	csum := doCsum(ipv6hdr, udphdr, pb)
	copy(udphdr[6:8], csum)
	h2.Csum = binary.BigEndian.Uint16(csum)

	pkt := make([]byte, ipv6HdrLen+udpHdrLen+len(pb))

	copy(pkt[0:ipv6HdrLen], ipv6hdr)
	copy(pkt[ipv6HdrLen:ipv6HdrLen+udpHdrLen], udphdr)
	copy(pkt[ipv6HdrLen+udpHdrLen:len(pkt)], pb)
	return pkt, nil
}

/***
* IPV6 header
 */
func marshalIPv6Hdr(h *ipv6.Header) []byte {
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

func unmarshalIPv6Hdr(hb []byte) *ipv6.Header {
	flowlbl := binary.BigEndian.Uint32([]byte{0x00, hb[1] & 0x0f, hb[2], hb[3]})
	return &ipv6.Header{
		Version:      ipv6.Version,
		TrafficClass: int((hb[0]&0x0f)<<4 | (hb[1]&0xf0)>>4),
		FlowLabel:    int(flowlbl),
		PayloadLen:   int(hb[4]<<8 | hb[5]),
		NextHeader:   int(hb[6]),
		HopLimit:     int(hb[7]),
		Src:          hb[8:24],
		Dst:          hb[24:40],
	}
}

/***
* UDP header
 */
func marshalUdpHdr(h *Udphdr) []byte {
	udphdr := make([]byte, udpHdrLen)

	// src port
	binary.BigEndian.PutUint16(udphdr[0:2], h.Src)
	// dest port
	binary.BigEndian.PutUint16(udphdr[2:4], h.Dst)
	// length
	binary.BigEndian.PutUint16(udphdr[4:6], h.Length)
	return udphdr
}

func unmarshalUdpHdr(hb []byte) *Udphdr {
	return &Udphdr{
		Src:    binary.BigEndian.Uint16(hb[0:2]),
		Dst:    binary.BigEndian.Uint16(hb[2:4]),
		Length: binary.BigEndian.Uint16(hb[4:6]),
		Csum:   binary.BigEndian.Uint16(hb[6:8]),
	}
}

/***
* DHCP packet
 */
func newSolicitOptions(mac *net.HardwareAddr) (dhcp6.Options, error) {
	// make options: iata
	var id = [4]byte{0x00, 0x00, 0x00, 0x0f}
	options := make(dhcp6.Options)
	if err := options.Add(dhcp6.OptionIANA, dhcp6.NewIANA(id, 0, 0, nil)); err != nil {
		return nil, err
	}
	// make options: rapid commit
	if err := options.Add(dhcp6.OptionRapidCommit, nil); err != nil {
		return nil, err
	}
	// make options: elapsed time
	var et dhcp6.ElapsedTime
	et.UnmarshalBinary([]byte{0x00, 0x00})
	if err := options.Add(dhcp6.OptionElapsedTime, et); err != nil {
		return nil, err
	}
	// make options: option request option
	oro := make(dhcp6.OptionRequestOption, 4)
	oro.UnmarshalBinary([]byte{0x00, 0x17, 0x00, 0x18})
	if err := options.Add(dhcp6.OptionORO, oro); err != nil {
		return nil, err
	}
	// make options: duid with mac address
	duid := dhcp6.NewDUIDLL(6, *mac)
	db, err := duid.MarshalBinary()
	if err != nil {
		return nil, err
	}
	// add row
	options[dhcp6.OptionClientID] = append(options[dhcp6.OptionClientID], db)

	return options, nil
}

func newSolicitPacket(mac *net.HardwareAddr) (*dhcp6.Packet, error) {
	options, err := newSolicitOptions(mac)
	if err != nil {
		return nil, err
	}

	return &dhcp6.Packet{
		MessageType:   dhcp6.MessageTypeSolicit,
		TransactionID: [3]byte{0x00, 0x01, 0x02},
		Options:       options,
	}, nil
}

/***
* Helper functions
 */
func doCsum(ipv6hdr []byte, udphdr []byte, packet []byte) []byte {
	csum := make([]byte, 2)
	// psuedoheader = srcip + dstip + udpcode + udplen
	psuedoHdr := append(ipv6hdr[8:], []byte{0x00, 0x11 /*udpcode=17*/}...)
	psuedoHdr = append(psuedoHdr, udphdr[4:6]...)
	// udp header + data (excluding checksum)
	udpData := append(udphdr, packet...)
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
	return csum
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
