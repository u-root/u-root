package client

import (
	"net"
	"time"

	"github.com/google/netstack/tcpip"
	"github.com/google/netstack/tcpip/header"
	"golang.org/x/sys/unix"
)

const (
	srcPort = 68
	dstPort = 67
)

var (
	bcastMAC = []byte{255, 255, 255, 255, 255, 255}
)

type packetSock struct {
	fd      int
	ifindex int
}

func newPacketSock(ifindex int) (*packetSock, error) {
	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_DGRAM, int(swap16(unix.ETH_P_IP)))
	if err != nil {
		return nil, err
	}

	addr := unix.SockaddrLinklayer{
		Ifindex:  ifindex,
		Protocol: swap16(unix.ETH_P_IP),
	}
	if err = unix.Bind(fd, &addr); err != nil {
		return nil, err
	}

	return &packetSock{
		fd:      fd,
		ifindex: ifindex,
	}, nil
}

func (pc *packetSock) Close() error {
	return unix.Close(pc.fd)
}

func (pc *packetSock) Write(packet []byte) error {
	lladdr := unix.SockaddrLinklayer{
		Ifindex:  pc.ifindex,
		Protocol: swap16(unix.ETH_P_IP),
		Halen:    uint8(len(bcastMAC)),
	}
	copy(lladdr.Addr[:], bcastMAC)

	return unix.Sendto(pc.fd, udp4pkt(packet), 0, &lladdr)
}

func (pc *packetSock) Read(p []byte) (int, net.IP, error) {
	ipLen := header.IPv4MaximumHeaderSize
	udpLen := header.UDPMinimumSize

	for {
		pkt := make([]byte, ipLen+udpLen+len(p))
		n, _, err := unix.Recvfrom(pc.fd, pkt, 0)
		if err != nil {
			return 0, nil, err
		}
		pkt = pkt[:n]
		buf := &buffer{pkt}

		// To read the header length, access data directly.
		ipHdr := header.IPv4(buf.data)
		ipHdr = header.IPv4(buf.next(int(ipHdr.HeaderLength())))

		udpHdr := header.UDP(buf.next(udpLen))

		if udpHdr.DestinationPort() != srcPort {
			// Not for the port we're looking for.
			continue
		}

		// TODO: This is ugly for now. For one, we should check the
		// length in the UDP header to see if we got the full packet
		// (and read more if not).

		return copy(p, buf.remaining()), net.IP(ipHdr.SourceAddress()), nil
	}
}

type buffer struct {
	data []byte
}

func (b *buffer) next(n int) []byte {
	p := b.data[:n]
	b.data = b.data[n:]
	return p
}

func (b *buffer) remaining() []byte {
	return b.next(len(b.data))
}

func udp4pkt(packet []byte) []byte {
	ipLen := header.IPv4MinimumSize
	udpLen := header.UDPMinimumSize

	h := make([]byte, ipLen+udpLen)
	hdr := &buffer{h}

	ipv4fields := &header.IPv4Fields{
		IHL:         header.IPv4MinimumSize,
		TotalLength: uint16(ipLen + udpLen + len(packet)),
		TTL:         30,
		Protocol:    uint8(header.UDPProtocolNumber),
		DstAddr:     tcpip.Address(net.IPv4bcast.To4()),
	}
	ipv4hdr := header.IPv4(hdr.next(ipLen))
	ipv4hdr.Encode(ipv4fields)
	ipv4hdr.SetChecksum(^ipv4hdr.CalculateChecksum())

	udphdr := header.UDP(hdr.next(udpLen))
	udphdr.Encode(&header.UDPFields{
		SrcPort: srcPort,
		DstPort: dstPort,
		Length:  uint16(udpLen + len(packet)),
	})

	xsum := header.Checksum(packet, header.PseudoHeaderChecksum(
		ipv4hdr.TransportProtocol(), ipv4fields.SrcAddr, ipv4fields.DstAddr))
	udphdr.SetChecksum(^udphdr.CalculateChecksum(xsum, udphdr.Length()))

	return append(h, packet...)
}

func (pc *packetSock) SetReadTimeout(t time.Duration) error {
	tv := unix.NsecToTimeval(t.Nanoseconds())
	return unix.SetsockoptTimeval(pc.fd, unix.SOL_SOCKET, unix.SO_RCVTIMEO, &tv)
}

func swap16(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}
