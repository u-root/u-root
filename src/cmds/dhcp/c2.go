package main
/* sample good packet.
29:26:3e:37 (oui Unknown), length 300, xid 0x1b83823b, Flags [none] (0x0000)
	  Client-Ethernet-Address 00:0c:29:26:3e:37 (oui Unknown)
	  Vendor-rfc1048 Extensions
	    Magic Cookie 0x63825363
	    DHCP-Message Option 53, length 1: Request
	    Requested-IP Option 50, length 4: 192.168.28.203
	    Hostname Option 12, length 24: "rminnich-virtual-machine"
	    Parameter-Request Option 55, length 13: 
	      Subnet-Mask, BR, Time-Zone, Default-Gateway
	      Domain-Name, Domain-Name-Server, Option 119, Hostname
	      Netbios-Name-Server, Netbios-Scope, MTU, Classless-Static-Route
	      NTP
	    END Option 255, length 0
	    PAD Option 0, length 0, occurs 9
	0x0000:  4510 0148 0000 0000 8011 3996 0000 0000
	0x0010:  ffff ffff 0044 0043 0134 7422 0101 0600
	0x0020:  1b83 823b 0000 0000 0000 0000 0000 0000
	0x0030:  0000 0000 0000 0000 000c 2926 3e37 0000
	0x0040:  0000 0000 0000 0000 0000 0000 0000 0000
	0x0050:  0000 0000 0000 0000 0000 0000 0000 0000
	0x0060:  0000 0000 0000 0000 0000 0000 0000 0000
	0x0070:  0000 0000 0000 0000 0000 0000 0000 0000
	0x0080:  0000 0000 0000 0000 0000 0000 0000 0000
	0x0090:  0000 0000 0000 0000 0000 0000 0000 0000
	0x00a0:  0000 0000 0000 0000 0000 0000 0000 0000
	0x00b0:  0000 0000 0000 0000 0000 0000 0000 0000
	0x00c0:  0000 0000 0000 0000 0000 0000 0000 0000
	0x00d0:  0000 0000 0000 0000 0000 0000 0000 0000
	0x00e0:  0000 0000 0000 0000 0000 0000 0000 0000
	0x00f0:  0000 0000 0000 0000 0000 0000 0000 0000
	0x0100:  0000 0000 0000 0000 6382 5363 3501 0332
	0x0110:  04c0 a81c cb0c 1872 6d69 6e6e 6963 682d
	0x0120:  7669 7274 7561 6c2d 6d61 6368 696e 6537
	0x0130:  0d01 1c02 030f 0677 0c2c 2f1a 792a ff00
	0x0140:  0000 0000 0000 0000


// BAD
11:36:18.213690 IP (tos 0x0, ttl 64, id 0, offset 0, flags [none], proto UDP (17), length 300, bad cksum 0 (->69c3)!)
    0.0.0.0.0 > 15.255.255.255.bootps: [no cksum] BOOTP/DHCP, Request from 00:0c:29:26:3e:37 (oui Unknown), length 272, xid 0x1020300, Flags [Broadcast] (0x8000)
	  Client-IP 255.255.255.255
	  Client-Ethernet-Address 00:0c:29:26:3e:37 (oui Unknown)
	  Vendor-rfc1048 Extensions
	    Magic Cookie 0x63825363
	    DHCP-Message Option 53, length 1: Discover
	    END Option 255, length 0
	    PAD Option 0, length 0, occurs 14
	0x0000:  4500 012c 0000 0000 4011 0000 0000 0000
	0x0010:  0fff ffff 0000 0043 012c 0000 0101 0600
	0x0020:  0102 0300 0000 8000 ffff ffff 0000 0000
	0x0030:  0000 0000 0000 0000 000c 2926 3e37 0000
	0x0040:  0000 0000 0000 0000 0000 0000 0000 0000
	0x0050:  0000 0000 0000 0000 0000 0000 0000 0000
	0x0060:  0000 0000 0000 0000 0000 0000 0000 0000
	0x0070:  0000 0000 0000 0000 0000 0000 0000 0000
	0x0080:  0000 0000 0000 0000 0000 0000 0000 0000
	0x0090:  0000 0000 0000 0000 0000 0000 0000 0000
	0x00a0:  0000 0000 0000 0000 0000 0000 0000 0000
	0x00b0:  0000 0000 0000 0000 0000 0000 0000 0000
	0x00c0:  0000 0000 0000 0000 0000 0000 0000 0000
	0x00d0:  0000 0000 0000 0000 0000 0000 0000 0000
	0x00e0:  0000 0000 0000 0000 0000 0000 0000 0000
	0x00f0:  0000 0000 0000 0000 0000 0000 0000 0000
	0x0100:  0000 0000 0000 0000 6382 5363 3501 01ff
	0x0110:  0000 0000 0000 0000 0000 0000 0000

strace of dhclient shows a write of this packet:

9955  socket(PF_PACKET, SOCK_RAW, 768)  = 5
9955  ioctl(5, SIOCGIFINDEX, {ifr_name="eth0", ifr_index=2}) = 0
9955  bind(5, {sa_family=AF_PACKET, proto=0x03, if2, pkttype=PACKET_HOST, addr(0)={0, }, 20) = 0
9955  setsockopt(5, SOL_PACKET, PACKET_AUXDATA, [1], 4) = 0
9955  setsockopt(5, SOL_SOCKET, SO_ATTACH_FILTER, "\v\0\0\0\0\0\0\0\200\363\3761,\177\0\0", 16) = 0
9955  sendto(3, "<30>Dec  2 17:52:12 dhclient: Listening on LPF/eth0/00:0c:29:26:3e:37", 69, MSG_NOSIGNAL, NULL, 0) = 69
9955  write(2, "Listening on LPF/eth0/00:0c:29:26:3e:37", 39) = 39
9955  write(2, "\n", 1)                 = 1
9955  sendto(3, "<30>Dec  2 17:52:12 dhclient: Sending on   LPF/eth0/00:0c:29:26:3e:37", 69, MSG_NOSIGNAL, NULL, 0) = 69
9955  write(2, "Sending on   LPF/eth0/00:0c:29:26:3e:37", 39) = 39
9955  write(2, "\n", 1)                 = 1
9955  fcntl(5, F_SETFD, FD_CLOEXEC)     = 0
9955  socket(PF_INET, SOCK_DGRAM, IPPROTO_UDP) = 6
9955  setsockopt(6, SOL_SOCKET, SO_REUSEADDR, [1], 4) = 0
9955  bind(6, {sa_family=AF_INET, sin_port=htons(68), sin_addr=inet_addr("0.0.0.0")}, 16) = 0
9955  sendto(3, "<30>Dec  2 17:52:12 dhclient: Sending on   Socket/fallback", 58, MSG_NOSIGNAL, NULL, 0) = 58
9955  write(2, "Sending on   Socket/fallback", 28) = 28
9955  write(2, "\n", 1)                 = 1
9955  fcntl(6, F_SETFD, FD_CLOEXEC)     = 0
9955  uname({sys="Linux", node="rminnich-virtual-machine", ...}) = 0
9955  sendto(3, "<30>Dec  2 17:52:12 dhclient: DHCPREQUEST of 192.168.28.203 on eth0 to 255.255.255.255 port 67 (xid=0x3b82831b)", 111, MSG_NOSIGNAL, NULL, 0) = 111
9955  write(2, "DHCPREQUEST of 192.168.28.203 on eth0 to 255.255.255.255 port 67 (xid=0x3b82831b)", 81) = 81
9955  write(2, "\n", 1)                 = 1

H = 0x48 = 3 words of options (HL is 8 words).
9955  write(5, "\377\377\377\377\377\377 \0\f)&>7 \10\0E \20\1 H \0\0\0\0\200\0219\226\0\0\0\0\377\377\377\377 (options) \0D\0C\0014t\"\1\1\6\0\33\203\202;\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\f)&>7\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0c\202Sc5\1\0032\4\300\250\34\313\f\30rminnich-virtual-machine7\r\1\34\2\3\17\6w\f,/\32y*\377\0\0\0\0\0\0\0\0\0", 342) = 342


9955  write(5, "\377\377\377\377\377\377\0\f)&>7\10\0E\20\1H\0\0\0\0\200\0219\226\0\0\0\0\377\377\377\377\0D\0C\0014t\"\1\1\6\0\33\203\202;\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\f)&>7\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0c\202Sc5\1\0032\4\300\250\34\313\f\30rminnich-virtual-machine7\r\1\34\2\3\17\6w\f,/\32y*\377\0\0\0\0\0\0\0\0\0", 342) = 342

hmm.

dest is ffffff, and we don't have that. do we need to add the enet header too?


Known good sequence from udhcpc

14444 sendto(3, "<15>Dec  3 09:19:10 udhcpc: Send"..., 47, MSG_NOSIGNAL, NULL, 0) = 47
14444 socket(PF_PACKET, SOCK_DGRAM, 8)  = 6
14444 bind(6, {sa_family=AF_PACKET, proto=0x800, if2, pkttype=PACKET_HOST, addr(6)={0, ffffffffffff}, 20) = 0
14444 sendto(6, "E\0\2@\0\0\0\0@\21x\256\0\0\0\0\377\377\377\377\0D\0C\2,]z\1\1\6\0"..., 576, 0, {sa_family=AF_PACKET, proto=0x800, if2, pkttype=PACKET_HOST, addr(6)={0, ffffffffffff}, 20) = 576

Note the enet header is not filled in, but space is left for it.
*/

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
	"unsafe"
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
		log.Printf("Let's check %v", v)
		if !re.Match([]byte(v.Name)) {
			continue
		}
		log.Printf("Let's USE iface  %v", v)
		go one(v, r)
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

