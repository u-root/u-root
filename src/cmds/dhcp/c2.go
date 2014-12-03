package main

// what we've learned.
// must send packets on link layer
// Our packet assembly is wrong.

import (
	dhcp "dhcp4"
	"fmt"
	"log"
	"net"
	"netlink"
	"os"
	"regexp"
	"syscall"
	"time"
)

type dhcpInfo struct {
	i *net.Interface
	dhcp.Packet
}

func c2(re *regexp.Regexp) {
	fails := 0
	r := make(chan *dhcpInfo)
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("client: Can't enumerate interfaces? %v", err)
		return
	}
	for _, v := range ifaces {
		if !re.Match([]byte(v.Name)) {
			continue
		}
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
	// possibly bogus packet created. I think they are not creating an IP header.
	p := dhcp.RequestPacket(dhcp.Discover, i.HardwareAddr, addr, []byte{1, 2, 3}, true, nil)
	fmt.Printf("client: %q\n", p)
	u := &IPUDPHeader {
	Version: 4,
	DPort: 67,
	}
	raw := u.Marshal(p)
	s, err := syscall.LsfSocket(i.Index, syscall.ETH_P_IP)
	if err != nil {
		fmt.Printf("lsfsocket: got %v\n", err)
		r <- nil
		return
	}

	// we don't set family; Sendto does.
	bcast := &syscall.SockaddrLinklayer{
		Protocol: syscall.ETH_P_IP,
		Ifindex:  i.Index,
		Halen:    6,
		Addr:     [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}
	for tries := 0; tries < 10; tries++ {
		fmt.Printf("Try it\n")
		err = syscall.Sendto(s, raw, 0, bcast)
		//err = pc.WriteTo(p, nil, addr)
		if err != nil {
			log.Printf("client: WriteToUDP failed: %v", err)
			r <- nil
			return
		}
		log.Printf("wrote it; get it")
		fmt.Printf("Client: sleep the read\n")
		time.Sleep(time.Second)

		/*
			b := [512]byte{}
			n, err := syscall.Read(s, b[:])
			if err != nil {
					log.Printf("client: %v\n", err)
							r <- nil
					return
				}
					fmt.Printf("client: Data %v amt %v \n", b, n)
					r <- &dhcpInfo{i, dhcp.Packet(b[:])}
		*/

	}
}
