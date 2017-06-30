package dhcp6client

import (
	"math/rand"
	"net"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/header"
	"github.com/mdlayher/dhcp6"
)

// Should send a PR to dhcp6 to add these and others.
const (
	OptionDefaultIPTTL dhcp6.OptionCode = 23
	OptionMTUTimeout   dhcp6.OptionCode = 24
)

func newSolicitOptions(mac net.HardwareAddr) (dhcp6.Options, error) {
	// make options: iana
	options := make(dhcp6.Options)

	var id = [4]byte{'r', 'o', 'o', 't'}
	if err := options.Add(dhcp6.OptionIANA, dhcp6.NewIANA(id, 0, 0, nil)); err != nil {
		return nil, err
	}
	if err := options.Add(dhcp6.OptionRapidCommit, nil); err != nil {
		return nil, err
	}
	if err := options.Add(dhcp6.OptionElapsedTime, dhcp6.ElapsedTime(0)); err != nil {
		return nil, err
	}

	oro := dhcp6.OptionRequestOption{OptionDefaultIPTTL, OptionMTUTimeout}
	if err := options.Add(dhcp6.OptionORO, oro); err != nil {
		return nil, err
	}

	// HardwareType *should* be 1?
	if err := options.Add(dhcp6.OptionClientID, dhcp6.NewDUIDLL(6, mac)); err != nil {
		return nil, err
	}

	return options, nil
}

func newSolicitPacket(mac net.HardwareAddr) (*dhcp6.Packet, error) {
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

// mac2ipv6 gives the EUI-64 IPv6 address corresponding to the mac.
func mac2ipv6(mac net.HardwareAddr) []byte {
	var v6addr []byte
	v6addr = append(mac[:3], append([]byte{0xff, 0xfe}, mac[3:]...)...)

	// Invert 7th bit from left.
	if v6addr[0]&0x02 == 0x02 {
		v6addr[0] -= 0x02
	} else {
		v6addr[0] += 0x02
	}

	return append([]byte{0xfe, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, v6addr...)
}

// ipv6UDPPacket wraps a dhcp6 packet in a IPv6 and UDP packet for sending on a
// packet socket.
func ipv6UDPPacket(p *dhcp6.Packet, srcMAC net.HardwareAddr) ([]byte, error) {
	pb, err := p.MarshalBinary()
	if err != nil {
		return nil, err
	}

	length := uint16(header.UDPMinimumSize + len(pb))
	ipv6fields := &header.IPv6Fields{
		FlowLabel:     rand.Uint32(),
		PayloadLength: length,
		NextHeader:    uint8(header.UDPProtocolNumber),
		HopLimit:      1,
		SrcAddr:       tcpip.Address(mac2ipv6(srcMAC)),
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
	return append(pkt, pb...), nil
}
