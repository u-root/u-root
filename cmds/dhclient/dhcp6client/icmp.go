package dhcp6client

import (
	"golang.org/x/net/ipv6"
)

type icmpPack struct {
	icmpType ipv6.ICMPType
	options  []byte
}
