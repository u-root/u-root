package main

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

type packetSock struct {
	fd      int
	ifindex int
}

var bcastMAC = []byte{0x33, 0x33, 0x00, 0x01, 0x00, 0x02}

func ifup(ifname string) (netlink.Link, error) {
	iface, err := netlink.LinkByName(ifname)
	if err != nil {
		return nil, fmt.Errorf("%s: netlink.LinkByName failed: %v", ifname, err)
	}
	if err := netlink.LinkSetUp(iface); err != nil {
		return nil, fmt.Errorf("%v: %v can't make it up: %v", ifname, iface, err)
	}
	return iface, nil
}

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

func (pc *packetSock) write(pb []byte, mac net.HardwareAddr) error {
	// Define linke layer
	lladdr := unix.SockaddrLinklayer{
		Ifindex:  pc.ifindex,
		Protocol: swap16(unix.ETH_P_IPV6),
		Halen:    uint8(len(mac)),
	}
	copy(lladdr.Addr[:], mac)

	pkt := []byte{0x01}
	return unix.Sendto(pc.fd, pkt, 0, &lladdr)
}

func (pc *packetSock) ReadFrom() ([]byte, error) {
	pb := make([]byte, 500)
	n, _, err := unix.Recvfrom(pc.fd, pb, 0)
	if err != nil {
		return nil, err
	}
	return pb[:n], nil
}

func swap16(x uint16) uint16 {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], x)
	return binary.LittleEndian.Uint16(b[:])
}
