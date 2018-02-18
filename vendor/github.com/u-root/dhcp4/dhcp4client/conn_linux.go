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
