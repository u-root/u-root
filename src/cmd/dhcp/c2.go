package main

import (
	dhcp "dhcp4"
	"fmt"
	"log"
	"net"
	"time"
)

func c2() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("client: Can't enumerate interfaces? %v", err)
		return
	}
	addr, _, err := net.ParseCIDR("0.0.0.0/32")
	if err != nil {
		log.Printf("client: Can't parse to ip: %v", err)
		return
	}

	p := dhcp.RequestPacket(dhcp.Discover, ifaces[0].HardwareAddr, addr, []byte{1, 2, 3}, true, nil)
	fmt.Printf("client: %q\n", p)

	d, err := net.ListenPacket("udp", "")
	if err != nil {
		fmt.Printf("listen packet: %v\n", err)
		return
	}
	defer d.Close()
	for {
		fmt.Printf("Try it\n")
		fmt.Printf("client: d is %q\n", d)
		ra, err := net.ResolveUDPAddr("udp", "127.0.0.1:67")
		if err != nil {
			log.Printf("client: ResolveUDPAddr failed: %v", err)
			return
		}

		fmt.Printf("client: ra %v\n", ra)
		if err := d.SetDeadline(time.Now().Add(10000 * time.Millisecond)); err != nil {
			log.Printf("client: Can't set deadline: %v\n", err)
			return
		}
		if _, err := d.WriteTo(p, ra); err != nil {
			log.Printf("client: WriteToUDP failed: %v", err)
			return
		} else {
			b := [512]byte{}
			if err := d.SetReadDeadline(time.Now().Add(10000 * time.Millisecond)); err != nil {
				log.Printf("client: Can't set deadline: %v\n", err)
				return
			}
			fmt.Printf("Client: sleep the read\n")
			time.Sleep(time.Second)
			if n, a, err := d.ReadFrom(b[:]); err != nil {
				log.Printf("client: Read  from UDP failed: %v", err)
				continue
			} else {
				fmt.Printf("client: Data %v amt %v a %v\n", b, n, a)
				return
			}
		}
	}
}
