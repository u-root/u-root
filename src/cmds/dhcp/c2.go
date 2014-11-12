package main

import (
	dhcp "dhcp4"
	"fmt"
	"log"
	"net"
	"netlink"
	"os"
	"time"
)

type dhcpInfo struct {
	i *net.Interface
	dhcp.Packet
}

func c2() {
	fails := 0
	r := make(chan *dhcpInfo)
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("client: Can't enumerate interfaces? %v", err)
		return
	}
	for _, v := range ifaces {
		go one(&v, r)
	}
	for p := range r {
		if p == nil {
			fails++
			if fails > 10 {
				os.Exit(1)
			}
			continue
		}
		fmt.Printf("GOT ONE! %v\n", p)
		if p.OpCode() != dhcp.BootReply {
			fmt.Printf("not a reply?\n")
			continue
		}
		options := p.ParseOptions()
		switch dhcp.MessageType(options[dhcp.OptionDHCPMessageType][0]) {
		case dhcp.Offer:
			fmt.Printf("reply. flags %v HLen %v XId %v CI %v YI %v SI %v GI %v CH %v\n",
				p.Flags(), p.HLen(), p.XId(),
				p.CIAddr(),
				p.YIAddr(),
				p.SIAddr(),
				p.GIAddr(),
				p.CHAddr())
			addr := p.YIAddr()

			netmask := options[dhcp.OptionSubnetMask]
			if netmask != nil {
				fmt.Printf("OptionSubnetMask is %v\n", netmask)
			} else {
				// what do to?
				netmask = addr
			}
			// they better be the same len. I'm happy explode if not.
			network := addr
			for i := range addr {
				network[i] = addr[i] & netmask[i]
			}
			if false {
				netlink.NetworkLinkAddIp(p.i, addr, &net.IPNet{network, netmask})
			}
			gwData := options[dhcp.OptionRouter]
			if gwData != nil {
				fmt.Printf("router %v\n", gwData)
			}
			if err := netlink.AddRouteIP(p.i, []byte{}, []byte{}, gwData); err != nil {
				fmt.Printf("Can't add route: %v\n", err)
			}

		default:
			fmt.Printf("not what we hoped: %v\n", dhcp.MessageType(p.HType()))
		}
	}
}

func one(i *net.Interface, r chan *dhcpInfo) {
     // the link has to be uppable
	if err := netlink.NetworkLinkUp(i); err != nil {
		log.Printf("%v can't make it up: %v", i, err)
		return
	}

	addr, _, err := net.ParseCIDR("255.255.255.255/32")
	if err != nil {
		log.Printf("client: Can't parse to ip: %v", err)
		r <- nil
		return
	}
	p := dhcp.RequestPacket(dhcp.Discover, i.HardwareAddr, addr, []byte{1, 2, 3}, true, nil)
	fmt.Printf("client: %q\n", p)

	d, err := net.ListenPacket("udp", "")
	if err != nil {
		fmt.Printf("listen packet: %v\n", err)
		r <- nil
		return
	}
	defer d.Close()
	for {
		fmt.Printf("Try it\n")
		fmt.Printf("client: d is %q\n", d)
		//ra, err := net.ResolveUDPAddr("udp", "127.0.0.1:67")
		ra, err := net.ResolveUDPAddr("udp", "255.255.255.255:67")
		if err != nil {
			log.Printf("client: ResolveUDPAddr failed: %v", err)
		r <- nil
			return
		}

		fmt.Printf("client: ra %v\n", ra)
		if err := d.SetDeadline(time.Now().Add(10000 * time.Millisecond)); err != nil {
			log.Printf("client: Can't set deadline: %v\n", err)
		r <- nil
			return
		}
		if _, err := d.WriteTo(p, ra); err != nil {
			log.Printf("client: WriteToUDP failed: %v", err)
		r <- nil
			return
		} else {
			b := [512]byte{}
			if err := d.SetReadDeadline(time.Now().Add(10000 * time.Millisecond)); err != nil {
				log.Printf("client: Can't set deadline: %v\n", err)
		r <- nil
				return
			}
			fmt.Printf("Client: sleep the read\n")
			time.Sleep(time.Second)
			if n, a, err := d.ReadFrom(b[:]); err != nil {
				log.Printf("client: Read  from UDP failed: %v", err)
				continue
			} else {
				fmt.Printf("client: Data %v amt %v a %v\n", b, n, a)
				r <- &dhcpInfo{i, dhcp.Packet(b[:])}
			}
		}
	}
}
