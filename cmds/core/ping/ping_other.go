// +build !linux

package main

import (
	"errors"
	"net"
)

func setupICMPv6Socket(c *net.IPConn) error {
	return errors.New("setting up ICMPv6 socket only supported on Linux")
}
