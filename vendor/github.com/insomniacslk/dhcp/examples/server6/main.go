package main

import (
	"log"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/dhcpv6/server6"
)

func handler(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6) {
	// this function will just print the received DHCPv6 message, without replying
	log.Print(m.Summary())
}

func main() {
	laddr := &net.UDPAddr{
		IP:   net.ParseIP("::1"),
		Port: dhcpv6.DefaultServerPort,
	}
	server, err := server6.NewServer("", laddr, handler)
	if err != nil {
		log.Fatal(err)
	}

	server.Serve()
}
