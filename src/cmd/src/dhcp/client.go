package main

import (
	dhcp "dhcp4"
	"fmt"
	"log"
	"net"
	"time"
)

func client() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Can't enumerate interfaces? %v", err)
		return
	}
	addr, _, err := net.ParseCIDR("0.0.0.0/32")
	if err != nil {
		log.Printf("Can't parse to ip: %v", err)
		return
	}

	p := dhcp.RequestPacket(dhcp.Discover, ifaces[0].HardwareAddr, addr, []byte{1, 2, 3}, true, nil)
	fmt.Printf("%q\n", p)

	d, err := net.Dial("udp", "127.0.0.1:67")
	if err != nil {
		log.Printf("dial bad %v", err)
	}
	fmt.Printf("d is %q\n", d)
	ra, err := net.ResolveUDPAddr("udp", "127.0.0.1:67")
	if err != nil {
		log.Printf("ResolveUDPAddr failed: %v", err)
		return
	}

	if err = d.SetDeadline(time.Second); err != nil {
		log.Printf("Can't set deadline: %v\n", err)
		return
	}
	fmt.Printf("ra %v\n", ra)
	for {
		_, err = d.Write(p)
		if err != nil {
			log.Printf("WriteToUDP failed: %v", err)
			return
		}
	}
}
