package dhcp4client

import (
	"net"
	"os"

	"golang.org/x/sys/unix"
)

// NewIPv4UDPConn returns a UDP connection bound to both the interface and port
// given. The UDP connection allows broadcasting.
func NewIPv4UDPConn(iface string, port int) (net.PacketConn, error) {
	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, unix.IPPROTO_UDP)
	if err != nil {
		return nil, err
	}
	f := os.NewFile(uintptr(fd), "")
	defer f.Close()

	// Allow broadcasting.
	if err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_BROADCAST, 1); err != nil {
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

	// This dups the FD. We still have to close f.
	return net.FilePacketConn(f)
}