func one(i net.Interface, r chan *dhcpInfo) {
	// the link has to be uppable
	if err := netlink.NetworkLinkUp(&i); err != nil {
		log.Printf("%v can't make it up: %v", i, err)
		return
	}

	addr, _, err := net.ParseCIDR("0.0.0.0/32")
	if err != nil {
		log.Printf("client: Can't parse to ip: %v", err)
		r <- nil
		return
	}
	// possibly bogus packet created. I think they are not creating an IP header.
	p := dhcp.RequestPacket(dhcp.Discover, i.HardwareAddr, addr, []byte{1, 2, 3}, false, nil)
	fmt.Printf("client: len %d\n", len(p))
	u := &EtherIPUDPHeader {
	Version: 4,
	IHL: 5,
	DPort: 67,
	SPort: 68,
	TotalLength: uint16(len(p)),
	Length:uint16(len(p)),
	DIP: 0xffffffff,
	Protocol: syscall.IPPROTO_UDP,
	TTL: 64,
	}
	raw := u.Marshal(p)
/* goddamn. if only this had worked.
	s, err := syscall.LsfSocket(i.Index, syscall.ETH_P_IP)
 */
	// yegads, the socket interface sucks so hard for over 30 years now ...
	// htons for a LOCAL RESOURCE? Riiiiiight.
	// How I miss Plan 9
	s, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_DGRAM, 0x8) //syscall.ETH_P_IP)
	if err != nil {
		fmt.Printf("lsfsocket: got %v\n", err)
		r <- nil
		return
	}
	var lsall syscall.SockaddrLinklayer
	pp := (*[2]byte)(unsafe.Pointer(&lsall.Protocol))
	pp[0] = 8
	pp[1] = 0
	lsall.Ifindex = i.Index
	lsall.Halen = 6
	lsall.Addr = [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	if err = syscall.Bind(s, &lsall); err != nil {
		fmt.Printf("lsfsocket: bind got %v\n", err)
		r <- nil
		return
	}

	// we don't set family; Sendto does.
	bcast := &syscall.SockaddrLinklayer{
		Protocol: 0x8, //syscall.ETH_P_IP,
		Ifindex:  i.Index,
		Halen:    6,
		Addr:     [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}
	log.Printf("bcast is %v", bcast)
	for tries := 0; tries < 1; tries++ {
		fmt.Printf("Try it\n")
		err = syscall.Sendto(s, raw, 0, bcast)
		//err = pc.WriteTo(p, nil, addr)
		//n, err := syscall.Write(s, raw)
		if err != nil {
			log.Printf("client: WriteToUDP failed: %v", err)
			r <- nil
			return
		}
		//log.Printf("wrote it; %v bytes", n)
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
