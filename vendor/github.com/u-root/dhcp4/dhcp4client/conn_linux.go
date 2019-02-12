// Copyright 2018 the u-root Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dhcp4client

import (
	"fmt"
	"net"
	"os"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
	"github.com/u-root/u-root/pkg/uio"
	"golang.org/x/sys/unix"
)

var (
	BroadcastMac = net.HardwareAddr([]byte{255, 255, 255, 255, 255, 255})
)

// NewIPv4UDPConn returns a UDP connection bound to both the interface and port
// given based on a IPv4 DGRAM socket. The UDP connection allows broadcasting.
func NewIPv4UDPConn(iface string, port int) (net.PacketConn, error) {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, unix.IPPROTO_UDP)
	if err != nil {
		return nil, err
	}
	f := os.NewFile(uintptr(fd), "")
	// net.FilePacketConn dups the FD, so we have to close this in any case.
	defer f.Close()

	// Allow broadcasting.
	if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_BROADCAST, 1); err != nil {
		return nil, err
	}
	// Allow reusing the addr to aid debugging.
	if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
		return nil, err
	}
	// Bind directly to the interface.
	if err := unix.BindToDevice(fd, iface); err != nil {
		return nil, err
	}
	// Bind to the port.
	if err := unix.Bind(fd, &unix.SockaddrInet4{Port: port}); err != nil {
		return nil, err
	}

	return net.FilePacketConn(f)
}

// NewPacketUDPConn returns a UDP connection bound to the interface and port
// given based on a raw packet socket. All packets are broadcasted.
func NewPacketUDPConn(iface string, port int) (net.PacketConn, error) {
	ifc, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}
	rawConn, err := raw.ListenPacket(ifc, uint16(ethernet.EtherTypeIPv4), &raw.Config{LinuxSockDGRAM: true})
	if err != nil {
		return nil, err
	}
	return NewBroadcastUDPConn(rawConn, &net.UDPAddr{Port: port}), nil
}

// UDPPacketConn implements net.PacketConn and marshals and unmarshals UDP
// packets.
type UDPPacketConn struct {
	net.PacketConn

	// boundAddr is the address this UDPPacketConn is "bound" to.
	//
	// Calls to ReadFrom will only return packets destined to this address.
	boundAddr *net.UDPAddr
}

// NewBroadcastUDPConn returns a PacketConn that marshals and unmarshals UDP
// packets, sending them to the broadcast MAC at on rawPacketConn.
//
// Calls to ReadFrom will only return packets destined to boundAddr.
func NewBroadcastUDPConn(rawPacketConn net.PacketConn, boundAddr *net.UDPAddr) net.PacketConn {
	return &UDPPacketConn{
		PacketConn: rawPacketConn,
		boundAddr:  boundAddr,
	}
}

func udpMatch(addr *net.UDPAddr, bound *net.UDPAddr) bool {
	if bound == nil {
		return true
	}
	if bound.IP != nil && !bound.IP.Equal(addr.IP) {
		return false
	}
	return bound.Port == addr.Port
}

// ReadFrom implements net.PacketConn.ReadFrom.
//
// ReadFrom reads raw IP packets and will try to match them against
// upc.boundAddr. Any matching packets are returned via the given buffer.
func (upc *UDPPacketConn) ReadFrom(b []byte) (int, net.Addr, error) {
	ipLen := IPv4MaximumHeaderSize
	udpLen := UDPMinimumSize

	for {
		pkt := make([]byte, ipLen+udpLen+len(b))
		n, _, err := upc.PacketConn.ReadFrom(pkt)
		if err != nil {
			return 0, nil, err
		}
		pkt = pkt[:n]
		buf := uio.NewBigEndianBuffer(pkt)

		// To read the header length, access data directly.
		ipHdr := IPv4(buf.Data())
		ipHdr = IPv4(buf.Consume(int(ipHdr.HeaderLength())))

		if ipHdr.TransportProtocol() != UDPProtocolNumber {
			continue
		}
		udpHdr := UDP(buf.Consume(udpLen))

		addr := &net.UDPAddr{
			IP:   net.IP(ipHdr.DestinationAddress()),
			Port: int(udpHdr.DestinationPort()),
		}
		if !udpMatch(addr, upc.boundAddr) {
			continue
		}
		return copy(b, buf.ReadAll()), addr, nil
	}
}

// WriteTo implements net.PacketConn.WriteTo and broadcasts all packets at the
// raw socket level.
//
// WriteTo wraps the given packet in the appropriate UDP and IP header before
// sending it on the packet conn.
func (upc *UDPPacketConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	udpAddr, ok := addr.(*net.UDPAddr)
	if !ok {
		return 0, fmt.Errorf("must supply UDPAddr")
	}

	// Using the boundAddr is not quite right here, but it works.
	packet := udp4pkt(b, udpAddr, upc.boundAddr)
	return upc.PacketConn.WriteTo(packet, &raw.Addr{HardwareAddr: BroadcastMac})
}
