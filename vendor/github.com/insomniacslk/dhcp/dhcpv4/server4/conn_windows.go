//go:build windows

package server4

import (
	"fmt"
	"net"
)

// NewIPv4UDPConn returns an UDPv4 connection bound to the IP and port provider
func NewIPv4UDPConn(iface string, addr *net.UDPAddr) (*net.UDPConn, error) {
	connection, err := net.ListenPacket("udp4", addr.String())
	if err != nil {
		return nil, fmt.Errorf("We cannot listen on %s and port %d: %v", addr.IP, addr.Port, err)
	}
	udpConn, ok := connection.(*net.UDPConn)
	if !ok {
		return nil, fmt.Errorf("The connection is not of the proper type")
	}
	return udpConn, nil
}
