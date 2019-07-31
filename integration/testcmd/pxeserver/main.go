// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net"
	"sync"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"pack.ag/tftp"
)

var (
	selfIP    = flag.String("ip", "192.168.0.1", "IP of self")
	yourIP    = flag.String("your-ip", "192.168.0.2", "The one and only IP to give to all clients")
	directory = flag.String("dir", "", "Directory to serve")
)

func dhcpHandler(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	self := net.ParseIP(*selfIP)
	you := net.ParseIP(*yourIP)
	log.Printf("Handling request %v", m)

	var replyType dhcpv4.MessageType
	switch mt := m.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		replyType = dhcpv4.MessageTypeOffer
	case dhcpv4.MessageTypeRequest:
		replyType = dhcpv4.MessageTypeAck
	default:
		log.Printf("Can't handle type %v", mt)
		return
	}
	reply, err := dhcpv4.NewReplyFromRequest(m,
		dhcpv4.WithMessageType(replyType),
		dhcpv4.WithServerIP(self),
		dhcpv4.WithRouter(self),
		dhcpv4.WithNetmask(net.CIDRMask(24, 32)),
		dhcpv4.WithYourIP(you),
	)
	reply.BootFileName = "pxelinux.0"
	if err != nil {
		log.Printf("Could not create reply for %v: %v", m, err)
		return
	}
	if _, err := conn.WriteTo(reply.ToBytes(), &net.UDPAddr{IP: net.IP{255, 255, 255, 255}, Port: 68}); err != nil {
		log.Printf("Could not write %v: %v", reply, err)
	}
}

func main() {
	flag.Parse()

	var wg sync.WaitGroup
	if len(*directory) != 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			server, err := tftp.NewServer(":69")
			if err != nil {
				log.Fatalf("Could not start TFTP server: %v", err)
			}

			log.Println("starting file server")
			server.ReadHandler(tftp.FileServer(*directory))
			log.Fatal(server.ListenAndServe())
		}()
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		laddr := &net.UDPAddr{Port: 67}
		server, err := server4.NewServer(laddr, dhcpHandler)
		if err != nil {
			log.Fatal(err)
		}
		server.Serve()
	}()

	wg.Wait()
}
