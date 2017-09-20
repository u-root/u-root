package dhcp6

import (
	"net"
	"syscall"
	"time"
)

// TODO: Make packetSock implement net.PacketConn?
type packetSock struct {
	fd      int
	ifindex int
}

// NewPacketSock creates a new socket that sends and receives packets.
func NewPacketSock(ifindex int) (*packetSock, error) {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_DGRAM, int(swap16(syscall.ETH_P_IPV6)))
	if err != nil {
		return nil, err
	}
	addr := syscall.SockaddrLinklayer{
		Ifindex:  ifindex,
		Protocol: swap16(syscall.ETH_P_IPV6),
	}
	if err := syscall.Bind(fd, &addr); err != nil {
		return nil, err
	}
	return &packetSock{
		fd:      fd,
		ifindex: ifindex,
	}, nil
}

// Write a packet.
func (pc packetSock) WriteTo(p []byte, mac net.HardwareAddr) error {
	lladdr := syscall.SockaddrLinklayer{
		Ifindex:  pc.ifindex,
		Protocol: swap16(syscall.ETH_P_IPV6),
		Halen:    uint8(len(mac)),
	}
	copy(lladdr.Addr[:], mac)

	// Send out request from link layer
	return syscall.Sendto(pc.fd, p, 0, &lladdr)
}

// Read reads packets from the socket.
func (pc packetSock) ReadFrom(p []byte) (int, error) {
	n, _, err := syscall.Recvfrom(pc.fd, p, 0)
	return n, err
}

// Close socket.
func (pc packetSock) Close() error {
	return syscall.Close(pc.fd)
}

// Set a read timeout
func (pc *packetSock) SetReadTimeout(t time.Duration) error {
	tv := syscall.NsecToTimeval(t.Nanoseconds())
	return syscall.SetsockoptTimeval(pc.fd, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv)
}

func swap16(x uint16) uint16 {
	return (x<<8)&0xff00 | (x >> 8)
}
