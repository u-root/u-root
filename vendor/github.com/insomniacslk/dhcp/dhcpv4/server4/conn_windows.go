package server4

import (
	"errors"
	"net"
)

// NewIPv4UDPConn fails on Windows. Use WithConn() to pass the connection.
func NewIPv4UDPConn(iface string, addr *net.UDPAddr) (*net.UDPConn, error) {
	return nil, errors.New("not implemented on Windows")
}
